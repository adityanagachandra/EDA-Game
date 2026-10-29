package em

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerSignal(reg *generator.Registry) {
	reg.Register("em.signal.lti_convolution", edakit.Func(ltiConvolution))
	reg.Register("em.signal.fourier_spectrum", edakit.Func(fourierSpectrum))
	reg.Register("em.signal.laplace_transfer", edakit.Func(laplaceTransfer))
	reg.Register("em.signal.poles_zeros", edakit.Func(polesZeros))
	reg.Register("em.signal.bode", edakit.Func(bode))
	reg.Register("em.signal.feedback_stability", edakit.Func(feedbackStability))
	reg.Register("em.signal.z_transform", edakit.Func(zTransform))
	reg.Register("em.signal.moment_matching", edakit.Func(momentMatching))
}

func ltiConvolution(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a, b := edakit.Int(-3, 4), edakit.Int(-3, 4)
		if a == 0 {
			a = 2
		}
		if b == 0 {
			b = -1
		}
		y1, y2 := rnd(0.1, 0.05, 30), rnd(0.1, 0.05, 30)
		t := edakit.Int(5, 90)
		y := float64(a)*y1 + float64(b)*y2
		return edakit.Numeric(
			sp("A linear time-invariant interconnect model responds to input $x_1(t)$ with an output of %s V at $t = %d$ ps, and to input $x_2(t)$ with %s V at the same time. What is its output at $t = %d$ ps for the input $%d\\,x_1(t) %s\\,x_2(t)$, in volts? Round to 2 decimal places.", num(y1, 2), t, num(y2, 2), t, a, signed(b)),
			num(y, 2),
			sp("Linearity gives superposition with scaling: $%d \\cdot %s %s \\cdot %s = %s$ V.", a, num(y1, 2), signed(b), num(y2, 2), num(y, 2)),
			facts("y", num(y, 2)))
	case 2:
		h := []int{edakit.Int(1, 5), edakit.Int(-2, 4), edakit.Int(-2, 3)}
		nx := edakit.Scaled(ctx, 3, 5)
		x := make([]int, nx)
		for i := range x {
			x[i] = edakit.Int(-3, 5)
		}
		n := edakit.Int(1, nx+1)
		y := 0
		var terms []string
		for k := 0; k < len(h); k++ {
			if j := n - k; j >= 0 && j < nx {
				y += h[k] * x[j]
				terms = append(terms, sp("h[%d]x[%d] = %s \\cdot %s", k, j, paren(h[k]), paren(x[j])))
			}
		}
		return edakit.Numeric(
			sp("A discrete-time system has impulse response $h = [%s]$ starting at $n = 0$. The input is $x = [%s]$ starting at $n = 0$, and zero elsewhere. What is the output $y[%d]$?", ints(h), ints(x), n),
			edakit.IntStr(y),
			sp("$y[n] = \\sum_k h[k]\\,x[n-k]$. For $n = %d$ the nonzero terms are $%s$, which sum to %d.", n, strings.Join(terms, ",\\ "), y),
			facts("y", edakit.IntStr(y)))
	default:
		v := edakit.Pick(supplies)
		tau := float64(5 * edakit.Int(2, 20))
		w := float64(5 * edakit.Int(1, 20))
		t := w + float64(5*edakit.Int(1, 20))
		y := v * (math.Exp(-(t-w)/tau) - math.Exp(-t/tau))
		return edakit.Numeric(
			sp("A single RC stage with $\\tau = %s$ ps, initially at 0 V, is driven by a rectangular pulse of %s V that lasts from $t = 0$ to $t = %s$ ps. What is the output voltage at $t = %s$ ps, in volts? Round to 3 decimal places.", num(tau, 0), num(v, 1), num(w, 0), num(t, 0)),
			num(y, 3),
			sp("The pulse is a step at 0 minus a step at %s ps, so by time invariance the output is $V(1 - e^{-t/\\tau}) - V(1 - e^{-(t-%s)/\\tau}) = %s(e^{-%s/%s} - e^{-%s/%s}) = %s$ V.", num(w, 0), num(w, 0), num(v, 1), num(t-w, 0), num(tau, 0), num(t, 0), num(tau, 0), num(y, 3)),
			facts("y", num(y, 3)))
	}
}

