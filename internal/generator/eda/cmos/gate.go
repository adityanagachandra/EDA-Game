package cmos

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const netLegend = "Write S(...) for transistors in series and P(...) for transistors in parallel."

// staticCMOS asks for the pull-down network (KP 1), the pull-up network (KP 2) or the function
// of a given pull-down network (KP 3) of a random two-level complex gate.
func staticCMOS(ctx generator.GeneratorContext) generator.Problem {
	e := randCompound(ctx)
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	f := "$Y = \\overline{" + e.latex() + "}$"
	switch kp {
	case 1:
		q := "A static CMOS gate computes " + f + ". " + netLegend + " Which nMOS pull-down network implements it?"
		return edakit.Choice(q, e.pdn(),
			"The pull-down network conducts exactly when the un-inverted expression is true: AND becomes series and OR becomes parallel, giving "+e.pdn()+".",
			e.pun(), e.net(topOnly, 0), e.net(innerOnly, 0))
	case 2:
		q := "A static CMOS gate computes " + f + ". " + netLegend + " Which pMOS pull-up network implements it?"
		return edakit.Choice(q, e.pun(),
			"The pull-up network is the dual of the pull-down network: every series connection becomes parallel and every parallel connection becomes series, at every level, giving "+e.pun()+".",
			e.pdn(), e.net(topOnly, 0), e.net(innerOnly, 0))
	}
	q := "The nMOS pull-down network of a static CMOS gate is " + e.pdn() + ". " + netLegend + " What function does the gate compute?"
	dual := e.swapped(always, 0)
	top := e.swapped(topOnly, 0)
	return edakit.Choice(q, f,
		"Series nMOS devices conduct when all their inputs are 1 (AND), parallel devices when any is 1 (OR). The network pulls the output low when "+e.latex()+" is true, so the output is its complement.",
		"$Y = \\overline{"+dual.latex()+"}$", "$Y = "+e.latex()+"$", "$Y = \\overline{"+top.latex()+"}$")
}

// transistorCount asks for the device count (KP 1) or the tallest series stack of the
// pull-down (KP 2) or pull-up (KP 3) network.
func transistorCount(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	var e *expr
	if kp == 1 && edakit.Int(0, 2) == 0 {
		op := edakit.Pick([]byte{'&', '|'})
		e = flat(op, edakit.Scaled(ctx, 2, 5))
	} else {
		e = randCompound(ctx)
	}
	f := "$Y = \\overline{" + e.latex() + "}$"
	n := e.literals()
	switch kp {
	case 1:
		ans := 2 * n
		return edakit.Numeric("How many transistors does a single-stage static CMOS gate for "+f+" use?",
			edakit.IntStr(ans),
			edakit.Sprintf("Every input literal drives one nMOS in the pull-down network and one pMOS in the pull-up network. With %d literals the gate has %d transistors.", n, ans),
			map[string]string{"literals": edakit.IntStr(n)})
	case 2:
		h := e.stack(false)
		return edakit.Numeric("For the static CMOS gate "+f+", how many nMOS transistors are in the tallest series stack of the pull-down network?",
			edakit.IntStr(h),
			edakit.Sprintf("In the pull-down network AND terms are series and OR terms parallel. A series group adds its children's stack heights and a parallel group takes the largest, giving %d.", h),
			map[string]string{"pdn_stack": edakit.IntStr(h)})
	}
	h := e.stack(true)
	return edakit.Numeric("For the static CMOS gate "+f+", how many pMOS transistors are in the tallest series stack of the pull-up network?",
		edakit.IntStr(h),
		edakit.Sprintf("The pull-up network is the dual, so OR terms become series and AND terms parallel. Adding along series groups and taking the maximum over parallel groups gives %d.", h),
		map[string]string{"pun_stack": edakit.IntStr(h)})
}

// aoiCell is an AOI or OAI cell named by its group sizes.
type aoiCell struct {
	oai    bool
	groups []int
}

