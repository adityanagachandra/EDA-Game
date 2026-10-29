package hdl

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// bnStmt is one procedural assignment target <= or = rhs, where rhs is src1 [+ src2 | + k].
type bnStmt struct {
	target   string
	blocking bool
	src1     string
	src2     string // "" when absent
	k        int    // constant addend when src2 is ""
}

func (s bnStmt) rhs() string {
	switch {
	case s.src2 != "":
		return s.src1 + " + " + s.src2
	case s.k != 0:
		return edakit.Sprintf("%s + %d", s.src1, s.k)
	}
	return s.src1
}

func (s bnStmt) text() string {
	op := "<="
	if s.blocking {
		op = "="
	}
	return edakit.Sprintf("%s %s %s;", s.target, op, s.rhs())
}

// runBN simulates edges rising clock edges of one always @(posedge clk) block.
// mode 0 uses the declared operators, 1 treats all as nonblocking, 2 all as blocking.
func runBN(stmts []bnStmt, init map[string]int, edges, mode int) map[string]int {
	cur := map[string]int{}
	for k, v := range init {
		cur[k] = v
	}
	for e := 0; e < edges; e++ {
		pending := map[string]int{}
		for _, s := range stmts {
			v := cur[s.src1]
			if s.src2 != "" {
				v += cur[s.src2]
			} else {
				v += s.k
			}
			v %= 256
			blocking := s.blocking
			if mode == 1 {
				blocking = false
			} else if mode == 2 {
				blocking = true
			}
			if blocking {
				cur[s.target] = v
			} else {
				pending[s.target] = v
			}
		}
		for k, v := range pending {
			cur[k] = v
		}
	}
	return cur
}

// blockingNonblocking asks for a register value after one or more clock edges. KP 1 uses
// only nonblocking assignments, KP 2 only blocking, KP 3 a mix.
func blockingNonblocking(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 3, 4)
	names := varNames[:n]
	for tries := 0; ; tries++ {
		init := map[string]int{}
		for _, v := range names {
			init[v] = edakit.Int(1, 9)
		}
		var stmts []bnStmt
		for _, t := range edakit.Shuffle(names) {
			s := bnStmt{target: t}
			others := []string{}
			for _, v := range names {
				if v != t {
					others = append(others, v)
				}
			}
			s.src1 = edakit.Pick(others)
			switch edakit.Int(0, 2) {
			case 0:
				s.src2 = edakit.Pick(names)
			case 1:
				s.k = edakit.Int(1, 5)
			}
			switch k {
			case 2:
				s.blocking = true
			case 3:
				s.blocking = edakit.Int(0, 1) == 0
			}
			stmts = append(stmts, s)
		}
		if k == 3 {
			nb, bl := 0, 0
			for _, s := range stmts {
				if s.blocking {
					bl++
				} else {
					nb++
				}
			}
			if nb == 0 || bl == 0 {
				continue
			}
		}
		edges := 1
		if ctx.Difficulty > 0.5 && edakit.Int(0, 1) == 0 {
			edges = 2
		}
		ask := edakit.Pick(names)
		got := runBN(stmts, init, edges, 0)[ask]
		alt := runBN(stmts, init, edges, 1)[ask]
		if k == 1 {
			alt = runBN(stmts, init, edges, 2)[ask]
		}
		if got == alt && tries < 40 {
			continue
		}
		var lines []string
		lines = append(lines, "always @(posedge clk) begin")
		for _, s := range stmts {
			lines = append(lines, "  "+s.text())
		}
		lines = append(lines, "end")
		var inits []string
		for _, v := range names {
			inits = append(inits, edakit.Sprintf("%s = %d", v, init[v]))
		}
		edgeWord := "one rising clock edge"
		if edges == 2 {
			edgeWord = "two rising clock edges"
		}
		q := "All variables are 8-bit regs, starting at " + strings.Join(inits, ", ") + ". Nothing else writes them." +
			code(lines...) + edakit.Sprintf("What is the value of %s after %s? Answer in decimal.", ask, edgeWord)
		why := edakit.Sprintf("Execute the statements in order. A blocking assignment (=) updates its variable immediately, so later statements read the new value; a nonblocking assignment (<=) samples its right-hand side now and updates only after the whole block has run. Doing that for %s gives %s = %d. Treating every statement the other way would give %d.",
			edgeWord, ask, got, alt)
		return edakit.Numeric(q, edakit.IntStr(got), why,
			map[string]string{"edges": edakit.IntStr(edges), "naive": edakit.IntStr(alt)})
	}
}

