package opt

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func dualFunction(ctx generator.GeneratorContext) generator.Problem {
	// p and q from {1, 2, 4} keep 1/(4p) + 1/(4q) a terminating decimal.
	p, q := edakit.Pick([]int{1, 2, 4}), edakit.Pick([]int{1, 2, 4})
	c := edakit.Int(2, 12)
	s := 1/(4*float64(p)) + 1/(4*float64(q))
	g := func(l float64) float64 { return l*float64(c) - s*l*l }
	obj := poly([]int{p, q}, []string{"x^2", "y^2"})
	stem := edakit.Sprintf("A power-minimizing sizing model solves $\\min\\ %s$ subject to a speed requirement $x + y \\ge %d$. Its Lagrangian is $L = %s + \\lambda(%d - x - y)$ with $\\lambda \\ge 0$.", obj, c, obj, c)
	expl := edakit.Sprintf("Minimizing L over x and y gives $x = \\lambda/%d$, $y = \\lambda/%d$, so $g(\\lambda) = %d\\lambda - (\\tfrac{1}{%d} + \\tfrac{1}{%d})\\lambda^2 = %d\\lambda - %s\\lambda^2$.", 2*p, 2*q, c, 4*p, 4*q, c, f(s, 4))
	switch kp(ctx) {
	case 1:
		l := float64(edakit.Int(1, 16)) / 2
		return num(stem+edakit.Sprintf(" What is the dual function value $g(\\lambda) = \\min_{x,y} L$ at $\\lambda = %s$?", f(l, 1))+roundNote(3), f(g(l), 3),
			expl+edakit.Sprintf(" At $\\lambda = %s$ that is %s, a lower bound on the minimum power.", f(l, 1), f(g(l), 3)),
			"g", f(g(l), 3))
	case 2:
		ls := float64(c) / (2 * s)
		gs := float64(c*c) / (4 * s)
		if edakit.Int(0, 1) == 0 {
			return num(stem+" What multiplier $\\lambda^*$ maximizes the dual function?"+roundNote(3), f(ls, 3),
				expl+edakit.Sprintf(" Setting $g'(\\lambda) = %d - 2 \\cdot %s\\lambda = 0$ gives $\\lambda^* = %s$.", c, f(s, 4), f(ls, 3)),
				"lambda_star", f(ls, 3), "g_star", f(gs, 3))
		}
		return num(stem+" What is the dual optimal value $\\max_\\lambda g(\\lambda)$?"+roundNote(3), f(gs, 3),
			expl+edakit.Sprintf(" The maximum is at $\\lambda^* = %s$ with $g^* = %s$. The problem is convex with a strictly feasible point, so this equals the primal minimum: zero duality gap.", f(ls, 3), f(gs, 3)),
			"lambda_star", f(ls, 3), "g_star", f(gs, 3))
	default:
		x := edakit.Int(c/2, c)
		y := c - x + edakit.Int(0, 2)
		fx := float64(p*x*x + q*y*y)
		l := float64(edakit.Int(1, 16)) / 2
		gap := fx - g(l)
		return num(stem+edakit.Sprintf(" A heuristic returns the feasible point $(x, y) = (%d, %d)$. Using the dual value at $\\lambda = %s$, how far, at most, can this point be from optimal?", x, y, f(l, 1))+roundNote(3), f(gap, 3),
			expl+edakit.Sprintf(" The point costs %s and $g(%s) = %s$ is a lower bound on the optimum, so the point is within %s of optimal.", f(fx, 0), f(l, 1), f(g(l), 3), f(gap, 3)),
			"primal", f(fx, 0), "g", f(g(l), 3), "gap", f(gap, 3))
	}
}

