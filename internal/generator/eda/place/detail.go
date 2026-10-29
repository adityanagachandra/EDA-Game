package place

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// netWith returns the HPWL of a net whose movable pin is at p and whose other pins are others.
func netWith(p pt, others []pt) int { return hpwl(append([]pt{p}, others...)) }

// swap computes the HPWL change of exchanging two cells.
func swap(ctx generator.GeneratorContext) generator.Problem {
	pa := pt{edakit.Int(0, 50), edakit.Int(0, 50)}
	pb := pt{edakit.Int(0, 50), edakit.Int(0, 50)}
	for try := 0; pa == pb && try < generator.MaxDrawAttempts; try++ {
		pb = pt{edakit.Int(0, 50), edakit.Int(0, 50)}
	}
	if pa == pb {
		pb = pt{(pa.x + 17) % 51, pa.y}
	}
	na, nb := 1, 1
	if ctx.KPOrdinal != 1 {
		na, nb = edakit.Int(1, 2), edakit.Int(1, 2)
	}
	before, after := 0, 0
	var desc, parts []string
	add := func(owner string, p, q pt, cnt int) {
		for i := 0; i < cnt; i++ {
			o := randPts(edakit.Int(1, 2), 0, 50)
			b, a := netWith(p, o), netWith(q, o)
			before += b
			after += a
			desc = append(desc, edakit.Sprintf("a net joins %s with fixed pins at %s", owner, fmtPts(o)))
			parts = append(parts, edakit.Sprintf("%d to %d", b, a))
		}
	}
	add("a", pa, pb, na)
	add("b", pb, pa, nb)
	if ctx.KPOrdinal == 3 {
		o := randPts(1, 0, 50)
		h := hpwl([]pt{pa, pb, o[0]})
		before += h
		after += h
		desc = append(desc, edakit.Sprintf("a shared net joins a, b and a fixed pin at %s", fmtPt(o[0])))
		parts = append(parts, edakit.Sprintf("shared net %d to %d", h, h))
	}
	d := after - before
	return edakit.Numeric(edakit.Sprintf("Detailed placement considers swapping equal-width cells a at %s and b at %s (um, treat each cell as one pin). In the placement, %s. What is the change in total HPWL if they swap, new minus old, in um?", fmtPt(pa), fmtPt(pb), strings.Join(desc, "; ")), edakit.IntStr(d),
		edakit.Sprintf("Only nets touching a or b change. Before and after: %s. The change is %d - %d = %d um; a swap is accepted only if this is negative.", strings.Join(parts, ", "), after, before, d),
		map[string]string{"before": edakit.IntStr(before), "after": edakit.IntStr(after)})
}

