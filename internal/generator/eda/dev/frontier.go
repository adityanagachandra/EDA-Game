package dev

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var extractionSkillItems = []edakit.Item{
	{KP: 1, Q: "A customer sends new wafer measurements and asks for the BSIM-CMG card to be refitted. The model equations stay the same. Which change tier is this?",
		A: "Tier 1: new parameter values in the model card, with no recompile", D: []string{"Tier 2: a new model plug-in", "Tier 3: an edit to the simulator core", "No tier: refits are not engine changes"},
		Why: "Parameter fitting changes only values that the model already exposes. It is configuration, which is why it is the most common and safest agent task."},
	{KP: 1, Q: "A team wants a threshold shift that depends on a layout-proximity quantity the model does not have. Which tier is this?",
		A: "Tier 2: a small model extension at the device-model seam, leaving the solver untouched", D: []string{"Tier 1: a parameter value", "Tier 3: rewrite the Newton solver", "Tier 1: change the simulator's tolerance"},
		Why: "A new dependence needs new code, but it belongs in the model plug-in. The circuit-solver core does not change."},
	{KP: 1, Q: "A request asks to change how the circuit simulator limits Newton steps for all devices. Which tier is this?",
		A: "Tier 3: a core solver edit that needs human approval and full regression", D: []string{"Tier 1: a model parameter", "Tier 2: a model plug-in", "No change needed: limiting is per model"},
		Why: "Step limiting lives in the solver loop and affects every circuit. It is core, even if a few models also provide their own limiting hooks."},
	{KP: 1, Q: "An agent fitting parameters finds the fit only works if it edits a model equation. What should it do?",
		A: "Stop and escalate the change to the plug-in tier, with its own tests, rather than editing during a Tier 1 task", D: []string{"Edit the equation silently and continue", "Delete the measurements that do not fit", "Widen every parameter bound to infinity"},
		Why: "Tier discipline means a configuration task never quietly becomes a code change. The new equation goes through the plug-in path and its golden tests."},
	{KP: 2, Q: "What should the manifest of a compact-model extraction skill declare?",
		A: "Supported model and version, required measurements, which parameters it fits, their bounds, and the fit-quality report it returns", D: []string{"The full source code of the model", "The customer's design netlist", "Only the skill's name"},
		Why: "An agent selects skills by reading manifests. Model version, inputs, tunables and outputs are what it needs to match a request to a skill."},
	{KP: 2, Q: "Two extraction skills exist: one for BSIM4 planar devices, one for BSIM-CMG FinFETs. How should an agent choose?",
		A: "Match the device type and model version in the request against what each manifest advertises", D: []string{"Run both and keep the lower error regardless of model", "Read both skills' source code in full", "Pick the most recently written skill"},
		Why: "Selection is a cheap text decision over manifests. Running a planar extraction on FinFET data would fit the wrong equations."},
	{KP: 2, Q: "Why does a good extraction skill fit parameters in stages, such as threshold first, then mobility, then short-channel terms?",
		A: "Each stage uses the bias region where that parameter dominates, which keeps the fit well conditioned and physical", D: []string{"Staging is slower and has no benefit", "Global fitting is impossible", "Staging avoids the need for measurements"},
		Why: "Fitting all parameters at once over all data invites compensating errors. Staged extraction pins each physical effect where it is visible, then a final global polish refines."},
	{KP: 2, Q: "An extraction skill exposes parameter bounds as a tunable. Why?",
		A: "Bounds keep fitted values physical, and different processes need different ranges", D: []string{"Bounds make the fit exact", "Bounds are needed only for display", "The simulator ignores out-of-bound values"},
		Why: "Without bounds an optimizer can drive a parameter to a nonphysical value that fits the data but extrapolates badly."},
	{KP: 3, Q: "What golden test should an extraction skill ship with?",
		A: "A reference dataset whose fit must reproduce a known-good card and error report within tolerance", D: []string{"A test that the skill's file names are sorted", "A check that the fit runs faster than one second", "A screenshot of an I-V plot"},
		Why: "Golden inputs and outputs let the harness confirm that an edit to the skill did not change its results."},
	{KP: 3, Q: "A Tier 1 refit lowers RMS error but the agent still needs one more check before promotion. Which?",
		A: "A differential check that circuit-level metrics, such as ring-oscillator delay and leakage, moved only as expected", D: []string{"None: lower RMS error is always enough", "A check that the card file grew", "A re-run of synthesis"},
		Why: "A better device fit can still shift circuit behaviour in an unintended region. Comparing benchmark circuits against the previous card exposes that."},
	{KP: 3, Q: "Measured wafer data is customer IP. Where should the extraction skill run?",
		A: "On-site, inside the sandbox, with only the reusable skill improvements leaving", D: []string{"On a public server for speed", "In the agent's chat transcript", "Anywhere, since measurements are not IP"},
		Why: "The IP boundary is architectural. Know-how about how to fit can flow back; the customer's data does not."},
	{KP: 3, Q: "An agent discovers a new extraction sequence that fits a tricky process much better. How does it become reusable?",
		A: "Iterate as a plug-in, validate with golden and differential tests, then promote it into the skill library", D: []string{"Copy it into every customer's model card", "Hard-code it into the simulator core", "Keep it only in the agent's memory"},
		Why: "The skill flywheel: identify, iterate, validate, promote. Promotion makes it selectable by manifest for future requests."},
}

