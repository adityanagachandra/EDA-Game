package engine

import (
	"errors"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/storage"
)

// TestStudySubmissionDoesNotPublishKPStateOnRejectedCommit proves the scaffold
// state is computed on a copy. The first submission advances the clone and
// publishes only after it commits; the mismatched retry would reset the clone,
// but its commit is refused, so the live state must stay at the first answer.
func TestStudySubmissionDoesNotPublishKPStateOnRejectedCommit(t *testing.T) {
	d, err := concepts.Build([]concepts.Concept{{
		ID: "a", Label: "A", Domain: "d", GradingType: "numeric",
		Prerequisites:    []string{},
		MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		KnowledgePoints: []concepts.KnowledgePoint{
			{Ordinal: 1, Title: "First"},
			{Ordinal: 2, Title: "Second"},
		},
	}})
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	defer store.Close()
	reg := generator.NewRegistry()
	reg.Register("a", &testGen{answer: "42"})
	e := New(store, d, reg, nil, nil)
	st, err := e.CreateStudent("kp-submission")
	if err != nil {
		t.Fatalf("student: %v", err)
	}
	e.SetStudyAnchor(st.ID, "a", "42")

	if _, err := e.SubmitStudyAnswerWithSubmission(st.ID, "a", "42", 5, "Q1", "sub-1"); err != nil {
		t.Fatalf("first answer: %v", err)
	}
	kp := e.GetKPState(st.ID, "a")
	if kp == nil || kp.CurrentStreak != 1 || kp.CurrentOrdinal != 1 {
		t.Fatalf("first answer KP = %+v, want streak 1 at ordinal 1", kp)
	}

	_, err = e.SubmitStudyAnswerWithSubmission(st.ID, "a", "wrong", 5, "Q1", "sub-1")
	if !errors.Is(err, storage.ErrStudySubmissionPayloadMismatch) {
		t.Fatalf("mismatched retry = %v, want payload mismatch", err)
	}
	after := e.GetKPState(st.ID, "a")
	if after == nil || after.CurrentStreak != 1 || after.CurrentOrdinal != 1 {
		t.Fatalf("rejected commit changed live KP to %+v", after)
	}
}
