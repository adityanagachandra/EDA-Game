package csim

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// fastAccurateSizing budgets a two-level sizing loop with a fast and an accurate simulator.
func fastAccurateSizing(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Pick([]int{200, 500, 1000, 2000, 5000})
	tf := pickF(0.5, 1, 2, 3, 6) // seconds per fast run
	k := edakit.Pick([]int{5, 10, 20, 25})
	corners := edakit.Pick([]int{5, 9, 15, 45})
	ta := pickF(1, 2, 3, 5) // minutes per accurate run
	total := float64(n)*tf/60 + float64(k*corners)*ta
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("A sizing agent screens %d candidate designs with a fast simulator at %s s per run (one corner), then verifies the best %d in accurate SPICE at %d corners, %s min per run. What is the total simulation time, in minutes?", n, f(tf, 1), k, corners, f(ta, 1)), total, 2,
			edakit.Sprintf("Screening: $%d \\times %s / 60 = %s$ min. Verification: $%d \\times %d \\times %s = %s$ min. Total %s min.", n, f(tf, 1), f(float64(n)*tf/60, 2), k, corners, f(ta, 1), f(float64(k*corners)*ta, 2), f(total, 2)), nil)
	case 2:
		all := float64(n*corners) * ta
		return num(edakit.Sprintf("A sizing agent screens %d candidates with a fast simulator at %s s per run, then verifies the best %d in accurate SPICE at %d corners, %s min per run. Running every candidate at every corner in accurate SPICE instead would cost how many times as much?", n, f(tf, 1), k, corners, f(ta, 1)), all/total, 1,
			edakit.Sprintf("All-accurate: $%d \\times %d \\times %s = %s$ min. Two-level: %s min. Ratio %s.", n, corners, f(ta, 1), f(all, 0), f(total, 2), f(all/total, 1)), nil)
	}
	spec := pickF(60, 65, 70)
	e := pickF(0.5, 1, 1.5, 2)
	m := edakit.Scaled(ctx, 6, 10)
	var vals []string
	count := 0
	for i := 0; i < m; i++ {
		v := spec + stepF(-4, 0.1, 80)
		vals = append(vals, f(v, 1))
		if math.Abs(v-spec) <= e+1e-9 {
			count++
		}
	}
	return edakit.Numeric(edakit.Sprintf("A fast simulator predicts gain within $\\pm %s$ dB of SPICE, validated on a golden set. The spec is gain $\\ge %s$ dB. Candidates whose fast gain clears the spec by more than the error bound are accepted, those below it by more than the bound are rejected, and the rest must be re-simulated accurately. The fast gains are %s dB. How many candidates must be re-simulated?", f(e, 1), f(spec, 0), strings.Join(vals, ", ")),
		edakit.IntStr(count),
		edakit.Sprintf("Re-simulate every value within %s dB of %s dB, the band $[%s, %s]$: %d candidates. Outside the band the fast answer cannot flip the decision.", f(e, 1), f(spec, 0), f(spec-e, 1), f(spec+e, 1), count), nil)
}

