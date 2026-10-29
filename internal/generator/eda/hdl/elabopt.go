package hdl

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// parameters asks for a port width after parameter overrides: named (KP 1), passed down a
// two-level hierarchy (KP 2) or positional against the declaration order (KP 3).
func parameters(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	d0, w0 := edakit.Pick([]int{4, 8, 16}), edakit.Pick([]int{8, 16, 32})
	wFirst := k == 3 && edakit.Int(0, 1) == 0
	var header string
	if wFirst {
		header = edakit.Sprintf("module buffer #(parameter W = %d, parameter DEPTH = %d) (", w0, d0)
	} else {
		header = edakit.Sprintf("module buffer #(parameter DEPTH = %d, parameter W = %d) (", d0, w0)
	}
	mod := []string{header, "  output [PW-1:0] bus", ");", "  localparam PW = W + 2*DEPTH;", "  // ...", "endmodule"}
	var inst []string
	var depth, w int
	var why string
	switch k {
	case 1:
		if edakit.Int(0, 1) == 0 {
			w = edakit.Pick([]int{4, 12, 20, 24})
			depth = d0
			inst = []string{edakit.Sprintf("buffer #(.W(%d)) u0 (.bus(b0));", w)}
			why = "Only W is overridden, so DEPTH keeps its default"
		} else {
			depth = edakit.Pick([]int{2, 6, 10, 32})
			w = w0
			inst = []string{edakit.Sprintf("buffer #(.DEPTH(%d)) u0 (.bus(b0));", depth)}
			why = "Only DEPTH is overridden, so W keeps its default"
		}
	case 2:
		nv := edakit.Int(2, 6)
		ma, mb := edakit.Int(2, 3), edakit.Int(1, 4)
		depth, w = ma*nv, nv+mb
		inst = []string{
			edakit.Sprintf("module core #(parameter N = 2) ();"),
			edakit.Sprintf("  buffer #(.DEPTH(N*%d), .W(N+%d)) u_buf (.bus(b0));", ma, mb),
			"endmodule",
			"",
			edakit.Sprintf("// in top: core #(.N(%d)) u_core ();", nv),
		}
		why = edakit.Sprintf("Elaboration resolves top down: u_core gets N = %d, so the overrides N*%d and N+%d are evaluated with that value", nv, ma, mb)
	default:
		p1, p2 := edakit.Pick([]int{3, 5, 6, 12}), edakit.Pick([]int{2, 7, 9, 10})
		inst = []string{edakit.Sprintf("buffer #(%d, %d) u0 (.bus(b0));", p1, p2)}
		if wFirst {
			w, depth = p1, p2
		} else {
			depth, w = p1, p2
		}
		why = "Positional overrides follow the order in which the parameters are declared"
	}
	pw := w + 2*depth
	q := "Given this module:" + code(mod...) + "and this instantiation:" + code(inst...) +
		"How many bits wide is the bus port of the buffer instance?"
	why += edakit.Sprintf(": DEPTH = %d and W = %d, so PW = %d + 2 x %d = %d.", depth, w, w, depth, pw)
	return edakit.Numeric(q, edakit.IntStr(pw), why, map[string]string{"depth": edakit.IntStr(depth), "w": edakit.IntStr(w)})
}

