package mleda

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// passK is the unbiased pass@k estimate from n samples with c correct.
func passK(n, c, k int) float64 {
	if n-c < k {
		return 1
	}
	return 1 - binom(n-c, k)/binom(n, k)
}

// passAtK: pass@1 (KP 1), unbiased pass@k for one problem (KP 2), and benchmark pass@k averaged
// over problems (KP 3).
func passAtK(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Pick([]int{10, 20, 50, 100, 200})
		var cs []int
		tot := 0.0
		m := edakit.Int(3, 5)
		for i := 0; i < m; i++ {
			c := edakit.Int(0, n)
			cs = append(cs, c)
			tot += float64(c) / float64(n)
		}
		ans := edakit.Num(100*tot/float64(m), 2)
		return edakit.Numeric(edakit.Sprintf("An RTL benchmark run draws %d samples for each of %d problems; the numbers of samples passing the testbench are %s. What is pass@1 for the run, in percent (the average over problems of the fraction of passing samples)? Round to 2 decimal places.", n, m, ints(cs)),
			ans, edakit.Sprintf("Per-problem pass@1 is $c/n$. The average of %s over %d problems is %s percent.", ints(cs), m, ans),
			map[string]string{"n": edakit.IntStr(n)})
	case 2:
		n := edakit.Pick([]int{10, 15, 20, 30, 50})
		k := edakit.Pick([]int{2, 3, 5, 10})
		c := edakit.Int(1, n/2)
		v := passK(n, c, k)
		ans := edakit.Num(v, 4)
		return edakit.Numeric(edakit.Sprintf("For one RTL-generation problem, %d samples were drawn and %d pass the testbench. Using the unbiased estimator $\\text{pass@}k = 1 - \\binom{n-c}{k}/\\binom{n}{k}$, what is pass@%d? Round to 4 decimal places.", n, c, k),
			ans, edakit.Sprintf("$\\binom{%d}{%d} = %s$ and $\\binom{%d}{%d} = %s$. pass@%d $= 1 - %s/%s = %s$: the probability that a random subset of %d samples contains at least one passing sample.", n-c, k, f(binom(n-c, k), 0), n, k, f(binom(n, k), 0), k, f(binom(n-c, k), 0), f(binom(n, k), 0), ans, k),
			map[string]string{"naive": f(1-math.Pow(1-float64(c)/float64(n), float64(k)), 4)})
	default:
		n := edakit.Pick([]int{10, 20})
		k := edakit.Pick([]int{3, 5})
		cs := []int{edakit.Int(0, n/2), edakit.Int(0, n), edakit.Int(1, 3)}
		tot := 0.0
		var terms []string
		for _, c := range cs {
			v := passK(n, c, k)
			tot += v
			terms = append(terms, f(v, 4))
		}
		ans := edakit.Num(100*tot/3, 2)
		return edakit.Numeric(edakit.Sprintf("A three-problem RTL benchmark draws %d samples per problem; %s samples pass on the three problems. Using the unbiased estimator per problem and averaging, what is the benchmark's pass@%d, in percent? Round to 2 decimal places.", n, ints(cs), k),
			ans, edakit.Sprintf("Per-problem pass@%d: %s. Average: %s percent. A problem with even one passing sample contributes a large pass@%d, which is why pass@%d far exceeds pass@1.", k, strings.Join(terms, ", "), ans, k, k),
			map[string]string{"k": edakit.IntStr(k)})
	}
}

