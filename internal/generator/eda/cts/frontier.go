package cts

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var irEcoItems = []edakit.Item{
	{KP: 1, Q: "An IR-drop ECO agent receives the signoff report. What should it read first to propose fixes?",
		A:   "The per-node or per-instance drop map, with the worst regions, the layers and the instances whose current causes them",
		D:   []string{"Only the total chip power number", "The RTL of the design", "The list of timing exceptions"},
		Why: "Fixes are local: the agent needs where the drop is, which resistance path causes it, and which loads draw the current."},
	{KP: 1, Q: "Which actions are typical IR-drop ECOs an agent may propose?",
		A:   "Add or widen straps, add via stacks, insert decap near hot spots, or spread high-current cells apart",
		D:   []string{"Change the clock period in the SDC", "Delete the hottest cells", "Lower the supply voltage of the whole chip"},
		Why: "Grid resistance, local charge and current density are the levers; the agent proposes edits through the PDN and placement engines."},
	{KP: 1, Q: "A hot spot shows large dynamic drop but small static drop. Which ECO is the better first proposal?",
		A:   "Insert decap next to the switching cluster",
		D:   []string{"Widen distant straps on the top layer", "Add more power bumps on the far side of the die", "Reduce the static leakage of the region"},
		Why: "Dynamic droop is a local charge problem, which decap solves; static drop is a resistance problem, which straps solve."},
	{KP: 1, Q: "In the loop read report, propose ECO, run engines, compare, what does the agent run after applying a strap ECO?",
		A:   "An IR re-analysis plus DRC and timing checks on the changed area, compared against the previous run",
		D:   []string{"Nothing until tape-out", "Only a compile of the agent's own code", "A full logic synthesis from RTL"},
		Why: "A strap can fix drop but steal routing tracks or create DRCs; the loop must re-verify every metric the edit could disturb."},
	{KP: 1, Q: "The agent proposes adding a strap through a region where signal routing is already congested. What should the loop do?",
		A:   "Weigh the routing cost: check congestion and timing after the ECO, and prefer decap or via fixes if they meet the budget",
		D:   []string{"Add the strap regardless", "Remove signal nets to make room", "Raise the IR budget instead"},
		Why: "IR fixes compete with routing; a good loop scores every candidate on all the metrics it touches."},
	{KP: 1, Q: "Which part of an IR ECO loop must remain a human decision?",
		A:   "Approving the final ECO for signoff, after the agent shows the before and after reports",
		D:   []string{"Choosing which solver iteration to run", "Every individual strap width trial", "Reading the drop map"},
		Why: "The agent iterates and verifies; a human approves the change that goes into the signed-off database."},
	{KP: 3, Q: "Each strap ECO is checked with a fast IR estimate inside the loop. When must the agent rerun the full signoff solver?",
		A:   "Before reporting a fix as done, and whenever the fast estimate is near the budget or outside the conditions it was calibrated on",
		D:   []string{"Never; the fast estimate is final", "Only if the design fails timing", "After tape-out"},
		Why: "Fast estimates steer; signoff numbers decide. Margins near the threshold are where approximation error matters."},
	{KP: 3, Q: "After an ECO the worst static drop improved by 8 mV, but EM violations appeared on a via array feeding the new strap. What went wrong?",
		A:   "The loop optimised drop without checking current density on the paths that now carry more current",
		D:   []string{"Nothing, because EM is unrelated to the grid", "The IR solver produced wrong voltages", "The via array is too large"},
		Why: "Lower resistance attracts current; new low-resistance paths must pass EM checks too."},
	{KP: 3, Q: "Two candidate ECOs each fix the hot spot alone. Applied together, the improvement is less than the sum. Why?",
		A:   "Grid fixes interact: both relieve the same resistive path, so their effects overlap rather than add",
		D:   []string{"The solver is non-deterministic", "ECOs always cancel each other", "Drop is unaffected by the grid"},
		Why: "IR drop is a solution of one linear system; sensitivities are local estimates and overlap when fixes share a current path."},
	{KP: 3, Q: "Why should an IR ECO agent express edits through the PDN and placement engines rather than by editing DEF text directly?",
		A:   "The engines enforce legality and keep the database consistent, so the edit is valid and reproducible",
		D:   []string{"Editing DEF is impossible", "Engines are faster than text editing in every case", "DEF files do not contain power nets"},
		Why: "The engines own the legality core; the agent works at the tool interface, which keeps edits checkable and reversible."},
	{KP: 3, Q: "An agent proposes fixing every hot spot by inserting decap into all free sites. What is the cost it is ignoring?",
		A:   "Decap leakage and area, and the loss of sites needed later for timing ECO buffers",
		D:   []string{"Decap increases static IR drop", "Decap cannot be placed in free sites", "Decap changes the clock frequency"},
		Why: "Decap is not free: gate-oxide decap leaks, and filled sites cannot absorb later fixes."},
	{KP: 3, Q: "Which record should accompany each ECO the agent applies?",
		A:   "The trigger (report and region), the edit, the before and after metrics, and the engine versions used",
		D:   []string{"Only the final drop number", "Nothing, because the database holds the result", "The agent's full prompt history and nothing else"},
		Why: "Provenance makes ECOs reviewable, reproducible and reversible, which the human approval step depends on."},
}

