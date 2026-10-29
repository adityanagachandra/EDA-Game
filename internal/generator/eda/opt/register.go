// Package opt generates optimization-for-EDA problems: formulation, linear and integer
// programming, convex and nonlinear optimization, first- and second-order methods, Lagrangian
// and penalty methods, metaheuristics, dynamic and constraint programming, multi-objective and
// Bayesian optimization, and the agent-first frontier of optimization engines.
package opt

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the optimization-for-EDA generators to reg.
func Register(reg *generator.Registry) {
	// formulate
	reg.Register("opt.formulate.variables_objective", &edakit.Bank{Items: variablesObjectiveItems})
	reg.Register("opt.formulate.problem_classes", &edakit.Mix{Parts: []generator.Generator{edakit.Func(searchSpaceSize), &edakit.Bank{Items: problemClassItems}, &edakit.Bank{Items: problemClassItems}}})
	reg.Register("opt.formulate.hardness", &edakit.Mix{Parts: []generator.Generator{edakit.Func(enumerationTime), &edakit.Bank{Items: hardnessItems}, &edakit.Bank{Items: hardnessItems}}})
	reg.Register("opt.formulate.strategy", &edakit.Bank{Items: strategyItems})
	// lp
	reg.Register("opt.lp.simplex_geometry", edakit.Func(lpGeometry))
	reg.Register("opt.lp.simplex_pivot", edakit.Func(simplexPivot))
	reg.Register("opt.lp.duality", edakit.Func(lpDuality))
	reg.Register("opt.lp.shadow_prices", edakit.Func(shadowPrices))
	reg.Register("opt.lp.difference_constraints", edakit.Func(differenceConstraints))
	// ip
	reg.Register("opt.ip.lp_relaxation", edakit.Func(lpRelaxation))
	reg.Register("opt.ip.branch_and_bound", edakit.Func(branchAndBound))
	reg.Register("opt.ip.big_m", edakit.Func(bigM))
	reg.Register("opt.ip.ilp_scheduling", edakit.Func(ilpScheduling))
	reg.Register("opt.ip.global_routing", edakit.Func(globalRouting))
	reg.Register("opt.ip.cutting_planes", edakit.Func(cuttingPlanes))
	reg.Register("opt.ip.benders", &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: bendersItems}, edakit.Func(bendersCut), &edakit.Bank{Items: bendersItems}}})
	// convex
	reg.Register("opt.convex.convexity_test", &edakit.Mix{Parts: []generator.Generator{edakit.Func(convexityClass), edakit.Func(hessianDet), &edakit.Bank{Items: convexityItems}}})
	reg.Register("opt.convex.kkt", edakit.Func(kkt))
	reg.Register("opt.convex.quadratic_placement", edakit.Func(quadraticPlacement))
	reg.Register("opt.convex.geometric_programming", edakit.Func(geometricProgramming))
	reg.Register("opt.convex.nonconvex", &edakit.Mix{Parts: []generator.Generator{edakit.Func(basinOfAttraction), &edakit.Bank{Items: nonconvexItems}, &edakit.Bank{Items: nonconvexItems}}})
	// gradient
	reg.Register("opt.gradient.descent", edakit.Func(gradientDescent))
	reg.Register("opt.gradient.step_size", edakit.Func(stepSize))
	reg.Register("opt.gradient.nesterov", edakit.Func(nesterov))
	reg.Register("opt.gradient.conjugate_gradient", edakit.Func(conjugateGradient))
	reg.Register("opt.gradient.newton", edakit.Func(newtonMin))
	reg.Register("opt.gradient.quasi_newton", edakit.Func(quasiNewton))
	reg.Register("opt.gradient.coordinate_descent", edakit.Func(coordinateDescent))
	reg.Register("opt.gradient.smooth_wirelength", edakit.Func(smoothWirelength))
	// lagrange
	reg.Register("opt.lagrange.dual_function", edakit.Func(dualFunction))
	reg.Register("opt.lagrange.relaxation", edakit.Func(lagrangianRelaxation))
	reg.Register("opt.lagrange.subgradient", edakit.Func(subgradient))
	reg.Register("opt.lagrange.penalty", edakit.Func(penalty))
	reg.Register("opt.lagrange.augmented", edakit.Func(augmentedLagrangian))
	reg.Register("opt.lagrange.admm", edakit.Func(admm))
	// meta
	reg.Register("opt.meta.local_search", edakit.Func(localSearch))
	reg.Register("opt.meta.simulated_annealing", edakit.Func(annealing))
	reg.Register("opt.meta.cooling_schedule", edakit.Func(cooling))
	reg.Register("opt.meta.polish_expression", edakit.Func(polishExpression))
	reg.Register("opt.meta.genetic", edakit.Func(genetic))
	// dp
	reg.Register("opt.dp.principle", edakit.Func(dpPrinciple))
	reg.Register("opt.dp.tree_mapping", edakit.Func(treeMapping))
	reg.Register("opt.dp.buffer_insertion", edakit.Func(bufferInsertion))
	reg.Register("opt.dp.constraint_propagation", edakit.Func(constraintPropagation))
	// multi
	reg.Register("opt.multi.pareto", edakit.Func(pareto))
	reg.Register("opt.multi.scalarization", edakit.Func(scalarization))
	reg.Register("opt.multi.robust_stochastic", edakit.Func(robustStochastic))
	// bayes
	reg.Register("opt.bayes.gaussian_process", edakit.Func(gaussianProcess))
	reg.Register("opt.bayes.bayesian_optimization", edakit.Func(bayesOpt))
	reg.Register("opt.bayes.eda_practice", &edakit.Bank{Items: boPracticeItems})
	// frontier
	reg.Register("opt.frontier.objective_weights", &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: weightTierItems}, edakit.Func(weightSwitch), &edakit.Bank{Items: weightTierItems}}})
	reg.Register("opt.frontier.move_set_policy", &edakit.Bank{Items: moveSetItems})
	reg.Register("opt.frontier.reward_design", &edakit.Mix{Parts: []generator.Generator{edakit.Func(rewardPick), &edakit.Bank{Items: rewardItems}, &edakit.Bank{Items: rewardItems}}})
	reg.Register("opt.frontier.dse_agent", &edakit.Mix{Parts: []generator.Generator{edakit.Func(successiveHalving), &edakit.Bank{Items: dseItems}, &edakit.Bank{Items: dseItems}}})
	reg.Register("opt.frontier.optimizer_verification", &edakit.Mix{Parts: []generator.Generator{edakit.Func(geomeanQoR), &edakit.Bank{Items: verifyItems}, &edakit.Bank{Items: verifyItems}}})
}

