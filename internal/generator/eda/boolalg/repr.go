package boolalg

import (
	"math/bits"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func hexWord(w uint32, n int) string {
	if n == 3 {
		return edakit.Sprintf("0x%02X", w&0xFF)
	}
	return edakit.Sprintf("0x%04X", w&0xFFFF)
}

// reverseTable reads the table in the opposite bit order, the classic convention slip.
func reverseTable(w uint32, n int) uint32 {
	r := bits.Reverse32(w)
	return r >> uint(32-rows(n))
}

func convention(n int) string {
	if n == 3 {
		return "Bit $i$ of the word holds $f$ at minterm $i$, where $i$ is written in binary as $abc$ with $a$ the most significant bit; bit 0 is the least significant bit of the word."
	}
	return "Bit $i$ of the word holds $f$ at minterm $i$, where $i$ is written in binary as $abcd$ with $a$ the most significant bit; bit 0 is the least significant bit of the word."
}

func bitParallel(ctx generator.GeneratorContext) generator.Problem {
	n := nVars(ctx)
	mask := fullMask(n)
	opts := func(correct uint32, alts ...uint32) (string, []string) {
		c := hexWord(correct, n)
		seen := map[string]bool{c: true}
		var ds []string
		for _, a := range alts {
			h := hexWord(a&mask, n)
			if !seen[h] {
				seen[h] = true
				ds = append(ds, h)
			}
		}
		for len(ds) < 3 {
			h := hexWord(correct^(1<<uint(rand.Intn(rows(n)))), n)
			if !seen[h] {
				seen[h] = true
				ds = append(ds, h)
			}
		}
		return c, ds
	}
	switch kpOf(ctx) {
	case 1:
		f, s := drawSOP(n)
		c, ds := opts(f.w, reverseTable(f.w, n), ^f.w, f.w^(1<<uint(f.minterms()[0])))
		q := edakit.Sprintf("A tool stores $f(%s) = %s$ as a %d-bit truth-table word. %s Which word is it?", varList(n), s.latex(), rows(n), convention(n))
		return edakit.Choice(q, c, "The ON-set is m("+mintermList(f.minterms())+"); set those bits and read the word in hex from bit "+edakit.IntStr(rows(n)-1)+" down to bit 0: "+c+".", ds...)
	case 2:
		// Combine projection words with word-level operations.
		u, v := 0, 0
		for u == v {
			u, v = rand.Intn(n), rand.Intn(n)
		}
		w := u
		for w == u || w == v {
			w = rand.Intn(n)
		}
		pu, pv, pw := varTT(n, u).w, varTT(n, v).w, varTT(n, w).w
		var res uint32
		var expr string
		var wrong uint32
		// Reversed-convention projection words, for the convention distractor.
		ru, rv, rw := varTT(n, n-1-u).w, varTT(n, n-1-v).w, varTT(n, n-1-w).w
		switch rand.Intn(3) {
		case 0:
			res = (pu & ^pv) | pw
			wrong = (ru & ^rv) | rw
			expr = varNames[u] + `\overline{` + varNames[v] + `} + ` + varNames[w]
		case 1:
			res = (pu ^ pv) & pw
			wrong = (ru ^ rv) & rw
			expr = `(` + varNames[u] + ` \oplus ` + varNames[v] + `)\,` + varNames[w]
		default:
			res = ^(pu & pv) ^ pw
			wrong = ^(ru & rv) ^ rw
			expr = `\overline{` + varNames[u] + varNames[v] + `} \oplus ` + varNames[w]
		}
		res &= mask
		c, ds := opts(res, wrong, ^res, res^pw)
		proj := make([]string, n)
		for i := 0; i < n; i++ {
			proj[i] = varNames[i] + " = " + hexWord(varTT(n, i).w, n)
		}
		q := edakit.Sprintf("%s The projection words are %s. Using only word-wide AND, OR, XOR and NOT, what is the word of $%s$?", convention(n), strings.Join(proj, ", "), expr)
		return edakit.Choice(q, c, "Apply the operations to the projection words bit-parallel, then mask to "+edakit.IntStr(rows(n))+" bits: "+c+". One instruction evaluates the expression on every row at once.", ds...)
	default:
		n = 4
		mask = fullMask(n)
		f, _ := drawFunc(n, 0)
		v := rand.Intn(n)
		val := rand.Intn(2)
		stride := uint(1) << uint(n-1-v)
		m := varTT(n, v).w
		var res, wrongDir, noCopy uint32
		if val == 1 {
			hi := f.w & m
			res = hi | hi>>stride
			wrongDir = hi | hi<<stride
			noCopy = hi
		} else {
			lo := f.w &^ m
			res = lo | lo<<stride
			wrongDir = lo | lo>>stride
			noCopy = lo
		}
		res &= mask
		other := f.w & m
		if val == 1 {
			other = f.w &^ m
			other = other | other<<stride
		} else {
			other = other | other>>stride
		}
		c, ds := opts(res, wrongDir, noCopy, other)
		q := edakit.Sprintf("%s A function has the word $t$ = %s. The projection word of $%s$ is %s. Compute the cofactor $f_{%s}$ as a 16-bit word that no longer depends on $%s$, by masking the half where $%s = %d$ and copying it over the other half with a shift of %d.",
			convention(n), hexWord(f.w, n), varNames[v], hexWord(m, n), lit(v, val), varNames[v], varNames[v], val, stride)
		why := edakit.Sprintf("Keep the bits where %s = %d, then shift them by %d toward the positions where %s = %d and OR: %s.", varNames[v], val, stride, varNames[v], 1-val, c)
		return edakit.Choice(q, c, why, ds...)
	}
}

// lit2 is a signal literal in an AIG: node id (0 = constant 0, 1..4 = inputs, then gates)
// and a complement flag.
type lit2 struct {
	id  int
	neg bool
}

func hashConsing(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	ng := edakit.Scaled(ctx, 5, 8)
	inputs := 4
	names := map[int]string{}
	for i := 1; i <= inputs; i++ {
		names[i] = varNames[i-1]
	}
	type gate struct{ a, b lit2 }
	var gates []gate
	render := func(l lit2) string {
		s := names[l.id]
		if l.neg {
			return `\overline{` + s + `}`
		}
		return s
	}
	randSig := func(upto int) lit2 {
		id := edakit.Int(1, upto)
		return lit2{id, kp >= 2 && rand.Intn(3) == 0}
	}
	for len(gates) < ng {
		id := inputs + len(gates) + 1
		names[id] = "g_{" + edakit.IntStr(len(gates)+1) + "}"
		r := rand.Float64()
		switch {
		case len(gates) > 0 && r < 0.3:
			// Duplicate an earlier gate, swapped at KP 2 and above.
			g := gates[rand.Intn(len(gates))]
			if kp >= 2 && rand.Intn(2) == 0 {
				g.a, g.b = g.b, g.a
			}
			gates = append(gates, g)
		case kp == 3 && r < 0.45:
			s := randSig(id - 1)
			if rand.Intn(2) == 0 {
				gates = append(gates, gate{s, s})
			} else {
				gates = append(gates, gate{s, lit2{s.id, !s.neg}})
			}
		default:
			a, b := randSig(id-1), randSig(id-1)
			for a.id == b.id {
				b = randSig(id - 1)
			}
			if kp == 1 && b.id < a.id {
				// At KP 1 operands are always listed in order, so only exact copies repeat.
				a, b = b, a
			}
			gates = append(gates, gate{a, b})
		}
	}
	// Structural hashing with canonical operand order and the trivial rules.
	canon := map[int]lit2{}
	for i := 1; i <= inputs; i++ {
		canon[i] = lit2{i, false}
	}
	resolve := func(l lit2) lit2 {
		c := canon[l.id]
		if l.neg {
			c.neg = !c.neg
		}
		return c
	}
	type key struct{ a, b lit2 }
	table := map[key]int{}
	next := inputs + 1
	less := func(x, y lit2) bool {
		if x.id != y.id {
			return x.id < y.id
		}
		return !x.neg && y.neg
	}
	for gi, g := range gates {
		id := inputs + gi + 1
		a, b := resolve(g.a), resolve(g.b)
		if less(b, a) {
			a, b = b, a
		}
		switch {
		case a.id == 0 && !a.neg:
			canon[id] = lit2{0, false}
			continue
		case a.id == 0 && a.neg:
			canon[id] = b
			continue
		case a == b:
			canon[id] = a
			continue
		case a.id == b.id:
			canon[id] = lit2{0, false}
			continue
		}
		k := key{a, b}
		if n, ok := table[k]; ok {
			canon[id] = lit2{n, false}
			continue
		}
		table[k] = next
		canon[id] = lit2{next, false}
		next++
	}
	var lines []string
	for gi, g := range gates {
		lines = append(lines, edakit.Sprintf("$g_{%d} = %s \\land %s$", gi+1, render(g.a), render(g.b)))
	}
	ans := edakit.IntStr(len(table))
	q := "An and-inverter graph over inputs $a, b, c, d$ is built gate by gate: " + strings.Join(lines, ", ") + "."
	switch kp {
	case 1:
		q += " The builder hash-conses: before creating an AND node it looks up the pair of operands in a unique table. How many AND nodes does the table hold at the end?"
	case 2:
		q += " The builder hash-conses with a canonical key (operands sorted, so AND is commutative) and rewrites every later reference to a duplicate through the node it was merged into. How many AND nodes does the unique table hold at the end?"
	default:
		q += " Before hashing, the builder simplifies $x \\land x = x$ and $x \\land \\overline{x} = 0$, propagates constants, sorts operands, and rewrites references to merged nodes. Dead nodes are kept. How many AND nodes does the unique table hold at the end?"
	}
	return edakit.Numeric(q, ans, edakit.Sprintf("Of the %d gate definitions, %s produce new keys; the rest are duplicates or simplify away. Hash-consing gives each distinct structure exactly one node, so equality of nodes is a pointer comparison.", len(gates), ans), map[string]string{"nodes": ans})
}

// robdd counts the internal nodes of the reduced ordered BDD of f under the order 0..n-1.
func robdd(f tt) int {
	total := 0
	for i := 0; i < f.n; i++ {
		seen := map[uint32]bool{}
		for assign := 0; assign < pow2(i); assign++ {
			g := f
			for v := i - 1; v >= 0; v-- {
				g = g.cof(v, (assign>>uint(i-1-v))&1)
			}
			// g is a function of variables i..n-1, renumbered from 0.
			if g.dependsOn(0) {
				seen[g.w] = true
			}
		}
		total += len(seen)
	}
	return total
}

func decisionDiagram(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	n := 3
	if kp >= 2 {
		n = 4
	}
	var f tt
	var fs string
	for {
		f, fs = drawFunc(n, 0.6)
		if f.supportSize() == n {
			break
		}
	}
	if kp <= 2 {
		k := robdd(f)
		ans := edakit.IntStr(k)
		q := edakit.Sprintf("Let %s. Build its reduced ordered BDD with the order %s (first variable at the top). How many internal (non-terminal) nodes does it have?", fs, strings.Join(varNames[:n], " < "))
		return edakit.Numeric(q, ans, "At each level, count the distinct subfunctions reached by fixing the variables above it that still depend on that level's variable; sum over levels: "+ans+" nodes, plus the two terminals.", map[string]string{"nodes": ans})
	}
	perm := rand.Perm(n) // perm[v] = position of variable v in the new order
	for tries := 0; tries < 20 && robdd(f.permute(perm)) == robdd(f); tries++ {
		perm = rand.Perm(n)
	}
	order := make([]string, n)
	for v := 0; v < n; v++ {
		order[perm[v]] = varNames[v]
	}
	g := f.permute(perm)
	k := robdd(g)
	k0 := robdd(f)
	ans := edakit.IntStr(k)
	q := edakit.Sprintf("Let %s. Build its reduced ordered BDD with the variable order %s (first variable at the top). How many internal nodes does it have?", fs, strings.Join(order, " < "))
	return edakit.Numeric(q, ans, edakit.Sprintf("Recount level by level in the new order: %s nodes. Under the order %s it would have %d.", ans, strings.Join(varNames[:n], " < "), k0), map[string]string{"nodes": ans})
}
