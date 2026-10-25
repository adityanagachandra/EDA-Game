package xp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The landing page used to advertise a 150 XP mastery check while the gate was 50. Nothing
// broke: the page was static, the number was a string literal, and ADR-020's rescale touched
// Go and the daily-goal backfill but not the marketing copy. It is the one number in the
// generated graph artifact that build-graph.mjs cannot derive from the corpus, so it is the
// one that can drift, and it did.
//
// No learner surface prints the gate now; the interface says when a check is ready. The
// artifact keeps carrying it as a policy mirror so the browser build and the Go constant
// cannot disagree silently. Editing the JSON by hand fails here; changing the engine constant
// without rebuilding the artifact fails here too.
func TestLandingQuizGateMatchesEngine(t *testing.T) {
	path := filepath.Join("..", "..", "web", "next-app", "data", "graph.meta.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v (run `npm run graph:build`)", path, err)
	}
	var meta struct {
		QuizGateXP *int `json:"quizGateXP"`
	}
	if err := json.Unmarshal(raw, &meta); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	if meta.QuizGateXP == nil {
		t.Fatalf("%s has no quizGateXP; the generated artifact lost its policy mirror", path)
	}
	if *meta.QuizGateXP != QuizGateXP {
		t.Errorf("graph artifact carries a %d XP mastery-check gate, the engine gates it at %d XP\n"+
			"fix: set QUIZ_GATE_XP in web/next-app/lib/graphPayload.ts to %d and run `npm run graph:build`",
			*meta.QuizGateXP, QuizGateXP, QuizGateXP)
	}
}
