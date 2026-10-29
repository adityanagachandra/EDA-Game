package cts

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the clock and power generators to reg.
func Register(reg *generator.Registry) {
	fn := func(g func(generator.GeneratorContext) generator.Problem) edakit.Func { return edakit.Func(g) }

	// Clock tree objectives.
	reg.Register("cts.objective.skew_latency", fn(skewLatency))
	reg.Register("cts.objective.clock_slew", fn(clockSlew))
	reg.Register("cts.objective.clock_power", fn(clockPower))

	// Clock topologies.
	reg.Register("cts.topology.h_tree", fn(hTree))
	reg.Register("cts.topology.mesh", fn(mesh))
	reg.Register("cts.topology.spine_hybrid", bank(spineItems))
	reg.Register("cts.topology.tree_skew", fn(treeSkew))

	// Zero-skew and bounded-skew trees.
	reg.Register("cts.zeroskew.linear_merge", fn(linearMerge))
	reg.Register("cts.zeroskew.tsay_merge", fn(tsayMerge))
	reg.Register("cts.zeroskew.merging_segment", fn(mergingSegment))
	reg.Register("cts.zeroskew.topology_generation", fn(topologyGeneration))
	reg.Register("cts.zeroskew.dme", mix3(fn(dmeBottomUp), fn(dmeTopDown), bank(dmeItems)))
	reg.Register("cts.zeroskew.bounded_skew", fn(boundedSkew))

	// Buffered clock trees.
	reg.Register("cts.buffered.buffer_levels", fn(bufferLevels))
	reg.Register("cts.buffered.slew_wire_length", fn(slewWireLength))
	reg.Register("cts.buffered.duty_cycle", fn(dutyCycle))
	reg.Register("cts.buffered.latency_balancing", fn(latencyBalancing))
	reg.Register("cts.buffered.cts_flow", bank(ctsFlowItems))

	// Useful skew and clock scheduling.
	reg.Register("cts.skew.useful_skew", fn(usefulSkew))
	reg.Register("cts.skew.difference_constraints", fn(differenceConstraints))
	reg.Register("cts.skew.scheduling_lp", fn(schedulingLP))

	// Clock gating.
	reg.Register("cts.gating.icg", mix3(bank(icgItems), fn(icgGlitch), bank(icgItems)))
	reg.Register("cts.gating.power_savings", fn(gatingSavings))
	reg.Register("cts.gating.enable_timing", fn(enableTiming))

	// Power distribution networks.
	reg.Register("cts.pdn.grid", fn(pdnGrid))
	reg.Register("cts.pdn.rail_ladder", fn(railLadder))
	reg.Register("cts.pdn.conductance_matrix", fn(conductanceMatrix))
	reg.Register("cts.pdn.ir_drop", fn(irDrop))
	reg.Register("cts.pdn.iterative_solve", fn(iterativeSolve))
	reg.Register("cts.pdn.strap_pitch", fn(strapPitch))
	reg.Register("cts.pdn.dynamic_ir", fn(dynamicIR))
	reg.Register("cts.pdn.decap", fn(decap))
	reg.Register("cts.pdn.ldi_dt", fn(ldiDt))
	reg.Register("cts.pdn.electromigration", fn(electromigration))
	reg.Register("cts.pdn.black_equation", fn(blackEquation))

	// Low-power design.
	reg.Register("cts.lowpower.multi_voltage", fn(multiVoltage))
	reg.Register("cts.lowpower.power_gating", fn(powerGating))
	reg.Register("cts.lowpower.switch_sizing", fn(switchSizing))
	reg.Register("cts.lowpower.rush_current", fn(rushCurrent))
	reg.Register("cts.lowpower.upf", bank(upfItems))
	reg.Register("cts.lowpower.isolation_level_shift", mix3(fn(crossingCells), fn(clampValue), fn(shifterDirection)))

	// Power analysis.
	reg.Register("cts.analysis.activity_propagation", fn(activityPropagation))
	reg.Register("cts.analysis.transition_density", fn(transitionDensity))
	reg.Register("cts.analysis.vector_vs_vectorless", mix3(bank(vectorItems), fn(saifPower), bank(vectorItems)))

	// Frontier.
	reg.Register("cts.frontier.ir_eco_loop", mix3(bank(irEcoItems), fn(irEcoChoice), bank(irEcoItems)))
	reg.Register("cts.frontier.fast_accurate_pdn", bank(fastPdnItems))
	reg.Register("cts.frontier.golden_ir", mix3(bank(goldenItems), fn(goldenCount), bank(goldenItems)))
	reg.Register("cts.frontier.cts_policy_plugin", bank(policyItems))
	reg.Register("cts.frontier.pdn_reward", bank(rewardItems))
}
