package fab

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// designRules: pitch from width and space, landed-via pitch, and width-dependent spacing.
func designRules(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		w := float64(edakit.Int(10, 40))
		s := float64(edakit.Int(10, 40))
		ans := edakit.Num(w+s, 0)
		q := edakit.Sprintf("A metal layer has a minimum width of %.0f nm and a minimum spacing of %.0f nm. What is the tightest legal line pitch, in nanometres?", w, s)
		return edakit.Numeric(q, ans, edakit.Sprintf("Pitch is width plus space: %.0f + %.0f = %s nm.", w, s, ans), map[string]string{"pitch": ans})
	case 2:
		v := float64(edakit.Int(10, 30))
		e := float64(edakit.Int(2, 10))
		s := float64(edakit.Int(10, 30))
		ans := edakit.Num(v+2*e+s, 0)
		q := edakit.Sprintf("Vias are %.0f nm square, the metal above must enclose a via by %.0f nm on each side across the line width, and line-to-line spacing is %.0f nm. If every line must be able to take a via anywhere, what is the minimum line pitch, in nanometres?", v, e, s)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The line must be at least $%.0f + 2 \\times %.0f = %.0f$ nm wide; add the %.0f nm space for a pitch of %s nm. Enclosure rules, not just width, often set pitch.", v, e, v+2*e, s, ans),
			map[string]string{"width": edakit.Num(v+2*e, 0)})
	default:
		wmin := float64(2 * edakit.Int(8, 20))
		s1 := float64(edakit.Int(14, 30))
		th1 := float64(10 * edakit.Int(6, 12))
		th2 := th1 + float64(10*edakit.Int(8, 20))
		s2 := s1 + float64(edakit.Int(6, 20))
		s3 := s2 + float64(edakit.Int(10, 30))
		W := float64(2 * edakit.Int(int(wmin/2), int(th2/2)+40))
		sp := s1
		if W > th2 {
			sp = s3
		} else if W > th1 {
			sp = s2
		}
		pitch := W/2 + sp + wmin/2
		ans := edakit.Num(pitch, 1)
		q := edakit.Sprintf("A spacing table says: if the wider of two wires is at most %.0f nm, spacing is %.0f nm; if it is wider than %.0f nm and at most %.0f nm, spacing is %.0f nm; above %.0f nm, spacing is %.0f nm. A %.0f nm wire runs beside a minimum %.0f nm wire. What is their minimum centre-to-centre distance, in nanometres? Round to 1 decimal place.", th1, s1, th1, th2, s2, th2, s3, W, wmin)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The wider wire, %.0f nm, selects a spacing of %.0f nm. Centre to centre is $%.0f/2 + %.0f + %.0f/2 = %s$ nm.", W, sp, W, sp, wmin, ans),
			map[string]string{"spacing": edakit.Num(sp, 0)})
	}
}

