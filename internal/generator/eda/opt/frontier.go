package opt

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var weightTierItems = []edakit.Item{
	{KP: 1, Q: "A user asks the placer agent to favor timing over wirelength on this design. Which change tier should this land in?",
		A: "Tier 1: raise the timing weight in the objective configuration", D: []string{"Tier 2: write a new placer plug-in", "Tier 3: edit the placer's legalization core", "No change is possible without retraining"},
		Why: "Objective weights are parameters exposed by the engine; adjusting them needs no recompile and no new code."},
	{KP: 1, Q: "A user wants the router to penalize via usage on M2 twice as much as elsewhere. The router exposes per-layer via costs. Which tier?",
		A: "Tier 1: a value change in the per-layer cost table", D: []string{"Tier 2: a new rip-up policy", "Tier 3: a core path-model edit", "Tier 2: a new detailed router"},
		Why: "When the cost model already has the knob, the request is configuration; a plug-in is needed only for a new cost form."},
	{KP: 1, Q: "A user wants a congestion penalty shaped like a cubic of overflow, but the engine only offers linear and quadratic penalties. Which tier?",
		A: "Tier 2: a new penalty function at the cost-function extension point", D: []string{"Tier 1: change a weight", "Tier 3: rewrite the legalizer", "It cannot be expressed"},
		Why: "A new functional form is new code, but it fits the existing cost-function seam, so it is a plug-in, not a core edit."},
	{KP: 1, Q: "Which of these belongs in the immutable core of a placement optimizer rather than in its parameters?",
		A: "The legality guarantee that the final placement has no overlaps and sits on sites", D: []string{"The density target", "The wirelength-versus-timing weight", "The number of Nesterov iterations"},
		Why: "Correctness invariants live in the core; tunable trade-offs are parameters the agent may change."},
	{KP: 1, Q: "An agent changes the target density from 0.7 to 0.6. What should the harness require before showing the result?",
		A: "A QoR comparison against the baseline on the reference designs", D: []string{"A full core regression and human code review", "Nothing; parameters cannot break anything", "A recompile of the engine"},
		Why: "Tier 1 changes need no rebuild, but they still change results, so the differential QoR check is the evidence."},
	{KP: 1, Q: "Why does pushing most requests into Tier 1 objective weights make an optimization engine easier for agents to drive?",
		A: "Weight changes are fast, safe and reversible, and need no compile-test cycle", D: []string{"Weights change the core algorithm", "Weights are invisible to QoR", "Tier 1 changes skip all testing"},
		Why: "A well-placed seam turns the common request into configuration, so the agent's loop is seconds long and low risk."},
	{KP: 3, Q: "An agent raises the timing weight from 0.3 to 0.6. Which result best demonstrates the change worked as intended?",
		A: "A differential report: WNS and TNS improve while wirelength and congestion change by a bounded amount", D: []string{"The run completes without errors", "Wirelength drops", "The configuration file is valid JSON"},
		Why: "The proof of a weight change is a measured move along the intended trade-off, compared with the baseline."},
	{KP: 3, Q: "Raising a weight moves the chosen design from one Pareto point to another. What should the agent report to the user?",
		A: "Both objective deltas, so the user sees what was given up for what was gained", D: []string{"Only the improved objective", "Only the new weight value", "Nothing if the run passed"},
		Why: "Scalarization trades objectives; reporting one side hides the cost."},
	{KP: 3, Q: "A user's preferred design sits in a nonconvex dent of the Pareto front. Sweeping weights never reaches it. What should the agent switch to?",
		A: "An epsilon-constraint formulation: optimize one objective with a bound on the other", D: []string{"A larger weight range", "Random weights", "A higher learning rate"},
		Why: "Weighted sums reach only the convex hull of the front; epsilon constraints reach every Pareto point."},
	{KP: 3, Q: "Two weight settings give QoR differing by 0.4%, while seed-to-seed noise is about 1%. What should the agent conclude?",
		A: "The difference is not significant; run more seeds before claiming an improvement", D: []string{"The second setting is better", "The first setting is better", "Weights do not matter"},
		Why: "An effect smaller than the noise floor needs repeated runs or paired comparisons before it means anything."},
	{KP: 3, Q: "Which golden test is appropriate for a Tier 1 weight change in a placer?",
		A: "Rerun the reference designs and check legality plus QoR within agreed bounds of the baseline", D: []string{"Check that the source compiles", "Run unit tests of the parser only", "Compare file sizes of the outputs"},
		Why: "Legality is a pass or fail invariant; QoR is checked against a baseline with tolerances."},
	{KP: 3, Q: "A user asks for timing focus only in one clock domain. The engine's timing weight is global. What is the right response?",
		A: "This needs a per-path-group weight: a Tier 2 extension, or Tier 1 if the engine already exposes group weights", D: []string{"Set the global weight very high", "Edit the timer core", "Ignore the clock domain"},
		Why: "Scoping a weight needs the context (path group) to be visible at the seam; if it is not, the seam is extended by a plug-in."},
}