// irEcoChoice: choose the cheapest candidate whose predicted reduction meets the gap.
func irEcoChoice(ctx generator.GeneratorContext) generator.Problem {
	names := []string{
		"widen the M6 strap over the region",
		"add a second M7 strap",
		"add via stacks between M2 rails and M5 straps",
		"insert decap cells around the cluster",
		"spread the high-current cells apart",
	}
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		drop := float64(edakit.Int(45, 90))
		budget := float64(edakit.Int(30, 45))
		if drop-budget < 8 {
			continue
		}
		nm := edakit.Shuffle(names)
		var cs []ecoCand
		for i := 0; i < 4; i++ {
			cs = append(cs, ecoCand{nm[i], math.Max(1, drop-budget+float64(edakit.Int(-6, 10))), float64(edakit.Int(1, 9))})
		}
		if best, ok := ecoBest(cs, drop-budget); ok {
			return ecoProblem(drop, budget, cs, best)
		}
	}
	// Known-good instance: the cheapest candidate falls short and the answer is unique.
	cs := []ecoCand{{names[2], 25, 2}, {names[1], 32, 5}, {names[0], 38, 7}, {names[3], 30, 6}}
	best, _ := ecoBest(cs, 30)
	return ecoProblem(70, 40, cs, best)
}

type ecoCand struct {
	name      string
	red, cost float64
}

// ecoBest returns the cheapest sufficient candidate, and whether the instance is a good drill:
// at least two candidates suffice, the cheapest sufficient one is unique, and the overall
// cheapest candidate falls short.
func ecoBest(cs []ecoCand, gap float64) (int, bool) {
	best, ok := -1, 0
	for i, c := range cs {
		if c.red >= gap {
			ok++
			if best < 0 || c.cost < cs[best].cost {
				best = i
			}
		}
	}
	if best < 0 || ok < 2 {
		return best, false
	}
	cheapest := 0
	for i, c := range cs {
		if i != best && c.red >= gap && c.cost == cs[best].cost {
			return best, false
		}
		if c.cost < cs[cheapest].cost {
			cheapest = i
		}
	}
	return best, cs[cheapest].red < gap
}

func ecoProblem(drop, budget float64, cs []ecoCand, best int) generator.Problem {
	gap := drop - budget
	var lines []string
	for _, c := range cs {
		lines = append(lines, edakit.Sprintf("%s (predicted reduction %s mV, cost %s)", c.name, f(c.red, 0), f(c.cost, 0)))
	}
	q := edakit.Sprintf("An IR ECO agent sees a worst drop of %s mV against a budget of %s mV. A fast solver predicts the effect of single candidate ECOs: %s. Costs combine routing tracks and area. Which single ECO should it propose first?", f(drop, 0), f(budget, 0), strings.Join(lines, "; "))
	var ds []string
	for i, c := range cs {
		if i != best {
			ds = append(ds, c.name)
		}
	}
	why := edakit.Sprintf("It must close a gap of %s mV. Of the candidates whose predicted reduction is at least %s mV, the cheapest is to %s. The overall cheapest falls short, and larger fixes cost more than needed. The choice is then confirmed by the accurate solver.", f(gap, 0), f(gap, 0), cs[best].name)
	return edakit.Choice(q, cs[best].name, why, ds...)
}

