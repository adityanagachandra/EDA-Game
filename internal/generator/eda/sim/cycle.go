package sim

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func levels(n netlist) []int {
	lv := make([]int, n.nIn+len(n.gates))
	for i, g := range n.gates {
		l := lv[g.a]
		if lv[g.b] > l {
			l = lv[g.b]
		}
		lv[n.nIn+i] = l + 1
	}
	return lv
}

func levelization(ctx generator.GeneratorContext) generator.Problem {
	n := randomNetlist(edakit.Int(3, 4), edakit.Scaled(ctx, 5, 9))
	// Name the last input as a register output, to show that flop outputs are level 0.
	n.names[n.nIn-1] = "q"
	lv := levels(n)
	maxL := 0
	for _, l := range lv {
		if l > maxL {
			maxL = l
		}
	}
	setup := edakit.Sprintf("A combinational block reads primary inputs %s and the flip-flop output q. Its gates are %s. Inputs and flop outputs are level 0, and a gate's level is one more than the highest level among its inputs.",
		strings.Join(n.names[:n.nIn-1], ", "), n.describe())
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	facts := map[string]string{"levels": edakit.IntStr(maxL)}
	switch kp {
	case 1:
		gi := edakit.Int(len(n.gates)/2, len(n.gates)-1)
		g := n.gates[gi]
		return edakit.Numeric(setup+edakit.Sprintf(" What is the level of %s?", n.names[n.nIn+gi]),
			edakit.IntStr(lv[n.nIn+gi]),
			edakit.Sprintf("%s reads %s (level %d) and %s (level %d), so its level is %d.", n.names[n.nIn+gi], n.names[g.a], lv[g.a], n.names[g.b], lv[g.b], lv[n.nIn+gi]), facts)
	case 2:
		return edakit.Numeric(setup+" How many levels of gates does the block have, that is, what is the highest gate level?",
			edakit.IntStr(maxL),
			edakit.Sprintf("Propagating levels in topological order, the deepest gate sits at level %d: the longest path from a level-0 source is %d gates.", maxL, maxL), facts)
	default:
		target := edakit.Int(1, maxL)
		cnt := 0
		for i := range n.gates {
			if lv[n.nIn+i] == target {
				cnt++
			}
		}
		return edakit.Numeric(setup+edakit.Sprintf(" A levelized simulator evaluates all gates of one level before the next. How many gates does it evaluate at level %d?", target),
			edakit.IntStr(cnt),
			edakit.Sprintf("Computing every gate's level and counting those equal to %d gives %d. Gates of one level have no dependences on each other, so they can be evaluated in any order.", target, cnt), facts)
	}
}

func cycleCost(ctx generator.GeneratorContext) generator.Problem {
	ce := edakit.Pick([]int{40, 50, 60, 80, 100})
	cc := edakit.Pick([]int{2, 4, 5, 8, 10})
	if edakit.Int(0, 1) == 0 {
		be := 100 * float64(cc) / float64(ce)
		ans := edakit.Num(be, 1) + "%"
		q := edakit.Sprintf("An event-driven simulator spends %d ns per gate evaluation, including queue management. A levelized cycle-based simulator of the same design spends %d ns per gate but evaluates every gate every cycle. Below what activity (the fraction of gates evaluated per cycle by the event-driven simulator) is event-driven simulation faster?", ce, cc)
		why := edakit.Sprintf("Per cycle, event-driven costs a N x %d ns and cycle-based costs N x %d ns. They are equal at a = %d / %d = %s.", ce, cc, cc, ce, ans)
		return edakit.Choice(q, ans, why,
			edakit.Num(100*float64(ce)/float64(cc), 1)+"%",
			edakit.Num(100*(1-float64(cc)/float64(ce)), 1)+"%",
			edakit.Num(100*float64(cc)/float64(ce+cc), 1)+"%")
	}
	a := edakit.Pick([]int{2, 5, 10, 15, 20, 30})
	r := float64(a) / 100 * float64(ce) / float64(cc)
	ans := edakit.Num(r, 2)
	q := edakit.Sprintf("A design has 200,000 gates. Each cycle the event-driven simulator evaluates %d%% of them at %d ns per evaluation, while a compiled cycle-based simulator evaluates all of them at %d ns each. What is the ratio of event-driven time per cycle to cycle-based time per cycle?", a, ce, cc)
	why := edakit.Sprintf("Event-driven: 0.%02d x 200,000 x %d ns. Cycle-based: 200,000 x %d ns. The ratio is %s; above 1 the cycle-based simulator wins.", a, ce, cc, ans)
	return edakit.Choice(q, ans, why,
		edakit.Num(1/r, 2),
		edakit.Num(float64(ce)/float64(cc), 2),
		edakit.Num(float64(a)/100*float64(cc)/float64(ce), 2))
}