// weightSwitch asks which design a weighted-sum objective selects after a weight change.
func weightSwitch(ctx generator.GeneratorContext) generator.Problem {
	for {
		n := 4
		pts := drawDesigns(n, 2)
		w := float64(edakit.Int(1, 9)) / 10
		best, bi, second := math.Inf(1), -1, math.Inf(1)
		for i, p := range pts {
			v := w*float64(p[0]) + (1-w)*float64(p[1])
			if v < best {
				second = best
				best, bi = v, i
			} else if v < second {
				second = v
			}
		}
		if second-best < 1e-9 {
			continue
		}
		var opts []string
		for i := range pts {
			opts = append(opts, edakit.Sprintf("Run %d", i+1))
		}
		var desc []string
		for i, p := range pts {
			desc = append(desc, edakit.Sprintf("Run %d: WNS penalty %d, wirelength %d", i+1, p[0], p[1]))
		}
		ans := opts[bi]
		var dis []string
		for i, o := range opts {
			if i != bi {
				dis = append(dis, o)
			}
		}
		q := "An agent exposes the placer's objective weights as Tier 1 parameters. Normalized results of four runs (lower is better): " + strings.Join(desc, "; ") + edakit.Sprintf(". The user sets the timing weight to %s and the wirelength weight to %s. Which run does the weighted objective pick?", f(w, 1), f(1-w, 1))
		return edakit.Choice(q, ans, edakit.Sprintf("Score each run as %s times WNS penalty plus %s times wirelength; Run %d is lowest at %s.", f(w, 1), f(1-w, 1), bi+1, f(best, 2)), dis...)
	}
}

