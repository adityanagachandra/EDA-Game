package agentic

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var goldenItems = []edakit.Item{
	{KP: 1, Q: "What is a golden test for a placer?",
		A: "A fixed small design with a verified expected placement result and legality checks", D: []string{"A random design run once to see what happens", "The largest customer design available", "A test that always passes"},
		Why: "Goldens pair fixed inputs with known-good outputs and invariants."},
	{KP: 1, Q: "A heuristic router is not bit-for-bit reproducible across machines. How can a golden still check it?",
		A: "Check invariants exactly and metrics within stated tolerances", D: []string{"It cannot be tested", "Compare logs character by character", "Accept any result that finishes"},
		Why: "Invariants (DRC-clean, connected) are exact; metrics get tolerances."},
	{KP: 1, Q: "Why are golden tests kept small?",
		A: "So they run in seconds inside the harness loop", D: []string{"So they never fail", "Because large designs are not legal", "So they can be stored in the manifest"},
		Why: "Goldens run on every attempt, so their runtime bounds iteration speed."},
	{KP: 1, Q: "What question does a golden suite answer?",
		A: "Is the engine still correct on cases we understand?", D: []string{"Is the new skill better on every design?", "How fast is the model?", "Which customer should use the skill?"},
		Why: "Goldens guard known-correct behavior; differentials compare against a baseline."},
	{KP: 3, Q: "A new derate provider passes goldens and a differential, but every reference design has logic depth above 5 and the provider's rule fires below 5. What has been proved?",
		A: "Nothing about the rule itself: its condition was never exercised", D: []string{"That the rule is correct", "That the rule never fires", "That the core is unchanged"},
		Why: "Reference designs must exercise each condition of the change."},
	{KP: 3, Q: "The baseline itself changes slightly from run to run with the same seed. What happens to differential checks?",
		A: "They show movement everywhere, so only loose tolerances can pass and real regressions hide", D: []string{"They become more sensitive", "Nothing: differentials ignore noise", "They run faster"},
		Why: "Differentials depend on a deterministic baseline."},
	{KP: 3, Q: "A bug escapes to a user although all goldens passed. What is the first fix to the test suite?",
		A: "Add a golden that would have caught it", D: []string{"Loosen the tolerances", "Remove the failing design", "Stop running goldens"},
		Why: "Each escaped bug should become a permanent test."},
	{KP: 3, Q: "Together, what do golden and differential tests prove about a change?",
		A: "It did what was intended on the reference set and broke nothing visible there", D: []string{"It is correct on every possible design", "It improves every metric", "It is faster than the baseline"},
		Why: "The proof is bounded by the reference set."},
	{KP: 2, Q: "A differential check classifies each object's change. Which classes cause a failure?",
		A: "Unintended movers and intended movers off by more than the tolerance", D: []string{"Only objects whose metric got worse", "Only intended movers", "None, if the aggregate QoR improved"},
		Why: "Pass requires zero unintended and zero wrong-amount movers."},
	{KP: 2, Q: "What must be the same between the baseline and candidate runs of a differential check?",
		A: "The inputs, configuration and seeds, everything except the change under test", D: []string{"Only the design name", "Only the engine version", "Nothing, since differences are expected"},
		Why: "Any other difference would be confused with the effect of the change."},
}

