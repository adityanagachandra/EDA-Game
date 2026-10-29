package synth

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var recipeNames = []string{
	"b; rw; rf; b; rw; rwz; b; rfz; rwz; b",
	"b; rw; rwz; b; rwz; b",
	"rw; rf; rs; b",
	"b -l; rw -l; rf -l; b -l",
	"fraig; rw; rf; b",
	"rs; rw; rs; rf; b",
}

// recipeReward: which recipe maximizes a weighted area-delay reward.
func recipeReward(ctx generator.GeneratorContext) generator.Problem {
	for {
		a0, l0 := edakit.Int(800, 2000), edakit.Int(20, 40)
		w := edakit.Pick([]float64{0.3, 0.5, 0.7})
		idx := edakit.Shuffle(seq(len(recipeNames)))[:4]
		type res struct {
			name   string
			a, l   int
			reward float64
		}
		var rs []res
		for _, i := range idx {
			a := a0 - edakit.Int(a0/20, a0/4)
			l := l0 - edakit.Int(0, l0/3)
			r := w*float64(a0-a)/float64(a0) + (1-w)*float64(l0-l)/float64(l0)
			rs = append(rs, res{recipeNames[i], a, l, r})
		}
		sort.Slice(rs, func(i, j int) bool { return rs[i].reward > rs[j].reward })
		if rs[0].reward-rs[1].reward < 0.01 {
			continue
		}
		var parts []string
		for _, r := range edakit.Shuffle(rs) {
			parts = append(parts, edakit.Sprintf("[%s] gives %d AND nodes and %d levels", r.name, r.a, r.l))
		}
		var ds []string
		for _, r := range rs[1:] {
			ds = append(ds, r.name)
		}
		q := edakit.Sprintf("An agent learns ABC recipes with the reward $R = %.1f \\cdot \\frac{A_0 - A}{A_0} + %.1f \\cdot \\frac{L_0 - L}{L_0}$, where the baseline AIG has $A_0 = %d$ nodes and $L_0 = %d$ levels. On this design, %s. Which recipe earns the highest reward?",
			w, 1-w, a0, l0, strings.Join(parts, "; "))
		return edakit.Choice(q, rs[0].name,
			edakit.Sprintf("Rewards: best %.3f for [%s], next %.3f. The weighting decides whether area or depth wins; neither is mapped QoR.", rs[0].reward, rs[0].name, rs[1].reward), ds...)
	}
}

// qorRank: discordant pairs between a fast proxy and the accurate metric.
func qorRank(ctx generator.GeneratorContext) generator.Problem {
	n := 4
	proxy := edakit.Shuffle([]int{edakit.Int(900, 950), edakit.Int(960, 1010), edakit.Int(1020, 1070), edakit.Int(1080, 1130)})
	acc := edakit.Shuffle([]int{edakit.Int(400, 420), edakit.Int(425, 445), edakit.Int(450, 470), edakit.Int(475, 495)})
	disc := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if (proxy[i]-proxy[j])*(acc[i]-acc[j]) < 0 {
				disc++
			}
		}
	}
	var parts []string
	for i := 0; i < n; i++ {
		parts = append(parts, edakit.Sprintf("recipe %d: %d AIG nodes, %d square microns of mapped area", i+1, proxy[i], acc[i]))
	}
	q := "A synthesis loop scores recipes with a fast proxy (AIG node count) and occasionally runs mapping for the accurate metric (area). Results: " + strings.Join(parts, "; ") + ". Of the 6 pairs of recipes, how many does the proxy rank in the opposite order from the accurate metric?"
	return edakit.Choice(q, edakit.IntStr(disc),
		edakit.Sprintf("Compare every pair: %d pairs are discordant. A loop trusting the proxy alone would pick a worse recipe whenever the best pair is discordant.", disc),
		nearInts(disc, 0, 6, 3)...)
}

