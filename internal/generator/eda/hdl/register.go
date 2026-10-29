// Package hdl generates problems on hardware description languages and the compilers that
// read them: RTL semantics that tools depend on (blocking and nonblocking assignment, flop and
// latch inference, sizing and signedness), elaboration, compiler internals (lexing, parsing,
// ASTs, control-flow graphs, SSA, dominators, data-flow analysis, classical optimizations) and
// the real intermediate representations (Yosys RTLIL, MLIR and CIRCT, FIRRTL).
package hdl

import (
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the HDL and compiler generators to reg.
func Register(reg *generator.Registry) {
	// RTL foundations.
	reg.Register("hdl.rtl.register_transfer", &edakit.Bank{Items: registerTransferItems})
	reg.Register("hdl.rtl.behavioral_structural", &edakit.Bank{Items: behavioralStructuralItems})
	reg.Register("hdl.rtl.four_state", edakit.Func(fourState))
	reg.Register("hdl.rtl.modules_hierarchy", &edakit.Bank{Items: modulesHierarchyItems})

	// Verilog semantics.
	reg.Register("hdl.verilog.event_scheduling", &edakit.Bank{Items: eventSchedulingItems})
	reg.Register("hdl.verilog.blocking_nonblocking", edakit.Func(blockingNonblocking))
	reg.Register("hdl.verilog.flop_inference", edakit.Func(flopInference))
	reg.Register("hdl.verilog.latch_inference", edakit.Func(latchInference))
	reg.Register("hdl.verilog.case_semantics", edakit.Func(caseSemantics))
	reg.Register("hdl.verilog.synthesizable_subset", &edakit.Bank{Items: synthSubsetItems})
	reg.Register("hdl.verilog.expression_sizing", edakit.Func(expressionSizing))
	reg.Register("hdl.verilog.signedness", edakit.Func(signedness))

	// SystemVerilog.
	reg.Register("hdl.sv.always_ff_comb", &edakit.Bank{Items: alwaysFFCombItems})
	reg.Register("hdl.sv.logic_types", edakit.Func(logicTypes))
	reg.Register("hdl.sv.interfaces_packages", &edakit.Bank{Items: interfacesItems})
	reg.Register("hdl.sv.unique_priority", &edakit.Bank{Items: uniquePriorityItems})

	// Other languages.
	reg.Register("hdl.lang.vhdl_signals_variables", edakit.Func(vhdlSignalsVariables))
	reg.Register("hdl.lang.vhdl_resolution", edakit.Func(vhdlResolution))
	reg.Register("hdl.lang.systemc", &edakit.Bank{Items: systemCItems})
	reg.Register("hdl.lang.ams", &edakit.Bank{Items: amsItems})
	reg.Register("hdl.lang.chisel_dsl", &edakit.Bank{Items: chiselItems})

	// Compiler front end and analysis.
	reg.Register("hdl.compile.lexing", edakit.Func(lexing))
	reg.Register("hdl.compile.parsing", edakit.Func(parsing))
	reg.Register("hdl.compile.ast", edakit.Func(astShape))
	reg.Register("hdl.compile.elaboration", &edakit.Bank{Items: elaborationItems})
	reg.Register("hdl.compile.ir", &edakit.Bank{Items: irItems})
	reg.Register("hdl.compile.cfg", edakit.Func(cfgCount))
	reg.Register("hdl.compile.ssa", edakit.Func(ssaRenaming))
	reg.Register("hdl.compile.dominators", edakit.Func(dominators))
	reg.Register("hdl.compile.phi_placement", edakit.Func(phiPlacement))
	reg.Register("hdl.compile.dataflow", edakit.Func(dataflow))

	// Elaboration.
	reg.Register("hdl.elab.parameters", edakit.Func(parameters))
	reg.Register("hdl.elab.generate", edakit.Func(generateLoops))
	reg.Register("hdl.elab.flatten", edakit.Func(flatten))

	// Optimization.
	reg.Register("hdl.opt.constant_propagation", edakit.Func(constantPropagation))
	reg.Register("hdl.opt.cse", edakit.Func(cse))
	reg.Register("hdl.opt.dce", edakit.Func(dce))

	// Real IRs and tools.
	reg.Register("hdl.ir.rtlil", &edakit.Bank{Items: rtlilItems})
	reg.Register("hdl.ir.yosys_proc", &edakit.Bank{Items: yosysProcItems})
	reg.Register("hdl.ir.yosys_opt", &edakit.Bank{Items: yosysOptItems})
	reg.Register("hdl.ir.frontends", &edakit.Bank{Items: frontendItems})
	reg.Register("hdl.ir.verilator", &edakit.Bank{Items: verilatorItems})
	reg.Register("hdl.ir.mlir_circt", &edakit.Bank{Items: mlirItems})
	reg.Register("hdl.ir.circt_dialects", &edakit.Bank{Items: circtDialectItems})
	reg.Register("hdl.ir.firrtl", &edakit.Bank{Items: firrtlItems})

	// Frontier.
	reg.Register("hdl.frontier.rtl_harness_loop", &edakit.Bank{Items: harnessItems})
	reg.Register("hdl.frontier.refactor_equivalence", &edakit.Bank{Items: refactorItems})
	reg.Register("hdl.frontier.rtl_reward", &edakit.Bank{Items: rewardItems})
	reg.Register("hdl.frontier.llm_friendly_ir", &edakit.Bank{Items: llmIRItems})
	reg.Register("hdl.frontier.agent_editable_compiler", &edakit.Bank{Items: agentCompilerItems})
}

// code wraps a snippet so that it sits on its own lines inside a stem.
func code(lines ...string) string {
	return "\n\n" + strings.Join(lines, "\n") + "\n\n"
}

// letters are short variable names for generated snippets. None is a single capital A to E,
// so an option line can never be mistaken for one.
var varNames = []string{"a", "b", "c", "d", "e", "f", "g", "h"}

// kp resolves the knowledge point, choosing one at random for unguided requests.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}
