package route

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type edgeDC struct{ d, c int }

func randEdges(n int) []edgeDC {
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		es := make([]edgeDC, n)
		over := 0
		for i := range es {
			c := edakit.Int(4, 16)
			d := c + edakit.Int(-5, 4)
			if d < 0 {
				d = 0
			}
			es[i] = edgeDC{d, c}
			if d > c {
				over++
			}
		}
		if over >= 2 && over < n {
			return es
		}
	}
	es := make([]edgeDC, n)
	for i := range es {
		es[i] = edgeDC{5 + i%3, 7}
	}
	es[0], es[1] = edgeDC{9, 7}, edgeDC{8, 6}
	return es
}

func edgesText(es []edgeDC) string {
	var parts []string
	for i, e := range es {
		parts = append(parts, edakit.Sprintf("e%d demand %d, capacity %d", i+1, e.d, e.c))
	}
	return strings.Join(parts, "; ")
}

func overflow(es []edgeDC) (tof, mof, n int) {
	for _, e := range es {
		if o := e.d - e.c; o > 0 {
			tof += o
			n++
			mof = maxI(mof, o)
		}
	}
	return
}

// congestion: overflowed-edge count (KP 1), total overflow (KP 2), total overflow after moving
// demand between edges (KP 3).
func congestion(ctx generator.GeneratorContext) generator.Problem {
	es := randEdges(edakit.Scaled(ctx, 6, 10))
	tof, mof, n := overflow(es)
	intro := "A congestion map lists, for each gcell edge, the routed demand (wires crossing it) and the capacity (usable tracks): " + edgesText(es) + ". Overflow on an edge is max(0, demand - capacity). "
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(intro+"How many edges are overflowed?", edakit.IntStr(n),
			edakit.Sprintf("An edge overflows when its demand exceeds its capacity. %d edges do, and the worst overflows by %d.", n, mof),
			map[string]string{"overflowed": edakit.IntStr(n), "max_overflow": edakit.IntStr(mof)})
	case 2:
		return edakit.Numeric(intro+"What is the total overflow, the sum over all edges, which global-routing contests report as TOF?", edakit.IntStr(tof),
			edakit.Sprintf("Summing max(0, d - c) edge by edge gives %d. Under-used edges do not offset over-used ones: the max(0, ...) clips them to zero. The maximum overflow is %d.", tof, mof),
			map[string]string{"tof": edakit.IntStr(tof), "mof": edakit.IntStr(mof)})
	default:
		var over, rest []int
		for i, e := range es {
			if e.d > e.c {
				over = append(over, i)
			} else {
				rest = append(rest, i)
			}
		}
		from := edakit.Pick(over)
		to := edakit.Pick(append(rest, over...))
		if to == from {
			to = rest[0]
		}
		k := edakit.Int(1, minI(3, es[from].d))
		es2 := append([]edgeDC(nil), es...)
		es2[from].d -= k
		es2[to].d += k
		t2, _, _ := overflow(es2)
		return edakit.Numeric(intro+edakit.Sprintf("Rip-up and reroute moves %d wires off e%d and onto e%d. What is the total overflow afterward?", k, from+1, to+1), edakit.IntStr(t2),
			edakit.Sprintf("e%d becomes %d over %d and e%d becomes %d over %d. Recomputing the sum of max(0, d - c) gives %d, against %d before. Moving demand helps only if the receiving edge has slack.", from+1, es2[from].d, es2[from].c, to+1, es2[to].d, es2[to].c, t2, tof),
			map[string]string{"before": edakit.IntStr(tof), "after": edakit.IntStr(t2)})
	}
}

type twoPin struct{ a, b pt }

// lUses reports whether the horizontal-first and vertical-first L shapes of n use the
// horizontal edge from (x, y) to (x+1, y).
func lUses(n twoPin, x, y int) (bool, bool) {
	x0, x1 := minI(n.a.x, n.b.x), maxI(n.a.x, n.b.x)
	if x < x0 || x >= x1 {
		return false, false
	}
	// Horizontal first runs along a's row; vertical first runs along b's row.
	return y == n.a.y, y == n.b.y
}

func netsText(ns []twoPin) string {
	var parts []string
	for i, n := range ns {
		parts = append(parts, edakit.Sprintf("net %d from %s to %s", i+1, ptText(n.a), ptText(n.b)))
	}
	return strings.Join(parts, "; ")
}

