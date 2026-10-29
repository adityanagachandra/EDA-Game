package hls

import (
	"math/bits"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func allocation(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		ka, km, kr := edakit.Int(1, 4), edakit.Int(1, 3), edakit.Int(3, 12)
		aa, am, ar := 10*edakit.Int(3, 12), 100*edakit.Int(3, 9), 10*edakit.Int(2, 6)
		tot := ka*aa + km*am + kr*ar
		q := edakit.Sprintf("An allocation provides %d adders of %d area units each, %d multipliers of %d units each and %d registers of %d units each. Ignoring interconnect, what is the datapath area in units?", ka, aa, km, am, kr, ar)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("%d × %d + %d × %d + %d × %d = %d units.", ka, aa, km, am, kr, ar, tot),
			map[string]string{"area": edakit.IntStr(tot)})
	case 2:
		na, nm := edakit.Int(4, 24), edakit.Int(3, 18)
		L := edakit.Int(3, 8)
		aa, am := 10*edakit.Int(3, 12), 100*edakit.Int(3, 9)
		ka, km := ceilDiv(na, L), ceilDiv(nm, L)
		tot := ka*aa + km*am
		q := edakit.Sprintf("A block has %d single-cycle additions and %d single-cycle multiplications and must finish in %d control steps. An adder costs %d area units and a multiplier %d. Ignoring dependences, what is the smallest functional-unit area any schedule could use?", na, nm, L, aa, am)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("At least the ceiling of %d / %d = %d adders and of %d / %d = %d multipliers are needed, costing %d × %d + %d × %d = %d units.", na, L, ka, nm, L, km, ka, aa, km, am, tot),
			map[string]string{"adders": edakit.IntStr(ka), "multipliers": edakit.IntStr(km)})
	default:
		na, ns := edakit.Int(2, 14), edakit.Int(2, 14)
		L := edakit.Int(2, 6)
		aa := 10 * edakit.Int(3, 8)
		as := aa + 10*edakit.Int(0, 2)
		aas := aa + 10*edakit.Int(1, 4)
		sep := ceilDiv(na, L)*aa + ceilDiv(ns, L)*as
		comb := ceilDiv(na+ns, L) * aas
		best := minInt(sep, comb)
		q := edakit.Sprintf("A block has %d additions and %d subtractions, all single-cycle, to fit in %d control steps. Option 1 allocates separate adders (%d area units) and subtractors (%d units). Option 2 allocates only combined add/subtract units (%d units each), each able to perform either operation. Using the counting lower bound on unit numbers for each option, what is the smaller of the two functional-unit areas?", na, ns, L, aa, as, aas)
		return num(q, edakit.IntStr(best),
			edakit.Sprintf("Option 1: %d adders and %d subtractors cost %d. Option 2: the ceiling of %d / %d = %d combined units cost %d. The smaller is %d.", ceilDiv(na, L), ceilDiv(ns, L), sep, na+ns, L, ceilDiv(na+ns, L), comb, best),
			map[string]string{"separate": edakit.IntStr(sep), "combined": edakit.IntStr(comb)})
	}
}

type ivl struct {
	name string
	s, e int // live in steps s..e inclusive
}

func randIntervals(n, steps int) []ivl {
	out := make([]ivl, n)
	for i := range out {
		s := edakit.Int(1, steps-1)
		e := edakit.Int(s, minInt(steps, s+edakit.Int(0, 4)))
		out[i] = ivl{edakit.Sprintf("v%d", i+1), s, e}
	}
	return out
}

func ivlText(vs []ivl) string {
	var p []string
	for _, v := range vs {
		p = append(p, edakit.Sprintf("%s: [%d, %d]", v.name, v.s, v.e))
	}
	return strings.Join(p, ", ")
}

func liveAt(vs []ivl, t int) int {
	c := 0
	for _, v := range vs {
		if v.s <= t && t <= v.e {
			c++
		}
	}
	return c
}

func maxOverlap(vs []ivl) (int, int) {
	best, at := 0, 1
	for t := 1; t <= 40; t++ {
		if c := liveAt(vs, t); c > best {
			best, at = c, t
		}
	}
	return best, at
}

