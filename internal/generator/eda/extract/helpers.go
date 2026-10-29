package extract

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// kp returns the knowledge point to drill: the requested one, or a random one when the
// request is unguided.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// num builds a numeric problem whose answer is x rounded to places, appending the rounding
// instruction to the stem.
func num(stem string, x float64, places int, expl string, facts map[string]string) generator.Problem {
	q := stem
	switch {
	case places == 1:
		q += " Round to 1 decimal place."
	case places > 1:
		q += edakit.Sprintf(" Round to %d decimal places.", places)
	}
	ans := edakit.Num(x, places)
	if facts == nil {
		facts = map[string]string{}
	}
	facts["answer"] = ans
	return edakit.Numeric(q, ans, expl, facts)
}

// f formats a value for prose at the given precision.
func f(x float64, places int) string { return edakit.Num(x, places) }

// numChoice builds a computed multiple-choice problem whose options are numbers with a unit.
// Candidates that coincide with the answer (at the shown precision) are dropped, and nearby
// values fill the list up to three distractors.
func numChoice(q string, ans float64, places int, unit, expl string, cands ...float64) generator.Problem {
	a := edakit.Num(ans, places) + unit
	seen := map[string]bool{a: true}
	var d []string
	for _, c := range cands {
		s := edakit.Num(c, places) + unit
		if !seen[s] && !math.IsNaN(c) && !math.IsInf(c, 0) {
			seen[s] = true
			d = append(d, s)
		}
	}
	step := math.Max(math.Abs(ans)*0.15, math.Pow(10, -float64(places)))
	for k := 1; len(d) < 3; k++ {
		for _, c := range []float64{ans + float64(k)*step, ans - float64(k)*step} {
			s := edakit.Num(c, places) + unit
			if !seen[s] && len(d) < 3 {
				seen[s] = true
				d = append(d, s)
			}
		}
	}
	return edakit.Choice(q, a, expl, d...)
}

// eps0 is the vacuum permittivity in fF/um.
const eps0 = 8.854e-3

// plural writes a count with its noun, adding an s unless the count is one.
func plural(n int, noun string) string {
	if n == 1 {
		return "one " + noun
	}
	return edakit.Sprintf("%d %ss", n, noun)
}
