package sta

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// pickInnerArc picks an arc whose head is not a sink and whose tail is not a source when possible,
// so both cones are non-trivial.
func pickInnerArc(g *dag) arc {
	var pool []arc
	for _, a := range g.arcs {
		if len(g.fanin(a.u)) > 0 && len(g.fanout(a.v)) > 0 {
			pool = append(pool, a)
		}
	}
	if len(pool) == 0 {
		for _, a := range g.arcs {
			if len(g.fanin(a.u)) > 0 || len(g.fanout(a.v)) > 0 {
				pool = append(pool, a)
			}
		}
	}
	if len(pool) == 0 {
		pool = g.arcs
	}
	return edakit.Pick(pool)
}

func coneNames(g *dag, m map[int]bool) string {
	var vs []int
	for v := 0; v < g.n; v++ {
		if m[v] {
			vs = append(vs, v)
		}
	}
	return g.nameList(vs)
}

func incrementalUpdate(ctx generator.GeneratorContext) generator.Problem {
	g := newDAG(edakit.Scaled(ctx, 7, 11), edakit.Int(1, 2), 1, 1)
	a := pickInnerArc(g)
	fwd := g.forwardCone(a.v)
	bwd := g.backwardCone(a.u)
	q := edakit.Sprintf("A timing graph has the arcs %s. An ECO resizes a cell, and only the delay of arc %s→%s changes. ", g.describe(false), g.names[a.u], g.names[a.v])
	facts := map[string]string{"forward": edakit.IntStr(countTrue(fwd)), "backward": edakit.IntStr(countTrue(bwd))}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+"How many nodes can have a changed arrival time (count every node whose arrival must be re-evaluated)?", edakit.IntStr(countTrue(fwd)),
			edakit.Sprintf("Arrivals flow forward, so only %s and its fanout cone can change: %s (%d in all).", g.names[a.v], coneNames(g, fwd), countTrue(fwd)), facts)
	case 2:
		return edakit.Numeric(q+"How many nodes can have a changed required time?", edakit.IntStr(countTrue(bwd)),
			edakit.Sprintf("Required times flow backward, so only %s and its fanin cone can change: %s (%d in all).", g.names[a.u], coneNames(g, bwd), countTrue(bwd)), facts)
	}
	union := map[int]bool{}
	for v := range fwd {
		union[v] = true
	}
	for v := range bwd {
		union[v] = true
	}
	facts["union"] = edakit.IntStr(countTrue(union))
	return edakit.Numeric(q+"How many nodes can have a changed slack?", edakit.IntStr(countTrue(union)),
		edakit.Sprintf("Slack changes where arrival or required time changes: the fanout cone of %s united with the fanin cone of %s, %s, which is %d of %d nodes.", g.names[a.v], g.names[a.u], coneNames(g, union), countTrue(union), g.n), facts)
}

func earlyTermination(ctx generator.GeneratorContext) generator.Problem {
	var g *dag
	var a arc
	var before, after []int
	var src map[int]int
	var srcText string
	var delta int
	var idx int
	for tries := 0; ; tries++ {
		g = newDAG(edakit.Scaled(ctx, 7, 10), edakit.Int(1, 2), 5, 40)
		src, srcText = sourceArrivals(g)
		before = g.arrivals(src, true)
		a = pickInnerArc(g)
		for i, x := range g.arcs {
			if x == a {
				idx = i
			}
		}
		delta = 5 * edakit.Int(1, 4)
		if edakit.Int(0, 1) == 0 {
			delta = -delta
		}
		g.arcs[idx].d += delta
		after = g.arrivals(src, true)
		g.arcs[idx].d -= delta
		changed := 0
		for v := range after {
			if after[v] != before[v] {
				changed++
			}
		}
		if (changed > 0 && changed < countTrue(g.forwardCone(a.v))) || tries > 30 {
			break
		}
	}
	changed := 0
	for v := range after {
		if after[v] != before[v] {
			changed++
		}
	}
	q := edakit.Sprintf("A timing graph has the arcs %s, and sources launch %s. An ECO changes the delay of arc %s→%s from %d ps to %d ps. ", g.describe(true), srcText, g.names[a.u], g.names[a.v], a.d, a.d+delta)
	exp := edakit.Sprintf("Arrivals before: %s. After: %s. Propagation stops wherever a node's new arrival equals its old one, because a max that did not change hides the edit from everything downstream.",
		g.valueList(allNodes(g.n), before), g.valueList(allNodes(g.n), after))
	facts := map[string]string{"changed": edakit.IntStr(changed)}
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(q+edakit.Sprintf("What is the new arrival time at %s, in ps?", g.names[a.v]), edakit.IntStr(after[a.v]), exp, facts)
	case 3:
		var pool []int
		cone := g.forwardCone(a.v)
		for v := 0; v < g.n; v++ {
			if cone[v] && len(g.fanout(v)) == 0 {
				pool = append(pool, v)
			}
		}
		if len(pool) == 0 {
			pool = []int{a.v}
		}
		v := edakit.Pick(pool)
		return edakit.Numeric(q+edakit.Sprintf("What is the new arrival time at endpoint %s, in ps?", g.names[v]), edakit.IntStr(after[v]), exp, facts)
	}
	return edakit.Numeric(q+"How many nodes end up with a different arrival time?", edakit.IntStr(changed), exp+edakit.Sprintf(" Changed arrivals: %d; nodes in the fanout cone: %d.", changed, countTrue(g.forwardCone(a.v))), facts)
}
