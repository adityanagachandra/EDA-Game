package hdl

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func isWord(t string) bool {
	c := t[0]
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// renderTokens joins tokens the way a person might type them: words are always separated,
// two symbols are always separated (so the lexer cannot fuse them into a longer operator),
// and a word next to a symbol is sometimes written without a space.
func renderTokens(toks []string) string {
	var b strings.Builder
	for i, t := range toks {
		if i > 0 {
			prev := toks[i-1]
			switch {
			case isWord(prev) && isWord(t):
				b.WriteString(" ")
			case !isWord(prev) && !isWord(t):
				b.WriteString(" ")
			default:
				if edakit.Int(0, 2) > 0 {
					b.WriteString(" ")
				}
			}
		}
		b.WriteString(t)
	}
	return b.String()
}

// lexing counts the tokens in one line of Verilog. KP 1 uses single-character operators,
// KP 2 multi-character operators that maximal munch keeps whole, KP 3 adds comments, bit
// selects and literals.
func lexing(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	ops1 := []string{"+", "-", "&", "|", "^", "*"}
	ops2 := []string{"==", "!=", "===", "<<", ">>>", "&&", "||", ">=", "**", "<=", "!==", "~^"}
	unary := []string{"~", "!", "~&", "~|", "&", "|"}
	var expr []string
	nOps := edakit.Scaled(ctx, 2, 4)
	operand := func() []string {
		v := edakit.Pick([]string{"a", "b", "c", "d", "en", "cnt", "rdy"})
		if k >= 2 && edakit.Int(0, 3) == 0 {
			return []string{edakit.Pick(unary), v}
		}
		if k == 3 && edakit.Int(0, 2) == 0 {
			return []string{v, "[", edakit.IntStr(edakit.Int(0, 7)), "]"}
		}
		if k == 3 && edakit.Int(0, 3) == 0 {
			return []string{edakit.IntStr(edakit.Int(1, 99))}
		}
		return []string{v}
	}
	expr = append(expr, operand()...)
	for i := 0; i < nOps; i++ {
		op := edakit.Pick(ops1)
		if k >= 2 && edakit.Int(0, 2) > 0 {
			op = edakit.Pick(ops2)
		}
		expr = append(expr, op)
		expr = append(expr, operand()...)
	}
	if edakit.Int(0, 1) == 0 {
		expr = append(append([]string{"("}, expr...), ")")
	}
	var toks []string
	switch edakit.Int(0, 2) {
	case 0:
		toks = append(append([]string{"assign", "y", "="}, expr...), ";")
	case 1:
		toks = append(append([]string{"always", "@", "(", "posedge", "clk", ")", "q", "<="}, expr...), ";")
	default:
		toks = append(append([]string{"if", "("}, expr...), ")", "y", "=", "1", ";")
	}
	line := renderTokens(toks)
	comments := 0
	if k == 3 {
		if edakit.Int(0, 1) == 0 {
			line += " // " + edakit.Pick([]string{"next state", "a <= b here", "x == y"})
		} else {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				line = parts[0] + " /* " + edakit.Pick([]string{"was a && b", "keep", "see spec"}) + " */ " + parts[1]
			}
		}
		comments = 1
	}
	n := len(toks)
	q := "A Verilog lexer reads this line:" + code(line) +
		"How many tokens does it produce? Count every keyword, identifier, number, operator and punctuation mark; comments and whitespace produce none."
	why := "The lexer uses maximal munch: at each position it takes the longest character sequence that forms a token, so operators such as ===, >>> and ~& are single tokens. The tokens are: " + strings.Join(toks, "  ") + edakit.Sprintf(", which is %d.", n)
	if comments > 0 {
		why += " The comment is discarded."
	}
	return edakit.Numeric(q, edakit.IntStr(n), why, map[string]string{"tokens": edakit.IntStr(n)})
}

// Verilog binary precedence, higher binds tighter.
var binPrec = map[string]int{
	"*": 9, "%": 9, "+": 8, "<<": 7, ">>": 7, "<": 6, ">": 6, "==": 5, "!=": 5,
	"&": 4, "^": 3, "|": 2, "&&": 1, "||": 0,
}

type pnode struct {
	op          string
	val         int
	left, right *pnode
	cond        *pnode // for ?:, left is the true arm and right the false arm
}

