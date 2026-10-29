package synth

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// internalNode is y = gate(u, v) over primary inputs.
type internalNode struct {
	g    gate2
	u, v int
}

func (y internalNode) eval(x []bool) bool { return y.g.f(x[y.u], x[y.v]) }

func randInternal(nIn int) internalNode {
	u := edakit.Int(0, nIn-1)
	v := edakit.Int(0, nIn-2)
	if v >= u {
		v++
	}
	return internalNode{edakit.Pick(gates2), u, v}
}

func inputs(m, n int) []bool {
	x := make([]bool, n)
	for i := 0; i < n; i++ {
		x[i] = m>>(n-1-i)&1 == 1
	}
	return x
}

// reachable returns the set of local patterns (y1..yk as a k-bit word, y1 most significant).
func reachable(ys []internalNode, nIn int) uint64 {
	var seen uint64
	k := len(ys)
	for m := 0; m < 1<<nIn; m++ {
		x := inputs(m, nIn)
		p := 0
		for i, y := range ys {
			if y.eval(x) {
				p |= 1 << (k - 1 - i)
			}
		}
		seen |= 1 << p
	}
	return seen
}

func cubeTexNames(c cube, n int, names []string) string {
	if c.mask == 0 {
		return "1"
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		p := uint32(1) << (n - 1 - i)
		if c.mask&p == 0 {
			continue
		}
		if c.val&p != 0 {
			b.WriteString(names[i])
		} else {
			b.WriteString("\\bar{" + names[i] + "}")
		}
	}
	return b.String()
}

// sdc: controllability (satisfiability-derived) don't cares at a node's local inputs.
func sdc(ctx generator.GeneratorContext) generator.Problem {
	nIn := 3
	k := 2
	which := kp(ctx)
	if which >= 2 {
		k = 3
	}
	if ctx.Difficulty > 0.7 {
		nIn = 4
	}
	for {
		ys := make([]internalNode, k)
		var desc []string
		for i := range ys {
			ys[i] = randInternal(nIn)
			desc = append(desc, edakit.Sprintf("$y_%d = %s$", i+1, ys[i].g.tex(varNames[ys[i].u], varNames[ys[i].v])))
		}
		reach := reachable(ys, nIn)
		cdc := (1 << k) - popcount(reach)
		if cdc == 0 && edakit.Int(0, 2) > 0 {
			continue
		}
		local := "$(y_1, y_2)$"
		if k == 3 {
			local = "$(y_1, y_2, y_3)$"
		}
		stem := "The inputs " + strings.Join(varNames[:nIn], ", ") + " are free. Internal nodes " + joinTex(desc) + " feed a node g whose local inputs are " + local + "."
		if which <= 2 {
			return edakit.Numeric(stem+edakit.Sprintf(" The satisfiability don't cares $y_i \\oplus f_i$ make some of the %d local patterns impossible. How many local patterns of %s can never occur?", 1<<k, local),
				edakit.IntStr(cdc), edakit.Sprintf("Simulate all %d input vectors and collect the patterns that appear: %d appear, so %d are controllability don't cares for g.", 1<<nIn, popcount(reach), cdc),
				map[string]string{"cdc": edakit.IntStr(cdc)})
		}
		// Minimize g over y1..y3 with the unreachable patterns as don't cares.
		gon := randFunc(3, 2, 5) & reach
		if gon == 0 || gon == reach || cdc == 0 {
			continue
		}
		dc := ^reach & 0xFF
		names := []string{"y_1", "y_2", "y_3"}
		cov := minCover(3, gon, primes(3, gon, dc))
		parts := make([]string, len(cov))
		for i, c := range cov {
			parts[i] = cubeTexNames(c, 3, names)
		}
		l := coverLits(cov)
		return edakit.Numeric(stem+" Node g must be 1 exactly on the reachable local patterns (writing the pattern as the binary number $y_1y_2y_3$) "+intsCSV(mintermList(gon, 3))+", and 0 on the other reachable ones. Using the unreachable patterns as don't cares, how many literals does a minimum SOP for g have?",
			edakit.IntStr(l), edakit.Sprintf("Unreachable patterns %s are don't cares. A minimum cover is $g = %s$ with %d literals.", intsCSV(mintermList(dc, 3)), strings.Join(parts, " + "), l),
			map[string]string{"literals": edakit.IntStr(l)})
	}
}

