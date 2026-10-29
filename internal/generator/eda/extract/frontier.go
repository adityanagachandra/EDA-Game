package extract

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var coreSurfaceItems = []edakit.Item{
	{KP: 1, Q: "In the four-layer model of an extractor, where does the discretized Laplace solve of the field solver belong?",
		A: "The immutable core, changed only with human approval and full regression", D: []string{"The parameter layer, tuned per design", "A plug-in strategy chosen per net", "An observation hook"},
		Why: "The field-solver mathematics is what makes the extractor correct. It is settled and shared by every customer, so it sits in the core."},
	{KP: 1, Q: "Which part of an extractor is change surface rather than core?",
		A: "The accuracy-versus-speed mode, net and layer selection, the pattern library and calibration coefficients", D: []string{"Netlist correctness: every node connected to the right net", "The capacitance-matrix assembly from solved charges", "The SPEF writer's syntax"},
		Why: "These vary by design, process and goal. Netlist correctness and solver math must not vary."},
	{KP: 1, Q: "What should an extractor's observation hook expose to an agent?",
		A: "Per-net parasitic summaries, coupling to each aggressor, and the estimated error of the mode used", D: []string{"The raw solver matrices for every net", "The customer's full layout database over the network", "Only a pass or fail flag"},
		Why: "An agent needs compact, typed state to decide what to re-extract or fix. Error estimates let it judge when fast results are trustworthy."},
	{KP: 1, Q: "Why must netlist correctness stay in the extractor's core, not in a plug-in?",
		A: "A parasitic attached to the wrong net corrupts every downstream analysis, whatever the accuracy mode", D: []string{"Netlist building is the slowest step", "Plug-ins cannot read layout", "Customers never ask to change it"},
		Why: "Connectivity is an invariant. A fast or approximate mode may trade accuracy, but never which net a resistor or capacitor belongs to."},
	{KP: 2, Q: "A customer asks to extract all clock nets with the 3D solver and everything else in rule-based mode. Which tier is this?",
		A: "Tier 1: a net-selection and mode setting", D: []string{"Tier 2: a new solver plug-in", "Tier 3: a core edit", "A new SPEF format"},
		Why: "Choosing a mode per net class is configuration, provided the extractor exposes net selection as a parameter."},
	{KP: 2, Q: "A foundry releases updated golden data for a new metal stack and asks to refit the pattern coefficients. Which tier is this?",
		A: "Tier 1: recalibration of tables and coefficients, validated against the golden data", D: []string{"Tier 2: a new pattern-matching engine", "Tier 3: a change to the field solver", "No change; calibration is fixed at release"},
		Why: "Calibration is data. Refitting it needs no recompile, only a validation run that shows the fitted tables match the golden solver."},
	{KP: 2, Q: "A new cell architecture introduces a via-rail structure that no existing pattern class matches, so the extractor needs new matching code. Which tier is this?",
		A: "Tier 2: a new pattern plug-in at the pattern-library seam", D: []string{"Tier 1: a coefficient change", "Tier 3: a core rewrite of connectivity", "It cannot be extracted at all"},
		Why: "A new geometric class needs new code to recognise it, but the core's fracturing and netlist assembly are untouched."},
	{KP: 2, Q: "An engineer proposes changing the extractor's mesh-refinement algorithm inside the 3D solver to speed it up. Which tier is this?",
		A: "Tier 3: a core edit needing a reviewed diff and full regression against golden results", D: []string{"Tier 1: a parameter change", "Tier 2: a plug-in", "No tier; it is a documentation change"},
		Why: "The solver's numerical kernel is core. Changes there can shift every golden result and must be proven against them."},
	{KP: 3, Q: "Where should the seam sit between the geometry engine and capacitance models in a rule-based extractor?",
		A: "At a capacitance-model interface that receives a fractured cross-section context and returns per-length ground and coupling values", D: []string{"Inside the SPEF writer", "Between the GDS reader and the layer map", "Inside each Laplace solve"},
		Why: "Fracturing and connectivity are fixed; how a cross-section maps to capacitance varies by process and accuracy mode. That is where models plug in."},
	{KP: 3, Q: "What context must the capacitance-model seam expose so that table, calibrated and surrogate models can all plug in?",
		A: "Layer, width, spacing to each neighbour and its net, the layers above and below, local density and the corner", D: []string{"Only the wire's length", "Only the net name", "The whole chip's GDS"},
		Why: "Different models key on different features. A seam that hides neighbour identity cannot produce coupling for SI."},
	{KP: 3, Q: "A plug-in capacitance model is called once per wire segment through a virtual function, and extraction slows tenfold. What does seam guidance recommend?",
		A: "Call the plug-in once per batch of segments with plain arrays, so the hot loop stays in the core", D: []string{"Move fracturing into the plug-in", "Remove the seam and hard-code one model", "Call the plug-in once per chip"},
		Why: "Coarse seams keep the compiler free to inline and vectorize the core's kernel."},
	{KP: 3, Q: "A surrogate capacitance model gives slightly different values on two runs of the same layout. Why is that unacceptable?",
		A: "Differential checks and incremental extraction need the same input to give the same parasitics", D: []string{"Surrogates are always less accurate", "SPEF cannot store the values", "Only field solvers may be nondeterministic"},
		Why: "Determinism is a seam rule. Without it, a differential check cannot tell an ECO's effect from noise."},
}

