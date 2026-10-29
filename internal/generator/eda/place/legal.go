package place

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// displacement measures how far legalization moved cells.
func displacement(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 6)
	var desc []string
	tot, mx, sq := 0, 0, 0
	for i := 0; i < n; i++ {
		g := pt{edakit.Int(0, 60), edakit.Int(0, 30)}
		l := pt{g.x + edakit.Int(-6, 6), 2 * ((g.y + edakit.Int(-2, 2)) / 2)}
		if l.y < 0 {
			l.y = 0
		}
		d := manhattan(g, l)
		tot += d
		mx = maxI(mx, d)
		sq += d * d
		desc = append(desc, edakit.Sprintf("c%d from %s to %s", i+1, fmtPt(g), fmtPt(l)))
	}
	head := edakit.Sprintf("A legalizer moved cells from their global-placement positions to legal sites (um): %s. Displacement is Manhattan distance. ", strings.Join(desc, "; "))
	switch ctx.KPOrdinal {
	case 2:
		return edakit.Numeric(head+"What is the maximum displacement, in um?", edakit.IntStr(mx),
			edakit.Sprintf("The largest single move is %d um. A large maximum often signals a local overflow the global placer left behind.", mx), map[string]string{"max": edakit.IntStr(mx)})
	case 3:
		return edakit.Numeric(head+"An Abacus-style legalizer minimises the sum of squared displacements. What is that sum for this result, in um^2?", edakit.IntStr(sq),
			edakit.Sprintf("Square each displacement and add: %d um^2. Squaring punishes one long move more than several short ones.", sq), map[string]string{"sum_sq": edakit.IntStr(sq)})
	}
	return edakit.Numeric(head+"What is the total displacement, in um?", edakit.IntStr(tot),
		edakit.Sprintf("Add the Manhattan moves: %d um.", tot), map[string]string{"total": edakit.IntStr(tot)})
}

type lcell struct {
	name string
	x, y int // desired lower-left
	w    int
	e    float64
}

// tetris legalizes cells in x order into the row of least displacement.
func tetris(ctx generator.GeneratorContext) generator.Problem {
	rows := 2
	if ctx.KPOrdinal == 1 {
		rows = 1
	} else if ctx.Difficulty > 0.6 {
		rows = 3
	}
	rh := 2
	n := edakit.Scaled(ctx, 4, 7)
	cells := make([]lcell, n)
	xs := distinctInts(n, 0, 24)
	for i := range cells {
		cells[i] = lcell{name: edakit.Sprintf("c%d", i+1), x: xs[i], y: edakit.Int(0, rh*rows-1), w: edakit.Int(2, 6)}
		if rows == 1 {
			cells[i].y = 0
		}
	}
	order := append([]lcell(nil), cells...)
	sort.Slice(order, func(i, j int) bool { return order[i].x < order[j].x })
	right := make([]int, rows)
	type res struct{ x, row, d int }
	out := map[string]res{}
	total := 0
	for _, c := range order {
		best, bd, bx := -1, 0, 0
		for r := 0; r < rows; r++ {
			x := maxI(c.x, right[r])
			d := absI(x-c.x) + absI(r*rh-c.y)
			if best < 0 || d < bd {
				best, bd, bx = r, d, x
			}
		}
		right[best] = bx + c.w
		out[c.name] = res{bx, best, bd}
		total += bd
	}
	var desc []string
	for _, c := range cells {
		if rows == 1 {
			desc = append(desc, edakit.Sprintf("%s (width %d) wants x = %d", c.name, c.w, c.x))
		} else {
			desc = append(desc, edakit.Sprintf("%s (width %d) wants (%d, %d)", c.name, c.w, c.x, c.y))
		}
	}
	rowDesc := "a single row at y = 0"
	if rows > 1 {
		var ys []string
		for r := 0; r < rows; r++ {
			ys = append(ys, edakit.IntStr(r*rh))
		}
		rowDesc = edakit.Sprintf("%d rows of height %d at y = %s", rows, rh, strings.Join(ys, ", "))
	}
	head := edakit.Sprintf("Tetris legalization processes cells in order of increasing desired x. Each cell goes into the row where its displacement $|\\Delta x| + |\\Delta y|$ is smallest (ties to the lower row), at x = max(desired x, the right edge of the cells already in that row). The region has %s, starting at x = 0 with no right limit. Cells (lower-left corners): %s. ", rowDesc, strings.Join(desc, "; "))
	pick := cells[edakit.Int(0, n-1)]
	switch {
	case ctx.KPOrdinal == 1:
		return edakit.Numeric(head+edakit.Sprintf("Where does %s end up: what is its legal x?", pick.name), edakit.IntStr(out[pick.name].x),
			edakit.Sprintf("Packing left to right in x order, each cell starts at the larger of its wish and the row's current right edge. %s lands at x = %d.", pick.name, out[pick.name].x),
			map[string]string{"x": edakit.IntStr(out[pick.name].x)})
	case ctx.KPOrdinal == 2:
		return edakit.Numeric(head+edakit.Sprintf("In which row index does %s end up (row 0 is the bottom row)?", pick.name), edakit.IntStr(out[pick.name].row),
			edakit.Sprintf("When %s is processed it compares its displacement in every row, given the cells already there; the cheapest is row %d, at x = %d.", pick.name, out[pick.name].row, out[pick.name].x),
			map[string]string{"row": edakit.IntStr(out[pick.name].row)})
	}
	return edakit.Numeric(head+"What is the total displacement of all cells?", edakit.IntStr(total),
		edakit.Sprintf("Process the cells in x order, place each greedily and add the displacements: %d. Greedy order means early cells never move for later ones, which is why Tetris is fast but can leave long moves at the end of a row.", total),
		map[string]string{"total": edakit.IntStr(total)})
}

