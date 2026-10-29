package cts

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// levelsFor is the smallest L with fan^L >= n.
func levelsFor(n, fan int) int {
	l, reach := 0, 1
	for reach < n {
		reach *= fan
		l++
	}
	return l
}

// bufferLevels: KP1 the number of buffer levels, KP2 the total buffer count, KP3 the fanout
// that minimises latency under a fanout-dependent stage delay.
func bufferLevels(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		fan := edakit.Int(3, 10)
		n := edakit.Int(50, 20000)
		l := levelsFor(n, fan)
		return edakit.Numeric(edakit.Sprintf("A clock tree must reach %d flops, and no buffer may drive more than %d loads (flops or buffers). The tree is built as a balanced fanout tree from one root buffer. How many buffer levels does it need?", n, fan), edakit.IntStr(l),
			edakit.Sprintf("$L$ levels reach at most $%d^L$ sinks, so $L = \\lceil \\log_{%d} %d \\rceil = %d$ ($%d^{%d} = %d \\ge %d$).", fan, fan, n, l, fan, l, int(math.Pow(float64(fan), float64(l))), n),
			map[string]string{"levels": edakit.IntStr(l)})
	case 2:
		fan := edakit.Int(3, 8)
		n := edakit.Int(30, 3000)
		total, cur := 0, n
		var counts []int
		for cur > 1 {
			cur = (cur + fan - 1) / fan
			counts = append(counts, cur)
			total += cur
		}
		return edakit.Numeric(edakit.Sprintf("A clock tree reaches %d flops with buffers that drive at most %d loads each. Working up from the flops, each level uses as few buffers as possible, until a single root buffer remains. How many buffers does the tree contain in total, root included?", n, fan), edakit.IntStr(total),
			edakit.Sprintf("Level counts from the leaves up are $\\lceil %d/%d \\rceil$ and so on: %s, a total of %d buffers.", n, fan, ints(counts), total),
			map[string]string{"buffers": edakit.IntStr(total)})
	}
	n := edakit.Pick([]int{256, 500, 1000, 2000, 4096, 10000, 50000})
	a := edakit.Int(8, 20)
	b := edakit.Int(2, 6)
	best, bf := math.MaxInt32, 0
	var rows []string
	for fan := 2; fan <= 8; fan++ {
		l := levelsFor(n, fan)
		lat := l * (a + b*fan)
		rows = append(rows, edakit.Sprintf("fanout %d: %d levels, %d ps", fan, l, lat))
		if lat < best {
			best, bf = lat, fan
		}
	}
	return edakit.Numeric(edakit.Sprintf("A clock tree must reach %d sinks. A buffer driving $F$ loads has delay $%d + %d F$ ps. The designer may choose any uniform fanout $F$ from 2 to 8, with $\\lceil \\log_F %d \\rceil$ levels. What is the smallest achievable latency, in ps?", n, a, b, n), edakit.IntStr(best),
		edakit.Sprintf("Latency is levels times stage delay: %s. The minimum, %d ps, is at fanout %d. Small fanout needs many levels; large fanout makes each stage slow.", strings.Join(rows, "; "), best, bf),
		map[string]string{"best_fanout": edakit.IntStr(bf)})
}

