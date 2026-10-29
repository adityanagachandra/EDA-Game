package synth

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// naf returns the canonical signed-digit (non-adjacent form) digits of k, least significant first.
func naf(k int) []int {
	var d []int
	for k > 0 {
		if k%2 == 1 {
			x := 2 - k%4
			d = append(d, x)
			k -= x
		} else {
			d = append(d, 0)
		}
		k /= 2
	}
	return d
}

func nafTex(d []int) string {
	var parts []string
	for i := len(d) - 1; i >= 0; i-- {
		parts = append(parts, edakit.IntStr(d[i]))
	}
	return "(" + strings.Join(parts, " ") + ")"
}

// csd: shift-add adders for binary and CSD constant multiplication.
func csd(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 7, 1000)
	for popcount(uint64(k)) < 3 {
		k = edakit.Scaled(ctx, 7, 1000)
	}
	bin := popcount(uint64(k))
	d := naf(k)
	nz := 0
	for _, x := range d {
		if x != 0 {
			nz++
		}
	}
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(edakit.Sprintf("Synthesis implements $y = %d x$ as a sum of shifted copies of x, one per 1 bit of the binary constant. How many two-input adders are needed?", k),
			edakit.IntStr(bin-1), edakit.Sprintf("%d in binary is %b, with %d ones, so %d shifted terms and %d adders.", k, k, bin, bin, bin-1),
			map[string]string{"ones": edakit.IntStr(bin)})
	case 2:
		return edakit.Numeric(edakit.Sprintf("How many nonzero digits does the canonical signed-digit (non-adjacent form) representation of %d have?", k),
			edakit.IntStr(nz), edakit.Sprintf("Scanning from the least significant bit, replace each run of ones by a +1 above the run and a -1 at its bottom. %d becomes the digit string %s, most significant first, with %d nonzero digits.", k, nafTex(d), nz),
			map[string]string{"csd": nafTex(d)})
	}
	return edakit.Numeric(edakit.Sprintf("Compare implementing $y = %d x$ with plain binary shift-and-add against canonical signed digits, where each nonzero digit after the first costs one adder or subtractor. How many adders and subtractors does CSD save?", k),
		edakit.IntStr(bin-nz), edakit.Sprintf("Binary %b has %d ones (%d adders); CSD %s has %d nonzero digits (%d adders or subtractors). Saving %d.", k, bin, bin-1, nafTex(d), nz, nz-1, bin-nz),
		map[string]string{"saving": edakit.IntStr(bin - nz)})
}

// carrySave: 3:2 reduction levels, number of word-level CSAs, total delay with a final adder.
func carrySave(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 4, 40)
	levels := 0
	for d := 2; d < k; d = d * 3 / 2 {
		levels++
	}
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(edakit.Sprintf("A multi-operand adder sums %d operands with a tree of 3:2 carry-save adders until two rows remain. What is the minimum number of carry-save levels?", k),
			edakit.IntStr(levels), edakit.Sprintf("Each level turns every 3 rows into 2, so the reachable row counts are 2, 3, 4, 6, 9, 13, 19, 28, 42 (each about 1.5 times the last). %d operands need %d levels.", k, levels),
			map[string]string{"levels": edakit.IntStr(levels)})
	case 2:
		return edakit.Numeric(edakit.Sprintf("How many word-level 3:2 carry-save adders are needed to reduce %d operands to two rows, before the final carry-propagate adder?", k),
			edakit.IntStr(k-2), edakit.Sprintf("Each 3:2 adder takes three rows and returns two, removing exactly one row. Going from %d rows to 2 removes %d.", k, k-2),
			map[string]string{"csas": edakit.IntStr(k - 2)})
	}
	n := edakit.Pick([]int{8, 16, 24, 32})
	t := levels + n
	return edakit.Numeric(edakit.Sprintf("A %d-operand, %d-bit adder uses a minimum-depth 3:2 carry-save tree followed by a %d-bit ripple-carry adder. Counting one full-adder delay per carry-save level and one per ripple stage, and ignoring word growth, what is the total delay in full-adder delays?", k, n, n),
		edakit.IntStr(t), edakit.Sprintf("The tree needs %d levels, each one full adder deep regardless of width, then the ripple adder adds %d. Total %d.", levels, n, t),
		map[string]string{"levels": edakit.IntStr(levels), "delay": edakit.IntStr(t)})
}

