package cmos

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// perturb returns a candidate table near the golden one.
func perturb(g table, spread int) table {
	c := table{name: "candidate", idx1: g.idx1, idx2: g.idx2, i1name: g.i1name, i2name: g.i2name, i2unit: g.i2unit}
	for _, row := range g.vals {
		var r []int
		for _, v := range row {
			r = append(r, v+edakit.Int(-spread, spread))
		}
		c.vals = append(c.vals, r)
	}
	return c
}

// goldenCheck: relative error, maximum relative error, and a tolerance count with a floor.
func goldenCheck(ctx generator.GeneratorContext) generator.Problem {
	slews := edakit.Pick(slewGrids)[:3]
	loads := edakit.Pick(loadGrids)[:3]
	gold := randDelayTable("golden SPICE cell_rise", slews, loads)
	switch kpOf(ctx) {
	case 1:
		g := float64(edakit.Int(15, 300))
		c := g + float64(edakit.Int(-12, 12))
		if c == g {
			c = g + 1
		}
		ans := edakit.Num(100*absf(c-g)/g, 2)
		q := edakit.Sprintf("A characterization skill under test reports a delay of %s ps at one grid point where the golden SPICE reference gives %s ps. What is the relative error, in percent? Round to 2 decimal places.", edakit.Num(c, 0), edakit.Num(g, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$|%s - %s| / %s = %s\\%%$, measured against the golden value.", edakit.Num(c, 0), edakit.Num(g, 0), edakit.Num(g, 0), ans), nil)
	case 2:
		cand := perturb(gold, 4)
		best := 0.0
		for i := range gold.vals {
			for j := range gold.vals[i] {
				e := 100 * absf(float64(cand.vals[i][j]-gold.vals[i][j])) / float64(gold.vals[i][j])
				if e > best {
					best = e
				}
			}
		}
		ans := edakit.Num(best, 2)
		q := gold.text() + " " + cand.text() + " Values are in ps. What is the largest relative error of the candidate against the golden table, in percent? Round to 2 decimal places."
		return edakit.Numeric(q, ans, edakit.Sprintf("Compute $|c - g|/g$ at every grid point and take the maximum: %s%%. Small entries (low slew, low load) amplify a fixed absolute error.", ans), map[string]string{"max_rel_pct": ans})
	}
	cand := perturb(gold, 3)
	rel := edakit.Pick([]int{2, 3, 5})
	abs := edakit.Pick([]int{1, 2})
	fails := 0
	for i := range gold.vals {
		for j := range gold.vals[i] {
			g := float64(gold.vals[i][j])
			d := absf(float64(cand.vals[i][j]) - g)
			tol := float64(rel) / 100 * g
			if float64(abs) > tol {
				tol = float64(abs)
			}
			if d > tol+1e-9 {
				fails++
			}
		}
	}
	q := gold.text() + " " + cand.text() + edakit.Sprintf(" Values are in ps. An entry passes if $|c - g| \\le \\max(%d\\%% \\cdot g, %d\\ \\text{ps})$. How many of the entries fail?", rel, abs)
	return edakit.Numeric(q, edakit.IntStr(fails), edakit.Sprintf("Check each entry against its own tolerance, the larger of %d%% of the golden value and %d ps. %d entries exceed it.", rel, abs, fails), map[string]string{"fails": edakit.IntStr(fails)})
}

// surrogatePart is the computed KP 2 of fast_char_surrogate.
func surrogatePart(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 10)
	th := float64(edakit.Pick([]int{2, 3, 4}))
	var parts []string
	cnt := 0
	for i := 0; i < n; i++ {
		s := float64(edakit.Int(5, 60)) / 10
		if s > th {
			cnt++
		}
		parts = append(parts, edakit.Num(s, 1))
	}
	q := edakit.Sprintf("A surrogate characterizer predicts a new corner's delay table and a predictive standard deviation at each of %d grid points: %s ps. The flow accepts a prediction when its standard deviation is at most %s ps and sends every other point to SPICE. How many points go to SPICE?", n, strings.Join(parts, ", "), edakit.Num(th, 0))
	ds := edakit.Distinct(cnt, 3, n-cnt, cnt+1, cnt-1)
	return edakit.Choice(q, edakit.IntStr(cnt), edakit.Sprintf("%d of the %d predictions have a standard deviation above %s ps, so they are re-simulated; the rest are accepted.", cnt, n, edakit.Num(th, 0)), ds...)
}

// paretoPart is the computed KP 2 of dtco_cell_agent.
func paretoPart(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	names := []string{"6T 2-fin", "6.5T 2-fin", "7.5T 3-fin", "9T 3-fin", "9T 4-fin", "10.5T 4-fin"}
	type cand struct {
		name        string
		area, delay int
	}
	var cs []cand
	var desc []string
	for i := 0; i < n; i++ {
		c := cand{names[i], edakit.Int(10, 30), edakit.Int(10, 30)}
		cs = append(cs, c)
		desc = append(desc, edakit.Sprintf("%s (area %d, delay %d)", c.name, c.area, c.delay))
	}
	cnt := 0
	for i, a := range cs {
		dom := false
		for j, b := range cs {
			if i != j && b.area <= a.area && b.delay <= a.delay && (b.area < a.area || b.delay < a.delay) {
				dom = true
			}
		}
		if !dom {
			cnt++
		}
	}
	q := "A DTCO agent evaluates candidate cell architectures on a block benchmark, scoring normalized area and critical-path delay (lower is better for both): " + strings.Join(desc, "; ") + ". How many candidates are Pareto-optimal?"
	ds := edakit.Distinct(cnt, 3, n, 1, n-cnt)
	return edakit.Choice(q, edakit.IntStr(cnt), edakit.Sprintf("A candidate is Pareto-optimal if no other is at least as good on both metrics and strictly better on one. %d candidates survive.", cnt), ds...)
}
