package opt

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func dominates(a, b []int) bool {
	strict := false
	for i := range a {
		if a[i] > b[i] {
			return false
		}
		if a[i] < b[i] {
			strict = true
		}
	}
	return strict
}

func drawDesigns(n, dim int) [][]int {
	pts := make([][]int, n)
	for i := range pts {
		pts[i] = make([]int, dim)
		for d := range pts[i] {
			pts[i][d] = edakit.Int(1, 12)
		}
		// Trade-off structure: in two objectives, faster designs tend to burn more power.
		if dim == 2 && edakit.Int(0, 3) > 0 {
			pts[i][1] = min(12, max(1, 13-pts[i][0]+edakit.Int(-3, 3)))
		}
	}
	return pts
}

func designStr(pts [][]int, labels []string) string {
	var out []string
	for i, p := range pts {
		var c []string
		for d, v := range p {
			c = append(c, edakit.Sprintf("%s %d", labels[d], v))
		}
		out = append(out, edakit.Sprintf("D%d (%s)", i+1, strings.Join(c, ", ")))
	}
	return strings.Join(out, "; ")
}

func paretoCount(pts [][]int) (int, []string) {
	cnt := 0
	var names []string
	for i := range pts {
		dom := false
		for j := range pts {
			if j != i && dominates(pts[j], pts[i]) {
				dom = true
				break
			}
		}
		if !dom {
			cnt++
			names = append(names, edakit.Sprintf("D%d", i+1))
		}
	}
	return cnt, names
}

func pareto(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 5, 9)
		pts := drawDesigns(n, 2)
		k := edakit.Int(0, n-1)
		cnt := 0
		var who []string
		for j := range pts {
			if j != k && dominates(pts[j], pts[k]) {
				cnt++
				who = append(who, edakit.Sprintf("D%d", j+1))
			}
		}
		ws := "none"
		if len(who) > 0 {
			ws = strings.Join(who, ", ")
		}
		return num("A flow sweep produced designs, both objectives to be minimized: "+designStr(pts, []string{"delay", "power"})+edakit.Sprintf(". How many designs dominate D%d (no worse in both objectives and strictly better in at least one)?", k+1)+roundNote(0), edakit.IntStr(cnt),
			edakit.Sprintf("The designs dominating D%d are %s, so the count is %d. Equal values in both objectives do not dominate.", k+1, ws, cnt),
			"count", edakit.IntStr(cnt))
	case 2:
		n := edakit.Scaled(ctx, 6, 10)
		pts := drawDesigns(n, 2)
		cnt, names := paretoCount(pts)
		return num("A flow sweep produced designs, both objectives to be minimized: "+designStr(pts, []string{"delay", "power"})+". How many designs are Pareto-optimal (dominated by no other design)?"+roundNote(0), edakit.IntStr(cnt),
			edakit.Sprintf("Check each design against all others. The nondominated set is %s, %d designs. Every other design can be improved in one objective at no cost in the other.", strings.Join(names, ", "), cnt),
			"pareto", edakit.IntStr(cnt))
	default:
		n := edakit.Scaled(ctx, 5, 8)
		pts := drawDesigns(n, 3)
		cnt, names := paretoCount(pts)
		return num("A PPA sweep produced designs with power, performance cost (delay) and area, all minimized: "+designStr(pts, []string{"power", "delay", "area"})+". How many designs are Pareto-optimal?"+roundNote(0), edakit.IntStr(cnt),
			edakit.Sprintf("With three objectives domination is harder to achieve, so more designs survive. The nondominated set is %s, %d designs.", strings.Join(names, ", "), cnt),
			"pareto", edakit.IntStr(cnt))
	}
}

