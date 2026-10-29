package infra

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerFrontier) }

func registerFrontier(reg *generator.Registry) {
	reg.Register("infra.frontier.four_layers_in_code", &edakit.Bank{Items: fourLayersInCodeItems})
	reg.Register("infra.frontier.extension_abi", &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: extensionAbiItems}, edakit.Func(extensionAbi)}})
	reg.Register("infra.frontier.golden_harness", edakit.Func(goldenHarness))
	reg.Register("infra.frontier.sandboxed_builds", &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: sandboxedBuildsItems}, edakit.Func(sandboxedBuilds)}})
	reg.Register("infra.frontier.deterministic_engines", &edakit.Bank{Items: deterministicEnginesItems})
	reg.Register("infra.frontier.mcp_openroad", &edakit.Bank{Items: mcpOpenroadItems})
}

// extensionAbi compares a per-cell plug-in call with a per-batch call at the same seam.
func extensionAbi(ctx generator.GeneratorContext) generator.Problem {
	n := 100000 * edakit.Scaled(ctx, 2, 20)
	c := float64(edakit.Int(4, 30))
	w := float64(edakit.Int(5, 40))
	b := edakit.Pick([]int{256, 512, 1024, 4096})
	batches := int(math.Ceil(float64(n) / float64(b)))
	stem := edakit.Sprintf("A placer calls a plug-in cost function at an extension point. Each call across the seam costs %s ns of overhead (indirect call, argument marshalling, no inlining), and the cost work itself is %s ns per cell. One pass evaluates %d cells. ",
		edakit.Num(c, 0), edakit.Num(w, 0), n)
	perCell := float64(n) * c / 1e6
	perBatch := float64(batches) * c / 1e3
	speed := (float64(n)*w + float64(n)*c) / (float64(n)*w + float64(batches)*c)
	switch edakit.Int(0, 2) {
	case 0:
		ans, d := numOptions(perCell, 2, "ms", perCell*1000, float64(n)*(c+w)/1e6, perCell/10)
		return edakit.Choice(stem+"If the core calls the plug-in once per cell, how much call overhead does one pass pay, in milliseconds? (Round to 2 decimal places.)", ans,
			edakit.Sprintf("%d calls × %s ns = %s ms of pure overhead per pass.", n, edakit.Num(c, 0), edakit.Num(perCell, 2)), d...)
	case 1:
		ans, d := numOptions(perBatch, 2, "microseconds", perCell*1000, float64(n)/float64(b)*w/1e3, perBatch*float64(b)/100)
		return edakit.Choice(stem+edakit.Sprintf("If instead the core passes cells in batches of %d per call, how much call overhead does one pass pay, in microseconds? (Round to 2 decimal places.)", b), ans,
			edakit.Sprintf("ceil(%d / %d) = %d calls × %s ns = %s microseconds.", n, b, batches, edakit.Num(c, 0), edakit.Num(perBatch, 2)), d...)
	}
	ans, d := numOptions(speed, 2, "", (w+c)/w, c/w+0.5, float64(b)/10)
	return edakit.Choice(stem+edakit.Sprintf("Moving from one call per cell to batches of %d per call, what is the speedup of the whole pass (overhead plus work)? (Round to 2 decimal places.)", b), ans,
		edakit.Sprintf("Per-cell: %d × (%s + %s) ns. Batched: %d × %s ns of work plus %d × %s ns of calls. The ratio is %s.", n, edakit.Num(w, 0), edakit.Num(c, 0), n, edakit.Num(w, 0), batches, edakit.Num(c, 0), edakit.Num(speed, 2)), d...)
}

