# Design-space exploration agents

A flow has dozens of parameters (utilization, density targets, timing-driven weights, synthesis recipes, clock-tree settings), and their best values depend on the design. A design-space exploration agent runs experiments in parallel over those parameters and uses a model to choose what to try next. Its value depends on arithmetic about wall time and budget, and its output is usually a set of trade-offs, not a single winner.

## Wall time of a parallel sweep

With $N$ experiments of duration $T$ each and $M$ machines, a sweep runs in $\lceil N / M \rceil$ waves and takes $\lceil N/M \rceil T$ of wall time. Fifty configurations of 3 hours each on 16 machines take 4 waves, 12 hours; on 25 machines, 2 waves, 6 hours; on 49 machines, still 2 waves. Exhaustive grids grow exponentially in the number of parameters (5 values for each of 6 parameters is 15,625 runs), so exploration agents choose experiments adaptively, with methods such as Bayesian optimization that fit a surrogate of QoR against parameters and pick the next point by an acquisition function.

## Successive halving

Many configurations can be judged early: a flow that is far behind after global placement rarely wins after routing. Successive halving exploits this. Start $n$ configurations with a small budget $r$ each, such as through placement only; keep the best $1/\eta$ of them; give the survivors $\eta$ times the budget; repeat until one remains. With $n = 27$, $\eta = 3$ and $r = 1$ unit, the rounds cost $27 \times 1$, $9 \times 3$, $3 \times 9$ and $1 \times 27$, a total of 108 units, against $27 \times 27 = 729$ units to run every configuration fully. Early stopping relies on early and final rankings agreeing, which is the same proxy-versus-signoff question as before.

## Pareto filtering of results

Objectives conflict: lower power often costs frequency, smaller area often costs routability. The useful output is the Pareto front, the configurations not dominated by any other. A point dominates another if it is at least as good on every objective and strictly better on one. For results of (area, delay) of (10, 9), (11, 7), (12, 7), (13, 5) and (14, 6), the point (12, 7) is dominated by (11, 7) and (14, 6) by (13, 5), so three points are on the front. The agent reports the front, and the designer, or a planner with stated priorities, chooses among it.

## Going deeper

Li et al., "Hyperband: A Novel Bandit-Based Approach to Hyperparameter Optimization" (JMLR, 2018), on successive halving; Shahriari et al., "Taking the Human Out of the Loop: A Review of Bayesian Optimization" (Proc. IEEE, 2016).
