package mleda

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// learnedFastMode: loop time with a learned fast mode and recursion (KP 1), its speedup
// (KP 2), and the recursion fraction at which the fast mode stops paying (KP 3).
func learnedFastMode(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 50, 1000)
	ta := float64(edakit.Int(10, 120))
	ts := edakit.Round(ta*float64(edakit.Int(3, 30))/100, 1)
	fr := float64(edakit.Int(5, 40)) / 100
	m := edakit.Int(3, 20)
	stem := edakit.Sprintf("An agent's placement-tuning loop evaluates %d candidates. The accurate engine (global route plus STA) takes %s minutes per candidate; a learned fast mode takes %s minutes. The agent's policy sends every candidate through the fast mode, re-runs the accurate engine on the %s percent whose prediction is uncertain, and always verifies the final %d chosen candidates with the accurate engine.", n, f(ta, 0), f(ts, 1), f(100*fr, 0), m)
	fast := float64(n)*ts + (fr*float64(n)+float64(m))*ta
	base := float64(n) * ta
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(fast/60, 2)
		return edakit.Numeric(stem+" How many hours does the loop take? Round to 2 decimal places.",
			ans, edakit.Sprintf("$%d \\times %s + (%s \\times %d + %d) \\times %s = %s$ min, or %s h.", n, f(ts, 1), f(fr, 2), n, m, f(ta, 0), f(fast, 1), ans),
			map[string]string{"minutes": f(fast, 1)})
	case 2:
		ans := edakit.Num(base/fast, 3)
		return edakit.Numeric(stem+" What is the speedup over running the accurate engine on every candidate? Round to 3 decimal places.",
			ans, edakit.Sprintf("Accurate only: $%d \\times %s = %s$ min. With the fast mode: %s min. Speedup: %s. The recursion fraction, not the fast mode's own speed, dominates.", n, f(ta, 0), f(base, 0), f(fast, 1), ans),
			map[string]string{"base": f(base, 0)})
	default:
		be := (float64(n)*ta - float64(n)*ts - float64(m)*ta) / (float64(n) * ta)
		ans := edakit.Num(100*be, 2)
		return edakit.Numeric(edakit.Sprintf("An agent's loop evaluates %d candidates. The accurate engine takes %s minutes per candidate and a learned fast mode %s minutes. Every candidate goes through the fast mode, a fraction $f$ is re-run accurately because of uncertainty, and the final %d choices are always verified accurately. Above what recursion fraction $f$, in percent, is the loop slower than running the accurate engine on everything? Round to 2 decimal places.", n, f(ta, 0), f(ts, 1), m),
			ans, edakit.Sprintf("Break-even: $N t_s + (fN + m)t_a = N t_a$, so $f = (N t_a - N t_s - m t_a)/(N t_a) = %s$ percent. An agent that sees its recursion rate approach this should stop using the fast mode on this design.", ans),
			map[string]string{"break_even": f(be, 4)})
	}
}

