package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/auth"
	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
	"github.com/chuma-beep/mathua/internal/timedquiz"
)

func timedQuizServer(t *testing.T) *Server {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60}},
	})
	if err != nil {
		t.Fatalf("build dag: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{})
	return New(engine.New(store, d, reg, nil, nil), store, auth.New(store))
}

func timedQuizGet(t *testing.T, mux *http.ServeMux, token, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	return rec
}

func TestTimedQuizLifecycleIsDurableAndIdempotent(t *testing.T) {
	s := timedQuizServer(t)
	mux := http.NewServeMux()
	s.Register(mux)
	token, studentID := quizGuest(t, mux)
	seedQuizLearned(t, s.repo, studentID, "a", string(mastery.StatusLearning), 0)

	// The gate is positive XP after the rollout baseline, not lifetime XP.
	quizSchedulingSeed(t, s.repo, studentID, timedquiz.GateXP-1, "timed-gate-a")
	summary := quizSchedulingDecode(t, timedQuizGet(t, mux, token, "/api/timed-quiz"))
	if summary["due"] != false {
		t.Fatalf("summary due at %v XP: %v", summary["eligible_xp"], summary)
	}
	quizSchedulingSeed(t, s.repo, studentID, 1, "timed-gate-b")
	summary = quizSchedulingDecode(t, timedQuizGet(t, mux, token, "/api/timed-quiz"))
	if summary["due"] != true {
		t.Fatalf("summary not due at gate: %v", summary)
	}

	start := quizSchedulingPost(t, mux, "/api/timed-quiz/session", token, `{}`)
	if start.Code != http.StatusOK {
		t.Fatalf("start timed quiz: %d %s", start.Code, start.Body.String())
	}
	started := quizSchedulingDecode(t, start)
	runID, _ := started["run_id"].(string)
	items, _ := started["items"].([]interface{})
	if runID == "" || len(items) != 1 {
		t.Fatalf("started run = %v", started)
	}
	question, _ := items[0].(map[string]interface{})
	itemID, _ := question["item_id"].(string)
	if itemID == "" {
		t.Fatalf("started item missing id: %v", question)
	}
	if _, leaked := question["expected_answer"]; leaked {
		t.Fatalf("session leaked the expected answer: %v", question)
	}

	store, ok := s.repo.(storage.TimedQuizStore)
	if !ok {
		t.Fatal("timed quiz store unavailable")
	}
	storedItems, err := store.GetTimedQuizItems(runID)
	if err != nil || len(storedItems) != 1 {
		t.Fatalf("stored items = %d, err = %v", len(storedItems), err)
	}
	expected := storedItems[0].ExpectedAnswer
	body := fmt.Sprintf(`{"run_id":%q,"item_id":%q,"answer":%q,"elapsed":5,"submission_id":"timed-sub-1"}`, runID, itemID, expected)
	answer := quizSchedulingPost(t, mux, "/api/timed-quiz/answer", token, body)
	if answer.Code != http.StatusOK {
		t.Fatalf("answer timed quiz: %d %s", answer.Code, answer.Body.String())
	}
	result := quizSchedulingDecode(t, answer)
	if result["correct"] != true || result["completed"] != true {
		t.Fatalf("answer result = %v", result)
	}
	if _, leaked := result["expected_answer"]; leaked {
		t.Fatalf("answer leaked the expected answer: %v", result)
	}

	// A replay must return the stored item and must not award XP or create a
	// second attempt.
	totalBefore, _, err := s.repo.GetXP(studentID)
	if err != nil {
		t.Fatalf("xp before replay: %v", err)
	}
	replay := quizSchedulingPost(t, mux, "/api/timed-quiz/answer", token, body)
	if replay.Code != http.StatusOK {
		t.Fatalf("replay timed quiz: %d %s", replay.Code, replay.Body.String())
	}
	totalAfter, _, err := s.repo.GetXP(studentID)
	if err != nil {
		t.Fatalf("xp after replay: %v", err)
	}
	if totalAfter != totalBefore {
		t.Fatalf("replay changed XP: %d -> %d", totalBefore, totalAfter)
	}
	progress, err := s.repo.GetProgress(studentID, "a")
	if err != nil || progress == nil {
		t.Fatalf("progress after timed quiz: %+v %v", progress, err)
	}
	if progress.FSRSStability != nil || progress.FSRSDifficulty != nil || progress.NextReviewDue != nil {
		t.Fatalf("timed quiz moved FSRS: stability=%v difficulty=%v due=%v", progress.FSRSStability, progress.FSRSDifficulty, progress.NextReviewDue)
	}

	// The completed run consumes the gate; the same earned XP cannot start a
	// second run.
	summary = quizSchedulingDecode(t, timedQuizGet(t, mux, token, "/api/timed-quiz"))
	if summary["due"] != false {
		t.Fatalf("gate was not consumed by completion: %v", summary)
	}
}

func TestTimedQuizAnswerRejectsLateSubmission(t *testing.T) {
	s := timedQuizServer(t)
	mux := http.NewServeMux()
	s.Register(mux)
	token, studentID := quizGuest(t, mux)
	seedQuizLearned(t, s.repo, studentID, "a", string(mastery.StatusLearning), 0)
	store := s.repo.(storage.TimedQuizStore)
	now := time.Now().UTC()
	run := storage.TimedQuizRun{
		ID:         fmt.Sprintf("late-run-%d", now.UnixNano()),
		StudentID:  studentID,
		StartedAt:  now.Add(-11 * time.Minute),
		DeadlineAt: now.Add(-time.Minute),
		Status:     "active",
		CreatedAt:  now.Add(-11 * time.Minute),
	}
	item := storage.TimedQuizItem{ID: run.ID + ":0", RunID: run.ID, ConceptID: "a", Position: 0, Question: "q", ExpectedAnswer: "42", ServedAt: run.StartedAt, Status: "served"}
	if err := store.CreateTimedQuizRun(run, []storage.TimedQuizItem{item}); err != nil {
		t.Fatalf("create late run: %v", err)
	}
	body := fmt.Sprintf(`{"run_id":%q,"item_id":%q,"answer":"42","elapsed":5,"submission_id":"late-1"}`, run.ID, item.ID)
	rec := quizSchedulingPost(t, mux, "/api/timed-quiz/answer", token, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("late answer status = %d, want 200 expired result: %s", rec.Code, rec.Body.String())
	}
	res := quizSchedulingDecode(t, rec)
	if res["expired"] != true {
		t.Fatalf("late answer result = %v, want expired", res)
	}
}
