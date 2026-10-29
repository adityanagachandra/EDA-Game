package logic

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerArith(reg *generator.Registry) {
	reg.Register("logic.arith.full_adder", retryTies(fullAdder))
	reg.Register("logic.arith.ripple_carry", retryTies(rippleCarry))
	reg.Register("logic.arith.carry_lookahead", retryTies(carryLookahead))
	reg.Register("logic.arith.carry_select", retryTies(carrySelect))
	reg.Register("logic.arith.prefix_adder", retryTies(prefixAdder))
	reg.Register("logic.arith.kogge_brent", retryTies(koggeBrent))
	reg.Register("logic.arith.subtract_compare", retryTies(subtractCompare))
	reg.Register("logic.arith.multiplier", retryTies(multiplier))
	reg.Register("logic.arith.booth", retryTies(booth))
	reg.Register("logic.arith.wallace_dadda", retryTies(wallaceDadda))
}

// carries returns the carry into each bit position 1..n for a + b with carry-in 0, as a bit mask
// where bit i is the carry out of column i.
func carries(a, b uint64, n int) uint64 {
	var c, out uint64
	for i := 0; i < n; i++ {
		ai, bi := a>>uint(i)&1, b>>uint(i)&1
		co := ai&bi | c&(ai^bi)
		out |= co << uint(i)
		c = co
	}
	return out
}

func operands(n int) (uint64, uint64) {
	mask := (uint64(1) << uint(n)) - 1
	return uint64(edakit.Int(1, int(mask))), uint64(edakit.Int(1, int(mask)))
}

func fullAdder(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := edakit.Pick([]int{8, 8, 12})
		a, b := operands(n)
		c := carries(a, b, n)
		cnt := popcount(c)
		return edakit.Numeric(edakit.Sprintf("Add the %d-bit numbers %s and %s with carry-in 0. In how many of the %d columns is a carry produced (a carry out of that column)?", n, groups4(bin(a, n)), groups4(bin(b, n)), n), edakit.IntStr(cnt),
			edakit.Sprintf("Column by column, the carry-out pattern (column %d first) is %s, which has %d ones.", n-1, groups4(bin(c, n)), cnt), map[string]string{"carries": bin(c, n)})
	case 2:
		n := edakit.Pick([]int{8, 8, 12, 16})
		a, b := operands(n)
		if edakit.Int(0, 1) == 0 {
			g := popcount(a & b)
			return edakit.Numeric(edakit.Sprintf("For $a$ = %s and $b$ = %s, how many columns generate a carry ($g_i = a_i b_i = 1$)?", groups4(bin(a, n)), groups4(bin(b, n))), edakit.IntStr(g),
				edakit.Sprintf("$a \\mathbin{\\&} b$ = %s has %d ones.", groups4(bin(a&b, n)), g), map[string]string{"g": bin(a&b, n)})
		}
		p := popcount(a ^ b)
		return edakit.Numeric(edakit.Sprintf("For $a$ = %s and $b$ = %s, how many columns propagate a carry ($p_i = a_i \\oplus b_i = 1$)?", groups4(bin(a, n)), groups4(bin(b, n))), edakit.IntStr(p),
			edakit.Sprintf("$a \\oplus b$ = %s has %d ones.", groups4(bin(a^b, n)), p), map[string]string{"p": bin(a^b, n)})
	default:
		n := edakit.Scaled(ctx, 4, 64)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("A full adder is built from two half adders and an OR gate (2 XOR, 2 AND, 1 OR). How many gates does a ripple-carry adder of %d columns of such full adders use, with a carry-in input?", n), edakit.IntStr(5*n),
				edakit.Sprintf("Five gates per column: $5 \\times %d = %d$.", n, 5*n), map[string]string{"per_column": "5"})
		}
		ans := 5*(n-1) + 2
		return edakit.Numeric(edakit.Sprintf("A ripple-carry adder of %d columns has no carry-in, so its least significant column is a half adder (1 XOR, 1 AND) and the rest are full adders of 2 XOR, 2 AND and 1 OR. How many gates does it use?", n), edakit.IntStr(ans),
			edakit.Sprintf("$5(%d - 1) + 2 = %d$.", n, ans), map[string]string{"full_adders": edakit.IntStr(n - 1)})
	}
}

