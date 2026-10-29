package analog

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// feasibleChoice draws four candidate op-amp sizings and asks which one a constrained
// formulation selects: minimum power subject to gain, GBW and phase-margin specs.
func feasibleChoice(ctx generator.GeneratorContext) generator.Problem {
	gainS := float64(edakit.Int(60, 75))
	gbwS := float64(10 * edakit.Int(2, 10))
	pmS := 60.0
	type cand struct{ gain, gbw, pm, pw float64 }
	for attempt := 0; ; attempt++ {
		var cs []cand
		for i := 0; i < 4; i++ {
			cs = append(cs, cand{gainS + unif(-6, 8, 1), gbwS * unif(0.7, 1.5, 0.05), unif(45, 80, 1), unif(0.2, 2, 0.05)})
		}
		if attempt >= generator.MaxDrawAttempts {
			// Known-good instance: D3 is the cheapest feasible design, D2 the cheapest overall.
			gainS, gbwS = 65, 50
			cs = []cand{{70, 60, 65, 0.9}, {62, 55, 70, 0.5}, {68, 52, 62, 0.7}, {66, 48, 75, 0.6}}
		}
		best, nf := -1, 0
		minAll := 0
		for i, c := range cs {
			if c.pw < cs[minAll].pw {
				minAll = i
			}
			if c.gain >= gainS && c.gbw >= gbwS && c.pm >= pmS {
				nf++
				if best < 0 || c.pw < cs[best].pw {
					best = i
				}
			}
		}
		// Require a feasible winner that is not simply the lowest-power candidate, and no ties.
		if best < 0 || best == minAll || nf == 4 {
			continue
		}
		tie := false
		for i := range cs {
			for j := i + 1; j < 4; j++ {
				if cs[i].pw == cs[j].pw {
					tie = true
				}
			}
		}
		if tie {
			continue
		}
		var rows []string
		for i, c := range cs {
			rows = append(rows, edakit.Sprintf("D%d: gain %s dB, GBW %s MHz, PM %s degrees, power %s mW", i+1, f(c.gain, 0), f(c.gbw, 1), f(c.pm, 0), f(c.pw, 2)))
		}
		q := edakit.Sprintf("A sizing run is posed as: minimize power subject to gain at least %s dB, GBW at least %s MHz and phase margin at least %s degrees. Simulation of four candidate sizings gives: ", f(gainS, 0), f(gbwS, 0), f(pmS, 0)) + strings.Join(rows, "; ") + ". Which candidate is the optimum of this formulation?"
		var ds []string
		for i := range cs {
			if i != best {
				ds = append(ds, edakit.Sprintf("D%d", i+1))
			}
		}
		return edakit.Choice(q, edakit.Sprintf("D%d", best+1),
			edakit.Sprintf("Discard every candidate that violates a constraint, then take the lowest power among the rest: D%d. D%d has the lowest power overall but fails a spec.", best+1, minAll+1), ds...)
	}
}

