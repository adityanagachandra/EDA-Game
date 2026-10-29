package dev

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var sourcesItems = []edakit.Item{
	{KP: 1, Q: "Two identically drawn transistors side by side have thresholds that differ by 12 mV. Which variation source is this?",
		A: "Local (within-die, random) variation, such as random dopant fluctuation", D: []string{"Global die-to-die variation", "Wafer-to-wafer variation", "Aging from bias temperature instability"},
		Why: "Neighbouring devices share the same global process conditions. A difference between them is local, random mismatch."},
	{KP: 1, Q: "What is random dopant fluctuation?",
		A: "The number and position of dopant atoms in a small channel vary from device to device, so the threshold varies", D: []string{"Dopants diffuse out of the channel during operation", "The doping profile differs between wafers in a lot", "The gate metal grains change size over time"},
		Why: "A small bulk channel holds only tens to hundreds of dopant atoms. Poisson statistics on that count give a threshold spread that grows as the device shrinks."},
	{KP: 1, Q: "Which variation source comes from the roughness of the patterned gate edge?",
		A: "Line-edge roughness, which varies the local gate length along the width", D: []string{"Random dopant fluctuation", "Oxide thickness drift across the wafer", "Metal gate granularity"},
		Why: "Lithography and etch leave a rough edge. The gate length varies along the width, and short-channel effects turn length variation into threshold variation."},
	{KP: 1, Q: "Why does metal gate granularity cause threshold variation in high-k metal-gate devices?",
		A: "The gate metal forms grains with different crystal orientations and work functions", D: []string{"Metal grains carry dopant atoms into the channel", "Grain boundaries short the gate to the channel", "Granularity changes the gate length"},
		Why: "Each orientation has its own work function, so the effective gate work function, and the threshold, depends on which grains sit over the channel."},
	{KP: 2, Q: "A wafer shows faster transistors near its centre than at its edge. How is this classified?",
		A: "Systematic, spatially correlated global variation", D: []string{"Random local mismatch", "Aging", "Measurement noise"},
		Why: "A smooth gradient across the wafer is a systematic process effect, correlated over long distances, unlike independent device-to-device mismatch."},
	{KP: 2, Q: "A layout places a transistor next to a dense block in one instance and in open space in another, and their currents differ. What kind of variation is this?",
		A: "Layout-dependent effects, such as stress and lithography proximity", D: []string{"Random dopant fluctuation", "Die-to-die variation", "Time-dependent dielectric breakdown"},
		Why: "Stress from wells and isolation, and lithographic proximity, depend on the surrounding layout. They are deterministic given the layout and are modeled by layout-dependent effect parameters."},
	{KP: 2, Q: "Why did undoped thin-body devices such as FinFETs reduce random dopant fluctuation?",
		A: "They control short-channel effects with geometry instead of heavy channel doping", D: []string{"They contain more dopant atoms, which averages out", "They have no threshold voltage", "Their gates are polysilicon"},
		Why: "With few dopants in the channel there is little to fluctuate. Remaining variation comes from fin geometry, work-function granularity and edge roughness."},
	{KP: 2, Q: "A Monte Carlo SPICE run samples local and global parameters. Why must the two be sampled differently?",
		A: "Global parameters are shared by all devices on a die, local ones are independent per device", D: []string{"Global parameters are not random", "Local parameters only apply to PMOS", "The simulator cannot sample more than one parameter"},
		Why: "Global variation shifts every device together, preserving matching. Local variation is drawn per instance. Mixing them up either destroys matching or hides it."},
	{KP: 3, Q: "A 6T SRAM bitcell fails read stability in a few cells per megabit. Which variation source matters most?",
		A: "Local threshold mismatch between the cell's transistors", D: []string{"Die-to-die variation that shifts all cells equally", "Interconnect resistance", "Package thermal resistance"},
		Why: "Read stability depends on the strength ratio of neighbouring transistors in one cell. Random mismatch in the far tail, around 5 to 6 sigma, sets the failure rate."},
	{KP: 3, Q: "Why does a timing engine treat random local variation along a deep logic path as partly averaging out?",
		A: "Independent per-gate delay deviations add in quadrature, so the relative spread of the path delay shrinks with depth", D: []string{"Random variation cancels exactly in pairs", "Local variation only affects the first gate", "Deep paths have no variation"},
		Why: "For $N$ independent stages the sigma grows as $\\sqrt N$ while the mean grows as $N$, so relative sigma falls as $1/\\sqrt N$. AOCV and POCV build this in."},
	{KP: 3, Q: "An analog designer makes a differential pair's transistors four times larger in area. What happens to their threshold mismatch sigma?",
		A: "It halves, because mismatch sigma scales as one over the square root of area", D: []string{"It quarters", "It is unchanged", "It doubles"},
		Why: "Pelgrom's law: $\\sigma_{\\Delta V_t} = A_{V_t}/\\sqrt{WL}$. Four times the area gives half the sigma."},
	{KP: 3, Q: "Why do statistical timing tools need correlation information, not just per-gate sigmas?",
		A: "Correlated variation does not average out along a path or between launch and capture paths", D: []string{"Correlation only affects power", "Sigmas cannot be added", "Correlation makes every path fail"},
		Why: "Fully correlated deviations add linearly, independent ones in quadrature. Without correlation the tool either overstates or understates the path sigma."},
}

