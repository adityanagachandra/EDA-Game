package fpga

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// crossCount is VPR's q(n) table for nets of 1 to 10 terminals.
var crossCount = []float64{1, 1, 1, 1, 1.0828, 1.1536, 1.2206, 1.2823, 1.3385, 1.3991, 1.4493}

type pt struct{ x, y int }

func ptsText(ps []pt) string {
	var s []string
	for _, p := range ps {
		s = append(s, edakit.Sprintf("(%d, %d)", p.x, p.y))
	}
	return strings.Join(s, ", ")
}

func spans(ps []pt) (int, int) {
	x0, x1, y0, y1 := ps[0].x, ps[0].x, ps[0].y, ps[0].y
	for _, p := range ps[1:] {
		x0, x1 = minInt(x0, p.x), maxInt(x1, p.x)
		y0, y1 = minInt(y0, p.y), maxInt(y1, p.y)
	}
	return x1 - x0 + 1, y1 - y0 + 1
}

func netCost(ps []pt) float64 {
	sx, sy := spans(ps)
	return crossCount[len(ps)] * float64(sx+sy)
}

func randPts(n, g int) []pt {
	seen := map[pt]bool{}
	var ps []pt
	for len(ps) < n {
		p := pt{edakit.Int(1, g), edakit.Int(1, g)}
		if !seen[p] {
			seen[p] = true
			ps = append(ps, p)
		}
	}
	return ps
}

func qText(ns ...int) string {
	seen := map[int]bool{}
	var s []string
	low := false
	for _, n := range ns {
		if n <= 3 {
			low = true
			continue
		}
		if !seen[n] {
			seen[n] = true
			s = append(s, edakit.Sprintf("q(%d) = %s", n, edakit.Num(crossCount[n], 4)))
		}
	}
	sortStrings(s)
	if low {
		s = append([]string{"q(n) = 1 for n ≤ 3"}, s...)
	}
	return strings.Join(s, ", ")
}

const bbRule = "VPR's bounding-box cost of a net with n terminals is q(n)·[(x_max - x_min + 1) + (y_max - y_min + 1)], counting the channels the bounding box spans; all channels have equal capacity."

func bboxCost(ctx generator.GeneratorContext) generator.Problem {
	g := edakit.Scaled(ctx, 6, 20)
	switch whichKP(ctx) {
	case 1:
		n := edakit.Int(2, 8)
		ps := randPts(n, g)
		sx, sy := spans(ps)
		c := netCost(ps)
		q := bbRule + edakit.Sprintf(" Use %s. A net has terminals at %s. What is its cost? Round to 2 decimal places.", qText(n), ptsText(ps))
		return num(q, edakit.Num(c, 2),
			edakit.Sprintf("The box spans %d columns and %d rows of channels, and the net has %d terminals, so the cost is %s × (%d + %d) = %s. The q(n) factor corrects half-perimeter for nets whose Steiner tree must cross the box more than once.", sx, sy, n, edakit.Num(crossCount[n], 4), sx, sy, edakit.Num(c, 2)),
			map[string]string{"x_span": edakit.IntStr(sx), "y_span": edakit.IntStr(sy)})
	case 2:
		var lines []string
		var ns []int
		tot := 0.0
		var each []string
		for i := 0; i < 3; i++ {
			n := edakit.Int(2, 7)
			ns = append(ns, n)
			ps := randPts(n, g)
			c := netCost(ps)
			tot += c
			lines = append(lines, edakit.Sprintf("net %d at %s", i+1, ptsText(ps)))
			each = append(each, edakit.Num(c, 4))
		}
		q := bbRule + edakit.Sprintf(" Use %s. The placement has three nets: %s. What is the total bounding-box cost? Round to 2 decimal places.", qText(ns...), strings.Join(lines, "; "))
		return num(q, edakit.Num(tot, 2),
			edakit.Sprintf("Net costs: %s. Sum: %s.", strings.Join(each, ", "), edakit.Num(tot, 2)),
			map[string]string{"total": edakit.Num(tot, 2)})
	default:
		from := pt{edakit.Int(1, g), edakit.Int(1, g)}
		to := pt{edakit.Int(1, g), edakit.Int(1, g)}
		for to == from {
			to = pt{edakit.Int(1, g), edakit.Int(1, g)}
		}
		var lines []string
		var ns []int
		d := 0.0
		var parts []string
		for i := 0; i < 2; i++ {
			n := edakit.Int(2, 5)
			ns = append(ns, n)
			others := randPts(n-1, g)
			for hasPt(others, from) || hasPt(others, to) {
				others = randPts(n-1, g)
			}
			before := append([]pt{from}, others...)
			after := append([]pt{to}, others...)
			cb, ca := netCost(before), netCost(after)
			d += ca - cb
			lines = append(lines, edakit.Sprintf("net %d connects B to %s", i+1, ptsText(others)))
			parts = append(parts, edakit.Sprintf("net %d: %s → %s", i+1, edakit.Num(cb, 4), edakit.Num(ca, 4)))
		}
		q := bbRule + edakit.Sprintf(" Use %s. Block B sits at (%d, %d) and is on two nets: %s. The annealer moves B to the empty location (%d, %d). What is the change in total cost, new minus old? Round to 2 decimal places.", qText(ns...), from.x, from.y, strings.Join(lines, "; "), to.x, to.y)
		return num(q, edakit.Num(d, 2),
			edakit.Sprintf("Only the nets on B change. %s. The change is %s; VPR computes exactly this incremental delta instead of re-costing the whole placement.", strings.Join(parts, "; "), edakit.Num(d, 2)),
			map[string]string{"delta": edakit.Num(d, 2)})
	}
}

