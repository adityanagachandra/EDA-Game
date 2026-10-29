package csim

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const vt = 0.025 // thermal voltage stated in every diode stem, volts

// isChoices are diode saturation currents in amperes with their printed form.
var isChoices = []struct {
	v float64
	s string
}{{1e-14, "10 fA"}, {2e-14, "20 fA"}, {5e-15, "5 fA"}, {1e-15, "1 fA"}, {1e-13, "100 fA"}}

// opReduction draws a circuit with capacitors and inductors and asks about its DC solution.
func opReduction(ctx generator.GeneratorContext) generator.Problem {
	vs := stepF(1, 0.5, 8)
	r1, r2, r3, r4 := edakit.Pick(kohmChoices), edakit.Pick(kohmChoices), edakit.Pick(kohmChoices), edakit.Pick(kohmChoices)
	withR3 := edakit.Int(0, 1) == 0
	desc := edakit.Sprintf("V1 drives node in at %s V. R1 = %s k$\\Omega$ runs from in to node a, capacitor C1 from a to ground, inductor L1 from a to node b, R2 = %s k$\\Omega$ from b to ground, and capacitor C2 from b to node c with R4 = %s k$\\Omega$ from c to ground", f(vs, 2), f(r1, 2), f(r2, 2), f(r4, 2))
	rbot := r2
	if withR3 {
		desc += edakit.Sprintf(", plus R3 = %s k$\\Omega$ from a to ground", f(r3, 2))
		rbot = r2 * r3 / (r2 + r3)
	}
	desc += "."
	va := vs * rbot / (r1 + rbot)
	switch kp(ctx) {
	case 1:
		node := edakit.Pick([]string{"a", "b", "a", "b", "c"})
		v := va
		if node == "c" {
			v = 0
		}
		return num(desc+edakit.Sprintf(" What is the DC operating-point voltage at node %s, in volts?", node), v, 3,
			edakit.Sprintf("At DC the capacitors are open and the inductor is a short, so a and b are one node and c sees no current through C2. The divider R1 into %s k$\\Omega$ gives $v_a = v_b = %s$ V, and $v_c = 0$.", f(rbot, 3), f(va, 4)),
			map[string]string{"va": f(va, 4)})
	case 2:
		il := va / r2
		return num(desc+" What DC current flows through L1 from a to b, in mA?", il, 3,
			edakit.Sprintf("With L1 shorted, $v_b = v_a = %s$ V, and all the inductor current flows on through R2 (C2 is open): $%s / %s = %s$ mA.", f(va, 4), f(va, 4), f(r2, 2), f(il, 4)), nil)
	}
	nodes := 4 // in, a, b, c
	u := nodes + 1 + 1
	return edakit.Numeric(desc+" SPICE keeps the inductor as a branch-current unknown with zero voltage across it at DC. How many unknowns does the DC MNA system have?",
		edakit.IntStr(u),
		edakit.Sprintf("Four node voltages (in, a, b, c), the V1 branch current and the L1 branch current: %d. The capacitors stamp nothing at DC; they stay in the netlist so the same matrix structure serves the transient.", u), nil)
}

// diodeAt returns the diode current and conductance at v.
func diodeAt(is, v float64) (id, gd float64) {
	e := math.Exp(v / vt)
	return is * (e - 1), is / vt * e
}

