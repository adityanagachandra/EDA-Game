package hls

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

const latencyRule = "Additions take 1 control step and multiplications take 2. Steps are numbered from 1; an operation that starts in step s with latency l occupies steps s through s+l-1, and a consumer may start in step s+l."

const unitRule = "Every operation takes one control step, and steps are numbered from 1."

// pickInner returns an operation index that has at least one predecessor when possible.
func pickInner(g *dfg) int {
	var cand []int
	for i, nd := range g.ops {
		if len(nd.preds) > 0 {
			cand = append(cand, i)
		}
	}
	if len(cand) == 0 {
		return len(g.ops) - 1
	}
	return edakit.Pick(cand)
}

func asapAlap(ctx generator.GeneratorContext) generator.Problem {
	g := newDFG(edakit.Scaled(ctx, 5, 10), 0.4)
	lat := g.lat(1, 2)
	as := g.asap(lat)
	crit := length(as, lat)
	switch whichKP(ctx) {
	case 1:
		x := pickInner(g)
		q := "Data-flow graph: " + g.text() + ". " + latencyRule +
			edakit.Sprintf(" What is the ASAP start step of %s?", opName(x))
		return num(q, edakit.IntStr(as[x]),
			edakit.Sprintf("ASAP starts every source in step 1 and every other operation at the maximum over its predecessors of start plus latency. The ASAP starts are %s, so %s starts in step %d.", joinInts(as), opName(x), as[x]),
			map[string]string{"asap": joinInts(as)})
	case 2:
		q := "Data-flow graph: " + g.text() + ". " + latencyRule +
			" With unlimited functional units, what is the minimum number of control steps the block needs?"
		return num(q, edakit.IntStr(crit),
			edakit.Sprintf("The ASAP starts are %s. The schedule ends when the last operation finishes, at step %d, which is the latency-weighted critical path.", joinInts(as), crit),
			map[string]string{"asap": joinInts(as), "length": edakit.IntStr(crit)})
	default:
		L := crit + edakit.Int(0, 2)
		al := g.alap(lat, L)
		x := edakit.Int(0, len(g.ops)-1)
		q := "Data-flow graph: " + g.text() + ". " + latencyRule +
			edakit.Sprintf(" The block must finish within %d control steps. What is the ALAP start step of %s?", L, opName(x))
		return num(q, edakit.IntStr(al[x]),
			edakit.Sprintf("ALAP works backwards: a sink starts at L - l + 1, and every other operation at the minimum over its successors of their ALAP start minus its own latency. With L = %d the ALAP starts are %s, so %s starts in step %d.", L, joinInts(al), opName(x), al[x]),
			map[string]string{"alap": joinInts(al)})
	}
}

func mobility(ctx generator.GeneratorContext) generator.Problem {
	g := newDFG(edakit.Scaled(ctx, 5, 10), 0.4)
	lat := g.lat(1, 2)
	as := g.asap(lat)
	crit := length(as, lat)
	al := g.alap(lat, crit)
	mob := make([]int, len(as))
	for i := range as {
		mob[i] = al[i] - as[i]
	}
	switch whichKP(ctx) {
	case 1:
		x := edakit.Int(0, len(g.ops)-1)
		q := "Data-flow graph: " + g.text() + ". " + latencyRule +
			edakit.Sprintf(" The latency bound equals the critical-path length, %d steps. What is the mobility (ALAP start minus ASAP start) of %s?", crit, opName(x))
		return num(q, edakit.IntStr(mob[x]),
			edakit.Sprintf("ASAP starts: %s. ALAP starts with L = %d: %s. For %s the mobility is %d - %d = %d.", joinInts(as), crit, joinInts(al), opName(x), al[x], as[x], mob[x]),
			map[string]string{"asap": joinInts(as), "alap": joinInts(al)})
	case 2:
		z := 0
		for _, m := range mob {
			if m == 0 {
				z++
			}
		}
		q := "Data-flow graph: " + g.text() + ". " + latencyRule +
			edakit.Sprintf(" The latency bound equals the critical-path length, %d steps. How many operations have zero mobility?", crit)
		return num(q, edakit.IntStr(z),
			edakit.Sprintf("ASAP starts: %s. ALAP starts: %s. The operations with equal ASAP and ALAP starts lie on a critical path; there are %d of them.", joinInts(as), joinInts(al), z),
			map[string]string{"zero_mobility": edakit.IntStr(z)})
	default:
		r := edakit.Int(1, 3)
		x := edakit.Int(0, len(g.ops)-1)
		q := "Data-flow graph: " + g.text() + ". " + latencyRule +
			edakit.Sprintf(" The latency bound is relaxed to %d steps. What is the mobility of %s?", crit+r, opName(x))
		return num(q, edakit.IntStr(mob[x]+r),
			edakit.Sprintf("The critical path is %d steps, so the bound adds %d steps of slack. ALAP shifts every start later by exactly %d while ASAP does not move, so %s has mobility %d + %d = %d.", crit, r, r, opName(x), mob[x], r, mob[x]+r),
			map[string]string{"critical": edakit.IntStr(crit)})
	}
}

