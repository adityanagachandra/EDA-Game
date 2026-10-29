package fpga

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const rrRule = "Model the fabric as an n × n array of logic blocks with a routing channel on every side of every block, so there are n + 1 horizontal channels and n + 1 vertical channels, each n tiles long. Every channel has W tracks of length-L wires laid end to end, and each wire is one routing-resource node."

func rrGraph(ctx generator.GeneratorContext) generator.Problem {
	l := edakit.Pick([]int{1, 2, 4})
	n := l * edakit.Int(1, 8/l+2)
	if n < 2 {
		n = 2
	}
	for n%l != 0 {
		n++
	}
	w := 2 * edakit.Int(2, 20)
	wires := 2 * (n + 1) * w * (n / l)
	switch whichKP(ctx) {
	case 1:
		q := rrRule + edakit.Sprintf(" How many wire nodes does the RR graph have for n = %d, W = %d and L = %d?", n, w, l)
		return num(q, edakit.IntStr(wires),
			edakit.Sprintf("Each direction has %d channels × %d tracks × %d/%d wires per track = %d nodes; both directions give %d.", n+1, w, n, l, wires/2, wires),
			map[string]string{"chanx": edakit.IntStr(wires / 2)})
	case 2:
		in, out := edakit.Int(8, 40), edakit.Int(4, 10)
		per := in + out + 2
		tot := wires + n*n*per
		q := rrRule + edakit.Sprintf(" Each logic block also contributes one IPIN node per input pin (%d inputs), one OPIN node per output pin (%d outputs), one SINK node shared by its logically equivalent inputs and one SOURCE node for its outputs. How many RR-graph nodes are there in total for n = %d, W = %d and L = %d?", in, out, n, w, l)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("Wires: 2 × %d × %d × %d = %d. Blocks: %d × (%d + %d + 2) = %d. Total %d. One shared SINK lets the router finish at any equivalent input pin.", n+1, w, n/l, wires, n*n, in, out, n*n*per, tot),
			map[string]string{"wire_nodes": edakit.IntStr(wires), "block_nodes": edakit.IntStr(n * n * per)})
	default:
		in, out := edakit.Int(8, 40), edakit.Int(4, 10)
		w = 8 * edakit.Int(1, 8)
		fin := edakit.Pick([]float64{0.25, 0.5, 1})
		fout := edakit.Pick([]float64{0.125, 0.25, 0.5})
		ci, co := int(fin*float64(w)), int(fout*float64(w))
		e := n * n * (in*ci + out*co)
		q := edakit.Sprintf("An RR graph for a %d × %d array of logic blocks, each with %d input pins and %d output pins, uses W = %d. Each input pin has an edge from Fc_in·W wires (Fc_in = %s) and each output pin an edge to Fc_out·W wires (Fc_out = %s). How many wire-to-pin and pin-to-wire edges does the graph contain?", n, n, in, out, w, edakit.Num(fin, 3), edakit.Num(fout, 3))
		return num(q, edakit.IntStr(e),
			edakit.Sprintf("Per block: %d × %d + %d × %d = %d edges; %d blocks give %d. Edge count, not node count, usually dominates RR-graph memory.", in, ci, out, co, in*ci+out*co, n*n, e),
			map[string]string{"per_block": edakit.IntStr(in*ci + out*co)})
	}
}

type wgraph struct {
	m     int
	cost  []int
	adj   [][]int
	src   []int
	sinks [][]int
}

func newWGraph(m int, nsinks int) *wgraph {
	g := &wgraph{m: m, cost: make([]int, m), adj: make([][]int, m)}
	for i := range g.cost {
		g.cost[i] = edakit.Int(1, 9)
	}
	has := map[[2]int]bool{}
	add := func(a, b int) {
		if a == b || has[[2]int{a, b}] {
			return
		}
		has[[2]int{a, b}], has[[2]int{b, a}] = true, true
		g.adj[a] = append(g.adj[a], b)
		g.adj[b] = append(g.adj[b], a)
	}
	for i := 1; i < m; i++ {
		add(i, edakit.Int(0, i-1))
	}
	for e := 0; e < m/2+1; e++ {
		add(edakit.Int(0, m-1), edakit.Int(0, m-1))
	}
	g.src = pickDistinct(seq(m), 2)
	for s := 0; s < nsinks; s++ {
		g.sinks = append(g.sinks, pickDistinct(seq(m), 2))
	}
	return g
}

