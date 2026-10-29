package place

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// hpwlGen asks for the HPWL of one net, a larger net, or the weighted total of several nets.
func hpwlGen(ctx generator.GeneratorContext) generator.Problem {
	if ctx.KPOrdinal == 3 || (ctx.KPOrdinal == 0 && ctx.Difficulty > 0.7) {
		m := edakit.Int(3, 4)
		total := 0
		var desc, parts []string
		for i := 0; i < m; i++ {
			pins := randPts(edakit.Int(2, 4), 0, 60)
			w := edakit.Int(1, 3)
			h := hpwl(pins)
			total += w * h
			desc = append(desc, edakit.Sprintf("net n%d (weight %d) has pins at %s", i+1, w, fmtPts(pins)))
			parts = append(parts, edakit.Sprintf("%d x %d", w, h))
		}
		return edakit.Numeric("Coordinates are in um. In a placement, "+strings.Join(desc, "; ")+". What is the total weighted half-perimeter wirelength, in um?", edakit.IntStr(total),
			edakit.Sprintf("Compute each net's bounding-box half-perimeter and multiply by its weight: %s, which sums to %d um.", strings.Join(parts, " + "), total),
			map[string]string{"total": edakit.IntStr(total)})
	}
	n := 3
	if ctx.KPOrdinal == 2 {
		n = edakit.Scaled(ctx, 4, 8)
	} else if ctx.KPOrdinal == 0 {
		n = edakit.Scaled(ctx, 3, 7)
	}
	pins := randPts(n, 0, 80)
	h := hpwl(pins)
	xs, ys := make([]int, n), make([]int, n)
	for i, p := range pins {
		xs[i], ys[i] = p.x, p.y
	}
	return edakit.Numeric(edakit.Sprintf("A net has %d pins at %s (um). What is its half-perimeter wirelength, in um?", n, fmtPts(pins)), edakit.IntStr(h),
		edakit.Sprintf("The bounding box spans x from %d to %d and y from %d to %d, so HPWL = %d + %d = %d um. Only the extreme pins matter.", minI(xs...), maxI(xs...), minI(ys...), maxI(ys...), span(xs), span(ys), h),
		map[string]string{"hpwl": edakit.IntStr(h)})
}

// incremental moves one cell and asks for the change in HPWL over its incident nets.
func incremental(ctx generator.GeneratorContext) generator.Problem {
	k := 1
	if ctx.KPOrdinal >= 2 || ctx.KPOrdinal == 0 {
		k = edakit.Scaled(ctx, 2, 3)
	}
	from := pt{edakit.Int(10, 50), edakit.Int(10, 50)}
	to := pt{edakit.Int(0, 60), edakit.Int(0, 60)}
	for try := 0; to == from && try < generator.MaxDrawAttempts; try++ {
		to = pt{edakit.Int(0, 60), edakit.Int(0, 60)}
	}
	if to == from {
		to = pt{from.x + 7, from.y}
	}
	before, after := 0, 0
	var desc, parts []string
	for i := 0; i < k; i++ {
		others := randPts(edakit.Int(1, 3), 0, 60)
		b := hpwl(append([]pt{from}, others...))
		a := hpwl(append([]pt{to}, others...))
		before += b
		after += a
		desc = append(desc, edakit.Sprintf("net n%d connects c to pins at %s", i+1, fmtPts(others)))
		parts = append(parts, edakit.Sprintf("n%d: %d to %d", i+1, b, a))
	}
	extra := ""
	if ctx.KPOrdinal == 3 {
		extra = edakit.Sprintf(" The design has %d other nets that do not touch c.", edakit.Int(1000, 90000))
	}
	d := after - before
	return edakit.Numeric(edakit.Sprintf("Cell c (treat it as a single pin) moves from %s to %s (um). In the placement, %s.%s What is the change in total HPWL, new minus old, in um?", fmtPt(from), fmtPt(to), strings.Join(desc, "; "), extra),
		edakit.IntStr(d),
		edakit.Sprintf("Only nets incident to c can change. Recomputing their bounding boxes: %s. The change is %d - %d = %d um.", strings.Join(parts, "; "), after, before, d),
		map[string]string{"before": edakit.IntStr(before), "after": edakit.IntStr(after)})
}

