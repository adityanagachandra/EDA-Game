package dft

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// dCalculus: one gate on five-valued inputs (KP 1), the value at a faulty gate output (KP 2), a
// two-gate chain (KP 3). The options are always the five values.
func dCalculus(ctx generator.GeneratorContext) generator.Problem {
	names := []string{"a", "b", "c"}
	switch kp(ctx) {
	case 1:
		t := edakit.Pick([]gtype{gAND, gNAND, gOR, gNOR, gXOR})
		ar := 2
		if t != gXOR {
			ar = edakit.Scaled(ctx, 2, 3)
		}
		in := make([]d5, ar)
		hasD := false
		for i := range in {
			in[i] = edakit.Pick(d5vals)
			if in[i].g != in[i].f && in[i].g != 2 {
				hasD = true
			}
		}
		if !hasD {
			in[edakit.Int(0, ar-1)] = edakit.Pick([]d5{{1, 0}, {0, 1}})
		}
		out := evalD(t, in)
		var p []string
		for i := range in {
			p = append(p, names[i]+" = "+in[i].String())
		}
		return edakit.Choice(edakit.Sprintf("In the five-valued D-calculus, what is the output of %s with %s?", aGate(t), strings.Join(p, ", ")),
			out.String(), edakit.Sprintf("Evaluate the fault-free and faulty machines separately: (%s) and (%s) give %s.", pairText(in, true), pairText(in, false), out.String()), d5Options(out.String())...)
	case 2:
		t := edakit.Pick([]gtype{gAND, gNAND, gOR, gNOR})
		in := make([]d5, 2)
		for i := range in {
			x := edakit.Pick([]int{0, 1, 1, 2})
			if t == gOR || t == gNOR {
				x = edakit.Pick([]int{0, 0, 1, 2})
			}
			in[i] = d5{x, x}
		}
		gv := evalD(t, in)
		s := edakit.Int(0, 1)
		var out d5
		switch {
		case gv.g == 2:
			out = d5{2, 2}
		default:
			out = d5{gv.g, s}
		}
		return edakit.Choice(edakit.Sprintf("Gate g = %s(a, b) has its output stuck-at-%d. With a = %s and b = %s, what five-valued value does line g carry?", gname[t], s, in[0].String(), in[1].String()),
			out.String(), edakit.Sprintf("The fault-free output is %s; the faulty machine holds %d. The pair (good, faulty) is written %s.", gv.String(), s, out.String()), d5Options(out.String())...)
	default:
		t1 := edakit.Pick([]gtype{gAND, gNAND, gOR, gNOR, gXOR})
		t2 := edakit.Pick([]gtype{gAND, gNAND, gOR, gNOR})
		a, b, cc := edakit.Pick(d5vals), edakit.Pick(d5vals), edakit.Pick(d5vals)
		if a.g == a.f && b.g == b.f {
			a = edakit.Pick([]d5{{1, 0}, {0, 1}})
		}
		g1 := evalD(t1, []d5{a, b})
		g2 := evalD(t2, []d5{g1, cc})
		return edakit.Choice(edakit.Sprintf("In the five-valued D-calculus, $g_1 = $ %s(a, b) and $g_2 = $ %s($g_1$, c), with a = %s, b = %s, c = %s. What value does $g_2$ carry?", gname[t1], gname[t2], a.String(), b.String(), cc.String()),
			g2.String(), edakit.Sprintf("$g_1$ = %s, then $g_2$ = %s(%s, %s) = %s.", g1.String(), gname[t2], g1.String(), cc.String(), g2.String()), d5Options(g2.String())...)
	}
}

func pairText(in []d5, good bool) string {
	var p []string
	for _, v := range in {
		x := v.f
		if good {
			x = v.g
		}
		if x == 2 {
			p = append(p, "X")
		} else {
			p = append(p, edakit.IntStr(x))
		}
	}
	return strings.Join(p, ", ")
}

// sensitization: pick an exciting vector (KP 1), a detecting vector (KP 2), a vector detecting
// a fanout-branch fault (KP 3).
func sensitization(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	k := kp(ctx)
	for attempt := 0; ; attempt++ {
		if attempt >= generator.MaxDrawAttempts {
			return sensitizationFallback(k)
		}
		c := randCircuit(nIn, nG+1, 2, typesBasic)
		var f fault
		if k == 3 {
			var br []fault
			for _, l := range c.lines() {
				if l.gate >= 0 {
					br = append(br, l)
				}
			}
			if len(br) == 0 {
				continue
			}
			f = edakit.Pick(br)
			f.v = edakit.Int(0, 1)
		} else {
			f = fault{node: edakit.Int(c.nIn, c.n()-1), gate: -1, v: edakit.Int(0, 1)}
		}
		var yes, excOnly, no []string
		for vk := 0; vk < 1<<c.nIn; vk++ {
			v := c.vec(vk)
			good := c.sim(v, nil)
			exc := good[f.node] != f.v
			det := c.detects(v, f)
			s := bits(v)
			switch {
			case k == 1 && exc:
				yes = append(yes, s)
			case k == 1:
				no = append(no, s)
			case det:
				yes = append(yes, s)
			case exc:
				excOnly = append(excOnly, s)
			default:
				no = append(no, s)
			}
		}
		order := strings.Join(c.names[:c.nIn], "")
		if k == 1 {
			if len(yes) == 0 || len(no) < 3 {
				continue
			}
			ds := edakit.Shuffle(no)[:3]
			return edakit.Choice(c.describe()+edakit.Sprintf(" Inputs in the order %s. Which vector activates (excites) %s, ignoring whether the error reaches an output?", order, c.faultName(f)),
				edakit.Pick(yes), edakit.Sprintf("Activation needs the fault-free value of %s to be %d, the opposite of the stuck value.", c.names[f.node], 1-f.v), ds...)
		}
		if len(yes) == 0 || len(excOnly)+len(no) < 3 || len(excOnly) == 0 {
			continue
		}
		pool := append(edakit.Shuffle(excOnly), edakit.Shuffle(no)...)
		ds := []string{pool[0]}
		rest := edakit.Shuffle(pool[1:])
		ds = append(ds, rest[0], rest[1])
		why := edakit.Sprintf("A test must activate the fault (%s = %d in the good circuit) and sensitize a path to an output by setting side inputs to non-controlling values.", c.names[f.node], 1-f.v)
		if k == 3 {
			why += " A branch fault affects only the gate that branch feeds, so the path must start through " + c.names[f.gate] + "."
		}
		return edakit.Choice(c.describe()+edakit.Sprintf(" Inputs in the order %s. Which vector detects %s?", order, c.faultName(f)),
			edakit.Pick(yes), why, ds...)
	}
}