func rippleCarry(ctx generator.GeneratorContext) generator.Problem {
	dc, ds := edakit.Int(2, 3), edakit.Int(2, 3)
	model := edakit.Sprintf("Each full adder's carry-out settles %d gate delays after its carry-in (column 0's carry-in and all operand bits are ready at time 0), and its sum settles %d gate delays after its carry-in.", dc, ds)
	switch variant(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 8, 64)
		if edakit.Int(0, 2) == 0 {
			return edakit.Numeric(edakit.Sprintf("In a ripple-carry adder of %d bits, ", n)+lowerFirst(model)+" When does the final carry-out settle, in gate delays?", edakit.IntStr(n*dc),
				edakit.Sprintf("The carry passes through all %d columns: $%d \\times %d = %d$.", n, n, dc, n*dc), map[string]string{"carry_out": edakit.IntStr(n * dc)})
		}
		ans := (n-1)*dc + ds
		return edakit.Numeric(edakit.Sprintf("In a ripple-carry adder of %d bits, ", n)+lowerFirst(model)+" In the worst case, when does the most significant sum bit settle, in gate delays?", edakit.IntStr(ans),
			edakit.Sprintf("The carry into the top column arrives after $(%d - 1) \\times %d$, then the sum takes %d more: %d.", n, dc, ds, ans), map[string]string{"carry_into_msb": edakit.IntStr((n - 1) * dc)})
	case 2:
		n := edakit.Scaled(ctx, 8, 64)
		tg := 5 * edakit.Int(2, 6)
		gd := (n-1)*dc + ds
		return edakit.Numeric(edakit.Sprintf("In a ripple-carry adder of %d bits, ", n)+lowerFirst(model)+edakit.Sprintf(" One gate delay is %d ps. What is the worst-case delay to the most significant sum bit, in ps?", tg), edakit.IntStr(gd*tg),
			edakit.Sprintf("$(%d - 1) \\cdot %d + %d = %d$ gate delays, times %d ps, is %d ps.", n, dc, ds, gd, tg, gd*tg), map[string]string{"gate_delays": edakit.IntStr(gd)})
	default:
		budget := edakit.Int(20, 100)
		for budget < ds+dc {
			budget += 10
		}
		n := (budget-ds)/dc + 1
		return edakit.Numeric("For a ripple-carry adder, "+lowerFirst(model)+edakit.Sprintf(" What is the widest adder, in bits, whose most significant sum bit settles within %d gate delays?", budget), edakit.IntStr(n),
			edakit.Sprintf("Need $(n - 1) \\cdot %d + %d \\le %d$, so $n - 1 \\le %d$ and n = %d.", dc, ds, budget, (budget-ds)/dc, n), map[string]string{"budget": edakit.IntStr(budget)})
	}
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'A' && b[0] <= 'Z' {
		b[0] += 'a' - 'A'
	}
	return string(b)
}

