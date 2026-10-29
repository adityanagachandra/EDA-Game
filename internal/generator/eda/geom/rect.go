package geom

import (
	"math"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func binom(n, k int) int {
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

// manhattan: L1 distance (KP 1), L1 minus L2 (KP 2), number of shortest grid routes (KP 3).
func manhattan(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		span := edakit.Scaled(ctx, 40, 400)
		a, b := pt{edakit.Int(0, span), edakit.Int(0, span)}, pt{edakit.Int(0, span), edakit.Int(0, span)}
		d := iabs(a.x-b.x) + iabs(a.y-b.y)
		q := edakit.Sprintf("A driver pin sits at %s and a sink pin at %s, in nm. What is the Manhattan (L1) distance between them, in nm?", a, b)
		return edakit.Numeric(q, edakit.IntStr(d),
			edakit.Sprintf("Add the absolute coordinate differences: |%d - %d| + |%d - %d| = %d + %d = %d nm.", a.x, b.x, a.y, b.y, iabs(a.x-b.x), iabs(a.y-b.y), d),
			map[string]string{"l1": edakit.IntStr(d)})
	case 2:
		dx, dy := edakit.Int(3, 60), edakit.Int(3, 60)
		a := pt{edakit.Int(0, 50), edakit.Int(0, 50)}
		b := pt{a.x + dx, a.y + dy}
		l1 := float64(dx + dy)
		l2 := math.Hypot(float64(dx), float64(dy))
		ans := edakit.Num(l1-l2, 2)
		q := edakit.Sprintf("Two pins sit at %s and %s, in um. A rectilinear wire must run at least the Manhattan distance, while the straight line is the Euclidean distance. By how many um does the Manhattan distance exceed the Euclidean distance? Round to 2 decimal places.", a, b)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("L1 = %d + %d = %d um. L2 = sqrt(%d^2 + %d^2) = %.4f um. The difference is %s um.", dx, dy, dx+dy, dx, dy, l2, ans),
			map[string]string{"l1": edakit.IntStr(dx + dy), "l2": edakit.Num(l2, 4)})
	default:
		dx, dy := edakit.Scaled(ctx, 1, 7), edakit.Int(1, 5)
		a := pt{edakit.Int(0, 9), edakit.Int(0, 9)}
		b := pt{a.x + dx, a.y - dy}
		if edakit.Int(0, 1) == 0 {
			b.y = a.y + dy
		}
		n := binom(dx+dy, dx)
		q := edakit.Sprintf("On an empty unit routing grid, a two-pin net connects grid points %s and %s. How many distinct shortest (Manhattan-length) grid paths connect them?", a, b)
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("Every shortest path makes %d horizontal and %d vertical unit steps in some order, so there are C(%d, %d) = %d of them.", dx, dy, dx+dy, dx, n),
			map[string]string{"paths": edakit.IntStr(n)})
	}
}

// bbox: half-perimeter or area of a pin bounding box (KP 1), intersection area of two
// rectangles (KP 2), of three rectangles (KP 3).
func bbox(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		ps := distinctPts(n, 60)
		minx, maxx, miny, maxy := ps[0].x, ps[0].x, ps[0].y, ps[0].y
		for _, p := range ps {
			minx, maxx = imin(minx, p.x), imax(maxx, p.x)
			miny, maxy = imin(miny, p.y), imax(maxy, p.y)
		}
		if edakit.Int(0, 1) == 0 {
			hp := (maxx - minx) + (maxy - miny)
			q := edakit.Sprintf("A net has pins at %s (um). What is the half-perimeter of their bounding box (the HPWL estimate), in um?", listPts(ps))
			return edakit.Numeric(q, edakit.IntStr(hp),
				edakit.Sprintf("The box spans x from %d to %d and y from %d to %d, so the half-perimeter is %d + %d = %d um.", minx, maxx, miny, maxy, maxx-minx, maxy-miny, hp),
				map[string]string{"hpwl": edakit.IntStr(hp)})
		}
		ar := (maxx - minx) * (maxy - miny)
		q := edakit.Sprintf("A net has pins at %s (um). What is the area of their bounding box, in square um?", listPts(ps))
		return edakit.Numeric(q, edakit.IntStr(ar),
			edakit.Sprintf("The box spans x from %d to %d and y from %d to %d, so its area is %d x %d = %d.", minx, maxx, miny, maxy, maxx-minx, maxy-miny, ar),
			map[string]string{"area": edakit.IntStr(ar)})
	case 2:
		a, b := randRect(20, 3, 12), randRect(20, 3, 12)
		for i := 0; i < 20 && interArea(a, b) == 0 && edakit.Int(0, 4) > 0; i++ {
			b = randRect(20, 3, 12)
		}
		ox, oy := overlap1(a.x1, a.x2, b.x1, b.x2), overlap1(a.y1, a.y2, b.y1, b.y2)
		ar := interArea(a, b)
		q := edakit.Sprintf("Two layout rectangles are A = %s and B = %s, in units of x-range times y-range. What is the area of their intersection?", a, b)
		return edakit.Numeric(q, edakit.IntStr(ar),
			edakit.Sprintf("The overlap in x is min(%d, %d) - max(%d, %d) = %d and in y is min(%d, %d) - max(%d, %d) = %d. A negative or zero overlap means no intersection; otherwise multiply. Area = %d.", a.x2, b.x2, a.x1, b.x1, ox, a.y2, b.y2, a.y1, b.y1, oy, ar),
			map[string]string{"ox": edakit.IntStr(ox), "oy": edakit.IntStr(oy)})
	default:
		var rs []rect
		var c rect
		for tries := 0; tries < 50; tries++ {
			rs = []rect{randRect(20, 6, 14), randRect(20, 6, 14), randRect(20, 6, 14)}
			c = rect{imax(imax(rs[0].x1, rs[1].x1), rs[2].x1), imax(imax(rs[0].y1, rs[1].y1), rs[2].y1),
				imin(imin(rs[0].x2, rs[1].x2), rs[2].x2), imin(imin(rs[0].y2, rs[1].y2), rs[2].y2)}
			if c.x2 > c.x1 && c.y2 > c.y1 || edakit.Int(0, 5) == 0 {
				break
			}
		}
		ar := 0
		if c.x2 > c.x1 && c.y2 > c.y1 {
			ar = c.area()
		}
		q := edakit.Sprintf("Three rectangles on one layer are %s. What is the area of the region covered by all three?", listRects("R", rs))
		return edakit.Numeric(q, edakit.IntStr(ar),
			edakit.Sprintf("Intersect the x-ranges and the y-ranges: the common box runs from max of the left edges %d to min of the right edges %d, and from %d to %d in y. Area = %d.", c.x1, c.x2, c.y1, c.y2, ar),
			map[string]string{"area": edakit.IntStr(ar)})
	}
}

