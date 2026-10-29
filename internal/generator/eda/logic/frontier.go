package logic

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerFrontier(reg *generator.Registry) {
	reg.Register("logic.frontier.spec_to_rtl", &edakit.Bank{Items: specToRTLItems})
	reg.Register("logic.frontier.golden_checks", &edakit.Bank{Items: goldenItems})
	dseBank := &edakit.Bank{Items: dseItems}
	reg.Register("logic.frontier.arith_dse", &edakit.Mix{Parts: []generator.Generator{retryTies(dseFilter), dseBank, dseBank}})
	reg.Register("logic.frontier.generator_tiers", &edakit.Bank{Items: tierItems})
	reg.Register("logic.frontier.rtl_reward", &edakit.Bank{Items: rewardItems})
}

type adderArch struct {
	name         string
	depth, nodes int
	fanout       int
}

// dseFilter asks which prefix architecture an agent's fast model should pick first: the fewest
// nodes among those meeting a depth budget and a fanout limit.
func dseFilter(ctx generator.GeneratorContext) generator.Problem {
	for {
		L := edakit.Pick([]int{4, 5, 6})
		n := 1 << uint(L)
		archs := []adderArch{
			{"Ripple-carry", n, n, 2},
			{"Brent-Kung", 2*L - 1, 2*n - 2 - L, 2},
			{"Sklansky", L, n / 2 * L, n / 2},
			{"Kogge-Stone", L, n*L - n + 1, 2},
		}
		d := 5 * edakit.Int(4, 10)
		budget := d * edakit.Int(L, 2*L+2)
		fanLimit := edakit.Pick([]int{4, 8, 64})
		best := -1
		for i, a := range archs {
			if a.depth*d <= budget && a.fanout <= fanLimit && (best < 0 || a.nodes < archs[best].nodes) {
				best = i
			}
		}
		if best < 0 {
			continue
		}
		var ds []string
		for i, a := range archs {
			if i != best {
				ds = append(ds, a.name)
			}
		}
		q := edakit.Sprintf("An agent explores %d-bit adder architectures. Its fast model treats each level as %d ps and uses: ripple-carry, %d levels and %d nodes; Brent-Kung, $2\\log_2 n - 1$ levels and $2n - 2 - \\log_2 n$ nodes, fanout 2; Sklansky, $\\log_2 n$ levels and $\\frac{n}{2}\\log_2 n$ nodes, fanout $n/2$; Kogge-Stone, $\\log_2 n$ levels and $n\\log_2 n - n + 1$ nodes, fanout 2. The budget is %d ps and the library allows a fanout of at most %d. Which architecture should it send to synthesis first, as the cheapest candidate meeting both limits?", n, d, n, n, budget, fanLimit)
		a := archs[best]
		why := edakit.Sprintf("At n = %d the candidates have (levels, nodes, fanout) of ripple (%d, %d, 2), Brent-Kung (%d, %d, 2), Sklansky (%d, %d, %d), Kogge-Stone (%d, %d, 2). The depth budget allows %d levels; the cheapest feasible one is %s with %d nodes.",
			n, n, n, archs[1].depth, archs[1].nodes, L, archs[2].nodes, n/2, L, archs[3].nodes, budget/d, a.name, a.nodes)
		return edakit.Choice(q, a.name, why, ds...)
	}
}

