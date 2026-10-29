package fv

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// ---------- CNF helpers ----------

type clause []int // literals: +v is x_v, -v is ¬x_v

func litText(l int) string {
	if l < 0 {
		return edakit.Sprintf("¬x%d", -l)
	}
	return edakit.Sprintf("x%d", l)
}

func clauseText(c clause) string {
	s := append(clause(nil), c...)
	sort.Slice(s, func(i, j int) bool {
		a, b := s[i], s[j]
		if a < 0 {
			a = -a
		}
		if b < 0 {
			b = -b
		}
		return a < b
	})
	var parts []string
	for _, l := range s {
		parts = append(parts, litText(l))
	}
	return "(" + strings.Join(parts, " ∨ ") + ")"
}

func cnfText(cs []clause) string {
	var parts []string
	for i, c := range cs {
		parts = append(parts, edakit.Sprintf("C%d = %s", i+1, clauseText(c)))
	}
	return strings.Join(parts, ", ")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// randomClause draws a clause of size k over vars 1..n with distinct variables.
func randomClause(n, k int) clause {
	vs := rand.Perm(n)[:k]
	c := make(clause, k)
	for i, v := range vs {
		if rand.Intn(2) == 0 {
			c[i] = v + 1
		} else {
			c[i] = -(v + 1)
		}
	}
	return c
}

// addUnique appends c unless the same clause (as a set of literals) is already present.
func addUnique(cs []clause, c clause) []clause {
	for _, d := range cs {
		if len(d) == len(c) && subset(c, d) {
			return cs
		}
	}
	return append(cs, c)
}

// litValue returns 1 (true), -1 (false) or 0 (unassigned) for literal l under assign.
func litValue(assign []int, l int) int {
	v := assign[abs(l)]
	if v == 0 {
		return 0
	}
	if (l > 0) == (v > 0) {
		return 1
	}
	return -1
}

// bcp runs unit propagation to a fixpoint. It returns false on a conflict.
func bcp(cs []clause, assign []int) bool {
	for changed := true; changed; {
		changed = false
		for _, c := range cs {
			sat, free, last := false, 0, 0
			for _, l := range c {
				switch litValue(assign, l) {
				case 1:
					sat = true
				case 0:
					free++
					last = l
				}
			}
			if sat {
				continue
			}
			if free == 0 {
				return false
			}
			if free == 1 {
				if last > 0 {
					assign[last] = 1
				} else {
					assign[-last] = -1
				}
				changed = true
			}
		}
	}
	return true
}

func assignedCount(assign []int) int {
	n := 0
	for _, v := range assign[1:] {
		if v != 0 {
			n++
		}
	}
	return n
}

func allSatisfied(cs []clause, assign []int) bool {
	for _, c := range cs {
		ok := false
		for _, l := range c {
			if litValue(assign, l) == 1 {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func bruteSat(cs []clause, n int) bool {
	for m := 0; m < 1<<n; m++ {
		assign := make([]int, n+1)
		for v := 1; v <= n; v++ {
			if m>>(v-1)&1 == 1 {
				assign[v] = 1
			} else {
				assign[v] = -1
			}
		}
		if allSatisfied(cs, assign) {
			return true
		}
	}
	return false
}

// ---------- Tseitin ----------

func tseitin(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		kind := edakit.Pick([]string{"AND", "OR", "NAND", "NOR", "XOR", "XNOR", "NOT"})
		switch kind {
		case "XOR", "XNOR":
			f := "$y = a \\oplus b$"
			if kind == "XNOR" {
				f = "$y = \\neg(a \\oplus b)$"
			}
			q := edakit.Sprintf("Tseitin-encode a single 2-input %s gate %s, introducing one variable for its output. How many clauses define the gate?", kind, f)
			return edakit.Numeric(q, "4", "A 2-input XOR or XNOR is false or true on exactly half of its four input rows, so its equivalence $y \\leftrightarrow (a \\oplus b)$ needs one 3-literal clause per row: 4 clauses.", map[string]string{"clauses": "4"})
		case "NOT":
			q := "Under the convention that inverters are absorbed into literals, how many clauses does the Tseitin encoding spend on a NOT gate?"
			return edakit.Numeric(q, "0", "A NOT gate gets no variable: every reader of its output simply uses the negated literal of its input, so it costs 0 clauses.", nil)
		}
		k := edakit.Scaled(ctx, 2, 6)
		ans := k + 1
		q := edakit.Sprintf("Tseitin-encode a single %d-input %s gate, introducing one variable $y$ for its output. How many clauses define $y$?", k, kind)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("A %d-input %s needs one binary clause per input (each input forces or is forced by the output) plus one long clause of %d literals: %d + 1 = %d clauses. Inverting the output only flips the polarity of $y$.", k, kind, k+1, k, ans),
			map[string]string{"clauses": edakit.IntStr(ans)})
	case 2:
		inputs := edakit.Scaled(ctx, 3, 8)
		and2 := edakit.Int(1, 4)
		and3 := edakit.Int(0, 2)
		or2 := edakit.Int(0, 3)
		xor := edakit.Int(1, 3)
		not := edakit.Int(0, 3)
		gates := and2 + and3 + or2 + xor
		vars := inputs + gates
		cls := and2*3 + and3*4 + or2*3 + xor*4 + 1
		desc := edakit.Sprintf("%d primary inputs, %d two-input AND gates, %d three-input AND gates, %d two-input OR gates, %d two-input XOR gates and %d NOT gates, with one output that is asserted to be 1", inputs, and2, and3, or2, xor, not)
		if rand.Intn(2) == 0 {
			q := "A circuit has " + desc + ". Using the Tseitin encoding with one variable per primary input and per gate, and inverters absorbed into literals, how many CNF variables are there?"
			return edakit.Numeric(q, edakit.IntStr(vars),
				edakit.Sprintf("Variables: %d inputs + %d non-inverter gates = %d. NOT gates add none.", inputs, gates, vars),
				map[string]string{"vars": edakit.IntStr(vars)})
		}
		q := "A circuit has " + desc + ". Using the Tseitin encoding (k-input AND or OR: k + 1 clauses, 2-input XOR: 4 clauses, inverters absorbed, one unit clause for the asserted output), how many clauses are there?"
		return edakit.Numeric(q, edakit.IntStr(cls),
			edakit.Sprintf("Clauses: %d·3 + %d·4 + %d·3 + %d·4 + 1 = %d.", and2, and3, or2, xor, cls),
			map[string]string{"clauses": edakit.IntStr(cls)})
	}
	n := edakit.Scaled(ctx, 3, 9)
	direct := 1 << (n - 1)
	ts := 4*(n-1) + 1
	switch rand.Intn(3) {
	case 0:
		q := edakit.Sprintf("The parity $x_1 \\oplus x_2 \\oplus \\dots \\oplus x_{%d}$ is asserted to be 1 and converted to CNF directly over the inputs, with no new variables. How many clauses does the smallest such CNF have?", n)
		return edakit.Numeric(q, edakit.IntStr(direct),
			edakit.Sprintf("Every clause over the original variables must contain all %d variables, so each clause excludes one even-parity row. There are $2^{%d} = %d$ of them.", n, n-1, direct),
			map[string]string{"direct": edakit.IntStr(direct)})
	case 1:
		q := edakit.Sprintf("The parity of %d inputs is built as a chain of %d two-input XOR gates, Tseitin-encoded (4 clauses per XOR), and the final output is asserted to be 1 with a unit clause. How many clauses are there?", n, n-1)
		return edakit.Numeric(q, edakit.IntStr(ts),
			edakit.Sprintf("%d XORs at 4 clauses each, plus 1 unit clause: 4·%d + 1 = %d.", n-1, n-1, ts),
			map[string]string{"tseitin": edakit.IntStr(ts)})
	}
	q := edakit.Sprintf("Asserting the parity of %d inputs: how many more clauses does the direct CNF (no new variables) have than a Tseitin chain of 2-input XORs at 4 clauses each plus one unit clause?", n)
	return edakit.Numeric(q, edakit.IntStr(direct-ts),
		edakit.Sprintf("Direct: $2^{%d} = %d$ clauses. Tseitin: 4·%d + 1 = %d. Difference: %d.", n-1, direct, n-1, ts, direct-ts),
		map[string]string{"direct": edakit.IntStr(direct), "tseitin": edakit.IntStr(ts)})
}

// ---------- unit propagation ----------

func unitProp(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	n := edakit.Scaled(ctx, 5, 7)
	switch kp {
	case 1:
		for {
			m := edakit.Int(5, 7)
			var cs []clause
			for i := 0; i < m; i++ {
				cs = addUnique(cs, randomClause(n, edakit.Int(2, 3)))
			}
			assign := make([]int, n+1)
			fixed := randomSubset(n, edakit.Int(2, 3))
			var at []string
			for _, v := range fixed {
				val := 1
				if rand.Intn(2) == 0 {
					val = -1
				}
				assign[v+1] = val
				at = append(at, edakit.Sprintf("x%d = %d", v+1, (val+1)/2))
			}
			units := 0
			var which []string
			for i, c := range cs {
				sat, free := false, 0
				for _, l := range c {
					switch litValue(assign, l) {
					case 1:
						sat = true
					case 0:
						free++
					}
				}
				if !sat && free == 1 {
					units++
					which = append(which, edakit.Sprintf("C%d", i+1))
				}
			}
			if units == 0 && rand.Intn(3) > 0 {
				continue
			}
			q := "Clauses: " + cnfText(cs) + ". Under the partial assignment " + strings.Join(at, ", ") +
				", how many clauses are unit (not satisfied, with exactly one unassigned literal)?"
			why := "No clause is unit."
			if units > 0 {
				why = "The unit clauses are " + strings.Join(which, ", ") + ": each has every literal false except one unassigned literal."
			}
			return edakit.Numeric(q, edakit.IntStr(units), why, map[string]string{"units": edakit.IntStr(units)})
		}
	case 2:
		for {
			m := edakit.Int(5, 8)
			cs := []clause{randomClause(n, 1)}
			if rand.Intn(2) == 0 {
				cs = addUnique(cs, randomClause(n, 1))
			}
			for i := 0; i < m; i++ {
				cs = addUnique(cs, randomClause(n, edakit.Int(2, 3)))
			}
			assign := make([]int, n+1)
			if !bcp(cs, assign) {
				continue
			}
			k := assignedCount(assign)
			if k < 3 {
				continue
			}
			cs = edakit.Shuffle(cs)
			q := edakit.Sprintf("Run unit propagation to a fixpoint on the CNF over $x_1, \\dots, x_{%d}$: %s. No conflict arises. How many variables are assigned at the fixpoint, counting those forced by the original unit clauses?", n, cnfText(cs))
			return edakit.Numeric(q, edakit.IntStr(k),
				edakit.Sprintf("Start from the unit clauses, then repeatedly find a clause with every literal false except one and set that literal true. The fixpoint assigns %d variables: %s.", k, assignText(assign)),
				map[string]string{"assigned": edakit.IntStr(k)})
		}
	}
	for {
		m := edakit.Int(6, 8)
		var cs []clause
		for i := 0; i < m; i++ {
			cs = addUnique(cs, randomClause(n, edakit.Int(2, 3)))
		}
		v := edakit.Int(1, n)
		val := edakit.Pick([]int{1, -1})
		assign := make([]int, n+1)
		assign[v] = val
		if !bcp(cs, assign) {
			continue
		}
		k := assignedCount(assign) - 1
		if k < 2 {
			continue
		}
		q := edakit.Sprintf("A solver decides $x_%d = %d$ on the CNF %s. How many further variables does unit propagation imply (do not count the decision)? No conflict arises.", v, (val+1)/2, cnfText(cs))
		return edakit.Numeric(q, edakit.IntStr(k),
			edakit.Sprintf("After the decision, propagation reaches the fixpoint %s, which is %d implied variables beyond the decision.", assignText(assign), k),
			map[string]string{"implied": edakit.IntStr(k)})
	}
}

func assignText(assign []int) string {
	var parts []string
	for v := 1; v < len(assign); v++ {
		if assign[v] != 0 {
			parts = append(parts, edakit.Sprintf("x%d = %d", v, (assign[v]+1)/2))
		}
	}
	return strings.Join(parts, ", ")
}

// ---------- DPLL ----------

// dpllCount runs DPLL with BCP, branching on the lowest unassigned variable with value 0
// first, and counts branching assignments.
func dpllCount(cs []clause, n int) (bool, int) {
	decisions := 0
	var solve func(assign []int) bool
	solve = func(assign []int) bool {
		if !bcp(cs, assign) {
			return false
		}
		if allSatisfied(cs, assign) {
			return true
		}
		v := 0
		for i := 1; i <= n; i++ {
			if assign[i] == 0 {
				v = i
				break
			}
		}
		if v == 0 {
			return false
		}
		for _, val := range []int{-1, 1} {
			decisions++
			a := append([]int(nil), assign...)
			a[v] = val
			if solve(a) {
				return true
			}
		}
		return false
	}
	ok := solve(make([]int, n+1))
	return ok, decisions
}

func dpll(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		var n, m int
		wantSat := true
		switch kp {
		case 1:
			n, m = 3, edakit.Int(4, 6)
		case 2:
			n, m = edakit.Scaled(ctx, 4, 5), edakit.Int(6, 9)
		default:
			n, m = 4, edakit.Int(12, 16)
			wantSat = false
		}
		var cs []clause
		for i := 0; i < m; i++ {
			k := 3
			if n == 3 || rand.Intn(3) == 0 {
				k = 2
			}
			cs = addUnique(cs, randomClause(n, k))
		}
		if bruteSat(cs, n) != wantSat {
			continue
		}
		ok, d := dpllCount(cs, n)
		if d < 2 || d > 10 {
			continue
		}
		res := "SAT"
		if !ok {
			res = "UNSAT"
		}
		q := edakit.Sprintf("Run DPLL on %s. After every assignment it runs unit propagation; it reports SAT as soon as every clause is satisfied; it branches on the lowest-index unassigned variable, tries 0 before 1, and backtracks chronologically. Count every branching assignment as a decision, including a second value tried after backtracking. How many decisions does it make before it finishes?", cnfText(cs))
		return edakit.Numeric(q, edakit.IntStr(d),
			edakit.Sprintf("Tracing the search tree with propagation at every node gives %d decisions, and the solver answers %s.", d, res),
			map[string]string{"decisions": edakit.IntStr(d), "result": res})
	}
}

// ---------- implication graphs ----------

type tnode struct {
	v        int
	val      bool
	level    int
	decision bool
	ante     []int // indices of antecedent nodes
}

type trail struct {
	nodes    []tnode
	conflict []int // antecedents of the conflict
	levels   int
}

// lit returns the literal that node i makes true.
func (t *trail) lit(i int) int {
	if t.nodes[i].val {
		return t.nodes[i].v
	}
	return -t.nodes[i].v
}

// makeTrail draws a synthetic trail with the given number of decision levels.
func makeTrail(levels, lastImplied int) *trail {
	t := &trail{levels: levels}
	vars := rand.Perm(20)
	next := 0
	newNode := func(level int, dec bool, ante []int) {
		t.nodes = append(t.nodes, tnode{v: vars[next] + 1, val: rand.Intn(2) == 0, level: level, decision: dec, ante: ante})
		next++
	}
	for l := 1; l <= levels; l++ {
		start := len(t.nodes)
		newNode(l, true, nil)
		k := edakit.Int(0, 2)
		if l == levels {
			k = lastImplied
		}
		for i := 0; i < k; i++ {
			cur := start + rand.Intn(len(t.nodes)-start)
			ante := []int{cur}
			extra := edakit.Int(0, 2)
			for j := 0; j < extra; j++ {
				a := rand.Intn(len(t.nodes))
				dup := false
				for _, x := range ante {
					if x == a {
						dup = true
					}
				}
				if !dup {
					ante = append(ante, a)
				}
			}
			sort.Ints(ante)
			newNode(l, false, ante)
		}
	}
	// Conflict: two or three current-level nodes plus up to two earlier ones.
	start := 0
	for i, nd := range t.nodes {
		if nd.level == levels {
			start = i
			break
		}
	}
	cur := rand.Perm(len(t.nodes) - start)
	k := minInt(len(cur), edakit.Int(2, 3))
	for i := 0; i < k; i++ {
		t.conflict = append(t.conflict, start+cur[i])
	}
	if start > 0 {
		extra := edakit.Int(0, 2)
		for j := 0; j < extra; j++ {
			a := rand.Intn(start)
			dup := false
			for _, x := range t.conflict {
				if x == a {
					dup = true
				}
			}
			if !dup {
				t.conflict = append(t.conflict, a)
			}
		}
	}
	sort.Ints(t.conflict)
	return t
}

// reasonText is the reason clause of implied node i.
func (t *trail) reasonText(i int) string {
	c := clause{t.lit(i)}
	for _, a := range t.nodes[i].ante {
		c = append(c, -t.lit(a))
	}
	return clauseText(c)
}

// nodesClause renders a learned clause given as a set of falsified nodes.
func (t *trail) nodesClause(set []int) string {
	var c clause
	for _, i := range set {
		c = append(c, -t.lit(i))
	}
	return clauseText(c)
}

func (t *trail) trailText() string {
	var lines []string
	for l := 1; l <= t.levels; l++ {
		var parts []string
		for i, nd := range t.nodes {
			if nd.level != l {
				continue
			}
			if nd.decision {
				parts = append(parts, edakit.Sprintf("decide x%d = %d", nd.v, b2i(nd.val)))
			} else {
				parts = append(parts, edakit.Sprintf("x%d = %d by %s", nd.v, b2i(nd.val), t.reasonText(i)))
			}
		}
		lines = append(lines, edakit.Sprintf("level %d: %s", l, strings.Join(parts, "; ")))
	}
	return strings.Join(lines, ". ")
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// resolveUntil resolves the conflict clause on current-level implied nodes, latest first,
// until stop reports true. It returns the clause as a sorted node set and the number of steps.
func (t *trail) resolveUntil(stop func(set map[int]bool) bool, maxSteps int) ([]int, int) {
	set := map[int]bool{}
	for _, i := range t.conflict {
		set[i] = true
	}
	steps := 0
	for !stop(set) && steps < maxSteps {
		pick := -1
		for i := range set {
			nd := t.nodes[i]
			if nd.level == t.levels && !nd.decision && i > pick {
				pick = i
			}
		}
		if pick < 0 {
			break
		}
		delete(set, pick)
		for _, a := range t.nodes[pick].ante {
			set[a] = true
		}
		steps++
	}
	var out []int
	for i := range set {
		out = append(out, i)
	}
	sort.Ints(out)
	return out, steps
}

func (t *trail) currentCount(set map[int]bool) int {
	n := 0
	for i := range set {
		if t.nodes[i].level == t.levels {
			n++
		}
	}
	return n
}

func (t *trail) firstUIP() ([]int, int) {
	return t.resolveUntil(func(s map[int]bool) bool { return t.currentCount(s) <= 1 }, 100)
}

func (t *trail) lastUIP() []int {
	out, _ := t.resolveUntil(func(s map[int]bool) bool {
		for i := range s {
			if t.nodes[i].level == t.levels && !t.nodes[i].decision {
				return false
			}
		}
		return true
	}, 100)
	return out
}

// decisionCone returns the decision nodes in the backward cone of the conflict.
func (t *trail) decisionCone() []int {
	seen := map[int]bool{}
	var walk func(i int)
	walk = func(i int) {
		if seen[i] {
			return
		}
		seen[i] = true
		for _, a := range t.nodes[i].ante {
			walk(a)
		}
	}
	for _, i := range t.conflict {
		walk(i)
	}
	var out []int
	for i := range seen {
		if t.nodes[i].decision {
			out = append(out, i)
		}
	}
	sort.Ints(out)
	return out
}

func implGraph(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	levels := edakit.Scaled(ctx, 3, 4)
	t := makeTrail(levels, edakit.Int(2, 4))
	// Decisions with levels, implied literals listed in shuffled order.
	var dec []string
	var imp []int
	for i, nd := range t.nodes {
		if nd.decision {
			dec = append(dec, edakit.Sprintf("x%d = %d at level %d", nd.v, b2i(nd.val), nd.level))
		} else {
			imp = append(imp, i)
		}
	}
	imp = edakit.Shuffle(imp)
	var impl []string
	for _, i := range imp {
		impl = append(impl, edakit.Sprintf("x%d = %d by %s", t.nodes[i].v, b2i(t.nodes[i].val), t.reasonText(i)))
	}
	ctxText := "Decisions: " + strings.Join(dec, ", ") + ". Implied literals, listed in no particular order, each with its reason clause: " + strings.Join(impl, "; ") + "."
	switch kp {
	case 1:
		// pick an implied node, preferably not at level 1
		i := imp[rand.Intn(len(imp))]
		nd := t.nodes[i]
		q := ctxText + edakit.Sprintf(" At what decision level is x%d assigned?", nd.v)
		return edakit.Numeric(q, edakit.IntStr(nd.level),
			edakit.Sprintf("An implied literal sits at the highest decision level among the literals that forced it. Tracing the reasons back from x%d gives level %d.", nd.v, nd.level),
			map[string]string{"level": edakit.IntStr(nd.level)})
	case 2:
		i := imp[rand.Intn(len(imp))]
		nd := t.nodes[i]
		k := len(nd.ante)
		q := ctxText + edakit.Sprintf(" In the implication graph, how many incoming edges does the node for x%d have?", nd.v)
		return edakit.Numeric(q, edakit.IntStr(k),
			edakit.Sprintf("Each other literal of the reason clause %s is false and points at the node: %d incoming edges.", t.reasonText(i), k),
			map[string]string{"indegree": edakit.IntStr(k)})
	}
	cone := t.decisionCone()
	q := ctxText + " The clause " + t.nodesClause(t.conflict) + " is now falsified. How many decision literals lie in the backward cone of this conflict in the implication graph?"
	return edakit.Numeric(q, edakit.IntStr(len(cone)),
		edakit.Sprintf("Walking back from the conflict through reason clauses reaches the decisions %s, so the all-decision cut has %d literals.", t.nodesClause(cone), len(cone)),
		map[string]string{"decisions": edakit.IntStr(len(cone))})
}

func firstUIP(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for tries := 0; ; tries++ {
		levels, last := 2, 2
		switch kp {
		case 2:
			levels, last = 3, 3
		case 3:
			levels, last = edakit.Scaled(ctx, 3, 4), 4
		}
		t := makeTrail(levels, last)
		uip, steps := t.firstUIP()
		if steps == 0 {
			continue
		}
		if kp == 1 && steps > 2 {
			continue
		}
		ans := t.nodesClause(uip)
		dist := []string{t.nodesClause(t.conflict), t.nodesClause(t.decisionCone()), t.nodesClause(t.lastUIP())}
		if steps > 1 {
			mid, _ := t.resolveUntil(func(s map[int]bool) bool { return false }, 1)
			dist = append(dist, t.nodesClause(mid))
		}
		uniq := map[string]bool{ans: true}
		var ds []string
		for _, d := range dist {
			if !uniq[d] {
				uniq[d] = true
				ds = append(ds, d)
			}
		}
		if len(ds) < 3 && tries < 50 {
			continue
		}
		if len(ds) < 2 {
			continue
		}
		q := "A CDCL solver's trail, in assignment order: " + t.trailText() + ". Now the clause " + t.nodesClause(t.conflict) +
			edakit.Sprintf(" is falsified at level %d. Which clause does first-UIP conflict analysis learn?", levels)
		return edakit.Choice(q, ans,
			edakit.Sprintf("Resolve the conflict clause with the reason of the most recently assigned level-%d literal, and repeat while more than one level-%d literal remains. After %d resolution step(s) the clause %s has a single level-%d literal, the first UIP.", levels, levels, steps, ans, levels),
			ds...)
	}
}

// ---------- backjumping ----------

func backjump(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 1 {
		cur := edakit.Scaled(ctx, 4, 9)
		k := edakit.Int(1, 4)
		vars := rand.Perm(15)
		var parts []string
		second := 0
		for i := 0; i < k; i++ {
			lv := cur
			if i > 0 {
				lv = edakit.Int(1, cur-1)
				second = maxInt(second, lv)
			}
			l := vars[i] + 1
			if rand.Intn(2) == 0 {
				l = -l
			}
			parts = append(parts, edakit.Sprintf("%s (level %d)", litText(l), lv))
		}
		if rand.Intn(4) == 0 {
			parts = parts[:1]
			second = 0
		}
		parts = edakit.Shuffle(parts)
		q := edakit.Sprintf("A conflict at decision level %d produces the learned clause with literals %s. To which decision level does a CDCL solver backjump?", cur, strings.Join(parts, ", "))
		why := edakit.Sprintf("Backjump to the second-highest level in the learned clause, %d: there the clause becomes unit and asserts its level-%d literal.", second, cur)
		if second == 0 {
			why = "A learned unit clause holds at the root, so the solver backjumps to level 0."
		}
		return edakit.Numeric(q, edakit.IntStr(second), why, map[string]string{"backjump": edakit.IntStr(second)})
	}
	for {
		levels := edakit.Scaled(ctx, 3, 5)
		t := makeTrail(levels, edakit.Int(2, 3))
		uip, steps := t.firstUIP()
		if steps == 0 {
			continue
		}
		bj := 0
		for _, i := range uip {
			if l := t.nodes[i].level; l < levels && l > bj {
				bj = l
			}
		}
		if kp == 2 {
			q := "Trail: " + t.trailText() + ". The clause " + t.nodesClause(t.conflict) +
				edakit.Sprintf(" is falsified at level %d. After learning the first-UIP clause, to which decision level does the solver backjump?", levels)
			return edakit.Numeric(q, edakit.IntStr(bj),
				edakit.Sprintf("The first-UIP clause is %s. Its second-highest decision level is %d, so the solver backjumps there.", t.nodesClause(uip), bj),
				map[string]string{"learned": t.nodesClause(uip), "backjump": edakit.IntStr(bj)})
		}
		skip := levels - 1 - bj
		q := "Trail: " + t.trailText() + ". The clause " + t.nodesClause(t.conflict) +
			edakit.Sprintf(" is falsified at level %d. Compared with chronological backtracking, which would return to level %d, how many more decision levels does backjumping on the first-UIP clause undo?", levels, levels-1)
		return edakit.Numeric(q, edakit.IntStr(skip),
			edakit.Sprintf("The first-UIP clause %s backjumps to level %d. Chronological backtracking returns to level %d, so backjumping skips %d - %d = %d levels.", t.nodesClause(uip), bj, levels-1, levels-1, bj, skip),
			map[string]string{"backjump": edakit.IntStr(bj)})
	}
}

// ---------- watched literals ----------

func watched(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 2 {
		for {
			k := edakit.Scaled(ctx, 4, 6)
			c := randomClause(9, k)
			w := [2]int{0, 1}
			order := rand.Perm(k)
			ev := order[:edakit.Int(2, minInt(4, k-1))]
			falseSet := map[int]bool{}
			moves := 0
			hits := 0
			for _, e := range ev {
				falseSet[e] = true
				for wi := 0; wi < 2; wi++ {
					if w[wi] != e {
						continue
					}
					hits++
					other := w[1-wi]
					for p := 0; p < k; p++ {
						if p != other && p != e && !falseSet[p] {
							w[wi] = p
							moves++
							break
						}
					}
				}
			}
			if hits == 0 {
				continue
			}
			var evs []string
			for _, e := range ev {
				evs = append(evs, litText(c[e]))
			}
			var lits []string
			for _, l := range c {
				lits = append(lits, litText(l))
			}
			q := edakit.Sprintf("The clause (%s) has its literals stored in that order and initially watches its first two literals. These literals become false, in order: %s. Each time a watched literal becomes false, the solver scans the clause from the start for the first literal that is neither false nor the other watch and moves the watch there; if none exists, the clause is unit or conflicting and the watch stays. How many watch moves happen?",
				strings.Join(lits, " ∨ "), strings.Join(evs, ", "))
			return edakit.Numeric(q, edakit.IntStr(moves),
				edakit.Sprintf("Only falsifying a watched literal triggers work. Replaying the sequence, %d of the events hit a watch and %d of those found a replacement literal.", hits, moves),
				map[string]string{"moves": edakit.IntStr(moves)})
		}
	}
	for {
		n := 6
		m := edakit.Scaled(ctx, 6, 9)
		var rows []string
		var cs []clause
		var ws [][2]int
		for i := 0; i < m; i++ {
			c := randomClause(n, edakit.Int(3, 4))
			p := rand.Perm(len(c))
			cs = append(cs, c)
			ws = append(ws, [2]int{c[p[0]], c[p[1]]})
			rows = append(rows, edakit.Sprintf("C%d %s watches %s, %s", i+1, clauseText(c), litText(c[p[0]]), litText(c[p[1]])))
		}
		v := edakit.Int(1, n)
		val := edakit.Pick([]int{0, 1})
		f := v // literal made false
		if val == 1 {
			f = -v
		}
		visits, occ := 0, 0
		for i, c := range cs {
			if ws[i][0] == f || ws[i][1] == f {
				visits++
			}
			for _, l := range c {
				if l == f {
					occ++
				}
			}
		}
		if visits == 0 || (kp == 3 && occ == visits && rand.Intn(3) > 0) {
			continue
		}
		head := "Watch lists: " + strings.Join(rows, "; ") + edakit.Sprintf(". The solver assigns $x_%d = %d$.", v, val)
		if kp == 1 {
			q := head + " With two watched literals, how many clauses must it visit?"
			return edakit.Numeric(q, edakit.IntStr(visits),
				edakit.Sprintf("The assignment falsifies %s, and only clauses watching %s are visited: %d of them.", litText(f), litText(f), visits),
				map[string]string{"visits": edakit.IntStr(visits)})
		}
		q := head + " An older scheme visits every clause that contains the falsified literal. How many fewer clause visits does the two-watched-literal scheme make?"
		return edakit.Numeric(q, edakit.IntStr(occ-visits),
			edakit.Sprintf("%s occurs in %d clauses but is watched in only %d, so the watched scheme saves %d visits.", litText(f), occ, visits, occ-visits),
			map[string]string{"occurrences": edakit.IntStr(occ), "visits": edakit.IntStr(visits)})
	}
}

// ---------- VSIDS ----------

func vsids(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		n := 6
		k := edakit.Scaled(ctx, 3, 4)
		decay := edakit.Pick([]float64{0.5, 0.8})
		var conf [][]int
		var ctext []string
		for i := 0; i < k; i++ {
			vs := randomSubset(n, edakit.Int(2, 3))
			for j := range vs {
				vs[j]++
			}
			conf = append(conf, vs)
			ctext = append(ctext, "{"+varList(vs)+"}")
		}
		act := make([]float64, n+1)
		exp := make([]float64, n+1)
		inc := 1.0
		for _, vs := range conf {
			for v := 1; v <= n; v++ {
				exp[v] *= decay
			}
			for _, v := range vs {
				act[v] += inc
				exp[v]++
			}
			inc /= decay
		}
		seq := strings.Join(ctext, ", then ")
		switch kp {
		case 1, 3:
			v := conf[k-1][rand.Intn(len(conf[k-1]))]
			if rand.Intn(2) == 0 {
				v = conf[0][0]
			}
			if kp == 1 {
				q := edakit.Sprintf("VSIDS with decay factor %s: all activities start at 0 and the increment starts at 1. On each conflict every variable of the learned clause is bumped by the increment, and then the increment is divided by %s. The learned clauses contain the variables %s. What is the activity of $x_%d$? Round to 2 decimal places.",
					edakit.Num(decay, 2), edakit.Num(decay, 2), seq, v)
				return edakit.Numeric(q, edakit.Num(act[v], 2),
					edakit.Sprintf("The increments are 1, %s, %s, ... Summing the increments of the conflicts that contain $x_%d$ gives %s.", edakit.Num(1/decay, 4), edakit.Num(1/decay/decay, 4), v, edakit.Num(act[v], 2)),
					map[string]string{"activity": edakit.Num(act[v], 4)})
			}
			q := edakit.Sprintf("Explicit-decay VSIDS with factor %s: activities start at 0. On each conflict the solver first multiplies every activity by %s and then adds 1 to each variable of the learned clause. The learned clauses contain the variables %s. What is the activity of $x_%d$? Round to 2 decimal places.",
				edakit.Num(decay, 2), edakit.Num(decay, 2), seq, v)
			return edakit.Numeric(q, edakit.Num(exp[v], 2),
				edakit.Sprintf("Each bump of 1 is decayed once per later conflict: the activity of $x_%d$ is %s. It is the growing-increment activity scaled by %s, so the ranking is the same.", v, edakit.Num(exp[v], 2), edakit.Num(exp[v]/act[v], 4)),
				map[string]string{"activity": edakit.Num(exp[v], 4)})
		}
		assigned := map[int]bool{}
		best, bestV, tie := -1.0, 0, false
		for _, x := range randomSubset(n, 2) {
			assigned[x+1] = true
		}
		for v := 1; v <= n; v++ {
			if assigned[v] {
				continue
			}
			if act[v] > best+1e-9 {
				best, bestV, tie = act[v], v, false
			} else if abs64(act[v]-best) < 1e-9 {
				tie = true
			}
		}
		if tie || best <= 0 {
			continue
		}
		var as []int
		for v := range assigned {
			as = append(as, v)
		}
		sort.Ints(as)
		q := edakit.Sprintf("VSIDS with decay factor %s (increment starts at 1 and is divided by %s after each conflict; activities start at 0). The learned clauses contain the variables %s. Variables %s are currently assigned. Which variable index does the solver pick for its next decision?",
			edakit.Num(decay, 2), edakit.Num(decay, 2), seq, varList(as))
		return edakit.Numeric(q, edakit.IntStr(bestV),
			edakit.Sprintf("Among the unassigned variables, $x_%d$ has the highest activity, %s.", bestV, edakit.Num(best, 4)),
			map[string]string{"next": edakit.IntStr(bestV)})
	}
}

func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func varList(vs []int) string {
	var s []string
	for _, v := range vs {
		s = append(s, edakit.Sprintf("x%d", v))
	}
	return strings.Join(s, ", ")
}

// ---------- restarts ----------

func luby(i int) int {
	for k := 1; ; k++ {
		if i == (1<<k)-1 {
			return 1 << (k - 1)
		}
		if i >= 1<<(k-1) && i < (1<<k)-1 {
			return luby(i - (1 << (k - 1)) + 1)
		}
	}
}

func restarts(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		i := edakit.Scaled(ctx, 4, 31)
		ans := luby(i)
		q := edakit.Sprintf("The Luby sequence begins 1, 1, 2, 1, 1, 2, 4, ... What is its term number %d (counting from 1)?", i)
		var terms []string
		for j := 1; j <= i; j++ {
			terms = append(terms, edakit.IntStr(luby(j)))
		}
		return edakit.Numeric(q, edakit.IntStr(ans), "The first terms are "+strings.Join(terms, ", ")+".", map[string]string{"term": edakit.IntStr(ans)})
	case 2:
		u := edakit.Pick([]int{32, 50, 64, 100, 128, 256})
		r := edakit.Scaled(ctx, 3, 12)
		sum := 0
		var terms []string
		for j := 1; j <= r; j++ {
			sum += luby(j)
			terms = append(terms, edakit.IntStr(luby(j)))
		}
		q := edakit.Sprintf("A solver restarts on the Luby schedule with a unit of %d conflicts: run $j$ lasts %d times the $j$-th Luby term. How many conflicts has it performed when run %d ends?", u, u, r)
		return edakit.Numeric(q, edakit.IntStr(u*sum),
			edakit.Sprintf("The first %d Luby terms are %s, summing to %d. Times %d gives %d conflicts.", r, strings.Join(terms, ", "), sum, u, u*sum),
			map[string]string{"sum": edakit.IntStr(sum)})
	}
	c := edakit.Pick([]int{50, 100, 200})
	g := edakit.Pick([][2]int{{3, 2}, {2, 1}, {5, 4}, {6, 5}})
	r := edakit.Scaled(ctx, 3, 6)
	total := 0
	var runs []string
	num, den := 1, 1
	for j := 1; j <= r; j++ {
		run := c * num / den
		total += run
		runs = append(runs, edakit.IntStr(run))
		num *= g[0]
		den *= g[1]
	}
	gs := edakit.Num(float64(g[0])/float64(g[1]), 2)
	q := edakit.Sprintf("Geometric restarts: run $j$ lasts $\\lfloor %d \\cdot %s^{j-1} \\rfloor$ conflicts. How many conflicts are performed in the first %d runs together?", c, gs, r)
	return edakit.Numeric(q, edakit.IntStr(total),
		edakit.Sprintf("The run lengths are %s, which sum to %d.", strings.Join(runs, ", "), total),
		map[string]string{"total": edakit.IntStr(total)})
}

// ---------- preprocessing ----------

func resolve(a, b clause, v int) (clause, bool) {
	set := map[int]bool{}
	for _, l := range a {
		if abs(l) != v {
			set[l] = true
		}
	}
	for _, l := range b {
		if abs(l) != v {
			set[l] = true
		}
	}
	var out clause
	for l := range set {
		if set[-l] {
			return nil, true
		}
		out = append(out, l)
	}
	return out, false
}

func subset(a, b clause) bool {
	for _, x := range a {
		found := false
		for _, y := range b {
			if x == y {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func preprocessing(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 3 {
		for {
			m := edakit.Scaled(ctx, 6, 8)
			var cs []clause
			for i := 0; i < m; i++ {
				cs = addUnique(cs, randomClause(5, edakit.Int(1, 3)))
			}
			// plant a few supersets
			for i := 0; i < 2; i++ {
				base := cs[rand.Intn(len(cs))]
				ext := append(clause(nil), base...)
				for _, l := range randomClause(5, 2) {
					dup := false
					for _, x := range ext {
						if abs(x) == abs(l) {
							dup = true
						}
					}
					if !dup {
						ext = append(ext, l)
						break
					}
				}
				cs = append(cs, ext)
			}
			dup := false
			for i := range cs {
				for j := range cs {
					if i < j && subset(cs[i], cs[j]) && subset(cs[j], cs[i]) {
						dup = true
					}
				}
			}
			if dup {
				continue
			}
			cs = edakit.Shuffle(cs)
			n := 0
			var names []string
			for j := range cs {
				for i := range cs {
					if i != j && subset(cs[i], cs[j]) {
						n++
						names = append(names, edakit.Sprintf("C%d", j+1))
						break
					}
				}
			}
			q := "Clauses: " + cnfText(cs) + ". A clause is subsumed if another clause of the list is a subset of it. How many clauses can subsumption delete?"
			why := "No clause is subsumed."
			if n > 0 {
				why = "The subsumed clauses are " + strings.Join(names, ", ") + ": each contains every literal of some shorter clause, so it is implied and can go."
			}
			return edakit.Numeric(q, edakit.IntStr(n), why, map[string]string{"subsumed": edakit.IntStr(n)})
		}
	}
	np, nn := edakit.Scaled(ctx, 2, 3), edakit.Int(2, 3)
	var pos, neg []clause
	for i := 0; i < np; i++ {
		c := randomClause(4, edakit.Int(1, 2))
		for j := range c {
			c[j] = sign(c[j]) * (abs(c[j]) + 1)
		}
		pos = append(pos, append(clause{1}, c...))
	}
	for i := 0; i < nn; i++ {
		c := randomClause(4, edakit.Int(1, 2))
		for j := range c {
			c[j] = sign(c[j]) * (abs(c[j]) + 1)
		}
		neg = append(neg, append(clause{-1}, c...))
	}
	res, taut := 0, 0
	for _, a := range pos {
		for _, b := range neg {
			if _, t := resolve(a, b, 1); t {
				taut++
			} else {
				res++
			}
		}
	}
	all := append(append([]clause(nil), pos...), neg...)
	head := "Bounded variable elimination on $x_1$. The clauses containing $x_1$ or its negation are " + cnfText(all) + "."
	if kp == 1 {
		q := head + " How many non-tautological resolvents on $x_1$ does elimination produce (count each pair of a positive and a negative clause once)?"
		return edakit.Numeric(q, edakit.IntStr(res),
			edakit.Sprintf("There are %d·%d = %d pairs; %d of them resolve to a tautology (they contain some $y$ and $\\neg y$) and are dropped, leaving %d.", np, nn, np*nn, taut, res),
			map[string]string{"resolvents": edakit.IntStr(res)})
	}
	net := res - np - nn
	q := head + " Eliminating $x_1$ deletes these clauses and adds the non-tautological resolvents (count each pair once). What is the net change in the number of clauses? (Negative means the formula shrinks.)"
	return edakit.Numeric(q, edakit.IntStr(net),
		edakit.Sprintf("%d resolvents replace %d + %d = %d clauses: the change is %d. A preprocessor such as SatELite eliminates the variable only when this change is not positive.", res, np, nn, np+nn, net),
		map[string]string{"resolvents": edakit.IntStr(res)})
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}