// functionalCorrectness: testbench escape probability (KP 1), vectors needed for a target
// escape rate (KP 2), and the inflation of a reported pass rate (KP 3).
func functionalCorrectness(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		d := float64(edakit.Int(1, 50)) / 1000
		t := edakit.Pick([]int{10, 20, 50, 100, 200, 500})
		e := math.Pow(1-d, float64(t))
		ans := edakit.Num(100*e, 2)
		return edakit.Numeric(edakit.Sprintf("A generated arithmetic module has a bug that changes its output on a fraction %s of all input vectors. The benchmark's testbench applies %d independent uniformly random vectors. What is the probability, in percent, that the buggy module passes the testbench? Round to 2 decimal places.", f(d, 3), t),
			ans, edakit.Sprintf("Each vector misses the bug with probability $1 - %s$. All %d miss: $(1 - %s)^{%d} = %s$ percent. An equivalence check against the reference would find it with certainty.", f(d, 3), t, f(d, 3), t, ans),
			map[string]string{"per_vector": f(1-d, 3)})
	case 2:
		d := float64(edakit.Int(1, 50)) / 1000
		eps := edakit.Pick([]float64{0.05, 0.01, 0.001})
		t := int(math.Ceil(math.Log(eps) / math.Log(1-d)))
		return edakit.Numeric(edakit.Sprintf("A bug corrupts the output on a fraction %s of random input vectors. How many independent random vectors must a testbench apply so that the probability of the bug escaping is at most %s? Give the smallest integer.", f(d, 3), f(eps, 3)),
			edakit.IntStr(t), edakit.Sprintf("Need $(1 - %s)^T \\le %s$, so $T \\ge \\ln %s / \\ln(1 - %s) = %s$, giving %d. Bugs that hit rarer inputs need many more vectors; random testing never gives a proof.", f(d, 3), f(eps, 3), f(eps, 3), f(d, 3), f(math.Log(eps)/math.Log(1-d), 3), t),
			map[string]string{"ratio": f(math.Log(eps)/math.Log(1-d), 3)})
	default:
		c := float64(edakit.Int(20, 70)) / 100
		e := float64(edakit.Int(5, 50)) / 100
		rep := c + (1-c)*e
		ans := edakit.Num(100*rep, 2)
		return edakit.Numeric(edakit.Sprintf("In a model's RTL samples, %s percent are functionally correct. The benchmark's testbenches are weak: an incorrect sample passes with probability %s. What pass rate, in percent, does the benchmark report? Round to 2 decimal places.", f(100*c, 0), f(e, 2)),
			ans, edakit.Sprintf("Reported $= c + (1 - c)e = %s + %s \\times %s = %s$ percent, an overstatement of %s points.", f(c, 2), f(1-c, 2), f(e, 2), ans, f(100*(rep-c), 2)),
			map[string]string{"inflation": f(100*(rep-c), 2)})
	}
}

var benchmarkItems = []edakit.Item{
	{KP: 1, Q: "How does VerilogEval judge a generated solution?",
		A: "By simulating it against a reference solution's behavior with a testbench", D: []string{"By checking that the code compiles", "By comparing its text to the reference", "By asking human graders to rate it"},
		Why: "VerilogEval problems come with reference designs and testbenches; a sample counts as correct when simulation agrees with the reference."},
	{KP: 1, Q: "Where do the VerilogEval problems come from, broadly?",
		A: "Exercises from an online Verilog practice collection, with natural-language problem descriptions", D: []string{"Industrial SoC blocks from tapeouts", "Randomly generated netlists", "Formal verification competitions"},
		Why: "The problems are adapted from a public tutorial-style exercise set, which bears on both difficulty and contamination risk."},
	{KP: 1, Q: "What does RTLLM add beyond checking whether generated code works?",
		A: "Design-level tasks with evaluation of syntax, functionality and design quality after synthesis", D: []string{"Only syntax checks", "Formal proofs of every design", "Measurements of fabricated silicon"},
		Why: "RTLLM targets larger design descriptions and also considers PPA-style quality, not only pass or fail."},
	{KP: 1, Q: "A benchmark result says pass@1 is 60 percent. What is being estimated?",
		A: "The probability that a single sampled solution to a random benchmark problem passes its tests", D: []string{"The share of problems the model has seen in training", "The fraction of generated code that compiles", "The probability that the model's best answer is formally correct"},
		Why: "pass@1 averages, over problems, the fraction of samples that pass the benchmark's functional tests."},
	{KP: 2, Q: "Two papers report pass@1 on the same benchmark, one with greedy decoding and one averaging samples at temperature 0.8. Can the numbers be compared directly?",
		A: "No: the decoding settings change what is measured, so they must match", D: []string{"Yes, pass@1 is always the same quantity", "Yes, temperature only affects runtime", "No, because the benchmark changed between papers"},
		Why: "Greedy decoding gives one deterministic sample; sampled pass@1 is an expectation over a distribution. Both are legitimate but not interchangeable."},
	{KP: 2, Q: "A model scores 90 percent pass@10 and 40 percent pass@1. What does the gap tell a tool builder?",
		A: "The model often finds a correct answer among several tries, so a checker that selects among samples is essential", D: []string{"The model is 90 percent reliable on single tries", "The benchmark is broken", "The model never produces correct code on the first try"},
		Why: "pass@k assumes an oracle picks the passing sample. In a tool, that oracle has to be a real verification step."},
	{KP: 2, Q: "Why is a benchmark of short exercises a weak predictor of a model's usefulness on a production block?",
		A: "Production RTL involves large interfaces, existing code, conventions and PPA targets that short exercises do not test", D: []string{"Short exercises use a different language", "Production blocks cannot be simulated", "Exercises are always harder than production code"},
		Why: "Benchmarks measure a narrow task distribution. Real use needs evaluation on representative tasks."},
	{KP: 2, Q: "A benchmark problem's testbench checks only a few directed cases. What kind of error does this create in reported scores?",
		A: "False passes: incorrect solutions that happen to agree on the tested cases are counted as correct", D: []string{"False failures for correct solutions", "Compile errors", "No error, since directed tests are complete"},
		Why: "Weak tests inflate pass rates. Stronger testbenches or equivalence checks against the reference reduce false passes."},
	{KP: 3, Q: "Which evaluation design best separates a model's RTL ability from memorization of a public benchmark?",
		A: "A held-out set of new problems written after the model's training cutoff, with hidden reference solutions", D: []string{"Re-running the public benchmark with a new seed", "Increasing k in pass@k", "Shortening the problem descriptions"},
		Why: "Fresh, unpublished problems cannot be in the training data."},
	{KP: 3, Q: "A team wants a benchmark for agentic RTL tools that edit existing code. What should it measure beyond pass@k on standalone modules?",
		A: "Correctness of edits within a larger codebase, regressions in unchanged behavior, and PPA after synthesis", D: []string{"The length of the generated code", "Token throughput of the model", "The number of comments in the code"},
		Why: "Agentic tools act on existing designs, so the benchmark must check both the intended change and the absence of collateral damage."},
	{KP: 3, Q: "A vendor reports that its RTL assistant improves productivity 'by up to 10x'. How should a tool designer treat this?",
		A: "As a claim against an unstated baseline until the task set, metric and baseline are known", D: []string{"As an established result", "As a lower bound", "As a property of the language model alone"},
		Why: "'Up to' figures without baselines describe the best case of an unknown comparison."},
	{KP: 3, Q: "Why should RTL benchmark scores be reported with the number of samples and a confidence interval?",
		A: "With few problems and samples, the estimate has sampling error large enough to reverse rankings between models", D: []string{"Because pass@k is always exact", "Because confidence intervals are required by Verilog", "Because more samples change the problems"},
		Why: "A benchmark of modest size gives a noisy estimate. Small differences between models may not be significant."},
}