func (c aoiCell) name() string {
	s := "AOI"
	if c.oai {
		s = "OAI"
	}
	for _, g := range c.groups {
		s += edakit.IntStr(g)
	}
	return s
}

func (c aoiCell) expr() *expr {
	var l letters
	top := &expr{op: '|'}
	inner := byte('&')
	if c.oai {
		top.op, inner = '&', '|'
	}
	for _, g := range c.groups {
		if g == 1 {
			top.kids = append(top.kids, &expr{op: 'L', lit: l.take()})
			continue
		}
		sub := &expr{op: inner}
		for i := 0; i < g; i++ {
			sub.kids = append(sub.kids, &expr{op: 'L', lit: l.take()})
		}
		top.kids = append(top.kids, sub)
	}
	return top
}

func (c aoiCell) inputs() int {
	n := 0
	for _, g := range c.groups {
		n += g
	}
	return n
}

var aoiShapes = [][]int{{2, 1}, {2, 2}, {2, 1, 1}, {2, 2, 1}, {3, 1}, {3, 2}, {2, 2, 2}, {3, 3}, {2, 2, 1, 1}, {3, 1, 1}}

// regroup returns a different group list with the same number of inputs, for a distractor.
func regroup(g []int) []int {
	n := 0
	for _, x := range g {
		n += x
	}
	for _, s := range edakit.Shuffle(aoiShapes) {
		t := 0
		for _, x := range s {
			t += x
		}
		if t == n && !sameInts(s, g) {
			return s
		}
	}
	return append([]int{n - 1}, 1)
}

func sameInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func aoiOai(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	c := aoiCell{oai: edakit.Int(0, 1) == 1, groups: edakit.Pick(aoiShapes)}
	other := aoiCell{oai: !c.oai, groups: c.groups}
	re := aoiCell{oai: c.oai, groups: regroup(c.groups)}
	f := func(x aoiCell) string { return "$Y = \\overline{" + x.expr().latex() + "}$" }
	plain := "$Y = " + c.expr().latex() + "$"
	switch kp {
	case 1:
		q := "A standard-cell library contains the cell " + c.name() + ", with inputs lettered A, B, C, ... in group order. Which function does it compute?"
		return edakit.Choice(q, f(c),
			"In AOI the groups are ANDed, the results ORed, and the output inverted; OAI swaps AND and OR. Each digit is the size of one group, so "+c.name()+" computes "+f(c)+".",
			f(other), plain, f(re))
	case 2:
		q := "Which compound cell computes " + f(c) + " in a single stage?"
		plainName := strings.Replace(c.name(), "I", "", 1)
		return edakit.Choice(q, c.name(),
			"Read the expression from the inside: the inner groups' operator gives the first letter (A for AND, O for OR), then the outer operator, then I for the inversion; the digits are the group sizes.",
			other.name(), re.name(), plainName)
	}
	n := c.inputs()
	ans := edakit.Sprintf("%d inputs, %d transistors", n, 2*n)
	q := "How many inputs and transistors does a single-stage static CMOS " + c.name() + " cell have?"
	return edakit.Choice(q, ans,
		edakit.Sprintf("The digits sum to %d inputs. Each input drives one nMOS and one pMOS, so the cell has %d transistors.", n, 2*n),
		edakit.Sprintf("%d inputs, %d transistors", n, n),
		edakit.Sprintf("%d inputs, %d transistors", n, 2*n+2),
		edakit.Sprintf("%d inputs, %d transistors", len(c.groups), 2*len(c.groups)))
}