// sandboxedBuilds computes harness iteration time with a compiler cache.
func sandboxedBuilds(ctx generator.GeneratorContext) generator.Problem {
	u := edakit.Scaled(ctx, 6, 40)
	h := float64(edakit.Int(6, 19)*5) / 100
	hit := float64(edakit.Int(1, 4)) / 10
	miss := float64(edakit.Int(8, 30))
	link := float64(edakit.Int(3, 15))
	test := float64(edakit.Int(10, 60))
	iter := float64(u)*(h*hit+(1-h)*miss) + link + test
	cold := float64(u)*miss + link + test
	stem := edakit.Sprintf("An agent harness rebuilds a sandboxed plug-in after each edit. The build has %d compilation units; the compiler cache hit rate is %s. A cache hit costs %s s per unit, a miss %s s per unit; the link takes %s s and the golden tests take %s s. ",
		u, edakit.Num(h*100, 0)+"%", edakit.Num(hit, 1), edakit.Num(miss, 0), edakit.Num(link, 0), edakit.Num(test, 0))
	switch edakit.Int(0, 2) {
	case 0:
		ans, d := numOptions(iter, 1, "s", cold, float64(u)*h*hit+link+test, float64(u)*(h*miss+(1-h)*hit)+link+test)
		return edakit.Choice(stem+"How long is one edit-build-test iteration, in seconds? (Round to 1 decimal place.)", ans,
			edakit.Sprintf("%d × (%s × %s + %s × %s) + %s + %s = %s s.", u, edakit.Num(h, 2), edakit.Num(hit, 1), edakit.Num(1-h, 2), edakit.Num(miss, 0), edakit.Num(link, 0), edakit.Num(test, 0), edakit.Num(iter, 1)), d...)
	case 1:
		k := int(math.Floor(3600 / iter))
		kc := int(math.Floor(3600 / cold))
		d := edakit.Distinct(k, 3, kc, int(math.Ceil(3600/iter))+1, int(math.Floor(3600/(iter-test))))
		return edakit.Choice(stem+"How many complete iterations can the harness run in one hour?", edakit.IntStr(k),
			edakit.Sprintf("One iteration takes %s s, so floor(3600 / %s) = %d iterations fit in an hour.", edakit.Num(iter, 1), edakit.Num(iter, 1), k), d...)
	}
	sp := cold / iter
	ans, d := numOptions(sp, 2, "", float64(u)*miss/(float64(u)*(h*hit+(1-h)*miss)), 1/(1-h), miss/hit)
	return edakit.Choice(stem+"What is the speedup of one iteration over a cold build with no cache hits? (Round to 2 decimal places.)", ans,
		edakit.Sprintf("Cold: %d × %s + %s + %s = %s s. Cached: %s s. Ratio %s; the fixed link and test time caps the gain.", u, edakit.Num(miss, 0), edakit.Num(link, 0), edakit.Num(test, 0), edakit.Num(cold, 1), edakit.Num(iter, 1), edakit.Num(sp, 2)), d...)
}

// goldenHarness asks for relative change (KP 1), tolerance violations (KP 2) and unintended
// movers (KP 3) in a differential QoR comparison.
func goldenHarness(ctx generator.GeneratorContext) generator.Problem {
	return goldenHarnessAttempt(ctx, 0)
}