// fiveSim simulates a partial input assignment (2 = X) with fault f injected, returning
// five-valued values per node.
func (c *circuit) fiveSim(in []int, f fault) []d5 {
	val := make([]d5, c.n())
	force := func(node int, v d5) d5 {
		if f.gate < 0 && f.node == node {
			if v.g == 2 {
				return d5{2, 2}
			}
			return d5{v.g, f.v}
		}
		return v
	}
	for i := 0; i < c.nIn; i++ {
		val[i] = force(i, d5{in[i], in[i]})
	}
	for g := c.nIn; g < c.n(); g++ {
		iv := make([]d5, len(c.ins[g]))
		for p, s := range c.ins[g] {
			iv[p] = val[s]
		}
		val[g] = force(g, evalD(c.typ[g], iv))
	}
	return val
}

func isD(v d5) bool { return v.g != 2 && v.g != v.f }

func gateSet(c *circuit, s []int) string {
	if len(s) == 0 {
		return "the empty set"
	}
	var p []string
	for _, g := range s {
		p = append(p, c.names[g])
	}
	return "{" + strings.Join(p, ", ") + "}"
}

// setOptions builds distinct gate-set distractors for a correct set.
func setOptions(c *circuit, correct []int, plausible []int) []string {
	want := gateSet(c, correct)
	seen := map[string]bool{want: true}
	var out []string
	add := func(s []int) {
		sort.Ints(s)
		t := gateSet(c, s)
		if !seen[t] && len(out) < 3 {
			seen[t] = true
			out = append(out, t)
		}
	}
	in := map[int]bool{}
	for _, g := range correct {
		in[g] = true
	}
	for _, g := range plausible {
		if !in[g] {
			add(append(append([]int(nil), correct...), g))
			add([]int{g})
		}
	}
	if len(correct) > 1 {
		add(append([]int(nil), correct[1:]...))
		add(append([]int(nil), correct[:len(correct)-1]...))
	}
	for tries := 0; len(out) < 3 && tries < 200; tries++ {
		var s []int
		for g := c.nIn; g < c.n(); g++ {
			if edakit.Int(0, 2) == 0 {
				s = append(s, g)
			}
		}
		add(s)
	}
	return out
}

// dAlgorithmCalc serves the D-frontier (KP 1) and the J-frontier (KP 2).
func dAlgorithmCalc(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	if kp(ctx) != 2 {
		for attempt := 0; ; attempt++ {
			if attempt >= generator.MaxDrawAttempts {
				return bank(dAlgItems).Generate(ctx)
			}
			c := randCircuit(nIn, nG+1, 2, typesBasic)
			f := fault{node: edakit.Int(0, c.n()-2), gate: -1, v: edakit.Int(0, 1)}
			in := make([]int, c.nIn)
			for i := range in {
				in[i] = edakit.Pick([]int{0, 1, 2})
			}
			val := c.fiveSim(in, f)
			if !isD(val[f.node]) {
				continue
			}
			var df, plaus []int
			for g := c.nIn; g < c.n(); g++ {
				hasD := false
				for _, s := range c.ins[g] {
					if isD(val[s]) {
						hasD = true
					}
				}
				if hasD && val[g].g == 2 {
					df = append(df, g)
				} else if hasD || val[g].g == 2 {
					plaus = append(plaus, g)
				}
			}
			if len(df) == 0 || len(plaus) == 0 {
				continue
			}
			var assign []string
			for i := 0; i < c.nIn; i++ {
				v := "X"
				if in[i] != 2 {
					v = edakit.IntStr(in[i])
				}
				assign = append(assign, c.names[i]+" = "+v)
			}
			return edakit.Choice(c.describe()+edakit.Sprintf(" Targeting %s, the D-algorithm has reached the assignment %s and implied every value it can. What is the D-frontier?", c.faultName(f), strings.Join(assign, ", ")),
				gateSet(c, df), "The D-frontier is the set of gates whose output is still X while at least one input carries D or $\\overline{D}$; they are the places where the error can be pushed forward.", setOptions(c, df, plaus)...)
		}
	}
	for attempt := 0; ; attempt++ {
		if attempt >= generator.MaxDrawAttempts {
			return bank(dAlgItems).Generate(ctx)
		}
		c := randCircuit(nIn, nG+1, 2, typesBasic)
		v := c.vec(edakit.Int(0, (1<<c.nIn)-1))
		full := c.sim(v, nil)
		assigned := map[int]bool{}
		for i := 0; i < c.n(); i++ {
			if edakit.Int(0, 9) < 4 {
				assigned[i] = true
			}
		}
		val := make([]int, c.n())
		for i := 0; i < c.nIn; i++ {
			val[i] = 2
			if assigned[i] {
				val[i] = full[i]
			}
		}
		for g := c.nIn; g < c.n(); g++ {
			iv := make([]int, len(c.ins[g]))
			for p, s := range c.ins[g] {
				iv[p] = val[s]
			}
			val[g] = eval3(c.typ[g], iv)
			if assigned[g] && val[g] == 2 {
				val[g] = full[g]
			}
		}
		var jf, plaus []int
		var desc []string
		for i := 0; i < c.n(); i++ {
			if assigned[i] {
				desc = append(desc, edakit.Sprintf("%s = %d", c.names[i], full[i]))
			}
		}
		for g := c.nIn; g < c.n(); g++ {
			iv := make([]int, len(c.ins[g]))
			for p, s := range c.ins[g] {
				iv[p] = val[s]
			}
			implied := eval3(c.typ[g], iv)
			if assigned[g] && implied == 2 {
				jf = append(jf, g)
			} else if assigned[g] || implied == 2 {
				plaus = append(plaus, g)
			}
		}
		if len(jf) == 0 || len(plaus) == 0 || len(desc) < 2 {
			continue
		}
		return edakit.Choice(c.describe()+edakit.Sprintf(" The D-algorithm has assigned %s; every other line is X, and values have been implied forward from the assigned lines. Which gates form the J-frontier?", strings.Join(desc, ", ")),
			gateSet(c, jf), "The J-frontier holds gates whose output value is assigned but not yet implied by their inputs; each still needs a justification decision on its inputs.", setOptions(c, jf, plaus)...)
	}
}