// passVoltage computes the level a pass transistor or transmission gate passes.
func passVoltage(ctx generator.GeneratorContext) generator.Problem {
	vdd := edakit.Pick([]float64{0.8, 0.9, 1.0, 1.1, 1.2, 1.8})
	vt := float64(edakit.Int(5, 9)) * 0.05
	v := func(x float64) string { return edakit.Num(x, 2) + " V" }
	switch edakit.Int(0, 3) {
	case 0:
		q := edakit.Sprintf("An nMOS pass transistor with its gate at $V_{DD} = %s$ V passes a logic 1 from a node held at $V_{DD}$ into a capacitive load that starts at 0 V. Its threshold is $V_{tn} = %s$ V; ignore the body effect. What voltage does the load settle at?", edakit.Num(vdd, 2), edakit.Num(vt, 2))
		return edakit.Choice(q, v(vdd-vt),
			"The nMOS stops conducting once its gate-source voltage falls to the threshold, so the source (the load) rises only to VDD minus Vtn: a weak 1.",
			v(vdd), v(vt), v(vdd-2*vt))
	case 1:
		q := edakit.Sprintf("A pMOS pass transistor with its gate at 0 V passes a logic 0 into a load that starts at $V_{DD} = %s$ V. Its threshold magnitude is $|V_{tp}| = %s$ V; ignore the body effect. What voltage does the load settle at?", edakit.Num(vdd, 2), edakit.Num(vt, 2))
		return edakit.Choice(q, v(vt),
			"The pMOS turns off once its source falls to one threshold above its gate, so the load discharges only to |Vtp|: a weak 0.",
			v(0), v(vdd-vt), v(2*vt))
	case 2:
		q := edakit.Sprintf("A transmission gate (an nMOS and a pMOS in parallel, both on) passes a logic 1 from a node held at $V_{DD} = %s$ V into a load at 0 V. Each device has a threshold magnitude of %s V. What voltage does the load settle at?", edakit.Num(vdd, 2), edakit.Num(vt, 2))
		return edakit.Choice(q, v(vdd),
			"When the nMOS cuts off near VDD minus Vt, the pMOS is still fully on and pulls the load the rest of the way, so a transmission gate passes a full-swing 1.",
			v(vdd-vt), v(vdd-2*vt), v(vdd/2))
	}
	q := edakit.Sprintf("Three nMOS pass transistors are chained source to drain, each with its gate tied to $V_{DD} = %s$ V, passing a 1 from a node at $V_{DD}$. Each has $V_{tn} = %s$ V; ignore the body effect. What voltage does the far end settle at?", edakit.Num(vdd, 2), edakit.Num(vt, 2))
	return edakit.Choice(q, v(vdd-vt),
		"Every device in the chain has its gate at VDD, so each node can rise to VDD minus Vt and no further. The drops do not accumulate; they would only if one pass transistor's output drove the next one's gate.",
		v(vdd-3*vt), v(vdd), v(vdd-2*vt))
}

