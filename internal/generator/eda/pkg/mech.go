package pkg

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type cteMat struct {
	name  string
	alpha float64
}

var lowCTE = []cteMat{{"silicon", 2.6}, {"glass", 3.5}, {"silicon nitride", 3.2}, {"alumina", 6.9}}
var highCTE = []cteMat{{"an organic substrate", 15}, {"FR-4 board", 17}, {"copper", 16.7}, {"mold compound", 12}, {"SAC solder", 21}, {"a build-up substrate", 14}}

// thermalStress: mismatch strain, biaxial film stress, swing to a stress limit.
func thermalStress(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		lo := edakit.Pick(lowCTE)
		hi := edakit.Pick(highCTE)
		dT := float64(10 * edakit.Int(5, 25))
		ans := f((hi.alpha-lo.alpha)*dT, 0)
		q := edakit.Sprintf("%s (%s ppm/K) is bonded to %s (%s ppm/K) and cooled %.0f K from its stress-free temperature. What mismatch strain must the bond accommodate, in ppm? Round to the nearest whole number.", capFirst(lo.name), f(lo.alpha, 1), hi.name, f(hi.alpha, 1), dT)
		return num(q, ans, edakit.Sprintf("$\\Delta\\alpha\\Delta T = %s \\times %.0f = %s$ ppm.", f(hi.alpha-lo.alpha, 1), dT, ans), "strain_ppm")
	case 2:
		E := float64(10 * edakit.Int(5, 20))
		nu := step(0.2, 0.4, 0.01)
		da := step(5, 20, 0.5)
		dT := float64(10 * edakit.Int(5, 25))
		s := E * 1000 / (1 - nu) * da * 1e-6 * dT
		ans := f(s, 1)
		q := edakit.Sprintf("A thin film with $E = %.0f$ GPa and Poisson's ratio %s is bonded to a thick substrate. Their CTEs differ by %s ppm/K, and the stack cools %.0f K from stress-free. What equal-biaxial stress develops in the film, in MPa? Round to 1 decimal place.", E, f(nu, 2), f(da, 1), dT)
		return num(q, ans, edakit.Sprintf("$\\sigma = \\frac{E}{1-\\nu}\\Delta\\alpha\\Delta T = %s \\times %s = %s$ MPa.", f(E*1000/(1-nu), 1), f(da*1e-6*dT, 6), ans), "stress_mpa")
	default:
		E := float64(10 * edakit.Int(5, 20))
		nu := step(0.2, 0.4, 0.01)
		da := step(5, 20, 0.5)
		sa := float64(10 * edakit.Int(10, 60))
		dT := sa * (1 - nu) / (E * 1000 * da * 1e-6)
		ans := f(dT, 1)
		q := edakit.Sprintf("A film with $E = %.0f$ GPa and Poisson's ratio %s on a thick substrate has a CTE mismatch of %s ppm/K. Its allowable biaxial stress is %.0f MPa. What temperature swing from stress-free reaches that stress, in kelvin? Round to 1 decimal place.", E, f(nu, 2), f(da, 1), sa)
		return num(q, ans, edakit.Sprintf("$\\Delta T = \\sigma(1-\\nu)/(E\\Delta\\alpha) = %.0f \\times %s/(%.0f \\times %s) = %s$ K.", sa, f(1-nu, 2), E*1000, f(da*1e-6, 7), ans), "dt_k")
	}
}

func capFirst(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 'a' - 'A'
	}
	return string(b)
}

// timoshenko returns the curvature in 1/m of a bilayer.
func timoshenko(da, dT, t1, t2, E1, E2 float64) float64 {
	m := t1 / t2
	n := E1 / E2
	h := t1 + t2
	return 6 * da * dT * (1 + m) * (1 + m) / (h * (3*(1+m)*(1+m) + (1+m*n)*(m*m+1/(m*n))))
}

