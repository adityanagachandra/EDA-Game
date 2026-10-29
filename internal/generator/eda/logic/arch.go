package logic

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerArch(reg *generator.Registry) {
	reg.Register("logic.arch.datapath_control", &edakit.Bank{Items: datapathItems})
	reg.Register("logic.arch.performance_equation", retryTies(performanceEquation))
	reg.Register("logic.arch.pipelining", retryTies(pipelining))
	reg.Register("logic.arch.hazards", retryTies(hazards))
	reg.Register("logic.arch.branch_prediction", retryTies(branchPrediction))
	reg.Register("logic.arch.cache_organization", retryTies(cacheOrganization))
	reg.Register("logic.arch.cache_performance", retryTies(cachePerformance))
}

func performanceEquation(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		ic := 10 * edakit.Int(5, 200)
		cpi := float64(edakit.Int(10, 30)) / 10
		f := edakit.Pick([]float64{1, 1.5, 2, 2.5, 3, 4})
		ms := float64(ic) * cpi / f
		return withExact(edakit.Numeric(edakit.Sprintf("A program executes %d million instructions with an average CPI of %s on a %s GHz clock. How long does it run, in milliseconds? Round to 2 decimal places.", ic, edakit.Num(cpi, 1), edakit.Num(f, 1)), edakit.Num(ms, 2),
			edakit.Sprintf("$T = IC \\cdot CPI / f = %d \\times 10^6 \\cdot %s / (%s \\times 10^9)$ s $= %s$ ms.", ic, edakit.Num(cpi, 1), edakit.Num(f, 1), edakit.Num(ms, 2)), map[string]string{"cycles_millions": edakit.Num(float64(ic)*cpi, 1)}), ms)
	case 2:
		k := edakit.Int(3, 4)
		var pct []int
		left := 100
		for i := 0; i < k-1; i++ {
			p := 5 * edakit.Int(2, (left/5)-2*(k-1-i))
			if p > left-10 {
				p = left - 10
			}
			pct = append(pct, p)
			left -= p
		}
		pct = append(pct, left)
		names := []string{"ALU", "load/store", "branch", "multiply"}
		var parts []string
		cpi := 0.0
		for i := 0; i < k; i++ {
			c := edakit.Int(1, 5)
			cpi += float64(pct[i]) * float64(c) / 100
			parts = append(parts, edakit.Sprintf("%d%% %s instructions at %d cycles", pct[i], names[i], c))
		}
		return withExact(edakit.Numeric("A workload's mix is "+strings.Join(parts, ", ")+". What is its average CPI? Round to 2 decimal places.", edakit.Num(cpi, 2),
			edakit.Sprintf("The frequency-weighted mean of the class CPIs is %s.", edakit.Num(cpi, 2)), map[string]string{"cpi": edakit.Num(cpi, 2)}), cpi)
	default:
		p := float64(5*edakit.Int(4, 19)) / 100
		s := edakit.Pick([]float64{2, 3, 4, 5, 8, 10, 20, 50})
		S := 1 / ((1 - p) + p/s)
		return withExact(edakit.Numeric(edakit.Sprintf("An optimization speeds up a part of a program that takes %s%% of its execution time by a factor of %s and leaves the rest unchanged. What is the overall speedup? Round to 2 decimal places.", edakit.Num(p*100, 0), edakit.Num(s, 0)), edakit.Num(S, 2),
			edakit.Sprintf("Amdahl: $1 / ((1 - %s) + %s/%s) = %s$; even an infinite speedup would cap at %s.", edakit.Num(p, 2), edakit.Num(p, 2), edakit.Num(s, 0), edakit.Num(S, 2), edakit.Num(1/(1-p), 2)),
			map[string]string{"cap": edakit.Num(1/(1-p), 2)}), S)
	}
}

