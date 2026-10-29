package pcb

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const milPerOz = 1.37 // copper thickness per ounce per square foot, in mil

var ozChoices = []float64{0.5, 1, 1, 2}

func ozText(oz float64) string { return edakit.Num(oz, 1) + " oz" }

// stackLayers: finished thickness (KP 1), core and prepreg counts (KP 2), symmetry check (KP 3).
func stackLayers(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		L := edakit.Pick([]int{4, 6, 8})
		var parts []string
		cu, diel := 0.0, 0.0
		for i := 1; i <= L; i++ {
			oz := 1.0
			if i == 1 || i == L {
				oz = edakit.Pick([]float64{1, 1.5, 2})
			} else {
				oz = edakit.Pick(ozChoices)
			}
			cu += oz
			parts = append(parts, edakit.Sprintf("L%d copper %s", i, ozText(oz)))
			if i < L {
				d := float64(edakit.Int(6, 50)) / 2
				diel += d
				parts = append(parts, edakit.Sprintf("dielectric %s mil", edakit.Num(d, 1)))
			}
		}
		total := diel + cu*milPerOz
		ans := edakit.Num(total, 2)
		q := edakit.Sprintf("%s %d-layer stackup from top to bottom is: %s. Take 1 oz copper as %s mil thick and ignore solder mask. What is the finished board thickness, in mil? Round to 2 decimal places.", strings.ToUpper(an(L)[:1])+an(L)[1:], L, strings.Join(parts, ", "), edakit.Num(milPerOz, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Dielectrics sum to %s mil. Copper totals %s oz, which is %s x %s = %s mil. The board is %s mil thick.", edakit.Num(diel, 1), edakit.Num(cu, 1), edakit.Num(cu, 1), edakit.Num(milPerOz, 2), edakit.Num(cu*milPerOz, 2), ans),
			map[string]string{"dielectric_mil": edakit.Num(diel, 1), "copper_oz": edakit.Num(cu, 1)})
	case 2:
		L := 2 * edakit.Scaled(ctx, 2, 10)
		cores := (L - 2) / 2
		prepreg := L / 2
		if edakit.Int(1, 2) == 1 {
			q := edakit.Sprintf("%s %d-layer board is built by foil lamination: every inner copper pair sits on a core (copper on both sides of a cured laminate), the cores are stacked with prepreg between them, and the two outer layers are copper foils pressed onto prepreg. How many cores does the build use?", strings.ToUpper(an(L)[:1])+an(L)[1:], L)
			return edakit.Numeric(q, edakit.IntStr(cores),
				edakit.Sprintf("The %d inner layers pair up on cores, so there are (%d - 2)/2 = %d cores.", L-2, L, cores), map[string]string{"cores": edakit.IntStr(cores)})
		}
		q := edakit.Sprintf("%s %d-layer board is built by foil lamination: every inner copper pair sits on a core, the cores are stacked with prepreg between them, and the two outer layers are copper foils pressed onto prepreg. Counting each prepreg-filled gap between adjacent copper layers once, how many of the dielectric layers are prepreg?", strings.ToUpper(an(L)[:1])+an(L)[1:], L)
		return edakit.Numeric(q, edakit.IntStr(prepreg),
			edakit.Sprintf("%s %d-layer board has %d dielectric layers. Cores account for %d of them, so %d are prepreg.", strings.ToUpper(an(L)[:1])+an(L)[1:], L, L-1, cores, prepreg), map[string]string{"cores": edakit.IntStr(cores)})
	}
	L := edakit.Pick([]int{6, 8, 10})
	cu := make([]float64, L)
	dl := make([]float64, L-1)
	for i := 0; i < L/2; i++ {
		oz := edakit.Pick(ozChoices)
		if i == 0 {
			oz = edakit.Pick([]float64{1, 2})
		}
		cu[i], cu[L-1-i] = oz, oz
	}
	for j := 0; j < L/2; j++ {
		d := float64(edakit.Int(3, 24))
		dl[j], dl[L-2-j] = d, d
	}
	flips := edakit.Int(0, 3)
	for f := 0; f < flips; f++ {
		if edakit.Int(1, 2) == 1 {
			i := edakit.Int(0, L-1)
			cu[i] = edakit.Pick([]float64{0.5, 1, 2, 3})
		} else {
			j := edakit.Int(0, L-2)
			dl[j] = float64(edakit.Int(3, 24))
		}
	}
	bad := 0
	for i := 0; i < L/2; i++ {
		if cu[i] != cu[L-1-i] {
			bad++
		}
	}
	for j := 0; j < L/2-1; j++ {
		if dl[j] != dl[L-2-j] {
			bad++
		}
	}
	var parts []string
	for i := 0; i < L; i++ {
		parts = append(parts, edakit.Sprintf("L%d %s", i+1, ozText(cu[i])))
		if i < L-1 {
			parts = append(parts, edakit.Sprintf("%s mil", edakit.Num(dl[i], 0)))
		}
	}
	q := edakit.Sprintf("%s %d-layer stackup, top to bottom, is: %s. A symmetry check compares each copper layer with its mirror (L1 with L%d, L2 with L%d, and so on) and each dielectric with its mirror about the board's center. How many mirrored pairs (copper pairs plus dielectric pairs) differ?", strings.ToUpper(an(L)[:1])+an(L)[1:], L, strings.Join(parts, ", "), L, L-1)
	return edakit.Numeric(q, edakit.IntStr(bad),
		edakit.Sprintf("There are %d copper pairs and %d dielectric pairs (the center dielectric mirrors onto itself). %d of them differ. An unbalanced stackup bows and twists in lamination and reflow.", L/2, L/2-1, bad),
		map[string]string{"mismatched_pairs": edakit.IntStr(bad)})
}