// zfun is an output built from two named two-input gates over y, c and d.
type zfun struct {
	tex string
	f   func(y, c, d bool) bool
}

func gateTex(g gate2, a, b string) string {
	wrap := func(x string) string {
		if strings.ContainsAny(x, "+\\") && !strings.HasPrefix(x, "\\overline") {
			return "(" + x + ")"
		}
		return x
	}
	switch g.name {
	case "AND":
		return g.tex(wrap(a), wrap(b))
	case "ANDN":
		return g.tex(wrap(a), b)
	}
	return g.tex(a, b)
}

func randZ() zfun {
	g1, g2 := edakit.Pick(gates2), edakit.Pick(gates2)
	switch edakit.Int(0, 2) {
	case 0:
		in := gateTex(g1, "y", "c")
		return zfun{gateTex(g2, in, "d"), func(y, c, d bool) bool { return g2.f(g1.f(y, c), d) }}
	case 1:
		in := gateTex(g1, "y", "d")
		return zfun{gateTex(g2, in, "c"), func(y, c, d bool) bool { return g2.f(g1.f(y, d), c) }}
	}
	in := gateTex(g1, "c", "d")
	return zfun{gateTex(g2, "y", in), func(y, c, d bool) bool { return g2.f(y, g1.f(c, d)) }}
}

// odc: observability don't cares of a node y through one or two outputs.
func odc(ctx generator.GeneratorContext) generator.Problem {
	which := kp(ctx)
	for {
		yg := edakit.Pick(gates2)
		yDesc := "$y = " + yg.tex("a", "b") + "$"
		nOut := 1
		if which == 3 {
			nOut = 2
		}
		var zs []zfun
		for i := 0; i < nOut; i++ {
			zs = append(zs, randZ())
		}
		count := 0
		for m := 0; m < 4; m++ {
			c, d := m>>1&1 == 1, m&1 == 1
			all := true
			for _, z := range zs {
				if z.f(false, c, d) != z.f(true, c, d) {
					all = false
				}
			}
			if all {
				count++
			}
		}
		if count == 4 || (count == 0 && edakit.Int(0, 1) == 0) {
			continue
		}
		total := count * 4 // a and b are free
		if which == 3 {
			return edakit.Numeric(edakit.Sprintf("Node %s fans out to two primary outputs, $z_1 = %s$ and $z_2 = %s$. Over the 16 combinations of a, b, c, d, how many are observability don't cares of y (y is observable at neither output)?",
				yDesc, zs[0].tex, zs[1].tex),
				edakit.IntStr(total), edakit.Sprintf("y is unobservable where both Boolean differences $\\partial z_1/\\partial y$ and $\\partial z_2/\\partial y$ are 0. That holds for %d of the 4 (c, d) values, and a, b are free, so %d minterms.", count, total),
				map[string]string{"odc": edakit.IntStr(total)})
		}
		if which == 2 {
			return edakit.Numeric(edakit.Sprintf("Node %s drives only the output $z = %s$. How many of the 16 combinations of a, b, c, d are in the observability care set of y (where flipping y flips z)?", yDesc, zs[0].tex),
				edakit.IntStr(16-total), edakit.Sprintf("The care set is where $\\partial z/\\partial y = z|_{y=1} \\oplus z|_{y=0} = 1$: %d of the 4 (c, d) values, times 4 for a and b, is %d.", 4-count, 16-total),
				map[string]string{"care": edakit.IntStr(16 - total)})
		}
		return edakit.Numeric(edakit.Sprintf("Node %s drives only the output $z = %s$. How many of the 16 combinations of a, b, c, d are observability don't cares of y (changing y would not change z)?", yDesc, zs[0].tex),
			edakit.IntStr(total), edakit.Sprintf("The ODC is the complement of the Boolean difference $z|_{y=1} \\oplus z|_{y=0}$. It is 1 for %d of the 4 (c, d) values; a and b are free, so %d minterms.", count, total),
			map[string]string{"odc": edakit.IntStr(total)})
	}
}

