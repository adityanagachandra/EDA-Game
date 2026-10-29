package stat

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// lsq returns slope, intercept and R^2 of the least-squares line.
func lsq(x, y []float64) (b1, b0, r2 float64) {
	n := float64(len(x))
	var sx, sy float64
	for i := range x {
		sx += x[i]
		sy += y[i]
	}
	mx, my := sx/n, sy/n
	var sxy, sxx, syy float64
	for i := range x {
		sxy += (x[i] - mx) * (y[i] - my)
		sxx += (x[i] - mx) * (x[i] - mx)
		syy += (y[i] - my) * (y[i] - my)
	}
	b1 = sxy / sxx
	b0 = my - b1*mx
	r2 = sxy * sxy / (sxx * syy)
	return
}

// regression: slope of delay against load (KP 1), predicted delay (KP 2), and R^2 (KP 3).
func regression(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	slope := float64(edakit.Int(5, 40)) / 10
	icpt := float64(edakit.Int(10, 40))
	x := make([]float64, n)
	y := make([]float64, n)
	var pts []string
	for i := 0; i < n; i++ {
		x[i] = float64((i + 1) * edakit.Pick([]int{2, 5, 10}))
		if i > 0 && x[i] <= x[i-1] {
			x[i] = x[i-1] + 2
		}
		y[i] = icpt + slope*x[i] + float64(edakit.Int(-8, 8))/2
		pts = append(pts, edakit.Sprintf("(%s, %s)", f(x[i], 0), f(y[i], 1)))
	}
	b1, b0, r2 := lsq(x, y)
	stem := edakit.Sprintf("A characterization run measures cell delay (ps) against output load (fF) at the points %s.", strings.Join(pts, ", "))
	facts := map[string]string{"slope": f(b1, 4), "intercept": f(b0, 4)}
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(b1, 3)
		return edakit.Numeric(stem+" What is the least-squares slope, in ps per fF? Round to 3 decimal places.",
			ans, edakit.Sprintf("$b_1 = S_{xy}/S_{xx} = %s$ ps per fF.", ans), facts)
	case 2:
		x0 := x[n-1] + float64(edakit.Int(1, 4))*x[0]/2
		ans := edakit.Num(b0+b1*x0, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" Using the least-squares line, what delay does the model predict at %s fF, in ps? Round to 2 decimal places.", f(x0, 1)),
			ans, edakit.Sprintf("$b_1 = %s$, $b_0 = \\bar{y} - b_1\\bar{x} = %s$, so $\\hat{y} = %s + %s(%s) = %s$ ps. This load is outside the measured range, so the prediction is an extrapolation.", f(b1, 4), f(b0, 4), f(b0, 4), f(b1, 4), f(x0, 1), ans), facts)
	default:
		ans := edakit.Num(r2, 4)
		return edakit.Numeric(stem+" What is the coefficient of determination $R^2$ of the least-squares line? Round to 4 decimal places.",
			ans, edakit.Sprintf("$R^2 = S_{xy}^2/(S_{xx}S_{yy}) = %s$.", ans), facts)
	}
}

