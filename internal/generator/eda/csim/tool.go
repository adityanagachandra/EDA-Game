package csim

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var nodePool = []string{"in", "out", "mid", "a", "b", "tap", "n7", "x1", "fb", "drv"}

// spiceNetlist draws a small deck and asks how many nodes or MNA unknowns it creates.
func spiceNetlist(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	n := edakit.Scaled(ctx, 3, 6)
	names := edakit.Shuffle(nodePool)[:n]
	count := map[byte]int{}
	var cards []string
	add := func(kind byte, a, b, val string) {
		count[kind]++
		cards = append(cards, edakit.Sprintf("%c%d %s %s %s", kind, count[kind], a, b, val))
	}
	add('V', names[0], "0", edakit.Pick([]string{"DC 1.8", "DC 3.3", "PULSE(0 1 0 10p 10p 1n 2n)", "SIN(0 0.5 1G)"}))
	for i := 1; i < n; i++ {
		if edakit.Int(0, 3) == 0 {
			add('L', names[i-1], names[i], edakit.Pick([]string{"10n", "1u", "2.2n"}))
		} else {
			add('R', names[i-1], names[i], edakit.Pick([]string{"1k", "2.2k", "470", "10k"}))
		}
	}
	add('C', names[n-1], "0", edakit.Pick([]string{"1p", "100f", "2p"}))
	for i := 1; i < n; i++ {
		switch edakit.Int(0, 3) {
		case 0:
			add('C', names[i], "0", edakit.Pick([]string{"5f", "20f", "1p"}))
		case 1:
			add('R', names[i], "0", edakit.Pick([]string{"100k", "1meg", "50k"}))
		}
	}
	if edakit.Int(0, 1) == 0 {
		add('I', "0", names[edakit.Int(1, n-1)], "DC 10u")
	}
	if edakit.Int(0, 1) == 0 && k >= 2 {
		add('V', names[edakit.Int(1, n-1)], names[0], "DC 0")
	}
	nodes := n
	if k == 3 {
		buf := "buf"
		add('E', buf, "0", names[edakit.Int(1, n-1)]+" 0 "+edakit.Pick([]string{"2", "10", "0.5"}))
		add('R', buf, "0", "1k")
		nodes++
		if edakit.Int(0, 1) == 0 {
			add('G', names[edakit.Int(1, n-1)], "0", buf+" 0 1m")
		}
	}
	deck := strings.Join(edakit.Shuffle(cards), "; ")
	facts := map[string]string{"nodes": edakit.IntStr(nodes), "v_sources": edakit.IntStr(count['V']), "inductors": edakit.IntStr(count['L']), "vcvs": edakit.IntStr(count['E'])}
	switch k {
	case 1:
		q := "A SPICE deck contains these element cards (node 0 is ground): " + deck + ". How many distinct non-ground nodes does the circuit have?"
		return edakit.Numeric(q, edakit.IntStr(nodes),
			edakit.Sprintf("Collect every node name on the cards and drop 0. The names are %s, so there are %d.", strings.Join(names, ", "), nodes), facts)
	case 2:
		u := nodes + count['V'] + count['L']
		q := "A SPICE deck contains these element cards (node 0 is ground): " + deck + ". In the modified nodal analysis that SPICE builds for the DC operating point, how many unknowns does the system have?"
		return edakit.Numeric(q, edakit.IntStr(u),
			edakit.Sprintf("There are %d node voltages, plus one branch current for each of the %d voltage sources and %d inductors. Resistors, capacitors and current sources add no unknowns: %d.", nodes, count['V'], count['L'], u), facts)
	}
	u := nodes + count['V'] + count['L'] + count['E']
	q := "A SPICE deck contains these element cards (node 0 is ground; E is a voltage-controlled voltage source, G a voltage-controlled current source): " + deck + ". How many unknowns does the modified nodal analysis system have?"
	return edakit.Numeric(q, edakit.IntStr(u),
		edakit.Sprintf("Node voltages: %d (including buf). Branch currents: %d voltage sources, %d inductors and %d VCVS. A VCCS and current sources add none. Total %d.", nodes, count['V'], count['L'], count['E'], u), facts)
}

