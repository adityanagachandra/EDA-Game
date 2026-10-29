package pcb

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type bomPart struct {
	prefix string
	desc   string
}

var bomPool = []bomPart{
	{"R", "10 kohm 0402 resistor"}, {"R", "10 kohm 0603 resistor"}, {"R", "1 kohm 0402 resistor"},
	{"R", "22 ohm 0402 resistor"}, {"R", "4.7 kohm 0402 resistor"},
	{"C", "100 nF 0402 X7R capacitor"}, {"C", "100 nF 0603 X7R capacitor"}, {"C", "4.7 uF 0805 X5R capacitor"},
	{"C", "10 uF 0805 X5R capacitor"}, {"C", "22 pF 0402 C0G capacitor"},
	{"D", "green 0603 LED"}, {"U", "QFN-48 microcontroller"}, {"U", "SOT-23-5 LDO regulator"},
	{"J", "USB-C receptacle"}, {"Y", "12 MHz 3225 crystal"},
}

// bom: count BOM lines (KP 1), order quantity with attrition (KP 2), placements of a variant (KP 3).
func bom(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 8, 16)
		counter := map[string]int{}
		distinct := map[string]bool{}
		var parts []string
		for i := 0; i < n; i++ {
			p := edakit.Pick(bomPool)
			counter[p.prefix]++
			distinct[p.desc] = true
			parts = append(parts, edakit.Sprintf("%s%d: %s", p.prefix, counter[p.prefix], p.desc))
		}
		ans := len(distinct)
		q := "A schematic contains these placements: " + strings.Join(parts, "; ") +
			". The BOM has one line per manufacturer part number, and each distinct value-and-package description above is one part number. How many lines does the BOM have?"
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Group the %d placements by description. Same value in a different package is a different part number, so it gets its own line. There are %d distinct descriptions.", n, ans),
			map[string]string{"placements": edakit.IntStr(n), "lines": edakit.IntStr(ans)})
	case 2:
		perBoard := edakit.Int(2, 40)
		boards := edakit.Pick([]int{25, 50, 100, 150, 250, 400, 500})
		attr := edakit.Pick([]int{2, 3, 5, 8, 10})
		part := edakit.Pick(bomPool[:10]).desc
		exact := float64(perBoard*boards) * (1 + float64(attr)/100)
		ans := ceilInt(exact)
		q := edakit.Sprintf("A BOM line for the %s has quantity %d per board. The build is %d boards, and the assembler asks for %d%% extra parts on every line for attrition (pick-up losses and machine setup). How many parts should be ordered for this line? Round up to a whole part.", part, perBoard, boards, attr)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("The build needs %d x %d = %d parts. Adding %d%% gives %s, rounded up to %d.", perBoard, boards, perBoard*boards, attr, edakit.Num(exact, 2), ans),
			map[string]string{"net_qty": edakit.IntStr(perBoard * boards)})
	}
	n := edakit.Scaled(ctx, 5, 8)
	idx := edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14})[:n]
	sum, dnp := 0, 0
	var lines []string
	for i, k := range idx {
		qty := edakit.Int(1, 12)
		mark := ""
		if i > 0 && edakit.Int(1, 3) == 1 {
			mark = " (DNP in this variant)"
			dnp += qty
		} else {
			sum += qty
		}
		lines = append(lines, edakit.Sprintf("%s, quantity %d%s", bomPool[k].desc, qty, mark))
	}
	q := "A BOM for the low-cost variant of a board lists: " + strings.Join(lines, "; ") +
		". DNP (do not populate) lines keep their footprints on the board but are not placed. How many component placements does the pick-and-place program make per board for this variant?"
	return edakit.Numeric(q, edakit.IntStr(sum),
		edakit.Sprintf("Sum the quantities of the populated lines only: %d placements. The %d DNP parts still have land patterns in the copper, which is why one board serves several variants.", sum, dnp),
		map[string]string{"placements": edakit.IntStr(sum), "dnp": edakit.IntStr(dnp)})
}

