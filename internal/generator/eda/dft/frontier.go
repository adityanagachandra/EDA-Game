package dft

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// goldenCoverage: failing golden checks (KP 1), the pattern-count ceiling (KP 2), coverage delta
// from raw counts (KP 3).
func goldenCoverage(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 7)
		dc := float64(edakit.Pick([]int{2, 5, 10})) / 100
		dp := edakit.Pick([]int{1, 2, 3, 5})
		fails := 0
		var rows []string
		var why []string
		for i := 0; i < n; i++ {
			bc := float64(edakit.Int(9500, 9990)) / 100
			bp := edakit.Int(800, 20000)
			nc := bc + float64(edakit.Int(-15, 10))/100
			np := int(float64(bp) * (1 + float64(edakit.Int(-40, 80))/1000))
			bad := bc-nc > dc+1e-9 || float64(np) > float64(bp)*(1+float64(dp)/100)+1e-9
			if bad {
				fails++
				why = append(why, edakit.Sprintf("D%d", i+1))
			}
			rows = append(rows, edakit.Sprintf("D%d: baseline %s percent with %d patterns, new %s percent with %d patterns", i+1, edakit.Num(bc, 2), bp, edakit.Num(nc, 2), np))
		}
		w := "none"
		if fails > 0 {
			w = strings.Join(why, ", ")
		}
		return edakit.Numeric(edakit.Sprintf("An agent changed the ATPG skill's fault-ordering policy. The golden suite fails a design if test coverage drops by more than %s percentage points or the pattern count grows by more than %d percent. Results: %s. How many designs fail?", edakit.Num(dc, 2), dp, strings.Join(rows, "; ")),
			edakit.IntStr(fails), "Failing: "+w+". Both limits are checked per design; a gain on one design does not excuse a regression on another.", nil)
	case 2:
		bp := edakit.Scaled(ctx, 500, 50000)
		dp := float64(edakit.Pick([]int{5, 10, 15, 20, 25, 30})) / 10
		mx := int(math.Floor(float64(bp) * (1 + dp/100)))
		return edakit.Numeric(edakit.Sprintf("A golden test for an ATPG skill allows the pattern count to grow by at most %s percent over the baseline of %d patterns. What is the largest pattern count that still passes?", edakit.Num(dp, 1), bp),
			edakit.IntStr(mx), edakit.Sprintf("$\\lfloor %d \\times (1 + %s/100) \\rfloor = %d$. Pattern count is tester time and memory, so it is a first-class QoR metric next to coverage.", bp, edakit.Num(dp, 1), mx), nil)
	default:
		tot := edakit.Scaled(ctx, 20000, 400000)
		u := tot * edakit.Int(10, 40) / 1000
		d0 := tot - u - tot*edakit.Int(5, 30)/1000
		d1 := d0 + edakit.Int(-400, 400)
		u1 := u + edakit.Int(-200, 200)
		tc0 := 100 * float64(d0) / float64(tot-u)
		tc1 := 100 * float64(d1) / float64(tot-u1)
		ans := edakit.Num(tc1-tc0, 3)
		return edakit.Numeric(edakit.Sprintf("Baseline ATPG on %d faults detects %d with %d proven untestable. After a skill change it detects %d with %d proven untestable. By how many percentage points did test coverage (detected over total minus untestable) change? Give a negative number for a drop. Round to 3 decimal places.", tot, d0, u, d1, u1),
			ans, edakit.Sprintf("Baseline TC = %s percent, new TC = %s percent, change %s points. Note that proving more faults untestable raises test coverage without detecting anything new, a reward a careless metric could be gamed by.", edakit.Num(tc0, 4), edakit.Num(tc1, 4), ans), nil)
	}
}