type suffix struct {
	text  string
	scale float64
}

// scaleSuffixes asks for the value of a SPICE number written with a scale suffix.
func scaleSuffixes(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		mant := edakit.Pick([]float64{1, 2.2, 3.3, 4.7, 10, 15, 22, 33, 47, 68, 100, 470})
		opts := []struct {
			elem, sfx, ansUnit string
			sfxScale, ansScale float64
		}{
			{"R1 a b", "k", "ohms", 1e3, 1},
			{"R1 a b", "meg", "kilohms", 1e6, 1e3},
			{"C1 a 0", "n", "picofarads", 1e-9, 1e-12},
			{"C1 a 0", "p", "femtofarads", 1e-12, 1e-15},
			{"C1 a 0", "u", "nanofarads", 1e-6, 1e-9},
			{"L1 a b", "n", "microhenries", 1e-9, 1e-6},
			{"L1 a b", "u", "nanohenries", 1e-6, 1e-9},
			{"C1 a 0", "f", "picofarads", 1e-15, 1e-12},
		}
		o := edakit.Pick(opts)
		v := mant * o.sfxScale / o.ansScale
		return num(edakit.Sprintf("A SPICE card reads %s %s%s. What is the element value in %s?", o.elem, f(mant, 2), o.sfx, o.ansUnit), v, 4,
			edakit.Sprintf("The suffix %s scales by %s, so in %s the value is %s times %s, which is %s.", o.sfx, sci(o.sfxScale), o.ansUnit, f(mant, 2), sci(o.sfxScale/o.ansScale), f(v, 6)), nil)
	case 2:
		mant := edakit.Pick([]float64{1, 2, 2.2, 3.3, 4.7, 5, 10})
		opts := []struct {
			text  string
			scale float64
		}{{"M", 1e-3}, {"MEG", 1e6}, {"Meg", 1e6}, {"m", 1e-3}, {"Mohm", 1e-3}, {"MEGohm", 1e6}}
		o := edakit.Pick(opts)
		if edakit.Int(0, 2) == 0 {
			// The F trap: F is femto, so a capacitor written as 10F is 10 fF.
			fo := edakit.Pick([]struct {
				text  string
				scale float64
			}{{"F", 1e-15}, {"pF", 1e-12}, {"nF", 1e-9}, {"fF", 1e-15}, {"uF", 1e-6}})
			v := mant * fo.scale / 1e-15
			return num(edakit.Sprintf("A SPICE card reads C1 x 0 %s%s. Remembering that SPICE reads only the scale letter and ignores trailing unit letters, what capacitance does the simulator use, in femtofarads?", f(mant, 2), fo.text), v, 3,
				edakit.Sprintf("In %s%s the first letter after the number sets the scale (%s); F alone is femto, and any letters after the scale letter are ignored. The value is %s fF.", f(mant, 2), fo.text, sci(fo.scale), f(v, 3)), nil)
		}
		v := mant * o.scale
		return num(edakit.Sprintf("A SPICE card reads R1 x y %s%s. SPICE suffixes are case-insensitive. What resistance does the simulator use, in ohms?", f(mant, 2), o.text), v, 4,
			edakit.Sprintf("M and m both mean milli; mega must be written MEG. So %s%s is %s times %s, which is %s ohms.", f(mant, 2), o.text, f(mant, 2), sci(o.scale), f(v, 4)), nil)
	}
	r := edakit.Pick([]float64{1, 2.2, 3.3, 4.7, 10, 22, 47})
	rs := edakit.Pick([]suffix{{"k", 1e3}, {"", 1}, {"meg", 1e6}})
	c := edakit.Pick([]float64{1, 2.2, 4.7, 10, 47, 100, 470})
	cs := edakit.Pick([]suffix{{"p", 1e-12}, {"f", 1e-15}, {"n", 1e-9}})
	tau := r * rs.scale * c * cs.scale / 1e-9
	if tau > 1e6 || tau < 1e-3 {
		rs, cs = suffix{"k", 1e3}, suffix{"p", 1e-12}
		tau = r * c * 1e3 * 1e-12 / 1e-9
	}
	pl := placesFor(tau)
	return num(edakit.Sprintf("A deck contains R1 in out %s%s and C1 out 0 %s%s. What is the RC time constant in nanoseconds?", f(r, 2), rs.text, f(c, 2), cs.text), tau, pl,
		edakit.Sprintf("R = %s ohms and C = %s F, so RC = %s s, which is %s ns.", sci(r*rs.scale), sci(c*cs.scale), sci(tau*1e-9), f(tau, pl)), map[string]string{"tau_ns": f(tau, pl)})
}

