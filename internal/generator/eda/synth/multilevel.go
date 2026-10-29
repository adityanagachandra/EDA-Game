package synth

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Algebraic helpers. A positive cube is a bitmask over the letters below; an algebraic
// expression is a list of distinct cubes, none containing another.

var mlNames = []string{"a", "b", "c", "d", "e", "f", "g", "h", "x", "y", "z", "w"}

type sop []uint32

func pcubeStr(c uint32) string {
	if c == 0 {
		return "1"
	}
	var b strings.Builder
	for i, nm := range mlNames {
		if c>>i&1 == 1 {
			b.WriteString(nm)
		}
	}
	return b.String()
}

func (f sop) norm() sop {
	m := map[uint32]bool{}
	var out sop
	for _, c := range f {
		if !m[c] {
			m[c] = true
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if popcountU(out[i]) != popcountU(out[j]) {
			return popcountU(out[i]) > popcountU(out[j])
		}
		return pcubeStr(out[i]) < pcubeStr(out[j])
	})
	return out
}

func (f sop) str() string {
	parts := make([]string, len(f))
	for i, c := range f {
		parts[i] = pcubeStr(c)
	}
	return strings.Join(parts, " + ")
}

func (f sop) key() string { return f.norm().str() }

func (f sop) lits() int {
	l := 0
	for _, c := range f {
		l += popcountU(c)
	}
	return l
}

func (f sop) support() uint32 {
	var s uint32
	for _, c := range f {
		s |= c
	}
	return s
}

// algebraic reports that no cube contains another and there are no duplicates.
func (f sop) algebraic() bool {
	for i, a := range f {
		for j, b := range f {
			if i != j && subset(a, b) {
				return false
			}
		}
	}
	return true
}

func (f sop) cubeFree() bool {
	if len(f) < 2 {
		return false
	}
	common := ^uint32(0)
	for _, c := range f {
		common &= c
	}
	return common == 0
}

func divCube(f sop, c uint32) sop {
	var q sop
	for _, x := range f {
		if subset(c, x) {
			q = append(q, x&^c)
		}
	}
	return q
}

func mul(a, b sop) sop {
	var out sop
	for _, x := range a {
		for _, y := range b {
			out = append(out, x|y)
		}
	}
	return out
}

// weakDivide returns quotient and remainder of algebraic division f / d.
func weakDivide(f, d sop) (sop, sop) {
	var q sop
	for i, di := range d {
		qi := divCube(f, di)
		if i == 0 {
			q = qi
			continue
		}
		set := map[uint32]bool{}
		for _, x := range qi {
			set[x] = true
		}
		var keep sop
		for _, x := range q {
			if set[x] {
				keep = append(keep, x)
			}
		}
		q = keep
	}
	prod := map[uint32]bool{}
	for _, x := range mul(q, d) {
		prod[x] = true
	}
	var r sop
	for _, x := range f {
		if !prod[x] {
			r = append(r, x)
		}
	}
	return q, r
}

// kernelsOf returns the distinct kernels of f with, for each, its co-kernels.
func kernelsOf(f sop) (map[string]sop, map[string][]uint32) {
	ks := map[string]sop{}
	cok := map[string][]uint32{}
	sup := f.support()
	var bitsIdx []int
	for i := 0; i < 32; i++ {
		if sup>>i&1 == 1 {
			bitsIdx = append(bitsIdx, i)
		}
	}
	for s := 0; s < 1<<len(bitsIdx); s++ {
		var c uint32
		for j, b := range bitsIdx {
			if s>>j&1 == 1 {
				c |= 1 << b
			}
		}
		q := divCube(f, c)
		if q.cubeFree() {
			k := q.key()
			ks[k] = q
			cok[k] = append(cok[k], c)
		}
	}
	return ks, cok
}

// formula is a read-once factored form over distinct letters.
type formula struct {
	op   byte // '+', '*', or 0 for a leaf
	leaf int
	kids []*formula
}

