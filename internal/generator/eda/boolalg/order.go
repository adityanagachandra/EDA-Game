package boolalg

import (
	"math/bits"
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func bitvec(x, n int) string {
	var b strings.Builder
	for v := 0; v < n; v++ {
		b.WriteByte(byte('0' + bitOf(x, v, n)))
	}
	return b.String()
}

func booleanLattice(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 8)
		v := rand.Intn(pow2(n))
		u := v & rand.Intn(pow2(n))
		cu, cv := bits.OnesCount(uint(u)), bits.OnesCount(uint(v))
		switch rand.Intn(3) {
		case 0:
			ans := edakit.IntStr(pow2(cv - cu))
			q := edakit.Sprintf("In the Boolean lattice $B^{%d}$ ordered componentwise, how many vectors $w$ satisfy %s $\\le w \\le$ %s?", n, bitvec(u, n), bitvec(v, n))
			return edakit.Numeric(q, ans, edakit.Sprintf("The lower vector has %d ones and the upper has %d, so %d positions are free: 2^%d = %s.", cu, cv, cv-cu, cv-cu, ans), nil)
		case 1:
			ans := edakit.IntStr(pow2(n - cv))
			q := edakit.Sprintf("In $B^{%d}$ ordered componentwise, how many vectors are greater than or equal to %s?", n, bitvec(v, n))
			return edakit.Numeric(q, ans, edakit.Sprintf("The %d ones must stay 1 and the %d zeros are free: 2^%d = %s.", cv, n-cv, n-cv, ans), nil)
		default:
			ans := edakit.IntStr(pow2(n-cv) + pow2(cv) - 1)
			q := edakit.Sprintf("In $B^{%d}$ ordered componentwise, how many vectors are comparable to %s (above it, below it, or equal)?", n, bitvec(v, n))
			return edakit.Numeric(q, ans, edakit.Sprintf("There are 2^%d vectors at or above and 2^%d at or below; the vector itself is counted twice: %s.", n-cv, cv, ans), nil)
		}
	case 2:
		n := edakit.Scaled(ctx, 3, 4)
		var f, g tt
		for {
			g = tt{n, rand.Uint32() & fullMask(n)}
			f = tt{n, g.w & rand.Uint32()}
			if f.ones() >= 1 && g.ones()-f.ones() >= 1 {
				break
			}
		}
		k := g.ones() - f.ones()
		ans := edakit.IntStr(pow2(k))
		q := "Let " + sigma("f", f) + " and " + sigma("g", g) + ", so $f \\le g$ ($f$ implies $g$). How many functions $h$ satisfy $f \\le h \\le g$?"
		return edakit.Numeric(q, ans, edakit.Sprintf("h must contain the ON-set of f and stay inside the ON-set of g, so the %d rows in g but not in f are free: 2^%d = %s. This interval is exactly an incompletely specified function.", k, k, ans), nil)
	default:
		n := nVars(ctx)
		f, fs := drawFunc(n, 0)
		g := tt{n, rand.Uint32() & fullMask(n)}
		if rand.Intn(2) == 0 {
			r := f.or(g)
			ans := edakit.IntStr(r.ones())
			return edakit.Numeric("Let "+fs+" and "+sigma("g", g)+". In the lattice of Boolean functions ordered by implication, how many ON rows does the join $f \\vee g$ have?", ans,
				"The join (least upper bound) is the OR, the union of the ON-sets: "+ans+" rows.", nil)
		}
		r := f.and(g)
		ans := edakit.IntStr(r.ones())
		return edakit.Numeric("Let "+fs+" and "+sigma("g", g)+". In the lattice of Boolean functions ordered by implication, how many ON rows does the meet $f \\wedge g$ have?", ans,
			"The meet (greatest lower bound) is the AND, the intersection of the ON-sets: "+ans+" rows.", nil)
	}
}

