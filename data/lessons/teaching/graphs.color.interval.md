# Interval graphs and track assignment

When conflicts come from overlapping intervals, such as nets spanning column ranges in a channel or values live over time ranges in a schedule, the conflict graph is an interval graph. Interval graphs are perfect: their chromatic number equals their largest clique, and a simple sweep colors them optimally. The left-edge algorithm of channel routing and the register count of a schedule both come from this fact.

## Density equals the chromatic number

Represent each net as a closed interval $[l, r]$ of columns. Two nets conflict, and need different tracks, when their intervals share a column. The **density** of the channel is the largest number of intervals covering a single column. All intervals covering a column pairwise overlap, so they form a clique and need distinct tracks: density is a lower bound on the track count.

For interval graphs the bound is tight, so the minimum number of tracks equals the density. With nets $[13, 15]$, $[8, 10]$, $[15, 18]$, $[15, 22]$, $[9, 10]$, $[10, 15]$, $[10, 15]$, column 15 is covered by five intervals, so five tracks are needed and suffice.

## The left-edge algorithm

Sort the intervals by left end. Take each in turn and place it on the lowest-numbered track whose last interval ends strictly before the new interval's left end; if no track fits, open a new one. This is Hashimoto and Stevens' left-edge algorithm.

With $[1, 4]$, $[2, 3]$, $[8, 9]$, $[8, 14]$, $[12, 14]$, $[14, 17]$ in left-end order: $[1,4]$ opens track 1, $[2,3]$ opens track 2, $[8,9]$ fits track 1, $[8,14]$ fits track 2, $[12,14]$ fits track 1 (track 1 last ended at 9), and $[14,17]$ cannot share a column 14 with either, so it opens track 3. Sorting dominates the cost, $O(n \log n)$. It is optimal for the interval-only problem: it never opens a track unless the current column's density forces it.

## Where the model breaks

Real channels add **vertical constraints**: at a column where one net has a pin on the top and another on the bottom, the top net's track must lie above the bottom net's. These constraints form a directed graph, and with them the problem is NP-hard; cycles in the vertical constraint graph force doglegs, splitting a net across tracks. Constrained left-edge variants process tracks top-down and only place nets whose vertical predecessors are already placed.

The interval-graph view also counts conflicts: the number of edges is the number of overlapping pairs, computed by a sweep in $O(n \log n + E)$. In high-level synthesis, the same structure gives the minimum number of registers for a schedule: values are intervals over control steps, and the left-edge sweep assigns registers optimally when lifetimes do not wrap around a loop.

## Going deeper

Hashimoto and Stevens, "Wire routing by optimizing channel assignment within large apertures" (Proc. Design Automation Workshop, 1971); Golumbic, "Algorithmic Graph Theory and Perfect Graphs" (2nd ed., 2004).
