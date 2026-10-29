package fv

import (
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// ---------- BMC unrolling ----------

func bmcUnrolling(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	n, m, g := edakit.Int(2, 12), edakit.Int(1, 8), edakit.Scaled(ctx, 5, 60)
	k := edakit.Scaled(ctx, 2, 20)
	conv := edakit.Sprintf("A circuit has %d latches, %d primary inputs and %d two-input AND gates; one gate output is the bad signal. BMC unrolls frames 0 to %d. Frame 0 latch outputs are fresh variables; in later frames each latch output is the previous frame's next-state gate output, so it adds no variable. Every frame copies the inputs and the gates.", n, m, g, k)
	switch kp {
	case 1:
		ans := n + (k+1)*(m+g)
		return edakit.Numeric(conv+" How many CNF variables does the unrolling have?", edakit.IntStr(ans),
			edakit.Sprintf("%d initial latch variables plus %d frames × (%d inputs + %d gates) = %d + %d = %d.", n, k+1, m, g, n, (k+1)*(m+g), ans),
			map[string]string{"vars": edakit.IntStr(ans)})
	case 2:
		ans := n + 3*g*(k+1) + 1
		return edakit.Numeric(conv+" The CNF has one unit clause per latch for the reset state, 3 clauses per AND gate per frame, and one clause saying the bad signal is 1 in some frame. How many clauses is that?", edakit.IntStr(ans),
			edakit.Sprintf("%d reset units + 3·%d·%d gate clauses + 1 = %d.", n, g, k+1, ans),
			map[string]string{"clauses": edakit.IntStr(ans)})
	}
	ans := 3 * g
	return edakit.Numeric(conv+edakit.Sprintf(" Clauses: one unit per latch for the reset state, 3 per AND gate per frame, and one clause saying the bad signal is 1 in some frame. How many clauses are added when the depth grows from %d to %d?", k, k+1), edakit.IntStr(ans),
		edakit.Sprintf("A new frame adds 3·%d = %d gate clauses. The bad-signal clause gains one literal but stays one clause, and the reset units do not change.", g, ans),
		map[string]string{"added": edakit.IntStr(ans)})
}

// ---------- shortest counterexample ----------

func bmcDepth(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	switch kp {
	case 1:
		for {
			n := edakit.Scaled(ctx, 6, 10)
			g := randomGraph(n, 1, 2)
			bad := edakit.Int(1, n-1)
			d := g.bfs([]int{0})[bad]
			if d < 2 {
				continue
			}
			q := "States s0 to s" + edakit.IntStr(n-1) + ", initial state s0, transitions " + g.edgesText() + edakit.Sprintf(". The bad state is s%d. BMC at depth $k$ asks for a path of $k$ transitions from s0 that ends in a bad state. What is the smallest $k$ for which the SAT query is satisfiable?", bad)
			return edakit.Numeric(q, edakit.IntStr(d),
				edakit.Sprintf("The shortest path from s0 to s%d has %d transitions, so depths 0 to %d are UNSAT and depth %d returns the counterexample.", bad, d, d-1, d),
				map[string]string{"depth": edakit.IntStr(d)})
		}
	case 2:
		nb := edakit.Int(3, 6)
		mod := 1 << nb
		s, t := rand.Intn(mod), rand.Intn(mod)
		c := 2*rand.Intn(mod/2) + 1
		k := 0
		for v := s; v != t; v = (v + c) % mod {
			k++
		}
		if k == 0 {
			t = (s + c) % mod
			k = 1
		}
		q := edakit.Sprintf("A %d-bit register starts at %d and adds %d (mod %d) every cycle. The assertion says it never equals %d. At what BMC depth (number of transitions) is the first counterexample found?", nb, s, c, mod, t)
		return edakit.Numeric(q, edakit.IntStr(k),
			edakit.Sprintf("Solve $%d + %d k \\equiv %d \\pmod{%d}$ for the smallest $k \\ge 0$; because %d is odd it generates every residue, and the answer is $k = %d$.", s, c, t, mod, c, k),
			map[string]string{"depth": edakit.IntStr(k)})
	}
	D := edakit.Scaled(ctx, 3, 40)
	step := edakit.Int(2, 7)
	calls := (D+step-1)/step + 1
	q := edakit.Sprintf("The shortest counterexample to an assertion has %d transitions. A BMC script checks depths 0, %d, %d, %d, ..., and each check asks whether the bad signal is 1 in any frame up to that depth. How many solver calls run, counting the depth-0 call, until the first SAT answer?", D, step, 2*step, 3*step)
	return edakit.Numeric(q, edakit.IntStr(calls),
		edakit.Sprintf("The first checked depth at or beyond %d is %d·%d = %d. Calls at depths 0, %d, ..., %d make $\\lceil %d/%d \\rceil + 1 = %d$ calls.", D, (D+step-1)/step, step, (D+step-1)/step*step, step, (D+step-1)/step*step, D, step, calls),
		map[string]string{"calls": edakit.IntStr(calls)})
}

// ---------- diameters ----------

func longestSimple(g *graph, start int) int {
	best := 0
	seen := make([]bool, g.n)
	var dfs func(s, d int)
	dfs = func(s, d int) {
		if d > best {
			best = d
		}
		seen[s] = true
		for _, t := range g.succ[s] {
			if !seen[t] {
				dfs(t, d+1)
			}
		}
		seen[s] = false
	}
	dfs(start, 0)
	return best
}

func completeness(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	if kp == 3 {
		if rand.Intn(2) == 0 {
			n := edakit.Scaled(ctx, 2, 8)
			ans := (1 << n) - 1
			q := edakit.Sprintf("A design is a single %d-bit counter that resets to 0 and increments every cycle, wrapping around. What is its reachability diameter from the reset state, the depth BMC needs to be complete for any safety property of the counter?", n)
			return edakit.Numeric(q, edakit.IntStr(ans),
				edakit.Sprintf("The states form one cycle of length $2^{%d}$; the farthest state, %d, is %d steps from 0.", n, ans, ans),
				map[string]string{"diameter": edakit.IntStr(ans)})
		}
		a, b := edakit.Int(2, 9), edakit.Int(2, 9)
		l := a / gcd(a, b) * b
		q := edakit.Sprintf("A design holds a mod-%d counter and a mod-%d counter. Both reset to 0 and increment every cycle. How many transitions does the longest shortest path from the reset state need (the reachability diameter)?", a, b)
		return edakit.Numeric(q, edakit.IntStr(l-1),
			edakit.Sprintf("The pair returns to (0, 0) after lcm(%d, %d) = %d cycles, and every state on that cycle is reached exactly once on the way, so the farthest is %d steps away. Only %d of the %d pairs are reachable.", a, b, l, l-1, l, a*b),
			map[string]string{"lcm": edakit.IntStr(l)})
	}
	for {
		n := edakit.Scaled(ctx, 5, 7)
		g := randomGraph(n, 1, 2)
		d := g.bfs([]int{0})
		diam, reach := 0, 0
		for _, x := range d {
			if x > diam {
				diam = x
			}
			if x >= 0 {
				reach++
			}
		}
		if reach < 4 || diam < 2 {
			continue
		}
		head := "States s0 to s" + edakit.IntStr(n-1) + ", initial state s0, transitions " + g.edgesText() + "."
		if kp == 1 {
			q := head + " What is the reachability diameter from s0: the largest breadth-first distance from s0 to a reachable state?"
			return edakit.Numeric(q, edakit.IntStr(diam),
				edakit.Sprintf("Breadth-first search from s0 reaches its last new state at distance %d. BMC that checks the bad signal in every frame up to depth %d is complete for this system.", diam, diam),
				map[string]string{"diameter": edakit.IntStr(diam)})
		}
		rd := longestSimple(g, 0)
		if rd == diam && rand.Intn(2) == 0 {
			continue
		}
		q := head + " What is the recurrence diameter from s0: the largest number of transitions on a path from s0 that never revisits a state?"
		return edakit.Numeric(q, edakit.IntStr(rd),
			edakit.Sprintf("A depth-first search over loop-free paths from s0 finds a longest one of %d transitions, against a reachability diameter of %d. Any longer path must repeat a state, which is why the recurrence diameter bounds BMC completeness without knowing shortest paths.", rd, diam),
			map[string]string{"recurrence": edakit.IntStr(rd), "diameter": edakit.IntStr(diam)})
	}
}

// ---------- inductive invariants ----------

const (
	invInductive = "P is an inductive invariant"
	invInit      = "P fails initiation: an initial state violates P"
	invNotInd    = "P is an invariant but not inductive: an unreachable P-state steps outside P"
	invNotInv    = "P is not an invariant: a reachable state violates P"
)

type sysP struct {
	g    *graph
	init []int
	p    []bool
}

func (s sysP) text() string {
	return "A system has states s0 to s" + edakit.IntStr(s.g.n-1) + ", initial states " + setText(s.init) + " and transitions " + s.g.edgesText() + ". The property is $P = $ " + setText(members(s.p)) + "."
}

func classify(s sysP) string {
	for _, i := range s.init {
		if !s.p[i] {
			return invInit
		}
	}
	d := s.g.bfs(s.init)
	for i, x := range d {
		if x >= 0 && !s.p[i] {
			return invNotInv
		}
	}
	for i := 0; i < s.g.n; i++ {
		if !s.p[i] {
			continue
		}
		for _, t := range s.g.succ[i] {
			if !s.p[t] {
				return invNotInd
			}
		}
	}
	return invInductive
}

func inductiveInvariant(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		n := 6
		if kp >= 2 {
			n = edakit.Scaled(ctx, 7, 8)
		}
		g := randomGraph(n, 1, 2)
		init := []int{0}
		reach := g.bfs(init)
		p := make([]bool, n)
		mode := rand.Intn(4)
		if kp == 3 && rand.Intn(3) > 0 {
			mode = 2
		}
		for i := range p {
			switch mode {
			case 0, 2:
				p[i] = reach[i] >= 0 || rand.Intn(2) == 0
			default:
				p[i] = rand.Intn(3) > 0
			}
		}
		s := sysP{g, init, p}
		c := classify(s)
		if (mode == 2 && c != invNotInd) || count(p) == n {
			continue
		}
		why := map[string]string{
			invInductive: "Every initial state is in P and every transition from a P-state lands in P, so P holds forever by induction.",
			invInit:      "An initial state lies outside P, so P fails at time 0 and the induction cannot start.",
			invNotInd:    "Every reachable state satisfies P, but some P-state that cannot be reached has a transition leaving P, so the step case fails. This is a counterexample to induction, not a bug.",
			invNotInv:    "Following transitions from the initial states reaches a state outside P, a real counterexample.",
		}[c]
		var ds []string
		for _, o := range []string{invInductive, invInit, invNotInd, invNotInv} {
			if o != c {
				ds = append(ds, o)
			}
		}
		return edakit.Choice(s.text()+" Which statement is true?", c, why, ds...)
	}
}

