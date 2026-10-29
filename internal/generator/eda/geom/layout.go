package geom

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// spacingCheck: Euclidean violations among many shapes (KP 1), pairs the square metric flags but
// the Euclidean metric passes (KP 2), a width- and run-length-dependent spacing table (KP 3).
func spacingCheck(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 7)
	k := kp(ctx)
	switch k {
	case 1, 2:
		s := edakit.Int(3, 6)
		rs := randRects(n, 36, 2, 8)
		eu, sqOnly := 0, 0
		for i := range rs {
			for j := i + 1; j < n; j++ {
				gx, gy := gaps(rs[i], rs[j])
				if gx == 0 && gy == 0 {
					continue
				}
				e := gx*gx+gy*gy < s*s
				sq := imax(gx, gy) < s
				if e {
					eu++
				}
				if sq && !e {
					sqOnly++
				}
			}
		}
		if k == 1 {
			q := edakit.Sprintf("Metal shapes are %s. The spacing rule is %d, measured Euclidean between separated shapes (touching or overlapping shapes merge and are not checked). How many pairs violate it?", listRects("M", rs), s)
			return edakit.Numeric(q, edakit.IntStr(eu),
				edakit.Sprintf("For each separated pair, compare gx^2 + gy^2 with %d. %d pairs fall short. A sweep or spatial index finds the candidates without testing all %d pairs.", s*s, eu, n*(n-1)/2),
				map[string]string{"violations": edakit.IntStr(eu)})
		}
		q := edakit.Sprintf("Metal shapes are %s, with spacing rule %d between separated shapes. How many pairs does the square metric (max of x-gap and y-gap) flag that the Euclidean metric passes?", listRects("M", rs), s)
		return edakit.Numeric(q, edakit.IntStr(sqOnly),
			edakit.Sprintf("Only corner-to-corner pairs can differ: max(gx, gy) < %d while sqrt(gx^2 + gy^2) >= %d. %d pairs do. The square metric is cheaper and more conservative; the Euclidean one matches the foundry rule.", s, s, sqOnly),
			map[string]string{"square_only": edakit.IntStr(sqOnly)})
	default:
		prl := []int{0, 200, 500}
		wid := []int{0, 150, 400}
		tab := [][]int{{70, 70, 70}, {70, 90, 110}, {70, 120, 160}}
		w1, w2 := 10*edakit.Int(5, 60), 10*edakit.Int(5, 60)
		for w1 == 150 || w1 == 400 || w2 == 150 || w2 == 400 {
			w1, w2 = 10*edakit.Int(5, 60), 10*edakit.Int(5, 60)
		}
		// Redraw the whole pair, not only x2: a short first wire inside a long second wire
		// overlaps by exactly its own length (200 or 500) for every x2, so redrawing x2 alone
		// never terminates.
		var x1, l1, x2, l2, run int
		for attempt := 0; attempt < generator.MaxDrawAttempts; attempt++ {
			x1, l1 = 10*edakit.Int(0, 40), 10*edakit.Int(20, 90)
			x2, l2 = 10*edakit.Int(0, 40), 10*edakit.Int(20, 90)
			run = overlap1(x1, x1+l1, x2, x2+l2)
			if run > 0 && run != 200 && run != 500 {
				break
			}
		}
		if run <= 0 || run == 200 || run == 500 {
			x1, l1, x2, l2 = 0, 300, 110, 300
			run = overlap1(x1, x1+l1, x2, x2+l2)
		}
		wmax := imax(w1, w2)
		row, col := 0, 0
		for i, t := range wid {
			if wmax > t {
				row = i
			}
		}
		for i, t := range prl {
			if run > t {
				col = i
			}
		}
		req := tab[row][col]
		q := edakit.Sprintf("A spacing table gives the required spacing in nm. Rows by the wider wire's width (above 0, above 150, above 400 nm) and columns by parallel run length (above 0, above 200, above 500 nm): row 1 is 70, 70, 70; row 2 is 70, 90, 110; row 3 is 70, 120, 160. Use the last row and the last column whose threshold is exceeded. Two horizontal wires of widths %d nm and %d nm run over x from %d to %d and from %d to %d. What spacing is required, in nm?", w1, w2, x1, x1+l1, x2, x2+l2)
		return edakit.Numeric(q, edakit.IntStr(req),
			edakit.Sprintf("The wider wire is %d nm (row %d) and the parallel run length is the x-overlap, %d nm (column %d), so the table gives %d nm. Wide, long-running neighbours need more space because of coupling and lithography effects.", wmax, row+1, run, col+1, req),
			map[string]string{"width": edakit.IntStr(wmax), "prl": edakit.IntStr(run)})
	}
}

