// Package sta generates timing-analysis problems: the timing graph, delay models, timing checks,
// clocks, constraints, path analysis, variation-aware timing, incremental timing, closure and the
// agent-first frontier. Files are split by subdomain; this one holds Register and the two
// reference concepts.
package sta

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the timing-analysis generators to reg.
func Register(reg *generator.Registry) {
	reg.Register("sta.graph.timing_graph", &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: timingGraphItems}, edakit.Func(timingGraphCount)}})
	reg.Register("sta.graph.unateness", edakit.Func(unateness))
	reg.Register("sta.graph.levelization", edakit.Func(levelization))
	reg.Register("sta.graph.arrival_time", edakit.Func(arrivalTime))
	reg.Register("sta.graph.arrival_required", edakit.Func(arrivalRequired))
	reg.Register("sta.graph.slack_meaning", &edakit.Bank{Items: slackMeaningItems})
	reg.Register("sta.graph.node_slack", edakit.Func(nodeSlack))
	reg.Register("sta.graph.wns_tns", edakit.Func(wnsTNS))

	reg.Register("sta.delay.elmore", edakit.Func(elmore))
	reg.Register("sta.delay.wire_rc", edakit.Func(wireRC))
	reg.Register("sta.delay.nldm_interpolation", edakit.Func(nldm))
	reg.Register("sta.delay.stage_delay", edakit.Func(stageDelay))
	reg.Register("sta.delay.slew_propagation", edakit.Func(slewProp))
	reg.Register("sta.delay.effective_capacitance", edakit.Func(ceff))

	reg.Register("sta.check.setup_hold", edakit.Func(setupHold))
	reg.Register("sta.check.max_frequency", edakit.Func(maxFrequency))
	reg.Register("sta.check.half_cycle", edakit.Func(halfCycle))
	reg.Register("sta.check.recovery_removal", edakit.Func(recoveryRemoval))
	reg.Register("sta.check.latch_borrowing", edakit.Func(latchBorrow))
	reg.Register("sta.check.check_types", &edakit.Bank{Items: checkTypeItems})

	reg.Register("sta.clock.waveform", edakit.Func(waveform))
	reg.Register("sta.clock.latency", edakit.Func(latency))
	reg.Register("sta.clock.skew_jitter", edakit.Func(skewJitter))
	reg.Register("sta.clock.uncertainty", edakit.Func(uncertainty))
	reg.Register("sta.clock.generated_clocks", edakit.Func(generatedClocks))
	reg.Register("sta.clock.launch_capture_edges", edakit.Func(launchCapture))
	reg.Register("sta.clock.cppr", edakit.Func(cppr))

	reg.Register("sta.constraint.sdc", &edakit.Bank{Items: sdcItems})
	reg.Register("sta.constraint.io_delay", edakit.Func(ioDelay))
	reg.Register("sta.constraint.exceptions", &edakit.Bank{Items: exceptionItems})
	reg.Register("sta.constraint.multicycle_edges", edakit.Func(multicycle))

	reg.Register("sta.path.path_count", edakit.Func(pathCount))
	reg.Register("sta.path.k_worst", edakit.Func(kWorst))
	reg.Register("sta.path.gba_vs_pba", &edakit.Bank{Items: gbaItems})
	reg.Register("sta.path.path_report", edakit.Func(pathReport))

	reg.Register("sta.var.corners", edakit.Func(corners))
	reg.Register("sta.var.mcmm", edakit.Func(mcmm))
	reg.Register("sta.var.ocv_derate", edakit.Func(ocvDerate))
	reg.Register("sta.var.aocv", edakit.Func(aocv))
	reg.Register("sta.var.pocv", edakit.Func(pocv))
	reg.Register("sta.var.ssta", edakit.Func(ssta))
	reg.Register("sta.var.clark_max", edakit.Func(clarkMax))
	reg.Register("sta.var.timing_yield", edakit.Func(timingYield))

	reg.Register("sta.incremental.update", edakit.Func(incrementalUpdate))
	reg.Register("sta.incremental.early_termination", edakit.Func(earlyTermination))

	reg.Register("sta.closure.timing_closure", &edakit.Bank{Items: closureItems})
	reg.Register("sta.closure.gate_sizing_eco", edakit.Func(gateSizing))
	reg.Register("sta.closure.hold_fixing", edakit.Func(holdFixing))
	reg.Register("sta.closure.useful_skew", edakit.Func(usefulSkew))
	reg.Register("sta.closure.opensta_flow", &edakit.Bank{Items: openstaItems})

	reg.Register("sta.frontier.derate_provider", &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: derateProviderItems}, edakit.Func(derateProviderCalc)}})
	reg.Register("sta.frontier.change_tiers", &edakit.Bank{Items: changeTierItems})
	reg.Register("sta.frontier.differential_timing", edakit.Func(differentialTiming))
	reg.Register("sta.frontier.surrogate_timer", edakit.Func(surrogateTimer))
	reg.Register("sta.frontier.agentic_eco_loop", &edakit.Bank{Items: ecoLoopItems})
	reg.Register("sta.frontier.slack_reward", edakit.Func(slackReward))
}

