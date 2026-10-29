package cts

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const pathLenModel = "Use the linear (path-length) delay model, in which the delay from a node to a sink equals the wire length between them in µm."

// linearMerge: x = (L + t2 - t1) / 2 under path-length delay, and the detour when |t1 - t2| > L.
func linearMerge(ctx generator.GeneratorContext) generator.Problem {
	k := kpOf(ctx)
	L := 10 * edakit.Int(10, 60)
	if k == 3 {
		t1 := 10 * edakit.Int(1, 10)
		t2 := t1 + L + 5*edakit.Int(1, 20)
		if edakit.Int(0, 1) == 1 {
			t1, t2 = t2, t1
		}
		extra := absInt(t1-t2) - L
		q := pathLenModel + edakit.Sprintf(" Subtree 1 has delay %d µm from its root to every sink and subtree 2 has %d µm. Their roots are %d µm apart. A zero-skew merge is impossible on the straight connection, so the tap sits at the root of the slower subtree and the wire to the faster one is lengthened. How much extra wire (beyond the %d µm straight distance) must be snaked in, in µm?", t1, t2, L, L)
		return edakit.Numeric(q, edakit.IntStr(extra),
			edakit.Sprintf("Zero skew needs the wire to the faster subtree to equal the delay difference, $|t_1 - t_2| = %d$ µm. The straight distance is %d µm, so %d µm of detour is needed. The formula $x = (L + t_2 - t_1)/2$ would give a value outside $[0, L]$, which is the signal that elongation is needed.", absInt(t1-t2), L, extra),
			map[string]string{"detour": edakit.IntStr(extra)})
	}
	t1 := 10 * edakit.Int(1, 15)
	t2 := t1 + 10*edakit.Int(-L/20, L/20)
	if t2 < 0 {
		t2 = -t2
	}
	x := float64(L+t2-t1) / 2
	if x < 0 || x > float64(L) {
		t2 = t1
		x = float64(L) / 2
	}
	base := pathLenModel + edakit.Sprintf(" Subtree 1 has delay %d µm from its root to every sink and subtree 2 has %d µm. Their roots are joined by a wire of %d µm. ", t1, t2, L)
	if k == 1 {
		return edakit.Numeric(base+"How far from the root of subtree 1 must the tapping point be for zero skew, in µm? Round to 1 decimal place.", f(x, 1),
			edakit.Sprintf("Balance $x + t_1 = (L - x) + t_2$: $x = (L + t_2 - t_1)/2 = (%d + %d - %d)/2 = %s$ µm.", L, t2, t1, f(x, 1)),
			map[string]string{"x": f(x, 1)})
	}
	md := x + float64(t1)
	return edakit.Numeric(base+"After the zero-skew merge, what is the delay from the tapping point to every sink, in µm? Round to 1 decimal place.", f(md, 1),
		edakit.Sprintf("The tap is at $x = (L + t_2 - t_1)/2 = %s$ µm from subtree 1, so the merged delay is $x + t_1 = %s$ µm, equal to $(L - x) + t_2$. In general it is $(L + t_1 + t_2)/2$.", f(x, 1), f(md, 1)),
		map[string]string{"merged_delay": f(md, 1)})
}

type tsay struct {
	r, c, L, c1, c2, t1, t2 float64
}

func (z tsay) describe() string {
	return edakit.Sprintf("Two zero-skew subtrees are merged by a wire of length $L = %s$ µm with $r = %s\\ \\Omega$/µm and $c = %s$ fF/µm. Subtree 1 has delay $t_1 = %s$ ps and load $C_1 = %s$ fF; subtree 2 has $t_2 = %s$ ps and $C_2 = %s$ fF. Use Elmore delay, with the wire on each side of the tap modelled as a $\\pi$ segment.",
		f(z.L, 0), f(z.r, 2), f(z.c, 2), f(z.t1, 1), f(z.c1, 0), f(z.t2, 1), f(z.c2, 0))
}

