# Parametric OCV and LVF sigma accumulation

Parametric OCV models each arc's delay not as a number scaled by a derate but as a distribution with a mean and a sigma. The timer accumulates sigma along a path statistically and margins only the path, at a chosen number of sigma. The per-arc sigmas come from Liberty Variation Format tables, characterized per cell, per arc, per transition and per slew and load point.

## Path sigma from independent arcs

If arc delays vary independently, their variances add. A path with arc sigmas $\sigma_1, \dots, \sigma_n$ has
$$\sigma_{path} = \sqrt{\textstyle\sum_i \sigma_i^2}.$$
For five arcs with sigmas 7, 9, 12, 8 and 12 ps, the sum of squares is 482 and the path sigma is 21.95 ps, much less than the 48 ps sum of the sigmas. LVF stores these sigmas in tables such as ocv_sigma_cell_rise and ocv_sigma_cell_fall, indexed like the NLDM tables, so each arc's sigma is looked up at its own slew and load. Independence is the modeling assumption for the local random component; correlated, systematic variation is handled by corners and, in statistical timing, by shared sources.

## The n-sigma arrival

The path mean is the sum of the arc means, and the late arrival used for setup is the mean plus $n$ path sigmas: $AT = \sum_i \mu_i + n\,\sigma_{path}$. With means summing to 211 ps, a path sigma of 18.73 ps and $n = 2$, the late arrival is 248.47 ps; at 3 sigma it would be 267.2 ps. The number of sigma is a signoff parameter, often 3. LVF can also give asymmetric early and late sigmas, since delay distributions are skewed at low voltage, and recent formats add mean-shift and skewness corrections; a POCV timer then uses the early sigma for the early side and the late sigma for the late side.

## Why per-arc margins are pessimistic

Margining each arc separately at $n$ sigma, which a derate effectively does, adds $n \sum_i \sigma_i$ to the path, while POCV adds $n \sqrt{\sum_i \sigma_i^2}$. The first is always at least the second, and the gap grows with depth. For three arcs with sigmas 5, 12 and 2 ps at 3 sigma, per-arc stacking adds 57 ps and POCV adds 39.46 ps, an overestimate of 17.54 ps. For $n$ equal arcs, the ratio is $\sqrt{n}$. This is the same averaging that AOCV approximates with depth tables, but POCV applies it exactly per path. In graph-based mode the timer propagates a mean and a variance per pin and merges them with a conservative rule, so path-based re-timing again recovers some pessimism.

## Going deeper

The Liberty reference manual on LVF (ocv_sigma tables); Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 10; Sapatnekar, "Timing" (2004), chapter on statistical timing.
