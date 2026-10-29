package fv

import (
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// ---------- bit-vectors ----------

func signedVal(x, n int) int {
	if x>>(n-1)&1 == 1 {
		return x - (1 << n)
	}
	return x
}

func bitvectors(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	n := edakit.Pick([]int{4, 8})
	mod := 1 << n
	mask := mod - 1
	a, b := rand.Intn(mod), rand.Intn(mod)
	bin := func(x int) string { return "#b" + bitsN(x, n) }
	switch kp {
	case 1:
		k := edakit.Int(1, n-1)
		type op struct {
			expr string
			val  int
		}
		ops := []op{
			{edakit.Sprintf("(bvadd a b)"), (a + b) & mask},
			{edakit.Sprintf("(bvsub a b)"), (a - b + mod) & mask},
			{edakit.Sprintf("(bvmul a b)"), (a * b) & mask},
			{edakit.Sprintf("(bvshl a %s)", bin(k)), (a << k) & mask},
			{edakit.Sprintf("(bvlshr a %s)", bin(k)), a >> k},
			{edakit.Sprintf("(bvashr a %s)", bin(k)), (signedVal(a, n) >> k) & mask},
			{edakit.Sprintf("(bvxor a b)"), a ^ b},
			{edakit.Sprintf("((_ extract %d %d) a)", n-1, k), a >> k},
			{edakit.Sprintf("(bvadd (bvshl a %s) b)", bin(1)), ((a << 1) + b) & mask},
		}
		o := ops[rand.Intn(len(ops))]
		q := edakit.Sprintf("In SMT-LIB, a and b are %d-bit vectors with a = %s (%d) and b = %s (%d). What is the value of %s, read as an unsigned integer?", n, bin(a), a, bin(b), b, o.expr)
		return edakit.Numeric(q, edakit.IntStr(o.val),
			edakit.Sprintf("Bit-vector arithmetic wraps modulo $2^{%d} = %d$; shifts drop bits off the end (bvashr copies the sign bit in); extract keeps the named bits. The result is %d.", n, mod, o.val),
			map[string]string{"value": edakit.IntStr(o.val)})
	case 2:
		w := edakit.Scaled(ctx, 2, 64)
		if rand.Intn(2) == 0 {
			ans := 7 + 17*(w-1)
			q := edakit.Sprintf("A solver bit-blasts a %d-bit bvadd as a ripple-carry adder: a half adder at bit 0 (one XOR, one AND) and full adders elsewhere (two XOR, two AND, one OR). With Tseitin clauses (2-input XOR: 4, 2-input AND or OR: 3), how many clauses does the adder take?", w)
			return edakit.Numeric(q, edakit.IntStr(ans),
				edakit.Sprintf("Half adder: 4 + 3 = 7. Each full adder: 2·4 + 2·3 + 3 = 17. Total 7 + 17·%d = %d.", w-1, ans),
				map[string]string{"clauses": edakit.IntStr(ans)})
		}
		ans := 2*w + 2 + 5*(w-1)
		q := edakit.Sprintf("A solver bit-blasts a %d-bit bvadd as a ripple-carry adder: a half adder at bit 0 (one XOR, one AND) and full adders elsewhere (two XOR, two AND, one OR). Counting one variable per input bit of both operands and one per gate, how many variables does it use?", w)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Inputs: 2·%d = %d. Gates: 2 + 5·%d = %d. Total %d.", w, 2*w, w-1, 2+5*(w-1), ans),
			map[string]string{"vars": edakit.IntStr(ans)})
	}
	var expr string
	var r int
	switch rand.Intn(3) {
	case 0:
		expr, r = "(bvadd a b)", (a+b)&mask
	case 1:
		expr, r = "(bvsub a b)", (a-b+mod)&mask
	default:
		expr, r = "(bvmul a b)", (a*b)&mask
	}
	q := edakit.Sprintf("a and b are %d-bit vectors with a = %s and b = %s. What is the value of %s read as a two's-complement signed integer?", n, bin(a), bin(b), expr)
	return edakit.Numeric(q, edakit.IntStr(signedVal(r, n)),
		edakit.Sprintf("The %d-bit result pattern is %s. Reading the top bit with weight $-2^{%d}$ gives the signed value %d (unsigned %d).", n, bitsN(r, n), n-1, signedVal(r, n), r),
		map[string]string{"pattern": bitsN(r, n)})
}

