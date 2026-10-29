package cts

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// usefulSkew: the feasible window for the clock arrival s at a flop F between A and B.
func usefulSkew(ctx generator.GeneratorContext) generator.Problem {
	for attempt := 0; ; attempt++ {
		T := float64(50 * edakit.Int(10, 20))
		tcq, tsu, th := float64(edakit.Int(40, 70)), float64(edakit.Int(20, 40)), float64(edakit.Int(10, 30))
		dmaxA := T - tcq - tsu + float64(edakit.Int(10, 60))
		dmaxB := T - tcq - tsu - float64(edakit.Int(80, 200))
		dminA := float64(edakit.Int(40, 200))
		dminB := float64(edakit.Int(20, 150))
		if attempt >= generator.MaxDrawAttempts {
			T, tcq, tsu, th, dmaxA, dminA, dmaxB, dminB = 1000, 43, 28, 21, 945, 158, 754, 24
		}
		L1 := tcq + dmaxA + tsu - T
		U1 := T - tcq - dmaxB - tsu
		U2 := tcq + dminA - th
		L2 := th - tcq - dminB
		lo, hi := math.Max(L1, L2), math.Min(U1, U2)
		if hi-lo < 5 && attempt < generator.MaxDrawAttempts {
			continue
		}
		desc := edakit.Sprintf("Flop F captures from flop A and launches to flop B. The clock reaches A and B at 0 ps and reaches F at $s$ ps, which CTS may choose. Period %s ps, clock-to-Q %s ps, setup %s ps, hold %s ps. Path A to F has a longest delay of %s ps and a shortest of %s ps; path F to B has a longest delay of %s ps and a shortest of %s ps.",
			f(T, 0), f(tcq, 0), f(tsu, 0), f(th, 0), f(dmaxA, 0), f(dminA, 0), f(dmaxB, 0), f(dminB, 0))
		facts := map[string]string{"lo": f(lo, 0), "hi": f(hi, 0)}
		switch kpOf(ctx) {
		case 1:
			return edakit.Numeric(desc+" With $s = 0$ the path A to F fails setup. What is the smallest $s$ that makes it meet setup, in ps?", f(L1, 0),
				edakit.Sprintf("Setup into F: $t_{cq} + d_{max} + t_{su} \\le s + T$, so $s \\ge %s + %s + %s - %s = %s$ ps. Delaying F's clock lends the path that much time.", f(tcq, 0), f(dmaxA, 0), f(tsu, 0), f(T, 0), f(L1, 0)), facts)
		case 2:
			return edakit.Numeric(desc+" Considering setup on F to B and hold on A to F, what is the largest $s$ allowed, in ps?", f(hi, 0),
				edakit.Sprintf("Setup out of F: $s + t_{cq} + d_{max} + t_{su} \\le T$ gives $s \\le %s$. Hold into F: $t_{cq} + d_{min} \\ge s + t_h$ gives $s \\le %s$. The tighter is %s ps.", f(U1, 0), f(U2, 0), f(hi, 0)), facts)
		}
		return edakit.Numeric(desc+" Taking all four setup and hold constraints into account, how wide is the window of feasible $s$, in ps?", f(hi-lo, 0),
			edakit.Sprintf("Lower bounds: setup into F gives %s and hold out of F ($s + t_{cq} + d_{min} \\ge t_h$) gives %s. Upper bounds: setup out of F gives %s and hold into F gives %s. The window is $[%s, %s]$, %s ps wide.", f(L1, 0), f(L2, 0), f(U1, 0), f(U2, 0), f(lo, 0), f(hi, 0), f(hi-lo, 0)), facts)
	}
}

