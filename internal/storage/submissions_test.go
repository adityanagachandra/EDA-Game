package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func studySubmissionTestCommit(studentID, submissionID, answer string) StudySubmissionCommit {
	return StudySubmissionCommit{
		StudentID:    studentID,
		SubmissionID: submissionID,
		Operation:    "study",
		RequestHash:  "hash-" + answer,
		ResponseJSON: `{"correct":true}`,
		Attempt: AttemptEntry{
			StudentID:      studentID,
			ConceptID:      "a",
			Answer:         answer,
			Expected:       "4",
			Correct:        true,
			ElapsedSeconds: 5,
			Timestamp:      time.Now().UTC(),
			Question:       "2+2",
			Source:         "practice",
		},
		Progress: &ConceptProgress{
			StudentID:       studentID,
			ConceptID:       "a",
			Status:          "LEARNING",
			Attempts:        1,
			AvgResponseTime: 5,
			SM2EFactor:      2.5,
		},
		TopicSpeed: &TopicSpeed{
			StudentID:     studentID,
			ConceptID:     "a",
			EFactor:       2.5,
			LearningSpeed: 1,
		},
		XPAmount: 1,
		XPReason: "answer",
	}
}

func studySubmissionStore(t testing.TB, store Repository) StudySubmissionStore {
	t.Helper()
	sub, ok := store.(StudySubmissionStore)
	if !ok {
		t.Fatalf("%T does not implement StudySubmissionStore", store)
	}
	return sub
}

func TestStudySubmission_ReplayAndPayloadMismatch(t *testing.T) {
	store := newTestStore(t)
	sub := studySubmissionStore(t, store)
	st, err := store.CreateStudent("submission")
	if err != nil {
		t.Fatal(err)
	}
	session, err := store.CreateSession(st.ID)
	if err != nil {
		t.Fatal(err)
	}
	commit := studySubmissionTestCommit(st.ID, "sub-1", "4")
	commit.Attempt.SessionID = session.ID

	first, err := sub.CommitStudySubmission(commit)
	if err != nil {
		t.Fatalf("first commit: %v", err)
	}
	if first.Replayed {
		t.Fatal("first commit reported replay")
	}
	second, err := sub.CommitStudySubmission(commit)
	if err != nil {
		t.Fatalf("second commit: %v", err)
	}
	if !second.Replayed || second.ResponseJSON != first.ResponseJSON {
		t.Fatalf("expected exact replay, got %+v", second)
	}

	attempts, _ := store.GetRecentAttemptsForConcept(st.ID, "a", 10)
	if len(attempts) != 1 {
		t.Fatalf("expected one attempt, got %d", len(attempts))
	}
	student, _ := store.GetStudent(st.ID)
	if student.XPTotal != 1 {
		t.Fatalf("expected one XP, got %d", student.XPTotal)
	}

	wrong := commit
	wrong.RequestHash = "hash-wrong"
	if _, err := sub.CommitStudySubmission(wrong); err != ErrStudySubmissionPayloadMismatch {
		t.Fatalf("expected payload mismatch, got %v", err)
	}
	attempts, _ = store.GetRecentAttemptsForConcept(st.ID, "a", 10)
	if len(attempts) != 1 {
		t.Fatalf("mismatch wrote an attempt: %d", len(attempts))
	}
	student, _ = store.GetStudent(st.ID)
	if student.XPTotal != 1 {
		t.Fatalf("mismatch moved XP to %d", student.XPTotal)
	}

	other := commit
	other.SubmissionID = "sub-2"
	out, err := sub.CommitStudySubmission(other)
	if err != nil {
		t.Fatalf("second logical commit: %v", err)
	}
	if out.Replayed {
		t.Fatal("a different submission id was treated as a replay")
	}
	attempts, _ = store.GetRecentAttemptsForConcept(st.ID, "a", 10)
	if len(attempts) != 2 {
		t.Fatalf("expected two attempts, got %d", len(attempts))
	}
	student, _ = store.GetStudent(st.ID)
	if student.XPTotal != 2 {
		t.Fatalf("expected two XP, got %d", student.XPTotal)
	}
}

