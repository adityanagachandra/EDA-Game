package geom

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// colPoly is an x-monotone rectilinear polygon built from n adjacent columns. Column i spans
// [xs[i], xs[i+1]] in x and [bot[i], top[i]] in y; neighbouring columns overlap in y and differ
// in both their bottom and their top, so every column boundary carries four vertices.
type colPoly struct {
	xs       []int
	bot, top []int
}

// newColPoly draws a column polygon with n columns whose widths lie in [minW, maxW].
func newColPoly(n, minW, maxW int) colPoly {
	for {
		p := colPoly{xs: []int{edakit.Int(0, 6)}}
		for i := 0; i < n; i++ {
			p.xs = append(p.xs, p.xs[i]+edakit.Int(minW, maxW))
			b := edakit.Int(0, 5)
			p.bot = append(p.bot, b)
			p.top = append(p.top, b+edakit.Int(3, 8))
		}
		ok := true
		for i := 0; i+1 < n; i++ {
			if p.bot[i] == p.bot[i+1] || p.top[i] == p.top[i+1] || imax(p.bot[i], p.bot[i+1]) >= imin(p.top[i], p.top[i+1]) {
				ok = false
			}
		}
		if ok {
			return p
		}
	}
}

// verts lists the vertices counter-clockwise: the bottom chain left to right, then the top
// chain right to left.
func (p colPoly) verts() []pt {
	n := len(p.bot)
	var v []pt
	for i := 0; i < n; i++ {
		v = append(v, pt{p.xs[i], p.bot[i]}, pt{p.xs[i+1], p.bot[i]})
	}
	for i := n - 1; i >= 0; i-- {
		v = append(v, pt{p.xs[i+1], p.top[i]}, pt{p.xs[i], p.top[i]})
	}
	return v
}

func (p colPoly) rects() []rect {
	var rs []rect
	for i := range p.bot {
		rs = append(rs, rect{p.xs[i], p.bot[i], p.xs[i+1], p.top[i]})
	}
	return rs
}

func (p colPoly) area() int {
	a := 0
	for _, r := range p.rects() {
		a += r.area()
	}
	return a
}

// ringArea2 is twice the signed shoelace area of a closed ring.
func ringArea2(v []pt) int {
	s := 0
	for i := range v {
		j := (i + 1) % len(v)
		s += v[i].x*v[j].y - v[j].x*v[i].y
	}
	return s
}

func ringPerimeter(v []pt) int {
	s := 0
	for i := range v {
		j := (i + 1) % len(v)
		s += iabs(v[j].x-v[i].x) + iabs(v[j].y-v[i].y)
	}
	return s
}

// reflexCount counts the reflex corners of a counter-clockwise ring.
func reflexCount(v []pt) int {
	r := 0
	for i := range v {
		a, b, c := v[(i+len(v)-1)%len(v)], v[i], v[(i+1)%len(v)]
		if cross(a, b, c) < 0 {
			r++
		}
	}
	return r
}

func reversed(v []pt) []pt {
	out := make([]pt, len(v))
	for i := range v {
		out[i] = v[len(v)-1-i]
	}
	return out
}

func ringText(v []pt) string {
	parts := make([]string, len(v))
	for i, p := range v {
		parts[i] = p.String()
	}
	return strings.Join(parts, ", ")
}

