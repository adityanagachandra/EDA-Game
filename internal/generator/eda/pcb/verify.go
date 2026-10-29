package pcb

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type rect struct {
	x1, y1, x2, y2 int
	net            string
}

func (r rect) text() string {
	return edakit.Sprintf("(%d, %d) to (%d, %d)", r.x1, r.y1, r.x2, r.y2)
}

// gap is the edge-to-edge distance between two axis-aligned rectangles, zero if they touch.
func gap(a, b rect) float64 {
	dx := max(0, b.x1-a.x2, a.x1-b.x2)
	dy := max(0, b.y1-a.y2, a.y1-b.y2)
	return sqrt(float64(dx*dx + dy*dy))
}

func randRect(net string) rect {
	x, y := edakit.Int(0, 80), edakit.Int(0, 80)
	return rect{x, y, x + edakit.Int(4, 30), y + edakit.Int(4, 30), net}
}

// boardDRC: axis gap (KP 1), corner-to-corner gap (KP 2), counting violations (KP 3).
func boardDRC(ctx generator.GeneratorContext) generator.Problem {
	if kp := kpOf(ctx); kp <= 2 {
		diag := kp == 2
		a := rect{edakit.Int(0, 20), edakit.Int(0, 20), 0, 0, "A"}
		a.x2, a.y2 = a.x1+edakit.Int(10, 40), a.y1+edakit.Int(10, 40)
		var b rect
		if diag {
			b.x1, b.y1 = a.x2+edakit.Int(2, 12), a.y2+edakit.Int(2, 12)
		} else {
			b.x1, b.y1 = a.x2+edakit.Int(2, 15), a.y1+edakit.Int(-5, 8)
		}
		b.x2, b.y2 = b.x1+edakit.Int(8, 30), b.y1+edakit.Int(8, 30)
		g := gap(a, b)
		ans := edakit.Num(g, 2)
		q := "Board DRC measures clearance between copper shapes as the shortest edge-to-edge distance. A pad on net A covers the rectangle " + a.text() + " and a trace segment on net B covers " + b.text() + ", in mil. What is their clearance, in mil? Round to 2 decimal places."
		dx := max(0, b.x1-a.x2, a.x1-b.x2)
		dy := max(0, b.y1-a.y2, a.y1-b.y2)
		exp := edakit.Sprintf("The x gap is %d mil and the y gap is %d mil. The shapes overlap in one axis, so the clearance is the other gap, %s mil.", dx, dy, ans)
		if dx > 0 && dy > 0 {
			exp = edakit.Sprintf("The x gap is %d mil and the y gap is %d mil. Both are positive, so the closest points are two corners and the clearance is $\\sqrt{%d^2 + %d^2}$ = %s mil, not the larger axis gap.", dx, dy, dx, dy, ans)
		}
		return edakit.Numeric(q, ans, exp, map[string]string{"dx": edakit.IntStr(dx), "dy": edakit.IntStr(dy)})
	}
	n := edakit.Scaled(ctx, 4, 6)
	nets := []string{"GND", "VCC", "CLK", "D0", "D1", "RST"}
	var rs []rect
	for i := 0; i < n; i++ {
		rs = append(rs, randRect(edakit.Pick(nets[:min(len(nets), n-1)])))
	}
	c := float64(edakit.Pick([]int{5, 6, 8, 10}))
	bad := 0
	var pairs []string
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rs[i].net == rs[j].net {
				continue
			}
			if gap(rs[i], rs[j]) < c-1e-9 {
				bad++
				pairs = append(pairs, edakit.Sprintf("S%d-S%d", i+1, j+1))
			}
		}
	}
	var desc []string
	for i, r := range rs {
		desc = append(desc, edakit.Sprintf("S%d on %s covers %s", i+1, r.net, r.text()))
	}
	q := edakit.Sprintf("A clearance rule requires at least %s mil, edge to edge, between copper on different nets; shapes on the same net are exempt, and overlapping shapes on different nets count as a violation. On one layer: %s (coordinates in mil). How many pairs of shapes violate the rule?", edakit.Num(c, 0), strings.Join(desc, "; "))
	exp := "Check each pair on different nets, using the axis gaps and the corner distance when both gaps are positive."
	if bad > 0 {
		exp += " Violations: " + strings.Join(pairs, ", ") + "."
	} else {
		exp += " No pair is closer than the rule."
	}
	exp += " A production DRC finds candidate pairs with a spatial index instead of testing all pairs."
	return edakit.Numeric(q, edakit.IntStr(bad), exp, map[string]string{"violations": edakit.IntStr(bad)})
}

