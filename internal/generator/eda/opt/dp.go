package opt

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func dpPrinciple(ctx generator.GeneratorContext) generator.Problem {
	s := edakit.Scaled(ctx, 3, 4) // gates in the chain
	k := edakit.Int(2, 3)         // sizes per gate
	sizes := []string{"S", "M", "L"}[:k]
	// node cost (delay contribution) of each size at each stage, plus transition cost
	node := make([][]int, s)
	for i := range node {
		node[i] = make([]int, k)
		for j := range node[i] {
			node[i][j] = edakit.Int(1, 9)
		}
	}
	trans := make([][][]int, s-1)
	for i := range trans {
		trans[i] = make([][]int, k)
		for a := range trans[i] {
			trans[i][a] = make([]int, k)
			for b := range trans[i][a] {
				trans[i][a][b] = edakit.Int(0, 6)
			}
		}
	}
	// cost-to-go V[i][a]
	V := make([][]int, s)
	V[s-1] = append([]int(nil), node[s-1]...)
	for i := s - 2; i >= 0; i-- {
		V[i] = make([]int, k)
		for a := 0; a < k; a++ {
			best := 1 << 30
			for b := 0; b < k; b++ {
				if c := trans[i][a][b] + V[i+1][b]; c < best {
					best = c
				}
			}
			V[i][a] = node[i][a] + best
		}
	}
	opt := 1 << 30
	for a := 0; a < k; a++ {
		opt = min(opt, V[0][a])
	}
	var desc []string
	for i := 0; i < s; i++ {
		var p []string
		for a := 0; a < k; a++ {
			p = append(p, edakit.Sprintf("%s = %d", sizes[a], node[i][a]))
		}
		desc = append(desc, edakit.Sprintf("gate %d costs %s", i+1, strings.Join(p, ", ")))
	}
	var tr []string
	for i := 0; i < s-1; i++ {
		var p []string
		for a := 0; a < k; a++ {
			for b := 0; b < k; b++ {
				p = append(p, edakit.Sprintf("%s%s = %d", sizes[a], sizes[b], trans[i][a][b]))
			}
		}
		tr = append(tr, edakit.Sprintf("gates %d to %d: %s", i+1, i+2, strings.Join(p, ", ")))
	}
	stem := edakit.Sprintf("A chain of %d gates picks one of %d sizes per gate (%s). Each choice has a delay cost (%s) and each adjacent pair adds a loading cost by (driver size, load size) (%s).", s, k, strings.Join(sizes, ", "), strings.Join(desc, "; "), strings.Join(tr, "; "))
	switch kp(ctx) {
	case 1:
		a := edakit.Int(0, k-1)
		var vs []string
		for b := 0; b < k; b++ {
			vs = append(vs, edakit.Sprintf("%s = %d", sizes[b], V[s-1][b]))
		}
		return num(stem+edakit.Sprintf(" The cost-to-go of the last gate is its own cost (%s). What is the cost-to-go of gate %d at size %s, $V_{%d}(%s) = c_{%d}(%s) + \\min_b [t(%s, b) + V_{%d}(b)]$?", strings.Join(vs, ", "), s-1, sizes[a], s-1, sizes[a], s-1, sizes[a], sizes[a], s)+roundNote(0), edakit.IntStr(V[s-2][a]),
			edakit.Sprintf("Add gate %d's own cost %d to the best transition-plus-cost-to-go over the next gate's sizes. That gives %d. Bellman's principle: the best completion from here does not depend on how we got here.", s-1, node[s-2][a], V[s-2][a]),
			"v", edakit.IntStr(V[s-2][a]))
	case 2:
		return num(stem+" What is the minimum total cost over all size assignments?"+roundNote(0), edakit.IntStr(opt),
			edakit.Sprintf("Compute cost-to-go backwards from the last gate; at gate 1 the values are %s, and the minimum is %d.", ints(V[0], ", "), opt),
			"opt", edakit.IntStr(opt))
	default:
		paths := 1
		for i := 0; i < s; i++ {
			paths *= k
		}
		evals := (s - 1) * k * k
		if edakit.Int(0, 1) == 0 {
			return num(edakit.Sprintf("A chain of %d gates picks one of %d sizes per gate, with a cost for each gate's size and for each adjacent pair of sizes. How many complete size assignments would brute force evaluate?", s, k)+roundNote(0), edakit.IntStr(paths),
				edakit.Sprintf("Each gate has %d choices independently, so there are $%d^{%d} = %d$ assignments; this grows exponentially with chain length.", k, k, s, paths),
				"paths", edakit.IntStr(paths))
		}
		return num(edakit.Sprintf("A chain of %d gates picks one of %d sizes per gate, with a cost for each gate's size and for each adjacent pair of sizes. Backward dynamic programming evaluates each (size at gate i, size at gate i+1) transition once. How many transition evaluations does it make?", s, k)+roundNote(0), edakit.IntStr(evals),
			edakit.Sprintf("There are %d adjacent pairs and $%d^2$ size combinations per pair, so %d evaluations, linear in chain length, against $%d^{%d} = %d$ brute-force assignments.", s-1, k, evals, k, s, paths),
			"evals", edakit.IntStr(evals))
	}
}