// ---------- k-induction ----------

// kInductionDepth returns the minimal k for plain k-induction (or -1 if none exists) and
// the minimal k under the simple-path constraint.
func kInductionDepth(g *graph, p []bool) (int, int) {
	pre := make([]bool, g.n) // P-states with a ¬P successor
	for s := 0; s < g.n; s++ {
		if !p[s] {
			continue
		}
		for _, t := range g.succ[s] {
			if !p[t] {
				pre[s] = true
			}
		}
	}
	r := g.reverse()
	// longest P-path ending at s, in states; -1 means unbounded (a cycle reaches s)
	memo := map[int]int{}
	onStack := map[int]bool{}
	inf := false
	var longest func(s int) int
	longest = func(s int) int {
		if v, ok := memo[s]; ok {
			return v
		}
		if onStack[s] {
			inf = true
			return 0
		}
		onStack[s] = true
		best := 1
		for _, u := range r.succ[s] {
			if p[u] {
				if v := longest(u) + 1; v > best {
					best = v
				}
			}
		}
		onStack[s] = false
		memo[s] = best
		return best
	}
	plain := 0
	simple := 0
	seen := make([]bool, g.n)
	var dfs func(s, d int)
	dfs = func(s, d int) {
		if d > simple {
			simple = d
		}
		seen[s] = true
		for _, u := range r.succ[s] {
			if p[u] && !seen[u] {
				dfs(u, d+1)
			}
		}
		seen[s] = false
	}
	for s := 0; s < g.n; s++ {
		if pre[s] {
			if v := longest(s); v > plain {
				plain = v
			}
			dfs(s, 1)
		}
	}
	if inf {
		return -1, simple + 1
	}
	return plain + 1, simple + 1
}

