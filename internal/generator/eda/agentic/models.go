package agentic

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var rolesItems = []edakit.Item{
	{KP: 1, Q: "Which role can a closed model reached only through a hosted API fill?",
		A: "Driving the harness, where the deployment allows prompts to leave the site", D: []string{"Running continuous selection on-site", "Becoming the per-customer domain model", "All three roles"},
		Why: "Roles two and three need fine-tuning and local hosting, which a remote-only API model cannot provide."},
	{KP: 1, Q: "What capability separates models most in the harness role?",
		A: "Correct, minimal diffs against existing source, not just fresh generation", D: []string{"Writing long essays", "Image understanding", "Knowing the most EDA trivia"},
		Why: "The harness model edits large existing codebases under strict rules."},
	{KP: 1, Q: "A harness model drifts and stops calling tools correctly after 30 turns. Which requirement does it fail?",
		A: "Reliable tool calling over a long multi-turn loop", D: []string{"Open weights", "A permissive license", "Small size"},
		Why: "The edit, build, test and fix loop runs many turns."},
	{KP: 1, Q: "Why does the harness model need a long context window?",
		A: "To hold the relevant headers, the file being edited and the build log at once", D: []string{"To store every customer design", "To memorize the skill library", "To speed up compilation"},
		Why: "Real engine sources and logs are large."},
	{KP: 2, Q: "Why is the continuous selection model usually smaller than the harness model?",
		A: "It runs constantly with low latency close to the data and is fine-tuned on selection records", D: []string{"Selection is unimportant", "Small models are always more accurate", "It never needs fine-tuning"},
		Why: "Frequency and latency favor a compact, specialized model."},
	{KP: 2, Q: "Why must the selection model run inside the customer's environment?",
		A: "It reads design-derived metrics continuously and must be fine-tuned there", D: []string{"Because it edits the core", "Because it compiles code", "Because APIs are slower to type"},
		Why: "Its inputs are design data, which stays on-site."},
	{KP: 2, Q: "Which requirement does the selection role add beyond tool calling?",
		A: "Weights the deployment controls, so it can be fine-tuned and run locally", D: []string{"Image generation", "A larger context than the harness model", "Internet access"},
		Why: "Fine-tuning and local hosting are the gate."},
	{KP: 2, Q: "A team plans a closed API model for continuous selection at a strict customer site. What is the problem?",
		A: "It cannot be fine-tuned and run on-site, so it cannot fill this role", D: []string{"It is too accurate", "It cannot read manifests", "It cannot output text"},
		Why: "The selection role needs local, tunable weights."},
	{KP: 3, Q: "What does a candidate need to become the on-site domain model?",
		A: "Open weights and a license allowing commercial use, fine-tuning and shipping a fine-tuned derivative", D: []string{"Only a high coding benchmark score", "A hosted API with low latency", "A small vocabulary"},
		Why: "The domain model is distilled and fine-tuned, then deployed into customer environments."},
	{KP: 3, Q: "Why does model size matter for the domain-model role in particular?",
		A: "Parameter count drives fine-tuning and training cost directly", D: []string{"Larger models cannot call tools", "Smaller models cannot be fine-tuned", "Size determines the license"},
		Why: "Training spend scales with size, and this role is trained."},
	{KP: 3, Q: "Which question should be asked before shipping a fine-tuned model to a customer abroad?",
		A: "Whether the license and export rules permit redistributing the fine-tuned derivative there", D: []string{"Whether the model can write poems", "Whether the customer has a GUI", "Whether the model is the newest release"},
		Why: "Redistribution and export terms gate the deployment."},
	{KP: 3, Q: "What is the common mistake when evaluating a model for agent-first EDA?",
		A: "Asking whether it is good enough as one question, instead of per role", D: []string{"Testing it on real engines", "Reading its license", "Measuring compile-pass rate"},
		Why: "Requirements differ sharply by role."},
}

