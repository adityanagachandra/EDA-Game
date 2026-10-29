// Package pv generates physical-verification and design-for-manufacturing problems: layout
// stream records and hierarchy, design-rule checks on drawn rectangles, the sweeps and
// hierarchy behind a DRC engine, device recognition and connectivity extraction, netlist
// comparison by partition refinement, antenna ratios and their repairs, density windows and
// fill, double-patterning conflict graphs, critical area, mask data preparation and the
// agentic loops that read and fix verification reports. Layout coordinates are integers, as in
// a layout database, so most concepts compute from a drawn instance.
package pv

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the physical-verification generators to reg.
func Register(reg *generator.Registry) {
	// format
	reg.Register("pv.format.gds_records", edakit.Func(gdsRecords))
	reg.Register("pv.format.gds_oasis", &edakit.Mix{Parts: []generator.Generator{edakit.Func(oasisVarint), edakit.Func(gdsFlatBytes), &edakit.Bank{Items: oasisItems}}})
	reg.Register("pv.format.units_grid", edakit.Func(unitsGrid))
	reg.Register("pv.format.hierarchy", edakit.Func(hierarchy))
	// drc
	reg.Register("pv.drc.rule_types", &edakit.Bank{Items: ruleTypeItems})
	reg.Register("pv.drc.width_spacing", edakit.Func(widthSpacing))
	reg.Register("pv.drc.enclosure_extension", edakit.Func(enclosureExtension))
	reg.Register("pv.drc.notch_eol", edakit.Func(notchEOL))
	reg.Register("pv.drc.area_via", edakit.Func(areaVia))
	reg.Register("pv.drc.engine", &edakit.Mix{Parts: []generator.Generator{edakit.Func(engineCandidates), edakit.Func(engineSizedMarker), &edakit.Bank{Items: engineItems}}})
	reg.Register("pv.drc.plane_sweep", edakit.Func(planeSweep))
	reg.Register("pv.drc.hierarchical", &edakit.Mix{Parts: []generator.Generator{edakit.Func(hierErrorCount), edakit.Func(hierInteractions), &edakit.Bank{Items: hierItems}}})
	reg.Register("pv.drc.rule_deck", &edakit.Mix{Parts: []generator.Generator{edakit.Func(deckBoolean), edakit.Func(deckSized), &edakit.Bank{Items: deckItems}}})
	// lvs
	reg.Register("pv.lvs.device_recognition", edakit.Func(deviceRecognition))
	reg.Register("pv.lvs.extraction", edakit.Func(extraction))
	reg.Register("pv.lvs.reduction", edakit.Func(reduction))
	reg.Register("pv.lvs.graph_matching", &edakit.Mix{Parts: []generator.Generator{edakit.Func(bipartiteSize), &edakit.Bank{Items: matchingItems}, &edakit.Bank{Items: matchingItems}}})
	reg.Register("pv.lvs.partition_refinement", edakit.Func(partitionRefinement))
	reg.Register("pv.lvs.ambiguity", &edakit.Mix{Parts: []generator.Generator{edakit.Func(ambiguousClasses), &edakit.Bank{Items: ambiguityItems}, &edakit.Bank{Items: ambiguityItems}}})
	reg.Register("pv.lvs.errors", &edakit.Bank{Items: lvsErrorItems})
	reg.Register("pv.lvs.erc", &edakit.Mix{Parts: []generator.Generator{edakit.Func(ercFloating), edakit.Func(ercWellTaps), &edakit.Bank{Items: ercItems}}})
	// antenna
	reg.Register("pv.antenna.rules", edakit.Func(antennaRules))
	reg.Register("pv.antenna.cumulative", edakit.Func(antennaCumulative))
	reg.Register("pv.antenna.diode", edakit.Func(antennaDiode))
	reg.Register("pv.antenna.jumper", edakit.Func(antennaJumper))
	// fill
	reg.Register("pv.fill.density", edakit.Func(fillDensity))
	reg.Register("pv.fill.insertion", edakit.Func(fillInsertion))
	reg.Register("pv.fill.windows", edakit.Func(fillWindows))
	reg.Register("pv.fill.side_effects", &edakit.Mix{Parts: []generator.Generator{edakit.Func(fillCoupling), &edakit.Bank{Items: fillItems}, &edakit.Bank{Items: fillItems}}})
	// dfm
	reg.Register("pv.dfm.hotspots", &edakit.Mix{Parts: []generator.Generator{edakit.Func(hotspotMetrics), &edakit.Bank{Items: hotspotItems}, &edakit.Bank{Items: hotspotItems}}})
	reg.Register("pv.dfm.pattern_matching", edakit.Func(patternMatching))
	reg.Register("pv.dfm.dp_coloring", edakit.Func(dpColoring))
	reg.Register("pv.dfm.critical_area", edakit.Func(criticalArea))
	reg.Register("pv.dfm.recommended_rules", &edakit.Mix{Parts: []generator.Generator{edakit.Func(dfmScore), edakit.Func(wireSpreading), &edakit.Bank{Items: recommendedItems}}})
	// tapeout
	reg.Register("pv.tapeout.handoff", &edakit.Bank{Items: handoffItems})
	reg.Register("pv.tapeout.mask_prep", edakit.Func(maskPrep))
	reg.Register("pv.tapeout.xor_check", edakit.Func(xorCheck))
	reg.Register("pv.tapeout.open_tools", &edakit.Bank{Items: openToolItems})
	// frontier
	reg.Register("pv.frontier.agentic_fix_loop", &edakit.Mix{Parts: []generator.Generator{edakit.Func(fixRouting), &edakit.Bank{Items: fixLoopItems}, &edakit.Bank{Items: fixLoopItems}}})
	reg.Register("pv.frontier.change_tiers", &edakit.Bank{Items: changeTierItems})
	reg.Register("pv.frontier.golden_clean", edakit.Func(goldenClean))
	reg.Register("pv.frontier.fast_hotspot", edakit.Func(fastHotspot))
	reg.Register("pv.frontier.ip_boundary", &edakit.Bank{Items: ipBoundaryItems})
}

