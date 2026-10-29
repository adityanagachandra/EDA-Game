# Minimum-cost flow

Minimum-cost flow sends a required amount of flow at the least total cost, with a cost per unit on every arc. It generalizes shortest paths (one unit), max flow (zero costs) and assignment, and it solves exactly a surprising number of EDA subproblems: spreading cells out of overfull bins in legalization, assigning pins or tracks, and layer assignment in routing. It is a polynomial-time tool hidden inside many heuristics.

## Cost per unit along paths

Each arc has a capacity and a cost per unit of flow. The cost of a flow is $\sum_{(u,v)} f(u,v) \cdot \text{cost}(u,v)$. For one unit, the cheapest flow is the cheapest $s$-$t$ path by cost, found with Dijkstra or Bellman-Ford.

With arcs $s \to a$ (cost 4), $a \to t$ (cost 7), $s \to c$ (cost 6), $c \to t$ (cost 7), one unit costs $\min(4 + 7, 6 + 7) = 11$. When more units are needed than the cheapest path's bottleneck allows, the next units must use the next-cheapest routes, and costs grow.

## Successive shortest paths

The successive-shortest-path algorithm sends flow along the cheapest residual path, pushes as much as its bottleneck (or the remaining demand) allows, updates the residual graph, and repeats. Backward residual arcs carry the negated cost, the saving from cancelling a unit. A flow is optimal exactly when its residual graph has no negative-cost cycle.

Because residual costs can be negative, the shortest-path step uses Bellman-Ford, or Dijkstra with node potentials (reduced costs) that keep all residual costs non-negative. For integer capacities the number of augmentations is at most the flow value. Network simplex and cost scaling are the faster algorithms used in production solvers, but successive shortest paths is the clearest to implement and verify.

## Legalization as a transportation problem

After global placement, some bins hold more cell area than fits and others have room. Model each overfull bin as a supply equal to its excess, each bin with room as a demand equal to its free capacity, and each supply-to-demand arc with a cost equal to the move distance. Min-cost flow then moves the excess at least total displacement.

On a row of bins with 4 excess cells in bin 4, 1 in bin 5, and room for 3 in bin 1 and 4 in bin 3, only four cells fit in bin 3, so one cell must travel to bin 1. The optimum costs 8, for example four cells from bin 4 to bin 3 (cost 4) and the cell from bin 5 to bin 1 (cost 4). Greedy nearest-room choices can be worse in general, because an early cell can take room that a later cell needed more. Flow-based spreading and legalization appear in several academic and industrial placers.

## Going deeper

Ahuja, Magnanti and Orlin, "Network Flows: Theory, Algorithms, and Applications" (1993), chapters 9 and 10; Brenner, Pauli and Vygen, "Almost optimum placement legalization by minimum cost flow and dynamic programming" (Proc. ISPD, 2004).
