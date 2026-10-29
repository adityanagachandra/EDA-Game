package route

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// channel is a two-sided channel: top[c] and bottom[c] hold a net number or 0, c from 1 to n.
type channel struct {
	n, k        int
	top, bottom []int
}

// fallbackChannel is a known-good channel: acyclic VCG 1 -> 2 -> 3 -> 4, or with cyclic set,
// a 2-cycle between nets 1 and 2.
func fallbackChannel(cyclic bool) channel {
	if cyclic {
		return channel{n: 8, k: 4, top: []int{0, 1, 2, 3, 4, 3, 0, 0, 0}, bottom: []int{0, 2, 1, 4, 0, 0, 0, 0, 0}}
	}
	return channel{n: 8, k: 4, top: []int{0, 1, 2, 0, 1, 3, 0, 4, 0}, bottom: []int{0, 2, 3, 0, 0, 4, 0, 0, 0}}
}

func randChannel(ctx generator.GeneratorContext) channel {
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		n := edakit.Scaled(ctx, 8, 12)
		k := edakit.Scaled(ctx, 4, 6)
		ch := channel{n: n, k: k, top: make([]int, n+1), bottom: make([]int, n+1)}
		type slot struct{ c, side int }
		var slots []slot
		for c := 1; c <= n; c++ {
			slots = append(slots, slot{c, 0}, slot{c, 1})
		}
		slots = edakit.Shuffle(slots)
		i := 0
		ok := true
		for net := 1; net <= k && ok; net++ {
			want := edakit.Int(2, 3)
			got := 0
			for got < want && i < len(slots) {
				s := slots[i]
				i++
				other := ch.bottom[s.c]
				if s.side == 1 {
					other = ch.top[s.c]
				}
				if other == net {
					continue
				}
				if s.side == 0 {
					ch.top[s.c] = net
				} else {
					ch.bottom[s.c] = net
				}
				got++
			}
			if got < 2 {
				ok = false
			}
		}
		if ok {
			return ch
		}
	}
	return fallbackChannel(false)
}

func (ch channel) text() string {
	return edakit.Sprintf("A channel has columns 1 to %d. Reading columns left to right, the top-side pins are %s and the bottom-side pins are %s, where a number names a net and 0 means no pin. ",
		ch.n, intsText(ch.top[1:]), intsText(ch.bottom[1:]))
}

func (ch channel) pinCols(net int) []int {
	var cs []int
	for c := 1; c <= ch.n; c++ {
		if ch.top[c] == net || ch.bottom[c] == net {
			cs = append(cs, c)
		}
	}
	return cs
}

func (ch channel) intervals() []iv {
	out := make([]iv, ch.k+1)
	for net := 1; net <= ch.k; net++ {
		cs := ch.pinCols(net)
		out[net] = iv{cs[0], cs[len(cs)-1]}
	}
	return out
}

func (ch channel) localDensity(c int) int {
	ivs := ch.intervals()
	d := 0
	for net := 1; net <= ch.k; net++ {
		if ivs[net].l <= c && c <= ivs[net].r {
			d++
		}
	}
	return d
}

func (ch channel) density() int {
	best := 0
	for c := 1; c <= ch.n; c++ {
		best = maxI(best, ch.localDensity(c))
	}
	return best
}

// vcg returns adjacency: adj[a][b] means net a must be above net b.
func (ch channel) vcg() map[int]map[int]bool {
	adj := map[int]map[int]bool{}
	for net := 1; net <= ch.k; net++ {
		adj[net] = map[int]bool{}
	}
	for c := 1; c <= ch.n; c++ {
		a, b := ch.top[c], ch.bottom[c]
		if a > 0 && b > 0 && a != b {
			adj[a][b] = true
		}
	}
	return adj
}

func countEdges(adj map[int]map[int]bool) int {
	n := 0
	for _, m := range adj {
		n += len(m)
	}
	return n
}

