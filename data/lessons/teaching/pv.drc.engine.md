# How a DRC engine works

A DRC engine turns a rule deck and a layout into a list of violations. Its core is a small set of geometric primitives (merging, Boolean operations, sizing, edge-pair distance queries) driven by spatial data structures so that each shape is compared only with nearby shapes. Engine designers trade between edge-based and polygon-based formulations and between batch and incremental use.

## Candidate pairs and spatial filtering

Comparing every pair of $n$ shapes costs $n(n-1)/2$ tests, impossible for billions of shapes. Engines first find candidate pairs whose bounding boxes, one expanded by the rule distance $S$, overlap: both the x gap and the y gap are below $S$. Only candidates get the exact Euclidean test. For six shapes there are 15 pairs; perhaps 4 survive the box filter, and the exact test confirms some of those, since a diagonal pair with gaps 15 and 15 nm passes the filter for $S = 20$ but is 21.2 nm apart. Spatial indexes (R-trees, grids) or a plane sweep find candidates in near-linear time.

## Polygon-based checks by sizing

A polygon-based engine expresses rules as Boolean algebra on layers. A spacing check can be written as: grow every shape by $S/2$ and report where grown shapes of different polygons overlap. For shapes facing across a gap $g$ with y overlap $o_y$, the marker is $(S - g)$ by $(o_y + S)$, because the growth also extends the shapes' ends. With $S = 24$, $g = 5$ and $o_y = 14$, the marker area is $19 \cdot 38 = 722$ square nm. Square growth measures the Chebyshev metric, so engines that need Euclidean spacing round corners or switch to edge-pair tests near corners.

## Edge-based and polygon-based engines

Edge-based engines, such as Magic's, express each rule as a relation between edges: for an edge of layer A, look outward a distance $d$ for edges of a forbidden layer. They are natural for incremental checking while editing. Polygon-based engines, such as the Boolean pipelines in KLayout and commercial signoff tools, build derived layers and run measurements on them, which suits rule decks with many layer combinations. Both need merged geometry, exact integer predicates and a halo when the die is cut into tiles for parallel work. Results come back as edge pairs or small marker polygons tagged with the rule name, which downstream tools and agents consume.

## Going deeper

The KLayout DRC reference manual; Ousterhout, "Corner Stitching: A Data-Structuring Technique for VLSI Layout Tools" (IEEE TCAD, 1984).
