package synth

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// aig is a small and-inverter graph. Nodes 0..nPI-1 are primary inputs, the rest are AND
// nodes with two fanins (node indices; complements are irrelevant to the structural drills).
type aig struct {
	nPI  int
	fin  [][2]int // fin[i] for i >= nPI
	outs []int
}

func (g *aig) name(i int) string {
	if i < g.nPI {
		return varNames[i]
	}
	return edakit.Sprintf("n%d", i-g.nPI+1)
}

func (g *aig) size() int { return len(g.fin) - g.nPI }

func (g *aig) desc() string {
	var parts []string
	for i := g.nPI; i < len(g.fin); i++ {
		parts = append(parts, edakit.Sprintf("%s = AND(%s, %s)", g.name(i), g.name(g.fin[i][0]), g.name(g.fin[i][1])))
	}
	var outs []string
	for _, o := range g.outs {
		outs = append(outs, g.name(o))
	}
	return strings.Join(parts, "; ") + ". Primary outputs: " + strings.Join(outs, ", ") + "."
}

// randAIG draws nAnd AND nodes over nPI inputs, each fanin pair distinct, every node used.
func randAIG(nPI, nAnd int) *aig {
	for {
		g := &aig{nPI: nPI, fin: make([][2]int, nPI)}
		used := map[int]bool{}
		seen := map[[2]int]bool{}
		for i := nPI; i < nPI+nAnd; i++ {
			for {
				lo := 0
				if i-nPI > 2 {
					lo = i - 4
				}
				a, b := edakit.Int(lo, i-1), edakit.Int(0, i-1)
				if a == b {
					continue
				}
				if a > b {
					a, b = b, a
				}
				if seen[[2]int{a, b}] {
					continue
				}
				seen[[2]int{a, b}] = true
				g.fin = append(g.fin, [2]int{a, b})
				used[a], used[b] = true, true
				break
			}
		}
		last := len(g.fin) - 1
		g.outs = []int{last}
		for i := nPI; i < last; i++ {
			if !used[i] {
				g.outs = append(g.outs, i)
			}
		}
		if len(g.outs) <= 2 {
			return g
		}
	}
}

func (g *aig) refs() []int {
	r := make([]int, len(g.fin))
	for i := g.nPI; i < len(g.fin); i++ {
		r[g.fin[i][0]]++
		r[g.fin[i][1]]++
	}
	for _, o := range g.outs {
		r[o]++
	}
	return r
}

// mffc returns the number of AND nodes in the maximum fanout-free cone of n.
func (g *aig) mffc(n int) int {
	r := g.refs()
	var deref func(v int) int
	deref = func(v int) int {
		if v < g.nPI {
			return 0
		}
		c := 1
		for _, f := range g.fin[v] {
			r[f]--
			if r[f] == 0 {
				c += deref(f)
			}
		}
		return c
	}
	return deref(n)
}

func (g *aig) cone(n int) int {
	seen := map[int]bool{}
	var walk func(v int)
	walk = func(v int) {
		if v < g.nPI || seen[v] {
			return
		}
		seen[v] = true
		walk(g.fin[v][0])
		walk(g.fin[v][1])
	}
	walk(n)
	return len(seen)
}

func (g *aig) level() []int {
	l := make([]int, len(g.fin))
	for i := g.nPI; i < len(g.fin); i++ {
		a, b := l[g.fin[i][0]], l[g.fin[i][1]]
		if b > a {
			a = b
		}
		l[i] = a + 1
	}
	return l
}

// cuts enumerates the k-feasible cuts of every node (as bitmasks over node indices).
func (g *aig) cuts(k int) [][]uint64 {
	cs := make([][]uint64, len(g.fin))
	for i := 0; i < len(g.fin); i++ {
		triv := uint64(1) << i
		if i < g.nPI {
			cs[i] = []uint64{triv}
			continue
		}
		set := map[uint64]bool{triv: true}
		for _, u := range cs[g.fin[i][0]] {
			for _, v := range cs[g.fin[i][1]] {
				w := u | v
				if popcount(w) <= k {
					set[w] = true
				}
			}
		}
		for w := range set {
			cs[i] = append(cs[i], w)
		}
		sort.Slice(cs[i], func(a, b int) bool { return cs[i][a] < cs[i][b] })
	}
	return cs
}

