package numer

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var kernelCoreItems = []edakit.Item{
	{KP: 1, Q: "In a circuit simulator built on the four-layer engine model, which of these belongs in the immutable core?",
		A: "The sparse LU refactorization and triangular solves", D: []string{"The choice between AMD and nested dissection ordering", "The reltol value used in Newton convergence", "A hook that logs Newton residuals for an agent"},
		Why: "The numerical kernel that every analysis relies on is core; orderings are policy, tolerances are parameters, logging is an observation hook."},
	{KP: 1, Q: "Which property of the solver core must every change preserve, because downstream engines and golden tests depend on it?",
		A: "The meaning of the residual and the convergence test, so a reported tolerance always means the same thing", D: []string{"The exact number of Newton iterations on every design", "The wall-clock time of each solve", "The ordering chosen for every matrix"},
		Why: "Iteration counts and orderings may change; the contract that a converged solution satisfies a stated residual tolerance may not."},
	{KP: 1, Q: "An engineer proposes to change the local truncation error estimator of the transient engine. Where does that change live?",
		A: "In the core: it changes the accuracy contract of every transient result", D: []string{"In the parameters, as a new reltol", "In an observation hook", "In the strategy layer, as a new ordering"},
		Why: "The LTE formula decides which steps are accepted, so it defines result accuracy; editing it is a core change needing approval and full regression."},
	{KP: 1, Q: "Why keep the hot numerical kernels (SpMV, triangular solve, element evaluation) inside the core rather than behind fine-grained plug-in calls?",
		A: "Per-entry calls across a plug-in boundary block inlining and vectorization and would dominate run time", D: []string{"Plug-ins cannot perform floating-point arithmetic", "The kernels change too often to be core", "Agents cannot call core code"},
		Why: "Seams should be coarse, per solve or per batch, so the compiler can optimize the inner loops."},
	{KP: 2, Q: "A team wants transient runs of their PLL to use BDF2 instead of trapezoidal integration. Which layer does this land in?",
		A: "Strategy or policy: choosing among integration methods the core already provides", D: []string{"The immutable core", "An observation hook", "The netlist parser"},
		Why: "Selecting an existing integration method is a policy choice; only writing a new formula would touch the core."},
	{KP: 2, Q: "A request reads: tighten reltol from 1e-3 to 1e-4 for a precision ADC block. Which layer handles it?",
		A: "Parameters: a configuration value with no recompile", D: []string{"The immutable core", "A new plug-in", "A new observation hook"},
		Why: "Tolerances are parameters; changing them changes run time and accuracy without touching code."},
	{KP: 2, Q: "An agent should see, after every Newton iteration, the residual norm and the node with the largest update. What provides that?",
		A: "An observation hook in the Newton loop", D: []string{"A change to the LU kernel", "A new parameter value", "A new integration method"},
		Why: "Hooks expose live solver state to agents without changing what the solver computes."},
	{KP: 2, Q: "A damping strategy that limits per-device junction voltage changes is to be swapped for a global line search. Which layer is that?",
		A: "Strategy: both are policies around the same Newton update", D: []string{"Parameters", "The immutable core", "An observation hook"},
		Why: "How to scale or limit the Newton step is a policy decision; the Newton update and Jacobian assembly stay in the core."},
	{KP: 3, Q: "A solver produces slightly different IR-drop results on two runs with the same input on the same machine. Which core property is broken?",
		A: "Determinism: fixed inputs and settings must give bit-identical results", D: []string{"Accuracy: the results must match signoff", "Stability: the matrix must be positive definite", "Sparsity: the fill must be minimal"},
		Why: "Without determinism, golden and differential tests cannot tell a real change from noise; nondeterministic parallel reductions are the usual cause."},
	{KP: 3, Q: "Where should a seedable random choice inside a solver (say, a random restart vector for an eigensolver) get its seed?",
		A: "From an explicit parameter recorded with the run, so the result can be reproduced", D: []string{"From the system clock", "From the thread scheduler", "From uninitialized memory"},
		Why: "Everything stochastic must be seeded explicitly so that runs reproduce and diffs are meaningful."},
	{KP: 3, Q: "A GPU solver offers a deterministic mode and a faster mode with atomic-add reductions. When should an agent loop use the deterministic mode?",
		A: "Whenever results are compared against a baseline or golden reference", D: []string{"Never, since the fast mode is always within tolerance", "Only for designs under 1000 nodes", "Only when the matrix is unsymmetric"},
		Why: "Comparisons need reproducible numbers; the nondeterministic fast mode can be used for exploration where only rough rankings matter."},
	{KP: 3, Q: "A new release reorders the floating-point operations of the triangular solve for speed. Results change in the 14th digit. What must the release process show?",
		A: "That the backward error and all golden waveforms stay within their stated tolerances, and that the new results are themselves deterministic", D: []string{"That results are bit-identical to the old release", "Nothing, since a 14th-digit change is invisible", "That the change was made by an agent"},
		Why: "Bitwise equality across a reordering is impossible; the contract is tolerance-based agreement plus determinism within the new version."},
}

