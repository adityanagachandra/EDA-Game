package em

import (
	"math"
	"math/cmplx"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerAC(reg *generator.Registry) {
	reg.Register("em.ac.phasors", edakit.Func(phasors))
	reg.Register("em.ac.impedance", edakit.Func(impedance))
	reg.Register("em.ac.resonance", edakit.Func(resonance))
}

var freqsMHz = []float64{10, 20, 50, 100, 200, 250, 500, 800, 1000, 2000}

func phasors(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a1, a2 := float64(edakit.Int(1, 10)), float64(edakit.Int(1, 10))
		p1, p2 := 15*edakit.Int(-6, 6), 15*edakit.Int(-6, 6)
		if p1 == p2 {
			p2 += 45
		}
		z := cmplx.Rect(a1, float64(p1)*math.Pi/180) + cmplx.Rect(a2, float64(p2)*math.Pi/180)
		return edakit.Numeric(
			sp("Two currents of the same frequency meet at a node: $%s\\cos(\\omega t %s)$ mA and $%s\\cos(\\omega t %s)$ mA. What is the amplitude of their sum, in mA? Round to 3 decimal places.", num(a1, 0), phaseText(p1), num(a2, 0), phaseText(p2)),
			num(cmplx.Abs(z), 3),
			sp("Add the phasors: $%s\\angle %d^\\circ + %s\\angle %d^\\circ = %s %s j%s$, whose magnitude is %s mA. Amplitudes add directly only when the phases match.", num(a1, 0), p1, num(a2, 0), p2, num(real(z), 3), plusMinus(imag(z)), num(math.Abs(imag(z)), 3), num(cmplx.Abs(z), 3)),
			facts("re", num(real(z), 3), "im", num(imag(z), 3)))
	case 2:
		f := edakit.Pick(freqsMHz)
		if edakit.Int(0, 1) == 0 {
			c := float64(edakit.Int(1, 50))
			x := 1 / (2 * math.Pi * f * 1e6 * c * 1e-12)
			return edakit.Numeric(
				sp("What is the magnitude of the impedance of a %s pF capacitor at %s MHz, in ohms? Round to 2 decimal places.", num(c, 0), num(f, 0)),
				num(x, 2),
				sp("$|Z_C| = 1/(\\omega C) = 1/(2\\pi \\cdot %s \\times 10^6 \\cdot %s \\times 10^{-12}) = %s$ ohm.", num(f, 0), num(c, 0), num(x, 2)),
				facts("x_ohm", num(x, 2)))
		}
		l := float64(edakit.Int(1, 50))
		x := 2 * math.Pi * f * 1e6 * l * 1e-9
		return edakit.Numeric(
			sp("What is the magnitude of the impedance of a %s nH inductor at %s MHz, in ohms? Round to 2 decimal places.", num(l, 0), num(f, 0)),
			num(x, 2),
			sp("$|Z_L| = \\omega L = 2\\pi \\cdot %s \\times 10^6 \\cdot %s \\times 10^{-9} = %s$ ohm.", num(f, 0), num(l, 0), num(x, 2)),
			facts("x_ohm", num(x, 2)))
	default:
		f := edakit.Pick(freqsMHz)
		r := float64(10 * edakit.Int(1, 20))
		c := float64(edakit.Int(1, 50))
		w := 2 * math.Pi * f * 1e6
		ph := deg(math.Atan(1 / (w * r * c * 1e-12)))
		return edakit.Numeric(
			sp("A sinusoidal voltage at %s MHz drives a %s ohm resistor in series with a %s pF capacitor. By how many degrees does the current lead the source voltage? Round to 1 decimal place.", num(f, 0), num(r, 0), num(c, 0)),
			num(ph, 1),
			sp("$Z = R - j/(\\omega C)$ with $1/(\\omega C) = %s$ ohm. The current's phase is minus the impedance angle: $\\arctan(1/(\\omega R C)) = %s$ degrees.", num(1/(w*c*1e-12), 2), num(ph, 1)),
			facts("xc", num(1/(w*c*1e-12), 2)))
	}
}

func phaseText(p int) string {
	switch {
	case p > 0:
		return sp("+ %d^\\circ", p)
	case p < 0:
		return sp("- %d^\\circ", -p)
	}
	return ""
}

func impedance(ctx generator.GeneratorContext) generator.Problem {
	f := edakit.Pick(freqsMHz)
	w := 2 * math.Pi * f * 1e6
	r := float64(5 * edakit.Int(1, 40))
	l := float64(edakit.Int(1, 100))
	c := float64(edakit.Int(1, 100))
	xl, xc := w*l*1e-9, 1/(w*c*1e-12)
	x := xl - xc
	stem := sp("A series circuit of $R = %s$ ohm, $L = %s$ nH and $C = %s$ pF is driven at %s MHz.", num(r, 0), num(l, 0), num(c, 0), num(f, 0))
	switch kp(ctx) {
	case 1:
		z := math.Hypot(r, x)
		return edakit.Numeric(stem+" What is the magnitude of its impedance, in ohms? Round to 2 decimal places.",
			num(z, 2),
			sp("$X_L = %s$ ohm and $X_C = %s$ ohm, so $|Z| = \\sqrt{R^2 + (X_L - X_C)^2} = \\sqrt{%s^2 + (%s)^2} = %s$ ohm.", num(xl, 2), num(xc, 2), num(r, 0), num(x, 2), num(z, 2)),
			facts("xl", num(xl, 2), "xc", num(xc, 2)))
	case 2:
		rp := float64(10 * edakit.Int(1, 100))
		y := complex(1/rp, w*c*1e-12)
		z := 1 / cmplx.Abs(y)
		return edakit.Numeric(
			sp("A %s ohm resistor is in parallel with a %s pF capacitor. What is the magnitude of the combination's impedance at %s MHz, in ohms? Round to 2 decimal places.", num(rp, 0), num(c, 0), num(f, 0)),
			num(z, 2),
			sp("Admittances add: $Y = 1/R + j\\omega C = %s + j%s$ mS. Then $|Z| = 1/|Y| = R/\\sqrt{1 + (\\omega R C)^2} = %s$ ohm.", num(1000/rp, 4), num(1000*w*c*1e-12, 4), num(z, 2)),
			facts("g_ms", num(1000/rp, 4), "b_ms", num(1000*w*c*1e-12, 4)))
	default:
		a := deg(math.Atan2(x, r))
		return edakit.Numeric(stem+" What is the phase angle of its impedance, in degrees (positive when inductive)? Round to 1 decimal place.",
			num(a, 1),
			sp("$X = X_L - X_C = %s - %s = %s$ ohm, so the angle is $\\arctan(X/R) = \\arctan(%s/%s) = %s$ degrees.", num(xl, 2), num(xc, 2), num(x, 2), num(x, 2), num(r, 0), num(a, 1)),
			facts("x_ohm", num(x, 2)))
	}
}

