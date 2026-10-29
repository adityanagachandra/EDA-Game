package route

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type iv struct{ l, r int }

func randIntervals(n, span int) []iv {
	out := make([]iv, n)
	for i := range out {
		l := edakit.Int(0, span-2)
		out[i] = iv{l, minI(span, l+edakit.Int(1, span/2))}
	}
	return out
}

func ivText(xs []iv, prefix string) string {
	var parts []string
	for i, x := range xs {
		parts = append(parts, edakit.Sprintf("%s%d [%d, %d]", prefix, i+1, x.l, x.r))
	}
	return strings.Join(parts, ", ")
}

func maxOverlap(xs []iv) int {
	best := 0
	for c := 0; c <= 40; c++ {
		n := 0
		for _, x := range xs {
			if x.l <= c && c <= x.r {
				n++
			}
		}
		best = maxI(best, n)
	}
	return best
}

func overlapLen(a, b iv) int { return maxI(0, minI(a.r, b.r)-maxI(a.l, b.l)+1) }

// trackAssignment: tracks a panel needs (KP 1), conflict pairs (KP 2), overlap cost of a given
// assignment (KP 3).
func trackAssignment(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	xs := randIntervals(n, 20)
	intro := edakit.Sprintf("A track-assignment panel is one row of gcells, numbered 0 to 20. Global routing left these horizontal segments in it, each covering an inclusive range of gcells: %s. Two segments conflict if they share a gcell. ", ivText(xs, "s"))
	switch kp(ctx) {
	case 1:
		m := maxOverlap(xs)
		return edakit.Numeric(intro+"What is the fewest tracks that can hold all the segments without two conflicting segments on the same track?", edakit.IntStr(m),
			edakit.Sprintf("The conflict graph is an interval graph, which is perfect: its chromatic number equals its largest clique, the most segments covering one gcell. That maximum is %d, and a left-edge sweep achieves it.", m),
			map[string]string{"tracks": edakit.IntStr(m)})
	case 2:
		c := 0
		for i := range xs {
			for j := i + 1; j < len(xs); j++ {
				if overlapLen(xs[i], xs[j]) > 0 {
					c++
				}
			}
		}
		return edakit.Numeric(intro+"How many pairs of segments conflict?", edakit.IntStr(c),
			edakit.Sprintf("Check each of the %d pairs for a shared gcell: %d do. These are the edges a track assigner must keep apart.", n*(n-1)/2, c),
			map[string]string{"conflicts": edakit.IntStr(c)})
	default:
		T := edakit.Int(2, 3)
		tr := make([]int, n)
		var parts []string
		for i := range tr {
			tr[i] = edakit.Int(1, T)
			parts = append(parts, edakit.Sprintf("s%d on track %d", i+1, tr[i]))
		}
		cost := 0
		for i := range xs {
			for j := i + 1; j < len(xs); j++ {
				if tr[i] == tr[j] {
					cost += overlapLen(xs[i], xs[j])
				}
			}
		}
		q := intro + edakit.Sprintf("A track assigner with only %d tracks proposes: %s. Its cost counts, for every pair of segments on the same track, the number of gcells they share. What is the total overlap cost?", T, strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(cost),
			edakit.Sprintf("Summing shared gcells over same-track pairs gives %d. Detailed routing must later resolve every overlap with a detour or a reroute, so track assignment minimises this cost.", cost),
			map[string]string{"overlap": edakit.IntStr(cost)})
	}
}

