package dft

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// This file holds the small gate-level engine every fault, simulation and ATPG generator draws
// on: a random combinational netlist, two-valued simulation with an injected stuck-at fault,
// the line model (stems and fanout branches), five-valued D-calculus values and SCOAP.

type gtype int

const (
	gAND gtype = iota
	gNAND
	gOR
	gNOR
	gNOT
	gBUF
	gXOR
)

var gname = map[gtype]string{gAND: "AND", gNAND: "NAND", gOR: "OR", gNOR: "NOR", gNOT: "NOT", gBUF: "BUF", gXOR: "XOR"}

// controlling returns the controlling input value of a gate and whether it has one.
func controlling(t gtype) (int, bool) {
	switch t {
	case gAND, gNAND:
		return 0, true
	case gOR, gNOR:
		return 1, true
	}
	return 0, false
}

// inverting reports whether the gate inverts (NAND, NOR, NOT).
func inverting(t gtype) bool { return t == gNAND || t == gNOR || t == gNOT }

// eval2 evaluates a gate on two-valued inputs.
func eval2(t gtype, in []int) int {
	switch t {
	case gNOT:
		return 1 - in[0]
	case gBUF:
		return in[0]
	case gXOR:
		x := 0
		for _, v := range in {
			x ^= v
		}
		return x
	}
	c, _ := controlling(t)
	out := 1 - c
	for _, v := range in {
		if v == c {
			out = c
		}
	}
	if inverting(t) {
		out = 1 - out
	}
	return out
}

// eval3 evaluates a gate on three-valued inputs (0, 1, 2 = X).
func eval3(t gtype, in []int) int {
	switch t {
	case gNOT:
		if in[0] == 2 {
			return 2
		}
		return 1 - in[0]
	case gBUF:
		return in[0]
	case gXOR:
		x := 0
		for _, v := range in {
			if v == 2 {
				return 2
			}
			x ^= v
		}
		return x
	}
	c, _ := controlling(t)
	out := 1 - c
	for _, v := range in {
		if v == c {
			out = c
			break
		}
		if v == 2 {
			out = 2
		}
	}
	if out != 2 && inverting(t) {
		out = 1 - out
	}
	return out
}

type circuit struct {
	nIn   int
	names []string
	typ   []gtype // typ[i] for i >= nIn
	ins   [][]int
	outs  []int
	fan   [][][2]int // fan[node] = list of (gate, pin)
}

func (c *circuit) n() int { return len(c.names) }

// isOut reports whether node is a primary output.
func (c *circuit) isOut(node int) bool {
	for _, o := range c.outs {
		if o == node {
			return true
		}
	}
	return false
}

// randCircuit draws a connected combinational netlist with nIn inputs and nGates gates in which
// every input and every gate is used and the gates with no fanout are the primary outputs.
func randCircuit(nIn, nGates, maxOut int, types []gtype) *circuit {
	for attempt := 0; ; attempt++ {
		if attempt >= generator.MaxDrawAttempts {
			return chainCircuit(nIn, nGates, types)
		}
		c := &circuit{nIn: nIn}
		for i := 0; i < nIn; i++ {
			c.names = append(c.names, string(rune('a'+i)))
			c.typ = append(c.typ, gBUF)
			c.ins = append(c.ins, nil)
		}
		used := make([]int, nIn+nGates)
		for k := 0; k < nGates; k++ {
			t := edakit.Pick(types)
			ar := 2
			if t == gNOT || t == gBUF {
				ar = 1
			}
			avail := len(c.names)
			var unused []int
			for i := 0; i < avail; i++ {
				if used[i] == 0 {
					unused = append(unused, i)
				}
			}
			chosen := map[int]bool{}
			var in []int
			for len(in) < ar {
				var x int
				if len(unused) > 0 && edakit.Int(0, 9) < 7 {
					x = edakit.Pick(unused)
				} else {
					x = edakit.Int(0, avail-1)
				}
				if chosen[x] {
					if len(chosen) >= avail {
						break
					}
					continue
				}
				chosen[x] = true
				in = append(in, x)
			}
			if len(in) < ar {
				break
			}
			sort.Ints(in)
			for _, x := range in {
				used[x]++
			}
			c.names = append(c.names, edakit.Sprintf("g%d", k+1))
			c.typ = append(c.typ, t)
			c.ins = append(c.ins, in)
		}
		if len(c.names) != nIn+nGates {
			continue
		}
		ok := true
		for i := 0; i < nIn; i++ {
			if used[i] == 0 {
				ok = false
			}
		}
		c.outs = nil
		for i := nIn; i < len(c.names); i++ {
			if used[i] == 0 {
				c.outs = append(c.outs, i)
			}
		}
		if !ok || len(c.outs) == 0 || len(c.outs) > maxOut {
			continue
		}
		c.fan = make([][][2]int, len(c.names))
		for g := nIn; g < len(c.names); g++ {
			for p, s := range c.ins[g] {
				c.fan[s] = append(c.fan[s], [2]int{g, p})
			}
		}
		return c
	}
}

