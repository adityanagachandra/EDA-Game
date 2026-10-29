package sta

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// regPath is a flop-to-flop path with its clock latencies.
type regPath struct {
	t, lLaunch, lCap, tcq, dmax, dmin, tsu, th int
}

func drawRegPath() regPath {
	return regPath{
		t:       50 * edakit.Int(10, 24),
		lLaunch: 5 * edakit.Int(10, 40),
		lCap:    5 * edakit.Int(10, 40),
		tcq:     edakit.Int(30, 80),
		dmax:    10 * edakit.Int(25, 90),
		dmin:    edakit.Int(10, 120),
		tsu:     edakit.Int(20, 60),
		th:      edakit.Int(10, 50),
	}
}

func (p regPath) describe() string {
	return edakit.Sprintf("A flop-to-flop path is clocked with period %d ps. The launch flop's clock arrives at %d ps and the capture flop's clock at %d ps after the source edge. Clock-to-Q is %d ps, the longest combinational delay is %d ps and the shortest is %d ps. The capture flop's setup time is %d ps and its hold time is %d ps.",
		p.t, p.lLaunch, p.lCap, p.tcq, p.dmax, p.dmin, p.tsu, p.th)
}

func setupHold(ctx generator.GeneratorContext) generator.Problem {
	p := drawRegPath()
	req := p.t + p.lCap - p.tsu
	arr := p.lLaunch + p.tcq + p.dmax
	setup := req - arr
	hold := (p.lLaunch + p.tcq + p.dmin) - (p.lCap + p.th)
	facts := map[string]string{"required": edakit.IntStr(req), "arrival": edakit.IntStr(arr), "setup_slack": edakit.IntStr(setup), "hold_slack": edakit.IntStr(hold)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(p.describe()+" What is the setup required time at the capture flop's D pin, in ps?", edakit.IntStr(req),
			edakit.Sprintf("The next capturing edge reaches the capture flop at %d + %d ps; setup must be met %d ps before it: %d ps.", p.t, p.lCap, p.tsu, req), facts)
	case 3:
		return edakit.Numeric(p.describe()+" What is the hold slack, in ps?", edakit.IntStr(hold),
			edakit.Sprintf("Hold compares the earliest arrival, %d + %d + %d = %d ps, with the same capturing edge plus hold, %d + %d = %d ps. Slack is %d ps. The period does not appear.",
				p.lLaunch, p.tcq, p.dmin, p.lLaunch+p.tcq+p.dmin, p.lCap, p.th, p.lCap+p.th, hold), facts)
	}
	return edakit.Numeric(p.describe()+" What is the setup slack, in ps?", edakit.IntStr(setup),
		edakit.Sprintf("Required %d + %d - %d = %d ps; arrival %d + %d + %d = %d ps; slack %d ps.", p.t, p.lCap, p.tsu, req, p.lLaunch, p.tcq, p.dmax, arr, setup), facts)
}

func maxFrequency(ctx generator.GeneratorContext) generator.Problem {
	type pth struct{ tcq, d, tsu, skew int }
	draw := func() pth {
		return pth{edakit.Int(30, 80), 10 * edakit.Int(20, 90), edakit.Int(20, 60), 5 * edakit.Int(-8, 8)}
	}
	tmin := func(p pth) int { return p.tcq + p.d + p.tsu - p.skew }
	line := func(i int, p pth) string {
		return edakit.Sprintf("path %d has clock-to-Q %d ps, logic delay %d ps, setup %d ps and skew %s ps", i, p.tcq, p.d, p.tsu, signed(p.skew))
	}
	note := " Skew is the capture clock arrival minus the launch clock arrival."
	if ctx.KPOrdinal == 1 {
		p := draw()
		v := tmin(p)
		return edakit.Numeric("In a register-to-register path, "+line(1, p)+"."+note+" What is the smallest clock period this path meets, in ps?", edakit.IntStr(v),
			edakit.Sprintf("Setup holds when $T \\ge t_{cq} + d + t_{su} - skew$ = %d + %d + %d - (%s) = %d ps.", p.tcq, p.d, p.tsu, signed(p.skew), v), map[string]string{"tmin": edakit.IntStr(v)})
	}
	n := edakit.Scaled(ctx, 2, 4)
	var parts []string
	worst := 0
	for i := 0; i < n; i++ {
		p := draw()
		parts = append(parts, line(i+1, p))
		if v := tmin(p); v > worst {
			worst = v
		}
	}
	intro := "A design has these register-to-register paths: " + strings.Join(parts, "; ") + "." + note
	if ctx.KPOrdinal == 3 {
		f := 1e6 / float64(worst)
		ans := edakit.Num(f, 1)
		return edakit.Numeric(intro+" What is the highest clock frequency the design meets, in MHz? Round to 1 decimal place.", ans,
			edakit.Sprintf("The design period is the largest per-path minimum, %d ps. Frequency is $10^6 / %d$ = %s MHz.", worst, worst, ans), map[string]string{"tmin": edakit.IntStr(worst), "fmax": ans})
	}
	return edakit.Numeric(intro+" What is the smallest clock period the whole design meets, in ps?", edakit.IntStr(worst),
		edakit.Sprintf("Each path needs $t_{cq} + d + t_{su} - skew$; the design needs the largest of these, %d ps.", worst), map[string]string{"tmin": edakit.IntStr(worst)})
}

