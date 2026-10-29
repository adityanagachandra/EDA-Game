package boolalg

import (
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// truthTable: ON-set size (KP 1), counting functions (KP 2), distance between two functions (KP 3).
func truthTable(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := nVars(ctx)
		f, s := drawSOP(n)
		q := edakit.Sprintf("Let $f(%s) = %s$. On how many of the %d rows of its truth table is $f = 1$?", varList(n), s.latex(), rows(n))
		ans := edakit.IntStr(f.ones())
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Evaluating the expression row by row gives the ON-set m(%s), which has %s rows.", mintermList(f.minterms()), ans),
			map[string]string{"on_set": mintermList(f.minterms())})
	case 2:
		switch rand.Intn(3) {
		case 0:
			n := edakit.Scaled(ctx, 2, 5)
			ans := edakit.IntStr(pow2(rows(n)))
			return edakit.Numeric(edakit.Sprintf("How many distinct Boolean functions of %d input variables are there?", n), ans,
				edakit.Sprintf("A truth table of %d variables has %d rows, and each row's output bit is free, so there are 2^%d = %s functions.", n, rows(n), rows(n), ans), nil)
		case 1:
			n := edakit.Int(3, 4)
			r := edakit.Int(1, rows(n)-2)
			ans := edakit.IntStr(pow2(rows(n) - r))
			return edakit.Numeric(edakit.Sprintf("A specification fixes the output of a %d-input Boolean function on %d of its rows and leaves the other rows free. How many functions meet the specification?", n, r), ans,
				edakit.Sprintf("Of the %d rows, %d are fixed and %d are free, so there are 2^%d = %s functions.", rows(n), r, rows(n)-r, rows(n)-r, ans), nil)
		default:
			n := edakit.Int(3, 4)
			k := edakit.Int(1, 4)
			if n == 3 {
				k = edakit.Int(1, 7)
			}
			ans := edakit.IntStr(binom(rows(n), k))
			return edakit.Numeric(edakit.Sprintf("How many Boolean functions of %d variables have exactly %d minterms in their ON-set?", n, k), ans,
				edakit.Sprintf("Choose which %d of the %d rows are 1: C(%d, %d) = %s.", k, rows(n), rows(n), k, ans), nil)
		}
	default:
		n := nVars(ctx)
		f, fs := drawFunc(n, 0.5)
		g := f
		for g == f {
			g = tt{n, f.w}
			for i := 0; i < edakit.Int(1, 4); i++ {
				g.w ^= 1 << uint(rand.Intn(rows(n)))
			}
		}
		gs := sigma("g", g)
		d := f.xor(g)
		ans := edakit.IntStr(d.ones())
		q := "Two Boolean functions are given: " + fs + " and " + gs + ". On how many input rows do $f$ and $g$ differ?"
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The rows where exactly one of them is 1 are m(%s): %s rows. This count is the number of counterexamples an equivalence check would find.", mintermList(d.minterms()), ans),
			map[string]string{"diff": mintermList(d.minterms())})
	}
}

// pickDependentVar returns a variable f depends on.
func pickDependentVar(f tt) int {
	for {
		v := rand.Intn(f.n)
		if f.dependsOn(v) {
			return v
		}
	}
}

func cofName(v, val int) string { return "f_{" + lit(v, val) + "}" }

func cofactor(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	kp := kpOf(ctx)
	if kp == 3 {
		n = 4
		f, fs := drawFunc(n, 0.5)
		vs := rand.Perm(n)[:2]
		if vs[0] > vs[1] {
			vs[0], vs[1] = vs[1], vs[0]
		}
		a, b := rand.Intn(2), rand.Intn(2)
		g := f.cof(vs[1], b).cof(vs[0], a)
		ans := edakit.IntStr(g.ones())
		q := edakit.Sprintf("Let %s. The cofactor of $f$ with respect to the cube $%s%s$ fixes both literals. How many of the %d rows of $f_{%s%s}$, over %s, are 1?",
			fs, lit(vs[0], a), lit(vs[1], b), rows(n-2), lit(vs[0], a), lit(vs[1], b), otherVars(n, vs...))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Keep the rows of f with %s = %d and %s = %d. Of those %d rows, %s have f = 1.", varNames[vs[0]], a, varNames[vs[1]], b, rows(n-2), ans),
			map[string]string{"cofactor_ones": ans})
	}
	var f tt
	var fs string
	if kp == 1 {
		f, fs = drawFunc(n, 1)
	} else {
		f, fs = drawFunc(n, 0)
	}
	v := pickDependentVar(f)
	val := rand.Intn(2)
	g := f.cof(v, val)
	ans := edakit.IntStr(g.ones())
	q := edakit.Sprintf("Let %s. Compute the cofactor $%s$ (set $%s = %d$), a function of %s. How many of its %d rows are 1?",
		fs, cofName(v, val), varNames[v], val, otherVars(n, v), rows(n-1))
	return edakit.Numeric(q, ans,
		edakit.Sprintf("Setting %s = %d keeps the %d rows with that value. Of those, %s have f = 1.", varNames[v], val, rows(n-1), ans),
		map[string]string{"cofactor_ones": ans})
}

