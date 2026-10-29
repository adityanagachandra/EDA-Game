package agentic

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var rlItems = []edakit.Item{
	{KP: 1, Q: "In tool-grounded reinforcement learning, who computes the reward for a model's output?",
		A: "Deterministic engines run on the output, such as a compiler, legality checker or timer", D: []string{"The model grading itself", "A human rating each output", "The manifest"},
		Why: "The reward comes from tools, which makes it objective."},
	{KP: 1, Q: "What is the action in tool-grounded RL for a skill-editing model?",
		A: "The code edit the model produces", D: []string{"The compiler's output", "The golden test", "The reward value"},
		Why: "The model's output is the action; engines score it."},
	{KP: 1, Q: "How does a policy-gradient method use tool-computed rewards?",
		A: "It raises the probability of outputs that earned above-average reward", D: []string{"It copies the highest-reward output into the engine", "It lowers every output's probability equally", "It ignores rewards below zero"},
		Why: "Policy gradients shift probability towards better-than-baseline actions."},
	{KP: 1, Q: "Which is a well-known example of RL with a reward computed from placement proxies?",
		A: "Macro placement with a reward built from proxy wirelength, density and congestion", D: []string{"Standard-cell legalization with Abacus", "Static timing analysis with OpenSTA", "LVS graph matching"},
		Why: "Macro placement by RL used proxy-based rewards."},
	{KP: 2, Q: "Which reward signal is sparse?",
		A: "Pass or fail of compilation", D: []string{"Relative wirelength improvement", "Change in total negative slack", "Overflow reduction"},
		Why: "Binary verdicts give no gradation between failures."},
	{KP: 2, Q: "Which reward component is hardest to argue with?",
		A: "A legality or equivalence verdict established by an exact checker", D: []string{"A fast congestion estimate", "A learned QoR predictor", "The model's self-assessment"},
		Why: "Exact checks are the strongest signals in this domain."},
	{KP: 2, Q: "How are binary verdicts and continuous metrics usually combined in a reward?",
		A: "Gate on the verdicts, then score passing outputs on the metrics", D: []string{"Add them with equal weight", "Use only the metrics", "Use only the verdicts"},
		Why: "Gating prevents illegal outputs from earning credit."},
	{KP: 2, Q: "Why are dense metric rewards easier to game than verdicts?",
		A: "Metrics are partial views of quality, so a policy can raise them while hurting what they ignore", D: []string{"Metrics are always wrong", "Verdicts change every run", "Dense rewards are computed by the model"},
		Why: "Every metric is a proxy."},
	{KP: 3, Q: "A rollout needs full place-and-route and signoff, taking hours. How is training made practical?",
		A: "Fast engines compute most rewards, with periodic accurate evaluation", D: []string{"Skip the reward", "Use the model's guess as reward", "Train only on compile success"},
		Why: "This creates the proxy-versus-signoff gap that must be watched."},
	{KP: 3, Q: "Rewards involve a customer's designs. Where must they be computed?",
		A: "On-site, inside the sandbox and the IP boundary", D: []string{"At the vendor after copying the designs", "On a public cloud", "On the model provider's servers"},
		Why: "Training on customer data happens where the data is."},
	{KP: 3, Q: "The reward engine is updated mid-training and the same action now earns a different reward. What goes wrong?",
		A: "Training chases a moving target; reward engines must be versioned and deterministic", D: []string{"Nothing: newer is better", "The policy learns faster", "The reward becomes sparse"},
		Why: "Inconsistent rewards are noise to the learner."},
	{KP: 3, Q: "Why must reward engines be deterministic?",
		A: "Otherwise identical actions earn different rewards and training learns noise", D: []string{"To make them faster", "Because RL forbids randomness", "So the model can predict them"},
		Why: "Reward noise directly degrades learning."},
}

