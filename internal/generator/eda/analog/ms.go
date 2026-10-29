package analog

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// vaCompute evaluates a small Verilog-A contribution for a drawn operating point.
func vaCompute(ctx generator.GeneratorContext) generator.Problem {
	switch edakit.Int(0, 2) {
	case 0:
		vs := unif(0.4, 1.2, 0.1)
		g := float64(edakit.Int(5, 100))
		vin := unif(0.001, 0.05, 0.001)
		out := vs * math.Tanh(g*vin/vs)
		q := edakit.Sprintf("A Verilog-A amplifier model contains `V(out) <+ vsat * tanh(gain * V(in) / vsat);` with `vsat = %s` and `gain = %s`. At a DC operating point with $V(in) = %s$ V, what does the simulator compute for $V(out)$?", f(vs, 1), f(g, 0), f(vin, 3))
		return choiceNums(q, out, 4, " V", edakit.Sprintf("$%s \\tanh(%s \\cdot %s/%s) = %s$ V. The tanh keeps the model smooth and bounded, which helps Newton converge.", f(vs, 1), f(g, 0), f(vin, 3), f(vs, 1), f(out, 4)),
			math.Min(g*vin, 99), vs, math.Tanh(g*vin), g*vin/vs)
	case 1:
		r := float64(edakit.Int(1, 100))
		c := unif(0.1, 10, 0.1)
		v := unif(0.1, 2, 0.1)
		dv := float64(edakit.Int(1, 50))
		i := v/r*1000 + c*dv
		q := edakit.Sprintf("A Verilog-A branch is `I(p,n) <+ V(p,n)/r + c*ddt(V(p,n));` where the parameter r is %s kilohms and c is %s pF. At an instant in a transient run, $V(p,n) = %s$ V and its slope is %s V/$\\mu$s. What branch current flows, in microamperes?", f(r, 0), f(c, 1), f(v, 1), f(dv, 0))
		return choiceNums(q, i, 3, " µA", edakit.Sprintf("The resistive term gives $%s/%s\\ \\mathrm{k\\Omega} = %s\\ \\mu$A and the `ddt` term gives $C\\,dV/dt = %s\\ \\mathrm{pF} \\cdot %s\\ \\mathrm{V}/\\mu\\mathrm{s} = %s\\ \\mu$A; contributions to one branch add, for a total of %s.", f(v, 1), f(r, 0), f(v/r*1000, 3), f(c, 1), f(dv, 0), f(c*dv, 3), f(i, 3)),
			v/r*1000, c*dv, v/r*1000-c*dv)
	default:
		g := float64(edakit.Int(2, 50))
		vos := unif(0.1, 5, 0.1)
		vin := unif(-0.05, 0.05, 0.001)
		vdd := unif(0.9, 1.8, 0.1)
		raw := g * (vin + vos/1000)
		out := math.Max(0, math.Min(vdd, vdd/2+raw))
		q := edakit.Sprintf("A behavioral op-amp model computes `vo = vdd/2 + gain*(V(inp,inn) + vos)` and then clamps it to the rails with `min(max(vo, 0), vdd)` before `V(out) <+ vo;`. With `gain = %s`, `vos = %s` mV, `vdd = %s` V and $V(inp,inn) = %s$ V, what is $V(out)$?", f(g, 0), f(vos, 1), f(vdd, 1), f(vin, 3))
		return choiceNums(q, out, 4, " V", edakit.Sprintf("$%s/2 + %s(%s + %s) = %s$ V before clamping; clamping to $[0, %s]$ gives %s V. A hard clamp has a kink in its derivative, which a smooth limiter avoids.", f(vdd, 1), f(g, 0), f(vin, 3), f(vos/1000, 4), f(vdd/2+raw, 4), f(vdd, 1), f(out, 4)),
			vdd/2+g*vin, vdd/2+raw+0.5, raw, vdd/2)
	}
}

func realNumber(ctx generator.GeneratorContext) generator.Problem {
	tau := float64(edakit.Int(5, 200)) // ns
	t := float64(edakit.Int(1, 20))    // ns
	alpha := 1 - math.Exp(-t/tau)
	switch whichKP(ctx) {
	case 1:
		q := edakit.Sprintf("A real-number model of an RC low-pass filter with $\\tau = %s$ ns is updated every $T = %s$ ns by the exact discrete step $y_{k+1} = y_k + \\alpha (x_k - y_k)$. What is $\\alpha$? Round to 4 decimal places.", f(tau, 0), f(t, 0))
		return num(q, f(alpha, 4), edakit.Sprintf("Holding $x$ constant over a step, the ODE $\\tau\\dot{y} = x - y$ solves exactly to $\\alpha = 1 - e^{-T/\\tau} = 1 - e^{-%s} = %s$.", f(t/tau, 4), f(alpha, 4)), nil)
	case 2:
		k := edakit.Int(2, 20)
		x := unif(0.5, 1.8, 0.1)
		y0 := unif(0, 0.4, 0.1)
		yk := x + (y0-x)*math.Exp(-float64(k)*t/tau)
		q := edakit.Sprintf("An event-driven real-number model of an RC filter ($\\tau = %s$ ns) starts at $y_0 = %s$ V, receives a step to $x = %s$ V, and updates every %s ns with the exact exponential step. What is $y$ after %d updates, in volts? Round to 4 decimal places.", f(tau, 0), f(y0, 1), f(x, 1), f(t, 0), k)
		return num(q, f(yk, 4), edakit.Sprintf("Repeated exact steps compose: $y_k = x + (y_0 - x)e^{-kT/\\tau} = %s + (%s - %s)e^{-%s} = %s$ V.", f(x, 1), f(y0, 1), f(x, 1), f(float64(k)*t/tau, 4), f(yk, 4)), nil)
	default:
		k := edakit.Int(2, 10)
		x := unif(0.5, 1.8, 0.1)
		tt := float64(edakit.Int(1, 20))
		tau2 := tt * unif(1.5, 8, 0.5)
		ex := x * (1 - math.Exp(-float64(k)*tt/tau2))
		eu := x * (1 - math.Pow(1-tt/tau2, float64(k)))
		d := math.Abs(ex-eu) * 1000
		q := edakit.Sprintf("Two real-number models of an RC filter with $\\tau = %s$ ns start at 0 V, see a step to %s V, and update every %s ns. One uses the exact step $\\alpha = 1 - e^{-T/\\tau}$; the other uses forward Euler, $\\alpha = T/\\tau$. How far apart are their outputs after %d updates, in millivolts? Round to 3 decimal places.", f(tau2, 1), f(x, 1), f(tt, 0), k)
		return num(q, f(d, 3), edakit.Sprintf("Exact: $%s(1 - e^{-%s}) = %s$ V. Euler: $%s(1 - (1 - %s)^{%d}) = %s$ V. The difference is %s mV; Euler rises faster than the exact response because its per-step factor $T/\\tau$ exceeds $1 - e^{-T/\\tau}$.", f(x, 1), f(float64(k)*tt/tau2, 4), f(ex, 5), f(x, 1), f(tt/tau2, 4), k, f(eu, 5), f(d, 3)), nil)
	}
}
