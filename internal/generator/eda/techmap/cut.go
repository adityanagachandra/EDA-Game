package techmap

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// lnet is a small logic network: signals 0..npi-1 are primary inputs, the rest are gates whose
// fanins are earlier signals. comp marks complemented fanin edges (AIG use only).
type lnet struct {
	npi  int
	fan  [][]int
	comp [][]bool
}

func (g *lnet) size() int { return len(g.fan) }

func (g *lnet) name(i int) string {
	if i < g.npi {
		return piName(i)
	}
	return edakit.Sprintf("n%d", i-g.npi+1)
}

// randNet draws ng gates, each with 2..maxFan distinct fanins biased toward recent signals so
// the network has depth and reconvergence.
func randNet(npi, ng, maxFan int, complements bool) *lnet {
	g := &lnet{npi: npi, fan: make([][]int, npi+ng), comp: make([][]bool, npi+ng)}
	for x := npi; x < npi+ng; x++ {
		nf := 2
		if maxFan > 2 && edakit.Int(1, 10) <= 4 {
			nf = edakit.Int(3, maxFan)
		}
		if nf > x {
			nf = x
		}
		used := map[int]bool{}
		for len(g.fan[x]) < nf {
			var y int
			if x > npi && edakit.Int(1, 10) <= 6 {
				lo := x - 3
				if lo < 0 {
					lo = 0
				}
				y = edakit.Int(lo, x-1)
			} else {
				y = edakit.Int(0, x-1)
			}
			if used[y] {
				continue
			}
			used[y] = true
			g.fan[x] = append(g.fan[x], y)
			g.comp[x] = append(g.comp[x], complements && edakit.Int(1, 10) <= 4)
		}
	}
	return g
}

// cone restricts g to the transitive fanin of root and renumbers it.
func (g *lnet) cone(root int) *lnet {
	in := map[int]bool{}
	var walk func(x int)
	walk = func(x int) {
		if in[x] {
			return
		}
		in[x] = true
		for _, f := range g.fan[x] {
			walk(f)
		}
	}
	walk(root)
	idx := map[int]int{}
	out := &lnet{}
	for x := 0; x < g.npi; x++ {
		if in[x] {
			idx[x] = len(out.fan)
			out.fan = append(out.fan, nil)
			out.comp = append(out.comp, nil)
		}
	}
	out.npi = len(out.fan)
	for x := g.npi; x < g.size(); x++ {
		if in[x] {
			idx[x] = len(out.fan)
			var f []int
			for _, y := range g.fan[x] {
				f = append(f, idx[y])
			}
			out.fan = append(out.fan, f)
			out.comp = append(out.comp, append([]bool{}, g.comp[x]...))
		}
	}
	return out
}

// randCone draws a network and returns the cone of its last gate with at least minGates gates.
func randCone(npi, ng, maxFan, minGates int, complements bool) *lnet {
	minPI := 3
	if npi < 3 {
		minPI = npi
	}
	for {
		c := randNet(npi, ng, maxFan, complements).cone(npi + ng - 1)
		if c.size()-c.npi >= minGates && c.npi >= minPI {
			return c
		}
	}
}

// describe prints the gates. AIG style shows AND with complemented inputs; generic style lists
// the fanins of each gate.
func (g *lnet) describe(aig bool) string {
	var lines []string
	for x := g.npi; x < g.size(); x++ {
		var args []string
		for i, f := range g.fan[x] {
			s := g.name(f)
			if g.comp[x][i] {
				s = "NOT " + s
			}
			args = append(args, s)
		}
		if aig {
			lines = append(lines, g.name(x)+" = AND("+strings.Join(args, ", ")+")")
		} else {
			lines = append(lines, g.name(x)+" has fanins "+strings.Join(args, ", "))
		}
	}
	return strings.Join(lines, "; ")
}

func (g *lnet) inputs() string {
	var s []string
	for i := 0; i < g.npi; i++ {
		s = append(s, g.name(i))
	}
	return strings.Join(s, ", ")
}