func goldenHarnessAttempt(ctx generator.GeneratorContext, attempt int) generator.Problem {
	designs := []string{"gcd", "aes", "ibex", "jpeg", "riscv32i", "uart", "spi", "dma", "fft"}
	switch variant(ctx) {
	case 1:
		base := float64(edakit.Int(200, 900) * 100)
		delta := float64(edakit.Int(-60, 60)*10) + float64(edakit.Int(1, 9))
		cand := base + delta
		metric := edakit.Pick([]string{"total wirelength (microns)", "standard-cell area (square microns)"})
		rel := (cand - base) / base * 100
		return edakit.Numeric(edakit.Sprintf("A golden harness compares a candidate build against the baseline on design %s. The baseline %s is %s and the candidate is %s. What is the relative change, in percent? (Round to 2 decimal places; negative means smaller.)",
			edakit.Pick(designs), metric, edakit.Num(base, 0), edakit.Num(cand, 0)), edakit.Num(rel, 2),
			edakit.Sprintf("(%s - %s) / %s × 100 = %s%%.", edakit.Num(cand, 0), edakit.Num(base, 0), edakit.Num(base, 0), edakit.Num(rel, 2)),
			map[string]string{"relative_percent": edakit.Num(rel, 2)})
	case 2:
		n := edakit.Scaled(ctx, 5, 9)
		tol := edakit.Pick([]float64{0.5, 1, 2})
		ds := edakit.Shuffle(designs)[:n]
		var parts []string
		flagged := 0
		var fl []string
		for _, d := range ds {
			base := float64(edakit.Int(100, 900) * 100)
			var pct float64
			switch edakit.Int(0, 3) {
			case 0:
				pct = 0
			case 1:
				pct = tol * float64(edakit.Int(1, 9)) / 10
			default:
				pct = tol * (1 + float64(edakit.Int(2, 30))/10)
			}
			if edakit.Int(0, 1) == 0 {
				pct = -pct
			}
			cand := math.Round(base * (1 + pct/100))
			rel := (cand - base) / base * 100
			if math.Abs(rel) > tol {
				flagged++
				fl = append(fl, d)
			}
			parts = append(parts, edakit.Sprintf("%s %s → %s", d, edakit.Num(base, 0), edakit.Num(cand, 0)))
		}
		return edakit.Numeric(edakit.Sprintf("A differential QoR gate flags any design whose total wirelength moves by more than %s%% (either direction) from its golden baseline. Baseline → candidate wirelength in microns: %s. How many designs are flagged?",
			edakit.Num(tol, 1), strings.Join(parts, "; ")), edakit.IntStr(flagged),
			edakit.Sprintf("Compute |candidate - baseline| / baseline for each design and compare with %s%%. Flagged: %s.", edakit.Num(tol, 1), names(fl)),
			map[string]string{"flagged": edakit.IntStr(flagged)})
	}
	n := edakit.Scaled(ctx, 5, 9)
	tol := float64(edakit.Int(1, 3))
	ds := edakit.Shuffle(designs)[:n]
	k := edakit.Int(1, 2)
	intended := map[string]bool{}
	for _, d := range ds[:k] {
		intended[d] = true
	}
	var parts, bad []string
	for _, d := range edakit.Shuffle(ds) {
		base := -float64(edakit.Int(0, 60))
		var dv float64
		if intended[d] {
			dv = float64(edakit.Int(5, 25))
		} else if edakit.Int(0, 2) == 0 {
			dv = tol + float64(edakit.Int(1, 8))
		} else {
			dv = float64(edakit.Int(0, int(tol)))
		}
		if edakit.Int(0, 1) == 0 {
			dv = -dv
		}
		if !intended[d] && math.Abs(dv) > tol {
			bad = append(bad, d)
		}
		parts = append(parts, edakit.Sprintf("%s %s → %s", d, edakit.Num(base, 0), edakit.Num(base+dv, 0)))
	}
	if len(bad) == 0 && attempt < generator.MaxDrawAttempts && edakit.Int(0, 4) > 0 {
		return goldenHarnessAttempt(ctx, attempt+1)
	}
	var in []string
	for _, d := range ds[:k] {
		in = append(in, d)
	}
	return edakit.Numeric(edakit.Sprintf("A new derate plug-in is meant to change timing only on %s. The differential harness treats a WNS change of more than %s ps as movement. Baseline → candidate WNS in ps: %s. How many designs moved that were not supposed to?",
		names(in), edakit.Num(tol, 0), strings.Join(parts, "; ")), edakit.IntStr(len(bad)),
		edakit.Sprintf("Skip the intended designs (%s) and count the others whose |ΔWNS| exceeds %s ps: %s.", names(in), edakit.Num(tol, 0), names(bad)),
		map[string]string{"unintended": edakit.IntStr(len(bad))})
}

