package pv

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// deviceRecognition: count gates from poly crossing diffusion (KP 1), W or L of a gate (KP 2),
// PMOS count from the well layer (KP 3).
func deviceRecognition(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		rows := edakit.Int(2, 3)
		var diffs []rect
		for j := 0; j < rows; j++ {
			a := edakit.Int(0, 100)
			diffs = append(diffs, rect{a, 300*j + 50, a + edakit.Int(250, 450), 300*j + 50 + edakit.Int(80, 160)})
		}
		np := edakit.Scaled(ctx, 3, 7)
		var polys []rect
		gates := 0
		for attempt := 0; len(polys) < np && attempt < 20*generator.MaxDrawAttempts; attempt++ {
			r1 := edakit.Int(0, rows-1)
			r2 := edakit.Int(r1, rows-1)
			x := edakit.Int(0, 560)
			p := rect{x, diffs[r1].y1 - edakit.Int(20, 60), x + 20, diffs[r2].y2 + edakit.Int(20, 60)}
			ok, g := true, 0
			for j := r1; j <= r2; j++ {
				d := diffs[j]
				inside := p.x1 > d.x1 && p.x2 < d.x2
				outside := p.x2 < d.x1 || p.x1 > d.x2
				if !inside && !outside {
					ok = false
				}
				if inside {
					g++
				}
			}
			clash := false
			for _, o := range polys {
				if overlap1(p.x1, p.x2, o.x1, o.x2) > -20 && overlap1(p.y1, p.y2, o.y1, o.y2) > 0 {
					clash = true
				}
			}
			if ok && !clash {
				polys = append(polys, p)
				gates += g
			}
		}
		if len(polys) == 0 {
			d := diffs[0]
			x := (d.x1 + d.x2) / 2
			polys = append(polys, rect{x, d.y1 - 30, x + 20, d.y2 + 30})
			gates = 1
		}
		np = len(polys)
		q := edakit.Sprintf("An LVS deck recognizes a transistor wherever POLY crosses DIFF completely (GATE = POLY AND DIFF). The diffusion regions are %s and the vertical poly lines are %s (nm). How many transistors does device recognition find?", listRects("D", diffs), listRects("P", polys))
		return edakit.Numeric(q, edakit.IntStr(gates),
			edakit.Sprintf("A poly line that spans several diffusion rows makes one gate in each row it crosses inside the row's x-range, and none where it passes beside a row. Counting crossings gives %d transistors from %d poly lines.", gates, np),
			map[string]string{"gates": edakit.IntStr(gates)})
	case 2:
		L := 5 * edakit.Int(3, 30)
		W := 10 * edakit.Int(10, 100)
		vertical := edakit.Int(0, 1) == 0
		var poly, diff rect
		if vertical {
			diff = rect{0, 0, edakit.Int(3, 6) * L * 2, W}
			px := edakit.Int(L, diff.x2-2*L)
			poly = rect{px, -edakit.Int(40, 100), px + L, W + edakit.Int(40, 100)}
		} else {
			diff = rect{0, 0, W, edakit.Int(3, 6) * L * 2}
			py := edakit.Int(L, diff.y2-2*L)
			poly = rect{-edakit.Int(40, 100), py, W + edakit.Int(40, 100), py + L}
		}
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("A poly shape %s crosses a diffusion shape %s (nm). What is the channel width W of the recognized transistor, in nm?", poly, diff)
			return edakit.Numeric(q, edakit.IntStr(W),
				edakit.Sprintf("The gate region is the intersection. Current flows across the poly between the two diffusion pieces, so L is the poly's width (%d nm) and W is the length of the gate along the poly, %d nm.", L, W),
				map[string]string{"W": edakit.IntStr(W), "L": edakit.IntStr(L)})
		}
		q := edakit.Sprintf("A poly shape %s crosses a diffusion shape %s (nm). What is the channel length L of the recognized transistor, in nm?", poly, diff)
		return edakit.Numeric(q, edakit.IntStr(L),
			edakit.Sprintf("The gate region is the intersection. L is measured in the direction of current flow, across the poly line, so it is the poly's width, %d nm; W along the poly is %d nm.", L, W),
			map[string]string{"W": edakit.IntStr(W), "L": edakit.IntStr(L)})
	default:
		nw := edakit.Int(1, 3)
		var wells []rect
		for i := 0; i < nw; i++ {
			x := 400 * i
			wells = append(wells, rect{x, edakit.Int(0, 100), x + edakit.Int(200, 350), edakit.Int(300, 500)})
		}
		n := edakit.Scaled(ctx, 4, 9)
		var parts []string
		pm := 0
		for i := 0; i < n; i++ {
			x, y := edakit.Int(0, 1200), edakit.Int(0, 600)
			in := false
			for _, w := range wells {
				if x > w.x1 && x < w.x2 && y > w.y1 && y < w.y2 {
					in = true
				}
			}
			if in {
				pm++
			}
			parts = append(parts, edakit.Sprintf("(%d, %d)", x, y))
		}
		q := edakit.Sprintf("Device recognition types a gate as PMOS when its channel lies inside NWELL, and NMOS otherwise (the substrate is p-type). The NWELL shapes are %s, and the gate channels are centred at %s (nm, no gate straddles a well edge). How many PMOS devices are recognized?", listRects("NW", wells), strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(pm),
			edakit.Sprintf("Test each channel centre for containment in a well: %d are inside, so %d PMOS and %d NMOS. A missing or misdrawn well layer silently changes device types, which LVS then reports as a type mismatch.", pm, pm, n-pm),
			map[string]string{"pmos": edakit.IntStr(pm)})
	}
}

