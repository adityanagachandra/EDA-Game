// Package extract generates parasitic-extraction and signoff problems: wire resistance,
// capacitance components, extraction methods, inductance, SPEF, RC reduction, signal
// integrity, signoff and ECOs, and the agent-first frontier. Files are split by subdomain;
// this one holds Register.
package extract

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func mix(items []edakit.Item, fn edakit.Func) generator.Generator {
	return &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: items}, fn}}
}

// Register adds the extraction and signoff generators to reg.
func Register(reg *generator.Registry) {
	reg.Register("extract.res.sheet_resistance", edakit.Func(sheetResistance))
	reg.Register("extract.res.wire_resistance", edakit.Func(wireResistance))
	reg.Register("extract.res.temperature", edakit.Func(temperature))
	reg.Register("extract.res.barrier_scattering", edakit.Func(barrierScattering))

	reg.Register("extract.cap.parallel_plate", edakit.Func(parallelPlate))
	reg.Register("extract.cap.fringe", edakit.Func(fringe))
	reg.Register("extract.cap.coupling", edakit.Func(coupling))
	reg.Register("extract.cap.cap_matrix", edakit.Func(capMatrix))
	reg.Register("extract.cap.total_net", edakit.Func(totalNet))
	reg.Register("extract.cap.field_solver", mix(fieldSolverItems, fieldSolverCalc))

	reg.Register("extract.method.pattern_matching", edakit.Func(patternMatching))
	reg.Register("extract.method.random_walk", edakit.Func(randomWalk))
	reg.Register("extract.method.bem_fem", edakit.Func(bemFem))
	reg.Register("extract.method.calibration", edakit.Func(calibration))

	reg.Register("extract.ind.partial_inductance", edakit.Func(partialInductance))
	reg.Register("extract.ind.significance", edakit.Func(significance))

	reg.Register("extract.format.spef", mix(spefItems, spefCalc))
	reg.Register("extract.format.spef_counts", edakit.Func(spefCounts))
	reg.Register("extract.format.name_map", edakit.Func(nameMap))

	reg.Register("extract.reduce.rc_reduction", mix(rcReductionItems, rcReductionCalc))
	reg.Register("extract.reduce.moments", edakit.Func(moments))
	reg.Register("extract.reduce.node_elimination", edakit.Func(nodeElimination))
	reg.Register("extract.reduce.elmore_preserving", edakit.Func(elmorePreserving))
	reg.Register("extract.reduce.pi_model", edakit.Func(piModel))

	reg.Register("extract.si.miller_factor", edakit.Func(millerFactor))
	reg.Register("extract.si.effective_cap", edakit.Func(effectiveCap))
	reg.Register("extract.si.crosstalk_delay", edakit.Func(crosstalkDelay))
	reg.Register("extract.si.glitch", edakit.Func(glitch))
	reg.Register("extract.si.noise_immunity", edakit.Func(noiseImmunity))
	reg.Register("extract.si.timing_windows", edakit.Func(timingWindows))

	reg.Register("extract.signoff.rc_corners", edakit.Func(rcCorners))
	reg.Register("extract.signoff.em_ir", edakit.Func(emIR))
	reg.Register("extract.signoff.checklist", &edakit.Bank{Items: checklistItems})
	reg.Register("extract.signoff.eco", &edakit.Bank{Items: ecoItems})
	reg.Register("extract.signoff.spare_cells", edakit.Func(spareCells))
	reg.Register("extract.signoff.ppa_closure", mix(ppaItems, ppaCalc))

	reg.Register("extract.frontier.core_vs_surface", &edakit.Bank{Items: coreSurfaceItems})
	reg.Register("extract.frontier.calibration_tier1", edakit.Func(calibrationTier1))
	reg.Register("extract.frontier.fast_accurate", edakit.Func(fastAccurate))
	reg.Register("extract.frontier.differential_eco_check", edakit.Func(differentialEcoCheck))
	reg.Register("extract.frontier.agentic_signoff", &edakit.Bank{Items: agenticSignoffItems})
}