var fourLayersInCodeItems = []edakit.Item{
	{KP: 1, Q: "In an OpenROAD-like placer, which code belongs in the immutable core?",
		A: "The legalizer that guarantees overlap-free cells on legal sites", D: []string{"The weight on timing-driven net criticality", "A Tcl hook that prints density each iteration", "A cost function for abutting regular cells"},
		Why: "Legality is an invariant every result must satisfy, so it lives in the core. Weights are parameters, hooks observe, cost functions are policies."},
	{KP: 1, Q: "A router's rip-up-and-reroute ordering is implemented as a class selected by name at startup. Which layer is it?",
		A: "Strategy or policy", D: []string{"Immutable core", "Parameters", "Observation hooks"},
		Why: "A selectable heuristic that decides how the engine searches is policy."},
	{KP: 1, Q: "A timer exposes a callback fired after every incremental update with the new WNS and TNS. Which layer is this?",
		A: "Observation hooks", D: []string{"Immutable core", "Strategy or policy", "Parameters"},
		Why: "A point where an agent reads live state is an observation hook."},
	{KP: 1, Q: "The placer's target density is read from a config file at run time. Which layer is it?",
		A: "Parameters", D: []string{"Strategy or policy", "Immutable core", "Observation hooks"},
		Why: "A value that changes behavior without recompiling is a parameter."},
	{KP: 2, Q: "Why does a seam use an abstract base class plus a factory instead of letting plug-ins include core headers freely?",
		A: "The core depends only on a small stable interface, not on each plug-in", D: []string{"Factories make code run faster than direct calls", "Abstract classes remove the need for tests", "It lets plug-ins change the core's invariants"},
		Why: "A narrow interface and a factory decouple the core from implementations, so new plug-ins never force core edits."},
	{KP: 2, Q: "A plug-in interface method takes a pointer to the whole design database. What is the main risk?",
		A: "Plug-ins can depend on, or mutate, anything, so the seam stops being narrow", D: []string{"The call becomes too fast to measure", "The compiler cannot compile pointers", "The database becomes read-only"},
		Why: "Passing the world couples plug-ins to every internal detail. A seam should pass only the context it needs."},
	{KP: 2, Q: "A data-oriented seam passes a plug-in which kind of argument?",
		A: "Plain arrays of the values it needs, such as cell widths and positions", D: []string{"Live iterators into the core's private containers", "A global pointer to the engine singleton", "A Tcl string to evaluate"},
		Why: "Plain data is easy to version, test and batch, and keeps the hot loop's layout in the core's control."},
	{KP: 2, Q: "Why should the hot numerical kernel stay inside the core rather than behind a plug-in call?",
		A: "So the compiler can inline and vectorize it", D: []string{"So plug-ins can replace it per cell", "Because kernels cannot be tested", "Because cores cannot call functions"},
		Why: "Crossing a seam per element blocks inlining and vectorization. Keep the kernel in the core and call plug-ins coarsely."},
	{KP: 3, Q: "A customer asks to raise the congestion weight from 0.3 to 0.5 for one design. Where does the change land?",
		A: "A parameter change, with no recompile", D: []string{"A new policy plug-in", "A core edit with human approval", "A new observation hook"},
		Why: "Changing an existing weight is configuration, the smallest tier."},
	{KP: 3, Q: "A customer wants macros kept 20 microns from the die edge, but the placer has no such rule. The constraint system has an extension point. Where does it land?",
		A: "A policy plug-in at the constraint extension point", D: []string{"A parameter change", "A core edit to legalization", "An observation hook"},
		Why: "New behavior at a defined extension point is a plug-in. The core stays untouched."},
	{KP: 3, Q: "A request requires changing how the legalizer defines a legal site. Where does it land?",
		A: "A core edit, needing human approval and full regression", D: []string{"A parameter change", "A policy plug-in", "An observation hook"},
		Why: "Legality is the core invariant. Changing it is the rare, approved tier."},
	{KP: 3, Q: "An agent wants to stop global placement early when overflow falls below 10%. What does it need from the engine?",
		A: "A hook exposing overflow each iteration plus a stop control", D: []string{"A core edit to the density solver", "A new legalizer", "A larger target density only"},
		Why: "Reading live state and injecting a decision is what observation hooks provide."},
}

