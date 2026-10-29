package analog

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func commonSource(ctx generator.GeneratorContext) generator.Problem {
	gm := unif(0.5, 5, 0.1)
	rd := float64(edakit.Int(5, 50))
	ro := float64(edakit.Int(20, 200))
	switch whichKP(ctx) {
	case 1:
		a := gm * par(rd, ro)
		q := edakit.Sprintf("A common-source stage has $g_m = %s$ mS, a resistive load $R_D = %s$ kilohms and device output resistance $r_o = %s$ kilohms. What is the magnitude of its small-signal voltage gain in V/V? Round to 2 decimal places.", f(gm, 1), f(rd, 0), f(ro, 0))
		return num(q, f(a, 2), edakit.Sprintf("$|A_v| = g_m (R_D \\parallel r_o) = %s \\cdot %s = %s$. The gain is inverting.", f(gm, 1), f(par(rd, ro), 3), f(a, 2)), map[string]string{"rout": f(par(rd, ro), 3)})
	case 2:
		rs := unif(0.2, 3, 0.1)
		a := gm * rd / (1 + gm*rs)
		q := edakit.Sprintf("A common-source stage with source degeneration has $g_m = %s$ mS, $R_D = %s$ kilohms and $R_S = %s$ kilohms. Ignoring $r_o$, what is the magnitude of its voltage gain in V/V? Round to 2 decimal places.", f(gm, 1), f(rd, 0), f(rs, 1))
		return num(q, f(a, 2), edakit.Sprintf("Degeneration divides the transconductance: $G_m = g_m/(1 + g_m R_S) = %s$ mS, so $|A_v| = G_m R_D = %s$. Gain falls, linearity improves.", f(gm/(1+gm*rs), 4), f(a, 2)), map[string]string{"Gm": f(gm/(1+gm*rs), 4)})
	default:
		cl := unif(0.1, 5, 0.1)
		r := par(rd, ro)
		fp := 1000 / (2 * math.Pi * r * cl)
		q := edakit.Sprintf("A common-source stage has $R_D = %s$ kilohms, $r_o = %s$ kilohms and drives a load capacitance $C_L = %s$ pF. Where is its output pole, in MHz? Round to 2 decimal places.", f(rd, 0), f(ro, 0), f(cl, 1))
		return num(q, f(fp, 2), edakit.Sprintf("The output node sees $R_{out} = R_D \\parallel r_o = %s$ kilohms, and $f_p = 1/(2\\pi R_{out} C_L)$. Kilohms times picofarads is nanoseconds, so $f_p = 1000/(2\\pi \\cdot %s \\cdot %s) = %s$ MHz.", f(r, 3), f(r, 3), f(cl, 1), f(fp, 2)), map[string]string{"rout": f(r, 3)})
	}
}

