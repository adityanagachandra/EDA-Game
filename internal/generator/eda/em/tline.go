package em

import (
	"math"
	"math/cmplx"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerTline(reg *generator.Registry) {
	reg.Register("em.tline.telegrapher", edakit.Func(telegrapher))
	reg.Register("em.tline.propagation_delay", edakit.Func(propagationDelay))
	reg.Register("em.tline.reflection", edakit.Func(reflection))
	reg.Register("em.tline.bounce_diagram", edakit.Func(bounceDiagram))
	reg.Register("em.tline.termination", edakit.Func(termination))
	reg.Register("em.tline.input_impedance", edakit.Func(inputImpedance))
	reg.Register("em.tline.lossy_line", edakit.Func(lossyLine))
	reg.Register("em.tline.s_parameters", edakit.Func(sParameters))
}

// psPerMM is the free-space delay per millimetre, 1/c, in ps/mm.
const psPerMM = 1e9 / cLight

func telegrapher(ctx generator.GeneratorContext) generator.Problem {
	z := float64(5 * edakit.Int(5, 20))
	ee := pick(1.5, 2, 2.7, 3, 3.4, 3.9, 4.3)
	lp := math.Round(z * math.Sqrt(ee) / cLight * 1e9)    // nH/m
	cp := math.Round(math.Sqrt(ee) / (z * cLight) * 1e12) // pF/m
	stem := sp("A lossless line has $L' = %s$ nH/m and $C' = %s$ pF/m.", num(lp, 0), num(cp, 0))
	z0 := math.Sqrt(lp * 1e-9 / (cp * 1e-12))
	v := 1 / math.Sqrt(lp*1e-9*cp*1e-12) // m/s
	switch kp(ctx) {
	case 1:
		l := pick(0.05, 0.1, 0.2, 0.25, 0.5, 1)
		lt, ct := lp*l, cp*l
		td := math.Sqrt(lt*1e-9*ct*1e-12) * 1e12
		return edakit.Numeric(
			sp("A line segment %s m long is measured as having a total series inductance of %s nH and a total shunt capacitance of %s pF, uniformly distributed. What is its one-way propagation delay, in picoseconds? Round to 1 decimal place.", num(l, 2), num(lt, 2), num(ct, 2)),
			num(td, 1),
			sp("Per unit length $L' = %s$ nH/m and $C' = %s$ pF/m, so the delay per metre is $\\sqrt{L'C'}$ and the delay is $\\ell\\sqrt{L'C'} = \\sqrt{L_{tot} C_{tot}} = %s$ ps.", num(lp, 0), num(cp, 0), num(td, 1)),
			facts("td_ps", num(td, 1)))
	case 2:
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(stem+" What is its characteristic impedance, in ohms? Round to 2 decimal places.",
				num(z0, 2),
				sp("$Z_0 = \\sqrt{L'/C'} = \\sqrt{%s \\times 10^{-9}/%s \\times 10^{-12}} = %s$ ohm.", num(lp, 0), num(cp, 0), num(z0, 2)),
				facts("z0", num(z0, 2)))
		}
		return edakit.Numeric(stem+" What is the wave velocity on it, in mm/ns? Round to 1 decimal place.",
			num(v*1e-6, 1),
			sp("$v = 1/\\sqrt{L'C'} = 1/\\sqrt{%s \\times 10^{-9} \\cdot %s \\times 10^{-12}} = %s$ m/s, which is %s mm/ns.", num(lp, 0), num(cp, 0), num(v/1e8, 4)+" \\times 10^8", num(v*1e-6, 1)),
			facts("v_mm_ns", num(v*1e-6, 1)))
	default:
		e := cLight * cLight * lp * 1e-9 * cp * 1e-12
		return edakit.Numeric(stem+" What effective dielectric constant does that imply? Round to 2 decimal places.",
			num(e, 2),
			sp("$v = c/\\sqrt{\\varepsilon_{eff}}$ and $v = 1/\\sqrt{L'C'}$, so $\\varepsilon_{eff} = c^2 L'C' = (3.00 \\times 10^8)^2 \\cdot %s \\times 10^{-9} \\cdot %s \\times 10^{-12} = %s$.", num(lp, 0), num(cp, 0), num(e, 2)),
			facts("eeff", num(e, 2)))
	}
}