// differenceConstraints: cycle weight, Bellman-Ford potentials, and the minimum period of a ring.
func differenceConstraints(ctx generator.GeneratorContext) generator.Problem {
	k := kpOf(ctx)
	if k == 2 {
		names := []string{"A", "B", "C", "D"}
		type edge struct {
			from, to int
			w        float64
		}
		for attempt := 0; ; attempt++ {
			var es []edge
			m := edakit.Scaled(ctx, 4, 6)
			pairs := edakit.Shuffle([][2]int{{0, 1}, {0, 2}, {0, 3}, {1, 0}, {1, 2}, {1, 3}, {2, 0}, {2, 1}, {2, 3}, {3, 0}, {3, 1}, {3, 2}})
			for _, pr := range pairs[:m] {
				es = append(es, edge{pr[0], pr[1], float64(5 * edakit.Int(-10, 12))})
			}
			last := attempt >= generator.MaxDrawAttempts
			if last { // known-good instance: no negative cycle, t_D = -15
				es = []edge{{2, 3, 5}, {3, 2, 30}, {0, 3, -15}, {0, 2, 20}}
			}
			dist := []float64{0, 0, 0, 0}
			for it := 0; it < 4; it++ {
				for _, e := range es {
					if dist[e.from]+e.w < dist[e.to] {
						dist[e.to] = dist[e.from] + e.w
					}
				}
			}
			neg := false
			for _, e := range es {
				if dist[e.from]+e.w < dist[e.to] {
					neg = true
				}
			}
			if neg && !last {
				continue
			}
			ask := edakit.Int(0, 3)
			if last {
				ask = 3
			} else if dist[ask] == 0 && edakit.Int(0, 2) > 0 {
				continue
			}
			var cs []string
			for _, e := range es {
				cs = append(cs, edakit.Sprintf("$t_%s - t_%s \\le %s$", names[e.to], names[e.from], f(e.w, 0)))
			}
			q := "Clock arrival times must satisfy " + strings.Join(cs, ", ") + " (ps). Each constraint $t_j - t_i \\le w$ is an edge from $i$ to $j$ of weight $w$. Bellman-Ford from a virtual source joined to every node by a zero-weight edge assigns each $t$ its shortest-path distance. What value does it give $t_" + names[ask] + "$, in ps?"
			return edakit.Numeric(q, f(dist[ask], 0),
				edakit.Sprintf("Start every distance at 0 (the source edges) and relax until nothing changes: $t$ = (%s) for (A, B, C, D). No edge can still be relaxed, so there is no negative cycle and these distances satisfy every constraint.", floats(dist, 0)),
				map[string]string{"t": f(dist[ask], 0)})
		}
	}
	n := edakit.Scaled(ctx, 2, 4)
	tcq, tsu := float64(edakit.Int(40, 70)), float64(edakit.Int(20, 40))
	var ds []float64
	sum := 0.0
	for i := 0; i < n; i++ {
		d := float64(10 * edakit.Int(40, 110))
		ds = append(ds, d)
		sum += tcq + d + tsu
	}
	ring := edakit.Sprintf("%d flops form a ring: each launches into the next and the last into the first. The longest path delays around the ring are %s ps; clock-to-Q is %s ps and setup %s ps for every flop.", n, floats(ds, 0), f(tcq, 0), f(tsu, 0))
	if k == 1 {
		T := math.Round(sum/float64(n)/10)*10 + float64(10*edakit.Int(-6, 6))
		w := float64(n)*T - sum
		return edakit.Numeric(ring+edakit.Sprintf(" At a period of %s ps, each setup constraint reads $t_i - t_j \\le T - t_{cq} - d_{ij} - t_{su}$. What is the total weight of the ring's cycle in the constraint graph, in ps?", f(T, 0)), f(w, 0),
			edakit.Sprintf("Summing the constraint weights around the ring cancels every $t$: $%d \\cdot %s - %s = %s$ ps. %s", n, f(T, 0), f(sum, 0), f(w, 0), map[bool]string{true: "It is negative, a negative cycle, so no clock schedule meets setup at this period.", false: "It is non-negative, so a feasible skew schedule exists for this ring."}[w < 0]),
			map[string]string{"cycle_weight": f(w, 0)})
	}
	tmin := sum / float64(n)
	return edakit.Numeric(ring+" With unconstrained useful skew, what is the minimum clock period at which setup can be met on every path, in ps? Round to 1 decimal place.", f(tmin, 1),
		edakit.Sprintf("The ring's cycle weight $nT - \\sum(t_{cq} + d + t_{su})$ must be non-negative, so $T \\ge %s/%d = %s$ ps: the mean stage delay around the cycle. Skew can redistribute time around a loop but cannot create it.", f(sum, 0), n, f(tmin, 1)),
		map[string]string{"tmin": f(tmin, 2)})
}