var modelSeamItems = []edakit.Item{
	{KP: 1, Q: "A circuit simulator treats device models as plug-ins. What should cross the seam on each call?",
		A: "Terminal voltages, temperature and instance parameters in; currents, charges and their derivatives out", D: []string{"The whole circuit matrix in both directions", "Only the drain current out", "The layout geometry in, a timing report out"},
		Why: "The solver assembles the matrix. The model answers one question: what are this device's currents, charges and Jacobian entries at this bias."},
	{KP: 1, Q: "Which part of a circuit simulator belongs in the immutable core rather than the model plug-in?",
		A: "Matrix assembly, the Newton loop and time-step control", D: []string{"The BSIM-CMG equations", "The mobility model of one process", "A customer's self-heating extension"},
		Why: "Solver mathematics is shared by every circuit. Device physics varies and belongs behind the seam."},
	{KP: 1, Q: "Which is an existing example of a compiled-model seam in open-source simulation?",
		A: "The OSDI interface that ngspice uses to load compiled Verilog-A models", D: []string{"The Liberty file format", "The DEF placement format", "The SDC constraint language"},
		Why: "OSDI defines a binary interface through which a compiled device model provides evaluations to the simulator."},
	{KP: 1, Q: "A plug-in must add an internal thermal node for self-heating. What must the seam allow?",
		A: "Models declaring internal nodes and their equations, so the solver adds unknowns", D: []string{"Nothing: models cannot have internal state", "A second simulator process", "Editing the core matrix code"},
		Why: "Real models have internal nodes for series resistance, body or temperature. A seam that only passes terminal values would force self-heating into the core."},
	{KP: 2, Q: "A plug-in computes derivatives by finite differences instead of analytically. What risk does that add?",
		A: "Inaccurate or inconsistent Jacobians that slow or break Newton convergence", D: []string{"None: finite differences are exact", "The model will run faster", "The currents become wrong at every bias"},
		Why: "Finite-difference derivatives carry truncation and roundoff error. Verilog-A compilers generate analytic derivatives for this reason."},
	{KP: 2, Q: "Why should the model seam evaluate devices in batches rather than one scalar call per device?",
		A: "Coarse calls let the model vectorize over many instances and keep call overhead small", D: []string{"Batches change the physics", "Single calls are not deterministic", "The solver cannot handle more than one device"},
		Why: "Model evaluation is hot. A per-device virtual call with scattered data defeats vectorization; batching by model type restores it."},
	{KP: 2, Q: "Parallel model evaluation gives different results run to run. What seam rule was broken?",
		A: "Plug-ins must be deterministic with no hidden global state", D: []string{"Plug-ins must be written in C", "Plug-ins must run on one thread", "Plug-ins must not return charges"},
		Why: "Hidden shared state makes results depend on thread scheduling. A clean seam passes plain data and keeps evaluations independent."},
	{KP: 2, Q: "A new model plug-in passes its own unit tests. What must the harness also run before accepting it?",
		A: "Golden I-V and C-V sweeps and benchmark circuits through the real simulator", D: []string{"Nothing more", "A synthesis run", "A layout DRC check"},
		Why: "Unit tests check the code. Golden sweeps through the solver check that the plug-in behaves inside Newton, across bias and temperature."},
	{KP: 3, Q: "A plug-in returns NaN at $V_{GS} = -2$ V, a bias no test covered. How should the harness catch this class of failure?",
		A: "Sweep stress biases beyond the normal range and check for finite, bounded outputs and derivatives", D: []string{"Only test at nominal bias", "Trust the plug-in's documentation", "Catch it in silicon"},
		Why: "Newton iterates visit wild voltages before converging. Robustness sweeps catch overflow, division by zero and log of negative numbers."},
	{KP: 3, Q: "The seam's binary interface changes in a new simulator release. What happens to existing plug-ins?",
		A: "They may break, which is why the extension ABI should be frozen and versioned", D: []string{"They speed up automatically", "Nothing, since ABIs never matter", "They are converted into core code"},
		Why: "A frozen, versioned ABI lets plug-ins built against it keep working. Changing it is a core change with migration cost."},
	{KP: 3, Q: "A team wants to try a neural-network device model. Where should it plug in?",
		A: "At the same device-model seam, providing currents, charges and smooth derivatives like any compact model", D: []string{"Inside the Newton solver", "In the Liberty files only", "As a replacement for the matrix solver"},
		Why: "The seam does not care how the model computes. A learned model must still meet the same contract, including smoothness and charge conservation."},
	{KP: 3, Q: "Why should a model plug-in run inside a sandbox during agent-driven iteration?",
		A: "Untrusted generated code must not access the network or files beyond its inputs, and must respect CPU and memory limits", D: []string{"Sandboxes make models more accurate", "The simulator requires a sandbox to load any model", "Sandboxing is only for TCAD"},
		Why: "Agent-written code is iterated quickly. Resource limits and no network keep failures contained and protect the IP boundary."},
}

