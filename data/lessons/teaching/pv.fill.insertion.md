# Fill insertion to a density target

Fill tools add electrically inactive metal shapes to windows below the minimum density. The amount needed is simple arithmetic; where the fill may go is a geometric problem, because fill must keep spacing to existing metal and stay out of keep-out regions. A fill engine is a constrained placement of many small squares, followed by a density re-check.

## Fill needed for a target

If a window of area $A_w$ holds metal area $A_m$ and the minimum density is $D$, fill must add at least $D A_w - A_m$. With fill squares of side $f$ placed only in empty space, the count is $\lceil (D A_w - A_m) / f^2 \rceil$. A $130 \times 130$ um window with 3339 square um of metal and a 20 percent minimum needs $3380 - 3339 = 41$ square um, so two $5 \times 5$ squares. Real tools overshoot slightly to keep margin against the maximum density and against windows that overlap this one.

## Legal fill sites

Fill candidates sit on a grid, and a candidate is legal only if it keeps the fill spacing $s$ to all metal: the square grown by $s$ must not overlap any metal shape. Equivalently the free region is the window NOT (metal SIZED BY $s$), and a candidate is legal if it lies inside it. On a $28 \times 28$ um window with 2 um squares on a 4 um grid and 1 um spacing, 40 of 49 candidates may survive three metal shapes. Engines test candidates with a spatial index over the sized metal, or rasterize the free region.

## The best reachable density

Filling every legal site gives the largest density fill can reach: the existing metal plus the legal squares, over the window area. With 293 square um of metal and 24 legal squares of 9 square um in a 1764 square um window, the best is $509 / 1764 = 28.9$ percent. If that is still below the minimum, no fill solution exists, and the design must change: open a blockage, move a macro, or obtain a waiver. Checking this bound first prevents a fill engine from searching for a solution that cannot exist.
