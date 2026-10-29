package geom

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type ival struct{ a, b int }

func randIvals(n, span int) []ival {
	out := make([]ival, n)
	for i := range out {
		a := edakit.Int(0, span-2)
		out[i] = ival{a, a + edakit.Int(1, imin(span-a, span/2))}
	}
	return out
}

func ivalText(iv []ival) string {
	parts := make([]string, len(iv))
	for i, v := range iv {
		parts[i] = edakit.Sprintf("[%d, %d]", v.a, v.b)
	}
	return strings.Join(parts, ", ")
}

// intervalTree: stabbing count (KP 1), intervals stored at the root of a centred interval tree
// (KP 2), intervals the root reports for a query (KP 3).
func intervalTree(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 9)
	iv := randIvals(n, 40)
	switch kp(ctx) {
	case 1:
		q0 := edakit.Int(2, 38)
		c := 0
		for _, v := range iv {
			if v.a <= q0 && q0 <= v.b {
				c++
			}
		}
		q := edakit.Sprintf("The y-extents of wires crossing a cut line are the closed intervals %s. A stabbing query asks which intervals contain y = %d. How many does it report?", ivalText(iv), q0)
		return edakit.Numeric(q, edakit.IntStr(c),
			edakit.Sprintf("An interval [a, b] contains %d when a <= %d <= b. %d intervals qualify.", q0, q0, c),
			map[string]string{"stabbed": edakit.IntStr(c)})
	case 2:
		c := edakit.Int(12, 28)
		at, left := 0, 0
		for _, v := range iv {
			if v.a <= c && c <= v.b {
				at++
			} else if v.b < c {
				left++
			}
		}
		q := edakit.Sprintf("A centred interval tree is built over %s with root centre x = %d. Intervals containing the centre stay at the root; the rest go to the left or right subtree. How many intervals are stored at the root?", ivalText(iv), c)
		return edakit.Numeric(q, edakit.IntStr(at),
			edakit.Sprintf("%d intervals contain %d and stay at the root; %d lie entirely to the left and %d entirely to the right.", at, c, left, n-at-left),
			map[string]string{"root": edakit.IntStr(at), "left": edakit.IntStr(left)})
	default:
		c := edakit.Int(14, 26)
		var root []ival
		for len(root) < 3 {
			root = nil
			iv = randIvals(n, 40)
			for _, v := range iv {
				if v.a <= c && c <= v.b {
					root = append(root, v)
				}
			}
		}
		var q0, rep int
		var how string
		if edakit.Int(0, 1) == 0 {
			q0 = edakit.Int(c-10, c-1)
			for _, v := range root {
				if v.a <= q0 {
					rep++
				}
			}
			how = "left endpoints in increasing order, stopping at the first left endpoint greater than the query"
		} else {
			q0 = edakit.Int(c+1, c+10)
			for _, v := range root {
				if v.b >= q0 {
					rep++
				}
			}
			how = "right endpoints in decreasing order, stopping at the first right endpoint less than the query"
		}
		q := edakit.Sprintf("The root of a centred interval tree has centre %d and stores the intervals %s. A stabbing query for %d scans the root's list sorted by %s. How many intervals does the root report?", c, ivalText(root), q0, how)
		return edakit.Numeric(q, edakit.IntStr(rep),
			edakit.Sprintf("Every root interval contains the centre %d, so on the query's side only one endpoint needs testing. %d of them reach %d. The scan stops after one failure, so the root costs O(1 + reported).", c, rep, q0),
			map[string]string{"reported": edakit.IntStr(rep)})
	}
}

// canon counts the canonical nodes of [a, b) in a segment tree over [lo, hi).
func canon(lo, hi, a, b int) int {
	if b <= lo || hi <= a {
		return 0
	}
	if a <= lo && hi <= b {
		return 1
	}
	mid := (lo + hi) / 2
	return canon(lo, mid, a, b) + canon(mid, hi, a, b)
}

