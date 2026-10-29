package pv

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// hotspotMetrics: precision or recall of a hotspot detector.
func hotspotMetrics(ctx generator.GeneratorContext) generator.Problem {
	tp := edakit.Scaled(ctx, 20, 400)
	fp := edakit.Int(5, 300)
	fn := edakit.Int(1, 80)
	prec := 100 * float64(tp) / float64(tp+fp)
	rec := 100 * float64(tp) / float64(tp+fn)
	f1 := 2 * prec * rec / (prec + rec)
	stem := edakit.Sprintf("A hotspot detector is run on clips whose true status is known from lithography simulation. It flags %d true hotspots and %d clips that print correctly, and it misses %d true hotspots.", tp, fp, fn)
	if edakit.Int(0, 1) == 0 {
		q := stem + " What is its recall, in percent, rounded to 1 decimal place?"
		why := edakit.Sprintf("Recall = TP / (TP + FN) = %d / %d = %s percent: the share of real hotspots caught. Missed hotspots become yield loss, so recall is the metric a signoff flow cannot trade away.", tp, tp+fn, edakit.Num(rec, 1))
		return numChoice(q, rec, 1, "percent", why, prec, f1, 100*float64(fn)/float64(tp+fn))
	}
	q := stem + " What is its precision, in percent, rounded to 1 decimal place?"
	why := edakit.Sprintf("Precision = TP / (TP + FP) = %d / %d = %s percent: the share of flags that are real. False alarms cost simulation time or needless layout changes.", tp, tp+fp, edakit.Num(prec, 1))
	return numChoice(q, prec, 1, "percent", why, rec, f1, 100*float64(fp)/float64(tp+fp))
}

type grid [][]int

func randGrid(r, c int) grid {
	g := make(grid, r)
	for i := range g {
		g[i] = make([]int, c)
		for j := range g[i] {
			if edakit.Int(0, 9) < 4 {
				g[i][j] = 1
			}
		}
	}
	return g
}

func (g grid) rows() string {
	var out []string
	for _, row := range g {
		var b strings.Builder
		for _, v := range row {
			switch v {
			case 0:
				b.WriteByte('0')
			case 1:
				b.WriteByte('1')
			default:
				b.WriteByte('x')
			}
		}
		out = append(out, b.String())
	}
	return strings.Join(out, ", ")
}

func (g grid) rot() grid {
	r, c := len(g), len(g[0])
	out := make(grid, c)
	for i := range out {
		out[i] = make([]int, r)
		for j := range out[i] {
			out[i][j] = g[r-1-j][i]
		}
	}
	return out
}

func matchAt(g, p grid, i, j int) bool {
	for a := range p {
		for b := range p[a] {
			if p[a][b] == 2 {
				continue
			}
			if g[i+a][j+b] != p[a][b] {
				return false
			}
		}
	}
	return true
}

func countMatches(g grid, ps []grid) int {
	n := 0
	for i := 0; i+len(ps[0]) <= len(g); i++ {
		for j := 0; j+len(ps[0][0]) <= len(g[0]); j++ {
			for _, p := range ps {
				if i+len(p) <= len(g) && j+len(p[0]) <= len(g[0]) && matchAt(g, p, i, j) {
					n++
					break
				}
			}
		}
	}
	return n
}