func lagrangianRelaxation(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		nin, nout := edakit.Int(1, 3), edakit.Int(2, 3)
		in := make([]float64, nin)
		sum := 0.0
		for i := range in {
			in[i] = float64(edakit.Int(1, 30)) / 10
			sum += in[i]
		}
		out := make([]float64, nout)
		rest := sum
		for i := 0; i < nout-1; i++ {
			out[i] = math.Round(rest*float64(edakit.Int(1, 6))/10*10) / 10
			if out[i] >= rest {
				out[i] = 0
			}
			rest -= out[i]
		}
		out[nout-1] = rest
		var is, os []string
		for _, v := range in {
			is = append(is, f(v, 1))
		}
		for _, v := range out[:nout-1] {
			os = append(os, f(v, 1))
		}
		return num(edakit.Sprintf("In Lagrangian-relaxation gate sizing, each timing arc carries a multiplier $\\lambda \\ge 0$, and optimality requires flow conservation at every gate: incoming multipliers sum to outgoing ones. A gate's %d incoming arcs carry %s; %d of its %d outgoing arcs carry %s. What multiplier must the last outgoing arc carry?", nin, strings.Join(is, ", "), nout-1, nout, strings.Join(os, ", "))+roundNote(1), f(out[nout-1], 1),
			edakit.Sprintf("The KKT conditions of the relaxed sizing problem force $\\sum \\lambda_{in} = \\sum \\lambda_{out}$ at each node. Incoming total %s, minus the known outgoing %s, leaves %s. That lets the solver keep only one multiplier set that lives on the flow polytope.", f(sum, 1), f(sum-out[nout-1], 1), f(out[nout-1], 1)),
			"missing", f(out[nout-1], 1))
	case 2:
		alpha := float64(edakit.Int(1, 4))
		mu := float64(edakit.Int(1, 20)) / 10
		R := float64(edakit.Int(2, 12))
		C := float64(edakit.Int(5, 40))
		x := math.Sqrt(mu * R * C / alpha)
		return num(edakit.Sprintf("In the Lagrangian subproblem for one gate, area costs $%sx$ and the gate's delay $R C_L / x$ (with $R = %s$ kOhm at unit size and $C_L = %s$ fF) is weighted by the summed multipliers $\\mu = %s$ of the arcs through it. Minimizing $%s x + %s \\cdot %s \\cdot %s / x$, what size x does the gate get?", f(alpha, 0), f(R, 0), f(C, 0), f(mu, 1), f(alpha, 0), f(mu, 1), f(R, 0), f(C, 0))+roundNote(3), f(x, 3),
			edakit.Sprintf("The term $a x + b/x$ is minimized at $x = \\sqrt{b/a} = \\sqrt{%s \\cdot %s \\cdot %s / %s} = %s$. A larger multiplier, meaning a more critical gate, buys a larger size.", f(mu, 1), f(R, 0), f(C, 0), f(alpha, 0), f(x, 3)),
			"x", f(x, 3))
	default:
		alpha := float64(edakit.Int(1, 4))
		mu := float64(edakit.Int(1, 20)) / 10
		mup := float64(edakit.Int(1, 20)) / 10
		R := float64(edakit.Int(2, 12))
		C := float64(edakit.Int(5, 40))
		Rp := float64(edakit.Int(1, 6))
		cin := float64(edakit.Int(1, 4))
		beta := mup * Rp * cin
		x := math.Sqrt(mu * R * C / (alpha + beta))
		return num(edakit.Sprintf("A gate of size x has area $%sx$ and delay $%s \\cdot %s / x$ weighted by its multipliers $\\mu = %s$. Its input capacitance $%sx$ fF also loads the upstream driver (resistance %s kOhm), whose arcs carry multipliers $\\mu_{up} = %s$. Minimizing the gate's share of the Lagrangian, what size x results?", f(alpha, 0), f(R, 0), f(C, 0), f(mu, 1), f(cin, 0), f(Rp, 0), f(mup, 1))+roundNote(3), f(x, 3),
			edakit.Sprintf("The upstream delay grows by $R_{up} c_{in} x$, weighted by $\\mu_{up}$, so the linear coefficient becomes $%s + %s \\cdot %s \\cdot %s = %s$. Then $x = \\sqrt{%s \\cdot %s \\cdot %s / %s} = %s$. Ignoring the upstream load would oversize the gate.", f(alpha, 0), f(mup, 1), f(Rp, 0), f(cin, 0), f(alpha+beta, 2), f(mu, 1), f(R, 0), f(C, 0), f(alpha+beta, 2), f(x, 3)),
			"lin", f(alpha+beta, 3), "x", f(x, 3))
	}
}

