package extract

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// rcCorners counts signoff scenarios, computes corner parasitics and compares corners.
func rcCorners(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		modes := edakit.Int(2, 5)
		ps, rs := edakit.Int(2, 5), edakit.Int(2, 3)
		ph, rh := edakit.Int(2, 5), edakit.Int(2, 3)
		tot := modes * (ps*rs + ph*rh)
		return num(edakit.Sprintf("A chip has %d functional modes. Setup is signed off at %d PVT corners, each with %d RC corners; hold is signed off at %d PVT corners, each with %d RC corners. Every mode is checked at every corner combination for its check type. How many timing scenarios must be run?", modes, ps, rs, ph, rh),
			float64(tot), 0, edakit.Sprintf("Setup: %d x %d x %d = %d. Hold: %d x %d x %d = %d. Total %d. Multi-mode multi-corner analysis multiplies quickly, which is why extraction runs once per RC corner and its results are shared across the PVT corners that use it.", modes, ps, rs, modes*ps*rs, modes, ph, rh, modes*ph*rh, tot), nil)
	case 2:
		dw := edakit.Pick([]float64{5, 8, 10})
		dt := edakit.Pick([]float64{5, 8, 10})
		r0 := float64(10 * edakit.Int(5, 50))
		r := r0 / ((1 - dw/100) * (1 - dt/100))
		return num(edakit.Sprintf("In the RCworst corner the metal width is %s%% below nominal and the thickness is %s%% below nominal. A net has %s ohm of wire resistance at nominal. What is its resistance in the RCworst corner, in ohms?", f(dw, 0), f(dt, 0), f(r0, 0)),
			r, 2, edakit.Sprintf("Resistance goes as $1/(WT)$: %s / (%s x %s) = %s ohm. Narrower, thinner metal raises resistance but lowers sidewall coupling, so RCworst is not the worst capacitance corner.", f(r0, 0), f(1-dw/100, 2), f(1-dt/100, 2), f(r, 2)), nil)
	}
	rd := float64(edakit.Int(5, 40)) / 10
	rw := float64(edakit.Int(2, 30)) / 10
	cw := float64(edakit.Int(10, 80))
	cl := float64(edakit.Int(2, 20))
	cR, cC := edakit.Pick([]float64{0.9, 0.92, 0.95}), edakit.Pick([]float64{1.08, 1.1, 1.12})
	rR, rC := edakit.Pick([]float64{1.12, 1.15, 1.2}), edakit.Pick([]float64{0.96, 0.98, 1.0})
	delay := func(kr, kc float64) float64 { return 0.69 * (rd*(kc*cw+cl) + kr*rw*(kc*cw/2+cl)) }
	dc, dr := delay(cR, cC), delay(rR, rC)
	return num(edakit.Sprintf("A net has driver resistance %s kohm, wire resistance %s kohm, wire capacitance %s fF and load %s fF; delay is $t = 0.69[R_d(C_w + C_L) + R_w(C_w/2 + C_L)]$ (kohm times fF is ps), and only the wire parasitics change with the RC corner. The Cworst corner scales wire R by %s and wire C by %s; the RCworst corner scales wire R by %s and wire C by %s. What is the larger of the two corner delays, in ps?", f(rd, 1), f(rw, 1), f(cw, 0), f(cl, 0), f(cR, 2), f(cC, 2), f(rR, 2), f(rC, 2)),
		math.Max(dc, dr), 2, edakit.Sprintf("Cworst: %s ps. RCworst: %s ps. The %s corner is worse here. Driver-dominated nets are worst at Cworst; long, resistive wires are worst at RCworst, which is why both must be signed off.", f(dc, 2), f(dr, 2), map[bool]string{true: "Cworst", false: "RCworst"}[dc >= dr]), nil)
}