// erc draws a handful of nets with typed pins and counts the nets a stated rule set flags.
func erc(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	types := []string{"output", "input", "bidirectional", "passive"}
	if kp == 3 {
		types = append(types, "power input", "power output")
	}
	n := edakit.Scaled(ctx, 4, 7)
	var nets []string
	flagged := 0
	var why []string
	for i := 0; i < n; i++ {
		k := edakit.Pick([]int{1, 2, 2, 3, 3, 4})
		var pins []string
		cnt := map[string]int{}
		used := map[string]bool{}
		for j := 0; j < k; j++ {
			ty := edakit.Pick(types)
			if kp == 1 && ty == "input" && edakit.Int(1, 2) == 1 {
				ty = "output"
			}
			cnt[ty]++
			ref := edakit.Sprintf("U%d pin %d", 6+j, 30+i)
			for a := 0; a < generator.MaxDrawAttempts; a++ {
				r := edakit.Sprintf("U%d pin %d", edakit.Int(1, 5), edakit.Int(1, 24))
				if !used[r] {
					ref = r
					break
				}
			}
			used[ref] = true
			pins = append(pins, edakit.Sprintf("%s (%s)", ref, ty))
		}
		name := edakit.Sprintf("N%d", i+1)
		nets = append(nets, name+": "+strings.Join(pins, ", "))
		bad := false
		drivers := cnt["output"] + cnt["power output"]
		if k == 1 {
			bad = true
		}
		if drivers >= 2 {
			bad = true
		}
		if kp >= 2 && cnt["input"] > 0 && cnt["output"]+cnt["bidirectional"]+cnt["passive"]+cnt["power output"] == 0 {
			bad = true
		}
		if kp == 3 && cnt["power input"] > 0 && cnt["power output"] == 0 {
			bad = true
		}
		if bad {
			flagged++
			why = append(why, name)
		}
	}
	rules := "(a) the net has only one pin, or (b) it has two or more pins of type output or power output"
	if kp >= 2 {
		rules += ", or (c) it has an input pin but no pin that can drive it (output, bidirectional, passive or power output)"
	}
	if kp == 3 {
		rules += ", or (d) it has a power input pin but no power output pin"
	}
	q := "An ERC flags a net if " + rules + ". The nets are: " + strings.Join(nets, "; ") + ". How many nets does the ERC flag?"
	exp := "Test every net against each rule and count a net once even if it breaks two rules."
	if flagged > 0 {
		sort.Strings(why)
		exp += " Flagged: " + strings.Join(why, ", ") + "."
	} else {
		exp += " No net breaks a rule."
	}
	return edakit.Numeric(q, edakit.IntStr(flagged), exp, map[string]string{"flagged": edakit.IntStr(flagged)})
}