// podemCalc serves objective selection (KP 1) and SCOAP-guided backtrace (KP 2).
func podemCalc(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	if kp(ctx) == 2 {
		c := randCircuit(nIn, nG+1, 2, typesInv)
		cc0, cc1, _ := c.scoap()
		node := c.outs[0]
		if len(c.outs) > 1 && edakit.Int(0, 1) == 1 {
			node = c.outs[1]
		}
		if edakit.Int(0, 2) == 0 && node > c.nIn+1 {
			node = edakit.Int(c.nIn, node-1)
		}
		v := edakit.Int(0, 1)
		l, x := node, v
		var steps []string
		for l >= c.nIn {
			t := c.typ[l]
			vin := x
			if inverting(t) {
				vin = 1 - x
			}
			in := c.ins[l]
			pick := in[0]
			if len(in) > 1 {
				cv, _ := controlling(t)
				cost := func(s int) int {
					if vin == 0 {
						return cc0[s]
					}
					return cc1[s]
				}
				best := in[0]
				for _, s := range in[1:] {
					if vin == cv && cost(s) < cost(best) {
						best = s
					}
					if vin != cv && cost(s) > cost(best) {
						best = s
					}
				}
				pick = best
			}
			steps = append(steps, edakit.Sprintf("%s=%d", c.names[pick], vin))
			l, x = pick, vin
		}
		var tab []string
		for i := 0; i < c.n(); i++ {
			tab = append(tab, edakit.Sprintf("%s (%d, %d)", c.names[i], cc0[i], cc1[i]))
		}
		ans := edakit.Sprintf("%s = %d", c.names[l], x)
		var ds []string
		ds = append(ds, edakit.Sprintf("%s = %d", c.names[l], 1-x))
		for _, i := range edakit.Shuffle(seq(c.nIn)) {
			if i != l {
				ds = append(ds, edakit.Sprintf("%s = %d", c.names[i], edakit.Int(0, 1)))
			}
		}
		return edakit.Choice(c.describe()+edakit.Sprintf(" SCOAP (CC0, CC1): %s. All inputs are X. PODEM backtraces the objective %s = %d. At each gate it computes the value needed on an input (inverting through NAND, NOR and NOT); if that value is the gate's controlling value it follows the easiest such input (lowest CC for the value), otherwise the hardest; ties go to the first listed input. Which input assignment does the backtrace return?", strings.Join(tab, "; "), c.names[node], v),
			ans, "Backtrace path: "+strings.Join(steps, ", then ")+". One controlling input suffices, so take the easiest; a non-controlling requirement needs every input, so attack the hardest first and fail fast.", ds...)
	}
	for attempt := 0; ; attempt++ {
		if attempt >= generator.MaxDrawAttempts {
			return bank(podemItems).Generate(ctx)
		}
		c := randCircuit(nIn, nG+1, 2, typesBasic)
		f := fault{node: edakit.Int(0, c.n()-2), gate: -1, v: edakit.Int(0, 1)}
		in := make([]int, c.nIn)
		for i := range in {
			in[i] = edakit.Pick([]int{0, 1, 2, 2})
		}
		val := c.fiveSim(in, f)
		site := val[f.node]
		detected := false
		for _, o := range c.outs {
			if isD(val[o]) {
				detected = true
			}
		}
		if detected {
			continue
		}
		var ans string
		var ds []string
		var why string
		if site.g == 2 {
			ans = edakit.Sprintf("%s = %d", c.names[f.node], 1-f.v)
			ds = []string{edakit.Sprintf("%s = %d", c.names[f.node], f.v)}
			why = "The fault is not yet activated, so the objective is to set the fault site to the opposite of the stuck value."
		} else if isD(site) {
			gsel := -1
			for g := c.nIn; g < c.n(); g++ {
				hasD := false
				for _, s := range c.ins[g] {
					if isD(val[s]) {
						hasD = true
					}
				}
				if hasD && val[g].g == 2 {
					gsel = g
					break
				}
			}
			if gsel < 0 {
				continue
			}
			cv, _ := controlling(c.typ[gsel])
			xin := -1
			for _, s := range c.ins[gsel] {
				if val[s].g == 2 {
					xin = s
					break
				}
			}
			if xin < 0 {
				continue
			}
			ans = edakit.Sprintf("%s = %d", c.names[xin], 1-cv)
			ds = []string{edakit.Sprintf("%s = %d", c.names[xin], cv), edakit.Sprintf("%s = %d", c.names[f.node], 1-f.v)}
			why = edakit.Sprintf("The fault is activated, so PODEM picks a D-frontier gate (%s, the first in netlist order) and sets one of its X inputs to the non-controlling value %d.", c.names[gsel], 1-cv)
		} else {
			continue
		}
		for _, g := range edakit.Shuffle(seq(c.n())) {
			if len(ds) >= 3 {
				break
			}
			d := edakit.Sprintf("%s = %d", c.names[g], edakit.Int(0, 1))
			dup := d == ans
			for _, e := range ds {
				if e == d {
					dup = true
				}
			}
			if !dup {
				ds = append(ds, d)
			}
		}
		var assign []string
		for i := 0; i < c.nIn; i++ {
			v := "X"
			if in[i] != 2 {
				v = edakit.IntStr(in[i])
			}
			assign = append(assign, c.names[i]+" = "+v)
		}
		return edakit.Choice(c.describe()+edakit.Sprintf(" PODEM targets %s. The current input assignment is %s, with all values implied. If the fault is activated, PODEM chooses the first D-frontier gate in netlist order and its first X input. What objective does PODEM set next?", c.faultName(f), strings.Join(assign, ", ")),
			ans, why, ds...)
	}
}

