# Differential testing of geometry engines

When an agent or an engineer changes a geometry engine, the proof that nothing broke is differential: run the old and new versions on reference layouts and compare the outputs geometrically. For polygon outputs the comparison is an XOR, and the expected result is usually empty. The size and shape of a nonzero XOR say exactly what changed, which makes it the acceptance test for every tier of change.

## XOR of two outputs

For a golden output $G$ and a new output $N$ on the same layer, the differential region is $G \oplus N$, with area $|G| + |N| - 2|G \cap N|$. Example: a golden shape $[2, 12] \times [3, 9]$ and a new shape $[2, 13] \times [3, 9]$ have areas 60 and 66 and overlap 60, so the XOR area is 6, a one-unit strip along the right edge. A performance change, a refactoring or a new back end must give an XOR of exactly zero on every reference layout. A rule change must give an XOR confined to where the rule was meant to act. Because the comparison is geometric, it ignores irrelevant differences such as vertex order or the starting point of a ring, which a text diff of output files would flag.

## Comparing whole layers

Real outputs are sets of shapes, and the comparison must be between unions, not shape lists. Two engines can produce different decompositions of the same region, for example one rectangle against two abutting ones, and still be equivalent; the union XOR is zero. Conversely, a change to one shape can be partly hidden by another shape that overlaps it, so the XOR can be smaller than a per-shape diff suggests. Example: golden shapes whose union covers 180 units and a new set in which one rectangle grew by 2 units along a 7-unit edge, half of that growth landing inside another shape, give an XOR of 7, not 14. The engine computes the XOR with the same sweep it uses for any Boolean, and reports markers for the differing regions.

## Reading the signature of a change

The shape of a nonzero XOR identifies the bug. A uniform shift of every shape by $\delta$ in x produces, for each rectangle of height $h$ and width greater than $\delta$, two slivers of $\delta \times h$, an XOR of $2\delta h$ per shape: a structured, layout-wide pattern that points at a coordinate transform or rounding change. Example: shifting shapes of heights 4, 7 and 10 by 2 units gives $2 \cdot 2 \cdot 21 = 84$. A sizing bug shows as thin rings around every shape, a Boolean bug as isolated slivers at special vertices. A good harness reports XOR area per layer and per region, the number of markers, and thresholds that block promotion, so an agent's change is accepted only when the differential matches its stated intent.

## Going deeper

McKeeman, "Differential testing for software" (Digital Technical Journal, 1998); the KLayout documentation on the XOR tool for layout comparison.