func propagationDelay(ctx generator.GeneratorContext) generator.Problem {
	d := edakit.Pick(dielectrics)
	tpd := psPerMM * math.Sqrt(d.er)
	switch kp(ctx) {
	case 1:
		l := float64(5 * edakit.Int(1, 60))
		return edakit.Numeric(
			sp("A %s mm stripline is embedded in %s. What is its one-way propagation delay, in picoseconds? Round to 1 decimal place.", num(l, 0), d.name),
			num(tpd*l, 1),
			sp("$t_{pd} = \\sqrt{\\varepsilon_r}/c = 3.336 \\cdot %s = %s$ ps/mm, so %s mm takes %s ps.", num(math.Sqrt(d.er), 4), num(tpd, 3), num(l, 0), num(tpd*l, 1)),
			facts("tpd_ps_mm", num(tpd, 3)))
	case 2:
		tr := float64(10 * edakit.Int(2, 50))
		l := tr / 6 / tpd
		return edakit.Numeric(
			sp("A team's rule says a stripline in %s needs a transmission-line model once its one-way delay exceeds one sixth of the signal's rise time. For a %s ps rise time, what is the longest line that may be modeled as lumped, in mm? Round to 2 decimal places.", d.name, num(tr, 0)),
			num(l, 2),
			sp("The delay budget is $%s/6 = %s$ ps at %s ps/mm, so the limit is %s mm.", num(tr, 0), num(tr/6, 2), num(tpd, 3), num(l, 2)),
			facts("tpd_ps_mm", num(tpd, 3)))
	default:
		l := float64(5 * edakit.Int(1, 40))
		f := pick(0.5, 1, 2, 2.5, 5, 10)
		td := tpd * l
		th := 360 * f * td / 1000
		return edakit.Numeric(
			sp("A %s mm stripline is embedded in %s. What is its electrical length at %s GHz, in degrees? Round to 1 decimal place.", num(l, 0), d.name, num(f, 1)),
			num(th, 1),
			sp("$T_d = %s$ ps. One period at %s GHz is %s ps, so the line spans $360 f T_d = %s$ degrees.", num(td, 2), num(f, 1), num(1000/f, 1), num(th, 1)),
			facts("td_ps", num(td, 2)))
	}
}

var z0s = []float64{25, 28, 40, 45, 50, 60, 75, 85, 90, 100}

func gamma(zl, z0 float64) float64 { return (zl - z0) / (zl + z0) }

