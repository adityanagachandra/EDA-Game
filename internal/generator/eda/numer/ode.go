package numer

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// rcDraw draws an RC low-pass: R in kohm, C in pF (tau in ns), a step h in ns, the present
// capacitor voltage and the input.
func rcDraw() (r, c, tau, h, v0, vin float64) {
	r = float64(edakit.Int(1, 10))
	c = float64(edakit.Int(1, 10)) / 10
	tau = r * c
	h = float64(edakit.Int(1, 20)) / 10
	v0 = float64(edakit.Int(0, 10)) / 10
	vin = float64(edakit.Int(5, 18)) / 10
	return
}

func backwardEuler(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		a := float64(edakit.Int(1, 20)) / 2
		b := float64(edakit.Int(-10, 20)) / 2
		y0 := float64(edakit.Int(-10, 10)) / 2
		h := float64(edakit.Int(1, 10)) / 10
		y1 := (y0 + h*b) / (1 + a*h)
		ans := f(y1, 4)
		q := sp("Apply one backward Euler step of size $h = %s$ to $y' = -%sy + %s$ from $y_0 = %s$. What is $y_1$? Round to 4 decimal places.", f(h, 1), f(a, 1), f(b, 1), f(y0, 1))
		q = strings.ReplaceAll(q, "+ -", "- ")
		return numeric(q, ans, sp("Backward Euler evaluates the slope at the new point: $y_1 = y_0 + h(-ay_1 + b)$, so $y_1 = (y_0 + hb)/(1 + ah) = (%s)/(%s) = %s$.", f(y0+h*b, 4), f(1+a*h, 4), ans),
			map[string]string{"y1": ans})
	case 2:
		r, c, tau, h, v0, vin := rcDraw()
		v1 := (v0 + h/tau*vin) / (1 + h/tau)
		ans := f(v1, 4)
		q := sp("An RC low-pass has $R = %s$ kΩ and $C = %s$ pF; its capacitor sits at %s V when the input steps to %s V. Simulate with backward Euler, $h = %s$ ns. What is the capacitor voltage after one step, in volts? Round to 4 decimal places.", f(r, 0), f(c, 1), f(v0, 1), f(vin, 1), f(h, 1))
		return numeric(q, ans, sp("$\\tau = RC = %s$ ns and $h/\\tau = %s$. $v_1 = (v_0 + (h/\\tau)V_{in})/(1 + h/\\tau) = %s$ V.", f(tau, 2), f(h/tau, 4), ans),
			map[string]string{"tau": f(tau, 2)})
	default:
		a := [][]float64{{-float64(edakit.Int(2, 6)), float64(edakit.Int(0, 2))}, {float64(edakit.Int(0, 2)), -float64(edakit.Int(1, 6))}}
		x0 := []float64{float64(edakit.Int(-3, 3)), float64(edakit.Int(1, 4))}
		h := float64(edakit.Int(1, 5)) / 10
		m := [][]float64{{1 - h*a[0][0], -h * a[0][1]}, {-h * a[1][0], 1 - h*a[1][1]}}
		x1 := solve2(m, x0)
		i := edakit.Int(0, 1)
		ans := f(x1[i], 4)
		q := "Two coupled nodes obey $x' = Ax$ with $A = $ " + pmat(a, 0) + sp(" and $x(0) = (%s, %s)$. One backward Euler step of size $h = %s$ solves $(I - hA)x_1 = x_0$. What is component %d of $x_1$? Round to 4 decimal places.", f(x0[0], 0), f(x0[1], 0), f(h, 1), i+1)
		return numeric(q, ans, "$(I - hA) = $ "+pmat(m, 2)+sp(". Solving gives $x_1 = (%s, %s)$.", f(x1[0], 4), f(x1[1], 4)),
			map[string]string{"x1": f(x1[0], 4), "x2": f(x1[1], 4)})
	}
}

