package dft

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// faultDetection: count the named stem faults a vector detects (KP 1), count every single
// stuck-at fault it detects (KP 2), count the vectors detecting one fault (KP 3).
func faultDetection(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	c := randCircuit(nIn, nG, 2, typesBasic)
	v := c.vec(edakit.Int(0, (1<<c.nIn)-1))
	order := strings.Join(c.names[:c.nIn], "")
	switch kp(ctx) {
	case 1:
		fs := edakit.Shuffle(c.stemFaults())[:edakit.Scaled(ctx, 4, 6)]
		var hit []string
		for _, f := range fs {
			if c.detects(v, f) {
				hit = append(hit, c.faultName(f))
			}
		}
		why := "None of them is detected."
		if len(hit) > 0 {
			why = "Detected: " + strings.Join(hit, "; ") + "."
		}
		good := c.sim(v, nil)
		return edakit.Numeric(c.describe()+edakit.Sprintf(" Apply %s. Of the faults %s, how many does this vector detect?", c.vecText(v), c.faultList(fs)),
			edakit.IntStr(len(hit)), "A fault is detected when its line carries the opposite of the stuck value and the difference reaches an output. "+why,
			map[string]string{"good_outputs": goodOuts(c, good)})
	case 2:
		n := 0
		for _, f := range c.faults() {
			if c.detects(v, f) {
				n++
			}
		}
		return edakit.Numeric(c.describe()+edakit.Sprintf(" Counting each fanout branch as its own line, the circuit has %d single stuck-at faults. How many of them does the vector %s = %s detect?", len(c.faults()), order, bits(v)),
			edakit.IntStr(n), edakit.Sprintf("Each line can only be detected stuck at the opposite of its fault-free value, so at most half the faults are candidates. Simulating each candidate, %d reach an output.", n),
			map[string]string{"detected": edakit.IntStr(n)})
	default:
		f, ok := c.testableStemFault(1, (1<<c.nIn)-1)
		if !ok {
			f = fault{node: c.outs[0], gate: -1, v: 0}
		}
		n := c.detectCount(f)
		return edakit.Numeric(c.describe()+edakit.Sprintf(" How many of the %d input vectors detect %s?", 1<<c.nIn, c.faultName(f)),
			edakit.IntStr(n), edakit.Sprintf("A detecting vector must set %s = %d and sensitize a path to an output. Exhaustive simulation finds %d such vectors.", c.names[f.node], 1-f.v, n),
			map[string]string{"detecting": edakit.IntStr(n)})
	}
}

func goodOuts(c *circuit, val []int) string {
	var p []string
	for _, o := range c.outs {
		p = append(p, edakit.Sprintf("%s=%d", c.names[o], val[o]))
	}
	return strings.Join(p, ", ")
}

// firstDetect draws a first-detection table: for each fault, the 1-based pattern that first
// detects it, or 0 if no pattern does.
func firstDetect(ctx generator.GeneratorContext) (nf, np int, first []int) {
	nf = edakit.Scaled(ctx, 5, 9)
	np = edakit.Scaled(ctx, 4, 8)
	first = make([]int, nf)
	for i := range first {
		if edakit.Int(0, 9) < 2 {
			first[i] = 0
		} else {
			first[i] = edakit.Int(1, np)
		}
	}
	return
}

func firstText(first []int) string {
	var p []string
	for i, x := range first {
		if x == 0 {
			p = append(p, edakit.Sprintf("f%d never", i+1))
		} else {
			p = append(p, edakit.Sprintf("f%d at pattern %d", i+1, x))
		}
	}
	return strings.Join(p, ", ")
}