func equationBased(ctx generator.GeneratorContext) generator.Problem {
	gbw := float64(edakit.Int(5, 100))
	cc := unif(0.5, 4, 0.1)
	gm1 := 2 * math.Pi * gbw * cc
	switch whichKP(ctx) {
	case 1:
		q := edakit.Sprintf("A hand-sizing script for a Miller two-stage op-amp starts from the spec GBW = %s MHz with $C_c = %s$ pF. What input-pair transconductance $g_{m1}$ does $\\mathrm{GBW} = g_{m1}/(2\\pi C_c)$ require, in microsiemens? Round to 1 decimal place.", f(gbw, 0), f(cc, 1))
		return num(q, f(gm1, 1), edakit.Sprintf("$g_{m1} = 2\\pi \\cdot \\mathrm{GBW} \\cdot C_c = 2\\pi \\cdot %s \\cdot %s = %s\\ \\mu$S (MHz times pF gives microsiemens).", f(gbw, 0), f(cc, 1), f(gm1, 1)), nil)
	case 2:
		gmid := float64(edakit.Int(8, 22))
		tail := 2 * gm1 / gmid
		q := edakit.Sprintf("An equation-based sizing step needs $g_{m1} = 2\\pi \\cdot \\mathrm{GBW} \\cdot C_c$ for GBW = %s MHz and $C_c = %s$ pF, and the input pair is biased at $g_m/I_D = %s\\ \\mathrm{V}^{-1}$. What tail current does the first stage need, in microamperes? Round to 1 decimal place.", f(gbw, 0), f(cc, 1), f(gmid, 0))
		return num(q, f(tail, 1), edakit.Sprintf("$g_{m1} = %s\\ \\mu$S, so each side carries $I_D = %s/%s = %s\\ \\mu$A, and the tail is twice that, %s microamperes.", f(gm1, 2), f(gm1, 2), f(gmid, 0), f(gm1/gmid, 2), f(tail, 1)), nil)
	default:
		cl := unif(1, 10, 0.5)
		k := pick(2.2, 2.5, 3.0)
		gm6 := k * 2 * math.Pi * gbw * cl
		q := edakit.Sprintf("For phase margin, a sizing procedure places the non-dominant pole $g_{m6}/(2\\pi C_L)$ at %s times GBW. With GBW = %s MHz and $C_L = %s$ pF, what minimum $g_{m6}$ is needed, in microsiemens? Round to 1 decimal place.", f(k, 1), f(gbw, 0), f(cl, 1))
		return num(q, f(gm6, 1), edakit.Sprintf("$g_{m6} = %s \\cdot 2\\pi \\cdot %s \\cdot %s = %s\\ \\mu$S. With a dominant pole at $-90$ degrees this gives about $90 - \\arctan(1/%s) = %s$ degrees of margin.", f(k, 1), f(gbw, 0), f(cl, 1), f(gm6, 1), f(k, 1), f(90-deg(math.Atan(1/k)), 1)), nil)
	}
}

func simLoop(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		it := 10 * edakit.Int(5, 60)
		pop := edakit.Int(4, 32)
		corners := pick(5, 9, 15, 27, 45)
		t := float64(edakit.Int(5, 120))
		lic := pick(4, 8, 16, 32, 64)
		hours := float64(it*pop*corners) * t / float64(lic) / 3600
		q := edakit.Sprintf("A population-based optimizer runs %d generations of %d candidate sizings. Each candidate is simulated at %d PVT corners, each simulation takes %s seconds, and %d simulator licenses run in parallel with perfect packing. How many hours of wall-clock time does the run take? Round to 1 decimal place.", it, pop, corners, f(t, 0), lic)
		return num(q, f(hours, 1), edakit.Sprintf("$%d \\cdot %d \\cdot %d = %d$ simulations, times %s s, divided by %d licenses and 3600 s per hour: %s hours.", it, pop, corners, it*pop*corners, f(t, 0), lic, f(hours, 1)), nil)
	case 2:
		n := edakit.Int(4, 30)
		corners := pick(1, 3, 5, 9)
		central := edakit.Int(0, 1) == 0
		per := n + 1
		kind := "forward"
		if central {
			per = 2 * n
			kind = "central"
		}
		tot := per * corners
		q := edakit.Sprintf("A gradient-based sizer estimates the gradient of a performance with respect to %d design variables by %s finite differences, and does so at each of %d corners. How many SPICE simulations does one gradient evaluation cost?", n, kind, corners)
		return num(q, edakit.IntStr(tot), edakit.Sprintf("%s differences need %d simulations per corner (%s), times %d corners: %d.", strings.ToUpper(kind[:1])+kind[1:], per, map[bool]string{true: "two per variable", false: "one per variable plus the nominal point"}[central], corners, tot), nil)
	default:
		n := 100 * edakit.Int(2, 50)
		pass := edakit.Int(n/4, n-n/20)
		ylo := float64(pass) / float64(n)
		se := math.Sqrt(ylo*(1-ylo)/float64(n)) * 100
		q := edakit.Sprintf("A simulation-in-the-loop sizer estimates yield by Monte Carlo: %d of %d mismatch samples pass every spec. What is the standard error of the yield estimate, in percentage points? Round to 2 decimal places.", pass, n)
		return num(q, f(se, 2), edakit.Sprintf("$\\hat{Y} = %s$, and the binomial standard error is $\\sqrt{\\hat{Y}(1-\\hat{Y})/n} = %s$ percentage points. Halving it needs four times the samples.", f(ylo, 4), f(se, 2)), nil)
	}
}