func trapezoidal(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		a := float64(edakit.Int(1, 20)) / 2
		b := float64(edakit.Int(-10, 20)) / 2
		y0 := float64(edakit.Int(-10, 10)) / 2
		h := float64(edakit.Int(1, 10)) / 10
		y1 := ((1-a*h/2)*y0 + h*b) / (1 + a*h/2)
		ans := f(y1, 4)
		q := sp("Apply one trapezoidal step of size $h = %s$ to $y' = -%sy + %s$ from $y_0 = %s$. What is $y_1$? Round to 4 decimal places.", f(h, 1), f(a, 1), f(b, 1), f(y0, 1))
		q = strings.ReplaceAll(q, "+ -", "- ")
		return numeric(q, ans, sp("$y_1 = y_0 + (h/2)(f_0 + f_1)$ gives $y_1 = ((1 - ah/2)y_0 + hb)/(1 + ah/2) = %s$.", ans),
			map[string]string{"y1": ans})
	case 2:
		r, c, tau, h, v0, vin := rcDraw()
		z := h / tau
		v1 := ((1-z/2)*v0 + z*vin) / (1 + z/2)
		ans := f(v1, 4)
		q := sp("An RC low-pass has $R = %s$ kΩ and $C = %s$ pF; its capacitor sits at %s V and the input is %s V, held constant over the step. Simulate with the trapezoidal rule, $h = %s$ ns. What is the capacitor voltage after one step, in volts? Round to 4 decimal places.", f(r, 0), f(c, 1), f(v0, 1), f(vin, 1), f(h, 1))
		return numeric(q, ans, sp("$\\tau = %s$ ns, $h/\\tau = %s$. $v_1 = ((1 - h/2\\tau)v_0 + (h/\\tau)V_{in})/(1 + h/2\\tau) = %s$ V.", f(tau, 2), f(z, 4), ans),
			map[string]string{"tau": f(tau, 2)})
	default:
		tf := float64(edakit.Int(1, 9))
		h := float64(edakit.Int(5, 40)) * tf
		e0 := float64(edakit.Int(1, 20)) / 10
		k := edakit.Int(1, 5)
		rr := (1 - h/(2*tf)) / (1 + h/(2*tf))
		ek := e0 * math.Pow(rr, float64(k))
		ans := f(ek, 4)
		q := sp("A parasitic mode with time constant %s ps carries an error of %s mV. The trapezoidal rule takes steps of $h = %s$ ps, far larger than the time constant. What is that error after %d steps, in mV (with sign)? Round to 4 decimal places.", f(tf, 0), f(e0, 1), f(h, 0), k)
		return numeric(q, ans, sp("Each step multiplies the mode by $R(z) = (1 + z/2)/(1 - z/2)$ with $z = -h/\\tau = %s$, so $R = %s$, close to $-1$. After %d steps the error is $%s \\cdot R^{%d} = %s$ mV: it alternates sign and barely decays, the trapezoidal ringing.", f(-h/tf, 1), f(rr, 4), k, f(e0, 1), k, ans),
			map[string]string{"R": f(rr, 6)})
	}
}