func randFormula(leaves []int) *formula {
	if len(leaves) == 1 {
		return &formula{leaf: leaves[0]}
	}
	k := edakit.Int(1, len(leaves)-1)
	l, r := randFormula(leaves[:k]), randFormula(leaves[k:])
	op := byte('+')
	if edakit.Int(0, 1) == 1 {
		op = '*'
	}
	node := &formula{op: op}
	for _, ch := range []*formula{l, r} {
		if ch.op == op {
			node.kids = append(node.kids, ch.kids...)
		} else {
			node.kids = append(node.kids, ch)
		}
	}
	return node
}

func (f *formula) str() string {
	if f.op == 0 {
		return mlNames[f.leaf]
	}
	parts := make([]string, len(f.kids))
	for i, k := range f.kids {
		s := k.str()
		if f.op == '*' && k.op == '+' {
			s = "(" + s + ")"
		}
		parts[i] = s
	}
	if f.op == '+' {
		return strings.Join(parts, " + ")
	}
	return strings.Join(parts, "")
}

func (f *formula) leaves() int {
	if f.op == 0 {
		return 1
	}
	n := 0
	for _, k := range f.kids {
		n += k.leaves()
	}
	return n
}

func (f *formula) sop() sop {
	if f.op == 0 {
		return sop{1 << f.leaf}
	}
	out := f.kids[0].sop()
	for _, k := range f.kids[1:] {
		if f.op == '+' {
			out = append(out, k.sop()...)
		} else {
			out = mul(out, k.sop())
		}
	}
	return out
}

// randReadOnce draws a factored form whose SOP has between lo and hi cubes and that is not
// itself a plain SOP (so factoring matters).
func randReadOnce(nLeaves, lo, hi int) *formula {
	for {
		perm := edakit.Shuffle(seq(8))[:nLeaves]
		sort.Ints(perm)
		f := randFormula(edakit.Shuffle(perm))
		s := f.sop()
		if len(s) >= lo && len(s) <= hi && s.lits() > f.leaves() {
			return f
		}
	}
}

func randPCube(letters []int, lo, hi int) uint32 {
	k := edakit.Int(lo, hi)
	var c uint32
	for _, i := range edakit.Shuffle(letters)[:k] {
		c |= 1 << i
	}
	return c
}

// booleanNetwork: total literals, depth, literals after collapsing a node.
func booleanNetwork(ctx generator.GeneratorContext) generator.Problem {
	// x over {a,b,c}, y over {c,d}, z over {x,y,e}.
	for {
		x := sop{randPCube([]int{0, 1, 2}, 1, 2), randPCube([]int{0, 1, 2}, 1, 2)}.norm()
		y := sop{randPCube([]int{2, 3}, 1, 2), randPCube([]int{2, 3, 5}, 1, 2)}.norm()
		zLetters := []int{8, 9, 4}
		z := sop{1<<8 | randPCube([]int{9, 4}, 0, 1), randPCube(zLetters, 1, 2), randPCube(zLetters, 1, 2)}.norm()
		if !x.algebraic() || !y.algebraic() || !z.algebraic() || len(x) < 2 || len(z) < 2 || z.support()&(1<<9) == 0 {
			continue
		}
		withW := ctx.Difficulty > 0.5
		w := sop{1<<10 | 1<<5, 1 << 9}
		total := x.lits() + y.lits() + z.lits()
		desc := "A Boolean network has primary inputs a to f and internal nodes $x = " + x.str() + "$, $y = " + y.str() + "$, $z = " + z.str() + "$"
		if withW {
			total += w.lits()
			desc += ", $w = " + w.str() + "$; the primary outputs are z and w."
		} else {
			desc += "; the primary output is z."
		}
		switch kp(ctx) {
		case 1:
			return edakit.Numeric(desc+" What is the total literal count of the network (the sum over nodes of the literals in each node's SOP)?",
				edakit.IntStr(total), edakit.Sprintf("Add the SOP literal counts node by node, counting an internal-node name like x as a literal where it is used. Total %d.", total),
				map[string]string{"literals": edakit.IntStr(total)})
		case 2:
			d := 2
			if withW {
				d = 3
			}
			return edakit.Numeric(desc+" Give each primary input level 0 and each node one more than its highest fanin. What is the largest level in the network?",
				edakit.IntStr(d), edakit.Sprintf("x and y depend only on inputs (level 1), z depends on x and y (level 2)%s. The network depth is %d.", map[bool]string{true: ", and w depends on z (level 3)", false: ""}[withW], d),
				map[string]string{"depth": edakit.IntStr(d)})
		}
		// Collapse x into z.
		var nz sop
		for _, c := range z {
			if c>>8&1 == 1 {
				for _, xc := range x {
					nz = append(nz, c&^(1<<8)|xc)
				}
			} else {
				nz = append(nz, c)
			}
		}
		after := total - x.lits() - z.lits() + nz.lits()
		return edakit.Numeric(desc+" Node x fans out only to z. Eliminate x by substituting its SOP into z and multiplying out, without further simplification. What is the network's total literal count afterwards?",
			edakit.IntStr(after), edakit.Sprintf("z becomes $%s$ with %d literals, and x disappears. The total goes from %d to %d.", nz.str(), nz.lits(), total, after),
			map[string]string{"literals_after": edakit.IntStr(after)})
	}
}

