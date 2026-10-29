package mleda

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// speedupAccounting: speedup of surrogate screening plus verification (KP 1), with the
// labeling and training cost amortized (KP 2), and the whole-flow Amdahl bound (KP 3).
func speedupAccounting(ctx generator.GeneratorContext) generator.Problem {
	k := kp(ctx)
	switch k {
	case 1, 2:
		n := edakit.Scaled(ctx, 100, 5000)
		tt := float64(edakit.Int(20, 240))
		ts := float64(edakit.Int(1, 30)) / 10
		fr := edakit.Pick([]float64{0.01, 0.02, 0.05, 0.1, 0.2})
		base := float64(n) * tt
		with := float64(n)*ts + fr*float64(n)*tt
		stem := edakit.Sprintf("An exploration loop must score %d candidate sizings. One SPICE-based evaluation takes %s minutes. A surrogate scores a candidate in %s minutes, and the top %s percent of candidates by surrogate score are re-evaluated with SPICE before anything is accepted.", n, f(tt, 0), f(ts, 1), f(100*fr, 0))
		if k == 1 {
			s := base / with
			ans := edakit.Num(s, 2)
			return edakit.Numeric(stem+" What is the speedup over evaluating every candidate with SPICE? Round to 2 decimal places.",
				ans, edakit.Sprintf("SPICE only: $%d \\times %s = %s$ min. With the surrogate: $%d \\times %s + %s \\times %d \\times %s = %s$ min. Speedup $= %s$. The verification share caps it near $1/%s$.", n, f(tt, 0), f(base, 0), n, f(ts, 1), f(fr, 2), n, f(tt, 0), f(with, 1), ans, f(fr, 2)),
				map[string]string{"with": f(with, 1)})
		}
		c := float64(edakit.Int(50, 400)) * 60
		m := edakit.Int(2, 20)
		tot := with + c/float64(m)
		s := base / tot
		ans := edakit.Num(s, 2)
		return edakit.Numeric(stem+edakit.Sprintf(" Building the surrogate cost %s machine-hours of labels and training, shared equally across %d such exploration campaigns. Including that share, what is the speedup per campaign? Round to 2 decimal places.", f(c/60, 0), m),
			ans, edakit.Sprintf("Per campaign: screening and verification %s min plus $%s \\times 60 / %d = %s$ min of amortized build cost, total %s min. Speedup $= %s / %s = %s$.", f(with, 1), f(c/60, 0), m, f(c/float64(m), 1), f(tot, 1), f(base, 0), f(tot, 1), ans),
			map[string]string{"amortized": f(c/float64(m), 1)})
	default:
		p := float64(edakit.Int(10, 80)) / 100
		s := float64(edakit.Pick([]int{5, 10, 50, 100, 1000}))
		o := 1 / ((1 - p) + p/s)
		ans := edakit.Num(o, 3)
		return edakit.Numeric(edakit.Sprintf("Detailed routing plus DRC takes %s percent of a block's flow runtime. A learned model replaces the per-iteration DRC estimate inside this step and makes the step %s times faster; the rest of the flow is unchanged. What is the overall flow speedup? Round to 3 decimal places.", f(100*p, 0), f(s, 0)),
			ans, edakit.Sprintf("Amdahl: $1/((1 - %s) + %s/%s) = %s$. Even an infinite step speedup cannot exceed $1/(1-%s) = %s$.", f(p, 2), f(p, 2), f(s, 0), ans, f(p, 2), f(1/(1-p), 3)),
			map[string]string{"bound": f(1/(1-p), 3)})
	}
}

