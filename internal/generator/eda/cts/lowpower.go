package cts

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// multiVoltage: V^2 scaling, DVFS savings, alpha-power delay scaling.
func multiVoltage(ctx generator.GeneratorContext) generator.Problem {
	v1 := pick(0.8, 1.1, 0.05)
	v2 := v1 - pick(0.1, 0.3, 0.05)
	switch kpOf(ctx) {
	case 1:
		p1 := float64(5 * edakit.Int(4, 60))
		p2 := p1 * (v2 / v1) * (v2 / v1)
		return edakit.Numeric(edakit.Sprintf("A block burns %s mW of dynamic power at %s V. It is moved to its own voltage island at %s V and keeps the same clock frequency and activity. What is its new dynamic power, in mW? Round to 2 decimal places.", f(p1, 0), f(v1, 2), f(v2, 2)), f(p2, 2),
			edakit.Sprintf("Dynamic power scales as $V^2$: $%s \\cdot (%s/%s)^2 = %s$ mW.", f(p1, 0), f(v2, 2), f(v1, 2), f(p2, 2)),
			map[string]string{"p2": f(p2, 3)})
	case 2:
		f1 := pick(1, 3, 0.25)
		f2 := f1 * pick(0.5, 0.9, 0.05)
		s := 100 * (1 - (v2/v1)*(v2/v1)*(f2/f1))
		return edakit.Numeric(edakit.Sprintf("Dynamic voltage and frequency scaling moves a domain from %s V at %s GHz to %s V at %s GHz. By what percentage does its dynamic power fall? Round to 1 decimal place.", f(v1, 2), f(f1, 2), f(v2, 2), f(edakit.Round(f2, 3), 3)), f(s, 1),
			edakit.Sprintf("$P \\propto V^2 f$, so the new power is $(%s/%s)^2 (%s/%s) = %s$ of the old: a %s%% reduction.", f(v2, 2), f(v1, 2), f(edakit.Round(f2, 3), 3), f(f1, 2), f(1-s/100, 4), f(s, 1)),
			map[string]string{"saving": f(s, 2)})
	}
	vt := pick(0.25, 0.4, 0.05)
	a := 1.3
	r := (v2 / math.Pow(v2-vt, a)) / (v1 / math.Pow(v1-vt, a))
	return edakit.Numeric(edakit.Sprintf("Gate delay follows the alpha-power law, $t_d \\propto V/(V - V_t)^{\\alpha}$, with $V_t = %s$ V and $\\alpha = 1.3$. A voltage island is lowered from %s V to %s V. By what factor does gate delay grow? Round to 3 decimal places.", f(vt, 2), f(v1, 2), f(v2, 2)), f(r, 3),
		edakit.Sprintf("Ratio $= \\frac{%s/(%s)^{1.3}}{%s/(%s)^{1.3}} = %s$. The island saves power quadratically but must run slower, which is why the lower-voltage domain gets a lower frequency or only off-critical logic.", f(v2, 2), f(v2-vt, 2), f(v1, 2), f(v1-vt, 2), f(r, 3)),
		map[string]string{"ratio": f(r, 4)})
}