type wire struct {
	vertical bool
	track    int
	lo, hi   int
	name     string
}

type pt struct{ x, y int }

// randomRouting draws M1 (horizontal) and M2 (vertical) centerline wires and vias at a subset
// of their crossings.
func randomRouting(ctx generator.GeneratorContext) ([]wire, []pt) {
	var ws []wire
	nh, nv := edakit.Scaled(ctx, 3, 6), edakit.Scaled(ctx, 3, 5)
	for i := 0; i < nh; i++ {
		lo := edakit.Int(0, 8)
		ws = append(ws, wire{false, edakit.Int(0, 8), lo, lo + edakit.Int(2, 6), edakit.Sprintf("h%d", i+1)})
	}
	for i := 0; i < nv; i++ {
		lo := edakit.Int(0, 8)
		ws = append(ws, wire{true, edakit.Int(0, 10), lo, lo + edakit.Int(2, 6), edakit.Sprintf("v%d", i+1)})
	}
	seen := map[pt]bool{}
	var vias []pt
	for _, h := range ws {
		if h.vertical {
			continue
		}
		for _, v := range ws {
			if !v.vertical {
				continue
			}
			if v.track >= h.lo && v.track <= h.hi && h.track >= v.lo && h.track <= v.hi {
				p := pt{v.track, h.track}
				if !seen[p] && edakit.Int(0, 1) == 0 {
					seen[p] = true
					vias = append(vias, p)
				}
			}
		}
	}
	return ws, vias
}

// connect builds the connectivity of wires and vias, skipping via index skip (-1 for none).
func connect(ws []wire, vias []pt, skip int) *dsu {
	d := newDSU(len(ws))
	for i := range ws {
		for j := i + 1; j < len(ws); j++ {
			a, b := ws[i], ws[j]
			if a.vertical == b.vertical && a.track == b.track && a.lo <= b.hi && b.lo <= a.hi {
				d.union(i, j)
			}
		}
	}
	for k, p := range vias {
		if k == skip {
			continue
		}
		hi, vi := -1, -1
		for i, w := range ws {
			on := false
			if w.vertical {
				on = w.track == p.x && p.y >= w.lo && p.y <= w.hi
			} else {
				on = w.track == p.y && p.x >= w.lo && p.x <= w.hi
			}
			if on && w.vertical && vi < 0 {
				vi = i
			}
			if on && !w.vertical && hi < 0 {
				hi = i
			}
		}
		if hi >= 0 && vi >= 0 {
			d.union(hi, vi)
		}
	}
	return d
}

func describeRouting(ws []wire, vias []pt) string {
	var hs, vs, ps []string
	for _, w := range ws {
		if w.vertical {
			vs = append(vs, edakit.Sprintf("%s at x = %d from y = %d to %d", w.name, w.track, w.lo, w.hi))
		} else {
			hs = append(hs, edakit.Sprintf("%s at y = %d from x = %d to %d", w.name, w.track, w.lo, w.hi))
		}
	}
	for _, p := range vias {
		ps = append(ps, edakit.Sprintf("(%d, %d)", p.x, p.y))
	}
	vtxt := "none"
	if len(ps) > 0 {
		vtxt = strings.Join(ps, ", ")
	}
	return edakit.Sprintf("Metal 1 wires (horizontal centerlines): %s. Metal 2 wires (vertical): %s. Vias between metal 1 and metal 2 at: %s.", strings.Join(hs, "; "), strings.Join(vs, "; "), vtxt)
}