var pdkContentsItems = []edakit.Item{
	{KP: 1, Q: "Which item is part of a process design kit?",
		A:   "SPICE device models with corner and statistical variants",
		D:   []string{"The customer's RTL", "The synthesized netlist of the design", "The design's timing constraints"},
		Why: "The PDK carries everything about the process that a design needs: models, rules, technology files and parameterized devices. Design data is the customer's."},
	{KP: 1, Q: "What is a parameterized cell (PCell) in a PDK?",
		A:   "A generator that draws a device's layout from parameters such as width, length and fingers",
		D:   []string{"A standard cell with a fixed layout", "A timing arc in a Liberty file", "A design-rule violation marker"},
		Why: "PCells encode the foundry's device layout rules in code, so a designer sets parameters and gets a correct layout."},
	{KP: 1, Q: "A router needs layer directions, pitches, via definitions and minimum widths. In what PDK collateral does it find them?",
		A:   "The technology LEF or the router's technology file",
		D:   []string{"The SPICE model card", "The Liberty timing library", "The design-rule manual's figures only"},
		Why: "Place-and-route tools read an abstract of the rules, such as technology LEF, rather than the full signoff deck."},
	{KP: 1, Q: "Which PDK file tells extraction tools how thick each layer is and what dielectric surrounds it?",
		A:   "The interconnect technology or extraction rule file",
		D:   []string{"The PCell library", "The GDS layer map only", "The DRC deck"},
		Why: "Extractors use the layer stack description (thicknesses, permittivities, resistivities, width-dependent tables) to compute parasitics."},
	{KP: 2, Q: "Two tools disagree about which GDS layer number is the first metal. Which PDK item should settle it?",
		A:   "The layer map that assigns names and purposes to GDS layer and datatype numbers",
		D:   []string{"The SPICE models", "The cell library's timing", "The netlist"},
		Why: "GDS stores only integers; the layer map gives each layer-datatype pair its meaning, and every tool must use the same one."},
	{KP: 2, Q: "Are standard-cell libraries part of the PDK?",
		A:   "Usually not: libraries are built on top of the PDK by the foundry or an IP vendor",
		D:   []string{"Yes, the PDK is only the cell library", "Libraries replace the PDK", "Libraries contain the design rules"},
		Why: "The PDK describes the process. Cell libraries, memories and I/O are IP built and characterized using it."},
	{KP: 2, Q: "Why does a PDK ship both a design-rule manual and a DRC runset?",
		A:   "The manual is the human-readable definition; the runset encodes it as executable checks",
		D:   []string{"They contain unrelated rules", "The runset is only for analog designs", "The manual is generated from the netlist"},
		Why: "Designers read the manual; signoff tools run the deck. Keeping them consistent is a quality task for the foundry."},
	{KP: 2, Q: "An open-source PDK lets an EDA researcher do what that a closed one usually does not?",
		A:   "Publish full flows, layouts and results on a real process",
		D:   []string{"Skip design rules", "Fabricate without masks", "Avoid device models"},
		Why: "Closed PDKs come under non-disclosure agreements; open PDKs let tool results be reproduced and shared."},
	{KP: 3, Q: "A foundry releases a new PDK version that changes one spacing value. What must a design team do with existing layouts?",
		A:   "Rerun signoff checks with the new deck and fix any new violations",
		D:   []string{"Nothing, PDK versions are always compatible", "Redesign the RTL", "Only rerun simulation"},
		Why: "Rule changes can make previously clean layouts fail; teams track PDK versions and re-verify."},
	{KP: 3, Q: "Which PDK piece makes LVS possible?",
		A:   "Device recognition rules that extract transistors and their parameters from layout layers",
		D:   []string{"The timing library", "The mask writer format", "The package model"},
		Why: "LVS extracts a netlist from layout by recognizing devices from layer combinations; those rules come with the PDK."},
	{KP: 3, Q: "A PDK's models were fitted to an early version of the process. What should a designer expect from later releases?",
		A:   "Model updates as the process matures, which can shift simulated performance",
		D:   []string{"Models never change after release", "Only the design rules change", "Models change only for analog devices"},
		Why: "Early models are pre-silicon or from immature data; foundries re-fit them as the process stabilizes."},
	{KP: 3, Q: "Why must a tool designer treat PDK data as configuration rather than code?",
		A:   "The same engine must run many processes, so process facts come in through files it reads",
		D:   []string{"PDK data is always hard-coded", "Tools only support one process", "PDK files are executables"},
		Why: "Engines encode algorithms; technology files, rule decks and models describe the process and change per node and version."},
}