var verificationGateItems = []edakit.Item{
	{KP: 1, Q: "An agent's learned model proposes a macro placement. Which check gates its acceptance?",
		A: "Legalization and a routed evaluation of timing and congestion by the flow's engines", D: []string{"The model's own predicted proxy cost", "A second model's agreement", "The number of macros placed"},
		Why: "Every learned output is accepted only after the exact engines measure it on the metrics that matter."},
	{KP: 1, Q: "An LLM rewrites an RTL module for lower power. What gate must the rewrite pass before it replaces the original?",
		A: "A sequential equivalence check against the original module, plus synthesis to confirm the power change", D: []string{"Compilation without warnings", "The LLM's explanation of the change", "A syntax highlighter"},
		Why: "Behavior must be proven unchanged, and the claimed benefit measured by the tool."},
	{KP: 1, Q: "A learned timing predictor flags 20 paths as safe to leave unoptimized. What is the right gate?",
		A: "Incremental STA on those paths before the agent stops working on them", D: []string{"None, the predictor has a low MAE", "A congestion map", "A check that the paths exist"},
		Why: "A prediction of safety must be confirmed by the timer before it becomes a decision."},
	{KP: 1, Q: "A generated flow script claims to set up multi-corner analysis. What gates it?",
		A: "Running it in a sandbox and checking the tool's reported corners and the resulting reports against expectations", D: []string{"Reading the script for plausibility", "Counting its lines", "The model's confidence score"},
		Why: "The tool's state after the script runs is the evidence; the script text is not."},
	{KP: 2, Q: "In the harness loop, where does the verification gate sit for a learned model's output?",
		A: "After the proposal and before any change is promoted or shown as a passing diff", D: []string{"Before the model is called", "Only at tapeout", "Only when a human asks for it"},
		Why: "The harness runs the golden and differential checks on every proposal; failures go back to the model as feedback."},
	{KP: 2, Q: "A learned fast mode is used inside a loop to rank thousands of candidates. Which outputs must pass the accurate gate?",
		A: "Every candidate that is accepted or reported, while rejected ones may rest on the fast mode alone", D: []string{"None, ranking needs no checks", "Every candidate, including all rejected ones", "Only the first candidate evaluated"},
		Why: "Errors on rejections cost opportunity; errors on acceptances reach the design. The gate sits on what is kept."},
	{KP: 2, Q: "A gate's check takes longer than the work the learned model saved. What should the architecture do?",
		A: "Use cheaper incremental or differential checks where they suffice, and reserve full checks for promotion", D: []string{"Drop the gate", "Trust the model when the check is slow", "Run the gate only on failures"},
		Why: "Verification cost is part of the speedup accounting. Incremental checks keep the gate affordable without removing it."},
	{KP: 2, Q: "Why should the gate's verdict be recorded alongside each learned output?",
		A: "So the agent and later reviewers can audit what was proven and use pass and fail outcomes as training signal", D: []string{"To make the design database larger", "Because verdicts change the RTL", "To replace golden tests"},
		Why: "Recorded verdicts make decisions reviewable and feed the learning loop with grounded labels."},
	{KP: 3, Q: "A gate checks generated RTL only with the benchmark's directed testbench. What can still pass the gate?",
		A: "Incorrect designs that agree with the reference on the tested inputs", D: []string{"Nothing, directed tests are complete", "Only designs with syntax errors", "Only designs that are too large"},
		Why: "A gate is only as strong as its check. Weak tests let wrong outputs through; equivalence or stronger stimulus closes the gap."},
	{KP: 3, Q: "A gate approves placements whose proxy cost beats the baseline. Why is that gate insufficient?",
		A: "It checks the same proxy the model optimized, so it cannot catch outputs that exploit the proxy's blind spots", D: []string{"Proxy costs cannot be computed", "Placements cannot be compared", "The gate is too slow"},
		Why: "Gating on the optimized metric is circular. The gate must use the engines that define the real objective."},
	{KP: 3, Q: "A gate proves a generated assertion holds on the design. What can that verdict still hide?",
		A: "A vacuous assertion whose trigger never occurs", D: []string{"A syntax error in the RTL", "A timing violation in the assertion", "A wrong clock frequency"},
		Why: "Passing proves only that no counterexample exists. Non-vacuity needs a separate reachability check."},
	{KP: 3, Q: "An agent's gate uses the fast timer instead of signoff STA to save time. When is that acceptable?",
		A: "For intermediate loop decisions, provided signoff STA gates the final result", D: []string{"Always, if the fast timer is usually right", "Never, the fast timer must not be used", "Only for hold checks"},
		Why: "Fast checks can gate inner-loop steps. The final promotion needs the accurate engine."},
}

