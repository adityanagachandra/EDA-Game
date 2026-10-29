package mleda

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var rtlGenerationItems = []edakit.Item{
	{KP: 1, Q: "An LLM generates a Verilog module from a natural-language specification. What does the model actually compute?",
		A: "A sequence of tokens sampled from a learned distribution conditioned on the prompt", D: []string{"A synthesized netlist checked against the spec", "A proof that the module meets the spec", "A simulation trace of the module"},
		Why: "The model is a next-token predictor. Any correctness has to be established afterwards by tools."},
	{KP: 1, Q: "Why do RTL-generation studies usually sample several completions per problem at a nonzero temperature?",
		A: "Different samples succeed on different problems, and a checker can pick a working one among them", D: []string{"A single sample is always correct at temperature zero", "Temperature changes the Verilog standard used", "Several samples are needed to compile one"},
		Why: "Sampling gives diversity; with an automatic checker, more samples raise the chance that at least one passes."},
	{KP: 1, Q: "A prompt specifies a FIFO without saying whether reset is synchronous or asynchronous. What should a designer expect from the generated RTL?",
		A: "The model picks one convention, which may not match the intended interface", D: []string{"The model raises a compile error", "The model generates both versions in one module", "The convention does not affect behavior"},
		Why: "Ambiguity in the specification becomes an arbitrary choice in the output. Specifications for generation must be as precise as for a human."},
	{KP: 1, Q: "What does fine-tuning a general LLM on a corpus of open-source Verilog mainly change?",
		A: "It shifts the model's output distribution toward the syntax and idioms in that corpus, including its bugs", D: []string{"It guarantees synthesizable output", "It teaches the model to run a simulator internally", "It makes the model's output formally verified"},
		Why: "Fine-tuning imitates the corpus. Quality of the corpus, not just its size, limits the result."},
	{KP: 2, Q: "A generation loop sends each candidate to a Verilog compiler and returns error messages to the model. What does this loop improve directly?",
		A: "Syntactic and elaboration correctness", D: []string{"Functional correctness against the specification", "Timing closure of the synthesized design", "Equivalence to a reference model"},
		Why: "A compiler checks that code is legal Verilog. It says nothing about whether it does what was intended."},
	{KP: 2, Q: "Which feedback signal lets an RTL generator improve functional correctness without a human?",
		A: "Simulation results against a testbench, or an equivalence check against a reference model", D: []string{"The compiler's warning count", "The length of the generated module", "The model's own explanation of its code"},
		Why: "Function is checked by executing or proving the design against an independent specification of its behavior."},
	{KP: 2, Q: "A team gives the model the port list and parameter names of neighboring modules as context. What problem does this address?",
		A: "Interface mismatches: the generated module must connect to existing signals with the right names and widths", D: []string{"Timing violations inside the module", "Synthesis runtime", "Simulation speed"},
		Why: "Retrieving real interface definitions grounds generation in the actual design instead of plausible guesses."},
	{KP: 2, Q: "Generated RTL passes simulation but synthesizes to three times the area of a hand-written version. What has the evaluation missed?",
		A: "Design quality: area, timing and power after synthesis, which functional tests do not measure", D: []string{"Nothing, because functional correctness is the only goal", "A syntax error", "An equivalence failure"},
		Why: "RTL quality also means PPA. Synthesizing generated code and comparing PPA is a separate check."},
	{KP: 3, Q: "Generated code uses blocking assignments inside a clocked always block that models a shift register. What is the risk?",
		A: "Simulation can show the register shifting by more than one stage per clock, mismatching the intended hardware", D: []string{"A compile error in every simulator", "Lower clock frequency", "The module will not synthesize"},
		Why: "Blocking assignments update immediately, so ordering inside the block changes behavior. It compiles cleanly, which is why it slips past syntax checks."},
	{KP: 3, Q: "A generated combinational block omits a default in a case statement. What appears after synthesis?",
		A: "An inferred latch that holds a value the specification never intended", D: []string{"A faster multiplexer", "A compile error", "An extra flip-flop driven by the clock"},
		Why: "Incomplete assignments in combinational logic imply storage. Lint and synthesis reports catch it; the model may not."},
	{KP: 3, Q: "Why is a generated module that is correct for one parameter value not necessarily correct for others?",
		A: "Tests and examples often exercise one configuration, while width-dependent bugs appear only at other values", D: []string{"Parameters are ignored by synthesis", "Parameterized modules cannot be simulated", "Verilog does not allow parameters"},
		Why: "Generators and their tests commonly fix one width. Correctness has to be checked across the parameter range that will be used."},
	{KP: 3, Q: "In a design flow, where should LLM-generated RTL be allowed to enter without human review?",
		A: "Nowhere unchecked: it enters only after lint, simulation or formal checks against an independent specification, and is reviewed like any contribution", D: []string{"Directly into the tapeout branch if it compiles", "Into signoff if the model's confidence is high", "Anywhere, since the model was trained on RTL"},
		Why: "Generated code is a proposal. The verification flow, not the generator, establishes that it can be used."},
}