func companion(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		c := float64(edakit.Int(5, 200))
		h := float64(edakit.Int(1, 50))
		vn := float64(edakit.Int(0, 12)) / 10
		g := c / h
		if edakit.Int(0, 1) == 0 {
			ans := f(g, 4)
			q := sp("Backward Euler replaces a %s fF capacitor by a companion conductance in parallel with a current source. With timestep $h = %s$ ps, what is the companion conductance, in mS? Round to 4 decimal places.", f(c, 0), f(h, 0))
			return numeric(q, ans, sp("$G_{eq} = C/h = %s\\text{ fF}/%s\\text{ ps} = %s$ mS (1 fF/ps is 1 mS).", f(c, 0), f(h, 0), ans), nil)
		}
		ans := f(g*vn, 4)
		q := sp("Backward Euler models a %s fF capacitor, at %s V after the previous step, as a conductance $C/h$ in parallel with a history current source $(C/h)v_n$. With $h = %s$ ps, what is the history current, in mA? Round to 4 decimal places.", f(c, 0), f(vn, 1), f(h, 0))
		return numeric(q, ans, sp("$G_{eq} = %s$ mS, and $(C/h)v_n = %s \\cdot %s = %s$ mA.", f(g, 4), f(g, 4), f(vn, 1), ans), nil)
	case 2:
		c := float64(edakit.Int(5, 200))
		h := float64(edakit.Int(1, 50))
		vn := float64(edakit.Int(0, 12)) / 10
		in := float64(edakit.Int(-50, 50)) / 100
		g := 2 * c / h
		hist := g*vn + in
		ans := f(hist, 4)
		q := sp("The trapezoidal companion of a capacitor is $i_{n+1} = (2C/h)v_{n+1} - [(2C/h)v_n + i_n]$. For $C = %s$ fF, $h = %s$ ps, $v_n = %s$ V and capacitor current $i_n = %s$ mA, what is the history current $(2C/h)v_n + i_n$, in mA? Round to 4 decimal places.", f(c, 0), f(h, 0), f(vn, 1), f(in, 2))
		return numeric(q, ans, sp("$2C/h = %s$ mS, so the history current is $%s \\cdot %s + %s = %s$ mA. Unlike backward Euler, it needs the previous capacitor current.", f(g, 4), f(g, 4), f(vn, 1), f(in, 2), ans),
			map[string]string{"geq": f(g, 4)})
	default:
		if edakit.Int(0, 1) == 0 {
			l := float64(edakit.Int(1, 20))
			h := float64(edakit.Int(1, 50))
			in := float64(edakit.Int(-20, 20)) / 10
			v := float64(edakit.Int(-12, 12)) / 10
			g := h / l
			i1 := in + g*v
			ans := f(i1, 4)
			q := sp("Backward Euler turns an inductor $v = L\\,di/dt$ into a conductance $h/L$ in parallel with a current source $i_n$. For $L = %s$ nH, $h = %s$ ps and $i_n = %s$ mA, the solve gives $v_{n+1} = %s$ V across it. What is $i_{n+1}$, in mA? Round to 4 decimal places.", f(l, 0), f(h, 0), f(in, 1), f(v, 1))
			return numeric(q, ans, sp("$h/L = %s$ mS (1 ps/nH is 1 mS), so $i_{n+1} = i_n + (h/L)v_{n+1} = %s$ mA.", f(g, 4), ans), map[string]string{"geq": f(g, 4)})
		}
		vs := float64(edakit.Int(5, 18)) / 10
		r := float64(edakit.Int(1, 10)) // kohm -> 1/r mS
		c := float64(edakit.Int(5, 200))
		h := float64(edakit.Int(1, 50))
		vn := float64(edakit.Int(0, 10)) / 10
		g := c / h
		v1 := (vs/r + g*vn) / (1/r + g)
		ans := f(v1, 4)
		q := sp("A %s V source drives a %s kΩ resistor into a node with a %s fF capacitor to ground, at %s V. Replace the capacitor by its backward Euler companion with $h = %s$ ps and solve the resulting resistive circuit. What is the node voltage after the step, in volts? Round to 4 decimal places.", f(vs, 1), f(r, 0), f(c, 0), f(vn, 1), f(h, 0))
		return numeric(q, ans, sp("Companion: $G = %s$ mS with history current %s mA. KCL: $(V_s - v)/R + G v_n = G v$, so $v = (V_s/R + G v_n)/(1/R + G) = %s$ V.", f(g, 4), f(g*vn, 4), ans),
			map[string]string{"geq": f(g, 4)})
	}
}

