package sta

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func phi(x float64) float64    { return math.Exp(-x*x/2) / math.Sqrt(2*math.Pi) }
func bigPhi(x float64) float64 { return 0.5 * (1 + math.Erf(x/math.Sqrt2)) }

func corners(ctx generator.GeneratorContext) generator.Problem {
	type corner struct {
		name   string
		lo, hi float64
	}
	all := []corner{{"ss_0p72v_125c", 1.22, 1.36}, {"ss_0p72v_m40c", 1.2, 1.42}, {"tt_0p80v_25c", 0.97, 1.03}, {"ff_0p88v_m40c", 0.68, 0.78}, {"ff_0p88v_125c", 0.72, 0.82}}
	k := edakit.Scaled(ctx, 3, 5)
	cs := edakit.Shuffle(all)[:k]
	t := 50 * edakit.Int(14, 24)
	baseMax := float64(10 * edakit.Int(45, 70))
	baseMin := float64(edakit.Int(40, 90))
	var parts []string
	setups := make([]int, k)
	holds := make([]int, k)
	for i, c := range cs {
		f := c.lo + (c.hi-c.lo)*float64(edakit.Int(0, 10))/10
		dmax := int(math.Round(baseMax * f))
		dmin := int(math.Round(baseMin * f))
		su := int(math.Round(30 * f))
		th := int(math.Round(20 * f))
		setups[i] = t - su - dmax
		holds[i] = dmin - th
		parts = append(parts, edakit.Sprintf("%s: max path %d ps, min path %d ps, setup %d ps, hold %d ps", c.name, dmax, dmin, su, th))
	}
	q := edakit.Sprintf("A flop-to-flop path with ideal clocks at period %d ps is characterized at several PVT corners (path delays include clock-to-Q): %s. ", t, strings.Join(parts, "; "))
	var exp []string
	for i, c := range cs {
		exp = append(exp, edakit.Sprintf("%s setup %d, hold %d", c.name, setups[i], holds[i]))
	}
	e := "Per corner, setup slack = T - setup - max path and hold slack = min path - hold: " + strings.Join(exp, "; ") + " (ps)."
	switch ctx.KPOrdinal {
	case 1:
		i := edakit.Int(0, k-1)
		return edakit.Numeric(q+edakit.Sprintf("What is the setup slack at corner %s, in ps?", cs[i].name), edakit.IntStr(setups[i]), e, map[string]string{"setup_slack": edakit.IntStr(setups[i])})
	case 3:
		w := minInt(holds)
		return edakit.Numeric(q+"What is the worst hold slack across these corners, in ps?", edakit.IntStr(w), e, map[string]string{"worst_hold": edakit.IntStr(w)})
	}
	w := minInt(setups)
	return edakit.Numeric(q+"What is the worst setup slack across these corners, in ps?", edakit.IntStr(w), e, map[string]string{"worst_setup": edakit.IntStr(w)})
}

func mcmm(ctx generator.GeneratorContext) generator.Problem {
	if ctx.KPOrdinal == 1 {
		m, c := edakit.Int(2, 5), edakit.Int(3, 9)
		e := edakit.Int(0, m*c/3)
		v := m*c - e
		return edakit.Numeric(edakit.Sprintf("A chip has %d functional and test modes and is signed off at %d PVT and extraction corners. %d mode-corner pairs are pruned because they can never dominate. How many scenarios does the MCMM run analyze?", m, c, e), edakit.IntStr(v),
			edakit.Sprintf("Scenarios are mode-corner pairs: %d x %d - %d = %d.", m, c, e, v), map[string]string{"scenarios": edakit.IntStr(v)})
	}
	ns := edakit.Int(3, 4)
	ne := edakit.Int(3, 5)
	if ctx.KPOrdinal == 2 {
		ne = 1
	}
	sl := make([][]int, ne)
	var rows []string
	for i := range sl {
		sl[i] = make([]int, ns)
		var cells []string
		for j := range sl[i] {
			sl[i][j] = edakit.Int(-40, 60)
			cells = append(cells, edakit.Sprintf("S%d: %d", j+1, sl[i][j]))
		}
		rows = append(rows, edakit.Sprintf("endpoint E%d has %s", i+1, strings.Join(cells, ", ")))
	}
	q := edakit.Sprintf("An MCMM run reports setup slacks (ps) per scenario: %s. ", strings.Join(rows, "; "))
	if ctx.KPOrdinal == 2 {
		w := minInt(sl[0])
		return edakit.Numeric(q+"What is the merged setup slack of E1 that the closure flow must fix, in ps?", edakit.IntStr(w),
			edakit.Sprintf("An endpoint passes only if it passes in every scenario, so the merged slack is the minimum, %d ps.", w), map[string]string{"merged": edakit.IntStr(w)})
	}
	tns := 0
	var ws []string
	for i := range sl {
		w := minInt(sl[i])
		ws = append(ws, edakit.IntStr(w))
		if w < 0 {
			tns += w
		}
	}
	return edakit.Numeric(q+"Taking each endpoint's worst slack across scenarios, what is the merged TNS, in ps?", edakit.IntStr(tns),
		edakit.Sprintf("Merged slacks per endpoint are %s. Summing the negative ones gives %d ps. Summing per-scenario TNS instead would count an endpoint once per failing scenario.", strings.Join(ws, ", "), tns), map[string]string{"tns": edakit.IntStr(tns)})
}

