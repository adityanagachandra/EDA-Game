package infra

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerParse) }

func registerParse(reg *generator.Registry) {
	reg.Register("infra.parse.verilog_netlist", edakit.Func(verilogNetlist))
	reg.Register("infra.parse.liberty_parser", &edakit.Bank{Items: libertyParserItems})
	reg.Register("infra.parse.liberty_units", edakit.Func(libertyUnits))
	reg.Register("infra.parse.def_parser", &edakit.Bank{Items: defParserItems})
	reg.Register("infra.parse.def_units", edakit.Func(defUnits))
	reg.Register("infra.parse.sdc_tcl", &edakit.Bank{Items: sdcTclItems})
	reg.Register("infra.parse.spef_parser", edakit.Func(spefParser))
	reg.Register("infra.parse.gds_records", edakit.Func(gdsRecords))
	reg.Register("infra.parse.interchange_formats", &edakit.Bank{Items: interchangeFormatsItems})
}

// parseCell is a library cell as a netlist parser sees it: input pins in declared order, then
// the output pin.
type parseCell struct {
	name string
	ins  []string
	out  string
}

var parseCells = []parseCell{
	{"INV_X1", []string{"A"}, "ZN"},
	{"BUF_X1", []string{"A"}, "Z"},
	{"NAND2_X1", []string{"A1", "A2"}, "ZN"},
	{"NOR2_X1", []string{"A1", "A2"}, "ZN"},
	{"AND2_X1", []string{"A1", "A2"}, "ZN"},
	{"NAND3_X1", []string{"A1", "A2", "A3"}, "ZN"},
	{"AOI21_X1", []string{"A", "B1", "B2"}, "ZN"},
}

type parseInst struct {
	name string
	cell parseCell
	in   []string
	out  string
}

// parseNetlist draws an acyclic structural netlist: primary inputs a, b, ..., instances u1..un,
// wires n1.. and outputs y0.. (instance outputs nobody reads).
func parseNetlist(ctx generator.GeneratorContext) ([]string, []parseInst, []string, []string) {
	nIn := edakit.Scaled(ctx, 2, 4)
	nInst := edakit.Scaled(ctx, 4, 8)
	inputs := []string{"a", "b", "c", "d"}[:nIn]
	signals := append([]string(nil), inputs...)
	var insts []parseInst
	for i := 1; i <= nInst; i++ {
		c := edakit.Pick(parseCells)
		var in []string
		for range c.ins {
			// Prefer recent signals so the netlist is deep rather than flat.
			lo := len(signals) - 4
			if lo < 0 {
				lo = 0
			}
			if edakit.Int(0, 2) == 0 {
				lo = 0
			}
			in = append(in, signals[edakit.Int(lo, len(signals)-1)])
		}
		out := edakit.Sprintf("t%d", i)
		insts = append(insts, parseInst{edakit.Sprintf("u%d", i), c, in, out})
		signals = append(signals, out)
	}
	// Rename: nets with no reader become outputs y0.., the others wires n1..
	read := map[string]bool{}
	for _, x := range insts {
		for _, s := range x.in {
			read[s] = true
		}
	}
	ren := map[string]string{}
	var wires, outputs []string
	for _, x := range insts {
		if read[x.out] {
			ren[x.out] = edakit.Sprintf("n%d", len(wires)+1)
			wires = append(wires, ren[x.out])
		} else {
			ren[x.out] = edakit.Sprintf("y%d", len(outputs))
			outputs = append(outputs, ren[x.out])
		}
	}
	for i := range insts {
		for j, s := range insts[i].in {
			if r, ok := ren[s]; ok {
				insts[i].in[j] = r
			}
		}
		insts[i].out = ren[insts[i].out]
	}
	return inputs, insts, wires, outputs
}

func parseRenderNetlist(inputs []string, insts []parseInst, wires, outputs []string, positional map[int]bool) string {
	var b strings.Builder
	ports := append(append([]string(nil), inputs...), outputs...)
	b.WriteString("`module top (" + strings.Join(ports, ", ") + "); input " + strings.Join(inputs, ", ") + "; output " + strings.Join(outputs, ", ") + ";")
	if len(wires) > 0 {
		b.WriteString(" wire " + strings.Join(wires, ", ") + ";")
	}
	for i, x := range insts {
		var conns []string
		if positional[i] {
			conns = append(append(conns, x.in...), x.out)
		} else {
			for j, p := range x.cell.ins {
				conns = append(conns, "."+p+"("+x.in[j]+")")
			}
			conns = append(conns, "."+x.cell.out+"("+x.out+")")
		}
		b.WriteString(" " + x.cell.name + " " + x.name + " (" + strings.Join(conns, ", ") + ");")
	}
	b.WriteString(" endmodule`")
	return b.String()
}