func bdf2(ctx generator.GeneratorContext) generator.Problem {
	a := float64(edakit.Int(1, 20)) / 2
	b := float64(edakit.Int(-10, 20)) / 2
	h := float64(edakit.Int(1, 10)) / 10
	ode := strings.ReplaceAll(sp("$y' = -%sy + %s$", f(a, 1), f(b, 1)), "+ -", "- ")
	switch kpOf(ctx) {
	case 1:
		ym, y0, y1 := float64(edakit.Int(-20, 20))/10, float64(edakit.Int(-20, 20))/10, float64(edakit.Int(-20, 20))/10
		d := (3*y1 - 4*y0 + ym) / (2 * h)
		ans := f(d, 4)
		q := sp("BDF2 estimates $y'(t_{n+1})$ by differentiating the quadratic through three equally spaced samples. With step $h = %s$ and $y_{n-1} = %s$, $y_n = %s$, $y_{n+1} = %s$, what is that derivative estimate $(3y_{n+1} - 4y_n + y_{n-1})/(2h)$? Round to 4 decimal places.", f(h, 1), f(ym, 1), f(y0, 1), f(y1, 1))
		return numeric(q, ans, sp("$(3 \\cdot %s - 4 \\cdot %s + %s)/(2 \\cdot %s) = %s$.", f(y1, 1), f(y0, 1), f(ym, 1), f(h, 1), ans), nil)
	case 2:
		ym, y0 := float64(edakit.Int(-20, 20))/10, float64(edakit.Int(-20, 20))/10
		y1 := (4*y0 - ym + 2*h*b) / (3 + 2*a*h)
		ans := f(y1, 4)
		q := "Apply one BDF2 step $y_{n+1} = (4y_n - y_{n-1} + 2h f(y_{n+1}))/3$ to " + ode + sp(" with $h = %s$, $y_{n-1} = %s$ and $y_n = %s$. What is $y_{n+1}$? Round to 4 decimal places.", f(h, 1), f(ym, 1), f(y0, 1))
		return numeric(q, ans, sp("Substitute $f = -ay_{n+1} + b$: $(3 + 2ah)y_{n+1} = 4y_n - y_{n-1} + 2hb$, so $y_{n+1} = %s/%s = %s$.", f(4*y0-ym+2*h*b, 4), f(3+2*a*h, 4), ans),
			map[string]string{"y": ans})
	default:
		y0 := float64(edakit.Int(-20, 20)) / 10
		y1 := (y0 + h*b) / (1 + a*h)
		y2 := (4*y1 - y0 + 2*h*b) / (3 + 2*a*h)
		ans := f(y2, 4)
		q := "BDF2 needs two past values, so a simulator starts with one backward Euler step. Integrate " + ode + sp(" from $y_0 = %s$ with $h = %s$: one backward Euler step to $y_1$, then one BDF2 step to $y_2$. What is $y_2$? Round to 4 decimal places.", f(y0, 1), f(h, 1))
		return numeric(q, ans, sp("Backward Euler: $y_1 = (y_0 + hb)/(1 + ah) = %s$. BDF2: $y_2 = (4y_1 - y_0 + 2hb)/(3 + 2ah) = %s$.", f(y1, 4), ans),
			map[string]string{"y1": f(y1, 6)})
	}
}