var solverSeamItems = []edakit.Item{
	{KP: 1, Q: "Which interface makes a good seam for plugging new linear solvers into a circuit simulator?",
		A: "analyze(pattern), factor(values), solve(rhs): coarse calls per matrix and per solve", D: []string{"get(i, j) and set(i, j) called for each matrix entry", "One function that parses the netlist and returns waveforms", "A global variable holding the current matrix"},
		Why: "Splitting symbolic analysis from numeric factorization lets the simulator reuse the analysis across Newton iterations; coarse calls keep the hot loops inside the solver."},
	{KP: 1, Q: "Why should the linear-solver seam separate the symbolic analysis call from the numeric factorization call?",
		A: "The pattern is fixed across Newton iterations and timesteps, so the ordering is computed once and the numeric factorization repeated", D: []string{"Symbolic analysis is the cheapest step", "Numeric factorization does not need the pattern", "Some solvers have no symbolic phase, so it must be hidden"},
		Why: "Separating the phases exposes the reuse that makes sparse direct solvers fast in simulation."},
	{KP: 1, Q: "A preconditioner seam offers setup(A) and apply(r). What does apply(r) return?",
		A: "An approximation of A^-1 r, used once per Krylov iteration", D: []string{"The exact solution of Ax = r", "The residual r - Ax", "The condition number of A"},
		Why: "A preconditioner applies an approximate inverse; the Krylov method stays responsible for convergence."},
	{KP: 1, Q: "A plug-in preconditioner keeps a hidden global cache between calls. Why is that a seam violation?",
		A: "Results then depend on call history, which breaks determinism and makes golden tests unreliable", D: []string{"Global variables are slower to read", "The compiler cannot compile global variables", "Caches always increase memory use"},
		Why: "Seams should pass plain data and keep no hidden state, so the same inputs always give the same outputs."},
	{KP: 2, Q: "Which context should the solver seam expose so a plug-in can choose between Cholesky and LU?",
		A: "Whether the matrix is symmetric and positive definite (for example, a pure RC or power-grid matrix)", D: []string{"The name of the designer", "The output waveform file path", "The number of netlist comments"},
		Why: "Structure hints decide the algorithm: SPD matrices allow Cholesky or CG; general MNA matrices need LU or GMRES."},
	{KP: 2, Q: "A preconditioner plug-in wants to reuse its setup when the matrix values change only slightly between timesteps. What must the seam tell it?",
		A: "Whether the pattern is unchanged and how much the values changed since the last setup", D: []string{"The total simulation time requested", "Which agent requested the run", "The number of output nodes"},
		Why: "Reuse decisions need to know what changed; exposing a change indicator lets the plug-in skip expensive setups safely."},
	{KP: 2, Q: "An iterative-solver plug-in is called inside the Newton loop. Which piece of context lets it pick a sensible tolerance?",
		A: "The outer Newton tolerance and the current nonlinear residual, so the inner solve is only as accurate as needed", D: []string{"The circuit's clock frequency", "The number of devices in the netlist", "The identity of the previous plug-in"},
		Why: "Inexact Newton methods tie the linear tolerance to the nonlinear residual; solving more accurately than that wastes work."},
	{KP: 2, Q: "A power-grid preconditioner wants to coarsen by metal layer. What must the seam provide beyond the matrix?",
		A: "A map from unknowns to grid layers and coordinates", D: []string{"The SPICE netlist text", "The timing constraints", "The die cost"},
		Why: "Geometric or structural metadata lets a plug-in build better hierarchies than it could infer from matrix entries alone."},
	{KP: 3, Q: "A customer asks to switch their IR-drop runs from Jacobi to incomplete-Cholesky preconditioning, which the tool already ships. Which change tier is this?",
		A: "Tier 1: a configuration choice", D: []string{"Tier 2: a plug-in", "Tier 3: a core edit", "No tier: it cannot be done"},
		Why: "Selecting an existing preconditioner is a parameter or policy setting, with no new code."},
	{KP: 3, Q: "A team wants a block-Jacobi preconditioner whose blocks are the power-grid layers. It does not exist yet. Which tier?",
		A: "Tier 2: a new plug-in behind the preconditioner seam", D: []string{"Tier 1: a configuration value", "Tier 3: a core edit", "No change: Jacobi already covers it"},
		Why: "It is new code, but it fits the existing setup/apply seam, so the core stays untouched."},
	{KP: 3, Q: "A request asks that Newton convergence be declared on the residual norm alone, dropping the update test. Which tier?",
		A: "Tier 3: it changes the convergence semantics in the core and needs human approval and full regression", D: []string{"Tier 1: a tolerance value", "Tier 2: a damping plug-in", "Tier 1: an observation hook"},
		Why: "The meaning of convergence is a core contract; every downstream result depends on it."},
	{KP: 3, Q: "Over a year, most solver change requests have needed Tier 3 core edits. What does that indicate?",
		A: "The seams are in the wrong place: behaviour that varies is trapped inside the core", D: []string{"The core is well designed", "Customers should stop requesting changes", "Plug-ins are too powerful"},
		Why: "A good seam makes the common change a configuration or plug-in; frequent core edits signal that the variable part is not exposed."},
}

