package sta

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// derateProviderCalc applies a Tier 2 provider rule to a path and asks for its derated delay.
func derateProviderCalc(ctx generator.GeneratorContext) generator.Problem {
	type a struct {
		name, kind, vt string
		d              float64
	}
	vts := []string{"HVT", "SVT", "LVT"}
	n := edakit.Int(4, 6)
	arcs := make([]a, n)
	for i := range arcs {
		arcs[i] = a{edakit.Sprintf("U%d", i+1), "data", edakit.Pick(vts), float64(10 * edakit.Int(2, 9))}
	}
	arcs[0].kind, arcs[0].vt = "clock", "HVT"
	arcs[1].vt = "HVT"
	arcs[2].vt = edakit.Pick([]string{"SVT", "LVT"})
	arcs = append(arcs[:1], edakit.Shuffle(arcs[1:])...)
	depth := edakit.Int(2, 8)
	base := 1 + float64(edakit.Int(3, 8))/100
	var parts []string
	for _, x := range arcs {
		parts = append(parts, edakit.Sprintf("%s (%s arc, %s, %s ps)", x.name, x.kind, x.vt, edakit.Num(x.d, 0)))
	}
	fires := depth < 5
	total, allHVT, none, mult, clockToo := 0.0, 0.0, 0.0, 0.0, 0.0
	for _, x := range arcs {
		hvtData := x.kind == "data" && x.vt == "HVT"
		none += x.d * base
		allHVT += x.d * (base + 0.03)
		if hvtData && fires {
			total += x.d * (base + 0.03)
			mult += x.d * base * 1.03
		} else {
			total += x.d * base
			mult += x.d * base
		}
		if x.vt == "HVT" && fires {
			clockToo += x.d * (base + 0.03)
		} else {
			clockToo += x.d * base
		}
	}
	q := edakit.Sprintf("A team registers a Tier 2 derate provider: on data arcs through HVT cells, when the path's logic depth is below 5, it adds 0.03 to the late derate factor; everywhere else it falls through to the AOCV table, which gives %s at this depth. A late path of logic depth %d has the arcs %s. What late path delay does the timer compute with this provider?",
		edakit.Num(base, 2), depth, strings.Join(parts, ", "))
	exp := edakit.Sprintf("The rule fires only when depth < 5 (here depth %d) and only on data arcs through HVT cells. Summing delay times factor arc by arc gives %s ps.", depth, edakit.Num(total, 2))
	cand := []string{edakit.Num(allHVT, 2), edakit.Num(none, 2), edakit.Num(mult, 2), edakit.Num(clockToo, 2), edakit.Num(total+0.03*100, 2)}
	ans := edakit.Num(total, 2)
	var d []string
	seen := map[string]bool{ans: true}
	for _, c := range cand {
		if !seen[c] {
			seen[c] = true
			d = append(d, c+" ps")
		}
	}
	for k := 1; len(d) < 3; k++ {
		c := edakit.Num(total+float64(k)*1.5, 2)
		if !seen[c] {
			seen[c] = true
			d = append(d, c+" ps")
		}
	}
	return edakit.Choice(q, ans+" ps", exp, d...)
}