var recipeRLItems = []edakit.Item{
	{KP: 2, Q: "An RL agent picks ABC commands one at a time. What is the most natural state for its policy?",
		A: "Features of the current AIG, such as node count, level count and recent command history", D: []string{"The raw Verilog text of the design", "The final mapped netlist, available only at the end", "A random seed, since recipes do not depend on the design"},
		Why: "The agent acts on the evolving network, so cheap structural statistics plus history form the observation."},
	{KP: 2, Q: "The reward is the AIG node reduction after each step. The agent learns to end every episode with many rwz commands. What is the concern?",
		A: "It may be exploiting the proxy: zero-gain moves cost run time and need not improve mapped QoR", D: []string{"rwz is not a valid ABC command", "rwz always increases the node count", "The agent has found the global optimum"},
		Why: "Reward on a fast proxy can be gamed; the real target is mapped area and delay after timing."},
	{KP: 2, Q: "Why is a learned recipe policy validated on held-out designs?",
		A: "To check that it generalizes rather than memorizing what works on the training circuits", D: []string{"Because ABC runs faster on new designs", "Because held-out designs have no AND nodes", "Because training designs cannot be mapped"},
		Why: "Recipe quality is design-dependent; overfitting to a benchmark suite is a common failure."},
	{KP: 2, Q: "An agent's recipe reduces AIG nodes by 12% on average over resyn2. What else must be reported before adopting it?",
		A: "Mapped area and delay after timing, run time, and equivalence of every result", D: []string{"Only the number of commands in the recipe", "Only the node reduction on the best design", "Nothing else, node count is the final metric"},
		Why: "Node count is an intermediate proxy; adoption requires signoff-relevant QoR, cost and correctness."},
	{KP: 3, Q: "In the four-layer model, where does a learned recipe policy plug into the synthesis engine?",
		A: "At the strategy layer: it chooses the sequence of compiled transformations, while the AIG core is untouched", D: []string{"In the immutable core, replacing structural hashing", "In the observation hooks, as a read-only probe", "In the Liberty library, as a new cell"},
		Why: "Recipes are policy; the transformations and the AIG data structure stay in the core."},
	{KP: 3, Q: "A tool-grounded reward for synthesis should be computed by:",
		A: "Running the actual mapper and timer on the result and checking equivalence", D: []string{"Asking a language model to rate the recipe", "Counting the characters of the script", "Comparing the recipe with resyn2 by edit distance"},
		Why: "Deterministic tools give a reward that reflects real QoR and correctness."},
	{KP: 3, Q: "The agent's policy outputs one script per design, chosen before running anything. Why might a closed-loop policy do better?",
		A: "It observes the AIG after each command and adapts, instead of committing to a fixed sequence", D: []string{"Closed-loop policies never need a reward", "Fixed scripts cannot contain balance", "ABC runs commands in random order"},
		Why: "Feedback lets the agent react to the actual effect of each transformation on this design."},
	{KP: 3, Q: "Two recipes tie on the training reward. One takes 3 seconds and the other 40. What should the reward design include?",
		A: "A run-time penalty, since the loop and the user pay for every command", D: []string{"Nothing, run time is irrelevant to synthesis", "A bonus for longer scripts", "A random tie-breaker only"},
		Why: "Practical recipe search optimizes QoR per unit of compute, not QoR alone."},
}