// divisionInstance draws F = D*Q + R with disjoint supports for D and Q.
func divisionInstance() (f, d, q, r sop) {
	for {
		d = sop{randPCube([]int{0, 1, 2}, 1, 2), randPCube([]int{0, 1, 2}, 1, 2)}.norm()
		if len(d) < 2 || !d.algebraic() {
			continue
		}
		nq := edakit.Int(1, 3)
		q = nil
		for i := 0; i < nq; i++ {
			q = append(q, randPCube([]int{3, 4, 5}, 1, 2))
		}
		q = q.norm()
		if !q.algebraic() {
			continue
		}
		nr := edakit.Int(1, 2)
		r = nil
		for i := 0; i < nr; i++ {
			r = append(r, randPCube([]int{0, 1, 3, 6, 7}, 1, 3))
		}
		f = append(mul(d, q), r...).norm()
		if len(f) != len(d)*len(q)+len(r.norm()) || !f.algebraic() {
			continue
		}
		return f, d, q, r.norm()
	}
}

// algebraicDivision: F/c for a cube, F/D, remainder literals.
func algebraicDivision(ctx generator.GeneratorContext) generator.Problem {
	f, d, _, _ := divisionInstance()
	switch kp(ctx) {
	case 1:
		c := edakit.Pick(d)
		q := divCube(f, c)
		return edakit.Numeric(edakit.Sprintf("Divide $F = %s$ algebraically by the single cube $%s$. How many cubes are in the quotient $F / %s$?", f.str(), pcubeStr(c), pcubeStr(c)),
			edakit.IntStr(len(q)), edakit.Sprintf("Keep each cube of F that contains %s and strike %s from it: $%s$, %d cubes.", pcubeStr(c), pcubeStr(c), q.str(), len(q)),
			map[string]string{"quotient": q.str()})
	case 2:
		q, _ := weakDivide(f, d)
		return edakit.Numeric(edakit.Sprintf("Perform weak (algebraic) division of $F = %s$ by $D = %s$. How many cubes are in the quotient Q?", f.str(), d.str()),
			edakit.IntStr(len(q)), edakit.Sprintf("Intersect the cube quotients $F/d_i$ over the cubes of D. The quotient is $Q = %s$, %d cubes.", q.str(), len(q)),
			map[string]string{"quotient": q.str()})
	}
	q, r := weakDivide(f, d)
	return edakit.Numeric(edakit.Sprintf("Divide $F = %s$ by $D = %s$ algebraically, so that $F = DQ + R$. How many literals are in the remainder R?", f.str(), d.str()),
		edakit.IntStr(r.lits()), edakit.Sprintf("The quotient is $%s$; removing the cubes of DQ from F leaves $R = %s$ with %d literals.", q.str(), r.str(), r.lits()),
		map[string]string{"remainder": r.str()})
}

