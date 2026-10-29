# Voronoi diagrams

The Voronoi diagram of a set of sites divides the plane into cells, each holding the points closer to one site than to any other. It answers nearest-site questions in one structure and appears in layout as the basis of critical-area analysis for yield (how far a defect must reach to short two wires), of proximity and clearance computation, and of meshing. Its size is linear in the number of sites, which follows from Euler's formula.

## Linear size from Euler's formula

The Voronoi diagram of $n$ sites, not all collinear, is a planar subdivision: vertices are points equidistant from three or more sites, edges are pieces of perpendicular bisectors. Every vertex has degree at least 3. Adding a vertex at infinity to close the unbounded edges and applying Euler's formula $V - E + F = 2$ with $F = n$ cells gives at most $2n - 5$ vertices and at most $3n - 6$ edges. Example: 100 sites produce at most 195 Voronoi vertices and 294 edges. So the diagram has linear size, and Fortune's sweep-line algorithm builds it in $O(n \log n)$, optimal because sorting reduces to it. A diagram of quadratic size would be useless for a billion-shape layout; a linear one is practical.

## Unbounded cells and the hull

A site's cell is unbounded exactly when the site lies on the boundary of the convex hull of all sites. For a hull site there is a direction in which it stays the nearest site forever; an interior site is surrounded on all sides. So the number of unbounded cells equals the number of sites on the hull, $h$, which for points in general position is the number of hull vertices. Example: sites $(0, 0)$, $(10, 0)$, $(5, 9)$, $(4, 3)$, $(6, 2)$ have hull vertices $(0, 0)$, $(10, 0)$, $(5, 9)$, so three cells are unbounded and two are bounded. The unbounded edges, $h$ of them, are the bisectors of consecutive hull vertices, which is how a hull computation seeds a Voronoi construction.

## Exact counts and metrics

In general position (no four sites on a circle), every Voronoi vertex has degree exactly 3 and the bounds become exact counts in terms of $n$ and the number $h$ of hull sites: $2n - 2 - h$ vertices and $3n - 3 - h$ edges, matching the triangles and edges of the dual Delaunay triangulation. Example: 50 sites with 8 on the hull give 90 vertices and 139 edges. Layout often needs other metrics. Under the $L_\infty$ or $L_1$ metric, bisectors are polygonal lines with horizontal, vertical and 45-degree pieces, which suits rectilinear geometry. Voronoi diagrams of segments, not points, give the medial axis of shapes and clearance between wires, the basis of critical-area extraction for shorts and opens.

## Going deeper

Aurenhammer, "Voronoi diagrams: a survey of a fundamental geometric data structure" (ACM Computing Surveys, 1991); Fortune, "A sweepline algorithm for Voronoi diagrams" (Algorithmica, 1987); Papadopoulou and Lee, "Critical area computation via Voronoi diagrams" (IEEE TCAD, 1999).