// diagnosisAgent: posterior of the top candidate (KP 1), candidates to send to failure analysis
// (KP 2), posterior with mismatch-based likelihoods (KP 3).
func diagnosisAgent(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 5)
	prior := make([]float64, n)
	lik := make([]float64, n)
	w := make([]int, n)
	for i := range prior {
		w[i] = edakit.Int(1, 9)
		prior[i] = float64(w[i])
	}
	k := kp(ctx)
	var lt []string
	if k == 3 {
		ms := make([]int, n)
		for attempt := 0; ; attempt++ {
			for i := range ms {
				ms[i] = edakit.Int(0, 3)
			}
			if ms[0] != ms[1] || (n > 2 && ms[1] != ms[2]) {
				break
			}
			if attempt >= generator.MaxDrawAttempts {
				ms[0], ms[1] = 0, 1
				break
			}
		}
		for i := range lik {
			lik[i] = math.Pow(0.1, float64(ms[i]))
			lt = append(lt, edakit.Sprintf("C%d: prior weight %d, %d mismatching patterns", i+1, w[i], ms[i]))
		}
	} else {
		for i := range lik {
			lik[i] = float64(edakit.Int(5, 95)) / 100
			lt = append(lt, edakit.Sprintf("C%d: prior weight %d, likelihood %s", i+1, w[i], edakit.Num(lik[i], 2)))
		}
	}
	post := make([]float64, n)
	z := 0.0
	for i := range post {
		post[i] = prior[i] * lik[i]
		z += post[i]
	}
	for i := range post {
		post[i] /= z
	}
	best := 0
	for i := range post {
		if post[i] > post[best] {
			best = i
		}
	}
	switch k {
	case 1, 3:
		extra := ""
		if k == 3 {
			extra = " The agent models the likelihood of the observed failures as $0.1^m$ for m mismatching patterns."
		}
		ans := pct(100 * post[best])
		return edakit.Numeric(edakit.Sprintf("A diagnosis agent reads a failure log and ranks candidate defects. Prior weights (proportional, not yet normalized) come from layout-weighted defect statistics; likelihoods from how well each candidate accounts for the log. %s.%s What is the posterior probability of the most probable candidate, in percent? Round to 2 decimal places.", strings.Join(lt, "; "), extra),
			ans, edakit.Sprintf("Posterior is prior times likelihood, normalized: C%d gets %s percent.", best+1, ans), nil)
	default:
		thr := edakit.Pick([]int{80, 90, 95})
		idx := seq(n)
		sort.Slice(idx, func(a, b int) bool { return post[idx[a]] > post[idx[b]] })
		cum, cnt := 0.0, 0
		for _, i := range idx {
			cum += post[i]
			cnt++
			if cum*100 >= float64(thr)-1e-9 {
				break
			}
		}
		var pt []string
		for _, i := range idx {
			pt = append(pt, edakit.Sprintf("C%d %s", i+1, edakit.Num(100*post[i], 1)))
		}
		return edakit.Numeric(edakit.Sprintf("A diagnosis agent ranks candidates for physical failure analysis (prior weights are proportional, posterior is proportional to prior times likelihood): %s. Each candidate sent to failure analysis costs a day of lab time. Sending candidates in order of posterior probability, how many must be sent so that their total posterior is at least %d percent?", strings.Join(lt, "; "), thr),
			edakit.IntStr(cnt), "Posteriors in percent, sorted: "+strings.Join(pt, ", ")+edakit.Sprintf(". The first %d reach %d percent.", cnt, thr), nil)
	}
}