func pipelining(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 3, 6)
	var d []int
	D, mx := 0, 0
	for i := 0; i < k; i++ {
		x := 10 * edakit.Int(10, 40)
		d = append(d, x)
		D += x
		mx = max(mx, x)
	}
	treg := 5 * edakit.Int(6, 14)
	T := mx + treg
	stages := joinInts(d, ", ")
	switch variant(ctx) {
	case 1:
		return edakit.Numeric(edakit.Sprintf("Logic is split into %d pipeline stages with delays %s ps. Each pipeline register adds %d ps (clock-to-q plus setup). What is the clock period, in ps?", k, stages, treg), edakit.IntStr(T),
			edakit.Sprintf("The slowest stage is %d ps; plus %d ps of register overhead gives %d ps (against %d ps unpipelined).", mx, treg, T, D+treg), map[string]string{"unpipelined": edakit.IntStr(D + treg)})
	case 2:
		N := edakit.Pick([]int{5, 10, 20, 100, 1000})
		tot := (k + N - 1) * T
		return edakit.Numeric(edakit.Sprintf("A %d-stage pipeline has stage delays %s ps and %d ps of register overhead per stage. How long does it take to process %d independent items, from the first entering to the last leaving, in ps?", k, stages, treg, N), edakit.IntStr(tot),
			edakit.Sprintf("The period is %d ps; %d items take $(%d + %d - 1) \\times %d = %d$ ps.", T, N, k, N, T, tot), map[string]string{"period": edakit.IntStr(T)})
	default:
		N := edakit.Pick([]int{10, 50, 100, 1000})
		sp := float64(N*(D+treg)) / float64((k+N-1)*T)
		return withExact(edakit.Numeric(edakit.Sprintf("Logic with total delay %d ps (stages %s ps) is pipelined into %d stages, each register adding %d ps. Unpipelined, it would be one stage with one register. What is the speedup for %d items? Round to 2 decimal places.", D, stages, k, treg, N), edakit.Num(sp, 2),
			edakit.Sprintf("Unpipelined: $%d \\times %d$ ps. Pipelined: $(%d + %d - 1) \\times %d$ ps. Ratio %s.", N, D+treg, k, N, T, edakit.Num(sp, 2)), map[string]string{"period": edakit.IntStr(T)}), sp)
	}
}

type instr struct {
	load bool
	dst  int
	srcs []int
}

func (in instr) String() string {
	if in.load {
		return edakit.Sprintf("lw r%d, %d(r%d)", in.dst, 4*edakit.Int(0, 8), in.srcs[0])
	}
	op := edakit.Pick([]string{"add", "sub", "and", "or"})
	return edakit.Sprintf("%s r%d, r%d, r%d", op, in.dst, in.srcs[0], in.srcs[1])
}

// pipeCycles returns the total cycles of a classic 5-stage pipeline (register file written in the
// first half of WB, read in the second half of ID) with or without full forwarding.
func pipeCycles(prog []instr, fwd bool) int {
	D := make([]int, len(prog))
	for j, in := range prog {
		if j == 0 {
			D[j] = 2
			continue
		}
		D[j] = D[j-1] + 1
		for _, s := range in.srcs {
			for i := j - 1; i >= 0; i-- {
				if prog[i].dst == s {
					switch {
					case !fwd:
						D[j] = max(D[j], D[i]+3)
					case prog[i].load:
						D[j] = max(D[j], D[i]+2)
					}
					break
				}
			}
		}
	}
	return D[len(D)-1] + 3
}