type mapCosts struct{ inv, nand, and, buf int }

func drawMapCosts() mapCosts {
	return mapCosts{inv: edakit.Int(1, 2), nand: edakit.Int(2, 3), and: edakit.Int(3, 4), buf: edakit.Int(2, 3)}
}

func (m mapCosts) String() string {
	return edakit.Sprintf("INV = %d, NAND2 = %d, AND2 (an INV fed by a NAND2) = %d, BUF (an INV fed by an INV) = %d", m.inv, m.nand, m.and, m.buf)
}

// chainDP maps a chain of node types (from input to output) and returns best[i] for every node.
func chainDP(chain []string, m mapCosts) []int {
	best := make([]int, len(chain))
	get := func(i int) int {
		if i < 0 {
			return 0
		}
		return best[i]
	}
	for i, t := range chain {
		if t == "NAND" {
			best[i] = get(i-1) + m.nand
			continue
		}
		best[i] = get(i-1) + m.inv
		if i >= 1 {
			if chain[i-1] == "NAND" {
				best[i] = min(best[i], get(i-2)+m.and)
			} else {
				best[i] = min(best[i], get(i-2)+m.buf)
			}
		}
	}
	return best
}

func randomChain(n int) []string {
	c := make([]string, n)
	for i := range c {
		c[i] = edakit.Pick([]string{"NAND", "INV", "INV"})
	}
	return c
}

func chainStr(c []string) string {
	out := make([]string, len(c))
	for i, t := range c {
		if t == "NAND" {
			out[i] = edakit.Sprintf("n%d NAND2", i+1)
		} else {
			out[i] = edakit.Sprintf("n%d INV", i+1)
		}
	}
	return strings.Join(out, ", ")
}