// schedulingLP: maximise the minimum slack M over clock arrivals.
func schedulingLP(ctx generator.GeneratorContext) generator.Problem {
	k := kpOf(ctx)
	if k == 1 {
		n := edakit.Scaled(ctx, 2, 4)
		tcq, tsu := float64(edakit.Int(40, 70)), float64(edakit.Int(20, 40))
		var ds []float64
		sum := 0.0
		for i := 0; i < n; i++ {
			d := float64(10 * edakit.Int(40, 110))
			ds = append(ds, d)
			sum += tcq + d + tsu
		}
		T := math.Ceil(sum/float64(n)/50)*50 + float64(50*edakit.Int(0, 3))
		m := T - sum/float64(n)
		return edakit.Numeric(edakit.Sprintf("%d flops form a ring, each launching into the next. Longest path delays are %s ps, clock-to-Q %s ps, setup %s ps, period %s ps. Ignoring hold, the scheduler chooses clock arrivals to maximise $M$, the minimum setup slack over the ring's paths. What is the optimal $M$, in ps? Round to 2 decimal places.", n, floats(ds, 0), f(tcq, 0), f(tsu, 0), f(T, 0)), f(m, 2),
			edakit.Sprintf("Each path needs $t_i - t_j \\le T - S_{ij} - M$. Summed around the ring, $0 \\le nT - \\sum S - nM$, so $M \\le T - \\sum S/n = %s - %s/%d = %s$ ps, and a schedule reaches it by spreading the time evenly.", f(T, 0), f(sum, 0), n, f(m, 2)),
			map[string]string{"m": f(m, 3)})
	}
	for attempt := 0; ; attempt++ {
		T := float64(50 * edakit.Int(10, 20))
		tcq, tsu, th := float64(edakit.Int(40, 70)), float64(edakit.Int(20, 40)), float64(edakit.Int(10, 30))
		dAB := T - tcq - tsu + float64(edakit.Int(-60, 60))
		dBA := T - tcq - tsu - float64(edakit.Int(50, 250))
		mAB := float64(edakit.Int(40, 250))
		mBA := float64(edakit.Int(40, 250))
		if attempt >= generator.MaxDrawAttempts {
			T, tcq, tsu, th, dAB, mAB, dBA, mBA = 800, 61, 21, 29, 736, 106, 622, 97
		}
		sAB, sBA := tcq+dAB+tsu, tcq+dBA+tsu
		hAB, hBA := tcq+mAB-th, tcq+mBA-th
		P := math.Min(T-sAB, hBA)
		Q := math.Min(T-sBA, hAB)
		M := (P + Q) / 2
		u := (P - Q) / 2
		if M <= 0 && attempt < generator.MaxDrawAttempts {
			continue
		}
		desc := edakit.Sprintf("Flops A and B form a loop: A launches to B and B launches back to A. Period %s ps, clock-to-Q %s ps, setup %s ps, hold %s ps. A to B has longest delay %s ps and shortest %s ps; B to A has longest %s ps and shortest %s ps. The scheduler chooses clock arrivals $t_A$ and $t_B$ to maximise $M$, the minimum of all four setup and hold slacks.", f(T, 0), f(tcq, 0), f(tsu, 0), f(th, 0), f(dAB, 0), f(mAB, 0), f(dBA, 0), f(mBA, 0))
		why := edakit.Sprintf("With $u = t_A - t_B$: setup A to B and hold B to A bound $u \\le \\min(%s, %s) - M = %s - M$; setup B to A and hold A to B bound $-u \\le \\min(%s, %s) - M = %s - M$. Adding gives $M \\le (%s + %s)/2 = %s$ ps, reached at $u = (%s - %s)/2 = %s$ ps.", f(T-sAB, 0), f(hBA, 0), f(P, 0), f(T-sBA, 0), f(hAB, 0), f(Q, 0), f(P, 0), f(Q, 0), f(M, 1), f(P, 0), f(Q, 0), f(u, 1))
		if k == 2 {
			return edakit.Numeric(desc+" What is the optimal $M$, in ps? Round to 1 decimal place.", f(M, 1), why, map[string]string{"m": f(M, 2), "u": f(u, 2)})
		}
		return edakit.Numeric(desc+" At the optimum, what is $t_A - t_B$, in ps? Round to 1 decimal place.", f(u, 1), why, map[string]string{"m": f(M, 2), "u": f(u, 2)})
	}
}