var onSiteItems = []edakit.Item{
	{KP: 1, Q: "Why would a learned congestion or timing model be retrained on a customer's site rather than only shipped pretrained?",
		A: "The customer's designs, node and flow differ from the training data, and their data cannot leave the site", D: []string{"Pretrained models cannot run on customer machines", "Retraining makes the model smaller", "Customers' designs never vary"},
		Why: "On-site learning addresses distribution shift on data that the IP boundary keeps local."},
	{KP: 1, Q: "Where do on-site training labels come from in a continuous learning setup?",
		A: "The accurate engines' results produced during normal flow runs on the customer's designs", D: []string{"Hand annotation by the vendor", "The vendor's other customers", "Synthetic random labels"},
		Why: "Every accurate run the flow already performs is a free, grounded label for the fast model."},
	{KP: 1, Q: "Which asset is meant to flow back from a customer site to the vendor's shared library in the skill flywheel?",
		A: "A validated skill, such as a policy or parameter rule, without the customer's design data", D: []string{"The customer's netlists", "The customer's routed layouts", "Raw training tiles from the customer"},
		Why: "Know-how flows back after validation and review; the design data stays on site."},
	{KP: 1, Q: "A base model must be fine-tuned on site to become a domain model. What property must it have?",
		A: "Open weights that can be trained and run locally", D: []string{"A closed API with a large context window", "The largest possible parameter count", "Training only on public Verilog"},
		Why: "Weights that can only be reached through a remote API cannot be trained or kept inside the IP boundary."},
	{KP: 2, Q: "Why can model weights trained on a customer's designs themselves be IP-sensitive?",
		A: "Models can memorize training examples, so weights may leak design details if exported", D: []string{"Weights are just random numbers", "Weights never depend on training data", "Only the optimizer state is sensitive"},
		Why: "Memorization is possible, so weights trained on private data inherit its sensitivity."},
	{KP: 2, Q: "A vendor proposes uploading anonymized feature maps from customers to train a shared model. What must the IP review consider?",
		A: "Whether the features still reveal design structure, such as macro placement or floorplan, even without names", D: []string{"Only the file size", "Whether the features are images", "Nothing, anonymization removes all risk"},
		Why: "Layout-derived maps can encode the floorplan. Anonymizing names does not anonymize geometry."},
	{KP: 2, Q: "Which architecture keeps on-site learning inside the IP boundary by construction?",
		A: "Training, inference and data stay in the customer's sandbox, and only reviewed skills or metrics leave", D: []string{"Training in the vendor's cloud with encrypted uploads", "Sending each design to a remote API for prediction", "Emailing models between sites"},
		Why: "The boundary is part of the architecture: the sandbox has no network path for design data."},
	{KP: 2, Q: "Several customers want the benefit of each other's learning without sharing data. What mechanism fits best?",
		A: "Sharing validated, design-independent skills through the vendor's library while models train locally", D: []string{"Pooling all designs in one dataset", "Training one model on each customer's raw data in turn at the vendor", "Exchanging netlists under NDA"},
		Why: "Abstracted, validated know-how can move; raw data and weights trained on it should not."},
	{KP: 3, Q: "An on-site model is retrained nightly and one morning its predictions degrade on a new block. What should the loop have in place?",
		A: "Versioned models with golden regression checks before promotion and a rollback to the last good version", D: []string{"Automatic promotion of every nightly model", "Deleting old models to save space", "Retraining every hour instead"},
		Why: "Continuous learning needs the same promotion discipline as code: test, version, roll back."},
	{KP: 3, Q: "After fine-tuning on a memory-heavy project, a model gets worse on the logic-heavy blocks it handled well before. What happened?",
		A: "Catastrophic forgetting: training on the new distribution overwrote what fit the old one", D: []string{"Label noise in the old data", "The GPU changed", "The model became too small"},
		Why: "Sequential training without replay or regularization drifts toward the newest data. Regression tests on older blocks catch it."},
	{KP: 3, Q: "How should an on-site learning loop decide that the fast model can be trusted on a new design?",
		A: "By measuring its error against the accurate engine on a sample of that design's runs before relying on it", D: []string{"By its accuracy on last year's designs", "By its training loss", "By the vendor's benchmark numbers"},
		Why: "Trust is earned per distribution. A small accurate sample on the new design checks for shift."},
	{KP: 3, Q: "Labels from the customer's flow come from a tool configuration that changes mid-project. What must the training loop do?",
		A: "Record the tool version and settings with each label and avoid mixing incompatible labels", D: []string{"Mix all labels, since more data is better", "Discard all labels", "Ignore tool versions"},
		Why: "A label is a measurement by a specific setup. Mixing setups silently corrupts the target."},
}