var derateProviderItems = []edakit.Item{
	{KP: 1, Q: "In a timer built around a DerateProvider seam, what does the propagation core ask the provider?",
		A: "What derate applies to this arc in this context", D: []string{"What the slack of this endpoint should be", "Which paths to report", "How to propagate arrival times"},
		Why: "The core keeps propagation and delay calculation. At one point it asks the provider for a factor, or a mean and sigma adjustment, for one arc in its context."},
	{KP: 1, Q: "Which context must the seam pass so that AOCV, POCV and custom margin providers can all be written against it?",
		A: "Arc, cell and Vt class, early or late, rise or fall, logic depth, distance, slew, corner, clock or data, and path group", D: []string{"Only the arc delay", "Only the corner name", "The full netlist and SDC as text"},
		Why: "Each scheme keys on different fields. A seam that hides depth or distance cannot host AOCV; one without clock-versus-data cannot host clock-specific margins."},
	{KP: 1, Q: "What does a POCV provider return through the seam, compared with a flat OCV provider?",
		A: "A mean and sigma adjustment rather than a single multiplicative factor", D: []string{"A required time", "A list of false paths", "The same factor for every arc"},
		Why: "Flat OCV and AOCV scale a delay. POCV works with distributions, so the seam must carry a mean and a sigma that the core combines statistically."},
	{KP: 1, Q: "Why is derating a good place for a seam in a timer?",
		A: "The propagation math is settled while every team margins differently", D: []string{"Derates are computed once per chip", "Derating is the slowest part of timing", "Derates never affect slack"},
		Why: "A seam belongs where behavior varies across customers but the surrounding core does not. Margin schemes vary; arrival and slack arithmetic do not."},
	{KP: 3, Q: "A provider is called once per arc evaluation through a virtual function. Profiling shows the hot loop no longer vectorizes. What does the seam design guidance recommend?",
		A: "Make the plug-in call coarse, per batch of arcs, so the core's kernel stays inlined", D: []string{"Move propagation into the plug-in", "Call the provider once per path instead", "Remove the seam and hard-code AOCV"},
		Why: "Seams should be coarse-grained so the hot kernel stays in the core and the compiler can still optimize it."},
	{KP: 3, Q: "A provider keeps a global counter and returns slightly different derates on the second run of the same design. What rule did it break?",
		A: "Providers must be deterministic with no hidden global state", D: []string{"Providers must be written in the same language as the core", "Providers must run only on clock arcs", "Providers must return factors above 1"},
		Why: "Reproducible timing and differential tests require the same inputs to give the same derates. Hidden state breaks both."},
	{KP: 3, Q: "Why freeze the provider's extension ABI?",
		A: "So plug-ins built against it keep working as the core evolves", D: []string{"To stop customers writing providers", "To make derates constant", "Because the ABI holds the derate tables"},
		Why: "A frozen interface lets providers be compiled and promoted independently of core releases."},
	{KP: 3, Q: "A new provider must be accepted. What evidence shows it is correct?",
		A: "A differential timing report on a reference design showing that only the intended arcs and paths moved, by the intended amounts", D: []string{"The provider compiles without warnings", "The worst slack improved", "The provider runs faster than AOCV"},
		Why: "A change to derating must be proven by its effect: intended paths shift by the expected amount and nothing else moves."},
	{KP: 3, Q: "Where should the AOCV table lookup live in a well-designed timer?",
		A: "In an AOCV provider behind the seam, not inside the propagation loop", D: []string{"Inside the arrival propagation code", "In the SDC reader", "In the Liberty parser"},
		Why: "Keeping scheme-specific logic in providers keeps the core small, so swapping schemes is a configuration or plug-in change rather than a core edit."},
}

