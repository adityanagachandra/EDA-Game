package hdl

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// ---------- structured programs for basic-block counting ----------

type sitem struct {
	kind  string // "stmt", "if", "ifelse", "while"
	text  string
	cond  string
	then  []sitem
	other []sitem
}

func randStmt() sitem {
	v := edakit.Pick([]string{"a", "b", "c", "d", "s", "t"})
	u := edakit.Pick([]string{"a", "b", "c", "d", "i"})
	return sitem{kind: "stmt", text: edakit.Sprintf("%s = %s + %d;", v, u, edakit.Int(1, 9))}
}

func randCond() string {
	return edakit.Sprintf("%s %s %d", edakit.Pick([]string{"a", "b", "c", "i"}), edakit.Pick([]string{">", "<", "!="}), edakit.Int(0, 9))
}

// genSeq builds a statement sequence that starts and ends with a simple statement and keeps
// constructs apart with simple statements, so no basic block is ever empty.
func genSeq(constructs, depth int) []sitem {
	out := []sitem{randStmt()}
	if edakit.Int(0, 1) == 0 {
		out = append(out, randStmt())
	}
	for i := 0; i < constructs; i++ {
		kind := edakit.Pick([]string{"if", "ifelse", "while"})
		it := sitem{kind: kind, cond: randCond()}
		inner := 0
		if depth > 0 && edakit.Int(0, 2) == 0 {
			inner = 1
		}
		it.then = genSeq(inner, depth-1)
		if kind == "ifelse" {
			it.other = genSeq(0, depth-1)
		}
		out = append(out, it, randStmt())
	}
	return out
}

func renderSeq(items []sitem, pad string) []string {
	var out []string
	for _, it := range items {
		switch it.kind {
		case "stmt":
			out = append(out, pad+it.text)
		case "if":
			out = append(out, pad+"if ("+it.cond+") begin")
			out = append(out, renderSeq(it.then, pad+"  ")...)
			out = append(out, pad+"end")
		case "ifelse":
			out = append(out, pad+"if ("+it.cond+") begin")
			out = append(out, renderSeq(it.then, pad+"  ")...)
			out = append(out, pad+"end else begin")
			out = append(out, renderSeq(it.other, pad+"  ")...)
			out = append(out, pad+"end")
		case "while":
			out = append(out, pad+"while ("+it.cond+") begin")
			out = append(out, renderSeq(it.then, pad+"  ")...)
			out = append(out, pad+"end")
		}
	}
	return out
}

type cfgBuild struct {
	blocks int
	edges  int
	decis  int
}

func (g *cfgBuild) newBlock() int { g.blocks++; return g.blocks }

// walk threads a sequence through the CFG starting in block cur and returns the open block.
func (g *cfgBuild) walk(items []sitem, cur int) int {
	for _, it := range items {
		switch it.kind {
		case "stmt":
		case "if":
			g.decis++
			t := g.newBlock()
			g.edges++ // cur -> then
			tEnd := g.walk(it.then, t)
			j := g.newBlock()
			g.edges += 2 // tEnd -> join, cur -> join
			_ = tEnd
			cur = j
		case "ifelse":
			g.decis++
			t := g.newBlock()
			e := g.newBlock()
			g.edges += 2
			g.walk(it.then, t)
			g.walk(it.other, e)
			j := g.newBlock()
			g.edges += 2
			cur = j
		case "while":
			g.decis++
			g.newBlock() // header
			g.edges++    // cur -> header
			b := g.newBlock()
			g.edges++ // header -> body
			g.walk(it.then, b)
			g.edges++ // body end -> header
			x := g.newBlock()
			g.edges++ // header -> exit
			cur = x
		}
		_ = cur
	}
	return cur
}