// generateLoops counts the instances a generate construct elaborates: one loop (KP 1), a
// nested loop with a condition (KP 2) or a halving reduction tree (KP 3).
func generateLoops(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1:
		start := edakit.Int(0, 3)
		end := edakit.Int(start+5, start+20)
		step := edakit.Pick([]int{1, 2, 3, 4})
		le := edakit.Int(0, 1) == 0
		cmp := "<"
		if le {
			cmp = "<="
		}
		count := 0
		for i := start; (le && i <= end) || (!le && i < end); i += step {
			count++
		}
		lines := []string{"genvar i;", "generate", edakit.Sprintf("  for (i = %d; i %s %d; i = i + %d) begin : g", start, cmp, end, step), "    slice u (.in(d[i]), .out(q[i]));", "  end", "endgenerate"}
		q := "Consider this generate loop:" + code(lines...) + "How many instances of slice does it elaborate?"
		why := edakit.Sprintf("The loop runs at elaboration time with i = %d, %d, ... while i %s %d, which is %d iterations, each creating g[i].u.", start, start+step, cmp, end, count)
		return edakit.Numeric(q, edakit.IntStr(count), why, map[string]string{"count": edakit.IntStr(count)})
	case 2:
		n, m := edakit.Int(3, 6), edakit.Int(3, 6)
		conds := []struct {
			text string
			f    func(i, j int) bool
		}{
			{"j <= i", func(i, j int) bool { return j <= i }},
			{"j < i", func(i, j int) bool { return j < i }},
			{"(i + j) % 2 == 0", func(i, j int) bool { return (i+j)%2 == 0 }},
			{"i != j", func(i, j int) bool { return i != j }},
		}
		c := edakit.Pick(conds)
		count := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if c.f(i, j) {
					count++
				}
			}
		}
		lines := []string{"genvar i, j;", "generate",
			edakit.Sprintf("  for (i = 0; i < %d; i = i + 1) begin : row", n),
			edakit.Sprintf("    for (j = 0; j < %d; j = j + 1) begin : col", m),
			"      if (" + c.text + ") begin : pe",
			"        mac u (.a(a[i]), .b(b[j]));",
			"      end", "    end", "  end", "endgenerate"}
		q := "Consider this generate block:" + code(lines...) + "How many mac instances does it elaborate?"
		why := edakit.Sprintf("The loops enumerate %d x %d = %d (i, j) pairs at elaboration time, and the generate-if keeps the pairs with %s: %d of them.", n, m, n*m, c.text, count)
		return edakit.Numeric(q, edakit.IntStr(count), why, map[string]string{"pairs": edakit.IntStr(n * m)})
	}
	n := edakit.Pick([]int{8, 12, 16, 20, 24, 32})
	l := edakit.Int(2, 4)
	count := 0
	var parts []string
	for lv := 0; lv < l; lv++ {
		c := n >> (lv + 1)
		count += c
		parts = append(parts, edakit.IntStr(c))
	}
	lines := []string{"module tree #(parameter N = 8, parameter L = 3) (...);", "  genvar l, k;", "  generate",
		"    for (l = 0; l < L; l = l + 1) begin : lvl",
		"      for (k = 0; k < (N >> (l + 1)); k = k + 1) begin : node",
		"        add2 u (...);", "      end", "    end", "  endgenerate", "endmodule"}
	q := edakit.Sprintf("The module below is instantiated with N = %d and L = %d.", n, l) + code(lines...) + "How many add2 instances does elaboration create?"
	why := edakit.Sprintf("Level l creates N >> (l + 1) adders, so the levels contribute %s, a total of %d.", strings.Join(parts, " + "), count)
	return edakit.Numeric(q, edakit.IntStr(count), why, map[string]string{"count": edakit.IntStr(count)})
}

