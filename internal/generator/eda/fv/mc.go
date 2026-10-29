package fv

import (
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// ---------- Kripke structures ----------

type nsFunc struct {
	text string
	f    func(x [3]int, in int) int
}

var nsTemplates = []nsFunc{
	{"¬x0", func(x [3]int, in int) int { return 1 - x[0] }},
	{"x0 ⊕ x1", func(x [3]int, in int) int { return x[0] ^ x[1] }},
	{"x0 ∧ x1", func(x [3]int, in int) int { return x[0] & x[1] }},
	{"x2 ∨ (x1 ∧ x0)", func(x [3]int, in int) int { return x[2] | (x[1] & x[0]) }},
	{"x2 ⊕ (x1 ∧ x0)", func(x [3]int, in int) int { return x[2] ^ (x[1] & x[0]) }},
	{"x1", func(x [3]int, in int) int { return x[1] }},
	{"x2", func(x [3]int, in int) int { return x[2] }},
	{"¬x2 ∧ x1", func(x [3]int, in int) int { return (1 - x[2]) & x[1] }},
	{"x0 ∨ x2", func(x [3]int, in int) int { return x[0] | x[2] }},
	{"i", func(x [3]int, in int) int { return in }},
	{"x0 ⊕ i", func(x [3]int, in int) int { return x[0] ^ in }},
	{"x1 ∧ i", func(x [3]int, in int) int { return x[1] & in }},
	{"¬x1 ∨ x0", func(x [3]int, in int) int { return (1 - x[1]) | x[0] }},
}

func kripke(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		for {
			n := edakit.Scaled(ctx, 6, 9)
			g := randomGraph(n, 1, 2)
			init := randomSubset(n, edakit.Int(1, 2))
			d := g.bfs(init)
			r := 0
			for _, x := range d {
				if x >= 0 {
					r++
				}
			}
			if r == n || r < 3 {
				continue
			}
			q := "A Kripke structure has states s0 to s" + edakit.IntStr(n-1) + ", initial states " + setText(init) + " and transitions " + g.edgesText() + ". How many states are reachable, counting the initial states?"
			return edakit.Numeric(q, edakit.IntStr(r),
				edakit.Sprintf("Search forward from the initial states along the transitions: %d of the %d states are reached.", r, n),
				map[string]string{"reachable": edakit.IntStr(r)})
		}
	case 2:
		for {
			var fs [3]nsFunc
			usesInput := false
			for i := 0; i < 3; i++ {
				fs[i] = nsTemplates[rand.Intn(len(nsTemplates))]
				if strings.Contains(fs[i].text, "i") {
					usesInput = true
				}
			}
			seen := map[int]bool{0: true}
			q := []int{0}
			for len(q) > 0 {
				s := q[0]
				q = q[1:]
				x := [3]int{s & 1, s >> 1 & 1, s >> 2 & 1}
				for in := 0; in < 2; in++ {
					t := 0
					for i := 0; i < 3; i++ {
						t |= fs[i].f(x, in) << i
					}
					if !seen[t] {
						seen[t] = true
						q = append(q, t)
					}
				}
			}
			r := len(seen)
			if r < 3 || r == 8 && rand.Intn(2) == 0 {
				continue
			}
			inText := " There are no inputs."
			if usesInput {
				inText = " The input i is free: it can be 0 or 1 in every cycle."
			}
			stem := "A circuit has latches x2, x1, x0, all reset to 0, with next-state functions x0' = " + fs[0].text + ", x1' = " + fs[1].text + ", x2' = " + fs[2].text + "." + inText + " Its Kripke structure has the 8 latch valuations as states. How many states are reachable from 000, counting 000?"
			return edakit.Numeric(stem, edakit.IntStr(r),
				edakit.Sprintf("Apply the next-state functions from 000 (for both input values where the input matters) until no new valuation appears: %d states are reachable.", r),
				map[string]string{"reachable": edakit.IntStr(r)})
		}
	}
	for {
		n := edakit.Scaled(ctx, 6, 9)
		g := newGraph(n)
		dead := 0
		for s := 0; s < n; s++ {
			if rand.Intn(3) == 0 {
				dead++
				continue
			}
			for i := 0; i < edakit.Int(1, 2); i++ {
				g.add(s, rand.Intn(n))
			}
		}
		if dead == 0 {
			continue
		}
		q := "A transition system has states s0 to s" + edakit.IntStr(n-1) + " and transitions " + g.edgesText() + ". A Kripke structure needs a total transition relation. How many states need an added self-loop (or other edge) to make it total?"
		return edakit.Numeric(q, edakit.IntStr(dead),
			edakit.Sprintf("Exactly the deadlock states, those with no successor, violate totality: %d of them.", dead),
			map[string]string{"deadlocks": edakit.IntStr(dead)})
	}
}