var contaminationItems = []edakit.Item{
	{KP: 1, Q: "What is benchmark contamination for an RTL-generation model?",
		A: "Benchmark problems or their solutions appear in the model's training data", D: []string{"The benchmark's testbenches have bugs", "The model produces code with syntax errors", "The benchmark uses a different Verilog standard"},
		Why: "If the model saw the solutions, a high score can reflect recall rather than ability."},
	{KP: 1, Q: "Why are RTL benchmarks built from public tutorial exercises especially exposed to contamination?",
		A: "The exercises and many solutions are posted online, where web-scale training corpora can collect them", D: []string{"Tutorial exercises are too easy to grade", "Public exercises cannot be simulated", "Tutorials use only VHDL"},
		Why: "Publicly posted problems and solutions are exactly what large crawls ingest."},
	{KP: 1, Q: "A model trained on a large crawl of public Verilog scores far higher on a benchmark than on similar new problems. What is the most likely explanation?",
		A: "Contamination: the benchmark items or close variants were in the crawl", D: []string{"The new problems are mislabeled", "The model improved after training", "The benchmark testbenches are stronger"},
		Why: "A gap between a public set and fresh, equivalent problems is the signature of memorization."},
	{KP: 1, Q: "Contamination also happens in learned EDA predictors. Which situation is an example?",
		A: "The test design for a congestion model was also used, with other flow settings, in training", D: []string{"The model was trained on a GPU", "The training set had more tiles than the test set", "The labels came from a router"},
		Why: "Leaking the test design through variants lets the model recognize it rather than generalize."},
	{KP: 2, Q: "How can an evaluator test whether a model has memorized a benchmark problem?",
		A: "Check whether it completes the reference solution verbatim from a partial prefix, or test on semantically equivalent rewrites", D: []string{"Ask the model whether it has seen the problem", "Measure how long the model takes to answer", "Increase the sampling temperature"},
		Why: "Verbatim completion and a score drop on rewritten but equivalent problems are behavioral evidence of memorization."},
	{KP: 2, Q: "Why is n-gram overlap between training data and benchmark solutions an imperfect contamination test?",
		A: "Paraphrased or renamed solutions evade it, while short common idioms create false matches", D: []string{"N-grams cannot be computed on code", "Overlap always means contamination", "Benchmarks have no solutions"},
		Why: "Text overlap catches exact copies but misses renamed variants and flags generic boilerplate."},
	{KP: 2, Q: "A team rewrites benchmark descriptions in new words and renames signals, keeping the hardware behavior. The model's score drops sharply. What does that suggest?",
		A: "The model was relying on surface features it had memorized, not on understanding the task", D: []string{"The rewritten problems are harder in hardware terms", "The testbenches broke", "The model's temperature changed"},
		Why: "Equivalent problems should be equally solvable. A sharp drop points to memorization of the originals."},
	{KP: 2, Q: "Why is it not enough for a model vendor to state that a benchmark was 'excluded from training'?",
		A: "Copies, solutions and close variants posted elsewhere may still have been ingested", D: []string{"Exclusion is impossible to implement", "Benchmarks are always in training sets", "Statements are never accurate"},
		Why: "Excluding the official files does not exclude the many reposted solutions. Behavioral tests are still needed."},
	{KP: 3, Q: "Which benchmark practice most reduces contamination risk going forward?",
		A: "Keep a private, regularly refreshed problem set with hidden references, evaluated through a submission process", D: []string{"Publish all solutions with the problems", "Use only problems from textbooks", "Report pass@100 instead of pass@1"},
		Why: "Problems that are never published and are refreshed over time cannot be trained on."},
	{KP: 3, Q: "A paper fine-tunes on synthetic Verilog generated by a larger model, then evaluates on a public benchmark. What extra contamination path exists?",
		A: "The larger model may have memorized the benchmark and passed solutions into the synthetic data", D: []string{"Synthetic data cannot contain benchmark content", "Fine-tuning removes contamination", "Synthetic data always lowers scores"},
		Why: "Distillation can launder memorized benchmark content into a training set that looks clean."},
	{KP: 3, Q: "An agentic RTL tool is tuned by repeatedly checking its score on the benchmark and adjusting prompts. What is the risk?",
		A: "Overfitting to the benchmark through repeated use, so the score no longer estimates performance on new tasks", D: []string{"No risk, since the model weights do not change", "The benchmark testbenches will change", "The tool will stop compiling code"},
		Why: "Using a test set for tuning turns it into a training set, even without gradient updates."},
	{KP: 3, Q: "How should a company evaluate an RTL assistant before deploying it on its own designs?",
		A: "On an internal set of representative tasks that never left the company, scored by its own verification flow", D: []string{"By the public benchmark leaderboard alone", "By the vendor's reported numbers", "By the number of parameters in the model"},
		Why: "Internal tasks are uncontaminated by construction and represent the real distribution of work."},
}

