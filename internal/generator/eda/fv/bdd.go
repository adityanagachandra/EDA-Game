package fv

import (
	"math/bits"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// A Boolean function of n <= 5 variables is a truth table: bit a of the table is f at the
// assignment a, where variable i (shown as x_{i+1}) is bit i of a.

func fullTable(n int) uint64 { return (uint64(1) << (uint(1) << uint(n))) - 1 }

// cof returns the cofactor of g with variable v fixed to b, still as an n-variable table.
func cof(g uint64, n, v, b int) uint64 {
	var out uint64
	for a := 0; a < 1<<n; a++ {
		src := a &^ (1 << v)
		if b == 1 {
			src |= 1 << v
		}
		if g>>uint(src)&1 == 1 {
			out |= 1 << uint(a)
		}
	}
	return out
}

func varTable(n, v int) uint64 {
	var out uint64
	for a := 0; a < 1<<n; a++ {
		if a>>v&1 == 1 {
			out |= 1 << uint(a)
		}
	}
	return out
}

type bnode struct{ v, lo, hi int }

// bdd builds ROBDDs over one variable order; nodes 0 and 1 are the terminals.
type bdd struct {
	n     int
	order []int
	uniq  map[[3]int]int
	memo  map[uint64]int
	nodes []bnode // index i holds node i+2
}

func newBDD(n int, order []int) *bdd {
	return &bdd{n: n, order: order, uniq: map[[3]int]int{}, memo: map[uint64]int{}}
}

func (b *bdd) build(g uint64) int {
	full := fullTable(b.n)
	if g == 0 {
		return 0
	}
	if g == full {
		return 1
	}
	if id, ok := b.memo[g]; ok {
		return id
	}
	for _, v := range b.order {
		lo, hi := cof(g, b.n, v, 0), cof(g, b.n, v, 1)
		if lo == hi {
			continue
		}
		l, h := b.build(lo), b.build(hi)
		key := [3]int{v, l, h}
		id, ok := b.uniq[key]
		if !ok {
			b.nodes = append(b.nodes, bnode{v, l, h})
			id = len(b.nodes) + 1
			b.uniq[key] = id
		}
		b.memo[g] = id
		return id
	}
	return 0
}

func (b *bdd) size() int { return len(b.nodes) }

func robddSize(g uint64, n int, order []int) int {
	b := newBDD(n, order)
	b.build(g)
	return b.size()
}

// quasiSize counts the nodes of the quasi-reduced BDD: distinct subfunctions per level.
func quasiSize(g uint64, n int, order []int) int {
	level := map[uint64]bool{g: true}
	total := 0
	for l := 0; l < n; l++ {
		total += len(level)
		next := map[uint64]bool{}
		for h := range level {
			next[cof(h, n, order[l], 0)] = true
			next[cof(h, n, order[l], 1)] = true
		}
		level = next
	}
	return total
}

func orderText(order []int, names []string) string {
	var s []string
	for _, v := range order {
		s = append(s, names[v])
	}
	return strings.Join(s, " < ")
}

func xNames(n int) []string {
	var s []string
	for i := 0; i < n; i++ {
		s = append(s, edakit.Sprintf("x%d", i+1))
	}
	return s
}

// onSetText lists the assignments where g is 1, as bit strings x1x2...xn.
func onSetText(g uint64, n int) string {
	var s []string
	for a := 0; a < 1<<n; a++ {
		if g>>uint(a)&1 == 1 {
			var bs strings.Builder
			for i := 0; i < n; i++ {
				bs.WriteByte(byte('0' + a>>i&1))
			}
			s = append(s, bs.String())
		}
	}
	return strings.Join(s, ", ")
}

func header(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(edakit.Sprintf("x%d", i+1))
	}
	return b.String()
}

// randomOnSet draws a nonconstant function of n variables that depends on every variable.
func randomOnSet(n int) uint64 {
	for {
		var g uint64
		for a := 0; a < 1<<n; a++ {
			if rand.Intn(2) == 0 {
				g |= 1 << uint(a)
			}
		}
		k := bits.OnesCount64(g)
		if k < 2 || k > (1<<n)-2 || (n == 4 && (k < 3 || k > 7)) {
			continue
		}
		ok := true
		for v := 0; v < n; v++ {
			if cof(g, n, v, 0) == cof(g, n, v, 1) {
				ok = false
			}
		}
		if ok {
			return g
		}
	}
}

// formula is a read-once Boolean formula over a subset of the variables.
type formula struct {
	text  string
	table uint64
}