func resourceBound(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		n := edakit.Int(5, 30)
		k := edakit.Int(2, 5)
		b := ceilDiv(n, k)
		q := edakit.Sprintf("A block contains %d single-cycle multiplications and the allocation provides %d multipliers. Whatever the dependences, what is the smallest number of control steps the multiplications alone force?", n, k)
		return num(q, edakit.IntStr(b),
			edakit.Sprintf("Each step serves at most %d multiplications, so at least the ceiling of %d / %d = %d steps are needed.", k, n, k, b),
			map[string]string{"bound": edakit.IntStr(b)})
	case 2:
		n := edakit.Int(4, 16)
		k := edakit.Int(1, 3)
		l := edakit.Int(2, 4)
		b := ceilDiv(n*l, k)
		q := edakit.Sprintf("A block contains %d multiplications. Each takes %d control steps on a non-pipelined multiplier, which stays busy for all of them. There are %d multipliers. What lower bound on the number of control steps do the multiplications force?", n, l, k)
		return num(q, edakit.IntStr(b),
			edakit.Sprintf("The multiplications need %d × %d = %d unit-steps of multiplier time and %d units supply %d per step, so at least the ceiling of %d / %d = %d steps.", n, l, n*l, k, k, n*l, k, b),
			map[string]string{"bound": edakit.IntStr(b)})
	default:
		g := newDFG(edakit.Scaled(ctx, 7, 12), 0.45)
		lat := g.lat(1, 1)
		cp := maxInt(g.heights(lat)...)
		ka, km := edakit.Int(1, 2), edakit.Int(1, 2)
		na, nm := g.count(opAdd), g.count(opMul)
		ba, bm := ceilDiv(na, ka), ceilDiv(nm, km)
		b := maxInt(cp, ba, bm)
		q := "Data-flow graph: " + g.text() + ". " + unitRule +
			edakit.Sprintf(" There are %d adders and %d multipliers. What is the largest of the simple lower bounds on schedule length: the critical path and, for each unit type, the operation count divided by the unit count, rounded up?", ka, km)
		return num(q, edakit.IntStr(b),
			edakit.Sprintf("The critical path is %d steps. %d additions on %d adders need %d steps; %d multiplications on %d multipliers need %d. The bound is the maximum, %d.", cp, na, ka, ba, nm, km, bm, b),
			map[string]string{"critical": edakit.IntStr(cp), "add_bound": edakit.IntStr(ba), "mul_bound": edakit.IntStr(bm)})
	}
}

