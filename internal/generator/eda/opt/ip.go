package opt

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type knap struct {
	v, w []int
	cap  int
}

func drawKnap(ctx generator.GeneratorContext) knap {
	for {
		n := edakit.Scaled(ctx, 3, 6)
		k := knap{}
		tot := 0
		for i := 0; i < n; i++ {
			k.v = append(k.v, edakit.Int(4, 30))
			k.w = append(k.w, edakit.Int(2, 12))
			tot += k.w[i]
		}
		k.cap = tot/2 + edakit.Int(-2, 2)
		if k.cap < 3 {
			continue
		}
		// distinct ratios keep the greedy order unambiguous
		seen := map[float64]bool{}
		ok := true
		for i := range k.v {
			r := float64(k.v[i]) / float64(k.w[i])
			if seen[r] {
				ok = false
			}
			seen[r] = true
		}
		if ok {
			return k
		}
	}
}

func (k knap) stem() string {
	var parts []string
	for i := range k.v {
		parts = append(parts, edakit.Sprintf("upsize %d recovers %d ps of TNS for %d area units", i+1, k.v[i], k.w[i]))
	}
	return "A timing-repair pass chooses a subset of candidate upsizes under an area budget of " + edakit.IntStr(k.cap) + " units: " + strings.Join(parts, "; ") + ". Each upsize is taken whole ($x_i \\in \\{0,1\\}$)."
}

// lpBound is the fractional knapsack bound with some items forced: fix[i] = -1 free, 0 out, 1 in.
// It returns the bound, the fractional item (or -1) and the greedy order.
func (k knap) lpBound(fix []int) (float64, int) {
	capLeft := float64(k.cap)
	val := 0.0
	var free []int
	for i := range k.v {
		switch fix[i] {
		case 1:
			capLeft -= float64(k.w[i])
			val += float64(k.v[i])
		case -1:
			free = append(free, i)
		}
	}
	if capLeft < 0 {
		return math.Inf(-1), -1
	}
	sort.Slice(free, func(a, b int) bool {
		return float64(k.v[free[a]])/float64(k.w[free[a]]) > float64(k.v[free[b]])/float64(k.w[free[b]])
	})
	for _, i := range free {
		if float64(k.w[i]) <= capLeft {
			capLeft -= float64(k.w[i])
			val += float64(k.v[i])
		} else {
			val += float64(k.v[i]) * capLeft / float64(k.w[i])
			if capLeft > 0 {
				return val, i
			}
			return val, -1
		}
	}
	return val, -1
}

func (k knap) intOpt() int {
	n := len(k.v)
	best := 0
	for m := 0; m < 1<<n; m++ {
		w, v := 0, 0
		for i := 0; i < n; i++ {
			if m&(1<<i) != 0 {
				w += k.w[i]
				v += k.v[i]
			}
		}
		if w <= k.cap && v > best {
			best = v
		}
	}
	return best
}

func freeFix(n int) []int {
	fix := make([]int, n)
	for i := range fix {
		fix[i] = -1
	}
	return fix
}

func lpRelaxation(ctx generator.GeneratorContext) generator.Problem {
	k := drawKnap(ctx)
	lp, _ := k.lpBound(freeFix(len(k.v)))
	ip := k.intOpt()
	switch kp(ctx) {
	case 1:
		return num(k.stem()+" Relax each $x_i$ to $0 \\le x_i \\le 1$. What is the LP relaxation's optimal value, an upper bound on the recoverable TNS, in ps?"+roundNote(2),
			f(lp, 2),
			edakit.Sprintf("The relaxed knapsack is solved greedily: take upsizes in decreasing order of ps per area unit and take a fraction of the first one that does not fit. That gives %s ps.", f(lp, 2)),
			"lp_bound", f(lp, 2))
	case 2:
		return num(k.stem()+" What is the best TNS recovery, in ps, achievable with whole upsizes?"+roundNote(0),
			edakit.IntStr(ip),
			edakit.Sprintf("Enumerating the feasible subsets gives an integer optimum of %d ps. The LP bound is %s ps, and the integer answer can never exceed it.", ip, f(lp, 2)),
			"ip_opt", edakit.IntStr(ip), "lp_bound", f(lp, 2))
	default:
		return num(k.stem()+" What is the integrality gap, the LP relaxation value minus the integer optimum, in ps?"+roundNote(2),
			f(lp-float64(ip), 2),
			edakit.Sprintf("The fractional greedy bound is %s ps and the best whole-upsize subset recovers %d ps, so the gap is %s ps. A solver must close that gap by branching or cutting.", f(lp, 2), ip, f(lp-float64(ip), 2)),
			"lp_bound", f(lp, 2), "ip_opt", edakit.IntStr(ip))
	}
}