// cfgCount asks for the number of basic blocks (KP 1), edges (KP 2) or the cyclomatic
// complexity (KP 3) of a structured procedural snippet.
func cfgCount(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 1, 3)
	prog := genSeq(n, 1)
	g := &cfgBuild{}
	g.walk(prog, g.newBlock())
	lines := renderSeq(prog, "  ")
	conv := "Use maximal basic blocks. An if condition is the last statement of the block that evaluates it; a while condition sits alone in its own header block, which the loop body jumps back to. There are no separate empty entry or exit nodes."
	var ask string
	var ans int
	switch k {
	case 1:
		ask = "How many basic blocks does its control-flow graph have?"
		ans = g.blocks
	case 2:
		ask = "How many edges does its control-flow graph have?"
		ans = g.edges
	default:
		ask = "What is the cyclomatic complexity E - N + 2 of its control-flow graph?"
		ans = g.edges - g.blocks + 2
	}
	q := "This code runs inside a procedural block." + code(append(append([]string{"begin"}, lines...), "end")...) + conv + " " + ask
	why := edakit.Sprintf("A new block starts at the target of every branch and at every join. An if-then adds a then block and a join block with 3 edges; an if-else adds two arms and a join with 4 edges; a while adds a header, a body and an exit block with 4 edges (into the header, header to body, back edge, header to exit). Here N = %d blocks and E = %d edges, so E - N + 2 = %d, which also equals the %d decisions plus one.",
		g.blocks, g.edges, g.edges-g.blocks+2, g.decis)
	return edakit.Numeric(q, edakit.IntStr(ans), why,
		map[string]string{"blocks": edakit.IntStr(g.blocks), "edges": edakit.IntStr(g.edges)})
}

// ---------- SSA renaming ----------

type ssaStmt struct {
	target string
	a, b   string
}

func (s ssaStmt) text() string { return s.target + " = " + s.a + " + " + s.b + ";" }

// ssaRenaming asks which version a straight-line read sees (KP 1), how many phi functions a
// diamond needs in minimal SSA (KP 2) or in pruned SSA (KP 3).
func ssaRenaming(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	vars := []string{"x", "y", "z", "w"}
	operand := func() string {
		if edakit.Int(0, 3) == 0 {
			return edakit.IntStr(edakit.Int(1, 9))
		}
		return edakit.Pick(vars)
	}
	if k == 1 {
		n := edakit.Scaled(ctx, 5, 7)
		var stmts []ssaStmt
		for i := 0; i < n; i++ {
			t := edakit.Pick(vars[:3])
			if i < 3 {
				t = "x"
			}
			stmts = append(stmts, ssaStmt{target: t, a: edakit.Pick(vars), b: operand()})
		}
		stmts = edakit.Shuffle(stmts)
		// Pick the last statement that reads some variable assigned earlier.
		for {
			idx := edakit.Int(n/2, n-1)
			s := stmts[idx]
			v := s.a
			defs := 0
			for _, p := range stmts[:idx] {
				if p.target == v {
					defs++
				}
			}
			if defs == 0 && edakit.Int(0, 3) > 0 {
				continue
			}
			var lines []string
			for i, p := range stmts {
				lines = append(lines, edakit.Sprintf("S%d: %s", i+1, p.text()))
			}
			q := "This straight-line code is converted to SSA. Each variable's incoming value is version 0 and every assignment to it creates the next version (1, 2, 3, ...)." +
				code(lines...) + edakit.Sprintf("Which version of %s does the first operand of S%d read? Answer with the version number.", v, idx+1)
			why := edakit.Sprintf("In straight-line code the reaching definition of a read is the most recent assignment above it. %s is assigned %d time(s) before S%d, so that read is %s_%d.", v, defs, idx+1, v, defs)
			return edakit.Numeric(q, edakit.IntStr(defs), why, map[string]string{"version": edakit.IntStr(defs)})
		}
	}
	for {
		then := map[string]bool{}
		els := map[string]bool{}
		var thenS, elseS []ssaStmt
		for i := 0; i < edakit.Int(1, 3); i++ {
			t := edakit.Pick(vars)
			then[t] = true
			thenS = append(thenS, ssaStmt{target: t, a: edakit.Pick(vars), b: operand()})
		}
		for i := 0; i < edakit.Int(1, 3); i++ {
			t := edakit.Pick(vars)
			els[t] = true
			elseS = append(elseS, ssaStmt{target: t, a: edakit.Pick(vars), b: operand()})
		}
		assigned := map[string]bool{}
		for v := range then {
			assigned[v] = true
		}
		for v := range els {
			assigned[v] = true
		}
		// Code after the join: reads and possibly redefinitions.
		var after []ssaStmt
		live := map[string]bool{}
		killed := map[string]bool{}
		for i := 0; i < edakit.Int(2, 3); i++ {
			s := ssaStmt{target: edakit.Pick(vars), a: edakit.Pick(vars), b: edakit.Pick(vars)}
			for _, r := range []string{s.a, s.b} {
				if !killed[r] {
					live[r] = true
				}
			}
			killed[s.target] = true
			after = append(after, s)
		}
		minimal := len(assigned)
		pruned := 0
		for v := range assigned {
			if live[v] {
				pruned++
			}
		}
		if k == 3 && (pruned == minimal || pruned == 0) && edakit.Int(0, 3) > 0 {
			continue
		}
		lines := []string{"x = in0; y = in1; z = in2; w = in3;", "if (sel) begin"}
		for _, s := range thenS {
			lines = append(lines, "  "+s.text())
		}
		lines = append(lines, "end else begin")
		for _, s := range elseS {
			lines = append(lines, "  "+s.text())
		}
		lines = append(lines, "end")
		for _, s := range after {
			lines = append(lines, s.text())
		}
		lines = append(lines, "// nothing after this point reads x, y, z or w")
		var q, why string
		var ans int
		if k == 2 {
			ans = minimal
			q = "This code is converted to minimal SSA (a phi wherever two different definitions meet, live or not)." + code(lines...) +
				"How many phi functions are placed at the join after the if/else?"
			why = edakit.Sprintf("Every variable assigned in at least one arm reaches the join with two different definitions, the arm's and the one from before the if, so it gets a phi: %s, which is %d.", strings.Join(sortedKeys(assigned), ", "), ans)
		} else {
			ans = pruned
			q = "This code is converted to pruned SSA, which places a phi only for a variable that is live at the join." + code(lines...) +
				"How many phi functions are placed at the join after the if/else?"
			liveAssigned := map[string]bool{}
			for v := range assigned {
				if live[v] {
					liveAssigned[v] = true
				}
			}
			names := "none"
			if len(liveAssigned) > 0 {
				names = strings.Join(sortedKeys(liveAssigned), ", ")
			}
			why = edakit.Sprintf("Minimal SSA would place %d phis, one per variable assigned in an arm (%s). Pruning keeps only those read after the join before being overwritten: %s, which is %d.", minimal, strings.Join(sortedKeys(assigned), ", "), names, ans)
		}
		return edakit.Numeric(q, edakit.IntStr(ans), why, map[string]string{"minimal": edakit.IntStr(minimal), "pruned": edakit.IntStr(pruned)})
	}
}

