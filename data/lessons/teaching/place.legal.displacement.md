# Legalization and displacement

Legalization takes the overlapping, off-grid output of global placement and produces a legal placement: every cell in a row, on a site, without overlap. It should change the global solution as little as possible, so its quality is measured by displacement, how far cells moved. Total, maximum and squared displacement each capture a different risk to wirelength and timing.

## What legalization must achieve

A legal placement puts every standard cell with its bottom edge on a row, its left edge on a site boundary, in the row's orientation, with no overlap with other cells, macros or blockages, and inside any region constraint. Global placement guarantees none of this. Legalization finds the nearest legal arrangement: it assigns each cell a row and an x position. Because the global placer worked hard to minimise wirelength, the legalizer's objective is to preserve its solution, which it measures as displacement between the global position $(x_i, y_i)$ and the legal position $(x_i', y_i')$. Cells keep their relative order where possible, since reordering lengthens wires unpredictably.

## Measuring displacement

The displacement of a cell is usually the Manhattan distance $|x_i - x_i'| + |y_i - y_i'|$. Total displacement sums it over cells; maximum displacement reports the worst single move; average displacement is the total over the cell count. For moves of 5, 6 and 6 um, the total is 17 um and the maximum 6 um. Some legalizers minimise the sum of squared displacements, $\sum_i (x_i - x_i')^2$ (Abacus does, within a row): squaring penalises a single long move more than several short ones, which tends to keep the maximum small. For moves of 6, 1 and 2 um the squared sum is $36 + 1 + 4 = 41$.

## Why the metric matters downstream

Total displacement correlates with the wirelength lost in legalization, but timing is sensitive to individual moves: one critical cell moved 80 rows can break a path even when the total is small. Large maximum displacement often signals that global placement left a local overflow, for example near a macro edge, that legalization had to resolve by pushing cells far away. Good flows therefore report both metrics, bound the maximum, weight timing-critical cells more heavily, and feed hot spots back to global placement rather than accept long moves. Legalization is also run repeatedly during optimisation, after buffering and sizing, so its speed and incremental behaviour matter as much as its quality.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 4.4 on legalization; Spindler, Schlichtmann and Johannes, "Abacus: fast legalization of standard cell circuits with minimal movement" (ISPD 2008).
