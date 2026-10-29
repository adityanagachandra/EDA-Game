package dft

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// scanCycles is the tester cycle count for p patterns through a chain of length l with
// overlapped unload and load and one capture cycle per pattern.
func scanCycles(p, l int) int { return p*(l+1) + l }

// scanChain: shift cycles for one load (KP 1), total cycles for a pattern set (KP 2), test time
// (KP 3).
func scanChain(ctx generator.GeneratorContext) generator.Problem {
	l := edakit.Scaled(ctx, 50, 4000)
	p := edakit.Scaled(ctx, 20, 3000)
	switch kp(ctx) {
	case 1:
		f := edakit.Scaled(ctx, 200, 20000)
		ch := edakit.Pick([]int{1, 2, 4, 8, 16})
		ll := int(math.Ceil(float64(f) / float64(ch)))
		return edakit.Numeric(edakit.Sprintf("A design has %d mux-D scan flip-flops stitched into %d equal-as-possible chains. With scan enable held at 1, how many shift clock cycles does it take to load one complete pattern into every chain?", f, ch),
			edakit.IntStr(ll), edakit.Sprintf("Chains shift in parallel, so the load takes as many cycles as the longest chain: $\\lceil %d/%d \\rceil = %d$.", f, ch, ll), map[string]string{"chain_length": edakit.IntStr(ll)})
	case 2:
		n := scanCycles(p, l)
		return edakit.Numeric(edakit.Sprintf("A scan test has %d patterns on chains whose longest is %d flip-flops. The unload of each response overlaps the load of the next pattern, each pattern needs one capture cycle, and the last response must be unloaded at the end. How many tester cycles does the scan test take?", p, l),
			edakit.IntStr(n), edakit.Sprintf("$p(L+1) + L = %d(%d+1) + %d = %d$ cycles.", p, l, l, n), map[string]string{"cycles": edakit.IntStr(n)})
	default:
		fr := edakit.Pick([]int{10, 20, 25, 50, 100})
		n := scanCycles(p, l)
		ms := float64(n) / (float64(fr) * 1000)
		ans := edakit.Num(ms, 2)
		return edakit.Numeric(edakit.Sprintf("A scan test has %d patterns, a longest chain of %d flip-flops and one capture cycle per pattern, with unload overlapped with the next load. The tester shifts at %d MHz and captures at the same rate. How long is the scan test, in milliseconds? Round to 2 decimal places.", p, l, fr),
			ans, edakit.Sprintf("Cycles = $%d(%d+1) + %d = %d$; at %d MHz that is %s ms. Shift dominates, which is why chain length is the main test-time lever.", p, l, l, n, fr, ans), map[string]string{"cycles": edakit.IntStr(n)})
	}
}

// scanInsertion: longest balanced chain (KP 1), cycles saved by rebalancing (KP 2), chains needed
// for a cycle budget (KP 3).
func scanInsertion(ctx generator.GeneratorContext) generator.Problem {
	f := edakit.Scaled(ctx, 500, 50000)
	p := edakit.Scaled(ctx, 100, 3000)
	switch kp(ctx) {
	case 1:
		c := edakit.Int(3, 40)
		l := int(math.Ceil(float64(f) / float64(c)))
		return edakit.Numeric(edakit.Sprintf("Scan insertion must stitch %d scan flip-flops into %d chains as evenly as possible. How long is the longest chain?", f, c),
			edakit.IntStr(l), edakit.Sprintf("Balanced chains differ by at most one flop: $\\lceil %d/%d \\rceil = %d$.", f, c, l), map[string]string{"length": edakit.IntStr(l)})
	case 2:
		c := edakit.Int(3, 6)
		lens := make([]int, c)
		tot := 0
		mx := 0
		for i := range lens {
			lens[i] = edakit.Int(100, 900)
			tot += lens[i]
			if lens[i] > mx {
				mx = lens[i]
			}
		}
		bal := int(math.Ceil(float64(tot) / float64(c)))
		var t []string
		for _, x := range lens {
			t = append(t, edakit.IntStr(x))
		}
		saved := scanCycles(p, mx) - scanCycles(p, bal)
		return edakit.Numeric(edakit.Sprintf("Clock-domain-driven insertion produced %d chains of lengths %s. A test of %d patterns uses overlapped load and unload with one capture cycle per pattern. How many tester cycles would rebalancing the same flops into %d equal-as-possible chains save?", c, strings.Join(t, ", "), p, c),
			edakit.IntStr(saved), edakit.Sprintf("Test time follows the longest chain. Unbalanced: L = %d gives %d cycles; balanced: L = %d gives %d; saved %d.", mx, scanCycles(p, mx), bal, scanCycles(p, bal), saved), map[string]string{"balanced": edakit.IntStr(bal)})
	default:
		budget := scanCycles(p, int(math.Ceil(float64(f)/float64(edakit.Int(4, 64))))) + edakit.Int(0, 500)
		c := 1
		for scanCycles(p, int(math.Ceil(float64(f)/float64(c)))) > budget {
			c++
		}
		l := int(math.Ceil(float64(f) / float64(c)))
		return edakit.Numeric(edakit.Sprintf("A design has %d scan flops and a test of %d patterns (overlapped load and unload, one capture cycle per pattern, balanced chains). The tester budget is %d cycles. What is the smallest number of chains that meets the budget?", f, p, budget),
			edakit.IntStr(c), edakit.Sprintf("Need $p(L+1)+L \\le %d$ with $L = \\lceil %d/C \\rceil$. C = %d gives L = %d and %d cycles; one fewer chain exceeds the budget.", budget, f, c, l, scanCycles(p, l)), map[string]string{"length": edakit.IntStr(l)})
	}
}