// pelgrom asks for mismatch sigma from area (KP 1), the area for a target sigma (KP 2) or the
// three-sigma spread of a single device (KP 3).
func pelgrom(ctx generator.GeneratorContext) generator.Problem {
	a := 0.5 * float64(edakit.Int(2, 8))
	w := pickF(0.1, 0.2, 0.5, 1, 2, 4)
	l := pickF(0.02, 0.05, 0.1, 0.2, 0.5)
	sig := a / math.Sqrt(w*l)
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(sig, 2)
		q := edakit.Sprintf("A process has a Pelgrom mismatch coefficient $A_{V_t} = %s$ mV µm. What is the standard deviation of the threshold difference between two adjacent identical transistors with $W = %s$ µm and $L = %s$ µm, in mV? Round to 2 decimal places.", edakit.Num(a, 1), edakit.Num(w, 1), edakit.Num(l, 2))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\sigma_{\\Delta V_t} = A_{V_t}/\\sqrt{WL} = %s/\\sqrt{%.4g} = %s$ mV.", edakit.Num(a, 1), w*l, ans), map[string]string{"WL": edakit.Sprintf("%.4g", w*l)})
	case 2:
		target := pickF(1, 2, 3, 5)
		area := math.Pow(a/target, 2)
		ans, rt := sigAnswer(area, 3)
		q := edakit.Sprintf("A comparator needs the input-pair threshold mismatch $\\sigma_{\\Delta V_t}$ to be at most %s mV. With $A_{V_t} = %s$ mV µm, what gate area $WL$ is required, in µm$^2$? %s", edakit.Num(target, 0), edakit.Num(a, 1), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("$WL = (A_{V_t}/\\sigma)^2 = (%s/%s)^2 = %s$ µm$^2$. Halving the sigma would take four times the area.", edakit.Num(a, 1), edakit.Num(target, 0), ans), map[string]string{"area": ans})
	default:
		single := sig / math.Sqrt2
		ans := edakit.Num(3*single, 1)
		q := edakit.Sprintf("With $A_{V_t} = %s$ mV µm, a pair of transistors with $W = %s$ µm and $L = %s$ µm has $\\sigma_{\\Delta V_t} = A_{V_t}/\\sqrt{WL}$. Treating each device's local deviation as independent and identically distributed, what is the 3-sigma local threshold deviation of a single device, in mV? Round to 1 decimal place.", edakit.Num(a, 1), edakit.Num(w, 1), edakit.Num(l, 2))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\sigma_{\\Delta V_t} = %.4f$ mV. The difference of two independent deviations has $\\sqrt 2$ times the single-device sigma, so $\\sigma_{V_t} = %.4f$ mV and $3\\sigma = %s$ mV.", sig, single, ans), map[string]string{"sigma_pair": edakit.Num(sig, 4)})
	}
}