func rewardDesign(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 9)
		var rows []string
		ok := 0
		yn := map[bool]string{true: "yes", false: "no"}
		for i := 0; i < n; i++ {
			leg := edakit.Int(0, 3) > 0
			drc := edakit.Int(0, 3) > 0
			hp := float64(edakit.Int(-4, 8))
			if leg && drc {
				ok++
			}
			rows = append(rows, edakit.Sprintf("P%d: overlap-free %s, DRC-clean after routing %s, HPWL change %s%%", i+1, yn[leg], yn[drc], edakit.Num(hp, 0)))
		}
		q := edakit.Sprintf("A placement reward gates on legality (overlap-free) and routability (DRC-clean after routing) before scoring HPWL. Candidates: %s. How many candidates receive a nonzero score?", strings.Join(rows, "; "))
		return edakit.Numeric(q, edakit.IntStr(ok), edakit.Sprintf("Only candidates passing both gates are scored, whatever their HPWL: %d.", ok), nil)
	case 2:
		if edakit.Int(0, 1) == 0 {
			b := float64(edakit.Int(100, 400)) / 100
			x := b * (1 - float64(edakit.Int(-20, 60))/1000)
			x = math.Round(x*1000) / 1000
			r := (b - x) / b
			q := edakit.Sprintf("A reward normalizes wirelength per design as $r = (x_{base} - x)/x_{base}$. The baseline wirelength is %s m and the candidate's is %s m. What is $r$? Round to 4 decimal places.", edakit.Num(b, 2), edakit.Num(x, 3))
			return edakit.Numeric(q, edakit.Num(r, 4), edakit.Sprintf("$(%s - %s)/%s = %s$.", edakit.Num(b, 2), edakit.Num(x, 3), edakit.Num(b, 2), edakit.Num(r, 4)), nil)
		}
		t := float64(100 * edakit.Int(5, 20))
		wb := float64(-edakit.Int(10, 120))
		w := wb + float64(edakit.Int(-20, 60))
		r := (w - wb) / t
		q := edakit.Sprintf("A reward normalizes slack changes by the clock period: $r = (\\mathrm{WNS} - \\mathrm{WNS}_{base}) / T_{clk}$. The baseline WNS is %s ps, the candidate's is %s ps and the period is %s ps. What is $r$? Round to 4 decimal places.", edakit.Num(wb, 0), edakit.Num(w, 0), edakit.Num(t, 0))
		return edakit.Numeric(q, edakit.Num(r, 4), edakit.Sprintf("$(%s - (%s))/%s = %s$.", edakit.Num(w, 0), edakit.Num(wb, 0), edakit.Num(t, 0), edakit.Num(r, 4)), nil)
	}
	ws := [][3]float64{{0.5, 0.3, 0.2}, {0.4, 0.4, 0.2}, {0.3, 0.3, 0.4}, {0.6, 0.2, 0.2}}
	w := edakit.Pick(ws)
	hb := float64(edakit.Int(150, 300)) / 100
	h := math.Round(hb*(1-float64(edakit.Int(-10, 60))/1000)*1000) / 1000
	ob := float64(100 * edakit.Int(5, 20))
	o := ob - float64(50*edakit.Int(-2, 8))
	t := float64(100 * edakit.Int(5, 15))
	wb := float64(-edakit.Int(20, 120))
	wn := wb + float64(edakit.Int(-10, 40))
	r1, r2, r3 := (hb-h)/hb, (ob-o)/ob, (wn-wb)/t
	R := w[0]*r1 + w[1]*r2 + w[2]*r3
	q := edakit.Sprintf("A gated candidate is scored as $R = %s\\,r_{WL} + %s\\,r_{OF} + %s\\,r_{WNS}$, with $r_{WL}$ and $r_{OF}$ relative improvements over the baseline and $r_{WNS}$ the WNS change divided by the clock period. Baseline: wirelength %s m, overflow %s, WNS %s ps. Candidate: wirelength %s m, overflow %s, WNS %s ps. The period is %s ps. What is $R$? Round to 4 decimal places.",
		edakit.Num(w[0], 1), edakit.Num(w[1], 1), edakit.Num(w[2], 1), edakit.Num(hb, 2), edakit.Num(ob, 0), edakit.Num(wb, 0), edakit.Num(h, 3), edakit.Num(o, 0), edakit.Num(wn, 0), edakit.Num(t, 0))
	exp := edakit.Sprintf("$r_{WL} = %s$, $r_{OF} = %s$, $r_{WNS} = %s$, so $R = %s$.", edakit.Num(r1, 4), edakit.Num(r2, 4), edakit.Num(r3, 4), edakit.Num(R, 4))
	return edakit.Numeric(q, edakit.Num(R, 4), exp, nil)
}