func seq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

// dijkstra returns the cheapest node-cost distance to every wire from the given start wires,
// where a start wire in free costs nothing, and the number of cheapest paths.
func (g *wgraph) dijkstra(start []int, free map[int]bool) ([]int, []int) {
	inf := 1 << 30
	d := make([]int, g.m)
	cnt := make([]int, g.m)
	done := make([]bool, g.m)
	for i := range d {
		d[i] = inf
	}
	for _, s := range start {
		c := g.cost[s]
		if free[s] {
			c = 0
		}
		if c < d[s] {
			d[s], cnt[s] = c, 1
		} else if c == d[s] {
			cnt[s]++
		}
	}
	for {
		u := -1
		for i := 0; i < g.m; i++ {
			if !done[i] && d[i] < inf && (u < 0 || d[i] < d[u]) {
				u = i
			}
		}
		if u < 0 {
			break
		}
		done[u] = true
		for _, v := range g.adj[u] {
			c := g.cost[v]
			if free[v] {
				c = 0
			}
			if d[u]+c < d[v] {
				d[v], cnt[v] = d[u]+c, cnt[u]
			} else if d[u]+c == d[v] {
				cnt[v] += cnt[u]
			}
		}
	}
	return d, cnt
}

func (g *wgraph) text() string {
	var cs, es []string
	for i, c := range g.cost {
		cs = append(cs, edakit.Sprintf("w%d: %d", i+1, c))
	}
	for a := 0; a < g.m; a++ {
		for _, b := range g.adj[a] {
			if a < b {
				es = append(es, edakit.Sprintf("w%d-w%d", a+1, b+1))
			}
		}
	}
	return "wire costs " + strings.Join(cs, ", ") + "; bidirectional switches " + strings.Join(es, ", ")
}

func wlist(xs []int) string {
	var s []string
	for _, x := range xs {
		s = append(s, edakit.Sprintf("w%d", x+1))
	}
	return strings.Join(s, " and ")
}

