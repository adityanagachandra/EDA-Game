package dft

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// testFlow: tester cost per die (KP 1), multi-site throughput (KP 2), Arrhenius burn-in
// acceleration (KP 3).
func testFlow(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		cyc := edakit.Scaled(ctx, 200000, 20000000)
		f := edakit.Pick([]int{25, 50, 100})
		ovh := edakit.Int(20, 300)
		rate := edakit.Pick([]int{2, 3, 4, 5, 8})
		sites := edakit.Pick([]int{1, 2, 4, 8})
		t := float64(cyc)/(float64(f)*1000) + float64(ovh)
		cost := t / 1000 * float64(rate) / 100 / float64(sites)
		ans := edakit.Num(cost*100, 3)
		return edakit.Numeric(edakit.Sprintf("Wafer probe applies %d scan cycles at %d MHz plus %d ms of parametric and handling overhead per touchdown, testing %d dies in parallel. The tester and prober cost %d cents per second. What is the test cost per die, in cents? Round to 3 decimal places.", cyc, f, ovh, sites, rate),
			ans, edakit.Sprintf("Touchdown time = %s ms + %d ms = %s ms; cost per touchdown = %s cents, shared by %d dies: %s cents.", edakit.Num(float64(cyc)/(float64(f)*1000), 3), ovh, edakit.Num(t, 3), edakit.Num(t/1000*float64(rate), 4), sites, ans), nil)
	case 2:
		t := float64(edakit.Int(200, 3000)) / 1000
		idx := float64(edakit.Int(100, 800)) / 1000
		sites := edakit.Pick([]int{1, 2, 4, 8, 16})
		eff := float64(edakit.Int(85, 100)) / 100
		dph := 3600 * float64(sites) * eff / (t + idx)
		ans := edakit.Num(dph, 0)
		return edakit.Numeric(edakit.Sprintf("Final test takes %s s per touchdown plus a %s s handler index time, with %d sites per touchdown and a multi-site efficiency of %s (the fraction of sites that do useful work on average). How many units per hour does the cell test? Round to the nearest integer.", edakit.Num(t, 3), edakit.Num(idx, 3), sites, edakit.Num(eff, 2)),
			ans, edakit.Sprintf("$3600 \\times %d \\times %s / (%s + %s) = %s$ units per hour.", sites, edakit.Num(eff, 2), edakit.Num(t, 3), edakit.Num(idx, 3), ans), nil)
	default:
		ea := float64(edakit.Int(5, 10)) / 10
		tu := float64(edakit.Pick([]int{40, 55, 70}))
		ts := float64(edakit.Pick([]int{105, 125, 150}))
		k := 8.617e-5
		af := math.Exp(ea / k * (1/(tu+273.15) - 1/(ts+273.15)))
		ans := edakit.Num(af, 1)
		return edakit.Numeric(edakit.Sprintf("Burn-in stresses parts at %s degrees C to accelerate a failure mechanism with activation energy %s eV, against a use temperature of %s degrees C. Using the Arrhenius model with $k = 8.617 \\times 10^{-5}$ eV/K, what is the acceleration factor? Round to 1 decimal place.", edakit.Num(ts, 0), edakit.Num(ea, 1), edakit.Num(tu, 0)),
			ans, edakit.Sprintf("$AF = \\exp\\left(\\frac{E_a}{k}\\left(\\frac{1}{%s} - \\frac{1}{%s}\\right)\\right) = %s$ with temperatures in kelvin; one burn-in hour ages the part like %s use hours.", edakit.Num(tu+273.15, 2), edakit.Num(ts+273.15, 2), ans, ans), nil)
	}
}