func lifetimes(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 8)
	S := edakit.Int(5, 8)
	type def struct {
		d    int
		uses []int
	}
	ds := make([]def, n)
	var vs []ivl
	var p []string
	k := whichKP(ctx)
	multi := k == 3
	for i := range ds {
		d := edakit.Int(1, S-1)
		nu := 1
		if multi {
			nu = edakit.Int(1, 3)
		}
		var uses []int
		for j := 0; j < nu; j++ {
			uses = append(uses, edakit.Int(d+1, minInt(S, d+4)))
		}
		sort.Ints(uses)
		ds[i] = def{d, uses}
		last := uses[len(uses)-1]
		vs = append(vs, ivl{edakit.Sprintf("v%d", i+1), d + 1, last})
		if multi {
			du := dedup(uses)
			word := "step"
			if len(du) > 1 {
				word = "steps"
			}
			p = append(p, edakit.Sprintf("v%d written in step %d, read in %s %s", i+1, d, word, joinInts(du)))
		} else {
			p = append(p, edakit.Sprintf("v%d written in step %d, read in step %d", i+1, d, last))
		}
	}
	head := "A schedule produces intermediate values; a value written at the end of step d and last read in step u occupies a register during steps d+1 through u. The values are: " + strings.Join(p, "; ") + "."
	if k == 1 {
		t := edakit.Int(2, S)
		c := liveAt(vs, t)
		return num(head+edakit.Sprintf(" How many values are live (held in registers) during step %d?", t), edakit.IntStr(c),
			edakit.Sprintf("The register lifetimes are %s. %d of them contain step %d.", ivlText(vs), c, t),
			map[string]string{"live": edakit.IntStr(c)})
	}
	m, at := maxOverlap(vs)
	return num(head+" What is the minimum number of registers needed to hold these values?", edakit.IntStr(m),
		edakit.Sprintf("The register lifetimes, from the step after the write to the last read, are %s. The densest step is step %d with %d live values, and for intervals the maximum overlap is always achievable, so %d registers suffice.", ivlText(vs), at, m, m),
		map[string]string{"registers": edakit.IntStr(m)})
}

func dedup(xs []int) []int {
	var out []int
	for i, x := range xs {
		if i == 0 || x != xs[i-1] {
			out = append(out, x)
		}
	}
	return out
}

// leftEdgeAssign sorts by start (ties to the lower index) and places each interval in the
// lowest-numbered register whose last interval ended before it starts.
func leftEdgeAssign(vs []ivl, sorted bool) []int {
	order := make([]int, len(vs))
	for i := range order {
		order[i] = i
	}
	if sorted {
		sort.SliceStable(order, func(a, b int) bool { return vs[order[a]].s < vs[order[b]].s })
	}
	reg := make([]int, len(vs))
	var last [][]ivl
	for _, i := range order {
		placed := false
		for r := range last {
			ok := true
			for _, w := range last[r] {
				if !(w.e < vs[i].s || vs[i].e < w.s) {
					ok = false
					break
				}
			}
			if ok {
				last[r] = append(last[r], vs[i])
				reg[i] = r + 1
				placed = true
				break
			}
		}
		if !placed {
			last = append(last, []ivl{vs[i]})
			reg[i] = len(last)
		}
	}
	return reg
}

func shuffledNames(vs []ivl) []ivl {
	out := edakit.Shuffle(vs)
	return out
}

func leftEdge(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 9)
	vs := randIntervals(n, edakit.Int(6, 9))
	listed := shuffledNames(vs)
	head := "Register lifetimes (live steps, inclusive) are " + ivlText(listed) + "."
	rule := " The left-edge algorithm sorts the values by start step (ties to the lower-numbered value) and gives each one to the lowest-numbered register whose previous value ended before it starts, opening a new register when none is free."
	reg := leftEdgeAssign(vs, true)
	nreg := maxInt(reg...)
	switch whichKP(ctx) {
	case 1:
		return num(head+rule+" How many registers does it use?", edakit.IntStr(nreg),
			edakit.Sprintf("Processing in start order gives registers %s for v1 to v%d, %d registers in all, which equals the maximum number of overlapping lifetimes.", joinInts(reg), n, nreg),
			map[string]string{"registers": edakit.IntStr(nreg)})
	case 2:
		x := edakit.Int(0, n-1)
		return num(head+rule+edakit.Sprintf(" Registers are numbered 1, 2, ... in the order they are opened. Which register does %s get?", vs[x].name), edakit.IntStr(reg[x]),
			edakit.Sprintf("Processing in start order gives registers %s for v1 to v%d, so %s gets register %d.", joinInts(reg), n, vs[x].name, reg[x]),
			map[string]string{"assignment": joinInts(reg)})
	default:
		ff := leftEdgeAssign(listed, false)
		used := maxInt(ff...)
		return num(head+" A careless binder skips the sort: it takes the values in exactly the listed order and gives each to the lowest-numbered compatible register, opening a new one when none is free. How many registers does it use?", edakit.IntStr(used),
			edakit.Sprintf("First-fit in the listed order assigns registers %s, using %d. Sorted left-edge would use %d, the maximum overlap; without the sort, first-fit can exceed it.", joinInts(ff), used, nreg),
			map[string]string{"unsorted": edakit.IntStr(used), "sorted": edakit.IntStr(nreg)})
	}
}