func branchAndBound(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 8)
		maxSense := edakit.Int(0, 1) == 0
		inc := edakit.Int(40, 60)
		bounds := make([]int, n)
		count := 0
		for i := range bounds {
			bounds[i] = inc + edakit.Int(-8, 8)
			if maxSense && bounds[i] <= inc || !maxSense && bounds[i] >= inc {
				count++
			}
		}
		var q, why string
		if maxSense {
			q = edakit.Sprintf("A branch-and-bound search maximizing recovered slack has an incumbent integer solution worth %d. The open nodes have LP upper bounds %s. A node is pruned when its bound cannot strictly beat the incumbent. How many open nodes can be pruned?", inc, ints(bounds, ", "))
			why = edakit.Sprintf("For maximization a node whose upper bound is at most the incumbent %d cannot contain a better solution. %d of the bounds satisfy that.", inc, count)
		} else {
			q = edakit.Sprintf("A branch-and-bound search for a minimum test-set cover (covering in ATPG compaction) has an incumbent cover of %d patterns. The open nodes have LP lower bounds %s. A node is pruned when its bound cannot strictly beat the incumbent. How many open nodes can be pruned?", inc, ints(bounds, ", "))
			why = edakit.Sprintf("For minimization a node whose lower bound is at least the incumbent %d cannot contain a smaller cover. %d of the bounds satisfy that.", inc, count)
		}
		return num(q+roundNote(0), edakit.IntStr(count), why, "pruned", edakit.IntStr(count))
	case 2:
		for {
			k := drawKnap(ctx)
			root, fi := k.lpBound(freeFix(len(k.v)))
			if fi < 0 {
				continue
			}
			side := edakit.Int(0, 1)
			fix := freeFix(len(k.v))
			fix[fi] = side
			child, _ := k.lpBound(fix)
			if math.IsInf(child, -1) {
				continue
			}
			return num(k.stem()+edakit.Sprintf(" The root LP relaxation is worth %s ps and leaves upsize %d fractional. The search branches on it. What is the LP bound of the child with $x_%d = %d$, in ps?", f(root, 2), fi+1, fi+1, side)+roundNote(2),
				f(child, 2),
				edakit.Sprintf("Fix $x_%d = %d$, then solve the remaining fractional knapsack greedily by ps per area unit with the budget that is left. The child bound is %s ps, never above the root's %s ps.", fi+1, side, f(child, 2), f(root, 2)),
				"root", f(root, 2), "child", f(child, 2))
		}
	default:
		inc := edakit.Int(200, 900)
		bound := inc + edakit.Int(1, 120)
		gap := 100 * float64(bound-inc) / float64(inc)
		return num(edakit.Sprintf("A MILP global router reports an incumbent of %d (maximized routed-net score) and a best bound of %d over all open nodes. Its relative gap is (bound minus incumbent) divided by the incumbent. What is the gap, in percent?", inc, bound)+roundNote(2),
			f(gap, 2),
			edakit.Sprintf("Gap = (%d - %d) / %d = %s%%. No solution in the tree can be better than the bound, so the incumbent is within that percentage of optimal; solvers stop when the gap falls under a tolerance.", bound, inc, inc, f(gap, 2)),
			"gap", f(gap, 2))
	}
}