func hazards(ctx generator.GeneratorContext) generator.Problem {
	v := variant(ctx)
	for {
		n := edakit.Scaled(ctx, 3, 6)
		var prog []instr
		for j := 0; j < n; j++ {
			in := instr{load: edakit.Int(0, 9) < 4, dst: edakit.Int(1, 7)}
			ns := 2
			if in.load {
				ns = 1
			}
			for s := 0; s < ns; s++ {
				r := edakit.Int(1, 7)
				if j > 0 && edakit.Int(0, 9) < 6 {
					r = prog[j-1-edakit.Int(0, min(1, j-1))].dst
				}
				in.srcs = append(in.srcs, r)
			}
			prog = append(prog, in)
		}
		nf, wf := pipeCycles(prog, false), pipeCycles(prog, true)
		base := n + 4
		if nf == base {
			continue
		}
		if v == 2 && wf == base && edakit.Int(0, 2) != 0 {
			continue
		}
		lines := make([]string, n)
		for i, in := range prog {
			lines[i] = in.String()
		}
		code := strings.Join(lines, "; ")
		model := "A classic five-stage pipeline (IF, ID, EX, MEM, WB) writes the register file in the first half of a cycle and reads it in the second half. "
		switch v {
		case 1:
			return edakit.Numeric(model+"There is no forwarding. How many stall cycles does it insert for the sequence: "+code+"?", edakit.IntStr(nf-base),
				edakit.Sprintf("A dependent instruction must decode at least 3 cycles after its producer decodes. The sequence takes %d cycles against %d without hazards: %d stalls.", nf, base, nf-base),
				map[string]string{"cycles": edakit.IntStr(nf)})
		case 2:
			return edakit.Numeric(model+"Full forwarding is present, so only a use immediately after a load stalls (1 cycle). How many cycles does the sequence take from the first fetch to the last write-back: "+code+"?", edakit.IntStr(wf),
				edakit.Sprintf("Without stalls %d instructions take %d cycles; load-use stalls add %d, giving %d.", n, base, wf-base, wf), map[string]string{"stalls": edakit.IntStr(wf - base)})
		default:
			return edakit.Numeric(model+"There is no forwarding. How many cycles does the sequence take from the first fetch to the last write-back: "+code+"?", edakit.IntStr(nf),
				edakit.Sprintf("%d instructions take %d cycles without hazards; the stalls add %d, giving %d. With forwarding it would take %d.", n, base, nf-base, nf, wf), map[string]string{"forwarded": edakit.IntStr(wf)})
		}
	}
}

func outcomes(n int) []bool {
	o := make([]bool, n)
	for i := range o {
		o[i] = edakit.Int(0, 9) < 6
	}
	return o
}

func tn(o []bool) string {
	var b strings.Builder
	for i, x := range o {
		if i > 0 {
			b.WriteString(", ")
		}
		if x {
			b.WriteByte('T')
		} else {
			b.WriteByte('N')
		}
	}
	return b.String()
}

func branchPrediction(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		fb := 5 * edakit.Int(2, 6)
		m := edakit.Int(2, 15)
		P := edakit.Int(2, 20)
		base := float64(edakit.Int(10, 15)) / 10
		cpi := base + float64(fb)/100*float64(m)/100*float64(P)
		return withExact(edakit.Numeric(edakit.Sprintf("A pipeline has base CPI %s. Branches are %d%% of instructions, %d%% of them are mispredicted, and each misprediction costs %d cycles. What is the CPI? Round to 3 decimal places.", edakit.Num(base, 1), fb, m, P), edakit.Num(cpi, 3),
			edakit.Sprintf("$%s + %s \\cdot %s \\cdot %d = %s$.", edakit.Num(base, 1), edakit.Num(float64(fb)/100, 2), edakit.Num(float64(m)/100, 2), P, edakit.Num(cpi, 3)), map[string]string{"penalty_cpi": edakit.Num(cpi-base, 3)}), cpi)
	case 2:
		o := outcomes(edakit.Scaled(ctx, 8, 12))
		st := edakit.Int(0, 3)
		s0 := st
		miss := 0
		for _, t := range o {
			if (st >= 2) != t {
				miss++
			}
			if t && st < 3 {
				st++
			} else if !t && st > 0 {
				st--
			}
		}
		return edakit.Numeric(edakit.Sprintf("A 2-bit saturating counter predictor (states 0 and 1 predict not taken, 2 and 3 predict taken; taken increments, not taken decrements) starts in state %d. The branch's outcomes are %s. How many are mispredicted?", s0, tn(o)), edakit.IntStr(miss),
			edakit.Sprintf("Tracing the counter through the %d outcomes gives %d mispredictions, ending in state %d.", len(o), miss, st), map[string]string{"final_state": edakit.IntStr(st)})
	default:
		var o []bool
		if edakit.Int(0, 1) == 0 {
			k := edakit.Int(3, 6)
			reps := edakit.Int(2, 4)
			for r := 0; r < reps; r++ {
				for i := 0; i < k-1; i++ {
					o = append(o, true)
				}
				o = append(o, false)
			}
		} else {
			o = outcomes(edakit.Scaled(ctx, 8, 12))
		}
		last := false
		miss := 0
		for _, t := range o {
			if last != t {
				miss++
			}
			last = t
		}
		return edakit.Numeric(edakit.Sprintf("A 1-bit predictor predicts that a branch repeats its last outcome; it starts predicting not taken. The outcomes are %s. How many are mispredicted?", tn(o)), edakit.IntStr(miss),
			edakit.Sprintf("A 1-bit predictor misses whenever the outcome differs from the previous one (and on a first taken outcome): %d misses.", miss), map[string]string{"misses": edakit.IntStr(miss)})
	}
}