// pinAccess: on-track access points of a pin (KP 1), compatible access pairs for two pins
// (KP 2), access points once half-track positions are allowed (KP 3).
func pinAccess(ctx generator.GeneratorContext) generator.Problem {
	p := edakit.Pick([]int{40, 48, 56, 64, 80})
	o := p / 2
	e := edakit.Pick([]int{10, 15, 20})
	count := func(x1, x2 int, half bool) []int {
		var out []int
		step := p
		start := o
		if half {
			step = p / 2
			start = o % step
		}
		for x := start; x <= x2; x += step {
			if x >= x1+e && x <= x2-e {
				out = append(out, x)
			}
		}
		return out
	}
	x1 := 10 * edakit.Int(5, 40)
	x2 := x1 + 10*edakit.Int(p/10+5, 4*p/10)
	intro := edakit.Sprintf("Vertical M2 tracks sit at x = %d + %dk nm. A via from an M1 pin to M2 lands on a track only if the via's metal fits inside the pin with %d nm to spare on each side, that is at $x_1 + %d \\le x \\le x_2 - %d$ for a pin spanning $[x_1, x_2]$. ", o, p, e, e, e)
	switch kp(ctx) {
	case 1:
		a := count(x1, x2, false)
		q := intro + edakit.Sprintf("A pin spans x from %d to %d nm. How many on-track access points does it have?", x1, x2)
		return edakit.Numeric(q, edakit.IntStr(len(a)),
			edakit.Sprintf("Tracks with %d <= x <= %d: %s. That is %d access points.", x1+e, x2-e, intsText(a), len(a)),
			map[string]string{"points": edakit.IntStr(len(a))})
	case 2:
		var a, b []int
		var y1, y2 int
		for try := 0; try < generator.MaxDrawAttempts; try++ {
			a = count(x1, x2, false)
			y1 = x2 + 10*edakit.Int(1, 4)
			y2 = y1 + 10*edakit.Int(p/10+5, 4*p/10)
			b = count(y1, y2, false)
			if len(a) > 0 && len(b) > 0 {
				break
			}
			x2 += p
		}
		d := p * edakit.Int(1, 2)
		pairs := 0
		for _, xa := range a {
			for _, xb := range b {
				if absI(xa-xb) >= d {
					pairs++
				}
			}
		}
		q := intro + edakit.Sprintf("Neighbouring pins A, spanning %d to %d nm, and B, spanning %d to %d nm, each need one access via, and the two vias must be at least %d nm apart. How many (A, B) pairs of access points are compatible?", x1, x2, y1, y2, d)
		return edakit.Numeric(q, edakit.IntStr(pairs),
			edakit.Sprintf("A's points are %s and B's are %s. Of the %d pairs, %d keep the vias %d nm apart. Pin-access analysis enumerates these choices per cell before routing so the router never boxes itself in.", intsText(a), intsText(b), len(a)*len(b), pairs, d),
			map[string]string{"pairs": edakit.IntStr(pairs)})
	default:
		a := count(x1, x2, false)
		b := count(x1, x2, true)
		q := intro + edakit.Sprintf("A pin spans x from %d to %d nm. A pin-access engine also allows off-track points halfway between tracks, which a jog then connects to a track. How many access points does the pin have in total?", x1, x2)
		return edakit.Numeric(q, edakit.IntStr(len(b)),
			edakit.Sprintf("Candidate positions now repeat every %d nm. Those within [%d, %d] are %s: %d points, against %d on track only. Off-track points rescue pins narrower than a pitch at the price of a jog that must itself be DRC-clean.", p/2, x1+e, x2-e, intsText(b), len(b), len(a)),
			map[string]string{"points": edakit.IntStr(len(b)), "on_track": edakit.IntStr(len(a))})
	}
}