func reflection(ctx generator.GeneratorContext) generator.Problem {
	z0 := edakit.Pick(z0s)
	zl := float64(5 * edakit.Int(0, 50))
	switch kp(ctx) {
	case 1:
		g := gamma(zl, z0)
		return edakit.Numeric(
			sp("A %s ohm line is terminated in a %s ohm resistor. What is the voltage reflection coefficient at the load? Round to 3 decimal places.", num(z0, 0), num(zl, 0)),
			num(g, 3),
			sp("$\\Gamma = (Z_L - Z_0)/(Z_L + Z_0) = (%s - %s)/(%s + %s) = %s$.", num(zl, 0), num(z0, 0), num(zl, 0), num(z0, 0), num(g, 3)),
			facts("gamma", num(g, 3)))
	case 2:
		vs := edakit.Pick(supplies)
		rs := float64(5 * edakit.Int(1, 20))
		if zl == 0 {
			zl = 10
		}
		vi := vs * z0 / (rs + z0)
		g := gamma(zl, z0)
		vl := vi * (1 + g)
		return edakit.Numeric(
			sp("A %s V step with %s ohm source resistance drives a long %s ohm line terminated in %s ohm. What voltage first appears at the load when the wave arrives, in volts? Round to 3 decimal places.", num(vs, 1), num(rs, 0), num(z0, 0), num(zl, 0)),
			num(vl, 3),
			sp("The source launches $V_s Z_0/(R_s + Z_0) = %s$ V; it cannot see the load yet. At the load $\\Gamma_L = %s$, so the load sees $%s(1 + (%s)) = %s$ V.", num(vi, 3), num(g, 3), num(vi, 3), num(g, 3), num(vl, 3)),
			facts("v_launch", num(vi, 3), "gamma_l", num(g, 3)))
	default:
		if zl == z0 {
			zl = z0 + 15
		}
		if zl == 0 {
			zl = 5
		}
		g := math.Abs(gamma(zl, z0))
		if edakit.Int(0, 1) == 0 {
			sw := (1 + g) / (1 - g)
			return edakit.Numeric(
				sp("A %s ohm line drives a %s ohm resistive load with a sinusoid. What is the voltage standing wave ratio (VSWR) on the line? Round to 3 decimal places.", num(z0, 0), num(zl, 0)),
				num(sw, 3),
				sp("$|\\Gamma| = %s$, so $\\mathrm{VSWR} = (1 + |\\Gamma|)/(1 - |\\Gamma|) = %s$; for a resistive load it equals $\\max(Z_L, Z_0)/\\min(Z_L, Z_0)$.", num(g, 4), num(sw, 3)),
				facts("abs_gamma", num(g, 4)))
		}
		rl := -20 * math.Log10(g)
		return edakit.Numeric(
			sp("A %s ohm line drives a %s ohm resistive load. What is the return loss at the load, in dB (a positive number)? Round to 2 decimal places.", num(z0, 0), num(zl, 0)),
			num(rl, 2),
			sp("$|\\Gamma| = %s$ and return loss is $-20\\log_{10}|\\Gamma| = %s$ dB. Larger return loss means a better match.", num(g, 4), num(rl, 2)),
			facts("abs_gamma", num(g, 4)))
	}
}

func bounceDiagram(ctx generator.GeneratorContext) generator.Problem {
	vs := edakit.Pick(supplies)
	z0 := edakit.Pick(z0s)
	rs := float64(5 * edakit.Int(1, 30))
	rl := float64(5 * edakit.Int(1, 60))
	if rl == z0 {
		rl += 50
	}
	td := float64(10 * edakit.Int(1, 50))
	vi := vs * z0 / (rs + z0)
	gl, gs := gamma(rl, z0), gamma(rs, z0)
	stem := sp("A %s V step (applied at $t = 0$) with source resistance %s ohm drives a %s ohm line of one-way delay %s ps, terminated in %s ohm.", num(vs, 1), num(rs, 0), num(z0, 0), num(td, 0), num(rl, 0))
	f := facts("v_launch", num(vi, 4), "gamma_l", num(gl, 4), "gamma_s", num(gs, 4))
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+" What voltage appears at the source end of the line just after $t = 0$, in volts? Round to 3 decimal places.",
			num(vi, 3),
			sp("Until a reflection returns, the line looks like a resistor $Z_0$, so the source end divides: $%s \\cdot %s/(%s + %s) = %s$ V.", num(vs, 1), num(z0, 0), num(rs, 0), num(z0, 0), num(vi, 3)), f)
	case 2:
		k := edakit.Int(1, 3)
		v, w := 0.0, vi
		for j := 0; j < k; j++ {
			v += w * (1 + gl)
			w *= gl * gs
		}
		t := float64(2*k-1) * td
		return edakit.Numeric(stem+sp(" What is the load voltage just after $t = %s$ ps, in volts? Round to 3 decimal places.", num(t, 0)),
			num(v, 3),
			sp("Launched wave %s V, $\\Gamma_L = %s$, $\\Gamma_S = %s$. Each arrival at the load adds $(1 + \\Gamma_L)$ times the incident wave, and the incident wave shrinks by $\\Gamma_L\\Gamma_S$ per round trip. After arrival %d the load sits at %s V (steady state %s V).", num(vi, 3), num(gl, 3), num(gs, 3), k, num(v, 3), num(vs*rl/(rs+rl), 3)), f)
	default:
		v := vi * (1 + gl + gl*gs)
		return edakit.Numeric(stem+sp(" What is the voltage at the source end of the line just after $t = %s$ ps, in volts? Round to 3 decimal places.", num(2*td, 0)),
			num(v, 3),
			sp("At $2T_d$ the reflection $\\Gamma_L V_i = %s$ V returns and is itself reflected by the source with $\\Gamma_S = %s$. The source end holds $V_i(1 + \\Gamma_L + \\Gamma_L\\Gamma_S) = %s$ V.", num(gl*vi, 3), num(gs, 3), num(v, 3)), f)
	}
}

