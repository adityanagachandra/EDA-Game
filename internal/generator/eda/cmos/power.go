package cmos

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const alphaDef = "Here $\\alpha$ is the probability that a node makes a 0-to-1 transition in a cycle, so $P = \\alpha C V^2 f$."

// dynamicPower computes alpha C V^2 f, voltage-frequency scaling, and backs out alpha.
func dynamicPower(ctx generator.GeneratorContext) generator.Problem {
	alpha := edakit.Pick([]float64{0.05, 0.08, 0.1, 0.12, 0.15, 0.2, 0.25, 0.3})
	c := float64(10 * edakit.Scaled(ctx, 5, 200))
	v := edakit.Pick([]float64{0.6, 0.7, 0.75, 0.8, 0.9, 1.0, 1.1, 1.2})
	f := edakit.Pick([]float64{0.5, 0.8, 1, 1.2, 1.5, 2, 2.5, 3})
	p := alpha * c * v * v * f
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(p, 2)
		q := edakit.Sprintf("A block has %s pF of switched capacitance, runs at $V_{DD} = %s$ V and $f = %s$ GHz, and has an activity factor $\\alpha = %s$. %s What is its dynamic power in milliwatts? Round to 2 decimal places.", edakit.Num(c, 0), edakit.Num(v, 2), edakit.Num(f, 1), edakit.Num(alpha, 2), alphaDef)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$P = %s \\times %s\\ \\text{pF} \\times %s^2 \\times %s\\ \\text{GHz}$. Since pF times GHz is $10^{-3}$, the result is %s mW.", edakit.Num(alpha, 2), edakit.Num(c, 0), edakit.Num(v, 2), edakit.Num(f, 1), ans),
			map[string]string{"p_mw": ans})
	case 2:
		v2 := v - edakit.Pick([]float64{0.05, 0.1, 0.15, 0.2})
		f2 := edakit.Round(f*v2/v, 4)
		pr := edakit.Round(p, 3)
		p2 := pr * (v2 / v) * (v2 / v) * (f2 / f)
		ans := edakit.Num(p2, 2)
		q := edakit.Sprintf("A block dissipates %s mW of dynamic power at $V_{DD} = %s$ V and $f = %s$ GHz. Voltage scaling lowers the supply to %s V, and the frequency is scaled in proportion to the voltage, to %s GHz. Activity and capacitance are unchanged. What is the new dynamic power in milliwatts? Round to 2 decimal places.", edakit.Num(p, 3), edakit.Num(v, 2), edakit.Num(f, 1), edakit.Num(v2, 2), edakit.Num(f2, 4))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Power scales as $V^2 f$: the factor is $(%s/%s)^2 \\times (%s/%s) = %s$, so the new power is %s mW.", edakit.Num(v2, 2), edakit.Num(v, 2), edakit.Num(f2, 4), edakit.Num(f, 1), edakit.Num(p2/pr, 4), ans),
			map[string]string{"factor": edakit.Num(p2/pr, 4), "p_mw": ans})
	}
	ans := edakit.Num(alpha, 3)
	q := edakit.Sprintf("A power simulation reports %s mW of dynamic power for a block with %s pF of switched capacitance at $V_{DD} = %s$ V and $f = %s$ GHz. %s What activity factor does that imply? Round to 3 decimal places.", edakit.Num(p, 4), edakit.Num(c, 0), edakit.Num(v, 2), edakit.Num(f, 1), alphaDef)
	return edakit.Numeric(q, ans,
		edakit.Sprintf("$\\alpha = P/(C V^2 f) = %s/(%s \\times %s \\times %s) = %s$.", edakit.Num(p, 4), edakit.Num(c, 0), edakit.Num(v*v, 4), edakit.Num(f, 1), ans),
		map[string]string{"alpha": ans})
}

type probGate struct {
	name string
	f    func(a, b float64) float64
	expr string
}

