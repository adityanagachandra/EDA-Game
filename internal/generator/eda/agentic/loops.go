package agentic

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var promotionItems = []edakit.Item{
	{KP: 1, Q: "What are the four steps of the skill flywheel, in order?",
		A: "Identify, iterate fast as a plug-in, validate, promote", D: []string{"Promote, validate, iterate, identify", "Train, deploy, monitor, retire", "Identify, promote, iterate, validate"},
		Why: "A skill is iterated and validated before it is promoted into the library."},
	{KP: 1, Q: "Who can identify a new skill in the flywheel?",
		A: "The model, an agent, a field engineer or a customer", D: []string{"Only the engine's original author", "Only the vendor's research team", "Only the model"},
		Why: "Discovery can come from anyone working on a design."},
	{KP: 1, Q: "Why is a new behavior first iterated as a runtime-loaded plug-in?",
		A: "Each attempt rebuilds only the plug-in in seconds and a broken one is simply unloaded", D: []string{"Because plug-ins skip testing", "Because plug-ins run faster than compiled code", "Because the core forbids static libraries"},
		Why: "Fast iteration needs isolation from the core build."},
	{KP: 1, Q: "A new abutment rule improves wirelength on one run. Is it ready to promote?",
		A: "No: it must pass golden and differential tests and hold its gain across seeds", D: []string{"Yes: one improvement is enough", "Yes, if the model is confident", "No: abutment rules can never be promoted"},
		Why: "Validation comes before promotion."},
	{KP: 2, Q: "What changes when a validated plug-in is promoted?",
		A: "It is compiled into the static skill library with a manifest, shippable goldens and a version", D: []string{"It is deleted and rewritten by hand", "It moves into the core", "It stays a runtime plug-in forever"},
		Why: "Promotion makes it a first-class library skill selectable at full speed."},
	{KP: 2, Q: "Why compile a promoted skill statically rather than keep loading it at runtime?",
		A: "The compiler can optimize across it and every deployment can select it at full speed", D: []string{"Static code cannot have bugs", "Runtime loading is forbidden after promotion for legal reasons", "Static skills need no manifest"},
		Why: "Promotion trades flexibility for speed once the skill is proven."},
	{KP: 2, Q: "A plug-in works on one customer's design but cannot be validated anywhere else. What happens?",
		A: "It stays a local plug-in at that customer", D: []string{"It is promoted with the customer's design as its golden", "It is promoted without goldens", "It is merged into the core"},
		Why: "Promotion requires validation that can leave the site."},
	{KP: 2, Q: "What must a promoted skill have that a quick plug-in did not need?",
		A: "A manifest with applicability and limits, goldens on vendor-held designs and owner review", D: []string{"A GUI", "A trained model", "A customer's signoff report"},
		Why: "Library skills are selected by manifest and must be testable anywhere."},
	{KP: 3, Q: "Where does a promoted skill land?",
		A: "In the shared skill library and in the domain model's knowledge of when to pick it", D: []string{"Only in the customer's repository", "Only in the vendor's archive", "In the engine core"},
		Why: "The library makes it selectable; the brain learns when to use it."},
	{KP: 3, Q: "What flows back to the vendor when a skill is promoted from a customer site?",
		A: "The skill's code, manifest, tests and know-how about its use", D: []string{"The design it was developed on", "The customer's constraints", "The customer's process data"},
		Why: "Know-how flows back; the design never does."},
	{KP: 3, Q: "Why does the flywheel get stronger with use?",
		A: "Each turn adds a reusable skill and teaches the selection model when to use it", D: []string{"Each turn makes the core larger", "Each turn collects more customer designs", "Each turn removes old skills"},
		Why: "Library growth and better selection compound across designs."},
	{KP: 3, Q: "Why is the flywheel acceptable to customers?",
		A: "Their discoveries improve the library without exposing their designs", D: []string{"Because they receive the vendor's core source", "Because promotion is optional for the vendor", "Because skills are free"},
		Why: "The IP boundary holds while know-how is shared."},
}

