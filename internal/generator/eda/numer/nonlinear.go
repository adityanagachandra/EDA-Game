package numer

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func newtonSystems(ctx generator.GeneratorContext) generator.Problem {
	a, b := float64(edakit.Int(1, 3)), float64(edakit.Int(1, 3))
	xr, yr := float64(edakit.Int(1, 3)), float64(edakit.Int(1, 3))
	p, q0 := xr*xr+a*yr, b*xr+yr*yr
	x0 := xr + float64(edakit.Pick([]int{-5, -3, 2, 4, 5}))/10
	y0 := yr + float64(edakit.Pick([]int{-4, -2, 3, 5}))/10
	for math.Abs(4*x0*y0-a*b) < 1 {
		x0 += 0.5
	}
	F := func(x, y float64) (float64, float64) { return x*x + a*y - p, b*x + y*y - q0 }
	J := [][]float64{{2 * x0, a}, {b, 2 * y0}}
	det := J[0][0]*J[1][1] - J[0][1]*J[1][0]
	f1, f2 := F(x0, y0)
	dx := solve2(J, []float64{-f1, -f2})
	x1, y1 := x0+dx[0], y0+dx[1]
	stem := sp("Solve $F_1 = x^2 + %sy - %s = 0$, $F_2 = %sx + y^2 - %s = 0$ by Newton-Raphson from $(x_0, y_0) = (%s, %s)$. ", f(a, 0), f(p, 0), f(b, 0), f(q0, 0), f(x0, 1), f(y0, 1))
	stem = strings.ReplaceAll(stem, " 1y", " y")
	stem = strings.ReplaceAll(stem, "= 1x", "= x")
	switch kpOf(ctx) {
	case 1:
		ans := f(det, 4)
		q := stem + "What is the determinant of the Jacobian at the starting point? Round to 4 decimal places."
		return numeric(q, ans, sp("$J = \\begin{pmatrix} 2x & %s \\\\ %s & 2y \\end{pmatrix}$, so $\\det J = 4x_0y_0 - %s = %s$.", f(a, 0), f(b, 0), f(a*b, 0), ans),
			map[string]string{"det": ans})
	case 2:
		v, name := x1, "x_1"
		if edakit.Int(0, 1) == 0 {
			v, name = y1, "y_1"
		}
		ans := f(v, 4)
		q := stem + sp("Take one Newton step, solving $J(x_0)\\,\\Delta = -F(x_0)$. What is $%s$? Round to 4 decimal places.", name)
		return numeric(q, ans, sp("$F(x_0) = (%s, %s)$, $\\det J = %s$, and the solve gives $\\Delta = (%s, %s)$, so $(x_1, y_1) = (%s, %s)$.", f(f1, 4), f(f2, 4), f(det, 4), f(dx[0], 4), f(dx[1], 4), f(x1, 4), f(y1, 4)),
			map[string]string{"dx": f(dx[0], 6), "dy": f(dx[1], 6)})
	default:
		r1, r2 := F(x1, y1)
		res := math.Max(math.Abs(r1), math.Abs(r2))
		ans := f(res, 4)
		q := stem + "Take one Newton step. What is the max-norm $\\|F(x_1, y_1)\\|_\\infty$ of the residual at the new point? Round to 4 decimal places."
		return numeric(q, ans, sp("The step lands at $(%s, %s)$. There $F = (%s, %s)$, so the max-norm is %s, down from %s at the start.", f(x1, 4), f(y1, 4), f(r1, 4), f(r2, 4), ans, f(math.Max(math.Abs(f1), math.Abs(f2)), 4)),
			map[string]string{"res0": f(math.Max(math.Abs(f1), math.Abs(f2)), 4)})
	}
}