var flowItems = []edakit.Item{
	{KP: 1, Q: "Which ordering matches a typical printed-circuit-board flow?",
		A:   "Schematic capture, footprint assignment, netlist transfer, placement, routing, DRC, manufacturing outputs",
		D:   []string{"Placement, schematic capture, routing, footprint assignment, DRC, manufacturing outputs", "Schematic capture, routing, placement, footprint assignment, manufacturing outputs, DRC", "Footprint assignment, routing, schematic capture, placement, DRC, manufacturing outputs"},
		Why: "Connectivity is captured first, each part gets a physical footprint, the netlist moves to layout, parts are placed and then routed, and the checked board is written out for fabrication and assembly."},
	{KP: 1, Q: "A chip flow synthesizes its gate netlist from RTL. Where does a board flow's netlist come from?",
		A:   "A designer draws it as a schematic of purchased components; there is no synthesis step",
		D:   []string{"The board tool synthesizes it from HDL", "It is extracted from the finished copper layout", "The autorouter generates it from the BOM"},
		Why: "Board connectivity is authored directly in schematic capture. The parts already exist as packaged chips and passives, so nothing is synthesized."},
	{KP: 1, Q: "How do the objects a board placer moves differ from those a chip placer moves?",
		A:   "Hundreds to a few thousand packages of very different sizes, versus millions of near-uniform cells in rows",
		D:   []string{"Millions of standard cells in rows in both cases", "Single transistors on boards, whole packages on chips", "Nets rather than components on boards"},
		Why: "A board holds a modest number of heterogeneous packages placed freely in two dimensions, so board placement is mostly constraint-driven and interactive rather than analytical over millions of cells."},
	{KP: 1, Q: "Who sets the design rules a board must meet, in place of a chip's foundry PDK?",
		A:   "The chosen fabricator's and assembler's capabilities, usually stated against IPC classes",
		D:   []string{"A foundry process design kit for the board's silicon", "The schematic symbol library", "The autorouter's default grid"},
		Why: "Minimum trace and space, drill sizes, annular rings and mask clearances come from the board shop's process capability, with IPC standards supplying common definitions and acceptance classes."},
	{KP: 1, Q: "Which statement draws the boundary between chip EDA and board EDA correctly?",
		A:   "Board EDA interconnects off-the-shelf packages on a few to a few dozen layers; chip EDA creates the devices and wiring on the die",
		D:   []string{"Board EDA and chip EDA use the same standard-cell placer and differ only in units", "Chip EDA ends at the package; board EDA designs the transistors", "Board EDA has no design rules because fabrication is cheap"},
		Why: "The board flow starts where the chip flow ends: the die is already packaged, and the board tool wires packages together and models the board's own interconnect."},
	{KP: 2, Q: "Which artifact carries connectivity from schematic capture into layout?",
		A:   "The netlist: each net with its component pins, plus the footprint chosen for each component",
		D:   []string{"The Gerber files", "The BOM on its own", "The pick-and-place centroid file"},
		Why: "Layout needs to know which pads must be joined and which footprint each reference designator uses. The netlist (or a direct schematic-to-layout update) carries exactly that."},
	{KP: 2, Q: "Which data set does a bare-board fabricator need?",
		A:   "Copper, solder mask, silkscreen and outline layers, drill data and a fabrication drawing with the stackup",
		D:   []string{"The schematic PDF and the BOM", "The pick-and-place file and paste layers only", "The IBIS models of every driver"},
		Why: "The bare-board shop images, drills and laminates the board. It needs the layer images, the drill program and the build notes, not the component data."},
	{KP: 2, Q: "Which outputs go to the assembly house rather than the bare-board fabricator?",
		A:   "The BOM, the component placement (centroid) data and the solder paste layers",
		D:   []string{"The drill files and the impedance table", "The copper layers and the bare-board test netlist", "The stackup drawing only"},
		Why: "Assembly buys and places parts. It needs what to place (BOM), where and at what rotation (centroid), and where to print paste (stencil apertures from the paste layers)."},
	{KP: 2, Q: "Both flows end in a geometry database for manufacturing. Which pairing is right?",
		A:   "Board: Gerber or ODB++ layer data plus drill files; chip: a GDSII or OASIS stream of mask layers",
		D:   []string{"Board: GDSII; chip: Gerber", "Both: Liberty files", "Board: a SPICE netlist; chip: Gerber files"},
		Why: "Board fabrication consumes photoplotter-style layer images and drill programs. Mask making consumes a hierarchical polygon stream."},
	{KP: 3, Q: "During layout a connector pinout turns out to be mirrored relative to the mating cable. Where should the fix go?",
		A:   "Into the schematic or the footprint library, then forward to layout, so both views stay consistent",
		D:   []string{"Only into the copper, by rerouting traces by hand", "Into the Gerber files after output", "Into the BOM quantity column"},
		Why: "If the copper is patched alone, the schematic no longer describes the board, and the next netlist update will undo or conflict with the fix."},
	{KP: 3, Q: "Routing a dense board stalls because a BGA sits near the edge with its escape side facing a connector. What does this illustrate?",
		A:   "A placement decision limits routability later, so the flow iterates back to placement",
		D:   []string{"Routing is independent of placement", "DRC caused the routing failure", "The schematic must be redrawn from scratch"},
		Why: "Routing can only use the channels placement leaves. Poor orientation or crowding shows up later as unroutable regions, and the fix is upstream."},
	{KP: 3, Q: "Why are high-speed nets often routed interactively while the rest of the board is autorouted?",
		A:   "Their constraints on impedance, length matching and reference planes are hard to express and verify in a general autorouter",
		D:   []string{"Autorouters cannot route nets shorter than one inch", "Interactive routing is faster for every net", "IPC standards forbid autorouting high-speed nets"},
		Why: "Generic routers optimize completion and wirelength. High-speed nets need topology, return-path and matching decisions that designers usually make and check themselves."},
	{KP: 3, Q: "After layout, the fabricator's stackup proposal changes a dielectric thickness. What must be rechecked?",
		A:   "The trace widths for controlled-impedance nets, since impedance depends on height to the reference plane",
		D:   []string{"Nothing, because impedance depends only on trace width", "Only the silkscreen legend", "Only the BOM quantities"},
		Why: "Impedance rises with dielectric height. A thicker dielectric needs wider traces to hold the target, so the impedance table and routing must be revisited."},
}