// slewWireLength: Elmore at the end of a buffer-driven wire, the ln 9 slew estimate, and the
// longest wire that meets a slew limit.
func slewWireLength(ctx generator.GeneratorContext) generator.Problem {
	rb := float64(50 * edakit.Int(2, 8))
	r := pick(0.5, 2.0, 0.1)
	c := pick(0.15, 0.25, 0.01)
	cin := pick(2, 8, 0.5)
	desc := edakit.Sprintf("A clock buffer with output resistance $%s\\ \\Omega$ drives a wire of $r = %s\\ \\Omega$/µm and $c = %s$ fF/µm into the next buffer's input of %s fF.", f(rb, 0), f(r, 1), f(c, 2), f(cin, 1))
	el := func(L float64) float64 { return (rb*(c*L+cin) + r*L*(c*L/2+cin)) / 1000 }
	switch kpOf(ctx) {
	case 1:
		L := float64(50 * edakit.Int(2, 12))
		e := el(L)
		return edakit.Numeric(desc+edakit.Sprintf(" The wire is %s µm long. What is the Elmore delay from the buffer output to the far end, in ps? Round to 2 decimal places.", f(L, 0)), f(e, 2),
			edakit.Sprintf("$R_b(cL + C_{in}) + rL(cL/2 + C_{in}) = %s(%s) + %s(%s)$ ohm-fF $= %s$ ps.", f(rb, 0), f(c*L+cin, 2), f(r*L, 1), f(c*L/2+cin, 2), f(e, 2)),
			map[string]string{"elmore": f(e, 3)})
	case 2:
		L := float64(50 * edakit.Int(2, 12))
		e := el(L)
		s := ln9 * e
		return edakit.Numeric(desc+edakit.Sprintf(" The wire is %s µm long. Using the single-pole estimate (10-90%% slew $\\approx \\ln 9$ times the Elmore delay), what is the slew at the far end, in ps? Round to 1 decimal place.", f(L, 0)), f(s, 1),
			edakit.Sprintf("Elmore delay: %s ps. Slew $\\approx 2.197 \\cdot %s = %s$ ps.", f(e, 3), f(e, 3), f(s, 1)),
			map[string]string{"elmore": f(e, 3), "slew": f(s, 2)})
	}
	sl := float64(10 * edakit.Int(5, 15))
	a, b, cc := r*c/2, rb*c+r*cin, rb*cin-1000*sl/ln9
	if cc >= 0 { // unreachable for the drawn ranges; keep the instance valid anyway
		sl = 10 * math.Ceil((rb*cin*ln9/1000+10)/10)
		cc = rb*cin - 1000*sl/ln9
	}
	L := (-b + math.Sqrt(b*b-4*a*cc)) / (2 * a)
	return edakit.Numeric(desc+edakit.Sprintf(" The clock slew limit is %s ps, estimated as $\\ln 9$ times the Elmore delay. What is the longest wire the buffer may drive before the next buffer must be inserted, in µm? Round to the nearest whole number.", f(sl, 0)), f(L, 0),
		edakit.Sprintf("Require Elmore $\\le %s/2.197 = %s$ ps. That is the quadratic $(rc/2)L^2 + (R_b c + r C_{in})L + R_b C_{in} - %s = 0$ in ohm-fF, with coefficients %s, %s, %s. The positive root is %s µm.", f(sl, 0), f(sl/ln9, 2), f(1000*sl/ln9, 0), f(a, 3), f(b, 2), f(cc, 1), f(L, 1)),
		map[string]string{"lmax": f(L, 1)})
}

// dutyCycle: pulse-width change through buffers and inverters with unequal rise and fall delays.
func dutyCycle(ctx generator.GeneratorContext) generator.Problem {
	T := float64(50 * edakit.Int(8, 40))
	W := T / 2
	lh := float64(edakit.Int(25, 60))
	hl := lh + float64(edakit.Int(-8, 8))
	if hl == lh {
		hl = lh + 3
	}
	n := edakit.Scaled(ctx, 4, 16)
	bufDesc := edakit.Sprintf("Each clock buffer delays a rising output by %s ps and a falling output by %s ps.", f(lh, 0), f(hl, 0))
	switch kpOf(ctx) {
	case 1:
		w := W + float64(n)*(hl-lh)
		return edakit.Numeric(edakit.Sprintf("A %s ps clock enters a chain of %d identical buffers with a high pulse of %s ps. %s What is the width of the high pulse at the chain output, in ps?", f(T, 0), n, f(W, 0), bufDesc), f(w, 0),
			edakit.Sprintf("Each buffer delays the rising edge by %s and the falling edge by %s, so the high pulse changes by $t_{HL} - t_{LH} = %s$ ps per stage. After %d stages: $%s + %d(%s) = %s$ ps.", f(lh, 0), f(hl, 0), f(hl-lh, 0), n, f(W, 0), n, f(hl-lh, 0), f(w, 0)),
			map[string]string{"width": f(w, 0)})
	case 2:
		w := W + float64(n)*(hl-lh)
		dc := 100 * w / T
		return edakit.Numeric(edakit.Sprintf("A %s ps clock with a 50%% duty cycle passes %d identical buffers. %s What is the duty cycle at the output, in percent? Round to 2 decimal places.", f(T, 0), n, bufDesc), f(dc, 2),
			edakit.Sprintf("The high pulse becomes $%s + %d(%s - %s) = %s$ ps, so the duty cycle is $%s/%s = %s\\%%$.", f(W, 0), n, f(hl, 0), f(lh, 0), f(w, 0), f(w, 0), f(T, 0), f(dc, 2)),
			map[string]string{"duty": f(dc, 3)})
	}
	ilh := float64(edakit.Int(15, 40))
	ihl := ilh + float64(edakit.Int(-8, 8))
	if ihl == ilh {
		ihl = ilh - 4
	}
	k := edakit.Int(3, 9)
	w := W + float64(n)*(hl-lh)
	odd := k%2 == 1
	if odd {
		w += ilh - ihl
	}
	q := edakit.Sprintf("A %s ps clock with a %s ps high pulse passes %d buffers and then %d inverters. %s Each inverter has a rising-output delay of %s ps and a falling-output delay of %s ps. What is the width at the output of the pulse that began as the input's high pulse (it ends up low if the chain inverts), in ps?", f(T, 0), f(W, 0), n, k, bufDesc, f(ilh, 0), f(ihl, 0))
	why := edakit.Sprintf("The buffers change the pulse by %d(%s) = %s ps. Inverters alternate which output edge is rising, so each consecutive pair cancels exactly. ", n, f(hl-lh, 0), f(float64(n)*(hl-lh), 0))
	if odd {
		why += edakit.Sprintf("With %d inverters, one is left over: it turns the high pulse into a low pulse whose width changes by $t_{LH} - t_{HL} = %s$ ps. The result is %s ps.", k, f(ilh-ihl, 0), f(w, 0))
	} else {
		why += edakit.Sprintf("With %d inverters, all pairs cancel, so the result is %s ps.", k, f(w, 0))
	}
	return edakit.Numeric(q, f(w, 0), why, map[string]string{"width": f(w, 0)})
}