var probGates = []probGate{
	{"NAND2", func(a, b float64) float64 { return 1 - a*b }, "1 - p_A p_B"},
	{"NOR2", func(a, b float64) float64 { return (1 - a) * (1 - b) }, "(1 - p_A)(1 - p_B)"},
	{"AND2", func(a, b float64) float64 { return a * b }, "p_A p_B"},
	{"OR2", func(a, b float64) float64 { return 1 - (1-a)*(1-b) }, "1 - (1 - p_A)(1 - p_B)"},
	{"XOR2", func(a, b float64) float64 { return a*(1-b) + b*(1-a) }, "p_A(1 - p_B) + p_B(1 - p_A)"},
}

// activityFactor: signal probability through a gate, then alpha = p0 p1, then power.
func activityFactor(ctx generator.GeneratorContext) generator.Problem {
	g := edakit.Pick(probGates)
	pa := float64(edakit.Int(1, 9)) / 10
	pb := float64(edakit.Int(1, 9)) / 10
	p1 := g.f(pa, pb)
	alpha := p1 * (1 - p1)
	pre := edakit.Sprintf("The inputs of %s are independent, with probabilities $p_A = %s$ and $p_B = %s$ of being 1 in any cycle, and successive cycles are independent. ", article(g.name), edakit.Num(pa, 1), edakit.Num(pb, 1))
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(p1, 3)
		return edakit.Numeric(pre+"What is the probability that the output is 1? Round to 3 decimal places.", ans,
			edakit.Sprintf("For %s the output-1 probability is $%s = %s$.", article(g.name), g.expr, ans), map[string]string{"p1": ans})
	case 2:
		ans := edakit.Num(alpha, 4)
		return edakit.Numeric(pre+"What is the output's activity factor $\\alpha$, the probability of a 0-to-1 transition per cycle? Round to 4 decimal places.", ans,
			edakit.Sprintf("The output is 1 with probability %s. A rising transition needs a 0 in one cycle and a 1 in the next: $\\alpha = p_0 p_1 = %s \\times %s = %s$.", edakit.Num(p1, 3), edakit.Num(1-p1, 3), edakit.Num(p1, 3), ans), map[string]string{"p1": edakit.Num(p1, 4), "alpha": ans})
	}
	c := float64(edakit.Int(2, 30))
	v := edakit.Pick([]float64{0.7, 0.8, 0.9, 1.0})
	f := edakit.Pick([]float64{0.5, 1, 1.5, 2})
	p := alpha * c * v * v * f
	ans := edakit.Num(p, 3)
	return edakit.Numeric(pre+edakit.Sprintf("The output drives %s fF at $V_{DD} = %s$ V with a %s GHz clock. What is the dynamic power of that output node in microwatts? %s Round to 3 decimal places.", edakit.Num(c, 0), edakit.Num(v, 1), edakit.Num(f, 1), alphaDef), ans,
		edakit.Sprintf("$p_1 = %s$, so $\\alpha = p_0p_1 = %s$. Then $P = \\alpha C V^2 f = %s$ uW (fF times GHz is $10^{-6}$).", edakit.Num(p1, 3), edakit.Num(alpha, 4), ans),
		map[string]string{"alpha": edakit.Num(alpha, 4), "p_uw": ans})
}

