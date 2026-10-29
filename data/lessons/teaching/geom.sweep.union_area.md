# Area of a union of rectangles by sweep

The area covered by overlapping rectangles, counting overlaps once, is needed for density checks, fill targets and capacitance estimates on unmerged data. Inclusion-exclusion over many shapes is hopeless, but a sweep solves it in $O(n \log n)$: the area is the integral over x of the covered length of the sweep line. The same machinery, with counts instead of flags, measures area covered at least $k$ times.

## Slabs and covered length

Sort the distinct x-coordinates of all rectangle edges. Between consecutive ones, the set of rectangles cut by the sweep line is constant, so the covered region within that slab is (slab width) times (covered length of the line). The union area is the sum over slabs:

$$A = \sum_i (x_{i+1} - x_i)\, L(x_i).$$

Example: $[0, 4] \times [0, 3]$ and $[2, 6] \times [1, 5]$. Slab $[0, 2]$ has covered length 3, slab $[2, 4]$ has the union of $[0, 3]$ and $[1, 5]$, length 5, and slab $[4, 6]$ has length 4. Total $6 + 10 + 8 = 24 = 12 + 16 - 4$, agreeing with inclusion-exclusion. This is the discrete form of integrating covered length over x. The number of slabs is at most $2n - 1$, so the work is dominated by computing each slab's covered length quickly.

## Maintaining the covered length

At each stop, the status must report $L$, the length of the union of active y-intervals, under insertions and deletions. Recomputing it from scratch costs $O(n)$ per stop and $O(n^2)$ overall. A segment tree over the compressed y-coordinates stores in each node a cover count (how many active intervals cover the node's whole range) and the covered length within the node: the full range if the count is positive, otherwise the sum of the children. Each insertion or deletion touches $O(\log n)$ nodes, and $L$ is read at the root. Example: active y-ranges $[0, 3]$, $[1, 5]$ and $[7, 9]$ merge into $[0, 5]$ and $[7, 9]$, so $L = 7$.

## Area covered at least twice

Replace the covered flag by a cover count and the same sweep measures the area covered by at least $k$ rectangles: in each slab, sum the lengths where the count is at least $k$. With $k = 2$ this finds double coverage, which matters for checks such as overlapping fill, double-patterning conflicts and duplicated geometry. Example: three rectangles pairwise overlapping in regions of area 4, 6 and 2, with a common triple region of area 1, cover $4 + 6 + 2 - 2 \cdot 1 = 10$ at least twice. Inclusion-exclusion needs every intersection; the sweep needs only counts. A segment tree supports this by storing, per node, the length covered at least once and at least twice.

## Going deeper

Bentley, "Algorithms for Klee's rectangle problems" (Carnegie Mellon University report, 1977); de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 10 on segment trees.