func exponential(ctx generator.GeneratorContext) generator.Problem {
	vs := float64(edakit.Int(2, 10)) / 2
	r := float64(edakit.Int(1, 10)) / 2 // kohm
	is := 1e-14
	vt := 0.025
	v0 := float64(edakit.Int(55, 75)) / 100
	// Work in mA and kohm: current (V - Vs)/R is in mA.
	fv := (v0-vs)/r + is*(math.Exp(v0/vt)-1)*1000
	fp := 1/r + is/vt*math.Exp(v0/vt)*1000
	stem := sp("A source of %s V drives a %s kΩ resistor into a diode to ground with $I_s = 10^{-14}$ A and $V_T = 25$ mV. The node equation is $f(V) = (V - %s)/R + I_s(e^{V/V_T} - 1) = 0$. ", f(vs, 1), f(r, 1), f(vs, 1))
	switch kpOf(ctx) {
	case 1:
		ans := f(fv, 4)
		q := stem + sp("At the guess $V_0 = %s$ V, what is the residual $f(V_0)$ in mA? Round to 4 decimal places.", f(v0, 2))
		return numeric(q, ans, sp("Resistor term $(%s - %s)/%s = %s$ mA; diode term $10^{-14}(e^{%s/0.025} - 1)$ A $= %s$ mA. Sum %s mA.", f(v0, 2), f(vs, 1), f(r, 1), f((v0-vs)/r, 4), f(v0, 2), f(is*(math.Exp(v0/vt)-1)*1000, 4), ans),
			map[string]string{"f": ans})
	case 2:
		v1 := v0 - fv/fp
		ans := f(v1, 4)
		q := stem + sp("Take one Newton step from $V_0 = %s$ V. What is $V_1$, in volts? Round to 4 decimal places.", f(v0, 2))
		return numeric(q, ans, sp("$f(V_0) = %s$ mA and $f'(V_0) = 1/R + (I_s/V_T)e^{V_0/V_T} = %s$ mS, so $V_1 = V_0 - f/f' = %s$ V.", f(fv, 4), f(fp, 4), ans),
			map[string]string{"f": f(fv, 6), "fprime": f(fp, 6)})
	default:
		if edakit.Int(0, 1) == 0 {
			dv := float64(edakit.Int(1, 12)) * 10
			ratio := math.Exp(dv / 25)
			ans := f(ratio, 2)
			q := sp("A diode follows $I = I_s(e^{V/V_T} - 1)$ with $V_T = 25$ mV, well into forward bias. A Newton iterate overshoots the solution by %s mV. By what factor is the diode current at the iterate larger than at the solution? Round to 2 decimal places.", f(dv, 0))
			return numeric(q, ans, sp("Far into forward bias the $-1$ is negligible, so the ratio is $e^{%s/25} = %s$.", f(dv, 0), ans), nil)
		}
		k := edakit.Int(1, 6)
		dv := 25 * math.Log(math.Pow(10, float64(k)))
		ans := f(dv, 1)
		q := sp("A diode follows $I = I_s(e^{V/V_T} - 1)$ with $V_T = 25$ mV, well into forward bias. How large a voltage change, in mV, multiplies the current by $10^{%d}$? Round to 1 decimal place.", k)
		return numeric(q, ans, sp("$\\Delta V = V_T \\ln 10^{%d} = 25 \\cdot %s = %s$ mV.", k, f(math.Log(math.Pow(10, float64(k))), 4), ans), nil)
	}
}

