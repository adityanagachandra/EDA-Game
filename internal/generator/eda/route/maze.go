package route

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// lee: the target's wavefront label (KP 1), cells labelled before the target (KP 2), label
// memory for a full-label, mod-3 or Akers coding (KP 3).
func lee(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k == 3 {
		W := 100 * edakit.Scaled(ctx, 3, 60)
		H := 100 * edakit.Scaled(ctx, 3, 60)
		scheme := edakit.Int(0, 2)
		var bits int
		var desc, why string
		switch scheme {
		case 0:
			maxLabel := W*H - 1
			bits = int(math.Ceil(math.Log2(float64(maxLabel + 3))))
			desc = edakit.Sprintf("stores the full wavefront label in every cell. A label can be as large as the number of cells minus 1, and two more codes mark empty and blocked cells, so each cell needs the fewest bits that can hold %d distinct values", maxLabel+3)
			why = edakit.Sprintf("Labels 0 to %d plus empty and blocked make %d codes, which need ceil(log2 %d) = %d bits.", maxLabel, maxLabel+3, maxLabel+3, bits)
		case 1:
			bits = 3
			desc = "stores each label modulo 3 (the sequence 1, 2, 3, 1, 2, 3, ...), which is enough to backtrace because a cell's predecessor label is always the one before it in the cycle. With codes for empty and blocked, each cell needs the fewest bits that can hold all its codes"
			why = "Three label values plus empty and blocked make 5 codes, which need 3 bits."
		default:
			bits = 2
			desc = "uses Akers' coding, labelling the wavefronts 0, 0, 1, 1, 0, 0, 1, 1, ... With codes for empty and blocked, each cell needs the fewest bits that can hold all its codes"
			why = "Two label values plus empty and blocked make 4 codes, which fit in 2 bits."
		}
		total := float64(W*H*bits) / 1e6
		ans := edakit.Num(total, 2)
		q := edakit.Sprintf("A Lee router on a %d by %d detailed-routing grid %s. How much label memory does the grid need, in millions of bits? Round to 2 decimal places.", W, H, desc)
		return edakit.Numeric(q, ans,
			why+edakit.Sprintf(" The grid has %d cells, so the memory is %d x %d bits = %s million bits.", W*H, W*H, bits, ans),
			map[string]string{"bits_per_cell": edakit.IntStr(bits)})
	}
	g, s, t, dist := randGrid(ctx)
	D := dist[t]
	intro := gridIntro(g, s, t)
	if k == 1 {
		return edakit.Numeric(intro+"Lee's algorithm labels the source 0 and expands one wavefront at a time. What label does the target receive, which is the length of the shortest route in grid steps?", edakit.IntStr(D),
			edakit.Sprintf("The Manhattan distance is %d, but the blockages force a detour. Breadth-first wavefronts reach the target at label %d.", manh(s, t), D),
			map[string]string{"label": edakit.IntStr(D), "manhattan": edakit.IntStr(manh(s, t))})
	}
	n := 0
	for _, d := range dist {
		if d < D {
			n++
		}
	}
	return edakit.Numeric(intro+"Lee's algorithm labels the source 0 and stops as soon as the target is labelled, so every cell with a smaller label than the target has been labelled by then. How many cells, including the source, carry a label strictly smaller than the target's?", edakit.IntStr(n),
		edakit.Sprintf("The target's label is %d. Counting all reachable cells with labels 0 to %d gives %d cells. This is the expansion cost of Lee's search: it grows with the area of the wavefront, not with the length of the route.", D, D-1, n),
		map[string]string{"target_label": edakit.IntStr(D), "expanded": edakit.IntStr(n)})
}