func mazeExpansion(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		g := newWGraph(edakit.Scaled(ctx, 6, 10), 1)
		d, _ := g.dijkstra(g.src, nil)
		best := minInt(d[g.sinks[0][0]], d[g.sinks[0][1]])
		q := edakit.Sprintf("A router expands a wavefront over an RR graph with %s. The source's output pin reaches %s; the sink's input pin is reached from %s. A path's cost is the sum of the costs of the wires it uses. What is the cost of the cheapest route?", g.text(), wlist(g.src), wlist(g.sinks[0]))
		return num(q, edakit.IntStr(best),
			edakit.Sprintf("Dijkstra from the source wires, adding each wire's cost when it is reached, gives cost %d to w%d and %d to w%d, so the cheapest route costs %d.", d[g.sinks[0][0]], g.sinks[0][0]+1, d[g.sinks[0][1]], g.sinks[0][1]+1, best),
			map[string]string{"cost": edakit.IntStr(best)})
	case 2:
		for {
			g := newWGraph(edakit.Scaled(ctx, 7, 10), 2)
			d, cnt := g.dijkstra(g.src, nil)
			t1 := g.sinks[0]
			end := t1[0]
			if d[t1[1]] < d[end] {
				end = t1[1]
			}
			if d[t1[0]] == d[t1[1]] || cnt[end] != 1 {
				continue
			}
			path := g.tracePath(d, end, g.src)
			if path == nil {
				continue
			}
			free := map[int]bool{}
			for _, p := range path {
				free[p] = true
			}
			d2, _ := g.dijkstra(path, free)
			c2 := minInt(d2[g.sinks[1][0]], d2[g.sinks[1][1]])
			if c2 == 0 {
				continue
			}
			tot := d[end] + c2
			q := edakit.Sprintf("A two-sink net is routed on an RR graph with %s. The source's output pin reaches %s. Sink 1 is reached from %s and sink 2 from %s. The router first finds the cheapest route to sink 1; then every wire already in the net's tree costs nothing, and it finds the cheapest branch from the tree to sink 2. Costs are sums of wire costs. What is the total cost of the routed tree?", g.text(), wlist(g.src), wlist(t1), wlist(g.sinks[1]))
			return num(q, edakit.IntStr(tot),
				edakit.Sprintf("The route to sink 1 uses %s at cost %d. Restarting the wavefront from that whole tree at zero cost, the cheapest branch to sink 2 adds %d. The tree costs %d, cheaper than routing the sinks independently whenever they share wires.", wlist(path), d[end], c2, tot),
				map[string]string{"first": edakit.IntStr(d[end]), "branch": edakit.IntStr(c2)})
		}
	default:
		for {
			alpha := edakit.Pick([]float64{1, 1.2, 1.5, 2})
			cmin := float64(edakit.Int(1, 4))
			type fr struct{ g, d float64 }
			var fs []fr
			var desc []string
			best, bi, second := math.Inf(1), 0, math.Inf(1)
			for i := 0; i < 3; i++ {
				f := fr{float64(edakit.Int(5, 40)), float64(edakit.Int(1, 12))}
				fs = append(fs, f)
				v := f.g + alpha*f.d*cmin
				desc = append(desc, edakit.Sprintf("node %d has path cost %s and a remaining distance of %s tiles", i+1, edakit.Num(f.g, 0), edakit.Num(f.d, 0)))
				if v < best {
					second, best, bi = best, v, i
				} else if v < second {
					second = v
				}
			}
			if second-best < 0.5 {
				continue
			}
			q := edakit.Sprintf("A directed-search router orders its priority queue by f = g + α·h, where g is the path cost so far and h = (tiles to the sink) × %s, the cheapest cost per tile. With α = %s, the queue holds: %s. Which node is expanded next? Answer with its number.", edakit.Num(cmin, 0), edakit.Num(alpha, 1), strings.Join(desc, "; "))
			var sc []string
			for i, f := range fs {
				sc = append(sc, edakit.Sprintf("node %d: %s", i+1, edakit.Num(f.g+alpha*f.d*cmin, 1)))
			}
			return num(q, edakit.IntStr(bi+1),
				edakit.Sprintf("f values: %s. The smallest is node %d. With α > 1 the estimate can overshoot, which speeds the search but can miss the cheapest path.", strings.Join(sc, ", "), bi+1),
				map[string]string{"winner": edakit.IntStr(bi + 1)})
		}
	}
}

// tracePath rebuilds the unique cheapest path that ends at end.
func (g *wgraph) tracePath(d []int, end int, start []int) []int {
	isStart := map[int]bool{}
	for _, s := range start {
		isStart[s] = true
	}
	path := []int{end}
	cur := end
	for steps := 0; steps < g.m+1; steps++ {
		if isStart[cur] && d[cur] == g.cost[cur] {
			return path
		}
		found := -1
		for _, p := range g.adj[cur] {
			if d[p]+g.cost[cur] == d[cur] {
				found = p
				break
			}
		}
		if found < 0 {
			return nil
		}
		cur = found
		path = append(path, cur)
	}
	return nil
}

type rnode struct {
	b, h      float64
	occ, cap_ int
}

func randNode() rnode {
	return rnode{float64(edakit.Int(1, 4)), float64(edakit.Int(0, 12)) / 4, edakit.Int(0, 3), edakit.Int(1, 2)}
}