var specToRTLItems = []edakit.Item{
	{KP: 1, Q: "An agent has produced RTL for a FIFO. In what order should the harness run its checks?",
		A: "Parse and elaborate, lint, simulate, then formal checks and synthesis", D: []string{"Synthesis first, since area is the main goal", "Formal checks first, then parsing", "Only simulation, since it covers the others"},
		Why: "Cheap checks that catch gross errors run first, so expensive engines only see designs that already compile and lint cleanly."},
	{KP: 1, Q: "A lint run on generated RTL reports an inferred latch in a combinational always block. What feedback should the harness return to the model?",
		A: "The signal and line with the incomplete assignment, and the rule it breaks", D: []string{"A generic message that the design failed", "The full synthesis log", "Only a pass or fail bit"},
		Why: "Specific, actionable feedback lets the model fix the missing else branch or default assignment on the next iteration."},
	{KP: 1, Q: "A simulation of generated RTL fails. What is the most useful feedback for the next iteration?",
		A: "The first failing cycle, the inputs that led to it, and the expected and actual outputs", D: []string{"The total number of cycles simulated", "The waveform of every signal in the design", "A request to regenerate the whole module from scratch"},
		Why: "A minimal failing trace localizes the bug, which is what the model needs to propose a small edit."},
	{KP: 1, Q: "When should the spec-to-RTL loop hand its result to a human?",
		A: "When all checks pass, as a reviewable diff with the evidence that it passed", D: []string{"As soon as the RTL parses", "Only if the model reports high confidence", "Never; passing tests is sufficient for tape-out"},
		Why: "The loop ends with a reviewable diff plus proof, so the human approves evidence rather than trusting the model."},
	{KP: 2, Q: "An agent writes a 12-input combinational decoder. Which golden check gives complete assurance at low cost?",
		A: "Exhaustive simulation of all 4096 input combinations against a reference", D: []string{"A thousand random vectors", "Reviewing the RTL for style", "Checking that synthesis area is small"},
		Why: "Twelve inputs give 4096 rows, cheap to enumerate, so the check covers every case."},
	{KP: 2, Q: "An agent writes a 32-bit Kogge-Stone adder in gate-level RTL. Which check proves it correct?",
		A: "Combinational equivalence against `assign s = a + b`", D: []string{"Exhaustive simulation of all input pairs", "Simulation of the carry-out for a few operands", "Comparing its node count with the formula"},
		Why: "Exhaustive simulation needs $2^{64}$ vectors, but an equivalence checker proves the structure equals the language's adder for all inputs."},
	{KP: 2, Q: "An agent re-implements a sequence detector FSM with a different state encoding. What check proves the behavior is unchanged?",
		A: "Sequential equivalence from reset against the original", D: []string{"Combinational equivalence with matched registers", "Comparing the number of flip-flops", "Running the original testbench once"},
		Why: "The registers no longer correspond one to one, so a sequential check from reset is needed."},
	{KP: 2, Q: "Random simulation of a generated 16-bit signed saturating adder passes a million vectors. What should still worry you?",
		A: "Rare corner cases such as the most negative value may not have been hit", D: []string{"Nothing, a million vectors is exhaustive for 16 bits", "The adder must be wrong since it passed so quickly", "Random tests always miss overflow"},
		Why: "Two 16-bit operands give about four billion pairs; a million random ones sample a small fraction, so corners need directed tests or a formal check."},
	{KP: 3, Q: "The spec says 'detect 1011' and does not say whether overlapping matches count. What should the agent do?",
		A: "Ask the human which reading is intended and record the answer", D: []string{"Pick overlapping silently, since it is more common", "Generate both and submit whichever passes its own testbench", "Leave the ambiguity in a code comment"},
		Why: "An ambiguous spec lets the RTL and the testbench disagree; resolving it with the human prevents converging on the wrong reading."},
	{KP: 3, Q: "The model writes both the RTL and the testbench that checks it. What is the weakness?",
		A: "Both can share the same misunderstanding, so the test passes a wrong design", D: []string{"Testbenches cannot be generated by models", "The simulation will run too slowly", "Lint cannot check testbench code"},
		Why: "A golden reference must be independent of the generated design, or a shared mistake goes undetected."},
	{KP: 3, Q: "The harness loop keeps failing on the same check after five attempts, each a different rewrite. What should it do?",
		A: "Stop and escalate with the failing evidence, since the spec or the check may be wrong", D: []string{"Keep regenerating until something passes", "Disable the failing check", "Lower the simulation coverage target"},
		Why: "Thrashing usually signals an ambiguity or a broken golden reference, which needs human judgment rather than more attempts."},
	{KP: 3, Q: "A confirmed spec decision (reset value, signedness) should be kept where?",
		A: "In a recorded specification the next request starts from", D: []string{"Only in the model's conversation context", "In a comment inside the generated RTL alone", "Nowhere; each request should start fresh"},
		Why: "Recording decisions avoids re-asking and keeps later generated blocks consistent."},
}