// powerGating: wake energy, break-even sleep time, net energy saved.
func powerGating(ctx generator.GeneratorContext) generator.Problem {
	c := pick(0.5, 10, 0.5)
	v := pick(0.6, 1.0, 0.05)
	pl := float64(10 * edakit.Int(5, 100))
	ps := pl * pick(0.02, 0.1, 0.01)
	ew := c * v * v
	switch kpOf(ctx) {
	case 1:
		return edakit.Numeric(edakit.Sprintf("A power-gated domain has %s nF of capacitance on its virtual supply rail, which discharges to 0 during sleep. On wake-up the header switches recharge it to %s V. How much energy is drawn from the supply to do so, in nJ? Round to 3 decimal places.", f(c, 1), f(v, 2)), f(ew, 3),
			edakit.Sprintf("Charging $C$ to $V$ through a resistance draws $QV = CV^2 = %s \\cdot %s^2 = %s$ nJ from the supply: half ends up stored and half is dissipated in the switches.", f(c, 1), f(v, 2), f(ew, 3)),
			map[string]string{"e_nj": f(ew, 4)})
	case 2:
		t := 1000 * ew / (pl - ps)
		return edakit.Numeric(edakit.Sprintf("A domain leaks %s µW when powered. Gated, only the switches leak, %s µW. Each wake-up costs %s nJ (recharging %s nF to %s V). What is the shortest sleep period that saves energy, in µs? Round to 2 decimal places.", f(pl, 0), f(ps, 1), f(ew, 3), f(c, 1), f(v, 2)), f(t, 2),
			edakit.Sprintf("Sleep saves $(P_{leak} - P_{switch}) t = %s\\,t$ µW. Break even when that equals %s nJ: $t = %s / %s$ nJ/µW $= %s$ µs (nJ over µW is ms, times 1000).", f(pl-ps, 1), f(ew, 3), f(ew, 3), f(pl-ps, 1), f(t, 2)),
			map[string]string{"t_us": f(t, 3)})
	}
	t := float64(10 * edakit.Int(1, 100))
	e := (pl-ps)*t/1000 - ew
	return edakit.Numeric(edakit.Sprintf("A domain leaks %s µW when powered and %s µW through its switches when gated. Each wake-up costs %s nJ. It sleeps for %s µs and then wakes. What net energy does gating save over that interval, in nJ? Round to 3 decimal places (negative if gating loses).", f(pl, 0), f(ps, 1), f(ew, 3), f(t, 0)), f(e, 3),
		edakit.Sprintf("Leakage avoided: $%s \\cdot %s/1000 = %s$ nJ. Minus the %s nJ wake-up cost gives %s nJ.", f(pl-ps, 1), f(t, 0), f((pl-ps)*t/1000, 3), f(ew, 3), f(e, 3)),
		map[string]string{"e": f(e, 4)})
}

// switchSizing: drop across parallel switches, count for an IR budget, off-state leakage.
func switchSizing(ctx generator.GeneratorContext) generator.Problem {
	ron := float64(100 * edakit.Int(5, 30))
	i := float64(10 * edakit.Int(2, 30))
	switch kpOf(ctx) {
	case 1:
		n := 100 * edakit.Int(10, 100)
		d := i * ron / float64(n)
		return edakit.Numeric(edakit.Sprintf("A power-gated domain draws %s mA through %d header switches in parallel, each with an on-resistance of $%s\\ \\Omega$. What is the voltage drop across the switch network, in mV? Round to 2 decimal places.", f(i, 0), n, f(ron, 0)), f(d, 2),
			edakit.Sprintf("Parallel switches give $R_{on}/N = %s/%d = %s$ ohm, so the drop is $%s \\cdot %s = %s$ mV. It adds directly to the domain's IR drop.", f(ron, 0), n, f(ron/float64(n), 3), f(i, 0), f(ron/float64(n), 3), f(d, 2)),
			map[string]string{"drop": f(d, 3)})
	case 2:
		b := float64(5 * edakit.Int(2, 10))
		n := int(math.Ceil(i*ron/b - 1e-9))
		return edakit.Numeric(edakit.Sprintf("A power-gated domain draws a peak of %s mA. Each header switch has an on-resistance of $%s\\ \\Omega$, and the switch network may drop at most %s mV. How many switches are needed?", f(i, 0), f(ron, 0), f(b, 0)), edakit.IntStr(n),
			edakit.Sprintf("Need $I R_{on}/N \\le %s$: $N \\ge %s \\cdot %s/%s = %s$, so %d switches.", f(b, 0), f(i, 0), f(ron, 0), f(b, 0), f(i*ron/b, 2), n),
			map[string]string{"n": edakit.IntStr(n)})
	}
	n := 10 * edakit.Int(20, 200)
	ioff := pick(0.5, 5, 0.5)
	v := pick(0.7, 1.0, 0.05)
	p := float64(n) * ioff * v / 1000
	return edakit.Numeric(edakit.Sprintf("A gated domain uses %d header switches. When off, each leaks %s nA from a %s V supply. What leakage power remains during sleep, in µW? Round to 2 decimal places.", n, f(ioff, 1), f(v, 2)), f(p, 2),
		edakit.Sprintf("$N I_{off} V = %d \\cdot %s$ nA $\\cdot %s$ V $= %s$ µW. More switches lower the on-state drop but raise this floor, so switch count is a trade-off.", n, f(ioff, 1), f(v, 2), f(p, 2)),
		map[string]string{"p_uw": f(p, 3)})
}