func bayesOpt(ctx generator.GeneratorContext) generator.Problem {
	best := unif(0.5, 3, 0.01)
	switch whichKP(ctx) {
	case 1:
		mu := best + unif(-0.4, 0.4, 0.01)
		sg := unif(0.05, 0.5, 0.01)
		z := (best - mu) / sg
		p := phi(z)
		q := edakit.Sprintf("Bayesian optimization minimizes an op-amp's power. The best feasible design so far uses %s mW. At a candidate sizing, the Gaussian-process surrogate predicts mean %s mW and standard deviation %s mW. What is the probability of improvement over the best design? Round to 3 decimal places.", f(best, 2), f(mu, 2), f(sg, 2))
		return num(q, f(p, 3), edakit.Sprintf("$\\mathrm{PI} = \\Phi((f^* - \\mu)/\\sigma) = \\Phi(%s) = %s$.", f(z, 4), f(p, 3)), map[string]string{"z": f(z, 4)})
	case 2:
		mu := best + unif(-0.3, 0.3, 0.01)
		sg := unif(0.05, 0.5, 0.01)
		z := (best - mu) / sg
		ei := (best-mu)*phi(z) + sg*pdf(z)
		q := edakit.Sprintf("A Bayesian sizer minimizes power; the incumbent is %s mW. The surrogate at a candidate predicts mean %s mW and standard deviation %s mW. What is the expected improvement $\\mathrm{EI} = (f^*-\\mu)\\Phi(z) + \\sigma\\varphi(z)$, $z = (f^*-\\mu)/\\sigma$, in mW? Round to 4 decimal places.", f(best, 2), f(mu, 2), f(sg, 2))
		return num(q, f(ei, 4), edakit.Sprintf("$z = %s$, $\\Phi(z) = %s$, $\\varphi(z) = %s$, so $\\mathrm{EI} = %s \\cdot %s + %s \\cdot %s = %s$ mW.", f(z, 4), f(phi(z), 4), f(pdf(z), 4), f(best-mu, 2), f(phi(z), 4), f(sg, 2), f(pdf(z), 4), f(ei, 4)), nil)
	default:
		kappa := pick(1.0, 1.5, 2.0, 2.5, 3.0)
		fbMu := []float64{1.5, 1.2, 1.6, 2.0}
		fbSg := []float64{0.05, 0.02, 0.3, 0.1}
		for attempt := 0; ; attempt++ {
			if attempt >= generator.MaxDrawAttempts {
				kappa = 2
			}
			var rows []string
			bi, bv := -1, 0.0
			mi := 0
			mus := make([]float64, 4)
			vals := make([]float64, 4)
			for i := 0; i < 4; i++ {
				mu := unif(1, 3, 0.01)
				sg := unif(0.02, 0.6, 0.01)
				if attempt >= generator.MaxDrawAttempts {
					mu, sg = fbMu[i], fbSg[i]
				}
				mus[i] = mu
				vals[i] = mu - kappa*sg
				rows = append(rows, edakit.Sprintf("C%d: mean %s, sigma %s", i+1, f(mu, 2), f(sg, 2)))
				if bi < 0 || vals[i] < bv {
					bi, bv = i, vals[i]
				}
				if mus[i] < mus[mi] {
					mi = i
				}
			}
			ok := bi != mi
			for i := 0; i < 4; i++ {
				if i != bi && math.Abs(vals[i]-bv) < 0.02 {
					ok = false
				}
			}
			if !ok {
				continue
			}
			q := edakit.Sprintf("A sizer picks its next SPICE run by the lower confidence bound $\\mu - \\kappa\\sigma$ with $\\kappa = %s$, minimizing power in mW. The surrogate gives: %s. Which candidate number does it simulate next?", f(kappa, 1), strings.Join(rows, "; "))
			return num(q, edakit.IntStr(bi+1), edakit.Sprintf("The bounds are %s, %s, %s and %s; the smallest is C%d. Pure exploitation would pick C%d, the lowest mean.", f(vals[0], 3), f(vals[1], 3), f(vals[2], 3), f(vals[3], 3), bi+1, mi+1), nil)
		}
	}
}