func goldenDiffCalc(ctx generator.GeneratorContext) generator.Problem {
	kinds := []struct{ obj, metric, unit string }{
		{"net", "via count", ""},
		{"region", "overflow", " tracks"},
		{"endpoint", "hold slack", " ps"},
		{"cell group", "displacement", " sites"},
	}
	kd := edakit.Pick(kinds)
	m := edakit.Scaled(ctx, 7, 11)
	delta := float64(edakit.Int(2, 6))
	before := make([]float64, m)
	after := make([]float64, m)
	target := make([]bool, m)
	for i := range before {
		before[i] = float64(edakit.Int(8, 60))
		target[i] = edakit.Int(0, 2) == 0
	}
	target[0], target[1] = true, false
	for i := range before {
		after[i] = before[i]
		if target[i] {
			after[i] -= delta
			if edakit.Int(0, 4) == 0 {
				after[i] += float64(edakit.Int(1, 2))
			}
		} else if edakit.Int(0, 3) == 0 {
			after[i] += float64(edakit.Pick([]int{-2, -1, 1, 2, 3}))
		}
	}
	var rows, tnames []string
	bad := 0
	for i := range before {
		name := edakit.Sprintf("%s%d", strings.ToUpper(kd.obj[:1]), i+1)
		rows = append(rows, edakit.Sprintf("%s: %s → %s", name, edakit.Num(before[i], 0), edakit.Num(after[i], 0)))
		if target[i] {
			tnames = append(tnames, name)
			if math.Abs(after[i]-before[i]+delta) > 0 {
				bad++
			}
		} else if after[i] != before[i] {
			bad++
		}
	}
	q := edakit.Sprintf("A plug-in should lower the %s of %ss %s by exactly %s and leave every other %s unchanged (tolerance 0). The differential report lists before → after values: %s. How many %ss fail the differential check?",
		kd.metric, kd.obj, strings.Join(tnames, ", "), edakit.Num(delta, 0), kd.obj, strings.Join(rows, ", "), kd.obj)
	exp := edakit.Sprintf("Count unintended movers plus targeted %ss that did not change by exactly %s: %d. The plug-in is promoted only when this is zero.", kd.obj, edakit.Num(delta, 0), bad)
	return choice(q, float64(bad), 0, "", exp, 1, float64(bad+1), float64(bad+2))
}

func qorBaseline(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		mu := float64(10000 * edakit.Int(50, 300))
		sd := float64(500 * edakit.Int(2, 20))
		z := float64(edakit.Int(-8, 12)) / 4
		x := mu + z*sd
		q := edakit.Sprintf("Routed wirelength on a reference block has a baseline mean of %s um and a standard deviation of %s um across seeds. A candidate change gives %s um. What is the z-score of the candidate against the baseline? Round to 2 decimal places.", edakit.Num(mu, 0), edakit.Num(sd, 0), edakit.Num(x, 0))
		exp := edakit.Sprintf("$z = (%s - %s)/%s = %s$. Values beyond about 2 or 3 are unlikely to be seed noise.", edakit.Num(x, 0), edakit.Num(mu, 0), edakit.Num(sd, 0), edakit.Num(z, 2))
		return edakit.Numeric(q, edakit.Num(z, 2), exp, nil)
	case 2:
		n := edakit.Scaled(ctx, 5, 9)
		th := edakit.Pick([]float64{0.5, 1, 2})
		var rows []string
		cnt := 0
		for i := 0; i < n; i++ {
			b := float64(10000 * edakit.Int(20, 150))
			d := float64(edakit.Int(-30, 40)) / 10 * th / 2
			for t := 0; math.Abs(d-th) < 0.06 && t < generator.MaxDrawAttempts; t++ {
				d = float64(edakit.Int(-30, 40)) / 10 * th / 2
			}
			if math.Abs(d-th) < 0.06 {
				d = th + 0.5
			}
			a := math.Round(b * (1 + d/100))
			pct := 100 * (a - b) / b
			if pct > th {
				cnt++
			}
			rows = append(rows, edakit.Sprintf("D%d %s → %s", i+1, edakit.Num(b, 0), edakit.Num(a, 0)))
		}
		q := edakit.Sprintf("A change may not worsen wirelength by more than %s%% on any design. Baseline → candidate wirelengths in um: %s. On how many designs does the change exceed the threshold?", edakit.Num(th, 1), strings.Join(rows, ", "))
		exp := edakit.Sprintf("Compute $100 (x_{new} - x_{base}) / x_{base}$ for each design and count deltas above %s%%: %d.", edakit.Num(th, 1), cnt)
		return edakit.Numeric(q, edakit.IntStr(cnt), exp, nil)
	}
	type met struct {
		name    string
		smaller bool
		th      float64
		unit    string
		abs     bool
	}
	ms := []met{{"WNS", false, 2, "ps", true}, {"TNS", false, 20, "ps", true}, {"area", true, 0.5, "square um", false}, {"power", true, 1, "mW", false}, {"DRC count", true, 0, "", true}}
	ms = edakit.Shuffle(ms)[:4]
	var rows []string
	cnt := 0
	for _, m := range ms {
		var b, a float64
		bad := false
		switch m.name {
		case "WNS":
			b = float64(-edakit.Int(0, 40))
			a = b + float64(edakit.Int(-5, 5))
			bad = b-a > m.th
		case "TNS":
			b = float64(-10 * edakit.Int(5, 60))
			a = b + float64(10*edakit.Int(-4, 4))
			bad = b-a > m.th
		case "area":
			b = float64(1000 * edakit.Int(20, 90))
			a = b * (1 + float64(edakit.Int(-10, 12))/1000)
			a = math.Round(a)
			bad = 100*(a-b)/b > m.th
		case "power":
			b = float64(edakit.Int(50, 300))
			a = b + float64(edakit.Int(-6, 6))
			bad = 100*(a-b)/b > m.th
		case "DRC count":
			b = float64(edakit.Int(0, 3))
			a = b + float64(edakit.Int(-1, 2))
			if a < 0 {
				a = 0
			}
			bad = a > b
		}
		if bad {
			cnt++
		}
		rule := ""
		if m.abs {
			rule = edakit.Sprintf("worse by more than %s %s", edakit.Num(m.th, 0), m.unit)
			if m.name == "DRC count" {
				rule = "any increase"
			}
		} else {
			rule = edakit.Sprintf("worse by more than %s%%", edakit.Num(m.th, 1))
		}
		rows = append(rows, edakit.Sprintf("%s %s → %s (regression if %s)", m.name, edakit.Num(b, 0), edakit.Num(a, 0), rule))
	}
	q := edakit.Sprintf("A differential report on one design lists baseline → candidate values: %s. Slack metrics (WNS, TNS) are better when larger; area, power and DRC count are better when smaller. How many metrics regress?", strings.Join(rows, "; "))
	exp := edakit.Sprintf("Apply each metric's direction and threshold: %d regress.", cnt)
	return edakit.Numeric(q, edakit.IntStr(cnt), exp, nil)
}

