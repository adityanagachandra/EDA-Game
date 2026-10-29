# Present and historical congestion costs

PathFinder's two congestion costs play different roles. The present cost reacts to this iteration's sharing and grows by a schedule; the history cost accumulates every past overuse and never decreases. History is what breaks oscillations, where two nets keep trading the same contested node, and the schedules of both are the main tuning knobs of a negotiated-congestion router.

## Accumulating history

After each iteration, PathFinder updates $h_n \leftarrow h_n + h_{fac} \cdot \max(0, \mathrm{occ}_n - \mathrm{cap}_n)$. A node of capacity 2 with occupancies 2, 0, 4, 4 and 2 over five iterations, at $h_{fac} = 1.25$, is overused by 2 twice, so its history grows to $1.25 \times (2 + 2) = 5$. History never falls, so a node that was contested for many iterations stays expensive even when it is momentarily free. This makes the router learn which resources are scarce: nets gradually abandon chronically congested regions rather than flowing back into them as soon as the present cost drops. VPR's original form multiplies the base cost by a history term that starts at 1 instead of adding to it, but the idea is the same.

## Growing the present cost

The present-cost factor starts small and is multiplied by a constant after every iteration: $p_{fac}^{(i)} = p_0 \cdot m^{i-1}$. With $p_0 = 1$ and $m = 1.5$, iteration 6 uses $p_{fac} = 1.5^5 = 7.59$. A node with $b = 2$, $h = 2.25$, occupancy 2 and capacity 1 then has $p = 1 + 7.59 \times 2 = 16.19$ and costs $4.25 \times 16.19 = 68.8$. Early iterations tolerate sharing, so every net finds its preferred region; late iterations make sharing so expensive that only necessary conflicts remain. A gentle multiplier converges slowly but to better quality; an aggressive one converges fast but freezes nets into detours too early.

## History breaks ties for good

Consider two nets that both want node A, of capacity 1 and base cost 4, while one of them has a detour costing 10. With $p_{fac}$ fixed at 1, the net that sees the other on A prices A at $(4 + h)(1 + 1)$. In iteration $i$, history is $h = h_{fac}(i - 1)$; with $h_{fac} = 0.5$, A costs 8, 9, 10, then 11 in iteration 4, the first time it exceeds the detour. Present cost alone would never change in this scenario, and the two nets could alternate on A forever. History turns a symmetric standoff into a decision. The schedule of $h_{fac}$ sets how many iterations such standoffs take to resolve.

## Going deeper

McMurchie and Ebeling, FPGA (1995); Betz, Rose and Marquardt (1999), chapter 6; the VPR router documentation for pres_fac and acc_fac.
