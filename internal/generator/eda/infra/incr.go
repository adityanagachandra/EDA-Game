package infra

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerIncr) }

func registerIncr(reg *generator.Registry) {
	reg.Register("infra.incr.dependency_graph", edakit.Func(depGraph))
	reg.Register("infra.incr.dirty_flags", edakit.Func(dirtyFlags))
	reg.Register("infra.incr.incremental_sta", &edakit.Bank{Items: incrementalStaItems})
	reg.Register("infra.incr.incremental_legalization", edakit.Func(incrementalLegalization))
	reg.Register("infra.incr.result_caching", edakit.Func(resultCaching))
}

// incrDAG is a dependency graph of derived results. Edges run from an input to a result
// computed from it, and every edge goes from a lower index to a higher one, so index order is
// a topological order.
type incrDAG struct {
	n     int
	names []string
	succ  [][]int
	pred  [][]int
	edges [][2]int
}

var incrNamePool = []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R"}

// newIncrDAG draws n nodes; every non-first node gets one or two inputs from earlier nodes,
// and a few extra edges are added.
func newIncrDAG(n int) *incrDAG {
	g := &incrDAG{n: n, succ: make([][]int, n), pred: make([][]int, n)}
	g.names = edakit.Shuffle(incrNamePool)[:n]
	has := map[[2]int]bool{}
	add := func(u, v int) {
		if u == v || has[[2]int{u, v}] {
			return
		}
		has[[2]int{u, v}] = true
		g.succ[u] = append(g.succ[u], v)
		g.pred[v] = append(g.pred[v], u)
		g.edges = append(g.edges, [2]int{u, v})
	}
	for v := 2; v < n; v++ {
		add(edakit.Int(0, v-1), v)
		if edakit.Int(0, 2) == 0 {
			add(edakit.Int(0, v-1), v)
		}
	}
	add(0, 1)
	return g
}

func (g *incrDAG) describe() string {
	es := edakit.Shuffle(g.edges)
	var parts []string
	for _, e := range es {
		parts = append(parts, g.names[e[0]]+"→"+g.names[e[1]])
	}
	return strings.Join(parts, ", ")
}

// reach returns the nodes reachable from the start set through successor edges, excluding the
// start nodes themselves unless they are reachable from another start node. A blocked node is
// visited but not expanded.
func (g *incrDAG) reach(starts []int, blocked int) map[int]bool {
	seen := map[int]bool{}
	stack := []int{}
	for _, s := range starts {
		stack = append(stack, g.succ[s]...)
	}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[v] {
			continue
		}
		seen[v] = true
		if v == blocked {
			continue
		}
		stack = append(stack, g.succ[v]...)
	}
	return seen
}

func (g *incrDAG) ancestorsOf(q int) map[int]bool {
	seen := map[int]bool{q: true}
	stack := []int{q}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, u := range g.pred[v] {
			if !seen[u] {
				seen[u] = true
				stack = append(stack, u)
			}
		}
	}
	return seen
}

func (g *incrDAG) list(m map[int]bool) string {
	var xs []string
	for v := 0; v < g.n; v++ {
		if m[v] {
			xs = append(xs, g.names[v])
		}
	}
	sort.Strings(xs)
	if len(xs) == 0 {
		return "none"
	}
	return strings.Join(xs, ", ")
}

// depGraph asks for invalidation counts, query-driven recomputation counts and level counts.
func depGraph(ctx generator.GeneratorContext) generator.Problem { return depGraphAttempt(ctx, 0) }

