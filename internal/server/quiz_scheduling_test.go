package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/quiz"
	"github.com/chuma-beep/mathua/internal/storage"
)

// The routine mastery check is scheduled by the existing XP gate. These tests
// pin the edge cases of the enforcement introduced in handleQuizSession: below,
// at and above the threshold, the snapshot moving only on completion, duplicate
// and concurrent starts joining one live run, and the modes that must not be
// gated (the anonymous quiz, diagnostics and placement).
//
// SQLite always runs. TEST_POSTGRES_DSN also runs every case on Postgres, the
// engine production uses: the gate itself is arithmetic over two repository
// reads, and this is the only way to know both stores answer them identically.
func quizSchedulingStores(t *testing.T) map[string]storage.Repository {
	t.Helper()
	dir := t.TempDir()
	sq, err := storage.NewSQLiteStore(filepath.Join(dir, "quiz-scheduling.db"))
	if err != nil {
		t.Fatalf("create sqlite store: %v", err)
	}
	t.Cleanup(func() { sq.Close() })
	out := map[string]storage.Repository{"sqlite": sq}
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		pg, err := storage.NewPostgresStore(dsn)
		if err != nil {
			t.Fatalf("open postgres: %v", err)
		}
		t.Cleanup(func() { pg.Close() })
		out["postgres"] = pg
	} else {
		t.Log("TEST_POSTGRES_DSN not set; Postgres leg skipped")
	}
	return out
}

func quizSchedulingServer(t *testing.T, store storage.Repository) (*Server, *http.ServeMux) {
	t.Helper()
	return quizSchedulingServerWith(t, store, []concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
		{ID: "b", Label: "B", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
}

func quizSchedulingServerWith(t *testing.T, store storage.Repository, cs []concepts.Concept) (*Server, *http.ServeMux) {
	t.Helper()
	d, err := concepts.Build(cs)
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	reg := generator.NewRegistry()
	for _, c := range cs {
		reg.Register(c.ID, &testGen{})
	}
	s := New(engine.New(store, d, reg, nil, planning.New(d)), store, auth.New(store))
	// The tests POST faster than a learner can, so the write limiter would
	// answer 429 for reasons unrelated to scheduling.
	s.writeLimiter = newRateLimiter(100000, 100000, time.Second, "")
	mux := http.NewServeMux()
	s.Register(mux)
	return s, mux
}

func quizSchedulingPost(t *testing.T, mux *http.ServeMux, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", path, bytes.NewReader([]byte(body)))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	mux.ServeHTTP(rec, req)
	return rec
}

// quizSchedulingLearn marks a then b as learned concepts. They are
// independent, so both remain routine-check candidates even though mastered
// concepts are reserved for retention review.
func quizSchedulingLearn(t *testing.T, store storage.Repository, studentID string) {
	t.Helper()
	seedQuizLearned(t, store, studentID, "a", string(mastery.StatusPracticing), 0)
	seedQuizLearned(t, store, studentID, "b", string(mastery.StatusPracticing), 0)
}

// seedQuizUnseen writes the propagated-weakness shape: an UNSEEN row with a
// score but no learning behind it. It is not a quiz candidate.
func seedQuizUnseen(t *testing.T, store storage.Repository, studentID, conceptID string) {
	t.Helper()
	if err := store.UpsertProgress(&storage.ConceptProgress{
		StudentID:     studentID,
		ConceptID:     conceptID,
		Status:        string(mastery.StatusUnseen),
		WeaknessScore: 0.9,
	}); err != nil {
		t.Fatalf("seed unseen progress for %s: %v", conceptID, err)
	}
}

func quizSchedulingSeed(t *testing.T, store storage.Repository, studentID string, amount int, key string) {
	t.Helper()
	if err := store.AwardXP(studentID, amount, "test", key); err != nil {
		t.Fatalf("seed %d XP: %v", amount, err)
	}
}

func quizSchedulingDecode(t *testing.T, rec *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var out map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode %s: %v", rec.Body.String(), err)
	}
	return out
}

func quizSchedulingScores(t *testing.T, mux *http.ServeMux, token string) map[string]interface{} {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/scores/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("scores: %d %s", rec.Code, rec.Body.String())
	}
	return quizSchedulingDecode(t, rec)
}