func verilogNetlist(ctx generator.GeneratorContext) generator.Problem {
	inputs, insts, wires, outputs := parseNetlist(ctx)
	k := variant(ctx)
	pos := map[int]bool{}
	if k == 3 {
		for i := range insts {
			if edakit.Int(0, 1) == 0 {
				pos[i] = true
			}
		}
		pos[edakit.Int(0, len(insts)-1)] = true
	}
	text := parseRenderNetlist(inputs, insts, wires, outputs, pos)
	iterms := 0
	fan := map[string]int{}
	for _, x := range insts {
		iterms += len(x.cell.ins) + 1
		for _, s := range x.in {
			fan[s]++
		}
	}
	nets := len(inputs) + len(wires) + len(outputs)
	facts := map[string]string{"instances": edakit.IntStr(len(insts)), "nets": edakit.IntStr(nets), "iterms": edakit.IntStr(iterms), "bterms": edakit.IntStr(len(inputs) + len(outputs))}
	if k == 1 {
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric("A parser reads this netlist into a design database: "+text+" How many nets does the block contain? Count every distinct signal, ports included.",
				edakit.IntStr(nets),
				edakit.Sprintf("Each port is also a net: %d inputs and %d outputs, plus %d declared wires, gives %d nets. Every instance drives exactly one net here, so the count is inputs plus instances.", len(inputs), len(outputs), len(wires), nets), facts)
		}
		var parts []string
		for _, x := range insts {
			parts = append(parts, edakit.Sprintf("%s has %d", x.name, len(x.cell.ins)+1))
		}
		return edakit.Numeric("A parser reads this netlist into a design database: "+text+" How many instance terminals (ITerms) does it create? Count one per pin of every instance, inputs and outputs.",
			edakit.IntStr(iterms),
			edakit.Sprintf("One ITerm per instance pin: %s, for %d in all. Ports become BTerms, not ITerms.", strings.Join(parts, ", "), iterms), facts)
	}
	// KP 2 and 3: fanout of a net that has readers.
	var cands []string
	for _, s := range append(append([]string(nil), inputs...), wires...) {
		if fan[s] > 0 {
			cands = append(cands, s)
		}
	}
	best := cands[0]
	for _, s := range cands {
		if fan[s] > fan[best] {
			best = s
		}
	}
	net := best
	if edakit.Int(0, 2) == 0 {
		net = edakit.Pick(cands)
	}
	var sinks []string
	for _, x := range insts {
		for j, s := range x.in {
			if s == net {
				sinks = append(sinks, x.name+"/"+x.cell.ins[j])
			}
		}
	}
	facts["fanout"] = edakit.IntStr(fan[net])
	if k == 2 {
		return edakit.Numeric("A parser reads this netlist: "+text+edakit.Sprintf(" How many instance input pins does net %s drive (its fanout)?", net),
			edakit.IntStr(fan[net]),
			edakit.Sprintf("Scan every connection for %s and keep the input pins: %s, which is %d.", net, names(sinks), fan[net]), facts)
	}
	used := map[string]bool{}
	var decl []string
	for _, x := range insts {
		if !used[x.cell.name] {
			used[x.cell.name] = true
			decl = append(decl, x.cell.name+" ("+strings.Join(append(append([]string(nil), x.cell.ins...), x.cell.out), ", ")+")")
		}
	}
	return edakit.Numeric("Some instances below connect by position. The library declares the port orders "+strings.Join(decl, "; ")+". Netlist: "+text+edakit.Sprintf(" How many instance input pins does net %s drive?", net),
		edakit.IntStr(fan[net]),
		edakit.Sprintf("Map each positional connection onto the declared port order, where the last position is the output pin. The input pins on %s are %s, which is %d.", net, names(sinks), fan[net]), facts)
}

func libertyUnits(ctx generator.GeneratorContext) generator.Problem {
	type tu struct {
		text string
		ps   float64
	}
	units := []tu{{"1ns", 1000}, {"1ps", 1}, {"10ps", 10}, {"100ps", 100}}
	switch variant(ctx) {
	case 1:
		u := edakit.Pick(units)
		dps := float64(edakit.Int(150, 4800)) / 10
		v := dps / u.ps
		tbl := []string{"cell_rise", "cell_fall", "rise_transition", "fall_transition"}
		q := edakit.Sprintf("A Liberty library header declares `time_unit : \"%s\" ;`. A %s table entry for one arc reads %s. What is that delay in ps? Round to 1 decimal place.", u.text, edakit.Pick(tbl), edakit.Num(v, 5))
		return edakit.Numeric(q, edakit.Num(dps, 1),
			edakit.Sprintf("Table values are in the library time unit, %s = %s ps, so %s × %s = %s ps.", u.text, edakit.Num(u.ps, 0), edakit.Num(v, 5), edakit.Num(u.ps, 0), edakit.Num(dps, 1)),
			map[string]string{"unit_ps": edakit.Num(u.ps, 0), "delay_ps": edakit.Num(dps, 1)})
	case 2:
		pf := edakit.Int(0, 2) > 0
		unit, scale := "(1,ff)", 1.0
		if pf {
			unit, scale = "(1,pf)", 1000.0
		}
		cff := float64(edakit.Int(4, 90)) / 10
		if edakit.Int(0, 1) == 0 {
			v := cff / scale
			q := edakit.Sprintf("A Liberty library declares `capacitive_load_unit %s ;` and an input pin has `capacitance : %s ;`. What is the pin capacitance in fF? Round to 1 decimal place.", unit, edakit.Num(v, 5))
			return edakit.Numeric(q, edakit.Num(cff, 1),
				edakit.Sprintf("Capacitance values are in the load unit %s, which is %s fF, so the pin is %s fF.", unit, edakit.Num(scale, 0), edakit.Num(cff, 1)),
				map[string]string{"unit_ff": edakit.Num(scale, 0), "cap_ff": edakit.Num(cff, 1)})
		}
		load := float64(edakit.Int(10, 400)) / 10
		v := load / scale
		places := 4
		if !pf {
			places = 1
		}
		q := edakit.Sprintf("A timer must look up an output load of %s fF in a delay table whose index_2 is total_output_net_capacitance. The library declares `capacitive_load_unit %s ;`. What index_2 value, in library units, corresponds to that load? Round to %d decimal places.", edakit.Num(load, 1), unit, places)
		if places == 1 {
			q = strings.Replace(q, "1 decimal places", "1 decimal place", 1)
		}
		return edakit.Numeric(q, edakit.Num(v, places),
			edakit.Sprintf("Index values share the load unit, %s fF per unit, so %s fF is %s library units.", edakit.Num(scale, 0), edakit.Num(load, 1), edakit.Num(v, places)),
			map[string]string{"unit_ff": edakit.Num(scale, 0), "index": edakit.Num(v, places)})
	}
	ua := units[0]
	ub := edakit.Pick(units[1:])
	n := edakit.Scaled(ctx, 3, 5)
	total := 0.0
	raw := 0.0
	var parts []string
	for i := 0; i < n; i++ {
		var u tu
		if i%2 == 0 {
			u = ua
		} else {
			u = ub
		}
		if i == n-1 && edakit.Int(0, 1) == 0 {
			u = ub
		}
		dps := float64(edakit.Int(20, 300))
		v := dps / u.ps
		total += dps
		raw += v
		lib := "A"
		if u == ub {
			lib = "B"
		}
		parts = append(parts, edakit.Sprintf("%s from library %s", edakit.Num(v, 4), lib))
	}
	q := edakit.Sprintf("A path crosses cells from two libraries. Library A declares `time_unit : \"%s\" ;` and library B declares `time_unit : \"%s\" ;`. The arc delays read from the tables are %s. What is the path delay in ps? Round to 1 decimal place.", ua.text, ub.text, strings.Join(parts, ", "))
	return edakit.Numeric(q, edakit.Num(total, 1),
		edakit.Sprintf("Each value is scaled by its own library's unit before summing: A values × %s, B values × %s, for %s ps. Adding the raw numbers (%s) mixes units.", edakit.Num(ua.ps, 0), edakit.Num(ub.ps, 0), edakit.Num(total, 1), edakit.Num(raw, 4)),
		map[string]string{"total_ps": edakit.Num(total, 1)})
}