// goldenPoint is one bias point of a golden I-V table.
type goldenPoint struct {
	vgs, vds float64
	gold     float64
}

// goldenIV builds a golden-versus-candidate table with exactly one point outside tolerance.
func goldenIV(ctx generator.GeneratorContext) generator.Problem {
	tol := pickF(0.5, 1, 2)
	vdss := []float64{0.05, 0.4, 0.8}
	vgss := []float64{0.2, 0.4, 0.6, 0.8}
	var pts []goldenPoint
	for _, i := range edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})[:edakit.Int(4, 5)] {
		vg, vd := vgss[i%4], vdss[i/4]
		// A plausible current in uA: exponential below 0.35 V, roughly quadratic above.
		var cur float64
		if vg < 0.35 {
			cur = 30 * math.Pow(10, (vg-0.35)/0.075)
		} else {
			cur = 0.5 * 600 * (vg - 0.3) * (vg - 0.3)
		}
		cur *= (1 - math.Exp(-vd/0.05))
		pts = append(pts, goldenPoint{vg, vd, cur})
	}
	bad := edakit.Int(0, len(pts)-1)
	var rows []string
	var opts []string
	var answer string
	for i, p := range pts {
		var err float64
		if i == bad {
			err = tol * pickF(1.6, 2, 3, 4)
		} else {
			err = tol * pickF(0.1, 0.2, 0.4, 0.6)
		}
		if edakit.Int(0, 1) == 0 {
			err = -err
		}
		cand := p.gold * (1 + err/100)
		label := edakit.Sprintf("$V_{GS} = %s$ V, $V_{DS} = %s$ V", edakit.Num(p.vgs, 1), edakit.Num(p.vds, 2))
		rows = append(rows, edakit.Sprintf("%s: golden %s µA, candidate %s µA", label, edakit.Sprintf("%.4g", p.gold), edakit.Sprintf("%.4g", cand)))
		if i == bad {
			answer = label
		} else {
			opts = append(opts, label)
		}
	}
	q := edakit.Sprintf("A compact-model plug-in edit is checked against a golden I-V table with a relative tolerance of %s%%. The results are: %s. Which bias point fails the golden test?", edakit.Num(tol, 1), strings.Join(rows, "; "))
	return edakit.Choice(q, answer, "The failing point is the one whose relative error $|I_{cand} - I_{gold}|/I_{gold}$ exceeds the tolerance. Relative, not absolute, error is used because currents span decades: a large absolute error at high current can pass while a tiny absolute error in subthreshold fails.", opts...)
}