// repairLoops: success within r independent attempts (KP 1), expected attempts with a cap
// (KP 2), and functional success when compiling is only the first gate (KP 3).
func repairLoops(ctx generator.GeneratorContext) generator.Problem {
	p := float64(edakit.Int(10, 70)) / 100
	r := edakit.Scaled(ctx, 2, 10)
	switch kp(ctx) {
	case 1:
		s := 1 - math.Pow(1-p, float64(r))
		ans := edakit.Num(100*s, 2)
		return edakit.Numeric(edakit.Sprintf("An LLM produces a compiling Verilog module on a given task with probability %s per attempt. A harness feeds compiler errors back and allows up to %d attempts; treat attempts as independent. What is the probability, in percent, that some attempt compiles? Round to 2 decimal places.", f(p, 2), r),
			ans, edakit.Sprintf("$1 - (1 - %s)^{%d} = 1 - %s = %s$ percent. Real attempts are correlated, so this is optimistic.", f(p, 2), r, f(math.Pow(1-p, float64(r)), 6), ans),
			map[string]string{"fail_all": f(math.Pow(1-p, float64(r)), 6)})
	case 2:
		e := (1 - math.Pow(1-p, float64(r))) / p
		ans := edakit.Num(e, 3)
		return edakit.Numeric(edakit.Sprintf("A harness calls an LLM and then the compiler, stopping at the first module that compiles or after %d attempts. Each attempt compiles independently with probability %s. What is the expected number of attempts (compiler calls)? Round to 3 decimal places.", r, f(p, 2)),
			ans, edakit.Sprintf("$E[N] = \\sum_{i=0}^{%d} (1 - p)^i = (1 - (1 - p)^{%d})/p = %s$. Without a cap it would be $1/p = %s$.", r-1, r, ans, f(1/p, 3)),
			map[string]string{"uncapped": f(1/p, 3)})
	default:
		q := float64(edakit.Int(20, 80)) / 100
		s := 1 - math.Pow(1-p*q, float64(r))
		ans := edakit.Num(100*s, 2)
		return edakit.Numeric(edakit.Sprintf("In each round of a repair loop, the model's module compiles with probability %s, and a compiling module passes the functional testbench with probability %s. The harness allows %d rounds, each with a fresh attempt; treat rounds as independent. What is the probability, in percent, of ending with a module that passes the testbench? Round to 2 decimal places.", f(p, 2), f(q, 2), r),
			ans, edakit.Sprintf("Per round success $= %s \\times %s = %s$. Over %d rounds: $1 - (1 - %s)^{%d} = %s$ percent, well below the compile-success rate of %s percent.", f(p, 2), f(q, 2), f(p*q, 4), r, f(p*q, 4), r, ans, f(100*(1-math.Pow(1-p, float64(r))), 2)),
			map[string]string{"per_round": f(p*q, 4)})
	}
}