func seq(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func tseitin(t gtype, ar int) int {
	switch t {
	case gNOT, gBUF:
		return 2
	case gXOR:
		return 4
	}
	return ar + 1
}

// satATPG: Tseitin clauses of the circuit (KP 1), size of the faulty copy (KP 2), and the full
// miter clause count (KP 3).
func satATPG(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	c := randCircuit(nIn, nG+1, 2, typesAll)
	good := 0
	for g := c.nIn; g < c.n(); g++ {
		good += tseitin(c.typ[g], len(c.ins[g]))
	}
	conv := " Tseitin encoding: a k-input AND, NAND, OR or NOR gives k + 1 clauses, a NOT gives 2, a two-input XOR gives 4."
	k := kp(ctx)
	switch k {
	case 1:
		return edakit.Numeric(c.describe()+conv+" How many clauses encode the fault-free circuit?",
			edakit.IntStr(good), edakit.Sprintf("Summing per gate gives %d clauses; the encoding is linear in circuit size.", good), map[string]string{"good_clauses": edakit.IntStr(good)})
	}
	site := edakit.Int(0, c.n()-1)
	for tries := 0; k == 2 && len(c.tfo(site)) == 0; tries++ {
		site = edakit.Int(0, c.n()-1)
		if tries >= generator.MaxDrawAttempts {
			site = 0
		}
	}
	cone := c.tfo(site)
	if k == 2 {
		return edakit.Numeric(c.describe()+edakit.Sprintf(" A SAT-based ATPG miter for a stuck-at fault on %s duplicates only the gates whose function the fault can change. How many gates does the faulty copy contain?", c.names[site]),
			edakit.IntStr(len(cone)), edakit.Sprintf("Only the transitive fanout of %s can differ: %d gates. Everything else is shared with the fault-free copy.", c.names[site], len(cone)), map[string]string{"cone": edakit.IntStr(len(cone))})
	}
	fc := 0
	for g := range cone {
		fc += tseitin(c.typ[g], len(c.ins[g]))
	}
	m := 0
	for _, o := range c.outs {
		if cone[o] || o == site {
			m++
		}
	}
	or := 0
	if m >= 2 {
		or = m + 1
	}
	total := good + fc + 1 + 4*m + or + 1
	return edakit.Numeric(c.describe()+conv+edakit.Sprintf(" The miter for %s stuck-at-%d encodes the whole fault-free circuit, a faulty copy of only the transitive fanout of %s, one unit clause fixing the faulty %s, a two-input XOR per affected primary output, an OR of those XORs when there are two or more, and one unit clause asserting that the outputs differ. How many clauses does the miter have?", c.names[site], edakit.Int(0, 1), c.names[site], c.names[site]),
		edakit.IntStr(total), edakit.Sprintf("Fault-free %d + faulty cone %d + fault unit 1 + XORs $4 \\times %d$ + OR %d + final unit 1 = %d. A satisfying assignment is a test; UNSAT proves the fault redundant.", good, fc, m, or, total),
		map[string]string{"good_clauses": edakit.IntStr(good), "cone_clauses": edakit.IntStr(fc), "outputs": edakit.IntStr(m)})
}

// scoapCC: controllability of a line.
func scoapCC(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	types := typesBasic
	switch kp(ctx) {
	case 2:
		types = typesInv
	case 3:
		types = typesAll
	}
	c := randCircuit(nIn, nG, 2, types)
	cc0, cc1, _ := c.scoap()
	node := edakit.Int(c.n()-2, c.n()-1)
	v := edakit.Int(0, 1)
	ans := cc0[node]
	if v == 1 {
		ans = cc1[node]
	}
	var tr []string
	for g := c.nIn; g <= node; g++ {
		tr = append(tr, edakit.Sprintf("%s (%d, %d)", c.names[g], cc0[g], cc1[g]))
	}
	return edakit.Numeric(c.describe()+edakit.Sprintf(" Using SCOAP combinational controllability (inputs have CC0 = CC1 = 1, and each gate adds 1), what is CC%d of %s?", v, c.names[node]),
		edakit.IntStr(ans), "Propagating (CC0, CC1) forward in topological order: "+strings.Join(tr, "; ")+".",
		map[string]string{"cc0": edakit.IntStr(cc0[node]), "cc1": edakit.IntStr(cc1[node])})
}

// scoapCO: observability of a line (KP 1), of a fanout stem (KP 2), and the SCOAP difficulty of
// a stuck-at fault (KP 3).
func scoapCO(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	k := kp(ctx)
	for attempt := 0; ; attempt++ {
		c := randCircuit(nIn, nG+1, 2, typesInv)
		cc0, cc1, co := c.scoap()
		var cands []int
		for i := 0; i < c.n(); i++ {
			if c.isOut(i) {
				continue
			}
			if k == 2 && len(c.fan[i]) < 2 && attempt < generator.MaxDrawAttempts {
				continue
			}
			cands = append(cands, i)
		}
		if len(cands) == 0 {
			continue
		}
		node := edakit.Pick(cands)
		var tab []string
		for i := 0; i < c.n(); i++ {
			tab = append(tab, edakit.Sprintf("%s (%d, %d)", c.names[i], cc0[i], cc1[i]))
		}
		var cot []string
		for i := c.n() - 1; i >= 0; i-- {
			if co[i] < 1<<29 {
				cot = append(cot, edakit.Sprintf("CO(%s) = %d", c.names[i], co[i]))
			}
		}
		rule := " Primary outputs have CO = 0; an input of an AND or NAND has CO = CO(output) + the CC1 of the other inputs + 1 (CC0 for OR and NOR); NOT adds 1; a fanout stem takes the minimum over its branches."
		if k == 3 {
			v := edakit.Int(0, 1)
			d := co[node] + cc1[node]
			if v == 1 {
				d = co[node] + cc0[node]
			}
			return edakit.Numeric(c.describe()+" SCOAP (CC0, CC1): "+strings.Join(tab, "; ")+"."+rule+edakit.Sprintf(" The SCOAP estimate of the difficulty of detecting %s stuck-at-%d is the controllability of the opposite value plus the observability. What is it?", c.names[node], v),
				edakit.IntStr(d), edakit.Sprintf("Detection needs %s = %d and observation: CC%d(%s) + CO(%s) = %d. Working backward: %s.", c.names[node], 1-v, 1-v, c.names[node], c.names[node], d, strings.Join(cot, ", ")),
				map[string]string{"co": edakit.IntStr(co[node])})
		}
		return edakit.Numeric(c.describe()+" SCOAP (CC0, CC1): "+strings.Join(tab, "; ")+"."+rule+edakit.Sprintf(" What is CO(%s)?", c.names[node]),
			edakit.IntStr(co[node]), "Working backward from the outputs: "+strings.Join(cot, ", ")+".",
			map[string]string{"co": edakit.IntStr(co[node])})
	}
}

// randomPattern: detection probability by enumeration (KP 1), patterns for a target escape
// probability (KP 2), expected coverage after N patterns (KP 3).
func randomPattern(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		nIn, nG := sizeFor(ctx)
		c := randCircuit(nIn+1, nG+1, 2, typesBasic)
		f, ok := c.testableStemFault(1, 1<<(c.nIn-1))
		if !ok {
			f = fault{node: c.outs[0], gate: -1, v: 1}
		}
		n := c.detectCount(f)
		p := float64(n) / float64(int(1)<<c.nIn)
		ans := edakit.Num(p, 4)
		return edakit.Numeric(c.describe()+edakit.Sprintf(" Inputs are independent and uniformly random. What is the probability that one random vector detects %s? Round to 4 decimal places.", c.faultName(f)),
			ans, edakit.Sprintf("%d of the %d vectors detect it, so $p = %d/%d = %s$.", n, 1<<c.nIn, n, 1<<c.nIn, ans), map[string]string{"detecting": edakit.IntStr(n)})
	case 2:
		k := edakit.Scaled(ctx, 5, 16)
		eps := edakit.Pick([]float64{0.05, 0.01, 0.001})
		p := math.Pow(2, -float64(k))
		n := int(math.Ceil(math.Log(eps) / math.Log(1-p)))
		return edakit.Numeric(edakit.Sprintf("The output stuck-at-0 fault of a %d-input AND gate driven directly by primary inputs needs all inputs at 1. How many independent uniformly random patterns are needed so that the probability of never detecting it is at most %s?", k, edakit.Num(eps, 3)),
			edakit.IntStr(n), edakit.Sprintf("$p = 2^{-%d}$. Escape after N patterns is $(1-p)^N \\le %s$, so $N = \\lceil \\ln %s / \\ln(1-p) \\rceil = %d$.", k, edakit.Num(eps, 3), edakit.Num(eps, 3), n), map[string]string{"p": edakit.Num(p, 8)})
	default:
		m := edakit.Scaled(ctx, 3, 5)
		nPat := edakit.Pick([]int{16, 32, 64, 100, 128, 256})
		var ps []string
		sum := 0.0
		for i := 0; i < m; i++ {
			e := edakit.Int(1, 10)
			p := math.Pow(2, -float64(e))
			sum += 1 - math.Pow(1-p, float64(nPat))
			ps = append(ps, edakit.Sprintf("$2^{-%d}$", e))
		}
		ans := pct(100 * sum / float64(m))
		return edakit.Numeric(edakit.Sprintf("A circuit has %d faults with random-pattern detection probabilities %s. What fault coverage do you expect after %d independent random patterns, in percent? Round to 2 decimal places.", m, strings.Join(ps, ", "), nPat),
			ans, edakit.Sprintf("Each fault is detected with probability $1 - (1-p_i)^{%d}$; the expected coverage is their mean, %s percent. The low-probability faults dominate the gap.", nPat, ans), nil)
	}
}

