# Point-in-polygon tests and fill rules

Deciding whether a point is inside a polygon underlies probing a net at a location, connectivity through vias, and the definition of what a self-overlapping or multi-ring shape even covers. The classic answer casts a ray and counts boundary crossings. For layout data the subtle part is the fill rule: even-odd and nonzero winding agree on simple polygons and disagree on overlapping ones, and the choice changes what a layer means.

## Counting ray crossings

Cast a ray from the query point, say towards $+x$, and count how many polygon edges it crosses. Each crossing switches between inside and outside, and the ray ends outside, so an odd count means the point is inside. For a rectilinear polygon only vertical edges can be crossed by a horizontal ray, and an edge at $x = e$ spanning $[y_1, y_2]$ is crossed when $e > q_x$ and $y_1 < q_y < y_2$. Example: a point at $(2.5, 3.5)$ with vertical edges to its right at $x = 4$ spanning $[0, 5]$ and at $x = 7$ spanning $[3, 8]$ crosses both, count 2, so it is outside. Degenerate cases, a ray through a vertex, are handled by a half-open rule such as counting an edge when $y_1 \le q_y < y_2$.

## Testing many points

Testing one point against an $n$-edge polygon costs $O(n)$. Testing $m$ points naively costs $O(mn)$, which is too slow for, say, probing every via against every shape. Two standard speed-ups apply. First, a bounding-box rejection test removes most points in $O(1)$. Second, sort the query points by y and sweep: the active edges at each y form an ordered list, and the crossing count for a point is its rank in that list, found by binary search. Example: with active vertical edges at $x = 1, 4, 7, 9$ on the query row, a point at $x = 5$ has two edges to its right and is outside, while a point at $x = 8$ has one and is inside. A slab or trapezoid structure gives $O(\log n)$ per query after preprocessing.

## Even-odd versus nonzero winding

When a boundary overlaps itself, or several rings overlap, parity and winding disagree. The winding number of a point counts how many times the boundary wraps around it, with counter-clockwise loops adding 1 and clockwise loops subtracting 1. Under the nonzero rule a point is inside when its winding number is not 0; under the even-odd rule, when its crossing count is odd. For a path drawn as several overlapping counter-clockwise rectangles, the winding number of a point equals the number of rectangles covering it. A point covered twice is inside under nonzero but outside under even-odd, which punches an accidental hole in the overlap. Layout semantics are normally a union, which corresponds to nonzero winding, so engines merge overlapping input before applying parity-based algorithms.

## Going deeper

O'Rourke, "Computational Geometry in C" (2nd ed., 1998), chapter 7; Haines, "Point in Polygon Strategies" in Graphics Gems IV (1994).