var compiledItems = []edakit.Item{
	{KP: 1, Q: "What does a cycle-based simulator assume that an event-driven simulator does not?",
		A: "Only values at clock edges matter, so each gate is evaluated once per cycle in levelized order", D: []string{"Every gate has a unit delay", "Signals can take only the values X and Z", "The design has no registers"},
		Why: "Cycle-based simulation drops intra-cycle timing for synchronous logic, which is what lets it skip the event queue."},
	{KP: 1, Q: "A design has a glitch on a combinational net between clock edges. What does a cycle-based simulator show?",
		A: "Nothing: it computes only the settled value at the edge", D: []string{"The glitch, at its exact time", "An X for the whole cycle", "A delta-cycle warning"},
		Why: "With one evaluation per gate per cycle there is no intermediate value to show."},
	{KP: 1, Q: "Why must the design be levelized before cycle-based evaluation?",
		A: "So that every gate is evaluated after all of its inputs in the same cycle", D: []string{"So that the event queue stays sorted by time", "So that X values are removed", "So that registers can be removed"},
		Why: "A topological (levelized) order guarantees one evaluation per gate gives the settled value."},
	{KP: 1, Q: "Which design style breaks a pure cycle-based simulator?",
		A: "Logic that depends on delays, such as a pulse generator built from gate delays", D: []string{"A pipeline of registers on one clock", "A wide adder", "A finite-state machine with synchronous reset"},
		Why: "Cycle-based simulation has no notion of gate delay, so delay-dependent behavior cannot be represented."},
	{KP: 1, Q: "A design has two asynchronous clock domains. How does a cycle-based simulator typically handle it?",
		A: "It evaluates each domain's logic on that domain's edges, scheduling edges from both clocks in time order", D: []string{"It merges both clocks into one", "It cannot simulate registers at all", "It falls back to analog simulation"},
		Why: "Edges are still scheduled in time; between edges each domain's combinational logic is evaluated in levelized order."},
	{KP: 1, Q: "Why do cycle-based simulators usually run two-state (0 and 1) values?",
		A: "Two-state values map onto machine words and bitwise instructions, which is much faster", D: []string{"Four-state values are illegal in synthesizable RTL", "Two-state values are more accurate", "X values cannot be levelized"},
		Why: "A 64-bit word holds 64 two-state signals; four-state needs two bits per signal and extra logic per operation."},
	{KP: 3, Q: "What does a compiled simulator generate from the netlist, compared with an interpreted one?",
		A: "Straight-line host code with one statement per operation, so evaluation needs no data-structure walk", D: []string{"A table of events sorted by time", "An analog matrix to factor each cycle", "A BDD of the whole design"},
		Why: "Compilation turns the levelized netlist into code the host CPU runs directly, removing interpretation overhead."},
	{KP: 3, Q: "A compiled cycle-based simulator rebuilds its executable after every small RTL change. What cost does this add?",
		A: "Compile time, which can dominate for short tests", D: []string{"Lower accuracy for long tests", "Loss of all coverage data", "A need to rerun synthesis"},
		Why: "Compiled simulation trades build time for run speed; for many short runs the build can cost more than the runs."},
	{KP: 3, Q: "A design runs at 2% activity per cycle. A team switches from event-driven to compiled cycle-based simulation and sees a slowdown. What is the likely reason?",
		A: "At low activity the event-driven simulator skips almost all gates, while cycle-based evaluates all of them", D: []string{"Compiled code is always slower", "The design has too few registers", "Cycle-based simulation adds delta cycles"},
		Why: "Cycle-based cost is proportional to design size; event-driven cost is proportional to activity."},
	{KP: 3, Q: "Which check can a cycle-based simulator not perform that a timing-annotated event-driven simulator can?",
		A: "Setup and hold timing checks against SDF delays", D: []string{"Comparing outputs with a reference model", "Collecting functional coverage", "Evaluating assertions on clock edges"},
		Why: "Timing checks need actual signal arrival times, which a cycle-based model does not compute."},
	{KP: 3, Q: "A combinational loop exists through a latch-free mux path that is never enabled at the same time. What does a levelizing compiler need to do?",
		A: "Break the apparent loop, for example by evaluating that logic iteratively until it settles", D: []string{"Ignore it, since levelization tolerates cycles", "Insert a register into the RTL", "Convert the loop to four-state values"},
		Why: "A topological order does not exist for a cyclic graph, so the compiler must split it or iterate to convergence."},
	{KP: 3, Q: "Why are compiled simulators favored for running software on a processor model for billions of cycles?",
		A: "Their per-cycle cost is low and predictable, and long runs amortize the compile time", D: []string{"They model gate delays more accurately", "They need no testbench", "They detect X propagation better"},
		Why: "Throughput dominates long runs; compiled cycle-based code maximizes cycles per second."},
}