func bigM(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		wi := edakit.Int(5, 40)
		ui := edakit.Int(50, 200)
		lj := edakit.Int(0, 40)
		m := ui + wi - lj
		return num(edakit.Sprintf("An ILP floorplanner keeps block i left of block j unless binary $p = 1$: $x_i + w_i \\le x_j + Mp$. Block i has width %d um and its fence allows $0 \\le x_i \\le %d$; block j's fence requires $x_j \\ge %d$. What is the smallest M, in um, that makes $p = 1$ switch the constraint off for every feasible placement?", wi, ui, lj)+roundNote(0),
			edakit.IntStr(m),
			edakit.Sprintf("With p = 1 the constraint must hold for the worst case, largest $x_i$ and smallest $x_j$: $M \\ge %d + %d - %d = %d$. A larger M is still valid but weakens the LP relaxation.", ui, wi, lj, m),
			"m", edakit.IntStr(m))
	case 2:
		n := edakit.Scaled(ctx, 3, 12)
		rot := edakit.Int(0, 1) == 1
		vars := n * (n - 1)
		extra := ""
		if rot {
			vars += n
			extra = " and one rotation binary per block"
		}
		return num(edakit.Sprintf("A sequence of big-M disjunctions keeps %d blocks from overlapping: each pair uses two binaries (p, q) to pick one of left-of, right-of, below or above%s. How many binary variables does the model need?", n, extra)+roundNote(0),
			edakit.IntStr(vars),
			edakit.Sprintf("There are $\\binom{%d}{2} = %d$ pairs and 2 binaries each, %d in all%s. The quadratic growth in pairs, together with weak big-M relaxations, is why exact ILP floorplanning is practical only for small block counts.", n, n*(n-1)/2, n*(n-1), map[bool]string{true: edakit.Sprintf(", plus %d rotation binaries gives %d", n, vars), false: ""}[rot]),
			"binaries", edakit.IntStr(vars))
	default:
		wi := edakit.Int(10, 40)
		xi := edakit.Int(20, 80)
		xj := xi + edakit.Int(1, wi-1)
		m := edakit.Int(150, 400)
		p := float64(xi+wi-xj) / float64(m)
		return num(edakit.Sprintf("In the LP relaxation of $x_i + w_i \\le x_j + Mp$ with $w_i = %d$ um and $M = %d$ um, the solver places $x_i = %d$ and $x_j = %d$, so the blocks overlap. What is the smallest fractional p that satisfies the relaxed constraint?", wi, m, xi, xj)+roundNote(3),
			f(p, 3),
			edakit.Sprintf("Need $%d + %d \\le %d + %d p$, so $p \\ge %d/%d = %s$. A tiny fractional p buys a real overlap at no cost, which is why big-M relaxations give weak bounds.", xi, wi, xj, m, xi+wi-xj, m, f(p, 3)),
			"p", f(p, 3))
	}
}

func ilpScheduling(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 8)
		L := edakit.Int(5, 9)
		total := 0
		var parts []string
		for i := 0; i < n; i++ {
			a := edakit.Int(1, L-1)
			b := edakit.Int(a, L)
			total += b - a + 1
			parts = append(parts, edakit.Sprintf("op%d [%d, %d]", i+1, a, b))
		}
		return num(edakit.Sprintf("A time-indexed ILP schedules %d operations into %d control steps with binaries $x_{i,t} = 1$ when operation i starts in step t. Variables exist only inside each operation's [ASAP, ALAP] window: ", n, L)+strings.Join(parts, ", ")+". How many binary variables does the model have?"+roundNote(0),
			edakit.IntStr(total),
			edakit.Sprintf("Each operation contributes ALAP - ASAP + 1 variables. Summing the windows gives %d, against %d without mobility pruning.", total, n*L),
			"binaries", edakit.IntStr(total))
	case 2:
		n := edakit.Scaled(ctx, 5, 8)
		preds := make([][]int, n)
		for j := 1; j < n; j++ {
			for i := 0; i < j; i++ {
				if edakit.Int(0, 99) < 35 {
					preds[j] = append(preds[j], i)
				}
			}
		}
		asap := make([]int, n)
		for j := 0; j < n; j++ {
			asap[j] = 1
			for _, i := range preds[j] {
				if asap[i]+1 > asap[j] {
					asap[j] = asap[i] + 1
				}
			}
		}
		crit := 0
		for _, a := range asap {
			if a > crit {
				crit = a
			}
		}
		L := crit + edakit.Int(0, 2)
		alap := make([]int, n)
		for i := range alap {
			alap[i] = L
		}
		for j := n - 1; j >= 0; j-- {
			for _, i := range preds[j] {
				if alap[j]-1 < alap[i] {
					alap[i] = alap[j] - 1
				}
			}
		}
		k := edakit.Int(0, n-1)
		var edges []string
		for j := 0; j < n; j++ {
			for _, i := range preds[j] {
				edges = append(edges, edakit.Sprintf("op%d to op%d", i+1, j+1))
			}
		}
		es := "none"
		if len(edges) > 0 {
			es = strings.Join(edges, ", ")
		}
		mob := alap[k] - asap[k]
		return num(edakit.Sprintf("A dataflow graph has %d single-cycle operations op1 to op%d with dependences %s. The latency bound is %d steps (steps 1 to %d). What is the mobility, ALAP minus ASAP, of op%d?", n, n, es, L, L, k+1)+roundNote(0),
			edakit.IntStr(mob),
			edakit.Sprintf("ASAP times are %s and ALAP times under the bound are %s. For op%d the window is [%d, %d], mobility %d. Mobility sets how many $x_{i,t}$ binaries the ILP needs.", ints(asap, ", "), ints(alap, ", "), k+1, asap[k], alap[k], mob),
			"asap", edakit.IntStr(asap[k]), "alap", edakit.IntStr(alap[k]))
	default:
		n := edakit.Scaled(ctx, 4, 8)
		var parts []string
		use := map[int]int{}
		for i := 0; i < n; i++ {
			s := edakit.Int(1, 6)
			d := edakit.Int(1, 3)
			for t := s; t < s+d; t++ {
				use[t]++
			}
			parts = append(parts, edakit.Sprintf("mul%d starts in step %d and occupies %d step(s)", i+1, s, d))
		}
		mx := 0
		for _, u := range use {
			if u > mx {
				mx = u
			}
		}
		return num("An ILP schedule's resource constraint is $\\sum_i \\sum_{\\tau = t-d_i+1}^{t} x_{i,\\tau} \\le a_{mul}$ for every step t. The solution assigns: "+strings.Join(parts, "; ")+". What is the smallest multiplier count $a_{mul}$ this schedule satisfies?"+roundNote(0),
			edakit.IntStr(mx),
			edakit.Sprintf("Count the multipliers busy in each step, including multi-cycle ones started earlier. The peak is %d, so at least %d multipliers are needed.", mx, mx),
			"peak", edakit.IntStr(mx))
	}
}

