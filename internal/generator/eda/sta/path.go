package sta

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func pathCount(ctx generator.GeneratorContext) generator.Problem {
	if ctx.KPOrdinal == 3 {
		n := edakit.Scaled(ctx, 3, 8)
		var ws []string
		total := 1
		for i := 0; i < n; i++ {
			w := edakit.Int(2, 3)
			total *= w
			ws = append(ws, edakit.IntStr(w))
		}
		q := edakit.Sprintf("A datapath is a cascade of %d reconvergent stages. In each stage the signal fans out into parallel branches that rejoin at a single gate; the branch counts are %s. How many distinct input-to-output paths does the cascade have?", n, strings.Join(ws, ", "))
		return edakit.Numeric(q, edakit.IntStr(total), edakit.Sprintf("Each stage multiplies the number of paths by its branch count: %s = %d. Path counts grow exponentially with depth, while the graph grows linearly.", strings.Join(ws, " x "), total), map[string]string{"paths": edakit.IntStr(total)})
	}
	g := newDAG(edakit.Scaled(ctx, 6, 10), edakit.Int(1, 2), 1, 1)
	pc := g.pathCounts()
	intro := "A timing graph has the arcs " + g.describe(false) + ". "
	if ctx.KPOrdinal == 2 {
		pool := g.internal()
		if len(pool) == 0 {
			pool = []int{g.n - 1}
		}
		v := edakit.Pick(pool)
		out := make([]int, g.n) // paths from node to any sink
		for w := g.n - 1; w >= 0; w-- {
			fo := g.fanout(w)
			if len(fo) == 0 {
				out[w] = 1
				continue
			}
			for _, a := range fo {
				out[w] += out[a.v]
			}
		}
		ans := pc[v] * out[v]
		return edakit.Numeric(intro+edakit.Sprintf("How many source-to-sink paths pass through node %s?", g.names[v]), edakit.IntStr(ans),
			edakit.Sprintf("Paths from the sources into %s: %d. Paths from %s to the sinks: %d. Every pairing is a distinct path: %d.", g.names[v], pc[v], g.names[v], out[v], ans), map[string]string{"in": edakit.IntStr(pc[v]), "out": edakit.IntStr(out[v])})
	}
	v := g.sinks()[0]
	for _, w := range g.sinks() {
		if pc[w] > pc[v] {
			v = w
		}
	}
	return edakit.Numeric(intro+edakit.Sprintf("How many distinct paths run from the sources to node %s?", g.names[v]), edakit.IntStr(pc[v]),
		edakit.Sprintf("Count in topological order: a source has one path and every other node sums its fanins. Counts: %s.", g.valueList(allNodes(g.n), pc)), map[string]string{"paths": edakit.IntStr(pc[v])})
}

func kWorst(ctx generator.GeneratorContext) generator.Problem {
	var g *dag
	var v int
	var ds []int
	for {
		g = newDAG(edakit.Scaled(ctx, 6, 9), 1, 5, 40)
		v = edakit.Pick(g.sinks())
		ds = g.pathDelays(v)
		if len(ds) >= 3 {
			break
		}
	}
	q := edakit.Sprintf("A timing graph with a single source %s has the arcs %s. ", g.names[0], g.describe(true))
	var all []string
	for _, d := range ds {
		all = append(all, edakit.IntStr(d))
	}
	exp := edakit.Sprintf("The path delays into %s, from worst down, are %s.", g.names[v], strings.Join(all, ", "))
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+edakit.Sprintf("What is the delay of the worst (longest) path from %s to %s, in ps?", g.names[0], g.names[v]), edakit.IntStr(ds[0]), exp, map[string]string{"worst": edakit.IntStr(ds[0])})
	case 3:
		delta := 5 * edakit.Int(2, 6)
		cnt := 0
		for _, d := range ds {
			if d >= ds[0]-delta {
				cnt++
			}
		}
		return edakit.Numeric(q+edakit.Sprintf("How many distinct paths from %s to %s have a delay within %d ps of the worst path (the worst path included)?", g.names[0], g.names[v], delta), edakit.IntStr(cnt),
			exp+edakit.Sprintf(" Those at or above %d ps number %d.", ds[0]-delta, cnt), map[string]string{"count": edakit.IntStr(cnt)})
	}
	return edakit.Numeric(q+edakit.Sprintf("Paths that tie count separately. What is the delay of the second-worst path from %s to %s, in ps?", g.names[0], g.names[v]), edakit.IntStr(ds[1]), exp, map[string]string{"second": edakit.IntStr(ds[1])})
}

