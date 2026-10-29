package stat

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var paramNames = []string{"$V_{th}$", "$L_{eff}$", "$t_{ox}$", "$W$"}
var paramUnits = []string{"mV", "nm", "nm", "nm"}

// linearModel draws n sensitivities (ps per unit) and sigmas for the first n parameters.
func linearModel(n int) (a, s []float64, desc string) {
	var parts []string
	for i := 0; i < n; i++ {
		ai := float64(edakit.Int(-12, 12)) / 4
		for ai == 0 {
			ai = float64(edakit.Int(-12, 12)) / 4
		}
		si := float64(edakit.Int(2, 20)) / 2
		a = append(a, ai)
		s = append(s, si)
		parts = append(parts, edakit.Sprintf("%s: sensitivity %s ps per %s, $\\sigma = %s$ %s", paramNames[i], f(ai, 2), paramUnits[i], f(si, 1), paramUnits[i]))
	}
	return a, s, strings.Join(parts, "; ")
}

// sensitivity: first-order sigma of a delay (KP 1), the variance share of one parameter (KP 2),
// and a central-difference sensitivity turned into a sigma contribution (KP 3).
func sensitivity(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1, 2:
		n := edakit.Scaled(ctx, 2, 4)
		a, s, desc := linearModel(n)
		v := 0.0
		var terms []string
		for i := range a {
			v += a[i] * a[i] * s[i] * s[i]
			terms = append(terms, f(a[i]*a[i]*s[i]*s[i], 4))
		}
		stem := "A gate delay depends on independent process parameters, linearized about nominal. " + desc + "."
		if ctx.KPOrdinal == 1 || (ctx.KPOrdinal == 0 && edakit.Int(0, 1) == 0) {
			ans := edakit.Num(math.Sqrt(v), 3)
			return edakit.Numeric(stem+" What is the first-order standard deviation of the delay, in ps? Round to 3 decimal places.",
				ans, edakit.Sprintf("$\\sigma_D^2 = \\sum_i (\\partial D/\\partial p_i)^2\\sigma_i^2 = %s = %s$, so $\\sigma_D = %s$ ps.", strings.Join(terms, " + "), f(v, 4), ans),
				map[string]string{"var": f(v, 4)})
		}
		j := edakit.Int(0, n-1)
		share := 100 * a[j] * a[j] * s[j] * s[j] / v
		ans := edakit.Num(share, 1)
		return edakit.Numeric(stem+edakit.Sprintf(" What percentage of the delay variance comes from %s? Round to 1 decimal place.", paramNames[j]),
			ans, edakit.Sprintf("Its term is $(%s \\cdot %s)^2 = %s$ out of a total variance of %s, which is %s percent. The sign of the sensitivity does not matter.", f(a[j], 2), f(s[j], 1), f(a[j]*a[j]*s[j]*s[j], 4), f(v, 4), ans),
			map[string]string{"var": f(v, 4)})
	default:
		d0 := float64(edakit.Int(400, 900)) / 10
		h := float64(edakit.Pick([]int{1, 2, 5, 10}))
		slope := float64(edakit.Int(2, 30)) / 20
		curv := float64(edakit.Int(-5, 5)) / 1000
		dp := d0 + slope*h + curv*h*h
		dm := d0 - slope*h + curv*h*h
		sig := float64(edakit.Int(3, 30))
		// Use the rounded readings the stem shows.
		dpr, dmr := edakit.Round(dp, 2), edakit.Round(dm, 2)
		g := (dpr - dmr) / (2 * h)
		c := math.Abs(g) * sig
		ans := edakit.Num(c, 3)
		return edakit.Numeric(edakit.Sprintf("A cell's delay is simulated with $V_{th}$ raised by %s mV (delay %s ps) and lowered by %s mV (delay %s ps). $V_{th}$ has $\\sigma = %s$ mV. Using the central-difference sensitivity, what is the first-order contribution of $V_{th}$ to the delay sigma, in ps? Round to 3 decimal places.", f(h, 0), f(dpr, 2), f(h, 0), f(dmr, 2), f(sig, 0)),
			ans, edakit.Sprintf("$\\partial D/\\partial V_{th} \\approx (%s - %s)/(2 \\cdot %s) = %s$ ps per mV. The contribution is $|%s| \\times %s = %s$ ps.", f(dpr, 2), f(dmr, 2), f(h, 0), f(g, 5), f(g, 5), f(sig, 0), ans),
			map[string]string{"slope": f(g, 5)})
	}
}