var moveSetItems = []edakit.Item{
	{KP: 1, Q: "In a simulated-annealing placer, which part is the natural Tier 2 policy seam?",
		A: "The move generator that proposes the next perturbation", D: []string{"The legality checker", "The netlist database", "The random-number generator's seed"},
		Why: "Move sets decide how the search explores; they vary by design style and fit a plug-in interface, while legality stays in the core."},
	{KP: 1, Q: "A customer's designs have many regular datapath arrays. An agent adds a move that shifts a whole aligned column of cells at once. Which tier is this?",
		A: "Tier 2: a new move plug-in behind the move-set interface", D: []string{"Tier 1: a weight change", "Tier 3: a core legality edit", "Not a change to the optimizer"},
		Why: "A new move is new code at a defined extension point; the annealer's acceptance rule and legality core are unchanged."},
	{KP: 1, Q: "What context must a move-set seam expose to a plug-in so it can propose useful moves?",
		A: "Read access to current positions, connectivity and local cost deltas, plus a way to return proposed moves", D: []string{"Write access to the legalizer's internals", "The engine's source code", "Nothing but a random seed"},
		Why: "A narrow interface gives the plug-in the state it reads and a typed way to propose; the core still applies and checks the move."},
	{KP: 1, Q: "Why should the move-set plug-in call be coarse, for example one call per batch of moves rather than per move?",
		A: "So the hot inner loop stays compiled and inlined in the core", D: []string{"So the plug-in can change legality", "Because per-move calls are illegal", "To make results nondeterministic"},
		Why: "Fine-grained virtual calls in the innermost loop cost performance; batching keeps the kernel fast."},
	{KP: 2, Q: "A plug-in move swaps two cells only if both are in the same timing-critical path group. What must the core still guarantee?",
		A: "That every accepted move leaves a legal, overlap-free placement", D: []string{"That the move improves timing", "That the move is random", "That the plug-in is never called"},
		Why: "Policies decide what to try; the core owns correctness, so legality holds whatever the plug-in proposes."},
	{KP: 2, Q: "An agent proposes a new rip-up-and-reroute ordering for a router. Which tier and what evidence?",
		A: "Tier 2 strategy plug-in, proven by DRC-clean golden designs and a QoR comparison", D: []string{"Tier 1, no evidence needed", "Tier 3, proven by compiling", "Tier 1, proven by a unit test of the parser"},
		Why: "Ordering is strategy; it lands at the rip-up policy seam and is judged by legality plus QoR against the baseline."},
	{KP: 2, Q: "A move-set plug-in uses a global variable to remember its last move across calls. Why is that a problem?",
		A: "Hidden state breaks determinism and reproducibility of runs", D: []string{"Global variables are slower", "It changes the cooling schedule", "It makes moves illegal"},
		Why: "Engines must be deterministic and seedable; state belongs in the explicit context passed through the seam."},
	{KP: 2, Q: "A new cost function for a router is supplied as a plug-in. What must its golden test confirm in addition to QoR?",
		A: "Determinism: the same inputs and seed reproduce the same result", D: []string{"That it is written in Python", "That it ignores congestion", "That it reduces runtime"},
		Why: "Reproducibility lets differential testing attribute QoR changes to the plug-in rather than to noise."},
	{KP: 3, Q: "A move plug-in found for one customer's design improves QoR by 2% on it and 0% on others. In the skill flywheel, what happens next?",
		A: "Validate it on golden designs, then promote it to the library with a manifest saying where it helps", D: []string{"Make it the default for all designs immediately", "Delete it", "Merge it into the core"},
		Why: "Promotion comes after validation, and the manifest records the design characteristics that predict a gain, so selection can be selective."},
	{KP: 3, Q: "An agent can choose among five compiled move-set policies for an annealer. What should it read to choose?",
		A: "Each policy's manifest: what it targets and on which design characteristics it was validated", D: []string{"The policies' source code line by line", "Nothing; choose at random", "The engine's legalization code"},
		Why: "Selection is a cheap text decision over manifests; the expensive compute stays in the engines."},
	{KP: 3, Q: "Which change would wrongly push a Tier 2 move-set request into Tier 3?",
		A: "Editing the annealer's acceptance and legality code to add the new move inline", D: []string{"Registering a new move class with the factory", "Adding a move-probability parameter", "Writing a manifest for the new move"},
		Why: "If the seam exists, adding a move should never touch the core; inlining it is a core edit with full regression cost."},
	{KP: 3, Q: "A floorplanning annealer's moves M1 to M3 are exposed as plug-ins with probabilities as parameters. A user wants more operand swaps. Which tier?",
		A: "Tier 1: raise the M1 probability", D: []string{"Tier 2: a new move", "Tier 3: a new Polish-expression evaluator", "Tier 2: a new cooling schedule"},
		Why: "Changing how often an existing move is used is a parameter; no new code is needed."},
}

// rewardPick asks which candidate a tool-grounded reward prefers.
func rewardPick(ctx generator.GeneratorContext) generator.Problem {
	for {
		lam := float64(edakit.Int(2, 10))
		mu := float64(edakit.Int(1, 8)) / 100
		type run struct{ wl, ovf, viol float64 }
		var rs []run
		var desc []string
		best, bi := math.Inf(-1), -1
		tie := false
		for i := 0; i < 4; i++ {
			r := run{float64(edakit.Int(85, 110)) / 100, float64(edakit.Int(0, 6)), float64(edakit.Int(0, 3))}
			rs = append(rs, r)
			v := -r.wl - lam*r.ovf/100 - mu*r.viol
			if v > best+1e-9 {
				best, bi, tie = v, i, false
			} else if math.Abs(v-best) < 1e-9 {
				tie = true
			}
			desc = append(desc, edakit.Sprintf("Run %d: wirelength ratio %s, overflow %s edges, DRC violations %s", i+1, f(r.wl, 2), f(r.ovf, 0), f(r.viol, 0)))
		}
		if tie {
			continue
		}
		var dis []string
		for i := range rs {
			if i != bi {
				dis = append(dis, edakit.Sprintf("Run %d", i+1))
			}
		}
		q := edakit.Sprintf("An RL routing agent is rewarded by the tool with $R = -\\text{WL ratio} - %s \\cdot \\text{overflow}/100 - %s \\cdot \\text{DRC violations}$. Candidate results: ", f(lam, 0), f(mu, 2)) + strings.Join(desc, "; ") + ". Which run earns the highest reward?"
		return edakit.Choice(q, edakit.Sprintf("Run %d", bi+1),
			edakit.Sprintf("Compute R for each run; Run %d scores highest at %s. Note how the penalty weights decide whether a short but dirty route can beat a longer clean one.", bi+1, f(best, 3)),
			dis...)
	}
}

