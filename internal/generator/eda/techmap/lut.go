package techmap

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// minNodeCut computes the FlowMap test: collapse every node of label p (and the root) into the
// sink and return the minimum number of cuttable nodes (inputs and lower-label gates) that
// separate the primary inputs from the sink.
func minNodeCut(g *lnet, root int, lab []int, p int) int {
	n := g.size()
	S, T := 2*n, 2*n+1
	N := 2*n + 2
	const inf = 1 << 20
	cap := make([][]int, N)
	for i := range cap {
		cap[i] = make([]int, N)
	}
	inX := func(x int) bool { return x == root || (x >= g.npi && lab[x] == p) }
	for x := 0; x < n; x++ {
		if inX(x) {
			continue
		}
		cap[2*x][2*x+1] = 1
		if x < g.npi {
			cap[S][2*x] = inf
		}
	}
	for x := g.npi; x < n; x++ {
		for _, f := range g.fan[x] {
			if inX(f) {
				continue
			}
			if inX(x) {
				cap[2*f+1][T] = inf
			} else {
				cap[2*f+1][2*x] = inf
			}
		}
	}
	flow := 0
	for {
		prev := make([]int, N)
		for i := range prev {
			prev[i] = -1
		}
		prev[S] = S
		queue := []int{S}
		for len(queue) > 0 && prev[T] < 0 {
			u := queue[0]
			queue = queue[1:]
			for w := 0; w < N; w++ {
				if prev[w] < 0 && cap[u][w] > 0 {
					prev[w] = u
					queue = append(queue, w)
				}
			}
		}
		if prev[T] < 0 {
			return flow
		}
		b := inf
		for w := T; w != S; w = prev[w] {
			if cap[prev[w]][w] < b {
				b = cap[prev[w]][w]
			}
		}
		for w := T; w != S; w = prev[w] {
			cap[prev[w]][w] -= b
			cap[w][prev[w]] += b
		}
		flow += b
	}
}

func labelList(g *lnet, lab []int, skip int) string {
	var s []string
	for x := g.npi; x < g.size(); x++ {
		if x != skip {
			s = append(s, edakit.Sprintf("%s: %d", g.name(x), lab[x]))
		}
	}
	return strings.Join(s, ", ")
}

func flowmap(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	for {
		K := 3
		var g *lnet
		switch k {
		case 1:
			g = randCone(edakit.Int(3, 4), edakit.Int(4, 5), 2, 3, false)
		case 2:
			K = edakit.Int(3, 4)
			g = randCone(edakit.Int(3, 5), edakit.Scaled(ctx, 4, 7), 3, 4, false)
		default:
			K = edakit.Int(3, 4)
			g = randCone(edakit.Int(4, 5), edakit.Scaled(ctx, 5, 8), 3, 5, false)
		}
		root := g.size() - 1
		cuts := enumCuts(g, K)
		lab, best := depthLabels(g, cuts)
		p := 0
		for _, f := range g.fan[root] {
			if lab[f] > p {
				p = lab[f]
			}
		}
		if p == 0 || lab[root] < 2 {
			continue
		}
		mc := minNodeCut(g, root, lab, p)
		if (mc <= K) != (lab[root] == p) {
			continue // cannot happen; guards the generator against a modelling slip
		}
		q := edakit.Sprintf("A combinational network over inputs %s is %s. It is mapped into %d-input LUTs with unit LUT delay; primary inputs have label 0.", g.inputs(), g.describe(false), K)
		facts := map[string]string{"label": edakit.IntStr(lab[root]), "p": edakit.IntStr(p), "min_cut": edakit.IntStr(mc)}
		verdict := edakit.Sprintf("p + 1 = %d", p+1)
		if mc <= K {
			verdict = edakit.Sprintf("p = %d", p)
		}
		if k == 2 {
			q += edakit.Sprintf(" The FlowMap labels of the other gates are %s. For node %s, p is the largest fanin label. Collapse %s together with every node of label p in its cone into one sink, and find a minimum node cut separating the primary inputs from that sink (inputs and gates of label below p may be cut). What is the size of that minimum cut?", labelList(g, lab, root), g.name(root), g.name(root))
			return edakit.Numeric(q, edakit.IntStr(mc),
				edakit.Sprintf("Here p = %d. After collapsing, the minimum node cut has %d nodes, found as a maximum flow with unit node capacities. Since %d %s %d, the label of %s is %s.", p, mc, mc, leq(mc <= K), K, g.name(root), verdict),
				facts)
		}
		return edakit.Numeric(q+edakit.Sprintf(" What is the FlowMap label (optimal LUT depth) of node %s?", g.name(root)),
			edakit.IntStr(lab[root]),
			edakit.Sprintf("Labels bottom-up: %s, and %s gets %d with cut %s. Its largest fanin label is p = %d, and the minimum cut after collapsing the label-p nodes has %d nodes, so the label is %s.", labelList(g, lab, root), g.name(root), lab[root], g.cutString(best[root]), p, mc, verdict),
			facts)
	}
}

