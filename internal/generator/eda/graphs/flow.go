package graphs

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// net is a flow network: node 0 is the source s, node n-1 the sink t.
type network struct {
	n   int
	cap [][]int
}

func fName(n, i int) string {
	switch i {
	case 0:
		return "s"
	case n - 1:
		return "t"
	}
	return string(rune('a' + i - 1))
}

func randNetwork(n, extra, cmax int) network {
	nw := network{n: n, cap: make([][]int, n)}
	for i := range nw.cap {
		nw.cap[i] = make([]int, n)
	}
	// Every middle node lies on some s-t route.
	for v := 1; v < n-1; v++ {
		u := edakit.Int(0, v-1)
		nw.cap[u][v] = edakit.Int(1, cmax)
		w := edakit.Int(v+1, n-1)
		nw.cap[v][w] = edakit.Int(1, cmax)
	}
	for t := 0; t < extra; t++ {
		u := edakit.Int(0, n-2)
		v := edakit.Int(u+1, n-1)
		if u == 0 && v == n-1 {
			continue
		}
		if nw.cap[u][v] == 0 {
			nw.cap[u][v] = edakit.Int(1, cmax)
		}
	}
	return nw
}

func netArcsText(nw network, flow [][]int) string {
	var parts []string
	for u := 0; u < nw.n; u++ {
		for v := 0; v < nw.n; v++ {
			if nw.cap[u][v] > 0 {
				if flow != nil {
					parts = append(parts, edakit.Sprintf("%s→%s (flow %d of %d)", fName(nw.n, u), fName(nw.n, v), flow[u][v], nw.cap[u][v]))
				} else {
					parts = append(parts, edakit.Sprintf("%s→%s (%d)", fName(nw.n, u), fName(nw.n, v), nw.cap[u][v]))
				}
			}
		}
	}
	return strings.Join(parts, ", ")
}

// bfsAugment finds a shortest augmenting path in the residual graph (neighbours in index order).
func bfsAugment(nw network, flow [][]int) []int {
	n := nw.n
	prev := make([]int, n)
	for i := range prev {
		prev[i] = -1
	}
	prev[0] = 0
	q := []int{0}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for v := 0; v < n; v++ {
			if prev[v] < 0 && nw.cap[u][v]-flow[u][v]+flow[v][u] > 0 {
				prev[v] = u
				q = append(q, v)
			}
		}
	}
	if prev[n-1] < 0 {
		return nil
	}
	var path []int
	for v := n - 1; v != 0; v = prev[v] {
		path = append([]int{v}, path...)
	}
	return append([]int{0}, path...)
}

func residual(nw network, flow [][]int, u, v int) int {
	return nw.cap[u][v] - flow[u][v] + flow[v][u]
}

func push(nw network, flow [][]int, path []int) int {
	b := math.MaxInt32
	for i := 0; i+1 < len(path); i++ {
		b = min(b, residual(nw, flow, path[i], path[i+1]))
	}
	for i := 0; i+1 < len(path); i++ {
		u, v := path[i], path[i+1]
		// Cancel opposite flow first, then add forward flow.
		c := min(b, flow[v][u])
		flow[v][u] -= c
		flow[u][v] += b - c
	}
	return b
}

func newFlow(n int) [][]int {
	f := make([][]int, n)
	for i := range f {
		f[i] = make([]int, n)
	}
	return f
}

func maxFlow(nw network) (int, [][]int) {
	flow := newFlow(nw.n)
	total := 0
	for {
		p := bfsAugment(nw, flow)
		if p == nil {
			return total, flow
		}
		total += push(nw, flow, p)
	}
}

func pathText(n int, p []int) string {
	var s []string
	for _, v := range p {
		s = append(s, fName(n, v))
	}
	return strings.Join(s, " → ")
}

func maxFlowGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 7)
	nw := randNetwork(n, edakit.Scaled(ctx, 2, 5), 9)
	mf, flow := maxFlow(nw)
	intro := "A flow network has arcs with capacities " + netArcsText(nw, nil) + ". "
	switch kpOf(ctx) {
	case 1:
		reach := residualReach(nw, flow)
		return edakit.Numeric(intro+"What is the maximum flow value from s to t?", edakit.IntStr(mf),
			edakit.Sprintf("Augmenting until no residual s-t path remains gives %d. The nodes still reachable from s in the residual graph, %s, form a minimum cut whose capacity is also %d.", mf, reachText(n, reach), mf),
			map[string]string{"max_flow": edakit.IntStr(mf)})
	case 2:
		S := map[int]bool{0: true}
		for v := 1; v < n-1; v++ {
			if edakit.Int(0, 1) == 0 {
				S[v] = true
			}
		}
		c := 0
		var terms []string
		for u := range S {
			for v := 0; v < n; v++ {
				if !S[v] && nw.cap[u][v] > 0 {
					c += nw.cap[u][v]
					terms = append(terms, edakit.Sprintf("%s→%s", fName(n, u), fName(n, v)))
				}
			}
		}
		return edakit.Numeric(intro+"Take the s-t cut with source side "+reachText(n, S)+". What is its capacity?", edakit.IntStr(c),
			edakit.Sprintf("Only arcs from the source side to the sink side count (%s), summing to %d. Arcs pointing back do not. Every cut capacity bounds the max flow (%d here) from above.", strings.Join(terms, ", "), c, mf),
			map[string]string{"cut": edakit.IntStr(c), "max_flow": edakit.IntStr(mf)})
	default:
		type arc struct{ u, v int }
		var arcs []arc
		for u := 0; u < n; u++ {
			for v := 0; v < n; v++ {
				if nw.cap[u][v] > 0 {
					arcs = append(arcs, arc{u, v})
				}
			}
		}
		e := edakit.Pick(arcs)
		d := edakit.Int(2, 6)
		nw2 := network{n: n, cap: make([][]int, n)}
		for i := range nw.cap {
			nw2.cap[i] = append([]int(nil), nw.cap[i]...)
		}
		nw2.cap[e.u][e.v] += d
		mf2, _ := maxFlow(nw2)
		return edakit.Numeric(intro+edakit.Sprintf("Arc %s→%s is widened by %d. By how much does the maximum flow increase?", fName(n, e.u), fName(n, e.v), d),
			edakit.IntStr(mf2-mf),
			edakit.Sprintf("Max flow goes from %d to %d. Widening helps only if the arc crosses every minimum cut, and then by at most the added capacity.", mf, mf2),
			map[string]string{"before": edakit.IntStr(mf), "after": edakit.IntStr(mf2)})
	}
}

func residualReach(nw network, flow [][]int) map[int]bool {
	reach := map[int]bool{0: true}
	q := []int{0}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for v := 0; v < nw.n; v++ {
			if !reach[v] && residual(nw, flow, u, v) > 0 {
				reach[v] = true
				q = append(q, v)
			}
		}
	}
	return reach
}

func reachText(n int, S map[int]bool) string {
	var names []string
	for v := 0; v < n; v++ {
		if S[v] {
			names = append(names, fName(n, v))
		}
	}
	return "{" + strings.Join(names, ", ") + "}"
}

func residualGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		c := edakit.Int(3, 15)
		f := edakit.Int(0, c)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(edakit.Sprintf("Arc u→v has capacity %d and carries flow %d. What is the residual capacity of the forward arc u→v?", c, f),
				edakit.IntStr(c-f), edakit.Sprintf("Forward residual = capacity - flow = %d - %d = %d.", c, f, c-f), map[string]string{"forward": edakit.IntStr(c - f)})
		}
		return edakit.Numeric(edakit.Sprintf("Arc u→v has capacity %d and carries flow %d, and there is no arc v→u. What is the residual capacity of the backward arc v→u?", c, f),
			edakit.IntStr(f), edakit.Sprintf("The backward residual arc lets the algorithm cancel flow already sent, so its capacity is the flow, %d.", f), map[string]string{"backward": edakit.IntStr(f)})
	case 2:
		for {
			n := edakit.Scaled(ctx, 5, 7)
			nw := randNetwork(n, edakit.Scaled(ctx, 3, 6), 9)
			flow := newFlow(n)
			p1 := bfsAugment(nw, flow)
			if p1 == nil {
				continue
			}
			push(nw, flow, p1)
			p2 := bfsAugment(nw, flow)
			if p2 == nil {
				continue
			}
			b := math.MaxInt32
			var rs []string
			for i := 0; i+1 < len(p2); i++ {
				r := residual(nw, flow, p2[i], p2[i+1])
				b = min(b, r)
				rs = append(rs, edakit.IntStr(r))
			}
			return edakit.Numeric("A flow network currently carries "+netArcsText(nw, flow)+". Consider the augmenting path "+pathText(n, p2)+" in the residual graph (a step against an arc uses that arc's flow as residual capacity). How much flow can be pushed along it?",
				edakit.IntStr(b),
				"The residual capacities along the path are "+strings.Join(rs, ", ")+"; the bottleneck is "+edakit.IntStr(b)+".",
				map[string]string{"bottleneck": edakit.IntStr(b)})
		}
	default:
		n := edakit.Scaled(ctx, 5, 7)
		nw := randNetwork(n, edakit.Scaled(ctx, 2, 5), 9)
		mf, flow := maxFlow(nw)
		reach := residualReach(nw, flow)
		return edakit.Numeric("A flow network has capacities "+netArcsText(nw, nil)+edakit.Sprintf(". A maximum flow of value %d is found, giving the flows ", mf)+netArcsText(nw, flow)+". How many nodes, counting s, are reachable from s in the final residual graph?",
			edakit.IntStr(len(reach)),
			"The reachable set is "+reachText(n, reach)+". Its outgoing arcs are all saturated, so it is the source side of a minimum cut of capacity "+edakit.IntStr(mf)+".",
			map[string]string{"source_side": edakit.IntStr(len(reach))})
	}
}

// minNetCut is the minimum (weighted) number of nets cut over all partitions with cell s on side
// 0 and cell t on side 1, by enumeration.
func minNetCut(h hyper, s, t int, weighted bool) int {
	best := math.MaxInt32
	var free []int
	for i := 0; i < h.n; i++ {
		if i != s && i != t {
			free = append(free, i)
		}
	}
	side := make([]int, h.n)
	side[t] = 1
	for mask := 0; mask < 1<<len(free); mask++ {
		for b, c := range free {
			side[c] = (mask >> b) & 1
		}
		best = min(best, cutSize(h, side, weighted))
	}
	return best
}

func hyperMinCutGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	h := randHyper(n, edakit.Scaled(ctx, 4, 8), 4)
	switch kpOf(ctx) {
	case 1:
		pins := 0
		for _, net := range h.nets {
			pins += len(net)
		}
		m := len(h.nets)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric("Nets: "+netText(h, false)+edakit.Sprintf(". Lawler's construction turns this %d-cell netlist into a flow network: each net becomes two nodes joined by one arc of capacity equal to the net weight, with infinite-capacity arcs from each of its cells into the first node and from the second node back to each cell. How many nodes does the network have?", n),
				edakit.IntStr(n+2*m), edakit.Sprintf("%d cell nodes plus 2 per net: %d + 2 × %d = %d.", n, n, m, n+2*m), map[string]string{"nodes": edakit.IntStr(n + 2*m)})
		}
		return edakit.Numeric("Nets: "+netText(h, false)+edakit.Sprintf(". Lawler's construction turns this %d-cell netlist into a flow network: each net becomes two nodes joined by one arc of capacity equal to the net weight, with infinite-capacity arcs from each of its cells into the first node and from the second node back to each cell. How many arcs does the network have?", n),
			edakit.IntStr(m+2*pins), edakit.Sprintf("One bridge per net plus two arcs per pin: %d + 2 × %d = %d.", m, pins, m+2*pins), map[string]string{"arcs": edakit.IntStr(m + 2*pins)})
	default:
		weighted := ctx.KPOrdinal == 3 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0)
		if weighted {
			for j := range h.w {
				h.w[j] = edakit.Int(1, 5)
			}
		}
		s, t := 0, n-1
		c := minNetCut(h, s, t, weighted)
		q := "Nets: " + netText(h, weighted) + ". "
		if weighted {
			q += "What is the minimum total weight of nets that must be cut to separate " + cellName(s) + " from " + cellName(t) + " (no balance constraint)?"
		} else {
			q += "What is the minimum number of nets that must be cut to separate " + cellName(s) + " from " + cellName(t) + " (no balance constraint)?"
		}
		return edakit.Numeric(q, edakit.IntStr(c),
			"This is a minimum s-t cut in the hypergraph. In Lawler's network every finite-capacity arc is a net bridge, so a minimum cut there removes whole nets, and max-flow computes it. Here the minimum is "+edakit.IntStr(c)+".",
			map[string]string{"min_cut": edakit.IntStr(c)})
	}
}