var rewardItems = []edakit.Item{
	{KP: 2, Q: "An RL placement agent is rewarded by HPWL alone and learns to pile cells into one corner. What happened?",
		A: "Reward hacking: the proxy ignores density, so an illegal clumped placement scores well", D: []string{"The agent found the true optimum", "HPWL is a signoff metric", "The legalizer improved the result"},
		Why: "A reward that omits a constraint invites the agent to violate it; density or legality must be in the reward or enforced by the environment."},
	{KP: 2, Q: "A sizing agent is rewarded by a fast surrogate timer's WNS. Signoff STA later shows WNS worse than the baseline. What is the likely cause?",
		A: "The agent exploited errors of the fast proxy that the accurate timer does not share", D: []string{"Signoff STA is wrong", "Rewards cannot use timing", "The agent was not trained long enough, nothing else"},
		Why: "Optimizing hard against a proxy finds the proxy's blind spots; periodic recursion to the accurate engine catches the gap."},
	{KP: 2, Q: "A routing agent's reward counts DRC violations reported by the checker. The agent learns to produce outputs that crash the checker, which then reports zero. What is the fix?",
		A: "Treat checker failure as the worst reward and validate the output format before scoring", D: []string{"Remove DRC from the reward", "Reward crashes", "Run the checker less often"},
		Why: "The environment must fail closed: an unscorable output gets the worst score, never a free pass."},
	{KP: 2, Q: "Which reward is hardest to hack for a placement agent?",
		A: "QoR measured after legalization and routing by the real tools", D: []string{"HPWL of the unlegalized placement", "The agent's own predicted wirelength", "The number of moves the agent made"},
		Why: "Downstream, tool-grounded metrics close the loopholes of early proxies, at the price of slower and costlier reward evaluation."},
	{KP: 2, Q: "A buffer-insertion agent gets reward for improved slack but none for area. What should be expected?",
		A: "It will insert far more buffers than needed", D: []string{"It will insert no buffers", "Area will shrink", "Slack will get worse"},
		Why: "Unpriced resources are free to the agent; every real cost must appear in the reward or as a constraint."},
	{KP: 2, Q: "An agent's reward is the improvement in TNS over the starting netlist. It learns to slightly degrade the start before optimizing. Why can it?",
		A: "The baseline is under the agent's control, so it can inflate its own improvement", D: []string{"TNS is not measurable", "The reward is too sparse", "The tool is nondeterministic"},
		Why: "Rewards must be anchored to baselines the agent cannot modify, such as a fixed reference run."},
	{KP: 3, Q: "A fast congestion estimator is used as the reward for most episodes. How should the accurate global router enter the loop?",
		A: "Score a sample of episodes with the accurate router to measure and correct proxy drift", D: []string{"Never: it is too slow", "Replace the agent with the router", "Only at the end of training"},
		Why: "Periodic accurate evaluation detects when the agent starts exploiting the fast proxy and recalibrates the reward."},
	{KP: 3, Q: "Rewards from a placement tool vary by 2% across seeds for the same action. What does that do to learning?",
		A: "It adds variance; repeated or paired evaluations are needed to separate real gains from noise", D: []string{"Nothing", "It makes the reward exact", "It removes reward hacking"},
		Why: "Noisy rewards slow learning and invite overfitting to lucky seeds."},
	{KP: 3, Q: "Why are deterministic solvers attractive as reward sources for training optimization agents?",
		A: "They score outputs reproducibly against objectives and constraints, with ground truth", D: []string{"They are always fast", "They cannot be gamed in any way", "They need no inputs"},
		Why: "Tool-grounded rewards are checkable and repeatable, though agents can still game whatever the tool fails to check."},
	{KP: 3, Q: "An RL agent improves the fast-proxy reward by 8% but signoff QoR by 1%. How should results be reported?",
		A: "Report signoff QoR as the result and the proxy gain only as a diagnostic", D: []string{"Report the 8% as the result", "Average the two numbers", "Report neither"},
		Why: "The proxy-to-signoff gap is the central risk; claims must be made in signoff terms."},
	{KP: 3, Q: "What is a sparse reward problem for an agent doing a full place-and-route flow?",
		A: "Reward arrives only after hours-long runs, so few learning signals per day", D: []string{"The reward is always zero", "The reward has too many terms", "The reward is computed per cell"},
		Why: "Long episodes give few samples; shaping with intermediate tool metrics or fast engines adds signal, at a proxy risk."},
	{KP: 3, Q: "An agent's reward includes a hard legality check that returns minus infinity on any overlap. What is a practical downside?",
		A: "Most early actions score the same, giving no gradient toward legality", D: []string{"It allows overlaps", "It makes rewards nondeterministic", "It removes the need for a legalizer"},
		Why: "A flat penalty carries no information about how close to legal an output was; a graded penalty or an environment-side legalizer helps learning."},
}

