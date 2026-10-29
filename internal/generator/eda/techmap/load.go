package techmap

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Units: resistance in kilohms and capacitance in femtofarads, so their product is in ps.

func libertyDelay(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Int(2, 3)
		var ins []string
		worst, worstPin := 0.0, ""
		for i := 0; i < n; i++ {
			a, d := float64(edakit.Int(10, 80)), float64(edakit.Int(8, 25))
			ins = append(ins, edakit.Sprintf("input %s arrives at %s ps with intrinsic pin delay %s ps", piName(i), edakit.Num(a, 0), edakit.Num(d, 0)))
			if a+d > worst {
				worst, worstPin = a+d, piName(i)
			}
		}
		R := float64(edakit.Int(5, 30)) / 10
		nf := edakit.Int(1, 4)
		var caps []string
		ct := 0.0
		for i := 0; i < nf; i++ {
			c := float64(edakit.Int(2, 8))
			ct += c
			caps = append(caps, edakit.Num(c, 0))
		}
		w := float64(edakit.Int(0, 10))
		ct += w
		ans := worst + R*ct
		q := edakit.Sprintf("A candidate match at node n uses a cell whose delay from each input pin is the pin's intrinsic delay plus R x C_load, with drive resistance R = %s kOhm. In the match, %s. The node drives %d fanout pins with input capacitances %s fF plus %s fF of estimated wire. What is the arrival time at n's output for this match, in ps? Round to 1 decimal place.",
			edakit.Num(R, 1), strings.Join(ins, "; "), nf, strings.Join(caps, ", "), edakit.Num(w, 0))
		return edakit.Numeric(q, edakit.Num(ans, 1),
			edakit.Sprintf("C_load = %s fF, so the load term is %s x %s = %s ps. The latest input plus intrinsic delay is through %s at %s ps. Arrival = %s ps.", edakit.Num(ct, 0), edakit.Num(R, 1), edakit.Num(ct, 0), edakit.Num(R*ct, 1), worstPin, edakit.Num(worst, 0), edakit.Num(ans, 1)),
			map[string]string{"c_load": edakit.Num(ct, 0), "arrival": edakit.Num(ans, 1)})
	case 2:
		s1, s2 := 20.0, float64(edakit.Pick([]int{60, 80, 100}))
		c1, c2 := 2.0, float64(edakit.Pick([]int{8, 10, 16}))
		v11 := float64(edakit.Int(15, 25))
		v12 := v11 + float64(edakit.Int(10, 30))
		v21 := v11 + float64(edakit.Int(5, 15))
		v22 := v12 + (v21 - v11) + float64(edakit.Int(0, 8))
		s := s1 + float64(edakit.Int(1, 9))*(s2-s1)/10
		c := c1 + float64(edakit.Int(1, 9))*(c2-c1)/10
		ts, tc := (s-s1)/(s2-s1), (c-c1)/(c2-c1)
		ans := (1-ts)*(1-tc)*v11 + (1-ts)*tc*v12 + ts*(1-tc)*v21 + ts*tc*v22
		q := edakit.Sprintf("A Liberty NLDM delay table for the arc of a matched cell has rows for input slew %s and %s ps and columns for load %s and %s fF. The entries are %s ps (slew %s, load %s), %s ps (slew %s, load %s), %s ps (slew %s, load %s) and %s ps (slew %s, load %s). The mapper evaluates the match at slew %s ps and load %s fF. What delay does bilinear interpolation give, in ps? Round to 2 decimal places.",
			edakit.Num(s1, 0), edakit.Num(s2, 0), edakit.Num(c1, 0), edakit.Num(c2, 0),
			edakit.Num(v11, 0), edakit.Num(s1, 0), edakit.Num(c1, 0), edakit.Num(v12, 0), edakit.Num(s1, 0), edakit.Num(c2, 0),
			edakit.Num(v21, 0), edakit.Num(s2, 0), edakit.Num(c1, 0), edakit.Num(v22, 0), edakit.Num(s2, 0), edakit.Num(c2, 0),
			edakit.Num(s, 0), edakit.Num(c, 1))
		return edakit.Numeric(q, edakit.Num(ans, 2),
			edakit.Sprintf("The fractional positions are %s along slew and %s along load. Weighting the four corners by (1 - t_s)(1 - t_c), (1 - t_s) t_c, t_s (1 - t_c) and t_s t_c gives %s ps.", edakit.Num(ts, 3), edakit.Num(tc, 3), edakit.Num(ans, 2)),
			map[string]string{"t_slew": edakit.Num(ts, 3), "t_load": edakit.Num(tc, 3), "delay": edakit.Num(ans, 2)})
	}
	rp := float64(edakit.Int(10, 30)) / 10
	cin1 := float64(edakit.Int(2, 4))
	ratio := float64(edakit.Pick([]int{2, 3}))
	cin2 := cin1 * ratio
	r2 := float64(edakit.Int(5, 10)) / 10
	r1 := r2 * ratio
	p1 := float64(edakit.Int(8, 14))
	p2 := p1 + float64(edakit.Int(0, 3))
	cl := float64(edakit.Int(10, 40))
	d := rp*(cin2-cin1) + (p2 - p1) + (r2-r1)*cl
	q := edakit.Sprintf("After mapping, a gate is swapped for a larger version of the same cell. Its input pin capacitance rises from %s to %s fF, its drive resistance falls from %s to %s kOhm and its intrinsic delay rises from %s to %s ps. It drives %s fF, and the gate before it has drive resistance %s kOhm. Using delay = intrinsic + R x C_load for both stages, by how many ps does the two-stage path delay change? Give a negative number if the path gets faster. Round to 1 decimal place.",
		edakit.Num(cin1, 0), edakit.Num(cin2, 0), edakit.Num(r1, 1), edakit.Num(r2, 1), edakit.Num(p1, 0), edakit.Num(p2, 0), edakit.Num(cl, 0), edakit.Num(rp, 1))
	return edakit.Numeric(q, edakit.Num(d, 1),
		edakit.Sprintf("The previous stage slows by %s x %s = %s ps. The gate itself changes by (%s - %s) + (%s - %s) x %s = %s ps. The net change is %s ps.",
			edakit.Num(rp, 1), edakit.Num(cin2-cin1, 0), edakit.Num(rp*(cin2-cin1), 2), edakit.Num(p2, 0), edakit.Num(p1, 0), edakit.Num(r2, 1), edakit.Num(r1, 1), edakit.Num(cl, 0), edakit.Num((p2-p1)+(r2-r1)*cl, 2), edakit.Num(d, 1)),
		map[string]string{"delta": edakit.Num(d, 1)})
}