// fastCoverage: sampled estimate (KP 1), confidence half-width (KP 2), required sample size
// (KP 3).
func fastCoverage(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		tot := edakit.Scaled(ctx, 100000, 5000000)
		n := edakit.Pick([]int{1000, 2000, 5000, 10000})
		d := n - edakit.Int(5, n/20)
		ans := pct(100 * float64(d) / float64(n))
		return edakit.Numeric(edakit.Sprintf("Inside an agent's DFT exploration loop, a fast mode fault-simulates a uniform random sample of %d faults out of %d instead of all of them, and %d of the sample are detected. What coverage does the fast mode report, in percent? Round to 2 decimal places.", n, tot, d),
			ans, edakit.Sprintf("The sample proportion %d/%d = %s percent is an unbiased estimate of the full coverage.", d, n, ans), nil)
	case 2:
		n := edakit.Pick([]int{500, 1000, 2000, 5000, 10000})
		p := float64(edakit.Int(900, 995)) / 1000
		h := 1.96 * math.Sqrt(p*(1-p)/float64(n)) * 100
		ans := edakit.Num(h, 3)
		return edakit.Numeric(edakit.Sprintf("A fault-sampled coverage estimate is %s percent from %d sampled faults. Using the normal approximation, what is the half-width of the 95 percent confidence interval, in percentage points? Round to 3 decimal places.", edakit.Num(100*p, 1), n),
			ans, edakit.Sprintf("$1.96\\sqrt{p(1-p)/n} = 1.96\\sqrt{%s \\times %s/%d}$ = %s points.", edakit.Num(p, 3), edakit.Num(1-p, 3), n, ans), nil)
	default:
		p := float64(edakit.Int(900, 995)) / 1000
		h := float64(edakit.Pick([]int{5, 10, 20, 25, 50})) / 100
		n := int(math.Ceil(1.96 * 1.96 * p * (1 - p) / math.Pow(h/100, 2)))
		return edakit.Numeric(edakit.Sprintf("An agent compares two DFT configurations by fault-sampled coverage, expected near %s percent, and needs a 95 percent confidence half-width of at most %s percentage points before trusting the fast estimate instead of running full fault simulation. How many faults must it sample (normal approximation)?", edakit.Num(100*p, 1), edakit.Num(h, 2)),
			edakit.IntStr(n), edakit.Sprintf("$n = \\lceil 1.96^2 p(1-p)/h^2 \\rceil$ with $h = %s$: %d faults. If the two configurations differ by less than about twice the half-width, the agent must recurse to the accurate engine.", edakit.Num(h/100, 4), n), nil)
	}
}

