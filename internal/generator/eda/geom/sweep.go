package geom

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func randRects(n, span, minS, maxS int) []rect {
	rs := make([]rect, n)
	for i := range rs {
		rs[i] = randRect(span, minS, maxS)
	}
	return rs
}

// planeSweep: distinct event stops (KP 1), maximum status size (KP 2), x-overlapping pairs the
// sweep hands to the y-test (KP 3).
func planeSweep(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 8)
	rs := randRects(n, 30, 3, 12)
	desc := listRects("R", rs)
	switch kp(ctx) {
	case 1:
		var xs []int
		for _, r := range rs {
			xs = append(xs, r.x1, r.x2)
		}
		u := uniq(xs)
		q := edakit.Sprintf("A vertical sweep line moves left to right over the rectangles %s. Each left and right edge is an event; events at the same x are handled at one stop. How many stops does the sweep make?", desc)
		return edakit.Numeric(q, edakit.IntStr(len(u)),
			edakit.Sprintf("The %d events have distinct x-coordinates %s, so the sweep stops %d times.", 2*n, joinInts(u), len(u)),
			map[string]string{"events": edakit.IntStr(2 * n), "stops": edakit.IntStr(len(u))})
	case 2:
		best := 0
		for x2 := 0; x2 <= 60; x2++ {
			c := 0
			for _, r := range rs {
				if 2*r.x1 < x2 && x2 < 2*r.x2 {
					c++
				}
			}
			best = imax(best, c)
		}
		q := edakit.Sprintf("A vertical sweep line moves left to right over the rectangles %s. The status structure holds the rectangles whose x-range the line is strictly inside. What is the largest number of rectangles the status ever holds?", desc)
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("Insert at each left edge and delete at each right edge. The count peaks at %d, the maximum depth of overlapping x-ranges. That peak bounds the status size and the cost of each update.", best),
			map[string]string{"max_status": edakit.IntStr(best)})
	default:
		c := 0
		for i := range rs {
			for j := i + 1; j < len(rs); j++ {
				if overlap1(rs[i].x1, rs[i].x2, rs[j].x1, rs[j].x2) > 0 {
					c++
				}
			}
		}
		q := edakit.Sprintf("A sweep over the rectangles %s compares two rectangles only when both are in the status together, that is, when their x-ranges overlap with positive length. How many pairs does it compare?", desc)
		return edakit.Numeric(q, edakit.IntStr(c),
			edakit.Sprintf("Brute force would test all %d pairs. Only the %d pairs with overlapping x-ranges meet in the status; the others are never compared.", n*(n-1)/2, c),
			map[string]string{"pairs": edakit.IntStr(c), "all_pairs": edakit.IntStr(n * (n - 1) / 2)})
	}
}

type hseg struct{ x1, x2, y int }
type vseg struct{ x, y1, y2 int }

// orthoSegs draws horizontal segments on odd y with even endpoints and vertical segments on
// odd x with even endpoints, so no endpoint ever touches another segment.
func orthoSegs(nh, nv int) ([]hseg, []vseg) {
	var hs []hseg
	var vs []vseg
	ys := edakit.Shuffle([]int{1, 3, 5, 7, 9, 11})
	xs := edakit.Shuffle([]int{1, 3, 5, 7, 9, 11})
	for i := 0; i < nh; i++ {
		a, b := 2*edakit.Int(0, 6), 2*edakit.Int(0, 6)
		for a == b {
			b = 2 * edakit.Int(0, 6)
		}
		hs = append(hs, hseg{imin(a, b), imax(a, b), ys[i%len(ys)]})
	}
	for i := 0; i < nv; i++ {
		a, b := 2*edakit.Int(0, 6), 2*edakit.Int(0, 6)
		for a == b {
			b = 2 * edakit.Int(0, 6)
		}
		vs = append(vs, vseg{xs[i%len(xs)], imin(a, b), imax(a, b)})
	}
	return hs, vs
}