func paren(x int) string {
	if x < 0 {
		return sp("(%d)", x)
	}
	return edakit.IntStr(x)
}

func signed(b int) string {
	if b < 0 {
		return sp("- %d", -b)
	}
	return sp("+ %d", b)
}

func ints(xs []int) string {
	var s []string
	for _, x := range xs {
		s = append(s, edakit.IntStr(x))
	}
	return strings.Join(s, ", ")
}

func fourierSpectrum(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		v := edakit.Pick(supplies)
		f := pick(10, 25, 50, 100, 125, 200, 250)
		n := edakit.Int(1, 9)
		a := 0.0
		if n%2 == 1 {
			a = 2 * v / (float64(n) * math.Pi)
		}
		expl := sp("A 0 to V square wave at 50%% duty has only odd harmonics, of amplitude $2V/(n\\pi)$. %s MHz is harmonic $n = %d$", num(f*float64(n), 0), n)
		if n%2 == 0 {
			expl += ", which is even, so its amplitude is 0."
		} else {
			expl += sp(", so its amplitude is $2 \\cdot %s/(%d\\pi) = %s$ V.", num(v, 1), n, num(a, 3))
		}
		return edakit.Numeric(
			sp("A clock is an ideal square wave from 0 to %s V at %s MHz with 50%% duty cycle. What is the amplitude of its spectral component at %s MHz, in volts? Round to 3 decimal places.", num(v, 1), num(f, 0), num(f*float64(n), 0)),
			num(a, 3), expl, facts("harmonic", edakit.IntStr(n)))
	case 2:
		tr := float64(5 * edakit.Int(2, 100))
		fk := 500 / tr
		return edakit.Numeric(
			sp("A digital edge has a 10%% to 90%% rise time of %s ps. Using the knee-frequency rule of thumb $f_{knee} = 0.5/t_r$, up to what frequency must an interconnect model be accurate, in GHz? Round to 2 decimal places.", num(tr, 0)),
			num(fk, 2),
			sp("$f_{knee} = 0.5/(%s \\times 10^{-12}\\text{ s}) = %s$ GHz. Most of the edge's energy lies below it.", num(tr, 0), num(fk, 2)),
			facts("fknee_ghz", num(fk, 2)))
	default:
		tr := float64(5 * edakit.Int(2, 60))
		pw := tr * float64(edakit.Int(4, 40))
		if edakit.Int(0, 1) == 0 {
			f2 := 1000 / (math.Pi * tr)
			return edakit.Numeric(
				sp("A trapezoidal pulse train has pulse width %s ps and rise time %s ps. Its spectral envelope falls at 20 dB per decade above the first corner and at 40 dB per decade above the second. Where is the second corner, in GHz? Round to 3 decimal places.", num(pw, 0), num(tr, 0)),
				num(f2, 3),
				sp("The corners are $1/(\\pi\\tau)$ for the pulse width and $1/(\\pi t_r)$ for the rise time. The second is $1/(\\pi \\cdot %s\\text{ ps}) = %s$ GHz.", num(tr, 0), num(f2, 3)),
				facts("f2_ghz", num(f2, 3)))
		}
		f1 := 1000 / (math.Pi * pw)
		return edakit.Numeric(
			sp("A trapezoidal pulse train has pulse width %s ps and rise time %s ps. Its spectral envelope is flat at low frequency and starts falling at 20 dB per decade at the first corner. Where is the first corner, in GHz? Round to 3 decimal places.", num(pw, 0), num(tr, 0)),
			num(f1, 3),
			sp("The first corner is set by the pulse width: $1/(\\pi\\tau) = 1/(\\pi \\cdot %s\\text{ ps}) = %s$ GHz. The rise time sets the second, $1/(\\pi t_r) = %s$ GHz.", num(pw, 0), num(f1, 3), num(1000/(math.Pi*tr), 3)),
			facts("f1_ghz", num(f1, 3)))
	}
}