func pathReport(ctx generator.GeneratorContext) generator.Problem {
	t := 50 * edakit.Int(10, 24)
	lclk := 5 * edakit.Int(20, 50)
	cclk := 5 * edakit.Int(20, 50)
	tcq := edakit.Int(40, 90)
	n := edakit.Int(3, 5)
	cells := []string{"NAND2_X1", "INV_X2", "AOI21_X1", "NOR2_X1", "BUF_X4", "OAI22_X1", "XOR2_X1"}
	var lines []string
	arr := lclk + tcq
	lines = append(lines, edakit.Sprintf("clock clk (rise edge) 0"), edakit.Sprintf("clock network delay (propagated) %d", lclk), edakit.Sprintf("FF1/CK -> FF1/Q (DFF_X1) %d", tcq))
	for i := 0; i < n; i++ {
		net := edakit.Int(2, 25)
		cell := edakit.Int(20, 90)
		arr += net + cell
		lines = append(lines, edakit.Sprintf("U%d (%s) %d + %d (net + cell)", i+1, edakit.Pick(cells), net, cell))
	}
	last := edakit.Int(2, 25)
	arr += last
	lines = append(lines, edakit.Sprintf("FF2/D net %d", last))
	unc := 5 * edakit.Int(2, 12)
	su := edakit.Int(20, 60)
	cr := edakit.Int(0, 25)
	req := t + cclk - unc + cr - su
	lines = append(lines, "---", edakit.Sprintf("clock clk (rise edge) %d", t), edakit.Sprintf("clock network delay (propagated) %d", cclk), edakit.Sprintf("clock reconvergence pessimism %d", cr), edakit.Sprintf("clock uncertainty -%d", unc), edakit.Sprintf("library setup time -%d", su))
	q := "A setup path report lists these increments in ps:\n\n" + strings.Join(lines, "\n") + "\n\n"
	slack := req - arr
	facts := map[string]string{"arrival": edakit.IntStr(arr), "required": edakit.IntStr(req), "slack": edakit.IntStr(slack)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+"What is the data arrival time at FF2/D, in ps?", edakit.IntStr(arr),
			edakit.Sprintf("Add the launch clock network delay, clock-to-Q and every net and cell increment: %d ps.", arr), facts)
	case 3:
		return edakit.Numeric(q+"What is the slack, in ps?", edakit.IntStr(slack),
			edakit.Sprintf("Required %d ps minus arrival %d ps = %d ps.", req, arr, slack), facts)
	}
	return edakit.Numeric(q+"What is the data required time, in ps?", edakit.IntStr(req),
		edakit.Sprintf("Capture edge %d + clock network %d + CPPR %d - uncertainty %d - setup %d = %d ps.", t, cclk, cr, unc, su, req), facts)
}