// activeLearning: ensemble variance (KP 1), uncertainty sampling choice (KP 2), and net hours
// saved by active learning (KP 3).
func activeLearning(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		m := edakit.Int(3, 5)
		base := float64(edakit.Int(40, 90))
		p := make([]float64, m)
		mean := 0.0
		for i := range p {
			p[i] = base + float64(edakit.Int(-8, 8))
			mean += p[i]
		}
		mean /= float64(m)
		v := 0.0
		for _, x := range p {
			v += (x - mean) * (x - mean)
		}
		v /= float64(m)
		ans := edakit.Num(v, 3)
		return edakit.Numeric(edakit.Sprintf("An ensemble of %d surrogates predicts the phase margin of one candidate op-amp sizing as %s degrees. Using the population variance (divide by the number of models), what is the ensemble's disagreement, in square degrees? Round to 3 decimal places.", m, floats(p, 0)),
			ans, edakit.Sprintf("Mean $= %s$. Variance $= \\frac{1}{%d}\\sum (p_i - %s)^2 = %s$.", f(mean, 4), m, f(mean, 4), ans),
			map[string]string{"mean": f(mean, 4)})
	case 2:
		var rows []string
		best, bestV := -1, -1.0
		for tries := 0; tries < 50; tries++ {
			rows = nil
			best, bestV = -1, -1.0
			second := -1.0
			for i := 0; i < 5; i++ {
				c := float64(edakit.Int(-20, 40))
				s := []float64{c + float64(edakit.Int(-6, 6)), c + float64(edakit.Int(-6, 6)), c + float64(edakit.Int(-6, 6))}
				mu := (s[0] + s[1] + s[2]) / 3
				v := ((s[0]-mu)*(s[0]-mu) + (s[1]-mu)*(s[1]-mu) + (s[2]-mu)*(s[2]-mu)) / 3
				rows = append(rows, edakit.Sprintf("candidate %d: %s", i+1, floats(s, 0)))
				if v > bestV {
					second = bestV
					best, bestV = i, v
				} else if v > second {
					second = v
				}
			}
			if bestV-second > 0.5 {
				break
			}
		}
		return edakit.Numeric(edakit.Sprintf("An active learner for slack prediction has five unlabeled candidate paths. Three ensemble members predict their slacks (ps) as follows: %s. With uncertainty sampling, which candidate is sent to the timer next? Give its number.", strings.Join(rows, "; ")),
			edakit.IntStr(best+1), edakit.Sprintf("Uncertainty sampling labels the point of greatest ensemble variance. Candidate %d has variance %s, the largest; its predicted slack value does not matter for this rule.", best+1, f(bestV, 3)),
			map[string]string{"variance": f(bestV, 3)})
	default:
		np := edakit.Scaled(ctx, 500, 5000)
		na := np * edakit.Int(15, 60) / 100
		t := float64(edakit.Int(1, 8))
		b := edakit.Pick([]int{10, 20, 50, 100})
		r := float64(edakit.Int(1, 10)) / 2
		rounds := int(math.Ceil(float64(na) / float64(b)))
		saved := float64(np-na)*t - float64(rounds)*r
		ans := edakit.Num(saved, 1)
		return edakit.Numeric(edakit.Sprintf("Random sampling needs %d labeled runs to reach a target error; active learning reaches it with %d. Each label is a %s-hour tool run. Active learning labels in batches of %d and retrains after every batch, at %s hours per retrain. How many machine-hours does active learning save? Round to 1 decimal place.", np, na, f(t, 0), b, f(r, 1)),
			ans, edakit.Sprintf("Labels saved: $(%d - %d) \\times %s = %s$ h. Retrains: $\\lceil %d/%d \\rceil = %d$ rounds at %s h, %s h. Net: %s h.", np, na, f(t, 0), f(float64(np-na)*t, 0), na, b, rounds, f(r, 1), f(float64(rounds)*r, 1), ans),
			map[string]string{"rounds": edakit.IntStr(rounds)})
	}
}

// se is the squared-exponential kernel.
func se(sf2, ell, a, b float64) float64 {
	return sf2 * math.Exp(-(a-b)*(a-b)/(2*ell*ell))
}