// emIR checks current density, IR drop and electromigration lifetime scaling.
func emIR(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		w := edakit.Pick([]float64{0.05, 0.1, 0.2, 0.4})
		t := edakit.Pick([]float64{0.05, 0.08, 0.1, 0.2})
		i := float64(edakit.Int(5, 80)) / 10 // mA
		jmax := edakit.Pick([]float64{50, 100, 200})
		j := i / (w * t)
		return num(edakit.Sprintf("A power strap %s um wide and %s um thick carries a DC current of %s mA. Its electromigration limit is %s mA/um^2. What is the current density in the strap, in mA/um^2?", f(w, 2), f(t, 2), f(i, 1), f(jmax, 0)),
			j, 1, edakit.Sprintf("$J = I/(WT)$ = %s / (%s x %s) = %s mA/um^2, which %s the %s mA/um^2 limit.", f(i, 1), f(w, 2), f(t, 2), f(j, 1), map[bool]string{true: "exceeds", false: "is within"}[j > jmax], f(jmax, 0)), nil)
	case 2:
		rs := edakit.Pick([]float64{0.02, 0.04, 0.05, 0.1})
		w := edakit.Pick([]float64{0.5, 1, 2, 4})
		l := float64(50 * edakit.Int(2, 20))
		i := float64(edakit.Int(5, 50)) // mA
		vdd := edakit.Pick([]float64{0.75, 0.8, 0.9})
		drop := i / 1000 * rs * l / w
		pct := 100 * drop / vdd
		return num(edakit.Sprintf("A power strap of sheet resistance %s ohm/sq, width %s um and length %s um delivers %s mA to a block at its far end. The supply is %s V. What is the IR drop along the strap as a percentage of the supply?", f(rs, 2), f(w, 1), f(l, 0), f(i, 0), f(vdd, 2)),
			pct, 2, edakit.Sprintf("R = %s x %s / %s = %s ohm. Drop = %s mA x %s ohm = %s mV, which is %s%% of %s V. A typical signoff budget allows a few percent of the supply for static plus dynamic drop.", f(rs, 2), f(l, 0), f(w, 1), f(rs*l/w, 3), f(i, 0), f(rs*l/w, 3), f(1000*drop, 2), f(pct, 2), f(vdd, 2)), nil)
	}
	n := edakit.Pick([]float64{1, 1.5, 2})
	k := edakit.Pick([]float64{1.2, 1.5, 2, 3})
	ratio := math.Pow(k, -n)
	return num(edakit.Sprintf("Black's equation gives electromigration lifetime $\\text{MTTF} \\propto J^{-n} e^{E_a/kT}$. With current exponent $n = %s$ and temperature unchanged, an ECO raises the current density in a wire by a factor of %s. What is the new lifetime as a fraction of the old?", f(n, 1), f(k, 1)),
		ratio, 4, edakit.Sprintf("The temperature term cancels, leaving $(J_{new}/J_{old})^{-n}$ = %s^{-%s} = %s. Upsizing a driver for timing can quietly cut the EM lifetime of its output wire, which is why EM is re-checked after every timing ECO.", f(k, 1), f(n, 1), f(ratio, 4)), nil)
}

