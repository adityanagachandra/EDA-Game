package hls

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func qorSurrogate(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		a := 10 * edakit.Int(20, 500)
		p := a + 10*edakit.Int(-a/40, a/40)
		if p == a {
			p += 30
		}
		e := math.Abs(float64(p-a)) / float64(a) * 100
		q := edakit.Sprintf("A fast QoR surrogate predicts %d cycles of latency for a pragma configuration; full HLS synthesis reports %d cycles. What is the absolute percentage error of the prediction? Round to 1 decimal place.", p, a)
		return num(q, edakit.Num(e, 1),
			edakit.Sprintf("|%d - %d| / %d × 100 = %s%%.", p, a, a, edakit.Num(e, 1)),
			map[string]string{"ape": edakit.Num(e, 1)})
	case 2:
		var p []string
		tot := 0.0
		for i := 0; i < 4; i++ {
			a := 10 * edakit.Int(20, 300)
			pr := a + 10*edakit.Int(-a/50-1, a/50+1)
			tot += math.Abs(float64(pr-a)) / float64(a) * 100
			p = append(p, edakit.Sprintf("D%d predicted %d, actual %d", i+1, pr, a))
		}
		m := tot / 4
		q := "A QoR surrogate is checked against full synthesis on four designs (latency in cycles): " + strings.Join(p, "; ") + ". What is its mean absolute percentage error? Round to 1 decimal place."
		return num(q, edakit.Num(m, 1),
			edakit.Sprintf("Average the four values of |predicted - actual| / actual × 100: %s%%.", edakit.Num(m, 1)),
			map[string]string{"mape": edakit.Num(m, 1)})
	default:
		n := 5
		act := edakit.Shuffle([]int{1, 2, 3, 4, 5})
		pred := edakit.Shuffle([]int{1, 2, 3, 4, 5})
		av := make([]int, n)
		pv := make([]int, n)
		for i := 0; i < n; i++ {
			av[i] = 100*act[i] + 10*edakit.Int(0, 9)
			pv[i] = 100*pred[i] + 10*edakit.Int(0, 9)
		}
		disc := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if (av[i]-av[j])*(pv[i]-pv[j]) < 0 {
					disc++
				}
			}
		}
		var p []string
		for i := 0; i < n; i++ {
			p = append(p, edakit.Sprintf("D%d predicted %d, actual %d", i+1, pv[i], av[i]))
		}
		q := "A design-space explorer uses a surrogate only to rank candidates. On five designs (latency in cycles): " + strings.Join(p, "; ") + ". How many of the 10 design pairs does the surrogate order differently from full synthesis?"
		return num(q, edakit.IntStr(disc),
			edakit.Sprintf("A pair is discordant when the predicted order and the actual order disagree. Checking all 10 pairs finds %d, so Kendall's tau is (10 - 2 × %d) / 10 = %s.", disc, disc, edakit.Num(float64(10-2*disc)/10, 1)),
			map[string]string{"discordant": edakit.IntStr(disc)})
	}
}