// shoelace: area of a CCW rectilinear polygon (KP 1), the signed sum for a CW listing (KP 2),
// net area of a polygon with a hole (KP 3).
func shoelace(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 2, 4)
	switch kp(ctx) {
	case 1:
		p := newColPoly(n, 1, 5)
		v := p.verts()
		a := ringArea2(v) / 2
		q := edakit.Sprintf("A rectilinear polygon has vertices, in order, %s. Use the shoelace formula to find its area.", ringText(v))
		return edakit.Numeric(q, edakit.IntStr(a),
			edakit.Sprintf("Summing $x_i y_{i+1} - x_{i+1} y_i$ around the ring gives %d, and half of it is the area %d. The listing is counter-clockwise, so the sum is positive.", 2*a, a),
			map[string]string{"twice_area": edakit.IntStr(2 * a)})
	case 2:
		p := newColPoly(n, 1, 5)
		v := reversed(p.verts())
		s := ringArea2(v)
		ans := edakit.Num(float64(s)/2, 1)
		q := edakit.Sprintf("A polygon is stored with vertices, in order, %s. Compute the signed shoelace area $\\frac{1}{2}\\sum_i (x_i y_{i+1} - x_{i+1} y_i)$, keeping its sign.", ringText(v))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The sum of cross terms is %d, so the signed area is %s. The negative sign says the ring is listed clockwise, which layout formats use to mark holes or reversed outlines.", s, ans),
			map[string]string{"twice_signed_area": edakit.IntStr(s)})
	default:
		var p colPoly
		var hole rect
		for {
			p = newColPoly(n, 3, 6)
			var cand []int
			for i := range p.bot {
				if p.xs[i+1]-p.xs[i] >= 3 && p.top[i]-p.bot[i] >= 3 {
					cand = append(cand, i)
				}
			}
			if len(cand) == 0 {
				continue
			}
			i := edakit.Pick(cand)
			x1 := edakit.Int(p.xs[i]+1, p.xs[i+1]-2)
			x2 := edakit.Int(x1+1, p.xs[i+1]-1)
			y1 := edakit.Int(p.bot[i]+1, p.top[i]-2)
			y2 := edakit.Int(y1+1, p.top[i]-1)
			hole = rect{x1, y1, x2, y2}
			break
		}
		hv := []pt{{hole.x1, hole.y1}, {hole.x1, hole.y2}, {hole.x2, hole.y2}, {hole.x2, hole.y1}}
		outer := p.area()
		net := outer - hole.area()
		q := edakit.Sprintf("A layout polygon has an outer ring %s (counter-clockwise) and a hole ring %s (clockwise). Summing the shoelace formula over both rings, what is the net area?", ringText(p.verts()), ringText(hv))
		return edakit.Numeric(q, edakit.IntStr(net),
			edakit.Sprintf("The outer ring gives +%d and the clockwise hole ring gives -%d, so the signed sum is %d.", outer, hole.area(), net),
			map[string]string{"outer": edakit.IntStr(outer), "hole": edakit.IntStr(hole.area())})
	}
}

// crossingsRight counts polygon edges crossed by the ray from (qx, qy) towards +x. Query
// coordinates are doubled so that half-integer queries stay integral.
func crossingsRight(v []pt, qx2, qy2 int) int {
	c := 0
	for i := range v {
		a, b := v[i], v[(i+1)%len(v)]
		if a.x != b.x {
			continue
		}
		lo, hi := imin(a.y, b.y), imax(a.y, b.y)
		if 2*a.x > qx2 && 2*lo < qy2 && qy2 < 2*hi {
			c++
		}
	}
	return c
}

func halfStr(v2 int) string { return edakit.Num(float64(v2)/2, 1) }

