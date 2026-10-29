# The Hanan grid

A Steiner tree may add junction points that are not pins, and in the plane there are infinitely many candidate points. Hanan's theorem, from 1966, removes that infinity for rectilinear trees: some optimal tree uses only intersections of horizontal and vertical lines through the pins. That finite grid is the search space of exact rectilinear Steiner algorithms and the candidate set of Steiner heuristics.

## Drawing the grid

Draw a horizontal line through every pin and a vertical line through every pin. Their intersections form the **Hanan grid**. If the pins have $X$ distinct x-coordinates and $Y$ distinct y-coordinates, the grid has $X \cdot Y$ points.

Five pins with all-distinct coordinates give $5 \times 5 = 25$ points. Pins sharing a coordinate share a line: six pins with only 4 distinct x values and 5 distinct y values give $4 \times 5 = 20$ points. The grid is defined purely by coordinates, so it is built in $O(n \log n)$ by sorting and deduplicating.

## Candidate Steiner points

Every pin is a Hanan point, so the candidate Steiner points are the $XY - n$ grid points that are not pins. Hanan's theorem states that some rectilinear Steiner minimum tree uses only Hanan points as Steiner points. Combined with the fact that an optimal tree on $n$ pins needs at most $n - 2$ Steiner points, this makes the exact problem a finite search.

Twenty grid points and six pins leave 14 candidates. The search is still exponential, which is consistent with the problem being NP-hard, but for small nets it is practical, and it is how lookup-table methods such as FLUTE precompute optimal topologies for nets of up to nine pins.

## The Hanan grid as a graph

The Hanan grid is also a graph: grid points are vertices, and segments between neighbouring points are edges weighted by their length. Each of the $Y$ horizontal lines contributes $X - 1$ edges and each of the $X$ vertical lines $Y - 1$, so there are $(X - 1)Y + X(Y - 1)$ edges. With $X = 3$ and $Y = 4$, that is $8 + 9 = 17$ edges.

A rectilinear Steiner minimum tree is then a minimum Steiner tree in this graph spanning the pin vertices, and graph Steiner algorithms (dynamic programming over terminal subsets, or integer programming) solve it exactly for small $n$. Routers use the same idea with obstacles, where an escape graph built from pin and obstacle edges plays the Hanan grid's role.

## Going deeper

Hanan, "On Steiner's problem with rectilinear distance" (SIAM J. Applied Mathematics, 1966); Chu and Wong, "FLUTE: fast lookup table based rectilinear Steiner minimal tree algorithm for VLSI design" (IEEE TCAD, 2008).