// latencyBalancing: balance clock groups, including a macro's internal insertion delay.
func latencyBalancing(ctx generator.GeneratorContext) generator.Problem {
	names := []string{"core flops", "IO flops", "RAM macro", "divider flops"}
	pin := []int{edakit.Int(300, 420), edakit.Int(250, 400), edakit.Int(150, 300), edakit.Int(250, 400)}
	internal := 10 * edakit.Int(6, 18)
	eff := append([]int(nil), pin...)
	eff[2] += internal
	target := maxInt(eff)
	d := 5 * edakit.Int(4, 8)
	tol := 5 * edakit.Int(1, 4)
	need := func(l int) int {
		gap := target - tol - l
		if gap <= 0 {
			return 0
		}
		return (gap + d - 1) / d
	}
	var desc []string
	for i, nm := range names {
		desc = append(desc, edakit.Sprintf("%s %d ps", nm, pin[i]))
	}
	q := edakit.Sprintf("Before balancing, the clock pins of four sink groups see latencies of %s. The RAM macro also has %d ps of internal insertion delay from its clock pin to its internal registers. ", strings.Join(desc, ", "), internal)
	switch kpOf(ctx) {
	case 1:
		sk := maxInt(eff) - minInt(eff)
		return edakit.Numeric(q+"Counting the macro's internal delay, what is the skew across the four groups' registers, in ps?", edakit.IntStr(sk),
			edakit.Sprintf("The macro's registers see %d + %d = %d ps. The register-level latencies are %s, so skew is %d - %d = %d ps. Balancing pin latencies alone would miss the internal delay.", pin[2], internal, eff[2], ints(eff), maxInt(eff), minInt(eff), sk),
			map[string]string{"skew": edakit.IntStr(sk)})
	case 2:
		lo := 0
		for i := range eff {
			if eff[i] < eff[lo] {
				lo = i
			}
		}
		n := need(eff[lo])
		return edakit.Numeric(q+edakit.Sprintf("CTS balances register-level latency to the largest one, within a tolerance of %d ps below it, by inserting delay cells of %d ps. How many delay cells does the earliest group (%s) need?", tol, d, names[lo]), edakit.IntStr(n),
			edakit.Sprintf("Target %d ps, so the group must reach at least %d ps. It is at %d ps; the gap %d ps needs $\\lceil %d/%d \\rceil = %d$ cells.", target, target-tol, eff[lo], target-tol-eff[lo], target-tol-eff[lo], d, n),
			map[string]string{"cells": edakit.IntStr(n)})
	}
	total := 0
	var per []string
	for i := range eff {
		n := need(eff[i])
		total += n
		per = append(per, edakit.Sprintf("%s %d", names[i], n))
	}
	return edakit.Numeric(q+edakit.Sprintf("CTS balances register-level latency to the largest one, within a tolerance of %d ps below it, by inserting delay cells of %d ps on the early groups. How many delay cells are inserted in total?", tol, d), edakit.IntStr(total),
		edakit.Sprintf("Register-level latencies %s, target %d ps, floor %d ps. Cells per group: %s, total %d.", ints(eff), target, target-tol, strings.Join(per, ", "), total),
		map[string]string{"cells": edakit.IntStr(total)})
}