func vprAnnealing(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		dc := float64(edakit.Int(1, 80)) / 4
		t := float64(edakit.Int(2, 60)) / 2
		p := math.Exp(-dc / t)
		q := edakit.Sprintf("At temperature T = %s, the annealer proposes a swap that raises the placement cost by ΔC = %s. With the Metropolis rule, what is the probability that the swap is accepted? Round to 3 decimal places.", edakit.Num(t, 1), edakit.Num(dc, 2))
		return num(q, edakit.Num(p, 3),
			edakit.Sprintf("An uphill move is accepted with probability e^(-ΔC/T) = e^(-%s/%s) = %s. A move with ΔC ≤ 0 would always be accepted.", edakit.Num(dc, 2), edakit.Num(t, 1), edakit.Num(p, 3)),
			map[string]string{"ratio": edakit.Num(dc/t, 4)})
	case 2:
		g := 8
		a := pt{edakit.Int(1, g), edakit.Int(1, g)}
		b := pt{edakit.Int(1, g), edakit.Int(1, g)}
		for a == b {
			b = pt{edakit.Int(1, g), edakit.Int(1, g)}
		}
		type net struct {
			hasA, hasB bool
			fixed      []pt
		}
		var nets []net
		nets = append(nets, net{true, false, randPts(edakit.Int(1, 2), g)})
		nets = append(nets, net{false, true, randPts(edakit.Int(1, 2), g)})
		if edakit.Int(0, 1) == 0 {
			nets = append(nets, net{true, true, randPts(edakit.Int(0, 1), g)})
		} else {
			nets = append(nets, net{true, false, randPts(1, g)})
		}
		cost := func(pa, pb pt) int {
			s := 0
			for _, n := range nets {
				var ps []pt
				if n.hasA {
					ps = append(ps, pa)
				}
				if n.hasB {
					ps = append(ps, pb)
				}
				ps = append(ps, n.fixed...)
				sx, sy := spans(ps)
				s += sx + sy
			}
			return s
		}
		before, after := cost(a, b), cost(b, a)
		var desc []string
		for i, n := range nets {
			var t []string
			if n.hasA {
				t = append(t, "A")
			}
			if n.hasB {
				t = append(t, "B")
			}
			if len(n.fixed) > 0 {
				t = append(t, "fixed pins at "+ptsText(n.fixed))
			}
			desc = append(desc, edakit.Sprintf("net %d joins %s", i+1, strings.Join(t, " and ")))
		}
		q := edakit.Sprintf("Every net here has at most 3 terminals, so its cost is (x_max - x_min + 1) + (y_max - y_min + 1). Block A is at (%d, %d) and block B at (%d, %d); %s. The annealer proposes swapping A and B. What is ΔC, the new cost minus the old cost of these nets?", a.x, a.y, b.x, b.y, strings.Join(desc, "; "))
		return num(q, edakit.IntStr(after-before),
			edakit.Sprintf("Before the swap the nets cost %d; after it they cost %d, so ΔC = %d. A net containing both A and B and nothing else cannot change, because the same two locations remain occupied.", before, after, after-before),
			map[string]string{"before": edakit.IntStr(before), "after": edakit.IntStr(after)})
	default:
		dc := float64(edakit.Int(1, 40)) / 2
		p := float64(edakit.Int(5, 95)) / 100
		t := -dc / math.Log(p)
		q := edakit.Sprintf("An annealer should accept an uphill move of ΔC = %s with probability %s under the Metropolis rule. What temperature achieves that? Round to 2 decimal places.", edakit.Num(dc, 1), edakit.Num(p, 2))
		return num(q, edakit.Num(t, 2),
			edakit.Sprintf("Solve e^(-ΔC/T) = p for T: T = -ΔC / ln p = -%s / ln %s = %s.", edakit.Num(dc, 1), edakit.Num(p, 2), edakit.Num(t, 2)),
			map[string]string{"ln_p": edakit.Num(math.Log(p), 4)})
	}
}

