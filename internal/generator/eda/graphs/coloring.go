package graphs

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func randSimple(n int, p int) wgraph {
	g := newWGraph(n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if edakit.Int(1, 100) <= p {
				g.w[i][j], g.w[j][i] = 1, 1
			}
		}
	}
	return g
}

func greedyColor(g wgraph, order []int) []int {
	col := make([]int, g.n)
	for _, v := range order {
		used := map[int]bool{}
		for u := 0; u < g.n; u++ {
			if g.w[v][u] > 0 && col[u] > 0 {
				used[col[u]] = true
			}
		}
		c := 1
		for used[c] {
			c++
		}
		col[v] = c
	}
	return col
}

func chromatic(g wgraph) int {
	for k := 1; ; k++ {
		col := make([]int, g.n)
		var try func(v int) bool
		try = func(v int) bool {
			if v == g.n {
				return true
			}
			for c := 1; c <= k; c++ {
				ok := true
				for u := 0; u < v; u++ {
					if g.w[v][u] > 0 && col[u] == c {
						ok = false
						break
					}
				}
				if ok {
					col[v] = c
					if try(v + 1) {
						return true
					}
				}
			}
			return false
		}
		if try(0) {
			return k
		}
	}
}

func coloringGen(ctx generator.GeneratorContext) generator.Problem {
	for {
		n := edakit.Scaled(ctx, 6, 8)
		g := randSimple(n, 45)
		if len(edgeText(g, false)) == 0 {
			continue
		}
		order := edakit.Shuffle(seq(n))
		col := greedyColor(g, order)
		k := 0
		for _, c := range col {
			k = max(k, c)
		}
		var ot []string
		for _, v := range order {
			ot = append(ot, vName(v))
		}
		intro := "A register allocator builds an interference graph: vertices v1 to v" + edakit.IntStr(n) + " are values, and an edge joins two values that are live at the same time: " + edgeText(g, false) + ". Greedy coloring visits " + strings.Join(ot, ", ") + " and gives each vertex the smallest color number (1, 2, ...) not used by an already-colored neighbour. "
		var ct []string
		for v := 0; v < n; v++ {
			ct = append(ct, edakit.Sprintf("%s:%d", vName(v), col[v]))
		}
		expl := "Greedy colors: " + strings.Join(ct, ", ") + "."
		switch kpOf(ctx) {
		case 1:
			v := order[edakit.Int(n/2, n-1)]
			return edakit.Numeric(intro+"What color does "+vName(v)+" receive?", edakit.IntStr(col[v]), expl, map[string]string{"color": edakit.IntStr(col[v])})
		case 2:
			return edakit.Numeric(intro+"How many colors (registers) does the greedy pass use?", edakit.IntStr(k), expl, map[string]string{"colors": edakit.IntStr(k)})
		default:
			chi := chromatic(g)
			if chi == k && edakit.Int(0, 2) > 0 {
				continue
			}
			return edakit.Numeric(intro+edakit.Sprintf("The greedy pass uses %d colors. What is the chromatic number, the fewest colors any proper coloring needs?", k), edakit.IntStr(chi),
				expl+edakit.Sprintf(" An exhaustive search finds a proper coloring with %d colors and none with fewer. Greedy depends on the visiting order and is only an upper bound.", chi),
				map[string]string{"greedy": edakit.IntStr(k), "chromatic": edakit.IntStr(chi)})
		}
	}
}

type ivl struct{ l, r int }

func intervalGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 9)
	iv := make([]ivl, n)
	var it []string
	for i := range iv {
		l := edakit.Int(1, 16)
		iv[i] = ivl{l, l + edakit.Int(1, 7)}
		it = append(it, edakit.Sprintf("N%d [%d, %d]", i+1, iv[i].l, iv[i].r))
	}
	intro := "In a routing channel, each net occupies a closed column interval: " + strings.Join(it, ", ") + ". Two nets conflict when their intervals share a column, so they cannot use the same track. "
	switch kpOf(ctx) {
	case 1:
		dens := 0
		col := 0
		for c := 1; c <= 23; c++ {
			d := 0
			for _, x := range iv {
				if x.l <= c && c <= x.r {
					d++
				}
			}
			if d > dens {
				dens, col = d, c
			}
		}
		return edakit.Numeric(intro+"What is the minimum number of tracks, equal to the channel density (the largest number of intervals covering one column)?", edakit.IntStr(dens),
			edakit.Sprintf("Column %d is covered by %d intervals. Interval graphs are perfect, so the chromatic number equals this largest clique, and the left-edge algorithm achieves it.", col, dens),
			map[string]string{"density": edakit.IntStr(dens)})
	case 2:
		idx := seq(n)
		sort.SliceStable(idx, func(a, b int) bool { return iv[idx[a]].l < iv[idx[b]].l })
		var last []int
		track := make([]int, n)
		for _, i := range idx {
			placed := false
			for t := range last {
				if iv[i].l > last[t] {
					last[t] = iv[i].r
					track[i] = t + 1
					placed = true
					break
				}
			}
			if !placed {
				last = append(last, iv[i].r)
				track[i] = len(last)
			}
		}
		k := edakit.Int(0, n-1)
		var tt []string
		for i := range track {
			tt = append(tt, edakit.Sprintf("N%d:%d", i+1, track[i]))
		}
		return edakit.Numeric(intro+"The left-edge algorithm sorts nets by left end (ties keep the listed order) and puts each on the lowest-numbered track whose last net ends before the new net's left end, opening a new track if none fits. Which track does N"+edakit.IntStr(k+1)+" get?",
			edakit.IntStr(track[k]),
			"Track assignments: "+strings.Join(tt, ", ")+".",
			map[string]string{"track": edakit.IntStr(track[k])})
	default:
		e := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if iv[i].l <= iv[j].r && iv[j].l <= iv[i].r {
					e++
				}
			}
		}
		return edakit.Numeric(intro+"How many edges does the conflict (interval) graph have?", edakit.IntStr(e),
			"Two closed intervals overlap when each starts no later than the other ends. Counting overlapping pairs gives "+edakit.IntStr(e)+".",
			map[string]string{"edges": edakit.IntStr(e)})
	}
}

func components(g wgraph) int {
	seen := make([]bool, g.n)
	c := 0
	for s := 0; s < g.n; s++ {
		if seen[s] {
			continue
		}
		c++
		st := []int{s}
		seen[s] = true
		for len(st) > 0 {
			u := st[len(st)-1]
			st = st[:len(st)-1]
			for v := 0; v < g.n; v++ {
				if g.w[u][v] > 0 && !seen[v] {
					seen[v] = true
					st = append(st, v)
				}
			}
		}
	}
	return c
}

func minMono(g wgraph) int {
	best := 1 << 30
	for mask := 0; mask < 1<<g.n; mask++ {
		m := 0
		for i := 0; i < g.n; i++ {
			for j := i + 1; j < g.n; j++ {
				if g.w[i][j] > 0 && (mask>>i)&1 == (mask>>j)&1 {
					m++
				}
			}
		}
		best = min(best, m)
	}
	return best
}

func fName2(i int) string { return edakit.Sprintf("F%d", i+1) }

func dpEdgeText(g wgraph) string {
	var parts []string
	for i := 0; i < g.n; i++ {
		for j := i + 1; j < g.n; j++ {
			if g.w[i][j] > 0 {
				parts = append(parts, fName2(i)+"-"+fName2(j))
			}
		}
	}
	return strings.Join(parts, ", ")
}