// diodeCompanion asks for the companion conductance, equivalent current or prediction.
func diodeCompanion(ctx generator.GeneratorContext) generator.Problem {
	is := isChoices[edakit.Int(0, len(isChoices)-1)]
	vk := stepF(0.55, 0.01, 15)
	id, gd := diodeAt(is.v, vk)
	for id > 0.05 || id < 2e-6 {
		vk = stepF(0.55, 0.01, 15)
		id, gd = diodeAt(is.v, vk)
	}
	idm, gdm := id*1e3, gd*1e3 // mA, mS
	ieq := idm - gdm*vk
	stem := edakit.Sprintf("A diode with $I_S$ = %s and $V_T = 0.025$ V obeys $I = I_S(e^{V/V_T} - 1)$. Newton's current iterate is $V_k = %s$ V.", is.s, f(vk, 2))
	facts := map[string]string{"id_ma": f(idm, 5), "gd_ms": f(gdm, 5), "ieq_ma": f(ieq, 5)}
	expl := edakit.Sprintf("$I_d(V_k) = %s$ mA, $g_d = (I_S/V_T)e^{V_k/V_T} = %s$ mS, $I_{eq} = I_d - g_d V_k = %s$ mA.", f(idm, 5), f(gdm, 5), f(ieq, 5))
	switch kp(ctx) {
	case 1:
		return num(stem+" What companion conductance $g_d$ does the simulator stamp, in mS?", gdm, placesFor(gdm), expl, facts)
	case 2:
		return num(stem+" The companion model is $g_d$ in parallel with a current source $I_{eq}$ flowing from anode to cathode, so that $I = g_d V + I_{eq}$. What is $I_{eq}$ in mA?", ieq, placesFor(ieq), expl, facts)
	}
	dv := stepF(0.02, 0.01, 3)
	pred := gdm*(vk+dv) + ieq
	trueI, _ := diodeAt(is.v, vk+dv)
	return num(stem+edakit.Sprintf(" The linearized companion model $I = g_d V + I_{eq}$ built at $V_k$ is used at $V = %s$ V. What current does it predict, in mA?", f(vk+dv, 2)), pred, placesFor(pred),
		expl+edakit.Sprintf(" The tangent line gives $I_d + g_d \\Delta V = %s$ mA, while the true diode current is %s mA: the exponential outruns its tangent, which is why Newton needs several iterations.", f(pred, 5), f(trueI*1e3, 5)), facts)
}

// newtonCompanion runs one Newton step on source - resistor - diode.
func newtonCompanion(ctx generator.GeneratorContext) generator.Problem {
	is := isChoices[edakit.Int(0, len(isChoices)-1)]
	vs := stepF(1, 0.5, 8)
	r := edakit.Pick([]float64{0.5, 1, 2, 4, 5, 10})
	vk := stepF(0.55, 0.01, 15)
	id, gd := diodeAt(is.v, vk)
	for id > 0.05 || id < 2e-6 {
		vk = stepF(0.55, 0.01, 15)
		id, gd = diodeAt(is.v, vk)
	}
	idm, gdm := id*1e3, gd*1e3
	ieq := idm - gdm*vk
	gtot := 1/r + gdm
	vnext := (vs/r - ieq) / gtot
	res := (vk-vs)/r + idm
	stem := edakit.Sprintf("A source of %s V drives a %s k$\\Omega$ resistor into node d; a diode from d to ground has $I_S$ = %s, $V_T = 0.025$ V and $I = I_S(e^{V/V_T}-1)$. The current Newton iterate is $v_d = %s$ V. Work in mS, mA and V.", f(vs, 2), f(r, 2), is.s, f(vk, 2))
	expl := edakit.Sprintf("At $V_k$: $I_d = %s$ mA, $g_d = %s$ mS, $I_{eq} = %s$ mA. The linearized KCL row is $(1/R + g_d)\\, v = V_s/R - I_{eq}$, so $%s\\, v = %s$ and $v_{k+1} = %s$ V.", f(idm, 5), f(gdm, 5), f(ieq, 5), f(gtot, 5), f(vs/r-ieq, 5), f(vnext, 5))
	facts := map[string]string{"gd_ms": f(gdm, 5), "ieq_ma": f(ieq, 5), "v_next": f(vnext, 5)}
	switch kp(ctx) {
	case 1:
		return num(stem+" What is the single entry of the linearized MNA matrix for node d, in mS?", gtot, placesFor(gtot), expl, facts)
	case 2:
		return num(stem+" What is the next iterate $v_{k+1}$ in volts?", vnext, 4, expl, facts)
	}
	return num(stem+" What is the KCL residual at node d, the current leaving through the resistor and the diode, $(v_k - V_s)/R + I_d(v_k)$, in mA?", res, placesFor(res),
		edakit.Sprintf("$(%s - %s)/%s + %s = %s$ mA. Newton drives this residual to zero; its size, together with the voltage update, decides convergence.", f(vk, 2), f(vs, 2), f(r, 2), f(idm, 5), f(res, 5)), facts)
}