func defUnits(ctx generator.GeneratorContext) generator.Problem {
	ns := []int{100, 200, 1000, 2000, 4000}
	switch variant(ctx) {
	case 1:
		n := edakit.Pick(ns)
		x := edakit.Int(200, 90000) * n / 100
		y := edakit.Int(200, 90000) * n / 100
		inst := edakit.Sprintf("u%d", edakit.Int(10, 9999))
		cell := edakit.Pick(parseCells).name
		axis, v := "x", x
		if edakit.Int(0, 1) == 0 {
			axis, v = "y", y
		}
		q := edakit.Sprintf("A DEF file declares `UNITS DISTANCE MICRONS %d ;` and its COMPONENTS section has `- %s %s + PLACED ( %d %d ) N ;`. What is the %s coordinate of the instance origin in microns? Round to 2 decimal places.", n, inst, cell, x, y, axis)
		um := float64(v) / float64(n)
		return edakit.Numeric(q, edakit.Num(um, 2),
			edakit.Sprintf("DEF stores integers in database units; one micron is %d units, so %d / %d = %s microns.", n, v, n, edakit.Num(um, 2)),
			map[string]string{"dbu_per_micron": edakit.IntStr(n), "microns": edakit.Num(um, 2)})
	case 2:
		n := edakit.Pick(ns)
		x0, y0 := 0, 0
		if edakit.Int(0, 1) == 0 {
			x0, y0 = edakit.Int(1, 20)*n/10, edakit.Int(1, 20)*n/10
		}
		w := edakit.Int(50, 900) * n / 10
		h := edakit.Int(50, 900) * n / 10
		area := float64(w) / float64(n) * float64(h) / float64(n)
		q := edakit.Sprintf("A DEF file declares `UNITS DISTANCE MICRONS %d ;` and `DIEAREA ( %d %d ) ( %d %d ) ;`. What is the die area in square microns? Round to 2 decimal places.", n, x0, y0, x0+w, y0+h)
		return edakit.Numeric(q, edakit.Num(area, 2),
			edakit.Sprintf("Width is %d dbu = %s um and height %d dbu = %s um, so the area is %s square microns. Dividing the product of dbu by %d once, not %d squared, is the classic error.", w, edakit.Num(float64(w)/float64(n), 2), h, edakit.Num(float64(h)/float64(n), 2), edakit.Num(area, 2), n, n),
			map[string]string{"area_um2": edakit.Num(area, 2)})
	}
	sites := []struct {
		w float64
		s string
	}{{0.19, "0.19"}, {0.46, "0.46"}, {0.2, "0.2"}, {0.054, "0.054"}, {0.28, "0.28"}}
	st := edakit.Pick(sites)
	n := edakit.Pick([]int{1000, 2000, 4000})
	step := int(math.Round(st.w * float64(n)))
	k := edakit.Scaled(ctx, 50, 900)
	x0 := edakit.Int(5, 40) * step
	end := x0 + k*step
	q := edakit.Sprintf("A LEF SITE has width %s um. A DEF with `UNITS DISTANCE MICRONS %d ;` places a row at x = %d with `DO %d BY 1` and a STEP equal to the site width. What is the x coordinate, in DEF database units, of the right edge of the row?", st.s, n, x0, k)
	return edakit.Numeric(q, edakit.IntStr(end),
		edakit.Sprintf("LEF dimensions are microns, DEF integers: the step is %s × %d = %d dbu, and the row ends at %d + %d × %d = %d.", st.s, n, step, x0, k, step, end),
		map[string]string{"step_dbu": edakit.IntStr(step), "end_dbu": edakit.IntStr(end)})
}

