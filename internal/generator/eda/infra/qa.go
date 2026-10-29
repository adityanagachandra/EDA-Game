package infra

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerQA) }

func registerQA(reg *generator.Registry) {
	reg.Register("infra.qa.regression_qor", edakit.Func(regressionQor))
	reg.Register("infra.qa.qor_significance", edakit.Func(qorSignificance))
	reg.Register("infra.qa.build_systems", edakit.Func(buildSystems))
	reg.Register("infra.qa.ci_pipelines", edakit.Func(ciPipelines))
	reg.Register("infra.qa.bisection", edakit.Func(bisection))
	reg.Register("infra.qa.cpp_rust", &edakit.Bank{Items: cppRustItems})
}

var qaDesigns = []string{"gcd", "aes", "ibex", "jpeg", "riscv32i", "swerv", "uart", "spi", "dynamic_node", "tinyRocket"}

// regressionQor asks for the geometric mean of QoR ratios, its percent form and its inverse.
func regressionQor(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 6)
	ds := edakit.Shuffle(qaDesigns)[:n]
	mi := edakit.Int(0, 3)
	metric := []string{"routed wirelength in mm", "standard-cell area in square microns", "total power in mW", "run time in s"}[mi]
	var parts, rs []string
	logSum := 0.0
	arith := 0.0
	for _, d := range ds {
		old := 10 * edakit.Int(20, 200)
		pct := edakit.Int(-15, 15)
		nw := int(math.Round(float64(old) * (1 + float64(pct)/100)))
		r := float64(nw) / float64(old)
		logSum += math.Log(r)
		arith += r
		parts = append(parts, edakit.Sprintf("%s %d → %d", d, old, nw))
		rs = append(rs, edakit.Num(r, 4))
	}
	g := math.Exp(logSum / float64(n))
	q := edakit.Sprintf("A nightly regression compares a candidate build of a placer against the baseline on %d designs. The %s (baseline → candidate, lower is better) is: %s. ", n, metric, strings.Join(parts, ", "))
	facts := map[string]string{"geomean": edakit.Num(g, 4)}
	switch variant(ctx) {
	case 1:
		ans := edakit.Num(g, 3)
		return edakit.Numeric(q+"What is the geometric mean of the candidate-to-baseline ratios? Round to 3 decimal places.", ans,
			edakit.Sprintf("The ratios candidate/baseline are %s. Multiply them and take the root of order n = %d: %s.", strings.Join(rs, ", "), n, ans), facts)
	case 2:
		ans := edakit.Num((g-1)*100, 2)
		return edakit.Numeric(q+"Report the geometric-mean change of the candidate relative to the baseline as a signed percentage (negative means better). Round to 2 decimal places.", ans,
			edakit.Sprintf("The geometric mean ratio is %s, so the change is (%s - 1) × 100 = %s%%.", edakit.Num(g, 4), edakit.Num(g, 4), ans), facts)
	}
	ans := edakit.Num(1/g, 3)
	facts["arith_mean"] = edakit.Num(arith/float64(n), 4)
	return edakit.Numeric(q+"Now compare in the other direction: what is the geometric mean of the baseline-to-candidate ratios? Round to 3 decimal places.", ans,
		edakit.Sprintf("The geometric mean of reciprocals is the reciprocal of the geometric mean: 1/%s = %s. The arithmetic mean has no such symmetry: the mean of candidate/baseline is %s, and its reciprocal is not the mean of baseline/candidate.", edakit.Num(g, 4), ans, edakit.Num(arith/float64(n), 4)), facts)
}

// qorSignificance asks for the paired mean, standard deviation and t statistic across seeds.
func qorSignificance(ctx generator.GeneratorContext) generator.Problem {
	return qorSignificanceAttempt(ctx, 0)
}