// widthCheck: area lost to an opening (KP 1), opening an L with a thin arm (KP 2), the overlap
// created by growing two shapes by half the spacing (KP 3).
func widthCheck(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		d := edakit.Int(2, 5)
		k := edakit.Scaled(ctx, 3, 6)
		var rs []rect
		lost := 0
		for i := 0; i < k; i++ {
			w, h := edakit.Int(1, 5*d), edakit.Int(1, 5*d)
			for w == 2*d || h == 2*d {
				w, h = edakit.Int(1, 5*d), edakit.Int(1, 5*d)
			}
			x := 40 * i
			rs = append(rs, rect{x, 0, x + w, h})
			if imin(w, h) < 2*d {
				lost += w * h
			}
		}
		q := edakit.Sprintf("Separate shapes %s are shrunk by %d and then grown by %d (an opening). How much area is lost?", listRects("S", rs), d, d)
		return edakit.Numeric(q, edakit.IntStr(lost),
			edakit.Sprintf("A rectangle survives the shrink only if both sides exceed 2d = %d; then the grow restores it exactly. Rectangles narrower than %d vanish and stay gone: %d area lost. The lost shapes are the width violations for a rule of %d.", 2*d, 2*d, lost, 2*d),
			map[string]string{"lost": edakit.IntStr(lost)})
	case 2:
		d := edakit.Int(2, 5)
		T := rect{0, 0, edakit.Int(2*d+2, 6*d), edakit.Int(2*d+4, 6*d)}
		armW := edakit.Int(1, 4*d)
		for armW == 2*d || armW >= T.h() {
			armW = edakit.Int(1, 4*d)
		}
		y := edakit.Int(0, T.h()-armW)
		arm := rect{T.x2, y, T.x2 + edakit.Int(d, 6*d), y + armW}
		keep := T.area()
		if armW > 2*d {
			keep += arm.area()
		}
		q := edakit.Sprintf("A polygon is the union of %s and the arm %s. It is shrunk by %d and grown by %d. What area remains?", T, arm, d, d)
		return edakit.Numeric(q, edakit.IntStr(keep),
			edakit.Sprintf("The opening keeps every point reachable by a 2d-by-2d square (%d by %d) lying inside the shape. The arm is %d wide, %s 2d, so %s. Remaining area %d.", 2*d, 2*d, armW, map[bool]string{true: "above", false: "below"}[armW > 2*d],
				map[bool]string{true: "it survives", false: "it disappears while the body survives"}[armW > 2*d], keep),
			map[string]string{"remaining": edakit.IntStr(keep)})
	default:
		s := 2 * edakit.Int(5, 20)
		g := edakit.Int(1, s-1)
		h1, h2 := edakit.Int(20, 80), edakit.Int(20, 80)
		a := rect{0, 0, edakit.Int(20, 60), h1}
		y := edakit.Int(-h2+5, h1-5)
		b := rect{a.x2 + g, y, a.x2 + g + edakit.Int(20, 60), y + h2}
		L := overlap1(a.y1, a.y2, b.y1, b.y2)
		ov := (s - g) * (L + s)
		q := edakit.Sprintf("To find spacing violations, shapes A = %s and B = %s (nm) are each grown by %d nm, half the %d nm rule, and then ANDed. What is the area of the AND, in square nm?", a, b, s/2, s)
		return edakit.Numeric(q, edakit.IntStr(ov),
			edakit.Sprintf("Grown, the x-ranges overlap by s - g = %d - %d = %d and the y-ranges by L + s = %d + %d = %d. The AND is a %d by %d marker of area %d; any non-empty AND flags a gap below %d nm.", s, g, s-g, L, s, L+s, s-g, L+s, ov, s),
			map[string]string{"x_overlap": edakit.IntStr(s - g), "y_overlap": edakit.IntStr(L + s)})
	}
}