// extraction: number of nets (KP 1), nets shorted by conflicting labels (KP 2), vias whose
// removal opens a net (KP 3).
func extraction(ctx generator.GeneratorContext) generator.Problem {
	ws, vias := randomRouting(ctx)
	desc := describeRouting(ws, vias)
	rule := "Two wires on the same layer and the same track connect when their ranges overlap or touch; a metal 1 and a metal 2 wire connect only through a via at a point on both, never by merely crossing."
	switch kp(ctx) {
	case 1:
		n := connect(ws, vias, -1).count()
		q := desc + " " + rule + " How many nets does connectivity extraction find?"
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("Union the wires that overlap on a track and the wire pairs joined by each via; %d wires collapse into %d connected components, each one net. Crossings without a via stay separate.", len(ws), n),
			map[string]string{"nets": edakit.IntStr(n), "wires": edakit.IntStr(len(ws))})
	case 2:
		d := connect(ws, vias, -1)
		names := []string{"A", "B", "Y"}
		var parts []string
		netLabels := map[int]map[string]bool{}
		for i := range ws {
			if edakit.Int(0, 1) == 0 {
				nm := edakit.Pick(names)
				parts = append(parts, edakit.Sprintf("%s on %s", nm, ws[i].name))
				r := d.find(i)
				if netLabels[r] == nil {
					netLabels[r] = map[string]bool{}
				}
				netLabels[r][nm] = true
			}
		}
		if len(parts) == 0 {
			parts = append(parts, "A on "+ws[0].name)
			netLabels[d.find(0)] = map[string]bool{"A": true}
		}
		shorts := 0
		for _, m := range netLabels {
			if len(m) > 1 {
				shorts++
			}
		}
		q := desc + " " + rule + edakit.Sprintf(" Text labels name the nets: %s. How many extracted nets carry two or more different label names (each reported as a short)?", strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(shorts),
			edakit.Sprintf("Extract the nets first, then collect the label names that land on each net. %d nets carry more than one distinct name. Repeated labels with the same name on one net are fine.", shorts),
			map[string]string{"shorts": edakit.IntStr(shorts)})
	default:
		base := connect(ws, vias, -1).count()
		crit := 0
		for k := range vias {
			if connect(ws, vias, k).count() > base {
				crit++
			}
		}
		q := desc + " " + rule + " A via is critical if deleting it alone splits a net into two. How many vias are critical?"
		return edakit.Numeric(q, edakit.IntStr(crit),
			edakit.Sprintf("Remove each via in turn and recount the nets (%d with every via present). %d vias are bridges of the connectivity graph; the rest have a parallel path, which is also what makes redundant vias valuable for yield.", base, crit),
			map[string]string{"nets": edakit.IntStr(base), "critical": edakit.IntStr(crit)})
	}
}

