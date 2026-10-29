package agentic

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var callableItems = []edakit.Item{
	{KP: 1, Q: "An agent drives a legacy placer through its interactive Tcl session. The same script gives different results depending on which commands a previous user ran in the session. What does an agent-first engine change to fix this?",
		A: "Each operation becomes a function of explicit inputs: design handle, parameters and seed", D: []string{"The agent restarts the GUI before every command", "The agent reads the console log more carefully", "The placer stops accepting parameters"},
		Why: "Hidden session state makes a call's result depend on history. Explicit inputs make the same call give the same answer."},
	{KP: 1, Q: "In an agent-first flow, what happens to the tool's GUI?",
		A: "It becomes a viewer of the same design database the agents use", D: []string{"It is removed, since no human looks at designs any more", "It becomes the only way agents can run engines", "It stores the design state the engines read"},
		Why: "Humans still inspect results. The GUI shows the shared database instead of being the place where state lives."},
	{KP: 1, Q: "Which existing interface is the natural starting point for exposing OpenROAD engines to agents?",
		A: "Its Tcl and Python scripting interfaces over the design database", D: []string{"Its GUI menus driven by screen automation", "Its log files parsed after each run", "Its source code, which the agent edits for every run"},
		Why: "Scripting APIs already expose engine operations. The agent-first step is to make them explicit, stateless calls with structured results."},
	{KP: 1, Q: "Why can an agent not reliably reuse a flow written as a long interactive session with many set and get commands?",
		A: "The implicit state between commands makes each step's behavior depend on hidden history", D: []string{"Agents cannot execute Tcl", "Session flows run too fast for an agent", "Interactive sessions cannot load designs"},
		Why: "An agent reasons about the effect of each call. Implicit state hides the inputs that determine that effect."},
	{KP: 2, Q: "A timer returns its results only as a 3,000-line text report. What is the main problem for an agent calling it?",
		A: "The agent must scrape prose, which is slow and breaks when the report format changes", D: []string{"The report is too accurate", "Text reports cannot contain slack values", "The timer must then run twice"},
		Why: "Agents need named, structured values. Scraping formatted text is brittle and costly."},
	{KP: 2, Q: "Which return value lets an agent decide its next step after a global placement call?",
		A: "Named metrics (HPWL, overflow, density, runtime), a status and the handle of the produced state", D: []string{"A screenshot of the placement", "The full DEF file inlined in the response", "A success message with no numbers"},
		Why: "Metrics, status and a handle let the agent compare, branch and pass the result to the next engine."},
	{KP: 2, Q: "A routing call fails. Which failure report is most useful to an agent?",
		A: "A structured record naming the failed check, the objects involved and their values", D: []string{"The exit code 1 and nothing else", "The last 500 lines of the router log", "A message saying routing was unsuccessful"},
		Why: "A located, structured failure lets the agent choose a fix. Raw logs and bare codes do not."},
	{KP: 2, Q: "Why should every engine result carry provenance (engine version, parameters, seed)?",
		A: "So any result can be reproduced and compared with others on equal terms", D: []string{"So the result file is larger", "Because the model needs it to compute slack", "So that the GUI can display it"},
		Why: "Provenance makes results reproducible and comparable, which differential tests and agents both rely on."},
	{KP: 3, Q: "An agent is asked for the setup slack of a path. What should it do?",
		A: "Call the timer and use the slack it returns", D: []string{"Compute the slack itself by reasoning over the netlist", "Estimate it from earlier runs", "Ask the model for its best guess"},
		Why: "Every number the agent relies on comes from an engine. The model orchestrates; it does not replace exact computation."},
	{KP: 3, Q: "Making engines callable by agents increases how often they are invoked. What does this imply for engine performance?",
		A: "Engine speed matters as much as ever, because an agent loop may call the engine hundreds of times", D: []string{"Engine speed stops mattering, because the model does the work", "Engines can be slower, since agents work overnight", "Only the model's speed matters"},
		Why: "The compute stays in the engines. A faster engine means more agent iterations per day."},
	{KP: 3, Q: "Which split of work between model and engine is correct in an agent-first flow?",
		A: "The model chooses engines, skills and parameters; compiled engines do the exact computation", D: []string{"The model legalizes placements; engines only report", "Engines choose skills; the model propagates timing", "The model and engines each compute slack and vote"},
		Why: "Cheap decisions belong to the model and expensive exact work to the engines."},
	{KP: 3, Q: "An agent proposes to legalize a placement by describing cell moves in text and applying them one by one. Why is this a poor design?",
		A: "Legalization is exact, expensive work the compiled legalizer does correctly and far faster", D: []string{"Text cannot describe cell moves", "Legalization is never needed after global placement", "The agent cannot read cell coordinates"},
		Why: "Reimplementing engine algorithms in the model is slow and error-prone; the agent should call the legalizer."},
}

