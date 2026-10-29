// Package fv generates formal-verification problems: SAT solving, BDDs, equivalence checking,
// model checking, bounded model checking, induction, IC3 and abstraction, SMT, assertions and
// the agent-first frontier of formal engines.
package fv

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// ---------- explicit transition systems ----------

// graph is an explicit transition system over states 0..n-1.
type graph struct {
	n    int
	succ [][]int
}

func newGraph(n int) *graph { return &graph{n: n, succ: make([][]int, n)} }

func (g *graph) add(a, b int) {
	for _, x := range g.succ[a] {
		if x == b {
			return
		}
	}
	g.succ[a] = append(g.succ[a], b)
	sort.Ints(g.succ[a])
}

// randomGraph draws n states with out-degree between minOut and maxOut.
func randomGraph(n, minOut, maxOut int) *graph {
	g := newGraph(n)
	for s := 0; s < n; s++ {
		k := edakit.Int(minOut, maxOut)
		for i := 0; i < k; i++ {
			g.add(s, rand.Intn(n))
		}
	}
	return g
}

// bfs returns the BFS distance from the set src to every state (-1 if unreachable).
func (g *graph) bfs(src []int) []int {
	d := make([]int, g.n)
	for i := range d {
		d[i] = -1
	}
	var q []int
	for _, s := range src {
		if d[s] < 0 {
			d[s] = 0
			q = append(q, s)
		}
	}
	for len(q) > 0 {
		s := q[0]
		q = q[1:]
		for _, t := range g.succ[s] {
			if d[t] < 0 {
				d[t] = d[s] + 1
				q = append(q, t)
			}
		}
	}
	return d
}

// reverse returns the graph with every edge reversed.
func (g *graph) reverse() *graph {
	r := newGraph(g.n)
	for s := 0; s < g.n; s++ {
		for _, t := range g.succ[s] {
			r.add(t, s)
		}
	}
	return r
}

// stateName names state i as s<i>.
func stateName(i int) string { return edakit.Sprintf("s%d", i) }

// edgesText lists the transition relation as "s0 → s1, s2; s1 → s0; ...".
func (g *graph) edgesText() string {
	var parts []string
	for s := 0; s < g.n; s++ {
		if len(g.succ[s]) == 0 {
			parts = append(parts, stateName(s)+" → (none)")
			continue
		}
		var ts []string
		for _, t := range g.succ[s] {
			ts = append(ts, stateName(t))
		}
		parts = append(parts, stateName(s)+" → "+strings.Join(ts, ", "))
	}
	return strings.Join(parts, "; ")
}

// setText writes a state set as {s0, s3}.
func setText(xs []int) string {
	var ns []string
	for _, x := range xs {
		ns = append(ns, stateName(x))
	}
	return "{" + strings.Join(ns, ", ") + "}"
}

// members returns the sorted indices i with in[i].
func members(in []bool) []int {
	var out []int
	for i, b := range in {
		if b {
			out = append(out, i)
		}
	}
	return out
}

func count(in []bool) int {
	n := 0
	for _, b := range in {
		if b {
			n++
		}
	}
	return n
}

// randomSubset returns k distinct values of 0..n-1, sorted.
func randomSubset(n, k int) []int {
	p := rand.Perm(n)[:k]
	sort.Ints(p)
	return p
}

// bitsN writes state s as an n-bit vector, most significant bit first.
func bitsN(s, n int) string {
	var b strings.Builder
	for i := n - 1; i >= 0; i-- {
		if s>>i&1 == 1 {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}

func boolsFrom(n int, xs []int) []bool {
	out := make([]bool, n)
	for _, x := range xs {
		out[x] = true
	}
	return out
}

func joinInts(xs []int, sep string) string {
	var s []string
	for _, x := range xs {
		s = append(s, edakit.IntStr(x))
	}
	return strings.Join(s, sep)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func binom(n, k int) int {
	if k < 0 || k > n {
		return 0
	}
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
