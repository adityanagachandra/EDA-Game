package sta

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// timingGraphCount draws a small netlist and asks how many arcs its timing graph has.
func timingGraphCount(ctx generator.GeneratorContext) generator.Problem {
	type cell struct {
		name string
		k    int
	}
	lib := []cell{{"INV", 1}, {"BUF", 1}, {"NAND2", 2}, {"NOR2", 2}, {"XOR2", 2}, {"NAND3", 3}, {"AOI21", 3}, {"OAI22", 4}}
	ng := edakit.Scaled(ctx, 2, 5)
	np := edakit.Int(1, 3)
	var parts []string
	cellArcs, netArcs := 0, 0
	for i := 0; i < ng; i++ {
		c := edakit.Pick(lib)
		f := edakit.Int(1, 4)
		cellArcs += c.k
		netArcs += f
		parts = append(parts, edakit.Sprintf("cell U%d (%s) drives a net of fanout %d", i+1, c.name, f))
	}
	for i := 0; i < np; i++ {
		f := edakit.Int(1, 3)
		netArcs += f
		parts = append(parts, edakit.Sprintf("input port in%d drives a net of fanout %d", i+1, f))
	}
	kind := edakit.Pick([]string{"total", "cell", "net"})
	var ans int
	switch kind {
	case "cell":
		ans = cellArcs
	case "net":
		ans = netArcs
	default:
		ans = cellArcs + netArcs
	}
	what := map[string]string{"total": "arcs (cell arcs plus net arcs)", "cell": "cell arcs", "net": "net arcs"}[kind]
	q := "Build the timing graph of this fragment, with one arc per input-to-output pin pair of a cell and one arc per driver-to-sink pin pair of a net (ignore rise and fall). " +
		strings.Join(parts, "; ") + edakit.Sprintf(". Every cell has a single output. How many %s does the fragment contribute?", what)
	exp := edakit.Sprintf("A single-output cell with k inputs contributes k cell arcs, here %d in all. A net with f sinks contributes f net arcs, here %d in all. The requested count is %d.", cellArcs, netArcs, ans)
	d := edakit.Distinct(ans, 4, cellArcs, netArcs, cellArcs+netArcs, 2*ans, cellArcs+netArcs+ng)
	return edakit.Choice(q, edakit.IntStr(ans), exp, d...)
}

var timingGraphItems = []edakit.Item{
	{KP: 1, Q: "In the timing graph a timer such as OpenSTA builds, what is a node?",
		A: "A pin of a cell instance or a top-level port", D: []string{"A cell instance as a whole", "A net as a whole", "A clock domain"},
		Why: "Nodes are pins and ports. Cell arcs connect the input and output pins of one instance, and net arcs connect a driver pin to each sink pin."},
	{KP: 1, Q: "Which of these is a timing startpoint?",
		A: "The clock pin of a flip-flop", D: []string{"The D pin of a flip-flop", "An output port", "The output pin of an inverter in the middle of a path"},
		Why: "Paths start at input ports and at sequential clock pins, where the clock-to-Q arc launches data. D pins and output ports are endpoints."},
	{KP: 1, Q: "Where do the delays on cell arcs come from?",
		A: "The Liberty timing arcs of each cell, evaluated at the pin's slew and load", D: []string{"The SPEF file alone", "The SDC file", "The Verilog netlist"},
		Why: "Liberty describes each input-to-output timing arc with delay and slew tables. Parasitics feed the load, and SDC sets clocks and constraints."},
	{KP: 1, Q: "Why does the timing graph stop at a flip-flop's D pin rather than continuing through to Q?",
		A: "D to Q is not a propagation arc; Q is launched by the clock edge, so D is an endpoint with a check arc", D: []string{"Flip-flops have no Liberty data", "Q always has zero delay", "The D pin is always a false path"},
		Why: "A flop samples D at the clock edge and launches Q from the clock. The graph models a setup and hold check at D and a clock-to-Q arc from CK to Q."},
	{KP: 3, Q: "A netlist contains a combinational loop through two NAND gates. What must the timer do before propagating arrivals?",
		A: "Break the loop by disabling one arc so the graph is acyclic", D: []string{"Iterate arrival times until they converge", "Report an infinite arrival time at every pin", "Treat the loop as a multicycle path"},
		Why: "Propagation in topological order needs a DAG. Timers break combinational loops by disabling an arc and warn about it."},
	{KP: 3, Q: "Why does a block-based timer propagate in topological order?",
		A: "Each pin's arrival then depends only on fanins that are already final, so one pass suffices", D: []string{"It finds the shortest path first", "It keeps memory constant", "It visits every path explicitly"},
		Why: "In topological order every fanin is complete before its fanout is evaluated, so the whole graph is timed in time linear in its arcs."},
	{KP: 3, Q: "Why does a timer keep separate rise and fall arrival times at every pin?",
		A: "Rising and falling transitions have different delays, and inverting arcs swap them", D: []string{"To double-check each other", "Because Liberty stores only rise data", "To model setup and hold separately"},
		Why: "Delay tables differ for rise and fall, and a negative-unate arc maps an input rise to an output fall, so both must be tracked."},
	{KP: 3, Q: "A designer marks an arc with set_disable_timing. What happens to the timing graph?",
		A: "That arc is removed, so no path propagates through it", D: []string{"The arc's delay becomes zero", "The arc becomes a false path only for hold", "The cell is removed from the netlist"},
		Why: "Disabling an arc deletes it from propagation. Paths that relied on it disappear from reports, which is why disabled arcs must be audited."},
	{KP: 3, Q: "How does the size of the timing graph grow with the netlist?",
		A: "Linearly in the number of pins and arcs", D: []string{"Exponentially in the number of paths", "Quadratically in the number of cells", "With the number of clock domains only"},
		Why: "The graph holds one node per pin and one edge per arc. It is the number of paths, not the graph, that grows exponentially."},
}