var fourLayerItems = []edakit.Item{
	{KP: 1, Q: "Which layer of the four-layer engine model holds a router's DRC-correct path model?",
		A: "The immutable core", D: []string{"Strategy or policy", "Parameters", "Observation hooks"},
		Why: "Legality and correctness are core. Changing them needs human approval and full regression."},
	{KP: 1, Q: "Which layer holds the cost function a router uses to rank candidate paths?",
		A: "Strategy or policy", D: []string{"The immutable core", "Parameters", "Observation hooks"},
		Why: "Cost functions decide how the engine searches, which is policy; new ones are added as plug-ins."},
	{KP: 1, Q: "Which layer holds the density target of a global placer?",
		A: "Parameters", D: []string{"The immutable core", "Strategy or policy", "Observation hooks"},
		Why: "A target value is configuration and needs no recompile."},
	{KP: 1, Q: "Which layer is where an agent reads overflow during global placement and can stop the run early?",
		A: "Observation hooks", D: []string{"The immutable core", "Strategy or policy", "Parameters"},
		Why: "Hooks are where an agent reads live state and injects behavior at defined points."},
	{KP: 2, Q: "A request says: weight timing twice as heavily as wirelength in placement. Which layer does it change?",
		A: "Parameters", D: []string{"The immutable core", "Strategy or policy", "Observation hooks"},
		Why: "It changes a weight in an existing objective, so it is a configuration change."},
	{KP: 2, Q: "A request says: let two cells overlap by half a site so the router can sort it out. Which layer would this touch?",
		A: "The immutable core, and it should be refused", D: []string{"Parameters, as a tolerance", "Strategy or policy, as a new move set", "Observation hooks, as an injected callback"},
		Why: "Overlap-free legalization is a correctness guarantee. Weakening it breaks the core contract."},
	{KP: 2, Q: "A request says: try a swap-and-shift move set in detailed placement instead of the current one. Which layer does it change?",
		A: "Strategy or policy", D: []string{"The immutable core", "Parameters", "Observation hooks"},
		Why: "Move sets are heuristics that decide how the engine searches; they are selected or added as policies."},
	{KP: 2, Q: "A request says: report the worst negative slack every 20 iterations of timing-driven placement. Which layer serves it?",
		A: "Observation hooks", D: []string{"The immutable core", "Strategy or policy", "Parameters"},
		Why: "Reading live state at defined points is what hooks are for."},
	{KP: 2, Q: "A request says: change how the extractor's field solver discretizes conductor surfaces. Which layer is that?",
		A: "The immutable core", D: []string{"Parameters", "Strategy or policy", "Observation hooks"},
		Why: "The field-solver mathematics is the numerical kernel. Choosing an existing accuracy mode would be a parameter; changing the kernel is core."},
	{KP: 3, Q: "An engine's change log shows that most routine tuning requests needed core edits. What does the four-layer rule say?",
		A: "The line between core and change surface is in the wrong place; variation must move out of the core", D: []string{"Core edits are the normal way to tune", "The engine needs fewer parameters", "Requests should be refused"},
		Why: "The design rule is a small, rarely touched core, with everything that varies in policy, parameters and hooks."},
	{KP: 3, Q: "What is the test for whether behavior belongs in the core?",
		A: "Whether changing it could make a result illegal or wrong, rather than better or worse", D: []string{"Whether it is written in C++", "Whether it runs in the hot loop", "Whether a customer asked for it"},
		Why: "The core guarantees correctness; quality trade-offs live outside it."},
	{KP: 3, Q: "Why is a core that is too small also a problem?",
		A: "It exposes legality to plug-ins that can break it", D: []string{"It makes the engine slower", "It forces every request to be Tier 3", "It prevents configuration changes"},
		Why: "Correctness guarantees must stay in the core so no plug-in or parameter can violate them."},
	{KP: 3, Q: "What does the four-layer design rule aim for in the common case?",
		A: "The common change never recompiles the core", D: []string{"Every change recompiles the whole engine", "All changes go through observation hooks", "Policies are never added after release"},
		Why: "Pushing variation into policy, parameters and hooks keeps routine changes fast and safe."},
}

