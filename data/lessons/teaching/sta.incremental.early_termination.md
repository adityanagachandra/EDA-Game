# Early termination in incremental propagation

The fanout cone of a change is an upper bound on what must be recomputed, not the actual amount. Arrival propagation uses a maximum, and a maximum often hides a change: if the changed input was not the latest one, or still is not after the change, the output does not move. An incremental timer exploits this by stopping propagation at every node whose value comes out unchanged.

## The new arrival at the changed node

Start at the head of the changed arc. Its new arrival is the maximum over its fanins with the new delay: $AT'(v) = \max\big(AT(u) + d'(u,v), \max_{w \ne u} (AT(w) + d(w,v))\big)$. If arc $J \to C$ drops from 30 to 25 ps with $AT(J) = 63$ ps and $C$ has no other fanin, $AT'(C) = 88$ ps, down from 93. If instead $C$ had another fanin arriving at 91 ps, its arrival would drop only to 91 ps, and if the changed arc was not the maximum before and did not become it after, $C$ would not change at all.

## Counting the nodes that actually change

Process the cone in level order. At each node, recompute the arrival from its fanins; if it equals the old value, do not schedule its fanouts on its account. Only nodes whose arrival really changed push work downstream. In a typical example, a change at $T$ alters $T$'s arrival, but at the next node another fanin dominates and the change stops there, so one node changes out of the two in the cone. For speed-ups, which lower delays, masking is very common, because a node's latest fanin rarely belongs to a small edit. Large designs see big savings from this, since optimizers mostly make small changes far from the critical paths.

## Arrivals at downstream endpoints

The endpoint result follows from the same rule. If every path from the change to an endpoint passes through a node where the change is masked, the endpoint's arrival is unchanged even though it lies in the cone. For instance, a 5 ps speed-up on $M \to K$ moves $K$ from 80 to 75 ps, but endpoint $J$ takes its arrival from $N$ at 82 ps plus 23 ps, so $J$ stays at 105 ps. This is also the reason an ECO that speeds up a gate off the critical path improves nothing: the max at the merge point absorbs it. Tool designers use the same masking to prune which endpoints an optimizer must re-examine.

## Going deeper

The OpenSTA source, Search and its arrival invalidation; Sapatnekar, "Timing" (2004), the chapter on incremental timing.
