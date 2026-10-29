package pcb

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var thtSmtItems = []edakit.Item{
	{KP: 1, Q: "How is a through-hole component attached to a board?",
		A:   "Its leads pass through plated holes and are soldered, usually by wave or selective soldering",
		D:   []string{"Its terminations sit on surface lands and are reflowed", "It is glued in place with no solder", "It is pressed into blind vias only"},
		Why: "Through-hole technology relies on plated barrels. The joint fills the hole and forms fillets on both sides."},
	{KP: 1, Q: "How is a surface-mount component attached to a board?",
		A:   "Its terminations sit on surface lands printed with solder paste and are joined in a reflow oven",
		D:   []string{"Its leads pass through plated holes and are wave soldered", "It is clipped into a socket", "Its leads are crimped to the traces"},
		Why: "SMT prints paste through a stencil, places parts on the wet paste and melts it in a controlled oven profile."},
	{KP: 1, Q: "What made surface mount the dominant assembly technology?",
		A:   "Smaller parts on both sides of the board, no hole per lead, and fast automated placement",
		D:   []string{"Stronger joints than through-hole under mechanical load", "Lower soldering temperatures", "No need for solder at all"},
		Why: "Removing a drilled hole per lead freed routing space on every layer and let machines place tens of thousands of parts per hour."},
	{KP: 1, Q: "Which parts are still commonly through-hole?",
		A:   "Large connectors, power components and parts that take mechanical stress",
		D:   []string{"0402 chip resistors", "Ball grid arrays", "QFN regulators"},
		Why: "A lead anchored through the board carries insertion and vibration loads that would crack a surface fillet."},
	{KP: 2, Q: "A board is all surface mount except one through-hole connector. Which process avoids a separate wave-solder step?",
		A:   "Pin-in-paste: print paste over the plated hole and reflow the connector with the SMT parts",
		D:   []string{"Wave soldering the whole board twice", "Reflowing the board twice with no paste in the hole", "Gluing the connector instead of soldering it"},
		Why: "Intrusive reflow fills the barrel from the printed paste, provided the part's body survives the reflow temperature."},
	{KP: 2, Q: "Why does a through-hole part cost routing space on every layer?",
		A:   "Each plated hole passes through all layers, with a pad or clearance on each one",
		D:   []string{"Through-hole pads exist only on the top layer", "Its holes are buried between inner layers", "Through-hole pads need no clearance"},
		Why: "The barrel and its antipads block every layer, unlike a surface land that occupies only the outer layer."},
	{KP: 2, Q: "A mixed-technology board sends bottom-side SMT parts over a wave-solder machine. What does that require?",
		A:   "Adhesive to hold the bottom SMT parts and a layout oriented and spaced for wave soldering",
		D:   []string{"Nothing beyond normal reflow rules", "Pin-in-paste for the bottom SMT parts", "Removing every bottom-side SMT part"},
		Why: "The parts are inverted through molten solder, so they are glued, and orientation and spacing rules prevent shadowing and bridges."},
	{KP: 2, Q: "Which surface-mount package hides its joints under the body, so inspection needs X-ray?",
		A: "A ball grid array", D: []string{"An SOIC", "An 0603 chip resistor", "A SOT-23"},
		Why: "BGA balls are entirely under the package, out of reach of optical inspection."},
	{KP: 3, Q: "After reflow, a small chip capacitor stands up on one end. What is this, and what causes it?",
		A:   "Tombstoning, from unequal wetting forces when one land wets or melts before the other",
		D:   []string{"Head-in-pillow, from BGA warpage", "Solder balling, from moisture in the paste", "Bridging, from excess paste"},
		Why: "Asymmetric pad size, trace thermal mass or paste volume makes one end's surface tension win and lift the part."},
	{KP: 3, Q: "Why is a 0.5 mm pitch QFP more prone to solder bridges than a 1.27 mm pitch SOIC?",
		A:   "The smaller gap between leads lets excess paste or slump join neighboring leads",
		D:   []string{"QFP leads are made of a different metal", "Bridges depend only on reflow peak temperature", "SOIC leads are always gold plated"},
		Why: "Bridging risk scales with paste volume relative to the gap, so fine pitch needs thinner stencils and tighter print control."},
	{KP: 3, Q: "Why do through-hole joints survive repeated connector insertion better than SMT joints?",
		A:   "The lead is anchored through the board, so force is carried by the barrel rather than a thin surface fillet",
		D:   []string{"Through-hole solder has a higher melting point", "SMT joints contain no tin", "SMT parts cannot be used for connectors"},
		Why: "Mechanical load on an SMT joint concentrates in the pad-to-laminate bond, which can lift the pad."},
	{KP: 3, Q: "Why are moisture-sensitive packages dried or kept in dry storage before reflow?",
		A:   "Absorbed moisture turns to steam during reflow and can crack or delaminate the package",
		D:   []string{"Moisture lowers the solder's melting point", "Moisture shorts the leads during placement", "Moisture changes the BOM"},
		Why: "The rapid vapor expansion inside a plastic package is the so-called popcorn failure; J-STD-020 classifies sensitivity levels."},
}

