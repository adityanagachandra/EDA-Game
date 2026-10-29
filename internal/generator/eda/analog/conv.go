package analog

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func adc(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 14)
	vref := pick(0.8, 1.0, 1.2, 1.8, 2.5, 3.3)
	lsb := vref / math.Pow(2, float64(n)) * 1000
	switch whichKP(ctx) {
	case 1:
		q := edakit.Sprintf("%s %d-bit ADC has a full-scale input range of 0 to %s V. What is one LSB, in millivolts? Round to 4 decimal places.", capA(n), n, f(vref, 1))
		return num(q, f(lsb, 4), edakit.Sprintf("$\\mathrm{LSB} = V_{FS}/2^N = %s/%d = %s$ mV.", f(vref, 1), 1<<n, f(lsb, 4)), nil)
	case 2:
		n = edakit.Scaled(ctx, 3, 10)
		lsb = vref / math.Pow(2, float64(n))
		code := edakit.Int(1, (1<<n)-2)
		vin := (float64(code) + unif(0.1, 0.9, 0.1)) * lsb
		vin = edakit.Round(vin, 4)
		c := int(math.Floor(vin / lsb))
		q := edakit.Sprintf("An ideal %d-bit ADC with range 0 to %s V truncates: code $k$ means $k \\cdot \\mathrm{LSB} \\le V_{in} < (k+1)\\mathrm{LSB}$. What decimal output code does it produce for $V_{in} = %s$ V?", n, f(vref, 1), f(vin, 4))
		return num(q, edakit.IntStr(c), edakit.Sprintf("$\\mathrm{LSB} = %s$ V and $V_{in}/\\mathrm{LSB} = %s$, so the code is %d (binary %s).", f(lsb, 6), f(vin/lsb, 3), c, binStr(c, n)), nil)
	default:
		n = edakit.Scaled(ctx, 4, 10)
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("How many comparators does %s %d-bit flash ADC need?", lowA(n), n)
			return num(q, edakit.IntStr((1<<n)-1), edakit.Sprintf("A flash converter compares the input against every threshold at once: $2^N - 1 = %d$ comparators. Exponential growth is why flash stops near 6 to 8 bits.", (1<<n)-1), nil)
		}
		if n%2 == 1 {
			n++
		}
		h := n / 2
		c := 2 * ((1 << h) - 1)
		q := edakit.Sprintf("%s %d-bit two-step (subranging) ADC resolves %d coarse bits with a flash, subtracts the coarse DAC output, then resolves %d fine bits with a second flash (no redundancy). How many comparators does it use in total?", capA(n), n, h, h)
		return num(q, edakit.IntStr(c), edakit.Sprintf("Each %d-bit flash needs $2^{%d} - 1 = %d$ comparators, so two of them need %d, against %d for a single %d-bit flash.", h, h, (1<<h)-1, c, (1<<n)-1, n), nil)
	}
}