func termination(ctx generator.GeneratorContext) generator.Problem {
	z0 := edakit.Pick(z0s)
	switch kp(ctx) {
	case 1:
		rd := float64(edakit.Int(5, int(z0)-3))
		return edakit.Numeric(
			sp("A driver with %s ohm output resistance drives a %s ohm point-to-point line to a high-impedance receiver. What series resistor at the driver makes the source matched, in ohms?", num(rd, 0), num(z0, 0)),
			num(z0-rd, 0),
			sp("Series termination makes $R_{drv} + R_s = Z_0$, so $R_s = %s - %s = %s$ ohm. The launched wave is half the swing, doubles at the open far end, and the returning wave is absorbed at the source.", num(z0, 0), num(rd, 0), num(z0-rd, 0)),
			facts("rs", num(z0-rd, 0)))
	case 2:
		v := edakit.Pick(supplies)
		p := v * v / z0 * 1000
		return edakit.Numeric(
			sp("A %s ohm line is parallel-terminated at the receiver by a %s ohm resistor to ground. While the driver holds the line at %s V, how much DC power does the terminator dissipate, in mW? Round to 2 decimal places.", num(z0, 0), num(z0, 0), num(v, 1)),
			num(p, 2),
			sp("$P = V^2/R = %s^2/%s = %s$ mW, paid for as long as the line is high. That static power is the cost of parallel termination.", num(v, 1), num(z0, 0), num(p, 2)),
			facts("p_mw", num(p, 2)))
	default:
		vdd := pick(1.8, 2.5, 3.3)
		vt := pick(0.5, 0.6, 0.75, 0.9, 1.2)
		r1, r2 := z0*vdd/vt, z0*vdd/(vdd-vt)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A %s ohm line gets a Thevenin termination: $R_1$ from the receiver to the %s V supply and $R_2$ from the receiver to ground, chosen so that $R_1 \\parallel R_2 = Z_0$ and the open-circuit voltage is %s V. What is $R_1$, in ohms? Round to 1 decimal place.", num(z0, 0), num(vdd, 1), num(vt, 2)),
				num(r1, 1),
				sp("From $R_1 \\parallel R_2 = Z_0$ and $V_{dd} R_2/(R_1 + R_2) = V_t$: $R_1 = Z_0 V_{dd}/V_t = %s$ ohm and $R_2 = Z_0 V_{dd}/(V_{dd} - V_t) = %s$ ohm.", num(r1, 1), num(r2, 1)),
				facts("r1", num(r1, 1), "r2", num(r2, 1)))
		}
		return edakit.Numeric(
			sp("A %s ohm line gets a Thevenin termination: $R_1$ from the receiver to the %s V supply and $R_2$ from the receiver to ground, chosen so that $R_1 \\parallel R_2 = Z_0$ and the open-circuit voltage is %s V. What is $R_2$, in ohms? Round to 1 decimal place.", num(z0, 0), num(vdd, 1), num(vt, 2)),
			num(r2, 1),
			sp("From $R_1 \\parallel R_2 = Z_0$ and $V_{dd} R_2/(R_1 + R_2) = V_t$: $R_1 = Z_0 V_{dd}/V_t = %s$ ohm and $R_2 = Z_0 V_{dd}/(V_{dd} - V_t) = %s$ ohm.", num(r1, 1), num(r2, 1)),
			facts("r1", num(r1, 1), "r2", num(r2, 1)))
	}
}