var checklistItems = []edakit.Item{
	{KP: 1, Q: "Which set of checks makes up a complete digital signoff?",
		A: "Timing (setup, hold, SI) across all scenarios, power and EM/IR, physical verification (DRC, LVS, antenna, density) and logical equivalence", D: []string{"Setup timing at the typical corner and DRC", "Simulation of the RTL and a clean synthesis log", "LVS and IR drop only"},
		Why: "Signoff means every independent failure mode has been checked at its worst conditions. Each check catches defects the others cannot see."},
	{KP: 1, Q: "Which signoff check proves that the layout implements the netlist that was timed?",
		A: "Layout versus schematic (LVS)", D: []string{"Design rule checking (DRC)", "Static IR-drop analysis", "Logical equivalence checking between RTL and netlist"},
		Why: "LVS extracts devices and connectivity from the layout and compares them with the netlist. DRC checks manufacturability, not connectivity."},
	{KP: 1, Q: "Why is antenna checking part of physical signoff?",
		A: "Long metal connected only to a gate can collect charge during plasma etching and damage the gate oxide", D: []string{"Antennas cause crosstalk at high frequency", "Antennas increase leakage power at runtime", "Antenna rules limit the number of vias"},
		Why: "During manufacturing, a large metal area connected only to a gate oxide charges up. Antenna rules, fixed with diodes or layer jumps, prevent that damage."},
	{KP: 1, Q: "What does logical equivalence checking (LEC) prove at signoff?",
		A: "That the final netlist, after optimization and ECOs, computes the same function as the reference RTL or netlist", D: []string{"That the layout matches the netlist", "That all paths meet timing", "That the clock tree is balanced"},
		Why: "Physical optimization and ECOs change the netlist. LEC proves formally that the function did not change; LVS separately ties the layout to that netlist."},
	{KP: 2, Q: "A block is clean at the typical corner but has never been timed at the fast corner. What signoff risk remains?",
		A: "Hold violations, which are worst when cells and wires are fast", D: []string{"Setup violations, which are worst when everything is fast", "Antenna violations", "LVS mismatches"},
		Why: "Hold fails when data arrives too early, which happens at fast process, high voltage and the Cbest or RCbest extraction corner."},
	{KP: 2, Q: "Which signoff check needs the coupling capacitances in the SPEF, not just each net's total capacitance?",
		A: "Signal-integrity analysis: crosstalk delay and glitch", D: []string{"Static IR drop", "LVS", "Antenna checking"},
		Why: "SI needs to know which aggressor couples to which victim and by how much. Total capacitance alone is enough for basic delay but not for crosstalk."},
	{KP: 2, Q: "Timing signoff uses voltage-drop-aware delays. Which analysis produces the voltage each cell actually sees?",
		A: "Dynamic IR-drop analysis of the power grid under realistic switching activity", D: []string{"Static timing at nominal voltage", "Design rule checking", "Equivalence checking"},
		Why: "Local supply droop slows cells. IR analysis supplies per-instance voltages that the timer can use to derate delays."},
	{KP: 2, Q: "A design passes DRC with the foundry deck but fails metal density checks. Why does density matter?",
		A: "Chemical-mechanical polishing needs roughly uniform metal density, or thickness varies and parasitics shift", D: []string{"Density affects only aesthetics of the layout", "Density rules limit leakage", "Density only affects antenna ratios"},
		Why: "CMP dishes sparse regions and erodes dense ones. Fill shapes are inserted to meet density windows, and that fill must be included in extraction."},
	{KP: 3, Q: "Metal fill is inserted after the final extraction. What should the signoff engineer do?",
		A: "Re-extract with the fill present, because fill adds capacitance to nearby nets", D: []string{"Nothing; fill is electrically floating so it has no effect", "Remove the fill to keep timing", "Only re-run DRC"},
		Why: "Floating or grounded fill changes the field between wires and raises capacitance. Signoff extraction must see the final geometry."},
	{KP: 3, Q: "An ECO fixes a setup path and the team re-runs only setup timing. Which checks have they skipped that the ECO could break?",
		A: "Hold, SI, EM/IR, DRC, LVS and equivalence, all of which the changed cells and wires can affect", D: []string{"None; an ECO only affects the path it fixes", "Only DRC", "Only the clock tree"},
		Why: "Every ECO changes geometry and netlist. Signoff is re-run across the checklist, often incrementally, before tape-out."},
	{KP: 3, Q: "Why is signoff timing run with extracted parasitics rather than the router's internal estimates?",
		A: "The signoff extractor is calibrated to the foundry's golden solver and process, while router estimates are fast approximations", D: []string{"Router estimates are always pessimistic, so they waste margin", "Routers do not compute capacitance at all", "Signoff timers cannot read router data"},
		Why: "Implementation tools use fast models to steer optimization. Signoff uses the qualified extractor whose accuracy the foundry has validated."},
	{KP: 3, Q: "A chip meets timing at signoff but fails in silicon only when cold, at low supply voltage. Which signoff gap is most likely?",
		A: "Temperature inversion: at low voltage cells can be slowest when cold, and the cold slow corner was not analysed", D: []string{"An LVS error that appears only when cold", "A missing antenna diode", "A name-map error in the SPEF"},
		Why: "Near threshold, the threshold-voltage rise with cooling can outweigh the mobility gain, so the slow corner is cold. Signoff must cover every corner where a path can be worst."},
}

