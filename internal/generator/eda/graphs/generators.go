// Package graphs generates EDA graph-algorithm problems: netlist hypergraphs and net models,
// DAG passes, shortest paths, partitioning, spectral methods, flows and matchings, Steiner
// trees and tours, coloring and planarity, and the agent-first frontier of graph engines.
package graphs

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the graph-algorithm generators to reg.
func Register(reg *generator.Registry) {
	// model
	reg.Register("graphs.model.hypergraph", edakit.Func(hypergraphGen))
	reg.Register("graphs.model.pin_graph", edakit.Func(pinGraphGen))
	reg.Register("graphs.model.net_models", edakit.Func(netModelsGen))
	reg.Register("graphs.model.csr", edakit.Func(csrGen))
	// dag
	reg.Register("graphs.tree.rooted_traversal", edakit.Func(rootedTraversal))
	reg.Register("graphs.dag.topo_sort", edakit.Func(topoSortGen))
	reg.Register("graphs.dag.levelize", edakit.Func(levelizeGen))
	reg.Register("graphs.dag.scc", edakit.Func(sccGen))
	reg.Register("graphs.dag.longest_path", edakit.Func(longestPathGen))
	reg.Register("graphs.dag.reverse_pass", edakit.Func(reversePassGen))
	reg.Register("graphs.dag.path_count", edakit.Func(pathCountGen))
	reg.Register("graphs.dag.incremental", edakit.Func(incrementalGen))
	// path
	reg.Register("graphs.path.grid_graph", edakit.Func(gridGraphGen))
	reg.Register("graphs.path.bfs_maze", edakit.Func(bfsMazeGen))
	reg.Register("graphs.path.dijkstra", edakit.Func(dijkstraGen))
	reg.Register("graphs.path.bellman_ford", edakit.Func(bellmanFordGen))
	reg.Register("graphs.path.difference_constraints", edakit.Func(differenceConstraintsGen))
	reg.Register("graphs.path.astar", edakit.Func(astarGen))
	// partition
	reg.Register("graphs.partition.cut_size", edakit.Func(cutSizeGen))
	reg.Register("graphs.partition.clique_distortion", edakit.Func(cliqueDistortionGen))
	reg.Register("graphs.partition.balance", edakit.Func(balanceGen))
	reg.Register("graphs.partition.kl_gain", edakit.Func(klGainGen))
	reg.Register("graphs.partition.fm_gain", edakit.Func(fmGainGen))
	reg.Register("graphs.partition.fm_bucket", edakit.Func(fmBucketGen))
	fmBank := &edakit.Bank{Items: fmItems}
	reg.Register("graphs.partition.fm", &edakit.Mix{Parts: []generator.Generator{edakit.Func(fmPass), fmBank, fmBank}})
	reg.Register("graphs.partition.coarsen_matching", edakit.Func(coarsenMatchingGen))
	mlBank := &edakit.Bank{Items: multilevelItems}
	reg.Register("graphs.partition.multilevel", &edakit.Mix{Parts: []generator.Generator{edakit.Func(multilevelLevels), mlBank, mlBank}})
	// spectral
	reg.Register("graphs.spectral.laplacian", edakit.Func(laplacianGen))
	reg.Register("graphs.spectral.cut_quadratic", edakit.Func(cutQuadraticGen))
	reg.Register("graphs.spectral.ratio_cut", edakit.Func(ratioCutGen))
	reg.Register("graphs.spectral.fiedler", edakit.Func(fiedlerGen))
	// flow and matching
	reg.Register("graphs.flow.max_flow_min_cut", edakit.Func(maxFlowGen))
	reg.Register("graphs.flow.residual", edakit.Func(residualGen))
	reg.Register("graphs.flow.hypergraph_mincut", edakit.Func(hyperMinCutGen))
	reg.Register("graphs.flow.min_cost", edakit.Func(minCostGen))
	reg.Register("graphs.match.bipartite", edakit.Func(bipartiteGen))
	reg.Register("graphs.match.assignment", edakit.Func(assignmentGen))
	// steiner and tours
	reg.Register("graphs.steiner.mst", edakit.Func(mstGen))
	reg.Register("graphs.steiner.hanan", edakit.Func(hananGen))
	reg.Register("graphs.steiner.rsmt", edakit.Func(rsmtGen))
	reg.Register("graphs.steiner.ratio", edakit.Func(ratioGen))
	reg.Register("graphs.steiner.one_steiner", edakit.Func(oneSteinerGen))
	reg.Register("graphs.steiner.scan_tsp", edakit.Func(scanTSPGen))
	// coloring and planarity
	reg.Register("graphs.color.coloring", edakit.Func(coloringGen))
	reg.Register("graphs.color.interval", edakit.Func(intervalGen))
	reg.Register("graphs.color.double_patterning", edakit.Func(doublePatterningGen))
	plBank := &edakit.Bank{Items: planarityItems}
	reg.Register("graphs.embed.planarity", &edakit.Mix{Parts: []generator.Generator{edakit.Func(planarityCheck), plBank, plBank}})
	// frontier
	reg.Register("graphs.frontier.engine_layers", &edakit.Bank{Items: engineLayerItems})
	nwBank := &edakit.Bank{Items: netWeightItems}
	reg.Register("graphs.frontier.net_weight_seam", &edakit.Mix{Parts: []generator.Generator{edakit.Func(netWeightCut), nwBank, nwBank}})
	dcBank := &edakit.Bank{Items: diffCutItems}
	reg.Register("graphs.frontier.diff_cut_tests", &edakit.Mix{Parts: []generator.Generator{edakit.Func(diffVerdict), dcBank, dcBank}})
	fmBankF := &edakit.Bank{Items: fastModeItems}
	reg.Register("graphs.frontier.fast_multilevel", &edakit.Mix{Parts: []generator.Generator{edakit.Func(fastDecision), fmBankF, fmBankF}})
	reg.Register("graphs.frontier.graph_rewards", &edakit.Bank{Items: rewardItems})
}