// probabilisticEstimate: expected demand of one net (KP 1), of several (KP 2), expected overflow
// (KP 3).
func probabilisticEstimate(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	nNets := 1
	if k >= 2 {
		nNets = edakit.Scaled(ctx, 3, 5)
	}
	for try := 0; ; try++ {
		var ns []twoPin
		for len(ns) < nNets {
			a := pt{edakit.Int(0, 6), edakit.Int(0, 4)}
			b := pt{(a.x + edakit.Int(1, 6)) % 7, edakit.Int(0, 4)}
			ns = append(ns, twoPin{a, b})
		}
		x, y := edakit.Int(0, 5), edakit.Int(0, 4)
		if try >= generator.MaxDrawAttempts {
			// Known-good fallback: every net touches the edge (1, 0)-(2, 0).
			ns = []twoPin{{pt{0, 0}, pt{3, 2}}, {pt{1, 0}, pt{4, 3}}, {pt{0, 0}, pt{5, 0}}, {pt{2, 3}, pt{0, 0}}, {pt{1, 1}, pt{6, 0}}}[:nNets]
			x, y = 1, 0
		}
		probs := make([]float64, len(ns))
		exp := 0.0
		touched := 0
		for i, n := range ns {
			u1, u2 := lUses(n, x, y)
			p := 0.0
			if n.a.y == n.b.y {
				if u1 {
					p = 1
				}
			} else {
				if u1 {
					p += 0.5
				}
				if u2 {
					p += 0.5
				}
			}
			probs[i] = p
			exp += p
			if p > 0 {
				touched++
			}
		}
		if exp == 0 || (k >= 2 && touched < 2) {
			continue
		}
		intro := edakit.Sprintf("A fast congestion estimator assumes each two-pin net takes one of its two L-shaped routes with probability 1/2 each (a net in a single row or column has only its straight route). Nets: %s. Consider the horizontal gcell edge between (%d, %d) and (%d, %d). ", netsText(ns), x, y, x+1, y)
		if k <= 2 {
			ans := edakit.Num(exp, 2)
			return edakit.Numeric(intro+"What is its expected demand? Round to 2 decimal places.", ans,
				edakit.Sprintf("Each net adds 1/2 for every L shape through the edge (1 for a straight net through it). The per-net contributions are %s, which sum to %s.", floatsText(probs), ans),
				map[string]string{"expected": ans})
		}
		c := int(math.Floor(exp))
		if c < 1 {
			c = 1
		}
		// Enumerate the equally likely L choices of the nets that have two shapes.
		base := 0
		var bent []twoPin
		for _, n := range ns {
			if n.a.y == n.b.y {
				if u1, _ := lUses(n, x, y); u1 {
					base++
				}
				continue
			}
			bent = append(bent, n)
		}
		eo := 0.0
		pr := math.Pow(0.5, float64(len(bent)))
		for mask := 0; mask < 1<<len(bent); mask++ {
			d := base
			for i, n := range bent {
				u1, u2 := lUses(n, x, y)
				if (mask&(1<<i) != 0 && u1) || (mask&(1<<i) == 0 && u2) {
					d++
				}
			}
			if d > c {
				eo += pr * float64(d-c)
			}
		}
		ans := edakit.Num(eo, 3)
		naive := math.Max(0, exp-float64(c))
		return edakit.Numeric(intro+edakit.Sprintf("The edge has capacity %d. Treating the nets' choices as independent, what is the expected overflow max(0, demand - %d) on this edge? Round to 3 decimal places.", c, c), ans,
			edakit.Sprintf("Enumerate the equally likely combinations of L choices, compute demand in each and average the overflow: %s. Plugging the expected demand %s into max(0, d - c) instead would give %s; overflow is convex in demand, so the expectation of overflow is at least the overflow of the expectation.", ans, edakit.Num(exp, 2), edakit.Num(naive, 3)),
			map[string]string{"expected_demand": edakit.Num(exp, 2), "expected_overflow": ans})
	}
}

func floatsText(xs []float64) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.Num(x, 2)
	}
	return strings.Join(parts, ", ")
}