func listSched(ctx generator.GeneratorContext) generator.Problem {
	g := newDFG(edakit.Scaled(ctx, 6, 11), 0.45)
	k := whichKP(ctx)
	rule := "Priority is the length of the longest path from the operation to a sink, counted in steps and including the operation itself; ties go to the lower-numbered operation. In each step the scheduler fills free units with ready operations in priority order."
	if k == 1 {
		lat := g.lat(1, 1)
		h := g.heights(lat)
		x := edakit.Int(0, len(g.ops)-2)
		q := "Data-flow graph: " + g.text() + ". " + unitRule + " A list scheduler uses as its priority the length of the longest path from an operation to a sink, counted in operations and including the operation itself. What is the priority of " + opName(x) + "?"
		return num(q, edakit.IntStr(h[x]),
			edakit.Sprintf("Computing heights from the sinks backwards gives %s, so %s has priority %d.", joinInts(h), opName(x), h[x]),
			map[string]string{"heights": joinInts(h)})
	}
	lm := 1
	latText := unitRule
	if k == 3 {
		lm = 2
		latText = latencyRule + " Multipliers are not pipelined: a multiplier is busy for both steps of its operation."
	}
	lat := g.lat(1, lm)
	units := [2]int{edakit.Int(1, 2), 1}
	if k == 2 {
		units[1] = edakit.Int(1, 2)
	}
	st := g.listSchedule(lat, units)
	L := length(st, lat)
	q := "Data-flow graph: " + g.text() + ". " + latText +
		edakit.Sprintf(" The allocation is %d adder(s) and %d multiplier(s). ", units[0], units[1]) + rule +
		" How many control steps does the resulting schedule take?"
	return num(q, edakit.IntStr(L),
		edakit.Sprintf("Running the scheduler step by step gives start steps %s (%s). The last operation finishes in step %d.", joinInts(st), stepTable(st), L),
		map[string]string{"starts": joinInts(st), "length": edakit.IntStr(L)})
}

func minResources(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		g := newDFG(edakit.Scaled(ctx, 6, 11), 0.5)
		lat := g.lat(1, 1)
		as := g.asap(lat)
		kind := edakit.Pick([]int{opAdd, opMul})
		need := perStepMax(g, as, kind)
		q := "Data-flow graph: " + g.text() + ". " + unitRule + " The block is scheduled ASAP: " + stepTable(as) +
			edakit.Sprintf(". How many %ss does this exact schedule need?", unitName(kind))
		return num(q, edakit.IntStr(need),
			edakit.Sprintf("A fixed schedule needs as many units of a type as the largest number of operations of that type sharing one step. Counting %ss per step, the maximum is %d.", kindName(kind), need),
			map[string]string{"units": edakit.IntStr(need)})
	case 2:
		n := edakit.Int(6, 30)
		L := edakit.Int(3, 9)
		b := ceilDiv(n, L)
		q := edakit.Sprintf("A latency-constrained schedule must fit %d single-cycle multiplications into %d control steps. Ignoring dependences, what is the fewest multipliers that could possibly suffice?", n, L)
		return num(q, edakit.IntStr(b),
			edakit.Sprintf("%d multipliers serve at most %d operations in %d steps; the smallest count that covers %d is the ceiling of %d / %d = %d.", b, b*L, L, n, n, L, b),
			map[string]string{"bound": edakit.IntStr(b)})
	default:
		g := newDFG(edakit.Scaled(ctx, 6, 11), 0.5)
		lat := g.lat(1, 1)
		cp := maxInt(g.heights(lat)...)
		L := cp + edakit.Int(0, 2)
		al := g.alap(lat, L)
		kind := edakit.Pick([]int{opAdd, opMul})
		need := perStepMax(g, al, kind)
		q := "Data-flow graph: " + g.text() + ". " + unitRule +
			edakit.Sprintf(" The latency bound is %d steps and every operation is placed at its ALAP step. How many %ss does that schedule need?", L, unitName(kind))
		return num(q, edakit.IntStr(need),
			edakit.Sprintf("ALAP starts with L = %d are %s (%s). The largest number of %ss in one step is %d.", L, joinInts(al), stepTable(al), kindName(kind), need),
			map[string]string{"alap": joinInts(al), "units": edakit.IntStr(need)})
	}
}

func perStepMax(g *dfg, start []int, kind int) int {
	c := map[int]int{}
	m := 0
	for i, s := range start {
		if g.ops[i].kind == kind {
			c[s]++
			if c[s] > m {
				m = c[s]
			}
		}
	}
	return m
}

type frame struct{ a, b int }

