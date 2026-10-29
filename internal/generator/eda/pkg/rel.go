package pkg

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const boltzmannEV = 8.617e-5 // eV/K

// electromigration: current density, Black's current scaling, allowed density for a life.
func electromigration(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		I := step(0.5, 20, 0.5)
		w := step(0.1, 2, 0.05)
		t := step(0.1, 1, 0.05)
		J := I * 0.1 / (w * t)
		ans := f(J, 3)
		q := edakit.Sprintf("A copper line %s um wide and %s um thick carries %s mA. What is its current density, in MA/cm^2? Round to 3 decimal places.", f(w, 2), f(t, 2), f(I, 1))
		return num(q, ans, edakit.Sprintf("Area $%s$ um^2 $= %s \\times 10^{-8}$ cm^2. $J = %s \\times 10^{-3}/(%s \\times 10^{-8}) = %s$ MA/cm^2.", f(w*t, 4), f(w*t, 4), f(I, 1), f(w*t, 4), ans), "j_ma")
	case 2:
		M1 := step(5, 30, 0.5)
		J1 := step(0.5, 4, 0.1)
		J2 := step(0.5, 4, 0.1)
		if J2 == J1 {
			J2 = J1 + 0.5
		}
		n := edakit.Pick([]float64{1, 1.5, 2})
		ans := f(M1*math.Pow(J1/J2, n), 2)
		q := edakit.Sprintf("At fixed temperature, a power strap has a median EM lifetime of %s years at %s MA/cm^2. With Black's current exponent $n = %s$, what is its median lifetime at %s MA/cm^2, in years? Round to 2 decimal places.", f(M1, 1), f(J1, 1), f(n, 1), f(J2, 1))
		return num(q, ans, edakit.Sprintf("$MTTF_2 = %s \\times (%s/%s)^{%s} = %s$ years.", f(M1, 1), f(J1, 1), f(J2, 1), f(n, 1), ans), "mttf_years")
	default:
		M1 := step(2, 10, 0.5)
		M2 := edakit.Round(M1*step(1.5, 4, 0.25), 2)
		J1 := step(1, 5, 0.1)
		n := edakit.Pick([]float64{1, 1.5, 2})
		ans := f(J1*math.Pow(M1/M2, 1/n), 3)
		q := edakit.Sprintf("A line reaches a median EM life of %s years at %s MA/cm^2. The product needs %s years at the same temperature. With current exponent $n = %s$, what current density is allowed, in MA/cm^2? Round to 3 decimal places.", f(M1, 1), f(J1, 1), f(M2, 2), f(n, 1))
		return num(q, ans, edakit.Sprintf("$J_2 = %s \\times (%s/%s)^{1/%s} = %s$ MA/cm^2.", f(J1, 1), f(M1, 1), f(M2, 2), f(n, 1), ans), "j_allowed")
	}
}

func arrhenius(ea, tUseC, tStressC float64) float64 {
	return math.Exp(ea / boltzmannEV * (1/(tUseC+273.15) - 1/(tStressC+273.15)))
}

// emTemperature: Arrhenius factor, extrapolated use life, Blech length.
func emTemperature(ctx generator.GeneratorContext) generator.Problem {
	ea := step(0.7, 0.95, 0.05)
	tu := float64(5 * edakit.Int(17, 25))
	ts := float64(25 * edakit.Int(8, 13))
	af := arrhenius(ea, tu, ts)
	switch kpOf(ctx) {
	case 1:
		ans := f(af, 0)
		q := edakit.Sprintf("Electromigration in copper has an activation energy of %s eV. At fixed current density, by what factor is the median life at %.0f °C longer than at a %.0f °C stress? Use $k = 8.617 \\times 10^{-5}$ eV/K. Round to the nearest whole number.", f(ea, 2), tu, ts)
		return num(q, ans, edakit.Sprintf("$AF = \\exp\\left(\\frac{%s}{8.617 \\times 10^{-5}}\\left(\\frac{1}{%s} - \\frac{1}{%s}\\right)\\right) = %s$.", f(ea, 2), f(tu+273.15, 2), f(ts+273.15, 2), ans), "af")
	case 2:
		Ms := float64(10 * edakit.Int(5, 50))
		Js := step(2, 5, 0.5)
		Ju := step(0.5, 2, 0.25)
		n := edakit.Pick([]float64{1, 2})
		yrs := Ms * math.Pow(Js/Ju, n) * af / 8760
		ans := f(yrs, 1)
		q := edakit.Sprintf("An EM stress test at %s MA/cm^2 and %.0f °C gives a median life of %.0f hours. Use conditions are %s MA/cm^2 and %.0f °C. With $n = %s$ and $E_a = %s$ eV ($k = 8.617 \\times 10^{-5}$ eV/K), what is the median use life, in years of 8760 hours? Round to 1 decimal place.", f(Js, 1), ts, Ms, f(Ju, 2), tu, f(n, 0), f(ea, 2))
		return num(q, ans, edakit.Sprintf("Current factor $(%s/%s)^{%s} = %s$, temperature factor %s. $%.0f \\times %s \\times %s/8760 = %s$ years.", f(Js, 1), f(Ju, 2), f(n, 0), f(math.Pow(Js/Ju, n), 3), f(af, 1), Ms, f(math.Pow(Js/Ju, n), 3), f(af, 1), ans), "years")
	default:
		jl := float64(500 * edakit.Int(3, 10))
		j := step(0.5, 5, 0.1)
		L := jl / (j * 1e6) * 1e4
		ans := f(L, 2)
		q := edakit.Sprintf("A copper line carries %s MA/cm^2. The Blech threshold for this stack is $(jL)_c = %.0f$ A/cm. Below what segment length, in micrometres, is the line immune to electromigration? Round to 2 decimal places.", f(j, 1), jl)
		return num(q, ans, edakit.Sprintf("$L_c = %.0f/(%s \\times 10^6) = %s$ cm $= %s$ um.", jl, f(j, 1), f(L/1e4, 7), ans), "blech_um")
	}
}

