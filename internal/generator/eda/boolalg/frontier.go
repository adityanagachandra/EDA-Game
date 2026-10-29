package boolalg

import (
	"math"
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// fastExactCompute serves KP 1 of the fast-versus-exact concept: what a clean random
// simulation run does and does not prove.
func fastExactCompute(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 10, 24)
	d := pow2(edakit.Int(0, n/2))
	p := float64(d) / float64(pow2(n))
	t := int(math.Round(edakit.Pick([]float64{0.5, 1, 2, 3}) / p))
	miss := math.Pow(1-p, float64(t))
	f := func(x float64) string { return edakit.Num(x, 4) }
	ans := "A difference that size would escape with probability about " + f(miss) + ", so equivalence is still unproven"
	ds := []string{
		"A difference that size would escape with probability about " + f(1-miss) + ", so equivalence is still unproven",
		"A difference that size would escape with probability about " + f(p) + ", so equivalence is still unproven",
		"The clean run proves equivalence, so the exact check can be skipped",
	}
	q := edakit.Sprintf("An agent proposes a rewrite of a %d-input block and runs %d independent uniformly random vectors through the original and the rewrite with no mismatch. Suppose the rewrite were in fact wrong on %d of the $2^{%d}$ inputs. Which statement is correct?", n, t, d, n)
	return edakit.Choice(q, ans, edakit.Sprintf("Each vector misses with probability 1 - %d/2^%d, so all %d miss with probability (1 - p)^%d = %s. Simulation can refute equivalence, never prove it; an exact SAT or BDD check must follow.", d, n, t, t, f(miss)), ds...)
}

type circuitSpec struct {
	gates []gateSpec
}

func randCircuit(k int) circuitSpec {
	var c circuitSpec
	for i := 0; i < k; i++ {
		c.gates = append(c.gates, randGate(3))
	}
	return c
}

func (c circuitSpec) clauses() int {
	t := 0
	for _, g := range c.gates {
		t += tseitinClauses(g)
	}
	return t
}

func (c circuitSpec) describe() string {
	counts := map[string]int{}
	var order []string
	for _, g := range c.gates {
		d := g.describe()
		if counts[d] == 0 {
			order = append(order, d)
		}
		counts[d]++
	}
	var parts []string
	for _, d := range order {
		parts = append(parts, gateCount(counts[d], d))
	}
	return strings.Join(parts, ", ")
}

// mutate returns a cover that is a small edit of s, the kind of slip a proposed transform makes.
func mutate(s sop, n int) sop {
	out := make(sop, len(s))
	for i, c := range s {
		out[i] = append(cube(nil), c...)
	}
	i := rand.Intn(len(out))
	switch rand.Intn(4) {
	case 0: // drop a literal
		for _, v := range rand.Perm(n) {
			if out[i][v] != 2 && out[i].lits() > 1 {
				out[i][v] = 2
				break
			}
		}
	case 1: // flip a literal
		for _, v := range rand.Perm(n) {
			if out[i][v] != 2 {
				out[i][v] = 1 - out[i][v]
				break
			}
		}
	case 2: // drop a cube
		if len(out) > 1 {
			out = append(out[:i], out[i+1:]...)
		}
	default: // add a literal
		for _, v := range rand.Perm(n) {
			if out[i][v] == 2 {
				out[i][v] = rand.Intn(2)
				break
			}
		}
	}
	return out
}

func equivalenceGolden(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := nVars(ctx)
		f, s := drawSOP(n)
		var g sop
		if rand.Intn(3) == 0 {
			g = sop(primes(f)) // a different cover of the same function
		} else {
			g = mutate(s, n)
		}
		gt := g.tt(n)
		diff := f.xor(gt)
		ans := edakit.IntStr(diff.ones())
		q := edakit.Sprintf("Over %s, an agent proposes replacing $f = %s$ by $g = %s$. The golden test builds the miter $f \\oplus g$ and asks a SAT solver for inputs that make it 1. How many of the %d input vectors satisfy the miter? (Zero means the transform is correct.)", varList(n), s.latex(), g.latex(), rows(n))
		list := "none"
		if diff.ones() > 0 {
			list = "m(" + mintermList(diff.minterms()) + ")"
		}
		return edakit.Numeric(q, ans, "Compare the two truth tables: they differ on "+list+". Any one of these is a counterexample the solver can return, and the transform is rejected unless the count is 0.", map[string]string{"counterexamples": list})
	case 2:
		a := randCircuit(edakit.Scaled(ctx, 3, 6))
		b := randCircuit(edakit.Scaled(ctx, 2, 5))
		total := a.clauses() + b.clauses() + 4 + 1
		ans := edakit.IntStr(total)
		q := tseitinRule + " An equivalence check compares an original block (" + a.describe() + ") with an agent's rewrite (" + b.describe() + "), both with one output and the same primary inputs. The miter encodes both blocks, one 2-input XOR of the two outputs, and the unit clause asserting the XOR. How many clauses does it have?"
		return edakit.Numeric(q, ans, edakit.Sprintf("Original: %d clauses. Rewrite: %d. XOR: 4. Unit clause: 1. Total %s. The check grows linearly with the two circuits, which is why it can run on every proposed transform.", a.clauses(), b.clauses(), ans), nil)
	default:
		n := nVars(ctx)
		var f, gt tt
		var s, g sop
		for {
			f, s = drawSOP(n)
			g = mutate(s, n)
			gt = g.tt(n)
			if f != gt {
				break
			}
		}
		diff := f.xor(gt)
		// The don't-care set covers some of the mismatches and some other rows.
		dc := tt{n, 0}
		ms := edakit.Shuffle(diff.minterms())
		for _, m := range ms[:edakit.Int(0, len(ms))] {
			dc.w |= 1 << uint(m)
		}
		for i := 0; i < edakit.Int(1, 2); i++ {
			dc.w |= 1 << uint(rand.Intn(rows(n)))
		}
		bad := diff.and(dc.not())
		ans := edakit.IntStr(bad.ones())
		q := edakit.Sprintf("Over %s, an agent replaces $f = %s$ by $g = %s$, citing don't cares: the inputs m(%s) can never occur in the design. Checking equivalence on the care set only, how many input vectors are real counterexamples?", varList(n), s.latex(), g.latex(), mintermList(dc.minterms()))
		return edakit.Numeric(q, ans, edakit.Sprintf("f and g differ on %d vectors; removing those in the don't-care set leaves %s. The don't-care set itself must come from a trusted analysis, or the check proves nothing.", diff.ones(), ans), nil)
	}
}