// boardDFT: test access coverage (KP 1), JTAG instruction length (KP 2), bypass chain length (KP 3).
func boardDFT(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		N := edakit.Int(20, 120) * 10
		tp := edakit.Int(N/3, N*8/10)
		bs := edakit.Int(0, (N-tp)*8/10)
		ans := edakit.Num(100*float64(tp+bs)/float64(N), 1)
		q := edakit.Sprintf("A board has %d nets. %d nets have a probe-able test point for in-circuit or flying-probe test, and %d other nets have no test point but are reachable through boundary-scan cells. What percentage of nets has some test access? Round to 1 decimal place.", N, tp, bs)
		return edakit.Numeric(q, ans, edakit.Sprintf("$(%d + %d)/%d = %s\\%%$. Boundary scan recovers access to nets under BGAs that no probe can reach.", tp, bs, N, ans), nil)
	case 2:
		n := edakit.Scaled(ctx, 3, 6)
		var irs []string
		sum := 0
		for i := 0; i < n; i++ {
			ir := edakit.Pick([]int{2, 4, 5, 6, 8, 10, 14})
			sum += ir
			irs = append(irs, edakit.Sprintf("U%d with a %d-bit instruction register", i+1, ir))
		}
		q := "A JTAG chain links " + joinAnd(irs) + " from TDI to TDO. How many bits must be shifted in Shift-IR to load an instruction into every device at once?"
		return edakit.Numeric(q, edakit.IntStr(sum), edakit.Sprintf("The instruction registers are in series, so the scan length is their sum, %d bits.", sum), nil)
	}
	n := edakit.Scaled(ctx, 3, 7)
	bsr := edakit.Int(40, 900)
	q := edakit.Sprintf("A JTAG chain has %d devices. To run EXTEST on one FPGA with a %d-cell boundary-scan register, every other device is put in BYPASS, which selects its 1-bit bypass register. How many TCK cycles in Shift-DR shift one complete test pattern through the chain?", n, bsr)
	return edakit.Numeric(q, edakit.IntStr(bsr+n-1), edakit.Sprintf("The data path is the target's %d cells plus one bit for each of the %d bypassed devices: %d cycles.", bsr, n-1, bsr+n-1), nil)
}