func treeMapping(ctx generator.GeneratorContext) generator.Problem {
	m := drawMapCosts()
	switch kp(ctx) {
	case 1:
		n := edakit.Int(4, 7)
		c := randomChain(n)
		c[n-1] = "INV"
		best := chainDP(c, m)
		var known []string
		for i := 0; i < n-1; i++ {
			known = append(known, edakit.Sprintf("best(n%d) = %d", i+1, best[i]))
		}
		return num("A tree mapper covers a NAND2/INV subject chain, listed from the input side: "+chainStr(c)+". Each NAND2's other input is a primary input. Library areas: "+m.String()+". The DP has "+strings.Join(known, ", ")+edakit.Sprintf(". What is best(n%d), the cheapest cover of the chain up to and including n%d?", n, n)+roundNote(0), edakit.IntStr(best[n-1]),
			edakit.Sprintf("At n%d either a 1-node INV covers it (best(n%d) + %d) or a 2-node pattern rooted at n%d covers it and n%d (best(n%d) + its area). The cheaper option gives %d.", n, n-1, m.inv, n, n-1, n-2, best[n-1]),
			"best", edakit.IntStr(best[n-1]))
	case 2:
		n := edakit.Scaled(ctx, 4, 8)
		c := randomChain(n)
		best := chainDP(c, m)
		return num("A tree mapper covers a NAND2/INV subject chain, listed from the input side: "+chainStr(c)+". Each NAND2's other input is a primary input. Library areas: "+m.String()+". What is the minimum total area of a cover?"+roundNote(0), edakit.IntStr(best[n-1]),
			edakit.Sprintf("Process nodes from the input: best(i) is the minimum over patterns rooted at node i of the pattern's area plus best() of the node below it. The values are %s, so the minimum area is %d.", ints(best, ", "), best[n-1]),
			"area", edakit.IntStr(best[n-1]))
	default:
		la, lb := edakit.Int(1, 3), edakit.Int(1, 3)
		A, B := randomChain(la), randomChain(lb)
		ba, bb := chainDP(A, m), chainDP(B, m)
		nandBest := m.nand + ba[la-1] + bb[lb-1]
		root := min(m.inv+nandBest, m.and+ba[la-1]+bb[lb-1])
		return num("A subject tree has an output INV fed by a NAND2 whose two inputs are chains. Chain a (from its primary input up to the NAND2): "+chainStr(A)+". Chain b: "+chainStr(B)+". Library areas: "+m.String()+". What is the minimum area to cover the whole tree?"+roundNote(0), edakit.IntStr(root),
			edakit.Sprintf("Solve each chain first: %d and %d. The NAND2 alone costs %d + both chains = %d. At the root, INV + that is %d, while AND2 covering the INV and the NAND2 is %d + both chains = %d. The minimum is %d; subtrees are solved independently because a tree has no shared nodes.", ba[la-1], bb[lb-1], m.nand, nandBest, m.inv+nandBest, m.and, m.and+ba[la-1]+bb[lb-1], root),
			"area", edakit.IntStr(root))
	}
}