// patternMatching: exact matches (KP 1), matches under rotation (KP 2), matches with
// don't-care cells (KP 3).
func patternMatching(ctx generator.GeneratorContext) generator.Problem {
	R, C := edakit.Scaled(ctx, 5, 8), edakit.Scaled(ctx, 5, 9)
	g := randGrid(R, C)
	k := kp(ctx)
	ps := 2
	if k == 3 || edakit.Int(0, 2) == 0 {
		ps = 3
	}
	i0, j0 := edakit.Int(0, R-ps), edakit.Int(0, C-ps)
	p := make(grid, ps)
	for a := range p {
		p[a] = append([]int(nil), g[i0+a][j0:j0+ps]...)
	}
	stem := edakit.Sprintf("A layer is rasterized into a %d by %d grid of pixels (1 = metal), rows top to bottom: %s. A hotspot pattern is the %d by %d block %s.", R, C, g.rows(), ps, ps, p.rows())
	switch k {
	case 1:
		n := countMatches(g, []grid{p})
		q := stem + " At how many positions does the pattern match the grid exactly (no rotation)?"
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("Slide the block over all %d positions and compare every pixel: %d match. Production matchers hash each window so a lookup replaces the pixel comparison.", (R-ps+1)*(C-ps+1), n),
			map[string]string{"matches": edakit.IntStr(n)})
	case 2:
		rs := []grid{p, p.rot(), p.rot().rot(), p.rot().rot().rot()}
		n := countMatches(g, rs)
		q := stem + " Lithography does not care about orientation, so the matcher also accepts the pattern rotated by 90, 180 or 270 degrees. At how many positions does some orientation match?"
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("Generate the four rotations, then count the positions where at least one of them matches: %d, against %d for the unrotated pattern alone.", n, countMatches(g, []grid{p})),
			map[string]string{"matches": edakit.IntStr(n)})
	default:
		nd := edakit.Int(2, 4)
		for _, c := range edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8})[:nd] {
			p[c/3][c%3] = 2
		}
		n := countMatches(g, []grid{p})
		stem = edakit.Sprintf("A layer is rasterized into a %d by %d grid of pixels (1 = metal), rows top to bottom: %s. A hotspot pattern is the 3 by 3 block %s, where x marks a don't-care pixel that matches either value.", R, C, g.rows(), p.rows())
		q := stem + " At how many positions does the pattern match (no rotation)?"
		return edakit.Numeric(q, edakit.IntStr(n),
			edakit.Sprintf("Compare only the %d constrained pixels at each of the %d positions: %d match. Don't-care pixels make a pattern cover a family of layouts, at the price of more false alarms.", 9-nd, (R-2)*(C-2), n),
			map[string]string{"matches": edakit.IntStr(n)})
	}
}

// graphComponents returns, for each component of an undirected graph, whether it is bipartite.
func graphComponents(n int, edges [][2]int) []bool {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	col := make([]int, n)
	for i := range col {
		col[i] = -1
	}
	var out []bool
	for s := 0; s < n; s++ {
		if col[s] >= 0 {
			continue
		}
		ok := true
		col[s] = 0
		queue := []int{s}
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for _, u := range adj[v] {
				if col[u] < 0 {
					col[u] = 1 - col[v]
					queue = append(queue, u)
				} else if col[u] == col[v] {
					ok = false
				}
			}
		}
		out = append(out, ok)
	}
	return out
}

// conflictGraph draws a graph whose vertices carry hidden colours; most edges join different
// colours, and an odd edge joins equal colours with the given probability out of 10.
func conflictGraph(n, m, oddTenths int) [][2]int {
	hidden := make([]int, n)
	for i := range hidden {
		hidden[i] = edakit.Int(0, 1)
	}
	seen := map[[2]int]bool{}
	var es [][2]int
	for tries := 0; len(es) < m && tries < 500; tries++ {
		a, b := edakit.Int(0, n-1), edakit.Int(0, n-1)
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		if seen[[2]int{a, b}] {
			continue
		}
		if hidden[a] == hidden[b] && edakit.Int(0, 9) >= oddTenths {
			continue
		}
		seen[[2]int{a, b}] = true
		es = append(es, [2]int{a, b})
	}
	return es
}

func describeEdges(es [][2]int) string {
	var parts []string
	for _, e := range es {
		parts = append(parts, edakit.Sprintf("F%d-F%d", e[0]+1, e[1]+1))
	}
	return strings.Join(parts, ", ")
}