func amortization(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		i := edakit.Int(5, 20)
		c := float64(edakit.Int(1, 6)) / 2
		r := 50 * edakit.Int(2, 20)
		v := float64(i) * c / float64(r) * 60
		q := edakit.Sprintf("Setting up an agentic flow takes %d iterations of %s hours each. The flow then serves %d production runs. How many minutes of setup does each run carry? Round to 2 decimal places.", i, edakit.Num(c, 1), r)
		exp := edakit.Sprintf("$%d \\times %s / %d = %s$ h per run, or %s minutes.", i, edakit.Num(c, 1), r, edakit.Num(float64(i)*c/float64(r), 4), edakit.Num(v, 2))
		return edakit.Numeric(q, edakit.Num(v, 2), exp, nil)
	case 2:
		u := edakit.Int(2, 15)
		uc := float64(edakit.Int(2, 12))
		l := float64(250 * edakit.Int(2, 16))
		f := 100 * float64(u) * uc / (float64(u)*uc + l)
		q := edakit.Sprintf("A deployment makes %d skill updates a month at %s compute-hours each, and its modeling loop and signoff runs consume %s compute-hours a month. What percentage of total compute goes to skill updates? Round to 2 decimal places.", u, edakit.Num(uc, 0), edakit.Num(l, 0))
		exp := edakit.Sprintf("$%s / (%s + %s) = %s\\%%$.", edakit.Num(float64(u)*uc, 0), edakit.Num(float64(u)*uc, 0), edakit.Num(l, 0), edakit.Num(f, 2))
		return edakit.Numeric(q, edakit.Num(f, 2), exp, nil)
	}
	uc := float64(edakit.Int(2, 20))
	s := float64(edakit.Int(1, 9)) / 10
	n := math.Ceil(uc/s - 1e-9)
	q := edakit.Sprintf("A skill update costs %s hours of compute and makes each subsequent run %s hours faster. After how many runs does the update break even?", edakit.Num(uc, 0), edakit.Num(s, 1))
	exp := edakit.Sprintf("$\\lceil %s / %s \\rceil = %s$ runs.", edakit.Num(uc, 0), edakit.Num(s, 1), edakit.Num(n, 0))
	return edakit.Numeric(q, edakit.Num(n, 0), exp, nil)
}

var brainItems = []edakit.Item{
	{KP: 1, Q: "What does a domain brain model learn?",
		A: "A mapping from situation to action: which skill, parameters, customization or recursion fits", D: []string{"How to replace the engines' algorithms", "The customer's designs verbatim", "How to legalize placements"},
		Why: "The brain learns when to pick each skill, not how to do the engines' work."},
	{KP: 1, Q: "Which record is a training example for the brain?",
		A: "Design traits and loop metrics, the skill chosen, and the engine-measured QoR outcome", D: []string{"The full GDS of the design", "The model's prompt only", "A list of skill names"},
		Why: "Situation, action and outcome form the learning signal."},
	{KP: 1, Q: "Where does the brain's outcome signal come from?",
		A: "The engines, which measure the QoR each choice produced", D: []string{"The model's own estimate", "A human rating each run", "The manifest's claims"},
		Why: "Outcomes are engine-measured, which keeps the signal objective."},
	{KP: 1, Q: "How do promoted skills change what the brain can do?",
		A: "They add new actions it can choose, and their validation runs add outcomes", D: []string{"They replace the brain", "They are ignored", "They remove old training data"},
		Why: "The flywheel feeds the brain."},
	{KP: 2, Q: "Which data may train a shared base brain model at the vendor?",
		A: "Promoted skills, their manifests and tests, and outcomes on public or synthetic designs", D: []string{"Every customer's design-derived features", "Customer timing constraints", "Customer process data"},
		Why: "Only what may leave a site can train the shared model."},
	{KP: 2, Q: "Where is a per-customer brain fine-tuned?",
		A: "On-site, from that customer's own records, and it stays there", D: []string{"At the vendor, from uploaded records", "On a public cloud", "On another customer's site"},
		Why: "Design-derived training data stays inside the boundary."},
	{KP: 2, Q: "Why does a per-customer brain need a model with controllable weights?",
		A: "Fine-tuning and hosting inside the customer's environment is impossible with a remote-only API model", D: []string{"API models cannot answer questions", "Controllable weights are always larger", "Open models need no hardware"},
		Why: "On-site fine-tuning requires the weights."},
	{KP: 3, Q: "A brain trained mostly on one customer's designs picks the same datapath skill for every design with a few adders. What is this?",
		A: "Overfitting to the training distribution", D: []string{"Good generalization", "Reward hacking by the engines", "A determinism bug"},
		Why: "Held-out evaluation on unlike designs catches it."},
	{KP: 3, Q: "A new process node arrives and the brain chooses confidently, but outcomes are poor. What failure is this?",
		A: "Distribution shift: the situation lies outside what it has seen", D: []string{"A frozen ABI", "Skill promotion", "Seed noise"},
		Why: "An uncertainty estimate should trigger exploration or recursion here."},
	{KP: 3, Q: "A skill the brain never picks never collects outcome data. What does this create?",
		A: "A feedback loop that can hide good skills; a fraction of exploratory choices fixes it", D: []string{"Perfect selection", "A smaller library", "Faster training"},
		Why: "Learning only from its own choices biases a selector."},
	{KP: 3, Q: "How should a brain's selection quality be evaluated?",
		A: "On held-out designs unlike its training set, against measured outcomes", D: []string{"On its training designs", "By the model's confidence", "By the number of skills it knows"},
		Why: "Evaluation must test generalization."},
	{KP: 3, Q: "Why should the brain report uncertainty along with its choice?",
		A: "So low-confidence situations trigger exploration or recursion to accurate engines", D: []string{"To slow selection down", "Because manifests require it", "So humans can ignore it"},
		Why: "Calibrated uncertainty tells the loop when not to trust the choice."},
}

