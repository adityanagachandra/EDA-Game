package route

import (
	"math"
	"sort"
	"strings"

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

// lz: Z-shape count (KP 1), all L and Z options (KP 2), monotone routes that patterns miss
// (KP 3).
func lz(ctx generator.GeneratorContext) generator.Problem {
	dx, dy := edakit.Scaled(ctx, 2, 9), edakit.Scaled(ctx, 2, 9)
	x1, y1 := edakit.Int(0, 20), edakit.Int(0, 20)
	a, b := pt{x1, y1}, pt{x1 + dx, y1 + dy}
	intro := edakit.Sprintf("A two-pin net joins gcells %s and %s on the global-routing grid. ", ptText(a), ptText(b))
	zdef := "A Z-shaped route is monotone with exactly two bends: horizontal, vertical, horizontal with the vertical jog at a column strictly between the pins, or vertical, horizontal, vertical with the horizontal jog at a row strictly between them."
	switch kp(ctx) {
	case 1:
		z := dx + dy - 2
		return edakit.Numeric(intro+zdef+" How many distinct Z-shaped routes does a pattern router consider?", edakit.IntStr(z),
			edakit.Sprintf("HVH shapes can jog at any of the %d interior columns and VHV shapes at any of the %d interior rows: %d + %d = %d.", dx-1, dy-1, dx-1, dy-1, z),
			map[string]string{"hvh": edakit.IntStr(dx - 1), "vhv": edakit.IntStr(dy - 1)})
	case 2:
		n := dx + dy
		return edakit.Numeric(intro+zdef+" An L-shaped route has exactly one bend. How many distinct L- and Z-shaped routes together can a pattern router evaluate for this net?", edakit.IntStr(n),
			edakit.Sprintf("There are 2 L shapes and (%d - 1) + (%d - 1) = %d Z shapes, %d in all. Pattern routing evaluates a handful of candidates, not every path.", dx, dy, n-2, n),
			map[string]string{"l": "2", "z": edakit.IntStr(n - 2)})
	default:
		if dx > 6 {
			dx = edakit.Int(2, 6)
		}
		if dy > 6 {
			dy = edakit.Int(2, 6)
		}
		b = pt{x1 + dx, y1 + dy}
		intro = edakit.Sprintf("A two-pin net joins gcells %s and %s on the global-routing grid. ", ptText(a), ptText(b))
		all := binom(dx+dy, dx)
		ans := all - (dx + dy)
		return edakit.Numeric(intro+zdef+" A monotone route never moves away from the target. How many monotone routes have three or more bends, so that neither an L nor a Z pattern produces them?", edakit.IntStr(ans),
			edakit.Sprintf("Monotone routes are lattice paths: C(%d, %d) = %d. Of these, 2 are L shapes and %d are Z shapes, leaving %d. Monotonic routing by dynamic programming searches all of them at linear cost in the box area.", dx+dy, dx, all, dx+dy-2, ans),
			map[string]string{"monotone": edakit.IntStr(all), "patterns": edakit.IntStr(dx + dy)})
	}
}

// edgeMap holds edge costs in a (dx+1) by (dy+1) box of gcells. h[y][x] joins (x, y) and
// (x+1, y); v[x][y] joins (x, y) and (x, y+1).
type edgeMap struct {
	dx, dy int
	h, v   [][]int
}

func randEdgeMap(dx, dy int) edgeMap {
	m := edgeMap{dx: dx, dy: dy}
	m.h = make([][]int, dy+1)
	for y := range m.h {
		m.h[y] = make([]int, dx)
		for x := range m.h[y] {
			m.h[y][x] = edakit.Int(1, 9)
		}
	}
	m.v = make([][]int, dx+1)
	for x := range m.v {
		m.v[x] = make([]int, dy)
		for y := range m.v[x] {
			m.v[x][y] = edakit.Int(1, 9)
		}
	}
	return m
}

func (m edgeMap) text(full bool) string {
	var rows, cols []string
	for y := 0; y <= m.dy; y++ {
		if full || y == 0 || y == m.dy {
			rows = append(rows, edakit.Sprintf("row %d: %s", y, intsText(m.h[y])))
		}
	}
	for x := 0; x <= m.dx; x++ {
		if full || x == 0 || x == m.dx {
			cols = append(cols, edakit.Sprintf("column %d: %s", x, intsText(m.v[x])))
		}
	}
	return edakit.Sprintf("The pins are at gcells (0, 0) and (%d, %d). Horizontal edge costs, listed west to east, are %s. Vertical edge costs, listed south to north, are %s. ", m.dx, m.dy, strings.Join(rows, "; "), strings.Join(cols, "; "))
}

func (m edgeMap) hrun(y, x0, x1 int) int {
	s := 0
	for x := x0; x < x1; x++ {
		s += m.h[y][x]
	}
	return s
}

func (m edgeMap) vrun(x, y0, y1 int) int {
	s := 0
	for y := y0; y < y1; y++ {
		s += m.v[x][y]
	}
	return s
}

// patterns returns the two L costs and every Z cost.
func (m edgeMap) patterns() (l []int, z []int) {
	l = []int{m.hrun(0, 0, m.dx) + m.vrun(m.dx, 0, m.dy), m.vrun(0, 0, m.dy) + m.hrun(m.dy, 0, m.dx)}
	for j := 1; j < m.dx; j++ {
		z = append(z, m.hrun(0, 0, j)+m.vrun(j, 0, m.dy)+m.hrun(m.dy, j, m.dx))
	}
	for i := 1; i < m.dy; i++ {
		z = append(z, m.vrun(0, 0, i)+m.hrun(i, 0, m.dx)+m.vrun(m.dx, i, m.dy))
	}
	return
}

func (m edgeMap) monotone() int {
	best := make([][]int, m.dx+1)
	for x := range best {
		best[x] = make([]int, m.dy+1)
	}
	for x := 0; x <= m.dx; x++ {
		for y := 0; y <= m.dy; y++ {
			if x == 0 && y == 0 {
				continue
			}
			b := math.MaxInt32
			if x > 0 {
				b = minI(b, best[x-1][y]+m.h[y][x-1])
			}
			if y > 0 {
				b = minI(b, best[x][y-1]+m.v[x][y-1])
			}
			best[x][y] = b
		}
	}
	return best[m.dx][m.dy]
}

func minSlice(xs []int) int {
	b := math.MaxInt32
	for _, x := range xs {
		b = minI(b, x)
	}
	return b
}

// costEval: the cheaper L (KP 1), the best L or Z (KP 2), the best pattern with a bend cost
// (KP 3).
func costEval(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		m := randEdgeMap(edakit.Scaled(ctx, 2, 5), edakit.Scaled(ctx, 1, 4))
		l, _ := m.patterns()
		ans := minI(l[0], l[1])
		q := "A pattern router prices each route as the sum of its gcell-edge costs, which grow with congestion. " + m.text(false) + "What is the cost of the cheaper of the two L-shaped routes?"
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Horizontal first along row 0, then up column %d: %d. Vertical first up column 0, then along row %d: %d. The router keeps %d.", m.dx, l[0], m.dy, l[1], ans),
			map[string]string{"l_hv": edakit.IntStr(l[0]), "l_vh": edakit.IntStr(l[1])})
	case 2:
		m := randEdgeMap(edakit.Int(2, 3), edakit.Int(1, 2))
		l, z := m.patterns()
		ans := minI(minSlice(l), minSlice(z))
		q := "A pattern router prices each route as the sum of its gcell-edge costs. " + m.text(true) + "Over both L shapes and every Z shape (two bends, jog strictly between the pins), what is the cheapest pattern cost?"
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("The L shapes cost %s and the Z shapes cost %s. The best is %d.", intsText(l), intsText(z), ans),
			map[string]string{"best_l": edakit.IntStr(minSlice(l)), "best": edakit.IntStr(ans)})
	default:
		m := randEdgeMap(edakit.Int(2, 3), edakit.Int(1, 2))
		bc := edakit.Int(2, 6)
		l, z := m.patterns()
		bl, bz := minSlice(l)+bc, minSlice(z)+2*bc
		ans := minI(bl, bz)
		q := "A pattern router prices each route as the sum of its gcell-edge costs plus " + edakit.IntStr(bc) + " per bend, a proxy for the via each bend will need. " + m.text(true) + "Over both L shapes (one bend) and every Z shape (two bends, jog strictly between the pins), what is the cheapest total cost?"
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("The best L costs %d + %d = %d and the best Z costs %d + 2(%d) = %d. The router keeps %d. A bend penalty pushes the choice toward fewer vias.", minSlice(l), bc, bl, minSlice(z), bc, bz, ans),
			map[string]string{"best_l": edakit.IntStr(bl), "best_z": edakit.IntStr(bz)})
	}
}