// warpage: symmetric bilayer curvature, Timoshenko's general case, deflection over a span.
func warpage(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		da := step(5, 18, 0.5)
		dT := float64(10 * edakit.Int(5, 25))
		h := step(0.3, 2, 0.1)
		k := 3 * da * 1e-6 * dT / (2 * h / 1000)
		ans := f(1/k, 2)
		q := edakit.Sprintf("Two bonded layers of equal thickness and equal modulus, %s mm thick in total, differ in CTE by %s ppm/K and cool %.0f K from stress-free. Using $\\kappa = 3\\Delta\\alpha\\Delta T/(2h)$, what is the radius of curvature, in metres? Round to 2 decimal places.", f(h, 1), f(da, 1), dT)
		return num(q, ans, edakit.Sprintf("$\\kappa = 3 \\times %s \\times %.0f/(2 \\times %s) = %s$ per m; radius $1/\\kappa = %s$ m.", f(da*1e-6, 7), dT, f(h/1000, 4), f(k, 4), ans), "radius_m")
	case 2:
		t1 := step(0.1, 0.8, 0.05)
		t2 := step(0.2, 1.2, 0.05)
		E1 := float64(edakit.Pick([]int{130, 160, 70, 120}))
		E2 := float64(edakit.Pick([]int{10, 15, 20, 25, 30}))
		da := step(8, 16, 0.2)
		dT := float64(10 * edakit.Int(5, 20))
		k := timoshenko(da*1e-6, dT, t1/1000, t2/1000, E1, E2)
		ans := f(1/k, 2)
		m := t1 / t2
		n := E1 / E2
		q := edakit.Sprintf("A die layer %s mm thick ($E_1 = %.0f$ GPa) is bonded to a layer %s mm thick ($E_2 = %.0f$ GPa). The CTEs differ by %s ppm/K and the stack cools %.0f K. Using Timoshenko's bimaterial formula $\\kappa = \\frac{6\\Delta\\alpha\\Delta T(1+m)^2}{h[3(1+m)^2 + (1+mn)(m^2 + 1/(mn))]}$ with $m = t_1/t_2$, $n = E_1/E_2$, $h = t_1 + t_2$, what is the radius of curvature, in metres? Round to 2 decimal places.", f(t1, 2), E1, f(t2, 2), E2, f(da, 1), dT)
		return num(q, ans, edakit.Sprintf("$m = %s$, $n = %s$, $h = %s$ m, $\\kappa = %s$ per m, radius %s m.", f(m, 4), f(n, 4), f((t1+t2)/1000, 5), f(k, 4), ans), "radius_m")
	default:
		R := step(0.2, 3, 0.05)
		L := float64(edakit.Int(10, 60))
		d := (L / 1000) * (L / 1000) / (8 * R) * 1e6
		ans := f(d, 1)
		q := edakit.Sprintf("A package bends to a radius of curvature of %s m. What out-of-plane warpage $\\delta = L^2/(8R)$ does it show across a %.0f mm span, in micrometres? Round to 1 decimal place.", f(R, 2), L)
		return num(q, ans, edakit.Sprintf("$\\delta = %s^2/(8 \\times %s) = %s$ m $= %s$ um.", f(L/1000, 3), f(R, 2), f(d/1e6, 7), ans), "warp_um")
	}
}

// solderFatigue: joint shear strain, Coffin-Manson scaling, DNP effect on life.
func solderFatigue(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		da := step(5, 15, 0.1)
		dT := float64(10 * edakit.Int(5, 18))
		D := step(2, 20, 0.5)
		h := float64(10 * edakit.Int(3, 15))
		g := da * 1e-6 * dT * D * 1000 / h * 100
		ans := f(g, 2)
		q := edakit.Sprintf("A solder joint of height %.0f um sits %s mm from the neutral point between parts whose CTEs differ by %s ppm/K. For a %.0f K temperature swing with no underfill, what is the joint's shear strain $\\Delta\\alpha\\Delta T \\cdot DNP/h$, as a percentage? Round to 2 decimal places.", h, f(D, 1), f(da, 1), dT)
		return num(q, ans, edakit.Sprintf("Displacement $%s \\times %.0f \\times %.0f = %s$ um over %.0f um gives %s%%.", f(da*1e-6, 7), dT, D*1000, f(da*1e-6*dT*D*1000, 3), h, ans), "gamma_pct")
	case 2:
		n := step(1.8, 2.5, 0.1)
		N1 := 100 * edakit.Int(5, 80)
		r := step(1.1, 2.0, 0.05)
		N2 := float64(N1) * math.Pow(1/r, n)
		ans := f(N2, 0)
		q := edakit.Sprintf("A joint survives %d thermal cycles at a given strain. A design change raises its cyclic shear strain by a factor of %s. With Coffin-Manson scaling $N_f \\propto \\gamma^{-n}$ and $n = %s$, how many cycles should it survive? Round to the nearest whole number.", N1, f(r, 2), f(n, 1))
		return num(q, ans, edakit.Sprintf("$N_2 = %d \\times (1/%s)^{%s} = %s$ cycles.", N1, f(r, 2), f(n, 1), ans), "cycles")
	default:
		n := step(1.8, 2.5, 0.1)
		N1 := 100 * edakit.Int(5, 80)
		D1 := float64(edakit.Int(5, 15))
		D2 := D1 + float64(edakit.Int(1, 10))
		h1 := float64(10 * edakit.Int(5, 12))
		h2 := h1
		extra := ""
		if edakit.Int(0, 1) == 1 {
			h2 = h1 + float64(10*edakit.Int(1, 4))
			extra = edakit.Sprintf(", and the standoff height rises from %.0f um to %.0f um", h1, h2)
		}
		N2 := float64(N1) * math.Pow((h2/D2)/(h1/D1), n)
		ans := f(N2, 0)
		q := edakit.Sprintf("The corner joint of a package lasts %d cycles at %.0f mm from the neutral point. A larger package puts its corner joint at %.0f mm%s. With $\\gamma \\propto DNP/h$ and $N_f \\propto \\gamma^{-n}$, $n = %s$, how many cycles does the new corner joint last? Round to the nearest whole number.", N1, D1, D2, extra, f(n, 1))
		return num(q, ans, edakit.Sprintf("Life scales as $(h/DNP)^n$: $%d \\times ((%.0f/%.0f)/(%.0f/%.0f))^{%s} = %s$ cycles.", N1, h2, D2, h1, D1, f(n, 1), ans), "cycles")
	}
}
