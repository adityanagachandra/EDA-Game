package storage

import (
	"database/sql"
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
		})
	}
}

func tableForTrigger(target string) string {
	return target
}
