package dft

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// stuckAt: count single stuck-at faults including fanout branches (KP 1), count multiple
// stuck-at faults (KP 2), and count the vectors that excite a fault (KP 3).
func stuckAt(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	switch kp(ctx) {
	case 1:
		c := randCircuit(nIn, nG, 2, typesInv)
		ls := c.lines()
		br := len(ls) - c.n()
		ans := edakit.IntStr(2 * len(ls))
		return edakit.Numeric(c.describe()+" Treat every fanout branch of a signal that feeds two or more gate inputs as a line of its own. How many single stuck-at faults does the circuit have?",
			ans, edakit.Sprintf("There are %d stems (%d inputs and %d gate outputs) and %d fanout branches, so %d lines. Each line can be stuck-at-0 or stuck-at-1: %s faults.", c.n(), c.nIn, c.n()-c.nIn, br, len(ls), ans),
			map[string]string{"lines": edakit.IntStr(len(ls)), "branches": edakit.IntStr(br)})
	case 2:
		l := edakit.Scaled(ctx, 3, 9)
		m := int(math.Pow(3, float64(l))) - 1
		return edakit.Numeric(edakit.Sprintf("A small circuit has %d lines. In the multiple stuck-at model each line is independently fault-free, stuck-at-0 or stuck-at-1, and at least one line must be faulty. How many multiple stuck-at faults are there?", l),
			edakit.IntStr(m), edakit.Sprintf("Each of the %d lines has 3 states, giving $3^{%d}$ combinations; removing the fault-free one leaves $3^{%d} - 1 = %d$. The single-fault model keeps only $2 \\cdot %d = %d$.", l, l, l, m, l, 2*l),
			map[string]string{"multiple": edakit.IntStr(m)})
	default:
		c := randCircuit(nIn, nG, 2, typesBasic)
		node := edakit.Int(c.nIn, c.n()-1)
		v := edakit.Int(0, 1)
		n := 0
		for k := 0; k < 1<<c.nIn; k++ {
			if c.sim(c.vec(k), nil)[node] != v {
				n++
			}
		}
		return edakit.Numeric(c.describe()+edakit.Sprintf(" How many of the %d input vectors excite (activate) the fault %s stuck-at-%d, that is, drive the fault-free value of %s to the opposite of the stuck value? Propagation to an output is not required.", 1<<c.nIn, c.names[node], v, c.names[node]),
			edakit.IntStr(n), edakit.Sprintf("Excitation needs %s = %d in the fault-free circuit. Enumerating the truth table, %d vectors do that.", c.names[node], 1-v, n),
			map[string]string{"exciting": edakit.IntStr(n)})
	}
}

