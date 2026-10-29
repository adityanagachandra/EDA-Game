package synth

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// kp returns the knowledge point to serve: the requested one, or a random one when unguided.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

func ceilLog2(k int) int {
	d := 0
	for (1 << d) < k {
		d++
	}
	return d
}

// costMetrics: literal count, two-input gate count and two-level depth of a random SOP.
func costMetrics(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 5)
	k := edakit.Scaled(ctx, 3, 5)
	var cs []cube
	seen := map[cube]bool{}
	for len(cs) < k {
		c := randCube(n, 2)
		if !seen[c] {
			seen[c] = true
			cs = append(cs, c)
		}
	}
	lits := coverLits(cs)
	expr := "$f = " + sopTex(n, cs) + "$"
	switch kp(ctx) {
	case 1:
		return edakit.Numeric("A synthesis tool scores the cover "+expr+" by its literal count. How many literals does it have?",
			edakit.IntStr(lits), edakit.Sprintf("Count every variable occurrence, complemented or not, across the %d cubes: %d literals.", k, lits),
			map[string]string{"literals": edakit.IntStr(lits)})
	case 2:
		g := 0
		for _, c := range cs {
			g += c.lits() - 1
		}
		g += k - 1
		return edakit.Numeric("Implement "+expr+" as a two-level AND-OR circuit using only two-input AND and OR gates, with complemented inputs available for free. How many two-input gates are needed?",
			edakit.IntStr(g), edakit.Sprintf("A cube with l literals needs l - 1 two-input ANDs and %d products need %d two-input ORs. Total %d, which equals literals minus 1 = %d - 1.", k, k-1, g, lits),
			map[string]string{"gates": edakit.IntStr(g)})
	}
	maxl := 0
	for _, c := range cs {
		if c.lits() > maxl {
			maxl = c.lits()
		}
	}
	d := ceilLog2(maxl) + ceilLog2(k)
	return edakit.Numeric("Build "+expr+" from two-input gates, with each product and the final sum as balanced trees and complemented inputs free. What is the depth of the circuit in gate levels?",
		edakit.IntStr(d), edakit.Sprintf("The widest product has %d literals, a balanced AND tree of depth %d. The sum of %d products is an OR tree of depth %d. Depth %d.", maxl, ceilLog2(maxl), k, ceilLog2(k), d),
		map[string]string{"depth": edakit.IntStr(d)})
}

// cubeNotation: intersection size, distance and consensus in positional cube notation.
func cubeNotation(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 5)
	switch kp(ctx) {
	case 1:
		a, b := randCube(n, 1), randCube(n, 1)
		inter := popcount(a.table(n) & b.table(n))
		return edakit.Numeric(edakit.Sprintf("In positional cube notation over variables %s, cube A is %s and cube B is %s. How many minterms lie in the intersection of A and B?",
			strings.Join(varNames[:n], ", "), a.pos(n), b.pos(n)),
			edakit.IntStr(inter), edakit.Sprintf("Intersect variable by variable: a 0 against a 1 makes the intersection empty, otherwise the more specific symbol wins. The result covers %d minterms.", inter),
			map[string]string{"intersection": edakit.IntStr(inter)})
	case 2:
		a, b := randCube(n, 2), randCube(n, 2)
		d := 0
		for i := 0; i < n; i++ {
			p := uint32(1) << i
			if a.mask&b.mask&p != 0 && (a.val^b.val)&p != 0 {
				d++
			}
		}
		return edakit.Numeric(edakit.Sprintf("Cubes %s and %s are written over variables %s. What is the distance between them (the number of variables in which one has 0 and the other has 1)?",
			a.pos(n), b.pos(n), strings.Join(varNames[:n], ", ")),
			edakit.IntStr(d), edakit.Sprintf("Only positions where both cubes are specified and disagree count. There are %d. Distance 0 means they intersect; distance 1 means they have a consensus.", d),
			map[string]string{"distance": edakit.IntStr(d)})
	}
	// Two cubes at distance exactly 1.
	for {
		a, b := randCube(n, 2), randCube(n, 2)
		conf := a.mask & b.mask & (a.val ^ b.val)
		if popcountU(conf) != 1 {
			continue
		}
		mask := (a.mask | b.mask) &^ conf
		val := (a.val | b.val) & mask
		cons := cube{mask, val}
		l := cons.lits()
		return edakit.Numeric(edakit.Sprintf("Cubes %s and %s (variables %s) conflict in exactly one variable. How many literals does their consensus cube have?",
			a.pos(n), b.pos(n), strings.Join(varNames[:n], ", ")),
			edakit.IntStr(l), edakit.Sprintf("The consensus drops the conflicting variable and intersects the rest, giving %s with %d literals.", cons.pos(n), l),
			map[string]string{"consensus": cons.pos(n)})
	}
}