func overlaps(a, b ivl) bool { return !(a.e < b.s || b.e < a.s) }

func intervalColoring(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	vs := randIntervals(n, edakit.Int(6, 8))
	head := "Register lifetimes (live steps, inclusive) are " + ivlText(vs) + ". In the interference graph, two values are adjacent when their lifetimes share a step."
	switch whichKP(ctx) {
	case 1:
		e := 0
		for i := range vs {
			for j := i + 1; j < len(vs); j++ {
				if overlaps(vs[i], vs[j]) {
					e++
				}
			}
		}
		return num(head+" How many edges does the interference graph have?", edakit.IntStr(e),
			edakit.Sprintf("Checking each of the %d pairs for a shared step finds %d overlapping pairs.", n*(n-1)/2, e),
			map[string]string{"edges": edakit.IntStr(e)})
	case 2:
		m, at := maxOverlap(vs)
		return num(head+" What is the chromatic number of the interference graph, the fewest registers that hold all values?", edakit.IntStr(m),
			edakit.Sprintf("An interval graph is perfect, so its chromatic number equals its largest clique, the most lifetimes sharing one step: %d values at step %d.", m, at),
			map[string]string{"chromatic": edakit.IntStr(m)})
	default:
		reg := leftEdgeAssign(vs, true)
		R := maxInt(reg...)
		if R < 2 {
			R = 2
		}
		for k := edakit.Int(1, 2); k > 0; k-- {
			reg[edakit.Int(0, n-1)] = edakit.Int(1, R)
		}
		var p []string
		for i := range reg {
			p = append(p, edakit.Sprintf("%s → R%d", vs[i].name, reg[i]))
		}
		bad := 0
		var pairs []string
		for i := range vs {
			for j := i + 1; j < n; j++ {
				if reg[i] == reg[j] && overlaps(vs[i], vs[j]) {
					bad++
					pairs = append(pairs, vs[i].name+"-"+vs[j].name)
				}
			}
		}
		why := "No adjacent pair shares a register, so the coloring is proper."
		if bad > 0 {
			why = "The adjacent pairs that share a register are " + strings.Join(pairs, ", ") + "."
		}
		return num(head+" A proposed binding is "+strings.Join(p, ", ")+". How many interference edges join two values bound to the same register?", edakit.IntStr(bad),
			why+edakit.Sprintf(" The count is %d.", bad),
			map[string]string{"conflicts": edakit.IntStr(bad)})
	}
}

func cliquePartition(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 9)
	S := edakit.Int(3, 6)
	steps := make([]int, n)
	var p []string
	for i := range steps {
		steps[i] = edakit.Int(1, S)
		p = append(p, edakit.Sprintf("m%d in step %d", i+1, steps[i]))
	}
	head := "A schedule places the single-cycle multiplications as follows: " + strings.Join(p, ", ") + ". In the compatibility graph two multiplications are adjacent when they could share one multiplier."
	switch whichKP(ctx) {
	case 1:
		e := 0
		for i := range steps {
			for j := i + 1; j < n; j++ {
				if steps[i] != steps[j] {
					e++
				}
			}
		}
		return num(head+" How many edges does the compatibility graph have?", edakit.IntStr(e),
			edakit.Sprintf("Two single-cycle operations are compatible exactly when they are in different steps. Of the %d pairs, %d are in different steps.", n*(n-1)/2, e),
			map[string]string{"edges": edakit.IntStr(e)})
	case 2:
		c := map[int]int{}
		m := 0
		for _, s := range steps {
			c[s]++
			m = maxInt(m, c[s])
		}
		return num(head+" Binding partitions the graph into cliques, one per multiplier. What is the fewest cliques that cover all the multiplications?", edakit.IntStr(m),
			edakit.Sprintf("Operations in the same step are mutually incompatible, so each needs its own clique; the busiest step holds %d. Distributing every step's operations over %d units achieves it, so the answer is %d.", m, m, m),
			map[string]string{"units": edakit.IntStr(m)})
	default:
		L := edakit.Int(2, 3)
		var vs []ivl
		var q []string
		for i := 0; i < n; i++ {
			s := edakit.Int(1, S+2)
			vs = append(vs, ivl{edakit.Sprintf("m%d", i+1), s, s + L - 1})
			q = append(q, edakit.Sprintf("m%d starts in step %d", i+1, s))
		}
		m, at := maxOverlap(vs)
		return num(edakit.Sprintf("Multiplications run on non-pipelined multipliers with a latency of %d steps: %s. Two multiplications are compatible when their busy intervals do not overlap. What is the fewest multipliers a clique partition of the compatibility graph can use?", L, strings.Join(q, ", ")), edakit.IntStr(m),
			edakit.Sprintf("Each operation is busy for %d steps from its start. The compatibility graph is the complement of an interval graph, so the minimum clique cover equals the maximum overlap: %d operations busy in step %d.", L, m, at),
			map[string]string{"units": edakit.IntStr(m)})
	}
}