// dpColoring: conflict edges among features (KP 1), components that cannot be 2-coloured
// (KP 2), number of valid mask assignments (KP 3).
func dpColoring(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		S := edakit.Int(25, 45)
		n := edakit.Scaled(ctx, 4, 8)
		rs := separated(n, 160, 10, 40)
		e := 0
		for i := range rs {
			for j := i + 1; j < n; j++ {
				if dist2(rs[i], rs[j]) < S*S {
					e++
				}
			}
		}
		q := edakit.Sprintf("In double patterning, two features closer than %d nm (Euclidean) cannot be printed on the same mask, so a conflict edge joins them. The features are %s (nm). How many conflict edges does the graph have?", S, listRects("F", rs))
		return edakit.Numeric(q, edakit.IntStr(e),
			edakit.Sprintf("Test all %d pairs against the same-mask spacing %d nm: %d pairs are closer, each one an edge that forces the two features onto different masks.", n*(n-1)/2, S, e),
			map[string]string{"edges": edakit.IntStr(e)})
	case 2:
		n := edakit.Scaled(ctx, 6, 10)
		es := conflictGraph(n, edakit.Int(n-1, n+3), 3)
		comps := graphComponents(n, es)
		bad := 0
		for _, ok := range comps {
			if !ok {
				bad++
			}
		}
		q := edakit.Sprintf("A double-patterning conflict graph has features F1 to F%d and conflict edges %s. How many connected components cannot be split across two masks?", n, describeEdges(es))
		return edakit.Numeric(q, edakit.IntStr(bad),
			edakit.Sprintf("A component is 2-colourable exactly when it has no odd cycle; a breadth-first search that alternates colours finds a same-colour edge only in a component with an odd cycle. %d of the %d components contain one and need a layout change or a stitch.", bad, len(comps)),
			map[string]string{"components": edakit.IntStr(len(comps)), "odd": edakit.IntStr(bad)})
	default:
		n := edakit.Scaled(ctx, 5, 9)
		es := conflictGraph(n, edakit.Int(n-3, n), 1)
		comps := graphComponents(n, es)
		ways := 1
		for _, ok := range comps {
			if !ok {
				ways = 0
				break
			}
			ways *= 2
		}
		q := edakit.Sprintf("A double-patterning conflict graph has features F1 to F%d and conflict edges %s. Every feature must go on mask A or mask B, and the two ends of every edge must differ. How many valid assignments are there (assignments that swap A and B count as different)?", n, describeEdges(es))
		return edakit.Numeric(q, edakit.IntStr(ways),
			edakit.Sprintf("There are %d components, counting isolated features. A bipartite component has exactly 2 colourings and a component with an odd cycle has none, so the count is %d. A decomposer exploits the freedom between components to balance mask density.", len(comps), ways),
			map[string]string{"components": edakit.IntStr(len(comps))})
	}
}

// criticalArea: short critical area at one defect size (KP 1), averaged over the defect size
// distribution (KP 2), yield from the average critical area (KP 3).
func criticalArea(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		L := 10 * edakit.Int(5, 100)
		s := 5 * edakit.Int(8, 30)
		x := s + 5*edakit.Int(-4, 40)
		a := 0.0
		if x > s {
			a = float64(L) * float64(x-s) / 1000
		}
		ans := edakit.Num(a, 2)
		q := edakit.Sprintf("Two parallel wires run side by side for %d um with a %d nm gap. A circular conducting defect of diameter %d nm shorts them if it touches both. Ignoring end effects, what is the critical area for shorts at this defect size, in square um? Round to 2 decimal places.", L, s, x)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The defect centre must fall in a band of width x - s = %d nm along the run (zero if the defect is smaller than the gap): %d um x %d nm = %s square um.", imax(0, x-s), L, imax(0, x-s), ans),
			map[string]string{"band_nm": edakit.IntStr(imax(0, x-s))})
	case 2:
		n := edakit.Scaled(ctx, 2, 8)
		L := 10 * edakit.Int(10, 200)
		x0 := 5 * edakit.Int(3, 8)
		s := x0 + 5*edakit.Int(1, 20)
		a := float64(n-1) * float64(L) * float64(x0*x0) / float64(s) / 1000
		ans := edakit.Num(a, 2)
		q := edakit.Sprintf("A bus of %d parallel wires, each %d um long, has %d nm gaps. Defect diameters follow $f(x) = 2 x_0^2 / x^3$ for $x \\ge x_0$, with $x_0 = %d$ nm, which gives an average short critical area of $L x_0^2 / s$ per gap. Counting each gap separately, what is the bus's average critical area for shorts, in square um? Round to 2 decimal places.", n, L, s, x0)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Integrating L (x - s) f(x) from s to infinity gives L x0^2 / s per gap. With %d gaps: %d x %d um x %d^2 / %d nm = %s square um. Doubling the spacing halves it, which is why wire spreading pays.", n-1, n-1, L, x0, s, ans),
			map[string]string{"gaps": edakit.IntStr(n - 1)})
	default:
		A := float64(edakit.Int(5, 200)) / 10
		D := float64(edakit.Int(5, 50)) / 100
		y := 100 * math.Exp(-D*A/100)
		ans := edakit.Num(y, 2)
		q := edakit.Sprintf("A die's average critical area for shorts on metal 2 is %s square mm, and the defect density for that mechanism is %s defects per square cm. Using the Poisson model $Y = e^{-D A_c}$, what is the yield for this mechanism, in percent? Round to 2 decimal places.", edakit.Num(A, 1), edakit.Num(D, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Convert: %s square mm = %s square cm. D x A = %s, and e^(-%s) = %s percent. Yields of independent mechanisms and layers multiply.", edakit.Num(A, 1), edakit.Num(A/100, 3), edakit.Num(D*A/100, 4), edakit.Num(D*A/100, 4), ans),
			map[string]string{"lambda": edakit.Num(D*A/100, 4)})
	}
}

