package cmos

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func kpOf(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal == 0 {
		return edakit.Int(1, 3)
	}
	return ctx.KPOrdinal
}

// rcModel: Elmore delay of a gate driving a load, width scaling, and the ln 2 factor.
func rcModel(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		r := float64(edakit.Int(2, 12))
		cp := float64(edakit.Int(1, 6))
		cl := float64(edakit.Scaled(ctx, 2, 40))
		t := r * (cp + cl)
		ans := edakit.Num(t, 1)
		q := edakit.Sprintf("An inverter's pull-down has an effective resistance of %s kOhm. Its own drain diffusion adds %s fF at the output and it drives a load of %s fF. Using the RC (Elmore) model $t = R(C_{par} + C_{load})$, what is the falling delay in picoseconds?",
			edakit.Num(r, 0), edakit.Num(cp, 0), edakit.Num(cl, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The resistance charges both capacitances: $%s \\times (%s + %s) = %s$ kOhm fF, and one kOhm times one fF is one ps.", edakit.Num(r, 0), edakit.Num(cp, 0), edakit.Num(cl, 0), ans),
			map[string]string{"delay_ps": ans})
	case 2:
		rn := float64(edakit.Int(4, 12))
		mu := edakit.Pick([]float64{2, 2.5, 3})
		wp := float64(edakit.Int(1, 6))
		cp := float64(edakit.Int(1, 8))
		cl := float64(edakit.Scaled(ctx, 4, 60))
		rp := mu * rn / wp
		t := rp * (cp + cl)
		ans := edakit.Num(t, 2)
		q := edakit.Sprintf("A unit-width nMOS has an effective resistance of %s kOhm, and a pMOS of the same width has %s times that resistance. Resistance scales as 1/width. An inverter uses a pMOS of width %s units, has %s fF of parasitic output capacitance, and drives %s fF. What is its rising delay in picoseconds under the RC model? Round to 2 decimal places.",
			edakit.Num(rn, 0), edakit.Num(mu, 1), edakit.Num(wp, 0), edakit.Num(cp, 0), edakit.Num(cl, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The rising edge goes through the pMOS: $R_p = %s \\times %s / %s = %s$ kOhm. Then $t = R_p (C_{par} + C_{load}) = %s$ ps.", edakit.Num(mu, 1), edakit.Num(rn, 0), edakit.Num(wp, 0), edakit.Num(rp, 3), ans),
			map[string]string{"rp_kohm": edakit.Num(rp, 3), "delay_ps": ans})
	}
	r := float64(edakit.Int(2, 15))
	c := float64(edakit.Scaled(ctx, 5, 80))
	t := ln(2) * r * c
	ans := edakit.Num(t, 2)
	q := edakit.Sprintf("A gate is modeled as a step source behind %s kOhm driving %s fF to ground. How long, in picoseconds, does the output take to cross 50%% of its swing? Round to 2 decimal places.",
		edakit.Num(r, 0), edakit.Num(c, 0))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("A single RC stage reaches 50%% when $e^{-t/RC} = 1/2$, so $t = \\ln 2 \\cdot RC = 0.693 \\times %s = %s$ ps.", edakit.Num(r*c, 0), ans),
		map[string]string{"rc_ps": edakit.Num(r*c, 0), "t50_ps": ans})
}