var changeTierItems = []edakit.Item{
	{KP: 1, Q: "A customer asks to run signoff at 3.5 sigma instead of 3 sigma with the existing POCV provider. Which tier is this?",
		A: "Tier 1, a configuration change", D: []string{"Tier 2, a new plug-in", "Tier 3, a core edit", "It cannot be done without retraining a model"},
		Why: "The sigma count is a parameter of the provider. No code changes and no recompile are needed."},
	{KP: 1, Q: "A team wants to switch from flat OCV to an AOCV provider that already ships with the timer, pointing at their own tables. Which tier is this?",
		A: "Tier 1: choose the provider and the table paths", D: []string{"Tier 2: write a new provider", "Tier 3: change the propagation core", "Tier 3: change the delay-model API"},
		Why: "Selecting an existing provider and its data is configuration."},
	{KP: 1, Q: "A project needs an extra 10 ps of margin on one path group only. Which tier is this?",
		A: "Tier 1, a per-path-group margin setting", D: []string{"Tier 2, a new provider", "Tier 3, a core edit", "A change to the Liberty library"},
		Why: "Per-group margins are a parameter of the custom margin layer, so this is configuration."},
	{KP: 1, Q: "A customer asks the timer to report the 20 worst paths per endpoint instead of 5. Which tier is this?",
		A: "Tier 1, a reporting parameter", D: []string{"Tier 2, a new path enumerator plug-in", "Tier 3, a change to propagation", "It requires a new SDC command"},
		Why: "Which paths and how many to report is part of the timer's change surface, set by parameters."},
	{KP: 2, Q: "A foundry correlation shows that HVT data paths with logic depth under 5 need 3% more late derate, while clocks keep AOCV. Which tier is this?",
		A: "Tier 2, a new provider at the derate seam", D: []string{"Tier 1, a global derate change", "Tier 3, a change to arrival propagation", "Tier 1, a new sigma count"},
		Why: "No existing parameter expresses a rule keyed on Vt, depth and data versus clock. A small provider implements it without touching the core."},
	{KP: 2, Q: "A team wants derates that depend on the arc's input slew bucket, a field the seam already passes. Which tier is this?",
		A: "Tier 2, a provider that reads the slew from the context", D: []string{"Tier 3, because slew is a core quantity", "Tier 1, a global derate", "It cannot be expressed"},
		Why: "The seam already exposes slew, so a plug-in can key on it. The core is untouched."},
	{KP: 2, Q: "A customer asks for derates keyed on a per-instance attribute that the seam's context does not carry. What is the likely path?",
		A: "A reviewed core change to extend the context, then a provider: Tier 3 once, Tier 2 thereafter", D: []string{"Tier 1, a parameter", "Tier 2 alone, by reading global state", "Edit the SDC to carry the attribute"},
		Why: "If the seam lacks the information, a plug-in cannot see it without hidden state. Extending the context is a core edit that needs human approval."},
	{KP: 2, Q: "A new provider stacks a silicon-correlation margin on top of AOCV. How should it be built?",
		A: "As a custom margin layer that wraps the AOCV provider's result", D: []string{"By editing the AOCV table format in the core", "By changing every Liberty file", "By adding uncertainty to all clocks"},
		Why: "Stackable margin layers compose with existing providers, keeping each piece small and separately testable."},
	{KP: 3, Q: "A request asks to change how arrival times combine at a multi-input gate. Which tier is this?",
		A: "Tier 3, a core edit that needs a human-approved diff and full regression", D: []string{"Tier 1, a parameter", "Tier 2, a derate provider", "Tier 2, a reporting plug-in"},
		Why: "Propagation math is the immutable core. Changes to it affect every result and are rare and reviewed."},
	{KP: 3, Q: "Requests for derate changes keep requiring Tier 3 edits. What does that say about the design?",
		A: "The seam is in the wrong place or exposes too little context", D: []string{"Derates should always be core code", "Tier 3 is the normal tier for derating", "The customers are asking for the wrong things"},
		Why: "A well-placed seam makes core edits essentially never necessary for that kind of change."},
	{KP: 3, Q: "An agent proposes changing the delay-model API to return an extra field for its new provider. What should the harness do?",
		A: "Route it as a Tier 3 change: produce a reviewable diff, run full regression, and wait for human approval", D: []string{"Merge it automatically if it compiles", "Treat it as Tier 1 because it is small", "Reject all API changes permanently"},
		Why: "The delay-model API is core. Small diffs there still change the contract for every provider."},
	{KP: 3, Q: "Which change is least risky to apply without review?",
		A: "Choosing a different shipped provider and its table file", D: []string{"Editing required-time propagation", "Changing slack sign conventions", "Altering the delay-model API"},
		Why: "Tier 1 changes alter configuration only, are reversible and do not touch compiled code."},
}

func differentialTiming(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Scaled(ctx, 8, 12)
	delta := float64(edakit.Int(2, 6))
	eps := 0.5
	before := make([]float64, m)
	after := make([]float64, m)
	inS := make([]bool, m)
	for i := 0; i < m; i++ {
		before[i] = float64(edakit.Int(-30, 80))
		inS[i] = edakit.Int(0, 2) == 0
	}
	inS[0] = true
	for i := 0; i < m; i++ {
		after[i] = before[i]
		if inS[i] {
			after[i] -= delta
			if edakit.Int(0, 5) == 0 {
				after[i] -= float64(edakit.Int(1, 3))
			}
		} else if edakit.Int(0, 4) == 0 {
			after[i] -= float64(edakit.Int(1, 4))
		}
	}
	var rows, sNames []string
	moved, unexpected, wrong := 0, 0, 0
	for i := 0; i < m; i++ {
		rows = append(rows, edakit.Sprintf("E%d: %s → %s", i+1, edakit.Num(before[i], 1), edakit.Num(after[i], 1)))
		if inS[i] {
			sNames = append(sNames, edakit.Sprintf("E%d", i+1))
		}
		diff := after[i] - before[i]
		if math.Abs(diff) > eps {
			moved++
			if !inS[i] {
				unexpected++
			}
		}
		if inS[i] && math.Abs(diff+delta) > eps {
			wrong++
		}
	}
	q := edakit.Sprintf("A new derate provider should lower the slack of the paths ending at %s by exactly %s ps and leave every other endpoint unchanged, with a tolerance of %s ps. The differential report lists before → after slacks in ps: %s. ",
		strings.Join(sNames, ", "), edakit.Num(delta, 0), edakit.Num(eps, 1), strings.Join(rows, ", "))
	facts := map[string]string{"moved": edakit.IntStr(moved), "unexpected": edakit.IntStr(unexpected), "wrong_amount": edakit.IntStr(wrong)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+"How many endpoints moved by more than the tolerance?", edakit.IntStr(moved),
			edakit.Sprintf("Compare after with before at each endpoint: %d moved.", moved), facts)
	case 2:
		return edakit.Numeric(q+"How many endpoints moved that should not have?", edakit.IntStr(unexpected),
			edakit.Sprintf("Endpoints outside the intended set must not move; %d did. Each one is a failure of the golden test.", unexpected), facts)
	}
	tot := unexpected + wrong
	return edakit.Numeric(q+"How many endpoints violate the expectation in total (unexpected movers plus intended endpoints whose change is not the intended amount)?", edakit.IntStr(tot),
		edakit.Sprintf("%d unexpected movers plus %d intended endpoints off by more than %s ps gives %d. The provider is promoted only when this is zero.", unexpected, wrong, edakit.Num(eps, 1), tot), facts)
}

