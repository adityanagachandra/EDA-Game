package flow

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type format struct {
	name, keeps, producer, need string
}

var formats = []format{
	{"RTL (Verilog or VHDL)", "cycle-by-cycle register behaviour, with no gates and no geometry", "The designer, during RTL design", "simulate the design's behaviour cycle by cycle before synthesis"},
	{"Gate-level netlist", "library cell instances and the nets between their pins, with no geometry", "Logic synthesis", "list every cell instance and which pins each net connects"},
	{"Liberty (.lib)", "per-cell delay, slew and power tables indexed by input slew and output load", "Library characterization", "look up the delay of a NAND2 arc at a given input slew and output load"},
	{"SDC", "clock periods, input and output delays and timing exceptions", "The designer, as constraint authoring", "know the clock period that every setup check is measured against"},
	{"LEF", "cell outlines, pin shapes, obstructions and routing-layer rules, with no transistors", "The library and technology provider", "know a standard cell's width and where its pins can be reached by a router"},
	{"DEF", "one design's die area, rows, placed components and routed nets", "Floorplanning, placement and routing tools", "read the placed coordinates of every cell instance in this design"},
	{"SPEF", "extracted resistance and capacitance of every routed net, including coupling", "Parasitic extraction", "obtain the coupling capacitance between two routed nets for crosstalk delay"},
	{"GDSII or OASIS", "polygons on numbered mask layers, with no connectivity", "Layout stream-out after routing and fill", "measure the exact polygons that a design-rule checker checks"},
	{"SDF", "computed delays of every cell arc and interconnect arc", "Delay calculation in the timing engine", "back-annotate delays into a gate-level timing simulation"},
	{"UPF", "power domains, supply nets, isolation and retention intent", "The designer, as power-intent authoring", "know which blocks may power off and where isolation cells belong"},
	{"SPICE netlist", "transistors and passive devices with their model parameters", "Schematic capture, or layout extraction for LVS", "run a transistor-level transient simulation of a cell"},
	{"VCD or SAIF activity file", "signal transitions or toggle counts recorded from a simulation", "Logic simulation", "use measured switching activity in a dynamic power estimate"},
}

// formatMap: what a format keeps (KP 1), who writes it (KP 2), what a computation must read (KP 3).
func formatMap(ctx generator.GeneratorContext) generator.Problem {
	i := edakit.Int(0, len(formats)-1)
	f := formats[i]
	var others []format
	for _, k := range edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}) {
		if k != i {
			others = append(others, formats[k])
		}
		if len(others) == 3 {
			break
		}
	}
	switch kp(ctx, 3) {
	case 1:
		q := "Which format keeps " + f.keeps + "?"
		return edakit.Choice(q, f.name, f.name+" keeps "+f.keeps+".", others[0].name, others[1].name, others[2].name)
	case 2:
		q := "Which stage or party normally writes the " + f.name + " file?"
		return edakit.Choice(q, f.producer, f.name+" is written by: "+lowerFirst(f.producer)+".", others[0].producer, others[1].producer, others[2].producer)
	default:
		q := "A tool needs to " + f.need + ". Which file holds that information?"
		return edakit.Choice(q, f.name, f.name+" keeps "+f.keeps+", which is what this computation needs.", others[0].name, others[1].name, others[2].name)
	}
}

// netlist: total fanout, total pin connections, and clique or star edge counts.
func netlist(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 8)
	var lines []string
	sizes := make([]int, n)
	for i := 0; i < n; i++ {
		p := edakit.Int(2, 6)
		sizes[i] = p
		insts := edakit.Shuffle([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})[:p]
		pins := []string{edakit.Sprintf("U%d/Y", insts[0])}
		for _, u := range insts[1:] {
			pins = append(pins, edakit.Sprintf("U%d/%s", u, edakit.Pick([]string{"A", "B", "C"})))
		}
		lines = append(lines, edakit.Sprintf("n%d: %s", i+1, strings.Join(pins, " ")))
	}
	list := "A netlist has these nets, each written with its driver pin first: " + strings.Join(lines, "; ") + "."
	fan, pins, clique, star := 0, 0, 0, 0
	for _, p := range sizes {
		fan += p - 1
		pins += p
		clique += p * (p - 1) / 2
		star += p
	}
	switch kp(ctx, 3) {
	case 1:
		q := list + " What is the total fanout, summed over all nets?"
		return edakit.Numeric(q, edakit.IntStr(fan),
			edakit.Sprintf("Each net's fanout is its pin count minus the driver; the sum is %d.", fan), map[string]string{"fanout": edakit.IntStr(fan)})
	case 2:
		q := list + " Viewing the netlist as a hypergraph with one hyperedge per net, what is the sum of the hyperedge sizes (the total number of pin connections)?"
		return edakit.Numeric(q, edakit.IntStr(pins),
			edakit.Sprintf("Sum of pin counts over the nets: %d.", pins), map[string]string{"pins": edakit.IntStr(pins)})
	default:
		if edakit.Pick([]bool{true, false}) {
			q := list + " A placer converts every net to a clique, joining every pair of pins on it. How many graph edges result?"
			return edakit.Numeric(q, edakit.IntStr(clique),
				edakit.Sprintf("A p-pin net gives $p(p-1)/2$ edges; summed over the nets that is %d.", clique), map[string]string{"clique_edges": edakit.IntStr(clique)})
		}
		q := list + " A placer converts every net to a star, adding one new vertex per net joined to each of its pins. How many graph edges result?"
		return edakit.Numeric(q, edakit.IntStr(star),
			edakit.Sprintf("A p-pin net gives p star edges; summed over the nets that is %d.", star), map[string]string{"star_edges": edakit.IntStr(star)})
	}
}

