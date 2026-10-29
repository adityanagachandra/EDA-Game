# Incremental computation and dependency graphs

Every EDA tool caches derived results: parasitics computed from wires, delays from parasitics, arrivals from delays, slacks from arrivals. When the design changes, recomputing everything is too slow, and recomputing too little gives wrong answers. A dependency graph makes the cache honest. Each derived result is a node, each edge says which input it was computed from, and an edit tells the tool exactly which cached values may now be stale. This structure sits underneath incremental timing, incremental legalization and incremental builds alike.

## Invalidation follows forward edges

Draw an edge $u \to v$ whenever result $v$ is computed from $u$. The graph is acyclic, because no result can depend on itself, so it has a topological order. When an edit modifies node $u$, every node reachable from $u$ along edges may be stale, and no other node can be: its inputs are untouched. The invalid set is therefore the forward reachable set of the edited nodes, found by one depth-first or breadth-first search.

Example: edges $A \to B$, $A \to C$, $B \to D$, $C \to D$, $D \to E$ and $F \to E$. An edit to $B$ invalidates $D$ and $E$, two results; $A$, $C$ and $F$ keep their cached values. With several edits, the invalid set is the union of their reachable sets, so a result reachable from two edits is invalidated once, not twice.

## A query needs only its dirty ancestors

Invalidation is cheap; recomputation is the expensive part, and a lazy tool postpones it until someone asks. A query for node $q$ must bring $q$ up to date, and that requires every input of $q$ to be valid first, recursively. So the work for one query is the set of dirty nodes that are $q$ itself or ancestors of $q$. Dirty nodes elsewhere stay dirty.

Continuing the example, edit $B$ and then query $D$. The dirty set is $\{D, E\}$. The ancestors of $D$ are $A$, $B$ and $C$, all clean, so only $D$ is recomputed: one evaluation, although two results are dirty. A query for $E$ would recompute $D$ and then $E$. This is why a timer asked for the slack of one endpoint touches only that endpoint's fanin cone.

## Recompute in levels

When everything dirty must be refreshed, the topological order says what may run together. Give each dirty node a level: one more than the highest level among its dirty inputs, and 1 if all its inputs are clean. Nodes on the same level do not depend on one another, so a parallel engine can recompute a whole level at once, and the number of rounds equals the deepest level.

Example: edit $A$ in the graph above. Dirty nodes are $B$, $C$, $D$ and $E$. $B$ and $C$ have only the clean input $A$, so both are level 1. $D$ depends on $B$ and $C$, so level 2. $E$ depends on $D$ and the clean $F$, so level 3. Three rounds, the longest dirty path, however many workers there are. Long dependency chains limit incremental parallelism just as they limit full recomputation.

## Going deeper

Acar, "Self-Adjusting Computation" (PhD thesis, Carnegie Mellon, 2005); the Salsa and Adapton papers on demand-driven incremental computation; the OpenSTA source for how a timer tracks invalid arrivals.