// compaction: greedy cube merging (KP 1), minimum covering test set (KP 2), reverse-order fault
// simulation (KP 3).
func compaction(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		w := edakit.Scaled(ctx, 5, 8)
		n := edakit.Scaled(ctx, 5, 8)
		cubes := make([][]byte, n)
		var txt []string
		for i := range cubes {
			cubes[i] = make([]byte, w)
			for j := range cubes[i] {
				switch edakit.Int(0, 9) {
				case 0, 1:
					cubes[i][j] = '0'
				case 2, 3:
					cubes[i][j] = '1'
				default:
					cubes[i][j] = 'X'
				}
			}
			txt = append(txt, edakit.Sprintf("t%d = %s", i+1, string(cubes[i])))
		}
		var merged [][]byte
		for _, cb := range cubes {
			placed := false
			for _, m := range merged {
				ok := true
				for j := range cb {
					if cb[j] != 'X' && m[j] != 'X' && cb[j] != m[j] {
						ok = false
					}
				}
				if ok {
					for j := range cb {
						if m[j] == 'X' {
							m[j] = cb[j]
						}
					}
					placed = true
					break
				}
			}
			if !placed {
				merged = append(merged, append([]byte(nil), cb...))
			}
		}
		var mt []string
		for _, m := range merged {
			mt = append(mt, string(m))
		}
		return edakit.Numeric(edakit.Sprintf("ATPG produced test cubes %s (X is don't care). Static compaction processes them in order and merges each cube into the first already-merged cube it is compatible with (no position where one has 0 and the other 1), otherwise it starts a new pattern. How many patterns result?", strings.Join(txt, ", ")),
			edakit.IntStr(len(merged)), "Merged patterns: "+strings.Join(mt, ", ")+". Unspecified bits are what make compaction possible.", map[string]string{"patterns": edakit.IntStr(len(merged))})
	case 2:
		for attempt := 0; ; attempt++ {
			m := edakit.Scaled(ctx, 4, 6)
			nf := edakit.Scaled(ctx, 6, 9)
			det := make([]int, m)
			var txt []string
			for i := range det {
				for f := 0; f < nf; f++ {
					if edakit.Int(0, 9) < 3 {
						det[i] |= 1 << f
					}
				}
				txt = append(txt, edakit.Sprintf("T%d detects {%s}", i+1, maskText(det[i])))
			}
			all := 0
			for _, d := range det {
				all |= d
			}
			best := m + 1
			for s := 1; s < 1<<m; s++ {
				cov := 0
				for i := 0; i < m; i++ {
					if s>>i&1 == 1 {
						cov |= det[i]
					}
				}
				if cov == all {
					if b := popcount(s); b < best {
						best = b
					}
				}
			}
			if (best == m || best < 2) && attempt < generator.MaxDrawAttempts {
				continue
			}
			return edakit.Numeric(edakit.Sprintf("Fault simulation of %d tests over faults f1 to f%d gives: %s. What is the smallest number of tests that still detects every fault detected by the full set?", m, nf, strings.Join(txt, "; ")),
				edakit.IntStr(best), edakit.Sprintf("This is a set-cover problem; checking all subsets, %d tests suffice. Production tools use greedy or ILP heuristics because set cover is NP-hard.", best), nil)
		}
	default:
		m := edakit.Scaled(ctx, 5, 7)
		nf := edakit.Scaled(ctx, 6, 9)
		det := make([]int, m)
		var txt []string
		for i := range det {
			for f := 0; f < nf; f++ {
				if edakit.Int(0, 9) < 3 {
					det[i] |= 1 << f
				}
			}
			if det[i] == 0 {
				det[i] = 1 << edakit.Int(0, nf-1)
			}
			txt = append(txt, edakit.Sprintf("T%d detects {%s}", i+1, maskText(det[i])))
		}
		cov, kept := 0, 0
		var keep []string
		for i := m - 1; i >= 0; i-- {
			if det[i]&^cov != 0 {
				kept++
				cov |= det[i]
				keep = append(keep, edakit.Sprintf("T%d", i+1))
			}
		}
		return edakit.Numeric(edakit.Sprintf("Tests were generated in the order T1 to T%d: %s. Reverse-order fault simulation simulates them from T%d back to T1 with fault dropping and keeps a test only if it detects a fault not yet detected. How many tests are kept?", m, strings.Join(txt, "; "), m),
			edakit.IntStr(kept), "Kept: "+strings.Join(keep, ", ")+". Late tests target hard faults and incidentally detect many easy ones, so early tests often become redundant.", nil)
	}
}

