package fab

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the fabrication and DTCO generators to reg.
func Register(reg *generator.Registry) {
	bank := func(items []edakit.Item) *edakit.Bank { return &edakit.Bank{Items: items} }

	// Unit processes.
	reg.Register("fab.process.flow_overview", bank(flowOverviewItems))
	reg.Register("fab.process.wafers", edakit.Func(wafers))
	reg.Register("fab.process.oxidation", edakit.Func(oxidation))
	reg.Register("fab.process.diffusion", edakit.Func(diffusion))
	reg.Register("fab.process.implant", edakit.Func(implant))
	reg.Register("fab.process.epitaxy", bank3(epitaxyItems, edakit.Func(epitaxyCompute)))
	reg.Register("fab.process.deposition", bank3(depositionItems, edakit.Func(depositionCompute)))
	reg.Register("fab.process.etch", edakit.Func(etch))
	reg.Register("fab.process.cmp", edakit.Func(cmpGen))

	// Lithography and patterning.
	reg.Register("fab.litho.resolution", edakit.Func(resolution))
	reg.Register("fab.litho.depth_of_focus", edakit.Func(depthOfFocus))
	reg.Register("fab.litho.photomask", edakit.Func(photomask))
	reg.Register("fab.litho.ret", bank(retItems))
	reg.Register("fab.litho.opc", bank(opcItems))
	reg.Register("fab.litho.model_based_opc", edakit.Func(modelOPC))
	reg.Register("fab.litho.euv", edakit.Func(euv))
	reg.Register("fab.litho.multiple_patterning", edakit.Func(multiPatterning))
	reg.Register("fab.litho.decomposition", edakit.Func(decomposition))
	reg.Register("fab.litho.sadp", edakit.Func(sadp))
	reg.Register("fab.litho.mask_data_prep", bank3(mdpItems, edakit.Func(mdpCompute)))

	// Process integration and interconnect.
	reg.Register("fab.beol.feol_mol_beol", bank(feolMolBeolItems))
	reg.Register("fab.beol.interconnect_stack", edakit.Func(interconnectStack))
	reg.Register("fab.beol.damascene", edakit.Func(damascene))
	reg.Register("fab.beol.wire_rc", edakit.Func(wireRC))
	reg.Register("fab.beol.low_k", edakit.Func(lowK))
	reg.Register("fab.beol.vias_contact", edakit.Func(viasContact))

	// Process design kits.
	reg.Register("fab.pdk.contents", bank(pdkContentsItems))
	reg.Register("fab.pdk.design_rules", edakit.Func(designRules))
	reg.Register("fab.pdk.drm", bank(drmItems))
	reg.Register("fab.pdk.device_models", bank(deviceModelItems))

	// Yield.
	reg.Register("fab.yield.dies_per_wafer", edakit.Func(diesPerWafer))
	reg.Register("fab.yield.defect_density", edakit.Func(defectDensity))
	reg.Register("fab.yield.negative_binomial", edakit.Func(negBinomial))
	reg.Register("fab.yield.critical_area", edakit.Func(criticalArea))
	reg.Register("fab.yield.average_critical_area", edakit.Func(avgCriticalArea))
	reg.Register("fab.yield.cost_per_die", edakit.Func(costPerDie))
	reg.Register("fab.yield.dfm", bank3(dfmItems, edakit.Func(dfmCompute)))

	// Design-technology co-optimization.
	reg.Register("fab.dtco.cell_architecture", edakit.Func(cellArchitecture))
	reg.Register("fab.dtco.cpp_cell_width", edakit.Func(cppWidth))
	reg.Register("fab.dtco.fin_depopulation", edakit.Func(finDepop))
	reg.Register("fab.dtco.buried_power_rail", edakit.Func(bpr))
	reg.Register("fab.dtco.backside_power", bank3(backsideItems, edakit.Func(bspdnCompute)))
	reg.Register("fab.dtco.dtco_loop", bank(dtcoLoopItems))
	reg.Register("fab.dtco.stco", edakit.Func(stco))

	// Frontier.
	reg.Register("fab.frontier.pdk_aware_agent", bank(pdkAgentItems))
	reg.Register("fab.frontier.fast_accurate_litho", bank(fastLithoItems))
	reg.Register("fab.frontier.rule_deck_golden", bank(ruleDeckItems))
	reg.Register("fab.frontier.process_ip_boundary", bank(processIPItems))
	reg.Register("fab.frontier.agentic_dtco_loop", bank(agenticDtcoItems))
}