// inverterVM computes the switching threshold of a long-channel inverter.
func inverterVM(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	vdd := edakit.Pick([]float64{1.0, 1.2, 1.5, 1.8, 2.5})
	vtn := float64(edakit.Int(6, 10)) * 0.05
	vtp := float64(edakit.Int(6, 10)) * 0.05
	switch kp {
	case 1:
		ratio := edakit.Pick([]float64{0.25, 0.5, 1, 2, 4})
		r := sqrt(ratio)
		vm := (vtn + r*(vdd-vtp)) / (1 + r)
		ans := edakit.Num(vm, 3)
		q := edakit.Sprintf("A long-channel CMOS inverter obeys the square law with $V_{DD} = %s$ V, $V_{tn} = %s$ V, $|V_{tp}| = %s$ V and $k_p / k_n = %s$. At its switching threshold $V_M$ both devices are saturated and carry equal current. What is $V_M$ in volts? Round to 3 decimal places.",
			edakit.Num(vdd, 2), edakit.Num(vtn, 2), edakit.Num(vtp, 2), edakit.Num(ratio, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Equating the saturation currents gives $V_M = (V_{tn} + r(V_{DD} - |V_{tp}|))/(1 + r)$ with $r = \\sqrt{k_p/k_n} = %s$. That is %s V.", edakit.Num(r, 3), ans),
			map[string]string{"r": edakit.Num(r, 4), "vm": ans})
	case 2:
		mu := edakit.Pick([]float64{2, 2.5, 3})
		wr := edakit.Pick([]float64{1, 1.5, 2, 3, 4})
		ratio := wr / mu
		r := sqrt(ratio)
		vm := (vtn + r*(vdd-vtp)) / (1 + r)
		ans := edakit.Num(vm, 3)
		q := edakit.Sprintf("A long-channel inverter has $W_p / W_n = %s$ with equal lengths, and electron mobility is %s times hole mobility. With $V_{DD} = %s$ V, $V_{tn} = %s$ V and $|V_{tp}| = %s$ V, what is the switching threshold $V_M$ in volts? Round to 3 decimal places.",
			edakit.Num(wr, 2), edakit.Num(mu, 2), edakit.Num(vdd, 2), edakit.Num(vtn, 2), edakit.Num(vtp, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Here $k_p/k_n = (W_p/W_n)/(\\mu_n/\\mu_p) = %s$, so $r = %s$ and $V_M = (V_{tn} + r(V_{DD} - |V_{tp}|))/(1+r) = %s$ V.", edakit.Num(ratio, 4), edakit.Num(r, 4), ans),
			map[string]string{"kp_over_kn": edakit.Num(ratio, 4), "vm": ans})
	}
	mu := edakit.Pick([]float64{2, 2.5, 3})
	vm := vdd / 2
	r := (vm - vtn) / (vdd - vm - vtp)
	wr := r * r * mu
	ans := edakit.Num(wr, 2)
	q := edakit.Sprintf("You want a long-channel inverter to switch at exactly $V_M = V_{DD}/2$ with $V_{DD} = %s$ V, $V_{tn} = %s$ V and $|V_{tp}| = %s$ V. Electron mobility is %s times hole mobility and the lengths are equal. What width ratio $W_p / W_n$ is required? Round to 2 decimal places.",
		edakit.Num(vdd, 2), edakit.Num(vtn, 2), edakit.Num(vtp, 2), edakit.Num(mu, 2))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Equal saturation currents at $V_M$ need $r = \\sqrt{k_p/k_n} = (V_M - V_{tn})/(V_{DD} - V_M - |V_{tp}|) = %s$. Then $k_p/k_n = r^2 = %s$ and $W_p/W_n = r^2 \\mu_n/\\mu_p = %s$.", edakit.Num(r, 4), edakit.Num(r*r, 4), ans),
		map[string]string{"r": edakit.Num(r, 4), "wp_over_wn": ans})
}

// noiseMargins computes NMH and NML from VTC points, a piecewise-linear VTC, or a driver and
// receiver pair.
func noiseMargins(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	vdd := edakit.Pick([]float64{0.8, 1.0, 1.2, 1.8})
	mv := func(x float64) int { return int(x*1000 + 0.5) }
	switch kp {
	case 1:
		voh := vdd - float64(edakit.Int(0, 4))*0.01
		vol := float64(edakit.Int(0, 4)) * 0.01
		vil := vdd * float64(edakit.Int(35, 44)) / 100
		vih := vdd * float64(edakit.Int(56, 65)) / 100
		high := edakit.Int(0, 1) == 1
		ans := mv(vil) - mv(vol)
		name := "NM_L"
		if high {
			ans = mv(voh) - mv(vih)
			name = "NM_H"
		}
		q := edakit.Sprintf("An inverter's VTC has $V_{OH} = %d$ mV, $V_{OL} = %d$ mV, and unity-gain points $V_{IL} = %d$ mV and $V_{IH} = %d$ mV. What is $%s$ in millivolts?",
			mv(voh), mv(vol), mv(vil), mv(vih), name)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("$NM_H = V_{OH} - V_{IH} = %d$ mV and $NM_L = V_{IL} - V_{OL} = %d$ mV: the room between what a driver guarantees and what a receiver needs.", mv(voh)-mv(vih), mv(vil)-mv(vol)),
			map[string]string{"margin_mv": edakit.IntStr(ans)})
	case 2:
		vm := vdd * float64(edakit.Int(42, 58)) / 100
		g := float64(edakit.Int(4, 12))
		vil := vm - (vdd-vm)/g
		vih := vm + vm/g
		high := edakit.Int(0, 1) == 1
		ans, name := vil, "NM_L"
		if high {
			ans, name = vdd-vih, "NM_H"
		}
		q := edakit.Sprintf("Model an inverter VTC as piecewise linear: the output sits at $V_{DD} = %s$ V, then falls along a straight line of slope $-%d$ that passes through $(V_M, V_M)$ with $V_M = %d$ mV, until it reaches 0 V. Take $V_{IL}$ and $V_{IH}$ as the corners of that line. What is $%s$ in millivolts? Round to the nearest millivolt.",
			edakit.Num(vdd, 2), int(g), mv(vm), name)
		return edakit.Numeric(q, edakit.IntStr(mv(ans)),
			edakit.Sprintf("The line meets $V_{DD}$ at $V_{IL} = V_M - (V_{DD} - V_M)/%d = %d$ mV and meets 0 at $V_{IH} = V_M + V_M/%d = %d$ mV. With $V_{OH} = V_{DD}$ and $V_{OL} = 0$, $NM_L = V_{IL}$ and $NM_H = V_{DD} - V_{IH}$.", int(g), mv(vil), int(g), mv(vih)),
			map[string]string{"vil_mv": edakit.IntStr(mv(vil)), "vih_mv": edakit.IntStr(mv(vih))})
	}
	vohD := vdd - float64(edakit.Int(1, 8))*0.01
	volD := float64(edakit.Int(1, 8)) * 0.01
	vilR := vdd * float64(edakit.Int(30, 45)) / 100
	vihR := vdd * float64(edakit.Int(55, 70)) / 100
	nmh, nml := mv(vohD)-mv(vihR), mv(vilR)-mv(volD)
	ans := nmh
	if nml < nmh {
		ans = nml
	}
	q := edakit.Sprintf("A driver cell guarantees $V_{OH} = %d$ mV and $V_{OL} = %d$ mV. The receiver cell it drives has $V_{IH} = %d$ mV and $V_{IL} = %d$ mV. What is the smaller of the two noise margins of this connection, in millivolts?",
		mv(vohD), mv(volD), mv(vihR), mv(vilR))
	return edakit.Numeric(q, edakit.IntStr(ans),
		edakit.Sprintf("Noise margins compare the driver's output levels with the receiver's input thresholds: $NM_H = %d - %d = %d$ mV and $NM_L = %d - %d = %d$ mV. The smaller one limits the connection.", mv(vohD), mv(vihR), nmh, mv(vilR), mv(volD), nml),
		map[string]string{"nmh_mv": edakit.IntStr(nmh), "nml_mv": edakit.IntStr(nml)})
}