var fastPdnItems = []edakit.Item{
	{KP: 1, Q: "Which is a typical fast PDN analysis and which the accurate one?",
		A:   "Fast: a coarsened or reduced grid with average currents; accurate: the full extracted grid solved with per-instance currents",
		D:   []string{"Fast: full extracted transient; accurate: a spreadsheet estimate", "Both solve exactly the same system at the same resolution", "Fast: silicon measurement; accurate: a reduced grid"},
		Why: "Model reduction, coarse meshes and averaged loads trade resolution for speed; signoff solves the full grid."},
	{KP: 1, Q: "Why is a fast IR estimate useful inside a placement or floorplanning loop?",
		A:   "It can be rerun after every change, steering decisions long before the full grid exists",
		D:   []string{"It is more accurate than signoff", "Placement cannot run without it", "It replaces EM analysis"},
		Why: "Early, frequent feedback is what a modeling loop needs; precision is spent later."},
	{KP: 1, Q: "A fast IR solver uses the conjugate gradient method with a loose tolerance. What does it trade away?",
		A:   "Residual accuracy of the solved voltages, in exchange for fewer iterations",
		D:   []string{"Symmetry of the conductance matrix", "The ability to handle more than one pad", "Kirchhoff's current law in the model"},
		Why: "Stopping early leaves a residual; the drops are approximate but much cheaper."},
	{KP: 1, Q: "Which structural property of the grid matrix lets both fast and accurate solvers use conjugate gradient?",
		A:   "It is symmetric positive definite once at least one node connects to the supply",
		D:   []string{"It is dense and diagonal", "It is upper triangular", "It has negative eigenvalues"},
		Why: "Conductance stamps are symmetric, and supply connections make the matrix positive definite, which is the CG requirement."},
	{KP: 2, Q: "The fast and accurate solvers agree within 1 mV almost everywhere, but differ by 9 mV near a macro edge. What should the loop do?",
		A:   "Send macro-edge regions to the accurate solver in future and record the gap to improve the fast model",
		D:   []string{"Average the two numbers", "Trust the fast number because it agrees elsewhere", "Disable the accurate solver"},
		Why: "Where the engines disagree is where recursion to the accurate engine is required; the pattern of disagreement is the signal."},
	{KP: 2, Q: "A fast estimate says the worst drop is 38 mV against a 40 mV budget. The fast model's typical error is 4 mV. What should the agent conclude?",
		A:   "The result is inconclusive; run the accurate solver before calling it a pass",
		D:   []string{"It passes with margin", "It certainly fails", "Raise the budget to 44 mV"},
		Why: "A margin smaller than the estimator's error does not separate pass from fail."},
	{KP: 2, Q: "A fast solver models each region's current as an average. When does that become most misleading?",
		A:   "For dynamic droop, when many cells switch together in a short window",
		D:   []string{"For static drop with steady loads", "For leakage-only analysis", "When the clock is gated off everywhere"},
		Why: "Averaging hides simultaneity; dynamic analysis needs time-resolved current."},
	{KP: 2, Q: "An agent runs the accurate dynamic IR solver after every small ECO and the loop becomes too slow. What is a better policy?",
		A:   "Use the fast solver for every trial and the accurate solver on candidates that pass and on borderline results",
		D:   []string{"Stop running any solver", "Run only the accurate solver, but less often, at random", "Use the fast solver as the signoff result"},
		Why: "Fast engines screen, accurate engines confirm; the selection rule decides where accuracy is spent."},
	{KP: 3, Q: "A surrogate model predicts IR drop from placement density maps. On a new design with an unusual macro arrangement, how should its output be treated?",
		A:   "As untrusted until checked: inputs outside the training distribution need the accurate solver",
		D:   []string{"As exact, since it was trained on many designs", "As an upper bound on the true drop", "As irrelevant to the PDN"},
		Why: "Learned surrogates are reliable only near their training data; distribution shift is when to recurse."},
	{KP: 3, Q: "Where is the fast-versus-accurate choice best exposed in a PDN analysis engine?",
		A:   "As a parameter or policy selecting the solver mode and tolerance, with the grid model and stamping in the core",
		D:   []string{"As a core edit for every new design", "As a hidden constant inside the solver loop", "In the RTL"},
		Why: "Accuracy mode is change surface; the conductance model and KCL stamping are invariant core."},
	{KP: 3, Q: "Why must a fast solver's speed-up be measured against the accurate solver on the same design and corner?",
		A:   "Speed is meaningful only together with the error it costs on comparable inputs",
		D:   []string{"Speed does not depend on the design", "The accurate solver has no runtime", "Errors are always zero for fast solvers"},
		Why: "A fast mode is a point on a speed-accuracy curve; both coordinates must be reported against the same reference."},
	{KP: 3, Q: "An agent's decisions are driven by the variance between fast estimates across nearby candidate edits. What does high variance tell it?",
		A:   "The region is sensitive, so it should evaluate those candidates with the accurate engine",
		D:   []string{"The fast solver has crashed", "The candidates are all equivalent", "The budget should be relaxed"},
		Why: "High sensitivity means small modeling errors can change the ranking, which is where the accurate engine earns its cost."},
}