// newtonConvergence covers SPICE's tolerance tests and quadratic convergence.
func newtonConvergence(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := edakit.Round(stepF(0.3, 0.0137, 300), 4)
		b := edakit.Round(a+stepF(-0.002, 0.0001, 40), 4)
		tol := 1e-3*math.Max(math.Abs(a), math.Abs(b))*1e6 + 1
		return num(edakit.Sprintf("Two successive Newton iterates of a node voltage are %s V and %s V. With SPICE defaults RELTOL = $10^{-3}$ and VNTOL = 1 uV, the voltage test passes when $|v_k - v_{k-1}| \\le \\text{RELTOL} \\cdot \\max(|v_k|, |v_{k-1}|) + \\text{VNTOL}$. What is the allowed band, in microvolts?", f(b, 4), f(a, 4)),
			tol, 1,
			edakit.Sprintf("$10^{-3} \\times %s$ V $= %s$ uV, plus 1 uV, gives %s uV. The actual change is %s uV.", f(math.Max(math.Abs(a), math.Abs(b)), 4), f(tol-1, 1), f(tol, 1), f(math.Abs(a-b)*1e6, 1)), nil)
	case 2:
		c := pickF(0.5, 1, 2, 4, 5)
		e0 := pickF(0.05, 0.1, 0.02, 0.08, 0.15)
		for c*e0 >= 0.8 {
			e0 /= 2
		}
		tol := pickF(1e-6, 1e-9, 1e-12)
		e := e0
		k := 0
		var seq []string
		for e >= tol && k < 20 {
			e = c * e * e
			k++
			seq = append(seq, edakit.Sprintf("%.3g", e))
		}
		return edakit.Numeric(edakit.Sprintf("Near a solution, a Newton iteration's error obeys $e_{k+1} = %s\\, e_k^2$ (in volts). The initial error is %s V. How many iterations are needed before the error first falls below %s V?", f(c, 2), f(e0, 4), sci(tol)),
			edakit.IntStr(k),
			edakit.Sprintf("Squaring each time: %s. The error first drops below the tolerance after %d iterations; quadratic convergence roughly doubles the correct digits per step.", strings.Join(seq, ", "), k), nil)
	}
	a := edakit.Round(stepF(0.05, 0.037, 100), 5)
	b := edakit.Round(a*(1+stepF(-0.01, 0.001, 20)), 5)
	tol := 1e-3*math.Max(a, b)*1e6 + 1e-3 // nA, with a and b in mA
	return num(edakit.Sprintf("Two successive iterates of a device branch current are %s mA and %s mA. With RELTOL = $10^{-3}$ and ABSTOL = 1 pA, the current test allows $\\text{RELTOL} \\cdot \\max(|i_k|, |i_{k-1}|) + \\text{ABSTOL}$. What is the allowed band, in nanoamperes?", f(b, 5), f(a, 5)),
		tol, 3,
		edakit.Sprintf("$10^{-3} \\times %s$ mA $= %s$ nA, plus 0.001 nA, gives %s nA.", f(math.Max(a, b), 5), f(tol-1e-3, 4), f(tol, 3)), nil)
}

// mosfetLinearization asks for gm, gds or the equivalent current of a square-law MOSFET.
func mosfetLinearization(ctx generator.GeneratorContext) generator.Problem {
	k := stepF(0.2, 0.1, 18)    // mA/V^2
	vth := stepF(0.35, 0.05, 5) // V
	vgs := vth + stepF(0.15, 0.05, 10)
	vds := (vgs - vth) + stepF(0.2, 0.1, 10)
	lam := stepF(0.04, 0.02, 7)
	vov := vgs - vth
	id := k / 2 * vov * vov * (1 + lam*vds)
	gm := k * vov * (1 + lam*vds)
	gds := k / 2 * vov * vov * lam
	ieq := id - gm*vgs - gds*vds
	stem := edakit.Sprintf("An NMOS in saturation follows $I_D = \\tfrac{k}{2}(V_{GS} - V_T)^2 (1 + \\lambda V_{DS})$ with $k = %s$ mA/V$^2$, $V_T = %s$ V and $\\lambda = %s$ 1/V. Newton's iterate has $V_{GS} = %s$ V and $V_{DS} = %s$ V (source grounded).", f(k, 2), f(vth, 2), f(lam, 2), f(vgs, 2), f(vds, 2))
	expl := edakit.Sprintf("$V_{ov} = %s$ V. $I_D = %s$ mA, $g_m = k V_{ov}(1 + \\lambda V_{DS}) = %s$ mS, $g_{ds} = \\tfrac{k}{2} V_{ov}^2 \\lambda = %s$ mS, $I_{eq} = I_D - g_m V_{GS} - g_{ds} V_{DS} = %s$ mA.", f(vov, 2), f(id, 5), f(gm, 5), f(gds, 5), f(ieq, 5))
	facts := map[string]string{"id_ma": f(id, 5), "gm_ms": f(gm, 5), "gds_ms": f(gds, 5), "ieq_ma": f(ieq, 5)}
	switch kp(ctx) {
	case 1:
		return num(stem+" What transconductance $g_m = \\partial I_D/\\partial V_{GS}$ is stamped, in mS?", gm, 4, expl, facts)
	case 2:
		return num(stem+" What output conductance $g_{ds} = \\partial I_D/\\partial V_{DS}$ is stamped, in microsiemens?", gds*1e3, 2, expl, facts)
	}
	return num(stem+" The companion model is $I_D \\approx g_m V_{GS} + g_{ds} V_{DS} + I_{eq}$. What is $I_{eq}$ in mA?", ieq, 4, expl, facts)
}