// enumCuts computes the k-feasible cuts of every signal by bottom-up merging, with the trivial
// cut included and dominated cuts (strict supersets of another cut) removed.
func enumCuts(g *lnet, k int) [][]uint32 {
	cuts := make([][]uint32, g.size())
	for x := 0; x < g.size(); x++ {
		if x < g.npi {
			cuts[x] = []uint32{1 << uint(x)}
			continue
		}
		acc := []uint32{0}
		for _, f := range g.fan[x] {
			var nxt []uint32
			seen := map[uint32]bool{}
			for _, a := range acc {
				for _, b := range cuts[f] {
					u := a | b
					if popcount(u) <= k && !seen[u] {
						seen[u] = true
						nxt = append(nxt, u)
					}
				}
			}
			acc = nxt
		}
		acc = append(acc, 1<<uint(x))
		cuts[x] = pruneDominated(acc)
	}
	return cuts
}

func pruneDominated(cs []uint32) []uint32 {
	var out []uint32
	for i, c := range cs {
		dom := false
		for j, d := range cs {
			if i != j && d != c && d&c == d {
				dom = true
				break
			}
		}
		if !dom {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if popcount(out[i]) != popcount(out[j]) {
			return popcount(out[i]) < popcount(out[j])
		}
		return out[i] < out[j]
	})
	return out
}

func (g *lnet) cutString(m uint32) string {
	var s []string
	for x := 0; x < g.size(); x++ {
		if m&(1<<uint(x)) != 0 {
			s = append(s, g.name(x))
		}
	}
	return "{" + strings.Join(s, ", ") + "}"
}

func (g *lnet) cutList(cs []uint32) string {
	var s []string
	for _, c := range cs {
		s = append(s, g.cutString(c))
	}
	return strings.Join(s, ", ")
}

func leaves(m uint32, n int) []int {
	var out []int
	for x := 0; x < n; x++ {
		if m&(1<<uint(x)) != 0 {
			out = append(out, x)
		}
	}
	return out
}

// depthLabels returns the unit-delay depth-optimal label of every signal and its best cut
// (smallest depth, then fewer leaves, then smaller mask).
func depthLabels(g *lnet, cuts [][]uint32) ([]int, []uint32) {
	lab := make([]int, g.size())
	best := make([]uint32, g.size())
	for x := g.npi; x < g.size(); x++ {
		lab[x] = 1 << 30
		for _, c := range cuts[x] {
			if c == 1<<uint(x) {
				continue
			}
			d := 0
			for _, l := range leaves(c, g.size()) {
				if lab[l] > d {
					d = lab[l]
				}
			}
			d++
			if d < lab[x] || (d == lab[x] && (popcount(c) < popcount(best[x]) || (popcount(c) == popcount(best[x]) && c < best[x]))) {
				lab[x], best[x] = d, c
			}
		}
	}
	return lab, best
}

