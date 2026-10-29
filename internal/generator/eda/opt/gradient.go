package opt

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func gradientDescent(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		k := ctx.KPOrdinal
		if k == 0 {
			k = edakit.Int(1, 2)
		}
		a, b := edakit.Int(1, 6), edakit.Int(1, 6)
		h := edakit.Pick([]int{-2, -1, 1, 2})
		u, v := edakit.Int(-10, 10), edakit.Int(-10, 10)
		x0, y0 := edakit.Int(-5, 5), edakit.Int(-5, 5)
		gx := a*x0 + h*y0 - u
		gy := h*x0 + b*y0 - v
		stem := edakit.Sprintf("A quadratic wirelength model for one cell is $f(x, y) = \\tfrac{1}{2}(%s) %s$. The cell sits at $(%d, %d)$.", poly([]int{a, 2 * h, b}, []string{"x^2", "xy", "y^2"}), signed(poly([]int{-u, -v}, []string{"x", "y"})), x0, y0)
		comp := edakit.Pick([]string{"x", "y"})
		g, p0 := gx, x0
		if comp == "y" {
			g, p0 = gy, y0
		}
		if k == 1 {
			return num(stem+edakit.Sprintf(" What is the %s-component of the gradient there?", comp)+roundNote(0), edakit.IntStr(g),
				edakit.Sprintf("$\\partial f/\\partial x = %s$ and $\\partial f/\\partial y = %s$. At the cell this is $(%d, %d)$.", poly([]int{a, h, -u}, []string{"x", "y", ""}), poly([]int{h, b, -v}, []string{"x", "y", ""}), gx, gy),
				"gx", edakit.IntStr(gx), "gy", edakit.IntStr(gy))
		}
		eta := float64(edakit.Int(1, 4)) / 20
		x1 := float64(p0) - eta*float64(g)
		return num(stem+edakit.Sprintf(" One gradient-descent step with step size $\\eta = %s$ moves it. What is the new %s-coordinate?", f(eta, 2), comp)+roundNote(2), f(x1, 2),
			edakit.Sprintf("The gradient at the cell is $(%d, %d)$. The update is $%s_1 = %d - %s \\cdot (%d) = %s$.", gx, gy, comp, p0, f(eta, 2), g, f(x1, 2)),
			"gx", edakit.IntStr(gx), "gy", edakit.IntStr(gy), "next", f(x1, 3))
	default:
		a := edakit.Int(1, 8)
		c := edakit.Int(-20, 20)
		x0 := c + edakit.Int(-30, 30)
		eta := float64(edakit.Int(1, 9)) / 20
		for eta*float64(a) >= 1.9 {
			eta = float64(edakit.Int(1, 9)) / 20
		}
		n := edakit.Scaled(ctx, 2, 5)
		r := 1 - eta*float64(a)
		xn := float64(c) + math.Pow(r, float64(n))*float64(x0-c)
		return num(edakit.Sprintf("Gradient descent minimizes $f(x) = \\tfrac{%d}{2}%s^2$ from $x_0 = %d$ with step size $\\eta = %s$. What is $x_%d$?", a, lin(c), x0, f(eta, 2), n)+roundNote(3), f(xn, 3),
			edakit.Sprintf("Each step multiplies the error $%s$ by $1 - \\eta a = %s$. After %d steps $x_%d = %d + (%s)^%d \\cdot (%d) = %s$. A factor with magnitude above 1 would diverge.", strings.Trim(lin(c), "()"), f(r, 3), n, n, c, f(r, 3), n, x0-c, f(xn, 3)),
			"factor", f(r, 3), "xn", f(xn, 3))
	}
}

func eig2(a, h, b float64) (float64, float64) {
	m := (a + b) / 2
	d := math.Sqrt((a-b)*(a-b)/4 + h*h)
	return m - d, m + d
}

