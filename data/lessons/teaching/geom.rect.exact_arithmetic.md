# Integer coordinates and exact predicates

A geometry engine for layout is only as trustworthy as its arithmetic. Layout data lives on an integer grid in database units, and the engine's decisions (is this point inside, do these edges cross, are these shapes touching) are signs of predicates. If a sign is computed wrongly the algorithm's decisions contradict each other and it produces garbage or crashes. This concept is about the representation choices that keep layout geometry exact and reproducible.

## Database units and the integer grid

A layout database stores every coordinate as an integer count of database units (DBU), for example 1 DBU = 1 nm. GDSII stores coordinates as 4-byte signed integers together with the unit size; OASIS uses variable-length integers. A coordinate of 12345 at 1 nm per DBU means 12.345 um. On top of the DBU grid the foundry imposes a manufacturing grid, say 5 nm, and every vertex must be a multiple of it. Integers make geometry exact: two edges are either at the same coordinate or not, areas are integers, and the same input always gives the same output. Converting to a finer DBU multiplies coordinates exactly; converting to a coarser one requires rounding and can move edges, so tools refuse or flag it.

## Exact predicates and overflow

Algorithms branch on predicate signs, such as the orientation determinant. In floating point the sign can be wrong when the true value is near zero, and then the algorithm can conclude both that a point is left of an edge and that it is right of it, producing broken polygons. With integer coordinates the predicate can be exact, but the arithmetic must be wide enough: with 32-bit coordinates a difference needs 33 bits and the determinant can need 65 bits. The standard engineering answer is a filtered predicate: evaluate quickly in floating point with a rigorous error bound, trust the sign when the value clears the bound, and recompute exactly in the rare uncertain cases. The result is exact and nearly as fast as the naive version.

## Where non-integer geometry sneaks in

Manhattan geometry is closed under Booleans: a vertical edge at $x = a$ meets a horizontal edge at $y = b$ exactly at $(a, b)$, so no new coordinates appear. Exactness breaks when non-axis edges enter. Two 45-degree or arbitrary-angle edges intersect at rational points that are generally off the grid, and sizing a diagonal edge by $d$ moves it by $d/\sqrt{2}$ in each axis. The engine must snap such points to the grid, and snapping perturbs edges: nearly coincident edges can collapse into zero-width slivers or start to cross. Transforms are another source: mirroring or rotating a cell about a half-unit coordinate maps integers to half-integers. Robust engines snap with care, then re-run the topology on the snapped result.

## Going deeper

Shewchuk, "Adaptive Precision Floating-Point Arithmetic and Fast Robust Geometric Predicates" (1997); Hobby, "Practical segment intersection with finite precision output" (Computational Geometry, 1999), on snap rounding; the Boost.Polygon documentation on its integer coordinate model.