func scalarization(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 4, 7)
	pts := drawDesigns(n, 2)
	lab := designStr(pts, []string{"delay", "power"})
	switch kp(ctx) {
	case 1:
		w := float64(edakit.Int(1, 9)) / 10
		best := math.Inf(1)
		bi := 0
		for i, p := range pts {
			if v := w*float64(p[0]) + (1-w)*float64(p[1]); v < best-1e-12 {
				best, bi = v, i
			}
		}
		return num("Designs (minimize both): "+lab+edakit.Sprintf(". A flow script scores $s = %s \\cdot \\text{delay} + %s \\cdot \\text{power}$. What is the lowest score?", f(w, 1), f(1-w, 1))+roundNote(2), f(best, 2),
			edakit.Sprintf("Score every design; D%d is lowest at %s. A weighted sum with positive weights always picks a Pareto-optimal design, but only those on the convex hull of the front.", bi+1, f(best, 2)),
			"score", f(best, 2))
	case 2:
		var eps, best int
		for {
			eps = edakit.Int(3, 11)
			best = 1 << 30
			for _, p := range pts {
				if p[1] <= eps && p[0] < best {
					best = p[0]
				}
			}
			if best < 1<<30 {
				break
			}
			pts = drawDesigns(n, 2)
			lab = designStr(pts, []string{"delay", "power"})
		}
		return num("Designs (minimize both): "+lab+edakit.Sprintf(". The epsilon-constraint method minimizes delay subject to power $\\le %d$. What delay does it return?", eps)+roundNote(0), edakit.IntStr(best),
			edakit.Sprintf("Discard designs with power above %d, then take the smallest delay among the rest: %d. Sweeping epsilon traces the whole front, including nonconvex parts a weighted sum misses.", eps, best),
			"delay", edakit.IntStr(best))
	default:
		R1, R2 := 13, 13
		// nondominated front sorted by delay
		var front [][]int
		for i := range pts {
			dom := false
			for j := range pts {
				if j != i && dominates(pts[j], pts[i]) {
					dom = true
				}
			}
			dup := false
			for _, q := range front {
				if q[0] == pts[i][0] && q[1] == pts[i][1] {
					dup = true
				}
			}
			if !dom && !dup {
				front = append(front, pts[i])
			}
		}
		sort.Slice(front, func(a, b int) bool { return front[a][0] < front[b][0] })
		hv := 0
		prev := R2
		for _, p := range front {
			hv += (R1 - p[0]) * (prev - p[1])
			prev = p[1]
		}
		return num("Designs (minimize both): "+lab+edakit.Sprintf(". What is the hypervolume (area dominated by the set and bounded by the reference point (delay %d, power %d))?", R1, R2)+roundNote(0), edakit.IntStr(hv),
			edakit.Sprintf("Only the nondominated designs matter. Sort them by delay and add the rectangle each contributes below the previous one's power: the total is %d. Hypervolume rewards both convergence and spread, which is why multi-objective tuners report it.", hv),
			"hv", edakit.IntStr(hv))
	}
}

func robustStochastic(ctx generator.GeneratorContext) generator.Problem {
	for {
		nd := edakit.Scaled(ctx, 3, 4)
		ns := 3
		corners := []string{"SS", "TT", "FF"}
		cost := make([][]int, nd)
		for i := range cost {
			cost[i] = make([]int, ns)
			for j := range cost[i] {
				cost[i][j] = edakit.Int(2, 20)
			}
		}
		pr := edakit.Pick([][]float64{{0.2, 0.6, 0.2}, {0.1, 0.8, 0.1}, {0.3, 0.5, 0.2}, {0.25, 0.5, 0.25}})
		worst := make([]int, nd)
		exp := make([]float64, nd)
		for i := range cost {
			for j := range cost[i] {
				worst[i] = max(worst[i], cost[i][j])
				exp[i] += pr[j] * float64(cost[i][j])
			}
		}
		rb, rv, eb := 0, worst[0], 0
		ties := 0
		for i := range worst {
			if worst[i] < rv {
				rb, rv, ties = i, worst[i], 0
			} else if worst[i] == rv && i != rb {
				ties++
			}
			if exp[i] < exp[eb] {
				eb = i
			}
		}
		if ties > 0 {
			continue
		}
		var rows []string
		for i := range cost {
			rows = append(rows, edakit.Sprintf("option %d: %s", i+1, ints(cost[i], ", ")))
		}
		stem := "A sizing choice is evaluated at corners " + strings.Join(corners, ", ") + edakit.Sprintf(" with probabilities %s, %s, %s. Costs per corner, in that order: ", f(pr[0], 2), f(pr[1], 2), f(pr[2], 2)) + strings.Join(rows, "; ") + "."
		var es []string
		for i := range exp {
			es = append(es, f(exp[i], 2))
		}
		switch kp(ctx) {
		case 1:
			return num(stem+" A robust (min-max) formulation picks the option with the best worst case. What is that worst-case cost?"+roundNote(0), edakit.IntStr(rv),
				edakit.Sprintf("The worst cases are %s; the smallest is %d (option %d). Robust optimization protects the worst corner and ignores the probabilities.", ints(worst, ", "), rv, rb+1),
				"robust", edakit.IntStr(rv))
		case 2:
			return num(stem+" A stochastic formulation minimizes expected cost. What is the minimum expected cost?"+roundNote(2), f(exp[eb], 2),
				edakit.Sprintf("Expected costs are %s; the minimum is %s (option %d).", strings.Join(es, ", "), f(exp[eb], 2), eb+1),
				"expected", f(exp[eb], 2))
		default:
			pc := exp[rb] - exp[eb]
			return num(stem+" What is the price of robustness: the expected cost of the min-max option minus the minimum expected cost?"+roundNote(2), f(pc, 2),
				edakit.Sprintf("The min-max option %d has expected cost %s; the best expected cost is %s (option %d). The difference %s is what the guarantee on the worst corner costs on average.", rb+1, f(exp[rb], 2), f(exp[eb], 2), eb+1, f(pc, 2)),
				"price", f(pc, 2))
		}
	}
}