var footprintItems = []edakit.Item{
	{KP: 1, Q: "What does a schematic symbol carry that a footprint does not?",
		A:   "Logical pins with names and electrical types (input, output, power) used by ERC",
		D:   []string{"The copper pad geometry", "The courtyard and silkscreen outlines", "The 3D body model of the package"},
		Why: "The symbol is the logical view. Pin electrical types exist so the schematic checker can reason about drivers and loads; the footprint is purely physical."},
	{KP: 1, Q: "What links a symbol pin to a footprint pad?",
		A:   "The symbol pin number matching the footprint pad number",
		D:   []string{"The pin name matching the net name", "The pin's position on the schematic sheet", "The alphabetical order of the pin names"},
		Why: "Netlisting writes nets as reference designator plus pin number, and layout attaches each to the pad with that number. Names are for humans."},
	{KP: 1, Q: "In IPC terms, what is the land pattern of a component?",
		A:   "The copper pad geometry on the board that the component's terminations solder to",
		D:   []string{"The silkscreen outline drawn around the part", "The schematic symbol of the part", "The paste stencil aperture only"},
		Why: "A land is a pad; the land pattern is the set of pads for one package. The footprint adds mask, paste, silkscreen, courtyard and assembly data around it."},
	{KP: 1, Q: "A dual op-amp has a symbol with units A and B, both placed in the schematic. What does the tool do?",
		A:   "Assigns both units to one package (U1A, U1B) and maps each unit's pins to different pads of one footprint",
		D:   []string{"Places two separate packages on the board", "Merges both units' outputs into one net", "Rejects the second unit as a duplicate"},
		Why: "Multi-unit symbols split one package across the schematic. Each unit carries its own pin numbers, so the footprint stays one package."},
	{KP: 2, Q: "A transistor symbol uses pins 1 = B, 2 = E, 3 = C, but the vendor's SOT-23 package is 1 = B, 2 = C, 3 = E. What happens?",
		A:   "The board is built with collector and emitter swapped, and no layout check catches it",
		D:   []string{"DRC reports a clearance violation", "ERC reports an output conflict", "Gerber export fails"},
		Why: "Copper matches the netlist exactly, so every check passes. The error lives in the pin-to-pad map, which is why library review is a separate discipline."},
	{KP: 2, Q: "Why does each footprint carry a courtyard?",
		A:   "It is the keep-out envelope used to check component-to-component spacing during placement",
		D:   []string{"It defines the copper pads", "It is the paste stencil aperture", "It sets the trace width to the pads"},
		Why: "The courtyard bounds the body, the pads and an assembly margin. Placement DRC checks courtyards for overlap so parts can be placed and reworked."},
	{KP: 2, Q: "Which footprint layer becomes the solder stencil apertures?",
		A:   "The paste layer",
		D:   []string{"The solder mask layer", "The silkscreen layer", "The courtyard layer"},
		Why: "Paste openings define where solder paste is printed. Mask openings define where the board's mask is removed, which is a different layer for a different process."},
	{KP: 2, Q: "Why is the solder mask opening around a pad usually slightly larger than the pad?",
		A:   "To tolerate mask-to-copper registration error without mask covering part of the pad",
		D:   []string{"To make the pad electrically larger", "To raise the impedance of the trace", "To hold the silkscreen ink in place"},
		Why: "Mask is imaged in a separate step with its own registration error. A small expansion keeps the whole land solderable."},
	{KP: 3, Q: "IPC-7351 defines three density levels for land patterns. What distinguishes them?",
		A:   "The size of the toe, heel and side solder fillets: most, nominal and least material",
		D:   []string{"The number of copper layers in the board", "Round pads versus square pads only", "The copper weight of the pads"},
		Why: "The levels trade density against solder-joint robustness. Dense consumer boards use smaller lands; rugged designs use larger fillets."},
	{KP: 3, Q: "Why compute pad sizes from the package tolerances rather than the nominal dimensions?",
		A:   "The pads must form good fillets over the whole range of lead positions the tolerances allow",
		D:   []string{"Nominal dimensions are never published", "Tolerances matter only for through-hole parts", "It makes the footprint smaller than the part"},
		Why: "Land-pattern calculators combine package, fabrication and placement tolerances so that a worst-case part still solders properly."},
	{KP: 3, Q: "Why do many teams keep one verified component library instead of drawing parts per project?",
		A:   "A footprint or pinout error repeats on every board that uses the part, so verifying once pays off",
		D:   []string{"A central library makes the autorouter faster", "IPC forbids project-local parts", "Gerber export requires a central library"},
		Why: "Library errors are silent to DRC and ERC. Review, versioning and reuse turn one careful verification into many correct boards."},
	{KP: 3, Q: "A QFN's exposed thermal pad is drawn with one full-size paste opening. Which assembly defect is likely?",
		A:   "Excess solder that floats the part, causing voids or opens on the perimeter pins",
		D:   []string{"Tombstoning of the QFN", "No solder at all on the thermal pad", "Mask-to-copper registration error"},
		Why: "Large paste deposits lift the package. Footprints split the thermal paste into a window pane of smaller openings covering roughly half the area or a bit more."},
	{KP: 3, Q: "A footprint carries a 3D body model. What can the tool check with it?",
		A:   "Mechanical clearance against the enclosure and neighboring parts",
		D:   []string{"Controlled impedance of the attached traces", "Netlist connectivity", "The thermal-relief spoke width"},
		Why: "3D models feed collision checks with the mechanical CAD model of the enclosure. They do not affect the electrical checks."},
}