func inputImpedance(ctx generator.GeneratorContext) generator.Problem {
	z0 := edakit.Pick(z0s)
	switch kp(ctx) {
	case 1:
		zl := float64(5 * edakit.Int(1, 40))
		bl := float64(5 * edakit.Pick([]int{2, 3, 4, 6, 8, 9, 10, 12, 14, 15, 20, 22, 24, 27, 30, 33}))
		t := math.Tan(bl * math.Pi / 180)
		zin := complex(z0, 0) * (complex(zl, z0*t)) / complex(z0, zl*t)
		return edakit.Numeric(
			sp("A lossless %s ohm line with electrical length $\\beta\\ell = %s^\\circ$ is terminated in %s ohm. What is the magnitude of its input impedance, in ohms? Round to 2 decimal places.", num(z0, 0), num(bl, 0), num(zl, 0)),
			num(cmplx.Abs(zin), 2),
			sp("$Z_{in} = Z_0\\frac{Z_L + jZ_0\\tan\\beta\\ell}{Z_0 + jZ_L\\tan\\beta\\ell}$ with $\\tan %s^\\circ = %s$ gives $Z_{in} = %s %s j%s$ ohm, magnitude %s ohm.", num(bl, 0), num(t, 4), num(real(zin), 2), plusMinus(imag(zin)), num(math.Abs(imag(zin)), 2), num(cmplx.Abs(zin), 2)),
			facts("re", num(real(zin), 2), "im", num(imag(zin), 2)))
	case 2:
		if edakit.Int(0, 1) == 0 {
			zl := float64(5 * edakit.Int(2, 40))
			return edakit.Numeric(
				sp("A quarter-wavelength section of %s ohm line is terminated in %s ohm. What is its input impedance at the design frequency, in ohms? Round to 2 decimal places.", num(z0, 0), num(zl, 0)),
				num(z0*z0/zl, 2),
				sp("At $\\beta\\ell = 90^\\circ$ the tangent diverges and $Z_{in} = Z_0^2/Z_L = %s^2/%s = %s$ ohm.", num(z0, 0), num(zl, 0), num(z0*z0/zl, 2)),
				facts("zin", num(z0*z0/zl, 2)))
		}
		z1, z2 := float64(5*edakit.Int(2, 20)), float64(5*edakit.Int(2, 20))
		if z1 == z2 {
			z2 += 25
		}
		return edakit.Numeric(
			sp("A quarter-wave transformer must match a %s ohm line to a %s ohm load at one frequency. What characteristic impedance must the quarter-wave section have, in ohms? Round to 2 decimal places.", num(z1, 0), num(z2, 0)),
			num(math.Sqrt(z1*z2), 2),
			sp("A quarter-wave section maps $Z_L$ to $Z_T^2/Z_L$; setting that equal to %s ohm gives $Z_T = \\sqrt{%s \\cdot %s} = %s$ ohm.", num(z1, 0), num(z1, 0), num(z2, 0), num(math.Sqrt(z1*z2), 2)),
			facts("zt", num(math.Sqrt(z1*z2), 2)))
	default:
		bl := float64(5 * edakit.Pick([]int{2, 3, 4, 6, 8, 9, 10, 12, 14, 15, 20, 21, 24, 27, 30, 33}))
		r := bl * math.Pi / 180
		if edakit.Int(0, 1) == 0 {
			x := z0 * math.Tan(r)
			return edakit.Numeric(
				sp("A lossless %s ohm stub with electrical length %s degrees is shorted at its far end. What is its input reactance, in ohms (positive for inductive)? Round to 2 decimal places.", num(z0, 0), num(bl, 0)),
				num(x, 2),
				sp("With $Z_L = 0$ the input impedance formula reduces to $Z_{in} = jZ_0\\tan\\beta\\ell$, so $X = %s\\tan %s^\\circ = %s$ ohm.", num(z0, 0), num(bl, 0), num(x, 2)),
				facts("x", num(x, 2)))
		}
		x := -z0 / math.Tan(r)
		return edakit.Numeric(
			sp("A lossless %s ohm stub with electrical length %s degrees is open at its far end. What is its input reactance, in ohms (negative for capacitive)? Round to 2 decimal places.", num(z0, 0), num(bl, 0)),
			num(x, 2),
			sp("With $Z_L \\to \\infty$ the formula reduces to $Z_{in} = -jZ_0\\cot\\beta\\ell$, so $X = -%s\\cot %s^\\circ = %s$ ohm.", num(z0, 0), num(bl, 0), num(x, 2)),
			facts("x", num(x, 2)))
	}
}

