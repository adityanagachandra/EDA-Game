package logic

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerFSM(reg *generator.Registry) {
	reg.Register("logic.fsm.moore_mealy", retryTies(mooreMealy))
	reg.Register("logic.fsm.state_diagram", retryTies(stateDiagram))
	reg.Register("logic.fsm.encoding", retryTies(stateEncoding))
	reg.Register("logic.fsm.minimization", retryTies(minimization))
	reg.Register("logic.fsm.next_state_logic", retryTies(nextStateLogic))
	reg.Register("logic.fsm.reset_states", &edakit.Bank{Items: resetItems})
}

const stateLetters = "ABCDEFGH"

func sname(i int) string { return string(stateLetters[i]) }

// machine is a Moore machine (out per state) or a Mealy machine (mout per state and input) with
// a one-bit input.
type machine struct {
	n    int
	next [][2]int
	out  []int
	mout [][2]int
}

func randMoore(n int) machine {
	m := machine{n: n, next: make([][2]int, n), out: make([]int, n)}
	for {
		ones := 0
		for i := 0; i < n; i++ {
			m.out[i] = edakit.Int(0, 1)
			ones += m.out[i]
			m.next[i] = [2]int{edakit.Int(0, n-1), edakit.Int(0, n-1)}
		}
		if ones > 0 && ones < n {
			return m
		}
	}
}

func (m machine) mooreTable() string {
	var parts []string
	for i := 0; i < m.n; i++ {
		parts = append(parts, edakit.Sprintf("%s (output %d): on 0 go to %s, on 1 go to %s", sname(i), m.out[i], sname(m.next[i][0]), sname(m.next[i][1])))
	}
	return strings.Join(parts, "; ")
}

func randBits(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = edakit.Int(0, 1)
	}
	return out
}

func bitString(bs []int) string {
	var b strings.Builder
	for _, x := range bs {
		b.WriteByte(byte('0' + x))
	}
	return b.String()
}

func countOverlaps(stream, pat string) int {
	c := 0
	for i := 0; i+len(pat) <= len(stream); i++ {
		if stream[i:i+len(pat)] == pat {
			c++
		}
	}
	return c
}

func mooreMealy(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		m := randMoore(edakit.Scaled(ctx, 3, 4))
		in := randBits(edakit.Scaled(ctx, 6, 10))
		s, ones := 0, 0
		var visits []string
		for _, x := range in {
			s = m.next[s][x]
			ones += m.out[s]
			visits = append(visits, sname(s))
		}
		return edakit.Numeric("A Moore machine has states "+m.mooreTable()+". It starts in A and receives the inputs "+joinInts(in, ", ")+", one per clock edge. Counting the output after each edge, how many of those outputs are 1?",
			edakit.IntStr(ones), "The states visited after each edge are "+strings.Join(visits, ", ")+edakit.Sprintf("; %d of them output 1.", ones), map[string]string{"ones": edakit.IntStr(ones)})
	case 2:
		n := edakit.Scaled(ctx, 2, 3)
		var parts []string
		next := make([][2]int, n)
		mout := make([][2]int, n)
		for i := 0; i < n; i++ {
			for x := 0; x < 2; x++ {
				next[i][x] = edakit.Int(0, n-1)
				mout[i][x] = edakit.Int(0, 1)
			}
			parts = append(parts, edakit.Sprintf("from %s, input 0 goes to %s with output %d and input 1 goes to %s with output %d", sname(i), sname(next[i][0]), mout[i][0], sname(next[i][1]), mout[i][1]))
		}
		in := randBits(edakit.Scaled(ctx, 6, 10))
		s, ones := 0, 0
		var outs []int
		for _, x := range in {
			ones += mout[s][x]
			outs = append(outs, mout[s][x])
			s = next[s][x]
		}
		return edakit.Numeric("A Mealy machine has transitions: "+strings.Join(parts, "; ")+". It starts in A and receives the inputs "+joinInts(in, ", ")+". How many of the outputs produced on those transitions are 1?",
			edakit.IntStr(ones), "Each output belongs to the transition taken; the outputs are "+joinInts(outs, ", ")+".", map[string]string{"ones": edakit.IntStr(ones)})
	default:
		k := edakit.Scaled(ctx, 3, 4)
		pat := bitString(randBits(k))
		L := edakit.Scaled(ctx, 10, 16)
		stream := []byte(bitString(randBits(L)))
		pos := edakit.Int(0, L-k)
		copy(stream[pos:], pat)
		if edakit.Int(0, 1) == 0 && pos+2*k-1 <= L {
			// plant an overlapping second copy when the pattern allows one
			for sh := 1; sh < k; sh++ {
				if pat[sh:] == pat[:k-sh] && pos+sh+k <= L {
					copy(stream[pos+sh:], pat)
					break
				}
			}
		}
		c := countOverlaps(string(stream), pat)
		return edakit.Numeric(edakit.Sprintf("A Moore detector asserts its output whenever the last %d inputs equal %s, with overlapping occurrences allowed. Over the input stream %s, how many times is the output asserted?", k, pat, string(stream)), edakit.IntStr(c),
			edakit.Sprintf("Sliding a window of length %d along the stream finds %s, counting ones that share bits.", k, plural(c, "match")), map[string]string{"matches": edakit.IntStr(c)})
	}
}