var goldenItems = []edakit.Item{
	{KP: 1, Q: "Two netlists have the same registers, matched by name. Which check proves they compute the same functions?",
		A: "Combinational equivalence checking", D: []string{"Sequential equivalence checking from reset", "Random simulation", "A QoR comparison"},
		Why: "With corresponding registers, equivalence reduces to proving every next-state and output function equal, a combinational problem."},
	{KP: 1, Q: "What does a passing directed testbench prove about a design?",
		A: "That it behaves correctly on the vectors applied, and nothing more", D: []string{"That it is correct for all inputs", "That it matches the specification", "That it is equivalent to the previous version"},
		Why: "Simulation is evidence on the tested inputs only; completeness needs exhaustive or formal methods."},
	{KP: 1, Q: "A retimed pipeline moved registers across logic. Which check can prove it still matches the original?",
		A: "Sequential equivalence checking", D: []string{"Combinational equivalence with name matching", "Lint", "A static timing report"},
		Why: "Retiming breaks register correspondence, so equivalence must be proved over sequences from reset."},
	{KP: 1, Q: "How many vectors does exhaustive simulation of a 20-input combinational block need?",
		A: "About one million", D: []string{"About 400", "About one thousand", "About one billion"},
		Why: "$2^{20} = 1048576$ input combinations."},
	{KP: 2, Q: "An agent replaces a ripple adder with a Brent-Kung adder as a pure refactor. What must the differential check show?",
		A: "Full equivalence with the previous version, plus the intended QoR change", D: []string{"That some outputs changed, proving the edit took effect", "Only that the new design synthesizes", "That the testbench still compiles"},
		Why: "A refactor must not change behavior at all, and the QoR diff shows the intended timing and area effect."},
	{KP: 2, Q: "An agent adds a saturate mode to an adder, enabled by a new input. How is the edit checked?",
		A: "Prove equivalence with the old design when saturate is off, and check the new behavior when it is on", D: []string{"Require full equivalence with the old design", "Skip equivalence because behavior changed", "Only check the synthesized area"},
		Why: "Constraining equivalence to the unchanged mode proves nothing else moved, while the new mode is checked against the spec."},
	{KP: 2, Q: "After an FSM re-encoding, a differential QoR report shows area down 8% and worst slack unchanged. What else must hold before accepting it?",
		A: "Sequential equivalence from reset with the original", D: []string{"Nothing; the QoR improved", "A larger flip-flop count", "An identical netlist"},
		Why: "QoR shows the physical effect, but correctness needs its own proof."},
	{KP: 2, Q: "A change was meant to affect only the multiplier, yet the differential QoR shows the control FSM's area changed too. What is the right reaction?",
		A: "Investigate: the edit had an unintended effect outside its target", D: []string{"Accept it, since total area went down", "Ignore it as synthesis noise in all cases", "Revert the multiplier change only"},
		Why: "Differential checks exist to show that only the intended part moved; an unexpected change needs an explanation."},
	{KP: 3, Q: "Two FSM versions are sequentially equivalent from reset. The spec requires safe recovery from any illegal state. Is the equivalence proof enough?",
		A: "No: they may still differ when started from an illegal code", D: []string{"Yes, equivalence from reset covers all states", "Yes, illegal states cannot occur", "No, equivalence checking never applies to FSMs"},
		Why: "Reset-based equivalence ignores unreachable codes, which is exactly where the recovery requirement lives."},
	{KP: 3, Q: "Both the old and new versions of a block were derived from the same flawed reference model. What will the differential check show?",
		A: "A pass, because both share the same bug", D: []string{"A failure at the bug's location", "A QoR regression", "A lint error"},
		Why: "Differential checks prove sameness, not correctness; an independent golden reference is needed to catch a shared bug."},
	{KP: 3, Q: "Synthesis is free to choose values for X assignments. Why can that make an equivalence pass misleading?",
		A: "The two versions may resolve X differently and simulate differently, while the checker treats X as a free choice", D: []string{"X values are always resolved to 0 by every tool", "X assignments are illegal in synthesizable RTL", "Equivalence checkers reject any design containing X"},
		Why: "The X semantics of the check must be stated, or a pass may hide behavior differences that simulation would show."},
	{KP: 3, Q: "What should a trustworthy harness report alongside every equivalence pass?",
		A: "Its reset assumptions and its treatment of X and don't-cares", D: []string{"Only the runtime of the check", "The model's confidence score", "The number of lines changed"},
		Why: "Those assumptions define what the proof covers; without them a reviewer cannot judge its strength."},
}