func shannon(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := nVars(ctx)
		v := rand.Intn(n)
		g, gs := drawSOP(n - 1)
		h, hs := drawSOP(n - 1)
		ren := func(t sop) string { return liftSOP(t, v).latex() }
		ans := edakit.IntStr(g.ones() + h.ones())
		x := varNames[v]
		q := edakit.Sprintf("A function of %s is built by Shannon expansion on $%s$: $f = %s\\,(%s) + \\overline{%s}\\,(%s)$. On how many of the %d rows of $f$ is $f = 1$?",
			varList(n), x, x, ren(gs), x, ren(hs), rows(n))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The two terms are disjoint, since one needs %s = 1 and the other %s = 0. The positive cofactor has %d ON rows and the negative cofactor has %d, so f has %s.", x, x, g.ones(), h.ones(), ans),
			map[string]string{"pos": edakit.IntStr(g.ones()), "neg": edakit.IntStr(h.ones())})
	case 2:
		n := nVars(ctx)
		v := rand.Intn(n)
		var f tt
		for {
			f = tt{n, rand.Uint32() & fullMask(n)}
			if f.dependsOn(v) && f.ones() >= 2 {
				break
			}
		}
		p, q0 := f.cof(v, 1), f.cof(v, 0)
		rest := otherVars(n, v)
		x := varNames[v]
		q := edakit.Sprintf("Number minterms over %s in that order. The cofactors of $f(%s)$ are $f_{%s} = %s$ and $f_{\\overline{%s}} = %s$, each over %s. Rebuild $f$ by Shannon expansion and list its minterm numbers, separated by commas.",
			varList(n), varList(n), x, sumOrZero(p), x, sumOrZero(q0), rest)
		ans := mintermList(f.minterms())
		return edakit.Numeric(q, strings.ReplaceAll(ans, " ", ""),
			edakit.Sprintf("Insert %s = 1 into each minterm of the positive cofactor and %s = 0 into each minterm of the negative cofactor, at %s's bit position. The result is m(%s).", x, x, x, ans),
			map[string]string{"minterms": ans})
	default:
		n := 4
		f, fs := drawFunc(n, 0.5)
		k := edakit.Int(1, 3)
		seen := map[uint32]bool{}
		for assign := 0; assign < pow2(k); assign++ {
			g := f
			// Fix variables 0..k-1, taking the cofactor of the highest index first.
			for v := k - 1; v >= 0; v-- {
				g = g.cof(v, (assign>>uint(k-1-v))&1)
			}
			seen[g.w] = true
		}
		ans := edakit.IntStr(len(seen))
		q := edakit.Sprintf("Let %s. Expand $f$ by Shannon expansion on %s, in that order, down to depth %d. Of the %d subfunctions at that depth (functions of the remaining variables), how many are distinct?",
			fs, strings.Join(varNames[:k], ", then "), k, pow2(k))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Fixing the first %d variables in all %d ways gives %s distinct cofactors. A decision diagram stores each distinct one once.", k, pow2(k), ans),
			map[string]string{"distinct": ans})
	}
}