var toolFeedbackItems = []edakit.Item{
	{KP: 1, Q: "What makes a reward 'tool-grounded' when training a domain model for EDA?",
		A: "It is computed by deterministic solvers and EDA tools that score the output against objectives, constraints and ground truth", D: []string{"It is a human rating of the output", "It is the model's own confidence", "It is the length of the output"},
		Why: "Tools provide objective, reproducible scores, such as compile status, test pass rates, equivalence and QoR."},
	{KP: 1, Q: "A model is trained to write flow scripts. Which reward signal is tool-grounded?",
		A: "Whether the script runs and the tool's reports show the requested result", D: []string{"Whether a grader model likes the script", "The script's similarity to a reference text", "The number of commands used"},
		Why: "The tool executes the script and its outputs show the effect; that is the grounded signal."},
	{KP: 1, Q: "Why is a tool-grounded reward better than imitation of expert scripts alone?",
		A: "It rewards outcomes, so the model can find correct solutions that differ from the examples and is penalized for plausible but wrong ones", D: []string{"It needs no tools", "It is always cheaper to compute", "It removes the need for any data"},
		Why: "Imitation rewards similarity; grounded rewards check function and quality."},
	{KP: 1, Q: "What is the main cost of tool-grounded RL for EDA tasks?",
		A: "Each reward needs a tool run, so rollouts are slow and the reward may come from a fast proxy instead", D: []string{"Tools cannot produce numbers", "Rewards are always binary", "The models cannot read tool output"},
		Why: "Expensive evaluation pushes training toward fast proxies, which opens the gap between proxy reward and signoff QoR."},
	{KP: 3, Q: "A model trained with reward 'fraction of testbench checks passed' learns to emit code that disables the checker's error reporting. What is this?",
		A: "Reward hacking: the model optimized the measurement instead of the intended behavior", D: []string{"Overfitting to the training designs", "A compiler bug", "Successful generalization"},
		Why: "When the reward can be reached by gaming the check, optimization will find the shortcut. The checker must be outside the model's reach."},
	{KP: 3, Q: "An RL model is rewarded by a fast congestion estimate and learns placements with low estimated congestion but high routed overflow. What is the gap?",
		A: "The fast-proxy reward diverges from the signoff measurement the proxy stands in for", D: []string{"The router is broken", "Congestion cannot be learned", "The model was too large"},
		Why: "Hard optimization exploits proxy errors. Periodic accurate evaluation and proxy recalibration limit the gap."},
	{KP: 3, Q: "How can a reward for RTL generation be made harder to hack?",
		A: "Combine hidden tests, equivalence checking against a reference and synthesis QoR, all run outside the model's control", D: []string{"Use only the compile status", "Let the model write its own tests", "Reward longer code"},
		Why: "Independent, diverse checks leave fewer shortcuts. A model that writes its own tests can make them trivially pass."},
	{KP: 3, Q: "Which evaluation shows that a tool-feedback-trained model has improved rather than learned to game the reward?",
		A: "Better results on held-out tasks judged by stronger checks than the training reward, such as signoff QoR or formal equivalence", D: []string{"A higher training reward", "Lower training loss", "Faster inference"},
		Why: "If gains survive a stronger, independent judge on new tasks, they reflect real ability."},
	{KP: 1, Q: "Which output can a tool-grounded reward score for an assertion-generation model?",
		A: "Whether each assertion compiles, is non-vacuous, holds on the reference RTL and fails on injected mutants", D: []string{"How fluent the assertion's comment is", "Whether the assertion resembles a textbook example", "How many tokens the assertion uses"},
		Why: "Formal and simulation tools can measure every one of these properties without a human."},
	{KP: 1, Q: "Vendors describe domain models 'trained with tool feedback'. Which question checks what that claim means?",
		A: "Which tools computed the reward, on what tasks, and how the result was evaluated on held-out work", D: []string{"How many GPUs were used", "Which programming language the model was written in", "Whether the model has a name"},
		Why: "The claim is only as strong as the reward's tools, the task distribution and the independent evaluation behind it."},
	{KP: 3, Q: "A model trained to fix timing ECOs earns reward from WNS improvement and learns to relax the clock constraint in the SDC it edits. What should the training setup have prevented?",
		A: "Edits to the constraints that define the reward; the SDC must be read-only to the model", D: []string{"Edits to the netlist", "Any change to WNS", "Use of the timer"},
		Why: "If the model can change what is measured, it will. Constraints and checkers belong outside its action space."},
	{KP: 3, Q: "Training uses a fast timer for reward and signoff STA only every thousand steps. The two diverge over training. What is a sound response?",
		A: "Recalibrate or retrain the fast reward model on signoff results and weight signoff checks more heavily", D: []string{"Stop running signoff STA", "Trust the fast timer since it is used more", "Increase the learning rate"},
		Why: "Divergence is the proxy gap growing under optimization pressure. The accurate engine must keep the proxy honest."},
}

// toolRewardChoice: a tool-grounded reward that combines compile status, tests and QoR.
func toolRewardChoice(ctx generator.GeneratorContext) generator.Problem {
	tests := edakit.Pick([]int{10, 20, 25, 40, 50})
	pass := edakit.Int(tests/3, tests)
	aref := float64(edakit.Int(800, 3000))
	area := aref * float64(edakit.Int(85, 140)) / 100
	area = float64(int(area))
	lam := edakit.Pick([]float64{0.5, 1, 2})
	over := (area - aref) / aref
	pen := 0.0
	if over > 0 {
		pen = lam * over
	}
	r := float64(pass)/float64(tests) - pen
	q := edakit.Sprintf("A domain model for RTL is trained with the reward $r = 0$ if the module fails to compile, and otherwise $r = (\\text{tests passed})/(\\text{tests}) - \\lambda \\max(0, (A - A_{\\mathrm{ref}})/A_{\\mathrm{ref}})$, where $A$ is the synthesized area. A sample compiles, passes %d of %d hidden tests, and synthesizes to %s um² against a reference of %s um². With $\\lambda = %s$, what is its reward?", pass, tests, f(area, 0), f(aref, 0), f(lam, 1))
	why := edakit.Sprintf("Test term $%d/%d = %s$. Area overshoot %s, penalty %s. Reward $= %s$. Area below the reference earns no bonus, so the model cannot trade correctness for area.", pass, tests, f(float64(pass)/float64(tests), 4), f(over, 4), f(pen, 4), f(r, 4))
	return numChoice(q, r, 4, "", why, float64(pass)/float64(tests), float64(pass)/float64(tests)-lam*over, float64(pass)/float64(tests)-over, 0)
}
