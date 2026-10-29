package sta

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func waveform(ctx generator.GeneratorContext) generator.Problem {
	p := edakit.Int(4, 20)
	r := edakit.Int(0, p/2)
	f := edakit.Int(r+1, p-1)
	cmd := edakit.Sprintf("create_clock -name clk -period %d -waveform {%d %d} [get_ports clk]", p, r, f)
	q := "The SDC command `" + cmd + "` defines a clock in ns. "
	switch ctx.KPOrdinal {
	case 1:
		duty := 100 * float64(f-r) / float64(p)
		ans := edakit.Num(duty, 1)
		return edakit.Numeric(q+"What is its duty cycle, in percent? Round to 1 decimal place.", ans,
			edakit.Sprintf("The clock is high from %d to %d ns, %d ns of every %d ns: %s%%.", r, f, f-r, p, ans), map[string]string{"duty": ans})
	case 3:
		t0 := edakit.Int(p, 4*p)
		// first falling edge strictly after t0
		k := 0
		for f+k*p <= t0 {
			k++
		}
		ans := f + k*p
		return edakit.Numeric(q+edakit.Sprintf("What is the time of the first falling edge after %d ns, in ns?", t0), edakit.IntStr(ans),
			edakit.Sprintf("Falling edges occur at %d + k x %d ns. The first one after %d ns is k = %d, at %d ns.", f, p, t0, k, ans), map[string]string{"edge": edakit.IntStr(ans)})
	}
	n := edakit.Int(2, 6)
	ans := r + (n-1)*p
	return edakit.Numeric(q+edakit.Sprintf("Counting the rising edge at %d ns as the first, at what time is rising edge number %d, in ns?", r, n), edakit.IntStr(ans),
		edakit.Sprintf("Rising edges occur at %d + (n - 1) x %d ns, so edge %d is at %d ns.", r, p, n, ans), map[string]string{"edge": edakit.IntStr(ans)})
}

// clockTree draws a shared trunk and two branches to a launch and a capture flop.
type clockTree struct {
	src             int
	trunk, la, ca   []int
	launch, capture int
}

func drawClockTree() clockTree {
	ct := clockTree{src: 10 * edakit.Int(0, 30)}
	for i := 0; i < edakit.Int(1, 2); i++ {
		ct.trunk = append(ct.trunk, edakit.Int(20, 60))
	}
	for i := 0; i < edakit.Int(1, 3); i++ {
		ct.la = append(ct.la, edakit.Int(20, 60))
	}
	for i := 0; i < edakit.Int(1, 3); i++ {
		ct.ca = append(ct.ca, edakit.Int(20, 60))
	}
	sum := func(xs []int) int {
		s := 0
		for _, x := range xs {
			s += x
		}
		return s
	}
	ct.launch = ct.src + sum(ct.trunk) + sum(ct.la)
	ct.capture = ct.src + sum(ct.trunk) + sum(ct.ca)
	return ct
}

func ints(xs []int) string {
	var p []string
	for _, x := range xs {
		p = append(p, edakit.IntStr(x))
	}
	return strings.Join(p, ", ")
}

func (ct clockTree) describe() string {
	return edakit.Sprintf("A clock has source latency %d ps. Its tree first passes through shared buffers with delays %s ps. One branch then passes buffers of %s ps to the launch flop FF1; the other passes buffers of %s ps to the capture flop FF2.",
		ct.src, ints(ct.trunk), ints(ct.la), ints(ct.ca))
}

