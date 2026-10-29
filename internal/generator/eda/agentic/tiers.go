package agentic

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var tierItems = []edakit.Item{
	{KP: 1, Q: "A team asks the placer to target 70% density instead of 65%. Which tier is this?",
		A: "Tier 1, configuration", D: []string{"Tier 2, a plug-in", "Tier 3, a core edit", "It needs a new placer"},
		Why: "A threshold within the declared range is a parameter change with no recompile."},
	{KP: 1, Q: "A team asks the router to prefer metal 5 for a 128-bit bus. Which tier is this?",
		A: "Tier 1, a layer directive", D: []string{"Tier 2, a new cost function", "Tier 3, a change to layer legality", "A placer change"},
		Why: "Layer preferences are existing directives on the router's change surface."},
	{KP: 1, Q: "A team wants to switch from the flat OCV provider to the shipped AOCV provider with its own tables. Which tier is this?",
		A: "Tier 1: choose the provider and point it at the tables", D: []string{"Tier 2: write a new provider", "Tier 3: change propagation", "Tier 3: change the Liberty reader"},
		Why: "Selecting an existing skill and its data is configuration."},
	{KP: 1, Q: "An extraction run should use the field solver for clock nets and pattern tables elsewhere, both of which the extractor already supports per net class. Which tier is this?",
		A: "Tier 1, a per-net-class mode setting", D: []string{"Tier 2, a new solver", "Tier 3, a change to the solver mathematics", "A timer change"},
		Why: "Choosing among existing accuracy modes is a parameter."},
	{KP: 2, Q: "A customer wants vias penalized within 20 um of an analog block, and no existing parameter can express distance to a block. Which tier is this?",
		A: "Tier 2, a cost-function plug-in at the router's cost seam", D: []string{"Tier 1, a via weight", "Tier 3, a change to DRC legality", "A placement density change"},
		Why: "New behavior at an existing extension point is a plug-in; legality is untouched."},
	{KP: 2, Q: "Silicon correlation shows HVT data paths of logic depth under 5 need 3% more late derate, and no table expresses that condition. Which tier is this?",
		A: "Tier 2, a new derate provider", D: []string{"Tier 1, a global derate", "Tier 3, a propagation change", "Tier 1, a new sigma count"},
		Why: "The rule needs code, but only at the derate seam."},
	{KP: 2, Q: "A design needs detailed placement to keep each bit slice's cells abutted, a rule no shipped move set implements. Which tier is this?",
		A: "Tier 2, a move-evaluation plug-in", D: []string{"Tier 1, a density target", "Tier 3, a legalization change", "A routing directive"},
		Why: "A new placement preference is policy added at an extension point."},
	{KP: 2, Q: "A synthesis team wants a new rewriting script ordering for its arithmetic blocks, built from existing passes. Which tier is this?",
		A: "Tier 1 if it is a script of existing passes; Tier 2 only if a new pass is needed", D: []string{"Always Tier 3", "Tier 2, because scripts are code", "It cannot be expressed"},
		Why: "A recipe of existing passes is configuration; a new transformation is a plug-in."},
	{KP: 3, Q: "A request asks to change how arrival times combine at a multi-input gate. Which tier is this?",
		A: "Tier 3, a core edit needing a human-approved diff and full regression", D: []string{"Tier 1, a parameter", "Tier 2, a derate provider", "Tier 2, a reporting plug-in"},
		Why: "Propagation is the timer's core and affects every result."},
	{KP: 3, Q: "A customer needs derates keyed on an instance attribute the derate seam does not pass. What is the right path?",
		A: "One reviewed core edit to extend the seam's context, then Tier 2 providers from then on", D: []string{"A Tier 2 provider that reads the attribute from global state", "A Tier 1 parameter", "Refuse the request permanently"},
		Why: "Extending the context is a legitimate Tier 3 change that turns future requests into plug-ins."},
	{KP: 3, Q: "A request asks the router to skip spacing checks on a congested region to finish faster. What should happen?",
		A: "Refuse it: it would make results illegal", D: []string{"Accept it as Tier 1", "Accept it as a Tier 2 cost plug-in", "Accept it if the human is busy"},
		Why: "Weakening legality is not a customization; it breaks the core contract."},
	{KP: 3, Q: "Most derating requests at one site keep needing Tier 3 edits. What does that indicate?",
		A: "The seam is in the wrong place or exposes too little context", D: []string{"Derating is inherently a core feature", "The customer is misusing the tool", "Tier 3 is the normal tier for derating"},
		Why: "A well-placed seam makes core edits for that kind of change essentially never necessary."},
	{KP: 3, Q: "A request sounds like a core change: route clocks differently. What should the model do first?",
		A: "Check whether existing parameters or directives already express it", D: []string{"Open a Tier 3 proposal immediately", "Write a new router", "Refuse it"},
		Why: "Always try the smallest safe tier first; many core-sounding requests are parameters."},
}