var modelCardItems = []edakit.Item{
	{KP: 1, Q: "A deck has the line M1 out in 0 0 nch W=1u L=0.1u and a separate line .model nch nmos (level=54 ...). What does each line carry?",
		A:   "The instance carries connectivity and geometry; the model card carries the shared process parameters",
		D:   []string{"The instance carries the process parameters; the model card carries the node connections", "Both lines must repeat the full parameter set for every transistor", "The model card sets W and L for every device that names it, overriding the instance"},
		Why: "Instances say where a device sits and how big it is. Many instances share one .model card, so the simulator evaluates the same equations with per-instance geometry."},
	{KP: 1, Q: "Two hundred diodes in a deck name the same .model dmod d (is=1e-14 n=1.05). An engineer edits is on the model card. What changes?",
		A: "The saturation current of all two hundred diodes", D: []string{"Only the first diode that names dmod", "Nothing until each instance line is edited too", "Only diodes whose instance line also sets is"},
		Why: "A model card is shared state: every instance that names it is evaluated with its parameters, which is what makes process changes a one-line edit."},
	{KP: 1, Q: "The first letter of an element card is D, and the card is D1 a k dmod. Which nodes are which?",
		A: "a is the anode and k is the cathode", D: []string{"a is the cathode and k is the anode", "Both nodes are interchangeable for a diode", "k names the model and dmod is a node"},
		Why: "SPICE diode cards list the anode first, then the cathode, then the model name. Swapping them reverses the junction."},
	{KP: 1, Q: "A MOSFET card lists four nodes, as in M1 d g s b nch. Why four rather than three?",
		A: "The fourth is the body (bulk) terminal, which sets the body effect on the threshold", D: []string{"The fourth node is the model name", "The fourth is a sense node used only for output", "The fourth node is ground and can be any value"},
		Why: "SPICE MOSFETs are four-terminal devices: drain, gate, source and bulk. The source-to-bulk voltage shifts the threshold voltage."},
	{KP: 2, Q: "Two decks use the same transistor names but one model card says level=1 and the other names a BSIM4-class level. What does the level select?",
		A: "Which compact-model equations the simulator evaluates for that card", D: []string{"The number of Newton iterations allowed", "The metal layer the device is drawn on", "How many instances may share the card"},
		Why: "The level parameter picks the model implementation, from the Shichman-Hodges square law at level 1 to BSIM-class models, each with its own parameter set."},
	{KP: 2, Q: "A foundry ships model cards nch.1, nch.2 and nch.3 with lmin and lmax fields, and the deck instance names nch. What is the simulator doing?",
		A: "Binning: it picks the card whose length and width range contains the instance geometry", D: []string{"Averaging the three cards", "Running the instance three times and reporting the worst", "Using nch.1 and ignoring the others"},
		Why: "Binned models split the geometry space into ranges with separately fitted parameters. The simulator selects the bin from the instance W and L."},
	{KP: 2, Q: "A level-1 model card for an NMOS sets vto=0.5 kp=200u lambda=0.05. Which equation is the simulator using in saturation?",
		A: "The square law, drain current kp/2 times W/L times the overdrive squared times (1 + lambda Vds)", D: []string{"An exponential subthreshold law only", "A BSIM surface-potential model", "A linear resistor of value 1/kp"},
		Why: "Level 1 is the Shichman-Hodges square-law model with channel-length modulation through lambda."},
	{KP: 2, Q: "A model card from a newer process gives hundreds of parameters, and a designer asks why a square-law fit is not enough. What is the best answer?",
		A: "Short-channel effects such as velocity saturation, DIBL and mobility degradation bend the curves away from the square law", D: []string{"Square-law models cannot be stamped into MNA", "The square law has no derivative for Newton", "SPICE refuses any model card with fewer than a hundred parameters"},
		Why: "Compact models such as BSIM capture effects that dominate short-channel devices, which a square law with a few parameters cannot fit across bias and geometry."},
	{KP: 3, Q: "A deck says .lib models.lib tt and another run says .lib models.lib ss. What changes between the runs?",
		A: "The library section loaded: the typical or slow process parameter set", D: []string{"The simulation temperature only", "The supply voltage only", "The integration method"},
		Why: "Foundry libraries group parameter sets into named sections, one per process corner, and .lib with a section name selects one."},
	{KP: 3, Q: "A deck defines .subckt inv in out vdd vss with two MOSFETs and instantiates X1 a b vdd 0 inv and X2 b c vdd 0 inv. How many MOSFET instances does the flattened circuit have?",
		A: "Four", D: []string{"Two", "Six", "One per unique subcircuit, so two"},
		Why: "Each X line expands the subcircuit body with its own internal nodes and devices. Two instances of a two-transistor inverter give four MOSFETs."},
	{KP: 3, Q: "Inside .subckt amp, an internal node is named n1, and the deck has two amp instances X1 and X2. Why do they not short together?",
		A: "Internal nodes are local to each instance, so the flattener renames them per instance", D: []string{"SPICE forbids reusing node names", "Only the first instance keeps n1", "They do short together, which is a classic bug"},
		Why: "Subcircuit expansion prefixes internal names with the instance path, so X1's n1 and X2's n1 become distinct nodes."},
	{KP: 3, Q: "A Monte Carlo deck needs each transistor to receive its own random threshold shift, but every instance shares one model card. Where must the randomness enter?",
		A: "In per-instance parameters or mismatch statements evaluated separately for each device", D: []string{"In the shared model card value alone, which then gives each device its own draw", "In the .tran timestep", "In the node names of each instance"},
		Why: "A single random value on a shared card moves every device together, which models global variation. Local mismatch needs an independent draw per instance."},
}