// flatten counts instances (KP 1) and primitive cells (KP 2) after flattening a hierarchy,
// or the distinct module definitions after uniquifying parameterized instances (KP 3).
func flatten(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k < 3 {
		n1, n2 := edakit.Int(1, 4), edakit.Int(1, 3)
		n3, n4 := edakit.Int(1, 4), edakit.Int(0, 3)
		cT, cA, cB, cC := edakit.Int(2, 20), edakit.Int(5, 40), edakit.Int(5, 40), edakit.Int(10, 60)
		lines := []string{
			edakit.Sprintf("top:   %d primitive cells, %s of alu, %s of ctrl", cT, nInst(n1), nInst(n2)),
			edakit.Sprintf("alu:   %d primitive cells, %s of regfile", cA, nInst(n3)),
			edakit.Sprintf("ctrl:  %d primitive cells, %s of regfile", cB, nInst(n4)),
			edakit.Sprintf("regfile: %d primitive cells, no instances", cC),
		}
		insts := n1*n3 + n2*n4
		cells := cT + n1*(cA+n3*cC) + n2*(cB+n4*cC)
		if k == 1 {
			q := "An elaborated design has this hierarchy:" + code(lines...) + "After flattening, how many copies of regfile's logic does the netlist contain?"
			why := edakit.Sprintf("Flattening copies a module once per path from the top: %d alu x %d + %d ctrl x %d = %d.", n1, n3, n2, n4, insts)
			return edakit.Numeric(q, edakit.IntStr(insts), why, map[string]string{"regfile_copies": edakit.IntStr(insts)})
		}
		q := "An elaborated design has this hierarchy:" + code(lines...) + "After flattening, how many primitive cells does the single flat module contain?"
		why := edakit.Sprintf("Each instance contributes its own cells plus its children's: %d + %d x (%d + %d x %d) + %d x (%d + %d x %d) = %d.", cT, n1, cA, n3, cC, n2, cB, n4, cC, cells)
		return edakit.Numeric(q, edakit.IntStr(cells), why, map[string]string{"regfile_copies": edakit.IntStr(insts)})
	}
	var lines []string
	seen := map[string]bool{"top": true}
	nInst := edakit.Scaled(ctx, 4, 6)
	for i := 0; i < nInst; i++ {
		switch edakit.Int(0, 3) {
		case 0:
			lines = append(lines, edakit.Sprintf("fifo u%d (...);", i))
			seen["fifo"] = true
		case 1:
			lines = append(lines, edakit.Sprintf("ctrl u%d (...);", i))
			seen["ctrl"] = true
		default:
			w := edakit.Pick([]int{8, 16, 32})
			if edakit.Int(0, 1) == 0 {
				d := edakit.Pick([]int{2, 8})
				lines = append(lines, edakit.Sprintf("fifo #(.W(%d), .DEPTH(%d)) u%d (...);", w, d, i))
				seen[edakit.Sprintf("fifo W=%d DEPTH=%d", w, d)] = true
			} else {
				lines = append(lines, edakit.Sprintf("fifo #(.W(%d)) u%d (...);", w, i))
				seen[edakit.Sprintf("fifo W=%d DEPTH=4", w)] = true
			}
		}
	}
	keys := make([]string, 0, len(seen))
	for kk := range seen {
		keys = append(keys, kk)
	}
	sort.Strings(keys)
	q := "Module fifo has parameters W (default 4) and DEPTH (default 4); ctrl has none. Module top contains only these instances:" + code(lines...) +
		"The elaborator creates one specialized module per distinct parameter assignment. How many module definitions does the elaborated design contain, counting top?"
	why := edakit.Sprintf("Instances that agree on every parameter value share one specialized module, and every distinct assignment needs its own. The distinct modules are: %s, which is %d.", strings.Join(keys, "; "), len(keys))
	return edakit.Numeric(q, edakit.IntStr(len(keys)), why, map[string]string{"modules": edakit.IntStr(len(keys))})
}

// ---------- optimization ----------

type tac struct {
	dst  string
	op   string
	a, b string
}

func isConstTok(s string) bool { return s != "" && (s[0] == '-' || (s[0] >= '0' && s[0] <= '9')) }

func atoi(s string) int {
	neg := false
	if s[0] == '-' {
		neg, s = true, s[1:]
	}
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	if neg {
		return -n
	}
	return n
}

func apply(op string, x, y int) int {
	switch op {
	case "+":
		return x + y
	case "-":
		return x - y
	case "*":
		return x * y
	case "&":
		return x & y
	}
	return 0
}