// copperResistance: R = rho L / (w t) with t from copper weight; drop at a current; hot resistance.
func copperResistance(ctx generator.GeneratorContext) generator.Problem {
	oz := edakit.Pick([]float64{0.5, 1, 2})
	w := edakit.Pick([]float64{0.1, 0.15, 0.2, 0.25, 0.3, 0.5, 1.0})
	L := float64(10 * edakit.Int(2, 30))
	t := 35e-6 * oz
	r := 1.72e-8 * (L / 1000) / ((w / 1000) * t) // ohm
	base := edakit.Sprintf("A trace is %s mm long and %s mm wide in %s copper. Use $\\rho = 1.72 \\times 10^{-8}\\ \\Omega\\cdot$m and take 1 oz copper as 35 micrometres thick, so $R = \\rho L / (w t)$.", edakit.Num(L, 0), edakit.Num(w, 2), ozText(oz))
	facts := map[string]string{"R_mohm": edakit.Num(r*1000, 2)}
	switch kpOf(ctx) {
	case 1:
		ans := edakit.Num(r*1000, 1)
		return edakit.Numeric(base+" What is its resistance, in milliohms? Round to 1 decimal place.", ans,
			edakit.Sprintf("Thickness is %s micrometres. $R = 1.72 \\times 10^{-8} \\cdot %s / (%s \\cdot %s) = %s$ milliohms, with every length in metres.", edakit.Num(t*1e6, 1), edakit.Num(L/1000, 3), edakit.Num(w/1000, 5), edakit.Num(t, 7), ans), facts)
	case 2:
		I := float64(edakit.Int(1, 10)) / 2
		ans := edakit.Num(r*I*1000, 1)
		return edakit.Numeric(base+edakit.Sprintf(" It carries %s A. What is the voltage drop along it, in millivolts? Round to 1 decimal place.", edakit.Num(I, 1)), ans,
			edakit.Sprintf("$R = %s$ milliohms, so the drop is $I R = %s \\cdot %s = %s$ mV.", edakit.Num(r*1000, 2), edakit.Num(I, 1), edakit.Num(r*1000, 2), ans), facts)
	}
	T := float64(edakit.Pick([]int{40, 50, 60, 70, 85, 100, 105}))
	rt := r * (1 + 0.00393*(T-20))
	ans := edakit.Num(rt*1000, 1)
	return edakit.Numeric(base+edakit.Sprintf(" That gives the resistance at 20 C. Copper's temperature coefficient is 0.00393 per C, so $R_T = R_{20}(1 + 0.00393 (T - 20))$. What is its resistance at %s C, in milliohms? Round to 1 decimal place.", edakit.Num(T, 0)), ans,
		edakit.Sprintf("$R_{20} = %s$ milliohms. Multiplying by $1 + 0.00393 \\cdot %s = %s$ gives %s milliohms.", edakit.Num(r*1000, 2), edakit.Num(T-20, 0), edakit.Num(1+0.00393*(T-20), 4), ans), facts)
}

