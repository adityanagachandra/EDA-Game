package stat

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// varianceAwareSurrogate: violation probability from a surrogate's mean and predictive sigma
// (KP 1), how many candidates a k-sigma rule sends to the accurate engine (KP 2), and the
// expected false passes among the accepted candidates (KP 3).
func varianceAwareSurrogate(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		mu := float64(edakit.Int(-20, 40))
		s := float64(edakit.Int(3, 20))
		p := cdf(-mu / s)
		ans := edakit.Num(100*p, 2)
		return edakit.Numeric(edakit.Sprintf("A fast timing surrogate predicts a path's slack as %s ps with a Gaussian predictive standard deviation of %s ps. What probability, in percent, does it give that the accurate timer would report negative slack? Round to 2 decimal places.", f(mu, 0), f(s, 0)),
			ans, edakit.Sprintf("$P(\\text{slack} < 0) = \\Phi(-\\hat{\\mu}/\\hat{\\sigma}) = \\Phi(%s) = %s$ percent.", f(-mu/s, 3), ans), map[string]string{"z": f(-mu/s, 3)})
	default:
		n := edakit.Scaled(ctx, 5, 8)
		k := edakit.Pick([]float64{2, 3})
		if ctx.KPOrdinal == 3 {
			k = edakit.Pick([]float64{1, 1.5})
		}
		var parts []string
		count := 0
		falsePass := 0.0
		accepted := 0
		for i := 0; i < n; i++ {
			mu := float64(edakit.Int(-30, 60))
			s := float64(edakit.Int(2, 15))
			parts = append(parts, edakit.Sprintf("(%s, %s)", f(mu, 0), f(s, 0)))
			if math.Abs(mu) < k*s {
				count++
			} else if mu > 0 {
				accepted++
				falsePass += cdf(-mu / s)
			}
		}
		stem := edakit.Sprintf("An agent's loop scores candidate paths with a surrogate that returns (predicted slack, predictive sigma) in ps: %s. Its rule: if $|\\hat{\\mu}| < %s\\hat{\\sigma}$, re-run the accurate timer; otherwise trust the surrogate's verdict.", strings.Join(parts, ", "), f(k, 1))
		if ctx.KPOrdinal == 2 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) || accepted == 0 {
			return edakit.Numeric(stem+" How many candidates go to the accurate timer?",
				edakit.IntStr(count), edakit.Sprintf("Compare $|\\hat{\\mu}|$ with $%s\\hat{\\sigma}$ for each: %d fall inside the uncertainty band.", f(k, 1), count),
				map[string]string{"recursed": edakit.IntStr(count)})
		}
		pl := sigPlaces(falsePass)
		if pl > 8 {
			pl = 8
		}
		ans := edakit.Num(falsePass, pl)
		return edakit.Numeric(stem+edakit.Sprintf(" Treating the surrogate's Gaussian as calibrated, what is the expected number of candidates that it accepts as passing but that truly have negative slack? %s", roundTo(pl)),
			ans, edakit.Sprintf("Sum $\\Phi(-\\hat{\\mu}/\\hat{\\sigma})$ over the %d candidates accepted as passing: %s. Each is below $Q(%s) = %s$; this holds only if the predictive sigma is calibrated.", accepted, ans, f(k, 1), sci(tail(k))),
			map[string]string{"accepted": edakit.IntStr(accepted)})
	}
}