// transition: count transition faults (KP 1), choose a two-pattern test (KP 2), count the stem
// transition faults a pair detects (KP 3). All options are computed.
func transition(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	switch kp(ctx) {
	case 1:
		c := randCircuit(nIn, nG, 2, typesInv)
		l := len(c.lines())
		ans := 2 * l
		return edakit.Choice(c.describe()+" Counting every fanout branch as its own line, how many transition faults (slow-to-rise and slow-to-fall) does the transition fault model define?",
			edakit.IntStr(ans), edakit.Sprintf("Each of the %d lines gets a slow-to-rise and a slow-to-fall fault: %d, the same count as single stuck-at faults.", l, ans),
			edakit.Distinct(ans, 3, 4*l, l, 2*c.n(), ans+2)...)
	case 2:
		for attempt := 0; ; attempt++ {
			if attempt >= generator.MaxDrawAttempts {
				return transition(generator.GeneratorContext{Difficulty: ctx.Difficulty, KPOrdinal: 1})
			}
			c := randCircuit(nIn, nG, 2, typesBasic)
			node := edakit.Int(0, c.n()-1)
			rise := edakit.Int(0, 1) == 1
			init := 0
			if !rise {
				init = 1
			}
			sa := fault{node: node, gate: -1, v: init}
			var good, badInit, badLaunch []string
			nv := 1 << c.nIn
			for i := 0; i < nv; i++ {
				for j := 0; j < nv; j++ {
					if i == j {
						continue
					}
					v1, v2 := c.vec(i), c.vec(j)
					ok1 := c.sim(v1, nil)[node] == init
					ok2 := c.detects(v2, sa)
					s := edakit.Sprintf("$V_1$ = %s, $V_2$ = %s", bits(v1), bits(v2))
					switch {
					case ok1 && ok2:
						good = append(good, s)
					case !ok1 && ok2:
						badInit = append(badInit, s)
					case ok1 && !ok2:
						badLaunch = append(badLaunch, s)
					}
				}
			}
			if len(good) == 0 || len(badInit) == 0 || len(badLaunch) < 2 {
				continue
			}
			kind := "slow-to-rise"
			if !rise {
				kind = "slow-to-fall"
			}
			bl := edakit.Shuffle(badLaunch)
			ans := edakit.Pick(good)
			order := strings.Join(c.names[:c.nIn], "")
			return edakit.Choice(c.describe()+edakit.Sprintf(" Vectors list the inputs in the order %s. Which two-pattern test $\\langle V_1, V_2 \\rangle$ detects the %s fault on %s?", order, kind, c.names[node]),
				ans, edakit.Sprintf("$V_1$ must set %s = %d, and $V_2$ must launch the opposite value and propagate it, which is exactly a test for %s stuck-at-%d.", c.names[node], init, c.names[node], init),
				edakit.Pick(badInit), bl[0], bl[1])
		}
	default:
		c := randCircuit(nIn, nG, 2, typesBasic)
		i1 := edakit.Int(0, (1<<c.nIn)-1)
		v1, v2 := c.vec(i1), c.vec((i1+edakit.Int(1, (1<<c.nIn)-1))%(1<<c.nIn))
		g1 := c.sim(v1, nil)
		n := 0
		var hit []string
		for node := 0; node < c.n(); node++ {
			init := g1[node]
			if c.detects(v2, fault{node: node, gate: -1, v: init}) {
				n++
				dir := "slow-to-rise"
				if init == 1 {
					dir = "slow-to-fall"
				}
				hit = append(hit, c.names[node]+" "+dir)
			}
		}
		why := "No stem transition fault is detected."
		if n > 0 {
			why = "Detected: " + strings.Join(hit, ", ") + ". A stem is tested when $V_1$ sets its initial value and $V_2$ detects the matching stuck-at fault."
		}
		order := strings.Join(c.names[:c.nIn], "")
		return edakit.Choice(c.describe()+edakit.Sprintf(" Inputs in the order %s, the pair $V_1$ = %s then $V_2$ = %s is applied at speed. Counting only the %d stems (inputs and gate outputs), how many stem transition faults does the pair detect?", order, bits(v1), bits(v2), c.n()),
			edakit.IntStr(n), why, edakit.Distinct(n, 3, n+1, n-1, n+2, 2*n)...)
	}
}

// pathCount counts the input-to-output paths through the gate graph.
func (c *circuit) pathCount() (total int, to []int, from []int) {
	n := c.n()
	to = make([]int, n)
	from = make([]int, n)
	for i := 0; i < c.nIn; i++ {
		to[i] = 1
	}
	for g := c.nIn; g < n; g++ {
		for _, s := range c.ins[g] {
			to[g] += to[s]
		}
	}
	for g := n - 1; g >= 0; g-- {
		if c.isOut(g) {
			from[g] = 1
		}
		for _, gp := range c.fan[g] {
			from[g] += from[gp[0]]
		}
	}
	for _, o := range c.outs {
		total += to[o]
	}
	return
}