func damping(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		c := float64(edakit.Int(2, 30))
		x0 := float64(edakit.Int(1, 8)) / 2
		lam := edakit.Pick([]float64{0.25, 0.5, 0.75})
		dx := -(x0*x0*x0 - c) / (3 * x0 * x0)
		x1 := x0 + lam*dx
		ans := f(x1, 4)
		q := sp("Damped Newton solves $x^3 - %s = 0$ from $x_0 = %s$ with damping factor $\\lambda = %s$, so $x_1 = x_0 + \\lambda\\,\\Delta x$ with the full Newton step $\\Delta x$. What is $x_1$? Round to 4 decimal places.", f(c, 0), f(x0, 1), f(lam, 2))
		return numeric(q, ans, sp("$\\Delta x = -(x_0^3 - %s)/(3x_0^2) = %s$, so $x_1 = %s + %s \\cdot %s = %s$.", f(c, 0), f(dx, 4), f(x0, 1), f(lam, 2), f(dx, 4), ans),
			map[string]string{"dx": f(dx, 6)})
	case 2:
		x0 := float64(edakit.Int(15, 40)) / 10
		dx := -math.Atan(x0) * (1 + x0*x0)
		lam := 1.0
		for i := 0; i < 30 && math.Abs(math.Atan(x0+lam*dx)) >= math.Abs(math.Atan(x0)); i++ {
			lam /= 2
		}
		x1 := x0 + lam*dx
		ans := f(x1, 4)
		q := sp("Newton's method on $f(x) = \\arctan x$ from $x_0 = %s$ overshoots. A backtracking line search tries $\\lambda = 1, 1/2, 1/4, \\dots$ and accepts the first $\\lambda$ with $|f(x_0 + \\lambda\\Delta x)| < |f(x_0)|$. What is the accepted iterate $x_1$? Round to 4 decimal places.", f(x0, 1))
		return numeric(q, ans, sp("$\\Delta x = -\\arctan(x_0)(1 + x_0^2) = %s$. The first accepted factor is $\\lambda = %s$, giving $x_1 = %s$.", f(dx, 4), f(lam, 6), ans),
			map[string]string{"lambda": f(lam, 6)})
	default:
		n := edakit.Int(2, 4)
		x0 := make([]float64, n)
		dx := make([]float64, n)
		var parts []string
		mx := 0.0
		for i := range x0 {
			x0[i] = float64(edakit.Int(0, 120)) / 100
			dx[i] = float64(edakit.Int(-300, 300)) / 100
			if dx[i] == 0 {
				dx[i] = 0.4
			}
			mx = math.Max(mx, math.Abs(dx[i]))
			parts = append(parts, sp("node %d: $v = %s$ V, $\\Delta v = %s$ V", i+1, f(x0[i], 2), f(dx[i], 2)))
		}
		lim := edakit.Pick([]float64{0.1, 0.2, 0.3, 0.5})
		lam := math.Min(1, lim/mx)
		j := edakit.Int(0, n-1)
		ans := f(x0[j]+lam*dx[j], 4)
		q := "A Newton iteration proposes these updates: " + strings.Join(parts, "; ") + sp(". The simulator scales the whole step so that no node voltage changes by more than %s V. What is the new voltage at node %d, in volts? Round to 4 decimal places.", f(lim, 1), j+1)
		return numeric(q, ans, sp("The largest proposed change is %s V, so $\\lambda = \\min(1, %s/%s) = %s$. Node %d moves to $%s + %s \\cdot %s = %s$ V.", f(mx, 2), f(lim, 1), f(mx, 2), f(lam, 4), j+1, f(x0[j], 2), f(lam, 4), f(dx[j], 2), ans),
			map[string]string{"lambda": f(lam, 6)})
	}
}