// ---------- reachability ----------

func reachability(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		n := edakit.Scaled(ctx, 7, 10)
		g := randomGraph(n, 1, 2)
		init := []int{0}
		head := "States s0 to s" + edakit.IntStr(n-1) + ", initial state s0, transitions " + g.edgesText() + "."
		d := g.bfs(init)
		depth := 0
		r2 := 0
		for _, x := range d {
			if x > depth {
				depth = x
			}
			if x >= 0 && x <= 2 {
				r2++
			}
		}
		switch kp {
		case 1:
			if depth < 2 {
				continue
			}
			q := head + " Symbolic reachability computes $R_0 = \\text{Init}$ and $R_{i+1} = R_i \\cup \\text{Img}(R_i)$. How many image steps add at least one new state before the fixpoint?"
			return edakit.Numeric(q, edakit.IntStr(depth),
				edakit.Sprintf("Each step adds the states at the next breadth-first distance from s0. The farthest reachable state is at distance %d, so %d steps add states and step %d confirms the fixpoint.", depth, depth, depth+1),
				map[string]string{"depth": edakit.IntStr(depth)})
		case 2:
			if depth < 3 {
				continue
			}
			q := head + " With $R_0 = \\{s0\\}$ and $R_{i+1} = R_i \\cup \\text{Img}(R_i)$, how many states does $R_2$ contain?"
			return edakit.Numeric(q, edakit.IntStr(r2),
				edakit.Sprintf("$R_2$ holds the states within two steps of s0: %d of them.", r2),
				map[string]string{"R2": edakit.IntStr(r2)})
		}
		bad := edakit.Int(1, n-1)
		bd := g.reverse().bfs([]int{bad})
		b := 0
		for _, x := range bd {
			if x >= 0 {
				b++
			}
		}
		if b < 3 || b == n {
			continue
		}
		q := head + edakit.Sprintf(" Backward reachability starts from the bad state $B_0 = \\{s%d\\}$ and computes $B_{i+1} = B_i \\cup \\text{Pre}(B_i)$. How many states does the fixpoint contain?", bad)
		return edakit.Numeric(q, edakit.IntStr(b),
			edakit.Sprintf("Follow the transitions backwards from s%d: %d states can reach it, counting s%d itself. The property fails exactly when s0 is among them.", bad, b, bad),
			map[string]string{"backward": edakit.IntStr(b)})
	}
}

// ---------- state explosion ----------

func stateExplosion(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		k, n := edakit.Int(2, 6), edakit.Scaled(ctx, 2, 6)
		ans := 1
		for i := 0; i < n; i++ {
			ans *= k
		}
		q := edakit.Sprintf("A model has %d independent copies of a component with %d local states each. How many global states does the product have?", n, k)
		return edakit.Numeric(q, edakit.IntStr(ans), edakit.Sprintf("The global state is a tuple of local states: $%d^{%d} = %d$.", k, n, ans), nil)
	case 2:
		var sizes []int
		prod := 1
		for i := 0; i < edakit.Int(2, 3); i++ {
			s := edakit.Int(3, 9)
			sizes = append(sizes, s)
			prod *= s
		}
		r := edakit.Scaled(ctx, 4, 16)
		ans := prod << r
		q := edakit.Sprintf("A block has FSMs with %s states and a free-running %d-bit counter, all independent. How many states does its state space have?", joinInts(sizes, ", "), r)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Multiply: %s × $2^{%d}$ = %d × %d = %d.", joinInts(sizes, " × "), r, prod, 1<<r, ans),
			map[string]string{"product": edakit.IntStr(ans)})
	}
	n, k := edakit.Scaled(ctx, 2, 8), edakit.Int(2, 5)
	ans := binom(n+k-1, k-1)
	q := edakit.Sprintf("%d identical processes each have %d local states. Symmetry reduction treats global states that differ only by a permutation of the processes as one. How many symmetry-reduced states are there?", n, k)
	full := 1
	for i := 0; i < n; i++ {
		full *= k
	}
	return edakit.Numeric(q, edakit.IntStr(ans),
		edakit.Sprintf("A reduced state records how many processes are in each local state: a multiset of size %d from %d kinds, $\\binom{%d}{%d} = %d$, against $%d^{%d} = %d$ unreduced.", n, k, n+k-1, k-1, ans, k, n, full),
		map[string]string{"reduced": edakit.IntStr(ans), "full": edakit.IntStr(full)})
}

