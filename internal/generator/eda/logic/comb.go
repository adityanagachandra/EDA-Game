package logic

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerComb(reg *generator.Registry) {
	reg.Register("logic.comb.gates_mux", retryTies(gatesMux))
	reg.Register("logic.comb.gate_delay", retryTies(gateDelay))
	reg.Register("logic.comb.decoder", retryTies(decoder))
	reg.Register("logic.comb.encoder", retryTies(encoder))
	reg.Register("logic.comb.mux_shannon", retryTies(muxShannon))
	reg.Register("logic.comb.hazards", &edakit.Bank{Items: hazardItems})
	reg.Register("logic.comb.shifter", retryTies(shifter))
}

var varNames = []string{"a", "b", "c", "d"}

// expr is a small gate network: a leaf (input index) or a two-input gate over two children.
type expr struct {
	op   string
	l, r *expr
	leaf int
}

var gateOps = []string{"AND", "OR", "NAND", "NOR", "XOR", "XNOR"}

func gate(op string, x, y int) int {
	switch op {
	case "AND":
		return x & y
	case "OR":
		return x | y
	case "NAND":
		return 1 - x&y
	case "NOR":
		return 1 - (x | y)
	case "XOR":
		return x ^ y
	default:
		return 1 - (x ^ y)
	}
}

func (e *expr) eval(in []int) int {
	if e.op == "" {
		return in[e.leaf]
	}
	return gate(e.op, e.l.eval(in), e.r.eval(in))
}

func (e *expr) String() string {
	if e.op == "" {
		return varNames[e.leaf]
	}
	return e.op + "(" + e.l.String() + ", " + e.r.String() + ")"
}

// randExpr draws a two-level network over nv inputs that uses every input at least once.
func randExpr(nv int) *expr {
	for {
		leaf := func() *expr { return &expr{leaf: edakit.Int(0, nv-1)} }
		g := func() *expr {
			a, b := edakit.Int(0, nv-1), edakit.Int(0, nv-1)
			for b == a {
				b = edakit.Int(0, nv-1)
			}
			return &expr{op: edakit.Pick(gateOps), l: &expr{leaf: a}, r: &expr{leaf: b}}
		}
		var l, r *expr
		l = g()
		if edakit.Int(0, 3) == 0 {
			r = leaf()
		} else {
			r = g()
		}
		e := &expr{op: edakit.Pick(gateOps), l: l, r: r}
		s := e.String()
		ok := true
		for i := 0; i < nv; i++ {
			if !strings.Contains(s, varNames[i]+")") && !strings.Contains(s, varNames[i]+",") {
				ok = false
			}
		}
		ones := 0
		for m := 0; m < 1<<uint(nv); m++ {
			ones += e.eval(inputs(m, nv))
		}
		if ok && ones > 0 && ones < 1<<uint(nv) {
			return e
		}
	}
}

// inputs expands row m into input values, with the first variable as the most significant bit.
func inputs(m, nv int) []int {
	in := make([]int, nv)
	for i := 0; i < nv; i++ {
		in[i] = (m >> uint(nv-1-i)) & 1
	}
	return in
}

func assign(in []int) string {
	parts := make([]string, len(in))
	for i, v := range in {
		parts[i] = edakit.Sprintf("%s = %d", varNames[i], v)
	}
	return strings.Join(parts, ", ")
}