func upClose(n int, gens []int) tt {
	f := tt{n, 0}
	for x := 0; x < rows(n); x++ {
		for _, g := range gens {
			if x&g == g {
				f.w |= 1 << uint(x)
				break
			}
		}
	}
	return f
}

// minimalTrue lists the ON points with no ON point strictly below them; the bit pattern x
// uses bit (n-1-v) for variable v, matching minterm numbering.
func minimalTrue(f tt) []int {
	var out []int
	for _, x := range f.minterms() {
		minimal := true
		for v := 0; v < f.n; v++ {
			b := 1 << uint(f.n-1-v)
			if x&b != 0 && f.at(x&^b) {
				minimal = false
				break
			}
		}
		if minimal {
			out = append(out, x)
		}
	}
	return out
}

func maximalFalse(f tt) []int {
	var out []int
	for x := 0; x < rows(f.n); x++ {
		if f.at(x) {
			continue
		}
		maximal := true
		for v := 0; v < f.n; v++ {
			b := 1 << uint(f.n-1-v)
			if x&b == 0 && !f.at(x|b) {
				maximal = false
				break
			}
		}
		if maximal {
			out = append(out, x)
		}
	}
	return out
}

func positiveTerm(x, n int) string {
	var b strings.Builder
	for v := 0; v < n; v++ {
		if bitOf(x, v, n) == 1 {
			b.WriteString(varNames[v])
		}
	}
	return b.String()
}

func monotone(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 5)
	var f tt
	for {
		var gens []int
		for i := 0; i < edakit.Int(2, 4); i++ {
			g := rand.Intn(rows(n))
			if bits.OnesCount(uint(g)) >= 1 && bits.OnesCount(uint(g)) <= n-1 {
				gens = append(gens, g)
			}
		}
		if len(gens) == 0 {
			continue
		}
		f = upClose(n, gens)
		if len(minimalTrue(f)) >= 2 {
			break
		}
	}
	mins := minimalTrue(f)
	switch kpOf(ctx) {
	case 1:
		var terms []string
		for _, m := range mins {
			terms = append(terms, positiveTerm(m, n))
		}
		sort.Strings(terms)
		ans := edakit.IntStr(f.ones())
		q := edakit.Sprintf("The monotone function $f(%s) = %s$ uses only positive literals. On how many of its %d rows is $f = 1$?", varList(n), strings.Join(terms, " + "), rows(n))
		return edakit.Numeric(q, ans, "The ON-set is the up-set of the minimal true points: every vector at or above one of them. It has "+ans+" rows.", nil)
	case 2:
		ans := edakit.IntStr(len(mins))
		var terms []string
		for _, m := range mins {
			terms = append(terms, positiveTerm(m, n))
		}
		q := "The function " + sigma("f", f) + " is monotone. How many minimal true points (ON rows with no ON row strictly below them) does it have?"
		return edakit.Numeric(q, ans, "The minimal true points are the vectors of the primes $"+strings.Join(terms, "$, $")+"$. For a monotone function each one is a positive prime implicant, so f has "+ans+" primes.", nil)
	default:
		mf := maximalFalse(f)
		ans := edakit.IntStr(len(mf))
		var vs []string
		for _, x := range mf {
			vs = append(vs, bitvec(x, n))
		}
		q := "The function " + sigma("f", f) + " is monotone. How many maximal false points (OFF rows with no OFF row strictly above them) does it have, which is the number of clauses in its minimal CNF?"
		return edakit.Numeric(q, ans, "The maximal false points are "+strings.Join(vs, ", ")+". Each gives the clause that ORs the variables at 0 in it.", nil)
	}
}

func fixpoint(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	for {
		if p, ok := fixpointOnce(ctx, kp); ok {
			return p
		}
	}
}