// elmore draws an RC chain (KP 1 and 2) or a two-branch tree (KP 3, and harder unguided draws)
// and asks for the Elmore delay at one sink in picoseconds.
func elmore(ctx generator.GeneratorContext) generator.Problem {
	tree := ctx.KPOrdinal == 3 || (ctx.KPOrdinal == 0 && ctx.Difficulty > 0.6)
	if !tree {
		n := edakit.Scaled(ctx, 2, 4)
		rs := make([]int, n)
		cs := make([]int, n)
		var parts []string
		for i := 0; i < n; i++ {
			rs[i] = 50 * edakit.Int(1, 8)
			cs[i] = 5 * edakit.Int(1, 8)
			parts = append(parts, edakit.Sprintf("$R_%d = %d\\ \\Omega$ then $C_%d = %d$ fF", i+1, rs[i], i+1, cs[i]))
		}
		// T_n = sum_k (sum_{j<=k} R_j) C_k for a chain.
		total, shared := 0, 0
		for i := 0; i < n; i++ {
			shared += rs[i]
			total += shared * cs[i]
		}
		ans := edakit.Num(float64(total)/1000, 3)
		q := "An RC chain runs from the driver through " + strings.Join(parts, ", ") +
			edakit.Sprintf(", each capacitor to ground. What is the Elmore delay at the far end, node %d, in picoseconds?", n)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Each capacitor sees all the resistance between it and the driver. Summing shared resistance times capacitance gives %d ohm-femtofarads, which is %s ps.", total, ans),
			map[string]string{"sum_rc": edakit.IntStr(total)})
	}
	// Trunk R0 into node t with Ct, then branch a (Ra, Ca) and branch b (Rb, Cb).
	r0, ct := 50*edakit.Int(1, 6), 5*edakit.Int(1, 6)
	ra, ca := 50*edakit.Int(1, 6), 5*edakit.Int(1, 8)
	rb, cb := 50*edakit.Int(1, 6), 5*edakit.Int(1, 8)
	sink := edakit.Pick([]string{"a", "b"})
	var total int
	if sink == "a" {
		total = r0*(ct+ca+cb) + ra*ca
	} else {
		total = r0*(ct+ca+cb) + rb*cb
	}
	ans := edakit.Num(float64(total)/1000, 3)
	q := edakit.Sprintf("A driver feeds a trunk of $%d\\ \\Omega$ into node t ($C_t = %d$ fF). From t, branch a has $%d\\ \\Omega$ to sink a ($C_a = %d$ fF) and branch b has $%d\\ \\Omega$ to sink b ($C_b = %d$ fF). What is the Elmore delay at sink %s, in picoseconds?",
		r0, ct, ra, ca, rb, cb, sink)
	return edakit.Numeric(q, ans,
		edakit.Sprintf("The trunk resistance charges all downstream capacitance, %d fF, and the branch to sink %s charges only its own sink. The total is %d ohm-femtofarads, which is %s ps.", ct+ca+cb, sink, total, ans),
		map[string]string{"sum_rc": edakit.IntStr(total)})
}