func depGraphAttempt(ctx generator.GeneratorContext, attempt int) generator.Problem {
	var g *incrDAG
	var edits []int
	var dirty map[int]bool
	for tries := 0; tries < generator.MaxDrawAttempts; tries++ {
		g = newIncrDAG(edakit.Scaled(ctx, 7, 12))
		edits = []int{edakit.Int(0, g.n/2)}
		if ctx.Difficulty > 0.6 && edakit.Int(0, 1) == 0 {
			e2 := edakit.Int(0, g.n-2)
			if e2 != edits[0] {
				edits = append(edits, e2)
			}
		}
		dirty = g.reach(edits, -1)
		for _, e := range edits {
			delete(dirty, e)
		}
		if len(dirty) >= 2 && len(dirty) < g.n-len(edits) || tries > 40 {
			break
		}
	}
	var en []string
	for _, e := range edits {
		en = append(en, g.names[e])
	}
	q := edakit.Sprintf("A tool caches derived results in a dependency graph whose edges run from an input to a result computed from it: %s. An edit directly modifies %s. ", g.describe(), names(en))
	facts := map[string]string{"invalidated": edakit.IntStr(len(dirty))}
	switch variant(ctx) {
	case 1:
		return edakit.Numeric(q+"How many results, not counting the edited node itself, must be marked invalid?", edakit.IntStr(len(dirty)),
			edakit.Sprintf("Invalidation follows edges forward: every result reachable from %s may now be stale. Those are %s, %d results.", names(en), g.list(dirty), len(dirty)), facts)
	case 2:
		// Query a dirty node with the most dirty ancestors.
		best, bestN := -1, -1
		for v := range dirty {
			anc := g.ancestorsOf(v)
			c := 0
			for u := range anc {
				if dirty[u] {
					c++
				}
			}
			if c > bestN || (c == bestN && v < best) {
				best, bestN = v, c
			}
		}
		need := map[int]bool{}
		for u := range g.ancestorsOf(best) {
			if dirty[u] {
				need[u] = true
			}
		}
		if len(need) < 2 && attempt < generator.MaxDrawAttempts && edakit.Int(0, 3) > 0 {
			return depGraphAttempt(ctx, attempt+1)
		}
		facts["recomputed"] = edakit.IntStr(len(need))
		return edakit.Numeric(q+edakit.Sprintf("The invalid results are marked dirty and nothing is recomputed yet. Then one query asks for %s. How many results must be recomputed to answer it?", g.names[best]), edakit.IntStr(len(need)),
			edakit.Sprintf("Dirty results are %s. The query needs %s and every dirty result it depends on: %s, %d recomputations. The other dirty results stay dirty until someone asks.", g.list(dirty), g.names[best], g.list(need), len(need)), facts)
	}
	level := map[int]int{}
	maxL := 0
	for v := 0; v < g.n; v++ {
		if !dirty[v] {
			continue
		}
		l := 1
		for _, u := range g.pred[v] {
			if dirty[u] && level[u]+1 > l {
				l = level[u] + 1
			}
		}
		level[v] = l
		if l > maxL {
			maxL = l
		}
	}
	facts["levels"] = edakit.IntStr(maxL)
	var lv []string
	for v := 0; v < g.n; v++ {
		if dirty[v] {
			lv = append(lv, edakit.Sprintf("%s at level %d", g.names[v], level[v]))
		}
	}
	return edakit.Numeric(q+"All invalidated results are recomputed in rounds; in each round, every result whose dirty inputs are already recomputed runs in parallel. How many rounds are needed?", edakit.IntStr(maxL),
		edakit.Sprintf("A dirty result's level is one more than the highest level of its dirty inputs (clean inputs are ready at once): %s. The number of rounds is the deepest level, %d.", strings.Join(lv, ", "), maxL), facts)
}

