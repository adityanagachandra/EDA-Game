// Package fab generates fabrication and design-technology co-optimization problems: unit
// processes, lithography and patterning, interconnect, process design kits, yield and DTCO.
package fab

import (
	"math"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// kpOf returns the requested knowledge point, or a uniform one when the request is unguided.
func kpOf(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return 1 + rand.Intn(3)
}

// step draws a value lo, lo+inc, ..., up to hi.
func step(lo, hi, inc float64) float64 {
	n := int(math.Round((hi - lo) / inc))
	return lo + inc*float64(edakit.Int(0, n))
}

// numOptions formats a computed value and its wrong computations as options with a unit
// suffix, dropping any wrong value that rounds to the correct one or to another option.
func numOptions(correct float64, places int, unit string, wrong ...float64) (string, []string) {
	c := edakit.Num(correct, places) + unit
	seen := map[string]bool{c: true}
	var out []string
	for _, w := range wrong {
		s := edakit.Num(w, places) + unit
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	// Pad with nearby values if the classic mistakes collided.
	for k := 1; len(out) < 3; k++ {
		for _, f := range []float64{1 + 0.15*float64(k), 1 - 0.12*float64(k)} {
			s := edakit.Num(correct*f+float64(k), places) + unit
			if !seen[s] && len(out) < 3 {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return c, out
}

// bank3 builds the bank-compute-bank mix used by the mixed concepts: knowledge points 1 and 3
// come from an authored bank, knowledge point 2 from a computed multiple-choice drill.
func bank3(items []edakit.Item, compute edakit.Func) *edakit.Mix {
	b := &edakit.Bank{Items: items}
	return &edakit.Mix{Parts: []generator.Generator{b, compute, b}}
}