func doublePatterningGen(ctx generator.GeneratorContext) generator.Problem {
	for {
		n := edakit.Scaled(ctx, 6, 9)
		var g wgraph
		bip := edakit.Int(0, 9) < 6
		if bip {
			g = newWGraph(n)
			side := randSides(n)
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					if side[i] != side[j] && edakit.Int(0, 9) < 3 {
						g.w[i][j], g.w[j][i] = 1, 1
					}
				}
			}
		} else {
			g = randSimple(n, 30)
		}
		if dpEdgeText(g) == "" {
			continue
		}
		comp := components(g)
		mono := minMono(g)
		intro := "On a double-patterning layer, features F1 to F" + edakit.IntStr(n) + " that are closer than the same-mask spacing get a conflict edge: " + dpEdgeText(g) + ". Each feature must go on mask A or mask B, and conflicting features need different masks. "
		switch kpOf(ctx) {
		case 1:
			return edakit.Numeric(intro+"How many connected components does the conflict graph have (an isolated feature counts as one)?", edakit.IntStr(comp),
				"Components can be colored independently, which is how decomposition tools split the problem. There are "+edakit.IntStr(comp)+".",
				map[string]string{"components": edakit.IntStr(comp)})
		case 2:
			cnt := 0
			if mono == 0 {
				cnt = 1 << comp
			}
			return edakit.Numeric(intro+"How many valid mask assignments are there in total (count assignments that differ in any feature as different)?", edakit.IntStr(cnt),
				edakit.Sprintf("A valid assignment exists only if the graph is bipartite (no odd cycle)%s. Each bipartite component then has exactly 2 assignments (swap A and B), so the count is 2^components when every component is bipartite.", map[bool]string{true: "; here it is, with " + edakit.IntStr(comp) + " components", false: "; here an odd cycle exists, so the count is 0"}[mono == 0]),
				map[string]string{"components": edakit.IntStr(comp), "assignments": edakit.IntStr(cnt)})
		default:
			if mono == 0 && edakit.Int(0, 2) > 0 {
				continue
			}
			return edakit.Numeric(intro+"What is the smallest number of conflict edges that must be left violated (both ends on the same mask) by the best assignment?", edakit.IntStr(mono),
				"Every odd cycle forces at least one violation. Searching all assignments, the minimum is "+edakit.IntStr(mono)+"; each such edge must be fixed by a stitch, a layout change or a third mask.",
				map[string]string{"violations": edakit.IntStr(mono)})
		}
	}
}

// planarityCheck is the computed part of planarity: what the Euler edge bound does and does not say.
func planarityCheck(ctx generator.GeneratorContext) generator.Problem {
	v := edakit.Int(5, 14)
	bip := edakit.Int(0, 2) == 0
	var e int
	bound := 3*v - 6
	if bip {
		bound = 2*v - 4
	}
	if edakit.Int(0, 1) == 0 {
		e = bound + edakit.Int(1, 4)
	} else {
		e = bound - edakit.Int(0, 4)
	}
	kind := "a conflict graph"
	bname := "3v - 6"
	if bip {
		kind = "a bipartite connection graph (no odd cycles, no triangles)"
		bname = "2v - 4"
	}
	q := edakit.Sprintf("A single-layer routing problem reduces to %s with v = %d vertices and e = %d edges. What does Euler's edge bound let you conclude?", kind, v, e)
	not := edakit.Sprintf("Not planar: e exceeds %s = %d", bname, bound)
	yes := edakit.Sprintf("Planar: e is within %s = %d", bname, bound)
	unk := edakit.Sprintf("Undecided: e is within %s = %d, which is necessary but not sufficient", bname, bound)
	other := edakit.Sprintf("Not planar: e exceeds v - 1 = %d, so the graph has a cycle", v-1)
	if e > bound {
		return edakit.Choice(q, not,
			edakit.Sprintf("A simple planar graph with v ≥ 3 has at most %s edges; %d > %d proves non-planarity, so a crossing (a via or second layer) is unavoidable.", bname, e, bound),
			yes, unk, other)
	}
	return edakit.Choice(q, unk,
		edakit.Sprintf("Every planar graph satisfies e ≤ %s, but the converse fails: K3,3 has 9 ≤ 3·6 - 6 = 12 edges and is non-planar, and subdividing the edges of K5 or K3,3 gives sparse non-planar graphs. Only a planarity test (Kuratowski's theorem, or a linear-time algorithm such as Hopcroft-Tarjan) decides.", bname),
		yes, not, other)
}