// ---------- random CFGs, dominators, frontiers ----------

type cfg struct {
	n     int
	succ  [][]int
	pred  [][]int
	edges [][2]int
}

// randCFG draws a CFG on nodes 1..n with entry 1. Every node i >= 2 gets a random earlier
// predecessor (so all are reachable), plus extra forward edges and optional back edges.
func randCFG(n int, backEdges int) *cfg {
	g := &cfg{n: n, succ: make([][]int, n+1), pred: make([][]int, n+1)}
	has := map[[2]int]bool{}
	add := func(a, b int) {
		if a == b || has[[2]int{a, b}] {
			return
		}
		has[[2]int{a, b}] = true
		g.succ[a] = append(g.succ[a], b)
		g.pred[b] = append(g.pred[b], a)
		g.edges = append(g.edges, [2]int{a, b})
	}
	for i := 2; i <= n; i++ {
		p := edakit.Int(max(1, i-2), i-1)
		add(p, i)
	}
	for e := 0; e < n/2; e++ {
		a := edakit.Int(1, n-2)
		b := edakit.Int(a+2, n)
		add(a, b)
	}
	for e := 0; e < backEdges; e++ {
		a := edakit.Int(3, n)
		b := edakit.Int(2, a-1)
		add(a, b)
	}
	sort.Slice(g.edges, func(i, j int) bool {
		if g.edges[i][0] != g.edges[j][0] {
			return g.edges[i][0] < g.edges[j][0]
		}
		return g.edges[i][1] < g.edges[j][1]
	})
	return g
}

func (g *cfg) edgeText() string {
	var parts []string
	for _, e := range g.edges {
		parts = append(parts, edakit.Sprintf("%d→%d", e[0], e[1]))
	}
	return strings.Join(parts, ", ")
}

