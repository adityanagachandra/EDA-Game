package boolalg

import (
	"math"
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// A clause is a list of non-zero literals: +v for x_v and -v for not x_v, v from 1.
type clause []int

func randClause(n, k int) clause {
	vs := rand.Perm(n)[:k]
	sort.Ints(vs)
	c := make(clause, k)
	for i, v := range vs {
		c[i] = v + 1
		if rand.Intn(2) == 0 {
			c[i] = -c[i]
		}
	}
	return c
}

func litTeX(l int) string {
	if l < 0 {
		return `\neg x_{` + edakit.IntStr(-l) + `}`
	}
	return `x_{` + edakit.IntStr(l) + `}`
}

func (c clause) tex() string {
	parts := make([]string, len(c))
	for i, l := range c {
		parts[i] = litTeX(l)
	}
	return "(" + strings.Join(parts, ` \lor `) + ")"
}

func cnfTeX(cs []clause) string {
	parts := make([]string, len(cs))
	for i, c := range cs {
		parts[i] = c.tex()
	}
	return "$" + strings.Join(parts, ` \land `) + "$"
}

// sat reports whether assignment a (bit v-1 holds x_v) satisfies c.
func (c clause) sat(a int) bool {
	for _, l := range c {
		v := l
		if v < 0 {
			v = -v
		}
		val := a>>uint(v-1)&1 == 1
		if (l > 0) == val {
			return true
		}
	}
	return false
}

func countModels(n int, cs []clause) int {
	k := 0
	for a := 0; a < pow2(n); a++ {
		ok := true
		for _, c := range cs {
			if !c.sat(a) {
				ok = false
				break
			}
		}
		if ok {
			k++
		}
	}
	return k
}

func clausesGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 6)
		m := edakit.Scaled(ctx, 4, 7)
		var cs []clause
		for i := 0; i < m; i++ {
			cs = append(cs, randClause(n, edakit.Int(2, 3)))
		}
		a := rand.Intn(pow2(n))
		var asg []string
		for v := 1; v <= n; v++ {
			asg = append(asg, edakit.Sprintf("x_{%d} = %d", v, a>>uint(v-1)&1))
		}
		bad := 0
		for _, c := range cs {
			if !c.sat(a) {
				bad++
			}
		}
		ans := edakit.IntStr(bad)
		q := "Under the assignment $" + strings.Join(asg, ", ") + "$, how many clauses of " + cnfTeX(cs) + " are falsified?"
		return edakit.Numeric(q, ans, "A clause is falsified only when every one of its literals is false. "+ans+" clauses are.", nil)
	case 2:
		n := edakit.Scaled(ctx, 4, 10)
		k := edakit.Int(2, minInt(4, n-1))
		c := randClause(n, k)
		if rand.Intn(2) == 0 {
			ans := edakit.IntStr(pow2(n - k))
			q := edakit.Sprintf("Over the %d variables $x_1$ to $x_{%d}$, how many complete assignments falsify the clause $%s$?", n, n, c.tex())
			return edakit.Numeric(q, ans, edakit.Sprintf("Exactly one assignment of its %d variables falsifies the clause, and the other %d variables are free: 2^%d = %s.", k, n-k, n-k, ans), nil)
		}
		ans := edakit.IntStr(pow2(n) - pow2(n-k))
		q := edakit.Sprintf("Over the %d variables $x_1$ to $x_{%d}$, how many complete assignments satisfy the clause $%s$?", n, n, c.tex())
		return edakit.Numeric(q, ans, edakit.Sprintf("2^%d assignments falsify it (one setting of its %d variables, the rest free), so 2^%d - 2^%d = %s satisfy it.", n-k, k, n, n-k, ans), nil)
	default:
		if rand.Intn(2) == 0 {
			n := edakit.Scaled(ctx, 2, 8)
			if rand.Intn(2) == 0 {
				ans := edakit.IntStr(pow2(n - 1))
				q := edakit.Sprintf("How many clauses does the smallest CNF of the parity function $x_1 \\oplus x_2 \\oplus \\dots \\oplus x_{%d}$ (with no auxiliary variables) have?", n)
				return edakit.Numeric(q, ans, edakit.Sprintf("Every clause of a parity CNF must mention all %d variables, since fixing fewer leaves the parity undetermined; so each clause excludes exactly one OFF row. There are 2^%d = %s OFF rows.", n, n-1, ans), nil)
			}
			ans := edakit.IntStr(n * pow2(n-1))
			q := edakit.Sprintf("The smallest CNF of $x_1 \\oplus \\dots \\oplus x_{%d}$ without auxiliary variables has one clause per OFF row. How many literals does it contain in total?", n)
			return edakit.Numeric(q, ans, edakit.Sprintf("There are 2^%d OFF rows, each a clause of %d literals: %s.", n-1, n, ans), nil)
		}
		n := nVars(ctx)
		f, fs := drawFunc(n, 0.6)
		ans := edakit.IntStr(rows(n) - f.ones())
		return edakit.Numeric("Let "+fs+". How many clauses does its canonical CNF (one maxterm per OFF row) have?", ans,
			edakit.Sprintf("f is 1 on %d of %d rows, so %s rows are OFF, each excluded by one clause.", f.ones(), rows(n), ans), nil)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type gateSpec struct {
	kind string
	k    int
}