func surrogateTimer(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Scaled(ctx, 8, 14)
	pred := make([]float64, m)
	var rows []string
	for i := range pred {
		pred[i] = float64(edakit.Int(-40, 60)) + 0.5
		rows = append(rows, edakit.Sprintf("E%d: %s", i+1, edakit.Num(pred[i], 1)))
	}
	q := "A fast surrogate timer inside an agent's optimization loop predicts these endpoint setup slacks, in ps: " + strings.Join(rows, ", ") + ". "
	count := func(eps float64) int {
		c := 0
		for _, p := range pred {
			if math.Abs(p) < eps {
				c++
			}
		}
		return c
	}
	switch ctx.KPOrdinal {
	case 2:
		eps := float64(edakit.Int(3, 15))
		w := minInt(func() []int {
			out := make([]int, 0, m)
			for _, p := range pred {
				out = append(out, int(math.Round(p*10)))
			}
			return out
		}())
		lb := float64(w)/10 - eps
		ans := edakit.Num(lb, 1)
		return edakit.Numeric(q+edakit.Sprintf("Validation shows its error never exceeds %s ps. What is the lowest value the true worst slack could have, in ps? Round to 1 decimal place.", edakit.Num(eps, 0)), ans,
			edakit.Sprintf("The worst prediction is %s ps; the truth can be up to %s ps lower: %s ps.", edakit.Num(float64(w)/10, 1), edakit.Num(eps, 0), ans), map[string]string{"bound": ans})
	case 3:
		sg := float64(edakit.Int(2, 6))
		k := float64(edakit.Pick([]int{2, 3}))
		eps := k * sg
		c := count(eps)
		return edakit.Numeric(q+edakit.Sprintf("Its validation error has sigma %s ps, and the loop treats %s sigma as the error bound. An endpoint must be re-timed by the accurate timer when its true slack could have either sign. How many endpoints must be re-timed?", edakit.Num(sg, 0), edakit.Num(k, 0)), edakit.IntStr(c),
			edakit.Sprintf("The bound is %s x %s = %s ps. Endpoints with |predicted slack| below %s ps are uncertain: %d of them.", edakit.Num(k, 0), edakit.Num(sg, 0), edakit.Num(eps, 0), edakit.Num(eps, 0), c), map[string]string{"bound": edakit.Num(eps, 0), "count": edakit.IntStr(c)})
	}
	eps := float64(edakit.Int(3, 15))
	c := count(eps)
	return edakit.Numeric(q+edakit.Sprintf("Validation shows its error never exceeds %s ps. An endpoint must be re-timed by the accurate timer when its true slack could have either sign. How many endpoints must be re-timed?", edakit.Num(eps, 0)), edakit.IntStr(c),
		edakit.Sprintf("True slack lies within %s ps of the prediction, so the sign is uncertain when |predicted slack| < %s ps: %d endpoints.", edakit.Num(eps, 0), edakit.Num(eps, 0), c), map[string]string{"count": edakit.IntStr(c)})
}