// ---------- EUF and congruence closure ----------

type term struct {
	name string
	arg  int // index of argument term, -1 for a constant
}

func euf(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 3 {
		m := edakit.Scaled(ctx, 3, 20)
		ans := m * (m - 1) / 2
		q := edakit.Sprintf("Ackermann's reduction removes an uninterpreted unary function by replacing its %d distinct applications $f(t_1), \\dots, f(t_{%d})$ with fresh variables $v_i$ and adding $t_i = t_j \\rightarrow v_i = v_j$ for every pair. How many such constraints are added?", m, m)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("One constraint per unordered pair: $\\binom{%d}{2} = %d$. The quadratic growth is why solvers prefer congruence closure inside DPLL(T).", m, ans),
			map[string]string{"pairs": edakit.IntStr(ans)})
	}
	for {
		consts := []string{"a", "b", "c", "d"}
		var terms []term
		idx := map[string]int{}
		addT := func(name string, arg int) int {
			if i, ok := idx[name]; ok {
				return i
			}
			terms = append(terms, term{name, arg})
			idx[name] = len(terms) - 1
			return len(terms) - 1
		}
		for _, c := range consts {
			addT(c, -1)
		}
		for _, c := range consts[:3] {
			addT("f("+c+")", idx[c])
		}
		if kp == 2 {
			for _, c := range consts[:2] {
				addT("f(f("+c+"))", idx["f("+c+")"])
			}
		}
		nEq := edakit.Int(2, 3)
		var eqs [][2]int
		var eqText []string
		for len(eqs) < nEq {
			x, y := rand.Intn(len(terms)), rand.Intn(len(terms))
			if x == y {
				continue
			}
			eqs = append(eqs, [2]int{x, y})
			eqText = append(eqText, terms[x].name+" = "+terms[y].name)
		}
		parent := make([]int, len(terms))
		for i := range parent {
			parent[i] = i
		}
		var find func(int) int
		find = func(x int) int {
			for parent[x] != x {
				x = parent[x]
			}
			return x
		}
		union := func(x, y int) bool {
			rx, ry := find(x), find(y)
			if rx == ry {
				return false
			}
			parent[rx] = ry
			return true
		}
		for _, e := range eqs {
			union(e[0], e[1])
		}
		congr := 0
		for changed := true; changed; {
			changed = false
			for i, s := range terms {
				for j, t := range terms {
					if i < j && s.arg >= 0 && t.arg >= 0 && find(s.arg) == find(t.arg) && union(i, j) {
						changed = true
						congr++
					}
				}
			}
		}
		if congr == 0 && rand.Intn(4) > 0 {
			continue
		}
		roots := map[int]bool{}
		var names []string
		for i, t := range terms {
			roots[find(i)] = true
			names = append(names, t.name)
		}
		q := "With constants a, b, c, d and an uninterpreted function f, assume " + strings.Join(eqText, ", ") + ". Into how many equivalence classes does congruence closure split the terms " + strings.Join(names, ", ") + "?"
		return edakit.Numeric(q, edakit.IntStr(len(roots)),
			edakit.Sprintf("Merge the asserted equalities, then repeatedly merge f(s) with f(t) whenever s and t are already in one class (%d congruence merge(s) here). The terms end in %d classes.", congr, len(roots)),
			map[string]string{"classes": edakit.IntStr(len(roots)), "congruence": edakit.IntStr(congr)})
	}
}

// ---------- arrays ----------