// icgGlitch: computed choice on the runt pulse an AND-only gate passes.
func icgGlitch(ctx generator.GeneratorContext) generator.Problem {
	T := float64(100 * edakit.Int(5, 20))
	h := T / 2
	te := float64(10 * edakit.Int(2, int(h/10)-2))
	rise := edakit.Int(0, 1) == 1
	var w float64
	var event string
	if rise {
		w = h - te
		event = "rises"
	} else {
		w = te
		event = "falls"
	}
	c, ds := numOptions(w, 0, " ps", h, h-w, T-te)
	q := edakit.Sprintf("A %s ps clock is high from 0 to %s ps in each cycle. A designer gates it with a plain AND gate (no latch). In one cycle the enable %s at %s ps, while the clock is high. What is the width of the high pulse the gated clock produces in that cycle?", f(T, 0), f(h, 0), event, f(te, 0))
	return edakit.Choice(q, c,
		edakit.Sprintf("The AND output is high only while both inputs are high. The enable %s at %s ps, so the output is high for %s ps of the %s ps phase: a runt pulse that can clock some flops and not others. A latch-based ICG samples the enable while the clock is low and holds it through the high phase, so no partial pulse appears.", event, f(te, 0), f(w, 0), f(h, 0)), ds...)
}

// gatingSavings: clock-gating power with enable probability, ICG overhead and break-even width.
func gatingSavings(ctx generator.GeneratorContext) generator.Problem {
	v := pick(0.7, 1.0, 0.05)
	fr := pick(0.5, 2.5, 0.25)
	pe := pick(0.1, 0.6, 0.05)
	switch kpOf(ctx) {
	case 1:
		cg := float64(5 * edakit.Int(4, 40))
		p := pe * cg * v * v * fr
		return edakit.Numeric(edakit.Sprintf("A clock gate drives %s pF of clock subtree. The enable is true in %s%% of cycles. At %s GHz and %s V, what is the average power of the gated subtree, in mW? Round to 2 decimal places.", f(cg, 0), f(100*pe, 0), f(fr, 2), f(v, 2)), f(p, 2),
			edakit.Sprintf("The subtree toggles only in enabled cycles: $P = p_{en} C V^2 f = %s \\cdot %s \\cdot %s^2 \\cdot %s = %s$ mW.", f(pe, 2), f(cg, 0), f(v, 2), f(fr, 2), f(p, 2)),
			map[string]string{"p_mw": f(p, 3)})
	case 2:
		cg := float64(5 * edakit.Int(4, 40))
		n := edakit.Int(10, 200)
		ci := pick(2, 6, 0.5)
		save := (1-pe)*cg*v*v*fr - float64(n)*ci/1000*v*v*fr
		return edakit.Numeric(edakit.Sprintf("Synthesis inserts %d clock gates that together gate %s pF of clock subtree. Each gate presents %s fF of switched capacitance on the ungated clock, which toggles every cycle. The enables are true in %s%% of cycles. At %s GHz and %s V, what is the net clock power saved, in mW? Round to 2 decimal places.", n, f(cg, 0), f(ci, 1), f(100*pe, 0), f(fr, 2), f(v, 2)), f(save, 2),
			edakit.Sprintf("Saved: $(1 - %s) \\cdot %s$ pF of gated load. Added: $%d \\cdot %s$ fF $= %s$ pF of gate input that always toggles. Net $= (%s - %s) \\cdot %s^2 \\cdot %s = %s$ mW.", f(pe, 2), f(cg, 0), n, f(ci, 1), f(float64(n)*ci/1000, 3), f((1-pe)*cg, 3), f(float64(n)*ci/1000, 3), f(v, 2), f(fr, 2), f(save, 2)),
			map[string]string{"save_mw": f(save, 3)})
	}
	cp := pick(0.6, 1.5, 0.1)
	ci := pick(2, 8, 0.5)
	n := int(math.Floor(ci/((1-pe)*cp))) + 1
	return edakit.Numeric(edakit.Sprintf("A clock gate costs %s fF of always-toggling capacitance on the ungated clock. Each flop it gates has a %s fF clock pin, and the enable is true in %s%% of cycles. Ignoring wire, what is the smallest number of flops sharing the gate for which gating saves power?", f(ci, 1), f(cp, 1), f(100*pe, 0)), edakit.IntStr(n),
		edakit.Sprintf("Gating $n$ flops removes $(1 - p_{en})\\,n\\,C_{pin} = %s\\,n$ fF of average switched load and adds %s fF. Savings need $n > %s/%s = %s$, so $n = %d$. This is why synthesis gates register banks, not single flops.", f((1-pe)*cp, 3), f(ci, 1), f(ci, 1), f((1-pe)*cp, 3), f(ci/((1-pe)*cp), 2), n),
		map[string]string{"n": edakit.IntStr(n)})
}

