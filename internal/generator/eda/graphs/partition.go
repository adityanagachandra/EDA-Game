package graphs

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func cutSizeGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 10)
	h := randHyper(n, edakit.Scaled(ctx, 5, 9), 4)
	switch kpOf(ctx) {
	case 1, 2:
		weighted := ctx.KPOrdinal == 2 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0)
		if weighted {
			for j := range h.w {
				h.w[j] = edakit.Int(1, 5)
			}
		}
		side := randSides(n)
		cut := cutSize(h, side, weighted)
		var cutNets []string
		for j, net := range h.nets {
			if isCut(net, side) {
				cutNets = append(cutNets, edakit.Sprintf("n%d", j+1))
			}
		}
		q := "Nets: " + netText(h, weighted) + ". The partition is A = " + sideText(side, 0) + ", B = " + sideText(side, 1) + ". "
		if weighted {
			q += "What is the weighted cut, the total weight of nets with pins on both sides?"
		} else {
			q += "How many nets are cut (have pins on both sides)?"
		}
		return edakit.Numeric(q, edakit.IntStr(cut),
			"A net is cut when its pins are not all on one side. The cut nets are "+strings.Join(cutNets, ", ")+", giving "+edakit.IntStr(cut)+". A net counts once however many of its pins cross.",
			map[string]string{"cut": edakit.IntStr(cut)})
	default:
		side := make([]int, n)
		for {
			cnt := [3]int{}
			for i := range side {
				side[i] = edakit.Int(0, 2)
				cnt[side[i]]++
			}
			if cnt[0] > 0 && cnt[1] > 0 && cnt[2] > 0 {
				break
			}
		}
		conn, cutNets := 0, 0
		for _, net := range h.nets {
			blocks := map[int]bool{}
			for _, p := range net {
				blocks[side[p]] = true
			}
			conn += len(blocks) - 1
			if len(blocks) > 1 {
				cutNets++
			}
		}
		var parts []string
		for b, nm := range []string{"A", "B", "C"} {
			parts = append(parts, nm+" = "+sideText(side, b))
		}
		return edakit.Numeric("Nets: "+netText(h, false)+". A 3-way partition is "+strings.Join(parts, ", ")+". What is the connectivity cost, the sum over nets of (number of blocks the net touches minus 1)?",
			edakit.IntStr(conn),
			edakit.Sprintf("Each net touching λ blocks costs λ - 1. The total is %d, while the plain cut-net count is %d: a net spanning all three blocks costs 2 under connectivity but 1 as a cut net.", conn, cutNets),
			map[string]string{"connectivity": edakit.IntStr(conn), "cut_nets": edakit.IntStr(cutNets)})
	}
}

func cliqueDistortionGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		k := edakit.Scaled(ctx, 4, 16)
		m := edakit.Int(1, k/2)
		return edakit.Numeric(edakit.Sprintf("A net with %d pins is modeled as a clique. A partition puts %d of its pins on side A and the rest on side B. How many clique edges cross the cut?", k, m),
			edakit.IntStr(m*(k-m)),
			edakit.Sprintf("Every pin on A connects to every pin on B: %d × %d = %d edges, though the hyperedge itself is cut only once.", m, k-m, m*(k-m)),
			map[string]string{"crossing": edakit.IntStr(m * (k - m))})
	case 2:
		k := edakit.Scaled(ctx, 4, 16)
		m := edakit.Int(1, k/2)
		v := float64(m*(k-m)) / float64(k-1)
		return edakit.Numeric(edakit.Sprintf("A net with k = %d pins is modeled as a clique with edge weight $1/(k-1)$. A partition splits its pins %d to %d between the two sides. What is the weight of the clique edges that cross? Round to 2 decimal places.", k, m, k-m),
			edakit.Num(v, 2),
			edakit.Sprintf("%d × %d crossing edges, each of weight 1/%d: %s. The hypergraph would charge 1 for any split, while the clique charges between 1 (one pin split off) and about k/4 (an even split).", m, k-m, k-1, edakit.Num(v, 2)),
			map[string]string{"weight": edakit.Num(v, 2)})
	default:
		n := edakit.Scaled(ctx, 6, 9)
		h := randHyper(n, edakit.Scaled(ctx, 4, 7), 5)
		side := randSides(n)
		g := 0.0
		for _, net := range h.nets {
			a := 0
			for _, p := range net {
				if side[p] == 0 {
					a++
				}
			}
			g += float64(a*(len(net)-a)) / float64(len(net)-1)
		}
		hc := cutSize(h, side, false)
		return edakit.Numeric("Nets: "+netText(h, false)+". The partition is A = "+sideText(side, 0)+", B = "+sideText(side, 1)+". Expand every k-pin net as a clique with edge weight $1/(k-1)$. What is the weighted graph cut? Round to 2 decimal places.",
			edakit.Num(g, 2),
			edakit.Sprintf("Each net with a pins on A contributes a(k - a)/(k - 1). The graph cut is %s, while the hypergraph cut is %d nets.", edakit.Num(g, 2), hc),
			map[string]string{"graph_cut": edakit.Num(g, 2), "hyper_cut": edakit.IntStr(hc)})
	}
}

func binom(n, k int) int {
	r := 1
	for i := 1; i <= k; i++ {
		r = r * (n - k + i) / i
	}
	return r
}

func balanceGen(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		A := 2 * edakit.Scaled(ctx, 50, 500)
		eps := edakit.Pick([]int{2, 5, 10, 20})
		hi := float64(A) / 2 * (1 + float64(eps)/100)
		return edakit.Numeric(edakit.Sprintf("A bipartitioner must keep each side's cell area between $(1-\\epsilon)A/2$ and $(1+\\epsilon)A/2$, where $A$ is the total area. With A = %d and an imbalance tolerance of %d%%, what is the largest area one side may hold? Round to 1 decimal place.", A, eps),
			edakit.Num(hi, 1),
			edakit.Sprintf("(1 + %s) × %d / 2 = %s.", edakit.Num(float64(eps)/100, 2), A, edakit.Num(hi, 1)),
			map[string]string{"upper": edakit.Num(hi, 1)})
	case 2:
		for {
			n := edakit.Scaled(ctx, 6, 10)
			area := make([]int, n)
			for i := range area {
				area[i] = edakit.Int(1, 12)
			}
			side := randSides(n)
			sa, total := 0, 0
			for i, a := range area {
				total += a
				if side[i] == 0 {
					sa += a
				}
			}
			eps := edakit.Pick([]int{10, 20})
			lo := float64(total) / 2 * (1 - float64(eps)/100)
			hi := float64(total) / 2 * (1 + float64(eps)/100)
			best, legal, illegal := -1, 0, 0
			for i, a := range area {
				if side[i] != 0 {
					continue
				}
				if float64(sa-a) >= lo && float64(total-sa+a) <= hi {
					legal++
					best = max(best, a)
				} else {
					illegal++
				}
			}
			if best < 0 || illegal == 0 {
				continue
			}
			var at []string
			for i, a := range area {
				at = append(at, edakit.Sprintf("%s:%d", cellName(i), a))
			}
			return edakit.Numeric("Cell areas: "+strings.Join(at, ", ")+". Side A = "+sideText(side, 0)+" and side B = "+sideText(side, 1)+edakit.Sprintf(". Each side must stay between %s and %s (total %d, tolerance %d%%). What is the largest area of a single cell on side A that can legally move to side B?", edakit.Num(lo, 1), edakit.Num(hi, 1), total, eps),
				edakit.IntStr(best),
				edakit.Sprintf("Side A holds %d. A move of area a leaves A with %d - a and B with %d + a, and both must stay in range. The largest legal cell has area %d.", sa, sa, total-sa, best),
				map[string]string{"area_a": edakit.IntStr(sa), "best": edakit.IntStr(best)})
		}
	default:
		n := 2 * edakit.Scaled(ctx, 3, 10)
		v := binom(n, n/2) / 2
		return edakit.Numeric(edakit.Sprintf("How many distinct balanced bisections do %d unit-area cells have, splitting them into two unlabeled halves of %d cells each?", n, n/2),
			edakit.IntStr(v),
			edakit.Sprintf("There are C(%d, %d) ways to pick the cells of one labelled half; dividing by 2 because swapping the sides gives the same bisection leaves %d. The count grows like $2^n/\\sqrt{n}$, which is why exhaustive search is hopeless.", n, n/2, v),
			map[string]string{"bisections": edakit.IntStr(v)})
	}
}