func spefParser(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		mods := []string{"u_core", "u_soc", "u_cpu", "u_dsp"}
		subs := []string{"u_alu", "u_lsu", "u_fetch", "u_regfile", "u_mul"}
		name := edakit.Sprintf("%s/%s/n_%d", edakit.Pick(mods), edakit.Pick(subs), edakit.Int(100, 99999))
		L := len(name)
		r := edakit.Scaled(ctx, 4, 30)
		idx := edakit.Int(1, 99999)
		tok := edakit.Sprintf("*%d", idx)
		t := len(tok)
		saved := r*L - r*t - (t + 1 + L + 1)
		q := edakit.Sprintf("A SPEF writer can replace the net name %s (%d characters) by the index %s through *NAME_MAP. The name appears %d times in the file. Each name-map line is the index, a space, the name and a newline. How many bytes does the map save?", name, L, tok, r)
		return edakit.Numeric(q, edakit.IntStr(saved),
			edakit.Sprintf("Without the map the name costs %d × %d = %d bytes. With it, references cost %d × %d = %d and the map line costs %d + 1 + %d + 1 = %d, so the saving is %d bytes.", r, L, r*L, r, t, r*t, t, L, t+L+2, saved),
			map[string]string{"saved": edakit.IntStr(saved)})
	case 2:
		pf := edakit.Int(0, 1) == 0
		unit, scale, places := "1 FF", 1.0, 1
		if pf {
			unit, scale, places = "1 PF", 1000.0, 4
		}
		net := edakit.Sprintf("n%d", edakit.Int(10, 999))
		ng := edakit.Scaled(ctx, 2, 5)
		nc := edakit.Int(1, 3)
		var lines []string
		total, ground := 0.0, 0.0
		i := 1
		for j := 0; j < ng; j++ {
			c := float64(edakit.Int(2, 60)) / 10
			total += c
			ground += c
			lines = append(lines, edakit.Sprintf("%d %s:%d %s", i, net, j+1, edakit.Num(c/scale, places)))
			i++
		}
		for j := 0; j < nc; j++ {
			c := float64(edakit.Int(1, 20)) / 10
			total += c
			lines = append(lines, edakit.Sprintf("%d %s:%d n%d:%d %s", i, net, edakit.Int(1, ng), edakit.Int(1000, 1999), edakit.Int(1, 4), edakit.Num(c/scale, places)))
			i++
		}
		q := edakit.Sprintf("A SPEF file declares `*C_UNIT %s`. The *CAP section of `*D_NET %s` lists: %s. Lines with one node are capacitance to ground and lines with two nodes are coupling capacitance. What total capacitance, in fF, should the *D_NET header carry for this net? Round to 1 decimal place.", unit, net, strings.Join(lines, "; "))
		return edakit.Numeric(q, edakit.Num(total, 1),
			edakit.Sprintf("The total includes ground and coupling entries: %s fF to ground plus %s fF of coupling, scaled by the C unit (%s fF), is %s fF.", edakit.Num(ground, 1), edakit.Num(total-ground, 1), edakit.Num(scale, 0), edakit.Num(total, 1)),
			map[string]string{"total_ff": edakit.Num(total, 1), "ground_ff": edakit.Num(ground, 1)})
	}
	kohm := edakit.Int(0, 1) == 0
	unit, scale, places := "1 OHM", 1.0, 0
	if kohm {
		unit, scale, places = "1 KOHM", 1000.0, 3
	}
	net := edakit.Sprintf("n%d", edakit.Int(10, 999))
	n := edakit.Scaled(ctx, 2, 5)
	var lines []string
	path := 0
	idx := 1
	prev := "u1:ZN"
	for j := 1; j <= n; j++ {
		r := edakit.Int(5, 90)
		path += r
		node := edakit.Sprintf("%s:%d", net, j)
		lines = append(lines, edakit.Sprintf("%d %s %s %s", idx, prev, node, edakit.Num(float64(r)/scale, places)))
		idx++
		prev = node
	}
	r := edakit.Int(5, 90)
	path += r
	lines = append(lines, edakit.Sprintf("%d %s u4:A %s", idx, prev, edakit.Num(float64(r)/scale, places)))
	idx++
	br := edakit.Int(1, n)
	rb := edakit.Int(5, 90)
	lines = append(lines, edakit.Sprintf("%d %s:%d u7:B %s", idx, net, br, edakit.Num(float64(rb)/scale, places)))
	lines = edakit.Shuffle(lines)
	q := edakit.Sprintf("A SPEF file declares `*R_UNIT %s`. The *RES section of `*D_NET %s`, driven by u1:ZN, lists (index, node, node, value): %s. What is the series resistance from the driver u1:ZN to the sink u4:A, in ohms?", unit, net, strings.Join(lines, "; "))
	return edakit.Numeric(q, edakit.IntStr(path),
		edakit.Sprintf("Follow the chain u1:ZN through %s:1 ... %s:%d to u4:A, skipping the branch to u7:B, and scale by the R unit: %d ohms.", net, net, n, path),
		map[string]string{"path_ohm": edakit.IntStr(path)})
}