// rushCurrent: average in-rush, minimum ramp time, staged switch groups.
func rushCurrent(ctx generator.GeneratorContext) generator.Problem {
	c := pick(0.5, 10, 0.5)
	v := pick(0.6, 1.0, 0.05)
	q := c * v // nC
	switch kpOf(ctx) {
	case 1:
		t := float64(5 * edakit.Int(1, 40))
		i := 1000 * q / t
		return edakit.Numeric(edakit.Sprintf("Waking a power-gated domain recharges %s nF of virtual-rail capacitance from 0 to %s V in %s ns. What is the average rush current, in mA? Round to 1 decimal place.", f(c, 1), f(v, 2), f(t, 0)), f(i, 1),
			edakit.Sprintf("$I = C\\,\\Delta V/\\Delta t = %s \\cdot %s/%s$ nC/ns $= %s$ A $= %s$ mA. Drawn from the shared grid, this current droops the neighbouring always-on domains.", f(c, 1), f(v, 2), f(t, 0), f(q/t, 4), f(i, 1)),
			map[string]string{"i_ma": f(i, 2)})
	case 2:
		im := float64(10 * edakit.Int(5, 50))
		t := 1000 * q / im
		return edakit.Numeric(edakit.Sprintf("A power-gated domain has %s nF on its virtual rail, recharged to %s V at wake-up. To protect neighbouring domains, the rush current must stay below %s mA on average. What is the shortest allowed wake-up ramp, in ns? Round to 2 decimal places.", f(c, 1), f(v, 2), f(im, 0)), f(t, 2),
			edakit.Sprintf("$t \\ge CV/I_{max} = %s$ nC $/ %s$ mA $= %s$ ns.", f(q, 3), f(im, 0), f(t, 2)),
			map[string]string{"t_ns": f(t, 3)})
	}
	im := float64(10 * edakit.Int(5, 50))
	td := pick(0.5, 5, 0.5)
	k := int(math.Ceil(1000*q/(im*td) - 1e-9))
	return edakit.Numeric(edakit.Sprintf("A domain has %s nF on its virtual rail, recharged to %s V at wake-up. Its switches are daisy-chained in groups turned on %s ns apart, so the rail charges over roughly (number of groups) times %s ns. The average rush current must stay below %s mA. What is the smallest number of groups?", f(c, 1), f(v, 2), f(td, 1), f(td, 1), f(im, 0)), edakit.IntStr(k),
		edakit.Sprintf("Need $CV/(k\\,t_d) \\le I_{max}$: $k \\ge %s$ nC$/(%s \\cdot %s)$ = %s, so %d groups. A common refinement turns on a few weak switches first and the strong ones once the rail is nearly charged.", f(q, 3), f(im/1000, 3), f(td, 1), f(1000*q/(im*td), 3), k),
		map[string]string{"k": edakit.IntStr(k)})
}

var domainNames = []string{"PD_CPU", "PD_GPU", "PD_DSP", "PD_MODEM", "PD_ISP", "PD_NPU"}

// crossingCells: which protection cells a signal crossing between two domains needs.
func crossingCells(ctx generator.GeneratorContext) generator.Problem {
	vs := []float64{0.65, 0.75, 0.8, 0.9, 1.0}
	names := edakit.Shuffle(domainNames)
	src, dst := names[0], names[1]
	v1, v2 := edakit.Pick(vs), edakit.Pick(vs)
	if edakit.Int(0, 2) == 0 {
		v2 = v1
	}
	srcSw := edakit.Int(0, 2) > 0
	dstSw := edakit.Int(0, 1) == 1
	follows := dstSw && edakit.Int(0, 1) == 1
	text := edakit.Sprintf("A signal leaves %s (supply %s V, %s) and drives logic in %s (supply %s V, %s).", src, f(v1, 2), map[bool]string{true: "can be powered down", false: "always on"}[srcSw], dst, f(v2, 2), map[bool]string{true: "can be powered down", false: "always on"}[dstSw])
	if dstSw && srcSw {
		if follows {
			text += edakit.Sprintf(" The power state table never allows %s to be on while %s is off.", dst, src)
		} else {
			text += " The power state table allows either domain to be off while the other is on."
		}
	}
	iso := srcSw && (!dstSw || !follows)
	ls := v1 != v2
	opts := []string{"An isolation cell only", "A level shifter only", "Both isolation and level shifting", "Neither"}
	var ans string
	switch {
	case iso && ls:
		ans = opts[2]
	case iso:
		ans = opts[0]
	case ls:
		ans = opts[1]
	default:
		ans = opts[3]
	}
	var ds []string
	for _, o := range opts {
		if o != ans {
			ds = append(ds, o)
		}
	}
	why := ""
	if iso {
		why += "The receiving logic can be on while the driver is off, so the crossing needs isolation to hold a defined value instead of a floating input. "
	} else {
		why += "The receiver is never on while the driver is off, so no floating input can reach live logic and isolation is not required. "
	}
	if ls {
		why += edakit.Sprintf("The supplies differ (%s V and %s V), so the signal must be level shifted.", f(v1, 2), f(v2, 2))
	} else {
		why += "Both supplies are equal, so no level shifter is needed."
	}
	return edakit.Choice(text+" What does the crossing need?", ans, why, ds...)
}