// ilpFormulation: binary variables (KP 1), constraints (KP 2), the optimum of a tiny instance
// by enumeration (KP 3).
func ilpFormulation(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		m := edakit.Scaled(ctx, 4, 8)
		ks := make([]int, m)
		sum := 0
		for i := range ks {
			ks[i] = edakit.Int(2, 6)
			sum += ks[i]
		}
		q := edakit.Sprintf("A global-routing integer program gives each net a set of candidate Steiner trees and one binary variable per candidate, equal to 1 if that tree is chosen. %d nets have %s candidates respectively. How many binary variables does the program have?", m, intsText(ks))
		return edakit.Numeric(q, edakit.IntStr(sum),
			edakit.Sprintf("One variable per candidate: %s sum to %d.", intsText(ks), sum),
			map[string]string{"vars": edakit.IntStr(sum)})
	case 2:
		X, Y := edakit.Int(3, 30), edakit.Int(3, 30)
		m := edakit.Int(5, 500)
		E := (X-1)*Y + X*(Y-1)
		q := edakit.Sprintf("A global-routing integer program has, for each of %d nets, a constraint that exactly one candidate tree is chosen, and for each edge of a %d by %d gcell grid, a constraint that the trees using it do not exceed its capacity. How many constraints does it have, not counting the 0-1 bounds?", m, X, Y)
		return edakit.Numeric(q, edakit.IntStr(m+E),
			edakit.Sprintf("%d selection constraints plus (%d - 1)(%d) + %d(%d - 1) = %d capacity constraints gives %d.", m, X, Y, X, Y, E, m+E),
			map[string]string{"edges": edakit.IntStr(E)})
	default:
		for try := 0; ; try++ {
			nE := 5
			caps := make([]int, nE)
			for i := range caps {
				caps[i] = edakit.Int(1, 2)
			}
			type cand struct {
				wl    int
				edges []int
			}
			nets := make([][]cand, 3)
			var desc []string
			for i := range nets {
				var cs []string
				for j := 0; j < 2; j++ {
					es := edakit.Shuffle([]int{0, 1, 2, 3, 4})[:edakit.Int(1, 3)]
					sort.Ints(es)
					wl := edakit.Int(3, 12)
					nets[i] = append(nets[i], cand{wl, es})
					var en []string
					for _, e := range es {
						en = append(en, edakit.Sprintf("e%d", e+1))
					}
					cs = append(cs, edakit.Sprintf("tree %c (wirelength %d, uses %s)", 'A'+j, wl, strings.Join(en, ", ")))
				}
				desc = append(desc, edakit.Sprintf("net %d: %s", i+1, strings.Join(cs, " or ")))
			}
			best, feas := math.MaxInt32, 0
			for mask := 0; mask < 8; mask++ {
				use := make([]int, nE)
				wl := 0
				for i := 0; i < 3; i++ {
					c := nets[i][(mask>>i)&1]
					wl += c.wl
					for _, e := range c.edges {
						use[e]++
					}
				}
				ok := true
				for e := range use {
					if use[e] > caps[e] {
						ok = false
					}
				}
				if ok {
					feas++
					best = minI(best, wl)
				}
			}
			if (feas == 0 || feas == 8) && try < generator.MaxDrawAttempts {
				continue
			}
			if feas == 0 {
				// Fallback: give every edge room for all three nets.
				for i := range caps {
					caps[i] = 3
				}
				feas, best = 8, 0
				for i := 0; i < 3; i++ {
					best += minI(nets[i][0].wl, nets[i][1].wl)
				}
			}
			var cp []string
			for i, c := range caps {
				cp = append(cp, edakit.Sprintf("e%d: %d", i+1, c))
			}
			q := edakit.Sprintf("A tiny global-routing ILP chooses one candidate tree per net to minimise total wirelength subject to edge capacities (%s). Candidates: %s. What is the optimal total wirelength?", strings.Join(cp, ", "), strings.Join(desc, "; "))
			return edakit.Numeric(q, edakit.IntStr(best),
				edakit.Sprintf("Of the 8 combinations, the number that respect every capacity is %d. The cheapest feasible one has total wirelength %d. Real instances have millions of variables, which is why routers relax or decompose this program instead of solving it exactly.", feas, best),
				map[string]string{"feasible": edakit.IntStr(feas), "optimum": edakit.IntStr(best)})
		}
	}
}