// quizSchedulingRun starts a run and answers every question, returning the
// concepts in the order they were served. It is the HTTP-level equivalent of a
// learner completing a mastery check.
func quizSchedulingRun(t *testing.T, mux *http.ServeMux, token string) []string {
	t.Helper()
	rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("start run: %d %s", rec.Code, rec.Body.String())
	}
	res := quizSchedulingDecode(t, rec)
	if res["done"] == true {
		t.Fatalf("run started already done: %v", res)
	}
	sessionID, _ := res["session_id"].(string)
	conceptID, _ := res["concept_id"].(string)
	if sessionID == "" || conceptID == "" {
		t.Fatalf("start run missing ids: %v", res)
	}
	seen := []string{conceptID}
	for i := 0; i < 20; i++ {
		body := fmt.Sprintf(`{"session_id":%q,"concept_id":%q,"answer":"4","elapsed":5}`, sessionID, conceptID)
		rec = quizSchedulingPost(t, mux, "/api/quiz/answer", token, body)
		if rec.Code != http.StatusOK {
			t.Fatalf("answer %s: %d %s", conceptID, rec.Code, rec.Body.String())
		}
		res = quizSchedulingDecode(t, rec)
		if res["done"] == true {
			return seen
		}
		conceptID, _ = res["concept_id"].(string)
		if conceptID == "" {
			t.Fatalf("next question missing concept_id: %v", res)
		}
		seen = append(seen, conceptID)
	}
	t.Fatalf("run did not finish: %v", seen)
	return nil
}

func TestQuizSession_ConsecutiveQuizzesAvoidRecentConcepts(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			var cs []concepts.Concept
			for i := 0; i < 12; i++ {
				id := fmt.Sprintf("c%02d", i)
				cs = append(cs, concepts.Concept{
					ID: id, Label: id, Domain: "d", Subdomain: fmt.Sprintf("s%02d", i),
					GradingType: "numeric", Prerequisites: []string{},
					MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
				})
			}
			_, mux := quizSchedulingServerWith(t, store, cs)
			token, sid := quizGuest(t, mux)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "consecutive-gate-1")
			for _, c := range cs {
				seedQuizLearned(t, store, sid, c.ID, string(mastery.StatusPracticing), 0.5)
			}

			first := quizSchedulingRun(t, mux, token)
			if len(first) != quiz.QuizSize {
				t.Fatalf("first run served %d concepts, want %d: %v", len(first), quiz.QuizSize, first)
			}

			// Completion reset the XP snapshot, so earn the gate again before
			// the second check. The first run's concepts are now the recent
			// quiz window; seven candidates remain unused.
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "consecutive-gate-2")
			second := quizSchedulingRun(t, mux, token)
			if len(second) != quiz.QuizSize {
				t.Fatalf("second run served %d concepts, want %d: %v", len(second), quiz.QuizSize, second)
			}
			seen := make(map[string]bool, len(first))
			for _, id := range first {
				seen[id] = true
			}
			for _, id := range second {
				if seen[id] {
					t.Errorf("consecutive quiz repeated %s; first %v, second %v", id, first, second)
				}
			}
		})
	}
}

func TestQuizSession_EnforcesXPGate(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			for _, tc := range []struct {
				name string
				xp   int
				want int
			}{
				{"below threshold", engine.QuizGateXP - 1, http.StatusConflict},
				{"exactly at threshold", engine.QuizGateXP, http.StatusOK},
				{"above threshold", engine.QuizGateXP + 1, http.StatusOK},
			} {
				t.Run(tc.name, func(t *testing.T) {
					token, sid := quizGuest(t, mux)
					quizSchedulingSeed(t, store, sid, tc.xp, fmt.Sprintf("gate-%d", tc.xp))
					if tc.want == http.StatusOK {
						quizSchedulingLearn(t, store, sid)
					}

					rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
					if rec.Code != tc.want {
						t.Fatalf("quiz session at %d XP: got %d, want %d: %s", tc.xp, rec.Code, tc.want, rec.Body.String())
					}
					if tc.want == http.StatusConflict {
						res := quizSchedulingDecode(t, rec)
						if res["error"] != quizNotDueMessage {
							t.Errorf("expected error %q, got %v", quizNotDueMessage, res["error"])
						}
						if _, ok := res["session_id"]; ok {
							t.Errorf("a refused start must not create a session: %v", res)
						}
						if p, _ := store.GetProgress(sid, "a"); p != nil {
							t.Errorf("a refused start wrote progress: %+v", p)
						}
					}
				})
			}
		})
	}
}