var ecoLoopItems = []edakit.Item{
	{KP: 1, Q: "What should an agentic timing-ECO loop read before proposing a fix?",
		A: "Timing reports across all active scenarios, plus the physical context of the failing paths", D: []string{"Only the worst path at the typical corner", "Only the SDC file", "Only the agent's previous proposals"},
		Why: "A fix that helps one scenario can hurt another, and whether sizing or buffering is possible depends on placement and routing context."},
	{KP: 1, Q: "Which actions belong in the ECO agent's proposal vocabulary?",
		A: "Gate sizing, Vt swaps, buffer insertion or removal, and useful-skew adjustments", D: []string{"Editing the clock period in the SDC", "Adding false paths to failing endpoints", "Changing the Liberty tables"},
		Why: "ECO moves change the netlist or clock tree. Changing constraints or libraries hides failures instead of fixing them."},
	{KP: 1, Q: "An agent proposes fixing a setup violation with set_false_path on the endpoint. What should the harness do?",
		A: "Reject it: constraint edits are not timing fixes and need separate human justification", D: []string{"Accept it, since the slack becomes positive", "Accept it at the fast corner only", "Convert it to a multicycle path automatically"},
		Why: "Removing a check makes the report pass without making the chip work. Constraint changes are outside the ECO agent's authority."},
	{KP: 1, Q: "Why must the loop re-time incrementally after each proposed change?",
		A: "Each change shifts arrivals and slews in its cone, and the next proposal must see the updated state", D: []string{"To recompute the whole netlist from RTL", "Because reports are cached forever otherwise", "To change the clock definitions"},
		Why: "Incremental timing gives the updated slacks cheaply, so the loop proposes against the real current state."},
	{KP: 2, Q: "A loop alternately fixes setup and breaks hold on the same paths for many iterations. What is the remedy?",
		A: "Evaluate setup and hold together in every scenario before accepting a change", D: []string{"Fix only hold", "Run more iterations", "Raise the hold uncertainty"},
		Why: "Oscillation means each move is judged on one check only. Accepting a change only when no other check degrades breaks the cycle."},
	{KP: 2, Q: "The loop uses a fast surrogate timer to rank candidate fixes. Before accepting the chosen fix, what must it do?",
		A: "Confirm the result with the accurate signoff timer", D: []string{"Accept the surrogate's slack", "Average the two timers", "Skip timing and run DRC only"},
		Why: "The surrogate is a fast proxy with error. Acceptance must rest on the accurate engine."},
	{KP: 2, Q: "A fix improves WNS but adds 2% area and violates a max-transition rule. How should the loop score it?",
		A: "As failing, since design-rule and area constraints are part of acceptance", D: []string{"As passing, since WNS improved", "As passing if TNS improved", "As neutral"},
		Why: "Timing fixes must keep electrical rules and budgets intact; otherwise the loop trades one violation for another."},
	{KP: 2, Q: "Which check proves an accepted ECO changed only what it claimed?",
		A: "A differential timing report against the pre-ECO baseline", D: []string{"A new run of synthesis", "A lower WNS alone", "A shorter run time"},
		Why: "The differential shows every endpoint that moved, so unintended side effects become visible."},
	{KP: 3, Q: "Where does the human sit in an agentic signoff ECO loop?",
		A: "Approving the final ECO, with the loop's evidence attached", D: []string{"Nowhere; the loop tapes out on its own", "Writing every individual fix", "Only reading the final WNS"},
		Why: "The loop proposes, executes and verifies at engine speed; a human approves the result before it becomes part of the design."},
	{KP: 3, Q: "Why should the ECO loop run inside a sandbox with resource limits and no network?",
		A: "It executes generated actions on proprietary design data, which must stay on site and be contained", D: []string{"Sandboxes make timing more accurate", "Agents cannot run without one", "To make the loop nondeterministic"},
		Why: "Design data is IP, and agent actions must be bounded. The sandbox is part of the architecture."},
	{KP: 3, Q: "After a project, the team wants the ECO tactics the agent discovered to help future designs. What flows back?",
		A: "The validated tactic as a skill, not the customer's design data", D: []string{"The entire netlist", "The parasitics of the design", "Nothing can be reused"},
		Why: "In the skill flywheel the know-how is promoted into a library skill while the design stays inside the IP boundary."},
	{KP: 3, Q: "An agent's ECO passes in the implementation tool's timer but fails at signoff. What should the loop learn?",
		A: "Its acceptance check must use the signoff engine, or keep margin for the correlation gap", D: []string{"Signoff is wrong", "Lower the signoff sigma", "Ignore signoff for ECOs"},
		Why: "Fast engines and signoff engines differ slightly. The loop must close on the engine that decides tapeout."},
}