func carryLookahead(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := 8
		a, b := operands(n)
		if edakit.Int(0, 2) == 0 {
			k := popcount(^(a | b) & 0xFF)
			return edakit.Numeric(edakit.Sprintf("For $a$ = %s and $b$ = %s, how many columns are kill columns ($a_i = b_i = 0$), which absorb any incoming carry?", bin(a, n), bin(b, n)), edakit.IntStr(k),
				edakit.Sprintf("The NOR of the operands is %s, with %d ones.", bin(^(a|b)&0xFF, n), k), map[string]string{"kill": edakit.IntStr(k)})
		}
		g, p := popcount(a&b), popcount(a^b)
		return edakit.Numeric(edakit.Sprintf("For $a$ = %s and $b$ = %s, how many columns either generate or propagate a carry?", bin(a, n), bin(b, n)), edakit.IntStr(g+p),
			edakit.Sprintf("Generate columns: %d (from $a \\mathbin{\\&} b$); propagate columns: %d (from $a \\oplus b$). They never overlap, so %d.", g, p, g+p), map[string]string{"g": edakit.IntStr(g), "p": edakit.IntStr(p)})
	case 2:
		var g, p uint64
		for i := 0; i < 4; i++ {
			switch edakit.Int(0, 2) {
			case 0:
				g |= 1 << uint(i)
			case 1:
				p |= 1 << uint(i)
			}
		}
		c0 := uint64(edakit.Int(0, 1))
		c := c0
		var vec uint64
		for i := 0; i < 4; i++ {
			c = g>>uint(i)&1 | p>>uint(i)&1&c
			vec |= c << uint(i)
		}
		return edakit.Numeric(edakit.Sprintf("A 4-bit lookahead block has $g_3 g_2 g_1 g_0$ = %s, $p_3 p_2 p_1 p_0$ = %s and carry-in $c_0 = %d$. Using $c_{i+1} = g_i + p_i c_i$ unrolled, enter the carry vector $c_4 c_3 c_2 c_1$ read as an unsigned decimal integer.", bin(g, 4), bin(p, 4), c0), edakit.IntStr(int(vec)),
			edakit.Sprintf("The carries are $c_4 c_3 c_2 c_1$ = %s, which is %d.", bin(vec, 4), vec), map[string]string{"carries": bin(vec, 4)})
	default:
		B := edakit.Scaled(ctx, 2, 16)
		db := edakit.Int(2, 3)
		n := 4 * B
		ans := db*B + 2
		return edakit.Numeric(edakit.Sprintf("An adder of %d bits is %d rippled 4-bit lookahead blocks. Model: g and p take 1 gate delay; each block produces its internal carries and carry-out %d gate delays after both its carry-in and its g, p are ready; each sum needs 1 more gate delay. When does the last sum bit settle, in gate delays?", n, B, db), edakit.IntStr(ans),
			edakit.Sprintf("Block k finishes at $1 + %d k$; the last block at $1 + %d \\cdot %d$, and its sums one later: %d.", db, db, B, ans), map[string]string{"blocks": edakit.IntStr(B)})
	}
}

func carrySelect(ctx generator.GeneratorContext) generator.Problem {
	tc := edakit.Int(2, 3)
	tm := edakit.Int(2, 3)
	switch variant(ctx) {
	case 1:
		n := edakit.Pick([]int{16, 24, 32, 48, 64})
		ks := []int{}
		for k := 2; k <= n/2; k++ {
			if n%k == 0 {
				ks = append(ks, k)
			}
		}
		k := edakit.Pick(ks)
		ans := k*tc + (n/k-1)*tm
		return edakit.Numeric(edakit.Sprintf("A carry-select adder of %d bits uses uniform %d-bit blocks. Each block ripples at %d gate delays per bit, and each select mux takes %d gate delays. When does the final carry-out settle, in gate delays?", n, k, tc, tm), edakit.IntStr(ans),
			edakit.Sprintf("All blocks finish rippling at $%d \\cdot %d$; the carry then crosses the remaining %d select muxes: $%d + %d \\cdot %d = %d$.", k, tc, n/k-1, k*tc, n/k-1, tm, ans), map[string]string{"blocks": edakit.IntStr(n / k)})
	case 2:
		n := edakit.Pick([]int{16, 32, 64, 128})
		best, bk := 1<<30, 0
		var cand []string
		for k := 2; k <= n/2; k *= 2 {
			d := k*tc + (n/k-1)*tm
			cand = append(cand, edakit.Sprintf("k = %d gives %d", k, d))
			if d < best {
				best, bk = d, k
			}
		}
		return edakit.Numeric(edakit.Sprintf("A carry-select adder of %d bits uses uniform blocks of k bits, with k a power of two between 2 and %d. Rippling costs %d gate delays per bit and each select mux %d gate delays, so the carry-out settles at $k \\cdot %d + (%d/k - 1) \\cdot %d$. What is the smallest achievable delay, in gate delays?", n, n/2, tc, tm, tc, n, tm), edakit.IntStr(best),
			edakit.Sprintf("Evaluating each block size: %s. The best is k = %d with %d gate delays.", joinStr(cand), bk, best), map[string]string{"best_k": edakit.IntStr(bk)})
	default:
		n := edakit.Pick([]int{16, 24, 32, 48, 64})
		k := edakit.Pick([]int{4, 8})
		ts := edakit.Int(2, 3)
		ans := 2*k*tc + (n/k-2)*ts
		return edakit.Numeric(edakit.Sprintf("A carry-skip adder of %d bits uses %d-bit blocks. In the worst case a carry ripples through the first block (%d gate delays per bit), skips each middle block (%d gate delays per skip), and ripples through the last block. What is that worst-case delay, in gate delays?", n, k, tc, ts), edakit.IntStr(ans),
			edakit.Sprintf("$%d \\cdot %d + (%d - 2) \\cdot %d + %d \\cdot %d = %d$.", k, tc, n/k, ts, k, tc, ans), map[string]string{"middle_blocks": edakit.IntStr(n/k - 2)})
	}
}

