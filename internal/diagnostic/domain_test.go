package diagnostic

import (
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	allgen "github.com/chuma-beep/mathua/internal/generator/all"
	"github.com/chuma-beep/mathua/internal/planning"
)

// The reported Abstract Algebra run opened on "What is the opposite of -2?".
// That concept is a real transitive prerequisite, but it is not the selected
// domain and it must not be the first probe when the whole domain was chosen.
// This loads the shipped curriculum because the defect was about the real
// closure, not a hand-built test graph.
func TestSelectedDomainIsProbedBeforeItsFoundations(t *testing.T) {
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load curriculum: %v", err)
	}
	var focus []string
	for _, c := range dag.Order() {
		if c.Domain == "abstract_algebra" {
			focus = append(focus, c.ID)
		}
	}
	if len(focus) == 0 {
		t.Fatal("curriculum has no abstract_algebra concepts")
	}
	path, err := planning.New(dag).PrerequisitesOf(focus)
	if err != nil {
		t.Fatalf("prerequisite closure: %v", err)
	}
	reg := generator.NewRegistry()
	allgen.Register(reg)
	eng := NewEngine(dag, reg)
	s := eng.StartWithPathFocus(path.Concepts, focus)

	inPath := make(map[string]bool, len(path.Concepts))
	for _, c := range path.Concepts {
		inPath[c.ID] = true
	}
	for i := 0; i < 3; i++ {
		_, cid, err := eng.NextQuestion(s)
		if err != nil {
			t.Fatalf("question %d: %v", i+1, err)
		}
		if cid == "" {
			t.Fatalf("question %d ended the session", i+1)
		}
		c := dag.Concept(cid)
		if c == nil {
			t.Fatalf("question %d named unknown concept %q", i+1, cid)
		}
		if c.Domain != "abstract_algebra" {
			t.Fatalf("question %d = %q (%s), want the selected domain before foundations", i+1, cid, c.Domain)
		}
		if !inPath[cid] {
			t.Fatalf("question %d = %q is outside the selected prerequisite closure", i+1, cid)
		}
		eng.SettleCurrent(s)
	}
}

// The arithmetic question is not unrelated: it is in the closure of a real
// Abstract Algebra target. The distinction matters, because the repair must
// keep legitimate cross-domain prerequisites rather than filtering them out.
func TestAbstractAlgebraClosureContainsTheArithmeticPrerequisite(t *testing.T) {
	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load curriculum: %v", err)
	}
	path, err := planning.New(dag).PrerequisitesOf([]string{"abstract.group.def"})
	if err != nil {
		t.Fatalf("prerequisite closure: %v", err)
	}
	found := false
	for _, c := range path.Concepts {
		if c.ID == "arith.neg.number_line" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("arith.neg.number_line is no longer a prerequisite of abstract.group.def; the reported question would then be unrelated")
	}
}
