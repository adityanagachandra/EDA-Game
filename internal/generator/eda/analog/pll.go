package analog

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func vco(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		n := pick(3, 5, 7, 9, 11, 13, 15)
		td := float64(edakit.Int(5, 80))
		fr := 1000 / (2 * float64(n) * td)
		q := edakit.Sprintf("A ring oscillator has %d identical inverting stages, each with a propagation delay of %s ps. What is its oscillation frequency, in GHz? Round to 3 decimal places.", n, f(td, 0))
		return num(q, f(fr, 3), edakit.Sprintf("An edge must travel the ring twice for one full period, so $T = 2Nt_d = %s$ ps and $f = 1/T = %s$ GHz.", f(2*float64(n)*td, 0), f(fr, 3)), nil)
	case 2:
		l := unif(0.2, 5, 0.1)
		c := unif(0.1, 5, 0.1)
		fr := 1 / (2 * math.Pi * math.Sqrt(l*1e-9*c*1e-12)) / 1e9
		q := edakit.Sprintf("An LC VCO tank has $L = %s$ nH and total capacitance $C = %s$ pF. What is its oscillation frequency, in GHz? Round to 3 decimal places.", f(l, 1), f(c, 1))
		return num(q, f(fr, 3), edakit.Sprintf("$f = 1/(2\\pi\\sqrt{LC}) = 1/(2\\pi\\sqrt{%s \\times 10^{-21}}) = %s$ GHz.", f(l*c, 2), f(fr, 3)), nil)
	default:
		v1 := unif(0.2, 0.6, 0.05)
		v2 := v1 + unif(0.1, 0.5, 0.05)
		f1 := float64(10 * edakit.Int(50, 300))
		k := float64(10 * edakit.Int(5, 150))
		f2 := f1 + k*(v2-v1)
		q := edakit.Sprintf("A VCO oscillates at %s MHz with control voltage %s V and at %s MHz with %s V, linear in between. What is its tuning gain $K_{VCO}$, in MHz/V? Round to 1 decimal place.", f(f1, 0), f(v1, 2), f(f2, 1), f(v2, 2))
		return num(q, f(k, 1), edakit.Sprintf("$K_{VCO} = \\Delta f/\\Delta V = (%s - %s)/(%s - %s) = %s$ MHz/V.", f(f2, 1), f(f1, 0), f(v2, 2), f(v1, 2), f(k, 1)), nil)
	}
}

func pllLoop(ctx generator.GeneratorContext) generator.Problem {
	fref := pick(10.0, 19.2, 20.0, 25.0, 26.0, 40.0, 50.0)
	switch whichKP(ctx) {
	case 1:
		n := edakit.Int(20, 250)
		fout := fref * float64(n)
		q := edakit.Sprintf("An integer-N PLL uses a %s MHz reference and must synthesize %s MHz. What feedback division ratio N does it lock with?", f(fref, 1), f(fout, 1))
		return num(q, edakit.IntStr(n), edakit.Sprintf("At lock the divided VCO matches the reference in frequency and phase, so $f_{out} = N f_{ref}$ and $N = %s/%s = %d$. Channel spacing equals $f_{ref}$.", f(fout, 1), f(fref, 1), n), nil)
	case 2:
		n := edakit.Int(20, 150)
		m := pick(4, 8, 16, 32, 64)
		k := edakit.Int(1, m-1)
		fout := fref * (float64(n) + float64(k)/float64(m))
		q := edakit.Sprintf("A fractional-N PLL with a %s MHz reference divides by %d for %d out of every %d reference cycles and by %d for the rest. What is its average output frequency, in MHz? Round to 4 decimal places.", f(fref, 1), n+1, k, m, n)
		return num(q, f(fout, 4), edakit.Sprintf("The average division ratio is $%d + %d/%d = %s$, so $f_{out} = %s \\cdot %s = %s$ MHz.", n, k, m, f(float64(n)+float64(k)/float64(m), 6), f(fref, 1), f(float64(n)+float64(k)/float64(m), 6), f(fout, 4)), nil)
	default:
		n := edakit.Int(20, 150)
		fout := fref * float64(n)
		kv := float64(10 * edakit.Int(10, 100))
		vc := unif(0.3, 1.2, 0.01)
		for attempt := 0; kv*vc > 0.6*fout; attempt++ {
			kv = float64(10 * edakit.Int(1, 100))
			vc = unif(0.3, 1.2, 0.01)
			if attempt >= generator.MaxDrawAttempts {
				kv, vc = 10, 0.3
			}
		}
		f0 := fout - kv*vc
		q := edakit.Sprintf("A PLL with a %s MHz reference and divide ratio %d drives a VCO that runs at %s MHz with zero control voltage and has $K_{VCO} = %s$ MHz/V. What control voltage does the loop settle to, in volts? Round to 3 decimal places.", f(fref, 1), n, f(f0, 1), f(kv, 0))
		return num(q, f(vc, 3), edakit.Sprintf("Lock forces $f_{out} = %d \\cdot %s = %s$ MHz, so $V_c = (f_{out} - f_0)/K_{VCO} = (%s - %s)/%s = %s$ V.", n, f(fref, 1), f(fout, 1), f(fout, 1), f(f0, 1), f(kv, 0), f(vc, 3)), nil)
	}
}