// dirtyFlags compares eager and lazy recomputation and shows early cutoff.
func dirtyFlags(ctx generator.GeneratorContext) generator.Problem {
	v := variant(ctx)
	if v == 3 {
		var g *incrDAG
		var e, x int
		var full, cut map[int]bool
		for tries := 0; tries < generator.MaxDrawAttempts; tries++ {
			g = newIncrDAG(edakit.Scaled(ctx, 7, 11))
			e = edakit.Int(0, 2)
			full = g.reach([]int{e}, -1)
			var cands []int
			for u := range full {
				if len(g.succ[u]) > 0 {
					cands = append(cands, u)
				}
			}
			sort.Ints(cands)
			if len(cands) == 0 {
				cands = []int{e}
			}
			x = edakit.Pick(cands)
			cut = g.reach([]int{e}, x)
			if len(cut) < len(full) || tries > 40 {
				break
			}
		}
		q := edakit.Sprintf("A dependency graph has edges from input to derived result: %s. An edit changes %s. A result is re-evaluated only if at least one of its inputs changed value. When %s is re-evaluated, its new value equals its old value; every other re-evaluated result does change. How many results are re-evaluated?", g.describe(), g.names[e], g.names[x])
		return edakit.Numeric(q, edakit.IntStr(len(cut)),
			edakit.Sprintf("Without the cutoff, all %d results downstream of %s would be re-evaluated (%s). Because %s comes out unchanged, change propagation stops there, and only results reachable from %s along some path that avoids passing through %s are re-evaluated, plus %s itself: %s, which is %d results.", len(full), g.names[e], g.list(full), g.names[x], g.names[e], g.names[x], g.names[x], g.list(cut), len(cut)),
			map[string]string{"without_cutoff": edakit.IntStr(len(full)), "with_cutoff": edakit.IntStr(len(cut))})
	}
	g := newIncrDAG(edakit.Scaled(ctx, 7, 11))
	k := edakit.Scaled(ctx, 2, 4)
	var cands []int
	for e := 0; e < g.n-1; e++ {
		if len(g.reach([]int{e}, -1)) > 0 {
			cands = append(cands, e)
		}
	}
	cands = edakit.Shuffle(cands)
	if k > len(cands) {
		k = len(cands)
	}
	edits := cands[:k]
	var en []string
	for _, e := range edits {
		en = append(en, g.names[e])
	}
	eager := 0
	var per []string
	for _, e := range edits {
		r := g.reach([]int{e}, -1)
		eager += len(r)
		per = append(per, edakit.Sprintf("%s reaches %d", g.names[e], len(r)))
	}
	union := g.reach(edits, -1)
	q := edakit.Sprintf("A dependency graph has edges from input to derived result: %s. A script applies %d edits in sequence, changing %s, and then queries every result once. Assume every downstream result changes. ", g.describe(), k, strings.Join(en, ", then "))
	facts := map[string]string{"eager": edakit.IntStr(eager), "lazy": edakit.IntStr(len(union))}
	if v == 1 {
		return edakit.Numeric(q+"With eager recomputation, each edit immediately recomputes every result downstream of the edited node. How many result evaluations happen in total?", edakit.IntStr(eager),
			edakit.Sprintf("Eager evaluation pays for each edit separately: %s, for %d evaluations. Results shared by several edits are recomputed each time.", strings.Join(per, ", "), eager), facts)
	}
	return edakit.Numeric(q+"With dirty flags, an edit only marks downstream results dirty, and the final query recomputes each dirty result once. How many result evaluations happen?", edakit.IntStr(len(union)),
		edakit.Sprintf("The dirty set is the union of the downstream sets: %s, %d results, each recomputed once at query time. Eager evaluation would have cost %d evaluations.", g.list(union), len(union), eager), facts)
}