// cornerSigma combines global and local sigma into a statistical corner.
func cornerSigma(ctx generator.GeneratorContext) generator.Problem {
	vt := 10 * edakit.Int(25, 45)
	sg := float64(edakit.Int(8, 25))
	sl := float64(edakit.Int(5, 25))
	k := edakit.Pick([]int{2, 3})
	tot := math.Sqrt(sg*sg + sl*sl)
	slow := float64(vt) + float64(k)*tot
	q := edakit.Sprintf("A device's threshold has nominal value %d mV, an independent global (die-to-die) sigma of %s mV and a local (within-die) sigma of %s mV. For a single device on a random die, what threshold sits at the $+%d\\sigma$ point of the combined distribution, in mV?", vt, edakit.Num(sg, 0), edakit.Num(sl, 0), k)
	return edakit.Choice(q, edakit.Num(slow, 1),
		edakit.Sprintf("Independent Gaussian sources add in variance: $\\sigma = \\sqrt{%s^2 + %s^2} = %.3f$ mV, so the point is $%d + %d \\times %.3f = %s$ mV. Adding the sigmas linearly would treat the sources as fully correlated.", edakit.Num(sg, 0), edakit.Num(sl, 0), tot, vt, k, tot, edakit.Num(slow, 1)),
		edakit.Num(float64(vt)+float64(k)*(sg+sl), 1), edakit.Num(float64(vt)+float64(k)*sg, 1), edakit.Num(float64(vt)+float64(k)*(sg*sg+sl*sl)/(sg+sl), 1))
}

var cornerItems = []edakit.Item{
	{KP: 1, Q: "What does the SF process corner of a CMOS technology describe?",
		A: "Slow NMOS combined with fast PMOS", D: []string{"Slow wires and fast transistors", "Slow at low temperature, fast at high temperature", "A corner for scan flip-flops"},
		Why: "Corner names list the NMOS first and the PMOS second. Skewed corners such as SF and FS stress circuits whose function depends on the n-to-p strength ratio."},
	{KP: 1, Q: "Why are process corners pessimistic for a design-level timing check?",
		A: "They push every device to the same extreme at once, which real variation rarely does", D: []string{"They ignore global variation", "They are always faster than silicon", "They exclude temperature"},
		Why: "A corner is a global shift. Combining it with worst local variation on every device stacks extremes, which is why statistical methods such as POCV recover margin."},
	{KP: 1, Q: "What does a foundry's statistical SPICE model provide that a corner model does not?",
		A: "Distributions of model parameters that a Monte Carlo run samples, separating global and local parts", D: []string{"A single fixed parameter set", "Only temperature dependence", "Layout geometry"},
		Why: "Corner models fix the parameters at chosen extremes. Statistical models carry parameter sigmas and correlations for Monte Carlo and for deriving sigma-aware library data."},
	{KP: 1, Q: "A timing signoff uses SS at low voltage and high or low temperature. Why can either temperature be worst for delay?",
		A: "Temperature inversion: at low supply, the falling threshold can outweigh falling mobility as temperature rises", D: []string{"Temperature never affects delay", "Delay always rises with temperature", "Low temperature always speeds wires only"},
		Why: "Mobility falls with temperature, slowing devices, but threshold falls too, speeding them. Near threshold the second effect can dominate, so cold becomes the slow corner."},
	{KP: 3, Q: "A liberty file carries LVF sigma tables. Where do those sigmas come from?",
		A: "Monte Carlo characterization of each cell using the foundry's local-variation statistical model", D: []string{"Measurements of the final chip", "The corner model alone", "A fixed percentage of nominal delay"},
		Why: "Library characterization runs Monte Carlo SPICE per arc, slew and load, and records the delay sigma due to local variation."},
	{KP: 3, Q: "An engineer runs 1000 Monte Carlo samples and sees no failure. Can they claim a 6-sigma yield?",
		A: "No: 1000 samples barely resolve the 3-sigma tail, so rare-event methods are needed", D: []string{"Yes, any failure-free run proves 6 sigma", "Yes, if the samples were run at the SS corner", "No, because Monte Carlo cannot sample local variation"},
		Why: "A 6-sigma failure has probability about $10^{-9}$. Plain sampling needs billions of runs. Importance sampling and other tail methods target those events directly."},
	{KP: 3, Q: "Why does a corner model for a FinFET process often have smaller local sigma than older bulk processes at similar area?",
		A: "Undoped fins suffer less random dopant fluctuation", D: []string{"FinFETs have no variation", "Corners do not include local variation at all", "FinFET channels are longer"},
		Why: "Lower channel doping removes the main RDF contribution. Fin geometry and work-function variation remain."},
	{KP: 3, Q: "A Monte Carlo run of a ring oscillator samples only global parameters. What does it miss?",
		A: "Stage-to-stage mismatch, which partly averages over the stages", D: []string{"All variation, since global parameters do not vary", "Temperature dependence", "Nothing at all"},
		Why: "Global sampling shifts every stage together. Local mismatch adds per-stage randomness whose effect on the period shrinks roughly as one over the square root of the stage count."},
}