func klD(g wgraph, side []int, v int) (int, int) {
	E, I := 0, 0
	for u := 0; u < g.n; u++ {
		if g.w[v][u] == 0 {
			continue
		}
		if side[u] == side[v] {
			I += g.w[v][u]
		} else {
			E += g.w[v][u]
		}
	}
	return E, I
}

func klGainGen(ctx generator.GeneratorContext) generator.Problem {
	n := 2 * edakit.Scaled(ctx, 3, 4)
	g := randWGraph(n, edakit.Scaled(ctx, 3, 7), 5)
	side := balancedSides(n)
	intro := "A weighted graph has edges " + edgeText(g, true) + ". The bisection is A = " + vSideText(side, 0) + ", B = " + vSideText(side, 1) + ". For a vertex v, D(v) = E(v) - I(v), the external minus internal edge weight. "
	pickIn := func(s int) int {
		var xs []int
		for i := range side {
			if side[i] == s {
				xs = append(xs, i)
			}
		}
		return edakit.Pick(xs)
	}
	switch kpOf(ctx) {
	case 1:
		v := edakit.Int(0, n-1)
		E, I := klD(g, side, v)
		return edakit.Numeric(intro+"What is D("+vName(v)+")?", edakit.IntStr(E-I),
			edakit.Sprintf("Edges from %s to the other side weigh %d and edges inside its side weigh %d, so D = %d - %d = %d.", vName(v), E, I, E, I, E-I),
			map[string]string{"E": edakit.IntStr(E), "I": edakit.IntStr(I)})
	case 2:
		a, b := pickIn(0), pickIn(1)
		Ea, Ia := klD(g, side, a)
		Eb, Ib := klD(g, side, b)
		gain := Ea - Ia + Eb - Ib - 2*g.w[a][b]
		return edakit.Numeric(intro+"What is the Kernighan-Lin gain of swapping "+vName(a)+" (in A) with "+vName(b)+" (in B)?", edakit.IntStr(gain),
			edakit.Sprintf("g = D(a) + D(b) - 2c(a,b) = %d + %d - 2 × %d = %d. The edge between the two stays cut after the swap, so it is subtracted twice.", Ea-Ia, Eb-Ib, g.w[a][b], gain),
			map[string]string{"Da": edakit.IntStr(Ea - Ia), "Db": edakit.IntStr(Eb - Ib), "gain": edakit.IntStr(gain)})
	default:
		best := math.MinInt32
		for a := 0; a < n; a++ {
			for b := 0; b < n; b++ {
				if side[a] == 0 && side[b] == 1 {
					Ea, Ia := klD(g, side, a)
					Eb, Ib := klD(g, side, b)
					best = max(best, Ea-Ia+Eb-Ib-2*g.w[a][b])
				}
			}
		}
		var dt []string
		for v := 0; v < n; v++ {
			E, I := klD(g, side, v)
			dt = append(dt, edakit.Sprintf("D(%s) = %d", vName(v), E-I))
		}
		return edakit.Numeric(intro+"Over all pairs (a in A, b in B), what is the largest swap gain?", edakit.IntStr(best),
			"The D values are "+strings.Join(dt, ", ")+". Checking g = D(a) + D(b) - 2c(a,b) for every pair, the maximum is "+edakit.IntStr(best)+".",
			map[string]string{"best": edakit.IntStr(best)})
	}
}

func fmGainGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 9)
	h := randHyper(n, edakit.Scaled(ctx, 5, 9), 4)
	side := randSides(n)
	intro := "Nets: " + netText(h, false) + ". The partition is A = " + sideText(side, 0) + ", B = " + sideText(side, 1) + ". "
	switch kpOf(ctx) {
	case 1:
		c := edakit.Int(0, n-1)
		_, fs, _ := fmGain(h, side, c)
		return edakit.Numeric(intro+"FS("+cellName(c)+") counts the nets on which "+cellName(c)+" is the only cell on its side. What is FS("+cellName(c)+")?",
			edakit.IntStr(fs),
			"On those nets every other pin is already across, so moving "+cellName(c)+" makes the net uncut. "+cellName(c)+" has "+edakit.IntStr(fs)+" such nets.",
			map[string]string{"fs": edakit.IntStr(fs)})
	case 2:
		c := edakit.Int(0, n-1)
		g, fs, te := fmGain(h, side, c)
		return edakit.Numeric(intro+"What is the FM gain of moving "+cellName(c)+" to the other side (the decrease in the number of cut nets)?",
			edakit.IntStr(g),
			edakit.Sprintf("FS = %d nets become uncut and TE = %d nets lying wholly on %s's side become cut, so the gain is %d - %d = %d.", fs, te, cellName(c), fs, te, g),
			map[string]string{"fs": edakit.IntStr(fs), "te": edakit.IntStr(te), "gain": edakit.IntStr(g)})
	default:
		best := math.MinInt32
		var gt []string
		for c := 0; c < n; c++ {
			g, _, _ := fmGain(h, side, c)
			gt = append(gt, edakit.Sprintf("%s:%d", cellName(c), g))
			best = max(best, g)
		}
		return edakit.Numeric(intro+"Ignoring balance, what is the largest FM gain over all cells?", edakit.IntStr(best),
			"Gains FS - TE per cell: "+strings.Join(gt, ", ")+". The maximum is "+edakit.IntStr(best)+".",
			map[string]string{"best": edakit.IntStr(best)})
	}
}

func fmBucketGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 9)
	h := randHyper(n, edakit.Scaled(ctx, 5, 9), 4)
	side := randSides(n)
	intro := "Nets: " + netText(h, false) + ". "
	switch kpOf(ctx) {
	case 1:
		deg := make([]int, n)
		for _, net := range h.nets {
			for _, p := range net {
				deg[p]++
			}
		}
		pmax := 0
		for _, d := range deg {
			pmax = max(pmax, d)
		}
		return edakit.Numeric(intro+"An FM bucket array keeps one list per possible gain value, and a cell's gain lies between -p and p where p is its number of nets. How many buckets does one side's array need?",
			edakit.IntStr(2*pmax+1),
			edakit.Sprintf("The largest cell degree is p_max = %d (degrees %s), so gains run from -%d to %d: 2 × %d + 1 = %d buckets.", pmax, intsText(deg), pmax, pmax, pmax, 2*pmax+1),
			map[string]string{"pmax": edakit.IntStr(pmax)})
	default:
		for tries := 0; ; tries++ {
			c := edakit.Int(0, n-1)
			after := append([]int(nil), side...)
			after[c] = 1 - after[c]
			changed := 0
			var nbs []int
			for d := 0; d < n; d++ {
				if d == c {
					continue
				}
				g0, _, _ := fmGain(h, side, d)
				g1, _, _ := fmGain(h, after, d)
				if g0 != g1 {
					changed++
					nbs = append(nbs, d)
				}
			}
			if len(nbs) == 0 && tries < 20 {
				continue
			}
			q := intro + "The partition is A = " + sideText(side, 0) + ", B = " + sideText(side, 1) + ". FM moves " + cellName(c) + " to the other side and locks it. "
			if kpOf(ctx) == 2 && len(nbs) > 0 {
				d := edakit.Pick(nbs)
				if edakit.Int(0, 2) == 0 {
					d = edakit.Int(0, n-1)
					for d == c {
						d = edakit.Int(0, n-1)
					}
				}
				g0, _, _ := fmGain(h, side, d)
				g1, _, _ := fmGain(h, after, d)
				return edakit.Numeric(q+"What is the gain of "+cellName(d)+" after the move?", edakit.IntStr(g1),
					edakit.Sprintf("Before the move %s had gain %d. Recounting FS - TE on the new partition gives %d. Only nets that contain both %s and %s can change it.", cellName(d), g0, g1, cellName(c), cellName(d)),
					map[string]string{"before": edakit.IntStr(g0), "after": edakit.IntStr(g1)})
			}
			var names []string
			for _, d := range nbs {
				names = append(names, cellName(d))
			}
			return edakit.Numeric(q+"How many of the other cells have a different gain after the move?", edakit.IntStr(changed),
				"Recomputing every gain shows changes at {"+strings.Join(names, ", ")+"}: "+edakit.IntStr(changed)+" cells. Only neighbours on critical nets (nets with zero or one pin on a side before or after the move) can change, which is why FM updates touch so few cells.",
				map[string]string{"changed": edakit.IntStr(changed)})
		}
	}
}

