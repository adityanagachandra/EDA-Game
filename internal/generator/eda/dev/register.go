package dev

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the semiconductor-device generators to reg.
func Register(reg *generator.Registry) {
	// physics
	reg.Register("dev.physics.bands_fermi", mix(bandsItems, fermiOccupancy))
	reg.Register("dev.physics.intrinsic", fn(intrinsic))
	reg.Register("dev.physics.mass_action", fn(massAction))
	reg.Register("dev.physics.fermi_position", fn(fermiPosition))
	reg.Register("dev.physics.drift_mobility", fn(driftMobility))
	reg.Register("dev.physics.diffusion_einstein", fn(diffusionEinstein))
	reg.Register("dev.physics.drift_diffusion", &edakit.Bank{Items: driftDiffusionItems})
	// junction
	reg.Register("dev.junction.pn", fn(pnJunction))
	reg.Register("dev.junction.depletion", fn(depletion))
	reg.Register("dev.junction.diode", fn(diode))
	reg.Register("dev.junction.bjt", fn(bjt))
	reg.Register("dev.junction.tunneling", fn(tunneling))
	// mos
	reg.Register("dev.mos.capacitor", fn(mosCapacitor))
	reg.Register("dev.mos.threshold", fn(threshold))
	reg.Register("dev.mos.body_effect", fn(bodyEffect))
	reg.Register("dev.mos.linear_region", fn(linearRegion))
	reg.Register("dev.mos.square_law", fn(squareLaw))
	reg.Register("dev.mos.clm", fn(clm))
	reg.Register("dev.mos.velocity_saturation", fn(velocitySaturation))
	reg.Register("dev.mos.subthreshold", fn(subthreshold))
	reg.Register("dev.mos.gate_capacitance", fn(gateCapacitance))
	reg.Register("dev.mos.cmos_pair", mix(cmosItems, cmosRatio))
	reg.Register("dev.mos.short_channel", &edakit.Bank{Items: shortChannelItems})
	reg.Register("dev.mos.leakage", &edakit.Bank{Items: leakageItems})
	// short-channel quantities
	reg.Register("dev.sce.dibl", fn(dibl))
	reg.Register("dev.sce.vt_rolloff", fn(vtRolloff))
	reg.Register("dev.sce.vt_flavors", fn(vtFlavors))
	reg.Register("dev.sce.gate_leakage", fn(gateLeakage))
	// advanced devices
	reg.Register("dev.advanced.soi", mix(soiItems, soiBackBias))
	reg.Register("dev.advanced.finfet", fn(finfet))
	reg.Register("dev.advanced.gaa", fn(gaa))
	reg.Register("dev.advanced.confinement", fn(confinement))
	reg.Register("dev.advanced.cfet", &edakit.Bank{Items: cfetItems})
	// variability
	reg.Register("dev.var.sources", &edakit.Bank{Items: sourcesItems})
	reg.Register("dev.var.pelgrom", fn(pelgrom))
	reg.Register("dev.var.corners", mix(cornerItems, cornerSigma))
	// reliability
	reg.Register("dev.rel.bti_hci", mix(btiHciItems, btiPowerLaw))
	reg.Register("dev.rel.tddb", fn(tddb))
	reg.Register("dev.rel.acceleration", fn(acceleration))
	reg.Register("dev.rel.aging_models", &edakit.Bank{Items: agingItems})
	// scaling
	reg.Register("dev.scaling.dennard", fn(dennard))
	reg.Register("dev.scaling.power_wall", fn(powerWall))
	reg.Register("dev.scaling.vt_floor", fn(vtFloor))
	// models
	reg.Register("dev.model.tcad", &edakit.Bank{Items: tcadItems})
	reg.Register("dev.model.compact", &edakit.Bank{Items: compactItems})
	reg.Register("dev.model.unified_iv", fn(unifiedIV))
	reg.Register("dev.model.bsim", &edakit.Bank{Items: bsimItems})
	reg.Register("dev.model.extraction", fn(extraction))
	// frontier
	reg.Register("dev.frontier.extraction_skill", &edakit.Bank{Items: extractionSkillItems})
	reg.Register("dev.frontier.model_seam", &edakit.Bank{Items: modelSeamItems})
	reg.Register("dev.frontier.golden_iv", mix(goldenItems, goldenIV))
	reg.Register("dev.frontier.tcad_surrogate", &edakit.Bank{Items: tcadSurrogateItems})
	reg.Register("dev.frontier.fit_reward", mix(fitRewardItems, fitRMS))
}