// pointInPolygon: crossings of one ray (KP 1), inside count over several points (KP 2),
// even-odd versus nonzero disagreement over overlapping rectangles (KP 3).
func pointInPolygon(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		p := newColPoly(edakit.Scaled(ctx, 2, 4), 1, 5)
		v := p.verts()
		qx2 := 2*edakit.Int(p.xs[0]-2, p.xs[len(p.xs)-1]) + 1
		qy2 := 2*edakit.Int(0, 11) + 1
		c := crossingsRight(v, qx2, qy2)
		q := edakit.Sprintf("A polygon has vertices %s. A ray is cast from the query point (%s, %s) towards +x. How many polygon edges does it cross?", ringText(v), halfStr(qx2), halfStr(qy2))
		return edakit.Numeric(q, edakit.IntStr(c),
			edakit.Sprintf("Only vertical edges to the right of the point whose y-span contains %s can be crossed. There are %d, so the point is %s.", halfStr(qy2), c, map[bool]string{true: "inside (odd)", false: "outside (even)"}[c%2 == 1]),
			map[string]string{"crossings": edakit.IntStr(c)})
	case 2:
		p := newColPoly(edakit.Scaled(ctx, 2, 4), 1, 5)
		v := p.verts()
		var qs []string
		inside := 0
		for i := 0; i < 5; i++ {
			qx2 := 2*edakit.Int(p.xs[0]-1, p.xs[len(p.xs)-1]) + 1
			qy2 := 2*edakit.Int(0, 11) + 1
			if crossingsRight(v, qx2, qy2)%2 == 1 {
				inside++
			}
			qs = append(qs, "("+halfStr(qx2)+", "+halfStr(qy2)+")")
		}
		q := edakit.Sprintf("A polygon has vertices %s. Using the crossing-number test, how many of the points %s lie inside it?", ringText(v), strings.Join(qs, ", "))
		return edakit.Numeric(q, edakit.IntStr(inside),
			edakit.Sprintf("Cast a ray from each point and count edge crossings; an odd count means inside. %d of the five points are inside.", inside),
			map[string]string{"inside": edakit.IntStr(inside)})
	default:
		k := edakit.Int(3, 4)
		var rs []rect
		for len(rs) < k {
			rs = append(rs, randRect(12, 4, 8))
		}
		var qs []string
		diff := 0
		for i := 0; i < 6; i++ {
			qx2, qy2 := 2*edakit.Int(0, 11)+1, 2*edakit.Int(0, 11)+1
			w := 0
			for _, r := range rs {
				if 2*r.x1 < qx2 && qx2 < 2*r.x2 && 2*r.y1 < qy2 && qy2 < 2*r.y2 {
					w++
				}
			}
			if w > 0 && w%2 == 0 {
				diff++
			}
			qs = append(qs, "("+halfStr(qx2)+", "+halfStr(qy2)+")")
		}
		q := edakit.Sprintf("One path is drawn as the overlapping rectangles %s, every one traced counter-clockwise. For the points %s, how many are inside under the nonzero winding rule but outside under the even-odd rule?", listRects("R", rs), strings.Join(qs, ", "))
		return edakit.Numeric(q, edakit.IntStr(diff),
			edakit.Sprintf("With every ring counter-clockwise, a point's winding number is the number of rectangles covering it. Nonzero calls it inside when that count is positive; even-odd only when it is odd. They disagree on points covered an even, positive number of times: %d here.", diff),
			map[string]string{"disagree": edakit.IntStr(diff)})
	}
}

// overlappingPair draws two rectangles that usually overlap.
func overlappingPair(span, minS, maxS int) (rect, rect) {
	a, b := randRect(span, minS, maxS), randRect(span, minS, maxS)
	for i := 0; i < 30 && interArea(a, b) == 0 && edakit.Int(0, 5) > 0; i++ {
		b = randRect(span, minS, maxS)
	}
	return a, b
}

// booleanOps: AND (KP 1), OR (KP 2), XOR or NOT (KP 3) areas of two rectangles.
func booleanOps(ctx generator.GeneratorContext) generator.Problem {
	a, b := overlappingPair(20, 4, 13)
	i := interArea(a, b)
	base := edakit.Sprintf("Layer A holds %s and layer B holds %s. ", a, b)
	facts := map[string]string{"area_a": edakit.IntStr(a.area()), "area_b": edakit.IntStr(b.area()), "and": edakit.IntStr(i)}
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(base+"What is the area of A AND B?", edakit.IntStr(i),
			edakit.Sprintf("AND keeps the region in both: the intersection box, with area %d.", i), facts)
	case 2:
		u := a.area() + b.area() - i
		return edakit.Numeric(base+"What is the area of A OR B?", edakit.IntStr(u),
			edakit.Sprintf("By inclusion-exclusion, OR = %d + %d - %d = %d.", a.area(), b.area(), i, u), facts)
	default:
		if edakit.Int(0, 1) == 0 {
			x := a.area() + b.area() - 2*i
			return edakit.Numeric(base+"What is the area of A XOR B?", edakit.IntStr(x),
				edakit.Sprintf("XOR keeps what is in exactly one layer: %d + %d - 2(%d) = %d.", a.area(), b.area(), i, x), facts)
		}
		x := a.area() - i
		return edakit.Numeric(base+"What is the area of A NOT B (A minus B)?", edakit.IntStr(x),
			edakit.Sprintf("NOT removes the overlap from A: %d - %d = %d.", a.area(), i, x), facts)
	}
}