func joinStr(xs []string) string {
	out := ""
	for i, x := range xs {
		if i > 0 {
			out += "; "
		}
		out += x
	}
	return out
}

func prefixAdder(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 128)
		return edakit.Numeric(edakit.Sprintf("What is the minimum number of levels of two-input prefix nodes needed to compute all carries of an adder of %d bits?", n), edakit.IntStr(clog2(n)),
			edakit.Sprintf("The top carry depends on all %d inputs through two-input nodes, so at least $\\lceil \\log_2 %d \\rceil = %d$ levels; Sklansky and Kogge-Stone achieve it.", n, n, clog2(n)), map[string]string{"n": edakit.IntStr(n)})
	case 2:
		n := 8
		a, b := operands(n)
		c := carries(a, b, n)
		cnt := popcount(c)
		return edakit.Numeric(edakit.Sprintf("A prefix adder computes $G_{i:0}$ for $i = 0, \\dots, 7$ with carry-in 0. For $a$ = %s and $b$ = %s, for how many i is $G_{i:0} = 1$?", bin(a, n), bin(b, n)), edakit.IntStr(cnt),
			edakit.Sprintf("$G_{i:0}$ is the carry out of column i. The carry-out pattern (column 7 first) is %s, with %d ones.", bin(c, n), cnt), map[string]string{"carries": bin(c, n)})
	default:
		L := edakit.Scaled(ctx, 3, 7)
		n := 1 << uint(L)
		if edakit.Int(0, 2) == 0 {
			return edakit.Numeric(edakit.Sprintf("In a Sklansky prefix network for %d bits, how many prefix nodes does the single most heavily loaded node drive at the last level?", n), edakit.IntStr(n/2),
				edakit.Sprintf("At the last level, the prefix of the lower half feeds every position of the upper half: %d nodes.", n/2), map[string]string{"n": edakit.IntStr(n)})
		}
		ans := n / 2 * L
		return edakit.Numeric(edakit.Sprintf("How many prefix nodes does a Sklansky (divide-and-conquer) prefix network for %d bits use?", n), edakit.IntStr(ans),
			edakit.Sprintf("Each of its $\\log_2 %d = %d$ levels has %d nodes: $%d \\times %d = %d$.", n, L, n/2, n/2, L, ans), map[string]string{"levels": edakit.IntStr(L)})
	}
}

func ksNodes(n, L int) int { return n*L - n + 1 }
func bkNodes(n, L int) int { return 2*n - 2 - L }

