package fpga

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func wideFunctions(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		k := edakit.Int(3, 6)
		n := edakit.Scaled(ctx, k+1, 64)
		op := edakit.Pick([]string{"AND", "OR", "XOR"})
		l := ceilDiv(n-1, k-1)
		q := edakit.Sprintf("What is the fewest %d-input LUTs that can implement a %d-input %s?", k, n, op)
		return num(q, edakit.IntStr(l),
			edakit.Sprintf("Each LUT takes k signals and returns one, so it reduces the number of live signals by at most k - 1 = %d. Going from %d signals to 1 needs ceil((%d - 1) / %d) = %d LUTs, and a tree of LUTs achieves it because %s is associative.", k-1, n, n, k-1, l, op),
			map[string]string{"luts": edakit.IntStr(l)})
	case 2:
		k := edakit.Int(3, 6)
		n := edakit.Scaled(ctx, k+1, 300)
		d, c := 0, 1
		for c < n {
			c *= k
			d++
		}
		q := edakit.Sprintf("A %d-input AND is built as a tree of %d-input LUTs. What is the fewest LUT levels on the longest input-to-output path?", n, k)
		return num(q, edakit.IntStr(d),
			edakit.Sprintf("d levels of k-LUTs can combine at most k^d inputs. %d^%d = %d ≥ %d but %d^%d = %d < %d, so the depth is %d.", k, d, c, n, k, d-1, c/k, n, d),
			map[string]string{"depth": edakit.IntStr(d)})
	default:
		k := edakit.Int(4, 6)
		n := edakit.Int(k+1, k+3)
		m := pow2(n - k)
		j := 0
		for (j+1)+pow2(j+1) <= k {
			j++
		}
		r := pow2(j)
		mux, cur := 0, m
		var lv []int
		for cur > 1 {
			nx := ceilDiv(cur, r)
			mux += nx
			lv = append(lv, nx)
			cur = nx
		}
		tot := m + mux
		q := edakit.Sprintf("An arbitrary %d-input function is built from %d-input LUTs by the plain Shannon construction: expand on %d of the inputs so that each cofactor (a function of %d inputs) fills one LUT, then select among the cofactors with a tree of multiplexers, each multiplexer packed into one LUT with as many data inputs as its select and data lines allow. How many LUTs does the construction use in total?", n, k, n-k, k)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("There are 2^%d = %d cofactor LUTs. A %d-LUT holds a multiplexer with s select and 2^s data inputs when s + 2^s ≤ %d, so each mux LUT is %d:1. By level, the mux tree needs %s LUTs, %d in all, for %d LUTs. Real mappers usually beat this by sharing logic between cofactors.", n-k, m, k, k, r, joinInts(lv), mux, tot),
			map[string]string{"cofactors": edakit.IntStr(m), "mux_luts": edakit.IntStr(mux)})
	}
}

type gate struct {
	name string
	op   string
	in   [2]string
}

type logicNet struct {
	pis   []string
	gates []gate
}

func newLogicNet(npi, ng int) *logicNet {
	pis := []string{"a", "b", "c", "d", "e", "f", "u", "v"}[:npi]
	n := &logicNet{pis: pis}
	ops := []string{"AND", "OR", "XOR", "NAND"}
	unused := edakit.Shuffle(pis)
	takePI := func() string {
		if len(unused) > 0 {
			x := unused[0]
			unused = unused[1:]
			return x
		}
		return edakit.Pick(pis)
	}
	for i := 0; i < ng; i++ {
		var x, y string
		if i == 0 {
			x = takePI()
		} else {
			lo := i - 2
			if lo < 0 {
				lo = 0
			}
			x = n.gates[edakit.Int(lo, i-1)].name
		}
		for y == "" || y == x {
			r := edakit.Int(0, 9)
			switch {
			case r < 6 || i < 2:
				y = takePI()
			case r < 8:
				y = n.gates[edakit.Int(0, i-1)].name
			default:
				y = edakit.Pick(pis)
			}
		}
		n.gates = append(n.gates, gate{edakit.Sprintf("g%d", i+1), edakit.Pick(ops), [2]string{x, y}})
	}
	return n
}

func (n *logicNet) text() string {
	var p []string
	for _, g := range n.gates {
		p = append(p, edakit.Sprintf("%s = %s(%s, %s)", g.name, g.op, g.in[0], g.in[1]))
	}
	return strings.Join(p, "; ")
}

func cutKey(c []string) string { return strings.Join(c, ",") }