// tsayMerge: KP1 the tapping fraction x, KP2 the merged delay, KP3 elongation when x > 1.
func tsayMerge(ctx generator.GeneratorContext) generator.Problem {
	z := tsay{r: pick(0.05, 0.3, 0.01), c: pick(0.1, 0.25, 0.01), L: float64(50 * edakit.Int(4, 16)),
		c1: float64(5 * edakit.Int(4, 30)), c2: float64(5 * edakit.Int(4, 30))}
	k := kpOf(ctx)
	if k == 3 {
		z.L = float64(10 * edakit.Int(10, 30))
		thr := z.r * z.L * (z.c*z.L/2 + z.c1) / 1000
		z.t1 = pick(5, 20, 0.5)
		z.t2 = edakit.Round(z.t1+thr+pick(0.5, 4, 0.5), 1)
		// (rc/2) L'^2 + r C1 L' - 1000 (t2 - t1) = 0
		a, b, cc := z.r*z.c/2, z.r*z.c1, -1000*(z.t2-z.t1)
		lp := (-b + math.Sqrt(b*b-4*a*cc)) / (2 * a)
		q := z.describe() + " Here Tsay's formula gives a tapping fraction above 1, so the tap moves onto the root of subtree 2 and the wire to subtree 1 is elongated. What wire length to subtree 1 gives zero skew, in µm? Round to 1 decimal place."
		return edakit.Numeric(q, f(lp, 1),
			edakit.Sprintf("Zero skew needs $r L' (c L'/2 + C_1)/1000 = t_2 - t_1 = %s$ ps, a quadratic $(rc/2)L'^2 + rC_1 L' - 1000(t_2 - t_1) = 0$. Its positive root is $L' = %s$ µm, longer than the %s µm straight distance.", f(z.t2-z.t1, 1), f(lp, 1), f(z.L, 0)),
			map[string]string{"elongated": f(lp, 2)})
	}
	var x float64
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		z.t1, z.t2 = pick(5, 40, 0.5), pick(5, 40, 0.5)
		x = (z.t2 - z.t1 + z.r*z.L*(z.c2+z.c*z.L/2)/1000) / (z.r * z.L * (z.c*z.L + z.c1 + z.c2) / 1000)
		if x > 0.05 && x < 0.95 {
			break
		}
	}
	if x <= 0 || x >= 1 { // known-good instance
		z = tsay{r: 0.1, c: 0.14, L: 200, c1: 75, c2: 25, t1: 35, t2: 36}
		x = (z.t2 - z.t1 + z.r*z.L*(z.c2+z.c*z.L/2)/1000) / (z.r * z.L * (z.c*z.L + z.c1 + z.c2) / 1000)
	}
	num := z.t2 - z.t1 + z.r*z.L*(z.c2+z.c*z.L/2)/1000
	den := z.r * z.L * (z.c*z.L + z.c1 + z.c2) / 1000
	if k == 1 {
		return edakit.Numeric(z.describe()+" At what fraction $x$ of the wire, measured from subtree 1, must the tap be placed for zero skew? Round to 3 decimal places.", f(x, 3),
			edakit.Sprintf("Tsay: $x = \\frac{(t_2 - t_1) + rL(C_2 + cL/2)}{rL(cL + C_1 + C_2)}$. With delays in ps (ohm-fF divided by 1000): numerator %s, denominator %s, so $x = %s$.", f(num, 3), f(den, 3), f(x, 3)),
			map[string]string{"x": f(x, 4)})
	}
	xl := x * z.L
	md := z.t1 + z.r*xl*(z.c*xl/2+z.c1)/1000
	return edakit.Numeric(z.describe()+" After merging at Tsay's zero-skew tapping point, what is the delay from the tap to every sink, in ps? Round to 2 decimal places.", f(md, 2),
		edakit.Sprintf("Tsay gives $x = %s$, so the tap is %s µm from subtree 1. The delay is $t_1 + r\\,xL\\,(c\\,xL/2 + C_1)/1000 = %s + %s = %s$ ps, and the same value comes out of the subtree-2 side.", f(x, 4), f(xl, 1), f(z.t1, 1), f(md-z.t1, 3), f(md, 2)),
		map[string]string{"x": f(x, 4), "merged_delay": f(md, 3)})
}

