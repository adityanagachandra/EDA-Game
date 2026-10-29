package agentic

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var harnessItems = []edakit.Item{
	{KP: 1, Q: "A request arrives: abut the regular cells of this datapath. What is the model's first move in the harness loop?",
		A: "Map the request to the smallest safe change, here a plug-in at the placer's extension point", D: []string{"Rewrite the legalizer to allow abutment", "Run the full regression suite", "Ask the human to write the change"},
		Why: "The model proposes the smallest safe edit; the harness then builds and tests it."},
	{KP: 1, Q: "In the harness loop, what happens immediately after the model proposes an edit?",
		A: "The harness writes the file and runs an incremental build", D: []string{"A human merges it", "The golden tests run on the unbuilt source", "The model runs the build itself on the host"},
		Why: "The harness owns every side effect: writing, building and testing."},
	{KP: 1, Q: "A candidate edit builds but fails one golden test. What does the harness do?",
		A: "Return the failure to the model as feedback and let it try again", D: []string{"Show the diff to a human for approval", "Promote the edit to the library", "Delete the golden test"},
		Why: "Failures loop back to the model; only passing changes reach human review."},
	{KP: 1, Q: "What does the human reviewer see when a change passes?",
		A: "A reviewable diff with its evidence: build log, tests passed and the QoR delta", D: []string{"Only the model's summary of the change", "The full engine source", "A notice that the change was merged"},
		Why: "Review is of the actual diff and the harness's own test results."},
	{KP: 2, Q: "The model writes in its reply that all tests pass. How should the harness treat this?",
		A: "As no evidence: only the harness's own test run counts", D: []string{"As sufficient for promotion", "As a reason to skip the differential", "As a Tier 1 change"},
		Why: "The model only proposes text. Evidence comes from the harness running the tests."},
	{KP: 2, Q: "Who may merge a passing change into the shared skill library?",
		A: "A human reviewer, after the harness has shown the evidence", D: []string{"The model, once tests pass", "The harness, automatically", "Any agent that calls the skill"},
		Why: "Acceptance belongs to a human; the model never gets the power to merge."},
	{KP: 2, Q: "Which component invokes the pinned compiler and enforces the sandbox's resource limits?",
		A: "The harness", D: []string{"The model", "The human reviewer", "The skill's manifest"},
		Why: "The harness owns builds, tests and limits; the model never touches the host directly."},
	{KP: 2, Q: "Why does the model never need direct write access to the shared repository?",
		A: "Because the harness writes files in a sandboxed workspace and a human promotes accepted changes", D: []string{"Because the model edits only manifests", "Because the repository is read-only for everyone", "Because the model works only on configuration"},
		Why: "The division of labor keeps every side effect under the harness's and the human's control."},
	{KP: 3, Q: "The model alternates between two edits: one fixes test A and breaks B, the other fixes B and breaks A. What should the harness do?",
		A: "Detect the oscillation, stop, and report both failures together", D: []string{"Keep retrying until the budget runs out", "Accept whichever edit passed more tests", "Merge both edits"},
		Why: "Oscillation means more retries will not converge; a human needs to see both failures."},
	{KP: 3, Q: "Working on a Tier 2 request, the model's edit touches the core legalization file. What should the harness do?",
		A: "Stop the automatic loop and turn it into a Tier 3 proposal for human review", D: []string{"Retry automatically, since it compiled", "Accept it if the goldens pass", "Silently revert the core part and keep the rest"},
		Why: "Edits outside the declared tier are not retried automatically; they need explicit approval."},
	{KP: 3, Q: "A request exhausts its attempt budget without passing. What should the harness report?",
		A: "The best attempt and its failures, so an engineer can pick it up", D: []string{"Nothing, to avoid noise", "Only that the request failed", "The first attempt, whatever its result"},
		Why: "A failed loop still produces useful partial work and diagnostics."},
	{KP: 3, Q: "Why does every request get both an attempt budget and a wall-time budget?",
		A: "Attempts bound how often the model retries; wall time bounds slow builds and tests", D: []string{"Because builds never fail", "So the model can choose which to ignore", "Because golden tests have no runtime"},
		Why: "Either can run away on its own; the loop needs both stopping rules."},
}