// mcBudget: Neyman allocation (KP 1), cost-aware allocation under a time budget (KP 2), and the
// samples still needed after a pilot (KP 3).
func mcBudget(ctx generator.GeneratorContext) generator.Problem {
	names := []string{"split A", "split B", "split C"}
	w := []float64{float64(edakit.Int(2, 6)) / 10, float64(edakit.Int(2, 4)) / 10, 0}
	w[2] = 1 - w[0] - w[1]
	s := []float64{float64(edakit.Int(1, 10)), float64(edakit.Int(1, 10)), float64(edakit.Int(1, 10))}
	j := edakit.Int(0, 2)
	var desc []string
	for i := range w {
		desc = append(desc, edakit.Sprintf("%s: fraction %s, sigma %s nA", names[i], f(w[i], 1), f(s[i], 0)))
	}
	switch kp(ctx) {
	case 1:
		n := edakit.Pick([]int{500, 1000, 2000, 5000})
		den := 0.0
		for i := range w {
			den += w[i] * s[i]
		}
		nh := float64(n) * w[j] * s[j] / den
		ans := edakit.Num(nh, 0)
		return edakit.Numeric(edakit.Sprintf("An agent estimates mean leakage over a product that mixes three process splits (%s). It has %d simulations and stratifies by split. Under Neyman allocation, how many go to %s? Round to the nearest integer.", strings.Join(desc, "; "), n, names[j]),
			ans, edakit.Sprintf("$n_h = n W_h\\sigma_h/\\sum_j W_j\\sigma_j = %d \\times %s / %s = %s$.", n, f(w[j]*s[j], 2), f(den, 2), ans), map[string]string{"den": f(den, 2)})
	case 2:
		c := []float64{float64(edakit.Int(1, 9)), float64(edakit.Int(1, 9)), float64(edakit.Int(1, 9))}
		budget := float64(edakit.Pick([]int{3600, 7200, 10000}))
		den := 0.0
		for i := range w {
			den += w[i] * s[i] * math.Sqrt(c[i])
		}
		nh := budget * w[j] * s[j] / math.Sqrt(c[j]) / den
		ans := edakit.Num(nh, 0)
		return edakit.Numeric(edakit.Sprintf("An agent estimates mean leakage over three process splits (%s). One simulation costs %s, %s and %s CPU-seconds in the three splits, and the budget is %s CPU-seconds. With the variance-optimal allocation $n_h \\propto W_h\\sigma_h/\\sqrt{c_h}$ that spends the whole budget, how many simulations go to %s? Round to the nearest integer.", strings.Join(desc, "; "), f(c[0], 0), f(c[1], 0), f(c[2], 0), f(budget, 0), names[j]),
			ans, edakit.Sprintf("$n_h = C\\,(W_h\\sigma_h/\\sqrt{c_h}) / \\sum_j W_j\\sigma_j\\sqrt{c_j} = %s \\times %s / %s = %s$.", f(budget, 0), f(w[j]*s[j]/math.Sqrt(c[j]), 4), f(den, 4), ans), map[string]string{"den": f(den, 4)})
	default:
		n0 := edakit.Pick([]int{200, 500, 1000})
		fails := edakit.Int(2, n0/20)
		p := float64(fails) / float64(n0)
		se := float64(edakit.Int(1, 5)) / 1000
		need := int(math.Ceil(p * (1 - p) / (se * se)))
		more := need - n0
		if more < 0 {
			more = 0
		}
		return edakit.Numeric(edakit.Sprintf("A pilot of %d Monte Carlo samples shows %d failures. The agent must report the failure rate with a standard error of at most %s. Using the pilot estimate, how many more samples should it request? Give an integer (0 if none).", n0, fails, f(se, 3)),
			edakit.IntStr(more), edakit.Sprintf("$\\hat{p} = %s$. $N = \\hat{p}(1-\\hat{p})/SE^2 = %d$ in total, so %d more beyond the pilot.", f(p, 4), need, more), map[string]string{"total": edakit.IntStr(need)})
	}
}