func latency(ctx generator.GeneratorContext) generator.Problem {
	ct := drawClockTree()
	skew := ct.capture - ct.launch
	facts := map[string]string{"launch": edakit.IntStr(ct.launch), "capture": edakit.IntStr(ct.capture), "skew": edakit.IntStr(skew)}
	switch ctx.KPOrdinal {
	case 1:
		who, v := "FF1", ct.launch
		if edakit.Int(0, 1) == 1 {
			who, v = "FF2", ct.capture
		}
		return edakit.Numeric(ct.describe()+edakit.Sprintf(" With propagated clocks, when does the clock edge launched at 0 ps reach %s, in ps?", who), edakit.IntStr(v),
			edakit.Sprintf("Clock arrival is source latency plus network latency, the sum of buffer delays along the path to %s: %d ps.", who, v), facts)
	case 3:
		t := 50 * edakit.Int(10, 20)
		tcq, d, su := edakit.Int(30, 70), 10*edakit.Int(30, 80), edakit.Int(20, 50)
		s := t + ct.capture - su - (ct.launch + tcq + d)
		return edakit.Numeric(ct.describe()+edakit.Sprintf(" The period is %d ps, clock-to-Q %d ps, the data path %d ps and setup %d ps. With propagated clocks, what is the setup slack from FF1 to FF2, in ps?", t, tcq, d, su), edakit.IntStr(s),
			edakit.Sprintf("Launch clock %d ps, capture clock %d ps. Slack = %d + %d - %d - (%d + %d + %d) = %d ps.", ct.launch, ct.capture, t, ct.capture, su, ct.launch, tcq, d, s), map[string]string{"setup_slack": edakit.IntStr(s)})
	}
	return edakit.Numeric(ct.describe()+" What is the clock skew from FF1 to FF2 (capture arrival minus launch arrival), in ps?", edakit.IntStr(skew),
		edakit.Sprintf("The shared trunk and source latency cancel. Skew = %d - %d = %d ps.", ct.capture, ct.launch, skew), facts)
}

func skewJitter(ctx generator.GeneratorContext) generator.Problem {
	p := drawRegPath()
	skew := p.lCap - p.lLaunch
	j := edakit.Int(5, 40)
	q := p.describe()
	facts := map[string]string{"skew": edakit.IntStr(skew)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+" What is the clock skew of this path (capture arrival minus launch arrival), in ps?", edakit.IntStr(skew),
			edakit.Sprintf("Skew = %d - %d = %d ps. Positive skew gives setup more time and hold less.", p.lCap, p.lLaunch, skew), facts)
	case 3:
		hold := p.tcq + p.dmin - p.th - skew
		facts["hold_slack"] = edakit.IntStr(hold)
		return edakit.Numeric(q+edakit.Sprintf(" The clock also has %d ps of cycle-to-cycle jitter. What is the hold slack, in ps?", j), edakit.IntStr(hold),
			edakit.Sprintf("Hold slack = $t_{cq} + d_{min} - t_h - skew$ = %d + %d - %d - (%d) = %d ps. Cycle-to-cycle jitter does not enter a same-edge hold check.", p.tcq, p.dmin, p.th, skew, hold), facts)
	}
	s := p.t - j + skew - p.tsu - p.tcq - p.dmax
	facts["setup_slack"] = edakit.IntStr(s)
	return edakit.Numeric(q+edakit.Sprintf(" The clock has %d ps of cycle-to-cycle jitter, which can shorten any one cycle by that much. What is the setup slack, in ps?", j), edakit.IntStr(s),
		edakit.Sprintf("Setup slack = $T - J + skew - t_{su} - t_{cq} - d_{max}$ = %d - %d + (%d) - %d - %d - %d = %d ps.", p.t, j, skew, p.tsu, p.tcq, p.dmax, s), facts)
}