// enableTiming: the ICG enable setup check and the gated-subtree delay that eats into it.
func enableTiming(ctx generator.GeneratorContext) generator.Problem {
	for attempt := 0; ; attempt++ {
		T := float64(50 * edakit.Int(8, 20))
		lf := float64(10 * edakit.Int(30, 50))
		dg := float64(10 * edakit.Int(5, 20))
		tcq, tsu := float64(edakit.Int(40, 70)), float64(edakit.Int(30, 60))
		de := float64(10 * edakit.Int(20, int(T/10)-10))
		if attempt >= generator.MaxDrawAttempts {
			T, lf, dg, tcq, tsu, de = 650, 360, 180, 61, 54, 350
		}
		slack := T + (lf - dg) - tsu - (lf + tcq + de)
		desc := edakit.Sprintf("A clock gate drives a subtree whose delay from the gate output to its flops is %s ps, and those flops have a clock latency of %s ps. The enable is launched by a flop whose clock latency is also %s ps, with clock-to-Q %s ps and %s ps of enable logic. The gate's enable setup requirement is %s ps before the rising edge at its clock pin, and the period is %s ps.", f(dg, 0), f(lf, 0), f(lf, 0), f(tcq, 0), f(de, 0), f(tsu, 0), f(T, 0))
		switch kpOf(ctx) {
		case 1:
			return edakit.Numeric(desc+" What is the setup slack on the enable, in ps?", f(slack, 0),
				edakit.Sprintf("The gate's clock pin sits upstream of its subtree, at $%s - %s = %s$ ps. Slack $= T + %s - t_{su} - (%s + t_{cq} + d_{en}) = %s$ ps. The gate sees its clock %s ps earlier than an ordinary flop would, which costs the enable path that much.", f(lf, 0), f(dg, 0), f(lf-dg, 0), f(lf-dg, 0), f(lf, 0), f(slack, 0), f(dg, 0)),
				map[string]string{"slack": f(slack, 0)})
		case 2:
			dmax := T - tsu - tcq - de
			if dmax <= 0 && attempt < generator.MaxDrawAttempts {
				continue
			}
			return edakit.Numeric(desc+" Keeping everything else fixed, what is the largest gated-subtree delay for which the enable still meets setup, in ps?", f(dmax, 0),
				edakit.Sprintf("Slack $= T - t_{su} - t_{cq} - d_{en} - D_g$, so $D_g \\le %s - %s - %s - %s = %s$ ps.", f(T, 0), f(tsu, 0), f(tcq, 0), f(de, 0), f(dmax, 0)),
				map[string]string{"dg_max": f(dmax, 0)})
		}
		dg2 := float64(10 * edakit.Int(1, int(dg/10)))
		if dg2 >= dg {
			dg2 = dg - 10
		}
		ns := slack + dg - dg2
		return edakit.Numeric(desc+edakit.Sprintf(" CTS clones the gate and moves the copies next to their flops, cutting the gated-subtree delay to %s ps. What is the new enable setup slack, in ps?", f(dg2, 0)), f(ns, 0),
			edakit.Sprintf("The old slack was %s ps. The gate's clock now arrives $%s - %s = %s$ ps later, so the slack becomes %s ps. Cloning trades clock-gate count and power for enable timing.", f(slack, 0), f(dg, 0), f(dg2, 0), f(dg-dg2, 0), f(ns, 0)),
			map[string]string{"slack": f(ns, 0)})
	}
}