func halfCycle(ctx generator.GeneratorContext) generator.Problem {
	t := 100 * edakit.Int(6, 16)
	h := t * edakit.Int(40, 60) / 100
	tcq := edakit.Int(30, 70)
	d := 10 * edakit.Int(5, h/10-5)
	dmin := edakit.Int(10, 80)
	tsu, th := edakit.Int(20, 50), edakit.Int(15, 50)
	riseToFall := edakit.Pick([]bool{true, false})
	var avail, holdEdge int
	var desc string
	if riseToFall {
		avail = h
		holdEdge = h - t
		desc = "launches on the rising edge at 0 ps and is captured by a falling-edge flop"
	} else {
		avail = t - h
		holdEdge = -h
		desc = edakit.Sprintf("launches on the falling edge at %d ps and is captured by a rising-edge flop", h)
	}
	setup := avail - tsu - (tcq + d)
	hold := (tcq + dmin) - (holdEdge + th)
	q := edakit.Sprintf("A clock has period %d ps and is high for the first %d ps of each cycle (rising at 0). A path %s. Clock-to-Q is %d ps, the longest logic delay %d ps and the shortest %d ps; setup is %d ps and hold %d ps. Clocks are ideal. ",
		t, h, desc, tcq, d, dmin, tsu, th)
	facts := map[string]string{"available": edakit.IntStr(avail), "setup_slack": edakit.IntStr(setup), "hold_slack": edakit.IntStr(hold)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+"How much time separates the launch edge from the setup capture edge, in ps?", edakit.IntStr(avail),
			edakit.Sprintf("The setup capture edge is the first opposite edge after launch, %d ps later.", avail), facts)
	case 3:
		return edakit.Numeric(q+"What is the hold slack, in ps?", edakit.IntStr(hold),
			edakit.Sprintf("The hold check uses the capturing edge just before the setup capture edge, %d ps relative to launch. Slack = (%d + %d) - (%d + %d) = %d ps, so half-cycle paths rarely fail hold.", holdEdge, tcq, dmin, holdEdge, th, hold), facts)
	}
	return edakit.Numeric(q+"What is the setup slack, in ps?", edakit.IntStr(setup),
		edakit.Sprintf("Only %d ps separate launch from capture. Slack = %d - %d - (%d + %d) = %d ps.", avail, avail, tsu, tcq, d, setup), facts)
}