func tierTurnaround(ctx generator.GeneratorContext) generator.Problem {
	f3 := float64(edakit.Int(2, 10)) / 100
	f2 := float64(edakit.Int(15, 35)) / 100
	f1 := 1 - f2 - f3
	t1 := float64(edakit.Int(2, 10)) / 10
	t2 := float64(edakit.Int(2, 10))
	t3 := float64(10 * edakit.Int(3, 10))
	e := f1*t1 + f2*t2 + f3*t3
	head := edakit.Sprintf("Of the change requests at a site, %s%% land in Tier 1 taking %s h each, %s%% in Tier 2 taking %s h and %s%% in Tier 3 taking %s h. ",
		edakit.Num(100*f1, 0), edakit.Num(t1, 1), edakit.Num(100*f2, 0), edakit.Num(t2, 0), edakit.Num(100*f3, 0), edakit.Num(t3, 0))
	switch kp(ctx) {
	case 1:
		q := head + "What is the expected turnaround per request in hours? Round to 2 decimal places."
		exp := edakit.Sprintf("$%s \\times %s + %s \\times %s + %s \\times %s = %s$ h.", edakit.Num(f1, 2), edakit.Num(t1, 1), edakit.Num(f2, 2), edakit.Num(t2, 0), edakit.Num(f3, 2), edakit.Num(t3, 0), edakit.Num(e, 2))
		return edakit.Numeric(q, edakit.Num(e, 2), exp, map[string]string{"expected": edakit.Num(e, 2)})
	case 2:
		share := 100 * f3 * t3 / e
		q := head + "What percentage of all change time is spent on Tier 3 requests? Round to 1 decimal place."
		exp := edakit.Sprintf("Tier 3 contributes $%s \\times %s = %s$ h of the expected %s h, which is %s%%.", edakit.Num(f3, 2), edakit.Num(t3, 0), edakit.Num(f3*t3, 2), edakit.Num(e, 2), edakit.Num(share, 1))
		return edakit.Numeric(q, edakit.Num(share, 1), exp, nil)
	}
	g := float64(edakit.Pick([]int{50, 60, 75, 80, 90})) / 100
	moved := f3 * g
	e2 := f1*t1 + (f2+moved)*t2 + (f3-moved)*t3
	q := head + edakit.Sprintf("Extending one seam's context converts %s%% of the Tier 3 requests into Tier 2 requests. What is the new expected turnaround in hours? Round to 2 decimal places.", edakit.Num(100*g, 0))
	exp := edakit.Sprintf("A fraction %s of all requests moves to Tier 2. New mean $%s \\times %s + %s \\times %s + %s \\times %s = %s$ h, down from %s h.",
		edakit.Num(moved, 4), edakit.Num(f1, 2), edakit.Num(t1, 1), edakit.Num(f2+moved, 4), edakit.Num(t2, 0), edakit.Num(f3-moved, 4), edakit.Num(t3, 0), edakit.Num(e2, 2), edakit.Num(e, 2))
	return edakit.Numeric(q, edakit.Num(e2, 2), exp, nil)
}

