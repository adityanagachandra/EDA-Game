package cmos

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// expr is a series-parallel Boolean expression over distinct literals. A static CMOS gate
// computes the complement of such an expression: its pull-down network maps AND to series and
// OR to parallel, and its pull-up network is the dual.
type expr struct {
	op   byte // 'L' literal, '&' AND, '|' OR
	lit  string
	kids []*expr
}

func (e *expr) isLit() bool { return e.op == 'L' }

// letters hands out A, B, C, ... in order.
type letters struct{ next int }

func (l *letters) take() string {
	s := string(rune('A' + l.next))
	l.next++
	return s
}

func opposite(op byte) byte {
	if op == '&' {
		return '|'
	}
	return '&'
}

// randCompound draws a two-level expression with at least one inner sub-expression, so that
// its pull-down and pull-up networks differ in more than one place.
func randCompound(ctx generator.GeneratorContext) *expr {
	var l letters
	top := edakit.Pick([]byte{'&', '|'})
	n := edakit.Scaled(ctx, 2, 3)
	root := &expr{op: top}
	inner := edakit.Int(0, n-1)
	for i := 0; i < n; i++ {
		if i == inner || (l.next < 4 && edakit.Int(0, 2) == 0) {
			k := 2
			if ctx.Difficulty > 0.7 && edakit.Int(0, 1) == 0 {
				k = 3
			}
			sub := &expr{op: opposite(top)}
			for j := 0; j < k; j++ {
				sub.kids = append(sub.kids, &expr{op: 'L', lit: l.take()})
			}
			root.kids = append(root.kids, sub)
		} else {
			root.kids = append(root.kids, &expr{op: 'L', lit: l.take()})
		}
	}
	return root
}

// flat builds AND or OR of k literals (the NAND-k or NOR-k gate).
func flat(op byte, k int) *expr {
	var l letters
	e := &expr{op: op}
	for i := 0; i < k; i++ {
		e.kids = append(e.kids, &expr{op: 'L', lit: l.take()})
	}
	return e
}

// latex formats the expression in LaTeX: juxtaposition for AND, + for OR.
func (e *expr) latex() string {
	if e.isLit() {
		return e.lit
	}
	var parts []string
	for _, k := range e.kids {
		s := k.latex()
		if e.op == '&' && !k.isLit() && k.op == '|' {
			s = "(" + s + ")"
		}
		parts = append(parts, s)
	}
	if e.op == '&' {
		return strings.Join(parts, "")
	}
	return strings.Join(parts, " + ")
}

// swapped returns the expression with AND and OR exchanged at the chosen depths.
func (e *expr) swapped(flip func(depth int) bool, depth int) *expr {
	if e.isLit() {
		return e
	}
	out := &expr{op: e.op}
	if flip(depth) {
		out.op = opposite(e.op)
	}
	for _, k := range e.kids {
		out.kids = append(out.kids, k.swapped(flip, depth+1))
	}
	return out
}

// net writes a transistor network: S(...) for series, P(...) for parallel. With pullUp false,
// AND is series (the pull-down network); flip exchanges series and parallel at chosen depths.
func (e *expr) net(flip func(depth int) bool, depth int) string {
	if e.isLit() {
		return e.lit
	}
	series := e.op == '&'
	if flip(depth) {
		series = !series
	}
	var parts []string
	for _, k := range e.kids {
		parts = append(parts, k.net(flip, depth+1))
	}
	tag := "P"
	if series {
		tag = "S"
	}
	return tag + "(" + strings.Join(parts, ", ") + ")"
}

var (
	never     = func(int) bool { return false }
	always    = func(int) bool { return true }
	topOnly   = func(d int) bool { return d == 0 }
	innerOnly = func(d int) bool { return d > 0 }
)

func (e *expr) pdn() string { return e.net(never, 0) }
func (e *expr) pun() string { return e.net(always, 0) }

// literals counts literal occurrences, which is the transistor count of each network.
func (e *expr) literals() int {
	if e.isLit() {
		return 1
	}
	n := 0
	for _, k := range e.kids {
		n += k.literals()
	}
	return n
}

// stack is the tallest series stack of the network, with series given by AND when pullUp is
// false and by OR when pullUp is true.
func (e *expr) stack(pullUp bool) int {
	if e.isLit() {
		return 1
	}
	series := (e.op == '&') != pullUp
	best := 0
	for _, k := range e.kids {
		h := k.stack(pullUp)
		if series {
			best += h
		} else if h > best {
			best = h
		}
	}
	return best
}

// widths sizes every transistor so that the worst-case series path has the resistance of one
// unit transistor: a device's width is the length of the longest series path through it. It
// returns literal -> width in units of the unit nMOS (pull-down) or unit pMOS (pull-up).
func (e *expr) widths(pullUp bool) map[string]int {
	out := map[string]int{}
	var walk func(n *expr, ext int)
	walk = func(n *expr, ext int) {
		if n.isLit() {
			out[n.lit] = ext + 1
			return
		}
		series := (n.op == '&') != pullUp
		if !series {
			for _, k := range n.kids {
				walk(k, ext)
			}
			return
		}
		total := 0
		for _, k := range n.kids {
			total += k.stack(pullUp)
		}
		for _, k := range n.kids {
			walk(k, ext+total-k.stack(pullUp))
		}
	}
	walk(e, 0)
	return out
}

// firstLits returns the literal names in order.
func (e *expr) lits() []string {
	if e.isLit() {
		return []string{e.lit}
	}
	var out []string
	for _, k := range e.kids {
		out = append(out, k.lits()...)
	}
	return out
}