func gaussianProcess(ctx generator.GeneratorContext) generator.Problem {
	s2 := float64(edakit.Int(1, 4))
	ell := float64(edakit.Int(1, 4)) / 2
	x1 := float64(edakit.Int(0, 20)) / 10
	x := x1 + float64(edakit.Int(-15, 15))/10
	for x == x1 {
		x = x1 + 0.5
	}
	y1 := float64(edakit.Int(-30, 30)) / 10
	noise := 0.0
	if edakit.Int(0, 1) == 0 {
		noise = float64(edakit.Int(1, 5)) / 10
	}
	k := s2 * math.Exp(-(x-x1)*(x-x1)/(2*ell*ell))
	stem := edakit.Sprintf("A Bayesian flow tuner models QoR versus a utilization knob with a zero-mean Gaussian process and squared-exponential kernel $k(x, x') = %s\\exp(-(x - x')^2 / (2 \\cdot %s^2))$.", f(s2, 0), f(ell, 1))
	nstr := "noise-free"
	if noise > 0 {
		nstr = edakit.Sprintf("with noise variance %s", f(noise, 1))
	}
	switch kp(ctx) {
	case 1:
		return num(stem+edakit.Sprintf(" What is the covariance $k(%s, %s)$?", f(x, 1), f(x1, 1))+roundNote(3), f(k, 3),
			edakit.Sprintf("The distance is %s, so $k = %s \\exp(-%s^2/(2 \\cdot %s^2)) = %s$. Nearby inputs are strongly correlated; at a few length scales the correlation vanishes.", f(math.Abs(x-x1), 1), f(s2, 0), f(math.Abs(x-x1), 1), f(ell, 1), f(k, 3)),
			"k", f(k, 3))
	case 2:
		mu := k / (s2 + noise) * y1
		return num(stem+edakit.Sprintf(" One run, %s, measured $y = %s$ at $x_1 = %s$. What is the posterior mean at $x = %s$?", nstr, f(y1, 1), f(x1, 1), f(x, 1))+roundNote(3), f(mu, 3),
			edakit.Sprintf("$\\mu(x) = k(x, x_1)\\,(k(x_1, x_1) + \\sigma_n^2)^{-1} y_1 = %s / %s \\cdot %s = %s$. Far from data the mean falls back to the prior mean 0.", f(k, 3), f(s2+noise, 1), f(y1, 1), f(mu, 3)),
			"k", f(k, 4), "mu", f(mu, 3))
	default:
		v := s2 - k*k/(s2+noise)
		return num(stem+edakit.Sprintf(" One run, %s, was made at $x_1 = %s$. What is the posterior variance at $x = %s$?", nstr, f(x1, 1), f(x, 1))+roundNote(3), f(v, 3),
			edakit.Sprintf("$\\sigma^2(x) = k(x, x) - k(x, x_1)^2/(k(x_1, x_1) + \\sigma_n^2) = %s - %s^2/%s = %s$. The variance does not depend on the measured value, only on where the data are.", f(s2, 0), f(k, 4), f(s2+noise, 1), f(v, 3)),
			"k", f(k, 4), "var", f(v, 3))
	}
}