// constantPropagation asks for a folded value (KP 1), the number of temporaries that fold
// (KP 2, with the annihilator identities) or a value that needs branch folding (KP 3).
func constantPropagation(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k == 3 {
		return sccp(ctx)
	}
	n := edakit.Scaled(ctx, 4, 6)
	for {
		var code3 []tac
		consts := map[string]int{}
		isC := map[string]bool{}
		pool := []string{"x"}
		for i := 1; i <= n; i++ {
			dst := edakit.Sprintf("t%d", i)
			op := edakit.Pick([]string{"+", "-", "*"})
			a := edakit.Pick(pool)
			b := edakit.IntStr(edakit.Int(1, 6))
			if edakit.Int(0, 2) == 0 {
				b = edakit.Pick(pool)
			}
			if i == 1 || edakit.Int(0, 3) == 0 {
				a = edakit.IntStr(edakit.Int(1, 9))
			}
			if a == b && !isConstTok(a) {
				b = edakit.IntStr(edakit.Int(1, 6))
			}
			if k == 2 && edakit.Int(0, 4) == 0 {
				op = edakit.Pick([]string{"*", "&"})
				a = edakit.Pick(pool)
				b = "0"
			}
			c := tac{dst: dst, op: op, a: a, b: b}
			val := func(s string) (int, bool) {
				if isConstTok(s) {
					return atoi(s), true
				}
				v, ok := consts[s]
				return v, ok && isC[s]
			}
			va, oka := val(a)
			vb, okb := val(b)
			switch {
			case oka && okb:
				consts[dst], isC[dst] = apply(op, va, vb), true
			case (op == "*" || op == "&") && ((okb && vb == 0) || (oka && va == 0)):
				consts[dst], isC[dst] = 0, true
			}
			code3 = append(code3, c)
			pool = append(pool, dst)
		}
		var lines []string
		for _, c := range code3 {
			lines = append(lines, edakit.Sprintf("%s = %s %s %s;", c.dst, c.a, c.op, c.b))
		}
		folded := 0
		var foldedNames []string
		for _, c := range code3 {
			if isC[c.dst] {
				folded++
				foldedNames = append(foldedNames, edakit.Sprintf("%s = %d", c.dst, consts[c.dst]))
			}
		}
		if k == 1 {
			var cands []string
			for _, c := range code3[1:] {
				if isC[c.dst] && !isConstTok(c.a) || isC[c.dst] && !isConstTok(c.b) {
					cands = append(cands, c.dst)
				}
			}
			if len(cands) == 0 {
				continue
			}
			t := cands[len(cands)-1]
			if consts[t] > 5000 || consts[t] < -5000 {
				continue
			}
			q := "In this SSA code, x is an input whose value is unknown at compile time." + code(lines...) +
				edakit.Sprintf("Constant propagation folds %s to a constant. What is it?", t)
			why := "Propagate forward: a temporary whose operands are all known constants is evaluated at compile time and its value substituted into later uses. The folded values are " + strings.Join(foldedNames, ", ") + "."
			return edakit.Numeric(q, edakit.IntStr(consts[t]), why, map[string]string{"value": edakit.IntStr(consts[t])})
		}
		if folded == 0 || folded == n {
			continue
		}
		q := "In this SSA code, x is an input whose value is unknown at compile time." + code(lines...) +
			"Constant propagation folds operations whose operands are all constants and also applies the identities v * 0 = 0 and v & 0 = 0. How many of the temporaries become compile-time constants?"
		why := "A temporary is constant when both operands are constants, or when it multiplies or ANDs anything with a constant 0. Constant: " + strings.Join(foldedNames, ", ") + edakit.Sprintf(". That is %d of %d.", folded, n)
		return edakit.Numeric(q, edakit.IntStr(folded), why, map[string]string{"folded": edakit.IntStr(folded)})
	}
}