// equalRiseFall sizes gates for the same worst-case drive as a unit inverter.
func equalRiseFall(ctx generator.GeneratorContext) generator.Problem {
	gamma := edakit.Pick([]int{2, 3})
	pre := edakit.Sprintf("A unit inverter has an nMOS of width 1 and a pMOS of width %d (holes are %d times slower). Size every gate so that its worst-case pull-down and pull-up resistances equal the unit inverter's. ", gamma, gamma)
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 2, 4)
		if edakit.Int(0, 1) == 0 {
			if edakit.Int(0, 1) == 0 {
				return edakit.Numeric(pre+edakit.Sprintf("What width does each nMOS of a NAND%d get?", n), edakit.IntStr(n),
					edakit.Sprintf("The %d nMOS devices are in series, so each must be %d times wider to give the unit resistance.", n, n), nil)
			}
			return edakit.Numeric(pre+edakit.Sprintf("What width does each pMOS of a NAND%d get?", n), edakit.IntStr(gamma),
				edakit.Sprintf("The pMOS devices are in parallel and the worst case is one of them on alone, so each keeps the inverter's width %d.", gamma), nil)
		}
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(pre+edakit.Sprintf("What width does each pMOS of a NOR%d get?", n), edakit.IntStr(n*gamma),
				edakit.Sprintf("The %d pMOS devices are in series, so each must be %d times the inverter's pMOS width: %d.", n, n, n*gamma), nil)
		}
		return edakit.Numeric(pre+edakit.Sprintf("What width does each nMOS of a NOR%d get?", n), "1",
			"The nMOS devices are in parallel and the worst case is one on alone, so each keeps width 1.", nil)
	case 2:
		n := edakit.Scaled(ctx, 2, 4)
		if edakit.Int(0, 1) == 0 {
			c := n + gamma
			return edakit.Numeric(pre+edakit.Sprintf("Gate capacitance is proportional to width. What is the input capacitance of one input of a NAND%d, in units where a width-1 transistor has capacitance 1?", n),
				edakit.IntStr(c), edakit.Sprintf("Each input drives one nMOS of width %d and one pMOS of width %d, for %d.", n, gamma, c), nil)
		}
		c := 1 + n*gamma
		return edakit.Numeric(pre+edakit.Sprintf("Gate capacitance is proportional to width. What is the input capacitance of one input of a NOR%d, in units where a width-1 transistor has capacitance 1?", n),
			edakit.IntStr(c), edakit.Sprintf("Each input drives one nMOS of width 1 and one pMOS of width %d, for %d.", n*gamma, c), nil)
	}
	e := randCompound(ctx)
	lit := edakit.Pick(e.lits())
	f := "$Y = \\overline{" + e.latex() + "}$"
	if edakit.Int(0, 1) == 0 {
		w := e.widths(false)[lit]
		return edakit.Numeric(pre+"Each transistor's width is the length of the longest series path through it times the unit width. For the gate "+f+", what width does the nMOS driven by input "+lit+" get?",
			edakit.IntStr(w), edakit.Sprintf("In the pull-down network AND is series. The longest series path through %s's nMOS has %d devices, so its width is %d.", lit, w, w),
			map[string]string{"width": edakit.IntStr(w)})
	}
	w := e.widths(true)[lit] * gamma
	return edakit.Numeric(pre+"Each transistor's width is the length of the longest series path through it times the unit width. For the gate "+f+", what width does the pMOS driven by input "+lit+" get?",
		edakit.IntStr(w), edakit.Sprintf("In the pull-up network OR is series. The longest series path through %s's pMOS has %d devices, so its width is %d times %d, which is %d.", lit, e.widths(true)[lit], e.widths(true)[lit], gamma, w),
		map[string]string{"width": edakit.IntStr(w)})
}

// stackElmore: Elmore delay of a series nMOS stack discharging the output.
func stackElmore(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	n := 2
	if kp > 1 {
		n = edakit.Scaled(ctx, 3, 4)
	}
	r := float64(edakit.Int(2, 8))
	cs := make([]float64, n)
	var parts []string
	for j := 1; j < n; j++ {
		cs[j] = float64(edakit.Int(1, 6))
		parts = append(parts, edakit.Sprintf("$C_%d = %s$ fF", j, edakit.Num(cs[j], 0)))
	}
	cout := float64(edakit.Scaled(ctx, 4, 30))
	internal := 0.0
	for j := 1; j < n; j++ {
		internal += float64(j) * r * cs[j]
	}
	full := internal + float64(n)*r*cout
	desc := edakit.Sprintf("A NAND%d pull-down stack has %d series nMOS transistors, each %s kOhm when on. Number the internal nodes from the ground end: %s (node $j$ sits above $j$ transistors). The output carries %s fF.",
		n, n, edakit.Num(r, 0), strings.Join(parts, ", "), edakit.Num(cout, 0))
	if kp == 3 {
		ans := edakit.Num(internal, 1)
		q := desc + " In the worst case all internal nodes start charged high and the bottom transistor switches last. How many picoseconds of Elmore delay are saved if instead the latest-arriving input drives the transistor nearest the output, so the internal nodes have already discharged?"
		return edakit.Numeric(q, ans,
			edakit.Sprintf("With the top transistor switching last only the output needs discharging, so the saving is the internal nodes' terms $\\sum_j jRC_j = %s$ ps.", ans),
			map[string]string{"internal_ps": ans, "full_ps": edakit.Num(full, 1)})
	}
	ans := edakit.Num(full, 1)
	q := desc + " All internal nodes start charged high. What is the Elmore delay of the falling output, in picoseconds?"
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Each node's capacitance is multiplied by the resistance between it and ground: internal node $j$ sees $jR$ and the output sees $%dR$. The internal terms give %s ps and the output term %s ps, for %s ps.", n, edakit.Num(internal, 1), edakit.Num(float64(n)*r*cout, 1), ans),
		map[string]string{"internal_ps": edakit.Num(internal, 1), "full_ps": ans})
}