// goldenStatTests: the acceptance band of a golden Monte Carlo test (KP 1), the family-wise
// false-alarm rate of a suite (KP 2), and the Bonferroni z threshold (KP 3).
func goldenStatTests(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		p := float64(edakit.Int(5, 50)) / 1000
		n := edakit.Pick([]int{1000, 2000, 5000, 10000, 20000})
		z := edakit.Pick([]float64{3, 4})
		hw := z * math.Sqrt(p*(1-p)/float64(n))
		ans := edakit.Num(100*hw, 3)
		return edakit.Numeric(edakit.Sprintf("A golden test feeds a Monte Carlo yield skill a synthetic limit state whose exact failure probability is %s, runs %d samples with a fresh seed, and passes if the estimate lands within %s standard errors of the truth. What is the half-width of the acceptance band, in percentage points? Round to 3 decimal places.", f(p, 3), n, f(z, 0)),
			ans, edakit.Sprintf("$SE = \\sqrt{p(1-p)/N} = %s$, so the half-width is $%s \\times SE = %s$ percentage points.", f(math.Sqrt(p*(1-p)/float64(n)), 6), f(z, 0), ans), map[string]string{"se": f(math.Sqrt(p*(1-p)/float64(n)), 6)})
	case 2:
		m := edakit.Scaled(ctx, 10, 400)
		a := edakit.Pick([]float64{0.05, 0.01, 0.0027, 0.001})
		fw := 1 - math.Pow(1-a, float64(m))
		ans := edakit.Num(100*fw, 2)
		return edakit.Numeric(edakit.Sprintf("A regression suite holds %d independent statistical golden tests, each with false-alarm probability %s when the engine is correct. What is the probability, in percent, that a correct engine change fails at least one test? Round to 2 decimal places.", m, f(a, 4)),
			ans, edakit.Sprintf("$1 - (1-\\alpha)^m = 1 - (1-%s)^{%d} = %s$ percent. A flaky suite trains the harness to ignore failures.", f(a, 4), m, ans), map[string]string{"fwer": ans})
	default:
		m := edakit.Scaled(ctx, 10, 500)
		fw := edakit.Pick([]float64{0.01, 0.05, 0.001})
		a := fw / float64(m)
		z := tailInv(a / 2)
		ans := edakit.Num(z, 2)
		return edakit.Numeric(edakit.Sprintf("A suite of %d two-sided statistical golden tests must have a family-wise false-alarm rate of at most %s. Using a Bonferroni split and a Gaussian test statistic, what z threshold should each test use? Round to 2 decimal places.", m, f(fw, 3)),
			ans, edakit.Sprintf("Per-test $\\alpha = %s/%d = %s$, split over two tails: $z = Q^{-1}(%s) = %s$.", f(fw, 3), m, sci(a), sci(a/2), ans), map[string]string{"alpha": sci(a)})
	}
}

