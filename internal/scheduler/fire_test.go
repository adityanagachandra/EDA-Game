package scheduler

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/mastery"
)

func fireDAG(t *testing.T) *concepts.DAG {
	t.Helper()
	d, err := concepts.Build([]concepts.Concept{
		{
			ID:               "prereq",
			Label:            "Prerequisite Topic",
			Domain:           "algebra",
			Subdomain:        "foundations",
			Prerequisites:    []string{},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
		{
			ID:               "adv",
			Label:            "Advanced Topic",
			Domain:           "calculus",
			Subdomain:        "derivatives",
			Prerequisites:    []string{"prereq"},
			Encompasses:      []string{"prereq"},
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 10},
		},
	})
	if err != nil {
		t.Fatalf("build fire DAG: %v", err)
	}
	return d
}

func dueSnapshots(now time.Time, banked int) map[string]*ConceptSnapshot {
	past := now.AddDate(0, 0, -5)
	return map[string]*ConceptSnapshot{
		"prereq": {
			Status:         mastery.StatusMastered,
			LastReviewed:   past,
			NextReviewDue:  &past,
			RequiredStreak: 3,
			Interval:       5,
			EFactor:        2.5,
			BankedOrdinals: banked,
			Encompassed:    true,
		},
		"adv": {
			Status:         mastery.StatusLearning,
			RequiredStreak: 3,
		},
	}
}

func TestFIRe_BankedOrdinalDoesNotSuppressReview(t *testing.T) {
	d := fireDAG(t)
	s := New(d)
	now := time.Now().UTC()
	snap := dueSnapshots(now, 3)
	due := *snap["prereq"].NextReviewDue

	review := s.NextReview(snap, "")
	if review == nil || review.Concept.ID != "prereq" {
		t.Fatalf("banked credit suppressed the review: got %+v", review)
	}
	prereq := snap["prereq"]
	if prereq.BankedOrdinals != 3 {
		t.Fatalf("banked ordinals changed to %d, want 3", prereq.BankedOrdinals)
	}
	if prereq.NextReviewDue == nil || !prereq.NextReviewDue.Equal(due) {
		t.Fatalf("banked credit moved the due date: got %v, want %v", prereq.NextReviewDue, due)
	}
	if prereq.Interval != 5 {
		t.Fatalf("banked credit moved the interval: got %d, want 5", prereq.Interval)
	}
}

func TestFIRe_NextSmartKeepsTheReview(t *testing.T) {
	d := fireDAG(t)
	s := New(d)
	now := time.Now().UTC()
	snap := dueSnapshots(now, 1)

	cands := s.NextSmart(snap, nil, 10, 0, nil)
	for _, c := range cands {
		if c.Concept.ID == "prereq" && c.IsReview {
			return
		}
	}
	t.Fatalf("banked credit removed the prereq review from NextSmart: %+v", cands)
}

func TestFIRe_BankedOrdinalDoesNotChangePriority(t *testing.T) {
	now := time.Now().UTC()
	past := now.AddDate(0, 0, -15)
	with := &ConceptSnapshot{
		Status:         mastery.StatusMastered,
		LastReviewed:   past,
		LastAttempted:  past,
		NextReviewDue:  &past,
		RequiredStreak: 3,
		BankedOrdinals: 1,
		Encompassed:    true,
	}
	without := &ConceptSnapshot{
		Status:         mastery.StatusMastered,
		LastReviewed:   past,
		LastAttempted:  past,
		NextReviewDue:  &past,
		RequiredStreak: 3,
		Encompassed:    true,
	}
	if a, b := computePriority(with, now, true), computePriority(without, now, true); a != b {
		t.Fatalf("banked credit changed priority: with %f, without %f", a, b)
	}
}

func TestFIRe_EligibilityIgnoresBankedOrdinal(t *testing.T) {
	d := fireDAG(t)
	now := time.Now().UTC()
	snap := dueSnapshots(now, 2)
	el := EligibilityOf(snap, now, d.Concept("prereq"))
	if el.State != StateDueForReview {
		t.Fatalf("eligibility = %s, want %s", el.State, StateDueForReview)
	}
}

func TestFIRe_ChainedOrdinalDoesNotAdvance(t *testing.T) {
	d := fireDAG(t)
	s := New(d)
	now := time.Now().UTC()
	snap := dueSnapshots(now, 2)
	due := *snap["prereq"].NextReviewDue

	s.NextReview(snap, "")
	s.NextReview(snap, "")
	prereq := snap["prereq"]
	if prereq.BankedOrdinals != 2 {
		t.Fatalf("banked ordinals changed to %d, want 2", prereq.BankedOrdinals)
	}
	if prereq.NextReviewDue == nil || !prereq.NextReviewDue.Equal(due) {
		t.Fatalf("banked credit moved the due date: got %v, want %v", prereq.NextReviewDue, due)
	}
}