// gateKind is a named gate with its logical effort and parasitic delay under gamma = 2,
// p_inv = 1 sizing.
type gateKind struct {
	name string
	g, p float64
}

var pathGates = []gateKind{
	{"inverter", 1, 1}, {"NAND2", 4.0 / 3, 2}, {"NAND3", 5.0 / 3, 3}, {"NOR2", 5.0 / 3, 2}, {"NOR3", 7.0 / 3, 3},
}

func nandG(n int, gamma float64) float64 { return (float64(n) + gamma) / (1 + gamma) }
func norG(n int, gamma float64) float64  { return (1 + float64(n)*gamma) / (1 + gamma) }

// logicalEffort computes g for NAND, NOR and compound gates.
func logicalEffort(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 2, 5)
		if edakit.Int(0, 1) == 0 {
			g := nandG(n, 2)
			ans := edakit.Num(g, 2)
			return edakit.Numeric(edakit.Sprintf("With pMOS twice as wide as nMOS for equal drive ($\\gamma = 2$), what is the logical effort of one input of a NAND%d? Round to 2 decimal places.", n), ans,
				edakit.Sprintf("Each input drives an nMOS of width %d and a pMOS of width 2, so its capacitance is %d against the inverter's 3: $g = %d/3 = %s$.", n, n+2, n+2, ans), map[string]string{"g": ans})
		}
		g := norG(n, 2)
		ans := edakit.Num(g, 2)
		return edakit.Numeric(edakit.Sprintf("With pMOS twice as wide as nMOS for equal drive ($\\gamma = 2$), what is the logical effort of one input of a NOR%d? Round to 2 decimal places.", n), ans,
			edakit.Sprintf("Each input drives an nMOS of width 1 and a pMOS of width %d, so its capacitance is %d against the inverter's 3: $g = %d/3 = %s$.", 2*n, 2*n+1, 2*n+1, ans), map[string]string{"g": ans})
	case 2:
		n := edakit.Scaled(ctx, 2, 4)
		gamma := edakit.Pick([]float64{1.5, 2.5, 3, 4})
		nand := edakit.Int(0, 1) == 0
		name, g := "NOR", norG(n, gamma)
		if nand {
			name, g = "NAND", nandG(n, gamma)
		}
		ans := edakit.Num(g, 3)
		return edakit.Numeric(edakit.Sprintf("In a process where the pMOS must be $\\gamma = %s$ times as wide as the nMOS for equal drive, what is the logical effort of one input of a %s%d sized for unit-inverter drive? Round to 3 decimal places.", edakit.Num(gamma, 1), name, n), ans,
			edakit.Sprintf("The unit inverter's input capacitance is $1 + \\gamma = %s$. For NAND$n$ an input sees $n + \\gamma$; for NOR$n$ it sees $1 + n\\gamma$. Dividing gives %s.", edakit.Num(1+gamma, 1), ans), map[string]string{"g": ans})
	}
	e := randCompound(ctx)
	lit := edakit.Pick(e.lits())
	wn, wp := e.widths(false)[lit], 2*e.widths(true)[lit]
	g := float64(wn+wp) / 3
	ans := edakit.Num(g, 2)
	q := "With $\\gamma = 2$ (unit inverter nMOS 1, pMOS 2), the gate $Y = \\overline{" + e.latex() + "}$ is sized so its worst-case pull-down and pull-up match the unit inverter. What is the logical effort of input " + lit + "? Round to 2 decimal places."
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Input %s drives an nMOS of width %d (longest series path in the pull-down) and a pMOS of width %d. Its capacitance %d divided by the inverter's 3 gives %s.", lit, wn, wp, wn+wp, ans),
		map[string]string{"cin": edakit.IntStr(wn + wp), "g": ans})
}