var abstractionItems = []edakit.Item{
	{KP: 1, Q: "What does a model-agnostic harness talk to?",
		A: "One internal tool-call interface, with a thin adapter per model behind it", D: []string{"Each model's API directly throughout the code", "Only one specific model", "The model's weights"},
		Why: "The harness logic never names a model."},
	{KP: 1, Q: "Where are the harness's tools (read a manifest, propose an edit, run a build) defined?",
		A: "Once, with typed arguments, independent of any model", D: []string{"Separately for each model", "Inside the model's prompt only", "In the engine core"},
		Why: "Tools are part of the harness, not the model."},
	{KP: 1, Q: "A team swaps the harness model. Which properties must survive unchanged?",
		A: "The tier policy, sandbox and test gates", D: []string{"The chat template", "The tool-call syntax", "The tokenizer"},
		Why: "Safety lives in the harness, so it is independent of the model."},
	{KP: 1, Q: "Why should the harness's safety logic never depend on which model is attached?",
		A: "So swapping models cannot weaken safety", D: []string{"Because models are always safe", "Because safety is the model's job", "To save memory"},
		Why: "Model-agnostic safety is what makes swapping safe."},
	{KP: 2, Q: "Which detail does a per-model adapter absorb?",
		A: "The chat template and the tool-call format", D: []string{"The golden tests", "The engine's legality rules", "The skill manifests"},
		Why: "Models differ in surface format, which the adapter translates."},
	{KP: 2, Q: "A model emits a malformed tool call. What should the adapter do?",
		A: "Reject it with an error the harness can feed back to the model", D: []string{"Guess the intended call and run it", "Crash the harness", "Execute it as a shell command"},
		Why: "Malformed calls become feedback, never guesses."},
	{KP: 2, Q: "How is a self-hosted open-weights model typically made available to the harness?",
		A: "Served by an inference runtime such as vLLM or SGLang behind the adapter", D: []string{"Linked directly into the placer", "Run inside the golden tests", "Through the vendor's public API"},
		Why: "A serving runtime exposes the local model to the adapter."},
	{KP: 2, Q: "Two models have different context limits and stop conditions. Where is this handled?",
		A: "In each model's adapter", D: []string{"In the engine core", "In the skill manifests", "In the sandbox"},
		Why: "Model-specific limits are adapter concerns."},
	{KP: 3, Q: "How should a team decide whether to swap the harness model?",
		A: "Rerun the skill-editing benchmark with the candidate behind the same harness and compare", D: []string{"Read a generic coding leaderboard", "Pick the newest release", "Ask the current model"},
		Why: "Swaps are decided empirically on the real task."},
	{KP: 3, Q: "One model makes C++ edits while a smaller fine-tuned model runs selection. Is this a compromise?",
		A: "No: mixing models per role is a design goal", D: []string{"Yes: one model should do everything", "Yes: it breaks the harness", "No, but it requires two harnesses"},
		Why: "Swappability makes per-role choices natural."},
	{KP: 3, Q: "Why is a model swap versioned and reversible like any other change?",
		A: "A new model can change the behavior of every agent that uses it", D: []string{"Models are code", "Swaps are rare", "The license requires it"},
		Why: "Model changes have a wide blast radius."},
	{KP: 3, Q: "A customer with stricter security requires a different base model. What does the abstraction layer allow?",
		A: "Running a different model for that customer behind the same harness", D: []string{"Nothing: all customers must share one model", "Rewriting the harness for that customer", "Removing the sandbox"},
		Why: "Per-customer model choice is a configuration of the adapter."},
}