// doms returns dom[v] as a set of nodes, by the iterative data-flow algorithm.
func (g *cfg) doms() []map[int]bool {
	dom := make([]map[int]bool, g.n+1)
	all := map[int]bool{}
	for v := 1; v <= g.n; v++ {
		all[v] = true
	}
	dom[1] = map[int]bool{1: true}
	for v := 2; v <= g.n; v++ {
		dom[v] = map[int]bool{}
		for x := range all {
			dom[v][x] = true
		}
	}
	for changed := true; changed; {
		changed = false
		for v := 2; v <= g.n; v++ {
			var nd map[int]bool
			for _, p := range g.pred[v] {
				if nd == nil {
					nd = map[int]bool{}
					for x := range dom[p] {
						nd[x] = true
					}
				} else {
					for x := range nd {
						if !dom[p][x] {
							delete(nd, x)
						}
					}
				}
			}
			nd[v] = true
			if len(nd) != len(dom[v]) {
				changed = true
			}
			dom[v] = nd
		}
	}
	return dom
}

func (g *cfg) idoms(dom []map[int]bool) []int {
	id := make([]int, g.n+1)
	for v := 2; v <= g.n; v++ {
		best, bestSize := 0, -1
		for d := range dom[v] {
			if d != v && len(dom[d]) > bestSize {
				best, bestSize = d, len(dom[d])
			}
		}
		id[v] = best
	}
	return id
}

func (g *cfg) frontiers(id []int) []map[int]bool {
	df := make([]map[int]bool, g.n+1)
	for v := range df {
		df[v] = map[int]bool{}
	}
	for b := 1; b <= g.n; b++ {
		if len(g.pred[b]) < 2 {
			continue
		}
		for _, p := range g.pred[b] {
			r := p
			for r != id[b] && r != 0 {
				df[r][b] = true
				r = id[r]
			}
		}
	}
	return df
}

func setList(s map[int]bool) []int {
	var out []int
	for x := range s {
		out = append(out, x)
	}
	sort.Ints(out)
	return out
}

func listStr(xs []int) string { return joinInts(xs, ",") }

// dominators asks for a dominator set (KP 1), an immediate dominator (KP 2) or a dominance
// frontier (KP 3) in a random CFG.
func dominators(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 5, 8)
	for {
		g := randCFG(n, edakit.Int(0, 1))
		dom := g.doms()
		id := g.idoms(dom)
		df := g.frontiers(id)
		v := edakit.Int(3, n)
		head := edakit.Sprintf("A control-flow graph has blocks 1 to %d with entry block 1 and edges %s.", n, g.edgeText())
		switch k {
		case 1:
			d := setList(dom[v])
			if len(d) < 2 {
				continue
			}
			q := head + edakit.Sprintf(" Which blocks dominate block %d? List them separated by commas, including %d itself.", v, v)
			why := edakit.Sprintf("Block d dominates v when every path from the entry to v passes through d. Solving Dom(v) = {v} ∪ ⋂ Dom(p) over the predecessors p gives Dom(%d) = {%s}.", v, listStr(d))
			return edakit.Numeric(q, listStr(d), why, map[string]string{"dom": listStr(d)})
		case 2:
			if len(g.pred[v]) < 2 && edakit.Int(0, 2) > 0 {
				continue
			}
			q := head + edakit.Sprintf(" What is the immediate dominator of block %d?", v)
			why := edakit.Sprintf("The immediate dominator is the strict dominator closest to the block, the one every other strict dominator also dominates. Dom(%d) = {%s}, and the closest strict dominator is %d.", v, listStr(setList(dom[v])), id[v])
			return edakit.Numeric(q, edakit.IntStr(id[v]), why, map[string]string{"idom": edakit.IntStr(id[v])})
		default:
			var cands []int
			for x := 2; x <= n; x++ {
				if len(df[x]) > 0 {
					cands = append(cands, x)
				}
			}
			if len(cands) == 0 {
				continue
			}
			v = edakit.Pick(cands)
			d := setList(df[v])
			q := head + edakit.Sprintf(" What is the dominance frontier of block %d? List its blocks separated by commas.", v)
			why := edakit.Sprintf("The dominance frontier of v is the set of blocks w where v dominates a predecessor of w but does not strictly dominate w: the first places where v's dominance ends. Walking up the dominator tree from each predecessor of every join gives DF(%d) = {%s}.", v, listStr(d))
			return edakit.Numeric(q, listStr(d), why, map[string]string{"df": listStr(d)})
		}
	}
}

func idf(df []map[int]bool, defs []int) map[int]bool {
	out := map[int]bool{}
	work := append([]int(nil), defs...)
	inWork := map[int]bool{}
	for _, d := range defs {
		inWork[d] = true
	}
	for len(work) > 0 {
		x := work[len(work)-1]
		work = work[:len(work)-1]
		for y := range df[x] {
			if !out[y] {
				out[y] = true
				if !inWork[y] {
					inWork[y] = true
					work = append(work, y)
				}
			}
		}
	}
	return out
}