func stepSize(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a, b := edakit.Int(2, 12), edakit.Int(2, 12)
		h := edakit.Int(-3, 3)
		lo, hi := eig2(float64(a), float64(h), float64(b))
		for lo <= 0 {
			a++
			lo, hi = eig2(float64(a), float64(h), float64(b))
		}
		eta := 2 / hi
		return num(edakit.Sprintf("A placer's quadratic objective has constant Hessian $\\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}$. Fixed-step gradient descent converges only if $\\eta < 2/L$, where L is the largest eigenvalue. What is $2/L$?", a, h, h, b)+roundNote(3), f(eta, 3),
			edakit.Sprintf("The eigenvalues are $\\frac{a + b}{2} \\pm \\sqrt{(\\frac{a - b}{2})^2 + h^2}$ = %s and %s. With $L = %s$, the limit is $2/L = %s$. Along the stiffest direction the error factor is $1 - \\eta L$, which leaves $(-1, 1)$ beyond that step.", f(lo, 3), f(hi, 3), f(hi, 3), f(eta, 3)),
			"l_max", f(hi, 3), "eta_max", f(eta, 3))
	case 2:
		mu := edakit.Int(1, 5)
		L := mu + edakit.Int(2, 40)
		if edakit.Int(0, 1) == 0 {
			eta := 2 / float64(L+mu)
			return num(edakit.Sprintf("A strongly convex quadratic has Hessian eigenvalues between $\\mu = %d$ and $L = %d$. What fixed step size minimizes the worst-case contraction of gradient descent?", mu, L)+roundNote(4), f(eta, 4),
				edakit.Sprintf("The error factors are $1 - \\eta\\lambda$ for every eigenvalue. Balancing the extremes, $1 - \\eta\\mu = -(1 - \\eta L)$, gives $\\eta = 2/(\\mu + L) = 2/%d = %s$.", L+mu, f(eta, 4)),
				"eta", f(eta, 4))
		}
		rho := float64(L-mu) / float64(L+mu)
		return num(edakit.Sprintf("A strongly convex quadratic has Hessian eigenvalues between $\\mu = %d$ and $L = %d$. With the best fixed step, by what factor does the worst-case error shrink per iteration?", mu, L)+roundNote(4), f(rho, 4),
			edakit.Sprintf("With $\\eta = 2/(\\mu + L)$ the factor is $(L - \\mu)/(L + \\mu) = (\\kappa - 1)/(\\kappa + 1)$ with $\\kappa = %s$: %d/%d = %s. Ill-conditioned placement systems have $\\kappa$ in the thousands, so plain descent crawls.", f(float64(L)/float64(mu), 2), L-mu, L+mu, f(rho, 4)),
			"rho", f(rho, 4))
	default:
		a := edakit.Int(2, 20)
		x0 := edakit.Int(1, 9)
		c := float64(edakit.Int(1, 5)) / 10
		eta0 := 1.0
		eta := eta0
		halv := 0
		fx := float64(a) / 2 * float64(x0*x0)
		g := float64(a * x0)
		for {
			xn := float64(x0) - eta*g
			if float64(a)/2*xn*xn <= fx-c*eta*g*g+1e-12 {
				break
			}
			eta /= 2
			halv++
		}
		return num(edakit.Sprintf("Backtracking line search on $f(x) = \\tfrac{%d}{2}x^2$ at $x = %d$ starts with $\\eta = 1$ and halves $\\eta$ until the Armijo condition $f(x - \\eta g) \\le f(x) - c\\,\\eta\\, g^2$ holds, with $c = %s$. How many halvings happen before a step is accepted?", a, x0, f(c, 1))+roundNote(0),
			edakit.IntStr(halv),
			edakit.Sprintf("For this quadratic the condition reduces to $\\eta \\le 2(1 - c)/a = %s$. Halving from 1 first satisfies it at $\\eta = %s$, after %d halvings.", f(2*(1-c)/float64(a), 4), f(eta, 4), halv),
			"halvings", edakit.IntStr(halv), "eta", f(eta, 4))
	}
}

