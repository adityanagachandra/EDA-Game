package engine

import (
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func commitOutboxEffects(t *testing.T, store storage.Repository, studentID, sessionID, submissionID string, effects []storage.StudyEffectInput) {
	t.Helper()
	sub, ok := store.(storage.StudySubmissionStore)
	if !ok {
		t.Fatal("store does not implement StudySubmissionStore")
	}
	now := time.Now().UTC()
	_, err := sub.CommitStudySubmission(storage.StudySubmissionCommit{
		StudentID:    studentID,
		SubmissionID: submissionID,
		Operation:    "study",
		RequestHash:  "hash-" + submissionID,
		ResponseJSON: `{"correct":false}`,
		Attempt: storage.AttemptEntry{
			SessionID:      sessionID,
			StudentID:      studentID,
			ConceptID:      "a",
			Answer:         "wrong",
			Expected:       "42",
			Correct:        false,
			ElapsedSeconds: 5,
			Timestamp:      now,
			Question:       "Q",
			Source:         "practice",
		},
		Effects: effects,
	})
	if err != nil {
		t.Fatalf("commit effects: %v", err)
	}
}

func TestOutbox_ProcessesCorrectAnswerEffects(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("outbox-correct")
	if err != nil {
		t.Fatal(err)
	}
	e.SetStudyAnchor(st.ID, "a", "42")
	if _, err := e.SubmitStudyAnswerWithSubmission(st.ID, "a", "42", 5, "Q", "outbox-correct"); err != nil {
		t.Fatalf("answer: %v", err)
	}
	effects := e.repo.(storage.StudyEffectStore)
	if err := e.ProcessStudyEffects(st.ID, 50); err != nil {
		t.Fatal(err)
	}
	pending, err := effects.ClaimStudyEffects(st.ID, time.Now().UTC().Add(time.Second), 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 0 {
		t.Fatalf("answer left %d pending effects", len(pending))
	}
	progress, err := e.repo.GetProgress(st.ID, "a")
	if err != nil || progress == nil {
		t.Fatalf("progress = %+v, %v", progress, err)
	}
	if progress.Attempts != 1 {
		t.Fatalf("attempts = %d, want 1", progress.Attempts)
	}
}

func TestOutbox_SurvivesRestartAndAppliesProvisionalDemotion(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("outbox-restart")
	if err != nil {
		t.Fatal(err)
	}
	session, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := e.repo.UpsertProgress(&storage.ConceptProgress{
		StudentID:       st.ID,
		ConceptID:       "a",
		Status:          "MASTERED",
		Streak:          3,
		Attempts:        3,
		AvgResponseTime: 5,
		WeaknessScore:   0,
		IsProvisional:   true,
		SM2EFactor:      2.5,
	}); err != nil {
		t.Fatal(err)
	}
	commitOutboxEffects(t, e.repo, st.ID, session.ID, "outbox-restart",
		[]storage.StudyEffectInput{
			{Effect: "remedial", PayloadJSON: `{"concept_ids":["a"]}`},
			{Effect: "provisional_fallback", PayloadJSON: `{"concept_id":"b"}`},
		})

	// A new process is a new engine over the same durable store.
	restarted := New(e.repo, e.dag, e.registry, nil, nil)
	if err := restarted.ProcessStudyEffects(st.ID, 50); err != nil {
		t.Fatalf("process: %v", err)
	}
	queue := restarted.QuizRemedial(st.ID)
	found := false
	for _, id := range queue {
		if id == "a" {
			found = true
		}
	}
	if !found {
		t.Fatalf("remedial queue = %v, want a", queue)
	}
	prereq, err := restarted.repo.GetProgress(st.ID, "a")
	if err != nil || prereq == nil {
		t.Fatalf("prereq progress = %+v, %v", prereq, err)
	}
	if prereq.Status != "LEARNING" || prereq.WeaknessScore < 0.79 {
		t.Fatalf("provisional prerequisite = %+v, want LEARNING weakness 0.8", prereq)
	}
	pending, err := restarted.repo.(storage.StudyEffectStore).ClaimStudyEffects(st.ID, time.Now().UTC().Add(time.Second), 50)
	if err != nil || len(pending) != 0 {
		t.Fatalf("pending after restart = %+v, %v", pending, err)
	}
}

func TestOutbox_UnknownEffectRetries(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("outbox-unknown")
	if err != nil {
		t.Fatal(err)
	}
	session, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatal(err)
	}
	commitOutboxEffects(t, e.repo, st.ID, session.ID, "outbox-unknown",
		[]storage.StudyEffectInput{{Effect: "not_a_real_effect"}})
	if err := e.ProcessStudyEffects(st.ID, 50); err != nil {
		t.Fatalf("process: %v", err)
	}
	effects := e.repo.(storage.StudyEffectStore)
	retry, err := effects.ClaimStudyEffects(st.ID, time.Now().UTC().Add(3*time.Second), 50)
	if err != nil || len(retry) != 1 {
		t.Fatalf("retry = %+v, %v", retry, err)
	}
	if retry[0].Attempts != 2 {
		t.Fatalf("attempts = %d, want 2", retry[0].Attempts)
	}
}

func TestOutbox_WeaknessEffectDoesNotRegressNewerProgress(t *testing.T) {
	e := testEngine(t)
	st, err := e.CreateStudent("outbox-weakness")
	if err != nil {
		t.Fatal(err)
	}
	session, err := e.repo.CreateSession(st.ID)
	if err != nil {
		t.Fatal(err)
	}
	commitOutboxEffects(t, e.repo, st.ID, session.ID, "outbox-weakness",
		[]storage.StudyEffectInput{{Effect: "propagate_weakness"}})
	for _, p := range []*storage.ConceptProgress{
		{StudentID: st.ID, ConceptID: "a", Status: "PRACTICING", Attempts: 1, WeaknessScore: 0.9, SM2EFactor: 2.5},
		{StudentID: st.ID, ConceptID: "b", Status: "PRACTICING", Attempts: 1, WeaknessScore: 0.5, SM2EFactor: 2.5},
	} {
		if err := e.repo.UpsertProgress(p); err != nil {
			t.Fatal(err)
		}
	}
	if err := e.ProcessStudyEffects(st.ID, 10); err != nil {
		t.Fatal(err)
	}
	b, err := e.repo.GetProgress(st.ID, "b")
	if err != nil || b == nil {
		t.Fatalf("b = %+v, %v", b, err)
	}
	if b.WeaknessScore != 0.5 {
		t.Fatalf("newer weakness regressed to %v, want 0.5", b.WeaknessScore)
	}
}
