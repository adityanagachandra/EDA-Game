package csim

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// sensitivity differentiates a resistive divider with respect to its resistors.
func sensitivity(ctx generator.GeneratorContext) generator.Problem {
	vs := stepF(1, 0.5, 6)
	r1 := pickF(1, 2, 3, 4, 5, 10)
	r2 := pickF(1, 2, 3, 4, 5, 10)
	vout := vs * r2 / (r1 + r2)
	which := edakit.Pick([]string{"R1", "R2"})
	wm := "R_1"
	if which == "R2" {
		wm = "R_2"
	}
	var d, s, rv float64
	if which == "R1" {
		d = -vs * r2 / ((r1 + r2) * (r1 + r2))
		s = -r1 / (r1 + r2)
		rv = r1
	} else {
		d = vs * r1 / ((r1 + r2) * (r1 + r2))
		s = r1 / (r1 + r2)
		rv = r2
	}
	stem := edakit.Sprintf("A divider has V1 = %s V at the top, R1 = %s k$\\Omega$ from the source to node out, and R2 = %s k$\\Omega$ from out to ground, so $v_{out} = V_1 R_2/(R_1 + R_2) = %s$ V.", f(vs, 2), f(r1, 2), f(r2, 2), f(vout, 4))
	expl := edakit.Sprintf("Differentiating, $\\partial v_{out}/\\partial %s = %s$ V/k$\\Omega$; the direct method gets the same number by solving $G\\, \\partial x/\\partial p = -(\\partial G/\\partial p)\\, x$ with the already-factored $G$.", wm, f(d, 5))
	switch kp(ctx) {
	case 1:
		return num(stem+edakit.Sprintf(" What is $\\partial v_{out}/\\partial %s$ in mV per k$\\Omega$?", wm), d*1e3, 2, expl, nil)
	case 2:
		return num(stem+edakit.Sprintf(" What is the normalized sensitivity $(%s/v_{out})\\, \\partial v_{out}/\\partial %s$?", wm, wm), s, 4,
			expl+edakit.Sprintf(" Normalizing by $%s/v_{out}$ gives %s: a 1%% change in %s moves $v_{out}$ by %s%%.", wm, f(s, 4), which, f(s, 4)), nil)
	}
	pct := pickF(1, 2, 5, 10, -1, -2, -5)
	dv := d * rv * pct / 100 * 1e3
	return num(stem+edakit.Sprintf(" If %s changes by %s%%, what first-order change in $v_{out}$ does the sensitivity predict, in mV?", which, f(pct, 0)), dv, 3,
		expl+edakit.Sprintf(" $\\Delta R = %s$ k$\\Omega$, so $\\Delta v \\approx %s \\times %s = %s$ mV.", f(rv*pct/100, 4), f(d, 5), f(rv*pct/100, 4), f(dv, 3)), nil)
}