func widthDecl(name string, w int) string {
	if w == 1 {
		return "reg " + name + ";"
	}
	return edakit.Sprintf("reg [%d:0] %s;", w-1, name)
}

// flopInference counts the flip-flop bits a clocked block infers before optimization.
// KP 1 counts nonblocking targets, KP 2 adds blocking temporaries that are written before
// they are read (no storage), KP 3 adds a temporary read before it is written (storage).
func flopInference(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	widths := []int{1, 4, 8, 16, 3, 12}
	nq := edakit.Scaled(ctx, 2, 4)
	var decls, body []string
	total := 0
	var parts []string
	qs := make([]string, nq)
	qw := make([]int, nq)
	for i := range qs {
		qs[i] = edakit.Sprintf("q%d", i+1)
		qw[i] = edakit.Pick(widths)
		decls = append(decls, widthDecl(qs[i], qw[i]))
	}
	type temp struct {
		name       string
		w          int
		readFirst  bool
		consumer   int
		expression string
	}
	var temps []temp
	switch k {
	case 2:
		for i := 0; i < edakit.Int(1, 2); i++ {
			temps = append(temps, temp{name: edakit.Sprintf("t%d", i+1), w: edakit.Pick(widths[1:]), consumer: i % nq})
		}
	case 3:
		temps = append(temps, temp{name: "t1", w: edakit.Pick(widths[1:]), readFirst: true, consumer: 0})
		if edakit.Int(0, 1) == 0 {
			temps = append(temps, temp{name: "t2", w: edakit.Pick(widths[1:]), consumer: 1 % nq})
		}
	}
	exprs := []string{"a + b", "a ^ b", "a & b", "b - a"}
	for i := range temps {
		temps[i].expression = edakit.Pick(exprs)
		decls = append(decls, widthDecl(temps[i].name, temps[i].w))
	}
	enable := edakit.Int(0, nq-1)
	for i := range qs {
		rhs := edakit.Pick([]string{"a", "b", qs[i] + " + 1", "a | b"})
		var consumer *temp
		for j := range temps {
			if temps[j].consumer == i && consumer == nil {
				consumer = &temps[j]
			}
		}
		if consumer != nil {
			rhs = consumer.name
			if consumer.readFirst {
				body = append(body, edakit.Sprintf("%s <= %s;", qs[i], rhs))
				body = append(body, edakit.Sprintf("%s = %s;", consumer.name, consumer.expression))
			} else {
				body = append(body, edakit.Sprintf("%s = %s;", consumer.name, consumer.expression))
				body = append(body, edakit.Sprintf("%s <= %s;", qs[i], rhs))
			}
		} else if i == enable {
			body = append(body, edakit.Sprintf("if (en) %s <= %s;", qs[i], rhs))
		} else {
			body = append(body, edakit.Sprintf("%s <= %s;", qs[i], rhs))
		}
		total += qw[i]
		parts = append(parts, edakit.Sprintf("%s: %d", qs[i], qw[i]))
	}
	for _, t := range temps {
		if t.readFirst {
			total += t.w
			parts = append(parts, edakit.Sprintf("%s: %d (read before written, so it holds last cycle's value)", t.name, t.w))
		} else {
			parts = append(parts, edakit.Sprintf("%s: 0 (written before read)", t.name))
		}
	}
	lines := append([]string{}, decls...)
	lines = append(lines, "wire [3:0] y;", "assign y = a[3:0] & b[3:0];", "always @(posedge clk) begin")
	lines = append(lines, indent(2, body)...)
	lines = append(lines, "end")
	note := ""
	if len(temps) > 0 {
		names := []string{}
		for _, t := range temps {
			names = append(names, t.name)
		}
		note = " The variable(s) " + strings.Join(names, ", ") + " are read nowhere outside this block."
	}
	q := "Inputs a and b are 16 bits wide." + note + code(lines...) +
		"Before any optimization, how many flip-flop bits does synthesis infer from this code?"
	why := "A variable assigned in a clocked block needs storage when its value must survive to the next cycle: every nonblocking target does, and a blocking temporary does only if some read sees the previous cycle's value. Continuous assigns infer no flops. Bits: " +
		strings.Join(parts, "; ") + edakit.Sprintf(". Total %d.", total)
	return edakit.Numeric(q, edakit.IntStr(total), why, map[string]string{"flop_bits": edakit.IntStr(total)})
}

