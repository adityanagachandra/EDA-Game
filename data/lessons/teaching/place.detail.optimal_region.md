# The optimal region of a cell

For a single cell with every other cell fixed, the positions that minimise the total HPWL of its nets form a rectangle, its optimal region, found from medians of the other pins' bounds. Detailed placers use it to decide where a cell wants to go before they search for a legal spot there, and the computation is a neat application of the median minimising a sum of distances.

## Why medians appear

Remove cell c and look at each of its nets: the remaining pins span an interval $[l_j, r_j]$ in x. With c at x, the net's x span is $(r_j - l_j) + \mathrm{dist}(x, [l_j, r_j])$, so the cell's total cost is a constant plus $\sum_j \mathrm{dist}(x, [l_j, r_j])$. This is a convex piecewise-linear function of x whose slope is the number of interval endpoints below x minus the number of nets, so it is minimised wherever exactly half of the $2k$ endpoints lie below x. Sort the $2k$ endpoints; any x between the $k$-th and the $(k+1)$-th is optimal. For two nets with others' intervals [11, 11] and [6, 46], the sorted bounds are 6, 11, 11, 46 and the optimal x is exactly 11.

## The gain from moving there

Comparing the cost at the current position with the cost in the optimal interval tells how much the cell could gain. For nets with others' intervals [9, 44] and [24, 51], the sorted bounds are 9, 24, 44, 51, so the optimal interval is [24, 44]. With the cell at x = 60, the spans are 51 and 36, total 87; anywhere in [24, 44] they are 35 and 27, total 62, a potential gain of 25. The y dimension is computed independently, so the optimal region is a rectangle, $[x_k, x_{k+1}] \times [y_k, y_{k+1}]$. The Manhattan distance from the cell to that rectangle measures how misplaced the cell is.

## Using optimal regions in detailed placement

FastDP's global swap uses the optimal region as a target: for a cell, find the optimal region, then try swapping with each cell in it or moving into a gap there, and accept the best improving move, which must account for the moved partner's own change in cost. The region ignores legality and other cells' preferences, so it is a guide, not a destination. Because every cell's region depends on the others, the computation is repeated as cells move. Cells whose region contains their current position are already locally optimal and can be skipped, which prunes most of the work in later passes. Timing-driven versions weight the intervals by net criticality, which turns medians into weighted medians.

## Going deeper

Pan, Viswanathan and Chu, "An efficient and effective detailed placement algorithm" (ICCAD 2005); Goto, "An efficient algorithm for the two-dimensional placement problem in electrical circuit layout" (IEEE Trans. Circuits and Systems, 1981).