type pt struct{ x, y float64 }

func l1(a, b pt) float64 { return absF(a.x-b.x) + absF(a.y-b.y) }

// segment returns the endpoints of the merging segment of sinks a and b when the tap must be
// at distance ea from a and d - ea from b (all points lie in the bounding box of a and b).
func segment(a, b pt, ea float64) (pt, pt) {
	dx, dy := absF(b.x-a.x), absF(b.y-a.y)
	sx, sy := 1.0, 1.0
	if b.x < a.x {
		sx = -1
	}
	if b.y < a.y {
		sy = -1
	}
	ulo := math.Max(0, ea-dy)
	uhi := math.Min(dx, ea)
	p := func(u float64) pt { return pt{a.x + sx*u, a.y + sy*(ea-u)} }
	return p(ulo), p(uhi)
}

// distToSeg is the Manhattan distance from p to the straight segment e1-e2, found by checking
// the endpoints and the breakpoints of the convex piecewise-linear distance.
func distToSeg(p, e1, e2 pt) float64 {
	best := math.Min(l1(p, e1), l1(p, e2))
	dx, dy := e2.x-e1.x, e2.y-e1.y
	var ts []float64
	if dx != 0 {
		ts = append(ts, (p.x-e1.x)/dx)
	}
	if dy != 0 {
		ts = append(ts, (p.y-e1.y)/dy)
	}
	for _, t := range ts {
		if t > 0 && t < 1 {
			q := pt{e1.x + t*dx, e1.y + t*dy}
			best = math.Min(best, l1(p, q))
		}
	}
	return best
}

func ptStr(p pt) string { return "(" + f(p.x, 1) + ", " + f(p.y, 1) + ")" }

type msDraw struct {
	a, b       pt
	ta, tb, ea float64
	e1, e2     pt
}

func drawMS() msDraw {
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		a := pt{float64(10 * edakit.Int(0, 6)), float64(10 * edakit.Int(0, 6))}
		b := pt{float64(10 * edakit.Int(0, 6)), float64(10 * edakit.Int(0, 6))}
		dx, dy := absF(a.x-b.x), absF(a.y-b.y)
		if dx < 10 || dy < 10 {
			continue
		}
		d := dx + dy
		ta, tb := float64(5*edakit.Int(0, 8)), float64(5*edakit.Int(0, 8))
		ea := (d + tb - ta) / 2
		if ea <= 0 || ea >= d {
			continue
		}
		e1, e2 := segment(a, b, ea)
		if l1(e1, e2) < 1 {
			continue
		}
		return msDraw{a, b, ta, tb, ea, e1, e2}
	}
	a, b := pt{50, 30}, pt{20, 20}
	e1, e2 := segment(a, b, 22.5)
	return msDraw{a, b, 15, 20, 22.5, e1, e2}
}

// parentPoint draws a placed parent that is not on the segment e1-e2.
func parentPoint(e1, e2 pt) pt {
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		p := pt{float64(10 * edakit.Int(-2, 8)), float64(10 * edakit.Int(-2, 8))}
		if distToSeg(p, e1, e2) > 0 {
			return p
		}
	}
	return pt{math.Min(e1.x, e2.x) - 20, math.Max(e1.y, e2.y) + 20}
}

func (m msDraw) describe() string {
	return pathLenModel + edakit.Sprintf(" Sink a sits at %s with downstream delay %s µm and sink b at %s with downstream delay %s µm (coordinates in µm).", ptStr(m.a), f(m.ta, 0), ptStr(m.b), f(m.tb, 0))
}