// reduction: devices after parallel merging (KP 1), merged width (KP 2), resistor value from
// squares (KP 3).
func reduction(ctx generator.GeneratorContext) generator.Problem {
	type dev struct {
		typ, g, s, d string
		w, l         int
	}
	switch k := kp(ctx); k {
	case 1, 2:
		gates := edakit.Shuffle([]string{"A", "B", "C", "D", "E"})
		nets := []string{"X", "Y", "Z", "VDD", "GND"}
		nl := edakit.Scaled(ctx, 2, 4)
		var devs []dev
		logical := 0
		type key struct{ typ, g, a, b string }
		groupW := map[key]int{}
		var firstKey key
		firstL := 0
		for i := 0; i < nl; i++ {
			typ := edakit.Pick([]string{"nmos", "pmos"})
			a, b := edakit.Pick(nets), edakit.Pick(nets)
			for b == a {
				b = edakit.Pick(nets)
			}
			L := edakit.Pick([]int{30, 40, 60})
			f := edakit.Int(1, 4)
			kk := key{typ, gates[i], a, b}
			if i == 0 {
				firstKey, firstL = kk, L
			}
			for j := 0; j < f; j++ {
				w := 10 * edakit.Int(10, 40)
				s, d := a, b
				if edakit.Int(0, 1) == 0 {
					s, d = b, a
				}
				devs = append(devs, dev{typ, gates[i], s, d, w, L})
				groupW[kk] += w
			}
			logical++
			if edakit.Int(0, 2) == 0 {
				devs = append(devs, dev{typ, gates[i], a, b, 10 * edakit.Int(10, 40), L + 20})
				logical++
			}
		}
		devs = edakit.Shuffle(devs)
		var parts []string
		for i, d := range devs {
			parts = append(parts, edakit.Sprintf("M%d %s g=%s s=%s d=%s W=%d L=%d", i+1, d.typ, d.g, d.s, d.d, d.w, d.l))
		}
		stem := edakit.Sprintf("An extracted netlist (W and L in nm) is: %s. The LVS tool merges devices in parallel: same type, same gate net, the same pair of source and drain nets in either order (source and drain are interchangeable), and equal L; the merged device's W is the sum.", strings.Join(parts, "; "))
		if k == 1 {
			q := stem + " How many devices remain after parallel merging?"
			return edakit.Numeric(q, edakit.IntStr(logical),
				edakit.Sprintf("Group the devices by (type, gate, unordered source/drain pair, L). There are %d groups; a device with a different L stays separate. Merging is what lets a 4-finger layout transistor match one schematic device.", logical),
				map[string]string{"merged": edakit.IntStr(logical)})
		}
		q := stem + edakit.Sprintf(" What is the W of the merged %s device with gate %s and L = %d, in nm?", firstKey.typ, firstKey.g, firstL)
		return edakit.Numeric(q, edakit.IntStr(groupW[firstKey]),
			edakit.Sprintf("Sum the W of every finger in that parallel group: %d nm. A device with the same gate but a different L is not part of it.", groupW[firstKey]),
			map[string]string{"W": edakit.IntStr(groupW[firstKey])})
	default:
		rs := edakit.Pick([]int{50, 100, 200, 400, 1000})
		n := edakit.Scaled(ctx, 1, 4)
		sq := 0.0
		var parts []string
		for i := 0; i < n; i++ {
			l, w := float64(edakit.Int(4, 40)), float64(edakit.Pick([]int{1, 2, 4}))/2
			sq += l / w
			parts = append(parts, edakit.Sprintf("%s by %s um", edakit.Num(l, 1), edakit.Num(w, 1)))
		}
		rl := float64(rs) * sq
		sch := rl * (1 + float64(edakit.Int(-12, 12))/100)
		sch = float64(int(sch/10+0.5)) * 10
		if sch <= 0 {
			sch = 10
		}
		pct := 100 * absf(rl-sch) / sch
		ans := edakit.Num(pct, 1)
		q := edakit.Sprintf("A poly resistor with sheet resistance %d ohm per square is drawn as %d segment(s) in series, each length by width: %s (contact and bend corrections ignored). The schematic resistor is %s ohm. By what percent does the extracted value differ from the schematic value, in absolute terms? Round to 1 decimal place.", rs, n, strings.Join(parts, ", "), edakit.Num(sch, 0))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Squares = sum of L / W = %s, so R = %d x %s = %s ohm. |%s - %s| / %s = %s percent. An LVS property check compares this against a tolerance set in the deck.", edakit.Num(sq, 2), rs, edakit.Num(sq, 2), edakit.Num(rl, 1), edakit.Num(rl, 1), edakit.Num(sch, 0), edakit.Num(sch, 0), ans),
			map[string]string{"squares": edakit.Num(sq, 2), "R": edakit.Num(rl, 1)})
	}
}

func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// mos is one transistor of a small CMOS netlist.
type mos struct {
	name, typ, g, s, d string
}

