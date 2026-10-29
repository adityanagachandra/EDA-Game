package cts

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const boltzmannEV = 8.617e-5 // eV per kelvin

// dynamicIR: a current step into a node held up by decap C behind grid resistance R.
func dynamicIR(ctx generator.GeneratorContext) generator.Problem {
	r := pick(0.05, 0.5, 0.05)
	c := float64(100 * edakit.Int(2, 20))
	i := float64(10 * edakit.Int(5, 50))
	tau := r * c
	desc := edakit.Sprintf("A region of the grid is modelled as a resistance of $%s\\ \\Omega$ from the supply to a node with %s pF of local decoupling capacitance. At $t = 0$ the node's load current steps from 0 to %s mA.", f(r, 2), f(c, 0), f(i, 0))
	switch kpOf(ctx) {
	case 1:
		t := float64(10 * edakit.Int(1, int(3*tau/10)+1))
		d := i * r * (1 - math.Exp(-t/tau))
		return edakit.Numeric(desc+edakit.Sprintf(" What is the voltage drop at the node %s ps after the step, in mV? Round to 2 decimal places.", f(t, 0)), f(d, 2),
			edakit.Sprintf("Charge first comes out of the decap, so the drop rises as $IR(1 - e^{-t/RC})$ with $RC = %s \\cdot %s = %s$ ps (ohm times pF gives ps). At %s ps: $%s(1 - e^{-%s}) = %s$ mV, against a static %s mV.", f(r, 2), f(c, 0), f(tau, 1), f(t, 0), f(i*r, 2), f(t/tau, 3), f(d, 2), f(i*r, 2)),
			map[string]string{"tau": f(tau, 2), "drop": f(d, 3)})
	case 2:
		frac := pick(0.3, 0.9, 0.05)
		x := frac * i * r
		t := -tau * math.Log(1-frac)
		return edakit.Numeric(desc+edakit.Sprintf(" How long after the step does the drop reach %s mV, in ps? Round to 1 decimal place.", f(x, 3)), f(t, 1),
			edakit.Sprintf("The final drop is $IR = %s$ mV and $RC = %s$ ps. Solve $%s = %s(1 - e^{-t/%s})$: $t = -RC \\ln(1 - %s) = %s$ ps.", f(i*r, 2), f(tau, 1), f(x, 3), f(i*r, 2), f(tau, 1), f(frac, 2), f(t, 1)),
			map[string]string{"t": f(t, 2)})
	}
	tb := float64(10 * edakit.Int(2, int(2*tau/10)+2))
	tp := tb * float64(edakit.Int(4, 20))
	peak := i * r * (1 - math.Exp(-tb/tau))
	static := i * r * tb / tp
	ratio := peak / static
	return edakit.Numeric(desc+edakit.Sprintf(" The current is in fact a burst: %s mA for %s ps, then zero, repeating every %s ps. A static analysis uses the average current. What is the ratio of the peak drop during the first burst to the static drop? Round to 2 decimal places.", f(i, 0), f(tb, 0), f(tp, 0)), f(ratio, 2),
		edakit.Sprintf("Static: average current $%s \\cdot %s/%s$ mA gives %s mV. Dynamic: the burst ends at $t_b = %s$ ps with drop $%s(1 - e^{-%s/%s}) = %s$ mV. Ratio %s. Average-current analysis can miss the real droop by this factor.", f(i, 0), f(tb, 0), f(tp, 0), f(static, 3), f(tb, 0), f(i*r, 2), f(tb, 0), f(tau, 1), f(peak, 3), f(ratio, 2)),
		map[string]string{"peak": f(peak, 3), "static": f(static, 3)})
}