func ocvDerate(ctx generator.GeneratorContext) generator.Problem {
	t := float64(50 * edakit.Int(12, 24))
	lc := float64(10 * edakit.Int(15, 40))
	cc := float64(10 * edakit.Int(15, 40))
	d := float64(10 * edakit.Int(35, 80))
	su := float64(edakit.Int(20, 50))
	x := float64(edakit.Int(3, 10)) / 100
	y := float64(edakit.Int(3, 10)) / 100
	q := edakit.Sprintf("A setup check uses a late derate of %s on the launch clock path and data path and an early derate of %s on the capture clock path. The launch clock path is %s ps, the data path (including clock-to-Q) %s ps, the capture clock path %s ps, the period %s ps and setup %s ps (setup is not derated). ",
		edakit.Num(1+x, 2), edakit.Num(1-y, 2), edakit.Num(lc, 0), edakit.Num(d, 0), edakit.Num(cc, 0), edakit.Num(t, 0), edakit.Num(su, 0))
	arr := (lc + d) * (1 + x)
	req := t + cc*(1-y) - su
	s := req - arr
	facts := map[string]string{"arrival": edakit.Num(arr, 2), "required": edakit.Num(req, 2), "slack": edakit.Num(s, 2)}
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(arr, 2)
		return edakit.Numeric(q+"What is the derated data arrival time at the capture flop, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("(%s + %s) x %s = %s ps.", edakit.Num(lc, 0), edakit.Num(d, 0), edakit.Num(1+x, 2), ans), facts)
	case 3:
		p := float64(10 * edakit.Int(5, int(math.Min(lc, cc))/10))
		cred := p * (x + y)
		s2 := s + cred
		ans := edakit.Num(s2, 2)
		facts["credit"] = edakit.Num(cred, 2)
		return edakit.Numeric(q+edakit.Sprintf("The first %s ps of both clock paths run through the same buffers. With clock reconvergence pessimism removal, what is the setup slack, in ps? Round to 2 decimal places.", edakit.Num(p, 0)), ans,
			edakit.Sprintf("Without CPPR the slack is %s ps. The common %s ps was taken at %s on one side and %s on the other, so CPPR adds back %s x (%s + %s) = %s ps, giving %s ps.",
				edakit.Num(s, 2), edakit.Num(p, 0), edakit.Num(1+x, 2), edakit.Num(1-y, 2), edakit.Num(p, 0), edakit.Num(x, 2), edakit.Num(y, 2), edakit.Num(cred, 2), ans), facts)
	}
	ans := edakit.Num(s, 2)
	return edakit.Numeric(q+"What is the derated setup slack, in ps? Round to 2 decimal places.", ans,
		edakit.Sprintf("Required %s + %s x %s - %s = %s ps; arrival (%s + %s) x %s = %s ps; slack %s ps.",
			edakit.Num(t, 0), edakit.Num(cc, 0), edakit.Num(1-y, 2), edakit.Num(su, 0), edakit.Num(req, 2), edakit.Num(lc, 0), edakit.Num(d, 0), edakit.Num(1+x, 2), edakit.Num(arr, 2), ans), facts)
}