// serialSim: runs with fault dropping (KP 1), percent saved (KP 2), coverage after k patterns
// (KP 3).
func serialSim(ctx generator.GeneratorContext) generator.Problem {
	nf, np, first := firstDetect(ctx)
	runs := 0
	for _, x := range first {
		if x == 0 {
			runs += np
		} else {
			runs += x
		}
	}
	stem := edakit.Sprintf("A serial fault simulator applies %d patterns to %d faults. One run simulates one faulty circuit for one pattern. The first pattern that detects each fault is: %s.", np, nf, firstText(first))
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+" With fault dropping (a fault is removed after its first detection), how many faulty-circuit runs are needed?",
			edakit.IntStr(runs), edakit.Sprintf("A detected fault costs runs up to and including its first detecting pattern; an undetected fault costs all %d. The sum is %d, against %d without dropping.", np, runs, nf*np),
			map[string]string{"runs": edakit.IntStr(runs)})
	case 2:
		s := 100 * float64(nf*np-runs) / float64(nf*np)
		ans := pct(s)
		return edakit.Numeric(stem+" What percentage of the faulty-circuit runs does fault dropping save compared with simulating every fault on every pattern? Round to 2 decimal places.",
			ans, edakit.Sprintf("Without dropping: %d runs. With dropping: %d. Saved: %d, which is %s percent.", nf*np, runs, nf*np-runs, ans),
			map[string]string{"runs": edakit.IntStr(runs), "full": edakit.IntStr(nf * np)})
	default:
		k := edakit.Int(2, np-1)
		d := 0
		for _, x := range first {
			if x != 0 && x <= k {
				d++
			}
		}
		ans := pct(100 * float64(d) / float64(nf))
		return edakit.Numeric(stem+edakit.Sprintf(" What is the fault coverage after the first %d patterns, in percent? Round to 2 decimal places.", k),
			ans, edakit.Sprintf("%d of the %d faults are first detected at or before pattern %d: %s percent.", d, nf, k, ans),
			map[string]string{"detected": edakit.IntStr(d)})
	}
}

// parallelSim: parallel-fault passes (KP 1), PPSFP passes (KP 2), bitwise word evaluation
// (KP 3).
func parallelSim(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		w := edakit.Pick([]int{16, 32, 64})
		f := edakit.Scaled(ctx, 100, 5000)
		p := edakit.Scaled(ctx, 10, 200)
		per := int(math.Ceil(float64(f) / float64(w-1)))
		ans := edakit.IntStr(per * p)
		return edakit.Numeric(edakit.Sprintf("A parallel-fault simulator packs machines into %d-bit words, one bit for the fault-free circuit and the rest for faulty circuits. With %d faults and %d patterns and no fault dropping, how many simulation passes over the circuit are needed?", w, f, p),
			ans, edakit.Sprintf("Each pass carries %d faults, so a pattern needs $\\lceil %d/%d \\rceil = %d$ passes; for %d patterns that is %s.", w-1, f, w-1, per, p, ans),
			map[string]string{"per_pattern": edakit.IntStr(per)})
	case 2:
		w := edakit.Pick([]int{32, 64})
		f := edakit.Scaled(ctx, 100, 5000)
		p := edakit.Scaled(ctx, 40, 2000)
		per := int(math.Ceil(float64(p) / float64(w)))
		ans := edakit.IntStr(per*f + per)
		return edakit.Numeric(edakit.Sprintf("A parallel-pattern single-fault simulator (PPSFP) packs %d patterns into each %d-bit word. With %d patterns and %d faults and no fault dropping, how many passes are needed, counting the fault-free passes too?", w, w, p, f),
			ans, edakit.Sprintf("Patterns form $\\lceil %d/%d \\rceil = %d$ words. Each word needs one fault-free pass and one pass per fault: $%d \\times (%d + 1) = %s$.", p, w, per, per, f, ans),
			map[string]string{"words": edakit.IntStr(per)})
	default:
		ops := []gtype{gAND, gNAND, gOR, gNOR, gXOR}
		t := edakit.Pick(ops)
		a, b := edakit.Int(0, 255), edakit.Int(0, 255)
		var r int
		switch t {
		case gAND:
			r = a & b
		case gNAND:
			r = ^(a & b) & 255
		case gOR:
			r = a | b
		case gNOR:
			r = ^(a | b) & 255
		case gXOR:
			r = a ^ b
		}
		bs := func(x int) string { return edakit.Sprintf("%08b", x) }
		return edakit.Numeric(edakit.Sprintf("In 8-bit parallel-pattern simulation, input A of %s carries the word %s and input B carries %s (bit $i$ is pattern $i$). What word does the output carry? Answer with the 8 bits, most significant first.", aGate(t), bs(a), bs(b)),
			bs(r), edakit.Sprintf("One bitwise %s instruction evaluates all 8 patterns at once: %s.", gname[t], bs(r)),
			map[string]string{"word": bs(r)})
	}
}