// latchInference counts latch bits inferred from a combinational block with incomplete
// assignment. KP 1 is an if without else, KP 2 an if/else, KP 3 a case statement.
func latchInference(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 2, 4)
	ys := make([]string, n)
	ws := make([]int, n)
	var decls []string
	for i := range ys {
		ys[i] = edakit.Sprintf("y%d", i+1)
		ws[i] = edakit.Pick([]int{1, 2, 4, 8})
		decls = append(decls, widthDecl(ys[i], ws[i]))
	}
	rhs := func() string { return edakit.Pick([]string{"a", "b", "a & b", "a | b", "~a", "a ^ b"}) }
	subset := func() []bool {
		s := make([]bool, n)
		for i := range s {
			s[i] = edakit.Int(0, 2) > 0
		}
		return s
	}
	deflt := make([]bool, n)
	for i := range deflt {
		deflt[i] = edakit.Int(0, 4) == 0
	}
	var body []string
	for i := range ys {
		if deflt[i] {
			body = append(body, edakit.Sprintf("%s = 0;", ys[i]))
		}
	}
	latched := make([]bool, n)
	assignedSomewhere := make([]bool, n)
	for i := range deflt {
		assignedSomewhere[i] = deflt[i]
	}
	emit := func(s []bool, pad string) []string {
		var out []string
		for i := range s {
			if s[i] {
				out = append(out, pad+edakit.Sprintf("%s = %s;", ys[i], rhs()))
				assignedSomewhere[i] = true
			}
		}
		if len(out) == 0 {
			i := edakit.Int(0, n-1)
			s[i] = true
			assignedSomewhere[i] = true
			out = append(out, pad+edakit.Sprintf("%s = %s;", ys[i], rhs()))
		}
		return out
	}
	switch k {
	case 1:
		then := subset()
		body = append(body, "if (s) begin")
		body = append(body, emit(then, "  ")...)
		body = append(body, "end")
		for i := range ys {
			latched[i] = !deflt[i]
		}
	case 2:
		then, els := subset(), subset()
		body = append(body, "if (s) begin")
		body = append(body, emit(then, "  ")...)
		body = append(body, "end else begin")
		body = append(body, emit(els, "  ")...)
		body = append(body, "end")
		for i := range ys {
			latched[i] = !deflt[i] && !(then[i] && els[i])
		}
	default:
		items := edakit.Pick([]int{3, 4})
		hasDefault := items == 3 && edakit.Int(0, 1) == 0
		body = append(body, "case (sel)  // sel is 2 bits")
		all := make([]bool, n)
		for i := range all {
			all[i] = true
		}
		for it := 0; it < items; it++ {
			s := subset()
			if edakit.Int(0, 1) == 0 {
				for i := range s {
					s[i] = true
				}
			}
			lines := emit(s, "    ")
			body = append(body, edakit.Sprintf("  2'd%d: begin", it))
			body = append(body, lines...)
			body = append(body, "  end")
			for i := range s {
				all[i] = all[i] && s[i]
			}
		}
		var dflt []bool
		if hasDefault {
			dflt = subset()
			body = append(body, "  default: begin")
			body = append(body, emit(dflt, "    ")...)
			body = append(body, "  end")
		}
		body = append(body, "endcase")
		for i := range ys {
			full := all[i] && (items == 4 || (hasDefault && dflt[i]))
			latched[i] = !deflt[i] && !full
		}
	}
	total := 0
	var parts []string
	for i := range ys {
		if !assignedSomewhere[i] {
			latched[i] = false
		}
		if latched[i] {
			total += ws[i]
			parts = append(parts, edakit.Sprintf("%s (%d bits)", ys[i], ws[i]))
		}
	}
	if total == 0 && edakit.Int(0, 3) > 0 {
		return latchInference(ctx)
	}
	lines := append([]string{}, decls...)
	lines = append(lines, "always @(*) begin")
	lines = append(lines, indent(2, body)...)
	lines = append(lines, "end")
	q := "Inputs a and b are 8 bits wide (truncated on assignment) and s is 1 bit." + code(lines...) +
		"How many latch bits does synthesis infer from this block?"
	list := "none"
	if len(parts) > 0 {
		list = strings.Join(parts, ", ")
	}
	why := "A variable in a combinational block that is not assigned on every path must hold its old value on the missing paths, which is a latch. A default assignment at the top of the block, or assignment in every branch of a complete if/else or case, removes it. Latched: " +
		list + edakit.Sprintf(". Total %d bits.", total)
	return edakit.Numeric(q, edakit.IntStr(total), why, map[string]string{"latch_bits": edakit.IntStr(total)})
}

