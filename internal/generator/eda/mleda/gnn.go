package mleda

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// netSizes draws n net sizes (pin counts) for a small netlist.
func netSizes(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = edakit.Pick([]int{2, 2, 2, 3, 3, 4, 5})
	}
	return out
}

// netlistGraph: clique or star edges of one net (KP 1), clique edges of a netlist with a
// high-fanout net (KP 2), and the clique-to-bipartite edge ratio (KP 3).
func netlistGraph(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		p := edakit.Scaled(ctx, 3, 40)
		if edakit.Int(0, 1) == 0 {
			e := p * (p - 1) / 2
			return edakit.Numeric(edakit.Sprintf("A net connects %d pins (one driver and %d sinks). Converting the netlist hypergraph to an ordinary graph with the clique model, how many edges does this net contribute?", p, p-1),
				edakit.IntStr(e), edakit.Sprintf("The clique model joins every pair of pins: $\\binom{%d}{2} = %d \\cdot %d / 2 = %d$ edges.", p, p, p-1, e), map[string]string{"pins": edakit.IntStr(p)})
		}
		return edakit.Numeric(edakit.Sprintf("A net connects %d pins (one driver and %d sinks). Converting the hypergraph with the star model, which adds one new node for the net and joins it to every pin, how many edges does this net contribute?", p, p-1),
			edakit.IntStr(p), edakit.Sprintf("The star model adds a net node and one edge per pin: %d edges, linear rather than quadratic in the pin count.", p), map[string]string{"pins": edakit.IntStr(p)})
	case 2:
		sizes := netSizes(edakit.Scaled(ctx, 4, 8))
		big := edakit.Pick([]int{16, 24, 32, 48, 64})
		sizes = append(sizes, big)
		sizes = edakit.Shuffle(sizes)
		tot := 0
		var terms []string
		for _, p := range sizes {
			tot += p * (p - 1) / 2
			terms = append(terms, edakit.IntStr(p*(p-1)/2))
		}
		return edakit.Numeric(edakit.Sprintf("A small netlist has nets with these pin counts: %s. Using the clique model and counting each net's edges separately, how many edges does the graph have?", ints(sizes)),
			edakit.IntStr(tot), edakit.Sprintf("Each net of $p$ pins gives $p(p-1)/2$ edges: %s, total %d. The %d-pin net alone gives %d.", strings.Join(terms, " + "), tot, big, big*(big-1)/2),
			map[string]string{"total": edakit.IntStr(tot)})
	default:
		sizes := netSizes(edakit.Scaled(ctx, 4, 8))
		big := edakit.Pick([]int{10, 12, 16, 20, 30})
		sizes = edakit.Shuffle(append(sizes, big))
		clique, pins := 0, 0
		for _, p := range sizes {
			clique += p * (p - 1) / 2
			pins += p
		}
		r := float64(clique) / float64(pins)
		ans := edakit.Num(r, 2)
		return edakit.Numeric(edakit.Sprintf("A netlist has nets with pin counts %s. A GNN can use the clique model (cells only, every pin pair joined) or a bipartite cell-net graph (one node per cell and per net, one edge per pin). What is the ratio of clique edges to bipartite edges? Round to 2 decimal places.", ints(sizes)),
			ans, edakit.Sprintf("Clique edges $= \\sum p(p-1)/2 = %d$. Bipartite edges $= \\sum p = %d$. Ratio $= %s$. The bipartite form keeps the net as an object and grows linearly with pins.", clique, pins, ans),
			map[string]string{"clique": edakit.IntStr(clique), "pins": edakit.IntStr(pins)})
	}
}