func kFeasible(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	rules := "Use the standard bottom-up enumeration: a primary input has only its trivial cut; the cuts of a node are its trivial cut plus every union of one cut from each fanin, keeping unions with at most K leaves and discarding any cut that strictly contains another cut of the same node."
	for {
		var g *lnet
		K := 3
		switch k {
		case 1:
			g = randCone(3, edakit.Int(3, 4), 2, 3, false)
		case 2:
			g = randCone(edakit.Int(3, 4), edakit.Scaled(ctx, 4, 6), 2, 4, false)
			K = edakit.Int(3, 4)
		default:
			g = randCone(4, edakit.Int(4, 5), 2, 4, false)
		}
		if k == 3 && g.npi < 4 {
			continue
		}
		root := g.size() - 1
		cuts := enumCuts(g, K)
		n := len(cuts[root])
		q := "An AIG over inputs " + g.inputs() + " is " + g.describe(true) + ". " + rules
		if k == 3 {
			c4 := enumCuts(g, 4)
			n3, n4 := len(enumCuts(g, 3)[root]), len(c4[root])
			if n4 > 18 || (n4 == n3 && edakit.Int(0, 3) > 0) {
				continue
			}
			return edakit.Numeric(q+" How many more cuts does node "+g.name(root)+" have with K = 4 than with K = 3?",
				edakit.IntStr(n4-n3),
				edakit.Sprintf("With K = 3, %s has %d cuts: %s. With K = 4 it has %d: %s. The difference is %d.", g.name(root), n3, g.cutList(enumCuts(g, 3)[root]), n4, g.cutList(c4[root]), n4-n3),
				map[string]string{"cuts_k3": edakit.IntStr(n3), "cuts_k4": edakit.IntStr(n4)})
		}
		if (k == 1 && n > 10) || n > 15 {
			continue
		}
		return edakit.Numeric(q+edakit.Sprintf(" With K = %d, how many cuts does node %s have, including its trivial cut?", K, g.name(root)),
			edakit.IntStr(n),
			edakit.Sprintf("Merging bottom-up, %s ends with %d cuts: %s.", g.name(root), n, g.cutList(cuts[root])),
			map[string]string{"cuts": edakit.IntStr(n)})
	}
}

func cutEnumeration(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	if k == 3 {
		return cutFunction(ctx)
	}
	for {
		K := edakit.Int(3, 4)
		g := randCone(edakit.Int(3, 4), edakit.Scaled(ctx, 4, 6), 2, 4, false)
		root := g.size() - 1
		u, w := g.fan[root][0], g.fan[root][1]
		if u < g.npi && w < g.npi {
			continue
		}
		cuts := enumCuts(g, K)
		seen := map[uint32]bool{}
		for _, a := range cuts[u] {
			for _, b := range cuts[w] {
				if m := a | b; popcount(m) <= K {
					seen[m] = true
				}
			}
		}
		final := len(cuts[root])
		if len(cuts[u])*len(cuts[w]) > 16 {
			continue
		}
		q := edakit.Sprintf("Node %s = AND(%s, %s) in an AIG. Its fanins have these cut sets, trivial cuts included: cuts(%s) = %s; cuts(%s) = %s. The cut-size limit is K = %d.",
			g.name(root), g.name(u), g.name(w), g.name(u), g.cutList(cuts[u]), g.name(w), g.cutList(cuts[w]), K)
		facts := map[string]string{"pairs": edakit.IntStr(len(cuts[u]) * len(cuts[w])), "feasible_unions": edakit.IntStr(len(seen)), "final": edakit.IntStr(final)}
		if k == 1 {
			return edakit.Numeric(q+edakit.Sprintf(" Form the union of every pair (one cut from each fanin). How many distinct unions have at most %d leaves?", K),
				edakit.IntStr(len(seen)),
				edakit.Sprintf("There are %d x %d = %d pairs. Their unions, with duplicates merged and those over %d leaves dropped, number %d.", len(cuts[u]), len(cuts[w]), len(cuts[u])*len(cuts[w]), K, len(seen)),
				facts)
		}
		return edakit.Numeric(q+" After merging, add the trivial cut and discard every cut that strictly contains another cut of the node. How many cuts does "+g.name(root)+" keep?",
			edakit.IntStr(final),
			edakit.Sprintf("The %d feasible unions plus the trivial cut, with dominated cuts removed, leave %d: %s.", len(seen), final, g.cutList(cuts[root])),
			facts)
	}
}

// evalCut returns the truth table (as a ones count and a list) of root over the cut leaves.
func evalCut(g *lnet, root int, lv []int) (ones int, minterms []int) {
	pos := map[int]int{}
	for i, l := range lv {
		pos[l] = i
	}
	for m := 0; m < 1<<uint(len(lv)); m++ {
		memo := map[int]bool{}
		var val func(x int) bool
		val = func(x int) bool {
			if p, ok := pos[x]; ok {
				return m&(1<<uint(p)) != 0
			}
			if r, ok := memo[x]; ok {
				return r
			}
			r := true
			for i, f := range g.fan[x] {
				b := val(f)
				if g.comp[x][i] {
					b = !b
				}
				r = r && b
			}
			memo[x] = r
			return r
		}
		if val(root) {
			ones++
			minterms = append(minterms, m)
		}
	}
	return
}