// pathDelay: count path delay faults (KP 1), those through one node (KP 2), and the exponential
// growth through cascaded reconvergent stages (KP 3).
func pathDelay(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	switch kp(ctx) {
	case 1:
		c := randCircuit(nIn, nG+1, 2, typesBasic)
		p, _, _ := c.pathCount()
		ans := edakit.IntStr(2 * p)
		return edakit.Numeric(c.describe()+" A path runs from a primary input through gates to a primary output. Each path carries a rising and a falling path delay fault. How many path delay faults does the circuit have?",
			ans, edakit.Sprintf("Counting paths by dynamic programming (paths into a gate are the sum over its inputs) gives %d input-to-output paths, so %s path delay faults.", p, ans),
			map[string]string{"paths": edakit.IntStr(p)})
	case 2:
		c := randCircuit(nIn, nG+1, 2, typesBasic)
		g := edakit.Int(c.nIn, c.n()-1)
		_, to, from := c.pathCount()
		p := to[g] * from[g]
		ans := edakit.IntStr(2 * p)
		return edakit.Numeric(c.describe()+edakit.Sprintf(" How many path delay faults (rising and falling for each input-to-output path) involve paths that pass through %s?", c.names[g]),
			ans, edakit.Sprintf("Paths through %s = (paths from inputs to %s) times (paths from %s to outputs) = %d times %d = %d; doubled for the two transitions, %s.", c.names[g], c.names[g], c.names[g], to[g], from[g], p, ans),
			map[string]string{"to": edakit.IntStr(to[g]), "from": edakit.IntStr(from[g])})
	default:
		k := edakit.Scaled(ctx, 3, 12)
		w := edakit.Pick([]int{2, 2, 3})
		p := int(math.Pow(float64(w), float64(k)))
		ans := edakit.IntStr(2 * p)
		return edakit.Numeric(edakit.Sprintf("A datapath is a cascade of %d stages. Each stage fans its input out to %d parallel gates whose outputs reconverge at one gate feeding the next stage. Counting rising and falling faults on every end-to-end path, how many path delay faults does the cascade have?", k, w),
			ans, edakit.Sprintf("Each stage multiplies the path count by %d, so there are $%d^{%d} = %d$ paths and %s path delay faults. This exponential growth is why tools test only selected critical paths.", w, w, k, p, ans),
			map[string]string{"paths": edakit.IntStr(p)})
	}
}

// bridging: count candidate bridges (KP 1), count exciting vectors (KP 2), count detecting
// vectors for a wired bridge (KP 3).
func bridging(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1:
		n := edakit.Scaled(ctx, 6, 60)
		models := edakit.Pick([]int{1, 2, 4})
		desc := map[int]string{1: "modeled only as a wired-AND", 2: "modeled as a wired-AND and as a wired-OR", 4: "modeled as a wired-AND, a wired-OR, and the two dominant bridges (each net dominating the other)"}
		pairs := n * (n - 1) / 2
		ans := edakit.IntStr(pairs * models)
		return edakit.Numeric(edakit.Sprintf("A block has %d nets and no layout data, so every unordered pair of distinct nets is a candidate bridge, each %s. How many bridging faults is that?", n, desc[models]),
			ans, edakit.Sprintf("$\\binom{%d}{2} = %d$ pairs times %d models is %s. Real flows extract only physically adjacent pairs from the layout, which cuts this quadratic count to roughly linear.", n, pairs, models, ans),
			map[string]string{"pairs": edakit.IntStr(pairs)})
	}
	nIn, nG := sizeFor(ctx)
	for attempt := 0; ; attempt++ {
		if attempt >= generator.MaxDrawAttempts {
			return bridging(generator.GeneratorContext{Difficulty: ctx.Difficulty, KPOrdinal: 1})
		}
		c := randCircuit(nIn, nG+1, 2, typesBasic)
		x, y := edakit.Int(c.nIn, c.n()-1), edakit.Int(0, c.n()-1)
		if x == y || c.tfi(x)[y] || c.tfi(y)[x] {
			continue
		}
		model := edakit.Pick([]string{"wired-AND", "wired-OR", "dominant"})
		desc := edakit.Sprintf("a %s bridge between %s and %s", model, c.names[x], c.names[y])
		if model == "dominant" {
			desc = edakit.Sprintf("a dominant bridge between %s and %s in which %s takes the value of %s", c.names[x], c.names[y], c.names[y], c.names[x])
		}
		excite, detect := 0, 0
		for vk := 0; vk < 1<<c.nIn; vk++ {
			v := c.vec(vk)
			good := c.sim(v, nil)
			if good[x] != good[y] {
				excite++
			}
			var b0, b1 int
			switch model {
			case "wired-AND":
				b0, b1 = good[x]&good[y], good[x]&good[y]
			case "wired-OR":
				b0, b1 = good[x]|good[y], good[x]|good[y]
			default:
				b0, b1 = good[x], good[x]
			}
			bad := c.simForced(v, map[int]int{x: b0, y: b1})
			for _, o := range c.outs {
				if bad[o] != good[o] {
					detect++
					break
				}
			}
		}
		if k == 2 {
			if excite > 0 {
				return edakit.Numeric(c.describe()+edakit.Sprintf(" A bridge is suspected between %s and %s. How many of the %d input vectors excite it, that is, drive the two nets to opposite values in the fault-free circuit?", c.names[x], c.names[y], 1<<c.nIn),
					edakit.IntStr(excite), edakit.Sprintf("A bridge between nets that carry equal values is invisible. Enumerating the truth table, %s and %s differ on %d vectors.", c.names[x], c.names[y], excite),
					map[string]string{"exciting": edakit.IntStr(excite)})
			}
			continue
		}
		if detect == 0 {
			continue
		}
		return edakit.Numeric(c.describe()+edakit.Sprintf(" Model %s (neither net lies in the other's fanin). How many of the %d input vectors detect it at a primary output?", desc, 1<<c.nIn),
			edakit.IntStr(detect), edakit.Sprintf("For each vector, compute the fault-free values, replace both bridged nets by the bridged value, resimulate the fanout, and compare outputs. %d of the %d vectors that excite the bridge also propagate the error.", detect, excite),
			map[string]string{"exciting": edakit.IntStr(excite), "detecting": edakit.IntStr(detect)})
	}
}