// sizing: grow or shrink one rectangle (KP 1), grow two and merge (KP 2), closing a gap (KP 3).
func sizing(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		r := rect{0, 0, 2 * edakit.Int(3, 30), 2 * edakit.Int(3, 30)}
		d := edakit.Int(1, 12)
		if edakit.Int(0, 1) == 0 {
			g := r.grow(d)
			q := edakit.Sprintf("A %d nm by %d nm rectangle is grown (sized up) by %d nm on every side. What is the area of the result, in square nm?", r.w(), r.h(), d)
			return edakit.Numeric(q, edakit.IntStr(g.area()),
				edakit.Sprintf("Each dimension grows by 2d: (%d + %d)(%d + %d) = %d.", r.w(), 2*d, r.h(), 2*d, g.area()),
				map[string]string{"w": edakit.IntStr(g.w()), "h": edakit.IntStr(g.h())})
		}
		w, h := imax(0, r.w()-2*d), imax(0, r.h()-2*d)
		q := edakit.Sprintf("A %d nm by %d nm rectangle is shrunk by %d nm on every side. What is the area of the result, in square nm? (A shape that collapses has area 0.)", r.w(), r.h(), d)
		return edakit.Numeric(q, edakit.IntStr(w*h),
			edakit.Sprintf("Each dimension loses 2d = %d nm, giving %d by %d, area %d. A dimension at or below 2d collapses to nothing.", 2*d, w, h, w*h),
			map[string]string{"w": edakit.IntStr(w), "h": edakit.IntStr(h)})
	case 2:
		a := rect{0, edakit.Int(0, 10), edakit.Int(10, 30), 0}
		a.y2 = a.y1 + edakit.Int(10, 30)
		g := edakit.Int(2, 16)
		b := rect{a.x2 + g, edakit.Int(0, 20), 0, 0}
		b.x2 = b.x1 + edakit.Int(10, 30)
		b.y2 = b.y1 + edakit.Int(10, 30)
		d := edakit.Int(1, 8)
		ga, gb := a.grow(d), b.grow(d)
		u := ga.area() + gb.area() - interArea(ga, gb)
		q := edakit.Sprintf("Shapes A = %s and B = %s (nm) are each grown by %d nm and the results are merged with OR. What is the area of the merged result, in square nm?", a, b, d)
		return edakit.Numeric(q, edakit.IntStr(u),
			edakit.Sprintf("Grown, A is %s and B is %s. Their overlap is %d, so the union is %d + %d - %d = %d. They merge only when 2d is at least the %d nm gap.", ga, gb, interArea(ga, gb), ga.area(), gb.area(), interArea(ga, gb), u, g),
			map[string]string{"overlap": edakit.IntStr(interArea(ga, gb))})
	default:
		h := edakit.Int(10, 40)
		w1, w2 := edakit.Int(10, 40), edakit.Int(10, 40)
		d := edakit.Int(2, 10)
		g := edakit.Int(1, 4*d)
		if g == 2*d {
			g++
		}
		var ar int
		if g < 2*d {
			ar = (w1 + g + w2) * h
		} else {
			ar = (w1 + w2) * h
		}
		a := rect{0, 0, w1, h}
		b := rect{w1 + g, 0, w1 + g + w2, h}
		q := edakit.Sprintf("Shapes A = %s and B = %s (nm) are grown by %d nm, merged, and the merged result is shrunk by %d nm (a closing). What is the final area, in square nm?", a, b, d, d)
		return edakit.Numeric(q, edakit.IntStr(ar),
			edakit.Sprintf("The %d nm gap is %s 2d = %d nm, so the grown shapes %s; shrinking then returns %s. Area %d.", g, map[bool]string{true: "below", false: "above"}[g < 2*d], 2*d,
				map[bool]string{true: "fuse and the bridge survives the shrink", false: "stay apart"}[g < 2*d],
				map[bool]string{true: "one bar across the gap", false: "the two original shapes"}[g < 2*d], ar),
			map[string]string{"gap": edakit.IntStr(g), "two_d": edakit.IntStr(2 * d)})
	}
}