// multicommodityFlow: flow variables (KP 1), conservation constraints (KP 2), the min-max
// utilisation of given flows (KP 3).
func multicommodityFlow(ctx generator.GeneratorContext) generator.Problem {
	X, Y := edakit.Int(3, 20), edakit.Int(3, 20)
	K := edakit.Int(2, 50)
	E := (X-1)*Y + X*(Y-1)
	switch kp(ctx) {
	case 1:
		q := edakit.Sprintf("Global routing is cast as a multicommodity flow: each of %d two-pin nets is a commodity, and each commodity has a flow variable on every directed arc of a %d by %d gcell grid (each undirected edge gives two arcs). How many flow variables are there?", K, X, Y)
		return edakit.Numeric(q, edakit.IntStr(K*2*E),
			edakit.Sprintf("The grid has %d edges, so %d arcs. %d commodities times %d arcs is %d variables.", E, 2*E, K, 2*E, K*2*E),
			map[string]string{"edges": edakit.IntStr(E)})
	case 2:
		q := edakit.Sprintf("In a multicommodity-flow model of global routing, every commodity (one per two-pin net) has a flow-conservation constraint at every vertex: outflow minus inflow equals +1 at its source, -1 at its sink and 0 elsewhere. With %d commodities on a %d by %d gcell grid, how many conservation constraints are there?", K, X, Y)
		return edakit.Numeric(q, edakit.IntStr(K*X*Y),
			edakit.Sprintf("%d vertices times %d commodities gives %d equations. The capacity constraints, one per edge, couple the commodities; without them the problem splits into independent shortest paths.", X*Y, K, K*X*Y),
			map[string]string{"vertices": edakit.IntStr(X * Y)})
	default:
		n := edakit.Int(3, 5)
		var parts []string
		worst := 0.0
		for i := 0; i < n; i++ {
			c := edakit.Int(2, 8)
			f := 0.5 * float64(edakit.Int(1, 4*c))
			u := f / float64(c)
			worst = math.Max(worst, u)
			parts = append(parts, edakit.Sprintf("e%d carries %s units with capacity %d", i+1, edakit.Num(f, 1), c))
		}
		ans := edakit.Num(worst, 2)
		q := "A fractional multicommodity-flow solution of a global-routing relaxation sums the flow of all commodities on each edge: " + strings.Join(parts, "; ") + ". The min-max formulation minimises the largest edge utilisation, flow divided by capacity. What is that utilisation for this solution? Round to 2 decimal places."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Divide each edge's flow by its capacity and take the maximum: %s. A value above 1 means no routing fits the capacities even fractionally under this flow.", ans),
			map[string]string{"lambda": ans})
	}
}

// pathfinderCost returns (b + h) * (1 + pf * max(0, occ + 1 - cap)).
func pathfinderCost(b, h, pf float64, occ, cp int) float64 {
	return (b + h) * (1 + pf*math.Max(0, float64(occ+1-cp)))
}