// doeRuns: full factorial run count (KP 1), fractional factorial (KP 2), and quadratic model or
// central composite design size (KP 3).
func doeRuns(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		lv := edakit.Int(2, 4)
		k := edakit.Int(2, 7)
		if lv == 4 && k > 5 {
			k = 5
		}
		n := int(math.Pow(float64(lv), float64(k)))
		return edakit.Numeric(edakit.Sprintf("A characterization study varies %d process and environment factors, each at %d levels. How many simulations does a full factorial design need?", k, lv),
			edakit.IntStr(n), edakit.Sprintf("Every combination is run: $%d^{%d} = %d$.", lv, k, n), map[string]string{"runs": edakit.IntStr(n)})
	case 2:
		k := edakit.Int(4, 10)
		p := edakit.Int(1, k-3)
		n := 1 << (k - p)
		return edakit.Numeric(edakit.Sprintf("A two-level screening experiment has %d factors and uses a $2^{%d-%d}$ fractional factorial. How many runs does it need?", k, k, p),
			edakit.IntStr(n), edakit.Sprintf("$2^{k-p} = 2^{%d} = %d$, a $1/2^{%d}$ fraction of the %d runs of the full design.", k-p, n, p, 1<<k), map[string]string{"runs": edakit.IntStr(n)})
	default:
		k := edakit.Int(2, 8)
		if edakit.Int(0, 1) == 0 {
			c := (k + 1) * (k + 2) / 2
			return edakit.Numeric(edakit.Sprintf("A full quadratic response surface in %d factors has an intercept, linear, pure quadratic and two-factor interaction terms. How many coefficients must be fitted?", k),
				edakit.IntStr(c), edakit.Sprintf("$1 + k + k + \\binom{k}{2} = (k+1)(k+2)/2 = %d$. The design needs at least that many distinct runs.", c), map[string]string{"coeffs": edakit.IntStr(c)})
		}
		nc := edakit.Int(1, 6)
		n := (1 << k) + 2*k + nc
		return edakit.Numeric(edakit.Sprintf("A central composite design for %d factors uses a full two-level factorial core, axial (star) points and %d centre runs. How many runs is that?", k, nc),
			edakit.IntStr(n), edakit.Sprintf("$2^k + 2k + n_c = %d + %d + %d = %d$.", 1<<k, 2*k, nc, n), map[string]string{"runs": edakit.IntStr(n)})
	}
}

// factorialEffects: a 2x2 design. Main effect (KP 1), interaction (KP 2), and a prediction from
// the coded regression model (KP 3).
func factorialEffects(ctx generator.GeneratorContext) generator.Problem {
	base := float64(edakit.Int(40, 80))
	ea := float64(edakit.Int(-20, 20)) / 2
	eb := float64(edakit.Int(-20, 20)) / 2
	eab := float64(edakit.Int(-8, 8)) / 2
	yv := func(a, b float64) float64 {
		return base + ea/2*a + eb/2*b + eab/2*a*b + float64(edakit.Int(-2, 2))/2
	}
	ymm, ypm, ymp, ypp := yv(-1, -1), yv(1, -1), yv(-1, 1), yv(1, 1)
	fa := edakit.Pick([]string{"supply voltage", "gate length", "temperature"})
	fb := edakit.Pick([]string{"oxide thickness", "threshold implant dose", "load capacitance"})
	stem := edakit.Sprintf("A $2^2$ experiment varies A (%s) and B (%s) between coded levels $-1$ and $+1$ and measures a delay in ps: $(-,-) = %s$, $(+,-) = %s$, $(-,+) = %s$, $(+,+) = %s$.", fa, fb, f(ymm, 2), f(ypm, 2), f(ymp, 2), f(ypp, 2))
	effA := (ypm+ypp)/2 - (ymm+ymp)/2
	effB := (ymp+ypp)/2 - (ymm+ypm)/2
	effAB := (ymm+ypp)/2 - (ypm+ymp)/2
	facts := map[string]string{"effect_a": f(effA, 3), "effect_b": f(effB, 3), "effect_ab": f(effAB, 3)}
	switch kp(ctx) {
	case 1:
		if edakit.Int(0, 1) == 0 {
			ans := edakit.Num(effA, 2)
			return edakit.Numeric(stem+" What is the main effect of A, in ps? Round to 2 decimal places.",
				ans, edakit.Sprintf("Average at A high minus average at A low: $(%s + %s)/2 - (%s + %s)/2 = %s$ ps.", f(ypm, 2), f(ypp, 2), f(ymm, 2), f(ymp, 2), ans), facts)
		}
		ans := edakit.Num(effB, 2)
		return edakit.Numeric(stem+" What is the main effect of B, in ps? Round to 2 decimal places.",
			ans, edakit.Sprintf("Average at B high minus average at B low: $(%s + %s)/2 - (%s + %s)/2 = %s$ ps.", f(ymp, 2), f(ypp, 2), f(ymm, 2), f(ypm, 2), ans), facts)
	case 2:
		ans := edakit.Num(effAB, 2)
		return edakit.Numeric(stem+" What is the AB interaction effect, in ps? Round to 2 decimal places.",
			ans, edakit.Sprintf("Average where $AB = +1$ minus average where $AB = -1$: $(%s + %s)/2 - (%s + %s)/2 = %s$ ps.", f(ymm, 2), f(ypp, 2), f(ypm, 2), f(ymp, 2), ans), facts)
	default:
		xa := float64(edakit.Int(-4, 4)) / 4
		xb := float64(edakit.Int(-4, 4)) / 4
		b0 := (ymm + ypm + ymp + ypp) / 4
		pred := b0 + effA/2*xa + effB/2*xb + effAB/2*xa*xb
		ans := edakit.Num(pred, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" Fit the coded model $y = b_0 + b_A x_A + b_B x_B + b_{AB} x_A x_B$. What delay does it predict at $x_A = %s$, $x_B = %s$, in ps? Round to 2 decimal places.", f(xa, 2), f(xb, 2)),
			ans, edakit.Sprintf("$b_0$ is the grand mean %s and each coefficient is half its effect: $b_A = %s$, $b_B = %s$, $b_{AB} = %s$. The prediction is %s ps.", f(b0, 3), f(effA/2, 3), f(effB/2, 3), f(effAB/2, 3), ans), facts)
	}
}