func segText(hs []hseg, vs []vseg) string {
	var parts []string
	for i, h := range hs {
		parts = append(parts, edakit.Sprintf("H%d from (%d, %d) to (%d, %d)", i+1, h.x1, h.y, h.x2, h.y))
	}
	for i, v := range vs {
		parts = append(parts, edakit.Sprintf("V%d from (%d, %d) to (%d, %d)", i+1, v.x, v.y1, v.x, v.y2))
	}
	return strings.Join(parts, "; ")
}

// segmentIntersection: orthogonal crossing count (KP 1), status size at a vertical segment
// (KP 2), brute force versus output-sensitive sweep cost (KP 3).
func segmentIntersection(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		hs, vs := orthoSegs(edakit.Scaled(ctx, 3, 5), edakit.Scaled(ctx, 3, 5))
		c := 0
		for _, h := range hs {
			for _, v := range vs {
				if h.x1 < v.x && v.x < h.x2 && v.y1 < h.y && h.y < v.y2 {
					c++
				}
			}
		}
		q := "Wire segments on one layer are " + segText(hs, vs) + ". How many horizontal-vertical crossings are there?"
		return edakit.Numeric(q, edakit.IntStr(c),
			edakit.Sprintf("A horizontal and a vertical segment cross when the vertical's x lies inside the horizontal's x-range and the horizontal's y lies inside the vertical's y-range. %d pairs qualify.", c),
			map[string]string{"crossings": edakit.IntStr(c)})
	case 2:
		hs, vs := orthoSegs(edakit.Scaled(ctx, 4, 6), 3)
		j := edakit.Int(0, len(vs)-1)
		active, hit := 0, 0
		for _, h := range hs {
			if h.x1 < vs[j].x && vs[j].x < h.x2 {
				active++
				if vs[j].y1 < h.y && h.y < vs[j].y2 {
					hit++
				}
			}
		}
		q := "A left-to-right sweep keeps the active horizontal segments in a balanced tree keyed by y. The segments are " + segText(hs, vs) + edakit.Sprintf(". When the sweep reaches V%d, how many horizontal segments are active?", j+1)
		return edakit.Numeric(q, edakit.IntStr(active),
			edakit.Sprintf("A horizontal segment is active while the sweep x lies inside its x-range: %d of them at x = %d. A range query over [%d, %d] in the tree then reports the %d it crosses.", active, vs[j].x, vs[j].y1, vs[j].y2, hit),
			map[string]string{"active": edakit.IntStr(active), "reported": edakit.IntStr(hit)})
	default:
		p := edakit.Int(10, 20)
		n := 1 << p
		k := edakit.Pick([]int{n / 4, n / 2, n, 2 * n, 4 * n})
		brute := float64(n) * float64(n-1) / 2
		sweep := float64(n+k) * float64(p)
		ratio := brute / sweep
		ans := edakit.Num(ratio, 1)
		q := edakit.Sprintf("A layer has $n = 2^{%d}$ segments with $k = %d$ intersections. Model brute force as $n(n-1)/2$ pair tests and the Bentley-Ottmann sweep as $(n + k)\\log_2 n$ operations. How many times fewer operations does the sweep perform? Round to 1 decimal place.", p, k)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Brute force: %.0f tests. Sweep: (%d + %d)(%d) = %.0f. Ratio %s. The sweep pays log n per event and per reported crossing, so it wins whenever k is far below n^2.", brute, n, k, p, sweep, ans),
			map[string]string{"brute": edakit.Num(brute, 0), "sweep": edakit.Num(sweep, 0)})
	}
}

