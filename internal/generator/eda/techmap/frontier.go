package techmap

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// goldenVerdict computes the gate decision for a candidate mapping change.
func goldenVerdict(ctx generator.GeneratorContext) generator.Problem {
	a0 := float64(edakit.Int(800, 2000))
	d0 := float64(edakit.Int(400, 900))
	budget := float64(edakit.Pick([]int{1, 2, 3}))
	eq := edakit.Pick([]string{"proved", "proved", "proved", "counterexample", "timeout"})
	da := float64(edakit.Int(-80, 20)) / 10
	dd := float64(edakit.Int(-30, 50)) / 10
	a1 := edakit.Round(a0*(1+da/100), 0)
	d1 := edakit.Round(d0*(1+dd/100), 0)
	areaPct := (a1 - a0) / a0 * 100
	delayPct := (d1 - d0) / d0 * 100
	const (
		accept = "Promote: equivalent, delay within budget, area not worse"
		rejEq  = "Reject: the mapped netlist is not equivalent"
		hold   = "Hold: equivalence is unproven, so the QoR numbers cannot be trusted yet"
		rejD   = "Reject: the delay regression exceeds the budget"
		rejA   = "Reject: area got worse with no delay gain to justify it"
	)
	var ans string
	switch {
	case eq == "counterexample":
		ans = rejEq
	case eq == "timeout":
		ans = hold
	case delayPct > budget:
		ans = rejD
	case areaPct > 0 && delayPct >= 0:
		ans = rejA
	default:
		ans = accept
	}
	var dis []string
	for _, o := range []string{accept, rejEq, hold, rejD, rejA} {
		if o != ans {
			dis = append(dis, o)
		}
	}
	q := edakit.Sprintf("An agent proposes a new mapping cost function. On the reference design the baseline mapping has area %s and critical delay %s ps; the candidate gives area %s and delay %s ps. The equivalence check of candidate against the pre-mapping network returned: %s. The promotion rule is: equivalence must be proved, delay may regress by at most %s%%, and area must not grow unless delay improves. What is the verdict?",
		edakit.Num(a0, 0), edakit.Num(d0, 0), edakit.Num(a1, 0), edakit.Num(d1, 0), eq, edakit.Num(budget, 0))
	return edakit.Choice(q, ans,
		edakit.Sprintf("Equivalence comes first: QoR deltas of a wrong netlist mean nothing. Here the check %s, area changed by %s%% and delay by %s%%, so the rule gives this verdict.", eq, edakit.Num(areaPct, 2), edakit.Num(delayPct, 2)),
		dis[:3]...)
}

func goldenDelta(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(3, 4)
	var rows []string
	sum := 0.0
	var pcts []float64
	for i := 0; i < n; i++ {
		a0 := float64(edakit.Int(500, 3000))
		a1 := edakit.Round(a0*(1+float64(edakit.Int(-60, 30))/1000), 0)
		p := (a1 - a0) / a0 * 100
		pcts = append(pcts, p)
		sum += p
		rows = append(rows, edakit.Sprintf("design %d: %s -> %s", i+1, edakit.Num(a0, 0), edakit.Num(a1, 0)))
	}
	mean := sum / float64(n)
	q := "A golden QoR suite reports mapped area, baseline then candidate: " + joinStr(rows) + ". The suite summarizes area as the mean of the per-design percentage changes. What is that mean, in percent? Negative means smaller."
	return edakit.Choice(q, edakit.Num(mean, 2),
		edakit.Sprintf("Per-design changes are %s%%; their mean is %s%%. Averaging percentages, rather than summing raw areas, keeps one large design from dominating the suite.", joinFloats(pcts, 2), edakit.Num(mean, 2)),
		numChoices(mean, 2, -mean, sum, mean/2, mean+1)...)
}

func joinStr(xs []string) string {
	out := ""
	for i, x := range xs {
		if i > 0 {
			out += "; "
		}
		out += x
	}
	return out
}

func joinFloats(xs []float64, p int) string {
	out := ""
	for i, x := range xs {
		if i > 0 {
			out += ", "
		}
		out += edakit.Num(x, p)
	}
	return out
}