var dfaItems = []edakit.Item{
	{KP: 1, Q: "What are fiducials for?",
		A:   "Optical reference marks that let the placement machine correct for board position, rotation and stretch",
		D:   []string{"Probe points for in-circuit test", "Mounting holes for the enclosure", "Thermal vias under hot parts"},
		Why: "The machine's camera finds the fiducials and transforms the centroid coordinates to where the board actually sits."},
	{KP: 1, Q: "Why keep a component-free border along two board edges, or add breakaway rails?",
		A:   "Conveyors and clamps grip those edges during printing, placement and reflow",
		D:   []string{"To set the trace impedance", "To leave room for the silkscreen", "Because DRC forbids copper near edges"},
		Why: "Assembly equipment transports boards by their edges, so parts there collide with the rails."},
	{KP: 1, Q: "Why do through-hole pads on planes use thermal reliefs?",
		A:   "Spokes limit heat flow into the plane so the joint can reach soldering temperature",
		D:   []string{"To lower the pad's impedance", "To raise its current capacity", "To keep solder mask off the pad"},
		Why: "A solid connection to a large plane sinks heat, giving cold joints and hard rework."},
	{KP: 1, Q: "Why must polarity and pin 1 be marked clearly on silkscreen and assembly drawings?",
		A:   "So programmers, inspectors and rework operators orient diodes, capacitors and ICs correctly",
		D:   []string{"Polarity marks set the netlist order", "The fabricator needs them for drilling", "They are only decorative"},
		Why: "A rotated polarized part is a functional defect that every electrical check before assembly passes."},
	{KP: 2, Q: "Why enforce a gap between tall parts and fine-pitch parts?",
		A:   "Tall bodies shadow automated optical inspection and block rework tools",
		D:   []string{"Tall parts change the trace impedance", "The gap sets the stencil thickness", "Fine-pitch parts need the gap for silkscreen"},
		Why: "Inspection cameras view at an angle and rework nozzles need access, so spacing grows with height."},
	{KP: 2, Q: "Why are small boards panelized for assembly?",
		A:   "Several boards per panel go through printing, placement and reflow together, raising throughput and easing handling",
		D:   []string{"Panelization improves signal integrity", "It removes the need for fiducials", "IPC requires every board to be panelized"},
		Why: "A panel amortizes machine loading and gives conveyors a large enough object; boards are separated afterwards."},
	{KP: 2, Q: "Why keep open vias out of SMT lands unless they are filled and capped?",
		A:   "Molten solder wicks down the via barrel and starves the joint",
		D:   []string{"Vias in lands change the land's net", "The pick-and-place camera cannot see them", "It is a clearance violation by definition"},
		Why: "Capillary action pulls solder into an open hole, leaving too little to form the fillet."},
	{KP: 2, Q: "How does rotating polarized parts to a few consistent orientations help?",
		A:   "It reduces programming, inspection and manual-placement errors",
		D:   []string{"It shortens every trace", "It lowers the board's layer count", "It removes the need for a centroid file"},
		Why: "Consistent orientation makes the board easy to check by eye and by AOI, and prevents rotation mistakes in the program."},
	{KP: 3, Q: "Where the circuit tolerates it, replacing three close 0402 resistor values with one helps assembly how?",
		A:   "Fewer BOM lines mean fewer feeders, setups and kitting errors",
		D:   []string{"It lowers trace impedance", "It improves signal integrity", "It reduces the copper area"},
		Why: "Feeder slots and changeovers cost machine time; consolidating values is a standard DFA and cost move."},
	{KP: 3, Q: "A design places parts on both sides of the board. What does that cost in assembly?",
		A:   "A second print and reflow pass, extra handling, and possibly adhesive for heavy bottom parts",
		D:   []string{"Nothing extra, since both sides reflow at once", "Only a second set of fiducials", "A second bare board"},
		Why: "Double-sided SMT doubles the line passes. It buys area, so it is a deliberate trade."},
	{KP: 3, Q: "An assembly quote jumps because one part is a 0.4 mm pitch BGA. Why?",
		A:   "It may need a finer stencil, tighter placement, via-in-pad boards and X-ray inspection",
		D:   []string{"BGAs require through-hole soldering", "Fine-pitch parts need a heavier copper weight", "BGAs cannot be placed by machine"},
		Why: "Fine pitch pushes both the board fabrication and the assembly process to tighter tolerances and adds inspection."},
	{KP: 3, Q: "Why should DFA review happen before layout is finished rather than after?",
		A:   "Fixing spacing, orientation and side assignments after routing means rerouting, so problems are cheapest early",
		D:   []string{"DFA rules do not depend on placement", "Assembly houses refuse finished layouts", "DFA can only be checked on Gerber files"},
		Why: "Assembly constraints are placement constraints. Late changes ripple into routing and signal integrity."},
}

