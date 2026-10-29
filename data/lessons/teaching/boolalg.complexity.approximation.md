# Approximation algorithms and heuristics

When a problem is NP-hard, a tool still has to return an answer quickly. An approximation algorithm runs in polynomial time and comes with a proven ratio: its cost is never worse than $\rho$ times the optimum. A heuristic runs fast and often does well but carries no such guarantee. EDA uses both: approximation ideas inside test compaction, covering and partitioning, and heuristics such as simulated annealing, local search and greedy rewriting almost everywhere. Knowing what a ratio promises, and what a heuristic does not, is basic engineering judgment.

## The matching-based vertex cover

Scan the edges in order; whenever an edge has neither endpoint in the cover, add both endpoints. The chosen edges form a maximal matching $M$, and the result is a cover of size $2|M|$. On the path with edges (1,2), (2,3), (3,4), the scan takes (1,2), skips (2,3), takes (3,4), and returns 4 vertices.

The proof of the ratio is two lines. The returned set covers every edge, or the scan would have taken that edge. Any cover, including the optimal one, needs a distinct vertex for each matched edge, since the matched edges share no endpoints, so $OPT \ge |M|$ and the algorithm's cost $2|M| \le 2 \cdot OPT$.

## What a ratio guarantees

A ratio is a worst-case bound over all instances, not a prediction for a given one. On the path above the optimum is $\{2, 3\}$, so the ratio achieved is $4/2 = 2$, the bound is tight. On other graphs the same algorithm may happen to be optimal. The matching also certifies a lower bound: $OPT \ge |M|$, so the algorithm can report how far from optimal it might be on this instance.

No polynomial algorithm with a ratio better than 2 is known for vertex cover. For independent set, by contrast, no polynomial algorithm achieves ratio $n^{1-\varepsilon}$ unless P = NP, so its exact complement relationship to vertex cover does not carry approximation over.

## Greedy set cover and heuristics without guarantees

Greedy set cover repeatedly picks the set covering the most uncovered elements. If the largest set has $d$ elements, greedy uses at most $H(d) \cdot OPT$ sets, where $H(d) = 1 + 1/2 + \dots + 1/d \le 1 + \ln d$. For $d = 8$, $H(8) \approx 2.718$: if the optimum is 4 tests, greedy returns at most $\lfloor 4 \times 2.718 \rfloor = 10$. Unless P = NP no polynomial algorithm does better than about $\ln n$ (Dinur and Steurer, 2014), so greedy is essentially the best possible.

Heuristics make no such promise. Simulated annealing for placement, Espresso's expand-irredundant-reduce loop and AIG rewriting scripts often beat any guaranteed algorithm on real designs, but their quality is measured on benchmarks, not proven. A disciplined tool reports the lower bound it knows (a matching, an LP relaxation) next to the heuristic's answer, so the gap is visible.

## Going deeper

Vazirani, "Approximation Algorithms" (2001), chapters 1 and 2. Williamson and Shmoys, "The Design of Approximation Algorithms" (2011). Dinur and Steurer, "Analytical approach to parallel repetition", STOC (2014).