var goldenItems = []edakit.Item{
	{KP: 1, Q: "Which set of sweeps belongs in a compact model's golden test suite?",
		A: "$I_D$-$V_{GS}$ and $I_D$-$V_{DS}$ families, C-V curves, derivatives such as $g_m$ and $g_{ds}$, and several temperatures", D: []string{"One $I_D$ value at nominal bias", "Only the threshold voltage", "Only the leakage at the hot corner"},
		Why: "Golden tests must exercise every region and every output the simulator uses, including derivatives that Newton and analog analysis depend on."},
	{KP: 1, Q: "Why should golden I-V comparisons use relative error, or error in log current, rather than absolute error?",
		A: "Currents span many decades, so absolute error ignores the subthreshold region entirely", D: []string{"Absolute error cannot be computed", "Relative error is always smaller", "Simulators report only relative currents"},
		Why: "A 1 nA error is negligible at 1 mA and catastrophic at 10 pA. Log or relative error weights every region."},
	{KP: 1, Q: "A golden suite contains the Gummel symmetry test. What failure does it guard against?",
		A: "Asymmetric or non-smooth behaviour around $V_{DS} = 0$ that breaks analog and RF simulation", D: []string{"High leakage at the hot corner", "Wrong gate resistance", "Too many model parameters"},
		Why: "Symmetry and smooth higher derivatives at zero drain bias are required by passive mixers and switches. Many otherwise good fits fail it."},
	{KP: 1, Q: "Where should the golden reference values come from for a model plug-in edit?",
		A: "The previously validated model version, run through the same simulator on fixed inputs", D: []string{"The edited plug-in itself", "Random numbers seeded per run", "A different technology's model"},
		Why: "Golden data is a frozen, trusted baseline. Generating it from the code under test would prove nothing."},
	{KP: 3, Q: "A plug-in edit adds a new GIDL term. What should the differential check show?",
		A: "Current changes only at negative $V_{GS}$ with high $V_{DS}$, and all other regions match the baseline", D: []string{"Every bias point changes by the same percentage", "Nothing changes anywhere", "Only the threshold voltage changes"},
		Why: "A differential test proves the edit moved only what it should. GIDL lives in off-state, high drain bias conditions."},
	{KP: 3, Q: "A model change passes every golden I-V point but a ring-oscillator benchmark's frequency shifts by 3%. What was probably missed?",
		A: "The charge or capacitance model changed, which I-V points alone do not exercise", D: []string{"Nothing: ring oscillators are unreliable", "The supply voltage of the benchmark", "The number of fins in the layout"},
		Why: "Delay depends on charge as well as current. Golden suites need C-V points and benchmark circuits, not just DC currents."},
	{KP: 3, Q: "Why is a tolerance of exactly zero a poor choice for golden model tests?",
		A: "Compiler and platform differences change the last bits of floating point, so tests fail without any real change", D: []string{"Zero tolerance passes everything", "Golden tests do not compare numbers", "Simulators always round to integers"},
		Why: "Tolerances should be tight enough to catch real changes and loose enough to absorb roundoff, typically relative tolerances well below fitting error."},
	{KP: 3, Q: "A golden suite covers only nominal temperature. An agent edit changes the temperature coefficient of mobility. What happens?",
		A: "The edit passes golden tests but corrupts hot and cold corners, which the suite never checks", D: []string{"The suite catches it automatically", "Temperature coefficients cannot change", "Only leakage is affected, which tests ignore anyway"},
		Why: "A test suite certifies only what it exercises. Corner temperatures belong in the golden set."},
}

// fitRMS computes the RMS log10 error of a candidate fit, a common extraction reward.
func fitRMS(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(3, 5)
	var ratios []float64
	var parts []string
	for i := 0; i < n; i++ {
		r := pickF(0.5, 0.8, 0.9, 1.1, 1.25, 1.5, 2, 3)
		ratios = append(ratios, r)
		parts = append(parts, edakit.Num(r, 2))
	}
	var ss, sabs, srel, mx float64
	for _, r := range ratios {
		e := math.Log10(r)
		ss += e * e
		sabs += math.Abs(e)
		srel += (r - 1) * (r - 1)
		if math.Abs(e) > mx {
			mx = math.Abs(e)
		}
	}
	rms := math.Sqrt(ss / float64(n))
	q := edakit.Sprintf("A fitting reward scores a compact-model candidate by the RMS of $\\log_{10}(I_{model}/I_{meas})$ over the measured points. The ratios $I_{model}/I_{meas}$ at %d points are %s. What is the RMS log error?", n, strings.Join(parts, ", "))
	return edakit.Choice(q, fixed(rms, 3),
		edakit.Sprintf("Take $\\log_{10}$ of each ratio, square, average over %d points and take the square root: %s. Log error treats a factor of 2 high and a factor of 2 low equally, at any current level.", n, fixed(rms, 3)),
		fixed(sabs/float64(n), 3), fixed(math.Sqrt(srel/float64(n)), 3), fixed(mx, 3), fixed(ss/float64(n), 3))
}