// readOnce builds a random read-once formula over the given variables (each used once).
func readOnce(n int, vars []int) formula {
	leaf := func(v int) formula {
		t := varTable(n, v)
		name := edakit.Sprintf("x%d", v+1)
		if rand.Intn(3) == 0 {
			return formula{"¬" + name, fullTable(n) &^ t}
		}
		return formula{name, t}
	}
	combine := func(a, b formula) formula {
		switch rand.Intn(3) {
		case 0:
			return formula{a.text + " ∧ " + b.text, a.table & b.table}
		case 1:
			return formula{a.text + " ∨ " + b.text, a.table | b.table}
		}
		return formula{a.text + " ⊕ " + b.text, a.table ^ b.table}
	}
	vs := edakit.Shuffle(vars)
	cur := leaf(vs[0])
	for i := 1; i < len(vs); i++ {
		nxt := leaf(vs[i])
		if strings.Contains(cur.text, " ") {
			cur.text = "(" + cur.text + ")"
		}
		cur = combine(cur, nxt)
	}
	return cur
}

func perm(n int) []int { return rand.Perm(n) }

// ---------- reduction rules ----------

func reductionRules(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 1 {
		n := edakit.Scaled(ctx, 2, 9)
		if rand.Intn(3) == 0 {
			q := edakit.Sprintf("A complete binary decision tree tests %d variables, one per level, on every path. How many terminal leaves does it have?", n)
			return edakit.Numeric(q, edakit.IntStr(1<<n), edakit.Sprintf("Every path fixes all %d variables, so there is one leaf per assignment: $2^{%d} = %d$.", n, n, 1<<n), nil)
		}
		q := edakit.Sprintf("Before any reduction, the decision tree of a function of %d variables tests every variable on every path. How many nonterminal (decision) nodes does it have?", n)
		ans := (1 << n) - 1
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Level $i$ has $2^i$ nodes for $i = 0, \\dots, %d$, a total of $2^{%d} - 1 = %d$.", n-1, n, ans),
			map[string]string{"nodes": edakit.IntStr(ans)})
	}
	n := 3
	if ctx.Difficulty > 0.6 {
		n = edakit.Int(3, 4)
	}
	var g uint64
	var desc string
	if rand.Intn(2) == 0 {
		f := readOnce(n, perm(n))
		g, desc = f.table, "$f = $ "+f.text
	} else {
		g = randomOnSet(n)
		desc = "the function $f$ that is 1 exactly on the inputs " + header(n) + " = " + onSetText(g, n)
	}
	order := perm(n)
	qr, r := quasiSize(g, n, order), robddSize(g, n, order)
	ot := orderText(order, xNames(n))
	if kp == 2 {
		q := "Take " + desc + ", with the variable order " + ot + ". Build the BDD that merges isomorphic subgraphs but keeps every level on every path (the quasi-reduced BDD). How many nonterminal nodes does it have?"
		return edakit.Numeric(q, edakit.IntStr(qr),
			edakit.Sprintf("Level $i$ needs one node per distinct cofactor of $f$ by the first $i$ variables of the order. Summing the distinct cofactors per level gives %d nodes.", qr),
			map[string]string{"quasi": edakit.IntStr(qr), "robdd": edakit.IntStr(r)})
	}
	q := "Take " + desc + ", with the variable order " + ot + ". Starting from the quasi-reduced BDD (isomorphic subgraphs merged, every level kept), how many nonterminal nodes does the redundant-test rule (delete a node whose two children are equal) remove on the way to the ROBDD?"
	return edakit.Numeric(q, edakit.IntStr(qr-r),
		edakit.Sprintf("The quasi-reduced BDD has %d nonterminal nodes and the ROBDD has %d, so %d redundant tests are removed.", qr, r, qr-r),
		map[string]string{"quasi": edakit.IntStr(qr), "robdd": edakit.IntStr(r)})
}

// ---------- ROBDD node count ----------