func frameText(fs []frame) string {
	var p []string
	for i, f := range fs {
		p = append(p, edakit.Sprintf("m%d: [%d, %d]", i+1, f.a, f.b))
	}
	return strings.Join(p, ", ")
}

func dg(fs []frame, s int) float64 {
	t := 0.0
	for _, f := range fs {
		if s >= f.a && s <= f.b {
			t += 1 / float64(f.b-f.a+1)
		}
	}
	return t
}

func forceDirected(ctx generator.GeneratorContext) generator.Problem {
	L := edakit.Int(4, 6)
	m := edakit.Scaled(ctx, 3, 6)
	fs := make([]frame, m)
	for i := range fs {
		a := edakit.Int(1, L)
		b := minInt(L, a+edakit.Int(0, 3))
		fs[i] = frame{a, b}
	}
	if fs[0].b == fs[0].a {
		fs[0] = frame{edakit.Int(1, L-2), L - edakit.Int(0, 1)}
	}
	head := edakit.Sprintf("Force-directed scheduling with a latency bound of %d steps. The multiplications have time frames [ASAP, ALAP]: %s. Each operation is equally likely to land in any step of its frame.", L, frameText(fs))
	switch whichKP(ctx) {
	case 1:
		x := edakit.Int(0, m-1)
		f := fs[x]
		s := edakit.Int(f.a, f.b)
		if edakit.Int(0, 4) == 0 {
			s = edakit.Int(1, L)
		}
		p := 0.0
		if s >= f.a && s <= f.b {
			p = 1 / float64(f.b-f.a+1)
		}
		q := head + edakit.Sprintf(" What is the probability that m%d executes in step %d? Round to 2 decimal places.", x+1, s)
		return num(q, edakit.Num(p, 2),
			edakit.Sprintf("m%d has frame [%d, %d] of width %d, so each step inside it gets probability 1/%d and steps outside get 0. Step %d gives %s.", x+1, f.a, f.b, f.b-f.a+1, f.b-f.a+1, s, edakit.Num(p, 2)),
			map[string]string{"p": edakit.Num(p, 2)})
	case 2:
		s := edakit.Int(1, L)
		v := dg(fs, s)
		q := head + edakit.Sprintf(" What is the value of the multiplier distribution graph DG(%d)? Round to 2 decimal places.", s)
		return num(q, edakit.Num(v, 2),
			edakit.Sprintf("DG(%d) sums, over the multiplications whose frame contains step %d, one over the frame width. The sum is %s.", s, s, edakit.Num(v, 2)),
			map[string]string{"dg": edakit.Num(v, 2)})
	default:
		var wide []int
		for i, f := range fs {
			if f.b > f.a {
				wide = append(wide, i)
			}
		}
		x := edakit.Pick(wide)
		f := fs[x]
		j := edakit.Int(f.a, f.b)
		w := float64(f.b - f.a + 1)
		sum := 0.0
		var dgs []string
		for s := f.a; s <= f.b; s++ {
			sum += dg(fs, s)
			dgs = append(dgs, edakit.Sprintf("DG(%d) = %s", s, edakit.Num(dg(fs, s), 3)))
		}
		force := dg(fs, j) - sum/w
		q := head + edakit.Sprintf(" What is the self-force of fixing m%d in step %d, defined as the sum over its frame of DG(s) times the change in its probability at s? Round to 2 decimal places.", x+1, j)
		return num(q, edakit.Num(force, 2),
			edakit.Sprintf("Inside m%d's frame, %s. Fixing it in step %d raises its probability there to 1 and drops every step from 1/%d. The force is DG(%d) minus the frame average of DG: %s.", x+1, strings.Join(dgs, ", "), j, int(w), j, edakit.Num(force, 2)),
			map[string]string{"force": edakit.Num(force, 2)})
	}
}