var qorItems = []edakit.Item{
	{KP: 2, Q: "A design-space exploration agent evaluates 500 recipes. Which plan uses fast and accurate engines well?",
		A: "Score all 500 with AIG metrics or a fast mapper, then run full mapping and STA on the top few", D: []string{"Run full place and route on all 500", "Use only AIG node count and never map", "Pick a recipe at random, then run STA"},
		Why: "Fast proxies prune the space; accurate engines confirm the shortlist."},
	{KP: 2, Q: "When must the loop recurse from the fast proxy to accurate mapping?",
		A: "When the candidates are within the proxy's error, or when the proxy and accurate ranks have disagreed on similar designs", D: []string{"Never, the proxy is always right", "Only after tapeout", "Only when the AIG has more than one output"},
		Why: "Escalate when the decision depends on differences the proxy cannot resolve."},
	{KP: 2, Q: "Which fast metric best tracks delay before mapping?",
		A: "AIG level count, or the depth of a quick LUT or cell mapping", D: []string{"Number of primary inputs", "AIG node count alone", "The number of commands in the script"},
		Why: "Level count measures logic depth; node count measures area."},
	{KP: 2, Q: "A fast QoR estimator is a learned model trained on mapped results. What is its main risk in an optimization loop?",
		A: "The optimizer drives designs outside the training distribution, where the model's predictions are unreliable", D: []string{"It is always slower than real mapping", "It cannot predict area at all", "It changes the circuit's function"},
		Why: "Optimizing against a surrogate exploits its errors; periodic accurate checks keep it honest."},
	{KP: 3, Q: "The proxy and accurate metric agree on 90% of recipe pairs. Is the proxy safe to choose the final recipe?",
		A: "Not by itself: the final choice is among the best candidates, exactly where small disagreements decide the winner", D: []string{"Yes, 90% agreement guarantees the best recipe", "No, a proxy must agree on 100% of pairs to be useful", "Yes, as long as node count decreased"},
		Why: "Rank correlation over all pairs can be high while the top-k ordering is still wrong."},
	{KP: 3, Q: "Which observation hook makes a fast-accurate loop possible in a synthesis engine?",
		A: "Reporting node count, levels and a quick mapped estimate after every pass", D: []string{"Writing the final netlist only at the end", "Hiding intermediate AIGs from the agent", "Logging only run time"},
		Why: "The agent decides the next skill from live metrics, which the engine must expose cheaply."},
	{KP: 3, Q: "An agent tunes recipes on a fast mapper with a unit-delay model. After real STA, the chosen recipe is slower. What is the most likely cause?",
		A: "Unit delay ignores fanout loads and wire delay, so the proxy rewarded structures that are slow in reality", D: []string{"STA is nondeterministic", "The recipe changed the function", "The fast mapper used the Liberty library exactly"},
		Why: "Proxy-reality gaps come from what the fast model leaves out."},
	{KP: 3, Q: "How should the variance of a fast metric across seeds influence the loop?",
		A: "High variance means differences below the noise floor should not drive decisions without accurate confirmation", D: []string{"Variance is irrelevant if the mean is good", "High variance means the proxy is accurate", "The loop should pick the recipe with the highest variance"},
		Why: "Decisions are only as reliable as the metric's resolution."},
}

var egraphSkillItems = []edakit.Item{
	{KP: 1, Q: "An e-graph rewriting engine is packaged as a skill. What belongs in its manifest?",
		A: "The rule sets it supports, the cost models for extraction, its input and output formats, and its tunable budgets", D: []string{"The full source of the e-graph library", "The customer's netlist", "A list of every e-class from the last run"},
		Why: "A manifest describes what the skill does and what can be tuned, so an agent can select it without reading code."},
	{KP: 1, Q: "A team adds a new rewrite rule for their custom multiply-accumulate. Which change tier is that?",
		A: "Tier 2: a plug-in rule set at a defined extension point, with the e-graph core untouched", D: []string{"Tier 1: a parameter change", "Tier 3: an edit to union-find and rebuild", "Not a change, rules are part of the design"},
		Why: "Rules are policy that the core loads; the saturation and congruence machinery stays fixed."},
	{KP: 1, Q: "Which part of an e-graph engine is the immutable core?",
		A: "Hashconsing, union-find, congruence rebuild and the extraction algorithm interface", D: []string{"The particular rule set for adders", "The node budget", "The cost assigned to a shift"},
		Why: "These guarantee that every represented term is equivalent; rule sets and costs vary per user."},
	{KP: 1, Q: "Which is a Tier 1 change to an e-graph skill?",
		A: "Raising the node budget from 100,000 to 400,000 e-nodes", D: []string{"Adding a rule for a new operator", "Replacing union-find with a different structure", "Writing a new extraction algorithm"},
		Why: "Budgets and cost weights are configuration and need no recompile."},
	{KP: 2, Q: "What must the golden test of a new rewrite rule check?",
		A: "That the rule is sound (both sides equal at every bit width it matches) and that designs extracted with it pass equivalence checking", D: []string{"That it fires at least once", "That it reduces cost on every design", "That it keeps the e-graph small"},
		Why: "Soundness is the invariant; one unsound rule silently corrupts every term it touches."},
	{KP: 2, Q: "An agent proposes the rule $x \\cdot 2^k \\to x \\ll k$ for signed operands. What should the harness require before promoting it?",
		A: "A proof or exhaustive check at small widths showing the rule is sound for the stated widths and signedness, plus a differential QoR report", D: []string{"A single example where it reduces area", "Only that it compiles", "Approval because the rule looks standard"},
		Why: "Width and signedness conditions are where datapath rules go wrong; promotion needs proof and evidence."},
	{KP: 2, Q: "Why should rule sets be versioned separately from the e-graph engine?",
		A: "So a QoR change can be traced to a rule change and reverted without touching the engine", D: []string{"Because rules cannot be compiled", "Because the engine changes every run", "Because versioning makes saturation faster"},
		Why: "Separate versions make differential testing and rollback precise."},
	{KP: 2, Q: "A differential test runs saturation with and without a new rule on 50 designs. What does it measure?",
		A: "Which designs' extracted cost changed and by how much, with equivalence confirmed for each", D: []string{"Only whether the e-graph got larger", "Only run time", "Only that no design failed to parse"},
		Why: "Differential tests show that only intended designs moved, and only in the intended direction."},
	{KP: 3, Q: "Equality saturation runs out of budget before a beneficial rewrite fires. What is a plug-in-level remedy?",
		A: "A scheduler policy that limits explosive rules such as associativity, so useful rules get budget", D: []string{"Removing congruence closure", "Disabling extraction", "Doubling every operator cost"},
		Why: "Rule scheduling is policy at an extension point; egg-style backoff schedulers do exactly this."},
	{KP: 3, Q: "Where should an agent's datapath e-graph skill hand results to the rest of the flow?",
		A: "As word-level RTLIL or an equivalent IR before bit-blasting, followed by equivalence checking", D: []string{"As a mapped netlist that skips synthesis", "As a Liberty file", "As placement constraints"},
		Why: "Word-level rewrites belong before lowering; downstream synthesis then optimizes the bits."},
	{KP: 3, Q: "A customer wants e-graph extraction to minimize their own power model. What is the right seam?",
		A: "A cost-function plug-in called by the extractor, given each e-node's operator, widths and context", D: []string{"Editing the union-find code", "Changing the rule set", "Hard-coding the costs in the core"},
		Why: "Extraction asks one question at one seam, the cost of an e-node; the model behind it varies per customer."},
	{KP: 3, Q: "The skill flywheel promotes a rule set found for one design into the library. What must stay behind?",
		A: "The customer's design data; only the validated rules and when to use them move into the shared library", D: []string{"The rules themselves", "The equivalence proofs", "The manifest"},
		Why: "Know-how flows back to the library; design data stays inside the IP boundary."},
}

