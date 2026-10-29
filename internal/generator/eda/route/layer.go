package route

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const stackText = "Routing layers alternate direction: M2, M4 and M6 are horizontal and M3, M5 and M7 are vertical. Pins sit on M1, and a via stack between layers a and b counts |a - b| vias."

func layerChoices(horizontal bool) []int {
	if horizontal {
		return []int{2, 4, 6}
	}
	return []int{3, 5, 7}
}

func segDir(h bool) string {
	if h {
		return "horizontal"
	}
	return "vertical"
}

// pathVias counts vias of a two-pin path whose segments sit on layers ls.
func pathVias(ls []int) int {
	v := ls[0] - 1 + ls[len(ls)-1] - 1
	for i := 1; i < len(ls); i++ {
		v += absI(ls[i] - ls[i-1])
	}
	return v
}

// viaMinimization: vias of an assignment (KP 1), vias saved by the best assignment (KP 2), vias
// of a tree with a junction (KP 3).
func viaMinimization(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		n := edakit.Scaled(ctx, 3, 6)
		h := edakit.Int(0, 1) == 0
		ls := make([]int, n)
		lo := make([]int, n)
		var parts []string
		for i := range ls {
			hi := (i%2 == 0) == h
			ls[i] = edakit.Pick(layerChoices(hi))
			lo[i] = layerChoices(hi)[0]
			parts = append(parts, edakit.Sprintf("S%d %s on M%d", i+1, segDir(hi), ls[i]))
		}
		v := pathVias(ls)
		q := stackText + edakit.Sprintf(" A two-pin net's route runs from a pin, through segments %s in that order, to a second pin. ", strings.Join(parts, ", "))
		if ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			return edakit.Numeric(q+"How many vias does the route use, counting the via stacks down to both pins?", edakit.IntStr(v),
				edakit.Sprintf("Pin vias: %d - 1 and %d - 1. Between consecutive segments, add the layer differences. The total is %d.", ls[0], ls[n-1], v),
				map[string]string{"vias": edakit.IntStr(v)})
		}
		best := pathVias(lo)
		return edakit.Numeric(q+"Keeping the same 2D route and each segment's direction, a layer assigner may move each segment to any layer of its direction. How many vias does the best assignment save compared with this one?", edakit.IntStr(v-best),
			edakit.Sprintf("This assignment uses %d vias. Putting every horizontal segment on M2 and every vertical segment on M3 uses %d, the minimum, since every bend needs at least one via and each pin at least its own. The saving is %d. Capacity, not via count, is what pushes real assigners upward.", v, best, v-best),
			map[string]string{"given": edakit.IntStr(v), "best": edakit.IntStr(best)})
	default:
		l1 := edakit.Pick(layerChoices(true))
		l2 := edakit.Pick(layerChoices(true))
		l3 := edakit.Pick(layerChoices(false))
		l4 := edakit.Pick(layerChoices(false))
		jv := maxI(maxI(l1, l2), l4) - minI(minI(l1, l2), l4)
		tot := (l1 - 1) + jv + absI(l2-l3) + (l3 - 1) + (l4 - 1)
		q := stackText + edakit.Sprintf(" A three-pin net is routed as a tree. From pin P1, horizontal segment S1 on M%d runs to junction J. From J, horizontal segment S2 on M%d continues to a bend, where vertical segment S3 on M%d runs to pin P2. Also from J, vertical segment S4 on M%d runs to pin P3. At a point where several segments meet, the via stack spans from the lowest to the highest of their layers. How many vias does the tree use?", l1, l2, l3, l4)
		return edakit.Numeric(q, edakit.IntStr(tot),
			edakit.Sprintf("P1: %d. J joins M%d, M%d and M%d, so its stack spans %d. The bend: |%d - %d| = %d. P2: %d. P3: %d. Total %d.", l1-1, l1, l2, l4, jv, l2, l3, absI(l2-l3), l3-1, l4-1, tot),
			map[string]string{"junction": edakit.IntStr(jv), "vias": edakit.IntStr(tot)})
	}
}

