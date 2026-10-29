package sta

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// arc is one timing arc u -> v of a drawn timing graph, with its delay in ps.
type arc struct{ u, v, d int }

// dag is a small random timing graph. Node indices are a topological order by construction;
// names are a random permutation of letters, so a learner cannot read the order off the names.
type dag struct {
	n     int
	names []string
	arcs  []arc
}

var nodeLetters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "W"}

// newDAG draws n nodes with nsrc sources. Every later node gets one or two fanin arcs from
// earlier nodes, biased toward recent nodes so the graph has depth. Delays are in [dlo, dhi].
func newDAG(n, nsrc, dlo, dhi int) *dag {
	g := &dag{n: n}
	g.names = edakit.Shuffle(nodeLetters)[:n]
	for v := nsrc; v < n; v++ {
		k := 1
		if v >= 2 && rand.Float64() < 0.55 {
			k = 2
		}
		used := map[int]bool{}
		for i := 0; i < k; i++ {
			lo := v - 3
			if lo < 0 {
				lo = 0
			}
			u := edakit.Int(lo, v-1)
			if i == 1 {
				u = edakit.Int(0, v-1)
			}
			if used[u] {
				continue
			}
			used[u] = true
			g.arcs = append(g.arcs, arc{u, v, edakit.Int(dlo, dhi)})
		}
	}
	// Every source drives something, so no node is isolated.
	for s := 0; s < nsrc && s < n-1; s++ {
		if len(g.fanout(s)) == 0 {
			g.arcs = append(g.arcs, arc{s, edakit.Int(nsrc, n-1), edakit.Int(dlo, dhi)})
		}
	}
	return g
}

func (g *dag) fanin(v int) []arc {
	var out []arc
	for _, a := range g.arcs {
		if a.v == v {
			out = append(out, a)
		}
	}
	return out
}

func (g *dag) fanout(u int) []arc {
	var out []arc
	for _, a := range g.arcs {
		if a.u == u {
			out = append(out, a)
		}
	}
	return out
}

func (g *dag) sources() []int {
	var out []int
	for v := 0; v < g.n; v++ {
		if len(g.fanin(v)) == 0 {
			out = append(out, v)
		}
	}
	return out
}

func (g *dag) sinks() []int {
	var out []int
	for v := 0; v < g.n; v++ {
		if len(g.fanout(v)) == 0 {
			out = append(out, v)
		}
	}
	return out
}

func (g *dag) internal() []int {
	var out []int
	for v := 0; v < g.n; v++ {
		if len(g.fanin(v)) > 0 && len(g.fanout(v)) > 0 {
			out = append(out, v)
		}
	}
	return out
}

// arrivals propagates latest (late=true) or earliest arrival times from the given source arrivals.
func (g *dag) arrivals(src map[int]int, late bool) []int {
	at := make([]int, g.n)
	for v := 0; v < g.n; v++ {
		in := g.fanin(v)
		if len(in) == 0 {
			at[v] = src[v]
			continue
		}
		best := 0
		for i, a := range in {
			t := at[a.u] + a.d
			if i == 0 || (late && t > best) || (!late && t < best) {
				best = t
			}
		}
		at[v] = best
	}
	return at
}

// required propagates required times backward from the given endpoint required times.
func (g *dag) required(sink map[int]int) []int {
	rat := make([]int, g.n)
	for v := g.n - 1; v >= 0; v-- {
		out := g.fanout(v)
		if len(out) == 0 {
			rat[v] = sink[v]
			continue
		}
		best := 0
		for i, a := range out {
			t := rat[a.v] - a.d
			if i == 0 || t < best {
				best = t
			}
		}
		rat[v] = best
	}
	return rat
}

func (g *dag) levels() []int {
	lv := make([]int, g.n)
	for v := 0; v < g.n; v++ {
		for _, a := range g.fanin(v) {
			if lv[a.u]+1 > lv[v] {
				lv[v] = lv[a.u] + 1
			}
		}
	}
	return lv
}

// forwardCone returns the nodes reachable from v, including v.
func (g *dag) forwardCone(v int) map[int]bool {
	seen := map[int]bool{v: true}
	for w := v; w < g.n; w++ {
		if !seen[w] {
			continue
		}
		for _, a := range g.fanout(w) {
			seen[a.v] = true
		}
	}
	return seen
}

// backwardCone returns the nodes that reach u, including u.
func (g *dag) backwardCone(u int) map[int]bool {
	seen := map[int]bool{u: true}
	for w := u; w >= 0; w-- {
		if !seen[w] {
			continue
		}
		for _, a := range g.fanin(w) {
			seen[a.u] = true
		}
	}
	return seen
}

// pathCounts returns, per node, the number of source-to-node paths.
func (g *dag) pathCounts() []int {
	pc := make([]int, g.n)
	for v := 0; v < g.n; v++ {
		in := g.fanin(v)
		if len(in) == 0 {
			pc[v] = 1
			continue
		}
		for _, a := range in {
			pc[v] += pc[a.u]
		}
	}
	return pc
}

// pathDelays returns the delays of every source-to-v path, sorted from worst (largest) down.
func (g *dag) pathDelays(v int) []int {
	var out []int
	var walk func(w, acc int)
	walk = func(w, acc int) {
		in := g.fanin(w)
		if len(in) == 0 {
			out = append(out, acc)
			return
		}
		for _, a := range in {
			walk(a.u, acc+a.d)
		}
	}
	walk(v, 0)
	sort.Sort(sort.Reverse(sort.IntSlice(out)))
	return out
}

// describe lists the arcs in a shuffled order, with delays when withDelay is set.
func (g *dag) describe(withDelay bool) string {
	var parts []string
	for _, a := range edakit.Shuffle(g.arcs) {
		if withDelay {
			parts = append(parts, edakit.Sprintf("%s→%s (%d ps)", g.names[a.u], g.names[a.v], a.d))
		} else {
			parts = append(parts, edakit.Sprintf("%s→%s", g.names[a.u], g.names[a.v]))
		}
	}
	return strings.Join(parts, ", ")
}

func (g *dag) nameList(vs []int) string {
	var parts []string
	for _, v := range vs {
		parts = append(parts, g.names[v])
	}
	return strings.Join(parts, ", ")
}

// valueList renders "A = 3, B = 5" for the given nodes.
func (g *dag) valueList(vs []int, vals []int) string {
	var parts []string
	for _, v := range vs {
		parts = append(parts, edakit.Sprintf("%s = %d", g.names[v], vals[v]))
	}
	return strings.Join(parts, ", ")
}

func allNodes(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func countTrue(m map[int]bool) int {
	n := 0
	for _, b := range m {
		if b {
			n++
		}
	}
	return n
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

func maxInt(xs []int) int {
	m := xs[0]
	for _, x := range xs {
		if x > m {
			m = x
		}
	}
	return m
}

func signed(x int) string {
	if x >= 0 {
		return edakit.Sprintf("+%d", x)
	}
	return edakit.IntStr(x)
}