var dseItems = []edakit.Item{
	{KP: 1, Q: "Why are closed-form depth and node counts useful to an architecture-exploration agent?",
		A: "They prune infeasible architectures in microseconds before any synthesis run", D: []string{"They replace synthesis entirely", "They include wire delay exactly", "They guarantee the post-layout winner"},
		Why: "Closed forms are a fast model: cheap enough to evaluate every candidate, but blind to wires and placement."},
	{KP: 1, Q: "On closed forms alone, Sklansky has the same depth as Kogge-Stone with fewer nodes. Why might an agent still pick Kogge-Stone?",
		A: "Sklansky's fanout of n/2 at the last level needs buffering that adds delay", D: []string{"Kogge-Stone has fewer wires", "Sklansky cannot be built for powers of two", "Kogge-Stone has fewer levels"},
		Why: "The node count hides Sklansky's high fanout, which costs buffers and delay once sized."},
	{KP: 1, Q: "What does the fast closed-form model of a prefix adder ignore that most often changes the ranking?",
		A: "Wire length and capacitance from placement and routing", D: []string{"The number of logic levels", "The number of prefix nodes", "The operand width"},
		Why: "Kogge-Stone's long wires and Brent-Kung's short ones only show up after placement."},
	{KP: 2, Q: "Five adder architectures pass the fast filter, but synthesis runs are expensive. Which should the agent synthesize?",
		A: "The few candidates within the fast model's error margin of the budget or of each other", D: []string{"Only the single cheapest one, since the fast model is exact", "All of them at every corner, regardless of cost", "None, since the fast model already decided"},
		Why: "Accurate runs are spent where the fast model cannot discriminate."},
	{KP: 2, Q: "The agent records fast estimates beside post-synthesis results and finds Kogge-Stone is consistently 15% slower than predicted. What should it do?",
		A: "Calibrate the fast model with a correction for that architecture", D: []string{"Discard the synthesis results as noise", "Stop using the fast model altogether", "Report the fast estimate as the result"},
		Why: "A systematic gap is knowledge: it improves the fast model for the next exploration."},
	{KP: 2, Q: "A post-synthesis result for a candidate is within 2 ps of the budget. What is the right next step?",
		A: "Recurse to placement and routing before trusting it, since wires can push it over", D: []string{"Accept it, since synthesis is signoff", "Reject it immediately", "Rerun the closed-form model"},
		Why: "Near the budget, only the more accurate engine can settle feasibility."},
	{KP: 2, Q: "Which split of work matches the fast-and-accurate loop for arithmetic exploration?",
		A: "Closed forms and quick synthesis rank candidates; placement and timing decide among the finalists", D: []string{"Placement and timing for every candidate, closed forms for the finalist", "Closed forms only, with no accurate runs", "Accurate runs only, with no fast model"},
		Why: "Cheap models cover the space broadly; expensive engines are reserved for the frontier."},
	{KP: 3, Q: "An agent minimizes area plus 0.1 times delay and picks a design that misses the clock target. What went wrong?",
		A: "A hard timing constraint was folded into a weighted sum that let area outweigh it", D: []string{"The delay weight was too large", "Area should never be optimized", "The agent should have used ripple-carry"},
		Why: "Hard constraints must be constraints; a scalar weight can trade them away."},
	{KP: 3, Q: "Why keep a Pareto set of (delay, area, power) rather than a single best design?",
		A: "Different contexts need different trade-offs, and dominated points can be dropped safely", D: []string{"A Pareto set always contains exactly one design", "It avoids running synthesis", "It guarantees the fastest design is also smallest"},
		Why: "The Pareto set keeps every non-dominated option, so a later constraint picks from it without re-exploring."},
	{KP: 3, Q: "The agent writes a Brent-Kung structure in RTL, but the synthesized netlist looks like a Sklansky adder. What happened?",
		A: "Synthesis re-optimized the arithmetic, so the agent must check what was built", D: []string{"The RTL was ignored entirely", "Brent-Kung and Sklansky are the same network", "The netlist is mislabeled"},
		Why: "Datapath synthesis may restructure explicit adders, so credit must go to the implemented architecture."},
	{KP: 3, Q: "The best adder at 16 bits is reused without re-exploration for a 64-bit datapath. Why is that risky?",
		A: "Depth, node and wire costs scale differently, so the ranking can change with width", D: []string{"64-bit adders cannot use prefix networks", "The 16-bit design will not synthesize at 64 bits", "Width does not affect the choice at all"},
		Why: "Ripple delay grows linearly while prefix depth grows logarithmically, and wire costs grow with width, so winners change."},
	{KP: 3, Q: "A tool-grounded reward for adder exploration uses logic depth only. Which design does it over-reward?",
		A: "A minimum-depth, wire-heavy network that loses after placement", D: []string{"A ripple-carry adder", "A Brent-Kung adder at a loose target", "Any design that fails equivalence"},
		Why: "A depth-only proxy ignores wires and fanout, so it prefers networks whose real delay is worse."},
}

