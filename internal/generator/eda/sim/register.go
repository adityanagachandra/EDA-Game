// Package sim generates functional-verification problems: event-driven and cycle-based
// simulation engines, testbenches, constrained-random stimulus, coverage, assertions, UVM,
// clock-domain crossing, emulation and agentic verification loops.
package sim

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the functional-verification generators to reg.
func Register(reg *generator.Registry) {
	// Event-driven simulation.
	reg.Register("sim.event.event_driven", edakit.Func(eventDriven))
	reg.Register("sim.event.event_wheel", edakit.Func(eventWheel))
	reg.Register("sim.event.delta_cycles", edakit.Func(deltaCycles))
	reg.Register("sim.event.scheduling_regions", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: schedulingItems}, edakit.Func(blockingNonblocking), &edakit.Bank{Items: schedulingItems}}})
	reg.Register("sim.event.delay_models", edakit.Func(delayModels))
	reg.Register("sim.event.four_state", edakit.Func(fourState))
	reg.Register("sim.event.x_pessimism", edakit.Func(xPessimism))

	// Cycle-based and compiled simulation.
	reg.Register("sim.cycle.levelization", edakit.Func(levelization))
	reg.Register("sim.cycle.compiled", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: compiledItems}, edakit.Func(cycleCost), &edakit.Bank{Items: compiledItems}}})
	reg.Register("sim.cycle.verilator", &edakit.Bank{Items: verilatorItems})
	reg.Register("sim.cycle.parallel", edakit.Func(parallelSim))

	// Testbenches.
	reg.Register("sim.tb.verification_gap", edakit.Func(verificationGap))
	reg.Register("sim.tb.testbench", &edakit.Bank{Items: testbenchItems})
	reg.Register("sim.tb.scoreboard", edakit.Func(scoreboard))
	reg.Register("sim.tb.cocotb", &edakit.Bank{Items: cocotbItems})

	// Constrained-random stimulus.
	reg.Register("sim.random.constrained", edakit.Func(constrained))
	reg.Register("sim.random.dist_solve_before", edakit.Func(distSolveBefore))
	reg.Register("sim.random.corner_probability", edakit.Func(cornerProbability))
	reg.Register("sim.random.constraint_solving", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: solverItems}, edakit.Func(rejectionSampling), &edakit.Bank{Items: solverItems}}})
	reg.Register("sim.random.seeds", &edakit.Bank{Items: seedItems})

	// Coverage.
	reg.Register("sim.cov.code_coverage", edakit.Func(codeCoverage))
	reg.Register("sim.cov.condition", edakit.Func(conditionCoverage))
	reg.Register("sim.cov.functional_coverage", edakit.Func(functionalCoverage))
	reg.Register("sim.cov.cross_coverage", edakit.Func(crossCoverage))
	reg.Register("sim.cov.coupon_collector", edakit.Func(couponCollector))
	reg.Register("sim.cov.closure", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: closureItems}, edakit.Func(greedyRanking), &edakit.Bank{Items: closureItems}}})

	// Assertions.
	reg.Register("sim.assert.in_simulation", &edakit.Bank{Items: assertItems})
	reg.Register("sim.assert.implication", edakit.Func(implication))

	// UVM and TLM.
	reg.Register("sim.uvm.architecture", &edakit.Bank{Items: uvmArchItems})
	reg.Register("sim.uvm.phasing", &edakit.Bank{Items: phasingItems})
	reg.Register("sim.uvm.sequences", &edakit.Bank{Items: sequenceItems})
	reg.Register("sim.uvm.factory_config", &edakit.Bank{Items: factoryItems})
	reg.Register("sim.uvm.tlm", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: tlmItems}, &edakit.Bank{Items: tlmItems}, edakit.Func(temporalDecoupling)}})

	// Clock- and reset-domain crossing.
	reg.Register("sim.cdc.synchronizers", &edakit.Bank{Items: synchronizerItems})
	reg.Register("sim.cdc.mtbf", edakit.Func(mtbf))
	reg.Register("sim.cdc.multibit", edakit.Func(multibit))
	reg.Register("sim.cdc.reset_domain", &edakit.Bank{Items: resetItems})

	// Emulation and prototyping.
	reg.Register("sim.accel.emulation", &edakit.Bank{Items: emulationItems})
	reg.Register("sim.accel.throughput", edakit.Func(throughput))
	reg.Register("sim.accel.partitioning", edakit.Func(partitioning))

	// Frontier: agentic verification.
	reg.Register("sim.frontier.coverage_agent", &edakit.Bank{Items: coverageAgentItems})
	reg.Register("sim.frontier.llm_testbench", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: llmTestbenchItems}, edakit.Func(mutationScore), &edakit.Bank{Items: llmTestbenchItems}}})
	reg.Register("sim.frontier.debug_agent", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: debugAgentItems}, edakit.Func(bisectReruns), &edakit.Bank{Items: debugAgentItems}}})
	reg.Register("sim.frontier.differential_sim", &edakit.Bank{Items: differentialItems})
	reg.Register("sim.frontier.simulator_seams", &edakit.Bank{Items: seamItems})
}
