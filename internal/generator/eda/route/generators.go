// Package route generates routing problems: the routing grid and its capacities, maze routing
// (Lee, Hadlock, Soukup, A*, line probing), pattern routing and net decomposition, global
// routing (congestion, integer and flow formulations, negotiation, rip-up and reroute), layer
// assignment, detailed routing, channel routing, timing-, crosstalk- and antenna-aware
// routing, and the agent-first frontier. Almost every concept draws a random instance and
// computes its answer; files are split by subdomain.
package route

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the routing generators to reg.
func Register(reg *generator.Registry) {
	reg.Register("route.grid.problem_formulation", &edakit.Bank{Items: problemFormulationItems})
	reg.Register("route.grid.global_vs_detailed", &edakit.Bank{Items: globalDetailedItems})
	reg.Register("route.grid.gcells_tracks", edakit.Func(gcellsTracks))
	reg.Register("route.grid.preferred_direction", edakit.Func(preferredDirection))
	reg.Register("route.grid.grid_graph", edakit.Func(gridGraph))

	reg.Register("route.maze.lee", edakit.Func(lee))
	reg.Register("route.maze.backtrace", edakit.Func(backtrace))
	reg.Register("route.maze.hadlock", edakit.Func(hadlock))
	reg.Register("route.maze.soukup", &edakit.Bank{Items: soukupItems})
	reg.Register("route.maze.astar_routing", edakit.Func(astarRouting))
	reg.Register("route.maze.bend_via_cost", edakit.Func(bendViaCost))
	reg.Register("route.maze.line_probing", edakit.Func(lineProbing))

	reg.Register("route.pattern.l_z", edakit.Func(lz))
	reg.Register("route.pattern.cost_eval", edakit.Func(costEval))
	reg.Register("route.pattern.monotonic", edakit.Func(monotonic))
	reg.Register("route.pattern.net_decomposition", edakit.Func(netDecomposition))
	reg.Register("route.pattern.flute", &edakit.Mix{Parts: []generator.Generator{edakit.Func(fluteSequences), edakit.Func(fluteEval), &edakit.Bank{Items: fluteItems}}})

	reg.Register("route.global.congestion", edakit.Func(congestion))
	reg.Register("route.global.probabilistic_estimate", edakit.Func(probabilisticEstimate))
	reg.Register("route.global.ilp_formulation", edakit.Func(ilpFormulation))
	reg.Register("route.global.multicommodity_flow", edakit.Func(multicommodityFlow))
	reg.Register("route.global.negotiation", edakit.Func(negotiation))
	reg.Register("route.global.history_cost", edakit.Func(historyCost))
	reg.Register("route.global.rip_up_reroute", edakit.Func(ripUpReroute))
	reg.Register("route.global.open_global_routers", &edakit.Bank{Items: openGlobalItems})

	reg.Register("route.layer.via_minimization", edakit.Func(viaMinimization))
	reg.Register("route.layer.assignment_dp", edakit.Func(assignmentDP))
	reg.Register("route.layer.timing_aware", edakit.Func(timingAware))
	reg.Register("route.layer.redundant_vias", edakit.Func(redundantVias))

	reg.Register("route.detail.track_assignment", edakit.Func(trackAssignment))
	reg.Register("route.detail.pin_access", edakit.Func(pinAccess))
	reg.Register("route.detail.drc_aware", edakit.Func(drcAware))
	reg.Register("route.detail.min_area_eol", edakit.Func(minAreaEOL))
	reg.Register("route.detail.drc_repair", edakit.Func(drcRepair))
	reg.Register("route.detail.multiple_patterning", edakit.Func(multiplePatterning))

	reg.Register("route.channel.density", edakit.Func(channelDensity))
	reg.Register("route.channel.hcg_vcg", edakit.Func(hcgVCG))
	reg.Register("route.channel.left_edge", edakit.Func(leftEdgeGen))
	reg.Register("route.channel.cyclic_vcg", edakit.Func(cyclicVCG))
	reg.Register("route.channel.dogleg", edakit.Func(dogleg))
	reg.Register("route.channel.switchbox", &edakit.Bank{Items: switchboxItems})

	reg.Register("route.driven.timing_driven", edakit.Func(timingDriven))
	reg.Register("route.driven.prim_dijkstra", edakit.Func(primDijkstraGen))
	reg.Register("route.driven.crosstalk", edakit.Func(crosstalk))
	reg.Register("route.driven.antenna", edakit.Func(antenna))

	reg.Register("route.frontier.core_vs_surface", &edakit.Bank{Items: coreSurfaceItems})
	reg.Register("route.frontier.cost_plugin", &edakit.Mix{Parts: []generator.Generator{edakit.Func(costPluginCalc), &edakit.Bank{Items: costPluginItems}, &edakit.Bank{Items: costPluginItems}}})
	reg.Register("route.frontier.layer_directives", &edakit.Bank{Items: layerDirectiveItems})
	reg.Register("route.frontier.drc_golden", edakit.Func(drcGolden))
	reg.Register("route.frontier.agentic_drc_repair", &edakit.Bank{Items: agenticRepairItems})
	reg.Register("route.frontier.fast_congestion", edakit.Func(fastCongestion))
}