// unionAreaGen: union area (KP 1), covered length of the status at one x (KP 2), area covered
// at least twice (KP 3).
func unionAreaGen(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		rs := randRects(edakit.Scaled(ctx, 2, 4), 20, 4, 12)
		u := unionArea(rs)
		q := edakit.Sprintf("Metal shapes %s may overlap. What is the area of their union?", listRects("R", rs))
		sum := 0
		for _, r := range rs {
			sum += r.area()
		}
		return edakit.Numeric(q, edakit.IntStr(u),
			edakit.Sprintf("The shapes sum to %d, but overlaps would be double counted. Sweeping in x and multiplying each slab's covered y-length by its width gives %d.", sum, u),
			map[string]string{"sum": edakit.IntStr(sum), "union": edakit.IntStr(u)})
	case 2:
		rs := randRects(edakit.Scaled(ctx, 4, 6), 24, 4, 14)
		x2 := 2*edakit.Int(1, 22) + 1
		var act []rect
		for _, r := range rs {
			if 2*r.x1 < x2 && x2 < 2*r.x2 {
				act = append(act, rect{0, r.y1, 1, r.y2})
			}
		}
		l := 0
		if len(act) > 0 {
			l = unionArea(act)
		}
		q := edakit.Sprintf("A sweep for the union area of %s stops at x = %s. What total y-length does the status cover there (the union of the active y-ranges)?", listRects("R", rs), halfStr(x2))
		return edakit.Numeric(q, edakit.IntStr(l),
			edakit.Sprintf("%d rectangles are active at that x. Merging their y-ranges gives covered length %d, which a segment tree maintains in O(log n) per update. Area is the sum of covered length times slab width.", len(act), l),
			map[string]string{"active": edakit.IntStr(len(act)), "covered": edakit.IntStr(l)})
	default:
		rs := randRects(edakit.Scaled(ctx, 3, 5), 20, 5, 12)
		a2 := coverAtLeast(rs, 2)
		q := edakit.Sprintf("Shapes %s are drawn on one layer. What area is covered by at least two of them (where a double-patterning or overlap check would look)?", listRects("R", rs))
		return edakit.Numeric(q, edakit.IntStr(a2),
			edakit.Sprintf("Keep a cover count per y-interval instead of a covered flag, and sum the length with count at least 2 in each slab: %d.", a2),
			map[string]string{"double": edakit.IntStr(a2)})
	}
}

// overlapPairs: positive-area overlapping pairs (KP 1), pairs closer than s in the square metric
// (KP 2), pairs that only touch (KP 3).
func overlapPairs(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 8)
	switch kp(ctx) {
	case 1:
		rs := randRects(n, 30, 3, 10)
		c := 0
		for i := range rs {
			for j := i + 1; j < n; j++ {
				if interArea(rs[i], rs[j]) > 0 {
					c++
				}
			}
		}
		q := edakit.Sprintf("Report all overlapping pairs among %s. Pairs that merely touch do not count. How many pairs overlap with positive area?", listRects("R", rs))
		return edakit.Numeric(q, edakit.IntStr(c),
			edakit.Sprintf("A pair overlaps when both its x-overlap and its y-overlap are positive. %d pairs qualify; a sweep finds them in O(n log n + k).", c),
			map[string]string{"pairs": edakit.IntStr(c)})
	case 2:
		s := edakit.Int(2, 5)
		rs := randRects(n, 34, 2, 8)
		c := 0
		for i := range rs {
			for j := i + 1; j < n; j++ {
				gx, gy := gaps(rs[i], rs[j])
				d := imax(gx, gy)
				if d > 0 && d < s {
					c++
				}
			}
		}
		q := edakit.Sprintf("Among the shapes %s, how many pairs are separated (distance greater than 0) but closer than %d in the square metric, where distance is the larger of the x-gap and the y-gap?", listRects("R", rs), s)
		return edakit.Numeric(q, edakit.IntStr(c),
			edakit.Sprintf("Growing every shape by %d/2 turns this into an overlap query: a separated pair with max(gx, gy) < %d now overlaps. %d pairs qualify.", s, s, c),
			map[string]string{"pairs": edakit.IntStr(c)})
	default:
		var rs []rect
		for len(rs) < n {
			r := randRect(16, 2, 6)
			if len(rs) > 0 && edakit.Int(0, 1) == 0 {
				o := rs[edakit.Int(0, len(rs)-1)]
				w, h := r.w(), r.h()
				switch edakit.Int(0, 3) {
				case 0:
					r = rect{o.x2, o.y1 + edakit.Int(-3, 3), o.x2 + w, 0}
				case 1:
					r = rect{o.x1 - w, o.y1 + edakit.Int(-3, 3), o.x1, 0}
				case 2:
					r = rect{o.x1 + edakit.Int(-3, 3), o.y2, 0, o.y2 + h}
				default:
					r = rect{o.x2, o.y2, o.x2 + w, o.y2 + h}
				}
				if r.x2 == 0 {
					r.x2 = r.x1 + w
				}
				if r.y2 == 0 {
					r.y2 = r.y1 + h
				}
			}
			rs = append(rs, r)
		}
		c := 0
		for i := range rs {
			for j := i + 1; j < n; j++ {
				if interArea(rs[i], rs[j]) == 0 && closedMeet(rs[i], rs[j]) {
					c++
				}
			}
		}
		q := edakit.Sprintf("Among the closed rectangles %s, how many pairs touch (share boundary points, along an edge or at a corner) without overlapping in area?", listRects("R", rs))
		return edakit.Numeric(q, edakit.IntStr(c),
			edakit.Sprintf("Touching means the overlaps in x and y are both at least 0 but not both positive. %d pairs qualify. Whether touching counts as overlap is a predicate choice the engine must make once and use everywhere.", c),
			map[string]string{"touching": edakit.IntStr(c)})
	}
}