// defectLevel: Williams-Brown defect level (KP 1), coverage for a target (KP 2), from Poisson
// yield (KP 3).
func defectLevel(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		y := float64(edakit.Int(40, 95)) / 100
		t := float64(edakit.Int(9000, 9995)) / 10000
		dl := 1 - math.Pow(y, 1-t)
		ans := edakit.Num(dl*1e6, 0)
		return edakit.Numeric(edakit.Sprintf("A process yields %s percent good dies and the test set has %s percent fault coverage. Using the Williams-Brown model, what is the defect level of shipped parts, in defective parts per million? Round to the nearest integer.", edakit.Num(100*y, 0), edakit.Num(100*t, 2)),
			ans, edakit.Sprintf("$DL = 1 - Y^{1-T} = 1 - %s^{%s} = %s$, which is %s DPPM.", edakit.Num(y, 2), edakit.Num(1-t, 4), edakit.Num(dl, 6), ans), nil)
	case 2:
		y := float64(edakit.Int(40, 95)) / 100
		dppm := float64(edakit.Pick([]int{10, 50, 100, 200, 500, 1000}))
		dl := dppm / 1e6
		t := 1 - math.Log(1-dl)/math.Log(y)
		ans := edakit.Num(100*t, 3)
		return edakit.Numeric(edakit.Sprintf("A product's yield is %s percent and the customer requires at most %s defective parts per million. Using the Williams-Brown model, what fault coverage is needed, in percent? Round to 3 decimal places.", edakit.Num(100*y, 0), edakit.Num(dppm, 0)),
			ans, edakit.Sprintf("$T = 1 - \\ln(1-DL)/\\ln Y = 1 - \\ln(1 - %s)/\\ln %s = %s$ percent.", edakit.Num(dl, 6), edakit.Num(y, 2), ans), nil)
	default:
		a := float64(edakit.Int(50, 600))
		d0 := float64(edakit.Int(5, 40)) / 100
		t := float64(edakit.Int(9500, 9990)) / 10000
		y := math.Exp(-a / 100 * d0)
		dl := 1 - math.Pow(y, 1-t)
		ans := edakit.Num(dl*1e6, 0)
		return edakit.Numeric(edakit.Sprintf("A %s mm$^2$ die is built in a process with %s random killer defects per cm$^2$ (Poisson yield) and tested to %s percent fault coverage. Using Williams-Brown, what is the defect level, in DPPM? Round to the nearest integer.", edakit.Num(a, 0), edakit.Num(d0, 2), edakit.Num(100*t, 2)),
			ans, edakit.Sprintf("$Y = e^{-%s} = %s$; $DL = 1 - Y^{%s} = %s$ DPPM. Since $Y^{1-T} = e^{-AD_0(1-T)}$, only the untested fraction of defects reaches the customer.", edakit.Num(a/100*d0, 4), edakit.Num(y, 4), edakit.Num(1-t, 4), ans), nil)
	}
}

// diagnosis: dictionary matching (KP 1), cone intersection (KP 2), yield learning (KP 3).
func diagnosis(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		for attempt := 0; ; attempt++ {
			last := attempt >= generator.MaxDrawAttempts
			np := edakit.Scaled(ctx, 8, 14)
			nc := edakit.Scaled(ctx, 3, 5)
			obs := 0
			for p := 0; p < np; p++ {
				if edakit.Int(0, 2) == 0 {
					obs |= 1 << p
				}
			}
			if obs == 0 {
				obs = 1 << edakit.Int(0, np-1)
			}
			cands := make([]int, nc)
			score := make([]int, nc)
			var txt []string
			for i := range cands {
				x := obs
				flips := edakit.Int(0, 4)
				for j := 0; j < flips; j++ {
					x ^= 1 << edakit.Int(0, np-1)
				}
				if x == 0 {
					x = 1 << edakit.Int(0, np-1)
				}
				if last && i == 0 {
					x = obs
				} else if last && x == obs {
					x ^= 1 << edakit.Int(0, np-1)
					if x == 0 {
						x = obs ^ 1
					}
				}
				cands[i] = x
				score[i] = popcount(x ^ obs)
				txt = append(txt, edakit.Sprintf("C%d predicts {%s}", i+1, patText(x, np)))
			}
			best, bi, tie := 1<<30, -1, false
			for i, s := range score {
				if s < best {
					best, bi, tie = s, i, false
				} else if s == best {
					tie = true
				}
			}
			if tie && !last {
				continue
			}
			var st []string
			for i, s := range score {
				st = append(st, edakit.Sprintf("C%d: %d", i+1, s))
			}
			return edakit.Numeric(edakit.Sprintf("A failing die fails patterns {%s} out of P1 to P%d. A pass/fail fault dictionary lists the failing patterns each candidate fault predicts: %s. Scoring each candidate by the number of patterns where prediction and observation disagree, which candidate number matches best?", patText(obs, np), np, strings.Join(txt, "; ")),
				edakit.IntStr(bi+1), "Mismatch counts (size of the symmetric difference): "+strings.Join(st, ", ")+".", nil)
		}
	case 2:
		nIn, nG := sizeFor(ctx)
		for attempt := 0; ; attempt++ {
			c := randCircuit(nIn+1, nG+2, 3, typesBasic)
			if len(c.outs) < 2 && attempt < generator.MaxDrawAttempts {
				continue
			}
			var fails []int
			for _, o := range c.outs {
				if edakit.Int(0, 1) == 1 {
					fails = append(fails, o)
				}
			}
			if len(fails) == 0 {
				fails = []int{c.outs[0]}
			}
			cand := c.tfi(fails[0])
			for _, o := range fails[1:] {
				t := c.tfi(o)
				for x := range cand {
					if !t[x] {
						delete(cand, x)
					}
				}
			}
			var ft []string
			for _, o := range fails {
				ft = append(ft, c.names[o])
			}
			var cl []string
			var keys []int
			for x := range cand {
				keys = append(keys, x)
			}
			sort.Ints(keys)
			for _, x := range keys {
				cl = append(cl, c.names[x])
			}
			return edakit.Numeric(c.describe()+edakit.Sprintf(" On one pattern the tester sees wrong values on %s only. Under the single-fault assumption, the defect lies on a stem in the fan-in cone of every failing output. How many stems (inputs and gate outputs) are candidates?", strings.Join(ft, " and ")),
				edakit.IntStr(len(cand)), "Intersecting the fan-in cones leaves "+strings.Join(cl, ", ")+edakit.Sprintf(": %d candidates. Passing outputs cannot simply be subtracted, because a fault in their cone need not reach them on this pattern.", len(cand)), nil)
		}
	default:
		names := []string{"via opens", "metal-2 shorts", "metal-1 shorts", "poly shorts", "contact opens"}
		n := edakit.Scaled(ctx, 3, 5)
		lam := make([]float64, n)
		tot := 0.0
		var p []string
		mx := 0
		for i := 0; i < n; i++ {
			lam[i] = float64(edakit.Int(2, 30)) / 100
			tot += lam[i]
			if lam[i] > lam[mx] {
				mx = i
			}
			p = append(p, edakit.Sprintf("%s %s", names[i], edakit.Num(lam[i], 2)))
		}
		for i := 0; i < n; i++ {
			if i != mx && lam[i] == lam[mx] {
				lam[mx] += 0.01
				tot += 0.01
				p[mx] = edakit.Sprintf("%s %s", names[mx], edakit.Num(lam[mx], 2))
			}
		}
		y0 := math.Exp(-tot)
		y1 := math.Exp(-(tot - lam[mx]))
		ans := pct(100 * y1)
		return edakit.Numeric(edakit.Sprintf("Volume diagnosis and failure analysis attribute the expected number of killer defects per die to mechanisms: %s. With Poisson yield, what would the yield become if the largest mechanism were eliminated, in percent? Round to 2 decimal places.", strings.Join(p, ", ")),
			ans, edakit.Sprintf("Today $Y = e^{-%s} = %s$ percent. Removing %s leaves $\\lambda = %s$, so $Y = %s$ percent. Ranking mechanisms this way is the Pareto step of yield learning.", edakit.Num(tot, 2), pct(100*y0), names[mx], edakit.Num(tot-lam[mx], 2), ans), nil)
	}
}