// cmosNetlist composes a few static CMOS gates into a transistor netlist. Inputs come from
// the primary inputs or earlier outputs. When independent is set, each gate gets its own
// primary input and output nets, which makes identical gates interchangeable.
func cmosNetlist(ngates int, independent bool, kinds []string) []mos {
	var ms []mos
	inputs := []string{"A", "B", "C", "D", "E", "F"}
	var avail []string
	inIdx := 0
	nextIn := func() string {
		s := inputs[inIdx%len(inputs)]
		inIdx++
		return s
	}
	internal := 0
	id := 0
	add := func(typ, g, s, d string) {
		id++
		ms = append(ms, mos{edakit.Sprintf("M%d", id), typ, g, s, d})
	}
	for gi := 0; gi < ngates; gi++ {
		kind := kinds[gi%len(kinds)]
		pickIn := func() string {
			if !independent && len(avail) > 0 && edakit.Int(0, 2) == 0 {
				return edakit.Pick(avail)
			}
			return nextIn()
		}
		out := edakit.Sprintf("Y%d", gi+1)
		switch kind {
		case "INV":
			a := pickIn()
			add("pmos", a, "VDD", out)
			add("nmos", a, "GND", out)
		case "NAND2":
			a, b := pickIn(), pickIn()
			for b == a {
				b = nextIn()
			}
			internal++
			x := edakit.Sprintf("n%d", internal)
			add("pmos", a, "VDD", out)
			add("pmos", b, "VDD", out)
			add("nmos", a, "GND", x)
			add("nmos", b, x, out)
		default:
			a, b := pickIn(), pickIn()
			for b == a {
				b = nextIn()
			}
			internal++
			x := edakit.Sprintf("n%d", internal)
			add("pmos", a, "VDD", x)
			add("pmos", b, x, out)
			add("nmos", a, "GND", out)
			add("nmos", b, "GND", out)
		}
		avail = append(avail, out)
	}
	return ms
}

func describeNetlist(ms []mos) string {
	var parts []string
	for _, m := range ms {
		parts = append(parts, edakit.Sprintf("%s %s g=%s s=%s d=%s", m.name, m.typ, m.g, m.s, m.d))
	}
	return strings.Join(parts, "; ")
}

// refineResult reports the class counts after each refinement round.
type refineResult struct {
	counts  []int // counts[0] is the initial partition
	classes map[string]int
	sizes   []int // sizes of the final classes
}

// refine runs Netgen-style colour refinement on the bipartite device-net graph. Devices start
// coloured by type, nets by degree, with VDD and GND in classes of their own. Each round
// recolours a vertex by its colour and the multiset of (terminal kind, neighbour colour), where
// gate is one kind and source/drain one interchangeable kind.
func refine(ms []mos) refineResult {
	type edge struct {
		to   int
		kind string
	}
	idx := map[string]int{}
	var names []string
	vertex := func(name string) int {
		if i, ok := idx[name]; ok {
			return i
		}
		idx[name] = len(names)
		names = append(names, name)
		return idx[name]
	}
	var adj [][]edge
	grow := func(n int) {
		for len(adj) < n {
			adj = append(adj, nil)
		}
	}
	isDev := map[int]string{}
	for _, m := range ms {
		dv := vertex("dev:" + m.name)
		isDev[dv] = m.typ
		for _, t := range []struct{ net, kind string }{{m.g, "g"}, {m.s, "sd"}, {m.d, "sd"}} {
			nv := vertex("net:" + t.net)
			grow(len(names))
			adj[dv] = append(adj[dv], edge{nv, t.kind})
			adj[nv] = append(adj[nv], edge{dv, t.kind})
		}
	}
	grow(len(names))
	col := make([]string, len(names))
	for i, nm := range names {
		switch {
		case isDev[i] != "":
			col[i] = "D" + isDev[i]
		case nm == "net:VDD" || nm == "net:GND":
			col[i] = nm
		default:
			col[i] = edakit.Sprintf("N%d", len(adj[i]))
		}
	}
	classCount := func(c []string) int {
		m := map[string]bool{}
		for _, x := range c {
			m[x] = true
		}
		return len(m)
	}
	res := refineResult{counts: []int{classCount(col)}}
	for round := 0; round < 20; round++ {
		// Compress colours to short ids so signatures stay small.
		ids := map[string]int{}
		for _, c := range col {
			if _, ok := ids[c]; !ok {
				ids[c] = len(ids)
			}
		}
		next := make([]string, len(col))
		for v := range col {
			var sig []string
			for _, e := range adj[v] {
				sig = append(sig, edakit.Sprintf("%s%d", e.kind, ids[col[e.to]]))
			}
			sort.Strings(sig)
			next[v] = edakit.Sprintf("%d|%s", ids[col[v]], strings.Join(sig, ","))
		}
		c := classCount(next)
		col = next
		if c == res.counts[len(res.counts)-1] {
			break
		}
		res.counts = append(res.counts, c)
	}
	sz := map[string]int{}
	for _, c := range col {
		sz[c]++
	}
	res.classes = sz
	for _, n := range sz {
		res.sizes = append(res.sizes, n)
	}
	return res
}