func currentMirror(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		iref := float64(5 * edakit.Int(2, 40))
		w1 := float64(edakit.Int(1, 8))
		w2 := float64(edakit.Int(1, 24))
		l := float64(edakit.Int(1, 4))
		io := iref * w2 / w1
		q := edakit.Sprintf("A simple NMOS current mirror has a diode-connected reference device with $W/L = %s/%s$ carrying $I_{REF} = %s\\ \\mu$A. The output device has $W/L = %s/%s$ (same $L$, both saturated). Ignoring channel-length modulation, what is $I_{OUT}$ in microamperes? Round to 2 decimal places.", f(w1, 0), f(l, 0), f(iref, 0), f(w2, 0), f(l, 0))
		return num(q, f(io, 2), edakit.Sprintf("Equal $V_{GS}$ forces equal current density, so $I_{OUT} = I_{REF}(W_2/W_1) = %s \\cdot %s/%s = %s\\ \\mu$A.", f(iref, 0), f(w2, 0), f(w1, 0), f(io, 2)), nil)
	case 2:
		iref := float64(10 * edakit.Int(2, 20))
		lam := unif(0.05, 0.3, 0.01)
		vgs := unif(0.4, 0.8, 0.05)
		vds := unif(0.3, 1.5, 0.05)
		io := iref * (1 + lam*vds) / (1 + lam*vgs)
		q := edakit.Sprintf("A 1:1 current mirror copies $I_{REF} = %s\\ \\mu$A. The reference diode sits at $V_{DS1} = V_{GS} = %s$ V, the output device at $V_{DS2} = %s$ V, and $\\lambda = %s\\ \\mathrm{V}^{-1}$. What is $I_{OUT}$ in microamperes? Round to 2 decimal places.", f(iref, 0), f(vgs, 2), f(vds, 2), f(lam, 2))
		return num(q, f(io, 2), edakit.Sprintf("$I_{OUT} = I_{REF}(1+\\lambda V_{DS2})/(1+\\lambda V_{DS1}) = %s \\cdot %s/%s = %s\\ \\mu$A. The copy error comes from the $V_{DS}$ mismatch.", f(iref, 0), f(1+lam*vds, 4), f(1+lam*vgs, 4), f(io, 2)), nil)
	default:
		gm := unif(0.2, 3, 0.1)
		ro1 := float64(edakit.Int(20, 150))
		ro2 := float64(edakit.Int(20, 150))
		r := gm * ro2 * ro1 / 1000
		q := edakit.Sprintf("A cascode current source stacks a cascode device ($g_{m2} = %s$ mS, $r_{o2} = %s$ kilohms) on a source device with $r_{o1} = %s$ kilohms. Using $R_{out} \\approx g_{m2} r_{o2} r_{o1}$, what is the output resistance in megohms? Round to 2 decimal places.", f(gm, 1), f(ro2, 0), f(ro1, 0))
		return num(q, f(r, 2), edakit.Sprintf("The cascode multiplies $r_{o1}$ by its intrinsic gain $g_{m2}r_{o2} = %s$: $R_{out} \\approx %s \\cdot %s$ kilohms $= %s$ megohms.", f(gm*ro2, 1), f(gm*ro2, 1), f(ro1, 0), f(r, 2)), nil)
	}
}

func diffPair(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		iss := float64(10 * edakit.Int(2, 50))
		gmid := float64(edakit.Int(6, 22))
		gm := gmid * iss / 2
		q := edakit.Sprintf("A differential pair is biased by a tail current $I_{SS} = %s\\ \\mu$A, and each input device runs at $g_m/I_D = %s\\ \\mathrm{V}^{-1}$. What is the transconductance of one input device, in microsiemens? Round to the nearest integer.", f(iss, 0), f(gmid, 0))
		return num(q, f(gm, 0), edakit.Sprintf("At balance each side carries $I_{SS}/2 = %s\\ \\mu$A, so $g_{m1} = %s \\cdot %s = %s\\ \\mu$S.", f(iss/2, 0), f(gmid, 0), f(iss/2, 0), f(gm, 0)), nil)
	case 2:
		gm := unif(0.1, 2, 0.05)
		ro2 := float64(edakit.Int(50, 400))
		ro4 := float64(edakit.Int(50, 400))
		a := gm * par(ro2, ro4)
		q := edakit.Sprintf("A five-transistor OTA has input-pair transconductance $g_{m1} = %s$ mS, input-device output resistance $r_{o2} = %s$ kilohms and mirror-load output resistance $r_{o4} = %s$ kilohms. What is its differential-to-single-ended gain in V/V? Round to 1 decimal place.", f(gm, 2), f(ro2, 0), f(ro4, 0))
		return num(q, f(a, 1), edakit.Sprintf("The mirror load folds both half-currents onto the output, so $A_d = g_{m1}(r_{o2} \\parallel r_{o4}) = %s \\cdot %s = %s$ V/V.", f(gm, 2), f(par(ro2, ro4), 2), f(a, 1)), nil)
	default:
		vov := unif(0.08, 0.35, 0.01)
		v := math.Sqrt2 * vov * 1000
		q := edakit.Sprintf("A square-law differential pair has input devices at overdrive $V_{OV} = %s$ V when balanced. Beyond what differential input voltage does one side carry the entire tail current, in millivolts? Round to 1 decimal place.", f(vov, 2))
		return num(q, f(v, 1), edakit.Sprintf("Full steering happens at $|V_{id}| = \\sqrt{2} V_{OV} = %s$ mV. A larger $V_{OV}$ (lower $g_m/I_D$) buys linear range at the cost of $g_m$ per unit current.", f(v, 1)), nil)
	}
}