func patText(x, np int) string {
	var p []string
	for i := 0; i < np; i++ {
		if x>>i&1 == 1 {
			p = append(p, edakit.Sprintf("P%d", i+1))
		}
	}
	return strings.Join(p, ", ")
}

// siliconDebugCalc (KP 2): the observation window of a trace buffer.
func siliconDebugCalc(ctx generator.GeneratorContext) generator.Problem {
	d := edakit.Pick([]int{512, 1024, 2048, 4096, 8192})
	w := edakit.Pick([]int{32, 64, 128})
	k := edakit.Int(w/4, 3*w)
	per := int(math.Ceil(float64(k) / float64(w)))
	cyc := d / per
	ans := cyc
	return edakit.Choice(edakit.Sprintf("An on-chip trace buffer has %d entries of %d bits. The debug engineer traces %d signals every cycle, and each cycle's sample occupies $\\lceil %d/%d \\rceil$ consecutive entries. After the trigger, how many consecutive cycles can be captured before the buffer is full?", d, w, k, k, w),
		edakit.IntStr(ans), edakit.Sprintf("Each sample takes %d of the %d-bit entries, so %d/%d = %d cycles fit. Tracing fewer signals or compressing them widens the window.", per, w, d, per, cyc),
		edakit.Distinct(ans, 3, d*w/k, d, d/(per+1), 2*ans)...)
}