var cecGoldenItems = []edakit.Item{
	{KP: 1, Q: "An agent proposes a new AIG transformation. What is the non-negotiable golden test?",
		A: "Combinational equivalence of the output against the input on every benchmark, proved with a miter and SAT", D: []string{"A node count reduction on most benchmarks", "Successful compilation", "Random simulation of 100 vectors"},
		Why: "QoR gains are worthless if the function changes; CEC proves it did not."},
	{KP: 1, Q: "Why is random simulation not enough as the golden test?",
		A: "It can miss rare differences; only a SAT or BDD proof covers every input vector", D: []string{"Simulation is slower than SAT on every circuit", "Simulation cannot handle AND gates", "Simulation changes the circuit"},
		Why: "Simulation finds counterexamples quickly but cannot prove their absence."},
	{KP: 1, Q: "A miter is built for two versions of a circuit. When are they equivalent?",
		A: "When the miter output is unsatisfiable, that is, constant 0", D: []string{"When the miter output is satisfiable", "When the miter has fewer nodes than either circuit", "When the miter output is 1 on half the inputs"},
		Why: "The miter XORs corresponding outputs; any satisfying assignment is a counterexample."},
	{KP: 1, Q: "Which transformation cannot be verified by combinational equivalence checking alone?",
		A: "Retiming, because registers move and lose their one-to-one correspondence", D: []string{"AIG rewriting", "Balancing", "Technology mapping of combinational logic"},
		Why: "Sequential transformations need sequential equivalence checking."},
	{KP: 2, Q: "An agent's transformation passes CEC on 99 benchmarks and fails on one. What should the harness do?",
		A: "Reject it, feed the counterexample back to the model, and keep the failing design as a new golden test", D: []string{"Accept it, since 99% is good enough", "Remove the failing benchmark from the suite", "Accept it with a warning"},
		Why: "One counterexample proves a bug; the failing case becomes part of the regression."},
	{KP: 2, Q: "CEC times out on a large design after an agent's change. What is a sound response?",
		A: "Treat the result as unproven and do not promote the change until it is proved, for example with partitioning or FRAIG-based sweeping", D: []string{"Treat a timeout as a pass", "Shorten the timeout until it passes", "Skip verification for large designs"},
		Why: "Absence of a counterexample within a time limit is not a proof."},
	{KP: 2, Q: "Why does checking each transformation step individually help an agentic loop?",
		A: "A failure pinpoints the step that introduced it, and small miters prove much faster", D: []string{"Individual checks are required by SAT solvers", "It avoids the need for a final check", "It makes the transformations faster"},
		Why: "Local checks give precise feedback; a final end-to-end check still guards composition."},
	{KP: 2, Q: "Internal equivalences between the original and optimized AIG speed up CEC. Why?",
		A: "Proving matched internal nodes first cuts the miter into small, easy SAT problems", D: []string{"They make the miter output constant 1", "They remove the need for SAT entirely", "They let CEC ignore primary outputs"},
		Why: "SAT sweeping merges proven internal equivalences, so the output check becomes trivial."},
	{KP: 3, Q: "Beyond CEC, what does a differential QoR check add for an agent-proposed pass?",
		A: "Evidence that the pass moves area and delay in the intended direction on the intended designs and nowhere else", D: []string{"A proof of equivalence", "A faster SAT solver", "A new Liberty library"},
		Why: "CEC proves correctness; differential QoR proves usefulness and absence of regressions."},
	{KP: 3, Q: "An agent proposes using don't cares from the testbench environment in synthesis. What must change in the golden test?",
		A: "Equivalence must be checked modulo those don't cares, with the environment constraints made explicit", D: []string{"Nothing, plain CEC still passes", "CEC must be disabled", "Only simulation is needed"},
		Why: "Using external don't cares legitimately changes outputs on excluded inputs, so the checker needs the same constraints."},
	{KP: 3, Q: "Why should the equivalence checker belong to a different code path than the optimizer it checks?",
		A: "A shared bug in a common component could make both the optimizer and the checker wrong in the same way", D: []string{"Because shared code cannot compile", "Because checkers must be slower", "Because optimizers cannot use SAT"},
		Why: "Independent verification avoids common-mode failures."},
	{KP: 3, Q: "The harness accepts a transformation only if CEC passes and QoR improves. What residual risk remains?",
		A: "Run-time blowups or nondeterminism on unseen designs, which need resource limits and determinism checks", D: []string{"The function could still change", "QoR could get worse on the benchmarks", "The miter could be satisfiable"},
		Why: "Correctness and QoR are covered; robustness and reproducibility need their own tests."},
}