func uncertainty(ctx generator.GeneratorContext) generator.Problem {
	t := 50 * edakit.Int(10, 24)
	tcq, d, dmin := edakit.Int(30, 80), 10*edakit.Int(25, 90), edakit.Int(10, 100)
	su, th := edakit.Int(20, 60), edakit.Int(10, 50)
	switch ctx.KPOrdinal {
	case 1:
		j, sk, m := edakit.Int(10, 40), 5*edakit.Int(4, 20), 5*edakit.Int(1, 6)
		u := j + sk + m
		return edakit.Numeric(edakit.Sprintf("Before clock tree synthesis, a team budgets setup uncertainty as jitter plus an estimate of the skew the tree will add plus a design margin. The jitter is %d ps, the skew estimate %d ps and the margin %d ps. What value should set_clock_uncertainty -setup carry, in ps?", j, sk, m), edakit.IntStr(u),
			edakit.Sprintf("The terms add: %d + %d + %d = %d ps. After CTS, propagated clocks model the real skew, so the skew term is removed.", j, sk, m, u), map[string]string{"uncertainty": edakit.IntStr(u)})
	case 3:
		uh := edakit.Int(5, 40)
		h := tcq + dmin - (th + uh)
		return edakit.Numeric(edakit.Sprintf("With ideal clocks, a path has clock-to-Q %d ps and shortest logic delay %d ps; the capture flop's hold time is %d ps. The SDC has set_clock_uncertainty -hold %d. What is the hold slack, in ps?", tcq, dmin, th, uh), edakit.IntStr(h),
			edakit.Sprintf("Hold uncertainty raises the hold requirement: (%d + %d) - (%d + %d) = %d ps.", tcq, dmin, th, uh, h), map[string]string{"hold_slack": edakit.IntStr(h)})
	}
	us := edakit.Int(20, 100)
	s := t - us - su - (tcq + d)
	return edakit.Numeric(edakit.Sprintf("With ideal clocks at period %d ps, a path has clock-to-Q %d ps and logic delay %d ps; setup is %d ps. The SDC has set_clock_uncertainty -setup %d. What is the setup slack, in ps?", t, tcq, d, su, us), edakit.IntStr(s),
		edakit.Sprintf("Setup uncertainty pulls the required time earlier: %d - %d - %d - (%d + %d) = %d ps.", t, us, su, tcq, d, s), map[string]string{"setup_slack": edakit.IntStr(s)})
}

// masterEdge returns the time of master edge k (1-based: 1 is the first rise, 2 the first fall).
func masterEdge(k, period, high int) int {
	if k%2 == 1 {
		return (k - 1) / 2 * period
	}
	return (k/2-1)*period + high
}