var goldenItems = []edakit.Item{
	{KP: 1, Q: "What is a golden IR test for a PDN solver?",
		A:   "A fixed grid and load set with known-correct node voltages that every solver version must reproduce within tolerance",
		D:   []string{"The fastest run ever recorded", "A grid with no loads", "A screenshot of an IR map"},
		Why: "Golden tests pin the solver's correctness; changes are checked against them before anything else."},
	{KP: 1, Q: "Which small grid makes a good golden test because its answer is known exactly?",
		A:   "A uniform resistor ladder or small mesh whose drops can be solved by hand or with a dense direct solver",
		D:   []string{"The full production chip", "A random grid regenerated on every run", "A grid with no supply connection"},
		Why: "Hand-checkable or exactly solvable cases catch stamping and solver errors unambiguously."},
	{KP: 1, Q: "A new preconditioner is added to the IR solver. What must the differential check show?",
		A:   "Node voltages on reference designs match the previous solver within tolerance, and runtime or iterations improve",
		D:   []string{"Only that the code compiles", "That the drop map looks similar", "That runtime improved, whatever the voltages"},
		Why: "A solver change must not move answers; the differential check proves that the only change is speed."},
	{KP: 1, Q: "Which tolerance makes sense for a golden IR comparison?",
		A:   "An absolute tolerance in mV well below the IR budget margin, plus a relative one for large drops",
		D:   []string{"Exactly zero difference on all floating-point bits", "Any difference under 50% is acceptable", "No tolerance is needed for IR"},
		Why: "Floating-point order changes cause tiny differences; the tolerance must be small relative to decisions the result drives."},
	{KP: 1, Q: "A PDN ECO is meant to touch one region only. What should the differential IR report confirm?",
		A:   "Drops improved in the target region and moved negligibly everywhere else",
		D:   []string{"Every node improved by the same amount", "The worst drop is unchanged", "The total chip power fell"},
		Why: "Differential checks prove the edit did what was intended and nothing else."},
	{KP: 1, Q: "Why must golden IR tests also include EM current checks?",
		A:   "The same solve yields branch currents, and a solver bug can corrupt currents while voltages look plausible",
		D:   []string{"EM is checked only in silicon", "Currents cannot be derived from voltages", "EM limits do not depend on current"},
		Why: "Currents come from voltage differences across small resistances, which amplify relative errors."},
	{KP: 3, Q: "After a solver change, all golden voltages match within tolerance but one test's iteration count doubled. What should happen?",
		A:   "Flag it as a performance regression for review even though correctness passed",
		D:   []string{"Ignore it because the voltages match", "Delete the slow test", "Raise the tolerance"},
		Why: "QoR baselines include runtime; a doubled iteration count is a regression to be explained."},
	{KP: 3, Q: "A golden test fails only on a grid with a single supply pad. Which bug is most likely?",
		A:   "Mishandled boundary conditions or near-singular conditioning when only one supply connection anchors the system",
		D:   []string{"A typo in the report header", "Too many supply pads", "The clock frequency"},
		Why: "With one pad the matrix is barely positive definite, which exposes preconditioner and boundary handling bugs."},
	{KP: 3, Q: "Why should golden IR data live outside the solver's own repository settings that an agent can edit?",
		A:   "So that an optimising agent cannot change the reference it is judged against",
		D:   []string{"Golden data is too large for any repository", "Solvers cannot read external files", "Agents cannot read golden data"},
		Why: "The oracle must be read-only to whatever is being optimised, or the test can be gamed."},
	{KP: 3, Q: "Two IR tools disagree by 6 mV on a production block. Both pass their own golden tests. What should the investigation compare first?",
		A:   "The inputs: extracted grid resistance, pad models and current assignment, before the solvers themselves",
		D:   []string{"The colour scales of the IR maps", "The tools' release dates", "Nothing, since both passed golden tests"},
		Why: "Most cross-tool differences come from model inputs, not the linear solve; differential checks must fix the inputs first."},
	{KP: 3, Q: "A golden check passes on the static solve but the dynamic solver was changed too. What extra golden data is needed?",
		A:   "Reference waveforms of node voltage over time for a known current profile, checked at chosen times or peaks",
		D:   []string{"Nothing: the static check covers it", "Only the final steady-state voltages", "The leakage report"},
		Why: "A dynamic change must be checked on dynamic quantities: droop peaks and their timing."},
	{KP: 3, Q: "Which property makes a differential IR test trustworthy across runs?",
		A:   "Determinism: the same inputs and seeds give the same voltages, so any change is attributable to the edit",
		D:   []string{"Random initialisation on every run", "Running on whichever machine is free with different thread counts and no seed", "Hiding the solver tolerance"},
		Why: "Non-determinism blurs real changes with run-to-run noise."},
}