// density: window density (KP 1), fill squares needed (KP 2), sliding-window count (KP 3).
func density(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		W := 10 * edakit.Int(5, 10)
		win := rect{0, 0, W, W}
		k := edakit.Scaled(ctx, 3, 5)
		rs := randRects(k, W+10, 5, W/2)
		var cl []rect
		for _, r := range rs {
			c := rect{imax(r.x1, 0), imax(r.y1, 0), imin(r.x2, W), imin(r.y2, W)}
			if c.x2 > c.x1 && c.y2 > c.y1 {
				cl = append(cl, c)
			}
		}
		m := 0
		if len(cl) > 0 {
			m = unionArea(cl)
		}
		pct := 100 * float64(m) / float64(win.area())
		ans := edakit.Num(pct, 1)
		q := edakit.Sprintf("A density window is %s (um). Metal shapes %s may overlap one another and may extend past the window. What is the metal density in the window, in percent? Round to 1 decimal place.", win, listRects("M", rs))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Clip each shape to the window, take the union area (%d square um, overlaps counted once) and divide by %d: %s percent.", m, win.area(), ans),
			map[string]string{"metal": edakit.IntStr(m)})
	case 2:
		W := 10 * edakit.Int(5, 20)
		D := 5 * edakit.Int(4, 8)
		cur := edakit.Int(5, D-3)
		metal := W * W * cur / 100
		f := edakit.Int(1, 4)
		need := W*W*D/100 - metal
		cnt := (need + f*f - 1) / (f * f)
		q := edakit.Sprintf("A %d by %d um window holds %d square um of metal, and the rule requires at least %d percent density. Fill is inserted as %d by %d um squares that never overlap existing metal. How many fill squares are needed?", W, W, metal, D, f, f)
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("The target is %d square um, so %d square um must be added: ceil(%d / %d) = %d squares.", W*W*D/100, need, need, f*f, cnt),
			map[string]string{"need": edakit.IntStr(need)})
	default:
		w := 10 * edakit.Int(2, 10)
		st := edakit.Pick([]int{w / 2, w})
		m := edakit.Int(3, 30)
		L := w + m*st
		cnt := (m + 1) * (m + 1)
		q := edakit.Sprintf("A density check slides a %d by %d um window over a %d by %d um die in steps of %d um in x and in y, keeping the window inside the die. How many window positions are checked?", w, w, L, L, st)
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("Per axis there are (%d - %d) / %d + 1 = %d positions, so %d windows. Each needs a clipped area, which is why density engines precompute per-tile areas and sum them.", L, w, st, m+1, cnt),
			map[string]string{"per_axis": edakit.IntStr(m + 1)})
	}
}