// kInductionSystem builds a system whose reachable part satisfies P, plus an unreachable
// P-region leading into ¬P states.
func kInductionSystem(chainLen int, branching, cyclic bool) (*graph, []int, []bool) {
	nr := edakit.Int(2, 3)
	nu := chainLen
	if branching {
		nu = chainLen + edakit.Int(1, 2)
	}
	nb := 1
	n := nr + nu + nb
	ids := rand.Perm(n)
	R := ids[:nr]
	U := ids[nr : nr+nu]
	B := ids[nr+nu]
	g := newGraph(n)
	p := make([]bool, n)
	for _, s := range R {
		p[s] = true
	}
	for _, s := range U {
		p[s] = true
	}
	for i, s := range R {
		g.add(s, R[(i+1)%nr])
		if rand.Intn(2) == 0 {
			g.add(s, R[rand.Intn(nr)])
		}
	}
	// main chain U[0] -> ... -> U[chainLen-1] -> B
	for i := 0; i < chainLen-1; i++ {
		g.add(U[i], U[i+1])
	}
	g.add(U[chainLen-1], B)
	for i := chainLen; i < nu; i++ {
		// side states feed into the chain at a random point
		g.add(U[i], U[rand.Intn(chainLen)])
	}
	for _, s := range U {
		if rand.Intn(3) == 0 {
			g.add(s, R[rand.Intn(nr)])
		}
	}
	if cyclic {
		j := rand.Intn(chainLen)
		g.add(U[j], U[rand.Intn(j+1)])
	}
	g.add(B, edakit.Pick([]int{B, U[0], R[0]}))
	return g, []int{R[0]}, p
}