// longestChain returns the most vertices on a directed path, or -1 if adj has a cycle.
func longestChain(adj map[int]map[int]bool) int {
	memo := map[int]int{}
	state := map[int]int{}
	cyc := false
	var f func(u int) int
	f = func(u int) int {
		if state[u] == 2 {
			return memo[u]
		}
		if state[u] == 1 {
			cyc = true
			return 0
		}
		state[u] = 1
		best := 1
		for v := range adj[u] {
			best = maxI(best, 1+f(v))
		}
		state[u] = 2
		memo[u] = best
		return best
	}
	best := 0
	for u := range adj {
		best = maxI(best, f(u))
	}
	if cyc {
		return -1
	}
	return best
}

func reach(adj map[int]map[int]bool, s int) map[int]bool {
	seen := map[int]bool{}
	stack := []int{s}
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for v := range adj[u] {
			if !seen[v] {
				seen[v] = true
				stack = append(stack, v)
			}
		}
	}
	return seen
}

// onCycle lists nets that lie on some directed cycle.
func onCycle(adj map[int]map[int]bool) []int {
	var out []int
	for u := range adj {
		if reach(adj, u)[u] {
			out = append(out, u)
		}
	}
	sort.Ints(out)
	return out
}

func without(adj map[int]map[int]bool, drop map[int]bool) map[int]map[int]bool {
	out := map[int]map[int]bool{}
	for u, m := range adj {
		if drop[u] {
			continue
		}
		out[u] = map[int]bool{}
		for v := range m {
			if !drop[v] {
				out[u][v] = true
			}
		}
	}
	return out
}

func minFVS(adj map[int]map[int]bool, k int) int {
	best := k
	for mask := 0; mask < 1<<k; mask++ {
		drop := map[int]bool{}
		c := 0
		for i := 0; i < k; i++ {
			if mask&(1<<i) != 0 {
				drop[i+1] = true
				c++
			}
		}
		if c < best && longestChain(without(adj, drop)) >= 0 {
			best = c
		}
	}
	return best
}

func shortestCycle(adj map[int]map[int]bool) int {
	best := 1 << 30
	for s := range adj {
		dist := map[int]int{s: 0}
		q := []int{s}
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			for v := range adj[u] {
				if v == s {
					best = minI(best, dist[u]+1)
				}
				if _, ok := dist[v]; !ok {
					dist[v] = dist[u] + 1
					q = append(q, v)
				}
			}
		}
	}
	return best
}

// leftEdge runs the constrained left-edge algorithm; with useVCG false it ignores vertical
// constraints. It returns the tracks (top first) as lists of nets, or nil on a cyclic VCG.
func leftEdge(ivs []iv, adj map[int]map[int]bool, k int, useVCG bool) [][]int {
	placed := map[int]bool{}
	var tracks [][]int
	order := make([]int, 0, k)
	for net := 1; net <= k; net++ {
		order = append(order, net)
	}
	sort.SliceStable(order, func(a, b int) bool {
		if ivs[order[a]].l != ivs[order[b]].l {
			return ivs[order[a]].l < ivs[order[b]].l
		}
		return order[a] < order[b]
	})
	for len(placed) < k {
		var tr []int
		last := -1
		for _, net := range order {
			if placed[net] {
				continue
			}
			if useVCG {
				ready := true
				for u := range adj {
					if adj[u][net] && !placed[u] {
						ready = false
					}
				}
				if !ready {
					continue
				}
			}
			if ivs[net].l > last {
				tr = append(tr, net)
				last = ivs[net].r
			}
		}
		if len(tr) == 0 {
			return nil
		}
		for _, net := range tr {
			placed[net] = true
		}
		tracks = append(tracks, tr)
	}
	return tracks
}

func ivList(ivs []iv, k int) string {
	var parts []string
	for net := 1; net <= k; net++ {
		parts = append(parts, edakit.Sprintf("net %d [%d, %d]", net, ivs[net].l, ivs[net].r))
	}
	return strings.Join(parts, ", ")
}