var seamItems = []edakit.Item{
	{KP: 1, Q: "In an HLS engine built on the four-layer model, which part belongs to the immutable core?",
		A:   "Dependence legality and the modulo reservation check that keep every schedule correct",
		D:   []string{"The priority function that list scheduling uses", "The default II target", "The hook that logs each scheduling decision"},
		Why: "Correctness rules are the core. Priorities are policy, targets are parameters and logging is an observation hook."},
	{KP: 1, Q: "Which change to an HLS engine is a strategy or policy change?",
		A:   "A new list-scheduling priority that favors operations feeding memory stores",
		D:   []string{"Raising the default unroll limit from 8 to 16", "Fixing dependence analysis to respect loop-carried memory dependences", "Adding a log line when a loop's II is chosen"},
		Why: "A priority decides how the legal schedule is chosen, which is policy. The others are a parameter, a core fix and a hook."},
	{KP: 1, Q: "A user sets a 4 ns target clock and II = 2 on a loop. Which layer do these values sit in?",
		A: "Parameters", D: []string{"Immutable core", "Strategy or policy", "Observation hooks"},
		Why: "They are values the engine reads, changed without recompiling anything."},
	{KP: 1, Q: "Which of these is an observation hook in an HLS engine?",
		A:   "A callback that exposes each loop's achieved II, its limiting resource or recurrence, and the chosen schedule",
		D:   []string{"The SDC solver", "The left-edge register binder", "The C front-end parser"},
		Why: "A hook lets an agent read live engine state; the others are engine components."},
	{KP: 2, Q: "A customer wants every loop in a kernel pipelined with a target II of 2 and unrolling capped at 4. Which tier is this?",
		A: "Tier 1, configuration", D: []string{"Tier 2, plug-in", "Tier 3, core edit", "None, because it cannot be expressed"},
		Why: "These are directive and parameter values the engine already accepts."},
	{KP: 2, Q: "A team wants a binding cost that penalizes sharing a multiplier between operations whose operands come from different memory banks. Which tier is this?",
		A:   "Tier 2, a plug-in at the binder's cost seam",
		D:   []string{"Tier 1, a parameter change", "Tier 3, an edit to the scheduler's legality check", "Tier 1, a new directive value"},
		Why: "It is new behavior, but it only changes how legal bindings are ranked, so it belongs at a cost-function extension point."},
	{KP: 2, Q: "A request would let the scheduler ignore a loop-carried memory dependence to reach a lower II. Which tier is it, and what does it require?",
		A:   "Tier 3: it changes correctness rules, so it needs a human-approved diff and full regression",
		D:   []string{"Tier 1, because II is a parameter", "Tier 2, a new priority plug-in", "No review, as long as co-simulation is skipped"},
		Why: "Dropping a dependence can make the hardware compute a different function. That is a core change, if it is allowed at all."},
	{KP: 2, Q: "Where should the seam for customer-specific operator latency and delay tables sit?",
		A:   "Behind an operator-characterization interface that the scheduler queries per operation, width and target",
		D:   []string{"As hard-coded constants inside the SDC solver loop", "Inside the C front-end parser", "In the string templates of the RTL emitter"},
		Why: "Scheduling, chaining and binding all need latencies; one narrow interface lets new tables plug in without touching those cores."},
	{KP: 3, Q: "A new binding plug-in is called through a virtual call once per operation inside the binder's innermost loop. What does the seam guidance recommend instead?",
		A:   "Call the plug-in coarsely, once per batch of operations, so the core loop stays fast",
		D:   []string{"Call it twice per operation for safety", "Move the plug-in into the core", "Remove the binder"},
		Why: "Coarse calls keep the hot loop inlinable while still letting the plug-in decide."},
	{KP: 3, Q: "What must the proof of a new list-scheduling priority plug-in show?",
		A:   "The golden kernels still pass co-simulation and the QoR diff moves latency and area only where intended",
		D:   []string{"Only that it compiles", "That the new priority is used on at least one loop", "That the plug-in has no parameters"},
		Why: "Golden tests prove correctness is unchanged; the differential QoR report proves the effect is the intended one."},
	{KP: 3, Q: "A plug-in reads a global variable left behind by an earlier run. Which principle does this violate?",
		A:   "Plug-ins take plain data, keep no hidden state, and stay deterministic and seedable",
		D:   []string{"Plug-ins must be written in C", "Plug-ins must take no parameters", "Plug-ins must call the core recursively"},
		Why: "Hidden state makes results depend on history, so differential tests stop being reproducible."},
	{KP: 3, Q: "Why keep the operator latency model behind one narrow, stable interface rather than letting each pass read the library tables directly?",
		A:   "So new targets and characterizations plug in without touching the scheduler, binder or estimator",
		D:   []string{"Because tables cannot be read twice", "So the scheduler can ignore latencies", "To make the core larger"},
		Why: "A frozen interface localizes change: a new FPGA family or cell library becomes a Tier 1 or Tier 2 change."},
}