func cutFunction(ctx generator.GeneratorContext) generator.Problem {
	for {
		g := randCone(edakit.Int(3, 4), edakit.Scaled(ctx, 3, 5), 2, 3, true)
		root := g.size() - 1
		var opts []uint32
		for _, c := range enumCuts(g, 3)[root] {
			if popcount(c) == 3 {
				opts = append(opts, c)
			}
		}
		if len(opts) == 0 {
			continue
		}
		c := edakit.Pick(opts)
		lv := leaves(c, g.size())
		ones, mins := evalCut(g, root, lv)
		q := "An AIG with complemented edges over inputs " + g.inputs() + " is " + g.describe(true) + ". " + g.cutString(c) +
			" is a 3-feasible cut of " + g.name(root) + ". How many of the 8 assignments to the cut leaves make " + g.name(root) + " equal to 1?"
		var order []string
		for _, l := range lv {
			order = append(order, g.name(l))
		}
		return edakit.Numeric(q, edakit.IntStr(ones),
			edakit.Sprintf("Treat the leaves %s as free variables and evaluate the cone above them for all 8 assignments. With the first leaf as the least significant bit, the function is 1 at minterms {%s}, %d assignments in all.", strings.Join(order, ", "), joinInts(mins), ones),
			map[string]string{"ones": edakit.IntStr(ones)})
	}
}

func priorityCuts(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		pool := []string{"a", "b", "c", "d", "n1", "n2", "n3", "n4", "n5"}
		depth := map[string]int{"a": 0, "b": 0, "c": 0, "d": 0}
		for _, s := range pool[4:] {
			depth[s] = edakit.Int(1, 4)
		}
		nc := edakit.Int(4, 5)
		var desc []string
		best := 99
		var all []int
		for i := 0; i < nc; i++ {
			lv := edakit.Shuffle(pool)[:edakit.Int(2, 4)]
			sort.Strings(lv)
			mx := 0
			for _, l := range lv {
				if depth[l] > mx {
					mx = depth[l]
				}
			}
			all = append(all, mx+1)
			if mx+1 < best {
				best = mx + 1
			}
			desc = append(desc, edakit.Sprintf("C%d = {%s}", i+1, strings.Join(lv, ", ")))
		}
		var dl []string
		for _, s := range pool[4:] {
			dl = append(dl, edakit.Sprintf("%s: %d", s, depth[s]))
		}
		q := "A cut-based LUT mapper is choosing among the candidate cuts of node t: " + strings.Join(desc, "; ") + ". Inputs a to d have depth 0 and the best depths of the internal nodes are " + strings.Join(dl, ", ") + ". Using unit LUT delay, what is the smallest depth t can reach with these candidates?"
		return edakit.Choice(q, edakit.IntStr(best),
			edakit.Sprintf("The depth of a cut is one more than the largest depth among its leaves. The candidates give %s, so the best is %d.", joinInts(all), best),
			edakit.Distinct(best, 3, best-1, best+1, best+2)...)
	case 2:
		return priorityKeep()
	}
	C := edakit.Int(3, 20)
	r := edakit.Int(2, 3)
	n := 1
	for i := 0; i < r; i++ {
		n *= C + 1
	}
	q := edakit.Sprintf("A priority-cut mapper stores C = %d cuts per node, plus the trivial cut, and merges cut sets of the fanins of a %d-input node by trying one stored cut from each fanin. How many candidate unions does it evaluate at that node?", C, r)
	return edakit.Choice(q, edakit.IntStr(n),
		edakit.Sprintf("Each fanin offers C + 1 = %d cuts, so there are %d^%d = %d combinations, independent of how many cuts exhaustive enumeration would have produced.", C+1, C+1, r, n),
		edakit.Distinct(n, 3, intPow(C, r), C*r, (C+1)*r, intPow(C, r)+1)...)
}