// parasiticDelay computes p from output diffusion, and the parasitic share of a stage delay.
func parasiticDelay(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 2, 5)
		pinv := edakit.Pick([]float64{0.8, 1, 1.2})
		name := edakit.Pick([]string{"NAND", "NOR"})
		ans := edakit.Num(float64(n)*pinv, 2)
		q := edakit.Sprintf("In a process with $\\gamma = 2$ and an inverter parasitic delay $p_{inv} = %s$, what is the parasitic delay of a %s%d sized for unit-inverter drive, counting only diffusion on the output node? Round to 2 decimal places.", edakit.Num(pinv, 1), name, n)
		expl := edakit.Sprintf("The NAND%d output touches %d parallel pMOS of width 2 and one nMOS of width %d: total %d, against the inverter's 3, so $p = %d p_{inv}$.", n, n, n, 3*n, n)
		if name == "NOR" {
			expl = edakit.Sprintf("The NOR%d output touches %d parallel nMOS of width 1 and one pMOS of width %d: total %d, against the inverter's 3, so $p = %d p_{inv}$.", n, n, 2*n, 3*n, n)
		}
		return edakit.Numeric(q, ans, expl, map[string]string{"p": ans})
	case 2:
		gamma := edakit.Pick([]int{2, 3})
		a := edakit.Int(1, 6)
		b := gamma * edakit.Int(1, 6)
		pinv := 1.0
		p := float64(a+b) / float64(1+gamma) * pinv
		ans := edakit.Num(p, 2)
		q := edakit.Sprintf("A unit inverter (nMOS width 1, pMOS width %d) has $p_{inv} = 1$. Diffusion capacitance is proportional to the width of the devices touching a node. A complex gate's output node touches nMOS diffusion of total width %d and pMOS diffusion of total width %d. What is the gate's parasitic delay? Round to 2 decimal places.", gamma, a, b)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Parasitic delay scales with output diffusion: $p = (%d + %d)/(1 + %d) \\cdot p_{inv} = %s$.", a, b, gamma, ans), map[string]string{"p": ans})
	}
	gk := edakit.Pick(pathGates)
	h := float64(edakit.Int(1, 8))
	d := gk.g*h + gk.p
	share := 100 * gk.p / d
	ans := edakit.Num(share, 1)
	q := edakit.Sprintf("%s ($g = %s$, $p = %s$) drives an electrical effort of $h = %s$. What percentage of its stage delay is parasitic delay? Round to 1 decimal place.", capArticle(gk.name), edakit.Num(gk.g, 3), edakit.Num(gk.p, 0), edakit.Num(h, 0))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Stage delay is $d = gh + p = %s$. The parasitic part is $p/d = %s/%s$, which is %s%%.", edakit.Num(d, 3), edakit.Num(gk.p, 0), edakit.Num(d, 3), ans),
		map[string]string{"d": edakit.Num(d, 3)})
}

