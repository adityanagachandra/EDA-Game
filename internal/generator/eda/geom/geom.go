// Package geom generates computational-geometry problems for layout tools: rectilinear
// distances and boxes, polygon area and Boolean operations, sizing, plane sweeps, spatial
// indexes, convex hulls and Delaunay counts, rectangle decomposition and the geometric checks
// behind DRC, extraction and fill. Everything is integer and rectilinear, as in a layout
// database, so almost every concept computes.
package geom

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the computational-geometry generators to reg.
func Register(reg *generator.Registry) {
	// rect
	reg.Register("geom.rect.manhattan", edakit.Func(manhattan))
	reg.Register("geom.rect.bbox", edakit.Func(bbox))
	reg.Register("geom.rect.orientation", edakit.Func(orientation))
	reg.Register("geom.rect.exact_arithmetic", &edakit.Bank{Items: exactArithmeticItems})
	reg.Register("geom.rect.spacing", edakit.Func(spacing))
	// poly
	reg.Register("geom.poly.shoelace", edakit.Func(shoelace))
	reg.Register("geom.poly.point_in_polygon", edakit.Func(pointInPolygon))
	reg.Register("geom.poly.boolean_ops", edakit.Func(booleanOps))
	reg.Register("geom.poly.sizing", edakit.Func(sizing))
	reg.Register("geom.poly.offset_area", edakit.Func(offsetArea))
	reg.Register("geom.poly.minkowski", edakit.Func(minkowski))
	reg.Register("geom.poly.representation", &edakit.Bank{Items: representationItems})
	// sweep
	reg.Register("geom.sweep.plane_sweep", edakit.Func(planeSweep))
	reg.Register("geom.sweep.segment_intersection", edakit.Func(segmentIntersection))
	reg.Register("geom.sweep.union_area", edakit.Func(unionAreaGen))
	reg.Register("geom.sweep.overlap_pairs", edakit.Func(overlapPairs))
	reg.Register("geom.sweep.point_location", edakit.Func(pointLocation))
	// index
	reg.Register("geom.index.interval_tree", edakit.Func(intervalTree))
	reg.Register("geom.index.segment_tree", edakit.Func(segmentTree))
	reg.Register("geom.index.rtree", edakit.Func(rtree))
	reg.Register("geom.index.rtree_split", edakit.Func(rtreeSplit))
	reg.Register("geom.index.kd_tree", edakit.Func(kdTree))
	reg.Register("geom.index.quadtree", edakit.Func(quadtree))
	reg.Register("geom.index.uniform_grid", edakit.Func(uniformGrid))
	reg.Register("geom.index.selection", &edakit.Bank{Items: selectionItems})
	// prox
	reg.Register("geom.prox.convex_hull", edakit.Func(convexHullGen))
	reg.Register("geom.prox.voronoi", edakit.Func(voronoi))
	reg.Register("geom.prox.delaunay", edakit.Func(delaunay))
	// decomp
	reg.Register("geom.decomp.rectangles", edakit.Func(decompRectangles))
	reg.Register("geom.decomp.slab", edakit.Func(decompSlab))
	// layout
	reg.Register("geom.layout.spacing_check", edakit.Func(spacingCheck))
	reg.Register("geom.layout.width_check", edakit.Func(widthCheck))
	reg.Register("geom.layout.density", edakit.Func(density))
	reg.Register("geom.layout.hierarchy", edakit.Func(hierarchy))
	reg.Register("geom.layout.parallel_run", edakit.Func(parallelRun))
	reg.Register("geom.layout.connectivity", edakit.Func(connectivity))
	// frontier
	reg.Register("geom.frontier.kernel_core", &edakit.Bank{Items: kernelCoreItems})
	reg.Register("geom.frontier.change_tiers", &edakit.Bank{Items: changeTierItems})
	reg.Register("geom.frontier.spatial_query_tool", &edakit.Bank{Items: spatialQueryItems})
	reg.Register("geom.frontier.parallel_tiling", edakit.Func(parallelTiling))
	reg.Register("geom.frontier.differential_check", edakit.Func(differentialCheck))
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
func (r rect) per() int  { return 2 * (r.w() + r.h()) }

func (r rect) String() string {
	return edakit.Sprintf("$[%d, %d] \\times [%d, %d]$", r.x1, r.x2, r.y1, r.y2)
}

func (r rect) grow(d int) rect { return rect{r.x1 - d, r.y1 - d, r.x2 + d, r.y2 + d} }

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

// overlap1 is the length of the overlap of [a1, a2] and [b1, b2], or the negative gap.
func overlap1(a1, a2, b1, b2 int) int { return imin(a2, b2) - imax(a1, b1) }

// interArea is the area of the intersection of a and b (0 when they only touch or are apart).
func interArea(a, b rect) int {
	ox, oy := overlap1(a.x1, a.x2, b.x1, b.x2), overlap1(a.y1, a.y2, b.y1, b.y2)
	if ox <= 0 || oy <= 0 {
		return 0
	}
	return ox * oy
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

// edgeConnected reports whether a and b overlap or share a boundary segment of positive length.
func edgeConnected(a, b rect) bool {
	ox, oy := overlap1(a.x1, a.x2, b.x1, b.x2), overlap1(a.y1, a.y2, b.y1, b.y2)
	if ox < 0 || oy < 0 {
		return false
	}
	return !(ox == 0 && oy == 0) && (ox > 0 || oy > 0)
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

// coverage compresses the coordinates of the rectangles in sets and reports, for every cell of
// the compressed grid, how many rectangles of each set cover it. f receives the cell area and
// the per-set cover counts.
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
func unionArea(rs []rect) int { return coverAtLeast(rs, 1) }

// coverAtLeast is the area covered by at least k rectangles.
func coverAtLeast(rs []rect, k int) int {
	total := 0
	coverage([][]rect{rs}, func(a int, c []int) {
		if c[0] >= k {
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

type pt struct{ x, y int }

func (p pt) String() string { return edakit.Sprintf("(%d, %d)", p.x, p.y) }

func listPts(ps []pt) string {
	parts := make([]string, len(ps))
	for i, p := range ps {
		parts[i] = p.String()
	}
	return strings.Join(parts, ", ")
}

// cross is (b - a) x (c - a).
func cross(a, b, c pt) int { return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x) }

// distinctPts draws n distinct lattice points in [0, span]^2.
func distinctPts(n, span int) []pt {
	seen := map[pt]bool{}
	var out []pt
	for len(out) < n {
		p := pt{edakit.Int(0, span), edakit.Int(0, span)}
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}

func joinInts(xs []int) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.IntStr(x)
	}
	return strings.Join(parts, ", ")
}
