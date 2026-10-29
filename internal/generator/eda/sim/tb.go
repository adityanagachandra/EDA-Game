package sim

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func verificationGap(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		if edakit.Int(0, 1) == 0 {
			n := edakit.Scaled(ctx, 8, 20)
			s := 1 << n
			return edakit.Numeric(edakit.Sprintf("A controller has %d flip-flops and no other state. How many distinct states can its register contents take?", n),
				edakit.IntStr(s), edakit.Sprintf("Each flip-flop holds one of two values, so there are 2^%d = %d states, whether or not all are reachable.", n, s),
				map[string]string{"states": edakit.IntStr(s)})
		}
		w1, w2 := edakit.Int(3, 8), edakit.Int(3, 8)
		c := edakit.Int(0, 2)
		total := 1 << (w1 + w2 + c)
		q := edakit.Sprintf("A combinational block has a %d-bit input A, a %d-bit input B and %d single-bit control inputs. How many distinct input vectors does it have?", w1, w2, c)
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("The input space is the product of the field sizes: 2^%d x 2^%d x 2^%d = 2^%d = %d.", w1, w2, c, w1+w2+c, total),
			map[string]string{"vectors": edakit.IntStr(total)})
	case 2:
		for {
			k := edakit.Int(34, 46)
			rateExp := edakit.Int(5, 7)
			rate := math.Pow(10, float64(rateExp))
			hours := math.Pow(2, float64(k)) / rate / 3600
			if hours < 1 || hours > 2e6 {
				continue
			}
			ans := edakit.Num(hours, 0)
			q := edakit.Sprintf("A block has %d input bits and no state. A simulator applies 10^%d input vectors per second. How many hours does exhaustive simulation of every input vector take? Round to the nearest hour.", k, rateExp)
			return edakit.Numeric(q, ans,
				edakit.Sprintf("2^%d vectors at 10^%d per second take %s seconds, which is %s hours. Exhaustive simulation stops scaling after a few dozen input bits.", k, rateExp, edakit.Num(math.Pow(2, float64(k))/rate, 0), ans),
				map[string]string{"hours": ans})
		}
	default:
		k := edakit.Int(12, 20)
		n := edakit.Int(1, 50) * 100
		pct := 100 * float64(n) / math.Pow(2, float64(k))
		ans := edakit.Num(pct, 3)
		q := edakit.Sprintf("A regression applies %d distinct random vectors to a block with %d input bits. What percentage of its input space has been exercised? Round to 3 decimal places.", n, k)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The space has 2^%d = %d vectors, so %d distinct vectors cover %s%% of it.", k, 1<<k, n, ans),
			map[string]string{"space": edakit.IntStr(1 << k)})
	}
}

func intsText(xs []int) string {
	s := make([]string, len(xs))
	for i, x := range xs {
		s[i] = edakit.IntStr(x)
	}
	return strings.Join(s, ", ")
}

func scoreboard(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	n := edakit.Scaled(ctx, 6, 10)
	exp := make([]int, n)
	for i := range exp {
		exp[i] = edakit.Int(10, 99)
	}
	switch kp {
	case 1:
		act := append([]int(nil), exp...)
		bad := edakit.Int(1, 3)
		for j := 0; j < bad; j++ {
			i := edakit.Int(0, n-1)
			act[i] = edakit.Int(10, 99)
		}
		mism := 0
		for i := range exp {
			if exp[i] != act[i] {
				mism++
			}
		}
		q := edakit.Sprintf("An in-order scoreboard pops one expected transaction from its queue for every DUT output and compares them. The reference model predicted %s; the DUT produced %s, in that order. How many comparisons fail?", intsText(exp), intsText(act))
		return edakit.Numeric(q, edakit.IntStr(mism),
			edakit.Sprintf("Comparing position by position, %d position(s) differ.", mism),
			map[string]string{"mismatches": edakit.IntStr(mism)})
	case 2:
		act := edakit.Shuffle(exp)
		bad := edakit.Int(1, 3)
		for j := 0; j < bad; j++ {
			act[edakit.Int(0, n-1)] = edakit.Int(10, 99)
		}
		// Multiset difference: actual items with no matching expected item.
		pool := map[int]int{}
		for _, e := range exp {
			pool[e]++
		}
		unmatched := 0
		for _, a := range act {
			if pool[a] > 0 {
				pool[a]--
			} else {
				unmatched++
			}
		}
		q := edakit.Sprintf("A DUT may complete requests out of order, so the scoreboard matches each output against any outstanding expected transaction with the same value and removes the match. Expected: %s. Actual: %s. How many actual outputs find no match?", intsText(exp), intsText(act))
		return edakit.Numeric(q, edakit.IntStr(unmatched),
			edakit.Sprintf("Treating both lists as multisets and removing matched pairs leaves %d unmatched outputs. An in-order comparison would also flag outputs that are merely reordered.", unmatched),
			map[string]string{"unmatched": edakit.IntStr(unmatched)})
	default:
		starts := make([]int, n)
		ends := make([]int, n)
		t := 0
		for i := 0; i < n; i++ {
			t += edakit.Int(1, 3)
			starts[i] = t
			e := t + edakit.Int(2, 8)
			if i > 0 && e <= ends[i-1] {
				e = ends[i-1] + 1
			}
			ends[i] = e
		}
		best := 0
		for c := 0; c <= ends[n-1]; c++ {
			cnt := 0
			for i := 0; i < n; i++ {
				if starts[i] <= c && c < ends[i] {
					cnt++
				}
			}
			if cnt > best {
				best = cnt
			}
		}
		q := edakit.Sprintf("A monitor on the DUT input pushes an expected transaction into the scoreboard queue at cycles %s. The output monitor pops them, in order, at cycles %s. A transaction occupies the queue from its push cycle up to, but not including, its pop cycle. What is the largest number of transactions in the queue at once?", intsText(starts), intsText(ends))
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("Counting, for each cycle, the transactions pushed but not yet popped, the maximum is %d. The queue must hold that many, and a DUT latency bug shows up as this number growing without bound.", best),
			map[string]string{"max_outstanding": edakit.IntStr(best)})
	}
}