func laplaceTransfer(ctx generator.GeneratorContext) generator.Problem {
	r1, r2 := edakit.Pick(kohms), edakit.Pick(kohms)
	c := float64(5 * edakit.Int(1, 40))
	h0 := r2 / (r1 + r2)
	stem := sp("A network has $R_1 = %s$ kohm in series from input to output, and $R_2 = %s$ kohm in parallel with $C = %s$ fF from output to ground.", num(r1, 1), num(r2, 1), num(c, 0))
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+" What is its DC gain $H(0)$? Round to 4 decimal places.",
			num(h0, 4),
			sp("At $s = 0$ the capacitor's impedance $1/(sC)$ is infinite, so it drops out and $H(0) = R_2/(R_1 + R_2) = %s$.", num(h0, 4)),
			facts("h0", num(h0, 4)))
	case 2:
		tau := c * par(r1, r2)
		return edakit.Numeric(stem+" Writing $H(s) = H(0)/(1 + s\\tau)$, what is $\\tau$, in picoseconds? Round to 2 decimal places.",
			num(tau, 2),
			sp("$H(s) = \\frac{Z_2}{R_1 + Z_2}$ with $Z_2 = R_2/(1 + sR_2C)$, which simplifies to $\\frac{R_2/(R_1+R_2)}{1 + sC(R_1 \\parallel R_2)}$. So $\\tau = %s \\cdot %s = %s$ ps.", num(c, 0), num(par(r1, r2), 4), num(tau, 2)),
			facts("tau_ps", num(tau, 2)))
	default:
		c1, c2 := float64(5*edakit.Int(1, 30)), float64(5*edakit.Int(1, 30))
		v := edakit.Pick(supplies)
		st := sp("A probe attenuator has $R_1 = %s$ kohm in parallel with $C_1 = %s$ fF in the series arm, and $R_2 = %s$ kohm in parallel with $C_2 = %s$ fF from output to ground. A %s V step is applied at $t = 0$.", num(r1, 1), num(c1, 0), num(r2, 1), num(c2, 0), num(v, 1))
		if edakit.Int(0, 1) == 0 {
			x := v * c1 / (c1 + c2)
			return edakit.Numeric(st+" What is the output voltage just after the step, at $t = 0^+$, in volts? Round to 3 decimal places.",
				num(x, 3),
				sp("By the initial value theorem the output starts at $\\lim_{s\\to\\infty} H(s) \\cdot V$. At high frequency the capacitors dominate: $C_1/(C_1 + C_2) \\cdot %s = %s$ V.", num(v, 1), num(x, 3)),
				facts("v0", num(x, 3), "vinf", num(v*h0, 3)))
		}
		x := v * h0
		return edakit.Numeric(st+" What voltage does the output settle to, in volts? Round to 3 decimal places.",
			num(x, 3),
			sp("By the final value theorem the output settles to $H(0) \\cdot V = R_2/(R_1+R_2) \\cdot %s = %s$ V. It starts at $C_1/(C_1+C_2) \\cdot V = %s$ V; the two match only when $R_1 C_1 = R_2 C_2$.", num(v, 1), num(x, 3), num(v*c1/(c1+c2), 3)),
			facts("v0", num(v*c1/(c1+c2), 3), "vinf", num(x, 3)))
	}
}

