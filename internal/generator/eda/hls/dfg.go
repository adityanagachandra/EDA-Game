// Package hls generates high-level-synthesis problems: data-flow graphs and their schedules,
// allocation and binding, loop pipelining and modulo scheduling, memory partitioning,
// dataflow, co-design and design-space exploration, and the agent-first frontier of HLS.
package hls

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const (
	opAdd = 0
	opMul = 1
)

// node is one operation of a straight-line data-flow graph.
type node struct {
	kind  int
	preds []int    // indices of the operations whose results it reads
	args  []string // the two operand names as printed
}

// dfg is a random data-flow graph in three-address form: n1 = a + b; n2 = n1 × c; ...
type dfg struct {
	ops []node
}

func opName(i int) string { return edakit.Sprintf("n%d", i+1) }

// newDFG draws n operations. The first two read only primary inputs; later ones read one or
// two earlier results, biased towards recent ones so the graph has depth as well as width.
func newDFG(n int, mulShare float64) *dfg {
	inputs := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	g := &dfg{}
	for i := 0; i < n; i++ {
		nd := node{kind: opAdd}
		if rand.Float64() < mulShare {
			nd.kind = opMul
		}
		if i >= 2 {
			k := 1
			if rand.Float64() < 0.55 {
				k = 2
			}
			picked := map[int]bool{}
			for len(picked) < k {
				lo := i - 4
				if lo < 0 {
					lo = 0
				}
				var p int
				if rand.Float64() < 0.75 {
					p = edakit.Int(lo, i-1)
				} else {
					p = edakit.Int(0, i-1)
				}
				picked[p] = true
			}
			for p := range picked {
				nd.preds = append(nd.preds, p)
			}
			sort.Ints(nd.preds)
		}
		for _, p := range nd.preds {
			nd.args = append(nd.args, opName(p))
		}
		for len(nd.args) < 2 {
			nd.args = append(nd.args, edakit.Pick(inputs))
		}
		if len(nd.preds) == 1 && rand.Intn(2) == 0 {
			nd.args[0], nd.args[1] = nd.args[1], nd.args[0]
		}
		g.ops = append(g.ops, nd)
	}
	return g
}

func (g *dfg) text() string {
	var parts []string
	for i, nd := range g.ops {
		sym := "+"
		if nd.kind == opMul {
			sym = "×"
		}
		parts = append(parts, edakit.Sprintf("%s = %s %s %s", opName(i), nd.args[0], sym, nd.args[1]))
	}
	return strings.Join(parts, "; ")
}

func (g *dfg) succs() [][]int {
	out := make([][]int, len(g.ops))
	for i, nd := range g.ops {
		for _, p := range nd.preds {
			out[p] = append(out[p], i)
		}
	}
	return out
}

func (g *dfg) edges() int {
	e := 0
	for _, nd := range g.ops {
		e += len(nd.preds)
	}
	return e
}

func (g *dfg) count(kind int) int {
	c := 0
	for _, nd := range g.ops {
		if nd.kind == kind {
			c++
		}
	}
	return c
}

// lat returns the latency of each operation for adder and multiplier latencies la and lm.
func (g *dfg) lat(la, lm int) []int {
	out := make([]int, len(g.ops))
	for i, nd := range g.ops {
		out[i] = la
		if nd.kind == opMul {
			out[i] = lm
		}
	}
	return out
}

// asap returns 1-based ASAP start steps.
func (g *dfg) asap(lat []int) []int {
	s := make([]int, len(g.ops))
	for i, nd := range g.ops {
		s[i] = 1
		for _, p := range nd.preds {
			if s[p]+lat[p] > s[i] {
				s[i] = s[p] + lat[p]
			}
		}
	}
	return s
}

// length is the number of control steps a schedule occupies.
func length(start, lat []int) int {
	m := 0
	for i := range start {
		if start[i]+lat[i]-1 > m {
			m = start[i] + lat[i] - 1
		}
	}
	return m
}

// alap returns 1-based ALAP start steps under a latency bound of L steps.
func (g *dfg) alap(lat []int, L int) []int {
	su := g.succs()
	s := make([]int, len(g.ops))
	for i := len(g.ops) - 1; i >= 0; i-- {
		s[i] = L - lat[i] + 1
		for _, c := range su[i] {
			if s[c]-lat[i] < s[i] {
				s[i] = s[c] - lat[i]
			}
		}
	}
	return s
}

// heights is the latency-weighted longest path from each operation to a sink, including itself.
func (g *dfg) heights(lat []int) []int {
	su := g.succs()
	h := make([]int, len(g.ops))
	for i := len(g.ops) - 1; i >= 0; i-- {
		best := 0
		for _, c := range su[i] {
			if h[c] > best {
				best = h[c]
			}
		}
		h[i] = lat[i] + best
	}
	return h
}

// listSchedule runs resource-constrained list scheduling with non-pipelined units. units[k]
// is the number of units of kind k. Priority is height, ties to the lower index.
func (g *dfg) listSchedule(lat []int, units [2]int) []int {
	n := len(g.ops)
	h := g.heights(lat)
	start := make([]int, n)
	done := 0
	for t := 1; done < n && t < 1000; t++ {
		var ready []int
		for i, nd := range g.ops {
			if start[i] != 0 {
				continue
			}
			ok := true
			for _, p := range nd.preds {
				if start[p] == 0 || start[p]+lat[p] > t {
					ok = false
					break
				}
			}
			if ok {
				ready = append(ready, i)
			}
		}
		sort.SliceStable(ready, func(a, b int) bool {
			if h[ready[a]] != h[ready[b]] {
				return h[ready[a]] > h[ready[b]]
			}
			return ready[a] < ready[b]
		})
		var busy [2]int
		for i := range g.ops {
			if start[i] != 0 && start[i] <= t && t < start[i]+lat[i] {
				busy[g.ops[i].kind]++
			}
		}
		for _, i := range ready {
			k := g.ops[i].kind
			if busy[k] < units[k] {
				busy[k]++
				start[i] = t
				done++
			}
		}
	}
	return start
}

// stepTable renders a schedule as "step 1: n1, n2; step 2: n3".
func stepTable(start []int) string {
	m := 0
	for _, s := range start {
		if s > m {
			m = s
		}
	}
	var parts []string
	for t := 1; t <= m; t++ {
		var in []string
		for i, s := range start {
			if s == t {
				in = append(in, opName(i))
			}
		}
		if len(in) == 0 {
			in = []string{"none"}
		}
		parts = append(parts, edakit.Sprintf("step %d: %s", t, strings.Join(in, ", ")))
	}
	return strings.Join(parts, "; ")
}

func ceilDiv(a, b int) int {
	if b <= 0 {
		return 0
	}
	return (a + b - 1) / b
}

func maxInt(xs ...int) int {
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

func minInt(xs ...int) int {
	m := xs[0]
	for _, x := range xs[1:] {
		if x < m {
			m = x
		}
	}
	return m
}

func joinInts(xs []int) string {
	var s []string
	for _, x := range xs {
		s = append(s, edakit.IntStr(x))
	}
	return strings.Join(s, ", ")
}

func kindName(k int) string {
	if k == opMul {
		return "multiplication"
	}
	return "addition"
}

func unitName(k int) string {
	if k == opMul {
		return "multiplier"
	}
	return "adder"
}

// whichKP resolves the knowledge point: the requested one, or a uniform pick when unguided.
func whichKP(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

func num(q, a, why string, facts map[string]string) generator.Problem {
	return edakit.Numeric(q, a, why, facts)
}
