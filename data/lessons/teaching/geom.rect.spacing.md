# Spacing between rectangles

Minimum spacing is the most common design rule: two shapes on a layer must be at least $s$ apart. The rule is stated as a distance, but a checker needs a formula for the distance between two rectangles, a choice of metric, and a way to avoid square roots. This concept computes rectangle-to-rectangle distance and shows how the metric changes the answer at corners, which is where most real spacing disputes happen.

## Facing edges

For rectangles $A$ and $B$ define the horizontal gap $g_x = \max(0, B.x_1 - A.x_2, A.x_1 - B.x_2)$ and the vertical gap $g_y$ likewise. If the x-projections overlap then $g_x = 0$, the shapes face each other across a horizontal or vertical channel, and the distance is simply the other gap. Example: $A = [0, 60] \times [0, 40]$ and $B = [85, 120] \times [10, 70]$ overlap in y, so the spacing is $g_x = 85 - 60 = 25$ nm. This is the common case in routed layouts, where parallel wires on a track grid face each other along their length. It is also where the parallel run length is measured, which more advanced spacing tables use. Most spacing violations in routed metal are of this facing kind, between neighbouring tracks.

## Corner to corner

When both gaps are positive, the closest points are two corners and the Euclidean distance is $\sqrt{g_x^2 + g_y^2}$. For $g_x = 30$ and $g_y = 40$ nm that is 50 nm, although neither gap alone is that large. Different metrics give different answers here: the square (Chebyshev) metric reports $\max(g_x, g_y) = 40$ nm, and a projection metric, which only measures shapes whose projections overlap, ignores the pair entirely. Foundry spacing rules are normally Euclidean, so corner-to-corner pairs need the full formula. Many DRC engines expose the metric as an option (Euclidean, square, projection), and choosing the wrong one either misses violations or floods the report with false ones at corners.

## Checking a rule without square roots

To test a rule $d \ge s$ the checker never takes a square root: it compares $g_x^2 + g_y^2$ with $s^2$ in exact integer arithmetic. A violation is $g_x^2 + g_y^2 < s^2$ for separated shapes. Example with $s = 50$: a neighbour with gaps $(0, 45)$ violates, $(30, 30)$ violates because $1800 < 2500$, and $(40, 35)$ passes because $2825 \ge 2500$. Shapes that overlap or touch merge into one polygon and are not spacing-checked against each other; they are handled by width and notch rules instead. The checker also needs a candidate set: only shapes within distance $s$ can violate, so a spatial query with the shape grown by $s$ finds the few neighbours worth testing.

## Going deeper

The KLayout DRC reference, on the euclidian, square and projection metrics of space and width checks; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on design rules.