// surrogateCheck compares surrogate predictions with SPICE golden results.
func surrogateCheck(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Scaled(ctx, 5, 8)
	gold := make([]float64, m)
	pred := make([]float64, m)
	var rows []string
	for i := range gold {
		gold[i] = float64(edakit.Int(40, 300))
		pred[i] = gold[i] + float64(edakit.Int(-12, 12))
		if pred[i] == gold[i] && i%2 == 0 {
			pred[i]++
		}
		rows = append(rows, edakit.Sprintf("(%s, %s)", f(gold[i], 0), f(pred[i], 0)))
	}
	stem := edakit.Sprintf("An ML surrogate predicts cell delays and is checked against SPICE golden results. The pairs (SPICE, surrogate), in ps, are: %s.", strings.Join(rows, ", "))
	switch kp(ctx) {
	case 1:
		worst, wi := 0.0, 0
		for i := range gold {
			r := math.Abs(pred[i]-gold[i]) / gold[i] * 100
			if r > worst {
				worst, wi = r, i
			}
		}
		return num(stem+" What is the worst-case relative error, in percent?", worst, 2,
			edakit.Sprintf("Relative error is $|s - g|/g$. The worst is the pair (%s, %s): %s%%. A worst-case bound, not an average, is what a sign-off decision can rely on.", f(gold[wi], 0), f(pred[wi], 0), f(worst, 2)), nil)
	case 2:
		tol := pickF(2, 3, 4, 5)
		c := 0
		for i := range gold {
			if math.Abs(pred[i]-gold[i])/gold[i]*100 > tol {
				c++
			}
		}
		return edakit.Numeric(stem+edakit.Sprintf(" The acceptance rule allows at most %s%% relative error per point. How many points fail?", f(tol, 0)),
			edakit.IntStr(c),
			edakit.Sprintf("Compute $|s - g|/g$ for each pair and count those above %s%%: %d. Each failing point marks a region where the agent must fall back to SPICE.", f(tol, 0), c), nil)
	}
	s := 0.0
	for i := range gold {
		d := pred[i] - gold[i]
		s += d * d
	}
	rms := math.Sqrt(s / float64(m))
	return num(stem+" What is the RMS error of the surrogate over this golden set, in ps?", rms, 3,
		edakit.Sprintf("Square the differences, average over %d points (%s ps$^2$) and take the root: %s ps.", m, f(s/float64(m), 3), f(rms, 3)), nil)
}

var convergenceAgentItems = []edakit.Item{
	{KP: 1, Q: "An agent runs ngspice on a deck and the log reports 'singular matrix' at a node named float1 during the operating point. What should the agent infer first?",
		A: "That node has no DC path to ground, for example it touches only capacitors or gates", D: []string{"The timestep is too large", "The deck needs method=gear", "The device models are missing"},
		Why: "A singular DC matrix with a named node usually means that node's row is empty at DC. The diagnosis comes from the topology, not from tolerances."},
	{KP: 1, Q: "A transient fails with 'timestep too small' at 3.2 ns. Which information should the agent extract before proposing anything?",
		A: "Which nodes or devices were changing fastest or failing Newton near 3.2 ns, and what source edge happens there", D: []string{"Only the total run time", "The number of comment lines", "The final value of every node at the stop time"},
		Why: "Timestep failures are local in time and space. The agent needs the offending nodes and the event that triggered them to choose a targeted fix."},
	{KP: 1, Q: "An agent sees that the operating point converged only after gmin stepping reduced gmin to 1e-9 S and then stopped. What does that report mean?",
		A: "The operating point at the requested final gmin was not reached, so the result is not valid yet", D: []string{"The deck is solved and the result is reliable", "Gmin stepping is unnecessary in this deck", "The DC point is exact and only transient remains"},
		Why: "Intermediate gmin solutions are of a modified circuit. The agent must recognize partial convergence as a failure, not a success."},
	{KP: 1, Q: "An agent classifies simulator failures to choose a fix. Which pair is correctly distinguished?",
		A: "DC non-convergence points to initial guess and continuation aids; timestep-too-small points to discontinuities, models or tolerances in transient", D: []string{"Both are fixed by raising reltol", "Both mean the netlist has a syntax error", "DC failures are fixed with method=gear; transient failures with nodesets"},
		Why: "The phases fail for different reasons, and the classification determines which knobs are relevant."},
	{KP: 2, Q: "An agent can fix a DC convergence failure by adding .nodeset v(out)=1.2, by raising reltol to 1e-2, or by editing the diode model's equations. Which should it propose first?",
		A: "The .nodeset: the smallest change that leaves tolerances and models intact", D: []string{"Raising reltol, because it always works", "Editing the diode model", "All three at once"},
		Why: "The agent edits at the smallest safe tier. A nodeset changes only the starting guess, so a converged answer still satisfies the original equations."},
	{KP: 2, Q: "An agent proposes .options gmin=1e-6 to make a deck converge. Why should the harness reject this change?",
		A: "A large permanent gmin changes the circuit being solved and corrupts high-impedance nodes", D: []string{"gmin cannot be set in an options card", "It slows the simulation", "It has no effect on DC"},
		Why: "Some fixes make convergence easier by solving a different problem. The harness must catch changes that alter the answer, not just the path to it."},
	{KP: 2, Q: "A transient fails at a PWL source corner with zero rise time. What fix fits the smallest-change rule?",
		A: "Give the corner a small finite rise time consistent with the real driver", D: []string{"Set reltol to 0.1", "Delete the source", "Switch every device to a level-1 model"},
		Why: "An ideal step is unphysical and creates a discontinuity. A realistic edge removes the cause without loosening global tolerances."},
	{KP: 2, Q: "An agent's harness loop proposes an option change, reruns the deck, and gets a new failure. What should the next iteration use?",
		A: "The new failure message and location as feedback for the next proposal", D: []string{"The original failure only", "A random option", "A larger change to the core solver"},
		Why: "The harness loop feeds each run's result back to the model, so the next proposal responds to what actually happened."},
	{KP: 3, Q: "After a convergence fix, the deck runs. What must the agent check before reporting success?",
		A: "That key results agree with a trusted reference run or golden values within tolerance", D: []string{"Only that the simulator returned exit code zero", "That the run took less time", "That the log is shorter"},
		Why: "A deck can converge to the wrong solution or with loosened accuracy. Comparison against golden results proves the fix did not change the answer."},
	{KP: 3, Q: "An agent fixed a latch's DC point with .nodeset q=0. A reference simulation with q=1 also exists. Why might the results differ and both be correct?",
		A: "The latch has two valid DC solutions, and the nodeset selected one", D: []string{"One simulator must be wrong", "Nodesets add error proportional to their value", "The latch has no DC solution"},
		Why: "Multiple solutions are physical. The agent must report which solution was selected and why, rather than calling the difference a bug."},
	{KP: 3, Q: "An agent tightened reltol from 1e-3 to 1e-4 to remove a glitch, and the run passed. What differential check shows the glitch was numerical?",
		A: "The waveform converges as tolerances tighten, while the glitch disappears; a physical glitch would persist", D: []string{"The run time doubled", "The log has no warnings", "The glitch moved to a later time"},
		Why: "Numerical artifacts change with tolerances and step size; physical behavior converges to a limit. Comparing runs at several accuracies separates the two."},
	{KP: 3, Q: "A reward for a convergence-fixing agent is 'the simulator finished without error'. How can an agent game it?",
		A: "By loosening tolerances or adding gmin until anything runs, regardless of accuracy", D: []string{"By making the deck longer", "By using fewer nodes", "It cannot be gamed"},
		Why: "A pass-or-fail reward ignores correctness. The reward must also include agreement with golden results and penalties for changes that alter the circuit."},
}