func bayesOpt(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Int(3, 5)
		kap := float64(edakit.Int(1, 6)) / 2
		best, bi := math.Inf(-1), 0
		var parts []string
		for i := 0; i < n; i++ {
			mu := float64(edakit.Int(0, 50)) / 10
			sd := float64(edakit.Int(1, 20)) / 10
			if u := mu + kap*sd; u > best+1e-12 {
				best, bi = u, i
			}
			parts = append(parts, edakit.Sprintf("x%d: mean %s, sd %s", i+1, f(mu, 1), f(sd, 1)))
		}
		return num("An analog-sizing BO loop maximizes gain margin. The surrogate predicts at candidate points "+strings.Join(parts, "; ")+edakit.Sprintf(". With the upper confidence bound $\\mu + \\kappa\\sigma$ and $\\kappa = %s$, what acquisition value does the chosen point have?", f(kap, 1))+roundNote(2), f(best, 2),
			edakit.Sprintf("Compute $\\mu + %s\\sigma$ for each and take the largest: x%d with %s. Larger $\\kappa$ favours uncertain points (exploration).", f(kap, 1), bi+1, f(best, 2)),
			"ucb", f(best, 2))
	case 2:
		fs := float64(edakit.Int(20, 40)) / 10
		mu := fs + float64(edakit.Int(-15, 15))/10
		sd := float64(edakit.Int(2, 15)) / 10
		z := (fs - mu) / sd
		ei := (fs-mu)*bigPhi(z) + sd*phi(z)
		return num(edakit.Sprintf("A BO tuner minimizes a flow's power. The best power so far is $f^* = %s$ mW. At a candidate the GP predicts mean %s and standard deviation %s. What is the expected improvement $EI = (f^* - \\mu)\\Phi(z) + \\sigma\\varphi(z)$ with $z = (f^* - \\mu)/\\sigma$?", f(fs, 1), f(mu, 1), f(sd, 1))+roundNote(3), f(ei, 3),
			edakit.Sprintf("$z = %s$, $\\Phi(z) = %s$, $\\varphi(z) = %s$, so $EI = %s \\cdot %s + %s \\cdot %s = %s$ mW. EI is positive even when the mean is worse than $f^*$, because uncertainty might hide an improvement.", f(z, 3), f(bigPhi(z), 4), f(phi(z), 4), f(fs-mu, 1), f(bigPhi(z), 4), f(sd, 1), f(phi(z), 4), f(ei, 3)),
			"z", f(z, 3), "ei", f(ei, 3))
	default:
		fs := float64(edakit.Int(10, 40)) / 10
		mu1 := fs - float64(edakit.Int(1, 5))/10
		sd1 := float64(edakit.Int(1, 3)) / 10
		mu2 := fs + float64(edakit.Int(0, 8))/10
		sd2 := float64(edakit.Int(8, 20)) / 10
		ei := func(m, s float64) float64 {
			z := (fs - m) / s
			return (fs-m)*bigPhi(z) + s*phi(z)
		}
		e1, e2 := ei(mu1, sd1), ei(mu2, sd2)
		return num(edakit.Sprintf("A BO tuner minimizes timing-closure runtime with incumbent $f^* = %s$ h. Candidate P has mean %s and sd %s (a safe exploit); candidate Q has mean %s and sd %s (an uncertain explore). EI is $(f^* - \\mu)\\Phi(z) + \\sigma\\varphi(z)$ with $z = (f^* - \\mu)/\\sigma$. What is the larger of the two EI values?", f(fs, 1), f(mu1, 1), f(sd1, 1), f(mu2, 1), f(sd2, 1))+roundNote(3), f(math.Max(e1, e2), 3),
			edakit.Sprintf("EI(P) = %s and EI(Q) = %s, so the tuner runs %s next. A wide predictive spread can outweigh a worse mean.", f(e1, 3), f(e2, 3), map[bool]string{true: "P", false: "Q"}[e1 >= e2]),
			"ei_p", f(e1, 3), "ei_q", f(e2, 3))
	}
}