var goldenItems = []edakit.Item{
	{KP: 3, Q: "Which check must every golden test of a mapping change include, before any QoR comparison?",
		A: "A formal equivalence check of the mapped netlist against the input network", D: []string{"A comparison of run times", "A check that the cell count did not change", "A diff of the netlist text against the baseline"},
		Why: "A mapper that produces a wrong netlist can show any QoR it likes. Equivalence is the correctness gate; QoR is the quality gate."},
	{KP: 3, Q: "A candidate change improves mean area by 2% on the golden suite but one design gets 9% worse. What should the differential report do?",
		A: "Show the per-design deltas and flag the outlier, rather than only the mean", D: []string{"Report only the mean, since it improved", "Drop the outlier design from the suite", "Average the outlier with the best design"},
		Why: "Means hide regressions. A differential check reports where results moved, so a reviewer can see that one design was hurt."},
	{KP: 3, Q: "Two runs of the same mapping change give different areas on the same design. What does this say about the golden test?",
		A: "The engine is nondeterministic, so it must be seeded and made reproducible before deltas mean anything", D: []string{"The change is good, since one run improved", "Golden tests should take the best of several runs", "Equivalence checking will catch it"},
		Why: "Differential testing assumes the baseline and the candidate are each reproducible. Nondeterminism turns noise into apparent QoR changes."},
	{KP: 3, Q: "A mapping change only alters the cost of one cell class. What should the differential test prove besides equivalence?",
		A: "That only the matches involving that cell class changed, and by plausible amounts", D: []string{"That every gate in the netlist changed", "That the run time is identical", "That the Liberty file parses"},
		Why: "A targeted change should have a targeted effect. Changes elsewhere point to a side effect in the seam or the core."},
	{KP: 3, Q: "Why do mapping golden suites usually report delay from a timer run on the mapped netlist rather than the mapper's own delay estimate?",
		A: "The mapper's internal model is approximate; an independent timer is the reference the change must hold up against", D: []string{"The mapper never computes delay", "Timers are faster than mappers", "Delay does not matter for mapping"},
		Why: "Load is unknown during mapping, so the mapper's numbers are estimates. QoR must be measured by the reference engine."},
	{KP: 3, Q: "The equivalence check of a candidate mapping times out on one design. What is the right action?",
		A: "Hold the change: unproven equivalence is not a pass, so rerun with more effort or a different method", D: []string{"Promote it, since no counterexample was found", "Reject the whole skill permanently", "Remove the design from the suite"},
		Why: "A timeout is not a proof. The gate stays closed until equivalence is established."},
	{KP: 3, Q: "Which baseline should a QoR comparison for a new mapper recipe use?",
		A: "The current promoted recipe run on the same designs, library and constraints", D: []string{"Published numbers from a paper on other designs", "The candidate's own previous run", "The unmapped AIG node count"},
		Why: "A differential comparison holds everything fixed except the change under test."},
	{KP: 3, Q: "An agent's change passes equivalence and improves area by 4% but adds 30% run time. How should the golden harness treat run time?",
		A: "Report it as a tracked metric with its own budget, since it is part of QoR for a tool", D: []string{"Ignore it, since only area and delay matter", "Reject any change that adds run time", "Measure it on a different machine"},
		Why: "Engine run time matters to users. A harness tracks it alongside area and delay, each with a stated tolerance."},
}