var verilatorItems = []edakit.Item{
	{KP: 1, Q: "What does Verilator produce from synthesizable SystemVerilog?",
		A: "C++ (or SystemC) classes that model the design, compiled with a host compiler", D: []string{"A netlist of FPGA primitives", "A bytecode file run by an interpreter", "A SPICE deck"},
		Why: "Verilator is a compiler: it translates the RTL into C++ code that becomes part of an executable."},
	{KP: 1, Q: "In a C++ harness around a Verilated model, what must the harness do after changing the clock input?",
		A: "Call the model's eval() so it recomputes its state and outputs", D: []string{"Nothing: the model updates on assignment", "Call the simulator's event queue", "Restart the model"},
		Why: "The model only computes when eval() is called; the harness drives inputs, calls eval() and reads outputs."},
	{KP: 1, Q: "Which value system does a Verilated model use for ordinary signals?",
		A: "Two-state: X and Z are mapped to 0 or 1", D: []string{"Full four-state with X and Z", "Nine-valued IEEE 1164 logic", "Real-valued voltages"},
		Why: "Verilator is a two-state simulator; X has to be resolved to a concrete value at compile or run time."},
	{KP: 1, Q: "A team wants a testbench written in SystemVerilog with delays and event controls to run on Verilator. What does it need?",
		A: "Timing support enabled with the --timing option, which uses C++ coroutines", D: []string{"Nothing: Verilator always supported every testbench construct", "A separate interpreter for the testbench", "It is impossible in any version"},
		Why: "Recent Verilator releases support timing constructs when built with --timing; otherwise the harness must be C++."},
	{KP: 2, Q: "Verilator offers --x-assign and --x-initial. What do they control?",
		A: "What concrete value replaces explicit X assignments and uninitialized variables", D: []string{"Whether X is printed in waveforms", "How many threads the model uses", "The delay of X transitions"},
		Why: "Because the model is two-state, these options pick fixed, random or unique values where X would appear."},
	{KP: 2, Q: "Why randomize initial values (for example --x-initial unique with a runtime seed) in a two-state simulator?",
		A: "To expose bugs where the design depends on uninitialized state, which a fixed 0 would hide", D: []string{"To make every run identical", "To speed up evaluation", "To emulate gate delays"},
		Why: "If every flop starts at 0 a missing reset can go unnoticed; random initial values make such dependences show up across seeds."},
	{KP: 2, Q: "A design passes in Verilator but fails in a four-state simulator with X on an output after reset. What is the most likely cause?",
		A: "A register without reset that Verilator initialized to a benign value", D: []string{"A Verilator compiler bug", "A missing timescale", "The four-state simulator is too optimistic"},
		Why: "Four-state simulation shows the uninitialized state as X; two-state simulation silently picks a value."},
	{KP: 2, Q: "Two-state simulation cannot show X-propagation. What complementary step catches reset and initialization bugs?",
		A: "Running reset sequences in a four-state simulator or with randomized initial values", D: []string{"Turning off lint warnings", "Increasing the thread count", "Using --x-assign 0 everywhere"},
		Why: "Either keep X (four-state) or randomize what X becomes; both reveal dependence on unknown state."},
	{KP: 3, Q: "How does Verilator run a model on several threads?",
		A: "It partitions the evaluation graph into macro-tasks and schedules them across threads", D: []string{"It runs one thread per always block with locks", "It forks one process per clock domain", "It cannot use threads"},
		Why: "Verilator's multithreading groups logic into coarse tasks to amortize synchronization costs."},
	{KP: 3, Q: "Verilator warns UNOPTFLAT on a signal. What does the warning mean?",
		A: "The signal appears to be in a combinational loop after flattening, which forces slower iterative evaluation", D: []string{"The signal is unused", "The signal has X values", "The signal is driven by two clocks"},
		Why: "Verilator found an apparent circular dependence that it cannot order statically; splitting the signal usually removes it."},
	{KP: 3, Q: "Why is Verilator's lint mode useful even when another simulator is used for sign-off?",
		A: "It applies strict static checks such as width mismatches and unused or undriven signals", D: []string{"It proves the design correct", "It measures timing slack", "It generates coverage"},
		Why: "Lint catches many RTL mistakes without simulation; Verilator is a fast, strict linter."},
	{KP: 3, Q: "A Verilated model is much faster than an interpreted event-driven simulator on the same design. What is the main reason?",
		A: "It is statically scheduled, compiled, two-state code with no runtime event queue for most logic", D: []string{"It skips the clock", "It simulates only the testbench", "It uses analog approximations"},
		Why: "Static ordering and compilation remove queue overhead and use machine words for signals."},
}