// kmpNext returns the longest suffix of s that is a prefix of p, capped below len(p)+1.
func kmpNext(s, p string) int {
	for l := min(len(s), len(p)); l > 0; l-- {
		if s[len(s)-l:] == p[:l] {
			return l
		}
	}
	return 0
}

func stateDiagram(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 3, 7)
	pat := bitString(randBits(k))
	switch variant(ctx) {
	case 1:
		switch edakit.Int(0, 2) {
		case 0:
			return edakit.Numeric(edakit.Sprintf("What is the minimum number of states of a Moore machine that detects the pattern %s (overlaps allowed) on a one-bit input?", pat), edakit.IntStr(k+1),
				edakit.Sprintf("States 0 to %d record the matched prefix length; state %d outputs 1. That is %d states, all distinguishable.", k, k, k+1), map[string]string{"k": edakit.IntStr(k)})
		case 1:
			return edakit.Numeric(edakit.Sprintf("What is the minimum number of states of a Mealy machine that detects the pattern %s (overlaps allowed) on a one-bit input?", pat), edakit.IntStr(k),
				edakit.Sprintf("The Mealy machine emits 1 on the transition that completes the pattern, so it needs states 0 to %d only: %d states.", k-1, k), map[string]string{"k": edakit.IntStr(k)})
		default:
			return edakit.Numeric(edakit.Sprintf("A minimal Moore detector for the pattern %s is binary-encoded. How many flip-flops does it need?", pat), edakit.IntStr(clog2(k+1)),
				edakit.Sprintf("It has %d states, so $\\lceil \\log_2 %d \\rceil = %d$ flip-flops.", k+1, k+1, clog2(k+1)), map[string]string{"states": edakit.IntStr(k + 1)})
		}
	case 2:
		i := edakit.Int(1, k-1)
		x := byte('1')
		if pat[i] == '1' {
			x = '0'
		}
		nx := kmpNext(pat[:i]+string(x), pat)
		return edakit.Numeric(edakit.Sprintf("In a detector for the pattern %s, state i means the longest suffix of the input so far that is a prefix of the pattern has length i. The machine is in state %d (matched %s) and the input is %c. What is the next state?", pat, i, pat[:i], x), edakit.IntStr(nx),
			edakit.Sprintf("The string %s%c has longest suffix that is a prefix of %s of length %d.", pat[:i], x, pat, nx), map[string]string{"next": edakit.IntStr(nx)})
	default:
		x := byte('0' + edakit.Int(0, 1))
		if edakit.Int(0, 3) == 0 {
			nx := 0
			if pat[0] == x {
				nx = 1
			}
			return edakit.Numeric(edakit.Sprintf("A Moore detector for the pattern %s does not allow overlapping matches: after a detection it starts afresh. State i records the matched prefix length. It is in state %d (just detected) and the input is %c. What is the next state?", pat, k, x), edakit.IntStr(nx),
				edakit.Sprintf("Without overlap the input starts a new match only if it equals the first pattern bit, %c.", pat[0]), map[string]string{"next": edakit.IntStr(nx)})
		}
		nx := kmpNext(pat[1:]+string(x), pat)
		return edakit.Numeric(edakit.Sprintf("A Moore detector for the pattern %s allows overlapping matches. State i records the longest suffix of the input that is a prefix of the pattern. It is in state %d (just detected) and the input is %c. What is the next state?", pat, k, x), edakit.IntStr(nx),
			edakit.Sprintf("The string %s%c has longest proper suffix that is a prefix of the pattern of length %d.", pat, x, nx), map[string]string{"next": edakit.IntStr(nx)})
	}
}