func gatesMux(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		nv := edakit.Scaled(ctx, 3, 4)
		e := randExpr(nv)
		m := edakit.Int(0, (1<<uint(nv))-1)
		in := inputs(m, nv)
		y := e.eval(in)
		return edakit.Numeric(edakit.Sprintf("Evaluate the gate network $y$ = %s at %s. What is $y$ (0 or 1)?", e, assign(in)), edakit.IntStr(y),
			edakit.Sprintf("Evaluate the inner gates first: %s gives %d and %s gives %d, then the output gate %s gives %d.", e.l, e.l.eval(in), e.r, e.r.eval(in), e.op, y),
			map[string]string{"y": edakit.IntStr(y)})
	case 2:
		if edakit.Int(0, 1) == 0 {
			data := uint64(edakit.Int(1, 254))
			sel := edakit.Int(0, 7)
			y := int(data >> uint(sel) & 1)
			return edakit.Numeric(edakit.Sprintf("An 8-to-1 multiplexer has data inputs $D_7 \\dots D_0$ = %s (so $D_7$ is the leftmost bit) and select lines $s_2 s_1 s_0$ = %s. What is its output (0 or 1)?", bin(data, 8), bin(uint64(sel), 3)), edakit.IntStr(y),
				edakit.Sprintf("The select lines read %d, so the output is $D_{%d}$ = %d.", sel, sel, y), map[string]string{"index": edakit.IntStr(sel)})
		}
		labels := []string{"0", "1", "c", "\\bar{c}"}
		weight := []int{0, 2, 1, 1}
		var picks []int
		total := 0
		for {
			picks = picks[:0]
			total = 0
			for i := 0; i < 4; i++ {
				p := edakit.Int(0, 3)
				picks = append(picks, p)
				total += weight[p]
			}
			if total > 0 && total < 8 {
				break
			}
		}
		return edakit.Numeric(edakit.Sprintf("A 4-to-1 multiplexer has select lines $a$ (high) and $b$ (low) and data inputs $(D_0, D_1, D_2, D_3) = (%s, %s, %s, %s)$. Over all eight combinations of $a, b, c$, on how many is the output 1?",
			labels[picks[0]], labels[picks[1]], labels[picks[2]], labels[picks[3]]), edakit.IntStr(total),
			"Each select value covers two rows (c = 0 and c = 1): a constant 1 contributes both rows, c or its complement contributes one, a constant 0 none. The total is "+edakit.IntStr(total)+".",
			map[string]string{"minterms": edakit.IntStr(total)})
	default:
		nv := edakit.Scaled(ctx, 3, 4)
		e := randExpr(nv)
		ones := 0
		for m := 0; m < 1<<uint(nv); m++ {
			ones += e.eval(inputs(m, nv))
		}
		return edakit.Numeric(edakit.Sprintf("The network $y$ = %s has %d inputs. On how many of the $2^{%d}$ input combinations is $y = 1$?", e, nv, nv), edakit.IntStr(ones),
			edakit.Sprintf("Enumerating all %d rows and evaluating the network on each gives %d rows with output 1.", 1<<uint(nv), ones),
			map[string]string{"onset": edakit.IntStr(ones)})
	}
}

// gateDelay draws a small DAG of two-input gates and asks for the latest (or earliest) output arrival.
func gateDelay(ctx generator.GeneratorContext) generator.Problem {
	v := variant(ctx)
	m := edakit.Scaled(ctx, 4, 7)
	names := []string{"a", "b", "c", "d"}
	arr := map[string]int{}
	early := map[string]int{}
	for _, n := range names {
		arr[n] = 0
		if v == 2 {
			arr[n] = 5 * edakit.Int(0, 12)
		}
		early[n] = arr[n]
	}
	var lines []string
	for i := 1; i <= m; i++ {
		g := edakit.Sprintf("g%d", i)
		pool := append([]string{}, names...)
		for j := 1; j < i; j++ {
			pool = append(pool, edakit.Sprintf("g%d", j))
		}
		x := pool[edakit.Int(0, len(pool)-1)]
		if i > 1 {
			x = edakit.Sprintf("g%d", edakit.Int(max(1, i-2), i-1))
		}
		y := pool[edakit.Int(0, len(pool)-1)]
		for y == x {
			y = pool[edakit.Int(0, len(pool)-1)]
		}
		d := 10 * edakit.Int(1, 6)
		arr[g] = max(arr[x], arr[y]) + d
		early[g] = min(early[x], early[y]) + d
		lines = append(lines, edakit.Sprintf("%s = %s(%s, %s) with delay %d ps", g, edakit.Pick(gateOps), x, y, d))
	}
	out := edakit.Sprintf("g%d", m)
	net := strings.Join(lines, "; ")
	switch v {
	case 1:
		return edakit.Numeric("Primary inputs a, b, c and d all arrive at 0 ps. The network is: "+net+". The output is "+out+". At what time, in ps, does the output settle (its latest arrival)?",
			edakit.IntStr(arr[out]), edakit.Sprintf("Propagate in topological order, taking the latest input arrival at each gate and adding its delay. The output %s settles at %d ps.", out, arr[out]),
			map[string]string{"arrival": edakit.IntStr(arr[out])})
	case 2:
		return edakit.Numeric(edakit.Sprintf("Primary inputs arrive at a = %d ps, b = %d ps, c = %d ps and d = %d ps. The network is: %s. The output is %s. At what time, in ps, does the output settle?", arr["a"], arr["b"], arr["c"], arr["d"], net, out),
			edakit.IntStr(arr[out]), edakit.Sprintf("Start each input at its arrival time and propagate the maximum through the gates. The output %s settles at %d ps.", out, arr[out]),
			map[string]string{"arrival": edakit.IntStr(arr[out])})
	default:
		return edakit.Numeric("Primary inputs a, b, c and d all change at 0 ps. The network is: "+net+". The output is "+out+". Using each gate's stated delay as its minimum delay too, what is the earliest time, in ps, at which the output can start to change (the contamination delay)?",
			edakit.IntStr(early[out]), edakit.Sprintf("Propagate the minimum instead of the maximum: the shortest path reaches %s at %d ps, while the longest reaches it at %d ps.", out, early[out], arr[out]),
			map[string]string{"min_arrival": edakit.IntStr(early[out]), "max_arrival": edakit.IntStr(arr[out])})
	}
}