func popcountU(x uint32) int { return popcount(uint64(x)) }

// kmap: largest implicant through a minterm, minimum SOP product count, minimum literals with DCs.
func kmap(ctx generator.GeneratorContext) generator.Problem {
	n := 3
	if ctx.Difficulty > 0.4 {
		n = 4
	}
	full := 1 << n
	switch kp(ctx) {
	case 1:
		on := randFunc(n, full/4+1, full/2+2)
		ms := mintermList(on, n)
		m := edakit.Pick(ms)
		best := 0
		for _, c := range allCubes(n) {
			if c.contains(m) && c.table(n)&^on == 0 && c.size(n) > best {
				best = c.size(n)
			}
		}
		return edakit.Numeric("On the Karnaugh map of "+sigma("f", on, 0, n)+edakit.Sprintf(", what is the largest group of 1-cells (a power-of-two rectangle, wrapping allowed) that contains minterm %d? Give its number of cells.", m),
			edakit.IntStr(best), edakit.Sprintf("The largest implicant containing m%d covers %d minterms; every cell of it is a 1.", m, best),
			map[string]string{"group": edakit.IntStr(best)})
	case 2:
		on := randFunc(n, 3, full-3)
		cov := minCover(n, on, primes(n, on, 0))
		return edakit.Numeric("Minimize "+sigma("f", on, 0, n)+" with a Karnaugh map. How many product terms does a minimum sum-of-products cover have?",
			edakit.IntStr(len(cov)), edakit.Sprintf("One minimum cover is $f = %s$, with %d products.", sopTex(n, cov), len(cov)),
			map[string]string{"products": edakit.IntStr(len(cov))})
	}
	on := randFunc(n, 3, full/2)
	dc := randFunc(n, 1, 3) &^ on
	if dc == 0 {
		dc = 1 << uint(edakit.Pick(mintermList(^on&(1<<full-1), n)))
	}
	cov := minCover(n, on, primes(n, on, dc))
	l := coverLits(cov)
	return edakit.Numeric("Minimize "+sigma("f", on, dc, n)+", using the don't cares wherever they help. How many literals does a minimum sum-of-products cover have (fewest products first, then fewest literals)?",
		edakit.IntStr(l), edakit.Sprintf("A minimum cover is $f = %s$: %d products, %d literals.", sopTex(n, cov), len(cov), l),
		map[string]string{"literals": edakit.IntStr(l)})
}

