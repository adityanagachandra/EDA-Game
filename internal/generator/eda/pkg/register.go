package pkg

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the packaging and integration generators to reg.
func Register(reg *generator.Registry) {
	// Package basics and wafer-level packaging.
	reg.Register("pkg.basic.package_roles", bank(packageRoleItems))
	reg.Register("pkg.basic.wire_bonding", edakit.Func(wireBonding))
	reg.Register("pkg.basic.flip_chip_bga", edakit.Func(flipChipBGA))
	reg.Register("pkg.basic.substrates", edakit.Func(substrates))
	reg.Register("pkg.basic.escape_routing", edakit.Func(escapeRouting))
	reg.Register("pkg.wlp.wafer_level", edakit.Func(waferLevel))
	reg.Register("pkg.wlp.rdl", edakit.Func(rdl))

	// Advanced packaging.
	reg.Register("pkg.adv.sip_mcm", edakit.Func(sipMCM))
	reg.Register("pkg.adv.interposer_25d", edakit.Func(interposer))
	reg.Register("pkg.adv.bridges", edakit.Func(bridges))
	reg.Register("pkg.adv.3dic_tsv", edakit.Func(tsvBasics))
	reg.Register("pkg.adv.tsv_keep_out", edakit.Func(tsvKeepOut))
	reg.Register("pkg.adv.tsv_capacitance", edakit.Func(tsvCapacitance))
	reg.Register("pkg.adv.microbumps", edakit.Func(microbumps))
	reg.Register("pkg.adv.hybrid_bonding", edakit.Func(hybridBonding))

	// Chiplets and die-to-die links.
	reg.Register("pkg.chiplet.disaggregation", edakit.Func(disaggregation))
	reg.Register("pkg.chiplet.d2d_interfaces", bank(d2dItems))
	reg.Register("pkg.chiplet.ucie", edakit.Func(ucie))
	reg.Register("pkg.chiplet.energy_per_bit", edakit.Func(energyPerBit))
	reg.Register("pkg.chiplet.package_board_codesign", edakit.Func(packageBoard))

	// Thermal.
	reg.Register("pkg.thermal.conduction", edakit.Func(conduction))
	reg.Register("pkg.thermal.resistance_network", edakit.Func(resistanceNetwork))
	reg.Register("pkg.thermal.junction_temperature", edakit.Func(junctionTemp))
	reg.Register("pkg.thermal.power_density", edakit.Func(powerDensity))
	reg.Register("pkg.thermal.joule_heating", edakit.Func(jouleHeating))
	reg.Register("pkg.thermal.simulation", edakit.Func(thermalSim))
	reg.Register("pkg.thermal.transient", edakit.Func(transient))
	reg.Register("pkg.thermal.management", edakit.Func(management))

	// Mechanical.
	reg.Register("pkg.mech.thermal_stress", edakit.Func(thermalStress))
	reg.Register("pkg.mech.cte_warpage", edakit.Func(warpage))
	reg.Register("pkg.mech.solder_fatigue", edakit.Func(solderFatigue))
	reg.Register("pkg.mech.fea", bank(feaItems))

	// Reliability and multiphysics.
	reg.Register("pkg.rel.electromigration", edakit.Func(electromigration))
	reg.Register("pkg.rel.em_temperature", edakit.Func(emTemperature))
	reg.Register("pkg.rel.stress_migration", bank(stressMigrationItems))
	reg.Register("pkg.rel.dielectric_breakdown", edakit.Func(tddb))
	reg.Register("pkg.rel.aging", edakit.Func(aging))
	reg.Register("pkg.rel.multiphysics", edakit.Func(multiphysics))

	// 3D-IC design.
	reg.Register("pkg.eda3d.partitioning", edakit.Func(partitioning))
	reg.Register("pkg.eda3d.inter_die_timing", edakit.Func(interDieTiming))
	reg.Register("pkg.eda3d.thermal_floorplan", edakit.Func(thermalFloorplan))
	reg.Register("pkg.eda3d.multi_die_verification", bank(multiDieItems))

	// Frontier.
	reg.Register("pkg.frontier.multiphysics_engine_layers", bank(engineLayerItems))
	reg.Register("pkg.frontier.thermal_surrogate", bank(surrogateItems))
	reg.Register("pkg.frontier.chiplet_partition_agent", bank(partitionAgentItems))
	reg.Register("pkg.frontier.golden_multiphysics", bank(goldenItems))
	reg.Register("pkg.frontier.codesign_loop", bank(codesignItems))
}