var extensionAbiItems = []edakit.Item{
	{KP: 1, Q: "A plug-in compiled against version 1 of an extension ABI stops loading after a minor core release. What promise was broken?",
		A: "The frozen ABI: layouts and signatures at the seam must not change", D: []string{"The promise that the core never gets faster", "The promise that plug-ins are written in Tcl", "The promise that tests always pass"},
		Why: "A frozen ABI means binaries built against it keep working. Changing struct layouts or signatures breaks them."},
	{KP: 1, Q: "Why do many plug-in systems expose a C ABI even when the core is written in C++?",
		A: "C calling conventions and struct layouts are stable across compilers", D: []string{"C is always faster than C++", "C++ cannot call C", "C ABIs allow exceptions to cross the boundary"},
		Why: "C++ name mangling, vtables and standard-library layouts vary by compiler and version. A C boundary is far more stable."},
	{KP: 1, Q: "An extension struct begins with a size field and a version number. What does that allow?",
		A: "Adding fields later while old plug-ins keep working", D: []string{"Removing fields at will", "Skipping validation of inputs", "Calling plug-ins without a header"},
		Why: "The core can tell which fields a plug-in knows about and fill or ignore the rest."},
	{KP: 1, Q: "Which change to a published extension interface is safe without a new major version?",
		A: "Adding a new optional function at the end of the table", D: []string{"Reordering the existing functions", "Changing an int argument to a double", "Removing a deprecated field"},
		Why: "Appending keeps existing offsets and signatures. Reordering, retyping or removing breaks compiled plug-ins."},
	{KP: 1, Q: "Why should C++ exceptions not cross a plug-in ABI boundary?",
		A: "Exception handling is not portable across compilers and runtimes", D: []string{"Exceptions are always faster than error codes", "C++ forbids exceptions in libraries", "Exceptions make the ABI frozen"},
		Why: "Unwinding across modules built by different toolchains can crash. Return error codes at the seam."},
	{KP: 3, Q: "A plug-in caches the last design it saw in a static variable. What goes wrong?",
		A: "Results depend on call history, breaking determinism and parallel use", D: []string{"Nothing, statics are always safe", "The plug-in becomes stateless", "The compiler rejects static variables"},
		Why: "Hidden state makes the output depend on what ran before and races under threads. Pass state explicitly."},
	{KP: 3, Q: "A seam passes a plug-in a std::vector of internal Cell objects by reference. Why is that risky for a frozen ABI?",
		A: "The vector and Cell layouts may change between compilers or core versions", D: []string{"Vectors cannot be passed by reference", "Plain arrays are always slower", "It makes the plug-in stateless"},
		Why: "Standard-library and class layouts are not ABI-stable. Pointers plus lengths to plain structs are."},
	{KP: 3, Q: "Why should a cost plug-in receive cell positions as read-only input and return costs in a separate output array?",
		A: "So the core can verify and test it as a pure function", D: []string{"So it can move cells without the legalizer", "Because outputs must be global", "To make it call the core back per cell"},
		Why: "A pure function of plain inputs is easy to golden-test, run in parallel and reason about."},
	{KP: 3, Q: "Two plug-ins give different answers when the core calls them in a different order. What does this reveal?",
		A: "At least one plug-in keeps hidden shared state", D: []string{"The ABI is frozen correctly", "The core has no seam", "Order never affects pure functions"},
		Why: "Pure plug-ins commute. Order dependence points to hidden state."},
	{KP: 3, Q: "Which design makes a derate plug-in easiest to run on many threads at once?",
		A: "A function of its explicit context arguments with no globals", D: []string{"One that reads a global corner variable", "One that logs to a shared buffer without locks", "One that caches per-thread results in a static map"},
		Why: "Reentrant code with explicit inputs is thread-safe by construction."},
}