// scanOrdering: wirelength of a given order (KP 1), nearest-neighbour order (KP 2), optimal order
// by enumeration (KP 3).
func scanOrdering(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	type pt struct{ x, y int }
	si := pt{0, 0}
	so := pt{edakit.Int(5, 12), edakit.Int(5, 12)}
	ps := make([]pt, n)
	var grid []pt
	for x := 0; x <= 12; x++ {
		for y := 0; y <= 12; y++ {
			if q := (pt{x, y}); q != si && q != so {
				grid = append(grid, q)
			}
		}
	}
	grid = edakit.Shuffle(grid)
	copy(ps, grid[:n])
	d := func(a, b pt) int { return absInt(a.x-b.x) + absInt(a.y-b.y) }
	wl := func(order []int) int {
		s := d(si, ps[order[0]])
		for i := 1; i < len(order); i++ {
			s += d(ps[order[i-1]], ps[order[i]])
		}
		return s + d(ps[order[len(order)-1]], so)
	}
	var desc []string
	for i, q := range ps {
		desc = append(desc, edakit.Sprintf("F%d (%d, %d)", i+1, q.x, q.y))
	}
	stem := edakit.Sprintf("Scan-in is at (0, 0) and scan-out at (%d, %d). The flops are %s. Wirelength is the Manhattan distance summed along the chain from scan-in, through every flop, to scan-out.", so.x, so.y, strings.Join(desc, ", "))
	switch kp(ctx) {
	case 1:
		order := edakit.Shuffle(seq(n))
		var o []string
		for _, i := range order {
			o = append(o, edakit.Sprintf("F%d", i+1))
		}
		w := wl(order)
		return edakit.Numeric(stem+" What is the wirelength of the order "+strings.Join(o, ", ")+"?",
			edakit.IntStr(w), edakit.Sprintf("Summing the Manhattan hops gives %d.", w), nil)
	case 2:
		used := make([]bool, n)
		cur := si
		var order []int
		for len(order) < n {
			best, bd := -1, 1<<30
			for i := 0; i < n; i++ {
				if !used[i] && d(cur, ps[i]) < bd {
					best, bd = i, d(cur, ps[i])
				}
			}
			used[best] = true
			order = append(order, best)
			cur = ps[best]
		}
		var o []string
		for _, i := range order {
			o = append(o, edakit.Sprintf("F%d", i+1))
		}
		w := wl(order)
		return edakit.Numeric(stem+" A greedy orderer starts at scan-in and always moves to the nearest unvisited flop (ties to the lower-numbered flop), then ends at scan-out. What wirelength does it produce?",
			edakit.IntStr(w), "Greedy order: "+strings.Join(o, ", ")+edakit.Sprintf(", wirelength %d.", w), nil)
	default:
		best := 1 << 30
		var bo []int
		permute(seq(n), 0, func(p []int) {
			if w := wl(p); w < best {
				best = w
				bo = append([]int(nil), p...)
			}
		})
		var o []string
		for _, i := range bo {
			o = append(o, edakit.Sprintf("F%d", i+1))
		}
		return edakit.Numeric(stem+" What is the minimum wirelength over all orders of the flops?",
			edakit.IntStr(best), "An optimal order is "+strings.Join(o, ", ")+edakit.Sprintf(", wirelength %d. This is a path version of the travelling-salesman problem, so real tools use heuristics within each placement region.", best), nil)
	}
}