// tinyGraph draws a connected undirected graph on n nodes as an adjacency list.
func tinyGraph(n, extra int) [][]int {
	adj := make([][]int, n)
	has := map[[2]int]bool{}
	add := func(a, b int) {
		if a == b {
			return
		}
		if a > b {
			a, b = b, a
		}
		if has[[2]int{a, b}] {
			return
		}
		has[[2]int{a, b}] = true
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	for i := 1; i < n; i++ {
		add(i, edakit.Int(0, i-1))
	}
	for k := 0; k < extra; k++ {
		add(edakit.Int(0, n-1), edakit.Int(0, n-1))
	}
	return adj
}

// edgeList renders the edges of adj with node names prefixed by name.
func edgeList(adj [][]int, name string) string {
	var es []string
	for a := range adj {
		for _, b := range adj[a] {
			if a < b {
				es = append(es, edakit.Sprintf("%s%d–%s%d", name, a+1, name, b+1))
			}
		}
	}
	return strings.Join(es, ", ")
}

// maxDegreeNode returns a node of largest degree.
func maxDegreeNode(adj [][]int) int {
	best := 0
	for i := range adj {
		if len(adj[i]) > len(adj[best]) {
			best = i
		}
	}
	return best
}

// messagePassing: one layer with sum aggregation (KP 1), mean aggregation (KP 2), and GCN
// symmetric normalization with self loops (KP 3).
func messagePassing(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 7)
	adj := tinyGraph(n, edakit.Int(1, 3))
	h := make([]float64, n)
	var hs []string
	for i := range h {
		h[i] = float64(edakit.Int(-3, 9))
		hs = append(hs, edakit.Sprintf("$h_{%d} = %s$", i+1, f(h[i], 0)))
	}
	v := maxDegreeNode(adj)
	if edakit.Int(0, 1) == 0 {
		v = edakit.Int(0, n-1)
		for tries := 0; len(adj[v]) < 2 && tries < n; tries++ {
			v = (v + 1) % n
		}
		if len(adj[v]) < 2 {
			v = maxDegreeNode(adj)
		}
	}
	stem := edakit.Sprintf("A netlist graph has nodes 1 to %d and undirected edges %s. Each node carries a scalar feature: %s.", n, edgeList(adj, ""), strings.Join(hs, ", "))
	sum := 0.0
	for _, u := range adj[v] {
		sum += h[u]
	}
	deg := float64(len(adj[v]))
	switch kp(ctx) {
	case 1:
		a := float64(edakit.Int(1, 3))
		b := float64(edakit.Int(1, 3))
		out := math.Max(0, a*h[v]+b*sum)
		ans := edakit.Num(out, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" One message-passing layer computes $h_v' = \\mathrm{ReLU}(%s\\,h_v + %s\\sum_{u \\in N(v)} h_u)$. What is $h_{%d}'$? Round to 2 decimal places.", f(a, 0), f(b, 0), v+1),
			ans, edakit.Sprintf("Neighbors of node %d: %s, with feature sum %s. $\\mathrm{ReLU}(%s \\cdot %s + %s \\cdot %s) = %s$.", v+1, ints(plus1(adj[v])), f(sum, 0), f(a, 0), pn(h[v], 0), f(b, 0), pn(sum, 0), ans),
			map[string]string{"agg": f(sum, 0)})
	case 2:
		a := float64(edakit.Pick([]int{-2, -1, 1, 2}))
		b := float64(edakit.Pick([]int{-2, -1, 1, 2, 3}))
		mean := sum / deg
		out := math.Max(0, a*h[v]+b*mean)
		ans := edakit.Num(out, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" One layer computes $h_v' = \\mathrm{ReLU}(%s\\,h_v + %s \\cdot \\mathrm{mean}_{u \\in N(v)} h_u)$. What is $h_{%d}'$? Round to 2 decimal places.", pn(a, 0), pn(b, 0), v+1),
			ans, edakit.Sprintf("Node %d has %d neighbors with sum %s, so the mean is %s. Pre-activation $%s \\cdot %s + %s \\cdot %s = %s$, and ReLU gives %s.", v+1, len(adj[v]), f(sum, 0), f(mean, 4), pn(a, 0), pn(h[v], 0), pn(b, 0), pn(mean, 4), f(a*h[v]+b*mean, 4), ans),
			map[string]string{"mean": f(mean, 4)})
	default:
		w := float64(edakit.Int(1, 3))
		dv := deg + 1
		tot := h[v] / dv
		var terms []string
		terms = append(terms, edakit.Sprintf("%s/%s", pn(h[v], 0), f(dv, 0)))
		for _, u := range adj[v] {
			du := float64(len(adj[u])) + 1
			tot += h[u] / math.Sqrt(du*dv)
			terms = append(terms, edakit.Sprintf("%s/\\sqrt{%s}", pn(h[u], 0), f(du*dv, 0)))
		}
		out := w * tot
		ans := edakit.Num(out, 3)
		return edakit.Numeric(stem+edakit.Sprintf(" A GCN layer with self loops computes $h_v' = %s\\sum_{u \\in N(v) \\cup \\{v\\}} h_u / \\sqrt{\\tilde{d}_u \\tilde{d}_v}$, where $\\tilde{d}$ is the degree plus one. No activation is applied. What is $h_{%d}'$? Round to 3 decimal places.", f(w, 0), v+1),
			ans, edakit.Sprintf("$\\tilde{d}_{%d} = %s$. The sum is $%s = %s$, times %s gives %s. High-degree neighbors are damped by their larger $\\tilde{d}_u$.", v+1, f(dv, 0), strings.Join(terms, " + "), f(tot, 4), f(w, 0), ans),
			map[string]string{"dv": f(dv, 0)})
	}
}

// plus1 returns the 1-based names of 0-based node indices.
func plus1(xs []int) []int {
	out := make([]int, len(xs))
	for i, x := range xs {
		out[i] = x + 1
	}
	return out
}