// ---------- LTL on lasso paths ----------

type ltl struct {
	op   string // "ap", "not", "and", "or", "imp", "X", "G", "F", "U"
	atom string
	a, b *ltl
}

func ap(s string) *ltl          { return &ltl{op: "ap", atom: s} }
func lnot(a *ltl) *ltl          { return &ltl{op: "not", a: a} }
func land(a, b *ltl) *ltl       { return &ltl{op: "and", a: a, b: b} }
func lor(a, b *ltl) *ltl        { return &ltl{op: "or", a: a, b: b} }
func limp(a, b *ltl) *ltl       { return &ltl{op: "imp", a: a, b: b} }
func lun(o string, a *ltl) *ltl { return &ltl{op: o, a: a} }
func luntil(a, b *ltl) *ltl     { return &ltl{op: "U", a: a, b: b} }

func (f *ltl) text() string {
	switch f.op {
	case "ap":
		return f.atom
	case "not":
		if f.a.op == "ap" {
			return "¬" + f.a.text()
		}
		return "¬(" + f.a.text() + ")"
	case "and":
		return f.a.text() + " ∧ " + f.b.text()
	case "or":
		return f.a.text() + " ∨ " + f.b.text()
	case "imp":
		return f.a.text() + " → " + f.b.text()
	case "U":
		return f.a.text() + " U " + f.b.text()
	}
	inner := f.a.text()
	if f.a.op == "and" || f.a.op == "or" || f.a.op == "imp" || f.a.op == "U" {
		inner = "(" + inner + ")"
	}
	return f.op + " " + inner
}

type lasso struct {
	labels []map[string]bool
	loop   int // index the path returns to after the last position
}

func (l *lasso) pos(t int) int {
	n := len(l.labels)
	if t < n {
		return t
	}
	per := n - l.loop
	return l.loop + (t-l.loop)%per
}

func (l *lasso) eval(f *ltl, t int) bool {
	h := len(l.labels) * 2
	switch f.op {
	case "ap":
		return l.labels[l.pos(t)][f.atom]
	case "not":
		return !l.eval(f.a, t)
	case "and":
		return l.eval(f.a, t) && l.eval(f.b, t)
	case "or":
		return l.eval(f.a, t) || l.eval(f.b, t)
	case "imp":
		return !l.eval(f.a, t) || l.eval(f.b, t)
	case "X":
		return l.eval(f.a, t+1)
	case "G":
		for u := t; u < t+h; u++ {
			if !l.eval(f.a, u) {
				return false
			}
		}
		return true
	case "F":
		for u := t; u < t+h; u++ {
			if l.eval(f.a, u) {
				return true
			}
		}
		return false
	case "U":
		for u := t; u < t+h; u++ {
			if l.eval(f.b, u) {
				return true
			}
			if !l.eval(f.a, u) {
				return false
			}
		}
		return false
	}
	return false
}

func (l *lasso) text() string {
	var parts []string
	for i, lab := range l.labels {
		var ps []string
		for _, a := range []string{"p", "q"} {
			if lab[a] {
				ps = append(ps, a)
			}
		}
		parts = append(parts, edakit.Sprintf("s%d {%s}", i, strings.Join(ps, ", ")))
	}
	return strings.Join(parts, " → ") + edakit.Sprintf(", and from s%d the path returns to s%d and repeats the loop forever", len(l.labels)-1, l.loop)
}