// orientation: the orientation determinant (KP 1), counting proper segment crossings (KP 2),
// the bit width that holds the determinant exactly (KP 3).
func orientation(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		var a, b, c pt
		for {
			a, b, c = pt{edakit.Int(-9, 9), edakit.Int(-9, 9)}, pt{edakit.Int(-9, 9), edakit.Int(-9, 9)}, pt{edakit.Int(-9, 9), edakit.Int(-9, 9)}
			if cross(a, b, c) != 0 || edakit.Int(0, 9) == 0 {
				break
			}
		}
		v := cross(a, b, c)
		q := edakit.Sprintf("For points $a = %s$, $b = %s$ and $c = %s$, compute the orientation determinant $(b_x - a_x)(c_y - a_y) - (b_y - a_y)(c_x - a_x)$.", a, b, c)
		return edakit.Numeric(q, edakit.IntStr(v),
			edakit.Sprintf("(%d)(%d) - (%d)(%d) = %d. A positive value means a left (counter-clockwise) turn, negative a right turn, zero collinear.", b.x-a.x, c.y-a.y, b.y-a.y, c.x-a.x, v),
			map[string]string{"orient": edakit.IntStr(v)})
	case 2:
		type seg struct{ p, q pt }
		var pairs [][2]seg
		count := 0
		for len(pairs) < 3 {
			s1 := seg{pt{edakit.Int(0, 10), edakit.Int(0, 10)}, pt{edakit.Int(0, 10), edakit.Int(0, 10)}}
			s2 := seg{pt{edakit.Int(0, 10), edakit.Int(0, 10)}, pt{edakit.Int(0, 10), edakit.Int(0, 10)}}
			d1, d2 := cross(s1.p, s1.q, s2.p), cross(s1.p, s1.q, s2.q)
			d3, d4 := cross(s2.p, s2.q, s1.p), cross(s2.p, s2.q, s1.q)
			if d1 == 0 || d2 == 0 || d3 == 0 || d4 == 0 {
				continue
			}
			pairs = append(pairs, [2]seg{s1, s2})
			if (d1 > 0) != (d2 > 0) && (d3 > 0) != (d4 > 0) {
				count++
			}
		}
		q := "Using orientation tests, decide which of these segment pairs cross: "
		for i, pr := range pairs {
			if i > 0 {
				q += "; "
			}
			q += edakit.Sprintf("pair %d: %s-%s with %s-%s", i+1, pr[0].p, pr[0].q, pr[1].p, pr[1].q)
		}
		q += ". No endpoint is collinear with the other segment. How many of the three pairs cross?"
		return edakit.Numeric(q, edakit.IntStr(count),
			"Two segments cross exactly when each one's endpoints lie on opposite sides of the other, that is, when both pairs of orientation signs differ. "+edakit.Sprintf("That holds for %d of the pairs.", count),
			map[string]string{"crossings": edakit.IntStr(count)})
	default:
		b := edakit.Pick([]int{12, 16, 20, 24, 28, 30, 32})
		ans := 2*b + 1
		q := edakit.Sprintf("Layout coordinates are signed %d-bit integers. The orientation determinant of three such points is computed exactly in integer arithmetic. What is the fewest bits a signed integer type needs so that every possible value of the determinant fits?", b)
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("A coordinate difference has magnitude up to 2^%d - 1, each product up to (2^%d - 1)^2, and the determinant (twice a triangle's area inside the coordinate square) is bounded by the same value, which is below 2^%d. A signed type therefore needs %d bits.", b, b, 2*b, ans),
			map[string]string{"bits": edakit.IntStr(ans)})
	}
}