var optionTierItems = []edakit.Item{
	{KP: 1, Q: "A designer asks to make a slow transient run faster by trading some accuracy. Which tier does the change land in?",
		A: "Tier 1: adjust options such as reltol, trtol or the integration method", D: []string{"Tier 2: write a new integrator plug-in", "Tier 3: edit the sparse LU core", "None: accuracy cannot be traded"},
		Why: "Tolerances and method selection are parameters of the engine, configuration only, with no recompile."},
	{KP: 1, Q: "A user wants a deck to run with Gear integration of order 2 instead of trapezoidal. What kind of change is this?",
		A: "A configuration change through simulator options", D: []string{"A new device model", "A core rewrite of the time loop", "A netlist topology change"},
		Why: "Simulators expose the integration method and maximum order as options, so choosing among built-in methods is Tier 1."},
	{KP: 1, Q: "Why should an agent prefer an option change over a code change when both fix a request?",
		A: "Options need no recompile, are easy to review and revert, and leave the verified core untouched", D: []string{"Code changes are always slower to run", "Options are more accurate", "Code changes cannot be tested"},
		Why: "Smallest safe tier first: configuration changes carry the least risk and the cheapest verification."},
	{KP: 1, Q: "Which of these requests belongs in Tier 1 for a circuit simulator?",
		A: "Run the same deck at five temperatures and three supply voltages", D: []string{"Add a new noise model for a device", "Replace the LU solver with an iterative one", "Change how the Jacobian is assembled"},
		Why: "Sweeps over temperature and supply are parameters of the analysis, handled by deck and option settings."},
	{KP: 2, Q: "A team needs a proprietary transistor model that is not in the simulator. Where should it enter?",
		A: "Tier 2: as a device plug-in at the device-model interface, for example compiled from Verilog-A", D: []string{"Tier 1: as an options card", "Tier 3: by editing the Newton loop", "By editing every netlist"},
		Why: "Device models are the classic extension point. A plug-in supplies currents, charges and derivatives without touching the solver core."},
	{KP: 2, Q: "A customer wants a custom convergence aid: ramp only the supply sources in a fixed order. What is the right tier and seam?",
		A: "Tier 2: a continuation strategy plug-in behind the operating-point homotopy interface", D: []string{"Tier 1: set gmin", "Tier 3: rewrite the LU factorization", "No change: source stepping already does this"},
		Why: "Convergence strategies are policies. A seam for homotopy strategies lets a new one be added while the Newton kernel stays fixed."},
	{KP: 2, Q: "A plug-in integrator is added at a seam. What must the seam's interface specify so that step control still works?",
		A: "The method's companion coefficients and its local truncation error estimate and order", D: []string{"Only the method's name", "The netlist format", "The plotting colors"},
		Why: "The time loop needs to stamp companions and to choose steps from error estimates, so a plug-in integrator must supply both."},
	{KP: 2, Q: "A device plug-in is called once per device per Newton iteration through a virtual call. A profile shows the call overhead dominates for millions of tiny devices. What design fixes it?",
		A: "A coarser seam: evaluate devices of one model in batches per call", D: []string{"Inline the plug-in into the netlist", "Remove derivatives from the interface", "Call the plug-in once per run"},
		Why: "Coarse plug-in calls amortize overhead and let the plug-in vectorize over many devices, keeping the hot loop efficient."},
	{KP: 3, Q: "Which change is genuinely Tier 3 in a SPICE-class simulator?",
		A: "Changing how the sparse matrix is ordered and factored", D: []string{"Choosing method=gear", "Adding a Verilog-A model", "Setting tstop"},
		Why: "Ordering and LU are the numerical kernel. Changing them affects every simulation and needs human review and full regression."},
	{KP: 3, Q: "A pull request edits the Newton convergence test in the core. What verification should the harness require before a human reviews it?",
		A: "The full golden regression suite and QoR comparison on reference decks", D: []string{"One deck that previously failed", "A compile-only check", "No tests, since the change is small"},
		Why: "A core change can alter results everywhere. The full golden suite and differential comparison show nothing unintended moved."},
	{KP: 3, Q: "Requests to change the transient time loop keep arriving, each handled as a Tier 3 edit. What does that pattern indicate?",
		A: "A seam is missing: the varying behavior should be moved behind an extension point", D: []string{"The core is well designed", "The requests should be refused", "Tier 1 options should be removed"},
		Why: "A seam is in the right place when Tier 3 edits essentially never happen for that kind of change. Repeated core edits signal where one belongs."},
	{KP: 3, Q: "An agent proposes a core edit to add a new device equation directly inside the matrix assembly loop. What should the harness steer it toward instead?",
		A: "Implementing the device behind the device-model interface as a plug-in", D: []string{"Accepting the edit if it compiles", "Adding the equation to the options card", "Hard-coding the device in every netlist"},
		Why: "The same behavior fits at a lower tier. Keeping device equations out of the assembly loop keeps the core small and stable."},
}