// optimalRegion finds the median interval that minimises a cell's incident HPWL.
func optimalRegion(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 2, 4)
	type nb struct{ lx, hx, ly, hy int }
	nets := make([]nb, k)
	var desc []string
	var xv, yv []int
	for i := range nets {
		o := randPts(edakit.Int(1, 3), 0, 60)
		xs, ys := make([]int, len(o)), make([]int, len(o))
		for j, p := range o {
			xs[j], ys[j] = p.x, p.y
		}
		nets[i] = nb{minI(xs...), maxI(xs...), minI(ys...), maxI(ys...)}
		xv = append(xv, nets[i].lx, nets[i].hx)
		yv = append(yv, nets[i].ly, nets[i].hy)
		if ctx.KPOrdinal == 3 {
			desc = append(desc, edakit.Sprintf("net n%d's other pins are at %s", i+1, fmtPts(o)))
		} else {
			desc = append(desc, edakit.Sprintf("net n%d's other pins are at x = %s", i+1, fmtInts(xs)))
		}
	}
	sx, sy := sortedCopy(xv), sortedCopy(yv)
	xlo, xhi := sx[k-1], sx[k]
	ylo, yhi := sy[k-1], sy[k]
	cost := func(x int) int {
		t := 0
		for _, n := range nets {
			t += maxI(n.hx, x) - minI(n.lx, x)
		}
		return t
	}
	cur := pt{edakit.Int(0, 60), edakit.Int(0, 60)}
	head := edakit.Sprintf("Cell c is on %d nets: %s. ", k, strings.Join(desc, "; "))
	switch ctx.KPOrdinal {
	case 2:
		return edakit.Numeric(head+edakit.Sprintf("Cell c is at x = %d. By how much does the total x span of its nets shrink if c moves to the best x position, in um (0 if it is already optimal)?", cur.x), edakit.IntStr(cost(cur.x)-cost(xlo)),
			edakit.Sprintf("The optimal x interval runs from the %s to the %s of the sorted net bounds %s: [%d, %d]. The span total is %d at x = %d and %d inside the interval.", ordinal(k), ordinal(k+1), fmtInts(sx), xlo, xhi, cost(cur.x), cur.x, cost(xlo)),
			map[string]string{"lo": edakit.IntStr(xlo), "hi": edakit.IntStr(xhi)})
	case 3:
		dist := func(v, lo, hi int) int {
			if v < lo {
				return lo - v
			}
			if v > hi {
				return v - hi
			}
			return 0
		}
		d := dist(cur.x, xlo, xhi) + dist(cur.y, ylo, yhi)
		return edakit.Numeric(head+edakit.Sprintf("Cell c is at %s. What is the Manhattan distance from c to the nearest point of its optimal region, in um?", fmtPt(cur)), edakit.IntStr(d),
			edakit.Sprintf("The optimal region is the rectangle x in [%d, %d], y in [%d, %d], from the medians of the net bounds (sorted x bounds %s, y bounds %s). From %s its nearest point is %d um away.", xlo, xhi, ylo, yhi, fmtInts(sx), fmtInts(sy), fmtPt(cur), d),
			map[string]string{"region": edakit.Sprintf("[%d,%d]x[%d,%d]", xlo, xhi, ylo, yhi)})
	}
	end, v := "lower", xlo
	if edakit.Int(0, 1) == 1 {
		end, v = "upper", xhi
	}
	return edakit.Numeric(head+edakit.Sprintf("Ignoring c, each net's other pins span an x interval. The x positions of c that minimise the total x span of its nets form an interval. What is its %s end?", end), edakit.IntStr(v),
		edakit.Sprintf("Collect each net's lower and upper x bound, %d values, and sort them: %s. Any x from the %s to the %s value is optimal: [%d, %d].", 2*k, fmtInts(sx), ordinal(k), ordinal(k+1), xlo, xhi),
		map[string]string{"lo": edakit.IntStr(xlo), "hi": edakit.IntStr(xhi)})
}

// slidingWindow reorders a window of cells in a row to minimise x wirelength.
func slidingWindow(ctx generator.GeneratorContext) generator.Problem {
	n := 3
	if ctx.Difficulty > 0.7 && ctx.KPOrdinal != 1 {
		n = 4
	}
	x0 := edakit.Int(10, 30)
	ws := make([]int, n)
	for i := range ws {
		ws[i] = 2 * edakit.Int(1, 4)
	}
	type net struct {
		cells []int
		fixed []int
	}
	var nets []net
	var desc []string
	for i := 0; i < n; i++ {
		f := []int{edakit.Int(0, 70)}
		if edakit.Int(0, 2) == 0 {
			f = append(f, edakit.Int(0, 70))
		}
		nets = append(nets, net{[]int{i}, f})
		desc = append(desc, edakit.Sprintf("%s to fixed x = %s", letter(i), fmtInts(f)))
	}
	i, j := edakit.Int(0, n-1), edakit.Int(0, n-1)
	if i != j {
		nets = append(nets, net{[]int{i, j}, nil})
		desc = append(desc, edakit.Sprintf("%s to %s", letter(i), letter(j)))
	}
	cost := func(order []int) int {
		center := make([]int, n)
		x := x0
		for _, c := range order {
			center[c] = x + ws[c]/2
			x += ws[c]
		}
		t := 0
		for _, nt := range nets {
			var xs []int
			for _, c := range nt.cells {
				xs = append(xs, center[c])
			}
			xs = append(xs, nt.fixed...)
			t += span(xs)
		}
		return t
	}
	id := make([]int, n)
	for k := range id {
		id[k] = k
	}
	best, bestOrd := cost(id), id
	for _, p := range permutations(n) {
		if c := cost(p); c < best {
			best, bestOrd = c, p
		}
	}
	var wd []string
	for k := range ws {
		wd = append(wd, edakit.Sprintf("%s %d", letter(k), ws[k]))
	}
	ordStr := func(o []int) string {
		s := make([]string, len(o))
		for k, c := range o {
			s[k] = letter(c)
		}
		return strings.Join(s, " ")
	}
	head := edakit.Sprintf("A sliding window covers %d adjacent cells in one row, packed from x = %d in the order %s, with widths %s um. Each cell's pin is at its centre. The nets touching the window are: %s. Only the x span of these nets changes when the window is reordered; the cells stay packed from x = %d. ", n, x0, ordStr(id), strings.Join(wd, ", "), strings.Join(desc, "; "), x0)
	switch ctx.KPOrdinal {
	case 1:
		return edakit.Numeric(head+"What is the total x span of these nets in the current order, in um?", edakit.IntStr(cost(id)),
			edakit.Sprintf("Place the centres, then sum each net's x span: %d um.", cost(id)), map[string]string{"current": edakit.IntStr(cost(id))})
	case 3:
		return edakit.Numeric(head+"The window tries every order of its cells and keeps the best. By how many um does the total x span drop (0 if the current order is best)?", edakit.IntStr(cost(id)-best),
			edakit.Sprintf("The current order costs %d um; the best order, %s, costs %d um, a gain of %d um.", cost(id), ordStr(bestOrd), best, cost(id)-best),
			map[string]string{"current": edakit.IntStr(cost(id)), "best": edakit.IntStr(best)})
	}
	return edakit.Numeric(head+"The window tries every order of its cells. What is the smallest total x span it can reach, in um?", edakit.IntStr(best),
		edakit.Sprintf("Of the %d orders, %s gives the minimum, %d um.", len(permutations(n)), ordStr(bestOrd), best),
		map[string]string{"best": edakit.IntStr(best)})
}