// chargeSharing computes the dynamic node voltage after charge sharing.
func chargeSharing(ctx generator.GeneratorContext) generator.Problem {
	vdd := edakit.Pick([]float64{0.8, 1.0, 1.2})
	cout := float64(5 * edakit.Int(3, 8))
	cx := float64(edakit.Int(2, 12))
	v := vdd * cout / (cout + cx)
	f := func(x float64) string { return edakit.Num(x, 3) + " V" }
	q := edakit.Sprintf("A dynamic gate's output node (%s fF) is precharged to $V_{DD} = %s$ V. During evaluation the top transistor of the pull-down stack turns on but the bottom one stays off, connecting the output to an internal node of %s fF that was at 0 V. Assume the two nodes share charge until they equalize. What voltage does the output fall to?",
		edakit.Num(cout, 0), edakit.Num(vdd, 2), edakit.Num(cx, 0))
	return edakit.Choice(q, f(v),
		edakit.Sprintf("Charge is conserved: $C_{out} V_{DD} = (C_{out} + C_x) V$, so $V = %s \\cdot %s/(%s + %s) = %s$ V. If that falls below the next stage's switching threshold the gate fails.", edakit.Num(vdd, 2), edakit.Num(cout, 0), edakit.Num(cout, 0), edakit.Num(cx, 0), edakit.Num(v, 3)),
		f(vdd*cx/(cout+cx)), f(vdd*(cout-cx)/cout), f(vdd/2))
}