var shiftItems = []edakit.Item{
	{KP: 1, Q: "A timing predictor trained on CPU blocks is applied to a network-on-chip design with very long wires. What kind of failure is most likely?",
		A: "Distribution shift: the feature-label relation on the new design is unlike the training designs", D: []string{"Overfitting on the new design", "A label error in the training set", "Exact predictions, because timing physics is universal"},
		Why: "The model learned correlations present in its training designs. A new design style changes those correlations."},
	{KP: 1, Q: "A congestion model is trained at one technology node and applied at a more advanced node with different metal stacks. What changes that the model cannot know?",
		A: "Routing resources per gcell and the rules that turn demand into violations", D: []string{"The netlist format", "The number of designs", "The meaning of the word congestion"},
		Why: "Labels depend on the node's capacity and rules. Without retraining or adaptation, the mapping from features to labels is wrong."},
	{KP: 1, Q: "Which split tests the generalization a production tool needs from a learned predictor?",
		A: "Train on some designs and test on designs never seen in any form", D: []string{"Random split of tiles across all designs", "Random split of flow runs across all designs", "Test on the training designs at a new seed"},
		Why: "A tool meets new designs. Only a design-level split estimates that."},
	{KP: 1, Q: "A predictor's inputs include features computed by one tool's placer. The team changes placer version. What should they check first?",
		A: "Whether the feature distributions shifted, because the model was trained on the old placer's outputs", D: []string{"Whether the GPU drivers changed", "Whether the labels are still integers", "Nothing; features are tool-independent"},
		Why: "Covariate shift from a changed upstream tool can degrade a model even when labels are unchanged."},
	{KP: 3, Q: "A model's error on a new design is three times its validation error. Which response is best for a tool?",
		A: "Detect the shift from input statistics or uncertainty and fall back to the accurate engine for that design", D: []string{"Report the validation error as the expected error", "Retrain on the new design's test labels and report the new error as generalization", "Ignore the model's error since it is only a heuristic"},
		Why: "A tool should know when its fast model is out of distribution and recurse to exact analysis there."},
	{KP: 3, Q: "Fine-tuning a pretrained predictor on a few runs of a new design is proposed. What must the evaluation do to stay honest?",
		A: "Report accuracy on held-out runs or regions of that design, not on the fine-tuning data", D: []string{"Report the fine-tuning loss", "Report accuracy on the original training designs", "Skip evaluation because fine-tuning always helps"},
		Why: "Adaptation is legitimate, but the test data must stay separate from the adaptation data."},
	{KP: 3, Q: "Which signal can flag a distribution shift at run time without labels?",
		A: "Input features or embeddings far from the training distribution, or high ensemble disagreement", D: []string{"The model's prediction being large", "The design having many cells", "The tool version string"},
		Why: "Out-of-distribution detection uses the inputs and the model's uncertainty, which are available before the expensive label."},
	{KP: 3, Q: "Why does a model trained on many flow-parameter variants of eight designs still generalize poorly to a ninth?",
		A: "The variants share each design's structure, so the model has effectively seen only eight independent designs", D: []string{"Eight designs exceed the model's capacity", "Parameter variants corrupt labels", "Generalization depends only on the number of samples"},
		Why: "Diversity in designs, not the number of correlated samples, limits generalization to new designs."},
	{KP: 1, Q: "An IR-drop model is trained on designs with one power-grid style and applied to a design with a much sparser grid. Which input-output relation changed?",
		A: "How a given current map turns into voltage drop, because the grid resistance differs", D: []string{"The definition of current", "The file format of the power map", "Nothing, since the current map is the only input"},
		Why: "The label depends on the grid as well as the currents. If grid features are missing or out of range, the learned mapping no longer holds."},
	{KP: 1, Q: "Training designs all use one clock frequency; the new design runs at twice that frequency. Why might a slack predictor fail even if the netlists look similar?",
		A: "Slack depends on the constraints, so labels shift with the clock period even for similar logic", D: []string{"Netlists change format with frequency", "Clock frequency affects only power", "Predictors ignore slack sign"},
		Why: "Labels are functions of netlist, placement and constraints. A constraint outside the training range is a shift."},
	{KP: 3, Q: "A team reports a predictor's accuracy averaged over leave-one-design-out folds, but one held-out design has five times the average error. What should the report add?",
		A: "The per-design spread or worst-case error, since a tool meets designs one at a time", D: []string{"Nothing, the average is enough", "Only the training loss", "The model's parameter count"},
		Why: "An average hides the tail. A tool designer needs to know how bad a single unlucky design can be."},
	{KP: 3, Q: "Which practice limits the damage when a deployed predictor meets an out-of-distribution design?",
		A: "Using the model only to steer search, with exact engines confirming every accepted result", D: []string{"Raising the model's confidence threshold to 100 percent", "Retraining on the test labels before reporting", "Disabling the accurate engines to save time"},
		Why: "If exact engines gate acceptance, a shifted model costs runtime or quality but cannot produce a wrong signoff result."},
}