// hierarchy: flat shape count (KP 1), flat-to-hierarchical ratio (KP 2), interacting instance
// pairs in an array (KP 3).
func hierarchy(ctx generator.GeneratorContext) generator.Problem {
	sC, sB, sA, sT := edakit.Int(10, 60), edakit.Int(5, 40), edakit.Int(5, 40), edakit.Int(5, 50)
	nBC, nAB, nAC := edakit.Int(2, 12), edakit.Int(2, 10), edakit.Int(0, 6)
	nTA, nTB := edakit.Int(10, 200), edakit.Int(0, 50)
	flatB := sB + nBC*sC
	flatA := sA + nAB*flatB + nAC*sC
	flat := sT + nTA*flatA + nTB*flatB
	desc := edakit.Sprintf("Cell C draws %d shapes. Cell B draws %d shapes and places %d instances of C. Cell A draws %d shapes and places %d instances of B and %d of C. The top cell draws %d shapes and places %d instances of A and %d of B.", sC, sB, nBC, sA, nAB, nAC, sT, nTA, nTB)
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(desc+" How many shapes does the flattened layout contain?", edakit.IntStr(flat),
			edakit.Sprintf("Flatten bottom-up: B holds %d + %d x %d = %d, A holds %d + %d x %d + %d x %d = %d, top holds %d + %d x %d + %d x %d = %d.", sB, nBC, sC, flatB, sA, nAB, flatB, nAC, sC, flatA, sT, nTA, flatA, nTB, flatB, flat),
			map[string]string{"flat_b": edakit.IntStr(flatB), "flat_a": edakit.IntStr(flatA), "flat": edakit.IntStr(flat)})
	case 2:
		hier := sC + sB + sA + sT
		r := float64(flat) / float64(hier)
		ans := edakit.Num(r, 1)
		return edakit.Numeric(desc+" A hierarchical checker examines each cell's own shapes once. What is the ratio of flat shapes to shapes stored once per cell? Round to 1 decimal place.", ans,
			edakit.Sprintf("Flat: %d shapes. Stored once per cell: %d + %d + %d + %d = %d. Ratio %s. The saving is real only where instances do not interact with their surroundings; interaction regions must still be checked in context.", flat, sC, sB, sA, sT, hier, ans),
			map[string]string{"flat": edakit.IntStr(flat), "hier": edakit.IntStr(hier)})
	default:
		k := edakit.Int(4, 12)
		w := 10 * edakit.Int(5, 20)
		p := w + 10*edakit.Int(1, 8)
		s := 10 * edakit.Int(1, 15)
		pairs := 0
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				if (j-i)*p-w < s {
					pairs++
				}
			}
		}
		q := edakit.Sprintf("A row of %d instances of a cell %d nm wide is placed at a pitch of %d nm. Two instances must be checked together when the gap between their bounding boxes is below the %d nm interaction distance. How many instance pairs must be checked together?", k, w, p, s)
		return edakit.Numeric(q, edakit.IntStr(pairs),
			edakit.Sprintf("Instances i and j are (j - i)(%d) - %d nm apart. Counting pairs below %d nm gives %d; every other pair can reuse the cell's own result.", p, w, s, pairs),
			map[string]string{"pairs": edakit.IntStr(pairs)})
	}
}

