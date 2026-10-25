package server

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/chuma-beep/mathua/internal/engine"
)

// The completion screen and the session endpoint must tell the same story.
// These tests drive a real two-concept mastery check through the HTTP surface
// and pin what the learner is offered afterwards: an answered run consumes the
// XP snapshot and withdraws the retake, while a fully skipped run does not.
func TestQuizLifecycle_AnsweredCompletionWithdrawsRetake(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "lifecycle-answered-gate")

			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			sessionID := sess["session_id"].(string)
			first := quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"a","answer":"4","elapsed":5}`, sessionID), token)
			if first["done"] == true {
				t.Fatalf("first answer finished a two-concept run: %v", first)
			}
			final := quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"b","answer":"4","elapsed":5}`, sessionID), token)
			if final["done"] != true {
				t.Fatalf("expected completion, got %v", final)
			}
			if final["retake_available"] != false {
				t.Errorf("answered completion offered a retake the next request refuses: %v", final)
			}

			// The answers still flowed through the evidence pipeline.
			if p, err := store.GetProgress(sid, "a"); err != nil || p == nil || p.Attempts < 2 {
				t.Errorf("quiz answer did not reach progress: %+v err=%v", p, err)
			}

			scores := quizSchedulingScores(t, mux, token)
			if scores["quiz_due"] != false || scores["xp_since_quiz"] != float64(0) {
				t.Errorf("completion did not consume the gate: %v", scores)
			}
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusConflict {
				t.Fatalf("immediate retake: got %d, want 409: %s", rec.Code, rec.Body.String())
			}

			// Earning the gate again makes a new check admissible, so the false
			// in the completion response was about the snapshot, not a permanent
			// lock.
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "lifecycle-answered-next-gate")
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusOK {
				t.Fatalf("retake after earning the gate: got %d, want 200: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestQuizLifecycle_FullySkippedRunKeepsRetake(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "lifecycle-skipped-gate")

			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			sessionID := sess["session_id"].(string)
			rec := quizSchedulingPost(t, mux, "/api/quiz/skip", token,
				fmt.Sprintf(`{"session_id":%q}`, sessionID))
			if rec.Code != http.StatusOK {
				t.Fatalf("first skip: %d %s", rec.Code, rec.Body.String())
			}
			if quizSchedulingDecode(t, rec)["done"] == true {
				t.Fatal("first skip finished a two-concept run")
			}
			rec = quizSchedulingPost(t, mux, "/api/quiz/skip", token,
				fmt.Sprintf(`{"session_id":%q}`, sessionID))
			if rec.Code != http.StatusOK {
				t.Fatalf("second skip: %d %s", rec.Code, rec.Body.String())
			}
			final := quizSchedulingDecode(t, rec)
			if final["done"] != true {
				t.Fatalf("expected completion, got %v", final)
			}
			if final["retake_available"] != true {
				t.Errorf("fully skipped run should leave the gate available: %v", final)
			}
			if qc, err := store.LastQuizCompletion(sid); err != nil || qc != nil {
				t.Errorf("fully skipped run recorded a completion: %+v err=%v", qc, err)
			}
			scores := quizSchedulingScores(t, mux, token)
			if scores["quiz_due"] != true {
				t.Errorf("fully skipped run consumed the gate: %v", scores)
			}
			if since, _ := scores["xp_since_quiz"].(float64); since < float64(engine.QuizGateXP) {
				t.Errorf("xp_since_quiz = %v, want at least %d", scores["xp_since_quiz"], engine.QuizGateXP)
			}
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusOK {
				t.Fatalf("retake after a fully skipped run: got %d, want 200: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestQuizLifecycle_PartialAnswerThenSkipWithdrawsRetake(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "lifecycle-partial-gate")

			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			sessionID := sess["session_id"].(string)
			first := quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"a","answer":"4","elapsed":5}`, sessionID), token)
			if first["done"] == true {
				t.Fatalf("first answer finished a two-concept run: %v", first)
			}
			rec := quizSchedulingPost(t, mux, "/api/quiz/skip", token,
				fmt.Sprintf(`{"session_id":%q}`, sessionID))
			if rec.Code != http.StatusOK {
				t.Fatalf("skip to completion: %d %s", rec.Code, rec.Body.String())
			}
			final := quizSchedulingDecode(t, rec)
			if final["done"] != true {
				t.Fatalf("expected completion, got %v", final)
			}
			if final["retake_available"] != false {
				t.Errorf("run with an answer moved the snapshot but offered a retake: %v", final)
			}
			if rec := quizSchedulingPost(t, mux, "/api/quiz/session", token, `{}`); rec.Code != http.StatusConflict {
				t.Fatalf("immediate retake: got %d, want 409: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestQuizLifecycle_RepeatedCompletionIsGone(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			_, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "lifecycle-repeat-gate")

			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			sessionID := sess["session_id"].(string)
			quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"a","answer":"4","elapsed":5}`, sessionID), token)
			final := quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"b","answer":"4","elapsed":5}`, sessionID), token)
			if final["done"] != true {
				t.Fatalf("expected completion, got %v", final)
			}

			// XP earned after completion belongs to the next check. Replaying
			// the closed run must not snapshot it again.
			quizSchedulingSeed(t, store, sid, 10, "lifecycle-after-completion")
			rec := quizSchedulingPost(t, mux, "/api/quiz/answer", token,
				fmt.Sprintf(`{"session_id":%q,"concept_id":"b","answer":"4","elapsed":5}`, sessionID))
			if rec.Code != http.StatusNotFound {
				t.Fatalf("replayed completion: got %d, want 404: %s", rec.Code, rec.Body.String())
			}
			scores := quizSchedulingScores(t, mux, token)
			if scores["xp_since_quiz"] != float64(10) || scores["quiz_due"] != false {
				t.Errorf("replayed completion moved the snapshot: %v", scores)
			}
		})
	}
}

func TestQuizLifecycle_MissQueuesRemediation(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			s, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "lifecycle-remedial-gate")

			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			sessionID := sess["session_id"].(string)
			miss := quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"a","answer":"999","elapsed":5}`, sessionID), token)
			if miss["correct"] != false {
				t.Fatalf("expected a miss, got %v", miss)
			}
			remedial, _ := miss["remedial"].([]interface{})
			foundResponse := false
			for _, id := range remedial {
				if id == "a" {
					foundResponse = true
				}
			}
			if !foundResponse {
				t.Errorf("miss did not return the missed concept for remediation: %v", miss)
			}
			foundQueue := false
			for _, id := range s.eng.QuizRemedial(sid) {
				if id == "a" {
					foundQueue = true
				}
			}
			if !foundQueue {
				t.Errorf("miss did not reach the remediation queue: %v", s.eng.QuizRemedial(sid))
			}

			final := quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"b","answer":"4","elapsed":5}`, sessionID), token)
			if final["done"] != true || final["retake_available"] != false {
				t.Errorf("failed check completion: %v", final)
			}
		})
	}
}

func TestQuizLifecycle_ScoresReportActiveSession(t *testing.T) {
	for tag, store := range quizSchedulingStores(t) {
		t.Run(tag, func(t *testing.T) {
			s, mux := quizSchedulingServer(t, store)
			token, sid := quizGuest(t, mux)
			quizSchedulingLearn(t, store, sid)
			quizSchedulingSeed(t, store, sid, engine.QuizGateXP, "lifecycle-active-gate")

			before := quizSchedulingScores(t, mux, token)
			if before["quiz_active"] != false || before["quiz_due"] != true {
				t.Fatalf("before start: %v", before)
			}

			sess := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			sessionID := sess["session_id"].(string)
			active := quizSchedulingScores(t, mux, token)
			if active["quiz_active"] != true || active["quiz_due"] != true {
				t.Fatalf("active session not reported: %v", active)
			}
			if active["xp_total"] != before["xp_total"] {
				t.Errorf("starting a session changed XP: %v -> %v", before["xp_total"], active["xp_total"])
			}

			// Refresh/reconnection reuses the run and does not create another.
			again := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			if again["session_id"] != sessionID {
				t.Fatalf("refresh created a second session: %v -> %v", sessionID, again["session_id"])
			}
			if after := quizSchedulingScores(t, mux, token); after["quiz_active"] != true {
				t.Fatalf("refresh lost the active session: %v", after)
			}

			// Another learner cannot see this run through their own scores.
			otherToken, _ := quizGuest(t, mux)
			if other := quizSchedulingScores(t, mux, otherToken); other["quiz_active"] != false {
				t.Errorf("another learner saw an active session: %v", other)
			}

			// An expired or invalidated live run is gone, but the gate remains due.
			s.dropStudentLiveSessions(sid)
			if expired := quizSchedulingScores(t, mux, token); expired["quiz_active"] != false || expired["quiz_due"] != true {
				t.Errorf("dropped session still reported active: %v", expired)
			}

			// Completion clears the live-session flag as well as the gate.
			next := quizPost(t, mux, "/api/quiz/session", `{}`, token)
			nextID := next["session_id"].(string)
			quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"a","answer":"4","elapsed":5}`, nextID), token)
			quizPost(t, mux, "/api/quiz/answer",
				fmt.Sprintf(`{"session_id":%q,"concept_id":"b","answer":"4","elapsed":5}`, nextID), token)
			if done := quizSchedulingScores(t, mux, token); done["quiz_active"] != false || done["quiz_due"] != false {
				t.Errorf("completed run still reported active: %v", done)
			}
		})
	}
}