var coreSurfaceItems = []edakit.Item{
	{KP: 1, Q: "For a placer, which item belongs to the core rather than the change surface?",
		A: "Overlap-free assignment of cells to legal rows and sites", D: []string{"The weight of congestion in the objective", "The density target", "A rule that abuts regular datapath cells"},
		Why: "Legalization guarantees a correct placement; the others trade quality."},
	{KP: 1, Q: "For a router, which item belongs to the change surface?",
		A: "The rip-up and reroute strategy", D: []string{"Spacing and width rule checks", "Grid and layer legality", "Via enclosure correctness"},
		Why: "Rip-up strategy decides how the router searches; legality rules are core."},
	{KP: 1, Q: "A customer wants the placer to keep the cells of each bit slice abutted. Where does this land?",
		A: "On the placer's change surface, as a rule at a placement extension point", D: []string{"In legalization, by changing row rules", "In the router's cost function", "Nowhere: it requires a new placer"},
		Why: "Abutment is a preference about where legal cells go, not a change to what legal means."},
	{KP: 1, Q: "Two routers use different penalty functions and both produce DRC-clean results. What does this illustrate?",
		A: "Penalty functions are change surface: different choices can all be correct", D: []string{"One router must be wrong", "Penalty functions are part of legality", "DRC rules are parameters"},
		Why: "Quality trade-offs vary; correctness does not."},
	{KP: 2, Q: "For a timer, which item belongs to the core?",
		A: "Arrival, required and slack propagation through the timing graph", D: []string{"Which corners to report", "How many paths per endpoint to report", "Which derating scheme to use"},
		Why: "Propagation and the delay-model API are the timer's core."},
	{KP: 2, Q: "For an extractor, which item belongs to the change surface?",
		A: "Using a pattern library on non-critical nets and the field solver on critical ones", D: []string{"The field-solver equations", "Making sure every conductor is accounted for", "Physical consistency of the capacitance matrix"},
		Why: "Accuracy-versus-speed choices per net are tunable; the solver mathematics and netlist correctness are core."},
	{KP: 2, Q: "A team wants the timer to run a fast surrogate mode during exploration. Where does this land?",
		A: "On the timer's change surface, as a mode", D: []string{"In graph propagation", "In the delay-model API", "In the SDC parser"},
		Why: "Fast or surrogate modes are listed on the timer's change surface."},
	{KP: 2, Q: "Calibrating an extractor's pattern library against a reference field solver is which kind of change?",
		A: "A change-surface change: calibration data for an existing mode", D: []string{"A core change to the solver", "A change to netlist correctness", "A change to the timer"},
		Why: "Calibration tunes how the fast mode approximates the solver; it does not change the solver."},
	{KP: 3, Q: "A request says: route clock nets at double width. What is it, really?",
		A: "A layer directive or non-default rule on the router's change surface", D: []string{"A change to the router's legality core", "A change to the timer", "Impossible without rewriting the router"},
		Why: "Wider clock routing is expressed with existing rules; legality is unchanged."},
	{KP: 3, Q: "A request says: skip coupling capacitances below 0.01 fF to speed up extraction. How should it be handled?",
		A: "As an explicit, reviewed accuracy mode, because it changes what the extractor guarantees", D: []string{"As a silent threshold set by the agent", "As a router change", "It needs no review because the values are small"},
		Why: "Dropping capacitance changes the accuracy contract, so it must be a declared mode, not a hidden shortcut."},
	{KP: 3, Q: "A request says: ignore timing on test-mode paths. What is it, really?",
		A: "A constraint or reporting choice, not a change to propagation", D: []string{"A core change to propagation", "A change to the delay-model API", "A change to the Liberty library"},
		Why: "Excluding paths is done with constraints such as false paths or mode settings; propagation is untouched."},
	{KP: 3, Q: "A request says: let the placer leave a few overlaps for the router to fix. What is the right response?",
		A: "Refuse it: it breaks the placer's legality contract", D: []string{"Accept it as a density parameter", "Accept it as a move-set policy", "Accept it if wirelength improves"},
		Why: "Overlap-free placement is the core guarantee. No quality gain justifies weakening it."},
}