var bringupItems = []edakit.Item{
	{KP: 1, Q: "Before a new board is powered for the first time, what should be measured?",
		A:   "The resistance from each supply rail to ground, to catch shorts",
		D:   []string{"The trace impedance of every net", "The clock frequency", "The firmware checksum"},
		Why: "A rail-to-ground short from a solder bridge or a rotated part is common on first articles and is cheap to find unpowered."},
	{KP: 1, Q: "Why power a new board from a current-limited bench supply at first?",
		A:   "A short or misplaced part then draws limited current instead of burning traces and parts",
		D:   []string{"It makes the clocks start faster", "It is required for JTAG", "It raises the rail voltage for margin"},
		Why: "Current limiting turns a destructive fault into a visible symptom, the supply folding back."},
	{KP: 1, Q: "Which order is sensible for first bring-up?",
		A:   "Rails and sequencing, then clocks and resets, then JTAG access, then firmware and interfaces",
		D:   []string{"Firmware first, then rails, then clocks", "Interfaces first, then rails, then JTAG", "JTAG first, before any rail is powered"},
		Why: "Each step depends on the one before. Debugging an interface on a board with a bad rail wastes time."},
	{KP: 1, Q: "Reading every device's IDCODE over JTAG is often the first digital test. What does it show?",
		A:   "The chain is wired correctly, the devices are powered and alive, and the expected parts are fitted",
		D:   []string{"All DDR timing is met", "Every net is free of crosstalk", "The board meets its EMC limits"},
		Why: "A wrong or missing IDCODE localizes chain breaks, unpowered devices or substituted parts quickly."},
	{KP: 2, Q: "A rail is shorted to ground with no visible cause. Which technique localizes the short?",
		A:   "Inject a limited current and find the hot spot with a thermal camera or follow the voltage-drop gradient with a milliohm meter",
		D:   []string{"Raise the supply voltage until the part fails open", "Re-run the layout DRC", "Probe the clock with an oscilloscope"},
		Why: "The shorted spot dissipates the power, and the voltage drop along the plane points toward it."},
	{KP: 2, Q: "A scope probe with a long ground lead shows heavy ringing on a fast edge. What is the likely cause?",
		A:   "The ground lead's inductance resonates with the probe's capacitance; a short ground spring removes it",
		D:   []string{"The trace is badly mismatched", "The scope's bandwidth is too high", "The power rail is noisy"},
		Why: "Measurement artifacts must be ruled out before redesigning termination."},
	{KP: 2, Q: "A TDR trace of a net shows a dip in impedance at one distance. What does it indicate?",
		A:   "A capacitive discontinuity there, such as a via, a pad or a wider section",
		D:   []string{"An inductive discontinuity there", "An open circuit at that point", "A perfectly matched section"},
		Why: "Excess capacitance lowers local impedance; excess inductance, such as a neck-down or a plane gap, raises it."},
	{KP: 2, Q: "An I2C bus fails intermittently and its rising edges are slow. What is the first suspect?",
		A: "Pull-up resistors too large for the bus capacitance", D: []string{"Termination resistors missing at the far end", "A ground plane that is too solid", "Pull-ups that are too small"},
		Why: "Open-drain rise time is set by the pull-up resistance times bus capacitance. Too large a resistor violates the rise-time limit."},
	{KP: 3, Q: "A fix found at bring-up is a bodge wire. Where must it go next?",
		A:   "Back into the schematic and layout as an ECO for the next revision, with the reason recorded",
		D:   []string{"Only into the bring-up notebook", "Into the Gerber files of the current revision", "Nowhere, since the board works now"},
		Why: "Unrecorded rework is lost knowledge and guarantees the bug on the next build."},
	{KP: 3, Q: "A DDR interface fails only at high temperature. What does this point to?",
		A:   "Marginal timing or signal integrity that the cold condition hides; check margins at corners",
		D:   []string{"A schematic pinout error", "A missing BOM line", "A wrong drill size"},
		Why: "A hard connectivity error fails everywhere. Temperature-dependent failures are margin problems in timing, SI or power."},
	{KP: 3, Q: "Why add test points and a debug header during design rather than at bring-up?",
		A:   "Bring-up needs probe access to rails, clocks and buses that BGAs and inner layers otherwise hide",
		D:   []string{"Test points raise trace impedance", "Headers are needed for fabrication", "They replace the netlist"},
		Why: "Once the board is built, a net with no access can only be reached by rework."},
	{KP: 3, Q: "A board works on the bench but fails inside its enclosure. What is a plausible cause?",
		A:   "Thermal rise, EMI coupling or chassis grounding that the bench setup did not reproduce",
		D:   []string{"A footprint pin-numbering error", "A missing net in the netlist", "A wrong reference designator"},
		Why: "Environment-dependent failures point to physics the bench lacks, not to static connectivity errors."},
}

// gerberCoord is the computed KP 2 of output.gerber_odb: decode two coordinates and measure a draw.
func gerberCoord(ctx generator.GeneratorContext) generator.Problem {
	x1 := edakit.Int(1000, 90000) * 100
	y1 := edakit.Int(1000, 90000) * 100
	horiz := edakit.Int(1, 2) == 1
	d := edakit.Int(500, 40000) * 100
	if edakit.Int(1, 2) == 1 {
		d = -d
	}
	x2, y2 := x1, y1
	if horiz {
		x2 = x1 + d
	} else {
		y2 = y1 + d
	}
	length := absf(float64(d)) / 1e6
	sum := (absf(float64(x1)) + absf(float64(x2))) / 1e6
	if !horiz {
		sum = (absf(float64(y1)) + absf(float64(y2))) / 1e6
	}
	ans, ds := options(length, 4, "mm", length*100, length/10, length*10, sum)
	q := edakit.Sprintf("A Gerber file declares %%FSLAX46Y46*%% (absolute coordinates, leading zeros omitted, 4 integer and 6 decimal digits) and %%MOMM*%% (millimetres). It contains X%dY%dD02* followed by X%dY%dD01*. How long is the drawn segment?", x1, y1, x2, y2)
	return edakit.Choice(q, ans, edakit.Sprintf("With 6 decimal digits each coordinate is the integer divided by $10^6$. D02 moves to the first point and D01 draws to the second, so the length is $|%d| / 10^6$, which is %s.", d, ans), ds...)
}