var testbenchItems = []edakit.Item{
	{KP: 1, Q: "In a layered testbench, which component converts a transaction into pin-level activity on the DUT interface?",
		A: "The driver", D: []string{"The monitor", "The scoreboard", "The coverage collector"},
		Why: "The driver turns abstract transactions into signal wiggles; the monitor does the reverse."},
	{KP: 1, Q: "Which component observes DUT pins and reconstructs transactions without driving anything?",
		A: "The monitor", D: []string{"The driver", "The sequencer", "The stimulus generator"},
		Why: "A monitor is passive: it samples the interface and publishes transactions to checkers and coverage."},
	{KP: 1, Q: "Where should the comparison of DUT outputs with expected results live, so that it is reused across tests?",
		A: "In a scoreboard fed by monitors and a reference model", D: []string{"Inside each directed test", "Inside the driver", "In the DUT's RTL"},
		Why: "Centralizing checking in the environment means every test, directed or random, is checked the same way."},
	{KP: 1, Q: "A testbench drives stimulus and dumps waveforms, and an engineer inspects them by eye. What is it missing?",
		A: "Self-checking: automatic comparison against expected behavior", D: []string{"A clock generator", "A DUT instance", "Waveform compression"},
		Why: "Without automatic checks, a regression cannot tell passing runs from failing ones."},
	{KP: 2, Q: "A team writes one directed test per feature from the specification. What is the main weakness of this approach as the design grows?",
		A: "It only finds bugs the test writer anticipated", D: []string{"Directed tests cannot check outputs", "Directed tests cannot be rerun", "Directed tests need formal tools"},
		Why: "Directed tests cover the scenarios someone thought of; random stimulus reaches combinations nobody wrote down."},
	{KP: 2, Q: "Why do testbenches work with transactions rather than individual signals at their upper layers?",
		A: "Stimulus, checking and coverage can be written once in terms of operations, independent of the pin protocol", D: []string{"Transactions simulate faster than RTL", "Signals cannot be randomized", "Transactions are synthesizable"},
		Why: "Abstraction separates what is being tested from how the interface encodes it."},
	{KP: 2, Q: "A constrained-random testbench reports no failures after 10,000 seeds. What else must be true before the block is considered verified?",
		A: "Coverage shows the intended scenarios were actually exercised", D: []string{"The seeds were all even", "The waveforms were viewed", "The simulator ran in four-state mode"},
		Why: "Random tests passing says nothing about what they reached; coverage measures that."},
	{KP: 2, Q: "Which item belongs in the reference model rather than in the testbench's checker?",
		A: "The computation of what the DUT should output for each input transaction", D: []string{"The comparison of expected and actual values", "The pin timing of the bus protocol", "The coverage bins"},
		Why: "The reference model predicts; the checker or scoreboard compares."},
	{KP: 3, Q: "A testbench samples a DUT output exactly at the clock edge in the same region where the DUT updates it. What symptom appears?",
		A: "Results flip between old and new values depending on simulator ordering", D: []string{"The output is always X", "The simulator stops with a fatal error", "Coverage drops to zero"},
		Why: "Sampling and updating in the same region is a race; sample before the edge or use clocking blocks."},
	{KP: 3, Q: "What does a clocking block's input skew specify?",
		A: "How long before the clock edge the testbench samples DUT outputs", D: []string{"How long after the edge the DUT updates", "The clock period", "The skew between two clock domains in the DUT"},
		Why: "Input skew moves the sampling point before the edge so the testbench reads settled values."},
	{KP: 3, Q: "Why should a driver apply new inputs after the active clock edge rather than at it?",
		A: "So the DUT samples stable inputs at the edge and never races with the driver", D: []string{"So that hold time is violated deliberately", "So that the simulation runs faster", "So that coverage is collected"},
		Why: "Driving after the edge mirrors real flop timing and removes ordering dependence."},
	{KP: 3, Q: "A monitor misses every other transaction on a bus that can complete back-to-back. What is the likely bug?",
		A: "The monitor waits an extra clock between transactions instead of sampling every cycle", D: []string{"The scoreboard is too small", "The driver uses nonblocking assignments", "The DUT has a reset bug"},
		Why: "A monitor must follow the protocol cycle by cycle; waiting extra cycles drops back-to-back transfers."},
}