var sandboxedBuildsItems = []edakit.Item{
	{KP: 1, Q: "Why does a harness build and run model-written plug-ins inside a sandbox with no network?",
		A: "So code cannot leak design data or fetch unvetted dependencies", D: []string{"Because compilers need no files", "To make builds slower and safer", "Because plug-ins never read input"},
		Why: "Design data is IP, and network access lets code exfiltrate it or pull in code nobody reviewed."},
	{KP: 1, Q: "A generated plug-in enters an infinite loop during golden tests. What sandbox feature contains it?",
		A: "A wall-clock or CPU time limit that kills the job", D: []string{"A larger memory limit", "Network access to report the hang", "Running it in the core process"},
		Why: "Time limits turn a hang into a reported failure the model can learn from."},
	{KP: 1, Q: "A plug-in allocates memory without bound during a test. What should the sandbox do?",
		A: "Enforce a memory limit so the job fails without harming the host", D: []string{"Give it all the host memory", "Silently restart the host", "Ignore it until the test ends"},
		Why: "Resource limits, for example through Linux control groups, protect the machine and make failures explicit."},
	{KP: 1, Q: "Why is a modular build, where one plug-in compiles on its own, important for the harness loop?",
		A: "An edit rebuilds only that plug-in, so iterations take seconds", D: []string{"It removes the need for tests", "It lets the plug-in edit the core", "It forces a full rebuild every time"},
		Why: "If a plug-in edit triggered a full core rebuild, the loop would take minutes per try."},
	{KP: 1, Q: "Should a plug-in build in the sandbox be able to write into the core's source tree?",
		A: "No, it should see the core only as read-only headers and libraries", D: []string{"Yes, so it can patch the core when needed", "Yes, as long as tests pass", "Only during the link step"},
		Why: "Read-only access keeps Tier 2 changes from silently becoming Tier 3 core edits."},
	{KP: 3, Q: "A build passes in the sandbox on Monday and fails Tuesday with no code change. The toolchain was fetched as latest each time. What fixes this?",
		A: "Pin the compiler and libraries to exact versions stored offline", D: []string{"Always fetch the newest compiler", "Disable the failing test", "Add network access to the sandbox"},
		Why: "Pinned, offline toolchains make builds reproducible and independent of upstream changes."},
	{KP: 3, Q: "Why should a promoted skill record the exact compiler version and flags that built it?",
		A: "So the same binary can be rebuilt and audited later", D: []string{"So users can choose random flags", "Because the version sets the clock period", "So the tests can be skipped"},
		Why: "Reproducible builds need the full recipe, not just the source."},
	{KP: 3, Q: "A container image holds the pinned toolchain for plug-in builds. What property makes it trustworthy across sites?",
		A: "It is referenced by an immutable content digest", D: []string{"It is tagged latest", "It downloads tools at start-up", "It includes the customer's design"},
		Why: "A digest names exact bytes; a moving tag can change underneath you."},
	{KP: 3, Q: "Two machines build the same plug-in source with the same pinned toolchain and get different binaries. What is a common cause?",
		A: "Embedded timestamps or absolute paths in the build", D: []string{"Different numbers of golden tests", "The plug-in uses a frozen ABI", "The sandbox has no network"},
		Why: "Reproducible builds strip timestamps and normalize paths so bytes match."},
	{KP: 3, Q: "Why keep a local mirror of the dependency sources a plug-in may use?",
		A: "Builds stay offline and unchanged even if upstream disappears", D: []string{"So the sandbox can browse the web", "To speed up the golden tests", "So the model can pick newer versions"},
		Why: "An offline mirror makes the build hermetic and reproducible."},
}