func stateEncoding(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		S := edakit.Scaled(ctx, 5, 40)
		switch edakit.Int(0, 3) {
		case 0:
			return edakit.Numeric(edakit.Sprintf("How many flip-flops does a binary-encoded state machine with %d states need?", S), edakit.IntStr(clog2(S)),
				edakit.Sprintf("$\\lceil \\log_2 %d \\rceil = %d$.", S, clog2(S)), map[string]string{"states": edakit.IntStr(S)})
		case 1:
			return edakit.Numeric(edakit.Sprintf("How many flip-flops does a one-hot state machine with %d states need?", S), edakit.IntStr(S),
				"One-hot uses one flip-flop per state.", map[string]string{"states": edakit.IntStr(S)})
		case 2:
			return edakit.Numeric(edakit.Sprintf("How many flip-flops does a Gray-encoded state machine with %d states need?", S), edakit.IntStr(clog2(S)),
				edakit.Sprintf("Gray codes are binary-width codes: $\\lceil \\log_2 %d \\rceil = %d$.", S, clog2(S)), map[string]string{"states": edakit.IntStr(S)})
		default:
			return edakit.Numeric(edakit.Sprintf("A cyclic sequencer with %d states is implemented as a Johnson (twisted-ring) counter. How many flip-flops does it need?", S), edakit.IntStr((S+1)/2),
				edakit.Sprintf("An n-flop Johnson counter has 2n states, so $\\lceil %d/2 \\rceil = %d$.", S, (S+1)/2), map[string]string{"states": edakit.IntStr(S)})
		}
	case 2:
		if edakit.Int(0, 1) == 0 {
			S := edakit.Int(5, 60)
			b := clog2(S)
			return edakit.Numeric(edakit.Sprintf("A state machine with %d states uses minimum-width binary encoding. How many codes are unused?", S), edakit.IntStr((1<<uint(b))-S),
				edakit.Sprintf("%d flip-flops give %d codes; %d are used.", b, 1<<uint(b), S), map[string]string{"flops": edakit.IntStr(b)})
		}
		S := edakit.Int(4, 12)
		return edakit.Numeric(edakit.Sprintf("A state machine with %d states uses one-hot encoding. How many bit patterns of its state register are illegal (not exactly one bit set)?", S), edakit.IntStr((1<<uint(S))-S),
			edakit.Sprintf("$2^{%d} - %d = %d$.", S, S, (1<<uint(S))-S), map[string]string{"flops": edakit.IntStr(S)})
	default:
		S := edakit.Scaled(ctx, 4, 12)
		kind := edakit.Int(0, 2)
		code := func(i int) uint64 {
			if kind == 1 {
				return uint64(i ^ (i >> 1))
			}
			return uint64(i)
		}
		tog := 0
		for i := 0; i < S; i++ {
			tog += popcount(code(i) ^ code((i+1)%S))
		}
		name := []string{"plain binary (state i gets code i)", "Gray code (state i gets the Gray code of i)", "one-hot"}[kind]
		if kind == 2 {
			tog = 2 * S
		}
		return edakit.Numeric(edakit.Sprintf("A state machine steps cyclically through states 0, 1, ..., %d and back to 0. With %s encoding, how many flip-flop toggles occur in one full pass of %d transitions, including the wrap to 0?", S-1, name, S), edakit.IntStr(tog),
			edakit.Sprintf("Sum the Hamming distances between consecutive codes, including the last-to-first wrap: %d toggles.", tog), map[string]string{"toggles": edakit.IntStr(tog)})
	}
}