var assertionItems = []edakit.Item{
	{KP: 1, Q: "An LLM is asked to turn the sentence 'every request is granted within 3 cycles' into SVA. Which property matches it?",
		A: "req |-> ##[1:3] gnt", D: []string{"req |-> gnt", "gnt |-> ##[1:3] req", "req ##3 gnt"},
		Why: "The implication starts at each request and requires a grant 1 to 3 cycles later. The sequence req ##3 gnt is not an implication and requires an exact delay."},
	{KP: 1, Q: "Why is assertion generation an attractive task for LLMs compared with RTL generation?",
		A: "Assertions are short, and each can be checked against the design with a formal or simulation tool", D: []string{"Assertions never need to be checked", "Assertions replace the RTL", "Formal tools cannot read hand-written assertions"},
		Why: "Short, checkable outputs suit generate-and-verify loops. A wrong assertion is caught by running it."},
	{KP: 1, Q: "What input should an assertion generator be given so that its properties capture intent rather than implementation?",
		A: "The specification text and interface descriptions, not only the RTL", D: []string{"Only the RTL of the module", "Only the synthesis log", "Only the waveform of one simulation"},
		Why: "Assertions derived from the RTL alone tend to restate what the code does, including its bugs."},
	{KP: 1, Q: "An engineer wants a property that the FIFO never overflows. Which SVA form expresses a safety property of this kind?",
		A: "assert property (@(posedge clk) disable iff (rst) !(push && full));", D: []string{"cover property (@(posedge clk) push && full);", "assume property (@(posedge clk) full);", "assert property (@(posedge clk) push |-> full);"},
		Why: "The safety property says a push never happens while full. A cover only shows reachability, and an assume constrains the environment instead of checking the design."},
	{KP: 2, Q: "A generated assertion passes formal verification on the design. What else must be checked before trusting it?",
		A: "That it is not vacuous: its antecedent can actually occur", D: []string{"That it compiles in a second simulator", "That it is shorter than ten tokens", "That the model was confident"},
		Why: "An implication whose trigger never fires is trivially true. A cover property on the antecedent shows it is reachable."},
	{KP: 2, Q: "How can a team measure whether generated assertions would catch real bugs?",
		A: "Inject mutations into the RTL and count how many mutants some assertion fails on", D: []string{"Count the number of assertions generated", "Measure the average assertion length", "Check that all assertions pass on the original RTL"},
		Why: "Mutation analysis measures detection power. Passing on correct RTL is necessary but says nothing about strength."},
	{KP: 2, Q: "A generated assertion fails under formal verification. What are the two possible explanations a verification engineer must distinguish?",
		A: "A real design bug, or an assertion that misstates the intended behavior", D: []string{"A compiler bug or a waveform bug", "The clock is too fast or the reset is too long", "The model is too small or the prompt too short"},
		Why: "A counterexample shows disagreement between design and property. Which one is wrong is decided against the specification."},
	{KP: 2, Q: "A formal tool proves a generated property only with an assumption the LLM also generated. What is the risk?",
		A: "An over-constrained assumption can exclude the inputs where the design fails, making the proof meaningless", D: []string{"Assumptions slow down simulation", "Assumptions are ignored by formal tools", "The property becomes a cover"},
		Why: "Assumptions restrict the environment. Generated assumptions must be reviewed, or checked as assertions on the neighboring block."},
	{KP: 3, Q: "A model writes 'a |=> b' where the specification says b must hold in the same cycle as a. What went wrong?",
		A: "It used the non-overlapping implication, which checks b one cycle later; the overlapping form is a |-> b", D: []string{"Nothing, the two forms are equivalent", "It should have used a cover property", "It should have used ##0 in the antecedent"},
		Why: "|=> starts the consequent on the next cycle. Off-by-one timing is a common, compile-clean error in generated SVA."},
	{KP: 3, Q: "Generated assertions all pass and none are vacuous, but every one mirrors an RTL expression exactly. Why are they weak?",
		A: "They restate the implementation, so they cannot detect a bug that the implementation shares", D: []string{"They will slow down synthesis", "They are illegal SVA", "They make formal proofs impossible"},
		Why: "An assertion is useful when it encodes intent independently of the code that implements it."},
	{KP: 3, Q: "Which metric is least informative about the quality of an assertion generator?",
		A: "The fraction of generated assertions that compile", D: []string{"The fraction of injected bugs caught", "The fraction of assertions that are non-vacuous and pass on correct RTL", "The fraction of specification clauses covered by some assertion"},
		Why: "Compiling is the lowest bar. Detection power, non-vacuity and intent coverage say whether the assertions are worth running."},
	{KP: 3, Q: "Where does an assertion generator fit in an agent-driven verification flow?",
		A: "It proposes properties that a formal or simulation engine checks, with failing and vacuous ones triaged before they join the regression", D: []string{"It replaces the formal engine", "It signs off the block when its assertions compile", "It edits the RTL until assertions pass"},
		Why: "Generated properties are inputs to verification engines. Letting a model edit the design until its own properties pass would be circular."},
}