func permute(a []int, k int, f func([]int)) {
	if k == len(a) {
		f(a)
		return
	}
	for i := k; i < len(a); i++ {
		a[k], a[i] = a[i], a[k]
		permute(a, k+1, f)
		a[k], a[i] = a[i], a[k]
	}
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// locLosCalc: the launch state under LOS (KP 1) or LOC (KP 2) for a three- or four-flop design.
func locLosCalc(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 4)
	type fn struct {
		txt string
		f   func(q []int) int
	}
	gen := func(i int) fn {
		a := edakit.Int(0, n-1)
		b := (a + edakit.Int(1, n-1)) % n
		switch edakit.Int(0, 4) {
		case 0:
			return fn{edakit.Sprintf("$Q_%d \\oplus Q_%d$", a+1, b+1), func(q []int) int { return q[a] ^ q[b] }}
		case 1:
			return fn{edakit.Sprintf("$Q_%d Q_%d$", a+1, b+1), func(q []int) int { return q[a] & q[b] }}
		case 2:
			return fn{edakit.Sprintf("$Q_%d + Q_%d$", a+1, b+1), func(q []int) int { return q[a] | q[b] }}
		case 3:
			return fn{edakit.Sprintf("$\\overline{Q_%d}$", a+1), func(q []int) int { return 1 - q[a] }}
		}
		return fn{edakit.Sprintf("$\\overline{Q_%d Q_%d}$", a+1, b+1), func(q []int) int { return 1 - q[a]&q[b] }}
	}
	fs := make([]fn, n)
	var ftxt []string
	for i := range fs {
		fs[i] = gen(i)
		ftxt = append(ftxt, edakit.Sprintf("$D_%d$ = %s", i+1, fs[i].txt))
	}
	v1 := make([]int, n)
	for i := range v1 {
		v1[i] = edakit.Int(0, 1)
	}
	sin := edakit.Int(0, 1)
	los := make([]int, n)
	los[0] = sin
	for i := 1; i < n; i++ {
		los[i] = v1[i-1]
	}
	loc := make([]int, n)
	for i := range loc {
		loc[i] = fs[i].f(v1)
	}
	rev := make([]int, n)
	for i := 0; i < n-1; i++ {
		rev[i] = v1[i+1]
	}
	rev[n-1] = sin
	inv := make([]int, n)
	for i := range inv {
		inv[i] = 1 - v1[i]
	}
	st := func(q []int) string { return bits(q) }
	stem := edakit.Sprintf("A scan chain runs scan-in → $Q_1$ → ... → $Q_%d$ → scan-out. In functional mode %s. The last shift leaves $V_1 = Q_1 \\ldots Q_%d$ = %s, and the scan-in pin holds %d.", n, strings.Join(ftxt, ", "), n, st(v1), sin)
	if kp(ctx) == 1 {
		ans := st(los)
		ds := uniq(ans, st(loc), st(rev), st(inv), st(v1))
		return edakit.Choice(stem+edakit.Sprintf(" Under launch-on-shift, what launch state $V_2 = Q_1 \\ldots Q_%d$ does the final at-speed shift create?", n),
			ans, "Launch-on-shift makes the second vector by shifting once more: $Q_1$ takes the scan-in bit and each other flop takes its predecessor's value.", ds...)
	}
	ans := st(loc)
	ds := uniq(ans, st(los), st(rev), st(inv), st(v1))
	return edakit.Choice(stem+edakit.Sprintf(" Under launch-on-capture, with scan enable dropped, what launch state $V_2 = Q_1 \\ldots Q_%d$ does the first functional clock create?", n),
		ans, "Launch-on-capture computes the second vector through the functional logic: each flop captures its next-state function of $V_1$.", ds...)
}

