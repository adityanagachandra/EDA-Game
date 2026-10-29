# Boolean operations on layout polygons

DRC rule decks, derived layers, extraction and mask preparation are written as Boolean expressions over layers: gate = poly AND diffusion, contacts NOT covered by metal, the XOR of two versions of a layer. A Boolean engine takes two sets of polygons and returns the region in both, either, one but not the other, or exactly one. It is the core primitive of a layout geometry kernel, and its correctness is non-negotiable.

## AND as intersection

A AND B is the region covered by both layers. For two rectangles it is the intersection box computed from the overlaps of the x- and y-ranges, with area $o_x o_y$ when both overlaps are positive. Example: $A = [2, 10] \times [1, 7]$ and $B = [6, 14] \times [4, 12]$ give $[6, 10] \times [4, 7]$, area 12. For general polygons the engine sweeps both layers' edges together, tracking at each point how many shapes of each layer cover it, and emits the boundary where the coverage condition changes. AND is how derived layers are formed: a transistor gate region is poly AND active, and the gate length and width checks run on that derived layer.

## OR as union

A OR B is the region covered by at least one layer. Its area follows from inclusion-exclusion,

$$|A \cup B| = |A| + |B| - |A \cap B|,$$

so for the rectangles above, $48 + 64 - 12 = 100$. Union also merges: overlapping or touching input shapes on one layer become a single polygon, which is the normal first step before width, spacing and notch checks, because a spacing rule between two drawn rectangles that overlap is meaningless once they are one piece of metal. For many shapes, inclusion-exclusion explodes combinatorially, so engines compute unions by a sweep that maintains coverage counts, not by subtracting pairwise intersections.

## NOT and XOR

A NOT B (set difference) removes B's region from A: $|A \setminus B| = |A| - |A \cap B|$. A XOR B keeps what is in exactly one layer: $|A \oplus B| = |A| + |B| - 2|A \cap B|$. With the rectangles above, NOT gives $48 - 12 = 36$ and XOR gives $48 + 64 - 24 = 88$. NOT expresses enclosure and coverage rules: vias NOT covered by metal are violations. XOR is the comparison operator: two versions of a layout, or the outputs of two engines, should XOR to nothing, and the size and shape of a nonzero XOR tells exactly where they differ. All four operations are instances of one sweep with a different coverage predicate, which is how engines implement them.

## Going deeper

Vatti, "A generic solution to polygon clipping" (Communications of the ACM, 1992); the Boost.Polygon and KLayout documentation on Boolean operations; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on physical verification.