type fuop struct{ l, r string }

func interconnect(ctx generator.GeneratorContext) generator.Problem {
	srcs := []string{"R1", "R2", "R3", "R4", "R5", "R6"}
	n := edakit.Scaled(ctx, 3, 5)
	ops := make([]fuop, n)
	var p []string
	for i := range ops {
		l := edakit.Pick(srcs)
		r := edakit.Pick(srcs)
		for r == l {
			r = edakit.Pick(srcs)
		}
		ops[i] = fuop{l, r}
		p = append(p, edakit.Sprintf("a%d = %s + %s", i+1, l, r))
	}
	head := "Adder A1 is bound to the additions " + strings.Join(p, ", ") + ", written as left operand + right operand. Each input port gets a multiplexer over the distinct registers that feed it; a k-input multiplexer costs k - 1 two-input multiplexers."
	left := map[string]bool{}
	right := map[string]bool{}
	for _, o := range ops {
		left[o.l] = true
		right[o.r] = true
	}
	cost := func(m map[string]bool) int {
		if len(m) == 0 {
			return 0
		}
		return len(m) - 1
	}
	switch whichKP(ctx) {
	case 1:
		port, set := "left", left
		if edakit.Int(0, 1) == 1 {
			port, set = "right", right
		}
		return num(head+edakit.Sprintf(" How many inputs does the multiplexer on the %s port have?", port), edakit.IntStr(len(set)),
			edakit.Sprintf("The %s port reads %s, which are %d distinct registers.", port, strings.Join(keys(set), ", "), len(set)),
			map[string]string{"inputs": edakit.IntStr(len(set))})
	case 2:
		tot := cost(left) + cost(right)
		return num(head+" Keeping every operand on the port written, how many two-input multiplexers do both ports need in total?", edakit.IntStr(tot),
			edakit.Sprintf("The left port reads %d distinct registers and the right port %d, so the cost is (%d - 1) + (%d - 1) = %d.", len(left), len(right), len(left), len(right), tot),
			map[string]string{"left": edakit.IntStr(len(left)), "right": edakit.IntStr(len(right))})
	default:
		best := 1 << 30
		for mask := 0; mask < 1<<n; mask++ {
			L, R := map[string]bool{}, map[string]bool{}
			for i, o := range ops {
				if mask&(1<<i) != 0 {
					L[o.r], R[o.l] = true, true
				} else {
					L[o.l], R[o.r] = true, true
				}
			}
			best = minInt(best, cost(L)+cost(R))
		}
		return num(head+" Addition is commutative, so the binder may swap the operands of any addition. What is the fewest two-input multiplexers both ports can need?", edakit.IntStr(best),
			edakit.Sprintf("As written the ports cost %d. Trying all %d operand orientations, the best one needs %d two-input multiplexers.", cost(left)+cost(right), 1<<n, best),
			map[string]string{"best": edakit.IntStr(best)})
	}
}

func keys(m map[string]bool) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func fsmStates(ctx generator.GeneratorContext) generator.Problem {
	S := edakit.Scaled(ctx, 5, 60)
	states := S + 1
	bin := bits.Len(uint(states - 1))
	q := edakit.Sprintf("An HLS controller is a Moore FSM with one state per control step of a %d-step schedule plus one idle state. How many flip-flops does its state register need with binary encoding?", S)
	return edakit.Choice(q, edakit.IntStr(bin),
		edakit.Sprintf("There are %d states, and binary encoding needs the ceiling of log2(%d) = %d bits. One-hot would need %d flip-flops.", states, states, bin, states),
		edakit.Distinct(bin, 3, states, bin-1, bin+1)...)
}