func intPow(b, e int) int {
	n := 1
	for i := 0; i < e; i++ {
		n *= b
	}
	return n
}

func priorityKeep() generator.Problem {
	for {
		nc := edakit.Int(5, 6)
		C := edakit.Int(2, 3)
		type cand struct {
			name  string
			depth int
			af    float64
		}
		var cs []cand
		used := map[int]bool{}
		for i := 0; i < nc; i++ {
			var a int
			for {
				a = edakit.Int(10, 60)
				if !used[a] {
					used[a] = true
					break
				}
			}
			cs = append(cs, cand{edakit.Sprintf("C%d", i+1), edakit.Int(2, 4), float64(a) / 10})
		}
		pickBy := func(less func(a, b cand) bool) string {
			s := append([]cand{}, cs...)
			sort.SliceStable(s, func(i, j int) bool { return less(s[i], s[j]) })
			var names []string
			for _, c := range s[:C] {
				names = append(names, c.name)
			}
			sort.Strings(names)
			return strings.Join(names, ", ")
		}
		byDelay := pickBy(func(a, b cand) bool {
			if a.depth != b.depth {
				return a.depth < b.depth
			}
			return a.af < b.af
		})
		byArea := pickBy(func(a, b cand) bool {
			if a.af != b.af {
				return a.af < b.af
			}
			return a.depth < b.depth
		})
		worst := pickBy(func(a, b cand) bool {
			if a.depth != b.depth {
				return a.depth > b.depth
			}
			return a.af > b.af
		})
		first := pickBy(func(a, b cand) bool { return a.name < b.name })
		mode := edakit.Pick([]string{"delay", "area"})
		correct, other := byDelay, byArea
		rule := "sorted by depth, ties broken by smaller area flow"
		if mode == "area" {
			correct, other = byArea, byDelay
			rule = "sorted by area flow, ties broken by smaller depth"
		}
		ds := distinctStrings([]string{other, worst, first})
		var dis []string
		for _, d := range ds {
			if d != correct {
				dis = append(dis, d)
			}
		}
		for tries := 0; len(dis) < 3 && tries < 50; tries++ {
			sh := edakit.Shuffle(cs)[:C]
			var names []string
			for _, c := range sh {
				names = append(names, c.name)
			}
			sort.Strings(names)
			s := strings.Join(names, ", ")
			dup := s == correct
			for _, d := range dis {
				if d == s {
					dup = true
				}
			}
			if !dup {
				dis = append(dis, s)
			}
		}
		if len(dis) < 3 || (correct == other && edakit.Int(0, 2) > 0) {
			continue
		}
		var desc []string
		for _, c := range cs {
			desc = append(desc, edakit.Sprintf("%s (depth %d, area flow %s)", c.name, c.depth, edakit.Num(c.af, 1)))
		}
		q := edakit.Sprintf("At node t a priority-cut mapper in %s-oriented mode has these candidates: %s. It keeps the best C = %d, %s. Which cuts does it keep?", mode, strings.Join(desc, "; "), C, rule)
		return edakit.Choice(q, correct,
			"With the candidates "+rule+", the first "+edakit.IntStr(C)+" are "+correct+". The other mode's ordering would keep "+other+".", dis...)
	}
}