func TestQuizSession_GateIsRelativeToLastCompletedQuiz(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "first-gate")

			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			created := quizSchedulingScores(t, mux, token)
			if created["xp_since_quiz"] != float64(engine.QuizGateXP) || created["quiz_due"] != true {
				t.Errorf("creating a session consumed its admission: %v", created)
			}
			sidSession := sess["session_id"].(string)
			first := quizPost(t, mux, "/api/quiz/answer",
				`{"session_id":"`+sidSession+`","concept_id":"a","answer":"4","elapsed":5}`, token)
			if first["done"] == true {
				t.Fatalf("expected the first of two answers to be mid-flight: %v", first)
			}
			last := quizPost(t, mux, "/api/quiz/answer",
				`{"session_id":"`+sidSession+`","concept_id":"b","answer":"4","elapsed":5}`, token)
			if last["done"] != true {
				t.Fatalf("expected completion: %v", last)
			}
			if qc, err := store.LastQuizCompletion(sid); err != nil || qc == nil {
				t.Fatalf("expected a completion snapshot, got %+v err=%v", qc, err)
			}

			// One below the gate, relative to the completion snapshot.
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP-1, "after-completion-49")
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusConflict {
				t.Fatalf("expected 409 at gate-1, got %d: %s", rec.Code, rec.Body.String())
			}
			refused := quizSchedulingScores(t, mux, token)
			if refused["xp_since_quiz"] != float64(engine.QuizGateXP-1) {
				t.Errorf("a refused start moved the snapshot: %v", refused)
			}
			if refused["quiz_due"] != false {
				t.Errorf("a refused start made the quiz due: %v", refused)
			}
			// The final point makes it due.
			quizSchedulingSeed(t, store, sid, 1, "after-completion-50")
			rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
			if rec.Code != http.StatusOK {
				t.Fatalf("expected 200 exactly at the gate, got %d: %s", rec.Code, rec.Body.String())
			}
			scores := quizSchedulingScores(t, mux, token)
			if scores["xp_since_quiz"] != float64(engine.QuizGateXP) {
				t.Errorf("expected xp_since_quiz %d, got %v", engine.QuizGateXP, scores["xp_since_quiz"])
			}
		})
	}
}

func TestQuizSession_RejectedAndAbandonedKeepSnapshot(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)

			// Refused below the gate: no completion row, no XP moved.
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusConflict {
				t.Fatalf("expected 409 below the gate, got %d: %s", rec.Code, rec.Body.String())
			}
			if qc, err := store.LastQuizCompletion(sid); err != nil || qc != nil {
				t.Fatalf("a refused start recorded a completion: %+v err=%v", qc, err)
			}
			if xp, _, _ := store.GetXP(sid); xp != 0 {
				t.Fatalf("a refused start changed XP: %d", xp)
			}

			// Abandoned mid-run: starting and answering part of a quiz is not
			// completion, so the snapshot does not move.
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "abandoned-gate")
			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			quizPost(t, mux, "/api/quiz/answer",
				`{"session_id":"`+sess["session_id"].(string)+`","concept_id":"a","answer":"4","elapsed":5}`, token)
			if qc, err := store.LastQuizCompletion(sid); err != nil || qc != nil {
				t.Fatalf("an abandoned quiz recorded completion: %+v err=%v", qc, err)
			}
			scores := quizSchedulingScores(t, mux, token)
			if scores["quiz_due"] != true {
				t.Errorf("an abandoned quiz consumed the gate: %v", scores)
			}
			if since, _ := scores["xp_since_quiz"].(float64); since < float64(engine.QuizGateXP) {
				t.Errorf("xp_since_quiz = %v, want at least %d", scores["xp_since_quiz"], engine.QuizGateXP)
			}
		})
	}
}