func recoveryRemoval(ctx generator.GeneratorContext) generator.Problem {
	t := 50 * edakit.Int(10, 24)
	ll, lc := 5*edakit.Int(10, 40), 5*edakit.Int(10, 40)
	tcq := edakit.Int(30, 80)
	dmax := 10 * edakit.Int(20, 80)
	dmin := edakit.Int(20, 150)
	trec, trem := edakit.Int(30, 90), edakit.Int(20, 70)
	q := edakit.Sprintf("A reset synchronizer flop, clocked at period %d ps with clock arrival %d ps, drives the asynchronous reset pin of a capture flop whose clock arrives at %d ps. The synchronizer's clock-to-Q is %d ps, and the reset net delay is at most %d ps and at least %d ps. The capture flop's recovery time is %d ps and its removal time is %d ps. ",
		t, ll, lc, tcq, dmax, dmin, trec, trem)
	req := t + lc - trec
	rec := req - (ll + tcq + dmax)
	rem := (ll + tcq + dmin) - (lc + trem)
	facts := map[string]string{"recovery_required": edakit.IntStr(req), "recovery_slack": edakit.IntStr(rec), "removal_slack": edakit.IntStr(rem)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+"What is the required time for the reset deassertion at the capture flop (the recovery check), in ps?", edakit.IntStr(req),
			edakit.Sprintf("Recovery is setup for the deasserting edge: %d + %d - %d = %d ps.", t, lc, trec, req), facts)
	case 3:
		return edakit.Numeric(q+"What is the removal slack, in ps?", edakit.IntStr(rem),
			edakit.Sprintf("Removal is hold for the deasserting edge: (%d + %d + %d) - (%d + %d) = %d ps.", ll, tcq, dmin, lc, trem, rem), facts)
	}
	return edakit.Numeric(q+"What is the recovery slack, in ps?", edakit.IntStr(rec),
		edakit.Sprintf("Required %d ps minus the latest deassertion arrival %d + %d + %d = %d ps gives %d ps.", req, ll, tcq, dmax, ll+tcq+dmax, rec), facts)
}

func latchBorrow(ctx generator.GeneratorContext) generator.Problem {
	t := 100 * edakit.Int(6, 14)
	h := t / 2
	tsu := edakit.Int(20, 50)
	dq := edakit.Int(20, 60)
	tsu2 := edakit.Int(20, 50)
	var a int
	if ctx.KPOrdinal == 1 && edakit.Int(0, 3) == 0 {
		a = t - 10*edakit.Int(1, 10)
	} else {
		a = t + 10*edakit.Int(1, (h-tsu)/10-1)
	}
	d2 := 10 * edakit.Int(10, (t-60)/10)
	borrow := int(math.Max(0, float64(a-t)))
	maxB := h - tsu
	depart := a
	if a < t {
		depart = t
	}
	slack2 := 2*t - tsu2 - (depart + dq + d2)
	q := edakit.Sprintf("A rising-edge flop launches data at 0 ps into a positive level-sensitive latch. The clock has period %d ps with a 50%% duty cycle, so the latch opens at %d ps and closes at %d ps; the latch setup time before closing is %d ps. Data arrives at the latch D pin at %d ps. ",
		t, t, t+h, tsu, a)
	facts := map[string]string{"borrow": edakit.IntStr(borrow), "max_borrow": edakit.IntStr(maxB), "next_slack": edakit.IntStr(slack2)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+"How much time does the path borrow from the next stage, in ps?", edakit.IntStr(borrow),
			edakit.Sprintf("Borrowing is the time past the opening edge: max(0, %d - %d) = %d ps.", a, t, borrow), facts)
	case 2:
		return edakit.Numeric(q+"What is the maximum time this latch can lend, in ps?", edakit.IntStr(maxB),
			edakit.Sprintf("Data may arrive as late as the closing edge minus setup, so the most it can borrow is the pulse width minus setup: %d - %d = %d ps.", h, tsu, maxB), facts)
	}
	return edakit.Numeric(q+edakit.Sprintf("The latch's D-to-Q delay is %d ps, and from Q the logic takes %d ps to a rising-edge flop with setup %d ps that captures at %d ps. What is the setup slack of that next stage, in ps?", dq, d2, tsu2, 2*t), edakit.IntStr(slack2),
		edakit.Sprintf("Data leaves the latch at %d + %d ps because the latch is already open. Slack = %d - %d - (%d + %d + %d) = %d ps. The borrowed time is charged to this stage.", depart, dq, 2*t, tsu2, depart, dq, d2, slack2), facts)
}

