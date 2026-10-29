# Move sets and cost functions as Tier 2 policies

Search-based optimizers, annealers, local search, rip-up and reroute, are defined by their move sets and cost functions. Those are policies, the second layer of the four-layer engine model, and they belong behind narrow extension points so that a new move or cost form is a Tier 2 plug-in, not an edit to the core. Designing that seam, and the evidence a new policy must produce, is how an optimization engine becomes safely customizable by agents.

## Where the policy seam goes

In a simulated-annealing placer the move generator is the natural seam. Moves vary by design style, such as shifting a whole aligned column of datapath cells, while the acceptance rule, cost bookkeeping and legality checks stay in the core. A new move is then Tier 2: a class registered with a factory behind the move-set interface.

The seam must expose what a plug-in needs and nothing more: read access to current positions, connectivity and local cost deltas, and a typed way to return proposed moves. It must not expose the legalizer's internals. Calls should be coarse, for example one call that proposes a batch of moves, so the core's inner loop stays compiled and inlined; a virtual call per move in the hottest loop costs real performance. Plain data structures, no hidden global state and seedable randomness keep plug-ins deterministic.

## What the core still guarantees

Policies decide what to try; the core decides what is allowed. A plug-in move that swaps cells only within a timing-critical path group still passes through the core, which applies the move, updates costs and guarantees a legal, overlap-free placement whatever the plug-in proposed. That division is what makes Tier 2 changes safe to accept from an agent.

Each new policy must produce evidence. A new rip-up-and-reroute ordering is proven by DRC-clean golden designs plus a QoR comparison against the baseline. A new router cost function must also show determinism: the same inputs and seed reproduce the same result, or differential testing cannot attribute a QoR change to the plug-in. A plug-in that keeps hidden state across calls, such as a global variable remembering its last move, breaks that reproducibility and should be rejected by review.

## Selecting and promoting policies

An agent choosing among five compiled move-set policies reads their manifests, short statements of what each targets and on which design characteristics it was validated, rather than their source. Selection is a cheap text decision; the compute stays in the engines.

A policy found for one design enters the skill flywheel: iterate as a plug-in, validate on golden designs, then promote to the library with a manifest that records where it helps. A move that gains 2% on one customer's datapath design and nothing elsewhere is promoted with that condition, not made the default. Tier placement must be honest too. Raising the probability of an existing operand-swap move is Tier 1, a parameter. Editing the annealer's acceptance and legality code to add a move inline would wrongly turn a Tier 2 request into a Tier 3 core edit with full regression cost.