// deductive: fault list at a gate output, excluding the output's own local fault.
func deductive(ctx generator.GeneratorContext) generator.Problem {
	for attempt := 0; ; attempt++ {
		if attempt >= generator.MaxDrawAttempts {
			return edakit.Numeric("In deductive fault simulation, a 3-input AND gate has inputs a = 1 with $L_a = \\{1, 3\\}$, b = 1 with $L_b = \\{2\\}$, c = 0 with $L_c = \\{2, 4\\}$. Each list holds the faults (numbered 1 to 4) that complement that input under the current vector. Which faults are in the output list, ignoring the output's own local fault? Answer with the fault numbers separated by commas.",
				"4", "Only c is at the controlling value 0: intersect $L_c$ and remove $L_a \\cup L_b$, leaving {4}.", map[string]string{"output_list": "4"})
		}
		types := []gtype{gAND, gNAND, gOR, gNOR}
		if kp(ctx) == 3 {
			types = []gtype{gAND, gNAND, gOR, gNOR, gXOR}
		}
		t := edakit.Pick(types)
		ar := 2
		if t != gXOR {
			ar = edakit.Scaled(ctx, 2, 3)
		}
		pool := edakit.Scaled(ctx, 6, 9)
		vals := make([]int, ar)
		lists := make([]map[int]bool, ar)
		names := []string{"a", "b", "c"}
		var desc []string
		for i := 0; i < ar; i++ {
			vals[i] = edakit.Int(0, 1)
			lists[i] = map[int]bool{}
			for f := 1; f <= pool; f++ {
				if edakit.Int(0, 9) < 4 {
					lists[i][f] = true
				}
			}
			desc = append(desc, edakit.Sprintf("%s = %d with $L_%s = \\{%s\\}$", names[i], vals[i], names[i], setText(lists[i])))
		}
		res := map[int]bool{}
		var rule string
		if t == gXOR {
			for f := 1; f <= pool; f++ {
				cnt := 0
				for i := 0; i < ar; i++ {
					if lists[i][f] {
						cnt++
					}
				}
				if cnt%2 == 1 {
					res[f] = true
				}
			}
			rule = "For XOR, a fault flips the output when it flips an odd number of inputs, so keep faults in exactly one list."
		} else {
			cv, _ := controlling(t)
			var S []int
			for i := 0; i < ar; i++ {
				if vals[i] == cv {
					S = append(S, i)
				}
			}
			if len(S) == 0 {
				for i := 0; i < ar; i++ {
					for f := range lists[i] {
						res[f] = true
					}
				}
				rule = "No input is at the controlling value, so any fault that flips any input flips the output: take the union."
			} else {
				for f := 1; f <= pool; f++ {
					in := true
					for _, i := range S {
						if !lists[i][f] {
							in = false
						}
					}
					for i := 0; i < ar; i++ {
						isS := false
						for _, j := range S {
							if j == i {
								isS = true
							}
						}
						if !isS && lists[i][f] {
							in = false
						}
					}
					if in {
						res[f] = true
					}
				}
				rule = edakit.Sprintf("Inputs at the controlling value %d: the output flips only if a fault flips all of them and none of the others, so intersect the controlling inputs' lists and remove the other lists.", cv)
			}
		}
		if len(res) == 0 {
			continue
		}
		ans := setAnswer(res)
		q := edakit.Sprintf("In deductive fault simulation, a %d-input %s gate has inputs %s. Each list holds the faults (numbered 1 to %d) that complement that input under the current vector. Which faults are in the output list, ignoring the output's own local fault? Answer with the fault numbers separated by commas.", ar, gname[t], strings.Join(desc, ", "), pool)
		return edakit.Numeric(q, ans, rule+" Result: {"+ans+"}.", map[string]string{"output_list": ans})
	}
}

func setText(m map[int]bool) string {
	var k []int
	for x := range m {
		k = append(k, x)
	}
	sort.Ints(k)
	var p []string
	for _, x := range k {
		p = append(p, edakit.IntStr(x))
	}
	return strings.Join(p, ", ")
}

func setAnswer(m map[int]bool) string {
	s := setText(m)
	return strings.ReplaceAll(s, " ", "")
}