// stageDelay computes d = gh + p.
func stageDelay(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		g := edakit.Pick([]float64{1, 4.0 / 3, 5.0 / 3, 2, 7.0 / 3})
		p := float64(edakit.Int(1, 4))
		cin := float64(edakit.Int(2, 10))
		cout := float64(edakit.Scaled(ctx, 4, 80))
		h := cout / cin
		d := g*h + p
		ans := edakit.Num(d, 2)
		q := edakit.Sprintf("A gate has logical effort $g = %s$ and parasitic delay $p = %s$. Its input capacitance is %s fF and it drives %s fF. What is its stage delay $d$ in units of $\\tau$? Round to 2 decimal places.", edakit.Num(g, 3), edakit.Num(p, 0), edakit.Num(cin, 0), edakit.Num(cout, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Electrical effort is $h = C_{out}/C_{in} = %s$, so $d = gh + p = %s \\times %s + %s = %s$.", edakit.Num(h, 3), edakit.Num(g, 3), edakit.Num(h, 3), edakit.Num(p, 0), ans),
			map[string]string{"h": edakit.Num(h, 4), "d": ans})
	case 2:
		n := edakit.Scaled(ctx, 2, 4)
		nand := edakit.Int(0, 1) == 0
		name, g := "NOR", norG(n, 2)
		if nand {
			name, g = "NAND", nandG(n, 2)
		}
		h := float64(edakit.Int(1, 8))
		d := g*h + float64(n)
		ans := edakit.Num(d, 2)
		q := edakit.Sprintf("With $\\gamma = 2$ and $p_{inv} = 1$, a %s%d drives a load %s times its own input capacitance. What is its stage delay in units of $\\tau$? Round to 2 decimal places.", name, n, edakit.Num(h, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The %s%d has $g = %s$ and $p = %d$, so $d = %s \\times %s + %d = %s$.", name, n, edakit.Num(g, 3), n, edakit.Num(g, 3), edakit.Num(h, 0), n, ans),
			map[string]string{"g": edakit.Num(g, 4), "d": ans})
	}
	tau := float64(edakit.Int(3, 15))
	a, b := edakit.Pick(pathGates), edakit.Pick(pathGates)
	h1, h2 := float64(edakit.Int(1, 6)), float64(edakit.Int(1, 6))
	d := a.g*h1 + a.p + b.g*h2 + b.p
	ans := edakit.Num(d*tau, 1)
	q := edakit.Sprintf("With $\\tau = %s$ ps, %s ($g = %s$, $p = %s$) with electrical effort %s drives %s ($g = %s$, $p = %s$) with electrical effort %s. What is the delay of the two stages in picoseconds? Round to 1 decimal place.",
		edakit.Num(tau, 0), article(a.name), edakit.Num(a.g, 3), edakit.Num(a.p, 0), edakit.Num(h1, 0), article(b.name), edakit.Num(b.g, 3), edakit.Num(b.p, 0), edakit.Num(h2, 0))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Sum the stage delays $g_ih_i + p_i$ to get %s in units of $\\tau$, then multiply by %s ps to get %s ps.", edakit.Num(d, 3), edakit.Num(tau, 0), ans),
		map[string]string{"d_tau": edakit.Num(d, 4), "d_ps": ans})
}

