package mleda

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// wa returns the weighted-average smooth wirelength of coordinates xs with smoothing gamma.
func wa(xs []float64, g float64) float64 {
	mx, mn := maxF(xs), xs[0]
	for _, x := range xs {
		mn = math.Min(mn, x)
	}
	var n1, d1, n2, d2 float64
	for _, x := range xs {
		e1 := math.Exp((x - mx) / g)
		e2 := math.Exp(-(x - mn) / g)
		n1 += x * e1
		d1 += e1
		n2 += x * e2
		d2 += e2
	}
	return n1/d1 - n2/d2
}

// differentiablePlacement: WA wirelength of a net (KP 1), its gap to HPWL (KP 2), and one
// gradient step on a cell coordinate (KP 3).
func differentiablePlacement(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		k := kp(ctx)
		if ctx.KPOrdinal == 0 {
			k = edakit.Int(1, 2)
		}
		n := edakit.Scaled(ctx, 3, 5)
		xs := make([]float64, n)
		for i := range xs {
			xs[i] = float64(edakit.Int(0, 20))
		}
		xs[0], xs[1] = float64(edakit.Int(0, 5)), float64(edakit.Int(12, 20))
		xs = edakit.Shuffle(xs)
		g := edakit.Pick([]float64{1, 2, 4, 5})
		w := wa(xs, g)
		hp := maxF(xs) - minF(xs)
		stem := edakit.Sprintf("A DREAMPlace-style placer models a net's x-extent with the weighted-average (WA) wirelength $\\frac{\\sum x_i e^{x_i/\\gamma}}{\\sum e^{x_i/\\gamma}} - \\frac{\\sum x_i e^{-x_i/\\gamma}}{\\sum e^{-x_i/\\gamma}}$. The net's pins sit at x = %s um and $\\gamma = %s$ um.", floats(xs, 0), f(g, 0))
		if k == 1 {
			ans := edakit.Num(w, 4)
			return edakit.Numeric(stem+" What is the WA wirelength in x, in um? Round to 4 decimal places.",
				ans, edakit.Sprintf("The first term is a soft maximum of the coordinates and the second a soft minimum. Their difference is %s um, close to the HPWL span of %s um.", ans, f(hp, 0)),
				map[string]string{"hpwl_x": f(hp, 0)})
		}
		ans := edakit.Num(hp-w, 4)
		return edakit.Numeric(stem+" By how much does the WA wirelength underestimate the exact x-span (max minus min), in um? Round to 4 decimal places.",
			ans, edakit.Sprintf("Exact span $= %s$ um, WA $= %s$ um, gap $= %s$ um. The gap shrinks as $\\gamma \\to 0$, but a small $\\gamma$ makes gradients nearly as sparse as HPWL's, so placers anneal $\\gamma$ downward.", f(hp, 0), f(w, 4), ans),
			map[string]string{"wa": f(w, 4)})
	default:
		x := float64(edakit.Int(10, 90))
		gw := float64(edakit.Int(-30, 30)) / 10
		gd := float64(edakit.Int(-30, 30)) / 10
		lam := edakit.Pick([]float64{0.1, 0.5, 1, 2, 4})
		eta := edakit.Pick([]float64{0.5, 1, 2})
		nx := x - eta*(gw+lam*gd)
		ans := edakit.Num(nx, 3)
		return edakit.Numeric(edakit.Sprintf("Placement as training: cell coordinates are the parameters and the loss is $\\mathrm{WL} + \\lambda D$, with smooth wirelength WL and electrostatic density penalty D. Autograd gives, for one cell at $x = %s$ um, $\\partial \\mathrm{WL}/\\partial x = %s$ and $\\partial D/\\partial x = %s$. With $\\lambda = %s$ and step size %s, what is the cell's x after one plain gradient-descent step, in um? Round to 3 decimal places.", f(x, 0), f(gw, 1), f(gd, 1), f(lam, 1), f(eta, 1)),
			ans, edakit.Sprintf("$x - \\eta(\\partial\\mathrm{WL}/\\partial x + \\lambda\\,\\partial D/\\partial x) = %s - %s(%s + %s \\cdot %s) = %s$ um. As $\\lambda$ grows over iterations, the density term pushes cells apart.", f(x, 0), f(eta, 1), f(gw, 1), f(lam, 1), pn(gd, 1), ans),
			map[string]string{"grad": f(gw+lam*gd, 3)})
	}
}

