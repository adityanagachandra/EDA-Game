package techmap

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func gateSizing(ctx generator.GeneratorContext) generator.Problem {
	for {
		rd := float64(edakit.Int(5, 20)) / 10
		c0 := float64(edakit.Int(1, 4))
		r0 := float64(edakit.Int(10, 40)) / 10
		cl := float64(edakit.Int(10, 120))
		p := float64(edakit.Int(5, 15))
		s := math.Sqrt(r0 * cl / (rd * c0))
		if s < 1.2 || s > 7.5 {
			continue
		}
		D := func(x float64) float64 { return rd*c0*x + p + r0*cl/x }
		q := edakit.Sprintf("A driver of resistance $R_d$ = %s kOhm drives one gate, which drives a fixed load $C_L$ = %s fF. A gate of size s has input capacitance $c_0 s$ with $c_0$ = %s fF, drive resistance $r_0 / s$ with $r_0$ = %s kOhm, and intrinsic delay p = %s ps, so the two-stage delay is $D(s) = R_d c_0 s + p + r_0 C_L / s$ ps.",
			edakit.Num(rd, 1), edakit.Num(cl, 0), edakit.Num(c0, 0), edakit.Num(r0, 1), edakit.Num(p, 0))
		sizes := []float64{1, 2, 4, 8}
		bestS, bestD := 0.0, 1e9
		var tried []string
		for _, x := range sizes {
			tried = append(tried, edakit.Sprintf("D(%s) = %s", edakit.Num(x, 0), edakit.Num(D(x), 2)))
			if D(x) < bestD {
				bestS, bestD = x, D(x)
			}
		}
		switch kp(ctx) {
		case 1:
			return edakit.Numeric(q+" What continuous size s minimizes D? Round to 2 decimal places.", edakit.Num(s, 2),
				edakit.Sprintf("Setting $dD/ds = R_d c_0 - r_0 C_L / s^2$ to zero gives $s = \\sqrt{r_0 C_L / (R_d c_0)}$ = sqrt(%s x %s / (%s x %s)) = %s. D is convex in s, so this is the minimum.", edakit.Num(r0, 1), edakit.Num(cl, 0), edakit.Num(rd, 1), edakit.Num(c0, 0), edakit.Num(s, 2)),
				map[string]string{"s_opt": edakit.Num(s, 2)})
		case 2:
			return edakit.Numeric(q+" The library offers only sizes 1, 2, 4 and 8. What is the smallest achievable D, in ps? Round to 2 decimal places.", edakit.Num(bestD, 2),
				strings.Join(tried, ", ")+edakit.Sprintf(". Size %s wins with %s ps.", edakit.Num(bestS, 0), edakit.Num(bestD, 2)),
				map[string]string{"best_size": edakit.Num(bestS, 0), "best_delay": edakit.Num(bestD, 2)})
		}
		pen := bestD - D(s)
		return edakit.Numeric(q+" The library offers only sizes 1, 2, 4 and 8. How much slower, in ps, is the best discrete size than the continuous optimum? Round to 2 decimal places.", edakit.Num(pen, 2),
			edakit.Sprintf("The continuous optimum s = %s gives D = %s ps. Among the discrete sizes, %s gives %s ps. The rounding penalty is %s ps.", edakit.Num(s, 2), edakit.Num(D(s), 2), edakit.Num(bestS, 0), edakit.Num(bestD, 2), edakit.Num(pen, 2)),
			map[string]string{"continuous": edakit.Num(D(s), 2), "discrete": edakit.Num(bestD, 2), "penalty": edakit.Num(pen, 2)})
	}
}

