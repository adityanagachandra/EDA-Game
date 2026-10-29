# Data-flow graphs

A data-flow graph (DFG) is the representation every HLS scheduler works on: one node per operation, one edge per value passed from a producer to a consumer. It is built from straight-line code, a basic block, and it says exactly which operations must wait for which. Everything later in the flow, schedule length, mobility, register count, reads off this graph, so a tool designer must be able to build it and measure it by eye.

## Nodes, edges and three-address code

Lower the code to three-address form, where each statement performs one operation on at most two operands: $n_1 = a + b$, $n_2 = n_1 \times c$, $n_3 = n_1 + n_2$. Each statement is a node, labelled with its operation type. Each operand that names an earlier result is an edge from that producer to this consumer; operands that are primary inputs, here $a$, $b$, $c$, are sources outside the graph and add no operation-to-operation edge. In the example, $n_2$ reads $n_1$ (one edge) and $n_3$ reads $n_1$ and $n_2$ (two edges), so the DFG has three operations and three data edges. Because SSA gives every value one definition, the edges can be read directly from the operand names. A DFG of a basic block is acyclic: a value cannot depend on itself inside one execution of the block.

## Depth: the longest dependence chain

The depth of a DFG is the largest number of operations on any path from a source to a sink. With every operation taking one control step, depth is the minimum schedule length even with unlimited hardware, because each operation on the chain must wait for the previous one. Compute it by processing nodes in topological order: the level of a node is one plus the largest level of its predecessors, and sources have level 1. For $n_1 = a + b$, $n_2 = c + d$, $n_3 = n_1 \times n_2$, $n_4 = n_3 + a$, $n_5 = n_2 + b$, the levels are 1, 1, 2, 3, 2, so the depth is 3. This is the same longest-path computation static timing analysis uses, with operation counts in place of delays.

## Sinks, sources and what the graph leaves out

Sources are operations that read only primary inputs; sinks are operations whose result no later operation reads. A sink's value must leave the block, as an output port, a store to memory or a value live into the next block, so sinks mark where results are consumed by the outside. Count them by checking which operation names never appear as an operand. The DFG deliberately omits several things: control flow (it cannot say that only one of two branches runs), the clock (it has no cycles yet), and memory ordering beyond explicit edges (two accesses to the same array need an extra dependence edge if one is a store). Width and operation type ride along as node attributes, which later decide which unit can execute the node and how long it takes.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 3 on graph models; Gajski, Dutt, Wu and Lin, "High-Level Synthesis: Introduction to Chip and System Design" (1992).