// calibrationTier1 refits a per-class scale factor and applies it.
func calibrationTier1(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 5)
	k := 1 + float64(edakit.Int(-10, 10))/100
	xs := make([]float64, n)
	gs := make([]float64, n)
	var rows []string
	sxy, sxx := 0.0, 0.0
	for i := 0; i < n; i++ {
		xs[i] = float64(edakit.Int(10, 80)) / 10
		gs[i] = edakit.Round(xs[i]*k*(1+float64(edakit.Int(-2, 2))/100), 2)
		sxy += xs[i] * gs[i]
		sxx += xs[i] * xs[i]
		rows = append(rows, edakit.Sprintf("%s fF fast vs %s fF golden", f(xs[i], 1), f(gs[i], 2)))
	}
	scale := sxy / sxx
	class := edakit.Pick([]string{"M2 coupling at minimum spacing", "M3 fringe over M2 dense fill", "M1-to-gate coupling in the cell rows", "M4 coupling over M3 crossings"})
	stem := edakit.Sprintf("An extractor's fast mode applies one scale factor per pattern class, stored in a configuration table. For the class %s, new golden 3D-solver samples give: %s. ", class, strings.Join(rows, "; "))
	switch kp(ctx) {
	case 1:
		return num(stem+"What least-squares scale factor $k$ (minimizing $\\sum (k x_i - g_i)^2$) should the Tier 1 table entry hold?", scale, 4,
			edakit.Sprintf("$k = \\sum x_i g_i / \\sum x_i^2$ = %s / %s = %s. Writing it into the table is a configuration change: no recompile, only a validation run.", f(sxy, 3), f(sxx, 3), f(scale, 4)), nil)
	case 2:
		x := float64(edakit.Int(10, 80)) / 10
		kr := edakit.Round(scale, 4)
		return num(stem+edakit.Sprintf("After the table is updated with the least-squares factor, rounded to 4 decimals (%s), what capacitance does fast mode report for a net of this class whose uncalibrated value is %s fF, in fF?", f(kr, 4), f(x, 1)), kr*x, 3,
			edakit.Sprintf("%s x %s = %s fF.", f(kr, 4), f(x, 1), f(kr*x, 3)), nil)
	}
	tol := edakit.Pick([]float64{1, 1.5, 2})
	kr := edakit.Round(scale, 4)
	m := edakit.Scaled(ctx, 4, 7)
	var vs []string
	bad := 0
	for i := 0; i < m; i++ {
		x := float64(edakit.Int(10, 80)) / 10
		g := edakit.Round(x*k*(1+float64(edakit.Int(-4, 4))/100), 2)
		if math.Abs(kr*x-g)/g*100 > tol {
			bad++
		}
		vs = append(vs, edakit.Sprintf("%s vs %s", f(x, 1), f(g, 2)))
	}
	return num(stem+edakit.Sprintf("The new factor, rounded to 4 decimals, is %s. Before promotion, the change is validated on held-out nets of the class (fast uncalibrated vs golden, fF): %s. How many held-out nets are off by more than %s%% of golden after calibration?", f(kr, 4), strings.Join(vs, "; "), f(tol, 1)), float64(bad), 0,
		edakit.Sprintf("Scale each fast value by %s and compare with golden: %d nets exceed %s%%. A Tier 1 change is promoted only if this validation passes; scatter that one factor cannot remove calls for a new pattern instead.", f(kr, 4), bad, f(tol, 1)), nil)
}