var deterministicEnginesItems = []edakit.Item{
	{KP: 1, Q: "A placer run with seed 7 on 8 threads gives a different HPWL than seed 7 on 16 threads. What determinism property fails?",
		A: "Results must not depend on the thread count", D: []string{"Results must not depend on the seed", "Results must be identical across designs", "Results must be faster with more threads"},
		Why: "A deterministic parallel engine fixes partitions and merge order so the thread count changes speed, not results."},
	{KP: 1, Q: "Where should an engine's random choices come from?",
		A: "A random generator created from an explicit seed and passed to the code", D: []string{"The global time of day", "A global generator shared by all threads in arrival order", "Addresses of allocated objects"},
		Why: "An explicit seeded generator, passed down, makes every random choice reproducible."},
	{KP: 1, Q: "A router breaks ties between equal-cost paths by comparing pointer addresses. Why is that a problem?",
		A: "Addresses change between runs, so ties resolve differently", D: []string{"Pointers are always equal", "It makes the router slower", "Ties never occur in routing"},
		Why: "Tie-breaking must use stable keys such as object ids, not memory addresses."},
	{KP: 1, Q: "An engine iterates a hash map of nets to decide routing order. Results vary between runs. What fix keeps the speed and restores determinism?",
		A: "Sort the nets by a stable key before ordering them", D: []string{"Use a larger hash table", "Add more threads", "Seed the hash with the clock"},
		Why: "Hash iteration order can vary; a stable sort makes the order a function of the data."},
	{KP: 2, Q: "What should a result manifest record so that a QoR number can be reproduced?",
		A: "Engine version, parameters, seed, inputs' hashes and thread count", D: []string{"Only the final QoR number", "The name of the person who ran it", "The machine's hostname only"},
		Why: "Reproduction needs every input that affects the result."},
	{KP: 2, Q: "Two results differ, and both manifests show the same version, seed and inputs. What is the likeliest explanation?",
		A: "The engine has a hidden source of nondeterminism", D: []string{"The manifests are always wrong", "Seeds have no effect", "QoR is random by design"},
		Why: "Identical recorded inputs with different outputs means something unrecorded, such as order or time, affects the result."},
	{KP: 2, Q: "Why hash the input design files in the manifest instead of storing only their paths?",
		A: "A path can point to changed content; a hash identifies the exact bytes", D: []string{"Hashes are shorter to type", "Paths are illegal in manifests", "Hashes speed up the engine"},
		Why: "Content hashes catch silent edits to inputs."},
	{KP: 2, Q: "A flow runs placement with five seeds and keeps the best. What must the manifest store for the chosen result?",
		A: "The seed of the chosen run, not just the seed list", D: []string{"Only the number of seeds", "The worst seed", "Nothing, since the best is obvious"},
		Why: "To reproduce the chosen result you need its exact seed."},
	{KP: 3, Q: "Why does a differential test of a plug-in change require a deterministic engine?",
		A: "Otherwise run-to-run noise looks like the change's effect", D: []string{"Differential tests need random noise", "Determinism makes tests unnecessary", "Only accurate engines need determinism"},
		Why: "If baseline reruns already differ, a difference after the change proves nothing."},
	{KP: 3, Q: "An RL reward comes from running a placer on the model's output. The placer is nondeterministic. What happens to learning?",
		A: "The reward is noisy, so the model learns slower and may chase luck", D: []string{"Learning becomes faster", "The reward becomes exact", "Nothing changes"},
		Why: "Noise in a tool-grounded reward blurs which actions were actually better."},
	{KP: 3, Q: "An agent reruns a failed step to debug it, and the failure disappears. Which engine property would have made this debuggable?",
		A: "Same inputs and seed reproduce the same failure", D: []string{"Faster runtime", "More threads", "A larger log file"},
		Why: "Reproducible failures can be bisected and fixed; vanishing ones cannot."},
	{KP: 3, Q: "An agent compares two parameter sets with one run each on a seed-sensitive engine. What should it do instead?",
		A: "Run several seeds per setting and compare distributions", D: []string{"Use a single seed, chosen at random", "Compare only the first iteration", "Disable the seed"},
		Why: "When results vary by seed, one run per setting cannot separate the parameter effect from seed noise."},
}

