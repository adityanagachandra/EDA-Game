# Global routing as multicommodity flow

Viewed as network flow, global routing is a multicommodity flow problem: each two-pin connection is a commodity that must send one unit from its source to its sink, and all commodities share edge capacities. Without integrality it is a linear program solvable in polynomial time, and its relaxations give lower bounds and fractional solutions that rounding turns into routes.

## Flow variables

For commodity $k$ (a two-pin connection) and directed arc $(u, v)$, let $f^k_{uv} \ge 0$ be the flow of $k$ on that arc. Each undirected gcell edge gives two arcs, so on a grid with $E$ edges there are $2E$ arcs and $K \cdot 2E$ flow variables for $K$ commodities.

Example: a 5 by 12 gcell grid has $4 \cdot 12 + 5 \cdot 11 = 103$ edges and 206 arcs, so 35 commodities need 7210 variables. Integer flows correspond to routes: if all $f^k$ are 0 or 1 and conserve flow, they trace a path (possibly with cycles that can be removed). The size grows with nets times edges, which is why exact flow formulations are used on small regions or with heavy decomposition.

## Conservation and capacity

Flow conservation holds for every commodity at every vertex: outflow minus inflow equals 1 at the commodity's source, -1 at its sink and 0 elsewhere. That is $K \cdot V$ equations. Capacity couples the commodities: for each edge, the total flow of all commodities in both directions is at most $c_e$.

Example: 28 commodities on a 15 by 3 grid of 45 vertices give 1260 conservation equations. Without the capacity constraints the problem splits into $K$ independent shortest-path problems, one per commodity. All the difficulty, and all the interaction between nets, lives in the capacity rows, which is exactly what Lagrangian relaxation prices with multipliers and what negotiated congestion prices with present and history costs.

## Min-max utilisation

A common variant minimises the worst utilisation $\lambda$: total flow on each edge must be at most $\lambda c_e$. If the optimum has $\lambda \le 1$, a fractional routing fits every capacity; if $\lambda > 1$, no routing fits, even fractionally, so the instance is provably congested. This gives the router a certificate that integral heuristics cannot.

Example: edge flows of 5.5, 7, 12.5 and 4.5 over capacities 6, 4, 7 and 6 give utilisations 0.92, 1.75, 1.79 and 0.75, so $\lambda = 1.79$ for this solution. Approximation schemes for min-max multicommodity flow (in the style of Garg and Koenemann) update edge lengths multiplicatively with utilisation and reroute, which looks very much like negotiated congestion with an exponential present cost.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); N. Garg and J. Koenemann, "Faster and Simpler Algorithms for Multicommodity Flow and Other Fractional Packing Problems" (1998).