func loadCrossover(ctx generator.GeneratorContext) generator.Problem {
	dA := float64(edakit.Int(8, 15))
	dB := dA + float64(edakit.Int(3, 12))
	rB := float64(edakit.Int(5, 12)) / 10
	rA := rB + float64(edakit.Int(5, 15))/10
	cx := (dB - dA) / (rA - rB)
	desc := edakit.Sprintf("Two library matches implement node n: match A has intrinsic delay %s ps and drive resistance %s kOhm; match B has intrinsic delay %s ps and drive resistance %s kOhm. Delay is intrinsic + R x C_load.", edakit.Num(dA, 0), edakit.Num(rA, 1), edakit.Num(dB, 0), edakit.Num(rB, 1))
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(desc+" Above what load does match B become faster than match A, in fF? Round to 1 decimal place.", edakit.Num(cx, 1),
			edakit.Sprintf("Set %s + %s C = %s + %s C, so C = (%s - %s) / (%s - %s) = %s fF.", edakit.Num(dA, 0), edakit.Num(rA, 1), edakit.Num(dB, 0), edakit.Num(rB, 1), edakit.Num(dB, 0), edakit.Num(dA, 0), edakit.Num(rA, 1), edakit.Num(rB, 1), edakit.Num(cx, 1)),
			map[string]string{"crossover": edakit.Num(cx, 1)})
	case 2:
		c := float64(edakit.Int(2, 40))
		a := float64(edakit.Int(20, 80))
		tA, tB := a+dA+rA*c, a+dB+rB*c
		best, which := tA, "A"
		if tB < tA {
			best, which = tB, "B"
		}
		return edakit.Numeric(desc+edakit.Sprintf(" The latest input of either match arrives at %s ps and the actual load on n is %s fF. What is the best arrival time at n, in ps? Round to 1 decimal place.", edakit.Num(a, 0), edakit.Num(c, 0)), edakit.Num(best, 1),
			edakit.Sprintf("Match A gives %s ps and match B gives %s ps, so match %s wins with %s ps. The winner flips at %s fF.", edakit.Num(tA, 1), edakit.Num(tB, 1), which, edakit.Num(best, 1), edakit.Num(cx, 1)),
			map[string]string{"best": edakit.Num(best, 1), "crossover": edakit.Num(cx, 1)})
	}
	// Assumed load on one side of the crossover, actual on the other.
	var cEst, cAct float64
	if edakit.Int(0, 1) == 0 {
		cEst = cx * float64(edakit.Int(3, 7)) / 10
		cAct = cx * float64(edakit.Int(14, 30)) / 10
	} else {
		cEst = cx * float64(edakit.Int(14, 30)) / 10
		cAct = cx * float64(edakit.Int(3, 7)) / 10
	}
	cEst, cAct = edakit.Round(cEst, 0), edakit.Round(cAct, 0)
	if cEst == cAct {
		cAct++
	}
	pick := func(c float64) (string, float64) {
		tA, tB := dA+rA*c, dB+rB*c
		if tB < tA {
			return "B", tB
		}
		return "A", tA
	}
	chosen, _ := pick(cEst)
	_, bestAct := pick(cAct)
	var chosenAct float64
	if chosen == "A" {
		chosenAct = dA + rA*cAct
	} else {
		chosenAct = dB + rB*cAct
	}
	pen := chosenAct - bestAct
	return edakit.Numeric(desc+edakit.Sprintf(" During bottom-up mapping the fanout of n is not yet known, so the mapper assumes a load of %s fF and picks the faster match for that load. After mapping, the real load turns out to be %s fF. How many ps slower is n than if the mapper had known the real load? Round to 1 decimal place.", edakit.Num(cEst, 0), edakit.Num(cAct, 0)),
		edakit.Num(pen, 1),
		edakit.Sprintf("At %s fF the mapper picks match %s. At the real %s fF that match takes %s ps, while the better match takes %s ps: a penalty of %s ps.", edakit.Num(cEst, 0), chosen, edakit.Num(cAct, 0), edakit.Num(chosenAct, 1), edakit.Num(bestAct, 1), edakit.Num(pen, 1)),
		map[string]string{"penalty": edakit.Num(pen, 1), "crossover": edakit.Num(cx, 1)})
}