// dfmScore: share of locations meeting a recommended rule.
func dfmScore(ctx generator.GeneratorContext) generator.Problem {
	if edakit.Int(0, 1) == 0 {
		n := edakit.Scaled(ctx, 50, 5000)
		d := edakit.Int(n/5, n)
		pct := 100 * float64(d) / float64(n)
		q := edakit.Sprintf("A routed block has %d single-cut via locations, and a via-doubling pass managed to add a redundant cut at %d of them. The DFM report scores redundant-via coverage as the percent of locations with two or more cuts. What is the score, rounded to 1 decimal place?", n, d)
		why := edakit.Sprintf("%d / %d = %s percent. Recommended rules are scored like this rather than pass/fail: the design is legal either way, but every single-cut via left is a yield risk.", d, n, edakit.Num(pct, 1))
		return numChoice(q, pct, 1, "percent", why, 100-pct, 100*float64(d)/float64(n+d), pct/2)
	}
	mn := 10 * edakit.Int(4, 6)
	rec := mn + 10*edakit.Int(1, 3)
	k := edakit.Scaled(ctx, 6, 12)
	var sp []int
	ok := 0
	for i := 0; i < k; i++ {
		s := mn + 5*edakit.Int(0, (rec-mn)/5+4)
		sp = append(sp, s)
		if s >= rec {
			ok++
		}
	}
	pct := 100 * float64(ok) / float64(k)
	q := edakit.Sprintf("The minimum metal spacing is %d nm and the recommended spacing is %d nm. A block's neighbouring-wire spacings, in nm, are %s. What percent of them meet the recommended value? Round to 1 decimal place.", mn, rec, joinInts(sp))
	why := edakit.Sprintf("%d of %d spacings are at least %d nm: %s percent. All of them pass the mandatory rule; the recommended rule measures margin.", ok, k, rec, edakit.Num(pct, 1))
	return numChoice(q, pct, 1, "percent", why, 100-pct, 100*float64(ok)/float64(k+1), 100)
}

// wireSpreading: reduction in average short critical area when a gap is widened.
func wireSpreading(ctx generator.GeneratorContext) generator.Problem {
	s := 10 * edakit.Int(4, 8)
	s2 := s + 10*edakit.Int(1, 6)
	red := 100 * (1 - float64(s)/float64(s2))
	q := edakit.Sprintf("Under the defect-size model where the average short critical area of a gap is $L x_0^2 / s$, a DFM pass spreads a pair of long parallel wires from %d nm to %d nm apart, keeping their length. By what percent does the pair's average short critical area fall? Round to 1 decimal place.", s, s2)
	why := edakit.Sprintf("The area scales as 1 / s, so the new value is %d / %d of the old: a reduction of %s percent. Spreading uses free routing tracks, so it is applied where there is slack, not everywhere.", s, s2, edakit.Num(red, 1))
	return numChoice(q, red, 1, "percent", why, 100*float64(s2-s)/float64(s), 100*(1-float64(s*s)/float64(s2*s2)), 100*float64(s)/float64(s2))
}