func resonance(ctx generator.GeneratorContext) generator.Problem {
	l := float64(edakit.Int(1, 50))
	c := float64(edakit.Int(1, 50))
	f0 := 1 / (2 * math.Pi * math.Sqrt(l*1e-9*c*1e-12)) / 1e6
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(
			sp("A series LC circuit has $L = %s$ nH and $C = %s$ pF. What is its resonant frequency, in MHz? Round to 1 decimal place.", num(l, 0), num(c, 0)),
			num(f0, 1),
			sp("At resonance $X_L = X_C$, so $f_0 = 1/(2\\pi\\sqrt{LC}) = 1/(2\\pi\\sqrt{%s \\times 10^{-9} \\cdot %s \\times 10^{-12}}) = %s$ MHz.", num(l, 0), num(c, 0), num(f0, 1)),
			facts("f0_mhz", num(f0, 1)))
	case 2:
		r := float64(edakit.Int(1, 20))
		q := math.Sqrt(1000*l/c) / r
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A series RLC circuit has $R = %s$ ohm, $L = %s$ nH and $C = %s$ pF. What is its quality factor $Q$? Round to 3 decimal places.", num(r, 0), num(l, 0), num(c, 0)),
				num(q, 3),
				sp("$Q = \\frac{1}{R}\\sqrt{L/C} = \\sqrt{%s \\times 10^{-9}/%s \\times 10^{-12}}/%s = %s$.", num(l, 0), num(c, 0), num(r, 0), num(q, 3)),
				facts("q", num(q, 3)))
		}
		bw := f0 / q
		return edakit.Numeric(
			sp("A series RLC circuit has $R = %s$ ohm, $L = %s$ nH and $C = %s$ pF. What is its 3 dB bandwidth, in MHz? Round to 2 decimal places.", num(r, 0), num(l, 0), num(c, 0)),
			num(bw, 2),
			sp("$f_0 = %s$ MHz and $Q = \\sqrt{L/C}/R = %s$, so the bandwidth is $f_0/Q = R/(2\\pi L) = %s$ MHz.", num(f0, 2), num(q, 3), num(bw, 2)),
			facts("f0_mhz", num(f0, 2), "q", num(q, 3)))
	default:
		for {
			lp := pick(0.1, 0.2, 0.25, 0.5, 1, 2)
			cp := pick(10, 20, 50, 100, 200, 500)
			rp := float64(edakit.Int(1, 20))
			q := math.Sqrt(lp*1e-9/(cp*1e-9)) / (rp * 1e-3)
			if q < 5 {
				continue
			}
			zpk := lp * 1e-9 / (rp * 1e-3 * cp * 1e-9)
			fr := 1 / (2 * math.Pi * math.Sqrt(lp*1e-9*cp*1e-9)) / 1e6
			if edakit.Int(0, 1) == 0 {
				return edakit.Numeric(
					sp("A power delivery network has %s nH of package inductance with %s milliohm of series resistance, in parallel (as seen from the die) with %s nF of on-die decoupling capacitance. At what frequency is the impedance seen by the die at its peak, in MHz? Round to 1 decimal place.", num(lp, 2), num(rp, 0), num(cp, 0)),
					num(fr, 1),
					sp("The parallel tank peaks near $f = 1/(2\\pi\\sqrt{LC}) = %s$ MHz (here $Q = %s$, so the high-Q formula is accurate).", num(fr, 1), num(q, 1)),
					facts("f_mhz", num(fr, 1), "q", num(q, 1)))
			}
			return edakit.Numeric(
				sp("A power delivery network has %s nH of package inductance with %s milliohm of series resistance, in parallel (as seen from the die) with %s nF of on-die decoupling capacitance. Using the high-Q approximation, what is the peak impedance at anti-resonance, in ohms? Round to 3 decimal places.", num(lp, 2), num(rp, 0), num(cp, 0)),
				num(zpk, 3),
				sp("At anti-resonance the tank looks like $L/(RC) = %s \\times 10^{-9}/(%s \\times 10^{-3} \\cdot %s \\times 10^{-9}) = %s$ ohm, which is $Q^2 R$ with $Q = %s$.", num(lp, 2), num(rp, 0), num(cp, 0), num(zpk, 3), num(q, 1)),
				facts("z_peak", num(zpk, 3), "q", num(q, 1)))
		}
	}
}