var stilItems = []edakit.Item{
	{KP: 1, Q: "What is STIL (IEEE 1450) used for?",
		A: "Describing test patterns, signals, timing and scan procedures in a tester-independent form for transfer from ATPG to ATE", D: []string{"Describing gate-level netlists for synthesis", "Describing layout geometry for DRC", "Describing standard-cell timing like Liberty"},
		Why: "STIL is the hand-off format: ATPG writes it, and tester pattern converters translate it into each ATE's native format."},
	{KP: 1, Q: "Which STIL block assigns, for each signal, how a logical 0, 1 or compare value is turned into edges within a tester cycle?",
		A: "A WaveformTable inside a Timing block", D: []string{"The Signals block", "The PatternBurst block", "The ScanStructures block"},
		Why: "Waveform tables map waveform characters to timed events such as drive-down at 0 ns or compare-high strobe at 40 ns."},
	{KP: 1, Q: "What does a STIL SignalGroups block let a pattern do?",
		A: "Refer to a named bus of signals and assign them a vector of values in one statement", D: []string{"Define the scan-chain order", "Define the period of the clock", "Declare the tester's pin electronics"},
		Why: "Grouping pins (for example all scan-in pins) keeps vector statements compact and readable."},
	{KP: 1, Q: "Why does STIL describe scan with procedures such as load_unload and shift instead of listing every shift cycle as a vector?",
		A: "A scan load is a repeated shift; describing it once and passing chain data keeps files compact and lets the tester use scan memory", D: []string{"Testers cannot apply individual vectors", "Scan chains have no timing", "Procedures are required by the netlist"},
		Why: "Expanding every shift would multiply file size by the chain length; procedures capture the repetition explicitly."},
	{KP: 2, Q: "In a STIL vector statement, what is the difference between the waveform characters for driving 0 and for expecting L?",
		A: "0 drives the pin low; L compares the device's output against low at the strobe", D: []string{"They are synonyms", "L drives the pin low with a weaker driver", "0 is an expected value and L is a drive"},
		Why: "Stimulus and response are distinguished by character: 0/1 drive, L/H compare, X don't-care or mask."},
	{KP: 2, Q: "An ATPG tool writes STIL with X in some expect positions. What does that tell the tester?",
		A: "Do not compare those outputs on that cycle, because their values are unknown", D: []string{"Drive them to an unknown value", "Fail the device if they toggle", "Repeat the vector"},
		Why: "Unknown simulation values must be masked at the tester, or good parts would fail randomly."},
	{KP: 2, Q: "What does a PatternBurst plus PatternExec structure define?",
		A: "Which pattern blocks run, in what order, under which timing and signal groups", D: []string{"The scan cell order of each chain", "The power supplies of the tester", "The layout of the probe card"},
		Why: "PatternExec is the entry point a tester executes; it references a burst of patterns and the timing and categories to use."},
	{KP: 2, Q: "Why is the tester period in a STIL Timing block usually fixed for a scan shift and different for capture?",
		A: "Shift runs at a slow, power-safe rate while launch and capture may need functional-speed pulses", D: []string{"STIL allows only one period per file", "Capture never needs timing", "Shift must be faster than capture"},
		Why: "Separate waveform tables let one pattern file use a slow shift period and an at-speed capture sequence."},
	{KP: 3, Q: "A pattern converter turns STIL into native ATE format, and the tester has a separate deep scan memory. What does the converter do with load_unload data?",
		A: "It maps the chain data into the scan memory and the surrounding vectors into the main vector memory", D: []string{"It discards the chain data", "It expands everything into main vector memory always", "It rewrites the scan chains in the netlist"},
		Why: "Testers have specialized memory for long serial scan streams, which is why STIL keeps scan data separate from parallel vectors."},
	{KP: 3, Q: "Why did the industry move from older formats such as WGL toward STIL?",
		A: "STIL is an open IEEE standard with explicit scan, timing and procedure constructs, not tied to one vendor", D: []string{"WGL cannot describe vectors", "STIL is a binary format and so smaller", "Testers cannot read text formats"},
		Why: "A vendor-neutral, standardized format reduces translation errors between EDA tools and testers."},
	{KP: 3, Q: "Patterns pass in simulation but fail on the tester on every die in the same cycle. What should be checked first?",
		A: "The pattern's timing and waveform translation, such as strobe placement relative to output delay", D: []string{"The fault coverage of the patterns", "The yield model", "The number of scan chains"},
		Why: "Systematic failure on every die in one cycle points to a pattern or timing translation error, not to a defect."},
	{KP: 3, Q: "Which extension of STIL describes test information for embedded cores so that it can be reused at chip level?",
		A: "CTL, the Core Test Language (IEEE 1450.6)", D: []string{"SDF", "Liberty", "SPEF"},
		Why: "CTL adds core test modes, wrapper and interface information to STIL, supporting IEEE 1500 core test reuse."},
}