// segmentTree: elementary intervals (KP 1), canonical nodes of one interval (KP 2), total
// storage for several intervals (KP 3).
func segmentTree(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		iv := randIvals(n, 20)
		var xs []int
		for _, v := range iv {
			xs = append(xs, v.a, v.b)
		}
		m := len(uniq(xs))
		q := edakit.Sprintf("A segment tree is built over the intervals %s. Its leaves are the elementary intervals: the open gaps between consecutive distinct endpoints, the endpoints themselves as points, and the two unbounded ends. How many leaves are there?", ivalText(iv))
		return edakit.Numeric(q, edakit.IntStr(2*m+1),
			edakit.Sprintf("There are %d distinct endpoints, giving %d point leaves and %d open intervals: 2m + 1 = %d.", m, m, m+1, 2*m+1),
			map[string]string{"distinct": edakit.IntStr(m)})
	case 2:
		L := edakit.Pick([]int{8, 16, 32})
		a := edakit.Int(0, L-2)
		b := edakit.Int(a+1, L)
		c := canon(0, L, a, b)
		q := edakit.Sprintf("A segment tree has %d unit leaves covering [0, %d), so each node covers an aligned power-of-two range. The interval [%d, %d) is stored at its canonical nodes: the highest nodes whose range lies inside it. How many canonical nodes does it use?", L, L, a, b)
		return edakit.Numeric(q, edakit.IntStr(c),
			edakit.Sprintf("Descend from the root, stopping at any node entirely inside [%d, %d) and discarding nodes entirely outside. %d nodes are stopped at, never more than two per level.", a, b, c),
			map[string]string{"canonical": edakit.IntStr(c)})
	default:
		L := edakit.Pick([]int{16, 32})
		k := edakit.Int(3, 4)
		tot := 0
		var parts []string
		for i := 0; i < k; i++ {
			a := edakit.Int(0, L-2)
			b := edakit.Int(a+1, L)
			tot += canon(0, L, a, b)
			parts = append(parts, edakit.Sprintf("[%d, %d)", a, b))
		}
		q := edakit.Sprintf("A segment tree has %d unit leaves over [0, %d). The intervals %s are each stored at their canonical nodes. How many interval entries does the tree store in total?", L, L, strings.Join(parts, ", "))
		return edakit.Numeric(q, edakit.IntStr(tot),
			edakit.Sprintf("Decompose each interval into canonical nodes and add the counts: %d entries. Each interval takes O(log n) entries, so storage is O(n log n), against O(n) for an interval tree.", tot),
			map[string]string{"entries": edakit.IntStr(tot)})
	}
}

func mbr(rs []rect) rect {
	m := rs[0]
	for _, r := range rs[1:] {
		m = rect{imin(m.x1, r.x1), imin(m.y1, r.y1), imax(m.x2, r.x2), imax(m.y2, r.y2)}
	}
	return m
}

// cluster draws k small boxes near a centre.
func cluster(k, cx, cy int) []rect {
	out := make([]rect, k)
	for i := range out {
		x, y := cx+edakit.Int(-6, 6), cy+edakit.Int(-6, 6)
		out[i] = rect{x, y, x + edakit.Int(1, 4), y + edakit.Int(1, 4)}
	}
	return out
}

// rtree: MBR of a node (KP 1), child nodes visited (KP 2), total overlap tests (KP 3).
func rtree(ctx generator.GeneratorContext) generator.Problem {
	centres := [][2]int{{10, 10}, {30, 12}, {18, 30}, {34, 32}}
	nn := edakit.Scaled(ctx, 3, 4)
	var nodes [][]rect
	for i := 0; i < nn; i++ {
		nodes = append(nodes, cluster(edakit.Int(3, 4), centres[i][0], centres[i][1]))
	}
	switch kp(ctx) {
	case 1:
		e := nodes[0]
		m := mbr(e)
		q := edakit.Sprintf("An R-tree leaf node holds the entries %s. What is the area of the node's minimum bounding rectangle?", listRects("E", e))
		return edakit.Numeric(q, edakit.IntStr(m.area()),
			edakit.Sprintf("Take the minimum and maximum of the entry edges: the MBR is %s, area %d.", m, m.area()),
			map[string]string{"mbr_area": edakit.IntStr(m.area())})
	case 2:
		qr := randRect(40, 6, 16)
		vis := 0
		var ms []rect
		for _, e := range nodes {
			m := mbr(e)
			ms = append(ms, m)
			if closedMeet(m, qr) {
				vis++
			}
		}
		q := edakit.Sprintf("An R-tree root has %d children with bounding rectangles %s. A window query asks for %s. How many children must the search descend into (touching counts as intersecting)?", nn, listRects("M", ms), qr)
		return edakit.Numeric(q, edakit.IntStr(vis),
			edakit.Sprintf("The search descends only into children whose MBR meets the window: %d of %d. The others are pruned with their whole subtrees.", vis, nn),
			map[string]string{"visited": edakit.IntStr(vis)})
	default:
		qr := randRect(40, 6, 16)
		tests, found := nn, 0
		var parts []string
		for i, e := range nodes {
			m := mbr(e)
			parts = append(parts, edakit.Sprintf("node %d with MBR %s holds %s", i+1, m, listRects("E", e)))
			if closedMeet(m, qr) {
				tests += len(e)
				for _, r := range e {
					if closedMeet(r, qr) {
						found++
					}
				}
			}
		}
		q := edakit.Sprintf("A two-level R-tree has a root whose children are: %s. A window query for %s tests the window against every root entry, then against every entry of each child it descends into. How many rectangle-intersection tests does it perform in total (touching counts as intersecting)?", strings.Join(parts, "; "), qr)
		return edakit.Numeric(q, edakit.IntStr(tests),
			edakit.Sprintf("%d tests at the root plus the entries of every child whose MBR meets the window gives %d tests, and %d entries are reported.", nn, tests, found),
			map[string]string{"tests": edakit.IntStr(tests), "found": edakit.IntStr(found)})
	}
}