type cluster struct {
	first, last int
	e, q, w, x  float64
}

// abacusRow runs Abacus clustering on one row and returns each cell's final x.
func abacusRow(cells []lcell, lo, hi float64) []float64 {
	var cs []cluster
	collapse := func() {
		for {
			c := &cs[len(cs)-1]
			c.x = c.q / c.e
			if c.x < lo {
				c.x = lo
			}
			if c.x > hi-c.w {
				c.x = hi - c.w
			}
			if len(cs) < 2 {
				return
			}
			p := &cs[len(cs)-2]
			if p.x+p.w <= c.x {
				return
			}
			p.e += c.e
			p.q += c.q - c.e*p.w
			p.w += c.w
			p.last = c.last
			cs = cs[:len(cs)-1]
		}
	}
	for i, c := range cells {
		if len(cs) == 0 || cs[len(cs)-1].x+cs[len(cs)-1].w <= float64(c.x) {
			cs = append(cs, cluster{first: i, last: i, e: c.e, q: c.e * float64(c.x), w: float64(c.w), x: float64(c.x)})
		} else {
			l := &cs[len(cs)-1]
			l.e += c.e
			l.q += c.e * (float64(c.x) - l.w)
			l.w += float64(c.w)
			l.last = i
		}
		collapse()
	}
	out := make([]float64, len(cells))
	for _, c := range cs {
		x := c.x
		for i := c.first; i <= c.last; i++ {
			out[i] = x
			x += float64(cells[i].w)
		}
	}
	return out
}

// abacus asks for cluster and cell positions after Abacus row legalization.
func abacus(ctx generator.GeneratorContext) generator.Problem {
	weighted := ctx.KPOrdinal == 3
	n := 2
	if ctx.KPOrdinal != 1 {
		n = edakit.Scaled(ctx, 4, 6)
	}
	hi := float64(edakit.Int(30, 40))
	var cells []lcell
	found := false
	for try := 0; try < generator.MaxDrawAttempts && !found; try++ {
		xs := sortedCopy(distinctInts(n, 0, 28))
		cells = make([]lcell, n)
		tw := 0
		for i := range cells {
			cells[i] = lcell{name: edakit.Sprintf("c%d", i+1), x: xs[i], w: edakit.Int(2, 7), e: 1}
			if weighted {
				cells[i].e = float64(edakit.Int(1, 4))
			}
			tw += cells[i].w
		}
		if float64(tw) > hi {
			continue
		}
		overlaps := false
		for i := 0; i+1 < n; i++ {
			if cells[i].x+cells[i].w > cells[i+1].x {
				overlaps = true
			}
		}
		found = overlaps
	}
	if !found {
		cells = make([]lcell, n)
		for i := range cells {
			cells[i] = lcell{name: edakit.Sprintf("c%d", i+1), x: 3 * i, w: 4, e: 1}
			if weighted {
				cells[i].e = float64(1 + i%3)
			}
		}
	}
	pos := abacusRow(cells, 0, hi)
	var desc []string
	for _, c := range cells {
		if weighted {
			desc = append(desc, edakit.Sprintf("%s: width %d, desired x %d, weight %s", c.name, c.w, c.x, edakit.Num(c.e, 0)))
		} else {
			desc = append(desc, edakit.Sprintf("%s: width %d, desired x %d", c.name, c.w, c.x))
		}
	}
	obj := "the sum of squared x displacements"
	if weighted {
		obj = "the weighted sum $\\sum_i e_i (x_i - x_i')^2$ of squared x displacements"
	}
	head := edakit.Sprintf("Abacus legalizes one row spanning x = 0 to %s, keeping the cells in their order and minimising %s, where $x_i'$ is the desired position. Cells are added left to right; when a cell overlaps the last cluster it joins it, the cluster moves to its optimal position (clamped to the row) and merges with its left neighbour if they then overlap. Cells: %s. ", edakit.Num(hi, 0), obj, strings.Join(desc, "; "))
	if ctx.KPOrdinal == 1 {
		ans := edakit.Num(pos[0], 2)
		return edakit.Numeric(head+"The two cells overlap, so they form one cluster. What is the cluster's optimal left edge x? Round to 2 decimal places.", ans,
			edakit.Sprintf("For a cluster the optimum is $x_c = \\sum e_i (x_i' - o_i)/\\sum e_i$, where $o_i$ is the cell's offset inside the cluster: (%d + (%d - %d))/2 = %s, then clamped to the row.", cells[0].x, cells[1].x, cells[0].w, ans),
			map[string]string{"xc": ans})
	}
	k := edakit.Int(0, n-1)
	ans := edakit.Num(pos[k], 2)
	return edakit.Numeric(head+edakit.Sprintf("What is the final x of %s? Round to 2 decimal places.", cells[k].name), ans,
		edakit.Sprintf("Running the clustering gives final positions %s for c1 to c%d. Each cluster sits at the weighted mean of its cells' desired positions minus their offsets.", fmtFloats(pos, 2), n),
		map[string]string{"positions": fmtFloats(pos, 2)})
}