// unateness propagates a transition through a chain of gates and asks what can reach the end.
func unateness(ctx generator.GeneratorContext) generator.Problem {
	type g struct {
		name  string
		sense int // +1 positive, -1 negative, 0 non-unate
	}
	lib := []g{{"buffer", 1}, {"inverter", -1}, {"NAND2", -1}, {"NOR2", -1}, {"AND2", 1}, {"OR2", 1}, {"XOR2", 0}}
	n := 1
	if ctx.KPOrdinal != 1 {
		n = edakit.Scaled(ctx, 2, 5)
	}
	var chain []g
	for i := 0; i < n; i++ {
		x := edakit.Pick(lib)
		if ctx.KPOrdinal == 2 && x.sense == 0 {
			x = lib[edakit.Int(0, 5)]
		}
		chain = append(chain, x)
	}
	if ctx.KPOrdinal == 3 {
		chain[edakit.Int(0, n-1)] = lib[6]
	}
	rise := edakit.Pick([]bool{true, false})
	// state: 1 rise only, -1 fall only, 0 both
	state := -1
	if rise {
		state = 1
	}
	var names []string
	for _, x := range chain {
		names = append(names, x.name)
		switch {
		case state == 0:
		case x.sense == 0:
			state = 0
		default:
			state *= x.sense
		}
	}
	edge := "falling"
	if rise {
		edge = "rising"
	}
	opts := map[int]string{1: "A rising transition only", -1: "A falling transition only", 0: "Both a rising and a falling transition"}
	q := edakit.Sprintf("A %s transition enters a chain of gates in this order: %s. It passes through one input of each gate, and side inputs may hold any value. Which output transitions must the timer propagate at the end of the chain?",
		edge, strings.Join(names, ", then "))
	exp := "Positive-unate arcs keep the direction, negative-unate arcs flip it, and a non-unate arc such as XOR can produce either, after which both directions travel on. The result is: " + strings.ToLower(opts[state]) + "."
	var d []string
	for _, k := range []int{1, -1, 0} {
		if k != state {
			d = append(d, opts[k])
		}
	}
	d = append(d, "No transition, because the side inputs block it")
	return edakit.Choice(q, opts[state], exp, d...)
}

// levelization asks for the level of a node, the depth, or the width of a level.
func levelization(ctx generator.GeneratorContext) generator.Problem {
	g := newDAG(edakit.Scaled(ctx, 6, 10), edakit.Int(1, 2), 1, 1)
	lv := g.levels()
	exp := "Sources sit at level 0 and every other node at one more than its deepest fanin. Levels: " + g.valueList(allNodes(g.n), lv) + "."
	intro := "A timing graph has the arcs " + g.describe(false) + ". Sources are at level 0, and every other node is at one plus the largest level among its fanins. "
	switch ctx.KPOrdinal {
	case 2:
		depth := maxInt(lv)
		return edakit.Numeric(intro+"How many levels deep is the graph, that is, what is the largest level?", edakit.IntStr(depth), exp, map[string]string{"depth": edakit.IntStr(depth)})
	case 3:
		target := lv[edakit.Int(1, g.n-1)]
		cnt := 0
		for _, l := range lv {
			if l == target {
				cnt++
			}
		}
		return edakit.Numeric(intro+edakit.Sprintf("How many nodes are at level %d?", target), edakit.IntStr(cnt), exp, map[string]string{"count": edakit.IntStr(cnt)})
	}
	v := edakit.Int(g.n/2, g.n-1)
	return edakit.Numeric(intro+edakit.Sprintf("What is the level of node %s?", g.names[v]), edakit.IntStr(lv[v]), exp, map[string]string{"level": edakit.IntStr(lv[v])})
}