// leakage: decades per subthreshold slope, current at another Vt, and DIBL.
func leakage(ctx generator.GeneratorContext) generator.Problem {
	s := float64(edakit.Pick([]int{70, 75, 80, 85, 90, 95, 100}))
	dvt := float64(edakit.Pick([]int{40, 50, 60, 75, 80, 100, 120, 150}))
	switch kpOf(ctx) {
	case 1:
		r := pow(10, dvt/s)
		ans := edakit.Num(r, 2)
		q := edakit.Sprintf("A library offers LVT and HVT flavours whose thresholds differ by %s mV. The subthreshold swing is %s mV per decade. By what factor does the off-current of an LVT device exceed that of the matching HVT device? Round to 2 decimal places.", edakit.Num(dvt, 0), edakit.Num(s, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Subthreshold current falls one decade for every $S$ of threshold, so the ratio is $10^{\\Delta V_t/S} = 10^{%s/%s} = %s$.", edakit.Num(dvt, 0), edakit.Num(s, 0), ans),
			map[string]string{"ratio": ans})
	case 2:
		i1 := float64(edakit.Int(5, 60))
		i2 := i1 / pow(10, dvt/s)
		ans := edakit.Num(i2, 3)
		q := edakit.Sprintf("An SVT transistor leaks %s nA when off. A higher-threshold flavour raises $V_t$ by %s mV, and the subthreshold swing is %s mV/decade. What is the off-current of the higher-threshold device, in nA? Round to 3 decimal places.", edakit.Num(i1, 0), edakit.Num(dvt, 0), edakit.Num(s, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$I = %s / 10^{%s/%s} = %s$ nA.", edakit.Num(i1, 0), edakit.Num(dvt, 0), edakit.Num(s, 0), ans),
			map[string]string{"ioff_na": ans})
	}
	eta := edakit.Pick([]float64{0.05, 0.08, 0.1, 0.12, 0.15})
	v1 := edakit.Pick([]float64{0.9, 1.0, 1.1, 1.2})
	v2 := v1 - edakit.Pick([]float64{0.2, 0.3, 0.4})
	dv := eta * (v1 - v2) * 1000
	r := pow(10, dv/s)
	ans := edakit.Num(r, 3)
	q := edakit.Sprintf("Drain-induced barrier lowering makes $V_t = V_{t0} - \\eta V_{DS}$ with $\\eta = %s$. The subthreshold swing is %s mV/decade. By what factor does the off-current of a device fall when the supply (and so $V_{DS}$ of the off device) drops from %s V to %s V? Round to 3 decimal places.", edakit.Num(eta, 2), edakit.Num(s, 0), edakit.Num(v1, 1), edakit.Num(v2, 1))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("The threshold rises by $\\eta \\Delta V = %s$ mV, so the current falls by $10^{%s/%s} = %s$.", edakit.Num(dv, 1), edakit.Num(dv, 1), edakit.Num(s, 0), ans),
		map[string]string{"dvt_mv": edakit.Num(dv, 1), "factor": ans})
}

// multiVt: mixed-flavour leakage totals, savings and the slack budget for swaps.
func multiVt(ctx generator.GeneratorContext) generator.Problem {
	r := float64(edakit.Pick([]int{5, 8, 10, 12, 15, 20}))
	pl := float64(edakit.Int(2, 20))
	switch kpOf(ctx) {
	case 1:
		nl := 1000 * edakit.Int(2, 40)
		nh := 1000 * edakit.Int(10, 200)
		tot := (float64(nl)*pl + float64(nh)*pl/r) / 1000
		ans := edakit.Num(tot, 2)
		q := edakit.Sprintf("A block has %d LVT cells leaking %s nW each and %d HVT cells. HVT cells leak %s times less than LVT cells. What is the block's leakage power in microwatts? Round to 2 decimal places.", nl, edakit.Num(pl, 0), nh, edakit.Num(r, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("LVT: %d x %s nW; HVT: %d x %s nW. The sum is %s uW.", nl, edakit.Num(pl, 0), nh, edakit.Num(pl/r, 3), ans),
			map[string]string{"leak_uw": ans})
	case 2:
		n := 1000 * edakit.Int(20, 100)
		m := n * edakit.Int(3, 9) / 10
		before := float64(n) * pl
		after := float64(n-m)*pl + float64(m)*pl/r
		saved := 100 * (before - after) / before
		ans := edakit.Num(saved, 1)
		q := edakit.Sprintf("A block is built entirely from %d LVT cells. A leakage-recovery pass swaps %d cells with positive slack to HVT, which leaks %s times less. What percentage of the leakage is saved? Round to 1 decimal place.", n, m, edakit.Num(r, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Each swapped cell saves $1 - 1/%s$ of its leakage. The saving is $(%d/%d)(1 - 1/%s)$, which is %s%%.", edakit.Num(r, 0), m, n, edakit.Num(r, 0), ans),
			map[string]string{"saved_pct": ans})
	}
	k := edakit.Scaled(ctx, 6, 20)
	d := float64(edakit.Int(10, 30))
	x := float64(edakit.Pick([]int{10, 15, 20, 25, 30}))
	slack := float64(edakit.Int(5, int(float64(k)*d*x/100)))
	per := d * x / 100
	m := floor(slack / per)
	if m > k {
		m = k
	}
	q := edakit.Sprintf("A path has %d LVT gates of %s ps each and %s ps of setup slack. Swapping a gate to HVT makes it %s%% slower. How many gates on this path can be swapped to HVT without making the slack negative?", k, edakit.Num(d, 0), edakit.Num(slack, 0), edakit.Num(x, 0))
	return edakit.Numeric(q, edakit.IntStr(m),
		edakit.Sprintf("Each swap costs %s ps. $\\lfloor %s / %s \\rfloor = %d$ swaps fit in the slack (and at most %d gates exist).", edakit.Num(per, 2), edakit.Num(slack, 0), edakit.Num(per, 2), m, k),
		map[string]string{"cost_ps": edakit.Num(per, 2), "swaps": edakit.IntStr(m)})
}

// energyPerOp: switching energy, total energy with leakage, and the leakage crossover.
func energyPerOp(ctx generator.GeneratorContext) generator.Problem {
	alpha := edakit.Pick([]float64{0.05, 0.1, 0.15, 0.2, 0.25})
	c := float64(edakit.Scaled(ctx, 20, 800))
	v := edakit.Pick([]float64{0.5, 0.6, 0.7, 0.8, 0.9, 1.0})
	edyn := alpha * c * v * v
	pre := edakit.Sprintf("A datapath switches %s pF of capacitance with activity $\\alpha = %s$ (0-to-1 probability per operation) at $V_{DD} = %s$ V, one operation per cycle. ", edakit.Num(c, 0), edakit.Num(alpha, 2), edakit.Num(v, 2))
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(edyn, 3)
		return edakit.Numeric(pre+"What is the dynamic energy per operation, in picojoules? Round to 3 decimal places.", ans,
			edakit.Sprintf("$E = \\alpha C V^2 = %s \\times %s \\times %s = %s$ pJ (pF times volts squared is pJ).", edakit.Num(alpha, 2), edakit.Num(c, 0), edakit.Num(v*v, 4), ans),
			map[string]string{"edyn_pj": ans})
	case 2:
		pl := float64(edakit.Int(5, 200))
		f := float64(edakit.Pick([]int{50, 100, 200, 250, 500, 1000}))
		el := pl / f
		ans := edakit.Num(edyn+el, 3)
		q := pre + edakit.Sprintf("It also leaks %s uW and runs at %s MHz. What is the total energy per operation, in picojoules? Round to 3 decimal places.", edakit.Num(pl, 0), edakit.Num(f, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Dynamic energy is %s pJ. Leakage energy is power times cycle time: %s uW / %s MHz = %s pJ. Total %s pJ.", edakit.Num(edyn, 4), edakit.Num(pl, 0), edakit.Num(f, 0), edakit.Num(el, 4), ans),
			map[string]string{"edyn_pj": edakit.Num(edyn, 4), "eleak_pj": edakit.Num(el, 4)})
	}
	pl := float64(edakit.Int(5, 200))
	f := pl / edyn
	ans := edakit.Num(f, 2)
	q := pre + edakit.Sprintf("It leaks %s uW regardless of clock rate. At what clock frequency, in MHz, does leakage energy per operation equal dynamic energy per operation? Round to 2 decimal places.", edakit.Num(pl, 0))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Dynamic energy is %s pJ per operation. Leakage energy per operation is $P_{leak}/f$, so they match at $f = %s / %s = %s$ MHz. Slower than that, leakage dominates.", edakit.Num(edyn, 4), edakit.Num(pl, 0), edakit.Num(edyn, 4), ans),
		map[string]string{"edyn_pj": edakit.Num(edyn, 4), "f_mhz": ans})
}
