# Advanced OCV derates by depth and distance

Flat OCV applies the same derate to a one-cell path and a forty-cell path. But random variation partly cancels along a path: the relative spread of a sum of $N$ independent delays shrinks like $1/\sqrt{N}$. Systematic variation, on the other hand, grows with physical distance. Advanced OCV replaces the flat factor with tables indexed by logic depth and by distance, which removes much of the pessimism on deep paths.

## Looking up the derate by depth

An AOCV table lists late derates against logic depth, for example 1.10 at depth 1, 1.08 at 2, 1.05 at 4, 1.03 at 8 and 1.02 at 16. Deeper paths get smaller derates. A data path of depth 4 with a non-derated delay of 460 ps becomes $460 \times 1.05 = 483$ ps; under a flat 1.10 derate it would be 506 ps. Separate tables exist for early and late, for cells and nets, and often per cell type or Vt class. The depth used is the number of cells on the path; graph-based analysis must use a conservative depth for each pin, while path-based analysis uses the true depth of the path it re-times.

## Interpolating between depths

Depths between table points are linearly interpolated. With 1.10 at depth 2 and 1.08 at depth 4, depth 3 lies halfway, so the derate is $1.10 + 0.5 \times (1.08 - 1.10) = 1.09$, and a 570 ps path becomes 621.3 ps. The interpolation is the same one-dimensional rule as in a Liberty table: $v = v_1 + t (v_2 - v_1)$ with $t = (d - d_1)/(d_2 - d_1)$. Derates are usually given to two or three decimals, so rounding the interpolated factor too early can shift a long path's delay by a picosecond. Beyond the last depth the table's final value is normally used.

## Depth and distance together

A two-dimensional AOCV table adds distance, typically the diagonal of the bounding box enclosing the path's cells, since cells far apart see more systematic variation. The derate falls with depth and rises with distance, and a query is interpolated bilinearly in the window around it. For depths 2 and 4, distances 500 and 1000 um, values 1.13, 1.16, 1.10 and 1.13, a path of depth 3 and distance 875 um gets 1.1375. The table is derived from statistical characterization: for each depth and distance, the derate is chosen so the derated delay covers a target number of sigma of the path delay distribution. POCV replaces this tabulated summary with sigma values per arc.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 10; the Synopsys PrimeTime documentation on AOCV table formats.
