package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
)

// quizSelectionStores runs the selection rule on both stores. The candidate
// query is a progress read and the recency signal is an attempt read, so a
// store that answers either differently changes who is assessed.
func quizSelectionStores(t *testing.T) map[string]storage.Repository {
	t.Helper()
	sq, err := storage.NewSQLiteStore(filepath.Join(t.TempDir(), "quiz-selection.db"))
	if err != nil {
		t.Fatalf("create sqlite store: %v", err)
	}
	t.Cleanup(func() { sq.Close() })
	out := map[string]storage.Repository{"sqlite": sq}
	if dsn := os.Getenv("TEST_POSTGRES_DSN"); dsn != "" {
		pg, err := storage.NewPostgresStore(dsn)
		if err != nil {
			t.Fatalf("open postgres: %v", err)
		}
		t.Cleanup(func() { pg.Close() })
		out["postgres"] = pg
	} else {
		t.Log("TEST_POSTGRES_DSN not set; Postgres leg skipped")
	}
	return out
}

func quizSelectionEngine(t *testing.T, store storage.Repository, cs []concepts.Concept) *Engine {
	t.Helper()
	d, err := concepts.Build(cs)
	if err != nil {
		t.Fatalf("build DAG: %v", err)
	}
	return New(store, d, generator.NewRegistry(), nil, nil)
}

func quizSelectionStudent(t *testing.T, store storage.Repository, suffix string) string {
	t.Helper()
	st, err := store.CreateStudent(fmt.Sprintf("quizsel_%s_%d", suffix, time.Now().UnixNano()))
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	return st.ID
}

func quizSelectionConcept(id, subdomain string, prereqs ...string) concepts.Concept {
	return concepts.Concept{
		ID:               id,
		Label:            id,
		Domain:           "d",
		Subdomain:        subdomain,
		GradingType:      "numeric",
		Prerequisites:    prereqs,
		MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
	}
}

func quizSelectionProgress(t *testing.T, store storage.Repository, studentID, conceptID, status string, attempts int, weakness float64, placement bool) {
	t.Helper()
	now := time.Now().UTC()
	p := &storage.ConceptProgress{
		StudentID:       studentID,
		ConceptID:       conceptID,
		Status:          status,
		Attempts:        attempts,
		WeaknessScore:   weakness,
		PlacementSeeded: placement,
	}
	if attempts > 0 {
		p.LastAttempted = &now
	}
	if status == string(mastery.StatusMastered) {
		p.LastReviewed = &now
	}
	if err := store.UpsertProgress(p); err != nil {
		t.Fatalf("seed progress for %s: %v", conceptID, err)
	}
}

func quizSelectionIDs(cs []*concepts.Concept) []string {
	out := make([]string, 0, len(cs))
	for _, c := range cs {
		out = append(out, c.ID)
	}
	return out
}