// assignmentDP: min vias with forbidden layers (KP 1), one DP entry (KP 2), the number of
// optimal assignments (KP 3).
func assignmentDP(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	h := edakit.Int(0, 1) == 0
	allowed := make([][]int, n)
	var parts []string
	for i := range allowed {
		hi := (i%2 == 0) == h
		ch := layerChoices(hi)
		var a []int
		for _, l := range ch {
			if edakit.Int(0, 2) > 0 {
				a = append(a, l)
			}
		}
		if len(a) == 0 {
			a = []int{edakit.Pick(ch)}
		}
		allowed[i] = a
		var ln []string
		for _, l := range a {
			ln = append(ln, edakit.Sprintf("M%d", l))
		}
		parts = append(parts, edakit.Sprintf("S%d %s, free on %s", i+1, segDir(hi), strings.Join(ln, " or ")))
	}
	// dp[i][l] = min vias for segments 0..i with segment i on layer l, including the first pin.
	dp := make([]map[int]int, n)
	ways := make([]map[int]int, n)
	for i := range dp {
		dp[i], ways[i] = map[int]int{}, map[int]int{}
		for _, l := range allowed[i] {
			if i == 0 {
				dp[i][l], ways[i][l] = l-1, 1
				continue
			}
			best := math.MaxInt32
			w := 0
			for _, p := range allowed[i-1] {
				c := dp[i-1][p] + absI(l-p)
				if c < best {
					best, w = c, ways[i-1][p]
				} else if c == best {
					w += ways[i-1][p]
				}
			}
			dp[i][l], ways[i][l] = best, w
		}
	}
	best, w := math.MaxInt32, 0
	for _, l := range allowed[n-1] {
		c := dp[n-1][l] + l - 1
		if c < best {
			best, w = c, ways[n-1][l]
		} else if c == best {
			w += ways[n-1][l]
		}
	}
	q := stackText + edakit.Sprintf(" A two-pin route has segments, in order from pin to pin: %s (the other layers are full there). ", strings.Join(parts, "; "))
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(q+"A dynamic program over the segments picks one free layer per segment to minimise vias, including both pin stacks. What is the minimum via count?", edakit.IntStr(best),
			edakit.Sprintf("Let D(i, l) be the fewest vias up to segment i on layer l: D(1, l) = l - 1 and D(i, l) = min over free p of D(i - 1, p) + |l - p|. Closing with the last pin's stack, the minimum is %d.", best),
			map[string]string{"min_vias": edakit.IntStr(best)})
	case 2:
		i := edakit.Int(1, n-1)
		l := edakit.Pick(allowed[i])
		return edakit.Numeric(q+edakit.Sprintf("Let D(i, l) be the fewest vias for segments S1 to Si with Si on layer l, counting the first pin's stack but not the second's: D(1, l) = l - 1 and D(i, l) = min over free layers p of S(i-1) of D(i - 1, p) + |l - p|. What is D(%d, %d)?", i+1, l), edakit.IntStr(dp[i][l]),
			edakit.Sprintf("Filling the table segment by segment gives D(%d, %d) = %d. The table has one entry per segment and free layer, so the program is linear in the route length.", i+1, l, dp[i][l]),
			map[string]string{"entry": edakit.IntStr(dp[i][l])})
	default:
		return edakit.Numeric(q+edakit.Sprintf("The minimum via count, including both pin stacks, is %d. How many distinct layer assignments achieve it?", best), edakit.IntStr(w),
			edakit.Sprintf("Carry a count of optimal predecessors alongside each table entry, adding counts on ties. %d assignments reach %d vias. A real assigner breaks such ties by capacity or timing.", w, best),
			map[string]string{"min_vias": edakit.IntStr(best), "ways": edakit.IntStr(w)})
	}
}

