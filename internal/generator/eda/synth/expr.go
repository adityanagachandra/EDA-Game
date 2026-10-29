package synth

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// expr is a small Boolean expression tree over named variables, used by the don't-care,
// equivalence-checking and AIG drills.
type expr struct {
	op   byte // 'v' variable, '!' not, '&', '|', '^'
	v    int
	l, r *expr
}

func ev(e *expr, x []bool) bool {
	switch e.op {
	case 'v':
		return x[e.v]
	case '!':
		return !ev(e.l, x)
	case '&':
		return ev(e.l, x) && ev(e.r, x)
	case '|':
		return ev(e.l, x) || ev(e.r, x)
	}
	return ev(e.l, x) != ev(e.r, x)
}

func (e *expr) tex(names []string) string {
	switch e.op {
	case 'v':
		return names[e.v]
	case '!':
		return "\\overline{" + e.l.tex(names) + "}"
	}
	wrap := func(c *expr) string {
		s := c.tex(names)
		if c.op != 'v' && c.op != '!' && c.op != e.op {
			return "(" + s + ")"
		}
		return s
	}
	sym := map[byte]string{'&': " \\cdot ", '|': " + ", '^': " \\oplus "}[e.op]
	return wrap(e.l) + sym + wrap(e.r)
}

// ands is the AND-node count of the expression in an AIG without sharing.
func (e *expr) ands() int {
	switch e.op {
	case 'v':
		return 0
	case '!':
		return e.l.ands()
	case '^':
		return e.l.ands() + e.r.ands() + 3
	}
	return e.l.ands() + e.r.ands() + 1
}

func (e *expr) uses(v int) bool {
	if e == nil {
		return false
	}
	if e.op == 'v' {
		return e.v == v
	}
	return e.l.uses(v) || e.r.uses(v)
}

func (e *expr) clone() *expr {
	if e == nil {
		return nil
	}
	return &expr{op: e.op, v: e.v, l: e.l.clone(), r: e.r.clone()}
}

func vr(i int) *expr { return &expr{op: 'v', v: i} }

// randExpr builds an expression with the given number of binary operators over vars.
func randExpr(vars []int, ops int, allowXor bool) *expr {
	if ops == 0 {
		e := vr(edakit.Pick(vars))
		if edakit.Int(0, 3) == 0 {
			return &expr{op: '!', l: e}
		}
		return e
	}
	k := edakit.Int(0, ops-1)
	choices := []byte{'&', '|', '&', '|'}
	if allowXor {
		choices = append(choices, '^')
	}
	e := &expr{op: edakit.Pick(choices), l: randExpr(vars, k, allowXor), r: randExpr(vars, ops-1-k, allowXor)}
	if edakit.Int(0, 6) == 0 {
		return &expr{op: '!', l: e}
	}
	return e
}

// table evaluates e over n variables (minterm bit n-1-i holds variable i).
func table(e *expr, n int) uint64 {
	var t uint64
	x := make([]bool, n)
	for m := 0; m < 1<<n; m++ {
		for i := 0; i < n; i++ {
			x[i] = m>>(n-1-i)&1 == 1
		}
		if ev(e, x) {
			t |= 1 << m
		}
	}
	return t
}

// mutate returns a copy of e with one local change: an operator swap, a flipped literal or
// a different variable.
func mutate(e *expr, vars []int) *expr {
	c := e.clone()
	var nodes []*expr
	var walk func(x *expr)
	walk = func(x *expr) {
		if x == nil {
			return
		}
		nodes = append(nodes, x)
		walk(x.l)
		walk(x.r)
	}
	walk(c)
	t := edakit.Pick(nodes)
	switch t.op {
	case 'v':
		t.v = edakit.Pick(vars)
	case '!':
		*t = *t.l
	case '&':
		t.op = '|'
	case '|':
		t.op = '&'
	default:
		t.op = '&'
	}
	return c
}

// deMorgan rewrites the first binary AND or OR it finds into its De Morgan dual, keeping the
// function: x & y = !( !x | !y ).
func deMorgan(e *expr) *expr {
	c := e.clone()
	var done bool
	var walk func(x *expr)
	walk = func(x *expr) {
		if x == nil || done {
			return
		}
		if x.op == '&' || x.op == '|' {
			dual := byte('|')
			if x.op == '|' {
				dual = '&'
			}
			inner := &expr{op: dual, l: &expr{op: '!', l: x.l}, r: &expr{op: '!', l: x.r}}
			*x = expr{op: '!', l: inner}
			done = true
			return
		}
		walk(x.l)
		walk(x.r)
	}
	walk(c)
	return c
}

// gate2 is a named two-input gate for don't-care drills.
type gate2 struct {
	name string
	f    func(a, b bool) bool
	tex  func(a, b string) string
}

var gates2 = []gate2{
	{"AND", func(a, b bool) bool { return a && b }, func(a, b string) string { return a + b }},
	{"OR", func(a, b bool) bool { return a || b }, func(a, b string) string { return a + " + " + b }},
	{"XOR", func(a, b bool) bool { return a != b }, func(a, b string) string { return a + " \\oplus " + b }},
	{"NAND", func(a, b bool) bool { return !(a && b) }, func(a, b string) string { return "\\overline{" + a + b + "}" }},
	{"NOR", func(a, b bool) bool { return !(a || b) }, func(a, b string) string { return "\\overline{" + a + " + " + b + "}" }},
	{"ANDN", func(a, b bool) bool { return a && !b }, func(a, b string) string {
		if len(b) == 1 {
			return a + "\\bar{" + b + "}"
		}
		return a + "\\overline{" + b + "}"
	}},
}

func joinTex(xs []string) string { return strings.Join(xs, ", ") }

// randExprSpread builds an expression whose leaves cycle through a shuffled variable list, so
// no variable repeats until all have been used.
func randExprSpread(vars []int, ops int) *expr {
	order := edakit.Shuffle(vars)
	next := 0
	var build func(k int) *expr
	build = func(k int) *expr {
		if k == 0 {
			v := vr(order[next%len(order)])
			next++
			if edakit.Int(0, 3) == 0 {
				return &expr{op: '!', l: v}
			}
			return v
		}
		left := edakit.Int(0, k-1)
		return &expr{op: edakit.Pick([]byte{'&', '|', '&', '|', '^'}), l: build(left), r: build(k - 1 - left)}
	}
	return build(ops)
}