// congruence: e-classes under hashconsing, unions and congruence closure.
func congruence(ctx generator.GeneratorContext) generator.Problem {
	type enode struct {
		op   byte
		l, r int
	}
	which := kp(ctx)
	for {
		var nodes []enode
		ids := map[enode]int{}
		add := func(e enode) int {
			if id, ok := ids[e]; ok {
				return id
			}
			ids[e] = len(nodes)
			nodes = append(nodes, e)
			return ids[e]
		}
		leafNames := []string{"a", "b", "c"}
		var termID []int
		var names []string
		for i := range leafNames {
			termID = append(termID, add(enode{op: byte('a' + i), l: -1, r: -1}))
			names = append(names, leafNames[i])
		}
		nt := edakit.Scaled(ctx, 4, 6)
		var defs []string
		twin := edakit.Int(0, 1) == 0
		twinOp := edakit.Pick([]byte{'+', '*'})
		for t := 0; t < nt; t++ {
			x, y := edakit.Int(0, len(termID)-1), edakit.Int(0, len(termID)-1)
			op := edakit.Pick([]byte{'+', '*'})
			if twin && t < 2 {
				// Twin terms over a and b invite congruence once a and b are merged.
				x, y, op = t, 2, twinOp
			}
			termID = append(termID, add(enode{op: op, l: termID[x], r: termID[y]}))
			nm := edakit.Sprintf("t_%d", t+1)
			sym := map[byte]string{'+': " + ", '*': " \\cdot "}[op]
			defs = append(defs, "$"+nm+" = "+names[x]+sym+names[y]+"$")
			names = append(names, nm)
		}
		initial := len(nodes)
		parent := make([]int, len(nodes))
		for i := range parent {
			parent[i] = i
		}
		var find func(x int) int
		find = func(x int) int {
			for parent[x] != x {
				parent[x] = parent[parent[x]]
				x = parent[x]
			}
			return x
		}
		union := func(a, b int) bool {
			ra, rb := find(a), find(b)
			if ra == rb {
				return false
			}
			parent[ra] = rb
			return true
		}
		roots := func() int {
			m := map[int]bool{}
			for i := range nodes {
				m[find(i)] = true
			}
			return len(m)
		}
		nu := 1
		if ctx.Difficulty > 0.6 {
			nu = 2
		}
		var us []string
		for i := 0; i < nu; i++ {
			x, y := edakit.Int(0, 2), edakit.Int(0, len(names)-1)
			if twin && i == 0 {
				x, y = 0, 1
			}
			if x == y {
				y = (y + 1) % len(names)
			}
			union(termID[x], termID[y])
			us = append(us, "$"+names[x]+" \\equiv "+names[y]+"$")
		}
		afterUnion := roots()
		for changed := true; changed; {
			changed = false
			seen := map[enode]int{}
			for i, e := range nodes {
				k := e
				if e.l >= 0 {
					k = enode{e.op, find(e.l), find(e.r)}
				}
				if j, ok := seen[k]; ok {
					if union(i, j) {
						changed = true
					}
				} else {
					seen[k] = i
				}
			}
		}
		closed := roots()
		stem := "An e-graph is built from the leaves a, b, c and the terms " + strings.Join(defs, ", ") + ", with hashconsing (an identical operator applied to identical e-classes reuses the existing e-node)."
		switch which {
		case 1:
			return edakit.Numeric(stem+" Before any equality is asserted, how many e-classes does the e-graph hold (leaves included)?",
				edakit.IntStr(initial), edakit.Sprintf("Each distinct e-node starts in its own e-class; duplicate definitions collapse into one. %d e-classes.", initial),
				map[string]string{"classes": edakit.IntStr(initial)})
		case 2:
			return edakit.Numeric(stem+" Now assert "+strings.Join(us, " and ")+" and restore congruence (if $x \\equiv y$ then $f(x) \\equiv f(y)$) until fixpoint. How many e-classes remain?",
				edakit.IntStr(closed), edakit.Sprintf("The asserted unions leave %d classes; congruence closure merges further, ending with %d.", afterUnion, closed),
				map[string]string{"classes": edakit.IntStr(closed)})
		}
		if afterUnion == closed && edakit.Int(0, 2) > 0 {
			continue
		}
		return edakit.Numeric(stem+" Assert "+strings.Join(us, " and ")+". How many additional e-class merges does congruence closure (the e-graph rebuild) perform beyond the asserted unions themselves?",
			edakit.IntStr(afterUnion-closed), edakit.Sprintf("The unions alone leave %d e-classes; repairing congruence brings it to %d, so %d merges came from congruence.", afterUnion, closed, afterUnion-closed),
			map[string]string{"merges": edakit.IntStr(afterUnion - closed)})
	}
}