func areaFlow(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	lutArea := edakit.Pick([]int{1, 1, 2, 3})
	afs := func(n int) []float64 {
		var out []float64
		for i := 0; i < n; i++ {
			if edakit.Int(1, 4) == 1 {
				out = append(out, 0)
			} else {
				out = append(out, float64(edakit.Int(2, 30))/10)
			}
		}
		return out
	}
	fmtLeaves := func(names []string, af []float64) string {
		var s []string
		for i := range names {
			s = append(s, edakit.Sprintf("%s with area flow %s", names[i], edakit.Num(af[i], 1)))
		}
		return strings.Join(s, ", ")
	}
	sum := func(xs []float64) float64 {
		t := 0.0
		for _, x := range xs {
			t += x
		}
		return t
	}
	def := edakit.Sprintf("Area flow is AF(n) = (A + sum of AF over the cut leaves) / F(n), where A = %d is the area of the gate or LUT that implements the cut and F(n) is the estimated fanout of n. Primary inputs have area flow 0.", lutArea)
	switch k {
	case 1:
		nl := edakit.Int(2, 4)
		names := []string{"p", "q", "r", "s"}[:nl]
		af := afs(nl)
		F := edakit.Int(1, 4)
		ans := (float64(lutArea) + sum(af)) / float64(F)
		q := def + edakit.Sprintf(" Node n has fanout %d and its cut has leaves %s. What is AF(n) for this cut? Round to 2 decimal places.", F, fmtLeaves(names, af))
		return edakit.Numeric(q, edakit.Num(ans, 2),
			edakit.Sprintf("(%d + %s) / %d = %s.", lutArea, edakit.Num(sum(af), 1), F, edakit.Num(ans, 2)),
			map[string]string{"area_flow": edakit.Num(ans, 2)})
	case 2:
		F := edakit.Int(1, 3)
		names := []string{"p", "q", "r", "s", "u"}
		leafAF := afs(len(names))
		var ld []string
		for i, nm := range names {
			ld = append(ld, edakit.Sprintf("%s: %s", nm, edakit.Num(leafAF[i], 1)))
		}
		var lines []string
		best := 1e9
		var vals []string
		for c := 0; c < 3; c++ {
			idx := edakit.Shuffle([]int{0, 1, 2, 3, 4})[:edakit.Int(2, 4)]
			sort.Ints(idx)
			var lv []string
			t := 0.0
			for _, i := range idx {
				lv = append(lv, names[i])
				t += leafAF[i]
			}
			v := (float64(lutArea) + t) / float64(F)
			if v < best {
				best = v
			}
			vals = append(vals, edakit.Num(v, 2))
			lines = append(lines, edakit.Sprintf("cut %d = {%s}", c+1, strings.Join(lv, ", ")))
		}
		q := def + edakit.Sprintf(" Node n has fanout %d. The leaf area flows are %s. The candidate cuts are %s. What is the smallest area flow n can get from these cuts? Round to 2 decimal places.", F, strings.Join(ld, ", "), strings.Join(lines, "; "))
		return edakit.Numeric(q, edakit.Num(best, 2),
			"The three cuts give area flows "+strings.Join(vals, ", ")+", so the best is "+edakit.Num(best, 2)+".",
			map[string]string{"best_area_flow": edakit.Num(best, 2)})
	}
	// Two-level computation: x and y over inputs, root over x, y and an input.
	fx, fy := edakit.Int(1, 3), edakit.Int(1, 3)
	afx := float64(lutArea) / float64(fx)
	afy := (float64(lutArea) + afx) / float64(fy)
	useX := edakit.Int(0, 1) == 1
	fr := edakit.Int(1, 2)
	var root float64
	var rootCut string
	if useX {
		root = (float64(lutArea) + afx + afy) / float64(fr)
		rootCut = "{x, y, c}"
	} else {
		root = (float64(lutArea) + afy) / float64(fr)
		rootCut = "{y, c, d}"
	}
	q := def + edakit.Sprintf(" Node x uses cut {a, b} and has fanout %d. Node y uses cut {x, c} and has fanout %d. Node n uses cut %s and has fanout %d. All of a, b, c, d are primary inputs. What is AF(n)? Round to 2 decimal places.", fx, fy, rootCut, fr)
	return edakit.Numeric(q, edakit.Num(root, 2),
		edakit.Sprintf("AF(x) = %d/%d = %s. AF(y) = (%d + AF(x))/%d = %s. AF(n) = %s.", lutArea, fx, edakit.Num(afx, 3), lutArea, fy, edakit.Num(afy, 3), edakit.Num(root, 2)),
		map[string]string{"af_x": edakit.Num(afx, 3), "af_y": edakit.Num(afy, 3), "af_n": edakit.Num(root, 2)})
}