// goldenCount: how many nodes of a differential IR report exceed tolerance.
func goldenCount(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(6, 8)
	tol := pick(0.5, 2, 0.5)
	bad := 0
	var rows []string
	for i := 0; i < n; i++ {
		g := pick(5, 45, 0.1)
		var d float64
		if edakit.Int(0, 2) == 0 {
			d = tol + pick(0.1, 3, 0.1)
		} else {
			d = pick(0, tol-0.1, 0.1)
		}
		if edakit.Int(0, 1) == 1 {
			d = -d
		}
		c := edakit.Round(g+d, 1)
		if absF(c-g) > tol+1e-9 {
			bad++
		}
		rows = append(rows, edakit.Sprintf("n%d: %s vs %s", i+1, f(g, 1), f(c, 1)))
	}
	q := edakit.Sprintf("A golden IR test compares a new solver's drops with reference drops (mV, reference vs new): %s. The tolerance is %s mV absolute. How many nodes fail?", strings.Join(rows, "; "), f(tol, 1))
	ds := edakit.Distinct(bad, 3, bad+1, bad-1, n-bad)
	return edakit.Choice(q, edakit.IntStr(bad),
		edakit.Sprintf("A node fails when $|d_{new} - d_{ref}| > %s$ mV. Counting the differences gives %d failing nodes, so this change cannot be promoted until they are explained.", f(tol, 1), bad), ds...)
}