func polesZeros(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		if edakit.Int(0, 1) == 0 {
			p1, p2 := edakit.Int(1, 12), edakit.Int(1, 12)
			if p1 == p2 {
				p2 += 3
			}
			a, b := p1+p2, p1*p2
			lo := p1
			if p2 < lo {
				lo = p2
			}
			return edakit.Numeric(
				sp("A transfer function has denominator $s^2 + %ds + %d$, with $s$ in units of $10^9$ per second. What is the pole closest to the origin, in those units?", a, b),
				edakit.IntStr(-lo),
				sp("The roots of $s^2 + %ds + %d$ are $s = -%d$ and $s = -%d$. The closer one, $-%d$, is the dominant (slowest) pole.", a, b, p1, p2, lo),
				facts("p1", edakit.IntStr(-p1), "p2", edakit.IntStr(-p2)))
		}
		sg, w := edakit.Int(1, 9), edakit.Int(1, 12)
		return edakit.Numeric(
			sp("A transfer function has denominator $s^2 + %ds + %d$, with $s$ in units of $10^9$ per second. Its poles are a complex pair. What is their real part, in those units?", 2*sg, sg*sg+w*w),
			edakit.IntStr(-sg),
			sp("$s = -%d \\pm \\sqrt{%d^2 - %d} = -%d \\pm j%d$. The real part is half the negated linear coefficient, $-%d$.", sg, sg, sg*sg+w*w, sg, w, sg),
			facts("sigma", edakit.IntStr(-sg), "omega", edakit.IntStr(w)))
	case 2:
		p := edakit.Int(-6, 8)
		if p == 0 {
			p = 2
		}
		b := edakit.Int(-6, 8)
		if b == 0 {
			b = -1
		}
		c := edakit.Int(-5, 20)
		if c == 0 {
			c = 3
		}
		rhp := 0
		if p < 0 {
			rhp++
		}
		switch {
		case c < 0:
			rhp++
		case b < 0:
			rhp += 2
		}
		return edakit.Numeric(
			sp("A system's denominator factors as $(s %s)(s^2 %s s %s)$. How many of its poles lie in the open right half-plane?", signed(p), signedCoef(b), signed(c)),
			edakit.IntStr(rhp),
			sp("The first factor has its pole at $s = %d$. For $s^2 + bs + c$: if $c < 0$ the roots have opposite signs (one in the right half-plane); if $c > 0$ both roots share the sign of $-b$. Here that gives %d right-half-plane pole(s) in total, so the system is %s.", -p, rhp, stability(rhp)),
			facts("rhp", edakit.IntStr(rhp)))
	default:
		sg := float64(edakit.Int(1, 10))
		w := float64(edakit.Int(1, 15))
		wn := math.Hypot(sg, w)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A second-order system has poles at $s = -%s \\pm j%s$ (units of $10^9$ per second). What is its damping ratio $\\zeta$? Round to 3 decimal places.", num(sg, 0), num(w, 0)),
				num(sg/wn, 3),
				sp("$\\omega_n$ is the poles' distance from the origin, $\\sqrt{%s^2 + %s^2} = %s$, and $\\zeta = -\\sigma/\\omega_n = %s/%s = %s$.", num(sg, 0), num(w, 0), num(wn, 3), num(sg, 0), num(wn, 3), num(sg/wn, 3)),
				facts("wn", num(wn, 3)))
		}
		return edakit.Numeric(
			sp("A second-order system has poles at $s = -%s \\pm j%s$ (units of $10^9$ per second). What is its undamped natural frequency $\\omega_n$, in the same units? Round to 3 decimal places.", num(sg, 0), num(w, 0)),
			num(wn, 3),
			sp("$\\omega_n = |s| = \\sqrt{%s^2 + %s^2} = %s$.", num(sg, 0), num(w, 0), num(wn, 3)),
			facts("wn", num(wn, 3)))
	}
}

func signedCoef(b int) string {
	switch b {
	case 1:
		return "+"
	case -1:
		return "-"
	}
	return signed(b)
}

func stability(rhp int) string {
	if rhp == 0 {
		return "stable"
	}
	return "unstable"
}

