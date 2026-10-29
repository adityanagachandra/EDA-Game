package csim

import (
	"math"
	"strings"

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
	return edakit.Numeric(q, ans, expl, facts)
}

// f formats a value for a stem or an explanation with up to places decimals, trailing zeros
// trimmed.
func f(x float64, places int) string { return edakit.Num(x, places) }

// pickF returns a uniform element of a float list.
func pickF(xs ...float64) float64 { return edakit.Pick(xs) }

// stepF draws a value lo + k*step for k in [0, n].
func stepF(lo, step float64, n int) float64 { return lo + step*float64(edakit.Int(0, n)) }

// solve solves a small dense system by Gaussian elimination with partial pivoting.
func solve(a [][]float64, b []float64) []float64 {
	n := len(b)
	m := make([][]float64, n)
	for i := range a {
		m[i] = append(append([]float64(nil), a[i]...), b[i])
	}
	for c := 0; c < n; c++ {
		p := c
		for r := c + 1; r < n; r++ {
			if math.Abs(m[r][c]) > math.Abs(m[p][c]) {
				p = r
			}
		}
		m[c], m[p] = m[p], m[c]
		for r := c + 1; r < n; r++ {
			k := m[r][c] / m[c][c]
			for j := c; j <= n; j++ {
				m[r][j] -= k * m[c][j]
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

// zeros returns an n by n zero matrix.
func zeros(n int) [][]float64 {
	m := make([][]float64, n)
	for i := range m {
		m[i] = make([]float64, n)
	}
	return m
}

// resistor is a two-terminal conductance between nodes a and b (0 is ground), with its
// resistance in kilohms.
type resistor struct {
	a, b int
	r    float64
}

// kohmChoices are resistances whose conductances in mS have short decimals.
var kohmChoices = []float64{0.5, 1, 2, 4, 5, 0.25, 2.5, 10}

// randomNetwork draws a connected resistor network on nodes 1..n with at least one path to
// ground, plus extra branches.
func randomNetwork(n, extra int) []resistor {
	var rs []resistor
	for i := 2; i <= n; i++ {
		rs = append(rs, resistor{edakit.Int(1, i-1), i, edakit.Pick(kohmChoices)})
	}
	rs = append(rs, resistor{edakit.Int(1, n), 0, edakit.Pick(kohmChoices)})
	for k := 0; k < extra; k++ {
		a := edakit.Int(1, n)
		b := edakit.Int(0, n)
		if a == b {
			b = 0
		}
		rs = append(rs, resistor{a, b, edakit.Pick(kohmChoices)})
	}
	return rs
}

// stampG builds the n by n nodal conductance matrix in mS.
func stampG(n int, rs []resistor) [][]float64 {
	g := zeros(n)
	for _, r := range rs {
		c := 1 / r.r
		if r.a > 0 {
			g[r.a-1][r.a-1] += c
		}
		if r.b > 0 {
			g[r.b-1][r.b-1] += c
		}
		if r.a > 0 && r.b > 0 {
			g[r.a-1][r.b-1] -= c
			g[r.b-1][r.a-1] -= c
		}
	}
	return g
}

// describeResistors renders a resistor list in prose.
func describeResistors(rs []resistor) string {
	var parts []string
	for i, r := range rs {
		if r.b == 0 {
			parts = append(parts, edakit.Sprintf("R%d from node %d to ground, %s k$\\Omega$", i+1, r.a, f(r.r, 2)))
		} else {
			parts = append(parts, edakit.Sprintf("R%d between nodes %d and %d, %s k$\\Omega$", i+1, r.a, r.b, f(r.r, 2)))
		}
	}
	return strings.Join(parts, "; ")
}

// nodeName renders a node index, 0 being ground.
func nodeName(i int) string {
	if i == 0 {
		return "0"
	}
	return edakit.IntStr(i)
}

// ordinal places for an answer of rough magnitude x so that about four significant digits
// remain, between 1 and 6 places.
func placesFor(x float64) int {
	a := math.Abs(x)
	switch {
	case a >= 100:
		return 1
	case a >= 10:
		return 2
	case a >= 1:
		return 3
	case a >= 0.1:
		return 4
	case a >= 0.01:
		return 5
	default:
		return 6
	}
}

// sci renders a value in scientific notation for prose, as inline LaTeX, keeping plain
// decimals for moderate magnitudes.
func sci(x float64) string {
	s := sciM(x)
	if strings.ContainsAny(s, "^\\") {
		return "$" + s + "$"
	}
	return s
}

// sciM is sci for use inside an existing math span: the same notation without delimiters.
func sciM(x float64) string {
	if x == 0 {
		return "0"
	}
	a := math.Abs(x)
	if a >= 0.01 && a < 1e5 {
		return f(x, 4)
	}
	e := int(math.Floor(math.Log10(a)))
	m := x / math.Pow(10, float64(e))
	if math.Abs(m-math.Round(m)) < 1e-9 && math.Round(math.Abs(m)) == 1 {
		if m < 0 {
			return edakit.Sprintf("-10^{%d}", e)
		}
		return edakit.Sprintf("10^{%d}", e)
	}
	return edakit.Sprintf("%s \\times 10^{%d}", f(m, 4), e)
}