// pnjlim is SPICE's junction voltage limiter.
func pnjlim(vnew, vold, vcrit float64) float64 {
	if vnew > vcrit && math.Abs(vnew-vold) > 2*vt {
		if vold > 0 {
			arg := 1 + (vnew-vold)/vt
			if arg > 0 {
				return vold + vt*math.Log(arg)
			}
			return vcrit
		}
		return vt * math.Log(vnew/vt)
	}
	return vnew
}

// junctionLimiting asks for the critical voltage or the limited junction voltage.
func junctionLimiting(ctx generator.GeneratorContext) generator.Problem {
	is := isChoices[edakit.Int(0, len(isChoices)-1)]
	vcrit := vt * math.Log(vt/(math.Sqrt2*is.v))
	rule := edakit.Sprintf("SPICE limits a junction update with $v_{crit} = V_T \\ln(V_T / (\\sqrt{2} I_S))$: if $v_{new} > v_{crit}$ and $|v_{new} - v_{old}| > 2V_T$, then for $v_{old} > 0$ it uses $v_{old} + V_T \\ln(1 + (v_{new} - v_{old})/V_T)$, and for $v_{old} \\le 0$ it uses $V_T \\ln(v_{new}/V_T)$; otherwise $v_{new}$ is kept. Take $V_T = 0.025$ V and $I_S$ = %s.", is.s)
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("A diode has $I_S$ = %s and $V_T = 0.025$ V. SPICE's junction limiter switches on above the critical voltage $v_{crit} = V_T \\ln(V_T/(\\sqrt{2} I_S))$, the voltage where the diode's current-voltage curve has its minimum radius of curvature. What is $v_{crit}$ in volts?", is.s), vcrit, 4,
			edakit.Sprintf("$V_T/(\\sqrt{2} I_S) = %s$, its natural log is %s, times 0.025 gives %s V.", sciM(vt/(math.Sqrt2*is.v)), f(math.Log(vt/(math.Sqrt2*is.v)), 4), f(vcrit, 4)), map[string]string{"vcrit": f(vcrit, 4)})
	case 2:
		vold := stepF(0.5, 0.01, 25)
		vnew := stepF(1.2, 0.1, 40)
		lim := pnjlim(vnew, vold, vcrit)
		return num(rule+edakit.Sprintf(" The previous junction voltage is %s V and the raw Newton update proposes %s V. What junction voltage is used for the next device evaluation?", f(vold, 2), f(vnew, 2)), lim, 4,
			edakit.Sprintf("$v_{crit} = %s$ V, the proposal exceeds it and the step is far larger than $2V_T$, so $%s + 0.025 \\ln(1 + %s/0.025) = %s$ V. The exponential is then evaluated at a sane voltage.", f(vcrit, 4), f(vold, 2), f(vnew-vold, 2), f(lim, 4)), nil)
	}
	var vold, vnew float64
	switch edakit.Int(0, 3) {
	case 0: // below vcrit: unchanged
		vold = stepF(0.3, 0.01, 20)
		vnew = math.Round((vcrit-stepF(0.05, 0.01, 10))*100) / 100
	case 1: // small step above vcrit: unchanged
		vnew = math.Round((vcrit+stepF(0.01, 0.01, 5))*100) / 100
		vold = vnew - stepF(0.01, 0.005, 6)
	case 2: // reverse-biased old point
		vold = -stepF(0.1, 0.1, 9)
		vnew = stepF(1, 0.1, 30)
	default:
		vold = stepF(0.5, 0.01, 25)
		vnew = stepF(1.2, 0.1, 40)
	}
	lim := pnjlim(vnew, vold, vcrit)
	return num(rule+edakit.Sprintf(" The previous junction voltage is %s V and Newton proposes %s V. What voltage is used?", f(vold, 3), f(vnew, 3)), lim, 4,
		edakit.Sprintf("$v_{crit} = %s$ V and the step is %s V. The rule gives %s V.", f(vcrit, 4), f(vnew-vold, 3), f(lim, 4)), map[string]string{"vcrit": f(vcrit, 4)})
}