var extensionItems = []edakit.Item{
	{KP: 1, Q: "What is the standard construction of an extension point?",
		A: "An abstract base declaring the few operations the core needs, plus a factory that creates implementations by name", D: []string{"A global function pointer any plug-in may overwrite", "A copy of the core for each customer", "A script that patches the core at load time"},
		Why: "The core calls the base; the factory maps a configured name to an implementation."},
	{KP: 1, Q: "A router's cost seam has implementations named negotiated and analog_guard. How does a run choose between them?",
		A: "By configuration: the factory creates the named implementation", D: []string{"By recompiling the core", "By editing the base class", "By the order the plug-ins were installed"},
		Why: "Choosing among registered implementations is a Tier 1 change."},
	{KP: 1, Q: "What must change in the core to add a new implementation of an existing seam?",
		A: "Nothing: a new class is written and registered with the factory", D: []string{"The base class gains a method for it", "The core's dispatch switch gets a new case", "The propagation loop is rewritten"},
		Why: "The core never knows which implementation it calls."},
	{KP: 1, Q: "Why does the core hold a pointer to the base rather than to a concrete class?",
		A: "So implementations can be swapped without the core knowing which one it calls", D: []string{"Because concrete classes cannot be compiled", "To make the core slower", "So plug-ins can modify the core"},
		Why: "Dependence on the abstraction is what makes the seam a seam."},
	{KP: 2, Q: "In a timer with a DerateProvider seam, what question does the core ask at the seam?",
		A: "What derate applies to this arc in this context", D: []string{"Which paths should be reported", "What the slack of this endpoint is", "How arrival times propagate"},
		Why: "Propagation stays in the core; only the derate is delegated."},
	{KP: 2, Q: "What does a POCV provider return through the derate seam, compared with flat OCV or AOCV?",
		A: "A mean and sigma adjustment rather than a single factor", D: []string{"A required time", "The same factor for every arc", "A list of false paths"},
		Why: "POCV works with distributions, so the seam carries a mean and a sigma."},
	{KP: 2, Q: "Which skills sit behind a derate seam?",
		A: "Flat OCV, AOCV and POCV providers, with a custom margin layer stacking on top", D: []string{"Placement, routing and extraction engines", "SDC readers and Liberty parsers", "The arrival and required propagators"},
		Why: "Each margining scheme is a provider; propagation is untouched."},
	{KP: 2, Q: "With the derate seam in place, where do most customer derating differences end up?",
		A: "In configuration and tables: which provider, which tables, how many sigma", D: []string{"In core edits", "In new Liberty files", "In SDC exceptions"},
		Why: "Most differences are data and parameters; the rest are small providers."},
	{KP: 3, Q: "A plug-in keeps a global counter and returns different results on a second run of the same design. Which rule did it break?",
		A: "No hidden mutable global state; plug-ins must be deterministic", D: []string{"Plug-ins must be small", "Plug-ins must use the core's language", "Plug-ins must run only once"},
		Why: "Hidden state breaks reproducibility and differential tests."},
	{KP: 3, Q: "A plug-in reads another plug-in's internal tables directly. Which rule does this break?",
		A: "No cross-skill coupling: plug-ins communicate only through the seam's plain data", D: []string{"Frozen ABI", "Coarse batching", "Factory registration"},
		Why: "Coupled skills cannot be versioned, tested or promoted independently."},
	{KP: 3, Q: "Why is the hot kernel kept inside the core rather than in plug-ins?",
		A: "So the compiler can inline and vectorize it, with plug-ins called coarsely per batch", D: []string{"Because plug-ins cannot loop", "Because the kernel is secret", "Because kernels change often"},
		Why: "Performance-critical loops must not be split by a call boundary."},
	{KP: 3, Q: "What form should data take when it crosses a seam?",
		A: "Plain structures such as the design state or a batch of contexts", D: []string{"Pointers into the core's private objects", "Text that the plug-in parses", "Global variables set before the call"},
		Why: "Plain data keeps the boundary explicit, testable and stable."},
}