func kernelInstance() sop {
	for {
		var f sop
		if edakit.Int(0, 1) == 0 {
			f = randReadOnce(edakit.Int(5, 7), 3, 8).sop().norm()
		} else {
			for i := 0; i < edakit.Int(4, 6); i++ {
				f = append(f, randPCube(seq(6), 1, 3))
			}
			f = f.norm()
		}
		if f.algebraic() && len(f) >= 3 && len(f) <= 8 {
			ks, _ := kernelsOf(f)
			if len(ks) >= 2 {
				return f
			}
		}
	}
}

// kernels: number of kernels, number of co-kernels, level-0 kernels.
func kernels(ctx generator.GeneratorContext) generator.Problem {
	f := kernelInstance()
	ks, cok := kernelsOf(f)
	var names []string
	for k := range ks {
		names = append(names, "$"+k+"$")
	}
	sort.Strings(names)
	stem := "$F = " + f.str() + "$"
	switch kp(ctx) {
	case 1:
		return edakit.Numeric("A kernel of F is a cube-free quotient $F/c$ (at least two cubes, no literal common to all) for some cube c, and F itself counts if it is cube-free. How many distinct kernels does "+stem+" have?",
			edakit.IntStr(len(ks)), "The kernels are "+strings.Join(names, ", ")+edakit.Sprintf(": %d in all.", len(ks)),
			map[string]string{"kernels": edakit.IntStr(len(ks))})
	case 2:
		n := 0
		for _, cs := range cok {
			n += len(cs)
		}
		return edakit.Numeric("Count the co-kernels of "+stem+": the cubes c (including the cube 1) for which $F/c$ is a kernel. Several co-kernels can give the same kernel. How many co-kernels are there?",
			edakit.IntStr(n), edakit.Sprintf("Dividing by every cube built from the support and keeping the cube-free quotients gives %d co-kernels for %d distinct kernels.", n, len(ks)),
			map[string]string{"cokernels": edakit.IntStr(n)})
	}
	lvl0 := 0
	for _, k := range ks {
		sub, _ := kernelsOf(k)
		if len(sub) == 1 {
			lvl0++
		}
	}
	return edakit.Numeric("A level-0 kernel has no kernels except itself. How many level-0 kernels does "+stem+" have?",
		edakit.IntStr(lvl0), "The kernels are "+strings.Join(names, ", ")+edakit.Sprintf("; %d of them contain no smaller kernel.", lvl0),
		map[string]string{"level0": edakit.IntStr(lvl0)})
}

// factoring: SOP cubes, SOP literals, literal saving of a factored form.
func factoring(ctx generator.GeneratorContext) generator.Problem {
	f := randReadOnce(edakit.Scaled(ctx, 5, 8), 3, 12)
	s := f.sop()
	fl := f.leaves()
	stem := "$F = " + f.str() + "$"
	switch kp(ctx) {
	case 1:
		return edakit.Numeric("Multiply the factored form "+stem+" out into a sum of products. How many cubes does the SOP have?",
			edakit.IntStr(len(s)), edakit.Sprintf("A product multiplies cube counts and a sum adds them. The SOP has %d cubes.", len(s)),
			map[string]string{"cubes": edakit.IntStr(len(s))})
	case 2:
		return edakit.Numeric("Multiply the factored form "+stem+" out into a sum of products. How many literals does the SOP have?",
			edakit.IntStr(s.lits()), edakit.Sprintf("The SOP is $%s$ with %d literals.", s.str(), s.lits()),
			map[string]string{"sop_literals": edakit.IntStr(s.lits())})
	}
	return edakit.Numeric("The SOP $"+s.str()+"$ can be written as the factored form "+stem+". How many literals does factoring save?",
		edakit.IntStr(s.lits()-fl), edakit.Sprintf("The SOP has %d literals and the factored form %d, one per leaf. Saving %d.", s.lits(), fl, s.lits()-fl),
		map[string]string{"saving": edakit.IntStr(s.lits() - fl)})
}