// booleanNetworks: truth-table sizes, AIG node counts and BDD node counts.
func booleanNetworks(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		if edakit.Pick([]bool{true, false}) {
			n := edakit.Scaled(ctx, 3, 16)
			ans := 1 << n
			q := edakit.Sprintf("How many rows does the truth table of a %d-input Boolean function have?", n)
			return edakit.Numeric(q, edakit.IntStr(ans), edakit.Sprintf("$2^{%d} = %d$ rows.", n, ans), nil)
		}
		n := edakit.Int(1, 4)
		ans := 1 << (1 << n)
		q := edakit.Sprintf("How many distinct Boolean functions of %d inputs exist?", n)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Each of the $2^{%d} = %d$ rows can be 0 or 1, so there are $2^{%d} = %d$ functions.", n, 1<<n, 1<<n, ans), nil)
	case 2:
		k := edakit.Scaled(ctx, 2, 4)
		total := 0
		var parts []string
		for i := 0; i < k; i++ {
			switch edakit.Int(0, 3) {
			case 0:
				m := edakit.Int(2, 8)
				total += m - 1
				parts = append(parts, edakit.Sprintf("a %d-input AND", m))
			case 1:
				m := edakit.Int(2, 8)
				total += m - 1
				parts = append(parts, edakit.Sprintf("a %d-input OR", m))
			case 2:
				total += 3
				parts = append(parts, "a two-input XOR")
			default:
				total += 3
				parts = append(parts, "a two-to-one multiplexer")
			}
		}
		q := "A logic network consists of " + strings.Join(parts, ", ") + ", with no logic shared between them. How many two-input AND nodes does its and-inverter graph need, with inverters free on edges?"
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("An m-input AND or OR needs m - 1 nodes (OR by De Morgan); an XOR or a multiplexer needs 3. The total is %d.", total), map[string]string{"aig_nodes": edakit.IntStr(total)})
	default:
		n := edakit.Scaled(ctx, 3, 12)
		switch edakit.Int(0, 2) {
		case 0:
			q := edakit.Sprintf("How many internal (non-terminal) nodes does the reduced ordered BDD of the %d-input AND function have, under any variable order?", n)
			return edakit.Numeric(q, edakit.IntStr(n), edakit.Sprintf("Each variable needs one node that sends 0 to the 0-terminal: %d nodes.", n), nil)
		case 1:
			q := edakit.Sprintf("How many internal (non-terminal) nodes does the reduced ordered BDD of the %d-input OR function have, under any variable order?", n)
			return edakit.Numeric(q, edakit.IntStr(n), edakit.Sprintf("Each variable needs one node that sends 1 to the 1-terminal: %d nodes.", n), nil)
		default:
			q := edakit.Sprintf("How many internal (non-terminal) nodes does the reduced ordered BDD (without complement edges) of %d-input parity (XOR of all inputs) have?", n)
			return edakit.Numeric(q, edakit.IntStr(2*n-1),
				edakit.Sprintf("The top level has one node and each of the other %d levels has two, one per parity so far: $2 \\cdot %d - 1 = %d$.", n-1, n, 2*n-1), nil)
		}
	}
}