func leq(b bool) string {
	if b {
		return "<="
	}
	return ">"
}

// lutMapping is a depth-optimal cover of a multi-output network.
type lutMapping struct {
	g       *lnet
	K       int
	lab     []int
	best    []uint32
	outs    []int
	cover   []int // LUT roots, ascending
	inCover map[int]bool
}

func randLutMapping(ctx generator.GeneratorContext) lutMapping {
	for {
		K := edakit.Int(3, 4)
		npi := 4
		g := randNet(npi, edakit.Scaled(ctx, 6, 9), 3, false)
		fo := make([]int, g.size())
		for x := g.npi; x < g.size(); x++ {
			for _, f := range g.fan[x] {
				fo[f]++
			}
		}
		var outs []int
		for x := g.npi; x < g.size(); x++ {
			if fo[x] == 0 {
				outs = append(outs, x)
			}
		}
		if len(outs) > 3 {
			continue
		}
		cuts := enumCuts(g, K)
		lab, best := depthLabels(g, cuts)
		in := map[int]bool{}
		stack := append([]int{}, outs...)
		for len(stack) > 0 {
			r := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if in[r] {
				continue
			}
			in[r] = true
			for _, l := range leaves(best[r], g.size()) {
				if l >= g.npi && !in[l] {
					stack = append(stack, l)
				}
			}
		}
		cover := sortedKeys(in)
		depth := 0
		for _, o := range outs {
			if lab[o] > depth {
				depth = lab[o]
			}
		}
		if len(cover) < 3 || depth < 2 {
			continue
		}
		return lutMapping{g: g, K: K, lab: lab, best: best, outs: outs, cover: cover, inCover: in}
	}
}

func (m lutMapping) names(xs []int) string {
	var s []string
	for _, x := range xs {
		s = append(s, m.g.name(x))
	}
	return strings.Join(s, ", ")
}

func (m lutMapping) bestCuts() string {
	var s []string
	for x := m.g.npi; x < m.g.size(); x++ {
		s = append(s, m.g.name(x)+": "+m.g.cutString(m.best[x]))
	}
	return strings.Join(s, "; ")
}

func (m lutMapping) coverText() string {
	var s []string
	for _, r := range m.cover {
		s = append(s, "the LUT rooted at "+m.g.name(r)+" with inputs "+m.g.cutString(m.best[r]))
	}
	return strings.Join(s, "; ")
}

// coneOf lists the gates a LUT rooted at r with cut c absorbs, r included.
func (m lutMapping) coneOf(r int) []int {
	cut := m.best[r]
	seen := map[int]bool{}
	var walk func(x int)
	walk = func(x int) {
		if seen[x] || cut&(1<<uint(x)) != 0 || x < m.g.npi {
			return
		}
		seen[x] = true
		for _, f := range m.g.fan[x] {
			walk(f)
		}
	}
	walk(r)
	return sortedKeys(seen)
}