// negotiation: cost of one resource (KP 1), the cheaper of two paths (KP 2), cost after the
// present factor grows (KP 3).
func negotiation(ctx generator.GeneratorContext) generator.Problem {
	formula := "A negotiated-congestion router (PathFinder) prices a resource n for the net being routed as $c_n = (b_n + h_n)\\,p_n$, with base cost $b_n$, history cost $h_n$ and present-congestion factor $p_n = 1 + p_{fac}\\max(0,\\ occ_n + 1 - cap_n)$, where $occ_n$ counts the other nets already using n. "
	switch kp(ctx) {
	case 1:
		b := float64(edakit.Int(1, 4))
		h := 0.5 * float64(edakit.Int(0, 8))
		pf := edakit.Pick([]float64{0.5, 1, 1.5, 2, 3})
		cp := edakit.Int(1, 4)
		occ := cp + edakit.Int(-1, 2)
		c := pathfinderCost(b, h, pf, occ, cp)
		ans := edakit.Num(c, 2)
		q := formula + edakit.Sprintf("An edge has $b_n = %s$, $h_n = %s$, capacity %d and occupancy %d, and $p_{fac} = %s$. What is its cost for the next net? Round to 2 decimal places.", edakit.Num(b, 2), edakit.Num(h, 2), cp, occ, edakit.Num(pf, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$p_n = 1 + %s \\cdot \\max(0, %d + 1 - %d) = %s$, so $c_n = (%s + %s)(%s) = %s$.", edakit.Num(pf, 2), occ, cp, edakit.Num(1+pf*math.Max(0, float64(occ+1-cp)), 2), edakit.Num(b, 2), edakit.Num(h, 2), edakit.Num(1+pf*math.Max(0, float64(occ+1-cp)), 2), ans),
			map[string]string{"cost": ans})
	case 2:
		pf := edakit.Pick([]float64{0.5, 1, 2})
		var parts []string
		var costs []float64
		for _, name := range []string{"A", "B"} {
			n := edakit.Int(2, 4)
			if name == "B" {
				n += edakit.Int(1, 2)
			}
			tot := 0.0
			var es []string
			for i := 0; i < n; i++ {
				b := 1.0
				h := 0.5 * float64(edakit.Int(0, 4))
				cp := edakit.Int(1, 3)
				occ := edakit.Int(0, cp+1)
				if name == "B" {
					occ = edakit.Int(0, cp-1)
				}
				tot += pathfinderCost(b, h, pf, occ, cp)
				es = append(es, edakit.Sprintf("(h %s, cap %d, occ %d)", edakit.Num(h, 1), cp, occ))
			}
			costs = append(costs, tot)
			parts = append(parts, edakit.Sprintf("path %s uses %d edges %s", name, n, strings.Join(es, " ")))
		}
		best := math.Min(costs[0], costs[1])
		ans := edakit.Num(best, 2)
		q := formula + edakit.Sprintf("Every edge has $b_n = 1$ and $p_{fac} = %s$. Two candidate paths remain: %s. What is the cost of the path the router picks? Round to 2 decimal places.", edakit.Num(pf, 2), strings.Join(parts, "; "))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Path A costs %s and path B costs %s. The router takes the cheaper, %s. A longer path through free edges can beat a short one through full edges.", edakit.Num(costs[0], 2), edakit.Num(costs[1], 2), ans),
			map[string]string{"a": edakit.Num(costs[0], 2), "b": edakit.Num(costs[1], 2)})
	default:
		pf1 := edakit.Pick([]float64{0.5, 1})
		mult := edakit.Pick([]float64{1.5, 2})
		it := edakit.Int(2, 5)
		pf := pf1 * math.Pow(mult, float64(it-1))
		b := 1.0
		h := 0.5 * float64(edakit.Int(0, 6))
		cp := edakit.Int(1, 3)
		occ := cp + edakit.Int(0, 2)
		c := pathfinderCost(b, h, pf, occ, cp)
		ans := edakit.Num(c, 2)
		q := formula + edakit.Sprintf("The router starts with $p_{fac} = %s$ in iteration 1 and multiplies it by %s at every new iteration. In iteration %d an edge has $b_n = 1$, $h_n = %s$, capacity %d and occupancy %d. What is its cost? Round to 2 decimal places.", edakit.Num(pf1, 2), edakit.Num(mult, 2), it, edakit.Num(h, 2), cp, occ)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("In iteration %d, $p_{fac} = %s \\cdot %s^{%d} = %s$. Then $p_n = 1 + %s \\cdot %d$ and $c_n = %s$. Raising $p_{fac}$ each iteration makes sharing ever more expensive, so the nets that can move do.", it, edakit.Num(pf1, 2), edakit.Num(mult, 2), it-1, edakit.Num(pf, 3), edakit.Num(pf, 3), maxI(0, occ+1-cp), ans),
			map[string]string{"pfac": edakit.Num(pf, 3), "cost": ans})
	}
}

