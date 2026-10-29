package dev

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// pnJunction asks for a built-in potential (KP 1), its change with one doping (KP 2) or the
// built-in potential at another temperature (KP 3).
func pnJunction(ctx generator.GeneratorContext) generator.Problem {
	ma, ea := pickF(1, 2, 5), edakit.Int(15, 18)
	md, ed := pickF(1, 2, 5), edakit.Int(15, 18)
	k := kpOf(ctx)
	if k == 2 {
		ed = edakit.Int(19, 20)
	}
	na := ma * math.Pow(10, float64(ea))
	nd := md * math.Pow(10, float64(ed))
	switch k {
	case 1:
		v := vT300 * math.Log(na*nd/1e20)
		ans := edakit.Num(v, 3)
		q := edakit.Sprintf("A silicon p-n junction at 300 K has $N_A = %s$ cm$^{-3}$ on the p side and $N_D = %s$ cm$^{-3}$ on the n side. With $kT/q = 0.02585$ V and $n_i = 10^{10}$ cm$^{-3}$, what is the built-in potential, in volts? Round to 3 decimal places.",
			sci(ma, ea), sci(md, ed))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$V_{bi} = (kT/q)\\ln(N_A N_D/n_i^2) = 0.02585 \\times \\ln(%s) = 0.02585 \\times %.3f = %s$ V.", sciOf(na*nd/1e20), math.Log(na*nd/1e20), ans),
			map[string]string{"vbi": ans})
	case 2:
		f := pickF(2, 3, 4, 5, 10, 20, 100)
		dv := vT300 * math.Log(f) * 1000
		ans := edakit.Num(dv, 1)
		q := edakit.Sprintf("A one-sided $n^+$-p junction at 300 K has $N_D = %s$ cm$^{-3}$ and $N_A = %s$ cm$^{-3}$. A process change raises the p-side doping by a factor of %s. By how much does the built-in potential increase, in mV? Take $kT/q = 25.85$ mV. Round to 1 decimal place.",
			sci(md, ed), sci(ma, ea), edakit.Num(f, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$V_{bi}$ depends on $\\ln(N_A N_D)$, so multiplying $N_A$ by %s adds $(kT/q)\\ln %s = 25.85 \\times %.4f = %s$ mV, whatever the starting doping.", edakit.Num(f, 0), edakit.Num(f, 0), math.Log(f), ans),
			map[string]string{"delta_mV": ans})
	default:
		T := float64(edakit.Pick([]int{250, 350, 375, 400}))
		ni := 1e10 * math.Pow(T/300, 1.5) * math.Exp(1.12/(2*kB)*(1/300.0-1/T))
		niM, niE := mantissa(ni)
		niM = edakit.Round(niM, 2)
		niS := niM * math.Pow(10, float64(niE))
		kT := kB * T
		v := kT * math.Log(na*nd/(niS*niS))
		ans := edakit.Num(v, 3)
		q := edakit.Sprintf("A junction with $N_A = %s$ cm$^{-3}$ and $N_D = %s$ cm$^{-3}$ operates at %d K, where $kT/q = %.5f$ V and $n_i = %s$ cm$^{-3}$. What is the built-in potential, in volts? Round to 3 decimal places.",
			sci(ma, ea), sci(md, ed), int(T), kT, sci(niM, niE))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$V_{bi} = %.5f \\times \\ln(N_A N_D / n_i^2) = %.5f \\times %.3f = %s$ V. Higher temperature raises $n_i$ much faster than $kT$, so $V_{bi}$ falls as the device heats.", kT, kT, math.Log(na*nd/(niS*niS)), ans),
			map[string]string{"vbi": ans})
	}
}

