package diagnostic

import "testing"

// A miss on an advanced concept must descend one prerequisite edge, not jump
// to the easiest concept in the whole closure. Report 117 was the inverse
// complaint at session start; this is the answer-dependent form of the same
// prerequisite-awareness requirement.
func TestEmptyScopeIsDoneWithoutSubstitution(t *testing.T) {
	d := chainDAG(t, 3)
	eng := NewEngine(d, chainRegistry(t, 3))
	s := eng.StartWithPathFocus(nil, nil)
	if s.State != StateDone {
		t.Fatalf("empty scope state = %q, want done", s.State)
	}
	p, cid, err := eng.NextQuestion(s)
	if err != nil {
		t.Fatalf("next question: %v", err)
	}
	if p != nil || cid != "" {
		t.Fatalf("empty scope produced question %v for %q, want no substitution", p, cid)
	}
}

func TestMissDescendsToTheNearestPrerequisite(t *testing.T) {
	d := chainDAG(t, 6)
	eng := NewEngine(d, chainRegistry(t, 6))
	s := eng.StartWithPathFocus(pathOf(t, d), []string{"c05"})

	_, cid, err := eng.NextQuestion(s)
	if err != nil {
		t.Fatalf("first question: %v", err)
	}
	if cid != "c05" {
		t.Fatalf("first question = %q, want the focused concept c05", cid)
	}
	eng.RecordAnswer(s, cid, false, false)

	// The concept keeps its second probe before the model moves on.
	_, cid, err = eng.NextQuestion(s)
	if err != nil {
		t.Fatalf("second probe: %v", err)
	}
	if cid != "c05" {
		t.Fatalf("second probe = %q, want c05 until its two probes settle it", cid)
	}
	eng.SettleCurrent(s)

	_, cid, err = eng.NextQuestion(s)
	if err != nil {
		t.Fatalf("after the miss: %v", err)
	}
	if cid != "c04" {
		t.Fatalf("after a miss on c05 the next question = %q, want the nearest prerequisite c04 (not the easiest concept)", cid)
	}
}