func slackReward(ctx generator.GeneratorContext) generator.Problem {
	a := float64(edakit.Pick([]int{1, 2, 5}))
	b := edakit.Pick([]float64{0.01, 0.05, 0.1})
	gm := float64(edakit.Pick([]int{2, 5, 10}))
	formula := edakit.Sprintf("An RL environment scores an agent's timing ECO with reward $R = %s \\min(WNS, 0) + %s\\,TNS - %s\\,\\Delta A$, where WNS and TNS are in ps and $\\Delta A$ is the area increase in percent. ", edakit.Num(a, 0), edakit.Num(b, 2), edakit.Num(gm, 0))
	da := float64(edakit.Int(0, 30)) / 10
	switch ctx.KPOrdinal {
	case 1:
		w := float64(edakit.Int(-60, 10))
		t := w * float64(edakit.Int(1, 30))
		if w > 0 {
			t = 0
		}
		r := a*math.Min(w, 0) + b*t - gm*da
		ans := edakit.Num(r, 2)
		return edakit.Numeric(formula+edakit.Sprintf("The signoff timer reports WNS %s ps and TNS %s ps, and area grew %s%%. What is the reward? Round to 2 decimal places.", edakit.Num(w, 0), edakit.Num(t, 0), edakit.Num(da, 1)), ans,
			edakit.Sprintf("R = %s x %s + %s x %s - %s x %s = %s.", edakit.Num(a, 0), edakit.Num(math.Min(w, 0), 0), edakit.Num(b, 2), edakit.Num(t, 0), edakit.Num(gm, 0), edakit.Num(da, 1), ans), map[string]string{"reward": ans})
	}
	m := edakit.Int(5, 9)
	sl := make([]int, m)
	var rows []string
	for i := range sl {
		sl[i] = edakit.Int(-40, 50)
		rows = append(rows, edakit.Sprintf("E%d: %d", i+1, sl[i]))
	}
	sl[0] = -edakit.Int(5, 40)
	rows[0] = edakit.Sprintf("E1: %d", sl[0])
	score := func(xs []int) (float64, float64, float64) {
		w := float64(minInt(xs))
		t := 0.0
		for _, x := range xs {
			if x < 0 {
				t += float64(x)
			}
		}
		return a*math.Min(w, 0) + b*t - gm*da, w, t
	}
	if ctx.KPOrdinal == 3 {
		k := edakit.Int(1, 2)
		order := allNodes(m)
		sort.Slice(order, func(i, j int) bool { return sl[order[i]] < sl[order[j]] })
		hide := map[int]bool{}
		for _, i := range order[:k] {
			hide[i] = true
		}
		var hidden []string
		var vis []int
		for i, x := range sl {
			if hide[i] {
				hidden = append(hidden, edakit.Sprintf("E%d", i+1))
				continue
			}
			vis = append(vis, x)
		}
		r, w, t := score(sl)
		rp, _, _ := score(vis)
		ans := edakit.Num(r, 2)
		return edakit.Numeric(formula+edakit.Sprintf("Signoff slacks after the agent's edit are, in ps: %s, and area grew %s%%. The edit also added set_false_path on %s, so the agent's own fast timer scored it %s. The grounded reward re-times with the original constraints. What is the grounded reward? Round to 2 decimal places.", strings.Join(rows, ", "), edakit.Num(da, 1), strings.Join(hidden, " and "), edakit.Num(rp, 2)), ans,
			edakit.Sprintf("With the original constraints every endpoint counts: WNS %s ps, TNS %s ps, so R = %s. The gap to %s is the reward hack that constraint edits buy.", edakit.Num(w, 0), edakit.Num(t, 0), ans, edakit.Num(rp, 2)), map[string]string{"reward": ans, "proxy": edakit.Num(rp, 2)})
	}
	r, w, t := score(sl)
	ans := edakit.Num(r, 2)
	return edakit.Numeric(formula+edakit.Sprintf("After the agent's edit the signoff endpoint slacks are, in ps: %s, and area grew %s%%. What is the reward? Round to 2 decimal places.", strings.Join(rows, ", "), edakit.Num(da, 1)), ans,
		edakit.Sprintf("WNS %s ps and TNS %s ps (negatives only), so R = %s.", edakit.Num(w, 0), edakit.Num(t, 0), ans), map[string]string{"wns": edakit.Num(w, 0), "tns": edakit.Num(t, 0), "reward": ans})
}