var tierItems = []edakit.Item{
	{KP: 1, Q: "A customer wants the FIFO generator's depth changed from 16 to 64. Which tier does the change land in?",
		A: "Tier 1: a parameter change, re-run with golden tests", D: []string{"Tier 2: a new plug-in", "Tier 3: a core edit", "None: it needs a new generator"},
		Why: "Depth is a parameter the generator exposes; no code changes."},
	{KP: 1, Q: "A team asks for one more pipeline stage in a generated multiplier that already exposes a stage count. Which tier?",
		A: "Tier 1, configuration", D: []string{"Tier 2, plug-in", "Tier 3, core", "It cannot be done with a generator"},
		Why: "When the stage count is a parameter, adding a stage is a configuration change."},
	{KP: 1, Q: "An agent edits the generator's source to change a data width from 32 to 48. What is wrong?",
		A: "It is working at the wrong tier: width should be a validated parameter", D: []string{"Nothing, source edits are always preferred", "Widths cannot be changed at all", "It should edit the netlist instead"},
		Why: "Parameter changes belong at tier 1 so they need no code change and no recompile."},
	{KP: 1, Q: "Why should a generator validate parameter values before elaboration?",
		A: "So a configuration change cannot produce an illegal design", D: []string{"So it runs faster", "So parameters can be ignored", "So the core never needs testing"},
		Why: "Validation keeps tier 1 safe: a bad value is rejected rather than turned into broken RTL."},
	{KP: 2, Q: "A team wants a Han-Carlson prefix network added to an adder generator that already offers Kogge-Stone and Brent-Kung. Which tier?",
		A: "Tier 2: a new policy behind the prefix-network seam", D: []string{"Tier 1: a parameter change", "Tier 3: a core edit to the generate-propagate logic", "None: it needs a new tool"},
		Why: "The new network implements the existing interface; the core's g, p and sum logic is untouched."},
	{KP: 2, Q: "What should the seam for adder architectures in an arithmetic generator take and return?",
		A: "Take the width and the g, p signals; return a prefix network producing the group generates", D: []string{"Take the whole RTL module and return a netlist", "Take the clock period and return a gate count", "Take the sum bits and return the carries"},
		Why: "A narrow interface keeps the shared logic in the core and lets each policy supply only the prefix network."},
	{KP: 2, Q: "A new FSM encoding policy is added as a plug-in to a controller generator. What proves it safe to promote?",
		A: "Sequential equivalence with the existing encodings on the golden FSMs, plus a QoR comparison", D: []string{"A successful compile", "A smaller flip-flop count only", "Model confidence"},
		Why: "The differential proof shows behavior is unchanged; QoR shows where the policy sits."},
	{KP: 2, Q: "A plug-in replacement policy for a cache generator passes its tests. What must remain untouched?",
		A: "The generator core, including tag comparison and the interface to the policy", D: []string{"The cache's data array size parameter", "The golden tests", "The QoR baseline"},
		Why: "Tier 2 changes live behind the seam; the core and its contracts stay fixed."},
	{KP: 3, Q: "A request says: make every generated flip-flop use an asynchronous reset. Which tier?",
		A: "Tier 3: it changes reset semantics in every generated design", D: []string{"Tier 1: a single option", "Tier 2: a plug-in for one module", "None: reset is outside the generator"},
		Why: "The reach of the change, not its size, sets the tier; reset semantics affect timing, CDC and every user."},
	{KP: 3, Q: "What does a tier 3 change to a generator require before it is merged?",
		A: "A human-approved diff and full regression across configurations", D: []string{"Only a passing compile", "A single golden test", "Nothing beyond the agent's review"},
		Why: "Core changes can affect every design the generator produces."},
	{KP: 3, Q: "A change renames how the generator emits module ports. Why is that core?",
		A: "Every downstream script and integration that references the ports can break", D: []string{"Port names have no effect on anything", "It only affects simulation speed", "It is a parameter by definition"},
		Why: "Naming and connectivity conventions are invariants that other tools depend on."},
	{KP: 3, Q: "Requests for one kind of change keep landing in tier 3. What does that indicate?",
		A: "A missing seam: the variable behavior should be moved behind an extension point", D: []string{"The core is well designed", "Agents should not handle such requests", "The tier model does not apply to generators"},
		Why: "A seam is in the right place when that kind of change essentially never needs a core edit."},
}