// concurrentCalc (KP 2): count the faulty gate copies whose output diverges from the good gate.
func concurrentCalc(ctx generator.GeneratorContext) generator.Problem {
	t := edakit.Pick([]gtype{gAND, gNAND, gOR, gNOR})
	ar := edakit.Scaled(ctx, 2, 3)
	names := []string{"a", "b", "c"}
	good := make([]int, ar)
	for i := range good {
		good[i] = edakit.Int(0, 1)
	}
	gout := eval2(t, good)
	nf := edakit.Scaled(ctx, 4, 7)
	var desc []string
	div := 0
	var divs []string
	for k := 0; k < nf; k++ {
		v := append([]int(nil), good...)
		v[edakit.Int(0, ar-1)] ^= 1
		if edakit.Int(0, 2) == 0 {
			v[edakit.Int(0, ar-1)] ^= 1
		}
		same := true
		for i := range v {
			if v[i] != good[i] {
				same = false
			}
		}
		if same {
			v[0] = 1 - v[0]
		}
		var p []string
		for i := 0; i < ar; i++ {
			p = append(p, edakit.Sprintf("%s=%d", names[i], v[i]))
		}
		desc = append(desc, edakit.Sprintf("f%d (%s)", k+1, strings.Join(p, ", ")))
		if eval2(t, v) != gout {
			div++
			divs = append(divs, edakit.Sprintf("f%d", k+1))
		}
	}
	var gp []string
	for i := 0; i < ar; i++ {
		gp = append(gp, edakit.Sprintf("%s=%d", names[i], good[i]))
	}
	why := edakit.Sprintf("The good gate outputs %d. ", gout)
	if div > 0 {
		why += "Diverging copies: " + strings.Join(divs, ", ") + "; the rest produce the good output and stop propagating."
	} else {
		why += "Every faulty copy produces the good output, so nothing propagates."
	}
	return edakit.Choice(edakit.Sprintf("A concurrent fault simulator holds, at a %d-input %s gate, the good gate with inputs %s and faulty copies %s. How many faulty copies produce an output different from the good gate, and so must be scheduled at the gate's fanout?", ar, gname[t], strings.Join(gp, ", "), strings.Join(desc, ", ")),
		edakit.IntStr(div), why, edakit.Distinct(div, 3, nf-div, nf, div+1)...)
}

// coverage: fault coverage (KP 1), test coverage (KP 2), ATPG efficiency (KP 3).
func coverage(ctx generator.GeneratorContext) generator.Problem {
	total := edakit.Scaled(ctx, 2000, 90000)
	u := total * edakit.Int(5, 40) / 1000
	ab := total * edakit.Int(1, 15) / 1000
	nd := total * edakit.Int(1, 20) / 1000
	d := total - u - ab - nd
	stem := edakit.Sprintf("After ATPG on %d collapsed stuck-at faults: %d detected, %d proven untestable (redundant or blocked by tied logic), %d aborted at the backtrack limit, and %d not detected for other reasons.", total, d, u, ab, nd)
	facts := map[string]string{"detected": edakit.IntStr(d), "untestable": edakit.IntStr(u), "total": edakit.IntStr(total)}
	switch kp(ctx) {
	case 1:
		ans := pct(100 * float64(d) / float64(total))
		return edakit.Numeric(stem+" What is the fault coverage, detected faults over all faults, in percent? Round to 2 decimal places.",
			ans, edakit.Sprintf("FC = %d/%d = %s percent.", d, total, ans), facts)
	case 2:
		ans := pct(100 * float64(d) / float64(total-u))
		return edakit.Numeric(stem+" What is the test coverage, which removes the untestable faults from the denominator, in percent? Round to 2 decimal places.",
			ans, edakit.Sprintf("TC = %d/(%d - %d) = %d/%d = %s percent; it is higher than fault coverage because untestable faults cannot be detected by any test.", d, total, u, d, total-u, ans), facts)
	default:
		ans := pct(100 * float64(d+u) / float64(total))
		return edakit.Numeric(stem+" What is the ATPG efficiency, the share of faults the tool resolved (detected or proven untestable), in percent? Round to 2 decimal places.",
			ans, edakit.Sprintf("Efficiency = (%d + %d)/%d = %s percent. The aborted faults are the gap a higher backtrack limit or a SAT engine might close.", d, u, total, ans), facts)
	}
}