func bode(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		k := float64(edakit.Int(1, 100))
		fp := pick(1, 2, 5, 10, 20, 50, 100)
		f := fp * pick(0.1, 0.5, 1, 2, 3, 5, 10, 20, 100)
		db := 20*math.Log10(k) - 10*math.Log10(1+(f/fp)*(f/fp))
		return edakit.Numeric(
			sp("An amplifier has $H(f) = %s/(1 + jf/f_p)$ with $f_p = %s$ MHz. What is its exact gain at %s MHz, in dB? Round to 2 decimal places.", num(k, 0), num(fp, 0), num(f, 1)),
			num(db, 2),
			sp("$|H| = %s/\\sqrt{1 + (%s/%s)^2}$, so $20\\log_{10}|H| = %s - 10\\log_{10}(1 + %s) = %s$ dB.", num(k, 0), num(f, 1), num(fp, 0), num(20*math.Log10(k), 2), num((f/fp)*(f/fp), 4), num(db, 2)),
			facts("k_db", num(20*math.Log10(k), 2)))
	case 2:
		kdb := float64(10 * edakit.Int(2, 10))
		p1 := pick(1, 2, 5, 10)
		p2 := p1 * pick(10, 20, 50, 100)
		f := p2 * pick(2, 5, 10, 20)
		db := kdb - 20*math.Log10(f/p1) - 20*math.Log10(f/p2)
		return edakit.Numeric(
			sp("A transfer function has a DC gain of %s dB and real poles at %s MHz and %s MHz. Using the straight-line Bode asymptotes, what is the magnitude at %s MHz, in dB? Round to 2 decimal places.", num(kdb, 0), num(p1, 0), num(p2, 0), num(f, 0)),
			num(db, 2),
			sp("Each pole subtracts 20 dB per decade above its frequency: $%s - 20\\log_{10}(%s/%s) - 20\\log_{10}(%s/%s) = %s$ dB.", num(kdb, 0), num(f, 0), num(p1, 0), num(f, 0), num(p2, 0), num(db, 2)),
			facts("db", num(db, 2)))
	default:
		p1 := pick(1, 2, 5, 10, 20)
		p2 := p1 * pick(3, 5, 10, 20)
		f := p1 * pick(0.5, 1, 2, 4, 5, 10, 20)
		ph := -deg(math.Atan(f/p1)) - deg(math.Atan(f/p2))
		return edakit.Numeric(
			sp("A transfer function has real left-half-plane poles at %s MHz and %s MHz and no zeros. What is its phase at %s MHz, in degrees? Round to 1 decimal place.", num(p1, 0), num(p2, 0), num(f, 1)),
			num(ph, 1),
			sp("Each pole contributes $-\\arctan(f/f_p)$: $-\\arctan(%s/%s) - \\arctan(%s/%s) = %s$ degrees.", num(f, 1), num(p1, 0), num(f, 1), num(p2, 0), num(ph, 1)),
			facts("phase", num(ph, 1)))
	}
}

func feedbackStability(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := pick(100, 200, 500, 1000, 5000, 10000, 100000)
		b := pick(0.01, 0.02, 0.05, 0.1, 0.2, 0.25, 0.5)
		g := a / (1 + a*b)
		return edakit.Numeric(
			sp("An amplifier with open-loop gain %s is wrapped in negative feedback with feedback factor $\\beta = %s$. What is the closed-loop gain? Round to 3 decimal places.", num(a, 0), num(b, 2)),
			num(g, 3),
			sp("$A_{cl} = A/(1 + A\\beta) = %s/(1 + %s) = %s$, close to $1/\\beta = %s$ because the loop gain $A\\beta = %s$ is large.", num(a, 0), num(a*b, 2), num(g, 3), num(1/b, 2), num(a*b, 2)),
			facts("loop_gain", num(a*b, 2)))
	case 2:
		t0 := float64(10 * edakit.Int(2, 100))
		fp := pick(1, 2, 5, 10, 20, 50)
		fc := fp * math.Sqrt(t0*t0-1) / 1000
		return edakit.Numeric(
			sp("A feedback loop has loop gain $T(f) = %s/(1 + jf/f_p)$ with $f_p = %s$ kHz. At what frequency does $|T|$ fall to 1, in MHz? Round to 3 decimal places.", num(t0, 0), num(fp, 0)),
			num(fc, 3),
			sp("$|T| = 1$ when $1 + (f/f_p)^2 = T_0^2$, so $f_c = f_p\\sqrt{T_0^2 - 1} \\approx T_0 f_p = %s$ MHz.", num(fc, 3)),
			facts("fc_mhz", num(fc, 3)))
	default:
		t0 := float64(10 * edakit.Int(2, 100))
		p1 := pick(0.01, 0.02, 0.05, 0.1)
		p2 := pick(1, 2, 5, 10, 20, 50)
		mag := func(f float64) float64 {
			return t0 / math.Sqrt((1+(f/p1)*(f/p1))*(1+(f/p2)*(f/p2)))
		}
		lo, hi := 0.0, t0*p1*10
		for i := 0; i < 200; i++ {
			m := (lo + hi) / 2
			if mag(m) > 1 {
				lo = m
			} else {
				hi = m
			}
		}
		fc := (lo + hi) / 2
		pm := 180 - deg(math.Atan(fc/p1)) - deg(math.Atan(fc/p2))
		return edakit.Numeric(
			sp("A loop gain is $T(f) = %s/((1 + jf/f_1)(1 + jf/f_2))$ with $f_1 = %s$ MHz and $f_2 = %s$ MHz. What is the phase margin, in degrees? Round to 1 decimal place.", num(t0, 0), num(p1, 2), num(p2, 0)),
			num(pm, 1),
			sp("Solve $|T(f_c)| = 1$ numerically: $f_c = %s$ MHz. The phase there is $-\\arctan(f_c/f_1) - \\arctan(f_c/f_2) = %s$ degrees, so the margin to $-180$ is %s degrees.", num(fc, 3), num(pm-180, 1), num(pm, 1)),
			facts("fc_mhz", num(fc, 3), "pm", num(pm, 1)))
	}
}