// channelDensity: local density at a column (KP 1), channel density (KP 2), columns at maximum
// density (KP 3).
func channelDensity(ctx generator.GeneratorContext) generator.Problem {
	ch := randChannel(ctx)
	ivs := ch.intervals()
	intro := ch.text() + "Each net's horizontal trunk spans from its leftmost to its rightmost pin column, inclusive. "
	switch kp(ctx) {
	case 1:
		c := edakit.Int(1, ch.n)
		d := ch.localDensity(c)
		return edakit.Numeric(intro+edakit.Sprintf("What is the local density at column %d, the number of nets whose span includes that column?", c), edakit.IntStr(d),
			edakit.Sprintf("The spans are %s. %d of them include column %d.", ivList(ivs, ch.k), d, c),
			map[string]string{"local": edakit.IntStr(d)})
	case 2:
		d := ch.density()
		return edakit.Numeric(intro+"What is the channel density, the lower bound on tracks for any routing that gives each net one trunk?", edakit.IntStr(d),
			edakit.Sprintf("The spans are %s. Sweeping the columns, the most spans covering one column is %d. Those nets need distinct tracks at that column.", ivList(ivs, ch.k), d),
			map[string]string{"density": edakit.IntStr(d)})
	default:
		d := ch.density()
		n := 0
		for c := 1; c <= ch.n; c++ {
			if ch.localDensity(c) == d {
				n++
			}
		}
		return edakit.Numeric(intro+"The channel density is the largest local density over the columns. How many columns attain it?", edakit.IntStr(n),
			edakit.Sprintf("The spans are %s, the density is %d, and %d columns reach it. Every one of those columns is a bottleneck that a router cannot squeeze.", ivList(ivs, ch.k), d, n),
			map[string]string{"density": edakit.IntStr(d), "columns": edakit.IntStr(n)})
	}
}

// hcgVCG: VCG edges (KP 1), the longest VCG chain (KP 2), HCG edges (KP 3).
func hcgVCG(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	ch := fallbackChannel(false)
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		c := randChannel(ctx)
		if k != 2 || longestChain(c.vcg()) > 0 {
			ch = c
			break
		}
	}
	adj := ch.vcg()
	intro := ch.text()
	switch k {
	case 1:
		e := countEdges(adj)
		return edakit.Numeric(intro+"In the vertical constraint graph, a column with net a on top and a different net b at the bottom forces a's trunk above b's, an edge from a to b. How many distinct edges does the VCG have?", edakit.IntStr(e),
			edakit.Sprintf("Read each column with two different nets; repeated pairs give the same edge. There are %d distinct edges.", e),
			map[string]string{"vcg_edges": edakit.IntStr(e)})
	case 2:
		L := longestChain(adj)
		return edakit.Numeric(intro+"Vertical constraints (net on top above net at the bottom in each column) form an acyclic VCG here. Without doglegs, how many tracks does the longest chain of vertical constraints force, that is, how many nets lie on the longest directed path?", edakit.IntStr(L),
			edakit.Sprintf("Each edge forces a strictly lower track, so a chain of %d nets needs %d tracks. The channel density is %d, and the tracks needed are at least the larger of the two bounds.", L, L, ch.density()),
			map[string]string{"chain": edakit.IntStr(L), "density": edakit.IntStr(ch.density())})
	default:
		ivs := ch.intervals()
		e := 0
		for a := 1; a <= ch.k; a++ {
			for b := a + 1; b <= ch.k; b++ {
				if overlapLen(ivs[a], ivs[b]) > 0 {
					e++
				}
			}
		}
		return edakit.Numeric(intro+"Each net's trunk spans its leftmost to rightmost pin column, inclusive. In the horizontal constraint graph, two nets are joined if their spans share a column, so they cannot share a track. How many edges does the HCG have?", edakit.IntStr(e),
			edakit.Sprintf("The spans are %s. %d pairs overlap. The HCG is an interval graph; its largest clique is the channel density, %d.", ivList(ivs, ch.k), e, ch.density()),
			map[string]string{"hcg_edges": edakit.IntStr(e)})
	}
}