func ltlPool(kp int) []*ltl {
	p, q := ap("p"), ap("q")
	switch kp {
	case 1:
		return []*ltl{lun("X", p), lun("X", q), lun("X", lun("X", p)), lun("X", lun("X", q)), lun("G", p), lun("G", q),
			lun("F", p), lun("F", q), lun("G", lnot(q)), lun("G", lnot(p)), lun("F", land(p, q)), lun("F", lnot(p))}
	case 2:
		return []*ltl{lun("G", lun("F", p)), lun("G", lun("F", q)), lun("F", lun("G", p)), lun("F", lun("G", q)),
			luntil(p, q), luntil(q, p), lun("F", lun("G", lnot(p))), lun("G", lun("F", land(p, q))), luntil(lnot(p), q), lun("F", lun("G", lnot(q)))}
	}
	return []*ltl{lun("G", limp(p, lun("F", q))), lun("G", limp(q, lun("F", p))), lun("G", limp(p, lun("X", q))), lun("G", limp(q, lun("X", p))),
		lun("G", lnot(land(p, q))), lun("G", limp(p, q)), lun("G", limp(q, p)), lun("F", lun("G", lor(p, q))), lun("G", limp(p, lun("X", lnot(p))))}
}

func ltlChoice(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		n := edakit.Scaled(ctx, 3, 5)
		l := &lasso{loop: rand.Intn(n)}
		for i := 0; i < n; i++ {
			l.labels = append(l.labels, map[string]bool{"p": rand.Intn(2) == 0, "q": rand.Intn(2) == 0})
		}
		pool := edakit.Shuffle(ltlPool(kp))
		var tr, fa []*ltl
		for _, f := range pool {
			if l.eval(f, 0) {
				tr = append(tr, f)
			} else {
				fa = append(fa, f)
			}
		}
		askHolds := rand.Intn(2) == 0
		var ans *ltl
		var ds []*ltl
		if askHolds && len(tr) >= 1 && len(fa) >= 3 {
			ans, ds = tr[0], fa[:3]
		} else if !askHolds && len(fa) >= 1 && len(tr) >= 3 {
			ans, ds = fa[0], tr[:3]
		} else {
			continue
		}
		verb := "holds"
		if !askHolds {
			verb = "fails"
		}
		var dtext []string
		for _, d := range ds {
			dtext = append(dtext, d.text())
		}
		q := "A path visits " + l.text() + ". Which LTL formula " + verb + " at the first position s0?"
		why := edakit.Sprintf("%s %s on this lasso, while each of the other three formulas %s. Evaluate G and F over every state the path visits from that point on, including the loop.", ans.text(), verb, map[bool]string{true: "fails", false: "holds"}[askHolds])
		return edakit.Choice(q, ans.text(), why, dtext...)
	}
}

// ---------- CTL ----------

type kripkeLab struct {
	g    *graph
	p, q []bool
}

func randomKripke(n int) kripkeLab {
	g := randomGraph(n, 1, 2)
	k := kripkeLab{g: g, p: make([]bool, n), q: make([]bool, n)}
	for s := 0; s < n; s++ {
		k.p[s] = rand.Intn(2) == 0
		k.q[s] = rand.Intn(3) == 0
	}
	return k
}

func (k kripkeLab) text() string {
	var labs []string
	for s := 0; s < k.g.n; s++ {
		var ps []string
		if k.p[s] {
			ps = append(ps, "p")
		}
		if k.q[s] {
			ps = append(ps, "q")
		}
		labs = append(labs, edakit.Sprintf("s%d {%s}", s, strings.Join(ps, ", ")))
	}
	return "States and labels: " + strings.Join(labs, ", ") + ". Transitions: " + k.g.edgesText() + "."
}

func (k kripkeLab) ex(s []bool) []bool {
	out := make([]bool, k.g.n)
	for x := 0; x < k.g.n; x++ {
		for _, t := range k.g.succ[x] {
			if s[t] {
				out[x] = true
			}
		}
	}
	return out
}

func (k kripkeLab) ax(s []bool) []bool {
	out := make([]bool, k.g.n)
	for x := 0; x < k.g.n; x++ {
		out[x] = true
		for _, t := range k.g.succ[x] {
			if !s[t] {
				out[x] = false
			}
		}
	}
	return out
}

func neg(s []bool) []bool {
	out := make([]bool, len(s))
	for i, b := range s {
		out[i] = !b
	}
	return out
}

func eqSet(a, b []bool) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// fix iterates step from z0 and returns the fixpoint and the number of iterations i at which
// Z_{i+1} = Z_i first holds.
func fix(z0 []bool, step func([]bool) []bool) ([]bool, int) {
	z := z0
	for i := 0; ; i++ {
		nz := step(z)
		if eqSet(nz, z) {
			return z, i
		}
		z = nz
	}
}