func kInduction(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		L := edakit.Scaled(ctx, 2, 4)
		g, init, p := kInductionSystem(L, kp >= 2, kp == 3)
		reach := g.bfs(init)
		ok := true
		for i, d := range reach {
			if d >= 0 && !p[i] {
				ok = false
			}
		}
		if !ok {
			continue
		}
		plain, simple := kInductionDepth(g, p)
		s := sysP{g, init, p}
		if kp == 3 {
			if plain != -1 {
				continue
			}
			q := s.text() + " P holds on every reachable state, but a cycle of unreachable P-states leads out of P, so plain k-induction never succeeds. With the simple-path constraint (the k states assumed in the step case are pairwise distinct), what is the smallest k that proves P?"
			return edakit.Numeric(q, edakit.IntStr(simple),
				edakit.Sprintf("The longest loop-free path of P-states that ends in a state with a successor outside P has %d states, so the step case fails up to k = %d and holds at k = %d.", simple-1, simple-1, simple),
				map[string]string{"k": edakit.IntStr(simple)})
		}
		if plain < 2 {
			continue
		}
		q := s.text() + " P holds on every reachable state, so every base case passes. The k-induction step assumes P on k consecutive states and must prove P on the next. What is the smallest k that proves P?"
		return edakit.Numeric(q, edakit.IntStr(plain),
			edakit.Sprintf("The step case fails while some path of k P-states can be followed by a ¬P state. The longest such P-path has %d states, so k = %d is the first depth whose step case holds.", plain-1, plain),
			map[string]string{"k": edakit.IntStr(plain)})
	}
}