func TestEngine_PickQuizConcepts_CandidateRule(t *testing.T) {
	for tag, store := range quizSelectionStores(t) {
		t.Run(tag, func(t *testing.T) {
			e := quizSelectionEngine(t, store, []concepts.Concept{
				quizSelectionConcept("a", "s1"),
				quizSelectionConcept("b", "s2", "a"),
				quizSelectionConcept("c", "s3"),
				quizSelectionConcept("locked", "s4", "b"),
				quizSelectionConcept("unseen", "s5"),
				quizSelectionConcept("placement", "s6"),
				quizSelectionConcept("unstarted", "s7"),
			})
			sid := quizSelectionStudent(t, store, "candidate")
			// a is mastered and deliberately weak: mastery excludes it from the
			// assessment pool, but it still satisfies b's prerequisite.
			quizSelectionProgress(t, store, sid, "a", string(mastery.StatusMastered), 3, 1.0, false)
			quizSelectionProgress(t, store, sid, "b", string(mastery.StatusPracticing), 2, 0.4, false)
			quizSelectionProgress(t, store, sid, "c", string(mastery.StatusPracticing), 2, 0.7, false)
			quizSelectionProgress(t, store, sid, "locked", string(mastery.StatusPracticing), 2, 0.9, false)
			quizSelectionProgress(t, store, sid, "unseen", string(mastery.StatusUnseen), 0, 0.9, false)
			quizSelectionProgress(t, store, sid, "placement", string(mastery.StatusPracticing), 0, 0.9, true)

			picked, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick: %v", err)
			}
			if got, want := quizSelectionIDs(picked), []string{"c", "b"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("picked %v, want %v", got, want)
			}
			for _, id := range quizSelectionIDs(picked) {
				if id == "a" {
					t.Error("a MASTERED concept was selected for a routine mastery check")
				}
			}

			// Selection is a read. It must not promote anything, write a row for
			// the unstarted concept or turn placement into evidence.
			progress, err := store.GetAllProgress(sid)
			if err != nil {
				t.Fatalf("reload progress: %v", err)
			}
			if len(progress) != 6 {
				t.Errorf("selection wrote progress: %d rows", len(progress))
			}
			if progress["locked"].Status != string(mastery.StatusPracticing) {
				t.Errorf("selection changed the locked concept: %+v", progress["locked"])
			}
			if progress["b"].Status != string(mastery.StatusPracticing) {
				t.Errorf("selection changed the assessed concept: %+v", progress["b"])
			}
			if progress["a"].Status != string(mastery.StatusMastered) {
				t.Errorf("selection changed the mastered prerequisite: %+v", progress["a"])
			}
		})
	}
}

func TestEngine_PickQuizConcepts_DecayedMasteryIsNotACandidate(t *testing.T) {
	for tag, store := range quizSelectionStores(t) {
		t.Run(tag, func(t *testing.T) {
			e := quizSelectionEngine(t, store, []concepts.Concept{
				quizSelectionConcept("a", "s1"),
				quizSelectionConcept("b", "s2"),
			})
			sid := quizSelectionStudent(t, store, "decayed")
			old := time.Now().UTC().Add(-30 * 24 * time.Hour)
			if err := store.UpsertProgress(&storage.ConceptProgress{
				StudentID:     sid,
				ConceptID:     "a",
				Status:        string(mastery.StatusMastered),
				Attempts:      3,
				WeaknessScore: 0.9,
				LastReviewed:  &old,
			}); err != nil {
				t.Fatalf("seed decayed mastery: %v", err)
			}
			quizSelectionProgress(t, store, sid, "b", string(mastery.StatusPracticing), 1, 0.2, false)

			picked, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick: %v", err)
			}
			if got, want := quizSelectionIDs(picked), []string{"b"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("picked %v, want %v; a decayed concept belongs to review, not this check", got, want)
			}
		})
	}
}

func TestEngine_PickQuizConcepts_WeakFirstThenHealthyAndDiverse(t *testing.T) {
	for tag, store := range quizSelectionStores(t) {
		t.Run(tag, func(t *testing.T) {
			e := quizSelectionEngine(t, store, []concepts.Concept{
				quizSelectionConcept("a", "s1"),
				quizSelectionConcept("b", "s2"),
				quizSelectionConcept("c", "s3"),
				quizSelectionConcept("d", "s4"),
				quizSelectionConcept("e", "s5"),
				quizSelectionConcept("f", "s6"),
				quizSelectionConcept("g", "s1"),
			})
			sid := quizSelectionStudent(t, store, "order")
			for id, weakness := range map[string]float64{
				"a": 0.9, "b": 0.8, "c": 0.7, "d": 0, "e": 0, "f": 0, "g": 0,
			} {
				quizSelectionProgress(t, store, sid, id, string(mastery.StatusPracticing), 1, weakness, false)
			}

			picked, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick: %v", err)
			}
			// The three weak concepts come first, then healthy concepts fill
			// the remaining places. g shares a subdomain with a, so diversity
			// keeps it off the list while d and e are available.
			if got, want := quizSelectionIDs(picked), []string{"a", "b", "c", "d", "e"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("picked %v, want %v", got, want)
			}
			seenSub := map[string]bool{}
			for _, c := range picked {
				if seenSub[c.Subdomain] {
					t.Errorf("duplicate subdomain %q in %v", c.Subdomain, quizSelectionIDs(picked))
				}
				seenSub[c.Subdomain] = true
			}
		})
	}
}

