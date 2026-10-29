package analog

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func photonics(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		pl := unif(0, 13, 0.5)
		gc := unif(1.5, 5, 0.1)
		wl := unif(0.5, 3, 0.1)
		ln := unif(0.2, 2, 0.1)
		mod := unif(2, 7, 0.5)
		pr := pl - 2*gc - wl*ln - mod
		q := edakit.Sprintf("A silicon-photonic link launches %s dBm from an off-chip laser, couples in and out through two grating couplers of %s dB loss each, crosses %s cm of waveguide at %s dB/cm and passes one modulator with %s dB insertion loss. What optical power reaches the photodiode, in dBm? Round to 2 decimal places.", f(pl, 1), f(gc, 1), f(ln, 1), f(wl, 1), f(mod, 1))
		return num(q, f(pr, 2), edakit.Sprintf("Losses in dB add along the path: $%s - 2(%s) - %s \\cdot %s - %s = %s$ dBm.", f(pl, 1), f(gc, 1), f(wl, 1), f(ln, 1), f(mod, 1), f(pr, 2)), nil)
	case 2:
		lam := pick(1310.0, 1550.0)
		ng := unif(3.8, 4.4, 0.05)
		r := unif(3, 50, 0.5)
		l := 2 * math.Pi * r * 1000 // nm
		fsr := lam * lam / (ng * l)
		q := edakit.Sprintf("A ring resonator of radius %s $\\mu$m has group index $n_g = %s$ near %s nm. What is its free spectral range $\\lambda^2/(n_g L)$, with $L$ the ring circumference, in nanometres? Round to 3 decimal places.", f(r, 1), f(ng, 2), f(lam, 0))
		return num(q, f(fsr, 3), edakit.Sprintf("$L = 2\\pi R = %s\\ \\mu$m, so $\\mathrm{FSR} = %s^2/(%s \\cdot %s \\times 10^3) = %s$ nm.", f(l/1000, 2), f(lam, 0), f(ng, 2), f(l/1000, 2), f(fsr, 3)), nil)
	default:
		lam := pick(1310.0, 1550.0)
		fw := unif(10, 500, 1)
		qf := lam * 1000 / fw
		q := edakit.Sprintf("A ring resonance at %s nm has a full width at half maximum of %s pm. What is its loaded quality factor? Round to the nearest integer.", f(lam, 0), f(fw, 0))
		return num(q, f(qf, 0), edakit.Sprintf("$Q = \\lambda_0/\\Delta\\lambda = %s\\ \\mathrm{nm}/%s\\ \\mathrm{pm} = %s$. A high-Q ring is also very sensitive to fabrication and temperature shifts, which is why photonic flows need thermal tuning and variation-aware simulation.", f(lam, 0), f(fw, 0), f(qf, 0)), nil)
	}
}