// fo4 converts between tau, the FO4 delay and path delays in FO4 units.
func fo4(ctx generator.GeneratorContext) generator.Problem {
	pinv := edakit.Pick([]float64{0.8, 1, 1.2})
	switch kpOf(ctx) {
	case 1:
		tau := float64(edakit.Int(2, 15))
		v := (4 + pinv) * tau
		ans := edakit.Num(v, 2)
		q := edakit.Sprintf("A process has $\\tau = %s$ ps and $p_{inv} = %s$. What is the delay of an inverter driving four copies of itself (its FO4 delay), in picoseconds? Round to 2 decimal places.", edakit.Num(tau, 0), edakit.Num(pinv, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("An inverter has $g = 1$, and four identical loads give $h = 4$, so $d = 4 + p_{inv} = %s$ and the delay is %s ps.", edakit.Num(4+pinv, 1), ans), nil)
	case 2:
		f := float64(edakit.Int(8, 60))
		tau := f / (4 + pinv)
		ans := edakit.Num(tau, 2)
		q := edakit.Sprintf("A ring-oscillator measurement gives an FO4 inverter delay of %s ps in a process with $p_{inv} = %s$. What is $\\tau$ in picoseconds? Round to 2 decimal places.", edakit.Num(f, 0), edakit.Num(pinv, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("FO4 delay is $(4 + p_{inv})\\tau$, so $\\tau = %s/%s = %s$ ps.", edakit.Num(f, 0), edakit.Num(4+pinv, 1), ans), nil)
	}
	f := float64(edakit.Int(8, 40))
	period := float64(edakit.Scaled(ctx, 150, 1500))
	ans := edakit.Num(period/f, 2)
	q := edakit.Sprintf("A design's clock period is %s ps and the process FO4 delay is %s ps. How many FO4 delays fit in one cycle? Round to 2 decimal places.", edakit.Num(period, 0), edakit.Num(f, 0))
	return edakit.Numeric(q, ans, edakit.Sprintf("Dividing gives %s FO4 per cycle, a process-independent measure of how much logic each pipeline stage holds.", ans), nil)
}

func drawPath(n int) []gateKind {
	out := make([]gateKind, n)
	for i := range out {
		out[i] = edakit.Pick(pathGates)
	}
	return out
}

func pathText(gs []gateKind) string {
	var names []string
	for _, g := range gs {
		names = append(names, g.name)
	}
	for i := range names {
		names[i] = article(names[i])
	}
	return strings.Join(names, ", then ")
}

const pathLegend = "Use $\\gamma = 2$ and $p_{inv} = 1$: inverter $g = 1, p = 1$; NAND2 $g = 4/3, p = 2$; NAND3 $g = 5/3, p = 3$; NOR2 $g = 5/3, p = 2$; NOR3 $g = 7/3, p = 3$."

// pathEffort computes F = GBH, with branching, and the minimum path delay.
func pathEffort(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	n := edakit.Scaled(ctx, 2, 4)
	gs := drawPath(n)
	G, P := 1.0, 0.0
	for _, g := range gs {
		G *= g.g
		P += g.p
	}
	cin := float64(edakit.Int(2, 8))
	cout := float64(edakit.Scaled(ctx, 20, 400))
	H := cout / cin
	switch kp {
	case 1:
		F := G * H
		ans := edakit.Num(F, 2)
		q := edakit.Sprintf("A path is %s. Its first input is %s fF and its last output drives %s fF, with no branching. %s What is the path effort $F$? Round to 2 decimal places.", pathText(gs), edakit.Num(cin, 0), edakit.Num(cout, 0), pathLegend)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$G = \\prod g_i = %s$ and $H = C_{out}/C_{in} = %s$, so $F = GH = %s$.", edakit.Num(G, 4), edakit.Num(H, 4), ans),
			map[string]string{"G": edakit.Num(G, 4), "H": edakit.Num(H, 4), "F": ans})
	case 2:
		ratio := edakit.Pick([]float64{0.5, 1, 1.5, 2, 3, 4})
		b := 1 + ratio
		F := G * b * H
		ans := edakit.Num(F, 2)
		q := edakit.Sprintf("A path is %s. Its first input is %s fF and its last output drives %s fF. At the output of the first gate the wire also feeds off-path gates whose capacitance is %s times the on-path gate's input. %s What is the path effort $F = GBH$? Round to 2 decimal places.", pathText(gs), edakit.Num(cin, 0), edakit.Num(cout, 0), edakit.Num(ratio, 1), pathLegend)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Branching effort is $b = (C_{on} + C_{off})/C_{on} = %s$. With $G = %s$ and $H = %s$, $F = GBH = %s$.", edakit.Num(b, 4), edakit.Num(G, 4), edakit.Num(H, 4), ans),
			map[string]string{"B": edakit.Num(b, 4), "F": ans})
	}
	F := G * H
	fhat := pow(F, 1/float64(n))
	D := float64(n)*fhat + P
	ans := edakit.Num(D, 2)
	q := edakit.Sprintf("A path is %s. Its first input is %s fF and its last output drives %s fF, with no branching. %s If every stage is sized for equal stage effort, what is the minimum path delay $D$ in units of $\\tau$? Round to 2 decimal places.", pathText(gs), edakit.Num(cin, 0), edakit.Num(cout, 0), pathLegend)
	return edakit.Numeric(q, ans,
		edakit.Sprintf("$F = GH = %s$, the best stage effort is $\\hat f = F^{1/%d} = %s$, and $P = %s$. So $D = N\\hat f + P = %s$.", edakit.Num(F, 4), n, edakit.Num(fhat, 4), edakit.Num(P, 0), ans),
		map[string]string{"F": edakit.Num(F, 4), "fhat": edakit.Num(fhat, 4), "D": ans})
}