var pragmaAgentItems = []edakit.Item{
	{KP: 1, Q: "Why can an agent not simply run HLS on every pragma configuration of a real kernel?",
		A:   "The space grows multiplicatively with each pragma, and each synthesis run takes minutes to hours",
		D:   []string{"HLS tools accept only one pragma per kernel", "Pragmas do not change QoR", "The space never exceeds ten points"},
		Why: "A few loops and arrays with a handful of options each already give thousands of points, so search must be guided."},
	{KP: 1, Q: "An agent proposes ARRAY_PARTITION complete on a one-million-element array. What should the harness do?",
		A:   "Reject it with a cheap resource check before spending a synthesis run",
		D:   []string{"Run full synthesis and see what happens", "Accept it because pragmas are only hints", "Shrink the array in the C source"},
		Why: "Complete partitioning turns every element into registers. A cheap estimate shows it cannot fit, which saves a run that would fail."},
	{KP: 1, Q: "What should a pragma-search agent read after each HLS run to choose the next experiment?",
		A:   "The achieved II, latency and resource use, and the reported reason a loop missed its target II",
		D:   []string{"Only the pass or fail exit code", "The generated Verilog line by line", "The host CPU load"},
		Why: "The report names the binding constraint, which tells the agent which pragma could actually help."},
	{KP: 1, Q: "Which change does a pragma-search agent make at Tier 1?",
		A:   "Setting directive values such as II targets, unroll factors and partition factors",
		D:   []string{"Rewriting the scheduler", "Editing the binder's legality checks", "Changing the RTL emitter"},
		Why: "Directives are configuration; the engine's code is untouched."},
	{KP: 2, Q: "A run reports 'II = 4, limited by the memory ports of array A'. What is the most targeted next move?",
		A:   "Partition A cyclically or restructure its accesses, then re-run",
		D:   []string{"Raise the target clock frequency", "Unroll a different loop", "Remove the PIPELINE directive"},
		Why: "The report names a resource bound on A's ports, so adding banks attacks exactly that bound."},
	{KP: 2, Q: "A run reports 'II = 5, limited by a carried dependence through s'. What will partitioning arrays do?",
		A:   "Nothing for II: the recurrence bound remains, so the agent must break or reassociate the recurrence",
		D:   []string{"Lower the II to 1", "Raise the II to 10", "Remove the dependence automatically"},
		Why: "RecMII comes from the cycle through s, which no amount of memory bandwidth changes."},
	{KP: 2, Q: "Why does the harness keep a baseline QoR for each kernel?",
		A:   "So each proposal is judged by its latency and area change against a known reference",
		D:   []string{"Because HLS tools refuse to run without one", "To store the C source", "To set the clock period"},
		Why: "A differential comparison against the baseline shows whether a proposal helped and by how much."},
	{KP: 2, Q: "Two consecutive proposals give identical reports. What should a good search policy conclude?",
		A:   "Those pragmas did not touch the binding constraint, so it should search elsewhere",
		D:   []string{"Repeat the same proposal", "Raise every unroll factor at once", "Stop the search immediately"},
		Why: "No change in QoR means the limiting resource or recurrence lies outside what was varied."},
	{KP: 3, Q: "An agent finds a configuration with half the latency but never ran co-simulation. What is the status of the result?",
		A:   "Unverified: the QoR is an estimate and the RTL's behavior is unchecked",
		D:   []string{"Ready for signoff", "Proven correct by the HLS report", "Final, because the latency improved"},
		Why: "Only co-simulation against the C reference shows the new microarchitecture still computes the right function."},
	{KP: 3, Q: "Rewarding a search agent only on reported latency leads it to unroll everything maximally. What is missing from the reward?",
		A:   "Area or resource limits, so the objective trades latency against cost",
		D:   []string{"A larger weight on latency", "A penalty for using any pragma", "Nothing, since minimum latency is the goal"},
		Why: "A single-metric reward is gamed by spending unlimited area; the objective must reflect the real constraints."},
	{KP: 3, Q: "Reported HLS latencies look far better than post-implementation results. What explains it, and what should the agent do?",
		A:   "HLS estimates miss routing delay and congestion, so finalists must be validated with implementation runs",
		D:   []string{"The implementation tool is wrong and should be ignored", "Implementation results should be discarded", "The clock in the report should be lowered"},
		Why: "The fast estimate is a proxy. The loop recurses into the accurate engine for the candidates that matter."},
	{KP: 3, Q: "Which proposal should a sandboxed pragma agent never apply on its own?",
		A:   "Editing the C source to delete a dependence it judges unnecessary",
		D:   []string{"Changing an unroll factor", "Setting a loop's II target", "Choosing a partition factor"},
		Why: "Directive values keep the function; source edits that remove dependences can change it and need human review."},
}