var slackMeaningItems = []edakit.Item{
	{KP: 1, Q: "A path launches on a clock edge at 0 ps. The capture edge is at 800 ps, setup time is 30 ps and data arrives at 790 ps. What is the setup slack?",
		A: "-20 ps", D: []string{"10 ps", "20 ps", "-50 ps"},
		Why: "Required time is 800 - 30 = 770 ps. Slack is required minus arrival: 770 - 790 = -20 ps."},
	{KP: 1, Q: "Which expression defines setup slack at a flip-flop data pin?",
		A: "Required time minus arrival time", D: []string{"Arrival time minus required time", "Clock period minus logic delay, with setup ignored", "Arrival time minus hold time"},
		Why: "Slack is how much later the data could arrive and still be captured, which is required minus arrival."},
	{KP: 1, Q: "Clock skew makes the capture clock arrive 25 ps later than the launch clock. With everything else fixed, what happens to setup slack?",
		A: "It grows by 25 ps", D: []string{"It shrinks by 25 ps", "It does not change", "It grows by 50 ps"},
		Why: "A later capture edge pushes the required time out by 25 ps, so setup slack grows by the same amount. Hold slack shrinks by 25 ps."},
	{KP: 1, Q: "An endpoint's required time is 640 ps and its data arrives at 655 ps. A buffer on the path is removed, saving 22 ps. What is the new setup slack?",
		A: "7 ps", D: []string{"-15 ps", "37 ps", "-37 ps"},
		Why: "The arrival drops to 633 ps, so slack is 640 - 633 = 7 ps. Before the change it was -15 ps."},
	{KP: 2, Q: "A design has endpoint slacks of -12, -3, 4 and 9 ps. What are WNS and TNS?",
		A: "WNS = -12 ps, TNS = -15 ps", D: []string{"WNS = -3 ps, TNS = -15 ps", "WNS = -12 ps, TNS = -2 ps", "WNS = 9 ps, TNS = -15 ps"},
		Why: "WNS is the most negative slack, -12 ps. TNS sums only the negative slacks: -12 + -3 = -15 ps."},
	{KP: 2, Q: "Two designs share a WNS of -8 ps. Design X has a TNS of -9 ps and design Y has a TNS of -4000 ps. What does the difference tell you?",
		A: "Y has many failing endpoints; X has a few local violations", D: []string{"X is harder to close than Y", "Both need exactly the same fix", "TNS is irrelevant once WNS is known"},
		Why: "TNS adds up every negative slack. A large TNS with the same WNS means the violations are widespread rather than concentrated on one path."},
	{KP: 2, Q: "An endpoint shows a slack of +35 ps. What does that mean?",
		A: "The path could slow down by 35 ps and still meet setup", D: []string{"The path fails by 35 ps", "The clock period must grow by 35 ps", "The path violates hold by 35 ps"},
		Why: "Positive setup slack is unused margin at that endpoint."},
	{KP: 2, Q: "Why do timing-driven optimizers watch TNS and not only WNS?",
		A: "Fixing only the worst path can leave hundreds of nearly-as-bad paths untouched", D: []string{"TNS is cheaper to compute than WNS", "WNS ignores the clock period", "TNS measures hold violations only"},
		Why: "WNS reports one path. TNS rewards progress on every violating endpoint, so it guides optimization across the whole design."},
	{KP: 3, Q: "Data can change at the capture flop at the earliest 18 ps after the clock edge, and the hold time is 25 ps. What is the hold slack?",
		A: "-7 ps", D: []string{"7 ps", "43 ps", "-43 ps"},
		Why: "Hold slack is earliest arrival minus the hold requirement: 18 - 25 = -7 ps, a hold violation."},
	{KP: 3, Q: "An engineer upsizes a shared gate to fix a setup violation and a hold violation appears. Why?",
		A: "The faster gate also speeds the short paths through it, so data arrives too early", D: []string{"Upsizing always adds clock skew", "Hold time depends on the clock period", "Setup and hold slack always sum to zero"},
		Why: "Hold fails when data arrives too soon. Speeding a gate helps setup on the long path and hurts hold on any short path through that gate."},
	{KP: 3, Q: "Which change helps a hold violation without touching the clock?",
		A: "Insert delay buffers on the short path", D: []string{"Upsize every gate on the short path", "Reduce the clock period", "Remove the capture flop's setup check"},
		Why: "Hold needs the earliest arrival to grow, so adding delay on the short path is the standard fix."},
	{KP: 3, Q: "Does increasing the clock period fix a hold violation?",
		A: "No: hold compares edges of the same cycle, so the period does not enter", D: []string{"Yes, always", "Only when skew is negative", "Yes, by half the added period"},
		Why: "The hold check is between the launch edge and the same capture edge, so the clock period cancels out of the check."},
}