// extraction: literals before, after and saved when a common kernel becomes a node.
func extraction(ctx generator.GeneratorContext) generator.Problem {
	for {
		d := sop{randPCube([]int{0, 1, 2}, 1, 2), randPCube([]int{0, 1, 2}, 1, 2)}.norm()
		if len(d) < 2 || !d.algebraic() {
			continue
		}
		nf := 2
		if ctx.Difficulty > 0.7 {
			nf = 3
		}
		var fs []sop
		ok := true
		for i := 0; i < nf; i++ {
			q := sop{randPCube([]int{3, 4, 5}, 0, 2)}
			if edakit.Int(0, 1) == 1 {
				q = append(q, randPCube([]int{3, 4, 5}, 1, 2))
			}
			q = q.norm()
			r := sop{randPCube([]int{6, 7, 3}, 1, 2)}
			f := append(mul(d, q), r...).norm()
			if !q.algebraic() || !f.algebraic() {
				ok = false
				break
			}
			fs = append(fs, f)
		}
		if !ok {
			continue
		}
		before, after := 0, d.lits()
		var desc []string
		for i, f := range fs {
			before += f.lits()
			q, r := weakDivide(f, d)
			if len(q) == 0 {
				ok = false
				break
			}
			for _, c := range q {
				after += popcountU(c) + 1
			}
			after += r.lits()
			desc = append(desc, edakit.Sprintf("$F_%d = %s$", i+1, f.str()))
		}
		if !ok {
			continue
		}
		stem := "The functions " + strings.Join(desc, ", ") + " share the divisor $D = " + d.str() + "$."
		switch kp(ctx) {
		case 1:
			return edakit.Numeric(stem+" What is their total literal count before extraction?",
				edakit.IntStr(before), edakit.Sprintf("Summing the literals of every cube gives %d.", before),
				map[string]string{"before": edakit.IntStr(before)})
		case 2:
			return edakit.Numeric(stem+" Extract D as a new node X and rewrite each function as $XQ + R$. What is the total literal count afterwards, including the literals of X itself?",
				edakit.IntStr(after), edakit.Sprintf("Each quotient cube costs its literals plus one for X; X costs %d. Total %d.", d.lits(), after),
				map[string]string{"after": edakit.IntStr(after)})
		}
		return edakit.Numeric(stem+" How many literals are saved by extracting D as a new node X shared by all of them?",
			edakit.IntStr(before-after), edakit.Sprintf("Before: %d literals. After: %d, counting X. Saving %d.", before, after, before-after),
			map[string]string{"saving": edakit.IntStr(before - after)})
	}
}

// cubeExtraction: occurrences of a cube, literals after extraction, best single-cube value.
func cubeExtraction(ctx generator.GeneratorContext) generator.Problem {
	letters := seq(6)
	for {
		var fs []sop
		nf := 2
		var all sop
		for i := 0; i < nf; i++ {
			var f sop
			for j := 0; j < edakit.Int(3, 4); j++ {
				f = append(f, randPCube(letters, 2, 4))
			}
			f = f.norm()
			if !f.algebraic() {
				f = nil
				break
			}
			fs = append(fs, f)
			all = append(all, f...)
		}
		if len(fs) < nf {
			continue
		}
		occ := func(c uint32) int {
			m := 0
			for _, x := range all {
				if subset(c, x) {
					m++
				}
			}
			return m
		}
		best, bestC := -1<<30, uint32(0)
		var cands []uint32
		for c := uint32(1); c < 1<<6; c++ {
			k := popcountU(c)
			if k < 2 {
				continue
			}
			m := occ(c)
			if m < 2 {
				continue
			}
			cands = append(cands, c)
			v := m*k - m - k
			if v > best {
				best, bestC = v, c
			}
		}
		if len(cands) == 0 {
			continue
		}
		stem := edakit.Sprintf("Consider $F_1 = %s$ and $F_2 = %s$.", fs[0].str(), fs[1].str())
		c := edakit.Pick(cands)
		m := occ(c)
		k := popcountU(c)
		switch kp(ctx) {
		case 1:
			return edakit.Numeric(stem+edakit.Sprintf(" How many cubes, across both functions, contain the cube $%s$?", pcubeStr(c)),
				edakit.IntStr(m), edakit.Sprintf("Cube %s divides %d cubes; that is the number of places a new node for it would be used.", pcubeStr(c), m),
				map[string]string{"occurrences": edakit.IntStr(m)})
		case 2:
			tot := all.lits() - m*k + m + k
			return edakit.Numeric(stem+edakit.Sprintf(" Extract the cube $%s$ as a new node $X = %s$ and substitute it everywhere it divides a cube. What is the total literal count afterwards, including X?", pcubeStr(c), pcubeStr(c)),
				edakit.IntStr(tot), edakit.Sprintf("Each of the %d cubes trades %d literals for one, and X adds %d: %d - %d + %d + %d = %d.", m, k, k, all.lits(), m*k, m, k, tot),
				map[string]string{"after": edakit.IntStr(tot)})
		}
		return edakit.Numeric(stem+" Over all cubes of two or more literals, what is the largest literal saving obtainable by extracting a single common cube as a new node? A cube of k literals used in m cubes saves $mk - m - k$.",
			edakit.IntStr(best), edakit.Sprintf("The best choice is $%s$, used %d times, saving %d literals.", pcubeStr(bestC), occ(bestC), best),
			map[string]string{"best_cube": pcubeStr(bestC)})
	}
}