func stabilityRegions(ctx generator.GeneratorContext) generator.Problem {
	method := edakit.Pick([]string{"forward Euler", "backward Euler", "trapezoidal"})
	R := func(z float64) float64 {
		switch method {
		case "forward Euler":
			return 1 + z
		case "backward Euler":
			return 1 / (1 - z)
		}
		return (1 + z/2) / (1 - z/2)
	}
	switch kpOf(ctx) {
	case 1:
		z := -float64(edakit.Int(1, 60)) / 4
		ans := f(R(z), 4)
		q := sp("For the test equation $y' = \\lambda y$, one step of a method multiplies $y$ by its amplification factor $R(z)$, $z = h\\lambda$. What is $R(z)$ for %s at $z = %s$? Round to 4 decimal places.", method, f(z, 2))
		return numeric(q, ans, sp("Forward Euler $R = 1 + z$, backward Euler $R = 1/(1 - z)$, trapezoidal $R = (1 + z/2)/(1 - z/2)$. For %s at $z = %s$: %s.", method, f(z, 2), ans),
			map[string]string{"z": f(z, 2)})
	case 2:
		lam := float64(edakit.Int(1, 50)) // 1/ns
		if edakit.Int(0, 1) == 0 {
			hmax := 2 / lam * 1000
			ans := f(hmax, 2)
			q := sp("A circuit mode decays like $e^{-%st}$ with $t$ in ns. What is the largest step forward Euler can take without the mode growing, in ps? Round to 2 decimal places.", f(lam, 0))
			return numeric(q, ans, sp("Forward Euler needs $|1 + h\\lambda| \\le 1$, i.e. $h \\le 2/|\\lambda| = %s$ ps.", ans), nil)
		}
		h := float64(edakit.Int(1, 100)) / 1000 * 2
		z := -lam * h
		ans := f(math.Abs(R(z)), 4)
		q := sp("A circuit mode has $\\lambda = -%s$ per ns. Using %s with $h = %s$ ns, what is $|R(h\\lambda)|$, the factor by which the mode is scaled per step? Round to 4 decimal places.", f(lam, 0), method, f(h, 3))
		return numeric(q, ans, sp("$z = h\\lambda = %s$ and $|R(z)| = %s$. A value above 1 means the computed mode grows although the true one decays.", f(z, 3), ans),
			map[string]string{"z": f(z, 3)})
	default:
		w := float64(edakit.Int(1, 20))
		h := float64(edakit.Int(1, 20)) / 100
		n := edakit.Int(5, 50)
		wh := w * h
		var amp float64
		switch method {
		case "forward Euler":
			amp = math.Pow(math.Sqrt(1+wh*wh), float64(n))
		case "backward Euler":
			amp = math.Pow(1/math.Sqrt(1+wh*wh), float64(n))
		default:
			amp = 1
		}
		ans := f(amp, 4)
		q := sp("An undamped LC tank oscillates at $\\omega = %s$ rad/ns, so its modes have $\\lambda = \\pm i\\omega$. Simulated with %s and $h = %s$ ns for %d steps, by what factor has the oscillation amplitude changed? Round to 4 decimal places.", f(w, 0), method, f(h, 2), n)
		return numeric(q, ans, sp("$\\omega h = %s$. Forward Euler gives $|R| = \\sqrt{1 + (\\omega h)^2}$ (growth), backward Euler $1/\\sqrt{1 + (\\omega h)^2}$ (artificial damping), trapezoidal exactly 1. For %s, $|R|^{%d} = %s$.", f(wh, 2), method, n, ans),
			map[string]string{"wh": f(wh, 2)})
	}
}

func stiffness(ctx generator.GeneratorContext) generator.Problem {
	a := float64(edakit.Int(500, 3000))
	d := float64(edakit.Int(1, 10))
	b := float64(edakit.Int(1, 8))
	mid := (a + d) / 2
	rad := math.Sqrt((a-d)*(a-d)/4 + b*b)
	lmax, lmin := mid+rad, mid-rad
	m := [][]float64{{-a, b}, {b, -d}}
	stem := "Two coupled RC nodes obey $x' = Ax$ with $A = $ " + pmat(m, 0) + " per ns. "
	switch kpOf(ctx) {
	case 1:
		ans := f(lmax/lmin, 1)
		q := stem + "Its eigenvalues are real and negative. What is the stiffness ratio $|\\lambda|_{max}/|\\lambda|_{min}$? Round to 1 decimal place."
		return numeric(q, ans, sp("For the symmetric matrix the eigenvalues are $-(%s \\pm %s)$, i.e. $-%s$ and $-%s$ per ns. Ratio %s.", f(mid, 2), f(rad, 4), f(lmax, 4), f(lmin, 4), ans),
			map[string]string{"lmax": f(lmax, 4), "lmin": f(lmin, 4)})
	case 2:
		hmax := 2 / lmax * 1000
		ans := f(hmax, 4)
		q := stem + "What is the largest stable forward Euler step, in ps? Round to 4 decimal places."
		return numeric(q, ans, sp("The fastest eigenvalue is $-%s$ per ns. Forward Euler needs $h \\le 2/|\\lambda_{max}| = %s$ ps, however slowly the solution itself changes.", f(lmax, 4), ans),
			map[string]string{"lmax": f(lmax, 4)})
	default:
		tf := float64(edakit.Int(1, 20))      // ps
		ts := float64(edakit.Int(1, 50)) * 10 // ns
		steps := math.Ceil(5 * ts * 1000 / (2 * tf))
		ans := f(steps, 0)
		q := sp("A circuit has a parasitic time constant of %s ps and a dominant time constant of %s ns. Simulating to $5\\tau_{slow}$ with forward Euler at its stability limit $h = 2\\tau_{fast}$, how many steps are needed? (Backward Euler could take about 100 accuracy-limited steps.)", f(tf, 0), f(ts, 0))
		return numeric(q, ans, sp("$5 \\cdot %s$ ns $/ (2 \\cdot %s$ ps$) = %s$ steps, rounded up.", f(ts, 0), f(tf, 0), ans), nil)
	}
}