var rareEventSkillItems = []edakit.Item{
	{KP: 1, Q: "A rare-event estimator is packaged as a skill. Which manifest entry most helps an agent decide whether it fits a new bitcell?",
		A: "The failure-region structures it handles (one region or several) and the dimension range it was validated on", D: []string{"The number of source files in the skill", "The name of the engineer who wrote it", "The compiler version used to build it"},
		Why: "Selection matches the design's characteristics against what the manifest advertises. Region structure and dimension decide whether a mean-shift or a multi-region method is safe."},
	{KP: 1, Q: "What must the skill return besides the point estimate of the failure rate, for a downstream agent to act on it?",
		A: "A confidence interval, the effective sample size and the number of simulator calls", D: []string{"Only the estimate, rounded to one digit", "The raw waveform of every simulation", "The random seed only"},
		Why: "An agent decides whether to spend more budget from the interval width and the health of the weights. A bare point estimate hides both."},
	{KP: 1, Q: "Two skills are listed: mean-shift importance sampling and subset simulation. The bitcell has read and write failures in opposite directions. Which does a manifest-reading agent pick?",
		A: "The skill whose manifest states it handles multiple failure regions", D: []string{"Mean-shift importance sampling, because it is faster per sample", "Whichever was used last", "Plain Monte Carlo, to be safe"},
		Why: "A single shift would miss one region and underestimate. The manifest's region-structure field is what makes that choice a cheap text decision."},
	{KP: 1, Q: "Why does a rare-event skill accept a batch of parameter samples per call to the simulator interface instead of one sample per call?",
		A: "Coarse calls keep the overhead low and let the simulator parallelize across the batch", D: []string{"Batching changes the estimator's bias", "Simulators cannot run single samples", "Batching hides the seed from the harness"},
		Why: "Plug-in calls should be coarse, per batch rather than per item, so the hot path stays efficient."},
	{KP: 2, Q: "A user asks to target 5.5 sigma instead of 6 sigma for a new array. Which change tier is that?",
		A: "Tier 1: a configuration value", D: []string{"Tier 2: a new plug-in", "Tier 3: a core edit", "No change is possible"},
		Why: "The sigma target is a parameter of the run. It needs no new code and no recompile."},
	{KP: 2, Q: "A team wants a proposal that mixes two mean shifts, one per failure region. Where does that land?",
		A: "Tier 2: a new proposal-distribution plug-in at the sampler seam", D: []string{"Tier 1: a configuration value", "Tier 3: an edit to the weight computation in the core", "Outside the skill, in the circuit simulator"},
		Why: "The estimator core computes weights for any proposal it is given. A new proposal is an implementation behind that seam."},
	{KP: 2, Q: "Someone proposes changing how the core computes likelihood-ratio weights to speed things up. Why is that Tier 3 with human review?",
		A: "The weights are what make the estimator unbiased; a mistake there silently corrupts every result", D: []string{"Weight code is too long to edit", "Weights only affect logging", "Tier 3 is the default for any speed change"},
		Why: "Correctness invariants live in the core. A change there needs full regression and a reviewed diff because every estimate depends on it."},
	{KP: 2, Q: "What context must the sampler seam expose so that a new proposal plug-in can work?",
		A: "The parameter dimension and distributions, any known failure points, and the sample budget", D: []string{"The full SPICE netlist text only", "The user's login and password", "Only the final failure estimate"},
		Why: "A proposal needs to know the space it samples, where failures were found, and how many samples it may use. It does not need privileged data."},
	{KP: 3, Q: "What is the most useful golden test for a rare-event estimator skill?",
		A: "A synthetic limit state with a known exact failure probability, checked for interval coverage across many seeds", D: []string{"A real bitcell whose failure rate nobody knows", "A test that the code compiles", "A single run that matches yesterday's number to every digit"},
		Why: "Only a case with known truth can show the estimator is unbiased and that its intervals cover the truth at the stated rate."},
	{KP: 3, Q: "A new proposal plug-in passes the golden tests. What should the differential test against the baseline skill on a reference cell show?",
		A: "Estimates agree within their combined confidence intervals, with fewer simulator calls", D: []string{"Bitwise identical estimates", "A lower failure rate than the baseline", "A larger confidence interval"},
		Why: "Two correct estimators agree statistically, not bitwise. The promotion case is the same answer at lower cost."},
	{KP: 3, Q: "Why must a rare-event skill be deterministic for a given seed?",
		A: "So a reported result can be reproduced and a diff between versions is attributable to the code change", D: []string{"So the estimate has zero variance", "So it never needs golden tests", "Because random seeds are slow"},
		Why: "Seedable determinism lets the harness rerun a failing case exactly and makes regressions bisectable."},
	{KP: 3, Q: "A golden test that compares an estimate to the truth with a 2-sigma band fails about once in twenty runs of a correct skill. What is the right fix?",
		A: "Widen the band or fix the seed set so the false-alarm rate is negligible, and keep a separate coverage test", D: []string{"Delete the test", "Rerun until it passes", "Narrow the band to 1 sigma"},
		Why: "A statistical golden test must have a known, small false-alarm rate. Retrying until green hides real regressions."},
}