var siliconDebugItems = []edakit.Item{
	{KP: 1, Q: "What distinguishes post-silicon validation from manufacturing test?",
		A: "Validation hunts design bugs and electrical issues in early silicon; manufacturing test screens every part for defects in a design assumed correct", D: []string{"They are the same activity", "Manufacturing test finds design bugs", "Validation runs only on the tester"},
		Why: "Validation asks whether the design works; manufacturing test asks whether this particular die was built correctly."},
	{KP: 1, Q: "A bug appears only after billions of cycles of an operating system boot on real silicon. Why was it not found in pre-silicon simulation?",
		A: "Simulation runs many orders of magnitude slower, so such long scenarios were never executed", D: []string{"Simulation cannot model processors", "The bug is always a manufacturing defect", "Pre-silicon verification never runs software"},
		Why: "Silicon runs at full speed and reaches deep states; that speed is why validation finds bugs simulation missed, and why it observes so little."},
	{KP: 1, Q: "What is the core problem of silicon debug?",
		A: "Limited observability: only a few internal signals can be seen while the chip runs at speed", D: []string{"Too many probe points", "Simulation is faster than silicon", "Scan chains cannot be read"},
		Why: "Debug infrastructure (trace buffers, scan dumps, triggers) exists to recover internal visibility after the fact."},
	{KP: 1, Q: "How is a scan dump used in silicon debug?",
		A: "The functional clock is stopped at a trigger and the scan chains are shifted out to read every flop's state", D: []string{"Scan chains are used to load new firmware", "The dump records every cycle continuously", "It measures power supply noise"},
		Why: "Scan gives complete state at one instant; it is destructive to execution, so it is paired with triggers to stop at the right moment."},
	{KP: 1, Q: "A trigger unit compares bus values every cycle and fires when a pattern appears. Why is it essential?",
		A: "Trace and scan capture are limited in size, so they must start or stop near the event of interest", D: []string{"It increases the clock frequency", "It replaces the trace buffer", "It runs ATPG on chip"},
		Why: "A small window of observation is only useful if it is placed around the failure; programmable triggers place it."},
	{KP: 1, Q: "A failure appears at low voltage and high frequency but disappears at nominal. What kind of issue is suspected?",
		A: "An electrical or timing marginality, often investigated with shmoo plots", D: []string{"A logic design bug independent of conditions", "A broken scan chain", "A March test failure"},
		Why: "Condition-dependent failures point to speed paths, noise or power integrity; shmoo plots map pass and fail across voltage and frequency."},
	{KP: 3, Q: "What is trace signal selection trying to maximize?",
		A: "The number of untraced internal states that can be reconstructed from the traced signals", D: []string{"The number of signals per entry regardless of content", "The trace buffer depth", "The fault coverage"},
		Why: "With restoration, traced flops let tools infer neighbours' values; good selection maximizes that restoration ratio."},
	{KP: 3, Q: "Why do debug engineers reproduce a silicon failure in emulation or simulation after localizing it?",
		A: "Once the failing window is known, a pre-silicon model can give full visibility to root-cause the bug", D: []string{"To measure the power consumption", "Because silicon cannot repeat a failure", "To generate manufacturing patterns"},
		Why: "Silicon finds the bug; full-visibility models explain it. Reproducing the trigger condition bridges the two."},
	{KP: 3, Q: "A design bug is found in production silicon. Which fix avoids a full respin?",
		A: "A metal-only ECO or a firmware or microcode patch, if spare cells or patch hooks exist", D: []string{"Running more ATPG patterns", "Raising the fault coverage target", "Burn-in"},
		Why: "Spare cells, metal-only ECOs and patchable microcode are debug-for-fix provisions that turn a respin into a cheaper change."},
	{KP: 3, Q: "Why is non-determinism a major obstacle in silicon debug of multi-clock systems?",
		A: "Asynchronous clock domains make the same test take different cycle-level paths, so a failure may not reproduce", D: []string{"Scan chains randomize their contents", "Trace buffers lose data", "Triggers fire randomly"},
		Why: "Crossing timing varies run to run; debug flows force deterministic modes or record enough to replay."},
	{KP: 3, Q: "What is the difference between a test escape found in the field and a validation escape?",
		A: "A test escape is a defective part that passed manufacturing test; a validation escape is a design bug that every part shares", D: []string{"They are the same", "A validation escape affects only one part", "A test escape is always a software bug"},
		Why: "Defects vary part to part; bugs are in every part. The responses differ: better tests and screening versus design fixes."},
	{KP: 3, Q: "Why do debug architectures often reuse the JTAG TAP?",
		A: "It is already present, standard and reachable on the board, so debug registers can be accessed without new pins", D: []string{"JTAG is faster than any functional interface", "Debug logic must be on the boundary register", "IEEE 1149.1 requires debug features"},
		Why: "The TAP is the universal low-pin-count access port; debug and IJTAG instruments ride on it."},
}