// decap: Q = C dV sizing.
func decap(ctx generator.GeneratorContext) generator.Problem {
	dv := float64(5 * edakit.Int(4, 20))
	switch kpOf(ctx) {
	case 1:
		i := float64(50 * edakit.Int(2, 40))
		dt := pick(0.1, 2, 0.1)
		c := i * dt / dv
		return edakit.Numeric(edakit.Sprintf("A block's current rises by %s mA for %s ns before the package and grid can respond, so local decoupling must supply the charge. The droop budget is %s mV. How much decoupling capacitance is needed, in nF? Round to 2 decimal places.", f(i, 0), f(dt, 1), f(dv, 0)), f(c, 2),
			edakit.Sprintf("$Q = I\\,\\Delta t = %s \\cdot %s = %s$ pC must come from the decap while its voltage falls by at most %s mV, so $C = Q/\\Delta V = %s$ nF (pC over mV is nF).", f(i, 0), f(dt, 1), f(i*dt, 1), f(dv, 0), f(c, 2)),
			map[string]string{"q_pc": f(i*dt, 2)})
	case 2:
		ip := float64(50 * edakit.Int(2, 40))
		tw := pick(0.1, 2, 0.1)
		q := ip * tw / 2
		c := q / dv
		return edakit.Numeric(edakit.Sprintf("A switching event draws a triangular current pulse that peaks at %s mA and lasts %s ns. All of its charge must come from local decap with at most %s mV of droop. How much decap is needed, in nF? Round to 2 decimal places.", f(ip, 0), f(tw, 1), f(dv, 0)), f(c, 2),
			edakit.Sprintf("The charge is the area of the triangle, $\\tfrac{1}{2} \\cdot %s \\cdot %s = %s$ pC. $C = Q/\\Delta V = %s/%s = %s$ nF.", f(ip, 0), f(tw, 1), f(q, 2), f(q, 2), f(dv, 0), f(c, 2)),
			map[string]string{"q_pc": f(q, 2)})
	}
	need := pick(2, 20, 0.5)
	intr := need * pick(0.2, 0.7, 0.05)
	dens := float64(5 * edakit.Int(2, 8))
	area := (need - intr) * 1e6 / dens
	return edakit.Numeric(edakit.Sprintf("A block needs %s nF of decoupling. The non-switching cells and the grid already provide %s nF of intrinsic decap. Explicit decap cells give %s fF/µm². How much decap-cell area must be added, in µm²? Round to the nearest whole number.", f(need, 1), f(intr, 2), f(dens, 0)), f(area, 0),
		edakit.Sprintf("Explicit decap covers the shortfall $%s - %s = %s$ nF $= %s$ fF, which at %s fF/µm² takes %s µm².", f(need, 1), f(intr, 2), f(need-intr, 3), f((need-intr)*1e6, 0), f(dens, 0), f(area, 0)),
		map[string]string{"area": f(area, 0)})
}

