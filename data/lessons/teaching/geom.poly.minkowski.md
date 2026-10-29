# Minkowski sums and keep-out regions

The Minkowski sum $A \oplus B = \{a + b : a \in A, b \in B\}$ is the operation behind sizing, keep-out regions and configuration spaces. Growing a shape by $d$ is its sum with a square or disc. The set of positions where a via centre would violate spacing to an obstacle is a sum of the obstacle with the via's shape. Knowing how sums behave for rectangles and convex polygons turns these questions into arithmetic.

## Sums of rectangles

For axis-parallel rectangles the sum is again a rectangle whose ranges are the sums of the ranges:

$$[a_1, a_2] \times [b_1, b_2] \oplus [c_1, c_2] \times [d_1, d_2] = [a_1 + c_1, a_2 + c_2] \times [b_1 + d_1, b_2 + d_2].$$

Widths and heights add. Example: $[3, 9] \times [2, 6]$ plus $[-2, 1] \times [0, 5]$ is $[1, 10] \times [2, 11]$, area 81. Sizing by $d$ is the special case $B = [-d, d]^2$, which adds $2d$ to each dimension. Because the Minkowski sum distributes over union, $(A_1 \cup A_2) \oplus B = (A_1 \oplus B) \cup (A_2 \oplus B)$, a rectilinear polygon decomposed into rectangles can be grown by growing each rectangle and taking the union. This is also why sizing a union of rectangles never needs anything more than per-rectangle arithmetic followed by one union, a pattern that grow-based checks rely on throughout the rule deck.

## Convex sums and edge directions

For convex polygons $P$ and $Q$, the boundary of $P \oplus Q$ is obtained by merging their edges sorted by direction, like merging two sorted lists. Edges with the same outward normal fuse into one longer edge. So with $m$ and $n$ edges and $k$ shared normal directions, the sum has $m + n - k$ edges, and it is computed in $O(m + n)$ time. Example: a rectangle (normals 0, 90, 180, 270 degrees) plus an octagon with all eight 45-degree normals shares 4 directions, so the sum has $4 + 8 - 4 = 8$ edges. This is why growing an octagonal via by a square keeps it an octagon, and why sums of non-convex shapes, which can have quadratic complexity, are computed by decomposing into convex pieces first.

## Configuration space for placement

Ask where the centre of a via may go without violating spacing $s$ to an obstacle. Shrink the via to its centre point and expand the obstacle by the via's half-dimensions plus the spacing: the forbidden region is the obstacle summed with the via shape reflected through its centre and grown by $s$. For an obstacle $w \times h$, a via $a \times b$ and square-metric spacing $s$, the forbidden region is $(w + a + 2s) \times (h + b + 2s)$. Example: obstacle $200 \times 120$ nm, via $40 \times 40$ nm, $s = 50$ nm gives $340 \times 260$ nm. Legal positions are the complement, so detailed routers and legalizers test a single point against grown obstacles instead of testing two shapes against each other.

## Going deeper

de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 13 on robot motion planning and Minkowski sums; Lozano-Perez, "Spatial planning: a configuration space approach" (IEEE Transactions on Computers, 1983).