// clampValue: the safe isolation clamp value for the receiving pin's polarity.
func clampValue(ctx generator.GeneratorContext) generator.Problem {
	type sig struct {
		name, role string
		clamp      int
	}
	sigs := []sig{
		{"rst_n", "an active-low reset of a register bank", 1},
		{"req", "an active-high request input of a bus arbiter", 0},
		{"cs_n", "an active-low chip select of a memory", 1},
		{"irq", "an active-high interrupt input of the interrupt controller", 0},
		{"cg_en", "the enable of a clock gate", 0},
		{"we_n", "an active-low write enable of a register file", 1},
		{"valid", "an active-high valid flag of a FIFO write port", 0},
		{"pwr_ok_n", "an active-low fault input of the power controller", 1},
	}
	s := edakit.Pick(sigs)
	names := edakit.Shuffle(domainNames)
	ans := map[int]string{0: "Clamp to 0", 1: "Clamp to 1"}[s.clamp]
	other := map[int]string{0: "Clamp to 1", 1: "Clamp to 0"}[s.clamp]
	q := edakit.Sprintf("The output %s of switchable domain %s drives %s in always-on domain %s. While %s is off, what should its isolation cell do?", s.name, names[0], s.role, names[1], names[0])
	return edakit.Choice(q, ans,
		edakit.Sprintf("Isolation must hold the receiving pin at its inactive value. For %s that value is %d, so the cell clamps to %d. Holding the last value needs a latch-based cell and can freeze an asserted control; leaving the input floating is exactly what isolation exists to prevent.", s.role, s.clamp, s.clamp),
		other, "Hold the last value before power-down", "Leave it undriven; the always-on logic ignores it")
}

// shifterDirection: which level shifter a voltage crossing needs.
func shifterDirection(ctx generator.GeneratorContext) generator.Problem {
	vs := []float64{0.6, 0.7, 0.75, 0.8, 0.9, 1.0, 1.1}
	pair := edakit.Shuffle(vs)[:2] // two distinct voltages, no redraw needed
	v1, v2 := pair[0], pair[1]
	names := edakit.Shuffle(domainNames)
	up := "A low-to-high shifter, which needs both the source and destination supplies"
	down := "A high-to-low shifter, which can be built from a buffer on the destination supply"
	q := edakit.Sprintf("A signal crosses from %s at %s V into %s at %s V, both always on. Which level shifter does the crossing need?", names[0], f(v1, 2), names[1], f(v2, 2))
	if v1 < v2 {
		return edakit.Choice(q, up,
			edakit.Sprintf("A %s V high level cannot fully turn off the PMOS of a gate supplied at %s V, which causes static current and a degraded output. A low-to-high shifter (usually a cross-coupled structure) needs both supplies to restore full swing.", f(v1, 2), f(v2, 2)),
			down, "No shifter: the higher supply always reads the lower swing correctly", "An isolation cell with a clamp value of 1")
	}
	return edakit.Choice(q, down,
		edakit.Sprintf("Going from %s V down to %s V, the input overdrives the receiving gate rather than underdriving it, so a simple buffer on the destination supply restores the lower swing. Flows still insert an explicit high-to-low shifter cell so that the crossing is visible to verification.", f(v1, 2), f(v2, 2)),
		up, "No shifter and no cell: the lower supply clips the signal automatically", "An isolation cell with a clamp value of 0")
}