func TestEngine_HasQuizCandidatesRequiresAnAcquiringConcept(t *testing.T) {
	for tag, store := range quizSelectionStores(t) {
		t.Run(tag, func(t *testing.T) {
			cs := []concepts.Concept{
				quizSelectionConcept("a", "s1"),
				quizSelectionConcept("b", "s2"),
			}
			e := quizSelectionEngine(t, store, cs)
			sid := quizSelectionStudent(t, store, "ready")

			if ready, err := e.HasQuizCandidates(sid); err != nil || ready {
				t.Fatalf("empty progress: ready=%v err=%v, want false", ready, err)
			}

			quizSelectionProgress(t, store, sid, "a", string(mastery.StatusMastered), 3, 0, false)
			if ready, err := e.HasQuizCandidates(sid); err != nil || ready {
				t.Fatalf("mastered-only progress: ready=%v err=%v, want false", ready, err)
			}

			quizSelectionProgress(t, store, sid, "b", string(mastery.StatusPracticing), 1, 0.5, false)
			if ready, err := e.HasQuizCandidates(sid); err != nil || !ready {
				t.Fatalf("practicing progress: ready=%v err=%v, want true", ready, err)
			}
		})
	}
}

func TestEngine_PickQuizConcepts_AvoidsRecentQuizConcepts(t *testing.T) {
	for tag, store := range quizSelectionStores(t) {
		t.Run(tag, func(t *testing.T) {
			var cs []concepts.Concept
			for i := 0; i < 12; i++ {
				id := string(rune('a' + i))
				cs = append(cs, quizSelectionConcept(id, fmt.Sprintf("s%d", i)))
			}
			e := quizSelectionEngine(t, store, cs)
			sid := quizSelectionStudent(t, store, "recent")
			for _, c := range cs {
				quizSelectionProgress(t, store, sid, c.ID, string(mastery.StatusPracticing), 1, 0.5, false)
			}
			// The latest five distinct concepts are the recent window. The
			// practice attempt on f proves only quiz attempts count.
			sess, err := store.CreateSession(sid)
			if err != nil {
				t.Fatalf("create attempt session: %v", err)
			}
			base := time.Now().UTC().Add(-time.Hour)
			for i, id := range []string{"a", "b", "c", "d", "e"} {
				if _, err := store.RecordAttempt(storage.AttemptEntry{
					SessionID: sess.ID, StudentID: sid, ConceptID: id, Answer: "4", Expected: "4",
					Correct: true, Timestamp: base.Add(time.Duration(i) * time.Minute), Source: TaskQuiz,
				}); err != nil {
					t.Fatalf("record quiz attempt: %v", err)
				}
			}
			if _, err := store.RecordAttempt(storage.AttemptEntry{
				SessionID: sess.ID, StudentID: sid, ConceptID: "f", Answer: "4", Expected: "4",
				Correct: true, Timestamp: base.Add(2 * time.Hour), Source: "practice",
			}); err != nil {
				t.Fatalf("record practice attempt: %v", err)
			}

			picked, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick: %v", err)
			}
			if got, want := quizSelectionIDs(picked), []string{"f", "g", "h", "i", "j"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("picked %v, want five non-recent concepts %v", got, want)
			}
			for _, id := range quizSelectionIDs(picked) {
				if id >= "a" && id <= "e" {
					t.Errorf("recent concept %s repeated while alternatives existed", id)
				}
			}
		})
	}
}