var drmItems = []edakit.Item{
	{KP: 1, Q: "A rule says spacing between two metal lines must be 32 nm if the parallel run length exceeds 100 nm, else 28 nm. What kind of rule is this?",
		A:   "A conditional, context-dependent spacing rule",
		D:   []string{"A minimum width rule", "An enclosure rule", "A density rule"},
		Why: "Modern decks are full of rules whose value depends on width, run length or neighbourhood, not single numbers."},
	{KP: 1, Q: "An end-of-line rule demands extra space in front of a line end. What process effect motivates it?",
		A:   "Line ends pull back and round in printing, so tip-to-tip and tip-to-side spaces need extra margin",
		D:   []string{"Copper resistivity rises at line ends", "Vias cannot land near line ends", "Line ends carry more current"},
		Why: "Line-end shortening and corner rounding make the printed gap differ from the drawn gap, so EOL spacing is larger."},
	{KP: 1, Q: "What does an antenna rule limit?",
		A:   "The ratio of metal area connected to a gate during processing, to avoid plasma charging damage",
		D:   []string{"The length of a clock net", "The number of vias per net", "The package bond-wire length"},
		Why: "During plasma etch, exposed metal collects charge that can damage a thin gate oxide; antenna ratios bound that charge."},
	{KP: 1, Q: "A rule requires metal density between 20% and 80% in every window. Which process step needs it?",
		A:   "Chemical-mechanical polishing",
		D:   []string{"Ion implantation", "Thermal oxidation", "Wafer dicing"},
		Why: "CMP removal depends on local pattern density; density limits keep dishing and erosion within the depth of focus budget."},
	{KP: 2, Q: "How does a DRC runset typically compute a spacing check?",
		A:   "By boolean and sizing operations on layer polygons, then measuring edge-to-edge distances",
		D:   []string{"By simulating the circuit", "By routing a test net", "By counting cells"},
		Why: "DRC languages are geometric algebras: derive layers, select edges, measure, and output violation markers."},
	{KP: 2, Q: "A DRC deck flags a gate-poly extension violation. Which layers did the check combine?",
		A:   "Poly and active, measuring how far poly extends past the active edge",
		D:   []string{"Metal 1 and metal 2", "Via and metal only", "N-well and pad"},
		Why: "Gate extension (endcap) rules measure poly beyond active so the transistor stays fully formed despite line-end pull-back."},
	{KP: 2, Q: "A manual describes a rule in a figure with words; the deck implements it in code. A violation appears in a case the figure does not show. What should a tool developer conclude?",
		A:   "The deck is the signoff authority, but the mismatch should be reported to the foundry",
		D:   []string{"The manual is always wrong", "Ignore the violation", "Change the deck locally and ship"},
		Why: "Signoff uses the deck. A disagreement with the manual is a deck or manual bug that the foundry must resolve, not a local edit."},
	{KP: 2, Q: "Why are multi-patterning colour rules included in the deck at nodes with LELE metal?",
		A:   "Same-colour spacing is larger than different-colour spacing, so the check must know mask assignment",
		D:   []string{"Colour changes copper resistivity", "Colours are cosmetic", "Coloured layers are only for documentation"},
		Why: "Two shapes on different masks can be closer than two on the same mask; the deck checks both spacings and odd cycles."},
	{KP: 3, Q: "Recommended rules in a DRM are not mandatory. Why do yield-aware flows still check them?",
		A:   "Following them improves yield or variability, so tools score and fix them when area allows",
		D:   []string{"They are required for LVS", "They are only for analog", "They are ignored by all tools"},
		Why: "Recommended rules (wider spaces, double vias, longer extensions) lower defect sensitivity; DFM flows apply them opportunistically."},
	{KP: 3, Q: "Rule counts in the manual have grown from hundreds to thousands across nodes. What is the main consequence for routers?",
		A:   "Routers must model many conditional rules internally or rely on DRC repair loops",
		D:   []string{"Routers can ignore all rules", "Routers now run DRC only after tapeout", "Rule count does not affect routing"},
		Why: "A router that does not model a rule creates violations that a later fixer must repair; correct-by-construction routing needs the rules in its cost model."},
	{KP: 3, Q: "A unidirectional, gridded metal layer drops many two-dimensional jog rules. What is the trade-off?",
		A:   "Simpler rules and better printability, at the cost of routing flexibility",
		D:   []string{"Lower density with no benefit", "Higher resistance with no benefit", "No trade-off at all"},
		Why: "Restricting shapes to one direction on a grid simplifies lithography and rules but forbids jogs, so routers need more vias and layer changes."},
	{KP: 3, Q: "Why are some rules expressed only as pattern libraries of forbidden layouts rather than as dimensions?",
		A:   "Some lithography hotspots depend on two-dimensional configurations that no single dimension rule captures",
		D:   []string{"Pattern rules are faster to write", "Dimensions cannot be measured", "Pattern rules replace LVS"},
		Why: "Hotspots arise from specific shapes in combination; pattern matching catches them where width and space rules cannot."},
}

