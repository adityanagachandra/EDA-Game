# A partitioner's core versus its change surface

In an agent-first EDA flow, a model edits and configures engines on request, so each engine must be built so that the common requests never touch its correctness-critical code. Applying the four-layer engine model to a multilevel partitioner shows concretely what belongs in the immutable core, what is policy, what is a parameter and what is an observation hook, and which change tier a given request lands in.

## Four layers of a partitioner

The **immutable core** holds what must always be right: the netlist data structures, exact cut and balance bookkeeping, legal moves, projection between levels that preserves the cut, and determinism under a seed. **Strategy** holds the swappable heuristics: the matching or clustering rating, the initial-partitioning method, FM tie-breaking (LIFO or FIFO buckets), and the refinement schedule. **Parameters** are numbers: the imbalance tolerance, the coarsening stop size, the number of V-cycles and initial tries, the seed. **Observation hooks** expose state, such as the cut and balance after each level, to an outside controller.

The test for placing an item: if a mistake in it could produce a wrong cut or an illegal partition, it is core; if it only changes which legal answer is found, it is strategy or a parameter.

## Which tier a request lands in

Requests map to the smallest safe tier. "Use 5 percent imbalance on this block" is **Tier 1**, a parameter change with no recompile. "Cluster each register bank before coarsening" is **Tier 2** if a clustering-policy extension point exists: a small plug-in, with the core untouched. "Weight cut nets by our timing criticality" is Tier 2 through a net-weight seam. "Gains overflow the buckets when weights reach a million, so switch to a heap" is **Tier 3**, a core edit, because the gain data structure is part of the move kernel; it needs a human-approved diff and full regression.

A request that seems to need Tier 3 often reveals a missing seam. If weighting requests kept forcing core edits, the right fix is to add a weight seam once, after which they all become Tier 2.

## Keeping the core small and fast

A small core keeps the expensive verification path rare: every core edit triggers human review and full regression, so the less often it happens, the cheaper the engine is to evolve. But seams must not cost speed. A plug-in called once per pin on every move puts a virtual call in FM's inner loop, which blocks inlining and slows every pass. Seams in a partitioner should be coarse: compute all net weights once per level, return a plain array, and let the kernel read it.

Plug-ins must also obey the core's invariants: no hidden global state, deterministic output for the same input and seed, and no side channels such as wall-clock time. A plug-in that keeps a global counter changing its tie-breaking between runs makes results irreproducible, which defeats the golden and differential tests that make agent-made changes safe to promote.

## Going deeper

The OpenROAD documentation for its partitioning and clustering tools; Schlag et al., "High-quality hypergraph partitioning" (ACM Journal of Experimental Algorithmics), for the configurable components of KaHyPar.