func koggeBrent(ctx generator.GeneratorContext) generator.Problem {
	L := edakit.Scaled(ctx, 3, 7)
	n := 1 << uint(L)
	switch variant(ctx) {
	case 1:
		if edakit.Int(0, 2) == 0 {
			return edakit.Numeric(edakit.Sprintf("How many levels deep is a Kogge-Stone prefix network for %d bits?", n), edakit.IntStr(L),
				edakit.Sprintf("Kogge-Stone has minimum depth, $\\log_2 %d = %d$.", n, L), map[string]string{"nodes": edakit.IntStr(ksNodes(n, L))})
		}
		return edakit.Numeric(edakit.Sprintf("How many prefix nodes does a Kogge-Stone network for %d bits use?", n), edakit.IntStr(ksNodes(n, L)),
			edakit.Sprintf("Level l has $%d - 2^l$ nodes; summing over %d levels gives $%d \\cdot %d - %d + 1 = %d$.", n, L, n, L, n, ksNodes(n, L)), map[string]string{"depth": edakit.IntStr(L)})
	case 2:
		if edakit.Int(0, 2) == 0 {
			return edakit.Numeric(edakit.Sprintf("How many levels deep is a Brent-Kung prefix network for %d bits?", n), edakit.IntStr(2*L-1),
				edakit.Sprintf("%d levels of up-sweep plus %d of down-sweep: $2 \\cdot %d - 1 = %d$.", L, L-1, L, 2*L-1), map[string]string{"nodes": edakit.IntStr(bkNodes(n, L))})
		}
		return edakit.Numeric(edakit.Sprintf("How many prefix nodes does a Brent-Kung network for %d bits use?", n), edakit.IntStr(bkNodes(n, L)),
			edakit.Sprintf("Up-sweep %d nodes plus down-sweep %d nodes: $2 \\cdot %d - 2 - %d = %d$.", n-1, n-1-L, n, L, bkNodes(n, L)), map[string]string{"depth": edakit.IntStr(2*L - 1)})
	default:
		if edakit.Int(0, 1) == 0 {
			d := ksNodes(n, L) - bkNodes(n, L)
			return edakit.Numeric(edakit.Sprintf("For an adder of %d bits, how many more prefix nodes does Kogge-Stone use than Brent-Kung?", n), edakit.IntStr(d),
				edakit.Sprintf("Kogge-Stone: %d nodes; Brent-Kung: %d nodes; difference %d.", ksNodes(n, L), bkNodes(n, L), d), map[string]string{"ks": edakit.IntStr(ksNodes(n, L)), "bk": edakit.IntStr(bkNodes(n, L))})
		}
		t := 5 * edakit.Int(3, 8)
		d := (2*L - 1 - L) * t
		return edakit.Numeric(edakit.Sprintf("Each prefix node costs %d ps and wires are ignored. For an adder of %d bits, how much longer, in ps, is the Brent-Kung prefix network's delay than Kogge-Stone's?", t, n), edakit.IntStr(d),
			edakit.Sprintf("Brent-Kung has %d levels and Kogge-Stone %d, so the difference is $%d \\times %d = %d$ ps.", 2*L-1, L, L-1, t, d), map[string]string{"levels_gap": edakit.IntStr(L - 1)})
	}
}

