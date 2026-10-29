package csim

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const (
	kB   = 1.380649e-23
	tRef = 300.0
)

// linearizeOP asks for small-signal parameters at a bias point.
func linearizeOP(ctx generator.GeneratorContext) generator.Problem {
	id := stepF(0.1, 0.05, 38)
	vov := stepF(0.1, 0.02, 15)
	lam := stepF(0.05, 0.01, 15)
	gm := 2 * id / vov
	ro := 1 / (lam * id)
	stem := edakit.Sprintf("After the DC operating point, an NMOS in saturation (square law) carries $I_D = %s$ mA with overdrive $V_{ov} = %s$ V and channel-length modulation $\\lambda = %s$ 1/V.", f(id, 3), f(vov, 2), f(lam, 2))
	switch kp(ctx) {
	case 1:
		return num(stem+" What transconductance $g_m$ does AC analysis use, in mS?", gm, 3,
			edakit.Sprintf("For the square law $g_m = 2 I_D / V_{ov} = 2 \\times %s / %s = %s$ mS.", f(id, 3), f(vov, 2), f(gm, 3)), nil)
	case 2:
		return num(stem+" What small-signal output resistance $r_o \\approx 1/(\\lambda I_D)$ results, in k$\\Omega$?", ro, 2,
			edakit.Sprintf("$r_o = 1/(%s \\times %s) = %s$ k$\\Omega$ (1/mS is k$\\Omega$).", f(lam, 2), f(id, 3), f(ro, 2)), nil)
	}
	return num(stem+" What is the intrinsic gain $g_m r_o$?", gm*ro, 2,
		edakit.Sprintf("$g_m = %s$ mS and $r_o = %s$ k$\\Omega$, so $g_m r_o = %s$, which equals $2/(\\lambda V_{ov})$: the current cancels.", f(gm, 3), f(ro, 2), f(gm*ro, 2)), nil)
}

// complexMNA asks about the admittance matrix G + j omega C of an RC node.
func complexMNA(ctx generator.GeneratorContext) generator.Problem {
	fr := pickF(1, 2, 5, 10, 20, 50, 100, 200)
	c := pickF(0.5, 1, 2, 5, 10)
	r := pickF(1, 2, 5, 10, 20)
	w := 2 * math.Pi * fr * c // uS, with f in MHz and C in pF
	g := 1000 / r             // uS
	switch kp(ctx) {
	case 1:
		return num(edakit.Sprintf("AC analysis stamps a capacitor as the admittance $j\\omega C$. A %s pF capacitor connects node 1 to ground. At %s MHz, what is the imaginary part of its contribution to the diagonal entry $Y_{11}$, in microsiemens?", f(c, 2), f(fr, 2)), w, 2,
			edakit.Sprintf("$\\omega C = 2\\pi \\times %s \\times 10^6 \\times %s \\times 10^{-12} = %s$ uS.", f(fr, 2), f(c, 2), f(w, 2)), nil)
	case 2:
		i := pickF(1, 2, 5, 10, 20)
		mag := i / math.Hypot(g, w)
		return num(edakit.Sprintf("An AC current source of amplitude %s uA drives node 1, which has R = %s k$\\Omega$ and C = %s pF to ground. Solving the one-node complex system $(G + j\\omega C) V = I$ at %s MHz, what is $|V|$ in volts?", f(i, 2), f(r, 2), f(c, 2), f(fr, 2)), mag, 4,
			edakit.Sprintf("$G = %s$ uS and $\\omega C = %s$ uS, so $|Y| = %s$ uS and $|V| = %s / %s = %s$ V.", f(g, 2), f(w, 2), f(math.Hypot(g, w), 2), f(i, 2), f(math.Hypot(g, w), 2), f(mag, 4)), nil)
	}
	ph := -math.Atan2(w, g) * 180 / math.Pi
	return num(edakit.Sprintf("An AC current source with zero phase drives node 1, which has R = %s k$\\Omega$ and C = %s pF to ground. At %s MHz, what is the phase of the node voltage $V = I/(G + j\\omega C)$, in degrees?", f(r, 2), f(c, 2), f(fr, 2)), ph, 2,
		edakit.Sprintf("The phase is $-\\arctan(\\omega C/G) = -\\arctan(%s/%s) = %s$ degrees.", f(w, 2), f(g, 2), f(ph, 2)), nil)
}

