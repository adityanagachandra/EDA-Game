# Change tiers for a synthesis engine

A synthesis engine fits the four-layer model cleanly. The AIG with structural hashing and the contract that every pass preserves function is the core. Scripts and pass selection are strategy, numbers such as cut size and iteration counts are parameters, and per-pass statistics are observation hooks. Classifying a customer request into the smallest safe tier is the first decision of an agentic change loop.

## The layers of a synthesis engine

The immutable core is the AIG data structure, structural hashing, the transformation contract (every pass preserves function) and the equivalence checker. Strategy is the choice and order of passes, the recipe, plus pluggable policies such as which cones to balance. Parameters are numbers: cut size, number of rewriting iterations, whether to accept zero-gain moves, level preservation. Observation hooks report node count, levels and estimated mapped QoR after each pass. A request to use 6-input cuts instead of 4-input cuts in rewriting is Tier 1 if the core supports it; a request to make AIG nodes carry a new attribute through every pass is Tier 3.

## Seams for customer-specific behavior

The behavior that varies most between customers is the cost judgment inside passes: whether a replacement is better. The right seam is a cost provider that the rewriter calls for each candidate replacement, given the cut, MFFC size, the candidate's node count and levels, and the node's timing criticality. A customer who penalizes XOR-heavy structures, or who weights depth on critical cones only, then writes a Tier 2 plug-in. The seam should be called per batch of candidates, not per node visit, so the core's hot loops stay inlined. FSM encoding heuristics and rule sets for datapath rewriting are similar Tier 2 extension points.

## Proof and promotion

A Tier 2 plug-in passes the golden tests (equivalence on every benchmark) and a differential QoR test showing that results changed only on designs where the new cost should matter, by the intended amount. Tier 3 changes, even pure speedups of the hash table, need a human-reviewed diff and full regression, because every pass depends on the core. If customer requests keep landing in Tier 3, the seam is in the wrong place. A recipe or plug-in that works for one customer is validated, promoted into the library with guidance on when to select it, and the customer's design stays private.

## Going deeper

The four-layer engine model and change tiers in this curriculum's agentic domain; Brayton and Mishchenko, "ABC: an academic industrial-strength verification tool", CAV (2010).