type leGate struct {
	name string
	g    float64
	gs   string
	p    float64
}

var leGates = map[string]leGate{
	"INV":   {"INV", 1, "1", 1},
	"NAND2": {"NAND2", 4.0 / 3, "4/3", 2},
	"NAND3": {"NAND3", 5.0 / 3, "5/3", 3},
	"NAND4": {"NAND4", 2, "2", 4},
	"NOR2":  {"NOR2", 5.0 / 3, "5/3", 2},
	"NOR3":  {"NOR3", 7.0 / 3, "7/3", 3},
}

func leTable(names []string) string {
	var s []string
	for _, n := range names {
		g := leGates[n]
		s = append(s, edakit.Sprintf("%s (g = %s, p = %s)", n, g.gs, edakit.Num(g.p, 0)))
	}
	return strings.Join(s, ", ")
}

func pathDelay(path []string, h float64) float64 {
	t := 0.0
	for _, n := range path {
		t += leGates[n].g*h + leGates[n].p
	}
	return t
}

func gainBased(ctx generator.GeneratorContext) generator.Problem {
	h := float64(edakit.Int(4, 10)) / 2
	intro := edakit.Sprintf("In gain-based mapping every gate is held at a fixed electrical effort (gain) h = C_out / C_in = %s, so a gate's delay is d = g h + p in units of tau, independent of the load it drives.", edakit.Num(h, 1))
	all := []string{"INV", "NAND2", "NAND3", "NAND4", "NOR2", "NOR3"}
	switch kp(ctx) {
	case 1:
		n := edakit.Int(2, 4)
		var path []string
		for i := 0; i < n; i++ {
			path = append(path, edakit.Pick(all))
		}
		d := pathDelay(path, h)
		return edakit.Numeric(intro+" Logical efforts and parasitics: "+leTable(distinctStrings(path))+". What is the delay of the path "+strings.Join(path, " -> ")+", in tau? Round to 2 decimal places.",
			edakit.Num(d, 2),
			edakit.Sprintf("Sum g h + p over the stages with h = %s: %s tau.", edakit.Num(h, 1), edakit.Num(d, 2)),
			map[string]string{"delay": edakit.Num(d, 2)})
	case 2:
		nf := edakit.Int(1, 4)
		var cs []string
		t := 0.0
		for i := 0; i < nf; i++ {
			c := float64(edakit.Int(2, 12))
			t += c
			cs = append(cs, edakit.Num(c, 0))
		}
		w := float64(edakit.Int(0, 6))
		t += w
		cin := t / h
		return edakit.Numeric(intro+edakit.Sprintf(" After mapping, a gate drives fanout pins of %s fF and %s fF of wire. To keep its gain at h, what input capacitance must the sizer give it, in fF? Round to 2 decimal places.", strings.Join(cs, ", "), edakit.Num(w, 0)),
			edakit.Num(cin, 2),
			edakit.Sprintf("C_out = %s fF, so C_in = C_out / h = %s / %s = %s fF. Sizing then realizes the delays the mapper assumed.", edakit.Num(t, 0), edakit.Num(t, 0), edakit.Num(h, 1), edakit.Num(cin, 2)),
			map[string]string{"c_out": edakit.Num(t, 0), "c_in": edakit.Num(cin, 2)})
	}
	type alt struct {
		fn   string
		a, b []string
	}
	alts := []alt{
		{"a 3-input NAND", []string{"NAND3"}, []string{"NAND2", "INV", "NAND2"}},
		{"a 4-input NAND", []string{"NAND4"}, []string{"NAND2", "INV", "NAND2"}},
		{"a 3-input NOR", []string{"NOR3"}, []string{"NOR2", "INV", "NOR2"}},
		{"a 2-input AND", []string{"NAND2", "INV"}, []string{"INV", "NOR2"}},
		{"a 2-input OR", []string{"NOR2", "INV"}, []string{"INV", "NAND2"}},
		{"a 4-input AND", []string{"NAND4", "INV"}, []string{"NAND2", "NOR2"}},
	}
	x := edakit.Pick(alts)
	da, db := pathDelay(x.a, h), pathDelay(x.b, h)
	best := da
	if db < best {
		best = db
	}
	names := distinctStrings(append(append([]string{}, x.a...), x.b...))
	return edakit.Numeric(intro+" Logical efforts and parasitics: "+leTable(names)+". The mapper can implement "+x.fn+" with the critical path "+strings.Join(x.a, " -> ")+" or with the critical path "+strings.Join(x.b, " -> ")+". What is the delay of the faster choice, in tau? Round to 2 decimal places.",
		edakit.Num(best, 2),
		edakit.Sprintf("The first path takes %s tau and the second %s tau. Because gain is fixed, the comparison depends only on logical effort and parasitic delay, not on the load.", edakit.Num(da, 2), edakit.Num(db, 2)),
		map[string]string{"delay_a": edakit.Num(da, 2), "delay_b": edakit.Num(db, 2)})
}
