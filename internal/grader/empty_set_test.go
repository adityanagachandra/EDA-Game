package grader

import "testing"

// Report 271: the set-operations lesson teaches that the empty set is written
// ∅ or {}, but the graded answer was the Unicode symbol only. Braces are the
// form a beginner is most likely to type, and the on-screen keyboard had no ∅.
func TestEmptySetAcceptsEveryStandardForm(t *testing.T) {
	r := NewRouter()
	accepted := []string{"∅", "{}", "{ }", "\\emptyset", "\\varnothing", "\\{\\}", "\\lbrace\\rbrace", "\\left\\{\\right\\}", "empty set", "Empty Set"}
	for _, answer := range accepted {
		res := r.Grade(GradingMultipleChoice, "∅", answer)
		if !res.Correct {
			t.Errorf("expected %q to be accepted for ∅", answer)
		}
	}
	for _, answer := range []string{"", "{1}", "0", "empty", "\\{1\\}"} {
		res := r.Grade(GradingMultipleChoice, "∅", answer)
		if res.Correct {
			t.Errorf("expected %q to be rejected for ∅", answer)
		}
	}
}

func TestEmptySetExpectedFormAlsoWorks(t *testing.T) {
	r := NewRouter()
	if res := r.Grade(GradingMultipleChoice, "{}", "∅"); !res.Correct {
		t.Fatal("expected ∅ to be accepted when the key is {}")
	}
}