// drcAware: cost of a route (KP 1), best of three (KP 2), tracks blocked by a wide wire (KP 3).
func drcAware(ctx generator.GeneratorContext) generator.Problem {
	wl, cv, cd := 1, edakit.Int(2, 6), edakit.Pick([]int{20, 50, 100, 200})
	type cand struct{ L, v, m int }
	mk := func() cand { return cand{edakit.Int(10, 40), edakit.Int(1, 6), edakit.Int(0, 2)} }
	cost := func(c cand) int { return wl*c.L + cv*c.v + cd*c.m }
	rule := edakit.Sprintf("A DRC-aware detailed router prices a route as %d per track step, %d per via and %d per design-rule violation it would create (a short, a spacing or an end-of-line violation). ", wl, cv, cd)
	switch kp(ctx) {
	case 1:
		c := mk()
		return edakit.Numeric(rule+edakit.Sprintf("A candidate route has %d steps, %s and %s. What is its cost?", c.L, plural(c.v, "via"), plural(c.m, "violation")), edakit.IntStr(cost(c)),
			edakit.Sprintf("%d + %d x %d + %d x %d = %d.", c.L, cv, c.v, cd, c.m, cost(c)),
			map[string]string{"cost": edakit.IntStr(cost(c))})
	case 2:
		cs := []cand{mk(), mk(), mk()}
		cs[0].m = 0
		cs[0].L += 20
		var parts []string
		best := math.MaxInt32
		for i, c := range cs {
			parts = append(parts, edakit.Sprintf("route %c: %d steps, %d vias, %d violations", 'A'+i, c.L, c.v, c.m))
			best = minI(best, cost(c))
		}
		return edakit.Numeric(rule+"Three candidates join the same pins: "+strings.Join(parts, "; ")+". What is the cost of the route the router chooses?", edakit.IntStr(best),
			edakit.Sprintf("The costs are %d, %d and %d; the router takes %d. With a large violation penalty, a longer clean route usually wins, which is how the router trades wirelength for legality.", cost(cs[0]), cost(cs[1]), cost(cs[2]), best),
			map[string]string{"best": edakit.IntStr(best)})
	default:
		p := edakit.Pick([]int{40, 48, 56, 64})
		w := p / 2
		W := w * edakit.Int(3, 8)
		s := p - w
		sw := s + 10*edakit.Int(2, 8)
		k := 0
		for j := 1; j < 50; j++ {
			if float64(j*p)-float64(W)/2-float64(w)/2 < float64(sw) {
				k++
			}
		}
		q := edakit.Sprintf("Tracks are %d nm apart and normal wires are %d nm wide, so normal neighbours sit %d nm apart, exactly the minimum spacing. A power-like wire %d nm wide is centred on one track. Next to a wire this wide the spacing table requires %d nm. On one side of it, how many adjacent tracks can no longer hold a normal wire?", p, w, s, W, sw)
		return edakit.Numeric(q, edakit.IntStr(k),
			edakit.Sprintf("A wire on the k-th track over has an edge gap of %dk - %d/2 - %d/2 nm. It is illegal while that gap is below %d nm, which holds for k = 1 to %d. The router models this by blocking those tracks, not by checking afterwards.", p, W, w, sw, k),
			map[string]string{"blocked": edakit.IntStr(k)})
	}
}

// minAreaEOL: extension for minimum area (KP 1), end-of-line violations on a track (KP 2),
// the stub length a via landing needs (KP 3).
func minAreaEOL(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		w := edakit.Pick([]int{20, 24, 28, 32, 40, 50})
		A := 1000 * edakit.Int(3, 20)
		need := ceilDiv(A, w)
		l := edakit.Int(w, need-1)
		q := edakit.Sprintf("A metal layer's minimum-area rule requires every isolated shape to have area at least %d $\\text{nm}^2$. A detailed router leaves a %d nm wide stub only %d nm long. By how many nm, rounded up to a whole nm, must it extend the stub?", A, w, l)
		return edakit.Numeric(q, edakit.IntStr(need-l),
			edakit.Sprintf("The stub needs length at least %d / %d = %s nm, so %d nm after rounding up, an extension of %d nm.", A, w, edakit.Num(float64(A)/float64(w), 2), need, need-l),
			map[string]string{"min_len": edakit.IntStr(need)})
	case 2:
		eol := edakit.Pick([]int{40, 50, 60, 70, 80})
		n := edakit.Scaled(ctx, 4, 7)
		x := 0
		var segs []string
		v := 0
		for i := 0; i < n; i++ {
			l := edakit.Int(60, 300)
			segs = append(segs, edakit.Sprintf("[%d, %d]", x, x+l))
			x += l
			if i < n-1 {
				g := edakit.Pick([]int{eol - 20, eol - 10, eol, eol + 10, eol + 30, eol + 60})
				if g < eol {
					v++
				}
				x += g
			}
		}
		q := edakit.Sprintf("Wire segments on one track, in nm along the track, are %s. The end-of-line rule requires facing line ends to be at least %d nm apart. How many end-of-line violations are there?", strings.Join(segs, ", "), eol)
		return edakit.Numeric(q, edakit.IntStr(v),
			edakit.Sprintf("Subtract each segment's end from the next one's start and count gaps below %d nm: %d. The router avoids them by charging a cost when it ends a wire too close to another line end.", eol, v),
			map[string]string{"violations": edakit.IntStr(v)})
	default:
		w := edakit.Pick([]int{20, 24, 32, 40})
		cut := edakit.Pick([]int{20, 24, 32})
		enc := edakit.Pick([]int{10, 15, 20, 30})
		A := 1000 * edakit.Int(2, 12)
		need := maxI(cut+2*enc, ceilDiv(A, w))
		q := edakit.Sprintf("A via lands on a %d nm wide metal stub that connects to nothing else on that layer. The via cut is %d nm long in the wire direction and needs %d nm of metal enclosure beyond each end of the cut. The layer's minimum area is %d $\\text{nm}^2$. What is the shortest legal stub length, in whole nm?", w, cut, enc, A)
		return edakit.Numeric(q, edakit.IntStr(need),
			edakit.Sprintf("Enclosure needs %d + 2(%d) = %d nm; minimum area needs ceil(%d / %d) = %d nm. The stub must satisfy both: %d nm. Rules interact, which is why routers check them against a rule deck rather than one formula.", cut, enc, cut+2*enc, A, w, ceilDiv(A, w), need),
			map[string]string{"stub": edakit.IntStr(need)})
	}
}