func lte(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		h := float64(edakit.Int(1, 20)) / 100
		ypp := float64(edakit.Int(1, 50)) / 10
		l := h * h / 2 * ypp * 1000
		ans := f(l, 4)
		q := sp("A node voltage has $|v''| \\approx %s$ V/ns$^2$ near $t_n$. What is the backward Euler local truncation error of a step $h = %s$ ns, in mV? Round to 4 decimal places.", f(ypp, 1), f(h, 2))
		return numeric(q, ans, sp("$\\text{LTE}_{BE} \\approx (h^2/2)|v''| = (%s^2/2)(%s)$ V $= %s$ mV.", f(h, 2), f(ypp, 1), ans), nil)
	case 2:
		h0 := float64(edakit.Int(5, 20)) / 100
		h1 := float64(edakit.Int(5, 20)) / 100
		hn := float64(edakit.Int(5, 20)) / 100
		x0 := float64(edakit.Int(0, 500)) / 1000
		x1 := x0 + float64(edakit.Int(-100, 200))/1000
		x2 := x1 + float64(edakit.Int(-100, 200))/1000
		d1 := (x1 - x0) / h0
		d2 := (x2 - x1) / h1
		dd2 := (d2 - d1) / (h0 + h1)
		ypp := 2 * dd2
		l := hn * hn / 2 * math.Abs(ypp) * 1000
		ans := f(l, 4)
		q := sp("A simulator has accepted $v = %s$, %s, %s V at times separated by %s ns and then %s ns. It estimates $v'' \\approx 2\\,v[t_0, t_1, t_2]$ from the second divided difference. What backward Euler LTE does that predict for a next step of %s ns, in mV? Round to 4 decimal places.", f(x0, 3), f(x1, 3), f(x2, 3), f(h0, 2), f(h1, 2), f(hn, 2))
		return numeric(q, ans, sp("First differences %s and %s V/ns; second divided difference $(%s - %s)/%s = %s$; $v'' \\approx %s$. LTE $\\approx (h^2/2)|v''| = %s$ mV.", f(d1, 4), f(d2, 4), f(d2, 4), f(d1, 4), f(h0+h1, 2), f(dd2, 4), f(ypp, 4), ans),
			map[string]string{"vpp": f(ypp, 6)})
	default:
		h := float64(edakit.Int(1, 20)) / 100
		y3 := float64(edakit.Int(1, 100)) / 10
		if edakit.Int(0, 1) == 0 {
			l := h * h * h / 12 * y3 * 1000
			ans := f(l, 4)
			q := sp("A node voltage has $|v'''| \\approx %s$ V/ns$^3$. What is the trapezoidal local truncation error of a step $h = %s$ ns, in mV? Round to 4 decimal places.", f(y3, 1), f(h, 2))
			return numeric(q, ans, sp("$|\\text{LTE}_{TR}| \\approx (h^3/12)|v'''| = %s$ mV.", ans), nil)
		}
		l := 2 * h * h * h / 9 * y3 * 1000
		ans := f(l, 4)
		q := sp("A node voltage has $|v'''| \\approx %s$ V/ns$^3$. What is the BDF2 local truncation error of an equal step $h = %s$ ns, in mV? Round to 4 decimal places.", f(y3, 1), f(h, 2))
		return numeric(q, ans, sp("$|\\text{LTE}_{BDF2}| \\approx (2/9)h^3|v'''| = %s$ mV, larger than trapezoidal's $h^3/12$ for the same step.", ans), nil)
	}
}