var diffTestItems = []edakit.Item{
	{KP: 3, Q: "A golden test suite for a transient engine includes an RC step response. What should the test compare?",
		A: "The waveform against the analytic exponential within stated reltol and abstol bands", D: []string{"Bit-identical output against the previous release", "Only the final value at the end time", "The number of lines in the output file"},
		Why: "Analytic references check accuracy directly, and tolerance bands allow legitimate changes in step selection."},
	{KP: 3, Q: "After a new ordering plug-in, IR-drop results match the baseline to 1e-12 V, but Newton iteration counts on a nonlinear suite doubled. How should the differential report treat it?",
		A: "As a QoR regression to review: results agree but cost changed, so counts belong in the report", D: []string{"As a pass with nothing to report", "As a correctness failure", "As proof the plug-in is better"},
		Why: "Differential checks cover accuracy and cost; iteration and timestep counts are QoR metrics that reveal performance regressions."},
	{KP: 3, Q: "Why should a solver golden suite contain deliberately ill-conditioned matrices?",
		A: "A change that is harmless on well-conditioned cases can lose many digits when kappa is large", D: []string{"Ill-conditioned matrices run faster", "They make every solver agree exactly", "They are needed to test the netlist parser"},
		Why: "Stability problems show up only under ill-conditioning; well-conditioned tests hide them."},
	{KP: 3, Q: "A differential test compares two transient engines by their waveforms at every timestep of the run. Why does that fail spuriously?",
		A: "The engines choose different timesteps, so waveforms must be compared at common times by interpolation or at fixed sample points", D: []string{"Waveforms cannot be compared at all", "Transient engines always pick the same timesteps", "The comparison must use bitwise equality"},
		Why: "Adaptive step control makes time grids differ; comparisons must align the grids before applying tolerances."},
	{KP: 3, Q: "A solver change passes all golden tests but was only run on small designs. What must still be shown before promotion?",
		A: "Results and run time on representative large designs from the QoR baseline", D: []string{"Nothing more, because golden tests passed", "That the change compiles on one machine", "That the change is short"},
		Why: "Scaling behaviour (fill, iteration counts, memory) appears on large designs; promotion needs the baseline comparison too."},
	{KP: 3, Q: "Which IR-drop golden case gives the strongest single check of a power-grid solver?",
		A: "A grid with a known analytic or independently verified worst-node voltage drop", D: []string{"A grid with no load currents", "A grid with one node", "Any grid, compared only to itself"},
		Why: "An independently known answer checks accuracy; a trivial grid or a self-comparison checks nothing."},
}