func nesterov(ctx generator.GeneratorContext) generator.Problem {
	a := edakit.Int(1, 4)
	c := edakit.Int(-10, 10)
	eta := edakit.Pick([]float64{0.05, 0.1, 0.15, 0.2, 0.25})
	for eta*float64(a) > 1 {
		eta = edakit.Pick([]float64{0.05, 0.1, 0.15, 0.2, 0.25})
	}
	beta := float64(edakit.Int(3, 9)) / 10
	grad := func(x float64) float64 { return float64(a) * (x - float64(c)) }
	stem := edakit.Sprintf("An ePlace-style placer runs Nesterov's method on $f(x) = \\tfrac{%d}{2}%s^2$ with step $\\eta = %s$ and momentum $\\beta = %s$: $y_k = x_k + \\beta(x_k - x_{k-1})$, then $x_{k+1} = y_k - \\eta f'(y_k)$.", a, lin(c), f(eta, 4), f(beta, 1))
	switch kp(ctx) {
	case 1:
		xk, xp := float64(edakit.Int(-20, 20)), float64(edakit.Int(-20, 20))
		y := xk + beta*(xk-xp)
		return num(stem+edakit.Sprintf(" With $x_k = %s$ and $x_{k-1} = %s$, what is the look-ahead point $y_k$?", f(xk, 0), f(xp, 0))+roundNote(2), f(y, 2),
			edakit.Sprintf("$y_k = %s + %s(%s - (%s)) = %s$. The gradient is then evaluated at this extrapolated point, not at $x_k$; that is the difference from heavy-ball momentum.", f(xk, 0), f(beta, 1), f(xk, 0), f(xp, 0), f(y, 2)),
			"y", f(y, 3))
	case 2:
		xk, xp := float64(edakit.Int(-20, 20)), float64(edakit.Int(-20, 20))
		y := xk + beta*(xk-xp)
		x1 := y - eta*grad(y)
		return num(stem+edakit.Sprintf(" With $x_k = %s$ and $x_{k-1} = %s$, what is $x_{k+1}$?", f(xk, 0), f(xp, 0))+roundNote(3), f(x1, 3),
			edakit.Sprintf("Look ahead: $y_k = %s$. Gradient there: $%d(%s - (%d)) = %s$. Step: $x_{k+1} = %s - %s \\cdot (%s) = %s$.", f(y, 3), a, f(y, 3), c, f(grad(y), 3), f(y, 3), f(eta, 4), f(grad(y), 3), f(x1, 3)),
			"y", f(y, 3), "next", f(x1, 3))
	default:
		x0 := float64(c + edakit.Int(-25, 25))
		x1 := x0 - eta*grad(x0)
		y1 := x1 + beta*(x1-x0)
		x2 := y1 - eta*grad(y1)
		return num(stem+edakit.Sprintf(" Start at $x_0 = x_{-1} = %s$. What is $x_2$?", f(x0, 0))+roundNote(3), f(x2, 3),
			edakit.Sprintf("With $x_{-1} = x_0$ the first look-ahead is $y_0 = x_0$, so $x_1 = %s - %s \\cdot (%s) = %s$. Then $y_1 = %s + %s(%s - (%s)) = %s$ and $x_2 = y_1 - \\eta f'(y_1) = %s$.", f(x0, 0), f(eta, 4), f(grad(x0), 3), f(x1, 3), f(x1, 3), f(beta, 1), f(x1, 3), f(x0, 0), f(y1, 3), f(x2, 3)),
			"x1", f(x1, 3), "y1", f(y1, 3), "x2", f(x2, 3))
	}
}