// mergingSegment: KP1 the edge length to a, KP2 the x-extent of the merging segment, KP3 the
// Manhattan distance from a parent point to the segment.
func mergingSegment(ctx generator.GeneratorContext) generator.Problem {
	m := drawMS()
	d := l1(m.a, m.b)
	switch kpOf(ctx) {
	case 1:
		return edakit.Numeric(m.describe()+" For a zero-skew merge with no detour, how long is the edge from the merge point to sink a, in µm? Round to 1 decimal place.", f(m.ea, 1),
			edakit.Sprintf("The Manhattan distance is $d = %s$ µm. Zero skew needs $e_a + t_a = e_b + t_b$ with $e_a + e_b = d$, so $e_a = (d + t_b - t_a)/2 = %s$ µm.", f(d, 0), f(m.ea, 1)),
			map[string]string{"ea": f(m.ea, 1)})
	case 2:
		ext := absF(m.e2.x - m.e1.x)
		return edakit.Numeric(m.describe()+" Every point at Manhattan distance $e_a$ from a and $d - e_a$ from b is a valid zero-skew merge point; together they form the merging segment, a line of slope $\\pm 1$. How long is its projection on the x axis, in µm? Round to 1 decimal place.", f(ext, 1),
			edakit.Sprintf("Here $e_a = (d + t_b - t_a)/2 = %s$ µm. The points lie inside the bounding box of a and b, on the diagonal at distance $e_a$ from a. Its ends are %s and %s, so the x extent is %s µm.", f(m.ea, 1), ptStr(m.e1), ptStr(m.e2), f(ext, 1)),
			map[string]string{"e1": ptStr(m.e1), "e2": ptStr(m.e2)})
	}
	p := parentPoint(m.e1, m.e2)
	dist := distToSeg(p, m.e1, m.e2)
	q := m.describe() + edakit.Sprintf(" Their merging segment runs from %s to %s. In the top-down pass the parent of this merge has been placed at %s. DME places the merge point at the point of the segment nearest the parent. What is the Manhattan length of the parent's edge to it, in µm? Round to 1 decimal place.", ptStr(m.e1), ptStr(m.e2), ptStr(p))
	return edakit.Numeric(q, f(dist, 1),
		edakit.Sprintf("Manhattan distance to a slope-1 segment is a convex piecewise-linear function along it, so check the endpoints and the points level with the parent in x or y. The minimum is %s µm.", f(dist, 1)),
		map[string]string{"dist": f(dist, 1)})
}