// depletion asks for a depletion width (KP 1), a junction capacitance per area (KP 2) or a
// capacitance ratio across reverse bias (KP 3).
func depletion(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		ma, ea := pickF(1, 2, 5), edakit.Int(15, 17)
		md, ed := pickF(1, 2, 5), edakit.Int(15, 17)
		na := ma * math.Pow(10, float64(ea))
		nd := md * math.Pow(10, float64(ed))
		vbi := edakit.Round(vT300*math.Log(na*nd/1e20), 2)
		vr := pickF(0, 0, 0.5, 1, 2)
		W := math.Sqrt(2*epsSi*(vbi+vr)/qC*(1/na+1/nd)) * 1e7
		ans, rt := sigAnswer(W, 3)
		bias := "at zero bias"
		if vr > 0 {
			bias = edakit.Sprintf("under a reverse bias of %s V", edakit.Num(vr, 1))
		}
		q := edakit.Sprintf("A silicon junction has $N_A = %s$ cm$^{-3}$, $N_D = %s$ cm$^{-3}$ and $V_{bi} = %.2f$ V. Using $\\varepsilon_{Si} = 1.036 \\times 10^{-12}$ F/cm and $q = 1.602 \\times 10^{-19}$ C, what is the total depletion width %s, in nm? %s",
			sci(ma, ea), sci(md, ed), vbi, bias, rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$W = \\sqrt{\\frac{2\\varepsilon_{Si}(V_{bi} + V_R)}{q}\\left(\\frac{1}{N_A} + \\frac{1}{N_D}\\right)}$ with $V_{bi} + V_R = %s$ V gives $W = %.4g$ cm $= %s$ nm.", edakit.Num(vbi+vr, 2), W*1e-7, ans),
			map[string]string{"W_nm": ans})
	case 2:
		ma, ea := pickF(1, 2, 5), edakit.Int(16, 18)
		na := ma * math.Pow(10, float64(ea))
		vbi := pickF(0.85, 0.9, 0.95, 1.0)
		vr := pickF(0, 0.5, 1, 1.5, 2, 3)
		W := math.Sqrt(2 * epsSi * (vbi + vr) / (qC * na))
		c := epsSi / W * 1e7 // fF/um^2
		ans, rt := sigAnswer(c, 3)
		q := edakit.Sprintf("An $n^+$-p source junction has $N_A = %s$ cm$^{-3}$ on the lightly doped side and $V_{bi} = %.2f$ V. Treat it as one-sided. With $\\varepsilon_{Si} = 1.036 \\times 10^{-12}$ F/cm and $q = 1.602 \\times 10^{-19}$ C, what is the junction capacitance per unit area at a reverse bias of %s V, in fF/µm$^2$? %s",
			sci(ma, ea), vbi, edakit.Num(vr, 1), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$W = \\sqrt{2\\varepsilon_{Si}(V_{bi}+V_R)/(qN_A)} = %.4g$ cm. Then $C_j = \\varepsilon_{Si}/W = %.4g$ F/cm$^2$, and 1 F/cm$^2$ equals $10^7$ fF/µm$^2$, so $C_j = %s$ fF/µm$^2$.", W, epsSi/W, ans),
			map[string]string{"W_cm": edakit.Sprintf("%.4g", W)})
	default:
		vbi := pickF(0.7, 0.8, 0.9, 1.0)
		v1 := pickF(0, 0.5, 1)
		v2 := v1 + pickF(1, 2, 3, 4, 5)
		r := math.Sqrt((vbi + v2) / (vbi + v1))
		ans := edakit.Num(r, 3)
		q := edakit.Sprintf("An abrupt junction has $V_{bi} = %.1f$ V. A parasitic extractor tabulates its capacitance at reverse biases of %s V and %s V. What is the ratio $C_j(%s\\ \\text{V}) / C_j(%s\\ \\text{V})$? Round to 3 decimal places.",
			vbi, edakit.Num(v1, 1), edakit.Num(v2, 1), edakit.Num(v1, 1), edakit.Num(v2, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("For an abrupt junction $C_j \\propto 1/W \\propto (V_{bi} + V_R)^{-1/2}$, so the ratio is $\\sqrt{(%.1f + %s)/(%.1f + %s)} = %s$.", vbi, edakit.Num(v2, 1), vbi, edakit.Num(v1, 1), ans),
			map[string]string{"ratio": ans})
	}
}