func plusMinus(x float64) string {
	if x < 0 {
		return "-"
	}
	return "+"
}

func lossyLine(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		z0 := edakit.Pick(z0s)
		r := float64(5 * edakit.Int(1, 40))
		g := pick(0, 0, 0.0005, 0.001, 0.002)
		l := pick(0.1, 0.2, 0.25, 0.5, 1)
		a := r/(2*z0) + g*z0/2
		db := 8.686 * a * l
		gtxt := "negligible shunt conductance"
		if g > 0 {
			gtxt = sp("$G' = %s$ S/m", num(g, 4))
		}
		return edakit.Numeric(
			sp("A low-loss %s ohm line has $R' = %s$ ohm/m and %s. How much does a signal attenuate over %s m, in dB? Round to 3 decimal places.", num(z0, 0), num(r, 0), gtxt, num(l, 2)),
			num(db, 3),
			sp("$\\alpha \\approx R'/(2Z_0) + G'Z_0/2 = %s$ Np/m. Over %s m that is %s Np, or %s dB.", num(a, 4), num(l, 2), num(a*l, 4), num(db, 3)),
			facts("alpha_np_m", num(a, 4)))
	case 2:
		a := pick(0.1, 0.2, 0.5, 1, 2)
		l := pick(0.1, 0.2, 0.3, 0.5, 1)
		v := float64(100 * edakit.Int(2, 10))
		out := v * math.Exp(-a*l)
		return edakit.Numeric(
			sp("A line has an attenuation constant of %s Np/m. A %s mV sinusoid enters a %s m length. What amplitude leaves it, in mV? Round to 2 decimal places.", num(a, 1), num(v, 0), num(l, 1)),
			num(out, 2),
			sp("Amplitude decays as $e^{-\\alpha\\ell} = e^{-%s} = %s$, so %s mV leaves; the loss is $8.686 \\cdot %s = %s$ dB.", num(a*l, 3), num(math.Exp(-a*l), 4), num(out, 2), num(a*l, 3), num(8.686*a*l, 2)),
			facts("loss_db", num(8.686*a*l, 2)))
	default:
		f0 := pick(1, 2, 4, 5)
		f1 := f0 * pick(2, 4, 5, 8, 10)
		ac, ad := pick(1, 2, 3, 4), pick(0.5, 1, 1.5, 2)
		tot := ac*math.Sqrt(f1/f0) + ad*f1/f0
		return edakit.Numeric(
			sp("At %s GHz a channel's loss is %s dB from the conductor (skin effect) and %s dB from the dielectric. Assuming conductor loss scales as $\\sqrt f$ and dielectric loss as $f$, what is the total loss at %s GHz, in dB? Round to 2 decimal places.", num(f0, 0), num(ac, 1), num(ad, 1), num(f1, 0)),
			num(tot, 2),
			sp("Conductor: $%s\\sqrt{%s/%s} = %s$ dB. Dielectric: $%s \\cdot %s/%s = %s$ dB. Total %s dB.", num(ac, 1), num(f1, 0), num(f0, 0), num(ac*math.Sqrt(f1/f0), 2), num(ad, 1), num(f1, 0), num(f0, 0), num(ad*f1/f0, 2), num(tot, 2)),
			facts("total_db", num(tot, 2)))
	}
}