// damping covers step clipping, uniform vector damping and backtracking.
func damping(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		v := stepF(0.2, 0.05, 16)
		d := stepF(-2, 0.1, 40)
		if math.Abs(d) < 0.05 {
			d = 0.85
		}
		dmax := pickF(0.1, 0.2, 0.3, 0.5)
		step := math.Copysign(math.Min(math.Abs(d), dmax), d)
		return num(edakit.Sprintf("A Newton solve at iterate $v = %s$ V proposes the update $\\Delta v = %s$ V. The simulator clips any node update to at most %s V in magnitude. What is the next iterate, in volts?", f(v, 2), f(d, 2), f(dmax, 2)), v+step, 3,
			edakit.Sprintf("The applied step is %s V, so the iterate becomes %s V. Clipping keeps the step in the region where the linearization is trustworthy.", f(step, 3), f(v+step, 3)), nil)
	case 2:
		x := []float64{stepF(0.1, 0.1, 15), stepF(0.1, 0.1, 15), stepF(0.1, 0.1, 15)}
		d := []float64{stepF(-1.5, 0.1, 30), stepF(-1.5, 0.1, 30), stepF(-1.5, 0.1, 30)}
		dmax := pickF(0.2, 0.25, 0.4, 0.5)
		inf := 0.0
		for _, v := range d {
			inf = math.Max(inf, math.Abs(v))
		}
		if inf == 0 {
			d[0], inf = 1, 1
		}
		alpha := math.Min(1, dmax/inf)
		j := edakit.Int(0, 2)
		return num(edakit.Sprintf("Node voltages $(%s, %s, %s)$ V receive the Newton update $(%s, %s, %s)$ V. The solver scales the whole update by $\\alpha = \\min(1, %s / \\|\\Delta\\|_\\infty)$ so that no component moves more than %s V. What is the new value of component %d, in volts?", f(x[0], 2), f(x[1], 2), f(x[2], 2), f(d[0], 2), f(d[1], 2), f(d[2], 2), f(dmax, 2), f(dmax, 2), j+1),
			x[j]+alpha*d[j], 4,
			edakit.Sprintf("$\\|\\Delta\\|_\\infty = %s$, $\\alpha = %s$, so component %d becomes $%s + %s \\times %s = %s$ V. Scaling all components together keeps the Newton direction.", f(inf, 2), f(alpha, 4), j+1, f(x[j], 2), f(alpha, 4), f(d[j], 2), f(x[j]+alpha*d[j], 4)), nil)
	}
	f0 := stepF(1, 0.5, 18)
	accept := edakit.Int(0, 3)
	alphas := []float64{1, 0.5, 0.25, 0.125}
	var norms []float64
	for i, a := range alphas {
		bound := (1 - a/2) * f0
		var v float64
		if i < accept {
			v = bound * (1 + stepF(0.05, 0.05, 20))
		} else {
			v = bound * stepF(0.3, 0.05, 10)
		}
		norms = append(norms, math.Round(v*100)/100)
	}
	ans := 0.0
	for i, a := range alphas {
		if norms[i] <= (1-a/2)*f0 {
			ans = a
			break
		}
	}
	var parts []string
	for i, a := range alphas {
		parts = append(parts, edakit.Sprintf("$\\alpha = %s$: %s", f(a, 3), f(norms[i], 2)))
	}
	return num(edakit.Sprintf("Damped Newton tries $x + \\alpha \\Delta$ for $\\alpha = 1, 1/2, 1/4, 1/8$ and accepts the first $\\alpha$ with $\\|f(x + \\alpha \\Delta)\\| \\le (1 - \\alpha/2) \\|f(x)\\|$. The current residual norm is %s mA, and the trial residual norms are %s (mA). Which $\\alpha$ is accepted? Answer as a decimal.", f(f0, 2), strings.Join(parts, "; ")),
		ans, 3,
		edakit.Sprintf("Compare each trial with $(1 - \\alpha/2) \\times %s$: the first that passes is $\\alpha = %s$.", f(f0, 2), f(ans, 3)), nil)
}

