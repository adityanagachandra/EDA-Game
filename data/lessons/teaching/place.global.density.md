# Bin density and overflow

A global placer controls spreading through density: the cell area in each bin of a grid laid over the placement region, compared with a target. Overflow, the excess area summed over bins, is the number that drives the penalty weight, decides when to stop and predicts how hard legalization will be. Computing it exactly and fast is a core placer kernel.

## Bin density

Divide the region into an $m \times m$ grid of bins. The density of bin $b$ is the area of movable cells overlapping it, divided by its area (or its free area, after subtracting fixed macros and blockages). A cell that straddles bins contributes to each in proportion to its overlap, computed by clipping its rectangle to the bin: a 6 by 2 cell at (4, 11) in a grid of 8 by 8 bins puts 4 by 2 = 8 um$^2$ in the upper-left bin and 2 by 2 = 4 um$^2$ in the upper-right one. During global placement cells may overlap one another, so a bin's density can exceed 1. The bin size is a trade-off: fine bins resolve local crowding, coarse bins are cheaper and smoother.

## Overflow of a bin and of the design

With target density $\rho_t$, the overflow of bin $b$ is $\max(0, A_b^{cells} - \rho_t A_b)$, the area in excess of what the bin may hold. A 12 by 12 bin with 130 um$^2$ of cells at target 0.8 has $130 - 115.2 = 14.8$ um$^2$ of overflow; a bin with 100 um$^2$ has none. The design's total overflow is usually normalised:
$$\tau = \frac{\sum_b \max\left(0, A_b^{cells} - \rho_t A_b\right)}{\sum_i A_i},$$
the fraction of cell area sitting in excess. Electrostatic placers such as RePlAce and OpenROAD's gpl stop global placement when $\tau$ reaches about 0.1 by default.

## Using overflow well

Overflow is an average, so it can hide trouble: a total of 0.08 can come from a few badly over-full bins near a macro, which legalization will resolve only with long moves. Placers therefore also watch the maximum bin density, and routability-driven flows use local targets. Overflow is not differentiable (the max and the clipping create kinks), so nonlinear placers penalise a smoothed density, such as bell-shaped cell potentials (NTUplace) or electrostatic energy (ePlace), and report true overflow only as a stopping metric. Implementation details matter for speed: bin accumulation is a scatter operation over cells, and GPU placers use atomic adds or sorting to accumulate it in parallel.

## Going deeper

Cheng, Kahng, Kang and Wang, "RePlAce: advancing solution quality and routability validation in global placement" (IEEE TCAD, 2019); the OpenROAD documentation for gpl.