var policyItems = []edakit.Item{
	{KP: 1, Q: "In a CTS engine built on the four-layer model, which part is the immutable core?",
		A:   "Legal buffer placement, the delay and slew evaluation API, and the tree data structure's invariants",
		D:   []string{"The choice of clustering heuristic", "The skew target value", "The order in which clock domains are processed"},
		Why: "Correctness and legality stay fixed; heuristics, targets and orderings vary per design."},
	{KP: 1, Q: "A customer wants CTS to target 30 ps of skew instead of 50 ps. Which tier is this?",
		A:   "Tier 1: a parameter change",
		D:   []string{"Tier 2: a new plug-in", "Tier 3: a core edit", "It cannot be expressed"},
		Why: "The skew target is an existing knob; changing it needs no new code."},
	{KP: 1, Q: "A team wants sinks clustered by timing criticality as well as by location, which the engine's clustering interface allows but no library policy implements. Which tier?",
		A:   "Tier 2: a new clustering policy plug-in at the clustering extension point",
		D:   []string{"Tier 1: change the cluster size", "Tier 3: rewrite the tree data structure", "No change; clustering cannot use timing"},
		Why: "New behaviour at a defined seam is a plug-in; the core that builds and legalises the tree is untouched."},
	{KP: 1, Q: "Which change would be a Tier 3 core edit to a CTS engine?",
		A:   "Changing how the engine computes stage delay and slew for every buffer it inserts",
		D:   []string{"Selecting a different buffer list", "Adding a topology policy at the extension point", "Raising the maximum fanout"},
		Why: "The delay and slew evaluation underpins every decision and every result; changing it needs human review and full regression."},
	{KP: 2, Q: "What context must a CTS topology policy plug-in receive at its seam?",
		A:   "Sink locations and loads, clock-domain and gating structure, blockages, and the current skew and latency targets",
		D:   []string{"Only the number of sinks", "The RTL source", "The router's internal grid"},
		Why: "A policy decides structure; it needs the geometry, the loads and the targets, not the internals of other engines."},
	{KP: 2, Q: "Why should a CTS policy plug-in be called once per clock net or cluster rather than once per buffer evaluation?",
		A:   "Coarse calls keep the hot delay kernel inside the compiled core, where it can be optimised",
		D:   []string{"Plug-ins cannot be called more than once", "Per-buffer calls give different answers", "The clock net has only one buffer"},
		Why: "The seam rule: plug-ins make coarse decisions, and the inner loops stay in the core."},
	{KP: 2, Q: "A customer requests that all ICGs be cloned when their enable slack is below 20 ps. The engine has a gating-policy seam. Which tier?",
		A:   "Tier 2: a gating policy plug-in that reads enable slack and requests clones through the core's API",
		D:   []string{"Tier 1: a skew parameter", "Tier 3: edit the timer core", "It must be done by hand after CTS"},
		Why: "It is a decision rule over available data, implemented at the policy seam; cloning itself is a core service."},
	{KP: 2, Q: "A new CTS policy plug-in passes its own unit tests. What must happen before it is promoted to the shared library?",
		A:   "It must pass the engine's golden CTS tests and show a skew, latency, power and wirelength comparison against the baseline",
		D:   []string{"Nothing more", "Only a code style review", "It must be merged into the core"},
		Why: "Promotion needs proof that the policy is legal and improves QoR on reference designs without regressions."},
	{KP: 3, Q: "A plug-in improves skew on one design by moving a branch point closer to the root, but latency grows 15% on others. How should the result be reported?",
		A:   "As a trade-off with per-design deltas, so that selection can match the policy to designs where it helps",
		D:   []string{"As a universal improvement", "As a failure to be deleted", "Only the best design's result"},
		Why: "Policies are selected by manifest and context; honest per-design deltas let the selector choose well."},
	{KP: 3, Q: "What should a CTS policy's manifest state so that an agent can select it?",
		A:   "What kinds of designs it suits (sink distribution, gating depth, mesh or tree), its inputs and its tunable parameters",
		D:   []string{"Its full source code", "The author's name only", "The binary size"},
		Why: "Selection is a cheap text decision over manifests; it needs the policy's domain of applicability."},
	{KP: 3, Q: "A plug-in reads a global variable set by another policy. What rule does that break?",
		A:   "Plug-ins must have no hidden global state, so results are deterministic and composable",
		D:   []string{"Plug-ins must be written in assembly", "Plug-ins must never read data", "Policies must run on separate machines"},
		Why: "Hidden coupling makes results depend on call order and breaks reproducibility."},
	{KP: 3, Q: "A differential CTS check after a policy change shows identical skew and latency but 3% more clock buffer area. What is the verdict?",
		A:   "A regression in power and area that must be justified before promotion",
		D:   []string{"A pass, since timing is unchanged", "An improvement", "Irrelevant to CTS"},
		Why: "The QoR baseline covers power and area too; identical timing with more buffers is worse."},
}