var fitRewardItems = []edakit.Item{
	{KP: 1, Q: "A tool-grounded reward for parameter extraction is built from simulator output. What should it measure?",
		A: "Error between simulated and measured curves across all regions, plus penalties for nonphysical parameters and broken smoothness", D: []string{"Only the number of parameters changed", "The agent's confidence in its answer", "The run time of the extraction only"},
		Why: "The simulator is the ground truth for what a card does. The reward must also encode the constraints a fit must satisfy, not just the error."},
	{KP: 1, Q: "Why should an extraction reward include errors in $g_m$ and $g_{ds}$, not only in $I_D$?",
		A: "A fit can match currents while getting their slopes wrong, which ruins gain, delay and Newton behaviour", D: []string{"Derivatives are always correct when currents are", "Derivatives are not used by simulators", "Including them makes the reward easier to game"},
		Why: "Analog gain depends on $g_m/g_{ds}$, and the solver needs consistent derivatives. Currents alone under-constrain the fit."},
	{KP: 1, Q: "Measured data comes from five device geometries. How should the reward be validated?",
		A: "Hold out some geometries, fit on the rest and score the held-out ones, so scalability is tested", D: []string{"Score only the devices used for fitting", "Use a different technology for validation", "Skip validation when training error is low"},
		Why: "A card must predict every legal W and L. Holding out geometries tests generalization, the way held-out data does for any model."},
	{KP: 1, Q: "Why are fast proxies dangerous as rewards in extraction?",
		A: "A proxy, such as error at a few bias points, can be optimized while the full-curve or circuit behaviour gets worse", D: []string{"Fast proxies are always exact", "Proxies make the simulator slower", "Rewards cannot be computed quickly"},
		Why: "This is reward hacking: the agent improves the measured number, not the goal. Periodic scoring with the full evaluation closes the gap."},
	{KP: 3, Q: "An agent's fit scores a near-perfect reward but one mobility parameter sits at its lower bound. What should a reviewer suspect?",
		A: "The parameter is compensating for a missing effect, and the card may extrapolate badly", D: []string{"The fit is certainly correct", "The bound has no effect on the fit", "The measurement data is perfect"},
		Why: "Parameters pinned at bounds often signal that the model is absorbing an effect into the wrong knob. Check physical plausibility and held-out behaviour."},
	{KP: 3, Q: "A reward uses only data up to $V_{DS} = 0.8$ V. The agent finds a card that fits perfectly but whose current falls with $V_{DS}$ above 0.9 V. Why was this not penalized?",
		A: "The reward never evaluated that region, so a non-monotone extrapolation was free", D: []string{"Currents always fall above 0.9 V", "The simulator rejected the card", "The reward penalized it automatically"},
		Why: "Add physical sanity checks over the full qualified range, such as monotonicity and positive output conductance, to the reward."},
	{KP: 3, Q: "An agent learns to lower the extraction error by smoothing the measured data before fitting. What is wrong?",
		A: "It is altering the ground truth instead of improving the model, which the reward must forbid", D: []string{"Nothing: smoothing is always allowed", "Smoothing makes the data more accurate", "Rewards cannot depend on data"},
		Why: "The reward's inputs must be immutable to the agent. Otherwise the cheapest way to improve is to change the reference."},
	{KP: 3, Q: "Why does a tool-grounded extraction reward eventually include a circuit benchmark such as ring-oscillator delay?",
		A: "Device-level error does not weight the bias regions that matter for circuits, and the benchmark checks what the card is for", D: []string{"Benchmarks are cheaper than device curves", "Benchmarks remove the need for measurements", "Ring oscillators measure leakage only"},
		Why: "Circuit metrics integrate current and charge over real transitions. They expose fits that are good on average but wrong where circuits operate."},
}