// fastAccurate decides which nets to re-extract accurately and what the loop saves.
func fastAccurate(ctx generator.GeneratorContext) generator.Problem {
	e := edakit.Pick([]float64{3, 4, 5, 8, 10})
	switch kp(ctx) {
	case 1:
		d := float64(edakit.Int(20, 300))
		return num(edakit.Sprintf("Fast extraction is within %s%% of the accurate extractor on net capacitance. A net's interconnect-driven delay is %s ps and scales linearly with its capacitance. What delay error bound, in ps, should an agent attach to this net's fast-mode timing?", f(e, 0), f(d, 0)),
			e/100*d, 2, edakit.Sprintf("%s%% of %s ps = %s ps. Slack within this band of zero cannot be trusted from fast mode.", f(e, 0), f(d, 0), f(e/100*d, 2)), nil)
	case 2:
		n := edakit.Scaled(ctx, 5, 8)
		var desc []string
		cnt := 0
		for i := 0; i < n; i++ {
			s := float64(edakit.Int(-15, 30))
			d := float64(10 * edakit.Int(2, 30))
			b := e / 100 * d
			if math.Abs(s) <= b {
				cnt++
			}
			desc = append(desc, edakit.Sprintf("net %d: slack %s ps, delay %s ps", i+1, f(s, 0), f(d, 0)))
		}
		return num(edakit.Sprintf("An agent's loop times a design with fast extraction, which is within %s%% of accurate on capacitance; each net's delay scales linearly with capacitance. It re-extracts a net accurately when the magnitude of its fast-mode slack is no more than its delay error bound. The nets are %s. How many nets does it send to accurate extraction?", f(e, 0), strings.Join(desc, "; ")),
			float64(cnt), 0, edakit.Sprintf("Each bound is %s%% of the net's delay. A net is ambiguous when |slack| is within its bound, since fast-mode error could flip the sign: %d nets qualify. Clearly passing or clearly failing nets keep their fast values.", f(e, 0), cnt), nil)
	}
	nets := edakit.Pick([]int{100000, 200000, 500000})
	frac := edakit.Pick([]float64{0.5, 1, 2, 5})
	tf := edakit.Pick([]float64{0.2, 0.5, 1})
	ta := tf * float64(edakit.Pick([]int{20, 50, 100}))
	full := float64(nets) * ta / 1000
	mixed := (float64(nets)*tf + frac/100*float64(nets)*ta) / 1000
	return num(edakit.Sprintf("A design has %d nets. Fast extraction costs %s ms per net and accurate extraction %s ms per net. An agent runs fast mode on every net, then accurate mode on the %s%% of nets near critical. What is the speed-up of this loop over running accurate extraction on every net?", nets, f(tf, 1), f(ta, 1), f(frac, 1)),
		full/mixed, 2, edakit.Sprintf("All accurate: %s s. Mixed: %s s fast + %s s accurate = %s s. Speed-up %s. The fraction recursed to the accurate engine dominates the cost, so a tight error bound pays directly.", f(full, 1), f(float64(nets)*tf/1000, 1), f(frac/100*float64(nets)*ta/1000, 1), f(mixed, 1), f(full/mixed, 2)), nil)
}