func seedNoise(ctx generator.GeneratorContext) generator.Problem {
	sd := float64(edakit.Int(2, 8)) / 10
	switch kp(ctx) {
	case 1:
		n := edakit.Pick([]int{4, 9, 16, 25, 36, 64})
		se := sd / math.Sqrt(float64(n))
		q := edakit.Sprintf("A router's wirelength varies across seeds with a standard deviation of %s%%. What is the standard error, in percent, of the mean of %d independent seeds? Round to 3 decimal places.", edakit.Num(sd, 1), n)
		exp := edakit.Sprintf("$\\sigma/\\sqrt{n} = %s/%s = %s\\%%$.", edakit.Num(sd, 1), edakit.Num(math.Sqrt(float64(n)), 0), edakit.Num(se, 3))
		return edakit.Numeric(q, edakit.Num(se, 3), exp, nil)
	case 2:
		d := float64(edakit.Int(1, 6)) / 10
		z := edakit.Pick([]float64{2, 2.5, 3})
		need := 2 * math.Pow(z*sd/d, 2)
		n := math.Ceil(need - 1e-9)
		q := edakit.Sprintf("A placer's QoR metric varies across seeds with a standard deviation of %s%%. To resolve a mean improvement of %s%% at $z = %s$ standard errors, using $n$ independent runs each of baseline and candidate, what is the smallest $n$?", edakit.Num(sd, 1), edakit.Num(d, 1), edakit.Num(z, 1))
		exp := edakit.Sprintf("$n \\ge 2 (z\\sigma/\\delta)^2 = 2 (%s \\times %s / %s)^2 = %s$, so %s runs of each.", edakit.Num(z, 1), edakit.Num(sd, 1), edakit.Num(d, 1), edakit.Num(need, 2), edakit.Num(n, 0))
		return edakit.Numeric(q, edakit.Num(n, 0), exp, nil)
	}
	d := float64(edakit.Int(1, 6)) / 10
	sdd := float64(edakit.Int(1, 5)) / 10
	z := edakit.Pick([]float64{2, 2.5, 3})
	need := math.Pow(z*sdd/d, 2)
	n := math.Ceil(need - 1e-9)
	q := edakit.Sprintf("Baseline and candidate run on the same seeds, and the per-seed differences in QoR have a standard deviation of %s%%. How many paired runs resolve a mean improvement of %s%% at $z = %s$?", edakit.Num(sdd, 1), edakit.Num(d, 1), edakit.Num(z, 1))
	exp := edakit.Sprintf("$n \\ge (z\\sigma_d/\\delta)^2 = (%s \\times %s / %s)^2 = %s$, so %s paired runs.", edakit.Num(z, 1), edakit.Num(sdd, 1), edakit.Num(d, 1), edakit.Num(need, 2), edakit.Num(n, 0))
	return edakit.Numeric(q, edakit.Num(n, 0), exp, nil)
}

