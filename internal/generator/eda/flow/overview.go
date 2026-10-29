package flow

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// scaleVLSI: doubling arithmetic (KP 1), exhaustive-simulation time (KP 2) and the design
// productivity gap (KP 3).
func scaleVLSI(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		n0 := edakit.Pick([]int{1, 2, 3, 5, 6})
		t := edakit.Pick([]int{2, 3})
		p := edakit.Scaled(ctx, 2, 6)
		y := p * t
		ans := n0 << p
		q := edakit.Sprintf("A chip family holds %d billion transistors today. If the count doubles every %d years, how many billion transistors will a chip of the same cost hold after %d years?", n0, t, y)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("%d years is %d doubling periods, so the count is $%d \\cdot 2^{%d} = %d$ billion.", y, p, n0, p, ans),
			map[string]string{"doublings": edakit.IntStr(p)})
	case 2:
		n := edakit.Scaled(ctx, 40, 52)
		r := edakit.Pick([]int{8, 9, 10})
		hours := math.Pow(2, float64(n)) / math.Pow(10, float64(r)) / 3600
		ans := edakit.Num(hours, 2)
		q := edakit.Sprintf("A combinational block has %d input bits. A simulator evaluates $10^{%d}$ input vectors per second. How many hours does it take to apply every possible input vector once? Round to 2 decimal places.", n, r)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("There are $2^{%d}$ vectors. Dividing by $10^{%d}$ per second and by 3600 seconds per hour gives %s hours.", n, r, ans),
			map[string]string{"vectors": edakit.Sprintf("2^%d", n)})
	default:
		rate := edakit.Pick([]int{50, 80, 100, 150, 200})
		days := edakit.Pick([]int{220, 230, 250})
		m := edakit.Pick([]int{2, 5, 8, 10, 20, 40, 50, 100})
		years := float64(m) * 1e6 / float64(rate*days)
		ans := edakit.Num(years, 1)
		q := edakit.Sprintf("An engineer drawing layout by hand completes %d transistors per working day, and a year has %d working days. How many engineer-years would it take to draw a design of %d million transistors? Round to 1 decimal place.", rate, days, m)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$%d \\times 10^6 / (%d \\cdot %d) = %s$ engineer-years, before any checking or rework.", m, rate, days, ans),
			map[string]string{"per_year": edakit.IntStr(rate * days)})
	}
}

// hierarchy: flattened cell counts, instance counts and the saving of storing definitions once.
func hierarchy(ctx generator.GeneratorContext) generator.Problem {
	m1, m2 := edakit.Int(2, 6), edakit.Int(1, 5)
	k1, j := edakit.Int(5, 30), edakit.Scaled(ctx, 1, 4)
	k3, k2 := edakit.Int(10, 40), edakit.Int(20, 80)
	desc := edakit.Sprintf("The top module instantiates %d copies of block A and %d %s of block B. Block A contains %d standard cells of its own plus %d %s of block C. Block C contains %d standard cells. Block B contains %d standard cells.", m1, m2, plural(m2, "copy", "copies"), k1, j, plural(j, "instance", "instances"), k3, k2)
	flat := m1*(k1+j*k3) + m2*k2
	switch kp(ctx, 3) {
	case 1:
		q := desc + " How many standard cells does the fully flattened design contain?"
		return edakit.Numeric(q, edakit.IntStr(flat),
			edakit.Sprintf("Multiply along each path from the top: $%d(%d + %d \\cdot %d) + %d \\cdot %d = %d$ cells.", m1, k1, j, k3, m2, k2, flat),
			map[string]string{"flat_cells": edakit.IntStr(flat)})
	case 2:
		inst := m1 + m2 + m1*j
		q := desc + " Counting only block instances (A, B and C, not standard cells), how many block instances does the flattened design contain?"
		return edakit.Numeric(q, edakit.IntStr(inst),
			edakit.Sprintf("There are %d A instances and %d B instances, and each A holds %d C instances, so $%d + %d + %d \\cdot %d = %d$.", m1, m2, j, m1, m2, m1, j, inst),
			map[string]string{"c_instances": edakit.IntStr(m1 * j)})
	default:
		defs := k1 + k2 + k3
		ratio := float64(flat) / float64(defs)
		ans := edakit.Num(ratio, 2)
		q := desc + " A hierarchical tool stores each block definition once. How many times more standard cells does the flattened design hold than the definitions of A, B and C together? Round to 2 decimal places."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Flattened: %d cells. Definitions: $%d + %d + %d = %d$ cells. Ratio: %s.", flat, k1, k2, k3, defs, ans),
			map[string]string{"flat_cells": edakit.IntStr(flat), "definition_cells": edakit.IntStr(defs)})
	}
}