func (g *aig) lutDepth(k int) []int {
	cs := g.cuts(k)
	lab := make([]int, len(g.fin))
	for i := g.nPI; i < len(g.fin); i++ {
		best := 1 << 30
		for _, c := range cs[i] {
			if c == 1<<i {
				continue
			}
			m := 0
			for j := 0; j < i; j++ {
				if c>>j&1 == 1 && lab[j] > m {
					m = lab[j]
				}
			}
			if m+1 < best {
				best = m + 1
			}
		}
		lab[i] = best
	}
	return lab
}

// aigGate is a gate whose AIG cost and depth are fixed.
type aigGate struct {
	name        string
	ins, ands   int
	depth       int
	twoLevelSet bool // allowed in the depth drill
}

var aigGates = []aigGate{
	{"AND", 2, 1, 1, true}, {"OR", 2, 1, 1, true}, {"NAND", 2, 1, 1, true}, {"NOR", 2, 1, 1, true},
	{"XOR", 2, 3, 2, true}, {"XNOR", 2, 3, 2, true}, {"MUX", 3, 3, 2, true},
	{"AND3", 3, 2, 2, false}, {"OR3", 3, 2, 2, false}, {"AND4", 4, 3, 2, false}, {"XOR3", 3, 6, 4, false},
}

// aigStructure: AND count of a gate netlist, its AIG depth, AND count of wide gates.
func aigStructure(ctx generator.GeneratorContext) generator.Problem {
	which := kp(ctx)
	nG := edakit.Scaled(ctx, 3, 5)
	var pool []aigGate
	for _, g := range aigGates {
		if which == 2 && !g.twoLevelSet {
			continue
		}
		if which == 1 && g.ins > 3 {
			continue
		}
		pool = append(pool, g)
	}
	if which == 3 {
		pool = aigGates[7:]
		pool = append(pool, aigGates[4], aigGates[6])
	}
	signals := []string{"a", "b", "c", "d"}
	lvl := map[string]int{"a": 0, "b": 0, "c": 0, "d": 0}
	total := 0
	var lines []string
	last := ""
	for i := 0; i < nG; i++ {
		gt := edakit.Pick(pool)
		var ins []string
		if last != "" {
			ins = append(ins, last)
		}
		for len(ins) < gt.ins {
			s := edakit.Pick(signals)
			dup := false
			for _, x := range ins {
				if x == s {
					dup = true
				}
			}
			if !dup {
				ins = append(ins, s)
			}
		}
		ins = edakit.Shuffle(ins)
		name := edakit.Sprintf("g%d", i+1)
		m := 0
		for _, x := range ins {
			if lvl[x] > m {
				m = lvl[x]
			}
		}
		lvl[name] = m + gt.depth
		total += gt.ands
		lines = append(lines, edakit.Sprintf("%s = %s(%s)", name, gt.name, strings.Join(ins, ", ")))
		signals = append(signals, name)
		last = name
	}
	net := strings.Join(lines, "; ")
	switch which {
	case 2:
		d := lvl[last]
		return edakit.Numeric("Convert the netlist "+net+" into an AIG, building each XOR, XNOR and MUX from three AND nodes in two levels and each two-input gate from one AND node with complemented edges. Primary inputs are at level 0. What is the level of "+last+"?",
			edakit.IntStr(d), edakit.Sprintf("Inverters are free edge attributes, so AND, OR, NAND and NOR add one level and XOR, XNOR and MUX add two. %s sits at level %d.", last, d),
			map[string]string{"depth": edakit.IntStr(d)})
	}
	return edakit.Numeric("Convert the netlist "+net+" into an and-inverter graph without sharing any logic between gates. Each k-input AND or OR costs k - 1 AND nodes, a two-input XOR, XNOR or a MUX costs 3, and a three-input XOR costs 6. How many AND nodes does the AIG have?",
		edakit.IntStr(total), edakit.Sprintf("OR and NAND are ANDs with complemented edges, so only the AND count matters. Summing gate by gate gives %d AND nodes.", total),
		map[string]string{"ands": edakit.IntStr(total)})
}