func robdd(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		n := 3
		if ctx.Difficulty > 0.6 && rand.Intn(2) == 0 {
			n = 4
		}
		f := readOnce(n, perm(n))
		order := perm(n)
		r := robddSize(f.table, n, order)
		q := "How many nonterminal nodes does the ROBDD of $f = $ " + f.text + " have under the variable order " + orderText(order, xNames(n)) + "?"
		return edakit.Numeric(q, edakit.IntStr(r),
			edakit.Sprintf("Cofactor $f$ variable by variable in the order given, keep one node per distinct cofactor that still depends on the variable being tested, and share equal cofactors. That leaves %d nonterminal nodes.", r),
			map[string]string{"nodes": edakit.IntStr(r)})
	case 2:
		n := edakit.Scaled(ctx, 3, 4)
		g := randomOnSet(n)
		order := perm(n)
		r := robddSize(g, n, order)
		q := "The function $f$ is 1 exactly on the inputs " + header(n) + " = " + onSetText(g, n) + " and 0 elsewhere. How many nonterminal nodes does its ROBDD have under the order " + orderText(order, xNames(n)) + "?"
		return edakit.Numeric(q, edakit.IntStr(r),
			edakit.Sprintf("Split the truth table by the variables in order; distinct non-constant cofactors that depend on the tested variable each need one node. The count is %d.", r),
			map[string]string{"nodes": edakit.IntStr(r)})
	}
	n := 3
	f := readOnce(n, perm(n))
	vs := perm(n)[:2]
	g := readOnce(n, vs)
	if rand.Intn(2) == 0 {
		g = readOnce(n, perm(n))
	}
	order := perm(n)
	b := newBDD(n, order)
	b.build(f.table)
	sf := b.size()
	b.build(g.table)
	shared := b.size()
	sg := robddSize(g.table, n, order)
	q := "A BDD package stores $f = $ " + f.text + " and $g = $ " + g.text + " in one shared ROBDD (one unique table) under the order " + orderText(order, xNames(n)) + ". How many nonterminal nodes does the shared graph hold?"
	return edakit.Numeric(q, edakit.IntStr(shared),
		edakit.Sprintf("Alone, $f$ needs %d nodes and $g$ needs %d. Subgraphs that represent the same subfunction are stored once, so the shared graph holds %d.", sf, sg, shared),
		map[string]string{"f": edakit.IntStr(sf), "g": edakit.IntStr(sg), "shared": edakit.IntStr(shared)})
}

// ---------- variable order ----------

func variableOrder(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 1 || kp == 2 {
		m := edakit.Scaled(ctx, 2, 5)
		var terms []string
		for i := 1; i <= m; i++ {
			terms = append(terms, edakit.Sprintf("x_%d y_%d", i, i))
		}
		f := "$f = " + strings.Join(terms, " + ") + "$"
		var inter, sep []string
		for i := 1; i <= m; i++ {
			inter = append(inter, edakit.Sprintf("x%d < y%d", i, i))
		}
		for i := 1; i <= m; i++ {
			sep = append(sep, edakit.Sprintf("x%d", i))
		}
		for i := 1; i <= m; i++ {
			sep = append(sep, edakit.Sprintf("y%d", i))
		}
		if kp == 1 {
			ans := 2 * m
			q := "How many nonterminal nodes does the ROBDD of " + f + " have under the interleaved order " + strings.Join(inter, " < ") + "?"
			return edakit.Numeric(q, edakit.IntStr(ans),
				edakit.Sprintf("Each pair needs one node for $x_i$ and one for $y_i$; on failure the path falls through to the next pair, on success it goes to 1. That is %d nodes.", ans),
				map[string]string{"nodes": edakit.IntStr(ans)})
		}
		ans := (1 << (m + 1)) - 2
		q := "How many nonterminal nodes does the ROBDD of " + f + " have under the separated order " + strings.Join(sep, " < ") + "?"
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("After reading all %d $x$ variables the graph must remember which subset of them is 1. The $x$ levels form a full tree of $2^{%d} - 1$ nodes and the $y$ levels add $2^{%d} - 1$ more: $2^{%d} - 2 = %d$.", m, m, m, m+1, ans),
			map[string]string{"nodes": edakit.IntStr(ans)})
	}
	n := 4
	for {
		g := randomOnSet(n)
		o1, o2 := perm(n), perm(n)
		a, b := robddSize(g, n, o1), robddSize(g, n, o2)
		if a == b && rand.Intn(3) > 0 {
			continue
		}
		q := "The function $f$ is 1 exactly on the inputs " + header(n) + " = " + onSetText(g, n) + ". Let $A$ be its ROBDD node count (nonterminal nodes) under " + orderText(o1, xNames(n)) + " and $B$ under " + orderText(o2, xNames(n)) + ". What is $A - B$?"
		return edakit.Numeric(q, edakit.IntStr(a-b),
			edakit.Sprintf("Building both: $A = %d$ and $B = %d$, so $A - B = %d$.", a, b, a-b),
			map[string]string{"A": edakit.IntStr(a), "B": edakit.IntStr(b)})
	}
}