// aocvTable draws a late-derate table that falls with depth and rises with distance.
func aocvTable() ([]int, []float64) {
	depths := []int{1, 2, 4, 8, 16}
	ders := make([]float64, len(depths))
	v := 1.10 + float64(edakit.Int(0, 6))/100
	for i := range ders {
		ders[i] = v
		v -= float64(edakit.Int(1, 3)) / 100
	}
	return depths, ders
}

func aocv(ctx generator.GeneratorContext) generator.Problem {
	depths, ders := aocvTable()
	var row []string
	for i := range depths {
		row = append(row, edakit.Sprintf("depth %d: %s", depths[i], edakit.Num(ders[i], 2)))
	}
	d := float64(10 * edakit.Int(30, 90))
	switch ctx.KPOrdinal {
	case 3:
		i := edakit.Int(1, len(depths)-2)
		dists := []int{500, 1000, 2000}
		j := edakit.Int(0, 1)
		inc := float64(edakit.Int(1, 3)) / 100
		inc2 := inc + float64(edakit.Int(1, 2))/100
		v11, v21 := ders[i], ders[i+1]
		if j == 1 {
			v11, v21 = v11+inc, v21+inc
		}
		v12, v22 := v11+inc2, v21+inc2
		d1, d2 := float64(depths[i]), float64(depths[i+1])
		x1, x2 := float64(dists[j]), float64(dists[j+1])
		dep := float64(edakit.Int(depths[i]+1, depths[i+1]-1))
		dist := x1 + (x2-x1)*float64(edakit.Int(1, 3))/4
		v := (v11*(d2-dep)*(x2-dist) + v12*(d2-dep)*(dist-x1) + v21*(dep-d1)*(x2-dist) + v22*(dep-d1)*(dist-x1)) / ((d2 - d1) * (x2 - x1))
		ans := edakit.Num(v, 4)
		q := edakit.Sprintf("A two-dimensional AOCV late-derate table is indexed by logic depth and by the path's bounding-box distance. Its window holds: depth %s at %s um: %s; depth %s at %s um: %s; depth %s at %s um: %s; depth %s at %s um: %s. With bilinear interpolation, what derate applies to a path of depth %s and distance %s um? Round to 4 decimal places.",
			edakit.Num(d1, 0), edakit.Num(x1, 0), edakit.Num(v11, 2), edakit.Num(d1, 0), edakit.Num(x2, 0), edakit.Num(v12, 2), edakit.Num(d2, 0), edakit.Num(x1, 0), edakit.Num(v21, 2), edakit.Num(d2, 0), edakit.Num(x2, 0), edakit.Num(v22, 2), edakit.Num(dep, 0), edakit.Num(dist, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("Fractions: depth %s, distance %s. Interpolating gives %s. Deeper paths average out random variation, so the derate falls; longer distances add systematic variation, so it rises.",
			edakit.Num((dep-d1)/(d2-d1), 2), edakit.Num((dist-x1)/(x2-x1), 2), ans), map[string]string{"derate": ans})
	case 1:
		i := edakit.Int(0, len(depths)-1)
		v := d * ders[i]
		ans := edakit.Num(v, 2)
		q := edakit.Sprintf("An AOCV late-derate table by logic depth reads %s. A data path of depth %d has a non-derated delay of %s ps. What is its derated delay, in ps? Round to 2 decimal places.", strings.Join(row, ", "), depths[i], edakit.Num(d, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("Depth %d is on the table axis, so the derate is %s: %s x %s = %s ps.", depths[i], edakit.Num(ders[i], 2), edakit.Num(d, 0), edakit.Num(ders[i], 2), ans), map[string]string{"derate": edakit.Num(ders[i], 2), "delay": ans})
	}
	i := edakit.Int(1, len(depths)-2)
	lo, hi := depths[i], depths[i+1]
	dep := edakit.Int(lo+1, hi-1)
	f := float64(dep-lo) / float64(hi-lo)
	der := ders[i] + f*(ders[i+1]-ders[i])
	v := d * der
	ans := edakit.Num(v, 2)
	q := edakit.Sprintf("An AOCV late-derate table by logic depth reads %s, with linear interpolation between depths. A data path of depth %d has a non-derated delay of %s ps. What is its derated delay, in ps? Round to 2 decimal places.", strings.Join(row, ", "), dep, edakit.Num(d, 0))
	return edakit.Numeric(q, ans, edakit.Sprintf("Depth %d lies %s of the way from %d to %d, so the derate is %s + %s x (%s - %s) = %s. The delay is %s x %s = %s ps.",
		dep, edakit.Num(f, 4), lo, hi, edakit.Num(ders[i], 2), edakit.Num(f, 4), edakit.Num(ders[i+1], 2), edakit.Num(ders[i], 2), edakit.Num(der, 4), edakit.Num(d, 0), edakit.Num(der, 4), ans), map[string]string{"derate": edakit.Num(der, 4), "delay": ans})
}

func pocv(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 3, 6)
	var mu, sig []float64
	var parts []string
	for i := 0; i < n; i++ {
		m := float64(edakit.Int(20, 90))
		s := float64(edakit.Int(2, 12))
		mu = append(mu, m)
		sig = append(sig, s)
		parts = append(parts, edakit.Sprintf("(%s, %s)", edakit.Num(m, 0), edakit.Num(s, 0)))
	}
	sumMu, sumSq, sumSig := 0.0, 0.0, 0.0
	for i := range mu {
		sumMu += mu[i]
		sumSq += sig[i] * sig[i]
		sumSig += sig[i]
	}
	ps := math.Sqrt(sumSq)
	nsig := float64(edakit.Pick([]int{3, 3, 2, 4}))
	q := edakit.Sprintf("A data path has %d arcs. Their LVF (mean, sigma) delays in ps are %s, and their random variations are independent. ", n, strings.Join(parts, ", "))
	facts := map[string]string{"mean": edakit.Num(sumMu, 0), "sigma": edakit.Num(ps, 2)}
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(ps, 2)
		return edakit.Numeric(q+"What is the sigma of the path delay, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("Independent variances add: $\\sigma = \\sqrt{%s}$ = %s ps.", edakit.Num(sumSq, 0), ans), facts)
	case 3:
		pess := nsig*sumSig - nsig*ps
		ans := edakit.Num(pess, 2)
		return edakit.Numeric(q+edakit.Sprintf("A flow that margins each arc separately at %s sigma adds the per-arc margins. By how many ps does that overestimate the %s-sigma path delay computed from the root-sum-square? Round to 2 decimal places.", edakit.Num(nsig, 0), edakit.Num(nsig, 0)), ans,
			edakit.Sprintf("Per-arc stacking adds %s x %s = %s ps; POCV adds %s x %s = %s ps. The difference is %s ps, and it grows with path depth.",
				edakit.Num(nsig, 0), edakit.Num(sumSig, 0), edakit.Num(nsig*sumSig, 2), edakit.Num(nsig, 0), edakit.Num(ps, 2), edakit.Num(nsig*ps, 2), ans), facts)
	}
	v := sumMu + nsig*ps
	ans := edakit.Num(v, 2)
	return edakit.Numeric(q+edakit.Sprintf("What is the %s-sigma late delay of the path, in ps? Round to 2 decimal places.", edakit.Num(nsig, 0)), ans,
		edakit.Sprintf("Means add to %s ps; sigma is $\\sqrt{%s}$ = %s ps; %s + %s x %s = %s ps.", edakit.Num(sumMu, 0), edakit.Num(sumSq, 0), edakit.Num(ps, 2), edakit.Num(sumMu, 0), edakit.Num(nsig, 0), edakit.Num(ps, 2), ans), facts)
}

type canon struct{ a0, a1, a2, ar float64 }

func drawCanon() canon {
	return canon{float64(10 * edakit.Int(10, 60)), float64(edakit.Int(-4, 12)), float64(edakit.Int(-4, 12)), float64(edakit.Int(1, 10))}
}

func (c canon) sigma() float64 { return math.Sqrt(c.a1*c.a1 + c.a2*c.a2 + c.ar*c.ar) }

func term(x float64, v string) string {
	if x < 0 {
		return edakit.Sprintf(" - %s %s", edakit.Num(-x, 0), v)
	}
	return edakit.Sprintf(" + %s %s", edakit.Num(x, 0), v)
}

func (c canon) String() string {
	return edakit.Num(c.a0, 0) + term(c.a1, "X_1") + term(c.a2, "X_2") + term(c.ar, "R")
}

func ssta(ctx generator.GeneratorContext) generator.Problem {
	a, b := drawCanon(), drawCanon()
	intro := "In first-order canonical form, $X_1$ and $X_2$ are shared global variation sources and each R is a private random term; all are independent standard normals, and coefficients are in ps. "
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(a.sigma(), 2)
		return edakit.Numeric(intro+edakit.Sprintf("An arrival time is $A = %s$. What is its standard deviation, in ps? Round to 2 decimal places.", a), ans,
			edakit.Sprintf("$\\sigma_A = \\sqrt{%s^2 + %s^2 + %s^2}$ = %s ps.", edakit.Num(a.a1, 0), edakit.Num(a.a2, 0), edakit.Num(a.ar, 0), ans), map[string]string{"sigma": ans})
	case 3:
		cov := a.a1*b.a1 + a.a2*b.a2
		rho := cov / (a.sigma() * b.sigma())
		ans := edakit.Num(rho, 3)
		return edakit.Numeric(intro+edakit.Sprintf("Two arrival times are $A = %s$ and $B = %s$, with independent random terms. What is their correlation coefficient? Round to 3 decimal places.", a, b), ans,
			edakit.Sprintf("Only shared sources correlate: cov = %s x %s + %s x %s = %s. Divide by $\\sigma_A \\sigma_B$ = %s x %s to get %s.",
				edakit.Num(a.a1, 0), edakit.Num(b.a1, 0), edakit.Num(a.a2, 0), edakit.Num(b.a2, 0), edakit.Num(cov, 0), edakit.Num(a.sigma(), 3), edakit.Num(b.sigma(), 3), ans), map[string]string{"cov": edakit.Num(cov, 0), "rho": ans})
	}
	c := canon{a.a0 + b.a0, a.a1 + b.a1, a.a2 + b.a2, math.Sqrt(a.ar*a.ar + b.ar*b.ar)}
	ans := edakit.Num(c.sigma(), 2)
	return edakit.Numeric(intro+edakit.Sprintf("An arrival time $A = %s$ feeds an arc with delay $D = %s$, whose random term is independent of A's. What is the standard deviation of $A + D$, in ps? Round to 2 decimal places.", a, b), ans,
		edakit.Sprintf("Sum: mean %s, sensitivities %s and %s, random part $\\sqrt{%s^2 + %s^2}$ = %s. Then $\\sigma$ = %s ps.",
			edakit.Num(c.a0, 0), edakit.Num(c.a1, 0), edakit.Num(c.a2, 0), edakit.Num(a.ar, 0), edakit.Num(b.ar, 0), edakit.Num(c.ar, 3), ans), map[string]string{"sigma": ans})
}

