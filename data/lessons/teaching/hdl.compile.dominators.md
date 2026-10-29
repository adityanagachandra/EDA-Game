# Dominators and dominance frontiers

In a control-flow graph, block d dominates block v when every path from the entry to v passes through d. Dominance gives a tree over the CFG and a frontier where each block's influence ends. These are the inputs to SSA construction, loop detection and many optimizations, so every compiler, hardware compilers included, computes them early and reuses them.

## Dominator sets

Block d dominates v, written d dom v, when every path from the entry block to v includes d. Every block dominates itself, and the entry dominates every block. The sets satisfy a fixed-point equation: $\text{Dom}(\text{entry}) = \{\text{entry}\}$ and, for every other block, $\text{Dom}(v) = \{v\} \cup \bigcap_{p \in \text{pred}(v)} \text{Dom}(p)$.

Solve it by iteration: start with every non-entry set equal to all blocks, then repeatedly apply the equation until nothing changes. With edges 1→2, 1→3, 1→4, 2→4, 3→5, 4→5, 4→6, block 4 has predecessors 1 and 2, so $\text{Dom}(4) = \{4\} \cup (\{1\} \cap \{1, 2\}) = \{1, 4\}$, and block 6, whose only predecessor is 4, has $\text{Dom}(6) = \{1, 4, 6\}$.

## Immediate dominators and the tree

The strict dominators of v (its dominators other than itself) are totally ordered along every path from the entry. The closest one is the immediate dominator, idom(v): the strict dominator that every other strict dominator of v also dominates. In a set such as $\text{Dom}(3) = \{1, 2, 3\}$, the immediate dominator is 2, the strict dominator with the largest dominator set of its own.

Linking each block to its immediate dominator gives the dominator tree, rooted at the entry. Practical compilers compute it directly with the Lengauer-Tarjan algorithm or with the simpler iterative algorithm of Cooper, Harvey and Kennedy, which runs over blocks in reverse postorder and is fast on real CFGs. A back edge is an edge whose target dominates its source, which is how natural loops are found.

## Dominance frontiers

The dominance frontier of v, DF(v), is the set of blocks w where v dominates a predecessor of w but does not strictly dominate w itself: the first places where v's dominance ends. With edges 1→2, 2→3, 2→4, 3→5, 4→5, block 3 dominates the predecessor 3 of block 5 but not 5, so DF(3) = {5}; likewise DF(4) = {5}, and DF(2) is empty because 2 dominates 5.

A simple algorithm finds frontiers from joins: for every block b with two or more predecessors, walk up the dominator tree from each predecessor until reaching idom(b), adding b to the frontier of every block passed. A block on a loop can be in its own frontier, through the back edge. Frontiers are exactly where SSA construction must place phi functions, which is why they matter.

## Going deeper

Cooper, Harvey and Kennedy, "A Simple, Fast Dominance Algorithm" (2001); Lengauer and Tarjan, "A Fast Algorithm for Finding Dominators in a Flowgraph" (TOPLAS 1979).
