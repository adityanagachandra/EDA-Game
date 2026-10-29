package infra

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerDB) }

func registerDB(reg *generator.Registry) {
	reg.Register("infra.db.design_database", &edakit.Bank{Items: designDatabaseItems})
	reg.Register("infra.db.object_model", edakit.Func(objectModel))
	reg.Register("infra.db.hierarchy_instancing", edakit.Func(hierarchyInstancing))
	reg.Register("infra.db.opendb", &edakit.Bank{Items: opendbItems})
	reg.Register("infra.db.memory_footprint", edakit.Func(memoryFootprint))
	reg.Register("infra.db.data_oriented_layout", edakit.Func(dataOrientedLayout))
	reg.Register("infra.db.memory_management", &edakit.Bank{Items: memoryManagementItems})
	reg.Register("infra.db.serialization", edakit.Func(serialization))
}

// dbMaster is a library cell with its signal pin count (power pins are not ITerms here).
type dbMaster struct {
	name string
	pins int
}

var dbMasters = []dbMaster{
	{"INV_X1", 2}, {"BUF_X2", 2}, {"NAND2_X1", 3}, {"NOR2_X1", 3}, {"AOI21_X1", 4},
	{"OAI22_X1", 5}, {"XOR2_X1", 3}, {"MUX2_X1", 4}, {"DFF_X1", 3}, {"DFFR_X1", 4}, {"FA_X1", 5},
}

func objectModel(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		k := edakit.Scaled(ctx, 3, 5)
		ms := edakit.Shuffle(dbMasters)[:k]
		total := 0
		var parts, work []string
		for _, m := range ms {
			n := edakit.Int(2, 40)
			total += n * m.pins
			parts = append(parts, edakit.Sprintf("%d instances of %s (%d signal pins)", n, m.name, m.pins))
			work = append(work, edakit.Sprintf("%d x %d", n, m.pins))
		}
		bterms := edakit.Int(4, 30)
		q := edakit.Sprintf("A block holds %s, and the block itself has %d ports. In an OpenDB-style object model, every signal pin of every instance is one instance terminal (ITerm) and every block port is one block terminal (BTerm). How many ITerms does the block hold?",
			names(parts), bterms)
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("An ITerm is the pairing of one instance with one master pin, so the count is the sum of instances times pins: %s = %d. The %d ports are BTerms, a different object kind.", strings.Join(work, " + "), total, bterms),
			map[string]string{"iterms": edakit.IntStr(total), "bterms": edakit.IntStr(bterms)})
	case 2:
		iterms := 10 * edakit.Int(20, 90)
		floating := edakit.Int(0, 12)
		bterms := edakit.Int(6, 40)
		nets := edakit.Int(iterms/4, iterms/3)
		conn := iterms - floating + bterms
		avg := float64(conn) / float64(nets)
		q := edakit.Sprintf("A block has %d ITerms, of which %d are left unconnected (spare outputs), %d BTerms, each on a net, and %d nets. Every connected terminal sits on exactly one net. What is the average number of terminals per net? Round to 2 decimal places.",
			iterms, floating, bterms, nets)
		return edakit.Numeric(q, edakit.Num(avg, 2),
			edakit.Sprintf("Each connected ITerm and each BTerm contributes one terminal to one net: %d - %d + %d = %d terminals. Divided over %d nets that is %s terminals per net, the mean hyperedge size.", iterms, floating, bterms, conn, nets, edakit.Num(avg, 2)),
			map[string]string{"terminals": edakit.IntStr(conn), "avg_degree": edakit.Num(avg, 2)})
	}
	k := edakit.Scaled(ctx, 3, 5)
	var parts, work []string
	clique, pins := 0, 0
	for i := 0; i < k; i++ {
		p := edakit.Int(2, 9)
		clique += p * (p - 1) / 2
		pins += p
		parts = append(parts, edakit.Sprintf("n%d with %d terminals", i+1, p))
		work = append(work, edakit.Sprintf("%d", p*(p-1)/2))
	}
	q := edakit.Sprintf("The database stores each net as a hyperedge: one net object with a list of its terminals. An algorithm that needs an ordinary graph expands every net into a clique, one edge per pair of terminals. The nets are %s. How many graph edges does the clique expansion create?", names(parts))
	return edakit.Numeric(q, edakit.IntStr(clique),
		edakit.Sprintf("A net with p terminals becomes p(p-1)/2 edges: %s = %d edges, against %d terminal entries in the hyperedge form. The quadratic growth is why the database keeps the hyperedge and lets each algorithm choose its own model.", strings.Join(work, " + "), clique, pins),
		map[string]string{"clique_edges": edakit.IntStr(clique), "pins": edakit.IntStr(pins)})
}