func subtractCompare(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Pick([]int{8, 8, 12})
	mask := (uint64(1) << uint(n)) - 1
	switch variant(ctx) {
	case 1:
		a, b := operands(n)
		r := (a - b) & mask
		return edakit.Numeric(edakit.Sprintf("An adder-subtractor of width %d bits computes $a + \\bar{b} + 1$ for unsigned $a = %d$ and $b = %d$. What %d-bit result does it produce, as an unsigned integer?", n, a, b, n), edakit.IntStr(int(r)),
			edakit.Sprintf("$\\bar{b} = %d$, and $%d + %d + 1 = %d$, reduced modulo $2^{%d}$ to %d.", (^b)&mask, a, (^b)&mask, a+((^b)&mask)+1, n, r), map[string]string{"inverted_b": edakit.IntStr(int((^b) & mask))})
	case 2:
		a, b := operands(n)
		sum := a + ((^b) & mask) + 1
		co := int(sum >> uint(n) & 1)
		return edakit.Numeric(edakit.Sprintf("A subtractor of width %d bits computes $a + \\bar{b} + 1$ for unsigned $a = %d$ and $b = %d$. What is its carry-out (0 or 1)?", n, a, b), edakit.IntStr(co),
			edakit.Sprintf("The full sum is %d, which %s $2^{%d}$, so the carry-out is %d: this says $a %s b$.", sum, map[bool]string{true: "reaches", false: "is below"}[co == 1], n, co, map[bool]string{true: "\\ge", false: "<"}[co == 1]),
			map[string]string{"sum": edakit.IntStr(int(sum))})
	default:
		lim := 1 << uint(n-1)
		a := int64(edakit.Int(-lim, lim-1))
		b := int64(edakit.Int(-lim, lim-1))
		if edakit.Int(0, 1) == 0 {
			// force opposite signs so overflow is likely
			if (a >= 0) == (b >= 0) {
				b = -b - 1
				if b < int64(-lim) {
					b = int64(-lim)
				}
			}
		}
		diff := a - b
		w := wrap(diff, n)
		N := bit(w < 0)
		V := bit(w != diff)
		ans := N ^ V
		return edakit.Numeric(edakit.Sprintf("An ALU of width %d bits subtracts signed $b = %d$ from signed $a = %d$ and sets the sign flag N (top bit of the result) and the overflow flag V. What is $N \\oplus V$, the signed less-than result (0 or 1)?", n, b, a), edakit.IntStr(ans),
			edakit.Sprintf("The true difference is %d and the stored result is %d, so N = %d and V = %d. $N \\oplus V = %d$, and indeed $a %s b$.", diff, w, N, V, ans, map[int]string{1: "<", 0: "\\ge"}[ans]),
			map[string]string{"N": edakit.IntStr(N), "V": edakit.IntStr(V)})
	}
}

func multiplier(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 8)
		A := uint64(edakit.Int(3, (1<<uint(n))-1))
		B := uint64(edakit.Int(3, (1<<uint(n))-1))
		var ones []int
		for j := 0; j < n; j++ {
			if B>>uint(j)&1 == 1 {
				ones = append(ones, j)
			}
		}
		j := edakit.Pick(ones)
		row := A << uint(j)
		return edakit.Numeric(edakit.Sprintf("An unsigned multiplier forms partial products of $A$ = %s and $B$ = %s, one row per bit of B, with row j shifted left by j. What is the value of row %d, as an unsigned decimal integer?", bin(A, n), bin(B, n), j), edakit.IntStr(int(row)),
			edakit.Sprintf("$b_{%d} = 1$, so row %d is $A \\cdot 2^{%d} = %d \\cdot %d = %d$.", j, j, j, A, 1<<uint(j), row), map[string]string{"product": edakit.IntStr(int(A * B))})
	case 2:
		n, m := edakit.Int(4, 32), edakit.Int(4, 32)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("How many partial-product bits (AND gates) does an unsigned %d by %d bit multiplier form?", n, m), edakit.IntStr(n*m),
				edakit.Sprintf("Every pair $a_i b_j$ is one bit: $%d \\times %d = %d$.", n, m, n*m), map[string]string{"width": edakit.IntStr(n + m)})
		}
		return edakit.Numeric(edakit.Sprintf("How many bits wide must the product of an unsigned %d-bit number and an unsigned %d-bit number be to hold every possible product?", n, m), edakit.IntStr(n+m),
			edakit.Sprintf("$(2^{%d} - 1)(2^{%d} - 1) < 2^{%d}$, so %d bits suffice and fewer do not.", n, m, n+m, n+m), map[string]string{"pp_bits": edakit.IntStr(n * m)})
	default:
		n, m := edakit.Int(3, 10), edakit.Int(3, 10)
		if edakit.Int(0, 1) == 0 {
			ans := ((1 << uint(n)) - 1) * ((1 << uint(m)) - 1)
			return edakit.Numeric(edakit.Sprintf("What is the largest product an unsigned %d-bit by %d-bit multiplier can produce?", n, m), edakit.IntStr(ans),
				edakit.Sprintf("$(2^{%d} - 1)(2^{%d} - 1) = %d$.", n, m, ans), map[string]string{"width": edakit.IntStr(n + m)})
		}
		ans := 1 << uint(n+m-2)
		return edakit.Numeric(edakit.Sprintf("Two's-complement operands of %d and %d bits are multiplied. What is the largest positive product that can occur?", n, m), edakit.IntStr(ans),
			edakit.Sprintf("The extreme is $(-2^{%d})(-2^{%d}) = 2^{%d} = %d$, which needs all %d bits as a signed value.", n-1, m-1, n+m-2, ans, n+m), map[string]string{"width": edakit.IntStr(n + m)})
	}
}