// btiPowerLaw extrapolates BTI degradation with a power law in time.
func btiPowerLaw(ctx generator.GeneratorContext) generator.Problem {
	d1 := float64(edakit.Int(5, 20))
	t1 := pickF(10, 100, 1000)
	years := pickF(5, 10)
	t2 := years * 8760
	n := pickF(0.15, 0.17, 0.2, 0.25)
	d2 := d1 * math.Pow(t2/t1, n)
	q := edakit.Sprintf("A stress test shows a %s mV threshold shift after %s hours. BTI degradation follows $\\Delta V_t = A t^{n}$ with $n = %s$. Under the same conditions, what shift is expected after %s years (8760 hours per year), in mV?",
		edakit.Num(d1, 0), edakit.Num(t1, 0), edakit.Num(n, 2), edakit.Num(years, 0))
	return edakit.Choice(q, edakit.Num(d2, 1),
		edakit.Sprintf("The ratio of times is $%s/%s = %.1f$. With a power law, $\\Delta V_t$ grows by $%.1f^{%s} = %.3f$, giving %s mV. Linear extrapolation in time would wildly overstate it.", edakit.Num(t2, 0), edakit.Num(t1, 0), t2/t1, t2/t1, edakit.Num(n, 2), math.Pow(t2/t1, n), edakit.Num(d2, 1)),
		edakit.Num(d1*math.Pow(t2/t1, 2*n), 1), edakit.Num(d1*math.Pow(t2/t1, n/2), 1), edakit.Num(d1*(1+n*math.Log10(t2/t1)), 1))
}

var btiHciItems = []edakit.Item{
	{KP: 1, Q: "Which stress condition produces negative bias temperature instability (NBTI) in a PMOS?",
		A: "Gate negative relative to the source and body, at elevated temperature, even with no drain current", D: []string{"High drain voltage with the gate half on", "Zero gate bias at low temperature", "Reverse bias on the drain junction only"},
		Why: "NBTI needs a strong vertical field across the dielectric and heat. A PMOS whose input is held low sits in exactly this state."},
	{KP: 1, Q: "What happens to a transistor's threshold and drive current under BTI stress?",
		A: "The threshold magnitude rises and the drive current falls", D: []string{"The threshold falls and current rises", "Both are unchanged; only leakage rises", "The gate capacitance doubles"},
		Why: "Interface and oxide traps capture charge that opposes the gate. $|V_t|$ rises, the overdrive shrinks and the device slows."},
	{KP: 1, Q: "Why does BTI recover partially when the stress is removed?",
		A: "Part of the trapped charge is released when the field is removed", D: []string{"The gate oxide regrows", "Dopants diffuse back into the channel", "Recovery is a measurement artifact only"},
		Why: "Some traps empty on relaxation, so measured degradation depends on the delay between stress and measurement. Duty-cycle-aware models capture this."},
	{KP: 1, Q: "Which BTI variant became significant for NMOS with high-k metal-gate stacks?",
		A: "Positive bias temperature instability (PBTI)", D: []string{"Negative BTI in NMOS", "Hot-hole injection", "Electromigration"},
		Why: "With high-k dielectrics, electron trapping in the high-k layer under positive gate bias gives an NMOS threshold shift."},
	{KP: 3, Q: "Hot-carrier injection degrades an NMOS most under which condition?",
		A: "High drain voltage with the device switching, so carriers gain energy from the lateral field near the drain", D: []string{"Gate at full voltage with $V_{DS} = 0$", "Device fully off at room temperature", "Only during power-up"},
		Why: "HCI needs current and a high lateral field together. Energetic carriers near the drain create interface traps or inject into the oxide."},
	{KP: 3, Q: "Why is HCI often worse at low temperature while BTI is worse at high temperature?",
		A: "Carriers scatter less with phonons when cold, so more reach high energies; BTI is thermally activated", D: []string{"Both are worse at low temperature", "HCI does not depend on temperature", "BTI is worse at low temperature because mobility is higher"},
		Why: "Lower phonon scattering lets carriers heat up more in the lateral field. BTI trap creation follows an Arrhenius-like temperature dependence. Advanced nodes complicate this with self-heating."},
	{KP: 3, Q: "Which circuits age most from hot-carrier injection?",
		A: "Circuits that switch often with large drain voltages, such as clock buffers and output drivers", D: []string{"Idle circuits held in one logic state", "SRAM cells that are never accessed", "Power gates that are always off"},
		Why: "HCI damage accumulates per transition, when current flows while $V_{DS}$ is high. Activity factor matters, unlike static BTI stress."},
	{KP: 3, Q: "Where along the channel does HCI damage concentrate?",
		A: "Near the drain, where the lateral field peaks", D: []string{"Near the source", "Uniformly along the channel", "In the substrate far from the surface"},
		Why: "The field is highest in the pinch-off region at the drain end, so carrier energy and damage peak there. This makes the damage asymmetric between forward and reverse operation."},
}

