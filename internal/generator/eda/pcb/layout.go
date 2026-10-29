package pcb

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type pt struct{ x, y int }

func manhattan(a, b pt) int {
	dx, dy := a.x-b.x, a.y-b.y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// mstLength is Prim's algorithm on the complete Manhattan graph.
func mstLength(ps []pt) int {
	n := len(ps)
	in := make([]bool, n)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = 1 << 30
	}
	dist[0] = 0
	total := 0
	for k := 0; k < n; k++ {
		best := -1
		for i := 0; i < n; i++ {
			if !in[i] && (best < 0 || dist[i] < dist[best]) {
				best = i
			}
		}
		in[best] = true
		total += dist[best]
		for i := 0; i < n; i++ {
			if d := manhattan(ps[best], ps[i]); !in[i] && d < dist[i] {
				dist[i] = d
			}
		}
	}
	return total
}

// ratsnest is the computed KP 2 of layout.placement.
func ratsnest(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	var ps []pt
	seen := map[pt]bool{}
	for a := 0; len(ps) < n && a < generator.MaxDrawAttempts; a++ {
		p := pt{edakit.Int(0, 40), edakit.Int(0, 30)}
		if !seen[p] {
			seen[p] = true
			ps = append(ps, p)
		}
	}
	for len(ps) < n {
		ps = append(ps, pt{41 + len(ps), 31})
	}
	mst := mstLength(ps)
	star, chain := 0, 0
	minx, maxx, miny, maxy := ps[0].x, ps[0].x, ps[0].y, ps[0].y
	for i := 1; i < n; i++ {
		star += manhattan(ps[0], ps[i])
		chain += manhattan(ps[i-1], ps[i])
		minx, maxx = min(minx, ps[i].x), max(maxx, ps[i].x)
		miny, maxy = min(miny, ps[i].y), max(maxy, ps[i].y)
	}
	hpwl := (maxx - minx) + (maxy - miny)
	var pins []string
	for i, p := range ps {
		pins = append(pins, edakit.Sprintf("P%d (%d, %d)", i+1, p.x, p.y))
	}
	ans, ds := options(float64(mst), 0, "mm", float64(star), float64(chain), float64(hpwl))
	q := "A board tool draws each net's ratsnest as a minimum spanning tree over its pins, with Manhattan distance between pins. One net has pins at " + strings.Join(pins, ", ") + ", in mm. What is the ratsnest length of this net?"
	return edakit.Choice(q, ans, edakit.Sprintf("Prim's algorithm on the Manhattan distances gives %d mm. A star from P1 gives %d mm, chaining the pins in listed order %d mm, and the bounding-box half-perimeter %d mm, which only lower-bounds the tree.", mst, star, chain, hpwl), ds...)
}

var placementItems = []edakit.Item{
	{KP: 1, Q: "Which components are usually placed first on a board?",
		A:   "Mechanically fixed parts: connectors, mounting holes, switches and displays set by the enclosure",
		D:   []string{"Decoupling capacitors", "Pull-up resistors", "Test points"},
		Why: "The enclosure fixes these positions. Everything else is placed around them."},
	{KP: 1, Q: "Where should a decoupling capacitor be placed?",
		A:   "As close as possible to the power pin it serves, with a short low-inductance path to its vias and plane",
		D:   []string{"Near the board edge to save space", "Next to the voltage regulator only", "Anywhere on the net, since the plane is ideal"},
		Why: "The capacitor's benefit at high frequency is limited by the loop inductance of its connection, which grows with distance."},
	{KP: 1, Q: "Why is board placement a larger share of the designer's effort than chip placement is?",
		A:   "It fixes routability, signal integrity, thermal and mechanical fit for a few hundred varied parts, largely by hand",
		D:   []string{"Boards hold more components than chips", "Board placers are always analytical", "Routing is fully automatic on boards"},
		Why: "Board placement is a constrained, heterogeneous problem where designer judgment and many non-wirelength objectives dominate."},
	{KP: 1, Q: "How does cross-probing between schematic and layout help placement?",
		A:   "Selecting a functional block in the schematic highlights or gathers its parts on the board",
		D:   []string{"It simulates the schematic", "It writes the Gerber files", "It computes trace impedance"},
		Why: "Grouping parts by function keeps their local nets short and makes the board easier to route and debug."},
	{KP: 3, Q: "A switching regulator's inductor is placed beside an analog front end. What problem is likely?",
		A:   "Magnetic and switching-noise coupling into the sensitive analog inputs",
		D:   []string{"Higher trace resistance", "A DRC clearance error", "Better thermal behavior"},
		Why: "The inductor's field and the regulator's fast current edges couple into nearby high-impedance nodes. Separation and orientation are placement decisions."},
	{KP: 3, Q: "Two hot parts are placed side by side near the board center. What should the designer check?",
		A:   "The combined thermal load and airflow, since heat spreads through copper and raises both junction temperatures",
		D:   []string{"Nothing, since each part's datasheet rating applies alone", "Only the silkscreen overlap", "The impedance of the traces between them"},
		Why: "Thermal resistance to ambient is shared. Clustering heat sources raises the local board temperature for both."},
	{KP: 3, Q: "A placement cuts total ratsnest length by 15% but moves a DDR memory 40 mm farther from the processor. Why might it be worse?",
		A:   "Ratsnest length weights all nets equally, while the DDR nets carry tight timing and signal-integrity limits",
		D:   []string{"Ratsnest length is not a usable metric", "A shorter ratsnest always means a better board", "DDR nets have no length limits"},
		Why: "Total wirelength is a proxy. The real objective weights critical nets, so a placer needs per-net weights or constraints."},
	{KP: 3, Q: "Why do placement checks use courtyards rather than component bodies?",
		A:   "Courtyards include the pads and an assembly margin, so clearing them leaves room for nozzles, inspection and rework",
		D:   []string{"Bodies are not stored in footprints", "Courtyards define impedance", "Component bodies are meant to overlap"},
		Why: "Two bodies can clear each other while their pads or the placement head cannot. The courtyard captures the real keep-out."},
	{KP: 3, Q: "A board is reflowed twice, top side second. Which placement rule follows for the bottom side?",
		A:   "Put light parts on the bottom so surface tension holds them during the second reflow",
		D:   []string{"Put the heaviest parts on the bottom", "Allow no parts on the bottom", "Use only through-hole parts on the bottom"},
		Why: "During the second pass the bottom joints melt again while the board is upside down. Heavy parts there need glue or must move to the top."},
}

