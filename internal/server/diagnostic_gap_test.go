package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/engine"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/planning"
	"github.com/chuma-beep/mathua/internal/storage"
)

func diagnosticChainServer(t *testing.T, n int) (*http.ServeMux, storage.Repository) {
	t.Helper()
	raw := make([]concepts.Concept, n)
	for i := range raw {
		id := fmt.Sprintf("c%02d", i)
		raw[i] = concepts.Concept{
			ID:               id,
			Label:            id,
			Domain:           "diagnostic_test",
			Subdomain:        "diagnostic_test.chain",
			GradingType:      "numeric",
			MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
		}
		if i > 0 {
			raw[i].Prerequisites = []string{fmt.Sprintf("c%02d", i-1)}
		}
	}
	d, err := concepts.Build(raw)
	if err != nil {
		t.Fatalf("build chain: %v", err)
	}
	store, err := storage.NewSQLiteStore(":memory:")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	reg := generator.NewRegistry()
	for i := range raw {
		if err := reg.Register(fmt.Sprintf("c%02d", i), &testGen{}); err != nil {
			t.Fatalf("register: %v", err)
		}
	}
	s := New(engine.New(store, d, reg, nil, planning.New(d)), store, nil)
	mux := http.NewServeMux()
	s.Register(mux)
	return mux, store
}

func diagnosticPost(t *testing.T, mux *http.ServeMux, path, body string) (int, map[string]interface{}) {
	t.Helper()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", path, bytes.NewReader([]byte(body))))
	var res map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &res)
	return rec.Code, res
}

// Through the real start and answer handlers: a miss on the selected advanced
// concept descends to its immediate prerequisite, not to the easiest concept
// in the closure.
func TestGoalDiagnosticMissDescendsOnePrerequisite(t *testing.T) {
	mux, _ := diagnosticChainServer(t, 6)
	code, start := diagnosticPost(t, mux, "/api/goal/diagnostic", `{"name":"tester","concept_ids":["c05"]}`)
	if code != 200 {
		t.Fatalf("start: %d %v", code, start)
	}
	sid, _ := start["session_id"].(string)
	cid, _ := start["concept_id"].(string)
	if sid == "" || cid != "c05" {
		t.Fatalf("start = session %q concept %q, want the focused c05", sid, cid)
	}

	// First miss keeps the two-probe rule.
	code, first := diagnosticPost(t, mux, "/api/goal/diagnostic/answer",
		`{"session_id":"`+sid+`","concept_id":"c05","answer":"0","elapsed":5.0}`)
	if code != 200 {
		t.Fatalf("first answer: %d %v", code, first)
	}
	if next, _ := first["concept_id"].(string); next != "c05" {
		t.Fatalf("first miss staged %q, want the second c05 probe", next)
	}

	// Second miss settles c05; the next concept is its immediate prerequisite.
	code, second := diagnosticPost(t, mux, "/api/goal/diagnostic/answer",
		`{"session_id":"`+sid+`","concept_id":"c05","answer":"0","elapsed":5.0}`)
	if code != 200 {
		t.Fatalf("second answer: %d %v", code, second)
	}
	if next, _ := second["concept_id"].(string); next != "c04" {
		t.Fatalf("after a settled miss on c05 the next concept = %q, want c04", next)
	}
}

// A diagnostic answer is placement evidence, not ordinary practice. It must
// not award XP or write a mastery row before the completed report is applied.
func TestGoalDiagnosticAnswerDoesNotAwardXPOrMastery(t *testing.T) {
	mux, store := diagnosticChainServer(t, 6)
	code, start := diagnosticPost(t, mux, "/api/goal/diagnostic", `{"name":"tester","concept_ids":["c05"]}`)
	if code != 200 {
		t.Fatalf("start: %d %v", code, start)
	}
	sid, _ := start["session_id"].(string)
	studentID, _ := start["student_id"].(string)
	if sid == "" || studentID == "" {
		t.Fatalf("start missing session/student: %v", start)
	}
	code, answer := diagnosticPost(t, mux, "/api/goal/diagnostic/answer",
		`{"session_id":"`+sid+`","concept_id":"c05","answer":"4","elapsed":5.0}`)
	if code != 200 {
		t.Fatalf("answer: %d %v", code, answer)
	}

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/api/scores/"+studentID, nil))
	if rec.Code != 200 {
		t.Fatalf("scores: %d %s", rec.Code, rec.Body.String())
	}
	var scores map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &scores); err != nil {
		t.Fatalf("scores json: %v", err)
	}
	if xp, _ := scores["xp_total"].(float64); xp != 0 {
		t.Fatalf("diagnostic answer awarded XP: %v", scores["xp_total"])
	}
	progress, err := store.GetProgress(studentID, "c05")
	if err != nil {
		t.Fatalf("get progress: %v", err)
	}
	if progress != nil {
		t.Fatalf("diagnostic answer wrote progress %+v before the report was applied", progress)
	}
}
