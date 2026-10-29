// Package pkg generates packaging and heterogeneous-integration problems: package structures,
// wafer-level and advanced packaging, chiplets and die-to-die links, thermal, mechanical and
// reliability physics, 3D-IC design and the agent-first frontier of package co-design.
package pkg

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

// fl is floor(a/b) with a guard against binary rounding just below an integer.
func fl(a, b float64) int {
	return int(math.Floor(a/b + 1e-9))
}

// cl is ceil(a/b) with a guard against binary rounding just above an integer.
func cl(a, b float64) int {
	return int(math.Ceil(a/b - 1e-9))
}

// f formats a value at the given number of decimals with trailing zeros trimmed.
func f(x float64, places int) string { return edakit.Num(x, places) }

// num builds a numeric problem with one named fact.
func num(q, ans, expl, key string) generator.Problem {
	return edakit.Numeric(q, ans, expl, map[string]string{key: ans})
}

// bank wraps authored items into a multiple-choice generator.
func bank(items []edakit.Item) *edakit.Bank { return &edakit.Bank{Items: items} }