// sobol: first-order index of an additive model (KP 1), first-order index with an interaction
// (KP 2), and the total index with an interaction (KP 3).
func sobol(ctx generator.GeneratorContext) generator.Problem {
	a1, a2 := float64(edakit.Int(1, 6)), float64(edakit.Int(1, 6))
	s1, s2 := float64(edakit.Int(1, 4)), float64(edakit.Int(1, 4))
	switch kp(ctx) {
	case 1:
		a3, s3 := float64(edakit.Int(1, 6)), float64(edakit.Int(1, 4))
		v := a1*a1*s1*s1 + a2*a2*s2*s2 + a3*a3*s3*s3
		si := a1 * a1 * s1 * s1 / v
		ans := edakit.Num(si, 3)
		return edakit.Numeric(edakit.Sprintf("A response is $Y = %s X_1 + %s X_2 + %s X_3$ with independent zero-mean inputs of standard deviations %s, %s and %s. What is the first-order Sobol index of $X_1$? Round to 3 decimal places.", f(a1, 0), f(a2, 0), f(a3, 0), f(s1, 0), f(s2, 0), f(s3, 0)),
			ans, edakit.Sprintf("For an additive linear model $S_1 = a_1^2\\sigma_1^2/\\mathrm{Var}(Y) = %s/%s = %s$.", f(a1*a1*s1*s1, 0), f(v, 0), ans), map[string]string{"var": f(v, 0)})
	default:
		c := float64(edakit.Int(1, 5))
		inter := c * c * s1 * s1 * s2 * s2
		v := a1*a1*s1*s1 + a2*a2*s2*s2 + inter
		stem := edakit.Sprintf("A response is $Y = %s X_1 + %s X_2 + %s X_1 X_2$ with independent zero-mean Gaussian inputs, $\\sigma_1 = %s$ and $\\sigma_2 = %s$.", f(a1, 0), f(a2, 0), f(c, 0), f(s1, 0), f(s2, 0))
		work := edakit.Sprintf("$\\mathrm{Var}(Y) = a_1^2\\sigma_1^2 + a_2^2\\sigma_2^2 + c^2\\sigma_1^2\\sigma_2^2 = %s + %s + %s = %s$.", f(a1*a1*s1*s1, 0), f(a2*a2*s2*s2, 0), f(inter, 0), f(v, 0))
		if kp(ctx) == 2 && ctx.KPOrdinal != 3 {
			ans := edakit.Num(a1*a1*s1*s1/v, 3)
			return edakit.Numeric(stem+" What is the first-order Sobol index of $X_1$? Round to 3 decimal places.",
				ans, work+edakit.Sprintf(" $E[Y\\mid X_1] = %s X_1$ because the other terms average to zero, so $S_1 = %s/%s = %s$.", f(a1, 0), f(a1*a1*s1*s1, 0), f(v, 0), ans), map[string]string{"var": f(v, 0)})
		}
		ans := edakit.Num((a1*a1*s1*s1+inter)/v, 3)
		return edakit.Numeric(stem+" What is the total Sobol index of $X_1$ (its first-order effect plus every interaction it takes part in)? Round to 3 decimal places.",
			ans, work+edakit.Sprintf(" $S_{T1} = (a_1^2\\sigma_1^2 + c^2\\sigma_1^2\\sigma_2^2)/\\mathrm{Var}(Y) = %s/%s = %s$.", f(a1*a1*s1*s1+inter, 0), f(v, 0), ans), map[string]string{"var": f(v, 0)})
	}
}