// parallelRun: parallel run length of two wires (KP 1), crossover overlap area between layers
// (KP 2), total coupled run length within a spacing window (KP 3).
func parallelRun(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		x1, l1 := edakit.Int(0, 200), edakit.Int(40, 300)
		x2, l2 := edakit.Int(0, 200), edakit.Int(40, 300)
		y, w, g := edakit.Int(0, 50), edakit.Int(4, 20), edakit.Int(4, 30)
		a := rect{x1, y, x1 + l1, y + w}
		b := rect{x2, y + w + g, x2 + l2, y + 2*w + g}
		run := imax(0, overlap1(a.x1, a.x2, b.x1, b.x2))
		q := edakit.Sprintf("Two parallel horizontal wires are A = %s and B = %s (nm). What is their parallel run length, in nm?", a, b)
		return edakit.Numeric(q, edakit.IntStr(run),
			edakit.Sprintf("The run length is the overlap of the x-projections: min(%d, %d) - max(%d, %d) = %d (zero if negative). Coupling capacitance grows with it.", a.x2, b.x2, a.x1, b.x1, run),
			map[string]string{"prl": edakit.IntStr(run)})
	case 2:
		var m1, m2 []rect
		for i := 0; i < edakit.Int(2, 3); i++ {
			y := 20*i + edakit.Int(0, 6)
			x := edakit.Int(0, 20)
			m1 = append(m1, rect{x, y, x + edakit.Int(30, 70), y + edakit.Int(3, 8)})
		}
		for i := 0; i < edakit.Int(2, 3); i++ {
			x := 25*i + edakit.Int(0, 10)
			y := edakit.Int(0, 15)
			m2 = append(m2, rect{x, y, x + edakit.Int(3, 8), y + edakit.Int(25, 50)})
		}
		tot := 0
		for _, a := range m1 {
			for _, b := range m2 {
				tot += interArea(a, b)
			}
		}
		q := edakit.Sprintf("Metal-1 wires are %s and metal-2 wires are %s. A plate-capacitance extractor needs the total M1-M2 overlap area. What is it?", listRects("A", m1), listRects("B", m2))
		return edakit.Numeric(q, edakit.IntStr(tot),
			edakit.Sprintf("AND the two layers and sum the intersection areas over all crossing pairs: %d.", tot),
			map[string]string{"overlap": edakit.IntStr(tot)})
	default:
		S := edakit.Int(20, 40)
		v := rect{0, 100, edakit.Int(150, 300), 100 + edakit.Int(8, 16)}
		var nb []rect
		tot := 0
		for i := 0; i < 4; i++ {
			gap := edakit.Int(6, 2*S)
			x := edakit.Int(-50, v.x2)
			l := edakit.Int(30, 150)
			w := edakit.Int(6, 14)
			var r rect
			if i%2 == 0 {
				r = rect{x, v.y2 + gap + 30*(i/2), x + l, v.y2 + gap + 30*(i/2) + w}
			} else {
				r = rect{x, v.y1 - gap - w - 30*(i/2), x + l, v.y1 - gap - 30*(i/2)}
			}
			nb = append(nb, r)
			_, gy := gaps(v, r)
			run := overlap1(v.x1, v.x2, r.x1, r.x2)
			if gy <= S && run > 0 {
				tot += run
			}
		}
		q := edakit.Sprintf("A victim wire %s has neighbours %s (nm). The coupling model counts run length only for neighbours within %d nm edge-to-edge. What total parallel run length does it count, in nm?", v, listRects("N", nb), S)
		return edakit.Numeric(q, edakit.IntStr(tot),
			edakit.Sprintf("For each neighbour with vertical gap at most %d nm, add the overlap of its x-range with the victim's. The total is %d nm.", S, tot),
			map[string]string{"prl": edakit.IntStr(tot)})
	}
}

type dsu []int

func newDSU(n int) dsu {
	d := make(dsu, n)
	for i := range d {
		d[i] = i
	}
	return d
}

func (d dsu) find(x int) int {
	for d[x] != x {
		d[x] = d[d[x]]
		x = d[x]
	}
	return x
}

func (d dsu) union(a, b int) bool {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return false
	}
	d[ra] = rb
	return true
}