// caseSemantics asks which item of a case, casez or casex statement executes.
func caseSemantics(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	kind := []string{"", "casez", "casez", "casex"}[k]
	if k == 2 && edakit.Int(0, 2) == 0 {
		kind = "case"
	}
	for {
		sel := make([]byte, 4)
		for i := range sel {
			sel[i] = edakit.Pick([]byte{'0', '1'})
		}
		switch {
		case k == 2:
			sel[edakit.Int(0, 3)] = 'z'
		case k == 3:
			sel[edakit.Int(0, 3)] = 'x'
		}
		dc := byte('?')
		if kind == "casex" {
			dc = 'x'
		}
		nItems := 4
		items := make([][]byte, nItems)
		for i := range items {
			p := make([]byte, 4)
			for j := range p {
				switch edakit.Int(0, 4) {
				case 0, 1:
					p[j] = '0'
				case 2, 3:
					p[j] = '1'
				default:
					p[j] = dc
				}
			}
			if kind == "case" && i == 1 {
				// A literal z in a plain case item matches only a z in the selector.
				copy(p, sel)
				p[edakit.Int(0, 3)] = dc
			}
			items[i] = p
		}
		match := func(p []byte, kindNow string) bool {
			for j := range p {
				pb, sb := p[j], sel[j]
				if pb == '?' {
					pb = 'z'
				}
				switch kindNow {
				case "casez":
					if pb == 'z' || sb == 'z' {
						continue
					}
				case "casex":
					if pb == 'x' || pb == 'z' || sb == 'x' || sb == 'z' {
						continue
					}
				}
				if pb != sb {
					return false
				}
			}
			return true
		}
		first := func(kindNow string) int {
			for i, p := range items {
				if match(p, kindNow) {
					return i
				}
			}
			return -1
		}
		got := first(kind)
		// Want the answer to differ from the reading that treats the selector's unknown
		// bit as a mismatch, or from plain case when kind is casez.
		naive := first("case")
		if kind == "case" {
			naive = first("casez")
		}
		if got == naive && edakit.Int(0, 3) > 0 {
			continue
		}
		label := func(i int) string {
			if i < 0 {
				return "y = 0 (the default item)"
			}
			return edakit.Sprintf("y = %d", i+1)
		}
		var lines []string
		lines = append(lines, kind+" (sel)")
		for i, p := range items {
			lines = append(lines, edakit.Sprintf("  4'b%s: y = %d;", string(p), i+1))
		}
		lines = append(lines, "  default: y = 0;", "endcase")
		q := edakit.Sprintf("A simulator evaluates this statement with sel = 4'b%s.", string(sel)) + code(lines...) + "Which assignment executes?"
		var ds []string
		for i := -1; i < nItems; i++ {
			if i != got {
				ds = append(ds, label(i))
			}
		}
		rule := map[string]string{
			"case":  "Plain case compares like ===: every bit, x and z included, must match exactly, and ? is just another way to write z.",
			"casez": "casez treats z (or ?) as a wildcard in both the item and the selector.",
			"casex": "casex treats both x and z as wildcards in the item and in the selector, so an x on sel matches anything in that lane.",
		}[kind]
		return edakit.Choice(q, label(got), rule+" Items are tried in order and the first match wins, which here is "+label(got)+".", ds...)
	}
}