func hierarchyInstancing(ctx generator.GeneratorContext) generator.Problem {
	// top instantiates a x core and b x periph; core instantiates c x alu and holds lc leaf
	// cells; alu holds la leaf cells; periph holds lp leaf cells.
	a, b, c := edakit.Int(2, 8), edakit.Int(1, 4), edakit.Int(2, 6)
	lc, la, lp := 10*edakit.Int(5, 40), 10*edakit.Int(20, 90), 10*edakit.Int(5, 50)
	lt := edakit.Int(5, 60)
	desc := edakit.Sprintf("Module top holds %d leaf cells, %d instances of module core and %d instances of module periph. Module core holds %d leaf cells and %d instances of module alu. Module alu holds %d leaf cells. Module periph holds %d leaf cells.",
		lt, a, b, lc, c, la, lp)
	coreLeaves := lc + c*la
	flat := lt + a*coreLeaves + b*lp
	facts := map[string]string{"flat_leaves": edakit.IntStr(flat), "core_leaves": edakit.IntStr(coreLeaves)}
	switch variant(ctx) {
	case 1:
		return edakit.Numeric(desc+" How many leaf-cell instances does the flattened design contain?", edakit.IntStr(flat),
			edakit.Sprintf("Flatten bottom up: one core contains %d + %d x %d = %d leaves. Top contains %d + %d x %d + %d x %d = %d leaves.", lc, c, la, coreLeaves, lt, a, coreLeaves, b, lp, flat), facts)
	case 2:
		hier := (lt + a + b) + (lc + c) + la + lp
		saved := flat - hier
		facts["hier_objects"] = edakit.IntStr(hier)
		return edakit.Numeric(desc+" A hierarchical database stores each module definition once: one object per leaf cell or child-module instance inside each definition. A flat database stores one object per flattened leaf cell. How many fewer instance objects does the hierarchical database store?", edakit.IntStr(saved),
			edakit.Sprintf("Hierarchical: top %d + core %d + alu %d + periph %d = %d objects. Flat: %d leaves. Saving: %d - %d = %d.", lt+a+b, lc+c, la, lp, hier, flat, flat, hier, saved), facts)
	}
	occ := a * c
	facts["alu_occurrences"] = edakit.IntStr(occ)
	return edakit.Numeric(desc+" An ECO must change one particular alu, reached through one particular core. How many distinct occurrences (hierarchical paths from top) of module alu exist, all of which share the one alu definition until the design is uniquified?", edakit.IntStr(occ),
		edakit.Sprintf("Each of the %d core instances holds %d alu instances, so there are %d x %d = %d occurrences. Editing the shared definition would change all of them, so the tool first uniquifies: it copies the definition for the occurrence being edited.", a, c, a, c, occ), facts)
}