// simForced simulates with the listed nodes forced to fixed values after evaluation.
func (c *circuit) simForced(in []int, force map[int]int) []int {
	val := make([]int, c.n())
	for i := 0; i < c.nIn; i++ {
		val[i] = in[i]
		if v, ok := force[i]; ok {
			val[i] = v
		}
	}
	for g := c.nIn; g < c.n(); g++ {
		iv := make([]int, len(c.ins[g]))
		for p, s := range c.ins[g] {
			iv[p] = val[s]
		}
		val[g] = eval2(c.typ[g], iv)
		if v, ok := force[g]; ok {
			val[g] = v
		}
	}
	return val
}

// cellAware KP 2: pick the input pattern that exposes a cell-internal defect characterized as a
// faulty truth table.
func cellAwareCalc(ctx generator.GeneratorContext) generator.Problem {
	type cell struct {
		name, expr string
		f          func(a, b, c int) int
	}
	cells := []cell{
		{"NAND3", "$Y = \\overline{ABC}$", func(a, b, c int) int { return 1 - a&b&c }},
		{"NOR3", "$Y = \\overline{A + B + C}$", func(a, b, c int) int { return 1 - (a | b | c) }},
		{"AOI21", "$Y = \\overline{AB + C}$", func(a, b, c int) int { return 1 - ((a & b) | c) }},
		{"OAI21", "$Y = \\overline{(A + B)C}$", func(a, b, c int) int { return 1 - ((a | b) & c) }},
		{"MAJ3", "$Y = AB + BC + CA$", func(a, b, c int) int { return (a & b) | (b & c) | (c & a) }},
		{"AO21", "$Y = AB + C$", func(a, b, c int) int { return (a & b) | c }},
	}
	cl := edakit.Pick(cells)
	good := make([]int, 8)
	for k := 0; k < 8; k++ {
		good[k] = cl.f(k>>2&1, k>>1&1, k&1)
	}
	flip := edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7})
	nf := edakit.Int(1, 2)
	bad := append([]int(nil), good...)
	isFlip := map[int]bool{}
	for i := 0; i < nf; i++ {
		bad[flip[i]] = 1 - bad[flip[i]]
		isFlip[flip[i]] = true
	}
	var col []string
	for k := 0; k < 8; k++ {
		col = append(col, edakit.IntStr(bad[k]))
	}
	pat := func(k int) string { return edakit.Sprintf("ABC = %d%d%d", k>>2&1, k>>1&1, k&1) }
	ans := pat(flip[0])
	var ds []string
	for _, k := range edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7}) {
		if !isFlip[k] && len(ds) < 3 {
			ds = append(ds, pat(k))
		}
	}
	q := edakit.Sprintf("Cell-aware characterization of the %s cell (%s) simulates an internal defect at transistor level. The defective cell's output column, for ABC = 000, 001, ..., 111 in order, is %s. Which input pattern detects this defect at the cell output?", cl.name, cl.expr, strings.Join(col, " "))
	return edakit.Choice(q, ans, edakit.Sprintf("Compare against the fault-free column of %s; the defect changes the output only on the rows where the columns differ, and %s is one of them.", cl.name, ans), ds...)
}

