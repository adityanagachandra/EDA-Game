package engine

import (
	"errors"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
)

// reviewFailureStore delegates every repository method to the real store but
// fails the atomic review commit. It models a database error after the answer
// has been graded, which is the point where the old code had already committed
// FSRS state and could lose its review event.
type reviewFailureStore struct {
	storage.Repository
	err error
}

func (s reviewFailureStore) CommitReviewSubmission(storage.ReviewSubmissionCommit) (storage.ReviewSubmissionOutcome, error) {
	return storage.ReviewSubmissionOutcome{}, s.err
}

func masteredReviewRow(studentID string, due time.Time) *storage.ConceptProgress {
	stability := 100.0
	difficulty := 5.0
	lastReview := due.AddDate(0, 0, -10)
	return &storage.ConceptProgress{
		StudentID:       studentID,
		ConceptID:       "a",
		Status:          string(mastery.StatusMastered),
		Streak:          1,
		BestStreak:      5,
		AvgResponseTime: 6,
		Attempts:        6,
		LastReviewed:    &lastReview,
		NextReviewDue:   &due,
		SM2Repetitions:  4,
		SM2Interval:     100,
		SM2EFactor:      2.5,
		FSRSStability:   &stability,
		FSRSDifficulty:  &difficulty,
		FSRSLastReview:  &lastReview,
	}
}

func TestReviewCommitFailureLeavesProgressAttemptsAndXPNil(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("review-commit-failure")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	p := masteredReviewRow(st.ID, now.Add(-time.Hour))
	if err := e.repo.UpsertProgress(p); err != nil {
		t.Fatalf("seed progress: %v", err)
	}

	sess, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	q, err := e.NextReviewQuestion(sess.ID, st.ID)
	if err != nil || q == nil {
		t.Fatalf("next review: %v, %v", q, err)
	}

	before, err := e.repo.GetProgress(st.ID, "a")
	if err != nil || before == nil {
		t.Fatalf("read progress before: %v", err)
	}
	xpBefore, _, err := e.repo.GetXP(st.ID)
	if err != nil {
		t.Fatalf("read xp before: %v", err)
	}

	injected := errors.New("injected review commit failure")
	e.repo = reviewFailureStore{Repository: e.repo, err: injected}

	if _, err := e.SubmitAnswer(sess.ID, st.ID, q.AttemptID, "42", 4); !errors.Is(err, injected) {
		t.Fatalf("SubmitAnswer error = %v, want injected failure", err)
	}

	after, err := e.repo.GetProgress(st.ID, "a")
	if err != nil || after == nil {
		t.Fatalf("read progress after: %v", err)
	}
	if after.FSRSStability == nil || *after.FSRSStability != *before.FSRSStability {
		t.Errorf("stability changed after failed commit: before %v after %v", before.FSRSStability, after.FSRSStability)
	}
	if after.NextReviewDue == nil || !after.NextReviewDue.Equal(*before.NextReviewDue) {
		t.Errorf("due date changed after failed commit: before %v after %v", before.NextReviewDue, after.NextReviewDue)
	}
	attempts, err := e.repo.GetAttemptsForStudent(st.ID)
	if err != nil {
		t.Fatalf("read attempts: %v", err)
	}
	if len(attempts) != 0 {
		t.Errorf("attempts after failed commit = %d, want 0", len(attempts))
	}
	xpAfter, _, err := e.repo.GetXP(st.ID)
	if err != nil {
		t.Fatalf("read xp after: %v", err)
	}
	if xpAfter != xpBefore {
		t.Errorf("xp after failed commit = %d, want %d", xpAfter, xpBefore)
	}
}

func TestNonReviewAttemptsDoNotMoveMasteredFSRSState(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("fsrs-non-review")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	p := masteredReviewRow(st.ID, now.AddDate(0, 0, 20))
	beforeStability := *p.FSRSStability
	beforeDifficulty := *p.FSRSDifficulty
	beforeDue := *p.NextReviewDue
	beforeLastReview := *p.FSRSLastReview

	if _, _, _, _, err := e.computeGradedAttempt(p, gradedAttempt{
		studentID: st.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 40, correct: true, review: false,
	}); err != nil {
		t.Fatalf("compute non-review attempt: %v", err)
	}

	if *p.FSRSStability != beforeStability || *p.FSRSDifficulty != beforeDifficulty {
		t.Errorf("non-review changed FSRS state: stability %v -> %v, difficulty %v -> %v",
			beforeStability, *p.FSRSStability, beforeDifficulty, *p.FSRSDifficulty)
	}
	if !p.NextReviewDue.Equal(beforeDue) {
		t.Errorf("non-review changed due date: %v -> %v", beforeDue, *p.NextReviewDue)
	}
	if !p.FSRSLastReview.Equal(beforeLastReview) {
		t.Errorf("non-review changed FSRS last review: %v -> %v", beforeLastReview, *p.FSRSLastReview)
	}
	if !p.LastReviewed.Equal(beforeLastReview) {
		t.Errorf("non-review changed LastReviewed: %v -> %v", beforeLastReview, *p.LastReviewed)
	}
}