// ism solves a small independent-set matching instance.
func ism(ctx generator.GeneratorContext) generator.Problem {
	n := 3
	slots := make([]pt, n)
	for i := range slots {
		slots[i] = pt{edakit.Int(0, 60), edakit.Int(0, 60)}
	}
	others := make([][][]pt, n)
	var desc []string
	for i := 0; i < n; i++ {
		for k := 0; k < edakit.Int(1, 2); k++ {
			o := randPts(edakit.Int(1, 2), 0, 60)
			others[i] = append(others[i], o)
			desc = append(desc, edakit.Sprintf("%s has a net to fixed pins at %s", letter(i), fmtPts(o)))
		}
	}
	c := func(i, j int) int {
		t := 0
		for _, o := range others[i] {
			t += netWith(slots[j], o)
		}
		return t
	}
	id := []int{0, 1, 2}
	total := func(p []int) int {
		t := 0
		for i, j := range p {
			t += c(i, j)
		}
		return t
	}
	best := total(id)
	var bestP []int = id
	for _, p := range permutations(n) {
		if v := total(p); v < best {
			best, bestP = v, p
		}
	}
	var sd []string
	for i, s := range slots {
		sd = append(sd, edakit.Sprintf("%s at %s", letter(i), fmtPt(s)))
	}
	head := edakit.Sprintf("Independent-set matching picks equal-size cells that share no net, so their costs are independent, and reassigns them to their own locations by minimum-cost matching. The cells are %s (um). %s. The cost of putting cell i in slot j is the HPWL of i's nets with i at j. ", strings.Join(sd, ", "), strings.Join(desc, "; "))
	switch ctx.KPOrdinal {
	case 1:
		i, j := edakit.Int(0, n-1), edakit.Int(0, n-1)
		for i == j {
			j = edakit.Int(0, n-1)
		}
		return edakit.Numeric(head+edakit.Sprintf("What is the cost of putting cell %s in %s's current location?", letter(i), letter(j)), edakit.IntStr(c(i, j)),
			edakit.Sprintf("Move %s's pin to %s and recompute each of its nets: %d um.", letter(i), fmtPt(slots[j]), c(i, j)), map[string]string{"cost": edakit.IntStr(c(i, j))})
	case 3:
		return edakit.Numeric(head+"By how many um does the optimal matching reduce the total HPWL of these nets compared with leaving every cell where it is?", edakit.IntStr(total(id)-best),
			edakit.Sprintf("Leaving cells in place costs %d um. The best assignment sends A, B, C to the locations of %s, %s, %s for %d um, a gain of %d.", total(id), letter(bestP[0]), letter(bestP[1]), letter(bestP[2]), best, total(id)-best),
			map[string]string{"current": edakit.IntStr(total(id)), "best": edakit.IntStr(best)})
	}
	return edakit.Numeric(head+"What is the minimum total cost over all assignments of the three cells to the three locations?", edakit.IntStr(best),
		edakit.Sprintf("Build the 3 x 3 cost matrix and take the cheapest of the 6 assignments: A, B, C go to the locations of %s, %s, %s for %d um.", letter(bestP[0]), letter(bestP[1]), letter(bestP[2]), best),
		map[string]string{"best": edakit.IntStr(best)})
}

func ordinal(n int) string {
	return []string{"zeroth", "first", "second", "third", "fourth", "fifth", "sixth", "seventh", "eighth", "ninth"}[n]
}
