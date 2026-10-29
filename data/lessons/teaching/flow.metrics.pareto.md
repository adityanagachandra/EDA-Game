# PPA trade-offs and Pareto fronts

There is no single best implementation of a design, because power, performance and area cannot all be minimized at once. Different tool settings, constraints and architectures produce different design points, and comparing them is a multi-objective problem. The Pareto front, the set of points that no other point beats on every metric at once, is the right way to think about it, and design-space exploration tools search for it.

## Dominance

Express every metric so that lower is better: power in mW, clock period or critical delay in ns, area in square millimetres. Design point $A$ dominates point $B$ if $A$ is no worse than $B$ on every metric and strictly better on at least one. Dominance is a partial order: two points can be incomparable, each better on something. Point $A = (12, 1.0, 2.0)$ dominates $B = (14, 1.1, 2.0)$, while $A$ and $C = (10, 1.2, 1.9)$ are incomparable, because $C$ uses less power and area but is slower.

## The Pareto front

A point is Pareto-optimal if no other point dominates it. The Pareto front is the set of all such points. To find it among $n$ points, compare each pair, which costs $O(n^2)$ comparisons, and keep the points that nothing dominates. Every point off the front is a strictly worse choice than some point on it, so it can be discarded no matter what the product needs. Points on the front represent genuine trade-offs, and choosing among them is a product decision, not a tool decision. In practice the front is often plotted two metrics at a time, for example power against delay, to show the team which trades are available.

## Weighted costs and their limits

Optimizers often collapse PPA into one scalar, such as $w_P P + w_D D + w_A A$, and minimize it. Each choice of positive weights picks out a point on the Pareto front, so weighted sums are a practical way to steer a tool. But the chosen point depends strongly on units and weights: changing power from mW to W changes the answer unless the weight changes too. And a weighted sum can reach only points on the convex hull of the front; points in a concave dent of the front are never selected by any weights. Design-space exploration therefore keeps the whole front and lets the team choose.

## Going deeper

Deb, "Multi-Objective Optimization Using Evolutionary Algorithms" (2001), chapter 2 on dominance and Pareto optimality.
