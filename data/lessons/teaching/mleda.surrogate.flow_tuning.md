# Flow parameter tuning with multi-fidelity search

A physical-design flow has dozens of parameters: target density, timing-driven weights, routing layer limits, clock-tree options. Tuning them per design is a black-box search where each full run takes hours. Autotuners, such as the AutoTuner shipped with OpenROAD-flow-scripts, combine model-based search with early stopping of unpromising runs, which is where most of the saving comes from.

## Successive halving

Successive halving starts many configurations on a small budget and repeatedly keeps the best fraction $1/\eta$, giving each survivor $\eta$ times more budget. With 27 configurations, $\eta = 3$ and a first budget of 1 CPU-hour, 9 survive the first elimination, 3 the second and 1 the third. In a flow, budget can be stages reached (stop after placement, after clock-tree synthesis, after global route) or optimizer iterations. The method assumes that a configuration's early result ranks it roughly where its final result would. If that assumption fails, the method discards good configurations early and never sees them again.

## Counting the cost

Each rung costs the number of survivors times their budget. With 27 configurations, $\eta = 3$ and 1 hour: $27 \cdot 1 + 9 \cdot 3 + 3 \cdot 9 + 1 \cdot 27 = 108$ CPU-hours, if each rung restarts from scratch. Every rung costs the same 27 hours, so the total is the number of rungs times $n r$. Running all 27 configurations at the final budget of 27 hours would cost 729 hours, so the search is $729/108 = 6.75$ times cheaper. Warm-starting survivors from their previous checkpoint saves more, and Hyperband runs several such brackets with different starting budgets to hedge against a bad early ranking.

## When early results mislead

The saving is real only if short runs rank configurations as full runs would. In a flow, a configuration that looks good after placement can fail after routing, for example a high target density with low HPWL and later routing congestion. Low-fidelity metrics must therefore correlate with the final objective, and the correlation should be measured, for example by rank correlation between stage-level and final QoR over past runs. Multi-objective goals (timing, power and area together) need a Pareto front rather than one score, and a tuned configuration found for one design is a starting point for the next, not a guarantee.

## Going deeper

Li et al., "Hyperband: A Novel Bandit-Based Approach to Hyperparameter Optimization," JMLR (2018); the OpenROAD-flow-scripts documentation for AutoTuner.
