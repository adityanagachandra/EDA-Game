package sta

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func gateSizing(ctx generator.GeneratorContext) generator.Problem {
	rd := edakit.Pick([]float64{0.5, 1, 1.5, 2})
	rg := edakit.Pick([]float64{1, 1.5, 2, 3, 4})
	d0 := float64(edakit.Int(5, 20))
	cin := float64(edakit.Int(1, 4))
	cl := float64(edakit.Int(8, 40))
	k := float64(edakit.Pick([]int{2, 3, 4}))
	slack := float64(-edakit.Int(5, 60))
	q := edakit.Sprintf("On a failing path, driver D (drive resistance %s kohm) drives only gate G's input (%s fF). G has intrinsic delay %s ps and drive resistance %s kohm and drives a load of %s fF. Cell delay is intrinsic delay plus drive resistance times load (kohm times fF is ps). The ECO replaces G with a size %sx larger, which divides its drive resistance by %s and multiplies its input capacitance by %s; the intrinsic delay is unchanged. ",
		edakit.Num(rd, 1), edakit.Num(cin, 0), edakit.Num(d0, 0), edakit.Num(rg, 1), edakit.Num(cl, 0), edakit.Num(k, 0), edakit.Num(k, 0), edakit.Num(k, 0))
	newG := d0 + rg/k*cl
	dG := newG - (d0 + rg*cl)
	dD := rd * (k - 1) * cin
	s2 := slack - (dG + dD)
	facts := map[string]string{"new_g": edakit.Num(newG, 2), "delta_g": edakit.Num(dG, 2), "delta_driver": edakit.Num(dD, 2), "slack": edakit.Num(s2, 2)}
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(newG, 2)
		return edakit.Numeric(q+"What is G's new delay, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("%s + %s / %s x %s = %s ps, down from %s ps.", edakit.Num(d0, 0), edakit.Num(rg, 1), edakit.Num(k, 0), edakit.Num(cl, 0), ans, edakit.Num(d0+rg*cl, 2)), facts)
	case 3:
		ans := edakit.Num(s2, 2)
		return edakit.Numeric(q+edakit.Sprintf("The path's setup slack was %s ps. What is it after the ECO, in ps? Round to 2 decimal places.", edakit.Num(slack, 0)), ans,
			edakit.Sprintf("G gets faster by %s ps, and D gets slower by %s ps because its load grew. Slack = %s - (%s + %s) = %s ps.", edakit.Num(-dG, 2), edakit.Num(dD, 2), edakit.Num(slack, 0), edakit.Num(dG, 2), edakit.Num(dD, 2), ans), facts)
	}
	ans := edakit.Num(dD, 2)
	return edakit.Numeric(q+"By how many ps does driver D slow down? Round to 2 decimal places.", ans,
		edakit.Sprintf("D's load grows by (%s - 1) x %s fF; times %s kohm that is %s ps. Sizing up moves delay upstream, which is why it is not free.", edakit.Num(k, 0), edakit.Num(cin, 0), edakit.Num(rd, 1), ans), facts)
}

func holdFixing(ctx generator.GeneratorContext) generator.Problem {
	h := edakit.Int(5, 80)
	bmin := edakit.Int(8, 20)
	bmax := bmin + edakit.Int(5, 20)
	n := int(math.Ceil(float64(h) / float64(bmin)))
	s := n*bmax + edakit.Int(-40, 150)
	if s < bmax {
		s = bmax + edakit.Int(0, 30)
	}
	q := edakit.Sprintf("A path fails hold by %d ps at the fast corner (hold slack -%d ps) and has %d ps of setup slack at the slow corner. A delay buffer adds %d ps at the fast corner and %d ps at the slow corner. ", h, h, s, bmin, bmax)
	facts := map[string]string{"buffers": edakit.IntStr(n), "setup_after": edakit.IntStr(s - n*bmax), "max_buffers": edakit.IntStr(s / bmax)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+"What is the fewest buffers that fixes the hold violation?", edakit.IntStr(n),
			edakit.Sprintf("Each buffer adds %d ps at the fast corner: ceil(%d / %d) = %d.", bmin, h, bmin, n), facts)
	case 3:
		m := s / bmax
		return edakit.Numeric(q+"What is the largest number of buffers the setup slack allows?", edakit.IntStr(m),
			edakit.Sprintf("Each buffer costs %d ps of setup slack at the slow corner: floor(%d / %d) = %d. The fix is feasible only if this is at least %d.", bmax, s, bmax, m, n), facts)
	}
	v := s - n*bmax
	return edakit.Numeric(q+"After inserting the fewest buffers that fix hold, what is the setup slack, in ps?", edakit.IntStr(v),
		edakit.Sprintf("%d buffers are needed (ceil(%d / %d)). They cost %d x %d = %d ps at the slow corner: %d - %d = %d ps.", n, h, bmin, n, bmax, n*bmax, s, n*bmax, v), facts)
}

