# Floorplanning by simulated annealing

Simulated annealing is the classic search engine for floorplanning: propose a move on the representation, evaluate the new floorplan's cost, accept improvements always and worsenings with a probability that falls as the temperature cools. Its quality depends on the cost function, the move set and the cooling schedule, all of which a floorplanner designer chooses and tunes.

## A cost that mixes area and wirelength

A typical floorplan cost is $C = \alpha A + \beta W$, the bounding-box area plus the total wirelength between blocks, with terms for outline violation, aspect ratio or timing as needed. The weights balance quantities with different units, so they are usually normalised by the values of a reference floorplan. Each move is evaluated by the change $\Delta C = C_{new} - C_{old}$. With $C = A + 2W$, a move from A = 526, W = 499 to A = 470, W = 435 changes the cost from 1524 to 1340, $\Delta C = -184$, a clear improvement. Wirelength between blocks is estimated by HPWL over block centres or pin positions after packing.

## The Metropolis acceptance rule

An improving move ($\Delta C \le 0$) is always accepted. A worsening move is accepted with probability $e^{-\Delta C/T}$. At T = 80 a move that raises cost by 56 is accepted with probability $e^{-0.7} = 0.4966$; at T = 10 the same move is accepted with probability $e^{-5.6} = 0.0037$. Early, at high temperature, the search wanders freely and escapes poor local minima; late, at low temperature, it behaves like greedy descent. A random number drawn uniformly in [0, 1) is compared against the probability. The rule is the same for Polish expressions, sequence pairs and B*-trees: only the move set and packing change.

## The cooling schedule

The temperature starts high enough that most uphill moves are accepted and is lowered by a schedule, most simply geometric: $T_{k+1} = r T_k$ with $r$ around 0.85 to 0.95, holding each temperature for a number of moves proportional to the problem size. The number of temperature steps from $T_0$ to $T_f$ is the smallest $k$ with $T_0 r^k < T_f$, that is $k > \ln(T_f/T_0)/\ln r$: from 2000 to below 0.5 at $r = 0.9$ takes 79 steps. Slower cooling gives better results and longer runs. Fixed-outline floorplanners such as Parquet add outline-violation penalties and moves aimed at pulling the packing into the outline. Annealing is easy to implement and to extend with new cost terms, which is why it remains common for macro-scale problems.

## Going deeper

Kirkpatrick, Gelatt and Vecchi, "Optimization by simulated annealing" (Science, 1983); Adya and Markov, "Fixed-outline floorplanning: enabling hierarchical design" (IEEE TVLSI, 2003), on the Parquet floorplanner.