// tddb: field, E-model acceleration, Weibull area scaling.
func tddb(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		V := step(0.5, 1.5, 0.05)
		t := step(1, 20, 0.5)
		ans := f(V/t*10, 3)
		q := edakit.Sprintf("A dielectric %s nm thick has %s V across it. What is the electric field, in MV/cm? Round to 3 decimal places.", f(t, 1), f(V, 2))
		return num(q, ans, edakit.Sprintf("1 V/nm is 10 MV/cm, so $E = 10 \\times %s/%s = %s$ MV/cm.", f(V, 2), f(t, 1), ans), "e_mvcm")
	case 2:
		g := step(1, 3.5, 0.5)
		Eu := step(2, 6, 0.5)
		Es := Eu + step(1, 3, 0.5)
		ans := f(math.Exp(g*(Es-Eu)), 1)
		q := edakit.Sprintf("In the E-model of dielectric breakdown, time to failure scales as $e^{-\\gamma E}$ with $\\gamma = %s$ cm/MV. Tests run at %s MV/cm; use is at %s MV/cm. By what factor is the use life longer than the test life? Round to 1 decimal place.", f(g, 1), f(Es, 1), f(Eu, 1))
		return num(q, ans, edakit.Sprintf("$e^{\\gamma(E_s - E_u)} = e^{%s \\times %s} = %s$.", f(g, 1), f(Es-Eu, 1), ans), "factor")
	default:
		eta := float64(edakit.Int(5, 50))
		ar := edakit.Pick([]float64{10, 30, 100, 300, 1000})
		beta := step(1.2, 3, 0.2)
		e2 := eta * math.Pow(1/ar, 1/beta)
		ans := f(e2, 4)
		q := edakit.Sprintf("A line-to-line dielectric test structure has a Weibull characteristic life of %.0f years with shape $\\beta = %s$. The chip contains %s times as much minimum-spaced line length as the test structure. What is the chip's characteristic life, in years? Round to 4 decimal places.", eta, f(beta, 1), f(ar, 0))
		return num(q, ans, edakit.Sprintf("Weakest link: $\\eta_2 = %.0f \\times (1/%s)^{1/%s} = %s$ years.", eta, f(ar, 0), f(beta, 1), ans), "eta_years")
	}
}