var fastAccItems = []edakit.Item{
	{KP: 1, Q: "What does a fast congestion estimator trade away compared with global routing?",
		A: "It estimates routing demand probabilistically or by patterns instead of routing", D: []string{"It ignores the design's nets", "It routes every net twice", "It skips legality of placement"},
		Why: "Each fast mode trades a specific piece of fidelity for speed."},
	{KP: 1, Q: "A fast timer uses Elmore delay on estimated trees and skips crosstalk. When is it trustworthy?",
		A: "Within the regime where its error against signoff has been measured", D: []string{"Always, since Elmore is exact", "Never", "Only on clock nets"},
		Why: "A fast engine is useful only where it has been calibrated."},
	{KP: 1, Q: "Which is a legitimate way to build a fast mode of an engine?",
		A: "A surrogate model trained on accurate engine results", D: []string{"Removing legality checks", "Reporting last week's results", "Skipping all computation"},
		Why: "Surrogates are one construction of fast modes."},
	{KP: 1, Q: "Which pair is a fast and accurate version of the same analysis?",
		A: "A pattern-table extractor and a field solver", D: []string{"A placer and a router", "A timer and a DRC checker", "Synthesis and simulation"},
		Why: "Both compute parasitics, at different fidelity and cost."},
	{KP: 1, Q: "A fast placer stops global placement at a looser overflow and skips detailed placement. What does the agent gain?",
		A: "Many more place-route-analyze iterations per day", D: []string{"A signoff-quality placement", "A guarantee of routability", "Exact timing"},
		Why: "Speed buys exploration; the accurate engine confirms."},
	{KP: 3, Q: "Fast-loop WNS improves steadily with low variance across seeds. What should the agent do?",
		A: "Continue with the current skill", D: []string{"Recurse to signoff on every endpoint", "Switch skills immediately", "Stop the flow"},
		Why: "Stable improvement means the fast loop is informative and the skill is working."},
	{KP: 3, Q: "Fast-loop congestion has plateaued for many iterations. What is the next move?",
		A: "Try another skill or customize the current one", D: []string{"Keep iterating unchanged", "Delete the congestion estimator", "Lower the density target to zero"},
		Why: "Plateaus signal that the current strategy has stopped helping."},
	{KP: 3, Q: "The variance of fast-timer WNS across seeds suddenly grows. What should the agent do?",
		A: "Recurse to the accurate engine, because the fast numbers have stopped being informative", D: []string{"Trust the best seed", "Average the seeds and continue", "Ignore variance"},
		Why: "Growing variance is a signal to recurse."},
	{KP: 3, Q: "The design moves to a new layer stack on which the fast extractor was never correlated. What should happen?",
		A: "Recurse to the accurate extractor more often or recalibrate", D: []string{"Trust the fast extractor as before", "Stop extracting", "Use the old layer stack's results"},
		Why: "Outside the calibrated regime, fast results are unearned."},
	{KP: 3, Q: "Why must a fast engine report its uncertainty as well as its estimate?",
		A: "Because the agent decides where to recurse from the uncertainty", D: []string{"So reports are longer", "Because the accurate engine needs it", "To satisfy the manifest format"},
		Why: "The recursion decision is made from data."},
}

