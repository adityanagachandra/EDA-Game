# Optimality limits of tree and DAG covering

Every mapper claims optimality of some kind. A tool designer must know exactly what is proved and what is not: which problems the DP solves exactly, which are NP-hard, and which heuristics fill the gap. These limits explain the architecture of modern mappers, which compute delay-optimal labels exactly and then recover area heuristically.

## What tree covering guarantees

Tree covering is exact for each tree of the partition, relative to the given subject graph, the pattern set and additive costs. It is linear in the tree size for a fixed library, because each node tries a bounded number of patterns. The guarantee is not global: the partition forces each multi-fanout node to be a gate output, the decomposition fixes which structures exist, and the patterns fix which cells can be recognized. A sum of per-tree optima is an upper bound on the best mapping of the network, not a lower bound, since covers that cross tree boundaries can be cheaper.

## Hard and easy DAG problems

Minimum-area covering of a general DAG is NP-hard: a shared node can be built once or absorbed into several matches, and those choices interact across cones. Delay-optimal covering under a load-independent delay model is polynomial: a node's best arrival depends only on its fanin cone, so a single forward labeling pass computes it everywhere. FlowMap is the LUT instance of this fact, computing depth-optimal labels with max-flow. Load-dependent delay breaks the clean picture, because a node's delay depends on how many consumers the final cover gives it, which is decided later.

## How practical mappers respond

The response is a two-phase recipe. First compute delay-optimal labels exactly, allowing duplication. Then propagate required times from the outputs and, at every node with slack, switch to cheaper matches with heuristics such as area flow and exact local area, never exceeding the required time. Structural choices widen the search, and inserting INV-INV pairs with a zero-cost INV(INV(x)) pattern lets each consumer choose signal polarity freely. When a network has few multi-fanout points, tree covering is already close to the DAG optimum; dense reconvergence is where the cut-based methods win.

## Going deeper

Kukimoto, Brayton and Sawkar, "Delay-optimal technology mapping by DAG covering", DAC 1998; Cong and Ding, "FlowMap", IEEE TCAD 1994; Garey and Johnson, "Computers and Intractability", for the reduction style behind DAG covering hardness.