func buildBudget(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		u := edakit.Scaled(ctx, 4, 30)
		tc := float64(edakit.Int(3, 12))
		tl := float64(edakit.Int(2, 8))
		hdr := float64(u)*tc + tl
		q := edakit.Sprintf("An edit changes a header included by %d translation units. Each unit takes %s s to compile and linking takes %s s. How many seconds does the rebuild take?", u, edakit.Num(tc, 0), edakit.Num(tl, 0))
		exp := edakit.Sprintf("$%d \\times %s + %s = %s$ s. Making the same change inside the plug-in's own source file would cost %s s.", u, edakit.Num(tc, 0), edakit.Num(tl, 0), edakit.Num(hdr, 0), edakit.Num(tc+tl, 0))
		return edakit.Numeric(q, edakit.Num(hdr, 0), exp, map[string]string{"plugin_only": edakit.Num(tc+tl, 0)})
	case 2:
		n := edakit.Scaled(ctx, 20, 120)
		c := edakit.Int(1, 4)
		h := float64(edakit.Int(50, 98)) / 100
		th := float64(edakit.Int(1, 5)) / 10
		tm := float64(edakit.Int(3, 10))
		tl := float64(edakit.Int(2, 6))
		e := float64(c)*tm + float64(n-c)*(h*th+(1-h)*tm) + tl
		q := edakit.Sprintf("A fresh sandbox rebuilds a skill of %d translation units with a compiler cache. %d units changed and always miss; the other units hit the cache with probability %s. A hit costs %s s, a miss %s s, and linking %s s. What is the expected build time in seconds? Round to 2 decimal places.",
			n, c, edakit.Num(h, 2), edakit.Num(th, 1), edakit.Num(tm, 0), edakit.Num(tl, 0))
		exp := edakit.Sprintf("$E[T] = %d \\times %s + %d \\times (%s \\times %s + %s \\times %s) + %s = %s$ s.", c, edakit.Num(tm, 0), n-c, edakit.Num(h, 2), edakit.Num(th, 1), edakit.Num(1-h, 2), edakit.Num(tm, 0), edakit.Num(tl, 0), edakit.Num(e, 2))
		return edakit.Numeric(q, edakit.Num(e, 2), exp, map[string]string{"expected_build": edakit.Num(e, 2)})
	}
	mt := edakit.Int(10, 40)
	bt := edakit.Int(5, 90)
	tt := edakit.Int(10, 60)
	it := mt + bt + tt
	n := 3600 / it
	q := edakit.Sprintf("Each harness attempt spends %d s waiting for the model, %d s building and %d s running golden tests. How many complete attempts fit in one hour?", mt, bt, tt)
	exp := edakit.Sprintf("One attempt takes $%d + %d + %d = %d$ s, and $\\lfloor 3600 / %d \\rfloor = %d$.", mt, bt, tt, it, it, n)
	return edakit.Numeric(q, edakit.IntStr(n), exp, map[string]string{"attempt_seconds": edakit.IntStr(it)})
}