// cuts enumerates every k-feasible cut of every node, the trivial cut included.
func (n *logicNet) cuts(k int) map[string][][]string {
	out := map[string][][]string{}
	for _, p := range n.pis {
		out[p] = [][]string{{p}}
	}
	for _, g := range n.gates {
		seen := map[string]bool{}
		list := [][]string{{g.name}}
		seen[g.name] = true
		for _, c1 := range out[g.in[0]] {
			for _, c2 := range out[g.in[1]] {
				u := union(c1, c2)
				if len(u) > k {
					continue
				}
				key := cutKey(u)
				if !seen[key] {
					seen[key] = true
					list = append(list, u)
				}
			}
		}
		out[g.name] = list
	}
	return out
}

// labels returns the depth-optimal LUT level of every node.
func (n *logicNet) labels(k int) map[string]int {
	cs := n.cuts(k)
	lab := map[string]int{}
	for _, p := range n.pis {
		lab[p] = 0
	}
	for _, g := range n.gates {
		best := 1 << 30
		for _, c := range cs[g.name] {
			if len(c) == 1 && c[0] == g.name {
				continue
			}
			m := 0
			for _, u := range c {
				if lab[u] > m {
					m = lab[u]
				}
			}
			if m+1 < best {
				best = m + 1
			}
		}
		lab[g.name] = best
	}
	return lab
}

func (n *logicNet) levels() map[string]int {
	lv := map[string]int{}
	for _, g := range n.gates {
		lv[g.name] = maxInt(lv[g.in[0]], lv[g.in[1]]) + 1
	}
	return lv
}

func (n *logicNet) labelText(lab map[string]int) string {
	var p []string
	for _, g := range n.gates {
		p = append(p, edakit.Sprintf("%s: %d", g.name, lab[g.name]))
	}
	return strings.Join(p, ", ")
}

func lutCuts(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		for {
			n := newLogicNet(edakit.Int(4, 5), edakit.Scaled(ctx, 3, 5))
			k := 3
			cs := n.cuts(k)
			g := n.gates[edakit.Int(1, len(n.gates)-1)]
			cnt := len(cs[g.name])
			if cnt > 9 || cnt < 3 {
				continue
			}
			var shown []string
			for _, c := range cs[g.name] {
				shown = append(shown, setText(c))
			}
			q := edakit.Sprintf("Gate network over primary inputs %s: %s. How many 3-feasible cuts does %s have, counting the trivial cut {%s} and counting each distinct leaf set once?", strings.Join(n.pis, ", "), n.text(), g.name, g.name)
			return num(q, edakit.IntStr(cnt),
				edakit.Sprintf("Merge every cut of %s with every cut of %s, keep the unions with at most 3 leaves, and add the trivial cut. The cuts of %s are %s.", g.in[0], g.in[1], g.name, strings.Join(shown, ", ")),
				map[string]string{"cuts": edakit.IntStr(cnt)})
		}
	case 2:
		n := newLogicNet(edakit.Int(6, 8), edakit.Scaled(ctx, 5, 8))
		k := edakit.Int(3, 4)
		lab := n.labels(k)
		g := n.gates[len(n.gates)-1]
		q := edakit.Sprintf("Gate network over primary inputs %s: %s. Primary inputs have label 0, and every gate's label is the minimum, over its non-trivial %d-feasible cuts, of 1 plus the largest label among the cut's leaves. What is the label of %s, its LUT depth in a depth-optimal %d-LUT mapping?", strings.Join(n.pis, ", "), n.text(), k, g.name, k)
		return num(q, edakit.IntStr(lab[g.name]),
			edakit.Sprintf("Compute labels in topological order. With k = %d they are %s, so %s has label %d. FlowMap finds the same labels with a max-flow test per node instead of enumerating cuts.", k, n.labelText(lab), g.name, lab[g.name]),
			map[string]string{"label": edakit.IntStr(lab[g.name])})
	default:
		for {
			n := newLogicNet(edakit.Int(6, 8), edakit.Scaled(ctx, 6, 9))
			k := edakit.Int(3, 4)
			lab := n.labels(k)
			lv := n.levels()
			g := n.gates[len(n.gates)-1]
			save := lv[g.name] - lab[g.name]
			if save < 1 && edakit.Int(0, 3) > 0 {
				continue
			}
			q := edakit.Sprintf("Gate network over primary inputs %s: %s. Mapping each gate into its own LUT gives %s a depth equal to its gate level. How many LUT levels does a depth-optimal %d-LUT mapping save at %s compared with that one-gate-per-LUT mapping?", strings.Join(n.pis, ", "), n.text(), g.name, k, g.name)
			return num(q, edakit.IntStr(save),
				edakit.Sprintf("The gate level of %s is %d. The depth-optimal labels with k = %d are %s, so %s sits at level %d and the mapping saves %d.", g.name, lv[g.name], k, n.labelText(lab), g.name, lab[g.name], save),
				map[string]string{"gate_level": edakit.IntStr(lv[g.name]), "label": edakit.IntStr(lab[g.name])})
		}
	}
}