var costSeamItems = []edakit.Item{
	{KP: 1, Q: "In a mapping engine, which part belongs in the immutable core rather than behind a seam?",
		A: "Cut enumeration, matching and the covering dynamic program, which guarantee a legal, equivalent cover", D: []string{"The weights that trade area against delay", "The per-cell cost adjustments a customer requests", "The choice of which library cells to exclude"},
		Why: "Correctness of the cover lives in the core. The costs that rank legal matches vary by customer and belong at a seam."},
	{KP: 1, Q: "Where should a library-specific cost function plug into a cut-based mapper?",
		A: "At the point where the DP asks for the cost of a match: a cost-provider interface called with the match context", D: []string{"Inside the cut enumeration loop, as an edit to the merge code", "In the Liberty parser", "In the equivalence checker"},
		Why: "The DP asks one question per candidate match. Putting the seam there lets cost policies change without touching enumeration or matching."},
	{KP: 1, Q: "A match-cost seam receives a context. Which item is least useful to include in it?",
		A: "The source code of the cell's layout generator", D: []string{"The cell and its Vt class", "The estimated load and required time at the node", "The arrival times of the match's inputs"},
		Why: "A cost provider ranks matches by area, delay, power and slack. It needs the cell, timing context and load estimate, not layout code."},
	{KP: 1, Q: "Why should the cost seam be called once per node with all candidate matches, rather than once per match?",
		A: "Coarse calls keep the hot loop inside the core, where it can be optimized, and cut call overhead", D: []string{"Per-match calls would change the mapping result", "The core cannot enumerate matches itself", "Per-node calls make the mapper nondeterministic"},
		Why: "Seams should be coarse-grained: a batch of matches per call lets the core keep its tight loops and caches."},
	{KP: 2, Q: "A customer wants high-Vt cells penalized on paths with less than 50 ps of slack. Where does this belong?",
		A: "In a cost-provider plug-in that adds a penalty using the slack in the match context", D: []string{"In a core edit to the covering DP", "In the cut enumeration limit K", "In the Liberty file, by editing the cell areas"},
		Why: "The behaviour depends on context (Vt class and slack) that the seam already exposes, so it is a small plug-in, not a core change."},
	{KP: 2, Q: "A new cost provider makes the mapper choose different cells, but the mapped netlist fails equivalence. What does this indicate?",
		A: "The seam leaks: cost policy was able to affect legality, which the core should guarantee regardless of cost", D: []string{"The cost provider is working as intended", "The library is incomplete", "Equivalence checking is too strict for mapping"},
		Why: "Costs should only rank legal matches. If a cost change breaks equivalence, the core let policy decide something it must own."},
	{KP: 2, Q: "What must the cost seam's interface be, to let plug-ins be compiled separately and swapped safely?",
		A: "Narrow and stable: plain data in, cost numbers out, no hidden global state", D: []string{"Wide, exposing all internal data structures", "Changed every release to add features", "Based on parsing log files"},
		Why: "A frozen, narrow interface is what lets plug-ins be developed, tested and replaced without recompiling the core."},
	{KP: 2, Q: "A plug-in cost provider reads a random number generator without a seed. What breaks?",
		A: "Determinism: the same design can map differently run to run, so golden and differential tests become unreliable", D: []string{"Nothing, randomness improves results", "The Liberty parser", "The NPN canonicalization"},
		Why: "Seams must be deterministic and seedable, or differential QoR checks cannot separate change from noise."},
	{KP: 3, Q: "Mapping should use different delay models for the early fast pass and the final pass. How should the engine support this?",
		A: "Two implementations of the same delay-model interface, chosen by configuration", D: []string{"Two copies of the mapper core", "An if-statement on the design name inside the DP", "Editing the Liberty file between passes"},
		Why: "One seam with interchangeable providers lets a fast model and an accurate model coexist, selected without code edits."},
	{KP: 3, Q: "Which test proves a new cost provider did what was intended?",
		A: "Equivalence plus a differential QoR report showing that the changed matches are the ones the policy targets", D: []string{"A clean compile", "A larger number of cuts enumerated", "A shorter log file"},
		Why: "Equivalence shows the cover is legal; the differential report shows the policy moved the intended decisions and nothing else."},
	{KP: 3, Q: "A team wants a cost provider that calls a learned power model. What is the main risk to manage at the seam?",
		A: "The model's latency and accuracy: it is called for many matches, and its errors steer the mapper", D: []string{"It cannot return numbers", "It changes the subject graph", "It forbids area recovery"},
		Why: "A seam can host a learned model, but the call volume and the model's error both feed directly into covering decisions."},
	{KP: 3, Q: "Why keep the cost seam at the match level rather than letting plug-ins rewrite the final netlist?",
		A: "Costs at the match level steer the optimizer, while netlist rewrites bypass it and can break the guarantees of the cover", D: []string{"Netlist rewrites are always faster", "Match-level costs cannot express area", "Plug-ins cannot read netlists"},
		Why: "The seam should feed the optimization, not patch its output, so legality and optimality reasoning stay inside the core."},
}