func decoder(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 5)
		x := uint64(edakit.Int(0, (1<<uint(n))-1))
		low := edakit.Int(0, 1) == 1
		kind := "active-high"
		if low {
			kind = "active-low"
		}
		return edakit.Numeric(edakit.Sprintf("A %d-to-%d decoder with %s outputs $Y_0 \\dots Y_{%d}$ receives the input %s (most significant bit first). Which output index is asserted?", n, 1<<uint(n), kind, (1<<uint(n))-1, bin(x, n)), edakit.IntStr(int(x)),
			edakit.Sprintf("The input reads %d in binary, so $Y_{%d}$ is asserted (driven %s) and all the others are not.", x, x, map[bool]string{true: "low", false: "high"}[low]),
			map[string]string{"index": edakit.IntStr(int(x))})
	case 2:
		k := edakit.Pick([]int{2, 3})
		levels := edakit.Scaled(ctx, 2, 3)
		n := k * levels
		cnt, term := 0, 1
		var parts []string
		for i := 0; i < levels; i++ {
			cnt += term
			parts = append(parts, edakit.IntStr(term))
			term <<= uint(k)
		}
		return edakit.Numeric(edakit.Sprintf("A %d-to-%d decoder is built as a tree of %d-to-%d decoders with enable inputs. How many of the small decoders does it need?", n, 1<<uint(n), k, 1<<uint(k)), edakit.IntStr(cnt),
			edakit.Sprintf("Each level enables $2^{%d}$ decoders below every output of the level above: %s, which sums to %d.", k, strings.Join(parts, " + "), cnt),
			map[string]string{"levels": edakit.IntStr(levels)})
	default:
		a := edakit.Pick([]int{16, 20, 24})
		k := edakit.Int(2, 8)
		if edakit.Int(0, 1) == 0 {
			ans := 1 << uint(a-k)
			return edakit.Numeric(edakit.Sprintf("On a %d-bit byte-address bus, a device is selected when the top %d address bits match a fixed pattern, and the decoder compares no other bits. How many addresses select the device?", a, k), edakit.IntStr(ans),
				edakit.Sprintf("The remaining %d bits are free, so $2^{%d} = %d$ addresses select it.", a-k, a-k, ans), map[string]string{"free_bits": edakit.IntStr(a - k)})
		}
		d := edakit.Int(4, a-k-2)
		al := 1 << uint(a-k-d)
		return edakit.Numeric(edakit.Sprintf("On a %d-bit address bus, a device uses %d low address bits internally, and its select logic compares only the top %d address bits. At how many distinct base addresses does the device appear (its aliases)?", a, d, k), edakit.IntStr(al),
			edakit.Sprintf("%d address bits are neither compared nor used by the device, so it appears $2^{%d} = %d$ times.", a-k-d, a-k-d, al),
			map[string]string{"ignored_bits": edakit.IntStr(a - k - d)})
	}
}

func groups4(s string) string {
	var parts []string
	for i := 0; i < len(s); i += 4 {
		parts = append(parts, s[i:i+4])
	}
	return strings.Join(parts, " ")
}

