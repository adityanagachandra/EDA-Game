# Data-flow analysis

Data-flow analysis computes facts that hold at every point of a program, such as which definitions may reach a point or which variables may still be read later, by solving equations over the control-flow graph to a fixed point. It is the engine behind constant propagation, dead-code elimination, register inference and latch detection, and its theory (lattices and monotone transfer functions) explains why iteration terminates.

## Reaching definitions

A definition d reaches a point p if some path from d to p does not redefine d's variable. For each block, GEN is the last definition of each variable in the block, and KILL is every other definition of those variables. The equations are $\text{IN}[b] = \bigcup_{p \in \text{pred}(b)} \text{OUT}[p]$ and $\text{OUT}[b] = \text{GEN}[b] \cup (\text{IN}[b] - \text{KILL}[b])$, with IN of the entry empty. It is a forward may-analysis: facts flow along edges, and a union merges them.

Take blocks 1 (d1: b), 2 (d2: c, then d3: b), 3 (d4: b) and 4 (d5: b), with edges 1→2, 2→3, 2→4, 4→2. Iterating gives OUT[1] = {d1}, OUT[2] = {d2, d3}, OUT[4] = {d2, d5}, so IN[2] = {d1, d2, d5}. This is the analysis that decides whether a variable in a process is assigned on every path, which is exactly latch detection.

## Liveness

A variable is live at a point if some path from that point reads it before writing it. Liveness runs backward: $\text{OUT}[b] = \bigcup_{s \in \text{succ}(b)} \text{IN}[s]$ and $\text{IN}[b] = \text{USE}[b] \cup (\text{OUT}[b] - \text{DEF}[b])$, where USE is the set of upward-exposed reads, read before any write in the block, and DEF the set of variables written. Nothing is live after the exit.

Liveness answers practical questions in hardware compilers. A blocking temporary in a clocked block that is live on entry to the block is read before it is written, so it needs a flip-flop. A phi whose variable is dead at a join can be pruned. And a register whose value is never live anywhere can be removed.

## Lattices and termination

Why does iteration stop? The facts at each point form a lattice, here the subsets of a finite set ordered by inclusion, and the transfer functions are monotone: a larger input never yields a smaller output. Starting from the empty set and applying monotone functions, each set can only grow, and it can grow only finitely often, so a fixed point is reached; it is the least solution of the equations.

The same framework covers many analyses by changing the lattice, the direction and the merge. Available expressions is a forward must-analysis with intersection; constant propagation uses a lattice of unknown, a constant value and not-a-constant. Visiting blocks in a good order (reverse postorder for forward problems) makes most real CFGs converge in a few passes, and loops are what require more than one. A worklist algorithm revisits only the blocks whose inputs changed.

## Going deeper

Aho, Lam, Sethi and Ullman, "Compilers", chapter 9; Nielson, Nielson and Hankin, "Principles of Program Analysis", chapter 2; Kildall, "A Unified Approach to Global Program Optimization" (POPL 1973).