func stepControl(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		tol := float64(edakit.Int(1, 9)) / 10
		n := edakit.Int(4, 6)
		var parts []string
		rej := 0
		for i := 0; i < n; i++ {
			e := tol * float64(edakit.Int(2, 25)) / 10
			if e > tol {
				rej++
			}
			parts = append(parts, f(e, 3))
		}
		q := sp("A timestep controller accepts a step when its LTE estimate is at most the tolerance of %s mV. The last %d attempted steps had estimates %s mV. How many were rejected?", f(tol, 1), n, strings.Join(parts, ", "))
		return numeric(q, itoa(rej), sp("Each estimate above %s mV fails: %d rejections, each retried with a smaller step.", f(tol, 1), rej), nil)
	case 2:
		p := edakit.Int(1, 2)
		h := float64(edakit.Int(5, 100))
		tol := float64(edakit.Int(1, 9)) / 10
		err := tol * float64(edakit.Int(2, 40)) / 10
		hn := h * 0.9 * math.Pow(tol/err, 1/float64(p+1))
		ans := f(hn, 3)
		name := map[int]string{1: "backward Euler (order 1)", 2: "trapezoidal (order 2)"}[p]
		q := sp("Using %s, a %s ps step produced an LTE estimate of %s mV against a tolerance of %s mV. With safety factor 0.9, the controller proposes $h_{new} = 0.9\\,h\\,(\\text{tol}/\\text{err})^{1/(p+1)}$. What is $h_{new}$, in ps? Round to 3 decimal places.", name, f(h, 0), f(err, 3), f(tol, 1))
		return numeric(q, ans, sp("LTE scales as $h^{p+1} = h^{%d}$, so $h_{new} = 0.9 \\cdot %s \\cdot (%s/%s)^{1/%d} = %s$ ps.", p+1, f(h, 0), f(tol, 1), f(err, 3), p+1, ans),
			map[string]string{"hnew": ans})
	default:
		h := float64(edakit.Int(5, 100))
		tol := float64(edakit.Int(1, 9)) / 10
		err := tol * float64(edakit.Int(1, 20)) / 100
		hn := h * 0.9 * math.Pow(tol/err, 1.0/3)
		t := float64(edakit.Int(100, 900))
		tb := t + float64(edakit.Int(5, 300))
		next := math.Min(hn, math.Min(2*h, tb-t))
		ans := f(next, 3)
		q := sp("A trapezoidal step of %s ps ends at $t = %s$ ps with an LTE estimate of %s mV against a tolerance of %s mV. The controller proposes $h_{new} = 0.9\\,h\\,(\\text{tol}/\\text{err})^{1/3}$, caps growth at twice the last step, and never steps past the next source breakpoint at %s ps. What is the next step, in ps? Round to 3 decimal places.", f(h, 0), f(t, 0), f(err, 3), f(tol, 1), f(tb, 0))
		return numeric(q, ans, sp("Formula: %s ps. Growth cap: %s ps. Distance to breakpoint: %s ps. The next step is the smallest, %s ps.", f(hn, 3), f(2*h, 0), f(tb-t, 0), ans),
			map[string]string{"hnew": f(hn, 3)})
	}
}

