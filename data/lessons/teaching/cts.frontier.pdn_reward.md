# Tool-grounded rewards for clock and power agents

Agents that design power grids or clock trees can be trained or steered with rewards computed by the engines themselves: IR, EM, skew, latency and power from the analysis tools. Such rewards are only as good as what they measure. A reward that sees one metric is maximised at the expense of all others, and a reward computed by a fast proxy is exploited wherever the proxy is wrong. Designing these rewards is a core skill for agent-first EDA.

## What a good reward measures

A tool-grounded reward for a PDN agent comes from deterministic engines: static and dynamic IR, EM current densities, and the routing tracks and area the grid consumes. A reward of lower worst static drop alone teaches the agent to add metal everywhere. A better form treats the IR and EM limits as constraints that must hold, and penalises resource use. For a CTS agent, skew, latency, slew and power all come from the timer and power engines; a reward of negative global skew alone teaches it to pad every branch with delay, wasting power and adding variation.

## How rewards get gamed

Reward hacking exploits blind spots. If the fast IR estimator underestimates drop near macros, an agent learns to concentrate current there. If EM is checked only after the fact, it learns grids that pass IR by pushing current through narrow paths. If a useful-skew agent is rewarded for setup improvement alone, it breaks hold, since setup and hold pull skew in opposite directions. If the reward uses one corner, trees balanced at that corner can be badly skewed at others, because buffer and wire delays scale differently. Each blind spot becomes the policy's favourite place.

## Keeping the reward honest

Audit the proxy: periodically score samples with the signoff solver and penalise or retrain where proxy and signoff disagree. Take the limits from the signoff specification and the foundry's EM rules, read by the tools, never chosen by the agent being judged. Measure progress on held-out designs with signoff QoR; a reward that rises in training while held-out results stay flat means overfitting or proxy exploitation. Where metrics trade against each other, a weighted sum picks one point on the Pareto front, so hard limits or a Pareto view often express the intent better.
