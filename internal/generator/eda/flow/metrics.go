package flow

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// ppa: dynamic power, voltage and frequency scaling, and maximum frequency.
func ppa(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		a := float64(edakit.Int(5, 30)) / 100
		c := 50 * edakit.Int(2, 40)
		v := float64(edakit.Int(6, 12)) / 10
		f := edakit.Pick([]float64{0.5, 0.8, 1, 1.2, 1.5, 2, 2.5, 3})
		p := a * float64(c) * 1e-12 * v * v * f * 1e9 * 1000
		ans := edakit.Num(p, 2)
		q := edakit.Sprintf("A block switches an effective capacitance of %d pF with activity factor %s, at %s V and %s GHz. What is its dynamic power $\\alpha C V^2 f$, in mW? Round to 2 decimal places.", c, f2(a, 2), f2(v, 1), f2(f, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$%s \\cdot %d \\times 10^{-12} \\cdot %s^2 \\cdot %s \\times 10^9$ W is %s mW.", f2(a, 2), c, f2(v, 1), f2(f, 1), ans), nil)
	case 2:
		p0 := float64(edakit.Int(20, 500))
		v0 := float64(edakit.Int(8, 12)) / 10
		v1 := v0 - float64(edakit.Int(1, 3))/10
		f0 := edakit.Pick([]float64{1, 1.5, 2, 2.5})
		f1 := f0
		if ctx.Difficulty > 0.5 && edakit.Pick([]bool{true, false}) {
			f1 = f0 * edakit.Pick([]float64{0.6, 0.75, 0.8})
		}
		p1 := p0 * (v1 / v0) * (v1 / v0) * (f1 / f0)
		ans := edakit.Num(p1, 2)
		clk := edakit.Sprintf("its clock unchanged at %s GHz", f2(f0, 2))
		if f1 != f0 {
			clk = edakit.Sprintf("its clock lowered to %s GHz", f2(f1, 2))
		}
		q := edakit.Sprintf("A core burns %s mW of dynamic power at %s V and %s GHz. Its supply is lowered to %s V, with %s and activity and capacitance unchanged. What is its new dynamic power, in mW? Round to 2 decimal places.", f2(p0, 0), f2(v0, 1), f2(f0, 2), f2(v1, 1), clk)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$%s \\cdot (%s/%s)^2 \\cdot (%s/%s) = %s$ mW.", f2(p0, 0), f2(v1, 1), f2(v0, 1), f2(f1, 2), f2(f0, 2), ans), nil)
	default:
		tcq, td, tsu := edakit.Int(30, 90), edakit.Int(200, 1500), edakit.Int(20, 60)
		t := tcq + td + tsu
		ans := edakit.Num(1e6/float64(t), 1)
		q := edakit.Sprintf("The critical register-to-register path has %d ps clock-to-output, %d ps of logic and wire delay and %d ps setup. What is the maximum clock frequency, in MHz? Round to 1 decimal place.", tcq, td, tsu)
		return edakit.Numeric(q, ans, edakit.Sprintf("$T_{min} = %d$ ps, so $f_{max} = 10^6/%d = %s$ MHz.", t, t, ans), nil)
	}
}

// dpoint is a design point; d and a are held in tenths so that comparisons are exact.
type dpoint struct {
	name    string
	p, d, a int
}

func dominates(x, y dpoint) bool {
	return x.p <= y.p && x.d <= y.d && x.a <= y.a && (x.p < y.p || x.d < y.d || x.a < y.a)
}

func randPoint(name string) dpoint {
	return dpoint{name, edakit.Int(8, 20), edakit.Int(8, 16), edakit.Int(15, 30)}
}

func pointStr(x dpoint) string {
	return edakit.Sprintf("%s: %d mW, %s ns, %s square mm", x.name, x.p, f2(float64(x.d)/10, 1), f2(float64(x.a)/10, 1))
}

// pareto: pairwise dominance (KP 1), Pareto-front size (KP 2) and weighted selection (KP 3).
func pareto(ctx generator.GeneratorContext) generator.Problem {
	const lower = " Lower is better for power, critical-path delay and area."
	switch kp(ctx, 3) {
	case 1:
		a, b := randPoint("A"), randPoint("B")
		switch edakit.Int(0, 2) {
		case 0:
			b = dpoint{"B", a.p + edakit.Int(0, 3), a.d + edakit.Int(0, 2), a.a + edakit.Int(1, 4)}
		case 1:
			a = dpoint{"A", b.p + edakit.Int(1, 3), b.d + edakit.Int(0, 2), b.a + edakit.Int(0, 3)}
		}
		if a.p == b.p && a.d == b.d && a.a == b.a {
			b.p++
		}
		ans := "Neither dominates the other"
		if dominates(a, b) {
			ans = "A dominates B"
		} else if dominates(b, a) {
			ans = "B dominates A"
		}
		q := "Two implementations of a block are " + pointStr(a) + "; " + pointStr(b) + "." + lower + " Which statement is true?"
		return edakit.Choice(q, ans, "A point dominates another when it is no worse on every metric and strictly better on at least one.",
			pickOthers([]string{"A dominates B", "B dominates A", "Neither dominates the other", "Each dominates the other"}, 3, ans)...)
	case 2:
		n := edakit.Scaled(ctx, 4, 6)
		pts := make([]dpoint, n)
		for i := range pts {
			pts[i] = randPoint(string(rune('P' + i)))
		}
		count := 0
		var front []string
		for i := range pts {
			dom := false
			for j := range pts {
				if i != j && dominates(pts[j], pts[i]) {
					dom = true
				}
			}
			if !dom {
				count++
				front = append(front, pts[i].name)
			}
		}
		var lines []string
		for _, x := range pts {
			lines = append(lines, pointStr(x))
		}
		q := "Design-space exploration produced " + strings.Join(lines, "; ") + "." + lower + " How many of these points are Pareto-optimal?"
		var cands []int
		for c := 1; c <= n; c++ {
			cands = append(cands, c)
		}
		return edakit.Choice(q, edakit.IntStr(count),
			"The points that no other point dominates are "+strings.Join(front, ", ")+edakit.Sprintf(", so %d are Pareto-optimal.", count),
			edakit.Distinct(count, 3, edakit.Shuffle(cands)...)...)
	default:
		for {
			pts := make([]dpoint, 4)
			for i := range pts {
				pts[i] = randPoint(string(rune('P' + i)))
			}
			wp := float64(edakit.Int(1, 4))
			wd := float64(edakit.Pick([]int{5, 10, 20}))
			wa := float64(edakit.Int(1, 5))
			best, bi, second := math.Inf(1), -1, math.Inf(1)
			for i, x := range pts {
				c := wp*float64(x.p) + wd*float64(x.d)/10 + wa*float64(x.a)/10
				if c < best {
					second, best, bi = best, c, i
				} else if c < second {
					second = c
				}
			}
			if second-best < 0.05 {
				continue
			}
			var lines, names []string
			for _, x := range pts {
				lines = append(lines, pointStr(x))
				names = append(names, x.name)
			}
			q := edakit.Sprintf("A flow scores each design with cost $%s P + %s D + %s A$, where P is power in mW, D is delay in ns and A is area in square mm. The candidates are %s. Which design has the lowest cost?", f2(wp, 0), f2(wd, 0), f2(wa, 0), strings.Join(lines, "; "))
			return edakit.Choice(q, pts[bi].name,
				edakit.Sprintf("%s has the lowest weighted cost, %s.", pts[bi].name, f2(best, 2)), pickOthers(names, 3, pts[bi].name)...)
		}
	}
}

func dpw(d, s float64) float64 {
	return math.Pi*(d/2)*(d/2)/s - math.Pi*d/math.Sqrt(2*s)
}

// diePerWafer: area ratio (KP 1), edge-corrected count (KP 2) and count with scribe lanes (KP 3).
func diePerWafer(ctx generator.GeneratorContext) generator.Problem {
	d := edakit.Pick([]float64{150, 200, 300})
	switch kp(ctx, 3) {
	case 1:
		s := float64(edakit.Int(20, 600))
		ans := edakit.Num(math.Pi*(d/2)*(d/2)/s, 1)
		q := edakit.Sprintf("Ignoring edge loss, how many dies of %s square mm fit by area on a %s mm wafer? Round to 1 decimal place.", f2(s, 0), f2(d, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\pi (%s/2)^2 / %s = %s$.", f2(d, 0), f2(s, 0), ans), nil)
	case 2:
		for {
			s := float64(edakit.Int(20, 600))
			x := dpw(d, s)
			if nearInt(x, 0.05) {
				continue
			}
			ans := edakit.IntStr(int(math.Floor(x)))
			q := edakit.Sprintf("Using $\\mathrm{DPW} = \\pi (d/2)^2/S - \\pi d/\\sqrt{2S}$, how many whole dies of area %s square mm fit on a %s mm wafer? Round down.", f2(s, 0), f2(d, 0))
			return edakit.Numeric(q, ans,
				edakit.Sprintf("$%s - %s = %s$, so %s dies.", f2(math.Pi*(d/2)*(d/2)/s, 2), f2(math.Pi*d/math.Sqrt(2*s), 2), f2(x, 2), ans), nil)
		}
	default:
		for {
			w := float64(edakit.Int(30, 200)) / 10
			h := float64(edakit.Int(30, 200)) / 10
			sc := edakit.Pick([]float64{0.1, 0.2})
			s := (w + sc) * (h + sc)
			x := dpw(d, s)
			if nearInt(x, 0.05) {
				continue
			}
			ans := edakit.IntStr(int(math.Floor(x)))
			q := edakit.Sprintf("A die measures %s mm by %s mm and the scribe lane adds %s mm to each dimension. Using $\\mathrm{DPW} = \\pi (d/2)^2/S - \\pi d/\\sqrt{2S}$ with S the stepped area, how many whole dies fit on a %s mm wafer? Round down.", f2(w, 1), f2(h, 1), f2(sc, 1), f2(d, 0))
			return edakit.Numeric(q, ans,
				edakit.Sprintf("$S = %s \\cdot %s = %s$ square mm, and the formula gives %s, so %s dies.", f2(w+sc, 1), f2(h+sc, 1), f2(s, 2), f2(x, 2), ans), nil)
		}
	}
}

// yieldCost: Poisson yield, good dies per wafer, cost per good die.
func yieldCost(ctx generator.GeneratorContext) generator.Problem {
	d0 := float64(edakit.Int(1, 10)) * 0.05
	amm := float64(edakit.Int(20, 600))
	y := math.Exp(-d0 * amm / 100)
	switch kp(ctx, 3) {
	case 1:
		ans := edakit.Num(100*y, 1)
		q := edakit.Sprintf("Defects are Poisson-distributed at %s per square cm. Under $Y = e^{-D_0 A}$, what fraction of %s square mm dies is defect-free, in percent? Round to 1 decimal place.", f2(d0, 2), f2(amm, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("%s square mm is %s square cm, so $Y = e^{-%s \\cdot %s} = %s\\%%$.", f2(amm, 0), f2(amm/100, 2), f2(d0, 2), f2(amm/100, 2), ans), nil)
	case 2:
		for {
			n := int(math.Floor(dpw(300, amm)))
			x := float64(n) * y
			if nearInt(x, 0.05) || math.Floor(x) != math.Floor(float64(n)*edakit.Round(y, 4)) {
				amm = float64(edakit.Int(20, 600))
				y = math.Exp(-d0 * amm / 100)
				continue
			}
			ans := edakit.IntStr(int(math.Floor(x)))
			q := edakit.Sprintf("A 300 mm wafer holds %d candidate dies of %s square mm. Defect density is %s per square cm and yield follows $Y = e^{-D_0 A}$. How many good dies does the wafer give? Round down.", n, f2(amm, 0), f2(d0, 2))
			return edakit.Numeric(q, ans,
				edakit.Sprintf("$Y = e^{-%s \\cdot %s} = %s$, and $%d \\cdot %s = %s$, so %s good dies.", f2(d0, 2), f2(amm/100, 2), f2(y, 4), n, f2(y, 4), f2(x, 2), ans), nil)
		}
	default:
		for {
			n := int(math.Floor(dpw(300, amm)))
			x := float64(n) * y
			good := int(math.Floor(x))
			if nearInt(x, 0.05) || good < 1 || math.Floor(x) != math.Floor(float64(n)*edakit.Round(y, 4)) {
				amm = float64(edakit.Int(20, 600))
				y = math.Exp(-d0 * amm / 100)
				continue
			}
			wc := 500 * edakit.Int(6, 40)
			ans := edakit.Num(float64(wc)/float64(good), 2)
			q := edakit.Sprintf("A 300 mm wafer costs %s dollars and holds %d candidate dies of %s square mm. Defect density is %s per square cm and yield follows $Y = e^{-D_0 A}$. Rounding good dies down to a whole number, what is the cost per good die, in dollars? Round to 2 decimal places.", comma(wc), n, f2(amm, 0), f2(d0, 2))
			return edakit.Numeric(q, ans,
				edakit.Sprintf("$Y = %s$, good dies $= \\lfloor %d \\cdot %s \\rfloor = %d$, and $%d/%d = %s$ dollars.", f2(y, 4), n, f2(y, 4), good, wc, good, ans), nil)
		}
	}
}

// respinCost: direct cost, total with delay, and the expected value of lowering respin risk.
func respinCost(ctx generator.GeneratorContext) generator.Problem {
	mask := float64(edakit.Int(2, 30))
	eng, months := edakit.Int(10, 60), edakit.Int(2, 6)
	rate := edakit.Pick([]int{15, 20, 25, 30})
	direct := mask + float64(eng*months*rate)/1000
	setup := edakit.Sprintf("A respin needs a new mask set costing %s million dollars and %d engineers for %d months at %d thousand dollars per engineer-month.", f2(mask, 0), eng, months, rate)
	switch kp(ctx, 3) {
	case 1:
		ans := edakit.Num(direct, 2)
		q := setup + " What is its direct cost, in millions of dollars? Round to 2 decimal places."
		return edakit.Numeric(q, ans, edakit.Sprintf("$%s + %d \\cdot %d \\cdot %d/1000 = %s$ million dollars.", f2(mask, 0), eng, months, rate, ans), nil)
	case 2:
		delay := edakit.Int(3, 9)
		rev := float64(edakit.Int(1, 20))
		total := direct + float64(delay)*rev
		ans := edakit.Num(total, 2)
		q := setup + edakit.Sprintf(" It also delays a product that would earn %s million dollars a month by %d months. What is the total cost of the respin, in millions of dollars? Round to 2 decimal places.", f2(rev, 0), delay)
		return edakit.Numeric(q, ans, edakit.Sprintf("Direct cost %s plus lost revenue $%d \\cdot %s$ gives %s million dollars.", f2(direct, 2), delay, f2(rev, 0), ans), nil)
	default:
		c := float64(edakit.Int(10, 80))
		p1 := float64(edakit.Int(15, 50)) / 100
		p2 := p1 - float64(edakit.Int(5, 14))/100
		ans := edakit.Num((p1-p2)*c, 2)
		q := edakit.Sprintf("A respin would cost %s million dollars in total. Extra formal verification would lower the probability of a respin from %s to %s. What is that reduction worth in expected cost, in millions of dollars? Round to 2 decimal places.", f2(c, 0), f2(p1, 2), f2(p2, 2))
		return edakit.Numeric(q, ans, edakit.Sprintf("$(%s - %s) \\cdot %s = %s$ million dollars.", f2(p1, 2), f2(p2, 2), f2(c, 0), ans), nil)
	}
}