var hookItems = []edakit.Item{
	{KP: 1, Q: "Why must an observation hook take its snapshot at a defined point between phases?",
		A: "So the state it reports is consistent, with data structures that agree with each other", D: []string{"So the snapshot is smaller", "Because hooks cannot run during placement", "So the agent can change legality"},
		Why: "A snapshot taken mid-update can mix old and new values and mislead the agent."},
	{KP: 1, Q: "What should hooks return to an agent?",
		A: "Numbers and identifiers, such as iteration, HPWL, overflow and the densest bins", D: []string{"A rendered image of the layout", "Free-form prose describing progress", "The engine's internal pointers"},
		Why: "Agents act on structured values, not prose or internal addresses."},
	{KP: 1, Q: "What is the default access mode of an observation hook, and why?",
		A: "Read-only, so an observer cannot corrupt the run", D: []string{"Read-write, so agents can fix problems directly", "Write-only, so agents can inject behavior", "No access until the run ends"},
		Why: "Most hooks only observe. Injection is a separate, restricted capability."},
	{KP: 1, Q: "A placer hook shows overflow has stayed at 18% for 200 iterations. What can the agent do with this?",
		A: "Stop the run or switch strategy rather than waste iterations", D: []string{"Edit the legalizer", "Ignore it, since hooks are for logging only", "Restart with the same settings"},
		Why: "Hooks let the agent see convergence and act on it."},
	{KP: 1, Q: "Which is a sensible observation hook for a router?",
		A: "Congestion and DRC counts after each rip-up and reroute pass", D: []string{"The cell coordinates before global placement", "The SDC file contents", "The Liberty tables"},
		Why: "Hooks report the engine's own live state at natural boundaries."},
	{KP: 3, Q: "An injected callback changes net weights between global placement iterations. Which tier is that change?",
		A: "Tier 1 if it only sets existing weights; injected code would be Tier 2", D: []string{"Always Tier 3", "Not a change at all", "Tier 3 unless the agent is confident"},
		Why: "Injection through hooks still follows the change tiers."},
	{KP: 3, Q: "Where may behavior be injected through a hook?",
		A: "Only at defined points between phases, never inside the hot kernel", D: []string{"Anywhere the agent chooses", "Inside the legalization loop", "Inside the timing propagation loop"},
		Why: "Defined injection points keep the kernel fast and the state consistent."},
	{KP: 3, Q: "What must the engine do after an injected change?",
		A: "Re-check its invariants, so a bad callback produces an error rather than an illegal result", D: []string{"Nothing: injected changes are trusted", "Restart the whole flow", "Disable all hooks"},
		Why: "Invariant checks keep the core's guarantees intact whatever is injected."},
	{KP: 3, Q: "Why is every injection logged with its values?",
		A: "So the run stays reproducible", D: []string{"To slow the engine down", "So the model can see its own output", "Because hooks cannot return results"},
		Why: "A run with injected behavior can be replayed only if the injections are recorded."},
	{KP: 3, Q: "An agent wants to inject a callback that moves cells to new rows directly. Why must the engine refuse?",
		A: "Moving cells outside legalization could break legality, which no hook may alter", D: []string{"Callbacks cannot move cells for performance reasons", "Rows are a parameter", "Only humans may move cells"},
		Why: "Hooks act only on the change surface; legality stays in the core."},
}