// bfsWithin counts nodes within k hops of s, s included.
func bfsWithin(adj [][]int, s, k int) int {
	dist := map[int]int{s: 0}
	q := []int{s}
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if dist[x] == k {
			continue
		}
		for _, y := range adj[x] {
			if _, ok := dist[y]; !ok {
				dist[y] = dist[x] + 1
				q = append(q, y)
			}
		}
	}
	return len(dist)
}

// netlistGNN: receptive field of a k-layer GNN (KP 1), layers needed to span a timing path in
// a bipartite cell-net graph (KP 2), and the parameter count, which is independent of graph
// size (KP 3).
func netlistGNN(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 7, 10)
		adj := tinyGraph(n, edakit.Int(0, 2))
		s := edakit.Int(0, n-1)
		k := edakit.Int(1, 3)
		if k == 3 && n < 9 {
			k = 2
		}
		c := bfsWithin(adj, s, k)
		return edakit.Numeric(edakit.Sprintf("A graph has nodes 1 to %d and undirected edges %s. A GNN with %d message-passing %s computes an embedding for node %d. How many nodes, node %d included, can influence that embedding?", n, edgeList(adj, ""), k, map[bool]string{true: "layer", false: "layers"}[k == 1], s+1, s+1),
			edakit.IntStr(c), edakit.Sprintf("After $k$ layers a node sees exactly its $k$-hop neighborhood. A breadth-first search from node %d to depth %d reaches %d nodes.", s+1, k, c),
			map[string]string{"hops": edakit.IntStr(k)})
	case 2:
		l := edakit.Scaled(ctx, 3, 30)
		hops := 2 * (l + 1)
		return edakit.Numeric(edakit.Sprintf("A timing path runs from a launching flip-flop through %d combinational gates to a capturing flip-flop, one net between each consecutive pair of cells. The netlist is encoded as a bipartite graph with a node per cell and a node per net. How many message-passing layers does a GNN need before the capturing flop's embedding can depend on the launching flop?", l),
			edakit.IntStr(hops), edakit.Sprintf("The path has %d cells and %d nets between them. Each net costs two hops (cell to net, net to cell), so the distance is $2 \\times %d = %d$ layers. Deep logic needs deep GNNs or a different encoding, which is why timing GNNs often propagate along the timing graph in topological order instead.", l+2, l+1, l+1, hops),
			map[string]string{"nets": edakit.IntStr(l + 1)})
	default:
		din := edakit.Pick([]int{8, 12, 16, 24, 32})
		d := edakit.Pick([]int{32, 64, 128})
		layers := edakit.Int(3, 5)
		cells := edakit.Pick([]string{"200 thousand", "2 million", "15 million"})
		p := (2*din*d + d) + (layers-1)*(2*d*d+d) + (d + 1)
		return edakit.Numeric(edakit.Sprintf("A GNN for a %s-cell netlist has %d layers. Each layer computes $h' = \\sigma(W_{\\mathrm{self}} h_v + W_{\\mathrm{nb}} \\cdot \\mathrm{mean}_{u} h_u + b)$. Input features have dimension %d and every hidden layer has width %d. A linear readout with bias maps the final embedding to one number. How many trainable parameters does the model have?", cells, layers, din, d),
			edakit.IntStr(p), edakit.Sprintf("Layer 1: $2 \\cdot %d \\cdot %d + %d = %d$. Each of the remaining %d layers: $2 \\cdot %d^2 + %d = %d$. Readout: $%d + 1$. Total %d. The cell count does not appear: the same weights apply to any netlist, which is what makes the model inductive.", din, d, d, 2*din*d+d, layers-1, d, d, 2*d*d+d, d, p),
			map[string]string{"params": edakit.IntStr(p)})
	}
}

