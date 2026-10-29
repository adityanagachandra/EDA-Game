# Polygon representations in layout databases

The same region can be stored as a vertex ring, a set of oriented edges, a stack of trapezoids, a bitmap or a hierarchy of placed cells. Each representation makes some operations cheap and others expensive, and a geometry engine converts between them constantly: rings for storage and area, edges for Booleans, trapezoids for mask writing, hierarchy for repetition. Knowing what each one preserves is how a tool designer decides where an operation should run.

## Vertex rings and their orientation

The canonical storage form is a ring of integer vertices with an orientation convention: outer boundaries counter-clockwise, holes clockwise. A ring gives area and perimeter in one pass and is compact. For Manhattan polygons each step changes only x or only y, so storing alternating coordinates halves the data. Rings should be normalized: no consecutive collinear vertices, no repeated points, a fixed starting vertex. Then two tools that produce the same region produce the same ring, which makes comparison and hashing trivial. The weakness is that rings say nothing about how shapes relate to each other: overlaps between rings, or a ring that touches itself, are invisible until an algorithm looks for them.

## Edge sets, trapezoids and scanlines

Boolean engines break rings into oriented edges, each carrying a direction that says which side is inside. A sweep over all edges, counting coverage, evaluates any Boolean expression for any number of overlapping shapes, and output edges are chained back into rings at the end. A scanline (trapezoid) form cuts each region into bands between consecutive vertex y-values, each band a list of x-intervals; clipping and Booleans become interval operations, and mask writers consume trapezoids directly as exposure figures. A bitmap makes coverage and density trivial but quantizes edges to the pixel size. The trade-off is always the same: the form that makes one operation linear makes reconstruction of another expensive.

## Holes, cut lines and file formats

GDSII stores a boundary as one vertex list with no notion of holes, so a ring-shaped region is written with a cut line: a zero-width slit from the outer boundary to the hole and back (the keyhole). A reader must merge the polygon before checking it, or the slit looks like two edges at zero spacing. OASIS adds rectangle and trapezoid records, repetitions for arrays and modal values that avoid restating shared fields, with variable-length integers, so regular structures compress dramatically. Above shapes sits hierarchy: a cell stored once and placed many times, with transforms. An engine that works on hierarchy directly processes each cell once and only revisits instances where their surroundings interact.

## Going deeper

The SEMI P39 OASIS specification; the Calma GDSII Stream Format manual; the KLayout documentation on polygon and edge classes.
