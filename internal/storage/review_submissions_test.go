package storage

import (
	"testing"
	"time"
)

func reviewProgress(studentID string) *ConceptProgress {
	stability := 10.0
	difficulty := 4.0
	lastReview := time.Now().UTC().Add(-24 * time.Hour)
	due := time.Now().UTC().Add(-time.Hour)
	return &ConceptProgress{
		StudentID:       studentID,
		ConceptID:       "a",
		Status:          "MASTERED",
		Streak:          1,
		AvgResponseTime: 5,
		LastReviewed:    &lastReview,
		NextReviewDue:   &due,
		SM2Repetitions:  2,
		SM2Interval:     6,
		SM2EFactor:      2.5,
		FSRSStability:   &stability,
		FSRSDifficulty:  &difficulty,
		FSRSLastReview:  &lastReview,
	}
}

func TestReviewSubmissionIsAtomicAndIdempotent(t *testing.T) {
	store := newTestStore(t)
	st := seedStudent(t, store)
	rs, ok := store.(ReviewSubmissionStore)
	if !ok {
		t.Fatalf("store %T does not implement ReviewSubmissionStore", store)
	}
	now := time.Now().UTC()
	sess, err := store.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	progress := reviewProgress(st.ID)
	event := ReviewEvent{
		StudentID:      st.ID,
		ConceptID:      "a",
		AttemptID:      "sess:attempt",
		Rating:         3,
		Correct:        true,
		ElapsedSeconds: 4,
		ReviewedAt:     now,
	}
	commit := ReviewSubmissionCommit{
		StudentID:    st.ID,
		SubmissionID: event.AttemptID,
		Event:        event,
		Attempt: AttemptEntry{
			SessionID: sess.ID,
			StudentID: st.ID, ConceptID: "a", Answer: "42", Expected: "42",
			Correct: true, ElapsedSeconds: 4, Timestamp: now, Source: "review",
		},
		Progress: progress,
		XPAmount: 5,
		XPReason: "study",
	}
	out, err := rs.CommitReviewSubmission(commit)
	if err != nil {
		t.Fatalf("commit review: %v", err)
	}
	if out.Replayed || out.AttemptID == 0 {
		t.Fatalf("first commit outcome = %+v, want a new attempt", out)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM review_events WHERE attempt_id = ?`, event.AttemptID); n != 1 {
		t.Fatalf("review events = %d, want 1", n)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM attempts WHERE student_id = ?`, st.ID); n != 1 {
		t.Fatalf("attempts = %d, want 1", n)
	}
	if total, _, err := store.GetXP(st.ID); err != nil || total != 5 {
		t.Fatalf("xp total = %d, err = %v, want 5", total, err)
	}
	got, err := store.GetProgress(st.ID, "a")
	if err != nil || got == nil {
		t.Fatalf("progress after commit: %+v, err = %v", got, err)
	}

	// A retry with the same attempt key must not write a second event, attempt
	// or XP award. This is the idempotency contract the engine relies on.
	replay := commit
	replay.Attempt.Answer = "different"
	if out, err := rs.CommitReviewSubmission(replay); err != nil {
		t.Fatalf("replay review: %v", err)
	} else if !out.Replayed {
		t.Fatalf("replay outcome = %+v, want Replayed", out)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM review_events WHERE attempt_id = ?`, event.AttemptID); n != 1 {
		t.Fatalf("review events after replay = %d, want 1", n)
	}
	if n := seedInt(t, store, `SELECT COUNT(*) FROM attempts WHERE student_id = ?`, st.ID); n != 1 {
		t.Fatalf("attempts after replay = %d, want 1", n)
	}
	if total, _, err := store.GetXP(st.ID); err != nil || total != 5 {
		t.Fatalf("xp total after replay = %d, err = %v, want 5", total, err)
	}
}

func TestReviewSubmissionRollsBackWhenTheAttemptCannotBeWritten(t *testing.T) {
	sq, ok := newTestStore(t).(*SQLiteStore)
	if !ok {
		t.Skip("failure trigger is SQLite-specific; the Postgres suite exercises the same transaction")
	}
	st := seedStudent(t, sq)
	if _, err := sq.db.Exec(`CREATE TRIGGER fail_review_attempt BEFORE INSERT ON attempts
		BEGIN SELECT RAISE(ABORT, 'injected attempt failure'); END`); err != nil {
		t.Fatalf("create failure trigger: %v", err)
	}
	t.Cleanup(func() { _, _ = sq.db.Exec(`DROP TRIGGER IF EXISTS fail_review_attempt`) })

	now := time.Now().UTC()
	event := ReviewEvent{
		StudentID: st.ID, ConceptID: "a", AttemptID: "sess:rollback",
		Rating: 3, Correct: true, ElapsedSeconds: 4, ReviewedAt: now,
	}
	_, err := sq.CommitReviewSubmission(ReviewSubmissionCommit{
		StudentID:    st.ID,
		SubmissionID: event.AttemptID,
		Event:        event,
		Attempt: AttemptEntry{
			StudentID: st.ID, ConceptID: "a", Answer: "42", Expected: "42",
			Correct: true, ElapsedSeconds: 4, Timestamp: now, Source: "review",
		},
		Progress: reviewProgress(st.ID),
		XPAmount: 5,
		XPReason: "study",
	})
	if err == nil {
		t.Fatal("commit succeeded despite the injected attempt failure")
	}
	if n := seedInt(t, sq, `SELECT COUNT(*) FROM review_events WHERE attempt_id = ?`, event.AttemptID); n != 0 {
		t.Fatalf("review events after rollback = %d, want 0", n)
	}
	if n := seedInt(t, sq, `SELECT COUNT(*) FROM attempts WHERE student_id = ?`, st.ID); n != 0 {
		t.Fatalf("attempts after rollback = %d, want 0", n)
	}
	if p, _ := sq.GetProgress(st.ID, "a"); p != nil {
		t.Fatalf("progress after rollback = %+v, want nil", p)
	}
	if total, _, err := sq.GetXP(st.ID); err != nil || total != 0 {
		t.Fatalf("xp total after rollback = %d, err = %v, want 0", total, err)
	}
}