var autorouterItems = []edakit.Item{
	{KP: 1, Q: "How do most modern board autorouters represent the routing space, compared with chip detailed routers?",
		A:   "Shape-based (gridless) geometry with arbitrary pad shapes and 45-degree or any-angle segments",
		D:   []string{"A uniform gcell grid with one preferred direction per layer", "A channel graph between rows of cells", "Orthogonal segments on a fixed track grid only"},
		Why: "Board obstacles are irregular pads and keep-outs with mixed rules, so board routers work on exact shapes rather than a uniform track grid."},
	{KP: 1, Q: "Why can a board router not give each layer a single preferred direction the way chip routers do?",
		A:   "Boards have few signal layers and large obstacles such as BGAs and connectors, so each layer must route in any direction",
		D:   []string{"Fabrication forbids vertical traces", "The IPC class fixes the directions", "Boards have hundreds of layers"},
		Why: "With 2 to perhaps a few dozen layers and big blockages, a strict horizontal-vertical split wastes too much capacity."},
	{KP: 1, Q: "What does a push-and-shove interactive router do?",
		A:   "Moves existing traces and vias aside, within the clearance rules, to make room for the trace being drawn",
		D:   []string{"Deletes all existing traces before routing", "Runs a batch autoroute overnight", "Pushes components apart in placement"},
		Why: "It keeps the board DRC-clean while the designer draws, which is how most critical board routing is done."},
	{KP: 1, Q: "Which early step dominates routing a large BGA board?",
		A:   "Fanout and escape: bringing every ball out to a via and a routing channel",
		D:   []string{"Placing the silkscreen", "Drilling the mounting holes", "Choosing the copper weight"},
		Why: "Escape determines how many layers the board needs. Once the balls are out, the remaining routing is comparatively open."},
	{KP: 2, Q: "A board autorouter routes all nets, rips up those in conflict and reroutes with higher costs in crowded areas, over many passes. Which chip-routing idea is this?",
		A:   "Rip-up and reroute with negotiated congestion costs",
		D:   []string{"Left-edge channel routing", "Simulated-annealing placement", "Static timing analysis"},
		Why: "Board routers such as the Specctra lineage use the same negotiation idea: shared resources get more expensive until nets spread out."},
	{KP: 2, Q: "A maze router runs on a multilayer board graph. How should it treat a via?",
		A:   "As an extra cost on the layer-change edge, so it prefers fewer vias at equal length",
		D:   []string{"As free, since vias cost nothing", "As forbidden, with infinite cost", "As a cost equal to the board thickness only"},
		Why: "Vias add fabrication cost, block other layers and add discontinuities, so routers price them explicitly."},
	{KP: 2, Q: "Why does a board router handle power and ground nets differently from signals?",
		A:   "They use planes or wide pours, so pins connect to them through vias rather than routed thin traces",
		D:   []string{"Power nets are ignored entirely", "They must be routed first as thin traces", "IPC forbids power planes"},
		Why: "On multilayer boards most power distribution is solid copper. The router only fans pins out to the plane."},
	{KP: 2, Q: "A topological router first decides how each net threads between obstacles and only then fixes geometry. What does that buy?",
		A:   "It separates which side of each obstacle a wire passes from its exact shape, so wires can be spread and smoothed afterwards",
		D:   []string{"It makes DRC unnecessary", "It removes all vias", "It guarantees minimum wirelength"},
		Why: "Choosing the homotopy class first is a coarse decision; geometry is then a continuous spreading problem."},
	{KP: 3, Q: "An autorouter reports 100% completion. What has it not shown?",
		A:   "That impedance, length-matching, return-path and crosstalk constraints are met",
		D:   []string{"That every net is connected", "That the clearance rules are met", "That the board has copper"},
		Why: "Completion means connected and clearance-clean. Electrical constraints need their own checks and simulation."},
	{KP: 3, Q: "Why do designers route critical nets by hand and lock them before autorouting?",
		A:   "Locked routes keep their topology and matching while the autorouter fills in the remaining nets",
		D:   []string{"Autorouters cannot route signals", "Locked routes are cheaper to fabricate", "DRC skips locked routes"},
		Why: "Locking turns the critical routes into fixed obstacles, so the router cannot trade them away for completion."},
	{KP: 3, Q: "A differential pair is autorouted as two independent nets. What goes wrong?",
		A:   "The traces take different paths and lengths, breaking their coupling and adding skew",
		D:   []string{"Nothing goes wrong", "The two nets short together", "The pair's impedance exactly doubles"},
		Why: "A pair must be routed as one object with a fixed gap so that its differential impedance and timing hold."},
	{KP: 3, Q: "What do net classes hold in a board tool?",
		A:   "Per-group rules such as trace width, clearance, via size and differential-pair gap",
		D:   []string{"The BOM lines of the group", "The Gerber aperture list", "The pick-and-place rotations"},
		Why: "Net classes attach routing rules to groups of nets, which the router and DRC both read."},
}