func tseitinClauses(g gateSpec) int {
	switch g.kind {
	case "AND", "OR", "NAND", "NOR":
		return g.k + 1
	case "XOR", "XNOR":
		return 4
	default: // NOT, BUF
		return 2
	}
}

var numberWords = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

func (g gateSpec) describe() string {
	switch g.kind {
	case "NOT":
		return "inverter"
	case "BUF":
		return "buffer"
	case "XOR", "XNOR":
		return "two-input " + g.kind
	}
	return numberWords[g.k] + "-input " + g.kind
}

func randGate(maxK int) gateSpec {
	kind := edakit.Pick([]string{"AND", "OR", "NAND", "NOR", "XOR", "XNOR", "NOT", "AND", "OR"})
	switch kind {
	case "XOR", "XNOR":
		return gateSpec{kind, 2}
	case "NOT":
		return gateSpec{kind, 1}
	}
	return gateSpec{kind, edakit.Int(2, maxK)}
}

const tseitinRule = "Use the standard Tseitin encoding: one fresh variable per gate output, the clauses of output if and only if gate function, and no simplification."

func tseitin(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		g := randGate(8)
		if rand.Intn(4) == 0 {
			g = gateSpec{"BUF", 1}
		}
		ans := edakit.IntStr(tseitinClauses(g))
		q := tseitinRule + " How many clauses encode a single " + g.describe() + " gate?"
		why := map[string]string{
			"AND":  "For y = AND of k inputs: k clauses (not y or x_i), one per input, plus one clause (y or not x_1 or ... or not x_k).",
			"NAND": "A NAND is an AND with the output literal negated: k two-literal clauses plus one long clause.",
			"OR":   "For y = OR of k inputs: k clauses (y or not x_i) plus one clause (not y or x_1 or ... or x_k).",
			"NOR":  "A NOR is an OR with the output literal negated: k two-literal clauses plus one long clause.",
			"XOR":  "y = a XOR b needs one clause per falsifying row of the 3-variable relation: 4 clauses of 3 literals.",
			"XNOR": "y = a XNOR b, like XOR, needs 4 three-literal clauses.",
			"NOT":  "y = not a is (y or a) and (not y or not a): 2 clauses.",
			"BUF":  "y = a is (not y or a) and (y or not a): 2 clauses.",
		}[g.kind]
		return edakit.Numeric(q, ans, why+" Total: "+ans+".", nil)
	case 2:
		k := edakit.Scaled(ctx, 3, 7)
		var gs []gateSpec
		var desc []string
		total := 0
		for i := 0; i < k; i++ {
			g := randGate(4)
			gs = append(gs, g)
			total += tseitinClauses(g)
		}
		counts := map[string]int{}
		var order []string
		for _, g := range gs {
			d := g.describe()
			if counts[d] == 0 {
				order = append(order, d)
			}
			counts[d]++
		}
		for _, d := range order {
			desc = append(desc, gateCount(counts[d], d))
		}
		total++ // the unit clause asserting the output
		ans := edakit.IntStr(total)
		q := tseitinRule + " A netlist has one output and the gates: " + strings.Join(desc, ", ") + ". To ask a SAT solver whether the output can be 1, you encode every gate and add the unit clause asserting the output. How many clauses is that?"
		return edakit.Numeric(q, ans, edakit.Sprintf("Sum the per-gate counts (k + 1 for a k-input AND, OR, NAND or NOR; 4 for a two-input XOR or XNOR; 2 for an inverter) and add 1 for the output unit clause: %s.", ans), nil)
	default:
		pi := edakit.Scaled(ctx, 3, 12)
		ng := edakit.Scaled(ctx, 4, 40)
		inv := edakit.Int(0, ng/3)
		ans := edakit.IntStr(pi + ng)
		q := edakit.Sprintf("%s A circuit has %d primary inputs and %d gates, %d of them inverters. Every gate, inverters included, gets its own variable. How many Boolean variables does the CNF use?", tseitinRule, pi, ng, inv)
		return edakit.Numeric(q, ans, edakit.Sprintf("One variable per primary input plus one per gate output: %d + %d = %s. The encoding is linear in circuit size, which is why circuit SAT and CNF-SAT are equally hard.", pi, ng, ans), nil)
	}
}