// successiveHalving asks for the compute spent by a successive-halving exploration.
func successiveHalving(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Int(2, 5)
	n := 1 << k
	b := edakit.Int(1, 6) * 10
	total := 0
	cfg := n
	var rounds []string
	for r := 0; cfg >= 1; r++ {
		cost := cfg * b * (1 << r)
		total += cost
		rounds = append(rounds, edakit.IntStr(cost))
		if cfg == 1 {
			break
		}
		cfg /= 2
	}
	q := edakit.Sprintf("A design-space exploration agent starts %d flow configurations with %d minutes each, then repeatedly keeps the better half and doubles each survivor's budget, until one configuration remains and gets its run. How many machine-minutes does the whole exploration use?", n, b)
	return edakit.Choice(q, edakit.IntStr(total),
		edakit.Sprintf("Each round costs (configs times minutes): %s. With halving and doubling every round costs the same, %d rounds of %d, total %d, against %d for running every configuration at the final budget.", strings.Join(rounds, ", "), k+1, n*b, total, n*b*n),
		edakit.Distinct(total, 3, n*b*n, n*b, n*b*k, total+n*b)...)
}

var dseItems = []edakit.Item{
	{KP: 2, Q: "An exploration agent has budget for 30 full flow runs over 6 knobs. Which next-experiment rule uses the budget best?",
		A: "Fit a surrogate to completed runs and pick the next point by an acquisition function", D: []string{"A full grid over 6 knobs", "Repeat the best run 30 times", "Pick knobs alphabetically"},
		Why: "Model-guided selection makes each run informative; a 6-knob grid with only a few levels already needs hundreds of runs."},
	{KP: 2, Q: "The agent can run a fast placer in 2 minutes or the full flow in 3 hours. How should it explore?",
		A: "Screen many configurations with the fast engine and recurse to the full flow on the promising few", D: []string{"Use only the full flow", "Use only the fast placer and ship its winner", "Alternate randomly"},
		Why: "Fast engines map the space cheaply; accurate engines confirm, because fast-engine rankings can disagree with signoff."},
	{KP: 2, Q: "When should a DSE agent stop trusting the fast engine's ranking and run the accurate flow?",
		A: "When candidates differ by less than the fast engine's known error, or the fast metrics' variance is high", D: []string{"Never", "Only after every configuration is tried", "When the fast engine finishes"},
		Why: "If the fast estimate cannot separate candidates, only the accurate engine can."},
	{KP: 2, Q: "A DSE agent returns one best configuration for a power, performance and area study. What should it return instead?",
		A: "The Pareto set of configurations, so the designer picks the trade-off", D: []string{"The average configuration", "The fastest run", "The configuration with the most knobs changed"},
		Why: "Multi-objective exploration produces a front; collapsing it to one point hides the trade-offs."},
	{KP: 2, Q: "Many of the agent's runs fail with routing DRC errors. How should the exploration treat them?",
		A: "Record them as infeasible and model the feasibility boundary to avoid that region", D: []string{"Drop them silently", "Score them as best", "Restart the whole exploration"},
		Why: "Failed runs are information: a feasibility model keeps the agent from wasting budget in infeasible regions."},
	{KP: 2, Q: "The agent plans 8 parallel runs per round. What should the batch look like?",
		A: "A mix of high-predicted-value points and high-uncertainty points", D: []string{"Eight copies of the current best", "Eight random points every round", "Only the most uncertain points"},
		Why: "Balancing exploitation and exploration within a batch keeps both progress and learning."},
	{KP: 3, Q: "Two configurations differ by 0.5% in WNS, and seed noise is about 1%. What should the agent do before reporting a winner?",
		A: "Run multiple seeds of both and compare distributions", D: []string{"Report the higher one", "Average their knobs", "Stop exploring"},
		Why: "A difference inside the noise band is not a result; replication separates signal from seed luck."},
	{KP: 3, Q: "Each flow run touches the customer's design database. Where must the exploration agent's runs execute?",
		A: "On-site inside the IP boundary, in sandboxes with resource limits", D: []string{"On any public cloud the agent prefers", "On the agent vendor's servers", "Wherever runs fastest, regardless of data"},
		Why: "Design data stays on-site; the agent orchestrates runs within the boundary rather than exporting the design."},
	{KP: 3, Q: "Why should every exploration run record its exact configuration, tool versions and seed?",
		A: "So any result can be reproduced and differential comparisons are meaningful", D: []string{"To make runs faster", "To fill the disk", "Because the surrogate needs file names"},
		Why: "Reproducibility is the basis of trust in an automatically found configuration."},
	{KP: 3, Q: "An agent's best configuration on design A transfers poorly to design B. How should the shared domain model learn from this?",
		A: "Record which design characteristics predicted success, so selection becomes conditional on them", D: []string{"Discard design A's results", "Make A's configuration the default anyway", "Stop exploring new designs"},
		Why: "The know-how worth keeping is when a configuration helps, not a single universal setting."},
	{KP: 3, Q: "Grid search over 4 knobs at 5 levels each needs how many runs, and what does that imply?",
		A: "625 runs; grids scale exponentially in the number of knobs", D: []string{"20 runs; grids scale linearly", "9 runs; grids are cheap", "4 runs; one per knob"},
		Why: "$5^4 = 625$; each extra knob multiplies the count, which is why model-guided or random search wins in higher dimensions."},
	{KP: 3, Q: "Which output makes a DSE agent's recommendation reviewable by a human?",
		A: "A summary of runs tried, the front found, and the evidence for the recommended point", D: []string{"Only the final configuration file", "The surrogate's raw weights", "The number of runs"},
		Why: "A human approves with evidence: what was tried, what the trade-offs were, and why this point."},
}

