package boolalg

import (
	"math/bits"
	"math/rand"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func satTo3sat(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	switch kp {
	case 1:
		k := edakit.Scaled(ctx, 4, 12)
		if rand.Intn(2) == 0 {
			ans := edakit.IntStr(k - 2)
			q := edakit.Sprintf("The standard SAT-to-3-SAT reduction replaces a clause of %d literals by a chain of three-literal clauses linked by fresh variables, preserving satisfiability. How many three-literal clauses replace it?", k)
			return edakit.Numeric(q, ans, edakit.Sprintf("(l1 or l2 or y1)(not y1 or l3 or y2)...(not y_{k-3} or l_{k-1} or l_k): the first and last clauses take two original literals each and every middle clause one, so k - 2 = %s clauses.", ans), nil)
		}
		ans := edakit.IntStr(k - 3)
		q := edakit.Sprintf("The standard SAT-to-3-SAT reduction splits a clause of %d literals into a chain of three-literal clauses. How many fresh linking variables does it introduce?", k)
		return edakit.Numeric(q, ans, edakit.Sprintf("A chain of k - 2 = %d clauses has k - 3 = %s links between consecutive clauses, one fresh variable each.", k-2, ans), nil)
	default:
		m := edakit.Scaled(ctx, 3, 7)
		sizes := make([]int, m)
		for i := range sizes {
			sizes[i] = edakit.Int(3, 9)
		}
		var ss []string
		newVars, clauses := 0, 0
		for _, s := range sizes {
			ss = append(ss, edakit.IntStr(s))
			newVars += s - 3
			clauses += s - 2
		}
		q0 := "A CNF has clauses with " + strings.Join(ss, ", ") + " literals. The standard reduction keeps the three-literal clauses and replaces each clause of k > 3 literals by a chain of k - 2 three-literal clauses with k - 3 fresh variables."
		if kp == 2 {
			ans := edakit.IntStr(newVars)
			return edakit.Numeric(q0+" How many fresh variables does the 3-CNF use?", ans, "Sum k - 3 over the clauses (zero for three-literal clauses): "+ans+".", nil)
		}
		ans := edakit.IntStr(clauses)
		return edakit.Numeric(q0+" How many clauses does the 3-CNF have?", ans, "Each clause of k >= 3 literals becomes k - 2 clauses (a three-literal clause stays as 1). The sum is "+ans+", linear in the original size.", nil)
	}
}

func graphReductions(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		V := edakit.Scaled(ctx, 6, 60)
		c := edakit.Int(1, V-1)
		switch rand.Intn(3) {
		case 0:
			ans := edakit.IntStr(V - c)
			q := edakit.Sprintf("A graph has %d vertices and its minimum vertex cover has %d vertices. How large is its maximum independent set?", V, c)
			return edakit.Numeric(q, ans, edakit.Sprintf("S is independent exactly when V minus S is a vertex cover, so the maximum independent set has %d - %d = %s vertices.", V, c, ans), nil)
		case 1:
			ans := edakit.IntStr(V - c)
			q := edakit.Sprintf("A graph $G$ has %d vertices and its maximum independent set has %d vertices. How large is the minimum vertex cover of $G$?", V, c)
			return edakit.Numeric(q, ans, edakit.Sprintf("Complement the largest independent set to get the smallest cover: %d - %d = %s.", V, c, ans), nil)
		default:
			ans := edakit.IntStr(V - c)
			q := edakit.Sprintf("A graph $G$ has %d vertices and its minimum vertex cover has %d vertices. What is the largest clique in the complement graph $\\overline{G}$?", V, c)
			return edakit.Numeric(q, ans, edakit.Sprintf("Cliques of the complement graph are independent sets of G, and the largest independent set of G is %d - %d = %s.", V, c, ans), nil)
		}
	case 2:
		n := edakit.Scaled(ctx, 3, 30)
		m := edakit.Scaled(ctx, 2, 60)
		q0 := edakit.Sprintf("The Garey-Johnson reduction from 3-SAT to vertex cover builds, for a formula with %d variables and %d clauses, one edge per variable (its two literals), one triangle per clause, and an edge from each triangle corner to the matching literal vertex.", n, m)
		switch rand.Intn(3) {
		case 0:
			ans := edakit.IntStr(2*n + 3*m)
			return edakit.Numeric(q0+" How many vertices does the graph have?", ans, edakit.Sprintf("2 per variable plus 3 per clause: 2(%d) + 3(%d) = %s.", n, m, ans), nil)
		case 1:
			ans := edakit.IntStr(n + 6*m)
			return edakit.Numeric(q0+" How many edges does the graph have?", ans, edakit.Sprintf("1 per variable, 3 triangle edges and 3 links per clause: %d + 6(%d) = %s.", n, m, ans), nil)
		default:
			ans := edakit.IntStr(n + 2*m)
			return edakit.Numeric(q0+" What cover size $k$ makes the instance equivalent to the formula being satisfiable?", ans, edakit.Sprintf("A cover needs one vertex per variable edge and two per triangle: k = %d + 2(%d) = %s.", n, m, ans), nil)
		}
	default:
		V := edakit.Scaled(ctx, 5, 40)
		maxE := V * (V - 1) / 2
		E := edakit.Int(V-1, maxE-1)
		ans := edakit.IntStr(maxE - E)
		q := edakit.Sprintf("A graph $G$ has %d vertices and %d edges. To reduce clique on $G$ to independent set, a tool builds the complement graph. How many edges does the complement have?", V, E)
		return edakit.Numeric(q, ans, edakit.Sprintf("The complete graph has C(%d, 2) = %d edges, so the complement has %d - %d = %s. A clique of G is an independent set of the complement, so the reduction is polynomial and preserves the answer.", V, maxE, maxE, E, ans), nil)
	}
}