func hookOverhead(ctx generator.GeneratorContext) generator.Problem {
	engines := []struct{ name, iter, snap string }{
		{"global placer", "placement iteration", "density snapshot"},
		{"router", "rip-up and reroute pass", "congestion map"},
		{"timing-driven placer", "iteration", "timing summary"},
	}
	e := edakit.Pick(engines)
	t := float64(5 * edakit.Int(4, 16))
	c := float64(10 * edakit.Int(4, 30))
	if edakit.Int(0, 1) == 0 {
		ks := []int{5, 10, 20, 25, 50}
		k := float64(edakit.Pick(ks))
		ans := 100 * c / (k * t)
		q := edakit.Sprintf("A %s spends %s ms per %s. An agent's observation hook takes a %s costing %s ms every %s iterations. By what percentage does the hook increase runtime? Round to 1 decimal place.",
			e.name, edakit.Num(t, 0), e.iter, e.snap, edakit.Num(c, 0), edakit.Num(k, 0))
		exp := edakit.Sprintf("The overhead fraction is $c/(Kt) = %s / (%s \\times %s) = %s\\%%$.", edakit.Num(c, 0), edakit.Num(k, 0), edakit.Num(t, 0), edakit.Num(ans, 1))
		return choice(q, ans, 1, "%", exp, 2.5, 100*c/t, 100*c/(k*t+c), 100*c/(k*(t+c))+1)
	}
	x := float64(edakit.Pick([]int{2, 5, 10}))
	kmin := c / (x / 100 * t)
	ans := float64(int(kmin))
	if ans < kmin-1e-9 {
		ans++
	}
	q := edakit.Sprintf("A %s spends %s ms per %s, and a %s costs %s ms. What is the smallest snapshot interval $K$, in iterations, that keeps the hook's overhead at or below %s%% of runtime?",
		e.name, edakit.Num(t, 0), e.iter, e.snap, edakit.Num(c, 0), edakit.Num(x, 0))
	exp := edakit.Sprintf("Require $c/(Kt) \\le %s\\%%$, so $K \\ge %s / (%s \\times %s) = %s$; round up to %s.", edakit.Num(x, 0), edakit.Num(c, 0), edakit.Num(x/100, 2), edakit.Num(t, 0), edakit.Num(kmin, 2), edakit.Num(ans, 0))
	return choice(q, ans, 0, " iterations", exp, 3, ans-1, c/t, kmin*2, x*t/c*10)
}