// gpPosterior: posterior mean (KP 1) and variance (KP 2) after one observation, and the
// posterior mean after two (KP 3).
func gpPosterior(ctx generator.GeneratorContext) generator.Problem {
	sf2 := float64(edakit.Int(1, 4))
	ell := float64(edakit.Int(1, 4))
	sn2 := edakit.Pick([]float64{0, 0.01, 0.1, 0.25})
	kern := edakit.Sprintf("$k(x, x') = %s\\exp(-(x - x')^2 / (2 \\cdot %s^2))$, noise variance $\\sigma_n^2 = %s$", f(sf2, 0), f(ell, 0), f(sn2, 2))
	k := kp(ctx)
	switch k {
	case 1, 2:
		m := float64(edakit.Int(-5, 5))
		x1 := float64(edakit.Int(0, 8))
		xs := x1 + float64(edakit.Pick([]int{-3, -2, -1, 1, 2, 3}))
		y1 := m + float64(edakit.Int(-10, 10))/2
		ks := se(sf2, ell, xs, x1)
		den := sf2 + sn2
		stem := edakit.Sprintf("A Gaussian-process surrogate of a flow's power, in mW, has constant prior mean %s and kernel %s. One run at $x_1 = %s$ returned $y_1 = %s$.", f(m, 0), kern, f(x1, 0), f(y1, 1))
		if k == 1 {
			mu := m + ks/den*(y1-m)
			ans := edakit.Num(mu, 4)
			return edakit.Numeric(stem+edakit.Sprintf(" What is the posterior mean at $x^* = %s$? Round to 4 decimal places.", f(xs, 0)),
				ans, edakit.Sprintf("$k(x^*, x_1) = %s$. $\\mu^* = m + \\frac{k(x^*,x_1)}{k(x_1,x_1) + \\sigma_n^2}(y_1 - m) = %s + \\frac{%s}{%s}(%s) = %s$. Far from the data the mean returns to the prior.", f(ks, 6), f(m, 0), f(ks, 6), f(den, 2), f(y1-m, 1), ans),
				map[string]string{"k_star": f(ks, 6)})
		}
		v := sf2 - ks*ks/den
		ans := edakit.Num(v, 4)
		return edakit.Numeric(stem+edakit.Sprintf(" What is the posterior variance at $x^* = %s$? Round to 4 decimal places.", f(xs, 0)),
			ans, edakit.Sprintf("$k(x^*, x_1) = %s$. $\\sigma^{*2} = k(x^*,x^*) - \\frac{k(x^*,x_1)^2}{k(x_1,x_1) + \\sigma_n^2} = %s - %s/%s = %s$. The variance does not depend on $y_1$.", f(ks, 6), f(sf2, 0), f(ks*ks, 6), f(den, 2), ans),
			map[string]string{"k_star": f(ks, 6)})
	default:
		x1 := float64(edakit.Int(0, 4))
		x2 := x1 + float64(edakit.Int(1, 3))
		xs := x1 + float64(edakit.Int(-1, 4))
		y1 := float64(edakit.Int(-6, 6))
		y2 := float64(edakit.Int(-6, 6))
		a := sf2 + sn2
		b := se(sf2, ell, x1, x2)
		det := a*a - b*b
		k1, k2 := se(sf2, ell, xs, x1), se(sf2, ell, xs, x2)
		// alpha = (K + sn2 I)^{-1} y
		al1 := (a*y1 - b*y2) / det
		al2 := (-b*y1 + a*y2) / det
		mu := k1*al1 + k2*al2
		ans := edakit.Num(mu, 4)
		return edakit.Numeric(edakit.Sprintf("A zero-mean Gaussian-process surrogate has kernel %s. Two runs gave $y = %s$ at $x = %s$ and $y = %s$ at $x = %s$. What is the posterior mean at $x^* = %s$? Round to 4 decimal places.", kern, f(y1, 0), f(x1, 0), f(y2, 0), f(x2, 0), f(xs, 0)),
			ans, edakit.Sprintf("$K + \\sigma_n^2 I = \\begin{pmatrix} %s & %s \\\\ %s & %s \\end{pmatrix}$, so $\\alpha = (K + \\sigma_n^2 I)^{-1} y = (%s, %s)$. With $k_* = (%s, %s)$, $\\mu^* = k_*^\\top \\alpha = %s$.", f(a, 2), f(b, 4), f(b, 4), f(a, 2), f(al1, 4), f(al2, 4), f(k1, 4), f(k2, 4), ans),
			map[string]string{"alpha1": f(al1, 4), "alpha2": f(al2, 4)})
	}
}