func loopBudget(ctx generator.GeneratorContext) generator.Problem {
	tf := float64(edakit.Int(3, 15))
	ta := float64(15 * edakit.Int(3, 12))
	k := edakit.Pick([]int{4, 5, 8, 10, 12})
	b := float64(60 * edakit.Int(6, 16))
	n := 0
	for float64(n+1)*tf+math.Floor(float64(n+1)/float64(k))*ta <= b {
		n++
	}
	naive := math.Floor(b / tf)
	allAcc := math.Floor(b / ta)
	noFloor := math.Floor(b / (tf + ta/float64(k)))
	q := edakit.Sprintf("An agent runs a fast loop whose iterations take %s minutes and confirms with an accurate run of %s minutes after every %d fast iterations. How many fast iterations fit in a budget of %s minutes?", edakit.Num(tf, 0), edakit.Num(ta, 0), k, edakit.Num(b, 0))
	exp := edakit.Sprintf("The largest $N$ with $N \\times %s + \\lfloor N/%d \\rfloor \\times %s \\le %s$ is %d.", edakit.Num(tf, 0), k, edakit.Num(ta, 0), edakit.Num(b, 0), n)
	return choice(q, float64(n), 0, "", exp, 1, naive, allAcc, noFloor, float64(n+1))
}

type pair struct{ x, y float64 }

func calibration(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	bias := float64(edakit.Int(-6, 8))
	slope := edakit.Pick([]float64{0.9, 1, 1.1, 1.2})
	ps := make([]pair, n)
	used := map[int]bool{}
	for i := range ps {
		x := 5 * edakit.Int(-8, 12)
		for t := 0; used[x] && t < generator.MaxDrawAttempts; t++ {
			x = 5 * edakit.Int(-8, 12)
		}
		for used[x] {
			x++
		}
		used[x] = true
		y := math.Round(slope*float64(x) - bias + float64(edakit.Int(-6, 6)))
		ps[i] = pair{float64(x), y}
	}
	var rows []string
	for _, p := range ps {
		rows = append(rows, edakit.Sprintf("(%s, %s)", edakit.Num(p.x, 0), edakit.Num(p.y, 0)))
	}
	head := edakit.Sprintf("Paired slacks in ps, written as (fast timer, signoff timer), on %d endpoints: %s. ", n, strings.Join(rows, ", "))
	switch kp(ctx) {
	case 1:
		var es []float64
		mean := 0.0
		for _, p := range ps {
			es = append(es, p.x-p.y)
			mean += p.x - p.y
		}
		mean /= float64(n)
		ss := 0.0
		for _, e := range es {
			ss += (e - mean) * (e - mean)
		}
		sd := math.Sqrt(ss / float64(n-1))
		if edakit.Int(0, 1) == 0 {
			q := head + "What is the bias of the fast timer, the mean of fast minus signoff, in ps? Round to 2 decimal places."
			return edakit.Numeric(q, edakit.Num(mean, 2), edakit.Sprintf("The errors fast minus signoff average to %s ps.", edakit.Num(mean, 2)), nil)
		}
		q := head + "What is the sample standard deviation (dividing by $n-1$) of the error fast minus signoff, in ps? Round to 2 decimal places."
		return edakit.Numeric(q, edakit.Num(sd, 2), edakit.Sprintf("Errors have mean %s ps; the sum of squared deviations divided by %d gives a standard deviation of %s ps.", edakit.Num(mean, 2), n-1, edakit.Num(sd, 2)), nil)
	case 2:
		mx, my := 0.0, 0.0
		for _, p := range ps {
			mx += p.x
			my += p.y
		}
		mx /= float64(n)
		my /= float64(n)
		sxx, sxy := 0.0, 0.0
		for _, p := range ps {
			sxx += (p.x - mx) * (p.x - mx)
			sxy += (p.x - mx) * (p.y - my)
		}
		b := sxy / sxx
		a := my - b*mx
		x0 := float64(5 * edakit.Int(-4, 6))
		pred := a + b*x0
		q := head + edakit.Sprintf("Fit signoff $\\approx a + b \\cdot$ fast by least squares. What signoff slack does the fit predict for a fast slack of %s ps? Round to 2 decimal places.", edakit.Num(x0, 0))
		exp := edakit.Sprintf("$\\bar x = %s$, $\\bar y = %s$, $S_{xx} = %s$, $S_{xy} = %s$, so $b = %s$ and $a = %s$; the prediction is %s ps.", edakit.Num(mx, 2), edakit.Num(my, 2), edakit.Num(sxx, 2), edakit.Num(sxy, 2), edakit.Num(b, 4), edakit.Num(a, 4), edakit.Num(pred, 2))
		return edakit.Numeric(q, edakit.Num(pred, 2), exp, nil)
	}
	conc, tot := 0, 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := ps[i].x - ps[j].x
			dy := ps[i].y - ps[j].y
			if dy == 0 {
				continue
			}
			tot++
			if dx*dy > 0 {
				conc++
			}
		}
	}
	q := head + "Ignoring any pair tied in signoff slack, how many pairs of endpoints are ordered the same way by both timers?"
	exp := edakit.Sprintf("%d of the %d untied pairs are concordant.", conc, tot)
	return edakit.Numeric(q, edakit.IntStr(conc), exp, nil)
}