// mcArc is a min-cost-flow arc with residual bookkeeping.
type mcArc struct{ to, cap, cost, rev int }

type mcGraph struct{ adj [][]mcArc }

func newMC(n int) *mcGraph { return &mcGraph{adj: make([][]mcArc, n)} }

func (g *mcGraph) add(u, v, cap, cost int) {
	g.adj[u] = append(g.adj[u], mcArc{v, cap, cost, len(g.adj[v])})
	g.adj[v] = append(g.adj[v], mcArc{u, 0, -cost, len(g.adj[u]) - 1})
}

// minCostFlow sends up to need units from s to t by successive shortest paths (Bellman-Ford),
// returning the flow sent and its cost.
func (g *mcGraph) minCostFlow(s, t, need int) (int, int) {
	n := len(g.adj)
	flow, cost := 0, 0
	for flow < need {
		dist := make([]int, n)
		pv, pe := make([]int, n), make([]int, n)
		for i := range dist {
			dist[i] = math.MaxInt32
		}
		dist[s] = 0
		for it := 0; it < n; it++ {
			upd := false
			for u := 0; u < n; u++ {
				if dist[u] == math.MaxInt32 {
					continue
				}
				for i, a := range g.adj[u] {
					if a.cap > 0 && dist[u]+a.cost < dist[a.to] {
						dist[a.to] = dist[u] + a.cost
						pv[a.to], pe[a.to] = u, i
						upd = true
					}
				}
			}
			if !upd {
				break
			}
		}
		if dist[t] == math.MaxInt32 {
			break
		}
		f := need - flow
		for v := t; v != s; v = pv[v] {
			f = min(f, g.adj[pv[v]][pe[v]].cap)
		}
		for v := t; v != s; v = pv[v] {
			a := &g.adj[pv[v]][pe[v]]
			a.cap -= f
			g.adj[v][a.rev].cap += f
		}
		flow += f
		cost += f * dist[t]
	}
	return flow, cost
}

func minCostGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1, 2:
		for {
			n := edakit.Scaled(ctx, 4, 6)
			nw := randNetwork(n, edakit.Scaled(ctx, 2, 4), 6)
			cost := make([][]int, n)
			var parts []string
			g := newMC(n)
			for u := 0; u < n; u++ {
				cost[u] = make([]int, n)
				for v := 0; v < n; v++ {
					if nw.cap[u][v] > 0 {
						cost[u][v] = edakit.Int(1, 7)
						g.add(u, v, nw.cap[u][v], cost[u][v])
						parts = append(parts, edakit.Sprintf("%s→%s (capacity %d, cost %d)", fName(n, u), fName(n, v), nw.cap[u][v], cost[u][v]))
					}
				}
			}
			mf, _ := maxFlow(nw)
			one := ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0)
			d := 1
			if !one {
				if mf < 3 {
					continue
				}
				d = edakit.Int(2, mf)
			}
			sent, c := g.minCostFlow(0, n-1, d)
			if sent < d {
				continue
			}
			if one {
				return edakit.Numeric("A network has arcs "+strings.Join(parts, ", ")+". Costs are per unit of flow. What is the minimum cost of sending 1 unit from s to t?",
					edakit.IntStr(c), "With one unit, min-cost flow is the cheapest s-t path by cost: "+edakit.IntStr(c)+".", map[string]string{"cost": edakit.IntStr(c)})
			}
			return edakit.Numeric("A network has arcs "+strings.Join(parts, ", ")+edakit.Sprintf(". Costs are per unit of flow. What is the minimum total cost of sending %d units from s to t?", d),
				edakit.IntStr(c),
				edakit.Sprintf("Successive shortest paths: repeatedly send flow along the cheapest residual path (backward arcs carry negative cost), until %d units are sent. The total is %d.", d, c),
				map[string]string{"cost": edakit.IntStr(c), "units": edakit.IntStr(d)})
		}
	default:
		for {
			k := edakit.Scaled(ctx, 5, 7)
			sup := make([]int, k)
			dem := make([]int, k)
			ts, td := 0, 0
			for i := 0; i < k; i++ {
				switch edakit.Int(0, 2) {
				case 0:
					sup[i] = edakit.Int(1, 4)
					ts += sup[i]
				case 1:
					dem[i] = edakit.Int(1, 4)
					td += dem[i]
				}
			}
			if ts == 0 || td < ts {
				continue
			}
			g := newMC(k + 2)
			for i := 0; i < k; i++ {
				if sup[i] > 0 {
					g.add(0, i+1, sup[i], 0)
				}
				if dem[i] > 0 {
					g.add(i+1, k+1, dem[i], 0)
				}
			}
			for i := 0; i < k; i++ {
				for j := 0; j < k; j++ {
					if sup[i] > 0 && dem[j] > 0 {
						g.add(i+1, j+1, 1000, iabs(i-j))
					}
				}
			}
			_, c := g.minCostFlow(0, k+1, ts)
			var bt []string
			for i := 0; i < k; i++ {
				switch {
				case sup[i] > 0:
					bt = append(bt, edakit.Sprintf("bin %d: excess %d", i+1, sup[i]))
				case dem[i] > 0:
					bt = append(bt, edakit.Sprintf("bin %d: room for %d", i+1, dem[i]))
				default:
					bt = append(bt, edakit.Sprintf("bin %d: full", i+1))
				}
			}
			return edakit.Numeric("A legalizer spreads cells along a row of bins: "+strings.Join(bt, "; ")+". Moving a cell from bin i to bin j costs |i - j|. Modeled as a min-cost flow (source to excess bins, excess bins to bins with room, room to sink), what is the minimum total movement to place every excess cell?",
				edakit.IntStr(c),
				edakit.Sprintf("All %d excess cells must find room; min-cost flow assigns them to minimize the total distance, %d. Greedy nearest-room choices can be worse because one move can block a cheaper one.", ts, c),
				map[string]string{"cost": edakit.IntStr(c), "excess": edakit.IntStr(ts)})
		}
	}
}

func maxBipartite(adj [][]int, nr int) (int, []int) {
	matchR := make([]int, nr)
	for i := range matchR {
		matchR[i] = -1
	}
	var try func(u int, seen []bool) bool
	try = func(u int, seen []bool) bool {
		for _, v := range adj[u] {
			if seen[v] {
				continue
			}
			seen[v] = true
			if matchR[v] < 0 || try(matchR[v], seen) {
				matchR[v] = u
				return true
			}
		}
		return false
	}
	m := 0
	for u := range adj {
		if try(u, make([]bool, nr)) {
			m++
		}
	}
	return m, matchR
}