const refineRule = "Treat the netlist as a bipartite graph of devices and nets. Initially devices are coloured by type, VDD and GND each form a class of their own (they are matched by name), and every other net is coloured by its number of device terminals. In each round, every vertex gets a new colour made of its current colour plus the multiset of (terminal kind, neighbour colour) over its edges, where the terminal kinds are gate and source/drain (source and drain are interchangeable)."

// partitionRefinement: classes after one round (KP 1), rounds that split classes (KP 2),
// classes in the stable partition (KP 3).
func partitionRefinement(ctx generator.GeneratorContext) generator.Problem {
	ng := edakit.Scaled(ctx, 1, 3)
	kinds := edakit.Shuffle([]string{"INV", "NAND2", "NOR2", edakit.Pick([]string{"INV", "NAND2", "NOR2"})})
	ms := cmosNetlist(ng, false, kinds)
	r := refine(ms)
	stem := edakit.Sprintf("A transistor netlist is: %s. %s", describeNetlist(ms), refineRule)
	switch kp(ctx) {
	case 1:
		one := r.counts[0]
		if len(r.counts) > 1 {
			one = r.counts[1]
		}
		q := stem + edakit.Sprintf(" The initial partition has %d classes. How many classes are there after one round?", r.counts[0])
		return edakit.Numeric(q, edakit.IntStr(one),
			edakit.Sprintf("Recolouring by (colour, neighbour multiset) splits classes whose members see different neighbourhoods. The class counts by round are %s.", joinInts(r.counts)),
			map[string]string{"initial": edakit.IntStr(r.counts[0]), "after_one": edakit.IntStr(one)})
	case 2:
		rounds := len(r.counts) - 1
		q := stem + " Rounds continue until a round leaves the number of classes unchanged. How many rounds increase the number of classes?"
		return edakit.Numeric(q, edakit.IntStr(rounds),
			edakit.Sprintf("The class counts by round are %s, then a round with no change ends the refinement: %d splitting rounds. Each round costs time linear in the number of edges.", joinInts(r.counts), rounds),
			map[string]string{"counts": joinInts(r.counts)})
	default:
		final := r.counts[len(r.counts)-1]
		singles := 0
		for _, s := range r.sizes {
			if s == 1 {
				singles++
			}
		}
		q := stem + " Refinement stops when a round leaves the number of classes unchanged. How many classes does the stable partition have?"
		return edakit.Numeric(q, edakit.IntStr(final),
			edakit.Sprintf("The class counts by round are %s. The stable partition has %d classes, %d of them singletons; a singleton in both netlists is matched outright.", joinInts(r.counts), final, singles),
			map[string]string{"final": edakit.IntStr(final), "singletons": edakit.IntStr(singles)})
	}
}

// bipartiteSize: vertices or edges of the device-net graph.
func bipartiteSize(ctx generator.GeneratorContext) generator.Problem {
	ng := edakit.Scaled(ctx, 1, 3)
	ms := cmosNetlist(ng, false, edakit.Shuffle([]string{"INV", "NAND2", "NOR2"}))
	nets := map[string]bool{}
	for _, m := range ms {
		nets[m.g], nets[m.s], nets[m.d] = true, true, true
	}
	d, n := len(ms), len(nets)
	if edakit.Int(0, 1) == 0 {
		q := edakit.Sprintf("LVS models a netlist as a bipartite graph with one vertex per device and one per net, and one edge per device terminal (gate, source, drain; bulk ignored). For the netlist %s, how many vertices does the graph have?", describeNetlist(ms))
		why := edakit.Sprintf("%d devices plus %d distinct nets (VDD and GND included) = %d vertices. The graph has %d edges, three per transistor.", d, n, d+n, 3*d)
		return numChoice(q, float64(d+n), 0, "", why, float64(3*d), float64(d), float64(n+d+2))
	}
	q := edakit.Sprintf("LVS models a netlist as a bipartite graph with one vertex per device and one per net, and one edge per device terminal (gate, source, drain; bulk ignored). For the netlist %s, how many edges does the graph have?", describeNetlist(ms))
	why := edakit.Sprintf("Each of the %d transistors contributes three terminal edges: %d. The graph has %d net vertices.", d, 3*d, n)
	return numChoice(q, float64(3*d), 0, "", why, float64(d+n), float64(2*d), float64(4*d))
}