// ---------- relative induction and generalization over 3-bit states ----------

type sys3 struct {
	g    *graph
	init []bool
	f    []bool
}

func random3() sys3 {
	g := randomGraph(8, 1, 2)
	init := make([]bool, 8)
	init[0] = true
	f := make([]bool, 8)
	f[0] = true
	for i := 1; i < 8; i++ {
		f[i] = rand.Intn(2) == 0
	}
	return sys3{g, init, f}
}

func (s sys3) edgesText() string {
	var parts []string
	for a := 0; a < 8; a++ {
		var ts []string
		for _, b := range s.g.succ[a] {
			ts = append(ts, bitsN(b, 3))
		}
		parts = append(parts, bitsN(a, 3)+" → "+strings.Join(ts, ", "))
	}
	return strings.Join(parts, "; ")
}

func bitSet(xs []bool) string {
	var s []string
	for i, b := range xs {
		if b {
			s = append(s, bitsN(i, 3))
		}
	}
	return "{" + strings.Join(s, ", ") + "}"
}

// cube is a partial assignment to x2 x1 x0: mask of fixed bits and their values.
type cube struct{ mask, val int }

func (c cube) has(s int) bool { return s&c.mask == c.val&c.mask }

func (c cube) text() string {
	var lits []string
	for i := 2; i >= 0; i-- {
		if c.mask>>i&1 == 0 {
			continue
		}
		if c.val>>i&1 == 1 {
			lits = append(lits, edakit.Sprintf("x%d", i))
		} else {
			lits = append(lits, edakit.Sprintf("¬x%d", i))
		}
	}
	return strings.Join(lits, " ∧ ")
}

// ctiCount counts states in F outside the cube that have a successor inside it.
func (s sys3) ctiCount(c cube) int {
	n := 0
	for a := 0; a < 8; a++ {
		if !s.f[a] || c.has(a) {
			continue
		}
		for _, b := range s.g.succ[a] {
			if c.has(b) {
				n++
				break
			}
		}
	}
	return n
}

func (s sys3) initHits(c cube) int {
	n := 0
	for a := 0; a < 8; a++ {
		if s.init[a] && c.has(a) {
			n++
		}
	}
	return n
}

func relativeInduction(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		s := random3()
		c := cube{mask: 7, val: rand.Intn(8)}
		if kp >= 2 {
			c.mask = edakit.Pick([]int{3, 5, 6})
		}
		if kp == 3 && rand.Intn(2) == 0 {
			c.val = 0
		}
		cti := s.ctiCount(c)
		if kp < 3 && (s.initHits(c) > 0 || (cti == 0 && rand.Intn(2) == 0)) {
			continue
		}
		head := "States are bit vectors x2x1x0. Transitions: " + s.edgesText() + ". The initial state is 000 and the frame is $F = $ " + bitSet(s.f) + ". The candidate clause is $c = \\neg(" + strings.ReplaceAll(strings.ReplaceAll(c.text(), "¬", "\\neg "), "∧", "\\wedge") + ")$."
		if kp == 3 {
			ih := s.initHits(c)
			q := head + " IC3 needs initiation (no initial state violates $c$) and consecution ($F \\wedge c \\wedge T \\rightarrow c'$). How many violations are there in total: initial states violating $c$ plus states of $F \\wedge c$ with a successor violating $c$?"
			return edakit.Numeric(q, edakit.IntStr(ih+cti),
				edakit.Sprintf("%d initial state(s) lie in the cube, and %d state(s) of $F$ outside the cube step into it: %d in total. Zero would mean $c$ can be added to the next frame.", ih, cti, ih+cti),
				map[string]string{"init": edakit.IntStr(ih), "cti": edakit.IntStr(cti)})
		}
		q := head + " How many states of $F \\wedge c$ have a successor that violates $c$? (Zero means $c$ is inductive relative to $F$.)"
		return edakit.Numeric(q, edakit.IntStr(cti),
			edakit.Sprintf("The states of $F$ outside the cube %s are checked one by one; %d of them have a transition into the cube.", c.text(), cti),
			map[string]string{"cti": edakit.IntStr(cti)})
	}
}