func alphaFor(r float64) float64 {
	switch {
	case r > 0.96:
		return 0.5
	case r > 0.8:
		return 0.9
	case r > 0.15:
		return 0.95
	default:
		return 0.8
	}
}

const alphaRule = "VPR updates T_new = α·T_old with α = 0.5 if the acceptance rate R > 0.96, 0.9 if 0.8 < R ≤ 0.96, 0.95 if 0.15 < R ≤ 0.8, and 0.8 if R ≤ 0.15."

func temperatureSchedule(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		n := edakit.Int(5, 7)
		var cs []float64
		var s []string
		mean := 0.0
		for i := 0; i < n; i++ {
			c := float64(edakit.Int(200, 400))
			cs = append(cs, c)
			s = append(s, edakit.Num(c, 0))
			mean += c
		}
		mean /= float64(n)
		v := 0.0
		for _, c := range cs {
			v += (c - mean) * (c - mean)
		}
		sd := math.Sqrt(v / float64(n))
		t0 := 20 * sd
		q := edakit.Sprintf("To choose its starting temperature, VPR makes random moves from a random placement, accepts them all, and records the costs: %s. It sets T_0 = 20 times the standard deviation of those costs. Using the population standard deviation, what is T_0? Round to 1 decimal place.", strings.Join(s, ", "))
		return num(q, edakit.Num(t0, 1),
			edakit.Sprintf("Mean = %s, population standard deviation = %s, so T_0 = 20 × %s = %s. Scaling by the cost spread makes almost every early move acceptable.", edakit.Num(mean, 3), edakit.Num(sd, 4), edakit.Num(sd, 4), edakit.Num(t0, 1)),
			map[string]string{"sigma": edakit.Num(sd, 4)})
	case 2:
		t := float64(edakit.Int(20, 400)) / 4
		r1 := float64(edakit.Int(1, 99)) / 100
		r2 := float64(edakit.Int(1, 99)) / 100
		a1, a2 := alphaFor(r1), alphaFor(r2)
		t2 := t * a1 * a2
		q := alphaRule + edakit.Sprintf(" The annealer is at T = %s. The next temperature step finishes with R = %s, and the one after with R = %s. What is the temperature after both updates? Round to 3 decimal places.", edakit.Num(t, 2), edakit.Num(r1, 2), edakit.Num(r2, 2))
		return num(q, edakit.Num(t2, 3),
			edakit.Sprintf("R = %s gives α = %s and R = %s gives α = %s, so T = %s × %s × %s = %s. The schedule lingers (α = 0.95) where the acceptance rate is moderate, which is where the cost falls fastest.", edakit.Num(r1, 2), edakit.Num(a1, 2), edakit.Num(r2, 2), edakit.Num(a2, 2), edakit.Num(t, 2), edakit.Num(a1, 2), edakit.Num(a2, 2), edakit.Num(t2, 3)),
			map[string]string{"alpha1": edakit.Num(a1, 2), "alpha2": edakit.Num(a2, 2)})
	default:
		if edakit.Int(0, 1) == 0 {
			in := edakit.Pick([]int{1, 2, 5, 10})
			nb := edakit.Scaled(ctx, 20, 3000)
			m := float64(in) * math.Pow(float64(nb), 4.0/3.0)
			q := edakit.Sprintf("VPR makes InnerNum·N^(4/3) moves at each temperature, where N is the number of blocks. How many moves per temperature does it make for N = %d blocks with InnerNum = %d? Round to the nearest integer.", nb, in)
			return num(q, edakit.Num(m, 0),
				edakit.Sprintf("%d × %d^(4/3) = %d × %s = %s. Raising InnerNum buys quality with run time almost linearly.", in, nb, in, edakit.Num(math.Pow(float64(nb), 4.0/3.0), 2), edakit.Num(m, 0)),
				map[string]string{"n_pow": edakit.Num(math.Pow(float64(nb), 4.0/3.0), 2)})
		}
		cost := float64(edakit.Int(50, 400))
		nn := edakit.Int(50, 400)
		thr := 0.005 * cost / float64(nn)
		t := float64(edakit.Int(1, 40)) / 100
		a := 0.8
		k := 0
		for tt := t; tt >= thr; tt *= a {
			k++
		}
		q := edakit.Sprintf("VPR stops annealing when T < 0.005·Cost / N_nets. The placement cost is %s over %d nets, the temperature is %s, and from here every update uses α = 0.8. How many more temperature updates happen before the exit test passes?", edakit.Num(cost, 0), nn, edakit.Num(t, 2))
		return num(q, edakit.IntStr(k),
			edakit.Sprintf("The threshold is 0.005 × %s / %d = %s. Multiplying %s by 0.8 repeatedly first drops below it after %d updates. Normalizing by cost per net means the test asks whether a typical uphill move per net is still affordable.", edakit.Num(cost, 0), nn, edakit.Num(thr, 6), edakit.Num(t, 2), k),
			map[string]string{"threshold": edakit.Num(thr, 6)})
	}
}