func seamGranularity(ctx generator.GeneratorContext) generator.Problem {
	k := float64(edakit.Int(3, 20))
	c := float64(5 * edakit.Int(2, 12))
	what := edakit.Pick([]string{"derate lookup per timing arc", "cost evaluation per routing edge", "move evaluation per cell", "capacitance lookup per wire segment"})
	switch kp(ctx) {
	case 1:
		f := 100 * c / (k + c)
		q := edakit.Sprintf("A seam calls a plug-in once per element for a %s. The core's own work costs %s ns per element and each crossing of the plug-in boundary costs %s ns. What percentage of runtime is spent crossing the boundary? Round to 1 decimal place.", what, edakit.Num(k, 0), edakit.Num(c, 0))
		exp := edakit.Sprintf("$c/(k+c) = %s/%s = %s\\%%$, before counting the vectorization the per-element call prevents.", edakit.Num(c, 0), edakit.Num(k+c, 0), edakit.Num(f, 1))
		return edakit.Numeric(q, edakit.Num(f, 1), exp, nil)
	case 2:
		n := edakit.Pick([]int{200000, 500000, 1000000, 2000000})
		b := edakit.Pick([]int{64, 128, 256, 512, 1024})
		calls := (n + b - 1) / b
		f := 100 * float64(calls) * c / (float64(n) * k)
		q := edakit.Sprintf("A seam for a %s is batched: the core passes %d elements per call over %d elements. Each crossing costs %s ns and the kernel work is %s ns per element. What is the boundary overhead as a percentage of kernel time? Round to 3 decimal places.", what, b, n, edakit.Num(c, 0), edakit.Num(k, 0))
		exp := edakit.Sprintf("$\\lceil %d/%d \\rceil = %d$ crossings cost %s ns against %s ns of kernel work: %s%%.", n, b, calls, edakit.Num(float64(calls)*c, 0), edakit.Num(float64(n)*k, 0), edakit.Num(f, 3))
		return edakit.Numeric(q, edakit.Num(f, 3), exp, nil)
	}
	x := edakit.Pick([]float64{0.5, 1, 2, 5})
	bmin := c / (x / 100 * k)
	b := math.Ceil(bmin - 1e-9)
	q := edakit.Sprintf("A seam for a %s crosses the plug-in boundary once per batch. Each crossing costs %s ns and the kernel work is %s ns per element. What is the smallest batch size that keeps the boundary overhead at or below %s%% of kernel time?", what, edakit.Num(c, 0), edakit.Num(k, 0), edakit.Num(x, 1))
	exp := edakit.Sprintf("Require $c/B \\le x k$: $B \\ge %s / (%s \\times %s) = %s$, so %s.", edakit.Num(c, 0), edakit.Num(x/100, 3), edakit.Num(k, 0), edakit.Num(bmin, 2), edakit.Num(b, 0))
	return edakit.Numeric(q, edakit.Num(b, 0), exp, nil)
}

