# Bounding boxes and rectangle intersection

The axis-parallel rectangle is the workhorse of layout geometry. Every shape has a bounding box, every spatial index stores boxes, and the first test in nearly every check is whether two boxes intersect. The operations are trivially cheap, a few comparisons, but they are executed billions of times in DRC, extraction and placement, and getting their edge cases right (touching versus overlapping, empty results) decides whether higher-level algorithms are correct.

## Boxes from points

The bounding box of a point set is $[\min x, \max x] \times [\min y, \max y]$. For a net's pins it gives the half-perimeter wirelength, $\text{HPWL} = (\max x - \min x) + (\max y - \min y)$, the standard placement estimate of routed length. It is exact for two- and three-pin nets (a rectilinear Steiner tree of three points has exactly that length) and a lower bound for larger nets. Example: pins $(3, 8)$, $(15, 2)$, $(9, 20)$ give a box from $x = 3$ to 15 and $y = 2$ to 20: HPWL $= 12 + 18 = 30$, area $12 \times 18 = 216$. The box can be updated incrementally when a pin moves, as long as the tool also tracks how many pins sit on each extreme, which is how placers keep HPWL current during millions of moves.

## Intersecting two rectangles

Two rectangles $A = [a_1, a_2] \times [b_1, b_2]$ and $B = [c_1, c_2] \times [d_1, d_2]$ intersect in the box whose x-range is $[\max(a_1, c_1), \min(a_2, c_2)]$ and likewise in y. The overlap lengths are

$$o_x = \min(a_2, c_2) - \max(a_1, c_1), \quad o_y = \min(b_2, d_2) - \max(b_1, d_1),$$

and the intersection area is $o_x o_y$ when both are positive, else 0. Example: $[2, 9] \times [1, 6]$ and $[5, 12] \times [4, 10]$ give $o_x = 9 - 5 = 4$, $o_y = 6 - 4 = 2$, area 8. The separation into independent x and y tests is the key structural fact: rectilinear intersection decomposes into one-dimensional interval problems, which is why sweeps and interval structures work so well here. The same two formulas, applied to grown boxes, answer distance questions, and applied to node bounding boxes, they drive every R-tree and bin-grid query in the engine.

## Edge cases and many rectangles

When $o_x = 0$ the rectangles touch along an edge; when both are 0 they touch at a corner. Whether touching counts as intersecting is a predicate choice: connectivity extraction must count shared edges as connected, while an overlap check for positive area must not. Each engine fixes the rule once and uses it everywhere. Intersection is associative, so the common region of several boxes is $[\max$ of left edges, $\min$ of right edges$]$ by the same in y. For $[0, 10] \times [0, 8]$, $[4, 14] \times [2, 12]$ and $[6, 9] \times [-3, 5]$, the common box is $[6, 9] \times [2, 5]$, area 9. The empty case appears as a negative extent, which code must test explicitly rather than multiplying two negatives into a positive area.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4 on wirelength estimation; the KLayout documentation on the Box class and DRC functions.