func (n rnode) cost(pfac float64) float64 {
	over := n.occ + 1 - n.cap_
	if over < 0 {
		over = 0
	}
	return (n.b + n.h) * (1 + pfac*float64(over))
}

func (n rnode) text() string {
	return edakit.Sprintf("b = %s, h = %s, occupancy %d, capacity %d", edakit.Num(n.b, 0), edakit.Num(n.h, 2), n.occ, n.cap_)
}

const pfRule = "PathFinder prices a node as c_n = (b_n + h_n)·p_n, where b_n is the base cost, h_n the history cost and p_n = 1 + p_fac·max(0, occ_n + 1 - cap_n) the present-congestion cost of adding one more net."

func pathfinder(ctx generator.GeneratorContext) generator.Problem {
	pf := float64(edakit.Int(1, 8)) / 2
	switch whichKP(ctx) {
	case 1:
		n := randNode()
		c := n.cost(pf)
		q := pfRule + edakit.Sprintf(" With p_fac = %s, a node has %s. What is its cost to the net now being routed? Round to 2 decimal places.", edakit.Num(pf, 1), n.text())
		return num(q, edakit.Num(c, 2),
			edakit.Sprintf("p_n = 1 + %s × max(0, %d + 1 - %d) = %s, so c_n = (%s + %s) × %s = %s.", edakit.Num(pf, 1), n.occ, n.cap_, edakit.Num(c/(n.b+n.h), 2), edakit.Num(n.b, 0), edakit.Num(n.h, 2), edakit.Num(c/(n.b+n.h), 2), edakit.Num(c, 2)),
			map[string]string{"p": edakit.Num(c/(n.b+n.h), 2)})
	case 2:
		for {
			mk := func() ([]rnode, float64) {
				var ns []rnode
				t := 0.0
				for i := 0; i < edakit.Int(2, 3); i++ {
					n := randNode()
					ns = append(ns, n)
					t += n.cost(pf)
				}
				return ns, t
			}
			a, ca := mk()
			b, cb := mk()
			if math.Abs(ca-cb) < 0.05 {
				continue
			}
			txt := func(ns []rnode) string {
				var s []string
				for _, n := range ns {
					s = append(s, "("+n.text()+")")
				}
				return strings.Join(s, ", ")
			}
			best := math.Min(ca, cb)
			q := pfRule + edakit.Sprintf(" With p_fac = %s, a net can reach its sink through path A, whose nodes are %s, or through path B, whose nodes are %s. The router takes the cheaper path. What is that path's cost? Round to 2 decimal places.", edakit.Num(pf, 1), txt(a), txt(b))
			return num(q, edakit.Num(best, 2),
				edakit.Sprintf("Path A costs %s and path B costs %s, so the router pays %s. Congestion is negotiated through these prices, not forbidden outright.", edakit.Num(ca, 2), edakit.Num(cb, 2), edakit.Num(best, 2)),
				map[string]string{"path_a": edakit.Num(ca, 2), "path_b": edakit.Num(cb, 2)})
		}
	default:
		k := edakit.Scaled(ctx, 4, 9)
		over := 0
		var s []string
		for i := 0; i < k; i++ {
			o, c := edakit.Int(0, 4), edakit.Int(1, 2)
			if o > c {
				over += o - c
			}
			s = append(s, edakit.Sprintf("%d/%d", o, c))
		}
		q := edakit.Sprintf("After a PathFinder iteration the nodes that any net uses have occupancy/capacity %s. PathFinder stops when no node is overused. What is the total overuse, Σ max(0, occ - cap)?", strings.Join(s, ", "))
		return num(q, edakit.IntStr(over),
			edakit.Sprintf("Add occupancy minus capacity over the overused nodes: %d. The routing is legal only when this reaches 0.", over),
			map[string]string{"overuse": edakit.IntStr(over)})
	}
}