// pickPlace: machine time (KP 1), line balance (KP 2), centroid rotation (KP 3).
func pickPlace(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		p := edakit.Int(10, 150) * 10
		b := edakit.Pick([]int{10, 25, 50, 100, 200, 500})
		cph := edakit.Int(8, 40) * 1000
		h := float64(p*b) / float64(cph)
		ans := edakit.Num(h, 2)
		q := edakit.Sprintf("A board has %d SMT placements. A build of %d boards runs on a placement machine that sustains %d components per hour (CPH) on this job. How many hours of placement time does the build need? Round to 2 decimal places.", p, b, cph)
		return edakit.Numeric(q, ans, edakit.Sprintf("$%d \\cdot %d / %d = %s$ hours. Rated CPH is measured on ideal parts; real jobs with large or fine-pitch parts run slower.", p, b, cph, ans), nil)
	case 2:
		p1 := edakit.Int(5, 60) * 10
		p2 := edakit.Int(5, 60) * 10
		c1 := edakit.Int(15, 40) * 1000
		c2 := edakit.Int(5, 20) * 1000
		t1 := float64(p1) / float64(c1) * 3600
		t2 := float64(p2) / float64(c2) * 3600
		t := t1
		if t2 > t {
			t = t2
		}
		ans := edakit.Num(t, 1)
		q := edakit.Sprintf("A line has a chip shooter at %d CPH placing %d parts per board, followed by a fine-pitch placer at %d CPH placing %d parts per board. Boards flow through both in sequence. Ignoring transfer time, what is the line's cycle time per board, in seconds? Round to 1 decimal place.", c1, p1, c2, p2)
		return edakit.Numeric(q, ans, edakit.Sprintf("Machine 1 needs %s s and machine 2 needs %s s per board. The slower machine sets the pace: %s s. Moving placements to balance the two shortens the cycle.", edakit.Num(t1, 1), edakit.Num(t2, 1), ans), map[string]string{"t1": edakit.Num(t1, 1), "t2": edakit.Num(t2, 1)})
	}
	r := edakit.Pick([]int{0, 90, 180, 270})
	o := edakit.Pick([]int{-90, 90, 180, 270, -180})
	m := ((r+o)%360 + 360) % 360
	q := edakit.Sprintf("A centroid file lists a top-side part at %d degrees. The machine's feeder library holds this package at a zero orientation that differs from the CAD library's by %d degrees, so the machine angle is (CAD angle + %d) mod 360, as a value from 0 to 359. What angle does the machine program use, in degrees?", r, o, o)
	return edakit.Numeric(q, edakit.IntStr(m), edakit.Sprintf("$(%d + %d) \\bmod 360 = %d$. A wrong offset in this table rotates every instance of the package, which is why first-article inspection checks polarity.", r, o, m), nil)
}