// ei is expected improvement for minimization below incumbent fstar.
func ei(fstar, mu, s float64) float64 {
	z := (fstar - mu) / s
	return (fstar-mu)*cdf(z) + s*phi(z)
}

// acquisition: EI of one candidate (KP 1), the LCB choice (KP 2), and the EI choice among
// several candidates (KP 3).
func acquisition(ctx generator.GeneratorContext) generator.Problem {
	fstar := float64(edakit.Int(40, 80))
	switch kp(ctx) {
	case 1:
		mu := fstar + float64(edakit.Int(-10, 10))
		s := float64(edakit.Int(1, 10))
		z := (fstar - mu) / s
		v := ei(fstar, mu, s)
		ans := edakit.Num(v, 4)
		return edakit.Numeric(edakit.Sprintf("Bayesian optimization minimizes a block's total power. The best power observed so far is $f^* = %s$ mW. At a candidate setting, the GP posterior is Gaussian with mean %s mW and standard deviation %s mW. What is the expected improvement $\\mathrm{EI} = (f^* - \\mu)\\Phi(z) + \\sigma\\varphi(z)$ with $z = (f^* - \\mu)/\\sigma$, in mW? Round to 4 decimal places.", f(fstar, 0), f(mu, 0), f(s, 0)),
			ans, edakit.Sprintf("$z = (%s - %s)/%s = %s$. $\\Phi(z) = %s$, $\\varphi(z) = %s$. EI $= %s \\cdot %s + %s \\cdot %s = %s$ mW.", f(fstar, 0), f(mu, 0), f(s, 0), f(z, 4), f(cdf(z), 4), f(phi(z), 4), pn(fstar-mu, 0), f(cdf(z), 4), f(s, 0), f(phi(z), 4), ans),
			map[string]string{"z": f(z, 4)})
	case 2:
		kappa := edakit.Pick([]float64{1, 2, 3})
		var rows []string
		best, bestV := 0, math.Inf(1)
		for tries := 0; tries < 50; tries++ {
			rows = nil
			best, bestV = 0, math.Inf(1)
			vals := []float64{}
			for i := 0; i < 4; i++ {
				mu := fstar + float64(edakit.Int(-5, 15))
				s := float64(edakit.Int(1, 9))
				v := mu - kappa*s
				vals = append(vals, v)
				rows = append(rows, edakit.Sprintf("%d: mean %s, sd %s", i+1, f(mu, 0), f(s, 0)))
				if v < bestV {
					best, bestV = i, v
				}
			}
			ties := 0
			for _, v := range vals {
				if v == bestV {
					ties++
				}
			}
			if ties == 1 {
				break
			}
		}
		return edakit.Numeric(edakit.Sprintf("A flow-tuning loop minimizes worst negative slack magnitude (ps) with the lower confidence bound $\\mu - \\kappa\\sigma$, $\\kappa = %s$. The GP posteriors at four candidate settings are %s. Which candidate does the loop run next? Give its number.", f(kappa, 0), strings.Join(rows, "; ")),
			edakit.IntStr(best+1), edakit.Sprintf("Compute $\\mu - %s\\sigma$ for each; candidate %d has the smallest bound, %s. A larger $\\kappa$ favors uncertain candidates (exploration).", f(kappa, 0), best+1, f(bestV, 0)),
			map[string]string{"lcb": f(bestV, 0)})
	default:
		var rows []string
		best := 0
		var vals []float64
		for tries := 0; tries < 50; tries++ {
			rows, vals = nil, nil
			best = 0
			for i := 0; i < 3; i++ {
				mu := fstar + float64(edakit.Int(-6, 8))
				s := float64(edakit.Int(1, 10))
				v := ei(fstar, mu, s)
				vals = append(vals, v)
				rows = append(rows, edakit.Sprintf("%d: mean %s, sd %s", i+1, f(mu, 0), f(s, 0)))
				if v > vals[best] {
					best = i
				}
			}
			ok := true
			for i, v := range vals {
				if i != best && vals[best]-v < 0.05 {
					ok = false
				}
			}
			if ok {
				break
			}
		}
		return edakit.Numeric(edakit.Sprintf("Bayesian optimization minimizes leakage power, and the best observed value is $f^* = %s$ uW. GP posteriors at three candidates are %s (uW). Which candidate maximizes expected improvement? Give its number.", f(fstar, 0), strings.Join(rows, "; ")),
			edakit.IntStr(best+1), edakit.Sprintf("EI values: %s. Candidate %d is largest. A candidate with a worse mean can win on EI when its standard deviation is large.", floats(vals, 4), best+1),
			map[string]string{"ei_best": f(vals[best], 4)})
	}
}