var tcadSurrogateItems = []edakit.Item{
	{KP: 1, Q: "What does an ML surrogate for TCAD learn?",
		A: "A fast map from device design inputs, such as geometry and doping, to outputs such as I-V curves or threshold, trained on TCAD runs", D: []string{"The layout of a full chip", "A replacement for Poisson's equation inside SPICE", "The foundry's measured data only"},
		Why: "A surrogate replaces repeated PDE solves with a fitted function, so a design-technology loop can evaluate thousands of options."},
	{KP: 1, Q: "In a DTCO loop, which role should the surrogate play and which should TCAD keep?",
		A: "The surrogate screens many candidates quickly; TCAD re-simulates the promising ones before decisions", D: []string{"The surrogate replaces TCAD completely", "TCAD screens and the surrogate confirms", "Neither: DTCO uses only measurements"},
		Why: "This is the fast-versus-accurate pattern. The fast engine explores, the accurate engine confirms."},
	{KP: 1, Q: "Why are surrogates attractive for device exploration specifically?",
		A: "Each 3D TCAD point can take minutes to hours, while a trained surrogate evaluates in milliseconds", D: []string{"Surrogates are always more accurate than TCAD", "TCAD cannot simulate FinFETs", "Surrogates need no training data"},
		Why: "The speed gap enables sweeps, optimization and uncertainty studies that would be unaffordable with TCAD alone."},
	{KP: 1, Q: "Which output representation makes a surrogate easiest to use in an extraction or circuit flow?",
		A: "Compact-model parameters or I-V curves on a fixed bias grid that downstream tools already consume", D: []string{"Raw mesh node values", "A single pass or fail flag", "Images of the device cross-section only"},
		Why: "Outputs aligned with existing interfaces let the surrogate drop into extraction and characterization without new tooling."},
	{KP: 2, Q: "A DTCO optimizer drives the surrogate to a fin geometry outside the training range and reports a record drive current. What should the loop do?",
		A: "Treat it as untrusted extrapolation and recurse to TCAD at that point before believing it", D: []string{"Accept the result and publish it", "Retrain the optimizer instead", "Lower the tolerance and continue"},
		Why: "Optimizers seek surrogate errors. Points outside the training distribution need accurate confirmation, and confirmed points can be added to training data."},
	{KP: 2, Q: "How can a surrogate tell the loop when to recurse to TCAD?",
		A: "By reporting an uncertainty estimate, for example from an ensemble, and flagging inputs far from training data", D: []string{"By always returning the same confidence", "By running slower when unsure", "It cannot, so the loop must always run TCAD"},
		Why: "Uncertainty-aware surrogates support active learning: spend TCAD runs where the surrogate is least sure."},
	{KP: 2, Q: "The TCAD deck is recalibrated to new silicon. What happens to the surrogate?",
		A: "It is stale, because it learned the old deck, and must be checked against new runs and retrained", D: []string{"Nothing: surrogates are independent of TCAD", "It improves automatically", "It must be deleted with no replacement"},
		Why: "A surrogate is only as current as its training data. Version the surrogate with the deck it was trained on."},
	{KP: 2, Q: "Which validation best shows that a TCAD surrogate is fit for a DTCO decision?",
		A: "Error on held-out TCAD runs in the region the decision depends on, plus agreement on the ranking of candidates", D: []string{"Training loss alone", "The number of network layers", "Speed of evaluation"},
		Why: "A decision depends on the right ordering of options in the region of interest. Held-out error and rank agreement test exactly that."},
	{KP: 3, Q: "A surrogate predicts drain current that sometimes decreases as gate voltage increases. Why is this unacceptable, and what fixes it?",
		A: "It violates physics and breaks downstream solvers; constrain monotonicity or train on derived quantities that guarantee it", D: []string{"It is acceptable if the average error is low", "Simulators ignore monotonicity", "Adding more layers always fixes it"},
		Why: "Physical constraints such as monotonicity, positivity and smoothness should be built into the surrogate's form or loss, not hoped for."},
	{KP: 3, Q: "A reinforcement-learning agent uses the surrogate as its reward for device design. What is the main risk?",
		A: "The agent exploits surrogate errors, so the reward rises while true TCAD performance does not", D: []string{"The reward is too accurate", "The agent cannot read the surrogate", "TCAD becomes faster"},
		Why: "Gaming a fast proxy is reward hacking. Periodic TCAD scoring of the agent's best designs keeps the reward honest."},
	{KP: 3, Q: "Why do surrogate-based device studies still need TCAD calibrated against silicon?",
		A: "The surrogate inherits whatever errors the TCAD deck has, so its accuracy is bounded by the calibration", D: []string{"Surrogates correct TCAD errors automatically", "Silicon data is not needed for TCAD", "Calibration only matters for metals"},
		Why: "A surrogate approximates TCAD, not silicon. Its best case is to reproduce a calibrated deck."},
	{KP: 3, Q: "How should training data for a TCAD surrogate be chosen within a fixed budget of TCAD runs?",
		A: "Space-filling designs to cover the input range, then active learning near the decision boundary or optimum", D: []string{"All runs at the nominal design", "Random runs with no record of inputs", "Only the extreme corners of the space"},
		Why: "Coverage first prevents blind spots; targeted sampling then spends accuracy where decisions are made."},
}