// diode asks for a forward current (KP 1), the voltage for a target current (KP 2) or the
// voltage step per current ratio (KP 3).
func diode(ctx generator.GeneratorContext) generator.Problem {
	isM, isE := pickF(1, 2, 5), -edakit.Int(12, 15)
	is := isM * math.Pow(10, float64(isE))
	n := pickF(1.0, 1.0, 1.1, 1.2, 1.5)
	switch kpOf(ctx) {
	case 1:
		v := 0.05 * float64(edakit.Int(7, 14))
		i := is * (math.Exp(v/(n*vT300)) - 1)
		val, unit := currentUnit(i)
		ans, rt := sigAnswer(val, 3)
		q := edakit.Sprintf("A diode has saturation current $I_S = %s$ A and ideality factor $n = %s$. At 300 K ($kT/q = 25.85$ mV), what current flows at a forward bias of %s V, in %s? %s",
			sci(isM, isE), edakit.Num(n, 1), edakit.Num(v, 2), unit, rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$I = I_S(e^{V/(nkT/q)} - 1)$ with $V/(n kT/q) = %.3f$, so $I = %.4g$ A $= %s$ %s.", v/(n*vT300), i, ans, unit),
			map[string]string{"I_A": edakit.Sprintf("%.4g", i)})
	case 2:
		iM, iE := pickF(1, 2, 5), -edakit.Int(3, 6)
		i := iM * math.Pow(10, float64(iE))
		v := n * vT300 * math.Log(i/is+1)
		ans := edakit.Num(v, 3)
		q := edakit.Sprintf("A diode at 300 K has $I_S = %s$ A and $n = %s$. What forward voltage gives a current of $%s$ A? Take $kT/q = 0.02585$ V. Round to 3 decimal places, in volts.",
			sci(isM, isE), edakit.Num(n, 1), sci(iM, iE))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$V = n (kT/q) \\ln(I/I_S + 1) = %s \\times 0.02585 \\times %.3f = %s$ V.", edakit.Num(n, 1), math.Log(i/is+1), ans),
			map[string]string{"V": ans})
	default:
		r := pickF(2, 3, 5, 10, 100, 1000)
		dv := n * vT300 * math.Log(r) * 1000
		ans := edakit.Num(dv, 1)
		q := edakit.Sprintf("A forward-biased diode with ideality factor $n = %s$ carries currents far above $I_S$ at 300 K ($kT/q = 25.85$ mV). By how much must the voltage rise to multiply the current by %s, in mV? Round to 1 decimal place.",
			edakit.Num(n, 1), edakit.Num(r, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Well above $I_S$, $I \\approx I_S e^{V/(nkT/q)}$, so $\\Delta V = n (kT/q)\\ln %s = %s \\times 25.85 \\times %.4f = %s$ mV. For $n = 1$ a factor of 10 costs 59.5 mV.", edakit.Num(r, 0), edakit.Num(n, 1), math.Log(r), ans),
			map[string]string{"dV_mV": ans})
	}
}

// bjt asks for a collector current (KP 1), a base or emitter current from beta (KP 2) or a
// transconductance (KP 3).
func bjt(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		isM, isE := pickF(1, 2, 5), -edakit.Int(14, 16)
		is := isM * math.Pow(10, float64(isE))
		vbe := 0.01 * float64(edakit.Int(55, 75))
		ic := is * math.Exp(vbe/vT300)
		val, unit := currentUnit(ic)
		ans, rt := sigAnswer(val, 3)
		q := edakit.Sprintf("An npn transistor in forward-active mode has $I_S = %s$ A. Ignoring the Early effect, what collector current flows at $V_{BE} = %s$ V and 300 K ($kT/q = 25.85$ mV), in %s? %s",
			sci(isM, isE), edakit.Num(vbe, 2), unit, rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$I_C = I_S e^{V_{BE}/(kT/q)} = %s \\times e^{%.3f} = %.4g$ A $= %s$ %s.", sci(isM, isE), vbe/vT300, ic, ans, unit),
			map[string]string{"Ic_A": edakit.Sprintf("%.4g", ic)})
	case 2:
		ic := pickF(0.1, 0.2, 0.5, 1, 2, 5)
		beta := float64(edakit.Pick([]int{50, 80, 100, 150, 200, 250}))
		ib := ic / beta * 1000
		if edakit.Int(0, 1) == 0 {
			ans, rt := sigAnswer(ib, 3)
			q := edakit.Sprintf("A bipolar transistor carries $I_C = %s$ mA with current gain $\\beta = %d$. What base current does it draw, in µA? %s", edakit.Num(ic, 1), int(beta), rt)
			return edakit.Numeric(q, ans, edakit.Sprintf("$I_B = I_C/\\beta = %s\\ \\text{mA} / %d = %s$ µA.", edakit.Num(ic, 1), int(beta), ans), map[string]string{"Ib_uA": ans})
		}
		ie := ic * (1 + 1/beta)
		ans := edakit.Num(ie, 4)
		q := edakit.Sprintf("A bipolar transistor carries $I_C = %s$ mA with current gain $\\beta = %d$. What is the emitter current, in mA? Round to 4 decimal places.", edakit.Num(ic, 1), int(beta))
		return edakit.Numeric(q, ans, edakit.Sprintf("$I_E = I_C + I_B = I_C(1 + 1/\\beta) = %s$ mA.", ans), map[string]string{"Ib_uA": edakit.Num(ib, 3)})
	default:
		ic := pickF(0.05, 0.1, 0.2, 0.5, 1, 2, 5)
		T := float64(edakit.Pick([]int{300, 300, 350, 400}))
		vt := kB * T
		gm := ic / vt // mA/V = mS
		ans, rt := sigAnswer(gm, 3)
		q := edakit.Sprintf("A bipolar transistor is biased at $I_C = %s$ mA and $T = %d$ K ($kT/q = %.5f$ V). What is its transconductance $g_m$, in mS? %s", edakit.Num(ic, 2), int(T), vt, rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$g_m = dI_C/dV_{BE} = I_C/(kT/q) = %s\\ \\text{mA} / %.5f\\ \\text{V} = %s$ mS. It depends only on the bias current and temperature, not on device size.", edakit.Num(ic, 2), vt, ans),
			map[string]string{"gm_mS": ans})
	}
}

