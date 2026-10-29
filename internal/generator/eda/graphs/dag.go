package graphs

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// rootedTraversal draws a rooted tree by parent pointers and asks for a pre-order or post-order
// position, or a post-order subtree sum (downstream capacitance).
func rootedTraversal(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 10)
	parent := make([]int, n)
	children := make([][]int, n)
	parent[0] = -1
	var ptxt []string
	for v := 1; v < n; v++ {
		lo := 0
		if v > 4 {
			lo = v - 4
		}
		parent[v] = edakit.Int(lo, v-1)
		children[parent[v]] = append(children[parent[v]], v)
		ptxt = append(ptxt, edakit.Sprintf("parent(%d) = %d", v+1, parent[v]+1))
	}
	var pre, post []int
	var walk func(v int)
	walk = func(v int) {
		pre = append(pre, v)
		for _, c := range children[v] {
			walk(c)
		}
		post = append(post, v)
	}
	walk(0)
	intro := "A rooted tree has nodes 1 to " + edakit.IntStr(n) + " with root 1 and " + strings.Join(ptxt, ", ") + ". Children are visited in increasing label order. "
	names := func(xs []int) string {
		var s []string
		for _, x := range xs {
			s = append(s, edakit.IntStr(x+1))
		}
		return strings.Join(s, ", ")
	}
	switch kpOf(ctx) {
	case 1:
		k := edakit.Int(3, n)
		return edakit.Numeric(intro+edakit.Sprintf("Which node is visited in position %d of the pre-order traversal? Answer with its label.", k),
			edakit.IntStr(pre[k-1]+1),
			"Pre-order visits a node before its children: "+names(pre)+". Position "+edakit.IntStr(k)+" holds node "+edakit.IntStr(pre[k-1]+1)+".",
			map[string]string{"preorder": names(pre)})
	case 2:
		k := edakit.Int(2, n-1)
		return edakit.Numeric(intro+edakit.Sprintf("Which node is visited in position %d of the post-order traversal? Answer with its label.", k),
			edakit.IntStr(post[k-1]+1),
			"Post-order visits a node after all of its children: "+names(post)+". Position "+edakit.IntStr(k)+" holds node "+edakit.IntStr(post[k-1]+1)+".",
			map[string]string{"postorder": names(post)})
	default:
		cap := make([]int, n)
		var ctxt []string
		for v := 0; v < n; v++ {
			cap[v] = edakit.Int(1, 9)
			ctxt = append(ctxt, edakit.Sprintf("C%d = %d", v+1, cap[v]))
		}
		down := make([]int, n)
		for _, v := range post {
			down[v] = cap[v]
			for _, c := range children[v] {
				down[v] += down[c]
			}
		}
		var inner []int
		for v := 1; v < n; v++ {
			if len(children[v]) > 0 {
				inner = append(inner, v)
			}
		}
		if len(inner) == 0 {
			inner = []int{0}
		}
		v := edakit.Pick(inner)
		return edakit.Numeric(intro+"Each node carries a capacitance in fF: "+strings.Join(ctxt, ", ")+". A post-order pass computes the downstream capacitance of every node (its own plus everything below it). What is the downstream capacitance of node "+edakit.IntStr(v+1)+", in fF?",
			edakit.IntStr(down[v]),
			"Post-order guarantees every child is finished before its parent, so each node adds its own value to its children's totals. Node "+edakit.IntStr(v+1)+" collects "+edakit.IntStr(down[v])+" fF.",
			map[string]string{"downstream": edakit.IntStr(down[v])})
	}
}

func topoSortGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	g := randDAGSources(n, edakit.Int(1, 2), edakit.Scaled(ctx, 2, 5), 1)
	// Relabel nodes randomly so that index order is not itself the answer.
	perm := edakit.Shuffle(seq(n))
	rg := dag{n: n, adj: make([][]int, n), d: map[[2]int]int{}}
	for u := 0; u < n; u++ {
		for _, v := range g.adj[u] {
			rg.adj[perm[u]] = append(rg.adj[perm[u]], perm[v])
		}
	}
	indeg := make([]int, n)
	for u := 0; u < n; u++ {
		for _, v := range rg.adj[u] {
			indeg[v]++
		}
	}
	intro := "A combinational netlist graph has nodes v1 to v" + edakit.IntStr(n) + " and arcs " + plainArcs(rg) + ". "
	switch kpOf(ctx) {
	case 1, 2:
		// Kahn's algorithm with the smallest-label ready node first.
		deg := append([]int(nil), indeg...)
		done := make([]bool, n)
		var order []int
		for len(order) < n {
			for v := 0; v < n; v++ {
				if !done[v] && deg[v] == 0 {
					done[v] = true
					order = append(order, v)
					for _, w := range rg.adj[v] {
						deg[w]--
					}
					break
				}
			}
		}
		k := edakit.Int(2, n-1)
		var on []string
		for _, v := range order {
			on = append(on, vName(v))
		}
		return edakit.Numeric(intro+edakit.Sprintf("Run Kahn's algorithm, always removing the ready node (in-degree 0) with the smallest label. Which node is output in position %d? Answer with its number (for v3, answer 3).", k),
			edakit.IntStr(order[k-1]+1),
			"Kahn's algorithm repeatedly outputs a node with no remaining predecessors and deletes its arcs. With smallest-label tie-breaking the order is "+strings.Join(on, ", ")+".",
			map[string]string{"order": strings.Join(on, ",")})
	default:
		anc := make([]map[int]bool, n)
		desc := make([]map[int]bool, n)
		for v := 0; v < n; v++ {
			anc[v], desc[v] = map[int]bool{}, map[int]bool{}
		}
		var dfs func(root, u int)
		dfs = func(root, u int) {
			for _, w := range rg.adj[u] {
				if !desc[root][w] {
					desc[root][w] = true
					anc[w][root] = true
					dfs(root, w)
				}
			}
		}
		for v := 0; v < n; v++ {
			dfs(v, v)
		}
		early := edakit.Int(0, 1) == 0
		v := edakit.Int(0, n-1)
		for tries := 0; tries < 30; tries++ {
			if (early && len(anc[v]) > 0) || (!early && len(desc[v]) > 0) {
				break
			}
			v = edakit.Int(0, n-1)
		}
		if early {
			ans := len(anc[v]) + 1
			return edakit.Numeric(intro+"Over all valid topological orders, what is the earliest position (counting from 1) that "+vName(v)+" can occupy?",
				edakit.IntStr(ans),
				vName(v)+" must follow every one of its ancestors, and there are "+edakit.IntStr(len(anc[v]))+" of them. Placing them first, then "+vName(v)+", is always valid, so the earliest position is "+edakit.IntStr(ans)+".",
				map[string]string{"ancestors": edakit.IntStr(len(anc[v]))})
		}
		ans := n - len(desc[v])
		return edakit.Numeric(intro+"Over all valid topological orders, what is the latest position (counting from 1) that "+vName(v)+" can occupy?",
			edakit.IntStr(ans),
			vName(v)+" must precede its "+edakit.IntStr(len(desc[v]))+" descendants, so at best it sits just before them: position "+edakit.IntStr(n)+" - "+edakit.IntStr(len(desc[v]))+" = "+edakit.IntStr(ans)+".",
			map[string]string{"descendants": edakit.IntStr(len(desc[v]))})
	}
}

func plainArcs(g dag) string {
	var parts []string
	for u := 0; u < g.n; u++ {
		for _, v := range g.adj[u] {
			parts = append(parts, vName(u)+"→"+vName(v))
		}
	}
	return strings.Join(parts, ", ")
}

func levelizeGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 10)
	g := randDAGSources(n, edakit.Int(2, 3), edakit.Scaled(ctx, 2, 6), 1)
	lvl := make([]int, n)
	for u := 0; u < n; u++ {
		for _, v := range g.adj[u] {
			if lvl[u]+1 > lvl[v] {
				lvl[v] = lvl[u] + 1
			}
		}
	}
	maxL := 0
	for _, l := range lvl {
		maxL = max(maxL, l)
	}
	intro := "A combinational netlist has nodes v1 to v" + edakit.IntStr(n) + " and arcs " + plainArcs(g) + ". Primary inputs (nodes with no incoming arc) are at level 0, and every other node sits one level above its highest fanin. "
	var lt []string
	for v := 0; v < n; v++ {
		lt = append(lt, edakit.Sprintf("%s:%d", vName(v), lvl[v]))
	}
	expl := "Levels in topological order: " + strings.Join(lt, ", ") + "."
	switch kpOf(ctx) {
	case 1:
		v := edakit.Int(n/2, n-1)
		return edakit.Numeric(intro+"What is the level of "+vName(v)+"?", edakit.IntStr(lvl[v]), expl, map[string]string{"level": edakit.IntStr(lvl[v])})
	case 2:
		return edakit.Numeric(intro+"How many levels does the netlist have, counting level 0?", edakit.IntStr(maxL+1), expl+" The deepest level is "+edakit.IntStr(maxL)+".", map[string]string{"levels": edakit.IntStr(maxL + 1)})
	default:
		L := edakit.Int(1, maxL)
		cnt := 0
		for _, l := range lvl {
			if l == L {
				cnt++
			}
		}
		return edakit.Numeric(intro+edakit.Sprintf("A levelized simulator evaluates all nodes of one level together. How many nodes are at level %d?", L),
			edakit.IntStr(cnt), expl, map[string]string{"width": edakit.IntStr(cnt)})
	}
}

// sccs returns the component id of every node (Kosaraju) and the component count.
func sccs(n int, adj [][]int) ([]int, int) {
	seen := make([]bool, n)
	var order []int
	var dfs1 func(u int)
	dfs1 = func(u int) {
		seen[u] = true
		for _, v := range adj[u] {
			if !seen[v] {
				dfs1(v)
			}
		}
		order = append(order, u)
	}
	for u := 0; u < n; u++ {
		if !seen[u] {
			dfs1(u)
		}
	}
	radj := make([][]int, n)
	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			radj[v] = append(radj[v], u)
		}
	}
	comp := make([]int, n)
	for i := range comp {
		comp[i] = -1
	}
	c := 0
	var dfs2 func(u int)
	dfs2 = func(u int) {
		comp[u] = c
		for _, v := range radj[u] {
			if comp[v] < 0 {
				dfs2(v)
			}
		}
	}
	for i := n - 1; i >= 0; i-- {
		if comp[order[i]] < 0 {
			dfs2(order[i])
			c++
		}
	}
	return comp, c
}

func sccGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 9)
	for {
		g := randDAG(n, edakit.Int(1, 3), 1)
		adj := make([][]int, n)
		for u := range g.adj {
			adj[u] = append(adj[u], g.adj[u]...)
		}
		backs := edakit.Int(1, 2)
		for b := 0; b < backs; b++ {
			u := edakit.Int(1, n-1)
			v := edakit.Int(0, u-1)
			dup := false
			for _, w := range adj[u] {
				dup = dup || w == v
			}
			if !dup {
				adj[u] = append(adj[u], v)
			}
		}
		perm := edakit.Shuffle(seq(n))
		radj := make([][]int, n)
		for u := 0; u < n; u++ {
			for _, v := range adj[u] {
				radj[perm[u]] = append(radj[perm[u]], perm[v])
			}
		}
		comp, c := sccs(n, radj)
		size := make([]int, c)
		for _, x := range comp {
			size[x]++
		}
		big, inLoop := 0, 0
		for _, s := range size {
			big = max(big, s)
			if s >= 2 {
				inLoop += s
			}
		}
		if big < 2 {
			continue
		}
		var parts []string
		for u := 0; u < n; u++ {
			for _, v := range radj[u] {
				parts = append(parts, vName(u)+"→"+vName(v))
			}
		}
		groups := make([][]string, c)
		for v := 0; v < n; v++ {
			groups[comp[v]] = append(groups[comp[v]], vName(v))
		}
		var gt []string
		for _, gr := range groups {
			gt = append(gt, "{"+strings.Join(gr, ", ")+"}")
		}
		expl := "The strongly connected components are " + strings.Join(gt, " ") + "."
		intro := "A gate-level netlist graph has nodes v1 to v" + edakit.IntStr(n) + " and arcs " + strings.Join(parts, ", ") + ". "
		switch kpOf(ctx) {
		case 1:
			return edakit.Numeric(intro+"How many strongly connected components does it have (a node on no cycle is a component by itself)?", edakit.IntStr(c), expl, map[string]string{"sccs": edakit.IntStr(c)})
		case 2:
			return edakit.Numeric(intro+"How many nodes does the largest strongly connected component contain?", edakit.IntStr(big), expl, map[string]string{"largest": edakit.IntStr(big)})
		default:
			return edakit.Numeric(intro+"A loop checker flags every node that lies on a combinational cycle. How many nodes are flagged?", edakit.IntStr(inLoop),
				expl+" A node is on a cycle exactly when its component has two or more nodes (there are no self-loops), so the flagged count is "+edakit.IntStr(inLoop)+".",
				map[string]string{"flagged": edakit.IntStr(inLoop)})
		}
	}
}

func longestPathGen(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	n := edakit.Scaled(ctx, 5, 9)
	if kp == 1 {
		n = edakit.Int(5, 6)
	}
	g := randDAGSources(n, 2, edakit.Scaled(ctx, 2, 6), 9)
	intro := "A timing graph has nodes v1 to v" + edakit.IntStr(n) + " and arcs with delays in ps: " + dagText(g, true) + ". "
	if kp == 3 {
		at := make([]int, n)
		at[0], at[1] = 5*edakit.Int(0, 6), 5*edakit.Int(0, 6)
		// A non-source arrival is the max over predecessors.
		a := make([]int, n)
		a[0], a[1] = at[0], at[1]
		for v := 2; v < n; v++ {
			best := -1
			for _, u := range g.preds(v) {
				best = max(best, a[u]+g.d[[2]int{u, v}])
			}
			a[v] = best
		}
		sink := n - 1
		return edakit.Numeric(intro+edakit.Sprintf("The inputs arrive at v1 at %d ps and at v2 at %d ps. What is the latest arrival time at %s, in ps?", at[0], at[1], vName(sink)),
			edakit.IntStr(a[sink]),
			"Process nodes in topological order and take, at each node, the maximum over incoming arcs of the predecessor's arrival plus the arc delay. Arrivals: "+intsText(a)+".",
			map[string]string{"arrival": edakit.IntStr(a[sink])})
	}
	a := arrivals(g)
	v := n - 1
	if kp == 1 {
		v = edakit.Int(n-3, n-1)
	}
	return edakit.Numeric(intro+"Both v1 and v2 are inputs with arrival 0 ps. What is the latest arrival time (longest-path length) at "+vName(v)+", in ps?",
		edakit.IntStr(a[v]),
		"Relax arcs in topological order, keeping the maximum: arrivals are "+intsText(a)+" for v1 onward.",
		map[string]string{"arrival": edakit.IntStr(a[v])})
}

func reversePassGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 9)
	g := randDAGSources(n, 2, edakit.Scaled(ctx, 2, 6), 9)
	a := arrivals(g)
	isSink := make([]bool, n)
	maxA := 0
	for v := 0; v < n; v++ {
		isSink[v] = len(g.adj[v]) == 0
		if isSink[v] {
			maxA = max(maxA, a[v])
		}
	}
	kp := kpOf(ctx)
	T := maxA + 5*edakit.Int(0, 4)
	if kp == 3 {
		T = maxA
	}
	req := make([]int, n)
	for v := n - 1; v >= 0; v-- {
		if isSink[v] {
			req[v] = T
			continue
		}
		best := 1 << 30
		for _, w := range g.adj[v] {
			best = min(best, req[w]-g.d[[2]int{v, w}])
		}
		req[v] = best
	}
	intro := "A timing graph has nodes v1 to v" + edakit.IntStr(n) + " and arcs with delays in ps: " + dagText(g, true) + ". Inputs v1 and v2 arrive at 0 ps, and every endpoint (node with no outgoing arc) is required by " + edakit.IntStr(T) + " ps. "
	expl := "Forward pass (max): arrivals " + intsText(a) + ". Backward pass in reverse topological order (min of successor required time minus arc delay): required " + intsText(req) + "."
	switch kp {
	case 1:
		var inner []int
		for v := 0; v < n; v++ {
			if !isSink[v] {
				inner = append(inner, v)
			}
		}
		v := edakit.Pick(inner)
		return edakit.Numeric(intro+"What is the required time at "+vName(v)+", in ps?", edakit.IntStr(req[v]), expl, map[string]string{"required": edakit.IntStr(req[v])})
	case 2:
		v := edakit.Int(0, n-2)
		s := req[v] - a[v]
		return edakit.Numeric(intro+"What is the slack (required minus arrival) at "+vName(v)+", in ps?", edakit.IntStr(s), expl+edakit.Sprintf(" Slack at %s = %d - %d = %d.", vName(v), req[v], a[v], s), map[string]string{"slack": edakit.IntStr(s)})
	default:
		z := 0
		for v := 0; v < n; v++ {
			if req[v]-a[v] == 0 {
				z++
			}
		}
		return edakit.Numeric(intro+"The required time equals the latest endpoint arrival. How many nodes have zero slack?", edakit.IntStr(z),
			expl+" Zero-slack nodes lie on a longest path; there are "+edakit.IntStr(z)+".",
			map[string]string{"zero_slack": edakit.IntStr(z)})
	}
}

func pathCountGen(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	if kp == 3 {
		k := edakit.Scaled(ctx, 3, 7)
		prod := 1
		var ws []int
		nodes := 1
		for i := 0; i < k; i++ {
			w := edakit.Int(2, 4)
			ws = append(ws, w)
			prod *= w
			nodes += w + 1
		}
		return edakit.Numeric(edakit.Sprintf("A datapath has %d stages in series. Stage i fans out from one node into $w_i$ parallel gates that reconverge into one node, with widths %s. How many distinct input-to-output paths does the graph have?", k, intsText(ws)),
			edakit.IntStr(prod),
			"Each stage multiplies the number of paths by its width, so the count is the product of "+intsText(ws)+", which is "+edakit.IntStr(prod)+", while the graph has only "+edakit.IntStr(nodes)+" nodes. This is why timers propagate per node instead of enumerating paths.",
			map[string]string{"paths": edakit.IntStr(prod)})
	}
	n := edakit.Scaled(ctx, 6, 9)
	g := randDAGSources(n, 2, edakit.Scaled(ctx, 3, 7), 1)
	from := make([]int, n)
	from[0], from[1] = 1, 1
	for v := 2; v < n; v++ {
		for _, u := range g.preds(v) {
			from[v] += from[u]
		}
	}
	to := make([]int, n)
	for v := n - 1; v >= 0; v-- {
		if len(g.adj[v]) == 0 {
			to[v] = 1
			continue
		}
		for _, w := range g.adj[v] {
			to[v] += to[w]
		}
	}
	intro := "A netlist DAG has nodes v1 to v" + edakit.IntStr(n) + " and arcs " + plainArcs(g) + ". Sources are v1 and v2; sinks are the nodes with no outgoing arc. "
	if kp == 1 {
		total := 0
		for v := 0; v < n; v++ {
			if len(g.adj[v]) == 0 {
				total += from[v]
			}
		}
		return edakit.Numeric(intro+"How many distinct source-to-sink paths are there?", edakit.IntStr(total),
			"Count paths from the sources forward: N(v) = sum of N over predecessors, with N = 1 at each source. The counts are "+intsText(from)+"; summing over sinks gives "+edakit.IntStr(total)+".",
			map[string]string{"paths": edakit.IntStr(total)})
	}
	v := edakit.Int(2, n-2)
	for tries := 0; tries < 20 && from[v]*to[v] < 2; tries++ {
		v = edakit.Int(2, n-2)
	}
	ans := from[v] * to[v]
	return edakit.Numeric(intro+"How many source-to-sink paths pass through "+vName(v)+"?", edakit.IntStr(ans),
		edakit.Sprintf("Paths into %s from the sources: %d. Paths from %s to any sink: %d. Every pair combines, so %d × %d = %d.", vName(v), from[v], vName(v), to[v], from[v], to[v], ans),
		map[string]string{"in": edakit.IntStr(from[v]), "out": edakit.IntStr(to[v])})
}

func incrementalGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 10)
	for {
		g := randDAGSources(n, 2, edakit.Scaled(ctx, 3, 7), 9)
		type arc struct{ u, v int }
		var arcs []arc
		for u := 0; u < n; u++ {
			for _, v := range g.adj[u] {
				arcs = append(arcs, arc{u, v})
			}
		}
		e := edakit.Pick(arcs)
		delta := edakit.Int(2, 8)
		before := arrivals(g)
		g2 := dag{n: n, adj: g.adj, d: map[[2]int]int{}}
		for k, x := range g.d {
			g2.d[k] = x
		}
		g2.d[[2]int{e.u, e.v}] += delta
		after := arrivals(g2)
		cone := map[int]bool{e.v: true}
		for u := e.v; u < n; u++ {
			if cone[u] {
				for _, w := range g.adj[u] {
					cone[w] = true
				}
			}
		}
		changed := 0
		for v := 0; v < n; v++ {
			if after[v] != before[v] {
				changed++
			}
		}
		kp := kpOf(ctx)
		if kp == 2 && changed == len(cone) && edakit.Int(0, 2) > 0 {
			continue
		}
		if kp == 1 && len(cone) < 3 {
			continue
		}
		if kp == 2 && changed == 0 && edakit.Int(0, 1) == 0 {
			continue
		}
		intro := "A timing graph has nodes v1 to v" + edakit.IntStr(n) + " and arcs with delays in ps: " + dagText(g, true) + ". Inputs v1 and v2 arrive at 0 ps. An ECO slows arc " + vName(e.u) + "→" + vName(e.v) + edakit.Sprintf(" by %d ps. ", delta)
		expl := "Arrivals before: " + intsText(before) + ". After: " + intsText(after) + "."
		switch kp {
		case 1:
			return edakit.Numeric(intro+"An incremental timer marks dirty "+vName(e.v)+" and everything in its fanout cone. How many nodes are marked, including "+vName(e.v)+"?",
				edakit.IntStr(len(cone)),
				"Only nodes reachable from "+vName(e.v)+" can see a new arrival. Walking forward from it reaches "+edakit.IntStr(len(cone))+" nodes, including itself.",
				map[string]string{"cone": edakit.IntStr(len(cone))})
		case 2:
			return edakit.Numeric(intro+"How many nodes actually get a different arrival time?", edakit.IntStr(changed),
				expl+" A node in the cone keeps its old arrival when another fanin still dominates it, so propagation can stop early there.",
				map[string]string{"changed": edakit.IntStr(changed)})
		default:
			s := n - 1
			return edakit.Numeric(intro+"What is the new arrival time at "+vName(s)+", in ps?", edakit.IntStr(after[s]), expl, map[string]string{"arrival": edakit.IntStr(after[s])})
		}
	}
}