func usefulSkew(ctx generator.GeneratorContext) generator.Problem {
	sa := float64(-edakit.Int(10, 80))
	sb := float64(edakit.Int(10, 120))
	q := edakit.Sprintf("In a pipeline FF1 → FF2 → FF3, the stage into FF2 has setup slack %s ps and the stage out of FF2 has setup slack %s ps. Delaying FF2's clock by $\\delta$ gives the first stage $\\delta$ more and the second stage $\\delta$ less. ", edakit.Num(sa, 0), edakit.Num(sb, 0))
	bal := (sb - sa) / 2
	switch ctx.KPOrdinal {
	case 1:
		d := float64(5 * edakit.Int(1, 20))
		w := math.Min(sa+d, sb-d)
		return edakit.Numeric(q+edakit.Sprintf("With $\\delta$ = %s ps, what is the worse of the two setup slacks, in ps?", edakit.Num(d, 0)), edakit.Num(w, 0),
			edakit.Sprintf("Stage 1: %s + %s = %s ps. Stage 2: %s - %s = %s ps. The worse is %s ps.", edakit.Num(sa, 0), edakit.Num(d, 0), edakit.Num(sa+d, 0), edakit.Num(sb, 0), edakit.Num(d, 0), edakit.Num(sb-d, 0), edakit.Num(w, 0)), map[string]string{"worst": edakit.Num(w, 0)})
	case 3:
		ha := float64(edakit.Int(5, 80))
		d := math.Min(bal, ha)
		w := math.Min(sa+d, sb-d)
		ans := edakit.Num(w, 1)
		return edakit.Numeric(q+edakit.Sprintf("Delaying FF2's clock also eats the hold slack of the stage into FF2, which is %s ps. Choosing $\\delta$ to maximize the worse setup slack without breaking hold, what is that worse setup slack, in ps? Round to 1 decimal place.", edakit.Num(ha, 0)), ans,
			edakit.Sprintf("Balancing needs $\\delta$ = %s ps, but hold allows at most %s ps, so $\\delta$ = %s ps and the worse setup slack is %s ps.", edakit.Num(bal, 1), edakit.Num(ha, 0), edakit.Num(d, 1), ans), map[string]string{"delta": edakit.Num(d, 1), "worst": ans})
	}
	ans := edakit.Num(bal, 1)
	return edakit.Numeric(q+"What $\\delta$ makes the two setup slacks equal, in ps? Round to 1 decimal place.", ans,
		edakit.Sprintf("Solve %s + $\\delta$ = %s - $\\delta$: $\\delta$ = (%s - (%s)) / 2 = %s ps, leaving both at %s ps.", edakit.Num(sa, 0), edakit.Num(sb, 0), edakit.Num(sb, 0), edakit.Num(sa, 0), ans, edakit.Num((sa+sb)/2, 1)), map[string]string{"delta": ans})
}