// kp resolves the knowledge point to serve: the requested one, or a uniform draw when the
// request is unguided.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// rect is a closed axis-parallel rectangle [x1, x2] x [y1, y2] with x1 < x2 and y1 < y2.
type rect struct{ x1, y1, x2, y2 int }

func (r rect) w() int    { return r.x2 - r.x1 }
func (r rect) h() int    { return r.y2 - r.y1 }
func (r rect) area() int { return r.w() * r.h() }

func (r rect) String() string {
	return edakit.Sprintf("$[%d, %d] \\times [%d, %d]$", r.x1, r.x2, r.y1, r.y2)
}

func (r rect) grow(d int) rect { return rect{r.x1 - d, r.y1 - d, r.x2 + d, r.y2 + d} }

func (r rect) shift(dx, dy int) rect { return rect{r.x1 + dx, r.y1 + dy, r.x2 + dx, r.y2 + dy} }

// randRect draws a rectangle with corners in [0, span] and sides in [minS, maxS].
func randRect(span, minS, maxS int) rect {
	w, h := edakit.Int(minS, maxS), edakit.Int(minS, maxS)
	if w > span {
		w = span
	}
	if h > span {
		h = span
	}
	x, y := edakit.Int(0, span-w), edakit.Int(0, span-h)
	return rect{x, y, x + w, y + h}
}

// overlap1 is the overlap length of [a1, a2] and [b1, b2], negative for a gap.
func overlap1(a1, a2, b1, b2 int) int { return imin(a2, b2) - imax(a1, b1) }

// interArea is the area of the intersection of a and b.
func interArea(a, b rect) int {
	ox, oy := overlap1(a.x1, a.x2, b.x1, b.x2), overlap1(a.y1, a.y2, b.y1, b.y2)
	if ox <= 0 || oy <= 0 {
		return 0
	}
	return ox * oy
}

// intersect returns the intersection of a and b and whether it has positive area.
func intersect(a, b rect) (rect, bool) {
	r := rect{imax(a.x1, b.x1), imax(a.y1, b.y1), imin(a.x2, b.x2), imin(a.y2, b.y2)}
	return r, r.x2 > r.x1 && r.y2 > r.y1
}

// gaps returns the horizontal and vertical gaps between a and b (0 when projections overlap).
func gaps(a, b rect) (int, int) {
	gx := imax(0, imax(b.x1-a.x2, a.x1-b.x2))
	gy := imax(0, imax(b.y1-a.y2, a.y1-b.y2))
	return gx, gy
}

// closedMeet reports whether the closed rectangles share at least one point.
func closedMeet(a, b rect) bool {
	return overlap1(a.x1, a.x2, b.x1, b.x2) >= 0 && overlap1(a.y1, a.y2, b.y1, b.y2) >= 0
}

// dist2 is the squared Euclidean distance between two disjoint rectangles.
func dist2(a, b rect) int {
	gx, gy := gaps(a, b)
	return gx*gx + gy*gy
}

