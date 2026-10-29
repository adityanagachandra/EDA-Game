# Rectilinear Steiner minimum trees

The shortest way to connect a net's pins with horizontal and vertical wires is a rectilinear Steiner minimum tree, which may branch at points that are not pins. It is the ideal topology for routing, the best wirelength estimate for placement, and a classic NP-hard problem. Tool designers need its exact behaviour for small nets, which dominate real designs, and its relation to the cheaper estimates.

## Three pins and the median point

For three pins the RSMT has at most one Steiner point, at the median x-coordinate and the median y-coordinate of the pins. Its length equals the half-perimeter of the bounding box: $\max x - \min x + \max y - \min y$. Every wire from the median point to a pin is monotone, and together they cover the full width and height of the box exactly once.

Pins $(4, 6)$, $(9, 9)$ and $(1, 1)$ have medians $x = 4$, $y = 6$, which here is the first pin, so no extra point is needed, and the length is $8 + 8 = 16$. For two and three pins, half-perimeter wirelength is therefore exact, which is why HPWL works well as a placement objective for the many small nets.

## Exact search for small nets

For four or more pins the RSMT can exceed the half-perimeter, and its topology is not fixed. By Hanan's theorem it suffices to try sets of at most $n - 2$ non-pin Hanan points: for each set $S$, compute the rectilinear MST of the pins plus $S$, and take the minimum. That minimum is the RSMT length, because an optimal tree is a spanning tree on its own pins and Steiner points.

For four pins this means trying the empty set, every single candidate and every pair, at most $1 + 12 + 66$ MSTs. A four-pin net with MST 17 and half-perimeter 13 might have RSMT 15. Exact methods scale exponentially; GeoSteiner solves nets with hundreds of pins exactly, but routers use FLUTE's lookup tables for nets up to nine pins and heuristics beyond.

## Savings and hardness

The saving of the RSMT over the rectilinear MST is often small but not negligible: a four-pin net might save 2 units out of 19. On random pin sets the average saving is roughly ten percent, and the worst case is a third of the MST, reached by four pins in a plus shape.

The rectilinear Steiner tree problem is NP-hard (Garey and Johnson, 1977), so large nets are handled by heuristics such as iterated 1-Steiner, edge-substitution methods, and FLUTE's net breaking. For routing, a Steiner tree also has to respect obstacles, congestion and timing: the minimum-length tree is not always the best tree, and timing-driven routers trade length for shorter source-to-critical-sink paths.

## Going deeper

Hwang, Richards and Winter, "The Steiner Tree Problem" (1992); Garey and Johnson, "The rectilinear Steiner tree problem is NP-complete" (SIAM J. Applied Mathematics, 1977); Chu and Wong (IEEE TCAD, 2008).