func encoder(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1, 2:
		hi := ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0)
		w := edakit.Pick([]int{8, 16})
		var r uint64
		for popcount(r) < 2 {
			r = uint64(edakit.Int(1, (1<<uint(w))-1))
		}
		top, low := 63-leadingZeros(r), trailingZeros(r)
		if hi {
			return edakit.Numeric(edakit.Sprintf("A priority encoder where the highest index wins receives the request vector $R_{%d} \\dots R_0$ = %s. What index does it output?", w-1, groups4(bin(r, w))), edakit.IntStr(top),
				edakit.Sprintf("The highest asserted bit is bit %d; the lower requests are ignored.", top), map[string]string{"lowest": edakit.IntStr(low)})
		}
		return edakit.Numeric(edakit.Sprintf("A fixed-priority arbiter grants the lowest-indexed request. It receives $R_{%d} \\dots R_0$ = %s. Which index does it grant?", w-1, groups4(bin(r, w))), edakit.IntStr(low),
			edakit.Sprintf("The lowest asserted bit is bit %d; $x \\mathbin{\\&} (-x)$ isolates exactly that bit.", low), map[string]string{"highest": edakit.IntStr(top)})
	default:
		w := edakit.Pick([]int{16, 32})
		p := edakit.Int(0, w-1)
		x := uint64(1) << uint(p)
		x |= uint64(edakit.Int(0, (1<<uint(p))-1))
		lz := w - 1 - p
		return edakit.Numeric(edakit.Sprintf("A floating-point normalizer must shift the %d-bit word %s left until its top bit is 1. What is the leading-zero count (the shift amount)?", w, groups4(bin(x, w))), edakit.IntStr(lz),
			edakit.Sprintf("The highest set bit is bit %d, so the count is %d - 1 - %d = %d.", p, w, p, lz), map[string]string{"top_bit": edakit.IntStr(p)})
	}
}

func leadingZeros(x uint64) int {
	n := 0
	for i := 63; i >= 0 && x>>uint(i)&1 == 0; i-- {
		n++
	}
	return n
}

func trailingZeros(x uint64) int {
	n := 0
	for x&1 == 0 && n < 64 {
		x >>= 1
		n++
	}
	return n
}