func (k kripkeLab) ef(s []bool) ([]bool, int) {
	return fix(s, func(z []bool) []bool {
		e := k.ex(z)
		out := make([]bool, len(z))
		for i := range z {
			out[i] = z[i] || e[i]
		}
		return out
	})
}

func (k kripkeLab) eg(s []bool) ([]bool, int) {
	return fix(s, func(z []bool) []bool {
		e := k.ex(z)
		out := make([]bool, len(z))
		for i := range z {
			out[i] = z[i] && e[i]
		}
		return out
	})
}

func (k kripkeLab) eu(a, b []bool) ([]bool, int) {
	return fix(b, func(z []bool) []bool {
		e := k.ex(z)
		out := make([]bool, len(z))
		for i := range z {
			out[i] = z[i] || (a[i] && e[i])
		}
		return out
	})
}

func (k kripkeLab) ag(s []bool) []bool { r, _ := k.ef(neg(s)); return neg(r) }
func (k kripkeLab) af(s []bool) []bool { r, _ := k.eg(neg(s)); return neg(r) }

func ctl(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		n := edakit.Scaled(ctx, 5, 7)
		k := randomKripke(n)
		var name string
		var set []bool
		switch kp {
		case 1:
			if rand.Intn(2) == 0 {
				name, set = "EX p", k.ex(k.p)
			} else {
				name, set = "AX p", k.ax(k.p)
			}
		case 2:
			switch rand.Intn(5) {
			case 0:
				name = "EF q"
				set, _ = k.ef(k.q)
			case 1:
				name, set = "AG p", k.ag(k.p)
			case 2:
				name = "EG p"
				set, _ = k.eg(k.p)
			case 3:
				name, set = "AF q", k.af(k.q)
			default:
				name = "E[p U q]"
				set, _ = k.eu(k.p, k.q)
			}
		default:
			switch rand.Intn(3) {
			case 0:
				name = "AG EF q"
				e, _ := k.ef(k.q)
				set = k.ag(e)
			case 1:
				name = "EF AG p"
				set, _ = k.ef(k.ag(k.p))
			default:
				name = "AG (p → AF q)"
				af := k.af(k.q)
				imp := make([]bool, n)
				for i := range imp {
					imp[i] = !k.p[i] || af[i]
				}
				set = k.ag(imp)
			}
		}
		c := count(set)
		if c == 0 && rand.Intn(3) > 0 || c == n && rand.Intn(3) > 0 {
			continue
		}
		q := k.text() + " In how many states does the CTL formula " + name + " hold?"
		why := edakit.Sprintf("Computing the satisfying set bottom-up (EX and AX look one step ahead, EF and EU are least fixpoints, EG is a greatest fixpoint, AG and AF are their duals) gives %s, so the count is %d.", setText(members(set)), c)
		return edakit.Numeric(q, edakit.IntStr(c), why, map[string]string{"set": setText(members(set))})
	}
}

func ctlFixpoint(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		n := edakit.Scaled(ctx, 6, 8)
		k := randomKripke(n)
		var it int
		var rule, name string
		var res []bool
		switch kp {
		case 1:
			res, it = k.ef(k.q)
			name, rule = "EF q", "$Z_0 = [[q]]$, $Z_{i+1} = Z_i \\cup \\text{EX}\\, Z_i$"
		case 2:
			res, it = k.eg(k.p)
			name, rule = "EG p", "$Z_0 = [[p]]$, $Z_{i+1} = Z_i \\cap \\text{EX}\\, Z_i$"
		default:
			res, it = k.eu(k.p, k.q)
			name, rule = "E[p U q]", "$Z_0 = [[q]]$, $Z_{i+1} = Z_i \\cup ([[p]] \\cap \\text{EX}\\, Z_i)$"
		}
		if it < 1 || (it < 2 && rand.Intn(2) == 0) {
			continue
		}
		q := k.text() + " A model checker computes " + name + " by the iteration " + rule + ". What is the smallest $i$ with $Z_{i+1} = Z_i$?"
		return edakit.Numeric(q, edakit.IntStr(it),
			edakit.Sprintf("Iterating from $Z_0$, the set first repeats at $i = %d$, stabilizing at %s; the operator is monotone on a finite lattice, so the iteration must stop.", it, setText(members(res))),
			map[string]string{"iterations": edakit.IntStr(it), "fixpoint": setText(members(res))})
	}
}