func BenchmarkStudySubmissionCommit(b *testing.B) {
	store, err := NewSQLiteStore(":memory:")
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()
	st, err := store.CreateStudent("bench")
	if err != nil {
		b.Fatal(err)
	}
	session, err := store.CreateSession(st.ID)
	if err != nil {
		b.Fatal(err)
	}
	sub := studySubmissionStore(b, store)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		commit := studySubmissionTestCommit(st.ID, "bench-"+strconv.Itoa(i), "4")
		commit.Attempt.SessionID = session.ID
		if _, err := sub.CommitStudySubmission(commit); err != nil {
			b.Fatal(err)
		}
	}
}

func TestStudySubmission_RollbackOnInjectedFailures(t *testing.T) {
	store := newTestStore(t)
	sub := studySubmissionStore(t, store)
	st, _ := store.CreateStudent("rollback")
	session, _ := store.CreateSession(st.ID)

	dialect := "sqlite"
	var db *sql.DB
	switch s := store.(type) {
	case *SQLiteStore:
		db = s.db
	case *PostgresStore:
		db = s.db
		dialect = "postgres"
	default:
		t.Fatalf("unsupported store %T", store)
	}

	failure := func(name, target, timing string) (string, string) {
		if dialect == "sqlite" {
			trigger := "CREATE TRIGGER " + name + " " + timing + " ON " + target +
				" BEGIN SELECT RAISE(ABORT, 'injected " + name + " failure'); END"
			return trigger, "DROP TRIGGER " + name
		}
		fn := name + "_fn"
		function := "CREATE OR REPLACE FUNCTION " + fn + "() RETURNS trigger AS $$ BEGIN RAISE EXCEPTION 'injected " + name + " failure'; END; $$ LANGUAGE plpgsql"
		trigger := "CREATE TRIGGER " + name + " " + timing + " ON " + target +
			" FOR EACH ROW EXECUTE FUNCTION " + fn + "()"
		return function + "; " + trigger,
			"DROP TRIGGER IF EXISTS " + name + " ON " + tableForTrigger(target) + "; DROP FUNCTION IF EXISTS " + fn + "()"
	}

	cases := []struct {
		name   string
		target string
		timing string
	}{
		{"fail_progress", "concept_progress", "BEFORE INSERT"},
		{"fail_xp", "xp_events", "BEFORE INSERT"},
		{"fail_after_xp", "xp_events", "AFTER INSERT"},
	}
	for _, tc := range cases {
		t.Run(dialect+" "+tc.name, func(t *testing.T) {
			create, drop := failure(tc.name, tc.target, tc.timing)
			if _, err := db.Exec(create); err != nil {
				t.Fatalf("create trigger: %v", err)
			}
			defer func() { _, _ = db.Exec(drop) }()
			commit := studySubmissionTestCommit(st.ID, "rollback-"+tc.name, "4")
			commit.Attempt.SessionID = session.ID
			commit.Effects = []StudyEffectInput{{Effect: "propagate_weakness"}}
			if _, err := sub.CommitStudySubmission(commit); err == nil {
				t.Fatal("expected injected failure")
			}
			attempts, _ := store.GetRecentAttemptsForConcept(st.ID, "a", 10)
			if len(attempts) != 0 {
				t.Fatalf("rollback left %d attempts", len(attempts))
			}
			student, _ := store.GetStudent(st.ID)
			if student.XPTotal != 0 {
				t.Fatalf("rollback left XP %d", student.XPTotal)
			}
			progress, _ := store.GetProgress(st.ID, "a")
			if progress != nil {
				t.Fatalf("rollback left progress: %+v", progress)
			}
			saved, _ := sub.GetStudySubmission(st.ID, "study", "rollback-"+tc.name)
			if saved != nil {
				t.Fatalf("rollback left submission: %+v", saved)
			}
			pending, _ := studyEffectStore(t, store).ClaimStudyEffects(st.ID, time.Now().UTC().Add(time.Second), 10)
			if len(pending) != 0 {
				t.Fatalf("rollback left %d outbox rows", len(pending))
			}
		})
	}
}

func tableForTrigger(target string) string {
	return target
}

func studyEffectStore(t testing.TB, store Repository) StudyEffectStore {
	t.Helper()
	out, ok := store.(StudyEffectStore)
	if !ok {
		t.Fatalf("%T does not implement StudyEffectStore", store)
	}
	return out
}