// reflow draws a piecewise-linear profile and asks for a ramp rate (KP 1), the time above
// liquidus (KP 2) or the time from 150 C to 200 C on the way up (KP 3).
func reflow(ctx generator.GeneratorContext) generator.Problem {
	t1 := float64(edakit.Int(60, 100))
	T2 := float64(edakit.Int(180, 195))
	t2 := t1 + float64(edakit.Int(60, 100))
	Tp := float64(edakit.Int(235, 250))
	t3 := t2 + float64(edakit.Int(30, 60))
	T4 := float64(edakit.Int(195, 210))
	t4 := t3 + float64(edakit.Int(20, 45))
	t5 := t4 + float64(edakit.Int(40, 80))
	prof := edakit.Sprintf("A reflow profile is piecewise linear through (0 s, 25 C), (%s s, 150 C), (%s s, %s C), (%s s, %s C peak), (%s s, %s C) and (%s s, 100 C).", edakit.Num(t1, 0), edakit.Num(t2, 0), edakit.Num(T2, 0), edakit.Num(t3, 0), edakit.Num(Tp, 0), edakit.Num(t4, 0), edakit.Num(T4, 0), edakit.Num(t5, 0))
	cross := func(ta, Ta, tb, Tb, T float64) float64 { return ta + (T-Ta)/(Tb-Ta)*(tb-ta) }
	switch kpOf(ctx) {
	case 1:
		switch edakit.Int(1, 3) {
		case 1:
			rate := 125 / t1
			ans := edakit.Num(rate, 2)
			return edakit.Numeric(prof+" What is the preheat ramp rate from 25 C to 150 C, in C per second? Round to 2 decimal places.", ans, edakit.Sprintf("$(150 - 25)/%s = %s$ C per second. Ramps much above a few C per second risk thermal shock to ceramic capacitors.", edakit.Num(t1, 0), ans), nil)
		case 2:
			rate := (Tp - T2) / (t3 - t2)
			ans := edakit.Num(rate, 2)
			return edakit.Numeric(prof+edakit.Sprintf(" What is the ramp rate from the end of soak (%s C) to the peak, in C per second? Round to 2 decimal places.", edakit.Num(T2, 0)), ans, edakit.Sprintf("$(%s - %s)/(%s - %s) = %s$ C per second.", edakit.Num(Tp, 0), edakit.Num(T2, 0), edakit.Num(t3, 0), edakit.Num(t2, 0), ans), nil)
		}
		rate := (Tp - T4) / (t4 - t3)
		ans := edakit.Num(rate, 2)
		return edakit.Numeric(prof+" What is the cooling rate just after the peak, as a positive number in C per second? Round to 2 decimal places.", ans, edakit.Sprintf("$(%s - %s)/(%s - %s) = %s$ C per second of cooling.", edakit.Num(Tp, 0), edakit.Num(T4, 0), edakit.Num(t4, 0), edakit.Num(t3, 0), ans), nil)
	case 2:
		up := cross(t2, T2, t3, Tp, 217)
		down := cross(t3, Tp, t4, T4, 217)
		ans := edakit.Num(down-up, 1)
		return edakit.Numeric(prof+" Take the liquidus of the SAC305 solder as 217 C. How long is the joint above liquidus, in seconds? Round to 1 decimal place.", ans,
			edakit.Sprintf("The profile crosses 217 C going up at %s s and coming down at %s s, by linear interpolation on those segments. Time above liquidus is %s s; too short gives poor wetting, too long grows brittle intermetallics.", edakit.Num(up, 2), edakit.Num(down, 2), ans), map[string]string{"up": edakit.Num(up, 2), "down": edakit.Num(down, 2)})
	}
	t200 := cross(t2, T2, t3, Tp, 200)
	ans := edakit.Num(t200-t1, 1)
	return edakit.Numeric(prof+" How long does the board take to go from 150 C to 200 C on the way up, in seconds? Round to 1 decimal place.", ans,
		edakit.Sprintf("150 C is reached at %s s. 200 C lies on the segment toward the peak and is reached at %s s. The difference is %s s, the soak in which flux activates and the board's temperature evens out.", edakit.Num(t1, 0), edakit.Num(t200, 2), ans), map[string]string{"t200": edakit.Num(t200, 2)})
}
