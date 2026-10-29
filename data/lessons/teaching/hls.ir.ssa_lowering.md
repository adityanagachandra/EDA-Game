# Lowering C to a CDFG: SSA, memories and widths

Before any scheduling, the HLS front end turns C into a clean intermediate form: static single assignment, explicit memories for arrays, known bit widths and simplified arithmetic. These lowering decisions fix the size of every unit and the dependences the scheduler must respect, so a tool designer must know what each one buys and what it can break.

## SSA form and phi nodes

In static single assignment form every variable is assigned exactly once in the program text: reassignments get fresh names $x_1, x_2, \dots$ At a join of control flow, where different paths define different versions, SSA inserts a phi node, $x_3 = \phi(x_1, x_2)$, which picks the version from the edge actually taken. SSA makes data-flow graphs trivial to build, because each use names one definition, so def-use edges are explicit. It also maps cleanly onto hardware. A phi at a join that has been if-converted becomes a multiplexer selected by the branch condition, and a phi at a loop header becomes the register that carries a value from one iteration into the next, with a multiplexer choosing between the initial value and the value fed back.

## Arrays, pointers, calls and widths

Each C construct gets a hardware meaning. A large local array becomes an on-chip RAM with one or two ports, not a bank of registers, so the number of accesses per cycle is limited; small arrays may be split into registers by partitioning. Pointers that may alias force the tool to keep their loads and stores in program order, since a store through one may change what a load through the other reads; restrict qualifiers or separate arrays remove that ordering. Function calls are usually inlined so the callee's operations join the caller's DFG and can overlap with it, although a large callee may be kept as a shared sub-module. Bit width is a first-class cost: arbitrary-precision types such as a 12-bit integer make 12-bit adders and registers instead of 32-bit ones, and the area and delay of a multiplier grow quickly with its operand widths.

## Optimizations that reshape the graph

Several compiler passes change the graph the scheduler sees. Constant propagation and strength reduction turn $y \times 8$ into a shift by 3 bits, which is only wiring. Common-subexpression elimination removes repeated computations. Tree-height reduction rebalances associative chains: $((a + b) + c) + d$ has depth 3, while $(a + b) + (c + d)$ has depth 2 with the same three additions, shortening the critical path. On integers this is exact, because two's-complement addition is associative even with wraparound. On floating point it is not, since every addition rounds, so tools reassociate floating-point only when the user permits it. Bitwidth analysis proves ranges: a counter that never exceeds 99 needs $\lceil \log_2 100 \rceil = 7$ bits. Each pass is small, but together they decide the operation count and depth that bound every schedule.

## Going deeper

Cytron et al., "Efficiently computing static single assignment form and the control dependence graph" (ACM TOPLAS, 1991); Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), chapters on arrays and bit widths.