func maskText(x int) string {
	var p []string
	for f := 0; x>>f > 0; f++ {
		if x>>f&1 == 1 {
			p = append(p, edakit.Sprintf("f%d", f+1))
		}
	}
	return strings.Join(p, ", ")
}

func popcount(x int) int {
	n := 0
	for ; x > 0; x &= x - 1 {
		n++
	}
	return n
}

var dAlgItems = []edakit.Item{
	{KP: 3, Q: "The D-algorithm assigns values to internal lines, not only to primary inputs. What is the main consequence for its search?",
		A: "It can make many decisions on internal lines that later prove unjustifiable, so the search space is larger than the input space", D: []string{"It never needs to backtrack", "It cannot handle reconvergent fanout", "It only works for single-output circuits"},
		Why: "Decisions on internal lines must later be justified back to the inputs; conflicts there force backtracking over a space much larger than $2^n$ input vectors."},
	{KP: 3, Q: "During the D-algorithm the D-frontier becomes empty while no output carries D or $\\overline{D}$. What must happen?",
		A: "Backtrack, because the error can no longer reach an output under the current decisions", D: []string{"Declare the fault detected", "Declare the fault redundant immediately", "Justify the J-frontier first and continue"},
		Why: "An empty D-frontier means every path from the fault site is blocked; the algorithm must undo a decision, and only an exhausted search proves redundancy."},
	{KP: 3, Q: "What is a primitive D-cube of failure for a two-input AND gate whose output is stuck-at-0?",
		A: "Inputs 1, 1 with output D", D: []string{"Inputs 0, X with output D", "Inputs 1, 1 with output $\\overline{D}$", "Inputs 0, 0 with output 0"},
		Why: "Activation needs the fault-free output to be 1, which requires both inputs at 1; the good value 1 against the faulty 0 is D."},
	{KP: 3, Q: "What is a propagation D-cube for a two-input NOR gate with D on input a?",
		A: "a = D, b = 0, output = $\\overline{D}$", D: []string{"a = D, b = 1, output = D", "a = D, b = 0, output = D", "a = D, b = X, output = $\\overline{D}$"},
		Why: "The side input must be at the non-controlling value 0 for NOR, and NOR inverts, so D becomes $\\overline{D}$."},
	{KP: 3, Q: "Why does the D-algorithm need consistency checks after each implication?",
		A: "An implied value can contradict an earlier assignment on the same line, signalling a conflict that requires backtracking", D: []string{"Implications never conflict in combinational circuits", "Consistency checks only matter for XOR gates", "They are needed only for multiple faults"},
		Why: "Implication propagates the consequences of a decision; a line implied to both 0 and 1 means the current decision set is infeasible."},
	{KP: 3, Q: "Multiple-path sensitization, where D travels along two reconverging paths at once, matters because:",
		A: "Some faults can only be detected when the error propagates along several paths simultaneously", D: []string{"Single-path sensitization always suffices", "It reduces the number of backtracks to zero", "It is required only for bridging faults"},
		Why: "Schneider's example shows faults that no single sensitized path detects; the D-algorithm's D-frontier formulation handles multiple paths naturally."},
	{KP: 3, Q: "A J-frontier gate is an AND gate with output assigned 0 and both inputs X. How many alternatives does justification have to try in the worst case?",
		A: "Two: set either input to 0", D: []string{"One: both inputs must be 0", "Four: every input combination", "Zero: it is already justified"},
		Why: "One controlling input justifies the output, so the choice is which input to set; if one choice fails later, the other is tried."},
	{KP: 3, Q: "Once the D-algorithm has an output carrying D and an empty J-frontier, what remains?",
		A: "Nothing: the current primary input assignment, with X as don't-care, is a test", D: []string{"It must still justify the D-frontier", "It must re-simulate every fault", "It must assign every X input to 0"},
		Why: "Error at an output plus all assignments justified means the inputs form a test cube; leftover X bits are don't-cares that compaction can use."},
	{KP: 3, Q: "Why was the D-algorithm historically important even though PODEM later outperformed it?",
		A: "It was the first ATPG algorithm guaranteed to find a test for any detectable fault in a combinational circuit", D: []string{"It was the first fault simulator", "It never backtracks", "It works only on sequential circuits"},
		Why: "Roth's D-calculus gave a complete algorithm; later methods kept its value system and changed the search space."},
	{KP: 3, Q: "The D-algorithm is applied to a fault that turns out to be redundant. What does the algorithm eventually report?",
		A: "That no test exists, after exhausting every decision alternative", D: []string{"A test that sets every input to X", "That the fault is detected by every vector", "A conflict on the first implication"},
		Why: "Completeness cuts both ways: an exhausted search is a proof that the fault is untestable, usually because the logic is redundant."},
	{KP: 3, Q: "Which value combination at a gate's inputs blocks a D from propagating through a two-input OR gate?",
		A: "The side input at 1", D: []string{"The side input at 0", "The side input at X", "The side input at D"},
		Why: "1 is the controlling value of OR; it forces the output to 1 in both machines and erases the difference."},
	{KP: 3, Q: "D and $\\overline{D}$ arrive at the two inputs of an AND gate. What is the output?",
		A: "0", D: []string{"D", "$\\overline{D}$", "X"},
		Why: "Good machine: 1 AND 0 = 0. Faulty machine: 0 AND 1 = 0. Both are 0, so the error is cancelled; this is how reconvergence can mask a fault."},
}

