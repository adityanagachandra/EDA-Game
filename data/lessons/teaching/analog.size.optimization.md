# Analog sizing as optimization

Sizing chooses device dimensions, bias currents and passive values for a fixed topology so that every specification holds. Written as mathematics it is a constrained, usually multi-objective, nonconvex optimization problem whose functions are evaluated by simulation. Every automated sizer, from geometric programming to Bayesian optimization to an agent, is a choice of how to pose and solve this problem, so the formulation comes first.

## Variables, objective and constraints

The decision variables are widths, lengths (or $g_m/I_D$ and $L$ per device), bias currents, and capacitor and resistor values: perhaps 10 to 50 continuous numbers for an op-amp. Specifications become constraints on simulated performances: gain at least 70 dB, GBW at least 50 MHz, phase margin at least 60 degrees, noise, offset, swing. One or two quantities that should be as small as possible, typically power and area, become the objective. Operating-point constraints keep every device in its intended region with margin, so small-signal specs remain meaningful. Writing a spec as an objective ("maximize GBW") is a common mistake: it spends power beyond what the system needs. A spec with a threshold is a constraint.

## Choosing among candidates

Given simulated candidates, the formulation decides which one wins. First discard every candidate that violates any constraint; then rank the rest by the objective. A candidate with the lowest power overall is worthless if it misses GBW by 1 MHz. When constraints are hard to satisfy, optimizers often work with a penalty: minimize the objective plus a weighted sum of normalized violations $\max(0, (\mathrm{spec} - \mathrm{measured})/\mathrm{spec})$. Normalization makes specs in different units comparable. A penalty method can return a slightly infeasible point with a better objective, so a careful flow re-checks feasibility exactly and treats any remaining violation as failure.

## Robustness, trade-offs and solution methods

A design that meets specs at the typical corner and fails at the slow corner is not a solution, so robust formulations constrain the worst case over corners or a yield target under Monte Carlo. With two objectives, such as power and GBW, the answer is a Pareto front of non-dominated designs, which lets the system designer choose the trade-off knowingly. Solution methods split by how they evaluate performance. Equation-based methods use analytical or fitted models, sometimes posynomials for geometric programming, which is convex and solves to global optimality in seconds, but only for the model. Simulation-based methods treat SPICE as a black box and search with evolutionary algorithms, simulated annealing or Bayesian optimization, which are accurate but limited by the simulation budget.

## Going deeper

Gielen and Rutenbar, "Computer-aided design of analog and mixed-signal integrated circuits" (Proceedings of the IEEE, 2000); del Mar Hershenson, Boyd and Lee, "Optimal design of a CMOS op-amp via geometric programming" (IEEE TCAD, 2001).