var cocotbItems = []edakit.Item{
	{KP: 1, Q: "How does a cocotb test talk to the simulator running the HDL?",
		A: "Through the simulator's procedural interface, such as VPI for Verilog or VHPI for VHDL", D: []string{"By editing the HDL source before compilation", "By parsing the waveform file after the run", "Through a network socket to the FPGA"},
		Why: "cocotb loads into the simulator as a library and reads and writes signals through its programming interface."},
	{KP: 1, Q: "In cocotb, what does await RisingEdge(dut.clk) do?",
		A: "Suspends the coroutine until the simulator reports a rising edge on clk", D: []string{"Drives a rising edge onto clk", "Blocks the whole simulator until Python returns", "Waits one second of wall-clock time"},
		Why: "Triggers suspend Python coroutines; the simulator advances and resumes them when the event happens."},
	{KP: 1, Q: "Why are cocotb tests written as async functions?",
		A: "Each test and helper is a coroutine that yields control to the simulator while waiting for time or events", D: []string{"Python requires async for any file I/O", "Async functions run on several CPU cores", "The HDL is compiled into Python"},
		Why: "Cooperative coroutines let many testbench threads interleave with simulated time."},
	{KP: 1, Q: "A cocotb test uses await Timer(10, 'ns'). What happens?",
		A: "The coroutine resumes after 10 ns of simulated time", D: []string{"The coroutine resumes after 10 ns of real time", "The DUT clock period becomes 10 ns", "Every signal is delayed by 10 ns"},
		Why: "Timer waits on simulation time, not wall-clock time."},
	{KP: 2, Q: "Which open-source simulator compiles Verilog into an intermediate form that an event-driven runtime executes?",
		A: "Icarus Verilog, with its vvp runtime", D: []string{"Verilator", "GHDL", "Yosys"},
		Why: "Icarus compiles to vvp code and runs it event-driven, with four-state values and delays."},
	{KP: 2, Q: "A team must simulate a VHDL design with an open-source tool under cocotb. Which simulator fits?",
		A: "GHDL", D: []string{"Icarus Verilog", "Verilator", "nextpnr"},
		Why: "GHDL is an open-source VHDL simulator; Icarus and Verilator target Verilog and SystemVerilog."},
	{KP: 2, Q: "A cocotb regression runs too slowly on Icarus for a large design. Which switch is most likely to speed up the simulation itself?",
		A: "Running the same tests on Verilator, a compiled two-state simulator", D: []string{"Running on GHDL", "Adding more waveform dumping", "Switching to a four-state simulator with more delays"},
		Why: "Verilator compiles the design to C++ and usually runs much faster, at the cost of two-state values and some constructs."},
	{KP: 2, Q: "Why might a design that simulates cleanly in Icarus fail to build in Verilator?",
		A: "Verilator accepts mainly synthesizable constructs and is stricter about some code", D: []string{"Verilator only supports VHDL", "Icarus is two-state", "Verilator needs a waveform file"},
		Why: "Verilator compiles a cycle-oriented model and rejects or warns on constructs it cannot schedule statically."},
	{KP: 3, Q: "Right after await RisingEdge(dut.clk), a test reads a register output and sees the value from before the edge. Why?",
		A: "The nonblocking update of that register has not happened yet in the current time step", D: []string{"cocotb always reads the previous cycle", "The register is broken", "VPI cannot read registers"},
		Why: "The coroutine resumes at the edge, before the NBA region; awaiting ReadOnly() reads settled values."},
	{KP: 3, Q: "What does awaiting ReadOnly() guarantee in cocotb?",
		A: "The test resumes after the time step's values have settled, but it may not write signals then", D: []string{"The test can write signals with no race", "Simulation time advances by one unit", "All coroutines are cancelled"},
		Why: "ReadOnly corresponds to the end of the time step, where reads are stable and writes are not allowed."},
	{KP: 3, Q: "A cocotb testbench does heavy Python computation on every clock cycle. What is the cost?",
		A: "Each simulator-to-Python transition and Python step adds overhead that can dominate run time", D: []string{"The DUT behavior changes", "Coverage is lost", "No cost: Python runs in parallel with the simulator"},
		Why: "The simulator waits while Python runs; moving work to transaction level reduces crossings."},
	{KP: 3, Q: "Why do teams use a Python reference model in a cocotb scoreboard?",
		A: "The model can be written quickly with standard libraries and compared transaction by transaction", D: []string{"Python models are synthesizable", "Python models run faster than RTL", "VPI requires a Python model"},
		Why: "cocotb's appeal is the Python ecosystem; a golden model in Python is easy to write and check against."},
}