func subgradient(ctx generator.GeneratorContext) generator.Problem {
	T := 10 * edakit.Int(50, 120)
	switch kp(ctx) {
	case 1, 2:
		k := ctx.KPOrdinal
		if k == 0 {
			k = edakit.Int(1, 2)
		}
		lam := float64(edakit.Int(0, 30)) / 10
		t := float64(edakit.Int(1, 10)) / 1000
		var a int
		if k == 1 {
			a = T + 10*edakit.Int(1, 30)
		} else {
			// negative violation large enough to cross zero
			need := int(math.Ceil(lam/t/10)) + edakit.Int(1, 10)
			a = T - 10*need
			if a < 0 {
				a = 0
			}
		}
		raw := lam + t*float64(a-T)
		nl := math.Max(0, raw)
		return num(edakit.Sprintf("A Lagrangian-relaxation sizer updates the multiplier of the timing constraint $a \\le T$ at a sink by $\\lambda \\leftarrow \\max(0, \\lambda + t(a - T))$. Now $\\lambda = %s$, step $t = %s$ per ps, the required time is T = %d ps and the current arrival is a = %d ps. What is the new $\\lambda$?", f(lam, 1), f(t, 3), T, a)+roundNote(3), f(nl, 3),
			edakit.Sprintf("The constraint value $a - T = %d$ ps is a subgradient of the dual. The raw update is $%s + %s \\cdot %d = %s$; the projection onto $\\lambda \\ge 0$ gives %s. Slack paths lose weight and violating paths gain it.", a-T, f(lam, 1), f(t, 3), a-T, f(raw, 3), f(nl, 3)),
			"raw", f(raw, 3), "lambda", f(nl, 3))
	default:
		lam := float64(edakit.Int(0, 30)) / 10
		t0 := float64(edakit.Int(1, 10)) / 1000
		v0, v1 := 10*edakit.Int(-20, 30), 10*edakit.Int(-30, 20)
		l1 := math.Max(0, lam+t0*float64(v0))
		l2 := math.Max(0, l1+t0/2*float64(v1))
		return num(edakit.Sprintf("A multiplier starts at $\\lambda_0 = %s$. The sizer uses diminishing steps $t_k = t_0/k$ with $t_0 = %s$ per ps and $k = 1, 2$. The constraint values $a - T$ observed are %d ps in iteration 1 and %d ps in iteration 2. With projection onto $\\lambda \\ge 0$, what is $\\lambda_2$?", f(lam, 1), f(t0, 3), v0, v1)+roundNote(4), f(l2, 4),
			edakit.Sprintf("$\\lambda_1 = \\max(0, %s + %s \\cdot %d) = %s$. $\\lambda_2 = \\max(0, %s + %s \\cdot %d) = %s$. Diminishing steps with divergent sum are what make subgradient ascent converge, slowly.", f(lam, 1), f(t0, 3), v0, f(l1, 4), f(l1, 4), f(t0/2, 4), v1, f(l2, 4)),
			"lambda1", f(l1, 4), "lambda2", f(l2, 4))
	}
}