func modelCount(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 4)
		m := edakit.Int(2, 4)
		var cs []clause
		for i := 0; i < m; i++ {
			cs = append(cs, randClause(n, edakit.Int(1, minInt(3, n))))
		}
		ans := edakit.IntStr(countModels(n, cs))
		q := edakit.Sprintf("Over $x_1$ to $x_{%d}$, how many assignments satisfy %s?", n, cnfTeX(cs))
		return edakit.Numeric(q, ans, edakit.Sprintf("Enumerate all %d assignments and test every clause: %s satisfy the formula.", pow2(n), ans), nil)
	case 2:
		n := edakit.Scaled(ctx, 4, 6)
		m := edakit.Int(2, 3)
		var cs []clause
		for i := 0; i < m; i++ {
			cs = append(cs, randClause(n, edakit.Int(2, 3)))
		}
		ans := edakit.IntStr(countModels(n, cs))
		q := edakit.Sprintf("Over $x_1$ to $x_{%d}$, count the satisfying assignments of %s by inclusion and exclusion over the sets of assignments that falsify each clause.", n, cnfTeX(cs))
		return edakit.Numeric(q, ans, edakit.Sprintf("A clause with k literals is falsified by 2^(%d-k) assignments. Two clauses can be falsified together only if they share no complementary literal; then the count is 2^(%d - variables in their union). Subtracting the union of the falsifying sets from 2^%d leaves %s.", n, n, n, ans), nil)
	default:
		// Two variable-disjoint groups plus free variables.
		a := edakit.Int(2, 3)
		b := edakit.Int(2, 3)
		free := edakit.Int(0, 2)
		n := a + b + free
		perm := rand.Perm(n)
		groupA := perm[:a]
		groupB := perm[a : a+b]
		mk := func(group []int, m int) []clause {
			var out []clause
			for i := 0; i < m; i++ {
				k := edakit.Int(2, len(group))
				vs := edakit.Shuffle(group)[:k]
				sort.Ints(vs)
				c := make(clause, k)
				for j, v := range vs {
					c[j] = v + 1
					if rand.Intn(2) == 0 {
						c[j] = -c[j]
					}
				}
				out = append(out, c)
			}
			return out
		}
		ca := mk(groupA, edakit.Int(1, 2))
		cb := mk(groupB, edakit.Int(1, 2))
		var all []clause
		seenC := map[string]bool{}
		for _, c := range append(append([]clause(nil), ca...), cb...) {
			if k := clauseKey(c); !seenC[k] {
				seenC[k] = true
				all = append(all, c)
			}
		}
		all = edakit.Shuffle(all)
		ans := edakit.IntStr(countModels(n, all))
		q := edakit.Sprintf("Over $x_1$ to $x_{%d}$, how many assignments satisfy %s? (Some variables may appear in no clause.)", n, cnfTeX(all))
		return edakit.Numeric(q, ans, edakit.Sprintf("The clauses split into groups with disjoint variables. Count each group over its own variables, multiply, and multiply by 2 for every variable in no clause (%d of them here): %s.", free, ans), nil)
	}
}

