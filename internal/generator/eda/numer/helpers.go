// Package numer generates numerical-methods-for-EDA problems: sparse storage and orderings,
// direct factorization and conditioning, Newton iteration, time integration, Krylov and
// multigrid solvers, PDE discretization for thermal and field problems, model order
// reduction, automatic differentiation and parallel kernels, and the agent-first frontier of
// numerical engines.
package numer

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// kpOf returns the knowledge point to serve: the requested one, or a random one when the
// request is unguided.
func kpOf(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// f formats a value for a stem or an answer, trimmed to places decimals.
func f(x float64, places int) string { return edakit.Num(x, places) }

// itoa formats an integer.
func itoa(n int) string { return edakit.IntStr(n) }

// sp is edakit.Sprintf.
var sp = edakit.Sprintf

// numeric builds a numeric problem.
func numeric(q, ans, why string, facts map[string]string) generator.Problem {
	return edakit.Numeric(q, ans, why, facts)
}

// pmat renders a matrix of numbers as an inline KaTeX pmatrix.
func pmat(m [][]float64, places int) string {
	var rows []string
	for _, r := range m {
		var cells []string
		for _, v := range r {
			cells = append(cells, f(v, places))
		}
		rows = append(rows, strings.Join(cells, " & "))
	}
	return `$\begin{pmatrix} ` + strings.Join(rows, ` \\ `) + ` \end{pmatrix}$`
}

// pvec renders a column vector as an inline KaTeX pmatrix.
func pvec(v []float64, places int) string {
	m := make([][]float64, len(v))
	for i, x := range v {
		m[i] = []float64{x}
	}
	return pmat(m, places)
}

// ints renders an integer list as "a, b, c".
func ints(xs []int) string {
	var s []string
	for _, x := range xs {
		s = append(s, itoa(x))
	}
	return strings.Join(s, ", ")
}

// sci renders x as a KaTeX scientific-notation string with one decimal of mantissa, and
// returns the value the reader sees, so computations use exactly what is shown.
func sci(mant float64, exp int) (string, float64) {
	return sp(`$%s \times 10^{%d}$`, f(mant, 2), exp), mant * math.Pow(10, float64(exp))
}

// solve2 solves a 2x2 linear system by Cramer's rule.
func solve2(a [][]float64, b []float64) []float64 {
	det := a[0][0]*a[1][1] - a[0][1]*a[1][0]
	return []float64{(b[0]*a[1][1] - a[0][1]*b[1]) / det, (a[0][0]*b[1] - b[0]*a[1][0]) / det}
}

// matvec multiplies a dense matrix by a vector.
func matvec(a [][]float64, x []float64) []float64 {
	y := make([]float64, len(a))
	for i := range a {
		for j := range x {
			y[i] += a[i][j] * x[j]
		}
	}
	return y
}

func dot(a, b []float64) float64 {
	s := 0.0
	for i := range a {
		s += a[i] * b[i]
	}
	return s
}

func norm2(a []float64) float64 { return math.Sqrt(dot(a, a)) }

// graph is an undirected simple graph on vertices 1..n.
type graph struct {
	n   int
	adj map[int]map[int]bool
}

func newGraph(n int) *graph {
	g := &graph{n: n, adj: map[int]map[int]bool{}}
	for v := 1; v <= n; v++ {
		g.adj[v] = map[int]bool{}
	}
	return g
}

func (g *graph) add(u, v int) {
	if u == v {
		return
	}
	g.adj[u][v] = true
	g.adj[v][u] = true
}

func (g *graph) clone() *graph {
	h := newGraph(g.n)
	for u, nb := range g.adj {
		for v := range nb {
			h.adj[u][v] = true
		}
	}
	return h
}

// edges lists the edges as "u-v" with u < v, sorted.
func (g *graph) edges() []string {
	var out []string
	for u := 1; u <= g.n; u++ {
		var vs []int
		for v := range g.adj[u] {
			if v > u {
				vs = append(vs, v)
			}
		}
		sort.Ints(vs)
		for _, v := range vs {
			out = append(out, sp("%d-%d", u, v))
		}
	}
	return out
}

func (g *graph) edgeCount() int { return len(g.edges()) }

// randomGraph draws a connected graph on n vertices: a random spanning tree plus extra edges.
func randomGraph(n int, extra float64) *graph {
	g := newGraph(n)
	perm := edakit.Shuffle(seq(1, n))
	for i := 1; i < n; i++ {
		g.add(perm[i], perm[edakit.Int(0, i-1)])
	}
	for u := 1; u <= n; u++ {
		for v := u + 1; v <= n; v++ {
			if !g.adj[u][v] && randFloat() < extra {
				g.add(u, v)
			}
		}
	}
	return g
}

func seq(lo, hi int) []int {
	var out []int
	for i := lo; i <= hi; i++ {
		out = append(out, i)
	}
	return out
}

func randFloat() float64 { return float64(edakit.Int(0, 9999)) / 10000 }

// eliminate removes vertex v from the (remaining) graph, joining its neighbours into a
// clique, and returns the number of fill edges created.
func (g *graph) eliminate(v int, removed map[int]bool) int {
	var nb []int
	for u := range g.adj[v] {
		if !removed[u] {
			nb = append(nb, u)
		}
	}
	fill := 0
	for i := 0; i < len(nb); i++ {
		for j := i + 1; j < len(nb); j++ {
			if !g.adj[nb[i]][nb[j]] {
				g.add(nb[i], nb[j])
				fill++
			}
		}
	}
	removed[v] = true
	return fill
}

// degree counts the neighbours of v that have not been eliminated.
func (g *graph) degree(v int, removed map[int]bool) int {
	d := 0
	for u := range g.adj[v] {
		if !removed[u] {
			d++
		}
	}
	return d
}

// minDegreeOrder runs minimum degree with lowest-label tie-breaking and returns the order
// and the total fill.
func minDegreeOrder(g0 *graph) ([]int, int) {
	g := g0.clone()
	removed := map[int]bool{}
	var order []int
	total := 0
	for len(order) < g.n {
		best, bd := 0, 1<<30
		for v := 1; v <= g.n; v++ {
			if removed[v] {
				continue
			}
			if d := g.degree(v, removed); d < bd {
				best, bd = v, d
			}
		}
		total += g.eliminate(best, removed)
		order = append(order, best)
	}
	return order, total
}

// fillForOrder eliminates in the given order and returns the total fill.
func fillForOrder(g0 *graph, order []int) int {
	g := g0.clone()
	removed := map[int]bool{}
	total := 0
	for _, v := range order {
		total += g.eliminate(v, removed)
	}
	return total
}

// distinctFloats returns formatted distractors that differ from the correct formatted value.
func distinctFloats(correct string, places int, cands ...float64) []string {
	seen := map[string]bool{correct: true}
	var out []string
	for _, c := range cands {
		s := f(c, places)
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