// quineMcCluskey: merged pairs in the first column, number of primes, primes with DCs.
func quineMcCluskey(ctx generator.GeneratorContext) generator.Problem {
	n := 4
	if ctx.Difficulty > 0.75 && ctx.KPOrdinal != 1 {
		n = 5
	}
	full := 1 << n
	switch kp(ctx) {
	case 1:
		on := randFunc(n, 4, full/2+1)
		pairs := 0
		for _, c := range allCubes(n) {
			if c.lits() == n-1 && c.table(n)&^on == 0 {
				pairs++
			}
		}
		return edakit.Numeric("Run the first Quine-McCluskey merge step on "+sigma("f", on, 0, n)+": combine every pair of ON minterms that differ in exactly one bit. How many merged one-dash implicants result?",
			edakit.IntStr(pairs), edakit.Sprintf("Each pair of adjacent ON minterms gives one implicant with a single dash. There are %d such pairs.", pairs),
			map[string]string{"pairs": edakit.IntStr(pairs)})
	case 2:
		on := randFunc(n, 4+(n-4)*5, full/2+2)
		ps := primes(n, on, 0)
		return edakit.Numeric("Apply the Quine-McCluskey method to "+sigma("f", on, 0, n)+". How many prime implicants does f have?",
			edakit.IntStr(len(ps)), edakit.Sprintf("Merging until no merge is possible leaves the unmerged implicants: $%s$, which is %d primes.", sopTex(n, ps), len(ps)),
			map[string]string{"primes": edakit.IntStr(len(ps))})
	}
	on := randFunc(n, 4, full/2)
	dc := randFunc(n, 2, 4) &^ on
	ps := primes(n, on, dc)
	return edakit.Numeric("Apply Quine-McCluskey to "+sigma("f", on, dc, n)+", merging don't cares as if they were 1s. How many prime implicants cover at least one ON minterm (the rows of the prime implicant table)?",
		edakit.IntStr(len(ps)), edakit.Sprintf("Primes of ON plus DC that touch the ON set: $%s$, which is %d rows.", sopTex(n, ps), len(ps)),
		map[string]string{"primes": edakit.IntStr(len(ps))})
}

// primeTable: column count for a minterm, essential primes, uncovered after essentials.
func primeTable(ctx generator.GeneratorContext) generator.Problem {
	n := 4
	on := randFunc(n, 5, 10)
	ps := primes(n, on, 0)
	es := essentials(n, on, ps)
	switch kp(ctx) {
	case 1:
		m := edakit.Pick(mintermList(on, n))
		cnt := 0
		for _, p := range ps {
			if p.contains(m) {
				cnt++
			}
		}
		return edakit.Numeric(edakit.Sprintf("The primes of %s are %s. In the prime implicant table, how many primes cover the column for minterm %d?", sigma("f", on, 0, n), "$"+sopTex(n, ps)+"$", m),
			edakit.IntStr(cnt), edakit.Sprintf("Minterm %d lies in %d of the primes. A column with exactly one mark identifies an essential prime.", m, cnt),
			map[string]string{"column": edakit.IntStr(cnt)})
	case 2:
		esText := "none"
		if len(es) > 0 {
			esText = "$" + sopTex(n, es) + "$"
		}
		return edakit.Numeric("Build the prime implicant table of "+sigma("f", on, 0, n)+". How many essential prime implicants does f have?",
			edakit.IntStr(len(es)), edakit.Sprintf("The primes are %s. The essential ones, each the only cover of some minterm, are %s: %d of them.", "$"+sopTex(n, ps)+"$", esText, len(es)),
			map[string]string{"essentials": edakit.IntStr(len(es))})
	}
	var cov uint64
	for _, e := range es {
		cov |= e.table(n)
	}
	left := popcount(on &^ cov)
	return edakit.Numeric("For "+sigma("f", on, 0, n)+", select every essential prime implicant. How many ON minterms remain uncovered after that?",
		edakit.IntStr(left), edakit.Sprintf("The essential primes (%d of them) cover all but %d ON minterms; those form the residual covering problem.", len(es), left),
		map[string]string{"uncovered": edakit.IntStr(left)})
}

// coverTable is a unate covering matrix: rows[i] is the set of columns row i covers.
type coverTable struct {
	rows []uint32
	cols int
}

func randCoverTable(nr, nc int) coverTable {
	for {
		t := coverTable{rows: make([]uint32, nr), cols: nc}
		for i := range t.rows {
			for t.rows[i] == 0 {
				for j := 0; j < nc; j++ {
					if edakit.Int(0, 99) < 38 {
						t.rows[i] |= 1 << j
					}
				}
			}
		}
		var all uint32
		for _, r := range t.rows {
			all |= r
		}
		if all == 1<<nc-1 {
			return t
		}
	}
}