// drcRepair: marker cost in an iteration (KP 1), the iteration in which a detour wins (KP 2),
// clips that a repair pass must reroute (KP 3).
func drcRepair(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		b := edakit.Int(1, 3)
		m0 := edakit.Pick([]int{2, 4, 5, 8, 10})
		g := edakit.Pick([]int{2, 3})
		i := edakit.Int(1, 5)
		mi := m0
		for k := 1; k < i; k++ {
			mi *= g
		}
		q := edakit.Sprintf("A search-and-repair detailed router marks every grid cell touched by a DRC violation. Each cell costs %d to cross, plus a marker penalty on marked cells that is %d in repair iteration 1 and is multiplied by %d at each further iteration. How much does crossing one marked cell cost in iteration %d?", b, m0, g, i)
		return edakit.Numeric(q, edakit.IntStr(b+mi),
			edakit.Sprintf("The penalty in iteration %d is %d x %d^%d = %d, so the cell costs %d + %d = %d. Growing penalties push the nets causing the violation away from it a little more each pass.", i, m0, g, i-1, mi, b, mi, b+mi),
			map[string]string{"penalty": edakit.IntStr(mi)})
	case 2:
		m0 := edakit.Pick([]int{1, 2, 3})
		g := 2
		nm := edakit.Int(1, 3)
		d := edakit.Int(5, 40)
		i, pen := 1, m0
		for nm*pen <= d {
			pen *= g
			i++
		}
		q := edakit.Sprintf("In a search-and-repair router a net's current route crosses %s. A clean detour is %d cells longer. The marker penalty per marked cell is %d in iteration 1 and doubles at each further iteration; unmarked cells cost the same on both routes. In which iteration does the detour first become strictly cheaper?", plural(nm, "marked cell"), d, m0)
		return edakit.Numeric(q, edakit.IntStr(i),
			edakit.Sprintf("The detour wins when %d x penalty > %d. The penalty is %d x 2^(i - 1), which first exceeds %d / %d in iteration %d (penalty %d). Until then the router keeps the short, dirty route.", nm, d, m0, d, nm, i, pen),
			map[string]string{"iteration": edakit.IntStr(i)})
	default:
		k := edakit.Pick([]int{5, 7, 10})
		off := edakit.Int(0, k-1)
		W := edakit.Int(30, 60)
		n := edakit.Scaled(ctx, 5, 9)
		clips := map[pt]bool{}
		var ms []pt
		for len(ms) < n {
			m := pt{edakit.Int(0, W-1), edakit.Int(0, W-1)}
			ms = append(ms, m)
			cx := int(math.Floor(float64(m.x+off) / float64(k)))
			cy := int(math.Floor(float64(m.y+off) / float64(k)))
			clips[pt{cx, cy}] = true
		}
		sortPts(ms)
		q := edakit.Sprintf("A repair pass cuts a %d by %d gcell region into %d by %d gcell clips whose boundaries are shifted by %s in x and y this iteration, so gcell (x, y) falls in clip (floor((x + %d) / %d), floor((y + %d) / %d)). Only clips containing a violation marker are ripped up and rerouted. Markers sit in gcells %s. How many clips are rerouted?", W, W, k, k, plural(off, "gcell"), off, k, off, k, pinsText(ms))
		return edakit.Numeric(q, edakit.IntStr(len(clips)),
			edakit.Sprintf("Map each marker to its clip and count distinct clips: %d. Shifting clip boundaries between iterations lets a violation that sat on a boundary be repaired inside one clip next time.", len(clips)),
			map[string]string{"clips": edakit.IntStr(len(clips))})
	}
}