func gdsRecords(ctx generator.GeneratorContext) generator.Problem {
	shapes := map[int]string{4: "rectangle", 6: "L-shaped polygon", 8: "U-shaped polygon", 10: "polygon", 12: "cross-shaped polygon"}
	switch variant(ctx) {
	case 1:
		n := edakit.Pick([]int{4, 6, 8, 10, 12})
		if edakit.Int(0, 1) == 0 {
			n = 2 * edakit.Scaled(ctx, 2, 20)
		}
		desc := shapes[n]
		if desc == "" {
			desc = "polygon"
		}
		ln := 4 + 8*(n+1)
		q := edakit.Sprintf("A GDSII writer emits a BOUNDARY for a %s with %d distinct vertices. Every record starts with a 4-byte header (2-byte length that includes the header, 1-byte record type, 1-byte data type), each coordinate is a 4-byte signed integer, and a boundary repeats its first point to close. How many bytes long is the XY record?", desc, n)
		return edakit.Numeric(q, edakit.IntStr(ln),
			edakit.Sprintf("The closed boundary stores %d + 1 = %d points, each an (x, y) pair of 8 bytes, plus the 4-byte header: 4 + 8 × %d = %d bytes.", n, n+1, n+1, ln),
			map[string]string{"points": edakit.IntStr(n + 1), "xy_bytes": edakit.IntStr(ln)})
	case 2:
		m := edakit.Scaled(ctx, 1, 4)
		total := 0
		var parts []string
		var work []string
		for i := 0; i < m; i++ {
			n := edakit.Pick([]int{4, 4, 6, 8})
			b := 24 + 8*(n+1)
			total += b
			parts = append(parts, edakit.Sprintf("a %d-vertex boundary", n))
			work = append(work, edakit.Sprintf("%d", b))
		}
		hasPath := edakit.Int(0, 1) == 0 || ctx.Difficulty > 0.7
		pathText := ""
		if hasPath {
			k := edakit.Int(2, 5)
			b := 4 + 6 + 6 + 8 + (4 + 8*k) + 4
			total += b
			parts = append(parts, edakit.Sprintf("a %d-point PATH with a WIDTH record", k))
			work = append(work, edakit.Sprintf("%d", b))
			pathText = " A PATH is PATH (4 bytes), LAYER, DATATYPE, WIDTH (header plus a 4-byte integer), XY with its points (an open path does not repeat a point) and ENDEL."
		}
		q := edakit.Sprintf("Each GDSII record has a 4-byte header. A BOUNDARY element is the records BOUNDARY (no data), LAYER and DATATYPE (one 2-byte integer each), XY (closed, so the first point repeats; 8 bytes per point) and ENDEL (no data).%s A structure holds %s. How many bytes do these elements occupy?", pathText, names(parts))
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("A boundary with n vertices costs 4 + 6 + 6 + (4 + 8(n + 1)) + 4 = 24 + 8(n + 1) bytes. The element sizes are %s, totalling %d bytes.", strings.Join(work, " + "), total),
			map[string]string{"bytes": edakit.IntStr(total)})
	}
	type du struct {
		m    float64
		text string
	}
	dus := []du{{1e-9, "1e-9"}, {1e-8, "1e-8"}, {5e-10, "5e-10"}, {2.5e-10, "2.5e-10"}}
	d := edakit.Pick(dus)
	user := d.m / 1e-6
	v := edakit.Int(1000, 2000000)
	um := float64(v) * d.m / 1e-6
	q := edakit.Sprintf("A GDSII library's UNITS record holds the size of a database unit in user units, %s, and in meters, %s. An XY coordinate stores the integer %d. What is that coordinate in microns? Round to 4 decimal places.", edakit.Num(user, 5), d.text, v)
	return edakit.Numeric(q, edakit.Num(um, 4),
		edakit.Sprintf("One database unit is %s m = %s um, so %d units are %s um.", d.text, edakit.Num(user, 5), v, edakit.Num(um, 4)),
		map[string]string{"microns": edakit.Num(um, 4)})
}

var libertyParserItems = []edakit.Item{
	{KP: 1, Q: "A Liberty parser meets `capacitive_load_unit (1,pf) ;` in the library header. Which construct is it?",
		A: "A complex attribute: a name with a parenthesized value list", D: []string{"A simple attribute: a name, a colon and one value", "A group: a name, arguments and a braced body", "A define statement that declares a new attribute"},
		Why: "Liberty has three statement shapes. Simple attributes use `name : value ;`, complex attributes use `name (v1, v2, ...) ;`, and groups add a braced body."},
	{KP: 1, Q: "Which statement is a Liberty group rather than an attribute?",
		A: "`pin (A1) { direction : input ; capacitance : 0.0016 ; }`", D: []string{"`direction : input ;`", "`index_1 (\"0.01, 0.04, 0.16\") ;`", "`time_unit : \"1ns\" ;`"},
		Why: "A group has a name, an argument list and a braced body holding further statements; the others are simple or complex attributes."},
	{KP: 1, Q: "A generic Liberty parser builds a tree from the grammar alone, before knowing what `cell_rise` means. What does the tree contain?",
		A: "Nested groups, each with its arguments, attributes and child groups", D: []string{"A flat table of every timing arc in the library", "A delay-calculation graph with arcs between pins", "Only the cells, with the header attributes discarded"},
		Why: "Liberty's syntax is uniform, so a first pass can build a generic tree of groups and attributes; a second pass interprets it into cells, pins and arcs."},
	{KP: 1, Q: "Why can one grammar parse every Liberty file even though vendors add their own attributes?",
		A: "Every statement is a group, a simple attribute or a complex attribute, whatever its name", D: []string{"Vendors must register each attribute with the parser at build time", "Unknown attributes are rejected, so files only use standard names", "The grammar lists every legal attribute name explicitly"},
		Why: "The syntax does not depend on the attribute vocabulary. A parser can keep or skip unknown names and still recover the structure."},
	{KP: 2, Q: "In a Liberty cell, where does the timing group for the arc from input A1 to output ZN live?",
		A: "Inside pin (ZN), with `related_pin : \"A1\"`", D: []string{"Inside pin (A1), with `related_pin : \"ZN\"`", "Directly under the cell, outside any pin", "In the library header, keyed by cell name"},
		Why: "Timing groups sit in the pin they arrive at, and related_pin names the pin the arc starts from."},
	{KP: 2, Q: "A `cell_rise (delay_template_7x7)` group has index_1, index_2 and values. What must the parser resolve to interpret the indices?",
		A: "The lu_table_template named delay_template_7x7, which says which variable each index is", D: []string{"The cell's area attribute", "The related_pin of the next timing group", "The default_max_transition of the library"},
		Why: "The template declares variable_1 and variable_2, for example input transition and output load, so the parser knows which axis is which and can fall back to its default indices."},
	{KP: 2, Q: "A timing group carries `timing_sense : negative_unate`. What does the timer use it for?",
		A: "To pair a rising input with a falling output when propagating the arc", D: []string{"To choose the setup check instead of the hold check", "To negate the table values before interpolation", "To skip the arc during early analysis"},
		Why: "Negative unate means an input rise causes an output fall, so the timer combines input rise transitions with cell_fall tables, and vice versa."},
	{KP: 2, Q: "The values attribute of a 7 × 7 table arrives as seven quoted strings of comma-separated numbers. What must the parser check before storing it?",
		A: "That it has as many rows as index_1 entries and as many columns as index_2 entries", D: []string{"That the values increase along both axes", "That every value is positive", "That index_1 and index_2 have equal values"},
		Why: "The shape is fixed by the indices; a mismatch is a malformed table. Monotonicity and sign are not required, since negative delays can appear in some tables."},
	{KP: 2, Q: "Which pin attribute lets a logic tool recover what a cell computes?",
		A: "`function : \"!(A1 & A2)\"` on the output pin", D: []string{"`related_pin : \"A1\"` in the timing group", "`direction : output` on the output pin", "`area : 0.798` on the cell"},
		Why: "The function string is a Boolean expression of the input pins; synthesis and equivalence checkers parse it into logic."},
	{KP: 3, Q: "A library for one corner is several hundred megabytes, but a placer needs only pin capacitances and cell areas. What parser design fits best?",
		A: "A streaming parser that interprets the groups it needs and skips the bodies it does not", D: []string{"A full in-memory tree of every group, then a query pass", "Converting the file to XML first, then parsing the XML", "Reading only the first megabyte, where the cells are"},
		Why: "Most of the bytes are timing and power tables. Skipping unneeded groups while still tracking brace nesting saves time and memory."},
	{KP: 3, Q: "A parser reads `cell_rise (tmpl_A)` before it has seen any lu_table_template named tmpl_A. What should it do?",
		A: "Report an error, because templates must be declared before tables refer to them", D: []string{"Assume a 7 × 7 template with default indices", "Silently treat the table as scalar", "Search the next library file for the template"},
		Why: "Templates are declared at library level before use. Guessing a template would misread the axes and corrupt every delay looked up from that table."},
	{KP: 3, Q: "Two libraries are read into one session, one with time_unit \"1ns\" and one with time_unit \"1ps\". Where must the scaling happen?",
		A: "When each library is read, converting its values into the tool's internal units", D: []string{"Only when a report is printed", "Never; the timer uses the first library's unit for both", "When the netlist is linked, from the cell names"},
		Why: "Units are per library. Converting at read time means every engine downstream sees one consistent unit system."},
	{KP: 3, Q: "A user complains that a cell's timing in the tool differs from the vendor's datasheet. The parser kept only the first timing group per pin. What was lost?",
		A: "The arcs from other related pins, and conditional (when) arcs for the same pin", D: []string{"The cell area and leakage power", "The library's time unit", "The pin's direction attribute"},
		Why: "An output pin carries one timing group per related pin and often several state-dependent groups. Keeping one drops real arcs."},
}