func backwardErrorChoice(ctx generator.GeneratorContext) generator.Problem {
	a, _ := spd2()
	x := []float64{float64(edakit.Int(-5, 5)), float64(edakit.Int(1, 5))}
	b := matvec(a, x)
	pert := float64(edakit.Int(1, 9)) * math.Pow(10, -float64(edakit.Int(3, 6)))
	xh := []float64{x[0] + pert, x[1] - pert}
	axh := matvec(a, xh)
	r := []float64{b[0] - axh[0], b[1] - axh[1]}
	inf := func(v []float64) float64 { return math.Max(math.Abs(v[0]), math.Abs(v[1])) }
	nA := math.Max(math.Abs(a[0][0])+math.Abs(a[0][1]), math.Abs(a[1][0])+math.Abs(a[1][1]))
	eta := inf(r) / (nA*inf(xh) + inf(b))
	sciS := func(v float64) string {
		if v == 0 {
			return "$0$"
		}
		e := int(math.Floor(math.Log10(v)))
		return sp(`$%s \times 10^{%d}$`, f(v/math.Pow(10, float64(e)), 2), e)
	}
	q := "A solver change is tested on $A = $ " + pmat(a, 0) + ", $b = $ " + pvec(b, 0) + sp(". It returns $\\hat x = (%s, %s)$. What is the normwise backward error $\\eta = \\|b - A\\hat x\\|_\\infty/(\\|A\\|_\\infty\\|\\hat x\\|_\\infty + \\|b\\|_\\infty)$?", f(xh[0], 7), f(xh[1], 7))
	return edakit.Choice(q, sciS(eta), sp("$\\|r\\|_\\infty = %s$, $\\|A\\|_\\infty = %s$, $\\|\\hat x\\|_\\infty = %s$, $\\|b\\|_\\infty = %s$. A stable solver in double precision gives $\\eta$ near $10^{-16}$; this one is far above, so the test fails.", f(inf(r), 9), f(nA, 0), f(inf(xh), 6), f(inf(b), 0)),
		sciS(inf(r)/inf(b)), sciS(inf(r)), sciS(inf(r)/(nA*inf(xh))), sciS(eta*10))
}