func penalty(ctx generator.GeneratorContext) generator.Problem {
	c := edakit.Int(10, 60)
	a := c + edakit.Int(5, 40)
	stem := edakit.Sprintf("Wirelength pulls a cell toward x = %d um, but a fence requires $x \\le %d$. A quadratic-penalty placer minimizes $(x - %d)^2 + \\rho\\,\\max(0, x - %d)^2$.", a, c, a, c)
	switch kp(ctx) {
	case 1:
		rho := float64(edakit.Pick([]int{1, 2, 4, 5, 9, 10, 20, 50, 100}))
		x := (float64(a) + rho*float64(c)) / (1 + rho)
		return num(stem+edakit.Sprintf(" With $\\rho = %s$, where is the penalized minimizer, in um?", f(rho, 0))+roundNote(3), f(x, 3),
			edakit.Sprintf("The minimizer is beyond the fence, so set $2(x - %d) + 2\\rho(x - %d) = 0$: $x = (%d + %s \\cdot %d)/(1 + %s) = %s$ um.", a, c, a, f(rho, 0), c, f(rho, 0), f(x, 3)),
			"x", f(x, 3))
	case 2:
		rho := float64(edakit.Pick([]int{1, 2, 4, 5, 9, 10, 20, 50, 100}))
		v := float64(a-c) / (1 + rho)
		return num(stem+edakit.Sprintf(" With $\\rho = %s$, by how many um does the minimizer still violate the fence?", f(rho, 0))+roundNote(3), f(v, 3),
			edakit.Sprintf("$x(\\rho) - %d = (%d - %d)/(1 + \\rho) = %s$ um. The violation shrinks only like $1/\\rho$: a penalty method is never exactly feasible at finite $\\rho$.", c, a, c, f(v, 3)),
			"violation", f(v, 3))
	default:
		eps := edakit.Pick([]float64{0.1, 0.2, 0.25, 0.5, 1})
		rho := float64(a-c)/eps - 1
		return num(stem+edakit.Sprintf(" What is the smallest $\\rho$ that brings the violation down to %s um?", f(eps, 2))+roundNote(2), f(rho, 2),
			edakit.Sprintf("Need $(%d - %d)/(1 + \\rho) \\le %s$, so $\\rho \\ge %d/%s - 1 = %s$. Such large weights make the problem ill-conditioned, which is the motivation for augmented Lagrangians.", a, c, f(eps, 2), a-c, f(eps, 2), f(rho, 2)),
			"rho", f(rho, 2))
	}
}

func augmentedLagrangian(ctx generator.GeneratorContext) generator.Problem {
	c := edakit.Int(0, 20)
	a := c + edakit.Int(-10, 10)
	rho := float64(edakit.Int(1, 6))
	xOf := func(l float64) float64 { return (2*float64(a) - l + rho*float64(c)) / (2 + rho) }
	stem := edakit.Sprintf("A pin must sit exactly on a track at x = %d while wirelength pulls it toward %d. The augmented Lagrangian method with $\\rho = %s$ minimizes $%s^2 + \\lambda%s + \\tfrac{\\rho}{2}%s^2$ over x, then updates $\\lambda \\leftarrow \\lambda + \\rho%s$.", c, a, f(rho, 0), lin(a), lin(c), lin(c), lin(c))
	switch kp(ctx) {
	case 1:
		l := float64(edakit.Int(-20, 20)) / 2
		x := xOf(l)
		return num(stem+edakit.Sprintf(" With $\\lambda = %s$, what x minimizes the augmented Lagrangian?", f(l, 1))+roundNote(3), f(x, 3),
			edakit.Sprintf("Set the derivative to zero: $2%s + \\lambda + \\rho%s = 0$, so $x = (2a - \\lambda + \\rho c)/(2 + \\rho) = (2 \\cdot (%d) - (%s) + %s \\cdot %d)/(2 + %s) = %s$.", lin(a), lin(c), a, f(l, 1), f(rho, 0), c, f(rho, 0), f(x, 3)),
			"x", f(x, 3))
	case 2:
		l := float64(edakit.Int(-20, 20)) / 2
		x := xOf(l)
		nl := l + rho*(x-float64(c))
		return num(stem+edakit.Sprintf(" Starting from $\\lambda = %s$, do one minimization and one multiplier update. What is the new $\\lambda$?", f(l, 1))+roundNote(3), f(nl, 3),
			edakit.Sprintf("First $x = %s$, so the residual is $x - %d = %s$ and $\\lambda \\leftarrow %s + %s \\cdot (%s) = %s$. The exact multiplier is $2(a - c) = %d$.", f(x, 3), c, f(x-float64(c), 3), f(l, 1), f(rho, 0), f(x-float64(c), 3), f(nl, 3), 2*(a-c)),
			"x", f(x, 3), "lambda", f(nl, 3))
	default:
		l0 := 0.0
		x1 := xOf(l0)
		l1 := l0 + rho*(x1-float64(c))
		x2 := xOf(l1)
		l2 := l1 + rho*(x2-float64(c))
		return num(stem+" Starting from $\\lambda_0 = 0$, run two full iterations. What is $\\lambda_2$?"+roundNote(3), f(l2, 3),
			edakit.Sprintf("Iteration 1: $x_1 = %s$, $\\lambda_1 = %s$. Iteration 2: $x_2 = %s$, $\\lambda_2 = %s$. The multiplier approaches the exact value $2(a - c) = %d$ with $\\rho$ held fixed; a pure penalty method would need $\\rho \\to \\infty$.", f(x1, 3), f(l1, 3), f(x2, 3), f(l2, 3), 2*(a-c)),
			"lambda1", f(l1, 3), "lambda2", f(l2, 3))
	}
}