// smallSignal asks about the gain of a common-source stage with a load capacitance.
func smallSignal(ctx generator.GeneratorContext) generator.Problem {
	gm := stepF(0.5, 0.25, 14)
	rd := pickF(5, 10, 20, 50)
	ro := pickF(20, 50, 100, 200)
	cl := pickF(0.1, 0.2, 0.5, 1, 2)
	rout := rd * ro / (rd + ro)
	a0 := gm * rout
	fp := 1 / (2 * math.Pi * rout * cl) * 1000 // MHz, rout*cl in ns
	stem := edakit.Sprintf("AC analysis linearizes a common-source stage at its operating point: $g_m = %s$ mS, drain resistor $R_D = %s$ k$\\Omega$, $r_o = %s$ k$\\Omega$, load $C_L = %s$ pF. The gain is $A(j\\omega) = -g_m R_{out}/(1 + j\\omega R_{out} C_L)$ with $R_{out} = R_D \\parallel r_o$.", f(gm, 2), f(rd, 2), f(ro, 2), f(cl, 2))
	switch kp(ctx) {
	case 1:
		db := 20 * math.Log10(a0)
		return num(stem+" What is the low-frequency gain magnitude in dB?", db, 2,
			edakit.Sprintf("$R_{out} = %s$ k$\\Omega$, $|A_0| = %s$, $20\\log_{10} %s = %s$ dB.", f(rout, 3), f(a0, 3), f(a0, 3), f(db, 2)), nil)
	case 2:
		fr := fp * pickF(0.5, 1, 2, 3, 5, 10)
		mag := a0 / math.Sqrt(1+math.Pow(fr/fp, 2))
		return num(stem+edakit.Sprintf(" What is $|A|$ at %s MHz?", f(fr, 2)), mag, 3,
			edakit.Sprintf("$|A_0| = %s$ and the pole is at %s MHz, so $|A| = %s/\\sqrt{1 + (%s/%s)^2} = %s$.", f(a0, 3), f(fp, 3), f(a0, 3), f(fr, 2), f(fp, 3), f(mag, 3)), nil)
	}
	return num(stem+" At what frequency is the gain 3 dB below its low-frequency value, in MHz?", fp, 2,
		edakit.Sprintf("The pole is at $1/(2\\pi R_{out} C_L)$ with $R_{out} C_L = %s$ ns, which is %s MHz.", f(rout*cl, 4), f(fp, 2)), nil)
}

// noise asks about thermal noise densities and kT/C noise.
func noise(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		r := pickF(0.05, 0.1, 0.5, 1, 2, 5, 10, 50, 100)
		e := math.Sqrt(4*kB*tRef*r*1e3) * 1e9
		return num(edakit.Sprintf("Noise analysis models a resistor's thermal noise as a voltage source with density $4kTR$ in V$^2$/Hz. For R = %s k$\\Omega$ at T = 300 K ($k = 1.380649 \\times 10^{-23}$ J/K), what is the noise voltage density in nV/$\\sqrt{\\text{Hz}}$?", f(r, 2)), e, 3,
			edakit.Sprintf("$4kTR = %s$ V$^2$/Hz; its square root is %s nV/$\\sqrt{\\text{Hz}}$.", sciM(4*kB*tRef*r*1e3), f(e, 3)), nil)
	case 2:
		e1, e2 := stepF(2, 0.5, 16), stepF(2, 0.5, 16)
		g1, g2 := stepF(1, 0.5, 18), stepF(0.5, 0.5, 8)
		out := math.Hypot(g1*e1, g2*e2)
		return num(edakit.Sprintf("Two uncorrelated noise sources reach the output: a resistor with %s nV/$\\sqrt{\\text{Hz}}$ through a gain magnitude of %s, and a transistor's input-referred noise of %s nV/$\\sqrt{\\text{Hz}}$ through a gain of %s. Noise analysis adds the output power densities. What is the total output noise density in nV/$\\sqrt{\\text{Hz}}$?", f(e1, 2), f(g1, 2), f(e2, 2), f(g2, 2)), out, 3,
			edakit.Sprintf("Contributions are %s and %s nV/$\\sqrt{\\text{Hz}}$; powers add, so $\\sqrt{%s^2 + %s^2} = %s$.", f(g1*e1, 3), f(g2*e2, 3), f(g1*e1, 3), f(g2*e2, 3), f(out, 3)), nil)
	}
	c := pickF(0.05, 0.1, 0.2, 0.5, 1, 2, 5)
	v := math.Sqrt(kB*tRef/(c*1e-12)) * 1e6
	return num(edakit.Sprintf("Integrating the thermal noise of any resistor across a capacitor over all frequencies gives a total variance $kT/C$, independent of R. For C = %s pF at 300 K, what is the rms noise voltage in microvolts?", f(c, 2)), v, 2,
		edakit.Sprintf("$kT/C = %s$ V$^2$, and the square root is %s uV.", sciM(kB*tRef/(c*1e-12)), f(v, 2)), nil)
}