func toleranceChoice(ctx generator.GeneratorContext) generator.Problem {
	ke := edakit.Int(3, 9)
	te := edakit.Int(8, 12)
	allowed := 2 * math.Pow(10, float64(ke-te))
	ok := edakit.Int(0, 1) == 0
	var diff float64
	if ok {
		diff = allowed * float64(edakit.Int(1, 8)) / 10
	} else {
		diff = allowed * float64(edakit.Int(20, 90))
	}
	ds := sciPlain(diff)
	q := sp("An old and a new linear solver each meet a relative residual tolerance of $10^{-%d}$ on a matrix with $\\kappa \\approx 10^{%d}$. Their solutions differ by a relative amount %s. What should the differential test conclude?", te, ke, ds)
	pass := sp("Consistent: forward errors of up to about $\\kappa \\cdot 10^{-%d}$ each are allowed, so a difference up to about %s is expected", te, sciPlain(allowed))
	fail := sp("Inconsistent: the difference exceeds the roughly %s that conditioning allows, so one solver is wrong", sciPlain(allowed))
	wrongPass := "Consistent: any difference is acceptable once both residuals pass"
	wrongFail := sp("Inconsistent: the solutions must agree to $10^{-%d}$, the residual tolerance", te)
	if ok {
		return edakit.Choice(q, pass, "Each solution's relative error is bounded by kappa times its relative residual, so the two may differ by about the sum of those bounds; this difference is within it.", fail, wrongFail, wrongPass)
	}
	return edakit.Choice(q, fail, "The bound kappa times the residual tolerance, summed over both solvers, is far smaller than the observed difference, so at least one solution is wrong.", pass, wrongFail, wrongPass)
}

func sciPlain(v float64) string {
	e := int(math.Floor(math.Log10(v)))
	return sp(`$%s \times 10^{%d}$`, f(v/math.Pow(10, float64(e)), 1), e)
}

var gpuLoopItems = []edakit.Item{
	{KP: 3, Q: "An agent's fast GPU IR-drop mode ranks two strap patterns 0.4 mV apart, and the fast mode's error band is about 1 mV. What should the agent do?",
		A: "Recurse to the accurate solver for both candidates before choosing", D: []string{"Pick the one the fast mode prefers", "Average the two designs", "Tighten the fast mode by running it twice"},
		Why: "The difference is inside the fast mode's error band, so the ranking is not trustworthy; the accurate engine must decide."},
	{KP: 3, Q: "The best fast-mode candidate's worst IR drop is 49 mV against a 50 mV limit, with a 2 mV fast-mode error band. What follows?",
		A: "It must be confirmed by the accurate solver, because the margin is inside the error band", D: []string{"It passes, since 49 is below 50", "It fails, since the error band is positive", "Nothing: the limit applies only at signoff"},
		Why: "Near a constraint limit a fast estimate cannot decide pass or fail; the accurate engine is required."},
	{KP: 3, Q: "Consecutive candidates in a decap-placement search differ in a few hundred capacitors. Which trick cuts GPU CG iterations most?",
		A: "Warm-start each solve from the previous candidate's solution", D: []string{"Start every solve from zero for fairness", "Switch to forward Euler", "Reduce the precision of the matrix indices"},
		Why: "A nearby solution makes the initial residual small, so CG needs far fewer iterations to reach the same tolerance."},
	{KP: 3, Q: "A fast mode was calibrated on 7 nm designs and the agent now explores a 3 nm design with different grid topology. What is the risk?",
		A: "The fast mode's error band may not hold on the new design class, so it must be re-measured against the accurate solver", D: []string{"None: numerical error does not depend on the design", "The GPU will refuse the new matrix", "The accurate solver will become slower"},
		Why: "Error estimates for surrogate or reduced modes are empirical; they transfer only to designs like the calibration set."},
	{KP: 3, Q: "Why is a GPU especially effective for the inner loop of a thermal floorplan agent that evaluates many candidate floorplans?",
		A: "Many independent solves of similar size can be batched to fill the device and amortize launches", D: []string{"GPUs make sparse direct factorization trivially parallel", "Thermal matrices are dense", "GPUs remove the need for a convergence test"},
		Why: "Batched, regular, bandwidth-bound work is what GPUs do well; single small solves leave most of the device idle."},
	{KP: 3, Q: "What is the main reason not to sign off a power grid on the fast GPU mode's numbers, even if they look comfortably within limits?",
		A: "Signoff requires the accurate engine with its validated tolerances; the fast mode's error is only estimated", D: []string{"GPU numbers are always too optimistic", "Signoff tools cannot read GPU output", "The fast mode uses a different netlist"},
		Why: "Fast modes guide exploration; the accuracy contract for signoff belongs to the accurate engine."},
}