// timingAware: Elmore delay on a lower layer (KP 1), on an upper layer reached by via stacks
// (KP 2), the break-even length (KP 3).
func timingAware(ctx generator.GeneratorContext) generator.Problem {
	rd := float64(50 * edakit.Int(4, 20))
	rlo := edakit.Pick([]float64{2, 3, 4, 5})
	rhi := edakit.Pick([]float64{0.2, 0.3, 0.4, 0.5})
	c := edakit.Pick([]float64{0.15, 0.2, 0.25})
	cl := float64(edakit.Int(2, 12))
	rv := float64(5 * edakit.Int(2, 8))
	k := edakit.Int(2, 4)
	lower := func(L float64) float64 { return rd*(c*L+cl) + rlo*L*(c*L/2+cl) }
	upper := func(L float64) float64 {
		ru := float64(k) * rv
		return rd*(c*L+cl) + ru*(c*L+cl) + rhi*L*(c*L/2+cl) + ru*cl
	}
	model := edakit.Sprintf("Delay is Elmore: a driver of %s ohm, a distributed wire of capacitance %s fF/um, and a %s fF load. On the low layer the wire has %s ohm/um. On the upper layer it has %s ohm/um, but reaching it takes a stack of %d vias of %s ohm each at both ends, modelled as lumped resistances at the start and end of the wire. One ohm times one fF is one fs. ",
		edakit.Num(rd, 0), edakit.Num(c, 2), edakit.Num(cl, 0), edakit.Num(rlo, 1), edakit.Num(rhi, 1), k, edakit.Num(rv, 0))
	switch kp(ctx) {
	case 1, 2:
		L := float64(50 * edakit.Int(2, 40))
		if ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			v := lower(L) / 1000
			ans := edakit.Num(v, 2)
			return edakit.Numeric(model+edakit.Sprintf("What is the delay of a %s um route on the low layer, in ps? Round to 2 decimal places.", edakit.Num(L, 0)), ans,
				edakit.Sprintf("$R_d(cL + C_L) + rL(cL/2 + C_L)$ = %s ohm-fF = %s ps.", edakit.Num(lower(L), 1), ans),
				map[string]string{"lower_ps": ans})
		}
		v := upper(L) / 1000
		ans := edakit.Num(v, 2)
		return edakit.Numeric(model+edakit.Sprintf("What is the delay of a %s um route on the upper layer, in ps? Round to 2 decimal places.", edakit.Num(L, 0)), ans,
			edakit.Sprintf("The start stack adds its resistance to the driver's, charging all of $cL + C_L$; the wire term uses the upper-layer resistance; the end stack charges only the load. Total %s ohm-fF = %s ps, against %s ps on the low layer.", edakit.Num(upper(L), 1), ans, edakit.Num(lower(L)/1000, 2)),
			map[string]string{"upper_ps": ans, "lower_ps": edakit.Num(lower(L)/1000, 2)})
	default:
		L := 1
		for ; L < 20000; L++ {
			if upper(float64(L)) < lower(float64(L)) {
				break
			}
		}
		return edakit.Numeric(model+"What is the shortest route length, in whole um, for which the upper layer is strictly faster?", edakit.IntStr(L),
			edakit.Sprintf("The upper layer saves $(r_{lo} - r_{hi})L(cL/2 + C_L)$ and pays $kR_v(cL + 2C_L)$ for the stacks. The saving grows with $L^2$ and the penalty with L, so they cross once; at L = %d um the upper layer first wins. Short nets belong on low layers, long critical nets high.", L),
			map[string]string{"break_even_um": edakit.IntStr(L)})
	}
}

// redundantVias: expected failing vias (KP 1), with a share doubled (KP 2), the via-limited
// yield (KP 3).
func redundantVias(ctx generator.GeneratorContext) generator.Problem {
	Nm := float64(edakit.Int(2, 40))
	pe := edakit.Int(1, 9)
	p := float64(pe) * 1e-7
	N := Nm * 1e6
	r := float64(edakit.Int(50, 95)) / 100
	ptxt := edakit.Sprintf("$%d \\times 10^{-7}$", pe)
	switch kp(ctx) {
	case 1:
		e := N * p
		ans := edakit.Num(e, 2)
		q := edakit.Sprintf("A routed design has %s million single-cut vias. Each fails (opens) independently with probability %s. What is the expected number of failing vias? Round to 2 decimal places.", edakit.Num(Nm, 0), ptxt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Expectation is linear: %s million x %s = %s.", edakit.Num(Nm, 0), ptxt, ans),
			map[string]string{"expected": ans})
	case 2:
		e := N * ((1-r)*p + r*p*p)
		ans := edakit.Num(e, 3)
		q := edakit.Sprintf("A design has %s million via locations. Each cut fails independently with probability %s. Redundant-via insertion doubles %d%% of the vias; a doubled via opens only if both of its cuts fail. What is the expected number of open vias? Round to 3 decimal places.", edakit.Num(Nm, 0), ptxt, int(math.Round(r*100)))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Single vias contribute %s x %s; doubled ones contribute %s x $p^2$, which is negligible. The total is %s, against %s with no redundancy.", edakit.Num(N*(1-r), 0), ptxt, edakit.Num(N*r, 0), ans, edakit.Num(N*p, 2)),
			map[string]string{"expected": ans})
	default:
		y := math.Pow(1-p, N*(1-r)) * math.Pow(1-p*p, N*r) * 100
		ans := edakit.Num(y, 2)
		q := edakit.Sprintf("A design has %s million via locations, each cut failing independently with probability %s. After redundant-via insertion, %d%% of the vias are doubled (open only if both cuts fail). What percentage of dies have no open via? Round to 2 decimal places.", edakit.Num(Nm, 0), ptxt, int(math.Round(r*100)))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Yield is $(1 - p)^{N_s}(1 - p^2)^{N_d}$ with %s single and %s doubled vias, which is about $e^{-N_s p}$ = %s%%. With no doubling it would be %s%%.", edakit.Num(N*(1-r), 0), edakit.Num(N*r, 0), ans, edakit.Num(math.Pow(1-p, N)*100, 2)),
			map[string]string{"yield_pct": ans})
	}
}