func coverSelection(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	m := randLutMapping(ctx)
	q := edakit.Sprintf("A network over inputs %s is %s. Its primary outputs are %s. A depth-optimal %d-LUT mapper has chosen these best cuts: %s.",
		m.g.inputs(), m.g.describe(false), m.names(m.outs), m.K, m.bestCuts())
	facts := map[string]string{"luts": edakit.IntStr(len(m.cover))}
	switch k {
	case 1:
		return edakit.Numeric(q+" The cover is selected backward: start from the outputs, implement each needed node as a LUT over its best cut, and mark that cut's gate leaves as needed. How many LUTs does the mapping use?",
			edakit.IntStr(len(m.cover)),
			edakit.Sprintf("Starting from %s and following best-cut leaves gives LUTs at %s: %d LUTs. Gates not reached are absorbed inside other LUTs.", m.names(m.outs), m.names(m.cover), len(m.cover)),
			facts)
	case 2:
		count := map[int]int{}
		for _, r := range m.cover {
			for _, x := range m.coneOf(r) {
				count[x]++
			}
		}
		var dup []int
		for x, c := range count {
			if c >= 2 {
				dup = append(dup, x)
			}
		}
		sort.Ints(dup)
		why := edakit.Sprintf("The selected LUTs are at %s. ", m.names(m.cover))
		for _, r := range m.cover {
			why += edakit.Sprintf("%s absorbs {%s}. ", m.g.name(r), m.names(m.coneOf(r)))
		}
		why += edakit.Sprintf("Gates inside two or more LUTs: %d.", len(dup))
		return edakit.Numeric(q+" Select the cover backward from the outputs. How many gates of the original network end up inside two or more of the selected LUTs (that is, are duplicated)?",
			edakit.IntStr(len(dup)), why, map[string]string{"luts": edakit.IntStr(len(m.cover)), "duplicated": edakit.IntStr(len(dup))})
	}
	depth := 0
	for _, o := range m.outs {
		if m.lab[o] > depth {
			depth = m.lab[o]
		}
	}
	return edakit.Numeric(q+" Select the cover backward from the outputs. How many LUT levels does the deepest output of the mapped network have?",
		edakit.IntStr(depth),
		edakit.Sprintf("With unit LUT delay the depth of a LUT is one more than the deepest of its LUT inputs. The outputs %s reach depth %d.", m.names(m.outs), depth),
		map[string]string{"depth": edakit.IntStr(depth)})
}

func lutAreaRecovery(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k == 3 {
		R := edakit.Int(2, 4)
		type cand struct {
			d  int
			af float64
		}
		for {
			var cs []cand
			used := map[int]bool{}
			for i := 0; i < 4; i++ {
				a := edakit.Int(10, 60)
				if used[a] {
					continue
				}
				used[a] = true
				cs = append(cs, cand{edakit.Int(R-1, R+1), float64(a) / 10})
			}
			bestAll, bestOK := 1e9, 1e9
			for _, c := range cs {
				if c.af < bestAll {
					bestAll = c.af
				}
				if c.d <= R && c.af < bestOK {
					bestOK = c.af
				}
			}
			if bestOK > 1e8 || bestOK == bestAll {
				continue
			}
			var desc []string
			for i, c := range cs {
				desc = append(desc, edakit.Sprintf("cut %d (depth %d, area flow %s)", i+1, c.d, edakit.Num(c.af, 1)))
			}
			q := edakit.Sprintf("During area recovery, node n has required depth %d. Its candidate cuts are %s. Which area flow does n end up with when the mapper takes the smallest-area-flow cut that still meets the required depth?", R, strings.Join(desc, "; "))
			return edakit.Numeric(q, edakit.Num(bestOK, 1),
				edakit.Sprintf("Cuts deeper than %d would delay the critical outputs, so the cheapest cut of area flow %s is not allowed. The cheapest legal cut has area flow %s.", R, edakit.Num(bestAll, 1), edakit.Num(bestOK, 1)),
				map[string]string{"area_flow": edakit.Num(bestOK, 1)})
		}
	}
	for {
		m := randLutMapping(ctx)
		depth := 0
		for _, o := range m.outs {
			if m.lab[o] > depth {
				depth = m.lab[o]
			}
		}
		D := depth + edakit.Int(0, 1)
		isOut := map[int]bool{}
		for _, o := range m.outs {
			isOut[o] = true
		}
		req := map[int]int{}
		for _, r := range m.cover {
			req[r] = 1 << 20
			if isOut[r] {
				req[r] = D
			}
		}
		for i := len(m.cover) - 1; i >= 0; i-- {
			r := m.cover[i]
			for _, l := range leaves(m.best[r], m.g.size()) {
				if l >= m.g.npi && req[r]-1 < req[l] {
					req[l] = req[r] - 1
				}
			}
		}
		var inner []int
		for _, r := range m.cover {
			if !isOut[r] {
				inner = append(inner, r)
			}
		}
		if len(inner) == 0 {
			continue
		}
		x := edakit.Pick(inner)
		q := edakit.Sprintf("A %d-LUT mapping over inputs %s consists of %s. The primary outputs are %s and must be ready by depth %d, with unit LUT delay.",
			m.K, m.g.inputs(), m.coverText(), m.names(m.outs), D)
		if k == 1 {
			return edakit.Numeric(q+" What is the required depth of the LUT rooted at "+m.g.name(x)+"?", edakit.IntStr(req[x]),
				edakit.Sprintf("Required depths propagate backward: an output gets %d, and each LUT input gets one less than the smallest required depth among the LUTs it feeds. %s gets %d.", D, m.g.name(x), req[x]),
				map[string]string{"required": edakit.IntStr(req[x])})
		}
		slack := req[x] - m.lab[x]
		return edakit.Numeric(q+" What is the slack (required depth minus arrival depth) of the LUT rooted at "+m.g.name(x)+"?", edakit.IntStr(slack),
			edakit.Sprintf("Arrival depth of %s is %d (one more than its deepest LUT input). Its required depth is %d, so the slack is %d. A node with positive slack may switch to a deeper but cheaper cut.", m.g.name(x), m.lab[x], req[x], slack),
			map[string]string{"arrival": edakit.IntStr(m.lab[x]), "required": edakit.IntStr(req[x]), "slack": edakit.IntStr(slack)})
	}
}

