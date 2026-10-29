// Package em generates circuits-and-electromagnetics problems: circuit theory, AC analysis,
// signals and systems, fields, transmission lines, EMC, field solvers and the agent-first
// frontier of EM engines.
package em

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Physical constants in SI units.
const (
	eps0     = 8.854187817e-12 // F/m
	mu0      = 4e-7 * math.Pi  // H/m
	cLight   = 2.99792458e8    // m/s
	rhoCu    = 1.72e-8         // ohm m, bulk copper at room temperature
	coulombK = 8.987551787e9   // 1/(4 pi eps0), m/F
)

// Register adds every circuits-and-electromagnetics generator to reg.
func Register(reg *generator.Registry) {
	registerCircuit(reg)
	registerAC(reg)
	registerSignal(reg)
	registerField(reg)
	registerTline(reg)
	registerEMC(reg)
	registerSolver(reg)
	registerFrontier(reg)
}

// kp returns the knowledge point to serve: the requested one, or a uniform draw when the
// request is unguided.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// num is a short alias for edakit.Num.
func num(x float64, places int) string { return edakit.Num(x, places) }

// sp is a short alias for edakit.Sprintf.
func sp(format string, a ...any) string { return edakit.Sprintf(format, a...) }

// par is the parallel combination of resistances (or any reciprocal-additive quantity).
func par(xs ...float64) float64 {
	s := 0.0
	for _, x := range xs {
		s += 1 / x
	}
	return 1 / s
}

// pick draws one float from xs.
func pick(xs ...float64) float64 { return edakit.Pick(xs) }

// rnd returns lo + step*k for a uniform k in [0, n].
func rnd(lo, step float64, n int) float64 { return lo + step*float64(edakit.Int(0, n)) }

// solve solves the dense system a x = b by Gaussian elimination with partial pivoting.
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
			f := m[r][col] / m[col][col]
			for k := col; k <= n; k++ {
				m[r][k] -= f * m[col][k]
			}
		}
	}
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		s := m[i][n]
		for k := i + 1; k < n; k++ {
			s -= m[i][k] * x[k]
		}
		x[i] = s / m[i][i]
	}
	return x
}

// deg converts radians to degrees.
func deg(x float64) float64 { return x * 180 / math.Pi }

// facts builds a facts map from alternating key, value pairs.
func facts(kv ...string) map[string]string {
	m := map[string]string{}
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i]] = kv[i+1]
	}
	return m
}