// offsetArea: the A + P d + 4 d^2 law from given A and P (KP 1), from a vertex list (KP 2),
// and the round-corner Euclidean offset (KP 3).
func offsetArea(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		for {
			p := newColPoly(edakit.Int(2, 4), 6, 12)
			a, per := p.area(), ringPerimeter(p.verts())
			d := edakit.Int(1, 2)
			if edakit.Int(0, 1) == 0 {
				ans := a + per*d + 4*d*d
				q := edakit.Sprintf("A rectilinear polygon has area %d and perimeter %d (in um and square um). It is grown by %d um with square corners, and no notch is narrower than %d um. What is the grown area, in square um?", a, per, d, 2*d)
				return edakit.Numeric(q, edakit.IntStr(ans),
					edakit.Sprintf("Each edge sweeps a strip of length times d, convex corners add a d-by-d square and reflex corners remove one; convex minus reflex is always 4. So %d + %d(%d) + 4(%d) = %d.", a, per, d, d*d, ans),
					map[string]string{"area": edakit.IntStr(a), "perimeter": edakit.IntStr(per)})
			}
			ans := a - per*d + 4*d*d
			if erodedArea(p.rects(), d) != ans {
				continue // a neck or arm too thin for the law; draw again
			}
			q := edakit.Sprintf("A rectilinear polygon has area %d and perimeter %d (in um and square um). It is shrunk by %d um, and every part of it is wider than %d um so nothing collapses. What is the shrunk area, in square um?", a, per, d, 2*d)
			return edakit.Numeric(q, edakit.IntStr(ans),
				edakit.Sprintf("Shrinking removes a strip of length times d along each edge; strips overlap at convex corners and leave a d-by-d square at reflex corners, a net +4 d^2. So %d - %d(%d) + 4(%d) = %d.", a, per, d, d*d, ans),
				map[string]string{"area": edakit.IntStr(a), "perimeter": edakit.IntStr(per)})
		}
	case 2:
		d := edakit.Int(1, 2)
		p := newColPoly(edakit.Scaled(ctx, 2, 4), 2*d+1, 2*d+4)
		var grown []rect
		for _, r := range p.rects() {
			grown = append(grown, r.grow(d))
		}
		ans := unionArea(grown)
		a, per := p.area(), ringPerimeter(p.verts())
		q := edakit.Sprintf("A polygon with vertices %s is grown by %d with square corners. No pocket is narrower than %d. What is the grown area?", ringText(p.verts()), d, 2*d)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Shoelace gives area %d and the edge lengths sum to perimeter %d, so the grown area is %d + %d(%d) + 4(%d) = %d.", a, per, a, per, d, d*d, ans),
			map[string]string{"area": edakit.IntStr(a), "perimeter": edakit.IntStr(per)})
	default:
		d := edakit.Int(1, 3)
		var p colPoly
		for {
			p = newColPoly(edakit.Scaled(ctx, 2, 4), 2*d+1, 2*d+5)
			if minEdge(p.verts()) >= d {
				break
			}
		}
		v := p.verts()
		a, per, r := p.area(), ringPerimeter(v), reflexCount(v)
		c := r + 4
		val := float64(a) + float64(per*d) + (float64(c)*math.Pi/4-float64(r))*float64(d*d)
		ans := edakit.Num(val, 2)
		q := edakit.Sprintf("A polygon with vertices %s is grown by %d with round (Euclidean) corners: every point within distance %d of the polygon is added. No pocket is narrower than %d. What is the grown area? Round to 2 decimal places.", ringText(v), d, d, 2*d)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Area %d, perimeter %d, %d reflex and %d convex corners. Edge strips add %d; each convex corner adds a quarter disc of area pi d^2 / 4 and each reflex corner removes a d-by-d overlap: %d + %d + (%d pi/4 - %d)(%d) = %s.", a, per, r, c, per*d, a, per*d, c, r, d*d, ans),
			map[string]string{"area": edakit.IntStr(a), "perimeter": edakit.IntStr(per), "reflex": edakit.IntStr(r)})
	}
}