var llmHLSItems = []edakit.Item{
	{KP: 1, Q: "An LLM translates a Python reference into C++ for HLS and uses std::vector with push_back. What will the HLS front end typically do?",
		A:   "Reject it, because a growing container has no static hardware size",
		D:   []string{"Synthesize a RAM that grows at run time", "Convert it into a FIFO automatically", "Silently ignore the calls"},
		Why: "Dynamic allocation is outside the synthesizable subset; the code must use fixed-size arrays or streams."},
	{KP: 1, Q: "Which construct is most likely to make an LLM-generated HLS kernel unsynthesizable?",
		A:   "Recursion with data-dependent depth",
		D:   []string{"A for-loop with a constant bound", "A fixed-size local array", "An arbitrary-precision integer type"},
		Why: "The others are standard HLS-C; unbounded recursion needs an unbounded stack."},
	{KP: 1, Q: "An LLM writes 'while (err > tol)' for an iterative solver. What should be added for HLS?",
		A:   "A bound on the trip count so the tool can budget and report worst-case latency",
		D:   []string{"Nothing, because HLS ignores while-loops", "A directive to unroll it fully", "A switch to double precision"},
		Why: "A data-dependent loop can be synthesized, but without a bound its latency is unknown and it may never terminate on bad data."},
	{KP: 1, Q: "Which check costs least and should run first on an LLM's HLS-C?",
		A:   "Compile and run the C testbench against the reference outputs",
		D:   []string{"Full place and route", "Gate-level simulation", "A silicon tape-out"},
		Why: "C simulation takes seconds and catches functional errors before any synthesis time is spent."},
	{KP: 2, Q: "In an LLM-to-HLS-C harness, what feedback goes back to the model when synthesis fails?",
		A:   "The tool's error messages with the offending source lines",
		D:   []string{"Only a score of zero", "The generated Verilog", "Nothing, and the model retries blind"},
		Why: "Specific errors let the model make a small targeted edit, which is the point of the harness loop."},
	{KP: 2, Q: "The C testbench passes but HLS reports II = 8 instead of the requested 1. What should the next prompt include?",
		A:   "The reported reason for the II violation and the loop it refers to",
		D:   []string{"A request to rewrite the testbench", "The full RTL netlist", "An instruction to lower the clock"},
		Why: "The report names the resource or recurrence that sets the II; that is the information the model needs to restructure the loop."},
	{KP: 2, Q: "Why is a C testbench pass not enough to accept an LLM's HLS-C?",
		A:   "The C can rely on behavior the hardware does not reproduce, such as uninitialized reads, so RTL co-simulation must pass too",
		D:   []string{"C testbenches never actually run", "HLS ignores C semantics", "Testbenches check only timing"},
		Why: "Co-simulation runs the same vectors on the generated RTL and catches mismatches the C execution hid."},
	{KP: 2, Q: "Which order of checks gives the cheapest failing signal first?",
		A:   "C compile, C simulation, HLS synthesis, co-simulation",
		D:   []string{"Co-simulation, C compile, HLS synthesis, C simulation", "HLS synthesis, co-simulation, C compile, C simulation", "C simulation, co-simulation, C compile, HLS synthesis"},
		Why: "Each stage costs more than the one before and needs it to pass, so the loop runs them in increasing cost."},
	{KP: 3, Q: "An LLM 'optimizes' a kernel by removing a bounds check that the testbench never triggers. What must catch this?",
		A:   "Tests that cover boundary inputs, compared against the reference model",
		D:   []string{"The HLS resource report", "A higher unroll factor", "Lint warnings alone"},
		Why: "A check that the tests never reach can be deleted without any test failing, so coverage of the boundary cases is the defense."},
	{KP: 3, Q: "Which metric best measures an LLM-to-HLS-C system?",
		A:   "The fraction of tasks that pass C simulation and co-simulation, with QoR against a baseline",
		D:   []string{"The number of lines of code generated", "Text similarity to a human-written kernel", "How fast the model responds"},
		Why: "Tool-grounded checks measure correctness and quality; text similarity rewards code that merely looks right."},
	{KP: 3, Q: "An LLM uses float where the reference uses 16-bit fixed point. What is the risk?",
		A:   "Results differ numerically and the float hardware is much larger, so the comparison needs a tolerance agreed with the spec",
		D:   []string{"None, because float is always exact", "The design will not compile", "Fixed point and float give bit-identical results"},
		Why: "Changing number format changes both arithmetic results and area; the acceptance test must reflect the specification."},
	{KP: 3, Q: "A team drives its HLS harness with an external API model. What does the IP-boundary principle require?",
		A:   "Deciding explicitly which design data may leave the site, and keeping tools, vectors and results in the sandbox",
		D:   []string{"Sending the whole repository to the model for context", "Disabling the sandbox for speed", "Skipping co-simulation to reduce data movement"},
		Why: "The boundary is part of the architecture: the harness controls what crosses it rather than leaving it to chance."},
}

