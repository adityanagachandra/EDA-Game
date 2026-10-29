// Package place generates problems for the floorplanning_and_placement domain: the placement
// region, LEF and DEF, wirelength models, floorplan representations, macro and pin planning,
// global placement (min-cut, quadratic, force-directed, nonlinear and electrostatic),
// legalization, detailed placement, timing- and routability-driven placement, and the
// agent-first frontier of placer engines. Files are split by subdomain.
package place

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds every floorplanning_and_placement generator to reg.
func Register(reg *generator.Registry) {
	// overview and formats
	reg.Register("place.overview.flow", &edakit.Bank{Items: flowItems})
	reg.Register("place.overview.rows_sites", edakit.Func(rowsSites))
	reg.Register("place.overview.utilization", edakit.Func(utilization))
	reg.Register("place.format.lef_def", &edakit.Bank{Items: lefDefItems})
	reg.Register("place.format.orientation", edakit.Func(orientation))

	// wirelength models
	reg.Register("place.wl.hpwl", edakit.Func(hpwlGen))
	reg.Register("place.wl.incremental", edakit.Func(incremental))
	reg.Register("place.wl.clique_star", edakit.Func(cliqueStar))
	reg.Register("place.wl.b2b", edakit.Func(b2b))
	reg.Register("place.wl.steiner", edakit.Func(steiner))

	// floorplanning
	reg.Register("place.floorplan.whitespace", edakit.Func(whitespace))
	reg.Register("place.floorplan.slicing", edakit.Func(slicing))
	reg.Register("place.floorplan.polish", edakit.Func(polish))
	reg.Register("place.floorplan.polish_moves", edakit.Func(polishMoves))
	reg.Register("place.floorplan.shape_curves", edakit.Func(shapeCurves))
	reg.Register("place.floorplan.sequence_pair", edakit.Func(sequencePair))
	reg.Register("place.floorplan.bstar", edakit.Func(bstar))
	reg.Register("place.floorplan.annealing", edakit.Func(annealing))

	// macros, pins and power
	macroBank := &edakit.Bank{Items: macroItems}
	reg.Register("place.macro.placement", &edakit.Mix{Parts: []generator.Generator{macroBank, edakit.Func(macroChannel), macroBank}})
	reg.Register("place.macro.pin_assignment", edakit.Func(pinAssignment))
	reg.Register("place.macro.power_planning", edakit.Func(powerPlanning))

	// global placement
	reg.Register("place.global.formulation", &edakit.Bank{Items: formulationItems})
	reg.Register("place.global.mincut", edakit.Func(mincut))
	reg.Register("place.global.quadratic", edakit.Func(quadratic))
	reg.Register("place.global.laplacian_system", edakit.Func(laplacianSystem))
	reg.Register("place.global.cg_solve", edakit.Func(cgSolve))
	reg.Register("place.global.spreading", edakit.Func(spreading))
	reg.Register("place.global.force_directed", edakit.Func(forceDirected))
	reg.Register("place.global.density", edakit.Func(density))
	reg.Register("place.global.smooth_wirelength", edakit.Func(smoothWL))
	reg.Register("place.global.lse_gradient", edakit.Func(lseGradient))
	reg.Register("place.global.wa_wirelength", edakit.Func(waWirelength))
	reg.Register("place.global.nonlinear", edakit.Func(nonlinear))
	reg.Register("place.global.electrostatic", edakit.Func(electrostatic))
	reg.Register("place.global.poisson_fft", edakit.Func(poissonFFT))
	reg.Register("place.global.nesterov", edakit.Func(nesterov))
	reg.Register("place.global.dreamplace", &edakit.Bank{Items: dreamplaceItems})

	// legalization
	reg.Register("place.legal.displacement", edakit.Func(displacement))
	reg.Register("place.legal.tetris", edakit.Func(tetris))
	reg.Register("place.legal.abacus", edakit.Func(abacus))
	reg.Register("place.legal.rules", &edakit.Bank{Items: legalRulesItems})

	// detailed placement
	reg.Register("place.detail.swap", edakit.Func(swap))
	reg.Register("place.detail.optimal_region", edakit.Func(optimalRegion))
	reg.Register("place.detail.sliding_window", edakit.Func(slidingWindow))
	reg.Register("place.detail.ism", edakit.Func(ism))

	// timing- and routability-driven placement
	reg.Register("place.driven.net_delay", edakit.Func(netDelay))
	reg.Register("place.driven.timing", edakit.Func(timingWeights))
	reg.Register("place.driven.routability", edakit.Func(rudy))
	reg.Register("place.driven.inflation", edakit.Func(inflation))
	reg.Register("place.driven.buffering", edakit.Func(buffering))

	// frontier
	reg.Register("place.frontier.core_vs_surface", &edakit.Bank{Items: coreSurfaceItems})
	ruleBank := &edakit.Bank{Items: rulePluginItems}
	reg.Register("place.frontier.rule_plugin", &edakit.Mix{Parts: []generator.Generator{ruleBank, edakit.Func(rulePluginCalc), ruleBank}})
	reg.Register("place.frontier.objective_weights", edakit.Func(objectiveWeights))
	reg.Register("place.frontier.qor_differential", edakit.Func(qorDifferential))
	fastBank := &edakit.Bank{Items: fastLoopItems}
	reg.Register("place.frontier.fast_loop", &edakit.Mix{Parts: []generator.Generator{fastBank, edakit.Func(fastLoopCalc), fastBank}})
}