var checkTypeItems = []edakit.Item{
	{KP: 1, Q: "A data path from one rising-edge flop to another must settle before the next capturing edge. Which check guards this?",
		A: "Setup", D: []string{"Hold", "Removal", "Minimum pulse width"},
		Why: "Setup checks the latest arrival against the next capturing edge minus the setup time."},
	{KP: 1, Q: "New data from a fast path could overwrite the value a flop is capturing on the same edge. Which check guards this?",
		A: "Hold", D: []string{"Setup", "Recovery", "Clock-gating setup"},
		Why: "Hold checks the earliest arrival against the same capturing edge plus the hold time."},
	{KP: 1, Q: "A reset signal deasserts close to a clock edge at a flop's asynchronous clear pin. Which check makes sure the flop leaves reset cleanly on that edge?",
		A: "Recovery, with removal as its hold-like partner", D: []string{"Setup on the D pin", "Minimum pulse width on the clear pin only", "No check: asynchronous pins are untimed"},
		Why: "Deassertion of an asynchronous control is timed against the clock: recovery is the setup-like check and removal the hold-like one."},
	{KP: 1, Q: "Why does a timer check only the deasserting edge of an asynchronous reset and not the asserting edge?",
		A: "Assertion takes effect immediately regardless of the clock, while deassertion races the clock edge", D: []string{"Assertion is always a false path in SDC", "Liberty files lack data for assertion", "Assertion is checked as hold"},
		Why: "Asserting reset forces the output without the clock. Releasing it near a clock edge can leave the flop metastable, so that edge is checked."},
	{KP: 2, Q: "An AND gate gates a clock with an enable signal. What does the timer check on the enable pin?",
		A: "Clock-gating setup and hold, so the enable changes only while the clock is low", D: []string{"Recovery and removal", "Minimum pulse width on the enable", "Nothing, because the enable is a data pin"},
		Why: "If the enable changed while the clock was high, the gated clock would glitch. Timers infer clock-gating checks on such gates."},
	{KP: 2, Q: "A clock divider produces a narrow high pulse after a long clock tree. Which check can fail even if every setup and hold check passes?",
		A: "Minimum pulse width at the sink's clock pin", D: []string{"Recovery", "Clock-gating hold", "Data-to-data"},
		Why: "Rise and fall delays differ along a tree, so pulses shrink. Liberty gives a minimum pulse width that the clock pin needs."},
	{KP: 2, Q: "Two data signals entering a custom block must change within a fixed window of each other, and neither is a clock. Which check expresses this?",
		A: "A data-to-data check", D: []string{"A setup check on a clock pin", "A recovery check", "A clock-gating check"},
		Why: "A data-to-data check (set_data_check in SDC) times one data pin against another, without a clock edge."},
	{KP: 2, Q: "A latch is the capturing element. Against which clock edge is its setup checked?",
		A: "The closing edge, with borrowing allowed after the opening edge", D: []string{"The opening edge, with no borrowing", "The edge one full period later", "Latches have no setup check"},
		Why: "A transparent latch accepts data until it closes. Arrival after the opening edge borrows time; the hard limit is closing minus setup."},
	{KP: 3, Q: "A path runs from a rising-edge flop to a falling-edge flop on the same clock. Which statement is right?",
		A: "Setup has only part of a cycle, and hold gets an easy check", D: []string{"Setup gets a full cycle and hold is unchanged", "The path is never checked", "Hold gets only half a cycle and setup a full cycle"},
		Why: "The setup capture edge is the next falling edge, part of a cycle away, and the hold edge is the falling edge before launch, which is far in the past."},
	{KP: 3, Q: "A flop's hold check fails by 5 ps at the fast corner, while its setup passes by 300 ps at the slow corner. What is the usual fix?",
		A: "Add delay on the short path, sized at the fast corner and checked at the slow corner", D: []string{"Lengthen the clock period", "Upsize the gates on the short path", "Mark the path as a multicycle path"},
		Why: "Hold is independent of the period, so delay must be added to the short path. The slow-corner setup margin bounds how much."},
	{KP: 3, Q: "A synchronizer's first flop captures an asynchronous input. What does the timer do with the path into its D pin?",
		A: "It is excluded with an exception, since no fixed phase relationship exists", D: []string{"It runs normal setup and hold against the faster clock", "It applies recovery and removal", "It treats it as a half-cycle path"},
		Why: "A signal from an unrelated clock has no defined launch-capture relationship; the path is declared false or the clocks asynchronous, and metastability is handled by the synchronizer design."},
	{KP: 3, Q: "A setup check and a hold check are both reported on the same flop. What distinguishes their capture edges?",
		A: "Setup uses the next capturing edge; hold uses the same edge the data was launched relative to", D: []string{"Both use the next capturing edge", "Hold uses the next edge and setup the same edge", "Both use the launching edge"},
		Why: "For a single-cycle path, setup is checked one period after launch and hold zero periods after launch."},
}
