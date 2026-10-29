package boolalg

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the Boolean-reasoning generators to reg.
func Register(reg *generator.Registry) {
	f := func(id string, fn func(generator.GeneratorContext) generator.Problem) {
		reg.Register(id, edakit.Func(fn))
	}
	bank := func(id string, items []edakit.Item) { reg.Register(id, &edakit.Bank{Items: items}) }

	f("boolalg.func.truth_table", truthTable)
	f("boolalg.func.cofactor", cofactor)
	f("boolalg.func.shannon", shannon)
	f("boolalg.func.boolean_difference", booleanDifference)
	f("boolalg.func.unateness", unateness)
	f("boolalg.func.symmetry", symmetry)
	f("boolalg.func.quantification", quantification)

	f("boolalg.decomp.support", support)
	f("boolalg.decomp.disjoint", disjointDecomp)

	f("boolalg.cover.cubes", cubes)
	f("boolalg.cover.sop_pos", sopPos)
	f("boolalg.cover.implicants", implicants)
	f("boolalg.cover.essential_primes", essentialPrimes)
	f("boolalg.cover.min_cover", minCoverGen)
	f("boolalg.cover.dont_cares", dontCares)
	f("boolalg.cover.tautology", tautology)

	f("boolalg.repr.bitparallel", bitParallel)
	f("boolalg.repr.hash_consing", hashConsing)
	f("boolalg.repr.decision_diagram", decisionDiagram)

	f("boolalg.order.boolean_lattice", booleanLattice)
	f("boolalg.order.monotone", monotone)
	f("boolalg.order.fixpoint", fixpoint)

	f("boolalg.cnf.clauses", clausesGen)
	f("boolalg.cnf.tseitin", tseitin)
	f("boolalg.cnf.model_count", modelCount)
	f("boolalg.cnf.resolution", resolution)
	f("boolalg.cnf.csp", csp)
	bank("boolalg.cnf.tractable", tractableItems)

	f("boolalg.algo.analysis", analysis)
	f("boolalg.algo.amortized", amortized)
	f("boolalg.algo.randomized", randomized)

	bank("boolalg.complexity.p_np", pnpItems)
	redBank := &edakit.Bank{Items: reductionItems}
	reg.Register("boolalg.complexity.reduction", &edakit.Mix{Parts: []generator.Generator{redBank, edakit.Func(reductionCompute), redBank}})
	bank("boolalg.complexity.np_complete", npcItems)
	f("boolalg.complexity.sat_to_3sat", satTo3sat)
	f("boolalg.complexity.graph_reductions", graphReductions)
	f("boolalg.complexity.approximation", approximation)
	f("boolalg.complexity.parameterized", parameterized)

	bank("boolalg.frontier.reasoning_tools", reasoningToolItems)
	feBank := &edakit.Bank{Items: fastExactItems}
	reg.Register("boolalg.frontier.fast_exact_modes", &edakit.Mix{Parts: []generator.Generator{edakit.Func(fastExactCompute), feBank, feBank}})
	f("boolalg.frontier.equivalence_golden", equivalenceGolden)
	bank("boolalg.frontier.change_tiers", changeTierItems)
	bank("boolalg.frontier.tool_grounded_reward", rewardItems)
}