// refine returns the number of equivalence classes of a Moore machine by partition refinement.
func refine(m machine) int {
	block := make([]int, m.n)
	copy(block, m.out)
	for {
		sig := map[[3]int]int{}
		nb := make([]int, m.n)
		for i := 0; i < m.n; i++ {
			key := [3]int{block[i], block[m.next[i][0]], block[m.next[i][1]]}
			id, ok := sig[key]
			if !ok {
				id = len(sig)
				sig[key] = id
			}
			nb[i] = id
		}
		old := map[int]bool{}
		for _, b := range block {
			old[b] = true
		}
		if len(sig) == len(old) {
			return len(sig)
		}
		block = nb
	}
}

func minimization(ctx generator.GeneratorContext) generator.Problem {
	v := variant(ctx)
	base := edakit.Int(2, 4)
	dups := edakit.Scaled(ctx, 2, 3)
	if v == 1 {
		dups = 2
		base = edakit.Int(2, 3)
	}
	n := base + dups
	src := make([]int, n) // which base state each state copies
	for i := 0; i < n; i++ {
		if i < base {
			src[i] = i
		} else {
			src[i] = edakit.Int(0, base-1)
		}
	}
	b := randMoore(base)
	m := machine{n: n, next: make([][2]int, n), out: make([]int, n)}
	for i := 0; i < n; i++ {
		m.out[i] = b.out[src[i]]
		for x := 0; x < 2; x++ {
			t := b.next[src[i]][x]
			var copies []int
			for j := 0; j < n; j++ {
				if src[j] == t {
					copies = append(copies, j)
				}
			}
			m.next[i][x] = edakit.Pick(copies)
		}
	}
	// relabel randomly
	perm := edakit.Shuffle(func() []int {
		p := make([]int, n)
		for i := range p {
			p[i] = i
		}
		return p
	}())
	r := machine{n: n, next: make([][2]int, n), out: make([]int, n)}
	for i := 0; i < n; i++ {
		r.out[perm[i]] = m.out[i]
		r.next[perm[i]] = [2]int{perm[m.next[i][0]], perm[m.next[i][1]]}
	}
	cls := refine(r)
	if v == 3 {
		return edakit.Numeric("A completely specified Moore machine has states "+r.mooreTable()+". After merging all equivalent states, how many states are removed?", edakit.IntStr(n-cls),
			edakit.Sprintf("Partition refinement starting from the output classes stabilizes with %d blocks, so %d of the %d states are removed.", cls, n-cls, n), map[string]string{"classes": edakit.IntStr(cls)})
	}
	return edakit.Numeric("A completely specified Moore machine has states "+r.mooreTable()+". How many states does the equivalent minimal machine have?", edakit.IntStr(cls),
		edakit.Sprintf("Split by output, then repeatedly split blocks whose members' successors fall in different blocks. The partition stabilizes with %d blocks.", cls), map[string]string{"classes": edakit.IntStr(cls)})
}