var manifestItems = []edakit.Item{
	{KP: 1, Q: "Which four things make up a skill folder?",
		A: "Source, a build, golden tests and a plain-language manifest", D: []string{"Source, a GUI, a log and a license", "A binary, a design, a report and a seed", "A model checkpoint, prompts, a reward and a dataset"},
		Why: "A skill is ordinary engine code packaged with its build, tests and a description."},
	{KP: 1, Q: "Why is a skill built as its own compile unit, such as a shared library at an extension point?",
		A: "So changing it rebuilds only that unit, in seconds", D: []string{"So it can bypass the core's legality checks", "So it can be written in any language without an interface", "So it does not need tests"},
		Why: "Isolated builds keep the harness loop fast."},
	{KP: 1, Q: "What do the golden tests in a skill folder consist of?",
		A: "Known-good inputs and outputs plus a QoR baseline for differential comparison", D: []string{"Random inputs with no expected outputs", "A list of customers who use the skill", "The model prompts used to write the skill"},
		Why: "Goldens and a baseline let the harness prove a change correct."},
	{KP: 1, Q: "An engine developer asks what is new about writing a skill. What is the honest answer?",
		A: "Very little: it is the same code, with a description and tests beside it", D: []string{"Skills must be written by the model", "Skills replace the engine's algorithms", "Skills are trained, not written"},
		Why: "The packaging is ordinary; the algorithm work is unchanged."},
	{KP: 2, Q: "A manifest lists a parameter max_displacement but no type or range. What risk does this create?",
		A: "The agent may set values the skill was never designed or tested for", D: []string{"The skill will not compile", "The golden tests will be skipped", "The agent cannot read the manifest"},
		Why: "Declared types and ranges let both the agent and the core validate parameters."},
	{KP: 2, Q: "Which manifest line best states applicability and limits?",
		A: "Assumes a single row height; not tested with multi-height cells", D: []string{"Works well on most designs", "Written in modern C++", "Fast and accurate"},
		Why: "Concrete limits let an agent rule the skill out where it does not apply."},
	{KP: 2, Q: "Which item belongs in a manifest because an agent needs it to choose between skills?",
		A: "Typical runtime per million cells", D: []string{"The author's coding style", "The full list of internal functions", "The build directory layout"},
		Why: "Cost is part of the choice, especially inside a fast loop."},
	{KP: 2, Q: "What does a manifest's tunable surface section list?",
		A: "Named parameters with types and ranges and the extension points the skill offers", D: []string{"Every source file in the skill", "The customer designs it was used on", "The model that wrote it"},
		Why: "The tunable surface tells the agent what may be changed without touching the core."},
	{KP: 3, Q: "A code change removed a skill's support for macros, but the manifest still advertises it. What kind of failure is this?",
		A: "A stale manifest", D: []string{"An overclaiming manifest", "A missing-limits manifest", "A build failure"},
		Why: "A manifest that describes removed behavior is stale; selection will pick the skill for designs it can no longer handle."},
	{KP: 3, Q: "A manifest says 'works for all designs', but its golden tests cover only one datapath style. What kind of failure is this?",
		A: "An overclaiming manifest", D: []string{"A stale manifest", "A frozen ABI", "A determinism bug"},
		Why: "The claim exceeds the evidence; every claimed capability should be exercised by a golden test."},
	{KP: 3, Q: "How can a team keep manifests honest over time?",
		A: "Tie every claimed capability to a golden test and version the manifest with the code", D: []string{"Let the model rewrite manifests after each run", "Keep manifests short by omitting limits", "Store manifests separately from the source"},
		Why: "Claims backed by tests and versioned with the code cannot silently drift."},
	{KP: 3, Q: "Why is a wrong manifest as harmful as a wrong algorithm in an agent-first system?",
		A: "Selection is made from manifests alone, so a wrong one sends the skill to designs it cannot handle", D: []string{"Manifests are compiled into the engine", "The manifest computes the QoR", "Agents edit the algorithm through the manifest"},
		Why: "The agent never reads the source to select; the manifest is its only evidence."},
}