func qorSignificanceAttempt(ctx generator.GeneratorContext, attempt int) generator.Problem {
	n := edakit.Scaled(ctx, 4, 7)
	var parts []string
	ds := make([]float64, n)
	sum := 0.0
	bias := edakit.Int(-6, 6)
	for i := 0; i < n; i++ {
		old := edakit.Int(400, 600)
		d := bias + edakit.Int(-5, 5)
		nw := old + d
		ds[i] = float64(d)
		sum += ds[i]
		parts = append(parts, edakit.Sprintf("seed %d: %d → %d", i+1, old, nw))
	}
	mean := sum / float64(n)
	ss := 0.0
	for _, d := range ds {
		ss += (d - mean) * (d - mean)
	}
	sd := math.Sqrt(ss / float64(n-1))
	if sd == 0 && attempt < generator.MaxDrawAttempts {
		return qorSignificanceAttempt(ctx, attempt+1)
	}
	t := mean / (sd / math.Sqrt(float64(n)))
	q := edakit.Sprintf("A placer change is run on one design with %d seeds, before and after the change, using the same seed for each pair. The routed wirelength in mm is %s. ", n, strings.Join(parts, ", "))
	facts := map[string]string{"mean_diff": edakit.Num(mean, 4), "sd": edakit.Num(sd, 4), "t": edakit.Num(t, 4)}
	var dl []string
	for _, d := range ds {
		dl = append(dl, edakit.Num(d, 0))
	}
	switch variant(ctx) {
	case 1:
		ans := edakit.Num(mean, 2)
		return edakit.Numeric(q+"What is the mean paired difference, after minus before, in mm? Round to 2 decimal places.", ans,
			edakit.Sprintf("The per-seed differences are %s; their mean is %s / %d = %s mm.", strings.Join(dl, ", "), edakit.Num(sum, 0), n, ans), facts)
	case 2:
		ans := edakit.Num(sd, 2)
		return edakit.Numeric(q+"What is the sample standard deviation of the paired differences (after minus before, dividing by n - 1), in mm? Round to 2 decimal places.", ans,
			edakit.Sprintf("Differences %s have mean %s. The squared deviations sum to %s; dividing by %d and taking the square root gives %s mm.", strings.Join(dl, ", "), edakit.Num(mean, 3), edakit.Num(ss, 3), n-1, ans), facts)
	}
	ans := edakit.Num(t, 2)
	return edakit.Numeric(q+"What is the paired t statistic of the differences (after minus before)? Round to 2 decimal places.", ans,
		edakit.Sprintf("Differences %s: mean %s, sample standard deviation %s. t = mean / (s / sqrt(n)) = %s / (%s / sqrt(%d)) = %s. Compare |t| with the t critical value for %d degrees of freedom before calling the change real.", strings.Join(dl, ", "), edakit.Num(mean, 3), edakit.Num(sd, 3), edakit.Num(mean, 3), edakit.Num(sd, 3), n, ans, n-1), facts)
}

// buildSystems draws an include graph and asks what a header edit rebuilds.
func buildSystems(ctx generator.GeneratorContext) generator.Problem {
	return buildSystemsAttempt(ctx, 0)
}

