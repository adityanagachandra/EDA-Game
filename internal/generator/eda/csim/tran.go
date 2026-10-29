package csim

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var (
	capPF  = []float64{0.5, 1, 2, 4, 5, 10}
	stepNS = []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1}
	indUH  = []float64{0.1, 0.2, 0.5, 1, 2, 5, 10}
)

// companionC asks for backward-Euler and trapezoidal capacitor companion quantities.
func companionC(ctx generator.GeneratorContext) generator.Problem {
	c := edakit.Pick(capPF)
	h := edakit.Pick(stepNS)
	vn := stepF(-1, 0.1, 40)
	in := stepF(-2, 0.1, 40)
	conv := "Units: C/h with C in pF and h in ns is in mS; currents are in mA. The companion is $i_C(t_{n+1}) = G_{eq} v_{n+1} - I_{eq}$."
	switch kp(ctx) {
	case 1:
		geq := c / h
		if edakit.Int(0, 1) == 0 {
			return num(edakit.Sprintf("A %s pF capacitor is integrated with backward Euler at step $h = %s$ ns. %s What is $G_{eq}$ in mS?", f(c, 2), f(h, 2), conv), geq, 2,
				edakit.Sprintf("Backward Euler replaces $C\\,dv/dt$ by $C(v_{n+1} - v_n)/h$, so $G_{eq} = C/h = %s$ mS.", f(geq, 2)), nil)
		}
		ieq := geq * vn
		return num(edakit.Sprintf("A %s pF capacitor is integrated with backward Euler at step $h = %s$ ns, and its voltage at $t_n$ is %s V. %s What is the history source $I_{eq}$ in mA?", f(c, 2), f(h, 2), f(vn, 2), conv), ieq, 2,
			edakit.Sprintf("$I_{eq} = (C/h) v_n = %s \\times %s = %s$ mA.", f(geq, 2), f(vn, 2), f(ieq, 2)), nil)
	case 2:
		geq := 2 * c / h
		if edakit.Int(0, 1) == 0 {
			return num(edakit.Sprintf("A %s pF capacitor is integrated with the trapezoidal rule at step $h = %s$ ns. %s What is $G_{eq}$ in mS?", f(c, 2), f(h, 2), conv), geq, 2,
				edakit.Sprintf("The trapezoidal rule averages the currents at both ends of the step, giving $G_{eq} = 2C/h = %s$ mS.", f(geq, 2)), nil)
		}
		ieq := geq*vn + in
		return num(edakit.Sprintf("A %s pF capacitor is integrated with the trapezoidal rule at step $h = %s$ ns. At $t_n$ its voltage is %s V and its current is %s mA. %s What is $I_{eq}$ in mA?", f(c, 2), f(h, 2), f(vn, 2), f(in, 2), conv), ieq, 2,
			edakit.Sprintf("$I_{eq} = (2C/h) v_n + i_n = %s \\times %s + %s = %s$ mA.", f(geq, 2), f(vn, 2), f(in, 2), f(ieq, 2)), nil)
	}
	vn1 := vn + stepF(-0.5, 0.05, 20)
	geq := 2 * c / h
	i1 := geq*(vn1-vn) - in
	return num(edakit.Sprintf("A %s pF capacitor is integrated with the trapezoidal rule at $h = %s$ ns. At $t_n$ it had $v_n = %s$ V and $i_n = %s$ mA; the solve at $t_{n+1}$ gives $v_{n+1} = %s$ V. What capacitor current $i_{n+1}$ must the simulator store for the next step, in mA?", f(c, 2), f(h, 2), f(vn, 2), f(in, 2), f(vn1, 2)), i1, 2,
		edakit.Sprintf("From $v_{n+1} = v_n + \\tfrac{h}{2C}(i_n + i_{n+1})$: $i_{n+1} = (2C/h)(v_{n+1} - v_n) - i_n = %s \\times %s - (%s) = %s$ mA.", f(geq, 2), f(vn1-vn, 2), f(in, 2), f(i1, 2)), nil)
}