var selectionItems = []edakit.Item{
	{KP: 1, Q: "A design has a 64-bit ALU built from repeated bit slices. Which manifest claim matches it best?",
		A: "Placement for regular bit-sliced datapath logic", D: []string{"Placement for designs dominated by large macros", "Routing for analog-sensitive nets", "Extraction for long global wires"},
		Why: "Selection matches design traits against what each manifest advertises."},
	{KP: 1, Q: "A skill's manifest says it requires single-height rows. The design has multi-height cells. What should selection do?",
		A: "Exclude the skill, whatever its other strengths", D: []string{"Choose it if its score is highest", "Choose it and let legalization fix problems", "Edit the design to remove multi-height cells"},
		Why: "Hard requirements filter candidates before any scoring."},
	{KP: 1, Q: "What does the agent read to choose a skill?",
		A: "The manifests, matched against a summary of the design and the current objective", D: []string{"The source of every skill", "The compiled binaries", "The full design netlist"},
		Why: "Selection is a text decision over catalog entries; the source is never opened to select."},
	{KP: 1, Q: "Which design summary is most useful for skill selection?",
		A: "Utilization, regularity, macro count, clock and corner count, current congestion and timing, and the step's objective", D: []string{"The full gate-level netlist", "The names of the designers", "The RTL comments"},
		Why: "The summary uses the terms manifests advertise."},
	{KP: 3, Q: "Every eligible placement skill scores poorly on a congestion-dominated design. What should the agent do next?",
		A: "Customize: tune a parameter or write a plug-in at a skill's extension point", D: []string{"Pick the best-scoring ineligible skill", "Give up on the design", "Edit the placer's legalizer"},
		Why: "When no skill fits, selection turns into customization."},
	{KP: 3, Q: "The fast loop shows the chosen router skill failing on the same metric run after run. What does this signal?",
		A: "The skill needs customizing, or another skill should be chosen", D: []string{"The fast loop is broken", "The manifest must be deleted", "Nothing: one metric does not matter"},
		Why: "The decision to customize comes from measured metrics in the live loop."},
	{KP: 3, Q: "After many runs, a skill's measured strength on timing differs from what its manifest claims. Which should selection trust?",
		A: "The measured strength, because the manifest states intent and the loop measures fact", D: []string{"The manifest, because it was written by the engine owner", "Neither: drop the skill", "Whichever is higher"},
		Why: "Evidence accumulated in the loop should replace advertised values."},
	{KP: 3, Q: "Where does the decision to customize a skill come from in this architecture?",
		A: "The metrics and variance of the live fast-skill loop", D: []string{"A static rule written at design time", "A human reading every report", "The model's general knowledge"},
		Why: "Selection and the decision to customize are both driven by the fast loop's data."},
	{KP: 2, Q: "Two eligible skills have equal total scores, but one is much faster. What does a weighted scheme with a runtime objective do?",
		A: "It already favors the faster one through the runtime weight", D: []string{"It ignores runtime", "It picks at random", "It rejects both"},
		Why: "Runtime is one objective in the weighted sum."},
	{KP: 2, Q: "Why are hard requirements applied before scoring rather than as a large negative score?",
		A: "So no combination of strengths can make an inapplicable skill win", D: []string{"Because scores cannot be negative", "To make selection slower", "Because requirements are optional"},
		Why: "A filter guarantees exclusion; a penalty can be outweighed."},
}

type skillCand struct {
	name    string
	s       [3]int
	limited bool
	score   float64
}

func selectionScore(ctx generator.GeneratorContext) generator.Problem {
	facts := []struct{ design, limit string }{
		{"the design has multi-height cells", "requires single-height rows"},
		{"the design has 6 clock domains", "supports one clock domain only"},
		{"the design has 40 hard macros", "assumes no macros"},
		{"timing must be analyzed at 5 corners", "evaluates a single corner"},
	}
	ws := [][3]float64{{0.5, 0.3, 0.2}, {0.2, 0.5, 0.3}, {0.4, 0.4, 0.2}, {0.6, 0.2, 0.2}, {0.3, 0.3, 0.4}, {0.2, 0.6, 0.2}}
	objs := [3]string{"timing", "congestion", "runtime"}
	f := edakit.Pick(facts)
	w := edakit.Pick(ws)
	var cs []skillCand
	for tries := 0; tries < generator.MaxDrawAttempts; tries++ {
		cs = nil
		for i, n := range []string{"Skill A", "Skill B", "Skill C", "Skill D"} {
			c := skillCand{name: n, s: [3]int{edakit.Int(2, 9), edakit.Int(2, 9), edakit.Int(2, 9)}}
			c.limited = i == 0 || (i == 1 && edakit.Int(0, 1) == 0)
			c.score = w[0]*float64(c.s[0]) + w[1]*float64(c.s[1]) + w[2]*float64(c.s[2])
			cs = append(cs, c)
		}
		cs = edakit.Shuffle(cs)
		sorted := append([]skillCand(nil), cs...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].score > sorted[j].score })
		var elig []skillCand
		for _, c := range sorted {
			if !c.limited {
				elig = append(elig, c)
			}
		}
		if sorted[0].limited && len(elig) >= 2 && elig[0].score-elig[1].score > 0.15 {
			break
		}
	}
	sort.Slice(cs, func(i, j int) bool { return cs[i].name < cs[j].name })
	var rows, scores []string
	best := ""
	bestScore := -1.0
	for _, c := range cs {
		lim := ""
		if c.limited {
			lim = ", manifest: " + f.limit
		}
		rows = append(rows, edakit.Sprintf("%s (timing %d, congestion %d, runtime %d%s)", c.name, c.s[0], c.s[1], c.s[2], lim))
		scores = append(scores, edakit.Sprintf("%s %s%s", c.name, edakit.Num(c.score, 2), map[bool]string{true: " (excluded)", false: ""}[c.limited]))
		if !c.limited && c.score > bestScore {
			best, bestScore = c.name, c.score
		}
	}
	q := edakit.Sprintf("An agent selects a skill by filtering on hard requirements and then scoring eligible skills as $\\sum_i w_i s_i$. The objective weights are %s %s, %s %s and %s %s, and %s. The manifests advertise strengths on a 1 to 9 scale: %s. Which skill does the agent select?",
		objs[0], edakit.Num(w[0], 1), objs[1], edakit.Num(w[1], 1), objs[2], edakit.Num(w[2], 1), f.design, strings.Join(rows, "; "))
	exp := edakit.Sprintf("Weighted scores: %s. The highest-scoring skill that passes the hard requirement is %s.", strings.Join(scores, ", "), best)
	var d []string
	for _, c := range cs {
		if c.name != best {
			d = append(d, c.name)
		}
	}
	return edakit.Choice(q, best, exp, d...)
}