// ---------- ite and apply ----------

func iteApply(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 3 {
		a, b := edakit.Scaled(ctx, 4, 40), edakit.Int(4, 30)
		q := edakit.Sprintf("BDD $F$ has %d nodes and BDD $G$ has %d nodes, both counting the two terminals. With a computed table, at most how many distinct recursive calls can apply($F$, $G$, op) make?", a, b)
		return edakit.Numeric(q, edakit.IntStr(a*b),
			edakit.Sprintf("Every call is indexed by a pair of nodes, one from each graph, and the computed table answers repeats in constant time: at most %d·%d = %d calls.", a, b, a*b),
			map[string]string{"bound": edakit.IntStr(a * b)})
	}
	n := 3
	if kp == 2 && ctx.Difficulty > 0.6 {
		n = 4
	}
	for {
		f := readOnce(n, perm(n)[:edakit.Int(2, n)])
		g := readOnce(n, perm(n)[:edakit.Int(2, n)])
		var op string
		var t uint64
		if kp == 1 {
			if rand.Intn(2) == 0 {
				op, t = "f \\wedge g", f.table&g.table
			} else {
				op, t = "f \\vee g", f.table|g.table
			}
		} else {
			op, t = "f \\oplus g", f.table^g.table
		}
		if t == 0 || t == fullTable(n) {
			continue
		}
		order := []int{0, 1, 2, 3}[:n]
		r := robddSize(t, n, order)
		q := "Let $f = $ " + f.text + " and $g = $ " + g.text + ". Apply computes the ROBDD of $" + op + "$ under the order " + orderText(order, xNames(n)) + ". How many nonterminal nodes does the result have?"
		return edakit.Numeric(q, edakit.IntStr(r),
			edakit.Sprintf("Apply splits both arguments on the top variable, recurses on the cofactors and rebuilds through the unique table, so the result is the canonical ROBDD of $%s$: %d nonterminal nodes.", op, r),
			map[string]string{"nodes": edakit.IntStr(r)})
	}
}

// ---------- sat count ----------

func satCount(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	n := 3
	if kp >= 2 {
		n = edakit.Scaled(ctx, 4, 5)
	}
	for {
		g := randomOnSet(n)
		if n == 5 {
			// thin out to keep the BDD readable
			g &= randomOnSet(n) | randomOnSet(n)
			if g == 0 {
				continue
			}
		}
		order := perm(n)
		b := newBDD(n, order)
		root := b.build(g)
		if root < 2 || b.size() > 9 {
			continue
		}
		table := bddTable(b, root, order)
		cnt := bits.OnesCount64(g)
		head := "An ROBDD over $x_1, \\dots, x_" + edakit.IntStr(n) + "$ with order " + orderText(order, xNames(n)) + " has root N1 and nodes: " + table + ". "
		if kp == 3 {
			p := float64(cnt) / float64(int(1)<<n)
			q := head + "What is the probability that a uniformly random input satisfies the function? Round to 4 decimal places."
			return edakit.Numeric(q, edakit.Num(p, 4),
				edakit.Sprintf("Counting paths to 1, weighting each by $2^k$ for the $k$ levels it skips, gives %d satisfying inputs of %d: probability %s.", cnt, 1<<n, edakit.Num(p, 4)),
				map[string]string{"count": edakit.IntStr(cnt)})
		}
		q := head + edakit.Sprintf("How many of the %d assignments to the %d variables satisfy the function?", 1<<n, n)
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("Process nodes bottom-up: a node's count is the low child's count times $2^{\\text{skipped levels}}$ plus the same for the high child, and the root's count is scaled by the levels above it. The total is %d.", cnt),
			map[string]string{"count": edakit.IntStr(cnt)})
	}
}

// bddTable names nodes N1.. in breadth-first order from the root and lists them.
func bddTable(b *bdd, root int, order []int) string {
	name := map[int]string{0: "0", 1: "1"}
	q := []int{root}
	var seq []int
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if _, ok := name[x]; ok {
			continue
		}
		name[x] = edakit.Sprintf("N%d", len(seq)+1)
		seq = append(seq, x)
		nd := b.nodes[x-2]
		q = append(q, nd.lo, nd.hi)
	}
	var parts []string
	for _, x := range seq {
		nd := b.nodes[x-2]
		parts = append(parts, edakit.Sprintf("%s tests x%d (low → %s, high → %s)", name[x], nd.v+1, name[nd.lo], name[nd.hi]))
	}
	return strings.Join(parts, "; ")
}