// ldiDt: inductive noise, bump count, and package-die resonance.
func ldiDt(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		l := float64(10 * edakit.Int(2, 30))
		di := float64(50 * edakit.Int(2, 30))
		dt := float64(50 * edakit.Int(1, 10))
		v := l * di / dt
		return edakit.Numeric(edakit.Sprintf("The supply path through the package has an effective inductance of %s pH. The chip's current rises by %s mA in %s ps. What is the inductive voltage noise $L\\,di/dt$, in mV? Round to 1 decimal place.", f(l, 0), f(di, 0), f(dt, 0)), f(v, 1),
			edakit.Sprintf("$L\\,\\Delta I/\\Delta t = %s \\cdot %s/%s = %s$ mV (pH times mA over ps gives mV).", f(l, 0), f(di, 0), f(dt, 0), f(v, 1)),
			map[string]string{"v": f(v, 2)})
	case 2:
		lb := float64(10 * edakit.Int(5, 30))
		di := float64(100 * edakit.Int(5, 50))
		dt := float64(100 * edakit.Int(1, 10))
		vb := float64(5 * edakit.Int(2, 10))
		n := int(math.Ceil(lb*di/(dt*vb) - 1e-9))
		return edakit.Numeric(edakit.Sprintf("Each power bump and its package path has %s pH of inductance. The chip's current rises by %s mA in %s ps, and the inductive noise budget is %s mV. Treating bumps as independent parallel inductors, how many power bumps are needed?", f(lb, 0), f(di, 0), f(dt, 0), f(vb, 0)), edakit.IntStr(n),
			edakit.Sprintf("$N$ bumps in parallel give $L_b/N$. Require $(L_b/N)\\,\\Delta I/\\Delta t \\le %s$ mV: $N \\ge %s \\cdot %s/(%s \\cdot %s) = %s$, so %d bumps (the same count is needed again for ground). Mutual inductance makes real bumps less effective than this.", f(vb, 0), f(lb, 0), f(di, 0), f(dt, 0), f(vb, 0), f(lb*di/(dt*vb), 3), n),
			map[string]string{"n": edakit.IntStr(n)})
	}
	l := float64(10 * edakit.Int(2, 50))
	c := float64(10 * edakit.Int(2, 50))
	fr := 1 / (2 * math.Pi * math.Sqrt(l*1e-12*c*1e-9)) / 1e6
	return edakit.Numeric(edakit.Sprintf("The package presents %s pH of supply inductance to the die, which has %s nF of on-die decoupling. At what frequency does this LC pair resonate, in MHz? Round to 1 decimal place.", f(l, 0), f(c, 0)), f(fr, 1),
		edakit.Sprintf("$f = 1/(2\\pi\\sqrt{LC}) = 1/(2\\pi\\sqrt{%s \\times 10^{-12} \\cdot %s \\times 10^{-9}}) = %s$ MHz. Near this frequency the PDN impedance peaks, and current activity there produces the largest droop.", f(l, 0), f(c, 0), f(fr, 1)),
		map[string]string{"f_mhz": f(fr, 2)})
}

// electromigration: current density, minimum width, via count.
func electromigration(ctx generator.GeneratorContext) generator.Problem {
	t := pick(0.2, 1.0, 0.05)
	jmax := pick(2, 20, 0.5)
	switch kpOf(ctx) {
	case 1:
		i := pick(1, 50, 0.5)
		w := pick(0.5, 10, 0.5)
		j := i / (w * t)
		return edakit.Numeric(edakit.Sprintf("A power strap %s µm wide and %s µm thick carries a DC current of %s mA. What is its current density, in mA/µm²? Round to 2 decimal places.", f(w, 1), f(t, 2), f(i, 1)), f(j, 2),
			edakit.Sprintf("$J = I/(wt) = %s/(%s \\cdot %s) = %s$ mA/µm².", f(i, 1), f(w, 1), f(t, 2), f(j, 2)),
			map[string]string{"j": f(j, 3)})
	case 2:
		i := pick(1, 50, 0.5)
		w := i / (jmax * t)
		return edakit.Numeric(edakit.Sprintf("A power strap of thickness %s µm must carry %s mA of DC current. The electromigration limit for this layer at the signoff temperature is %s mA/µm². What is the minimum strap width, in µm? Round to 2 decimal places.", f(t, 2), f(i, 1), f(jmax, 1)), f(w, 2),
			edakit.Sprintf("$w \\ge I/(J_{max} t) = %s/(%s \\cdot %s) = %s$ µm.", f(i, 1), f(jmax, 1), f(t, 2), f(w, 2)),
			map[string]string{"w": f(w, 3)})
	}
	i := pick(5, 80, 0.5)
	iv := pick(0.2, 2, 0.1)
	n := int(math.Ceil(i/iv - 1e-9))
	return edakit.Numeric(edakit.Sprintf("A strap-to-strap connection must carry %s mA of DC current. Each via of that cut layer is rated for %s mA at the signoff temperature. What is the smallest number of vias in the via array?", f(i, 1), f(iv, 1)), edakit.IntStr(n),
		edakit.Sprintf("$\\lceil %s/%s \\rceil = %d$ vias. Vias are a frequent EM weak point because current crowds through a small cross-section at the via-wire interface.", f(i, 1), f(iv, 1), n),
		map[string]string{"vias": edakit.IntStr(n)})
}