func globalRouting(ctx generator.GeneratorContext) generator.Problem {
	for {
		E := edakit.Scaled(ctx, 4, 6)
		N := edakit.Scaled(ctx, 3, 4)
		capE := make([]int, E)
		for i := range capE {
			capE[i] = edakit.Int(1, 2)
		}
		type cand struct {
			edges []int
			wl    int
		}
		cands := make([][2]cand, N)
		for n := 0; n < N; n++ {
			for c := 0; c < 2; c++ {
				k := edakit.Int(1, 3)
				es := distinctInts(k, 0, E-1)
				sort.Ints(es)
				cands[n][c] = cand{es, k*10 + edakit.Int(0, 6)}
			}
		}
		edgeStr := func(es []int) string {
			var p []string
			for _, e := range es {
				p = append(p, edakit.Sprintf("e%d", e+1))
			}
			return strings.Join(p, "+")
		}
		var desc []string
		for n := 0; n < N; n++ {
			desc = append(desc, edakit.Sprintf("net %d: route A uses %s (length %d), route B uses %s (length %d)", n+1, edgeStr(cands[n][0].edges), cands[n][0].wl, edgeStr(cands[n][1].edges), cands[n][1].wl))
		}
		var caps []string
		for i, c := range capE {
			caps = append(caps, edakit.Sprintf("e%d = %d", i+1, c))
		}
		stem := "A global router picks one candidate route per net. Edge capacities: " + strings.Join(caps, ", ") + ". Candidates: " + strings.Join(desc, "; ") + "."
		overflow := func(choice []int) int {
			use := make([]int, E)
			for n, c := range choice {
				for _, e := range cands[n][c].edges {
					use[e]++
				}
			}
			o := 0
			for i := range use {
				if use[i] > capE[i] {
					o += use[i] - capE[i]
				}
			}
			return o
		}
		switch kp(ctx) {
		case 1:
			choice := make([]int, N)
			var names []string
			for n := range choice {
				choice[n] = edakit.Int(0, 1)
				names = append(names, edakit.Sprintf("net %d takes %s", n+1, []string{"A", "B"}[choice[n]]))
			}
			o := overflow(choice)
			return num(stem+" The current solution: "+strings.Join(names, ", ")+". What is the total overflow, the sum over edges of usage above capacity?"+roundNote(0),
				edakit.IntStr(o),
				edakit.Sprintf("Count how many chosen routes cross each edge, subtract the capacity where usage exceeds it, and add. The total overflow is %d.", o),
				"overflow", edakit.IntStr(o))
		case 2:
			best := -1
			for m := 0; m < 1<<N; m++ {
				choice := make([]int, N)
				wl := 0
				for n := 0; n < N; n++ {
					choice[n] = (m >> n) & 1
					wl += cands[n][choice[n]].wl
				}
				if overflow(choice) == 0 && (best < 0 || wl < best) {
					best = wl
				}
			}
			if best < 0 {
				continue
			}
			return num(stem+" Among assignments with zero overflow, what is the minimum total wirelength?"+roundNote(0),
				edakit.IntStr(best),
				edakit.Sprintf("This is the 0-1 program: minimize total length subject to one route per net and edge usage within capacity. Checking all $2^{%d}$ assignments, the best overflow-free one has length %d.", N, best),
				"wl", edakit.IntStr(best))
		default:
			e := edakit.Int(0, E-1)
			xs := make([]float64, N)
			exp := 0.0
			var fr []string
			for n := 0; n < N; n++ {
				xs[n] = float64(edakit.Int(1, 9)) / 10
				fr = append(fr, edakit.Sprintf("net %d: A = %s, B = %s", n+1, f(xs[n], 1), f(1-xs[n], 1)))
				for _, ee := range cands[n][0].edges {
					if ee == e {
						exp += xs[n]
					}
				}
				for _, ee := range cands[n][1].edges {
					if ee == e {
						exp += 1 - xs[n]
					}
				}
			}
			return num(stem+" The LP relaxation returns fractional route weights "+strings.Join(fr, "; ")+edakit.Sprintf(". Randomized rounding picks each net's route with these probabilities. What is the expected usage of edge e%d?", e+1)+roundNote(2),
				f(exp, 2),
				edakit.Sprintf("By linearity of expectation, the expected usage of e%d is the sum of the weights of the candidate routes that cross it: %s. This equals the LP's own usage of the edge, so rounding preserves capacity in expectation.", e+1, f(exp, 2)),
				"expected_usage", f(exp, 2))
		}
	}
}

