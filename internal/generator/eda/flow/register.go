// Package flow generates problems for the eda_landscape domain: the map of electronic design
// automation, the RTL-to-GDSII flow stage by stage, design representations and formats, PPA and
// manufacturing economics, adjacent branches, the industry, research history, and the agent-first
// frontier of the flow.
package flow

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds every eda_landscape generator to reg.
func Register(reg *generator.Registry) {
	// overview
	reg.Register("flow.overview.what_is_eda", &edakit.Bank{Items: whatIsEDAItems})
	reg.Register("flow.overview.scale_vlsi", edakit.Func(scaleVLSI))
	reg.Register("flow.overview.abstraction_levels", &edakit.Bank{Items: abstractionItems})
	reg.Register("flow.overview.hierarchy", edakit.Func(hierarchy))
	asicBank := &edakit.Bank{Items: asicItems}
	reg.Register("flow.overview.asic_fpga_soc", &edakit.Mix{Parts: []generator.Generator{asicBank, edakit.Func(breakEven), asicBank}})
	reg.Register("flow.overview.ip_libraries", &edakit.Bank{Items: ipItems})
	reg.Register("flow.overview.codesign", edakit.Func(amdahl))

	// stages
	reg.Register("flow.stages.rtl_to_gds", edakit.Func(rtlToGDS))
	reg.Register("flow.stages.spec_to_rtl", &edakit.Bank{Items: specItems})
	reg.Register("flow.stages.functional_verification", edakit.Func(verification))
	reg.Register("flow.stages.logic_synthesis", &edakit.Bank{Items: synthItems})
	reg.Register("flow.stages.dft", edakit.Func(scanTest))
	reg.Register("flow.stages.floorplan", edakit.Func(floorplan))
	reg.Register("flow.stages.placement", edakit.Func(placement))
	reg.Register("flow.stages.cts", edakit.Func(clockTree))
	reg.Register("flow.stages.routing", edakit.Func(routing))
	reg.Register("flow.stages.extraction_signoff", edakit.Func(signoff))
	reg.Register("flow.stages.tapeout_to_test", &edakit.Bank{Items: tapeoutItems})
	reg.Register("flow.stages.iteration_closure", &edakit.Bank{Items: iterationItems})
	reg.Register("flow.stages.error_propagation", edakit.Func(errorPropagation))

	// representations
	reg.Register("flow.repr.format_map", edakit.Func(formatMap))
	reg.Register("flow.repr.netlist", edakit.Func(netlist))
	reg.Register("flow.repr.boolean_networks", edakit.Func(booleanNetworks))
	reg.Register("flow.repr.liberty", edakit.Func(liberty))
	reg.Register("flow.repr.sdc", edakit.Func(sdc))
	reg.Register("flow.repr.lef_def", edakit.Func(lefDef))
	reg.Register("flow.repr.spef", edakit.Func(spef))
	reg.Register("flow.repr.gdsii_oasis", edakit.Func(gdsii))
	reg.Register("flow.repr.spice_pdk", &edakit.Bank{Items: pdkItems})

	// metrics
	reg.Register("flow.metrics.ppa", edakit.Func(ppa))
	reg.Register("flow.metrics.pareto", edakit.Func(pareto))
	reg.Register("flow.metrics.die_per_wafer", edakit.Func(diePerWafer))
	reg.Register("flow.metrics.yield_cost", edakit.Func(yieldCost))
	reg.Register("flow.metrics.respin_cost", edakit.Func(respinCost))

	// branches
	reg.Register("flow.branches.analog_ams", &edakit.Bank{Items: analogItems})
	fpgaBank := &edakit.Bank{Items: fpgaItems}
	reg.Register("flow.branches.fpga_flow", &edakit.Mix{Parts: []generator.Generator{fpgaBank, edakit.Func(lutCount), fpgaBank}})
	reg.Register("flow.branches.pcb_package", &edakit.Bank{Items: pcbItems})
	reg.Register("flow.branches.manufacturing", &edakit.Bank{Items: mfgItems})

	// industry
	reg.Register("flow.industry.foundry_fabless", &edakit.Bank{Items: foundryItems})
	reg.Register("flow.industry.vendors", &edakit.Bank{Items: vendorItems})
	reg.Register("flow.industry.tool_landscape", edakit.Func(toolLandscape))
	reg.Register("flow.industry.open_source", edakit.Func(openSource))
	reg.Register("flow.industry.standards", &edakit.Bank{Items: standardsItems})

	// history
	reg.Register("flow.history.landmarks", &edakit.Bank{Items: landmarkItems})
	reg.Register("flow.history.communities", &edakit.Bank{Items: communityItems})

	// frontier
	reg.Register("flow.frontier.engine_layers", &edakit.Bank{Items: engineLayerItems})
	reg.Register("flow.frontier.change_tiers", &edakit.Bank{Items: changeTierItems})
	reg.Register("flow.frontier.fast_accurate", edakit.Func(fastAccurate))
	reg.Register("flow.frontier.differential_qor", edakit.Func(differentialQoR))
	reg.Register("flow.frontier.agentic_flow", &edakit.Bank{Items: agenticFlowItems})
}