// strash: AND nodes after structural hashing, optionally with trivial-node rules.
func strash(ctx generator.GeneratorContext) generator.Problem {
	which := kp(ctx)
	nPI := 4
	n := edakit.Scaled(ctx, 6, 9)
	// Literals: 2*node + complement; node 0 is constant 0, nodes 1..nPI are inputs.
	type def struct{ a, b int }
	var defs []def
	name := func(lit int) string {
		v := lit / 2
		s := ""
		if v == 0 {
			s = "0"
		} else if v <= nPI {
			s = varNames[v-1]
		} else {
			s = edakit.Sprintf("n%d", v-nPI)
		}
		if lit%2 == 1 {
			if v == 0 {
				return "1"
			}
			return "!" + s
		}
		return s
	}
	for i := 0; i < n; i++ {
		top := 2 * (nPI + 1 + i)
		r := edakit.Int(0, 9)
		switch {
		case r < 3 && i > 0:
			d := edakit.Pick(defs)
			defs = append(defs, def{d.b, d.a})
		case r == 3 && which == 3:
			x := edakit.Int(2, top-1)
			if edakit.Int(0, 1) == 0 {
				defs = append(defs, def{x, x})
			} else {
				defs = append(defs, def{x, x ^ 1})
			}
		default:
			a, b := edakit.Int(2, top-1), edakit.Int(2, top-1)
			for a/2 == b/2 {
				b = edakit.Int(2, top-1)
			}
			defs = append(defs, def{a, b})
		}
	}
	// Hash.
	mp := make([]int, nPI+1+n) // node -> literal after hashing
	for v := 0; v <= nPI; v++ {
		mp[v] = 2 * v
	}
	tr := func(lit int) int { return mp[lit/2] ^ (lit & 1) }
	tab := map[[2]int]int{}
	created := 0
	for i, d := range defs {
		a, b := tr(d.a), tr(d.b)
		if a > b {
			a, b = b, a
		}
		var res int
		switch {
		case a == 0:
			res = 0
		case a == 1:
			res = b
		case a == b:
			res = a
		case a == b^1:
			res = 0
		default:
			if v, ok := tab[[2]int{a, b}]; ok {
				res = v
			} else {
				created++
				res = 2 * (nPI + 1 + i)
				tab[[2]int{a, b}] = res
			}
		}
		mp[nPI+1+i] = res
	}
	var lines []string
	for i, d := range defs {
		lines = append(lines, edakit.Sprintf("n%d = AND(%s, %s)", i+1, name(d.a), name(d.b)))
	}
	rule := " Structural hashing keys each AND on its unordered pair of fanin literals (after earlier merges)."
	if which == 3 {
		rule += " It also applies $x \\cdot x = x$, $x \\cdot \\bar{x} = 0$ and constant propagation."
	}
	return edakit.Numeric("An AIG builder receives, in order, "+strings.Join(lines, "; ")+" (! marks a complemented edge)."+rule+" How many distinct AND nodes remain?",
		edakit.IntStr(created), edakit.Sprintf("Walking the list in order and looking each normalized fanin pair up in the hash table, %d nodes are new; the rest map to existing nodes%s.", created, map[bool]string{true: " or vanish by the trivial rules", false: ""}[which == 3]),
		map[string]string{"nodes": edakit.IntStr(created)})
}

// aigCuts: number of k-feasible cuts, and of non-dominated cuts.
func aigCuts(ctx generator.GeneratorContext) generator.Problem {
	which := kp(ctx)
	k := 3
	if which == 2 {
		k = 4
	}
	g := randAIG(k+1, edakit.Scaled(ctx, 4, 6))
	cs := g.cuts(k)
	n := len(g.fin) - 1
	if which == 1 {
		n = edakit.Int(g.nPI+1, len(g.fin)-1)
	}
	set := cs[n]
	if which == 3 {
		cnt := 0
		for _, c := range set {
			dom := false
			for _, d := range set {
				if d != c && d&^c == 0 {
					dom = true
					break
				}
			}
			if !dom {
				cnt++
			}
		}
		return edakit.Numeric("In the AIG "+g.desc()+edakit.Sprintf(" Enumerate the 3-feasible cuts of %s, then drop every cut that strictly contains another cut of %s. How many cuts remain (count the trivial cut {%s})?", g.name(n), g.name(n), g.name(n)),
			edakit.IntStr(cnt), edakit.Sprintf("%s has %d 3-feasible cuts; %d are not supersets of another cut. A dominated cut never gives a better mapping.", g.name(n), len(set), cnt),
			map[string]string{"cuts": edakit.IntStr(cnt)})
	}
	return edakit.Numeric("In the AIG "+g.desc()+edakit.Sprintf(" A k-feasible cut of a node is a set of at most k nodes such that every path from a primary input to the node passes through the set. How many %d-feasible cuts does %s have, including the trivial cut {%s}?", k, g.name(n), g.name(n)),
		edakit.IntStr(len(set)), edakit.Sprintf("Merge the cut sets of the two fanins pairwise, keep unions of size at most %d, and add the trivial cut: %d cuts.", k, len(set)),
		map[string]string{"cuts": edakit.IntStr(len(set))})
}