// rtreeSplit: dead area of one pair (KP 1), the PickSeeds maximum (KP 2), R* split metrics
// (KP 3).
func rtreeSplit(ctx generator.GeneratorContext) generator.Problem {
	dead := func(a, b rect) int { return mbr([]rect{a, b}).area() - a.area() - b.area() }
	switch kp(ctx) {
	case 1:
		a, b := randRect(20, 2, 7), randRect(20, 2, 7)
		d := dead(a, b)
		q := edakit.Sprintf("Guttman's quadratic split scores a pair of entries by the area wasted if they shared a node: $d = \\text{area}(\\text{MBR}(E_1 \\cup E_2)) - \\text{area}(E_1) - \\text{area}(E_2)$. Compute d for $E_1$ = %s and $E_2$ = %s.", a, b)
		return edakit.Numeric(q, edakit.IntStr(d),
			edakit.Sprintf("The joint MBR has area %d; subtracting %d and %d gives %d. A negative d means the entries overlap.", mbr([]rect{a, b}).area(), a.area(), b.area(), d),
			map[string]string{"dead": edakit.IntStr(d)})
	case 2:
		es := randRects(4, 24, 2, 6)
		best := -1 << 30
		for i := range es {
			for j := i + 1; j < 4; j++ {
				best = imax(best, dead(es[i], es[j]))
			}
		}
		q := edakit.Sprintf("An overflowing R-tree node holds %s. PickSeeds chooses the pair with the largest dead area $d = \\text{area}(\\text{MBR}) - \\text{area}(E_i) - \\text{area}(E_j)$ as the seeds of the two new nodes. What is that largest d?", listRects("E", es))
		return edakit.Numeric(q, edakit.IntStr(best),
			edakit.Sprintf("Scoring all six pairs, the largest dead area is %d. Those two entries are the worst pair to keep together, so they seed different nodes.", best),
			map[string]string{"best_dead": edakit.IntStr(best)})
	default:
		es := randRects(4, 24, 2, 8)
		g1, g2 := mbr(es[:2]), mbr(es[2:])
		if edakit.Int(0, 1) == 0 {
			m := g1.per() + g2.per()
			q := edakit.Sprintf("An R*-tree considers splitting %s into {E1, E2} and {E3, E4}. Its axis choice sums the margins (perimeters) of the two group MBRs. What is that margin sum?", listRects("E", es))
			return edakit.Numeric(q, edakit.IntStr(m),
				edakit.Sprintf("The group MBRs are %s and %s with perimeters %d and %d, total %d. Small margins favour square-ish nodes, which answer window queries with fewer visits.", g1, g2, g1.per(), g2.per(), m),
				map[string]string{"margin": edakit.IntStr(m)})
		}
		o := interArea(g1, g2)
		q := edakit.Sprintf("An R*-tree considers splitting %s into {E1, E2} and {E3, E4}. It prefers the split whose two group MBRs overlap least. What is the overlap area of the two group MBRs?", listRects("E", es))
		return edakit.Numeric(q, edakit.IntStr(o),
			edakit.Sprintf("The group MBRs are %s and %s; their intersection has area %d. Overlap matters because a point query inside it must descend into both nodes.", g1, g2, o),
			map[string]string{"overlap": edakit.IntStr(o)})
	}
}