func sourceArrivals(g *dag) (map[int]int, string) {
	src := map[int]int{}
	var parts []string
	for _, s := range g.sources() {
		src[s] = 5 * edakit.Int(0, 6)
		parts = append(parts, edakit.Sprintf("%s at %d ps", g.names[s], src[s]))
	}
	return src, strings.Join(parts, ", ")
}

// arrivalTime propagates late (or early, at KP 3) arrivals on a random DAG.
func arrivalTime(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 9)
	if ctx.KPOrdinal == 1 {
		n = edakit.Int(4, 5)
	}
	g := newDAG(n, edakit.Int(1, 2), 5, 40)
	src, srcText := sourceArrivals(g)
	late := ctx.KPOrdinal != 3
	at := g.arrivals(src, late)
	sinks := g.sinks()
	v := edakit.Pick(sinks)
	word, op := "latest", "max"
	if !late {
		word, op = "earliest", "min"
	}
	q := edakit.Sprintf("A timing graph has the arcs %s. Signals leave the sources %s. What is the %s arrival time at node %s, in ps?", g.describe(true), srcText, word, g.names[v])
	exp := edakit.Sprintf("Visit nodes in topological order and take the %s over fanins of arrival plus arc delay. Arrivals: %s. At %s the %s arrival is %d ps.", op, g.valueList(allNodes(g.n), at), g.names[v], word, at[v])
	return edakit.Numeric(q, edakit.IntStr(at[v]), exp, map[string]string{"arrival": edakit.IntStr(at[v])})
}

func endpointRequired(g *dag) (map[int]int, string) {
	req := map[int]int{}
	var parts []string
	for _, s := range g.sinks() {
		req[s] = 10 * edakit.Int(8, 16)
		parts = append(parts, edakit.Sprintf("%s at %d ps", g.names[s], req[s]))
	}
	return req, strings.Join(parts, ", ")
}

// arrivalRequired asks for an endpoint required time (KP 1) or a backward-propagated one.
func arrivalRequired(ctx generator.GeneratorContext) generator.Problem {
	if ctx.KPOrdinal == 1 {
		t := 50 * edakit.Int(10, 30)
		lat := 5 * edakit.Int(0, 30)
		su := edakit.Int(15, 60)
		ans := t + lat - su
		q := edakit.Sprintf("A capture flip-flop is clocked with period %d ps. Its clock pin sees the capturing edge %d ps late (clock latency), and its setup time is %d ps. Data launches from a flop on the edge at 0 ps with ideal latency. What is the required time at the capture flop's D pin, in ps?", t, lat, su)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Required time is the capturing edge plus its latency minus setup: %d + %d - %d = %d ps.", t, lat, su, ans), map[string]string{"required": edakit.IntStr(ans)})
	}
	g := newDAG(edakit.Scaled(ctx, 6, 9), edakit.Int(1, 2), 5, 40)
	req, reqText := endpointRequired(g)
	rat := g.required(req)
	var pool []int
	if ctx.KPOrdinal == 3 {
		pool = g.sources()
	} else {
		pool = g.internal()
		if len(pool) == 0 {
			pool = g.sources()
		}
	}
	v := edakit.Pick(pool)
	q := edakit.Sprintf("A timing graph has the arcs %s. The endpoints have required times %s. Propagating backward, what is the required time at node %s, in ps?", g.describe(true), reqText, g.names[v])
	exp := edakit.Sprintf("Visit nodes in reverse topological order and take the minimum over fanouts of required time minus arc delay, because the tightest downstream endpoint governs. Required times: %s.", g.valueList(allNodes(g.n), rat))
	return edakit.Numeric(q, edakit.IntStr(rat[v]), exp, map[string]string{"required": edakit.IntStr(rat[v])})
}

