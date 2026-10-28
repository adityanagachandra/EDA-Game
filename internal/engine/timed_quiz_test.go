package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
	"github.com/chuma-beep/mathua/internal/timedquiz"
)

func TestTimedQuizGateAndFSRSIsolation(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("timed-quiz-learner")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:       st.ID,
		ConceptID:       "a",
		Status:          string(mastery.StatusLearning),
		Streak:          2,
		Attempts:        2,
		AvgResponseTime: 8,
		SM2EFactor:      2.5,
		SM2Interval:     3,
	}); err != nil {
		t.Fatalf("seed progress: %v", err)
	}

	if due, err := e.TimedQuizDue(st.ID); err != nil || due {
		t.Fatalf("new learner due = %v, err = %v; want false", due, err)
	}
	if err := e.repo.AwardXP(st.ID, timedquiz.GateXP-1, "test", "timed-gate-1"); err != nil {
		t.Fatalf("award xp: %v", err)
	}
	if due, err := e.TimedQuizDue(st.ID); err != nil || due {
		t.Fatalf("learner at %d XP due = %v, err = %v; want false", timedquiz.GateXP-1, due, err)
	}
	if err := e.repo.AwardXP(st.ID, 1, "test", "timed-gate-2"); err != nil {
		t.Fatalf("award final xp: %v", err)
	}
	if due, err := e.TimedQuizDue(st.ID); err != nil || !due {
		t.Fatalf("learner at gate due = %v, err = %v; want true", due, err)
	}

	run, items, err := e.StartTimedQuiz(st.ID)
	if err != nil {
		t.Fatalf("start timed quiz: %v", err)
	}
	if run == nil || len(items) == 0 {
		t.Fatalf("start returned run=%+v items=%d", run, len(items))
	}
	if run.DeadlineAt.Before(run.StartedAt.Add(timedquiz.Duration - time.Second)) {
		t.Fatalf("deadline = %v, want roughly %v after start", run.DeadlineAt, timedquiz.Duration)
	}
	resumed, resumedItems, err := e.StartTimedQuiz(st.ID)
	if err != nil {
		t.Fatalf("resume timed quiz: %v", err)
	}
	if resumed.ID != run.ID || len(resumedItems) != len(items) {
		t.Fatalf("resume returned a different run: %s vs %s", resumed.ID, run.ID)
	}

	item := items[0]
	res, err := e.SubmitTimedQuizAnswer(st.ID, item.ConceptID, item.ExpectedAnswer, item.ExpectedAnswer, 5, item.Question, item.Difficulty, "timed-sub-1")
	if err != nil {
		t.Fatalf("submit timed answer: %v", err)
	}
	if !res.Correct {
		t.Fatalf("timed answer was not correct: %+v", res)
	}
	progress, err := e.repo.GetProgress(st.ID, item.ConceptID)
	if err != nil || progress == nil {
		t.Fatalf("get progress: %+v %v", progress, err)
	}
	if progress.FSRSStability != nil || progress.FSRSDifficulty != nil || progress.NextReviewDue != nil {
		t.Fatalf("timed quiz moved FSRS: stability=%v difficulty=%v due=%v", progress.FSRSStability, progress.FSRSDifficulty, progress.NextReviewDue)
	}
	stored := progress
	totalBefore, _, err := e.repo.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp: %v", err)
	}
	// Replaying the same submission must not add a second attempt or XP.
	if _, err := e.SubmitTimedQuizAnswer(st.ID, item.ConceptID, item.ExpectedAnswer, item.ExpectedAnswer, 5, item.Question, item.Difficulty, "timed-sub-1"); err != nil {
		t.Fatalf("replay timed answer: %v", err)
	}
	totalAfter, _, err := e.repo.GetXP(st.ID)
	if err != nil {
		t.Fatalf("get xp after replay: %v", err)
	}
	if totalAfter != totalBefore {
		t.Fatalf("replay changed XP: %d -> %d", totalBefore, totalAfter)
	}
	progressAfter, err := e.repo.GetProgress(st.ID, item.ConceptID)
	if err != nil {
		t.Fatalf("get progress after replay: %v", err)
	}
	if progressAfter.Attempts != stored.Attempts {
		t.Fatalf("replay added an attempt: %d -> %d", stored.Attempts, progressAfter.Attempts)
	}

	// Complete the run and verify the gate is consumed: XP earned before the
	// completion does not make the next timed quiz due.
	store := e.repo.(storage.TimedQuizStore)
	if _, err := store.FinalizeTimedQuizRun(run.ID, st.ID, time.Now().UTC()); err != nil {
		t.Fatalf("finalize run: %v", err)
	}
	if due, err := e.TimedQuizDue(st.ID); err != nil || due {
		t.Fatalf("after completion due = %v, err = %v; want false", due, err)
	}
}
