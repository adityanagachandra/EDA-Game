# Reporting overlapping rectangle pairs

Many checks begin by finding every pair of shapes that overlap, or that come closer than a rule distance: shorts between nets, spacing candidates, enclosure pairs. Reporting all $k$ such pairs among $n$ rectangles takes $O(n \log n + k)$ with a sweep, and the same algorithm, run on grown shapes, finds near pairs. The predicate details (overlap, touch, closer than $s$) decide what the report means.

## Overlap means both projections overlap

Two rectangles overlap with positive area exactly when their x-ranges overlap with positive length and their y-ranges do too. A sweep in x keeps active rectangles in an interval structure on y; when a rectangle is inserted, a query for intervals overlapping its y-range reports every active rectangle that overlaps it. Each pair is reported once, when the later-starting rectangle enters. Example: $[0, 6] \times [0, 4]$, $[4, 9] \times [2, 7]$ and $[7, 12] \times [0, 3]$: the first two overlap ($2 \times 2$), the last two overlap ($2 \times 1$), and the first and third have disjoint x-ranges, so 2 pairs. With an interval tree as the status, the cost is $O(n \log n + k)$.

## Near pairs by growing

To find separated pairs closer than $s$ in the square metric, where distance is $\max(g_x, g_y)$, grow every rectangle by $s/2$ and report overlaps: two shapes closer than $s$ in both axes have grown versions that overlap. Pairs that already touched or overlapped before growing are excluded, since they merge. Example with $s = 4$: shapes with gaps $(3, 0)$ and $(2, 3)$ are flagged, a shape with gaps $(5, 1)$ is not, because its x-gap alone reaches the rule. For a Euclidean rule the grown-overlap report is a superset, and each candidate is then confirmed with $g_x^2 + g_y^2 < s^2$. This filter-then-refine pattern is the general structure of distance queries.

## Touching is a separate predicate

Two closed rectangles touch when their overlaps in x and y are both at least 0 but not both positive: they share an edge segment or a corner. Touching pairs have zero overlap area, so an area-overlap report misses them, yet they matter: shapes that share an edge are electrically connected after merging, and a corner-only touch is a classic ambiguous case that design rules often forbid. Example: $[0, 3] \times [0, 3]$ and $[3, 5] \times [1, 4]$ share the edge from $(3, 1)$ to $(3, 3)$; $[0, 3] \times [0, 3]$ and $[3, 6] \times [3, 6]$ touch at one corner. An engine must define, once, whether its overlap predicate uses closed or open intervals, and apply that choice consistently in sweeps, indexes and Booleans.

## Going deeper

Edelsbrunner, "A new approach to rectangle intersections" (International Journal of Computer Mathematics, 1983); de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 10.