func zTransform(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := pick(0.5, 0.6, 0.8, 0.9, -0.5, -0.8)
		b := pick(0.1, 0.2, 0.25, 0.5, 1)
		k := edakit.Int(2, 6)
		y := 0.0
		var seq []string
		for n := 0; n <= k; n++ {
			y = a*y + b
			seq = append(seq, num(y, 4))
		}
		return edakit.Numeric(
			sp("A discrete-time filter obeys $y[n] = %s\\,y[n-1] + %s\\,x[n]$ with $y[-1] = 0$. The input is a unit step, $x[n] = 1$ for $n \\ge 0$. What is $y[%d]$? Round to 4 decimal places.", num(a, 2), num(b, 2), k),
			num(y, 4),
			sp("Iterate from $n = 0$: %s. The pole of $H(z) = %s/(1 - %s z^{-1})$ is $z = %s$, inside the unit circle, so the output converges to $%s/(1 - %s) = %s$.", strings.Join(seq, ", "), num(b, 2), num(a, 2), num(a, 2), num(b, 2), num(a, 2), num(b/(1-a), 4)),
			facts("y", num(y, 4)))
	case 2:
		tau := float64(5 * edakit.Int(1, 40))
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A simulator integrates $dv/dt = -v/\\tau$ with $\\tau = %s$ ps using forward Euler. What is the largest time step, in picoseconds, for which the numerical solution does not grow?", num(tau, 0)),
				num(2*tau, 0),
				sp("Forward Euler gives $v_{n+1} = (1 - h/\\tau) v_n$, a pole at $z = 1 - h/\\tau$. Stability needs $|z| \\le 1$, so $h \\le 2\\tau = %s$ ps.", num(2*tau, 0)),
				facts("hmax", num(2*tau, 0)))
		}
		h := tau * pick(0.1, 0.25, 0.5, 1, 1.5, 2.5, 3)
		z := 1 - h/tau
		return edakit.Numeric(
			sp("A simulator integrates $dv/dt = -v/\\tau$ with $\\tau = %s$ ps using forward Euler at a step of %s ps. What is the pole $z$ of the resulting difference equation? Round to 3 decimal places.", num(tau, 0), num(h, 2)),
			num(z, 3),
			sp("$v_{n+1} = v_n + h(-v_n/\\tau) = (1 - h/\\tau)v_n$, so $z = 1 - %s/%s = %s$; $|z| > 1$ means the solution blows up.", num(h, 2), num(tau, 0), num(z, 3)),
			facts("z", num(z, 3)))
	default:
		tau := float64(5 * edakit.Int(1, 40))
		h := tau * pick(0.1, 0.25, 0.5, 1, 2, 5, 10)
		if edakit.Int(0, 1) == 0 {
			z := 1 / (1 + h/tau)
			return edakit.Numeric(
				sp("A simulator integrates $dv/dt = -v/\\tau$ with $\\tau = %s$ ps using backward Euler at a step of %s ps. What is the pole $z$ of the resulting difference equation? Round to 4 decimal places.", num(tau, 0), num(h, 2)),
				num(z, 4),
				sp("$v_{n+1} = v_n - (h/\\tau)v_{n+1}$ gives $z = 1/(1 + h/\\tau) = %s$, inside the unit circle for every positive step.", num(z, 4)),
				facts("z", num(z, 4)))
		}
		z := (1 - h/(2*tau)) / (1 + h/(2*tau))
		return edakit.Numeric(
			sp("A simulator integrates $dv/dt = -v/\\tau$ with $\\tau = %s$ ps using the trapezoidal rule at a step of %s ps. What is the pole $z$ of the resulting difference equation? Round to 4 decimal places.", num(tau, 0), num(h, 2)),
			num(z, 4),
			sp("Trapezoidal integration maps $s$ to $z = (1 + hs/2)/(1 - hs/2)$; with $s = -1/\\tau$, $z = (1 - %s)/(1 + %s) = %s$. A large step pushes $z$ toward $-1$, the source of trapezoidal ringing.", num(h/(2*tau), 4), num(h/(2*tau), 4), num(z, 4)),
			facts("z", num(z, 4)))
	}
}