// companionL asks for inductor companion quantities.
func companionL(ctx generator.GeneratorContext) generator.Problem {
	l := edakit.Pick(indUH)
	h := edakit.Pick(stepNS)
	in := stepF(-5, 0.5, 20)
	vn := stepF(-2, 0.1, 40)
	conv := "With h in ns and L in uH, h/L is in mS; currents are in mA. The companion is $i_{n+1} = G_{eq} v_{n+1} + I_{eq}$."
	switch kp(ctx) {
	case 1:
		geq := h / l
		if edakit.Int(0, 1) == 0 {
			return num(edakit.Sprintf("A %s uH inductor is integrated with backward Euler at $h = %s$ ns. %s What is $G_{eq}$ in mS?", f(l, 2), f(h, 2), conv), geq, 4,
				edakit.Sprintf("Backward Euler on $v = L\\,di/dt$ gives $i_{n+1} = i_n + (h/L) v_{n+1}$, so $G_{eq} = h/L = %s$ mS and $I_{eq} = i_n$.", f(geq, 4)), nil)
		}
		v1 := stepF(-1, 0.1, 20)
		i1 := in + geq*v1
		return num(edakit.Sprintf("A %s uH inductor carries %s mA at $t_n$. Backward Euler at $h = %s$ ns gives a voltage of %s V across it at $t_{n+1}$. %s What is $i_{n+1}$ in mA?", f(l, 2), f(in, 2), f(h, 2), f(v1, 2), conv), i1, 4,
			edakit.Sprintf("$i_{n+1} = i_n + (h/L) v_{n+1} = %s + %s \\times %s = %s$ mA.", f(in, 2), f(geq, 4), f(v1, 2), f(i1, 4)), nil)
	case 2:
		geq := h / (2 * l)
		return num(edakit.Sprintf("A %s uH inductor is integrated with the trapezoidal rule at $h = %s$ ns. %s What is $G_{eq}$ in mS?", f(l, 2), f(h, 2), conv), geq, 4,
			edakit.Sprintf("The trapezoidal rule gives $i_{n+1} = i_n + \\tfrac{h}{2L}(v_n + v_{n+1})$, so $G_{eq} = h/(2L) = %s$ mS.", f(geq, 4)), nil)
	}
	geq := h / (2 * l)
	ieq := in + geq*vn
	return num(edakit.Sprintf("A %s uH inductor is integrated with the trapezoidal rule at $h = %s$ ns. At $t_n$ it carries %s mA with %s V across it. %s What is the history source $I_{eq}$ in mA?", f(l, 2), f(h, 2), f(in, 2), f(vn, 2), conv), ieq, 4,
		edakit.Sprintf("$I_{eq} = i_n + \\tfrac{h}{2L} v_n = %s + %s \\times %s = %s$ mA.", f(in, 2), f(geq, 4), f(vn, 2), f(ieq, 4)), nil)
}

// rcSetup draws a step-driven RC node: source vs through R (kohm) into C (pF) to ground.
func rcSetup() (vs, r, c, h, tau float64) {
	vs = stepF(1, 0.5, 6)
	r = edakit.Pick([]float64{0.5, 1, 2, 5})
	c = edakit.Pick([]float64{0.5, 1, 2, 5})
	tau = r * c
	h = tau * edakit.Pick([]float64{0.1, 0.2, 0.25, 0.5, 1})
	return
}

func rcStem(vs, r, c, h float64) string {
	return edakit.Sprintf("A source steps to %s V at $t = 0$ and drives R = %s k$\\Omega$ into a node with C = %s pF to ground. The step is $h = %s$ ns (k$\\Omega$ times pF is ns; 1/R is in mS and C/h in mS).", f(vs, 2), f(r, 2), f(c, 2), f(h, 3))
}