func (p *pnode) eval() (int, bool) {
	if p.op == "" {
		return p.val, true
	}
	if p.op == "?:" {
		c, ok := p.cond.eval()
		if !ok {
			return 0, false
		}
		if c != 0 {
			return p.left.eval()
		}
		return p.right.eval()
	}
	l, ok1 := p.left.eval()
	r, ok2 := p.right.eval()
	if !ok1 || !ok2 {
		return 0, false
	}
	b2i := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}
	var v int
	switch p.op {
	case "*":
		v = l * r
	case "%":
		if r == 0 {
			return 0, false
		}
		v = l % r
	case "+":
		v = l + r
	case "<<":
		if r > 20 {
			return 0, false
		}
		v = l << r
	case ">>":
		if r > 30 {
			return 0, false
		}
		v = l >> r
	case "<":
		v = b2i(l < r)
	case ">":
		v = b2i(l > r)
	case "==":
		v = b2i(l == r)
	case "!=":
		v = b2i(l != r)
	case "&":
		v = l & r
	case "^":
		v = l ^ r
	case "|":
		v = l | r
	case "&&":
		v = b2i(l != 0 && r != 0)
	case "||":
		v = b2i(l != 0 || r != 0)
	}
	if v < 0 || v > 1<<30 {
		return 0, false
	}
	return v, true
}

func (p *pnode) paren() string {
	if p.op == "" {
		return edakit.IntStr(p.val)
	}
	if p.op == "?:" {
		return "(" + p.cond.paren() + " ? " + p.left.paren() + " : " + p.right.paren() + ")"
	}
	return "(" + p.left.paren() + " " + p.op + " " + p.right.paren() + ")"
}

// parsing asks for the value of an unparenthesized expression under Verilog precedence.
func parsing(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	var opPool []string
	nOps := 2
	switch k {
	case 1:
		opPool = []string{"+", "*", "<<"}
	case 2:
		opPool = []string{"+", "*", "%", "<<", ">>", "&", "|", "^", "<", ">", "=="}
		nOps = edakit.Scaled(ctx, 3, 4)
	default:
		opPool = []string{"+", "*", "&", "|", "^", "<", ">", "==", "!=", "&&", "||"}
		nOps = edakit.Scaled(ctx, 2, 3)
	}
	for {
		mk := func(m int) ([]int, []string) {
			vals := make([]int, m+1)
			ops := make([]string, m)
			for i := range vals {
				vals[i] = edakit.Int(1, 9)
			}
			for i := range ops {
				ops[i] = edakit.Pick(opPool)
				if ops[i] == "<<" || ops[i] == ">>" {
					vals[i+1] = edakit.Int(1, 3)
				}
			}
			return vals, ops
		}
		vals, ops := mk(nOps)
		// Precedence must matter: the left-to-right reading must give another value.
		tree := parseFlatSeq(vals, ops)
		text := flatText(vals, ops)
		var root *pnode = tree
		if k == 3 {
			cv, cops := mk(1)
			cops[0] = edakit.Pick([]string{"<", ">", "==", "!="})
			fv, fops := mk(1)
			cond := parseFlatSeq(cv, cops)
			falseArm := parseFlatSeq(fv, fops)
			root = &pnode{op: "?:", cond: cond, left: tree, right: falseArm}
			text = flatText(cv, cops) + " ? " + text + " : " + flatText(fv, fops)
		}
		got, ok := root.eval()
		if !ok {
			continue
		}
		naive := leftToRight(vals, ops)
		if k < 3 && naive == got {
			continue
		}
		q := "Evaluate this Verilog expression. Every literal is an unsized decimal, comparisons yield 0 or 1, and nothing overflows." +
			code(text) + "What is its value? Answer in decimal."
		why := "Verilog binds * and % tightest, then +, then the shifts, then the relational operators, then == and !=, then &, ^, |, &&, || and finally ?:. Parenthesized, the expression is " +
			root.paren() + edakit.Sprintf(", which evaluates to %d.", got)
		return edakit.Numeric(q, edakit.IntStr(got), why, map[string]string{"left_to_right": edakit.IntStr(naive)})
	}
}

// parseFlatSeq parses a flat operand/operator list with precedence climbing.
func parseFlatSeq(vals []int, ops []string) *pnode {
	pos := 0
	return climbRest(&pos, vals, ops, 0)
}

// climbRest parses the operand at vals[*pos] and every following operator that binds at
// least minPrec.
func climbRest(pos *int, vals []int, ops []string, minPrec int) *pnode {
	lhs := &pnode{val: vals[*pos]}
	for *pos < len(ops) && binPrec[ops[*pos]] >= minPrec {
		op := ops[*pos]
		*pos++
		rhs := climbRest(pos, vals, ops, binPrec[op]+1)
		lhs = &pnode{op: op, left: lhs, right: rhs}
	}
	return lhs
}

func flatText(vals []int, ops []string) string {
	var b strings.Builder
	for i, v := range vals {
		if i > 0 {
			b.WriteString(" " + ops[i-1] + " ")
		}
		b.WriteString(edakit.IntStr(v))
	}
	return b.String()
}