// fmPass is the computed part of the FM algorithm: committing the best prefix of a move sequence.
func fmPass(ctx generator.GeneratorContext) generator.Problem {
	for {
		k := edakit.Scaled(ctx, 5, 8)
		c0 := edakit.Int(8, 20)
		gs := make([]int, k)
		for i := range gs {
			gs[i] = edakit.Int(-3, 3)
		}
		cum, best, arg, unique := 0, 0, 0, true
		cums := []int{0}
		for i, g := range gs {
			cum += g
			cums = append(cums, cum)
			if cum > best {
				best, arg, unique = cum, i+1, true
			} else if cum == best && best > 0 {
				unique = false
			}
		}
		if best <= 0 || !unique || arg == k {
			continue
		}
		opt := func(j int) string {
			if j == 1 {
				return edakit.Sprintf("Keep only the first move, cut %d", c0-cums[1])
			}
			return edakit.Sprintf("Keep the first %d moves, cut %d", j, c0-cums[j])
		}
		// Greedy stop at the first negative gain.
		stop := 0
		for stop < k && gs[stop] >= 0 {
			stop++
		}
		var ds []string
		for _, j := range []int{k, stop, arg - 1, arg + 1, 1} {
			if j != arg && j >= 1 && j <= k {
				ds = append(ds, opt(j))
			}
		}
		if len(ds) < 3 {
			continue
		}
		return edakit.Choice(edakit.Sprintf("An FM pass starts from a cut of %d nets and tentatively moves cells (each then locked) with gains %s, in that order. Which outcome does the pass commit?", c0, intsText(gs)),
			opt(arg),
			edakit.Sprintf("The cumulative gains are %s. FM rolls back to the prefix with the largest cumulative gain, %d after %d moves, giving cut %d - %d = %d. Stopping at the first negative gain would miss later recoveries.", intsText(cums[1:]), best, arg, c0, best, c0-best),
			ds...)
	}
}

func heavyEdgeMatching(g wgraph, order []int) ([]int, int) {
	mate := make([]int, g.n)
	for i := range mate {
		mate[i] = -1
	}
	mw := 0
	for _, v := range order {
		if mate[v] >= 0 {
			continue
		}
		best := -1
		for u := 0; u < g.n; u++ {
			if u != v && mate[u] < 0 && g.w[v][u] > 0 && (best < 0 || g.w[v][u] > g.w[v][best]) {
				best = u
			}
		}
		if best >= 0 {
			mate[v], mate[best] = best, v
			mw += g.w[v][best]
		}
	}
	return mate, mw
}