func bufferInsertion(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		C := float64(edakit.Int(5, 60))
		q := float64(edakit.Int(100, 600))
		Rw := float64(edakit.Int(20, 200))
		Cw := float64(edakit.Int(5, 40))
		nq := q - Rw*(Cw/2+C)/1000
		return num(edakit.Sprintf("Van Ginneken's DP carries candidates (load C, required time q) up the tree. A candidate (C = %s fF, q = %s ps) passes up through a wire with R = %s ohm and C = %s fF (pi model). What is the candidate's new required time, in ps?", f(C, 0), f(q, 0), f(Rw, 0), f(Cw, 0))+roundNote(3), f(nq, 3),
			edakit.Sprintf("The wire's Elmore delay is $R(C_w/2 + C) = %s(%s + %s)$ ohm fF = %s ps, so $q' = %s - %s = %s$ ps and the load grows to %s fF.", f(Rw, 0), f(Cw/2, 1), f(C, 0), f(Rw*(Cw/2+C)/1000, 3), f(q, 0), f(Rw*(Cw/2+C)/1000, 3), f(nq, 3), f(C+Cw, 0)),
			"q", f(nq, 3), "c", f(C+Cw, 0))
	case 2:
		n := edakit.Int(2, 4)
		Rb := float64(edakit.Int(100, 800))
		Db := float64(edakit.Int(10, 40))
		Cb := float64(edakit.Int(2, 8))
		best := -1e18
		var parts, vals []string
		for i := 0; i < n; i++ {
			C := float64(edakit.Int(10, 120))
			q := float64(edakit.Int(200, 500))
			v := q - Db - Rb*C/1000
			if v > best {
				best = v
			}
			parts = append(parts, edakit.Sprintf("(%s fF, %s ps)", f(C, 0), f(q, 0)))
			vals = append(vals, f(v, 3))
		}
		return num("At a buffer site the DP holds candidates "+strings.Join(parts, ", ")+edakit.Sprintf(". A buffer has intrinsic delay %s ps, output resistance %s ohm and input capacitance %s fF. Buffering a candidate gives (%s fF, $q - D_b - R_b C$). What is the best required time among the buffered candidates, in ps?", f(Db, 0), f(Rb, 0), f(Cb, 0), f(Cb, 0))+roundNote(3), f(best, 3),
			"Buffering each candidate gives required times "+strings.Join(vals, ", ")+edakit.Sprintf(" ps, all with load %s fF. Only the best survives, since they share the same load: %s ps.", f(Cb, 0), f(best, 3)),
			"best", f(best, 3))
	default:
		if edakit.Int(0, 1) == 0 {
			n := edakit.Scaled(ctx, 5, 9)
			Cs := distinctInts(n, 5, 60)
			Qs := distinctInts(n, 100, 400)
			type cand struct{ c, q int }
			cs := make([]cand, n)
			var parts []string
			for i := range cs {
				cs[i] = cand{Cs[i], Qs[i]}
				parts = append(parts, edakit.Sprintf("(%d, %d)", Cs[i], Qs[i]))
			}
			alive := 0
			var keep []string
			for i := range cs {
				dom := false
				for j := range cs {
					if j != i && cs[j].c <= cs[i].c && cs[j].q >= cs[i].q {
						dom = true
					}
				}
				if !dom {
					alive++
					keep = append(keep, edakit.Sprintf("(%d, %d)", cs[i].c, cs[i].q))
				}
			}
			return num("A node's candidate list holds (load in fF, required time in ps) pairs "+strings.Join(parts, ", ")+". A candidate is pruned if another has no more load and no less required time. How many candidates survive?"+roundNote(0), edakit.IntStr(alive),
				edakit.Sprintf("Sorting by load, a survivor must have strictly larger q than every lighter candidate. The survivors are %s: %d. Pruning keeps the list linear in the number of buffer positions instead of exponential.", strings.Join(keep, ", "), alive),
				"survivors", edakit.IntStr(alive))
		}
		n := edakit.Int(2, 4)
		Rd := float64(edakit.Int(100, 1000))
		best := -1e18
		var parts, vals []string
		type cq struct{ c, q float64 }
		var cs []cq
		for i := 0; i < n; i++ {
			cs = append(cs, cq{float64(edakit.Int(5, 150)), float64(edakit.Int(150, 500))})
		}
		sort.Slice(cs, func(i, j int) bool { return cs[i].c < cs[j].c })
		for _, x := range cs {
			v := x.q - Rd*x.c/1000
			best = max(best, v)
			parts = append(parts, edakit.Sprintf("(%s fF, %s ps)", f(x.c, 0), f(x.q, 0)))
			vals = append(vals, f(v, 3))
		}
		return num("At the source the DP holds candidates "+strings.Join(parts, ", ")+edakit.Sprintf(". The driver has output resistance %s ohm. Slack at the driver input is $q - R_d C$. What is the best achievable slack, in ps?", f(Rd, 0))+roundNote(3), f(best, 3),
			"Subtracting the driver delay $R_d C$ from each gives "+strings.Join(vals, ", ")+edakit.Sprintf(" ps; the DP returns the best, %s ps, and traces back the buffer positions that built it.", f(best, 3)),
			"slack", f(best, 3))
	}
}

