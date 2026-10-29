# Dead-code elimination

Dead-code elimination (DCE) removes computations whose results cannot affect any output. In software that means instructions; in hardware it means cells, wires and registers left behind by constant folding, by unused outputs of IP blocks and by debug logic. Yosys runs it as opt_clean after almost every pass, and getting it right on cyclic hardware requires the mark-and-sweep form.

## Marking from the outputs

Mark-and-sweep DCE starts from what must be kept, the outputs (and any object marked keep), and walks backward through operands, marking each definition reached. Everything left unmarked is swept: it cannot influence an output on any path. With outputs depending on t6, which reads t4 and t5, which read t2 and t1, a temporary t3 that nothing marked reads is removed.

The walk is a graph search over use-def edges, linear in the size of the program. In a netlist the same search runs over nets and cells: start at the module's output ports and kept objects, and traverse from each cell's outputs back to its inputs. Unmarked cells and the wires that only they drive are deleted.

## Transitive dead code

Removing one dead statement can make others dead. If t3 = b & t1; and t4 = t2 + b; are both unused, removing them leaves t2 still used by t5, but in another program removing t4 might leave its operand t2 unused as well. A simple pass that deletes definitions with zero uses must therefore iterate until nothing changes; the marking approach gets the whole transitive set in one pass, because anything not reachable from an output is never marked.

This interacts with constant propagation. Folding a multiplexer select to a constant disconnects one input cone, which then becomes dead; folding a register to a constant disconnects its input logic. That is why synthesis scripts interleave opt_expr (folding) and opt_clean (sweeping) until a fixed point.

## Dead cycles in hardware

Hardware adds a case software DCE rarely meets: dead cycles. Two registers r1 <= r2 + a; and r2 <= r1 ^ b; that feed only each other each have a use, so a pass that removes only definitions with zero uses keeps both forever. Marking from the outputs never reaches them, so mark-and-sweep removes them correctly.

The opposite risk is removing logic that matters for reasons the IR does not show: a register observed only by a debug probe, a redundant flip-flop deliberately added for reliability, logic used by a formal property. Tools protect such objects with keep attributes or by treating properties and probes as outputs. A DCE implementation must define its roots precisely, because the roots, not the algorithm, decide what survives.

## Going deeper

Cooper and Torczon, "Engineering a Compiler", chapter 10; the Yosys manual, opt_clean.
