package cts

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type logicNet struct {
	expr string
	p    func(a, b, c float64) float64
}

var probNets = []logicNet{
	{"$y = (a \\wedge b) \\vee c$", func(a, b, c float64) float64 { return 1 - (1-a*b)*(1-c) }},
	{"$y = (a \\vee b) \\wedge c$", func(a, b, c float64) float64 { return (1 - (1-a)*(1-b)) * c }},
	{"$y = (a \\oplus b) \\wedge c$", func(a, b, c float64) float64 { return (a*(1-b) + b*(1-a)) * c }},
	{"$y = \\overline{a \\wedge b} \\wedge c$", func(a, b, c float64) float64 { return (1 - a*b) * c }},
	{"$y = (a \\wedge b) \\oplus c$", func(a, b, c float64) float64 { return a*b*(1-c) + c*(1-a*b) }},
	{"$y = \\overline{(a \\vee b) \\vee c}$", func(a, b, c float64) float64 { return (1 - a) * (1 - b) * (1 - c) }},
}

// activityPropagation: signal probability, switching activity, and node power.
func activityPropagation(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Pick(probNets)
	pa, pb, pc := pick(0.1, 0.9, 0.05), pick(0.1, 0.9, 0.05), pick(0.1, 0.9, 0.05)
	py := n.p(pa, pb, pc)
	alpha := py * (1 - py)
	desc := edakit.Sprintf("Inputs a, b and c are independent, with signal probabilities (probability of being 1) %s, %s and %s, and each is temporally independent from cycle to cycle. Node %s.", f(pa, 2), f(pb, 2), f(pc, 2), n.expr)
	switch kpOf(ctx) {
	case 1:
		return edakit.Numeric(desc+" What is the signal probability of y? Round to 4 decimal places.", f(py, 4),
			edakit.Sprintf("Propagate with independence: AND multiplies probabilities, OR is one minus the product of complements, XOR is $p_a(1-p_b) + p_b(1-p_a)$, NOT complements. That gives $P(y) = %s$. Reconvergent fanout would break the independence assumption.", f(py, 4)),
			map[string]string{"py": f(py, 5)})
	case 2:
		return edakit.Numeric(desc+" What is the switching activity of y, the probability of a 0-to-1 transition in a cycle? Round to 4 decimal places.", f(alpha, 4),
			edakit.Sprintf("$P(y) = %s$. With temporal independence, a rising transition needs 0 then 1: $\\alpha = (1 - p)p = %s$.", f(py, 4), f(alpha, 4)),
			map[string]string{"alpha": f(alpha, 5)})
	}
	c := pick(2, 20, 0.5)
	v := pick(0.7, 1.0, 0.05)
	fr := pick(0.5, 3, 0.25)
	p := alpha * c * v * v * fr
	return edakit.Numeric(desc+edakit.Sprintf(" Node y drives %s fF at %s V with a %s GHz clock. Using $P = \\alpha C V^2 f$, where $\\alpha$ is the probability of a 0-to-1 transition per cycle, what is its dynamic power, in µW? Round to 3 decimal places.", f(c, 1), f(v, 2), f(fr, 2)), f(p, 3),
		edakit.Sprintf("$P(y) = %s$ and $\\alpha = p(1-p) = %s$. Then $%s \\cdot %s \\cdot %s^2 \\cdot %s = %s$ µW (fF times V squared times GHz gives µW).", f(py, 4), f(alpha, 4), f(alpha, 4), f(c, 1), f(v, 2), f(fr, 2), f(p, 3)),
		map[string]string{"p_uw": f(p, 4)})
}