func (t coverTable) text() string {
	parts := make([]string, len(t.rows))
	for i, r := range t.rows {
		var cs []string
		for j := 0; j < t.cols; j++ {
			if r>>j&1 == 1 {
				cs = append(cs, edakit.Sprintf("c%d", j+1))
			}
		}
		parts[i] = edakit.Sprintf("P%d covers {%s}", i+1, strings.Join(cs, ", "))
	}
	return strings.Join(parts, "; ")
}

// minRows returns the minimum number of rows (or minimum total weight when w != nil) that
// cover every column.
func (t coverTable) minRows(w []int) int {
	best := 1 << 30
	nr := len(t.rows)
	full := uint32(1)<<t.cols - 1
	for s := 0; s < 1<<nr; s++ {
		var cov uint32
		cost := 0
		for i := 0; i < nr; i++ {
			if s>>i&1 == 1 {
				cov |= t.rows[i]
				if w == nil {
					cost++
				} else {
					cost += w[i]
				}
			}
		}
		if cov == full && cost < best {
			best = cost
		}
	}
	return best
}

func subset(a, b uint32) bool { return a&^b == 0 }

// reduce runs essential selection, row dominance and column dominance to a fixpoint and
// returns the number of columns left in the cyclic core.
func (t coverTable) reduce() (core int, essentials int) {
	alive := make([]bool, len(t.rows))
	for i := range alive {
		alive[i] = true
	}
	cols := uint32(1)<<t.cols - 1
	for changed := true; changed; {
		changed = false
		// Essential rows.
		for j := 0; j < t.cols; j++ {
			if cols>>j&1 == 0 {
				continue
			}
			only := -1
			cnt := 0
			for i, r := range t.rows {
				if alive[i] && r>>j&1 == 1 {
					cnt++
					only = i
				}
			}
			if cnt == 1 {
				cols &^= t.rows[only]
				alive[only] = false
				essentials++
				changed = true
			}
		}
		// Row dominance: drop a row whose remaining columns are a subset of another's.
		for q := range t.rows {
			if !alive[q] {
				continue
			}
			for p := range t.rows {
				if p == q || !alive[p] {
					continue
				}
				rq, rp := t.rows[q]&cols, t.rows[p]&cols
				if subset(rq, rp) && (rq != rp || p < q) {
					alive[q] = false
					changed = true
					break
				}
			}
		}
		// Column dominance: drop a column whose row set contains another column's row set.
		rowsOf := func(j int) uint32 {
			var s uint32
			for i, r := range t.rows {
				if alive[i] && r>>j&1 == 1 {
					s |= 1 << i
				}
			}
			return s
		}
		for i := 0; i < t.cols; i++ {
			if cols>>i&1 == 0 {
				continue
			}
			for j := 0; j < t.cols; j++ {
				if j == i || cols>>j&1 == 0 {
					continue
				}
				ri, rj := rowsOf(i), rowsOf(j)
				if subset(rj, ri) && (ri != rj || j < i) {
					cols &^= 1 << i
					changed = true
					break
				}
			}
		}
	}
	return popcountU(cols), essentials
}