// topologyGeneration: greedy nearest-pair matching and the method of means and medians.
func topologyGeneration(ctx generator.GeneratorContext) generator.Problem {
	k := kpOf(ctx)
	n := edakit.Scaled(ctx, 4, 7)
	if k == 3 {
		n = edakit.Pick([]int{4, 6})
	}
	var ps []pt
	var dists map[[2]int]float64
	for attempt := 0; ; attempt++ {
		ps = nil
		for i := 0; i < n; i++ {
			ps = append(ps, pt{float64(edakit.Int(0, 60)), float64(edakit.Int(0, 60))})
		}
		if attempt >= generator.MaxDrawAttempts {
			// Known-good instance: all pair distances distinct.
			n = 4
			ps = []pt{{18, 22}, {49, 37}, {6, 41}, {15, 47}}
		}
		dists = map[[2]int]float64{}
		seen := map[float64]bool{}
		ok := true
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				dd := l1(ps[i], ps[j])
				if dd == 0 || (k == 3 && seen[dd]) {
					ok = false
				}
				seen[dd] = true
				dists[[2]int{i, j}] = dd
			}
		}
		if ok {
			break
		}
	}
	var desc []string
	for i, p := range ps {
		desc = append(desc, edakit.Sprintf("s%d %s", i+1, ptStr(p)))
	}
	sinks := "Clock sinks sit at (in µm) " + strings.Join(desc, ", ") + "."
	switch k {
	case 1:
		best, bi, bj := math.Inf(1), 0, 0
		for key, dd := range dists {
			if dd < best {
				best, bi, bj = dd, key[0], key[1]
			}
		}
		return edakit.Numeric(sinks+" A greedy bottom-up topology generator merges the closest pair of sinks first, by Manhattan distance. What is the distance between the first pair it merges, in µm?", f(best, 0),
			edakit.Sprintf("Over all pairs the smallest Manhattan distance is between s%d and s%d: %s µm.", bi+1, bj+1, f(best, 0)),
			map[string]string{"first_pair_dist": f(best, 0)})
	case 2:
		xs := make([]float64, n)
		for i, p := range ps {
			xs[i] = p.x
		}
		sorted := append([]float64(nil), xs...)
		for i := range sorted {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[j] < sorted[i] {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		var med float64
		if n%2 == 1 {
			med = sorted[n/2]
		} else {
			med = (sorted[n/2-1] + sorted[n/2]) / 2
		}
		return edakit.Numeric(sinks+" The method of means and medians (MMM) splits the set at the median x coordinate to build the top level of the tree. At what x value is the cut made, in µm? Round to 1 decimal place.", f(med, 1),
			edakit.Sprintf("The sorted x coordinates are %s. The median is %s µm, which puts half the sinks on each side and so balances the two subtrees.", floats(sorted, 0), f(med, 1)),
			map[string]string{"median": f(med, 1)})
	}
	used := make([]bool, n)
	total := 0.0
	var log []string
	for r := 0; r < n/2; r++ {
		best, bi, bj := math.Inf(1), -1, -1
		for key, dd := range dists {
			if !used[key[0]] && !used[key[1]] && dd < best {
				best, bi, bj = dd, key[0], key[1]
			}
		}
		used[bi], used[bj] = true, true
		total += best
		log = append(log, edakit.Sprintf("s%d-s%d (%s)", bi+1, bj+1, f(best, 0)))
	}
	return edakit.Numeric(sinks+" One round of greedy matching repeatedly pairs the two closest unmatched sinks (Manhattan distance) until every sink is paired. What is the total Manhattan length of the pairs chosen in this round, in µm?", f(total, 0),
		edakit.Sprintf("The greedy order is %s, totalling %s µm. Each pair becomes one merge node in the next level of the topology.", strings.Join(log, ", "), f(total, 0)),
		map[string]string{"round_total": f(total, 0)})
}

// boundedSkew: the feasible tapping interval when the merged skew may be up to B.
func boundedSkew(ctx generator.GeneratorContext) generator.Problem {
	var L, a1, b1, a2, b2, B, lo, hi float64
	for attempt := 0; ; attempt++ {
		L = float64(10 * edakit.Int(10, 40))
		B = float64(5 * edakit.Int(2, 8))
		a1 = float64(10 * edakit.Int(5, 20))
		b1 = a1 + float64(5*edakit.Int(0, int(B/5)-1))
		a2 = float64(10 * edakit.Int(5, 20))
		b2 = a2 + float64(5*edakit.Int(0, int(B/5)-1))
		if attempt >= generator.MaxDrawAttempts {
			L, B, a1, b1, a2, b2 = 330, 40, 80, 110, 70, 95
		}
		lo = math.Max(0, (L+b2-a1-B)/2)
		hi = math.Min(L, (B+L+a2-b1)/2)
		if hi-lo >= 2.5 || attempt >= generator.MaxDrawAttempts {
			break
		}
	}
	desc := pathLenModel + edakit.Sprintf(" Subtree 1's sink delays range from %s to %s µm and subtree 2's from %s to %s µm. Their roots are joined by a %s µm wire, and the tap sits $x$ µm from subtree 1. The skew bound is $B = %s$ µm.", f(a1, 0), f(b1, 0), f(a2, 0), f(b2, 0), f(L, 0), f(B, 0))
	skewAt := func(x float64) float64 {
		return math.Max(x+b1, L-x+b2) - math.Min(x+a1, L-x+a2)
	}
	switch kpOf(ctx) {
	case 1:
		return edakit.Numeric(desc+" What is the smallest $x$ that meets the skew bound, in µm? Round to 1 decimal place.", f(lo, 1),
			edakit.Sprintf("The binding constraint at small $x$ is the slow sink of subtree 2 against the fast sink of subtree 1: $(L - x + %s) - (x + %s) \\le B$, so $x \\ge (L + %s - %s - B)/2$. Clamped to $[0, L]$ this is %s µm.", f(b2, 0), f(a1, 0), f(b2, 0), f(a1, 0), f(lo, 1)),
			map[string]string{"lo": f(lo, 2), "hi": f(hi, 2)})
	case 2:
		return edakit.Numeric(desc+" How wide is the interval of tapping positions $x$ that meet the skew bound, in µm? Round to 1 decimal place.", f(hi-lo, 1),
			edakit.Sprintf("Lower end: $x \\ge (L + b_2 - a_1 - B)/2$, upper end: $x \\le (L + a_2 - b_1 + B)/2$, both clamped to $[0, L]$. That gives $[%s, %s]$, a width of %s µm. A zero-skew merge pins $x$ to one point; the bound turns it into an interval that a bounded-skew tree spends on shorter wire.", f(lo, 1), f(hi, 1), f(hi-lo, 1)),
			map[string]string{"lo": f(lo, 2), "hi": f(hi, 2)})
	}
	x := 10 * math.Round(pick(math.Max(0, lo-50), math.Min(L, hi+50), 1)/10)
	s := skewAt(x)
	return edakit.Numeric(desc+edakit.Sprintf(" If the tap is placed at $x = %s$ µm, what is the skew of the merged tree, in µm?", f(x, 0)), f(s, 1),
		edakit.Sprintf("From the tap, the latest sink is at $\\max(x + %s, L - x + %s) = %s$ and the earliest at $\\min(x + %s, L - x + %s) = %s$. The skew is %s µm, which %s the bound of %s.", f(b1, 0), f(b2, 0), f(math.Max(x+b1, L-x+b2), 0), f(a1, 0), f(a2, 0), f(math.Min(x+a1, L-x+a2), 0), f(s, 1), map[bool]string{true: "meets", false: "violates"}[s <= B], f(B, 0)),
		map[string]string{"skew": f(s, 1)})
}

// dmeBottomUp: computed choice for the bottom-up edge length of a merge.
func dmeBottomUp(ctx generator.GeneratorContext) generator.Problem {
	m := drawMS()
	d := l1(m.a, m.b)
	eu := math.Sqrt((m.a.x-m.b.x)*(m.a.x-m.b.x) + (m.a.y-m.b.y)*(m.a.y-m.b.y))
	c, ds := numOptions(m.ea, 1, " µm", d/2, d-m.ea, (eu+m.tb-m.ta)/2, (d+m.ta-m.tb)/2+5)
	q := m.describe() + " In the bottom-up phase of DME, how long is the edge from their merge point to sink a?"
	return edakit.Choice(q, c,
		edakit.Sprintf("DME works in the Manhattan metric: $d = %s$ µm and $e_a = (d + t_b - t_a)/2 = %s$ µm. Halving $d$ ignores the delay difference, and the Euclidean distance is the wrong metric for rectilinear wires.", f(d, 0), f(m.ea, 1)), ds...)
}

// dmeTopDown: computed choice for the top-down embedding distance.
func dmeTopDown(ctx generator.GeneratorContext) generator.Problem {
	m := drawMS()
	p := parentPoint(m.e1, m.e2)
	dist := distToSeg(p, m.e1, m.e2)
	mid := pt{(m.e1.x + m.e2.x) / 2, (m.e1.y + m.e2.y) / 2}
	far := math.Max(l1(p, m.e1), l1(p, m.e2))
	near := math.Min(math.Hypot(p.x-m.e1.x, p.y-m.e1.y), math.Hypot(p.x-m.e2.x, p.y-m.e2.y))
	c, ds := numOptions(dist, 1, " µm", l1(p, mid), far, near)
	q := edakit.Sprintf("In the top-down phase of DME, a node's merging segment runs from %s to %s (µm). Its parent has already been embedded at %s. DME embeds the node at the point of its segment nearest the parent. How long is the rectilinear edge from the parent to the node?", ptStr(m.e1), ptStr(m.e2), ptStr(p))
	return edakit.Choice(q, c,
		edakit.Sprintf("The nearest point minimises Manhattan distance along the segment; checking the endpoints and the points level with the parent gives %s µm. The midpoint and the far endpoint are not nearest in general, and Euclidean distance is the wrong metric.", f(dist, 1)), ds...)
}