func retryExpectation(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		p := float64(edakit.Int(15, 60)) / 100
		if edakit.Int(0, 1) == 0 {
			k := edakit.Scaled(ctx, 2, 8)
			v := 1 - math.Pow(1-p, float64(k))
			q := edakit.Sprintf("Each harness attempt at a skill edit succeeds independently with probability %s. What is the probability that at least one of %d attempts succeeds? Round to 4 decimal places.", edakit.Num(p, 2), k)
			exp := edakit.Sprintf("$1 - (1 - %s)^{%d} = %s$.", edakit.Num(p, 2), k, edakit.Num(v, 4))
			return edakit.Numeric(q, edakit.Num(v, 4), exp, map[string]string{"all_fail": edakit.Num(1-v, 4)})
		}
		target := edakit.Pick([]float64{0.9, 0.95, 0.99})
		k := math.Ceil(math.Log(1-target) / math.Log(1-p))
		q := edakit.Sprintf("Each harness attempt succeeds independently with probability %s. What is the smallest attempt budget that gives at least a %s chance of success?", edakit.Num(p, 2), edakit.Num(target, 2))
		exp := edakit.Sprintf("Require $(1-%s)^k \\le %s$, so $k \\ge \\ln %s / \\ln %s = %s$; round up to %s.", edakit.Num(p, 2), edakit.Num(1-target, 2), edakit.Num(1-target, 2), edakit.Num(1-p, 2), edakit.Num(math.Log(1-target)/math.Log(1-p), 2), edakit.Num(k, 0))
		return edakit.Numeric(q, edakit.Num(k, 0), exp, nil)
	case 2:
		p := float64(edakit.Int(15, 60)) / 100
		k := edakit.Scaled(ctx, 3, 8)
		e := (1 - math.Pow(1-p, float64(k))) / p
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("Each harness attempt succeeds independently with probability %s, and the harness stops at the first success or after %d attempts. What is the expected number of attempts? Round to 4 decimal places.", edakit.Num(p, 2), k)
			exp := edakit.Sprintf("$E[\\min(G,k)] = (1 - %s^{%d}) / %s = %s$.", edakit.Num(1-p, 2), k, edakit.Num(p, 2), edakit.Num(e, 4))
			return edakit.Numeric(q, edakit.Num(e, 4), exp, nil)
		}
		t := 10 * edakit.Int(4, 18)
		q := edakit.Sprintf("Each harness attempt takes %d s and succeeds independently with probability %s. The harness stops at the first success or after %d attempts. What is the expected time per request in seconds? Round to 1 decimal place.", t, edakit.Num(p, 2), k)
		exp := edakit.Sprintf("Expected attempts $(1 - %s^{%d}) / %s = %s$, times %d s gives %s s.", edakit.Num(1-p, 2), k, edakit.Num(p, 2), edakit.Num(e, 4), t, edakit.Num(e*float64(t), 1))
		return edakit.Numeric(q, edakit.Num(e*float64(t), 1), exp, nil)
	}
	n := edakit.Scaled(ctx, 2, 5)
	var ps []string
	prod := 1.0
	p := float64(edakit.Int(15, 40)) / 100
	for i := 0; i < n; i++ {
		ps = append(ps, edakit.Num(p, 2))
		prod *= 1 - p
		p = math.Min(0.9, p+float64(edakit.Int(5, 20))/100)
	}
	v := 1 - prod
	q := edakit.Sprintf("Feedback makes later attempts likelier to succeed: the success probabilities of attempts 1 to %d, given that earlier ones failed, are %s. What is the probability that the request succeeds within %d attempts? Round to 4 decimal places.", n, strings.Join(ps, ", "), n)
	exp := edakit.Sprintf("Multiply the failure probabilities to get %s, and subtract from 1: %s.", edakit.Num(prod, 4), edakit.Num(v, 4))
	return edakit.Numeric(q, edakit.Num(v, 4), exp, nil)
}

var feedbackItems = []edakit.Item{
	{KP: 1, Q: "A plug-in fails to compile with 40 errors. What should the harness send back first?",
		A: "The first diagnostic verbatim, with file, line and error class", D: []string{"The last error, since it is the most recent", "A count of the errors only", "The entire build log"},
		Why: "Later errors are often consequences of the first; the exact first diagnostic is what the model can fix."},
	{KP: 1, Q: "A plug-in's method signature does not match the extension-point interface. What extra context helps the model most?",
		A: "The interface declaration the error refers to", D: []string{"The full engine source", "The golden test outputs", "The model's previous prompt"},
		Why: "Showing the expected signature makes the fix obvious."},
	{KP: 1, Q: "AddressSanitizer reports a heap overflow with a 60-frame stack trace. How should the harness return it?",
		A: "Trimmed to the frames in the skill's own code, with the error type", D: []string{"As a pass, since the build succeeded", "With all 60 frames and the memory map", "As a single word, 'crash'"},
		Why: "The model can act only on the frames it can change."},
	{KP: 1, Q: "Why is returning a 2,000-line build log a poor choice?",
		A: "It buries the error, wastes context and makes the model guess which line matters", D: []string{"Logs are always wrong", "The model cannot read logs", "Long logs make the build slower"},
		Why: "Precise, small feedback leads to better next attempts."},
	{KP: 2, Q: "A golden placement test fails. Which feedback is most actionable?",
		A: "Test dp_alu32: displacement 14.2 sites, limit 8", D: []string{"Golden tests failed", "The placement looks wrong", "Exit code 3"},
		Why: "Naming the test, metric, actual value and limit tells the model what to fix."},
	{KP: 2, Q: "A differential timing check fails. What should the feedback list?",
		A: "The endpoints that moved unexpectedly, with before and after slack", D: []string{"Every endpoint in the design", "Only the worst negative slack", "The SDC file"},
		Why: "The unexpected movers are the evidence of collateral damage."},
	{KP: 2, Q: "Why should test feedback also say which checks passed?",
		A: "So the model does not undo a correct part of its edit", D: []string{"To make the report longer", "Because passing checks are optional", "So the model can delete them"},
		Why: "Knowing what works keeps the next attempt focused."},
	{KP: 2, Q: "Ten downstream checks fail because one upstream check failed. What should the feedback lead with?",
		A: "The first failing check, the likely cause of the rest", D: []string{"All ten failures in random order", "The last failing check", "Only a count of failures"},
		Why: "Consequences distract from the cause."},
	{KP: 3, Q: "The harness returns each golden test's full expected output. What can a model learn to do?",
		A: "Special-case those inputs and pass without solving the problem", D: []string{"Write faster code", "Fix the build system", "Nothing, since expected outputs are harmless"},
		Why: "Exposed answers invite hard-coding; hidden tests catch it."},
	{KP: 3, Q: "A golden run times out because the test machine was overloaded, and the harness reports it as the model's failure. What goes wrong?",
		A: "The model learns to change correct code to fix a problem it did not cause", D: []string{"Nothing: timeouts are always the model's fault", "The golden test becomes stricter", "The build cache is cleared"},
		Why: "Infrastructure errors must be separated from skill errors."},
	{KP: 3, Q: "The harness reports only pass or fail with no detail. How do models typically respond?",
		A: "With near-random changes, because they cannot tell what to fix", D: []string{"With smaller, better edits", "By asking for more tests", "By converging faster"},
		Why: "Without a located failure, the next attempt is guesswork."},
	{KP: 3, Q: "Why should some tests stay hidden from the model?",
		A: "To catch edits that pass the visible tests without solving the task", D: []string{"To save compute", "Because hidden tests are faster", "To prevent the model from building"},
		Why: "Hidden tests measure whether the change generalizes beyond what the model saw."},
}