// describe renders the netlist as text.
func (c *circuit) describe() string {
	var parts []string
	for g := c.nIn; g < c.n(); g++ {
		var a []string
		for _, s := range c.ins[g] {
			a = append(a, c.names[s])
		}
		parts = append(parts, edakit.Sprintf("%s = %s(%s)", c.names[g], gname[c.typ[g]], strings.Join(a, ", ")))
	}
	var o []string
	for _, x := range c.outs {
		o = append(o, c.names[x])
	}
	var pi []string
	for i := 0; i < c.nIn; i++ {
		pi = append(pi, c.names[i])
	}
	po := "output " + o[0]
	if len(o) > 1 {
		po = "outputs " + strings.Join(o, " and ")
	}
	return "Inputs " + strings.Join(pi, ", ") + "; " + strings.Join(parts, "; ") + "; primary " + po + "."
}

// fault is a single stuck-at fault on a stem (gate < 0) or on the fanout branch that feeds
// input pin `pin` of gate `gate`.
type fault struct {
	node, gate, pin, v int
}

func (c *circuit) faultName(f fault) string {
	return c.lineName(f) + edakit.Sprintf(" stuck-at-%d", f.v)
}

func (c *circuit) lineName(f fault) string {
	if f.gate < 0 {
		return c.names[f.node]
	}
	return edakit.Sprintf("the branch of %s into %s", c.names[f.node], c.names[f.gate])
}

// lines lists every line: a stem per node, plus one branch per fanout connection of a node that
// fans out to two or more gate inputs.
func (c *circuit) lines() []fault {
	var out []fault
	for i := 0; i < c.n(); i++ {
		out = append(out, fault{node: i, gate: -1})
		if len(c.fan[i]) > 1 {
			for _, gp := range c.fan[i] {
				out = append(out, fault{node: i, gate: gp[0], pin: gp[1]})
			}
		}
	}
	return out
}

// faults lists both stuck-at faults on every line.
func (c *circuit) faults() []fault {
	var out []fault
	for _, l := range c.lines() {
		for v := 0; v <= 1; v++ {
			l.v = v
			out = append(out, l)
		}
	}
	return out
}

// sim simulates input vector in with an optional fault and returns every node value.
func (c *circuit) sim(in []int, f *fault) []int {
	val := make([]int, c.n())
	for i := 0; i < c.nIn; i++ {
		val[i] = in[i]
		if f != nil && f.gate < 0 && f.node == i {
			val[i] = f.v
		}
	}
	for g := c.nIn; g < c.n(); g++ {
		iv := make([]int, len(c.ins[g]))
		for p, s := range c.ins[g] {
			iv[p] = val[s]
			if f != nil && f.gate == g && f.pin == p {
				iv[p] = f.v
			}
		}
		val[g] = eval2(c.typ[g], iv)
		if f != nil && f.gate < 0 && f.node == g {
			val[g] = f.v
		}
	}
	return val
}

// detects reports whether vector in detects fault f at some primary output.
func (c *circuit) detects(in []int, f fault) bool {
	good := c.sim(in, nil)
	bad := c.sim(in, &f)
	for _, o := range c.outs {
		if good[o] != bad[o] {
			return true
		}
	}
	return false
}

// vec returns the input vector numbered k (bit nIn-1-i of k is input i).
func (c *circuit) vec(k int) []int {
	v := make([]int, c.nIn)
	for i := 0; i < c.nIn; i++ {
		v[i] = (k >> (c.nIn - 1 - i)) & 1
	}
	return v
}

// detectCount counts the input vectors that detect f.
func (c *circuit) detectCount(f fault) int {
	n := 0
	for k := 0; k < 1<<c.nIn; k++ {
		if c.detects(c.vec(k), f) {
			n++
		}
	}
	return n
}

