package graphs

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func mstGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1, 2:
		n := 4
		if ctx.KPOrdinal == 2 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			n = edakit.Scaled(ctx, 5, 6)
		}
		pins := distinctPins(n, 12)
		l := mstLen(pins)
		return edakit.Numeric("A net has pins "+ptsText(pins, "p")+". Using Manhattan distance, what is the length of its rectilinear minimum spanning tree?",
			edakit.IntStr(l),
			edakit.Sprintf("Prim's or Kruskal's algorithm on the complete graph of Manhattan distances gives total length %d (half-perimeter of the bounding box: %d).", l, hpwl(pins)),
			map[string]string{"mst": edakit.IntStr(l), "hpwl": edakit.IntStr(hpwl(pins))})
	default:
		for {
			n := edakit.Scaled(ctx, 5, 6)
			pins := distinctPins(n, 12)
			es, ok := primEdges(pins)
			if !ok {
				continue
			}
			k := edakit.Int(2, n-1)
			return edakit.Numeric("A net has pins "+ptsText(pins, "p")+edakit.Sprintf(". Prim's algorithm grows the rectilinear MST from p1, each step adding the shortest Manhattan edge from the tree to a new pin. What is the length of the edge added at step %d?", k),
				edakit.IntStr(es[k-1]),
				"The edges added in order have lengths "+intsText(es)+".",
				map[string]string{"edges": intsText(es)})
		}
	}
}

func hananGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 8)
	pins := make([]pt, 0, n)
	seen := map[pt]bool{}
	for len(pins) < n {
		p := pt{edakit.Int(0, 9), edakit.Int(0, 9)}
		if !seen[p] {
			seen[p] = true
			pins = append(pins, p)
		}
	}
	xs, ys := map[int]bool{}, map[int]bool{}
	for _, p := range pins {
		xs[p.x], ys[p.y] = true, true
	}
	X, Y := len(xs), len(ys)
	intro := "A net has pins " + ptsText(pins, "p") + ". The Hanan grid is formed by drawing a horizontal and a vertical line through every pin. "
	switch kpOf(ctx) {
	case 1:
		return edakit.Numeric(intro+"How many Hanan grid points (line intersections) are there?", edakit.IntStr(X*Y),
			edakit.Sprintf("There are %d distinct x values and %d distinct y values, so %d × %d = %d points. Pins sharing a coordinate share a line.", X, Y, X, Y, X*Y),
			map[string]string{"x_lines": edakit.IntStr(X), "y_lines": edakit.IntStr(Y)})
	case 2:
		return edakit.Numeric(intro+"How many Hanan points are candidate Steiner points, that is, not pins?", edakit.IntStr(X*Y-n),
			edakit.Sprintf("%d grid points minus %d pins = %d. By Hanan's theorem some optimal rectilinear Steiner tree uses only these candidates.", X*Y, n, X*Y-n),
			map[string]string{"candidates": edakit.IntStr(X*Y - n)})
	default:
		e := (X-1)*Y + X*(Y-1)
		return edakit.Numeric(intro+"Treating the Hanan grid as a graph whose vertices are the grid points and whose edges are the segments between neighbouring points, how many edges does it have?", edakit.IntStr(e),
			edakit.Sprintf("Each of %d horizontal lines has %d segments and each of %d vertical lines has %d: %d × %d + %d × %d = %d. An exact Steiner solver can search this graph.", Y, X-1, X, Y-1, Y, X-1, X, Y-1, e),
			map[string]string{"edges": edakit.IntStr(e)})
	}
}

func rsmtGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		pins := distinctPins(3, 12)
		l := hpwl(pins)
		return edakit.Numeric("A 3-pin net has pins "+ptsText(pins, "p")+". What is the length of its rectilinear Steiner minimum tree?",
			edakit.IntStr(l),
			edakit.Sprintf("For three pins the RSMT meets at the Steiner point (median x, median y) and its length equals the half-perimeter of the bounding box: %d. The rectilinear MST is %d.", l, mstLen(pins)),
			map[string]string{"rsmt": edakit.IntStr(l), "mst": edakit.IntStr(mstLen(pins))})
	case 2:
		pins := distinctPins(4, 10)
		l := rsmtLen(pins)
		return edakit.Numeric("A 4-pin net has pins "+ptsText(pins, "p")+". What is the length of its rectilinear Steiner minimum tree?",
			edakit.IntStr(l),
			edakit.Sprintf("Searching Steiner points on the Hanan grid (at most two are needed for four pins) gives %d. For comparison the MST is %d and the half-perimeter %d.", l, mstLen(pins), hpwl(pins)),
			map[string]string{"rsmt": edakit.IntStr(l), "mst": edakit.IntStr(mstLen(pins)), "hpwl": edakit.IntStr(hpwl(pins))})
	default:
		for {
			pins := distinctPins(4, 10)
			l, m := rsmtLen(pins), mstLen(pins)
			if m == l && edakit.Int(0, 2) > 0 {
				continue
			}
			return edakit.Numeric("A 4-pin net has pins "+ptsText(pins, "p")+". How much shorter is its rectilinear Steiner minimum tree than its rectilinear minimum spanning tree?",
				edakit.IntStr(m-l),
				edakit.Sprintf("MST %d, RSMT %d, saving %d.", m, l, m-l),
				map[string]string{"rsmt": edakit.IntStr(l), "mst": edakit.IntStr(m)})
		}
	}
}

func ratioGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		for {
			pins := distinctPins(4, 10)
			l, m := rsmtLen(pins), mstLen(pins)
			if m == l && edakit.Int(0, 3) > 0 {
				continue
			}
			r := float64(m) / float64(l)
			return edakit.Numeric("A 4-pin net has pins "+ptsText(pins, "p")+". What is the ratio of its rectilinear MST length to its RSMT length? Round to 2 decimal places.",
				edakit.Num(r, 2),
				edakit.Sprintf("MST %d / RSMT %d = %s. Hwang's theorem caps this ratio at 3/2 for any rectilinear net.", m, l, edakit.Num(r, 2)),
				map[string]string{"mst": edakit.IntStr(m), "rsmt": edakit.IntStr(l)})
		}
	case 2:
		n := edakit.Scaled(ctx, 4, 5)
		pins := distinctPins(n, 10)
		l, hp := rsmtLen(pins), hpwl(pins)
		for l == hp && edakit.Int(0, 2) > 0 {
			pins = distinctPins(n, 10)
			l, hp = rsmtLen(pins), hpwl(pins)
		}
		r := float64(l) / float64(hp)
		return edakit.Numeric(edakit.Sprintf("A %d-pin net has pins %s. What is the ratio of its RSMT length to its half-perimeter wirelength? Round to 2 decimal places.", n, ptsText(pins, "p")),
			edakit.Num(r, 2),
			edakit.Sprintf("RSMT %d / HPWL %d = %s. HPWL is a lower bound on any tree that connects the pins and is exact for two or three pins; the gap grows with pin count.", l, hp, edakit.Num(r, 2)),
			map[string]string{"rsmt": edakit.IntStr(l), "hpwl": edakit.IntStr(hp)})
	default:
		if edakit.Int(0, 1) == 0 {
			L := 2 * edakit.Int(10, 200)
			return edakit.Numeric(edakit.Sprintf("A net's rectilinear Steiner minimum tree has length %d. What is the largest length its rectilinear minimum spanning tree could have?", L),
				edakit.Num(1.5*float64(L), 1),
				edakit.Sprintf("Hwang (1976): RMST ≤ (3/2) RSMT for rectilinear nets, so at most %s. The bound is tight: four pins arranged as a plus sign reach it.", edakit.Num(1.5*float64(L), 1)),
				map[string]string{"bound": edakit.Num(1.5*float64(L), 1)})
		}
		M := 3 * edakit.Int(10, 200)
		return edakit.Numeric(edakit.Sprintf("A net's rectilinear minimum spanning tree has length %d. What is the smallest length its rectilinear Steiner minimum tree could have?", M),
			edakit.Num(float64(M)*2/3, 1),
			edakit.Sprintf("From RMST ≤ (3/2) RSMT, RSMT ≥ (2/3) × %d = %s. So an MST-based wirelength estimate overshoots the best Steiner tree by at most 50%%.", M, edakit.Num(float64(M)*2/3, 1)),
			map[string]string{"bound": edakit.Num(float64(M)*2/3, 1)})
	}
}

func oneSteinerGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 5)
	for {
		pins := distinctPins(n, 10)
		base := mstLen(pins)
		cands := hananPoints(pins)
		gain := func(cur []pt, s pt) int { return mstLen(cur) - mstLen(append(append([]pt(nil), cur...), s)) }
		bestG, bestP := 0, pt{}
		for _, s := range cands {
			if g := gain(pins, s); g > bestG {
				bestG, bestP = g, s
			}
		}
		if bestG == 0 {
			continue
		}
		intro := "A net has pins " + ptsText(pins, "p") + edakit.Sprintf(" and a rectilinear MST of length %d. ", base)
		switch kpOf(ctx) {
		case 1:
			s := edakit.Pick(cands)
			if edakit.Int(0, 1) == 0 {
				s = bestP
			}
			g := gain(pins, s)
			return edakit.Numeric(intro+edakit.Sprintf("Adding the Hanan point (%d, %d) as a Steiner point and recomputing the MST, by how much does the tree length drop?", s.x, s.y),
				edakit.IntStr(g),
				edakit.Sprintf("MST with the point: %d. Gain %d - %d = %d. A point that does not help gives gain 0 or less.", base-g, base, base-g, g),
				map[string]string{"gain": edakit.IntStr(g)})
		case 2:
			return edakit.Numeric(intro+"Over all Hanan points that are not pins, what is the largest drop in MST length from adding a single Steiner point?",
				edakit.IntStr(bestG),
				edakit.Sprintf("Trying every candidate, the best is (%d, %d), which lowers the MST from %d to %d.", bestP.x, bestP.y, base, base-bestG),
				map[string]string{"best_gain": edakit.IntStr(bestG)})
		default:
			cur := append([]pt(nil), pins...)
			for {
				bg, bp := 0, pt{}
				for _, s := range hananPoints(pins) {
					dup := false
					for _, c := range cur {
						dup = dup || c == s
					}
					if dup {
						continue
					}
					if g := gain(cur, s); g > bg {
						bg, bp = g, s
					}
				}
				if bg == 0 {
					break
				}
				cur = append(cur, bp)
			}
			final := mstLen(cur)
			return edakit.Numeric(intro+"The iterated 1-Steiner heuristic repeatedly adds the Hanan point with the largest positive MST gain until no point helps. What is the final tree length?",
				edakit.IntStr(final),
				edakit.Sprintf("It adds %d Steiner point(s) and ends at %d. The exact RSMT here is %d.", len(cur)-len(pins), final, rsmtLen(pins)),
				map[string]string{"final": edakit.IntStr(final), "rsmt": edakit.IntStr(rsmtLen(pins))})
		}
	}
}

func chainLen(start pt, cells []pt, order []int) int {
	l := 0
	cur := start
	for _, i := range order {
		l += manh(cur, cells[i])
		cur = cells[i]
	}
	return l
}

func orderText(order []int) string {
	var s []string
	for _, i := range order {
		s = append(s, edakit.Sprintf("f%d", i+1))
	}
	return strings.Join(s, ", ")
}

func scanTSPGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 5, 7)
	start := pt{0, 0}
	for {
		cells := make([]pt, n)
		for i := range cells {
			cells[i] = pt{edakit.Int(1, 12), edakit.Int(1, 12)}
		}
		intro := "Scan flip-flops sit at " + ptsText(cells, "f") + ". The chain starts at the scan-in port at (0, 0) and wires each flop to the next with Manhattan length. "
		switch kpOf(ctx) {
		case 1:
			order := edakit.Shuffle(seq(n))
			l := chainLen(start, cells, order)
			return edakit.Numeric(intro+"What is the total chain wirelength for the order "+orderText(order)+"?", edakit.IntStr(l),
				"Add the Manhattan distance of each hop, starting from the port: "+edakit.IntStr(l)+".",
				map[string]string{"length": edakit.IntStr(l)})
		case 2:
			used := make([]bool, n)
			cur := start
			var order []int
			ok := true
			for k := 0; k < n; k++ {
				b := -1
				for i := 0; i < n; i++ {
					if used[i] {
						continue
					}
					if b < 0 || manh(cur, cells[i]) < manh(cur, cells[b]) {
						b = i
					}
				}
				for i := 0; i < n; i++ {
					if !used[i] && i != b && manh(cur, cells[i]) == manh(cur, cells[b]) {
						ok = false
					}
				}
				used[b] = true
				order = append(order, b)
				cur = cells[b]
			}
			if !ok {
				continue
			}
			l := chainLen(start, cells, order)
			return edakit.Numeric(intro+"The nearest-neighbour heuristic always hops to the closest unvisited flop. What total wirelength does it produce?", edakit.IntStr(l),
				"Nearest neighbour visits "+orderText(order)+", for a total of "+edakit.IntStr(l)+". It is fast but can strand far flops for a long final hop.",
				map[string]string{"length": edakit.IntStr(l), "order": orderText(order)})
		default:
			order := edakit.Shuffle(seq(n))
			i := edakit.Int(0, n-3)
			j := edakit.Int(i+1, n-1)
			nw := append([]int(nil), order...)
			for a, b := i, j; a < b; a, b = a+1, b-1 {
				nw[a], nw[b] = nw[b], nw[a]
			}
			l0, l1 := chainLen(start, cells, order), chainLen(start, cells, nw)
			if l0 == l1 && edakit.Int(0, 3) > 0 {
				continue
			}
			return edakit.Numeric(intro+"The current order is "+orderText(order)+edakit.Sprintf(". A 2-opt move reverses the segment from position %d to position %d. By how much does the wirelength decrease (negative if it grows)?", i+1, j+1),
				edakit.IntStr(l0-l1),
				edakit.Sprintf("Old length %d; the new order %s has length %d; the change is %d. Only the hop into the segment and the hop out of it (if any) change, since a reversed segment keeps its internal hop lengths.", l0, orderText(nw), l1, l0-l1),
				map[string]string{"before": edakit.IntStr(l0), "after": edakit.IntStr(l1)})
		}
	}
}
