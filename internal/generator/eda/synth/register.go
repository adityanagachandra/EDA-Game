// Package synth generates logic-synthesis problems: two-level minimization, multi-level
// algebraic optimization, don't cares, and-inverter graphs, FSM and sequential synthesis,
// datapath optimization with e-graphs, and the Yosys and ABC flow.
package synth

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the logic-synthesis generators to reg.
func Register(reg *generator.Registry) {
	bank := func(items []edakit.Item) generator.Generator { return &edakit.Bank{Items: items} }
	f := func(fn func(generator.GeneratorContext) generator.Problem) generator.Generator {
		return edakit.Func(fn)
	}

	reg.Register("synth.overview.flow", bank(flowItems))
	reg.Register("synth.overview.cost_metrics", f(costMetrics))
	reg.Register("synth.overview.netlist_abstraction", bank(netlistItems))

	reg.Register("synth.twolevel.cube_notation", f(cubeNotation))
	reg.Register("synth.twolevel.kmap", f(kmap))
	reg.Register("synth.twolevel.quine_mccluskey", f(quineMcCluskey))
	reg.Register("synth.twolevel.prime_table", f(primeTable))
	reg.Register("synth.twolevel.table_reduction", f(tableReduction))
	reg.Register("synth.twolevel.unate_covering", f(unateCovering))
	reg.Register("synth.twolevel.petrick", f(petrick))
	reg.Register("synth.twolevel.tautology", f(tautology))
	reg.Register("synth.twolevel.expand", f(expand))
	reg.Register("synth.twolevel.irredundant", f(irredundant))
	reg.Register("synth.twolevel.espresso", bank(espressoItems))

	reg.Register("synth.multilevel.boolean_network", f(booleanNetwork))
	reg.Register("synth.multilevel.algebraic_division", f(algebraicDivision))
	reg.Register("synth.multilevel.kernels", f(kernels))
	reg.Register("synth.multilevel.factoring", f(factoring))
	reg.Register("synth.multilevel.extraction", f(extraction))
	reg.Register("synth.multilevel.cube_extraction", f(cubeExtraction))
	reg.Register("synth.multilevel.elimination", f(elimination))
	reg.Register("synth.multilevel.boolean_division", bank(booleanDivisionItems))

	reg.Register("synth.dc.dont_cares", bank(dontCareItems))
	reg.Register("synth.dc.sdc", f(sdc))
	reg.Register("synth.dc.odc", f(odc))
	reg.Register("synth.dc.resubstitution", f(resubstitution))
	reg.Register("synth.dc.compatibility", bank(compatibilityItems))

	reg.Register("synth.aig.structure", f(aigStructure))
	reg.Register("synth.aig.strash", f(strash))
	reg.Register("synth.aig.cuts", f(aigCuts))
	reg.Register("synth.aig.mffc", f(mffcDrill))
	reg.Register("synth.aig.rewriting", f(rewriting))
	reg.Register("synth.aig.refactoring", f(refactoring))
	reg.Register("synth.aig.balancing", f(balancing))
	reg.Register("synth.aig.fraig", f(fraig))
	reg.Register("synth.aig.tech_mapping", f(techMapping))
	reg.Register("synth.aig.lut_mapping", f(lutMapping))

	reg.Register("synth.fsm.state_minimization", f(stateMinimization))
	reg.Register("synth.fsm.equivalence", f(fsmEquivalence))
	reg.Register("synth.fsm.encoding", f(fsmEncoding))

	reg.Register("synth.seq.retiming", f(retiming))
	reg.Register("synth.seq.wd_matrices", f(wdMatrices))
	reg.Register("synth.seq.min_period", f(minPeriod))
	reg.Register("synth.seq.min_area", bank(minAreaItems))
	reg.Register("synth.seq.pipelining", bank(pipeliningItems))

	reg.Register("synth.arith.csd", f(csd))
	reg.Register("synth.arith.carry_save", f(carrySave))
	reg.Register("synth.egraph.equality_saturation", bank(eqSatItems))
	reg.Register("synth.egraph.congruence", f(congruence))
	reg.Register("synth.egraph.extraction", f(egExtraction))

	reg.Register("synth.tool.yosys_abc", bank(yosysAbcItems))
	reg.Register("synth.tool.abc_scripts", bank(abcScriptItems))
	reg.Register("synth.tool.cec", f(cec))

	registerFrontier(reg)
}