func bits(v []int) string {
	var b strings.Builder
	for _, x := range v {
		b.WriteByte(byte('0' + x))
	}
	return b.String()
}

// vecText renders an input assignment as "a=1, b=0, ...".
func (c *circuit) vecText(v []int) string {
	var p []string
	for i := 0; i < c.nIn; i++ {
		p = append(p, edakit.Sprintf("%s=%d", c.names[i], v[i]))
	}
	return strings.Join(p, ", ")
}

// tfo returns the set of gates in the transitive fanout of node (excluding node itself).
func (c *circuit) tfo(node int) map[int]bool {
	out := map[int]bool{}
	var walk func(int)
	walk = func(x int) {
		for _, gp := range c.fan[x] {
			if !out[gp[0]] {
				out[gp[0]] = true
				walk(gp[0])
			}
		}
	}
	walk(node)
	return out
}

// tfi returns the set of nodes in the transitive fanin of node (including node itself).
func (c *circuit) tfi(node int) map[int]bool {
	out := map[int]bool{}
	var walk func(int)
	walk = func(x int) {
		if out[x] {
			return
		}
		out[x] = true
		for _, s := range c.ins[x] {
			walk(s)
		}
	}
	walk(node)
	return out
}

// scoap computes Goldstein's combinational controllabilities and observabilities.
func (c *circuit) scoap() (cc0, cc1, co []int) {
	n := c.n()
	cc0, cc1, co = make([]int, n), make([]int, n), make([]int, n)
	for i := 0; i < c.nIn; i++ {
		cc0[i], cc1[i] = 1, 1
	}
	for g := c.nIn; g < n; g++ {
		in := c.ins[g]
		sum0, sum1, min0, min1 := 0, 0, 1<<30, 1<<30
		for _, s := range in {
			sum0 += cc0[s]
			sum1 += cc1[s]
			if cc0[s] < min0 {
				min0 = cc0[s]
			}
			if cc1[s] < min1 {
				min1 = cc1[s]
			}
		}
		switch c.typ[g] {
		case gAND:
			cc0[g], cc1[g] = min0+1, sum1+1
		case gNAND:
			cc0[g], cc1[g] = sum1+1, min0+1
		case gOR:
			cc0[g], cc1[g] = sum0+1, min1+1
		case gNOR:
			cc0[g], cc1[g] = min1+1, sum0+1
		case gNOT:
			cc0[g], cc1[g] = cc1[in[0]]+1, cc0[in[0]]+1
		case gBUF:
			cc0[g], cc1[g] = cc0[in[0]]+1, cc1[in[0]]+1
		case gXOR:
			a, b := in[0], in[1]
			cc0[g] = minInt(cc0[a]+cc0[b], cc1[a]+cc1[b]) + 1
			cc1[g] = minInt(cc0[a]+cc1[b], cc1[a]+cc0[b]) + 1
		}
	}
	// Observability: branch observability per (gate, pin), stem = min over branches.
	const inf = 1 << 30
	for i := range co {
		co[i] = inf
	}
	for _, o := range c.outs {
		co[o] = 0
	}
	for g := n - 1; g >= c.nIn; g-- {
		// co[g] is final once every gate after g has been processed.
		in := c.ins[g]
		for p, s := range in {
			var b int
			switch c.typ[g] {
			case gAND, gNAND:
				b = co[g] + 1
				for q, t := range in {
					if q != p {
						b += cc1[t]
					}
				}
			case gOR, gNOR:
				b = co[g] + 1
				for q, t := range in {
					if q != p {
						b += cc0[t]
					}
				}
			case gNOT, gBUF:
				b = co[g] + 1
			case gXOR:
				t := in[1-p]
				b = co[g] + minInt(cc0[t], cc1[t]) + 1
			}
			if b < co[s] {
				co[s] = b
			}
		}
	}
	return
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// kp resolves the knowledge point to serve: the requested one, or a uniform draw when unguided.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// basic gate type sets used by the generators.
var (
	typesBasic = []gtype{gAND, gNAND, gOR, gNOR}
	typesInv   = []gtype{gAND, gNAND, gOR, gNOR, gNOT}
	typesAll   = []gtype{gAND, gNAND, gOR, gNOR, gNOT, gXOR}
)

// sizeFor picks inputs and gates from the difficulty.
func sizeFor(ctx generator.GeneratorContext) (int, int) {
	nIn := edakit.Scaled(ctx, 3, 4)
	nG := edakit.Scaled(ctx, 3, 6)
	if nG < nIn-1 {
		nG = nIn - 1
	}
	return nIn, nG
}

// five-valued values: pairs (good, faulty) over {0,1,X}.
type d5 struct{ g, f int }

func (v d5) String() string {
	switch {
	case v.g == 2 || v.f == 2:
		if v.g == v.f {
			return "X"
		}
		return "X"
	case v.g == 1 && v.f == 0:
		return "D"
	case v.g == 0 && v.f == 1:
		return "$\\overline{D}$"
	case v.g == 1:
		return "1"
	}
	return "0"
}

var d5vals = []d5{{0, 0}, {1, 1}, {2, 2}, {1, 0}, {0, 1}}

func evalD(t gtype, in []d5) d5 {
	gs := make([]int, len(in))
	fs := make([]int, len(in))
	for i, v := range in {
		gs[i], fs[i] = v.g, v.f
	}
	r := d5{eval3(t, gs), eval3(t, fs)}
	if r.g == 2 || r.f == 2 {
		return d5{2, 2}
	}
	return r
}

// d5Options is the full five-valued option list.
func d5Options(correct string) []string {
	var out []string
	for _, v := range d5vals {
		if s := v.String(); s != correct {
			out = append(out, s)
		}
	}
	return out
}

// faultsDistinct returns the textual names of faults.
func (c *circuit) faultList(fs []fault) string {
	var p []string
	for _, f := range fs {
		p = append(p, c.faultName(f))
	}
	return strings.Join(p, "; ")
}

// stemFaults returns stuck-at faults on stems only.
func (c *circuit) stemFaults() []fault {
	var out []fault
	for i := 0; i < c.n(); i++ {
		for v := 0; v <= 1; v++ {
			out = append(out, fault{node: i, gate: -1, v: v})
		}
	}
	return out
}

// testableStemFault draws a stem fault detected by at least one and at most max vectors.
func (c *circuit) testableStemFault(lo, hi int) (fault, bool) {
	fs := edakit.Shuffle(c.stemFaults())
	for _, f := range fs {
		n := c.detectCount(f)
		if n >= lo && n <= hi {
			return f, true
		}
	}
	return fault{}, false
}

func pct(x float64) string { return edakit.Num(x, 2) }

// aGate returns "a NAND gate" or "an AND gate" for a gate type.
func aGate(t gtype) string {
	n := gname[t]
	if n[0] == 'A' || n[0] == 'O' || n[0] == 'X' {
		return "an " + n + " gate"
	}
	return "a " + n + " gate"
}

// chainCircuit is the deterministic fallback netlist: a chain of two-input gates that uses every
// input once and ends in a single primary output.
func chainCircuit(nIn, nGates int, types []gtype) *circuit {
	var two []gtype
	for _, t := range types {
		if t != gNOT && t != gBUF {
			two = append(two, t)
		}
	}
	if len(two) == 0 {
		two = []gtype{gAND}
	}
	if nGates < nIn-1 {
		nGates = nIn - 1
	}
	c := &circuit{nIn: nIn}
	for i := 0; i < nIn; i++ {
		c.names = append(c.names, string(rune('a'+i)))
		c.typ = append(c.typ, gBUF)
		c.ins = append(c.ins, nil)
	}
	for k := 0; k < nGates; k++ {
		var in []int
		switch {
		case k == 0:
			in = []int{0, 1 % nIn}
		case k+1 < nIn:
			in = []int{k + 1, nIn + k - 1}
		default:
			in = []int{(k + 1) % nIn, nIn + k - 1}
		}
		if in[0] == in[1] {
			in[1] = nIn + k - 1
		}
		if in[0] > in[1] {
			in[0], in[1] = in[1], in[0]
		}
		c.names = append(c.names, edakit.Sprintf("g%d", k+1))
		c.typ = append(c.typ, two[k%len(two)])
		c.ins = append(c.ins, in)
	}
	c.outs = []int{c.n() - 1}
	c.fan = make([][][2]int, c.n())
	for g := nIn; g < c.n(); g++ {
		for p, s := range c.ins[g] {
			c.fan[s] = append(c.fan[s], [2]int{g, p})
		}
	}
	return c
}