func minF(xs []float64) float64 {
	m := xs[0]
	for _, x := range xs {
		m = math.Min(m, x)
	}
	return m
}

// lse returns gamma * log sum exp(a_i / gamma), computed stably.
func lse(a []float64, g float64) float64 {
	m := maxF(a)
	s := 0.0
	for _, x := range a {
		s += math.Exp((x - m) / g)
	}
	return m + g*math.Log(s)
}

// softw returns the softmax weight of a[i] at temperature g.
func softw(a []float64, i int, g float64) float64 {
	m := maxF(a)
	s := 0.0
	for _, x := range a {
		s += math.Exp((x - m) / g)
	}
	return math.Exp((a[i]-m)/g) / s
}

// smoothMax: log-sum-exp value (KP 1), its gap to the true max (KP 2), and the gradient
// weight of one input (KP 3).
func smoothMax(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 2, 4)
	a := make([]float64, n)
	g := edakit.Pick([]float64{1, 2, 5, 10})
	base := edakit.Int(100, 130)
	for i := range a {
		a[i] = float64(base + edakit.Int(-int(2*g), int(2*g)))
	}
	v := lse(a, g)
	stem := edakit.Sprintf("A differentiable timer replaces the max at a gate input by $\\mathrm{LSE}_\\gamma(a) = \\gamma \\ln \\sum_i e^{a_i/\\gamma}$. The candidate arrival times are %s ps and $\\gamma = %s$ ps.", floats(a, 0), f(g, 0))
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(v, 4)
		return edakit.Numeric(stem+" What is the smoothed arrival time, in ps? Round to 4 decimal places.",
			ans, edakit.Sprintf("Factor out the max %s: $%s + %s\\ln\\sum e^{(a_i - %s)/%s} = %s$ ps.", f(maxF(a), 0), f(maxF(a), 0), f(g, 0), f(maxF(a), 0), f(g, 0), ans),
			map[string]string{"max": f(maxF(a), 0)})
	case 2:
		gap := v - maxF(a)
		ans := edakit.Num(gap, 4)
		return edakit.Numeric(stem+" By how many ps does the smoothed value exceed the true maximum? Round to 4 decimal places.",
			ans, edakit.Sprintf("LSE $= %s$, max $= %s$, gap $= %s$ ps. The gap is always between 0 and $\\gamma\\ln n = %s$ ps: the smooth max overestimates, so a smoothed timer is pessimistic.", f(v, 4), f(maxF(a), 0), ans, f(g*math.Log(float64(n)), 4)),
			map[string]string{"bound": f(g*math.Log(float64(n)), 4)})
	default:
		i := edakit.Int(0, n-1)
		w := softw(a, i, g)
		ans := edakit.Num(w, 4)
		return edakit.Numeric(stem+edakit.Sprintf(" What is $\\partial\\,\\mathrm{LSE}/\\partial a_{%d}$, the share of the gradient that flows back to input %d (arrival %s ps)? Round to 4 decimal places.", i+1, i+1, f(a[i], 0)),
			ans, edakit.Sprintf("The gradient is the softmax weight $e^{a_i/\\gamma}/\\sum_j e^{a_j/\\gamma} = %s$. With a hard max it would be 1 for the latest input and 0 for the rest.", ans),
			map[string]string{"weight": ans})
	}
}