func twoStage(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		gm1 := unif(0.1, 1, 0.05)
		ro2 := float64(edakit.Int(100, 400))
		ro4 := float64(edakit.Int(100, 400))
		gm6 := unif(0.5, 4, 0.1)
		ro6 := float64(edakit.Int(20, 100))
		ro7 := float64(edakit.Int(20, 100))
		a1 := gm1 * par(ro2, ro4)
		a2 := gm6 * par(ro6, ro7)
		tot := db20(a1 * a2)
		q := edakit.Sprintf("A two-stage op-amp has a first stage with $g_{m1} = %s$ mS and $r_{o2} \\parallel r_{o4}$ from %s and %s kilohms, and a common-source second stage with $g_{m6} = %s$ mS and $r_{o6} \\parallel r_{o7}$ from %s and %s kilohms. What is the DC open-loop gain in dB? Round to 1 decimal place.", f(gm1, 2), f(ro2, 0), f(ro4, 0), f(gm6, 1), f(ro6, 0), f(ro7, 0))
		return num(q, f(tot, 1), edakit.Sprintf("$A_1 = %s \\cdot %s = %s$ and $A_2 = %s \\cdot %s = %s$. $A_0 = A_1 A_2 = %s$, which is %s dB.", f(gm1, 2), f(par(ro2, ro4), 2), f(a1, 2), f(gm6, 1), f(par(ro6, ro7), 2), f(a2, 2), f(a1*a2, 0), f(tot, 1)), map[string]string{"A1": f(a1, 2), "A2": f(a2, 2)})
	case 2:
		gm1 := float64(10 * edakit.Int(5, 100))
		cc := unif(0.5, 5, 0.1)
		gbw := gm1 / (2 * math.Pi * cc)
		q := edakit.Sprintf("A Miller-compensated two-stage op-amp has input-pair transconductance $g_{m1} = %s\\ \\mu$S and compensation capacitor $C_c = %s$ pF. What is its gain-bandwidth product in MHz? Round to 2 decimal places.", f(gm1, 0), f(cc, 1))
		return num(q, f(gbw, 2), edakit.Sprintf("$\\mathrm{GBW} = g_{m1}/(2\\pi C_c)$. Microsiemens over picofarads is $10^6$ per second, so $\\mathrm{GBW} = %s/(2\\pi \\cdot %s) = %s$ MHz.", f(gm1, 0), f(cc, 1), f(gbw, 2)), nil)
	default:
		a0 := unif(60, 100, 1)
		gbw := unif(1, 200, 1)
		fp := gbw * 1e6 / undb20(a0) / 1000
		q := edakit.Sprintf("A two-stage op-amp has DC gain %s dB and a single-pole roll-off with gain-bandwidth product %s MHz. Where is its dominant pole, in kHz? Round to 3 decimal places.", f(a0, 0), f(gbw, 0))
		return num(q, f(fp, 3), edakit.Sprintf("$A_0 = 10^{%s/20} = %s$, and $f_{p1} = \\mathrm{GBW}/A_0 = %s$ kHz.", f(a0, 0), f(undb20(a0), 0), f(fp, 3)), nil)
	}
}