func convergenceTest(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		vold := float64(edakit.Int(100, 3300)) / 1000
		vnew := vold + float64(edakit.Int(-300, 300))/100000
		reltol, vntol := 1e-3, 1e-6
		tol := reltol*math.Max(math.Abs(vold), math.Abs(vnew)) + vntol
		ans := f(tol*1e6, 2)
		q := sp("A SPICE-style test accepts a node update when $|v_{k+1} - v_k| \\le \\text{reltol} \\cdot \\max(|v_{k+1}|, |v_k|) + \\text{vntol}$ with reltol = 0.001 and vntol = 1 µV. Two iterates are $v_k = %s$ V and $v_{k+1} = %s$ V. What is the allowed change, in µV? Round to 2 decimal places.", f(vold, 5), f(vnew, 5))
		return numeric(q, ans, sp("$0.001 \\cdot %s$ V $+ 1$ µV $= %s$ µV. The actual change is %s µV.", f(math.Max(math.Abs(vold), math.Abs(vnew)), 5), ans, f(math.Abs(vnew-vold)*1e6, 2)),
			map[string]string{"tol_uV": ans})
	case 2:
		n := edakit.Int(4, 6)
		var parts []string
		pass := 0
		for i := 0; i < n; i++ {
			imax := float64(edakit.Int(1, 500)) / 10 // uA
			tol := 1e-3*imax + 1e-6                  // uA (abstol 1 pA = 1e-6 uA)
			var res float64
			if edakit.Int(0, 1) == 0 {
				res = tol * float64(edakit.Int(1, 9)) / 10
			} else {
				res = tol * float64(edakit.Int(12, 40)) / 10
			}
			if res <= tol {
				pass++
			}
			parts = append(parts, sp("node %d: residual %s nA, largest branch current %s µA", i+1, f(res*1000, 2), f(imax, 1)))
		}
		q := "A KCL residual test accepts a node when $|\\sum i| \\le \\text{reltol} \\cdot \\max|i_{branch}| + \\text{abstol}$ with reltol = 0.001 and abstol = 1 pA. After an iteration: " + strings.Join(parts, "; ") + ". How many nodes pass?"
		return numeric(q, itoa(pass), sp("Each node's tolerance is 0.001 times its largest branch current (1 nA per µA) plus 1 pA. Comparing residuals, %d pass.", pass),
			map[string]string{"pass": itoa(pass)})
	default:
		p := edakit.Pick([]float64{1, 1.5, 2, 2})
		c := float64(edakit.Int(5, 30)) / 10
		m0 := float64(edakit.Int(11, 50)) / 10
		e0e := -edakit.Int(1, 2)
		if p > 1 {
			e0e = -edakit.Int(2, 3)
		}
		s0, e0 := sci(m0, e0e)
		e1 := c * math.Pow(e0, p)
		if p == 1 {
			e1 = c / 10 * e0
		}
		e2 := c * math.Pow(e1, p)
		if p == 1 {
			e2 = c / 10 * e1
		}
		ex1 := int(math.Floor(math.Log10(e1)))
		s1, e1v := sci(e1/math.Pow(10, float64(ex1)), ex1)
		ex2 := int(math.Floor(math.Log10(e2)))
		s2, e2v := sci(e2/math.Pow(10, float64(ex2)), ex2)
		order := math.Log(e2v/e1v) / math.Log(e1v/e0)
		ans := f(order, 2)
		q := "Three successive Newton errors on a DC operating point are " + s0 + ", " + s1 + " and " + s2 + " V. Estimate the order of convergence $p \\approx \\ln(e_2/e_1)/\\ln(e_1/e_0)$. Round to 2 decimal places."
		return numeric(q, ans, sp("$\\ln(e_2/e_1) = %s$ and $\\ln(e_1/e_0) = %s$, ratio %s. Near 2 means quadratic (Newton with an exact Jacobian); near 1 means linear (a stale or wrong Jacobian).", f(math.Log(e2v/e1v), 4), f(math.Log(e1v/e0), 4), ans),
			map[string]string{"p": ans})
	}
}