func memoryFootprint(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		ni := edakit.Scaled(ctx, 1, 20) * 100000
		pinsPer := float64(edakit.Int(25, 40)) / 10
		nn := ni + ni/10*edakit.Int(0, 2)
		si, sn, st := 8*edakit.Int(6, 12), 8*edakit.Int(4, 10), 4*edakit.Int(4, 8)
		terms := float64(ni) * pinsPer
		bytes := float64(ni*si+nn*sn) + terms*float64(st)
		mb := bytes / 1e6
		q := edakit.Sprintf("A netlist database stores each instance in %d bytes, each net in %d bytes and each ITerm in %d bytes. The design has %d instances with an average of %s ITerms each, and %d nets. How much memory do these objects take, in megabytes (1 MB = 10^6 bytes)? Round to 1 decimal place.",
			si, sn, st, ni, edakit.Num(pinsPer, 1), nn)
		return edakit.Numeric(q, edakit.Num(mb, 1),
			edakit.Sprintf("Instances: %d x %d = %d bytes. Nets: %d x %d = %d bytes. ITerms: %s x %d = %s bytes. Total %s bytes, which is %s MB.", ni, si, ni*si, nn, sn, nn*sn, edakit.Num(terms, 0), st, edakit.Num(terms*float64(st), 0), edakit.Num(bytes, 0), edakit.Num(mb, 1)),
			map[string]string{"bytes": edakit.Num(bytes, 0), "mb": edakit.Num(mb, 1)})
	case 2:
		n := edakit.Scaled(ctx, 2, 40) * 1000000
		refs := edakit.Int(3, 6)
		saved := float64(n*refs*4) / 1e6
		q := edakit.Sprintf("Each ITerm record holds %d references to other objects (its instance, its master pin, its net and neighbours in intrusive lists, for example). The design has %d ITerms. If every reference is stored as a 32-bit index into an object table instead of a 64-bit pointer, how much memory is saved, in megabytes (1 MB = 10^6 bytes)? Round to 1 decimal place.", refs, n)
		return edakit.Numeric(q, edakit.Num(saved, 1),
			edakit.Sprintf("Each reference shrinks from 8 to 4 bytes, saving 4 bytes: %d ITerms x %d references x 4 bytes = %s MB. Indices also survive reallocation of the table and serialize without pointer fix-up.", n, refs, edakit.Num(saved, 1)),
			map[string]string{"saved_mb": edakit.Num(saved, 1)})
	}
	type field struct {
		name string
		size int
	}
	pool := []field{{"flags", 1}, {"orient", 1}, {"master", 8}, {"x", 4}, {"y", 4}, {"net", 8}, {"weight", 8}, {"status", 2}, {"id", 4}, {"level", 2}}
	var fs []field
	for tries := 0; tries < generator.MaxDrawAttempts; tries++ {
		fs = edakit.Shuffle(pool)[:edakit.Scaled(ctx, 4, 6)]
		has8, has1 := false, false
		for _, f := range fs {
			if f.size == 8 {
				has8 = true
			}
			if f.size <= 2 {
				has1 = true
			}
		}
		if has8 && has1 {
			break
		}
		fs = nil
	}
	if fs == nil {
		fs = []field{pool[0], pool[2], pool[3], pool[1]}
	}
	off, maxA := 0, 1
	var lay []string
	for _, f := range fs {
		if off%f.size != 0 {
			off += f.size - off%f.size
		}
		lay = append(lay, edakit.Sprintf("%s at %d", f.name, off))
		off += f.size
		if f.size > maxA {
			maxA = f.size
		}
	}
	size := off
	if size%maxA != 0 {
		size += maxA - size%maxA
	}
	var decl []string
	for _, f := range fs {
		decl = append(decl, edakit.Sprintf("%s (%d bytes)", f.name, f.size))
	}
	q := edakit.Sprintf("A C++ record declares its fields in this order: %s. Each field is aligned to a multiple of its own size, and the record size is rounded up to a multiple of its largest field. What is the size of the record in bytes?", strings.Join(decl, ", "))
	return edakit.Numeric(q, edakit.IntStr(size),
		edakit.Sprintf("Lay the fields out in order, padding each to its alignment: %s, ending at byte %d; rounding up to a multiple of %d gives %d bytes. Sorting the fields from largest to smallest would remove most of the padding, which matters when the record is repeated millions of times.", strings.Join(lay, ", "), off, maxA, size),
		map[string]string{"size": edakit.IntStr(size)})
}