func ilpCount(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 8)
		L := edakit.Int(4, 7)
		fs := make([]frame, n)
		total := 0
		var p []string
		for i := range fs {
			a := edakit.Int(1, L)
			b := minInt(L, a+edakit.Int(0, 3))
			fs[i] = frame{a, b}
			total += b - a + 1
			p = append(p, edakit.Sprintf("o%d: [%d, %d]", i+1, a, b))
		}
		q := edakit.Sprintf("A time-indexed ILP schedule uses a binary variable x(i, s) = 1 when operation i starts in step s, and creates variables only for steps inside each operation's [ASAP, ALAP] frame. With latency bound %d the frames are %s. How many binary variables does the model have?", L, strings.Join(p, ", "))
		return num(q, edakit.IntStr(total),
			edakit.Sprintf("Each operation contributes one variable per step of its frame, its mobility plus one. Summing the widths gives %d variables, against %d for the naive model over all %d steps.", total, n*L, L),
			map[string]string{"vars": edakit.IntStr(total)})
	case 2:
		n := edakit.Int(5, 20)
		e := edakit.Int(n-1, 2*n)
		L := edakit.Int(4, 10)
		T := edakit.Int(2, 3)
		total := n + e + L*T
		q := edakit.Sprintf("A time-indexed ILP schedule has %d operations, %d dependence edges, a latency bound of %d steps and %d resource types. It has one unique-start constraint per operation, one precedence constraint per edge (written on the start times sum s · x(i, s)) and one resource constraint per step per resource type. How many constraints does it have?", n, e, L, T)
		return num(q, edakit.IntStr(total),
			edakit.Sprintf("%d unique-start + %d precedence + %d × %d resource constraints = %d.", n, e, L, T, total),
			map[string]string{"constraints": edakit.IntStr(total)})
	default:
		m := edakit.Scaled(ctx, 4, 6)
		T := edakit.Int(5, 9)
		ds := make([]int, m)
		for i := range ds {
			ds[i] = edakit.Int(1, 5)
		}
		cnt := 0
		var pairs []string
		for i := 0; i < m; i++ {
			s := 0
			for j := i; j < m; j++ {
				s += ds[j]
				if j > i && s > T {
					cnt++
					pairs = append(pairs, edakit.Sprintf("(o%d, o%d)", i+1, j+1))
				}
			}
		}
		q := edakit.Sprintf("An SDC scheduler handles chaining on a chain o1 → o2 → ... → o%d of operations with combinational delays %s ns and a clock period of %d ns. For every pair (oi, oj), i < j, whose total delay from the start of oi to the end of oj exceeds the period, it adds the difference constraint s_j - s_i ≥ 1 so the two cannot share a cycle. How many such constraints does it add?", m, joinInts(ds), T)
		ex := "No pair exceeds the period."
		if cnt > 0 {
			ex = "The pairs whose summed delay exceeds the period are " + strings.Join(pairs, ", ") + "."
		}
		return num(q, edakit.IntStr(cnt), ex+edakit.Sprintf(" That is %d constraints.", cnt),
			map[string]string{"constraints": edakit.IntStr(cnt)})
	}
}

type dc struct{ from, to, w int } // s_to - s_from >= w

// longestFeasible solves s >= 0 under difference constraints by Bellman-Ford longest paths.
func longestFeasible(n int, cs []dc) ([]int, bool) {
	s := make([]int, n)
	for it := 0; it <= n+1; it++ {
		changed := false
		for _, c := range cs {
			if s[c.from]+c.w > s[c.to] {
				s[c.to] = s[c.from] + c.w
				changed = true
			}
		}
		if !changed {
			return s, true
		}
	}
	return s, false
}

func dcText(cs []dc) string {
	var p []string
	for _, c := range cs {
		if c.w >= 0 {
			p = append(p, edakit.Sprintf("s%d - s%d ≥ %d", c.to+1, c.from+1, c.w))
		} else {
			p = append(p, edakit.Sprintf("s%d - s%d ≤ %d", c.from+1, c.to+1, -c.w))
		}
	}
	return strings.Join(p, "; ")
}

