package sta

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// wireRC asks for the Elmore delay of a distributed RC wire, alone or with driver and load.
func wireRC(ctx generator.GeneratorContext) generator.Problem {
	r := edakit.Pick([]float64{0.5, 1, 1.5, 2, 3})
	c := edakit.Pick([]float64{0.1, 0.15, 0.2, 0.25})
	l := 50 * float64(edakit.Scaled(ctx, 2, 20))
	rd := 50 * float64(edakit.Int(2, 20))
	cl := float64(edakit.Int(2, 20))
	wire := func(L float64) float64 { return r * c * L * L / 2 }
	full := func(L float64) float64 { return rd*(c*L+cl) + r*L*(c*L/2+cl) }
	switch ctx.KPOrdinal {
	case 1:
		fs := wire(l)
		ans := edakit.Num(fs/1000, 2)
		q := edakit.Sprintf("A wire of length %s um has resistance %s ohm/um and capacitance %s fF/um, distributed uniformly. Ignoring the driver and any load, what is its Elmore delay from one end to the other, in ps? Round to 2 decimal places.",
			edakit.Num(l, 0), edakit.Num(r, 2), edakit.Num(c, 2))
		return edakit.Numeric(q, ans, edakit.Sprintf("A distributed line has Elmore delay $rcL^2/2$, half the lumped value: %s x %s x %s^2 / 2 = %s ohm-fF, which is %s ps.", edakit.Num(r, 2), edakit.Num(c, 2), edakit.Num(l, 0), edakit.Num(fs, 1), ans),
			map[string]string{"wire_fs": edakit.Num(fs, 1)})
	case 3:
		l2 := 2 * l
		fs := full(l2)
		ans := edakit.Num(fs/1000, 2)
		q := edakit.Sprintf("A driver of resistance %s ohm drives a wire of resistance %s ohm/um and capacitance %s fF/um into a load of %s fF. The wire was %s um long, and a reroute doubles it to %s um. What is the new Elmore delay from driver to load, in ps? Round to 2 decimal places.",
			edakit.Num(rd, 0), edakit.Num(r, 2), edakit.Num(c, 2), edakit.Num(cl, 0), edakit.Num(l, 0), edakit.Num(l2, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("Use $R_d(cL + C_L) + rL(cL/2 + C_L)$ with L = %s um: the delay goes from %s ps to %s ps. The driver term doubles in its wire part, while the wire term grows with the square of the length.", edakit.Num(l2, 0), edakit.Num(full(l)/1000, 2), ans),
			map[string]string{"delay_fs": edakit.Num(fs, 1)})
	}
	fs := full(l)
	ans := edakit.Num(fs/1000, 2)
	q := edakit.Sprintf("A driver of resistance %s ohm drives a wire %s um long with resistance %s ohm/um and capacitance %s fF/um, ending in a load of %s fF. What is the Elmore delay from the driver to the load, in ps? Round to 2 decimal places.",
		edakit.Num(rd, 0), edakit.Num(l, 0), edakit.Num(r, 2), edakit.Num(c, 2), edakit.Num(cl, 0))
	return edakit.Numeric(q, ans, edakit.Sprintf("The driver charges all of it, $R_d(cL + C_L)$ = %s ohm-fF. The wire resistance charges half its own capacitance plus the load, $rL(cL/2 + C_L)$ = %s ohm-fF. The total is %s ohm-fF, which is %s ps.",
		edakit.Num(rd*(c*l+cl), 1), edakit.Num(r*l*(c*l/2+cl), 1), edakit.Num(fs, 1), ans), map[string]string{"delay_fs": edakit.Num(fs, 1)})
}

// nldm draws a 2x2 window of a Liberty delay table and a query point (on a grid line, inside,
// or beyond the edge) and asks for the interpolated value.
func nldm(ctx generator.GeneratorContext) generator.Problem {
	slews := []int{10, 20, 40, 80, 160, 320}
	loads := []int{1, 2, 4, 8, 16, 32}
	i := edakit.Int(0, len(slews)-2)
	j := edakit.Int(0, len(loads)-2)
	s1, s2 := slews[i], slews[i+1]
	c1, c2 := loads[j], loads[j+1]
	v11 := float64(edakit.Int(15, 60))
	v12 := v11 + float64(edakit.Int(6, 30))
	v21 := v11 + float64(edakit.Int(3, 20))
	v22 := v12 + (v21 - v11) + float64(edakit.Int(0, 10))
	var s, c float64
	mode := ctx.KPOrdinal
	if mode == 0 {
		mode = edakit.Int(1, 3)
	}
	switch mode {
	case 1:
		s = float64(edakit.Pick([]int{s1, s2}))
		c = float64(c1) + float64(c2-c1)*float64(edakit.Int(1, 7))/8
	case 2:
		s = float64(s1) + float64(s2-s1)*float64(edakit.Int(1, 7))/8
		c = float64(c1) + float64(c2-c1)*float64(edakit.Int(1, 7))/8
	default:
		s = float64(s1) + float64(s2-s1)*float64(edakit.Int(1, 7))/8
		c = float64(c2) + float64(c2-c1)*float64(edakit.Int(1, 4))/4
	}
	fs1, fs2, fc1, fc2 := float64(s1), float64(s2), float64(c1), float64(c2)
	v := (v11*(fs2-s)*(fc2-c) + v12*(fs2-s)*(c-fc1) + v21*(s-fs1)*(fc2-c) + v22*(s-fs1)*(c-fc1)) / ((fs2 - fs1) * (fc2 - fc1))
	ans := edakit.Num(v, 2)
	q := edakit.Sprintf("A Liberty cell_rise table is indexed by input slew (index_1) and output load (index_2). Its window around the query holds: at slew %d ps, load %d fF the delay is %s ps; at slew %d ps, load %d fF it is %s ps; at slew %d ps, load %d fF it is %s ps; at slew %d ps, load %d fF it is %s ps. ",
		s1, c1, edakit.Num(v11, 0), s1, c2, edakit.Num(v12, 0), s2, c1, edakit.Num(v21, 0), s2, c2, edakit.Num(v22, 0))
	if mode == 3 {
		q += edakit.Sprintf("Load %d fF is the last index_2 point of the table. ", c2)
	}
	q += edakit.Sprintf("Using bilinear interpolation (extrapolating linearly from this window if the point lies outside it), what delay does the timer compute at slew %s ps and load %s fF? Round to 2 decimal places.", edakit.Num(s, 2), edakit.Num(c, 2))
	ts := (s - fs1) / (fs2 - fs1)
	tc := (c - fc1) / (fc2 - fc1)
	exp := edakit.Sprintf("The fractional positions are $t_s = %s$ along slew and $t_c = %s$ along load. Interpolate along load at each slew, then along slew: the result is %s ps.", edakit.Num(ts, 4), edakit.Num(tc, 4), ans)
	if mode == 3 {
		exp += " With $t_c > 1$ the same formula extrapolates past the table edge, which is only as good as the linear trend of the last window."
	}
	return edakit.Numeric(q, ans, exp, map[string]string{"t_slew": edakit.Num(ts, 4), "t_load": edakit.Num(tc, 4), "delay": ans})
}

// stageDelay times one cell, one stage, or a short path of linear-model cells and pi wires.
func stageDelay(ctx generator.GeneratorContext) generator.Problem {
	type stage struct{ d0, r, cw, rw, cpin float64 }
	n := 1
	if ctx.KPOrdinal == 3 || (ctx.KPOrdinal == 0 && ctx.Difficulty > 0.6) {
		n = edakit.Int(2, 3)
	}
	st := make([]stage, n)
	for i := range st {
		st[i] = stage{
			d0:   float64(edakit.Int(5, 20)),
			r:    edakit.Pick([]float64{0.5, 1, 1.5, 2, 3}),
			cw:   float64(edakit.Int(2, 10)),
			rw:   edakit.Pick([]float64{0.1, 0.2, 0.3, 0.5}),
			cpin: float64(edakit.Int(1, 6)),
		}
	}
	cell := func(s stage) float64 { return s.d0 + s.r*(s.cw+s.cpin) }
	wire := func(s stage) float64 { return s.rw * (s.cw/2 + s.cpin) }
	desc := ""
	for i, s := range st {
		desc += edakit.Sprintf(" Stage %d: the cell has intrinsic delay %s ps and drive resistance %s kohm; its output wire has capacitance %s fF and resistance %s kohm; the sink pin it drives has capacitance %s fF.",
			i+1, edakit.Num(s.d0, 0), edakit.Num(s.r, 2), edakit.Num(s.cw, 0), edakit.Num(s.rw, 2), edakit.Num(s.cpin, 0))
	}
	model := "Model each cell delay as intrinsic delay plus drive resistance times its total load (wire plus pin), and each wire as a pi model whose Elmore delay is wire resistance times (half the wire capacitance plus the pin capacitance). One kohm times one fF is one ps."
	if ctx.KPOrdinal == 1 {
		v := cell(st[0])
		ans := edakit.Num(v, 2)
		return edakit.Numeric(model+desc+" What is the cell delay of stage 1, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("Cell delay = %s + %s x (%s + %s) = %s ps.", edakit.Num(st[0].d0, 0), edakit.Num(st[0].r, 2), edakit.Num(st[0].cw, 0), edakit.Num(st[0].cpin, 0), ans), map[string]string{"cell": ans})
	}
	if n == 1 {
		v := cell(st[0]) + wire(st[0])
		ans := edakit.Num(v, 2)
		return edakit.Numeric(model+desc+" What is the total stage delay, from the cell input to the sink pin, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("Cell arc %s ps plus net arc %s ps gives %s ps.", edakit.Num(cell(st[0]), 2), edakit.Num(wire(st[0]), 2), ans), map[string]string{"stage": ans})
	}
	start := float64(5 * edakit.Int(0, 10))
	total := start
	for _, s := range st {
		total += cell(s) + wire(s)
	}
	ans := edakit.Num(total, 2)
	return edakit.Numeric(model+desc+edakit.Sprintf(" The signal reaches the stage 1 cell input at %s ps, and each stage's sink pin is the next stage's cell input. What is the arrival time at the last sink pin, in ps? Round to 2 decimal places.", edakit.Num(start, 0)), ans,
		edakit.Sprintf("Sum cell and wire delays stage by stage on top of the %s ps start: the arrival is %s ps.", edakit.Num(start, 0), ans), map[string]string{"arrival": ans})
}

// slewProp contrasts worst-slew merging with the path's own slew at a two-input gate.
func slewProp(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Pick([]float64{0.2, 0.3, 0.4, 0.5})
	m := edakit.Pick([]float64{0.5, 0.6, 0.8})
	kh := edakit.Pick([]float64{0.2, 0.3, 0.4})
	var atA, atB, sA, sB, aA, aB, bA, bB, aH float64
	for {
		atA, atB = float64(5*edakit.Int(10, 40)), float64(5*edakit.Int(10, 40))
		sA, sB = float64(10*edakit.Int(2, 20)), float64(10*edakit.Int(2, 20))
		aA, aB = float64(edakit.Int(15, 40)), float64(edakit.Int(15, 40))
		bA, bB = float64(edakit.Int(10, 30)), float64(edakit.Int(10, 30))
		aH = float64(edakit.Int(15, 40))
		arrA, arrB := atA+aA+k*sA, atB+aB+k*sB
		soA, soB := bA+m*sA, bB+m*sB
		if arrA != arrB && soA != soB && (arrA > arrB) != (soA > soB) {
			break
		}
	}
	arrA, arrB := atA+aA+k*sA, atB+aB+k*sB
	soA, soB := bA+m*sA, bB+m*sB
	soMax := math.Max(soA, soB)
	arrG := math.Max(arrA, arrB)
	crit := "A"
	soCrit := soA
	if arrB > arrA {
		crit, soCrit = "B", soB
	}
	gba := arrG + aH + kh*soMax
	pba := arrG + aH + kh*soCrit
	q := edakit.Sprintf("Gate G has inputs A and B. Input A arrives at %s ps with slew %s ps; input B arrives at %s ps with slew %s ps. The arc from input X has delay $a_X + %s s_X$ and output slew $b_X + %s s_X$, where $s_X$ is the input slew, with $a_A = %s$, $a_B = %s$, $b_A = %s$, $b_B = %s$ (all ps). G drives gate H, whose arc delay is $%s + %s s$ for input slew s. ",
		edakit.Num(atA, 0), edakit.Num(sA, 0), edakit.Num(atB, 0), edakit.Num(sB, 0), edakit.Num(k, 1), edakit.Num(m, 1), edakit.Num(aA, 0), edakit.Num(aB, 0), edakit.Num(bA, 0), edakit.Num(bB, 0), edakit.Num(aH, 0), edakit.Num(kh, 1))
	facts := map[string]string{"slew_out": edakit.Num(soMax, 2), "gba": edakit.Num(gba, 2), "pba": edakit.Num(pba, 2)}
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(soMax, 2)
		return edakit.Numeric(q+"A graph-based timer keeps one slew per pin, the worst over all arcs into it. What slew does it store at G's output, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("The arc from A gives output slew %s ps and the arc from B gives %s ps. Worst-slew merging keeps the larger, %s ps, even though it may not belong to the latest-arriving input.", edakit.Num(soA, 2), edakit.Num(soB, 2), ans), facts)
	case 3:
		ans := edakit.Num(gba-pba, 2)
		return edakit.Numeric(q+edakit.Sprintf("The critical path comes through input %s, the latest arrival at G's output. Re-timing that path with its own slew (path-based analysis) instead of the merged worst slew, by how many ps does the arrival at H's output drop? Round to 2 decimal places.", crit), ans,
			edakit.Sprintf("GBA uses the worst slew %s ps at H: arrival %s ps. The path through %s carries slew %s ps: arrival %s ps. The pessimism removed is %s ps.", edakit.Num(soMax, 2), edakit.Num(gba, 2), crit, edakit.Num(soCrit, 2), edakit.Num(pba, 2), ans), facts)
	}
	ans := edakit.Num(gba, 2)
	return edakit.Numeric(q+"Under graph-based analysis (latest arrival and worst slew merged separately at G's output), what is the arrival time at H's output, in ps? Round to 2 decimal places.", ans,
		edakit.Sprintf("G's output arrival is max(%s, %s) = %s ps and its stored slew is max(%s, %s) = %s ps. H adds %s + %s x %s, giving %s ps.",
			edakit.Num(arrA, 2), edakit.Num(arrB, 2), edakit.Num(arrG, 2), edakit.Num(soA, 2), edakit.Num(soB, 2), edakit.Num(soMax, 2), edakit.Num(aH, 0), edakit.Num(kh, 1), edakit.Num(soMax, 2), ans), facts)
}

// ceff computes the charge-matched effective capacitance of a pi load under a ramp.
func ceff(ctx generator.GeneratorContext) generator.Problem {
	c1 := float64(edakit.Int(5, 30))
	c2 := float64(edakit.Int(10, 60))
	r := float64(50 * edakit.Int(4, 40))
	tau := r * c2 / 1000 // ps
	eff := func(T float64) float64 { return c1 + c2*(1-(tau/T)*(1-math.Exp(-T/tau))) }
	model := edakit.Sprintf("A driver sees a pi load: %s fF at its pin, then %s ohm of wire resistance, then %s fF at the far end. ", edakit.Num(c1, 0), edakit.Num(r, 0), edakit.Num(c2, 0))
	formula := "Matching the charge delivered during a ramp of duration T gives $C_{eff} = C_1 + C_2\\left[1 - \\frac{\\tau}{T}\\left(1 - e^{-T/\\tau}\\right)\\right]$ with $\\tau = R C_2$. "
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(tau, 2)
		return edakit.Numeric(model+"The resistance shields the far capacitance with time constant $\\tau = R C_2$. What is $\\tau$, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("%s ohm x %s fF = %s ohm-fF, which is %s ps. A ramp much shorter than this sees mostly the near capacitance.", edakit.Num(r, 0), edakit.Num(c2, 0), edakit.Num(r*c2, 0), ans), map[string]string{"tau": ans})
	case 3:
		t0 := float64(edakit.Int(5, 20))
		kd := edakit.Pick([]float64{0.5, 1, 1.5, 2})
		T0 := t0 + kd*(c1+c2)
		e1 := eff(T0)
		ans := edakit.Num(e1, 2)
		return edakit.Numeric(model+formula+edakit.Sprintf("The cell's output ramp time is $%s + %s C$ ps for a load of C fF. Start the iteration from the total capacitance: compute the ramp time it gives, then the new $C_{eff}$. What is that first updated $C_{eff}$, in fF? Round to 2 decimal places.", edakit.Num(t0, 0), edakit.Num(kd, 1)), ans,
			edakit.Sprintf("$C_{tot} = %s$ fF gives T = %s ps. With $\\tau = %s$ ps the charge match gives $C_{eff}$ = %s fF. A timer repeats this until T and $C_{eff}$ agree.", edakit.Num(c1+c2, 0), edakit.Num(T0, 2), edakit.Num(tau, 2), ans), map[string]string{"t0": edakit.Num(T0, 2), "ceff": ans})
	}
	T := float64(edakit.Int(20, 200))
	e := eff(T)
	ans := edakit.Num(e, 2)
	return edakit.Numeric(model+formula+edakit.Sprintf("The driver's output ramp lasts %s ps. What is $C_{eff}$, in fF? Round to 2 decimal places.", edakit.Num(T, 0)), ans,
		edakit.Sprintf("$\\tau$ = %s ps and $\\tau/T$ = %s. The far capacitor contributes %s fF of its %s fF, so $C_{eff}$ = %s fF, less than the total %s fF.",
			edakit.Num(tau, 2), edakit.Num(tau/T, 4), edakit.Num(e-c1, 2), edakit.Num(c2, 0), ans, edakit.Num(c1+c2, 0)), map[string]string{"tau": edakit.Num(tau, 2), "ceff": ans})
}
