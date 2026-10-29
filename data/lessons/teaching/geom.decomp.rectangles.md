# Decomposing rectilinear polygons into rectangles

Many consumers of layout geometry want rectangles, not polygons: mask writers, some extractors, fill generators, simple DRC engines and data compressors. Partitioning a rectilinear polygon into the fewest non-overlapping rectangles is solvable exactly in polynomial time, by a beautiful reduction to bipartite matching. Covering with overlapping rectangles, by contrast, is NP-hard, which is why tools partition.

## Reflex vertices

A vertex of a rectilinear polygon is convex (interior angle 90 degrees) or reflex (270 degrees). Walking the boundary counter-clockwise, every left turn is convex and every right turn reflex, and the turns must total one full rotation, so convex minus reflex equals 4. A simple rectilinear polygon with $n$ vertices therefore has $r = (n - 4)/2$ reflex vertices. Example: an L-shape has 6 vertices and 1 reflex vertex; a plus shape has 12 vertices and 4 reflex vertices. In a polygon with holes, each hole's corners count from the region's point of view: a rectangular hole contributes 4 reflex vertices. Reflex vertices are where any rectangle partition must place a cut, since no rectangle can have a 270-degree corner.

## Chords and the r minus g plus 1 formula

Every reflex vertex needs at least one cut. A single axis-parallel cut from a reflex vertex resolves one reflex vertex and adds one rectangle, giving $r + 1$ rectangles. A chord, an axis-parallel segment through the interior joining two reflex vertices, resolves two at once. The minimum partition of a simple rectilinear polygon has

$$r - g + 1$$

rectangles, where $g$ is the largest number of chords that pairwise do not intersect, not even at endpoints. Horizontal chords never cross each other, nor do vertical ones, so the chord intersection graph is bipartite, and its maximum independent set is the number of chords minus a maximum matching (Konig's theorem). Example: a plus shape has $r = 4$ and four chords forming a cycle; $g = 2$, so 3 rectangles. A T-shape has $r = 2$, one chord, and 2 rectangles; a U-shape has no chord and needs 3.

## Holes and degenerate cases

With $h$ holes the minimum becomes $r - g - h + 1$, where $r$ counts reflex vertices on all rings. Example: a square frame (one rectangular hole) has $r = 4$ and no chords, so $4 - 0 - 1 + 1 = 4$ rectangles, the four sides. A rectangle with two holes side by side at the same height has $r = 8$ and two chords between the holes, so $8 - 2 - 2 + 1 = 5$. The algorithm cuts along the chosen chords, then resolves each remaining reflex vertex with one cut, all in polynomial time. Collinear chains of reflex vertices and degenerate holes need care, and they are where implementations break. The contrast matters for tool design: the minimum cover, with overlaps allowed, is NP-hard, so mask and fill tools partition, and sometimes accept more pieces in exchange for speed.

## Going deeper

Lipski, Lodi, Luccio, Mugnai and Pagli, "On two-dimensional data organization II" (Fundamenta Informaticae, 1979); Ohtsuki, "Minimum dissection of rectilinear regions" (IEEE ISCAS, 1982); Eppstein, "Graph-theoretic solutions to computational geometry problems" (2009 survey).