// differentiableTiming: smoothed gate arrival (KP 1), its derivative with respect to one arc
// delay (KP 2), and the TNS gradient through two endpoints (KP 3).
func differentiableTiming(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(2, 3)
	a := make([]float64, n)
	d := make([]float64, n)
	s := make([]float64, n)
	var desc []string
	g := edakit.Pick([]float64{2, 4, 5, 10})
	base := edakit.Int(120, 150)
	for i := range a {
		d[i] = float64(edakit.Int(15, 40))
		a[i] = float64(base+edakit.Int(-int(2*g), int(2*g))) - d[i]
		s[i] = a[i] + d[i]
		desc = append(desc, edakit.Sprintf("input %d arrives at %s ps through an arc of delay $d_%d = %s$ ps", i+1, f(a[i], 0), i+1, f(d[i], 0)))
	}
	A := lse(s, g)
	stem := edakit.Sprintf("In a differentiable STA, gate G's output arrival is $A = \\mathrm{LSE}_\\gamma(a_i + d_i)$ with $\\gamma = %s$ ps, where %s.", f(g, 0), strings.Join(desc, "; "))
	w := softw(s, 0, g)
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(A, 4)
		return edakit.Numeric(stem+" What is $A$, in ps? Round to 4 decimal places.",
			ans, edakit.Sprintf("Arc sums: %s. $A = %s$ ps.", floats(s, 0), ans), map[string]string{"max": f(maxF(s), 0)})
	case 2:
		ans := edakit.Num(w, 4)
		return edakit.Numeric(stem+" What is $\\partial A/\\partial d_1$? Round to 4 decimal places.",
			ans, edakit.Sprintf("By the chain rule, $\\partial A/\\partial d_1 = \\partial\\,\\mathrm{LSE}/\\partial(a_1 + d_1) = e^{(a_1+d_1)/\\gamma}/\\sum_j e^{(a_j+d_j)/\\gamma} = %s$.", ans),
			map[string]string{"A": f(A, 4)})
	default:
		e1 := float64(edakit.Int(10, 40))
		e2 := float64(edakit.Int(10, 40))
		t1 := math.Round(A + e1 + float64(edakit.Int(-20, 15)))
		t2 := math.Round(A + e2 + float64(edakit.Int(-20, 15)))
		s1 := t1 - (A + e1)
		s2 := t2 - (A + e2)
		cnt := 0
		if s1 < 0 {
			cnt++
		}
		if s2 < 0 {
			cnt++
		}
		gr := float64(cnt) * w
		ans := edakit.Num(gr, 4)
		return edakit.Numeric(stem+edakit.Sprintf(" G's output reaches endpoint 1 after %s ps more and endpoint 2 after %s ps more, with required times %s ps and %s ps. The loss is $\\mathrm{TNS} = \\sum_j \\max(0, -\\mathrm{slack}_j)$. What is $\\partial\\,\\mathrm{TNS}/\\partial d_1$? Round to 4 decimal places.", f(e1, 0), f(e2, 0), f(t1, 0), f(t2, 0)),
			ans, edakit.Sprintf("$A = %s$. Slacks: $%s - (A + %s) = %s$ and $%s - (A + %s) = %s$ ps. Each violating endpoint contributes $\\partial A/\\partial d_1 = %s$; passing endpoints contribute 0. With %s, the gradient is %s.", f(A, 4), f(t1, 0), f(e1, 0), f(s1, 4), f(t2, 0), f(e2, 0), f(s2, 4), f(w, 4), map[int]string{0: "no endpoint violating", 1: "one endpoint violating", 2: "both endpoints violating"}[cnt], ans),
			map[string]string{"violating": edakit.IntStr(cnt)})
	}
}