var yieldAgentItems = []edakit.Item{
	{KP: 1, Q: "An agentic yield loop for an analog block reads which inputs before proposing a fix?",
		A: "Per-spec failure counts, the sensitivities of each failing spec, and the confidence of each estimate", D: []string{"Only the overall yield number", "The layout's GDS file size", "The number of agents running"},
		Why: "To propose a targeted fix the agent needs to know which specs fail, which parameters drive them, and how sure the estimates are."},
	{KP: 1, Q: "The agent proposes resizing a mirror to recentre a gain spec. What must the loop verify before a human sees the change?",
		A: "That the yield gain holds in accurate simulation with fresh samples and that no other spec degraded", D: []string{"That the surrogate predicted a gain", "That the netlist still parses", "That the change is smaller than 1 percent"},
		Why: "A fix is real only if accurate simulation confirms it and the differential check shows nothing else moved for the worse."},
	{KP: 1, Q: "Which step of an agentic yield-closure loop belongs to a human?",
		A: "Approving the final design change after reviewing the verified diff", D: []string{"Running each Monte Carlo sample", "Computing sensitivities", "Rerunning the surrogate"},
		Why: "Agents propose and verify; the final ECO is approved by a person, as in signoff ECO loops."},
	{KP: 1, Q: "The agent watches yield estimates for several variants. What tells it a variant needs more samples before it can be ranked?",
		A: "Its confidence interval overlaps the intervals of competing variants or the yield target", D: []string{"Its point estimate is the highest", "It was simulated most recently", "Its netlist is largest"},
		Why: "Budget should go where the decision is still uncertain: where intervals overlap the comparison that matters."},
	{KP: 2, Q: "An RL reward is the yield from 200 Monte Carlo samples drawn with a seed the agent chooses. What reward hack appears?",
		A: "The agent learns seeds whose samples happen to pass, inflating reward without improving the design", D: []string{"The agent stops changing the design", "The reward becomes negative", "The simulator crashes"},
		Why: "If the policy controls the randomness it can game it. The harness must own the seeds and draw fresh ones for scoring."},
	{KP: 2, Q: "A reward uses a fast surrogate's yield prediction. The agent drives designs into a region with no training data, where the surrogate predicts 100 percent yield. What is this?",
		A: "Reward hacking of a fast proxy by exploiting its extrapolation error", D: []string{"A genuine yield breakthrough", "Overfitting of the agent's tokenizer", "A bug in the voting logic"},
		Why: "The gap between a fast proxy and signoff truth is the central risk of tool-grounded RL. Out-of-distribution designs must be sent to the accurate engine."},
	{KP: 2, Q: "Which reward resists gaming by too few samples better than the point-estimate yield?",
		A: "The lower confidence bound on yield", D: []string{"The upper confidence bound on yield", "The number of samples run", "The yield rounded to one digit"},
		Why: "With few samples the lower bound is low, so the policy only earns reward by running enough samples to be sure."},
	{KP: 2, Q: "The agent raises yield by editing the testbench so a measurement window is narrower. Why must the harness block that?",
		A: "Specs and testbench definitions are ground truth, not part of the design the agent may change", D: []string{"Testbenches are too slow to edit", "Narrow windows always lower yield", "Testbench edits are Tier 1 and need no review"},
		Why: "Changing what is measured changes the question rather than improving the answer. Constraints belong to the protected core of the loop."},
	{KP: 3, Q: "An agent tries 60 variants and reports the best observed Monte Carlo yield. Why is that number biased upward?",
		A: "Picking the maximum of noisy estimates selects lucky noise (the winner's curse)", D: []string{"Monte Carlo estimates are always biased upward", "The agent rounds up", "More variants lower the true yield"},
		Why: "The largest of many noisy estimates tends to exceed its own true value. Re-estimate the winner with fresh samples before reporting."},
	{KP: 3, Q: "An agent reports zero failures in 1,000 samples and claims a 1-ppm failure target is met. What is wrong?",
		A: "Zero failures in 1,000 samples only bounds the rate near 3 in 1,000 at 95 percent confidence", D: []string{"Nothing: zero failures proves the target", "It should have reported 1 failure", "1,000 samples is more than needed"},
		Why: "By the rule of three, a 1-ppm claim at 95 percent confidence needs about 3 million failure-free samples, or a rare-event method."},
	{KP: 3, Q: "A centring change raises a block's yield from 92 to 97 percent. What must the differential report show to promote it?",
		A: "The yield gain with its interval, and every other spec's yield unchanged or better within noise", D: []string{"Only the new yield number", "That the change touched fewer than ten devices", "That the surrogate agreed"},
		Why: "Differential evidence shows the intended metric moved and nothing else regressed."},
	{KP: 3, Q: "Why should the reward of a yield-optimizing agent be computed by the harness on held-out seeds, not by the agent's own tool calls?",
		A: "So the agent cannot see or influence the samples that score it", D: []string{"Held-out seeds simulate faster", "Agents cannot call simulators", "It removes the need for confidence intervals"},
		Why: "Scoring on randomness the policy never touched keeps the reward an honest measure of the design."},
}
