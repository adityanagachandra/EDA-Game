package cmos

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// powerGating: switch IR drop, switch count, and wake-up time under a rush-current limit.
func powerGating(ctx generator.GeneratorContext) generator.Problem {
	i := float64(edakit.Scaled(ctx, 5, 400))
	r := float64(100 * edakit.Int(5, 60))
	switch kpOf(ctx) {
	case 1:
		target := float64(edakit.Int(5, 60))
		n := 10 * ceil(i*r/(10*target))
		dv := i * r / float64(n)
		ans := edakit.Num(dv, 2)
		q := edakit.Sprintf("A switchable domain draws %s mA through %d identical header switches in parallel, each with an on-resistance of %s Ohm. Assuming the current divides evenly, what is the voltage drop across the switches, in mV? Round to 2 decimal places.", edakit.Num(i, 0), n, edakit.Num(r, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("The switches act as one resistor of $%s / %d$ Ohm, so the drop is $%s \\times %s / %d = %s$ mV (mA times Ohm is mV).", edakit.Num(r, 0), n, edakit.Num(i, 0), edakit.Num(r, 0), n, ans), map[string]string{"drop_mv": ans})
	case 2:
		dmax := float64(edakit.Pick([]int{10, 15, 20, 25, 30, 40, 50}))
		n := ceil(i * r / dmax)
		q := edakit.Sprintf("A switchable domain draws %s mA. Each header switch has an on-resistance of %s Ohm, and the drop across the switch network must stay at or below %s mV. What is the minimum number of switches?", edakit.Num(i, 0), edakit.Num(r, 0), edakit.Num(dmax, 0))
		return edakit.Numeric(q, edakit.IntStr(n), edakit.Sprintf("Need $I R / N \\le \\Delta V$, so $N \\ge %s \\times %s / %s = %s$; round up to %d.", edakit.Num(i, 0), edakit.Num(r, 0), edakit.Num(dmax, 0), edakit.Num(i*r/dmax, 3), n), map[string]string{"switches": edakit.IntStr(n)})
	}
	c := float64(edakit.Int(1, 40)) / 2
	v := edakit.Pick([]float64{0.7, 0.75, 0.8, 0.9, 1.0})
	imax := float64(edakit.Int(5, 100))
	t := c * v / imax * 1000
	ans := edakit.Num(t, 1)
	q := edakit.Sprintf("A power-gated domain has %s nF of capacitance on its virtual supply. To limit rush current, the switches are enabled in a daisy chain so the charging current never exceeds %s mA. Treating the current as constant at that limit, how long does it take to charge the virtual rail from 0 to %s V, in ns? Round to 1 decimal place.", edakit.Num(c, 1), edakit.Num(imax, 0), edakit.Num(v, 2))
	return edakit.Numeric(q, ans, edakit.Sprintf("$t = C V / I = %s\\ \\text{nF} \\times %s\\ \\text{V} / %s\\ \\text{mA} = %s$ ns. Wake-up latency is the price of limiting rush current.", edakit.Num(c, 1), edakit.Num(v, 2), edakit.Num(imax, 0), ans), map[string]string{"t_ns": ans})
}

// decapPart is the computed KP 2 of physical_cells.
func decapPart(ctx generator.GeneratorContext) generator.Problem {
	i := float64(edakit.Int(2, 60))
	dt := float64(edakit.Pick([]int{20, 50, 100, 200}))
	dv := float64(edakit.Pick([]int{10, 20, 30, 40, 50}))
	vdd := edakit.Pick([]float64{0.75, 0.8, 0.9})
	c := i * dt / dv
	f := func(x float64) string { return edakit.Num(x, 1) + " pF" }
	q := edakit.Sprintf("A cluster of cells draws an extra %s mA for %s ps when it switches, before the power grid can respond. The local supply ($V_{DD} = %s$ V) may droop by at most %s mV during that time. How much decap must be placed nearby to supply the charge?", edakit.Num(i, 0), edakit.Num(dt, 0), edakit.Num(vdd, 2), edakit.Num(dv, 0))
	return edakit.Choice(q, f(c),
		edakit.Sprintf("The decap supplies $Q = I \\Delta t$ while its voltage falls by $\\Delta V$, so $C = I\\Delta t/\\Delta V = %s \\times %s / %s = %s$ pF (mA times ps over mV is pF).", edakit.Num(i, 0), edakit.Num(dt, 0), edakit.Num(dv, 0), edakit.Num(c, 1)),
		f(2*c), f(c/2), f(i*dt/(vdd*1000)))
}