func sdcSolve(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 6)
	var cs []dc
	for j := 1; j < n; j++ {
		k := 1
		if j >= 2 && edakit.Int(0, 1) == 1 {
			k = 2
		}
		used := map[int]bool{}
		for len(used) < k {
			i := edakit.Int(0, j-1)
			if !used[i] {
				used[i] = true
				cs = append(cs, dc{i, j, edakit.Int(1, 3)})
			}
		}
	}
	base, _ := longestFeasible(n, cs)
	switch whichKP(ctx) {
	case 1:
		x := edakit.Int(1, n-1)
		q := edakit.Sprintf("An SDC schedule has start-time variables s1 to s%d, all at least 0, and the difference constraints %s. What is the smallest feasible value of s%d?", n, dcText(cs), x+1)
		return num(q, edakit.IntStr(base[x]),
			edakit.Sprintf("Each constraint s_j - s_i ≥ w is an edge i → j of weight w, and the least solution is the longest path from a source at 0. The least starts are %s, so s%d = %d.", joinInts(base), x+1, base[x]),
			map[string]string{"starts": joinInts(base)})
	case 2:
		// A maximum-latency constraint s_j - s_i <= c with c at least the longest i-to-j path.
		// Prefer a pair whose latency constraint actually drags s_i later.
		i, j, lp := 1, n-1, 0
		for try := 0; try < 30; try++ {
			i = edakit.Int(0, n-2)
			j = edakit.Int(i+1, n-1)
			lp = pathLongest(n, cs, i, j)
			if lp < 0 {
				lp = 0
			}
			if base[j]-base[i] > lp {
				break
			}
		}
		c := lp + edakit.Int(0, 1)
		if gap := base[j] - base[i]; gap > lp {
			c = edakit.Int(lp, gap-1)
		}
		all := append(append([]dc(nil), cs...), dc{j, i, -c})
		s, ok := longestFeasible(n, all)
		if !ok {
			s = base
		}
		q := edakit.Sprintf("An SDC schedule has start-time variables s1 to s%d, all at least 0, and the constraints %s. The last one is a maximum-latency constraint. What is the smallest feasible value of s%d?", n, dcText(all), i+1)
		return num(q, edakit.IntStr(s[i]),
			edakit.Sprintf("Rewrite s%d - s%d ≤ %d as s%d - s%d ≥ %d, a backward edge of weight %d. Longest paths with it give %s, so s%d = %d (it was %d without the latency constraint).", j+1, i+1, c, i+1, j+1, -c, -c, joinInts(s), i+1, s[i], base[i]),
			map[string]string{"starts": joinInts(s)})
	default:
		i := edakit.Int(0, n-2)
		j := edakit.Int(i+1, n-1)
		lp := pathLongest(n, cs, i, j)
		for lp < 0 {
			cs = append(cs, dc{i, j, edakit.Int(1, 3)})
			lp = pathLongest(n, cs, i, j)
		}
		q := edakit.Sprintf("An SDC schedule has start-time variables s1 to s%d, all at least 0, and the constraints %s. A designer adds the maximum-latency constraint s%d - s%d ≤ c. What is the smallest integer c that keeps the system feasible?", n, dcText(cs), j+1, i+1)
		return num(q, edakit.IntStr(lp),
			edakit.Sprintf("The new constraint is an edge s%d → s%d of weight -c. It closes a cycle with the longest path s%d → s%d, of weight %d, and a system of difference constraints is feasible exactly when no cycle has positive weight. So c must be at least %d.", j+1, i+1, i+1, j+1, lp, lp),
			map[string]string{"longest": edakit.IntStr(lp)})
	}
}

// pathLongest is the longest-path weight from i to j over forward constraints, or -1.
func pathLongest(n int, cs []dc, i, j int) int {
	d := make([]int, n)
	for k := range d {
		d[k] = -1 << 30
	}
	d[i] = 0
	for it := 0; it < n; it++ {
		for _, c := range cs {
			if c.w >= 0 && d[c.from] > -1<<29 && d[c.from]+c.w > d[c.to] {
				d[c.to] = d[c.from] + c.w
			}
		}
	}
	if d[j] < -1<<29 {
		return -1
	}
	return d[j]
}