func miller(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		cc := unif(0.5, 5, 0.1)
		a2 := float64(edakit.Int(10, 120))
		c := cc * (1 + a2)
		q := edakit.Sprintf("A compensation capacitor $C_c = %s$ pF bridges a second stage whose gain is -%s V/V. What capacitance does the Miller effect present at the second stage's input, in pF? Round to 1 decimal place.", f(cc, 1), f(a2, 0))
		return num(q, f(c, 1), edakit.Sprintf("$C_{in} = C_c(1 + |A_2|) = %s \\cdot %s = %s$ pF. This large capacitance on the first-stage output creates the dominant pole.", f(cc, 1), f(1+a2, 0), f(c, 1)), nil)
	case 2:
		gm6 := float64(100 * edakit.Int(5, 50))
		cl := unif(0.5, 10, 0.5)
		p2 := gm6 / (2 * math.Pi * cl)
		q := edakit.Sprintf("After pole splitting, the non-dominant pole of a two-stage op-amp sits near $g_{m6}/C_L$ (with $C_c$ much larger than the first-stage node capacitance). For $g_{m6} = %s\\ \\mu$S and $C_L = %s$ pF, what is this pole in MHz? Round to 2 decimal places.", f(gm6, 0), f(cl, 1))
		return num(q, f(p2, 2), edakit.Sprintf("$f_{p2} \\approx g_{m6}/(2\\pi C_L) = %s/(2\\pi \\cdot %s) = %s$ MHz. Raising $g_{m6}$ pushes it out; more load pulls it in.", f(gm6, 0), f(cl, 1), f(p2, 2)), nil)
	default:
		gm6 := float64(100 * edakit.Int(5, 50))
		cc := unif(0.5, 5, 0.1)
		if edakit.Int(0, 1) == 0 {
			z := gm6 / (2 * math.Pi * cc)
			q := edakit.Sprintf("The feedforward path through $C_c$ in a two-stage op-amp creates a right-half-plane zero at $g_{m6}/C_c$. For $g_{m6} = %s\\ \\mu$S and $C_c = %s$ pF, where is the zero, in MHz? Round to 2 decimal places.", f(gm6, 0), f(cc, 1))
			return num(q, f(z, 2), edakit.Sprintf("$f_z = g_{m6}/(2\\pi C_c) = %s$ MHz. A right-half-plane zero adds gain like a left-half-plane zero but subtracts phase like a pole.", f(z, 2)), nil)
		}
		rz := 1000 / gm6
		q := edakit.Sprintf("A nulling resistor $R_z$ in series with $C_c$ moves the zero to $1/(C_c(1/g_{m6} - R_z))$. For $g_{m6} = %s\\ \\mu$S, what $R_z$ pushes the zero to infinity, in kilohms? Round to 3 decimal places.", f(gm6, 0))
		return num(q, f(rz, 3), edakit.Sprintf("The zero goes to infinity when $R_z = 1/g_{m6} = 1/(%s\\ \\mu\\mathrm{S}) = %s$ kilohms. A slightly larger $R_z$ moves it into the left half plane, where it can cancel $p_2$.", f(gm6, 0), f(rz, 3)), nil)
	}
}

func phaseMargin(ctx generator.GeneratorContext) generator.Problem {
	fu := unif(5, 100, 1)
	fp1 := unif(0.001, 0.05, 0.001)
	fp2 := fu * unif(0.8, 4, 0.1)
	switch whichKP(ctx) {
	case 1:
		a1 := deg(math.Atan(fu / fp1))
		a2 := deg(math.Atan(fu / fp2))
		pm := 180 - a1 - a2
		q := edakit.Sprintf("An op-amp loop gain has two left-half-plane poles, at %s MHz and %s MHz, and crosses unity magnitude at $f_u = %s$ MHz. What is the phase margin in degrees? Round to 1 decimal place.", f(fp1, 3), f(fp2, 1), f(fu, 0))
		return num(q, f(pm, 1), edakit.Sprintf("Each pole contributes $-\\arctan(f_u/f_p)$: %s and %s degrees. $\\mathrm{PM} = 180 - %s - %s = %s$ degrees.", f(a1, 2), f(a2, 2), f(a1, 2), f(a2, 2), f(pm, 1)), map[string]string{"pm": f(pm, 1)})
	case 2:
		fz := fu * unif(2, 10, 0.5)
		a1 := deg(math.Atan(fu / fp1))
		a2 := deg(math.Atan(fu / fp2))
		a3 := deg(math.Atan(fu / fz))
		pm := 180 - a1 - a2 - a3
		q := edakit.Sprintf("A two-stage op-amp's loop gain has poles at %s MHz and %s MHz and an uncancelled right-half-plane zero at %s MHz. It crosses unity at $f_u = %s$ MHz. What is the phase margin in degrees? Round to 1 decimal place.", f(fp1, 3), f(fp2, 1), f(fz, 1), f(fu, 0))
		return num(q, f(pm, 1), edakit.Sprintf("Poles subtract %s and %s degrees; the right-half-plane zero subtracts a further $\\arctan(%s/%s) = %s$ degrees. $\\mathrm{PM} = %s$ degrees.", f(a1, 2), f(a2, 2), f(fu, 0), f(fz, 1), f(a3, 2), f(pm, 1)), map[string]string{"pm": f(pm, 1)})
	default:
		pmT := float64(5 * edakit.Int(9, 15))
		gbw := unif(5, 200, 1)
		p2 := gbw / math.Tan((90-pmT)*math.Pi/180)
		q := edakit.Sprintf("Treat an op-amp as a dominant pole contributing $-90$ degrees at crossover plus one non-dominant pole $f_{p2}$, with crossover at the gain-bandwidth product %s MHz. What is the lowest $f_{p2}$, in MHz, that gives a phase margin of %s degrees? Round to 2 decimal places.", f(gbw, 0), f(pmT, 0))
		return num(q, f(p2, 2), edakit.Sprintf("$\\mathrm{PM} = 90 - \\arctan(\\mathrm{GBW}/f_{p2})$, so $f_{p2} = \\mathrm{GBW}/\\tan(90 - \\mathrm{PM}) = %s/\\tan(%s) = %s$ MHz, that is $%s \\times$ GBW.", f(gbw, 0), f(90-pmT, 0), f(p2, 2), f(p2/gbw, 2)), nil)
	}
}