// leftEdgeGen: unconstrained left-edge tracks (KP 1), constrained left-edge tracks (KP 2), nets
// on the top track (KP 3).
func leftEdgeGen(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1:
		n := edakit.Scaled(ctx, 5, 8)
		ivs := make([]iv, n+1)
		for i := 1; i <= n; i++ {
			l := edakit.Int(1, 14)
			ivs[i] = iv{l, minI(16, l+edakit.Int(1, 7))}
		}
		tr := leftEdge(ivs, nil, n, false)
		var flat []iv
		for i := 1; i <= n; i++ {
			flat = append(flat, ivs[i])
		}
		q := edakit.Sprintf("A channel has no vertical constraints. Its nets' trunks span the inclusive column ranges %s. The left-edge algorithm sorts the nets by left end and fills one track at a time, adding each net whose left end lies to the right of the last net placed on the track. How many tracks does it use?", ivList(ivs, n))
		return edakit.Numeric(q, edakit.IntStr(len(tr)),
			edakit.Sprintf("The tracks fill as %s. Without vertical constraints the left-edge algorithm is optimal: it uses exactly the density, %d.", tracksText(tr), maxOverlap(flat)),
			map[string]string{"tracks": edakit.IntStr(len(tr))})
	default:
		ch := fallbackChannel(false)
		tr := leftEdge(ch.intervals(), ch.vcg(), ch.k, true)
		for try := 0; try < generator.MaxDrawAttempts; try++ {
			c := randChannel(ctx)
			adj := c.vcg()
			if longestChain(adj) < 0 || countEdges(adj) < 2 {
				continue
			}
			if t := leftEdge(c.intervals(), adj, c.k, true); t != nil {
				ch, tr = c, t
				break
			}
		}
		intro := ch.text() + "Each net's trunk spans its leftmost to rightmost pin column, inclusive, and the vertical constraint graph is acyclic. The constrained left-edge algorithm fills tracks from the top. For each track it scans the unplaced nets in order of left end (ties by net number) and adds a net if every net that must lie above it is already on an earlier track and its left end lies to the right of the last net placed on this track. "
		if k == 2 {
			return edakit.Numeric(intro+"How many tracks does it use?", edakit.IntStr(len(tr)),
				edakit.Sprintf("Spans: %s. Tracks from the top: %s, so %d tracks against a density of %d and a longest VCG chain of %d.", ivList(ch.intervals(), ch.k), tracksText(tr), len(tr), ch.density(), longestChain(ch.vcg())),
				map[string]string{"tracks": edakit.IntStr(len(tr))})
		}
		return edakit.Numeric(intro+"How many nets does it place on the top track?", edakit.IntStr(len(tr[0])),
			edakit.Sprintf("Only nets with no VCG predecessor may go on the top track. Scanning those by left end gives %s, %d nets. All tracks: %s.", tracksText(tr[:1]), len(tr[0]), tracksText(tr)),
			map[string]string{"top": edakit.IntStr(len(tr[0]))})
	}
}

func tracksText(tr [][]int) string {
	var parts []string
	for i, t := range tr {
		parts = append(parts, edakit.Sprintf("track %d {%s}", i+1, intsText(t)))
	}
	return strings.Join(parts, ", ")
}

// cyclicVCG: nets on a cycle (KP 1), the minimum nets to split (KP 2), the shortest cycle
// (KP 3).
func cyclicVCG(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	ch := fallbackChannel(true)
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		c := randChannel(ctx)
		cyc := longestChain(c.vcg()) < 0
		if cyc || (k == 1 && edakit.Int(0, 3) == 0) {
			ch = c
			break
		}
	}
	adj := ch.vcg()
	intro := ch.text() + "In each column with two different nets, the top net's trunk must lie above the bottom net's: an edge of the vertical constraint graph. "
	switch k {
	case 1:
		on := onCycle(adj)
		exp := "A net is on a cycle exactly when it can reach itself. No net can, so the VCG is acyclic and a dogleg-free routing exists."
		if len(on) > 0 {
			exp = edakit.Sprintf("A net is on a cycle exactly when it can reach itself. The nets on cycles are {%s}, %d in all. Any such cycle makes a dogleg-free routing impossible.", intsText(on), len(on))
		}
		return edakit.Numeric(intro+"How many nets lie on at least one directed cycle of the VCG? Answer 0 if it is acyclic.", edakit.IntStr(len(on)),
			exp,
			map[string]string{"on_cycle": edakit.IntStr(len(on))})
	case 2:
		f := minFVS(adj, ch.k)
		return edakit.Numeric(intro+"The VCG has a cycle. A dogleg router must split at least one net on every cycle, so the nets it splits form a set that meets every cycle. What is the smallest such set's size, the fewest nets whose removal leaves the VCG acyclic?", edakit.IntStr(f),
			edakit.Sprintf("Try removing sets of nets in increasing size until no cycle remains: the minimum is %d. Splitting a net only helps if it has a pin column between the conflicting columns, so this is a lower bound, not a guarantee.", f),
			map[string]string{"fvs": edakit.IntStr(f)})
	default:
		s := shortestCycle(adj)
		return edakit.Numeric(intro+"The VCG has a cycle. How many nets are on its shortest directed cycle?", edakit.IntStr(s),
			edakit.Sprintf("Breadth-first search from each net for a path back to itself finds a shortest cycle of %d nets. A 2-cycle arises when two nets swap top and bottom positions in two columns.", s),
			map[string]string{"cycle": edakit.IntStr(s)})
	}
}

