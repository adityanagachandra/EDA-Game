# Timing-driven placement cost

Timing-driven placement adds a timing term to VPR's wirelength cost. Each connection gets a criticality from a static timing analysis, its delay is estimated from the placement, and the annealer minimizes a normalized blend of timing cost and bounding-box cost. This is T-VPlace, the timing-driven placer of VPR, and it shows how a timer and an optimizer are coupled.

## Connection criticality

After a timing analysis, connection $(i, j)$ from driver $i$ to sink $j$ has criticality $\mathrm{crit}(i,j) = 1 - \mathrm{slack}(i,j)/D_{max}$, where $D_{max}$ is the critical-path delay. A connection on the critical path has slack 0 and criticality 1; one with half the critical delay as slack has criticality 0.5. The placer uses $\mathrm{crit}^{\theta}$, with a criticality exponent $\theta$ that grows during the anneal (from 1 toward 8 in VPR), so late in the anneal only the most critical connections carry weight. With $\theta = 5$, a connection of criticality 0.5 weighs $0.5^5 = 0.031$. Criticalities are recomputed periodically, not after every move, because a full timing analysis is too slow to run per swap.

## Timing cost

The timing cost of a placement is $\sum_{(i,j)} \mathrm{delay}(i,j) \cdot \mathrm{crit}(i,j)^{\theta}$, where the delay of each connection comes from a lookup table indexed by the distance between the blocks, built once by routing sample connections through an empty fabric. With $D_{max} = 1660$ ps and $\theta = 2$, connections of delay 450 ps and slack 1440 ps, 220 ps and slack 220 ps, and 450 ps and slack 620 ps contribute $450 \times 0.0176 + 220 \times 0.7525 + 450 \times 0.3925 = 350.08$ ps. Moving a block changes the delays of its connections, so the timing cost updates incrementally like the wirelength cost. Delays are estimates; the true delays exist only after routing.

## Blending timing and wirelength

VPR scores a move by $\Delta C = \lambda \frac{\Delta \mathrm{TimingCost}}{\mathrm{PrevTimingCost}} + (1 - \lambda) \frac{\Delta \mathrm{BBCost}}{\mathrm{PrevBBCost}}$, with the previous costs recorded at the start of the temperature. Normalizing each term by its total makes the two dimensionless, so $\lambda$ sets their relative weight regardless of units; VPR's default is $\lambda = 0.5$. With $\lambda = 0.7$, previous costs 3704 and 2442, and a move that adds 34 to timing cost and removes 72 from wirelength, $\Delta C = 0.7 \times 34/3704 + 0.3 \times (-72)/2442 = -0.0024$, so the move is accepted. The normalization is why the temperature schedule still works when the cost mixes picoseconds and channel spans.

## Going deeper

Marquardt, Betz and Rose, "Timing-Driven Placement for FPGAs", FPGA (2000); the VPR placer documentation for timing_tradeoff and the criticality exponent.