// separated draws n rectangles in [0, span]^2 that pairwise do not touch.
func separated(n, span, minS, maxS int) []rect {
	for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
		var rs []rect
		for tries := 0; len(rs) < n && tries < 400; tries++ {
			r := randRect(span, minS, maxS)
			ok := true
			for _, s := range rs {
				if closedMeet(r, s) {
					ok = false
					break
				}
			}
			if ok {
				rs = append(rs, r)
			}
		}
		if len(rs) == n {
			return rs
		}
	}
	// Known-good fallback: a loose grid of minimum-size shapes.
	var rs []rect
	for i := 0; i < n; i++ {
		x, y := (i%4)*(minS+span/4), (i/4)*(minS+span/4)
		rs = append(rs, rect{x, y, x + minS, y + minS})
	}
	return rs
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func iabs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func ceilDiv(a, b int) int {
	if a <= 0 {
		return 0
	}
	return (a + b - 1) / b
}

// coverage compresses the coordinates of the rectangle sets and calls f for every cell of the
// compressed grid with the cell area and how many rectangles of each set cover it.
func coverage(sets [][]rect, f func(cellArea int, counts []int)) {
	var xs, ys []int
	for _, s := range sets {
		for _, r := range s {
			xs = append(xs, r.x1, r.x2)
			ys = append(ys, r.y1, r.y2)
		}
	}
	xs, ys = uniq(xs), uniq(ys)
	counts := make([]int, len(sets))
	for i := 0; i+1 < len(xs); i++ {
		for j := 0; j+1 < len(ys); j++ {
			for k, s := range sets {
				counts[k] = 0
				for _, r := range s {
					if r.x1 <= xs[i] && xs[i+1] <= r.x2 && r.y1 <= ys[j] && ys[j+1] <= r.y2 {
						counts[k]++
					}
				}
			}
			f((xs[i+1]-xs[i])*(ys[j+1]-ys[j]), counts)
		}
	}
}

// unionArea is the area covered by at least one rectangle.
func unionArea(rs []rect) int {
	total := 0
	coverage([][]rect{rs}, func(a int, c []int) {
		if c[0] > 0 {
			total += a
		}
	})
	return total
}

// xorArea is the area covered by exactly one of the two unions.
func xorArea(a, b []rect) int {
	total := 0
	coverage([][]rect{a, b}, func(ar int, c []int) {
		if (c[0] > 0) != (c[1] > 0) {
			total += ar
		}
	})
	return total
}

func uniq(xs []int) []int {
	sort.Ints(xs)
	out := xs[:0]
	for i, x := range xs {
		if i == 0 || x != xs[i-1] {
			out = append(out, x)
		}
	}
	return out
}

// listRects names rectangles R1, R2, ... in one sentence fragment.
func listRects(prefix string, rs []rect) string {
	parts := make([]string, len(rs))
	for i, r := range rs {
		parts[i] = edakit.Sprintf("%s%d = %s", prefix, i+1, r.String())
	}
	return strings.Join(parts, ", ")
}

func joinInts(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.IntStr(x)
	}
	return strings.Join(parts, ", ")
}

// dsu is a disjoint-set forest with path halving.
type dsu struct{ p []int }

func newDSU(n int) *dsu {
	d := &dsu{p: make([]int, n)}
	for i := range d.p {
		d.p[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	for d.p[x] != x {
		d.p[x] = d.p[d.p[x]]
		x = d.p[x]
	}
	return x
}

func (d *dsu) union(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	d.p[ra] = rb
	return true
}

func (d *dsu) count() int {
	c := 0
	for i := range d.p {
		if d.find(i) == i {
			c++
		}
	}
	return c
}

// numChoice builds a computed multiple-choice problem whose options are numbers with a unit
// suffix. The distractors are the given wrong values, deduplicated and padded with nearby values.
func numChoice(q string, correct float64, places int, unit, why string, wrong ...float64) generator.Problem {
	fmtv := func(x float64) string {
		s := edakit.Num(x, places)
		if unit != "" {
			s += " " + unit
		}
		return s
	}
	ans := fmtv(correct)
	seen := map[string]bool{ans: true}
	var ds []string
	for _, w := range wrong {
		s := fmtv(w)
		if !seen[s] && len(ds) < 4 {
			seen[s] = true
			ds = append(ds, s)
		}
	}
	step := 1.0
	for i := 0; i < places; i++ {
		step /= 10
	}
	if correct >= 20 {
		step = float64(imax(1, int(correct)/10))
	}
	for k := 1; len(ds) < 3; k++ {
		for _, c := range []float64{correct + float64(k)*step, correct - float64(k)*step} {
			s := fmtv(c)
			if c >= 0 && !seen[s] && len(ds) < 3 {
				seen[s] = true
				ds = append(ds, s)
			}
		}
	}
	return edakit.Choice(q, ans, why, ds...)
}