// ambiguousClasses: classes of size above one in a symmetric netlist.
func ambiguousClasses(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Int(2, 3)
	kind := edakit.Pick([]string{"INV", "NAND2", "NOR2"})
	kinds := make([]string, k)
	for i := range kinds {
		kinds[i] = kind
	}
	if edakit.Int(0, 1) == 0 {
		kinds = append(kinds, edakit.Pick([]string{"INV", "NAND2", "NOR2"}))
	}
	ms := cmosNetlist(len(kinds), true, kinds)
	r := refine(ms)
	amb, total := 0, 0
	for _, s := range r.sizes {
		total++
		if s > 1 {
			amb++
		}
	}
	q := edakit.Sprintf("A layout netlist is: %s. %s Refinement runs until the number of classes stops changing. Only VDD and GND are matched by name; no other port names are used. How many classes of the stable partition still contain more than one vertex?", describeNetlist(ms), refineRule)
	why := edakit.Sprintf("Identical gates with private inputs and outputs are automorphic images of each other, so refinement can never separate them. The stable partition has %d classes, %d of them with more than one member; the matcher must break each tie by an arbitrary choice and refine again.", total, amb)
	return numChoice(q, float64(amb), 0, "", why, float64(total), 0, float64(len(ms)/k))
}

// ercFloating: nets that connect only to gates and are not ports.
func ercFloating(ctx generator.GeneratorContext) generator.Problem {
	ng := edakit.Scaled(ctx, 2, 3)
	ms := cmosNetlist(ng, false, edakit.Shuffle([]string{"INV", "NAND2", "NOR2"}))
	gateOnly := map[string]bool{}
	sd := map[string]bool{}
	for _, m := range ms {
		gateOnly[m.g] = true
		sd[m.s], sd[m.d] = true, true
	}
	var cand []string
	for n := range gateOnly {
		if !sd[n] {
			cand = append(cand, n)
		}
	}
	sort.Strings(cand)
	var ports []string
	floating := 0
	for _, n := range cand {
		if edakit.Int(0, 2) > 0 {
			ports = append(ports, n)
		} else {
			floating++
		}
	}
	portTxt := "none"
	if len(ports) > 0 {
		portTxt = strings.Join(ports, ", ")
	}
	q := edakit.Sprintf("An ERC deck flags a floating gate net: a net that reaches only gate terminals and is not a port of the cell. The netlist is %s. Besides VDD, GND and the outputs, the cell's input ports are: %s. How many floating gate nets are flagged?", describeNetlist(ms), portTxt)
	why := edakit.Sprintf("The nets that touch only gates are %s. Those not listed as input ports have no driver at all: %d floating nets, which would leave both transistors partly on.", strings.Join(cand, ", "), floating)
	return numChoice(q, float64(floating), 0, "", why, float64(len(cand)), float64(len(ports)), float64(floating+1))
}

// ercWellTaps: wells with no tap.
func ercWellTaps(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 6)
	var wells []rect
	for i := 0; i < n; i++ {
		x := 500 * i
		wells = append(wells, rect{x, edakit.Int(0, 50), x + edakit.Int(200, 400), edakit.Int(250, 450)})
	}
	var taps []string
	tapped := make([]bool, n)
	nt := edakit.Int(n-1, n+3)
	for i := 0; i < nt; i++ {
		x, y := edakit.Int(0, 500*n), edakit.Int(0, 450)
		if edakit.Int(0, 1) == 0 {
			w := edakit.Pick(wells)
			x, y = edakit.Int(w.x1+1, w.x2-1), edakit.Int(w.y1+1, w.y2-1)
		}
		for j, w := range wells {
			if x > w.x1 && x < w.x2 && y > w.y1 && y < w.y2 {
				tapped[j] = true
			}
		}
		taps = append(taps, edakit.Sprintf("(%d, %d)", x, y))
	}
	bad := 0
	for _, t := range tapped {
		if !t {
			bad++
		}
	}
	q := edakit.Sprintf("An ERC check requires every NWELL region to contain at least one N+ well tap tying it to VDD. The wells are %s and the taps sit at %s (nm). How many wells are flagged as untapped?", listRects("NW", wells), strings.Join(taps, ", "))
	why := edakit.Sprintf("Test each tap for containment. %d wells contain no tap; an untapped well floats, which shifts PMOS thresholds and raises latch-up risk.", bad)
	return numChoice(q, float64(bad), 0, "", why, float64(n-bad), float64(imax(0, n-nt)), float64(bad+1))
}