func binStr(c, n int) string {
	var b strings.Builder
	for i := n - 1; i >= 0; i-- {
		if c&(1<<i) != 0 {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}

func quantSNR(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(4, 18)
	switch whichKP(ctx) {
	case 1:
		s := 6.02*float64(n) + 1.76
		q := edakit.Sprintf("What is the ideal signal-to-quantization-noise ratio of %s %d-bit converter for a full-scale sine wave, in dB? Round to 2 decimal places.", lowA(n), n)
		return num(q, f(s, 2), edakit.Sprintf("$\\mathrm{SQNR} = 6.02N + 1.76 = 6.02 \\cdot %d + 1.76 = %s$ dB.", n, f(s, 2)), nil)
	case 2:
		n = edakit.Int(6, 14)
		vref := pick(1.0, 1.2, 1.8, 2.0, 3.3)
		lsb := vref / math.Pow(2, float64(n))
		rms := lsb / math.Sqrt(12) * 1e6
		q := edakit.Sprintf("An ideal %d-bit ADC spans %s V. Modelling the quantization error as uniform over one LSB, what is its RMS value, in microvolts? Round to 2 decimal places.", n, f(vref, 1))
		return num(q, f(rms, 2), edakit.Sprintf("$\\mathrm{LSB} = %s\\ \\mu$V, and a uniform error on $[-\\Delta/2, \\Delta/2]$ has variance $\\Delta^2/12$, so the RMS is $\\Delta/\\sqrt{12} = %s\\ \\mu$V.", f(lsb*1e6, 3), f(rms, 2)), nil)
	default:
		osr := pick(2, 4, 8, 16, 32, 64, 128, 256)
		s := 6.02*float64(n) + 1.76 + db10(float64(osr))
		q := edakit.Sprintf("%s %d-bit quantizer samples at an oversampling ratio of %d, and the out-of-band quantization noise is filtered away digitally (no noise shaping). What is the in-band SQNR for a full-scale sine, in dB? Round to 2 decimal places.", capA(n), n, osr)
		return num(q, f(s, 2), edakit.Sprintf("White quantization noise spreads over $f_s/2$; keeping a band $f_s/(2\\cdot\\mathrm{OSR})$ removes a factor OSR of it: $6.02 \\cdot %d + 1.76 + 10\\log_{10}(%d) = %s$ dB, about half a bit per doubling.", n, osr, f(s, 2)), nil)
	}
}

// sarTrace runs an ideal binary-search SAR conversion and returns the trial voltages and bits.
func sarTrace(vin, vref float64, n int) ([]float64, []int) {
	code := 0
	var trials []float64
	var bits []int
	for i := n - 1; i >= 0; i-- {
		t := code | (1 << i)
		v := float64(t) * vref / math.Pow(2, float64(n))
		trials = append(trials, v)
		if vin >= v {
			code = t
			bits = append(bits, 1)
		} else {
			bits = append(bits, 0)
		}
	}
	return trials, bits
}

func sar(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 8)
	vref := pick(1.0, 1.024, 2.0, 2.048, 1.6)
	vin := edakit.Round(unif(0.02, 0.98, 0.001)*vref, 4)
	trials, bits := sarTrace(vin, vref, n)
	trace := func(k int) string {
		var s []string
		for i := 0; i < k; i++ {
			s = append(s, edakit.Sprintf("cycle %d trial %s mV gives %d", i+1, f(trials[i]*1000, 3), bits[i]))
		}
		return strings.Join(s, "; ")
	}
	switch whichKP(ctx) {
	case 1:
		k := edakit.Int(2, n)
		q := edakit.Sprintf("An ideal %d-bit SAR ADC with reference %s V converts $V_{in} = %s$ V. Each cycle sets the next bit to 1, compares $V_{in}$ against the DAC voltage, and keeps the bit if $V_{in} \\ge V_{DAC}$. What DAC trial voltage does cycle %d compare against, in millivolts? Round to 3 decimal places.", n, f(vref, 3), f(vin, 4), k)
		return num(q, f(trials[k-1]*1000, 3), edakit.Sprintf("Trace the binary search: %s.", trace(k)), nil)
	case 2:
		code := 0
		for _, b := range bits {
			code = 2*code + b
		}
		q := edakit.Sprintf("An ideal %d-bit SAR ADC with reference %s V converts $V_{in} = %s$ V by binary search, keeping each bit when $V_{in} \\ge V_{DAC}$. What decimal output code results?", n, f(vref, 3), f(vin, 4))
		return num(q, edakit.IntStr(code), edakit.Sprintf("%s. The bits read %s, which is %d.", trace(n), binStr(code, n), code), nil)
	default:
		n = edakit.Int(8, 16)
		samp := edakit.Int(1, 3)
		fclk := float64(10 * edakit.Int(5, 100))
		fs := fclk / float64(n+samp)
		q := edakit.Sprintf("%s %d-bit synchronous SAR ADC spends %d clock %s sampling and then one clock cycle per bit decision. With a %s MHz clock, what is its maximum sample rate, in MS/s? Round to 3 decimal places.", capA(n), n, samp, plural(samp, "cycle"), f(fclk, 0))
		return num(q, f(fs, 3), edakit.Sprintf("One conversion takes $%d + %d = %d$ cycles, so $f_s = %s/%d = %s$ MS/s. A flash converter would need one cycle but $2^{%d} - 1$ comparators.", samp, n, n+samp, f(fclk, 0), n+samp, f(fs, 3), n), nil)
	}
}

func sampling(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		c := unif(0.05, 10, 0.05)
		t := float64(edakit.Int(250, 400))
		v := math.Sqrt(kBoltzmann*t/(c*1e-12)) * 1e6
		q := edakit.Sprintf("A track-and-hold samples onto a %s pF capacitor at %s K. What is the RMS of the sampled kT/C noise, in microvolts? Use $k = 1.380649 \\times 10^{-23}$ J/K. Round to 2 decimal places.", f(c, 2), f(t, 0))
		return num(q, f(v, 2), edakit.Sprintf("$v_n = \\sqrt{kT/C} = \\sqrt{%s / %s} = %s\\ \\mu$V. The switch resistance cancels: a larger resistance gives more noise density but less bandwidth.", sci(kBoltzmann*t), sci(c*1e-12), f(v, 2)), nil)
	case 2:
		n := edakit.Int(8, 14)
		vfs := pick(1.0, 1.2, 1.8, 2.0)
		lsb := vfs / math.Pow(2, float64(n))
		cmin := 12 * kBoltzmann * 300 / (lsb * lsb) * 1e15
		q := edakit.Sprintf("%s %d-bit ADC spans %s V. What sampling capacitance makes the kT/C noise power equal the quantization noise power $\\mathrm{LSB}^2/12$ at 300 K, in femtofarads? Use $k = 1.380649 \\times 10^{-23}$ J/K. Round to 2 decimal places.", capA(n), n, f(vfs, 1))
		return num(q, f(cmin, 2), edakit.Sprintf("Set $kT/C = \\Delta^2/12$: $C = 12kT/\\Delta^2$ with $\\Delta = %s\\ \\mu$V, giving %s fF. Every extra bit quadruples it.", f(lsb*1e6, 3), f(cmin, 2)), nil)
	default:
		fin := float64(edakit.Int(1, 500))
		sj := unif(0.05, 5, 0.05)
		s := -db20(2 * math.Pi * fin * 1e6 * sj * 1e-12)
		q := edakit.Sprintf("A full-scale %s MHz sine is sampled with RMS clock jitter of %s ps. What is the jitter-limited SNR, in dB? Round to 2 decimal places.", f(fin, 0), f(sj, 2))
		return num(q, f(s, 2), edakit.Sprintf("$\\mathrm{SNR} = -20\\log_{10}(2\\pi f_{in}\\sigma_t) = -20\\log_{10}(%s) = %s$ dB. Doubling the input frequency costs 6 dB.", sci(2*math.Pi*fin*1e6*sj*1e-12), f(s, 2)), nil)
	}
}

