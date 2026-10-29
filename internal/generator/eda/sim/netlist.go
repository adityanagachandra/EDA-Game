package sim

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// gate is one two-input gate of a small combinational netlist. Inputs index into the net list:
// nets 0..nIn-1 are primary inputs, net nIn+i is the output of gate i.
type gate struct {
	op   string
	a, b int
}

type netlist struct {
	nIn   int
	gates []gate
	names []string
}

var gateOps = []string{"AND", "OR", "XOR", "NAND", "NOR"}

func evalOp(op string, x, y int) int {
	switch op {
	case "AND":
		return x & y
	case "OR":
		return x | y
	case "XOR":
		return x ^ y
	case "NAND":
		return 1 - (x & y)
	default:
		return 1 - (x | y)
	}
}

// randomNetlist builds nIn inputs and nGates gates, each gate reading two distinct earlier nets,
// biased toward recent nets so the graph has depth.
func randomNetlist(nIn, nGates int) netlist {
	n := netlist{nIn: nIn}
	for i := 0; i < nIn; i++ {
		n.names = append(n.names, string(rune('a'+i)))
	}
	for g := 0; g < nGates; g++ {
		avail := nIn + g
		pick := func() int {
			if g > 0 && edakit.Int(0, 2) > 0 {
				lo := avail - 3
				if lo < 0 {
					lo = 0
				}
				return edakit.Int(lo, avail-1)
			}
			return edakit.Int(0, avail-1)
		}
		a := pick()
		b := pick()
		for b == a {
			b = edakit.Int(0, avail-1)
		}
		if a > b {
			a, b = b, a
		}
		n.gates = append(n.gates, gate{op: edakit.Pick(gateOps), a: a, b: b})
		n.names = append(n.names, edakit.Sprintf("n%d", g+1))
	}
	return n
}

func (n netlist) settle(in []int) []int {
	v := make([]int, n.nIn+len(n.gates))
	copy(v, in)
	for i, g := range n.gates {
		v[n.nIn+i] = evalOp(g.op, v[g.a], v[g.b])
	}
	return v
}

func (n netlist) describe() string {
	var parts []string
	for i, g := range n.gates {
		parts = append(parts, edakit.Sprintf("%s = %s(%s, %s)", n.names[n.nIn+i], g.op, n.names[g.a], n.names[g.b]))
	}
	return strings.Join(parts, "; ")
}

func (n netlist) fanout() [][]int {
	fo := make([][]int, n.nIn+len(n.gates))
	for i, g := range n.gates {
		fo[g.a] = append(fo[g.a], i)
		fo[g.b] = append(fo[g.b], i)
	}
	return fo
}

// simResult records one event-driven run after a single input toggle with every gate at unit
// delay (or, read as delta cycles, at zero delay).
type simResult struct {
	events      int   // net value changes, the input toggle included
	evals       int   // gate evaluations
	steps       int   // time steps (or delta cycles) in which at least one net changed
	lastTime    int   // time of the last net change
	changes     []int // number of changes per net
	evaluatedOK []bool
}

// simulateToggle runs the classic two-list algorithm: apply the changes of this step, evaluate
// each fanout gate of a changed net once, and schedule any output that differs for the next step.
func (n netlist) simulateToggle(in []int, which int) simResult {
	v := n.settle(in)
	fo := n.fanout()
	res := simResult{changes: make([]int, len(v)), evaluatedOK: make([]bool, len(n.gates))}
	type ev struct{ net, val int }
	cur := []ev{{which, 1 - v[which]}}
	t := 0
	for len(cur) > 0 {
		var changed []int
		for _, e := range cur {
			if v[e.net] != e.val {
				v[e.net] = e.val
				changed = append(changed, e.net)
				res.events++
				res.changes[e.net]++
			}
		}
		if len(changed) == 0 {
			break
		}
		res.steps++
		res.lastTime = t
		toEval := map[int]bool{}
		var order []int
		for _, c := range changed {
			for _, g := range fo[c] {
				if !toEval[g] {
					toEval[g] = true
					order = append(order, g)
				}
			}
		}
		var next []ev
		for _, gi := range order {
			res.evals++
			res.evaluatedOK[gi] = true
			g := n.gates[gi]
			nv := evalOp(g.op, v[g.a], v[g.b])
			if nv != v[n.nIn+gi] {
				next = append(next, ev{n.nIn + gi, nv})
			}
		}
		cur = next
		t++
	}
	return res
}

func inputVector(n int) []int {
	in := make([]int, n)
	for i := range in {
		in[i] = edakit.Int(0, 1)
	}
	return in
}

func (n netlist) inputText(in []int) string {
	var parts []string
	for i := 0; i < n.nIn; i++ {
		parts = append(parts, edakit.Sprintf("%s = %d", n.names[i], in[i]))
	}
	return strings.Join(parts, ", ")
}
