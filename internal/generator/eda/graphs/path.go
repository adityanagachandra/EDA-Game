package graphs

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func gridGraphGen(ctx generator.GeneratorContext) generator.Problem {
	W := edakit.Scaled(ctx, 4, 40)
	H := edakit.Scaled(ctx, 4, 40)
	L := edakit.Int(2, 6)
	hl := (L + 1) / 2 // layers 1, 3, 5 are horizontal
	vl := L / 2
	nodes := W * H * L
	wire := hl*(W-1)*H + vl*W*(H-1)
	vias := W * H * (L - 1)
	intro := edakit.Sprintf("A detailed-routing grid has %d columns, %d rows and %d metal layers. Layer 1 runs horizontally and directions alternate layer by layer. Each layer has one vertex per grid point, wire edges join neighbouring points along that layer's preferred direction only, and via edges join the same point on adjacent layers. ", W, H, L)
	switch kpOf(ctx) {
	case 1:
		return edakit.Numeric(intro+"How many vertices does the grid graph have?", edakit.IntStr(nodes),
			edakit.Sprintf("One vertex per point per layer: %d × %d × %d = %d.", W, H, L, nodes),
			map[string]string{"vertices": edakit.IntStr(nodes)})
	case 2:
		return edakit.Numeric(intro+"How many wire edges (excluding vias) does it have?", edakit.IntStr(wire),
			edakit.Sprintf("Horizontal layers (%d of them) each have (%d - 1) × %d edges and vertical layers (%d of them) each have %d × (%d - 1) edges: %d in total.", hl, W, H, vl, W, H, wire),
			map[string]string{"wire_edges": edakit.IntStr(wire)})
	default:
		return edakit.Numeric(intro+"How many edges does it have in total, wires plus vias?", edakit.IntStr(wire+vias),
			edakit.Sprintf("Wire edges: %d. Via edges: %d × %d × (%d - 1) = %d. Total %d.", wire, W, H, L, vias, wire+vias),
			map[string]string{"wire_edges": edakit.IntStr(wire), "vias": edakit.IntStr(vias)})
	}
}

// maze is a grid with blocked cells, indexed (col, row) from (1,1).
type maze struct {
	w, h    int
	blocked map[pt]bool
}

func (m maze) bfs(s pt) map[pt]int {
	dist := map[pt]int{s: 0}
	q := []pt{s}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, d := range []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			v := pt{u.x + d.x, u.y + d.y}
			if v.x < 1 || v.y < 1 || v.x > m.w || v.y > m.h || m.blocked[v] {
				continue
			}
			if _, ok := dist[v]; !ok {
				dist[v] = dist[u] + 1
				q = append(q, v)
			}
		}
	}
	return dist
}

func randMaze(ctx generator.GeneratorContext) (maze, pt, pt, map[pt]int) {
	for {
		w, h := edakit.Scaled(ctx, 5, 7), edakit.Scaled(ctx, 5, 7)
		m := maze{w: w, h: h, blocked: map[pt]bool{}}
		s := pt{1, edakit.Int(1, h)}
		t := pt{w, edakit.Int(1, h)}
		nb := edakit.Scaled(ctx, 4, 9)
		// A wall segment forces detours.
		wx := edakit.Int(2, w-1)
		gap := edakit.Int(1, h)
		for y := 1; y <= h; y++ {
			if y != gap && len(m.blocked) < nb {
				m.blocked[pt{wx, y}] = true
			}
		}
		for len(m.blocked) < nb {
			p := pt{edakit.Int(2, w-1), edakit.Int(1, h)}
			m.blocked[p] = true
		}
		dist := m.bfs(s)
		if d, ok := dist[t]; ok && d > manh(s, t) {
			return m, s, t, dist
		}
	}
}

func blockedText(m maze) string {
	var ps []pt
	for p := range m.blocked {
		ps = append(ps, p)
	}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].x != ps[j].x {
			return ps[i].x < ps[j].x
		}
		return ps[i].y < ps[j].y
	})
	var parts []string
	for _, p := range ps {
		parts = append(parts, edakit.Sprintf("(%d, %d)", p.x, p.y))
	}
	return strings.Join(parts, ", ")
}