// geomeanQoR asks for the geometric-mean QoR ratio of a candidate against the baseline.
func geomeanQoR(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(3, 5)
	rs := make([]float64, n)
	prod := 1.0
	ar := 0.0
	var parts []string
	for i := range rs {
		rs[i] = float64(edakit.Int(90, 108)) / 100
		prod *= rs[i]
		ar += rs[i]
		parts = append(parts, f(rs[i], 2))
	}
	g := math.Pow(prod, 1/float64(n))
	q := edakit.Sprintf("A new optimizer skill is compared with the baseline on %d golden designs. The wirelength ratios (new over baseline) are %s. What is the geometric-mean ratio?", n, strings.Join(parts, ", "))
	return edakit.Choice(q, f(g, 3),
		edakit.Sprintf("Multiply the ratios and take the %d-th root: %s. The geometric mean treats a 10%% gain and a 10%% loss symmetrically in log terms, so one large design cannot dominate.", n, f(g, 3)),
		numChoices(g, 3, ar/float64(n)+0.004, prod, g+0.01, g-0.012)...)
}

var verifyItems = []edakit.Item{
	{KP: 2, Q: "An LP-based legalizer returns a solution. Which certificate lets a test check optimality without trusting the solver?",
		A: "A dual-feasible solution whose objective equals the primal objective", D: []string{"The solver's log saying optimal", "The runtime", "The number of iterations"},
		Why: "Primal feasibility, dual feasibility and equal objectives prove optimality by LP duality; the check is cheap and independent."},
	{KP: 2, Q: "For a convex QP sizing step, what does a golden test check on the returned point?",
		A: "That the KKT residuals (stationarity, feasibility, complementarity) are below a tolerance", D: []string{"That the objective is negative", "That the point is integer", "That the gradient is large"},
		Why: "For convex problems KKT conditions are sufficient for optimality, so small residuals certify the answer."},
	{KP: 2, Q: "A MILP floorplanner returns incumbent 1040 with best bound 1000 (minimization). What can a test assert?",
		A: "The solution is within 4% of optimal", D: []string{"The solution is optimal", "The solution is infeasible", "Nothing at all"},
		Why: "The bound brackets the optimum: (1040 - 1000)/1000 = 4%, a provable quality bound."},
	{KP: 2, Q: "An annealing placer gives different results on two runs with the same seed. What does that break?",
		A: "Differential testing: a QoR change can no longer be attributed to the code change", D: []string{"Nothing; annealing is random", "Only the runtime report", "The legality of the placement"},
		Why: "Seeded determinism is what makes before-and-after comparisons meaningful."},
	{KP: 2, Q: "Which invariant must a golden test check on every placer output regardless of QoR?",
		A: "Legality: no overlaps, cells on rows and sites, fixed cells unmoved", D: []string{"Wirelength below the baseline", "Runtime below a minute", "Exactly the baseline coordinates"},
		Why: "Correctness invariants are pass or fail; QoR is compared with tolerances."},
	{KP: 2, Q: "A branch-and-bound solver change is meant only to speed up node selection. What should the differential test show?",
		A: "The same optimal objective on every golden instance, with fewer nodes or less time", D: []string{"A better optimal objective", "Different optimal objectives", "More nodes explored"},
		Why: "Search-order changes cannot change the optimum of an exactly solved instance; only effort should move."},
	{KP: 3, Q: "A new move-set skill improves wirelength on 9 golden designs and worsens it 6% on 1. What should the promotion decision consider?",
		A: "Whether the regression is explained and bounded, and whether the manifest should exclude that design class", D: []string{"Promote: the average improved", "Reject: any regression is fatal", "Ignore the outlier"},
		Why: "Aggregate gains can hide a serious regression; the manifest and selection rules should capture where the skill helps."},
	{KP: 3, Q: "Why compare optimizer skills with the geometric mean of ratios rather than the arithmetic mean of absolute numbers?",
		A: "Ratios normalize design size, and the geometric mean treats gains and losses symmetrically", D: []string{"It is always larger", "It ignores regressions", "It needs fewer designs"},
		Why: "Absolute wirelength is dominated by the largest designs; the geometric mean of ratios weighs each design equally."},
	{KP: 3, Q: "A Tier 1 weight change and a Tier 2 plug-in both pass golden tests. Which needs more review before promotion?",
		A: "The plug-in, because it is new code that can affect behavior beyond the tested cases", D: []string{"The weight change", "Neither", "Both need a core regression"},
		Why: "New code carries more risk than a parameter; its review and test scope are correspondingly larger."},
	{KP: 3, Q: "A differential test shows a placer change moved only the cells it was meant to move, by the expected amounts. What has been proven?",
		A: "That the change's effect matches its intent on the reference designs", D: []string{"That the change is optimal", "That every future design benefits", "That the core is unchanged"},
		Why: "Differential checks prove intent on references; they do not prove optimality or generality."},
	{KP: 3, Q: "Golden tests for a Lagrangian-relaxation sizer pass, but the duality gap on one design is 15%. What does this suggest?",
		A: "The sizer may be far from optimal on that design, so the multiplier updates or discretization deserve scrutiny", D: []string{"The sizer is optimal", "The gap is irrelevant", "The design is infeasible"},
		Why: "A large gap between the relaxed bound and the feasible solution means quality may have been left behind."},
	{KP: 3, Q: "Why must an optimizer's golden tests run in a sandbox with time and memory limits?",
		A: "A bad change can loop or blow up memory, and the harness must fail it cleanly", D: []string{"To make results better", "To allow network access", "To skip legality checks"},
		Why: "Limits turn runaway behavior into a reported failure instead of a hung harness."},
}
