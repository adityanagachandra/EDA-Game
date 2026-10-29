package fv

import (
	"math/bits"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// ---------- miters ----------

func miter(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		n := 3
		f := readOnce(n, perm(n))
		var g formula
		if rand.Intn(2) == 0 {
			g = readOnce(n, perm(n))
		} else {
			g = readOnce(n, perm(n)[:2])
		}
		d := bits.OnesCount64((f.table ^ g.table) & fullTable(n))
		q := "Circuit A computes " + f.text + " and circuit B computes " + g.text + " over the inputs $x_1, x_2, x_3$. Their miter XORs the two outputs. On how many of the 8 input vectors is the miter output 1?"
		why := edakit.Sprintf("The miter is 1 exactly where the circuits disagree. Comparing the two truth tables, they differ on %d input vectors.", d)
		if d == 0 {
			why = "The truth tables agree everywhere, so the miter is unsatisfiable and the circuits are equivalent."
		}
		return edakit.Numeric(q, edakit.IntStr(d), why, map[string]string{"distinguishing": edakit.IntStr(d)})
	case 2:
		n := edakit.Scaled(ctx, 2, 16)
		if rand.Intn(2) == 0 {
			ans := 5*n + 2
			q := edakit.Sprintf("Two circuits with %d outputs each are joined into a miter: one XOR per output pair and one %d-input OR of the XORs, whose output is asserted to be 1. Using Tseitin encoding (2-input XOR: 4 clauses; k-input OR: k + 1 clauses; one unit clause for the assertion), how many clauses does the miter glue add on top of the two circuits?", n, n)
			return edakit.Numeric(q, edakit.IntStr(ans),
				edakit.Sprintf("XORs: 4·%d = %d. OR: %d + 1 = %d. Assertion: 1. Total %d.", n, 4*n, n, n+1, ans),
				map[string]string{"clauses": edakit.IntStr(ans)})
		}
		ans := n + 1
		q := edakit.Sprintf("Two circuits with %d outputs each are joined into a miter: one XOR per output pair and one %d-input OR of the XORs. With one Tseitin variable per gate, how many new variables does the miter glue add?", n, n)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("One variable per XOR (%d) plus one for the OR: %d.", n, ans),
			map[string]string{"vars": edakit.IntStr(ans)})
	}
	in := edakit.Scaled(ctx, 4, 32)
	ga, gb := edakit.Int(10, 120), edakit.Int(10, 120)
	o := edakit.Int(2, 8)
	ans := in + ga + gb + o + 1
	q := edakit.Sprintf("Circuit A has %d gates and circuit B has %d gates (inverters excluded); both read the same %d primary inputs and have %d outputs. The miter shares the inputs, XORs each output pair and ORs the XORs. With one Tseitin variable per primary input and per gate, how many variables does the whole miter CNF have?", ga, gb, in, o)
	return edakit.Numeric(q, edakit.IntStr(ans),
		edakit.Sprintf("Inputs are shared, so they count once: %d + %d + %d + (%d XORs + 1 OR) = %d.", in, ga, gb, o, ans),
		map[string]string{"vars": edakit.IntStr(ans)})
}

// ---------- structural hashing ----------

type aigNode struct{ a, b int } // literals: 2*id + complement; id 0 is constant 0

func strashCount(nodes []aigNode, firstNode int, simplify bool) int {
	canon := map[int]int{} // node id -> canonical literal
	table := map[[2]int]int{}
	next := firstNode
	lit := func(l int) int {
		id := l / 2
		if id < firstNode {
			return l
		}
		return canon[id] ^ (l & 1)
	}
	for i, nd := range nodes {
		id := firstNode + i
		p, q := lit(nd.a), lit(nd.b)
		if p > q {
			p, q = q, p
		}
		if simplify {
			switch {
			case p == q:
				canon[id] = p
				continue
			case p^1 == q || p == 0:
				canon[id] = 0
				continue
			case p == 1:
				canon[id] = q
				continue
			}
		}
		key := [2]int{p, q}
		if r, ok := table[key]; ok {
			canon[id] = r
			continue
		}
		table[key] = 2 * next
		canon[id] = 2 * next
		next++
	}
	return len(table)
}