var mcpOpenroadItems = []edakit.Item{
	{KP: 1, Q: "What does wrapping an OpenROAD action as an MCP tool add over calling its Tcl command directly?",
		A: "A typed input and output schema that any agent can discover and call", D: []string{"A faster placement algorithm", "A new database format", "A replacement for the timer"},
		Why: "MCP tools advertise names, descriptions and JSON schemas, so agents can call engines without knowing Tcl."},
	{KP: 1, Q: "An MCP tool report_timing returns its result. Which form is most useful to an agent?",
		A: "Structured fields such as WNS, TNS and a list of worst paths", D: []string{"A screenshot of the GUI", "The raw log with color codes", "Nothing; it only prints to stdout"},
		Why: "Structured output can be parsed and compared reliably."},
	{KP: 1, Q: "A tool's input schema declares utilization as a number between 0.1 and 0.9. What does that buy?",
		A: "Bad arguments are rejected before the engine runs", D: []string{"The engine runs twice as fast", "The agent never needs to read results", "The DB is saved automatically"},
		Why: "Validation at the boundary turns a malformed call into an immediate, clear error."},
	{KP: 1, Q: "Why should an MCP tool's description state its preconditions, such as requiring a placed design?",
		A: "The agent reads descriptions to choose and order tool calls", D: []string{"Descriptions are compiled into the engine", "Preconditions are enforced by MCP itself", "It hides the tool from agents"},
		Why: "Like a skill manifest, a description is how the caller decides what to call."},
	{KP: 2, Q: "A server exposes a single tool eval_tcl that runs any Tcl string. What is the main drawback?",
		A: "No typed contract and no safety: anything can run, and results are unstructured", D: []string{"It is too slow to start", "It cannot read the database", "It makes the tools deterministic"},
		Why: "A raw eval gives power but loses validation, structure and the ability to sandbox or audit actions."},
	{KP: 2, Q: "Which tool granularity is usually best for an agent driving a flow?",
		A: "Coarse actions such as run global placement or report timing", D: []string{"One tool per database object accessor", "One tool per cell move", "A single tool for the whole flow with no options"},
		Why: "Coarse tools keep calls few and meaningful; per-object tools flood the agent with round trips."},
	{KP: 2, Q: "An agent must move 10,000 cells. Why is a move_cell tool called once per cell a poor design?",
		A: "Ten thousand round trips dominate time and context", D: []string{"MCP forbids loops", "Cells cannot be moved", "It makes the placement exact"},
		Why: "Batch operations take a list in one call, as with coarse plug-in calls at a seam."},
	{KP: 2, Q: "A tool run_detailed_route can take an hour. How should its interface handle that?",
		A: "Start a job and return an id, with a separate tool to poll status", D: []string{"Block silently for an hour", "Return a fake result immediately", "Split it into one call per net"},
		Why: "Long jobs need asynchronous handling so the agent can wait, poll or cancel."},
	{KP: 3, Q: "Why does an OpenROAD MCP server keep a session that holds a loaded database between calls?",
		A: "So each call works on the current design without reloading it", D: []string{"So every call starts from RTL", "Because MCP forbids files", "To prevent the agent from reading results"},
		Why: "Loading a large ODB is costly; a session keeps the state for a sequence of tool calls."},
	{KP: 3, Q: "Where should the OpenROAD process behind an MCP server run when the design is customer IP?",
		A: "On-site, in a sandbox with limits and no outbound network", D: []string{"On any public server", "Inside the model provider's cloud by default", "In the agent's chat window"},
		Why: "The IP boundary is part of the architecture: engines and data stay on-site."},
	{KP: 3, Q: "An agent calls set_place_density twice with the same value. Why should the tool be idempotent?",
		A: "Retries then leave the design in the same state", D: []string{"Idempotent tools run faster", "The second call must undo the first", "Agents never retry"},
		Why: "Agents retry after timeouts; idempotent tools make retries safe."},
	{KP: 3, Q: "Why should the server log every tool call with its arguments and the result's checksum?",
		A: "So the run can be audited, replayed and reverted", D: []string{"To slow down the agent", "Because MCP requires logs", "To store the model's weights"},
		Why: "A versioned, replayable record keeps agent actions reviewable and reversible."},
}