func passAtK(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := 10 * edakit.Int(10, 40)
		c := edakit.Int(n/4, n-5)
		r := 100 * float64(c) / float64(n)
		q := edakit.Sprintf("A model produces %d edits on a skill-editing benchmark and %d of them build. What is its compile-pass rate in percent? Round to 2 decimal places.", n, c)
		return edakit.Numeric(q, edakit.Num(r, 2), edakit.Sprintf("$%d/%d = %s\\%%$.", c, n, edakit.Num(r, 2)), nil)
	case 2:
		n := edakit.Scaled(ctx, 8, 20)
		k := edakit.Int(2, 5)
		c := edakit.Int(1, n-k)
		v := 1 - binom(n-c, k)/binom(n, k)
		q := edakit.Sprintf("A model makes %d independent attempts at one task and %d succeed. Using the unbiased estimator $1 - \\binom{n-c}{k}/\\binom{n}{k}$, what is pass@%d? Round to 4 decimal places.", n, c, k)
		exp := edakit.Sprintf("$1 - \\binom{%d}{%d}/\\binom{%d}{%d} = 1 - %s/%s = %s$.", n-c, k, n, k, edakit.Num(binom(n-c, k), 0), edakit.Num(binom(n, k), 0), edakit.Num(v, 4))
		return edakit.Numeric(q, edakit.Num(v, 4), exp, nil)
	}
	k := edakit.Int(2, 5)
	m := edakit.Int(3, 4)
	var ps []float64
	var pt []string
	mean, pool := 0.0, 0.0
	for i := 0; i < m; i++ {
		p := float64(edakit.Int(1, 9)) / 10
		ps = append(ps, p)
		pt = append(pt, edakit.Num(p, 1))
		mean += 1 - math.Pow(1-p, float64(k))
		pool += p
	}
	mean /= float64(m)
	pool /= float64(m)
	if edakit.Int(0, 1) == 0 {
		q := edakit.Sprintf("A benchmark has %d tasks whose per-attempt success probabilities are %s. Computing pass@%d per task as $1 - (1-p)^k$ and averaging over tasks, what is the benchmark's pass@%d? Round to 4 decimal places.", m, strings.Join(pt, ", "), k, k)
		return edakit.Numeric(q, edakit.Num(mean, 4), edakit.Sprintf("Per-task values averaged give %s; pooling first would give %s.", edakit.Num(mean, 4), edakit.Num(1-math.Pow(1-pool, float64(k)), 4)), nil)
	}
	diff := (1 - math.Pow(1-pool, float64(k))) - mean
	q := edakit.Sprintf("A benchmark has %d tasks whose per-attempt success probabilities are %s. By how much does pooling the success rate first, then computing $1 - (1-\\bar p)^%d$, overstate the per-task average of pass@%d? Round to 4 decimal places.", m, strings.Join(pt, ", "), k, k)
	exp := edakit.Sprintf("Pooled: %s. Per-task average: %s. Difference %s, positive because $1-(1-p)^k$ is concave.", edakit.Num(1-math.Pow(1-pool, float64(k)), 4), edakit.Num(mean, 4), edakit.Num(diff, 4))
	return edakit.Numeric(q, edakit.Num(diff, 4), exp, nil)
}

func editCorrectness(ctx generator.GeneratorContext) generator.Problem {
	r := func() float64 { return float64(edakit.Int(50, 98)) / 100 }
	switch kp(ctx) {
	case 1:
		a, b, c := r(), r(), r()
		v := a * b * c
		q := edakit.Sprintf("On a skill-editing benchmark, %s%% of a model's edits build, %s%% of those pass the golden tests, and %s%% of those have a clean differential. What fraction of edits are correct? Round to 4 decimal places.", edakit.Num(100*a, 0), edakit.Num(100*b, 0), edakit.Num(100*c, 0))
		return edakit.Numeric(q, edakit.Num(v, 4), edakit.Sprintf("$%s \\times %s \\times %s = %s$.", edakit.Num(a, 2), edakit.Num(b, 2), edakit.Num(c, 2), edakit.Num(v, 4)), nil)
	case 2:
		a1 := float64(edakit.Int(88, 98)) / 100
		a2 := float64(edakit.Int(70, 87)) / 100
		b1, c1 := float64(edakit.Int(50, 70))/100, float64(edakit.Int(60, 80))/100
		b2, c2 := float64(edakit.Int(70, 90))/100, float64(edakit.Int(75, 95))/100
		for t := 0; a2*b2*c2 <= a1*b1*c1+0.01 && t < generator.MaxDrawAttempts; t++ {
			a1, a2 = float64(edakit.Int(88, 98))/100, float64(edakit.Int(70, 87))/100
			b1, c1 = float64(edakit.Int(50, 70))/100, float64(edakit.Int(60, 80))/100
			b2, c2 = float64(edakit.Int(70, 90))/100, float64(edakit.Int(75, 95))/100
		}
		if a2*b2*c2 <= a1*b1*c1+0.01 {
			a1, b1, c1, a2, b2, c2 = 0.95, 0.6, 0.75, 0.85, 0.75, 0.85
		}
		v1, v2 := a1*b1*c1, a2*b2*c2
		d := 100 * (v2 - v1)
		q := edakit.Sprintf("Model A builds %s%% of its edits, with %s%% of those passing goldens and %s%% of those differential-clean. Model B builds %s%%, then %s%% and %s%%. By how many percentage points does B's skill-edit correctness exceed A's? Round to 2 decimal places.",
			edakit.Num(100*a1, 0), edakit.Num(100*b1, 0), edakit.Num(100*c1, 0), edakit.Num(100*a2, 0), edakit.Num(100*b2, 0), edakit.Num(100*c2, 0))
		exp := edakit.Sprintf("A: %s. B: %s. Difference %s points, although A has the higher compile rate.", edakit.Num(v1, 4), edakit.Num(v2, 4), edakit.Num(d, 2))
		return edakit.Numeric(q, edakit.Num(d, 2), exp, nil)
	}
	n := 10 * edakit.Int(3, 10)
	s := edakit.Int(1, n/4)
	v := 100 * float64(n-s) / float64(n)
	q := edakit.Sprintf("Of %d edits that pass every test, %d touched files outside the tier their request declared. What percentage of the passing edits are correct and in scope? Round to 2 decimal places.", n, s)
	return edakit.Numeric(q, edakit.Num(v, 2), edakit.Sprintf("$(%d - %d)/%d = %s\\%%$.", n, s, n, edakit.Num(v, 2)), nil)
}