// responseSurface: stationary point of a one-factor quadratic through three coded runs (KP 1),
// the predicted optimum (KP 2), and the stationary point of a two-factor quadratic (KP 3).
func responseSurface(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		y0 := float64(edakit.Int(40, 80))
		b2 := float64(edakit.Int(2, 12)) / 2
		xs := float64(edakit.Int(-8, 8)) / 10
		if xs == 0 {
			xs = 0.3
		}
		b1 := -2 * b2 * xs
		ym := y0 - b1 + b2
		yp := y0 + b1 + b2
		stem := edakit.Sprintf("A leakage-versus-delay figure of merit is simulated at a sizing factor coded $x = -1, 0, +1$, giving %s, %s and %s.", f(ym, 2), f(y0, 2), f(yp, 2))
		bb1 := (yp - ym) / 2
		bb2 := (yp + ym - 2*y0) / 2
		x := -bb1 / (2 * bb2)
		if ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			ans := edakit.Num(x, 3)
			return edakit.Numeric(stem+" Fit $y = b_0 + b_1 x + b_2 x^2$ through the three points. At what coded $x$ is the fitted surface stationary? Round to 3 decimal places.",
				ans, edakit.Sprintf("$b_1 = (y_+ - y_-)/2 = %s$, $b_2 = (y_+ + y_- - 2y_0)/2 = %s$. Setting $b_1 + 2b_2x = 0$ gives $x^* = %s$.", f(bb1, 3), f(bb2, 3), ans),
				map[string]string{"b1": f(bb1, 3), "b2": f(bb2, 3)})
		}
		ans := edakit.Num(y0-bb1*bb1/(4*bb2), 3)
		return edakit.Numeric(stem+" Fit $y = b_0 + b_1 x + b_2 x^2$ through the three points. What is the fitted minimum value? Round to 3 decimal places.",
			ans, edakit.Sprintf("$b_0 = %s$, $b_1 = %s$, $b_2 = %s$. The minimum is $b_0 - b_1^2/(4b_2) = %s$, at $x^* = %s$.", f(y0, 2), f(bb1, 3), f(bb2, 3), ans, f(x, 3)),
			map[string]string{"b1": f(bb1, 3), "b2": f(bb2, 3)})
	default:
		b1 := float64(edakit.Int(-10, 10))
		b2 := float64(edakit.Int(-10, 10))
		b11 := float64(edakit.Int(2, 8))
		b22 := float64(edakit.Int(2, 8))
		b12 := float64(edakit.Int(-3, 3))
		// Gradient zero: [2b11 b12; b12 2b22] x = -[b1; b2].
		det := 4*b11*b22 - b12*b12
		x1 := (-b1*2*b22 + b2*b12) / det
		x2 := (-b2*2*b11 + b1*b12) / det
		which := edakit.Int(1, 2)
		x := x1
		if which == 2 {
			x = x2
		}
		ans := edakit.Num(x, 3)
		return edakit.Numeric(edakit.Sprintf("A fitted response surface in coded factors is $y = 50 %s x_1 %s x_2 + %s x_1^2 + %s x_2^2 %s x_1 x_2$. At its stationary point, what is $x_%d$? Round to 3 decimal places.", signed(b1), signed(b2), f(b11, 0), f(b22, 0), signed(b12), which),
			ans, edakit.Sprintf("Set the gradient to zero: $%s x_1 + %s x_2 = %s$ and $%s x_1 + %s x_2 = %s$. The determinant is %s, giving $x_1 = %s$ and $x_2 = %s$.", f(2*b11, 0), f(b12, 0), f(-b1, 0), f(b12, 0), f(2*b22, 0), f(-b2, 0), f(det, 0), f(x1, 4), f(x2, 4)),
			map[string]string{"x1": f(x1, 4), "x2": f(x2, 4)})
	}
}