var defParserItems = []edakit.Item{
	{KP: 1, Q: "A reader finds LAYER metal2 with DIRECTION, PITCH, WIDTH and SPACING statements. Which file is it reading?",
		A: "A technology LEF", D: []string{"A cell LEF", "The DEF NETS section", "A Liberty library"},
		Why: "Technology LEF defines the routing layers, vias, sites and units; cell LEF defines MACROs."},
	{KP: 1, Q: "Which LEF content does a placer need to put a standard cell into a row?",
		A: "The MACRO's SIZE, its SITE and its CLASS", D: []string{"The MACRO's PIN shapes on metal1 only", "The routing layer PITCH values", "The VIA definitions"},
		Why: "SIZE gives the footprint, SITE says which rows it fits, and CLASS separates core cells from pads and blocks."},
	{KP: 1, Q: "A cell LEF MACRO lists OBS with rectangles on metal1 and metal2. What must the router do with them?",
		A: "Treat them as blockages that no other net may overlap", D: []string{"Connect them to the nearest power net", "Ignore them, since only PIN shapes matter", "Copy them into the DEF SPECIALNETS"},
		Why: "OBS (obstructions) mark internal metal of the cell that routes must avoid."},
	{KP: 1, Q: "A PIN in a cell LEF has a PORT with RECT shapes on metal1. What does a parser store for it?",
		A: "The pin's access geometry, which the router uses to connect wires to the pin", D: []string{"The pin's capacitance for timing", "The pin's logic function", "The pin's placement in the die"},
		Why: "LEF describes physical abstraction: shapes and layers. Electrical data such as capacitance comes from Liberty."},
	{KP: 2, Q: "Which DEF section creates the instances of a design?",
		A: "COMPONENTS", D: []string{"PINS", "NETS", "ROWS"},
		Why: "COMPONENTS lists each instance with its master and, when placed, its location and orientation. PINS creates the block's ports."},
	{KP: 2, Q: "A DEF line reads `- u7 NAND2_X1 + PLACED ( 3800 5600 ) FS ;`. What does FS give the database?",
		A: "The orientation: mirrored about the x-axis, as used for flipped rows", D: []string{"The fixed status, meaning the placer must not move it", "The fill-cell status of the instance", "The first site of the row"},
		Why: "N, S, E, W, FN, FS, FE and FW are orientations. FS is the mirror used for flipped rows; FIXED and PLACED are placement statuses."},
	{KP: 2, Q: "Where does a DEF describe the power and ground straps of the PDN?",
		A: "SPECIALNETS", D: []string{"NETS", "BLOCKAGES", "TRACKS"},
		Why: "SPECIALNETS carries prerouted special wiring such as power straps and rings; NETS carries regular signal connectivity and routing."},
	{KP: 2, Q: "A DEF has TRACKS and GCELLGRID statements. Which engine relies on them most?",
		A: "The router, for its routing grid and its global-routing tiles", D: []string{"The placer, for its site legality", "The timer, for its clock definitions", "The parser, for its unit scaling"},
		Why: "TRACKS define legal wire positions per layer and GCELLGRID the coarse tiles of global routing."},
	{KP: 3, Q: "Why must a tool read the LEF before the DEF?",
		A: "DEF components and vias refer to masters, sites and layers that only the LEF defines", D: []string{"LEF sets the DEF database units", "DEF is compressed with a key stored in the LEF", "The DEF NETS section is written inside the LEF"},
		Why: "A DEF component names a master such as NAND2_X1. Without the LEF the reader cannot resolve it, its size or its pins."},
	{KP: 3, Q: "A DEF NETS entry connects `( u7 A3 )`, but the LEF MACRO for u7's master has only pins A1, A2 and ZN. What should a robust reader do?",
		A: "Report the unknown pin and keep reading, so the user sees every such error", D: []string{"Create a new pin A3 on the master", "Connect the net to A1 instead", "Drop the whole NETS section"},
		Why: "Cross-reference errors are design bugs. Reporting them with line numbers, without inventing pins, is what lets a user fix the input."},
	{KP: 3, Q: "The LEF declares `DATABASE MICRONS 2000` and a DEF says `UNITS DISTANCE MICRONS 1000`. Is that consistent?",
		A: "Yes: the DEF resolution may not exceed the LEF's, and 1000 is below 2000", D: []string{"No: the two values must be equal", "No: the DEF value must be larger", "Yes, but only if both are powers of two"},
		Why: "LEF/DEF require the DEF units to be no finer than the LEF database units, so every DEF integer can be represented exactly."},
	{KP: 3, Q: "A reader streams a 20 GB DEF. Which section benefits most from a callback-style parser that hands each entry to the database as it is read?",
		A: "NETS, which carries the connectivity and routing of millions of nets", D: []string{"UNITS, which is a single line", "DIEAREA, which is one rectangle", "VERSION, which is one token"},
		Why: "Callback parsers (the Si2 LEF/DEF reference parsers work this way) avoid holding an intermediate copy of the biggest sections."},
}