// tableReduction: essentials, column dominance count, cyclic core size.
func tableReduction(ctx generator.GeneratorContext) generator.Problem {
	nr := edakit.Scaled(ctx, 4, 7)
	nc := edakit.Scaled(ctx, 5, 8)
	t := randCoverTable(nr, nc)
	intro := "A unit-cost covering table has columns c1 to c" + edakit.IntStr(nc) + ": " + t.text() + "."
	switch kp(ctx) {
	case 1:
		var ess uint32
		for j := 0; j < nc; j++ {
			cnt, only := 0, 0
			for i, r := range t.rows {
				if r>>j&1 == 1 {
					cnt++
					only = i
				}
			}
			if cnt == 1 {
				ess |= 1 << only
			}
		}
		e := popcountU(ess)
		return edakit.Numeric(intro+" How many rows are essential in this table, before any reduction?",
			edakit.IntStr(e), edakit.Sprintf("A row is essential when it is the only row marking some column. %d rows qualify.", e),
			map[string]string{"essential_rows": edakit.IntStr(e)})
	case 2:
		rowsOf := func(j int) uint32 {
			var s uint32
			for i, r := range t.rows {
				if r>>j&1 == 1 {
					s |= 1 << i
				}
			}
			return s
		}
		del := 0
		for i := 0; i < nc; i++ {
			for j := 0; j < nc; j++ {
				if j != i && subset(rowsOf(j), rowsOf(i)) && (rowsOf(i) != rowsOf(j) || j < i) {
					del++
					break
				}
			}
		}
		return edakit.Numeric(intro+" Apply column dominance once to the original table: a column is deleted when its set of rows contains another column's set of rows (of two identical columns, keep one). How many columns are deleted?",
			edakit.IntStr(del), edakit.Sprintf("Covering the column with the smaller row set forces the larger one to be covered too, so the larger column is redundant. %d columns go.", del),
			map[string]string{"deleted_columns": edakit.IntStr(del)})
	}
	core, ess := t.reduce()
	return edakit.Numeric(intro+" Repeat essential-row selection, row dominance (unit cost, ties keep the lower index) and column dominance until nothing changes. How many columns remain in the cyclic core? Answer 0 if the reductions solve the table.",
		edakit.IntStr(core), edakit.Sprintf("The reductions select %d essential rows in total and leave %d columns. A nonzero core needs branch and bound or Petrick's method.", ess, core),
		map[string]string{"core_columns": edakit.IntStr(core)})
}

// unateCovering: uncovered columns of a candidate, minimum rows, minimum weighted cost.
func unateCovering(ctx generator.GeneratorContext) generator.Problem {
	nr := edakit.Scaled(ctx, 4, 7)
	nc := edakit.Scaled(ctx, 5, 8)
	t := randCoverTable(nr, nc)
	intro := "A covering matrix has columns c1 to c" + edakit.IntStr(nc) + ": " + t.text() + "."
	switch kp(ctx) {
	case 1:
		var pick []int
		var cov uint32
		var names []string
		for i := range t.rows {
			if edakit.Int(0, 1) == 1 {
				pick = append(pick, i)
				cov |= t.rows[i]
				names = append(names, edakit.Sprintf("P%d", i+1))
			}
		}
		if len(pick) == 0 {
			pick = []int{0}
			cov = t.rows[0]
			names = []string{"P1"}
		}
		left := nc - popcountU(cov)
		return edakit.Numeric(intro+" Set the binary variables of "+strings.Join(names, ", ")+" to 1 and the rest to 0. How many covering constraints (one per column) are violated?",
			edakit.IntStr(left), edakit.Sprintf("Each column needs at least one chosen row. The selection leaves %d columns unmarked.", left),
			map[string]string{"violated": edakit.IntStr(left)})
	case 2:
		m := t.minRows(nil)
		return edakit.Numeric(intro+" Solve the unate covering problem: what is the minimum number of rows that cover every column?",
			edakit.IntStr(m), edakit.Sprintf("Minimize the sum of the row variables subject to one at-least-one constraint per column. The optimum uses %d rows.", m),
			map[string]string{"min_rows": edakit.IntStr(m)})
	}
	w := make([]int, nr)
	var ws []string
	for i := range w {
		w[i] = edakit.Int(1, 4)
		ws = append(ws, edakit.Sprintf("P%d costs %d", i+1, w[i]))
	}
	m := t.minRows(w)
	return edakit.Numeric(intro+" Row costs (literal counts): "+strings.Join(ws, ", ")+". What is the minimum total cost of a cover?",
		edakit.IntStr(m), edakit.Sprintf("Weighted covering minimizes the sum of cost times row variable. The cheapest cover costs %d.", m),
		map[string]string{"min_cost": edakit.IntStr(m)})
}