var ecoItems = []edakit.Item{
	{KP: 1, Q: "What distinguishes a functional ECO from a timing ECO?",
		A: "A functional ECO changes what the logic computes; a timing ECO changes cells or wires without changing the function", D: []string{"A functional ECO changes only metal layers", "A timing ECO always adds new flip-flops", "There is no difference; both are bug fixes"},
		Why: "A functional ECO fixes a logic bug, so equivalence to the old netlist must fail in exactly the intended way. A timing ECO must stay logically equivalent."},
	{KP: 1, Q: "Why are late changes made as ECOs rather than by re-running synthesis and place-and-route?",
		A: "A full re-run perturbs the whole design and throws away its closure, while an ECO touches only what must change", D: []string{"Synthesis cannot be run twice on one design", "ECOs are always more optimal", "Re-running the flow is forbidden by the foundry"},
		Why: "Timing, SI and physical closure took weeks. A local change preserves them everywhere else."},
	{KP: 1, Q: "Which is a typical timing ECO?",
		A: "Upsizing a driver or inserting a buffer on a failing path", D: []string{"Changing a state machine's next-state logic", "Adding a new port to the block", "Changing a register's reset value"},
		Why: "Timing ECOs resize, buffer, swap threshold-voltage flavours or reroute, all without changing logic."},
	{KP: 1, Q: "Which is a functional ECO?",
		A: "Inverting the polarity of an enable signal that the specification got wrong", D: []string{"Swapping a cell for its low-Vt equivalent", "Adding a hold buffer", "Rerouting a net to a wider track"},
		Why: "Changing what a signal means changes the function. The other options change only timing or physical properties."},
	{KP: 2, Q: "What does a metal-only ECO change?",
		A: "Only metal and via layers, rewiring existing cells including pre-placed spare cells", D: []string{"Only the transistor layers", "The placement of every cell", "Only the clock tree buffers"},
		Why: "Base layers (transistors) are untouched, so the expensive front-end masks are reused and wafers held before metallization can be finished with the fix."},
	{KP: 2, Q: "Why are spare cells scattered through a design before tape-out?",
		A: "So that a later metal-only ECO can find unused gates nearby and wire them in", D: []string{"To meet metal density rules", "To act as decoupling capacitors only", "To increase the clock frequency"},
		Why: "A metal-only ECO cannot add transistors. Pre-placed NANDs, NORs, inverters and flops, with inputs tied off, are the raw material."},
	{KP: 2, Q: "A metal-only ECO needs a 4-input AND but the nearest spares are 2-input NANDs and inverters. What is the usual approach?",
		A: "Build the function from several spare cells, accepting extra delay and wiring", D: []string{"Add a new 4-input AND cell to the base layers", "Abandon the ECO", "Use a spare flip-flop instead"},
		Why: "Spare gates are generic. Mapping the fix onto available spare types is a small technology-mapping problem."},
	{KP: 2, Q: "Wafers were held after the base layers. A bug fix needs one new gate and some rewiring, and spare cells are available. What is the cheapest path?",
		A: "A metal-only ECO using spare cells, finishing the held wafers with new metal masks", D: []string{"A full respin with all new masks", "A new synthesis run with the same masks", "A firmware patch only"},
		Why: "Only the changed metal and via masks are made; the held wafers resume processing with them."},
	{KP: 3, Q: "After an ECO, the team re-extracts only the nets it touched. What else must be re-extracted?",
		A: "The neighbours of touched wires, because their coupling capacitance changed too", D: []string{"Nothing else", "The entire chip, always", "Only the clock nets"},
		Why: "Coupling is shared: moving a wire changes the capacitance of every wire beside it. Incremental extraction includes the geometric neighbourhood."},
	{KP: 3, Q: "An ECO tool inserts a buffer to fix setup and the hold check on a nearby path now fails. Why might a buffer on one path hurt another?",
		A: "The new buffer and its wire add coupling and load to neighbouring nets and may share logic with short paths", D: []string{"Buffers always speed up all paths", "Hold checks ignore buffers", "Buffers change the clock period"},
		Why: "ECO side effects propagate through shared logic, added coupling and placement displacement. Every ECO needs a full incremental re-check."},
	{KP: 3, Q: "How is a functional ECO proven correct before tape-out?",
		A: "Equivalence checking of the patched netlist against the corrected RTL, plus the full physical and timing re-check", D: []string{"A single simulation of the failing test", "Checking that the ECO tool reported no errors", "Comparing file sizes of the netlists"},
		Why: "Formal equivalence covers all inputs, unlike a test. The physical and timing checks ensure the patch did not break anything else."},
	{KP: 3, Q: "Why do teams prefer spare flip-flops with their clock pins already connected?",
		A: "Connecting a new clock pin late would perturb the balanced clock tree, so prewired clocks keep skew closure intact", D: []string{"Flip-flops without clocks leak more", "A spare flop cannot be used without a scan chain", "Clock pins cannot be wired in metal-only ECOs"},
		Why: "A clock connection changes the tree's load and skew. Pre-connecting spare flops to the clock avoids reworking the clock during the ECO."},
}

