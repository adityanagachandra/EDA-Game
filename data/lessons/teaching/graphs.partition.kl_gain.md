# Kernighan-Lin swap gain

The Kernighan-Lin algorithm of 1970 was the first effective partitioning heuristic, and its gain formula is the template for every move-based partitioner since. KL improves a balanced bisection of a weighted graph by swapping pairs of vertices, choosing each swap by how much it reduces the cut. Computing those gains exactly, including the correction for an edge between the swapped pair, is the core skill.

## External and internal cost

For a vertex $v$ in a bisection, let $E(v)$ be the total weight of its edges to the other side and $I(v)$ the total weight of its edges to its own side. Define $D(v) = E(v) - I(v)$. Moving $v$ alone to the other side would reduce the cut by exactly $D(v)$: its external edges become internal and its internal edges become external.

If $v_6$ has edges of weight 5 and 1 to the other side and edges of weight 2 and 4 to its own side, then $E = 6$, $I = 6$ and $D = 0$: moving it alone would not change the cut. A positive $D$ marks a vertex that is pulled toward the other side.

## The swap gain

KL keeps exact balance by swapping a vertex $a$ from side $A$ with a vertex $b$ from side $B$. The gain is

$$g(a, b) = D(a) + D(b) - 2c(a, b),$$

where $c(a, b)$ is the weight of the edge between them (0 if none). The correction appears because $D(a)$ and $D(b)$ each count the edge $a$-$b$ as external and therefore as removed from the cut, but after the swap the two vertices are again on opposite sides, so the edge stays cut. With $D(a) = 3$, $D(b) = 2$ and $c(a,b) = 2$, the gain is $3 + 2 - 4 = 1$.

## A KL pass

A pass repeatedly picks the unlocked pair with the largest gain, tentatively swaps it, locks both vertices and updates the $D$ values of the others: for a vertex $x$ on $a$'s former side, $D(x)$ increases by $2c(x, a) - 2c(x, b)$, and symmetrically. After $n/2$ swaps the sides have fully exchanged, so KL keeps the prefix of swaps with the largest cumulative gain and undoes the rest.

The best single gain is often negative, and the pass still proceeds, because later swaps may recover more. Finding the best pair naively costs $O(n^2)$ per swap and $O(n^3)$ per pass, which limited KL to small graphs. Its other limits, unit-size vertices, exact balance and graph edges rather than nets, are what Fiduccia and Mattheyses removed.

## Going deeper

Kernighan and Lin, "An efficient heuristic procedure for partitioning graphs" (Bell System Technical Journal, 1970); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 2.