func TestQuizSession_DuplicateCreationJoinsLiveRun(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			s, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "duplicate-gate")

			first := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			second := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			if first["session_id"] != second["session_id"] {
				t.Fatalf("duplicate creation made a second run: %v then %v", first, second)
			}

			// Answering part of the run and creating again resumes the same run
			// at its next question; the recorded answer is not discarded.
			quizPost(t, mux, "/api/quiz/answer",
				`{"session_id":"`+first["session_id"].(string)+`","concept_id":"a","answer":"4","elapsed":5}`, token)
			third := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			if third["session_id"] != first["session_id"] {
				t.Fatalf("creation after a partial answer replaced the run: %v", third)
			}
			if third["concept_id"] != "b" {
				t.Errorf("expected the resumed run to be on b, got %v", third["concept_id"])
			}

			// Completion closes the run, and the same request is then refused
			// rather than treated as a retake.
			quizPost(t, mux, "/api/quiz/answer",
				`{"session_id":"`+first["session_id"].(string)+`","concept_id":"b","answer":"4","elapsed":5}`, token)
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusConflict {
				t.Fatalf("expected 409 after completion, got %d: %s", rec.Code, rec.Body.String())
			}
			s.mu.Lock()
			live := len(s.quizSessions)
			s.mu.Unlock()
			if live != 0 {
				t.Errorf("expected no live sessions after completion, got %d", live)
			}
		})
	}
}

func TestQuizSession_ConcurrentCreationJoinsOneRun(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			s, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "concurrent-gate")

			const n = 8
			type result struct {
				code int
				id   string
			}
			results := make(chan result, n)
			start := make(chan struct{})
			var wg sync.WaitGroup
			for i := 0; i < n; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
					var body map[string]interface{}
					_ = json.Unmarshal(rec.Body.Bytes(), &body)
					id, _ := body["session_id"].(string)
					results <- result{code: rec.Code, id: id}
				}()
			}
			close(start)
			wg.Wait()
			close(results)

			var want string
			for res := range results {
				if res.code != http.StatusOK {
					t.Fatalf("concurrent start returned %d", res.code)
				}
				if res.id == "" {
					t.Fatal("concurrent start returned no session id")
				}
				if want == "" {
					want = res.id
				} else if res.id != want {
					t.Errorf("concurrent starts produced two runs: %s and %s", want, res.id)
				}
			}
			s.mu.Lock()
			live := len(s.quizSessions)
			s.mu.Unlock()
			if live != 1 {
				t.Errorf("expected exactly one live session, got %d", live)
			}
			if xp, _, _ := store.GetXP(sid); xp != engine.QuizGateXP {
				t.Errorf("concurrent creation changed XP accounting: %d", xp)
			}
		})
	}
}

func TestQuizSession_NoEligibleConceptsReturnsClearReason(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "no-candidates-gate")

			rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("expected 422 with no eligible concepts, got %d: %s", rec.Code, rec.Body.String())
			}
			res := quizSchedulingDecode(t, rec)
			if res["error"] != quizNoCandidatesMessage {
				t.Errorf("expected error %q, got %v", quizNoCandidatesMessage, res["error"])
			}
			if res["error"] == quizNotDueMessage {
				t.Error("no candidate concepts must not read as the schedule's not-due answer")
			}
			if _, ok := res["session_id"]; ok {
				t.Errorf("a refused start created a session: %v", res)
			}

			// An UNSEEN row is state without learning, so it is not a candidate.
			seedQuizUnseen(t, store, sid, "a")
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("expected 422 for an UNSEEN row, got %d: %s", rec.Code, rec.Body.String())
			}

			// A learned concept makes the same request valid, proving the 422
			// was about candidates rather than the gate.
			seedQuizLearned(t, store, sid, "a", string(mastery.StatusPracticing), 0.5)
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusOK {
				t.Fatalf("expected 200 after learning a concept, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestQuizSession_DifficultyUsesEvidenceWindow(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			s, mux := quizSchedulingServer(t, store)
			token, sid := quizDueGuest(t, store, mux)

			start := func() *quiz.Session {
				t.Helper()
				rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
				if rec.Code != http.StatusOK {
					t.Fatalf("start quiz: %d %s", rec.Code, rec.Body.String())
				}
				res := quizSchedulingDecode(t, rec)
				id, _ := res["session_id"].(string)
				s.mu.Lock()
				sess := s.quizSessions[id]
				s.mu.Unlock()
				if sess == nil {
					t.Fatalf("session %q not found", id)
				}
				return sess
			}

			sess := start()
			if got, want := sess.Difficulties["a"], mastery.NeutralDifficulty; math.Abs(got-want) > 1e-9 {
				t.Errorf("difficulty with one attempt = %v, want neutral %v", got, want)
			}

			now := time.Now().UTC()
			if err := store.UpsertProgress(&storage.ConceptProgress{
				StudentID:     sid,
				ConceptID:     "a",
				Status:        string(mastery.StatusPracticing),
				Attempts:      mastery.EvidenceWindow,
				LastAttempted: &now,
				WeaknessScore: 0.9,
			}); err != nil {
				t.Fatalf("seed full evidence window: %v", err)
			}
			s.dropStudentLiveSessions(sid)

			sess = start()
			if got, want := sess.Difficulties["a"], 0.35; math.Abs(got-want) > 1e-9 {
				t.Errorf("difficulty with full evidence window = %v, want %v", got, want)
			}
		})
	}
}

