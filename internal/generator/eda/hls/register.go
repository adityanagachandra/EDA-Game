package hls

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the high-level-synthesis generators to reg.
func Register(reg *generator.Registry) {
	// Representations.
	reg.Register("hls.ir.flow_overview", &edakit.Bank{Items: flowOverviewItems})
	reg.Register("hls.ir.dfg", edakit.Func(dfgProblem))
	reg.Register("hls.ir.cdfg", &edakit.Bank{Items: cdfgItems})
	reg.Register("hls.ir.ssa_lowering", &edakit.Bank{Items: ssaItems})
	reg.Register("hls.ir.op_latency", edakit.Func(opLatency))

	// Scheduling.
	reg.Register("hls.sched.asap_alap", edakit.Func(asapAlap))
	reg.Register("hls.sched.mobility", edakit.Func(mobility))
	reg.Register("hls.sched.resource_bound", edakit.Func(resourceBound))
	reg.Register("hls.sched.list", edakit.Func(listSched))
	reg.Register("hls.sched.min_resources", edakit.Func(minResources))
	reg.Register("hls.sched.force_directed", edakit.Func(forceDirected))
	reg.Register("hls.sched.ilp", edakit.Func(ilpCount))
	reg.Register("hls.sched.sdc_solve", edakit.Func(sdcSolve))
	reg.Register("hls.sched.chaining", edakit.Func(chaining))
	reg.Register("hls.sched.multicycle_fu", edakit.Func(multicycleFU))

	// Allocation and binding.
	reg.Register("hls.bind.allocation", edakit.Func(allocation))
	reg.Register("hls.bind.lifetimes", edakit.Func(lifetimes))
	reg.Register("hls.bind.left_edge", edakit.Func(leftEdge))
	reg.Register("hls.bind.interval_coloring", edakit.Func(intervalColoring))
	reg.Register("hls.bind.clique_partition", edakit.Func(cliquePartition))
	reg.Register("hls.bind.interconnect", edakit.Func(interconnect))
	reg.Register("hls.bind.datapath_fsm", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: fsmItems}, edakit.Func(fsmStates), &edakit.Bank{Items: fsmItems},
	}})

	// Pipelining.
	reg.Register("hls.pipe.initiation_interval", edakit.Func(initiationInterval))
	reg.Register("hls.pipe.loop_pipelining", edakit.Func(loopPipelining))
	reg.Register("hls.pipe.unrolling", edakit.Func(unrolling))
	reg.Register("hls.pipe.res_mii", edakit.Func(resMII))
	reg.Register("hls.pipe.rec_mii", edakit.Func(recMII))
	reg.Register("hls.pipe.modulo_scheduling", edakit.Func(moduloSched))
	reg.Register("hls.pipe.recurrence_breaking", edakit.Func(recurrenceBreaking))

	// Memory, interfaces and dataflow.
	reg.Register("hls.mem.partitioning", edakit.Func(partitioning))
	reg.Register("hls.mem.bank_conflicts", edakit.Func(bankConflicts))
	reg.Register("hls.mem.interfaces", &edakit.Bank{Items: interfaceItems})
	reg.Register("hls.mem.dataflow", edakit.Func(dataflow))
	reg.Register("hls.mem.line_buffers", edakit.Func(lineBuffers))

	// System level and tools.
	reg.Register("hls.sys.codesign", edakit.Func(codesign))
	reg.Register("hls.sys.dse", edakit.Func(dse))
	reg.Register("hls.sys.tools", &edakit.Bank{Items: toolItems})

	// Frontier.
	reg.Register("hls.frontier.seams", &edakit.Bank{Items: seamItems})
	reg.Register("hls.frontier.pragma_agent", &edakit.Bank{Items: pragmaAgentItems})
	reg.Register("hls.frontier.llm_to_hlsc", &edakit.Bank{Items: llmHLSItems})
	reg.Register("hls.frontier.cosim_golden", &edakit.Bank{Items: cosimItems})
	reg.Register("hls.frontier.qor_surrogate", edakit.Func(qorSurrogate))
}