func booleanDifference(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	f, fs := drawFunc(n, 0.6)
	v := pickDependentVar(f)
	d := f.cof(v, 0).xor(f.cof(v, 1))
	x := varNames[v]
	switch kpOf(ctx) {
	case 1:
		ans := edakit.IntStr(d.ones())
		q := edakit.Sprintf("Let %s. The Boolean difference $\\partial f / \\partial %s = f_{%s} \\oplus f_{\\overline{%s}}$ is a function of %s. On how many of its %d rows is it 1?",
			fs, x, x, x, otherVars(n, v), rows(n-1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("XOR the two cofactors row by row: they differ on %s of the %d rows.", ans, rows(n-1)),
			map[string]string{"diff_ones": ans})
	case 2:
		ans := edakit.IntStr(2 * d.ones())
		q := edakit.Sprintf("Let %s. For how many of the %d input vectors does flipping $%s$ alone change the value of $f$?", fs, rows(n), x)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The Boolean difference with respect to %s is 1 on %d assignments of the other variables. Each pairs with both values of %s, so %s vectors are sensitive to %s.", x, d.ones(), x, ans, x),
			map[string]string{"diff_ones": edakit.IntStr(d.ones())})
	default:
		stuck := rand.Intn(2)
		// A test-mode restriction fixes one other variable.
		w := v
		for w == v {
			w = rand.Intn(n)
		}
		wv := rand.Intn(2)
		count := 0
		for i := 0; i < rows(n); i++ {
			if bitOf(i, v, n) != 1-stuck || bitOf(i, w, n) != wv {
				continue
			}
			// Index of i in the cofactor domain.
			if f.at(i) != f.at(i^(1<<uint(n-1-v))) {
				count++
			}
		}
		ans := edakit.IntStr(count)
		q := edakit.Sprintf("Let %s, implemented as a single block. Input $%s$ may be stuck at %d. Test mode forces $%s = %d$. How many of the %d input vectors with $%s = %d$ detect the fault at the output of $f$?",
			fs, x, stuck, varNames[w], wv, rows(n-1), varNames[w], wv)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("A test must set %s = %d (the opposite of the stuck value) and make the Boolean difference with respect to %s equal 1 so the error reaches the output. With %s = %d as well, %s vectors qualify.", x, 1-stuck, x, varNames[w], wv, ans),
			map[string]string{"tests": ans})
	}
}

const (
	clsPos = iota
	clsNeg
	clsBinate
	clsIndep
)

func unateClass(f tt, v int) int {
	n0, n1 := f.cof(v, 0), f.cof(v, 1)
	switch {
	case n0 == n1:
		return clsIndep
	case n0.w&^n1.w == 0:
		return clsPos
	case n1.w&^n0.w == 0:
		return clsNeg
	}
	return clsBinate
}

func unateness(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	switch kpOf(ctx) {
	case 1:
		want := rand.Intn(3)
		for {
			f, s := drawSOP(n)
			v := rand.Intn(n)
			if unateClass(f, v) != want {
				continue
			}
			x := varNames[v]
			labels := []string{"Positive unate in $" + x + "$", "Negative unate in $" + x + "$", "Binate in $" + x + "$", "Independent of $" + x + "$"}
			why := []string{
				"The negative cofactor implies the positive one, so raising " + x + " never lowers f: positive unate.",
				"The positive cofactor implies the negative one, so raising " + x + " never raises f: negative unate.",
				"Neither cofactor implies the other: some rows rise with " + x + " and others fall, so f is binate in " + x + ".",
			}[want]
			var ds []string
			for i, l := range labels {
				if i != want {
					ds = append(ds, l)
				}
			}
			q := edakit.Sprintf("Let $f(%s) = %s$. How does $f$ depend on $%s$? Compare $f_{%s}$ with $f_{\\overline{%s}}$.", varList(n), s.latex(), x, x, x)
			return edakit.Choice(q, labels[want], why, ds...)
		}
	case 2:
		f, fs := drawFunc(n, 0.5)
		count := 0
		var parts []string
		names := []string{"positive unate", "negative unate", "binate", "independent"}
		for v := 0; v < n; v++ {
			c := unateClass(f, v)
			if c == clsBinate {
				count++
			}
			parts = append(parts, varNames[v]+" is "+names[c])
		}
		q := "Let " + fs + ". How many of its variables is $f$ binate in?"
		var ds []string
		for k := 0; k <= n; k++ {
			if k != count {
				ds = append(ds, edakit.IntStr(k))
			}
		}
		return edakit.Choice(q, edakit.IntStr(count), "Comparing the two cofactors of each variable: "+strings.Join(parts, "; ")+".", ds...)
	default:
		for {
			f, s := drawSOP(n)
			bin := -1
			nb := 0
			for v := 0; v < n; v++ {
				if unateClass(f, v) == clsBinate {
					bin = v
					nb++
				}
			}
			if nb != 1 {
				continue
			}
			// Prefer instances where some other variable appears in both polarities but is unate.
			opts := []string{}
			for v := 0; v < n; v++ {
				if f.dependsOn(v) || s.mentions(v) {
					opts = append(opts, "$"+varNames[v]+"$")
				}
			}
			opts = append(opts, "No variable: $f$ is unate")
			ans := "$" + varNames[bin] + "$"
			var ds []string
			for _, o := range opts {
				if o != ans {
					ds = append(ds, o)
				}
			}
			q := edakit.Sprintf("Let $f(%s) = %s$. Exactly one variable is binate in the function (not merely in the expression). Which one?", varList(n), s.latex())
			return edakit.Choice(q, ans, "Only for "+varNames[bin]+" does neither cofactor imply the other. A variable that appears in one polarity throughout a cover is unate, and the function decides, not the way it is written.", ds...)
		}
	}
}