// gminStepping covers the shunt conductance, the schedule and the residual error.
func gminStepping(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		i := stepF(0.05, 0.05, 9)
		r := edakit.Pick([]float64{1, 2, 5, 10, 20, 50})
		gmin := pickF(0.01, 0.05, 0.1, 0.2, 0.5)
		v := i / (1/r + gmin)
		return num(edakit.Sprintf("A current source pushes %s mA into a node that has a %s k$\\Omega$ resistor to ground. During gmin stepping the simulator also attaches $g_{min} = %s$ mS from every node to ground. What node voltage does that intermediate solve produce, in volts?", f(i, 2), f(r, 2), f(gmin, 2)), v, 4,
			edakit.Sprintf("The node sees $1/R + g_{min} = %s$ mS, so $v = %s / %s = %s$ V, against %s V once gmin is removed. Each stepping solve starts from the previous one as gmin shrinks.", f(1/r+gmin, 4), f(i, 2), f(1/r+gmin, 4), f(v, 4), f(i*r, 4)), nil)
	case 2:
		start := edakit.Int(-4, -2)
		final := edakit.Pick([]int{-12, -12, -14, -10})
		fac := edakit.Pick([]int{1, 1, 2})
		for (start-final)%fac != 0 {
			start--
		}
		n := (start-final)/fac + 1
		return edakit.Numeric(edakit.Sprintf("A gmin-stepping schedule starts with $g_{min} = 10^{%d}$ S at every node, divides it by $10^{%d}$ after each converged solve, and finishes with a solve at $10^{%d}$ S. How many operating-point solves does the schedule run, counting the first and the last?", start, fac, final),
			edakit.IntStr(n),
			edakit.Sprintf("The exponent falls from %d to %d in steps of %d: %d steps plus the first solve gives %d.", start, final, fac, (start-final)/fac, n), nil)
	}
	vs := stepF(0.5, 0.5, 6)
	r := edakit.Pick([]float64{1e8, 2e8, 5e8, 1e9, 2e9, 1e10})
	gmin := 1e-12
	v := vs / (1 + r*gmin)
	errUV := (vs - v) * 1e6
	return num(edakit.Sprintf("A %s V source drives a high-impedance node through %s ohms; the node has no other path to ground except the simulator's final $g_{min} = 10^{-12}$ S. By how much is the computed node voltage below %s V, in microvolts?", f(vs, 2), sci(r), f(vs, 2)), errUV, 1,
		edakit.Sprintf("The node is a divider of $R$ and $1/g_{min}$: $v = V_s/(1 + R g_{min}) = %s$ V, an error of %s uV. Gmin is negligible on ordinary nodes but not on gigaohm ones.", f(v, 7), f(errUV, 1)), nil)
}