// resubstitution: functional dependency of f on a set of divisors.
func resubstitution(ctx generator.GeneratorContext) generator.Problem {
	nIn := 4
	which := kp(ctx)
	for {
		gs := make([]internalNode, 4)
		var desc []string
		for i := range gs {
			gs[i] = randInternal(nIn)
			desc = append(desc, edakit.Sprintf("$g_%d = %s$", i+1, gs[i].g.tex(varNames[gs[i].u], varNames[gs[i].v])))
		}
		distinct := map[uint64]bool{}
		for _, g := range gs {
			var t uint64
			for m := 0; m < 16; m++ {
				if g.eval(inputs(m, nIn)) {
					t |= 1 << m
				}
			}
			distinct[t] = true
		}
		if len(distinct) < 4 {
			continue
		}
		i, j := edakit.Int(0, 3), edakit.Int(0, 2)
		if j >= i {
			j++
		}
		h := edakit.Pick(gates2[:3])
		var on uint64
		for m := 0; m < 16; m++ {
			x := inputs(m, nIn)
			if h.f(gs[i].eval(x), gs[j].eval(x)) {
				on |= 1 << m
			}
		}
		if on == 0 || on == 0xFFFF {
			continue
		}
		var dc uint64
		if which == 3 {
			dc = randFunc(4, 3, 5)
		}
		depends := func(s int) bool {
			seen := map[int]bool{}
			val := map[int]bool{}
			for m := 0; m < 16; m++ {
				if dc>>m&1 == 1 {
					continue
				}
				x := inputs(m, nIn)
				key := 0
				for t := 0; t < 4; t++ {
					if s>>t&1 == 1 && gs[t].eval(x) {
						key |= 1 << t
					}
				}
				fv := on>>m&1 == 1
				if seen[key] && val[key] != fv {
					return false
				}
				seen[key], val[key] = true, fv
			}
			return true
		}
		stem := "Target $f(a,b,c,d) = \\sum m(" + intsCSV(mintermList(on, 4)) + ")$ and candidate divisors " + joinTex(desc) + "."
		if which == 1 {
			pairs := map[int]bool{}
			for m := 0; m < 16; m++ {
				x := inputs(m, nIn)
				k := 0
				if gs[0].eval(x) {
					k |= 2
				}
				if gs[1].eval(x) {
					k |= 1
				}
				pairs[k] = true
			}
			return edakit.Numeric(stem+" Resubstitution asks which divisor values ever occur. How many distinct values of the pair $(g_1, g_2)$ occur over the 16 input vectors?",
				edakit.IntStr(len(pairs)), edakit.Sprintf("Simulate both divisors on all 16 vectors: %d of the 4 pairs appear. The missing ones are don't cares for any function of $g_1, g_2$.", len(pairs)),
				map[string]string{"pairs": edakit.IntStr(len(pairs))})
		}
		best := -1
		for size := 0; size <= 4 && best < 0; size++ {
			for s := 0; s < 16; s++ {
				if popcount(uint64(s)) == size && depends(s) {
					best = size
					break
				}
			}
		}
		if which == 2 {
			return edakit.Numeric(stem+" f can be re-expressed as a function of a set S of divisors exactly when no two input vectors agree on every divisor in S but disagree on f. What is the smallest size of such a set S?",
				edakit.IntStr(best), edakit.Sprintf("Checking subsets by size, the first set that separates every pair of vectors with different f values has %d divisors.", best),
				map[string]string{"min_divisors": edakit.IntStr(best)})
		}
		return edakit.Numeric(stem+" The input vectors "+intsCSV(mintermList(dc, 4))+" are don't cares for f (never applied or never observed). Ignoring them, what is the smallest number of divisors from which f can be computed?",
			edakit.IntStr(best), edakit.Sprintf("Only care vectors must be separated. The smallest separating set has %d divisors.", best),
			map[string]string{"min_divisors": edakit.IntStr(best)})
	}
}