// subnets splits every net at its pin columns into consecutive spans.
func (ch channel) subnets() ([]iv, []int) {
	var spans []iv
	var owner []int
	for net := 1; net <= ch.k; net++ {
		cs := ch.pinCols(net)
		for i := 0; i+1 < len(cs); i++ {
			spans = append(spans, iv{cs[i], cs[i+1]})
			owner = append(owner, net)
		}
	}
	return spans, owner
}

func (ch channel) subnetVCG() map[int]map[int]bool {
	spans, owner := ch.subnets()
	adj := map[int]map[int]bool{}
	for i := range spans {
		adj[i] = map[int]bool{}
	}
	for c := 1; c <= ch.n; c++ {
		a, b := ch.top[c], ch.bottom[c]
		if a == 0 || b == 0 || a == b {
			continue
		}
		for i := range spans {
			if owner[i] != a || spans[i].l > c || spans[i].r < c {
				continue
			}
			for j := range spans {
				if owner[j] == b && spans[j].l <= c && c <= spans[j].r {
					adj[i][j] = true
				}
			}
		}
	}
	return adj
}

// dogleg: subnet count (KP 1), longest chain after splitting (KP 2), how much splitting shortens
// the chain bound (KP 3).
func dogleg(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	ch := fallbackChannel(false)
	for try := 0; try < generator.MaxDrawAttempts; try++ {
		c := randChannel(ctx)
		if k == 1 {
			ch = c
			break
		}
		b := longestChain(c.vcg())
		if b >= 3 && (longestChain(c.subnetVCG()) < b || edakit.Int(0, 3) == 0) {
			ch = c
			break
		}
	}
	spans, _ := ch.subnets()
	intro := ch.text() + "A dogleg router splits every net at each of its pin columns, so a net with pins in m distinct columns becomes m - 1 subnets, each spanning two consecutive pin columns; subnets of different nets that cover the same column inherit that column's vertical constraint. "
	switch k {
	case 1:
		return edakit.Numeric(intro+"How many subnets result?", edakit.IntStr(len(spans)),
			edakit.Sprintf("Count each net's distinct pin columns and subtract one: %d subnets in all.", len(spans)),
			map[string]string{"subnets": edakit.IntStr(len(spans))})
	case 2:
		L := longestChain(ch.subnetVCG())
		return edakit.Numeric(intro+"After splitting, how many subnets lie on the longest directed path of the subnet VCG?", edakit.IntStr(L),
			edakit.Sprintf("Build the VCG over subnets, adding an edge only between subnets that both cover the constraining column. Its longest path has %d subnets, against %d nets before splitting.", L, longestChain(ch.vcg())),
			map[string]string{"after": edakit.IntStr(L), "before": edakit.IntStr(longestChain(ch.vcg()))})
	default:
		b, a := longestChain(ch.vcg()), longestChain(ch.subnetVCG())
		return edakit.Numeric(intro+"By how many tracks does splitting lower the longest-chain lower bound (the number of vertices on the longest VCG path, nets before and subnets after)?", edakit.IntStr(b-a),
			edakit.Sprintf("Before splitting the longest chain has %d nets; afterwards %d subnets. The bound drops by %d. Splitting removes constraints between parts of a net that never share a column with the other net.", b, a, b-a),
			map[string]string{"before": edakit.IntStr(b), "after": edakit.IntStr(a)})
	}
}