var continuationItems = []edakit.Item{
	{KP: 1, Q: "What makes a homotopy $H(x, s)$ useful for a DC operating point that plain Newton cannot find?",
		A: "H(x, 0) = 0 has a known or easy solution and H(x, 1) = F(x), so the solution can be tracked from s = 0 to s = 1", D: []string{"H is linear in x for every s", "H has the same Jacobian as F at every s", "H removes all nonlinear devices from the circuit"},
		Why: "Continuation embeds the hard problem in a family that starts at an easy one; each solve uses the previous solution as its initial guess."},
	{KP: 1, Q: "During continuation, a step from s = 0.4 to s = 0.6 fails to converge in Newton. What should the stepper do?",
		A: "Reject the step and retry with a smaller increment, such as s = 0.5", D: []string{"Accept the last iterate and continue to s = 0.8", "Restart from s = 0", "Switch off the nonlinear devices"},
		Why: "Newton converges when the previous solution is close enough to the new one, and a smaller parameter step makes it closer."},
	{KP: 1, Q: "Why does continuation use the solution at the previous parameter value as the Newton starting guess?",
		A: "Newton converges only locally, and the solution path moves continuously with s, so the previous point lies inside the basin of attraction", D: []string{"It guarantees quadratic convergence from any point", "It removes the need to factor the Jacobian", "It makes the Jacobian symmetric"},
		Why: "Newton's quadratic convergence is local. Continuation keeps each new problem close to one already solved."},
	{KP: 1, Q: "A continuation run reaches s = 1 with every step converging in 2 or 3 Newton iterations. What is a sensible step-size policy?",
		A: "Grow the step after quick convergence and shrink it after slow convergence or failure", D: []string{"Keep the step fixed at its first value", "Always halve the step after each success", "Choose steps at random to avoid bias"},
		Why: "Adaptive steppers spend effort only where the solution path bends sharply, judged by Newton iteration counts."},
	{KP: 2, Q: "In source stepping, what is the homotopy parameter?",
		A: "A factor s that scales every independent voltage and current source from 0 to full value", D: []string{"The conductance added from every node to ground", "The timestep of a pseudo-transient run", "The reltol of the Newton test"},
		Why: "With all sources at zero the circuit sits at the trivial solution; ramping s brings the sources up gradually."},
	{KP: 2, Q: "In gmin stepping, what is added to the circuit and how does it change?",
		A: "A conductance from every node to ground, started large and reduced step by step toward zero", D: []string{"A capacitor at every node, increased each step", "A series resistance in every source, increased each step", "A current source at every node, reduced to zero"},
		Why: "A large gmin makes the Jacobian strongly diagonally dominant and the problem nearly linear; shrinking it walks back to the real circuit."},
	{KP: 2, Q: "Pseudo-transient continuation adds a capacitor from every node to ground and runs a transient to steady state. Why does that help DC convergence?",
		A: "Backward Euler on the added capacitors adds C/h to the Jacobian diagonal, so each step is a well-conditioned, small move", D: []string{"Capacitors remove the nonlinearity of the diodes", "The transient is linear, so no Newton iterations are needed", "It changes the DC solution to an easier one"},
		Why: "The artificial dynamics damp each Newton problem, and at steady state the capacitor currents vanish, so the DC point is unchanged."},
	{KP: 2, Q: "A latch has two stable DC states. Source stepping from zero lands on one of them. What does that say about the result?",
		A: "The operating point is one valid solution; a different path or initial guess can land on the other", D: []string{"The other state cannot be a solution", "Source stepping always finds the lowest-power state", "The circuit has no DC solution"},
		Why: "Bistable circuits have multiple solutions, and continuation follows whichever branch its path leads to, so tools offer nodesets and initial conditions."},
	{KP: 3, Q: "A natural-parameter continuation fails at s = 0.73: the solution curve folds back, and no solution exists just beyond s = 0.73 on that branch. What fixes it?",
		A: "Pseudo-arclength continuation, which parameterizes the path by arclength so it can turn around the fold", D: []string{"Taking a larger step in s past the fold", "Tightening reltol", "Switching from backward Euler to trapezoidal"},
		Why: "At a turning point ds/dx changes sign, so stepping s monotonically cannot follow the curve; arclength parameterization can."},
	{KP: 3, Q: "Near a turning point of the continuation path, what happens to the Jacobian $\\partial H/\\partial x$?",
		A: "It becomes singular, so Newton in x alone stops converging", D: []string{"It becomes the identity", "It becomes strictly diagonally dominant", "It stays well conditioned"},
		Why: "At a fold the path is tangent to constant s, which means the Jacobian with respect to x has a null vector there."},
	{KP: 3, Q: "In pseudo-arclength continuation, what extra equation is appended to H(x, s) = 0?",
		A: "A condition that the new point lies a fixed distance along the tangent of the path", D: []string{"A condition that s increases by a fixed amount", "A condition that the residual is zero at s = 1", "A condition that all node voltages are positive"},
		Why: "The appended equation fixes progress along the curve rather than in s, which keeps the augmented Jacobian nonsingular at folds."},
	{KP: 3, Q: "A gmin stepping run converges at gmin = 1e-6 S but fails at the next step to 1e-7 S. Which response is most reasonable?",
		A: "Retry with a smaller factor, such as 1e-6 to 3e-7, before giving up", D: []string{"Report the gmin = 1e-6 solution as the final answer", "Jump straight to gmin = 0", "Raise reltol to 0.1"},
		Why: "A failed continuation step usually means the step was too aggressive; the 1e-6 solution includes artificial conductance and is not the circuit's answer."},
}
