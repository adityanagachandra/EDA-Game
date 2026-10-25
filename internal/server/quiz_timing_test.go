package server

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/storage"
)

// The quiz time limit is a slow-answer threshold, not a hard stop. Commit
// b70fd5c1 removed the visible countdown on purpose and kept the measurement:
// a learner who thinks for a minute must not be thrown into a failing state.
// These tests pin the server side of that decision. The server stamps the
// question when it is served, measures from that stamp when the answer arrives
// and never lets a client report a faster answer than the server observed.
// Late answers are still graded, because the limit classifies speed rather than
// expiring the run.

func quizTimingLearner(t *testing.T, store storage.Repository, mux *http.ServeMux) (token, studentID string) {
	t.Helper()
	token, studentID = quizGuest(t, mux)
	quizSchedulingSeed(t, store, studentID, engine.QuizGateXP, "timing-gate")
	quizSchedulingLearn(t, store, studentID)
	return token, studentID
}

func quizTimingStart(t *testing.T, mux *http.ServeMux, token string) (sessionID, conceptID string) {
	t.Helper()
	rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("start timing quiz: %d %s", rec.Code, rec.Body.String())
	}
	res := quizSchedulingDecode(t, rec)
	if res["done"] == true {
		t.Fatalf("timing quiz started done: %v", res)
	}
	sessionID, _ = res["session_id"].(string)
	conceptID, _ = res["concept_id"].(string)
	if sessionID == "" || conceptID == "" {
		t.Fatalf("timing quiz missing ids: %v", res)
	}
	return sessionID, conceptID
}

func backdateQuizQuestion(t *testing.T, s *Server, sessionID string, d time.Duration) {
	t.Helper()
	s.mu.Lock()
	sess := s.quizSessions[sessionID]
	s.mu.Unlock()
	if sess == nil {
		t.Fatalf("quiz session %s is not live", sessionID)
	}
	sess.Lock()
	sess.QuestionStarted = time.Now().UTC().Add(-d)
	sess.Unlock()
}

func quizQuestionStarted(t *testing.T, s *Server, sessionID string) time.Time {
	t.Helper()
	s.mu.Lock()
	sess := s.quizSessions[sessionID]
	s.mu.Unlock()
	if sess == nil {
		t.Fatalf("quiz session %s is not live", sessionID)
	}
	sess.Lock()
	defer sess.Unlock()
	return sess.QuestionStarted
}

func quizTimingAttempts(t *testing.T, store storage.Repository, studentID string) []storage.AttemptEntry {
	t.Helper()
	attempts, err := store.GetAttemptsForStudent(studentID)
	if err != nil {
		t.Fatalf("load attempts: %v", err)
	}
	return attempts
}

func TestQuizTiming_LateAnswerIsGradedWithServerElapsed(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				delay time.Duration
			}{
				// The concept threshold in this fixture is 60 seconds.
				{"at the limit", 60 * time.Second},
				{"past the limit", 75 * time.Second},
			} {
				t.Run(tc.name, func(t *testing.T) {
					s, mux := quizSchedulingServer(t, store)
					token, sid := quizTimingLearner(t, store, mux)
					sessionID, conceptID := quizTimingStart(t, mux, token)
					backdateQuizQuestion(t, s, sessionID, tc.delay)

					// The client claims half a second. The server has watched
					// the question for at least tc.delay and must use that.
					body := fmt.Sprintf(`{"session_id":%q,"concept_id":%q,"answer":"4","elapsed":0.5}`, sessionID, conceptID)
					rec := quizSchedulingPost(t, mux, "/api/quiz/answer", token, body)
					if rec.Code != http.StatusOK {
						t.Fatalf("a late answer must still be graded, got %d: %s", rec.Code, rec.Body.String())
					}
					res := quizSchedulingDecode(t, rec)
					if res["correct"] != true {
						t.Errorf("expected the late correct answer to be graded correct: %v", res)
					}
					if res["done"] == true {
						t.Errorf("two-question fixture ended after one answer: %v", res)
					}

					attempts := quizTimingAttempts(t, store, sid)
					if len(attempts) != 1 {
						t.Fatalf("expected one recorded attempt, got %d", len(attempts))
					}
					if attempts[0].Source != "quiz" {
						t.Errorf("attempt source = %q, want %q", attempts[0].Source, "quiz")
					}
					if attempts[0].ElapsedSeconds < tc.delay.Seconds() {
						t.Errorf("recorded elapsed = %.3fs, want at least %.1fs from the server clock",
							attempts[0].ElapsedSeconds, tc.delay.Seconds())
					}
				})
			}
		})
	}
}