var concurrentItems = []edakit.Item{
	{KP: 1, Q: "What does a concurrent fault simulator store at each gate?",
		A: "The good gate's values plus a list of faulty copies that differ from the good gate at that gate", D: []string{"A complete copy of the circuit for every fault", "Only the faults that have already been detected", "A bit-packed word of all patterns"},
		Why: "Concurrent simulation keeps a faulty-gate record only where a fault's values diverge from the good machine, so storage follows activity, not circuit size times faults."},
	{KP: 1, Q: "Why is concurrent fault simulation built on event-driven simulation?",
		A: "Faulty machines that agree with the good machine generate no events, so only divergent activity is evaluated", D: []string{"Concurrent simulation needs a levelized netlist with no events", "Event-driven simulation is required to pack bits into words", "It only simulates primary outputs"},
		Why: "The approach exploits the fact that most faulty machines are identical to the good one almost everywhere."},
	{KP: 1, Q: "A faulty copy at a gate converges: its inputs and output now match the good gate. What does the simulator do?",
		A: "It deletes that faulty record, since the fault's effect has disappeared there", D: []string{"It keeps the record and marks the fault detected", "It copies the record to every fanout gate", "It restarts the simulation of that fault"},
		Why: "Convergence means the fault is invisible at that point; keeping the record would waste memory and evaluations."},
	{KP: 1, Q: "Compared with deductive fault simulation, what does concurrent simulation handle more naturally?",
		A: "Sequential circuits, multi-valued logic and timing, because each faulty gate is evaluated explicitly", D: []string{"Only combinational circuits with two-valued logic", "Faults that are never activated", "Bit-parallel patterns"},
		Why: "Deductive simulation manipulates fault-list set algebra that becomes awkward with X values and delays; concurrent simulation just evaluates faulty gate copies."},
	{KP: 1, Q: "What is the main cost of concurrent fault simulation?",
		A: "Memory: the fault lists at busy gates can grow large and change dynamically", D: []string{"It needs one pass per fault", "It cannot drop detected faults", "It needs SAT solving per pattern"},
		Why: "Storage per gate grows with the number of divergent faults, and list management dominates run time in large designs."},
	{KP: 1, Q: "A local fault at a gate is, for example, its own output stuck-at-0. When is its faulty record created?",
		A: "At the start, at the fault site, because the fault is always present there", D: []string{"Only when the fault is detected", "Only when an input event reaches the gate", "Never: local faults are handled by the good machine"},
		Why: "Each fault begins as a record at its own site and spreads to fanout gates only when it causes divergence."},
	{KP: 3, Q: "Fault f7's record reaches a primary output with a value that differs from the good output. What happens next with fault dropping on?",
		A: "f7 is marked detected and its records are removed everywhere", D: []string{"f7 keeps propagating until the end of the pattern set", "f7 is reclassified as untestable", "The good machine is resimulated without f7"},
		Why: "Dropping detected faults shrinks the lists, which is the main speedup over the pattern set."},
	{KP: 3, Q: "Most commercial fault simulators for large combinational logic under scan use PPSFP instead of concurrent simulation. Why?",
		A: "With full scan the problem is combinational, and bit-parallel evaluation of many patterns per fault is faster and uses little memory", D: []string{"Concurrent simulation cannot handle stuck-at faults", "PPSFP models timing exactly", "Concurrent simulation needs no fault list"},
		Why: "Full scan removes state, so a compiled, bit-parallel engine with fault dropping wins; concurrent simulation shines on sequential and timing-aware problems."},
	{KP: 3, Q: "Two faults produce identical values at a gate. Could the simulator share one record for both?",
		A: "Not in the classic algorithm: each record belongs to one fault, so their effects can later diverge", D: []string{"Yes, always, because their values match", "Only if both faults are stuck-at-0", "Only at primary outputs"},
		Why: "Classic concurrent simulation keeps records per fault; the two faults may behave differently further downstream or on later patterns."},
	{KP: 3, Q: "Why does the number of faulty records tend to fall as simulation proceeds through a large pattern set?",
		A: "Fault dropping removes detected faults, and most faults are detected early", D: []string{"Records are deleted after a fixed time", "Gates become faster", "The circuit is collapsed further during simulation"},
		Why: "Coverage curves rise steeply: the easy faults go in the first patterns, leaving few, hard faults with small lists."},
	{KP: 3, Q: "A faulty record at a flip-flop in a sequential circuit diverges from the good flip-flop. What does that imply for the next clock cycle?",
		A: "The fault's state differs, so its effect can reappear in later cycles even if no output differs now", D: []string{"The fault is detected immediately", "The record is deleted at the clock edge", "The fault becomes a transition fault"},
		Why: "Divergent state carries the fault's effect forward in time, which is why sequential fault simulation must keep the record."},
	{KP: 3, Q: "Why is concurrent simulation said to simulate only the differences?",
		A: "The good machine is simulated once, and each faulty machine is represented only where it differs from it", D: []string{"It skips simulating the good machine", "It simulates only one fault at a time", "It simulates only the primary outputs"},
		Why: "Representing faulty machines as deltas over the good machine is the central idea, shared with incremental timing and incremental simulation."},
}