func buildSystemsAttempt(ctx generator.GeneratorContext, attempt int) generator.Problem {
	m := edakit.Scaled(ctx, 4, 6)
	k := edakit.Scaled(ctx, 5, 9)
	hn := edakit.Shuffle([]string{"db.h", "geom.h", "net.h", "timer.h", "util.h", "log.h", "graph.h", "lib.h"})[:m]
	tn := edakit.Shuffle([]string{"place.cpp", "route.cpp", "sta.cpp", "parse_def.cpp", "parse_lib.cpp", "cts.cpp", "gui.cpp", "drc.cpp", "main.cpp", "rsz.cpp"})[:k]
	// hinc[i] lists headers included by header i (always earlier ones); tinc for TUs.
	hinc := make([][]int, m)
	for i := 1; i < m; i++ {
		if edakit.Int(0, 2) > 0 {
			hinc[i] = append(hinc[i], edakit.Int(0, i-1))
		}
	}
	tinc := make([][]int, k)
	for t := 0; t < k; t++ {
		a := edakit.Int(0, m-1)
		tinc[t] = []int{a}
		if edakit.Int(0, 1) == 0 {
			b := edakit.Int(0, m-1)
			if b != a {
				tinc[t] = append(tinc[t], b)
			}
		}
	}
	// closure[i] is the set of headers reachable from header i, including i.
	var closure func(i int, acc map[int]bool)
	closure = func(i int, acc map[int]bool) {
		if acc[i] {
			return
		}
		acc[i] = true
		for _, j := range hinc[i] {
			closure(j, acc)
		}
	}
	affected := func(h int) []string {
		var out []string
		for t := 0; t < k; t++ {
			acc := map[int]bool{}
			for _, i := range tinc[t] {
				closure(i, acc)
			}
			if acc[h] {
				out = append(out, tn[t])
			}
		}
		return out
	}
	var desc []string
	for i := 0; i < m; i++ {
		for _, j := range hinc[i] {
			desc = append(desc, hn[i]+" includes "+hn[j])
		}
	}
	for t := 0; t < k; t++ {
		var hs []string
		for _, i := range tinc[t] {
			hs = append(hs, hn[i])
		}
		desc = append(desc, tn[t]+" includes "+strings.Join(hs, " and "))
	}
	q := edakit.Sprintf("In a tool's source tree, %s. ", strings.Join(desc, "; "))
	// Edit the header with the most dependents, but not always.
	best, bestN := 0, -1
	for h := 0; h < m; h++ {
		if n := len(affected(h)); n > bestN {
			best, bestN = h, n
		}
	}
	switch variant(ctx) {
	case 1:
		h := edakit.Int(0, m-1)
		a := affected(h)
		if len(a) == 0 && attempt < generator.MaxDrawAttempts && edakit.Int(0, 4) > 0 {
			return buildSystemsAttempt(ctx, attempt+1)
		}
		return edakit.Numeric(q+edakit.Sprintf("A developer edits %s. How many translation units must an incremental build recompile?", hn[h]), edakit.IntStr(len(a)),
			edakit.Sprintf("A translation unit depends on every header it includes directly or through other headers. The units that reach %s are %s: %d recompiles.", hn[h], names(a), len(a)),
			map[string]string{"recompiled": edakit.IntStr(len(a))})
	case 2:
		h := edakit.Int(0, m-1)
		a := affected(h)
		for tries := 0; len(a) == 0 && tries < generator.MaxDrawAttempts; tries++ {
			h = edakit.Int(0, m-1)
			a = affected(h)
		}
		if len(a) == 0 {
			h, a = best, affected(best)
		}
		p := edakit.Int(2, 4)
		c := edakit.Int(5, 40)
		link := edakit.Int(3, 20)
		waves := (len(a) + p - 1) / p
		total := waves*c + link
		return edakit.Numeric(q+edakit.Sprintf("Each translation unit takes %d s to compile and the final link takes %d s. The build runs %d compile jobs in parallel. After an edit to %s, how long does the incremental build take, in seconds?", c, link, p, hn[h]), edakit.IntStr(total),
			edakit.Sprintf("The edit dirties %s (%d units). They compile in ceil(%d/%d) = %d waves of %d s, then the link runs: %d × %d + %d = %d s.", names(a), len(a), len(a), p, waves, c, waves, c, link, total),
			map[string]string{"recompiled": edakit.IntStr(len(a)), "seconds": edakit.IntStr(total)})
	}
	return edakit.Numeric(q+"Which single header, if edited, forces the most recompilation? Give the number of translation units that edit recompiles.", edakit.IntStr(bestN),
		edakit.Sprintf("Count dependents per header through the transitive include graph. The worst is %s, reached by %s: %d units. Headers near the bottom of the include graph are the expensive ones to touch.", hn[best], names(affected(best)), bestN),
		map[string]string{"worst": edakit.IntStr(bestN)})
}

// ciPipelines asks for sharded wall time, runners needed and pass probability with retries.
func ciPipelines(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 12, 60)
	t := edakit.Int(2, 9)
	setup := edakit.Int(3, 12)
	switch variant(ctx) {
	case 1:
		r := edakit.Int(2, 8)
		w := setup + (n+r-1)/r*t
		q := edakit.Sprintf("A CI job for an open-source placer runs %d regression designs, each taking %d min, split evenly across %d runners. Each runner first spends %d min building the tool. Designs are handed out so no runner gets more than ceil(%d/%d) of them. What is the wall-clock time of the stage, in minutes?", n, t, r, setup, n, r)
		return edakit.Numeric(q, edakit.IntStr(w), edakit.Sprintf("Each runner builds for %d min, then runs at most ceil(%d/%d) = %d designs of %d min: %d + %d = %d min.", setup, n, r, (n+r-1)/r, t, setup, (n+r-1)/r*t, w),
			map[string]string{"wall": edakit.IntStr(w)})
	case 2:
		budget := setup + t*edakit.Int(2, 6)
		r := 1
		for setup+(n+r-1)/r*t > budget {
			r++
		}
		q := edakit.Sprintf("A pull-request check must finish in %d min. It runs %d regression designs of %d min each, spread evenly across identical runners, and every runner spends %d min building the tool first. What is the smallest number of runners that meets the budget?", budget, n, t, setup)
		return edakit.Numeric(q, edakit.IntStr(r), edakit.Sprintf("The designs per runner may be at most floor((%d - %d)/%d) = %d. So the runner count must be at least ceil(%d/%d) = %d; with %d runners the stage takes %d min.", budget, setup, t, (budget-setup)/t, n, (budget-setup)/t, r, r, setup+(n+r-1)/r*t),
			map[string]string{"runners": edakit.IntStr(r)})
	}
	f := edakit.Int(2, 8)
	qf := float64(edakit.Int(1, 10)) / 100
	pr := math.Pow(1-qf*qf, float64(f)) * 100
	ans := edakit.Num(pr, 2)
	q := edakit.Sprintf("A regression suite contains %d flaky tests, each of which fails independently with probability %s on a correct commit. CI reruns a failed flaky test once and counts it as passing if the rerun passes. What is the probability, as a percentage, that a correct commit passes CI? Round to 2 decimal places.", f, edakit.Num(qf, 2))
	return edakit.Numeric(q, ans, edakit.Sprintf("With one retry a flaky test fails only if both attempts fail: %s² = %s. All %d must pass: (1 - %s)^%d = %s%%. Without retries it would be %s%%.", edakit.Num(qf, 2), edakit.Num(qf*qf, 4), f, edakit.Num(qf*qf, 4), f, ans, edakit.Num(math.Pow(1-qf, float64(f))*100, 2)),
		map[string]string{"pass_percent": ans})
}