// tddb asks for a Weibull failure fraction (KP 1), the time to a failure fraction (KP 2) or a
// characteristic life scaled to a larger area (KP 3).
func tddb(ctx generator.GeneratorContext) generator.Problem {
	beta := pickF(1.0, 1.2, 1.5, 2.0, 2.5)
	eta := float64(10 * edakit.Int(2, 20))
	switch kpOf(ctx) {
	case 1:
		t := float64(edakit.Int(3, 15))
		f := 1 - math.Exp(-math.Pow(t/eta, beta))
		ans, rt := sigAnswer(100*f, 3)
		q := edakit.Sprintf("Gate-oxide breakdown times follow a Weibull distribution $F(t) = 1 - e^{-(t/\\eta)^{\\beta}}$ with $\\eta = %d$ years and $\\beta = %s$. What percentage of parts has failed after %d years? %s", int(eta), edakit.Num(beta, 1), int(t), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("$(t/\\eta)^\\beta = (%d/%d)^{%s} = %.5f$, so $F = 1 - e^{-%.5f} = %.5f$, which is %s%%.", int(t), int(eta), edakit.Num(beta, 1), math.Pow(t/eta, beta), math.Pow(t/eta, beta), f, ans), map[string]string{"F": edakit.Sprintf("%.5f", f)})
	case 2:
		fp := pickF(0.01, 0.1, 1)
		f := fp / 100
		t := eta * math.Pow(-math.Log(1-f), 1/beta)
		ans, rt := sigAnswer(t, 3)
		q := edakit.Sprintf("Oxide breakdown follows a Weibull distribution with $\\eta = %d$ years and $\\beta = %s$. At what time has %s%% of the population failed, in years? %s", int(eta), edakit.Num(beta, 1), edakit.Num(fp, 2), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("Invert $F = 1 - e^{-(t/\\eta)^\\beta}$: $t = \\eta(-\\ln(1-F))^{1/\\beta} = %d \\times (%.6g)^{1/%s} = %s$ years. Low failure fractions sit far below $\\eta$.", int(eta), -math.Log(1-f), edakit.Num(beta, 1), ans), map[string]string{"t": ans})
	default:
		at := pickF(1e-5, 1e-4)
		ac := pickF(0.01, 0.05, 0.1, 0.5)
		ec := eta * 100 // test-structure characteristic life in years
		e2 := ec * math.Pow(at/ac, 1/beta)
		ans, rt := sigAnswer(e2, 3)
		q := edakit.Sprintf("Test capacitors of area $%s$ cm$^2$ show Weibull oxide breakdown with characteristic life %s years (extrapolated to use conditions) and $\\beta = %s$. A chip has a total gate area of %s cm$^2$. Using weakest-link area scaling, what is the chip's characteristic life, in years? %s", sciOf(at), edakit.Num(ec, 0), edakit.Num(beta, 1), edakit.Num(ac, 2), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("Weakest-link scaling gives $\\eta_{chip} = \\eta_{test}(A_{test}/A_{chip})^{1/\\beta} = %s \\times (%.4g)^{1/%s} = %s$ years. A larger area has more weak spots, and a low $\\beta$ makes the penalty steeper.", edakit.Num(ec, 0), at/ac, edakit.Num(beta, 1), ans), map[string]string{"area_ratio": edakit.Sprintf("%.4g", at/ac)})
	}
}