var contextItems = []edakit.Item{
	{KP: 1, Q: "Every customer wants a different placement preference: grouping, abutment, region affinity. Where should the placer's seam sit?",
		A: "At the detailed-placement move evaluator, with legalization kept in the core", D: []string{"Inside the legalizer", "In the LEF reader", "In the router's cost function"},
		Why: "Preferences vary; legality does not. The seam goes where the variation is."},
	{KP: 1, Q: "Customers disagree on how to trade vias against wirelength and congestion. Where is the router's seam?",
		A: "At the edge cost and rip-up policy", D: []string{"At the DRC checker", "At the routing grid definition", "At the design database"},
		Why: "Cost trade-offs vary across teams while legality is settled."},
	{KP: 1, Q: "In an extractor, customers want different accuracy on different nets. What should the seam select?",
		A: "Per-net fidelity: field solver, pattern table or simplified model", D: []string{"The solver's equations", "Which conductors exist", "The netlist format"},
		Why: "Accuracy trade-offs vary; the solver mathematics does not."},
	{KP: 1, Q: "What question should an engine owner ask to find where a seam belongs for a customer behavior?",
		A: "What question would the core have to ask to accommodate this behavior", D: []string{"Which file is shortest", "Which function runs most often", "Which customer is largest"},
		Why: "The seam sits at the core's one question about the varying behavior."},
	{KP: 2, Q: "A derate seam passes cell, Vt class, early or late, rise or fall, slew and corner, but not logic depth or distance. Which scheme cannot be written as a provider?",
		A: "AOCV", D: []string{"Flat OCV", "A global late margin", "A per-corner factor"},
		Why: "AOCV tables are indexed by depth and distance, so the context must carry them."},
	{KP: 2, Q: "A derate seam does not say whether an arc is on a clock or a data path. Which customer rule becomes impossible?",
		A: "Keep AOCV for clocks but add margin on data paths", D: []string{"A flat 5% late derate everywhere", "A higher sigma count", "A per-corner global derate"},
		Why: "Clock-versus-data rules need that field in the context."},
	{KP: 2, Q: "A router cost seam passes the edge and its congestion but not the net's class. Which plug-in cannot be written cleanly?",
		A: "One that treats clock and analog nets differently from signal nets", D: []string{"One that penalizes congestion", "One that prefers shorter edges", "One that uses historical congestion"},
		Why: "Net-class-specific costs need the net class in the context."},
	{KP: 2, Q: "A seam lacks a field a plug-in needs, and the plug-in author reads it from a global. What is wrong?",
		A: "It creates hidden coupling; the field should be added to the seam's context through a reviewed change", D: []string{"Nothing, as long as it works", "Globals are the intended mechanism", "The plug-in should be moved into the core"},
		Why: "Missing context must be fixed at the seam, not worked around."},
	{KP: 3, Q: "A derate provider recomputes each path's logic depth by walking the graph itself. What does this suggest?",
		A: "The seam sits too low or passes too little context", D: []string{"The provider is well designed", "Depth belongs in Liberty", "The core should stop computing depth"},
		Why: "Plug-ins duplicating core work point to a missing context field."},
	{KP: 3, Q: "Profiling shows a seam's boundary overhead dominating a hot loop with tiny per-element payloads. What should change?",
		A: "Raise the seam to a batch, such as one call per level or per net", D: []string{"Inline the plug-in into the core permanently", "Remove the seam's context", "Call the plug-in twice per element"},
		Why: "A seam in the innermost loop should be lifted to a coarser natural unit."},
	{KP: 3, Q: "Two plug-ins at different seams must coordinate through a shared table to work. What does this indicate?",
		A: "The two seams were meant to be one, or the core is missing shared context", D: []string{"Shared tables are best practice", "Both plug-ins should be core", "The ABI should be unfrozen"},
		Why: "Plug-ins coordinating through side channels signal a mis-drawn boundary."},
	{KP: 3, Q: "Which request-log pattern shows a seam is well placed?",
		A: "Requests of that kind land in Tier 1 or Tier 2 and core edits for them essentially never happen", D: []string{"Every request needs a core edit", "No one uses the seam", "Plug-ins at the seam are rewritten every release"},
		Why: "The tier distribution is the measure of seam placement."},
}