func batchChoice(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Int(20, 500)
	s := float64(edakit.Int(5, 50)) // ms setup
	t := float64(edakit.Int(1, 20)) / 2
	un := float64(k) * (s + t)
	bt := s + float64(k)*t
	sp1 := un / bt
	ans := f(sp1, 2) + " times faster"
	q := sp("An agent evaluates %d candidate power-grid variants that share one sparsity pattern. Each solve on its own costs %s ms of setup (symbolic analysis, transfer, launch) plus %s ms of numeric work. Batched, the setup is paid once and the numeric work stays %s ms per candidate. What does batching gain?", k, f(s, 0), f(t, 1), f(t, 1))
	return edakit.Choice(q, ans, sp("Unbatched: %d(%s + %s) = %s ms. Batched: %s + %d · %s = %s ms. Ratio %s.", k, f(s, 0), f(t, 1), f(un, 0), f(s, 0), k, f(t, 1), f(bt, 1), f(sp1, 2)),
		f(s/t, 2)+" times faster", f(float64(k), 0)+" times faster", f((s+t)/s, 2)+" times faster", f(sp1*2, 2)+" times faster")
}

func budgetChoice(ctx generator.GeneratorContext) generator.Problem {
	T := float64(edakit.Int(10, 60)) * 60 // s
	F := float64(edakit.Int(10, 120))
	tf := float64(edakit.Int(5, 50)) / 10
	m := edakit.Int(2, 8)
	ta := float64(edakit.Int(30, 300))
	n := int(math.Floor((T - F - float64(m)*ta) / tf))
	if n < 1 {
		n = 1
	}
	q := sp("An agent has %s minutes to explore thermal floorplans. Loading the design costs %s s once. Each fast-mode candidate takes %s s, and the %d best candidates must then be re-run in the accurate solver at %s s each. How many fast-mode candidates fit in the budget?", f(T/60, 0), f(F, 0), f(tf, 1), m, f(ta, 0))
	return edakit.Choice(q, itoa(n), sp("$\\lfloor(%s - %s - %d \\cdot %s)/%s\\rfloor = %d$.", f(T, 0), f(F, 0), m, f(ta, 0), f(tf, 1), n),
		edakit.Distinct(n, 3, int(math.Floor((T-F)/tf)), int(math.Floor(T/tf)), int(math.Floor((T-float64(m)*ta)/tf)), n+m)...)
}

var fastAccurateItems = []edakit.Item{
	{KP: 3, Q: "An agent finds a PDN configuration using the fast mode (FP32, CG tolerance 1e-3, coarse mesh). What is the right promotion step?",
		A: "Re-run that configuration in the accurate mode (FP64, tight tolerance, full mesh) and promote only on its numbers", D: []string{"Promote it directly, since the fast mode found it", "Re-run it in the fast mode with a different seed", "Report the fast-mode result with a disclaimer"},
		Why: "Fast results guide the search; the decision and signoff rest on the accurate engine."},
	{KP: 3, Q: "Which fast-mode setting most directly risks a wrong ranking near a tight IR-drop limit?",
		A: "A loose iterative-solver tolerance whose error is comparable to the margin", D: []string{"A deterministic reduction order", "Warm starting from the previous solution", "Reusing the symbolic analysis"},
		Why: "Solver tolerance directly bounds the error of each estimate; the others save work without adding error."},
	{KP: 3, Q: "A transient fast mode loosens the LTE tolerance tenfold. Roughly how does the number of timesteps change for a trapezoidal integrator?",
		A: "It falls by about the cube root of 10, roughly a factor 2.2", D: []string{"It falls by 10", "It falls by 100", "It does not change"},
		Why: "Trapezoidal LTE scales as h^3, so the step grows like tol^(1/3) and the step count shrinks by about 10^(1/3)."},
	{KP: 3, Q: "A reduced-order thermal model is used as the fast mode. What must be checked before trusting its rankings on a new floorplan?",
		A: "That the reduction error on that floorplan, measured against the full model on a few cases, is below the ranking margins", D: []string{"That the reduced model has fewer states than the full one", "That the reduced model runs faster", "Nothing: reduced models are exact at DC"},
		Why: "A reduced model is accurate only for the inputs and geometries it was built for; its error must be measured where it is used."},
	{KP: 3, Q: "An agent's reward is computed from fast-mode IR drop. Over many iterations it learns layouts whose fast-mode drop is low but whose accurate drop is not. What is this?",
		A: "Reward hacking of the fast proxy; the reward must be periodically checked or computed with the accurate engine", D: []string{"Normal convergence of the agent", "A bug in the accurate engine", "Proof that the fast mode is more accurate"},
		Why: "Optimizing against a proxy finds the proxy's blind spots; grounding the reward in the accurate tool closes the gap."},
	{KP: 3, Q: "Which statistic lets an agent decide automatically when the fast mode is good enough to rank two candidates?",
		A: "An error bound or measured error band for the fast mode, compared with the difference between the candidates", D: []string{"The run time of the fast mode", "The number of candidates explored so far", "The sparsity of the matrix"},
		Why: "If the candidates differ by more than the combined error, the fast ranking is trustworthy; otherwise recurse."},
}