func ceilDiv(a, b int) int { return (a + b - 1) / b }

func wideFunctions(ctx generator.GeneratorContext) generator.Problem {
	K := edakit.Int(3, 6)
	switch kp(ctx) {
	case 1:
		m := edakit.Scaled(ctx, K+2, 26)
		fn := edakit.Pick([]string{"AND", "OR", "XOR (parity)"})
		n := ceilDiv(m-1, K-1)
		return edakit.Numeric(edakit.Sprintf("What is the minimum number of %d-input LUTs needed to implement a %d-input %s, with no other logic available?", K, m, fn),
			edakit.IntStr(n),
			edakit.Sprintf("Each LUT combines at most %d signals into one, removing at most %d signals from the set still to be combined. Going from %d signals to 1 takes at least ceil(%d / %d) = %d LUTs, and a tree of LUTs reaches that bound.", K, K-1, m, m-1, K-1, n),
			map[string]string{"luts": edakit.IntStr(n)})
	case 2:
		m := edakit.Scaled(ctx, K+2, 40)
		d, cap := 0, 1
		for cap < m {
			cap *= K
			d++
		}
		return edakit.Numeric(edakit.Sprintf("A %d-input AND is mapped into %d-input LUTs. What is the minimum possible depth, in LUT levels?", m, K),
			edakit.IntStr(d),
			edakit.Sprintf("A LUT tree of depth d can combine at most %d^d inputs. The smallest d with %d^d >= %d is %d.", K, K, m, d),
			map[string]string{"depth": edakit.IntStr(d)})
	}
	m := K + edakit.Int(1, 4)
	n := intPow(2, m-K+1) - 1
	return edakit.Numeric(edakit.Sprintf("An arbitrary %d-input function is mapped into %d-input LUTs by repeated Shannon expansion: f = x f_x + x' f_x', expanding one variable at a time until every cofactor has %d inputs, with each 2-to-1 multiplexer in its own LUT. How many LUTs does this construction use?", m, K, K),
		edakit.IntStr(n),
		edakit.Sprintf("Let N(j) be the count for j inputs. N(%d) = 1 and N(j) = 2 N(j - 1) + 1, so N(%d) = 2^(%d) - 1 = %d. This is an upper bound for an arbitrary function, not the minimum for a particular one.", K, m, m-K+1, n),
		map[string]string{"luts": edakit.IntStr(n)})
}