var ngspiceXyceItems = []edakit.Item{
	{KP: 1, Q: "A deck contains .tran 1p 2n. What do the two numbers mean in SPICE?",
		A: "The suggested print step and the stop time; the internal step is still chosen by error control", D: []string{"The fixed integration step and the stop time", "The start time and the stop time", "The minimum and maximum internal step"},
		Why: "In .tran tstep tstop, tstep is the output (and initial suggested) step. The simulator adapts its own internal steps from truncation error estimates."},
	{KP: 1, Q: "Which analysis card makes the simulator linearize the circuit at its operating point and sweep frequency?",
		A: ".ac", D: []string{".tran", ".dc", ".op"},
		Why: ".ac computes the DC operating point, linearizes every device there and solves the complex system at each frequency."},
	{KP: 1, Q: "A user runs .dc V1 0 1.8 0.01 on an inverter. What is computed at each point?",
		A: "A full nonlinear DC operating point with V1 set to that value, starting from the previous point's solution", D: []string{"One transient step of 0.01 s", "A small-signal gain at 0.01 Hz", "A linear solve with devices frozen at the first point"},
		Why: "A DC sweep repeats the Newton operating-point solve per value, and using the previous solution as the initial guess is a natural continuation."},
	{KP: 1, Q: "A deck has .noise v(out) V1 dec 10 1k 1g. What does it report?",
		A: "Output noise spectral density at v(out) and the equivalent input noise referred to source V1", D: []string{"The transient waveform of random noise", "The DC operating point with noise sources removed", "Monte Carlo statistics of v(out)"},
		Why: "Noise analysis is a small-signal frequency sweep that sums each device's noise power through its transfer function to the output, and refers it to the named input."},
	{KP: 2, Q: "A transient run takes tiny steps and the user adds .options reltol=1e-2. What is the effect?",
		A: "Looser error and convergence tolerances: fewer, larger steps and lower accuracy", D: []string{"Tighter tolerances and smaller steps", "No effect, because reltol only applies to AC", "The integration method switches to Gear"},
		Why: "reltol scales the relative tolerances in both the Newton convergence test and timestep control. Raising it from the default 1e-3 trades accuracy for speed."},
	{KP: 2, Q: "An LC oscillator deck shows trapezoidal ringing on a sharp edge. Which ngspice option addresses the integration method directly?",
		A: ".options method=gear", D: []string{".options gmin=1e-15", ".options itl1=500", ".options vntol=1e-9"},
		Why: "method=gear switches from trapezoidal to backward differentiation formulas, which damp the spurious oscillation trapezoidal produces on stiff components."},
	{KP: 2, Q: "A deck fails the DC operating point with too many iterations. Which pair of options relates to that failure rather than to transient accuracy?",
		A: "itl1 (the DC iteration limit) and gmin (the minimum shunt conductance)", D: []string{"trtol and chgtol", "tstep and tstop", "method and maxord"},
		Why: "itl1 bounds Newton iterations for the operating point and gmin sets the conductance used by gmin stepping. trtol and chgtol govern transient step control."},
	{KP: 2, Q: "A user sets .options abstol=1e-6 on a circuit whose bias currents are nanoamperes. What risk does that create?",
		A: "Current convergence tests pass while nanoampere currents are still wrong", D: []string{"The run cannot converge at all", "Timesteps become smaller", "Voltages are reported in microvolts"},
		Why: "abstol is the absolute current tolerance, default 1 pA. Setting it to 1 uA makes errors far larger than the signal currents acceptable."},
	{KP: 3, Q: "A script runs ngspice in batch mode and needs every node voltage of a transient for later analysis. What should it ask for?",
		A: "A raw output file, for example with ngspice -b deck.cir -r out.raw", D: []string{"A screenshot of the plot window", "The .model cards printed to the log", "Only the final time point from the log"},
		Why: "The raw file stores all vectors (node voltages, branch currents, time) in a documented format that post-processing tools read directly."},
	{KP: 3, Q: "A characterization flow needs the 50% propagation delay of each run as one number. Which deck feature computes it inside the simulator?",
		A: "A .measure statement with trig and targ conditions on the two waveforms", D: []string{"An extra .op analysis", "A tighter reltol", "A .print of every time point"},
		Why: ".measure evaluates expressions on the waveforms after the run, such as the time between a 50% crossing of the input and of the output, and logs the value."},
	{KP: 3, Q: "A team must simulate a very large post-layout analog block on a cluster with distributed memory. Which property of Xyce matters most?",
		A: "It was designed for parallel execution, with distributed-memory solvers", D: []string{"It is the only simulator that reads SPICE decks", "It uses fixed timesteps, which parallelize trivially", "It skips the DC operating point"},
		Why: "Xyce was written at Sandia National Laboratories as a parallel circuit simulator, so it partitions the problem and uses distributed linear solvers."},
	{KP: 3, Q: "A deck that converges in ngspice fails the operating point in another simulator. What is the most useful first comparison?",
		A: "The default options of both tools, such as tolerances, gmin and the convergence aids each enables", D: []string{"The order of the cards in the deck", "The number of comment lines", "The plotting backend"},
		Why: "Different simulators ship different defaults and homotopy strategies. Matching tolerances and aids isolates whether the difference is numerical or a model difference."},
}