func lagrangianSizing(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		nin, nout := edakit.Int(2, 3), edakit.Int(1, 3)
		var outs []string
		sumOut := 0.0
		for i := 0; i < nout; i++ {
			x := float64(edakit.Int(5, 40)) / 10
			sumOut += x
			outs = append(outs, edakit.Num(x, 1))
		}
		var ins []string
		sumIn := 0.0
		for i := 0; i < nin-1; i++ {
			x := float64(edakit.Int(1, 10)) / 10
			if sumIn+x >= sumOut {
				x = 0.1
			}
			sumIn += x
			ins = append(ins, edakit.Num(x, 1))
		}
		miss := sumOut - sumIn
		if miss <= 0 {
			miss = 0
		}
		return edakit.Numeric(edakit.Sprintf("In Lagrangian-relaxation sizing, each timing arc carries a multiplier, and at an optimum the multipliers satisfy flow conservation: at every gate, the multipliers on its incoming arcs sum to the multipliers on its outgoing arcs. A gate has outgoing-arc multipliers %s and %d incoming arcs; the known incoming multipliers are %s. What must the remaining incoming multiplier be? Round to 1 decimal place.", strings.Join(outs, ", "), nin, strings.Join(ins, ", ")),
			edakit.Num(miss, 1),
			edakit.Sprintf("Outgoing sum %s minus known incoming sum %s leaves %s. Large multipliers mark critical arcs: they act like the share of criticality flowing through the arc.", edakit.Num(sumOut, 1), edakit.Num(sumIn, 1), edakit.Num(miss, 1)),
			map[string]string{"multiplier": edakit.Num(miss, 1)})
	case 2:
		a := float64(edakit.Int(1, 4))
		mu := float64(edakit.Int(2, 30)) / 10
		r := float64(edakit.Int(10, 30)) / 10
		cl := float64(edakit.Int(10, 60))
		nu := float64(edakit.Int(1, 20)) / 10
		c := float64(edakit.Int(1, 3))
		x := math.Sqrt(mu * r * cl / (a + nu*c))
		return edakit.Numeric(edakit.Sprintf("For fixed multipliers, the Lagrangian separates into one subproblem per gate. Gate i's terms are a x + mu r C_L / x + nu c x, where x is its size, a = %s is its area per unit size, mu = %s is the sum of multipliers on its own arcs, r = %s kOhm, C_L = %s fF, nu = %s is the multiplier-weighted resistance of its driver and c = %s fF is its input capacitance per unit size. What size minimizes the subproblem? Round to 2 decimal places.",
			edakit.Num(a, 0), edakit.Num(mu, 1), edakit.Num(r, 1), edakit.Num(cl, 0), edakit.Num(nu, 1), edakit.Num(c, 0)),
			edakit.Num(x, 2),
			edakit.Sprintf("Setting the derivative to zero gives x = sqrt(mu r C_L / (a + nu c)) = sqrt(%s / %s) = %s. Iterating this local rule over all gates is the inner loop of Lagrangian sizing.", edakit.Num(mu*r*cl, 2), edakit.Num(a+nu*c, 2), edakit.Num(x, 2)),
			map[string]string{"size": edakit.Num(x, 2)})
	}
	lam := float64(edakit.Int(0, 30)) / 10
	theta := float64(edakit.Int(1, 5)) / 100
	arr := float64(edakit.Int(200, 400))
	req := arr + float64(edakit.Int(-60, 60))
	nl := lam + theta*(arr-req)
	if nl < 0 {
		nl = 0
	}
	return edakit.Numeric(edakit.Sprintf("A subgradient step updates an arc multiplier by lambda <- max(0, lambda + theta (arrival - required)). The arc's multiplier is %s, the step size is theta = %s per ps, and the arc's endpoint has arrival %s ps against a required time of %s ps. What is the updated multiplier? Round to 2 decimal places.",
		edakit.Num(lam, 1), edakit.Num(theta, 2), edakit.Num(arr, 0), edakit.Num(req, 0)),
		edakit.Num(nl, 2),
		edakit.Sprintf("lambda + theta (arrival - required) = %s + %s x (%s) = %s; clipping at zero gives %s. Violated arcs gain weight, and arcs with slack lose it.", edakit.Num(lam, 1), edakit.Num(theta, 2), edakit.Num(arr-req, 0), edakit.Num(lam+theta*(arr-req), 2), edakit.Num(nl, 2)),
		map[string]string{"lambda": edakit.Num(nl, 2)})
}

type vtGate struct {
	name string
	dd   int
	save int
}