// unionFind is a tiny disjoint-set structure for fault collapsing.
type unionFind []int

func newUF(n int) unionFind {
	u := make(unionFind, n)
	for i := range u {
		u[i] = i
	}
	return u
}

func (u unionFind) find(x int) int {
	for u[x] != x {
		u[x] = u[u[x]]
		x = u[x]
	}
	return x
}

func (u unionFind) union(a, b int) { u[u.find(a)] = u.find(b) }

// collapse performs structural equivalence collapsing and then dominance removal. It returns
// the total fault count, the equivalence class count and the count after dominance.
func (c *circuit) collapse() (total, eq, dom int) {
	fs := c.faults()
	idx := map[fault]int{}
	for i, f := range fs {
		idx[f] = i
	}
	total = len(fs)
	// inLine returns the line feeding pin p of gate g.
	inLine := func(g, p, v int) int {
		s := c.ins[g][p]
		if len(c.fan[s]) > 1 {
			return idx[fault{node: s, gate: g, pin: p, v: v}]
		}
		return idx[fault{node: s, gate: -1, v: v}]
	}
	out := func(g, v int) int { return idx[fault{node: g, gate: -1, v: v}] }
	u := newUF(len(fs))
	type domRel struct{ big, small int }
	var doms []domRel
	for g := c.nIn; g < c.n(); g++ {
		t := c.typ[g]
		switch t {
		case gNOT:
			u.union(inLine(g, 0, 0), out(g, 1))
			u.union(inLine(g, 0, 1), out(g, 0))
			continue
		case gBUF:
			u.union(inLine(g, 0, 0), out(g, 0))
			u.union(inLine(g, 0, 1), out(g, 1))
			continue
		case gXOR:
			continue
		}
		cv, _ := controlling(t)
		ov := cv
		if inverting(t) {
			ov = 1 - cv
		}
		for p := range c.ins[g] {
			u.union(inLine(g, p, cv), out(g, ov))
		}
		if len(c.ins[g]) >= 2 {
			for p := range c.ins[g] {
				doms = append(doms, domRel{out(g, 1-ov), inLine(g, p, 1-cv)})
			}
		}
	}
	classes := map[int]bool{}
	for i := range fs {
		classes[u.find(i)] = true
	}
	eq = len(classes)
	removed := map[int]bool{}
	for _, d := range doms {
		if u.find(d.big) != u.find(d.small) {
			removed[u.find(d.big)] = true
		}
	}
	dom = eq - len(removed)
	return
}