// minkowski: sum of two rectangles (KP 1), edge count of a convex sum (KP 2), the via keep-out
// region as a Minkowski sum (KP 3).
func minkowski(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := randRect(20, 2, 10)
		b := rect{edakit.Int(-5, 5), edakit.Int(-5, 5), 0, 0}
		b.x2, b.y2 = b.x1+edakit.Int(1, 8), b.y1+edakit.Int(1, 8)
		s := rect{a.x1 + b.x1, a.y1 + b.y1, a.x2 + b.x2, a.y2 + b.y2}
		q := edakit.Sprintf("What is the area of the Minkowski sum of A = %s and B = %s?", a, b)
		return edakit.Numeric(q, edakit.IntStr(s.area()),
			edakit.Sprintf("The sum of two boxes is the box whose ranges are the sums of the ranges: %s, area %d x %d = %d.", s, s.w(), s.h(), s.area()),
			map[string]string{"w": edakit.IntStr(s.w()), "h": edakit.IntStr(s.h())})
	case 2:
		dirs := []int{0, 45, 90, 135, 180, 225, 270, 315}
		pick := func() []int {
			for {
				var s []int
				for _, d := range dirs {
					if edakit.Int(0, 1) == 0 {
						s = append(s, d)
					}
				}
				if len(s) < 3 {
					continue
				}
				ok := true
				for i := range s {
					next := s[(i+1)%len(s)]
					gap := (next - s[i] + 360) % 360
					if gap == 0 {
						gap = 360
					}
					if gap >= 180 {
						ok = false
					}
				}
				if ok {
					return s
				}
			}
		}
		p, qd := pick(), pick()
		set := map[int]bool{}
		for _, d := range p {
			set[d] = true
		}
		for _, d := range qd {
			set[d] = true
		}
		q := edakit.Sprintf("Convex polygon P has edges with outward normals at %s degrees, and convex polygon Q has edges with outward normals at %s degrees. How many edges does the Minkowski sum P + Q have?", joinInts(p), joinInts(qd))
		return edakit.Numeric(q, edakit.IntStr(len(set)),
			edakit.Sprintf("The boundary of a convex Minkowski sum merges the two edge sequences by angle; edges with the same normal fuse into one. %d + %d edges minus %d shared directions = %d.", len(p), len(qd), len(p)+len(qd)-len(set), len(set)),
			map[string]string{"edges": edakit.IntStr(len(set))})
	default:
		w, h := 10*edakit.Int(5, 30), 10*edakit.Int(5, 30)
		a, b := 10*edakit.Int(2, 8), 10*edakit.Int(2, 8)
		s := 10 * edakit.Int(3, 9)
		W, H := w+a+2*s, h+b+2*s
		q := edakit.Sprintf("An obstacle is a %d nm by %d nm rectangle. A via of %d nm by %d nm must keep at least %d nm of square-metric spacing from it. The set of via-centre positions that violate the rule is a Minkowski sum. What is its area, in square nm?", w, h, a, b, s)
		return edakit.Numeric(q, edakit.IntStr(W*H),
			edakit.Sprintf("The forbidden region is the obstacle summed with the via shape grown by the spacing: (%d + %d + 2(%d)) by (%d + %d + 2(%d)) = %d by %d, area %d.", w, a, s, h, b, s, W, H, W*H),
			map[string]string{"w": edakit.IntStr(W), "h": edakit.IntStr(H)})
	}
}

// minEdge is the length of the shortest edge of a ring.
func minEdge(v []pt) int {
	m := 1 << 30
	for i := range v {
		j := (i + 1) % len(v)
		m = imin(m, iabs(v[j].x-v[i].x)+iabs(v[j].y-v[i].y))
	}
	return m
}

// erodedArea rasterises a union of integer rectangles on the unit grid and returns the area
// of its erosion by the square of half-side d: the unit cells whose d-neighbourhood is filled.
func erodedArea(rs []rect, d int) int {
	minx, miny, maxx, maxy := rs[0].x1, rs[0].y1, rs[0].x2, rs[0].y2
	for _, r := range rs {
		minx, miny = imin(minx, r.x1), imin(miny, r.y1)
		maxx, maxy = imax(maxx, r.x2), imax(maxy, r.y2)
	}
	filled := func(x, y int) bool {
		for _, r := range rs {
			if r.x1 <= x && x < r.x2 && r.y1 <= y && y < r.y2 {
				return true
			}
		}
		return false
	}
	n := 0
	for x := minx; x < maxx; x++ {
		for y := miny; y < maxy; y++ {
			ok := true
			for dx := -d; dx <= d && ok; dx++ {
				for dy := -d; dy <= d && ok; dy++ {
					if !filled(x+dx, y+dy) {
						ok = false
					}
				}
			}
			if ok {
				n++
			}
		}
	}
	return n
}