// TestStudyEffects_ClaimCompleteAndFail covers the durable outbox lease:
// a claimed row cannot be claimed again, completion retires it, and a failure
// returns it with a later retry time.
func TestStudyEffects_ClaimCompleteAndFail(t *testing.T) {
	store := newTestStore(t)
	sub := studySubmissionStore(t, store)
	effects := studyEffectStore(t, store)
	st, _ := store.CreateStudent("outbox")
	session, _ := store.CreateSession(st.ID)
	now := time.Now().UTC().Add(time.Second)

	commit := studySubmissionTestCommit(st.ID, "outbox-1", "4")
	commit.Attempt.SessionID = session.ID
	commit.Effects = []StudyEffectInput{{Effect: "propagate_weakness"}}
	if _, err := sub.CommitStudySubmission(commit); err != nil {
		t.Fatalf("commit: %v", err)
	}
	first, err := effects.ClaimStudyEffects(st.ID, now, 10)
	if err != nil || len(first) != 1 {
		t.Fatalf("first claim = %+v, %v", first, err)
	}
	if first[0].Attempts != 1 || first[0].Effect != "propagate_weakness" {
		t.Fatalf("first claim row = %+v", first[0])
	}
	again, err := effects.ClaimStudyEffects(st.ID, now, 10)
	if err != nil || len(again) != 0 {
		t.Fatalf("claimed row was claimable again: %+v, %v", again, err)
	}
	if err := effects.CompleteStudyEffect(first[0].ID); err != nil {
		t.Fatalf("complete: %v", err)
	}
	done, err := effects.ClaimStudyEffects(st.ID, now, 10)
	if err != nil || len(done) != 0 {
		t.Fatalf("completed row was claimable: %+v, %v", done, err)
	}

	second := studySubmissionTestCommit(st.ID, "outbox-2", "4")
	second.Attempt.SessionID = session.ID
	second.Effects = []StudyEffectInput{{Effect: "implicit_credit", PayloadJSON: `{"concept_id":"a","attempt_id":"x","correct":true}`}}
	if _, err := sub.CommitStudySubmission(second); err != nil {
		t.Fatalf("second commit: %v", err)
	}
	claimed, err := effects.ClaimStudyEffects(st.ID, now, 10)
	if err != nil || len(claimed) != 1 {
		t.Fatalf("second claim = %+v, %v", claimed, err)
	}
	if err := effects.FailStudyEffect(claimed[0].ID, "transient", now.Add(-time.Second)); err != nil {
		t.Fatalf("fail: %v", err)
	}
	retried, err := effects.ClaimStudyEffects(st.ID, now, 10)
	if err != nil || len(retried) != 1 {
		t.Fatalf("failed row was not retryable: %+v, %v", retried, err)
	}
	if retried[0].Attempts != 2 {
		t.Fatalf("retry attempts = %d, want 2", retried[0].Attempts)
	}
}

// TestStudyEffects_PruneOnlyTouchesCompletedRows proves retention pruning
// removes retired rows but never a pending or retrying effect.
func TestStudyEffects_PruneOnlyTouchesCompletedRows(t *testing.T) {
	store := newTestStore(t)
	sub := studySubmissionStore(t, store)
	effects := studyEffectStore(t, store)
	st, _ := store.CreateStudent("prune")
	session, _ := store.CreateSession(st.ID)
	now := time.Now().UTC().Add(time.Second)

	commit := studySubmissionTestCommit(st.ID, "prune-1", "4")
	commit.Attempt.SessionID = session.ID
	commit.Effects = []StudyEffectInput{
		{Effect: "propagate_weakness"},
		{Effect: "implicit_credit", PayloadJSON: `{"concept_id":"a","attempt_id":"x","correct":true}`},
		{Effect: "remove_remedial", PayloadJSON: `{"concept_id":"a"}`},
	}
	if _, err := sub.CommitStudySubmission(commit); err != nil {
		t.Fatal(err)
	}
	rows, err := effects.ClaimStudyEffects(st.ID, now, 10)
	if err != nil || len(rows) != 3 {
		t.Fatalf("claim = %d, %v; want 3", len(rows), err)
	}
	for _, row := range rows {
		if err := effects.CompleteStudyEffect(row.ID); err != nil {
			t.Fatal(err)
		}
	}
	removed, err := effects.PruneStudyEffects(now.Add(time.Hour), 100)
	if err != nil || removed != 3 {
		t.Fatalf("prune completed = %d, %v; want 3", removed, err)
	}

	pending := studySubmissionTestCommit(st.ID, "prune-2", "4")
	pending.Attempt.SessionID = session.ID
	pending.Effects = []StudyEffectInput{{Effect: "propagate_weakness"}}
	if _, err := sub.CommitStudySubmission(pending); err != nil {
		t.Fatal(err)
	}
	removed, err = effects.PruneStudyEffects(now.Add(time.Hour), 100)
	if err != nil || removed != 0 {
		t.Fatalf("prune pending = %d, %v; want 0", removed, err)
	}
	left, err := effects.ClaimStudyEffects(st.ID, now, 10)
	if err != nil || len(left) != 1 {
		t.Fatalf("pending effect missing after prune: %d, %v", len(left), err)
	}
}

