// Package dft generates manufacturing-test and design-for-test problems: fault models and
// collapsing, fault simulation, ATPG and testability, scan and compression, logic and memory
// BIST, boundary scan, the test floor (ATE, defect level, diagnosis, silicon debug) and the
// agent-first frontier. Files are split by subdomain; circuit.go holds the gate-level engine.
package dft

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func bank(items []edakit.Item) *edakit.Bank { return &edakit.Bank{Items: items} }

func mix(parts ...generator.Generator) *edakit.Mix { return &edakit.Mix{Parts: parts} }

// Register adds the test-and-DFT generators to reg.
func Register(reg *generator.Registry) {
	// fault models
	reg.Register("dft.fault.defects_vs_faults", bank(defectItems))
	reg.Register("dft.fault.stuck_at", edakit.Func(stuckAt))
	reg.Register("dft.fault.transition", edakit.Func(transition))
	reg.Register("dft.fault.path_delay", edakit.Func(pathDelay))
	reg.Register("dft.fault.bridging", edakit.Func(bridging))
	reg.Register("dft.fault.cell_aware", mix(bank(cellAwareItems), edakit.Func(cellAwareCalc), bank(cellAwareItems)))
	reg.Register("dft.fault.collapsing", edakit.Func(collapsing))
	// fault simulation
	reg.Register("dft.fsim.fault_detection", edakit.Func(faultDetection))
	reg.Register("dft.fsim.serial", edakit.Func(serialSim))
	reg.Register("dft.fsim.parallel", edakit.Func(parallelSim))
	reg.Register("dft.fsim.deductive", edakit.Func(deductive))
	reg.Register("dft.fsim.concurrent", mix(bank(concurrentItems), edakit.Func(concurrentCalc), bank(concurrentItems)))
	reg.Register("dft.fsim.coverage", edakit.Func(coverage))
	// ATPG
	reg.Register("dft.atpg.d_calculus", edakit.Func(dCalculus))
	reg.Register("dft.atpg.sensitization", edakit.Func(sensitization))
	reg.Register("dft.atpg.d_algorithm", mix(edakit.Func(dAlgorithmCalc), edakit.Func(dAlgorithmCalc), bank(dAlgItems)))
	reg.Register("dft.atpg.scoap", edakit.Func(scoapCC))
	reg.Register("dft.atpg.scoap_observability", edakit.Func(scoapCO))
	reg.Register("dft.atpg.podem", mix(edakit.Func(podemCalc), edakit.Func(podemCalc), bank(podemItems)))
	reg.Register("dft.atpg.fan", bank(fanItems))
	reg.Register("dft.atpg.sat_based", edakit.Func(satATPG))
	reg.Register("dft.atpg.random_pattern", edakit.Func(randomPattern))
	reg.Register("dft.atpg.compaction", edakit.Func(compaction))
	// scan
	reg.Register("dft.scan.chain", edakit.Func(scanChain))
	reg.Register("dft.scan.insertion", edakit.Func(scanInsertion))
	reg.Register("dft.scan.ordering", edakit.Func(scanOrdering))
	reg.Register("dft.scan.loc_los", mix(edakit.Func(locLosCalc), edakit.Func(locLosCalc), bank(locLosItems)))
	reg.Register("dft.scan.compression", edakit.Func(scanCompression))
	// compression
	reg.Register("dft.compress.decompressor", edakit.Func(decompressor))
	reg.Register("dft.compress.compactor", edakit.Func(compactor))
	// logic BIST
	reg.Register("dft.bist.lfsr_misr", edakit.Func(lfsrMISR))
	reg.Register("dft.bist.lfsr_period", edakit.Func(lfsrPeriod))
	reg.Register("dft.bist.aliasing", edakit.Func(aliasing))
	reg.Register("dft.bist.lbist", edakit.Func(lbist))
	// memory BIST
	reg.Register("dft.mbist.memory_faults", bank(memFaultItems))
	reg.Register("dft.mbist.march", edakit.Func(march))
	reg.Register("dft.mbist.march_simulation", edakit.Func(marchSim))
	reg.Register("dft.mbist.repair", edakit.Func(repair))
	reg.Register("dft.mbist.ecc", edakit.Func(ecc))
	// boundary scan
	reg.Register("dft.jtag.boundary_scan", edakit.Func(boundaryScan))
	reg.Register("dft.jtag.tap_fsm", edakit.Func(tapFSM))
	reg.Register("dft.jtag.instructions", mix(bank(instructionItems), edakit.Func(instructionsCalc), bank(instructionItems)))
	reg.Register("dft.jtag.ieee1500_1687", mix(bank(ijtagItems), edakit.Func(ijtagCalc), bank(ijtagItems)))
	// test floor
	reg.Register("dft.ate.test_flow", edakit.Func(testFlow))
	reg.Register("dft.ate.stil", bank(stilItems))
	reg.Register("dft.ate.defect_level", edakit.Func(defectLevel))
	reg.Register("dft.ate.diagnosis", edakit.Func(diagnosis))
	reg.Register("dft.ate.silicon_debug", mix(bank(siliconDebugItems), edakit.Func(siliconDebugCalc), bank(siliconDebugItems)))
	// frontier
	reg.Register("dft.frontier.atpg_skill", bank(atpgSkillItems))
	reg.Register("dft.frontier.dft_insertion_agent", bank(insertionAgentItems))
	reg.Register("dft.frontier.golden_coverage_tests", edakit.Func(goldenCoverage))
	reg.Register("dft.frontier.diagnosis_agent", edakit.Func(diagnosisAgent))
	reg.Register("dft.frontier.fast_coverage_estimate", edakit.Func(fastCoverage))
}
