# Density windows and fill

Chemical-mechanical polishing removes material unevenly unless every region of every layer has similar metal density. Foundries therefore require density between a minimum and a maximum in every window of a given size, and fill tools insert dummy metal to meet the minimum. The computation is a union area clipped to windows, repeated over a sliding grid of windows, which makes efficient area accounting the core of a density engine.

## Density of one window

The density of a window is the area of metal inside it divided by the window area. Shapes must be clipped to the window, and overlapping shapes must be merged, or overlaps are counted twice. Example: a $50 \times 50$ um window containing $[0, 20] \times [0, 10]$, $[10, 30] \times [5, 25]$ and $[40, 60] \times [40, 55]$: clip the third to $[40, 50] \times [40, 50]$ (area 100); the first two have areas 200 and 400 with overlap $10 \times 5 = 50$, so the union is 550. Total 650 out of 2500, density 26.0 percent. Clipping is an AND with the window and merging is an OR, so density is a Boolean expression plus an area, and a union-area sweep computes it.

## Fill to reach a minimum

If a window's density is below the minimum $D$, fill must add at least $D \cdot A_w - A_m$ of area, where $A_w$ is the window area and $A_m$ the existing metal. With fill squares of side $f$ that never overlap existing metal, the count is $\lceil (D A_w - A_m)/f^2 \rceil$. Example: a $100 \times 100$ um window with 1500 square um of metal and a 25 percent minimum needs $2500 - 1500 = 1000$ square um, so 63 squares of $4 \times 4$ um. Real fill must also respect spacing to signal wires, the maximum density, and coupling capacitance on critical nets, so the fill engine's job is a constrained placement problem whose feasibility depends on the free space, which is itself a Boolean NOT of the existing layer grown by the fill spacing.

## Sliding windows

Density rules apply to every position of a window sliding over the die with a given step, not to a fixed tiling. A $w \times w$ window stepping by $t$ over a $W \times W$ die takes $(W - w)/t + 1$ positions per axis, so $((W - w)/t + 1)^2$ windows. Example: a 100 um window stepping 50 um over a 2000 um die gives 39 per axis, 1521 windows. Computing each window's area from scratch would repeat work, so engines precompute the metal area per step-sized tile once (a union-area sweep per tile) and sum tiles for each window, with a two-dimensional prefix sum making every window an $O(1)$ lookup. A smaller step catches more local density extremes but multiplies the windows.

## Going deeper

Kahng, Robins, Singh and Zelikovsky, "Filling algorithms and analyses for layout density control" (IEEE TCAD, 1999); the OpenROAD documentation for its density fill tool.