// phiPlacement asks for the number of phi functions minimal SSA places for a variable, by
// the iterated dominance frontier of its definition blocks. KP 1 is acyclic, KP 2 has a loop,
// KP 3 sums over two variables.
func phiPlacement(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 5, 8)
	back := 0
	if k >= 2 {
		back = 1
	}
	for {
		g := randCFG(n, back)
		dom := g.doms()
		df := g.frontiers(g.idoms(dom))
		pickDefs := func() []int {
			m := edakit.Int(1, 2)
			d := edakit.Shuffle(seq(n - 1))[:m]
			var out []int
			for _, x := range d {
				out = append(out, x+2)
			}
			sort.Ints(out)
			return out
		}
		head := edakit.Sprintf("A control-flow graph has blocks 1 to %d with entry block 1 and edges %s.", n, g.edgeText())
		if k < 3 {
			defs := pickDefs()
			phis := setList(idf(df, append(defs, 1)))
			if len(phis) == 0 {
				continue
			}
			q := head + edakit.Sprintf(" Variable v has a value on entry and is assigned again in block(s) %s. How many phi functions for v does minimal SSA place?", listStr(defs))
			why := edakit.Sprintf("Minimal SSA places phis at the iterated dominance frontier of the definition blocks: start from DF of each defining block, and since a phi is itself a definition, add the frontiers of the new phi blocks until nothing changes. That gives phis in block(s) %s, %d in all.", listStr(phis), len(phis))
			return edakit.Numeric(q, edakit.IntStr(len(phis)), why, map[string]string{"phi_blocks": listStr(phis)})
		}
		du, dv := pickDefs(), pickDefs()
		pu := setList(idf(df, append(du, 1)))
		pv := setList(idf(df, append(dv, 1)))
		total := len(pu) + len(pv)
		if total == 0 {
			continue
		}
		q := head + edakit.Sprintf(" Variables u and v have values on entry. u is assigned again in block(s) %s and v in block(s) %s. How many phi functions does minimal SSA place in total for u and v?", listStr(du), listStr(dv))
		bl := func(xs []int) string {
			if len(xs) == 0 {
				return "none"
			}
			return listStr(xs)
		}
		why := edakit.Sprintf("Each variable gets phis at the iterated dominance frontier of its own definition blocks. For u that is %s; for v it is %s. Total %d.", bl(pu), bl(pv), total)
		return edakit.Numeric(q, edakit.IntStr(total), why, map[string]string{"u_phis": edakit.IntStr(len(pu)), "v_phis": edakit.IntStr(len(pv))})
	}
}

// ---------- data-flow analysis ----------