// monotonic: the min-cost monotone route by dynamic programming (KP 1 small, KP 2 larger), and
// how much it beats the best pattern (KP 3).
func monotonic(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		dx, dy := edakit.Int(2, 3), edakit.Int(1, 2)
		if ctx.KPOrdinal == 2 || (ctx.KPOrdinal == 0 && ctx.Difficulty > 0.6) {
			dx, dy = edakit.Int(3, 4), edakit.Int(2, 3)
		}
		m := randEdgeMap(dx, dy)
		ans := m.monotone()
		q := "Monotonic routing finds the cheapest route that never moves west or south, by dynamic programming over the bounding box: the best cost to a gcell is the smaller of (best to its west neighbour + that edge) and (best to its south neighbour + that edge). " + m.text(true) + "What is the cost of the best monotone route?"
		return edakit.Numeric(q, edakit.IntStr(ans),
			edakit.Sprintf("Fill the table row by row from (0, 0). The entry at (%d, %d) is %d. Each gcell is visited once, so the work is linear in the box area.", dx, dy, ans),
			map[string]string{"best": edakit.IntStr(ans)})
	default:
		var m edgeMap
		var gap int
		for try := 0; try < generator.MaxDrawAttempts; try++ {
			m = randEdgeMap(edakit.Int(3, 4), edakit.Int(2, 3))
			l, z := m.patterns()
			gap = minI(minSlice(l), minSlice(z)) - m.monotone()
			if gap > 0 || edakit.Int(0, 3) == 0 {
				break
			}
		}
		q := "A global router first tries L and Z patterns (at most two bends) and then monotonic routing (any route that never moves west or south). Costs are sums of gcell-edge costs. " + m.text(true) + "By how much does the best monotone route undercut the best L or Z pattern?"
		l, z := m.patterns()
		bp := minI(minSlice(l), minSlice(z))
		return edakit.Numeric(q, edakit.IntStr(gap),
			edakit.Sprintf("The best pattern costs %d and the dynamic program over all monotone routes finds %d, a saving of %d. Monotone routes include every pattern, so the saving is never negative.", bp, m.monotone(), gap),
			map[string]string{"pattern": edakit.IntStr(bp), "monotone": edakit.IntStr(m.monotone())})
	}
}

