# Search-and-repair DRC loops

No detailed router produces a clean layout in one pass. After initial routing, a search-and-repair loop finds the remaining violations, rips up the nets around them, raises the cost of the offending locations and reroutes, over and over, until the violation count reaches zero or stops falling. The schedule of this loop decides runtime and final quality.

## Marker penalties

After a pass, the design-rule checker leaves markers at every violation. The router adds a penalty to grid cells covered by markers, so that rerouted nets avoid them. To make the loop converge, the penalty grows across iterations: in iteration $i$ a marked cell might cost $b + m_0 g^{i-1}$ for base cost $b$, initial penalty $m_0$ and growth factor $g$.

Example: $b = 3$, $m_0 = 8$ and $g = 2$ give a marked cell a cost of $3 + 128 = 131$ in iteration 5. Early iterations tolerate violations so nets can find their way; later ones make violations so expensive that any clean detour wins. This mirrors negotiated congestion, applied to rule violations instead of capacity.

## When the detour wins

A net crossing $n_m$ marked cells keeps its short route until the marker penalty outweighs the extra length of a clean detour $d$. With a penalty $m_0 2^{i-1}$, the detour first wins in the smallest iteration $i$ with $n_m m_0 2^{i-1} > d$. Example: one marked cell, $m_0 = 1$, a detour 35 cells longer: the penalty first exceeds 35 at 64, in iteration 7.

This calculation explains loop behaviour. A violation with a cheap detour clears in a couple of iterations; one whose only fix is a long detour or the cooperation of several nets lingers. A geometric schedule bounds the wait logarithmically in the detour length, which is why repair loops typically run a few tens of iterations rather than hundreds.

## Clip-based repair

Repair works locally. The router cuts the region into clips of a few gcells, and in each iteration reroutes only the clips that contain markers, keeping wires crossing the clip boundary fixed so the clip behaves like a switchbox. Clip boundaries shift between iterations, so a violation that sat on a boundary lands inside a clip next time.

Example: with 5 by 5 clips offset by one gcell, a marker at gcell (6, 34) falls in clip (1, 7). Mapping every marker and counting distinct clips gives the work for the iteration. Local repair keeps each pass cheap and parallel, since clips can be routed independently, and it limits collateral damage to nets that were already clean. TritonRoute and OpenROAD's drt run search-and-repair in this clip-based style.

## Going deeper

The TritonRoute paper (IEEE TCAD, 2021); the OpenROAD drt documentation on detailed-routing iterations.
