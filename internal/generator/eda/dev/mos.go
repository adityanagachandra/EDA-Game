package dev

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// coxF is the oxide capacitance in F/cm^2 for an SiO2-equivalent thickness in nm.
func coxF(toxNm float64) float64 { return epsOx / (toxNm * 1e-7) }

// mosCapacitor asks for an oxide capacitance (KP 1), an equivalent oxide thickness (KP 2) or
// the inversion sheet density (KP 3).
func mosCapacitor(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		tox := pickF(1.0, 1.2, 1.5, 1.8, 2.0, 2.5, 3.0, 4.0, 5.0)
		c := coxF(tox) * 1e7
		ans, rt := sigAnswer(c, 3)
		q := edakit.Sprintf("A MOS gate stack has an SiO$_2$ dielectric $t_{ox} = %s$ nm thick. With $\\varepsilon_{ox} = 3.9\\varepsilon_0$ and $\\varepsilon_0 = 8.854 \\times 10^{-14}$ F/cm, what is the oxide capacitance per unit area, in fF/µm$^2$? %s", edakit.Num(tox, 1), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$C_{ox} = \\varepsilon_{ox}/t_{ox} = 3.453 \\times 10^{-13} / (%s \\times 10^{-7}) = %.4g$ F/cm$^2$. Since 1 F/cm$^2$ equals $10^7$ fF/µm$^2$, $C_{ox} = %s$ fF/µm$^2$.", edakit.Num(tox, 1), coxF(tox), ans),
			map[string]string{"cox": ans})
	case 2:
		k := pickF(9, 16, 20, 22, 25)
		t := pickF(1.5, 2.0, 2.5, 3.0, 3.5, 4.0)
		eot := t * 3.9 / k
		ans := edakit.Num(eot, 3)
		q := edakit.Sprintf("A high-k gate dielectric with relative permittivity %s is %s nm thick. Ignoring any interfacial layer, what is its equivalent oxide thickness (EOT), in nm? Round to 3 decimal places.", edakit.Num(k, 0), edakit.Num(t, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("EOT is the SiO$_2$ thickness with the same capacitance: $t_{phys} \\times 3.9/k = %s \\times 3.9/%s = %s$ nm.", edakit.Num(t, 1), edakit.Num(k, 0), ans),
			map[string]string{"eot": ans})
	default:
		eot := pickF(0.8, 0.9, 1.0, 1.2, 1.5, 2.0)
		vov := 0.05 * float64(edakit.Int(4, 14))
		n := coxF(eot) * vov / qC / 1e12
		ans := edakit.Num(n, 2)
		q := edakit.Sprintf("A transistor has EOT = %s nm ($\\varepsilon_{ox} = 3.453 \\times 10^{-13}$ F/cm) and is biased %s V above threshold. Using the charge-sheet approximation $Q_{inv} = C_{ox}(V_{GS} - V_t)$ and $q = 1.602 \\times 10^{-19}$ C, what is the inversion electron density, in units of $10^{12}$ cm$^{-2}$? Round to 2 decimal places.",
			edakit.Num(eot, 1), edakit.Num(vov, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$C_{ox} = %.4g$ F/cm$^2$, $Q_{inv} = %.4g$ C/cm$^2$, and $N_{inv} = Q_{inv}/q = %.4g$ cm$^{-2}$, which is %s in units of $10^{12}$ cm$^{-2}$.", coxF(eot), coxF(eot)*vov, n*1e12, ans),
			map[string]string{"Ninv": ans})
	}
}

// threshold asks for 2 phi_F (KP 1), the depletion-charge term (KP 2) or the full threshold
// voltage (KP 3) of a bulk NMOS.
func threshold(ctx generator.GeneratorContext) generator.Problem {
	m, e := pickF(1, 2, 3, 5), edakit.Int(16, 17)
	na := m * math.Pow(10, float64(e))
	tox := pickF(1.5, 2, 3, 4, 5)
	phi2 := 2 * vT300 * math.Log(na/1e10)
	consts := "Use $kT/q = 0.02585$ V, $n_i = 10^{10}$ cm$^{-3}$, $\\varepsilon_{Si} = 1.036 \\times 10^{-12}$ F/cm, $\\varepsilon_{ox} = 3.453 \\times 10^{-13}$ F/cm and $q = 1.602 \\times 10^{-19}$ C."
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(phi2, 3)
		q := edakit.Sprintf("Strong inversion in a bulk NMOS starts when the surface potential reaches $2\\phi_F$. For a p-type body with $N_A = %s$ cm$^{-3}$ at 300 K ($kT/q = 0.02585$ V, $n_i = 10^{10}$ cm$^{-3}$), what is $2\\phi_F$, in volts? Round to 3 decimal places.", sci(m, e))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\phi_F = (kT/q)\\ln(N_A/n_i) = 0.02585 \\times %.3f = %.4f$ V, so $2\\phi_F = %s$ V.", math.Log(na/1e10), phi2/2, ans),
			map[string]string{"2phiF": ans})
	case 2:
		p := edakit.Round(phi2, 2)
		qd := math.Sqrt(2 * epsSi * qC * na * p)
		term := qd / coxF(tox)
		ans, rt := sigAnswer(term, 3)
		q := edakit.Sprintf("A bulk NMOS has $N_A = %s$ cm$^{-3}$, $2\\phi_F = %.2f$ V and $t_{ox} = %s$ nm. What is the depletion-charge term of the threshold voltage, $\\sqrt{2\\varepsilon_{Si} q N_A (2\\phi_F)}/C_{ox}$, in volts? %s %s",
			sci(m, e), p, edakit.Num(tox, 1), consts, rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$Q_{dep} = \\sqrt{2 \\varepsilon_{Si} q N_A (2\\phi_F)} = %.4g$ C/cm$^2$ and $C_{ox} = %.4g$ F/cm$^2$, so the term is %s V.", qd, coxF(tox), ans),
			map[string]string{"Qdep": edakit.Sprintf("%.4g", qd), "Cox": edakit.Sprintf("%.4g", coxF(tox))})
	default:
		m, e = pickF(1, 2, 3, 5), 17
		na = m * 1e17
		tox = pickF(2, 3, 4, 5)
		phi2 = 2 * vT300 * math.Log(na/1e10)
		vfb := -pickF(0.8, 0.85, 0.9, 0.95, 1.0)
		qd := math.Sqrt(2 * epsSi * qC * na * phi2)
		vt := vfb + phi2 + qd/coxF(tox)
		ans := edakit.Num(vt, 3)
		q := edakit.Sprintf("A long-channel bulk NMOS has $N_A = %s$ cm$^{-3}$, $t_{ox} = %s$ nm and flat-band voltage $V_{FB} = %s$ V. With $V_t = V_{FB} + 2\\phi_F + \\sqrt{2\\varepsilon_{Si} q N_A (2\\phi_F)}/C_{ox}$, what is the threshold voltage, in volts? %s Round to 3 decimal places.",
			sci(m, e), edakit.Num(tox, 1), edakit.Num(vfb, 2), consts)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$2\\phi_F = %.4f$ V. The depletion term is $%.4g / %.4g = %.4f$ V. So $V_t = %s + %.4f + %.4f = %s$ V.", phi2, qd, coxF(tox), qd/coxF(tox), edakit.Num(vfb, 2), phi2, qd/coxF(tox), ans),
			map[string]string{"2phiF": edakit.Num(phi2, 4), "dep_term": edakit.Num(qd/coxF(tox), 4)})
	}
}

