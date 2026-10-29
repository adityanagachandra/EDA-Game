# Sliding density windows and gradients

Density rules apply to windows that slide across the die with a fixed step, not to a fixed tiling. Some rules also bound how fast density may change between neighbouring windows. Efficient engines precompute per-tile areas once and assemble windows from tiles, which turns a quadratic amount of clipping into prefix-sum lookups.

## Counting window positions

A $w \times w$ window stepping by $t$ over a $D_x \times D_y$ die takes $(D_x - w)/t + 1$ positions in x and $(D_y - w)/t + 1$ in y. A 20 um window with a 20 um step over a 520 by 740 um die gives $26 \cdot 37 = 962$ positions; halving the step roughly quadruples the count. Each position needs a density, so computing each from scratch would clip the same shapes many times. The step is a rule parameter: a smaller step catches more local extremes.

## The minimum window by tile sums

Cut the die into tiles of the step size and compute each tile's metal area once with a union sweep. A window covering $k \times k$ tiles has area equal to the sum of its tiles, so its density is the mean of their densities. A two-dimensional prefix sum $P$ gives any window's sum as $P(i_2, j_2) - P(i_1, j_2) - P(i_2, j_1) + P(i_1, j_1)$ in constant time. On a $5 \times 5$ tile grid with $2 \times 2$ windows there are 16 positions; the minimum mean, say 20.8 percent, is the window a fill engine must fix first.

## Density gradients

Polishing also suffers when density changes sharply between nearby regions, so some decks bound the difference between windows one step apart. Neighbouring windows share most of their tiles, so their difference depends only on the tiles that enter and leave. With $2 \times 2$ windows over a $5 \times 5$ grid, there are 24 horizontal and vertical neighbour pairs; the largest absolute difference, for example 14.5 percentage points, is compared with the limit. A gradient violation is fixed by filling the sparse side or slotting the dense side, not by meeting the minimum alone.