// resolve returns the resolvent of c1 and c2 on variable v (c1 has +v, c2 has -v), and
// whether it is tautological.
func resolveOn(c1, c2 clause, v int) (clause, bool) {
	set := map[int]bool{}
	for _, l := range c1 {
		if l != v {
			set[l] = true
		}
	}
	for _, l := range c2 {
		if l != -v {
			set[l] = true
		}
	}
	taut := false
	var out clause
	for l := range set {
		if set[-l] {
			taut = true
		}
		out = append(out, l)
	}
	sort.Slice(out, func(i, j int) bool {
		ai, aj := out[i], out[j]
		if ai < 0 {
			ai = -ai
		}
		if aj < 0 {
			aj = -aj
		}
		if ai != aj {
			return ai < aj
		}
		return out[i] > out[j]
	})
	return out, taut
}

func clauseKey(c clause) string {
	s := append(clause(nil), c...)
	sort.Ints(s)
	parts := make([]string, len(s))
	for i, l := range s {
		parts[i] = edakit.IntStr(l)
	}
	return strings.Join(parts, " ")
}

func withLit(n, k, l int) clause {
	for {
		c := randClause(n, k)
		has := false
		for _, x := range c {
			if x == l || x == -l {
				has = true
			}
		}
		if !has {
			continue
		}
		for i, x := range c {
			if x == -l {
				c[i] = l
			}
		}
		return c
	}
}

func resolution(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	v := edakit.Int(1, n)
	kp := kpOf(ctx)
	switch kp {
	case 1:
		for {
			c1 := withLit(n, edakit.Int(2, 4), v)
			c2 := withLit(n, edakit.Int(2, 4), -v)
			r, taut := resolveOn(c1, c2, v)
			if taut {
				continue
			}
			ans := edakit.IntStr(len(r))
			q := edakit.Sprintf("Resolve $%s$ with $%s$ on $x_{%d}$. How many distinct literals does the resolvent contain?", c1.tex(), c2.tex(), v)
			return edakit.Numeric(q, ans, "Drop the clashing pair and take the union of what remains, merging repeated literals: $"+r.tex()+"$, "+ans+" literals.", nil)
		}
	default:
		np, nn := edakit.Int(2, 3), edakit.Int(2, 3)
		var P, N, rest []clause
		for i := 0; i < np; i++ {
			P = append(P, withLit(n, edakit.Int(2, 3), v))
		}
		for i := 0; i < nn; i++ {
			N = append(N, withLit(n, edakit.Int(2, 3), -v))
		}
		for i := 0; i < edakit.Int(1, 3); i++ {
			for {
				c := randClause(n, edakit.Int(2, 3))
				ok := true
				for _, l := range c {
					if l == v || l == -v {
						ok = false
					}
				}
				if ok {
					rest = append(rest, c)
					break
				}
			}
		}
		seen := map[string]bool{}
		for _, c := range rest {
			seen[clauseKey(c)] = true
		}
		useful := map[string]bool{}
		for _, p := range P {
			for _, q := range N {
				r, taut := resolveOn(p, q, v)
				if !taut {
					useful[clauseKey(r)] = true
				}
			}
		}
		all := append(append(append([]clause(nil), P...), N...), rest...)
		all = edakit.Shuffle(all)
		if kp == 2 {
			ans := edakit.IntStr(len(useful))
			q := edakit.Sprintf("In the CNF %s, resolve every clause containing $x_{%d}$ with every clause containing $\\neg x_{%d}$. How many distinct non-tautological resolvents result?", cnfTeX(all), v, v)
			return edakit.Numeric(q, ans, edakit.Sprintf("There are %d x %d = %d pairs; discard resolvents that contain some literal and its negation, and merge duplicates: %s remain.", np, nn, np*nn, ans), nil)
		}
		after := map[string]bool{}
		for k := range seen {
			after[k] = true
		}
		for k := range useful {
			after[k] = true
		}
		ans := edakit.IntStr(len(after))
		q := edakit.Sprintf("Eliminate $x_{%d}$ from %s by Davis-Putnam resolution: remove every clause that mentions $x_{%d}$ and add all non-tautological resolvents on it. Counting distinct clauses and applying no subsumption, how many clauses remain?", v, cnfTeX(all), v)
		return edakit.Numeric(q, ans, edakit.Sprintf("%d clauses mention x_%d and are removed; %d clauses remain untouched and %d distinct useful resolvents are added, giving %s distinct clauses. Elimination can grow the formula, which is why preprocessors only eliminate a variable when the count does not grow.", np+nn, v, len(rest), len(useful), ans), nil)
	}
}

