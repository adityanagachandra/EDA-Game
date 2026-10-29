# Standard-cell architecture and DTCO

A standard-cell library fixes a cell height, a set of routing tracks inside each cell, power rails at the top and bottom, and a gate grid along the width. Those choices couple process pitches to logic density. Design-technology co-optimization (DTCO) treats them as variables, and track height is the first one: how many metal tracks tall a cell is.

## Cell height as tracks times pitch

Cell height is written in tracks of the horizontal routing layer: a 6-track (6T) library on a 28 nm metal pitch is $6 \times 28 = 168$ nm tall. Those tracks hold the cell's internal wiring, the pins the router reaches, and a share of the power rails at the cell boundaries, shared with the abutting rows above and below. Libraries at older nodes used 9 to 12 tracks; FinFET nodes moved to about 7.5 and 6, and recent nodes to around 5. Fractional values such as 7.5 arise when the rails are wider than one track or when the cell grid is not an integer number of pitches.

## Why track height is a DTCO variable

Fewer tracks shrink every cell, but each track removed takes away internal wiring resources and pin access for the router. A cell must still fit its transistors: fewer tracks force fewer fins per device or a different device architecture. Rails must still carry current: narrowing them raises IR drop and electromigration risk, which pushes toward buried or backside power. The cell height reaches the router as the row pitch and as the number of pin access points, so a library that looks smaller can route worse. DTCO evaluates track height by its effect on routed block area and timing, not by cell area alone.

## Area arithmetic of a track reduction

At fixed metal pitch and fixed cell widths, cell area scales with track height. Going from 7.5T to 6T cuts cell area by $1 - 6/7.5 = 20\%$; from 6T to 5T, by 16.7%. Whether a block shrinks as much depends on achievable placement utilization: if fewer tracks raise congestion, the placer must spread cells and the block may shrink only a fraction of the cell gain. A useful habit for a tool designer is to report both numbers, the cell-level reduction and the routed block reduction, and to look for the mechanism when they differ: pin access, power rails or local congestion.

## Going deeper

Clark et al., "ASAP7: A 7-nm finFET predictive process design kit", Microelectronics Journal (2016); Liebmann et al., papers on design-technology co-optimization at advanced nodes in Proc. SPIE.