var convergenceAidItems = []edakit.Item{
	{KP: 1, Q: "A deck contains a node connected only to capacitor plates and a MOSFET gate. The DC operating point reports a singular matrix. Why?",
		A: "At DC the node has no conductive path to ground, so its row in G is all zeros", D: []string{"The MOSFET exponential overflowed", "The timestep was too large", "Newton started too close to the solution"},
		Why: "Capacitors are open at DC and an ideal gate draws no current, so nothing stamps that node's row. A gmin shunt or a high-value resistor gives it a DC path."},
	{KP: 1, Q: "Newton starts a diode circuit at 0 V, and the first update proposes 8 V across the junction. What goes wrong if the simulator evaluates the diode there?",
		A: "The exponential overflows or gives a huge current and conductance, and the next iterate is meaningless", D: []string{"The matrix becomes symmetric", "The diode current becomes negative", "Nothing; Newton always converges from any start"},
		Why: "At 8 V, $e^{V/V_T}$ is astronomically large. Junction limiting replaces such proposals with a logarithmically compressed voltage."},
	{KP: 1, Q: "A latch has two stable states and a metastable one. Newton converges to the metastable point. Is the operating point wrong?",
		A: "It is a valid DC solution; the equations have several, and Newton found the one nearest its path", D: []string{"It is always a convergence failure", "It means gmin is too large", "DC analysis cannot have more than one solution"},
		Why: "Nonlinear circuits with positive feedback have multiple DC solutions. Which one Newton finds depends on the initial guess, which is why .nodeset exists."},
	{KP: 1, Q: "An operating point oscillates between two iterates and never meets the tolerance. Which property of Newton explains this?",
		A: "Newton converges only locally; far from a solution the full step can overshoot back and forth", D: []string{"Newton is globally convergent, so the matrix must be wrong", "The tolerance is applied only to currents", "LU factorization introduces the oscillation"},
		Why: "Quadratic convergence holds only near a root. Far away, full Newton steps can cycle, which damping and limiting are designed to stop."},
	{KP: 2, Q: "A large mixed circuit fails the operating point with a plain Newton solve from zero. Which aid adds a shunt conductance at every node and removes it gradually?",
		A: "Gmin stepping", D: []string{"Source stepping", "Junction limiting", "Timestep control"},
		Why: "Gmin stepping starts with a large conductance to ground at every node, which makes the problem nearly linear and well conditioned, then shrinks it in a continuation."},
	{KP: 2, Q: "Which aid ramps every independent source from zero to its full value, solving at each step from the previous solution?",
		A: "Source stepping", D: []string{"Gmin stepping", "Pseudo-transient analysis", "Damping"},
		Why: "At zero source value the all-zero solution is exact for most circuits. Raising the sources in small increments keeps each Newton solve close to its starting point."},
	{KP: 2, Q: "A designer knows a bandgap should start near 1.2 V at its output but the solver lands in the zero-current state. What is the cleanest fix?",
		A: "A .nodeset on the output to start Newton near the intended solution", D: []string{"Raising reltol", "Switching to Gear integration", "Deleting the startup circuit"},
		Why: "A nodeset is an initial guess for the operating point that is released after the first iterations, steering Newton toward the desired one of several solutions."},
	{KP: 2, Q: "Both gmin and source stepping fail on a circuit with strong positive feedback. Which method treats the operating point as the end of a fictitious transient?",
		A: "Pseudo-transient continuation", D: []string{"Junction limiting", "Increasing itl1 only", "AC analysis"},
		Why: "Pseudo-transient adds artificial capacitance (or inductance) and integrates in time until the circuit settles. Time-stepping follows a physical path that rarely gets stuck."},
	{KP: 3, Q: "Gmin stepping converged, but the final solve with gmin at its default value did not. What does the earlier converged point mean?",
		A: "Nothing for the real circuit: it solved a modified circuit with extra conductances", D: []string{"It is the answer to report", "It proves the circuit has no solution", "It means reltol was too tight"},
		Why: "Intermediate continuation solutions belong to a different circuit. Only a converged solve with the final, negligible gmin is the operating point."},
	{KP: 3, Q: "Source stepping on a Schmitt trigger stalls at 60% of the supply. What is the likely reason?",
		A: "The solution curve folds back at a turning point, so increasing the source parameter cannot follow it", D: []string{"The matrix has too many nonzeros", "The supply is too low for Newton", "Source stepping cannot handle MOSFETs"},
		Why: "Hysteretic circuits have solution curves with turning points. Naive parameter stepping jumps or fails there; arc-length continuation or pseudo-transient can follow the fold."},
	{KP: 3, Q: "An engineer sets gmin to 1e-6 S globally to make a deck converge. What is the risk?",
		A: "Every node now leaks a microsiemens to ground, corrupting high-impedance nodes and leakage results", D: []string{"Transient steps become smaller", "The matrix becomes singular", "No risk: gmin is removed after convergence"},
		Why: "The gmin option sets a conductance that stays in the circuit. A microsiemens across a nanoampere-biased node changes its voltage substantially."},
	{KP: 3, Q: "Damping by clipping every node update to 0.1 V fixes divergence but a 50-node solve now takes 200 iterations. What has been traded?",
		A: "Robustness for speed: small fixed steps lose Newton's quadratic convergence far from the root", D: []string{"Accuracy for robustness: the answer is now less accurate", "Nothing, since iteration count does not matter", "Memory for speed"},
		Why: "Clipped steps creep toward the solution linearly. The converged answer is as accurate as before, but reaching it costs many more factorizations."},
}