func nextStateLogic(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		S := edakit.Int(3, 40)
		in := edakit.Int(1, 3)
		b := clog2(S)
		rows := 1 << uint(b+in)
		return edakit.Numeric(edakit.Sprintf("A binary-encoded state machine has %d states and %d input bits. How many rows does its full next-state truth table have?", S, in), edakit.IntStr(rows),
			edakit.Sprintf("It has %d state bits, so the table has $2^{%d + %d} = %d$ rows.", b, b, in, rows), map[string]string{"state_bits": edakit.IntStr(b)})
	case 2:
		m := randMoore(4)
		bitIdx := edakit.Int(0, 1)
		cnt := 0
		var rows []string
		for s := 0; s < 4; s++ {
			for x := 0; x < 2; x++ {
				nb := (m.next[s][x] >> uint(bitIdx)) & 1
				cnt += nb
				if nb == 1 {
					rows = append(rows, bin(uint64(s<<1|x), 3))
				}
			}
		}
		var parts []string
		for s := 0; s < 4; s++ {
			parts = append(parts, edakit.Sprintf("%s: (%s, %s)", sname(s), sname(m.next[s][0]), sname(m.next[s][1])))
		}
		return edakit.Numeric(edakit.Sprintf("A four-state machine is encoded A = 00, B = 01, C = 10, D = 11 as $s_1 s_0$, with one input $x$. Its next states on $x = 0$ and $x = 1$ are %s. Over the eight rows of $(s_1, s_0, x)$, how many minterms does the next-state function $s_%d^+$ have?", strings.Join(parts, ", "), bitIdx), edakit.IntStr(cnt),
			edakit.Sprintf("Reading bit %d of each next-state code, the rows with $s_%d^+ = 1$ are %s.", bitIdx, bitIdx, strings.Join(rows, ", ")), map[string]string{"minterms": edakit.IntStr(cnt)})
	default:
		n := edakit.Scaled(ctx, 4, 5)
		m := randMoore(n)
		t := edakit.Int(0, n-1)
		preds := map[int]bool{}
		for s := 0; s < n; s++ {
			for x := 0; x < 2; x++ {
				if m.next[s][x] == t {
					preds[s] = true
				}
			}
		}
		if len(preds) == 0 {
			m.next[edakit.Int(0, n-1)][edakit.Int(0, 1)] = t
			for s := 0; s < n; s++ {
				for x := 0; x < 2; x++ {
					if m.next[s][x] == t {
						preds[s] = true
					}
				}
			}
		}
		var parts []string
		for s := 0; s < n; s++ {
			parts = append(parts, edakit.Sprintf("%s: (%s, %s)", sname(s), sname(m.next[s][0]), sname(m.next[s][1])))
		}
		return edakit.Numeric(edakit.Sprintf("A one-hot state machine with one input $x$ has next states on $x = 0$ and $x = 1$ of %s. Its next-state equation for state bit %s is written as a sum of products, one term per source state, with terms from the same source state merged. How many product terms does it have?", strings.Join(parts, ", "), sname(t)), edakit.IntStr(len(preds)),
			edakit.Sprintf("Each distinct state with an arc into %s contributes one term (just the state bit when both inputs lead there): %d terms.", sname(t), len(preds)), map[string]string{"terms": edakit.IntStr(len(preds))})
	}
}