// boothDigits recodes an n-bit two's-complement pattern into radix-4 digits, least significant first.
func boothDigits(y uint64, n int) []int {
	var d []int
	prev := 0
	for i := 0; i < n/2; i++ {
		b0 := int(y >> uint(2*i) & 1)
		b1 := int(y >> uint(2*i+1) & 1)
		d = append(d, -2*b1+b0+prev)
		prev = b1
	}
	return d
}

func booth(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		n := 2 * edakit.Scaled(ctx, 4, 32)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("How many partial products does radix-4 Booth recoding create for a two's-complement multiplier of %d bits?", n), edakit.IntStr(n/2),
				edakit.Sprintf("One per overlapping bit pair: %d / 2 = %d.", n, n/2), map[string]string{"radix2": edakit.IntStr(n)})
		}
		return edakit.Numeric(edakit.Sprintf("Compared with plain radix-2 partial products, how many fewer partial products does radix-4 Booth recoding create for a two's-complement multiplier of %d bits?", n), edakit.IntStr(n/2),
			edakit.Sprintf("Radix-2 creates %d rows and radix-4 Booth %d, a saving of %d.", n, n/2, n/2), map[string]string{"booth": edakit.IntStr(n / 2)})
	case 2:
		n := edakit.Pick([]int{8, 8, 12})
		y := uint64(edakit.Int(1, (1<<uint(n))-1))
		d := boothDigits(y, n)
		i := edakit.Int(0, n/2-1)
		var t string
		lo := 0
		if i > 0 {
			lo = int(y >> uint(2*i-1) & 1)
		}
		t = edakit.Sprintf("%d%d%d", y>>uint(2*i+1)&1, y>>uint(2*i)&1, lo)
		return edakit.Numeric(edakit.Sprintf("Recode the %d-bit two's-complement multiplier %s in radix-4 Booth form, $d_i = -2y_{2i+1} + y_{2i} + y_{2i-1}$ with $y_{-1} = 0$. What is digit $d_{%d}$?", n, groups4(bin(y, n)), i), edakit.IntStr(d[i]),
			edakit.Sprintf("Digit %d reads the triplet $y_{%d} y_{%d} y_{%d}$ = %s, giving %d. All digits, least significant first: %s.", i, 2*i+1, 2*i, 2*i-1, t, d[i], joinInts(d, ", ")),
			map[string]string{"value": i64(signed(y, n))})
	default:
		n := 8
		A := edakit.Int(3, 60)
		y := uint64(edakit.Int(1, 255))
		d := boothDigits(y, n)
		i := edakit.Int(0, 3)
		for d[i] == 0 {
			i = edakit.Int(0, 3)
			if d[0] == 0 && d[1] == 0 && d[2] == 0 && d[3] == 0 {
				y = uint64(edakit.Int(1, 255))
				d = boothDigits(y, n)
			}
		}
		return edakit.Numeric(edakit.Sprintf("A radix-4 Booth multiplier multiplies $A = %d$ by the 8-bit two's-complement multiplier %s ($d_i = -2y_{2i+1} + y_{2i} + y_{2i-1}$, $y_{-1} = 0$). Before shifting, what multiple of A does digit $d_{%d}$ select? Enter the signed value $d_{%d} \\cdot A$.", A, bin(y, n), i, i), edakit.IntStr(d[i]*A),
			edakit.Sprintf("The digits are %s (least significant first), so $d_{%d} = %d$ and the row is $%d \\cdot %d = %d$, later shifted by %d bits.", joinInts(d, ", "), i, d[i], d[i], A, d[i]*A, 2*i),
			map[string]string{"digit": edakit.IntStr(d[i])})
	}
}

