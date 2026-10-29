# Synthesis with Yosys and ABC

The standard open-source synthesis flow is Yosys for the front end and word-level optimization and ABC for combinational logic optimization and mapping. Yosys reads Verilog into its RTLIL representation, lowers processes, optimizes, maps memories and arithmetic, and hands each combinational region to ABC. Knowing what each pass does lets a tool designer place a new optimization where its information still exists.

## Front end and process lowering

Yosys parses Verilog into an abstract syntax tree and elaborates it into RTLIL: modules, wires and cells, plus processes for always blocks. The proc family turns processes into netlist cells: decision trees become multiplexers, clocked assignments become flip-flop cells (with asynchronous reset where present), and incompletely assigned combinational variables become latches. The opt family then simplifies: opt_expr folds constants, opt_merge merges identical cells, opt_muxtree removes unreachable multiplexer branches, and opt_clean deletes unused cells and wires. At this point the design is still word-level, which is where fsm extraction, memory inference and arithmetic grouping (alumacc, share) work.

## Lowering and the hand-off to ABC

techmap replaces word-level cells such as adders, comparators and multipliers with implementations from a map library, ending in Yosys's internal single-bit gates (AND, OR, NOT, XOR, MUX cell types). Flip-flops are mapped to library cells by dfflibmap with a Liberty file. The abc pass then extracts the combinational logic between registers, writes it to ABC, runs an ABC script (strash, rewriting and balancing passes, then mapping to the Liberty cells or to LUTs) and reads the mapped netlist back. Because ABC sees only combinational logic, registers are handled separately, and anything not exported to ABC, such as keep-marked nets, survives as is.

## Where new optimizations belong

Each representation keeps different information. Recognizing a special multiplier, re-encoding an FSM, or sharing an adder between exclusive branches needs word-level RTLIL, before techmap bit-blasts it. Generic Boolean optimization belongs in ABC on the AIG. Library-aware sizing belongs after mapping. FPGA scripts (for example synth_ice40) replace the Liberty mapping with LUT mapping and techmap to device primitives such as carry chains and block RAMs. When QoR changes between versions, the ABC script and ABC version used inside the abc pass are the usual cause, since most combinational optimization happens there.

## Going deeper

Wolf, Glaser and Kepler, "Yosys: a free Verilog synthesis suite", Austrochip (2013); Brayton and Mishchenko, "ABC: an academic industrial-strength verification tool", CAV (2010); the Yosys manual.