type edge struct{ u, v int }

func randGraph(V, E int) []edge {
	var all []edge
	for u := 1; u <= V; u++ {
		for v := u + 1; v <= V; v++ {
			all = append(all, edge{u, v})
		}
	}
	all = edakit.Shuffle(all)
	if E > len(all) {
		E = len(all)
	}
	return all[:E]
}

func edgeList(es []edge) string {
	parts := make([]string, len(es))
	for i, e := range es {
		parts[i] = edakit.Sprintf("(%d,%d)", e.u, e.v)
	}
	return strings.Join(parts, ", ")
}

func minVC(V int, es []edge) int {
	best := V
	for s := 0; s < 1<<uint(V); s++ {
		k := bits.OnesCount(uint(s))
		if k >= best {
			continue
		}
		ok := true
		for _, e := range es {
			if s>>uint(e.u-1)&1 == 0 && s>>uint(e.v-1)&1 == 0 {
				ok = false
				break
			}
		}
		if ok {
			best = k
		}
	}
	return best
}

func matchingCover(es []edge) int {
	in := map[int]bool{}
	for _, e := range es {
		if !in[e.u] && !in[e.v] {
			in[e.u], in[e.v] = true, true
		}
	}
	return len(in)
}

func harmonic(d int) float64 {
	h := 0.0
	for i := 1; i <= d; i++ {
		h += 1 / float64(i)
	}
	return h
}