// connectivity: components on one layer (KP 1), nets across two layers and vias (KP 2),
// successful union operations (KP 3).
func connectivity(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 8)
		rs := randRects(n, 30, 2, 9)
		d := newDSU(n)
		comps := n
		for i := range rs {
			for j := i + 1; j < n; j++ {
				if edgeConnected(rs[i], rs[j]) && d.union(i, j) {
					comps--
				}
			}
		}
		q := edakit.Sprintf("Shapes on one metal layer are %s. Shapes that overlap or share a boundary segment are electrically connected; touching only at a corner does not connect. How many separate nets (connected components) are there?", listRects("M", rs))
		return edakit.Numeric(q, edakit.IntStr(comps),
			edakit.Sprintf("Union every connected pair with a disjoint-set structure; %d sets remain.", comps),
			map[string]string{"nets": edakit.IntStr(comps)})
	case 2:
		var m1, m2, vias []rect
		for i := 0; i < edakit.Int(3, 4); i++ {
			y := 10*i + edakit.Int(0, 3)
			x := edakit.Int(0, 15)
			m1 = append(m1, rect{x, y, x + edakit.Int(10, 25), y + 3})
		}
		for i := 0; i < edakit.Int(3, 4); i++ {
			x := 10*i + edakit.Int(0, 3)
			y := edakit.Int(0, 15)
			m2 = append(m2, rect{x, y, x + 3, y + edakit.Int(10, 25)})
		}
		nv := edakit.Int(2, 4)
		for len(vias) < nv {
			a, b := m1[edakit.Int(0, len(m1)-1)], m2[edakit.Int(0, len(m2)-1)]
			x0, x1 := imax(a.x1, b.x1), imin(a.x2, b.x2)
			y0, y1 := imax(a.y1, b.y1), imin(a.y2, b.y2)
			if x1-x0 >= 1 && y1-y0 >= 1 {
				vias = append(vias, rect{x0, y0, x0 + 1, y0 + 1})
			} else if edakit.Int(0, 3) == 0 {
				vx, vy := edakit.Int(0, 40), edakit.Int(0, 40)
				vias = append(vias, rect{vx, vy, vx + 1, vy + 1})
			}
		}
		all := append(append([]rect(nil), m1...), m2...)
		d := newDSU(len(all))
		for i := range m1 {
			for j := i + 1; j < len(m1); j++ {
				if edgeConnected(m1[i], m1[j]) {
					d.union(i, j)
				}
			}
		}
		for i := range m2 {
			for j := i + 1; j < len(m2); j++ {
				if edgeConnected(m2[i], m2[j]) {
					d.union(len(m1)+i, len(m1)+j)
				}
			}
		}
		for _, v := range vias {
			for i, a := range m1 {
				for j, b := range m2 {
					if interArea(v, a) > 0 && interArea(v, b) > 0 {
						d.union(i, len(m1)+j)
					}
				}
			}
		}
		roots := map[int]bool{}
		for i := range all {
			roots[d.find(i)] = true
		}
		q := edakit.Sprintf("Metal-1 shapes are %s; metal-2 shapes are %s; vias are %s. Same-layer shapes connect when they overlap or share an edge segment, and a via connects the M1 and M2 shapes it overlaps (vias touching nothing on one layer connect nothing). How many nets are there?", listRects("A", m1), listRects("B", m2), listRects("V", vias))
		return edakit.Numeric(q, edakit.IntStr(len(roots)),
			edakit.Sprintf("Union same-layer contacts, then union through each via that lands on both layers. %d components remain, each a net.", len(roots)),
			map[string]string{"nets": edakit.IntStr(len(roots))})
	default:
		n := edakit.Scaled(ctx, 6, 12)
		k := edakit.Int(n/2, n+3)
		d := newDSU(n)
		merges := 0
		var es []string
		for i := 0; i < k; i++ {
			a, b := edakit.Int(1, n), edakit.Int(1, n)
			for a == b {
				b = edakit.Int(1, n)
			}
			es = append(es, edakit.Sprintf("(%d, %d)", a, b))
			if d.union(a-1, b-1) {
				merges++
			}
		}
		q := edakit.Sprintf("An extractor numbers %d shapes 1 to %d and, scanning for touching shapes, calls union on the pairs %s, in that order. How many of these calls merge two different sets?", n, n, strings.Join(es, ", "))
		return edakit.Numeric(q, edakit.IntStr(merges),
			edakit.Sprintf("A call merges only when its two shapes are not yet connected. %d calls merge, leaving %d nets: nets = shapes - merges.", merges, n-merges),
			map[string]string{"merges": edakit.IntStr(merges), "nets": edakit.IntStr(n - merges)})
	}
}