func history(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		hf := float64(edakit.Int(1, 8)) / 4
		r := edakit.Scaled(ctx, 3, 7)
		cp := edakit.Int(1, 2)
		h := 0.0
		var os []string
		for i := 0; i < r; i++ {
			o := edakit.Int(0, 4)
			os = append(os, edakit.IntStr(o))
			if o > cp {
				h += hf * float64(o-cp)
			}
		}
		q := edakit.Sprintf("A routing node has capacity %d and history cost h = 0 before the first iteration. After each iteration PathFinder updates h ← h + h_fac·max(0, occ - cap) with h_fac = %s. Over %d iterations the node's occupancy was %s. What is h now? Round to 2 decimal places.", cp, edakit.Num(hf, 2), r, strings.Join(os, ", "))
		return num(q, edakit.Num(h, 2),
			edakit.Sprintf("Only overused iterations add to h. Summing h_fac × overuse gives %s. History never decreases, so a chronically contested node stays expensive.", edakit.Num(h, 2)),
			map[string]string{"history": edakit.Num(h, 2)})
	case 2:
		p0 := float64(edakit.Int(1, 4)) / 2
		m := edakit.Pick([]float64{1.3, 1.5, 2})
		it := edakit.Int(2, 6)
		pf := p0 * math.Pow(m, float64(it-1))
		n := randNode()
		n.occ = n.cap_ + edakit.Int(0, 1)
		c := n.cost(pf)
		q := pfRule + edakit.Sprintf(" The present-cost factor starts at p_fac = %s in iteration 1 and is multiplied by %s after every iteration. In iteration %d a node has %s. What is its cost? Round to 2 decimal places.", edakit.Num(p0, 1), edakit.Num(m, 1), it, n.text())
		return num(q, edakit.Num(c, 2),
			edakit.Sprintf("p_fac = %s × %s^%d = %s. Then p_n = 1 + %s × %d = %s and c_n = (%s + %s) × %s = %s. Growing p_fac makes sharing steadily less tolerable.", edakit.Num(p0, 1), edakit.Num(m, 1), it-1, edakit.Num(pf, 4), edakit.Num(pf, 4), n.occ+1-n.cap_, edakit.Num(c/(n.b+n.h), 4), edakit.Num(n.b, 0), edakit.Num(n.h, 2), edakit.Num(c/(n.b+n.h), 4), edakit.Num(c, 2)),
			map[string]string{"p_fac": edakit.Num(pf, 4)})
	default:
		ba := float64(edakit.Int(1, 4))
		hf := float64(edakit.Int(1, 4)) / 2
		pf := float64(edakit.Int(1, 4)) / 2
		bb := ba*(1+pf) + float64(edakit.Int(1, 12))
		i := 1
		for (ba+hf*float64(i-1))*(1+pf) <= bb {
			i++
		}
		q := edakit.Sprintf("Two nets both want node A (capacity 1, base cost %s). In every iteration both use it, so it ends each iteration overused by 1, and its history grows by h_fac = %s per iteration starting from 0. The present-cost factor stays p_fac = %s, so a net that rips up and sees the other net still on A prices it at (%s + h)(1 + %s). Its detour costs %s and never congests. In which iteration (numbering from 1, with h = h_fac·(i - 1) in iteration i) does A first cost more than the detour?", edakit.Num(ba, 0), edakit.Num(hf, 1), edakit.Num(pf, 1), edakit.Num(ba, 0), edakit.Num(pf, 1), edakit.Num(bb, 2))
		return num(q, edakit.IntStr(i),
			edakit.Sprintf("In iteration %d, A costs (%s + %s)(1 + %s) = %s, above %s; one iteration earlier it cost %s. History is what breaks the tie for good: present cost alone, with p_fac fixed, would never change.", i, edakit.Num(ba, 0), edakit.Num(hf*float64(i-1), 2), edakit.Num(pf, 1), edakit.Num((ba+hf*float64(i-1))*(1+pf), 2), edakit.Num(bb, 2), edakit.Num((ba+hf*float64(i-2))*(1+pf), 2)),
			map[string]string{"iteration": edakit.IntStr(i)})
	}
}