// petrick: clause absorption, product terms after expansion, number of minimum covers.
func petrick(ctx generator.GeneratorContext) generator.Problem {
	nr := edakit.Scaled(ctx, 3, 5)
	nc := edakit.Scaled(ctx, 3, 5)
	var t coverTable
	for {
		t = randCoverTable(nr, nc)
		ok := true
		for j := 0; j < nc && ok; j++ {
			c := 0
			for _, r := range t.rows {
				if r>>j&1 == 1 {
					c++
				}
			}
			if c < 2 {
				ok = false
			}
		}
		if ok {
			break
		}
	}
	// Clauses: the row set of each column.
	clauses := make([]uint32, nc)
	for j := 0; j < nc; j++ {
		for i, r := range t.rows {
			if r>>j&1 == 1 {
				clauses[j] |= 1 << i
			}
		}
	}
	intro := "A prime implicant table (every column covered by at least two rows) has columns c1 to c" + edakit.IntStr(nc) + ": " + t.text() + "."
	k := kp(ctx)
	switch k {
	case 1:
		kept := absorbMin(clauses)
		return edakit.Numeric(intro+" Write Petrick's product of sums, one clause per column. After removing duplicate clauses and any clause that contains another clause, how many clauses remain?",
			edakit.IntStr(len(kept)), edakit.Sprintf("A clause that contains another is implied by it, so it can be dropped. %d clauses remain.", len(kept)),
			map[string]string{"clauses": edakit.IntStr(len(kept))})
	}
	// Minimal transversals = products after multiplying out with absorption.
	var terms []uint32
	for s := uint32(1); s < 1<<nr; s++ {
		hit := true
		for _, c := range clauses {
			if c&s == 0 {
				hit = false
				break
			}
		}
		if hit {
			terms = append(terms, s)
		}
	}
	minimal := absorbMin(terms)
	if k == 2 {
		return edakit.Numeric(intro+" Multiply out Petrick's product of sums and apply absorption ($X + XY = X$). How many product terms are in the resulting irredundant sum of products?",
			edakit.IntStr(len(minimal)), edakit.Sprintf("Each surviving product is an irredundant cover. There are %d of them.", len(minimal)),
			map[string]string{"products": edakit.IntStr(len(minimal))})
	}
	best := 1 << 30
	for _, s := range minimal {
		if popcountU(s) < best {
			best = popcountU(s)
		}
	}
	cnt := 0
	for _, s := range minimal {
		if popcountU(s) == best {
			cnt++
		}
	}
	return edakit.Numeric(intro+" Use Petrick's method with unit row cost. How many different minimum covers (fewest rows) are there?",
		edakit.IntStr(cnt), edakit.Sprintf("After expansion and absorption, the smallest products have %d rows, and %d products reach that size.", best, cnt),
		map[string]string{"min_covers": edakit.IntStr(cnt)})
}

// absorbMin keeps the inclusion-minimal distinct sets.
func absorbMin(sets []uint32) []uint32 {
	var out []uint32
	for i, s := range sets {
		keep := true
		for j, u := range sets {
			if j != i && subset(u, s) && (u != s || j < i) {
				keep = false
				break
			}
		}
		if keep {
			out = append(out, s)
		}
	}
	return out
}