// incrementalLegalization widens or inserts a cell in a packed row and shifts neighbors right.
func incrementalLegalization(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	xs := make([]int, n)
	ws := make([]int, n)
	x := edakit.Int(0, 3)
	for i := 0; i < n; i++ {
		ws[i] = edakit.Int(2, 6)
		xs[i] = x
		gap := edakit.Int(0, 3)
		if edakit.Int(0, 2) == 0 {
			gap = 0
		}
		x += ws[i] + gap
	}
	k := edakit.Int(0, n/2)
	grow := edakit.Int(3, 9)
	newX := append([]int(nil), xs...)
	end := xs[k] + ws[k] + grow
	moved := 0
	total := 0
	for i := k + 1; i < n; i++ {
		if newX[i] < end {
			newX[i] = end
		}
		if newX[i] != xs[i] {
			moved++
			total += newX[i] - xs[i]
		}
		end = newX[i] + ws[i]
	}
	var parts []string
	for i := 0; i < n; i++ {
		parts = append(parts, edakit.Sprintf("c%d at x = %d (width %d)", i+1, xs[i], ws[i]))
	}
	var after []string
	for i := k + 1; i < n; i++ {
		after = append(after, edakit.Sprintf("c%d → %d", i+1, newX[i]))
	}
	q := edakit.Sprintf("A placement row holds, left to right, %s, all in site units. A resize ECO widens c%d by %d sites, keeping its left edge fixed. The incremental legalizer only pushes cells to the right, each by the least amount that removes its overlap, and leaves every other cell alone. ", strings.Join(parts, ", "), k+1, grow)
	exp := edakit.Sprintf("c%d now ends at %d. Walking right, each cell starts at the larger of its old x and the end of its left neighbour: %s. Shifting stops at the first cell whose gap absorbs the push.", k+1, xs[k]+ws[k]+grow, strings.Join(after, ", "))
	facts := map[string]string{"moved": edakit.IntStr(moved), "displacement": edakit.IntStr(total)}
	switch variant(ctx) {
	case 1:
		return edakit.Numeric(q+"How many cells move?", edakit.IntStr(moved), exp+edakit.Sprintf(" %d cells move.", moved), facts)
	case 2:
		j := k + 1
		for i := k + 1; i < n; i++ {
			if newX[i] != xs[i] {
				j = i
			}
		}
		if j >= n {
			j = n - 1
		}
		return edakit.Numeric(q+edakit.Sprintf("What is the new x of c%d, in sites?", j+1), edakit.IntStr(newX[j]), exp, facts)
	}
	return edakit.Numeric(q+"What is the total displacement of all moved cells, in sites?", edakit.IntStr(total), exp+edakit.Sprintf(" The displacements sum to %d sites.", total), facts)
}

// resultCaching asks for expected run time, break-even hit rate and hit counting by key.
func resultCaching(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		lookup := float64(edakit.Int(2, 30))
		compute := float64(10 * edakit.Int(6, 60))
		h := float64(edakit.Int(40, 95))
		exp := lookup + (1-h/100)*compute
		q := edakit.Sprintf("A flow caches the result of an extraction step keyed by a hash of its inputs. A lookup costs %s s (hashing the inputs and probing the store), and on a miss the step runs for %s s after the lookup. The hit rate is %s%%. What is the expected time per invocation, in seconds? Round to 2 decimal places.", edakit.Num(lookup, 0), edakit.Num(compute, 0), edakit.Num(h, 0))
		ans := edakit.Num(exp, 2)
		return edakit.Numeric(q, ans, edakit.Sprintf("Every call pays the lookup; only misses pay the computation: %s + (1 - %s) × %s = %s s.", edakit.Num(lookup, 0), edakit.Num(h/100, 2), edakit.Num(compute, 0), ans),
			map[string]string{"expected": ans})
	case 2:
		lookup := float64(edakit.Int(3, 40))
		compute := float64(10 * edakit.Int(5, 40))
		be := lookup / compute * 100
		q := edakit.Sprintf("Caching a placement result costs %s s per lookup (hashing a large netlist and probing remote storage). Without the cache the step takes %s s; with it, a hit takes only the lookup and a miss takes the lookup plus the full %s s. Above what hit rate does the cache save time on average? Answer as a percentage, rounded to 2 decimal places.", edakit.Num(lookup, 0), edakit.Num(compute, 0), edakit.Num(compute, 0))
		ans := edakit.Num(be, 2)
		return edakit.Numeric(q, ans, edakit.Sprintf("Set the cached expectation equal to the uncached time: L + (1 - h) C = C gives h = L / C = %s / %s = %s%%.", edakit.Num(lookup, 0), edakit.Num(compute, 0), ans),
			map[string]string{"break_even_percent": ans})
	}
	type run struct{ net, sdc, tool, seed int }
	n := edakit.Scaled(ctx, 6, 9)
	var runs []run
	for i := 0; i < n; i++ {
		if i > 0 && edakit.Int(0, 2) == 0 {
			r := runs[edakit.Int(0, i-1)]
			runs = append(runs, r)
			continue
		}
		runs = append(runs, run{edakit.Int(1, 2), edakit.Int(1, 2), edakit.Int(1, 2), edakit.Int(1, 2)})
	}
	seen := map[run]bool{}
	seenNoTool := map[[3]int]bool{}
	hits, wrongHits := 0, 0
	var parts []string
	for i, r := range runs {
		parts = append(parts, edakit.Sprintf("run %d (netlist v%d, SDC v%d, tool v%d, seed %d)", i+1, r.net, r.sdc, r.tool, r.seed))
		if seen[r] {
			hits++
		}
		seen[r] = true
		k := [3]int{r.net, r.sdc, r.seed}
		if seenNoTool[k] {
			wrongHits++
		}
		seenNoTool[k] = true
	}
	q := edakit.Sprintf("A content-addressed cache keys each placement result by a hash of the netlist version, the SDC version, the tool version and the seed. Starting from an empty cache, the runs arrive in order: %s. How many runs are cache hits?", strings.Join(parts, "; "))
	return edakit.Numeric(q, edakit.IntStr(hits),
		edakit.Sprintf("A run hits only if an earlier run had all four inputs identical. There are %d distinct keys among %d runs, so the hit count is %d.", len(seen), n, hits)+incrToolNote(wrongHits-hits),
		map[string]string{"hits": edakit.IntStr(hits), "distinct_keys": edakit.IntStr(len(seen))})
}