var sdcTclItems = []edakit.Item{
	{KP: 1, Q: "An SDC file sets a Tcl variable with `set clk_period 2.5` and then passes that variable, through Tcl variable substitution, as the -period of create_clock. Is that legal SDC?",
		A: "Yes: SDC is Tcl, so variables and substitution work", D: []string{"No: SDC accepts only numeric literals", "Only if the variable is declared with set_units", "No: square brackets are comments in SDC"},
		Why: "SDC commands are Tcl commands. Variables, expressions, loops and procedures are all available to the file."},
	{KP: 1, Q: "Which value does `[get_ports clk]` do in `create_clock -period 2 [get_ports clk]`?",
		A: "Runs get_ports first and passes its result as an argument", D: []string{"Declares a new port named clk", "Comments out the port name", "Groups the arguments into a list without evaluating"},
		Why: "Square brackets are Tcl command substitution: the inner command runs and its result replaces the bracketed text."},
	{KP: 1, Q: "An SDC line reads `set_input_delay 0.4 -clock clk {in_a in_b}`. What do the braces do?",
		A: "Pass the two names as one list word, without substitution", D: []string{"Run a command named in_a", "Mark the names as optional", "Evaluate in_a and in_b as variables"},
		Why: "Braces quote a word literally; the result is a Tcl list, which the command reads as two port names."},
	{KP: 1, Q: "An SDC file contains a Tcl foreach loop over the list {a b c} whose body calls set_load 0.01 on get_ports of the loop variable. What does the tool see?",
		A: "Three set_load commands, one per port", D: []string{"One set_load on a port named after the variable", "A syntax error, since loops are not SDC", "One set_load applied to a port named {a b c}"},
		Why: "Because SDC is evaluated by a Tcl interpreter, the loop runs and the loop variable is substituted on each iteration."},
	{KP: 2, Q: "Which value does `get_pins u1/*` return inside a timing tool?",
		A: "A collection of pin objects whose names match the pattern", D: []string{"The text string u1/*", "The number of pins of u1", "A list of the nets attached to u1"},
		Why: "Object queries return handles to database objects; constraints then attach to those objects, not to names."},
	{KP: 2, Q: "`get_ports data*` matches nothing because the design's ports are named din[0] to din[7]. What does a careful tool do?",
		A: "Warn that the pattern matched no objects, so the constraint applies to nothing", D: []string{"Create ports named data*", "Apply the constraint to all ports", "Silently match the closest names"},
		Why: "A pattern that matches nothing silently drops a constraint. Tools warn because unconstrained paths are a classic source of escapes."},
	{KP: 2, Q: "A constraint names `[get_pins u_core/u_alu/add_0/A]`, but synthesis renamed the instance. What is the likely outcome?",
		A: "The query returns an empty collection and the constraint is lost", D: []string{"The tool finds the pin by its function", "The constraint moves to the nearest pin", "The design fails to link"},
		Why: "Constraints bind to objects through names. When netlists are renamed, constraints written against hierarchical names silently stop matching."},
	{KP: 2, Q: "Why do `get_cells`, `get_nets` and `get_clocks` return collections rather than plain Tcl lists of names?",
		A: "So commands can work on database objects directly, without re-looking up each name", D: []string{"Because Tcl lists cannot hold more than 1000 elements", "So the results print in sorted order", "Because names are not unique in a netlist"},
		Why: "A collection holds object handles. Re-resolving names for every command would be slow and ambiguous for objects that share names across types."},
	{KP: 3, Q: "An SDC file says `set_input_delay 300 -clock clk [get_ports a]` with no units. How does a tool interpret 300?",
		A: "In the time unit in effect, from set_units or the library, so 300 is wrong if the unit is ns", D: []string{"Always in picoseconds", "Always in nanoseconds", "In clock periods of clk"},
		Why: "SDC numbers are unitless; their meaning comes from set_units or the library's time unit. 300 ns where 300 ps was meant is a classic constraint bug."},
	{KP: 3, Q: "Why do some sign-off flows read SDC with a restricted interpreter instead of a full Tcl shell?",
		A: "To reject side effects and tool-specific commands, so the constraints mean the same everywhere", D: []string{"Because full Tcl cannot parse brackets", "Because SDC files are binary", "Because restricted interpreters run SDC faster on GPUs"},
		Why: "A full shell lets a constraint file run arbitrary code. A restricted reader keeps constraints portable and reproducible across tools."},
	{KP: 3, Q: "In OpenSTA, how is an SDC file read?",
		A: "It is sourced by the tool's Tcl interpreter, with SDC commands defined as Tcl commands", D: []string{"A separate compiler turns it into Liberty", "It is converted to Verilog attributes", "It is parsed by the DEF reader"},
		Why: "OpenSTA defines the SDC commands in its Tcl environment, so read_sdc evaluates the file as Tcl."},
	{KP: 3, Q: "A tool writes constraints back out after optimization. What must write_sdc do with collections such as `[get_pins u1/*]`?",
		A: "Write the matching objects by name in a form that resolves in the new netlist", D: []string{"Write the pattern unchanged, whatever the netlist became", "Drop every constraint that used a wildcard", "Write object handles as memory addresses"},
		Why: "Handles do not survive a session. Written constraints must name objects so that the next tool's queries find them."},
}

