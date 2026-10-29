// Package cts generates clock-network and power-delivery problems: clock tree objectives and
// topologies, zero-skew and bounded-skew merging (Tsay, DME), buffered trees, useful skew and
// clock scheduling, clock gating, power distribution networks and IR drop, decoupling and
// inductive noise, electromigration, power gating and multi-voltage intent, power analysis,
// and the agent-first frontier of clock and power signoff.
package cts

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// kpOf returns the requested knowledge point, or a uniform one when the request is unguided.
func kpOf(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// pick returns a uniform value lo, lo+inc, ..., hi.
func pick(lo, hi, inc float64) float64 {
	n := int(math.Round((hi - lo) / inc))
	return lo + inc*float64(edakit.Int(0, n))
}

// f formats a value at the given number of places for prose (trailing zeros trimmed).
func f(x float64, places int) string { return edakit.Num(x, places) }

// numOptions formats a computed value and its wrong computations as options with a unit
// suffix, dropping any wrong value that formats the same as the correct one or another option,
// and padding with nearby values so that there are always at least three distractors.
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
	for k := 1; len(out) < 3; k++ {
		for _, m := range []float64{1 + 0.2*float64(k), 1 - 0.15*float64(k)} {
			s := edakit.Num(correct*m+float64(k), places) + unit
			if !seen[s] && len(out) < 3 {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return c, out
}

// ints joins integers with commas for prose.
func ints(xs []int) string {
	var p []string
	for _, x := range xs {
		p = append(p, edakit.IntStr(x))
	}
	return strings.Join(p, ", ")
}

// floats joins values with commas for prose.
func floats(xs []float64, places int) string {
	var p []string
	for _, x := range xs {
		p = append(p, f(x, places))
	}
	return strings.Join(p, ", ")
}

func maxInt(xs []int) int {
	m := xs[0]
	for _, x := range xs {
		if x > m {
			m = x
		}
	}
	return m
}

func minInt(xs []int) int {
	m := xs[0]
	for _, x := range xs {
		if x < m {
			m = x
		}
	}
	return m
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func absF(x float64) float64 { return math.Abs(x) }

// solve solves the dense system a x = b by Gaussian elimination with partial pivoting. The
// systems here are the few-node conductance matrices of the IR-drop drills, so dense is fine.
func solve(a [][]float64, b []float64) []float64 {
	n := len(b)
	m := make([][]float64, n)
	for i := range a {
		m[i] = append(append([]float64(nil), a[i]...), b[i])
	}
	for col := 0; col < n; col++ {
		p := col
		for r := col + 1; r < n; r++ {
			if math.Abs(m[r][col]) > math.Abs(m[p][col]) {
				p = r
			}
		}
		m[col], m[p] = m[p], m[col]
		for r := col + 1; r < n; r++ {
			k := m[r][col] / m[col][col]
			for c := col; c <= n; c++ {
				m[r][c] -= k * m[col][c]
			}
		}
	}
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		s := m[i][n]
		for j := i + 1; j < n; j++ {
			s -= m[i][j] * x[j]
		}
		x[i] = s / m[i][i]
	}
	return x
}

// mix3 serves knowledge point k from parts[k-1].
func mix3(a, b, c generator.Generator) *edakit.Mix {
	return &edakit.Mix{Parts: []generator.Generator{a, b, c}}
}

func bank(items []edakit.Item) *edakit.Bank { return &edakit.Bank{Items: items} }