func momentMatching(ctx generator.GeneratorContext) generator.Problem {
	r1, r2 := edakit.Pick(kohms), edakit.Pick(kohms)
	c1, c2 := float64(5*edakit.Int(1, 20)), float64(5*edakit.Int(1, 20))
	m1n1 := r1 * (c1 + c2)
	m1n2 := r1*c1 + (r1+r2)*c2
	m2n2 := r1*c1*m1n1 + (r1+r2)*c2*m1n2
	stem := sp("An RC chain has $R_1 = %s$ kohm from an ideal source to node 1 ($C_1 = %s$ fF to ground) and $R_2 = %s$ kohm from node 1 to node 2 ($C_2 = %s$ fF to ground). Expand the transfer function to node $i$ as $H_i(s) = m_0 + m_1 s + m_2 s^2 + \\cdots$.", num(r1, 1), num(c1, 0), num(r2, 1), num(c2, 0))
	f := facts("m1_node1", num(m1n1, 2), "m1_node2", num(m1n2, 2), "m2_node2", num(m2n2, 1))
	switch kp(ctx) {
	case 1:
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(stem+" What is $-m_1$ at node 1, in picoseconds? Round to 2 decimal places.",
				num(m1n1, 2),
				sp("$-m_1$ at a node is $\\sum_k R_{ik} C_k$ (the Elmore delay). Node 1 shares only $R_1$ with both capacitors: $%s(%s + %s) = %s$ ps.", num(r1, 1), num(c1, 0), num(c2, 0), num(m1n1, 2)), f)
		}
		return edakit.Numeric(stem+" What is $-m_1$ at node 2, in picoseconds? Round to 2 decimal places.",
			num(m1n2, 2),
			sp("$-m_1 = \\sum_k R_{2k}C_k = R_1 C_1 + (R_1 + R_2) C_2 = %s$ ps, the Elmore delay to node 2.", num(m1n2, 2)), f)
	case 2:
		return edakit.Numeric(stem+" What is $m_2$ at node 2, in square picoseconds? Round to 1 decimal place.",
			num(m2n2, 1),
			sp("The recursion is $m_2^{(i)} = -\\sum_k R_{ik} C_k\\, m_1^{(k)}$, with $-m_1^{(1)} = %s$ and $-m_1^{(2)} = %s$ ps. So $m_2^{(2)} = %s \\cdot %s \\cdot %s + %s \\cdot %s \\cdot %s = %s$ square ps.", num(m1n1, 2), num(m1n2, 2), num(r1, 1), num(c1, 0), num(m1n1, 2), num(r1+r2, 1), num(c2, 0), num(m1n2, 2), num(m2n2, 1)), f)
	default:
		sd := math.Sqrt(2*m2n2 - m1n2*m1n2)
		return edakit.Numeric(stem+" Treating the impulse response at node 2 as a distribution, what is its standard deviation, in picoseconds? Round to 2 decimal places.",
			num(sd, 2),
			sp("The moments satisfy $\\int t\\,h = -m_1$ and $\\int t^2 h = 2m_2$. With $-m_1 = %s$ ps and $m_2 = %s$ square ps, the variance is $2m_2 - m_1^2 = %s$, so the standard deviation is %s ps.", num(m1n2, 2), num(m2n2, 1), num(2*m2n2-m1n2*m1n2, 1), num(sd, 2)), f)
	}
}