func dataOrientedLayout(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 1, 16) * 4096
	s := 8 * edakit.Int(5, 16)
	f := edakit.Pick([]int{4, 8})
	const line = 64
	aos := (n*s + line - 1) / line
	soa := (n*f + line - 1) / line
	stem := edakit.Sprintf("A placer stores %d cells as an array of %d-byte structs (array of structures). Its inner loop reads only the %d-byte x coordinate of every cell, in order. Cache lines are 64 bytes and the array starts on a line boundary. ", n, s, f)
	facts := map[string]string{"aos_lines": edakit.IntStr(aos), "soa_lines": edakit.IntStr(soa)}
	switch variant(ctx) {
	case 1:
		return edakit.Numeric(stem+"How many distinct cache lines does one pass of the loop touch?", edakit.IntStr(aos),
			edakit.Sprintf("The structs are contiguous, so the loop sweeps the whole array: %d x %d = %d bytes, which is %d lines of 64 bytes. Only %d of every %d bytes loaded are used.", n, s, n*s, aos, f, s), facts)
	case 2:
		ratio := float64(aos) / float64(soa)
		facts["ratio"] = edakit.Num(ratio, 2)
		return edakit.Numeric(stem+"The struct is split into a structure of arrays, so all x coordinates sit in one contiguous array. By what factor does the number of cache lines the loop touches fall? Round to 2 decimal places.", edakit.Num(ratio, 2),
			edakit.Sprintf("Array of structures: %d x %d / 64 = %d lines. Structure of arrays: %d x %d / 64 = %d lines. The factor is %d / %d = %s, the struct size over the field size.", n, s, aos, n, f, soa, aos, soa, edakit.Num(ratio, 2)), facts)
	}
	hot := f + 4*edakit.Int(1, 3)
	if hot >= s {
		hot = s - 8
	}
	hl := (n*hot + line - 1) / line
	facts["hot_lines"] = edakit.IntStr(hl)
	return edakit.Numeric(stem+edakit.Sprintf("Instead of a full split, the struct is divided into a hot part of %d bytes (x, y and the fields the inner loops read) stored contiguously, and a cold part (names, properties) stored elsewhere. How many cache lines does one pass over the hot array touch?", hot), edakit.IntStr(hl),
		edakit.Sprintf("The hot array holds %d x %d = %d bytes, which is %d lines, against %d lines for the full structs. A hot and cold split keeps most of the cache benefit while leaving related fields together.", n, hot, n*hot, hl, aos), facts)
}

// dbVarint returns the number of bytes in the unsigned LEB128 (7 bits per byte) encoding of v.
func dbVarint(v int) int {
	n := 1
	for v >= 128 {
		v >>= 7
		n++
	}
	return n
}

func dbZigzag(v int) int {
	if v >= 0 {
		return 2 * v
	}
	return -2*v - 1
}

func serialization(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		k := edakit.Scaled(ctx, 4, 7)
		var vs, work []string
		total := 0
		for i := 0; i < k; i++ {
			var v int
			switch edakit.Int(0, 3) {
			case 0:
				v = edakit.Int(0, 127)
			case 1:
				v = edakit.Int(128, 16383)
			case 2:
				v = edakit.Int(16384, 2097151)
			default:
				v = edakit.Int(2097152, 200000000)
			}
			b := dbVarint(v)
			total += b
			vs = append(vs, edakit.IntStr(v))
			work = append(work, edakit.Sprintf("%d takes %d", v, b))
		}
		q := edakit.Sprintf("A database writer stores unsigned integers as variable-length integers: 7 payload bits per byte, with the top bit of each byte marking that another byte follows. How many bytes does it write for the values %s?", strings.Join(vs, ", "))
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("A value below 2^7 = 128 takes 1 byte, below 2^14 = 16384 takes 2, below 2^21 = 2097152 takes 3, below 2^28 takes 4. Here %s, for %d bytes.", strings.Join(work, ", "), total),
			map[string]string{"bytes": edakit.IntStr(total)})
	case 2:
		k := edakit.Scaled(ctx, 5, 8)
		x := 1000 * edakit.Int(10, 900)
		var xs, work []string
		total := 0
		prev := 0
		for i := 0; i < k; i++ {
			if i > 0 {
				x += 10 * edakit.Pick([]int{19, 38, 57, 76, 95, 190, 1520})
			}
			d := x - prev
			b := dbVarint(d)
			total += b
			xs = append(xs, edakit.IntStr(x))
			work = append(work, edakit.Sprintf("%d (%d B)", d, b))
			prev = x
		}
		fixed := 4 * k
		q := edakit.Sprintf("The x coordinates of the cells in one row, in database units and sorted, are %s. A writer stores the first value and then each difference from the previous value, all as variable-length integers (7 payload bits per byte). How many bytes does this save against storing every coordinate as a fixed 4-byte integer?", strings.Join(xs, ", "))
		return edakit.Numeric(q, edakit.IntStr(fixed-total),
			edakit.Sprintf("Stored values: %s, %d bytes in all, against %d x 4 = %d fixed. Saving: %d bytes. Sorted coordinates have small gaps, which is why delta encoding pays.", strings.Join(work, ", "), total, k, fixed, fixed-total),
			map[string]string{"varint_bytes": edakit.IntStr(total), "fixed_bytes": edakit.IntStr(fixed)})
	}
	k := edakit.Scaled(ctx, 4, 7)
	var ds, work []string
	total := 0
	for i := 0; i < k; i++ {
		d := edakit.Pick([]int{1, 5, 30, 63, 64, 100, 700, 8191, 8192, 20000})
		if edakit.Int(0, 1) == 0 {
			d = -d
		}
		z := dbZigzag(d)
		b := dbVarint(z)
		total += b
		ds = append(ds, edakit.IntStr(d))
		work = append(work, edakit.Sprintf("%d maps to %d (%d B)", d, z, b))
	}
	q := edakit.Sprintf("Signed coordinate deltas are written with zigzag encoding, which maps n to 2n when n is at least 0 and to -2n - 1 when n is negative, and the result is stored as a variable-length integer (7 payload bits per byte). How many bytes do the deltas %s take?", strings.Join(ds, ", "))
	return edakit.Numeric(q, edakit.IntStr(total),
		edakit.Sprintf("%s, for %d bytes. Without zigzag, a negative number in two's complement would look huge and take the maximum number of bytes.", strings.Join(work, "; "), total),
		map[string]string{"bytes": edakit.IntStr(total)})
}