// kdTree: root split (KP 1), second-level split (KP 2), cells visited by a nearest-neighbour
// search with a known radius (KP 3).
func kdTree(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1, 2:
		n := edakit.Pick([]int{7, 9})
		var ps []pt
		usedX, usedY := map[int]bool{}, map[int]bool{}
		for len(ps) < n {
			p := pt{edakit.Int(0, 40), edakit.Int(0, 40)}
			if !usedX[p.x] && !usedY[p.y] {
				usedX[p.x], usedY[p.y] = true, true
				ps = append(ps, p)
			}
		}
		sorted := append([]pt(nil), ps...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].x < sorted[j].x })
		med := sorted[n/2]
		if k == 1 {
			q := edakit.Sprintf("A 2-d tree is built over the cell centres %s. The root splits on x at the median point, which goes to the root. What is the x-coordinate of the root's split?", listPts(ps))
			return edakit.Numeric(q, edakit.IntStr(med.x),
				edakit.Sprintf("Sorting by x, the middle (number %d of %d) point is %s, so the root splits at x = %d.", n/2+1, n, med, med.x),
				map[string]string{"split": edakit.IntStr(med.x)})
		}
		left := append([]pt(nil), sorted[:n/2]...)
		sort.Slice(left, func(i, j int) bool { return left[i].y < left[j].y })
		lm := left[len(left)/2]
		q := edakit.Sprintf("A 2-d tree is built over the cell centres %s. The root splits on x at the median point; each child then splits on y at the median of its own points (the middle one, or the upper middle for an even count). What is the y-coordinate of the split in the root's left child (the points with smaller x)?", listPts(ps))
		return edakit.Numeric(q, edakit.IntStr(lm.y),
			edakit.Sprintf("The root takes %s. The %d points left of it have y-values sorted %s; the split is at y = %d.", med, len(left), func() string {
				var ys []int
				for _, p := range left {
					ys = append(ys, p.y)
				}
				return joinInts(ys)
			}(), lm.y),
			map[string]string{"split": edakit.IntStr(lm.y)})
	default:
		W := 40
		s := edakit.Int(12, 28)
		t1, t2 := edakit.Int(10, 30), edakit.Int(10, 30)
		cells := []rect{{0, 0, s, t1}, {0, t1, s, W}, {s, 0, W, t2}, {s, t2, W, W}}
		qp := pt{edakit.Int(2, 38), edakit.Int(2, 38)}
		r := edakit.Int(3, 14)
		vis := 0
		for _, c := range cells {
			dx := imax(0, imax(c.x1-qp.x, qp.x-c.x2))
			dy := imax(0, imax(c.y1-qp.y, qp.y-c.y2))
			if dx+dy < r {
				vis++
			}
		}
		q := edakit.Sprintf("A 2-d tree over the region $[0, 40] \\times [0, 40]$ splits at x = %d, then the left half at y = %d and the right half at y = %d, giving four leaf cells. A Manhattan nearest-neighbour search from %s already holds a candidate at distance %d. A cell must be searched only if its L1 distance from the query is less than %d. How many of the four cells must be searched?", s, t1, t2, qp, r, r)
		return edakit.Numeric(q, edakit.IntStr(vis),
			edakit.Sprintf("The L1 distance to a cell is its x-gap plus its y-gap (zero when inside). %d cells come closer than %d, including the one holding the query; the rest are pruned.", vis, r),
			map[string]string{"visited": edakit.IntStr(vis)})
	}
}

