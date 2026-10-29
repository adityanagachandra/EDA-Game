package infra

import (
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// variant returns the knowledge point a request asks for (1, 2 or 3), or a uniform pick for an
// unguided request.
func variant(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

// numOptions formats a computed multiple-choice answer and its distractors at the given
// precision with a unit suffix (which may be empty), dropping any distractor that formats the
// same as the answer or as an earlier distractor, and padding with nearby values until there are
// at least three distractors.
func numOptions(correct float64, places int, unit string, cands ...float64) (string, []string) {
	f := func(x float64) string {
		s := edakit.Num(x, places)
		if unit != "" {
			s += " " + unit
		}
		return s
	}
	ans := f(correct)
	seen := map[string]bool{ans: true}
	var out []string
	for _, c := range cands {
		s := f(c)
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	step := 1.0
	for i := 0; i < places; i++ {
		step /= 10
	}
	if correct != 0 {
		step = absf(correct) * 0.13
	}
	for k := 1; len(out) < 3 && k < 50; k++ {
		for _, x := range []float64{correct + float64(k)*step, correct - float64(k)*step} {
			s := f(x)
			if !seen[s] && len(out) < 4 {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return ans, out
}

func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// names joins a list as "a, b and c".
func names(xs []string) string {
	switch len(xs) {
	case 0:
		return "none"
	case 1:
		return xs[0]
	}
	return strings.Join(xs[:len(xs)-1], ", ") + " and " + xs[len(xs)-1]
}