// sized expression tree for the sizing generator.
type sizeNode struct {
	text  string
	width int
	kind  string // "leaf", "arith", "cmp", "concat", "repl", "reduce", "shift", "cond"
}

// expressionSizing asks for bit widths under the IEEE 1364 sizing rules (KP 1 and 2) or the
// value of an expression whose evaluation width depends on its context (KP 3).
func expressionSizing(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k == 3 {
		return sizingContext(ctx)
	}
	opWidths := []int{3, 4, 5, 6, 8, 12, 16}
	leaves := map[string]int{}
	var order []string
	for _, v := range []string{"a", "b", "c", "d"} {
		leaves[v] = edakit.Pick(opWidths)
		order = append(order, v)
	}
	leaf := func() sizeNode {
		v := edakit.Pick(order)
		return sizeNode{text: v, width: leaves[v], kind: "leaf"}
	}
	wrap := func(n sizeNode) string {
		if n.kind == "leaf" || n.kind == "concat" || n.kind == "repl" {
			return n.text
		}
		return "(" + n.text + ")"
	}
	var build func(depth int) sizeNode
	build = func(depth int) sizeNode {
		if depth == 0 {
			return leaf()
		}
		var kinds []string
		if k == 1 {
			kinds = []string{"arith", "arith", "concat", "repl"}
		} else {
			kinds = []string{"arith", "cmp", "concat", "reduce", "shift", "cond", "repl"}
		}
		switch edakit.Pick(kinds) {
		case "arith":
			l, r := build(depth-1), build(depth-1)
			op := edakit.Pick([]string{"+", "-", "*", "&", "|", "^"})
			return sizeNode{text: wrap(l) + " " + op + " " + wrap(r), width: max(l.width, r.width), kind: "arith"}
		case "cmp":
			l, r := build(depth-1), build(depth-1)
			op := edakit.Pick([]string{"<", "==", "!=", ">="})
			return sizeNode{text: wrap(l) + " " + op + " " + wrap(r), width: 1, kind: "cmp"}
		case "concat":
			l, r := build(depth-1), build(depth-1)
			return sizeNode{text: "{" + l.text + ", " + r.text + "}", width: l.width + r.width, kind: "concat"}
		case "repl":
			c := build(depth - 1)
			n := edakit.Int(2, 4)
			return sizeNode{text: edakit.Sprintf("{%d{%s}}", n, wrap(c)), width: n * c.width, kind: "repl"}
		case "reduce":
			c := build(depth - 1)
			op := edakit.Pick([]string{"&", "|", "^"})
			return sizeNode{text: op + wrap(c), width: 1, kind: "reduce"}
		case "shift":
			l := build(depth - 1)
			r := leaf()
			op := edakit.Pick([]string{"<<", ">>"})
			return sizeNode{text: wrap(l) + " " + op + " " + r.text, width: l.width, kind: "shift"}
		default:
			s, l, r := leaf(), build(depth-1), build(depth-1)
			return sizeNode{text: s.text + " ? " + wrap(l) + " : " + wrap(r), width: max(l.width, r.width), kind: "cond"}
		}
	}
	depth := 2
	if ctx.Difficulty > 0.6 && k == 2 {
		depth = 3
	}
	var e sizeNode
	for {
		e = build(depth)
		if e.kind != "leaf" && len(e.text) < 90 {
			break
		}
	}
	var decl []string
	for _, v := range order {
		decl = append(decl, edakit.Sprintf("%s is %d bits", v, leaves[v]))
	}
	q := "All operands are unsigned: " + strings.Join(decl, ", ") + "." + code(e.text) +
		"What is the self-determined bit width of this expression?"
	why := edakit.Sprintf("Arithmetic and bitwise operators take the larger operand width, comparisons and reductions give 1 bit, a shift keeps its left operand's width, a conditional takes the larger of its two arms, a concatenation sums its parts and a replication multiplies. The expression is %d bits.", e.width)
	return edakit.Numeric(q, edakit.IntStr(e.width), why, map[string]string{"width": edakit.IntStr(e.width)})
}