// liberty: read a table entry, interpolate along load, and sum leakage.
func liberty(ctx generator.GeneratorContext) generator.Problem {
	slews := []float64{0.02, 0.05, 0.1}
	loads := []float64{0.005, 0.01, 0.02, 0.04}
	a := float64(edakit.Int(15, 35)) / 1000
	b := float64(edakit.Int(3, 8)) / 10
	c := float64(edakit.Int(15, 40)) / 10
	val := func(s, l float64) float64 { return edakit.Round(a+b*s+c*l, 4) }
	switch kp(ctx, 3) {
	case 1:
		var rows []string
		for _, s := range slews {
			var cells []string
			for _, l := range loads {
				cells = append(cells, f2(val(s, l), 4))
			}
			rows = append(rows, edakit.Sprintf("slew %s ns: %s", f2(s, 3), strings.Join(cells, ", ")))
		}
		si, li := edakit.Int(0, 2), edakit.Int(0, 3)
		ans := f2(val(slews[si], loads[li]), 4)
		q := "A Liberty cell_rise table has index_1 (input slew) 0.02, 0.05, 0.1 ns and index_2 (output load) 0.005, 0.01, 0.02, 0.04 pF. Its rows are " + strings.Join(rows, "; ") +
			edakit.Sprintf(". What rise delay does it give, in ns, at an input slew of %s ns and a load of %s pF?", f2(slews[si], 3), f2(loads[li], 3))
		return edakit.Numeric(q, ans, "The point lies on the grid, so the delay is the table entry "+ans+" ns.", nil)
	case 2:
		s := edakit.Pick(slews)
		li := edakit.Int(0, 2)
		l0, l1 := loads[li], loads[li+1]
		frac := float64(edakit.Int(1, 9)) / 10
		l := edakit.Round(l0+frac*(l1-l0), 4)
		d0, d1 := val(s, l0), val(s, l1)
		d := d0 + (l-l0)/(l1-l0)*(d1-d0)
		ans := edakit.Num(d, 4)
		q := edakit.Sprintf("At an input slew of %s ns, a Liberty delay table gives %s ns at a load of %s pF and %s ns at %s pF. Using linear interpolation in load, what is the delay at %s pF, in ns? Round to 4 decimal places.",
			f2(s, 3), f2(d0, 4), f2(l0, 4), f2(d1, 4), f2(l1, 4), f2(l, 4))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$%s + \\frac{%s - %s}{%s - %s}(%s - %s) = %s$ ns.", f2(d0, 4), f2(l, 4), f2(l0, 4), f2(l1, 4), f2(l0, 4), f2(d1, 4), f2(d0, 4), ans), nil)
	default:
		k := edakit.Scaled(ctx, 2, 4)
		names := edakit.Shuffle([]string{"INV_X1", "NAND2_X1", "NOR2_X1", "DFF_X1", "AOI21_X1", "BUF_X2"})[:k]
		total := 0.0
		var parts []string
		for _, nm := range names {
			cnt := 100 * edakit.Int(5, 80)
			leak := float64(edakit.Int(5, 120)) / 10
			total += float64(cnt) * leak
			parts = append(parts, edakit.Sprintf("%s instances of %s at %s nW each", comma(cnt), nm, f2(leak, 1)))
		}
		ans := edakit.Num(total/1000, 3)
		q := "A block contains " + strings.Join(parts, ", ") + ", using the Liberty leakage values at one corner. What is its total leakage power, in microwatts? Round to 3 decimal places."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Sum count times leakage: %s nW, which is %s microwatts.", f2(total, 1), ans), map[string]string{"leak_nw": f2(total, 1)})
	}
}