// blackEquation: MTTF = A J^-n exp(Ea / kT).
func blackEquation(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Pick([]float64{1, 2})
		k := pick(1.1, 2.5, 0.1)
		ratio := math.Pow(1/k, n)
		return edakit.Numeric(edakit.Sprintf("Black's equation models electromigration lifetime as $\\mathrm{MTTF} = A J^{-n} e^{E_a/kT}$ with $n = %s$. An ECO raises the current density in a strap by a factor of %s at the same temperature. By what factor does the median time to failure change? Round to 3 decimal places.", f(n, 0), f(k, 1)), f(ratio, 3),
			edakit.Sprintf("Only $J^{-n}$ changes: the ratio is $%s^{-%s} = %s$.", f(k, 1), f(n, 0), f(ratio, 3)),
			map[string]string{"ratio": f(ratio, 4)})
	case 2:
		ea := pick(0.7, 1.0, 0.05)
		t1c := float64(5 * edakit.Int(11, 21))
		t2c := t1c + float64(5*edakit.Int(1, 6))
		t1, t2 := t1c+273.15, t2c+273.15
		ratio := math.Exp(ea / boltzmannEV * (1/t2 - 1/t1))
		return edakit.Numeric(edakit.Sprintf("Black's equation: $\\mathrm{MTTF} = A J^{-n} e^{E_a/kT}$ with $E_a = %s$ eV and $k = 8.617 \\times 10^{-5}$ eV/K. A strap's temperature rises from %s to %s degrees C at the same current density. By what factor does the median time to failure change? Round to 3 decimal places.", f(ea, 2), f(t1c, 0), f(t2c, 0)), f(ratio, 3),
			edakit.Sprintf("Convert to kelvin: %s K and %s K. The ratio is $\\exp\\big(\\tfrac{E_a}{k}(\\tfrac{1}{T_2} - \\tfrac{1}{T_1})\\big) = \\exp(%s) = %s$. Self-heating and hot spots shorten EM life sharply.", f(t1, 2), f(t2, 2), f(ea/boltzmannEV*(1/t2-1/t1), 4), f(ratio, 3)),
			map[string]string{"ratio": f(ratio, 4)})
	}
	ea := pick(0.7, 1.0, 0.05)
	n := 2.0
	t1c := float64(5 * edakit.Int(11, 21))
	t2c := t1c + float64(5*edakit.Int(1, 6))
	t1, t2 := t1c+273.15, t2c+273.15
	frac := 100 * math.Exp(ea/(n*boltzmannEV)*(1/t2-1/t1))
	return edakit.Numeric(edakit.Sprintf("Black's equation: $\\mathrm{MTTF} = A J^{-n} e^{E_a/kT}$ with $n = 2$, $E_a = %s$ eV and $k = 8.617 \\times 10^{-5}$ eV/K. The EM limit was set for %s degrees C. To keep the same lifetime at %s degrees C, what percentage of the original current-density limit may the strap carry? Round to 1 decimal place.", f(ea, 2), f(t1c, 0), f(t2c, 0)), f(frac, 1),
		edakit.Sprintf("Equal lifetime needs $J_2^{-2} e^{E_a/kT_2} = J_1^{-2} e^{E_a/kT_1}$, so $J_2/J_1 = \\exp\\big(\\tfrac{E_a}{2k}(\\tfrac{1}{T_2} - \\tfrac{1}{T_1})\\big) = %s$, which is %s%% of the original limit.", f(frac/100, 4), f(frac, 1)),
		map[string]string{"frac": f(frac, 2)})
}