type stage struct {
	name string
	c, q float64
}

func expCost(st []stage) float64 {
	tot, pass := 0.0, 1.0
	for _, s := range st {
		tot += pass * s.c
		pass *= s.q
	}
	return tot
}

func checkOrdering(ctx generator.GeneratorContext) generator.Problem {
	pool := []string{"lint", "unit tests", "sanitizer build", "golden tests", "differential QoR run", "format check"}
	n := 3
	if ctx.Difficulty > 0.7 {
		n = 4
	}
	names := edakit.Shuffle(pool)[:n]
	var st []stage
	for _, nm := range names {
		st = append(st, stage{nm, float64(5 * edakit.Int(1, 60)), float64(edakit.Int(55, 95)) / 100})
	}
	opt := append([]stage(nil), st...)
	sort.Slice(opt, func(i, j int) bool { return opt[i].c/(1-opt[i].q) < opt[j].c/(1-opt[j].q) })
	if math.Abs(expCost(st)-expCost(opt)) < 1e-9 {
		for i, j := 0, len(st)-1; i < j; i, j = i+1, j-1 {
			st[i], st[j] = st[j], st[i]
		}
	}
	given := expCost(st)
	best := expCost(opt)
	var desc []string
	for _, s := range st {
		desc = append(desc, edakit.Sprintf("%s (%s s, passes with probability %s)", s.name, edakit.Num(s.c, 0), edakit.Num(s.q, 2)))
	}
	var order []string
	for _, s := range opt {
		order = append(order, s.name)
	}
	head := edakit.Sprintf("A harness runs independent checks on each candidate and stops at the first failure: %s. ", strings.Join(desc, ", "))
	switch kp(ctx) {
	case 1:
		q := head + "Run in the order listed, what is the expected cost per candidate in seconds? Round to 2 decimal places."
		exp := edakit.Sprintf("Each stage's cost is weighted by the probability that all earlier stages passed; the sum is %s s.", edakit.Num(given, 2))
		return edakit.Numeric(q, edakit.Num(given, 2), exp, nil)
	case 2:
		q := head + "If the checks are run in the order that minimizes expected cost, what is the expected cost per candidate in seconds? Round to 2 decimal places."
		exp := edakit.Sprintf("Sort by cost over failure probability: %s. The expected cost in that order is %s s.", strings.Join(order, ", "), edakit.Num(best, 2))
		return edakit.Numeric(q, edakit.Num(best, 2), exp, nil)
	}
	q := head + "How many seconds per candidate does the cost-minimizing order save compared with the order listed? Round to 2 decimal places."
	exp := edakit.Sprintf("Listed order: %s s. Best order (%s): %s s. Saving %s s.", edakit.Num(given, 2), strings.Join(order, ", "), edakit.Num(best, 2), edakit.Num(given-best, 2))
	return edakit.Numeric(q, edakit.Num(given-best, 2), exp, nil)
}