// mffcDrill: MFFC size, and shared logic in the cone.
func mffcDrill(ctx generator.GeneratorContext) generator.Problem {
	which := kp(ctx)
	g := randAIG(4, edakit.Scaled(ctx, 5, 8))
	n := len(g.fin) - 1
	if which == 1 {
		n = edakit.Int(g.nPI+2, len(g.fin)-1)
	}
	m := g.mffc(n)
	if which == 3 {
		cone := g.cone(n)
		return edakit.Numeric("In the AIG "+g.desc()+edakit.Sprintf(" How many AND nodes in the transitive fanin cone of %s (including %s) lie outside its maximum fanout-free cone, that is, are shared with logic outside it?", g.name(n), g.name(n)),
			edakit.IntStr(cone-m), edakit.Sprintf("The cone has %d AND nodes and the MFFC %d, so %d are shared and survive if %s is replaced.", cone, m, cone-m, g.name(n)),
			map[string]string{"shared": edakit.IntStr(cone - m)})
	}
	return edakit.Numeric("In the AIG "+g.desc()+edakit.Sprintf(" What is the size of the maximum fanout-free cone of %s, counting AND nodes only (including %s)?", g.name(n), g.name(n)),
		edakit.IntStr(m), edakit.Sprintf("Dereference %s: decrement the reference count of each fanin, and recurse into any AND fanin whose count reaches 0. %d nodes are freed, the MFFC.", g.name(n), m),
		map[string]string{"mffc": edakit.IntStr(m)})
}

// rewriting: gain of a replacement.
func rewriting(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		m := edakit.Int(2, 7)
		s := edakit.Int(1, 6)
		r := edakit.Int(0, s-1)
		gain := m - (s - r)
		return edakit.Numeric(edakit.Sprintf("Rewriting node n with a 4-input cut: its MFFC has %d AND nodes. The best precomputed structure for the cut function has %d AND nodes, of which %d already exist in the graph (found by structural hashing). What is the gain in AND nodes? (Negative means the graph grows.)", m, s, r),
			edakit.IntStr(gain), edakit.Sprintf("Replacing n frees the %d MFFC nodes and adds only the %d - %d = %d new nodes. Gain %d.", m, s, r, s-r, gain),
			map[string]string{"gain": edakit.IntStr(gain)})
	case 2:
		m := edakit.Int(3, 7)
		best := -1 << 30
		var parts []string
		for i := 0; i < 3; i++ {
			s := edakit.Int(1, 6)
			r := edakit.Int(0, s-1)
			if m-(s-r) > best {
				best = m - (s - r)
			}
			parts = append(parts, edakit.Sprintf("candidate %d has %d nodes, %d already present", i+1, s, r))
		}
		return edakit.Numeric(edakit.Sprintf("Node n has an MFFC of %d AND nodes. Three precomputed structures implement its cut function: %s. What is the best gain available?", m, strings.Join(parts, "; ")),
			edakit.IntStr(best), edakit.Sprintf("Gain is MFFC size minus newly added nodes for each candidate; the best is %d.", best),
			map[string]string{"gain": edakit.IntStr(best)})
	}
	g := randAIG(4, edakit.Scaled(ctx, 5, 8))
	n := len(g.fin) - 1
	m := g.mffc(n)
	s := edakit.Int(1, 4)
	r := edakit.Int(0, s-1)
	gain := m - (s - r)
	return edakit.Numeric("In the AIG "+g.desc()+edakit.Sprintf(" A rewrite of %s uses a structure of %d AND nodes, %d of which already exist outside the MFFC of %s. What is the gain in AND nodes? (Negative means the graph grows.)", g.name(n), s, r, g.name(n)),
		edakit.IntStr(gain), edakit.Sprintf("The MFFC of %s has %d nodes, all freed by the replacement; %d new nodes are added. Gain %d.", g.name(n), m, s-r, gain),
		map[string]string{"mffc": edakit.IntStr(m), "gain": edakit.IntStr(gain)})
}