// dataflow poses reaching definitions (KP 1) or liveness, live-in (KP 2) and live-out (KP 3),
// on a random CFG with one back edge.
func dataflow(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 4, 6)
	vars := []string{"a", "b", "c", "d"}
	for {
		g := randCFG(n, 1)
		head := edakit.Sprintf("A control-flow graph has blocks 1 to %d with entry block 1 and edges %s.", n, g.edgeText())
		if k == 1 {
			type def struct {
				id    int
				v     string
				block int
			}
			var defs []def
			byBlock := make([][]def, n+1)
			id := 0
			for b := 1; b <= n; b++ {
				for i := 0; i < edakit.Int(0, 2); i++ {
					id++
					d := def{id: id, v: edakit.Pick(vars[:3]), block: b}
					defs = append(defs, d)
					byBlock[b] = append(byBlock[b], d)
				}
			}
			if id < 4 {
				continue
			}
			gen := make([]map[int]bool, n+1)
			kill := make([]map[int]bool, n+1)
			for b := 1; b <= n; b++ {
				gen[b], kill[b] = map[int]bool{}, map[int]bool{}
				last := map[string]int{}
				for _, d := range byBlock[b] {
					last[d.v] = d.id
				}
				for _, id := range last {
					gen[b][id] = true
				}
				for _, d := range defs {
					if _, ok := last[d.v]; ok && !gen[b][d.id] {
						kill[b][d.id] = true
					}
				}
			}
			in := make([]map[int]bool, n+1)
			out := make([]map[int]bool, n+1)
			for b := range in {
				in[b], out[b] = map[int]bool{}, map[int]bool{}
			}
			for changed := true; changed; {
				changed = false
				for b := 1; b <= n; b++ {
					ni := map[int]bool{}
					for _, p := range g.pred[b] {
						for x := range out[p] {
							ni[x] = true
						}
					}
					no := map[int]bool{}
					for x := range gen[b] {
						no[x] = true
					}
					for x := range ni {
						if !kill[b][x] {
							no[x] = true
						}
					}
					if len(ni) != len(in[b]) || len(no) != len(out[b]) {
						changed = true
					}
					in[b], out[b] = ni, no
				}
			}
			var cands []int
			for b := 2; b <= n; b++ {
				if len(in[b]) >= 2 {
					cands = append(cands, b)
				}
			}
			if len(cands) == 0 {
				continue
			}
			t := edakit.Pick(cands)
			var dl []string
			for b := 1; b <= n; b++ {
				var parts []string
				for _, d := range byBlock[b] {
					parts = append(parts, edakit.Sprintf("d%d: %s = ...", d.id, d.v))
				}
				if len(parts) == 0 {
					parts = []string{"no definitions"}
				}
				dl = append(dl, edakit.Sprintf("block %d: %s", b, strings.Join(parts, "; ")))
			}
			ans := listStr(setList(in[t]))
			q := head + " The definitions, in order within each block, are:" + code(dl...) +
				edakit.Sprintf("Which definitions reach the entry of block %d? List their numbers separated by commas.", t)
			why := edakit.Sprintf("Reaching definitions is a forward may-analysis: IN[b] is the union of OUT over the predecessors, and OUT[b] = GEN[b] ∪ (IN[b] - KILL[b]), where GEN is the last definition of each variable in b and KILL every other definition of those variables. Iterating to the fixed point gives IN[%d] = {%s}.", t, ans)
			return edakit.Numeric(q, ans, why, map[string]string{"in": ans})
		}
		use := make([]map[string]bool, n+1)
		defn := make([]map[string]bool, n+1)
		var dl []string
		for b := 1; b <= n; b++ {
			use[b], defn[b] = map[string]bool{}, map[string]bool{}
			for _, v := range vars {
				switch edakit.Int(0, 4) {
				case 0:
					use[b][v] = true
				case 1:
					defn[b][v] = true
				}
			}
			show := func(s map[string]bool) string {
				if len(s) == 0 {
					return "{}"
				}
				return "{" + strings.Join(sortedKeys(s), ", ") + "}"
			}
			dl = append(dl, edakit.Sprintf("block %d: use %s, def %s", b, show(use[b]), show(defn[b])))
		}
		in := make([]map[string]bool, n+1)
		out := make([]map[string]bool, n+1)
		for b := range in {
			in[b], out[b] = map[string]bool{}, map[string]bool{}
		}
		for changed := true; changed; {
			changed = false
			for b := n; b >= 1; b-- {
				no := map[string]bool{}
				for _, s := range g.succ[b] {
					for x := range in[s] {
						no[x] = true
					}
				}
				ni := map[string]bool{}
				for x := range use[b] {
					ni[x] = true
				}
				for x := range no {
					if !defn[b][x] {
						ni[x] = true
					}
				}
				if len(ni) != len(in[b]) || len(no) != len(out[b]) {
					changed = true
				}
				in[b], out[b] = ni, no
			}
		}
		t := edakit.Int(1, n)
		var set map[string]bool
		var what string
		if k == 2 {
			set, what = in[t], "at the entry of"
		} else {
			set, what = out[t], "at the exit of"
		}
		if len(set) == 0 {
			continue
		}
		q := head + " Each block's upward-exposed uses (read before any write in the block) and definitions are:" + code(dl...) +
			edakit.Sprintf("Nothing is live after a block with no successors. How many variables are live %s block %d?", what, t)
		why := edakit.Sprintf("Liveness is a backward may-analysis: OUT[b] is the union of IN over the successors and IN[b] = USE[b] ∪ (OUT[b] - DEF[b]). Iterating to the fixed point gives {%s} live %s block %d, which is %d variables.", strings.Join(sortedKeys(set), ", "), what, t, len(set))
		return edakit.Numeric(q, edakit.IntStr(len(set)), why, map[string]string{"live": strings.Join(sortedKeys(set), " ")})
	}
}