// aging: power-law shift, time to a limit, delay guard band.
func aging(ctx generator.GeneratorContext) generator.Problem {
	dv1 := step(5, 30, 0.5)
	t1 := edakit.Pick([]float64{100, 500, 1000})
	n := step(0.12, 0.25, 0.01)
	switch kpOf(ctx) {
	case 1:
		yrs := edakit.Pick([]float64{5, 7, 10, 15})
		ans := f(dv1*math.Pow(yrs*8760/t1, n), 2)
		q := edakit.Sprintf("BTI shifts a transistor's threshold by %s mV after %.0f hours of stress, following $\\Delta V_{th} \\propto t^n$ with $n = %s$. What shift is predicted after %.0f years of continuous stress (8760 hours per year), in mV? Round to 2 decimal places.", f(dv1, 1), t1, f(n, 2), yrs)
		return num(q, ans, edakit.Sprintf("$%s \\times (%s/%.0f)^{%s} = %s$ mV.", f(dv1, 1), f(yrs*8760, 0), t1, f(n, 2), ans), "dvth_mv")
	case 2:
		n = step(0.15, 0.25, 0.01)
		lim := dv1 * step(1.5, 3, 0.1)
		yrs := t1 * math.Pow(lim/dv1, 1/n) / 8760
		ans := f(yrs, 2)
		q := edakit.Sprintf("A threshold shift of %s mV is measured after %.0f hours, and it follows $t^n$ with $n = %s$. The timing guard band allows %s mV. After how many years of continuous stress is the limit reached (8760 hours per year)? Round to 2 decimal places.", f(dv1, 1), t1, f(n, 2), f(lim, 2))
		return num(q, ans, edakit.Sprintf("$t = %.0f \\times (%s/%s)^{1/%s} = %s$ hours $= %s$ years.", t1, f(lim, 2), f(dv1, 1), f(n, 2), f(yrs*8760, 0), ans), "years")
	default:
		vdd := step(0.6, 1.0, 0.05)
		vth := step(0.2, 0.4, 0.01)
		dv := step(0.01, 0.06, 0.005)
		al := step(1.1, 1.6, 0.1)
		r := math.Pow((vdd-vth)/(vdd-vth-dv), al)
		ans := f((r-1)*100, 2)
		q := edakit.Sprintf("Gate delay follows the alpha-power law, delay $\\propto 1/(V_{DD} - V_{th})^{\\alpha}$ with $\\alpha = %s$. At $V_{DD} = %s$ V and a fresh $V_{th} = %s$ V, aging adds %s mV to $V_{th}$. By what percentage does the delay grow? Round to 2 decimal places.", f(al, 1), f(vdd, 2), f(vth, 2), f(dv*1000, 1))
		return num(q, ans, edakit.Sprintf("$(%s/%s)^{%s} = %s$, so the delay grows %s%%.", f(vdd-vth, 3), f(vdd-vth-dv, 3), f(al, 1), f(r, 5), ans), "delay_pct")
	}
}

// multiphysics: two coupled iterations, the converged temperature, the runaway current.
func multiphysics(ctx generator.GeneratorContext) generator.Problem {
	beta := step(0.0035, 0.0043, 0.0001)
	for attempt := 0; ; attempt++ {
		R0 := step(0.01, 0.2, 0.005)
		I := step(0.5, 6, 0.1)
		th := float64(edakit.Int(5, 80))
		Ta := float64(edakit.Int(25, 85))
		if attempt >= generator.MaxDrawAttempts {
			R0, I, th, Ta, beta = 0.05, 3, 40, 50, 0.0039
		}
		g := th * I * I * R0 * beta
		if g > 0.6 || g < 0.01 {
			continue
		}
		base := edakit.Sprintf("A conductor has $R(T) = R_0(1 + \\beta(T - T_a))$ with $R_0 = %s$ ohm and $\\beta = %s$ per K, and $T = T_a + \\theta I^2 R(T)$ with $\\theta = %.0f$ K/W and $T_a = %.0f$ °C.", f(R0, 3), f(beta, 4), th, Ta)
		switch kpOf(ctx) {
		case 1:
			T1 := Ta + th*I*I*R0
			T2 := Ta + th*I*I*R0*(1+beta*(T1-Ta))
			ans := f(T2, 2)
			q := base + edakit.Sprintf(" It carries %s A. Starting from $T = T_a$, alternate power and temperature updates. What temperature does the second iteration give, in °C? Round to 2 decimal places.", f(I, 1))
			return num(q, ans, edakit.Sprintf("Iteration 1: $T_1 = %.0f + %.0f \\times %s \\times %s = %s$. Iteration 2: $R = %s$ ohm, $T_2 = %s$ °C.", Ta, th, f(I*I, 2), f(R0, 3), f(T1, 3), f(R0*(1+beta*(T1-Ta)), 5), ans), "t2")
		case 2:
			dT := th * I * I * R0 / (1 - g)
			ans := f(Ta+dT, 2)
			q := base + edakit.Sprintf(" It carries %s A. What is the converged steady-state temperature, in °C? Round to 2 decimal places.", f(I, 1))
			return num(q, ans, edakit.Sprintf("Loop gain $g = \\theta I^2 R_0 \\beta = %s$. $\\Delta T = %s/(1 - %s) = %s$ K, so $T = %s$ °C.", f(g, 4), f(th*I*I*R0, 3), f(g, 4), f(dT, 3), ans), "t_conv")
		default:
			Ic := 1 / math.Sqrt(th*R0*beta)
			ans := f(Ic, 2)
			q := base + " Above what current, in amperes, does the coupled system have no steady state (thermal runaway)? Round to 2 decimal places."
			return num(q, ans, edakit.Sprintf("Runaway when $\\theta I^2 R_0 \\beta = 1$: $I_c = 1/\\sqrt{%.0f \\times %s \\times %s} = %s$ A.", th, f(R0, 3), f(beta, 4), ans), "i_crit")
		}
	}
}