func cacheGeometry() (A, cKB, B, ways int) {
	A = edakit.Pick([]int{32, 40, 48})
	cKB = 1 << uint(edakit.Int(2, 10))
	B = edakit.Pick([]int{16, 32, 64, 128})
	ways = edakit.Pick([]int{1, 2, 4, 8, 16})
	return
}

func waysName(w int) string {
	if w == 1 {
		return "direct-mapped"
	}
	return edakit.Sprintf("%d-way set-associative", w)
}

func cacheOrganization(ctx generator.GeneratorContext) generator.Problem {
	A, cKB, B, ways := cacheGeometry()
	sets := cKB * 1024 / (B * ways)
	off, idx := clog2(B), clog2(sets)
	tag := A - off - idx
	desc := edakit.Sprintf("A cache of %d KB, %s, with %d-byte blocks serves %d-bit byte addresses.", cKB, waysName(ways), B, A)
	switch variant(ctx) {
	case 1:
		switch edakit.Int(0, 2) {
		case 0:
			return edakit.Numeric(desc+" How many tag bits does each line store?", edakit.IntStr(tag),
				edakit.Sprintf("%d sets give %d index bits, %d-byte blocks give %d offset bits, so the tag is $%d - %d - %d = %d$ bits.", sets, idx, B, off, A, idx, off, tag), map[string]string{"index": edakit.IntStr(idx)})
		case 1:
			return edakit.Numeric(desc+" How many index bits does the address use?", edakit.IntStr(idx),
				edakit.Sprintf("$%d \\cdot 1024 / (%d \\cdot %d) = %d$ sets, so $\\log_2 %d = %d$ index bits.", cKB, B, ways, sets, sets, idx), map[string]string{"tag": edakit.IntStr(tag)})
		default:
			return edakit.Numeric(desc+" How many offset bits does the address use?", edakit.IntStr(off),
				edakit.Sprintf("$\\log_2 %d = %d$.", B, off), map[string]string{"tag": edakit.IntStr(tag)})
		}
	case 2:
		return edakit.Numeric(desc+" How many sets does it have?", edakit.IntStr(sets),
			edakit.Sprintf("$C / (B \\cdot a) = %d \\cdot 1024 / (%d \\cdot %d) = %d$.", cKB, B, ways, sets), map[string]string{"lines": edakit.IntStr(cKB * 1024 / B)})
	default:
		lines := cKB * 1024 / B
		per := 8*B + tag + 2
		return edakit.Numeric(desc+" It is write-back, so each line stores its data, its tag, a valid bit and a dirty bit. How many bits of storage does the cache hold in total?", edakit.IntStr(lines*per),
			edakit.Sprintf("%d lines of $%d + %d + 2 = %d$ bits: %d bits.", lines, 8*B, tag, per, lines*per), map[string]string{"tag": edakit.IntStr(tag)})
	}
}

// traceCache simulates an LRU set-associative cache over block numbers and returns hits.
func traceCache(blocks []int, sets, ways int) (int, []string) {
	content := make([][]int, sets) // most recent last
	hits := 0
	var res []string
	for _, b := range blocks {
		s := b % sets
		found := -1
		for i, x := range content[s] {
			if x == b {
				found = i
			}
		}
		if found >= 0 {
			hits++
			res = append(res, "hit")
			content[s] = append(append(content[s][:found:found], content[s][found+1:]...), b)
			continue
		}
		res = append(res, "miss")
		if len(content[s]) == ways {
			content[s] = content[s][1:]
		}
		content[s] = append(content[s], b)
	}
	return hits, res
}