// tunneling asks for the WKB decay constant (KP 1), a transmission probability (KP 2) or the
// growth of tunneling current when a barrier thins (KP 3).
func tunneling(ctx generator.GeneratorContext) generator.Problem {
	du := pickF(0.5, 1.0, 1.5, 2.0, 3.1, 3.2)
	mr := pickF(0.3, 0.4, 0.5, 0.6, 1.0)
	kappa := math.Sqrt(mr * du / 0.0381)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(kappa, 2)
		q := edakit.Sprintf("An electron meets a rectangular barrier %s eV above its energy. Its tunneling effective mass is $%s\\, m_0$. Using $\\hbar^2/(2m_0) = 0.0381$ eV nm$^2$, what is the decay constant $\\kappa = \\sqrt{2m^*(U - E)}/\\hbar$ inside the barrier, in nm$^{-1}$? Round to 2 decimal places.",
			edakit.Num(du, 1), edakit.Num(mr, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\kappa = \\sqrt{(m^*/m_0)(U-E)/0.0381} = \\sqrt{%s \\times %s / 0.0381} = %s$ nm$^{-1}$.", edakit.Num(mr, 1), edakit.Num(du, 1), ans),
			map[string]string{"kappa": ans})
	case 2:
		k := edakit.Round(kappa, 2)
		d := pickF(0.8, 1.0, 1.2, 1.5, 2.0, 2.5, 3.0)
		lt := -2 * k * d / ln10
		ans := edakit.Num(lt, 2)
		q := edakit.Sprintf("Inside a barrier the wavefunction decays with $\\kappa = %s$ nm$^{-1}$. In the WKB approximation the transmission is $T \\approx e^{-2\\kappa d}$. For a barrier $d = %s$ nm thick, what is $\\log_{10} T$? Round to 2 decimal places.",
			edakit.Num(k, 2), edakit.Num(d, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$2\\kappa d = 2 \\times %s \\times %s = %.3f$, and $\\log_{10} T = -2\\kappa d/\\ln 10 = %s$.", edakit.Num(k, 2), edakit.Num(d, 1), 2*k*d, ans),
			map[string]string{"two_kappa_d": edakit.Sprintf("%.3f", 2*k*d)})
	default:
		k := edakit.Round(kappa, 2)
		dd := pickF(0.1, 0.2, 0.3, 0.4, 0.5)
		f := math.Exp(2 * k * dd)
		ans, rt := sigAnswer(f, 3)
		q := edakit.Sprintf("A gate dielectric has a tunneling decay constant $\\kappa = %s$ nm$^{-1}$, and the current scales as $e^{-2\\kappa d}$. By what factor does the tunneling current grow when the dielectric is thinned by %s nm? %s",
			edakit.Num(k, 2), edakit.Num(dd, 1), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The prefactor cancels in the ratio: $e^{2\\kappa \\Delta d} = e^{2 \\times %s \\times %s} = e^{%.3f} = %s$.", edakit.Num(k, 2), edakit.Num(dd, 1), 2*k*dd, ans),
			map[string]string{"exponent": edakit.Sprintf("%.3f", 2*k*dd)})
	}
}