func vtSwap(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		lk := float64(edakit.Int(20, 200))
		dvt := float64(edakit.Pick([]int{60, 80, 100, 120}))
		n := float64(edakit.Pick([]int{12, 13, 14, 15})) / 10
		vt := 26.0
		f := math.Exp(-dvt / (n * vt))
		ans := lk * f
		return edakit.Numeric(edakit.Sprintf("Subthreshold leakage scales as exp(-V_t / (n v_T)). A cell leaks %s nW in its low-Vt version. Its high-Vt version has V_t higher by %s mV. With n = %s and v_T = 26 mV, how much does the high-Vt version leak, in nW? Round to 2 decimal places.", edakit.Num(lk, 0), edakit.Num(dvt, 0), edakit.Num(n, 1)),
			edakit.Num(ans, 2),
			edakit.Sprintf("The factor is exp(-%s / (%s x 26)) = %s, so the leakage is %s x %s = %s nW.", edakit.Num(dvt, 0), edakit.Num(n, 1), edakit.Num(f, 4), edakit.Num(lk, 0), edakit.Num(f, 4), edakit.Num(ans, 2)),
			map[string]string{"factor": edakit.Num(f, 4), "leakage": edakit.Num(ans, 2)})
	case 2:
		ng := edakit.Scaled(ctx, 3, 5)
		var gs []vtGate
		for i := 0; i < ng; i++ {
			gs = append(gs, vtGate{edakit.Sprintf("G%d", i+1), edakit.Int(3, 12), edakit.Int(5, 40)})
		}
		S := edakit.Int(8, 25)
		tot := 0
		for _, g := range gs {
			tot += g.dd
		}
		if S >= tot {
			S = tot - edakit.Int(2, 6)
		}
		best, bestSet := knap(gs, func(set int) bool {
			t := 0
			for i := range gs {
				if set&(1<<uint(i)) != 0 {
					t += gs[i].dd
				}
			}
			return t <= S
		})
		var desc []string
		for _, g := range gs {
			desc = append(desc, edakit.Sprintf("%s (+%d ps, saves %d nW)", g.name, g.dd, g.save))
		}
		return edakit.Numeric(edakit.Sprintf("Gates on one path, and on no other critical path, can each be swapped from low-Vt to high-Vt: %s. The path has %d ps of slack. What is the largest total leakage saving, in nW, without making the slack negative?", strings.Join(desc, "; "), S),
			edakit.IntStr(best),
			edakit.Sprintf("This is a 0-1 knapsack with the slack as capacity. The best set is %s, saving %d nW. Greedy by saving per ps is a common heuristic but is not guaranteed optimal.", setNames(gs, bestSet), best),
			map[string]string{"saving": edakit.IntStr(best)})
	}
	ng := edakit.Int(4, 5)
	var gs []vtGate
	for i := 0; i < ng; i++ {
		gs = append(gs, vtGate{edakit.Sprintf("G%d", i+1), edakit.Int(3, 10), edakit.Int(5, 40)})
	}
	split := edakit.Int(1, ng-2)
	share := edakit.Int(split, ng-1)
	// path 1: G1..G(share+1); path 2: G(split+1)..Gng
	in1 := func(i int) bool { return i <= share }
	in2 := func(i int) bool { return i >= split }
	S1, S2 := edakit.Int(8, 20), edakit.Int(8, 20)
	best, bestSet := knap(gs, func(set int) bool {
		t1, t2 := 0, 0
		for i := range gs {
			if set&(1<<uint(i)) != 0 {
				if in1(i) {
					t1 += gs[i].dd
				}
				if in2(i) {
					t2 += gs[i].dd
				}
			}
		}
		return t1 <= S1 && t2 <= S2
	})
	var p1, p2, desc []string
	for i, g := range gs {
		desc = append(desc, edakit.Sprintf("%s (+%d ps, saves %d nW)", g.name, g.dd, g.save))
		if in1(i) {
			p1 = append(p1, g.name)
		}
		if in2(i) {
			p2 = append(p2, g.name)
		}
	}
	return edakit.Numeric(edakit.Sprintf("Candidate high-Vt swaps: %s. Path P1 runs through %s and has %d ps of slack; path P2 runs through %s and has %d ps of slack. What is the largest total leakage saving, in nW, that keeps both slacks non-negative?", strings.Join(desc, "; "), strings.Join(p1, ", "), S1, strings.Join(p2, ", "), S2),
		edakit.IntStr(best),
		edakit.Sprintf("A shared gate's delay counts against both paths, so the problem is a knapsack with two capacity constraints. The best set is %s, saving %d nW.", setNames(gs, bestSet), best),
		map[string]string{"saving": edakit.IntStr(best)})
}