func bfsMazeGen(ctx generator.GeneratorContext) generator.Problem {
	m, s, t, dist := randMaze(ctx)
	intro := edakit.Sprintf("A routing grid has columns 1 to %d and rows 1 to %d; cells are written (column, row). Blocked cells: %s. The source is at (%d, %d) and the target at (%d, %d). Moves go to the 4 neighbours. ", m.w, m.h, blockedText(m), s.x, s.y, t.x, t.y)
	switch kpOf(ctx) {
	case 1:
		var cands []pt
		for p, d := range dist {
			if d >= 3 && d <= 7 {
				cands = append(cands, p)
			}
		}
		sort.Slice(cands, func(i, j int) bool {
			if cands[i].x != cands[j].x {
				return cands[i].x < cands[j].x
			}
			return cands[i].y < cands[j].y
		})
		p := edakit.Pick(cands)
		return edakit.Numeric(intro+edakit.Sprintf("Lee's algorithm labels every reachable cell with its wavefront number, the source being 0. What label does cell (%d, %d) receive?", p.x, p.y),
			edakit.IntStr(dist[p]),
			edakit.Sprintf("The wavefront expands one step at a time around blocked cells, so the label is the breadth-first distance: %d.", dist[p]),
			map[string]string{"label": edakit.IntStr(dist[p])})
	case 2:
		return edakit.Numeric(intro+"What is the length, in grid steps, of the shortest route from source to target?", edakit.IntStr(dist[t]),
			edakit.Sprintf("Without blockages the distance would be %d. The blocked cells force a detour, and breadth-first search finds length %d.", manh(s, t), dist[t]),
			map[string]string{"length": edakit.IntStr(dist[t]), "manhattan": edakit.IntStr(manh(s, t))})
	default:
		k := edakit.Int(2, 5)
		cnt := 0
		for _, d := range dist {
			if d == k {
				cnt++
			}
		}
		return edakit.Numeric(intro+edakit.Sprintf("How many cells are on wavefront %d (cells whose label is exactly %d)?", k, k), edakit.IntStr(cnt),
			edakit.Sprintf("Wavefront %d holds every unblocked cell at breadth-first distance %d from the source: %d cells.", k, k, cnt),
			map[string]string{"wavefront": edakit.IntStr(cnt)})
	}
}

// dijkstraTrace runs Dijkstra from s on g, breaking ties by the smaller vertex index. It returns
// the final distances, the settling order and the tentative distances after each settle.
func dijkstraTrace(g wgraph, s int) ([]int, []int, [][]int) {
	const inf = math.MaxInt32
	d := make([]int, g.n)
	for i := range d {
		d[i] = inf
	}
	d[s] = 0
	done := make([]bool, g.n)
	var order []int
	var snaps [][]int
	for k := 0; k < g.n; k++ {
		u := -1
		for i := 0; i < g.n; i++ {
			if !done[i] && d[i] < inf && (u < 0 || d[i] < d[u]) {
				u = i
			}
		}
		if u < 0 {
			break
		}
		done[u] = true
		order = append(order, u)
		for v := 0; v < g.n; v++ {
			if g.w[u][v] > 0 && !done[v] && d[u]+g.w[u][v] < d[v] {
				d[v] = d[u] + g.w[u][v]
			}
		}
		snaps = append(snaps, append([]int(nil), d...))
	}
	return d, order, snaps
}

func dijkstraGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	g := randWGraph(n, edakit.Scaled(ctx, 2, 6), 9)
	d, order, snaps := dijkstraTrace(g, 0)
	intro := "An undirected routing graph has vertices v1 to v" + edakit.IntStr(n) + " and edge costs " + edgeText(g, true) + ". Run Dijkstra's algorithm from v1, breaking ties by the smaller vertex number. "
	var on []string
	for _, v := range order {
		on = append(on, edakit.Sprintf("%s (%d)", vName(v), d[v]))
	}
	expl := "Settling order with final distances: " + strings.Join(on, ", ") + "."
	switch kpOf(ctx) {
	case 1:
		t := edakit.Int(2, n-1)
		return edakit.Numeric(intro+"What is the shortest-path cost from v1 to "+vName(t)+"?", edakit.IntStr(d[t]), expl, map[string]string{"distance": edakit.IntStr(d[t])})
	case 2:
		k := edakit.Int(3, n)
		return edakit.Numeric(intro+edakit.Sprintf("Which vertex is settled (removed from the priority queue) in position %d, counting v1 as position 1? Answer with its number.", k),
			edakit.IntStr(order[k-1]+1), expl, map[string]string{"order": edakit.IntStr(order[k-1] + 1)})
	default:
		for tries := 0; tries < 50; tries++ {
			k := edakit.Int(1, min(3, n-2))
			x := edakit.Int(1, n-1)
			settledBy := false
			for _, v := range order[:k] {
				settledBy = settledBy || v == x
			}
			if settledBy || snaps[k-1][x] == math.MaxInt32 {
				continue
			}
			return edakit.Numeric(intro+edakit.Sprintf("Right after the first %d vertices have been settled and their edges relaxed, what is the tentative distance of %s?", k, vName(x)),
				edakit.IntStr(snaps[k-1][x]),
				expl+edakit.Sprintf(" After %d settles the tentative distance of %s is %d; it may still drop later (its final value is %d).", k, vName(x), snaps[k-1][x], d[x]),
				map[string]string{"tentative": edakit.IntStr(snaps[k-1][x])})
		}
		return edakit.Numeric(intro+"What is the shortest-path cost from v1 to "+vName(n-1)+"?", edakit.IntStr(d[n-1]), expl, map[string]string{"distance": edakit.IntStr(d[n-1])})
	}
}

type darc struct{ u, v, w int }

func arcsText(arcs []darc) string {
	var parts []string
	for _, a := range arcs {
		parts = append(parts, edakit.Sprintf("%s→%s (%d)", vName(a.u), vName(a.v), a.w))
	}
	return strings.Join(parts, ", ")
}

// bellmanFord returns distances from s and the distances after the first round.
func bellmanFord(n int, arcs []darc, s int) ([]int, []int) {
	const inf = math.MaxInt32
	d := make([]int, n)
	for i := range d {
		d[i] = inf
	}
	d[s] = 0
	var first []int
	for r := 0; r < n-1; r++ {
		for _, a := range arcs {
			if d[a.u] < inf && d[a.u]+a.w < d[a.v] {
				d[a.v] = d[a.u] + a.w
			}
		}
		if r == 0 {
			first = append([]int(nil), d...)
		}
	}
	return d, first
}

// randNegArcs draws a DAG-like arc set with some negative weights (no cycles, so no negative
// cycles), shuffled so that the listed order is not topological.
func randNegArcs(n int, extra int) []darc {
	var arcs []darc
	has := map[[2]int]bool{}
	for v := 1; v < n; v++ {
		u := edakit.Int(max(0, v-2), v-1)
		has[[2]int{u, v}] = true
		arcs = append(arcs, darc{u, v, edakit.Int(1, 9)})
	}
	for t := 0; t < extra; t++ {
		u := edakit.Int(0, n-2)
		v := edakit.Int(u+1, n-1)
		if has[[2]int{u, v}] {
			continue
		}
		has[[2]int{u, v}] = true
		w := edakit.Int(1, 9)
		if edakit.Int(0, 1) == 0 {
			w = -edakit.Int(1, 6)
		}
		arcs = append(arcs, darc{u, v, w})
	}
	return edakit.Shuffle(arcs)
}

func bellmanFordGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1, 2:
		n := edakit.Scaled(ctx, 5, 7)
		for {
			arcs := randNegArcs(n, edakit.Scaled(ctx, 3, 6))
			neg := false
			for _, a := range arcs {
				neg = neg || a.w < 0
			}
			if !neg {
				continue
			}
			d, first := bellmanFord(n, arcs, 0)
			t := n - 1
			intro := "A directed constraint graph has vertices v1 to v" + edakit.IntStr(n) + " and arcs, listed in this order, " + arcsText(arcs) + ". There is no negative cycle. "
			if ctx.KPOrdinal == 2 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
				if first[t] == math.MaxInt32 || first[t] == d[t] {
					continue
				}
				return edakit.Numeric(intro+"Bellman-Ford from v1 relaxes every arc once per round, in the listed order. What is the distance estimate of "+vName(t)+" at the end of round 1?",
					edakit.IntStr(first[t]),
					"One round in this order gives estimates "+intsText(first)+" (unreached vertices stay infinite). More rounds lower "+vName(t)+" to its final value "+edakit.IntStr(d[t])+".",
					map[string]string{"round1": edakit.IntStr(first[t]), "final": edakit.IntStr(d[t])})
			}
			return edakit.Numeric(intro+"What is the shortest-path distance from v1 to "+vName(t)+"?", edakit.IntStr(d[t]),
				"After n - 1 rounds of relaxing every arc the distances are final: "+intsText(d)+".",
				map[string]string{"distance": edakit.IntStr(d[t])})
		}
	default:
		k := edakit.Int(3, 5)
		var ws []int
		sum := 0
		for i := 0; i < k-1; i++ {
			w := edakit.Int(-6, 9)
			ws = append(ws, w)
			sum += w
		}
		var names []string
		for i := 0; i < k; i++ {
			names = append(names, vName(i))
		}
		var arcs []string
		for i := 0; i < k-1; i++ {
			arcs = append(arcs, edakit.Sprintf("%s→%s (%d)", vName(i), vName(i+1), ws[i]))
		}
		ans := -sum
		return edakit.Numeric("A cycle in a constraint graph has arcs "+strings.Join(arcs, ", ")+" and a closing arc "+vName(k-1)+"→v1 with weight x. What is the smallest integer x for which the cycle is not negative, so Bellman-Ford converges?",
			edakit.IntStr(ans),
			edakit.Sprintf("The cycle weight is %d + x. It is non-negative when x ≥ %d.", sum, ans),
			map[string]string{"partial_sum": edakit.IntStr(sum)})
	}
}

func differenceConstraintsGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 5)
		for {
			type con struct{ i, j, c int }
			var cons []con
			var arcs []darc
			var text []string
			m := min(edakit.Int(n, n+2), n*(n-1)/2)
			has := map[[2]int]bool{}
			for len(cons) < m {
				i, j := edakit.Int(1, n), edakit.Int(1, n)
				if i == j || has[[2]int{i, j}] || has[[2]int{j, i}] {
					continue
				}
				has[[2]int{i, j}] = true
				c := edakit.Int(-4, 6)
				cons = append(cons, con{i, j, c})
				arcs = append(arcs, darc{i, j, c})
				text = append(text, edakit.Sprintf("$x_%d - x_%d \\le %d$", j, i, c))
			}
			for v := 1; v <= n; v++ {
				arcs = append(arcs, darc{0, v, 0})
			}
			d, _ := bellmanFord(n+1, arcs, 0)
			// Check feasibility: one more round must not improve.
			ok := true
			for _, a := range arcs {
				if d[a.u]+a.w < d[a.v] {
					ok = false
				}
			}
			if !ok {
				continue
			}
			v := edakit.Int(1, n)
			if d[v] == 0 && edakit.Int(0, 3) > 0 {
				continue
			}
			return edakit.Numeric("Solve the difference constraints "+strings.Join(text, ", ")+" with the standard construction: an arc $x_i \\to x_j$ of weight $c$ for each constraint $x_j - x_i \\le c$, plus a source $x_0$ with weight-0 arcs to every variable, and $x_k$ = shortest distance from $x_0$. "+edakit.Sprintf("What value does $x_%d$ get?", v),
				edakit.IntStr(d[v]),
				"Shortest distances satisfy every constraint because d(j) ≤ d(i) + c is exactly x_j - x_i ≤ c. Distances from the source: "+intsText(d[1:])+".",
				map[string]string{"x": edakit.IntStr(d[v])})
		}
	case 2:
		n := edakit.Scaled(ctx, 4, 7)
		w := make([]int, n)
		for i := range w {
			w[i] = edakit.Int(1, 5)
		}
		type sp struct{ a, b, s int }
		var cons []sp
		var text []string
		for i := 0; i < n; i++ {
			text = append(text, edakit.Sprintf("cell %d has width %d", i+1, w[i]))
		}
		// Spacing constraints: left edge of b at least left edge of a plus width of a plus gap.
		has := map[[2]int]bool{}
		for i := 1; i < n; i++ {
			a := edakit.Int(0, i-1)
			has[[2]int{a, i}] = true
			cons = append(cons, sp{a, i, edakit.Int(0, 3)})
		}
		for t := 0; t < n-2; t++ {
			b := edakit.Int(1, n-1)
			a := edakit.Int(0, b-1)
			if has[[2]int{a, b}] {
				continue
			}
			has[[2]int{a, b}] = true
			cons = append(cons, sp{a, b, edakit.Int(0, 3)})
		}
		var ctext []string
		for _, c := range cons {
			ctext = append(ctext, edakit.Sprintf("cell %d left of cell %d with gap at least %d", c.a+1, c.b+1, c.s))
		}
		x := make([]int, n)
		for b := 0; b < n; b++ {
			for _, c := range cons {
				if c.b == b {
					x[b] = max(x[b], x[c.a]+w[c.a]+c.s)
				}
			}
		}
		width := 0
		for i := 0; i < n; i++ {
			width = max(width, x[i]+w[i])
		}
		return edakit.Numeric("A one-dimensional compactor places cells on a row whose left edge is at 0: "+strings.Join(text, ", ")+". Constraints: "+strings.Join(ctext, "; ")+". Each constraint $x_b \\ge x_a + w_a + g$ is an arc of the constraint graph, and the smallest legal positions are longest-path distances from the left edge. What is the minimum row width (the largest right edge)?",
			edakit.IntStr(width),
			"Longest-path positions of the left edges: "+intsText(x)+". Adding each cell's width, the largest right edge is "+edakit.IntStr(width)+".",
			map[string]string{"width": edakit.IntStr(width)})
	default:
		n := edakit.Scaled(ctx, 4, 6)
		for {
			var arcs []darc
			var text []string
			has := map[[2]int]bool{}
			for len(arcs) < n+2 {
				i, j := edakit.Int(1, n), edakit.Int(1, n)
				if i == j || has[[2]int{i, j}] {
					continue
				}
				has[[2]int{i, j}] = true
				c := edakit.Int(-2, 7)
				arcs = append(arcs, darc{i - 1, j - 1, c})
				text = append(text, edakit.Sprintf("$x_%d - x_%d \\le %d$", j, i, c))
			}
			// Feasibility of the whole system: no negative cycle anywhere (virtual source n).
			all := append([]darc(nil), arcs...)
			for v := 0; v < n; v++ {
				all = append(all, darc{n, v, 0})
			}
			pd, _ := bellmanFord(n+1, all, n)
			ok := true
			for _, a := range all {
				if pd[a.u]+a.w < pd[a.v] {
					ok = false
				}
			}
			d, _ := bellmanFord(n, arcs, 0)
			if !ok || d[n-1] == math.MaxInt32 {
				continue
			}
			return edakit.Numeric("Constraints: "+strings.Join(text, ", ")+edakit.Sprintf(". What is the largest value that $x_%d - x_1$ can take?", n),
				edakit.IntStr(d[n-1]),
				edakit.Sprintf("Chaining constraints along any path from x_1 to x_%d bounds the difference by that path's weight, so the tightest bound is the shortest-path distance, %d, and the potentials d(v) achieve it.", n, d[n-1]),
				map[string]string{"bound": edakit.IntStr(d[n-1])})
		}
	}
}

func astarGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		t := pt{edakit.Int(6, 12), edakit.Int(6, 12)}
		p := pt{edakit.Int(1, 10), edakit.Int(1, 10)}
		for p == t {
			p = pt{edakit.Int(1, 10), edakit.Int(1, 10)}
		}
		g := manh(pt{1, 1}, p) + 2*edakit.Int(0, 3)
		h := manh(p, t)
		return edakit.Numeric(edakit.Sprintf("A* routes from (1, 1) to the target (%d, %d) on a grid with unit step cost, using the Manhattan-distance heuristic. Node (%d, %d) is in the open list with path cost g = %d. What is its priority f = g + h?", t.x, t.y, p.x, p.y, g),
			edakit.IntStr(g+h),
			edakit.Sprintf("h = |%d - %d| + |%d - %d| = %d, so f = %d + %d = %d.", p.x, t.x, p.y, t.y, h, g, h, g+h),
			map[string]string{"h": edakit.IntStr(h), "f": edakit.IntStr(g + h)})
	case 2:
		t := pt{edakit.Int(8, 14), edakit.Int(8, 14)}
		for {
			type cand struct {
				p    pt
				g, f int
			}
			var cs []cand
			var text []string
			for i := 0; i < 4; i++ {
				p := pt{edakit.Int(1, t.x), edakit.Int(1, t.y)}
				g := manh(pt{1, 1}, p) + 2*edakit.Int(0, 3)
				cs = append(cs, cand{p, g, g + manh(p, t)})
				text = append(text, edakit.Sprintf("node %d at (%d, %d) with g = %d", i+1, p.x, p.y, g))
			}
			best := 0
			unique := true
			for i := 1; i < 4; i++ {
				if cs[i].f < cs[best].f {
					best = i
				}
			}
			for i := 0; i < 4; i++ {
				if i != best && cs[i].f == cs[best].f {
					unique = false
				}
			}
			if !unique {
				continue
			}
			var ft []string
			for i, c := range cs {
				ft = append(ft, edakit.Sprintf("node %d: f = %d", i+1, c.f))
			}
			return edakit.Numeric(edakit.Sprintf("A* routes toward the target (%d, %d) with unit step cost and the Manhattan heuristic. The open list holds %s. Which node does A* expand next? Answer with its number.", t.x, t.y, strings.Join(text, "; ")),
				edakit.IntStr(best+1),
				"A* expands the smallest f = g + h, not the smallest g: "+strings.Join(ft, ", ")+".",
				map[string]string{"node": edakit.IntStr(best + 1)})
		}
	default:
		if edakit.Int(0, 1) == 0 {
			ws := []float64{1.2, 1.5, 2, 2.5, 3}
			w := edakit.Pick(ws)
			c := edakit.Int(20, 90)
			return edakit.Numeric(edakit.Sprintf("A router runs weighted A* with f = g + %s·h, where h is admissible. It returns a path of cost %d. What is the smallest possible cost of the optimal path? Round to 2 decimal places.", edakit.Num(w, 1), c),
				edakit.Num(float64(c)/w, 2),
				edakit.Sprintf("Weighted A* with an admissible h returns a path at most w times the optimum, so optimum ≥ %d / %s = %s.", c, edakit.Num(w, 1), edakit.Num(float64(c)/w, 2)),
				map[string]string{"bound": edakit.Num(float64(c)/w, 2)})
		}
		// Heuristic admissibility count on a small weighted graph toward target.
		n := edakit.Scaled(ctx, 5, 7)
		g := randWGraph(n, edakit.Int(2, 4), 6)
		d, _, _ := dijkstraTrace(g, n-1)
		var hs []string
		bad := 0
		for v := 0; v < n-1; v++ {
			h := d[v] + edakit.Int(-3, 2)
			if h < 0 {
				h = 0
			}
			if h > d[v] {
				bad++
			}
			hs = append(hs, edakit.Sprintf("h(%s) = %d", vName(v), h))
		}
		return edakit.Numeric("A search graph has vertices v1 to v"+edakit.IntStr(n)+" and undirected edge costs "+edgeText(g, true)+". The target is "+vName(n-1)+". A proposed heuristic gives "+strings.Join(hs, ", ")+". For how many vertices does it overestimate the true remaining cost (violating admissibility)?",
			edakit.IntStr(bad),
			"True distances to the target are "+intsText(d)+". A vertex violates admissibility when h exceeds its true distance; "+edakit.IntStr(bad)+" do.",
			map[string]string{"violations": edakit.IntStr(bad)})
	}
}