// spareCells finds the nearest spare cell, the wirelength to hook it up, and mask counts.
func spareCells(ctx generator.GeneratorContext) generator.Problem {
	type pt struct{ x, y int }
	switch kp(ctx) {
	case 1:
		e := pt{10 * edakit.Int(5, 40), 10 * edakit.Int(5, 40)}
		types := []string{"NAND2", "NOR2", "INV"}
		need := edakit.Pick(types)
		n := edakit.Scaled(ctx, 4, 7)
		var desc []string
		best := math.MaxInt32
		for i := 0; i < n; i++ {
			ty := types[i%3]
			if i == 0 {
				ty = need
			}
			p := pt{10 * edakit.Int(0, 50), 10 * edakit.Int(0, 50)}
			d := int(math.Abs(float64(p.x-e.x)) + math.Abs(float64(p.y-e.y)))
			if ty == need && d < best {
				best = d
			}
			desc = append(desc, edakit.Sprintf("%s at (%d, %d)", ty, p.x, p.y))
		}
		desc = edakit.Shuffle(desc)
		return num(edakit.Sprintf("A metal-only ECO needs one %s near the point (%d, %d) um. The unused spare cells are: %s. Using Manhattan distance, how far is the nearest usable spare cell, in um?", need, e.x, e.y, strings.Join(desc, ", ")),
			float64(best), 0, edakit.Sprintf("Only %s spares qualify. The nearest is %d um away by Manhattan distance, which is what a router would roughly have to wire.", need, best), nil)
	case 2:
		s := pt{10 * edakit.Int(0, 40), 10 * edakit.Int(0, 40)}
		k := edakit.Scaled(ctx, 2, 4)
		xs, ys := []int{s.x}, []int{s.y}
		var desc []string
		for i := 0; i < k; i++ {
			p := pt{10 * edakit.Int(0, 40), 10 * edakit.Int(0, 40)}
			xs = append(xs, p.x)
			ys = append(ys, p.y)
			desc = append(desc, edakit.Sprintf("(%d, %d)", p.x, p.y))
		}
		mm := func(v []int) int {
			lo, hi := v[0], v[0]
			for _, a := range v {
				if a < lo {
					lo = a
				}
				if a > hi {
					hi = a
				}
			}
			return hi - lo
		}
		h := mm(xs) + mm(ys)
		return num(edakit.Sprintf("A spare inverter at (%d, %d) um is wired in to drive %d sink pins at %s. What is the half-perimeter wirelength of the new net, in um?", s.x, s.y, k, strings.Join(desc, ", ")),
			float64(h), 0, edakit.Sprintf("The bounding box spans %d um in x and %d um in y, so HPWL is %d um. A spare far from the sinks adds wire, capacitance and delay to the fix.", mm(xs), mm(ys), h), nil)
	}
	top := edakit.Int(3, 6)
	dp := edakit.Int(1, top-1)
	masks := 0
	var parts []string
	for m := 1; m <= top; m++ {
		mk := 1
		if m <= dp {
			mk = 2
		}
		masks += mk
	}
	vias := top - 1
	masks += vias
	parts = append(parts, edakit.Sprintf("M1 to M%d are double patterned (2 masks each)", dp))
	return num(edakit.Sprintf("A metal-only ECO rewires layers M1 through M%d and via layers V1 through V%d. %s, the remaining changed metal layers need one mask each, and each via layer needs one mask. How many new masks must be made?", top, vias, parts[0]),
		float64(masks), 0, edakit.Sprintf("Metal: %d double-patterned layers x 2 + %d single = %d. Vias: %d. Total %d masks, against the full set for a base-layer respin. Confining an ECO to fewer, upper layers cuts mask cost further.", dp, top-dp, 2*dp+top-dp, vias, masks), nil)
}