// beStep runs backward-Euler steps on the RC node.
func beStep(ctx generator.GeneratorContext) generator.Problem {
	vs, r, c, h, tau := rcSetup()
	g, gc := 1/r, c/h
	step := func(v float64) float64 { return (vs*g + gc*v) / (g + gc) }
	switch kp(ctx) {
	case 1:
		v0 := stepF(0, 0.1, int(vs*10)-1)
		v1 := step(v0)
		return num(rcStem(vs, r, c, h)+edakit.Sprintf(" The capacitor voltage at $t_n$ is %s V. Which value does one backward-Euler step give for $v_{n+1}$, in volts?", f(v0, 2)), v1, 4,
			edakit.Sprintf("$(1/R + C/h) v_{n+1} = V_s/R + (C/h) v_n$: $(%s + %s) v_{n+1} = %s + %s$, so $v_{n+1} = %s$ V.", f(g, 4), f(gc, 4), f(vs*g, 4), f(gc*v0, 4), f(v1, 4)), nil)
	case 2:
		v1 := step(0)
		v2 := step(v1)
		return num(rcStem(vs, r, c, h)+" The capacitor starts at 0 V. What is its voltage after two backward-Euler steps, in volts?", v2, 4,
			edakit.Sprintf("Each step solves $(1/R + C/h) v_{n+1} = V_s/R + (C/h) v_n$. Step 1 gives %s V, step 2 gives %s V. The exact value at $2h$ is %s V.", f(v1, 4), f(v2, 4), f(vs*(1-math.Exp(-2*h/tau)), 4)), nil)
	}
	v1 := step(0)
	exact := vs * (1 - math.Exp(-h/tau))
	e := (exact - v1) * 1e3
	return num(rcStem(vs, r, c, h)+" The capacitor starts at 0 V. By how much does one backward-Euler step undershoot the exact solution $V_s(1 - e^{-h/RC})$ at $t = h$, in millivolts?", e, 2,
		edakit.Sprintf("Backward Euler gives %s V; the exact value is %s V; the difference is %s mV. The error per step is first order in $h$ globally, second order locally.", f(v1, 5), f(exact, 5), f(e, 2)), nil)
}

// trapStep runs trapezoidal steps on the RC node.
func trapStep(ctx generator.GeneratorContext) generator.Problem {
	vs, r, c, h, tau := rcSetup()
	g, gc := 1/r, 2*c/h
	step := func(v, i float64) (float64, float64) {
		v1 := (vs*g + gc*v + i) / (g + gc)
		i1 := gc*(v1-v) - i
		return v1, i1
	}
	k := kp(ctx)
	v0 := 0.0
	if k == 1 {
		v0 = stepF(0, 0.1, int(vs*10)-1)
	}
	i0 := (vs - v0) * g
	v1, i1 := step(v0, i0)
	stem := rcStem(vs, r, c, h) + edakit.Sprintf(" At $t_n$ the capacitor is at %s V, so its current is $i_n = (V_s - v_n)/R = %s$ mA.", f(v0, 2), f(i0, 4))
	expl := edakit.Sprintf("$(1/R + 2C/h) v_{n+1} = V_s/R + (2C/h) v_n + i_n$: $(%s + %s) v_{n+1} = %s + %s + %s$, so $v_{n+1} = %s$ V and $i_{n+1} = (2C/h)(v_{n+1} - v_n) - i_n = %s$ mA.", f(g, 4), f(gc, 4), f(vs*g, 4), f(gc*v0, 4), f(i0, 4), f(v1, 5), f(i1, 5))
	switch k {
	case 1:
		return num(stem+" Which value does one trapezoidal step give for $v_{n+1}$, in volts?", v1, 4, expl, nil)
	case 2:
		return num(stem+" After one trapezoidal step, what capacitor current $i_{n+1}$ is carried forward, in mA?", i1, 4,
			expl+edakit.Sprintf(" It matches the resistor current $(V_s - v_{n+1})/R = %s$ mA, as KCL requires.", f((vs-v1)*g, 5)), nil)
	}
	v2, _ := step(v1, i1)
	return num(stem+" What is the capacitor voltage after two trapezoidal steps, in volts?", v2, 4,
		expl+edakit.Sprintf(" Repeating once more gives %s V; the exact value at $2h$ is %s V.", f(v2, 5), f(vs*(1-math.Exp(-2*h/tau)), 5)), nil)
}