// refactoring: AIG size of a factored form, savings over SOP, gain against the MFFC.
func refactoring(ctx generator.GeneratorContext) generator.Problem {
	f := randReadOnce(edakit.Scaled(ctx, 4, 7), 2, 10)
	s := f.sop()
	l := f.leaves()
	stem := "$F = " + f.str() + "$"
	switch kp(ctx) {
	case 1:
		return edakit.Numeric("Refactoring rebuilds a cone from the factored form "+stem+". How many AND nodes does the AIG of this factored form need, with inversions free?",
			edakit.IntStr(l-1), edakit.Sprintf("Every binary AND or OR of a factored form is one AND node (OR by De Morgan), and a form with %d literals has %d binary operators.", l, l-1),
			map[string]string{"ands": edakit.IntStr(l - 1)})
	case 2:
		return edakit.Numeric("The cone of a node computes $"+s.str()+"$, built as a two-level SOP in the AIG. Refactoring replaces it with the factored form "+stem+". How many AND nodes are saved?",
			edakit.IntStr(s.lits()-l), edakit.Sprintf("An SOP with %d literals costs %d AND nodes; the factored form with %d literals costs %d. Saving %d.", s.lits(), s.lits()-1, l, l-1, s.lits()-l),
			map[string]string{"saving": edakit.IntStr(s.lits() - l)})
	}
	m := edakit.Int(l-2, l+4)
	if m < 1 {
		m = 1
	}
	gain := m - (l - 1)
	return edakit.Numeric(edakit.Sprintf("A node's large cut is collapsed and factored as %s. The node's MFFC has %d AND nodes and none of the new structure already exists. What is the refactoring gain in AND nodes? (Negative means rejected.)", stem, m),
		edakit.IntStr(gain), edakit.Sprintf("The factored form costs %d AND nodes and frees %d. Gain %d; refactoring accepts it only if positive (or zero with the zero-cost variant).", l-1, m, gain),
		map[string]string{"gain": edakit.IntStr(gain)})
}

// balancing: balanced depth of an AND supergate and the reduction against a chain.
func balancing(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 4, 9)
	which := kp(ctx)
	lv := make([]int, k)
	if which != 1 {
		for i := range lv {
			lv[i] = edakit.Int(0, 4)
		}
	}
	// Chain in the listed order.
	chain := lv[0]
	for i := 1; i < k; i++ {
		if lv[i] > chain {
			chain = lv[i]
		}
		chain++
	}
	h := append([]int(nil), lv...)
	for len(h) > 1 {
		sort.Ints(h)
		a, b := h[0], h[1]
		if b > a {
			a = b
		}
		h = append(h[2:], a+1)
	}
	bal := h[0]
	switch which {
	case 1:
		return edakit.Numeric(edakit.Sprintf("An AND supergate with %d inputs, all arriving at level 0, is implemented as a chain of two-input ANDs of depth %d. Balancing rebuilds it as a tree of two-input ANDs. What is the minimum depth?", k, k-1),
			edakit.IntStr(bal), edakit.Sprintf("A balanced binary tree over %d leaves has depth $\\lceil \\log_2 %d \\rceil = %d$.", k, k, bal),
			map[string]string{"depth": edakit.IntStr(bal)})
	case 2:
		return edakit.Numeric("An AND supergate has inputs arriving at levels "+intsCSV(lv)+". Balancing pairs the two lowest-level inputs repeatedly. What is the level of the supergate output after balancing?",
			edakit.IntStr(bal), edakit.Sprintf("Combining the two earliest signals each time (a Huffman-style merge with max + 1) gives output level %d, which is optimal.", bal),
			map[string]string{"depth": edakit.IntStr(bal)})
	}
	return edakit.Numeric("An AND supergate is built as a chain that absorbs inputs in the order given, with arrival levels "+intsCSV(lv)+" (the first two inputs feed the first AND). By how many levels does balancing reduce the output level?",
		edakit.IntStr(chain-bal), edakit.Sprintf("The chain output is at level %d. Pairing the earliest signals gives level %d. Reduction %d.", chain, bal, chain-bal),
		map[string]string{"chain": edakit.IntStr(chain), "balanced": edakit.IntStr(bal)})
}