func approximation(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	switch kp {
	case 1, 2:
		V := edakit.Scaled(ctx, 5, 8)
		es := randGraph(V, edakit.Int(V, V+4))
		c := matchingCover(es)
		q0 := edakit.Sprintf("A graph on vertices 1 to %d has edges %s. The 2-approximation for vertex cover scans the edges in this order and, whenever an edge has neither endpoint in the cover, adds both endpoints.", V, edgeList(es))
		if kp == 1 {
			ans := edakit.IntStr(c)
			return edakit.Numeric(q0+" How many vertices does it return?", ans, edakit.Sprintf("The edges it picks form a maximal matching of %d edges, and it returns both endpoints of each: %s vertices.", c/2, ans), nil)
		}
		opt := minVC(V, es)
		ans := edakit.Num(float64(c)/float64(opt), 2)
		return edakit.Numeric(q0+" The optimum cover is found by exhaustive search. What ratio does the algorithm achieve on this instance, its size divided by the optimum? Round to 2 decimal places.", ans, edakit.Sprintf("The algorithm returns %d vertices and the optimum is %d, a ratio of %s. The proof guarantees at most 2, because any cover needs a distinct vertex for each of the %d matched edges.", c, opt, ans, c/2), map[string]string{"opt": edakit.IntStr(opt)})
	default:
		d := edakit.Scaled(ctx, 2, 12)
		if rand.Intn(2) == 0 {
			ans := edakit.Num(harmonic(d), 3)
			q := edakit.Sprintf("Greedy set cover repeatedly takes the set that covers the most uncovered elements. If the largest set has %d elements, greedy is within a factor $H(%d) = 1 + 1/2 + \\dots + 1/%d$ of optimal. What is that factor? Round to 3 decimal places.", d, d, d)
			return edakit.Numeric(q, ans, "Summing the reciprocals gives "+ans+". The bound grows like ln d, and no polynomial algorithm can do better than about ln n unless P = NP.", nil)
		}
		k := edakit.Int(2, 12)
		bound := int(float64(k) * harmonic(d))
		ans := edakit.IntStr(bound)
		q := edakit.Sprintf("In a test-set compaction problem posed as set cover, the largest set covers %d faults and the optimum uses %d tests. What is the largest number of tests that greedy set cover can return, by the $H(d)$ guarantee? (Take the floor.)", d, k)
		return edakit.Numeric(q, ans, edakit.Sprintf("Greedy uses at most H(%d) x %d = %s tests, which floors to %s. Heuristics such as simulated annealing carry no such guarantee, even when they do better in practice.", d, k, edakit.Num(float64(k)*harmonic(d), 3), ans), nil)
	}
}