var missingGoldenItems = []edakit.Item{
	{KP: 1, Q: "A derate rule fires when logic depth is below 5. The reference designs have depths 3, 4, 8 and 12. Which test is missing?",
		A: "A path of depth exactly 5, the boundary", D: []string{"A path of depth 2", "A path of depth 20", "None: both sides are covered"},
		Why: "Boundary values catch off-by-one errors in the condition."},
	{KP: 1, Q: "A router plug-in penalizes vias within 20 um of an analog block. The tests have nets at 5 um and 60 um. What is missing?",
		A: "Nets just inside, at and just outside 20 um", D: []string{"A net at 200 um", "More nets at 5 um", "A test of the placer"},
		Why: "Each condition needs tests on both sides of its threshold."},
	{KP: 1, Q: "A placement rule applies only to multi-height cells, but no reference design contains any. What has been tested?",
		A: "Only that the rule does nothing when it should not fire", D: []string{"That the rule works", "That multi-height cells are legal", "Everything the rule does"},
		Why: "A condition never satisfied in the tests is never exercised."},
	{KP: 1, Q: "A rule keys on three conditions: Vt class, clock versus data, and depth. What must the reference set include?",
		A: "Paths that satisfy and paths that fail each condition", D: []string{"Only paths that satisfy all three", "Only the worst path", "One path per corner"},
		Why: "Coverage means each condition is exercised both ways."},
	{KP: 2, Q: "A test checks that targeted nets lost vias. What collateral damage could it miss?",
		A: "Changes on untargeted nets", D: []string{"Nothing: targeted nets are all that matter", "Compile warnings", "Manifest typos"},
		Why: "Differentials must list unintended movers across all objects."},
	{KP: 2, Q: "A setup-timing fix is tested only on setup slack. Which check is missing?",
		A: "Hold slack, and early as well as late analysis at every corner", D: []string{"Wirelength of unrelated blocks only", "The SDC syntax", "The Liberty file date"},
		Why: "Fixes for one check commonly degrade the opposite check."},
	{KP: 2, Q: "A cost plug-in improves wirelength in tests but nobody measured runtime. What could ship unnoticed?",
		A: "A plug-in that doubles runtime in the fast loop", D: []string{"A faster plug-in", "A manifest change", "An ABI freeze"},
		Why: "Runtime and memory belong to the metrics a change might affect."},
	{KP: 2, Q: "A placement plug-in is tested on wirelength only. Which check must run on every output?",
		A: "Legality: no overlaps and cells on legal sites", D: []string{"A spell check of the manifest", "A license check", "Model accuracy"},
		Why: "Improving a metric by producing an illegal result is the classic failure."},
	{KP: 3, Q: "After a change, the team regenerates the goldens from the new engine's output and all tests pass. What is wrong?",
		A: "The goldens now record the change instead of testing it", D: []string{"Nothing: goldens should track the engine", "The tests are too strict", "The baseline is too old"},
		Why: "Regenerating goldens from the changed engine makes the test vacuous."},
	{KP: 3, Q: "The QoR baseline was recorded three engine versions ago. What does this do to a differential?",
		A: "Mixes unrelated changes into the comparison", D: []string{"Makes it more accurate", "Has no effect", "Speeds it up"},
		Why: "Baselines must be pinned to the engine version being compared."},
	{KP: 3, Q: "A golden accepts WNS within 50 ps of the expected value, while seed noise is 2 ps. What is the problem?",
		A: "The tolerance is so loose that almost any regression passes", D: []string{"The tolerance is too tight", "WNS cannot be tested", "The golden is too small"},
		Why: "Tolerances should sit near the measured noise."},
	{KP: 3, Q: "When may a golden's expected output be changed?",
		A: "Only through a reviewed change with a stated reason", D: []string{"Whenever it fails", "Whenever the model asks", "Every release automatically"},
		Why: "Controlled regeneration keeps goldens meaningful."},
}

