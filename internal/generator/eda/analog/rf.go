package analog

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func sParams(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		zl := float64(edakit.Int(5, 300))
		if zl == 50 {
			zl = 75
		}
		g := math.Abs((zl - 50) / (zl + 50))
		rl := -db20(g)
		q := edakit.Sprintf("An LNA input presents a purely resistive %s ohm impedance to a 50 ohm system. What is its return loss, in dB? Round to 2 decimal places.", f(zl, 0))
		return num(q, f(rl, 2), edakit.Sprintf("$|\\Gamma| = |S_{11}| = |Z_L - Z_0|/(Z_L + Z_0) = %s$, so $\\mathrm{RL} = -20\\log_{10}|\\Gamma| = %s$ dB.", f(g, 4), f(rl, 2)), nil)
	case 2:
		s11 := -unif(3, 30, 0.5)
		g := undb20(s11)
		v := (1 + g) / (1 - g)
		q := edakit.Sprintf("A port measures $|S_{11}| = %s$ dB. What is its voltage standing-wave ratio? Round to 3 decimal places.", f(s11, 1))
		return num(q, f(v, 3), edakit.Sprintf("$|\\Gamma| = 10^{%s/20} = %s$, and $\\mathrm{VSWR} = (1+|\\Gamma|)/(1-|\\Gamma|) = %s$.", f(s11, 1), f(g, 4), f(v, 3)), nil)
	default:
		a := unif(-8, 8, 0.1)
		b := unif(-8, 8, 0.1)
		if math.Abs(a)+math.Abs(b) < 0.5 {
			a = 3.2
		}
		m := math.Hypot(a, b)
		bs := f(b, 1)
		sign := "+"
		if b < 0 {
			sign = "-"
			bs = f(-b, 1)
		}
		q := edakit.Sprintf("A two-port amplifier measured in a 50 ohm system has $S_{21} = %s %s j%s$ at the frequency of interest. What is its gain $20\\log_{10}|S_{21}|$, in dB? Round to 2 decimal places.", f(a, 1), sign, bs)
		return num(q, f(db20(m), 2), edakit.Sprintf("$|S_{21}| = \\sqrt{%s^2 + %s^2} = %s$, which is %s dB. The phase of $S_{21}$ is the insertion phase and does not affect gain.", f(a, 1), f(b, 1), f(m, 4), f(db20(m), 2)), nil)
	}
}

func noiseFigure(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		sin := unif(20, 60, 0.5)
		nf := unif(0.5, 10, 0.1)
		sout := sin - nf
		F := undb10(nf)
		q := edakit.Sprintf("With a 290 K source, a block's input SNR is %s dB and its output SNR is %s dB. What is its noise factor $F$ as a linear ratio? Round to 3 decimal places.", f(sin, 1), f(sout, 1))
		return num(q, f(F, 3), edakit.Sprintf("$\\mathrm{NF} = \\mathrm{SNR}_{in} - \\mathrm{SNR}_{out} = %s$ dB, and $F = 10^{\\mathrm{NF}/10} = %s$.", f(nf, 1), f(F, 3)), nil)
	case 2:
		nf := unif(0.3, 8, 0.1)
		te := (undb10(nf) - 1) * 290
		q := edakit.Sprintf("An amplifier has a noise figure of %s dB. What is its equivalent input noise temperature, in kelvin? Round to 1 decimal place.", f(nf, 1))
		return num(q, f(te, 1), edakit.Sprintf("$T_e = (F - 1) \\cdot 290 = (10^{%s/10} - 1) \\cdot 290 = %s$ K.", f(nf, 1), f(te, 1)), nil)
	default:
		nf := unif(1, 12, 0.5)
		bw := pick(0.2, 1.0, 5.0, 10.0, 20.0, 100.0)
		snr := unif(0, 20, 0.5)
		p := -174 + db10(bw*1e6) + nf + snr
		q := edakit.Sprintf("A receiver has noise figure %s dB over a %s MHz channel and needs %s dB of output SNR to demodulate. Taking the thermal floor as $-174$ dBm/Hz, what is its sensitivity (minimum input power), in dBm? Round to 1 decimal place.", f(nf, 1), f(bw, 1), f(snr, 1))
		return num(q, f(p, 1), edakit.Sprintf("$P_{min} = -174 + 10\\log_{10}(%s \\times 10^6) + %s + %s = %s$ dBm.", f(bw, 1), f(nf, 1), f(snr, 1), f(p, 1)), nil)
	}
}