func strash(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	names := []string{"0", "a", "b", "c", "d"}
	first := len(names)
	for {
		var nodes []aigNode
		base := edakit.Scaled(ctx, 3, 5)
		randLit := func(limit int) int {
			id := 1 + rand.Intn(limit-1)
			return 2*id + rand.Intn(2)
		}
		for i := 0; i < base; i++ {
			lim := first + len(nodes)
			a, b := randLit(lim), randLit(lim)
			if a/2 == b/2 {
				continue
			}
			nodes = append(nodes, aigNode{a, b})
		}
		if len(nodes) < 3 {
			continue
		}
		extra := edakit.Int(2, 4)
		for i := 0; i < extra; i++ {
			src := nodes[rand.Intn(len(nodes))]
			switch {
			case kp >= 2 && rand.Intn(2) == 0:
				// cascade: copy a node that reads a node, pointing at a duplicate of the fanin
				ref := -1
				for j, nd := range nodes {
					if nd.a/2 >= first && rand.Intn(2) == 0 {
						ref = j
					}
				}
				if ref < 0 {
					nodes = append(nodes, aigNode{src.b, src.a})
					continue
				}
				fan := nodes[ref].a/2 - first
				nodes = append(nodes, aigNode{nodes[fan].b, nodes[fan].a})
				dupID := first + len(nodes) - 1
				nodes = append(nodes, aigNode{nodes[ref].b, 2*dupID + nodes[ref].a&1})
			case kp == 3 && rand.Intn(2) == 0:
				l := randLit(first + len(nodes))
				if rand.Intn(2) == 0 {
					nodes = append(nodes, aigNode{l, l})
				} else {
					nodes = append(nodes, aigNode{l, l ^ 1})
				}
			default:
				if rand.Intn(2) == 0 {
					nodes = append(nodes, aigNode{src.b, src.a})
				} else {
					nodes = append(nodes, src)
				}
			}
		}
		if len(nodes) > 10 {
			continue
		}
		total := len(nodes)
		after := strashCount(nodes, first, kp == 3)
		if after == total {
			continue
		}
		litName := func(l int) string {
			id := l / 2
			s := ""
			if id < first {
				s = names[id]
			} else {
				s = edakit.Sprintf("n%d", id-first+1)
			}
			if l&1 == 1 {
				return "¬" + s
			}
			return s
		}
		var parts []string
		for i, nd := range nodes {
			parts = append(parts, edakit.Sprintf("n%d = AND(%s, %s)", i+1, litName(nd.a), litName(nd.b)))
		}
		rules := "Structural hashing keys each node by its pair of input literals, in either order, after replacing already-merged inputs by their representatives."
		if kp == 3 {
			rules += " Before hashing it also applies AND(x, x) = x and AND(x, ¬x) = 0."
		}
		q := edakit.Sprintf("An AIG over inputs a, b, c, d has, in topological order: %s. %s How many AND nodes remain after one pass in this order?", strings.Join(parts, "; "), rules)
		return edakit.Numeric(q, edakit.IntStr(after),
			edakit.Sprintf("Walking the nodes in order and looking each normalized input pair up in the hash table, %d of the %d nodes are merged or simplified away, leaving %d.", total-after, total, after),
			map[string]string{"before": edakit.IntStr(total), "after": edakit.IntStr(after)})
	}
}

// ---------- SAT sweeping ----------