var interchangeFormatsItems = []edakit.Item{
	{KP: 1, Q: "An FPGA vendor tool exports a netlist with nested parenthesized forms such as `(cell ... (view ... (interface ...) (contents ...)))`. Which format is it?",
		A: "EDIF", D: []string{"IP-XACT", "UPF", "SPEF"},
		Why: "EDIF uses LISP-like S-expressions to describe libraries, cells, views, interfaces and contents."},
	{KP: 1, Q: "What does an EDIF netlist describe?",
		A: "Cells with their ports and the instances and nets inside them", D: []string{"Register maps and bus interfaces of an IP block", "Power domains and isolation cells", "Wire parasitics per net"},
		Why: "EDIF is a netlist interchange format: hierarchy, interfaces and connectivity."},
	{KP: 1, Q: "Why is an S-expression format like EDIF easy to parse but bulky?",
		A: "Every construct is a parenthesized list, so one recursive reader handles it, but keywords repeat for every object", D: []string{"It stores coordinates as binary floats", "It requires an XML schema validator", "It compresses names with a name map"},
		Why: "A uniform list syntax needs one small reader; the price is verbose, keyword-heavy text."},
	{KP: 1, Q: "A tool imports an EDIF netlist and a Verilog netlist of the same design. What should be identical after import?",
		A: "The instances and their connectivity in the database", D: []string{"The file sizes", "The order of statements", "The comments"},
		Why: "Both are netlist representations; the parsers must build the same object model even though the syntax differs."},
	{KP: 2, Q: "An IP vendor ships an XML file describing a block's bus interfaces, memory map and registers. Which standard is it?",
		A: "IP-XACT (IEEE 1685)", D: []string{"UPF (IEEE 1801)", "EDIF", "Liberty"},
		Why: "IP-XACT is the XML schema for IP metadata: components, bus interfaces, memory maps and registers."},
	{KP: 2, Q: "What can a tool generate from an IP-XACT register description?",
		A: "Register RTL, software header files and documentation from one source", D: []string{"The block's placed layout", "The block's timing arcs", "The block's power domains"},
		Why: "Register descriptions are the single source from which hardware, software and documentation views are generated, so they stay consistent."},
	{KP: 2, Q: "An SoC integrator wires 40 IP blocks by matching their bus interfaces. Which representation lets a tool check that two interfaces are compatible?",
		A: "IP-XACT bus interface definitions", D: []string{"The blocks' GDS", "The blocks' SPEF", "The blocks' SDC"},
		Why: "IP-XACT names bus types and their port mappings, so an integration tool can check and connect interfaces automatically."},
	{KP: 2, Q: "Why does an EDA tool parse IP-XACT with an XML parser plus a schema rather than a hand-written grammar?",
		A: "IP-XACT is defined as an XML schema, so validation comes from the schema", D: []string{"IP-XACT files are binary", "XML parsers are the only way to read S-expressions", "IP-XACT has no fixed structure"},
		Why: "The standard is a schema; generic XML tooling handles syntax and the schema checks structure."},
	{KP: 3, Q: "A UPF file contains `create_power_domain PD_CPU -elements {u_cpu}`. How does a tool read it?",
		A: "With a Tcl interpreter, since UPF commands are Tcl commands", D: []string{"With the Liberty parser", "With an XML parser", "With the DEF reader"},
		Why: "UPF (IEEE 1801) is Tcl based, like SDC. Its commands build power domains, supplies, isolation, level shifting and retention."},
	{KP: 3, Q: "Why is power intent written in UPF rather than in the RTL itself?",
		A: "So the same RTL can be implemented under different power strategies, kept in a separate file", D: []string{"Because RTL cannot describe any logic", "Because UPF replaces the netlist", "Because power domains are only physical"},
		Why: "Separating power intent from function lets verification and implementation share one description without editing the RTL."},
	{KP: 3, Q: "What role does Si2 play in EDA standards?",
		A: "An industry consortium that hosts coalitions maintaining standards such as OpenAccess", D: []string{"A government agency that certifies EDA tools", "A foundry that publishes PDKs", "A company that sells a place-and-route tool"},
		Why: "Si2 (Silicon Integration Initiative) runs member coalitions; OpenAccess and the LEF/DEF reference parsers are among the standards it stewards."},
	{KP: 3, Q: "OpenDB's LEF and DEF readers trace back to reference parsers distributed through Si2. What does reusing a reference parser buy a tool?",
		A: "Consistent interpretation of the format's corner cases across tools", D: []string{"Faster placement", "Automatic timing closure", "A license to the foundry's PDK"},
		Why: "When several tools use the same reference parser, the files mean the same thing to each of them."},
}