// differentialEcoCheck counts unexpected parasitic changes after an ECO.
func differentialEcoCheck(ctx generator.GeneratorContext) generator.Problem {
	tol := edakit.Pick([]float64{0.01, 0.02, 0.05})
	n := edakit.Scaled(ctx, 6, 9)
	type row struct {
		name           string
		before, after  float64
		intended, nbor bool
	}
	rows := make([]row, n)
	var desc []string
	var intended, nbors []string
	for i := range rows {
		b := float64(edakit.Int(10, 120)) / 10
		r := row{name: edakit.Sprintf("n%d", 10+i), before: b}
		switch edakit.Int(0, 4) {
		case 0:
			r.intended = true
		case 1:
			r.nbor = true
		}
		change := 0.0
		switch edakit.Int(0, 3) {
		case 0:
			change = 0
		case 1:
			change = float64(edakit.Int(1, 5)) / 100 * tol / 0.05
		default:
			change = float64(edakit.Int(1, 30)) / 10
		}
		if edakit.Int(0, 1) == 0 {
			change = -change
		}
		r.after = edakit.Round(b+change, 2)
		rows[i] = r
		desc = append(desc, edakit.Sprintf("%s %s to %s", r.name, f(r.before, 1), f(r.after, 2)))
		if r.intended {
			intended = append(intended, r.name)
		}
		if r.nbor {
			nbors = append(nbors, r.name)
		}
	}
	if len(intended) == 0 {
		rows[0].intended, rows[0].nbor = true, false
		intended = append([]string{rows[0].name}, intended...)
		nbors = nil
		for _, r := range rows {
			if r.nbor {
				nbors = append(nbors, r.name)
			}
		}
	}
	moved := func(r row) bool { return math.Abs(edakit.Round(r.after-r.before, 4)) > tol }
	nb := "none"
	if len(nbors) > 0 {
		nb = strings.Join(nbors, ", ")
	}
	stem := edakit.Sprintf("An agent's metal-only ECO claims to reroute %s. Its geometric neighbours, whose coupling may legitimately change, are %s. Re-extraction gives total net capacitance before and after, in fF: %s. A net has moved if its capacitance changed by more than %s fF.", strings.Join(intended, ", "), nb, strings.Join(desc, "; "), f(tol, 2))
	cntMoved, unexpected, stuck := 0, 0, 0
	for _, r := range rows {
		if moved(r) {
			cntMoved++
			if !r.intended && !r.nbor {
				unexpected++
			}
		} else if r.intended {
			stuck++
		}
	}
	switch kp(ctx) {
	case 1:
		return num(stem+" How many nets moved?", float64(cntMoved), 0,
			edakit.Sprintf("Compare |after - before| with %s fF for each net: %d moved.", f(tol, 2), cntMoved), nil)
	case 2:
		return num(stem+" How many nets moved that are neither rerouted nor neighbours of a rerouted net?", float64(unexpected), 0,
			edakit.Sprintf("Remove the rerouted nets and their neighbours from the moved set: %d unexpected movers remain. Each is evidence that the ECO touched more than it claimed.", unexpected), nil)
	}
	return num(stem+" The check fails once for each unexpected mover and once for each rerouted net that did not move. How many failures does the check report?", float64(unexpected+stuck), 0,
		edakit.Sprintf("Unexpected movers: %d. Rerouted nets that did not move: %d. Total %d. The ECO is accepted only at zero, which turns the agent's claim into a verified one.", unexpected, stuck, unexpected+stuck), nil)
}

