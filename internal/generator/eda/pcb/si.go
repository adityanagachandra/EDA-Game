package pcb

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// reflection: coefficient (KP 1), first-incident load voltage (KP 2), termination or transmission (KP 3).
func reflection(ctx generator.GeneratorContext) generator.Problem {
	z0 := float64(edakit.Pick([]int{40, 45, 50, 50, 55, 60, 75, 85, 90, 100}))
	switch kpOf(ctx) {
	case 1:
		zl := z0 + 20
		for a := 0; a < generator.MaxDrawAttempts; a++ {
			if c := float64(edakit.Pick([]int{20, 25, 33, 40, 45, 50, 60, 68, 75, 100, 120, 150, 220, 1000})); c != z0 {
				zl = c
				break
			}
		}
		g := (zl - z0) / (zl + z0)
		ans := edakit.Num(g, 3)
		q := edakit.Sprintf("A %s ohm trace ends in a %s ohm load. The reflection coefficient is $\\Gamma = (Z_L - Z_0)/(Z_L + Z_0)$. What is $\\Gamma$? Round to 3 decimal places.", edakit.Num(z0, 0), edakit.Num(zl, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\Gamma = (%s - %s)/(%s + %s) = %s$. A load above $Z_0$ reflects a same-sign wave; a load below it reflects an inverted one.", edakit.Num(zl, 0), edakit.Num(z0, 0), edakit.Num(zl, 0), edakit.Num(z0, 0), ans), map[string]string{"gamma": ans})
	case 2:
		vs := edakit.Pick([]float64{1.0, 1.2, 1.5, 1.8, 2.5, 3.3})
		rs := float64(edakit.Int(8, 40))
		zl := float64(edakit.Pick([]int{25, 33, 75, 100, 150, 1000, 10000}))
		v1 := vs * z0 / (rs + z0)
		g := (zl - z0) / (zl + z0)
		vl := v1 * (1 + g) * 1000
		ans := edakit.Num(vl, 0)
		q := edakit.Sprintf("A driver with a %s V step and %s ohm output resistance launches into a %s ohm trace that ends in a %s ohm load. Using $V_1 = V_s Z_0/(R_s + Z_0)$ for the launched step and $\\Gamma = (Z_L - Z_0)/(Z_L + Z_0)$ at the load, what voltage appears at the load when the step first arrives, in mV? Round to the nearest mV.", edakit.Num(vs, 1), edakit.Num(rs, 0), edakit.Num(z0, 0), edakit.Num(zl, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$V_1 = %s \\cdot %s/%s = %s$ V. $\\Gamma_L = %s$, so the load sees $V_1(1 + \\Gamma_L) = %s$ mV until the reflection returns from the source.", edakit.Num(vs, 1), edakit.Num(z0, 0), edakit.Num(rs+z0, 0), edakit.Num(v1, 4), edakit.Num(g, 4), ans), map[string]string{"v1": edakit.Num(v1, 4), "gamma_load": edakit.Num(g, 4)})
	}
	if edakit.Int(1, 2) == 1 {
		rs := float64(edakit.Int(6, 30))
		ans := edakit.Num(z0-rs, 0)
		q := edakit.Sprintf("A point-to-point net on a %s ohm trace is series-terminated at the driver, whose output resistance is %s ohms. What series resistor makes the source match the line, in ohms?", edakit.Num(z0, 0), edakit.Num(rs, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("The source resistance plus the resistor should equal $Z_0$: $%s - %s = %s$ ohms. The launched wave is then half the swing, doubles at the open receiver, and the reflection is absorbed at the source.", edakit.Num(z0, 0), edakit.Num(rs, 0), ans), nil)
	}
	z2 := z0 + 15
	for a := 0; a < generator.MaxDrawAttempts; a++ {
		if c := float64(edakit.Pick([]int{30, 35, 40, 45, 60, 65, 70, 80})); c != z0 {
			z2 = c
			break
		}
	}
	vi := float64(edakit.Pick([]int{200, 400, 500, 800, 1000}))
	vt := vi * 2 * z2 / (z0 + z2)
	ans := edakit.Num(vt, 1)
	q := edakit.Sprintf("A %s mV step travels on a %s ohm trace into a %s ohm connector section. The transmitted fraction is $1 + \\Gamma = 2 Z_2/(Z_1 + Z_2)$. What step amplitude enters the connector, in mV? Round to 1 decimal place.", edakit.Num(vi, 0), edakit.Num(z0, 0), edakit.Num(z2, 0))
	return edakit.Numeric(q, ans, edakit.Sprintf("$1 + \\Gamma = 2 \\cdot %s / %s = %s$, so %s mV enters. The difference from %s mV is reflected back toward the driver.", edakit.Num(z2, 0), edakit.Num(z0+z2, 0), edakit.Num(2*z2/(z0+z2), 4), ans, edakit.Num(vi, 0)), nil)
}

// ibisImpedance is the computed KP 2 of si.ibis: the linear-region output resistance from a pulldown table.
func ibisImpedance(ctx generator.GeneratorContext) generator.Problem {
	r := float64(edakit.Int(15, 60))
	v1, v2 := 0.2, 0.4
	i1, i2 := v1/r*1000, v2/r*1000
	isat := i2 * (1.6 + float64(edakit.Int(0, 6))/10)
	v3 := edakit.Pick([]float64{1.2, 1.5, 1.8})
	rbig := v3 / (isat / 1000)
	ans, ds := options(r, 1, "ohms", rbig, r/2, 2*r, r+50)
	q := edakit.Sprintf("An IBIS [Pulldown] table (typ column) lists %s V at %s mA, %s V at %s mA and %s V at %s mA, the last point in saturation. What output resistance does the driver present in its linear region near 0 V?", edakit.Num(v1, 1), edakit.Num(i1, 2), edakit.Num(v2, 1), edakit.Num(i2, 2), edakit.Num(v3, 1), edakit.Num(isat, 2))
	return edakit.Choice(q, ans, edakit.Sprintf("The slope between the first two points is $%s\\ \\text{V} / %s\\ \\text{mA}$ = %s. Dividing the saturated voltage by its current gives a large-signal ratio, not the small-signal resistance that sets the source reflection.", edakit.Num(v2-v1, 1), edakit.Num(i2-i1, 2), ans), ds...)
}

var ibisItems = []edakit.Item{
	{KP: 1, Q: "What does an IBIS model describe about an I/O buffer?",
		A:   "Its behavior at the pins: V-I tables, V-t switching waveforms, die capacitance and package parasitics",
		D:   []string{"Its transistor netlist and process parameters", "Its RTL source", "Its Liberty timing arcs for the core"},
		Why: "IBIS is a behavioral, table-driven description of what the pin does, not how the transistors are built."},
	{KP: 1, Q: "Which IBIS tables describe the protection clamp diodes?",
		A: "[POWER Clamp] and [GND Clamp]", D: []string{"[Pullup] and [Pulldown]", "[Ramp] alone", "[Package] alone"},
		Why: "Clamp tables give the current that flows when the pin is driven above the supply or below ground, which shapes overshoot."},
	{KP: 1, Q: "What do the typ, min and max columns of an IBIS table represent?",
		A:   "Corners of process, voltage and temperature for the buffer's behavior",
		D:   []string{"Three different pins of the package", "Three package types", "Three operating frequencies"},
		Why: "Running SI at all three corners brackets the weakest-slowest and strongest-fastest behavior."},
	{KP: 1, Q: "Why do silicon vendors publish IBIS models instead of SPICE netlists?",
		A:   "IBIS hides the transistor design and simulates much faster while still predicting the waveform at the pin",
		D:   []string{"IBIS is more accurate than SPICE in every case", "SPICE cannot model I/O buffers", "IBIS contains the layout database"},
		Why: "Behavioral tables protect IP and make board-level simulation of many nets practical."},
	{KP: 1, Q: "What is C_comp in an IBIS model?",
		A: "The die capacitance at the pad, excluding the package", D: []string{"The package pin capacitance", "The capacitance of the load", "The decoupling capacitor value"},
		Why: "C_comp is the silicon pad and buffer capacitance. Package R, L and C are listed separately."},
	{KP: 3, Q: "Where do plain IBIS V-I and V-t tables fall short?",
		A:   "Modeling adaptive equalization in multi-gigabit SerDes, which needs algorithmic models (IBIS-AMI)",
		D:   []string{"Simulating a 100 MHz single-ended bus", "Representing package parasitics", "Representing clamp diodes"},
		Why: "Tables capture an analog driver well. Equalizers and clock recovery are algorithms and need executable models."},
	{KP: 3, Q: "A board SI simulation with an IBIS model at the typ corner only passes. What is missing?",
		A:   "Runs at the min and max corners, where edge rate and drive strength change ringing and timing",
		D:   []string{"Nothing, since typ represents the part", "A layout DRC run", "A BOM review"},
		Why: "The fast-strong corner causes the worst overshoot and crosstalk; the slow-weak corner causes the worst timing."},
	{KP: 3, Q: "Why do power-aware IBIS features matter for simultaneous switching noise?",
		A:   "Supply-rail droop changes the buffer's drive, so the model must tie output behavior to its supply current",
		D:   []string{"They are needed to write Gerber files", "They set the thermal relief geometry", "They generate footprints"},
		Why: "When many outputs switch at once the local rail sags, and an ideal-supply model misses the resulting edge slowdown and noise."},
	{KP: 3, Q: "When is the lumped R_pkg, L_pkg and C_pkg of an IBIS file too coarse?",
		A:   "For fast edges where package traces are electrically long and pin-to-pin coupling matters",
		D:   []string{"For any DC operating point", "For slow edges on short traces", "For clamp-current calculation"},
		Why: "Then per-pin package models or package S-parameters replace the single lumped section."},
}

// touchstone: parameter counts (KP 1), insertion loss in dB (KP 2), return loss of a mismatch (KP 3).
func touchstone(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Pick([]int{1, 2, 3, 4, 6, 8, 12, 16})
		if edakit.Int(1, 2) == 1 {
			q := edakit.Sprintf("A Touchstone file named channel.s%dp holds the data of a %d-port network. How many complex S-parameters does it store at each frequency point?", n, n)
			return edakit.Numeric(q, edakit.IntStr(n*n), edakit.Sprintf("An $N$-port network has an $N \\times N$ scattering matrix: $%d^2 = %d$ entries per frequency.", n, n*n), nil)
		}
		q := edakit.Sprintf("A Touchstone file for %s %d-port network stores each S-parameter as a magnitude and an angle. Counting the frequency value itself, how many real numbers does one frequency point carry?", an(n), n)
		return edakit.Numeric(q, edakit.IntStr(2*n*n+1), edakit.Sprintf("$%d^2 = %d$ complex entries, each two reals, plus the frequency: $2 \\cdot %d + 1 = %d$.", n, n*n, n*n, 2*n*n+1), nil)
	case 2:
		if edakit.Int(1, 2) == 1 {
			m := float64(edakit.Int(5, 95)) / 100
			ans := edakit.Num(-20*log10(m), 2)
			q := edakit.Sprintf("At the Nyquist frequency a channel's Touchstone data gives $|S_{21}| = %s$. Insertion loss is $-20 \\log_{10} |S_{21}|$. What is the insertion loss, in dB? Round to 2 decimal places.", edakit.Num(m, 2))
			return edakit.Numeric(q, ans, edakit.Sprintf("$-20 \\log_{10} %s = %s$ dB.", edakit.Num(m, 2), ans), nil)
		}
		il := float64(edakit.Int(5, 300)) / 10
		ans := edakit.Num(math.Pow(10, -il/20), 3)
		q := edakit.Sprintf("A channel shows %s dB of insertion loss at one frequency, where insertion loss is $-20 \\log_{10} |S_{21}|$. What is $|S_{21}|$? Round to 3 decimal places.", edakit.Num(il, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("$|S_{21}| = 10^{-%s/20} = %s$, the fraction of the incident voltage wave that arrives.", edakit.Num(il, 1), ans), nil)
	}
	zl := float64(edakit.Pick([]int{30, 35, 40, 42, 45, 47, 53, 55, 58, 60, 65, 70, 75, 80}))
	g := absf((zl - 50) / (zl + 50))
	ans := edakit.Num(-20*log10(g), 2)
	q := edakit.Sprintf("A one-port measurement with a 50 ohm reference sees a %s ohm load, so $|S_{11}| = |Z_L - 50|/(Z_L + 50)$. Return loss is $-20 \\log_{10}|S_{11}|$. What is the return loss, in dB? Round to 2 decimal places.", edakit.Num(zl, 0))
	return edakit.Numeric(q, ans, edakit.Sprintf("$|S_{11}| = %s$, so the return loss is %s dB. A larger return loss means a better match.", edakit.Num(g, 4), ans), map[string]string{"s11": edakit.Num(g, 4)})
}

// ffeOutput is the computed KP 2 of si.ibis_ami: one output sample of a 3-tap transmit FFE.
func ffeOutput(ctx generator.GeneratorContext) generator.Problem {
	pre := -float64(edakit.Int(0, 15)) / 100
	post := -float64(edakit.Int(5, 30)) / 100
	main := 1 + pre + post
	if main < 0.5 {
		main = 0.5
	}
	xs := []int{edakit.Pick([]int{-1, 1}), edakit.Pick([]int{-1, 1}), edakit.Pick([]int{-1, 1})}
	// xs[0] = x[n-1], xs[1] = x[n], xs[2] = x[n+1]
	y := pre*float64(xs[2]) + main*float64(xs[1]) + post*float64(xs[0])
	swapped := pre*float64(xs[0]) + main*float64(xs[1]) + post*float64(xs[2])
	noPre := main*float64(xs[1]) + post*float64(xs[0])
	sum := (pre + main + post) * float64(xs[1])
	sgn := func(v int) string {
		if v > 0 {
			return "+1"
		}
		return "-1"
	}
	f := func(x float64) string { return edakit.Num(x, 2) }
	ans := f(y)
	seen := map[string]bool{ans: true}
	var ds []string
	for _, w := range []float64{swapped, noPre, sum, -y, y + 0.2, y - 0.2} {
		s := f(w)
		if !seen[s] && len(ds) < 4 {
			seen[s] = true
			ds = append(ds, s)
		}
	}
	q := edakit.Sprintf("A transmitter FFE computes $y[n] = c_{-1} x[n+1] + c_0 x[n] + c_1 x[n-1]$ with taps $c_{-1} = %s$, $c_0 = %s$, $c_1 = %s$. The symbols are $x[n-1] = %s$, $x[n] = %s$, $x[n+1] = %s$. What is $y[n]$, in units of the full swing?", f(pre), f(main), f(post), sgn(xs[0]), sgn(xs[1]), sgn(xs[2]))
	return edakit.Choice(q, ans, edakit.Sprintf("$%s \\cdot (%s) + %s \\cdot (%s) + %s \\cdot (%s) = %s$. The pre-cursor tap multiplies the next symbol and the post-cursor tap the previous one; a transition is emphasized, a repeated symbol de-emphasized.", paren(f(pre)), sgn(xs[2]), f(main), sgn(xs[1]), paren(f(post)), sgn(xs[0]), ans), ds...)
}

var amiItems = []edakit.Item{
	{KP: 1, Q: "What does IBIS-AMI add to an IBIS buffer model?",
		A:   "Algorithmic models of SerDes equalization and clock recovery, as a compiled library plus a parameter file",
		D:   []string{"A transistor-level netlist of the SerDes", "A Touchstone model of the package", "A Liberty timing file"},
		Why: "The analog part stays an IBIS buffer; the AMI part is executable code the simulator calls, configured by an .ami parameter file."},
	{KP: 1, Q: "Which entry points does an AMI executable model expose to the channel simulator?",
		A: "AMI_Init, AMI_GetWave and AMI_Close", D: []string{"read_liberty and report_timing", "main and exit only", "init_spice and tran"},
		Why: "Init processes an impulse response, GetWave processes waveform blocks, and Close releases the model's memory."},
	{KP: 1, Q: "In a statistical, Init-only AMI flow, what does the model process?",
		A:   "The channel impulse response, returning an equalized impulse response, under a linear time-invariant assumption",
		D:   []string{"A long bit-by-bit waveform with adaptation", "The board's Gerber files", "The bill of materials"},
		Why: "If the equalizer is linear and fixed, its effect on the impulse response is enough to build an eye and BER statistically."},
	{KP: 1, Q: "When is a time-domain GetWave flow needed instead of Init-only?",
		A:   "When the equalizer adapts or acts nonlinearly, as with DFE adaptation and clock-data recovery",
		D:   []string{"When the channel is linear and the equalizer is fixed", "Only for DC analysis", "Never, since Init covers every case"},
		Why: "Decision feedback and CDR depend on the actual bit decisions over time, which an impulse response alone cannot capture."},
	{KP: 3, Q: "Why do SerDes vendors ship compiled AMI models rather than source code?",
		A:   "To protect proprietary equalization algorithms while letting any compliant simulator run them",
		D:   []string{"Source code cannot be executed", "Compiled code is always more accurate", "The IBIS standard forbids source code"},
		Why: "The standard API is a seam: simulators call the same interface, and the vendor's algorithm stays a black box."},
	{KP: 3, Q: "A channel's loss at Nyquist is too high for the link. Which AMI-modeled settings are tuned first?",
		A: "Transmitter FFE taps and receiver CTLE and DFE settings", D: []string{"The IBIS clamp tables", "The package R_pkg value", "The copper weight alone"},
		Why: "Equalization compensates frequency-dependent loss. The AMI parameter file exposes those settings for sweeps."},
	{KP: 3, Q: "A channel simulation reports an eye opening and a BER contour. Which inputs produced it?",
		A:   "Transmitter IBIS and AMI models, the passive channel as S-parameters, and receiver IBIS and AMI models",
		D:   []string{"Only the Gerber files", "Only the netlist", "Only the IBIS V-I tables"},
		Why: "The flow cascades transmitter, Touchstone channel and receiver, then measures the equalized eye."},
	{KP: 3, Q: "Why can the same AMI models give different eyes in two simulators?",
		A:   "The simulators differ in how they combine the analog channel, jitter and statistical or time-domain assumptions around the models",
		D:   []string{"AMI models cannot be shared between tools", "Simulators ignore the models", "AMI models are plain text only"},
		Why: "The model API is standard, but the surrounding analysis (jitter injection, sampling, BER extrapolation) is each tool's own."},
	{KP: 3, Q: "Which part of the pulse response does a DFE cancel?",
		A: "Post-cursor intersymbol interference, using decisions on previous bits", D: []string{"Pre-cursor intersymbol interference", "Random jitter", "Only DC offset"},
		Why: "A DFE subtracts weighted past decisions. It cannot see future bits, so pre-cursor ISI is left to the transmitter FFE or CTLE."},
}

// pdnTarget: target impedance (KP 1), capacitor resonance (KP 2), capacitor count for an inductive limit (KP 3).
func pdnTarget(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		v := edakit.Pick([]float64{0.8, 0.9, 1.0, 1.1, 1.2, 1.8, 3.3})
		r := float64(edakit.Pick([]int{2, 3, 5}))
		i := float64(edakit.Int(1, 40))
		z := v * r / 100 / i * 1000
		ans := edakit.Num(z, 2)
		q := edakit.Sprintf("A %s V rail may ripple by at most %s%% while its load current steps by %s A. The PDN target impedance is $Z_t = V \\cdot \\text{ripple} / \\Delta I$. What is $Z_t$, in milliohms? Round to 2 decimal places.", edakit.Num(v, 1), edakit.Num(r, 0), edakit.Num(i, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$%s \\cdot %s / %s = %s$ milliohms. The PDN must stay below this from DC up to the frequency where the package and die take over.", edakit.Num(v, 1), edakit.Num(r/100, 2), edakit.Num(i, 0), ans), nil)
	case 2:
		type cap struct {
			text string
			f    float64
		}
		c := edakit.Pick([]cap{{"10 nF", 10e-9}, {"100 nF", 100e-9}, {"1 uF", 1e-6}, {"4.7 uF", 4.7e-6}, {"10 uF", 10e-6}, {"22 uF", 22e-6}})
		l := float64(edakit.Int(3, 20)) / 10
		f := 1 / (2 * math.Pi * sqrt(l*1e-9*c.f)) / 1e6
		ans := edakit.Num(f, 2)
		q := edakit.Sprintf("A %s capacitor mounted on the board has %s nH of total inductance (its ESL plus the mounting loop). Its series resonance is $f = 1/(2\\pi\\sqrt{LC})$. What is $f$, in MHz? Round to 2 decimal places.", c.text, edakit.Num(l, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("With $L = %s$ nH and $C$ = %s, $f = 1/(2\\pi\\sqrt{LC}) = %s$ MHz. Above it the capacitor behaves as an inductor, so mounting inductance, not capacitance, limits it.", edakit.Num(l, 1), c.text, ans), nil)
	}
	f := float64(edakit.Int(2, 20)) * 10
	l := float64(edakit.Int(4, 15)) / 10
	zt := float64(edakit.Int(5, 50))
	need := 2 * math.Pi * f * 1e6 * l * 1e-9 / (zt / 1000)
	n := ceilInt(need)
	q := edakit.Sprintf("Above resonance each mounted decoupling capacitor looks like an inductance of %s nH, and $n$ identical capacitors in parallel give $L/n$. The PDN must stay below %s milliohms at %s MHz, where $|Z| = 2\\pi f L / n$. What is the smallest number of capacitors that meets the target?", edakit.Num(l, 1), edakit.Num(zt, 0), edakit.Num(f, 0))
	return edakit.Numeric(q, edakit.IntStr(n), edakit.Sprintf("One capacitor gives $2\\pi \\cdot %s\\ \\text{MHz} \\cdot %s\\ \\text{nH} = %s$ milliohms. Dividing by the %s milliohm target gives %s, so %d capacitors.", edakit.Num(f, 0), edakit.Num(l, 1), edakit.Num(need*zt, 1), edakit.Num(zt, 0), edakit.Num(need, 3), n), map[string]string{"n": edakit.IntStr(n)})
}