func cachePerformance(ctx generator.GeneratorContext) generator.Problem {
	v := variant(ctx)
	if v == 1 {
		th := edakit.Int(1, 4)
		m1 := edakit.Int(2, 15)
		pen := 10 * edakit.Int(5, 30)
		if edakit.Int(0, 1) == 0 {
			amat := float64(th) + float64(m1)/100*float64(pen)
			return withExact(edakit.Numeric(edakit.Sprintf("A cache has a hit time of %s, a miss rate of %d%% and a miss penalty of %d cycles. What is the average memory access time, in cycles? Round to 2 decimal places.", plural(th, "cycle"), m1, pen), edakit.Num(amat, 2),
				edakit.Sprintf("$%d + %s \\cdot %d = %s$ cycles.", th, edakit.Num(float64(m1)/100, 2), pen, edakit.Num(amat, 2)), map[string]string{"amat": edakit.Num(amat, 2)}), amat)
		}
		t2 := edakit.Int(8, 20)
		m2 := 5 * edakit.Int(2, 12)
		amat := float64(th) + float64(m1)/100*(float64(t2)+float64(m2)/100*float64(pen))
		return withExact(edakit.Numeric(edakit.Sprintf("L1 hits in %s and misses %d%% of the time. L2 hits in %d cycles with a local miss rate of %d%%, and memory takes %d cycles. What is the average memory access time, in cycles? Round to 2 decimal places.", plural(th, "cycle"), m1, t2, m2, pen), edakit.Num(amat, 2),
			edakit.Sprintf("$%d + %s(%d + %s \\cdot %d) = %s$ cycles.", th, edakit.Num(float64(m1)/100, 2), t2, edakit.Num(float64(m2)/100, 2), pen, edakit.Num(amat, 2)), map[string]string{"amat": edakit.Num(amat, 2)}), amat)
	}
	S := edakit.Pick([]int{4, 8})
	B := edakit.Pick([]int{16, 32})
	s1 := edakit.Int(0, S-1)
	s2 := (s1 + edakit.Int(1, S-1)) % S
	pool := []int{s1, s1 + S, s2, s1 + 2*S, s2 + S}
	pool = pool[:edakit.Scaled(ctx, 3, 5)]
	L := edakit.Scaled(ctx, 7, 10)
	var blocks, addrs []int
	for i := 0; i < L; i++ {
		b := edakit.Pick(pool)
		if i > 0 && edakit.Int(0, 3) == 0 {
			b = blocks[i-1]
		}
		blocks = append(blocks, b)
		addrs = append(addrs, b*B+4*edakit.Int(0, B/4-1))
	}
	if v == 2 {
		h, res := traceCache(blocks, S, 1)
		return edakit.Numeric(edakit.Sprintf("A direct-mapped cache has %d sets of %d-byte blocks and starts empty. It sees the byte addresses %s. How many accesses hit?", S, B, joinInts(addrs, ", ")), edakit.IntStr(h),
			edakit.Sprintf("Blocks %s map to sets %s; the results are %s.", joinInts(blocks, ", "), joinInts(mods(blocks, S), ", "), strings.Join(res, ", ")), map[string]string{"hits": edakit.IntStr(h)})
	}
	h, res := traceCache(blocks, S/2, 2)
	return edakit.Numeric(edakit.Sprintf("A 2-way set-associative cache with LRU replacement has %d sets of two %d-byte blocks (%d blocks in all) and starts empty. It sees the byte addresses %s. How many accesses hit?", S/2, B, S, joinInts(addrs, ", ")), edakit.IntStr(h),
		edakit.Sprintf("Blocks %s map to sets %s; the results are %s.", joinInts(blocks, ", "), joinInts(mods(blocks, S/2), ", "), strings.Join(res, ", ")), map[string]string{"hits": edakit.IntStr(h)})
}

func mods(xs []int, m int) []int {
	out := make([]int, len(xs))
	for i, x := range xs {
		out[i] = x % m
	}
	return out
}