// randPins draws n pins with distinct x and distinct y coordinates in [0, span]; it needs
// n <= span + 1, and it picks each coordinate from the values still unused, so it never redraws.
func randPins(n, span int) []pt {
	var xs, ys []int
	for v := 0; v <= span; v++ {
		xs = append(xs, v)
		ys = append(ys, v)
	}
	xs, ys = edakit.Shuffle(xs), edakit.Shuffle(ys)
	ps := make([]pt, n)
	for i := range ps {
		ps[i] = pt{xs[i], ys[i]}
	}
	return ps
}

func mstLen(ps []pt) int {
	n := len(ps)
	in := make([]bool, n)
	d := make([]int, n)
	for i := range d {
		d[i] = math.MaxInt32
	}
	d[0] = 0
	tot := 0
	for k := 0; k < n; k++ {
		u := -1
		for i := 0; i < n; i++ {
			if !in[i] && (u < 0 || d[i] < d[u]) {
				u = i
			}
		}
		in[u] = true
		tot += d[u]
		for i := 0; i < n; i++ {
			if !in[i] && manh(ps[u], ps[i]) < d[i] {
				d[i] = manh(ps[u], ps[i])
			}
		}
	}
	return tot
}

func hpwl(ps []pt) int {
	x0, x1, y0, y1 := ps[0].x, ps[0].x, ps[0].y, ps[0].y
	for _, p := range ps {
		x0, x1, y0, y1 = minI(x0, p.x), maxI(x1, p.x), minI(y0, p.y), maxI(y1, p.y)
	}
	return x1 - x0 + y1 - y0
}

func pinsText(ps []pt) string {
	var parts []string
	for _, p := range ps {
		parts = append(parts, ptText(p))
	}
	return strings.Join(parts, ", ")
}