// cliqueStar computes clique and star quadratic costs of a net in x.
func cliqueStar(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 3, 5)
	xs := make([]int, k)
	for i := range xs {
		xs[i] = edakit.Int(0, 20)
	}
	switch ctx.KPOrdinal {
	case 2:
		mean := 0.0
		for _, x := range xs {
			mean += float64(x)
		}
		mean /= float64(k)
		cost := 0.0
		for _, x := range xs {
			cost += (float64(x) - mean) * (float64(x) - mean)
		}
		ans := edakit.Num(cost, 2)
		return edakit.Numeric(edakit.Sprintf("A %d-pin net has pins at x = %s (um). In the star model a new star node s connects to every pin with unit-weight two-pin edges, and the quadratic cost in x is the sum of $(x_i - x_s)^2$. Place the star node where this cost is smallest. What is the minimum cost, in um^2? Round to 2 decimal places.", k, fmtInts(xs)), ans,
			edakit.Sprintf("The cost is smallest at the mean, $x_s$ = %s. The sum of squared deviations from it is %s um^2.", edakit.Num(mean, 4), ans),
			map[string]string{"x_s": edakit.Num(mean, 4)})
	case 3:
		big := edakit.Int(6, 60)
		d := big*(big-1)/2 - big
		return edakit.Numeric(edakit.Sprintf("A clock-enable net has %d pins. A quadratic placer models it either as a clique of two-pin edges or as a star with one extra node. How many more two-pin edges does the clique model add to the system than the star model?", big), edakit.IntStr(d),
			edakit.Sprintf("A clique on %d pins has %d x %d / 2 = %d edges; a star has %d. The difference is %d, which is why large nets are modeled as stars.", big, big, big-1, big*(big-1)/2, big, d),
			map[string]string{"clique_edges": edakit.IntStr(big * (big - 1) / 2)})
	}
	sum := 0
	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			sum += (xs[i] - xs[j]) * (xs[i] - xs[j])
		}
	}
	cost := float64(sum) / float64(k-1)
	ans := edakit.Num(cost, 2)
	return edakit.Numeric(edakit.Sprintf("A %d-pin net has pins at x = %s (um). The clique model joins every pair of pins with a two-pin edge of weight $1/(k-1)$, where k is the pin count. What is the quadratic wirelength in x, the weighted sum of $(x_i - x_j)^2$ over all pairs, in um^2? Round to 2 decimal places.", k, fmtInts(xs)), ans,
		edakit.Sprintf("The %d pairwise squared distances sum to %d. Times 1/%d gives %s um^2.", k*(k-1)/2, sum, k-1, ans),
		map[string]string{"pair_sum": edakit.IntStr(sum)})
}

// b2b exercises the bound-to-bound net model.
func b2b(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 3, 6)
	xs := distinctInts(k, 0, 40)
	s := sortedCopy(xs)
	lo, hi := s[0], s[k-1]
	head := edakit.Sprintf("A %d-pin net has pins at x = %s (um). The bound-to-bound model connects the leftmost and rightmost pins to each other and every inner pin to both of them, with edge weight $w_{ij} = 1/((k-1)|x_i - x_j|)$ computed from the current positions. ", k, fmtInts(xs))
	switch ctx.KPOrdinal {
	case 1:
		in := s[edakit.Int(1, k-2)]
		w := 1 / (float64(k-1) * float64(in-lo))
		ans := edakit.Num(w, 4)
		return edakit.Numeric(head+edakit.Sprintf("What is the weight of the edge between the pin at x = %d and the left bound? Round to 4 decimal places.", in), ans,
			edakit.Sprintf("The left bound is x = %d. The distance is %d, so the weight is 1/(%d x %d) = %s.", lo, in-lo, k-1, in-lo, ans), map[string]string{"w": ans})
	case 3:
		sizes := make([]int, edakit.Int(3, 5))
		total := 0
		for i := range sizes {
			sizes[i] = edakit.Int(2, 12)
			total += 2*sizes[i] - 3
		}
		return edakit.Numeric(edakit.Sprintf("A netlist has nets with pin counts %s. How many two-pin edges does the bound-to-bound model create in the x system for all of them together?", fmtInts(sizes)), edakit.IntStr(total),
			edakit.Sprintf("A k-pin net gives one bound-to-bound edge plus two edges per inner pin: 1 + 2(k - 2) = 2k - 3. Summing over the nets gives %d.", total),
			map[string]string{"edges": edakit.IntStr(total)})
	}
	cost := 0.0
	cost += float64(hi-lo) / float64(k-1)
	for _, x := range s[1 : k-1] {
		cost += float64(x-lo)/float64(k-1) + float64(hi-x)/float64(k-1)
	}
	ans := edakit.Num(cost, 2)
	return edakit.Numeric(head+"What is the model's quadratic cost, the sum of $w_{ij}(x_i - x_j)^2$ over its edges, in um? Round to 2 decimal places.", ans,
		edakit.Sprintf("Each edge contributes $w_{ij}(x_i-x_j)^2 = |x_i - x_j|/(k-1)$. The bound edge gives %d, and every inner pin gives its distances to both bounds, which sum to %d. That is %d x %d / %d = %s um, exactly the net's span: the model reproduces HPWL at the current positions.", hi-lo, hi-lo, k-1, hi-lo, k-1, ans),
		map[string]string{"span": edakit.IntStr(hi - lo)})
}