func areaFlow(ctx generator.GeneratorContext) generator.Problem {
	af := func() float64 { return float64(edakit.Int(0, 12)) / 4 }
	switch whichKP(ctx) {
	case 1:
		nl := edakit.Int(2, 4)
		var vs []float64
		var ps []string
		sum := 0.0
		for i := 0; i < nl; i++ {
			v := af()
			vs = append(vs, v)
			sum += v
			ps = append(ps, edakit.Num(v, 2))
		}
		fo := edakit.Int(1, 4)
		r := (1 + sum) / float64(fo)
		q := edakit.Sprintf("A LUT mapper uses unit LUT area and area flow AF(v) = (1 + Σ AF(leaf)) / fanout(v). Node v has fanout %d and its chosen cut has %d leaves with area flows %s. What is AF(v)? Round to 2 decimal places.", fo, nl, strings.Join(ps, ", "))
		return num(q, edakit.Num(r, 2),
			edakit.Sprintf("(1 + %s) / %d = %s. Dividing by fanout shares the cone's area among the nodes that will reuse it.", edakit.Num(sum, 2), fo, edakit.Num(r, 2)),
			map[string]string{"leaf_sum": edakit.Num(sum, 2)})
	case 2:
		fo := edakit.Int(1, 3)
		mk := func() ([]float64, float64) {
			n := edakit.Int(2, 4)
			var v []float64
			s := 0.0
			for i := 0; i < n; i++ {
				x := af()
				v = append(v, x)
				s += x
			}
			return v, s
		}
		a, sa := mk()
		b, sb := mk()
		for sa == sb {
			b, sb = mk()
		}
		fa, fb := (1+sa)/float64(fo), (1+sb)/float64(fo)
		fmtv := func(v []float64) string {
			var p []string
			for _, x := range v {
				p = append(p, edakit.Num(x, 2))
			}
			return strings.Join(p, ", ")
		}
		best := fa
		if fb < fa {
			best = fb
		}
		q := edakit.Sprintf("During area recovery, node v (fanout %d) chooses between two cuts that meet its required time. Cut 1 has leaves with area flows %s; cut 2 has leaves with area flows %s. With unit LUT area, AF = (1 + Σ AF(leaf)) / fanout. What area flow does v record after choosing the better cut? Round to 2 decimal places.", fo, fmtv(a), fmtv(b))
		return num(q, edakit.Num(best, 2),
			edakit.Sprintf("Cut 1: (1 + %s) / %d = %s. Cut 2: (1 + %s) / %d = %s. The mapper keeps the smaller, %s.", edakit.Num(sa, 2), fo, edakit.Num(fa, 2), edakit.Num(sb, 2), fo, edakit.Num(fb, 2), edakit.Num(best, 2)),
			map[string]string{"af1": edakit.Num(fa, 2), "af2": edakit.Num(fb, 2)})
	default:
		f1, f2, f3 := edakit.Int(1, 4), edakit.Int(1, 3), edakit.Int(1, 2)
		x := float64(edakit.Int(0, 8)) / 4
		a1 := 1 / float64(f1)
		a2 := (1 + a1 + x) / float64(f2)
		a3 := (1 + a2 + a1) / float64(f3)
		q := edakit.Sprintf("Unit-area LUTs, AF(v) = (1 + Σ AF(leaf)) / fanout(v), and primary inputs have AF 0. Node n1 has fanout %d and a cut of three primary inputs. Node n2 has fanout %d and the cut {n1, m, p}, where m has AF %s and p is a primary input. Node n3 has fanout %d and the cut {n2, n1, q}, with q a primary input. What is AF(n3)? Round to 2 decimal places.", f1, f2, edakit.Num(x, 2), f3)
		return num(q, edakit.Num(a3, 2),
			edakit.Sprintf("AF(n1) = 1 / %d = %s. AF(n2) = (1 + %s + %s) / %d = %s. AF(n3) = (1 + %s + %s) / %d = %s. Note n1 is counted again at n3: area flow is a heuristic that can double-count reconvergent logic.", f1, edakit.Num(a1, 4), edakit.Num(a1, 4), edakit.Num(x, 2), f2, edakit.Num(a2, 4), edakit.Num(a2, 4), edakit.Num(a1, 4), f3, edakit.Num(a3, 2)),
			map[string]string{"af_n1": edakit.Num(a1, 4), "af_n2": edakit.Num(a2, 4)})
	}
}