func dac(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		n := edakit.Int(6, 16)
		m := edakit.Int(0, n)
		el := (1 << m) - 1 + (n - m)
		var kind string
		switch m {
		case 0:
			kind = "fully binary-weighted (no thermometer MSBs)"
		case n:
			kind = "fully thermometer-coded"
		default:
			kind = edakit.Sprintf("segmented with the top %d bits thermometer-coded and the remaining %d bits binary-weighted", m, n-m)
		}
		q := edakit.Sprintf("%s %d-bit current-steering DAC is %s. How many separately switched current sources does it have?", capA(n), n, kind)
		return num(q, edakit.IntStr(el), edakit.Sprintf("The thermometer segment needs $2^{%d} - 1 = %d$ equal sources, and each binary bit one source: $%d + %d = %d$.", m, (1<<m)-1, (1<<m)-1, n-m, el), nil)
	case 2:
		levels := []float64{0}
		lsb := unif(1, 5, 0.5)
		for i := 1; i < 8; i++ {
			levels = append(levels, levels[i-1]+lsb*unif(0.7, 1.3, 0.05))
		}
		ideal := (levels[7] - levels[0]) / 7
		k := edakit.Int(1, 7)
		dnl := (levels[k]-levels[k-1])/ideal - 1
		var parts []string
		for i, v := range levels {
			parts = append(parts, edakit.Sprintf("code %d: %s mV", i, f(v, 3)))
		}
		q := "A 3-bit DAC's measured outputs are " + strings.Join(parts, ", ") + edakit.Sprintf(". Using the endpoint LSB (full span divided by 7), what is the DNL of the step from code %d to code %d, in LSB? Round to 3 decimal places.", k-1, k)
		return num(q, f(dnl, 3), edakit.Sprintf("Endpoint LSB $= %s/7 = %s$ mV. The step is %s mV, so $\\mathrm{DNL} = %s/%s - 1 = %s$ LSB.", f(levels[7], 3), f(ideal, 4), f(levels[k]-levels[k-1], 3), f(levels[k]-levels[k-1], 3), f(ideal, 4), f(dnl, 3)), nil)
	default:
		n := edakit.Int(4, 14)
		su := unif(0.1, 2, 0.1)
		s := su / 100 * math.Sqrt(math.Pow(2, float64(n))-1)
		q := edakit.Sprintf("%s %d-bit binary-weighted DAC is built from unit elements with independent relative mismatch $\\sigma_u = %s$%% (each bit $k$ uses $2^k$ units). At the mid-scale transition the MSB turns on and all lower bits turn off. What is the standard deviation of that step's error, in LSB? Round to 4 decimal places.", capA(n), n, f(su, 1))
		return num(q, f(s, 4), edakit.Sprintf("The step swaps $2^{%d}$ units for $2^{%d} - 1$ different units, so $2^{%d} - 1$ independent unit errors add in variance: $\\sigma = \\sigma_u\\sqrt{2^N - 1} = %s$ LSB. A thermometer DAC switches one unit per step, so its step error is just $\\sigma_u$.", n-1, n-1, n, f(s, 4)), nil)
	}
}