// uniq returns up to 3 distinct values not equal to ans.
func uniq(ans string, cands ...string) []string {
	seen := map[string]bool{ans: true}
	var out []string
	for _, c := range cands {
		if !seen[c] && len(out) < 3 {
			seen[c] = true
			out = append(out, c)
		}
	}
	for attempt := 0; len(out) < 3; attempt++ {
		if attempt >= generator.MaxDrawAttempts {
			for k := 0; k < 1<<len(ans) && len(out) < 3; k++ {
				s := edakit.Sprintf("%0*b", len(ans), k)
				if !seen[s] {
					seen[s] = true
					out = append(out, s)
				}
			}
			break
		}
		// flip bits of ans to make distinct fallbacks
		b := []byte(ans)
		i := edakit.Int(0, len(b)-1)
		if b[i] == '0' {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
		j := edakit.Int(0, len(b)-1)
		if j != i && edakit.Int(0, 1) == 1 {
			if b[j] == '0' {
				b[j] = '1'
			} else {
				b[j] = '0'
			}
		}
		if s := string(b); !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// scanCompression: shift-cycle reduction (KP 1), data-volume ratio (KP 2), effective ratio with
// pattern inflation (KP 3).
func scanCompression(ctx generator.GeneratorContext) generator.Problem {
	f := edakit.Scaled(ctx, 20000, 2000000)
	ch := edakit.Pick([]int{4, 8, 16, 32})
	nint := ch * edakit.Pick([]int{10, 20, 25, 50, 100})
	lu := int(math.Ceil(float64(f) / float64(ch)))
	lc := int(math.Ceil(float64(f) / float64(nint)))
	stem := edakit.Sprintf("A design has %d scan flops and %d scan-in channels on the tester. Without compression the channels drive %d chains directly; with an EDT-style decompressor they feed %d internal chains.", f, ch, ch, nint)
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(float64(lu)/float64(lc), 2)
		return edakit.Numeric(stem+" By what factor does the shift length per pattern fall? Round to 2 decimal places.",
			ans, edakit.Sprintf("Uncompressed chains are $\\lceil %d/%d \\rceil = %d$ long, internal chains $\\lceil %d/%d \\rceil = %d$; the ratio is %s.", f, ch, lu, f, nint, lc, ans), map[string]string{"lu": edakit.IntStr(lu), "lc": edakit.IntStr(lc)})
	case 2:
		p := edakit.Scaled(ctx, 1000, 20000)
		vu := float64(p) * float64(ch) * float64(lu)
		vc := float64(p) * float64(ch) * float64(lc)
		ans := edakit.Num(vu/vc, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" For %d patterns, the stimulus volume is patterns times channels times shift cycles. What is the stimulus compression ratio, uncompressed volume over compressed volume? Round to 2 decimal places.", p),
			ans, edakit.Sprintf("Uncompressed: $%d \\times %d \\times %d$ bits; compressed: $%d \\times %d \\times %d$. The ratio is %s, close to the chain ratio %d/%d.", p, ch, lu, p, ch, lc, ans, nint, ch), nil)
	default:
		p := edakit.Scaled(ctx, 1000, 20000)
		infl := edakit.Int(5, 40)
		pc := int(math.Ceil(float64(p) * (1 + float64(infl)/100)))
		vu := float64(p) * float64(ch) * float64(lu)
		vc := float64(pc) * float64(ch) * float64(lc)
		ans := edakit.Num(vu/vc, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" Uncompressed ATPG needs %d patterns; compressed ATPG needs %d percent more (%d) because some care-bit combinations cannot be encoded. Volume is patterns times channels times shift cycles. What is the effective compression ratio? Round to 2 decimal places.", p, infl, pc),
			ans, edakit.Sprintf("$(%d \\cdot %d \\cdot %d)/(%d \\cdot %d \\cdot %d) = %s$. Pattern inflation erodes the raw chain ratio of %d.", p, ch, lu, pc, ch, lc, ans, nint/ch), nil)
	}
}