var benchItems = []edakit.Item{
	{KP: 1, Q: "Why do popular fresh-code benchmarks separate strong models poorly for skill editing?",
		A: "Scores cluster near the top, and the tasks do not test editing large existing engines under invariants", D: []string{"They are too hard for any model", "They test only C++", "They run too slowly"},
		Why: "Saturation and the wrong skill both shrink the differences that matter here."},
	{KP: 1, Q: "Two models score within a point of each other on a function-writing benchmark. What can you conclude about their skill-edit correctness on a real placer?",
		A: "Little: they can differ widely on that task", D: []string{"They will be within a point there too", "The higher scorer will win", "Both will fail"},
		Why: "Generic benchmarks under-separate models on this task."},
	{KP: 1, Q: "Which skill does editing a placer's move evaluator require that a fresh-code benchmark rarely tests?",
		A: "Reading an existing interface and respecting invariants not stated in the prompt", D: []string{"Writing a sorting function", "Formatting output", "Explaining code in prose"},
		Why: "Real edits sit inside a large codebase with unstated contracts."},
	{KP: 1, Q: "Which number should decide the go or no-go for a harness model?",
		A: "Skill-edit correctness and compile-pass rate on the real engines", D: []string{"Its rank on a general leaderboard", "Its parameter count", "Its release date"},
		Why: "Decide on the task the model will actually do."},
	{KP: 2, Q: "What does one task in a skill-editing benchmark contain?",
		A: "A plain-language change request on a real engine with its tier, scope, goldens, differential spec and hidden tests", D: []string{"A function signature and three unit tests", "A multiple-choice question about EDA", "A design to place by hand"},
		Why: "Tasks mirror real change requests and their acceptance checks."},
	{KP: 2, Q: "Why build benchmark tasks on open-source engines such as OpenROAD and OpenSTA?",
		A: "They are real, large engines that can be shared and inspected", D: []string{"They are small toy programs", "They have no tests", "They are written in Python"},
		Why: "Realistic and shareable engines make a credible benchmark."},
	{KP: 2, Q: "Which metric set fits a skill-editing benchmark?",
		A: "Compile-pass rate, correctness through the funnel, scope violations, attempts per success and QoR deltas", D: []string{"Lines of code generated", "Tokens per second only", "Number of comments written"},
		Why: "These measure whether edits are correct, in scope and useful."},
	{KP: 2, Q: "Why should a benchmark include tasks of every tier?",
		A: "Configuration, plug-in and core tasks test different skills and failure modes", D: []string{"To make the benchmark longer", "Because Tier 3 tasks are the easiest", "Because tiers do not matter"},
		Why: "A model that handles Tier 1 well may fail at plug-ins."},
	{KP: 3, Q: "Benchmark tasks were drawn from a public repository's history. What risk must be addressed?",
		A: "Contamination: the tasks may be in the model's training data", D: []string{"The tasks are too recent", "The tasks are too small", "The tasks cannot compile"},
		Why: "Held-out tasks written after the training cutoff, or on private engines, avoid it."},
	{KP: 3, Q: "A report says a model achieves 0.72 pass@5. What else must it state to be interpretable?",
		A: "The number of samples per task, the estimator used and a confidence interval", D: []string{"The model's color scheme", "Nothing else", "The authors' names"},
		Why: "Sampling details determine how much the number means."},
	{KP: 3, Q: "A vendor reports up to 30% better QoR versus an unnamed baseline. How should this be treated?",
		A: "As a best-case claim, not a benchmark result", D: []string{"As a reproducible result", "As an average over all designs", "As a guarantee"},
		Why: "Without baseline, workload and method, it cannot be evaluated."},
	{KP: 3, Q: "Which detail must accompany a QoR improvement claim from an agent benchmark?",
		A: "The baseline, the designs, the seeds and whether the result holds at signoff", D: []string{"The model's prompt length", "The number of GPUs", "The benchmark's name only"},
		Why: "QoR claims follow the same discipline as any differential."},
}
