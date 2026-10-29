# Constant propagation and folding

Constant propagation replaces values known at compile time by their constants and folds the operations that use them. In hardware compilers it is the most productive single optimization: tied-off inputs, parameter-dependent logic and configuration bits all become constants, and folding them deletes logic. Yosys runs it as opt_expr; MLIR runs it as folding during canonicalization.

## Folding forward through SSA

In SSA form each value has one definition, so propagation is a single forward pass: if both operands of an operation are known constants, evaluate it at compile time and record the result as a constant for every later use. With x an unknown input, t1 = 1 + 2; t2 = t1 * 3; t3 = t2 + x; gives t1 = 3, then t2 = 9, while t3 stays an operation because x is unknown; it becomes t3 = 9 + x.

In hardware the effect is concrete. A constant operand turns an adder into a narrower or simpler circuit, a constant multiplexer select removes the multiplexer and one of its input cones, and a constant register input with a matching reset value removes the register. Synthesis flows run folding repeatedly, because each round can expose new constants.

## Algebraic identities

Plain folding needs every operand constant, but some identities produce constants from unknown values: $v \cdot 0 = 0$, $v \wedge 0 = 0$ (bitwise AND), $v \vee \bar{0} = $ all ones, $v \oplus v = 0$. Others simplify without producing constants: $v + 0 = v$, $v \cdot 1 = v$. A compiler applies them as rewrite rules alongside folding.

So with t1 = 3 * 2; t2 = x + 3; t3 = t1 + t1; t4 = x * 0; t5 = t1 - 2;, the constants are t1 = 6, t3 = 12, t4 = 0 (by the identity) and t5 = 4: four of five temporaries, while t2 depends on x. Yosys's opt_expr implements a set of such identities over its word-level cells, and also handles x bits, which may be chosen as either value to simplify further.

## Branches and sparse conditional propagation

Control flow limits plain propagation. If k = 9 + 3; and then if (k > 9) y = k * 3; else y = x + 1;, the phi for y at the join meets a constant (36) and an unknown, so plain propagation calls y unknown. But k folds to 12, so the condition is a constant true and the else arm can never run.

Sparse conditional constant propagation (SCCP), due to Wegman and Zadeck, tracks reachability and constants together. It marks only executable CFG edges, evaluates a branch condition to decide which successors become executable, and lets a phi meet only the values from executable predecessors. Here y = 36 and a later z = y + 1 folds to 37. The lattice has three levels per value (not yet known, a constant, not a constant), and the analysis runs to a fixed point exactly like other data-flow analyses.

## Going deeper

Wegman and Zadeck, "Constant Propagation with Conditional Branches" (TOPLAS 1991); the Yosys manual, opt_expr.
