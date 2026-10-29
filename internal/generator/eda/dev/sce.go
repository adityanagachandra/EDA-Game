package dev

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// dibl asks for a saturation threshold (KP 1), a DIBL coefficient from two thresholds (KP 2)
// or the off-current increase DIBL causes (KP 3).
func dibl(ctx generator.GeneratorContext) generator.Problem {
	eta := float64(10 * edakit.Int(3, 12)) // mV/V
	vdd := 0.05 * float64(edakit.Int(14, 20))
	switch kpOf(ctx) {
	case 1:
		vtl := 0.01 * float64(edakit.Int(30, 48))
		vts := vtl - eta*(vdd-0.05)/1000
		ans := edakit.Num(vts, 3)
		q := edakit.Sprintf("A transistor has $V_{t,lin} = %s$ V measured at $V_{DS} = 0.05$ V and a DIBL coefficient of %d mV/V. What is its threshold at $V_{DS} = %s$ V, in volts? Round to 3 decimal places.", edakit.Num(vtl, 2), int(eta), edakit.Num(vdd, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$V_t$ drops by $\\eta \\Delta V_{DS} = %d \\times (%s - 0.05) = %.1f$ mV, so $V_{t,sat} = %s$ V.", int(eta), edakit.Num(vdd, 2), eta*(vdd-0.05), ans),
			map[string]string{"shift_mV": edakit.Num(eta*(vdd-0.05), 1)})
	case 2:
		vtl := 0.01 * float64(edakit.Int(30, 48))
		drop := float64(5 * edakit.Int(4, 24))
		vts := vtl - drop/1000
		e := drop / (vdd - 0.05)
		ans := edakit.Num(e, 1)
		q := edakit.Sprintf("A compact-model extraction measures $V_t = %s$ V at $V_{DS} = 0.05$ V and $V_t = %s$ V at $V_{DS} = %s$ V. What is the DIBL coefficient, in mV/V? Round to 1 decimal place.", edakit.Num(vtl, 2), edakit.Num(vts, 3), edakit.Num(vdd, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\eta = -\\Delta V_t/\\Delta V_{DS} = %d\\ \\text{mV} / %s\\ \\text{V} = %s$ mV/V.", int(drop), edakit.Num(vdd-0.05, 2), ans),
			map[string]string{"eta": ans})
	default:
		s := float64(edakit.Pick([]int{65, 70, 75, 80, 90}))
		f := math.Pow(10, eta*(vdd-0.05)/s)
		ans, rt := sigAnswer(f, 3)
		q := edakit.Sprintf("A transistor has a DIBL coefficient of %d mV/V and a subthreshold swing of %d mV/decade. By what factor is its off-current at $V_{DS} = %s$ V larger than it would be if the threshold stayed at its $V_{DS} = 0.05$ V value? %s", int(eta), int(s), edakit.Num(vdd, 2), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("DIBL lowers $V_t$ by $%d \\times %s = %.1f$ mV, which is $%.1f/%d = %.3f$ decades of subthreshold current: a factor of %s.", int(eta), edakit.Num(vdd-0.05, 2), eta*(vdd-0.05), eta*(vdd-0.05), int(s), eta*(vdd-0.05)/s, ans),
			map[string]string{"decades": edakit.Num(eta*(vdd-0.05)/s, 3)})
	}
}

// vtRolloff asks for a roll-off shift from the exponential model (KP 1), a natural length
// (KP 2) or the shortest length that keeps roll-off below a limit (KP 3).
func vtRolloff(ctx generator.GeneratorContext) generator.Problem {
	A := pickF(0.5, 0.8, 1.0, 1.5, 2.0)
	lam := float64(edakit.Int(4, 15))
	switch kpOf(ctx) {
	case 1:
		L := float64(5 * edakit.Int(4, 20))
		d := A * math.Exp(-L/(2*lam)) * 1000
		ans := edakit.Num(d, 1)
		q := edakit.Sprintf("A roll-off model gives the short-channel threshold reduction as $\\Delta V_t = A\\, e^{-L/(2\\lambda)}$ with $A = %s$ V and natural length $\\lambda = %d$ nm. What is $\\Delta V_t$ at $L = %d$ nm, in mV? Round to 1 decimal place.", edakit.Num(A, 1), int(lam), int(L))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$L/(2\\lambda) = %d/%d = %.4f$, so $\\Delta V_t = %s \\times e^{-%.4f} = %s$ mV.", int(L), int(2*lam), L/(2*lam), edakit.Num(A, 1), L/(2*lam), ans),
			map[string]string{"exponent": edakit.Num(L/(2*lam), 4)})
	case 2:
		tsi := float64(edakit.Int(5, 20))
		tox := pickF(0.8, 1.0, 1.2, 1.5, 2.0)
		if edakit.Int(0, 1) == 0 {
			l := math.Sqrt(3 * tsi * tox)
			ans := edakit.Num(l, 2)
			q := edakit.Sprintf("For a single-gate fully depleted SOI device, one standard estimate of the natural length is $\\lambda = \\sqrt{(\\varepsilon_{Si}/\\varepsilon_{ox})\\, t_{Si}\\, t_{ox}}$. With $\\varepsilon_{Si}/\\varepsilon_{ox} = 3$, $t_{Si} = %d$ nm and $t_{ox} = %s$ nm, what is $\\lambda$, in nm? Round to 2 decimal places.", int(tsi), edakit.Num(tox, 1))
			return edakit.Numeric(q, ans, edakit.Sprintf("$\\lambda = \\sqrt{3 \\times %d \\times %s} = %s$ nm.", int(tsi), edakit.Num(tox, 1), ans), map[string]string{"lambda": ans})
		}
		l := math.Sqrt(3 * tsi * tox / 2)
		ans := edakit.Num(l, 2)
		q := edakit.Sprintf("For a symmetric double-gate device, one standard estimate of the natural length is $\\lambda = \\sqrt{\\varepsilon_{Si} t_{Si} t_{ox}/(2\\varepsilon_{ox})}$. With $\\varepsilon_{Si}/\\varepsilon_{ox} = 3$, $t_{Si} = %d$ nm and $t_{ox} = %s$ nm, what is $\\lambda$, in nm? Round to 2 decimal places.", int(tsi), edakit.Num(tox, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\lambda = \\sqrt{3 \\times %d \\times %s / 2} = %s$ nm. The second gate shrinks $\\lambda$ by $\\sqrt 2$.", int(tsi), edakit.Num(tox, 1), ans), map[string]string{"lambda": ans})
	default:
		lim := float64(5 * edakit.Int(2, 10))
		L := 2 * lam * math.Log(A*1000/lim)
		ans := edakit.Num(L, 1)
		q := edakit.Sprintf("With the roll-off model $\\Delta V_t = A\\, e^{-L/(2\\lambda)}$, $A = %s$ V and $\\lambda = %d$ nm, what is the shortest gate length that keeps $\\Delta V_t$ at or below %d mV, in nm? Round to 1 decimal place.", edakit.Num(A, 1), int(lam), int(lim))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Solve $A e^{-L/(2\\lambda)} = \\Delta V_{max}$: $L = 2\\lambda \\ln(A/\\Delta V_{max}) = %d \\times \\ln(%s/%d) = %s$ nm.", int(2*lam), edakit.Num(A*1000, 0), int(lim), ans),
			map[string]string{"L_min": ans})
	}
}

// vtFlavors asks for the leakage ratio between two threshold flavors (KP 1), the threshold
// step for a target ratio (KP 2) or the total leakage of a mixed-flavor block (KP 3).
func vtFlavors(ctx generator.GeneratorContext) generator.Problem {
	s := float64(edakit.Pick([]int{65, 70, 75, 80, 85}))
	switch kpOf(ctx) {
	case 1:
		dv := float64(10 * edakit.Int(3, 15))
		names := edakit.Pick([][2]string{{"LVT", "SVT"}, {"SVT", "HVT"}, {"ULVT", "LVT"}, {"LVT", "HVT"}})
		r := math.Pow(10, dv/s)
		ans, rt := sigAnswer(r, 3)
		q := edakit.Sprintf("A library's %s cells have a threshold %d mV below its %s cells. With a subthreshold swing of %d mV/decade and everything else equal, how many times more off-current does an %s cell leak than the matching %s cell? %s", names[0], int(dv), names[1], int(s), names[0], names[1], rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("%d mV is $%d/%d = %.3f$ decades, so the ratio is $10^{%.3f} = %s$.", int(dv), int(dv), int(s), dv/s, dv/s, ans),
			map[string]string{"decades": edakit.Num(dv/s, 3)})
	case 2:
		r := pickF(3, 5, 10, 20, 30, 50, 100)
		dv := s * math.Log10(r)
		ans := edakit.Num(dv, 1)
		q := edakit.Sprintf("A process team wants a new high-threshold flavor that leaks %s times less than the standard flavor. With a subthreshold swing of %d mV/decade, how much higher must its threshold be, in mV? Round to 1 decimal place.", edakit.Num(r, 0), int(s))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\Delta V_t = S \\log_{10} R = %d \\times %.4f = %s$ mV.", int(s), math.Log10(r), ans), map[string]string{"dVt": ans})
	default:
		base := pickF(0.5, 1, 2)
		d1 := float64(10 * edakit.Int(4, 10))
		d2 := float64(10 * edakit.Int(4, 10))
		nl := 1000 * edakit.Int(1, 20)
		ns := 1000 * edakit.Int(10, 80)
		nh := 1000 * edakit.Int(5, 60)
		il := base * math.Pow(10, d1/s)
		ih := base / math.Pow(10, d2/s)
		tot := (float64(nl)*il + float64(ns)*base + float64(nh)*ih) / 1000 // uA
		ans, rt := sigAnswer(tot, 3)
		q := edakit.Sprintf("A block has %d LVT, %d SVT and %d HVT cells. An SVT cell leaks %s nA. LVT is %d mV below SVT and HVT is %d mV above, with a subthreshold swing of %d mV/decade. Treating all leakage as subthreshold, what is the block's total leakage current, in µA? %s",
			nl, ns, nh, edakit.Num(base, 1), int(d1), int(d2), int(s), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Per cell: LVT leaks $%s \\times 10^{%d/%d} = %.4g$ nA and HVT $%s / 10^{%d/%d} = %.4g$ nA. The sum is $%d \\times %.4g + %d \\times %s + %d \\times %.4g$ nA $= %s$ µA. A few thousand LVT cells can dominate.",
				edakit.Num(base, 1), int(d1), int(s), il, edakit.Num(base, 1), int(d2), int(s), ih, nl, il, ns, edakit.Num(base, 1), nh, ih, ans),
			map[string]string{"I_lvt_nA": edakit.Sprintf("%.4g", il), "I_hvt_nA": edakit.Sprintf("%.4g", ih)})
	}
}

// gateLeakage asks for a high-k physical thickness (KP 1), the decades of tunneling reduction
// it buys (KP 2) or the total gate leakage of a chip (KP 3).
func gateLeakage(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		eot := pickF(0.7, 0.8, 0.9, 1.0, 1.1, 1.2)
		k := pickF(16, 20, 22, 25)
		t := eot * k / 3.9
		ans := edakit.Num(t, 2)
		q := edakit.Sprintf("A gate stack must reach EOT = %s nm using a high-k layer with relative permittivity %s and no interfacial oxide. How thick must the high-k layer be, in nm? Round to 2 decimal places.", edakit.Num(eot, 1), edakit.Num(k, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$t_{hk} = \\text{EOT} \\times k/3.9 = %s \\times %s/3.9 = %s$ nm.", edakit.Num(eot, 1), edakit.Num(k, 0), ans), map[string]string{"t_hk": ans})
	case 2:
		eot := pickF(0.8, 1.0, 1.2)
		k := pickF(16, 20, 25)
		kox := pickF(6.0, 6.4)
		khk := pickF(2.5, 2.8, 3.0, 3.3)
		thk := eot * k / 3.9
		dec := (2*khk*thk - 2*kox*eot) / ln10
		ans := edakit.Num(dec, 1)
		q := edakit.Sprintf("Model direct tunneling as $J \\propto e^{-2\\kappa t}$ with equal prefactors. An SiO$_2$ gate of %s nm has $\\kappa = %s$ nm$^{-1}$. A high-k layer ($k = %s$) of the same EOT has $\\kappa = %s$ nm$^{-1}$. By how many decades (powers of ten) does the high-k layer reduce gate leakage? Round to 1 decimal place.",
			edakit.Num(eot, 1), edakit.Num(kox, 1), edakit.Num(k, 0), edakit.Num(khk, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The high-k layer is $%s \\times %s/3.9 = %.3f$ nm thick. Exponents: $2 \\times %s \\times %.3f = %.3f$ versus $2 \\times %s \\times %s = %.3f$. The difference over $\\ln 10$ is %s decades. Its lower barrier gives a smaller $\\kappa$, but the extra thickness wins.",
				edakit.Num(eot, 1), edakit.Num(k, 0), thk, edakit.Num(khk, 1), thk, 2*khk*thk, edakit.Num(kox, 1), edakit.Num(eot, 1), 2*kox*eot, ans),
			map[string]string{"t_hk": edakit.Num(thk, 3)})
	default:
		jM, jE := pickF(1, 2, 5), -edakit.Int(0, 2)
		j := jM * math.Pow(10, float64(jE)) // A/cm^2
		nDev := pickF(1e8, 2e8, 5e8, 1e9)
		w := pickF(0.05, 0.1, 0.2)
		l := pickF(0.015, 0.02, 0.03)
		area := nDev * w * l * 1e-8 // cm^2
		i := j * area * 1000        // mA
		ans, rt := sigAnswer(i, 3)
		q := edakit.Sprintf("A chip has $%s$ transistors whose gates average $W = %s$ µm by $L = %s$ µm. The gate dielectric leaks $%s$ A/cm$^2$ at the supply voltage. Assuming every gate is biased, what is the total gate leakage current, in mA? %s",
			sciOf(nDev), edakit.Num(w, 2), edakit.Num(l, 3), sci(jM, jE), rt)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Total gate area is $%s \\times %s \\times %s$ µm$^2 = %.4g$ cm$^2$ (1 µm$^2 = 10^{-8}$ cm$^2$). Then $I = J \\times A = %.4g$ A $= %s$ mA.", sciOf(nDev), edakit.Num(w, 2), edakit.Num(l, 3), area, j*area, ans),
			map[string]string{"area_cm2": edakit.Sprintf("%.4g", area)})
	}
}