// collapsing: equivalence classes (KP 1), after dominance (KP 2), collapse ratio (KP 3).
func collapsing(ctx generator.GeneratorContext) generator.Problem {
	nIn, nG := sizeFor(ctx)
	c := randCircuit(nIn, nG, 2, typesInv)
	total, eq, dom := c.collapse()
	rules := " Count every fanout branch as its own line, so the circuit has " + edakit.IntStr(total) + " single stuck-at faults. Collapse gate by gate: an input stuck at the controlling value is equivalent to the output stuck at the resulting value, and a NOT gate's input and output faults pair up with inverted values; equivalence is transitive."
	facts := map[string]string{"total": edakit.IntStr(total), "equivalence_classes": edakit.IntStr(eq), "after_dominance": edakit.IntStr(dom)}
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(c.describe()+rules+" How many equivalence classes remain?",
			edakit.IntStr(eq), edakit.Sprintf("Merging the gate-local equivalences reduces %d faults to %d classes; one representative per class is enough for test generation.", total, eq), facts)
	case 2:
		return edakit.Numeric(c.describe()+rules+" Then drop every class that contains a gate-output fault dominating a fault on one of that gate's inputs (for an AND, the output stuck-at-1 dominates each input stuck-at-1). How many faults remain?",
			edakit.IntStr(dom), edakit.Sprintf("Equivalence leaves %d classes; dominance then drops %d output-fault classes of multi-input gates, leaving %d. Any test for the dominated input fault also detects the dropped output fault.", eq, eq-dom, dom), facts)
	default:
		r := float64(dom) / float64(total)
		ans := edakit.Num(r, 2)
		return edakit.Numeric(c.describe()+rules+" Then drop every class that contains a gate-output fault dominating an input fault of the same gate. What is the collapse ratio, collapsed count divided by the original count? Round to 2 decimal places.",
			ans, edakit.Sprintf("Equivalence gives %d classes, dominance leaves %d; the ratio is %d/%d = %s.", eq, dom, dom, total, ans), facts)
	}
}

var defectItems = []edakit.Item{
	{KP: 1, Q: "A metal particle shorts two adjacent wires on one die. In test terminology, what is the particle and what is the stuck-at or bridging description?",
		A: "The particle is a defect; the logical description is a fault model", D: []string{"The particle is a fault; the logical description is a defect", "Both are failures", "The particle is an error; the description is a failure"},
		Why: "A defect is the physical imperfection. A fault is its abstraction at a level the tools can reason about, such as a bridge or a stuck line."},
	{KP: 1, Q: "A faulty line produces a wrong value inside the chip, but the wrong value never reaches an output under the applied vector. Which term describes the wrong internal value?",
		A: "An error that was not propagated to a failure", D: []string{"A failure, because the line is faulty", "A defect, because it is physical", "Nothing: the fault is absent for that vector"},
		Why: "Fault, error, failure is a causal chain. An activated fault creates an error; only an error that reaches an observed output is a failure."},
	{KP: 1, Q: "Why do test tools use logical fault models instead of enumerating physical defects?",
		A: "Physical defects are unbounded in kind and size, while a fault model gives a finite, simulatable list", D: []string{"Physical defects never cause logical errors", "Fault models are always exact descriptions of defects", "Defects only occur in memories"},
		Why: "A model such as stuck-at turns an open-ended physical space into a finite fault list that ATPG and fault simulation can target and count."},
	{KP: 1, Q: "A die fails because its threshold voltages drifted and it runs 15 percent slow, with no particle or missing feature. What kind of yield loss is this?",
		A: "Parametric loss, which logical stuck-at tests do not target", D: []string{"A random spot defect", "A stuck-at fault on every line", "A bridging fault between power rails"},
		Why: "Parametric loss comes from process variation shifting device behaviour, not from a localized defect; it is screened by at-speed and parametric tests."},
	{KP: 2, Q: "Which statement about the single stuck-at model is true?",
		A: "It assumes one line is permanently fixed at 0 or 1 while the rest of the circuit is fault-free", D: []string{"It assumes every line may be faulty at once", "It models timing defects directly", "It models only shorts to the power rails, not opens"},
		Why: "The single fault assumption keeps the fault list linear in circuit size; tests for single stuck-at faults also catch most multiple faults in practice."},
	{KP: 2, Q: "A defect is a resistive open in a via that adds 300 ps of delay but leaves logic values correct at slow clocks. Which model captures it?",
		A: "A delay fault model such as transition or path delay", D: []string{"A stuck-at-0 fault on the via's net", "A bridging fault", "A memory coupling fault"},
		Why: "The logic function is intact; only timing changes, so only at-speed tests under a delay fault model can catch it."},
	{KP: 2, Q: "An engineer argues that 99 percent stuck-at coverage means 99 percent of defects are caught. What is wrong with the claim?",
		A: "Coverage is measured against a fault model, and many defects do not behave like the modeled faults", D: []string{"Coverage is always measured against defects directly", "Stuck-at coverage is always pessimistic", "Nothing: the two are identical by definition"},
		Why: "Fault coverage is a model metric. Delay, bridging and cell-internal defects can escape a stuck-at test set, which is why flows add other fault models."},
	{KP: 2, Q: "Why is a fault considered testable only relative to the circuit's observable outputs?",
		A: "A fault is detected only if some input vector makes the faulty and fault-free outputs differ", D: []string{"Every activated fault is automatically detected", "A fault is testable if it exists in the layout", "Testability depends only on the defect's size"},
		Why: "Detection requires both activation and propagation to an observed point; a fault with no such vector is undetectable or redundant."},
	{KP: 3, Q: "A test program finds a die that fails at an elevated temperature but passes at room temperature. Which kind of defect does this suggest, and why does test flow care?",
		A: "A marginal or latent defect, which is why stress and corner testing screen for it", D: []string{"A clean stuck-at fault, which is temperature independent", "A design bug that every die shares", "An ATE calibration error by definition"},
		Why: "Marginal defects such as resistive vias depend on conditions; testing only at nominal conditions lets them escape to the field."},
	{KP: 3, Q: "Every die from a lot fails the same functional vector, at the same output, in the same way. What is the most likely cause?",
		A: "A design or systematic issue rather than random defects", D: []string{"Random particle defects", "A single stuck-at fault on one die", "Normal yield loss from defect density"},
		Why: "Random defects scatter across dies and locations. Identical failure on every die points to a design bug, a test program error or a systematic process problem."},
	{KP: 3, Q: "A chip ships after passing test and fails in a customer system months later from a gate-oxide breakdown. What is this called?",
		A: "A reliability failure from a latent defect, the target of burn-in screening", D: []string{"A test escape caused by low stuck-at coverage only", "A parametric yield loss at wafer probe", "A diagnosis error"},
		Why: "Latent defects pass time-zero tests and wear out in use; burn-in accelerates them so infant-mortality failures happen before shipment."},
	{KP: 3, Q: "Why is the quality of shipped parts, rather than the yield, the figure a test engineer is judged on?",
		A: "Yield measures how many dies are good; quality measures how many bad dies escaped test", D: []string{"Yield and quality are the same number", "Quality depends only on the fab, not on test", "Test cannot affect the escape rate"},
		Why: "Test cannot raise yield, but its coverage sets the defect level of what ships, usually quoted in defective parts per million."},
}

