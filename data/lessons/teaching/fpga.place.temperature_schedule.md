# The adaptive annealing schedule

VPR's annealing schedule adapts to the problem instead of using a fixed cooling rate. It sets the initial temperature from the cost spread, decreases $T$ by a factor chosen from the acceptance rate, makes a fixed number of moves per temperature scaled to the design size, and stops when $T$ is small compared with the cost per net. These rules are what made VPR's annealer both robust and efficient.

## Initial temperature

VPR makes $N_{blocks}$ random moves from a random placement, accepts them all, records the costs, and sets $T_0 = 20\sigma$, where $\sigma$ is the standard deviation of those costs. With recorded costs 326, 260, 228, 386, 354, 263 and 301, the mean is 302.57 and the population standard deviation 52.23, so $T_0 = 1044.6$. A temperature many times the typical cost change makes nearly every early move acceptable, so the anneal starts from a genuinely random state. Scaling by $\sigma$ rather than using a fixed number makes the schedule independent of the cost function's units, which matters once timing terms are mixed in.

## Adaptive cooling

After each temperature, VPR updates $T_{new} = \alpha T_{old}$ with $\alpha$ chosen from the fraction $R$ of moves accepted: 0.5 if $R > 0.96$, 0.9 if $0.8 < R \le 0.96$, 0.95 if $0.15 < R \le 0.8$, and 0.8 if $R \le 0.15$. Starting at $T = 63.25$, a step with $R = 0.3$ and then one with $R = 0.74$ both use $\alpha = 0.95$, giving $63.25 \times 0.95^2 = 57.083$. The schedule races through the hot, almost random phase and the cold, frozen phase, and lingers at intermediate acceptance rates, where most of the cost reduction happens. Combined with range limiting, which keeps $R$ near 0.44, the schedule spends most moves in the productive regime.

## Effort and exit

VPR makes $\mathrm{InnerNum} \cdot N_{blocks}^{4/3}$ moves per temperature, so effort grows slightly faster than the design. With $N = 1000$ and $\mathrm{InnerNum} = 10$, that is $10 \times 10{,}000 = 100{,}000$ moves per temperature; lowering InnerNum trades quality for run time almost linearly. The anneal exits when $T < 0.005 \cdot \mathrm{Cost} / N_{nets}$, that is, when the temperature is small compared with the average cost per net, so no meaningful uphill move would be accepted. With cost 177 over 270 nets the threshold is 0.00328, and from $T = 0.08$ with $\alpha = 0.8$ it takes 15 more updates to drop below it. A final greedy pass at $T = 0$ then removes the last easy improvements.

## Going deeper

Betz and Rose, FPL (1997); Betz, Rose and Marquardt (1999), chapter 5; Lam and Delosme, "Performance of a New Annealing Schedule", DAC (1988).