// rmst is the rectilinear minimum spanning tree length by Prim's algorithm.
func rmst(pts []pt) int {
	n := len(pts)
	in := make([]bool, n)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = 1 << 30
	}
	dist[0] = 0
	total := 0
	for it := 0; it < n; it++ {
		b := -1
		for i := 0; i < n; i++ {
			if !in[i] && (b < 0 || dist[i] < dist[b]) {
				b = i
			}
		}
		in[b] = true
		total += dist[b]
		for i := 0; i < n; i++ {
			if !in[i] {
				if d := manhattan(pts[b], pts[i]); d < dist[i] {
					dist[i] = d
				}
			}
		}
	}
	return total
}

// rsmt4 finds the rectilinear Steiner minimum tree length of up to four terminals by trying
// every set of at most n-2 Steiner points on the Hanan grid.
func rsmt4(pts []pt) int {
	var xs, ys []int
	for _, p := range pts {
		xs = append(xs, p.x)
		ys = append(ys, p.y)
	}
	var hanan []pt
	for _, x := range xs {
		for _, y := range ys {
			hanan = append(hanan, pt{x, y})
		}
	}
	best := rmst(pts)
	for i := range hanan {
		if l := rmst(append(append([]pt(nil), pts...), hanan[i])); l < best {
			best = l
		}
		if len(pts) >= 4 {
			for j := i + 1; j < len(hanan); j++ {
				if l := rmst(append(append([]pt(nil), pts...), hanan[i], hanan[j])); l < best {
					best = l
				}
			}
		}
	}
	return best
}

// steiner compares HPWL, RMST and RSMT estimates.
func steiner(ctx generator.GeneratorContext) generator.Problem {
	switch ctx.KPOrdinal {
	case 1:
		p := randPts(3, 0, 40)
		h := hpwl(p)
		return edakit.Numeric(edakit.Sprintf("A three-pin net has pins at %s (um). What is the length of its rectilinear Steiner minimum tree, in um?", fmtPts(p)), edakit.IntStr(h),
			edakit.Sprintf("For three pins the RSMT joins them through the median point (median x, median y), and its length equals the HPWL: %d um.", h),
			map[string]string{"hpwl": edakit.IntStr(h)})
	case 3:
		p := []pt{{0, 5}, {5, 0}, {10, 5}, {5, 10}}
		for try := 0; try < generator.MaxDrawAttempts; try++ {
			if c := randPts(4, 0, 30); rmst(c) > rsmt4(c) {
				p = c
				break
			}
		}
		l := rsmt4(p)
		return edakit.Numeric(edakit.Sprintf("A four-pin net has pins at %s (um). Steiner points may be added anywhere on the Hanan grid. What is the length of its rectilinear Steiner minimum tree, in um?", fmtPts(p)), edakit.IntStr(l),
			edakit.Sprintf("The rectilinear MST is %d um and the HPWL is %d um. Adding the best Steiner points on the Hanan grid lowers the tree to %d um, which lies between the two.", rmst(p), hpwl(p), l),
			map[string]string{"rmst": edakit.IntStr(rmst(p)), "hpwl": edakit.IntStr(hpwl(p)), "rsmt": edakit.IntStr(l)})
	}
	n := edakit.Scaled(ctx, 4, 6)
	p := randPts(n, 0, 50)
	l := rmst(p)
	return edakit.Numeric(edakit.Sprintf("A net has pins at %s (um). What is the length of its rectilinear minimum spanning tree (Manhattan edge lengths, no Steiner points), in um?", fmtPts(p)), edakit.IntStr(l),
		edakit.Sprintf("Prim's algorithm from the first pin, always adding the nearest pin in Manhattan distance, builds a tree of %d um. The HPWL lower bound is %d um.", l, hpwl(p)),
		map[string]string{"rmst": edakit.IntStr(l), "hpwl": edakit.IntStr(hpwl(p))})
}