type egNode struct {
	op   string
	kids []int
}

// extractCost computes the minimum tree cost of every e-class by relaxation (handles cycles).
func extractCost(classes [][]egNode, cost map[string]int) []int {
	inf := 1 << 30
	c := make([]int, len(classes))
	for i := range c {
		c[i] = inf
	}
	for iter := 0; iter < 64; iter++ {
		changed := false
		for i, cl := range classes {
			for _, n := range cl {
				t := cost[n.op]
				ok := true
				for _, k := range n.kids {
					if c[k] == inf {
						ok = false
						break
					}
					t += c[k]
				}
				if ok && t < c[i] {
					c[i] = t
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	return c
}

func egDesc(classes [][]egNode) string {
	var parts []string
	for i, cl := range classes {
		var ns []string
		for _, n := range cl {
			switch len(n.kids) {
			case 0:
				ns = append(ns, n.op)
			case 2:
				sym := map[string]string{"*": " \\cdot ", "+": " + ", "<<": " \\ll "}[n.op]
				ns = append(ns, edakit.Sprintf("C_{%d}%sC_{%d}", n.kids[0], sym, n.kids[1]))
			}
		}
		parts = append(parts, edakit.Sprintf("$C_{%d} = \\{%s\\}$", i, strings.Join(ns, ", ")))
	}
	return strings.Join(parts, ", ")
}

// egExtraction: cheapest term in an e-graph under an operator cost model.
func egExtraction(ctx generator.GeneratorContext) generator.Problem {
	which := kp(ctx)
	cost := map[string]int{"*": edakit.Int(3, 8), "<<": edakit.Int(1, 2), "+": edakit.Int(1, 3)}
	leaf := func(s string) []egNode { return []egNode{{op: s}} }
	bin := func(op string, a, b int) egNode { return egNode{op: op, kids: []int{a, b}} }
	// a*3 + b*4 after saturation.
	classes := [][]egNode{
		leaf("a"), leaf("b"), leaf("3"), leaf("4"), leaf("2"), leaf("1"),
		{bin("*", 0, 4), bin("<<", 0, 5), bin("+", 0, 0)}, // C6 = a*2
		{bin("*", 0, 2), bin("+", 6, 0)},                  // C7 = a*3
		{bin("*", 1, 3), bin("<<", 1, 4)},                 // C8 = b*4
		{bin("+", 7, 8), bin("+", 8, 7)},                  // C9 = root
	}
	if which == 3 {
		classes[0] = append(classes[0], bin("*", 0, 5))  // a = a*1
		classes[7] = append(classes[7], bin("*", 7, 5))  // a*3 = (a*3)*1
		classes[9] = append(classes[9], bin("+", 9, 10)) // root = root + 0
		classes = append(classes, leaf("0"))
	}
	for _, k := range []string{"a", "b", "0", "1", "2", "3", "4"} {
		cost[k] = 0
	}
	c := extractCost(classes, cost)
	stem := edakit.Sprintf("After equality saturation of $a \\cdot 3 + b \\cdot 4$ the e-graph is %s. Costs: multiply %d, shift %d, add %d, leaves and constants 0. The cost of a term is the sum of its operator costs (tree cost).",
		egDesc(classes), cost["*"], cost["<<"], cost["+"])
	switch which {
	case 1:
		cl := edakit.Pick([]int{6, 7, 8})
		return edakit.Numeric(stem+edakit.Sprintf(" What is the cost of the cheapest term in e-class $C_%d$?", cl),
			edakit.IntStr(c[cl]), edakit.Sprintf("Bottom up: $C_6$ costs %d, $C_7$ costs %d, $C_8$ costs %d.", c[6], c[7], c[8]),
			map[string]string{"cost": edakit.IntStr(c[cl])})
	case 2:
		return edakit.Numeric(stem+" What is the cost of the cheapest term extracted from the root e-class $C_9$?",
			edakit.IntStr(c[9]), edakit.Sprintf("Pick the cheapest e-node in every class from the leaves up: $C_6$ = %d, $C_7$ = %d, $C_8$ = %d, so the root costs %d + %d + %d = %d.", c[6], c[7], c[8], cost["+"], c[7], c[8], c[9]),
			map[string]string{"cost": edakit.IntStr(c[9])})
	}
	return edakit.Numeric(stem+" Several e-classes contain themselves through identities such as $x \\cdot 1 = x$, so the e-graph has cycles. What is the cost of the cheapest finite term in the root e-class $C_9$?",
		edakit.IntStr(c[9]), edakit.Sprintf("A cyclic e-node can never be part of the cheapest term, because it costs more than the class it loops back to. Relaxation to a fixpoint gives root cost %d.", c[9]),
		map[string]string{"cost": edakit.IntStr(c[9])})
}

// cec: miter counts and miter size.
func cec(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 4)
	vars := seq(n)
	which := kp(ctx)
	full := uint64(1)<<(1<<n) - 1
	pair := func() (*expr, *expr) {
		var f *expr
		for {
			f = randExprSpread(vars, edakit.Int(2, 4))
			t := table(f, n)
			if t != 0 && t != full {
				break
			}
		}
		var g *expr
		if edakit.Int(0, 3) == 0 {
			g = deMorgan(f)
		} else {
			g = mutate(f, vars)
		}
		return f, g
	}
	switch which {
	case 1:
		f, g := pair()
		diff := popcount((table(f, n) ^ table(g, n)) & full)
		return edakit.Numeric(edakit.Sprintf("An optimizer rewrote $F = %s$ into $G = %s$ over inputs %s. Build the miter $F \\oplus G$. How many input vectors set the miter output to 1? (0 means the rewrite is correct.)",
			f.tex(varNames), g.tex(varNames), strings.Join(varNames[:n], ", ")),
			edakit.IntStr(diff), edakit.Sprintf("The miter is 1 exactly where F and G differ. Comparing truth tables gives %d such vectors.", diff),
			map[string]string{"counterexamples": edakit.IntStr(diff)})
	case 2:
		f1, g1 := pair()
		f2, g2 := pair()
		diff := popcount(((table(f1, n) ^ table(g1, n)) | (table(f2, n) ^ table(g2, n))) & full)
		return edakit.Numeric(edakit.Sprintf("A two-output circuit computes $F_1 = %s$ and $F_2 = %s$; the optimized version computes $G_1 = %s$ and $G_2 = %s$. The miter ORs the XORs of corresponding outputs. On how many input vectors is the miter 1?",
			f1.tex(varNames), f2.tex(varNames), g1.tex(varNames), g2.tex(varNames)),
			edakit.IntStr(diff), edakit.Sprintf("A vector is a counterexample if any output pair differs. %d vectors are.", diff),
			map[string]string{"counterexamples": edakit.IntStr(diff)})
	}
	f1, g1 := pair()
	f2, g2 := pair()
	size := f1.ands() + g1.ands() + f2.ands() + g2.ands() + 3*2 + 1
	return edakit.Numeric(edakit.Sprintf("Build the AIG miter of a two-output circuit, $F_1 = %s$, $F_2 = %s$, against its rewrite, $G_1 = %s$, $G_2 = %s$, with no structural sharing. Each two-input AND or OR is one AND node, an XOR is three, and the final OR of the two output XORs is one. How many AND nodes does the miter have?",
		f1.tex(varNames), f2.tex(varNames), g1.tex(varNames), g2.tex(varNames)),
		edakit.IntStr(size), edakit.Sprintf("The four cones cost %d, %d, %d and %d AND nodes; two output XORs cost 6 and the OR costs 1. Total %d.", f1.ands(), f2.ands(), g1.ands(), g2.ands(), size),
		map[string]string{"miter_ands": edakit.IntStr(size)})
}