var resetItems = []edakit.Item{
	{KP: 1, Q: "At power-up, before reset, what values do the flip-flops of an FSM's state register hold?",
		A: "Arbitrary values, possibly an illegal code", D: []string{"All zeros, guaranteed by the process", "The reset state, because synthesis initializes it", "The last state before power was removed"},
		Why: "Flip-flops power up in unpredictable states, which is why control state must be reset explicitly."},
	{KP: 1, Q: "A designer leaves some datapath registers without reset but resets every FSM state register. Is this acceptable?",
		A: "Yes, if the datapath registers are always written before they are read", D: []string{"No, every flip-flop in a design must be reset", "Yes, but only if the clock is gated during reset", "No, because unreset registers violate setup"},
		Why: "Skipping reset on registers that are written before use saves area and routing; control state must still be initialized."},
	{KP: 1, Q: "How does a synchronous reset differ from an asynchronous reset in an FSM?",
		A: "It acts only at a clock edge, as part of the next-state logic", D: []string{"It acts immediately, through dedicated flip-flop pins", "It needs no clock running to take effect", "It cannot be used with one-hot encoding"},
		Why: "A synchronous reset is a mux in front of D that forces the reset code on the next edge; an asynchronous reset bypasses the clock."},
	{KP: 1, Q: "Formal verification of an FSM starts from its reset state. What happens if one flip-flop is left out of reset?",
		A: "The set of initial states doubles, and properties relying on initialization may fail", D: []string{"Nothing, because formal tools ignore reset", "The flip-flop is removed by the tool", "The FSM becomes combinational"},
		Why: "An unreset bit can start at either value, so the tool must consider both initial states."},
	{KP: 2, Q: "How is the set of reachable states of an explicit FSM computed?",
		A: "A graph search from the reset state over the transition relation", D: []string{"Counting the states listed in the specification", "Enumerating all codes of the state register", "Minimizing the next-state logic"},
		Why: "Reachability is breadth-first or depth-first search from reset; symbolic model checkers do the same with BDDs or SAT."},
	{KP: 2, Q: "A binary-encoded FSM with 5 states has 3 flip-flops. How does synthesis usually treat the 3 unused codes?",
		A: "As don't-cares in the next-state logic, unless a safe FSM option is set", D: []string{"As extra states that must each go to reset", "As errors that stop synthesis", "As states that must hold their value"},
		Why: "Codes that are never reached have unspecified behavior, which lets the minimizer simplify the logic."},
	{KP: 2, Q: "A sequential optimization merges logic using the fact that some state codes are unreachable. What must the equivalence checker be told to prove it correct?",
		A: "The reset state, so it can compare behavior from reset only", D: []string{"Nothing: combinational equivalence suffices", "The clock period", "The list of illegal codes to treat as reachable"},
		Why: "The optimized design is equivalent only on reachable states, so the proof must start from reset."},
	{KP: 2, Q: "A state appears in the state diagram, but no transition enters it and it is not the reset state. What is it?",
		A: "Unreachable, so its behavior can never be observed in normal operation", D: []string{"Illegal, so synthesis must add recovery logic for it", "The default state that every illegal code goes to", "Equivalent to the reset state"},
		Why: "Without an incoming arc from a reachable state, search from reset never reaches it."},
	{KP: 3, Q: "A one-hot FSM suffers a single-event upset that sets a second bit. With unused codes treated as don't-cares, what can happen?",
		A: "The machine may lock up in illegal states and never return to a legal one", D: []string{"It always returns to reset on the next edge", "It automatically clears the extra bit", "It continues exactly as if nothing happened"},
		Why: "Don't-care treatment leaves the behavior of illegal codes arbitrary, so it may cycle among them forever."},
	{KP: 3, Q: "What does a safe FSM option in synthesis add?",
		A: "Logic that sends every illegal code to a known recovery state", D: []string{"A second copy of the state register for voting", "A watchdog timer that resets the chip", "Gray encoding of all states"},
		Why: "Safe FSMs give every unused code a defined transition, typically to reset or an error state."},
	{KP: 3, Q: "Which property would a formal tool check to confirm safe recovery?",
		A: "From any state code, a legal state is reached within a bounded number of cycles", D: []string{"The reset state is reachable from itself", "No two states have the same output", "The next-state logic has no don't-cares"},
		Why: "Bounded recovery from every code is the property that distinguishes a safe FSM."},
	{KP: 3, Q: "Why do safety-critical designs sometimes accept larger FSM logic?",
		A: "Explicit recovery from illegal states protects against soft errors and faults", D: []string{"Larger logic is always faster", "Smaller logic cannot be tested", "Don't-care minimization is illegal in safety standards"},
		Why: "The extra logic buys robustness: an upset cannot strand the controller in an undefined state."},
}