// signed renders a coefficient with its sign for a polynomial, "+ 3" or "- 3".
func signed(v float64) string {
	if v < 0 {
		return "- " + f(-v, 0)
	}
	return "+ " + f(v, 0)
}

// surrogateValidation: hold-out RMSE (KP 1), a leave-one-out residual from leverage (KP 2), and
// the coefficient count of a polynomial surrogate (KP 3).
func surrogateValidation(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 6)
		var parts []string
		ss := 0.0
		for i := 0; i < n; i++ {
			pred := float64(edakit.Int(500, 900)) / 10
			e := float64(edakit.Int(-30, 30)) / 10
			parts = append(parts, edakit.Sprintf("predicted %s, simulated %s", f(pred, 1), f(pred+e, 1)))
			ss += e * e
		}
		r := math.Sqrt(ss / float64(n))
		ans := edakit.Num(r, 3)
		return edakit.Numeric(edakit.Sprintf("A delay surrogate is checked on %d held-out SPICE runs (ps): %s. What is the hold-out root-mean-square error, in ps? Round to 3 decimal places.", n, strings.Join(parts, "; ")),
			ans, edakit.Sprintf("The sum of squared errors is %s; dividing by %d and taking the root gives %s ps.", f(ss, 2), n, ans), map[string]string{"sse": f(ss, 2)})
	case 2:
		e := float64(edakit.Int(-40, 40)) / 10
		for e == 0 {
			e = 0.7
		}
		h := float64(edakit.Int(5, 90)) / 100
		loo := e / (1 - h)
		ans := edakit.Num(loo, 3)
		return edakit.Numeric(edakit.Sprintf("In a least-squares surrogate fit, training point $i$ has residual %s ps and leverage (hat-matrix diagonal) $h_{ii} = %s$. What is its leave-one-out prediction residual, in ps? Round to 3 decimal places.", f(e, 1), f(h, 2)),
			ans, edakit.Sprintf("$e_{(i)} = e_i/(1 - h_{ii}) = %s/%s = %s$ ps, with no refit needed. High-leverage points hide their errors in the training residual.", f(e, 1), f(1-h, 2), ans), map[string]string{"loo": ans})
	default:
		k := edakit.Int(3, 12)
		d := edakit.Int(1, 3)
		c := binom(k+d, d)
		return edakit.Numeric(edakit.Sprintf("A polynomial surrogate of total degree %d in %d process parameters includes every monomial up to that degree. How many coefficients does it have?", d, k),
			edakit.IntStr(c), edakit.Sprintf("The number of monomials of degree at most $d$ in $k$ variables is $\\binom{k+d}{d} = \\binom{%d}{%d} = %d$. Fitting needs comfortably more simulations than that.", k+d, d, c), map[string]string{"coeffs": edakit.IntStr(c)})
	}
}

func binom(n, k int) int {
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}