// historyCost: history after several iterations (KP 1), the iteration in which a net leaves a
// congested edge (KP 2), the full cost in a given iteration (KP 3).
func historyCost(ctx generator.GeneratorContext) generator.Problem {
	upd := "After each routing iteration, a negotiated-congestion router updates every edge's history cost by $h_n \\leftarrow h_n + h_{fac}\\max(0,\\ occ_n - cap_n)$, starting from $h_n = 0$. "
	switch kp(ctx) {
	case 1:
		hf := edakit.Pick([]float64{0.5, 1, 2})
		cp := edakit.Int(2, 6)
		n := edakit.Int(3, 6)
		occ := make([]int, n)
		h := 0.0
		for i := range occ {
			occ[i] = cp + edakit.Int(-2, 3)
			if occ[i] < 0 {
				occ[i] = 0
			}
			h += hf * math.Max(0, float64(occ[i]-cp))
		}
		ans := edakit.Num(h, 2)
		q := upd + edakit.Sprintf("An edge of capacity %d ends iterations 1 to %d with occupancies %s, and $h_{fac} = %s$. What is its history cost after iteration %d? Round to 2 decimal places.", cp, n, intsText(occ), edakit.Num(hf, 2), n)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Only iterations that end overflowed add history; each adds %s times its overflow. The sum is %s. History never decreases, so an edge that is chronically overused stays expensive even in an iteration when it happens to fit.", edakit.Num(hf, 2), ans),
			map[string]string{"history": ans})
	case 2:
		hf := edakit.Pick([]float64{0.5, 1, 1.5, 2})
		ov := edakit.Int(1, 3)
		bA := edakit.Int(5, 12)
		bB := bA + edakit.Int(2, 10)
		step := hf * float64(ov)
		i := int(math.Floor(float64(bB-bA)/step)) + 2
		q := upd + edakit.Sprintf("Ignore present-congestion cost. A net's shorter path A costs %d plus the history of one hot edge on it; its alternative path B costs %d and has no history. The hot edge ends every iteration overflowed by %d, and $h_{fac} = %s$. Iteration i uses the history accumulated through iteration i - 1. In which iteration is path A first strictly more expensive than path B?", bA, bB, ov, edakit.Num(hf, 2))
		return edakit.Numeric(q, edakit.IntStr(i),
			edakit.Sprintf("In iteration i the hot edge carries history %s(i - 1). A loses when %d + %s(i - 1) > %d, that is i - 1 > %s, so i = %d. History is what eventually moves a net that present cost alone would leave in place.", edakit.Num(step, 2), bA, edakit.Num(step, 2), bB, edakit.Num(float64(bB-bA)/step, 3), i),
			map[string]string{"iteration": edakit.IntStr(i)})
	default:
		hf := edakit.Pick([]float64{0.5, 1})
		pf := edakit.Pick([]float64{0.5, 1, 2})
		cp := edakit.Int(1, 4)
		n := edakit.Int(2, 4)
		occ := make([]int, n)
		h := 0.0
		for i := range occ {
			occ[i] = cp + edakit.Int(-1, 2)
			if occ[i] < 0 {
				occ[i] = 0
			}
			h += hf * math.Max(0, float64(occ[i]-cp))
		}
		cur := cp + edakit.Int(-1, 1)
		b := float64(edakit.Int(1, 3))
		c := pathfinderCost(b, h, pf, cur, cp)
		ans := edakit.Num(c, 2)
		q := upd + edakit.Sprintf("The cost of using edge n is $(b_n + h_n)\\,p_n$ with $p_n = 1 + p_{fac}\\max(0,\\ occ_n + 1 - cap_n)$. An edge with $b_n = %s$ and capacity %d ended iterations 1 to %d with occupancies %s, and $h_{fac} = %s$. In iteration %d, while routing a net, %d other nets already occupy it and $p_{fac} = %s$. What is its cost? Round to 2 decimal places.", edakit.Num(b, 0), cp, n, intsText(occ), edakit.Num(hf, 2), n+1, cur, edakit.Num(pf, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("History through iteration %d is %s. The present factor is $1 + %s \\cdot %d$. The cost is (%s + %s) times that, %s.", n, edakit.Num(h, 2), edakit.Num(pf, 2), maxI(0, cur+1-cp), edakit.Num(b, 0), edakit.Num(h, 2), ans),
			map[string]string{"history": edakit.Num(h, 2), "cost": ans})
	}
}