// analogSizing: feasibility probability of a constraint (KP 1), constrained EI (KP 2), and the
// best feasible incumbent (KP 3).
func analogSizing(ctx generator.GeneratorContext) generator.Problem {
	pmMin := float64(edakit.Pick([]int{45, 55, 60, 65}))
	switch kp(ctx) {
	case 1:
		mu := pmMin + float64(edakit.Int(-10, 15))
		s := float64(edakit.Int(2, 8))
		p := cdf((mu - pmMin) / s)
		ans := edakit.Num(100*p, 2)
		return edakit.Numeric(edakit.Sprintf("A two-stage op-amp is sized by Bayesian optimization. A separate GP models phase margin; at a candidate sizing its posterior is Gaussian with mean %s degrees and standard deviation %s degrees. The specification is a phase margin of at least %s degrees. What is the probability, in percent, that the candidate is feasible? Round to 2 decimal places.", f(mu, 0), f(s, 0), f(pmMin, 0)),
			ans, edakit.Sprintf("$P(\\mathrm{PM} \\ge %s) = \\Phi((%s - %s)/%s) = \\Phi(%s) = %s$ percent.", f(pmMin, 0), f(mu, 0), f(pmMin, 0), f(s, 0), f((mu-pmMin)/s, 4), ans),
			map[string]string{"z": f((mu-pmMin)/s, 4)})
	case 2:
		fstar := float64(edakit.Int(200, 600)) / 100
		mp := fstar + float64(edakit.Int(-60, 40))/100
		sp := float64(edakit.Int(10, 60)) / 100
		mu := pmMin + float64(edakit.Int(-6, 12))
		s := float64(edakit.Int(2, 8))
		e := ei(fstar, mp, sp)
		pf := cdf((mu - pmMin) / s)
		v := e * pf
		ans := edakit.Num(v, 4)
		return edakit.Numeric(edakit.Sprintf("An op-amp sizing loop minimizes power subject to phase margin at least %s degrees. The best feasible power so far is %s mW. At a candidate, the power GP gives mean %s mW, sd %s mW; the phase-margin GP gives mean %s degrees, sd %s degrees, independent of power. What is the constrained expected improvement, EI times probability of feasibility, in mW? Round to 4 decimal places.", f(pmMin, 0), f(fstar, 2), f(mp, 2), f(sp, 2), f(mu, 0), f(s, 0)),
			ans, edakit.Sprintf("EI $= %s$ mW (with $z = %s$). $P(\\text{feasible}) = \\Phi(%s) = %s$. Product: %s mW.", f(e, 4), f((fstar-mp)/sp, 4), f((mu-pmMin)/s, 4), f(pf, 4), ans),
			map[string]string{"ei": f(e, 4), "pf": f(pf, 4)})
	default:
		n := edakit.Scaled(ctx, 5, 8)
		type obs struct{ p, pm float64 }
		var os []obs
		for len(os) < n {
			os = append(os, obs{float64(edakit.Int(150, 500)) / 100, pmMin + float64(edakit.Int(-12, 15))})
		}
		// make the overall minimum infeasible and ensure one feasible point
		os[0] = obs{float64(edakit.Int(100, 140)) / 100, pmMin - float64(edakit.Int(2, 10))}
		os[1].pm = pmMin + float64(edakit.Int(1, 10))
		os = edakit.Shuffle(os)
		best := math.Inf(1)
		var rows []string
		for i, o := range os {
			rows = append(rows, edakit.Sprintf("run %d: %s mW, PM %s", i+1, f(o.p, 2), f(o.pm, 0)))
			if o.pm >= pmMin && o.p < best {
				best = o.p
			}
		}
		ans := edakit.Num(best, 2)
		return edakit.Numeric(edakit.Sprintf("A constrained Bayesian optimization of an op-amp minimizes power with phase margin at least %s degrees. The SPICE runs so far gave %s (PM in degrees). What incumbent value $f^*$, in mW, should the expected-improvement computation use? Round to 2 decimal places.", f(pmMin, 0), strings.Join(rows, "; ")),
			ans, edakit.Sprintf("The incumbent is the best feasible observation: discard runs with PM below %s, then take the minimum power, %s mW. Using an infeasible run as $f^*$ would make every feasible candidate look hopeless.", f(pmMin, 0), ans),
			map[string]string{"incumbent": ans})
	}
}