func symmetricPair(f tt, i, j int) bool {
	return f.cof(j, 1).cof(i, 0) == f.cof(j, 0).cof(i, 1)
}

func symmetry(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		var ws []int
		for len(ws) == 0 || len(ws) == n+1 {
			ws = nil
			for k := 0; k <= n; k++ {
				if rand.Intn(2) == 0 {
					ws = append(ws, k)
				}
			}
		}
		total := 0
		var terms []string
		for _, k := range ws {
			total += binom(n, k)
			terms = append(terms, edakit.Sprintf("C(%d,%d)", n, k))
		}
		ans := edakit.IntStr(total)
		q := edakit.Sprintf("A totally symmetric function of %d inputs is 1 exactly when the number of inputs at 1 is one of {%s}. How many of its %d rows are 1?", n, mintermList(ws), rows(n))
		return edakit.Numeric(q, ans, "Count the rows of each allowed weight: "+strings.Join(terms, " + ")+" = "+ans+".", nil)
	case 2:
		n := 4
		var f tt
		f, _ = drawFunc(n, 0)
		// Symmetrize over a random pair with probability one half to get non-trivial counts.
		if rand.Intn(2) == 0 {
			i, j := 0, 0
			for i == j {
				i, j = rand.Intn(n), rand.Intn(n)
			}
			var w uint32
			for r := 0; r < rows(n); r++ {
				bi, bj := bitOf(r, i, n), bitOf(r, j, n)
				src := r
				if bi == 0 && bj == 1 {
					src = r ^ (1 << uint(n-1-i)) ^ (1 << uint(n-1-j))
				}
				if f.at(src) {
					w |= 1 << uint(r)
				}
			}
			f = tt{n, w}
		}
		if f.isConst() {
			f = tt{n, 0x0116}
		}
		count := 0
		var pairs []string
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if symmetricPair(f, i, j) {
					count++
					pairs = append(pairs, "("+varNames[i]+","+varNames[j]+")")
				}
			}
		}
		ans := edakit.IntStr(count)
		list := "none"
		if len(pairs) > 0 {
			list = strings.Join(pairs, ", ")
		}
		q := "Let " + sigma("f", f) + ". In how many of the 6 unordered pairs of variables is $f$ symmetric, meaning that swapping the two variables leaves $f$ unchanged?"
		return edakit.Numeric(q, ans, "A pair (x, y) is symmetric when the cofactor with x = 1, y = 0 equals the cofactor with x = 0, y = 1. The symmetric pairs are: "+list+".", map[string]string{"pairs": list})
	default:
		n := edakit.Scaled(ctx, 2, 8)
		switch rand.Intn(3) {
		case 0:
			ans := edakit.IntStr(pow2(n + 1))
			return edakit.Numeric(edakit.Sprintf("How many totally symmetric Boolean functions of %d variables are there?", n), ans,
				edakit.Sprintf("A totally symmetric function is fixed by one output bit per weight 0 to %d, which is %d bits, so there are 2^%d = %s.", n, n+1, n+1, ans), nil)
		case 1:
			ans := edakit.IntStr(pow2(n))
			return edakit.Numeric(edakit.Sprintf("How many totally symmetric Boolean functions of %d variables output 0 on the all-zeros input?", n), ans,
				edakit.Sprintf("The bit for weight 0 is fixed and the other %d weights are free: 2^%d = %s.", n, n, ans), nil)
		default:
			ans := edakit.IntStr(pow2(n - 1))
			return edakit.Numeric(edakit.Sprintf("How many totally symmetric Boolean functions of %d variables output 0 on the all-zeros input and 1 on the all-ones input?", n), ans,
				edakit.Sprintf("Two of the %d weight bits are fixed, leaving %d free: 2^%d = %s.", n+1, n-1, n-1, ans), nil)
		}
	}
}