type eaNode struct {
	area int
	cut  []string
	refs int
}

func exactArea(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	for {
		nodes := map[string]*eaNode{}
		var order []string
		// Second-level nodes over inputs and first-level nodes.
		pis := []string{"a", "b", "c", "d", "e"}
		lvl1 := []string{"u", "w"}
		for _, n := range lvl1 {
			cut := edakit.Shuffle(pis)[:2]
			sort.Strings(cut)
			refs := 0
			if edakit.Int(1, 2) == 1 {
				refs = edakit.Int(1, 2)
			}
			nodes[n] = &eaNode{area: edakit.Int(2, 4), cut: cut, refs: refs}
			order = append(order, n)
		}
		lvl2 := []string{"p", "q", "r"}
		for _, n := range lvl2 {
			pool := append(append([]string{}, pis...), lvl1...)
			cut := edakit.Shuffle(pool)[:2]
			sort.Strings(cut)
			refs := 0
			if edakit.Int(1, 3) == 1 {
				refs = edakit.Int(1, 2)
			}
			nodes[n] = &eaNode{area: edakit.Int(2, 5), cut: cut, refs: refs}
			order = append(order, n)
		}
		nl := 2 + edakit.Int(0, 1)
		if k == 1 {
			nl = 2
		}
		cand := edakit.Shuffle(append(append([]string{}, lvl2...), "c", "d"))[:nl]
		sort.Strings(cand)
		gate := edakit.Int(3, 6)
		refs := map[string]int{}
		for n, nd := range nodes {
			refs[n] = nd.refs
		}
		var trace []string
		shared := false
		var ref func(cut []string) int
		ref = func(cut []string) int {
			t := 0
			for _, l := range cut {
				nd, ok := nodes[l]
				if !ok {
					continue
				}
				refs[l]++
				if refs[l] == 2 && nd.refs == 0 {
					shared = true
				}
				if refs[l] == 1 {
					t += nd.area + ref(nd.cut)
					trace = append(trace, l)
				}
			}
			return t
		}
		total := gate + ref(cand)
		// Prefer instances where the recursion matters.
		if k >= 2 && len(trace) < 2 {
			continue
		}
		if k == 1 && len(trace) == 0 {
			continue
		}
		if k == 3 && !shared {
			continue
		}
		var desc []string
		for _, n := range order {
			nd := nodes[n]
			desc = append(desc, edakit.Sprintf("%s (best cut {%s}, gate area %d, reference count %d)", n, strings.Join(nd.cut, ", "), nd.area, nd.refs))
		}
		q := "A mapper stores, for each internal node, its best cut, the area of the gate that implements it, and its reference count in the current mapping (how many selected gates or outputs use it; 0 means it is not currently used). The internal nodes are " + strings.Join(desc, "; ") + ". Signals a to e are primary inputs. Node t is being re-evaluated with the cut {" + strings.Join(cand, ", ") + "}, implemented by a gate of area " + edakit.IntStr(gate) +
			". Its exact local area counts the gate plus, recursively, the best-cut gates of every leaf whose reference count would rise from 0, each counted once. What is the exact local area of this cut?"
		why := edakit.Sprintf("Referencing the leaves one by one makes these nodes newly used: %s. Adding their gate areas to the cut's own gate area %d gives %d.", strings.Join(trace, ", "), gate, total)
		if len(trace) == 0 {
			why = edakit.Sprintf("Every leaf is already referenced or is a primary input, so the exact area is just the gate area %d.", gate)
		}
		return edakit.Numeric(q, edakit.IntStr(total), why, map[string]string{"exact_area": edakit.IntStr(total), "new_nodes": edakit.IntStr(len(trace))})
	}
}
