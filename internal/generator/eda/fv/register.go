package fv

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the formal-verification generators to reg.
func Register(reg *generator.Registry) {
	// SAT
	reg.Register("fv.sat.cnf_tseitin", edakit.Func(tseitin))
	reg.Register("fv.sat.unit_propagation", edakit.Func(unitProp))
	reg.Register("fv.sat.dpll", edakit.Func(dpll))
	reg.Register("fv.sat.implication_graph", edakit.Func(implGraph))
	reg.Register("fv.sat.cdcl", &edakit.Bank{Items: cdclItems})
	reg.Register("fv.sat.first_uip", edakit.Func(firstUIP))
	reg.Register("fv.sat.backjump", edakit.Func(backjump))
	reg.Register("fv.sat.watched_literals", edakit.Func(watched))
	reg.Register("fv.sat.vsids", edakit.Func(vsids))
	reg.Register("fv.sat.restarts", edakit.Func(restarts))
	reg.Register("fv.sat.preprocessing", edakit.Func(preprocessing))
	// BDD
	reg.Register("fv.bdd.reduction_rules", edakit.Func(reductionRules))
	reg.Register("fv.bdd.robdd", edakit.Func(robdd))
	reg.Register("fv.bdd.variable_order", edakit.Func(variableOrder))
	reg.Register("fv.bdd.ite_apply", edakit.Func(iteApply))
	reg.Register("fv.bdd.sat_count", edakit.Func(satCount))
	reg.Register("fv.bdd.multiplier_blowup", &edakit.Bank{Items: multiplierBlowupItems})
	// equivalence
	reg.Register("fv.equiv.miter", edakit.Func(miter))
	reg.Register("fv.equiv.structural_hashing", edakit.Func(strash))
	reg.Register("fv.equiv.cec", &edakit.Bank{Items: cecItems})
	reg.Register("fv.equiv.sat_sweeping", edakit.Func(sweeping))
	reg.Register("fv.equiv.cut_points", &edakit.Bank{Items: cutPointsItems})
	reg.Register("fv.equiv.sequential", edakit.Func(sequential))
	// model checking
	reg.Register("fv.mc.kripke", edakit.Func(kripke))
	reg.Register("fv.mc.reachability", edakit.Func(reachability))
	reg.Register("fv.mc.state_explosion", edakit.Func(stateExplosion))
	reg.Register("fv.mc.ltl", edakit.Func(ltlChoice))
	reg.Register("fv.mc.safety_liveness", &edakit.Bank{Items: safetyLivenessItems})
	reg.Register("fv.mc.ctl", edakit.Func(ctl))
	reg.Register("fv.mc.ctl_fixpoint", edakit.Func(ctlFixpoint))
	// bounded model checking
	reg.Register("fv.bmc.unrolling", edakit.Func(bmcUnrolling))
	reg.Register("fv.bmc.shortest_counterexample", edakit.Func(bmcDepth))
	reg.Register("fv.bmc.completeness_threshold", edakit.Func(completeness))
	reg.Register("fv.bmc.incremental", &edakit.Bank{Items: bmcIncrementalItems})
	// induction
	reg.Register("fv.induction.inductive_invariant", edakit.Func(inductiveInvariant))
	reg.Register("fv.induction.k_induction", edakit.Func(kInduction))
	reg.Register("fv.induction.strengthening", &edakit.Bank{Items: strengtheningItems})
	// IC3 and abstraction
	reg.Register("fv.ic3.relative_induction", edakit.Func(relativeInduction))
	reg.Register("fv.ic3.pdr", &edakit.Bank{Items: pdrItems})
	genBank := &edakit.Bank{Items: generalizationItems}
	reg.Register("fv.ic3.generalization", &edakit.Mix{Parts: []generator.Generator{genBank, edakit.Func(generalizationCompute), genBank}})
	cegarBank := &edakit.Bank{Items: cegarItems}
	reg.Register("fv.ic3.cegar", &edakit.Mix{Parts: []generator.Generator{cegarBank, edakit.Func(cegarCompute), cegarBank}})
	// SMT and program reasoning
	reg.Register("fv.smt.bitvectors", edakit.Func(bitvectors))
	reg.Register("fv.smt.euf", edakit.Func(euf))
	reg.Register("fv.smt.arrays", &edakit.Mix{Parts: []generator.Generator{edakit.Func(arraysCompute), edakit.Func(arraysCompute), &edakit.Bank{Items: arraysItems}}})
	reg.Register("fv.smt.dpll_t", &edakit.Bank{Items: dplltItems})
	reg.Register("fv.smt.theorem_proving", &edakit.Bank{Items: theoremProvingItems})
	reg.Register("fv.smt.symbolic_execution", edakit.Func(symbolicExecution))
	reg.Register("fv.smt.abstract_interpretation", edakit.Func(absInterp))
	// assertions
	reg.Register("fv.sva.assertions", edakit.Func(svaAssertions))
	reg.Register("fv.sva.sequences", edakit.Func(svaSequences))
	reg.Register("fv.sva.assume_assert_cover", &edakit.Bank{Items: assumeAssertCoverItems})
	reg.Register("fv.sva.formal_apps", &edakit.Bank{Items: formalAppsItems})
	reg.Register("fv.sva.symbiyosys", &edakit.Bank{Items: symbiyosysItems})
	// frontier
	reg.Register("fv.frontier.agent_invariants", edakit.Func(houdini))
	reg.Register("fv.frontier.llm_sva_checking", &edakit.Bank{Items: llmSvaItems})
	reg.Register("fv.frontier.proof_guided_repair", &edakit.Bank{Items: repairItems})
	reg.Register("fv.frontier.formal_reward", &edakit.Bank{Items: formalRewardItems})
	reg.Register("fv.frontier.engine_seams", &edakit.Bank{Items: engineSeamsItems})
}