func conjugateGradient(ctx generator.GeneratorContext) generator.Problem {
	p, q := edakit.Int(2, 8), edakit.Int(2, 8)
	h := edakit.Int(-1, 1)
	b1, b2 := edakit.Int(-10, 10), edakit.Int(-10, 10)
	for b1 == 0 && b2 == 0 {
		b1 = edakit.Int(1, 10)
	}
	rr := float64(b1*b1 + b2*b2)
	Ar1 := float64(p*b1 + h*b2)
	Ar2 := float64(h*b1 + q*b2)
	rAr := float64(b1)*Ar1 + float64(b2)*Ar2
	alpha := rr / rAr
	stem := edakit.Sprintf("A quadratic placer solves $Ax = b$ with $A = \\begin{pmatrix} %d & %d \\\\ %d & %d \\end{pmatrix}$ and $b = (%d, %d)$ by conjugate gradient from $x_0 = 0$, so $r_0 = d_0 = b$.", p, h, h, q, b1, b2)
	switch kp(ctx) {
	case 1:
		return num(stem+" What is the first step length $\\alpha_0 = r_0^T r_0 / d_0^T A d_0$?"+roundNote(4), f(alpha, 4),
			edakit.Sprintf("$r_0^T r_0 = %s$, $A d_0 = (%s, %s)$, $d_0^T A d_0 = %s$, so $\\alpha_0 = %s$. This is the exact line minimum along $d_0$.", f(rr, 0), f(Ar1, 0), f(Ar2, 0), f(rAr, 0), f(alpha, 4)),
			"alpha0", f(alpha, 4))
	case 2:
		comp := edakit.Int(1, 2)
		v := alpha * float64(b1)
		if comp == 2 {
			v = alpha * float64(b2)
		}
		return num(stem+edakit.Sprintf(" What is component %d of the first iterate $x_1 = x_0 + \\alpha_0 d_0$?", comp)+roundNote(3), f(v, 3),
			edakit.Sprintf("$\\alpha_0 = %s/%s = %s$, so $x_1 = %s \\cdot (%d, %d) = (%s, %s)$.", f(rr, 0), f(rAr, 0), f(alpha, 4), f(alpha, 4), b1, b2, f(alpha*float64(b1), 3), f(alpha*float64(b2), 3)),
			"alpha0", f(alpha, 4))
	default:
		det := float64(p*q - h*h)
		x1 := (float64(q*b1) - float64(h*b2)) / det
		x2 := (float64(p*b2) - float64(h*b1)) / det
		comp := edakit.Int(1, 2)
		v := x1
		if comp == 2 {
			v = x2
		}
		return num(stem+edakit.Sprintf(" In exact arithmetic, what is component %d of the second iterate $x_2$?", comp)+roundNote(3), f(v, 3),
			edakit.Sprintf("CG on an n-by-n symmetric positive definite system terminates in at most n steps, because its search directions are A-conjugate. Here n = 2, so $x_2 = A^{-1}b = (%s, %s)$.", f(x1, 3), f(x2, 3)),
			"x_star_1", f(x1, 3), "x_star_2", f(x2, 3))
	}
}

func newtonMin(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := edakit.Int(2, 9)
		x0 := float64(edakit.Int(0, 4)) / 2
		x1 := x0 - 1 + float64(a)*math.Exp(-x0)
		return num(edakit.Sprintf("Newton's method minimizes $f(x) = e^x - %dx$ from $x_0 = %s$. What is $x_1 = x_0 - f'(x_0)/f''(x_0)$?", a, f(x0, 1))+roundNote(3), f(x1, 3),
			edakit.Sprintf("$f'(x) = e^x - %d$ and $f''(x) = e^x$, so $x_1 = x_0 - 1 + %d e^{-x_0} = %s$. The true minimizer is $\\ln %d = %s$.", a, a, f(x1, 3), a, f(math.Log(float64(a)), 3)),
			"x1", f(x1, 3))
	case 2:
		p, q := edakit.Int(1, 5), edakit.Int(1, 5)
		h := edakit.Pick([]int{-2, -1, 1, 2})
		for 4*p*q <= h*h {
			p++
		}
		u, v := edakit.Int(-12, 12), edakit.Int(-12, 12)
		x0, y0 := edakit.Int(-5, 5), edakit.Int(-5, 5)
		H11, H22, H12 := float64(2*p), float64(2*q), float64(h)
		det := H11*H22 - H12*H12
		xs := (float64(u)*H22 - H12*float64(v)) / det
		ys := (H11*float64(v) - H12*float64(u)) / det
		comp := edakit.Pick([]string{"x", "y"})
		ans := xs
		if comp == "y" {
			ans = ys
		}
		return num(edakit.Sprintf("Newton's method minimizes $f(x, y) = %s$ from $(%d, %d)$. What is the %s-coordinate after one full Newton step?", poly([]int{p, h, q, -u, -v}, []string{"x^2", "xy", "y^2", "x", "y"}), x0, y0, comp)+roundNote(3), f(ans, 3),
			edakit.Sprintf("f is quadratic, so its second-order Taylor model is exact and one Newton step lands on the minimizer from any start. Solve $\\nabla f = 0$: $%s = %d$, $%s = %d$, giving $(%s, %s)$.", poly([]int{2 * p, h}, []string{"x", "y"}), u, poly([]int{h, 2 * q}, []string{"x", "y"}), v, f(xs, 3), f(ys, 3)),
			"x_star", f(xs, 3), "y_star", f(ys, 3))
	default:
		s := edakit.Int(1, 4)
		x0 := -float64(edakit.Int(1, 9)) / 2
		x1 := x0 - (x0*x0-float64(s*s))/(2*x0)
		return num(edakit.Sprintf("Pure Newton is applied to $f(x) = x^3 - %dx$ from $x_0 = %s$. What is $x_1$?", 3*s*s, f(x0, 1))+roundNote(3), f(x1, 3),
			edakit.Sprintf("$f' = 3x^2 - %d$, $f'' = 6x$, so $x_1 = x_0 - (x_0^2 - %d)/(2x_0) = %s$. Here $f''(x_0) < 0$: the quadratic model is concave, and the step heads to the local maximum at $x = -%d$. Practical codes check curvature and fall back to a descent direction.", 3*s*s, s*s, f(x1, 3), s),
			"x1", f(x1, 3))
	}
}