var abiItems = []edakit.Item{
	{KP: 1, Q: "A core release adds a field in the middle of the context structure passed to plug-ins. Existing compiled plug-ins still compile from source. What happens to already-built plug-ins?",
		A: "They break, because the binary layout they were built against changed", D: []string{"Nothing, because the source still compiles", "They run faster", "They automatically use the new field"},
		Why: "The ABI includes structure layout; changing it breaks compiled plug-ins."},
	{KP: 1, Q: "What is the difference between an API and an ABI at a seam?",
		A: "The API is the source-level contract; the ABI is the binary contract of symbols, layouts and calling convention", D: []string{"They are the same thing", "The ABI is the documentation", "The API is only for Python"},
		Why: "Plug-ins loaded as shared libraries depend on the binary contract."},
	{KP: 1, Q: "How should new capability be added to a frozen extension interface?",
		A: "As a new versioned interface beside the old one, or a structure that carries its size and version", D: []string{"By reordering the virtual methods", "By changing parameter types in place", "By removing old methods"},
		Why: "Compatible extension keeps old plug-ins working."},
	{KP: 1, Q: "Which change is compatible with a frozen plug-in ABI?",
		A: "Adding a separate, newly versioned interface the core also accepts", D: []string{"Reordering virtual methods", "Changing an int parameter to double", "Inserting a field before existing ones"},
		Why: "Only additive, versioned changes preserve the binary contract."},
	{KP: 2, Q: "A strange result appears a month after deployment. What lets the team reproduce it?",
		A: "Run provenance naming the exact configuration, plug-in and core versions", D: []string{"The model's memory of the run", "The latest versions of everything", "A screenshot of the result"},
		Why: "Versioned changes and provenance make old runs reproducible."},
	{KP: 2, Q: "What should a version record for an agent-made plug-in change contain?",
		A: "What changed, the request, the evidence that admitted it and who approved it", D: []string{"Only the new binary", "Only the model's prompt", "Only the date"},
		Why: "The record ties a change to its reason and its proof."},
	{KP: 2, Q: "Does a Tier 1 configuration change need a version?",
		A: "Yes: every change gets a revision, whatever its tier", D: []string{"No, configuration is not code", "Only if it fails", "Only for core parameters"},
		Why: "Unversioned configuration changes make results irreproducible."},
	{KP: 3, Q: "A newly promoted plug-in degrades QoR on one design. What is the rollback for a Tier 2 change?",
		A: "Unload the new plug-in version and reload the previous one", D: []string{"Revert the core release", "Delete the golden tests", "Retrain the model"},
		Why: "Each tier has its own rollback; for plug-ins it is a version swap."},
	{KP: 3, Q: "Why must old plug-in artifacts be kept immutable?",
		A: "So rollback can restore exactly the previous known-good state", D: []string{"To save disk space", "So agents can edit them later", "Because they cannot be deleted"},
		Why: "Rollback is only as good as the artifacts it restores."},
	{KP: 3, Q: "A plug-in built for ABI version 2 is loaded by a core that implements version 3, which is incompatible. What should happen?",
		A: "The core refuses it with a clear error", D: []string{"The core loads it and hopes for the best", "The core silently computes defaults", "The plug-in rewrites itself"},
		Why: "Declared ABI versions let the core reject incompatible plug-ins instead of crashing or computing wrong results."},
	{KP: 3, Q: "Why do Tier 3 changes need more evidence than Tier 2 changes before release?",
		A: "Rolling back a core release affects every engine user, while a plug-in rollback affects one skill", D: []string{"Because core code is shorter", "Because Tier 3 changes cannot be rolled back at all", "Because plug-ins are never tested"},
		Why: "Blast radius and rollback cost grow with the tier."},
	{KP: 3, Q: "A newer core must keep loading older plug-ins. What makes this possible?",
		A: "A frozen ABI within the major version and versioned interfaces for new capability", D: []string{"Recompiling every plug-in at every site", "Removing old interfaces quickly", "Storing plug-ins as source only"},
		Why: "Compatibility is designed in at the seam."},
}