func qaCeilLog2(n int) int {
	k := 0
	for (1 << uint(k)) < n {
		k++
	}
	return k
}

// bisection asks for build counts and times to bisect a regression.
func bisection(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 10, 900)
	steps := qaCeilLog2(n)
	q := edakit.Sprintf("The nightly QoR run shows that worst slack on a reference design got 40 ps worse. The last good nightly build and the first bad one are separated by %d commits; the culprit is one of those %d commits (the bad nightly's commit included). ", n, n)
	switch variant(ctx) {
	case 1:
		return edakit.Numeric(q+"Each test builds the tool at one commit and runs the design. In the worst case, how many tests does bisection need to identify the culprit?", edakit.IntStr(steps),
			edakit.Sprintf("Each test halves the suspect range, so the count is ceil(log2 %d) = %d.", n, steps), map[string]string{"steps": edakit.IntStr(steps)})
	case 2:
		b := edakit.Int(5, 30)
		r := edakit.Int(10, 90)
		total := steps * (b + r)
		return edakit.Numeric(q+edakit.Sprintf("Each test takes %d min to build and %d min to run the design. In the worst case, how long does bisection take, in minutes?", b, r), edakit.IntStr(total),
			edakit.Sprintf("Bisection needs ceil(log2 %d) = %d tests of %d + %d min: %d min.", n, steps, b, r, total), map[string]string{"steps": edakit.IntStr(steps), "minutes": edakit.IntStr(total)})
	}
	return edakit.Numeric(q+"A linear scan tests commits one at a time from the oldest. In the worst case, how many more tests does the linear scan need than bisection?", edakit.IntStr(n-1-steps),
		edakit.Sprintf("A linear scan needs up to %d tests (the last suspect need not be tested once all others pass). Bisection needs ceil(log2 %d) = %d. Difference: %d.", n-1, n, steps, n-1-steps),
		map[string]string{"linear": edakit.IntStr(n - 1), "bisect": edakit.IntStr(steps)})
}