func TestQuizSession_MasteredConceptsAreNotRoutineCandidates(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "mastered-only-gate")
			seedQuizLearned(t, store, sid, "a", string(mastery.StatusMastered), 0.9)

			// Mastered material is retained through review, not re-demonstrated
			// in a routine mastery check, so one mastered concept is not a pool.
			rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("expected 422 for a mastered-only pool, got %d: %s", rec.Code, rec.Body.String())
			}
			if res := quizSchedulingDecode(t, rec); res["error"] != quizNoCandidatesMessage {
				t.Errorf("expected %q, got %v", quizNoCandidatesMessage, res["error"])
			}

			// A concept still being acquired makes the same request valid.
			seedQuizLearned(t, store, sid, "b", string(mastery.StatusPracticing), 0.5)
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusOK {
				t.Fatalf("expected 200 with an acquiring concept, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestQuizSession_AnonymousQuizModeRemainsUngated(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)

			// No token and no student_id: no XP, no mastery writes, and no gate
			// to meet. This is the existing anonymous quiz mode, not a bypass of
			// an identified learner's schedule.
			first := quizSchedulingPost(t, mux, "/api/quiz/session", "", `{}`)
			second := quizSchedulingPost(t, mux, "/api/quiz/session", "", `{}`)
			if first.Code != http.StatusOK || second.Code != http.StatusOK {
				t.Fatalf("anonymous quiz refused: %d / %d", first.Code, second.Code)
			}
			a := quizSchedulingDecode(t, first)
			b := quizSchedulingDecode(t, second)
			if a["session_id"] == b["session_id"] {
				t.Error("anonymous quiz sessions should not be shared between callers")
			}
		})
	}
}

func TestQuizSession_AuthenticationBehaviourUnchanged(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, _ := quizGuest(t, mux)
			_, otherSID := quizGuest(t, mux)
			quizSchedulingSeed(t, store, otherSID, engine.QuizGateXP, "other-gate")

			// The authenticated identity wins, and a body that names another
			// learner is refused before the gate is consulted.
			rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{"student_id":"`+otherSID+`"}`)
			if rec.Code != http.StatusForbidden {
				t.Fatalf("expected 403 for a foreign student_id, got %d: %s", rec.Code, rec.Body.String())
			}

			// An unauthenticated caller may name only a legacy guest_-prefixed
			// id. That path is identified too, so the gate still applies; an id
			// with no XP cannot start a run.
			rec = quizSchedulingPost(t, mux, "/api/quiz/session", "", `{"student_id":"guest_no_such_learner"}`)
			if rec.Code != http.StatusConflict {
				t.Fatalf("expected 409 for an unseeded guest id, got %d: %s", rec.Code, rec.Body.String())
			}

			// A registered id cannot be claimed by an unauthenticated caller.
			rec = quizSchedulingPost(t, mux, "/api/quiz/session", "", `{"student_id":"someone_registered"}`)
			if rec.Code != http.StatusForbidden {
				t.Fatalf("expected 403 for an unowned student_id, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestQuizSession_DiagnosticIsNotGated(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)

			// Zero XP is enough for the goal diagnostic, which is placement,
			// not a routine mastery check.
			diag := quizSchedulingPost(t, mux, "/api/goal/diagnostic", token, `{"concept_ids":["a","b"]}`)
			if diag.Code != http.StatusOK {
				t.Fatalf("diagnostic at 0 XP: %d %s", diag.Code, diag.Body.String())
			}
			if res := quizSchedulingDecode(t, diag); res["session_id"] == "" {
				t.Fatalf("diagnostic returned no session: %v", res)
			}
			// The same learner still cannot start the routine quiz.
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusConflict {
				t.Fatalf("expected 409 for the gated quiz, got %d: %s", rec.Code, rec.Body.String())
			}
			// Starting a diagnostic did not seed mastery from nothing.
			if p, _ := store.GetProgress(sid, "a"); p != nil && p.Status == "MASTERED" {
				t.Errorf("diagnostic start wrote MASTERED: %+v", p)
			}
		})
	}
}