func comparator(ctx generator.GeneratorContext) generator.Problem {
	c := float64(edakit.Int(5, 60))
	gm := float64(50 * edakit.Int(2, 40))
	tau := 1000 * c / gm
	switch whichKP(ctx) {
	case 1:
		vin := unif(0.1, 20, 0.1)
		vout := unif(0.4, 1.0, 0.05)
		t := tau * math.Log(vout*1000/vin)
		q := edakit.Sprintf("A cross-coupled latch has regeneration time constant $\\tau = C/g_m$ with $C = %s$ fF and $g_m = %s\\ \\mu$S. It starts from a %s mV imbalance. How long does it take to regenerate to %s V, in picoseconds? Round to 1 decimal place.", f(c, 0), f(gm, 0), f(vin, 1), f(vout, 2))
		return num(q, f(t, 1), edakit.Sprintf("$\\tau = %s$ ps. The imbalance grows as $e^{t/\\tau}$, so $t = \\tau\\ln(V_{out}/V_{in}) = %s \\cdot \\ln(%s) = %s$ ps.", f(tau, 3), f(tau, 3), f(vout*1000/vin, 2), f(t, 1)), map[string]string{"tau": f(tau, 3)})
	case 2:
		tau := float64(edakit.Int(5, 60))
		tr := math.Round(tau * unif(2, 12, 0.5))
		vout := unif(0.4, 1.0, 0.05)
		vmin := vout * math.Exp(-tr/tau) * 1e6
		q := edakit.Sprintf("A latch with $\\tau = %s$ ps gets %s ps to regenerate to a valid %s V output before the next stage samples it. What is the smallest initial imbalance that resolves in time, in microvolts? Round to 2 decimal places.", f(tau, 3), f(tr, 0), f(vout, 2))
		return num(q, f(vmin, 2), edakit.Sprintf("$V_{min} = V_{out} e^{-T/\\tau} = %s \\cdot e^{-%s} = %s\\ \\mu$V. Smaller inputs risk a metastable decision.", f(vout, 2), f(tr/tau, 3), f(vmin, 2)), nil)
	default:
		vp := unif(0.5, 5, 0.1)
		vl := unif(5, 40, 1)
		ap := float64(edakit.Int(2, 20))
		v := math.Sqrt(vp*vp + (vl/ap)*(vl/ap))
		q := edakit.Sprintf("A comparator puts a preamplifier with gain %s V/V in front of a latch. The preamp's input-referred offset has $\\sigma = %s$ mV and the latch's own offset has $\\sigma = %s$ mV, independent. What is the total input-referred offset sigma, in mV? Round to 2 decimal places.", f(ap, 0), f(vp, 1), f(vl, 0))
		return num(q, f(v, 2), edakit.Sprintf("Refer the latch offset to the input by dividing by the preamp gain, $%s/%s = %s$ mV, then add independent variances: $\\sqrt{%s^2 + %s^2} = %s$ mV.", f(vl, 0), f(ap, 0), f(vl/ap, 3), f(vp, 1), f(vl/ap, 3), f(v, 2)), nil)
	}
}