func sccp(ctx generator.GeneratorContext) generator.Problem {
	c1, c2 := edakit.Int(1, 9), edakit.Int(1, 9)
	kv := c1 + c2
	constInThen := edakit.Int(0, 1) == 0
	var thr int
	var cmp string
	if constInThen {
		thr = edakit.Int(1, kv-1)
		cmp = ">"
	} else {
		thr = edakit.Int(1, kv-1)
		cmp = "<"
	}
	m := edakit.Int(2, 4)
	c3 := edakit.Int(1, 9)
	constArm := edakit.Sprintf("y = k * %d;", m)
	inputArm := "y = x + 1;"
	thenArm, elseArm := constArm, inputArm
	if !constInThen {
		thenArm, elseArm = inputArm, constArm
	}
	lines := []string{
		edakit.Sprintf("k = %d + %d;", c1, c2),
		edakit.Sprintf("if (k %s %d)", cmp, thr),
		"  " + thenArm,
		"else",
		"  " + elseArm,
		edakit.Sprintf("z = y + %d;", c3),
	}
	z := kv*m + c3
	q := "Here x is an input unknown at compile time." + code(lines...) +
		"Sparse conditional constant propagation tracks which branches can execute. What constant does it find for z?"
	why := edakit.Sprintf("k folds to %d, so the condition k %s %d is decided at compile time and only the arm assigning y = k * %d is reachable. The phi for y then sees one executable input, y = %d, and z = %d + %d = %d. Plain constant propagation, which meets both arms at the phi, would call y and z unknown.", kv, cmp, thr, m, kv*m, kv*m, c3, z)
	return edakit.Numeric(q, edakit.IntStr(z), why, map[string]string{"k": edakit.IntStr(kv)})
}

// cse counts the operations local value numbering removes: exact repeats (KP 1), repeats up
// to commutativity (KP 2) and code where reassignment kills availability (KP 3).
func cse(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 5, 8)
	inputs := []string{"a", "b", "c"}
	for {
		vn := map[string]int{}
		next := 0
		for _, v := range inputs {
			next++
			vn[v] = next
		}
		table := map[string]int{}
		var lines []string
		var tmps []string
		redundant := 0
		var redNames []string
		type pair struct{ op, a, b string }
		var history []pair
		for i := 1; i <= n; i++ {
			// KP 3: occasionally overwrite an input with a copy before continuing.
			if k == 3 && i > 2 && edakit.Int(0, 2) == 0 && len(tmps) > 0 {
				v := edakit.Pick(inputs)
				src := edakit.Pick(tmps)
				vn[v] = vn[src]
				lines = append(lines, edakit.Sprintf("%s = %s;", v, src))
			}
			dst := edakit.Sprintf("t%d", i)
			var p pair
			if len(history) > 0 && edakit.Int(0, 2) == 0 {
				p = edakit.Pick(history)
				if k >= 2 && (p.op == "+" || p.op == "*" || p.op == "&") && edakit.Int(0, 1) == 0 {
					p.a, p.b = p.b, p.a
				}
			} else {
				pool := append(append([]string{}, inputs...), tmps...)
				p = pair{op: edakit.Pick([]string{"+", "*", "-", "&"}), a: edakit.Pick(pool), b: edakit.Pick(pool)}
				for p.b == p.a {
					p.b = edakit.Pick(pool)
				}
			}
			history = append(history, p)
			x, y := vn[p.a], vn[p.b]
			if p.op != "-" && x > y {
				x, y = y, x
			}
			key := edakit.Sprintf("%s %d %d", p.op, x, y)
			if v, ok := table[key]; ok {
				redundant++
				redNames = append(redNames, dst)
				vn[dst] = v
			} else {
				next++
				vn[dst] = next
				table[key] = next
			}
			lines = append(lines, edakit.Sprintf("%s = %s %s %s;", dst, p.a, p.op, p.b))
			tmps = append(tmps, dst)
		}
		if redundant == 0 {
			continue
		}
		rule := "Two operations match only if they have the same operator and the same operand values in the same order."
		if k >= 2 {
			rule = "+, * and & are commutative, so operand order does not matter for them; - is not commutative."
		}
		q := "Local value numbering runs over this basic block. " + rule + " A temporary that is replaced by a copy carries the same value number as the original." +
			code(lines...) + "How many of the t operations are redundant and replaced by a copy of an earlier result?"
		why := "Value numbering gives every distinct value a number and keys each operation by its operator and operand value numbers, so a reassigned variable gets a new number and stops matching. The redundant operations are " +
			strings.Join(redNames, ", ") + edakit.Sprintf(", which is %d.", redundant)
		return edakit.Numeric(q, edakit.IntStr(redundant), why, map[string]string{"redundant": edakit.IntStr(redundant)})
	}
}