// trapRinging asks about the trapezoidal amplification factor on a fast decaying mode.
func trapRinging(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	tau := pickF(0.01, 0.02, 0.05, 0.1, 0.2)
	ratio := pickF(0.5, 1, 1.5, 3, 4, 5, 10, 20, 40)
	if k == 3 {
		ratio = pickF(3, 4, 5, 10, 20, 40)
	}
	h := tau * ratio
	rf := (1 - ratio/2) / (1 + ratio/2)
	stem := edakit.Sprintf("A capacitor discharges through a resistor with time constant $\\tau = %s$ ns (source-free, $x' = -x/\\tau$). The simulator takes steps of $h = %s$ ns with the trapezoidal rule, whose one-step factor is $(1 - h/2\\tau)/(1 + h/2\\tau)$.", f(tau, 3), f(h, 3))
	switch k {
	case 1:
		return num(stem+" What is the factor multiplying $x$ at each step?", rf, 4,
			edakit.Sprintf("$h/\\tau = %s$, so the factor is $(1 - %s)/(1 + %s) = %s$. It is negative whenever $h > 2\\tau$, which flips the sign every step: ringing.", f(ratio, 3), f(ratio/2, 3), f(ratio/2, 3), f(rf, 4)), nil)
	case 2:
		v0 := stepF(0.5, 0.5, 5)
		m := edakit.Int(2, 5)
		v := v0 * math.Pow(rf, float64(m))
		return num(stem+edakit.Sprintf(" Starting from %s V, what is the computed voltage after %d steps, in volts?", f(v0, 2), m), v, 4,
			edakit.Sprintf("The factor is %s, so $%s \\times (%s)^{%d} = %s$ V. The true solution is %s V, positive and nearly zero.", f(rf, 4), f(v0, 2), f(rf, 4), m, f(v, 4), sci(v0*math.Exp(-float64(m)*ratio))), nil)
	}
	v0 := stepF(0.5, 0.5, 5)
	m := edakit.Int(2, 4)
	be := 1 / (1 + ratio)
	v := v0 * be * math.Pow(rf, float64(m))
	return num(stem+edakit.Sprintf(" To damp the ringing at a breakpoint, the simulator takes the first step with backward Euler (factor $1/(1 + h/\\tau)$) and then %d trapezoidal steps. Starting from %s V, what voltage results, in volts?", m, f(v0, 2)), v, 4,
		edakit.Sprintf("Backward Euler multiplies by %s, then each trapezoidal step by %s: $%s \\times %s \\times (%s)^{%d} = %s$ V. The L-stable first step removes most of the stiff component.", f(be, 4), f(rf, 4), f(v0, 2), f(be, 4), f(rf, 4), m, f(v, 4)), nil)
}

// gearBDF asks for second-order BDF companion quantities.
func gearBDF(ctx generator.GeneratorContext) generator.Problem {
	c := edakit.Pick(capPF)
	h := edakit.Pick(stepNS)
	vn := stepF(0, 0.1, 20)
	vm := vn - stepF(-0.3, 0.05, 12)
	geq := 3 * c / (2 * h)
	ieq := c / (2 * h) * (4*vn - vm)
	formula := "Second-order BDF (Gear 2) at constant step uses $x'_{n+1} \\approx (3x_{n+1} - 4x_n + x_{n-1})/(2h)$, so the capacitor companion is $i_{n+1} = G_{eq} v_{n+1} - I_{eq}$."
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("%s For C = %s pF and $h = %s$ ns, what is $G_{eq}$ in mS?", formula, f(c, 2), f(h, 2)), geq, 2,
			edakit.Sprintf("$G_{eq} = 3C/(2h) = %s$ mS, larger than backward Euler's $C/h$ by half.", f(geq, 2)), nil)
	case 2:
		return num(edakit.Sprintf("%s For C = %s pF and $h = %s$ ns, with $v_n = %s$ V and $v_{n-1} = %s$ V, what is $I_{eq} = \\tfrac{C}{2h}(4v_n - v_{n-1})$ in mA?", formula, f(c, 2), f(h, 2), f(vn, 2), f(vm, 2)), ieq, 2,
			edakit.Sprintf("$C/(2h) = %s$ mS and $4v_n - v_{n-1} = %s$ V, so $I_{eq} = %s$ mA. BDF2 needs two past points, which is why it must be started by a one-step method.", f(c/(2*h), 3), f(4*vn-vm, 2), f(ieq, 2)), nil)
	}
	vs := stepF(1, 0.5, 6)
	r := edakit.Pick([]float64{0.5, 1, 2, 5})
	v1 := (vs/r + ieq) / (1/r + geq)
	return num(edakit.Sprintf("%s A source of %s V drives R = %s k$\\Omega$ into C = %s pF. With $h = %s$ ns, $v_n = %s$ V and $v_{n-1} = %s$ V, what BDF2 value $v_{n+1}$ results, in volts?", formula, f(vs, 2), f(r, 2), f(c, 2), f(h, 2), f(vn, 2), f(vm, 2)), v1, 4,
		edakit.Sprintf("$(1/R + 3C/(2h)) v_{n+1} = V_s/R + I_{eq}$: $(%s + %s) v_{n+1} = %s + %s$, so $v_{n+1} = %s$ V.", f(1/r, 4), f(geq, 3), f(vs/r, 4), f(ieq, 4), f(v1, 4)), nil)
}