// decompressor: evaluate a cell (KP 1), count consistent seeds (KP 2), encodability probability
// (KP 3).
func decompressor(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k == 3 {
		nv := edakit.Scaled(ctx, 16, 64)
		s := nv - edakit.Int(0, 12)
		p := 1.0
		for i := 0; i < s; i++ {
			p *= 1 - math.Pow(2, float64(i-nv))
		}
		ans := edakit.Num(p, 4)
		return edakit.Numeric(edakit.Sprintf("A linear decompressor injects %d free variables per pattern; each care bit is one linear equation over GF(2) in those variables. Treat the coefficient vectors of %d care bits as independent and uniformly random. What is the probability that they are linearly independent, which guarantees the pattern is encodable whatever the care values? Round to 4 decimal places.", nv, s),
			ans, edakit.Sprintf("$P = \\prod_{i=0}^{%d} (1 - 2^{i-%d}) = %s$. Near $s = n$ the last factors approach 1/2, which is why encoders keep care bits a margin below the variable count.", s-1, nv, ans), nil)
	}
	nv := edakit.Scaled(ctx, 3, 5)
	nc := edakit.Scaled(ctx, 4, 7)
	eqs := make([]int, nc)
	var et []string
	for i := range eqs {
		eqs[i] = edakit.Int(1, (1<<nv)-1)
		et = append(et, edakit.Sprintf("$s_%d = %s$", i+1, xorText(eqs[i], nv)))
	}
	stem := "A linear decompressor drives scan cells with XORs of the free variables: " + strings.Join(et, ", ") + "."
	if k == 1 {
		x := edakit.Int(0, (1<<nv)-1)
		var vt []string
		for j := 0; j < nv; j++ {
			vt = append(vt, edakit.Sprintf("$x_%d = %d$", j+1, x>>j&1))
		}
		out := make([]int, nc)
		for i := range out {
			out[i] = popcount(eqs[i]&x) & 1
		}
		ans := bits(out)
		return edakit.Numeric(stem+" With "+strings.Join(vt, ", ")+edakit.Sprintf(", what values do $s_1 \\ldots s_%d$ receive? Answer with the bits in order.", nc),
			ans, "Each cell is the parity of the variables it lists: "+ans+".", map[string]string{"cells": ans})
	}
	// KP 2: care bits on a subset of cells.
	m := edakit.Int(2, minInt(nc, nv+1))
	cells := edakit.Shuffle(seq(nc))[:m]
	var rows [][2]int
	var ct []string
	for _, ci := range cells {
		v := edakit.Int(0, 1)
		rows = append(rows, [2]int{eqs[ci], v})
		ct = append(ct, edakit.Sprintf("$s_%d = %d$", ci+1, v))
	}
	rank, ok := gf2Solve(rows, nv)
	cnt := 0
	if ok {
		cnt = 1 << (nv - rank)
	}
	why := edakit.Sprintf("Gaussian elimination over GF(2) gives rank %d", rank)
	if ok {
		why += edakit.Sprintf(" with no contradiction, so $2^{%d - %d} = %d$ assignments satisfy the care bits.", nv, rank, cnt)
	} else {
		why += " and a contradiction (0 = 1), so the pattern cannot be encoded: 0 assignments."
	}
	return edakit.Numeric(stem+" A test cube needs "+strings.Join(ct, ", ")+edakit.Sprintf(" and leaves the other cells as don't-cares. How many of the $2^{%d}$ variable assignments produce it?", nv),
		edakit.IntStr(cnt), why, map[string]string{"rank": edakit.IntStr(rank)})
}

func xorText(mask, nv int) string {
	var p []string
	for j := 0; j < nv; j++ {
		if mask>>j&1 == 1 {
			p = append(p, edakit.Sprintf("x_%d", j+1))
		}
	}
	return strings.Join(p, " \\oplus ")
}

// gf2Solve returns the rank of the system and whether it is consistent.
func gf2Solve(rows [][2]int, nv int) (int, bool) {
	r := append([][2]int(nil), rows...)
	rank := 0
	for col := 0; col < nv; col++ {
		piv := -1
		for i := rank; i < len(r); i++ {
			if r[i][0]>>col&1 == 1 {
				piv = i
				break
			}
		}
		if piv < 0 {
			continue
		}
		r[rank], r[piv] = r[piv], r[rank]
		for i := range r {
			if i != rank && r[i][0]>>col&1 == 1 {
				r[i][0] ^= r[rank][0]
				r[i][1] ^= r[rank][1]
			}
		}
		rank++
	}
	for i := rank; i < len(r); i++ {
		if r[i][0] == 0 && r[i][1] == 1 {
			return rank, false
		}
	}
	return rank, true
}