func mems(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		k := unif(0.5, 50, 0.5)
		m := unif(0.5, 50, 0.5) // nanograms
		fr := math.Sqrt(k/(m*1e-12)) / (2 * math.Pi) / 1000
		q := edakit.Sprintf("A MEMS resonator behaves as a spring-mass system with stiffness $k = %s$ N/m and effective mass %s ng. What is its resonant frequency, in kHz? Round to 1 decimal place.", f(k, 1), f(m, 1))
		return num(q, f(fr, 1), edakit.Sprintf("$f_0 = \\frac{1}{2\\pi}\\sqrt{k/m} = \\frac{1}{2\\pi}\\sqrt{%s/(%s \\times 10^{-12})} = %s$ kHz.", f(k, 1), f(m, 1), f(fr, 1)), nil)
	case 2:
		a := float64(10 * edakit.Int(5, 50))
		g := unif(0.5, 4, 0.1)
		c := eps0 * a * a * 1e-12 / (g * 1e-6) * 1e15
		q := edakit.Sprintf("A square parallel-plate MEMS capacitor has plates %s $\\mu$m on a side separated by a %s $\\mu$m air gap. Ignoring fringing, what is its capacitance, in femtofarads? Use $\\varepsilon_0 = 8.854 \\times 10^{-12}$ F/m. Round to 2 decimal places.", f(a, 0), f(g, 1))
		return num(q, f(c, 2), edakit.Sprintf("$C = \\varepsilon_0 A/g = 8.854 \\times 10^{-12} \\cdot %s \\times 10^{-12} / (%s \\times 10^{-6}) = %s$ fF.", f(a*a, 0), f(g, 1), f(c, 2)), nil)
	default:
		k := unif(0.5, 20, 0.5)
		g := unif(0.5, 3, 0.1)
		a := float64(10 * edakit.Int(5, 40))
		v := math.Sqrt(8 * k * math.Pow(g*1e-6, 3) / (27 * eps0 * a * a * 1e-12))
		q := edakit.Sprintf("An electrostatic MEMS actuator is a plate of area %s $\\mu$m by %s $\\mu$m on a spring of stiffness %s N/m above a %s $\\mu$m gap. What is its pull-in voltage $\\sqrt{8kg^3/(27\\varepsilon_0 A)}$, in volts? Use $\\varepsilon_0 = 8.854 \\times 10^{-12}$ F/m. Round to 2 decimal places.", f(a, 0), f(a, 0), f(k, 1), f(g, 1))
		return num(q, f(v, 2), edakit.Sprintf("Beyond pull-in the electrostatic force grows faster than the spring force and the plate snaps down, at a deflection of one third of the gap. $V_{PI} = %s$ V.", f(v, 2)), nil)
	}
}

func power(ctx generator.GeneratorContext) generator.Problem {
	vin := pick(5.0, 12.0, 24.0, 48.0, 3.6, 3.3)
	vout := edakit.Round(vin*unif(0.1, 0.8, 0.05), 2)
	d := vout / vin
	switch whichKP(ctx) {
	case 1:
		q := edakit.Sprintf("An ideal buck converter in continuous conduction steps %s V down to %s V. What duty cycle does its high-side switch run at? Round to 4 decimal places.", f(vin, 1), f(vout, 2))
		return num(q, f(d, 4), edakit.Sprintf("Volt-second balance on the inductor gives $V_{out} = D V_{in}$, so $D = %s/%s = %s$.", f(vout, 2), f(vin, 1), f(d, 4)), nil)
	case 2:
		l := unif(0.5, 22, 0.5)
		fs := unif(0.2, 3, 0.1)
		di := (vin - vout) * d / (l * 1e-6 * fs * 1e6) * 1000
		q := edakit.Sprintf("A buck converter steps %s V down to %s V with a %s $\\mu$H inductor switching at %s MHz in continuous conduction. What is the peak-to-peak inductor current ripple, in milliamperes? Round to 1 decimal place.", f(vin, 1), f(vout, 2), f(l, 1), f(fs, 1))
		return num(q, f(di, 1), edakit.Sprintf("During the on-time $DT = %s\\ \\mu$s the inductor sees $V_{in} - V_{out} = %s$ V, so $\\Delta I = (V_{in}-V_{out})D/(Lf_s) = %s$ mA.", f(d/fs, 4), f(vin-vout, 2), f(di, 1)), nil)
	default:
		io := unif(1, 20, 0.5)
		ron := float64(edakit.Int(2, 50)) // mOhm
		tr := float64(edakit.Int(2, 30))  // ns
		fs := unif(0.1, 2, 0.1)           // MHz
		pc := io * io * ron * 1e-3 * d
		psw := 0.5 * vin * io * (2 * tr * 1e-9) * fs * 1e6
		q := edakit.Sprintf("The high-side MOSFET of a buck converter (%s V in, %s V out) carries %s A of nearly constant load current when on. It has $R_{on} = %s$ milliohms, rise and fall times of %s ns each, and switches at %s MHz. Estimating conduction loss as $I^2 R_{on} D$ and switching loss as $\\frac{1}{2} V_{in} I (t_r + t_f) f_s$, what is its total loss, in watts? Round to 3 decimal places.", f(vin, 1), f(vout, 2), f(io, 1), f(ron, 0), f(tr, 0), f(fs, 1))
		return num(q, f(pc+psw, 3), edakit.Sprintf("Conduction: $%s^2 \\cdot %s \\times 10^{-3} \\cdot %s = %s$ W. Switching: $\\frac{1}{2} \\cdot %s \\cdot %s \\cdot %s\\ \\mathrm{ns} \\cdot %s\\ \\mathrm{MHz} = %s$ W. Total %s W. Faster wide-bandgap switches (GaN, SiC) shrink the switching term.", f(io, 1), f(ron, 0), f(d, 4), f(pc, 4), f(vin, 1), f(io, 1), f(2*tr, 0), f(fs, 1), f(psw, 4), f(pc+psw, 3)), nil)
	}
}