var rewardItems = []edakit.Item{
	{KP: 1, Q: "An RL agent edits the PDN and is rewarded only for lower worst static IR drop. What will it learn to do?",
		A:   "Add as much metal as it can, ignoring the routing tracks and congestion it consumes",
		D:   []string{"Minimise wirelength", "Balance IR against routability automatically", "Lower the supply voltage"},
		Why: "A reward that sees one metric is maximised at the expense of everything it does not see."},
	{KP: 1, Q: "Which reward better reflects PDN quality?",
		A:   "Meeting static and dynamic IR and EM limits, with penalties for routing resource and area used",
		D:   []string{"Total metal area alone", "The number of ECOs applied", "The solver's runtime"},
		Why: "PDN design is constrained optimisation: limits must hold, and resources must be spent sparingly."},
	{KP: 1, Q: "What makes a reward tool-grounded for a CTS agent?",
		A:   "Skew, latency, slew and power are computed by the timing and power engines on the agent's tree",
		D:   []string{"A human rates each tree by looking at it", "The agent estimates its own skew", "The reward is the number of buffers"},
		Why: "Deterministic engines score the result against objectives and constraints, which is the reward signal."},
	{KP: 1, Q: "A CTS agent's reward is negative global skew. What can go wrong?",
		A:   "It may pad every branch with delay to equalise latency, wasting power and adding variation",
		D:   []string{"It will always minimise power", "Skew cannot be computed", "It will refuse to insert buffers"},
		Why: "Global skew ignores latency and power; the cheapest way to equalise is often to slow everything down."},
	{KP: 2, Q: "The fast IR estimator used in the reward underestimates drop near macros. What will an RL agent learn?",
		A:   "To concentrate current near macros, where the reward cannot see the problem",
		D:   []string{"To avoid macros entirely", "Nothing different", "To improve the estimator"},
		Why: "Reward hacking exploits proxy errors; the gap between proxy and signoff is where the policy goes."},
	{KP: 2, Q: "How should the training loop guard against the agent exploiting the fast proxy?",
		A:   "Periodically score samples with the signoff solver and penalise or retrain where proxy and signoff disagree",
		D:   []string{"Use only the proxy, more often", "Hide signoff results from the developers", "Add random noise to the proxy"},
		Why: "Auditing against the accurate engine detects and closes the gaps the agent is exploiting."},
	{KP: 2, Q: "A useful-skew agent is rewarded for total negative slack improvement. Which constraint must the reward also enforce?",
		A:   "Hold slack on every path, since skew that helps setup can break hold",
		D:   []string{"The die size", "The number of clock domains", "The RTL line count"},
		Why: "Setup and hold pull skew in opposite directions; rewarding one alone invites violations of the other."},
	{KP: 2, Q: "Why are EM checks part of the reward for a PDN agent and not only an after-the-fact check?",
		A:   "Otherwise the agent learns layouts that pass IR by pushing current through narrow paths that fail EM",
		D:   []string{"EM is cheaper to compute than IR", "EM has no relation to the grid", "Signoff tools do not check EM"},
		Why: "Constraints left out of the reward are violated whenever that helps the reward."},
	{KP: 3, Q: "A reward combines worst IR drop and metal usage as a weighted sum. Designs with similar scores look very different. What does that suggest?",
		A:   "The weights trade one metric against another; a Pareto view or hard limits may express the intent better",
		D:   []string{"The solver is broken", "Weighted sums are always exact", "Metal usage does not matter"},
		Why: "A weighted sum picks one point on a trade-off curve; hard limits keep each metric within spec."},
	{KP: 3, Q: "The reward for a CTS agent is computed on one corner only. What risk does that create?",
		A:   "Trees tuned for that corner can have much worse skew at other corners where buffer and wire delays scale differently",
		D:   []string{"None: skew is the same at every corner", "The agent cannot build trees", "Power is no longer defined"},
		Why: "Buffer delay and wire delay scale differently across corners, so balance at one corner does not imply balance at others."},
	{KP: 3, Q: "Where should the reward's IR and EM limits come from?",
		A:   "From the signoff specification and the foundry's EM rules, read by the tools, not chosen by the agent",
		D:   []string{"From the agent's own estimate of a reasonable limit", "From the previous chip's datasheet", "From the fastest run"},
		Why: "The limits are part of the oracle; an agent must not set the thresholds it is judged by."},
	{KP: 3, Q: "An agent's PDN reward rises during training while signoff IR results on held-out designs stay flat. What is happening?",
		A:   "The agent is overfitting the reward or the training designs, not learning better PDNs",
		D:   []string{"Signoff is wrong", "Training is converging correctly", "Held-out designs cannot be evaluated"},
		Why: "Held-out signoff QoR is the real measure; reward gains that do not transfer are proxy exploitation or overfitting."},
}