// compactor: XOR outputs (KP 1), aliasing double errors (KP 2), observability under X (KP 3).
func compactor(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Scaled(ctx, 6, 12)
	q := edakit.Int(2, 4)
	group := make([]int, m)
	for i := range group {
		group[i] = i % q
	}
	group = edakit.Shuffle(group)
	var gt []string
	for o := 0; o < q; o++ {
		var mem []string
		for i := 0; i < m; i++ {
			if group[i] == o {
				mem = append(mem, edakit.Sprintf("c%d", i+1))
			}
		}
		gt = append(gt, edakit.Sprintf("$z_%d$ = XOR(%s)", o+1, strings.Join(mem, ", ")))
	}
	stem := edakit.Sprintf("A space compactor folds %d scan-chain outputs c1 to c%d into %d tester outputs: %s.", m, m, q, strings.Join(gt, "; "))
	switch kp(ctx) {
	case 1:
		v := make([]int, m)
		for i := range v {
			v[i] = edakit.Int(0, 1)
		}
		z := make([]int, q)
		for i := range v {
			z[group[i]] ^= v[i]
		}
		ans := bits(z)
		return edakit.Numeric(stem+edakit.Sprintf(" In one shift cycle the chains present c1 to c%d = %s. What do $z_1 \\ldots z_%d$ show? Answer with the bits in order.", m, bits(v), q),
			ans, "Each output is the parity of its group: "+ans+".", nil)
	case 2:
		al := 0
		sizes := make([]int, q)
		for _, g := range group {
			sizes[g]++
		}
		var st []string
		for _, s := range sizes {
			al += s * (s - 1) / 2
			st = append(st, edakit.IntStr(s))
		}
		return edakit.Numeric(stem+" Exactly two chains capture an error in the same shift cycle. Of all pairs of chains, how many pairs cancel each other so that no output shows the error?",
			edakit.IntStr(al), edakit.Sprintf("Two errors cancel only inside the same XOR group. Group sizes %s give $\\sum \\binom{g}{2} = %d$ of the $\\binom{%d}{2} = %d$ pairs.", strings.Join(st, ", "), al, m, m*(m-1)/2), nil)
	default:
		nx := edakit.Int(1, 2)
		xs := edakit.Shuffle(seq(m))[:nx]
		bad := map[int]bool{}
		var xt []string
		for _, x := range xs {
			bad[group[x]] = true
			xt = append(xt, edakit.Sprintf("c%d", x+1))
		}
		obs := 0
		for i := 0; i < m; i++ {
			if !bad[group[i]] {
				obs++
			}
		}
		return edakit.Numeric(stem+edakit.Sprintf(" In one shift cycle chain(s) %s carry an unknown X value and no masking is applied. How many of the %d chains still have an observable bit in that cycle?", strings.Join(xt, " and "), m),
			edakit.IntStr(obs), edakit.Sprintf("An X entering an XOR makes that output X, hiding every chain in the group. %d chains remain observable; X-masking logic gates the X chain to a constant to recover the rest.", obs), nil)
	}
}