func TestEngine_PickQuizConcepts_RecentUsedOnlyToFill(t *testing.T) {
	for tag, store := range quizSelectionStores(t) {
		t.Run(tag, func(t *testing.T) {
			var cs []concepts.Concept
			for i := 0; i < 7; i++ {
				id := string(rune('a' + i))
				cs = append(cs, quizSelectionConcept(id, fmt.Sprintf("s%d", i)))
			}
			e := quizSelectionEngine(t, store, cs)
			sid := quizSelectionStudent(t, store, "fill")
			for _, c := range cs {
				quizSelectionProgress(t, store, sid, c.ID, string(mastery.StatusPracticing), 1, 0.5, false)
			}
			sess, err := store.CreateSession(sid)
			if err != nil {
				t.Fatalf("create attempt session: %v", err)
			}
			base := time.Now().UTC().Add(-time.Hour)
			for i, id := range []string{"a", "b", "c", "d"} {
				if _, err := store.RecordAttempt(storage.AttemptEntry{
					SessionID: sess.ID, StudentID: sid, ConceptID: id, Answer: "4", Expected: "4",
					Correct: true, Timestamp: base.Add(time.Duration(i) * time.Minute), Source: TaskQuiz,
				}); err != nil {
					t.Fatalf("record quiz attempt: %v", err)
				}
			}

			picked, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick: %v", err)
			}
			// Three unused candidates cannot fill five places, so two recent
			// concepts return. They fill after every non-recent one.
			if got, want := quizSelectionIDs(picked), []string{"e", "f", "g", "a", "b"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("picked %v, want %v", got, want)
			}
		})
	}
}

func TestEngine_PickQuizConcepts_NoCandidatesOrSmallPool(t *testing.T) {
	for tag, store := range quizSelectionStores(t) {
		t.Run(tag, func(t *testing.T) {
			e := quizSelectionEngine(t, store, []concepts.Concept{
				quizSelectionConcept("a", "s1"),
				quizSelectionConcept("b", "s2"),
			})
			sid := quizSelectionStudent(t, store, "empty")

			picked, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick with no history: %v", err)
			}
			if len(picked) != 0 {
				t.Fatalf("no history returned %v, want no candidates", quizSelectionIDs(picked))
			}

			quizSelectionProgress(t, store, sid, "a", string(mastery.StatusUnseen), 0, 0.9, false)
			picked, err = e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick with unseen only: %v", err)
			}
			if len(picked) != 0 {
				t.Fatalf("an UNSEEN row returned %v, want no candidates", quizSelectionIDs(picked))
			}

			quizSelectionProgress(t, store, sid, "a", string(mastery.StatusPracticing), 1, 0.2, false)
			first, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick small pool: %v", err)
			}
			second, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick small pool again: %v", err)
			}
			if got, want := quizSelectionIDs(first), []string{"a"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("small pool picked %v, want %v", got, want)
			}
			if !reflect.DeepEqual(quizSelectionIDs(first), quizSelectionIDs(second)) {
				t.Errorf("selection is not deterministic: %v then %v", quizSelectionIDs(first), quizSelectionIDs(second))
			}
		})
	}
}

func TestEngine_PickQuizConcepts_PrerequisiteGate(t *testing.T) {
	for tag, store := range quizSelectionStores(t) {
		t.Run(tag, func(t *testing.T) {
			e := quizSelectionEngine(t, store, []concepts.Concept{
				quizSelectionConcept("a", "s1"),
				quizSelectionConcept("b", "s2", "a"),
			})
			sid := quizSelectionStudent(t, store, "prereq")
			quizSelectionProgress(t, store, sid, "a", string(mastery.StatusPracticing), 1, 0.2, false)
			quizSelectionProgress(t, store, sid, "b", string(mastery.StatusPracticing), 1, 0.9, false)

			picked, err := e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick before prerequisite: %v", err)
			}
			if got, want := quizSelectionIDs(picked), []string{"a"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("locked dependent picked: %v, want %v", got, want)
			}

			quizSelectionProgress(t, store, sid, "a", string(mastery.StatusMastered), 3, 0, false)
			picked, err = e.PickQuizConcepts(sid)
			if err != nil {
				t.Fatalf("pick after prerequisite: %v", err)
			}
			// Mastery satisfies the prerequisite, but the mastered concept is
			// retention-review material and is not assessed again here.
			if got, want := quizSelectionIDs(picked), []string{"b"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("after mastery picked %v, want %v", got, want)
			}
		})
	}
}