// TestStudyEffects_ExpiredLeaseIsRecoverable proves a crashed worker does not
// strand an effect: the lease blocks another claim, then expires and the row is
// claimable again with its attempt count preserved.
func TestStudyEffects_ExpiredLeaseIsRecoverable(t *testing.T) {
	store := newTestStore(t)
	sub := studySubmissionStore(t, store)
	effects := studyEffectStore(t, store)
	st, _ := store.CreateStudent("lease")
	session, _ := store.CreateSession(st.ID)
	commit := studySubmissionTestCommit(st.ID, "lease-1", "4")
	commit.Attempt.SessionID = session.ID
	commit.Effects = []StudyEffectInput{{Effect: "propagate_weakness"}}
	if _, err := sub.CommitStudySubmission(commit); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Add(time.Second)
	first, err := effects.ClaimStudyEffects(st.ID, now, 1)
	if err != nil || len(first) != 1 {
		t.Fatalf("first claim = %+v, %v", first, err)
	}
	blocked, err := effects.ClaimStudyEffects(st.ID, now.Add(30*time.Second), 1)
	if err != nil || len(blocked) != 0 {
		t.Fatalf("live lease was claimable: %+v, %v", blocked, err)
	}
	recovered, err := effects.ClaimStudyEffects(st.ID, now.Add(61*time.Second), 1)
	if err != nil || len(recovered) != 1 {
		t.Fatalf("expired lease was not claimable: %+v, %v", recovered, err)
	}
	if recovered[0].Attempts != 2 {
		t.Fatalf("attempts = %d, want 2", recovered[0].Attempts)
	}
}

// TestStudyEffects_ConcurrentClaimsAreDisjoint runs two workers against the
// same pending set. A row may be claimed once, never by both.
func newConcurrentStore(t *testing.T) Repository {
	t.Helper()
	if os.Getenv("TEST_POSTGRES_DSN") != "" {
		return newTestStore(t)
	}
	store, err := NewSQLiteStore(filepath.Join(t.TempDir(), "concurrent.db"))
	if err != nil {
		t.Fatalf("create file sqlite store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestStudyEffects_ConcurrentClaimsAreDisjoint(t *testing.T) {
	store := newConcurrentStore(t)
	sub := studySubmissionStore(t, store)
	effects := studyEffectStore(t, store)
	st, _ := store.CreateStudent("claims")
	session, _ := store.CreateSession(st.ID)
	commit := studySubmissionTestCommit(st.ID, "claims-1", "4")
	commit.Attempt.SessionID = session.ID
	for i := 0; i < 20; i++ {
		commit.Effects = append(commit.Effects, StudyEffectInput{Effect: fmt.Sprintf("effect-%02d", i)})
	}
	if _, err := sub.CommitStudySubmission(commit); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC().Add(time.Second)
	type result struct {
		rows []StudyAnswerEffect
		err  error
	}
	ch := make(chan result, 2)
	for i := 0; i < 2; i++ {
		go func() {
			rows, err := effects.ClaimStudyEffects(st.ID, now, 20)
			ch <- result{rows, err}
		}()
	}
	seen := map[int64]bool{}
	for i := 0; i < 2; i++ {
		r := <-ch
		if r.err != nil {
			t.Fatal(r.err)
		}
		for _, row := range r.rows {
			if seen[row.ID] {
				t.Fatalf("effect %d claimed twice", row.ID)
			}
			seen[row.ID] = true
		}
	}
	if len(seen) != 20 {
		t.Fatalf("claimed %d of 20 effects", len(seen))
	}
}