func quasiNewton(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		for {
			x0, x1 := edakit.Int(0, 5), edakit.Int(6, 12)
			g0, g1 := edakit.Int(-20, -1), edakit.Int(-15, 15)
			if g1 == g0 {
				continue
			}
			B := float64(g1-g0) / float64(x1-x0)
			if B <= 0 {
				continue
			}
			x2 := float64(x1) - float64(g1)/B
			return num(edakit.Sprintf("In one dimension, BFGS reduces to the secant method. A sizing knob has derivative values $f'(%d) = %d$ and $f'(%d) = %d$. Using the secant curvature estimate, what is the next iterate $x_2$?", x0, g0, x1, g1)+roundNote(3), f(x2, 3),
				edakit.Sprintf("The curvature estimate is $B = (%d - (%d))/(%d - %d) = %s$, which satisfies the secant condition $B s = y$. Then $x_2 = %d - (%d)/%s = %s$.", g1, g0, x1, x0, f(B, 4), x1, g1, f(B, 4), f(x2, 3)),
				"b", f(B, 4), "x2", f(x2, 3))
		}
	case 2:
		for {
			s1, s2 := edakit.Int(-3, 3), edakit.Int(-3, 3)
			y1, y2 := edakit.Int(-4, 4), edakit.Int(-4, 4)
			ys := float64(y1*s1 + y2*s2)
			ss := float64(s1*s1 + s2*s2)
			if ys <= 0 || ss == 0 {
				continue
			}
			i, j := edakit.Int(0, 1), edakit.Int(0, 1)
			sv, yv := []float64{float64(s1), float64(s2)}, []float64{float64(y1), float64(y2)}
			id := 0.0
			if i == j {
				id = 1
			}
			v := id - sv[i]*sv[j]/ss + yv[i]*yv[j]/ys
			return num(edakit.Sprintf("BFGS starts from $B_0 = I$ (2 by 2). After a step $s = (%d, %d)$ the gradient changes by $y = (%d, %d)$. The update is $B_1 = B_0 - \\frac{B_0 s s^T B_0}{s^T B_0 s} + \\frac{y y^T}{y^T s}$. What is entry $(%d, %d)$ of $B_1$?", s1, s2, y1, y2, i+1, j+1)+roundNote(3), f(v, 3),
				edakit.Sprintf("$s^T s = %s$ and $y^T s = %s > 0$ (the curvature condition, which keeps $B_1$ positive definite). Entry $(%d, %d)$ is $%s - %s/%s + %s/%s = %s$.", f(ss, 0), f(ys, 0), i+1, j+1, f(id, 0), f(sv[i]*sv[j], 0), f(ss, 0), f(yv[i]*yv[j], 0), f(ys, 0), f(v, 3)),
				"ys", f(ys, 0), "entry", f(v, 3))
		}
	default:
		cells := edakit.Int(1, 20) * 100000
		m := edakit.Int(3, 20)
		n := 2 * cells
		mb := float64(2*m*n*8) / 1e6
		return num(edakit.Sprintf("A global placer optimizes x and y for %d cells (n = %d variables) with L-BFGS keeping m = %d correction pairs $(s_k, y_k)$ in 8-byte doubles. How many megabytes ($10^6$ bytes) do the stored pairs take?", cells, n, m)+roundNote(1), f(mb, 1),
			edakit.Sprintf("Each pair is two n-vectors, so storage is $2mn \\cdot 8 = 2 \\cdot %d \\cdot %d \\cdot 8$ bytes = %s MB. A dense BFGS matrix would need $n^2$ entries, which is impossible at this size.", m, n, f(mb, 1)),
			"mb", f(mb, 1))
	}
}

