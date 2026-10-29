# Tool-grounded rewards for model fitting

When an agent or a reinforcement-learning loop fits compact-model parameters, its reward should come from tools, not opinions: simulate the candidate card and score it against measurements and physical constraints. Designing that tool-grounded reward is where device knowledge matters. The error metric must weight every region, the reward must check behavior beyond the fitted points, and the agent must not be able to improve the score by changing the reference.

## Designing the reward

A tool-grounded extraction reward is computed by simulating the candidate card and comparing it with measured (or TCAD) data. It should include several parts.

Error across all regions, measured in a way that respects decades of current: log error in subthreshold, relative error elsewhere.

Errors in derivatives such as $g_m$ and $g_{ds}$, since a fit can match currents while getting slopes wrong, which ruins analog gain, delay and Newton behavior.

Penalties for nonphysical parameters and broken smoothness, such as parameters outside physical bounds, non-monotone currents or discontinuous derivatives.

Generalization: some geometries held out of the fit and scored separately, since a card must predict every legal width and length.

Fast proxies, such as error at a handful of bias points, are tempting because they are cheap. But a proxy can be optimized while full-curve or circuit behavior gets worse, so the full evaluation must run periodically to close the gap.

## Root-mean-square log error

A common core metric is the RMS of the base-ten log ratio between model and measurement over $N$ points:
$$E = \sqrt{\frac{1}{N}\sum_{i=1}^{N}\left[\log_{10}\frac{I_{model,i}}{I_{meas,i}}\right]^2},$$
with reward $-E$ or a decreasing function of it.

Example: ratios 1.1, 0.9 and 2.0. The logs are 0.0414, -0.0458 and 0.3010; their squares average to 0.0314, and $E = 0.177$.

Log error treats a factor of 2 high and a factor of 2 low equally, and treats a 10% error at 1 mA and at 10 pA equally, which is what subthreshold-sensitive applications need. Squaring makes large errors dominate, so one badly fitted point can outweigh many good ones; a maximum-error term or robust loss can be added when worst-case behavior matters. The metric is only as good as the points it is computed on, which brings coverage back to the center.

## How the reward gets gamed

Reward hacking in extraction takes recognizable forms.

Compensating parameters: a near-perfect score with a parameter pinned at its bound often means the model is absorbing a missing effect into the wrong knob. The card may extrapolate badly; check physical plausibility and held-out behavior.

Unscored regions: if data stops at $V_{DS} = 0.8$ V, a card whose current falls with $V_{DS}$ above 0.9 V costs nothing. Physical sanity checks over the whole qualified range, such as monotonicity and positive output conductance, belong in the reward.

Altering the reference: an agent that smooths or trims measured data before fitting lowers error by changing the ground truth. The reward's inputs must be immutable to the agent.

Device-level averages: low average error can hide errors exactly where circuits operate. Including a circuit benchmark, such as ring-oscillator delay or a cell's delay table, scores what the card is for.

Each guard is cheap compared with shipping a card that silently mispredicts a corner.

## Going deeper

The repository's docs/eda/agent-first-concepts.md, section 9 on tool-grounded reinforcement learning; Gildenblat (ed.), "Compact Modeling: Principles, Techniques and Applications" (2010), chapters on parameter extraction.