func arraysCompute(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp != 2 && (kp == 1 || rand.Intn(2) == 0) {
		k := edakit.Int(2, 4)
		addrs := []int{}
		vals := []int{}
		expr := "a"
		for i := 0; i < k; i++ {
			ad := edakit.Int(0, 7)
			if i > 0 && rand.Intn(3) == 0 {
				ad = addrs[rand.Intn(len(addrs))]
			}
			v := edakit.Int(10, 99)
			addrs = append(addrs, ad)
			vals = append(vals, v)
			expr = edakit.Sprintf("store(%s, %d, %d)", expr, ad, v)
		}
		r := edakit.Int(0, 7)
		if rand.Intn(3) > 0 {
			r = addrs[rand.Intn(k)]
		}
		ans := edakit.Sprintf("the initial content a[%d]", r)
		for i := k - 1; i >= 0; i-- {
			if addrs[i] == r {
				ans = edakit.IntStr(vals[i])
				break
			}
		}
		var ds []string
		for _, v := range vals {
			ds = append(ds, edakit.IntStr(v))
		}
		ds = append(ds, edakit.Sprintf("the initial content a[%d]", r), "0")
		q := edakit.Sprintf("In the theory of arrays, what is select(%s, %d)?", expr, r)
		return edakit.Choice(q, ans,
			"Read-over-write peels stores from the outside in: if the store address equals the read address the stored value is the answer, otherwise the read moves to the inner array. Only when every store misses does the read reach the initial array.",
			ds...)
	}
	type fact struct {
		text string
		ans  string
	}
	facts := []fact{
		{"k = j", "the value w"},
		{"k ≠ j and k = i", "the value v"},
		{"k ≠ j and k ≠ i", "select(a, k)"},
		{"k ≠ j, with nothing known about k and i", "ite(k = i, v, select(a, k))"},
		{"k = i, with nothing known about k and j", "ite(k = j, w, v)"},
	}
	f := facts[rand.Intn(len(facts))]
	opts := []string{"the value w", "the value v", "select(a, k)", "ite(k = i, v, select(a, k))", "ite(k = j, w, v)", "ite(k = j, w, ite(k = i, v, select(a, k)))"}
	var ds []string
	for _, o := range opts {
		if o != f.ans {
			ds = append(ds, o)
		}
	}
	ds = edakit.Shuffle(ds)[:3]
	q := "Simplify select(store(store(a, i, v), j, w), k) given " + f.text + "."
	return edakit.Choice(q, f.ans,
		"Apply read-over-write to the outer store first (compare k with j), then to the inner store (compare k with i). A known equality or disequality resolves a case; an unknown one leaves an ite that the solver must split on.",
		ds...)
}

// ---------- symbolic execution ----------

type branch struct {
	text string
	f    func(x int) bool
}

func randomBranch() branch {
	c := edakit.Int(1, 14)
	switch rand.Intn(4) {
	case 0:
		return branch{edakit.Sprintf("x < %d", c), func(x int) bool { return x < c }}
	case 1:
		return branch{edakit.Sprintf("x >= %d", c), func(x int) bool { return x >= c }}
	case 2:
		return branch{edakit.Sprintf("x == %d", c), func(x int) bool { return x == c }}
	}
	b := edakit.Int(0, 3)
	return branch{edakit.Sprintf("x[%d] == 1", b), func(x int) bool { return x>>b&1 == 1 }}
}

func symbolicExecution(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	k := edakit.Scaled(ctx, 2, 4)
	var bs []branch
	var lines []string
	for i := 0; i < k; i++ {
		b := randomBranch()
		bs = append(bs, b)
		lines = append(lines, edakit.Sprintf("if (%s) y = y + %d;", b.text, i+1))
	}
	prog := "A program reads a 4-bit unsigned input x (0 to 15) and runs, in sequence: " + strings.Join(lines, " ") + " (x[b] is bit b of x.)"
	if kp == 1 {
		if rand.Intn(2) == 0 {
			ans := 1 << k
			return edakit.Numeric(prog+" Ignoring feasibility, how many syntactic paths does the program have?", edakit.IntStr(ans),
				edakit.Sprintf("Each of the %d independent if-statements doubles the paths: $2^{%d} = %d$.", k, k, ans), nil)
		}
		ans := k + 1
		q := edakit.Sprintf("A program on a 4-bit input runs an if / else-if chain with %d conditions and a final else. Ignoring feasibility, how many paths does it have?", k)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Exactly one arm of the chain executes: %d conditions plus the final else give %d paths.", k, ans), nil)
	}
	outcomes := map[int]int{}
	for x := 0; x < 16; x++ {
		v := 0
		for i, b := range bs {
			if b.f(x) {
				v |= 1 << i
			}
		}
		outcomes[v]++
	}
	if kp == 2 {
		ans := len(outcomes)
		return edakit.Numeric(prog+" Symbolic execution forks at each branch and asks the solver whether the path condition is satisfiable. How many feasible paths are there?", edakit.IntStr(ans),
			edakit.Sprintf("Of the $2^{%d}$ syntactic paths, %d have a satisfiable path condition: they are the distinct branch-outcome patterns that some x in 0..15 produces.", k, ans),
			map[string]string{"feasible": edakit.IntStr(ans)})
	}
	var keys []int
	for v := range outcomes {
		keys = append(keys, v)
	}
	v := keys[rand.Intn(len(keys))]
	var conds []string
	for i, b := range bs {
		if v>>i&1 == 1 {
			conds = append(conds, "("+b.text+")")
		} else {
			conds = append(conds, "!("+b.text+")")
		}
	}
	ans := outcomes[v]
	return edakit.Numeric(prog+" How many inputs x follow the path with condition "+strings.Join(conds, " ∧ ")+"?", edakit.IntStr(ans),
		edakit.Sprintf("Intersect the sets of x allowed by each conjunct of the path condition; the count of remaining x values is %d.", ans),
		map[string]string{"models": edakit.IntStr(ans)})
}