var closureItems = []edakit.Item{
	{KP: 1, Q: "After placement a block has WNS -40 ps and TNS -9000 ps over 2000 endpoints. Which response targets the problem?",
		A: "Look for a systematic cause such as a bad floorplan, constraint or congestion, since the failures are widespread", D: []string{"Upsize the cells on the single worst path", "Add hold buffers everywhere", "Raise the clock uncertainty"},
		Why: "A large TNS with many failing endpoints signals a global cause. Fixing the worst path alone barely moves TNS."},
	{KP: 1, Q: "Why do closure flows fix setup before hold?",
		A: "Setup fixes change delays on many paths, and hold fixes add delay that is cheap to insert at the end", D: []string{"Hold violations do not matter until tapeout", "Hold fixing reduces setup slack to zero", "Setup and hold cannot be fixed in the same tool"},
		Why: "Sizing and buffering for setup reshape fast paths too, so hold is fixed afterwards with targeted delay insertion that consumes little setup slack."},
	{KP: 1, Q: "A design meets timing in the placer's estimate but fails after routing. What is the most likely reason?",
		A: "The placer's wire model was optimistic compared with routed parasitics and detours", D: []string{"Routing always speeds up a design", "The Liberty library changed", "The clock period was halved by the router"},
		Why: "Early stages estimate wires with Steiner lengths or HPWL. Detours, layer assignment and coupling only appear after routing, so margins carry forward."},
	{KP: 1, Q: "Why do teams carry extra margin, such as a tighter clock, in early implementation stages?",
		A: "Later stages add delay the early models cannot see, so early slack must cover it", D: []string{"To make the tools run faster", "Because Liberty data is less accurate at signoff", "To create hold violations deliberately"},
		Why: "Each stage's estimate is optimistic relative to the next. Over-constraining early leaves room for routing, crosstalk and variation effects."},
	{KP: 2, Q: "A path fails setup by 8 ps. Its gates are already large, and its wires are long. Which fix is most promising?",
		A: "Buffer or reroute the long wires, or move the cells closer", D: []string{"Upsize the gates further", "Add hold buffers", "Relax the hold constraint"},
		Why: "When wire delay dominates, gate sizing gives diminishing returns. Repeaters, layer promotion or placement changes attack the wire."},
	{KP: 2, Q: "Swapping a cell from high-Vt to low-Vt fixes a setup violation. What does it cost?",
		A: "Higher leakage power", D: []string{"Larger area only", "Lower dynamic power", "Worse hold on unrelated paths only"},
		Why: "Low-Vt cells switch faster but leak far more, so Vt swapping is applied only where slack demands it."},
	{KP: 2, Q: "Fixing the worst path improves WNS by 5 ps, but a different endpoint becomes the new worst at almost the same slack. What does this suggest?",
		A: "Many paths are near-critical, a slack wall, so broad optimization is needed", D: []string{"The timer is unstable", "The fix was applied to the wrong corner", "WNS cannot be improved further"},
		Why: "A slack wall means many endpoints sit near the worst value. Progress then shows up in TNS and the count of failing endpoints."},
	{KP: 2, Q: "An optimizer improves the slow-corner setup but breaks fast-corner hold. Why must closure be run across scenarios at once?",
		A: "A fix that helps one scenario can hurt another, so each change must be checked in every active scenario", D: []string{"Scenarios share one netlist only during signoff", "Hold is not checked at the fast corner", "MCMM runs only after closure"},
		Why: "MCMM optimization evaluates each change in all dominant scenarios, otherwise fixes in one scenario create violations in another."},
	{KP: 3, Q: "Clock tree synthesis adds 50 ps of skew to a block that met setup with ideal clocks. Which constraint change reflects the transition correctly?",
		A: "Switch to propagated clocks and drop the skew estimate from clock uncertainty", D: []string{"Keep ideal clocks and add 50 ps more uncertainty", "Remove all clock uncertainty", "Add multicycle paths on skewed endpoints"},
		Why: "After CTS the timer computes real latencies and skew. Leaving the pre-CTS skew budget in uncertainty would count skew twice."},
	{KP: 3, Q: "A late ECO must fix three violations without changing the placement of other cells. What kind of ECO is this?",
		A: "A small, local netlist change applied incrementally, with the rest of the layout frozen", D: []string{"A full re-synthesis of the block", "A new floorplan", "A change to the Liberty library"},
		Why: "Timing ECOs late in the flow resize, swap or insert a few cells and re-time incrementally so that verified work elsewhere stays intact."},
	{KP: 3, Q: "Signoff STA fails by 3 ps on a path that the implementation tool reported as passing. What is the usual cause?",
		A: "The implementation tool's timer is a faster approximation that differs from the signoff timer's delay calculation and parasitics", D: []string{"Signoff timers ignore wire delay", "The SDC was ignored at signoff", "Signoff runs only the typical corner"},
		Why: "Implementation tools use fast built-in timers. Correlation to signoff is close but not exact, which is why flows keep a little margin."},
	{KP: 3, Q: "After closure, a metal-only ECO is required. Which fixes remain available?",
		A: "Rewiring and using spare cells already placed in the layout", D: []string{"Inserting new cells anywhere", "Changing transistor sizes inside standard cells", "Moving the floorplan blocks"},
		Why: "A metal-only ECO keeps the base layers fixed, so only pre-placed spare cells and metal changes can be used."},
}

