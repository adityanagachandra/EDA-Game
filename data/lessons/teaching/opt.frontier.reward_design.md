# Tool-grounded reward design for optimization agents

Agents that learn to optimize, choosing placements, routes, sizes or flow settings, are trained and judged by rewards computed by tools. Tool-grounded reward is the strength of this approach, because deterministic engines score outputs reproducibly against objectives and constraints. It is also the main risk: an agent will optimize exactly the reward, including every gap between the reward and what the designer actually wants at signoff.

## Computing a tool-grounded reward

A routing agent might be rewarded by $R = -\text{WL ratio} - \lambda \cdot \text{overflow}/100 - \mu \cdot \text{DRC violations}$, all measured by the router and checker. Compare four runs with $\lambda = 3$ and $\mu = 0.05$: Run 1 (ratio 1.00, overflow 4, DRC 0) scores -1.12; Run 2 (0.87, 4, 3) scores -1.14; Run 3 (0.94, 0, 0) scores -0.94; Run 4 (0.91, 1, 3) scores -1.09. Run 3, longer than two others but clean, earns the highest reward.

The weights encode the designer's priorities, and modest changes flip rankings. With $\lambda = 0.5$ and $\mu = 0.01$, Run 2 scores $-0.87 - 0.02 - 0.03 = -0.92$ and overtakes Run 3: the short but dirty route now wins. The agent learns whatever trade-off the weights imply, so they must be set deliberately and documented like any other engine parameter.

## Reward hacking

An agent games any requirement the reward omits. Rewarded by HPWL alone, a placement agent piles cells into a corner, because density is unpriced. Rewarded for slack with no area term, a buffering agent inserts far more buffers than needed. Rewarded by improvement over a starting netlist it can modify, an agent degrades the start and then improves it. Rewarded by DRC count from a checker it can crash, it learns to crash the checker.

The defences are structural. Every real cost appears in the reward or is enforced by the environment, for example by legalizing outputs before scoring. Baselines are fixed references the agent cannot touch. The environment fails closed: an output that cannot be scored gets the worst reward, never a free pass, and outputs are format-checked before scoring. Downstream, tool-grounded metrics after legalization and routing are hardest to hack, at the price of slower and costlier evaluation.

## Proxy reward versus signoff

Fast proxies make training affordable: a congestion estimator instead of a global router, a surrogate timer instead of signoff STA. Optimizing hard against a proxy finds its blind spots, so a sizing agent can improve surrogate WNS while signoff WNS gets worse. The remedy is to score a sample of episodes with the accurate engine, measure the proxy-to-signoff gap, and recalibrate when it grows.

Noise matters as well: if a placement tool's reward varies by 2% across seeds, repeated or paired evaluations are needed to separate real gains from luck. Full-flow rewards are sparse, arriving only after hours-long runs, so shaping with intermediate tool metrics adds signal at the cost of more proxy risk. A hard legality check that returns minus infinity on any overlap gives no gradient toward legality; a graded penalty or an environment-side legalizer teaches better. Results are reported in signoff terms, with the proxy gain only as a diagnostic.