func neuromorphic(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		var parts []string
		tot := 0.0
		for i := 0; i < n; i++ {
			g := unif(1, 100, 1)
			v := unif(0, 0.3, 0.05)
			tot += g * v
			parts = append(parts, edakit.Sprintf("Row %d drives %s V through %s $\\mu$S", i+1, f(v, 2), f(g, 0)))
		}
		q := "One column of a resistive crossbar holds its bit line at virtual ground. " + strings.Join(parts, "; ") + ". What current flows out of the column, in microamperes? Round to 2 decimal places."
		return num(q, f(tot, 2), edakit.Sprintf("Ohm's law in each cell and Kirchhoff's current law at the bit line give $I = \\sum_i G_i V_i = %s\\ \\mu$A: a dot product computed in one step.", f(tot, 2)), nil)
	case 2:
		var parts []string
		tot := 0.0
		for i := 0; i < 3; i++ {
			gp := unif(5, 100, 1)
			gn := unif(5, 100, 1)
			v := unif(0.05, 0.3, 0.05)
			tot += (gp - gn) * v
			parts = append(parts, edakit.Sprintf("input %d is %s V with $G^+ = %s$ and $G^- = %s\\ \\mu$S", i+1, f(v, 2), f(gp, 0), f(gn, 0)))
		}
		q := "An in-memory multiply-accumulate stores each signed weight as a pair of conductances on two columns and subtracts the column currents: " + strings.Join(parts, "; ") + ". What is the differential output current $I^+ - I^-$, in microamperes? Round to 2 decimal places."
		return num(q, f(tot, 2), edakit.Sprintf("$I^+ - I^- = \\sum_i (G_i^+ - G_i^-)V_i = %s\\ \\mu$A. Device-to-device conductance variation enters as weight error, so the same Pelgrom-style statistics apply.", f(tot, 2)), nil)
	default:
		r := float64(edakit.Int(10, 500)) // MOhm
		c := unif(1, 50, 1)               // pF
		vth := unif(0.2, 1, 0.05)
		i := edakit.Round(vth/r*1000*unif(1.1, 4, 0.1), 1) / 1000 // uA, a whole tenth of a nA
		tau := r * c * 1e-6                                       // s (MOhm * pF = us)
		t := tau * math.Log(i*r/(i*r-vth)) * 1000
		q := edakit.Sprintf("A leaky integrate-and-fire neuron has membrane resistance %s megohms and capacitance %s pF, starts at rest (0 V), and receives a constant input current of %s nA. It spikes when its voltage reaches %s V. How long until the first spike, in milliseconds? Round to 4 decimal places.", f(r, 0), f(c, 0), f(i*1000, 1), f(vth, 2))
		return num(q, f(t, 4), edakit.Sprintf("$V(t) = IR(1 - e^{-t/\\tau})$ with $IR = %s$ V and $\\tau = RC = %s$ ms. Solving $V = V_{th}$ gives $t = \\tau\\ln(IR/(IR - V_{th})) = %s$ ms.", f(i*r, 4), f(tau*1000, 4), f(t, 4)), nil)
	}
}