// backtrace: choices at the target's first backtrace step (KP 1), the number of shortest routes
// (KP 2), the fewest bends among shortest routes (KP 3).
func backtrace(ctx generator.GeneratorContext) generator.Problem {
	g, s, t, dist := randGrid(ctx)
	D := dist[t]
	intro := gridIntro(g, s, t)
	switch kp(ctx) {
	case 1:
		n := 0
		for _, d := range dirs4 {
			v := pt{t.x + d.x, t.y + d.y}
			if dv, ok := dist[v]; ok && dv == D-1 {
				n++
			}
		}
		return edakit.Numeric(intro+"After Lee's wavefront labels the target, the backtrace steps from the target to a neighbour whose label is exactly one less. How many neighbours of the target qualify for that first step?", edakit.IntStr(n),
			edakit.Sprintf("The target has label %d, and %d of its free neighbours carry label %d. Each is the end of at least one shortest route, which is where a backtrace tie-break rule matters.", D, n, D-1),
			map[string]string{"target_label": edakit.IntStr(D), "choices": edakit.IntStr(n)})
	case 2:
		cnt := countShortest(g, s, dist)
		return edakit.Numeric(intro+"How many distinct shortest routes join the source and the target? A backtrace could return any one of them.", edakit.IntStr(cnt[t]),
			edakit.Sprintf("Count routes wavefront by wavefront: each cell's count is the sum of the counts of its neighbours labelled one less. The target, at label %d, collects %d shortest routes.", D, cnt[t]),
			map[string]string{"length": edakit.IntStr(D), "routes": edakit.IntStr(cnt[t])})
	default:
		b := minBendsShortest(g, s, t, dist)
		return edakit.Numeric(intro+"Among all shortest routes from source to target, what is the fewest number of bends (changes of direction)? A backtrace that prefers to keep its current direction aims for this.", edakit.IntStr(b),
			edakit.Sprintf("Restrict to moves that raise the label by one, then minimise direction changes over that shortest-route DAG. The best shortest route has %d bends. Every bend costs a via in a preferred-direction stack, which is why backtrace tie-breaking matters.", b),
			map[string]string{"length": edakit.IntStr(D), "bends": edakit.IntStr(b)})
	}
}

func countShortest(g grid, s pt, dist map[pt]int) map[pt]int {
	maxd := 0
	for _, d := range dist {
		if d > maxd {
			maxd = d
		}
	}
	layers := make([][]pt, maxd+1)
	for p, d := range dist {
		layers[d] = append(layers[d], p)
	}
	cnt := map[pt]int{s: 1}
	for d := 1; d <= maxd; d++ {
		for _, p := range layers[d] {
			c := 0
			for _, dv := range dirs4 {
				u := pt{p.x + dv.x, p.y + dv.y}
				if du, ok := dist[u]; ok && du == d-1 {
					c += cnt[u]
				}
			}
			cnt[p] = c
		}
	}
	return cnt
}

func minBendsShortest(g grid, s, t pt, dist map[pt]int) int {
	type st struct {
		p pt
		d int
	}
	const inf = 1 << 30
	best := map[st]int{}
	for d := range dirs4 {
		best[st{s, d}] = 0
	}
	maxd := dist[t]
	layers := make([][]pt, maxd+1)
	for p, d := range dist {
		if d <= maxd {
			layers[d] = append(layers[d], p)
		}
	}
	for d := 0; d < maxd; d++ {
		for _, p := range layers[d] {
			for dirIn := range dirs4 {
				c, ok := best[st{p, dirIn}]
				if !ok {
					continue
				}
				for dd, dv := range dirs4 {
					v := pt{p.x + dv.x, p.y + dv.y}
					if dvv, ok := dist[v]; !ok || dvv != d+1 {
						continue
					}
					nc := c
					if p != s && dd != dirIn {
						nc++
					}
					k := st{v, dd}
					if old, ok := best[k]; !ok || nc < old {
						best[k] = nc
					}
				}
			}
		}
	}
	ans := inf
	for d := range dirs4 {
		if v, ok := best[st{t, d}]; ok && v < ans {
			ans = v
		}
	}
	return ans
}