// breakEven: the ASIC-versus-FPGA break-even volume as a computed multiple choice.
func breakEven(ctx generator.GeneratorContext) generator.Problem {
	for {
		fa := 1000 * edakit.Pick([]int{2000, 3000, 4000, 5000, 6000, 8000, 10000})
		ff := 1000 * edakit.Pick([]int{50, 100, 200, 300, 400})
		ua := edakit.Int(3, 12)
		delta := edakit.Pick([]int{20, 25, 30, 40, 50, 60, 75, 80, 100})
		uf := ua + delta
		if (fa-ff)%delta != 0 {
			continue
		}
		n := (fa - ff) / delta
		d1 := fa / uf
		d2 := fa / delta
		d3 := (fa - ff) / uf
		d4 := 2 * n
		if d1 == n || d2 == n || d3 == n {
			continue
		}
		q := edakit.Sprintf("An ASIC costs %s dollars of non-recurring engineering and %d dollars per unit. An FPGA solution costs %s dollars of non-recurring engineering and %d dollars per unit. At what volume do the two have equal total cost?",
			comma(fa), ua, comma(ff), uf)
		return edakit.Choice(q, comma(n)+" units",
			edakit.Sprintf("Set the totals equal: %s + %d N = %s + %d N. Then N = (%s - %s)/(%d - %d) = %s units.", comma(fa), ua, comma(ff), uf, comma(fa), comma(ff), uf, ua, comma(n)),
			comma(d1)+" units", comma(d2)+" units", comma(d3)+" units", comma(d4)+" units")
	}
}

// amdahl: overall speedup (KP 1), new runtime (KP 2) and accelerator sizing (KP 3).
func amdahl(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		f := float64(edakit.Int(10, 19)) * 0.05
		s := edakit.Pick([]int{2, 4, 5, 8, 10, 20, 50})
		sp := 1 / ((1 - f) + f/float64(s))
		ans := edakit.Num(sp, 2)
		q := edakit.Sprintf("Profiling shows that %s%% of a program's runtime is a kernel that a hardware accelerator runs %d times faster. The rest stays in software. What is the overall speedup? Round to 2 decimal places.", f2(f*100, 0), s)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$S = 1/((1 - %s) + %s/%d) = %s$.", f2(f, 2), f2(f, 2), s, ans),
			map[string]string{"f": f2(f, 2), "s": edakit.IntStr(s)})
	case 2:
		t := 10 * edakit.Int(5, 60)
		f := float64(edakit.Int(8, 19)) * 0.05
		s := edakit.Pick([]int{3, 4, 5, 8, 10, 16, 25})
		nt := float64(t) * ((1 - f) + f/float64(s))
		ans := edakit.Num(nt, 2)
		q := edakit.Sprintf("A task takes %d ms in software. A fraction %s of that time moves to an accelerator that is %d times faster. How long does the task take now, in milliseconds? Round to 2 decimal places.", t, f2(f, 2), s)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$%d \\cdot ((1 - %s) + %s/%d) = %s$ ms.", t, f2(f, 2), f2(f, 2), s, ans),
			map[string]string{"new_ms": ans})
	default:
		f := edakit.Pick([]float64{0.6, 0.7, 0.75, 0.8, 0.9, 0.95})
		lim := 1 / (1 - f)
		target := edakit.Round(lim*edakit.Pick([]float64{0.5, 0.6, 0.7, 0.8}), 1)
		s := f / (1/target - (1 - f))
		ans := edakit.Num(s, 2)
		q := edakit.Sprintf("A fraction %s of a program's runtime can be accelerated in hardware. How many times faster must the accelerator run for the whole program to speed up by a factor of %s? Round to 2 decimal places.", f2(f, 2), f2(target, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Solve $%s = 1/((1 - %s) + %s/s)$: $s = %s/(1/%s - %s) = %s$. The limit for any accelerator is $1/(1 - %s) = %s$.", f2(target, 1), f2(f, 2), f2(f, 2), f2(f, 2), f2(target, 1), f2(1-f, 2), ans, f2(f, 2), f2(lim, 2)),
			map[string]string{"limit": f2(lim, 2)})
	}
}