var incrementalStaItems = []edakit.Item{
	{KP: 1, Q: "A resizer swaps a NAND2_X1 for a NAND2_X4 through the database API. How does an incremental timer embedded in the same process learn about the change?",
		A: "A database callback fires and the timer marks the affected pins invalid", D: []string{"It rereads the netlist file at the next report", "It re-times the full design on a timer every second", "The resizer must call a full update_timing after every swap"},
		Why: "Tools sharing one database register observers; the swap notifies the timer, which records what is now stale instead of recomputing on the spot."},
	{KP: 1, Q: "Which database edit should make an embedded timer invalidate arrival times downstream of a pin?",
		A: "Reconnecting an instance pin to a different net", D: []string{"Renaming a net without changing connectivity", "Changing a cell's placement status flag to locked", "Adding a comment property to the block"},
		Why: "Connectivity changes alter loads and arcs, so delays and arrivals downstream change. A rename or a lock flag changes no electrical fact."},
	{KP: 1, Q: "An in-process timer keeps its own graph of vertices for pins and edges for arcs. A net is deleted from the database. What must the timer's callback do before anything else touches the graph?",
		A: "Remove or detach the wire arcs of that net and mark its former loads invalid", D: []string{"Nothing, because nets are not timing objects", "Re-run levelization of the whole graph immediately", "Delete every vertex in the net's fanout cone"},
		Why: "The timer's graph must stay consistent with the database. Wire arcs of a deleted net would point at dead objects; the loads now see a different driver and are stale."},
	{KP: 1, Q: "A placer moves 50,000 cells in one detailed-placement pass, each move reported by a callback. What is the sensible way for the timer to respond?",
		A: "Accumulate invalidations and recompute once at the next timing query", D: []string{"Run a full incremental update after each of the 50,000 moves", "Ignore moves, since placement never changes delay", "Recompute only the slack of the moved cells' output pins"},
		Why: "Moves change wire parasitics and so delays. Batching the invalidations until a query amortizes the update across all moves."},
	{KP: 2, Q: "After a batch of edits, an incremental timer has a set of pins with invalid arrival times. In what order must it recompute them?",
		A: "In increasing level order, so every fanin is final before a pin is evaluated", D: []string{"In the order the edits arrived", "In decreasing level order from the endpoints", "In any order, because each arrival is computed independently"},
		Why: "An arrival is a max over fanin arrivals plus delays, so it can only be computed once its fanins are up to date; levelized order guarantees that."},
	{KP: 2, Q: "A query asks for the slack of one endpoint after an edit. Which values must be brought up to date for that answer?",
		A: "The invalid arrivals in its fanin cone and its own required time", D: []string{"Every invalid arrival and required time in the design", "Only the arrival at the endpoint itself", "The required times of every startpoint in its fanin cone"},
		Why: "Slack at an endpoint is its required minus its arrival. The arrival depends on the fanin cone; the endpoint's required time comes from its clock and constraints."},
	{KP: 2, Q: "An edit shortens an arc delay, and the new arrival at its head pin equals the old one because another fanin dominates the max. What should the timer do downstream?",
		A: "Stop propagating along that pin, since nothing downstream can change through it", D: []string{"Continue invalidating the whole fanout cone to be safe", "Re-time the whole design because maxes are non-monotone", "Recompute the required times in the fanout cone"},
		Why: "If a recomputed value is unchanged, its dependents see identical inputs. Early termination is what keeps incremental timing much cheaper than the cone bound."},
	{KP: 2, Q: "Why does an incremental timer track required times separately from arrival times when invalidating?",
		A: "Arrivals depend on the fanin cone, required times on the fanout cone, so an arc edit invalidates them in opposite directions", D: []string{"Required times never change after the SDC is read", "Arrival times are stored in the database and required times are not", "Both flow forward, so one set of flags would duplicate the other"},
		Why: "An arc change affects arrivals at and after its head, and required times at and before its tail, so the two sets of flags cover different pins."},
	{KP: 3, Q: "Which edit is typically the most expensive for an incremental timer to absorb?",
		A: "Changing a clock period in the constraints", D: []string{"Upsizing one cell on a critical path", "Inserting one buffer on a high-fanout net", "Moving one cell by a few sites"},
		Why: "A clock change alters required times at every endpoint of that clock and therefore across most of the design, whereas a netlist edit has a local cone."},
	{KP: 3, Q: "OpenTimer differs from a classic levelized incremental update in how it executes the recomputation. What does it use?",
		A: "A task graph of the affected work, run in parallel by a work-stealing scheduler", D: []string{"A complete re-timing on a GPU after every edit", "A SAT solver to find which arcs changed", "Monte Carlo sampling of the affected paths"},
		Why: "OpenTimer builds a task dependency graph for the invalid region and executes it with a parallel task scheduler, so independent parts update concurrently."},
	{KP: 3, Q: "An optimizer alternates one edit and one full worst-slack query, 10,000 times, and each query touches most of the design. What limits the speed of this loop?",
		A: "The query pattern: asking for global WNS after each edit forces a near-global update", D: []string{"Callback overhead, which dominates any timing work", "Levelization, which must be redone from scratch for each edit", "Nothing; incremental timing makes every query local"},
		Why: "Incrementality only pays when queries are local. A loop that asks a global question after each small edit should ask local questions, such as the slack of the affected endpoints, or batch edits."},
	{KP: 3, Q: "A buffer insertion adds a new instance and splits a net. What must happen to the timing graph's levels?",
		A: "The new pins need levels, and pins downstream may need their levels raised", D: []string{"Levels never change after the first levelization", "Every level in the design must be renumbered from scratch", "Only the buffer's input pin needs a level"},
		Why: "The buffer lengthens the path, so the levels of its fanout may rise. Timers update levels incrementally for the affected region rather than relevelizing everything."},
}

// incrToolNote explains what a key without the tool version would have done.
func incrToolNote(extra int) string {
	if extra <= 0 {
		return " Here no two runs differ only in tool version, but a key that omitted it would be wrong in general."
	}
	return edakit.Sprintf(" A key that left out the tool version would add %d false hit(s), each silently returning a result produced by a different tool build.", extra)
}