// lteEstimate asks for derivative and truncation-error estimates from samples.
func lteEstimate(ctx generator.GeneratorContext) generator.Problem {
	a := stepF(1, 0.5, 4)
	tau := pickF(1, 2, 3, 4, 5)
	h := pickF(0.1, 0.2, 0.25, 0.4, 0.5)
	t0 := pickF(0, 0.5, 1, 2)
	v := make([]float64, 4)
	var parts []string
	for j := range v {
		v[j] = edakit.Round(a*(1-math.Exp(-(t0+float64(j)*h)/tau)), 6)
		parts = append(parts, f(v[j], 6))
	}
	d2 := v[2] - 2*v[1] + v[0]
	d3 := v[3] - 3*v[2] + 3*v[1] - v[0]
	stem := edakit.Sprintf("A node voltage has been accepted at four equally spaced time points, $h = %s$ ns apart: %s V.", f(h, 2), strings.Join(parts, ", "))
	switch kp(ctx) {
	case 1:
		xpp := d2 / (h * h)
		return num(stem+" Using the first three samples, estimate the second derivative $x'' \\approx (v_2 - 2v_1 + v_0)/h^2$ in V/ns$^2$.", xpp, 4,
			edakit.Sprintf("$v_2 - 2v_1 + v_0 = %s$ V; divided by $h^2 = %s$ gives %s V/ns$^2$. This is twice the second divided difference.", f(d2, 6), f(h*h, 4), f(xpp, 4)), nil)
	case 2:
		lte := math.Abs(d2) / 2 * 1e3
		return num(stem+" Backward Euler's local truncation error is about $(h^2/2) x''$. Using $x'' \\approx (v_2 - 2v_1 + v_0)/h^2$ from the first three samples, what is the estimated error magnitude for a step of size $h$, in millivolts?", lte, 3,
			edakit.Sprintf("$(h^2/2) \\cdot (v_2 - 2v_1 + v_0)/h^2 = |%s|/2$ V $= %s$ mV.", f(d2, 6), f(lte, 3)), nil)
	}
	lte := math.Abs(d3) / 12 * 1e3
	return num(stem+" The trapezoidal rule's local truncation error is about $(h^3/12) x'''$. Using $x''' \\approx (v_3 - 3v_2 + 3v_1 - v_0)/h^3$, what is the estimated error magnitude for a step of size $h$, in millivolts?", lte, 4,
		edakit.Sprintf("$v_3 - 3v_2 + 3v_1 - v_0 = %s$ V, so the estimate is $|%s|/12$ V $= %s$ mV. The third difference is six times the third divided difference times $h^3$.", f(d3, 6), f(d3, 6), f(lte, 4)), nil)
}

// timestepControl asks for the next step from an error estimate.
func timestepControl(ctx generator.GeneratorContext) generator.Problem {
	h := pickF(0.1, 0.2, 0.5, 1, 2)
	eps := pickF(0.5, 1, 2)
	switch kp(ctx) {
	case 1:
		lte := eps * pickF(0.05, 0.1, 0.2, 0.4, 0.8, 0.125, 0.027)
		hn := h * math.Cbrt(eps/lte)
		return num(edakit.Sprintf("A trapezoidal step of $h = %s$ ns was accepted with an estimated local truncation error of %s mV against a tolerance of %s mV. Since the trapezoidal error scales as $h^3$, what step $h_{new} = h(\\varepsilon/|LTE|)^{1/3}$ should the simulator try next, in ns?", f(h, 2), f(lte, 4), f(eps, 2)), hn, 4,
			edakit.Sprintf("$\\varepsilon/LTE = %s$; its cube root is %s; $h_{new} = %s$ ns.", f(eps/lte, 4), f(math.Cbrt(eps/lte), 4), f(hn, 4)), nil)
	case 2:
		lte := eps * pickF(0.04, 0.09, 0.16, 0.25, 0.5, 0.64)
		hn := h * math.Sqrt(eps/lte)
		return num(edakit.Sprintf("A backward-Euler step of $h = %s$ ns was accepted with an estimated local truncation error of %s mV against a tolerance of %s mV. Backward Euler's local error scales as $h^2$. What next step $h(\\varepsilon/|LTE|)^{1/2}$ results, in ns?", f(h, 2), f(lte, 4), f(eps, 2)), hn, 4,
			edakit.Sprintf("For order $k$ the local error scales as $h^{k+1}$; here $k = 1$, so take the square root: $\\sqrt{%s} = %s$ and $h_{new} = %s$ ns.", f(eps/lte, 4), f(math.Sqrt(eps/lte), 4), f(hn, 4)), nil)
	}
	ratio := pickF(2, 4, 8, 10, 27, 64)
	lte := eps * ratio
	method := edakit.Pick([]string{"trapezoidal", "backward Euler"})
	p := 3.0
	if method == "backward Euler" {
		p = 2
	}
	hn := 0.9 * h * math.Pow(eps/lte, 1/p)
	return num(edakit.Sprintf("A %s step of $h = %s$ ns produces an estimated local truncation error of %s mV against a tolerance of %s mV. The simulator rejects the step and retries from $t_n$ with $h_{new} = 0.9\\, h(\\varepsilon/|LTE|)^{1/(k+1)}$, where $k$ is the method order. What is the retry step, in ns?", method, f(h, 2), f(lte, 3), f(eps, 2)), hn, 4,
		edakit.Sprintf("The order is %d, so the exponent is $1/%d$: $(1/%s)^{1/%d} = %s$, times $0.9 \\times %s$ gives %s ns. The rejected step's work is discarded.", int(p)-1, int(p), f(ratio, 2), int(p), f(math.Pow(1/ratio, 1/p), 4), f(h, 2), f(hn, 4)), nil)
}