func sizingContext(ctx generator.GeneratorContext) generator.Problem {
	w := edakit.Pick([]int{4, 8})
	top := 1 << w
	a := edakit.Int(top/2, top-1)
	b := edakit.Int(top/2, top-1)
	ly := edakit.Pick([]int{w, w + 1})
	concat := edakit.Int(0, 2) == 0
	evalW := max(ly, w)
	if concat {
		evalW = w
	}
	sum := (a + b) % (1 << evalW)
	y := (sum >> 1) % (1 << ly)
	expr := "(a + b) >> 1"
	if concat {
		expr = "{a + b} >> 1"
	}
	q := edakit.Sprintf("a and b are unsigned %d-bit values, a = %d and b = %d, and y is declared [%d:0].", w, a, b, ly-1) +
		code("assign y = "+expr+";") + "What value does y hold? Answer in decimal."
	why := edakit.Sprintf("The operands of + and >> are context-determined, so they are evaluated at the width of the largest operand or of the assignment target, %d bits here", max(ly, w))
	if concat {
		why = edakit.Sprintf("Inside a concatenation the sum is self-determined, so it is computed at %d bits whatever the width of y", w)
	}
	why += edakit.Sprintf(". a + b = %d becomes %d at %d bits; shifting right by one gives %d, and y keeps %d bits: %d.", a+b, sum, evalW, sum>>1, ly, y)
	return edakit.Numeric(q, edakit.IntStr(y), why, map[string]string{"eval_width": edakit.IntStr(evalW)})
}

// signedness asks for the 8-bit result of adding two 4-bit operands whose signedness decides
// between sign and zero extension.
func signedness(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	av := edakit.Int(8, 15) // negative as a signed nibble
	cv := edakit.Int(1, 15)
	sx := func(v int) int {
		if v >= 8 {
			return v - 16
		}
		return v
	}
	mod := func(v int) int { return ((v % 256) + 256) % 256 }
	bits := func(v int) string { return edakit.Sprintf("4'b%04b", v) }
	var decl []string
	var expr string
	var signedExpr bool
	var reason string
	switch k {
	case 1:
		decl = []string{"logic signed [3:0] a = " + bits(av) + ";", "logic signed [3:0] c = " + bits(cv) + ";"}
		expr = "a + c"
		signedExpr = true
		reason = "Both operands are signed, so the expression is signed and each is sign-extended to the 8-bit context."
	case 2:
		decl = []string{"logic signed [3:0] a = " + bits(av) + ";", "logic [3:0] c = " + bits(cv) + ";"}
		expr = "a + c"
		reason = "One operand is unsigned, so the whole expression is unsigned and a is zero-extended even though it was declared signed."
	default:
		decl = []string{"logic signed [3:0] a = " + bits(av) + ";", "logic signed [3:0] c = " + bits(cv) + ";"}
		switch edakit.Int(0, 2) {
		case 0:
			expr = "a[3:0] + c"
			reason = "A part-select is unsigned even when it covers a whole signed vector, so the expression is unsigned and both operands are zero-extended."
		case 1:
			expr = "a + 1'b1"
			cv = 1
			reason = "The sized literal 1'b1 is unsigned, which makes the whole expression unsigned, so a is zero-extended."
		default:
			expr = "a + 1"
			cv = 1
			signedExpr = true
			reason = "An unsized decimal literal is a signed 32-bit integer, so the expression stays signed and a is sign-extended before the add."
		}
	}
	var y, alt int
	if signedExpr {
		y = mod(sx(av) + sx(cv))
		alt = mod(av + cv)
	} else {
		y = mod(av + cv)
		alt = mod(sx(av) + sx(cv))
		if strings.Contains(expr, "1'b1") {
			alt = mod(sx(av) + 1)
		}
	}
	lines := append(decl, "logic [7:0] y;", "assign y = "+expr+";")
	q := "In SystemVerilog:" + code(lines...) + "What unsigned value does y hold? Answer in decimal."
	why := reason + edakit.Sprintf(" That gives y = %d. The other extension would give %d.", y, alt)
	return edakit.Numeric(q, edakit.IntStr(y), why, map[string]string{"y": edakit.IntStr(y)})
}

// sortedKeys returns the keys of a set in order.
func sortedKeys(m map[string]bool) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
