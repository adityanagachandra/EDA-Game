# Growing and shrinking polygons

Sizing moves every edge of a polygon outward (grow) or inward (shrink) by a distance $d$. It is how rule decks build halos and keep-outs, how bias is applied in mask preparation, and, combined with Booleans, how width and spacing checks are implemented. Morphologically, growing is a Minkowski sum with a square (or a disc) and shrinking is the corresponding erosion, so the results have exact geometric meaning.

## Sizing one rectangle

Growing a $w \times h$ rectangle by $d$ on every side gives $(w + 2d) \times (h + 2d)$; shrinking gives $(w - 2d) \times (h - 2d)$, and the rectangle disappears when either dimension is at most $2d$. Example: a $40 \times 26$ nm rectangle grown by 5 nm becomes $50 \times 36$, area 1800; shrunk by 5 nm it becomes $30 \times 16$, area 480; shrunk by 13 nm its height reaches 0 and it vanishes. Sizing rectilinear shapes with square corners keeps them rectilinear, which is why it is the default in layout tools. Growing is a Minkowski sum with a $2d \times 2d$ square centred at the origin, and shrinking is erosion by the same square. Engines treat the collapse as a removal, never as a negative rectangle.

## Sizing merges shapes

Grow two disjoint shapes and they merge when the gap between them is at most $2d$, because each side contributes $d$. That is the basis of spacing-by-sizing: grow each shape by $s/2$, and any separated pair closer than $s$ now overlaps. After growing, the union must be computed correctly. Example: $A = [0, 20] \times [0, 15]$ and $B = [26, 40] \times [5, 30]$ grown by 4 become $[-4, 24] \times [-4, 19]$ and $[22, 44] \times [1, 34]$; their overlap is $2 \times 18 = 36$, so the merged area is $28 \cdot 23 + 22 \cdot 33 - 36 = 1334$. Note that sizing a union differs from the union of sized shapes only in the overlap bookkeeping; the region itself is the same, because Minkowski sum distributes over union.

## Closing and opening

Growing then shrinking by the same $d$ is a closing: it fills gaps and notches narrower than $2d$ and otherwise returns the original shape. Shrinking then growing is an opening: it deletes parts narrower than $2d$ and restores the rest. They are not inverses. Example: two $h = 20$ bars with a 6 nm gap, closed with $d = 4$, become one bar spanning the gap, because the grown shapes fused and the bridge survives the shrink. With $d = 2$ the gap ($6 > 4$) never closes and the shapes return unchanged. Closing is used to merge nearly abutting shapes and detect notches; opening is the classic implementation of a minimum-width check, where whatever the opening removes is a width violation.

## Going deeper

Serra, "Image Analysis and Mathematical Morphology" (1982), on dilation, erosion, opening and closing; the KLayout DRC reference on sized, and on how width and space checks relate to sizing.