var gerberItems = []edakit.Item{
	{KP: 1, Q: "In an RS-274X Gerber file, what does a D03 operation do?",
		A:   "Flashes the current aperture at the coordinate, creating a pad image",
		D:   []string{"Draws a line from the current point with the current aperture", "Moves the current point without drawing", "Ends the file"},
		Why: "D01 draws (interpolates), D02 moves, D03 flashes. Pads are usually flashes, traces draws."},
	{KP: 1, Q: "In an RS-274X Gerber file, what does a D01 operation do?",
		A:   "Draws from the current point to the new coordinate with the current aperture",
		D:   []string{"Flashes the aperture at the coordinate", "Moves without drawing", "Selects aperture number 1"},
		Why: "D01 is interpolation, linear or circular depending on the G01, G02 or G03 mode."},
	{KP: 1, Q: "What does the statement %ADD10C,0.25*% define?",
		A: "Aperture D10 as a circle 0.25 units in diameter", D: []string{"A circle drawn at coordinates (10, 0.25)", "A coordinate format of 10 and 25 digits", "A selection of layer 10"},
		Why: "AD defines an aperture: a number from 10 upward, a template (C, R, O, P or a macro) and its parameters."},
	{KP: 1, Q: "What is Gerber's basic model of a board layer?",
		A:   "A 2D vector image built from flashes, draws and filled regions with dark and clear polarity, one file per layer",
		D:   []string{"A hierarchical cell database like GDSII", "A netlist with attached geometry", "A 3D solid model of the board"},
		Why: "Gerber describes what to image on one layer. Connectivity and stackup are not part of the classic image model."},
	{KP: 1, Q: "How does ODB++ differ from a set of Gerber files?",
		A:   "It is one structured database holding layers, drills, netlist, components and stackup",
		D:   []string{"It is a single Gerber layer with a new extension", "It carries the netlist and nothing else", "It is a schematic exchange format"},
		Why: "A unified job database removes the need to reassemble layer order, drill pairing and netlist from separate files."},
	{KP: 3, Q: "Gerber X2 adds attributes such as .FileFunction and .AperFunction. Why does a fabricator's CAM tool care?",
		A:   "It can identify each file's layer and each object's role (pad, via, conductor) without guessing from file names",
		D:   []string{"It makes the files smaller", "It replaces the drill files", "It encrypts the design"},
		Why: "Attributes turn an image into self-describing data, so CAM import is automatic and less error-prone."},
	{KP: 3, Q: "A CAM engineer receives plain RS-274X layers and Excellon drill files. What must be reconstructed from other documents?",
		A:   "Layer order and stackup from the fabrication drawing, and connectivity for test from an IPC-D-356 netlist",
		D:   []string{"Nothing, since the images contain everything", "Only the silkscreen", "Only the aperture list"},
		Why: "Plain Gerber has no stackup or netlist, so the job relies on side documents that must match the images."},
	{KP: 3, Q: "Which file carries the drill program in a Gerber-based release?",
		A: "An Excellon (NC drill) file, one per drill span", D: []string{"The paste layer Gerber", "The silkscreen Gerber", "The pick-and-place centroid file"},
		Why: "Drill data is a separate numerical-control format; blind and buried spans each get their own file."},
	{KP: 3, Q: "A design is released without the board outline layer. What happens?",
		A:   "The fabricator cannot cut the board profile and must stop to ask, delaying the order",
		D:   []string{"The fab uses the copper extents automatically and nothing is lost", "The drill file supplies the outline", "The netlist supplies the outline"},
		Why: "The profile defines the routed edge and the panel layout, so it is a required layer."},
	{KP: 3, Q: "Why does a fabricator run its own DFM checks on received data even after the designer's DRC passed?",
		A:   "Its real process limits and panel rules may be stricter or different from the rule set the designer used",
		D:   []string{"Designer DRC never checks clearance", "Gerber files cannot be checked by designers", "Fab DFM only checks the silkscreen"},
		Why: "The designer's rule deck is a model of the fab. The fab checks against the process itself."},
	{KP: 3, Q: "What is IPC-2581 in the context of manufacturing outputs?",
		A:   "An open XML format that carries layers, stackup, drills, netlist and BOM in one file",
		D:   []string{"A Gerber aperture macro", "A land-pattern standard", "A soldering acceptance standard"},
		Why: "IPC-2581 plays the same role as ODB++, as a single neutral database instead of a collection of images."},
}