func TestQuizTiming_ClientCannotClaimAFasterAnswer(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			s, mux := quizSchedulingServer(t, store)
			token, sid := quizTimingLearner(t, store, mux)
			sessionID, conceptID := quizTimingStart(t, mux, token)
			backdateQuizQuestion(t, s, sessionID, 10*time.Second)

			// 0.1 is below the human floor and above the server's own clock.
			// The request is not refused and the persisted time is the server's.
			body := fmt.Sprintf(`{"session_id":%q,"concept_id":%q,"answer":"4","elapsed":0.1}`, sessionID, conceptID)
			rec := quizSchedulingPost(t, mux, "/api/quiz/answer", token, body)
			if rec.Code != http.StatusOK {
				t.Fatalf("server clock should satisfy the too-quick check, got %d: %s", rec.Code, rec.Body.String())
			}
			attempts := quizTimingAttempts(t, store, sid)
			if len(attempts) != 1 {
				t.Fatalf("expected one recorded attempt, got %d", len(attempts))
			}
			if attempts[0].ElapsedSeconds < 10 {
				t.Errorf("client claimed 0.1s after 10s; recorded %.3fs", attempts[0].ElapsedSeconds)
			}
		})
	}
}

func TestQuizTiming_RefreshDoesNotResetTheQuestionClock(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			s, mux := quizSchedulingServer(t, store)
			token, sid := quizTimingLearner(t, store, mux)
			sessionID, conceptID := quizTimingStart(t, mux, token)

			// A refresh or reconnect gets the same run and the same question.
			// It must not restart the clock, or a learner could reload to make
			// every answer look fast.
			backdateQuizQuestion(t, s, sessionID, 40*time.Second)
			before := quizQuestionStarted(t, s, sessionID)
			rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`)
			if rec.Code != http.StatusOK {
				t.Fatalf("refresh: %d %s", rec.Code, rec.Body.String())
			}
			res := quizSchedulingDecode(t, rec)
			if res["session_id"] != sessionID || res["concept_id"] != conceptID {
				t.Fatalf("refresh did not resume the same question: %v", res)
			}
			after := quizQuestionStarted(t, s, sessionID)
			if !after.Equal(before) {
				t.Errorf("refresh reset the question clock: %v then %v", before, after)
			}

			// Time already spent stays spent across the refresh.
			backdateQuizQuestion(t, s, sessionID, 90*time.Second)
			body := fmt.Sprintf(`{"session_id":%q,"concept_id":%q,"answer":"4","elapsed":0.5}`, sessionID, conceptID)
			rec = quizSchedulingPost(t, mux, "/api/quiz/answer", token, body)
			if rec.Code != http.StatusOK {
				t.Fatalf("answer after refresh: %d %s", rec.Code, rec.Body.String())
			}
			attempts := quizTimingAttempts(t, store, sid)
			if len(attempts) != 1 {
				t.Fatalf("expected one recorded attempt, got %d", len(attempts))
			}
			if attempts[0].ElapsedSeconds < 90 {
				t.Errorf("elapsed after refresh = %.3fs, want at least 90s", attempts[0].ElapsedSeconds)
			}
		})
	}
}

func TestQuizTiming_SkippedQuestionsRecordNoAttemptOrSnapshot(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizTimingLearner(t, store, mux)
			sessionID, conceptID := quizTimingStart(t, mux, token)

			first := quizSchedulingPost(t, mux, "/api/quiz/skip", token,
				fmt.Sprintf(`{"session_id":%q,"concept_id":%q}`, sessionID, conceptID))
			if first.Code != http.StatusOK {
				t.Fatalf("skip first question: %d %s", first.Code, first.Body.String())
			}
			if res := quizSchedulingDecode(t, first); res["done"] == true {
				t.Fatalf("first skip ended a two-question run: %v", res)
			}
			nextConcept, _ := quizSchedulingDecode(t, first)["concept_id"].(string)
			second := quizSchedulingPost(t, mux, "/api/quiz/skip", token,
				fmt.Sprintf(`{"session_id":%q,"concept_id":%q}`, sessionID, nextConcept))
			if second.Code != http.StatusOK {
				t.Fatalf("skip second question: %d %s", second.Code, second.Body.String())
			}
			if res := quizSchedulingDecode(t, second); res["done"] != true {
				t.Fatalf("expected the second skip to finish the run: %v", res)
			}

			// Unanswered means untreated: no attempt, no XP, and because no
			// answer was recorded the gate snapshot does not move.
			if attempts := quizTimingAttempts(t, store, sid); len(attempts) != 0 {
				t.Errorf("skips recorded %d attempts, want none", len(attempts))
			}
			if qc, err := store.LastQuizCompletion(sid); err != nil || qc != nil {
				t.Fatalf("a fully skipped run recorded completion: %+v err=%v", qc, err)
			}
			scores := quizSchedulingScores(t, mux, token)
			if scores["quiz_due"] != true {
				t.Errorf("a fully skipped run consumed the gate: %v", scores)
			}
			if since, _ := scores["xp_since_quiz"].(float64); since < float64(engine.QuizGateXP) {
				t.Errorf("xp_since_quiz = %v, want at least %d", scores["xp_since_quiz"], engine.QuizGateXP)
			}
		})
	}
}

func TestQuizTiming_AnswerAndSkipAreOwnerOnly(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			ownerToken, _ := quizTimingLearner(t, store, mux)
			otherToken, otherSID := quizGuest(t, mux)
			sessionID, conceptID := quizTimingStart(t, mux, ownerToken)

			answer := fmt.Sprintf(`{"session_id":%q,"concept_id":%q,"answer":"4","elapsed":5,"student_id":%q}`, sessionID, conceptID, otherSID)
			if rec := quizSchedulingPost(t, mux, "/api/quiz/answer", otherToken, answer); rec.Code != http.StatusForbidden {
				t.Errorf("another learner answered a quiz session: %d %s", rec.Code, rec.Body.String())
			}
			skip := fmt.Sprintf(`{"session_id":%q,"concept_id":%q,"student_id":%q}`, sessionID, conceptID, otherSID)
			if rec := quizSchedulingPost(t, mux, "/api/quiz/skip", otherToken, skip); rec.Code != http.StatusForbidden {
				t.Errorf("another learner skipped a quiz session: %d %s", rec.Code, rec.Body.String())
			}

			// The refused requests did not consume the owner's question.
			own := fmt.Sprintf(`{"session_id":%q,"concept_id":%q,"answer":"4","elapsed":5}`, sessionID, conceptID)
			if rec := quizSchedulingPost(t, mux, "/api/quiz/answer", ownerToken, own); rec.Code != http.StatusOK {
				t.Fatalf("owner could not answer after foreign attempts: %d %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestEffectiveAnswerElapsed(t *testing.T) {
	for _, tc := range []struct {
		name   string
		client float64
		server float64
		want   float64
	}{
		{"server is later", 5, 8, 8},
		{"client is later", 12, 8, 12},
		{"client negative", -3, 4, 4},
		{"both negative", -3, -1, 0},
		{"equal", 7, 7, 7},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := effectiveAnswerElapsed(tc.client, tc.server); got != tc.want {
				t.Errorf("effectiveAnswerElapsed(%v, %v) = %v, want %v", tc.client, tc.server, got, tc.want)
			}
		})
	}
}