var hackItems = []edakit.Item{
	{KP: 1, Q: "A placement policy is rewarded on HPWL alone. What hack is most likely?",
		A: "Packing cells so tightly that the router cannot finish", D: []string{"Spreading cells evenly", "Improving timing at the cost of HPWL", "Refusing to place macros"},
		Why: "HPWL ignores congestion, so the policy trades routability for shorter wires."},
	{KP: 1, Q: "A policy is rewarded by a probabilistic congestion estimate that ignores pin access. What can it learn?",
		A: "Layouts that satisfy the estimate but create unroutable pin clusters", D: []string{"Layouts with lower HPWL only", "Layouts with fewer cells", "Nothing harmful"},
		Why: "Whatever the estimate ignores becomes the place to cheat."},
	{KP: 1, Q: "An ECO policy is rewarded by a fast timer that ignores crosstalk. What might it do?",
		A: "Place aggressor and victim nets side by side to shorten them", D: []string{"Add shielding everywhere", "Remove all buffers", "Lower the clock frequency"},
		Why: "Crosstalk-blind rewards favor coupling-heavy layouts."},
	{KP: 1, Q: "A code-editing reward counts golden tests passed. What hack appears?",
		A: "Special-casing the golden inputs", D: []string{"Writing more general code", "Adding more tests", "Improving performance"},
		Why: "Visible tests become targets; hidden tests catch it."},
	{KP: 2, Q: "A reward counts violations fixed by an ECO agent but not violations created. What will the agent learn?",
		A: "Fixes that create new violations elsewhere", D: []string{"Minimal, careful fixes", "To fix nothing", "To lower the threshold"},
		Why: "Uncounted side effects are free to the policy."},
	{KP: 2, Q: "A reward is improvement over a baseline that the policy's own run produces. What can go wrong?",
		A: "The policy can make the baseline worse to inflate its improvement", D: []string{"Nothing", "The reward becomes sparse", "The baseline becomes too good"},
		Why: "Baselines must be frozen outside the policy's influence."},
	{KP: 2, Q: "A reward measures global placement runtime alone. What hack is predictable?",
		A: "Terminating global placement early", D: []string{"Running more iterations", "Improving legality", "Adding more nets"},
		Why: "Reward runtime only, and the policy skips work."},
	{KP: 2, Q: "A plug-in is rewarded by a metric computed in a component the plug-in can edit. What hack becomes possible?",
		A: "Changing how the metric is reported instead of improving the result", D: []string{"Improving the result", "Making the build slower", "Writing better manifests"},
		Why: "Reward code must sit outside what the policy can change."},
	{KP: 3, Q: "Proxy reward keeps rising while periodic signoff evaluation stalls and then falls. What has happened?",
		A: "The policy has started optimizing the proxy rather than the task", D: []string{"Training is converging well", "Signoff is broken", "The proxy is more accurate than signoff"},
		Why: "Divergence between proxy and accurate scores is the signature of reward hacking."},
	{KP: 3, Q: "What is the first defense against a placement policy producing illegal results with great metrics?",
		A: "Gate on legality with an independent checker the policy cannot modify", D: []string{"Raise the metric weights", "Train longer", "Lower the learning rate"},
		Why: "Legality gates come before any metric."},
	{KP: 3, Q: "How is a proxy reward paired against its blind spot?",
		A: "Check samples with an engine that measures what the proxy ignores, such as routing a sample of placements", D: []string{"Use two copies of the same proxy", "Average the proxy over seeds", "Remove the proxy"},
		Why: "Each proxy needs a complementary check."},
	{KP: 3, Q: "Where should baselines and reward code live relative to the policy's sandbox?",
		A: "Outside it, frozen and versioned", D: []string{"Inside it, so the policy can update them", "In the policy's weights", "In the manifest"},
		Why: "What the policy can touch, it can game."},
}