// elimination: the literal value of eliminating a node.
func elimination(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n, l := edakit.Int(1, 5), edakit.Int(2, 7)
		v := n*l - n - l
		return edakit.Numeric(edakit.Sprintf("Node t has a factored form with %d literals and appears as a literal in %d fanout expressions. If t is eliminated, each use is replaced by t's factored form. By how many literals does the network's factored literal count change? (Negative means it shrinks.)", l, n),
			edakit.IntStr(v), edakit.Sprintf("After: %d copies of %d literals = %d. Before: %d for t plus %d uses = %d. Change %d, which is $(n-1)(l-1) - 1$.", n, l, n*l, l, n, l+n, v),
			map[string]string{"value": edakit.IntStr(v)})
	case 2:
		thr := edakit.Int(-1, 3)
		cnt := 0
		var parts []string
		for i := 0; i < 5; i++ {
			n, l := edakit.Int(1, 4), edakit.Int(1, 6)
			v := n*l - n - l
			if v <= thr {
				cnt++
			}
			parts = append(parts, edakit.Sprintf("t%d (%d literals, %d fanout uses)", i+1, l, n))
		}
		return edakit.Numeric("An eliminate pass collapses every node whose value $nl - n - l$ (literal increase when collapsed) is at most "+edakit.IntStr(thr)+". The candidates are "+strings.Join(parts, ", ")+". How many nodes are eliminated?",
			edakit.IntStr(cnt), edakit.Sprintf("Compute the value of each node and compare it with the threshold %d: %d nodes qualify.", thr, cnt),
			map[string]string{"eliminated": edakit.IntStr(cnt)})
	}
	// Exact SOP: x = SOP with c cubes, L literals; fanout cubes q_j x.
	x := sop{randPCube([]int{0, 1, 2}, 1, 2), randPCube([]int{0, 1, 2}, 1, 2), randPCube([]int{0, 1, 2}, 1, 2)}.norm()
	for !x.algebraic() || len(x) < 2 {
		x = sop{randPCube([]int{0, 1, 2}, 1, 2), randPCube([]int{0, 1, 2}, 1, 2)}.norm()
	}
	nu := edakit.Int(1, 3)
	var uses []string
	before := x.lits()
	after := 0
	for j := 0; j < nu; j++ {
		q := randPCube([]int{3, 4, 5, 6}, 1, 2)
		uses = append(uses, "$"+pcubeStr(q|1<<8)+"$")
		lq := popcountU(q)
		before += lq + 1
		after += lq*len(x) + x.lits()
	}
	return edakit.Numeric("Node $x = "+x.str()+"$ is used in the cubes "+strings.Join(uses, ", ")+", each in a different fanout node. Eliminate x by multiplying it into each of those cubes. By how many SOP literals does the total change (x's own literals included before; negative means a saving)?",
		edakit.IntStr(after-before), edakit.Sprintf("Before: %d literals (x plus its fanout cubes). Each cube qx becomes %d cubes, $l_q \\cdot %d + %d$ literals. After: %d. Change %d.", before, len(x), len(x), x.lits(), after, after-before),
		map[string]string{"change": edakit.IntStr(after - before)})
}
