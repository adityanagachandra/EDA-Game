# Spacing checks over many shapes

A full spacing check finds every pair of separated shapes on a layer closer than the rule, across billions of shapes, and with rules that depend on widths, run lengths and corners. It combines everything before it: rectangle distance, the choice of metric, sweeps or spatial indexes to find candidates, and table lookups for advanced rules. This is the check that dominates DRC run time and most violation reports.

## Pairwise Euclidean violations

For separated shapes $A$ and $B$ with gaps $g_x$ and $g_y$, a Euclidean spacing rule $s$ is violated when $g_x^2 + g_y^2 < s^2$. Touching or overlapping shapes are merged first and never spacing-checked against each other. The checker avoids testing all $n(n - 1)/2$ pairs: it grows each shape by $s$ (or $s/2$ each) and uses a sweep or spatial index to find only the pairs that come within the rule, then evaluates the exact test on those candidates. Example: among six shapes with rule 5, perhaps 15 pairs exist but only 4 come within 5 units in both axes, and 3 of those fail the Euclidean test. The output is a set of markers, usually the region between the offending edges, for each violation.

## Square versus Euclidean metrics

The square metric, $\max(g_x, g_y)$, is cheaper (no squaring) and never larger than the Euclidean distance, so it flags every Euclidean violation plus some corner-to-corner pairs that actually pass. A pair with gaps $(3, 3)$ under rule 4 is flagged by the square metric (3 < 4) but passes the Euclidean rule ($\sqrt{18} \approx 4.24 \ge 4$). Projection metrics only measure pairs whose edges face each other, so they miss corner pairs entirely. Engines therefore use the cheap metric as a filter and the Euclidean metric for sign-off, or expose the metric per rule, because some foundry rules are written for projection only (for example, end-of-line and parallel-run rules) while basic spacing is Euclidean.

## Width and run-length dependent spacing

At advanced nodes the required spacing depends on context. A spacing table, like the LEF SPACINGTABLE PARALLELRUNLENGTH construct, indexes the requirement by the wider wire's width and by the parallel run length, the overlap of the two wires' projections along their length. Example table: for widths above 0, 150 and 400 nm and run lengths above 0, 200 and 500 nm, requirements of 70, 70, 70 / 70, 90, 110 / 70, 120, 160 nm. Two wires of widths 180 and 90 nm running in parallel for 320 nm use the row for width above 150 and the column above 200: 90 nm. The checker must therefore compute widths (which needs the merged polygon) and run lengths for each candidate pair before it knows the threshold to compare against.

## Going deeper

The LEF/DEF Language Reference, section on SPACINGTABLE; the KLayout DRC reference on space, separation and metrics; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on physical verification.