var designDatabaseItems = []edakit.Item{
	{KP: 1, Q: "A placer, a router and a timer each read the design from files written by the previous tool. What is the main cost a shared in-memory design database removes?",
		A: "Re-parsing and rebuilding the design between every step", D: []string{"The need for a netlist at all", "The cost of timing analysis itself", "The need for technology data such as layers"},
		Why: "With file handoff every tool parses LEF, DEF, Verilog and Liberty again and builds its own model. A shared database is built once and every engine reads and edits the same objects."},
	{KP: 1, Q: "Which three kinds of data does a physical design database keep, and keep apart?",
		A: "Technology data, library cell data and the design itself", D: []string{"Timing reports, logs and scripts", "RTL, testbenches and waveforms", "Only the netlist and its placement"},
		Why: "Technology data (layers, vias, sites, rules), library masters (cell abstracts and pins) and the design (instances, nets, wires) have different lifetimes and are shared differently across designs."},
	{KP: 1, Q: "A netlist object model and a hypergraph describe the same connectivity. What does the database add that the plain hypergraph lacks?",
		A: "Typed objects with names, masters, geometry and properties", D: []string{"Edges between every pair of pins", "A guarantee that the graph is acyclic", "Fewer terminals per net"},
		Why: "A hypergraph only says which vertices share a hyperedge. The database attaches names, master cells, locations, orientations, wires and attributes that every engine needs."},
	{KP: 1, Q: "A design database is described as the single source of truth for a flow. What does that require of the engines?",
		A: "They edit the shared objects instead of keeping private copies", D: []string{"They never modify the design", "They each write their own output format", "They run in separate processes"},
		Why: "If an engine keeps a private copy, the copy and the database drift apart. Engines read and write the shared objects, so a later engine sees every earlier change."},
	{KP: 2, Q: "A tool must answer: which nets connect to instance U17? Which representation makes this a direct lookup?",
		A: "Each instance keeps its list of instance terminals, each pointing to its net", D: []string{"A flat list of nets scanned for U17", "A text search through the DEF file", "A list of instances sorted by name only"},
		Why: "Bidirectional links, instance to terminals to nets and net to terminals to instances, make both directions constant time per terminal."},
	{KP: 2, Q: "A router needs to know which tracks and layers exist and what width a wire on metal 2 must have. Where does that come from in the database?",
		A: "The technology section, loaded from technology LEF", D: []string{"The netlist, loaded from Verilog", "The timing library, loaded from Liberty", "The constraints, loaded from SDC"},
		Why: "Layer definitions, widths, spacings, vias and sites are technology data. Liberty carries timing and power, SDC carries constraints, and Verilog carries connectivity."},
	{KP: 2, Q: "A database stores master cells once and gives each instance a reference to its master. Which query does this make cheap?",
		A: "Changing all NAND2 instances' footprint by editing one master", D: []string{"Finding the slack of an instance", "Computing the wirelength of a net", "Listing the clocks of the design"},
		Why: "Shared masters hold pin geometry, size and obstructions once. Instances only store what differs: location, orientation, connections."},
	{KP: 2, Q: "Which information is persistent design data rather than a derived result that an engine recomputes?",
		A: "Instance placement and routed wires", D: []string{"Arrival times at every pin", "Congestion estimates per gcell", "A wirelength gradient"},
		Why: "Placement and wires are the design. Arrival times, congestion maps and gradients are derived from it and are rebuilt by the engines that need them."},
	{KP: 3, Q: "A timer caches arrival times inside its own structures. The placer then moves 200 cells in the shared database. What must the architecture provide so the timer stays correct?",
		A: "Change notifications so the timer invalidates what the moves affect", D: []string{"A full reload of the design from DEF", "A second copy of the database for the timer", "Nothing, because arrival times do not depend on placement"},
		Why: "Derived results live outside the database, so the database must tell them what changed. Callbacks let the timer mark only the affected wires and arcs dirty."},
	{KP: 3, Q: "A team stores net names as std::string inside every terminal object. Memory explodes on a 50-million-terminal design. What is the standard fix?",
		A: "Store names once in a string table and keep only an index in the objects", D: []string{"Compress the DEF file", "Store names in the timing library", "Drop names entirely and use only coordinates"},
		Why: "Names are needed for input and output but not in inner loops. Interning them in one table removes duplication and keeps the hot objects small."},
	{KP: 3, Q: "Two engines run on the same database. One deletes a net while the other holds a reference to it. What property of the database design prevents a crash?",
		A: "References are stable handles that can be checked for validity", D: []string{"Nets are never deleted", "Each engine runs on its own copy", "Pointers are always 64 bits"},
		Why: "Handles such as table indices plus a validity check let an engine detect that an object is gone, where a raw pointer would dangle."},
	{KP: 3, Q: "A database is saved to disk and reloaded. Which saved form reloads fastest?",
		A: "A binary dump of the object tables, read back directly", D: []string{"A DEF file re-parsed through the full parser", "A Verilog netlist plus a re-placement", "A Tcl script that recreates every object"},
		Why: "A native binary format skips tokenizing, name resolution and validation. Text formats like DEF are for interchange, not for fast checkpoints."},
}