func rangeLimit(ctx generator.GeneratorContext) generator.Problem {
	upd := func(r, acc float64, hi float64) float64 {
		n := r * (1 - 0.44 + acc)
		if n < 1 {
			n = 1
		}
		if n > hi {
			n = hi
		}
		return n
	}
	rule := "VPR updates the range limit as R_new = R_old·(1 - 0.44 + R_accept), clamped between 1 and the larger grid dimension."
	switch whichKP(ctx) {
	case 1:
		g := edakit.Int(10, 60)
		r := float64(edakit.Int(4, 4*g)) / 4
		acc := float64(edakit.Int(1, 99)) / 100
		n := upd(r, acc, float64(g))
		q := rule + edakit.Sprintf(" On a %d × %d grid the limit is %s and the last temperature accepted a fraction %s of moves. What is the new range limit? Round to 2 decimal places.", g, g, edakit.Num(r, 2), edakit.Num(acc, 2))
		return num(q, edakit.Num(n, 2),
			edakit.Sprintf("%s × (0.56 + %s) = %s, then clamp to [1, %d]: %s. The rule steers the acceptance rate toward 0.44.", edakit.Num(r, 2), edakit.Num(acc, 2), edakit.Num(r*(0.56+acc), 4), g, edakit.Num(n, 2)),
			map[string]string{"raw": edakit.Num(r*(0.56+acc), 4)})
	case 2:
		g := edakit.Int(10, 60)
		r := float64(edakit.Int(4, 4*g)) / 4
		a1 := float64(edakit.Int(1, 99)) / 100
		a2 := float64(edakit.Int(1, 99)) / 100
		n1 := upd(r, a1, float64(g))
		n2 := upd(n1, a2, float64(g))
		q := rule + edakit.Sprintf(" On a %d × %d grid the limit is %s. The next two temperatures accept fractions %s and then %s of their moves. What is the range limit after both updates? Round to 2 decimal places.", g, g, edakit.Num(r, 2), edakit.Num(a1, 2), edakit.Num(a2, 2))
		return num(q, edakit.Num(n2, 2),
			edakit.Sprintf("First update: %s (after clamping). Second: %s × (0.56 + %s) = %s, which after clamping is %s.", edakit.Num(n1, 4), edakit.Num(n1, 4), edakit.Num(a2, 2), edakit.Num(n1*(0.56+a2), 4), edakit.Num(n2, 2)),
			map[string]string{"first": edakit.Num(n1, 4)})
	default:
		nx, ny := edakit.Int(6, 20), edakit.Int(6, 20)
		x, y := edakit.Int(1, nx), edakit.Int(1, ny)
		rf := float64(edakit.Int(4, 40)) / 4
		r := int(math.Floor(rf))
		cx := minInt(nx, x+r) - maxInt(1, x-r) + 1
		cy := minInt(ny, y+r) - maxInt(1, y-r) + 1
		c := cx*cy - 1
		q := edakit.Sprintf("A block sits at (%d, %d) on a %d × %d grid of logic-block sites numbered from 1. With range limit %s, the move generator picks a target site within floor(R) sites of the block in both x and y, inside the grid and not the block's own site. How many target sites are possible?", x, y, nx, ny, edakit.Num(rf, 2))
		return num(q, edakit.IntStr(c),
			edakit.Sprintf("floor(R) = %d. The window covers x from %d to %d (%d columns) and y from %d to %d (%d rows), so %d × %d - 1 = %d sites.", r, maxInt(1, x-r), minInt(nx, x+r), cx, maxInt(1, y-r), minInt(ny, y+r), cy, cx, cy, c),
			map[string]string{"cols": edakit.IntStr(cx), "rows": edakit.IntStr(cy)})
	}
}