func pvtRobust(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		p := pick(3, 5)
		v := pick(2, 3)
		t := pick(2, 3, 4)
		mc := pick(0, 50, 100, 200)
		tot := p * v * t
		q := edakit.Sprintf("A verification plan sweeps %d process corners, %d supply voltages and %d temperatures as a full grid", p, v, t)
		if mc > 0 {
			q += edakit.Sprintf(", and at every PVT point runs %d mismatch Monte Carlo samples", mc)
			tot *= mc
		}
		q += ". How many simulations does one candidate sizing cost?"
		return num(q, edakit.IntStr(tot), edakit.Sprintf("A full grid multiplies the axis sizes: $%d \\cdot %d \\cdot %d = %d$ PVT points%s.", p, v, t, p*v*t, map[bool]string{true: edakit.Sprintf(", times %d samples is %d", mc, tot), false: ""}[mc > 0]), nil)
	case 2:
		names := []string{"TT 27C 1.8V", "FF -40C 1.98V", "SS 125C 1.62V", "FS 27C 1.62V", "SF 125C 1.98V", "SS -40C 1.62V"}
		spec := float64(edakit.Int(40, 70))
		var rows []string
		worst := 1e9
		for _, nm := range names {
			v := spec + unif(-6, 18, 0.1)
			if v < worst {
				worst = v
			}
			rows = append(rows, edakit.Sprintf("%s: %s", nm, f(v, 1)))
		}
		m := worst - spec
		q := edakit.Sprintf("An op-amp's phase margin is simulated at six corners (degrees): %s. The spec is PM at least %s degrees across all corners. What is the worst-case margin to the spec, in degrees (negative if it fails)? Round to 1 decimal place.", strings.Join(rows, "; "), f(spec, 0))
		return num(q, f(m, 1), edakit.Sprintf("Robust sizing optimizes the worst corner, %s degrees, so the margin is $%s - %s = %s$ degrees.", f(worst, 1), f(worst, 1), f(spec, 0), f(m, 1)), nil)
	default:
		mu := unif(50, 80, 0.5)
		sg := unif(0.5, 5, 0.1)
		spec := edakit.Round(mu-sg*unif(0.5, 3.5, 0.05), 2)
		y := phi((mu-spec)/sg) * 100
		q := edakit.Sprintf("Monte Carlo over process and mismatch gives an op-amp DC gain that is approximately normal with mean %s dB and standard deviation %s dB. The spec is gain at least %s dB. What is the parametric yield, in percent? Round to 2 decimal places.", f(mu, 1), f(sg, 1), f(spec, 2))
		return num(q, f(y, 2), edakit.Sprintf("$Y = \\Phi((\\mu - \\mathrm{spec})/\\sigma) = \\Phi(%s) = %s$%%. Moving the mean one sigma further from the spec is what robust sizing buys.", f((mu-spec)/sg, 4), f(y, 2)), nil)
	}
}
