# Residual graphs and augmenting paths

Every max-flow and min-cost-flow algorithm works on the residual graph: the network of capacities that remain after the current flow, including backward arcs that let the algorithm undo flow already sent. Getting residual capacities right, and seeing why backward arcs are necessary, is what separates a correct flow implementation from a greedy one that gets stuck.

## Forward and backward residual arcs

For an arc $u \to v$ with capacity $c$ carrying flow $f$, the residual graph has a forward arc $u \to v$ with residual capacity $c - f$ (how much more can be sent) and a backward arc $v \to u$ with residual capacity $f$ (how much can be cancelled). Saturated arcs disappear forward; arcs with zero flow have no backward arc.

An arc with capacity 10 and flow 7 gives forward residual 3 and backward residual 7. If the network also has an arc $v \to u$ of its own, the two residual capacities add: $c(v,u) - f(v,u) + f(u,v)$. Implementations store each arc and its reverse as a pair, so updating one updates the other in constant time.

## Pushing along an augmenting path

An augmenting path is an $s$-$t$ path in the residual graph. The amount that can be pushed along it is its **bottleneck**, the smallest residual capacity on the path. Pushing increases flow on forward steps and decreases it on backward steps, by the bottleneck.

After a first augmentation sends 9 units along $s \to b \to t$, a path $s \to c \to t$ with residual capacities 4 and 8 has bottleneck 4. A path may also step against an arc: if $a \to b$ carries 3 units, a path can travel $b \to a$ with residual capacity 3. Pushing 2 units along such a step lowers the flow on $a \to b$ from 3 to 1, rerouting flow that was sent earlier, which is how the algorithm corrects early greedy choices.

## Why backward arcs are needed

Without backward arcs, augmenting greedily can block the optimum: a first path may use an arc that two later paths would both need, leaving no augmenting path while the flow is below maximum. Backward arcs let a later path cancel part of the first path's flow and send it elsewhere, so the algorithm reaches the true maximum.

At termination, the set of vertices reachable from $s$ in the residual graph is the source side of a minimum cut. If even the arcs out of $s$ are saturated, that set is just $\{s\}$. The residual graph is also central to min-cost flow, where backward arcs carry negative cost, the saving of undoing a unit, and a negative-cost cycle in the residual graph means the current flow is not cheapest.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), chapter 24; Ahuja, Magnanti and Orlin, "Network Flows: Theory, Algorithms, and Applications" (1993), chapters 6 and 7.
