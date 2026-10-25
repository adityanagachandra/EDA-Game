package grader

import "strings"

type Router struct {
	numeric    *numericGrader
	choice     *choiceGrader
	comparison *comparisonGrader
	ordering   *orderingGrader
	tuple      *tupleGrader
	complex    *complexGrader
}

func NewRouter() *Router {
	return &Router{
		numeric:    &numericGrader{},
		choice:     &choiceGrader{},
		comparison: &comparisonGrader{},
		ordering:   &orderingGrader{},
		tuple:      &tupleGrader{},
		complex:    &complexGrader{},
	}
}

func (r *Router) Grade(gradingType GradingType, expected, answer string) Result {
	// The empty set has several standard written forms. A set-operation
	// question whose expected answer is ∅ must accept {}: that is the form the
	// lesson itself teaches, and a learner should not need LaTeX to answer it.
	if isEmptySet(expected) && isEmptySetAnswer(answer) {
		return Result{Correct: true, Score: 1}
	}
	// A yes/no answer is a word, so it is graded as a choice whatever the
	// concept's type says. Symbolic grading of "yes" asks SymPy to parse an
	// English word, and the registry now serves those answers as options
	// (report 118: "Can't type yes or no").
	if isBooleanWord(expected) {
		return r.choice.grade(expected, answer)
	}
	switch gradingType {
	case GradingNumeric:
		return r.numeric.grade(expected, answer)
	case GradingPolynomial, GradingExpression, GradingSymbolic:
		return gradeSymPy(expected, answer)
	case GradingMultipleChoice:
		return r.choice.grade(expected, answer)
	case GradingComparison:
		return r.comparison.grade(expected, answer)
	case GradingOrdering:
		return r.ordering.grade(expected, answer)
	case GradingTuple:
		return r.tuple.grade(expected, answer)
	case GradingComplex:
		return r.complex.grade(expected, answer)
	default:
		return r.tuple.grade(expected, answer)
	}
}

func (r *Router) Close() {}

// isEmptySet reports whether an expected answer is the empty set in any of the
// forms the corpus writes it.
func isEmptySet(s string) bool {
	switch strings.TrimSpace(s) {
	case "∅", "\\emptyset", "\\varnothing", "{}", "\\{\\}":
		return true
	}
	return false
}

// isEmptySetAnswer accepts the empty set written with braces, with the Unicode
// symbol, with LaTeX, or in words. Braces are how the set-operations lesson
// spells it, so requiring ∅ would fail a learner who read the lesson.
func isEmptySetAnswer(s string) bool {
	n := strings.ToLower(strings.TrimSpace(s))
	n = strings.ReplaceAll(n, " ", "")
	switch n {
	case "∅", "{}", "\\emptyset", "\\varnothing", "\\{\\}", "\\lbrace\\rbrace", "\\left\\{\\right\\}", "emptyset":
		return true
	}
	return false
}

// isBooleanWord reports whether an expected answer is one of the four words a
// yes/no question can have.
func isBooleanWord(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "yes", "no", "true", "false":
		return true
	}
	return false
}