func promotionGate(ctx generator.GeneratorContext) generator.Problem {
	type cand struct {
		build, golden bool
		unexp         int
		qor, rt       float64
		det, tier     bool
	}
	n := edakit.Scaled(ctx, 5, 8)
	qth := edakit.Pick([]float64{0.5, 1})
	rth := edakit.Pick([]float64{10, 15, 20})
	cs := make([]cand, n)
	var rows []string
	for i := range cs {
		c := cand{build: edakit.Int(0, 5) > 0, golden: edakit.Int(0, 4) > 0, unexp: 0,
			qor: float64(edakit.Int(-12, 15)) / 10, rt: float64(edakit.Int(-5, 24)), det: edakit.Int(0, 5) > 0, tier: edakit.Int(0, 5) > 0}
		if edakit.Int(0, 3) == 0 {
			c.unexp = edakit.Int(1, 4)
		}
		cs[i] = c
		yn := map[bool]string{true: "yes", false: "no"}
		rows = append(rows, edakit.Sprintf("C%d: builds %s, goldens %s, unintended movers %d, worst-design QoR change %s%%, runtime change %s%%, same output on rerun %s, diff within tier %s",
			i+1, yn[c.build], yn[c.golden], c.unexp, edakit.Num(c.qor, 1), edakit.Num(c.rt, 0), yn[c.det], yn[c.tier]))
	}
	hard := func(c cand) bool { return c.build && c.golden && c.unexp == 0 }
	thr := func(c cand) bool { return hard(c) && c.qor >= -qth && c.rt <= rth }
	full := func(c cand) bool { return thr(c) && c.det && c.tier }
	count := func(f func(cand) bool) int {
		k := 0
		for _, c := range cs {
			if f(c) {
				k++
			}
		}
		return k
	}
	head := edakit.Sprintf("Candidate changes (a negative QoR change is a regression): %s. ", strings.Join(rows, "; "))
	switch kp(ctx) {
	case 1:
		k := count(hard)
		q := head + "How many candidates pass the hard conditions: build, goldens and zero unintended movers?"
		return edakit.Numeric(q, edakit.IntStr(k), edakit.Sprintf("The hard conditions are a conjunction; %d candidates meet all three.", k), nil)
	case 2:
		k := count(thr)
		q := head + edakit.Sprintf("The gate requires the hard conditions (build, goldens, zero unintended movers), no worst-design regression beyond %s%% and a runtime increase of at most %s%%. How many candidates pass?", edakit.Num(qth, 1), edakit.Num(rth, 0))
		return edakit.Numeric(q, edakit.IntStr(k), edakit.Sprintf("Apply the hard conditions, then the QoR and runtime thresholds: %d pass.", k), nil)
	}
	k := count(full)
	q := head + edakit.Sprintf("The full gate requires build, goldens, zero unintended movers, no worst-design regression beyond %s%%, a runtime increase of at most %s%%, identical output on rerun with the same seed, and a diff within its declared tier. How many candidates reach human review?", edakit.Num(qth, 1), edakit.Num(rth, 0))
	return edakit.Numeric(q, edakit.IntStr(k), edakit.Sprintf("Every condition must hold; %d candidates pass them all.", k), nil)
}