var podemItems = []edakit.Item{
	{KP: 3, Q: "Why does PODEM make decisions only on primary inputs?",
		A: "Every assignment is then consistent by construction, and the search space is the $2^n$ input space", D: []string{"Internal lines cannot be assigned in any ATPG", "It avoids the need for implication", "It removes the need for a D-frontier"},
		Why: "With decisions only at inputs, forward simulation implies every internal value, so no justification conflicts arise; the decision tree has depth at most n."},
	{KP: 3, Q: "PODEM's latest input assignment leaves the D-frontier empty with no error at an output. What does it do?",
		A: "Flip the most recent untried input assignment, or pop it and backtrack further", D: []string{"Declare the fault redundant at once", "Continue backtracing from the old objective", "Restart with all inputs at 0"},
		Why: "The decision tree is explored depth-first; flipping the last decision is the cheapest repair, and exhausting the tree proves untestability."},
	{KP: 3, Q: "An objective is (g, 1) for a NAND gate g. What does backtrace need on g's inputs?",
		A: "At least one input at 0, so it follows the easiest input to set to 0", D: []string{"Every input at 1, following the hardest input", "Every input at 0", "At least one input at 1"},
		Why: "NAND outputs 1 when any input is 0, which is the controlling value; one input suffices, so pick the easiest by CC0."},
	{KP: 3, Q: "Why does backtrace pick the hardest input when every input must be set to the non-controlling value?",
		A: "If the hardest one cannot be set, the search fails early before effort is spent on the easy ones", D: []string{"The hardest input always justifies the others", "Hardest inputs have the lowest SCOAP values", "Easy inputs cannot be set to non-controlling values"},
		Why: "All inputs are needed anyway; tackling the most constrained first exposes conflicts sooner, a classic fail-first heuristic."},
	{KP: 3, Q: "What is the X-path check in PODEM?",
		A: "A test that some path of X values still connects a D-frontier gate to an output; if none exists, PODEM backtracks", D: []string{"A check that every input is X", "A check that the fault is activated", "A check that SCOAP values are finite"},
		Why: "If every path to the outputs is already fixed, the error cannot propagate under any completion, so the current branch is dead."},
	{KP: 3, Q: "How does PODEM's worst-case search tree compare with the D-algorithm's?",
		A: "PODEM's has at most $2^n$ leaves for n inputs; the D-algorithm can branch on every internal line", D: []string{"Both are bounded by the number of gates", "PODEM's is always larger", "Neither can backtrack"},
		Why: "Restricting decisions to inputs bounds the tree; the D-algorithm's internal decisions can create a much larger and more conflict-prone space."},
	{KP: 3, Q: "After backtrace returns a = 1, what does PODEM do before choosing the next objective?",
		A: "Simulate the circuit forward with a = 1 (implication) and recompute the fault site, D-frontier and X-paths", D: []string{"Assign every remaining input to 0", "Backtrace again from the same objective without simulating", "Justify the J-frontier"},
		Why: "Implication by simulation is the only propagation step PODEM needs; it determines whether the objective was met and what the next one is."},
	{KP: 3, Q: "PODEM exhausts its decision tree for a fault without finding a test. What has it shown?",
		A: "The fault is redundant: no input vector detects it", D: []string{"The fault needs a sequential test", "The backtrack limit was too low", "The circuit has a timing error"},
		Why: "PODEM is complete for combinational circuits, so an exhausted tree is a proof of untestability; an abort at the backtrack limit is different and proves nothing."},
	{KP: 3, Q: "A production ATPG sets a backtrack limit of 100 per fault. What happens to faults that hit the limit?",
		A: "They are aborted: neither detected nor proven untestable, and they lower test coverage", D: []string{"They are counted as untestable", "They are detected by random patterns automatically", "They are removed from the fault list permanently"},
		Why: "An abort is a resource decision, not a proof; a second pass with a higher limit or a SAT engine often resolves them."},
	{KP: 3, Q: "Why do SCOAP values guide PODEM's backtrace rather than an exact count of solutions?",
		A: "SCOAP is linear-time to compute and gives a usable estimate, while exact counting is as hard as the ATPG problem itself", D: []string{"SCOAP is exact for every circuit", "Exact counts are always smaller", "Backtrace needs no guidance"},
		Why: "Testability measures are heuristics: cheap, approximate, and wrong on reconvergent fanout, but good enough to order decisions."},
	{KP: 3, Q: "PODEM backtraces an objective through a fanout stem that is already 1 on one branch's path. What can go wrong?",
		A: "The new objective may conflict with values required elsewhere, which only shows up after implication", D: []string{"Backtrace cannot pass through fanout stems", "The stem value is ignored by simulation", "The D-frontier becomes larger automatically"},
		Why: "Backtrace is a heuristic that looks at one objective at a time; reconvergent fanout can make its choice incompatible with other requirements, causing backtracks."},
	{KP: 3, Q: "What does PODEM gain from a test cube that leaves some primary inputs at X?",
		A: "The X inputs are don't-cares that compaction and compression can fill to target other faults", D: []string{"Nothing: X inputs must be set to 0", "The test becomes invalid", "It doubles the fault coverage automatically"},
		Why: "Leaving inputs unspecified is what makes dynamic compaction and EDT-style compression effective."},
}