// heteroGraphs: parameters of a relational layer (KP 1), per-type input encoders (KP 2), and
// typed edge counts with reverse relations (KP 3).
func heteroGraphs(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		r := edakit.Int(2, 6)
		d := edakit.Pick([]int{16, 32, 64, 128})
		l := edakit.Int(2, 4)
		p := l * ((r+1)*d*d + d)
		return edakit.Numeric(edakit.Sprintf("A heterogeneous netlist GNN uses %d edge relations (for example cell-drives-net, net-feeds-cell and their reverses). Each layer has one $%d \\times %d$ weight matrix per relation, one self-loop matrix of the same size and a bias vector of length %d. How many parameters do %d such layers have?", r, d, d, d, l),
			edakit.IntStr(p), edakit.Sprintf("Per layer: $(%d + 1) \\cdot %d^2 + %d = %d$. Times %d layers: %d. Parameters grow linearly with the number of relations, which is why relation-heavy schemas use basis decomposition.", r, d, d, (r+1)*d*d+d, l, p),
			map[string]string{"per_layer": edakit.IntStr((r+1)*d*d + d)})
	case 2:
		d := edakit.Pick([]int{32, 64, 128})
		fc := edakit.Int(8, 40)
		fn := edakit.Int(3, 12)
		fp := edakit.Int(2, 10)
		p := fc*d + d + fn*d + d + fp*d + d
		return edakit.Numeric(edakit.Sprintf("A heterogeneous graph has three node types with different raw features: cells (%d features, such as one-hot cell type, drive strength and area), nets (%d features, such as fanout and bounding box) and pins (%d features). Each type has its own linear encoder with bias into a shared width of %d. How many parameters do the three encoders have?", fc, fn, fp, d),
			edakit.IntStr(p), edakit.Sprintf("Each encoder has $f_t \\cdot %d + %d$ parameters: $(%d + %d + %d) \\cdot %d + 3 \\cdot %d = %d$. Separate encoders let each type keep its own feature space.", d, d, fc, fn, fp, d, d, p),
			map[string]string{"features": edakit.IntStr(fc + fn + fp)})
	default:
		n := edakit.Scaled(ctx, 4, 9)
		var fan []int
		sinks := 0
		for i := 0; i < n; i++ {
			x := edakit.Int(1, 6)
			fan = append(fan, x)
			sinks += x
		}
		tot := 2 * (n + sinks)
		return edakit.Numeric(edakit.Sprintf("A netlist has %d nets with fanouts (sink counts) %s; each net also has one driver. The heterogeneous graph has relation drives (driver cell to net) and relation feeds (net to sink cell), and a reverse relation is added for each so messages flow both ways. How many edges does the graph have in total?", n, ints(fan)),
			edakit.IntStr(tot), edakit.Sprintf("drives: %d edges, one per net. feeds: %d edges, one per sink. Reverses double the total: $2(%d + %d) = %d$.", n, sinks, n, sinks, tot),
			map[string]string{"sinks": edakit.IntStr(sinks)})
	}
}

// fanoutScale: messages per layer in a bipartite netlist graph (KP 1), the sampled
// computation-tree bound (KP 2), and the clique-edge share of one high-fanout net (KP 3).
func fanoutScale(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		nets := edakit.Scaled(ctx, 1000, 50000)
		avg := float64(edakit.Int(25, 40)) / 10
		clk := edakit.Pick([]int{2000, 5000, 10000, 20000})
		pins := int(math.Round(float64(nets)*avg)) + clk
		msgs := 2 * pins
		return edakit.Numeric(edakit.Sprintf("A block has %d signal nets averaging %s pins each, plus one clock net with %d pins. It is encoded as a bipartite cell-net graph and messages flow both ways along every edge. How many messages does one full-graph message-passing layer send?", nets, f(avg, 1), clk),
			edakit.IntStr(msgs), edakit.Sprintf("Pins $= %d \\times %s + %d = %d$. Each pin is one edge, and each edge carries two messages: %d.", nets, f(avg, 1), clk, pins, msgs),
			map[string]string{"pins": edakit.IntStr(pins)})
	case 2:
		s := edakit.Int(5, 25)
		l := edakit.Int(2, 4)
		tot, pw := 0, 1
		for i := 0; i <= l; i++ {
			tot += pw
			pw *= s
		}
		return edakit.Numeric(edakit.Sprintf("To train on a large netlist in mini-batches, each node samples at most %d neighbors per layer. For a %d-layer GNN, what is the largest number of nodes in the sampled computation tree of one target node (target included)?", s, l),
			edakit.IntStr(tot), edakit.Sprintf("Level $i$ has at most $%d^i$ nodes: $\\sum_{i=0}^{%d} %d^i = %d$. Without the cap, one clock-net neighbor alone could pull in thousands of nodes.", s, l, s, tot),
			map[string]string{"layers": edakit.IntStr(l)})
	default:
		sizes := netSizes(edakit.Scaled(ctx, 30, 120))
		big := edakit.Pick([]int{40, 64, 100, 128, 200})
		small := 0
		for _, p := range sizes {
			small += p * (p - 1) / 2
		}
		bigE := big * (big - 1) / 2
		share := 100 * float64(bigE) / float64(bigE+small)
		ans := edakit.Num(share, 2)
		return edakit.Numeric(edakit.Sprintf("A netlist has %d small nets whose clique expansions contribute %d edges in total, plus one reset net with %d pins. What percentage of all clique-model edges does the reset net contribute? Round to 2 decimal places.", len(sizes), small, big),
			ans, edakit.Sprintf("The reset net gives $%d \\cdot %d / 2 = %d$ edges. Share $= %d / (%d + %d) = %s$ percent. One net dominates the graph, which is why GNNs drop, cap or star-expand high-fanout nets.", big, big-1, bigE, bigE, bigE, small, ans),
			map[string]string{"big_edges": edakit.IntStr(bigE)})
	}
}