// flowTuning: survivors after successive-halving rungs (KP 1), total cost (KP 2), and the
// saving against full-budget evaluation (KP 3).
func flowTuning(ctx generator.GeneratorContext) generator.Problem {
	eta := edakit.Pick([]int{2, 3, 4})
	rr := edakit.Int(2, 4)
	if eta == 4 && rr > 3 {
		rr = 3
	}
	n := 1
	for i := 0; i < rr; i++ {
		n *= eta
	}
	r := float64(edakit.Int(1, 6)) / 2
	stem := edakit.Sprintf("An autotuner explores %d OpenROAD flow configurations with successive halving: every configuration first gets %s CPU-hours, then at each rung only the best $1/%d$ continue with %d times the previous budget, until one remains.", n, f(r, 1), eta, eta)
	switch kp(ctx) {
	case 1:
		k := edakit.Int(1, rr)
		left := n
		for i := 0; i < k; i++ {
			left /= eta
		}
		return edakit.Numeric(stem+edakit.Sprintf(" How many configurations are still running after %d elimination rounds?", k),
			edakit.IntStr(left), edakit.Sprintf("$%d / %d^{%d} = %d$.", n, eta, k, left), map[string]string{"eta": edakit.IntStr(eta)})
	case 2:
		tot := float64(rr+1) * float64(n) * r
		var terms []string
		cnt, bud := n, r
		for i := 0; i <= rr; i++ {
			terms = append(terms, edakit.Sprintf("%d \\cdot %s", cnt, f(bud, 1)))
			cnt /= eta
			bud *= float64(eta)
		}
		ans := edakit.Num(tot, 1)
		return edakit.Numeric(stem+" Counting each rung's budget as a fresh run (no warm start), how many CPU-hours does the whole search use? Round to 1 decimal place.",
			ans, edakit.Sprintf("Rung costs: $%s$. Each rung costs $%d \\times %s$, and there are %d rungs: %s CPU-hours.", strings.Join(terms, " + "), n, f(r, 1), rr+1, ans),
			map[string]string{"rungs": edakit.IntStr(rr + 1)})
	default:
		full := float64(n) * r * math.Pow(float64(eta), float64(rr))
		sh := float64(rr+1) * float64(n) * r
		ratio := full / sh
		ans := edakit.Num(ratio, 3)
		return edakit.Numeric(stem+" Counting each rung as a fresh run, by what factor is this cheaper than running all configurations at the final budget? Round to 3 decimal places.",
			ans, edakit.Sprintf("Full budget per config $= %s \\cdot %d^{%d} = %s$ h, so all %d cost %s h. Successive halving costs %s h. Ratio $= %s$. The saving is real only if short runs rank configurations the way full runs would.", f(r, 1), eta, rr, f(r*math.Pow(float64(eta), float64(rr)), 1), n, f(full, 1), f(sh, 1), ans),
			map[string]string{"full": f(full, 1)})
	}
}