var fanItems = []edakit.Item{
	{KP: 1, Q: "In FAN, what is a head line?",
		A: "A line that is not reachable from any fanout stem but feeds a line that is", D: []string{"A line driven by a fanout stem", "Any primary output", "The line carrying the fault"},
		Why: "Head lines bound fanout-free regions; any value on them can be justified without conflict, so FAN stops backtracing there and justifies them last."},
	{KP: 1, Q: "Why can FAN stop backtrace at head lines instead of primary inputs?",
		A: "The cone feeding a head line is fanout-free, so any value on it can later be justified without conflict", D: []string{"Head lines are always primary inputs", "Head lines carry D values", "Head lines are never observable"},
		Why: "With no reconvergent fanout behind a head line, justification is guaranteed to succeed, so deciding there shrinks the search."},
	{KP: 1, Q: "What does FAN's unique sensitization step do?",
		A: "When every path from the D-frontier to the outputs passes through a single gate, it immediately sets that gate's side inputs to non-controlling values", D: []string{"It sensitizes every path at once", "It picks a random path", "It assigns all inputs to 0"},
		Why: "Dominator-like gates must be passed through anyway, so their side-input values are mandatory and can be implied before any decision."},
	{KP: 1, Q: "Which FAN technique propagates mandatory values both forward and backward as far as possible after each decision?",
		A: "Immediate implication", D: []string{"Unique sensitization", "Multiple backtrace", "Head-line justification"},
		Why: "Doing implications eagerly finds conflicts early and fixes values that would otherwise become decisions."},
	{KP: 2, Q: "What does multiple backtrace in FAN do that PODEM's backtrace does not?",
		A: "It traces several objectives at once and counts requests for 0 and 1 at fanout stems before deciding", D: []string{"It traces only one objective at a time", "It traces from primary inputs to outputs", "It ignores fanout stems"},
		Why: "Collecting requests from all objectives lets FAN pick the stem value most objectives want, reducing conflicts that single backtrace discovers late."},
	{KP: 2, Q: "During multiple backtrace, a fanout stem receives three requests for 1 and one request for 0. What value does FAN choose for the stem?",
		A: "1, the value requested more often", D: []string{"0, to break the tie", "X, because requests conflict", "Both, in two separate searches"},
		Why: "The stem is a decision point; the majority request satisfies more objectives and is less likely to need backtracking."},
	{KP: 2, Q: "Why does FAN usually backtrack less than PODEM on the same circuit?",
		A: "It decides at head lines and stems with more information, and implies mandatory values early", D: []string{"It explores a larger search space", "It disables implication", "It decides only at primary outputs"},
		Why: "Fewer decisions, made at better places with more information, mean fewer wrong guesses to undo."},
	{KP: 2, Q: "A circuit has no reconvergent fanout at all. How do FAN's head lines relate to the primary outputs?",
		A: "The whole circuit is fanout-free, so justification never conflicts and FAN needs no backtracking", D: []string{"Every line becomes a head line and FAN fails", "FAN cannot run without fanout", "The circuit has no testable faults"},
		Why: "Fanout-free circuits are the easy case: any objective can be justified directly, which is exactly why head lines are useful in general circuits."},
	{KP: 3, Q: "Which idea from FAN carried over most directly into modern SAT-based ATPG engines?",
		A: "Learning and propagating implications to find conflicts early", D: []string{"Stopping backtrace at head lines", "Using only primary-input decisions", "Avoiding implication entirely"},
		Why: "Static and dynamic learning of implications, as in SOCRATES, anticipates the unit propagation and clause learning of CDCL solvers."},
	{KP: 3, Q: "Why do production ATPG tools combine a structural engine such as FAN or PODEM with a SAT engine?",
		A: "Structural engines are fast on easy faults; SAT resolves hard and redundant faults that cause structural aborts", D: []string{"SAT engines cannot generate tests", "Structural engines prove redundancy faster than SAT", "They need two engines to handle X values"},
		Why: "Most faults are easy and structural search finds them quickly; SAT's learning shines on the hard tail and on proofs of untestability."},
	{KP: 3, Q: "How does FAN handle a fault whose D-frontier contains several gates?",
		A: "It picks one to propagate through, backtracing objectives for its side inputs, and backtracks to another if that fails", D: []string{"It propagates through all of them simultaneously by default", "It declares the fault untestable", "It ignores the D-frontier"},
		Why: "Like PODEM, FAN selects a frontier gate as its next objective; choices are revisited on backtrack."},
	{KP: 3, Q: "What problem of the D-algorithm does FAN's head-line idea specifically attack?",
		A: "Wasted decisions inside fanout-free regions whose justification could never fail", D: []string{"The lack of a five-valued algebra", "The inability to handle XOR gates", "The need for fault simulation"},
		Why: "Deciding inside a fanout-free cone adds branching without adding information; FAN defers those to a conflict-free final justification step."},
}

// sensitizationFallback is a fixed, verified instance served if random draws keep failing.
func sensitizationFallback(k int) generator.Problem {
	q := "Inputs a, b, c; g1 = AND(a, b); g2 = OR(g1, c); primary output g2. Inputs in the order abc. "
	if k == 1 {
		return edakit.Choice(q+"Which vector activates (excites) g1 stuck-at-0, ignoring whether the error reaches an output?",
			"111", "Activation needs g1 = 1, so a = b = 1.", "000", "010", "100")
	}
	return edakit.Choice(q+"Which vector detects g1 stuck-at-0?",
		"110", "Activation needs a = b = 1, and c = 0 keeps the OR gate from masking the error.", "111", "000", "010")
}