// tautology: cofactor size, uncovered minterms, unate variables.
func tautology(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 4)
	cs := randCover(n)
	list := make([]string, len(cs))
	for i, c := range cs {
		list[i] = c.pos(n)
	}
	intro := edakit.Sprintf("A cover over variables %s has the cubes %s.", strings.Join(varNames[:n], ", "), strings.Join(list, ", "))
	switch kp(ctx) {
	case 1:
		v := edakit.Int(0, n-1)
		b := edakit.Int(0, 1)
		p := uint32(1) << (n - 1 - v)
		cnt := 0
		for _, c := range cs {
			if c.mask&p == 0 || (c.val&p != 0) == (b == 1) {
				cnt++
			}
		}
		return edakit.Numeric(intro+edakit.Sprintf(" Take the cofactor with respect to %s = %d. How many cubes does the cofactor cover have?", varNames[v], b),
			edakit.IntStr(cnt), edakit.Sprintf("Cubes with the opposite literal of %s vanish; the others survive with %s raised to a dash. %d cubes remain.", varNames[v], varNames[v], cnt),
			map[string]string{"cofactor_cubes": edakit.IntStr(cnt)})
	case 2:
		var t uint64
		for _, c := range cs {
			t |= c.table(n)
		}
		miss := (1 << (1 << n)) - 1 - int(t)
		missing := popcount(uint64(miss))
		return edakit.Numeric(intro+" Decide whether the cover is a tautology by recursive cofactoring. How many of the minterms are not covered? Answer 0 if the cover is a tautology.",
			edakit.IntStr(missing), edakit.Sprintf("Splitting on variables until every leaf cover contains the all-dash cube or is empty shows %d uncovered minterms.", missing),
			map[string]string{"uncovered": edakit.IntStr(missing)})
	}
	un := 0
	for i := 0; i < n; i++ {
		p := uint32(1) << (n - 1 - i)
		pos, neg := false, false
		for _, c := range cs {
			if c.mask&p != 0 {
				if c.val&p != 0 {
					pos = true
				} else {
					neg = true
				}
			}
		}
		if !(pos && neg) {
			un++
		}
	}
	return edakit.Numeric(intro+" In how many variables is the cover unate (the variable appears in only one polarity, or not at all)?",
		edakit.IntStr(un), edakit.Sprintf("%d of the %d variables never appear in both polarities. A cover unate in every variable is a tautology exactly when it contains the all-dash cube.", un, n),
		map[string]string{"unate_vars": edakit.IntStr(un)})
}

// randCover builds a cover that is a tautology about half the time: a random Shannon split of
// the universe, optionally missing a leaf, plus a couple of extra overlapping cubes.
func randCover(n int) []cube {
	var leaves []cube
	var split func(c cube, depth int)
	split = func(c cube, depth int) {
		if c.lits() == n || (depth > 0 && edakit.Int(0, 2) == 0) {
			leaves = append(leaves, c)
			return
		}
		var free []int
		for i := 0; i < n; i++ {
			if c.mask&(1<<i) == 0 {
				free = append(free, i)
			}
		}
		p := uint32(1) << edakit.Pick(free)
		split(cube{c.mask | p, c.val}, depth+1)
		split(cube{c.mask | p, c.val | p}, depth+1)
	}
	for {
		leaves = nil
		split(cube{}, 0)
		if len(leaves) >= 4 && len(leaves) <= 7 {
			break
		}
	}
	if edakit.Int(0, 1) == 0 {
		i := edakit.Int(0, len(leaves)-1)
		leaves = append(leaves[:i], leaves[i+1:]...)
	}
	extra := edakit.Int(0, 2)
	for i := 0; i < extra; i++ {
		leaves = append(leaves, randCube(n, 2))
	}
	return edakit.Shuffle(leaves)
}