// currentCapacity uses the IPC-2221 chart fit I = k dT^0.44 A^0.725.
func currentCapacity(ctx generator.GeneratorContext) generator.Problem {
	oz := edakit.Pick([]float64{0.5, 1, 2})
	t := oz * milPerOz
	dT := float64(edakit.Pick([]int{10, 20, 30}))
	rule := edakit.Sprintf("Use the IPC-2221 rule of thumb $I = k\\, \\Delta T^{0.44} A^{0.725}$, with $I$ in amperes, temperature rise $\\Delta T$ in C, cross-section $A$ in square mils, $k = 0.048$ for outer layers and $k = 0.024$ for inner layers, and 1 oz copper %s mil thick.", edakit.Num(milPerOz, 2))
	kp := kpOf(ctx)
	if kp == 1 {
		w := float64(edakit.Pick([]int{5, 6, 8, 10, 12, 15, 20, 25, 30, 40, 50}))
		A := w * t
		I := 0.048 * pow(dT, 0.44) * pow(A, 0.725)
		ans := edakit.Num(I, 2)
		q := rule + edakit.Sprintf(" An outer-layer trace is %s mil wide in %s copper. What current gives a %s C rise? Answer in amperes, rounded to 2 decimal places.", edakit.Num(w, 0), ozText(oz), edakit.Num(dT, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("$A = %s \\cdot %s = %s$ square mils. $I = 0.048 \\cdot %s^{0.44} \\cdot %s^{0.725} = %s$ A.", edakit.Num(w, 0), edakit.Num(t, 3), edakit.Num(A, 2), edakit.Num(dT, 0), edakit.Num(A, 2), ans), map[string]string{"A_sqmil": edakit.Num(A, 2)})
	}
	k, where := 0.048, "outer"
	if kp == 3 {
		k, where = 0.024, "inner"
	}
	I := float64(edakit.Int(2, 16)) / 2
	A := pow(I/(k*pow(dT, 0.44)), 1/0.725)
	w := A / t
	ans := edakit.Num(w, 1)
	q := rule + edakit.Sprintf(" An %s-layer trace in %s copper must carry %s A with at most a %s C rise. What minimum width does it need, in mil? Round to 1 decimal place.", where, ozText(oz), edakit.Num(I, 1), edakit.Num(dT, 0))
	return edakit.Numeric(q, ans, edakit.Sprintf("Invert the rule: $A = (I / (k \\Delta T^{0.44}))^{1/0.725} = %s$ square mils with $k = %s$. Dividing by the %s mil copper thickness gives %s mil.", edakit.Num(A, 2), edakit.Num(k, 3), edakit.Num(t, 3), ans), map[string]string{"A_sqmil": edakit.Num(A, 2)})
}

func microstripZ(er, h, w, t float64) float64 { return 87 / sqrt(er+1.41) * ln(5.98*h/(0.8*w+t)) }
func striplineZ(er, b, w, t float64) float64  { return 60 / sqrt(er) * ln(1.9*b/(0.8*w+t)) }

var dkChoices = []float64{3.5, 3.7, 3.8, 4.0, 4.2, 4.4, 4.6}

// controlledImpedance: microstrip (KP 1), stripline (KP 2), sensitivity to a stackup change (KP 3).
func controlledImpedance(ctx generator.GeneratorContext) generator.Problem {
	msForm := "Use the IPC-2141 microstrip approximation $Z_0 = \\frac{87}{\\sqrt{\\varepsilon_r + 1.41}} \\ln\\frac{5.98 h}{0.8 w + t}$, where $h$ is the dielectric height to the plane, $w$ the trace width and $t$ the copper thickness."
	switch kpOf(ctx) {
	case 1:
		er := edakit.Pick(dkChoices)
		h := float64(edakit.Int(6, 16)) / 2
		w := edakit.Round(h*float64(edakit.Int(8, 19))/10, 1)
		t := edakit.Pick([]float64{0.7, 1.4, 2.1})
		z := microstripZ(er, h, w, t)
		ans := edakit.Num(z, 1)
		q := msForm + edakit.Sprintf(" A surface trace has $w = %s$ mil, $t = %s$ mil, $h = %s$ mil and $\\varepsilon_r = %s$. What is $Z_0$, in ohms? Round to 1 decimal place.", edakit.Num(w, 1), edakit.Num(t, 1), edakit.Num(h, 1), edakit.Num(er, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("$87/\\sqrt{%s} = %s$ and $\\ln(%s / %s) = %s$, so $Z_0 = %s$ ohms.", edakit.Num(er+1.41, 2), edakit.Num(87/sqrt(er+1.41), 3), edakit.Num(5.98*h, 2), edakit.Num(0.8*w+t, 2), edakit.Num(ln(5.98*h/(0.8*w+t)), 4), ans), map[string]string{"Z0": ans})
	case 2:
		er := edakit.Pick(dkChoices)
		b := float64(edakit.Int(10, 30))
		w := edakit.Round(b*float64(edakit.Int(15, 35))/100, 1)
		t := edakit.Pick([]float64{0.7, 1.4})
		z := striplineZ(er, b, w, t)
		ans := edakit.Num(z, 1)
		q := edakit.Sprintf("Use the IPC-2141 centered-stripline approximation $Z_0 = \\frac{60}{\\sqrt{\\varepsilon_r}} \\ln\\frac{1.9 b}{0.8 w + t}$, where $b$ is the plane-to-plane spacing. An inner trace has $w = %s$ mil, $t = %s$ mil, $b = %s$ mil and $\\varepsilon_r = %s$. What is $Z_0$, in ohms? Round to 1 decimal place.", edakit.Num(w, 1), edakit.Num(t, 1), edakit.Num(b, 0), edakit.Num(er, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("$60/\\sqrt{%s} = %s$ and $\\ln(%s / %s) = %s$, so $Z_0 = %s$ ohms.", edakit.Num(er, 1), edakit.Num(60/sqrt(er), 3), edakit.Num(1.9*b, 2), edakit.Num(0.8*w+t, 2), edakit.Num(ln(1.9*b/(0.8*w+t)), 4), ans), map[string]string{"Z0": ans})
	}
	er := edakit.Pick(dkChoices)
	h := float64(edakit.Int(8, 14)) / 2
	w := edakit.Round(h*float64(edakit.Int(10, 17))/10, 1)
	t := 1.4
	z1 := microstripZ(er, h, w, t)
	var change string
	h2, w2, er2 := h, w, er
	switch edakit.Int(1, 3) {
	case 1:
		h2 = h + edakit.Pick([]float64{-1, -0.5, 0.5, 1, 1.5})
		change = edakit.Sprintf("the fabricator's prepreg makes the height %s mil instead of %s mil", edakit.Num(h2, 1), edakit.Num(h, 1))
	case 2:
		w2 = w + edakit.Pick([]float64{-1.5, -1, -0.5, 0.5, 1, 2})
		change = edakit.Sprintf("the trace is drawn %s mil wide instead of %s mil", edakit.Num(w2, 1), edakit.Num(w, 1))
	default:
		er2 = er + 0.3
		for a := 0; a < generator.MaxDrawAttempts; a++ {
			if c := edakit.Pick(dkChoices); c != er {
				er2 = c
				break
			}
		}
		change = edakit.Sprintf("the laminate's $\\varepsilon_r$ is %s instead of %s", edakit.Num(er2, 1), edakit.Num(er, 1))
	}
	z2 := microstripZ(er2, h2, w2, t)
	d := edakit.Round(z2, 1) - edakit.Round(z1, 1)
	ans := edakit.Num(d, 1)
	q := msForm + edakit.Sprintf(" A trace designed with $w = %s$ mil, $t = 1.4$ mil, $h = %s$ mil and $\\varepsilon_r = %s$ is built, but %s. By how many ohms does $Z_0$ change (new minus old, negative if it falls)? Round each impedance to 1 decimal place before subtracting.", edakit.Num(w, 1), edakit.Num(h, 1), edakit.Num(er, 1), change)
	return edakit.Numeric(q, ans, edakit.Sprintf("Old $Z_0 = %s$ ohms, new $Z_0 = %s$ ohms, change %s ohms. Impedance rises with height and falls with width and with dielectric constant.", edakit.Num(z1, 1), edakit.Num(z2, 1), ans), map[string]string{"Z_old": edakit.Num(z1, 1), "Z_new": edakit.Num(z2, 1)})
}

// propDelay: stripline (KP 1), microstrip with an effective Dk (KP 2), a mixed route (KP 3).
func propDelay(ctx generator.GeneratorContext) generator.Problem {
	cText := "Light travels 11.8 inches per nanosecond in vacuum, so a line's delay is $t_{pd} = \\sqrt{\\varepsilon_{eff}} / 11.8$ ns per inch."
	er := edakit.Pick([]float64{3.0, 3.3, 3.5, 3.7, 3.9, 4.0, 4.2, 4.4, 4.6})
	switch kpOf(ctx) {
	case 1:
		v := 1000 * sqrt(er) / 11.8
		ans := edakit.Num(v, 1)
		q := cText + edakit.Sprintf(" A stripline is fully embedded in a dielectric with $\\varepsilon_r = %s$, so $\\varepsilon_{eff} = \\varepsilon_r$. What is its delay, in ps per inch? Round to 1 decimal place.", edakit.Num(er, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\sqrt{%s} = %s$, and $1000 \\cdot %s / 11.8 = %s$ ps per inch.", edakit.Num(er, 1), edakit.Num(sqrt(er), 4), edakit.Num(sqrt(er), 4), ans), map[string]string{"tpd": ans})
	case 2:
		ee := 0.475*er + 0.67
		v := 1000 * sqrt(ee) / 11.8
		ans := edakit.Num(v, 1)
		q := cText + edakit.Sprintf(" A surface microstrip has part of its field in air. Approximate $\\varepsilon_{eff} = 0.475\\,\\varepsilon_r + 0.67$ for a laminate with $\\varepsilon_r = %s$. What is its delay, in ps per inch? Round to 1 decimal place.", edakit.Num(er, 1))
		return edakit.Numeric(q, ans, edakit.Sprintf("$\\varepsilon_{eff} = %s$, so the delay is $1000\\sqrt{%s}/11.8 = %s$ ps per inch, less than the %s ps per inch of a stripline in the same laminate.", edakit.Num(ee, 4), edakit.Num(ee, 4), ans, edakit.Num(1000*sqrt(er)/11.8, 1)), map[string]string{"eeff": edakit.Num(ee, 4)})
	}
	ls := float64(edakit.Int(5, 60)) / 10
	lm := float64(edakit.Int(3, 20)) / 10
	ee := 0.475*er + 0.67
	ts := 1000 * sqrt(er) / 11.8 * ls
	tm := 1000 * sqrt(ee) / 11.8 * lm
	ans := edakit.Num(ts+tm, 1)
	q := cText + edakit.Sprintf(" A net runs %s in of surface microstrip (approximate $\\varepsilon_{eff} = 0.475\\,\\varepsilon_r + 0.67$) and %s in of stripline ($\\varepsilon_{eff} = \\varepsilon_r$), in a laminate with $\\varepsilon_r = %s$. Ignoring the vias, what is the total delay, in ps? Round to 1 decimal place.", edakit.Num(lm, 1), edakit.Num(ls, 1), edakit.Num(er, 1))
	return edakit.Numeric(q, ans, edakit.Sprintf("Microstrip: %s ps. Stripline: %s ps. Total %s ps. Equal lengths on different layers are not equal delays, which is why matching is done in time.", edakit.Num(tm, 2), edakit.Num(ts, 2), ans), map[string]string{"t_micro": edakit.Num(tm, 2), "t_strip": edakit.Num(ts, 2)})
}

// bendRadius is the computed KP 2 of flex_rigid_flex.
func bendRadius(ctx generator.GeneratorContext) generator.Problem {
	type rule struct {
		name string
		f    float64
	}
	rules := []rule{{"one copper layer", 6}, {"two copper layers", 12}, {"three or more copper layers", 24}}
	r := edakit.Pick(rules)
	thick := float64(edakit.Int(5, 40)) / 100
	dynamic := edakit.Int(1, 3) == 1
	f := r.f
	kind := "a static bend (flexed only at installation)"
	if dynamic {
		f = 100
		kind = "a dynamic bend (flexed repeatedly in service)"
	}
	correct := f * thick
	wrong := []float64{}
	for _, o := range rules {
		if o.f != f {
			wrong = append(wrong, o.f*thick)
		}
	}
	if !dynamic {
		wrong = append(wrong, 100*thick)
	}
	wrong = append(wrong, f*thick/2)
	ans, ds := options(correct, 2, "mm", wrong...)
	q := edakit.Sprintf("A fabricator's rule of thumb sets the minimum bend radius at 6 times the flex thickness for one copper layer, 12 times for two, and 24 times for three or more, for static bends, and at 100 times the thickness for any dynamic bend. A flex section with %s is %s mm thick and will see %s. What is the minimum bend radius?", r.name, edakit.Num(thick, 2), kind)
	return edakit.Choice(q, ans, edakit.Sprintf("The applicable factor is %s, so the radius is %s x %s mm = %s. Bending tighter strains the outer copper past its fatigue limit.", edakit.Num(f, 0), edakit.Num(f, 0), edakit.Num(thick, 2), ans), ds...)
}

var flexItems = []edakit.Item{
	{KP: 1, Q: "What is the usual base dielectric of a flexible circuit?",
		A: "Polyimide film", D: []string{"FR-4 glass-epoxy laminate", "Alumina ceramic", "Glass-reinforced PTFE"},
		Why: "Polyimide survives soldering temperatures and bends without cracking, which is why it is the standard flex substrate."},
	{KP: 1, Q: "In a rigid-flex board, what forms the flexible region?",
		A:   "Flexible polyimide layers that run out of the rigid sections with no rigid material laminated over them",
		D:   []string{"A cable soldered between two rigid boards", "FR-4 milled down to a thinner section", "Copper wires molded into silicone"},
		Why: "The flex layers are continuous through the rigid sections, where they are laminated with rigid material, and left bare between them."},
	{KP: 1, Q: "What protects flex copper in place of a liquid solder mask?",
		A:   "A coverlay: polyimide film with adhesive laminated over the circuit",
		D:   []string{"A thicker silkscreen layer", "Conformal coating applied only after assembly", "A sheet of prepreg"},
		Why: "Liquid masks crack when flexed. A laminated polyimide coverlay bends with the circuit, with openings cut where pads must be exposed."},
	{KP: 1, Q: "Why is rolled-annealed copper preferred over electrodeposited copper for dynamic flex?",
		A:   "Its elongated grain structure survives many more bend cycles before cracking",
		D:   []string{"Its resistivity is half that of electrodeposited copper", "It is always the cheaper foil", "It cannot oxidize"},
		Why: "Rolling produces long horizontal grains that resist fatigue cracking, which matters when a hinge flexes thousands of times."},
	{KP: 3, Q: "How should traces cross a bend line?",
		A:   "Perpendicular to the bend, staggered between layers rather than stacked",
		D:   []string{"Parallel to the bend line", "Stacked directly above each other on adjacent layers", "At 45 degrees with sharp corners"},
		Why: "Perpendicular traces bend uniformly. Stacked traces form an I-beam that stiffens the section and concentrates strain."},
	{KP: 3, Q: "Why keep vias and pads out of the bend area?",
		A:   "Plated holes and pads concentrate strain and crack when the flex bends",
		D:   []string{"Vias raise the impedance of the bend", "Fabricators cannot drill polyimide", "Vias block the coverlay adhesive"},
		Why: "A plated barrel is rigid. Placed in a bend it becomes a stress riser, so vias stay in the rigid or stiffened regions."},
	{KP: 3, Q: "Why are flex planes often cross-hatched rather than solid?",
		A:   "Hatching keeps the section flexible, at the cost of changing the impedance of traces referenced to it",
		D:   []string{"IPC forbids solid planes on flex", "Hatching lowers the plane's DC resistance", "Hatching raises the dielectric constant"},
		Why: "Solid copper stiffens a flex. A hatched reference plane bends better, but impedance must then be modeled with the hatch."},
	{KP: 3, Q: "Why does rigid-flex often replace a connector-and-cable assembly?",
		A:   "It removes connectors and their joints, saving space and failure points, at a higher board cost",
		D:   []string{"It is always cheaper per board", "It has no limit on layer count", "It needs no stackup design"},
		Why: "Connectors and cables are a common field failure. Rigid-flex trades a more expensive fabrication for fewer interconnects."},
	{KP: 3, Q: "What does a stiffener do on a flex circuit?",
		A:   "Locally supports the area under connectors or components so their solder joints are not flexed",
		D:   []string{"Adds a ground reference plane", "Raises the allowed bend radius", "Replaces the coverlay"},
		Why: "A stiffener, often FR-4 or polyimide bonded on the back, makes a local rigid zone for assembly and insertion forces."},
}