// adjoint covers the adjoint solve, adjoint sensitivities and solve counts.
func adjoint(ctx generator.GeneratorContext) generator.Problem {
	r1, r12, r2 := pickF(1, 2, 4, 5), pickF(1, 2, 4, 5), pickF(1, 2, 4, 5)
	g1, g12, g2 := 1/r1, 1/r12, 1/r2
	i := stepF(0.5, 0.5, 7)
	det := (g1+g12)*(g2+g12) - g12*g12
	v1 := i * (g2 + g12) / det
	v2 := i * g12 / det
	l1 := g12 / det
	l2 := (g1 + g12) / det
	stem := edakit.Sprintf("A current source pushes %s mA into node 1. $R_a$ = %s k$\\Omega$ from node 1 to ground, $R_b$ = %s k$\\Omega$ between nodes 1 and 2, $R_c$ = %s k$\\Omega$ from node 2 to ground. The output is $v_2$. The nominal solve gives $v_1 = %s$ V and $v_2 = %s$ V.", f(i, 2), f(r1, 2), f(r12, 2), f(r2, 2), f(v1, 4), f(v2, 4))
	switch kp(ctx) {
	case 1:
		which := edakit.Int(1, 2)
		l := l1
		if which == 2 {
			l = l2
		}
		return num(stem+edakit.Sprintf(" The adjoint system is $G^T \\lambda = e_2$ (a unit current injected at the output node). What is $\\lambda_%d$ in k$\\Omega$ (V per mA)?", which), l, 4,
			edakit.Sprintf("$G$ is symmetric for resistors, so $\\lambda = G^{-1} e_2$: $\\det G = %s$, $\\lambda_1 = g_b/\\det = %s$, $\\lambda_2 = (g_a + g_b)/\\det = %s$ k$\\Omega$.", f(det, 5), f(l1, 4), f(l2, 4)), nil)
	case 2:
		which := edakit.Pick([]string{"a", "b", "c"})
		var s float64
		switch which {
		case "a":
			s = -l1 * v1
		case "b":
			s = -(l1 - l2) * (v1 - v2)
		default:
			s = -l2 * v2
		}
		return num(stem+edakit.Sprintf(" The adjoint solve gives $\\lambda_1 = %s$ and $\\lambda_2 = %s$ k$\\Omega$. For a conductance $g$ between nodes p and q, $\\partial v_2/\\partial g = -(\\lambda_p - \\lambda_q)(v_p - v_q)$, with ground at 0. What is $\\partial v_2/\\partial g_%s$ (the conductance of $R_%s$), in V/mS?", f(l1, 4), f(l2, 4), which, which), s, 4,
			edakit.Sprintf("Plug the two solutions into the formula for $g_%s$: %s V/mS. Every element's sensitivity comes from the same pair of solutions.", which, f(s, 4)), nil)
	}
	p := edakit.Int(20, 400)
	m := edakit.Int(1, 4)
	if edakit.Int(0, 1) == 0 {
		return edakit.Numeric(edakit.Sprintf("A circuit has %d parameters and %d outputs of interest. After the nominal solve and its LU factorization, how many additional forward/back substitutions does the direct method need to get every output's sensitivity to every parameter?", p, m),
			edakit.IntStr(p),
			edakit.Sprintf("The direct method solves $G\\, \\partial x/\\partial p = -(\\partial G/\\partial p) x$ once per parameter: %d solves, each giving all outputs.", p), nil)
	}
	return edakit.Numeric(edakit.Sprintf("A circuit has %d parameters and %d outputs of interest. After the nominal solve and its LU factorization, how many additional substitutions (with the transposed factors) does the adjoint method need to get every output's sensitivity to every parameter?", p, m),
		edakit.IntStr(m),
		edakit.Sprintf("The adjoint method solves $G^T \\lambda = e_{out}$ once per output: %d solves, after which each parameter costs only an inner product.", m), nil)
}

// mismatch applies the Pelgrom area law.
func mismatch(ctx generator.GeneratorContext) generator.Problem {
	a := pickF(1, 1.5, 2, 2.5, 3, 4, 5)
	w := pickF(0.2, 0.5, 1, 2, 4, 5, 10)
	l := pickF(0.05, 0.1, 0.2, 0.5, 1, 2)
	sig := a / math.Sqrt(w*l)
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("A process has a Pelgrom coefficient $A_{VT} = %s$ mV um for threshold mismatch, $\\sigma(\\Delta V_T) = A_{VT}/\\sqrt{WL}$. For a pair of identical transistors with W = %s um and L = %s um, what is $\\sigma(\\Delta V_T)$ in mV?", f(a, 2), f(w, 2), f(l, 2)), sig, 3,
			edakit.Sprintf("$WL = %s$ um$^2$, $\\sqrt{WL} = %s$ um, so $\\sigma = %s$ mV.", f(w*l, 4), f(math.Sqrt(w*l), 4), f(sig, 3)), nil)
	case 2:
		target := pickF(0.5, 1, 2, 3, 5)
		lneed := (a / target) * (a / target) / w
		return num(edakit.Sprintf("With $A_{VT} = %s$ mV um, a designer needs $\\sigma(\\Delta V_T) \\le %s$ mV for a matched pair at W = %s um. What minimum gate length L does that require, in um?", f(a, 2), f(target, 2), f(w, 2)), lneed, 4,
			edakit.Sprintf("$WL \\ge (A_{VT}/\\sigma)^2 = %s$ um$^2$, so $L \\ge %s / %s = %s$ um. Halving the mismatch costs four times the area.", f((a/target)*(a/target), 4), f((a/target)*(a/target), 4), f(w, 2), f(lneed, 4)), nil)
	}
	dev := sig / math.Sqrt2
	return num(edakit.Sprintf("Measured pair mismatch follows $\\sigma(\\Delta V_T) = A_{VT}/\\sqrt{WL}$ with $A_{VT} = %s$ mV um. A Monte Carlo model gives each transistor an independent threshold shift. For W = %s um and L = %s um, what per-device standard deviation must the model use so that the simulated pair difference matches, in mV?", f(a, 2), f(w, 2), f(l, 2)), dev, 3,
		edakit.Sprintf("The difference of two independent draws has variance $2\\sigma_{dev}^2$, so $\\sigma_{dev} = \\sigma(\\Delta V_T)/\\sqrt2 = %s/\\sqrt2 = %s$ mV. Using the pair value per device would overstate mismatch by $\\sqrt2$.", f(sig, 3), f(dev, 3)), nil)
}