var daeItems = []edakit.Item{
	{KP: 1, Q: "A linear circuit's transient equations are written $G x + C \\dot{x} = b(t)$. What do $G$ and $C$ collect?",
		A: "G holds conductance and incidence stamps; C holds capacitance and inductance stamps", D: []string{"G holds capacitances; C holds conductances", "G is the Jacobian of the sources; C is the source vector", "G holds only resistors to ground; C holds all other elements"},
		Why: "MNA separates memoryless stamps into G and the coefficients of time derivatives (capacitor charges and inductor fluxes) into C."},
	{KP: 1, Q: "A nonlinear simulator writes the equations as $f(x) + \\tfrac{d}{dt} q(x) = b(t)$ rather than with a capacitance matrix. What does $q$ represent?",
		A: "Device charges and fluxes as functions of the unknowns", D: []string{"The branch currents only", "The Newton residual", "The local truncation error"},
		Why: "Charge-oriented formulations integrate charge $q(x)$ directly, so nonlinear capacitances conserve charge, unlike integrating $C(v)\\,dv/dt$."},
	{KP: 1, Q: "Why do accurate simulators prefer integrating a MOSFET's terminal charges rather than its voltage-dependent capacitances?",
		A: "Integrating C(v) dv/dt with a numerical rule can create or destroy charge, while integrating q(v) conserves it", D: []string{"Capacitances cannot be stamped into MNA", "Charges make the matrix symmetric", "Charge formulations need no Newton iteration"},
		Why: "For a nonlinear capacitance, the discretized C(v) dv/dt does not equal the exact change in charge, so charge leaks numerically. A charge-based model avoids that."},
	{KP: 1, Q: "An inductor in MNA contributes which kind of unknown and equation?",
		A: "A branch current unknown with the equation $v_+ - v_- - L\\, di/dt = 0$", D: []string{"No unknown; it stamps only a conductance $L$", "A node voltage unknown", "A charge unknown with $q = L v$"},
		Why: "Like a voltage source, an inductor needs its current as an unknown, and its branch equation contains the time derivative of that current."},
	{KP: 2, Q: "A circuit contains a voltage source and resistors with no capacitor at some nodes. Why is its $C$ matrix singular?",
		A: "Rows for nodes without capacitance and branch rows of sources contain no derivatives, so they are algebraic constraints", D: []string{"Because G is singular too", "Because the timestep is zero", "Because voltage sources make C negative"},
		Why: "A singular C means the system is a differential-algebraic equation: some unknowns are fixed by algebraic equations at every instant rather than by integration."},
	{KP: 2, Q: "Why must a transient analysis start from consistent initial conditions rather than arbitrary node voltages?",
		A: "The algebraic equations must hold at t = 0, so only some initial values are free", D: []string{"Because the integrator is explicit", "Because C must be inverted at t = 0", "It need not; any initial values are fine"},
		Why: "In a DAE the algebraic part constrains the state. SPICE computes a DC operating point first so the starting values satisfy every algebraic equation."},
	{KP: 2, Q: "A loop of capacitors and a voltage source, with no resistance in the loop, is simulated. What numerical difficulty does this create?",
		A: "A higher-index DAE: the capacitor voltages are constrained algebraically, and the currents depend on derivatives of the source", D: []string{"A stiff ODE that explicit methods handle well", "A purely resistive DC problem", "No difficulty, since capacitors are always state variables"},
		Why: "Capacitor-voltage-source loops (and inductor-current-source cutsets) raise the DAE index, which makes integration and initialization ill-posed for standard methods."},
	{KP: 2, Q: "Why can't SPICE simply invert C and integrate $\\dot{x} = C^{-1}(b - Gx)$ with a standard ODE solver?",
		A: "C is usually singular, and even when it is not, inverting it destroys sparsity", D: []string{"C is always diagonal, so inversion is pointless", "ODE solvers cannot handle linear systems", "Inverting C would make the method explicit and stable"},
		Why: "MNA's C has zero rows for algebraic constraints. Implicit DAE methods instead discretize the derivative and solve $(G + \\alpha C)$-type sparse systems."},
	{KP: 3, Q: "A circuit has time constants of 1 ps (a small parasitic) and 1 us (a bias filter) and is simulated for 10 us. Why would an explicit method be a poor choice?",
		A: "Stability would force steps near the 1 ps constant for the whole 10 us run", D: []string{"Explicit methods are less accurate on slow signals", "Explicit methods cannot represent capacitors", "The run would converge to the wrong DC point"},
		Why: "Stiffness: explicit methods are stable only for steps of the order of the fastest time constant, even long after that mode has died out."},
	{KP: 3, Q: "What property of backward Euler and the trapezoidal rule lets SPICE take steps much longer than the smallest circuit time constant?",
		A: "A-stability: the numerical solution of any decaying mode stays bounded for every step size", D: []string{"They are explicit", "They are fourth-order accurate", "They need no matrix factorization"},
		Why: "Both methods are A-stable, so stiff decaying modes cannot blow up. The price is solving an implicit, Newton-linearized system at each step."},
	{KP: 3, Q: "Each implicit step of a nonlinear transient requires what work?",
		A: "A Newton solve of the companion-model circuit, typically a few LU factorizations", D: []string{"One matrix-vector product", "A single forward substitution", "An eigenvalue decomposition of G"},
		Why: "Discretizing the derivatives turns each time point into a nonlinear algebraic system, solved by Newton with the same sparse machinery as the DC operating point."},
	{KP: 3, Q: "A transient run spends most of its time in device evaluation rather than LU factorization. What is the usual reason in large analog decks?",
		A: "Evaluating compact models for every device at every Newton iteration is expensive", D: []string{"LU factorization of sparse circuit matrices is free", "The integrator calls each device once per run", "Printing output dominates"},
		Why: "BSIM-class models are long sequences of floating-point code, evaluated per device per iteration. For moderate sizes this dominates; for very large post-layout decks the solver can dominate instead."},
}