// dce counts the statements a mark-and-sweep dead-code eliminator removes, including (KP 3)
// a dead register loop that use-count-based removal would keep.
func dce(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 5, 8)
	for {
		type st struct {
			dst  string
			srcs []string
			text string
		}
		var stmts []st
		pool := []string{"a", "b", "c"}
		for i := 1; i <= n; i++ {
			dst := edakit.Sprintf("t%d", i)
			x, y := edakit.Pick(pool), edakit.Pick(pool)
			op := edakit.Pick([]string{"+", "&", "^", "|"})
			stmts = append(stmts, st{dst: dst, srcs: []string{x, y}, text: edakit.Sprintf("%s = %s %s %s;", dst, x, op, y)})
			pool = append(pool, dst)
		}
		if k == 3 {
			x := edakit.Pick(pool)
			stmts = append(stmts,
				st{dst: "r1", srcs: []string{"r2", "a"}, text: "r1 <= r2 + a;   // register"},
				st{dst: "r2", srcs: []string{"r1", x}, text: edakit.Sprintf("r2 <= r1 ^ %s;   // register", x)})
		}
		// Outputs read one or two temporaries.
		nOut := edakit.Int(1, 2)
		var outs []string
		for i := 0; i < nOut; i++ {
			outs = append(outs, edakit.Sprintf("t%d", edakit.Int(n/2+1, n)))
		}
		live := map[string]bool{}
		work := append([]string(nil), outs...)
		def := map[string]st{}
		for _, s := range stmts {
			def[s.dst] = s
		}
		for len(work) > 0 {
			v := work[len(work)-1]
			work = work[:len(work)-1]
			if live[v] {
				continue
			}
			live[v] = true
			if s, ok := def[v]; ok {
				work = append(work, s.srcs...)
			}
		}
		removed := 0
		var dead []string
		for _, s := range stmts {
			if !live[s.dst] {
				removed++
				dead = append(dead, s.dst)
			}
		}
		if removed == 0 || (k == 1 && removed > 3) {
			continue
		}
		var lines []string
		for _, s := range stmts {
			lines = append(lines, s.text)
		}
		outText := "out = " + outs[0]
		if len(outs) == 2 {
			outText = "out = {" + outs[0] + ", " + outs[1] + "}"
		}
		lines = append(lines, outText+";   // the only output")
		q := "A dead-code eliminator marks everything the outputs depend on, transitively, and sweeps the rest." + code(lines...) +
			"How many statements, other than the output assignment, does it remove?"
		why := "Start from the output and walk backwards through operands, marking each definition reached; anything unmarked cannot affect an output. Removed: " + strings.Join(dead, ", ") + edakit.Sprintf(", which is %d.", removed)
		if k == 3 {
			why += " The two registers read each other, so each still has a use; a pass that only deletes definitions with zero uses would keep them, but marking from the outputs never reaches them."
		}
		return edakit.Numeric(q, edakit.IntStr(removed), why, map[string]string{"removed": edakit.IntStr(removed)})
	}
}

func nInst(n int) string {
	if n == 1 {
		return "1 instance"
	}
	return edakit.Sprintf("%d instances", n)
}
