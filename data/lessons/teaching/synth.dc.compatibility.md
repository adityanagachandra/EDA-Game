# Compatible observability don't cares

Each node's full observability don't care is valid only if no other node changes at the same time. Optimizers that simplify many nodes against one snapshot of the network need a weaker but safer set: compatible don't cares, chosen so that every node can use its own freedom simultaneously. The idea goes back to the compatible sets of permissible functions in Muroga's transduction method.

## Why full ODCs conflict

Consider $z = xy$. The input $x$ is unobservable whenever $y = 0$, and $y$ is unobservable whenever $x = 0$. At a point where both are 0, each one's ODC says it may change. If both change to 1, $z$ becomes 1 where it was 0: the circuit is wrong. The same happens at an OR gate when both inputs are 1. Full ODCs are each correct alone, because each assumes the rest of the network is fixed. Using several of them in one pass violates that assumption. Either the don't cares are recomputed after every change, or they are restricted so that they remain valid together.

## Building compatible don't cares

Compatible ODCs order the fanins of each gate. The first fanin keeps its full ODC. Each later fanin gets only the points where it is unobservable no matter how the earlier fanins use their freedom: for $z = xy$ with $x$ first, $x$ may change where $y = 0$, but $y$ may change only where $x = 0$ under every allowed choice of $x$, so the conflicting point is given to $x$ alone. CODCs are therefore a subset of the full ODCs. Their size depends on the ordering, so tools order fanins to give freedom to the nodes that benefit most. The cost is missed simplifications.

## Choosing a strategy

Three practical strategies exist. Change one node at a time and recompute exact don't cares after each change; this needs no compatibility. Compute compatible don't cares once and simplify many nodes together, which suits batch passes. Or verify each change locally with SAT on a window, which is what ABC's mfs does on mapped networks, making compatibility unnecessary. Large designs also force windowing, because exact global don't cares are too expensive. A node with a large full ODC but a nearly empty CODC has given its freedom to fanins processed earlier; reordering can return it.

## Going deeper

Muroga, Kambayashi, Lai and Culliney, "The transduction method: design of logic networks based on permissible functions", IEEE Transactions on Computers (1989); Savoj and Brayton, DAC (1990).