func sParameters(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		zl := float64(5 * edakit.Int(1, 40))
		if zl == 50 {
			zl = 55
		}
		g := math.Abs(gamma(zl, 50))
		db := 20 * math.Log10(g)
		return edakit.Numeric(
			sp("A one-port presents %s ohm to a 50 ohm reference system. What is $|S_{11}|$, in dB? Round to 2 decimal places.", num(zl, 0)),
			num(db, 2),
			sp("$S_{11}$ is the reflection coefficient against the reference impedance: $|S_{11}| = |%s - 50|/(%s + 50) = %s$, which is $20\\log_{10}(%s) = %s$ dB.", num(zl, 0), num(zl, 0), num(g, 4), num(g, 4), num(db, 2)),
			facts("abs_s11", num(g, 4)))
	case 2:
		if edakit.Int(0, 1) == 0 {
			s := float64(edakit.Int(10, 99)) / 100
			return edakit.Numeric(
				sp("A channel has $|S_{21}| = %s$ at its Nyquist frequency. What is its insertion loss expressed as $|S_{21}|$ in dB? Round to 2 decimal places.", num(s, 2)),
				num(20*math.Log10(s), 2),
				sp("S-parameters are ratios of wave amplitudes, so dB is $20\\log_{10}|S_{21}| = 20\\log_{10}(%s) = %s$ dB.", num(s, 2), num(20*math.Log10(s), 2)),
				facts("db", num(20*math.Log10(s), 2)))
		}
		pin := pick(1, 2, 5, 10)
		pout := pin * float64(edakit.Int(5, 95)) / 100
		return edakit.Numeric(
			sp("A matched two-port receives %s mW of incident power at port 1 and delivers %s mW into a matched load at port 2. What is $|S_{21}|$, in dB? Round to 2 decimal places.", num(pin, 0), num(pout, 2)),
			num(10*math.Log10(pout/pin), 2),
			sp("$|S_{21}|^2$ is the power ratio, so $|S_{21}|$ in dB is $10\\log_{10}(%s/%s) = %s$ dB.", num(pout, 2), num(pin, 0), num(10*math.Log10(pout/pin), 2)),
			facts("db", num(10*math.Log10(pout/pin), 2)))
	default:
		if edakit.Int(0, 1) == 0 {
			s11 := float64(edakit.Int(5, 50)) / 100
			s21 := float64(edakit.Int(30, 95)) / 100
			for s11*s11+s21*s21 > 1 {
				s21 -= 0.05
			}
			lost := 100 * (1 - s11*s11 - s21*s21)
			return edakit.Numeric(
				sp("A two-port driven at port 1 (port 2 matched) has $|S_{11}| = %s$ and $|S_{21}| = %s$. What percentage of the incident power is dissipated inside the network? Round to 1 decimal place.", num(s11, 2), num(s21, 2)),
				num(lost, 1),
				sp("Incident power is reflected ($|S_{11}|^2 = %s$), transmitted ($|S_{21}|^2 = %s$) or absorbed. Absorbed: $1 - %s - %s = %s$, so %s percent.", num(s11*s11, 4), num(s21*s21, 4), num(s11*s11, 4), num(s21*s21, 4), num(lost/100, 4), num(lost, 1)),
				facts("lost_pct", num(lost, 1)))
		}
		rl := float64(edakit.Int(3, 30))
		s11 := math.Pow(10, -rl/20)
		s21 := math.Sqrt(1 - s11*s11)
		return edakit.Numeric(
			sp("A lossless, reciprocal two-port has $|S_{11}| = -%s$ dB with port 2 matched. What is $|S_{21}|$, in dB? Round to 3 decimal places.", num(rl, 0)),
			num(20*math.Log10(s21), 3),
			sp("$|S_{11}| = 10^{-%s/20} = %s$. A lossless network has a unitary S-matrix, so $|S_{21}|^2 = 1 - |S_{11}|^2 = %s$ and $|S_{21}| = %s$ dB.", num(rl, 0), num(s11, 4), num(s21*s21, 4), num(20*math.Log10(s21), 3)),
			facts("s21", num(s21, 4)))
	}
}