var recipeItems = []edakit.Item{
	{KP: 1, Q: "An agent has three mapping skills with manifests: 'delay-first standard-cell mapping with area recovery', 'area-only standard-cell mapping', and 'depth-optimal K-LUT mapping'. The target is an FPGA. Which does it select?",
		A: "Depth-optimal K-LUT mapping", D: []string{"Delay-first standard-cell mapping", "Area-only standard-cell mapping", "Run all three and keep the smallest netlist"},
		Why: "Selection matches the target to what each manifest advertises. Only the LUT mapper produces an FPGA netlist."},
	{KP: 1, Q: "What does an agent read to choose a mapping recipe?",
		A: "The recipes' manifests: what each does, its inputs and outputs, and which parameters may be tuned", D: []string{"The source code of each mapper", "The full log of every past run", "The Liberty file of each cell"},
		Why: "Selection is a cheap text decision over manifests; the heavy compute stays in the engines."},
	{KP: 1, Q: "A block has 200 ps of negative slack after mapping with the area recipe. Which recipe should the agent try next?",
		A: "The delay-oriented recipe with area recovery", D: []string{"The same area recipe with a different seed only", "The LUT recipe", "No mapping, keeping the AIG"},
		Why: "The metric that is failing, timing, decides the next skill. A delay-first mapper with area recovery addresses it directly."},
	{KP: 1, Q: "Which manifest field matters most for an agent deciding whether a recipe can run on a block with a strict area cap?",
		A: "The objective it optimizes and the parameters that trade area for delay", D: []string{"The author's name", "The size of the binary", "The date it was compiled"},
		Why: "The agent matches objectives and tunables to the block's constraint."},
	{KP: 2, Q: "An agent may tune a recipe's parameters. Which is a sensible tunable to expose for a cut-based mapper?",
		A: "The cut size K and the number of priority cuts per node", D: []string{"The definition of a k-feasible cut", "The NPN canonical form algorithm", "The equivalence-checking engine"},
		Why: "K and C trade quality against run time and memory. Definitions and checkers are core and stay fixed."},
	{KP: 2, Q: "An agent runs four recipes in parallel on a block and picks the best by area. What must hold for the choice to be valid?",
		A: "Every candidate passed equivalence and met the timing constraint before area was compared", D: []string{"The recipes were written by the same author", "The run times were equal", "The smallest netlist is always chosen regardless of timing"},
		Why: "Area is only comparable among legal, constraint-meeting results."},
	{KP: 2, Q: "A recipe works well on datapath-heavy blocks and poorly on control logic. How should this knowledge be recorded?",
		A: "In the recipe's manifest and the domain model, as a condition for when to select it", D: []string{"In a comment inside the mapper's core loop", "Nowhere, since agents rediscover it", "By deleting the recipe"},
		Why: "The flywheel captures when a skill works, so the next selection uses that knowledge."},
	{KP: 2, Q: "An agent alternates between two mapping recipes because each fixes what the other breaks. What should the harness add?",
		A: "A stopping rule based on QoR progress and an iteration budget", D: []string{"More recipes to alternate between", "A larger cut size", "A random restart each iteration"},
		Why: "Agentic loops need termination criteria; otherwise they oscillate without converging."},
	{KP: 3, Q: "Why should an agent try a fast mapping estimate before an accurate mapping-plus-timing run?",
		A: "The fast estimate filters recipes cheaply, and only promising ones pay for the accurate run", D: []string{"Fast estimates are always more accurate", "Accurate runs cannot be automated", "Agents cannot call timers"},
		Why: "Fast and accurate engines work as a pair: the fast one prunes the search, the accurate one decides."},
	{KP: 3, Q: "An agent's chosen recipe produced the smallest area under the mapper's own delay estimate, but timing fails after placement. What went wrong in the selection?",
		A: "It trusted a fast proxy for delay instead of checking with the reference timer before committing", D: []string{"It used the wrong library", "Area mapping is always wrong", "Placement ignores mapping"},
		Why: "A proxy metric can be gamed by the choice. The decision must be checked with the engine the result will be judged by."},
	{KP: 3, Q: "A customer asks for 'smaller area at the same speed'. Which step comes first for the agent?",
		A: "Measure the current area and delay with the reference flow to fix the baseline", D: []string{"Pick the area-only recipe immediately", "Increase K for the LUT mapper", "Edit the mapper core"},
		Why: "A delta needs a baseline; without one the agent cannot prove it met the request."},
	{KP: 3, Q: "A tool-grounded reward for an agent that picks mapping recipes uses only cell count. What is the risk?",
		A: "The agent learns to minimize cell count with large complex cells, which can hurt area and delay", D: []string{"The reward cannot be computed", "Cell count is equivalent to area", "The agent will always pick the LUT mapper"},
		Why: "A reward that misses part of QoR invites gaming. Rewards should use the real objectives, measured by the reference engines."},
}