var breakpointItems = []edakit.Item{
	{KP: 1, Q: "A PWL source has a corner at 5 ns, and the step controller wants to step from 4.9 ns to 5.3 ns. What does SPICE do?",
		A: "It shortens the step to land exactly on 5 ns, a breakpoint", D: []string{"It steps to 5.3 ns and interpolates", "It ignores the corner because LTE is small", "It restarts the whole simulation"},
		Why: "Source corners are registered as breakpoints. Stepping across one would smear the discontinuity in the derivative, so the step is cut to hit it."},
	{KP: 1, Q: "Why does the step controller not simply detect a source corner from a large truncation error?",
		A: "The error estimate assumes smooth waveforms, so a step across a kink gives a misleading estimate and wasted rejections", D: []string{"LTE is never computed near sources", "Corners have zero truncation error", "Breakpoints are only used for printing"},
		Why: "Divided-difference estimates of higher derivatives are garbage across a discontinuity. Knowing the corner in advance is cheaper and more accurate."},
	{KP: 1, Q: "A PULSE source has 10 ps rise and fall times and a 1 ns period over a 100 ns run. How many breakpoints does it contribute, roughly?",
		A: "About 400: four corners per period over 100 periods", D: []string{"About 100: one per period", "Two: start and end of the run", "None: PULSE sources are smooth"},
		Why: "Each period has the start and end of the rise and of the fall, four slope discontinuities, each registered as a breakpoint."},
	{KP: 1, Q: "A digital-heavy deck has many fast PULSE sources with unrelated periods. What happens to the timestep?",
		A: "It is repeatedly cut to hit breakpoints, so the run takes many small steps even where signals are idle", D: []string{"It grows, because breakpoints relax tolerances", "Nothing; breakpoints only affect output", "It becomes fixed at tstep"},
		Why: "Every source corner forces a step boundary. Many unsynchronized edges pack breakpoints densely, which is one reason fast SPICE uses event-driven methods."},
	{KP: 2, Q: "Right after a breakpoint, SPICE typically shrinks the step and may use backward Euler for the first step. Why?",
		A: "The derivative history is invalid across the kink, and an L-stable step damps the transient a trapezoidal step would ring on", D: []string{"Backward Euler is more accurate than trapezoidal", "Trapezoidal cannot start from a DC point", "To save memory"},
		Why: "Multistep history and error estimates from before the breakpoint do not describe the new behavior. A small damped first step restarts integration cleanly."},
	{KP: 2, Q: "A step is rejected three times in a row at the same time point, and the step falls below the minimum allowed. What does the user see?",
		A: "A 'timestep too small' failure", D: []string{"A singular-matrix error at DC", "A silently wrong waveform", "A finished run with fewer output points"},
		Why: "When Newton or the error test keeps failing, the controller keeps shrinking the step. Below the minimum step the simulator aborts with a timestep-too-small message."},
	{KP: 2, Q: "A Newton solve at a new time point fails to converge. What does a typical simulator do next?",
		A: "Reject the step and retry from the last accepted point with a smaller step", D: []string{"Accept the unconverged solution", "Restart from the DC operating point", "Double the step to escape the difficulty"},
		Why: "A smaller step makes the previous solution a better initial guess and the companion conductances larger, which both help Newton converge."},
	{KP: 2, Q: "A switch model changes its conductance abruptly when its control crosses a threshold. How can a simulator handle the crossing accurately?",
		A: "Locate the crossing time, for example by interpolation or bisection, and place a breakpoint there", D: []string{"Ignore it; the step controller always hits it exactly", "Use a fixed step equal to the run length", "Switch to AC analysis"},
		Why: "Internal discontinuities are not known in advance. Event location turns them into breakpoints after the fact, at the cost of a rejected step."},
	{KP: 3, Q: "A user puts uic on a .tran line and sets .ic v(out)=0. What does the simulator skip?",
		A: "The DC operating point; integration starts from the given initial conditions, with unspecified nodes at zero", D: []string{"The first timestep", "The breakpoints of all sources", "The Newton iterations at every step"},
		Why: "uic means use initial conditions: no operating-point solve, so the initial state need not be consistent and may produce a large first transient."},
	{KP: 3, Q: "What is the difference between .ic without uic and .nodeset?",
		A: ".ic holds the node at the value during the operating-point solve and releases it at t = 0; .nodeset is only an initial guess", D: []string{"They are identical", ".nodeset fixes the value for the whole transient", ".ic is ignored unless the circuit is linear"},
		Why: "Without uic, .ic constrains the operating point to the requested values. .nodeset merely helps Newton start near a solution and imposes nothing."},
	{KP: 3, Q: "An oscillator deck starts at its exact DC operating point and never oscillates in simulation. Why?",
		A: "The equilibrium is a solution, and with no noise or kick nothing perturbs it", D: []string{"The trapezoidal rule cannot simulate oscillators", "The timestep was too small", "Breakpoints damp all oscillations"},
		Why: "A deterministic simulator stays at an unstable equilibrium. Designers kick the circuit with an initial condition or a brief source pulse."},
	{KP: 3, Q: "A run starts with uic and inconsistent initial conditions, with an inductor current set but no voltage. What can happen on the first step?",
		A: "Large spurious currents or voltages as the circuit jumps to satisfy its algebraic constraints", D: []string{"Nothing, since the integrator corrects initial conditions exactly", "The simulator recomputes a DC operating point automatically", "The run becomes purely linear"},
		Why: "Inconsistent initial values violate the algebraic equations, and the first implicit step forces an abrupt jump, which can look like a real transient."},
}