func chaining(ctx generator.GeneratorContext) generator.Problem {
	T := edakit.Pick([]int{40, 50, 60, 80, 100})
	r := edakit.Pick([]int{3, 5, 8, 10})
	use := T - r
	m := edakit.Scaled(ctx, 4, 8)
	ds := make([]int, m)
	var dstr []string
	for i := range ds {
		ds[i] = edakit.Int(5, use*2/3)
		dstr = append(dstr, edakit.Num(float64(ds[i])/10, 1))
	}
	cycles, acc, first := 1, 0, 0
	firstDone := false
	for _, d := range ds {
		if acc+d > use {
			cycles++
			acc = d
			firstDone = true
		} else {
			acc += d
			if !firstDone {
				first++
			}
		}
	}
	head := edakit.Sprintf("A chain of %d dependent operations has combinational delays %s ns. The clock period is %s ns and each register stage costs %s ns of clock-to-q plus setup, so at most %s ns of logic fits in one cycle. Operations are chained greedily: the next operation joins the current cycle if it still fits, and otherwise starts a new cycle.",
		m, strings.Join(dstr, ", "), edakit.Num(float64(T)/10, 1), edakit.Num(float64(r)/10, 1), edakit.Num(float64(use)/10, 1))
	switch whichKP(ctx) {
	case 1:
		return num(head+" How many operations are chained into the first cycle?", edakit.IntStr(first),
			edakit.Sprintf("Adding delays from the start of the chain while the running sum stays within %s ns admits %d operations.", edakit.Num(float64(use)/10, 1), first),
			map[string]string{"first": edakit.IntStr(first)})
	case 2:
		return num(head+" How many cycles does the chain take?", edakit.IntStr(cycles),
			edakit.Sprintf("Greedy packing of a single chain is optimal, and it uses %d cycles.", cycles),
			map[string]string{"cycles": edakit.IntStr(cycles)})
	default:
		return num(head+" Without chaining every operation would take its own cycle. How many cycles does chaining save?", edakit.IntStr(m-cycles),
			edakit.Sprintf("Unchained, the chain takes %d cycles; chained, %d. The saving is %d cycles.", m, cycles, m-cycles),
			map[string]string{"cycles": edakit.IntStr(cycles)})
	}
}

func multicycleFU(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 20)
	u := edakit.Int(1, 3)
	L := edakit.Int(2, 6)
	b := ceilDiv(n, u)
	switch whichKP(ctx) {
	case 1:
		ans := b * L
		q := edakit.Sprintf("%d independent multiplications are all ready in step 1. There are %d multipliers, each with a latency of %d steps and not pipelined, so a unit accepts a new operation only after the previous one finishes. In which step does the last result become available (the step in which the last multiplication completes)?", n, u, L)
		return num(q, edakit.IntStr(ans),
			edakit.Sprintf("The units process the operations in %d rounds of %d steps each, so the last one completes in step %d × %d = %d.", b, L, b, L, ans),
			map[string]string{"rounds": edakit.IntStr(b)})
	case 2:
		ans := b - 1 + L
		q := edakit.Sprintf("%d independent multiplications are all ready in step 1. There are %d multipliers, each with a latency of %d steps and fully pipelined, so a unit accepts a new operation every step. In which step does the last multiplication complete?", n, u, L)
		return num(q, edakit.IntStr(ans),
			edakit.Sprintf("Each unit issues one operation per step, so the last of %d issue rounds starts in step %d and completes %d - 1 steps later, in step %d.", b, b, L, ans),
			map[string]string{"rounds": edakit.IntStr(b)})
	default:
		p := edakit.Int(2, 3)
		if p >= L {
			L = p + 1
		}
		ans := (b-1)*p + L
		q := edakit.Sprintf("%d independent divisions are all ready in step 1. There are %d dividers, each with a latency of %d steps and partially pipelined: a divider accepts a new operation every %d steps. In which step does the last division complete?", n, u, L, p)
		return num(q, edakit.IntStr(ans),
			edakit.Sprintf("Each divider issues every %d steps, so the last of %d rounds starts in step 1 + %d × %d = %d and completes %d - 1 steps later, in step %d.", p, b, b-1, p, 1+(b-1)*p, L, ans),
			map[string]string{"rounds": edakit.IntStr(b)})
	}
}