var daeItems = []edakit.Item{
	{KP: 1, Q: "MNA for a circuit gives $C x' + G x = u(t)$. When is this a DAE rather than an ODE?",
		A: "When C is singular, for example at a node with no capacitor or for a voltage-source branch current", D: []string{"When G is singular", "When the circuit contains a nonlinear device", "When u(t) is discontinuous"},
		Why: "A singular C means some equations contain no derivatives: they are algebraic constraints, which makes the system differential-algebraic."},
	{KP: 1, Q: "An MNA system has 8 unknowns and its capacitance matrix C has rank 5. How many of its equations are purely algebraic (after a change of variables)?",
		A: "3", D: []string{"5", "8", "0"},
		Why: "The number of algebraic equations is n minus rank(C): 8 - 5 = 3."},
	{KP: 1, Q: "Why can a simulator not hand $C x' + G x = u$ to a standard explicit ODE solver by writing $x' = C^{-1}(u - Gx)$?",
		A: "C is generally singular in MNA, so its inverse does not exist", D: []string{"C is always diagonal", "Explicit solvers cannot handle linear systems", "G must be inverted instead of C"},
		Why: "Nodes without capacitance and source branch currents give zero rows in C, so the ODE form does not exist."},
	{KP: 1, Q: "A resistive node with no capacitor connects two capacitive nodes. What kind of equation is its KCL row?",
		A: "An algebraic constraint that fixes its voltage from its neighbours at every instant", D: []string{"A differential equation with zero time constant", "An equation that can be dropped", "A stiff differential equation with infinite eigenvalue that must be integrated explicitly"},
		Why: "With no capacitor, the KCL row has no derivative term, so the voltage follows the neighbouring voltages instantaneously."},
	{KP: 2, Q: "What is the differentiation index of a DAE, informally?",
		A: "The number of times the algebraic equations must be differentiated to obtain an ODE for all variables", D: []string{"The order of the highest derivative appearing", "The number of algebraic equations", "The rank of the capacitance matrix"},
		Why: "Index 1 means one differentiation of the constraints gives an ODE; higher index means hidden constraints and harder numerics."},
	{KP: 2, Q: "Which circuit structure raises the DAE index above 1?",
		A: "A loop made only of capacitors and voltage sources", D: []string{"A resistor in parallel with a capacitor", "A node with a capacitor to ground", "A diode in series with a resistor"},
		Why: "Capacitor-voltage-source loops (and inductor-current-source cutsets) create constraints among differential variables, which raises the index."},
	{KP: 2, Q: "A transient starts from initial conditions that violate an algebraic constraint of the circuit. What happens in an implicit integrator?",
		A: "The first step is forced onto the constraint, which shows up as a jump or an impulse in some currents", D: []string{"Nothing: algebraic constraints are only enforced at DC", "The simulator integrates the inconsistent state indefinitely", "The error decays smoothly with the circuit time constant"},
		Why: "Inconsistent initial conditions cannot be held; the solver jumps to a consistent state, and derivatives of the jump appear as spikes."},
	{KP: 2, Q: "How does a SPICE transient usually obtain consistent initial conditions?",
		A: "It solves the DC operating point first and starts the transient from it", D: []string{"It sets every node voltage to zero", "It integrates backward in time from the first breakpoint", "It uses the initial guess of the Newton solver"},
		Why: "The DC solution satisfies all algebraic constraints with derivatives zero, so it is consistent by construction."},
	{KP: 3, Q: "Backward Euler on $C x' + G x = u$ gives $(C/h + G)x_{n+1} = (C/h)x_n + u_{n+1}$. What must hold for each step to be solvable?",
		A: "The matrix C/h + G must be nonsingular", D: []string{"C must be invertible", "G must be symmetric", "h must be smaller than every time constant"},
		Why: "Implicit BDF steps never invert C alone; they solve with the pencil C/h + G, which is nonsingular for well-posed circuits."},
	{KP: 3, Q: "Why are BDF methods (backward Euler, Gear) the classic choice for index-1 circuit DAEs?",
		A: "They treat the algebraic variables implicitly at the new time point and damp errors in them", D: []string{"They are explicit and avoid matrix factorization", "They are the only methods with order above 1", "They make C invertible"},
		Why: "BDF methods evaluate everything at t_{n+1}, so algebraic constraints hold at each accepted point, and their L-stability damps transient errors."},
	{KP: 3, Q: "On an index-1 DAE, the trapezoidal rule can show oscillation in an algebraic variable that never decays. Why?",
		A: "Its amplification factor tends to -1 for infinitely stiff components, and algebraic equations behave like infinitely stiff ones", D: []string{"The trapezoidal rule is explicit", "The trapezoidal rule is only first order", "The algebraic equations are ignored by the trapezoidal rule"},
		Why: "An algebraic constraint is the limit of an infinitely fast mode; the trapezoidal factor approaches -1 there, so errors flip sign without decaying."},
	{KP: 3, Q: "After a source breakpoint, SPICE often takes a few backward Euler steps before returning to trapezoidal. Which DAE issue does that address?",
		A: "It damps the sign-alternating errors that trapezoidal leaves in stiff and algebraic components after a discontinuity", D: []string{"It raises the order of accuracy at the breakpoint", "It makes the capacitance matrix nonsingular", "It removes the need for a Newton solve"},
		Why: "Backward Euler is L-stable, so a few steps after a discontinuity remove the components that trapezoidal would leave ringing."},
}