var atpgSkillItems = []edakit.Item{
	{KP: 1, Q: "An agent must generate stuck-at patterns for a block. What should the ATPG skill's manifest tell it?",
		A: "Supported fault models, required inputs (netlist, scan and clock setup, fault list), outputs (patterns, coverage report) and the tunable knobs", D: []string{"The ATPG engine's full source code", "Only the vendor name", "The customer's design data"},
		Why: "Agents choose skills by reading manifests, not source; the manifest must state what the skill does, its typed I/O and what may be tuned."},
	{KP: 1, Q: "Which belongs in the immutable core of an ATPG engine rather than in its tunable surface?",
		A: "The fault simulator's correctness: a pattern reported as detecting a fault must truly detect it", D: []string{"The backtrack limit", "The order in which faults are targeted", "The compaction effort level"},
		Why: "Detection correctness is the invariant coverage numbers rest on; limits, orders and effort levels are parameters and policies."},
	{KP: 1, Q: "An ATPG skill is exposed as a typed tool to other agents. Which output field matters most for a downstream test-time agent?",
		A: "The pattern count and per-pattern shift length, alongside the coverage achieved", D: []string{"The ATPG run's log verbosity", "The source file names of the engine", "The license server address"},
		Why: "Test time is pattern count times shift cycles; a typed result lets the next agent compute cost without parsing logs."},
	{KP: 1, Q: "Why should an agent call the compiled ATPG engine instead of writing its own test generator for each request?",
		A: "The engine is fast, verified and deterministic; the agent's job is selection and configuration, not re-implementing search", D: []string{"Agents cannot generate any code", "ATPG engines have no parameters", "The engine is always slower than generated code"},
		Why: "Selection is a cheap text decision; the expensive, correctness-critical compute stays in the engine."},
	{KP: 2, Q: "A customer asks for fewer patterns at the same coverage. Which change tier should the agent try first?",
		A: "Tier 1: raise the compaction effort and adjust the fault ordering parameters", D: []string{"Tier 3: rewrite PODEM's backtrace", "Tier 2: add a new fault model", "Tier 3: change the fault simulator's value system"},
		Why: "Pattern count responds to existing configuration knobs; the smallest safe tier is configuration."},
	{KP: 2, Q: "A team wants ATPG to target a new defect type described by a cell-aware table. Which tier does this land in?",
		A: "Tier 2: a plug-in fault model at the engine's fault-model extension point", D: []string{"Tier 1: a single parameter", "Tier 3: a change to the D-calculus", "No change: stuck-at already covers it"},
		Why: "A new fault model is a defined extension (activation and propagation conditions) without touching the search core."},
	{KP: 2, Q: "Where should the seam for customer-specific ATPG behaviour usually sit?",
		A: "At the fault-model and fault-ordering interfaces, which receive the netlist context and return targets and priorities", D: []string{"Inside the five-valued gate evaluation", "Inside the netlist parser", "In the license check"},
		Why: "What varies by customer is which faults matter and in what order; gate evaluation is settled mathematics."},
	{KP: 2, Q: "An agent proposes a Tier 3 change to the implication engine to speed up one design. What must happen before it is accepted?",
		A: "Human review of the diff plus full golden and differential regression across many designs", D: []string{"Nothing: speedups are always safe", "Only a run on that one design", "A change to the manifest text only"},
		Why: "Core edits risk silent coverage errors everywhere; they need the strongest gate."},
	{KP: 3, Q: "An agent's reward for an ATPG experiment is coverage alone. What can it learn to game?",
		A: "Declaring hard faults untestable or excluding them, which raises test coverage without better tests", D: []string{"Nothing: coverage cannot be gamed", "Making the netlist larger", "Running ATPG twice"},
		Why: "A reward must check that untestable classifications are proven and that pattern count and fault lists stay comparable."},
	{KP: 3, Q: "Why must an ATPG skill be deterministic and seedable when used by agents?",
		A: "So a result can be reproduced, compared against a baseline and attributed to the change, not to randomness", D: []string{"Deterministic ATPG is always faster", "Agents cannot handle random numbers", "Seeds increase coverage"},
		Why: "Differential testing only means something if the same inputs give the same outputs."},
	{KP: 3, Q: "An ATPG tool skill runs in a sandbox with no network. Why?",
		A: "Design netlists are the customer's IP; the skill must not be able to send them anywhere", D: []string{"Network access slows ATPG", "ATPG needs a GPU", "Sandboxes improve coverage"},
		Why: "The IP boundary is architectural: skills see design data, so they must not have an exfiltration path."},
	{KP: 3, Q: "An agent must choose between a fast structural ATPG skill and a SAT-based one for a block with many aborted faults. What should drive the choice?",
		A: "The manifests' stated strengths: SAT-based ATPG resolves hard and redundant faults that structural search aborts", D: []string{"Always the structural one, because it is first in the list", "Whichever has the shorter manifest", "Always both, on every fault"},
		Why: "Skill selection matches the situation (many aborts) to the capability advertised (proving or finding hard faults)."},
}

