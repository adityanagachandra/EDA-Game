# QoR baselines and delta thresholds

A QoR baseline records what an engine achieved on a reference suite before a change: wirelength, worst and total negative slack, overflow, DRC count, power, runtime. A change is judged by its delta against that baseline. The judgment is statistical, because heuristic engines vary with seeds and small perturbations, and directional, because some metrics improve upward and others downward.

## A baseline with spread

A useful baseline stores both the mean $\mu$ and the standard deviation $\sigma$ of each metric across several seeds, not a single run. A new result $x$ is then judged by its z-score $z = (x - \mu)/\sigma$. Suppose routed wirelength on a reference block has $\mu = 1{,}250{,}000$ um and $\sigma = 4{,}000$ um over ten seeds. A change that gives $1{,}262{,}000$ um has $z = 3.0$: three standard deviations worse, very unlikely to be seed noise. A change that gives $1{,}254{,}000$ um has $z = 1.0$, well within the noise. Without the spread, both look like regressions of the same kind.

## Delta thresholds

Many teams state thresholds as relative deltas: a change may not worsen any metric by more than a set percentage on any design. The relative delta is $100 \times (x_{new} - x_{base}) / x_{base}$. With a threshold of 0.5% on wirelength, a block going from 820,000 to 825,000 um has a delta of $+0.61\%$ and fails; another going from 640,000 to 642,500 um has $+0.39\%$ and passes. Thresholds should exceed the seed noise, or every change fails at random, and be tighter than the effects that matter, or real regressions pass.

## Direction matters

Each metric has a direction. Wirelength, area, power, overflow, DRC count and runtime are better when smaller. WNS and TNS, as slacks, are better when larger (closer to zero or positive). A regression is a move in the bad direction beyond the threshold. A WNS change from -12 ps to -15 ps is a regression of 3 ps even though the number got smaller in magnitude terms, and relative deltas on quantities near zero, such as a WNS of -2 ps, are meaningless, so slack thresholds are stated in absolute picoseconds. A differential report should state each metric's direction and units explicitly.

## Going deeper

The OpenROAD-flow-scripts documentation on its metrics and regression checks; Jung, Kahng et al., "METRICS2.1 and Flow Tuning in the IEEE CEDA Robust Design Flow and OpenROAD" (ICCAD 2021).