// ---------- abstract interpretation with intervals ----------

func absInterp(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		a := edakit.Int(-5, 5)
		b := a + edakit.Int(1, 8)
		c := edakit.Int(-5, 5)
		d := c + edakit.Int(1, 8)
		upper := rand.Intn(2) == 0
		var expr string
		var lo, hi int
		switch rand.Intn(4) {
		case 0:
			expr, lo, hi = "x + y", a+c, b+d
		case 1:
			expr, lo, hi = "x - y", a-d, b-c
		case 2:
			ps := []int{a * c, a * d, b * c, b * d}
			lo, hi = ps[0], ps[0]
			for _, p := range ps {
				lo, hi = minInt(lo, p), maxInt(hi, p)
			}
			expr = "x * y"
		default:
			expr, lo, hi = "x - x", a-b, b-a
		}
		which, ans := "lower", lo
		if upper {
			which, ans = "upper", hi
		}
		q := edakit.Sprintf("Interval analysis knows $x \\in [%d, %d]$ and $y \\in [%d, %d]$. What %s bound does it compute for %s?", a, b, c, d, which, expr)
		why := edakit.Sprintf("Interval arithmetic gives $[%d, %d]$.", lo, hi)
		if expr == "x - x" {
			why += " The domain forgets that both operands are the same variable, so it cannot conclude 0: that is the dependency problem of non-relational domains."
		}
		return edakit.Numeric(q, edakit.IntStr(ans), why, map[string]string{"interval": edakit.Sprintf("[%d, %d]", lo, hi)})
	case 2:
		a := edakit.Int(0, 6)
		b := a + edakit.Int(1, 6)
		k1, k2 := edakit.Int(1, 9), edakit.Int(2, 3)
		lo1, hi1 := a+k1, b+k1
		lo2, hi2 := a*k2, b*k2
		lo, hi := minInt(lo1, lo2), maxInt(hi1, hi2)
		upper := rand.Intn(2) == 0
		which, ans := "lower", lo
		if upper {
			which, ans = "upper", hi
		}
		q := edakit.Sprintf("With $x \\in [%d, %d]$ and an unknown condition c, a program runs: if (c) x = x + %d; else x = x * %d. What %s bound does interval analysis give for x after the merge?", a, b, k1, k2, which)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("The then-branch gives $[%d, %d]$ and the else-branch $[%d, %d]$. The join is the smallest interval containing both: $[%d, %d]$.", lo1, hi1, lo2, hi2, lo, hi),
			map[string]string{"join": edakit.Sprintf("[%d, %d]", lo, hi)})
	}
	N := edakit.Scaled(ctx, 5, 40)
	s := edakit.Int(2, 7)
	// Kleene iteration at the loop head.
	hl, hh := 0, 0
	for it := 0; it < 1000; it++ {
		bl, bh := hl, minInt(hh, N-1)
		nl, nh := minInt(0, bl+s), maxInt(0, bh+s)
		if nl == hl && nh == hh {
			break
		}
		hl, hh = nl, nh
	}
	exitHi := hh
	q := edakit.Sprintf("Interval analysis without widening runs on: i = 0; while (i < %d) i = i + %d. At the fixpoint, what is the upper bound of the interval for i after the loop exits?", N, s)
	return edakit.Numeric(q, edakit.IntStr(exitHi),
		edakit.Sprintf("The loop-head interval grows by %d per iteration until the guard $i \\le %d$ clips it; the clipped body yields at most %d + %d = %d, so the head stabilizes at $[0, %d]$ and the exit interval is $[%d, %d]$. Intervals cannot record that i is a multiple of %d.", s, N-1, N-1, s, N-1+s, hh, N, exitHi, s),
		map[string]string{"head": edakit.Sprintf("[0, %d]", hh)})
}