func clarkMax(ctx generator.GeneratorContext) generator.Problem {
	var ma, mb, sa, sb, rho, theta float64
	for {
		ma, mb = float64(10*edakit.Int(30, 60)), float64(10*edakit.Int(30, 60))
		sa, sb = float64(edakit.Int(5, 30)), float64(edakit.Int(5, 30))
		rho = float64(edakit.Int(0, 8)) / 10
		theta = math.Sqrt(sa*sa + sb*sb - 2*rho*sa*sb)
		if theta > 2 && math.Abs(ma-mb) < 2.5*theta {
			break
		}
	}
	alpha := (ma - mb) / theta
	q := edakit.Sprintf("Two arrival times converge at a gate: A has mean %s ps and sigma %s ps, B has mean %s ps and sigma %s ps, and their correlation is %s. Clark's method uses $\\theta = \\sqrt{\\sigma_A^2 + \\sigma_B^2 - 2\\rho\\sigma_A\\sigma_B}$ and $\\alpha = (\\mu_A - \\mu_B)/\\theta$. ",
		edakit.Num(ma, 0), edakit.Num(sa, 0), edakit.Num(mb, 0), edakit.Num(sb, 0), edakit.Num(rho, 1))
	facts := map[string]string{"theta": edakit.Num(theta, 4), "alpha": edakit.Num(alpha, 4)}
	switch ctx.KPOrdinal {
	case 1:
		ans := edakit.Num(theta, 2)
		return edakit.Numeric(q+"What is $\\theta$, in ps? Round to 2 decimal places.", ans,
			edakit.Sprintf("$\\theta^2$ = %s + %s - 2 x %s x %s x %s = %s, so $\\theta$ = %s ps.", edakit.Num(sa*sa, 0), edakit.Num(sb*sb, 0), edakit.Num(rho, 1), edakit.Num(sa, 0), edakit.Num(sb, 0), edakit.Num(theta*theta, 2), ans), facts)
	case 3:
		P := edakit.Round(bigPhi(alpha), 4)
		pd := edakit.Round(phi(alpha), 4)
		m := ma*P + mb*(1-P) + theta*pd
		ans := edakit.Num(m, 2)
		return edakit.Numeric(q+edakit.Sprintf("For this pair, $\\Phi(\\alpha) = %s$ and $\\varphi(\\alpha) = %s$. Using $E[\\max] = \\mu_A \\Phi(\\alpha) + \\mu_B \\Phi(-\\alpha) + \\theta \\varphi(\\alpha)$, what is the mean of the latest arrival, in ps? Round to 2 decimal places.", edakit.Num(P, 4), edakit.Num(pd, 4)), ans,
			edakit.Sprintf("$\\theta$ = %s ps. %s x %s + %s x %s + %s x %s = %s ps, above the larger mean because the max of two uncertain arrivals is pushed later.",
				edakit.Num(theta, 4), edakit.Num(ma, 0), edakit.Num(P, 4), edakit.Num(mb, 0), edakit.Num(1-P, 4), edakit.Num(theta, 4), edakit.Num(pd, 4), ans), facts)
	}
	ans := edakit.Num(bigPhi(alpha), 2)
	return edakit.Numeric(q+"What is the tightness probability $P(A > B) = \\Phi(\\alpha)$? Round to 2 decimal places.", ans,
		edakit.Sprintf("$\\theta$ = %s ps and $\\alpha$ = %s, so $\\Phi(\\alpha)$ = %s. This is the weight A gets when the max is re-expressed in canonical form.", edakit.Num(theta, 3), edakit.Num(alpha, 3), ans), facts)
}