var ppaItems = []edakit.Item{
	{KP: 1, Q: "What does PPA closure mean at signoff?",
		A: "Meeting the timing target while keeping power and area within budget, all at signoff accuracy", D: []string{"Meeting timing at any power cost", "Minimizing area while ignoring timing", "Passing DRC and LVS"},
		Why: "Power, performance and area trade against each other. Closure means all three meet their targets together, measured by signoff tools."},
	{KP: 1, Q: "Swapping a cell from high-Vt to low-Vt fixes a setup path. What does it cost?",
		A: "Higher leakage power, often several times that of the high-Vt cell", D: []string{"More area", "Lower switching activity", "A longer clock tree"},
		Why: "Lower threshold voltage speeds the transistor but leaks exponentially more. Vt swaps are a power-for-timing trade."},
	{KP: 1, Q: "Why is upsizing every cell on a critical path rarely the best timing fix?",
		A: "Bigger cells present more input capacitance to their drivers, raise power and area, and can slow the previous stage", D: []string{"Bigger cells are always slower", "Upsizing changes the logic function", "Upsizing removes the path from analysis"},
		Why: "Logical-effort reasoning: a larger gate loads its driver. Gains saturate while power and area keep growing."},
	{KP: 1, Q: "Why do SI fixes such as spacing and shielding compete with area?",
		A: "They consume routing tracks that could otherwise hold wires, pushing congestion or block size up", D: []string{"They add transistors", "They increase leakage directly", "They reduce the clock frequency"},
		Why: "Doubling spacing or adding a shield wire uses track resources. A dense design may not have the room."},
	{KP: 1, Q: "A design meets timing only with a large positive margin on every path. What is the PPA implication?",
		A: "It is probably over-designed, spending power and area that a tighter closure would save", D: []string{"It is optimal in every metric", "It will fail hold", "It will fail LVS"},
		Why: "Excess slack means cells are bigger or faster than needed. Recovering power and area on non-critical paths is part of closure."},
	{KP: 1, Q: "Which fix for a crosstalk-induced setup violation costs no extra leakage?",
		A: "Increasing spacing to the aggressor or rerouting the victim", D: []string{"Swapping the victim driver to low-Vt", "Adding a larger driver", "Raising the supply voltage"},
		Why: "Reducing coupling attacks the cause without faster, leakier devices, at the cost of routing resources."},
	{KP: 3, Q: "Recovering leakage by swapping non-critical cells to high-Vt makes several near-critical paths fail. What went wrong?",
		A: "The recovery used too little slack margin, so the slower cells consumed the slack those paths needed", D: []string{"High-Vt cells have more area", "The swap changed the logic function", "Leakage recovery never affects timing"},
		Why: "Power recovery spends positive slack. It must respect a margin and re-time incrementally after each batch."},
	{KP: 3, Q: "An agent proposes raising the supply voltage by 5% to close timing. Which signoff checks must be revisited?",
		A: "Power and EM/IR, because dynamic power and current rise, along with timing at the new voltage corners", D: []string{"Only DRC", "Only LVS", "None, since voltage does not affect signoff"},
		Why: "Dynamic power scales with the square of voltage and currents rise, stressing EM and IR budgets. Corners must be re-characterized."},
	{KP: 3, Q: "Two fixes close a path equally: upsizing a driver (+2 uW leakage) or rerouting with wider spacing (+0 uW, +3% local congestion). The region is already near overflow. Which is the closure-aware choice?",
		A: "Upsizing the driver, because the region cannot absorb more congestion", D: []string{"Rerouting, because zero leakage always wins", "Both together, for safety", "Neither; relax the clock"},
		Why: "PPA closure weighs the cost of each resource where it is scarce. In a congested region, track resources are the binding constraint."},
	{KP: 3, Q: "Why does closing timing at a fast extraction mode and then signing off at the accurate mode often reopen violations?",
		A: "Fast-mode errors on critical nets shift slack, and paths closed with no margin flip negative", D: []string{"The accurate mode ignores coupling", "Accurate mode uses a different netlist", "Fast mode is always pessimistic"},
		Why: "Closure should leave margin proportional to the fast mode's error, or recurse to accurate extraction on near-critical nets."},
	{KP: 3, Q: "A late ECO trades 30 ps of WNS improvement for 4% more area in one block. When is that worth it?",
		A: "When that block's paths limit the chip frequency and its area does not grow the die or break the floorplan", D: []string{"Always, because timing is the only signoff metric", "Never, because area cannot grow after placement", "Only if leakage falls too"},
		Why: "The value of a PPA trade depends on which constraint binds the product. A frequency-limiting block justifies local area spent."},
	{KP: 3, Q: "Why must power analysis after a timing ECO use the post-ECO parasitics?",
		A: "Switching power depends on the capacitance being switched, which the ECO changed", D: []string{"Power depends only on cell count", "Parasitics affect only timing", "Leakage depends on wire capacitance"},
		Why: "Dynamic power is proportional to switched capacitance times voltage squared times frequency and activity. New wires and cells change it."},
}