// ripUpReroute: nets touching overflowed edges (KP 1), the reroute order (KP 2), total overflow
// after one reroute (KP 3).
func ripUpReroute(ctx generator.GeneratorContext) generator.Problem {
	for try := 0; ; try++ {
		nE := 8
		nN := edakit.Scaled(ctx, 5, 7)
		routes := make([][]int, nN)
		use := make([]int, nE)
		for i := range routes {
			es := edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7})[:edakit.Int(2, 4)]
			sort.Ints(es)
			routes[i] = es
			for _, e := range es {
				use[e]++
			}
		}
		caps := make([]int, nE)
		for e := range caps {
			caps[e] = maxI(1, use[e]+edakit.Int(-1, 1))
		}
		if try >= generator.MaxDrawAttempts {
			// Known-good fallback: only e1 overflows, used by N1 and N2.
			nN = 5
			routes = [][]int{{0, 1}, {0, 2}, {3, 4}, {5, 6}, {1, 7}}
			use = []int{2, 2, 1, 1, 1, 1, 1, 1}
			caps = []int{1, 2, 1, 1, 1, 1, 1, 1}
		}
		over := map[int]bool{}
		for e := range caps {
			if use[e] > caps[e] {
				over[e] = true
			}
		}
		var ripped []int
		for i, r := range routes {
			for _, e := range r {
				if over[e] {
					ripped = append(ripped, i)
					break
				}
			}
		}
		if len(over) < 1 || len(ripped) < 2 || len(ripped) == nN {
			continue
		}
		var rt []string
		for i, r := range routes {
			var en []string
			for _, e := range r {
				en = append(en, edakit.Sprintf("e%d", e+1))
			}
			rt = append(rt, edakit.Sprintf("N%d uses %s", i+1, strings.Join(en, ", ")))
		}
		var cp []string
		for e, c := range caps {
			cp = append(cp, edakit.Sprintf("e%d: %d", e+1, c))
		}
		intro := edakit.Sprintf("After initial routing, the edge capacities are %s, and the nets' routes are: %s. ", strings.Join(cp, ", "), strings.Join(rt, "; "))
		switch kp(ctx) {
		case 1:
			return edakit.Numeric(intro+"The router rips up every net that uses at least one overflowed edge. How many nets does it rip up?", edakit.IntStr(len(ripped)),
				edakit.Sprintf("Count each edge's demand and compare with its capacity to find the overflowed edges; %d nets touch at least one of them.", len(ripped)),
				map[string]string{"ripped": edakit.IntStr(len(ripped))})
		case 2:
			var pool []int
			for v := 3; v <= 30; v++ {
				pool = append(pool, v)
			}
			hp := edakit.Shuffle(pool)[:nN]
			var hs []string
			for i := range hp {
				hs = append(hs, edakit.Sprintf("N%d %d", i+1, hp[i]))
			}
			order := append([]int(nil), ripped...)
			small := edakit.Int(0, 1) == 0
			sort.Slice(order, func(a, b int) bool {
				if small {
					return hp[order[a]] < hp[order[b]]
				}
				return hp[order[a]] > hp[order[b]]
			})
			pos := edakit.Int(1, len(order))
			rule, why := "smallest half-perimeter first, so short nets with few alternatives claim resources early", "smallest"
			if !small {
				rule, why = "largest half-perimeter first, so long nets that block many others are settled early", "largest"
			}
			ord := []string{"first", "second", "third", "fourth", "fifth", "sixth", "seventh"}
			q := intro + edakit.Sprintf("The net bounding-box half-perimeters are %s. The router rips up every net that uses an overflowed edge and reroutes them %s. Which net is rerouted %s? Answer with its number.", strings.Join(hs, ", "), rule, ord[pos-1])
			return edakit.Numeric(q, edakit.IntStr(order[pos-1]+1),
				edakit.Sprintf("The ripped nets are those touching an overflowed edge. Sorted by half-perimeter, %s first, the %s is N%d.", why, ord[pos-1], order[pos-1]+1),
				map[string]string{"ripped": edakit.IntStr(len(ripped))})
		default:
			tof := 0
			for e := range caps {
				tof += maxI(0, use[e]-caps[e])
			}
			r := edakit.Pick(ripped)
			alt := edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7})[:edakit.Int(2, 4)]
			sort.Ints(alt)
			use2 := append([]int(nil), use...)
			for _, e := range routes[r] {
				use2[e]--
			}
			for _, e := range alt {
				use2[e]++
			}
			t2 := 0
			for e := range caps {
				t2 += maxI(0, use2[e]-caps[e])
			}
			var en []string
			for _, e := range alt {
				en = append(en, edakit.Sprintf("e%d", e+1))
			}
			q := intro + edakit.Sprintf("Net N%d is ripped up and rerouted on %s. What is the total overflow afterward?", r+1, strings.Join(en, ", "))
			return edakit.Numeric(q, edakit.IntStr(t2),
				edakit.Sprintf("Before the reroute the total overflow is %d. Removing N%d's old edges and adding its new ones and summing max(0, demand - capacity) gives %d.", tof, r+1, t2),
				map[string]string{"before": edakit.IntStr(tof), "after": edakit.IntStr(t2)})
		}
	}
}