func timingCost(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		dmax := float64(10 * edakit.Int(50, 300))
		sl := float64(10 * edakit.Int(0, int(dmax)/10))
		th := edakit.Int(1, 8)
		c := 1 - sl/dmax
		ce := math.Pow(c, float64(th))
		q := edakit.Sprintf("In timing-driven placement, the criticality of a connection is crit = 1 - slack / D_max, and its timing cost weight is crit^θ. The critical-path delay is D_max = %s ps, the connection's slack is %s ps and the criticality exponent is θ = %d. What is crit^θ? Round to 3 decimal places.", edakit.Num(dmax, 0), edakit.Num(sl, 0), th)
		return num(q, edakit.Num(ce, 3),
			edakit.Sprintf("crit = 1 - %s/%s = %s, and crit^%d = %s. A large θ makes the placer ignore all but the most critical connections.", edakit.Num(sl, 0), edakit.Num(dmax, 0), edakit.Num(c, 4), th, edakit.Num(ce, 3)),
			map[string]string{"crit": edakit.Num(c, 4)})
	case 2:
		dmax := float64(10 * edakit.Int(80, 200))
		th := edakit.Int(1, 4)
		tot := 0.0
		var parts, lines []string
		for i := 0; i < 3; i++ {
			sl := float64(10 * edakit.Int(0, int(dmax)/10))
			d := float64(10 * edakit.Int(5, 60))
			c := math.Pow(1-sl/dmax, float64(th))
			tot += d * c
			lines = append(lines, edakit.Sprintf("delay %s ps, slack %s ps", edakit.Num(d, 0), edakit.Num(sl, 0)))
			parts = append(parts, edakit.Sprintf("%s × %s", edakit.Num(d, 0), edakit.Num(c, 4)))
		}
		q := edakit.Sprintf("VPR's timing cost sums delay × crit^θ over connections, with crit = 1 - slack / D_max. With D_max = %s ps and θ = %d, three connections have: %s. What is their timing cost, in ps? Round to 2 decimal places.", edakit.Num(dmax, 0), th, strings.Join(lines, "; "))
		return num(q, edakit.Num(tot, 2),
			edakit.Sprintf("Sum: %s = %s ps.", strings.Join(parts, " + "), edakit.Num(tot, 2)),
			map[string]string{"timing_cost": edakit.Num(tot, 2)})
	default:
		lam := float64(edakit.Int(1, 9)) / 10
		pt := float64(edakit.Int(500, 5000))
		pb := float64(edakit.Int(500, 5000))
		dt := float64(edakit.Int(-80, 80))
		db := float64(edakit.Int(-80, 80))
		for dt == 0 && db == 0 {
			dt = float64(edakit.Int(-80, 80))
		}
		d := lam*dt/pt + (1-lam)*db/pb
		q := edakit.Sprintf("A timing-driven annealer scores a swap by ΔC = λ·ΔTimingCost / PrevTimingCost + (1 - λ)·ΔBBCost / PrevBBCost, where the Prev values were recorded at the start of this temperature. With λ = %s, PrevTimingCost = %s, PrevBBCost = %s, a swap changes the timing cost by %s and the bounding-box cost by %s. What is ΔC? Round to 4 decimal places.", edakit.Num(lam, 1), edakit.Num(pt, 0), edakit.Num(pb, 0), edakit.Num(dt, 0), edakit.Num(db, 0))
		return num(q, edakit.Num(d, 4),
			edakit.Sprintf("%s × %s / %s + %s × %s / %s = %s. Normalizing each term by its previous total lets λ trade the two objectives regardless of their units.", edakit.Num(lam, 1), edakit.Num(dt, 0), edakit.Num(pt, 0), edakit.Num(1-lam, 1), edakit.Num(db, 0), edakit.Num(pb, 0), edakit.Num(d, 4)),
			map[string]string{"delta": edakit.Num(d, 4)})
	}
}