// fraig: candidate equivalence classes from simulation signatures.
func fraig(ctx generator.GeneratorContext) generator.Problem {
	width := 8
	mask := uint64(1)<<width - 1
	type node struct {
		sig  uint64
		newb int
	}
	var nodes []node
	nb := edakit.Int(2, 3)
	for b := 0; b < nb; b++ {
		base := uint64(edakit.Int(1, int(mask)-1))
		baseNew := edakit.Int(0, 1)
		for c := 0; c < edakit.Int(2, 3); c++ {
			if edakit.Int(0, 2) == 0 {
				nodes = append(nodes, node{^base & mask, 1 - baseNew})
			} else {
				nodes = append(nodes, node{base, baseNew})
			}
		}
		if edakit.Int(0, 1) == 0 {
			// An impostor with the same signature but a different function.
			nodes = append(nodes, node{base, 1 - baseNew})
		}
	}
	for s := 0; s < edakit.Int(1, 2); s++ {
		nodes = append(nodes, node{uint64(edakit.Int(1, int(mask)-1)), edakit.Int(0, 1)})
	}
	nodes = edakit.Shuffle(nodes)
	canon := func(s uint64) uint64 {
		if s&1 == 1 {
			return ^s & mask
		}
		return s
	}
	classes := map[uint64]int{}
	for _, nd := range nodes {
		classes[canon(nd.sig)]++
	}
	var parts []string
	for i, nd := range nodes {
		parts = append(parts, edakit.Sprintf("n%d: %08b", i+1, nd.sig))
	}
	stem := "Random simulation of 8 input patterns gives the node signatures " + strings.Join(parts, ", ") + ". Nodes whose signatures are equal or bitwise complementary are candidate equivalences."
	multi, calls := 0, 0
	for _, c := range classes {
		if c >= 2 {
			multi++
			calls += c - 1
		}
	}
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+" How many candidate equivalence classes contain at least two nodes?",
			edakit.IntStr(multi), edakit.Sprintf("Normalize each signature so that its rightmost bit is 0 (complement it if needed) and group equal values. %d groups have two or more members.", multi),
			map[string]string{"classes": edakit.IntStr(multi)})
	case 2:
		return edakit.Numeric(stem+" A FRAIG builder proves each non-representative member against its class representative with one SAT call. How many SAT calls are needed?",
			edakit.IntStr(calls), edakit.Sprintf("Each class of size s needs s - 1 calls: %d in total.", calls),
			map[string]string{"sat_calls": edakit.IntStr(calls)})
	}
	var nbits []string
	refined := map[[2]uint64]int{}
	for i, nd := range nodes {
		nbits = append(nbits, edakit.Sprintf("n%d = %d", i+1, nd.newb))
		c := canon(nd.sig)
		nbv := uint64(nd.newb)
		if c != nd.sig {
			nbv ^= 1
		}
		refined[[2]uint64{c, nbv}]++
	}
	after := 0
	for _, c := range refined {
		if c >= 2 {
			after++
		}
	}
	return edakit.Numeric(stem+" A SAT call returns a counterexample, which is simulated as a ninth pattern and gives "+strings.Join(nbits, ", ")+". After refining by this bit, how many candidate classes still contain at least two nodes?",
		edakit.IntStr(after), edakit.Sprintf("Append the new bit to each signature (complemented along with the signature if the node was normalized) and regroup: %d classes remain.", after),
		map[string]string{"classes": edakit.IntStr(after)})
}

// lutMapping: depth-optimal K-LUT mapping by cut labels.
func lutMapping(ctx generator.GeneratorContext) generator.Problem {
	which := kp(ctx)
	k := 3
	if which == 2 {
		k = 4
	}
	g := randAIG(k+2, edakit.Scaled(ctx, 7, 10))
	lab := g.lutDepth(k)
	lv := g.level()
	depth, aDepth := 0, 0
	for _, o := range g.outs {
		if lab[o] > depth {
			depth = lab[o]
		}
		if lv[o] > aDepth {
			aDepth = lv[o]
		}
	}
	if which == 3 {
		return edakit.Numeric("In the AIG "+g.desc()+" Map it into 3-input LUTs for minimum depth. By how many levels is the LUT depth smaller than the AIG depth (the largest AND level)?",
			edakit.IntStr(aDepth-depth), edakit.Sprintf("The AIG depth is %d. Labeling each node with the best 3-feasible cut gives LUT depth %d. Reduction %d.", aDepth, depth, aDepth-depth),
			map[string]string{"aig_depth": edakit.IntStr(aDepth), "lut_depth": edakit.IntStr(depth)})
	}
	return edakit.Numeric("In the AIG "+g.desc()+edakit.Sprintf(" Map it into %d-input LUTs for minimum depth: label each primary input 0 and each AND node with the minimum over its nontrivial %d-feasible cuts of 1 + the largest leaf label. What is the LUT depth of the mapping (the largest output label)?", k, k),
		edakit.IntStr(depth), edakit.Sprintf("Labels are computed in topological order, each from the best cut. The outputs reach label %d.", depth),
		map[string]string{"lut_depth": edakit.IntStr(depth)})
}
