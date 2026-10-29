// Package cmos generates problems on CMOS circuits and standard cells: static and dynamic
// gates, the inverter, RC delay and logical effort, power, sequential cells, standard-cell
// layout, the Liberty library abstraction and its characterization, special low-power and
// physical cells, and the agent-first future of cell-library work.
package cmos

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func bank(items []edakit.Item) *edakit.Bank { return &edakit.Bank{Items: items} }

// mix serves a computed part on one knowledge point and a bank on the other two.
func mix(kp int, f edakit.Func, b *edakit.Bank) *edakit.Mix {
	parts := []generator.Generator{b, b, b}
	parts[kp-1] = f
	return &edakit.Mix{Parts: parts}
}

// Register adds the CMOS circuits and cells generators to reg.
func Register(reg *generator.Registry) {
	// gate
	reg.Register("cmos.gate.switch_model", bank(switchModelItems))
	reg.Register("cmos.gate.static_cmos", edakit.Func(staticCMOS))
	reg.Register("cmos.gate.transistor_count", edakit.Func(transistorCount))
	reg.Register("cmos.gate.aoi_oai", edakit.Func(aoiOai))
	reg.Register("cmos.gate.transmission_gate", mix(1, passVoltage, bank(tgateItems)))
	reg.Register("cmos.gate.inverter_vtc", edakit.Func(inverterVM))
	reg.Register("cmos.gate.noise_margins", edakit.Func(noiseMargins))
	reg.Register("cmos.gate.dynamic_logic", mix(2, chargeSharing, bank(dynamicItems)))
	// delay
	reg.Register("cmos.delay.rc_model", edakit.Func(rcModel))
	reg.Register("cmos.delay.equal_rise_fall", edakit.Func(equalRiseFall))
	reg.Register("cmos.delay.stack_elmore", edakit.Func(stackElmore))
	reg.Register("cmos.delay.logical_effort", edakit.Func(logicalEffort))
	reg.Register("cmos.delay.parasitic_delay", edakit.Func(parasiticDelay))
	reg.Register("cmos.delay.stage_delay", edakit.Func(stageDelay))
	reg.Register("cmos.delay.fo4", edakit.Func(fo4))
	reg.Register("cmos.delay.path_effort", edakit.Func(pathEffort))
	reg.Register("cmos.delay.optimal_stages", edakit.Func(optimalStages))
	reg.Register("cmos.delay.path_sizing", edakit.Func(pathSizing))
	reg.Register("cmos.delay.slew_dependence", bank(slewItems))
	// power
	reg.Register("cmos.power.dynamic", edakit.Func(dynamicPower))
	reg.Register("cmos.power.activity_factor", edakit.Func(activityFactor))
	reg.Register("cmos.power.short_circuit", bank(shortCircuitItems))
	reg.Register("cmos.power.leakage", edakit.Func(leakage))
	reg.Register("cmos.power.multi_vt", edakit.Func(multiVt))
	reg.Register("cmos.power.energy_per_op", edakit.Func(energyPerOp))
	// seq
	reg.Register("cmos.seq.latch_circuits", bank(latchItems))
	reg.Register("cmos.seq.setup_hold_characterization", edakit.Func(setupHoldChar))
	reg.Register("cmos.seq.clock_gating_cell", bank(icgItems))
	reg.Register("cmos.seq.metastability", edakit.Func(metastability))
	// cell
	reg.Register("cmos.cell.standard_cell_layout", edakit.Func(cellLayout))
	reg.Register("cmos.cell.library_views", bank(libraryViewItems))
	reg.Register("cmos.cell.row_arithmetic", edakit.Func(rowArithmetic))
	reg.Register("cmos.cell.drive_strength", edakit.Func(driveStrength))
	// lib
	reg.Register("cmos.lib.liberty_structure", mix(1, unatenessPart, bank(libStructureItems)))
	reg.Register("cmos.lib.liberty_nldm", edakit.Func(nldmLookup))
	reg.Register("cmos.lib.nldm_abstraction", bank(nldmAbstractionItems))
	reg.Register("cmos.lib.liberty_ccs_lvf", edakit.Func(ccsLvf))
	reg.Register("cmos.lib.power_noise_views", edakit.Func(powerViews))
	reg.Register("cmos.lib.characterization", edakit.Func(characterization))
	reg.Register("cmos.lib.corners", bank(cornerItems))
	reg.Register("cmos.lib.sequential_arcs", edakit.Func(sequentialArcs))
	// special
	reg.Register("cmos.special.lowpower_cells", bank(lowpowerItems))
	reg.Register("cmos.special.power_gating", edakit.Func(powerGating))
	reg.Register("cmos.special.physical_cells", mix(2, decapPart, bank(physicalItems)))
	// frontier
	reg.Register("cmos.frontier.characterization_skill", bank(charSkillItems))
	reg.Register("cmos.frontier.golden_spice_check", edakit.Func(goldenCheck))
	reg.Register("cmos.frontier.fast_char_surrogate", mix(2, surrogatePart, bank(surrogateItems)))
	reg.Register("cmos.frontier.dtco_cell_agent", mix(2, paretoPart, bank(dtcoItems)))
	reg.Register("cmos.frontier.library_reward", bank(libraryRewardItems))
}