func tdCrit(slack, dmax float64) float64 {
	c := 0.99 - slack/dmax
	if c < 0 {
		c = 0
	}
	return c
}

const tdRule = "The timing-driven router uses Crit = max(0, 0.99 - slack / D_max) and prices a node as Crit·delay(n) + (1 - Crit)·b(n)·h(n)·p(n)."

func timingDriven(ctx generator.GeneratorContext) generator.Problem {
	dmax := float64(10 * edakit.Int(100, 400))
	switch whichKP(ctx) {
	case 1:
		sl := float64(10 * edakit.Int(0, int(dmax)/10))
		c := tdCrit(sl, dmax)
		q := edakit.Sprintf("A timing-driven router assigns a connection the criticality Crit = max(0, 0.99 - slack / D_max). The critical path is D_max = %s ps and the connection's slack is %s ps. What is Crit? Round to 3 decimal places.", edakit.Num(dmax, 0), edakit.Num(sl, 0))
		return num(q, edakit.Num(c, 3),
			edakit.Sprintf("0.99 - %s / %s = %s, floored at 0. Capping at 0.99 keeps a little congestion pressure even on the most critical connection, so it can still be moved off an overused wire.", edakit.Num(sl, 0), edakit.Num(dmax, 0), edakit.Num(c, 3)),
			map[string]string{"crit": edakit.Num(c, 3)})
	case 2:
		sl := float64(10 * edakit.Int(0, int(dmax)/10))
		c := tdCrit(sl, dmax)
		d := float64(edakit.Int(1, 20)) / 4
		b := float64(edakit.Int(1, 4))
		h := 1 + float64(edakit.Int(0, 8))/4
		p := 1 + float64(edakit.Int(0, 8))/2
		v := c*d + (1-c)*b*h*p
		q := tdRule + edakit.Sprintf(" Delays are normalized to the same scale as the congestion term. With D_max = %s ps and slack %s ps, a node has delay %s, b = %s, h = %s and p = %s. What is its cost? Round to 3 decimal places.", edakit.Num(dmax, 0), edakit.Num(sl, 0), edakit.Num(d, 2), edakit.Num(b, 0), edakit.Num(h, 2), edakit.Num(p, 1))
		return num(q, edakit.Num(v, 3),
			edakit.Sprintf("Crit = %s. Cost = %s × %s + %s × %s × %s × %s = %s.", edakit.Num(c, 4), edakit.Num(c, 4), edakit.Num(d, 2), edakit.Num(1-c, 4), edakit.Num(b, 0), edakit.Num(h, 2), edakit.Num(p, 1), edakit.Num(v, 3)),
			map[string]string{"crit": edakit.Num(c, 4)})
	default:
		for {
			sl := float64(10 * edakit.Int(0, int(dmax)/10))
			c := tdCrit(sl, dmax)
			da, ga := float64(edakit.Int(2, 6)), float64(edakit.Int(6, 20))
			db, gb := da+float64(edakit.Int(2, 10)), float64(edakit.Int(1, 5))
			ca := c*da + (1-c)*ga
			cb := c*db + (1-c)*gb
			if math.Abs(ca-cb) < 0.2 {
				continue
			}
			ch := 1
			if cb < ca {
				ch = 2
			}
			q := tdRule + edakit.Sprintf(" Delays are normalized to the scale of the congestion term. A connection has slack %s ps with D_max = %s ps. Path 1 has total delay %s and congestion cost Σ b·h·p = %s; path 2 has delay %s and congestion cost %s. Which path does the router choose? Answer 1 or 2.", edakit.Num(sl, 0), edakit.Num(dmax, 0), edakit.Num(da, 0), edakit.Num(ga, 0), edakit.Num(db, 0), edakit.Num(gb, 0))
			return num(q, edakit.IntStr(ch),
				edakit.Sprintf("Crit = %s. Path 1 costs %s × %s + %s × %s = %s; path 2 costs %s × %s + %s × %s = %s. Critical connections buy the fast, congested path; non-critical ones detour and free it.", edakit.Num(c, 3), edakit.Num(c, 3), edakit.Num(da, 0), edakit.Num(1-c, 3), edakit.Num(ga, 0), edakit.Num(ca, 3), edakit.Num(c, 3), edakit.Num(db, 0), edakit.Num(1-c, 3), edakit.Num(gb, 0), edakit.Num(cb, 3)),
				map[string]string{"cost1": edakit.Num(ca, 3), "cost2": edakit.Num(cb, 3)})
		}
	}
}