// optimalStages: how many inverter stages minimize delay for a path effort F.
func optimalStages(ctx generator.GeneratorContext) generator.Problem {
	F := float64(edakit.Pick([]int{12, 20, 30, 45, 64, 100, 150, 250, 400, 640, 1000, 1600, 2500, 4000, 6400}))
	if ctx.Difficulty < 0.4 && F > 1000 {
		F = F / 10
	}
	best, bestD := 1, F+1
	for N := 1; N <= 14; N++ {
		d := float64(N)*pow(F, 1/float64(N)) + float64(N)
		if d < bestD {
			best, bestD = N, d
		}
	}
	switch kpOf(ctx) {
	case 1:
		est := log10(F) / log10(4)
		ans := edakit.Num(est, 2)
		q := edakit.Sprintf("A buffer chain must drive a load whose path effort is $F = %s$. Using the rule of thumb that each stage should bear an effort of about 4, estimate the number of stages $\\log_4 F$. Round to 2 decimal places.", edakit.Num(F, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\log_4 %s = \\ln %s/\\ln 4 = %s$. Round to a nearby integer and compare the neighbours.", edakit.Num(F, 0), edakit.Num(F, 0), ans), nil)
	case 2:
		q := edakit.Sprintf("A chain of inverters ($p_{inv} = 1$) must deliver a path effort $F = %s$. The delay with $N$ stages is $D(N) = N F^{1/N} + N$. Which integer $N$ minimizes $D$?", edakit.Num(F, 0))
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("Evaluating $D(N)$ at neighbouring integers, the minimum is at $N = %d$ with $D = %s$. The continuous optimum puts each stage effort near 3.59 when $p_{inv} = 1$.", best, edakit.Num(bestD, 2)),
			map[string]string{"n_best": edakit.IntStr(best), "d_best": edakit.Num(bestD, 3)})
	}
	ans := edakit.Num(bestD, 2)
	q := edakit.Sprintf("A chain of inverters ($p_{inv} = 1$) must deliver a path effort $F = %s$, and you may choose any integer number of stages $N$. The delay is $D(N) = N F^{1/N} + N$. What is the smallest achievable $D$ in units of $\\tau$? Round to 2 decimal places.", edakit.Num(F, 0))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("The best integer is $N = %d$, giving stage effort %s and $D = %s$.", best, edakit.Num(pow(F, 1/float64(best)), 3), ans),
		map[string]string{"n_best": edakit.IntStr(best), "d_best": ans})
}

// pathSizing sizes gates backwards from the load at equal stage effort.
func pathSizing(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	n := 3
	if kp == 1 {
		n = edakit.Scaled(ctx, 2, 4)
	}
	gs := drawPath(n)
	G := 1.0
	for _, g := range gs {
		G *= g.g
	}
	cin := float64(edakit.Int(2, 10))
	cout := float64(edakit.Scaled(ctx, 40, 600))
	F := G * cout / cin
	fhat := pow(F, 1/float64(n))
	caps := make([]float64, n+1)
	caps[n] = cout
	for i := n - 1; i >= 0; i-- {
		caps[i] = gs[i].g * caps[i+1] / fhat
	}
	base := edakit.Sprintf("A path is %s. The first input is %s fF and the last output drives %s fF, with no branching. %s", pathText(gs), edakit.Num(cin, 0), edakit.Num(cout, 0), pathLegend)
	switch kp {
	case 1:
		ans := edakit.Num(fhat, 3)
		return edakit.Numeric(base+" What stage effort $\\hat f$ minimizes the path delay? Round to 3 decimal places.", ans,
			edakit.Sprintf("$F = GH = %s \\times %s = %s$, and $\\hat f = F^{1/%d} = %s$.", edakit.Num(G, 4), edakit.Num(cout/cin, 4), edakit.Num(F, 4), n, ans),
			map[string]string{"F": edakit.Num(F, 4), "fhat": ans})
	case 2:
		ans := edakit.Num(caps[n-1], 2)
		return edakit.Numeric(base+" At the minimum-delay sizing, what is the input capacitance of the last gate, in fF? Round to 2 decimal places.", ans,
			edakit.Sprintf("$\\hat f = %s$. Each gate's input capacitance is $C_{in,i} = g_i C_{out,i}/\\hat f$, so the last gate has $%s \\times %s / %s = %s$ fF.", edakit.Num(fhat, 4), edakit.Num(gs[n-1].g, 3), edakit.Num(cout, 0), edakit.Num(fhat, 4), ans),
			map[string]string{"fhat": edakit.Num(fhat, 4), "c_last": ans})
	}
	ans := edakit.Num(caps[1], 2)
	return edakit.Numeric(base+" At the minimum-delay sizing, what is the input capacitance of the second gate, in fF? Round to 2 decimal places.", ans,
		edakit.Sprintf("$\\hat f = %s$. Working back from the load with $C_{in,i} = g_i C_{out,i}/\\hat f$: the third gate's input is %s fF and the second gate's is %s fF. A check: the first gate's comes out at %s fF, the given input.", edakit.Num(fhat, 4), edakit.Num(caps[2], 2), ans, edakit.Num(caps[0], 2)),
		map[string]string{"fhat": edakit.Num(fhat, 4), "c2": ans})
}