func loopDynamics(ctx generator.GeneratorContext) generator.Problem {
	icp := float64(10 * edakit.Int(2, 50))  // uA
	kv := float64(10 * edakit.Int(10, 100)) // MHz/V
	n := edakit.Int(10, 200)
	c := float64(10 * edakit.Int(2, 100)) // pF
	wn := math.Sqrt(icp * 1e-6 * kv * 1e6 * 2 * math.Pi / (2 * math.Pi * float64(n) * c * 1e-12))
	setup := edakit.Sprintf("A type-II charge-pump PLL has charge-pump current $I_{CP} = %s\\ \\mu$A, VCO gain $K_{VCO} = %s$ MHz/V, divide ratio $N = %d$, and a loop filter of $R$ in series with $C = %s$ pF. Take $\\omega_n = \\sqrt{I_{CP} K_{VCO}/(2\\pi N C)}$ with $K_{VCO}$ in rad/s/V, and $\\zeta = (R C/2)\\,\\omega_n$.", f(icp, 0), f(kv, 0), n, f(c, 0))
	switch whichKP(ctx) {
	case 1:
		fn := wn / (2 * math.Pi) / 1000
		q := setup + " What is the natural frequency $f_n = \\omega_n/(2\\pi)$, in kHz? Round to 2 decimal places."
		return num(q, f(fn, 2), edakit.Sprintf("$K_{VCO} = 2\\pi \\cdot %s \\times 10^6$ rad/s/V, so the $2\\pi$ cancels and $\\omega_n = \\sqrt{I_{CP}K/(NC)} = %s$ rad/s, giving $f_n = %s$ kHz.", f(kv, 0), f(wn, 0), f(fn, 2)), map[string]string{"wn": f(wn, 0)})
	case 2:
		r := edakit.Round(2*unif(0.4, 1.5, 0.01)/(c*1e-12*wn)/1000, 1)
		if r < 0.1 {
			r = 0.1
		}
		z := r * 1e3 * c * 1e-12 / 2 * wn
		q := setup + edakit.Sprintf(" With $R = %s$ kilohms, what is the damping factor $\\zeta$? Round to 3 decimal places.", f(r, 1))
		return num(q, f(z, 3), edakit.Sprintf("$\\omega_n = %s$ rad/s and $RC/2 = %s$ s, so $\\zeta = %s$. Values near 0.7 to 1 give a well-damped lock transient.", f(wn, 0), sci(r*1e3*c*1e-12/2), f(z, 3)), nil)
	default:
		zt := pick(0.6, 0.707, 0.8, 1.0, 1.2)
		r := 2 * zt / (c * 1e-12 * wn) / 1000
		q := setup + edakit.Sprintf(" What resistance $R$ gives $\\zeta = %s$, in kilohms? Round to 2 decimal places.", f(zt, 3))
		return num(q, f(r, 2), edakit.Sprintf("$\\omega_n = %s$ rad/s, and $R = 2\\zeta/(C\\omega_n) = %s$ kilohms. Raising $R$ adds the stabilizing zero's phase lead.", f(wn, 0), f(r, 2)), nil)
	}
}

func phaseNoise(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		l1 := -float64(edakit.Int(80, 120))
		f1 := pick(10.0, 100.0, 1000.0)
		f2 := f1 * float64(edakit.Int(2, 30))
		l2 := l1 - db20(f2/f1)
		q := edakit.Sprintf("A free-running VCO has phase noise %s dBc/Hz at a %s kHz offset, in its $1/f^2$ region. What is its phase noise at a %s kHz offset, in dBc/Hz? Round to 1 decimal place.", f(l1, 0), f(f1, 0), f(f2, 0))
		return num(q, f(l2, 1), edakit.Sprintf("In the $1/f^2$ region phase noise falls 20 dB per decade: $%s - 20\\log_{10}(%s/%s) = %s$ dBc/Hz.", f(l1, 0), f(f2, 0), f(f1, 0), f(l2, 1)), nil)
	case 2:
		lr := -float64(edakit.Int(130, 160))
		n := edakit.Int(10, 500)
		lo := lr + db20(float64(n))
		q := edakit.Sprintf("Inside its loop bandwidth, a PLL output tracks its reference. The reference phase noise is %s dBc/Hz at the offset of interest and the divide ratio is %d. What is the in-band output phase noise contributed by the reference, in dBc/Hz? Round to 1 decimal place.", f(lr, 0), n)
		return num(q, f(lo, 1), edakit.Sprintf("Multiplying frequency by $N$ multiplies phase deviations by $N$, adding $20\\log_{10}(%d) = %s$ dB: %s dBc/Hz.", n, f(db20(float64(n)), 2), f(lo, 1)), nil)
	default:
		l := -float64(edakit.Int(90, 130))
		b := float64(edakit.Int(1, 50)) // MHz
		f0 := unif(1, 10, 0.1)          // GHz
		phiRms := math.Sqrt(2 * undb10(l) * b * 1e6)
		j := phiRms / (2 * math.Pi * f0 * 1e9) * 1e15
		q := edakit.Sprintf("A %s GHz clock has flat single-sideband phase noise of %s dBc/Hz from near the carrier out to %s MHz offset, and negligible noise beyond. What is its RMS jitter, in femtoseconds? Use $\\sigma_\\phi^2 = 2\\int \\mathcal{L}(f)\\,df$. Round to 1 decimal place.", f(f0, 1), f(l, 0), f(b, 0))
		return num(q, f(j, 1), edakit.Sprintf("$\\sigma_\\phi = \\sqrt{2 \\cdot 10^{%s/10} \\cdot %s \\times 10^6} = %s$ rad, and $\\sigma_t = \\sigma_\\phi/(2\\pi f_0) = %s$ fs.", f(l, 0), f(b, 0), sci(phiRms), f(j, 1)), nil)
	}
}