var cellAwareItems = []edakit.Item{
	{KP: 1, Q: "Why can a cell pass a complete stuck-at test set yet still contain an undetected internal defect?",
		A: "Stuck-at faults sit on cell pins, and some transistor-level defects change the output only for specific pin patterns that the set never applies", D: []string{"Stuck-at tests always apply every input pattern of every cell", "Internal defects never affect the logic function", "Cell defects are only detectable by IDDQ"},
		Why: "A stuck-at set needs only a few patterns per cell. A defect that corrupts one specific row of the truth table escapes unless that row is applied and observed."},
	{KP: 1, Q: "What does cell-aware characterization produce for each library cell?",
		A: "A table of defect-to-detecting-pattern mappings from transistor-level simulation of each modeled defect", D: []string{"A new timing model of the cell", "A list of the cell's stuck-at faults only", "The cell's layout extraction for routing"},
		Why: "Each extracted defect (short, open, resistive bridge) is simulated in SPICE across input patterns, and the patterns that expose it become the cell's user-defined fault model."},
	{KP: 1, Q: "Where do the candidate defects for cell-aware analysis come from?",
		A: "The cell's layout: extracted opens and shorts between adjacent shapes and transistor terminals", D: []string{"Random gate-level stuck-at faults", "The RTL description of the cell", "The timing library arcs"},
		Why: "Layout extraction finds the physically plausible defect locations, which is what makes the model more realistic than pin-level faults."},
	{KP: 1, Q: "A drain-source short in one pull-down transistor of a NAND2 forces the output low whenever the other input is high. How does the defect behave at the pins?",
		A: "As a pattern-dependent fault visible only for specific input rows", D: []string{"As a stuck-at-1 on the output for every pattern", "As a pure delay fault", "As an undetectable fault"},
		Why: "The short only matters when the shorted transistor would have been off and the rest of the path conducts, so detection requires a particular input combination."},
	{KP: 1, Q: "Which model is cell-aware test closest to in spirit?",
		A: "A user-defined fault model attached to each cell type, targeted by ATPG like any other fault", D: []string{"A path delay fault on every timing path", "A March test for embedded memories", "A boundary-scan interconnect test"},
		Why: "Cell-aware faults are per-cell UDFM entries listing required input conditions; the ATPG treats them as targets alongside stuck-at faults."},
	{KP: 1, Q: "Why is cell-aware characterization done once per library cell and not per design instance?",
		A: "The defect behaviour depends on the cell's layout and function, which every instance shares", D: []string{"Each instance has a different layout", "Characterization requires the full chip netlist", "Instances never contain defects"},
		Why: "Characterization is expensive transistor-level simulation; amortizing it over every instance of the cell is what makes it practical."},
	{KP: 3, Q: "A transistor open in a CMOS cell leaves the output floating for one input pattern, so it holds its previous value. What test does that require?",
		A: "A two-pattern test: the first pattern sets the output to the opposite value, the second floats it", D: []string{"A single pattern applied twice", "Any stuck-at test of the output", "An IDDQ test only"},
		Why: "A stuck-open makes the cell sequential: the output keeps its old charge, so detection needs an initializing pattern followed by the exposing one."},
	{KP: 3, Q: "Adding cell-aware faults to a design's fault list typically increases the pattern count. Why?",
		A: "Many cell-aware faults require specific input combinations that existing stuck-at patterns did not happen to apply", D: []string{"Cell-aware faults require scan chains to be lengthened", "Each cell-aware fault needs a separate test mode", "Cell-aware ATPG cannot use compaction"},
		Why: "The extra conditions constrain more cells per pattern, so compaction merges fewer faults and the pattern count grows."},
	{KP: 3, Q: "Cell-aware characterization finds a defect whose only exposing pattern is unreachable in the design, because the inputs are tied together. What happens to that fault in the design?",
		A: "It is undetectable in this context and is excluded from test coverage", D: []string{"It must be detected by an extra pattern", "It turns into a stuck-at fault", "It is reported as a timing violation"},
		Why: "Context matters: if the required pin combination cannot occur, the defect cannot change behaviour, so it is classified as untestable."},
	{KP: 3, Q: "Field returns show defects inside cells that stuck-at and transition tests missed. Which change most directly addresses them?",
		A: "Generate patterns for cell-aware (UDFM) faults in addition to stuck-at and transition faults", D: []string{"Run stuck-at ATPG with a higher abort limit", "Increase the scan shift frequency", "Collapse more faults"},
		Why: "If the escapes are cell-internal, only a fault model that represents them steers ATPG to the exposing patterns."},
	{KP: 3, Q: "A resistive bridge inside a cell detected only at speed is characterized. What must the cell-aware entry record?",
		A: "That detection requires a two-cycle at-speed test with specified initial and final pin values", D: []string{"Only the final pin values", "That the defect is a stuck-at-0", "Nothing: resistive bridges are not modeled"},
		Why: "Some cell-internal defects only add delay, so their UDFM entries look like transition conditions on the cell pins."},
	{KP: 3, Q: "Why does cell-aware ATPG matter more at advanced nodes and for automotive parts?",
		A: "FinFET cells have more internal defect sites and automotive quality targets are near zero defective parts per million", D: []string{"Advanced nodes have no stuck-at faults", "Automotive parts are never scan tested", "Cell-aware tests replace burn-in"},
		Why: "Tighter quality targets and denser cells make the cell-internal escapes that stuck-at misses a significant share of the defect level."},
}