var opendbItems = []edakit.Item{
	{KP: 1, Q: "In OpenDB, what object represents one instance's connection point for one pin of its master?",
		A: "dbITerm", D: []string{"dbMTerm", "dbBTerm", "dbNet"},
		Why: "A dbMTerm is a pin of the master cell, a dbBTerm is a port of the block, and a dbITerm is the instance-specific terminal that a dbNet connects to."},
	{KP: 1, Q: "Which OpenDB object holds the design's instances, nets, rows and block terminals?",
		A: "dbBlock", D: []string{"dbTech", "dbLib", "dbMaster"},
		Why: "A dbBlock is the design container. Technology lives in dbTech, library cells in dbLib and dbMaster."},
	{KP: 1, Q: "OpenDB's schema closely follows which industry formats?",
		A: "LEF and DEF", D: []string{"Liberty and SDC", "GDSII and OASIS", "SPEF and SDF"},
		Why: "OpenDB was designed around LEF/DEF concepts: technology layers, sites, macros, pins, components, nets, special nets and wires."},
	{KP: 1, Q: "Where does OpenDB keep a metal layer's direction, pitch and minimum width?",
		A: "In a dbTechLayer under dbTech", D: []string{"In each dbNet", "In each dbInst", "In the dbBlock's rows"},
		Why: "Layer rules are technology data, read from technology LEF into dbTech and its dbTechLayer objects."},
	{KP: 2, Q: "A script wants every sink of net n1. Which traversal does the OpenDB object model support directly?",
		A: "From the dbNet to its dbITerms, then each dbITerm to its dbInst", D: []string{"From the dbNet to its dbMaster", "From dbTech to the net's layers", "From the dbBlock's rows to the net"},
		Why: "A net owns the set of instance terminals it connects. Each terminal knows its instance and its master terminal, which gives the pin direction."},
	{KP: 2, Q: "A plug-in reads a dbInst's location and wants the cell's width. Where does the width come from?",
		A: "The instance's dbMaster", D: []string{"The instance's dbITerms", "The block's DIEAREA", "The net the instance drives"},
		Why: "Size, pins and obstructions belong to the master. The instance stores placement, orientation and status."},
	{KP: 2, Q: "Which OpenDB object represents a port at the boundary of the block, such as a top-level input?",
		A: "dbBTerm", D: []string{"dbITerm", "dbMTerm", "dbRow"},
		Why: "A block terminal belongs to the block, connects to a net and may have pin shapes placed on the boundary by pin placement."},
	{KP: 2, Q: "An engine in OpenROAD wants to react when another engine moves an instance. What OpenDB mechanism exists for this?",
		A: "A block callback object registered on the block", D: []string{"Polling the DEF file for changes", "A timer interrupt", "Reading the Tcl history"},
		Why: "OpenDB lets a client register a callback object on a block so it is told about edits such as instance moves, net changes and terminal connections."},
	{KP: 3, Q: "OpenDB is a database of a single flattened block for most flows. What does a hierarchical Verilog netlist become when it is read in the usual OpenROAD flow?",
		A: "A flat block, with hierarchy kept only in instance names or optional hierarchy objects", D: []string{"One block per Verilog module, all fully linked", "A set of DEF files, one per module", "A GDS hierarchy"},
		Why: "Physical design in OpenROAD operates on a flat netlist. Hierarchical names survive, and newer support can retain module hierarchy, but the engines work on flat instances."},
	{KP: 3, Q: "Why does OpenDB distinguish dbWire from dbNet?",
		A: "A net is logical connectivity; a wire is its routed geometry", D: []string{"A wire is a net with only two pins", "Wires are only for clock nets", "A net stores timing, a wire stores names"},
		Why: "Connectivity exists before routing. The wire, written by the detailed router, encodes the shapes and vias that realize it."},
	{KP: 3, Q: "A tool saves an OpenDB checkpoint and another tool reloads it. Compared with writing and re-reading DEF, what is preserved more completely?",
		A: "The full database, including objects and properties that DEF cannot express", D: []string{"Only the netlist", "Only placement", "Nothing; they are identical by definition"},
		Why: "The native database file is a serialization of all objects. DEF is an interchange format with a fixed vocabulary."},
	{KP: 3, Q: "Special nets in OpenDB, from DEF SPECIALNETS, usually hold what?",
		A: "Power and ground nets with their stripes and rails", D: []string{"Clock nets after CTS", "Scan chains", "Nets with timing exceptions"},
		Why: "SPECIALNETS describe nets routed with special, usually prescribed geometry: power grids, rings and followpin rails."},
}