func cuttingPlanes(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		for {
			n := edakit.Scaled(ctx, 4, 6)
			a := make([]int, n)
			sum := 0
			for i := range a {
				a[i] = edakit.Int(2, 9)
				sum += a[i]
			}
			b := sum*2/3 + edakit.Int(-1, 1)
			// a cover: indices whose weights exceed b
			idx := edakit.Shuffle(rangeInts(0, n-1))
			var cov []int
			w := 0
			for _, i := range idx {
				cov = append(cov, i)
				w += a[i]
				if w > b {
					break
				}
			}
			if w <= b || len(cov) < 2 {
				continue
			}
			sort.Ints(cov)
			x := make([]float64, n)
			lhs := 0.0
			for i := range x {
				x[i] = float64(edakit.Int(0, 10)) / 10
				lhs += float64(a[i]) * x[i]
			}
			if lhs > float64(b) {
				continue
			}
			s := 0.0
			for _, i := range cov {
				s += x[i]
			}
			viol := s - float64(len(cov)-1)
			if viol <= 1e-9 {
				continue
			}
			var terms, xs, cs []string
			for i := range a {
				terms = append(terms, edakit.Sprintf("%dx_%d", a[i], i+1))
				xs = append(xs, edakit.Sprintf("$x_%d = %s$", i+1, f(x[i], 1)))
			}
			for _, i := range cov {
				cs = append(cs, edakit.Sprintf("x_%d", i+1))
			}
			return num(edakit.Sprintf("A binary knapsack row from a placement-legalization ILP reads $%s \\le %d$. The LP relaxation returns ", strings.Join(terms, " + "), b)+strings.Join(xs, ", ")+edakit.Sprintf(". The set {%s} is a cover (its weights exceed %d), giving the cover inequality $%s \\le %d$. By how much does the LP point violate it?", strings.Join(cs, ", "), b, strings.Join(cs, " + "), len(cov)-1)+roundNote(2),
				f(viol, 2),
				edakit.Sprintf("No integer solution can take every item of a cover, so at most %d of them are 1. The LP point sums to %s over the cover, exceeding %d by %s, so the cut removes it.", len(cov)-1, f(s, 2), len(cov)-1, f(viol, 2)),
				"violation", f(viol, 2))
		}
	case 2:
		a1, a2 := edakit.Int(2, 9), edakit.Int(2, 9)
		r := edakit.Int(5, 30)
		k := edakit.Int(2, 4)
		c1, c2, rr := a1/k, a2/k, r/k
		ask := edakit.Int(0, 1)
		ans := rr
		what := "right-hand side"
		if ask == 1 {
			ans = c1
			what = "coefficient of x"
		}
		return num(edakit.Sprintf("For integers $x, y \\ge 0$ satisfying $%dx + %dy \\le %d$, the Chvatal-Gomory procedure multiplies by $u = 1/%d$ and rounds down every coefficient and the right-hand side. What is the %s of the resulting cut?", a1, a2, r, k, what)+roundNote(0),
			edakit.IntStr(ans),
			edakit.Sprintf("Scaling gives $%s x + %s y \\le %s$. Rounding the left coefficients down keeps the inequality valid because x and y are nonnegative, and then the left side is an integer, so the right side can be rounded down: $%s \\le %d$.", f(float64(a1)/float64(k), 3), f(float64(a2)/float64(k), 3), f(float64(r)/float64(k), 3), cgLHS(c1, c2), rr),
			"cut_rhs", edakit.IntStr(rr), "cut_x", edakit.IntStr(c1))
	default:
		a1, a2 := edakit.Int(1, 5), edakit.Int(1, 5)
		r := edakit.Int(6, 20)
		B := edakit.Int(3, 6)
		k := edakit.Int(2, 6)
		cnt, feas := 0, 0
		for x := 0; x <= B; x++ {
			for y := 0; y <= B; y++ {
				if a1*x+a2*y <= r {
					feas++
					if x+y > k {
						cnt++
					}
				}
			}
		}
		return num(edakit.Sprintf("Integer points satisfy $%dx + %dy \\le %d$ and $0 \\le x, y \\le %d$. A cut generator proposes $x + y \\le %d$. How many feasible integer points does the proposed cut wrongly remove? (Zero means the cut is valid.)", a1, a2, r, B, k)+roundNote(0),
			edakit.IntStr(cnt),
			edakit.Sprintf("Enumerate the %d feasible integer points and count those with $x + y > %d$: %d. A valid cut removes only fractional LP points, never a feasible integer one.", feas, k, cnt),
			"removed", edakit.IntStr(cnt))
	}
}

