# Placing phi functions with dominance frontiers

Constructing SSA form requires knowing where to put phi functions. The classical answer, due to Cytron and colleagues, is the iterated dominance frontier of the blocks that assign each variable. It places the minimum number of phis that minimal SSA needs, and it is the algorithm behind SSA construction in many compilers; MLIR expresses the same merges as block arguments.

## One definition block: the dominance frontier

If variable v has a value on entry and is assigned again in block d, the new definition reaches every block that d dominates without competition. The first blocks where it competes with another definition are exactly the blocks of DF(d): there, some predecessors carry d's version and others carry an older one, so a phi is needed.

With edges 1→2, 1→4, 2→3, 3→4, 3→5 and assignments to v in blocks 3 and 5, DF(3) = {4} (block 4 also receives the entry version via 1→4) and DF(5) is empty, since block 5 has a single predecessor and no successors. Minimal SSA places one phi, in block 4. The entry block never contributes, because it dominates everything and its frontier is empty.

## Iterating the frontier

A phi is itself a new definition of the variable. Placing a phi in block w can therefore require more phis at DF(w), and so on. The set of phi blocks is the iterated dominance frontier $\text{DF}^+(S)$ of the definition set S: the limit of $\text{DF}_1 = \text{DF}(S)$, $\text{DF}_{i+1} = \text{DF}(S \cup \text{DF}_i)$. A worklist computes it: start with the definition blocks, and whenever a block gains a phi, add it to the worklist so that its frontier is processed too.

Loops are where iteration shows. With a back edge 3→2 and definitions in blocks 2 and 4, the back edge puts block 2 in DF(2) and DF(3), and the join at block 4 is in DF(2) as well, so phis go to blocks 2 and 4. The loop header needs a phi because it merges the entry version with the version coming around the loop.

## Several variables and pruning

Phi placement is done per variable, each with its own definition blocks, and the totals add. If u is assigned in blocks 3 and 7 and v in blocks 4 and 5 of the same CFG, their phi sets are computed independently, and the total count is the sum of the two set sizes. The dominance frontiers are computed once and shared, which is why the algorithm is efficient: roughly linear in the size of the program plus the number of phis placed.

Minimal SSA can place phis for variables that are dead at the join. Pruned SSA filters the iterated frontier by liveness, and semi-pruned SSA skips variables that are never live across a block boundary. MLIR represents the result differently but equivalently: instead of a phi at the top of block w, block w takes a block argument and each branch into w passes the value, which makes the edge-to-value pairing explicit.

## Going deeper

Cytron, Ferrante, Rosen, Wegman and Zadeck, "Efficiently Computing Static Single Assignment Form and the Control Dependence Graph" (TOPLAS 1991); Appel, "Modern Compiler Implementation", chapter 19.