func recursion(ctx generator.GeneratorContext) generator.Problem {
	b := float64(edakit.Int(0, 8))
	sd := float64(edakit.Int(3, 12))
	n := edakit.Scaled(ctx, 5, 8)
	var ms []float64
	used := map[int]bool{}
	for t := 0; len(ms) < n; t++ {
		v := edakit.Int(-40, 50)
		if t > generator.MaxDrawAttempts {
			v = 60 + len(ms)
		}
		if used[v] {
			continue
		}
		used[v] = true
		ms = append(ms, float64(v))
	}
	sort.Float64s(ms)
	var vs []string
	for _, m := range ms {
		vs = append(vs, edakit.Num(m, 0))
	}
	head := edakit.Sprintf("A fast timer is calibrated against signoff: its slack is optimistic by a bias of %s ps, with an error standard deviation of %s ps (approximately normal). ", edakit.Num(b, 0), edakit.Num(sd, 0))
	zs := make([]float64, n)
	for i, m := range ms {
		zs[i] = (m - b) / sd
	}
	switch kp(ctx) {
	case 1:
		m := math.Round(b + sd*float64(edakit.Int(-10, 10))/4)
		z := (m - b) / sd
		p := phi(-z)
		q := head + edakit.Sprintf("An endpoint has a fast slack of %s ps. What is the probability that its signoff slack is negative? Round to 4 decimal places.", edakit.Num(m, 0))
		exp := edakit.Sprintf("$z = (%s - %s)/%s = %s$ and $\\Phi(-z) = %s$.", edakit.Num(m, 0), edakit.Num(b, 0), edakit.Num(sd, 0), edakit.Num(z, 4), edakit.Num(p, 4))
		return edakit.Numeric(q, edakit.Num(p, 4), exp, nil)
	case 2:
		k := edakit.Pick([]float64{1.5, 2, 2.5, 3})
		cnt := 0
		for _, z := range zs {
			if math.Abs(z) < k {
				cnt++
			}
		}
		q := head + edakit.Sprintf("The agent recurses to signoff on every endpoint whose bias-corrected z-score has absolute value below %s. Fast slacks in ps: %s. How many endpoints need the accurate engine?", edakit.Num(k, 1), strings.Join(vs, ", "))
		var zt []string
		for _, z := range zs {
			zt = append(zt, edakit.Num(z, 2))
		}
		exp := edakit.Sprintf("z-scores $(m - %s)/%s$: %s. %d have $|z| < %s$.", edakit.Num(b, 0), edakit.Num(sd, 0), strings.Join(zt, ", "), cnt, edakit.Num(k, 1))
		return edakit.Numeric(q, edakit.IntStr(cnt), exp, nil)
	}
	r := edakit.Int(1, 3)
	abs := make([]float64, n)
	for i, z := range zs {
		abs[i] = math.Abs(z)
	}
	sort.Float64s(abs)
	sum := 0.0
	for _, a := range abs[r:] {
		sum += phi(-a)
	}
	q := head + edakit.Sprintf("Fast slacks in ps: %s. The accurate engine can run on only %d endpoints, chosen as those with the smallest absolute z-score; the fast verdict (pass if corrected slack is positive) is accepted for the rest. What is the expected number of wrong verdicts among the endpoints left to the fast timer? Round to 3 decimal places.", strings.Join(vs, ", "), r)
	exp := edakit.Sprintf("Sum $\\Phi(-|z|)$ over the endpoints not sent to signoff: %s.", edakit.Num(sum, 3))
	return edakit.Numeric(q, edakit.Num(sum, 3), exp, nil)
}