var cosimItems = []edakit.Item{
	{KP: 1, Q: "What does C/RTL co-simulation run?",
		A:   "The original C testbench, driving the generated RTL through adapters and comparing its outputs",
		D:   []string{"A new hand-written RTL testbench only", "The C code twice with different compilers", "A gate-level timing simulation"},
		Why: "The same vectors drive both models, so any output difference is attributable to the generated hardware."},
	{KP: 1, Q: "Why is the C model the golden reference in an HLS flow?",
		A:   "It is the specification the RTL was generated from, and it runs fast enough to cover many vectors",
		D:   []string{"It is cycle-accurate", "It models the clock tree", "It includes routing delays"},
		Why: "The RTL is supposed to compute exactly the C function, so the C execution defines the right answers."},
	{KP: 1, Q: "Co-simulation passes. What has it not shown?",
		A:   "That the design meets timing after implementation",
		D:   []string{"That outputs match the C model on the tested vectors", "That the interface handshakes work on those tests", "That the RTL finishes on those inputs"},
		Why: "Co-simulation is functional and cycle-level at RTL; physical timing needs implementation and static timing analysis."},
	{KP: 1, Q: "Besides pass or fail, what does co-simulation measure that the HLS report only estimates?",
		A:   "The actual latency and interval in cycles on the test inputs",
		D:   []string{"Silicon power", "Wire capacitance", "The layout's design-rule count"},
		Why: "Variable-latency loops and stalls make the report a range; running the RTL gives the measured cycle counts."},
	{KP: 2, Q: "A new scheduling heuristic is proposed for the HLS engine. What is the differential test?",
		A:   "Run the golden kernel suite before and after; every kernel must still pass co-simulation and QoR changes must be the intended ones",
		D:   []string{"Run one kernel once", "Require the Verilog text to be identical", "Check only that the tool still starts"},
		Why: "A scheduler change should change schedules, so text equality is wrong; correctness plus an explained QoR diff is the bar."},
	{KP: 2, Q: "When is byte-identical RTL the right golden criterion for an engine change?",
		A:   "When the change is meant to preserve behavior, such as a refactor of the RTL emitter",
		D:   []string{"Always, for every scheduler change", "Never", "Only for changes to pragmas"},
		Why: "A behavior-preserving refactor should produce the same output, so any diff is a regression."},
	{KP: 2, Q: "A binding change keeps every co-simulation passing but area rises 9% on two kernels. What should the harness do?",
		A:   "Flag the QoR regression for review even though correctness holds",
		D:   []string{"Accept it because the tests pass", "Revert all binding code", "Ignore area entirely"},
		Why: "Golden tests guard correctness; the QoR baseline guards quality, and both must be checked."},
	{KP: 2, Q: "What makes a kernel a good member of the golden suite?",
		A:   "It exercises distinct features, such as recurrences, partitioned arrays, streams and branches, with known-good outputs",
		D:   []string{"It is as large as possible", "Its QoR never changes", "It has no testbench"},
		Why: "Coverage of features, not size, is what makes a regression suite catch breakage."},
	{KP: 3, Q: "A testbench drives a stream interface but never stalls the output. Which class of bug can slip through?",
		A:   "Back-pressure bugs where data is lost when the consumer is not ready",
		D:   []string{"Arithmetic overflow in the adder", "Wrong loop bounds", "Register binding conflicts"},
		Why: "If ready is always high the hold-while-stalled path never runs, so its bugs stay hidden."},
	{KP: 3, Q: "Co-simulation passes on 10 vectors but the design fails in the field on others. What improves the golden test?",
		A:   "Randomized and coverage-directed vectors, including corner values, checked against the C model",
		D:   []string{"Removing the failing vectors", "Running the same 10 vectors more often", "Lowering the clock"},
		Why: "A golden test is only as strong as its inputs; the C model can label any number of new vectors."},
	{KP: 3, Q: "Floating-point results differ in the last bit between C and RTL after a permitted reassociation. What should the comparison use?",
		A:   "A tolerance justified by the numerical spec, keeping bit-exact comparison for integer paths",
		D:   []string{"Bit-exact comparison everywhere", "No comparison for floating point", "A tolerance of 50%"},
		Why: "Reassociation legitimately changes rounding, so the check must allow it without hiding real errors."},
	{KP: 3, Q: "Why is co-simulation, not the HLS report alone, the right reward signal for a learned HLS-C generator?",
		A:   "The report rewards QoR that a broken design can also achieve; co-simulation grounds the reward in correct behavior",
		D:   []string{"Co-simulation is faster than the report", "The report is not machine-readable", "Rewards must be binary"},
		Why: "A reward on estimated latency alone can be maximized by code that computes the wrong function, which is reward hacking."},
}