// hadlock: detours of a given move string (KP 1), the minimum detour number to the target on a
// blocked grid (KP 2), Hadlock's label of a given cell (KP 3).
func hadlock(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		dx, dy := edakit.Int(2, 6), edakit.Int(1, 5)
		a, b := edakit.Int(0, 2), edakit.Int(0, 2)
		if a+b == 0 {
			a = 1
		}
		var moves []string
		for i := 0; i < dx+a; i++ {
			moves = append(moves, "E")
		}
		for i := 0; i < a; i++ {
			moves = append(moves, "W")
		}
		for i := 0; i < dy+b; i++ {
			moves = append(moves, "N")
		}
		for i := 0; i < b; i++ {
			moves = append(moves, "S")
		}
		moves = edakit.Shuffle(moves)
		s := pt{edakit.Int(1, 4), edakit.Int(1, 4)}
		t := pt{s.x + dx, s.y + dy}
		L := len(moves)
		q := edakit.Sprintf("A route leaves the source at %s for the target at %s with the moves %s (E is +1 column, N is +1 row). Hadlock's router counts a detour for every move that takes the route farther from the target in Manhattan distance. How many detours does this route make?", ptText(s), ptText(t), strings.Join(moves, " "))
		return edakit.Numeric(q, edakit.IntStr(a+b),
			edakit.Sprintf("The route has %d moves and the Manhattan distance is %d. Every detour must be undone by an extra move toward the target, so length = distance + 2 x detours, giving (%d - %d) / 2 = %d.", L, dx+dy, L, dx+dy, a+b),
			map[string]string{"length": edakit.IntStr(L), "manhattan": edakit.IntStr(dx + dy)})
	case 2:
		g, s, t, dist := randGrid(ctx)
		d := (dist[t] - manh(s, t)) / 2
		return edakit.Numeric(gridIntro(g, s, t)+"Hadlock's router labels cells by the number of detours (moves away from the target) needed to reach them and expands the smallest label first. What is the detour number of the best route to the target?", edakit.IntStr(d),
			edakit.Sprintf("The shortest route has length %d and the Manhattan distance is %d. Since length = distance + 2d, the minimum detour number is (%d - %d) / 2 = %d.", dist[t], manh(s, t), dist[t], manh(s, t), d),
			map[string]string{"length": edakit.IntStr(dist[t]), "detours": edakit.IntStr(d)})
	default:
		g, s, t, dist := randGrid(ctx)
		var cands []pt
		for p, dp := range dist {
			if p != s && p != t && dp >= 3 {
				cands = append(cands, p)
			}
		}
		sortPts(cands)
		v := edakit.Pick(cands)
		lab := (dist[v] - manh(s, t) + manh(v, t)) / 2
		return edakit.Numeric(gridIntro(g, s, t)+edakit.Sprintf("Hadlock's router labels each cell with the fewest detours (moves that increase the Manhattan distance to the target) on any route from the source to that cell. What label does cell %s receive?", ptText(v)), edakit.IntStr(lab),
			edakit.Sprintf("A route of length l from the source to v makes d detours with l = (MD(s,t) - MD(v,t)) + 2d. The shortest route to v has length %d, MD(s,t) = %d and MD(v,t) = %d, so d = (%d - %d + %d) / 2 = %d.", dist[v], manh(s, t), manh(v, t), dist[v], manh(s, t), manh(v, t), lab),
			map[string]string{"dist": edakit.IntStr(dist[v]), "label": edakit.IntStr(lab)})
	}
}