// pointLocation: slab count (KP 1), slab-method storage (KP 2), trapezoidal-map bounds (KP 3).
func pointLocation(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 6)
	type seg struct{ a, b pt }
	var ss []seg
	for i := 0; i < n; i++ {
		x1, x2 := edakit.Int(0, 20), edakit.Int(0, 20)
		for x1 == x2 {
			x2 = edakit.Int(0, 20)
		}
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		ss = append(ss, seg{pt{x1, 3*i + edakit.Int(0, 2)}, pt{x2, 3*i + edakit.Int(0, 2)}})
	}
	var parts []string
	var xs []int
	for i, s := range ss {
		parts = append(parts, edakit.Sprintf("s%d from %s to %s", i+1, s.a, s.b))
		xs = append(xs, s.a.x, s.b.x)
	}
	u := uniq(xs)
	desc := strings.Join(parts, "; ")
	switch kp(ctx) {
	case 1:
		q := "The slab method for point location draws a vertical line through every segment endpoint. For the non-crossing segments " + desc + ", how many slabs result, counting the two unbounded ones?"
		return edakit.Numeric(q, edakit.IntStr(len(u)+1),
			edakit.Sprintf("The endpoints have %d distinct x-coordinates (%s), which cut the plane into %d slabs.", len(u), joinInts(u), len(u)+1),
			map[string]string{"distinct_x": edakit.IntStr(len(u))})
	case 2:
		tot := 0
		for i := 0; i+1 < len(u); i++ {
			for _, s := range ss {
				if s.a.x <= u[i] && u[i+1] <= s.b.x {
					tot++
				}
			}
		}
		q := "The slab method stores, for every bounded slab, the segments that span it, sorted bottom to top. For the non-crossing segments " + desc + ", how many segment entries are stored in total?"
		return edakit.Numeric(q, edakit.IntStr(tot),
			edakit.Sprintf("Between consecutive endpoint x-values %s, count the segments covering each slab and add: %d entries. In the worst case this grows as n^2, which is why trapezoidal maps replace slabs.", joinInts(u), tot),
			map[string]string{"entries": edakit.IntStr(tot)})
	default:
		m := edakit.Scaled(ctx, 5, 400)
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("A trapezoidal map is built over %d non-crossing segments in general position, inside a bounding box. What is the largest number of trapezoids it can have?", m)
			return edakit.Numeric(q, edakit.IntStr(3*m+1),
				edakit.Sprintf("Each segment endpoint shoots two vertical extensions; counting gives at most 3n + 1 = %d trapezoids, so the map has linear size, unlike the slab method.", 3*m+1),
				map[string]string{"trapezoids": edakit.IntStr(3*m + 1)})
		}
		q := edakit.Sprintf("A trapezoidal map is built over %d non-crossing segments in general position, inside a bounding box. What is the largest number of vertices it can have?", m)
		return edakit.Numeric(q, edakit.IntStr(6*m+4),
			edakit.Sprintf("Each of the 2n endpoints is a vertex and adds at most two extension endpoints, plus the four box corners: 6n + 4 = %d.", 6*m+4),
			map[string]string{"vertices": edakit.IntStr(6*m + 4)})
	}
}