// spacing: axis gap (KP 1), corner-to-corner Euclidean spacing (KP 2), counting violations of
// a Euclidean spacing rule around a reference shape (KP 3).
func spacing(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		a := rect{0, 0, edakit.Int(20, 80), edakit.Int(20, 80)}
		g := edakit.Int(5, 60)
		var b rect
		if edakit.Int(0, 1) == 0 {
			y := edakit.Int(-15, a.y2-5)
			b = rect{a.x2 + g, y, a.x2 + g + edakit.Int(10, 50), y + edakit.Int(10, 60)}
		} else {
			x := edakit.Int(-15, a.x2-5)
			b = rect{x, a.y2 + g, x + edakit.Int(10, 60), a.y2 + g + edakit.Int(10, 50)}
		}
		q := edakit.Sprintf("Two metal shapes are A = %s and B = %s, in nm. What is the Euclidean edge-to-edge spacing between them, in nm?", a, b)
		return edakit.Numeric(q, edakit.IntStr(g),
			edakit.Sprintf("The projections overlap on one axis, so the closest points face each other straight across the gap of %d nm.", g),
			map[string]string{"gap": edakit.IntStr(g)})
	case 2:
		a := rect{0, 0, edakit.Int(20, 60), edakit.Int(20, 60)}
		gx, gy := edakit.Int(3, 50), edakit.Int(3, 50)
		b := rect{a.x2 + gx, a.y2 + gy, a.x2 + gx + edakit.Int(10, 40), a.y2 + gy + edakit.Int(10, 40)}
		d := math.Hypot(float64(gx), float64(gy))
		ans := edakit.Num(d, 2)
		q := edakit.Sprintf("Two metal shapes are A = %s and B = %s, in nm. What is the Euclidean spacing between them, in nm? Round to 2 decimal places.", a, b)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Neither projection overlaps: the gaps are %d nm in x and %d nm in y, so the closest points are two corners, sqrt(%d^2 + %d^2) = %s nm apart.", gx, gy, gx, gy, ans),
			map[string]string{"gx": edakit.IntStr(gx), "gy": edakit.IntStr(gy)})
	default:
		s := 10 * edakit.Int(4, 8)
		ref := rect{0, 0, 100, 100}
		var nb []rect
		count := 0
		for len(nb) < 4 {
			gx, gy := edakit.Int(0, s+15), edakit.Int(0, s+15)
			if edakit.Int(0, 2) == 0 {
				gx = 0
			}
			if gx == 0 && gy == 0 {
				continue
			}
			sx, sy := edakit.Pick([]int{-1, 1}), edakit.Pick([]int{-1, 1})
			w, h := edakit.Int(20, 40), edakit.Int(20, 40)
			var r rect
			if gx == 0 {
				x := edakit.Int(10, 70)
				r = rect{x, 0, x + w, 0}
			} else if sx > 0 {
				r = rect{ref.x2 + gx, 0, ref.x2 + gx + w, 0}
			} else {
				r = rect{ref.x1 - gx - w, 0, ref.x1 - gx, 0}
			}
			if sy > 0 {
				r.y1, r.y2 = ref.y2+gy, ref.y2+gy+h
			} else {
				r.y1, r.y2 = ref.y1-gy-h, ref.y1-gy
			}
			if gy == 0 {
				y := edakit.Int(10, 70)
				r.y1, r.y2 = y, y+h
			}
			ok := true
			for _, o := range nb {
				if closedMeet(o, r) {
					ok = false
				}
			}
			if !ok {
				continue
			}
			nb = append(nb, r)
			if gx*gx+gy*gy < s*s {
				count++
			}
		}
		q := edakit.Sprintf("A spacing rule requires Euclidean spacing of at least %d nm. The reference shape is %s and its neighbours are %s, all in nm. How many neighbours violate the rule against the reference shape?", s, ref, listRects("N", nb))
		return edakit.Numeric(q, edakit.IntStr(count),
			"For each neighbour take the x-gap and y-gap (zero where projections overlap) and compare gx^2 + gy^2 with the squared rule, which avoids square roots. "+edakit.Sprintf("%d neighbours fall short of %d nm.", count, s),
			map[string]string{"violations": edakit.IntStr(count)})
	}
}