var cppRustItems = []edakit.Item{
	{KP: 1, Q: "Why are the hot kernels of most production EDA engines (timers, routers, placers) written in C++ rather than in Python?",
		A: "They need predictable native speed, control of memory layout and no interpreter overhead per object", D: []string{"Python cannot express graph algorithms", "C++ programs need no testing", "Python cannot read text files as large as Liberty"},
		Why: "Engines touch hundreds of millions of objects in tight loops. Native code with controlled layout is orders of magnitude faster than interpreted per-object work."},
	{KP: 1, Q: "A C++ timer library is distributed as a compiled shared library, and a plug-in built with a different compiler version crashes when it receives a std::string. What is the cause?",
		A: "The C++ ABI, including standard-library layouts, is not stable across compilers and versions", D: []string{"std::string cannot cross a function call", "Shared libraries cannot export functions", "The plug-in forgot to include the header"},
		Why: "Name mangling, class layout and standard-library internals may differ between toolchains. Stable plug-in boundaries use C ABI functions and plain data."},
	{KP: 1, Q: "A templated graph traversal in C++ is instantiated for each edge type. What is the main benefit compared with a virtual call per edge?",
		A: "The compiler can inline the per-edge code and optimize each instantiation", D: []string{"It uses less source code", "It makes the binary smaller", "It removes the need for a build system"},
		Why: "Templates resolve the call at compile time, so the per-edge work is inlined. A virtual call per edge blocks inlining and vectorization."},
	{KP: 1, Q: "Which recurring bug class in large C++ EDA code bases is a direct consequence of manual memory management?",
		A: "Use after free, for example a stale pointer to a deleted instance", D: []string{"Integer overflow in loop counters", "Slow compile times", "Incorrect operator precedence"},
		Why: "When an object is deleted, any raw pointer that still names it dangles. Using ids or handles that are checked on access is one mitigation."},
	{KP: 2, Q: "A Rust router shares a read-only routing grid among threads and gives each thread exclusive mutable access to its own region. What does the compiler check?",
		A: "That no data is mutated while another thread can read or write it", D: []string{"That the threads finish in a fixed order", "That the routing is DRC-clean", "That the program uses no heap memory"},
		Why: "Ownership and borrowing rules enforce shared-xor-mutable access at compile time, which rules out data races in safe code. They do not guarantee a deterministic order."},
	{KP: 2, Q: "A netlist graph has nets pointing to pins and pins pointing back to nets. Why is this awkward in safe Rust, and what do Rust EDA codes often do?",
		A: "Cyclic references conflict with single ownership, so they store objects in arenas and refer to them by index", D: []string{"Rust forbids graphs, so they use a C++ library instead", "Rust copies the whole graph on every access", "They use reference counting for every pin, which is free"},
		Why: "Single ownership makes cyclic pointer graphs hard. Arena storage with integer ids, the same pattern many C++ databases use, keeps the graph simple and cache friendly."},
	{KP: 2, Q: "Which guarantee does safe Rust NOT give a parallel EDA engine?",
		A: "That floating-point reductions produce the same bits on every run", D: []string{"That no two threads mutate the same data at once", "That references never outlive what they point to", "That a freed object cannot be accessed"},
		Why: "Memory safety and race freedom are enforced, but reduction order and other sources of nondeterminism are still the programmer's responsibility."},
	{KP: 2, Q: "A team rewrites a parser in Rust and keeps the timer in C++. Where do memory-safety guarantees stop?",
		A: "At the foreign-function boundary, where unsafe code must uphold the invariants by hand", D: []string{"Nowhere, because Rust checks the C++ code too", "At the first heap allocation", "At every function call inside the Rust code"},
		Why: "The Rust compiler cannot see into C++. Calls across the FFI are unsafe, and the wrapper must ensure lifetimes and aliasing rules hold."},
	{KP: 3, Q: "A Rust plug-in must be called from a C++ engine. What interface is the most robust to put between them?",
		A: "A C ABI with plain structs and explicit create and destroy functions", D: []string{"C++ classes with virtual functions shared directly", "Rust generic functions exported by name", "Passing std::vector and Vec values directly"},
		Why: "Both languages speak the C ABI. Plain data and explicit ownership functions avoid relying on either language's unstable layouts."},
	{KP: 3, Q: "OpenROAD exposes C++ engines to Tcl and Python through generated wrappers. What does a wrapper generator such as SWIG produce?",
		A: "Glue code that converts arguments and objects between the scripting language and C++", D: []string{"A complete rewrite of the engine in Python", "A new database format", "Machine code optimized for each script"},
		Why: "Wrapper generators read C++ declarations and emit bindings, so scripts call native functions and hold handles to native objects."},
	{KP: 3, Q: "A Python script holds a handle to a database instance, then a C++ engine deletes that instance. What can happen when the script uses the handle?",
		A: "The handle may dangle and crash or read garbage unless the binding checks validity", D: []string{"Python's garbage collector restores the instance", "The engine automatically re-creates the instance", "Nothing, because Python owns all database objects"},
		Why: "The binding wraps a native pointer that Python cannot track. Ownership across the boundary has to be designed, for example with ids checked on access."},
	{KP: 3, Q: "When is it reasonable to write a new EDA tool component in Python rather than C++ or Rust?",
		A: "When it orchestrates engines or post-processes reports, not when it runs a hot per-object loop", D: []string{"When it implements the inner loop of a detailed router", "When it must process a billion polygons per second", "Never, because Python cannot call native code"},
		Why: "Flow control, reporting and glue are cheap relative to the engines they call; inner loops over design objects belong in native code."},
}