var scriptItems = []edakit.Item{
	{KP: 1, Q: "Why are EDA flow scripts a natural target for LLM generation?",
		A: "They are programs against a documented command API, and their effect can be checked by running the tool", D: []string{"Scripts do not affect QoR", "Tools accept any text as a script", "Scripts need no testing"},
		Why: "Tcl and Python flow scripts call tool commands; a run shows whether they work and what they produced."},
	{KP: 1, Q: "A user asks an assistant to 'report the 10 worst setup paths after placement' in OpenROAD. What does a correct answer consist of?",
		A: "Calls to the timing-reporting command with the right options, run on the placed database with constraints loaded", D: []string{"A description of what setup paths are", "A Liberty file", "A modified SDC with tighter constraints"},
		Why: "Script generation turns intent into concrete tool commands in the right state of the flow."},
	{KP: 1, Q: "Which context most improves an LLM's flow scripts for a specific tool version?",
		A: "That version's command reference and example scripts retrieved into the prompt", D: []string{"A general description of chip design", "The model's memory of older versions", "The netlist of the design"},
		Why: "Commands and options change across versions. Grounding in the actual reference reduces invented or outdated commands."},
	{KP: 1, Q: "An LLM writes a Python script that uses a tool's database API to resize cells. What must the script respect that a Tcl report command does not?",
		A: "The database's invariants and the need to update dependent analyses, such as incremental timing, after edits", D: []string{"Nothing, since Python is safer than Tcl", "The script must avoid loops", "The script must run before reading the netlist"},
		Why: "Editing scripts change design state. They must use the API correctly and keep analyses consistent."},
	{KP: 2, Q: "A generated script calls a command option that does not exist in the installed tool. How should a harness catch this before running a long flow?",
		A: "Validate each command and option against the tool's own command registry or help output, or run a fast dry-run", D: []string{"Trust the model if the script looks plausible", "Run the full flow and check the next day", "Ask the model whether the option exists"},
		Why: "Invented options are a common failure. The tool itself is the authority on what it accepts."},
	{KP: 2, Q: "A generated flow script runs without errors. Which check shows that it did what the user asked?",
		A: "Inspecting the resulting design state or reports for the requested effect, such as the expected utilization or path count", D: []string{"The exit status alone", "The script's line count", "The model's summary of the script"},
		Why: "Running without errors is not the same as doing the right thing. The tool's outputs show the effect."},
	{KP: 2, Q: "A script generated for one design hardcodes a die size and a list of macro names. What happens when it is reused?",
		A: "It fails or silently does the wrong thing on another design, because the values are design-specific", D: []string{"It adapts automatically to the new design", "The tool ignores hardcoded values", "Nothing, as long as the PDK is the same"},
		Why: "Generated scripts should read design data from the database or parameters, not embed one design's values."},
	{KP: 2, Q: "An assistant's script deletes and recreates the power grid to change one strap pitch. Why should the harness flag it?",
		A: "It changes far more state than requested, so its blast radius exceeds the intent", D: []string{"Power grids cannot be scripted", "Strap pitch is fixed by the PDK", "It will always crash the tool"},
		Why: "A good change is the smallest one that achieves the intent. A broad rewrite needs review even if it works."},
	{KP: 3, Q: "Why should LLM-generated scripts run in a sandbox with no network access and limited file-system scope?",
		A: "A script can run arbitrary shell commands and could read or send design data or damage files", D: []string{"Sandboxes make tools run faster", "EDA tools need network access to run", "Sandboxing improves QoR"},
		Why: "Generated code is untrusted. The IP boundary and the workspace must be protected by construction."},
	{KP: 3, Q: "A flow-script assistant is evaluated by whether its scripts run without errors on a test suite. What does the metric miss?",
		A: "Whether the scripts produce the requested results and acceptable QoR", D: []string{"Whether the scripts use Tcl", "Whether the scripts are short", "Nothing; running cleanly is the goal"},
		Why: "Execution success is a low bar. Evaluation should check outcomes against expected reports or QoR baselines."},
	{KP: 3, Q: "A generated script sets a random seed differently each run, and QoR varies between runs. Why does this matter for an agent comparing two script variants?",
		A: "Run-to-run noise can be mistaken for an improvement unless seeds are fixed or several runs are compared", D: []string{"It does not matter for scripts", "Seeds only affect simulation", "QoR noise always averages to zero in one run"},
		Why: "Comparisons need deterministic runs or statistics over runs, otherwise noise drives decisions."},
	{KP: 3, Q: "In the change-tier model, what kind of change is a generated script that only sets flow parameters?",
		A: "A Tier 1 configuration change, reversible and needing no recompilation", D: []string{"A Tier 3 core edit", "A Tier 2 plug-in", "A change to the PDK"},
		Why: "Parameters and command sequences sit in the configuration surface. They still need a run and a QoR comparison before adoption."},
}