func fixpointOnce(ctx generator.GeneratorContext, kp int) (generator.Problem, bool) {
	ns := edakit.Scaled(ctx, 6, 9)
	succ := make([][]int, ns)
	var edges []string
	for s := 0; s < ns; s++ {
		k := edakit.Int(0, 2)
		seen := map[int]bool{}
		for i := 0; i < k; i++ {
			t := rand.Intn(ns)
			if seen[t] {
				continue
			}
			seen[t] = true
			succ[s] = append(succ[s], t)
			edges = append(edges, edakit.Sprintf("s_%d \\to s_%d", s, t))
		}
	}
	graph := "$" + strings.Join(edges, ",\\ ") + "$"
	if len(edges) == 0 {
		graph = "none"
	}
	switch kp {
	case 1, 2:
		init := map[int]bool{rand.Intn(ns): true}
		if rand.Intn(2) == 0 {
			init[rand.Intn(ns)] = true
		}
		reached := map[int]bool{}
		var frontier []int
		var initNames []string
		for s := 0; s < ns; s++ {
			if init[s] {
				reached[s] = true
				frontier = append(frontier, s)
				initNames = append(initNames, edakit.Sprintf("s_%d", s))
			}
		}
		steps := 0
		for len(frontier) > 0 {
			var nxt []int
			for _, s := range frontier {
				for _, t := range succ[s] {
					if !reached[t] {
						reached[t] = true
						nxt = append(nxt, t)
					}
				}
			}
			if len(nxt) > 0 {
				steps++
			}
			frontier = nxt
		}
		q0 := edakit.Sprintf("A transition system has states $s_0$ to $s_%d$ and transitions %s. Reachability iterates $R_0 = \\{%s\\}$, $R_{k+1} = R_k \\cup \\mathrm{Img}(R_k)$ until nothing changes.", ns-1, graph, strings.Join(initNames, ", "))
		if steps == 0 {
			return generator.Problem{}, false
		}
		if kp == 1 {
			ans := edakit.IntStr(steps)
			return edakit.Numeric(q0+" How many steps add at least one new state?", ans, "Breadth-first, each step adds the states one transition further out. The farthest reachable state is "+ans+" transitions from the initial set, so "+ans+" steps add states and the next step confirms the fixpoint.", nil), true
		}
		ans := edakit.IntStr(len(reached))
		return edakit.Numeric(q0+" How many states are in the least fixpoint, the reachable set?", ans, "The iteration is monotone on the finite subset lattice, so it stops at the least fixpoint containing the initial states: "+ans+" states.", nil), true
	default:
		P := map[int]bool{}
		var pn []string
		for s := 0; s < ns; s++ {
			if rand.Intn(3) != 0 {
				P[s] = true
				pn = append(pn, edakit.Sprintf("s_%d", s))
			}
		}
		if len(pn) == 0 {
			P[0] = true
			pn = []string{"s_0"}
		}
		Z := map[int]bool{}
		for s := range P {
			Z[s] = true
		}
		for {
			changed := false
			for s := 0; s < ns; s++ {
				if !Z[s] {
					continue
				}
				ok := false
				for _, t := range succ[s] {
					if Z[t] {
						ok = true
						break
					}
				}
				if !ok {
					delete(Z, s)
					changed = true
				}
			}
			if !changed {
				break
			}
		}
		if len(Z) == 0 && rand.Intn(4) != 0 {
			return generator.Problem{}, false
		}
		ans := edakit.IntStr(len(Z))
		q := edakit.Sprintf("A transition system has states $s_0$ to $s_%d$ and transitions %s. Let $P = \\{%s\\}$. Compute the greatest fixpoint of $Z = P \\cap \\mathrm{Pre}(Z)$, where $\\mathrm{Pre}(Z)$ is the set of states with at least one successor in $Z$, starting from $Z_0 = P$. How many states does it contain?",
			ns-1, graph, strings.Join(pn, ", "))
		return edakit.Numeric(q, ans, "Repeatedly remove states of Z with no successor in Z. What survives, "+ans+" states, is the set from which some infinite path stays inside P forever (the CTL property EG P).", nil), true
	}
}