// astarRouting: f = g + h of a cell (KP 1), cells A* must expand (KP 2), cells A* saves against
// Lee (KP 3).
func astarRouting(ctx generator.GeneratorContext) generator.Problem {
	g, s, t, dist := randGrid(ctx)
	D := dist[t]
	intro := gridIntro(g, s, t)
	k := kp(ctx)
	if k == 1 {
		var cands []pt
		for p, dp := range dist {
			if p != s && p != t && dp >= 2 {
				cands = append(cands, p)
			}
		}
		sortPts(cands)
		v := edakit.Pick(cands)
		f := dist[v] + manh(v, t)
		return edakit.Numeric(intro+edakit.Sprintf("An A* router uses unit step cost and the Manhattan distance to the target as its heuristic. Cell %s has been reached by its shortest route from the source. What is its priority f = g + h?", ptText(v)), edakit.IntStr(f),
			edakit.Sprintf("g is the shortest-route length from the source, %d, and h is the Manhattan distance to the target, %d. f = %d.", dist[v], manh(v, t), f),
			map[string]string{"g": edakit.IntStr(dist[v]), "h": edakit.IntStr(manh(v, t))})
	}
	nf, ng := 0, 0
	for p, dp := range dist {
		if dp < D {
			ng++
		}
		if dp+manh(p, t) < D {
			nf++
		}
	}
	if k == 2 {
		return edakit.Numeric(intro+edakit.Sprintf("An A* router uses unit step cost and the Manhattan heuristic, which is consistent. Whatever its tie-breaking, it must expand every cell whose f = g + h is strictly below the cost of the best route, %d. How many cells is that, including the source?", D), edakit.IntStr(nf),
			edakit.Sprintf("Compute g by breadth-first search and add each cell's Manhattan distance to the target. %d cells have f < %d. Cells with f = %d may or may not be expanded, depending on ties.", nf, D, D),
			map[string]string{"best": edakit.IntStr(D), "a_star_must": edakit.IntStr(nf)})
	}
	return edakit.Numeric(intro+edakit.Sprintf("The best route costs %d. Lee's breadth-first router must expand every cell with label below %d, and an A* router with the Manhattan heuristic must expand every cell with f = g + h below %d. How many more cells is Lee guaranteed to expand than A*?", D, D, D), edakit.IntStr(ng-nf),
		edakit.Sprintf("Lee must expand %d cells (label < %d); A* must expand %d (f < %d). The difference is %d. The heuristic prunes cells that lie away from the target, which is the whole point of goal-directed search in a router.", ng, D, nf, D, ng-nf),
		map[string]string{"lee": edakit.IntStr(ng), "a_star": edakit.IntStr(nf)})
}

// moveCost parses a move string and returns length and bends.
func moveStats(moves []string) (int, int) {
	b := 0
	for i := 1; i < len(moves); i++ {
		if moves[i] != moves[i-1] {
			b++
		}
	}
	return len(moves), b
}

// minBendCost runs Dijkstra over (cell, direction) states with unit step cost plus bend cost.
func minBendCost(g grid, s, t pt, bend int) int {
	type st struct {
		p pt
		d int
	}
	const inf = 1 << 30
	best := map[st]int{}
	for d := range dirs4 {
		best[st{s, d}] = 0
	}
	done := map[st]bool{}
	for {
		var u st
		bu := inf
		for k, v := range best {
			if !done[k] && (v < bu || (v == bu && (k.p.x < u.p.x || (k.p.x == u.p.x && (k.p.y < u.p.y || (k.p.y == u.p.y && k.d < u.d)))))) {
				u, bu = k, v
			}
		}
		if bu == inf {
			break
		}
		done[u] = true
		for d, dv := range dirs4 {
			v := pt{u.p.x + dv.x, u.p.y + dv.y}
			if !g.free(v) {
				continue
			}
			c := bu + 1
			if u.p != s && d != u.d {
				c += bend
			}
			k := st{v, d}
			if old, ok := best[k]; !ok || c < old {
				best[k] = c
			}
		}
	}
	ans := inf
	for d := range dirs4 {
		if v, ok := best[st{t, d}]; ok && v < ans {
			ans = v
		}
	}
	return ans
}

