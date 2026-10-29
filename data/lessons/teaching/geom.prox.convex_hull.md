# Convex hulls by monotone chain

The convex hull of a point set is the smallest convex polygon containing it. Layout tools use hulls for blockage and keep-out outlines, for quick containment and separation tests, as a building block of Voronoi and Delaunay computations, and as a lower bound in packing and floorplanning. Andrew's monotone chain computes it in $O(n \log n)$ with nothing but sorting and the orientation predicate, and it is the version to implement.

## Hull vertices

A hull vertex is an extreme point: it cannot be written as a convex combination of the others. Points strictly inside, and points lying on a hull edge between two corners, are not vertices. The hull of $n$ points has $h \le n$ vertices, and for points on a small lattice many points are collinear, so a careful implementation must decide whether to keep collinear boundary points; layout outlines normally drop them. Example: the points $(0, 0)$, $(4, 0)$, $(4, 4)$, $(0, 4)$, $(2, 2)$, $(2, 0)$ have a hull of four vertices: $(2, 2)$ is inside and $(2, 0)$ lies on the bottom edge. The hull of rectangles is the hull of their corners, since a convex set containing all corners contains every rectangle.

## Andrew's monotone chain

Sort the points by x, breaking ties by y. Build the lower hull left to right: append each point, but first pop the last point while the last two kept points and the new one fail to make a strict left turn, that is, while their orientation is at most 0. Build the upper hull the same way from right to left. Concatenate the two chains, dropping the duplicated endpoints. Each point is pushed once and popped at most once, so after the $O(n \log n)$ sort the scan is linear. Example: sorted points $(0, 0)$, $(1, 2)$, $(2, 1)$, $(3, 0)$: after $(2, 1)$ arrives, $(1, 2)$ is popped (a right turn), and after $(3, 0)$ arrives, $(2, 1)$ is popped, leaving the lower hull $(0, 0)$, $(3, 0)$ with two points.

## Area and the lower bound

Once the hull is a counter-clockwise ring, its area comes from the shoelace formula; for integer points it is a multiple of $\frac{1}{2}$. Example: the hull $(0, 0)$, $(6, 0)$, $(7, 3)$, $(2, 5)$ has shoelace sum $0 + 18 + 29 + 0 = 47$, area 23.5. The $O(n \log n)$ cost is optimal for comparison-based algorithms: sorting can be reduced to hull construction by lifting numbers $x_i$ onto the parabola $(x_i, x_i^2)$, whose hull lists them in order. Output-sensitive algorithms (Chan's, $O(n \log h)$) win when the hull is small. For rectilinear layout outlines, the orthogonal hull, the smallest region convex along every horizontal and vertical line, is often more useful than the convex hull and is computed by a similar sweep.

## Going deeper

Andrew, "Another efficient algorithm for convex hulls in two dimensions" (Information Processing Letters, 1979); de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 1; Chan, "Optimal output-sensitive convex hull algorithms in two and three dimensions" (Discrete and Computational Geometry, 1996).