func admm(ctx generator.GeneratorContext) generator.Problem {
	a, b := edakit.Int(0, 40), edakit.Int(0, 40)
	rho := float64(edakit.Int(1, 4))
	z0 := float64(edakit.Int(0, 40))
	u0 := 0.0
	x1 := (2*float64(a) + rho*(z0-u0)) / (2 + rho)
	z1 := (2*float64(b) + rho*(x1+u0)) / (2 + rho)
	u1 := u0 + x1 - z1
	stem := edakit.Sprintf("Two placement partitions each own a copy of a shared boundary cell's coordinate: partition 1 prefers %d um, partition 2 prefers %d um. ADMM solves $\\min %s^2 + %s^2$ subject to $x = z$ with scaled updates $x \\leftarrow \\arg\\min_x %s^2 + \\tfrac{\\rho}{2}(x - z + u)^2$, $z \\leftarrow \\arg\\min_z %s^2 + \\tfrac{\\rho}{2}(x - z + u)^2$, $u \\leftarrow u + x - z$, with $\\rho = %s$, $z_0 = %s$, $u_0 = 0$.", a, b, lin(a), linv("z", b), lin(a), linv("z", b), f(rho, 0), f(z0, 0))
	expl := edakit.Sprintf("$x_1 = (2 \\cdot %d + \\rho(z_0 - u_0))/(2 + \\rho) = %s$. $z_1 = (2 \\cdot %d + \\rho(x_1 + u_0))/(2 + \\rho) = %s$. $u_1 = u_0 + x_1 - z_1 = %s$. Each partition solves only its own small problem; u accumulates the disagreement.", a, f(x1, 3), b, f(z1, 3), f(u1, 3))
	switch kp(ctx) {
	case 1:
		return num(stem+" What is $x_1$?"+roundNote(3), f(x1, 3), expl, "x1", f(x1, 3))
	case 2:
		return num(stem+" What is $z_1$ (computed with the new $x_1$)?"+roundNote(3), f(z1, 3), expl, "x1", f(x1, 3), "z1", f(z1, 3))
	default:
		return num(stem+" What is the scaled dual $u_1$ after the first full iteration?"+roundNote(3), f(u1, 3), expl, "x1", f(x1, 3), "z1", f(z1, 3), "u1", f(u1, 3))
	}
}

// linv renders (v - r) with a tidy sign for any variable name.
func linv(v string, r int) string {
	switch {
	case r == 0:
		return v
	case r > 0:
		return edakit.Sprintf("(%s - %d)", v, r)
	default:
		return edakit.Sprintf("(%s + %d)", v, -r)
	}
}