func constraintPropagation(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 3:
		k := ctx.KPOrdinal
		if k == 0 {
			k = edakit.Pick([]int{1, 3})
		}
		n := edakit.Scaled(ctx, 3, 5)
		H := edakit.Int(8, 16)
		d := make([]int, n-1)
		for i := range d {
			d[i] = edakit.Int(1, 4)
		}
		lo, hi := make([]int, n), make([]int, n)
		init := make([]int, n)
		for i := range lo {
			lo[i], hi[i] = 0, H
			if edakit.Int(0, 2) == 0 {
				lo[i] = edakit.Int(0, 3)
			}
			if edakit.Int(0, 2) == 0 {
				hi[i] = H - edakit.Int(0, 3)
			}
			init[i] = hi[i] - lo[i] + 1
		}
		var dom, cons []string
		for i := 0; i < n; i++ {
			dom = append(dom, edakit.Sprintf("$s_%d \\in \\{%d..%d\\}$", i+1, lo[i], hi[i]))
		}
		for i := 0; i < n-1; i++ {
			cons = append(cons, edakit.Sprintf("$s_%d + %d \\le s_%d$", i+1, d[i], i+2))
		}
		// propagate to fixpoint
		for changed := true; changed; {
			changed = false
			for i := 0; i < n-1; i++ {
				if lo[i]+d[i] > lo[i+1] {
					lo[i+1] = lo[i] + d[i]
					changed = true
				}
				if hi[i+1]-d[i] < hi[i] {
					hi[i] = hi[i+1] - d[i]
					changed = true
				}
			}
		}
		feasible := true
		for i := range lo {
			if lo[i] > hi[i] {
				feasible = false
			}
		}
		if !feasible {
			return constraintPropagation(ctx)
		}
		stem := "A scheduling CSP has start-time domains " + strings.Join(dom, ", ") + " and precedence constraints " + strings.Join(cons, ", ") + "."
		var fin []string
		pruned := 0
		for i := range lo {
			fin = append(fin, edakit.Sprintf("$s_%d \\in \\{%d..%d\\}$", i+1, lo[i], hi[i]))
			pruned += init[i] - (hi[i] - lo[i] + 1)
		}
		expl := "Bounds propagation pushes lower bounds forward along the chain and upper bounds backward until nothing changes: " + strings.Join(fin, ", ") + ". For these binary precedence constraints that is exactly arc consistency."
		if k == 1 {
			v := edakit.Int(1, n)
			sz := hi[v-1] - lo[v-1] + 1
			return num(stem+edakit.Sprintf(" After arc consistency, how many values remain in the domain of $s_%d$?", v)+roundNote(0), edakit.IntStr(sz), expl, "size", edakit.IntStr(sz))
		}
		return num(stem+" In total, how many values does arc consistency remove from all the domains together?"+roundNote(0), edakit.IntStr(pruned), expl+edakit.Sprintf(" The domains shrink by %d values in all, with no search.", pruned), "pruned", edakit.IntStr(pruned))
	default:
		for {
			n := edakit.Scaled(ctx, 3, 4)
			T := edakit.Int(3, 5)
			names := []string{"a", "b", "c", "d"}[:n]
			type pair struct{ i, j int }
			var prec, conf []pair
			for i := 0; i < n; i++ {
				for j := i + 1; j < n; j++ {
					switch edakit.Int(0, 2) {
					case 0:
						prec = append(prec, pair{i, j})
					case 1:
						conf = append(conf, pair{i, j})
					}
				}
			}
			if len(prec)+len(conf) < 2 {
				continue
			}
			cnt := 0
			x := make([]int, n)
			var rec func(k int)
			rec = func(k int) {
				if k == n {
					for _, p := range prec {
						if x[p.i] >= x[p.j] {
							return
						}
					}
					for _, p := range conf {
						if x[p.i] == x[p.j] {
							return
						}
					}
					cnt++
					return
				}
				for t := 1; t <= T; t++ {
					x[k] = t
					rec(k + 1)
				}
			}
			rec(0)
			if cnt == 0 {
				continue
			}
			var cs []string
			for _, p := range prec {
				cs = append(cs, edakit.Sprintf("%s before %s", names[p.i], names[p.j]))
			}
			for _, p := range conf {
				cs = append(cs, edakit.Sprintf("%s and %s share a unit (different steps)", names[p.i], names[p.j]))
			}
			total := 1
			for i := 0; i < n; i++ {
				total *= T
			}
			return num(edakit.Sprintf("Operations %s each take one of control steps 1 to %d, subject to: %s. How many solutions does this CSP have?", strings.Join(names, ", "), T, strings.Join(cs, "; "))+roundNote(0), edakit.IntStr(cnt),
				edakit.Sprintf("Out of $%d^%d = %d$ assignments, exactly %d satisfy every precedence ($<$) and resource ($\\ne$) constraint. A CP solver finds them without enumerating: propagation prunes domains and search branches only on what is left.", T, n, total, cnt),
				"solutions", edakit.IntStr(cnt))
		}
	}
}