func knap(gs []vtGate, ok func(set int) bool) (int, int) {
	best, bestSet := 0, 0
	for set := 0; set < 1<<uint(len(gs)); set++ {
		if !ok(set) {
			continue
		}
		s := 0
		for i := range gs {
			if set&(1<<uint(i)) != 0 {
				s += gs[i].save
			}
		}
		if s > best {
			best, bestSet = s, set
		}
	}
	return best, bestSet
}

func setNames(gs []vtGate, set int) string {
	var s []string
	for i, g := range gs {
		if set&(1<<uint(i)) != 0 {
			s = append(s, g.name)
		}
	}
	if len(s) == 0 {
		return "empty"
	}
	return "{" + strings.Join(s, ", ") + "}"
}

func pinSwap(ctx generator.GeneratorContext) generator.Problem {
	n := 3
	if ctx.Difficulty > 0.6 && edakit.Int(0, 1) == 1 {
		n = 4
	}
	cell := map[int]string{3: "NAND3", 4: "NAND4"}[n]
	var d []int
	base := edakit.Int(10, 14)
	for i := 0; i < n; i++ {
		base += edakit.Int(1, 4)
		d = append(d, base)
	}
	pins := []string{"A", "B", "C", "D"}[:n]
	sigs := []string{"s", "t", "u", "w"}[:n]
	var arr []int
	for i := 0; i < n; i++ {
		arr = append(arr, edakit.Int(0, 40))
	}
	var pd []string
	for i := 0; i < n; i++ {
		pd = append(pd, edakit.Sprintf("%s %d ps", pins[i], d[i]))
	}
	var ad []string
	for i := 0; i < n; i++ {
		ad = append(ad, edakit.Sprintf("%s at %d ps", sigs[i], arr[i]))
	}
	q := edakit.Sprintf("A mapped %s has functionally symmetric inputs with different pin-to-output delays: %s (pin A sits nearest the output in the transistor stack). Its input signals arrive %s.", cell, strings.Join(pd, ", "), strings.Join(ad, ", "))
	// current assignment: signal i on pin i
	cur := 0
	for i := 0; i < n; i++ {
		if arr[i]+d[i] > cur {
			cur = arr[i] + d[i]
		}
	}
	// optimal: latest signal on fastest pin
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool { return arr[idx[a]] > arr[idx[b]] })
	opt := 0
	var asg []string
	for rank, s := range idx {
		if arr[s]+d[rank] > opt {
			opt = arr[s] + d[rank]
		}
		asg = append(asg, sigs[s]+" on "+pins[rank])
	}
	curText := "Currently " + func() string {
		var s []string
		for i := 0; i < n; i++ {
			s = append(s, sigs[i]+" is on pin "+pins[i])
		}
		return strings.Join(s, ", ")
	}() + "."
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(q+" "+curText+" When does the output settle, in ps?", edakit.IntStr(cur),
			edakit.Sprintf("The output settles at the latest of arrival plus pin delay over the inputs: %d ps.", cur),
			map[string]string{"current": edakit.IntStr(cur)})
	case 2:
		return edakit.Numeric(q+" Over all assignments of signals to pins, what is the earliest the output can settle, in ps?", edakit.IntStr(opt),
			"Put the latest-arriving signal on the fastest pin and continue in order: "+strings.Join(asg, ", ")+edakit.Sprintf(". The output settles at %d ps. An exchange argument shows that this sorted pairing is optimal.", opt),
			map[string]string{"optimal": edakit.IntStr(opt)})
	}
	return edakit.Numeric(q+" "+curText+" How many ps does the best pin swap save?", edakit.IntStr(cur-opt),
		edakit.Sprintf("The current assignment settles at %d ps and the sorted assignment (%s) at %d ps, saving %d ps without changing area or function.", cur, strings.Join(asg, ", "), opt, cur-opt),
		map[string]string{"current": edakit.IntStr(cur), "optimal": edakit.IntStr(opt)})
}
