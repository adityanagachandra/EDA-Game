# Levelizing the timing graph

Before a timer propagates anything it levelizes the graph: it assigns every node a level so that each arc goes from a lower level to a higher one. Levels are a topological order with extra structure. They let the timer process the graph as a sequence of independent batches, schedule incremental updates, and run nodes of one level in parallel. OpenSTA's search is built around this levelized order.

## The level of a node

Sources, the startpoints with no fanin, sit at level 0. Every other node's level is one more than the largest level among its fanins: $\ell(v) = 1 + \max_{u \to v} \ell(u)$. Compute it in any topological order. For the arcs $A \to B$, $A \to C$, $B \to C$, $C \to D$, node $A$ is level 0, $B$ is level 1, $C$ is level 2 because its deepest fanin $B$ is at level 1, and $D$ is level 3. Note that $C$ is not at level 1 even though it has a fanin at level 0; the maximum governs. This guarantees that when the timer reaches level $k$, every fanin of every node at that level is already final.

## Depth is the longest path in arcs

The largest level in the graph is its depth: the number of arcs on the longest path from any source, counting arcs rather than delays. It is exactly the longest-path problem on a DAG with unit weights, solved by the same one-pass recurrence. Depth tells a tool designer how many sequential steps a levelized propagation needs, even with unlimited parallel workers. A wide, shallow graph parallelizes well; a deep chain of carry logic does not. Logic depth also reappears in timing itself: advanced on-chip variation tables index their derates by the number of cells on a path. For the four-node example above, the depth is 3.

## Levels as batches of independent work

Nodes on the same level never feed each other, since an arc always climbs at least one level. A timer can therefore evaluate all nodes of one level in parallel, then move to the next, and the number of nodes on a level measures the parallelism available there. Levels also organize incremental timing: when an arc changes, the timer queues the affected nodes by level and processes them in increasing order, so no node is evaluated before its fanins settle and none is evaluated twice. OpenSTA keeps the levels up to date as the netlist changes, relevelizing only the region that moved, because a full relevelization after every edit would waste the savings of incremental analysis.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms", the chapter on topological sort and DAG shortest paths; the OpenSTA source, Levelize and BfsIterator.