func muxShannon(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 7)
		if edakit.Int(0, 1) == 0 {
			ans := (1 << uint(n)) - 1
			return edakit.Numeric(edakit.Sprintf("An arbitrary function of %d variables is built as a tree of 2-to-1 multiplexers whose leaves are the constants 0 and 1 (its truth table). How many 2-to-1 multiplexers does the tree need?", n), edakit.IntStr(ans),
				edakit.Sprintf("Expanding on one variable per level gives $1 + 2 + \\dots + 2^{%d} = 2^{%d} - 1 = %d$ muxes.", n-1, n, ans), map[string]string{"levels": edakit.IntStr(n)})
		}
		ans := (1 << uint(n-1)) - 1
		return edakit.Numeric(edakit.Sprintf("An arbitrary function of %d variables is built as a tree of 2-to-1 multiplexers, with the last variable fed into the data inputs (each leaf is 0, 1, the variable or its complement). How many 2-to-1 multiplexers does the tree need?", n), edakit.IntStr(ans),
			edakit.Sprintf("The select variables are the first %d, so the tree has $2^{%d} - 1 = %d$ muxes.", n-1, n-1, ans), map[string]string{"levels": edakit.IntStr(n - 1)})
	case 2:
		nv := edakit.Scaled(ctx, 3, 4)
		e := randExpr(nv)
		init := 0
		for m := 0; m < 1<<uint(nv); m++ {
			if e.eval(inputs(m, nv)) == 1 {
				init |= 1 << uint(m)
			}
		}
		idx := "4a + 2b + c"
		if nv == 4 {
			idx = "8a + 4b + 2c + d"
		}
		return edakit.Numeric(edakit.Sprintf("A %d-input LUT implements $f$ = %s. Bit $i$ of its configuration word holds $f$ on the row with index $i = %s$. What is the configuration word, as an unsigned decimal integer?", nv, e, idx), edakit.IntStr(init),
			edakit.Sprintf("The truth table, row 0 first, is %s; reading row i as bit i gives %d.", reverse(bin(uint64(init), 1<<uint(nv))), init),
			map[string]string{"init": edakit.IntStr(init)})
	default:
		k := edakit.Int(4, 6)
		n := edakit.Scaled(ctx, k+2, 64)
		if edakit.Int(0, 1) == 0 {
			ans := (n - 1 + k - 2) / (k - 1)
			return edakit.Numeric(edakit.Sprintf("An AND of %d inputs is mapped onto %d-input LUTs arranged as a tree. How many LUTs are needed?", n, k), edakit.IntStr(ans),
				edakit.Sprintf("Each LUT reduces the signal count by %d, so $\\lceil (%d - 1)/%d \\rceil = %d$ LUTs.", k-1, n, k-1, ans), map[string]string{"k": edakit.IntStr(k)})
		}
		lv, cap := 0, 1
		for cap < n {
			cap *= k
			lv++
		}
		return edakit.Numeric(edakit.Sprintf("An AND of %d inputs is mapped onto %d-input LUTs arranged as a tree. How many LUT levels deep is the tree?", n, k), edakit.IntStr(lv),
			edakit.Sprintf("Each level multiplies the number of inputs covered by %d; %d levels cover %d inputs, the first count at least %d.", k, lv, cap, n), map[string]string{"k": edakit.IntStr(k)})
	}
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func shifter(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := edakit.Pick([]int{8, 8, 16})
		mask := uint64(1)<<uint(n) - 1
		x := uint64(edakit.Int(1, int(mask)))
		s := edakit.Int(1, n-1)
		var r uint64
		op := edakit.Int(0, 3)
		name := []string{"rotated left", "rotated right", "logically shifted left", "logically shifted right"}[op]
		switch op {
		case 0:
			r = (x<<uint(s) | x>>uint(n-s)) & mask
		case 1:
			r = (x>>uint(s) | x<<uint(n-s)) & mask
		case 2:
			r = (x << uint(s)) & mask
		default:
			r = x >> uint(s)
		}
		return edakit.Numeric(edakit.Sprintf("The %d-bit word %s is %s by %d. What is the result as an unsigned decimal integer?", n, groups4(bin(x, n)), name, s), edakit.IntStr(int(r)),
			edakit.Sprintf("The result pattern is %s, which is %d.", groups4(bin(r, n)), r), map[string]string{"pattern": bin(r, n)})
	case 2:
		n := edakit.Pick([]int{8, 12, 16, 24, 32, 64})
		st := clog2(n)
		return edakit.Numeric(edakit.Sprintf("A logarithmic rotator for %d-bit words supports every rotate amount from 0 to %d. How many 2-to-1 multiplexers does it use?", n, n-1), edakit.IntStr(n*st),
			edakit.Sprintf("It needs $\\lceil \\log_2 %d \\rceil = %d$ stages of %d muxes each: %d.", n, st, n, n*st), map[string]string{"stages": edakit.IntStr(st)})
	default:
		n := edakit.Pick([]int{16, 32, 64})
		s := edakit.Int(3, n-1)
		return edakit.Numeric(edakit.Sprintf("A logarithmic shifter for %d-bit words shifts by %d. How many of its stages actually shift, rather than pass their input through?", n, s), edakit.IntStr(popcount(uint64(s))),
			edakit.Sprintf("%d is %s in binary; one stage acts per set bit, so %d stages shift.", s, bin(uint64(s), clog2(n)), popcount(uint64(s))),
			map[string]string{"amount": bin(uint64(s), clog2(n))})
	}
}

