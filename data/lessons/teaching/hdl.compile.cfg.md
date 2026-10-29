# Control-flow graphs and basic blocks

A control-flow graph (CFG) represents procedural code as basic blocks connected by possible transfers of control. It is the substrate of every compiler analysis that follows: dominators, SSA construction and data-flow analysis. Hardware compilers need it for the procedural parts of HDLs (always blocks, functions, tasks and testbench code) before those parts become multiplexers and registers.

## Basic blocks

A basic block is a maximal sequence of statements that execute together: control enters only at the first statement and leaves only after the last. A new block starts at the target of every branch and at every join point. A condition such as if (a > 2) is the last statement of the block that evaluates it, because control may go two ways after it.

Count blocks structurally. Straight-line code is one block. An if-then adds two: the then block and the join block that follows. An if-else adds three: two arms and the join. A while loop adds three: a header that evaluates the condition, the body and the exit block after the loop. For d = a + 1; a = a + 1; if (i < 6) begin s = c + 4; end c = d + 5;, there are three blocks: the entry block ending with the test, the then block, and the join block.

## Edges

An edge goes from block b to block c when control can pass directly from the end of b to the start of c. An if-then has three edges: test block to then block, test block to join (the false path) and then block to join. An if-else has four: two out of the test and one from each arm into the join. A while loop has four: into the header, header to body, the back edge from the end of the body to the header, and header to exit.

The back edge is what makes a loop: it points to a block that dominates its source. Nested constructs compose, because the inner construct's join or exit block becomes the current block of the outer one. In the CFG of an if-else whose then arm contains an if-then, there are six blocks and seven edges.

## Cyclomatic complexity

For a connected CFG with $N$ blocks and $E$ edges, the cyclomatic complexity is $M = E - N + 2$, the number of linearly independent paths through the code. For structured code with one entry and one exit, it equals the number of decisions plus one: every if and every while adds one. Six blocks and seven edges give $M = 3$, matching two decisions.

The number has practical uses: it bounds the number of test paths needed for branch coverage and flags overly complex procedures. In hardware, procedural complexity also predicts the size of the multiplexer trees that synthesis builds from a combinational block. Compilers build the CFG once and run every later analysis on it, so getting block boundaries right is the foundation of everything that follows.

## Going deeper

Aho, Lam, Sethi and Ullman, "Compilers", section 8.4 (basic blocks and flow graphs); McCabe, "A Complexity Measure" (IEEE Transactions on Software Engineering, 1976).