// quadtree: depth for a target cell size (KP 1), depth separating two points (KP 2), cells a
// rectangle overlaps at one depth (KP 3).
func quadtree(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Int(8, 12)
	W := 1 << k
	switch kp(ctx) {
	case 1:
		s := edakit.Int(3, W/4)
		d := 0
		for (W >> d) > s {
			d++
		}
		q := edakit.Sprintf("A region quadtree covers a %d by %d unit die; each level halves the cell side. What is the smallest depth (root at depth 0) whose cells have side at most %d units?", W, W, s)
		return edakit.Numeric(q, edakit.IntStr(d),
			edakit.Sprintf("Cells at depth d have side %d / 2^d. That first reaches %d or less at depth %d, where the side is %d.", W, s, d, W>>d),
			map[string]string{"depth": edakit.IntStr(d)})
	case 2:
		var a, b pt
		for a == b {
			a, b = pt{edakit.Int(0, W-1), edakit.Int(0, W-1)}, pt{edakit.Int(0, W-1), edakit.Int(0, W-1)}
			if edakit.Int(0, 1) == 0 {
				b = pt{imin(W-1, a.x+edakit.Int(0, W/16)), imin(W-1, a.y+edakit.Int(0, W/16))}
			}
		}
		d := 1
		for {
			c := W >> d
			if a.x/c != b.x/c || a.y/c != b.y/c {
				break
			}
			d++
		}
		q := edakit.Sprintf("A point-region quadtree covers $[0, %d)^2$ with half-open cells. Points %s and %s start in the root cell (depth 0). At what depth do they first fall into different cells?", W, a, b)
		return edakit.Numeric(q, edakit.IntStr(d),
			edakit.Sprintf("At depth d the cell index is the coordinate divided by %d / 2^d, rounded down. The indices first differ at depth %d (cell side %d).", W, d, W>>d),
			map[string]string{"depth": edakit.IntStr(d)})
	default:
		d := edakit.Int(3, 5)
		c := W >> d
		x1 := edakit.Int(0, W-c)
		x2 := edakit.Int(x1+1, imin(W, x1+3*c))
		y1 := edakit.Int(0, W-c)
		y2 := edakit.Int(y1+1, imin(W, y1+3*c))
		nx := (x2-1)/c - x1/c + 1
		ny := (y2-1)/c - y1/c + 1
		q := edakit.Sprintf("A quadtree over $[0, %d)^2$ is cut to depth %d, giving half-open square cells of side %d. How many depth-%d cells does the rectangle $[%d, %d) \\times [%d, %d)$ overlap?", W, d, c, d, x1, x2, y1, y2)
		return edakit.Numeric(q, edakit.IntStr(nx*ny),
			edakit.Sprintf("In x the rectangle runs from cell %d to cell %d (%d cells) and in y from %d to %d (%d cells): %d cells.", x1/c, (x2-1)/c, nx, y1/c, (y2-1)/c, ny, nx*ny),
			map[string]string{"nx": edakit.IntStr(nx), "ny": edakit.IntStr(ny)})
	}
}

// uniformGrid: bin index (KP 1), stored entries (KP 2), expected candidates (KP 3).
func uniformGrid(ctx generator.GeneratorContext) generator.Problem {
	b := edakit.Pick([]int{20, 25, 40, 50, 100})
	C := edakit.Int(8, 20)
	R := edakit.Int(8, 20)
	switch kp(ctx) {
	case 1:
		x, y := edakit.Int(0, C*b-1), edakit.Int(0, R*b-1)
		idx := (y/b)*C + x/b
		q := edakit.Sprintf("A uniform grid of %d by %d nm bins has %d columns and %d rows, bins numbered row by row from 0 at the origin. Bins are half-open. Which bin index holds the point (%d, %d)?", b, b, C, R, x, y)
		return edakit.Numeric(q, edakit.IntStr(idx),
			edakit.Sprintf("Column = floor(%d / %d) = %d, row = floor(%d / %d) = %d, index = %d x %d + %d = %d.", x, b, x/b, y, b, y/b, y/b, C, x/b, idx),
			map[string]string{"col": edakit.IntStr(x / b), "row": edakit.IntStr(y / b)})
	case 2:
		k := edakit.Scaled(ctx, 3, 5)
		tot := 0
		var rs []rect
		for i := 0; i < k; i++ {
			x1, y1 := edakit.Int(0, 6*b), edakit.Int(0, 6*b)
			r := rect{x1, y1, x1 + edakit.Int(5, 2*b), y1 + edakit.Int(5, 2*b)}
			rs = append(rs, r)
			tot += ((r.x2-1)/b - r.x1/b + 1) * ((r.y2-1)/b - r.y1/b + 1)
		}
		q := edakit.Sprintf("Shapes are registered in every %d nm half-open bin they overlap. For the shapes %s (half-open, in nm), how many bin entries are stored in total?", b, listRects("S", rs))
		return edakit.Numeric(q, edakit.IntStr(tot),
			edakit.Sprintf("Count the bins each shape spans in x and in y and multiply, then add: %d entries. A shape spanning several bins is stored several times, so a query must de-duplicate its results.", tot),
			map[string]string{"entries": edakit.IntStr(tot)})
	default:
		n := 1000 * edakit.Int(5, 200)
		qb := edakit.Int(2, 9)
		e := float64(n) * float64(qb) / float64(C*R)
		ans := edakit.Num(e, 1)
		q := edakit.Sprintf("A die is divided into %d by %d bins, and %d small shapes are spread uniformly over it, each in one bin. A window query touches %d bins. How many candidate shapes does it fetch on average? Round to 1 decimal place.", C, R, n, qb)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Each bin holds n / bins = %d / %d shapes on average; %d bins give %s. Clustered layouts break this average, which is why adaptive trees beat grids on uneven data.", n, C*R, qb, ans),
			map[string]string{"per_bin": edakit.Num(float64(n)/float64(C*R), 2)})
	}
}
