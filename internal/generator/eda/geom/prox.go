package geom

import (
	"sort"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// hull returns the strict convex hull, counter-clockwise, collinear points dropped, and the
// size of the lower hull (both end points included).
func hull(ps []pt) ([]pt, int) {
	s := append([]pt(nil), ps...)
	sort.Slice(s, func(i, j int) bool {
		if s[i].x != s[j].x {
			return s[i].x < s[j].x
		}
		return s[i].y < s[j].y
	})
	var lo, up []pt
	for _, p := range s {
		for len(lo) >= 2 && cross(lo[len(lo)-2], lo[len(lo)-1], p) <= 0 {
			lo = lo[:len(lo)-1]
		}
		lo = append(lo, p)
	}
	for i := len(s) - 1; i >= 0; i-- {
		p := s[i]
		for len(up) >= 2 && cross(up[len(up)-2], up[len(up)-1], p) <= 0 {
			up = up[:len(up)-1]
		}
		up = append(up, p)
	}
	h := append(append([]pt(nil), lo[:len(lo)-1]...), up[:len(up)-1]...)
	return h, len(lo)
}

// onHullBoundaryOnly reports whether some point lies on a hull edge without being a vertex.
func onHullBoundaryOnly(ps, h []pt) bool {
	isV := map[pt]bool{}
	for _, v := range h {
		isV[v] = true
	}
	for _, p := range ps {
		if isV[p] {
			continue
		}
		for i := range h {
			a, b := h[i], h[(i+1)%len(h)]
			if cross(a, b, p) == 0 && imin(a.x, b.x) <= p.x && p.x <= imax(a.x, b.x) && imin(a.y, b.y) <= p.y && p.y <= imax(a.y, b.y) {
				return true
			}
		}
	}
	return false
}

// generalPts draws n distinct points with no point on a hull edge between hull vertices.
func generalPts(n, span int) ([]pt, []pt) {
	for {
		ps := distinctPts(n, span)
		h, _ := hull(ps)
		if len(h) >= 3 && !onHullBoundaryOnly(ps, h) {
			return ps, h
		}
	}
}

// convexHullGen: hull vertex count (KP 1), lower-hull size in Andrew's chain (KP 2), hull area of
// rectangle corners (KP 3).
func convexHullGen(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		ps := distinctPts(edakit.Scaled(ctx, 6, 10), 10)
		h, _ := hull(ps)
		q := edakit.Sprintf("How many vertices does the convex hull of the points %s have? Points lying on a hull edge between two corners are not vertices.", listPts(ps))
		return edakit.Numeric(q, edakit.IntStr(len(h)),
			edakit.Sprintf("Andrew's monotone chain sorts by x and keeps only left turns. The hull is %s: %d vertices.", listPts(h), len(h)),
			map[string]string{"hull": edakit.IntStr(len(h))})
	case 2:
		ps := distinctPts(edakit.Scaled(ctx, 6, 10), 10)
		_, nl := hull(ps)
		q := edakit.Sprintf("Andrew's monotone chain sorts the points %s by x (then y) and builds the lower hull left to right, popping the last point while it does not make a strict left turn. How many points does the finished lower hull contain, counting both ends?", listPts(ps))
		return edakit.Numeric(q, edakit.IntStr(nl),
			edakit.Sprintf("Scanning the sorted points and popping on every right turn or collinear step leaves %d points on the lower chain, from the leftmost to the rightmost point.", nl),
			map[string]string{"lower": edakit.IntStr(nl)})
	default:
		rs := randRects(edakit.Int(2, 3), 12, 2, 6)
		var ps []pt
		seen := map[pt]bool{}
		for _, r := range rs {
			for _, p := range []pt{{r.x1, r.y1}, {r.x2, r.y1}, {r.x2, r.y2}, {r.x1, r.y2}} {
				if !seen[p] {
					seen[p] = true
					ps = append(ps, p)
				}
			}
		}
		h, _ := hull(ps)
		a2 := ringArea2(h)
		ans := edakit.Num(float64(a2)/2, 1)
		q := edakit.Sprintf("A placement blockage must enclose the cells %s. What is the area of the convex hull of these rectangles? Round to 1 decimal place.", listRects("C", rs))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The hull of rectangles is the hull of their corners: %s. Its shoelace area is %s.", listPts(h), ans),
			map[string]string{"twice_area": edakit.IntStr(a2), "hull": edakit.IntStr(len(h))})
	}
}