// netDecomposition: rectilinear MST length (KP 1), what a Steiner point saves on a 3-pin net
// (KP 2), how much a driver-star decomposition wastes over the MST (KP 3).
func netDecomposition(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		ps := randPins(edakit.Scaled(ctx, 4, 6), 20)
		L := mstLen(ps)
		q := edakit.Sprintf("A global router decomposes a net with pins at %s into two-pin connections along a rectilinear minimum spanning tree. What is the total Manhattan length of that tree?", pinsText(ps))
		return edakit.Numeric(q, edakit.IntStr(L),
			edakit.Sprintf("Prim's algorithm on Manhattan distances adds the nearest outside pin each step. The %d tree edges sum to %d.", len(ps)-1, L),
			map[string]string{"mst": edakit.IntStr(L)})
	case 2:
		ps := []pt{{0, 0}, {10, 2}, {4, 9}}
		for try := 0; try < generator.MaxDrawAttempts; try++ {
			if c := randPins(3, 20); mstLen(c) > hpwl(c) {
				ps = c
				break
			}
		}
		xs := []int{ps[0].x, ps[1].x, ps[2].x}
		ys := []int{ps[0].y, ps[1].y, ps[2].y}
		sort.Ints(xs)
		sort.Ints(ys)
		save := mstLen(ps) - hpwl(ps)
		q := edakit.Sprintf("A three-pin net has pins at %s. Decomposing it along the rectilinear MST gives one total length; decomposing it at the optimal Steiner point (the median x and median y of the pins) into three two-pin connections gives another. How much wirelength does the Steiner decomposition save?", pinsText(ps))
		return edakit.Numeric(q, edakit.IntStr(save),
			edakit.Sprintf("The MST has length %d. The Steiner point (%d, %d) gives a tree whose length equals the half-perimeter of the bounding box, %d, which is optimal for three pins. The saving is %d.", mstLen(ps), xs[1], ys[1], hpwl(ps), save),
			map[string]string{"mst": edakit.IntStr(mstLen(ps)), "rsmt": edakit.IntStr(hpwl(ps))})
	default:
		ps := randPins(edakit.Scaled(ctx, 4, 6), 20)
		star := 0
		for _, p := range ps[1:] {
			star += manh(ps[0], p)
		}
		L := mstLen(ps)
		q := edakit.Sprintf("A net's driver is at %s and its sinks at %s. A naive decomposition routes one two-pin connection from the driver to each sink. How much more total Manhattan length does that star use than a rectilinear MST decomposition?", ptText(ps[0]), pinsText(ps[1:]))
		return edakit.Numeric(q, edakit.IntStr(star-L),
			edakit.Sprintf("The star sums driver-to-sink distances: %d. The MST, which lets sinks connect to each other, is %d. The difference is %d of demand the global router would place on the grid for nothing, unless sink delay demands short source paths.", star, L, star-L),
			map[string]string{"star": edakit.IntStr(star), "mst": edakit.IntStr(L)})
	}
}

// fluteSequences asks how many position sequences index FLUTE's table for a degree.
func fluteSequences(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Int(3, 8)
	f := 1
	for i := 2; i <= n; i++ {
		f *= i
	}
	q := edakit.Sprintf("FLUTE sorts a net's pins by x, then records the rank of each pin's y coordinate. Nets with the same rank sequence (position sequence) share one entry of its lookup table. For a net of degree %d with distinct coordinates, how many position sequences are possible?", n)
	ans := edakit.IntStr(f)
	d := []string{edakit.IntStr(1 << n), edakit.IntStr(n * n), edakit.IntStr(f / n), edakit.IntStr(binom(2*n, n))}
	return edakit.Choice(q, ans,
		edakit.Sprintf("After sorting by x, the y ranks form a permutation of 1 to %d, and there are %d! = %d of them. This is why tables are precomputed only for small degrees and larger nets are broken up.", n, n, f), d...)
}

// fluteEval evaluates FLUTE wirelength vectors against a net's Hanan-grid gaps.
func fluteEval(ctx generator.GeneratorContext) generator.Problem {
	h := []int{edakit.Int(1, 9), edakit.Int(1, 9), edakit.Int(1, 9)}
	v := []int{edakit.Int(1, 9), edakit.Int(1, 9), edakit.Int(1, 9)}
	powv := [][]int{{1, 1, 1, 1, 2, 1}, {1, 2, 1, 1, 1, 1}}
	val := func(c []int) int {
		return c[0]*h[0] + c[1]*h[1] + c[2]*h[2] + c[3]*v[0] + c[4]*v[1] + c[5]*v[2]
	}
	a, b := val(powv[0]), val(powv[1])
	ans := minI(a, b)
	hp := h[0] + h[1] + h[2] + v[0] + v[1] + v[2]
	q := edakit.Sprintf("FLUTE writes the wirelength of a degree-4 net as a weighted sum of the gaps between adjacent Hanan-grid lines: horizontal gaps $h_1, h_2, h_3$ and vertical gaps $v_1, v_2, v_3$. For this net's position sequence the table holds two potentially optimal wirelength vectors, with coefficients (1, 1, 1 | 1, 2, 1) and (1, 2, 1 | 1, 1, 1) on (h | v). The net's gaps are h = (%s) and v = (%s). What wirelength does FLUTE return?", intsText(h), intsText(v))
	return edakit.Choice(q, edakit.IntStr(ans),
		edakit.Sprintf("The first vector gives %d and the second %d. FLUTE evaluates every vector stored for the sequence and returns the minimum, %d, which is the half-perimeter %d plus min(h2, v2) = %d.", a, b, ans, hp, minI(h[1], v[1])),
		edakit.Distinct(ans, 4, maxI(a, b), hp, hp+h[1]+v[1], a+b-hp)...)
}