var insertionAgentItems = []edakit.Item{
	{KP: 1, Q: "An agent inserts scan into a netlist. Which check proves it did not change the design's function?",
		A: "Combinational equivalence checking between the pre- and post-insertion netlists with scan enable tied to functional mode", D: []string{"Running ATPG and looking at coverage", "Comparing file sizes", "Static timing analysis alone"},
		Why: "With scan enable at 0 the scan muxes are transparent, so the two netlists must be equivalent; a miter-based checker proves it."},
	{KP: 1, Q: "After an agent stitches scan chains, which test catches a chain broken by a wrong connection?",
		A: "A chain integrity (flush) test that shifts a known pattern such as 0011 through each chain and checks it arrives intact", D: []string{"A stuck-at ATPG run only", "Equivalence checking in functional mode", "A March test"},
		Why: "Functional equivalence ignores the scan path; only shifting through the chain proves its continuity and length."},
	{KP: 1, Q: "Which scan design-rule violation must an insertion agent fix before ATPG can use a flop?",
		A: "A flop whose clock or asynchronous reset is not controllable from the test clocks and pins", D: []string{"A flop with high fanout", "A flop on a timing-critical path", "A flop with a long name"},
		Why: "Uncontrollable clocks or resets make scan shifting unpredictable, so the flop cannot be part of a chain until fixed."},
	{KP: 1, Q: "An agent orders scan flops by placement to cut wirelength. What must it still respect?",
		A: "Clock-domain grouping and lockup latches where chains cross domains, and hold timing on the shift path", D: []string{"Alphabetical order of instance names", "The order of flops in the RTL", "Nothing: placement order is always safe"},
		Why: "Physical reordering is a policy choice; domain and timing rules are invariants of a working chain."},
	{KP: 2, Q: "A customer wants test points inserted only in modules that miss coverage. Which tier?",
		A: "Tier 1 or 2: configuration of the test-point policy, or a plug-in selection rule", D: []string{"Tier 3: rewrite the netlist database", "Tier 3: change the fault simulator", "No tier: test points cannot be targeted"},
		Why: "Where to place test points is policy; the insertion mechanics and netlist database are core and need not change."},
	{KP: 2, Q: "Which harness loop fits an agent inserting compression logic?",
		A: "Propose the configuration, insert, run DRC, equivalence and ATPG, compare coverage and pattern counts with baseline, then show the diff to a human", D: []string{"Insert once and ship", "Run ATPG only after tapeout", "Edit the engine core first"},
		Why: "The loop is edit, build, test against goldens, and review; failures feed back to the agent's next proposal."},
	{KP: 2, Q: "An insertion agent proposes adding an observation point on a high-fanout net. What must it re-verify?",
		A: "Functional equivalence (the point must be transparent in functional mode), timing on that net, and the coverage gain", D: []string{"Only the coverage gain", "Only the chain length", "Nothing: observation points are passive"},
		Why: "A test point adds load and logic; it must not change function, must not break timing, and must actually help coverage."},
	{KP: 2, Q: "Why should the agent's scan insertion be reproducible from a configuration file?",
		A: "So the change is versioned, reviewable and reversible, and can be replayed on the next netlist", D: []string{"Configuration files are faster to read", "Netlists cannot be versioned", "Reproducibility lowers test coverage"},
		Why: "Every engine change must be versioned and reversible; a configuration-driven insertion is an auditable artifact."},
	{KP: 3, Q: "An agent inserts scan and all checks pass, but test coverage dropped 0.5 points versus baseline. What is the right action?",
		A: "Treat it as a regression: diagnose which faults became untestable, for example through new tied logic, before promoting", D: []string{"Promote it, since equivalence passed", "Lower the baseline", "Ignore coverage for scan insertion"},
		Why: "Equivalence proves function; it does not prove testability. A coverage drop is a golden-test failure."},
	{KP: 3, Q: "An agent's DFT reward is chain count balance. What could go wrong?",
		A: "It may balance chains by crossing clock domains without lockup latches, creating hold failures in shift", D: []string{"Nothing: balance is always good", "It will reduce coverage to zero", "It will remove scan entirely"},
		Why: "A proxy reward must be paired with the invariants (domain rules, timing) or the agent will trade them away."},
	{KP: 3, Q: "Which evidence should accompany an agent's DFT change shown to a human reviewer?",
		A: "The diff plus equivalence, chain integrity, DRC, coverage and pattern-count deltas against the baseline", D: []string{"The agent's chain of thought only", "A summary saying the change is safe", "The final GDS"},
		Why: "Reviewers approve evidence, not assertions; the golden and differential results are the proof."},
	{KP: 3, Q: "A successful insertion recipe found on one design is promoted into the library. What flows back, and what stays?",
		A: "The recipe and when to use it flow into the shared skill library; the customer's netlist stays on site", D: []string{"The customer's netlist is uploaded for retraining", "Nothing is kept", "The recipe is discarded after use"},
		Why: "The skill flywheel moves know-how, not design data, across the IP boundary."},
}