func quantification(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	kp := kpOf(ctx)
	if kp == 3 {
		n = 4
	}
	f, fs := drawFunc(n, 0.5)
	v := pickDependentVar(f)
	x := varNames[v]
	switch kp {
	case 1:
		g := f.cof(v, 0).or(f.cof(v, 1))
		ans := edakit.IntStr(g.ones())
		q := edakit.Sprintf("Let %s. The property \"there exists a value of $%s$ that makes $f$ true\" is $\\exists %s\\, f$, a function of %s. On how many of its %d rows is it 1?",
			fs, x, x, otherVars(n, v), rows(n-1))
		return edakit.Numeric(q, ans, edakit.Sprintf("Existential quantification is the OR of the cofactors: f with %s = 0 OR f with %s = 1. That OR is 1 on %s rows.", x, x, ans), nil)
	case 2:
		g := f.cof(v, 0).and(f.cof(v, 1))
		ans := edakit.IntStr(g.ones())
		q := edakit.Sprintf("Let %s. The property \"for every value of $%s$, $f$ is true\" is $\\forall %s\\, f$, a function of %s. On how many of its %d rows is it 1?",
			fs, x, x, otherVars(n, v), rows(n-1))
		return edakit.Numeric(q, ans, edakit.Sprintf("Universal quantification is the AND of the cofactors. f with %s = 0 AND f with %s = 1 is 1 on %s rows.", x, x, ans), nil)
	default:
		u := v
		for u == v {
			u = rand.Intn(n)
		}
		existsOuter := rand.Intn(2) == 0
		// Inner quantifier on v, outer on u. Cofactor indices shift after removing v.
		var inner tt
		if existsOuter {
			inner = f.cof(v, 0).and(f.cof(v, 1)) // forall v
		} else {
			inner = f.cof(v, 0).or(f.cof(v, 1)) // exists v
		}
		uu := u
		if u > v {
			uu = u - 1
		}
		var outer tt
		var phrase string
		if existsOuter {
			outer = inner.cof(uu, 0).or(inner.cof(uu, 1))
			phrase = edakit.Sprintf("$\\exists %s\\, \\forall %s\\, f$", varNames[u], x)
		} else {
			outer = inner.cof(uu, 0).and(inner.cof(uu, 1))
			phrase = edakit.Sprintf("$\\forall %s\\, \\exists %s\\, f$", varNames[u], x)
		}
		ans := edakit.IntStr(outer.ones())
		q := edakit.Sprintf("Let %s. A property quantifies two of its inputs: evaluate %s, a function of %s. On how many of its %d rows is it 1?",
			fs, phrase, otherVars(n, u, v), rows(n-2))
		return edakit.Numeric(q, ans, edakit.Sprintf("Quantify the inner variable %s first (AND of cofactors for for-all, OR for exists), then the outer variable %s. The result is 1 on %s rows.", x, varNames[u], ans), nil)
	}
}

// liftSOP inserts an absent variable at position v, so a cover over n-1 variables names the
// variables other than v.
func liftSOP(s sop, v int) sop {
	out := make(sop, len(s))
	for i, c := range s {
		d := make(cube, 0, len(c)+1)
		d = append(d, c[:v]...)
		d = append(d, 2)
		d = append(d, c[v:]...)
		out[i] = d
	}
	return out
}

// sumOrZero renders a cofactor as a minterm sum, or 0 when it is empty.
func sumOrZero(f tt) string {
	if f.ones() == 0 {
		return "0"
	}
	return `\sum m(` + mintermList(f.minterms()) + ")"
}