func enobFOM(ctx generator.GeneratorContext) generator.Problem {
	sndr := unif(30, 95, 0.1)
	enob := (sndr - 1.76) / 6.02
	switch whichKP(ctx) {
	case 1:
		q := edakit.Sprintf("A converter measures an SNDR (SINAD) of %s dB with a full-scale sine. What is its effective number of bits? Round to 2 decimal places.", f(sndr, 1))
		return num(q, f(enob, 2), edakit.Sprintf("$\\mathrm{ENOB} = (\\mathrm{SNDR} - 1.76)/6.02 = (%s - 1.76)/6.02 = %s$ bits.", f(sndr, 1), f(enob, 2)), nil)
	case 2:
		p := unif(0.1, 50, 0.1)
		fs := float64(edakit.Int(1, 500))
		fom := p * 1e-3 / (math.Pow(2, enob) * fs * 1e6) * 1e15
		q := edakit.Sprintf("An ADC sampling at %s MS/s consumes %s mW and achieves SNDR = %s dB. What is its Walden figure of merit $P/(2^{\\mathrm{ENOB}} f_s)$, in femtojoules per conversion step? Round to 2 decimal places.", f(fs, 0), f(p, 1), f(sndr, 1))
		return num(q, f(fom, 2), edakit.Sprintf("$\\mathrm{ENOB} = %s$, $2^{\\mathrm{ENOB}} = %s$, so $\\mathrm{FOM} = %s\\ \\mathrm{mW}/(%s \\cdot %s\\ \\mathrm{MS/s}) = %s$ fJ per step.", f(enob, 4), f(math.Pow(2, enob), 1), f(p, 1), f(math.Pow(2, enob), 1), f(fs, 0), f(fom, 2)), nil)
	default:
		p := unif(0.1, 50, 0.1)
		fs := float64(edakit.Int(1, 500))
		fom := sndr + db10(fs*1e6/2/(p*1e-3))
		q := edakit.Sprintf("An ADC sampling at %s MS/s (Nyquist bandwidth $f_s/2$) consumes %s mW and achieves SNDR = %s dB. What is its Schreier figure of merit $\\mathrm{SNDR} + 10\\log_{10}(\\mathrm{BW}/P)$, with BW in Hz and P in watts, in dB? Round to 2 decimal places.", f(fs, 0), f(p, 1), f(sndr, 1))
		return num(q, f(fom, 2), edakit.Sprintf("$\\mathrm{BW}/P = %s/%s = %s$, so the FOM is $%s + %s = %s$ dB.", sci(fs*1e6/2), sci(p*1e-3), sci(fs*1e6/2/(p*1e-3)), f(sndr, 1), f(db10(fs*1e6/2/(p*1e-3)), 2), f(fom, 2)), nil)
	}
}

// capA and lowA return the indefinite article for "<n>-bit".
func capA(n int) string {
	if n == 8 || n == 11 || n == 18 {
		return "An"
	}
	return "A"
}

func lowA(n int) string {
	if n == 8 || n == 11 || n == 18 {
		return "an"
	}
	return "a"
}
