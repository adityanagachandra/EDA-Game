# VPR placement by annealing

VPR places clusters by simulated annealing: start from a random legal placement, propose random swaps, always accept improvements and accept a worsening with probability $e^{-\Delta C / T}$, while the temperature $T$ falls. Annealing is simple, handles any cost function and gives high-quality FPGA placements, which is why VPR and nextpnr's original placer both use it.

## The Metropolis acceptance rule

A move that changes the cost by $\Delta C \le 0$ is always accepted. A move with $\Delta C > 0$ is accepted with probability $e^{-\Delta C/T}$. At $T = 15.5$ a swap with $\Delta C = 1.75$ is accepted with probability $e^{-1.75/15.5} = 0.893$; at $T = 1$ the same swap would be accepted with probability 0.174. High temperature lets the search escape local minima by climbing; low temperature turns it into greedy descent. The acceptance test needs only $\Delta C$, never the absolute cost, which is why incremental cost updates make annealing fast. A uniform random number $r \in [0, 1)$ decides: accept when $r < e^{-\Delta C/T}$.

## Evaluating a swap

VPR's basic move picks a block and a random target location within the range limit; if the target holds a block of the same type, the two swap, otherwise the block moves into the empty slot. $\Delta C$ sums the change in cost of every net touching either block. Take nets of at most 3 terminals with cost $(x_{max} - x_{min} + 1) + (y_{max} - y_{min} + 1)$. A net that joins exactly A and B is unchanged by swapping them, since the same two locations stay occupied. A net joining A and fixed pins changes as A moves to B's location, and likewise for B. The swap's $\Delta C$ is the new cost of those nets minus the old, computed without touching the rest of the placement.

## Choosing a temperature

Inverting the acceptance rule gives the temperature that accepts a move of size $\Delta C$ with probability $p$: $T = -\Delta C / \ln p$. To accept a move of $\Delta C = 14$ with probability 0.3 requires $T = -14 / \ln 0.3 = 11.63$. This is how to reason about a schedule: early on, $T$ should accept most moves of a typical size, and late in the anneal only small uphill moves should pass. VPR sets its starting temperature from the spread of costs of random moves and then adapts $T$ to the acceptance rate, so that the time spent at each temperature goes where the cost is changing. Annealing's quality comes from this schedule, not from the acceptance rule alone.

## Going deeper

Betz and Rose, "VPR: A New Packing, Placement and Routing Tool for FPGA Research", FPL (1997); Kirkpatrick, Gelatt and Vecchi, "Optimization by Simulated Annealing", Science (1983).