func timingYield(ctx generator.GeneratorContext) generator.Problem {
	mu := float64(10 * edakit.Int(60, 120))
	sg := float64(edakit.Int(10, 40))
	switch ctx.KPOrdinal {
	case 1:
		t := mu + float64(edakit.Int(-10, 60))
		z := (t - mu) / sg
		ans := edakit.Num(z, 2)
		return edakit.Numeric(edakit.Sprintf("SSTA gives a circuit delay with mean %s ps and sigma %s ps. The clock period is %s ps. What is the z-score of the period? Round to 2 decimal places.", edakit.Num(mu, 0), edakit.Num(sg, 0), edakit.Num(t, 0)), ans,
			edakit.Sprintf("z = (%s - %s) / %s = %s.", edakit.Num(t, 0), edakit.Num(mu, 0), edakit.Num(sg, 0), ans), map[string]string{"z": ans})
	case 3:
		type tgt struct {
			p, z string
			zv   float64
		}
		tg := edakit.Pick([]tgt{{"95%", "1.645", 1.645}, {"99%", "2.326", 2.326}, {"99.9%", "3.090", 3.090}, {"90%", "1.282", 1.282}, {"97.5%", "1.960", 1.960}})
		v := mu + tg.zv*sg
		ans := edakit.Num(v, 1)
		return edakit.Numeric(edakit.Sprintf("SSTA gives a circuit delay with mean %s ps and sigma %s ps. The product needs %s timing yield, and the standard normal quantile for %s is %s. What is the shortest clock period that meets the target, in ps? Round to 1 decimal place.", edakit.Num(mu, 0), edakit.Num(sg, 0), tg.p, tg.p, tg.z), ans,
			edakit.Sprintf("T = %s + %s x %s = %s ps.", edakit.Num(mu, 0), tg.z, edakit.Num(sg, 0), ans), map[string]string{"period": ans})
	}
	z := float64(edakit.Int(-5, 30)) / 10
	t := mu + z*sg
	y := 100 * bigPhi(z)
	ans := edakit.Num(y, 1)
	return edakit.Numeric(edakit.Sprintf("SSTA gives a circuit delay with mean %s ps and sigma %s ps, approximately Gaussian. What fraction of dies meet a clock period of %s ps, in percent? Round to 1 decimal place.", edakit.Num(mu, 0), edakit.Num(sg, 0), edakit.Num(t, 1)), ans,
		edakit.Sprintf("z = (%s - %s) / %s = %s, and $\\Phi(%s)$ = %s, so the yield is %s%%.", edakit.Num(t, 1), edakit.Num(mu, 0), edakit.Num(sg, 0), edakit.Num(z, 1), edakit.Num(z, 1), edakit.Num(bigPhi(z), 4), ans), map[string]string{"z": edakit.Num(z, 1), "yield": ans})
}