// bgaEscape: traces between pads (KP 1), layers to escape a full array (KP 2), balls in outer rings (KP 3).
func bgaEscape(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		for a := 0; ; a++ {
			p := edakit.Pick([]int{400, 500, 650, 800, 1000, 1270})
			d := 10 * ((p/2 + edakit.Int(-5, 5)*10) / 10)
			w := edakit.Pick([]int{75, 90, 100, 125, 150})
			s := edakit.Pick([]int{75, 90, 100, 125})
			if a >= generator.MaxDrawAttempts {
				p, d, w, s = 1000, 500, 100, 100
			} else if p-d-s < w+s {
				continue
			}
			n := (p - d - s) / (w + s)
			q := edakit.Sprintf("A BGA has a %d micrometre ball pitch and %d micrometre pads. Traces are %d micrometres wide with %d micrometre clearance to pads and to each other. $N$ traces fit between two adjacent pads when $N w + (N + 1) s \\le p - d$. How many traces fit between adjacent pads?", p, d, w, s)
			return edakit.Numeric(q, edakit.IntStr(n), edakit.Sprintf("The gap is $%d - %d = %d$ micrometres. $N \\le (%d - %d)/(%d + %d) = %s$, so $N = %d$.", p, d, p-d, p-d, s, w, s, edakit.Num(float64(p-d-s)/float64(w+s), 3), n), map[string]string{"gap": edakit.IntStr(p - d)})
		}
	case 2:
		n := edakit.Scaled(ctx, 10, 34)
		N := edakit.Int(1, 2)
		rings := ceilInt(float64(n) / 2)
		layers := ceilInt(float64(rings) / float64(N+1))
		q := edakit.Sprintf("A full %d x %d BGA array must be escaped. %s between adjacent balls, and assume each signal layer escapes $N + 1$ rows of the array from the outside in, where $N$ is the number of traces between balls. How many signal layers are needed to escape every ball?", n, n, plural(N, "trace fits", "traces fit"))
		return edakit.Numeric(q, edakit.IntStr(layers), edakit.Sprintf("A %d x %d array has %d concentric rows. Each layer escapes %d of them, so $\\lceil %d / %d \\rceil = %d$ layers.", n, n, rings, N+1, rings, N+1, layers), map[string]string{"rings": edakit.IntStr(rings)})
	}
	n := edakit.Scaled(ctx, 8, 30)
	k := edakit.Int(1, min(4, n/2-1))
	inner := n - 2*k
	ans := n*n - inner*inner
	q := edakit.Sprintf("A full %d x %d BGA array is escaped on the top layer for its outer %s of balls. How many balls lie in those outer rows?", n, n, plural(k, "row", "rows"))
	return edakit.Numeric(q, edakit.IntStr(ans), edakit.Sprintf("Remove the inner $%d \\times %d$ block from the $%d \\times %d$ array: $%d - %d = %d$ balls.", inner, inner, n, n, n*n, inner*inner, ans), nil)
}