var changeTierItems = []edakit.Item{
	{KP: 1, Q: "A user asks the mapping engine to exclude all cells with drive strength 1 from the library. Which change tier is this?",
		A: "Tier 1, configuration: a don't-use list, with no recompile", D: []string{"Tier 2, a new cost plug-in", "Tier 3, a core edit to matching", "It cannot be done"},
		Why: "Excluding cells is a table or list change that the engine reads at run time."},
	{KP: 1, Q: "A user wants the cut size raised from 5 to 6 for a LUT mapping. Which tier?",
		A: "Tier 1, configuration", D: []string{"Tier 2, plug-in", "Tier 3, core edit", "A new skill"},
		Why: "K is a parameter of the existing mapper."},
	{KP: 1, Q: "A user wants the area-delay weight of the mapper changed from 0.5 to 0.8. Which tier?",
		A: "Tier 1, configuration", D: []string{"Tier 2, plug-in", "Tier 3, core edit", "Tier 3, because it changes results"},
		Why: "Weights are parameters. Changing results is fine; changing code is what the tier measures."},
	{KP: 1, Q: "A user asks for a cost that penalizes cells whose pins are hard to reach in their specific layout style, using a per-cell table they supply. Which tier?",
		A: "Tier 2, a cost-provider plug-in that reads the table", D: []string{"Tier 1, by editing the Liberty area", "Tier 3, an edit to the DP", "No change needed"},
		Why: "A new cost term computed from new data is a small implementation at the cost seam."},
	{KP: 2, Q: "A user wants matches that cross a clock-gating boundary to be forbidden. The mapper has no hook for match legality. Which tier is this, and what does it suggest?",
		A: "Tier 3, a core edit, and it suggests the engine needs a legality seam for match filtering", D: []string{"Tier 1, a parameter", "Tier 2, using the cost seam with an infinite cost", "Tier 1, by deleting the clock gates"},
		Why: "Legality belongs to the core. If such requests recur, a match-filter extension point should be added so they become Tier 2."},
	{KP: 2, Q: "A user's request can be met with the cost seam by assigning an enormous cost to forbidden matches. Why is this risky?",
		A: "If every legal alternative is also expensive or missing, the mapper may still pick a forbidden match, so legality is not guaranteed", D: []string{"Large numbers slow the mapper", "Costs cannot be large", "It is not risky"},
		Why: "Costs rank; they do not forbid. Legality constraints need a mechanism that removes options."},
	{KP: 2, Q: "A user wants the mapper to call their in-house delay model instead of Liberty tables. The engine has a delay-model interface. Which tier?",
		A: "Tier 2, a new delay-model implementation behind the existing interface", D: []string{"Tier 1, a parameter", "Tier 3, rewriting the mapper", "It requires a new mapper"},
		Why: "The seam exists, so a new provider is a plug-in."},
	{KP: 2, Q: "A user asks for a different cut enumeration algorithm that enumerates cuts on the fly. Which tier?",
		A: "Tier 3, a core edit with human review and full regression", D: []string{"Tier 1, a parameter", "Tier 2, a cost plug-in", "Tier 1, by lowering K"},
		Why: "Cut enumeration is part of the core's search; changing it changes the engine itself."},
	{KP: 3, Q: "Over a year, most mapping requests land in Tier 3. What does this say about the engine?",
		A: "The seams are in the wrong place: common changes should land at configuration or plug-in tiers", D: []string{"The users are asking for too much", "Tier 3 edits are the safest", "The engine is well designed"},
		Why: "A seam is well placed when the common kind of change never needs a core edit."},
	{KP: 3, Q: "Which verification should a Tier 2 mapping plug-in pass before promotion?",
		A: "Golden tests with equivalence checks plus a differential QoR comparison against the baseline", D: []string{"Only a compile", "A code review of the core", "No tests, since the core is untouched"},
		Why: "Plug-ins change results, so they need the same correctness and QoR gates as any change."},
	{KP: 3, Q: "A Tier 1 change (a new don't-use list) causes the mapper to fail on a design. What failed?",
		A: "The library became incomplete for some function, for example no inverter remains, so no cover exists", D: []string{"The core was edited", "The plug-in ABI broke", "Equivalence checking failed to run"},
		Why: "Configuration can still create infeasible problems. The engine should validate that the library stays functionally complete."},
	{KP: 3, Q: "An agent proposes a core edit to fix a Tier 2 request because it is quicker for it. What should the harness do?",
		A: "Steer it to the smallest safe tier: implement the plug-in, and require human approval for any core edit", D: []string{"Accept the core edit if it compiles", "Reject the request", "Let the agent edit both"},
		Why: "The model edits at the smallest safe tier. Core edits need human-approved diffs and full regression."},
}