func analytical(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		n := edakit.Int(2, 5)
		var parts []string
		sw, swx := 0.0, 0.0
		for i := 0; i < n; i++ {
			x := float64(edakit.Int(0, 40))
			w := float64(edakit.Int(1, 4))
			sw += w
			swx += w * x
			parts = append(parts, edakit.Sprintf("x = %s (weight %s)", edakit.Num(x, 0), edakit.Num(w, 0)))
		}
		x := swx / sw
		q := edakit.Sprintf("A quadratic placer minimizes Σ w_i (x - x_i)² for one movable block connected to fixed pins at %s. What is the optimal x? Round to 2 decimal places.", strings.Join(parts, "; "))
		return num(q, edakit.Num(x, 2),
			edakit.Sprintf("Setting the derivative to zero gives the weighted mean: %s / %s = %s.", edakit.Num(swx, 0), edakit.Num(sw, 0), edakit.Num(x, 2)),
			map[string]string{"sum_w": edakit.Num(sw, 0), "sum_wx": edakit.Num(swx, 0)})
	case 2:
		p1 := pickDistinct([]int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20}, edakit.Int(1, 2))
		p2 := pickDistinct([]int{10, 14, 18, 22, 26, 30, 34, 38}, edakit.Int(1, 2))
		w := float64(edakit.Int(1, 4))
		s1, s2 := 0.0, 0.0
		for _, x := range p1 {
			s1 += float64(x)
		}
		for _, x := range p2 {
			s2 += float64(x)
		}
		a, b := float64(len(p1))+w, float64(len(p2))+w
		det := a*b - w*w
		x1 := (s1*b + w*s2) / det
		q := edakit.Sprintf("Two movable blocks are placed in one dimension by minimizing quadratic wirelength. Block 1 connects with weight 1 to fixed pins at x = %s; block 2 connects with weight 1 to fixed pins at x = %s; the two blocks connect to each other with weight %s. What is the optimal x of block 1? Round to 2 decimal places.", joinInts(p1), joinInts(p2), edakit.Num(w, 0))
		return num(q, edakit.Num(x1, 2),
			edakit.Sprintf("Zero gradient gives (%s) x1 - %s x2 = %s and -%s x1 + (%s) x2 = %s. Solving, x1 = %s. This linear system, the Laplacian plus the fixed-pin diagonal, is what an analytical placer solves for every block at once.", edakit.Num(a, 0), edakit.Num(w, 0), edakit.Num(s1, 0), edakit.Num(w, 0), edakit.Num(b, 0), edakit.Num(s2, 0), edakit.Num(x1, 2)),
			map[string]string{"det": edakit.Num(det, 0)})
	default:
		step := edakit.Int(4, 9)
		off := edakit.Int(1, step-1)
		var cols []int
		for c := off; c < 60; c += step {
			cols = append(cols, c)
		}
		x := float64(edakit.Int(0, 240)) / 4
		best := math.Inf(1)
		for _, c := range cols {
			if d := math.Abs(float64(c) - x); d < best {
				best = d
			}
		}
		q := edakit.Sprintf("On this FPGA, DSP blocks may sit only in the columns x = %d, %d, %d, ... (every %d columns). The analytical solution places a DSP block at x = %s. Legalizing it to the nearest DSP column, how far does it move in x? Round to 2 decimal places.", cols[0], cols[1], cols[2], step, edakit.Num(x, 2))
		return num(q, edakit.Num(best, 2),
			edakit.Sprintf("The nearest DSP column is %s columns away. Heterogeneous columns are why FPGA analytical placers spread and legalize each block type separately.", edakit.Num(best, 2)),
			map[string]string{"move": edakit.Num(best, 2)})
	}
}

func hasPt(ps []pt, p pt) bool {
	for _, q := range ps {
		if q == p {
			return true
		}
	}
	return false
}