var relaxationItems = []edakit.Item{
	{KP: 1, Q: "A differentiable placer converges to a low-loss solution. Why must a legalizer still run?",
		A: "The continuous solution ignores row, site and overlap constraints, which are discrete", D: []string{"Gradients always point toward overlaps", "The loss function includes legality exactly", "Legalization only affects runtime"},
		Why: "Density penalties discourage overlap on average but cannot enforce site-aligned, overlap-free positions. That is a discrete problem."},
	{KP: 1, Q: "Which placement decision is hardest to express as a differentiable loss?",
		A: "Choosing which discrete orientation or which row a cell takes", D: []string{"Moving a cell by a small distance", "Reducing a smooth wirelength", "Spreading cells by a density field"},
		Why: "Orientation and row choice are categorical. Gradient methods need continuous relaxations whose rounding can change the result."},
	{KP: 1, Q: "Why does a differentiable placer use WA or log-sum-exp wirelength rather than HPWL?",
		A: "HPWL is piecewise linear with gradients only at the extreme pins, so it gives sparse, non-smooth updates", D: []string{"HPWL cannot be computed on a GPU", "HPWL ignores two-pin nets", "WA is always larger than HPWL"},
		Why: "Smooth surrogates give every pin a gradient and make the objective friendly to first-order optimizers."},
	{KP: 1, Q: "Which part of global placement maps most directly onto deep-learning toolkit operators?",
		A: "Wirelength and density evaluation over all nets and bins, with gradients from autodifferentiation or hand-written backward kernels", D: []string{"Legalization into rows", "Reading LEF and DEF", "Clock tree synthesis"},
		Why: "These are large, regular, data-parallel computations, the same shape as neural-network layers, which is why GPU toolkits accelerate them."},
	{KP: 2, Q: "A differentiable timer uses log-sum-exp in place of max. What happens to the reported arrival times?",
		A: "They are slightly pessimistic, by at most gamma times the log of the number of inputs per max", D: []string{"They become exact", "They are always optimistic", "They no longer depend on delays"},
		Why: "LSE upper-bounds max. The error is controlled by gamma, which trades smoothness against accuracy."},
	{KP: 2, Q: "Why do differentiable placers start with a large smoothing parameter and shrink it?",
		A: "Large smoothing gives informative gradients early, and small smoothing makes the final objective match the true metric", D: []string{"A large value makes the result legal", "Shrinking it reduces memory", "The value is fixed by the netlist"},
		Why: "This continuation strategy follows a sequence of easier problems toward the real one."},
	{KP: 2, Q: "A gradient step reduces smoothed TNS but signoff TNS gets worse. What is a likely reason?",
		A: "The smoothed model mis-ranks critical paths or the delay model used for gradients differs from signoff", D: []string{"Gradients cannot reduce TNS", "Signoff TNS is always positive", "The step size was exactly right"},
		Why: "Gradients optimize the model being differentiated. If the model's approximations diverge from signoff, improvement does not transfer."},
	{KP: 2, Q: "Why are electrostatic density penalties used instead of counting overlaps directly?",
		A: "An overlap count is piecewise constant, so its gradient is zero almost everywhere", D: []string{"Overlap counts are too expensive to compute", "Electrostatics is required by LEF", "Overlaps do not affect routability"},
		Why: "A smooth potential field yields forces that push cells out of dense regions; a count yields no gradient."},
	{KP: 3, Q: "A differentiable placement reaches a good wirelength, but the legalizer moves many cells far from their positions. What should the tool designer look at?",
		A: "The density target and penalty schedule, since the continuous solution was too dense to legalize locally", D: []string{"The file writer", "The random seed only", "Nothing, because legalization is free"},
		Why: "If global placement leaves regions overfull, legalization displaces cells and destroys the optimized wirelength."},
	{KP: 3, Q: "What guarantee does gradient descent give on a non-convex placement objective?",
		A: "Only convergence toward a stationary point, not a global optimum", D: []string{"A globally optimal placement", "A legal placement", "The best possible timing"},
		Why: "The placement objective is non-convex. Quality depends on initialization, schedules and the downstream detailed steps."},
	{KP: 3, Q: "Why do differentiable placers keep a classical detailed placer after legalization?",
		A: "Discrete local moves such as swaps and reordering find gains that continuous gradients cannot express", D: []string{"Because GPUs cannot store placements", "To undo the global placement", "Because detailed placement replaces the density penalty"},
		Why: "Swaps and in-row reordering are combinatorial improvements on a legal placement, outside a gradient method's reach."},
	{KP: 3, Q: "In the four-layer engine model, where does a GPU-accelerated differentiable global placer belong?",
		A: "As a strategy or fast mode that feeds the core's legalizer and detailed placer", D: []string{"As a replacement for the legality checks in the core", "In the observation hooks only", "In the timing library"},
		Why: "The continuous optimizer decides approximate positions; legality and final positions stay with the core's exact steps."},
}