var changeTierItems = []edakit.Item{
	{KP: 1, Q: "A customer asks the synthesis engine to use 6-input cuts instead of 4-input cuts in rewriting. Which tier?",
		A: "Tier 1: a parameter of the rewriting skill", D: []string{"Tier 2: a new rewriting plug-in", "Tier 3: a change to the AIG core", "Not an engine change at all"},
		Why: "Cut size is a parameter, as long as the core supports it; no new code is needed."},
	{KP: 1, Q: "A customer asks for a recipe that runs balance only on timing-critical cones. Which tier?",
		A: "Tier 2: a new strategy plug-in that selects cones using a timing hook and calls the existing balance pass", D: []string{"Tier 1: a flag of balance", "Tier 3: rewriting the AIG data structure", "Tier 1: changing the Liberty library"},
		Why: "It composes existing transformations under a new policy at an extension point."},
	{KP: 1, Q: "A customer wants AIG nodes to carry a new attribute that every pass must preserve. Which tier?",
		A: "Tier 3: it changes the core data structure and every transformation's contract", D: []string{"Tier 1: a parameter", "Tier 2: a plug-in pass", "Not an engine change"},
		Why: "Changing the core representation touches every pass; it needs a human-approved diff and full regression."},
	{KP: 1, Q: "In the four-layer model for a synthesis engine, what is the immutable core?",
		A: "The AIG data structure with structural hashing and the equivalence-preserving transformation contract", D: []string{"The default resyn2 script", "The cut size and iteration counts", "The list of designs in the benchmark suite"},
		Why: "Scripts are policy and numbers are parameters; the representation and its invariants are the core."},
	{KP: 2, Q: "Where is the right seam for customer-specific cost in rewriting (for example, penalizing XOR-heavy structures)?",
		A: "A cost provider that the rewriter calls for each candidate replacement", D: []string{"Inside structural hashing", "In the Verilog parser", "In the netlist writer"},
		Why: "The rewriter asks one question, is this replacement better, at one seam; the answer varies by customer."},
	{KP: 2, Q: "What context should a rewriting cost seam expose?",
		A: "The cut, the MFFC size, the candidate's node count and levels, and the node's timing criticality", D: []string{"Only the node's name", "The full RTL source", "The customer's revenue target"},
		Why: "A cost decision needs local structure and live observations, nothing else."},
	{KP: 2, Q: "Why should a cost seam be called per candidate batch rather than per AND node visit?",
		A: "Coarse calls let the core keep its hot loops inlined and fast", D: []string{"Per-node calls are not deterministic", "Batches change the function", "The core cannot call plug-ins"},
		Why: "Seam granularity matters for performance; fine-grained virtual calls in the hot loop are expensive."},
	{KP: 2, Q: "A customer wants their own FSM encoding heuristic. Where does it belong?",
		A: "A Tier 2 encoder plug-in at the FSM pass's encoding extension point", D: []string{"A Tier 3 edit to the AIG", "A Tier 1 flag of ABC's strash", "A change to the Liberty library"},
		Why: "Encoding is a policy choice inside a pass with a stable interface."},
	{KP: 3, Q: "A plug-in rewriting cost passes golden CEC tests. What else must the differential test show?",
		A: "That QoR changed only on designs where the new cost should matter, by the intended amount", D: []string{"That the core code changed", "That run time doubled", "That every design got smaller"},
		Why: "Differential evidence shows the plug-in did what it claimed and nothing else."},
	{KP: 3, Q: "The harness reports that a requested change needs a Tier 3 edit every time a customer asks for a new cost. What does that say?",
		A: "The seam is in the wrong place: cost should be an extension point, not part of the core", D: []string{"The customers are asking for too much", "Tier 3 edits are the normal path", "The golden tests are too strict"},
		Why: "A well-placed seam makes common requests Tier 1 or Tier 2."},
	{KP: 3, Q: "An agent proposes speeding up structural hashing by changing the hash table. Which process applies?",
		A: "Tier 3: a human-reviewed diff plus full regression, since every pass depends on it", D: []string{"Tier 1: apply it as a configuration change", "Tier 2: drop it in as a plug-in", "No review, performance changes are safe"},
		Why: "The core's correctness is everyone's correctness; even pure speedups need review and regression."},
	{KP: 3, Q: "A recipe found for one customer works well on similar designs. What is the flywheel path?",
		A: "Validate it with golden and differential tests, promote it into the recipe library, and record when to select it", D: []string{"Hard-code it as the new default for everyone immediately", "Ship the customer's design with it", "Delete it after the project"},
		Why: "Validated know-how is promoted with selection guidance; design data stays private."},
}

func registerFrontier(reg *generator.Registry) {
	rl := &edakit.Bank{Items: recipeRLItems}
	reg.Register("synth.frontier.recipe_rl", &edakit.Mix{Parts: []generator.Generator{edakit.Func(recipeReward), rl, rl}})
	qb := &edakit.Bank{Items: qorItems}
	reg.Register("synth.frontier.qor_estimation", &edakit.Mix{Parts: []generator.Generator{edakit.Func(qorRank), qb, qb}})
	reg.Register("synth.frontier.egraph_skill", &edakit.Bank{Items: egraphSkillItems})
	reg.Register("synth.frontier.cec_golden", &edakit.Bank{Items: cecGoldenItems})
	reg.Register("synth.frontier.change_tiers", &edakit.Bank{Items: changeTierItems})
}

// nearInts returns n distinct integers in [lo, hi] other than x, nearest first.
func nearInts(x, lo, hi, n int) []string {
	var out []string
	for d := 1; len(out) < n && d <= hi-lo; d++ {
		for _, c := range []int{x + d, x - d} {
			if c >= lo && c <= hi && len(out) < n {
				out = append(out, edakit.IntStr(c))
			}
		}
	}
	return out
}