// worstCase: the k-sigma worst-case stack (KP 1), the statistical k-sigma spread (KP 2), and the
// effective sigma of the worst-case corner (KP 3).
func worstCase(ctx generator.GeneratorContext) generator.Problem {
	n := edakit.Scaled(ctx, 2, 4)
	a, s, desc := linearModel(n)
	k := edakit.Pick([]float64{2, 3})
	wc, v := 0.0, 0.0
	for i := range a {
		wc += math.Abs(a[i]) * s[i]
		v += a[i] * a[i] * s[i] * s[i]
	}
	stem := "A delay depends linearly on independent Gaussian parameters. " + desc + "."
	switch kp(ctx) {
	case 1:
		ans := edakit.Num(k*wc, 3)
		return edakit.Numeric(stem+edakit.Sprintf(" A worst-case analysis pushes every parameter %s sigma in its delay-increasing direction at once. How much does the delay rise, in ps? Round to 3 decimal places.", f(k, 0)),
			ans, edakit.Sprintf("Each parameter adds $k|a_i|\\sigma_i$: the total is $%s \\times %s = %s$ ps.", f(k, 0), f(wc, 4), ans), map[string]string{"sum_abs": f(wc, 4)})
	case 2:
		ans := edakit.Num(k*math.Sqrt(v), 3)
		return edakit.Numeric(stem+edakit.Sprintf(" What is the %s-sigma statistical spread of the delay, in ps? Round to 3 decimal places.", f(k, 0)),
			ans, edakit.Sprintf("$\\sigma_D = \\sqrt{\\sum a_i^2\\sigma_i^2} = %s$, so $%s\\sigma_D = %s$ ps, smaller than the worst-case stack of %s ps.", f(math.Sqrt(v), 4), f(k, 0), ans, f(k*wc, 3)), map[string]string{"sigma": f(math.Sqrt(v), 4)})
	default:
		beta := k * wc / math.Sqrt(v)
		ans := edakit.Num(beta, 3)
		return edakit.Numeric(stem+edakit.Sprintf(" The worst-case corner sets every parameter at %s sigma in the bad direction. How many sigmas of the actual delay distribution is that corner's delay above the mean? Round to 3 decimal places.", f(k, 0)),
			ans, edakit.Sprintf("The corner adds $%s \\times %s$ ps against a delay sigma of %s ps, so it sits at %s sigma: a much rarer event than %s sigma.", f(k, 0), f(wc, 4), f(math.Sqrt(v), 4), ans, f(k, 0)), map[string]string{"sigma": f(math.Sqrt(v), 4)})
	}
}