// designSplitChoice: the generalization error to report from leave-one-design-out results.
func designSplitChoice(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Int(4, 6)
	names := []string{"aes", "ibex", "jpeg", "swerv", "ariane", "bp"}
	errs := make([]float64, m)
	var rows []string
	for i := range errs {
		errs[i] = float64(edakit.Int(60, 250)) / 10
		rows = append(rows, edakit.Sprintf("%s %s ps", names[i], f(errs[i], 1)))
	}
	rnd := minF(errs) * float64(edakit.Int(30, 70)) / 100
	mean := sumF(errs) / float64(m)
	sorted := append([]float64(nil), errs...)
	sort.Float64s(sorted)
	q := edakit.Sprintf("A slack predictor is trained on %d designs. A random split of all endpoints gives a test MAE of %s ps. Leave-one-design-out (train on the others, test on the held-out design) gives MAE %s. What MAE best estimates the error on a new design?", m, f(rnd, 2), strings.Join(rows, ", "))
	why := edakit.Sprintf("The leave-one-design-out mean, %s ps, measures error on unseen designs. The random-split figure %s ps is optimistic because test endpoints share designs with training endpoints; the spread across designs shows the risk.", f(mean, 2), f(rnd, 2))
	return numChoice(q, mean, 2, "ps", why, rnd, sorted[0], (rnd+mean)/2, sorted[m-1])
}