func selectionCost(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		s := edakit.Scaled(ctx, 12, 60)
		m := 50 * edakit.Int(4, 12)
		d := 500 * edakit.Int(2, 6)
		src := 10000 * edakit.Int(3, 12)
		man := s*m + d
		full := s*src + d
		ratio := float64(full) / float64(man)
		q := edakit.Sprintf("A skill library has %d skills. Each manifest is about %d tokens, the design summary is %d tokens, and reading one skill's source would cost about %d tokens. How many times more tokens does selecting by reading every skill's source take than selecting by manifests (both include the design summary)? Round to 1 decimal place.",
			s, m, d, src)
		exp := edakit.Sprintf("Manifests: $%d \\times %d + %d = %d$ tokens. Source: $%d \\times %d + %d = %d$ tokens. Ratio %s.", s, m, d, man, s, src, d, full, edakit.Num(ratio, 1))
		return edakit.Numeric(q, edakit.Num(ratio, 1), exp, map[string]string{"manifest_tokens": edakit.IntStr(man), "source_tokens": edakit.IntStr(full)})
	case 2:
		in := 1000 * edakit.Int(6, 30)
		pre := 1000 * edakit.Int(2, 8)
		out := 50 * edakit.Int(2, 10)
		dec := 10 * edakit.Int(3, 10)
		loop := edakit.Int(2, 10)
		t := float64(in)/float64(pre) + float64(out)/float64(dec)
		pct := 100 * t / float64(60*loop)
		q := edakit.Sprintf("A selection call reads %d input tokens at a prefill rate of %d tokens per second and writes %d output tokens at %d tokens per second. One iteration of the fast place-route-analyze loop takes %d minutes. What percentage of one loop iteration does the selection call take? Round to 1 decimal place.",
			in, pre, out, dec, loop)
		exp := edakit.Sprintf("Call time $%d/%d + %d/%d = %s$ s; as a share of %d s it is $%s\\%%$.", in, pre, out, dec, edakit.Num(t, 2), 60*loop, edakit.Num(pct, 1))
		return edakit.Numeric(q, edakit.Num(pct, 1), exp, map[string]string{"call_seconds": edakit.Num(t, 2)})
	}
	w := 1000 * edakit.Pick([]int{16, 32, 64, 128})
	r := 1000 * edakit.Int(3, 12)
	m := 50 * edakit.Int(5, 16)
	n := (w - r) / m
	q := edakit.Sprintf("A model has a %d-token context window. The harness reserves %d tokens for instructions, the design summary and the answer. Manifests average %d tokens. How many manifests fit in one selection call?", w, r, m)
	exp := edakit.Sprintf("$\\lfloor (%d - %d) / %d \\rfloor = %d$ manifests. A larger library needs a tag-based prefilter.", w, r, m, n)
	return edakit.Numeric(q, edakit.IntStr(n), exp, map[string]string{"fit": edakit.IntStr(n)})
}
