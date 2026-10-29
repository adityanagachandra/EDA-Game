package graphs

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// kpOf returns the knowledge point a request targets, choosing one uniformly when the request is
// unguided, so an unguided draw still exercises every kind of question.
func kpOf(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// hyper is a netlist hypergraph: cells 0..n-1 and nets as sorted pin lists.
type hyper struct {
	n    int
	nets [][]int
	w    []int
}

// randHyper draws nNets nets of 2..maxDeg distinct cells over n cells, and makes sure every cell
// sits on at least one net.
func randHyper(n, nNets, maxDeg int) hyper {
	for {
		h := hyper{n: n}
		used := make([]bool, n)
		for j := 0; j < nNets; j++ {
			k := edakit.Int(2, maxDeg)
			if k > n {
				k = n
			}
			perm := edakit.Shuffle(seq(n))
			pins := append([]int(nil), perm[:k]...)
			sort.Ints(pins)
			for _, p := range pins {
				used[p] = true
			}
			h.nets = append(h.nets, pins)
			h.w = append(h.w, 1)
		}
		ok := true
		for _, u := range used {
			if !u {
				ok = false
			}
		}
		if ok {
			return h
		}
	}
}

func seq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func cellName(i int) string { return edakit.Sprintf("c%d", i+1) }

func netText(h hyper, weighted bool) string {
	var parts []string
	for j, net := range h.nets {
		var names []string
		for _, p := range net {
			names = append(names, cellName(p))
		}
		s := edakit.Sprintf("n%d = {%s}", j+1, strings.Join(names, ", "))
		if weighted {
			s += edakit.Sprintf(" (weight %d)", h.w[j])
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, "; ")
}

func sideText(side []int, which int) string {
	var names []string
	for i, s := range side {
		if s == which {
			names = append(names, cellName(i))
		}
	}
	return "{" + strings.Join(names, ", ") + "}"
}

// randSides assigns cells to sides 0 and 1 with both sides non-empty.
func randSides(n int) []int {
	for {
		side := make([]int, n)
		c := 0
		for i := range side {
			side[i] = edakit.Int(0, 1)
			c += side[i]
		}
		if c > 0 && c < n {
			return side
		}
	}
}

// balancedSides splits n cells into two halves (sizes differ by at most one).
func balancedSides(n int) []int {
	perm := edakit.Shuffle(seq(n))
	side := make([]int, n)
	for i, p := range perm {
		if i < n/2 {
			side[p] = 0
		} else {
			side[p] = 1
		}
	}
	return side
}

func isCut(net []int, side []int) bool {
	for _, p := range net[1:] {
		if side[p] != side[net[0]] {
			return true
		}
	}
	return false
}

func cutSize(h hyper, side []int, weighted bool) int {
	t := 0
	for j, net := range h.nets {
		if isCut(net, side) {
			if weighted {
				t += h.w[j]
			} else {
				t++
			}
		}
	}
	return t
}

// fmGain is the Fiduccia-Mattheyses gain of moving cell c to the other side: nets on which c
// is alone on its side (they become uncut) minus nets lying wholly on c's side (they become cut).
func fmGain(h hyper, side []int, c int) (gain, fs, te int) {
	for _, net := range h.nets {
		in := false
		for _, p := range net {
			if p == c {
				in = true
			}
		}
		if !in {
			continue
		}
		same := 0
		for _, p := range net {
			if side[p] == side[c] {
				same++
			}
		}
		if same == 1 {
			fs++
		}
		if same == len(net) {
			te++
		}
	}
	return fs - te, fs, te
}

// wgraph is an undirected weighted graph on n vertices; w[i][j] = 0 means no edge.
type wgraph struct {
	n int
	w [][]int
}

func newWGraph(n int) wgraph {
	g := wgraph{n: n, w: make([][]int, n)}
	for i := range g.w {
		g.w[i] = make([]int, n)
	}
	return g
}

// randWGraph draws a connected weighted graph: a random spanning tree plus extra edges.
func randWGraph(n, extra, wmax int) wgraph {
	g := newWGraph(n)
	perm := edakit.Shuffle(seq(n))
	for i := 1; i < n; i++ {
		a, b := perm[i], perm[edakit.Int(0, i-1)]
		g.w[a][b] = edakit.Int(1, wmax)
		g.w[b][a] = g.w[a][b]
	}
	for t := 0; t < extra*4 && extra > 0; t++ {
		a, b := edakit.Int(0, n-1), edakit.Int(0, n-1)
		if a == b || g.w[a][b] > 0 {
			continue
		}
		g.w[a][b] = edakit.Int(1, wmax)
		g.w[b][a] = g.w[a][b]
		extra--
	}
	return g
}

func vName(i int) string { return edakit.Sprintf("v%d", i+1) }

func edgeText(g wgraph, weighted bool) string {
	var parts []string
	for i := 0; i < g.n; i++ {
		for j := i + 1; j < g.n; j++ {
			if g.w[i][j] > 0 {
				if weighted {
					parts = append(parts, edakit.Sprintf("%s-%s (%d)", vName(i), vName(j), g.w[i][j]))
				} else {
					parts = append(parts, edakit.Sprintf("%s-%s", vName(i), vName(j)))
				}
			}
		}
	}
	return strings.Join(parts, ", ")
}

func graphCut(g wgraph, side []int) int {
	t := 0
	for i := 0; i < g.n; i++ {
		for j := i + 1; j < g.n; j++ {
			if g.w[i][j] > 0 && side[i] != side[j] {
				t += g.w[i][j]
			}
		}
	}
	return t
}

func vSideText(side []int, which int) string {
	var names []string
	for i, s := range side {
		if s == which {
			names = append(names, vName(i))
		}
	}
	return "{" + strings.Join(names, ", ") + "}"
}

// dag is a directed acyclic graph whose edges always go from a lower to a higher index, with an
// integer delay on every edge.
type dag struct {
	n   int
	adj [][]int
	d   map[[2]int]int
}

func (g dag) preds(v int) []int {
	var out []int
	for u := 0; u < g.n; u++ {
		for _, w := range g.adj[u] {
			if w == v {
				out = append(out, u)
			}
		}
	}
	return out
}

// randDAG draws a DAG on n nodes in which every non-first node has at least one predecessor,
// with about extra additional edges.
func randDAG(n, extra, dmax int) dag {
	g := dag{n: n, adj: make([][]int, n), d: map[[2]int]int{}}
	has := map[[2]int]bool{}
	add := func(u, v int) {
		if has[[2]int{u, v}] {
			return
		}
		has[[2]int{u, v}] = true
		g.adj[u] = append(g.adj[u], v)
		g.d[[2]int{u, v}] = edakit.Int(1, dmax)
	}
	for v := 1; v < n; v++ {
		lo := 0
		if v > 3 {
			lo = v - 3
		}
		add(edakit.Int(lo, v-1), v)
	}
	for t := 0; t < extra; t++ {
		u := edakit.Int(0, n-2)
		v := edakit.Int(u+1, n-1)
		add(u, v)
	}
	for u := range g.adj {
		sort.Ints(g.adj[u])
	}
	return g
}

// randDAGSources draws a DAG whose first ns nodes are the only sources: every later node gets a
// predecessor of lower index, and every source drives at least one node.
func randDAGSources(n, ns, extra, dmax int) dag {
	g := dag{n: n, adj: make([][]int, n), d: map[[2]int]int{}}
	has := map[[2]int]bool{}
	add := func(u, v int) {
		if has[[2]int{u, v}] {
			return
		}
		has[[2]int{u, v}] = true
		g.adj[u] = append(g.adj[u], v)
		g.d[[2]int{u, v}] = edakit.Int(1, dmax)
	}
	for v := ns; v < n; v++ {
		add(edakit.Int(0, v-1), v)
	}
	for t := 0; t < extra; t++ {
		v := edakit.Int(ns, n-1)
		u := edakit.Int(0, v-1)
		add(u, v)
	}
	for u := 0; u < ns; u++ {
		if len(g.adj[u]) == 0 {
			add(u, edakit.Int(ns, n-1))
		}
	}
	for u := range g.adj {
		sort.Ints(g.adj[u])
	}
	return g
}

func dagText(g dag, delays bool) string {
	var parts []string
	for u := 0; u < g.n; u++ {
		for _, v := range g.adj[u] {
			if delays {
				parts = append(parts, edakit.Sprintf("%s→%s (%d)", vName(u), vName(v), g.d[[2]int{u, v}]))
			} else {
				parts = append(parts, edakit.Sprintf("%s→%s", vName(u), vName(v)))
			}
		}
	}
	return strings.Join(parts, ", ")
}

// arrivals computes longest-path arrival times from all sources (arrival 0) in index order,
// which is a topological order for these DAGs.
func arrivals(g dag) []int {
	a := make([]int, g.n)
	for u := 0; u < g.n; u++ {
		for _, v := range g.adj[u] {
			if x := a[u] + g.d[[2]int{u, v}]; x > a[v] {
				a[v] = x
			}
		}
	}
	return a
}

type pt struct{ x, y int }

func manh(a, b pt) int { return iabs(a.x-b.x) + iabs(a.y-b.y) }

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// distinctPins draws n pins with pairwise distinct x and distinct y coordinates in [0, span].
func distinctPins(n, span int) []pt {
	xs := edakit.Shuffle(seq(span + 1))[:n]
	ys := edakit.Shuffle(seq(span + 1))[:n]
	out := make([]pt, n)
	for i := range out {
		out[i] = pt{xs[i], ys[i]}
	}
	return out
}

// mstLen is the rectilinear minimum spanning tree length of pts (Prim's algorithm).
func mstLen(pts []pt) int {
	n := len(pts)
	if n <= 1 {
		return 0
	}
	in := make([]bool, n)
	best := make([]int, n)
	for i := range best {
		best[i] = math.MaxInt32
	}
	best[0] = 0
	total := 0
	for k := 0; k < n; k++ {
		u := -1
		for i := 0; i < n; i++ {
			if !in[i] && (u < 0 || best[i] < best[u]) {
				u = i
			}
		}
		in[u] = true
		total += best[u]
		for i := 0; i < n; i++ {
			if !in[i] {
				if d := manh(pts[u], pts[i]); d < best[i] {
					best[i] = d
				}
			}
		}
	}
	return total
}

// primEdges returns the edge lengths of Prim's tree from pin 0 in the order they are added, and
// reports whether every choice was unique (no ties), so a question about the order is well posed.
func primEdges(pts []pt) ([]int, bool) {
	n := len(pts)
	in := make([]bool, n)
	best := make([]int, n)
	for i := range best {
		best[i] = math.MaxInt32
	}
	in[0] = true
	for i := 1; i < n; i++ {
		best[i] = manh(pts[0], pts[i])
	}
	var out []int
	unique := true
	for k := 1; k < n; k++ {
		u := -1
		for i := 0; i < n; i++ {
			if !in[i] && (u < 0 || best[i] < best[u]) {
				u = i
			}
		}
		for i := 0; i < n; i++ {
			if !in[i] && i != u && best[i] == best[u] {
				unique = false
			}
		}
		in[u] = true
		out = append(out, best[u])
		for i := 0; i < n; i++ {
			if !in[i] {
				if d := manh(pts[u], pts[i]); d < best[i] {
					best[i] = d
				}
			}
		}
	}
	return out, unique
}

func hpwl(pts []pt) int {
	minx, maxx, miny, maxy := pts[0].x, pts[0].x, pts[0].y, pts[0].y
	for _, p := range pts {
		minx, maxx = min(minx, p.x), max(maxx, p.x)
		miny, maxy = min(miny, p.y), max(maxy, p.y)
	}
	return maxx - minx + maxy - miny
}

// hananPoints returns the Hanan grid points that are not pins.
func hananPoints(pts []pt) []pt {
	xs, ys := map[int]bool{}, map[int]bool{}
	pin := map[pt]bool{}
	for _, p := range pts {
		xs[p.x], ys[p.y] = true, true
		pin[p] = true
	}
	var out []pt
	for x := range xs {
		for y := range ys {
			if !pin[pt{x, y}] {
				out = append(out, pt{x, y})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].x != out[j].x {
			return out[i].x < out[j].x
		}
		return out[i].y < out[j].y
	})
	return out
}

// rsmtLen is the exact rectilinear Steiner minimum tree length for up to five pins: by Hanan's
// theorem some optimal tree uses at most n-2 Steiner points, all on the Hanan grid, so the
// minimum over such point sets of the rectilinear MST is the RSMT.
func rsmtLen(pts []pt) int {
	h := hananPoints(pts)
	best := mstLen(pts)
	maxS := len(pts) - 2
	var rec func(start int, chosen []pt)
	rec = func(start int, chosen []pt) {
		if len(chosen) > 0 {
			if l := mstLen(append(append([]pt(nil), pts...), chosen...)); l < best {
				best = l
			}
		}
		if len(chosen) == maxS {
			return
		}
		for i := start; i < len(h); i++ {
			rec(i+1, append(chosen, h[i]))
		}
	}
	rec(0, nil)
	return best
}

func ptsText(pts []pt, prefix string) string {
	var parts []string
	for i, p := range pts {
		parts = append(parts, edakit.Sprintf("%s%d (%d, %d)", prefix, i+1, p.x, p.y))
	}
	return strings.Join(parts, ", ")
}

func intsText(xs []int) string {
	var parts []string
	for _, x := range xs {
		parts = append(parts, edakit.IntStr(x))
	}
	return strings.Join(parts, ", ")
}