// rankingMetrics: concordant pairs (KP 1), Kendall tau (KP 2), and top-k overlap (KP 3) between
// a predictor's ranking of flow configurations and the true ranking.
func rankingMetrics(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 7)
	actual := distinctVals(n, 5, 250, 0)
	pred := make([]float64, n)
	perm := distinctVals(n, 50, 3000, 1)
	// mostly agree: blend order with a little shuffle
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool { return actual[order[a]] < actual[order[b]] })
	sw := edakit.Int(1, 3)
	ranks := make([]int, n)
	for i, o := range order {
		ranks[o] = i
	}
	for s := 0; s < sw; s++ {
		i := edakit.Int(0, n-2)
		a, b := order[i], order[i+1]
		ranks[a], ranks[b] = ranks[b], ranks[a]
		order[i], order[i+1] = b, a
	}
	sort.Float64s(perm)
	for i := range pred {
		pred[i] = perm[ranks[i]]
	}
	stem := edakit.Sprintf("A predictor scores %d flow configurations before running them. True post-route TNS magnitudes (ps) are %s and the predicted values are %s, in the same order.", n, floats(actual, 0), floats(pred, 1))
	con, dis := 0, 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if (actual[i]-actual[j])*(pred[i]-pred[j]) > 0 {
				con++
			} else {
				dis++
			}
		}
	}
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+" How many of the configuration pairs are ordered the same way by the prediction and the truth (concordant pairs)?",
			edakit.IntStr(con), edakit.Sprintf("Of the $\\binom{%d}{2} = %d$ pairs, %d are concordant and %d discordant.", n, n*(n-1)/2, con, dis),
			map[string]string{"discordant": edakit.IntStr(dis)})
	case 2:
		tau := float64(con-dis) / float64(n*(n-1)/2)
		ans := edakit.Num(tau, 4)
		return edakit.Numeric(stem+" What is Kendall's tau, $(C - D)/\\binom{n}{2}$? Round to 4 decimal places.",
			ans, edakit.Sprintf("$C = %d$, $D = %d$, $\\binom{%d}{2} = %d$. $\\tau = %d/%d = %s$. For choosing configurations, rank agreement matters more than absolute error.", con, dis, n, n*(n-1)/2, con-dis, n*(n-1)/2, ans),
			map[string]string{"concordant": edakit.IntStr(con)})
	default:
		k := edakit.Int(2, 3)
		ia := topK(actual, k)
		ip := topK(pred, k)
		ov := 0
		for x := range ia {
			if ip[x] {
				ov++
			}
		}
		return edakit.Numeric(stem+edakit.Sprintf(" Lower TNS magnitude is better. The team will fully run only the %d configurations the predictor ranks best. How many of the truly best %d are among them?", k, k),
			edakit.IntStr(ov), edakit.Sprintf("True best %d: configurations %s. Predicted best %d: %s. Overlap: %d.", k, keysStr(ia), k, keysStr(ip), ov),
			map[string]string{"k": edakit.IntStr(k)})
	}
}

// distinctVals returns n distinct values drawn from [lo, hi] and divided by 10^places.
func distinctVals(n, lo, hi, places int) []float64 {
	seen := map[int]bool{}
	var out []float64
	for tries := 0; len(out) < n; tries++ {
		v := edakit.Int(lo, hi)
		if tries >= generator.MaxDrawAttempts {
			v = hi + tries
		}
		if !seen[v] {
			seen[v] = true
			out = append(out, float64(v)/math.Pow(10, float64(places)))
		}
	}
	return out
}

// topK returns the indices of the k smallest values.
func topK(xs []float64, k int) map[int]bool {
	idx := make([]int, len(xs))
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool { return xs[idx[a]] < xs[idx[b]] })
	out := map[int]bool{}
	for _, i := range idx[:k] {
		out[i] = true
	}
	return out
}

// keysStr lists 1-based indices in ascending order.
func keysStr(m map[int]bool) string {
	var ks []int
	for k := range m {
		ks = append(ks, k+1)
	}
	sort.Ints(ks)
	return ints(ks)
}