var netlistItems = []edakit.Item{
	{KP: 1, Q: "What does a board netlist list?",
		A:   "Each net with its component pins as reference designator and pin number, plus each component's footprint",
		D:   []string{"The trace geometry and widths", "The drill sizes and plating", "The stackup and dielectric constants"},
		Why: "The netlist is pure connectivity plus the footprint binding. Geometry is what layout adds."},
	{KP: 1, Q: "What does annotation do in schematic capture?",
		A:   "Assigns a unique reference designator such as R12 or U3 to every component",
		D:   []string{"Adds explanatory notes to the printed schematic", "Assigns net names to every pin", "Generates a footprint for each symbol"},
		Why: "Reference designators are the keys that join schematic, netlist, BOM, layout and assembly data."},
	{KP: 1, Q: "In a hierarchical schematic, how does a local net label differ from a global label?",
		A:   "A local label connects only within its own sheet; a global label connects everywhere it appears with the same name",
		D:   []string{"A local label connects across all sheets; a global label only within one sheet", "Both connect only within one sheet", "Local labels name buses and global labels name single wires"},
		Why: "Scoping lets a reused sheet use its own local names without shorting to other instances, while global labels and power symbols join nets across the design."},
	{KP: 1, Q: "Which mathematical object models a board netlist?",
		A:   "A hypergraph: components are vertices and each net is a hyperedge on the components it touches",
		D:   []string{"A tree rooted at the power connector", "A DAG ordered by signal flow", "A planar graph by construction"},
		Why: "A net can join many pins at once, so it is a hyperedge, the same representation used for chip netlists."},
	{KP: 1, Q: "Power symbols named GND and +3V3 appear on several sheets. What do they do?",
		A:   "They create global nets that join every pin attached to a power symbol of the same name",
		D:   []string{"They are graphics only and connect nothing", "They create a separate net on each sheet", "They place a regulator on the board"},
		Why: "Power symbols behave like global labels, which is why one misspelled power name silently creates a separate floating net."},
	{KP: 2, Q: "A designer changes a resistor value and adds a capacitor in the schematic while layout is in progress. What brings the change into the layout?",
		A:   "Forward annotation, an ECO or netlist update that changes components and nets but keeps placed parts",
		D:   []string{"Back annotation", "Regenerating the Gerber files", "Re-running DRC on the layout"},
		Why: "Forward annotation pushes schematic edits to layout as an incremental update, so existing placement and routing survive."},
	{KP: 2, Q: "During layout, gates of a quad NAND are swapped to shorten routes. What must happen next?",
		A:   "Back annotation updates the schematic so its gate units and pin numbers match the layout",
		D:   []string{"Nothing, because the schematic is unaffected", "Only the BOM is regenerated", "Forward annotation overwrites the swap"},
		Why: "The swap changed which pins carry which nets. Back annotation keeps the schematic a true description of the board."},
	{KP: 2, Q: "Why do tools match components by a stable unique identifier rather than by reference designator during an update?",
		A:   "Designators get renumbered, and matching on them would rip up placed parts or attach them to the wrong parts",
		D:   []string{"Designators are not stored in the netlist", "Unique identifiers make routing faster", "IPC standards require identifiers"},
		Why: "A renumbering pass changes R7 to R12 without changing the part. A hidden identifier lets the layout keep its placement through the rename."},
	{KP: 2, Q: "Which pin swaps may layout perform without changing the circuit?",
		A:   "Swaps between pins declared equivalent, such as the two inputs of one NAND gate",
		D:   []string{"Any two pins on the same package", "A power pin with a ground pin", "Pins with the same number on different parts"},
		Why: "Libraries declare swap groups for logically equivalent pins and gates. Anything outside those groups changes the circuit."},
	{KP: 3, Q: "After a netlist import, a ratsnest line runs to a pad of an unrelated part. What is the likely cause?",
		A:   "A pin-to-pad numbering mismatch between a symbol and its footprint",
		D:   []string{"A DRC clearance violation", "The wrong copper weight", "A Gerber aperture error"},
		Why: "Layout connects nets to pad numbers. If the footprint numbers pads differently from the symbol, nets land on the wrong pads."},
	{KP: 3, Q: "Which check compares the routed copper's connectivity against the schematic netlist?",
		A:   "A layout-versus-netlist connectivity check, the board analogue of LVS",
		D:   []string{"ERC on the schematic", "The impedance calculation", "A silkscreen overlap check"},
		Why: "The tool extracts which pads the copper joins and compares those nets with the netlist, flagging opens and shorts."},
	{KP: 3, Q: "A bare-board test house receives an IPC-D-356 netlist. What does it do with it?",
		A:   "Tests the fabricated bare board for opens and shorts against the intended connectivity",
		D:   []string{"Places the components", "Programs the microcontroller", "Cuts the solder stencil"},
		Why: "Flying-probe or fixture testers check every net's continuity and isolation before any part is mounted."},
	{KP: 3, Q: "A schematic from one tool must drive layout in another. What is the usual path?",
		A:   "Export a netlist the layout tool can import, with footprint names mapped to its library",
		D:   []string{"Send the Gerber files to the layout tool", "Copy the schematic PDF", "Export a GDSII file"},
		Why: "Connectivity crosses tools as a netlist. The footprint mapping is the fragile part, because each tool names its library differently."},
}

