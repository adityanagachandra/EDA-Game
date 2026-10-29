// Package edakit is the shared toolkit for the electronic-design-automation generators.
//
// The EDA corpus mixes two kinds of concept. A computational concept (Elmore delay, setup
// slack, half-perimeter wirelength, a Fiduccia-Mattheyses gain) draws its operands and
// computes its answer, exactly like the mathematics generators. A conceptual concept (what a
// Liberty file abstracts away, why a placer is analytical before it is legal) has no operands,
// so its supply is an authored bank of distinct situations. Both are served through the same
// generator.Generator interface, and both are drawn from the global math/rand that the
// registry seeds per question, so a seeded request reproduces.
//
// A bank is not a table of three questions. The variety floor measures distinct situations,
// and a bank with six items supplies six. Every bank in this corpus carries at least twelve
// items, and the domain contract test (CheckDomain) fails a concept below MinVariety.
package edakit

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
)

// Item is one authored multiple-choice situation.
//
// Q is the stem, A the correct option and D at least three distractors. Why is the
// explanation a learner sees after answering; it states the reason, not the letter. KP pins
// the item to a knowledge point (1, 2 or 3) so that guided practice on one step stays on that
// step; zero means the item serves any step.
type Item struct {
	Q   string
	A   string
	D   []string
	Why string
	KP  int
}

// Bank is a multiple-choice generator over authored items.
type Bank struct {
	Items []Item
}

// Generate picks an item, preferring the requested knowledge point, and renders it.
func (b *Bank) Generate(ctx generator.GeneratorContext) generator.Problem {
	if len(b.Items) == 0 {
		return generator.Problem{Question: "empty bank", Answer: "empty"}
	}
	pool := b.Items
	if ctx.KPOrdinal > 0 {
		var matched []Item
		for _, it := range b.Items {
			if it.KP == ctx.KPOrdinal || it.KP == 0 {
				matched = append(matched, it)
			}
		}
		if len(matched) > 0 {
			pool = matched
		}
	}
	idx := rand.Intn(len(pool))
	it := pool[idx]
	p := generator.MultipleChoice(it.Q, it.A, it.D...)
	p.Explanation = it.Why
	if p.Explanation == "" {
		p.Explanation = "The correct option is: " + it.A + "."
	}
	return tidyProblem(p)
}

// Func adapts a plain function into a generator, for computational concepts.
type Func func(ctx generator.GeneratorContext) generator.Problem

// Generate calls the function.
func (f Func) Generate(ctx generator.GeneratorContext) generator.Problem { return tidyProblem(f(ctx)) }

// Mix serves several generators for one concept, choosing by knowledge point when the request
// names one (ordinal k serves Parts[(k-1) % len]) and uniformly otherwise. It is how a concept
// combines a computational drill with an authored bank of conceptual situations.
type Mix struct {
	Parts []generator.Generator
}

// Generate delegates to one part.
func (m *Mix) Generate(ctx generator.GeneratorContext) generator.Problem {
	if len(m.Parts) == 0 {
		return generator.Problem{Question: "empty mix", Answer: "empty"}
	}
	if ctx.KPOrdinal > 0 {
		return m.Parts[(ctx.KPOrdinal-1)%len(m.Parts)].Generate(ctx)
	}
	return m.Parts[rand.Intn(len(m.Parts))].Generate(ctx)
}

// Int returns a uniform integer in [lo, hi].
func Int(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	return lo + rand.Intn(hi-lo+1)
}

// Scaled returns an integer in [lo, lo + span*difficulty], at least lo and at most hi.
// It lets a harder request draw a larger instance (more cells, more nets, more stages).
func Scaled(ctx generator.GeneratorContext, lo, hi int) int {
	d := ctx.Difficulty
	if d <= 0 {
		d = 0.5
	}
	top := lo + int(math.Round(float64(hi-lo)*d))
	if top < lo {
		top = lo
	}
	return Int(lo, top)
}

// Pick returns a uniform element of xs.
func Pick[T any](xs []T) T { return xs[rand.Intn(len(xs))] }

// Shuffle returns a shuffled copy of xs.
func Shuffle[T any](xs []T) []T {
	out := append([]T(nil), xs...)
	rand.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}

// Num formats a value for a numeric answer: rounded to places, trailing zeros trimmed, so the
// grader reads it at the precision the stem asked for ("round to 2 decimal places").
func Num(x float64, places int) string {
	if places <= 0 {
		return strconv.FormatInt(int64(math.Round(x)), 10)
	}
	s := strconv.FormatFloat(x, 'f', places, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimSuffix(s, ".")
	}
	if s == "-0" {
		s = "0"
	}
	return s
}

// Round rounds x to places decimal places.
func Round(x float64, places int) float64 {
	p := math.Pow(10, float64(places))
	return math.Round(x*p) / p
}

// IntStr formats an integer answer.
func IntStr(n int) string { return strconv.Itoa(n) }

// loneNumber is a math span holding nothing but a number, such as "$1200$". The LaTeX
// canonicaliser renders that as literal dollars, so such spans are written as plain text.
var loneNumber = regexp.MustCompile(`\$(-?\d+(?:\.\d+)?)\$`)

// Tidy unwraps lone-number math spans. Generators often wrap a formatted value in "$...$"
// because the value is sometimes scientific notation (which needs math mode) and sometimes a
// plain number (which must not have it); Tidy makes the plain case correct without every call
// site having to know which case it drew.
func Tidy(s string) string { return loneNumber.ReplaceAllString(s, "$1") }

func tidyProblem(p generator.Problem) generator.Problem {
	p.Question = Tidy(p.Question)
	p.Answer = Tidy(p.Answer)
	p.Explanation = Tidy(p.Explanation)
	return p
}

// Numeric builds a numeric problem with its facts. The answer must match the numeric keypad
// (digits, sign, decimal point, comma-separated lists), which TestAnswerFormatMatchesTheAnswer
// enforces corpus-wide.
func Numeric(question, answer, explanation string, facts map[string]string) generator.Problem {
	return generator.Problem{Question: question, Answer: answer, Explanation: explanation, Facts: facts}
}

// Choice builds a computed multiple-choice problem: a stem whose correct option is computed
// from drawn operands and whose distractors are the classic wrong computations.
func Choice(question, correct, explanation string, distractors ...string) generator.Problem {
	p := generator.MultipleChoice(question, correct, distractors...)
	p.Explanation = explanation
	return p
}

// Distinct returns up to n distinct wrong values near the correct integer answer, for a
// computed multiple-choice question whose options are numbers.
func Distinct(correct int, n int, candidates ...int) []string {
	seen := map[int]bool{correct: true}
	var out []string
	for _, c := range candidates {
		if !seen[c] {
			seen[c] = true
			out = append(out, strconv.Itoa(c))
		}
		if len(out) == n {
			return out
		}
	}
	for d := 1; len(out) < n; d++ {
		for _, c := range []int{correct + d, correct - d} {
			if !seen[c] && len(out) < n {
				seen[c] = true
				out = append(out, strconv.Itoa(c))
			}
		}
	}
	return out
}

// Sprintf is fmt.Sprintf, re-exported so a generator file needs one import for formatting.
func Sprintf(format string, a ...any) string { return fmt.Sprintf(format, a...) }