func friis(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		nf1 := unif(0.5, 4, 0.1)
		g1 := unif(8, 25, 0.5)
		nf2 := unif(4, 15, 0.5)
		F := undb10(nf1) + (undb10(nf2)-1)/undb10(g1)
		q := edakit.Sprintf("An LNA with noise figure %s dB and gain %s dB drives a mixer with noise figure %s dB. What is the cascade noise figure, in dB? Round to 2 decimal places.", f(nf1, 1), f(g1, 1), f(nf2, 1))
		return num(q, f(db10(F), 2), edakit.Sprintf("Friis: $F = F_1 + (F_2 - 1)/G_1 = %s + (%s - 1)/%s = %s$, which is %s dB.", f(undb10(nf1), 4), f(undb10(nf2), 4), f(undb10(g1), 3), f(F, 4), f(db10(F), 2)), nil)
	case 2:
		nf1 := unif(0.5, 3, 0.1)
		g1 := unif(10, 20, 0.5)
		nf2 := unif(5, 12, 0.5)
		g2 := unif(-8, 12, 0.5)
		nf3 := unif(8, 20, 0.5)
		F := undb10(nf1) + (undb10(nf2)-1)/undb10(g1) + (undb10(nf3)-1)/undb10(g1+g2)
		q := edakit.Sprintf("A receiver chain is LNA (NF %s dB, gain %s dB), then mixer (NF %s dB, gain %s dB), then IF amplifier (NF %s dB). What is the cascade noise figure, in dB? Round to 2 decimal places.", f(nf1, 1), f(g1, 1), f(nf2, 1), f(g2, 1), f(nf3, 1))
		return num(q, f(db10(F), 2), edakit.Sprintf("$F = F_1 + (F_2-1)/G_1 + (F_3-1)/(G_1G_2) = %s + %s + %s = %s$, so NF = %s dB. Each later stage is divided by all the gain ahead of it.", f(undb10(nf1), 4), f((undb10(nf2)-1)/undb10(g1), 4), f((undb10(nf3)-1)/undb10(g1+g2), 4), f(F, 4), f(db10(F), 2)), nil)
	default:
		loss := unif(0.5, 4, 0.1)
		nf2 := unif(0.8, 4, 0.1)
		F := undb10(loss) + (undb10(nf2)-1)*undb10(loss)
		q := edakit.Sprintf("A passive preselect filter with %s dB insertion loss (at 290 K) sits in front of an LNA with noise figure %s dB. What is the noise figure of the combination, in dB? Round to 2 decimal places.", f(loss, 1), f(nf2, 1))
		return num(q, f(db10(F), 2), edakit.Sprintf("A matched passive loss $L$ at 290 K has $F_1 = L$ and $G_1 = 1/L$. Friis gives $F = L + (F_2 - 1)L = L F_2$, so the noise figures add: $%s + %s = %s$ dB.", f(loss, 1), f(nf2, 1), f(db10(F), 2)), nil)
	}
}

func iip3(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		pin := -unif(10, 40, 1)
		g := unif(5, 20, 0.5)
		d := unif(20, 70, 1)
		v := pin + d/2
		q := edakit.Sprintf("In a two-tone test each tone enters an amplifier at %s dBm. At the output each fundamental is %s dBm and each third-order intermodulation product is %s dBm. What is the input-referred third-order intercept IIP3, in dBm? Round to 1 decimal place.", f(pin, 0), f(pin+g, 1), f(pin+g-d, 1))
		return num(q, f(v, 1), edakit.Sprintf("The fundamental rises 1 dB per dB and IM3 rises 3 dB per dB, so they meet $\\Delta/2$ above the input: $\\mathrm{IIP3} = P_{in} + \\Delta P/2 = %s + %s/2 = %s$ dBm.", f(pin, 0), f(d, 0), f(v, 1)), nil)
	case 2:
		ip := unif(-15, 15, 0.5)
		g := unif(5, 25, 0.5)
		pin := ip - unif(15, 35, 1)
		im := 3*pin - 2*ip + g
		q := edakit.Sprintf("An amplifier has IIP3 = %s dBm and gain %s dB. Two equal tones each enter at %s dBm. What output power does each IM3 product have, in dBm? Round to 1 decimal place.", f(ip, 1), f(g, 1), f(pin, 1))
		return num(q, f(im, 1), edakit.Sprintf("Input-referred IM3 is $3P_{in} - 2\\,\\mathrm{IIP3} = %s$ dBm; adding the gain gives %s dBm at the output.", f(3*pin-2*ip, 1), f(im, 1)), nil)
	default:
		ip1 := unif(-15, 5, 0.5)
		g1 := unif(10, 20, 0.5)
		ip2 := unif(-5, 15, 0.5)
		inv := 1/undb10(ip1) + undb10(g1)/undb10(ip2)
		tot := db10(1 / inv)
		q := edakit.Sprintf("An LNA (IIP3 %s dBm, gain %s dB) drives a mixer (IIP3 %s dBm). Assuming the IM3 products add in phase, what is the cascade IIP3, in dBm? Round to 2 decimal places.", f(ip1, 1), f(g1, 1), f(ip2, 1))
		return num(q, f(tot, 2), edakit.Sprintf("In milliwatts, $1/\\mathrm{IIP3} = 1/%s + %s/%s = %s$, so IIP3 $= %s$ mW, which is %s dBm. Gain ahead of a stage makes its linearity matter more, the reverse of noise.", f(undb10(ip1), 4), f(undb10(g1), 3), f(undb10(ip2), 4), f(inv, 4), f(1/inv, 4), f(tot, 2)), nil)
	}
}
