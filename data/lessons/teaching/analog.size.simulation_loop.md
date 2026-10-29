# Simulation-in-the-loop sizing

Simulation-based sizing calls SPICE inside the optimization loop, so every performance is evaluated with the real device models. It is accurate by construction and works for any topology, but its cost is measured in simulations. Planning that budget is the first engineering task of a simulation-in-the-loop sizer, and it shapes which optimizers are viable.

## The simulation budget

The cost of a run is iterations times candidates per iteration times corners times seconds per simulation, divided by the parallel licenses or cores. An evolutionary optimizer running 200 generations of 20 candidates at 15 corners, with 30-second simulations on 16 licenses, needs $200 \cdot 20 \cdot 15 = 60{,}000$ simulations and $60{,}000 \cdot 30/16/3600 = 31.3$ hours. Adding mismatch Monte Carlo at 100 samples per point multiplies that by 100. These numbers explain the field's design choices: sample-efficient optimizers, staged flows that optimize at the typical corner first and add corners later, and surrogate models that stand in for SPICE during most of the search.

## Gradients from a black box

SPICE returns performances, not gradients. A gradient-based sizer estimates them by finite differences: forward differences need one nominal simulation plus one per variable, $n + 1$ in total; central differences need $2n$ but are more accurate. With 20 variables at 5 corners, one forward-difference gradient costs 105 simulations. Simulators can also compute sensitivities directly by adjoint methods for some analyses, at roughly the cost of one extra solve, which is far cheaper when variables are many. The step size matters: too small and simulator tolerances dominate the difference; too large and the curvature does. Many simulation-based sizers avoid gradients altogether for this reason and use derivative-free search.

## Estimating yield in the loop

Yield-aware sizing estimates the fraction of Monte Carlo samples that meet every spec. With $n$ samples and $k$ passes, $\hat{Y} = k/n$ and its binomial standard error is $\sqrt{\hat{Y}(1-\hat{Y})/n}$. For 900 passes out of 1000, $\hat{Y} = 0.9$ and the standard error is 0.95 percentage points. Halving the error needs four times the samples, and high-yield targets (99.9% and above) need rare-event methods such as importance sampling or worst-case-distance analysis, because plain Monte Carlo sees almost no failures. A loop that compares candidates by yield must therefore account for estimation noise: two candidates whose estimates differ by less than the standard error are not distinguishable, and an optimizer that treats them as different is chasing noise.

## Going deeper

Rutenbar, Gielen and Roychowdhury, "Hierarchical modeling, optimization, and synthesis for system-level analog and RF designs" (Proceedings of the IEEE, 2007).