// expand: individually droppable literals, greedy expansion, minterms covered after expand.
func expand(ctx generator.GeneratorContext) generator.Problem {
	n := 4
	on := randFunc(n, 6, 11)
	start := cube{0xF, uint32(edakit.Pick(mintermList(on, n)))}
	if edakit.Int(0, 2) == 0 {
		// Start from a one-dash cube when one exists inside ON.
		p := uint32(1) << edakit.Int(0, 3)
		c := cube{start.mask &^ p, start.val &^ p}
		if c.table(n)&^on == 0 {
			start = c
		}
	}
	off := ^on & 0xFFFF
	intro := edakit.Sprintf("Function f of a, b, c, d has OFF-set minterms %s (everything else is ON). Espresso expands the cube %s, written in positional notation.",
		intsCSV(mintermList(off, n)), start.pos(n))
	hits := func(c cube) bool { return c.table(n)&off != 0 }
	k := kp(ctx)
	switch k {
	case 1:
		cnt := 0
		for i := 0; i < n; i++ {
			p := uint32(1) << i
			if start.mask&p != 0 && !hits(cube{start.mask &^ p, start.val &^ p}) {
				cnt++
			}
		}
		return edakit.Numeric(intro+" How many of its literals can each be raised to a dash on its own without the cube meeting the OFF-set?",
			edakit.IntStr(cnt), edakit.Sprintf("Test each literal separately: raising it doubles the cube, and the move is legal only if no new minterm is an OFF minterm. %d literals pass.", cnt),
			map[string]string{"droppable": edakit.IntStr(cnt)})
	}
	c := start
	for i := 0; i < n; i++ {
		p := uint32(1) << (n - 1 - i)
		if c.mask&p == 0 {
			continue
		}
		t := cube{c.mask &^ p, c.val &^ p}
		if !hits(t) {
			c = t
		}
	}
	if k == 2 {
		return edakit.Numeric(intro+" Raise literals greedily in the order a, b, c, d, keeping each raise that stays off the OFF-set. How many literals does the final cube have?",
			edakit.IntStr(c.lits()), edakit.Sprintf("The greedy pass ends at %s, a prime with %d literals. Another order can reach a different prime.", c.pos(n), c.lits()),
			map[string]string{"expanded": c.pos(n)})
	}
	cov := popcount(c.table(n))
	return edakit.Numeric(intro+" Raise literals greedily in the order a, b, c, d, keeping each legal raise. How many minterms does the expanded cube cover?",
		edakit.IntStr(cov), edakit.Sprintf("The expansion stops at %s, which covers %d minterms; every cube of the cover it now contains can be deleted.", c.pos(n), cov),
		map[string]string{"expanded": c.pos(n)})
}

// irredundant: relatively essential cubes, redundant cubes, minimum subcover.
func irredundant(ctx generator.GeneratorContext) generator.Problem {
	n := 4
	var on uint64
	var ps []cube
	for {
		on = randFunc(n, 5, 10)
		ps = primes(n, on, 0)
		if len(ps) >= 3 && len(ps) <= 8 {
			break
		}
	}
	list := make([]string, len(ps))
	for i, c := range ps {
		list[i] = c.pos(n)
	}
	intro := "A cover of f(a, b, c, d) holds the prime cubes " + strings.Join(list, ", ") + ", and f is exactly their union."
	relEss := 0
	for i, p := range ps {
		mine := p.table(n)
		for j, q := range ps {
			if j != i {
				mine &^= q.table(n)
			}
		}
		if mine != 0 {
			relEss++
		}
	}
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(intro+" How many cubes are relatively essential, that is, cover some minterm that no other cube of the cover covers?",
			edakit.IntStr(relEss), edakit.Sprintf("%d cubes each own a minterm and must stay in any subcover.", relEss),
			map[string]string{"rel_essential": edakit.IntStr(relEss)})
	case 2:
		red := len(ps) - relEss
		return edakit.Numeric(intro+" How many cubes are redundant, that is, contained in the union of the other cubes?",
			edakit.IntStr(red), edakit.Sprintf("%d of %d cubes are not relatively essential, so each is covered by the rest. They cannot all be removed at once in general.", red, len(ps)),
			map[string]string{"redundant": edakit.IntStr(red)})
	}
	m := minCover(n, on, ps)
	return edakit.Numeric(intro+" IRREDUNDANT keeps the relatively essential cubes and solves a covering problem over the rest. What is the minimum number of cubes in an irredundant subcover?",
		edakit.IntStr(len(m)), edakit.Sprintf("One minimum subcover is $%s$, with %d cubes.", sopTex(n, m), len(m)),
		map[string]string{"min_subcover": edakit.IntStr(len(m))})
}