func bipartiteGen(ctx generator.GeneratorContext) generator.Problem {
	for {
		nl := edakit.Scaled(ctx, 4, 7)
		nr := edakit.Int(nl-1, nl+1)
		adj := make([][]int, nl)
		var lt []string
		for u := 0; u < nl; u++ {
			k := edakit.Int(1, 3)
			opts := edakit.Shuffle(seq(nr))[:min(k, nr)]
			sortInts(opts)
			adj[u] = opts
			var rn []string
			for _, v := range opts {
				rn = append(rn, edakit.Sprintf("S%d", v+1))
			}
			lt = append(lt, edakit.Sprintf("P%d: %s", u+1, strings.Join(rn, ", ")))
		}
		m, _ := maxBipartite(adj, nr)
		greedy := 0
		used := make([]bool, nr)
		for u := 0; u < nl; u++ {
			for _, v := range adj[u] {
				if !used[v] {
					used[v] = true
					greedy++
					break
				}
			}
		}
		intro := edakit.Sprintf("A pin assigner must place %d pins P1 to P%d onto %d slots S1 to S%d. Each pin may use only these slots: %s. ", nl, nl, nr, nr, strings.Join(lt, "; "))
		switch kpOf(ctx) {
		case 1:
			return edakit.Numeric(intro+"What is the largest number of pins that can be assigned, one pin per slot?", edakit.IntStr(m),
				"This is a maximum bipartite matching; augmenting paths find size "+edakit.IntStr(m)+". Adding a source to every pin and every slot to a sink, all capacities 1, turns it into a max-flow problem with the same value.",
				map[string]string{"matching": edakit.IntStr(m)})
		case 2:
			if greedy == m && edakit.Int(0, 2) > 0 {
				continue
			}
			return edakit.Numeric(intro+"A greedy pass gives each pin, in order P1 onward, its lowest-numbered free allowed slot. How many more pins could augmenting paths add beyond the greedy result?",
				edakit.IntStr(m-greedy),
				edakit.Sprintf("Greedy assigns %d pins; the maximum matching has %d, so augmenting paths add %d. Each augmenting path reroutes earlier pins to free a slot.", greedy, m, m-greedy),
				map[string]string{"greedy": edakit.IntStr(greedy), "matching": edakit.IntStr(m)})
		default:
			return edakit.Numeric(intro+"What is the size of a minimum vertex cover of this pin-slot graph (fewest pins and slots touching every allowed pairing)?", edakit.IntStr(m),
				"By König's theorem the minimum vertex cover of a bipartite graph equals its maximum matching, "+edakit.IntStr(m)+". The cover certifies that no larger assignment exists.",
				map[string]string{"cover": edakit.IntStr(m)})
		}
	}
}

func sortInts(xs []int) {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j] < xs[j-1]; j-- {
			xs[j], xs[j-1] = xs[j-1], xs[j]
		}
	}
}

func permutations(n int) [][]int {
	if n == 0 {
		return [][]int{{}}
	}
	var out [][]int
	for _, p := range permutations(n - 1) {
		for i := 0; i <= len(p); i++ {
			q := append(append(append([]int(nil), p[:i]...), n-1), p[i:]...)
			out = append(out, q)
		}
	}
	return out
}

func assignmentGen(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	n := 3
	if kp >= 2 {
		n = 4
	}
	c := make([][]int, n)
	var rows []string
	for i := range c {
		c[i] = make([]int, n)
		for j := range c[i] {
			c[i][j] = edakit.Int(1, 9)
		}
		rows = append(rows, edakit.Sprintf("cell %d: %s", i+1, intsText(c[i])))
	}
	best := math.MaxInt32
	var bp []int
	for _, p := range permutations(n) {
		s := 0
		for i, j := range p {
			s += c[i][j]
		}
		if s < best {
			best, bp = s, p
		}
	}
	intro := edakit.Sprintf("A detailed placer reassigns %d cells to %d free slots in a window. The wirelength cost of putting each cell in slots 1 to %d is: %s. ", n, n, n, strings.Join(rows, "; "))
	var pt []string
	for i, j := range bp {
		pt = append(pt, edakit.Sprintf("cell %d → slot %d", i+1, j+1))
	}
	if kp == 3 {
		used := make([]bool, n)
		g := 0
		for i := 0; i < n; i++ {
			bj := -1
			for j := 0; j < n; j++ {
				if !used[j] && (bj < 0 || c[i][j] < c[i][bj]) {
					bj = j
				}
			}
			used[bj] = true
			g += c[i][bj]
		}
		return edakit.Numeric(intro+"A greedy pass lets each cell, in order, take its cheapest free slot (ties to the lower slot). How much more does the greedy assignment cost than the optimal one?",
			edakit.IntStr(g-best),
			edakit.Sprintf("Greedy costs %d. The optimum, %s, costs %d. The gap is %d.", g, strings.Join(pt, ", "), best, g-best),
			map[string]string{"greedy": edakit.IntStr(g), "optimal": edakit.IntStr(best)})
	}
	return edakit.Numeric(intro+"What is the minimum total cost of a one-to-one assignment?", edakit.IntStr(best),
		"The optimal assignment is "+strings.Join(pt, ", ")+", total "+edakit.IntStr(best)+". The Hungarian algorithm finds it in polynomial time; it is a min-cost perfect matching, a special min-cost flow.",
		map[string]string{"optimal": edakit.IntStr(best)})
}
