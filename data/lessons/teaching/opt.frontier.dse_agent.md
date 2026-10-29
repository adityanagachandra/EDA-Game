# Agent-driven design-space exploration

A design-space exploration agent runs many flow experiments in parallel, chooses the next ones from the results so far, and recommends configurations. It combines everything in this domain: Bayesian optimization or successive halving to spend the budget, fast and accurate engines in a loop, Pareto fronts for PPA trade-offs, and the reproducibility and IP-boundary discipline that make an automatically found configuration trustworthy.

## Budgeting with successive halving

Successive halving starts many configurations on a small budget, keeps the better half, doubles each survivor's budget, and repeats until one remains. Start 16 configurations at 30 minutes: round costs are $16 \cdot 30$, $8 \cdot 60$, $4 \cdot 120$, $2 \cdot 240$ and $1 \cdot 480$, each 480 machine-minutes, total 2400. Running all 16 at the final 480-minute budget would cost 7680.

The method assumes that early QoR predicts final QoR, for example that a placement's overflow at 30 minutes ranks configurations the way final timing would. Where that fails, halving discards the eventual winner, so the agent should validate the ranking correlation on past runs. Hyperband runs several halving brackets with different starting budgets to hedge against that risk. Every round's survivors and their scores should be logged, since they are the evidence for the final recommendation.

## Choosing the next experiment

With a budget of 30 full runs over 6 knobs, a grid is hopeless (5 levels per knob is $5^6 = 15625$ runs), so the agent fits a surrogate to completed runs and picks the next point by an acquisition function. With 8 parallel runs per round, the batch mixes high-predicted-value points with high-uncertainty ones.

Fast and accurate engines work in a loop: screen many configurations with a 2-minute fast placer, and recurse to the 3-hour full flow on the promising few. The agent should stop trusting the fast ranking when candidates differ by less than the fast engine's known error, or when the fast metrics' variance is high. Failed runs, such as routing DRC failures, are information: recorded as infeasible, they train a feasibility model that keeps the agent out of that region. For PPA the output is the Pareto set, not a single winner.

## Trustworthy recommendations

Two configurations 0.5% apart in WNS, with seed noise of about 1%, are not ranked until multiple seeds of both are compared as distributions. Every run records its exact configuration, tool versions and seed, so any result can be reproduced and differential comparisons mean something. Runs touch the customer's design database, so they execute on-site, inside the IP boundary, in sandboxes with resource limits; the agent orchestrates rather than exporting the design.

A configuration that wins on design A and transfers poorly to design B teaches something worth keeping: which design characteristics predicted success. That conditional know-how, not a single universal setting, is what flows into the shared domain model. Finally, the recommendation must be reviewable by a human: a summary of runs tried, the front found, and the evidence for the recommended point, so approval rests on data rather than on trust in the agent.
