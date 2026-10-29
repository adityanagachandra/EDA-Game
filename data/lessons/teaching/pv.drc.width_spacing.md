# Width and spacing checks

Width and spacing are the bulk of any rule deck and the first checks an engine must get exactly right. Both are distance computations on rectilinear shapes, but the details matter: which metric, whether corners count, and how the required value changes with the shapes involved. A signoff checker must report every violation and nothing else, so these definitions are exact, not approximate.

## Width violations

For a rectangle, width is the shorter side; for a general polygon it is the smallest distance between opposite edges facing inward. With a 50 nm width rule, a wire $[92, 129] \times [0, 575]$ (37 nm wide) fails and $[67, 445] \times [400, 464]$ (64 nm) passes, however long they are. Width must be measured on the merged layer: two overlapping 30 nm rectangles can form one legal 60 nm wire, while a single legal shape can have a narrow neck. Polygon engines find width violations as pairs of inward-facing edges closer than the rule.

## Euclidean spacing

Spacing between two disjoint rectangles is computed from the gaps of their projections. Let $g_x$ and $g_y$ be the x and y gaps, zero where the projections overlap. The Euclidean distance is $\sqrt{g_x^2 + g_y^2}$: if the projections overlap on one axis, the closest points face each other and the distance is the other gap; otherwise two corners are closest. With a 23 nm rule, shapes with gaps 10 and 15 nm are $\sqrt{325} \approx 18.0$ nm apart and violate, while gaps 20 and 15 give 25 nm and pass. Comparing squared distances with $S^2$ keeps the test in integers.

## Width-dependent spacing

Wide wires need more spacing, for lithography and for reliability. A typical rule is a table: 40 nm normally, 80 nm when either wire is at least 200 nm wide. Each pair's required spacing is looked up from the widths before the distance comparison. For parallel vertical wires, checking neighbours suffices when every wire is at least as wide as the largest rule value: anything farther away is separated by a whole wire. In the general case an engine expands each wide wire by the larger value and queries neighbours, so the table changes the search radius as well as the threshold.
