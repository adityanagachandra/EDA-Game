package agentic

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// opts returns at least three distinct distractors for a computed multiple-choice answer. The
// candidates are the classic wrong computations; when they collide with the answer or with each
// other, values are stepped away from the answer by step until three remain.
func opts(ans float64, places int, unit string, step float64, cands ...float64) []string {
	a := edakit.Num(ans, places)
	seen := map[string]bool{a: true}
	var out []string
	for _, c := range cands {
		if math.IsNaN(c) || math.IsInf(c, 0) {
			continue
		}
		s := edakit.Num(c, places)
		if !seen[s] {
			seen[s] = true
			out = append(out, s+unit)
		}
	}
	for k := 1; len(out) < 3; k++ {
		for _, c := range []float64{ans + float64(k)*step, ans - float64(k)*step} {
			s := edakit.Num(c, places)
			if !seen[s] && len(out) < 3 && c >= 0 {
				seen[s] = true
				out = append(out, s+unit)
			}
		}
		if k > 50 {
			break
		}
	}
	if len(out) > 4 {
		out = out[:4]
	}
	return out
}

// choice renders a computed multiple-choice problem whose options are numbers with a unit.
func choice(q string, ans float64, places int, unit, exp string, step float64, cands ...float64) generator.Problem {
	return edakit.Choice(q, edakit.Num(ans, places)+unit, exp, opts(ans, places, unit, step, cands...)...)
}

// phi is the standard normal distribution function.
func phi(z float64) float64 { return 0.5 * (1 + math.Erf(z/math.Sqrt2)) }

// binom returns the binomial coefficient C(n, k) as a float.
func binom(n, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	r := 1.0
	for i := 1; i <= k; i++ {
		r = r * float64(n-k+i) / float64(i)
	}
	return r
}

// kp returns the requested knowledge point, or a uniform one when the request is unguided.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}