// monteCarloCorners covers corner counts, Monte Carlo sample sizes and standard errors.
func monteCarloCorners(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		p := edakit.Pick([]int{3, 5})
		v := edakit.Int(2, 3)
		t := edakit.Int(2, 4)
		rc := edakit.Pick([]int{1, 3, 5})
		n := p * v * t * rc
		names := "TT, FF, SS"
		if p == 5 {
			names = "TT, FF, SS, FS, SF"
		}
		return edakit.Numeric(edakit.Sprintf("A sign-off plan simulates every combination of %d process corners (%s), %d supply voltages, %d temperatures and %d parasitic-extraction corners. How many simulator runs is that per testbench?", p, names, v, t, rc),
			edakit.IntStr(n),
			edakit.Sprintf("Corners multiply: $%d \\times %d \\times %d \\times %d = %d$. This product growth is why corner lists are pruned and why statistical methods replace some of them.", p, v, t, rc, n), nil)
	case 2:
		y := pickF(0.5, 0.8, 0.9, 0.95, 0.99)
		se := pickF(0.005, 0.01, 0.02)
		n := int(math.Ceil(y*(1-y)/(se*se) - 1e-9))
		return edakit.Numeric(edakit.Sprintf("A Monte Carlo yield run expects a yield near %s. Each sample is pass or fail, so the estimate's standard error is $\\sqrt{Y(1 - Y)/N}$. How many samples are needed for a standard error of %s?", f(y, 2), f(se, 3)),
			edakit.IntStr(n),
			edakit.Sprintf("$N = Y(1 - Y)/\\text{se}^2 = %s/%s = %d$ (rounded up). The cost grows as $1/\\text{se}^2$ and is independent of the number of parameters.", f(y*(1-y), 4), f(se*se, 6), n), nil)
	}
	sig := stepF(1, 0.5, 18)
	n := edakit.Pick([]int{25, 50, 100, 200, 400, 1000})
	se := sig / math.Sqrt(float64(n))
	return num(edakit.Sprintf("Monte Carlo of an amplifier's offset gives a sample standard deviation of %s mV over %d runs. What is the standard error of the estimated mean offset, in mV?", f(sig, 2), n), se, 4,
		edakit.Sprintf("$\\sigma/\\sqrt{N} = %s/\\sqrt{%d} = %s$ mV. Quadrupling N halves it.", f(sig, 2), n, f(se, 4)), nil)
}