func generalizationCompute(ctx generator.GeneratorContext) generator.Problem {
	for {
		s := random3()
		st := 1 + rand.Intn(7)
		s.f[st] = false
		full := cube{7, st}
		if s.initHits(full) > 0 || s.ctiCount(full) > 0 {
			continue
		}
		var valid []int
		for _, bit := range []int{2, 1, 0} {
			c := cube{7 &^ (1 << bit), st}
			if s.initHits(c) == 0 && s.ctiCount(c) == 0 {
				valid = append(valid, bit)
			}
		}
		if len(valid) > 1 {
			continue
		}
		opts := []string{"Drop the literal on x2", "Drop the literal on x1", "Drop the literal on x0", "No single literal can be dropped"}
		ans := opts[3]
		if len(valid) == 1 {
			ans = opts[2-valid[0]]
		}
		var ds []string
		for _, o := range opts {
			if o != ans {
				ds = append(ds, o)
			}
		}
		q := "States are bit vectors x2x1x0. Transitions: " + s.edgesText() + ". The initial state is 000 and the frame is $F = $ " + bitSet(s.f) + ". IC3 blocks the cube " + full.text() + " (the state " + bitsN(st, 3) + "), which is inductive relative to $F$. Generalization tries to drop one literal: the smaller cube must still contain no initial state, and no state of $F$ outside it may step into it. Which choice is correct?"
		why := "Dropping a literal doubles the cube. Only the cube that avoids 000 and receives no transition from the rest of F survives the check."
		if len(valid) == 0 {
			why = "Each 2-literal cube either contains 000 or receives a transition from a state of F outside it, so the full cube must be blocked as it is."
		}
		return edakit.Choice(q, ans, why, ds...)
	}
}

// ---------- CEGAR: existential abstraction hiding x0 ----------

func cegarCompute(ctx generator.GeneratorContext) generator.Problem {
	for {
		g := randomGraph(8, 1, 1)
		if rand.Intn(2) == 0 {
			g = randomGraph(8, 1, 2)
		}
		conc := g.bfs([]int{0})
		cr := 0
		cproj := map[int]bool{}
		for i, d := range conc {
			if d >= 0 {
				cr++
				cproj[i>>1] = true
			}
		}
		ag := newGraph(4)
		for a := 0; a < 8; a++ {
			for _, b := range g.succ[a] {
				ag.add(a>>1, b>>1)
			}
		}
		ad := ag.bfs([]int{0})
		ar := 0
		for _, d := range ad {
			if d >= 0 {
				ar++
			}
		}
		if ar == 1 && rand.Intn(2) == 0 {
			continue
		}
		var parts []string
		for a := 0; a < 8; a++ {
			var ts []string
			for _, b := range g.succ[a] {
				ts = append(ts, bitsN(b, 3))
			}
			parts = append(parts, bitsN(a, 3)+" → "+strings.Join(ts, ", "))
		}
		q := "A concrete system has states x2x1x0 with transitions " + strings.Join(parts, "; ") + " and initial state 000. Localization abstraction hides x0: an abstract state is the pair x2x1, and there is an abstract transition whenever some concrete transition maps onto it. How many abstract states are reachable from 00?"
		why := edakit.Sprintf("The abstraction has an edge for every concrete edge, whether or not its source is reachable, so it reaches %d abstract states, while the reachable concrete states project onto only %d. The extra ones are where spurious counterexamples come from.", ar, len(cproj))
		return edakit.Choice(q, edakit.IntStr(ar), why, edakit.Distinct(ar, 3, len(cproj), cr, 4, 2, 1)...)
	}
}