var toolItems = []edakit.Item{
	{KP: 1, Q: "Which tool is an open-source suite with schematic and board editors, an interactive push-and-shove router and S-expression text files?",
		A: "KiCad", D: []string{"Altium Designer", "Cadence Allegro", "Zuken CR-8000"},
		Why: "KiCad stores schematics and boards as readable S-expression files and ships an interactive push-and-shove router."},
	{KP: 1, Q: "What is Lepton EDA?",
		A:   "A community fork of the gEDA schematic and netlisting tools",
		D:   []string{"A commercial successor to EAGLE", "Siemens' entry-level board tool", "A Gerber viewer and nothing more"},
		Why: "Lepton continues the gEDA/gaf code base (schematic editor and netlister) with Scheme scripting."},
	{KP: 1, Q: "Which open-source tool builds its design around a versioned library in which every element carries a UUID?",
		A: "LibrePCB", D: []string{"gEDA PCB", "PADS", "OrCAD Capture"},
		Why: "LibrePCB identifies symbols, packages and devices by UUID so that libraries can be versioned and shared without name clashes."},
	{KP: 1, Q: "Freerouting is an open-source board autorouter. How does it exchange a board with an editor such as KiCad?",
		A:   "Through a Specctra DSN design file in and a Specctra SES session file back",
		D:   []string{"Through Gerber files", "Through an ODB++ archive", "Through a GDSII stream"},
		Why: "The Specctra DSN and SES formats describe the board, rules and resulting wires, and many editors support them for external routing."},
	{KP: 2, Q: "Which vendor and product pairing is correct?",
		A: "Siemens: Xpedition and PADS", D: []string{"Cadence: Xpedition and PADS", "Zuken: Allegro and OrCAD", "Altium: CR-8000"},
		Why: "Xpedition (enterprise) and PADS (mainstream) came to Siemens with Mentor Graphics."},
	{KP: 2, Q: "Which vendor and product pairing is correct?",
		A: "Cadence: Allegro and OrCAD", D: []string{"Siemens: Allegro and OrCAD", "Autodesk: Allegro", "Zuken: OrCAD"},
		Why: "Cadence offers OrCAD (schematic capture and mainstream layout) and Allegro (enterprise layout and constraint management)."},
	{KP: 2, Q: "Which products belong to Zuken?",
		A: "CR-8000 and eCADSTAR", D: []string{"Xpedition and PADS", "Allegro and OrCAD", "Altium Designer and Altium 365"},
		Why: "Zuken's board and package design lines include CR-8000 and eCADSTAR."},
	{KP: 2, Q: "What happened to EAGLE after Autodesk acquired it from CadSoft?",
		A:   "Its electronics functionality moved into Autodesk Fusion's electronics workspace",
		D:   []string{"It was open-sourced and became KiCad", "It was merged into Cadence Allegro", "It became Siemens PADS"},
		Why: "Autodesk folded EAGLE's schematic and board capability into Fusion, its combined mechanical and electronics product."},
	{KP: 3, Q: "A team wants to run a custom board check over hundreds of designs with an open toolchain. Which KiCad property helps most?",
		A:   "Readable text board files plus a Python scripting interface",
		D:   []string{"An encrypted binary database", "Cloud-only design storage", "Output limited to Gerber files"},
		Why: "Text files can be parsed, diffed and version-controlled, and scripting exposes the board model to custom checks."},
	{KP: 3, Q: "When a board moves between two commercial tools, what usually survives translation best?",
		A:   "Connectivity and placed geometry; constraint sets and rule definitions often need rebuilding",
		D:   []string{"Everything, since all tools share one database", "Only the constraint rules", "Only the BOM"},
		Why: "Geometry and nets are common to every tool, but each tool models constraints differently, so rules are the first thing lost."},
	{KP: 3, Q: "Which capability do enterprise tools such as Allegro and Xpedition stress for large teams?",
		A:   "Central constraint management and concurrent multi-user editing of one board",
		D:   []string{"Logic synthesis from RTL", "Mask data preparation", "Standard-cell row placement"},
		Why: "Large boards are routed by several designers at once against a shared set of electrical rules, so constraint management is the core feature."},
	{KP: 3, Q: "Why do board tools import and export Specctra DSN even when they have their own router?",
		A:   "It is a widely supported way to hand a placed board to an external autorouter and get the wires back",
		D:   []string{"Fabricators require it for manufacturing", "It encodes the stackup for the fabricator", "It is the IPC bare-board test format"},
		Why: "DSN carries the board outline, components, nets, rules and existing wiring, which is exactly what a router needs and nothing more."},
	{KP: 3, Q: "The gEDA and Lepton netlister uses back ends written in Scheme. What does this design allow?",
		A:   "A new netlist format can be added as a small back-end script without touching the schematic editor",
		D:   []string{"It makes the schematic editor run faster", "It removes the need for footprints", "It lets the netlister route the board"},
		Why: "Separating the connectivity model from output writers is a seam: each target format is a plug-in over one stable netlist."},
}