// bodyEffect asks for a threshold shift (KP 1), a shifted threshold (KP 2) or the body-effect
// coefficient (KP 3).
func bodyEffect(ctx generator.GeneratorContext) generator.Problem {
	gamma := 0.05 * float64(edakit.Int(4, 12))
	phi := pickF(0.7, 0.75, 0.8, 0.85, 0.9)
	vsb := 0.1 * float64(edakit.Int(2, 12))
	dvt := gamma * (math.Sqrt(phi+vsb) - math.Sqrt(phi))
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(dvt*1000, 1)
		q := edakit.Sprintf("An NMOS has body-effect coefficient $\\gamma = %s$ V$^{1/2}$ and $2\\phi_F = %s$ V. Its source sits %s V above the body. By how much does the threshold voltage rise, in mV? Round to 1 decimal place.",
			edakit.Num(gamma, 2), edakit.Num(phi, 2), edakit.Num(vsb, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\Delta V_t = \\gamma(\\sqrt{2\\phi_F + V_{SB}} - \\sqrt{2\\phi_F}) = %s \\times (%.4f - %.4f) = %s$ mV.", edakit.Num(gamma, 2), math.Sqrt(phi+vsb), math.Sqrt(phi), ans),
			map[string]string{"dVt_mV": ans})
	case 2:
		vt0 := 0.01 * float64(edakit.Int(25, 50))
		ans := edakit.Num(vt0+dvt, 3)
		q := edakit.Sprintf("In a two-high NMOS stack, the upper transistor's source rises to %s V while its body stays at ground. Its zero-bias threshold is $V_{t0} = %s$ V, with $\\gamma = %s$ V$^{1/2}$ and $2\\phi_F = %s$ V. What is its threshold voltage, in volts? Round to 3 decimal places.",
			edakit.Num(vsb, 1), edakit.Num(vt0, 2), edakit.Num(gamma, 2), edakit.Num(phi, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$V_t = V_{t0} + \\gamma(\\sqrt{2\\phi_F + V_{SB}} - \\sqrt{2\\phi_F}) = %s + %.4f = %s$ V. This is one reason the upper device of a stack is weaker.", edakit.Num(vt0, 2), dvt, ans),
			map[string]string{"dVt": edakit.Num(dvt, 4)})
	default:
		m, e := pickF(1, 2, 3, 5), edakit.Int(16, 17)
		na := m * math.Pow(10, float64(e))
		tox := pickF(1.5, 2, 3, 4, 5)
		g := math.Sqrt(2*qC*epsSi*na) / coxF(tox)
		ans, rt := sigAnswer(g, 3)
		q := edakit.Sprintf("A bulk NMOS has body doping $N_A = %s$ cm$^{-3}$ and $t_{ox} = %s$ nm. Using $\\varepsilon_{Si} = 1.036 \\times 10^{-12}$ F/cm, $\\varepsilon_{ox} = 3.453 \\times 10^{-13}$ F/cm and $q = 1.602 \\times 10^{-19}$ C, what is the body-effect coefficient $\\gamma = \\sqrt{2q\\varepsilon_{Si}N_A}/C_{ox}$, in V$^{1/2}$? %s",
			sci(m, e), edakit.Num(tox, 1), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\sqrt{2q\\varepsilon_{Si}N_A} = %.4g$ and $C_{ox} = %.4g$ F/cm$^2$, so $\\gamma = %s$ V$^{1/2}$. Heavier doping and thicker oxide both strengthen the body effect.", math.Sqrt(2*qC*epsSi*na), coxF(tox), ans),
			map[string]string{"gamma": ans})
	}
}

type mosBias struct {
	kp, wl, vt, vgs float64
}

func drawBias() mosBias {
	return mosBias{
		kp:  float64(50 * edakit.Int(2, 10)),
		wl:  float64(edakit.Pick([]int{1, 2, 4, 5, 8, 10, 16, 20})),
		vt:  0.05 * float64(edakit.Int(6, 9)),
		vgs: 0.1 * float64(edakit.Int(8, 12)),
	}
}

func (b mosBias) stem() string {
	return edakit.Sprintf("an NMOS with $k' = \\mu_n C_{ox} = %d$ µA/V$^2$, $W/L = %d$, $V_t = %s$ V and $V_{GS} = %s$ V", int(b.kp), int(b.wl), edakit.Num(b.vt, 2), edakit.Num(b.vgs, 1))
}

// linearRegion asks for a deep-triode resistance (KP 1), a linear-region current (KP 2) or a
// current that first needs the region identified (KP 3).
func linearRegion(ctx generator.GeneratorContext) generator.Problem {
	b := drawBias()
	vov := b.vgs - b.vt
	switch kpOf(ctx) {
	case 1:
		r := 1 / (b.kp * 1e-6 * b.wl * vov) / 1000
		ans, rt := sigAnswer(r, 3)
		q := edakit.Sprintf("Used as a pass switch at very small $V_{DS}$, %s behaves like a resistor. What is its on-resistance, in kilohms? %s", b.stem(), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("For $V_{DS} \\ll V_{GS} - V_t$, $I_D \\approx k'(W/L)(V_{GS} - V_t)V_{DS}$, so $R_{on} = 1/[k'(W/L)(V_{GS}-V_t)] = 1/(%d \\times 10^{-6} \\times %d \\times %s) = %s$ kilohms.", int(b.kp), int(b.wl), edakit.Num(vov, 2), ans),
			map[string]string{"Ron_kohm": ans})
	case 2:
		vds := 0.05 * float64(edakit.Int(1, int(math.Floor(vov/0.05))-1))
		id := b.kp * b.wl * (vov*vds - vds*vds/2)
		ans := edakit.Num(id, 1)
		q := edakit.Sprintf("Find the drain current of %s and $V_{DS} = %s$ V, using the long-channel linear-region model. Answer in µA. Round to 1 decimal place.", b.stem(), edakit.Num(vds, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$V_{DS} = %s < V_{GS} - V_t = %s$, so the device is linear: $I_D = k'(W/L)[(V_{GS}-V_t)V_{DS} - V_{DS}^2/2] = %d \\times %d \\times %.5f = %s$ µA.", edakit.Num(vds, 2), edakit.Num(vov, 2), int(b.kp), int(b.wl), vov*vds-vds*vds/2, ans),
			map[string]string{"Id_uA": ans})
	default:
		vds := 0.1 * float64(edakit.Int(1, 12))
		var id float64
		var why string
		if vds < vov {
			id = b.kp * b.wl * (vov*vds - vds*vds/2)
			why = edakit.Sprintf("$V_{DS} = %s$ V is below $V_{GS} - V_t = %s$ V, so the linear formula applies", edakit.Num(vds, 1), edakit.Num(vov, 2))
		} else {
			id = 0.5 * b.kp * b.wl * vov * vov
			why = edakit.Sprintf("$V_{DS} = %s$ V is at or above $V_{GS} - V_t = %s$ V, so the channel is pinched off and $I_D = \\tfrac12 k'(W/L)(V_{GS}-V_t)^2$", edakit.Num(vds, 1), edakit.Num(vov, 2))
		}
		ans := edakit.Num(id, 1)
		q := edakit.Sprintf("Find the drain current of %s and $V_{DS} = %s$ V, using the long-channel model and ignoring channel-length modulation. Answer in µA. Round to 1 decimal place.", b.stem(), edakit.Num(vds, 1))
		return edakit.Numeric(q, ans, why+edakit.Sprintf(". $I_D = %s$ µA.", ans), map[string]string{"Id_uA": ans})
	}
}

// squareLaw asks for a saturation current (KP 1), the W/L for a target current (KP 2) or a
// transconductance (KP 3).
func squareLaw(ctx generator.GeneratorContext) generator.Problem {
	b := drawBias()
	vov := b.vgs - b.vt
	switch kpOf(ctx) {
	case 1:
		id := 0.5 * b.kp * b.wl * vov * vov
		ans := edakit.Num(id, 1)
		q := edakit.Sprintf("What is the saturation drain current of %s, using the long-channel square law and ignoring channel-length modulation? Answer in µA. Round to 1 decimal place.", b.stem())
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$I_{D,sat} = \\tfrac12 k'(W/L)(V_{GS}-V_t)^2 = 0.5 \\times %d \\times %d \\times %s^2 = %s$ µA.", int(b.kp), int(b.wl), edakit.Num(vov, 2), ans),
			map[string]string{"Idsat_uA": ans})
	case 2:
		target := float64(10 * edakit.Int(5, 60))
		wl := 2 * target / (b.kp * vov * vov)
		ans := edakit.Num(wl, 2)
		q := edakit.Sprintf("A designer needs %d µA of saturation current from an NMOS with $k' = %d$ µA/V$^2$, $V_t = %s$ V and $V_{GS} = %s$ V. Using the long-channel square law, what $W/L$ is required? Round to 2 decimal places.",
			int(target), int(b.kp), edakit.Num(b.vt, 2), edakit.Num(b.vgs, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$W/L = 2 I_D / [k'(V_{GS}-V_t)^2] = 2 \\times %d / (%d \\times %s^2) = %s$.", int(target), int(b.kp), edakit.Num(vov, 2), ans),
			map[string]string{"WL": ans})
	default:
		if edakit.Int(0, 1) == 0 {
			gm := b.kp * b.wl * vov
			ans := edakit.Num(gm, 1)
			q := edakit.Sprintf("In saturation, what is the transconductance $g_m = \\partial I_D/\\partial V_{GS}$ of %s, under the square law? Answer in µA/V. Round to 1 decimal place.", b.stem())
			return edakit.Numeric(q, ans,
				edakit.Sprintf("$g_m = k'(W/L)(V_{GS} - V_t) = %d \\times %d \\times %s = %s$ µA/V.", int(b.kp), int(b.wl), edakit.Num(vov, 2), ans),
				map[string]string{"gm": ans})
		}
		id := float64(10 * edakit.Int(5, 50))
		ov := 0.05 * float64(edakit.Int(2, 10))
		gm := 2 * id / ov
		ans := edakit.Num(gm, 1)
		q := edakit.Sprintf("A square-law NMOS in saturation carries %d µA at an overdrive $V_{GS} - V_t = %s$ V. What is its transconductance, in µA/V? Round to 1 decimal place.", int(id), edakit.Num(ov, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("With $I_D = \\tfrac12 k'(W/L)V_{ov}^2$, $g_m = k'(W/L)V_{ov} = 2I_D/V_{ov} = 2 \\times %d / %s = %s$ µA/V.", int(id), edakit.Num(ov, 2), ans),
			map[string]string{"gm": ans})
	}
}

// clm asks for a current with channel-length modulation (KP 1), the output resistance (KP 2)
// or the intrinsic gain (KP 3).
func clm(ctx generator.GeneratorContext) generator.Problem {
	lambda := 0.05 * float64(edakit.Int(1, 8))
	switch kpOf(ctx) {
	case 1:
		idsat := float64(10 * edakit.Int(5, 50))
		vds := 0.1 * float64(edakit.Int(5, 12))
		id := idsat * (1 + lambda*vds)
		ans := edakit.Num(id, 1)
		q := edakit.Sprintf("A saturated NMOS would carry %d µA with no channel-length modulation. With $\\lambda = %s$ V$^{-1}$ and $V_{DS} = %s$ V, what current does the model $I_D = I_{D,sat}(1 + \\lambda V_{DS})$ give, in µA? Round to 1 decimal place.",
			int(idsat), edakit.Num(lambda, 2), edakit.Num(vds, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("$I_D = %d \\times (1 + %s \\times %s) = %s$ µA.", int(idsat), edakit.Num(lambda, 2), edakit.Num(vds, 1), ans), map[string]string{"Id": ans})
	case 2:
		id := float64(10 * edakit.Int(2, 50))
		ro := 1 / (lambda * id * 1e-6) / 1000
		ans, rt := sigAnswer(ro, 3)
		q := edakit.Sprintf("A saturated NMOS carries %d µA and has $\\lambda = %s$ V$^{-1}$. What is its small-signal output resistance $r_o \\approx 1/(\\lambda I_D)$, in kilohms? %s", int(id), edakit.Num(lambda, 2), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("$r_o = 1/(%s \\times %d \\times 10^{-6}) = %s$ kilohms.", edakit.Num(lambda, 2), int(id), ans), map[string]string{"ro": ans})
	default:
		vov := 0.05 * float64(edakit.Int(2, 10))
		g := 2 / (lambda * vov)
		ans := edakit.Num(g, 1)
		q := edakit.Sprintf("A square-law NMOS is biased in saturation at overdrive $V_{GS} - V_t = %s$ V with $\\lambda = %s$ V$^{-1}$. Using $r_o = 1/(\\lambda I_D)$, what is its intrinsic gain $g_m r_o$? Round to 1 decimal place.", edakit.Num(vov, 2), edakit.Num(lambda, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$g_m = 2I_D/V_{ov}$ and $r_o = 1/(\\lambda I_D)$, so the current cancels: $g_m r_o = 2/(\\lambda V_{ov}) = 2/(%s \\times %s) = %s$.", edakit.Num(lambda, 2), edakit.Num(vov, 2), ans),
			map[string]string{"gain": ans})
	}
}

// velocitySaturation asks for the drain voltage at which the lateral field reaches the critical
// field (KP 1), the fully saturated current per width (KP 2) or an alpha-power delay ratio (KP 3).
func velocitySaturation(ctx generator.GeneratorContext) generator.Problem {
	vsat := pickF(1e7, 8e6)
	switch kpOf(ctx) {
	case 1:
		mu := float64(50 * edakit.Int(3, 10))
		L := float64(10 * edakit.Int(2, 10))
		ec := vsat / mu
		v := ec * L * 1e-7
		ans := edakit.Num(v, 3)
		q := edakit.Sprintf("A channel has effective mobility %d cm$^2$/(V s) and saturation velocity $%s$ cm/s. Taking the critical field as $E_c = v_{sat}/\\mu$, what voltage across a %d nm channel produces an average lateral field of $E_c$, in volts? Round to 3 decimal places.",
			int(mu), sciOf(vsat), int(L))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$E_c = %s / %d = %.4g$ V/cm. Over $L = %d$ nm $= %d \\times 10^{-7}$ cm, $E_c L = %s$ V. Short channels reach velocity saturation at small drain voltages.", sciOf(vsat), int(mu), ec, int(L), int(L), ans),
			map[string]string{"Ec": edakit.Sprintf("%.4g", ec)})
	case 2:
		eot := pickF(0.8, 1.0, 1.2, 1.5)
		vov := 0.05 * float64(edakit.Int(4, 12))
		i := coxF(eot) * vsat * vov * 100 // uA/um
		ans, rt := sigAnswer(i, 3)
		q := edakit.Sprintf("In the fully velocity-saturated limit, the drain current per width is $I_D/W = C_{ox} v_{sat}(V_{GS} - V_t)$. For EOT = %s nm ($\\varepsilon_{ox} = 3.453 \\times 10^{-13}$ F/cm), $v_{sat} = %s$ cm/s and overdrive %s V, what is $I_D/W$ in µA/µm? %s",
			edakit.Num(eot, 1), sciOf(vsat), edakit.Num(vov, 2), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$C_{ox} = %.4g$ F/cm$^2$. The product is %.4g A/cm, and 1 A/cm equals 100 µA/µm, giving %s µA/µm. The current is linear in overdrive, not quadratic.", coxF(eot), coxF(eot)*vsat*vov, ans),
			map[string]string{"Cox": edakit.Sprintf("%.4g", coxF(eot))})
	default:
		alpha := pickF(1.1, 1.2, 1.3, 1.4, 1.5)
		vt := 0.05 * float64(edakit.Int(5, 8))
		v1 := 0.1 * float64(edakit.Int(9, 12))
		v2 := v1 - 0.1*float64(edakit.Int(1, 3))
		r := (v2 / math.Pow(v2-vt, alpha)) / (v1 / math.Pow(v1-vt, alpha))
		ans := edakit.Num(r, 3)
		q := edakit.Sprintf("Under the alpha-power law, gate delay scales as $t_d \\propto V_{DD}/(V_{DD} - V_t)^{\\alpha}$. With $\\alpha = %s$ and $V_t = %s$ V, by what factor does delay change when $V_{DD}$ drops from %s V to %s V? Round to 3 decimal places.",
			edakit.Num(alpha, 1), edakit.Num(vt, 2), edakit.Num(v1, 1), edakit.Num(v2, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$t_d(%s)/t_d(%s) = \\frac{%s/%.2f^{%s}}{%s/%.2f^{%s}} = %s$.", edakit.Num(v2, 1), edakit.Num(v1, 1), edakit.Num(v2, 1), v2-vt, edakit.Num(alpha, 1), edakit.Num(v1, 1), v1-vt, edakit.Num(alpha, 1), ans),
			map[string]string{"ratio": ans})
	}
}

// subthreshold asks for the swing (KP 1), a subthreshold current (KP 2) or the swing from the
// capacitive divider (KP 3).
func subthreshold(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := pickF(1.0, 1.05, 1.1, 1.2, 1.3, 1.4, 1.5)
		T := float64(edakit.Pick([]int{250, 300, 350, 400}))
		s := n * kB * T * ln10 * 1000
		ans := edakit.Num(s, 1)
		q := edakit.Sprintf("A transistor has subthreshold slope factor $n = %s$ and operates at %d K. With $k/q = 8.617 \\times 10^{-5}$ V/K, what is its subthreshold swing, in mV/decade? Round to 1 decimal place.", edakit.Num(n, 2), int(T))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$S = n (kT/q)\\ln 10 = %s \\times %.5f \\times 2.3026 = %s$ mV/decade.", edakit.Num(n, 2), kB*T, ans),
			map[string]string{"S": ans})
	case 2:
		s := float64(edakit.Pick([]int{65, 70, 75, 80, 85, 90, 100}))
		i0 := float64(edakit.Pick([]int{50, 100, 200, 300, 500}))
		dv := float64(10 * edakit.Int(5, 30))
		i := i0 * 1e-9 * math.Pow(10, -dv/s)
		val, unit := currentUnit(i)
		ans, rt := sigAnswer(val, 3)
		q := edakit.Sprintf("A transistor carries %d nA at $V_{GS} = V_t$ and has a subthreshold swing of %d mV/decade. What current flows at $V_{GS} = V_t - %d$ mV, in %s? %s", int(i0), int(s), int(dv), unit, rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("%d mV below threshold is $%d/%d = %.3f$ decades, so $I = %d\\ \\text{nA} \\times 10^{-%.3f} = %s$ %s.", int(dv), int(dv), int(s), dv/s, int(i0), dv/s, ans, unit),
			map[string]string{"decades": edakit.Num(dv/s, 3)})
	default:
		cox := float64(edakit.Pick([]int{10, 15, 20, 25, 30}))
		cd := float64(edakit.Int(1, 12))
		n := 1 + cd/cox
		s := n * vT300 * ln10 * 1000
		ans := edakit.Num(s, 1)
		q := edakit.Sprintf("A bulk MOSFET has $C_{ox} = %d$ fF/µm$^2$ and a depletion capacitance of %d fF/µm$^2$, with negligible interface-trap capacitance. At 300 K ($kT/q = 25.85$ mV), what is the subthreshold swing, in mV/decade? Round to 1 decimal place.", int(cox), int(cd))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The gate controls the surface through the divider $n = 1 + C_{dep}/C_{ox} = 1 + %d/%d = %.4f$. Then $S = n \\times 25.85 \\times 2.3026 = %s$ mV/decade.", int(cd), int(cox), n, ans),
			map[string]string{"n": edakit.Num(n, 4)})
	}
}

// gateCapacitance asks for Cgs in the linear region (KP 1), Cgs in saturation (KP 2) or the
// Miller input capacitance (KP 3).
func gateCapacitance(ctx generator.GeneratorContext) generator.Problem {
	cox := float64(edakit.Pick([]int{10, 15, 20, 25, 30}))
	w := 0.1 * float64(edakit.Int(2, 20))
	l := float64(10 * edakit.Int(2, 10))
	cov := 0.05 * float64(edakit.Int(2, 8))
	wlc := w * l * 1e-3 * cox
	stem := edakit.Sprintf("An NMOS has $W = %s$ µm, $L = %d$ nm, $C_{ox} = %d$ fF/µm$^2$ and overlap capacitance %s fF/µm per edge", edakit.Num(w, 1), int(l), int(cox), edakit.Num(cov, 2))
	switch kpOf(ctx) {
	case 1:
		c := 0.5*wlc + cov*w
		ans, rt := sigAnswer(c, 3)
		q := edakit.Sprintf("%s. In strong inversion with $V_{DS}$ near zero, what is $C_{gs}$ including overlap, in fF, under the Meyer model? %s", stem, rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$WLC_{ox} = %.4g$ fF. In the linear region the channel charge splits evenly, so $C_{gs} = \\tfrac12 WLC_{ox} + C_{ov}W = %s$ fF.", wlc, ans),
			map[string]string{"WLCox": edakit.Sprintf("%.4g", wlc)})
	case 2:
		c := 2.0/3*wlc + cov*w
		ans, rt := sigAnswer(c, 3)
		q := edakit.Sprintf("%s. In saturation, what is $C_{gs}$ including overlap, in fF, under the Meyer model? %s", stem, rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$WLC_{ox} = %.4g$ fF. In saturation the channel is pinched off at the drain, so $C_{gs} = \\tfrac23 WLC_{ox} + C_{ov}W = %s$ fF, while $C_{gd}$ falls to the overlap term alone.", wlc, ans),
			map[string]string{"WLCox": edakit.Sprintf("%.4g", wlc)})
	default:
		av := float64(edakit.Int(2, 20))
		cgs := 2.0/3*wlc + cov*w
		cgd := cov * w
		c := cgs + (1+av)*cgd
		ans, rt := sigAnswer(c, 3)
		q := edakit.Sprintf("%s. It is saturated in an amplifier with voltage gain magnitude %d from gate to drain. Using Meyer capacitances, what input capacitance $C_{gs} + (1 + |A_v|)C_{gd}$ does the gate present, in fF? %s", stem, int(av), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$C_{gs} = \\tfrac23 WLC_{ox} + C_{ov}W = %.4g$ fF and $C_{gd} = C_{ov}W = %.4g$ fF. The Miller effect multiplies $C_{gd}$ by $1 + %d$, giving %s fF.", cgs, cgd, int(av), ans),
			map[string]string{"Cgs": edakit.Sprintf("%.4g", cgs), "Cgd": edakit.Sprintf("%.4g", cgd)})
	}
}

// cmosRatio computes the PMOS-to-NMOS width ratio for equal saturation current.
func cmosRatio(ctx generator.GeneratorContext) generator.Problem {
	mr := pickF(1.5, 2, 2.5, 3)
	vdd := 0.1 * float64(edakit.Int(8, 12))
	vtn := 0.05 * float64(edakit.Int(6, 9))
	vtp := 0.05 * float64(edakit.Int(6, 9))
	r := mr * math.Pow((vdd-vtn)/(vdd-vtp), 2)
	q := edakit.Sprintf("In a CMOS inverter with $V_{DD} = %s$ V, $V_{tn} = %s$ V, $|V_{tp}| = %s$ V, equal channel lengths and $\\mu_n/\\mu_p = %s$, what ratio $W_p/W_n$ gives equal saturation currents at full gate drive, under the square law?",
		edakit.Num(vdd, 1), edakit.Num(vtn, 2), edakit.Num(vtp, 2), edakit.Num(mr, 1))
	return edakit.Choice(q, fixed(r, 2),
		edakit.Sprintf("Equal currents need $\\mu_n W_n (V_{DD} - V_{tn})^2 = \\mu_p W_p (V_{DD} - |V_{tp}|)^2$, so $W_p/W_n = %s \\times (%.2f/%.2f)^2 = %.2f$.", edakit.Num(mr, 1), vdd-vtn, vdd-vtp, r),
		fixed(mr, 2), fixed(1/r, 2), fixed(math.Sqrt(r), 2), fixed(mr*(vdd-vtn)/(vdd-vtp), 2), fixed(math.Pow((vdd-vtn)/(vdd-vtp), 2), 2))
}

var cmosItems = []edakit.Item{
	{KP: 1, Q: "A static CMOS gate holds its output high. Ideally, how much current flows from supply to ground?",
		A: "None apart from leakage, because the pull-down network is off", D: []string{"A current set by the ratio of pull-up to pull-down strength", "Half the switching current", "A constant bias current set by the PMOS"},
		Why: "In every static state one network is off, so there is no DC path. This is why CMOS replaced ratioed nMOS logic."},
	{KP: 1, Q: "Why does a CMOS pull-up network use PMOS transistors rather than NMOS?",
		A: "A PMOS passes a strong 1, while an NMOS pull-up stops a threshold below the supply", D: []string{"PMOS transistors are faster than NMOS", "NMOS transistors cannot connect to the supply", "PMOS transistors have no threshold voltage"},
		Why: "An NMOS with its gate at $V_{DD}$ turns off as its source approaches $V_{DD} - V_t$. A PMOS pulls all the way to $V_{DD}$, so CMOS outputs swing rail to rail."},
	{KP: 1, Q: "A NAND gate's pull-down network is two NMOS in series. What is its pull-up network?",
		A: "Two PMOS in parallel", D: []string{"Two PMOS in series", "One PMOS", "Two NMOS in parallel"},
		Why: "The pull-up is the dual of the pull-down: series becomes parallel. The output is pulled high when either input is low."},
	{KP: 1, Q: "Bulk silicon holes have lower mobility than electrons. How do standard cells usually compensate?",
		A: "They give the PMOS a wider device or more fins than the NMOS", D: []string{"They give the PMOS a longer channel", "They raise the PMOS threshold voltage", "They connect the PMOS body to ground"},
		Why: "Drive current scales with $\\mu W/L$. Widening the PMOS balances the rise and fall delays at the cost of input capacitance and area."},
	{KP: 3, Q: "A CMOS gate switches its output from 0 to 1 and back once. How much energy is drawn from the supply to charge a load C?",
		A: "$C V_{DD}^2$: half is dissipated in the pull-up, half is stored and later dissipated in the pull-down", D: []string{"$\\tfrac12 C V_{DD}^2$, all of it returned to the supply", "$C V_{DD}$", "$2 C V_{DD}^2$"},
		Why: "The supply delivers charge $CV_{DD}$ at voltage $V_{DD}$. Half the energy heats the PMOS during charging and the stored half heats the NMOS during discharge."},
	{KP: 3, Q: "During an input transition both networks conduct briefly. Which power component does this cause?",
		A: "Short-circuit power, which grows with slow input slews", D: []string{"Leakage power, which is independent of switching", "Gate tunneling power", "Glitch power only"},
		Why: "While the input is between $V_{tn}$ and $V_{DD} - |V_{tp}|$ both devices are on. Slow slews lengthen that window, which is why timers report slew limits."},
	{KP: 3, Q: "A bulk CMOS layout omits well and substrate taps near a large I/O driver. Which failure becomes likely?",
		A: "Latch-up through the parasitic pnpn thyristor formed by the wells", D: []string{"Gate-oxide breakdown", "Electromigration in the taps", "Hot-carrier injection in the PMOS"},
		Why: "The p+/n-well/p-substrate/n+ stack forms two coupled parasitic bipolar transistors. Taps and guard rings keep the well resistances low so they cannot turn each other on."},
	{KP: 3, Q: "In a CMOS inverter, why is the switching threshold $V_M$ often set near $V_{DD}/2$?",
		A: "It balances noise margins and makes rise and fall delays similar", D: []string{"It minimizes leakage at all inputs", "It is required for the gate to work at all", "It removes short-circuit current entirely"},
		Why: "$V_M$ is where the PMOS and NMOS currents match. Sizing the PMOS to centre it gives symmetric noise margins and transitions."},
}

var shortChannelItems = []edakit.Item{
	{KP: 1, Q: "As gate length shrinks, the threshold voltage of a planar MOSFET falls even at low drain bias. What causes this roll-off?",
		A: "The source and drain depletion regions supply part of the channel depletion charge, so the gate controls less of it", D: []string{"Velocity saturation lowers the threshold", "The oxide gets thinner automatically at short lengths", "Mobility rises in short channels"},
		Why: "Charge sharing: a growing fraction of the depletion charge under the gate is imaged by the source and drain junctions, so the gate needs less voltage to invert."},
	{KP: 1, Q: "Which effect makes the threshold voltage depend on the drain voltage in a short-channel device?",
		A: "Drain-induced barrier lowering", D: []string{"Channel-length modulation", "The body effect", "Gate-induced drain leakage"},
		Why: "The drain field reaches the source barrier and lowers it, so a higher $V_{DS}$ lowers $V_t$. Channel-length modulation instead shortens the effective channel in saturation."},
	{KP: 1, Q: "Under extreme drain bias, current flows deep below the surface and the gate loses control entirely. What is this?",
		A: "Punch-through, where the source and drain depletion regions merge", D: []string{"Velocity saturation", "Gate tunneling", "Mobility degradation"},
		Why: "When the drain depletion region reaches the source, the barrier vanishes in the bulk, where the gate field is weakest."},
	{KP: 1, Q: "Which device dimension mainly sets how strongly short-channel effects appear in a planar bulk MOSFET?",
		A: "The channel length relative to an electrostatic natural length set by oxide thickness and depletion depth", D: []string{"The gate width", "The metal pitch", "The contact resistance"},
		Why: "Short-channel effects scale as roughly $e^{-L/(2\\lambda)}$, where $\\lambda$ grows with oxide thickness and the depth of the controlled silicon. Width does not enter."},
	{KP: 2, Q: "A foundry adds halo (pocket) implants near the source and drain. Which short-channel effect is it fighting?",
		A: "Threshold roll-off, by raising the local doping as the channel shortens", D: []string{"Gate leakage", "Channel-length modulation at long lengths", "Electromigration"},
		Why: "Halos raise the effective channel doping more in short devices. Overcompensation even gives the reverse short-channel effect, where $V_t$ first rises as $L$ shrinks."},
	{KP: 2, Q: "A short-channel transistor's saturation current grows nearly linearly with overdrive instead of quadratically. Why?",
		A: "Carriers reach the saturation velocity, so current is set by charge times $v_{sat}$", D: []string{"The body effect raises the threshold", "Gate tunneling adds current", "The square law applies only to PMOS"},
		Why: "When the lateral field exceeds the critical field, drift velocity stops growing. $I_D \\approx W C_{ox} v_{sat}(V_{GS} - V_t)$ is linear in overdrive."},
	{KP: 2, Q: "Why do short-channel devices show lower effective mobility at high gate overdrive?",
		A: "The strong vertical field pushes carriers against the oxide interface, where surface roughness scattering rises", D: []string{"The channel gets longer at high gate bias", "Velocity saturation lowers the vertical field", "Dopant atoms move during operation"},
		Why: "Vertical-field mobility degradation is modeled as $\\mu_{eff} = \\mu_0/(1 + \\theta(V_{GS} - V_t))$ in simple models. It is separate from the lateral-field velocity saturation."},
	{KP: 2, Q: "Which design change most directly improves gate control of the channel in a short device?",
		A: "Thinning the effective oxide or wrapping the gate around a thin body", D: []string{"Widening the transistor", "Lengthening the metal interconnect", "Raising the supply voltage"},
		Why: "Both shrink the natural length $\\lambda$. High-k dielectrics thin the EOT, and FinFETs and nanosheets thin the body and add gates."},
	{KP: 3, Q: "Why does a timing library characterized at nominal channel length still need statistical data in a short-channel process?",
		A: "Threshold depends steeply on length near the roll-off region, so small length variation becomes large delay variation", D: []string{"Short-channel effects disappear when the process is mature", "Delay is independent of threshold", "Liberty files cannot hold numbers"},
		Why: "On the steep part of the roll-off curve, a nanometre of gate-length variation shifts $V_t$ by tens of mV, which subthreshold and drive current amplify."},
	{KP: 3, Q: "A circuit designer uses longer-than-minimum channel length for a sensitive analog current mirror. What does this buy?",
		A: "Smaller short-channel effects: higher output resistance and a threshold less sensitive to length variation", D: []string{"Higher drive current per width", "Lower gate capacitance", "Immunity to temperature"},
		Why: "A longer channel moves the device away from the roll-off region and reduces DIBL and channel-length modulation, at the cost of area and speed."},
	{KP: 3, Q: "A compact model fits long-channel devices well but overestimates the off-current of minimum-length devices at high $V_{DS}$. What did the fit miss?",
		A: "It underestimates the short-channel threshold reduction, including DIBL", D: []string{"It overestimates gate leakage", "It ignores the overlap capacitance", "It uses the wrong supply voltage"},
		Why: "Off-current is exponential in $V_t$. Any error in modeled roll-off or DIBL at the short end shows up as an order-of-magnitude leakage error."},
	{KP: 3, Q: "Why do short-channel effects make an EDA timing engine care about the drain voltage of the driving transistor during a transition?",
		A: "Through DIBL and channel-length modulation, drive current depends on $V_{DS}$ throughout the transition", D: []string{"The drain voltage never changes during a transition", "Timers only use the gate voltage", "Drain voltage only affects hold time"},
		Why: "The output voltage is the drive transistor's $V_{DS}$. Current-source cell models capture $I(V_{in}, V_{out})$ for exactly this reason."},
}

var leakageItems = []edakit.Item{
	{KP: 1, Q: "In a high-performance logic block at high temperature, which leakage component usually dominates the off-state current?",
		A: "Subthreshold conduction through off transistors", D: []string{"Gate tunneling through high-k dielectrics", "Reverse-bias diode leakage of wells", "Electromigration current"},
		Why: "Subthreshold current is exponential in $-V_t/(nkT/q)$ and grows steeply with temperature, so low-threshold devices at hot corners dominate."},
	{KP: 1, Q: "Which leakage component flows through the gate dielectric rather than the channel?",
		A: "Direct tunneling gate leakage", D: []string{"Subthreshold leakage", "Gate-induced drain leakage", "Junction leakage"},
		Why: "Gate leakage is carriers tunneling through the dielectric. High-k dielectrics reduce it by allowing a physically thicker layer at the same capacitance."},
	{KP: 1, Q: "An off NMOS has its gate at 0 V and drain at $V_{DD}$. Which leakage arises from band-to-band tunneling in the gate-to-drain overlap region?",
		A: "Gate-induced drain leakage (GIDL)", D: []string{"Direct gate tunneling", "Punch-through", "Subthreshold leakage"},
		Why: "The large gate-to-drain field bends the bands in the drain overlap until electrons tunnel from valence to conduction band. The current flows from drain to body."},
	{KP: 1, Q: "Why do heavily doped halo implants raise junction leakage?",
		A: "Steep, heavily doped junctions have high fields that drive band-to-band tunneling", D: []string{"Halos thin the gate oxide", "Halos raise the subthreshold swing", "Halos add a parallel channel"},
		Why: "Band-to-band tunneling across a reverse-biased junction grows steeply with field. Halos trade lower roll-off for this junction leakage."},
	{KP: 2, Q: "Two off NMOS transistors in series leak much less than one. What is this called, and why?",
		A: "The stack effect: the intermediate node rises, giving the upper device negative $V_{GS}$, body effect and less DIBL", D: []string{"The Miller effect: capacitance doubles", "Velocity saturation: carriers slow down", "The stack effect: each device sees twice the voltage"},
		Why: "The middle node settles at a small positive voltage. The upper device then has $V_{GS} < 0$, a raised threshold and lower $V_{DS}$, all of which cut subthreshold current exponentially."},
	{KP: 2, Q: "A chip's measured leakage triples between 25 °C and 85 °C. Which leakage component is most responsible?",
		A: "Subthreshold leakage, whose threshold falls and swing rises with temperature", D: []string{"Gate tunneling, which is almost temperature independent", "Capacitor charging current", "Clock tree power"},
		Why: "Subthreshold current depends exponentially on $V_t$ and on $kT$, and $V_t$ falls by roughly 1 mV/K. Gate tunneling barely changes with temperature."},
	{KP: 2, Q: "A power-management scheme applies reverse body bias in standby. What is the trade-off?",
		A: "Subthreshold leakage falls through the body effect, but junction and GIDL leakage rise with the larger reverse bias", D: []string{"All leakage falls with no penalty", "Gate leakage rises tenfold", "Active-mode frequency rises"},
		Why: "Reverse body bias raises $V_t$, cutting subthreshold current, but increases the field across the drain-body junction. Below some point total leakage starts to rise."},
	{KP: 2, Q: "Which power-reduction technique attacks leakage by disconnecting the supply of an idle block?",
		A: "Power gating with header or footer switches", D: []string{"Clock gating", "Multi-bit flip-flop banking", "Voltage-scaled clock trees"},
		Why: "Clock gating stops dynamic power but leaves the block leaking. A high-threshold switch in series with the supply removes almost all leakage while the block is off."},
	{KP: 3, Q: "Why must a leakage power report state the process corner and temperature it was computed at?",
		A: "Leakage varies by an order of magnitude or more across corners and temperature", D: []string{"Leakage does not depend on corner", "Only dynamic power depends on temperature", "Liberty leakage tables have no conditions"},
		Why: "A fast corner has lower $V_t$, and hot silicon raises subthreshold current steeply. A single leakage number without its condition is meaningless."},
	{KP: 3, Q: "Leakage heats a die, and heat raises leakage. What failure can this loop cause?",
		A: "Thermal runaway, if leakage power grows faster with temperature than the package can remove heat", D: []string{"Latch-up, which is a purely electrical effect", "Hold violations at low temperature", "Gate oxide thickening"},
		Why: "Positive electrothermal feedback is stable only while the slope of leakage power with temperature stays below the package's heat-removal slope."},
	{KP: 3, Q: "A Liberty file gives each cell a leakage value per input state. Why per state?",
		A: "Which transistors are off, and whether they form stacks, depends on the inputs", D: []string{"Leakage depends only on cell area", "Input state changes the supply voltage", "Liberty requires it for timing arcs only"},
		Why: "A NAND2 with both inputs low has a two-high off stack and leaks much less than with one input high. State-dependent tables let power tools weight by input probabilities."},
	{KP: 3, Q: "A device engineer lowers gate leakage by switching from SiO$_2$ to HfO$_2$ at the same EOT. Why does gate leakage fall?",
		A: "The high-k layer is physically thicker for the same capacitance, and tunneling falls exponentially with thickness", D: []string{"HfO$_2$ has a larger band gap than SiO$_2$", "HfO$_2$ has zero defects", "The gate voltage is lower"},
		Why: "Capacitance depends on $\\varepsilon/t$, tunneling on $e^{-2\\kappa t}$. A higher $\\varepsilon$ permits a larger $t$. HfO$_2$ actually has a smaller gap and barrier than SiO$_2$; thickness wins."},
}