func leftToRight(vals []int, ops []string) int {
	acc := &pnode{val: vals[0]}
	for i, op := range ops {
		acc = &pnode{op: op, left: acc, right: &pnode{val: vals[i+1]}}
	}
	v, _ := acc.eval()
	return v
}

// expression tree for the AST generator.
type enode struct {
	op       string // "" for a leaf
	leaf     string
	children []*enode
}

func (e *enode) key() string {
	if e.op == "" {
		return e.leaf
	}
	parts := []string{e.op}
	for _, c := range e.children {
		parts = append(parts, c.key())
	}
	return "(" + strings.Join(parts, " ") + ")"
}

func (e *enode) text(top bool) string {
	if e.op == "" {
		return e.leaf
	}
	var s string
	switch len(e.children) {
	case 1:
		s = e.op + e.children[0].text(false)
		return s
	case 2:
		s = e.children[0].text(false) + " " + e.op + " " + e.children[1].text(false)
	default:
		s = e.children[0].text(false) + " ? " + e.children[1].text(false) + " : " + e.children[2].text(false)
	}
	if top {
		return s
	}
	return "(" + s + ")"
}

func (e *enode) count() int {
	n := 1
	for _, c := range e.children {
		n += c.count()
	}
	return n
}

func (e *enode) height() int {
	h := 0
	for _, c := range e.children {
		h = max(h, 1+c.height())
	}
	return h
}

func (e *enode) dagCount(seen map[string]bool) {
	if seen[e.key()] {
		return
	}
	seen[e.key()] = true
	for _, c := range e.children {
		c.dagCount(seen)
	}
}

// astShape asks for the node count (KP 1), the height (KP 2) or the node count after
// hash-consing identical subtrees into a DAG (KP 3) of an expression's AST.
func astShape(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	leaves := []string{"a", "b", "c", "d", "1", "2"}
	binOps := []string{"+", "*", "&", "|", "^", "-"}
	var build func(depth int) *enode
	build = func(depth int) *enode {
		if depth == 0 || (depth < 3 && edakit.Int(0, 3) == 0) {
			return &enode{leaf: edakit.Pick(leaves)}
		}
		switch r := edakit.Int(0, 9); {
		case r == 0:
			return &enode{op: "~", children: []*enode{build(depth - 1)}}
		case r == 1 && k != 3:
			return &enode{op: "?:", children: []*enode{{leaf: edakit.Pick([]string{"s", "t"})}, build(depth - 1), build(depth - 1)}}
		default:
			return &enode{op: edakit.Pick(binOps), children: []*enode{build(depth - 1), build(depth - 1)}}
		}
	}
	depth := edakit.Scaled(ctx, 2, 4)
	var root *enode
	for {
		if k == 3 {
			// Reuse a subtree so the DAG genuinely shares structure.
			shared := build(1)
			for shared.op == "" {
				shared = build(1)
			}
			other := build(depth - 1)
			op := edakit.Pick(binOps)
			root = &enode{op: op, children: []*enode{
				{op: edakit.Pick(binOps), children: []*enode{shared, other}},
				{op: edakit.Pick(binOps), children: []*enode{build(1), shared}},
			}}
		} else {
			root = build(depth)
		}
		if root.op != "" && root.count() >= 5 && len(root.text(true)) < 80 {
			break
		}
	}
	text := root.text(true)
	var ask string
	var ans int
	var why string
	switch k {
	case 1:
		ans = root.count()
		ask = "How many nodes does its abstract syntax tree have? Every operator is one node (a ?: is one node with three children) and every operand occurrence is a leaf node; parentheses produce no node."
		why = edakit.Sprintf("An AST keeps operators and operands and drops the parentheses, whose only job was to fix the shape. Counting operators plus operand occurrences gives %d nodes.", ans)
	case 2:
		ans = root.height()
		ask = "What is the height of its abstract syntax tree, measured in edges from the root to the deepest leaf? Parentheses produce no node."
		why = edakit.Sprintf("The root is the operator applied last. Following the deepest chain of nested operators down to an operand gives a height of %d edges.", ans)
	default:
		seen := map[string]bool{}
		root.dagCount(seen)
		ans = len(seen)
		ask = "A compiler hash-conses this expression's AST into a DAG, so identical subtrees (written identically, operand order included) and repeated operands become one shared node. How many nodes does the DAG have?"
		why = edakit.Sprintf("The tree has %d nodes. Merging every repeated operand and every repeated subtree leaves %d distinct nodes, which is exactly the sharing that common subexpression elimination exploits.", root.count(), ans)
	}
	q := "Consider the expression" + code(text) + ask
	return edakit.Numeric(q, edakit.IntStr(ans), why, map[string]string{"tree_nodes": edakit.IntStr(root.count())})
}