var deviceSeamItems = []edakit.Item{
	{KP: 1, Q: "At each Newton iteration the simulator calls a device model through its interface. What must the model return?",
		A: "Terminal currents and charges, and their derivatives with respect to the terminal voltages", D: []string{"Only the terminal currents", "The full MNA matrix", "The next timestep"},
		Why: "The core stamps currents into the residual and the derivatives (conductances and capacitances) into the Jacobian. Charges serve transient and AC."},
	{KP: 1, Q: "Why does the device interface return charges rather than capacitances for transient?",
		A: "The integrator differentiates charge, which conserves it for nonlinear capacitances", D: []string{"Capacitances cannot be represented as numbers", "Charges make the matrix symmetric", "Transient analysis has no Jacobian"},
		Why: "A charge-based interface lets the core apply any integration rule to $q(v)$ and keeps charge conserved."},
	{KP: 1, Q: "Which context should a device-model seam expose to the model, beyond terminal voltages?",
		A: "Temperature, instance parameters and the analysis mode, such as DC, transient, AC or noise", D: []string{"The full netlist text", "The plotting settings", "The other devices' internal state"},
		Why: "Models depend on temperature and geometry and must provide extra data in some analyses, such as noise sources in noise analysis."},
	{KP: 1, Q: "Why should the device interface also return noise sources?",
		A: "Noise analysis needs each device's noise power densities to propagate them to the output", D: []string{"To randomize the transient", "To set gmin", "To report Monte Carlo statistics"},
		Why: "Noise analysis sums device noise through transfer functions, so each model reports its thermal, shot and flicker densities at the operating point."},
	{KP: 2, Q: "A team writes a new transistor model in Verilog-A. How does it typically reach a SPICE-class simulator without editing the simulator core?",
		A: "A Verilog-A compiler generates code that loads through the simulator's model interface", D: []string{"It is pasted into the Newton loop", "It is converted to a .options card", "It replaces the LU solver"},
		Why: "Compiled Verilog-A plug-ins, for example through standardized compiled-model interfaces, implement the device seam while the core stays fixed."},
	{KP: 2, Q: "A plug-in model computes currents correctly but its derivatives contain a sign error. What will the simulator show?",
		A: "Slow or failed Newton convergence, while any converged answer is still correct", D: []string{"Wrong converged voltages", "Faster convergence", "A syntax error at load time"},
		Why: "The residual uses the currents, so a converged solution satisfies KCL. A wrong Jacobian only degrades the Newton iterations."},
	{KP: 2, Q: "A plug-in model keeps a static counter that changes its output on every call. What rule of good seams does this break?",
		A: "No hidden global state: results must be deterministic for the same inputs", D: []string{"Coarse calls", "A frozen ABI", "Narrow interfaces"},
		Why: "Hidden state makes results depend on call order, so repeated runs and differential tests stop being reproducible."},
	{KP: 2, Q: "Why should the device-model ABI be frozen once published?",
		A: "So compiled plug-ins keep working across simulator releases without recompilation", D: []string{"So models cannot be improved", "To prevent golden testing", "Because device models cannot change"},
		Why: "A stable extension ABI lets the core evolve independently of the plug-ins that depend on it."},
	{KP: 3, Q: "Which automated check catches a derivative bug in a new device plug-in?",
		A: "Comparing the returned Jacobian with finite differences of the returned currents at many bias points", D: []string{"Checking that the model compiles", "Running one DC point", "Comparing node names"},
		Why: "A finite-difference Jacobian check exercises exactly the interface contract and flags inconsistent derivatives before they slow Newton."},
	{KP: 3, Q: "A new model is meant to reproduce an existing one with a speed optimization. What does the differential test compare?",
		A: "Currents, charges and derivatives from both models over a grid of biases and temperatures, plus golden circuit results", D: []string{"Only compile time", "The source line count", "The number of parameters"},
		Why: "A differential test proves the optimization changed nothing beyond a numerical tolerance, at the device and the circuit level."},
	{KP: 3, Q: "A model passes its device-level tests, but a ring oscillator built from it runs 5% faster than the golden result. What should happen?",
		A: "Block promotion and investigate, since circuit-level golden tests disagree", D: []string{"Promote it, since device tests passed", "Update the golden result to match", "Loosen the tolerance to 6%"},
		Why: "Golden circuit tests catch interactions device tests miss, such as charge errors. A failed golden check blocks promotion."},
	{KP: 3, Q: "A plug-in model shows a discontinuity in drain current at Vds = 0 when tested. Why test that point specifically?",
		A: "Symmetry and continuity at Vds = 0 are where many models fail, causing Newton and distortion problems", D: []string{"Vds = 0 is never simulated", "It is the only point with nonzero current", "It tests the timestep control"},
		Why: "Source-drain symmetry tests near Vds = 0 are standard model quality checks, because kinks there cause convergence failures and wrong harmonic distortion."},
}