var gbaItems = []edakit.Item{
	{KP: 1, Q: "At a multi-input gate, a graph-based timer stores one late slew at the output. Which slew does it keep?",
		A: "The worst slew over all input arcs, even if it comes from a different input than the latest arrival", D: []string{"The slew of the latest-arriving input only", "The average of the input slews", "The slew of the fastest arc"},
		Why: "GBA merges arrival and slew separately, each taking its worst value. The pair may not come from one real path, which is a source of pessimism."},
	{KP: 1, Q: "Why is graph-based analysis pessimistic but never optimistic for setup?",
		A: "Each node keeps the worst arrival and worst slew, which bound every path through it from above", D: []string{"It ignores wire delay", "It uses typical-corner libraries", "It removes clock reconvergence pessimism twice"},
		Why: "Delays grow with input slew, so propagating the worst slew can only overestimate a path's delay. The bound is safe but loose."},
	{KP: 1, Q: "A design has billions of paths but only a few million pins. Which analysis can time the whole design in one linear pass?",
		A: "Graph-based analysis", D: []string{"Path-based analysis of every path", "Monte Carlo SPICE of every path", "Exhaustive path enumeration with exact slews"},
		Why: "GBA keeps one value per pin and transition, so its work grows with the graph, not with the exponential number of paths."},
	{KP: 1, Q: "AOCV derates depend on logic depth. In graph-based mode, which depth does a pin's derate use?",
		A: "A conservative depth bound for the pin, such as the shallowest path through it", D: []string{"The exact depth of the path being reported", "The average depth of the design", "Always depth 1"},
		Why: "One pin lies on paths of many depths. GBA must pick a value safe for all of them; PBA uses the actual depth of the path it re-times."},
	{KP: 2, Q: "A path shows -12 ps under GBA and +4 ps under PBA. What should the closure flow conclude?",
		A: "The path meets timing; the GBA violation was pessimism from merged slews or depth", D: []string{"PBA is optimistic, so the path still fails", "The two analyses disagree because the netlist is wrong", "Both numbers must be averaged"},
		Why: "PBA re-times the specific path with its own slews and depth, which removes pessimism. Its result is still a safe bound for that path."},
	{KP: 2, Q: "Which relation between the PBA and GBA slack of the same path always holds?",
		A: "PBA slack is greater than or equal to GBA slack", D: []string{"PBA slack is less than or equal to GBA slack", "They are always equal", "PBA slack is always positive"},
		Why: "PBA removes pessimism from the same bound, so it can only raise slack."},
	{KP: 2, Q: "Why does a timing optimizer usually run on GBA rather than PBA?",
		A: "GBA updates incrementally in time proportional to the changed cone, while PBA must re-time paths one by one", D: []string{"PBA does not support hold checks", "GBA is exact while PBA is approximate", "Optimizers cannot read PBA results"},
		Why: "Optimization makes thousands of changes. GBA's per-pin values update cheaply; PBA is reserved for recovering pessimism on the worst paths."},
	{KP: 2, Q: "A signoff flow runs PBA only on paths whose GBA slack is below a threshold. Why is that safe?",
		A: "Paths that pass under GBA also pass under PBA, so only failing or near-failing paths need re-timing", D: []string{"PBA is wrong for passing paths", "GBA is optimistic for passing paths", "Passing paths are removed from the netlist"},
		Why: "GBA bounds every path. Re-timing can only improve slack, so a path already passing under GBA cannot fail under PBA."},
	{KP: 3, Q: "A timer reports the 10 worst paths to one endpoint, and all are near-duplicates through the same cone. Why?",
		A: "Path enumeration ranks by slack, and reconvergent logic yields many paths that differ by one segment", D: []string{"The timer has a bug in its sort", "Only ten paths exist", "Duplicates come from the SDC"},
		Why: "With reconvergence, the k worst paths often share most of their segments. Reports limit paths per endpoint to show distinct problems."},
	{KP: 3, Q: "Exhaustive PBA over every path to an endpoint would be exact. Why do tools not run it by default?",
		A: "The number of paths grows exponentially with reconvergent depth", D: []string{"PBA cannot handle more than one path", "PBA ignores derates", "PBA does not use Liberty data"},
		Why: "Enumerating every path is exponential. Tools enumerate the worst paths in slack order and stop at a limit or a slack threshold."},
	{KP: 3, Q: "How can the GBA arrival at a pin exceed the arrival of every real path into it?",
		A: "The delay of a downstream arc is computed from the latest arrival paired with a worse slew that belongs to a different path", D: []string{"GBA adds a fixed margin at every pin", "GBA counts each path twice", "It cannot: the GBA arrival always equals the worst real path"},
		Why: "Arrival and slew are maxed separately. A downstream arc then sees a combination no single path produces, so its delay and the arrival can exceed every real path."},
	{KP: 3, Q: "After PBA recovers 15 ps on a path, an ECO changes a cell on it. What must happen to the PBA result?",
		A: "It is stale and must be recomputed, because slews and depths along the path may have changed", D: []string{"It stays valid because PBA is exact", "It is shifted by the cell's delay change only", "It becomes the new GBA value automatically"},
		Why: "PBA is computed for a specific path state. Any change upstream can alter slews, so the re-timing must be repeated."},
}