// sdc: frequency from a period, input-path budget and multicycle budget.
func sdc(ctx generator.GeneratorContext) generator.Problem {
	periods := []float64{0.8, 1, 1.25, 1.6, 2, 2.5, 3.2, 4, 5, 6.4, 8, 10}
	p := edakit.Pick(periods)
	switch kp(ctx, 3) {
	case 1:
		ans := edakit.Num(1000/p, 1)
		q := edakit.Sprintf("An SDC file contains create_clock -name core_clk -period %s [get_ports clk], with times in ns. What is the clock frequency, in MHz? Round to 1 decimal place.", f2(p, 2))
		return edakit.Numeric(q, ans, edakit.Sprintf("$1000/%s = %s$ MHz.", f2(p, 2), ans), nil)
	case 2:
		in := edakit.Round(p*float64(edakit.Int(15, 45))/100, 2)
		su := float64(edakit.Int(3, 12)) / 100
		un := float64(edakit.Int(2, 15)) / 100
		ans := edakit.Num(p-in-su-un, 3)
		q := edakit.Sprintf("An SDC file declares a %s ns clock, set_input_delay %s -clock clk on port din and set_clock_uncertainty %s. The first flip-flop after din has a %s ns setup time. How much time, in ns, is left for the logic from din to that flip-flop? Round to 3 decimal places.", f2(p, 2), f2(in, 2), f2(un, 2), f2(su, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$%s - %s - %s - %s = %s$ ns.", f2(p, 2), f2(in, 2), f2(su, 2), f2(un, 2), ans), nil)
	default:
		n := edakit.Int(2, 4)
		su := float64(edakit.Int(3, 12)) / 100
		un := float64(edakit.Int(2, 15)) / 100
		ans := edakit.Num(float64(n)*p-su-un, 3)
		q := edakit.Sprintf("A path between two flip-flops on a %s ns clock carries set_multicycle_path %d -setup. The capture flip-flop has a %s ns setup time and the clock uncertainty is %s ns. How much time, in ns, may the path's logic take? Round to 3 decimal places.", f2(p, 2), n, f2(su, 2), f2(un, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The setup check moves to the edge %d periods after launch: $%d \\cdot %s - %s - %s = %s$ ns.", n, n, f2(p, 2), f2(su, 2), f2(un, 2), ans), nil)
	}
}

// lefDef: die area from DIEAREA, a component coordinate in microns, and row or site counts.
func lefDef(ctx generator.GeneratorContext) generator.Problem {
	u := edakit.Pick([]int{1000, 2000})
	switch kp(ctx, 3) {
	case 1:
		x, y := 100*edakit.Int(5, 40), 100*edakit.Int(5, 40)
		area := float64(x*y) / 1e6
		ans := edakit.Num(area, 4)
		q := edakit.Sprintf("A DEF file declares UNITS DISTANCE MICRONS %d and DIEAREA ( 0 0 ) ( %d %d ). What is the die area, in square millimetres? Round to 4 decimal places.", u, x*u, y*u)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Dividing by %d units per micron gives a %d by %d micron die, %s square millimetres.", u, x, y, ans), nil)
	case 2:
		xu, yu := edakit.Int(1000, 900000), edakit.Int(1000, 900000)
		axis, val := "x", xu
		if edakit.Pick([]bool{true, false}) {
			axis, val = "y", yu
		}
		ans := edakit.Num(float64(val)/float64(u), 3)
		q := edakit.Sprintf("With UNITS DISTANCE MICRONS %d, a DEF component line reads - u%d NAND2_X1 + PLACED ( %d %d ) N. What is the cell's %s coordinate, in microns? Round to 3 decimal places.", u, edakit.Int(1, 9999), xu, yu, axis)
		return edakit.Numeric(q, ans, edakit.Sprintf("$%d / %d = %s$ microns.", val, u, ans), nil)
	default:
		if edakit.Pick([]bool{true, false}) {
			h := edakit.Pick([]float64{1.2, 1.4, 1.8, 2.0, 2.72})
			rows := edakit.Int(100, 900)
			core := float64(rows) * h
			q := edakit.Sprintf("A core is %s microns tall and its LEF site is %s microns high. How many standard-cell rows fit in the core?", f2(core, 2), f2(h, 2))
			return edakit.Numeric(q, edakit.IntStr(rows), edakit.Sprintf("$%s / %s = %d$ rows.", f2(core, 2), f2(h, 2), rows), nil)
		}
		sw := edakit.Pick([]float64{0.19, 0.2, 0.46, 0.054, 0.136})
		w := 10 * edakit.Int(20, 200)
		raw := float64(w) / sw
		if nearInt(raw, 1e-6) {
			raw = math.Round(raw)
		}
		ans := int(math.Floor(raw + 1e-9))
		q := edakit.Sprintf("A placement row is %d microns wide and the site width is %s microns. How many whole sites does the row hold?", w, f2(sw, 3))
		return edakit.Numeric(q, edakit.IntStr(ans), edakit.Sprintf("$\\lfloor %d / %s \\rfloor = %d$ sites.", w, f2(sw, 3), ans), nil)
	}
}

// spef: total capacitance, unit conversion and path resistance.
func spef(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		total := 0.0
		var parts []string
		for i := 0; i < n; i++ {
			v := float64(edakit.Int(1, 40)) / 10
			total += v
			parts = append(parts, edakit.Sprintf("%d *8:%d %s", i+1, i+1, f2(v, 1)))
		}
		cc := float64(edakit.Int(1, 20)) / 10
		total += cc
		parts = append(parts, edakit.Sprintf("%d *8:2 *14:3 %s", n+1, f2(cc, 1)))
		ans := edakit.Num(total, 1)
		q := "A SPEF file with *C_UNIT 1 FF has a *D_NET for net *8 whose *CAP section reads: " + strings.Join(parts, "; ") + ". The last entry couples *8 to another net. What total capacitance should the *D_NET line report, in fF? Round to 1 decimal place."
		return edakit.Numeric(q, ans, "Ground and coupling capacitances add: "+ans+" fF.", nil)
	case 2:
		n := edakit.Scaled(ctx, 2, 5)
		total := 0.0
		var parts []string
		for i := 0; i < n; i++ {
			v := float64(edakit.Int(5, 90)) / 10000
			total += v
			parts = append(parts, f2(v, 4))
		}
		ans := edakit.Num(total*1000, 1)
		q := "A SPEF header declares *C_UNIT 1 PF. A net's *CAP entries to ground are " + strings.Join(parts, ", ") + ". What is the net's total ground capacitance, in fF? Round to 1 decimal place."
		return edakit.Numeric(q, ans, edakit.Sprintf("The sum is %s pF, which is %s fF.", f2(total, 4), ans), nil)
	default:
		r1, r2, r3 := edakit.Int(5, 60), edakit.Int(5, 60), edakit.Int(5, 60)
		rb, rb2 := edakit.Int(5, 60), edakit.Int(5, 60)
		q := edakit.Sprintf("A net's *RES section (with *R_UNIT 1 OHM) reads: 1 *5:D *5:1 %d; 2 *5:1 *5:2 %d; 3 *5:2 *U9:A %d; 4 *5:1 *5:3 %d; 5 *5:3 *U4:B %d. Node *5:D is the driver pin and the resistors form a tree. What is the resistance from the driver to sink pin *U9:A, in ohms?", r1, r2, r3, rb, rb2)
		ans := r1 + r2 + r3
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("The tree path is *5:D, *5:1, *5:2, *U9:A: $%d + %d + %d = %d$ ohms. The branch through *5:3 serves the other sink.", r1, r2, r3, ans), nil)
	}
}

// gdsii: database-unit conversion, coordinate data size, and records saved by repetition.
func gdsii(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		units := []struct {
			m    float64
			name string
			tex  string
		}{{1e-9, "1E-9", "10^{-9}"}, {5e-10, "5E-10", "5 \\times 10^{-10}"}, {1e-8, "1E-8", "10^{-8}"}, {2.5e-10, "2.5E-10", "2.5 \\times 10^{-10}"}}
		uu := edakit.Pick(units)
		c := edakit.Int(1000, 2000000)
		ans := edakit.Num(float64(c)*uu.m*1e6, 4)
		q := edakit.Sprintf("A GDSII UNITS record says one database unit is %s metres. A BOUNDARY vertex has the integer coordinate %d. What is that coordinate in microns? Round to 4 decimal places.", uu.name, c)
		return edakit.Numeric(q, ans, edakit.Sprintf("$%d \\cdot %s$ m is %s microns.", c, uu.tex, ans), nil)
	case 2:
		r := edakit.Scaled(ctx, 1, 40)
		o := edakit.Int(0, 10)
		bytes := float64(r)*1e6*5*8 + float64(o)*1e6*9*8
		ans := edakit.Num(bytes/1e6, 0)
		q := edakit.Sprintf("A GDSII layout holds %d million rectangles and %d million octagons, each stored as a BOUNDARY whose point list repeats the first vertex to close it. Each point is two 4-byte integers. How many megabytes ($10^6$ bytes) of coordinate data is that, ignoring record headers?", r, o)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("A rectangle stores 5 points (40 bytes) and an octagon 9 points (72 bytes): $%d \\cdot 40 + %d \\cdot 72 = %s$ MB.", r, o, ans), nil)
	default:
		r1, c1 := edakit.Int(2, 64), edakit.Int(2, 64)
		r2, c2 := edakit.Int(2, 32), edakit.Int(2, 32)
		ans := (r1*c1 - 1) + (r2*c2 - 1)
		q := edakit.Sprintf("A layout contains a %d by %d array of one SRAM bit cell and a %d by %d array of one decap cell. Written flat, each placement is its own record. In OASIS, each array becomes a single placement record with a repetition. How many placement records does OASIS save in total?", r1, c1, r2, c2)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Each array saves rows times columns minus one: $(%d - 1) + (%d - 1) = %d$.", r1*c1, r2*c2, ans), nil)
	}
}
