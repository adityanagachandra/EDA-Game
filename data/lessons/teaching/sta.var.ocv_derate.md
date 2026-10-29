# On-chip variation derates

Corners model chip-to-chip variation: a whole die is slow or fast. Within one die, delays still vary from cell to cell, because of local process variation, voltage drop and temperature gradients. On-chip variation analysis models this by derating: the paths that hurt a check are made slower and the paths that help it faster. Flat OCV applies one late and one early factor everywhere, the simplest and most pessimistic scheme.

## Derating the launch side

For a setup check, the arrival should be as late as possible, so the launch clock path and the data path are scaled by a late derate, a factor $1 + x$ slightly above 1. With a launch clock path of 230 ps, a data path of 630 ps including clock-to-Q, and a late derate of 1.05, the derated arrival is $(230 + 630) \times 1.05 = 903$ ps. In SDC-based flows the derates are set with set_timing_derate -late and -early, separately for cell and net delays and for clock and data paths if desired. Setup and hold times themselves are library values and are not derated.

## Derated setup slack

The required time should be as early as possible, so the capture clock path is scaled by an early derate $1 - y$ below 1: $S = T + C (1 - y) - t_{su} - (L + D)(1 + x)$. With $T = 850$, $C = 200$, $y = 0.03$, $t_{su} = 36$, $L + D = 880$ and $x = 0.03$ ps, the required time is 1008 ps, the arrival 906.4 ps and the slack 101.6 ps. Without derating the slack would have been 134 ps: OCV cost 32.4 ps. For hold the roles reverse: launch early, capture late. Flat derates scale with path delay, so long clock trees and long paths pay the most margin.

## Removing pessimism on the common clock path

The launch and capture clock paths usually share their first buffers. The flat derate scaled that shared segment late on one side and early on the other, though it is one physical path. Clock reconvergence pessimism removal adds back the difference. If the common segment has nominal delay $P$, the credit is $P(1 + x) - P(1 - y) = P(x + y)$. With $P = 160$ ps, $x = 0.08$ and $y = 0.09$, the credit is 27.2 ps, raising a raw slack of 30.8 ps to 58 ps. Flat OCV with CPPR is still pessimistic for long paths, because random variation partly averages out along many cells; that observation leads to AOCV and POCV.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 10 on on-chip variation; the Synopsys Design Constraints format specification, set_timing_derate.