var deviceModelItems = []edakit.Item{
	{KP: 1, Q: "What does an SS corner model in a PDK represent?",
		A:   "Slow NMOS and slow PMOS, near the slow edge of the process distribution",
		D:   []string{"Standard supply and standard temperature", "Small-signal analysis only", "The single-sided version of the cell"},
		Why: "Corner letters name NMOS then PMOS speed: SS, FF, SF, FS and TT bound the global process spread."},
	{KP: 1, Q: "Why do PDKs include skewed corners such as SF and FS?",
		A:   "Some circuits fail when NMOS and PMOS move in opposite directions, such as ratioed logic and SRAM",
		D:   []string{"They are needed for DRC", "They model packaging only", "They set the metal resistance"},
		Why: "Global variation of NMOS and PMOS is not perfectly correlated; skewed corners stress circuits sensitive to their ratio."},
	{KP: 1, Q: "Which model family is the usual foundation of PDK transistor models for FinFETs?",
		A:   "BSIM-CMG",
		D:   []string{"The Elmore model", "The NLDM table", "The Ebers-Moll model"},
		Why: "BSIM-CMG is the industry standard compact model for multi-gate devices."},
	{KP: 1, Q: "A PDK models global and local variation separately. What is local variation?",
		A:   "Random device-to-device mismatch between nearby identical transistors",
		D:   []string{"Lot-to-lot drift", "Wafer-to-wafer thickness change", "Temperature variation across the package"},
		Why: "Local, or mismatch, variation comes from effects such as random dopant fluctuation and line-edge roughness; it does not average out between neighbours."},
	{KP: 2, Q: "Why do statistical models matter beyond corners?",
		A:   "Corners stack worst cases and can be pessimistic or miss combinations; Monte Carlo samples the real distribution",
		D:   []string{"Corners are always optimistic", "Statistical models replace DRC", "Monte Carlo is only for yield of masks"},
		Why: "A corner is one point; statistical models give distributions so circuits can be designed to a target sigma."},
	{KP: 2, Q: "Pelgrom's law says mismatch standard deviation scales how with device area?",
		A:   "As one over the square root of the gate area",
		D:   []string{"Linearly with area", "As the square of area", "It does not depend on area"},
		Why: "Averaging random microscopic fluctuations over a larger area reduces sigma as $1/\\sqrt{WL}$."},
	{KP: 2, Q: "A Liberty library for the SS corner was characterized from which PDK data?",
		A:   "SPICE simulations of each cell with the SS models at the corner's voltage and temperature",
		D:   []string{"The DRC deck", "Measured package delays", "The layer map"},
		Why: "Library characterization runs the PDK models over input slews and loads to fill timing and power tables."},
	{KP: 2, Q: "Layout-dependent effects such as well proximity and stress change a transistor's threshold. How do PDK models capture them?",
		A:   "Instance parameters extracted from layout, such as distances to well edges and diffusion edges",
		D:   []string{"They cannot be modelled", "Through the timing constraints", "By a global corner shift only"},
		Why: "Extraction measures the geometry around each device and passes it as model parameters, so simulation sees layout context."},
	{KP: 3, Q: "A PDK's FF corner shows lower delay but much higher leakage than TT. Which design check uses FF at high temperature?",
		A:   "Leakage power and hold-time checks",
		D:   []string{"Setup time only", "Antenna rules", "Mask counts"},
		Why: "Fast devices leak more and make short paths faster, so FF corners bound leakage and hold."},
	{KP: 3, Q: "Why do reliability models (aging) appear in modern PDKs?",
		A:   "Threshold shifts from BTI and hot carriers over the lifetime must be included in timing margins",
		D:   []string{"Aging affects only masks", "They replace corners", "They are needed for LVS"},
		Why: "End-of-life models let designers simulate degraded devices and set guard-bands."},
	{KP: 3, Q: "Pre-silicon models for a new node have a known uncertainty. How do design teams typically respond?",
		A:   "Add margin and plan for model updates, re-verifying with each PDK release",
		D:   []string{"Ignore the uncertainty", "Freeze the design on the first model", "Use models from an older node"},
		Why: "Early models come from projections and test structures; margins absorb their error until silicon data tightens them."},
	{KP: 3, Q: "Why are model corners and metal (RC) corners separate in a PDK?",
		A:   "Device speed and interconnect parasitics vary from different process steps and are not fully correlated",
		D:   []string{"Metal never varies", "Device and metal corners are identical", "RC corners only apply to analog"},
		Why: "FEOL and BEOL variations come from different steps, so timing signoff combines device corners with RC corners such as Cworst and RCworst."},
}