func csp(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		d := edakit.Scaled(ctx, 3, 12)
		if rand.Intn(2) == 0 {
			ans := edakit.IntStr(1 + binom(d, 2))
			q := edakit.Sprintf("A CSP variable has a domain of %d values. With a one-hot (direct) encoding, one Boolean per value, an at-least-one clause and pairwise at-most-one clauses, how many clauses encode the domain?", d)
			return edakit.Numeric(q, ans, edakit.Sprintf("One at-least-one clause plus C(%d, 2) = %d pairwise clauses: %s.", d, binom(d, 2), ans), nil)
		}
		V := edakit.Int(3, 20)
		ans := edakit.IntStr(V * (1 + binom(d, 2)))
		q := edakit.Sprintf("A CSP has %d variables, each with a domain of %d values, encoded one-hot with an at-least-one clause and pairwise at-most-one clauses per variable. How many domain clauses are there in total, before any constraint clauses?", V, d)
		return edakit.Numeric(q, ans, edakit.Sprintf("Each variable needs 1 + C(%d, 2) = %d clauses; times %d variables gives %s.", d, 1+binom(d, 2), V, ans), nil)
	case 2:
		V := edakit.Scaled(ctx, 4, 12)
		E := edakit.Int(V-1, V*(V-1)/2)
		if E > 3*V {
			E = 3 * V
		}
		k := edakit.Int(3, 5)
		ans := edakit.IntStr(V*(1+binom(k, 2)) + E*k)
		q := edakit.Sprintf("Encode %d-coloring of a graph with %d vertices and %d edges as CNF: one Boolean per (vertex, color), each vertex gets at least one color and at most one (pairwise), and for every edge and every color the two endpoints are not both that color. How many clauses?", k, V, E)
		return edakit.Numeric(q, ans, edakit.Sprintf("Vertices: %d x (1 + C(%d, 2)) = %d. Edges: %d x %d = %d. Total %s.", V, k, V*(1+binom(k, 2)), E, k, E*k, ans), nil)
	default:
		var d int
		for {
			d = edakit.Scaled(ctx, 3, 30)
			if d&(d-1) != 0 {
				break
			}
		}
		b := int(math.Ceil(math.Log2(float64(d))))
		V := edakit.Int(1, 10)
		if rand.Intn(2) == 0 {
			ans := edakit.IntStr(V * (pow2(b) - d))
			q := edakit.Sprintf("A CSP has %d variables, each with a domain of %d values, in a binary (log) encoding with the fewest bits per variable. Each unused bit pattern is forbidden by one clause. How many such exclusion clauses are needed in total?", V, d)
			return edakit.Numeric(q, ans, edakit.Sprintf("%d values need %d bits, which have %d patterns; %d are unused per variable. Times %d variables: %s.", d, b, pow2(b), pow2(b)-d, V, ans), nil)
		}
		ans := edakit.IntStr(V * b)
		q := edakit.Sprintf("A CSP has %d variables, each with a domain of %d values. How many Boolean variables does a binary (log) encoding with the fewest bits per variable use?", V, d)
		return edakit.Numeric(q, ans, edakit.Sprintf("The ceiling of log2 %d is %d bits per variable; times %d is %s, against %d for one-hot.", d, b, V, ans, V*d), nil)
	}
}

// gateCount renders "3 two-input XOR gates" style counts.
func gateCount(k int, d string) string {
	noun := " gates"
	if k == 1 {
		noun = " gate"
	}
	if d == "inverter" || d == "buffer" {
		if k == 1 {
			return "1 " + d
		}
		return edakit.IntStr(k) + " " + d + "s"
	}
	return edakit.IntStr(k) + " " + d + noun
}