// voronoi: worst-case vertex or edge counts (KP 1), unbounded cells (KP 2), exact counts from
// n and h (KP 3).
func voronoi(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 2000)
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("The Voronoi diagram of %d sites (not all collinear) is a planar subdivision. What is the largest number of Voronoi vertices it can have?", n)
			return edakit.Numeric(q, edakit.IntStr(2*n-5),
				edakit.Sprintf("Euler's formula on the diagram, with every vertex of degree at least 3, bounds the vertices by 2n - 5 = %d.", 2*n-5),
				map[string]string{"bound": edakit.IntStr(2*n - 5)})
		}
		q := edakit.Sprintf("The Voronoi diagram of %d sites (not all collinear) is a planar subdivision. What is the largest number of Voronoi edges it can have?", n)
		return edakit.Numeric(q, edakit.IntStr(3*n-6),
			edakit.Sprintf("Euler's formula with every vertex of degree at least 3 bounds the edges by 3n - 6 = %d, so the diagram has linear size.", 3*n-6),
			map[string]string{"bound": edakit.IntStr(3*n - 6)})
	case 2:
		ps, h := generalPts(edakit.Scaled(ctx, 6, 10), 12)
		q := edakit.Sprintf("Sites are placed at %s. How many cells of their Voronoi diagram are unbounded?", listPts(ps))
		return edakit.Numeric(q, edakit.IntStr(len(h)),
			edakit.Sprintf("A site's cell is unbounded exactly when the site lies on the convex hull. The hull vertices are %s, so %d cells are unbounded.", listPts(h), len(h)),
			map[string]string{"hull": edakit.IntStr(len(h))})
	default:
		n := edakit.Scaled(ctx, 6, 500)
		h := edakit.Int(3, imax(3, n/3))
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("%d sites in general position (no four on a circle) have %d sites on their convex hull. How many vertices does their Voronoi diagram have?", n, h)
			return edakit.Numeric(q, edakit.IntStr(2*n-2-h),
				edakit.Sprintf("Voronoi vertices correspond one-to-one with Delaunay triangles, and a triangulation of n points with h on the hull has 2n - 2 - h = %d triangles.", 2*n-2-h),
				map[string]string{"vertices": edakit.IntStr(2*n - 2 - h)})
		}
		q := edakit.Sprintf("%d sites in general position (no four on a circle) have %d sites on their convex hull. How many edges does their Voronoi diagram have, bounded and unbounded?", n, h)
		return edakit.Numeric(q, edakit.IntStr(3*n-3-h),
			edakit.Sprintf("Voronoi edges are dual to Delaunay edges, and a triangulation with h hull points has 3n - 3 - h = %d edges.", 3*n-3-h),
			map[string]string{"edges": edakit.IntStr(3*n - 3 - h)})
	}
}

func incircle(a, b, c, d pt) int {
	ax, ay := a.x-d.x, a.y-d.y
	bx, by := b.x-d.x, b.y-d.y
	cx, cy := c.x-d.x, c.y-d.y
	a2, b2, c2 := ax*ax+ay*ay, bx*bx+by*by, cx*cx+cy*cy
	return ax*(by*c2-b2*cy) - ay*(bx*c2-b2*cx) + a2*(bx*cy-by*cx)
}

// delaunay: triangle or edge counts (KP 1), the incircle determinant (KP 2), counts after
// inserting interior points (KP 3).
func delaunay(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 400)
		h := edakit.Int(3, imax(3, n/2))
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("A Delaunay triangulation is built over %d points, %d of them on the convex hull. How many triangles does it have?", n, h)
			return edakit.Numeric(q, edakit.IntStr(2*n-2-h),
				edakit.Sprintf("Every triangulation of the same point set has 2n - 2 - h = %d triangles, by Euler's formula.", 2*n-2-h),
				map[string]string{"triangles": edakit.IntStr(2*n - 2 - h)})
		}
		q := edakit.Sprintf("A Delaunay triangulation is built over %d points, %d of them on the convex hull. How many edges does it have?", n, h)
		return edakit.Numeric(q, edakit.IntStr(3*n-3-h),
			edakit.Sprintf("Every triangulation of the same point set has 3n - 3 - h = %d edges.", 3*n-3-h),
			map[string]string{"edges": edakit.IntStr(3*n - 3 - h)})
	case 2:
		var a, b, c, d pt
		for {
			a, b, c, d = pt{edakit.Int(0, 8), edakit.Int(0, 8)}, pt{edakit.Int(0, 8), edakit.Int(0, 8)}, pt{edakit.Int(0, 8), edakit.Int(0, 8)}, pt{edakit.Int(0, 8), edakit.Int(0, 8)}
			o := cross(a, b, c)
			if o < 0 {
				b, c = c, b
				o = -o
			}
			if o > 0 && d != a && d != b && d != c {
				break
			}
		}
		v := incircle(a, b, c, d)
		q := edakit.Sprintf("Triangle $a = %s$, $b = %s$, $c = %s$ is counter-clockwise. Compute the incircle determinant for $d = %s$: the 3 by 3 determinant with rows $(a_x - d_x,\\ a_y - d_y,\\ (a_x - d_x)^2 + (a_y - d_y)^2)$ and likewise for b and c.", a, b, c, d)
		return edakit.Numeric(q, edakit.IntStr(v),
			edakit.Sprintf("The determinant is %d. Positive means d is inside the circumcircle of abc (the edge must flip), negative outside, zero on the circle.", v),
			map[string]string{"incircle": edakit.IntStr(v)})
	default:
		n := edakit.Scaled(ctx, 5, 300)
		h := edakit.Int(3, imax(3, n/2))
		m := edakit.Int(1, 50)
		if edakit.Int(0, 1) == 0 {
			q := edakit.Sprintf("An extraction mesh is a Delaunay triangulation of %d points, %d on the hull. The mesher inserts %d more points strictly inside the hull. How many triangles does the refined mesh have?", n, h, m)
			return edakit.Numeric(q, edakit.IntStr(2*(n+m)-2-h),
				edakit.Sprintf("Interior points leave the hull unchanged, and each adds two triangles: 2(n + m) - 2 - h = %d.", 2*(n+m)-2-h),
				map[string]string{"triangles": edakit.IntStr(2*(n+m) - 2 - h)})
		}
		q := edakit.Sprintf("An extraction mesh is a Delaunay triangulation of %d points, %d on the hull. The mesher inserts %d more points strictly inside the hull. How many edges does the refined mesh have?", n, h, m)
		return edakit.Numeric(q, edakit.IntStr(3*(n+m)-3-h),
			edakit.Sprintf("Interior points leave the hull unchanged, and each adds three edges: 3(n + m) - 3 - h = %d.", 3*(n+m)-3-h),
			map[string]string{"edges": edakit.IntStr(3*(n+m) - 3 - h)})
	}
}