func minChannelWidth(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		rows := edakit.Int(2, 3)
		best := 0
		var lines []string
		var dens []int
		for r := 0; r < rows; r++ {
			cols := 10
			cnt := make([]int, cols+1)
			var iv []string
			for k := 0; k < edakit.Scaled(ctx, 3, 7); k++ {
				a := edakit.Int(1, cols-1)
				b := edakit.Int(a+1, cols)
				iv = append(iv, edakit.Sprintf("%d-%d", a, b))
				for c := a; c <= b; c++ {
					cnt[c]++
				}
			}
			d := maxInt(cnt...)
			dens = append(dens, d)
			if d > best {
				best = d
			}
			lines = append(lines, edakit.Sprintf("channel %d carries nets spanning columns %s", r+1, strings.Join(iv, ", ")))
		}
		q := "After global routing, the horizontal channels have this demand: " + strings.Join(lines, "; ") + ". Spans are inclusive, and two nets crossing the same column of a channel need different tracks. What is the smallest channel width W that could hold this demand?"
		return num(q, edakit.IntStr(best),
			edakit.Sprintf("Each channel needs as many tracks as its densest column: %s. W must cover the worst channel, %d. A real detailed router may need more, because switch-box flexibility is limited.", joinInts(dens), best),
			map[string]string{"density": joinInts(dens)})
	case 2:
		wmin := edakit.Scaled(ctx, 6, 120)
		w0 := edakit.Pick([]int{2, 4, 8, 12})
		for w0 >= wmin {
			w0 /= 2
		}
		if w0 < 1 {
			w0 = 1
		}
		att := 0
		w := w0
		lo := 0
		for {
			att++
			if w >= wmin {
				break
			}
			lo = w
			w *= 2
		}
		hi := w
		for hi-lo > 1 {
			mid := (lo + hi) / 2
			att++
			if mid >= wmin {
				hi = mid
			} else {
				lo = mid
			}
		}
		q := edakit.Sprintf("A router finds the minimum channel width by trying W = %d, doubling W after every failure until a route succeeds, then binary searching between the last failing and first succeeding widths (try mid = floor((fail + success) / 2), and stop when they differ by 1). The design routes exactly when W ≥ %d. How many routing attempts does the search make in total?", w0, wmin)
		return num(q, edakit.IntStr(att),
			edakit.Sprintf("The doubling phase and the binary search together take %d complete routing runs, which is why minimum-channel-width experiments dominate VPR run time.", att),
			map[string]string{"attempts": edakit.IntStr(att)})
	default:
		wmin := edakit.Scaled(ctx, 10, 200)
		w := (13*wmin + 9) / 10
		if w%2 == 1 {
			w++
		}
		q := edakit.Sprintf("Architecture studies report delay at a low-stress channel width of 1.3 times the minimum, rounded up, and this fabric uses unidirectional wires, which come in pairs, so W must be even. The minimum channel width found is %d. What channel width is used for the low-stress run?", wmin)
		return num(q, edakit.IntStr(w),
			edakit.Sprintf("1.3 × %d = %s, rounded up to %d, then to the next even value: %d. Routing at exactly W_min would make the router struggle and inflate critical-path delay.", wmin, edakit.Num(1.3*float64(wmin), 1), (13*wmin+9)/10, w),
			map[string]string{"raw": edakit.IntStr((13*wmin + 9) / 10)})
	}
}