func generatedClocks(ctx generator.GeneratorContext) generator.Problem {
	p := 2 * edakit.Int(2, 10)
	h := p / 2
	switch ctx.KPOrdinal {
	case 1:
		n := edakit.Pick([]int{2, 3, 4, 8})
		v := n * p
		return edakit.Numeric(edakit.Sprintf("A master clock has period %d ns. A divider output is declared with create_generated_clock -source [get_ports clk] -divide_by %d. What is the generated clock's period, in ns?", p, n), edakit.IntStr(v),
			edakit.Sprintf("Dividing by %d multiplies the period: %d x %d = %d ns.", n, n, p, v), map[string]string{"period": edakit.IntStr(v)})
	case 3:
		src := 10 * edakit.Int(5, 40)
		net1, cq, net2 := 10*edakit.Int(10, 40), edakit.Int(40, 120), 10*edakit.Int(10, 40)
		v := src + net1 + cq + net2
		return edakit.Numeric(edakit.Sprintf("A master clock has source latency %d ps. Its network delivers it to a divider flop's clock pin %d ps later; the divider's clock-to-Q is %d ps; and the generated clock's own network takes %d ps to reach a sink flop. With propagated clocks, when does a master edge at 0 ps produce the generated clock edge at that sink, in ps?", src, net1, cq, net2), edakit.IntStr(v),
			edakit.Sprintf("The generated clock's source latency is the master's latency plus the path through the divider: %d + %d + %d = %d ps. Adding its network gives %d ps.", src, net1, cq, src+net1+cq, v), map[string]string{"arrival": edakit.IntStr(v)})
	}
	e1 := edakit.Int(1, 3)
	e2 := e1 + edakit.Int(1, 3)
	e3 := e2 + edakit.Int(1, 3)
	if (e3-e1)%2 == 1 {
		e3++
	}
	t1, t2, t3 := masterEdge(e1, p, h), masterEdge(e2, p, h), masterEdge(e3, p, h)
	askHigh := edakit.Int(0, 1) == 1
	q := edakit.Sprintf("A master clock has period %d ns and waveform {0 %d}. Its edges are numbered 1 (the rise at 0), 2 (the fall at %d), 3 (the rise at %d) and so on. A generated clock is declared with -edges {%d %d %d}. ", p, h, h, p, e1, e2, e3)
	exp := edakit.Sprintf("Edge %d is at %d ns, edge %d at %d ns and edge %d at %d ns. The generated clock rises at the first, falls at the second and rises again at the third.", e1, t1, e2, t2, e3, t3)
	if askHigh {
		return edakit.Numeric(q+"How long is the generated clock high in each cycle, in ns?", edakit.IntStr(t2-t1), exp+edakit.Sprintf(" High time = %d - %d = %d ns.", t2, t1, t2-t1), map[string]string{"high": edakit.IntStr(t2 - t1)})
	}
	return edakit.Numeric(q+"What is the generated clock's period, in ns?", edakit.IntStr(t3-t1), exp+edakit.Sprintf(" Period = %d - %d = %d ns.", t3, t1, t3-t1), map[string]string{"period": edakit.IntStr(t3 - t1)})
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// relationships returns the setup and hold relationships between two synchronous rising-edge
// clocks. Setup is the smallest gap from a launch edge to the first strictly later capture edge.
// Hold is derived from the most restrictive setup pair (launch t, capture c): the larger of
// (c - Tc) - t and c - (t + Tl). Tied setup pairs take the larger hold.
func relationships(tl, tc int) (int, int) {
	l := tl / gcd(tl, tc) * tc
	setup, hold := 1<<30, -(1 << 30)
	for t := 0; t < l; t += tl {
		c := (t/tc + 1) * tc
		h := (c - tc) - t
		if h2 := c - (t + tl); h2 > h {
			h = h2
		}
		switch {
		case c-t < setup:
			setup, hold = c-t, h
		case c-t == setup && h > hold:
			hold = h
		}
	}
	return setup, hold
}

func launchCapture(ctx generator.GeneratorContext) generator.Problem {
	var tl, tc int
	switch ctx.KPOrdinal {
	case 1:
		base := edakit.Int(2, 6)
		m := edakit.Int(2, 4)
		if edakit.Int(0, 1) == 0 {
			tl, tc = base, base*m
		} else {
			tl, tc = base*m, base
		}
	default:
		for {
			tl, tc = edakit.Int(3, 12), edakit.Int(3, 12)
			if tl != tc && tl%tc != 0 && tc%tl != 0 {
				break
			}
		}
	}
	setup, hold := relationships(tl, tc)
	q := edakit.Sprintf("Data launches from a rising-edge flop on clock CLKA (period %d ns) and is captured by a rising-edge flop on clock CLKB (period %d ns). Both clocks come from one source and rise together at 0 ns. ", tl, tc)
	lcm := tl / gcd(tl, tc) * tc
	facts := map[string]string{"setup": edakit.IntStr(setup), "hold": edakit.IntStr(hold), "lcm": edakit.IntStr(lcm)}
	if ctx.KPOrdinal == 3 {
		return edakit.Numeric(q+"The timer takes the most restrictive setup pair (launch edge, capture edge) and derives two hold checks from it: the capture edge one CLKB period earlier against the same launch edge, and the same capture edge against the next CLKA launch edge. What is the resulting hold relationship (the larger of the two capture-minus-launch values), in ns?", edakit.IntStr(hold),
			edakit.Sprintf("Over the common period %d ns the most restrictive setup gap is %d ns. From that pair, (capture - one CLKB period) - launch and capture - (launch + one CLKA period) are computed, and the larger, %d ns, is the hold relationship.", lcm, setup, hold), facts)
	}
	return edakit.Numeric(q+"Over all launch edges, the setup check uses the nearest later capture edge. What is the most restrictive setup relationship (capture time minus launch time), in ns?", edakit.IntStr(setup),
		edakit.Sprintf("The edge pattern repeats every lcm(%d, %d) = %d ns. For each CLKA edge in that window, find the first CLKB edge strictly later; the smallest gap is %d ns.", tl, tc, lcm, setup), facts)
}

func cppr(ctx generator.GeneratorContext) generator.Problem {
	type seg struct{ e, l int }
	draw := func(n int) []seg {
		var s []seg
		for i := 0; i < n; i++ {
			e := edakit.Int(20, 60)
			s = append(s, seg{e, e + edakit.Int(2, 12)})
		}
		return s
	}
	com := draw(edakit.Int(1, 3))
	la := draw(edakit.Int(1, 2))
	ca := draw(edakit.Int(1, 2))
	list := func(ss []seg) string {
		var p []string
		for _, s := range ss {
			p = append(p, edakit.Sprintf("%d/%d", s.e, s.l))
		}
		return strings.Join(p, ", ")
	}
	sumE, sumL := func(ss []seg) int {
		t := 0
		for _, s := range ss {
			t += s.e
		}
		return t
	}, func(ss []seg) int {
		t := 0
		for _, s := range ss {
			t += s.l
		}
		return t
	}
	q := edakit.Sprintf("A clock tree is analyzed with early/late (min/max) delays per buffer, written early/late in ps. The common segment from the clock root to the branch point has buffers %s. The launch branch adds %s to FF1, and the capture branch adds %s to FF2. ", list(com), list(la), list(ca))
	credit := sumL(com) - sumE(com)
	facts := map[string]string{"credit": edakit.IntStr(credit), "common_late": edakit.IntStr(sumL(com)), "common_early": edakit.IntStr(sumE(com))}
	switch ctx.KPOrdinal {
	case 1:
		late := edakit.Int(0, 1) == 1
		v, w := sumE(com), "early"
		if late {
			v, w = sumL(com), "late"
		}
		return edakit.Numeric(q+edakit.Sprintf("What is the %s arrival of the clock at the branch point (the common point), in ps?", w), edakit.IntStr(v),
			edakit.Sprintf("Sum the %s delays of the common buffers: %d ps.", w, v), facts)
	case 3:
		t := 50 * edakit.Int(10, 20)
		tcq, d, su := edakit.Int(30, 70), 10*edakit.Int(30, 80), edakit.Int(20, 50)
		launch := sumL(com) + sumL(la)
		capt := sumE(com) + sumE(ca)
		before := t + capt - su - (launch + tcq + d)
		after := before + credit
		facts["slack_before"] = edakit.IntStr(before)
		facts["slack_after"] = edakit.IntStr(after)
		return edakit.Numeric(q+edakit.Sprintf("For setup, the launch clock is taken late and the capture clock early. The period is %d ps, clock-to-Q %d ps (already late), data delay %d ps and setup %d ps. What is the setup slack after clock reconvergence pessimism removal, in ps?", t, tcq, d, su), edakit.IntStr(after),
			edakit.Sprintf("Launch clock (late) %d ps, capture clock (early) %d ps, so the raw slack is %d ps. The common segment cannot be late and early at once: add back %d - %d = %d ps, giving %d ps.", launch, capt, before, sumL(com), sumE(com), credit, after), facts)
	}
	return edakit.Numeric(q+"For a setup check from FF1 to FF2 (launch late, capture early), how much pessimism does CPPR remove, in ps?", edakit.IntStr(credit),
		edakit.Sprintf("Only the shared segment is double counted: late minus early over the common buffers, %d - %d = %d ps.", sumL(com), sumE(com), credit), facts)
}