// ---------- Houdini filtering of agent-proposed invariants ----------

func houdini(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		n := 8
		g := randomGraph(n, 1, 2)
		init := []int{0}
		reach := g.bfs(init)
		k := edakit.Scaled(ctx, 4, 6)
		var cands [][]bool
		for i := 0; i < k; i++ {
			c := make([]bool, n)
			for s := 0; s < n; s++ {
				if rand.Intn(3) == 0 {
					c[s] = reach[s] >= 0 || rand.Intn(2) == 0
				} else {
					c[s] = rand.Intn(4) > 0
				}
			}
			c[0] = c[0] || rand.Intn(4) > 0
			cands = append(cands, c)
		}
		alive := make([]bool, k)
		for i := range alive {
			alive[i] = true
		}
		rounds, afterFirst := 0, -1
		for {
			rounds++
			inv := make([]bool, n)
			for s := range inv {
				inv[s] = true
				for i, c := range cands {
					if alive[i] && !c[s] {
						inv[s] = false
					}
				}
			}
			var drop []int
			for i, c := range cands {
				if !alive[i] {
					continue
				}
				bad := !c[0]
				for s := 0; s < n && !bad; s++ {
					if !inv[s] {
						continue
					}
					for _, t := range g.succ[s] {
						if !c[t] {
							bad = true
						}
					}
				}
				if bad {
					drop = append(drop, i)
				}
			}
			for _, i := range drop {
				alive[i] = false
			}
			if afterFirst < 0 {
				afterFirst = count(alive)
			}
			if len(drop) == 0 {
				break
			}
		}
		final := count(alive)
		if rounds < 3 && kp != 1 || final == 0 && rand.Intn(2) == 0 {
			continue
		}
		var cs []string
		for i, c := range cands {
			cs = append(cs, edakit.Sprintf("c%d = %s", i+1, setText(members(c))))
		}
		head := "An agent proposes candidate invariants for a system with states s0 to s7, initial state s0 and transitions " + g.edgesText() + ". Candidates (each is the set of states it allows): " + strings.Join(cs, "; ") + ". Houdini runs rounds: it takes the conjunction (intersection) of all surviving candidates, and removes every candidate that excludes s0 or that some state of the conjunction can leave in one transition; it repeats until a round removes nothing."
		switch kp {
		case 1:
			return edakit.Numeric(head+" How many candidates survive the first round?", edakit.IntStr(afterFirst),
				edakit.Sprintf("In round 1 the conjunction is the intersection of all %d candidates; checking each candidate against its successors leaves %d.", k, afterFirst),
				map[string]string{"after_first": edakit.IntStr(afterFirst)})
		case 2:
			return edakit.Numeric(head+" How many candidates survive at the fixpoint?", edakit.IntStr(final),
				edakit.Sprintf("Removing candidates weakens the conjunction, which can expose new failures, so the rounds repeat. After %d rounds, %d candidates survive; their conjunction is inductive.", rounds, final),
				map[string]string{"final": edakit.IntStr(final), "rounds": edakit.IntStr(rounds)})
		}
		return edakit.Numeric(head+" How many rounds run, counting the final round that removes nothing?", edakit.IntStr(rounds),
			edakit.Sprintf("Rounds keep removing candidates until the survivors' conjunction is inductive; that takes %d rounds including the last, empty one, ending with %d survivors.", rounds, final),
			map[string]string{"final": edakit.IntStr(final), "rounds": edakit.IntStr(rounds)})
	}
}