var rewardItems = []edakit.Item{
	{KP: 1, Q: "Why should a QoR term in an RTL reward be granted only to designs that already pass functional checks?",
		A: "Otherwise an empty or wrong design wins on area", D: []string{"QoR cannot be measured on failing designs", "QoR is always zero for correct designs", "Functional checks are too expensive"},
		Why: "A design that does nothing is tiny; gating prevents rewarding it."},
	{KP: 1, Q: "Which tool signal gives the strongest functional reward when a golden RTL reference exists?",
		A: "A formal equivalence result against the reference", D: []string{"Lint warnings count", "Parse success", "Synthesized area"},
		Why: "Equivalence covers all inputs, unlike sampled simulation."},
	{KP: 1, Q: "What should a generated design that fails to elaborate receive?",
		A: "Zero reward regardless of other signals", D: []string{"Partial credit for its area", "Full credit if it lints cleanly", "Credit for passing simulation"},
		Why: "Nothing downstream is meaningful for a design the tools cannot elaborate."},
	{KP: 1, Q: "No golden RTL exists for a new block. What is the next-best functional signal?",
		A: "A software reference model with exhaustive or high-coverage simulation", D: []string{"The model's self-reported confidence", "A synthesis area target", "Lint cleanliness"},
		Why: "An independent executable reference still checks behavior directly."},
	{KP: 2, Q: "A model's design passes the testbench by matching exactly the eight vectors it applies, using a lookup table. What reward hack is this?",
		A: "Hard-coding the test vectors", D: []string{"Exploiting a fast QoR proxy", "Removing the reset", "Unchecked outputs"},
		Why: "The design fits the test rather than the specification; hidden or randomized tests and formal checks defeat it."},
	{KP: 2, Q: "A checker ignores outputs that are X. A model learns to leave a hard output undriven. What is the defense?",
		A: "X-aware comparison that fails on any unknown output", D: []string{"More lint rules", "A larger area penalty", "Shorter simulations"},
		Why: "The exploit relies on the checker accepting X; strict comparison closes it."},
	{KP: 2, Q: "The testbench never asserts reset, and a model's designs omit reset logic to save area. How should the reward change?",
		A: "Apply reset explicitly in every test and check post-reset values", D: []string{"Raise the area weight", "Remove the area term", "Accept designs without reset"},
		Why: "The reward must exercise the behavior the specification requires, including initialization."},
	{KP: 2, Q: "Which defense addresses rewards being gamed through a quick, unplaced synthesis proxy?",
		A: "Measure QoR with the same flow used for signoff, at least periodically", D: []string{"Drop QoR from the reward entirely", "Only reward parse success", "Use a smaller testbench"},
		Why: "Scoring with the accurate flow reveals designs that look good only on the proxy."},
	{KP: 3, Q: "A model trained on a logic-depth proxy starts producing high-fanout, wire-heavy circuits. What does this show?",
		A: "The proxy's ranking disagrees with post-layout results, and the model learned the proxy", D: []string{"The model has found the true optimum", "Fanout does not matter after placement", "Depth is the signoff metric"},
		Why: "Optimizing a proxy exploits its blind spots; here, wires and fanout."},
	{KP: 3, Q: "How should a training loop monitor the gap between a fast proxy reward and accurate QoR?",
		A: "Score a sample with the accurate flow and track rank correlation with the proxy", D: []string{"Assume the gap is constant", "Compare only the best design", "Stop measuring once training starts"},
		Why: "Falling rank agreement signals that the proxy needs recalibration or reweighting."},
	{KP: 3, Q: "Which result should be reported as the final evaluation of an RTL-generating model?",
		A: "Accurate-flow QoR and correctness on held-out specifications", D: []string{"Proxy QoR on the training specifications", "Average reward during training", "The number of designs generated"},
		Why: "Held-out specs and the accurate flow measure what the model actually achieves, not what it optimized."},
	{KP: 3, Q: "Rank correlation between proxy and accurate scores drops sharply late in training. What is the right response?",
		A: "Recalibrate or reweight the proxy before continuing", D: []string{"Increase the learning rate", "Report the proxy results anyway", "Remove the accurate checks to save time"},
		Why: "When the proxy no longer tracks the target, further optimization against it moves away from real QoR."},
}