var cornerItems = []edakit.Item{
	{KP: 1, Q: "What does the corner label SF mean for a CMOS process?",
		A: "Slow NMOS and fast PMOS", D: []string{"Slow PMOS and fast NMOS", "Slow at low temperature, fast at high temperature", "Slow front-end-of-line, fast back-end-of-line"},
		Why: "By convention the first letter is the NMOS and the second the PMOS, so SF is slow NMOS with fast PMOS."},
	{KP: 1, Q: "A process corner such as SS is defined how, statistically?",
		A: "Device model parameters shifted together to a chosen sigma level of the global distribution", D: []string{"The slowest die ever measured", "The mean of all dies at high temperature", "A Monte Carlo sample of local mismatch"},
		Why: "Corner models move the global parameters jointly, typically near 3 sigma, so every transistor on the die is equally slow."},
	{KP: 1, Q: "Which corner is used to check hold time in a digital signoff flow, in the usual case?",
		A: "A fast process corner, where short paths are fastest", D: []string{"A slow process corner, where paths are longest", "The typical corner only", "The SF skewed corner only"},
		Why: "Hold fails when data arrives too early, which happens when the logic is fastest."},
	{KP: 1, Q: "Why do skewed corners FS and SF matter for an SRAM cell and not much for a static inverter chain's delay?",
		A: "Ratioed circuits depend on the relative strength of NMOS and PMOS, which skewed corners stress", D: []string{"Inverter chains have no PMOS", "Skewed corners only affect wires", "SRAM cells only use NMOS devices"},
		Why: "Read and write margins compare pull-up, pull-down and access device strengths, so an imbalance between N and P is the dangerous case."},
	{KP: 2, Q: "A flow checks 5 process corners, 3 supply voltages, 3 temperatures and 5 interconnect corners. How many full combinations are there?",
		A: "225", D: []string{"16", "45", "75"},
		Why: "Independent choices multiply: 5 x 3 x 3 x 5 = 225, which is why flows prune to a dominant subset."},
	{KP: 2, Q: "At advanced nodes, delay at low supply can be worst at low temperature rather than high. What does that mean for corner selection?",
		A: "The slow corner must be checked at both temperature extremes, because the worst temperature depends on voltage", D: []string{"Temperature can be dropped from the corner set", "Only high temperature needs checking", "Only the fast corner is affected"},
		Why: "Temperature inversion makes the worst-case temperature voltage-dependent, so a single assumed hot slow corner can miss the true worst case."},
	{KP: 2, Q: "A design is clean at every global corner but fails in silicon on a few paths. Which variation did the corners leave out?",
		A: "Local within-die variation that makes devices on one path differ from each other", D: []string{"Global die-to-die variation", "Supply voltage variation", "Temperature variation"},
		Why: "Corners move every device together. Within-die random and spatial variation needs on-chip variation margins or statistical analysis."},
	{KP: 2, Q: "Why does an interconnect corner such as RC-worst differ from C-worst?",
		A: "Wire resistance and capacitance move in opposite directions with metal thickness, so the worst case depends on whether R or C dominates the path", D: []string{"They are identical names for one corner", "RC-worst only affects clock nets", "C-worst ignores coupling capacitance"},
		Why: "A thicker wire has lower resistance but more sidewall capacitance. Long resistive nets are worst at RC-worst, short gate-dominated nets at C-worst."},
	{KP: 3, Q: "Two independent parameters each sit at their 3-sigma worst value in a corner. If the delay depends equally on both, about how many sigmas of the delay distribution is the corner?",
		A: "About 4.24 sigma", D: []string{"3 sigma", "6 sigma", "1.5 sigma"},
		Why: "The corner shifts delay by 3a + 3a = 6a while the delay sigma is $a\\sqrt{2}$. The ratio is $6/\\sqrt{2} \\approx 4.24$: corners compound pessimism."},
	{KP: 3, Q: "A team sets every corner at 3 sigma for each of ten independent parameters. What is the consequence?",
		A: "The corner is far rarer than 3 sigma of the actual performance, so the design is over-margined", D: []string{"The corner exactly matches 3-sigma performance", "The corner is optimistic", "Ten parameters make the corner more typical"},
		Why: "Stacking many independent worst values is an event much less likely than any single one, which wastes area and power."},
	{KP: 3, Q: "Corners are derived for drive current. A designer uses the same corners to bound leakage power. What risk does that carry?",
		A: "The parameter combination that is worst for speed need not be worst for leakage, so the bound may not hold", D: []string{"None: corners bound every metric", "Leakage is independent of process", "The corner overstates leakage by exactly 3x"},
		Why: "A corner is worst-case for the performance it was built for. Another metric with different sensitivities can peak at a different point."},
	{KP: 3, Q: "When is a corner-based analysis the right tool instead of Monte Carlo?",
		A: "When global variation dominates and the metric is monotone in the corner parameters", D: []string{"When local mismatch dominates", "When the metric is a ratio of matched devices", "When the failure rate target is 6 sigma"},
		Why: "Corners are cheap and exact for monotone metrics driven by shared shifts. Mismatch-driven and non-monotone metrics need statistical simulation."},
}