var agenticSignoffItems = []edakit.Item{
	{KP: 1, Q: "In an agentic signoff loop, what does the agent read at each iteration?",
		A: "Typed reports from the engines: timing and SI violations, IR and EM hot spots, DRC markers", D: []string{"The raw layout polygons of the whole chip", "Only the final tape-out checklist", "The RTL source code"},
		Why: "The agent reasons over the engines' outputs. Compact, structured reports let it localize problems and propose targeted ECOs."},
	{KP: 1, Q: "Who approves the final ECO in an agentic signoff flow?",
		A: "A human engineer, after the agent's change has passed the engines' checks", D: []string{"The agent itself, once its reward is positive", "The fast extractor", "No one; ECOs are applied automatically"},
		Why: "Agents propose and verify; humans own tape-out decisions. The agent's job is to hand over a change already proven by the tools."},
	{KP: 1, Q: "An agent proposes a fix for an SI-induced setup violation. What must happen before it is shown to a human?",
		A: "The ECO is applied, the touched region re-extracted and timing, SI, EM/IR and DRC re-run on it", D: []string{"The agent explains its reasoning in prose", "The fix is applied directly to the tape-out database", "Only the failing path is re-timed with old parasitics"},
		Why: "Proposals are cheap. Evidence from the engines, on fresh parasitics, is what makes a change reviewable."},
	{KP: 1, Q: "Why does an agentic loop use fast extraction for most iterations?",
		A: "Each candidate fix needs parasitics quickly, and fast mode is accurate enough to rank candidates away from the margin", D: []string{"Fast extraction is more accurate than signoff extraction", "Accurate extraction cannot handle ECOs", "Fast mode is required by the foundry"},
		Why: "Speed lets the agent explore. It recurses to the accurate extractor on near-critical nets and for the final signoff check."},
	{KP: 2, Q: "The agent's loop alternates: a setup fix breaks hold, the hold fix breaks setup, for ten iterations. What should the loop do?",
		A: "Detect the oscillation and propose fixes that satisfy both checks together, or escalate to a human", D: []string{"Keep alternating until it converges by itself", "Disable hold checking", "Accept whichever result came last"},
		Why: "A loop needs convergence guards. Joint fixes, such as delaying the short path while speeding the long one, address the coupling between the checks."},
	{KP: 2, Q: "An agent is rewarded for reducing the count of timing violations. It learns to add false-path exceptions. What failure is this?",
		A: "Reward hacking: it games the metric instead of fixing the timing", D: []string{"A correct optimization", "A fast-mode accuracy problem", "An LVS failure"},
		Why: "Constraints are part of the design intent. Changes to them must be blocked or require human approval, or the reward will be gamed."},
	{KP: 2, Q: "Which change should an agentic signoff loop never make without explicit human approval?",
		A: "Changing timing constraints or signoff margins", D: []string{"Upsizing a cell", "Inserting a hold buffer", "Rerouting a net with wider spacing"},
		Why: "Physical fixes are verifiable by the engines. Constraints define what correct means, so editing them changes the target itself."},
	{KP: 2, Q: "An agent fixes a crosstalk violation by moving the victim to a free track. Which check proves the move did not hurt other nets?",
		A: "A differential re-extraction and SI report showing that only the victim and its old and new neighbours changed", D: []string{"The victim's slack is now positive", "The ECO compiled without errors", "The agent's confidence score"},
		Why: "A local fix can create new victims. The differential check compares the whole region before and after against the claimed effect."},
	{KP: 3, Q: "Fast extraction says a net has 2 ps of slack after an agent's ECO, and fast mode's error bound on that net is 5 ps. What should the agent do?",
		A: "Re-extract the net accurately before accepting the fix", D: []string{"Accept the fix, since slack is positive", "Reject the fix, since fast mode is always optimistic", "Add 5 ps of margin to the clock"},
		Why: "Slack inside the error band is ambiguous. The fast-accurate loop recurses to the accurate engine exactly there."},
	{KP: 3, Q: "Why should the signoff agent run inside a sandbox on the customer's site?",
		A: "Design data must not leave the IP boundary, and tool runs need CPU, memory and time limits", D: []string{"Sandboxes make extraction more accurate", "The foundry requires agents to be offline", "Agents cannot call engines over a network"},
		Why: "The IP boundary is architectural: layouts and parasitics stay on-site, and limits keep runaway loops from consuming the cluster."},
	{KP: 3, Q: "After many projects, an agent's spacing-based fix for coupling on clock nets keeps working. What is the flywheel step?",
		A: "Validate it with golden and differential tests and promote it to a library skill that any agent can select", D: []string{"Hard-code it into the extractor core", "Keep it as a private prompt", "Turn off SI checks for clock nets"},
		Why: "Skills found in the loop are validated, then promoted. The know-how moves into the library; the customer's design does not."},
	{KP: 3, Q: "An agent proposes 400 independent small ECOs at once to close the last violations. Why might a human reviewer push back?",
		A: "Each change interacts through coupling and shared logic, so a large batch is hard to verify and to roll back", D: []string{"ECOs can only be applied one per day", "The extractor cannot handle 400 nets", "Small ECOs never close timing"},
		Why: "Batches should be sized so that their effects can be checked and reversed. Large, entangled changes defeat the review step."},
}