func coordinateDescent(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		k := ctx.KPOrdinal
		if k == 0 {
			k = edakit.Int(1, 2)
		}
		a, b := edakit.Int(2, 6), edakit.Int(3, 6)
		h := edakit.Pick([]int{-2, -1, 1, 2})
		u, v := edakit.Int(-12, 12), edakit.Int(-12, 12)
		x0, y0 := edakit.Int(-4, 4), edakit.Int(-4, 4)
		x1 := float64(u-h*y0) / float64(a)
		y1 := (float64(v) - float64(h)*x1) / float64(b)
		stem := edakit.Sprintf("Coordinate descent minimizes $f(x, y) = \\tfrac{1}{2}(%s) %s$ exactly along one coordinate at a time, x first, from $(%d, %d)$.", poly([]int{a, 2 * h, b}, []string{"x^2", "xy", "y^2"}), signed(poly([]int{-u, -v}, []string{"x", "y"})), x0, y0)
		if k == 1 {
			return num(stem+" What is x after the first coordinate update?"+roundNote(3), f(x1, 3),
				edakit.Sprintf("Hold $y = %d$ and set $\\partial f/\\partial x = %dx %s %d \\cdot (%d) - (%d) = 0$: $x = %s$.", y0, a, sgn(h), abs(h), y0, u, f(x1, 3)),
				"x1", f(x1, 3))
		}
		return num(stem+" What is y after one full sweep (x, then y)?"+roundNote(3), f(y1, 3),
			edakit.Sprintf("First $x_1 = (%d - (%d)(%d))/%d = %s$. Then with x fixed, $y_1 = (%d - (%d)(%s))/%d = %s$. This is exactly one Gauss-Seidel sweep on $\\nabla f = 0$.", u, h, y0, a, f(x1, 3), v, h, f(x1, 3), b, f(y1, 3)),
			"x1", f(x1, 3), "y1", f(y1, 3))
	default:
		P := edakit.Int(50, 120)
		w := []int{edakit.Int(1, 3), edakit.Int(1, 3), edakit.Int(1, 3), edakit.Int(1, 3)}
		c := []float64{float64(edakit.Int(0, P)), float64(edakit.Int(0, P)), float64(edakit.Int(0, P))}
		init := append([]float64(nil), c...)
		c[0] = (float64(w[0])*0 + float64(w[1])*c[1]) / float64(w[0]+w[1])
		c[1] = (float64(w[1])*c[0] + float64(w[2])*c[2]) / float64(w[1]+w[2])
		c[2] = (float64(w[2])*c[1] + float64(w[3])*float64(P)) / float64(w[2]+w[3])
		which := edakit.Int(2, 3)
		return num(edakit.Sprintf("A chain pad(0) - c1 - c2 - c3 - pad(%d) has net weights %d, %d, %d, %d from left to right. Cells start at c1 = %s, c2 = %s, c3 = %s. One sweep moves c1, then c2, then c3 to the weighted average of its two neighbours' current positions. Where is c%d after the sweep, in um?", P, w[0], w[1], w[2], w[3], f(init[0], 0), f(init[1], 0), f(init[2], 0), which)+roundNote(3), f(c[which-1], 3),
			edakit.Sprintf("Each move is the exact minimizer of quadratic wirelength in that one coordinate. In order: c1 = %s, c2 = %s, c3 = %s. Each update uses the newest neighbour values, as Gauss-Seidel does.", f(c[0], 3), f(c[1], 3), f(c[2], 3)),
			"c1", f(c[0], 3), "c2", f(c[1], 3), "c3", f(c[2], 3))
	}
}

