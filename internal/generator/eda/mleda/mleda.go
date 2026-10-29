// Package mleda generates machine-learning-for-EDA problems: where learning helps and where it
// cannot replace algorithms and verification, graph neural networks on netlists, prediction
// tasks (congestion, timing, DRC hotspots, IR drop), reinforcement learning for placement and
// synthesis, surrogates and Bayesian optimization, differentiable optimization, learned routing
// and generative design, LLMs for RTL and their evaluation, and the agent-first frontier.
package mleda

import (
	"math"
	"strconv"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the ML-for-EDA generators to reg.
func Register(reg *generator.Registry) {
	// foundations
	reg.Register("mleda.found.where_ml_helps", &edakit.Bank{Items: whereMLHelpsItems})
	reg.Register("mleda.found.labels_from_tools", edakit.Func(labelsFromTools))
	// graph neural networks
	reg.Register("mleda.gnn.netlist_graph", edakit.Func(netlistGraph))
	reg.Register("mleda.gnn.message_passing", edakit.Func(messagePassing))
	reg.Register("mleda.gnn.netlist_gnn", edakit.Func(netlistGNN))
	reg.Register("mleda.gnn.hetero_graphs", edakit.Func(heteroGraphs))
	reg.Register("mleda.gnn.fanout_scale", edakit.Func(fanoutScale))
	// prediction tasks
	reg.Register("mleda.predict.congestion", edakit.Func(congestion))
	reg.Register("mleda.predict.timing", edakit.Func(timingPrediction))
	reg.Register("mleda.predict.drc_hotspot", edakit.Func(drcHotspot))
	reg.Register("mleda.predict.ir_drop", edakit.Func(irDrop))
	reg.Register("mleda.predict.datasets", &edakit.Bank{Items: datasetItems})
	// reinforcement learning
	reg.Register("mleda.rl.mdp_formulation", edakit.Func(mdpFormulation))
	reg.Register("mleda.rl.placement", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: rlPlacementItems}, edakit.Func(proxyCostChoice), &edakit.Bank{Items: rlPlacementItems}}})
	reg.Register("mleda.rl.placement_debate", &edakit.Bank{Items: placementDebateItems})
	reg.Register("mleda.rl.synthesis_recipes", edakit.Func(synthesisRecipes))
	// surrogates and Bayesian optimization
	reg.Register("mleda.surrogate.speedup_accounting", edakit.Func(speedupAccounting))
	reg.Register("mleda.surrogate.active_learning", edakit.Func(activeLearning))
	reg.Register("mleda.surrogate.gp_posterior", edakit.Func(gpPosterior))
	reg.Register("mleda.surrogate.acquisition", edakit.Func(acquisition))
	reg.Register("mleda.surrogate.analog_sizing", edakit.Func(analogSizing))
	reg.Register("mleda.surrogate.flow_tuning", edakit.Func(flowTuning))
	// differentiable optimization
	reg.Register("mleda.diff.differentiable_placement", edakit.Func(differentiablePlacement))
	reg.Register("mleda.diff.smooth_max", edakit.Func(smoothMax))
	reg.Register("mleda.diff.differentiable_timing", edakit.Func(differentiableTiming))
	reg.Register("mleda.diff.relaxation_limits", &edakit.Bank{Items: relaxationItems})
	// routing and generative design
	reg.Register("mleda.gen.learned_routing", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: learnedRoutingItems}, edakit.Func(weightedAStarChoice), &edakit.Bank{Items: learnedRoutingItems}}})
	reg.Register("mleda.gen.generative_design", edakit.Func(generativeDesign))
	// LLMs
	reg.Register("mleda.llm.rtl_generation", &edakit.Bank{Items: rtlGenerationItems})
	reg.Register("mleda.llm.repair_loops", edakit.Func(repairLoops))
	reg.Register("mleda.llm.assertion_generation", &edakit.Bank{Items: assertionItems})
	reg.Register("mleda.llm.script_generation", &edakit.Bank{Items: scriptItems})
	// evaluation
	reg.Register("mleda.eval.pass_at_k", edakit.Func(passAtK))
	reg.Register("mleda.eval.functional_correctness", edakit.Func(functionalCorrectness))
	reg.Register("mleda.eval.benchmarks", &edakit.Bank{Items: benchmarkItems})
	reg.Register("mleda.eval.contamination", &edakit.Bank{Items: contaminationItems})
	reg.Register("mleda.eval.distribution_shift", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: shiftItems}, edakit.Func(designSplitChoice), &edakit.Bank{Items: shiftItems}}})
	reg.Register("mleda.eval.ranking_metrics", edakit.Func(rankingMetrics))
	// frontier
	reg.Register("mleda.frontier.learned_fast_mode", edakit.Func(learnedFastMode))
	reg.Register("mleda.frontier.verification_gate", &edakit.Bank{Items: verificationGateItems})
	reg.Register("mleda.frontier.on_site_learning", &edakit.Bank{Items: onSiteItems})
	reg.Register("mleda.frontier.tool_feedback_training", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: toolFeedbackItems}, edakit.Func(toolRewardChoice), &edakit.Bank{Items: toolFeedbackItems}}})
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

// f formats a drawn operand with up to the given places, trailing zeros trimmed.
func f(x float64, places int) string { return edakit.Num(x, places) }

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

// ints formats a list of integers joined by commas.
func ints(xs []int) string {
	s := ""
	for i, x := range xs {
		if i > 0 {
			s += ", "
		}
		s += strconv.Itoa(x)
	}
	return s
}

// floats formats a list of values joined by commas.
func floats(xs []float64, places int) string {
	s := ""
	for i, x := range xs {
		if i > 0 {
			s += ", "
		}
		s += f(x, places)
	}
	return s
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

// binom returns the binomial coefficient C(n, k) as a float.
func binom(n, k int) float64 {
	if k < 0 || k > n {
		return 0
	}
	r := 1.0
	for i := 1; i <= k; i++ {
		r = r * float64(n-k+i) / float64(i)
	}
	return r
}