// transitionDensity: Najm's D(y) = sum_i P(dy/dx_i) D(x_i).
func transitionDensity(ctx generator.GeneratorContext) generator.Problem {
	pa, pb, pc := pick(0.1, 0.9, 0.05), pick(0.1, 0.9, 0.05), pick(0.1, 0.9, 0.05)
	da, db, dc := pick(0.1, 2, 0.1), pick(0.1, 2, 0.1), pick(0.1, 2, 0.1)
	in := edakit.Sprintf("Independent inputs a, b, c have signal probabilities %s, %s, %s and transition densities %s, %s, %s transitions per ns.", f(pa, 2), f(pb, 2), f(pc, 2), f(da, 1), f(db, 1), f(dc, 1))
	switch kpOf(ctx) {
	case 1:
		g := edakit.Pick([]string{"AND", "OR", "XOR"})
		var d float64
		var why string
		switch g {
		case "AND":
			d = pb*da + pa*db
			why = edakit.Sprintf("For $y = ab$, $\\partial y/\\partial a = b$ and $\\partial y/\\partial b = a$, so $D(y) = p_b D_a + p_a D_b = %s$ per ns.", f(d, 3))
		case "OR":
			d = (1-pb)*da + (1-pa)*db
			why = edakit.Sprintf("For $y = a \\vee b$, $\\partial y/\\partial a = \\bar b$, so $D(y) = (1-p_b) D_a + (1-p_a) D_b = %s$ per ns.", f(d, 3))
		default:
			d = da + db
			why = edakit.Sprintf("For XOR the Boolean difference with respect to either input is 1, so every input transition propagates: $D(y) = D_a + D_b = %s$ per ns.", f(d, 3))
		}
		return edakit.Numeric(in+edakit.Sprintf(" Using only a and b, what is the transition density at the output of a 2-input %s gate, in transitions per ns? Round to 3 decimal places.", g), f(d, 3), why,
			map[string]string{"d": f(d, 4)})
	case 2:
		g := edakit.Pick([]string{"AND", "OR"})
		var d float64
		var why string
		if g == "AND" {
			d = pb*pc*da + pa*pc*db + pa*pb*dc
			why = edakit.Sprintf("For $y = abc$, $\\partial y/\\partial a = bc$ and so on: $D(y) = p_b p_c D_a + p_a p_c D_b + p_a p_b D_c = %s$ per ns.", f(d, 3))
		} else {
			d = (1-pb)*(1-pc)*da + (1-pa)*(1-pc)*db + (1-pa)*(1-pb)*dc
			why = edakit.Sprintf("For $y = a \\vee b \\vee c$, an input matters only when the other two are 0: $D(y) = \\bar p_b \\bar p_c D_a + \\bar p_a \\bar p_c D_b + \\bar p_a \\bar p_b D_c = %s$ per ns.", f(d, 3))
		}
		return edakit.Numeric(in+edakit.Sprintf(" What is the transition density at the output of a 3-input %s gate, in transitions per ns? Round to 3 decimal places.", g), f(d, 3), why,
			map[string]string{"d": f(d, 4)})
	}
	pn := pa * pb
	dn := pb*da + pa*db
	d := (1-pc)*dn + (1-pn)*dc
	return edakit.Numeric(in+" Node $n = a \\wedge b$ and output $y = n \\vee c$. What is the transition density of y, in transitions per ns? Round to 3 decimal places.", f(d, 3),
		edakit.Sprintf("First $p_n = p_a p_b = %s$ and $D(n) = p_b D_a + p_a D_b = %s$. Then $D(y) = (1 - p_c) D(n) + (1 - p_n) D_c = %s$ per ns.", f(pn, 4), f(dn, 4), f(d, 3)),
		map[string]string{"pn": f(pn, 4), "dn": f(dn, 4)})
}

// saifPower: computed choice on average power from a toggle count.
func saifPower(ctx generator.GeneratorContext) generator.Problem {
	tc := 1000 * edakit.Int(20, 500)
	c := pick(5, 40, 1)
	v := pick(0.7, 1.0, 0.05)
	t := float64(10 * edakit.Int(5, 50))
	p := 0.5 * float64(tc) * c * v * v / t / 1000
	cs, ds := numOptions(p, 3, " µW", 2*p, p*v, 0.5*float64(tc)*c*v/t/1000)
	q := edakit.Sprintf("A gate-level simulation writes a SAIF file. Over a %s µs window, net n42 toggled %d times (rising plus falling). It drives %s fF from a %s V supply. What average switching power does a vector-based tool report for it?", f(t, 0), tc, f(c, 0), f(v, 2))
	return edakit.Choice(q, cs,
		edakit.Sprintf("Each rising transition draws $CV^2$ from the supply and the falling one draws none, so a toggle costs $\\tfrac{1}{2}CV^2$ on average: $P = \\tfrac{1}{2} \\cdot %d \\cdot %s$ fF $\\cdot %s^2/%s$ µs $= %s$ µW. Forgetting the half double-counts the energy.", tc, f(c, 0), f(v, 2), f(t, 0), f(p, 3)), ds...)
}