func smoothWirelength(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 2, 4)
	xs := distinctInts(n, 0, 40)
	g := float64(edakit.Int(2, 10))
	lse := func(sign float64) float64 {
		s := 0.0
		for _, x := range xs {
			s += math.Exp(sign * float64(x) / g)
		}
		return g * math.Log(s)
	}
	mx, mn := xs[0], xs[0]
	for _, x := range xs {
		mx = max(mx, x)
		mn = min(mn, x)
	}
	stem := edakit.Sprintf("A net's pins sit at x = %s um. An analytical placer smooths HPWL with the log-sum-exp model and $\\gamma = %s$ um.", ints(xs, ", "), f(g, 0))
	switch kp(ctx) {
	case 1:
		v := lse(1)
		return num(stem+" What is the smooth maximum $\\gamma \\ln \\sum_i e^{x_i/\\gamma}$, in um?"+roundNote(3), f(v, 3),
			edakit.Sprintf("Summing the exponentials and taking the log gives %s um, slightly above the true maximum %d. The smooth version is differentiable everywhere, which gradient-based placers need.", f(v, 3), mx),
			"lse_max", f(v, 3))
	case 2:
		v := lse(1) + lse(-1)
		return num(stem+" What is the LSE wirelength $\\gamma\\ln\\sum e^{x_i/\\gamma} + \\gamma\\ln\\sum e^{-x_i/\\gamma}$, in um?"+roundNote(3), f(v, 3),
			edakit.Sprintf("The first term smooths max x = %d, the second smooths minus min x = %d. Together %s + (%s) = %s um, against an HPWL of %d um.", mx, mn, f(lse(1), 3), f(lse(-1), 3), f(v, 3), mx-mn),
			"lse_wl", f(v, 3), "hpwl", edakit.IntStr(mx-mn))
	default:
		if edakit.Int(0, 1) == 0 {
			b := 2 * g * math.Log(float64(n))
			return num(stem+" The LSE wirelength never underestimates HPWL. What is the worst-case overestimate for a net with this many pins, in um?"+roundNote(3), f(b, 3),
				edakit.Sprintf("Each smooth max exceeds the true one by at most $\\gamma \\ln n$, so the wirelength error is at most $2\\gamma\\ln n = 2 \\cdot %s \\ln %d = %s$ um.", f(g, 0), n, f(b, 3)),
				"bound", f(b, 3))
		}
		eps := float64(edakit.Int(1, 8))
		gg := eps / (2 * math.Log(float64(n)))
		return num(edakit.Sprintf("A placer wants the LSE wirelength of every %d-pin net within %s um of HPWL. What is the largest $\\gamma$ that guarantees it, in um?", n, f(eps, 0))+roundNote(3), f(gg, 3),
			edakit.Sprintf("The error bound is $2\\gamma\\ln n$, so $\\gamma \\le %s/(2\\ln %d) = %s$ um. Smaller $\\gamma$ is more accurate but makes the gradient nearly nonsmooth and the problem stiffer, so placers shrink $\\gamma$ gradually.", f(eps, 0), n, f(gg, 3)),
			"gamma", f(gg, 3))
	}
}

var _ = strings.Join

// signed prefixes a polynomial with "+" unless it already starts with a minus sign.
func signed(p string) string {
	if p == "0" {
		return ""
	}
	if strings.HasPrefix(p, "-") {
		return "- " + p[1:]
	}
	return "+ " + p
}