var locLosItems = []edakit.Item{
	{KP: 3, Q: "Launch-on-shift usually reaches higher transition-fault coverage than launch-on-capture. Why?",
		A: "The second vector is a free shift of the first, so ATPG controls both vectors almost independently", D: []string{"LOS uses a slower clock", "LOS needs no scan enable", "LOC cannot launch transitions at all"},
		Why: "Under LOC the second vector is the circuit's own response to the first, which restricts which transitions can be launched."},
	{KP: 3, Q: "What design requirement does launch-on-shift impose that launch-on-capture does not?",
		A: "Scan enable must switch from shift to capture within one at-speed cycle, so it needs timing closure like a clock", D: []string{"A second scan chain", "Lockup latches between every flop", "A separate test clock for each flop"},
		Why: "In LOS the last shift is the launch edge and the next edge is capture, so scan enable must be distributed as a fast, timed signal."},
	{KP: 3, Q: "Why are some faults detected by LOS considered over-testing?",
		A: "The launch state can be a state the functional circuit never reaches, so a failure may never occur in the field", D: []string{"LOS tests run at half speed", "LOS detects only stuck-at faults", "LOS patterns are never compressed"},
		Why: "A shifted state is not constrained by the logic, so it can exercise functionally unreachable transitions and reject good parts."},
	{KP: 3, Q: "A design team cannot close timing on a fast scan-enable tree. Which delay-test scheme should they use?",
		A: "Launch-on-capture, which keeps scan enable static during the launch and capture cycles", D: []string{"Launch-on-shift", "Stuck-at only", "Boundary scan EXTEST"},
		Why: "LOC drops scan enable before the launch and capture pulses, so scan enable can be slow."},
	{KP: 3, Q: "In a launch-on-capture test, which clocks run at the functional speed?",
		A: "The launch and capture pulses; the shift can run slowly", D: []string{"Only the shift clocks", "Every clock, including shift", "None: LOC is a static test"},
		Why: "Only the interval between launch and capture must equal the functional period; shift speed is limited by power and scan-path timing."},
	{KP: 3, Q: "Why does on-chip clock control (an OCC) matter for at-speed scan tests?",
		A: "It generates the precise fast launch and capture pulses from a PLL, which the tester cannot deliver directly", D: []string{"It removes the need for scan chains", "It slows the functional clock for shift", "It stores the patterns on chip"},
		Why: "Testers drive slow clocks; an OCC gates a burst of PLL pulses so the launch-to-capture interval is the real functional period."},
	{KP: 3, Q: "Transition-fault ATPG with LOC is described as a two-time-frame problem. What does that mean?",
		A: "The circuit is unrolled into two copies, with the second frame's state equal to the first frame's next state", D: []string{"Two separate stuck-at tests are generated", "Two scan chains are used", "The fault is simulated twice in the same frame"},
		Why: "ATPG must find $V_1$ such that $V_1$ sets the initial value and its functional response $V_2$ launches and propagates the transition."},
	{KP: 3, Q: "A slow-to-rise fault at a flop's next-state input requires $V_1$ to set it to 0 and $V_2$ to set it to 1. Under LOC, what decides $V_2$?",
		A: "The combinational logic's response to $V_1$", D: []string{"The scan-in pin", "The tester's second vector, freely chosen", "The scan-enable timing"},
		Why: "With LOC the second vector is computed by the functional logic from the first, so ATPG must choose $V_1$ to make that response right."},
	{KP: 3, Q: "Why do many flows combine LOC with a small number of LOS patterns or other enhancements?",
		A: "To pick up transition faults that LOC alone cannot launch, where the extra design cost is acceptable", D: []string{"Because LOC cannot test any transition fault", "To reduce stuck-at coverage", "To avoid using at-speed clocks"},
		Why: "Each scheme leaves faults the other detects; a hybrid trades scan-enable effort against coverage."},
	{KP: 3, Q: "Which faults do transition patterns, under either launch scheme, still miss that path-delay tests target?",
		A: "Small distributed delays that add up along a long path without any single gate being slow", D: []string{"Gross stuck-at faults", "Bridges between power and ground", "Faults inside scan flip-flops only"},
		Why: "Transition tests launch along whatever path ATPG finds, often short; small-delay defects need the longest paths, which path-delay or timing-aware ATPG targets."},
	{KP: 3, Q: "During the capture cycle of a LOC test, what must be true of primary inputs on a low-cost tester?",
		A: "They are usually held constant, because the tester cannot change them at speed", D: []string{"They must toggle at speed every cycle", "They must be driven by the scan chain", "They must be tristated"},
		Why: "Low-cost ATE cannot switch inputs at functional speed, so ATPG assumes static primary inputs across launch and capture."},
	{KP: 3, Q: "Why must the final shift in LOS be timed relative to the capture clock exactly like a functional cycle?",
		A: "That shift is the launch event, so the shift-to-capture interval is the path delay being tested", D: []string{"Because shift and capture use different flops", "Because the scan-in pin is slow", "It need not be: LOS is a slow test"},
		Why: "The at-speed interval in LOS runs from the last shift edge to the capture edge; any extra time would let slow paths pass."},
}
