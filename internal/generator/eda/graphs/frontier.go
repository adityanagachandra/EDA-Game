package graphs

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// netWeightCut is the computed part of the net-weight seam: the cut a NetWeightProvider induces.
func netWeightCut(ctx generator.GeneratorContext) generator.Problem {
	for {
		n := edakit.Scaled(ctx, 6, 9)
		h := randHyper(n, edakit.Scaled(ctx, 5, 8), 4)
		side := randSides(n)
		crit := edakit.Int(2, 5)
		kind := make([]int, len(h.nets)) // 0 plain, 1 critical, 2 clock
		var ct, kt []string
		for j := range h.nets {
			switch edakit.Int(0, 5) {
			case 0, 1:
				kind[j] = 1
				ct = append(ct, edakit.Sprintf("n%d", j+1))
			case 2:
				kind[j] = 2
				kt = append(kt, edakit.Sprintf("n%d", j+1))
			}
		}
		if len(ct) == 0 || len(kt) == 0 {
			continue
		}
		correct, plain, noClock, noCrit := 0, 0, 0, 0
		for j, net := range h.nets {
			if !isCut(net, side) {
				continue
			}
			plain++
			switch kind[j] {
			case 0:
				correct++
				noClock++
				noCrit++
			case 1:
				correct += crit
				noClock += crit
				noCrit++
			case 2:
				noClock++
			}
		}
		vals := map[int]bool{correct: true}
		var ds []string
		for _, v := range []int{plain, noClock, noCrit, correct + crit, correct + 1} {
			if !vals[v] {
				vals[v] = true
				ds = append(ds, edakit.Sprintf("Weighted cut %d", v))
			}
		}
		if len(ds) < 3 {
			continue
		}
		q := "Nets: " + netText(h, false) + ". The partition is A = " + sideText(side, 0) + ", B = " + sideText(side, 1) + edakit.Sprintf(". A customer's NetWeightProvider plug-in returns weight %d for timing-critical nets (%s), weight 0 for clock nets (%s, which clock-tree synthesis rebuilds anyway) and weight 1 otherwise. What weighted cut does the unchanged core report?", crit, strings.Join(ct, ", "), strings.Join(kt, ", "))
		return edakit.Choice(q, edakit.Sprintf("Weighted cut %d", correct),
			edakit.Sprintf("Of %d cut nets, each contributes its provider weight: critical nets %d each, clock nets 0, others 1, for %d. The core's cut accounting is unchanged; only the weights arrive through the seam.", plain, crit, correct),
			ds...)
	}
}

// diffVerdict is the computed part of differential cut testing: recompute and compare.
func diffVerdict(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 6, 9)
	h := randHyper(n, edakit.Scaled(ctx, 5, 8), 4)
	side := randSides(n)
	c := cutSize(h, side, false)
	delta := edakit.Pick([]int{-2, -1, 0, 0, 1, 2})
	r := c + delta
	if r < 0 {
		r, delta = c, 0
	}
	q := "Nets: " + netText(h, false) + ". After a refactor of FM's incremental gain updates, the engine returns A = " + sideText(side, 0) + ", B = " + sideText(side, 1) + edakit.Sprintf(" and its internal counter reports a cut of %d. The differential test recomputes the cut from the returned partition. What should it report?", r)
	pass := edakit.Sprintf("Pass: the recomputed cut is %d, matching the counter", c)
	off := edakit.Sprintf("Fail: the recomputed cut is %d, so the counter is off by %d", c, iabs(delta))
	trust := "Pass: the counter is the engine's own measurement, so it is trusted"
	tol := "Pass: a difference of up to 2 nets is within normal tolerance"
	if delta == 0 {
		return edakit.Choice(q, pass,
			edakit.Sprintf("Counting nets with pins on both sides of the returned partition gives %d, equal to the counter. The check passes only because the recount was done independently.", c),
			edakit.Sprintf("Fail: the recomputed cut is %d, so the counter is off by 1", c+1), trust, tol)
	}
	return edakit.Choice(q, off,
		edakit.Sprintf("The independent recount gives %d but the counter says %d. A cut is an exact integer, so any mismatch is a bug in the incremental bookkeeping, not noise.", c, r),
		pass, trust, tol)
}

// fastDecision is the computed part of the fast-mode loop: keep the fast estimate or recurse.
func fastDecision(ctx generator.GeneratorContext) generator.Problem {
	ref := edakit.Int(300, 900)
	tau := edakit.Pick([]int{3, 5, 8})
	sigma := edakit.Pick([]int{10, 20, 30})
	runs := make([]int, 5)
	lo, hi, sum := 1<<30, 0, 0
	for {
		lo, hi, sum = 1<<30, 0, 0
		bias := edakit.Int(-2, 12)
		spread := edakit.Int(2, 40)
		for i := range runs {
			runs[i] = ref + ref*bias/100 + edakit.Int(0, spread)
			lo, hi, sum = min(lo, runs[i]), max(hi, runs[i]), sum+runs[i]
		}
		mean := float64(sum) / 5
		gap := (mean - float64(ref)) / float64(ref) * 100
		if gap > float64(tau)-0.3 && gap < float64(tau)+0.3 {
			continue
		}
		if hi-lo == sigma {
			continue
		}
		break
	}
	mean := float64(sum) / 5
	gap := (mean - float64(ref)) / float64(ref) * 100
	spread := hi - lo
	gs := edakit.Num(gap, 1)
	q := edakit.Sprintf("An agent's modeling loop uses a fast multilevel mode (one V-cycle, no refinement at the two finest levels). The last accurate run on this block cut %d nets. Five fast runs with different seeds cut %s. Policy: keep using the fast mode only if its mean cut is within %d%% of the accurate reference and its seed spread (max minus min) is at most %d nets; otherwise recurse to the accurate mode. What does the loop do?", ref, intsText(runs), tau, sigma)
	keep := edakit.Sprintf("Keep the fast mode: gap %s%% and spread %d are within limits", gs, spread)
	rGap := edakit.Sprintf("Recurse: the mean gap of %s%% exceeds %d%%, though the spread is fine", gs, tau)
	rSpr := edakit.Sprintf("Recurse: the spread of %d exceeds %d, though the gap is fine", spread, sigma)
	rBoth := edakit.Sprintf("Recurse: both the gap of %s%% and the spread of %d exceed their limits", gs, spread)
	gapBad := gap > float64(tau)
	sprBad := spread > sigma
	ans := keep
	switch {
	case gapBad && sprBad:
		ans = rBoth
	case gapBad:
		ans = rGap
	case sprBad:
		ans = rSpr
	}
	var ds []string
	for _, o := range []string{keep, rGap, rSpr, rBoth} {
		if o != ans {
			ds = append(ds, o)
		}
	}
	return edakit.Choice(q, ans,
		edakit.Sprintf("Mean fast cut %s, so the gap is (%s - %d)/%d = %s%%. Spread %d - %d = %d. Gap test %s, spread test %s.", edakit.Num(mean, 1), edakit.Num(mean, 1), ref, ref, gs, hi, lo, spread, passFail(!gapBad), passFail(!sprBad)),
		ds...)
}

func passFail(ok bool) string {
	if ok {
		return "passes"
	}
	return "fails"
}