func coarsenMatchingGen(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 9)
	g := randWGraph(n, edakit.Scaled(ctx, 3, 7), 9)
	order := edakit.Shuffle(seq(n))
	mate, mw := heavyEdgeMatching(g, order)
	var ot []string
	for _, v := range order {
		ot = append(ot, vName(v))
	}
	matched := 0
	for _, m := range mate {
		if m >= 0 {
			matched++
		}
	}
	var pairs []string
	for v := 0; v < n; v++ {
		if mate[v] > v {
			pairs = append(pairs, vName(v)+"-"+vName(mate[v]))
		}
	}
	intro := "A weighted graph has edges " + edgeText(g, true) + ". Heavy-edge matching visits vertices in the order " + strings.Join(ot, ", ") + "; an unmatched vertex is matched to its unmatched neighbour with the heaviest edge (ties go to the smaller number), or stays single if none is left. "
	expl := "The matching is {" + strings.Join(pairs, ", ") + "}."
	switch kpOf(ctx) {
	case 1:
		var cands []int
		for v := 0; v < n; v++ {
			if mate[v] >= 0 {
				cands = append(cands, v)
			}
		}
		v := edakit.Pick(cands)
		return edakit.Numeric(intro+"Which vertex is "+vName(v)+" matched with? Answer with its number.", edakit.IntStr(mate[v]+1), expl, map[string]string{"mate": edakit.IntStr(mate[v] + 1)})
	case 2:
		coarse := n - matched/2
		return edakit.Numeric(intro+"Contracting every matched pair into one vertex, how many vertices does the coarse graph have?", edakit.IntStr(coarse),
			expl+edakit.Sprintf(" %d pairs merge, so %d - %d = %d vertices remain.", matched/2, n, matched/2, coarse),
			map[string]string{"coarse": edakit.IntStr(coarse)})
	default:
		total := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				total += g.w[i][j]
			}
		}
		return edakit.Numeric(intro+"After contraction (parallel edges merged by adding weights, matched edges disappear inside the new vertices), what is the total edge weight of the coarse graph?",
			edakit.IntStr(total-mw),
			expl+edakit.Sprintf(" The original total weight is %d and the matched edges carry %d, so %d remains exposed. Hiding heavy edges inside clusters is the point: no partition of the coarse graph can cut them.", total, mw, total-mw),
			map[string]string{"total": edakit.IntStr(total), "matched": edakit.IntStr(mw)})
	}
}

// multilevelLevels is the computed part of multilevel partitioning: how many coarsening levels.
func multilevelLevels(ctx generator.GeneratorContext) generator.Problem {
	N := edakit.Pick([]int{10000, 50000, 200000, 1000000, 4000000})
	r := edakit.Pick([]float64{1.6, 1.8, 2.0})
	T := edakit.Pick([]int{100, 150, 200, 500})
	L := 0
	x := float64(N)
	for x > float64(T) {
		x /= r
		L++
	}
	opt := func(l int) string { return edakit.Sprintf("%d levels", l) }
	ds := []string{opt(L - 1), opt(L + 1), opt(L + 3)}
	if L-2 > 0 {
		ds = append(ds, opt(L-2))
	}
	return edakit.Choice(edakit.Sprintf("A multilevel partitioner coarsens a netlist of %d cells. Each level shrinks the vertex count by a factor of %s, and coarsening stops at the first level with at most %d vertices. How many coarsening levels run?", N, edakit.Num(r, 1), T),
		opt(L),
		edakit.Sprintf("The smallest L with %d / %s^L ≤ %d is L = %d. The hierarchy depth grows only logarithmically in the netlist size, and each level costs time linear in its size, so the whole V-cycle stays near-linear.", N, edakit.Num(r, 1), T, L),
		ds...)
}