// daddaLevels is the number of 3:2 reduction levels to bring a column height of r down to 2.
func daddaLevels(r int) int {
	lv, d := 0, 2
	for d < r {
		lv++
		d = d * 3 / 2
	}
	return lv
}

func daddaSeq(r int) string {
	var xs []int
	for d := 2; d < r; d = d * 3 / 2 {
		xs = append(xs, d)
	}
	return joinInts(xs, ", ")
}

func wallaceDadda(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		r := edakit.Scaled(ctx, 3, 64)
		lv := daddaLevels(r)
		return edakit.Numeric(edakit.Sprintf("A partial-product matrix has maximum column height %d. How many levels of 3:2 (full-adder) reduction does a Wallace or Dadda tree need to bring it to height 2?", r), edakit.IntStr(lv),
			edakit.Sprintf("The Dadda heights below %d are %s: %d levels.", r, daddaSeq(r), lv), map[string]string{"levels": edakit.IntStr(lv)})
	case 2:
		n := edakit.Scaled(ctx, 4, 32)
		if edakit.Int(0, 1) == 0 {
			fa := n*n - 4*n + 3
			return edakit.Numeric(edakit.Sprintf("How many full adders does a Dadda tree use to reduce the partial products of an unsigned %d by %d bit multiplier?", n, n), edakit.IntStr(fa),
				edakit.Sprintf("$n^2 - 4n + 3 = %d - %d + 3 = %d$.", n*n, 4*n, fa), map[string]string{"half_adders": edakit.IntStr(n - 1)})
		}
		return edakit.Numeric(edakit.Sprintf("How many half adders does a Dadda tree use to reduce the partial products of an unsigned %d by %d bit multiplier?", n, n), edakit.IntStr(n-1),
			edakit.Sprintf("A Dadda tree uses $n - 1 = %d$ half adders.", n-1), map[string]string{"full_adders": edakit.IntStr(n*n - 4*n + 3)})
	default:
		n := 2 * edakit.Scaled(ctx, 6, 32)
		before, after := daddaLevels(n), daddaLevels(n/2)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("A multiplier of %d by %d bits uses radix-4 Booth recoding, and its reduction tree is treated as starting from %d rows. How many 3:2 reduction levels does the tree need?", n, n, n/2), edakit.IntStr(after),
				edakit.Sprintf("The Dadda heights below %d are %s: %d levels.", n/2, daddaSeq(n/2), after), map[string]string{"without_booth": edakit.IntStr(before)})
		}
		return edakit.Numeric(edakit.Sprintf("For a multiplier of %d by %d bits, radix-4 Booth recoding reduces the starting height of the reduction tree from %d rows to %d. How many 3:2 reduction levels does that save?", n, n, n, n/2), edakit.IntStr(before-after),
			edakit.Sprintf("%d rows need %d levels and %d rows need %d, saving %d.", n, before, n/2, after, before-after), map[string]string{"before": edakit.IntStr(before), "after": edakit.IntStr(after)})
	}
}