var memoryManagementItems = []edakit.Item{
	{KP: 1, Q: "A database hands out object references as 32-bit indices into per-type tables instead of raw pointers. Which benefit follows directly?",
		A: "References stay valid when a table is reallocated and are half the size", D: []string{"Objects can never be deleted", "Lookups no longer need memory", "Indices make multithreading automatic"},
		Why: "An index is relative to the table, so growing the table moves memory without breaking references, and a 32-bit index is half a 64-bit pointer."},
	{KP: 1, Q: "Why do EDA databases allocate millions of small objects from arenas or pools rather than one heap call each?",
		A: "Allocation is a pointer bump and objects of a kind stay contiguous", D: []string{"The operating system forbids small allocations", "Arenas make objects immutable", "Pools check design rules"},
		Why: "General-purpose allocation has per-call overhead and scatters objects. Arenas cut both, which helps speed and cache locality."},
	{KP: 1, Q: "An arena frees all its objects at once. Which use fits that pattern best?",
		A: "Temporary structures built for one pass, such as a router's search state", D: []string{"The persistent netlist of the design", "Technology layers", "Objects that must be freed one at a time in random order"},
		Why: "Arenas suit data with a common lifetime. The persistent netlist needs individual deletion, which a free list or table with slot reuse handles."},
	{KP: 1, Q: "A table reuses the slot of a deleted object. A stale handle to the old object now points at a new one. What prevents the confusion?",
		A: "A generation counter stored in the handle and the slot", D: []string{"Larger indices", "Sorting the table", "A mutex around every read"},
		Why: "Each slot increments its generation on reuse. A handle carries the generation it was issued with, so a mismatch reveals a stale reference."},
	{KP: 2, Q: "A legalizer stores raw pointers to cells in a vector, then the database grows its cell storage and moves objects. What happens?",
		A: "The legalizer's pointers dangle and later reads are undefined behaviour", D: []string{"The pointers are updated automatically", "The database refuses to grow", "Nothing; pointers are relative"},
		Why: "Raw pointers name an address. When storage moves, they still name the old address. Stable handles or pointer-stable storage avoid this."},
	{KP: 2, Q: "A tool leaks 200 bytes per ECO iteration and runs 2 million iterations in an overnight job. What does that cost?",
		A: "About 400 MB of memory that is never returned", D: []string{"About 400 KB", "About 4 GB", "Nothing, because the operating system frees it per iteration"},
		Why: "200 bytes x 2,000,000 = 400,000,000 bytes, about 400 MB. Leaks in long agent or ECO loops accumulate even when each one is small."},
	{KP: 2, Q: "Which data should be kept by value in a dense array rather than behind a pointer per element?",
		A: "Cell x and y coordinates read in every placement iteration", D: []string{"A rarely read property map per cell", "The name string of each net", "A per-design log file"},
		Why: "Hot, small, fixed-size data belongs in contiguous arrays. Rarely used or variable-size data can live behind an indirection."},
	{KP: 2, Q: "A C++ engine uses reference counting (shared ownership) for every netlist object. Profiling shows heavy slowdown in parallel loops. Why?",
		A: "Atomic increments and decrements on shared counters contend between threads", D: []string{"Reference counts make objects larger than a cache line", "Reference counting forbids reading objects", "The compiler removes the counts"},
		Why: "Copying a shared pointer touches an atomic counter. Many threads touching the same counters bounce cache lines. Plain handles or borrowed references avoid it."},
	{KP: 3, Q: "Deleting 30% of nets leaves holes in the net table. Iteration over nets now skips many empty slots. What is a standard remedy?",
		A: "Compact the table and remap handles at a safe point", D: []string{"Never delete nets", "Store nets in a linked list", "Double the table size"},
		Why: "Compaction restores density. Because handles change, it must happen when no engine holds old handles, or with a remapping table."},
	{KP: 3, Q: "An engine keeps a derived per-net array indexed by net index. A new net is created. What must the engine do?",
		A: "Grow its array to cover the new index before using it", D: []string{"Nothing; the array resizes itself in C++", "Rebuild the whole database", "Renumber every existing net"},
		Why: "Side arrays keyed by object index are cheap and cache friendly, but each owner must grow them, usually in response to a creation callback."},
	{KP: 3, Q: "A tool frees a large design and the process memory stays high. What is the most likely reason?",
		A: "The allocator keeps freed pages for reuse instead of returning them", D: []string{"The design was not really freed", "Memory cannot be freed in C++", "The operating system leaks memory"},
		Why: "Allocators often retain freed memory in their own pools. Measured resident memory may fall only after the allocator trims or the process exits."},
	{KP: 3, Q: "Which choice makes a database easiest to share read-only between worker threads?",
		A: "Freeze edits during the parallel phase and let workers read plain handles", D: []string{"Give each thread a reference-counted copy of every object", "Lock every object on each read", "Let each worker edit its own part freely"},
		Why: "Reads of immutable data need no synchronization. Edits are batched outside the parallel phase or partitioned so that workers never write the same object."},
}