// nodeSlack asks for a slack at an endpoint, at an internal node, or the size of the critical set.
func nodeSlack(ctx generator.GeneratorContext) generator.Problem {
	g := newDAG(edakit.Scaled(ctx, 6, 9), edakit.Int(1, 2), 5, 40)
	src, srcText := sourceArrivals(g)
	at := g.arrivals(src, true)
	req := map[int]int{}
	var parts []string
	for _, s := range g.sinks() {
		req[s] = at[s] + 5*edakit.Int(-6, 6)
		parts = append(parts, edakit.Sprintf("%s at %d ps", g.names[s], req[s]))
	}
	rat := g.required(req)
	slack := make([]int, g.n)
	for v := range slack {
		slack[v] = rat[v] - at[v]
	}
	intro := edakit.Sprintf("A timing graph has the arcs %s. Sources launch %s. Endpoint required times are %s. ", g.describe(true), srcText, strings.Join(parts, ", "))
	exp := edakit.Sprintf("Arrivals (forward max): %s. Required times (backward min): %s. Slack is required minus arrival: %s.",
		g.valueList(allNodes(g.n), at), g.valueList(allNodes(g.n), rat), g.valueList(allNodes(g.n), slack))
	switch ctx.KPOrdinal {
	case 1:
		v := edakit.Pick(g.sinks())
		return edakit.Numeric(intro+edakit.Sprintf("What is the slack at endpoint %s, in ps?", g.names[v]), edakit.IntStr(slack[v]), exp, map[string]string{"slack": edakit.IntStr(slack[v])})
	case 3:
		w := minInt(slack)
		cnt := 0
		for _, s := range slack {
			if s == w {
				cnt++
			}
		}
		return edakit.Numeric(intro+"How many nodes have a slack equal to the worst slack in the graph (the nodes on the critical path or paths)?", edakit.IntStr(cnt), exp, map[string]string{"worst": edakit.IntStr(w), "count": edakit.IntStr(cnt)})
	}
	pool := g.internal()
	if len(pool) == 0 {
		pool = g.sources()
	}
	v := edakit.Pick(pool)
	return edakit.Numeric(intro+edakit.Sprintf("What is the slack at node %s, in ps?", g.names[v]), edakit.IntStr(slack[v]), exp, map[string]string{"slack": edakit.IntStr(slack[v])})
}

// wnsTNS draws endpoint slacks and asks for WNS, TNS, or TNS after a partial fix.
func wnsTNS(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Scaled(ctx, 5, 10)
	sl := make([]int, m)
	neg := false
	for i := range sl {
		sl[i] = edakit.Int(-60, 70)
		if sl[i] < 0 {
			neg = true
		}
	}
	if !neg {
		sl[edakit.Int(0, m-1)] = -edakit.Int(1, 40)
	}
	var parts []string
	for i, s := range sl {
		parts = append(parts, edakit.Sprintf("E%d: %d", i+1, s))
	}
	intro := "The setup slacks of a design's endpoints, in ps, are " + strings.Join(parts, ", ") + ". "
	tns := func(xs []int) int {
		t := 0
		for _, x := range xs {
			if x < 0 {
				t += x
			}
		}
		return t
	}
	switch ctx.KPOrdinal {
	case 1:
		w := minInt(sl)
		return edakit.Numeric(intro+"What is the worst negative slack (WNS), in ps?", edakit.IntStr(w),
			edakit.Sprintf("WNS is the most negative endpoint slack, %d ps.", w), map[string]string{"wns": edakit.IntStr(w)})
	case 3:
		x := 5 * edakit.Int(2, 8)
		k := edakit.Int(2, m)
		idx := edakit.Shuffle(allNodes(m))[:k]
		sort.Ints(idx)
		after := append([]int(nil), sl...)
		var names []string
		for _, i := range idx {
			after[i] += x
			names = append(names, edakit.Sprintf("E%d", i+1))
		}
		t := tns(after)
		return edakit.Numeric(intro+edakit.Sprintf("An ECO raises the slack of %s by %d ps each and leaves the rest unchanged. What is the new total negative slack (TNS), in ps?", strings.Join(names, ", "), x), edakit.IntStr(t),
			edakit.Sprintf("Add %d ps to the listed endpoints, then sum only the slacks that are still negative. Endpoints pushed above zero contribute nothing. The new TNS is %d ps (it was %d ps).", x, t, tns(sl)),
			map[string]string{"tns_before": edakit.IntStr(tns(sl)), "tns_after": edakit.IntStr(t)})
	}
	t := tns(sl)
	return edakit.Numeric(intro+"What is the total negative slack (TNS), in ps?", edakit.IntStr(t),
		edakit.Sprintf("TNS sums only the negative slacks; positive ones contribute zero. The sum is %d ps.", t), map[string]string{"tns": edakit.IntStr(t)})
}
