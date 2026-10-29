# Yosys opt passes and flatten

Yosys's opt command is a fixed-point loop of netlist optimizations that implement, on RTLIL cells, the classical compiler optimizations: constant folding, common subexpression elimination and dead-code elimination, plus hardware-specific passes for multiplexer trees and flip-flops. Flatten decides how far those optimizations can reach. Knowing which pass does what lets you read a synthesis log and design a flow.

## Folding, merging and cleaning

Three passes map directly onto compiler optimizations. opt_expr performs constant folding and simple algebraic rewrites: a cell with constant inputs is evaluated, and patterns such as AND with 0 or a multiplexer with equal inputs are simplified. opt_merge is common subexpression elimination on the netlist: cells with the same type, the same parameters and the same input signals compute the same value, so one is kept and drives both sets of loads. opt_clean is dead-code elimination: it removes cells and wires that drive nothing used, keeping ports and objects with the keep attribute.

Each pass exposes work for the others. Folding can make two cells identical, merging can leave a cell without loads, and cleaning can simplify fan-in so that more folding applies. The opt command therefore repeats its sub-passes until none of them changes the design.

## Multiplexers and flip-flops

Hardware adds passes with no direct software analogue. opt_muxtree removes multiplexer inputs that can never be selected: in a tree of multiplexers sharing a select signal, an inner multiplexer on the same select along one branch has an unreachable input. opt_reduce simplifies reduce and multiplexer cells, for example removing duplicate inputs of a wide OR.

opt_dff works on flip-flops. A flip-flop whose D input is a constant equal to its reset or initial value can only ever hold that value, so it becomes a constant. A feedback multiplexer in front of a flip-flop, D = en ? x : Q, is folded into a flip-flop with an enable, and synchronous reset logic is folded into a flip-flop with a synchronous reset, matching cells that libraries and FPGAs provide directly.

## Flattening and pass ordering

Optimization stops at module boundaries: a constant tied to a submodule's input cannot propagate inside while the boundary exists. flatten inlines submodules into their parents so that constants and duplicates become visible across former boundaries; modules marked keep_hierarchy, and black boxes, stay as instances. After flattening, two copies of one submodule tied to different constants are optimized differently, which is exactly the divergence that uniquification would otherwise have to handle.

Script order matters. A typical flow reads the design, checks the hierarchy, runs proc, optionally flattens, runs opt, extracts FSMs and memories, runs opt again, lowers with techmap, runs opt once more and hands the result to ABC. Running opt only once, before techmap, would leave the constant and duplicate gates that techmap creates when it expands word-level cells.

## Going deeper

The Yosys manual, the opt, opt_expr, opt_merge, opt_clean, opt_muxtree, opt_dff and flatten commands, and the synth command's script.