// parallelSim draws a task DAG and asks for work, span, or a speedup bound.
func parallelSim(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 9)
	cost := make([]int, n)
	preds := make([][]int, n)
	var lines []string
	for i := 0; i < n; i++ {
		cost[i] = edakit.Int(1, 9)
		if i >= 2 {
			k := edakit.Int(1, 2)
			seen := map[int]bool{}
			for j := 0; j < k; j++ {
				p := edakit.Int(0, i-1)
				if !seen[p] {
					seen[p] = true
					preds[i] = append(preds[i], p)
				}
			}
		} else if i == 1 && edakit.Int(0, 1) == 0 {
			preds[i] = []int{0}
		}
		desc := edakit.Sprintf("T%d (%d us)", i+1, cost[i])
		if len(preds[i]) > 0 {
			var ps []string
			for _, p := range preds[i] {
				ps = append(ps, edakit.Sprintf("T%d", p+1))
			}
			desc += " after " + strings.Join(ps, " and ")
		}
		lines = append(lines, desc)
	}
	fin := make([]int, n)
	W, S := 0, 0
	for i := 0; i < n; i++ {
		st := 0
		for _, p := range preds[i] {
			if fin[p] > st {
				st = fin[p]
			}
		}
		fin[i] = st + cost[i]
		W += cost[i]
		if fin[i] > S {
			S = fin[i]
		}
	}
	setup := "One simulated cycle of a partitioned design is a task graph: " + strings.Join(lines, "; ") + "."
	facts := map[string]string{"work": edakit.IntStr(W), "span": edakit.IntStr(S)}
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		return edakit.Numeric(setup+" What is the total work W of one cycle, in microseconds?", edakit.IntStr(W),
			edakit.Sprintf("Work is the sum of all task costs, %d us: the time one thread needs.", W), facts)
	case 2:
		return edakit.Numeric(setup+" What is the span S (the critical-path time) of one cycle, in microseconds?", edakit.IntStr(S),
			edakit.Sprintf("Taking each task's finish time as its cost plus the latest finish among its predecessors, the last task finishes at %d us. No number of threads can beat that.", S), facts)
	default:
		P := edakit.Pick([]int{2, 4, 8})
		if edakit.Int(0, 1) == 0 {
			b := float64(W)/float64(P) + float64(S)
			ans := edakit.Num(b, 2)
			return edakit.Numeric(setup+edakit.Sprintf(" By Brent's bound, a greedy scheduler on %d threads finishes a cycle within at most W/P + S microseconds. What is that bound? Round to 2 decimal places.", P), ans,
				edakit.Sprintf("W = %d us and S = %d us, so W/P + S = %d/%d + %d = %s us.", W, S, W, P, S, ans), facts)
		}
		sp := float64(W) / float64(S)
		if float64(P) < sp {
			sp = float64(P)
		}
		ans := edakit.Num(sp, 2)
		return edakit.Numeric(setup+edakit.Sprintf(" With %d threads and no synchronization cost, what is the best possible speedup over one thread, min(P, W/S)? Round to 2 decimal places.", P), ans,
			edakit.Sprintf("W = %d us and S = %d us, so the parallelism W/S is %s; with P = %d the bound is %s.", W, S, edakit.Num(float64(W)/float64(S), 2), P, ans), facts)
	}
}