func toleranceIterChoice(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(10, 80) * 10
	e1 := edakit.Int(8, 12)
	e2 := edakit.Int(2, 5)
	est := float64(n) * math.Log(2*math.Pow(10, float64(e2))) / math.Log(2*math.Pow(10, float64(e1)))
	ans := f(est, 0)
	q := sp("The accurate mode of an IR-drop solver runs CG to a relative tolerance of $10^{-%d}$ in %d iterations. A fast mode stops at $10^{-%d}$. If iterations scale like $\\ln(2/\\varepsilon)$ at a fixed condition number, about how many iterations does the fast mode take?", e1, n, e2)
	return edakit.Choice(q, ans, sp("$%d \\cdot \\ln(2 \\cdot 10^{%d})/\\ln(2 \\cdot 10^{%d}) = %s$ iterations.", n, e2, e1, ans),
		distinctFloats(ans, 0, float64(n)*float64(e2)/float64(e1)*0.5, float64(n), float64(n)/2, float64(n)*math.Sqrt(float64(e2)/float64(e1)))...)
}

func rankingChoice(ctx generator.GeneratorContext) generator.Problem {
	d1 := float64(edakit.Int(200, 600)) / 10
	gap := float64(edakit.Int(1, 60)) / 10
	d2 := d1 + gap
	e := float64(edakit.Int(5, 30)) / 10
	trust := gap > 2*e
	q := sp("A fast-mode IR-drop solve gives worst drops of %s mV for candidate A and %s mV for candidate B. Each fast estimate is within $\\pm$%s mV of the accurate value. Can the agent trust that A is better than B?", f(d1, 1), f(d2, 1), f(e, 1))
	yes := sp("Yes: the gap of %s mV exceeds the combined error of %s mV", f(gap, 1), f(2*e, 1))
	no := sp("No: the gap of %s mV is within the combined error of %s mV, so recurse to the accurate solver", f(gap, 1), f(2*e, 1))
	wrongYes := sp("Yes: the gap of %s mV exceeds a single error band of %s mV", f(gap, 1), f(e, 1))
	wrongNo := "No: fast-mode results can never rank candidates"
	if trust {
		return edakit.Choice(q, yes, "Both estimates can move by the error band in opposite directions, so the ranking is safe only if the gap exceeds twice the band.", no, wrongNo,
			"No: the estimates must agree with the accurate solver to 0.1 mV before any ranking")
	}
	return edakit.Choice(q, no, "Both estimates can move by the error band in opposite directions, so a gap below twice the band does not settle the ranking.", yes, wrongYes, wrongNo)
}