type wseg struct {
	track, l, r int
}

// multiplePatterning: masks needed (KP 1), conflict edges from geometry (KP 2), same-mask
// violations of a proposed colouring (KP 3).
func multiplePatterning(ctx generator.GeneratorContext) generator.Problem {
	g := edakit.Pick([]int{30, 40, 50})
	var ws []wseg
	for t := 1; t <= 3; t++ {
		x := edakit.Int(0, 30)
		for i := 0; i < edakit.Int(2, 3); i++ {
			l := edakit.Int(40, 120)
			ws = append(ws, wseg{t, x, x + l})
			x += l + edakit.Pick([]int{g - 20, g - 10, g + 10, g + 40})
		}
	}
	n := len(ws)
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	edges := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a, b := ws[i], ws[j]
			c := false
			if absI(a.track-b.track) == 1 && minI(a.r, b.r) > maxI(a.l, b.l) {
				c = true
			}
			if a.track == b.track {
				gap := maxI(a.l, b.l) - minI(a.r, b.r)
				if gap < g {
					c = true
				}
			}
			if c {
				adj[i][j], adj[j][i] = true, true
				edges++
			}
		}
	}
	var parts []string
	for i, w := range ws {
		parts = append(parts, edakit.Sprintf("w%d on track %d over [%d, %d]", i+1, w.track, w.l, w.r))
	}
	intro := edakit.Sprintf("A layer is printed by multiple patterning. Its wires are %s (x in nm). Two wires must go on different masks if they are on adjacent tracks with overlapping x ranges (positive overlap length), or on the same track with a tip-to-tip gap below %d nm. ", strings.Join(parts, "; "), g)
	switch kp(ctx) {
	case 1:
		chi := chromatic(adj)
		return edakit.Numeric(intro+"What is the fewest masks that print the layer with no same-mask conflict?", edakit.IntStr(chi),
			edakit.Sprintf("Build the conflict graph (%d edges) and find its chromatic number: %d. Adjacent-track conflicts alone are 2-colourable by track parity; it is tip-to-tip conflicts on one track that close odd cycles. A colouring-aware router avoids creating such cycles.", edges, chi),
			map[string]string{"edges": edakit.IntStr(edges), "masks": edakit.IntStr(chi)})
	case 2:
		return edakit.Numeric(intro+"How many pairs of wires conflict?", edakit.IntStr(edges),
			edakit.Sprintf("Check adjacent-track pairs for x overlap and same-track neighbours for a short tip-to-tip gap: %d conflicting pairs.", edges),
			map[string]string{"edges": edakit.IntStr(edges)})
	default:
		col := make([]int, n)
		var cs []string
		for i := range col {
			col[i] = ws[i].track % 2
			if edakit.Int(0, 3) == 0 {
				col[i] = 1 - col[i]
			}
			cs = append(cs, edakit.Sprintf("w%d %c", i+1, 'A'+col[i]))
		}
		v := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if adj[i][j] && col[i] == col[j] {
					v++
				}
			}
		}
		return edakit.Numeric(intro+"A decomposer proposes two masks: "+strings.Join(cs, ", ")+". How many conflicting pairs end up on the same mask?", edakit.IntStr(v),
			edakit.Sprintf("Of the %d conflicting pairs, %d share a mask. Each is a coloring violation that the router must fix by rerouting or by stitching.", edges, v),
			map[string]string{"edges": edakit.IntStr(edges), "violations": edakit.IntStr(v)})
	}
}

func chromatic(adj [][]bool) int {
	n := len(adj)
	for k := 1; k <= n; k++ {
		col := make([]int, n)
		for i := range col {
			col[i] = -1
		}
		var try func(i int) bool
		try = func(i int) bool {
			if i == n {
				return true
			}
			for c := 0; c < k; c++ {
				ok := true
				for j := 0; j < i; j++ {
					if adj[i][j] && col[j] == c {
						ok = false
						break
					}
				}
				if ok {
					col[i] = c
					if try(i + 1) {
						return true
					}
				}
			}
			col[i] = -1
			return false
		}
		if try(0) {
			return k
		}
	}
	return n
}

var _ = sort.Ints