func TestMasteredMissingFSRSIsNotInitializedByNonReview(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("fsrs-missing")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	p := &storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: string(mastery.StatusMastered),
		Streak: 1, BestStreak: 5, AvgResponseTime: 6, Attempts: 6,
		SM2Repetitions: 4, SM2Interval: 100, SM2EFactor: 2.5,
	}
	if err := e.repo.UpsertProgress(p); err != nil {
		t.Fatalf("seed progress: %v", err)
	}

	if _, _, _, _, err := e.computeGradedAttempt(p, gradedAttempt{
		studentID: st.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 40, correct: true, review: false,
	}); err != nil {
		t.Fatalf("compute non-review attempt: %v", err)
	}
	if p.FSRSStability != nil || p.FSRSDifficulty != nil || p.FSRSLastReview != nil || p.NextReviewDue != nil {
		t.Fatalf("non-review initialized missing FSRS: stability=%v difficulty=%v lastReview=%v due=%v",
			p.FSRSStability, p.FSRSDifficulty, p.FSRSLastReview, p.NextReviewDue)
	}
	if err := e.repo.UpsertProgress(p); err != nil {
		t.Fatalf("store progress: %v", err)
	}
	due, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("ReviewsDue: %v", err)
	}
	if due != 0 {
		t.Errorf("ReviewsDue after non-review missing FSRS = %d, want 0", due)
	}

	if _, _, _, _, err := e.computeGradedAttempt(p, gradedAttempt{
		studentID: st.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 40, correct: true, review: true,
	}); err != nil {
		t.Fatalf("compute review attempt: %v", err)
	}
	if p.FSRSStability == nil || p.FSRSDifficulty == nil || p.FSRSLastReview == nil || p.NextReviewDue == nil {
		t.Fatalf("review did not initialize FSRS: stability=%v difficulty=%v lastReview=%v due=%v",
			p.FSRSStability, p.FSRSDifficulty, p.FSRSLastReview, p.NextReviewDue)
	}
	if *p.FSRSStability <= 0 {
		t.Errorf("initialized stability = %v, want > 0", *p.FSRSStability)
	}
}

func TestReviewAttemptFillsPartiallyPopulatedFSRS(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("fsrs-partial")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	stability := 50.0
	due := now.AddDate(0, 0, 10)
	p := &storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "a", Status: string(mastery.StatusMastered),
		Streak: 1, BestStreak: 5, AvgResponseTime: 6, Attempts: 6,
		LastReviewed: &now, NextReviewDue: &due,
		SM2Repetitions: 4, SM2Interval: 100, SM2EFactor: 2.5,
		FSRSStability: &stability,
	}

	if _, _, _, _, err := e.computeGradedAttempt(p, gradedAttempt{
		studentID: st.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 40, correct: true, review: true,
	}); err != nil {
		t.Fatalf("compute review attempt: %v", err)
	}
	if p.FSRSDifficulty == nil {
		t.Fatal("review did not fill missing difficulty")
	}
	if p.FSRSStability == nil || *p.FSRSStability <= 50 {
		t.Errorf("partial FSRS stability = %v, want growth from 50", p.FSRSStability)
	}
	if p.NextReviewDue == nil || !p.NextReviewDue.After(now) {
		t.Errorf("partial FSRS due = %v, want a future date", p.NextReviewDue)
	}
}

func TestReviewsDueUnaffectedByNonReviewOnValidSchedule(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("fsrs-bulk")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	now := time.Now().UTC()
	future := now.AddDate(0, 0, 20)
	past := now.Add(-time.Hour)

	valid := masteredReviewRow(st.ID, future)
	if err := e.repo.UpsertProgress(valid); err != nil {
		t.Fatalf("seed valid progress: %v", err)
	}
	missing := &storage.ConceptProgress{
		StudentID: st.ID, ConceptID: "b", Status: string(mastery.StatusMastered),
		Streak: 1, BestStreak: 5, AvgResponseTime: 6, Attempts: 6,
		SM2Repetitions: 4, SM2Interval: 100, SM2EFactor: 2.5,
	}
	if err := e.repo.UpsertProgress(missing); err != nil {
		t.Fatalf("seed missing progress: %v", err)
	}
	dueMissing := *missing
	dueMissing.ConceptID = "b"
	dueMissing.NextReviewDue = &past
	if err := e.repo.UpsertProgress(&dueMissing); err != nil {
		t.Fatalf("seed due missing progress: %v", err)
	}

	before, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("ReviewsDue before: %v", err)
	}
	if before != 1 {
		t.Fatalf("ReviewsDue before = %d, want 1 (valid future is not due, missing due is)", before)
	}

	if _, _, _, _, err := e.computeGradedAttempt(valid, gradedAttempt{
		studentID: st.ID, conceptID: "a", requiredStreak: 3,
		timeThreshold: 60, elapsed: 40, correct: true, review: false,
	}); err != nil {
		t.Fatalf("compute valid non-review: %v", err)
	}
	if err := e.repo.UpsertProgress(valid); err != nil {
		t.Fatalf("store valid progress: %v", err)
	}
	after, err := e.ReviewsDue(st.ID)
	if err != nil {
		t.Fatalf("ReviewsDue after: %v", err)
	}
	if after != before {
		t.Errorf("ReviewsDue changed from non-review: %d -> %d", before, after)
	}
}