var datapathItems = []edakit.Item{
	{KP: 1, Q: "For `add r3, r1, r2`, which control setting routes the right value to the register file's write port?",
		A: "Select the ALU result, not the data-memory output, at the write-back mux", D: []string{"Select the data-memory output at the write-back mux", "Disable the register write enable", "Set the ALU to pass its first operand unchanged"},
		Why: "An ALU instruction writes its ALU result; only loads select the memory data."},
	{KP: 1, Q: "Which of these is a control signal rather than a datapath element?",
		A: "The register file's write enable", D: []string{"The ALU's adder", "The register file's storage array", "The data memory"},
		Why: "Datapath elements hold or transform data; control signals such as enables and mux selects tell them what to do each cycle."},
	{KP: 1, Q: "A load `lw r4, 8(r5)` and an add share the same datapath. What is different between them?",
		A: "Only the settings of the control signals", D: []string{"Each instruction has its own ALU", "The load uses a separate register file", "The add bypasses the instruction fetch"},
		Why: "One datapath serves every instruction; the decoder sets mux selects, ALU operation and enables differently for each."},
	{KP: 1, Q: "Which control signals does a store instruction `sw r2, 4(r1)` need asserted?",
		A: "ALU set to add for the address, memory write enabled, register write disabled", D: []string{"Register write enabled and memory read enabled", "ALU set to subtract and memory write disabled", "Write-back mux set to the memory output"},
		Why: "A store computes base plus offset, writes memory, and must not write the register file."},
	{KP: 2, Q: "Why must a single-cycle processor's clock period fit its slowest instruction?",
		A: "Every instruction completes in one cycle of the same length", D: []string{"The slowest instruction runs most often", "The control unit is an FSM with one state per instruction", "Faster instructions are stalled to match"},
		Why: "With one cycle per instruction and one clock, the period must cover the longest path, typically a load."},
	{KP: 2, Q: "What does a multicycle design gain by sharing one ALU across steps?",
		A: "Less hardware, at the cost of a controller FSM sequencing the steps", D: []string{"A shorter critical path with no control cost", "One instruction per cycle on average", "No need for intermediate registers"},
		Why: "Sharing units across cycles saves area but needs registers between steps and a state machine to sequence them."},
	{KP: 2, Q: "In a multicycle processor, what are the states of the control FSM?",
		A: "The steps of instruction execution, such as fetch, decode, execute, memory and write-back", D: []string{"The registers of the register file", "The instructions of the program", "The pipeline stages of other instructions"},
		Why: "The controller steps through the phases of each instruction, branching by opcode after decode."},
	{KP: 2, Q: "A high-level synthesis tool schedules a loop body into control steps and shares one multiplier. What hardware results?",
		A: "A datapath with a shared multiplier and muxes, driven by an FSM controller", D: []string{"A purely combinational circuit with no registers", "A separate multiplier for every multiplication", "A microcoded processor"},
		Why: "HLS produces exactly the multicycle structure: scheduling defines the FSM states and binding shares the units through muxes."},
	{KP: 3, Q: "What is the main advantage of microcoded control over hardwired control?",
		A: "The instruction behavior can be changed by rewriting the control store", D: []string{"It is always faster", "It needs no instruction decoder", "It eliminates the datapath"},
		Why: "Microcode keeps the control signals in a memory, which is flexible and suits complex instruction sets."},
	{KP: 3, Q: "Which kind of bug is more typical of the control unit than of the datapath?",
		A: "A missed stall when two instructions arrive in a particular order", D: []string{"An adder overflow at the most negative value", "A wrong bit width in a shifter", "A rounding error in a multiplier"},
		Why: "Control bugs are sequence-dependent; datapath bugs are usually data-dependent."},
	{KP: 3, Q: "Which verification method suits a control-unit FSM best?",
		A: "Model checking of temporal properties over its state transitions", D: []string{"Word-level arithmetic reasoning about adders", "Exhaustive simulation of a 64-bit multiplier", "Glitch power analysis"},
		Why: "Control logic has small data but complex sequencing, which model checking and FSM transition coverage target."},
	{KP: 3, Q: "Why do RISC processors and accelerators usually use hardwired control?",
		A: "Simple, regular instructions make fast synthesized control logic practical", D: []string{"Microcode cannot implement loads and stores", "Hardwired control needs no testing", "Hardwired control can be changed after fabrication"},
		Why: "A regular instruction set decodes with little logic, so hardwired control is smaller and faster than a control store."},
}