// acceleration asks for an Arrhenius factor (KP 1), a voltage acceleration factor (KP 2) or a
// use-condition lifetime from both (KP 3).
func acceleration(ctx generator.GeneratorContext) generator.Problem {
	ea := pickF(0.3, 0.5, 0.7, 0.9, 1.0)
	tu := float64(edakit.Pick([]int{50, 85, 105})) + 273
	ts := float64(edakit.Pick([]int{125, 150, 175})) + 273
	af := math.Exp(ea / kB * (1/tu - 1/ts))
	g := pickF(4, 6, 8, 10, 12)
	vu := 0.05 * float64(edakit.Int(14, 20))
	vs := vu + 0.1*float64(edakit.Int(2, 6))
	av := math.Exp(g * (vs - vu))
	switch kpOf(ctx) {
	case 1:
		ans, rt := sigAnswer(af, 3)
		q := edakit.Sprintf("A wear-out mechanism has activation energy $E_a = %s$ eV. Using the Arrhenius model with $k = 8.617 \\times 10^{-5}$ eV/K, what acceleration factor does a stress at %d °C give relative to use at %d °C? %s", edakit.Num(ea, 1), int(ts-273), int(tu-273), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("$AF = \\exp[(E_a/k)(1/T_{use} - 1/T_{stress})]$ with $T_{use} = %d$ K and $T_{stress} = %d$ K: $\\exp[%.1f \\times %.6f] = %s$.", int(tu), int(ts), ea/kB, 1/tu-1/ts, ans), map[string]string{"AF_T": ans})
	case 2:
		ans, rt := sigAnswer(av, 3)
		q := edakit.Sprintf("A dielectric's lifetime follows an exponential voltage model, $t \\propto e^{-\\gamma V}$ with $\\gamma = %s$ V$^{-1}$. What acceleration factor does stressing at %s V give relative to use at %s V? %s", edakit.Num(g, 0), edakit.Num(vs, 2), edakit.Num(vu, 2), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("$AF_V = e^{\\gamma(V_{stress} - V_{use})} = e^{%s \\times %s} = %s$.", edakit.Num(g, 0), edakit.Num(vs-vu, 2), ans), map[string]string{"AF_V": ans})
	default:
		ea = pickF(0.3, 0.5, 0.7)
		af = math.Exp(ea / kB * (1/tu - 1/ts))
		g = pickF(4, 5, 6)
		vs = vu + 0.1*float64(edakit.Int(2, 4))
		av = math.Exp(g * (vs - vu))
		th := float64(100 * edakit.Int(2, 20))
		life := th * af * av / 8760
		ans, rt := sigAnswer(life, 3)
		q := edakit.Sprintf("A stress test at %d °C and %s V reaches the failure criterion after %d hours. The mechanism has $E_a = %s$ eV (Arrhenius, $k = 8.617 \\times 10^{-5}$ eV/K) and exponential voltage acceleration with $\\gamma = %s$ V$^{-1}$. Treating the two factors as independent, what is the expected lifetime at %d °C and %s V, in years (8760 hours per year)? %s",
			int(ts-273), edakit.Num(vs, 2), int(th), edakit.Num(ea, 1), edakit.Num(g, 0), int(tu-273), edakit.Num(vu, 2), rt)
		return edakit.Numeric(q, ans, edakit.Sprintf("$AF_T = %.4g$ and $AF_V = %.4g$. Lifetime $= %d \\times %.4g \\times %.4g / 8760 = %s$ years.", af, av, int(th), af, av, ans), map[string]string{"AF_T": edakit.Sprintf("%.4g", af), "AF_V": edakit.Sprintf("%.4g", av)})
	}
}