func bendersCut(ctx generator.GeneratorContext) generator.Problem {
	u1, u2 := edakit.Int(0, 4), edakit.Int(1, 4)
	h1, h2 := edakit.Int(5, 20), edakit.Int(5, 20)
	g1, g2 := edakit.Int(1, 4), edakit.Int(1, 4)
	y0 := edakit.Int(0, 3)
	val := u1*(h1-g1*y0) + u2*(h2-g2*y0)
	wrongSign := u1*(h1+g1*y0) + u2*(h2+g2*y0)
	noY := u1*h1 + u2*h2
	oneTerm := u2 * (h2 - g2*y0)
	q := edakit.Sprintf("In a Benders scheme the master problem picks an integer macro decision y; the LP subproblem $\\min\\{c^T x : Wx \\ge h - Gy\\}$ has right-hand side $h = (%d, %d)$ and $G = (%d, %d)^T$. At an earlier master point the subproblem dual returned $u = (%d, %d)$, giving the optimality cut $\\theta \\ge u^T(h - Gy)$. What lower bound does this cut put on $\\theta$ at $y = %d$?", h1, h2, g1, g2, u1, u2, y0)
	return edakit.Choice(q, edakit.IntStr(val),
		edakit.Sprintf("Evaluate $u^T(h - Gy) = %d(%d - %d \\cdot %d) + %d(%d - %d \\cdot %d) = %d$. Any dual-feasible u gives a valid bound for every y by weak duality.", u1, h1, g1, y0, u2, h2, g2, y0, val),
		edakit.Distinct(val, 3, wrongSign, noY, oneTerm)...)
}

// cgLHS renders c1*x + c2*y, dropping zero terms and bare unit coefficients.
func cgLHS(c1, c2 int) string {
	var parts []string
	for _, t := range []struct {
		c int
		v string
	}{{c1, "x"}, {c2, "y"}} {
		switch t.c {
		case 0:
		case 1:
			parts = append(parts, t.v)
		default:
			parts = append(parts, edakit.IntStr(t.c)+t.v)
		}
	}
	if len(parts) == 0 {
		return "0"
	}
	return strings.Join(parts, " + ")
}