var openstaItems = []edakit.Item{
	{KP: 1, Q: "Which OpenSTA command loads the cell timing models that give arcs their delays?",
		A: "read_liberty", D: []string{"read_verilog", "read_sdc", "read_spef"},
		Why: "Liberty files hold the cells' timing arcs, delay and slew tables, pin capacitances and checks."},
	{KP: 1, Q: "What does link_design do after read_verilog?",
		A: "Binds the netlist's instances to library cells and builds the design the timer will analyze", D: []string{"Runs placement", "Reads the clock constraints", "Writes a timing report"},
		Why: "The Verilog netlist names cells; linking resolves them against the loaded Liberty libraries and elaborates the top module."},
	{KP: 1, Q: "Which file tells OpenSTA the clock periods and input and output delays?",
		A: "The SDC file read with read_sdc", D: []string{"The Liberty file", "The SPEF file", "The Verilog netlist"},
		Why: "SDC carries the constraints: clocks, I/O delays, exceptions and uncertainty."},
	{KP: 1, Q: "What does read_spef add to the analysis?",
		A: "Extracted parasitic resistance and capacitance of the routed nets", D: []string{"The cell delay tables", "The clock definitions", "The placement coordinates"},
		Why: "SPEF holds RC networks per net. Without it the timer uses wire-load or zero-wire estimates."},
	{KP: 2, Q: "A routed design is timed in OpenSTA without set_propagated_clock. What happens to clock latency?",
		A: "Clocks stay ideal, so the real clock-tree delays and skew are ignored", D: []string{"The timer propagates through the tree anyway", "Every flop gets a random latency", "Setup checks are disabled"},
		Why: "Clocks are ideal by default. After CTS the clocks must be marked propagated to use the buffer and wire delays of the tree."},
	{KP: 2, Q: "Which command prints the worst setup paths with their arrival, required time and slack?",
		A: "report_checks -path_delay max", D: []string{"report_checks -path_delay min", "report_power", "report_clock_properties"},
		Why: "report_checks reports timing paths. -path_delay max selects setup (late) checks, and min selects hold."},
	{KP: 2, Q: "A script needs only the design's total negative slack as a number. Which command gives it?",
		A: "report_tns", D: []string{"report_wns", "report_checks -group_count 1", "report_worst_slack -min"},
		Why: "report_tns prints total negative slack, report_wns worst negative slack, and report_worst_slack the worst slack itself."},
	{KP: 2, Q: "How does OpenSTA compute the delay of a routed net once SPEF is loaded?",
		A: "It reduces the RC network and runs a delay calculator on the driver and load model", D: []string{"It reads the net delay directly from the SDC", "It uses the wire length from the Verilog", "It sets every net delay to zero"},
		Why: "The delay calculator reduces parasitics, models the driver from Liberty and computes stage delay and slew. Simple Elmore-style or more accurate calculators are available."},
	{KP: 3, Q: "What does Liberty abstract away that a SPICE simulation of the same cell would keep?",
		A: "The transistor-level waveform, replaced by tables of delay and slew against input slew and load", D: []string{"The cell's pin names", "The function of the cell", "The cell's area"},
		Why: "NLDM tables reduce each arc to two numbers per slew and load point. That is what makes timing millions of cells fast, and what effective-capacitance models compensate for."},
	{KP: 3, Q: "A team wants their in-house placer to call the timer after every move. Which property of OpenSTA makes that practical?",
		A: "Incremental timing: after a netlist or parasitic change, only affected arrivals and required times are recomputed", D: []string{"It rereads all files on every query", "It runs SPICE for each move", "It caches only the final report"},
		Why: "OpenSTA invalidates the changed portion of the graph and recomputes on demand, so queries after small edits are cheap."},
	{KP: 3, Q: "Why does OpenROAD embed OpenSTA instead of writing timing files for a separate run?",
		A: "Optimizers such as resizing and buffering query slack in memory after each change", D: []string{"Separate runs would be more accurate", "OpenSTA cannot read files", "Embedding removes the need for SDC"},
		Why: "Timing-driven optimization needs fast in-memory queries and incremental updates. File round-trips per change would be far too slow."},
	{KP: 3, Q: "Two runs of the same design give different slacks. One loaded SPEF; the other did not. What does that show?",
		A: "Wire parasitics are a first-order part of delay, so timing without them is only an estimate", D: []string{"The timer is nondeterministic", "SPEF only affects hold", "Liberty overrides SPEF"},
		Why: "Without parasitics, net delays and pin loads are estimates. Signoff timing always uses extracted parasitics."},
}