// ppaCalc picks the candidate fix that closes timing at least leakage.
func ppaCalc(ctx generator.GeneratorContext) generator.Problem {
	wns := -float64(edakit.Int(5, 40))
	type cand struct {
		name       string
		gain, leak float64
	}
	names := []string{"upsize U12 to X4", "swap U7 and U9 to LVT", "insert a buffer on the long net", "reroute with double spacing", "restructure the AND tree", "swap U3 to LVT"}
	names = edakit.Shuffle(names)
	n := 4
	cs := make([]cand, n)
	found := false
	for try := 0; try < generator.MaxDrawAttempts && !found; try++ {
		wns = -float64(edakit.Int(5, 40))
		for i := range cs {
			cs[i] = cand{names[i], float64(edakit.Int(int(math.Max(1, -wns-12)), -int(wns)+15)), float64(edakit.Int(0, 30)) / 10}
		}
		ok := 0
		for _, c := range cs {
			if wns+c.gain >= 0 {
				ok++
			}
		}
		found = ok >= 1 && ok < n
	}
	if !found {
		// Known-good fallback: two fixes close, two do not.
		wns = -20
		cs[0] = cand{names[0], 25, 1.2}
		cs[1] = cand{names[1], 30, 0.4}
		cs[2] = cand{names[2], 12, 0.1}
		cs[3] = cand{names[3], 18, 0.6}
	}
	best := -1
	for i, c := range cs {
		if wns+c.gain >= 0 && (best < 0 || c.leak < cs[best].leak) {
			best = i
		}
	}
	// Make the best unique.
	for i, c := range cs {
		if i != best && wns+c.gain >= 0 && c.leak == cs[best].leak {
			cs[i].leak += 0.5
		}
	}
	var desc []string
	for _, c := range cs {
		desc = append(desc, edakit.Sprintf("%s: +%s ps, +%s uW leakage", c.name, f(c.gain, 0), f(c.leak, 1)))
	}
	var d []string
	for i, c := range cs {
		if i != best {
			d = append(d, c.name)
		}
	}
	q := edakit.Sprintf("A failing endpoint has WNS %s ps. Each candidate fix improves that slack and adds leakage as follows: %s. Fixes are applied one at a time. Which single fix closes the endpoint (slack at least 0) at the least added leakage?", f(wns, 0), strings.Join(desc, "; "))
	return edakit.Choice(q, cs[best].name, edakit.Sprintf("Only fixes with gain of at least %s ps close the endpoint. Among those, %s adds the least leakage. A fix that adds less leakage but leaves negative slack is not a closure.", f(-wns, 0), cs[best].name), d...)
}
