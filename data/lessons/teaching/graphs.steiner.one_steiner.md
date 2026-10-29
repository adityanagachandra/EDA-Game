# The iterated 1-Steiner heuristic

Iterated 1-Steiner, due to Kahng and Robins, improves a minimum spanning tree by adding Steiner points one at a time, each time the Hanan point that shortens the MST most. It is simple, uses only MST computations, and produces trees within a few percent of optimal on typical nets. It is a model for a whole class of greedy improvement heuristics in EDA.

## The gain of one Steiner point

For a pin set $P$ and a candidate point $s$, the 1-Steiner gain is $\Delta(s) = \text{MST}(P) - \text{MST}(P \cup \{s\})$. A positive gain means adding $s$ as a junction shortens the tree. The candidate set is the non-pin Hanan points, by Hanan's theorem.

For five pins with MST 27, adding the point $(4, 8)$ might give MST 23, a gain of 4. A point outside the useful region gives gain 0 or less: if it joins the tree as a leaf, the tree only gets longer. Computing one gain costs one MST, $O(n^2)$ with Prim's algorithm.

## Choosing the best point

One round evaluates the gain of every candidate Hanan point and picks the largest. With $O(n^2)$ candidates and $O(n^2)$ per MST, a naive round costs $O(n^4)$; incremental MST updates bring this down substantially, and practical implementations evaluate many candidates in parallel.

For five pins with MST 22, trying all candidates might show that $(5, 5)$ is best, lowering the tree to 18. The best single point already captures much of the saving available, because for three pins one point is all an optimal tree ever needs, and many multi-pin nets behave like groups of near-three-pin clusters.

## Iterating to a fixed point

The iterated version repeats: add the best point if its gain is positive, recompute, and stop when no candidate helps. Kahng and Robins also remove Steiner points whose degree drops to 2 or less, since such points no longer save length. The result has at most $n - 2$ useful Steiner points.

On a four-pin net with MST 13, one round might add a single point and reach 12, which is the exact RSMT. On random nets the heuristic averages within about one percent of optimal. Its greedy nature is the lesson: a locally best addition is not always part of an optimal tree, so it can stop short, but the gains are cheap to evaluate and every step is a valid tree, so it degrades gracefully.

## Going deeper

Kahng and Robins, "A new class of iterative Steiner tree heuristics with good performance" (IEEE TCAD, 1992); Kahng and Robins, "On Optimal Interconnections for VLSI" (1995).