var ipcItems = []edakit.Item{
	{KP: 1, Q: "Which IPC standard is the generic design standard for printed boards, with spacing and conductor tables?",
		A: "IPC-2221", D: []string{"IPC-A-610", "IPC-7351", "J-STD-001"},
		Why: "IPC-2221 is the generic design standard; the IPC-2222 to 2226 sectionals add rigid, flex and HDI detail."},
	{KP: 1, Q: "Which IPC standard covers surface-mount land pattern design?",
		A: "IPC-7351", D: []string{"IPC-2221", "IPC-6012", "IPC-A-600"},
		Why: "IPC-7351 gives the land-pattern calculation and naming convention for SMT footprints."},
	{KP: 1, Q: "Which IPC standard gives visual acceptance criteria for electronic assemblies, such as solder fillets?",
		A: "IPC-A-610", D: []string{"IPC-A-600", "IPC-2152", "IPC-2581"},
		Why: "IPC-A-610 is the assembly acceptability standard used by inspectors."},
	{KP: 1, Q: "Which IPC standard gives acceptance criteria for bare printed boards?",
		A: "IPC-A-600", D: []string{"IPC-A-610", "IPC-7351", "IPC-2226"},
		Why: "IPC-A-600 shows acceptable and rejectable conditions of the bare board, such as plating voids and annular rings."},
	{KP: 2, Q: "Which IPC specification covers the qualification and performance of rigid printed boards?",
		A: "IPC-6012", D: []string{"IPC-2223", "IPC-A-610", "IPC-D-356"},
		Why: "IPC-6012 is the performance spec a rigid board is built and tested to, by class."},
	{KP: 2, Q: "Which standard sets the requirements for soldered electrical and electronic assemblies?",
		A: "J-STD-001", D: []string{"IPC-7351", "IPC-2152", "IPC-2581"},
		Why: "J-STD-001 is the soldering process and materials requirement, paired with IPC-A-610 for inspection."},
	{KP: 2, Q: "Which IPC standard addresses conductor current-carrying capacity with newer test data than the IPC-2221 chart?",
		A: "IPC-2152", D: []string{"IPC-2141", "IPC-6012", "IPC-A-600"},
		Why: "IPC-2152 accounts for board thickness, planes and other factors that the older single chart ignores."},
	{KP: 2, Q: "Which IPC sectional design standard covers HDI boards and microvias?",
		A: "IPC-2226", D: []string{"IPC-2223", "IPC-2222", "IPC-6012"},
		Why: "IPC-2226 covers high-density interconnect structures; IPC-2223 covers flex."},
	{KP: 2, Q: "Which IPC sectional design standard covers flexible and rigid-flex boards?",
		A: "IPC-2223", D: []string{"IPC-2226", "IPC-2221", "IPC-7351"},
		Why: "IPC-2223 gives flex-specific rules such as bend radius guidance and coverlay design."},
	{KP: 3, Q: "What do IPC product Classes 1, 2 and 3 represent?",
		A:   "Rising reliability requirements: general electronics, dedicated service, and high-reliability products",
		D:   []string{"Ranges of layer count", "Ranges of copper weight", "Ranges of board size"},
		Why: "The class sets acceptance criteria across the design, fabrication and assembly standards."},
	{KP: 3, Q: "A life-critical medical board is specified as Class 3. What changes for the fabricator?",
		A:   "Tighter acceptance criteria, such as minimum annular ring and plating thickness, and more inspection",
		D:   []string{"Nothing, since class affects only assembly", "Only the silkscreen color", "Only the panel size"},
		Why: "Class 3 shrinks the allowed defects and raises minimums, which raises cost and lowers yield."},
	{KP: 3, Q: "Why should a board DRC rule deck be derived from IPC tables and the product class?",
		A:   "Rules tied to standard tables are portable across fabricators and auditable against the class the product needs",
		D:   []string{"IPC tables make DRC run faster", "Fabricators ignore all other rules", "DRC cannot check anything else"},
		Why: "A shared reference lets the designer, the fab and the customer agree on what a passing board means."},
	{KP: 3, Q: "What is IPC-D-356?",
		A: "A netlist format for electrical testing of bare boards", D: []string{"A land-pattern naming convention", "A soldering process standard", "A stackup material specification"},
		Why: "It lists nets and test-point locations so flying-probe or fixture testers can check for opens and shorts."},
}