// bendViaCost: cost of a described route (KP 1), the min-cost route with a bend penalty (KP 2),
// the via cost at which a longer, straighter route wins (KP 3).
func bendViaCost(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		var moves []string
		n := edakit.Int(4, 7)
		cur := edakit.Pick([]string{"E", "N"})
		for i := 0; i < n; i++ {
			run := edakit.Int(1, 4)
			for j := 0; j < run; j++ {
				moves = append(moves, cur)
			}
			if cur == "E" {
				cur = "N"
			} else {
				cur = "E"
			}
		}
		L, B := moveStats(moves)
		b := edakit.Int(2, 6)
		cost := L + b*B
		q := edakit.Sprintf("A maze router charges 1 per grid step and %d per bend, since in a preferred-direction stack each bend needs a via. A candidate route makes the moves %s. What is its cost?", b, strings.Join(moves, " "))
		return edakit.Numeric(q, edakit.IntStr(cost),
			edakit.Sprintf("The route has %d steps and changes direction %d times, so it costs %d + %d x %d = %d.", L, B, L, b, B, cost),
			map[string]string{"length": edakit.IntStr(L), "bends": edakit.IntStr(B)})
	case 2:
		g, s, t, dist := randGrid(ctx)
		b := edakit.Int(1, 4)
		c := minBendCost(g, s, t, b)
		return edakit.Numeric(gridIntro(g, s, t)+edakit.Sprintf("The router charges 1 per step and %d per bend, and searches over (cell, incoming direction) states so it can price bends. What is the cost of the cheapest route?", b), edakit.IntStr(c),
			edakit.Sprintf("The shortest route has %d steps and the fewest-bend route has %d bends. With a bend cost of %d, Dijkstra over direction-aware states finds a cheapest cost of %d.", dist[t], g.minBends(s, t), b, c),
			map[string]string{"length": edakit.IntStr(dist[t]), "cost": edakit.IntStr(c)})
	default:
		L1 := edakit.Int(10, 30)
		k1 := edakit.Int(3, 7)
		k2 := edakit.Int(1, k1-1)
		L2 := L1 + 2*edakit.Int(1, 6)
		vmin := (L2-L1)/(k1-k2) + 1
		q := edakit.Sprintf("Two candidate routes join the same pins. Route A is %d grid steps with %d vias; route B is %d steps with %d vias. The router charges 1 per step and v per via. What is the smallest integer via cost v at which route B is strictly cheaper?", L1, k1, L2, k2)
		return edakit.Numeric(q, edakit.IntStr(vmin),
			edakit.Sprintf("B is cheaper when %d + %dv < %d + %dv, that is (%d - %d)v > %d - %d, so v > %s. The smallest integer is %d.", L2, k2, L1, k1, k1, k2, L2, L1, edakit.Num(float64(L2-L1)/float64(k1-k2), 3), vmin),
			map[string]string{"v_min": edakit.IntStr(vmin)})
	}
}

// lineProbing: fewest bends on any route (KP 1), cells on the source's level-0 probe lines
// (KP 2), cells reachable with at most one bend (KP 3).
func lineProbing(ctx generator.GeneratorContext) generator.Problem {
	g, s, t, _ := randGrid(ctx)
	intro := gridIntro(g, s, t)
	switch kp(ctx) {
	case 1:
		b := g.minBends(s, t)
		return edakit.Numeric(intro+"A Mikami-Tabuchi line-probe router extends straight escape lines from both pins and adds a level of perpendicular lines at each round, so it finds a route with the fewest bends. How many bends does that route have?", edakit.IntStr(b),
			edakit.Sprintf("Searching over (cell, direction) with a cost only for turning gives a minimum of %d bends, that is %d straight segments.", b, b+1),
			map[string]string{"bends": edakit.IntStr(b)})
	case 2:
		line := probeLine(g, s)
		n := len(line)
		return edakit.Numeric(intro+"The source's level-0 probe lines run from the source in all four directions until they hit a blocked cell or the grid edge. How many cells, not counting the source, lie on those lines?", edakit.IntStr(n),
			edakit.Sprintf("Walking east, west, north and south from %s until blocked covers %d cells. A line-probe router stores these as a few segments rather than labelling every cell.", ptText(s), n),
			map[string]string{"cells": edakit.IntStr(n)})
	default:
		reach := map[pt]bool{}
		l0 := probeLine(g, s)
		for _, p := range l0 {
			reach[p] = true
		}
		for _, p := range append([]pt{s}, l0...) {
			for _, q := range probeLine(g, p) {
				reach[q] = true
			}
		}
		delete(reach, s)
		n := len(reach)
		return edakit.Numeric(intro+"Line probing reaches a cell with at most one bend if the cell lies on a level-0 line from the source or on a perpendicular level-1 line launched from some point of those lines. How many cells other than the source can a route reach from the source with at most one bend?", edakit.IntStr(n),
			edakit.Sprintf("Take the cells on the source's straight lines, then extend straight lines from every one of those points in all four directions. The union, without the source, has %d cells.", n),
			map[string]string{"cells": edakit.IntStr(n), "level0": edakit.IntStr(len(l0))})
	}
}

func probeLine(g grid, s pt) []pt {
	var out []pt
	for _, d := range dirs4 {
		p := pt{s.x + d.x, s.y + d.y}
		for g.free(p) {
			out = append(out, p)
			p = pt{p.x + d.x, p.y + d.y}
		}
	}
	return out
}