func proxyGap(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 6)
		var xs, ys []float64
		for i := 0; i < n; i++ {
			x := float64(edakit.Int(0, 10))
			xs = append(xs, x)
			ys = append(ys, math.Max(0, x+float64(edakit.Int(-4, 4))))
		}
		xs[0], xs[1] = 0, 10
		mx, my := 0.0, 0.0
		for i := range xs {
			mx += xs[i]
			my += ys[i]
		}
		mx /= float64(n)
		my /= float64(n)
		sxy, sxx, syy := 0.0, 0.0, 0.0
		for i := range xs {
			sxy += (xs[i] - mx) * (ys[i] - my)
			sxx += (xs[i] - mx) * (xs[i] - mx)
			syy += (ys[i] - my) * (ys[i] - my)
		}
		if syy == 0 {
			ys[1]++
			my += 1 / float64(n)
			sxy, syy = 0, 0
			for i := range xs {
				sxy += (xs[i] - mx) * (ys[i] - my)
				syy += (ys[i] - my) * (ys[i] - my)
			}
		}
		r := sxy / math.Sqrt(sxx*syy)
		var rows []string
		for i := range xs {
			rows = append(rows, edakit.Sprintf("(%s, %s)", edakit.Num(xs[i], 0), edakit.Num(ys[i], 0)))
		}
		q := edakit.Sprintf("Candidates evaluated by a fast proxy and by signoff give (proxy improvement, signoff improvement) pairs: %s. What is the Pearson correlation between proxy and signoff improvement? Round to 3 decimal places.", strings.Join(rows, ", "))
		exp := edakit.Sprintf("$S_{xy} = %s$, $S_{xx} = %s$, $S_{yy} = %s$, so $r = %s$.", edakit.Num(sxy, 3), edakit.Num(sxx, 3), edakit.Num(syy, 3), edakit.Num(r, 3))
		return edakit.Numeric(q, edakit.Num(r, 3), exp, nil)
	case 2:
		n := edakit.Scaled(ctx, 4, 7)
		type c struct{ p, s float64 }
		var cs []c
		usedP := map[int]bool{}
		for t := 0; len(cs) < n; t++ {
			p := edakit.Int(5, 50)
			if t > generator.MaxDrawAttempts {
				p = 51 + len(cs)
			}
			if usedP[p] {
				continue
			}
			usedP[p] = true
			cs = append(cs, c{float64(p) / 10, float64(edakit.Int(0, 40)) / 10})
		}
		bp, bs := 0, 0
		for i := range cs {
			if cs[i].p > cs[bp].p {
				bp = i
			}
			if cs[i].s > cs[bs].s {
				bs = i
			}
		}
		reg := cs[bs].s - cs[bp].s
		var rows []string
		for i, x := range cs {
			rows = append(rows, edakit.Sprintf("K%d (%s, %s)", i+1, edakit.Num(x.p, 1), edakit.Num(x.s, 1)))
		}
		q := edakit.Sprintf("Candidates with (proxy improvement, signoff improvement) in percent: %s. An agent picks the candidate with the best proxy improvement. What is its regret, the best signoff improvement minus that of the chosen candidate, in percentage points? Round to 1 decimal place.", strings.Join(rows, ", "))
		exp := edakit.Sprintf("The proxy picks K%d (signoff %s); the best at signoff is K%d (%s). Regret %s.", bp+1, edakit.Num(cs[bp].s, 1), bs+1, edakit.Num(cs[bs].s, 1), edakit.Num(reg, 1))
		return edakit.Numeric(q, edakit.Num(reg, 1), exp, nil)
	}
	n := edakit.Scaled(ctx, 7, 10)
	k := edakit.Int(2, 3)
	m := edakit.Int(k, k+2)
	sig := edakit.Shuffle(func() []int {
		v := make([]int, n)
		for i := range v {
			v[i] = i + 1
		}
		return v
	}())
	var rows []string
	for i := 0; i < n; i++ {
		rows = append(rows, edakit.Sprintf("proxy rank %d → signoff rank %d", i+1, sig[i]))
	}
	hit := 0
	for i := 0; i < m; i++ {
		if sig[i] <= k {
			hit++
		}
	}
	var picked []int
	for i := 0; i < m; i++ {
		picked = append(picked, sig[i])
	}
	sort.Ints(picked)
	q := edakit.Sprintf("Eight-to-ten candidates are ranked by a fast proxy and by signoff: %s. The agent sends the proxy's top %d to signoff. How many of the true top %d (by signoff) are in that shortlist?", strings.Join(rows, ", "), m, k)
	q = strings.Replace(q, "Eight-to-ten candidates", edakit.Sprintf("%d candidates", n), 1)
	exp := edakit.Sprintf("The shortlist has signoff ranks %v; %d of them are within the top %d, a recall of %d/%d.", picked, hit, k, hit, k)
	return edakit.Numeric(q, edakit.IntStr(hit), exp, nil)
}