var processCornerItems = []edakit.Item{
	{KP: 1, Q: "In the corner name FS, what do the two letters describe?",
		A: "Fast NMOS and slow PMOS", D: []string{"Fast PMOS and slow NMOS", "Fast supply and slow temperature", "Fast setup and slow hold"},
		Why: "Corner names list the NMOS speed first, then the PMOS: FS is fast NMOS, slow PMOS, a skewed corner."},
	{KP: 1, Q: "What does the TT corner represent?",
		A: "Typical NMOS and typical PMOS process parameters", D: []string{"The worst-case timing corner", "Typical temperature only", "A Monte Carlo average of all corners"},
		Why: "TT uses nominal (typical) global process parameters for both device types."},
	{KP: 1, Q: "A corner model card for SS shifts threshold voltages up and mobility down for both device types. What does that model?",
		A: "Global, die-to-die variation toward slow devices", D: []string{"Local mismatch between neighboring devices", "Aging of the devices", "Temperature at 125 C"},
		Why: "Process corners are shifted global parameter sets. They move every device on the die together, which is global variation, not mismatch."},
	{KP: 1, Q: "Why do skewed corners FS and SF exist in addition to FF and SS?",
		A: "Some circuits fail when NMOS and PMOS strengths are unbalanced, such as ratioed logic and inverter switching thresholds", D: []string{"They are always the fastest corners", "They model temperature inversion", "They replace Monte Carlo entirely"},
		Why: "An inverter's switching threshold and ratioed circuits depend on the NMOS-to-PMOS strength ratio, which the skewed corners stress."},
	{KP: 2, Q: "Which corner combination is usually checked for the worst setup timing of a digital path?",
		A: "SS process, low supply voltage, and the temperature that makes the devices slowest", D: []string{"FF process, high supply, low temperature", "TT process at nominal supply", "FS process at high supply"},
		Why: "Setup fails when the path is slow, so slow devices at low voltage are the standard pessimistic corner. The slowest temperature depends on the technology."},
	{KP: 2, Q: "Which corner is typically worst for hold timing and for dynamic power?",
		A: "FF process with high supply voltage", D: []string{"SS process with low supply voltage", "TT process at room temperature", "SF process at low voltage"},
		Why: "Hold fails when data paths are fast, and switching power rises with supply voltage and fast devices. Both point to the FF, high-voltage corner."},
	{KP: 2, Q: "Leakage power of a standard-cell block must be bounded. Which corner and temperature matter most?",
		A: "FF devices at high temperature", D: []string{"SS devices at low temperature", "TT devices at low temperature", "Any corner, since leakage is corner-independent"},
		Why: "Leakage grows exponentially as thresholds drop and as temperature rises, so fast devices at the hot corner give the worst case."},
	{KP: 2, Q: "In an advanced low-voltage process, a path is slower at -40 C than at 125 C. What effect is this?",
		A: "Temperature inversion: at low supply the threshold increase at cold dominates the mobility gain", D: []string{"Electromigration", "A simulator convergence artifact", "Negative bias temperature instability"},
		Why: "Cooling raises mobility but also threshold voltage. Near threshold, the threshold effect wins, so cold becomes the slow corner."},
	{KP: 3, Q: "A comparator's input offset is checked at all five process corners and passes. Why is that not evidence that offset is fine?",
		A: "Offset comes from local mismatch, which corners do not model; it needs Monte Carlo with mismatch", D: []string{"Corners are too pessimistic for offset", "Offset depends only on temperature", "Five corners is too few; seven would suffice"},
		Why: "Global corners shift both input transistors equally, so their difference stays zero. Mismatch is local and random and requires statistical simulation."},
	{KP: 3, Q: "Why is a ±3 sigma process corner not necessarily the ±3 sigma value of a circuit's performance?",
		A: "The performance depends on many correlated parameters nonlinearly, so the parameter corner maps to a different probability in performance space", D: []string{"Because sigma is undefined for circuits", "It always is the same value by definition", "Because corners include mismatch twice"},
		Why: "A corner pushes every parameter to an extreme together. For a given metric that combination may be too pessimistic or not the worst at all."},
	{KP: 3, Q: "A team replaces FF and SS for one analog block with corners derived from the block's own Monte Carlo distribution. What are these called?",
		A: "Statistical or performance-specific corners", D: []string{"Skewed corners", "Extraction corners", "Aging corners"},
		Why: "Statistical corners are chosen so that the block's metric reaches a target percentile, giving a sharper bound than generic device corners."},
	{KP: 3, Q: "A digital library is characterized at SS and FF. An analog designer reuses those two corners as the only checks for a bandgap reference. What is missing?",
		A: "Temperature sweep, mismatch Monte Carlo and the skewed corners that the bandgap's curvature and offset depend on", D: []string{"Nothing: SS and FF bound every circuit", "Only the TT corner", "Only a higher supply voltage"},
		Why: "Corner sets are chosen per metric. A reference voltage depends on temperature behavior and on mismatch, neither of which two global corners capture."},
}