var hazardItems = []edakit.Item{
	{KP: 1, Q: "In $f = ab + \\bar{a}c$ with $b = c = 1$, input $a$ falls from 1 to 0. The inverter producing $\\bar{a}$ is slower than the direct path. What can the output do?",
		A: "Dip briefly to 0 before returning to 1", D: []string{"Rise briefly to 1 before settling at 0", "Stay at 1 with no transition at all", "Settle at 0 permanently"},
		Why: "The term ab turns off at once while the complemented term turns on one inverter delay later; for that interval both terms are 0, a static-1 hazard."},
	{KP: 1, Q: "Which of these is a static-0 hazard?",
		A: "An output that should stay at 0 pulses briefly to 1", D: []string{"An output that should stay at 1 dips briefly to 0", "An output that should change once changes three times", "An output that never changes when it should"},
		Why: "Static-0 means the correct value is a steady 0 and a glitch to 1 appears; static-1 is the dual."},
	{KP: 1, Q: "What is the minimum structural condition for a dynamic hazard (several transitions where one is expected)?",
		A: "At least three paths of different delay from the changing input to the output", D: []string{"Two paths of equal delay", "A single path through an inverter", "Any gate with fan-in greater than two"},
		Why: "A dynamic hazard needs the output to be pulled back and forth, which takes at least three reconverging paths with different delays."},
	{KP: 1, Q: "A static-1 hazard is observed on a single-input change in a two-level SOP circuit. What does that say about the two adjacent minterms involved?",
		A: "They are covered by different product terms and no single term covers both", D: []string{"They are both covered by the same product term", "One of them is a don't-care", "They differ in two input variables"},
		Why: "When one term covers both minterms it stays on through the change; the hazard appears only when coverage passes from one term to another."},
	{KP: 2, Q: "Which term removes the static-1 hazard from $f = ab + \\bar{a}c$?",
		A: "$bc$", D: []string{"$\\bar{a}\\bar{b}$", "$ac$", "$\\bar{b}\\bar{c}$"},
		Why: "The consensus of ab and a-bar c on variable a is bc. With b = c = 1 it holds the output at 1 while a changes."},
	{KP: 2, Q: "Adding a consensus term to remove a hazard changes the function how?",
		A: "Not at all: the consensus term is logically redundant", D: []string{"It adds new minterms to the on-set", "It removes the minterm where the glitch occurred", "It complements the output during the transition"},
		Why: "By the consensus theorem, xy + x-bar z + yz equals xy + x-bar z, so the function is unchanged; only its timing behavior improves."},
	{KP: 2, Q: "A minimizer drops every redundant prime implicant from an SOP. What is the effect on static-1 hazards for single-input changes?",
		A: "Hazards can reappear, since the dropped primes may have covered adjacent minterm pairs", D: []string{"All hazards are removed, since the cover is minimal", "Static-0 hazards become static-1 hazards", "Hazards cannot exist in a minimal cover"},
		Why: "A cover containing all prime implicants is free of single-input static-1 hazards; minimization removes exactly the redundant primes that provided that guarantee."},
	{KP: 2, Q: "For $f = \\bar{a}\\bar{b} + bc$, which single-input change can produce a static-1 hazard?",
		A: "$b$ changing while $a = 0$ and $c = 1$", D: []string{"$a$ changing while $b = 1$ and $c = 1$", "$c$ changing while $a = 1$ and $b = 0$", "$a$ changing while $b = 0$ and $c = 0$"},
		Why: "With a = 0 and c = 1, the term a-bar b-bar holds f at 1 for b = 0 and bc holds it for b = 1; coverage passes between terms, and the consensus a-bar c is missing."},
	{KP: 3, Q: "Why does a synchronous register-to-register path usually tolerate glitches in its combinational logic?",
		A: "The capture flip-flop samples only after the logic has settled within the setup window", D: []string{"Flip-flops filter out any pulse shorter than their clock-to-q", "Static timing analysis removes glitches during optimization", "Glitches cannot occur in logic synthesized by a tool"},
		Why: "If the path meets setup, the data input is stable when sampled, so earlier transitions do not change the stored value."},
	{KP: 3, Q: "Where is a glitch a functional bug even in a synchronous design?",
		A: "On an asynchronous reset or a signal used as a clock", D: []string{"On a data path that meets setup and hold", "On the input of a pipeline register", "On a bus that is sampled at the clock edge"},
		Why: "Asynchronous resets and clocks act immediately, so a glitch on them resets or clocks flip-flops at the wrong time."},
	{KP: 3, Q: "Why does a standard integrated clock-gating cell contain a latch?",
		A: "To hold the enable stable while the clock is high, so the gated clock cannot glitch", D: []string{"To add clock skew that fixes hold violations", "To store the clock's previous value for power analysis", "To divide the clock frequency by two"},
		Why: "The latch is transparent while the clock is low and holds while it is high, so enable changes cannot chop a clock pulse."},
	{KP: 3, Q: "In which circuit is glitch power typically a large share of dynamic power?",
		A: "A deep array multiplier with many reconverging paths of different lengths", D: []string{"A shift register with no logic between flip-flops", "A single inverter driving a long wire", "An SRAM cell array holding its data"},
		Why: "Unequal path delays through deep arithmetic make outputs toggle several times per cycle, and each spurious toggle burns energy."},
	{KP: 3, Q: "A lint tool flags combinational logic driving an asynchronous reset pin. What is the risk it is guarding against?",
		A: "A glitch from that logic could reset flip-flops spuriously", D: []string{"The reset would violate the setup time of the flip-flop", "The logic would be removed by synthesis", "The reset would take two cycles to propagate"},
		Why: "The asynchronous reset acts on any pulse, so hazards in the driving logic become spurious resets."},
}
