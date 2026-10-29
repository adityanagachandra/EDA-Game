// Package stat generates variability-statistics problems: process variation, Gaussian algebra
// for statistical timing, Monte Carlo and rare-event estimation, sensitivity analysis, model
// fitting and design of experiments, yield, reliability, and the agent-first frontier.
package stat

import (
	"math"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the variability-statistics generators to reg.
func Register(reg *generator.Registry) {
	// variation
	reg.Register("stat.var.process_variation", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: processVariationItems}, edakit.Func(varianceSplitChoice), &edakit.Bank{Items: processVariationItems}}})
	reg.Register("stat.var.pelgrom", edakit.Func(pelgrom))
	reg.Register("stat.var.spatial_correlation", edakit.Func(spatialCorrelation))
	reg.Register("stat.var.pca", edakit.Func(pcaVariation))
	reg.Register("stat.var.spc", edakit.Func(spc))
	// gaussian algebra
	reg.Register("stat.gauss.sigma_to_prob", edakit.Func(sigmaToProb))
	reg.Register("stat.gauss.linear_combo", edakit.Func(linearCombo))
	reg.Register("stat.gauss.sum_max", edakit.Func(clarkMax))
	reg.Register("stat.gauss.canonical_form", edakit.Func(canonicalForm))
	reg.Register("stat.gauss.timing_yield", edakit.Func(timingYield))
	// monte carlo
	reg.Register("stat.mc.monte_carlo", edakit.Func(monteCarlo))
	reg.Register("stat.mc.sample_size", edakit.Func(sampleSize))
	reg.Register("stat.mc.quasi_mc", edakit.Func(quasiMC))
	reg.Register("stat.mc.importance_sampling", edakit.Func(importanceSampling))
	reg.Register("stat.mc.mpfp", edakit.Func(mpfp))
	reg.Register("stat.mc.rare_event", &edakit.Mix{Parts: []generator.Generator{
		edakit.Func(rareEventBudgetChoice), &edakit.Bank{Items: rareEventItems}, &edakit.Bank{Items: rareEventItems}}})
	reg.Register("stat.mc.extreme_value", edakit.Func(extremeValue))
	// sensitivity
	reg.Register("stat.sens.sensitivity", edakit.Func(sensitivity))
	reg.Register("stat.sens.sobol", edakit.Func(sobol))
	reg.Register("stat.sens.worst_case", edakit.Func(worstCase))
	reg.Register("stat.sens.corners", &edakit.Bank{Items: cornerItems})
	// fitting
	reg.Register("stat.fit.regression", edakit.Func(regression))
	reg.Register("stat.fit.doe_runs", edakit.Func(doeRuns))
	reg.Register("stat.fit.factorial_effects", edakit.Func(factorialEffects))
	reg.Register("stat.fit.response_surface", edakit.Func(responseSurface))
	reg.Register("stat.fit.surrogate_validation", edakit.Func(surrogateValidation))
	// yield
	reg.Register("stat.yield.yield_model", edakit.Func(yieldModel))
	reg.Register("stat.yield.poisson", edakit.Func(poissonYield))
	reg.Register("stat.yield.clustered", edakit.Func(clusteredYield))
	reg.Register("stat.yield.redundancy", edakit.Func(redundancy))
	reg.Register("stat.yield.design_centering", edakit.Func(designCentering))
	// reliability
	reg.Register("stat.rel.weibull", edakit.Func(weibull))
	reg.Register("stat.rel.hazard_fit", edakit.Func(hazardFIT))
	reg.Register("stat.rel.acceleration", edakit.Func(acceleration))
	reg.Register("stat.rel.weibull_fit", edakit.Func(weibullFit))
	reg.Register("stat.rel.system", edakit.Func(systemReliability))
	// frontier
	reg.Register("stat.frontier.variance_aware_surrogate", edakit.Func(varianceAwareSurrogate))
	reg.Register("stat.frontier.mc_budget", edakit.Func(mcBudget))
	reg.Register("stat.frontier.rare_event_skill", &edakit.Bank{Items: rareEventSkillItems})
	reg.Register("stat.frontier.golden_statistical_tests", edakit.Func(goldenStatTests))
	reg.Register("stat.frontier.yield_agent_loop", &edakit.Bank{Items: yieldAgentItems})
}

// kp resolves the knowledge point to serve: the requested one, or a uniform draw when unguided.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// phi is the standard normal density.
func phi(x float64) float64 { return math.Exp(-x*x/2) / math.Sqrt(2*math.Pi) }

// cdf is the standard normal distribution function.
func cdf(x float64) float64 { return 0.5 * math.Erfc(-x/math.Sqrt2) }

// tail is the upper tail Q(x) = 1 - cdf(x), computed without cancellation.
func tail(x float64) float64 { return 0.5 * math.Erfc(x/math.Sqrt2) }

// tailInv returns k with tail(k) = p, for 0 < p < 0.5, by bisection.
func tailInv(p float64) float64 {
	lo, hi := 0.0, 40.0
	for i := 0; i < 200; i++ {
		mid := (lo + hi) / 2
		if tail(mid) > p {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

// sigPlaces returns the decimal places that show about three significant figures of x.
func sigPlaces(x float64) int {
	if x <= 0 {
		return 2
	}
	p := 2 - int(math.Floor(math.Log10(x)))
	if p < 0 {
		p = 0
	}
	return p
}

// f formats a drawn operand with up to the given places, trailing zeros trimmed.
func f(x float64, places int) string { return edakit.Num(x, places) }

// sciT wraps sci in math delimiters for use in prose.
func sciT(x float64) string { return "$" + sci(x) + "$" }

// roundTo states the rounding a stem asks for.
func roundTo(places int) string {
	switch places {
	case 0:
		return "Round to the nearest integer."
	case 1:
		return "Round to 1 decimal place."
	}
	return edakit.Sprintf("Round to %d decimal places.", places)
}

// pn renders a number for use inside a product, parenthesizing negatives.
func pn(x float64, places int) string {
	if x < 0 {
		return "(" + f(x, places) + ")"
	}
	return f(x, places)
}

// sci formats x as bare LaTeX, a.bc \times 10^{n}, for use inside math.
func sci(x float64) string {
	if x == 0 {
		return "0"
	}
	e := int(math.Floor(math.Log10(math.Abs(x))))
	m := x / math.Pow(10, float64(e))
	if math.Abs(m) >= 9.995 {
		m /= 10
		e++
	}
	return edakit.Sprintf("%s \\times 10^{%d}", strconv.FormatFloat(m, 'f', 2, 64), e)
}

// numChoice builds a computed multiple-choice problem whose options are numbers formatted to
// places, followed by unit. Distractors that collide with the answer or each other after
// formatting are replaced by nudged values so that at least three distinct wrong options remain.
func numChoice(q string, correct float64, places int, unit, why string, wrong ...float64) generator.Problem {
	fmtv := func(x float64) string {
		s := edakit.Num(x, places)
		if unit != "" {
			s += " " + unit
		}
		return s
	}
	ans := fmtv(correct)
	seen := map[string]bool{ans: true}
	var ds []string
	add := func(x float64) {
		s := fmtv(x)
		if !seen[s] && !math.IsNaN(x) && !math.IsInf(x, 0) {
			seen[s] = true
			ds = append(ds, s)
		}
	}
	for _, w := range wrong {
		add(w)
	}
	step := math.Pow(10, -float64(places))
	for k := 1; len(ds) < 3 && k < 50; k++ {
		add(correct * (1 + 0.25*float64(k)))
		add(correct + float64(k)*step*5)
	}
	return edakit.Choice(q, ans, why, ds...)
}
