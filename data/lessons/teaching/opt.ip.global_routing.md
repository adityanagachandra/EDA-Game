# Global routing as an integer program

Global routing assigns each net a route through a coarse grid of routing regions so that no edge carries more wires than its capacity and total wirelength stays small. Written as an integer program, each net picks one candidate route, and edge capacities couple the nets. Routers built on this model solve its LP relaxation and round, and the same model underlies negotiated-congestion routers that price overflow. It is the cleanest example of relaxation plus rounding in physical design.

## Usage, capacity and overflow

Give each net n a set of candidate routes, and a binary $x_{n,r}$ that selects route r. Each net takes one: $\sum_r x_{n,r} = 1$. Edge e has usage $u_e = \sum_{(n,r) \ni e} x_{n,r}$ and capacity $c_e$. The overflow is $\sum_e \max(0, u_e - c_e)$.

Example: capacities e1 = 1, e2 = 1, e3 = 2, e4 = 1, e5 = 1. Net 1 takes e1, e2, e4; net 2 takes e2, e4, e5; net 3 takes e4. Usage: e1 = 1, e2 = 2, e3 = 0, e4 = 3, e5 = 1. Overflow: e2 has 1 excess and e4 has 2, total 3. Routers report total overflow and the maximum overflow per edge, because a few badly overflowed edges are much harder to fix in detailed routing than many slightly full ones.

## The 0-1 routing program

The exact model minimizes total wirelength $\sum_{n,r} \ell_{n,r} x_{n,r}$ subject to one route per net and $u_e \le c_e$ on every edge. With 3 nets and 2 candidates each there are $2^3 = 8$ assignments, few enough to check by hand: compute each assignment's usage, discard those that overflow, and keep the shortest.

Real instances have millions of nets and implicit candidate sets (all paths in the grid), so routers generate candidates on demand by shortest-path searches, which is column generation in LP terms. When the capacities cannot all be met, the objective weights overflow heavily and the result tells the placer where to spread cells. Pricing edges by their congestion and rerouting nets repeatedly, as negotiated-congestion routers do, is a heuristic for this same program.

## LP relaxation and randomized rounding

Relax each $x_{n,r}$ to $[0, 1]$. The LP spreads a net over several routes, for example net 1 with 0.2 on route A and 0.8 on route B. Randomized rounding then picks route A with probability 0.2 and route B with 0.8, independently per net.

By linearity of expectation the expected usage of an edge equals its LP usage: the sum of the fractional weights of the routes that cross it. If e2 is crossed by net 1's route A (0.2), net 2's route A (0.1) and net 3's route B (0.3), its expected usage is $0.2 + 0.1 + 0.3 = 0.6$. Because the LP respects capacities, rounding respects them in expectation, and Chernoff bounds show that overflow beyond a logarithmic factor is unlikely, which is the classic analysis of Raghavan and Thompson. In practice routers round, then repair the remaining overflow by rip-up and reroute.

## Going deeper

Raghavan and Thompson, "Randomized rounding: a technique for provably good algorithms and algorithmic proofs" (Combinatorica, 1987); McMurchie and Ebeling, "PathFinder: a negotiation-based performance-driven router for FPGAs" (FPGA, 1995); the OpenROAD documentation for its global router.