func parameterized(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	switch kp {
	case 1:
		k := edakit.Scaled(ctx, 3, 20)
		if rand.Intn(2) == 0 {
			ans := edakit.IntStr(pow2(k))
			q := edakit.Sprintf("The bounded search tree for vertex cover picks an uncovered edge $(u, v)$ and branches: put $u$ in the cover or put $v$ in the cover, each lowering the budget by 1. With budget $k = %d$, at most how many leaves does the tree have?", k)
			return edakit.Numeric(q, ans, edakit.Sprintf("Two branches per level and depth at most k: 2^%d = %s leaves, independent of the graph size n. The running time O(2^k n) is fixed-parameter tractable.", k, ans), nil)
		}
		L := []int{1, 1, 1}
		for i := 3; i <= k; i++ {
			L = append(L, L[i-1]+L[i-3])
		}
		ans := edakit.IntStr(L[k])
		q := edakit.Sprintf("A refined vertex-cover search branches on a vertex $v$ of degree at least 3: take $v$ (budget drops by 1) or take all its neighbours (budget drops by at least 3). Its leaf count satisfies $L(k) = L(k-1) + L(k-3)$ with $L(0) = L(1) = L(2) = 1$. What is $L(%d)$?", k)
		return edakit.Numeric(q, ans, edakit.Sprintf("Iterating gives L(%d) = %s, against 2^%d = %d for plain edge branching. Better branching lowers the base of the exponent, not the dependence on n.", k, ans, k, pow2(k)), nil)
	default:
		for {
			V := edakit.Scaled(ctx, 8, 12)
			k := edakit.Int(2, 4)
			// Hubs with many edges plus a sparse remainder.
			deg := map[int]map[int]bool{}
			add := func(u, v int) {
				if u == v {
					return
				}
				if deg[u] == nil {
					deg[u] = map[int]bool{}
				}
				if deg[v] == nil {
					deg[v] = map[int]bool{}
				}
				deg[u][v], deg[v][u] = true, true
			}
			hubs := edakit.Int(1, k)
			for h := 1; h <= hubs; h++ {
				for _, v := range rand.Perm(V)[:edakit.Int(k+1, V-1)] {
					add(h, v+1)
				}
			}
			for i := 0; i < edakit.Int(1, 4); i++ {
				add(edakit.Int(hubs+1, V), edakit.Int(hubs+1, V))
			}
			var es []edge
			for u := 1; u <= V; u++ {
				for v := range deg[u] {
					if u < v {
						es = append(es, edge{u, v})
					}
				}
			}
			sort.Slice(es, func(i, j int) bool {
				if es[i].u != es[j].u {
					return es[i].u < es[j].u
				}
				return es[i].v < es[j].v
			})
			// Apply Buss's rule to a fixpoint.
			alive := map[edge]bool{}
			for _, e := range es {
				alive[e] = true
			}
			budget := k
			forced := 0
			for {
				d := map[int]int{}
				for e := range alive {
					d[e.u]++
					d[e.v]++
				}
				pick := -1
				for v := 1; v <= V; v++ {
					if d[v] > budget {
						pick = v
						break
					}
				}
				if pick < 0 || budget == 0 {
					break
				}
				for e := range alive {
					if e.u == pick || e.v == pick {
						delete(alive, e)
					}
				}
				budget--
				forced++
			}
			if forced == 0 || (kp == 3 && budget == 0) {
				continue
			}
			q0 := edakit.Sprintf("A graph on vertices 1 to %d has edges %s. We ask for a vertex cover of size at most $k = %d$. Buss's kernel rule: any vertex of degree greater than the current budget must be in the cover, so take it, delete its edges and lower the budget by 1; repeat.", V, edgeList(es), k)
			if kp == 2 {
				ans := edakit.IntStr(forced)
				return edakit.Numeric(q0+" How many vertices does the rule force into the cover?", ans, edakit.Sprintf("Recompute degrees after each forced vertex. %s vertices are forced, leaving budget %d.", ans, budget), nil)
			}
			ans := edakit.IntStr(budget * budget)
			verdict := "within the bound, so the small kernel goes on to exact search"
			if len(alive) > budget*budget {
				verdict = "more than the bound, so the instance is rejected at once"
			}
			return edakit.Numeric(q0+" After the rule stops, every remaining vertex has degree at most the remaining budget $k'$, so a yes-instance has at most $k'^2$ edges left. What is that edge bound?", ans, edakit.Sprintf("%d forced vertices give k' = %d - %d = %d, so a yes-instance keeps at most %d^2 = %s edges. Here %d edges remain, %s.", forced, k, forced, budget, budget, ans, len(alive), verdict), map[string]string{"forced": edakit.IntStr(forced)})
		}
	}
}

// reductionCompute serves KP 2 of the reduction concept: the size of the 3-SAT to
// independent-set graph for a drawn formula.
func reductionCompute(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 5)
	m := edakit.Scaled(ctx, 3, 5)
	var cs []clause
	for i := 0; i < m; i++ {
		cs = append(cs, randClause(n, 3))
	}
	pairs := 0
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			for _, a := range cs[i] {
				for _, b := range cs[j] {
					if a == -b {
						pairs++
					}
				}
			}
		}
	}
	V, E, k := 3*m, 3*m+pairs, m
	f := func(v, e, kk int) string { return edakit.Sprintf("$|V| = %d$, $|E| = %d$, $k = %d$", v, e, kk) }
	ans := f(V, E, k)
	ds := []string{f(V, pairs, k), f(V, E, n), f(2*n+3*m, n+6*m, n+2*m), f(V, 3*m+2*pairs, k)}
	q := "Karp's reduction from 3-SAT to independent set makes one vertex per literal occurrence, joins the three literals of each clause in a triangle, joins every pair of complementary literals in different clauses, and asks for an independent set of size $k$. Apply it to " + cnfTeX(cs) + ". What does it produce?"
	return edakit.Choice(q, ans, edakit.Sprintf("Vertices: 3 x %d = %d. Edges: %d triangle edges plus %d complementary pairs = %d. The target is one literal per clause, k = %d.", m, V, 3*m, pairs, E, k), ds...)
}