func sweeping(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	k := 8
	full := (1 << k) - 1
	for {
		m := edakit.Scaled(ctx, 6, 10)
		nbase := edakit.Int(3, 5)
		var base []int
		for len(base) < nbase {
			s := rand.Intn(full + 1)
			ok := true
			for _, b := range base {
				if b == s || b == full^s {
					ok = false
				}
			}
			if ok {
				base = append(base, s)
			}
		}
		var sigs []int
		for i := 0; i < m; i++ {
			s := base[rand.Intn(nbase)]
			if kp == 3 && rand.Intn(3) == 0 {
				s ^= full
			}
			sigs = append(sigs, s)
		}
		classes := map[int]int{}
		for _, s := range sigs {
			key := s
			if kp == 3 && s&1 == 1 {
				key = s ^ full
			}
			classes[key]++
		}
		big, calls := 0, 0
		for _, c := range classes {
			if c >= 2 {
				big++
				calls += c - 1
			}
		}
		if big == 0 {
			continue
		}
		var parts []string
		for i, s := range sigs {
			parts = append(parts, edakit.Sprintf("n%d: %s", i+1, bitsN(s, k)))
		}
		head := "After 8 rounds of random simulation the AIG nodes have signatures " + strings.Join(parts, ", ") + "."
		switch kp {
		case 1:
			q := head + " Nodes with identical signatures are candidate equivalences. How many candidate classes contain at least two nodes?"
			return edakit.Numeric(q, edakit.IntStr(big),
				edakit.Sprintf("Grouping identical signatures gives %d classes with two or more members; singletons are already proved distinct by simulation.", big),
				map[string]string{"classes": edakit.IntStr(big)})
		case 2:
			q := head + " The sweeper picks one representative per class of identical signatures and proves every other member equal to it with one SAT call. How many SAT calls does it make if every candidate holds?"
			return edakit.Numeric(q, edakit.IntStr(calls),
				edakit.Sprintf("A class of size $s$ needs $s - 1$ calls. Summing over the classes gives %d.", calls),
				map[string]string{"calls": edakit.IntStr(calls)})
		}
		q := head + " The sweeper also treats a node and the complement of another as candidates, so a signature and its bitwise complement fall into one class. How many candidate classes contain at least two nodes?"
		return edakit.Numeric(q, edakit.IntStr(big),
			edakit.Sprintf("Normalize each signature up to complement (for example, so that its last bit is 0) and group: %d classes have two or more members.", big),
			map[string]string{"classes": edakit.IntStr(big)})
	}
}

// ---------- sequential equivalence ----------

type moore struct {
	out  []int
	next [][2]int
}

func randomMoore(n int) moore {
	m := moore{out: make([]int, n), next: make([][2]int, n)}
	for s := 0; s < n; s++ {
		m.out[s] = rand.Intn(2)
		m.next[s] = [2]int{rand.Intn(n), rand.Intn(n)}
	}
	return m
}

func (m moore) text(prefix string) string {
	var parts []string
	for s := range m.out {
		parts = append(parts, edakit.Sprintf("%s%d outputs %d, goes to %s%d on input 0 and %s%d on input 1", prefix, s, m.out[s], prefix, m.next[s][0], prefix, m.next[s][1]))
	}
	return strings.Join(parts, "; ")
}

func productBFS(a, b moore) (int, int) {
	type pair struct{ x, y int }
	dist := map[pair]int{{0, 0}: 0}
	q := []pair{{0, 0}}
	first := -1
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		if a.out[p.x] != b.out[p.y] && first < 0 {
			first = dist[p]
		}
		for in := 0; in < 2; in++ {
			n := pair{a.next[p.x][in], b.next[p.y][in]}
			if _, ok := dist[n]; !ok {
				dist[n] = dist[p] + 1
				q = append(q, n)
			}
		}
	}
	return len(dist), first
}

func sequential(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		na, nb := 3, edakit.Int(2, 3)
		if kp == 3 {
			na, nb = 4, edakit.Int(3, 4)
		}
		a, b := randomMoore(na), randomMoore(nb)
		if a.out[0] != b.out[0] {
			continue
		}
		reach, first := productBFS(a, b)
		head := "Machine A: " + a.text("A") + ". Machine B: " + b.text("B") + ". Both start in state 0 and read the same input bit each cycle."
		if kp == 1 {
			if reach < 3 {
				continue
			}
			q := head + " How many states of the product machine (pairs of A and B states) are reachable, including the start pair?"
			return edakit.Numeric(q, edakit.IntStr(reach),
				edakit.Sprintf("Breadth-first search over pairs from (A0, B0), applying both inputs, reaches %d pairs.", reach),
				map[string]string{"reachable": edakit.IntStr(reach)})
		}
		if first < 1 || (kp == 3 && first < 2) {
			continue
		}
		q := head + " Their outputs agree in the start states. What is the length of the shortest input sequence after which the outputs differ?"
		return edakit.Numeric(q, edakit.IntStr(first),
			edakit.Sprintf("Search the product machine breadth-first for a pair whose outputs differ; the nearest such pair is at depth %d from (A0, B0). That sequence is the counterexample a sequential equivalence checker reports.", first),
			map[string]string{"depth": edakit.IntStr(first), "reachable": edakit.IntStr(reach)})
	}
}