// kp returns the requested knowledge point, or a uniform one for an unguided request.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

// num builds a numeric problem with a single fact map built from alternating key, value pairs.
func num(q, ans, expl string, kv ...string) generator.Problem {
	facts := map[string]string{}
	for i := 0; i+1 < len(kv); i += 2 {
		facts[kv[i]] = kv[i+1]
	}
	return edakit.Numeric(q, ans, expl, facts)
}

// f formats x at the given number of decimal places.
func f(x float64, places int) string { return edakit.Num(x, places) }

// roundNote is the rounding sentence appended to a stem.
func roundNote(places int) string {
	switch places {
	case 0:
		return " Give an integer."
	case 1:
		return " Round to 1 decimal place."
	default:
		return edakit.Sprintf(" Round to %d decimal places.", places)
	}
}

// numChoices returns distinct formatted distractors for a numeric multiple-choice question,
// dropping any that format the same as the correct value or as each other.
func numChoices(correct float64, places int, cands ...float64) []string {
	seen := map[string]bool{f(correct, places): true}
	var out []string
	for _, c := range cands {
		s := f(c, places)
		if !seen[s] && !math.IsNaN(c) && !math.IsInf(c, 0) {
			seen[s] = true
			out = append(out, s)
		}
	}
	for d := 1.0; len(out) < 3; d++ {
		for _, c := range []float64{correct + d, correct - d, correct * (1 + d/2)} {
			s := f(c, places)
			if !seen[s] && len(out) < 4 {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return out
}

// ints joins integers with a separator.
func ints(xs []int, sep string) string {
	parts := make([]string, len(xs))
	for i, x := range xs {
		parts[i] = edakit.IntStr(x)
	}
	return strings.Join(parts, sep)
}

// distinctInts draws n distinct integers from [lo, hi].
func distinctInts(n, lo, hi int) []int {
	perm := edakit.Shuffle(rangeInts(lo, hi))
	return perm[:n]
}

func rangeInts(lo, hi int) []int {
	out := make([]int, 0, hi-lo+1)
	for i := lo; i <= hi; i++ {
		out = append(out, i)
	}
	return out
}

// phi and bigPhi are the standard normal density and distribution function.
func phi(z float64) float64    { return math.Exp(-z*z/2) / math.Sqrt(2*math.Pi) }
func bigPhi(z float64) float64 { return 0.5 * (1 + math.Erf(z/math.Sqrt2)) }