var agingItems = []edakit.Item{
	{KP: 1, Q: "How does a typical aging-aware timing flow represent degraded transistors?",
		A: "It characterizes aged cell libraries from device models whose parameters are shifted by the predicted end-of-life degradation", D: []string{"It multiplies every delay by a random number", "It reruns placement for each year of life", "It ignores aging until silicon returns"},
		Why: "Reliability simulators apply BTI and HCI shifts to model parameters for a mission profile. Aged libraries then let static timing check end-of-life performance."},
	{KP: 1, Q: "What inputs does a reliability simulator need to predict a transistor's end-of-life shift?",
		A: "The device's stress history: voltages, temperature, duty cycle and switching activity over the mission profile", D: []string{"Only the gate length", "Only the final layout", "The package pin count"},
		Why: "BTI depends on time under bias and temperature with recovery, HCI on switching events with high $V_{DS}$. Both need the waveform statistics of each device."},
	{KP: 1, Q: "Why does a single aged library for the whole chip overestimate degradation for many cells?",
		A: "Cells see very different stress: a gate held idle or rarely toggling ages differently from a clock buffer", D: []string{"Aging is identical for every cell", "Aged libraries always underestimate aging", "Aging only affects memories"},
		Why: "Stress is instance-specific. A uniform worst-case aged library is safe but pessimistic, which is why instance- or activity-aware aging analysis exists."},
	{KP: 1, Q: "What does a compact-model aging extension typically add to the model?",
		A: "Equations that compute threshold and mobility shifts from accumulated stress, applied during simulation", D: []string{"A second netlist", "Extra metal layers", "A new definition of the drain current"},
		Why: "Aging models track stress-dependent parameter drift and feed it into the base compact model, keeping the I-V equations unchanged."},
	{KP: 2, Q: "A clock tree is aged asymmetrically because one branch's buffers sit idle at logic 0 for long periods. What timing risk follows?",
		A: "Clock skew changes over life, which can create hold or setup violations that were not present at time zero", D: []string{"None: clock trees do not age", "Only power increases", "The clock period shrinks"},
		Why: "Different stress per branch shifts delays differently. Skew is a difference, so asymmetric aging moves it directly."},
	{KP: 2, Q: "A design team signs off timing with an aged library and a fresh library. Why both?",
		A: "Setup is often worst on slowed aged paths, while hold can be worst on fresh fast paths, or when aging slows the capture clock", D: []string{"Fresh libraries are only for power", "Signoff tools require exactly two libraries", "Aged libraries are only for analog"},
		Why: "Aging mostly slows devices, hurting setup, but hold checks depend on relative delays and need the fast, fresh case and asymmetric aging cases."},
	{KP: 2, Q: "Why does BTI recovery complicate the mission profile a reliability tool needs?",
		A: "The same total stress time produces different shifts depending on how stress and relaxation alternate", D: []string{"Recovery makes aging permanent", "Recovery only happens at zero temperature", "Recovery doubles the degradation"},
		Why: "Duty cycle and stress frequency matter, not just accumulated time. A static-probability profile captures much of this for logic."},
	{KP: 2, Q: "An agent proposes guardbanding aging by raising the supply voltage over the product's life. What device effect does this trade against?",
		A: "Higher voltage accelerates BTI, HCI and dielectric wear, so the compensation itself speeds aging", D: []string{"Higher voltage reverses aging", "Voltage has no effect on reliability", "Only leakage depends on voltage"},
		Why: "Adaptive voltage scaling recovers speed, but every wear mechanism is voltage-accelerated, so the schedule must be checked against lifetime models."},
	{KP: 3, Q: "Why are aging predictions extrapolated from accelerated stress tests rather than measured directly?",
		A: "A ten-year life cannot be observed before product release, so stress at high voltage and temperature compresses time", D: []string{"Use conditions cause no aging", "Accelerated tests are more accurate than use data", "Stress tests measure only leakage"},
		Why: "Models fitted at elevated stress are extrapolated with voltage and temperature acceleration laws. The extrapolation itself is a key model risk."},
	{KP: 3, Q: "What is the main risk when extrapolating a BTI power law from 1000 hours of stress to 10 years?",
		A: "The exponent or the mechanism may change outside the measured window, so the error grows with the extrapolation range", D: []string{"Power laws cannot be extrapolated at all", "The result is always pessimistic", "Time does not enter BTI models"},
		Why: "Extrapolating almost two decades in time magnifies exponent error. Saturation or new mechanisms can appear that the stress window never saw."},
	{KP: 3, Q: "Self-heating raises the channel temperature of FinFET and nanosheet devices above the ambient. Why does aging analysis care?",
		A: "BTI and other thermally activated mechanisms age devices according to the local channel temperature, not the ambient", D: []string{"Self-heating lowers aging", "Self-heating only matters for metal", "Aging uses the package temperature by definition"},
		Why: "Confined 3D channels conduct heat poorly, so the device runs hotter than its surroundings. Aging models need the local temperature."},
	{KP: 3, Q: "An SRAM bitcell stores the same value for years. Which aging effect threatens it most?",
		A: "Asymmetric BTI on the cross-coupled pair, which shifts the cell's stability margin over time", D: []string{"Hot-carrier injection from frequent switching", "Electromigration in the bitline", "Dielectric breakdown of the word line"},
		Why: "The static state stresses one PMOS and leaves its mirror unstressed. The growing mismatch erodes read and write margins."},
}
