# Tool-grounded rewards for board routing

Training a routing policy with reinforcement learning needs a reward computed by tools: DRC counts, via counts, wirelength against a baseline, skew from the constraint reporter. Such rewards are cheap and exact, but they measure only what they include, and a policy will exploit any gap between the reward and the real objective. Writing the reward, computing it, and detecting reward hacking against signoff are the three skills.

## A tool-grounded reward

A reward combines tool outputs with weights. Example: $R = 100 - 20 V - 0.1 N - 2 W - 5 S$, where $V$ is the DRC violation count, $N$ the via count, $W$ the wirelength excess over a baseline router in percent, and $S$ the worst intra-pair skew in ps. A candidate with no violations, 120 vias, 1.5% excess wirelength and 2 ps skew scores $100 - 12 - 3 - 10 = 75$. Every term is computed by a deterministic tool, the DRC engine, the database and the constraint reporter, so the reward is reproducible and cannot be talked up by the policy. Normalizing wirelength against a baseline keeps the reward comparable across designs of different size.

## Weights encode priorities

The weights say what the policy should care about. Comparing two candidates exposes them: one DRC violation costs 20 points, as much as 200 vias, which encodes that legality comes first; one picosecond of skew costs as much as 50 vias. Example: candidate A with 59 vias, 2.4% excess and 1 ps skew scores 84.3, candidate B with 175 vias, 0.1% excess and 5.5 ps skew scores 54.8, so A wins by 29.5 despite more wirelength. Badly scaled weights give perverse policies: if vias are nearly free, the policy sprinkles them; if one violation is barely penalized, the policy trades violations for completion. Hard constraints are often better as gates, with the reward defined only on legal candidates, than as large penalties.

## Reward hacking against signoff

A reward can only reward what it measures. Here there is no signal-integrity term, so a policy can win by routing pairs over plane splits or through long via stubs, which costs nothing in the reward and fails the eye-mask simulation. The test is to rank candidates by the proxy reward and check the top ones against the real signoff: if several of the highest-reward candidates fail SI, the policy is exploiting the gap. Fixes are to add a fast SI estimate to the reward, gate the reward on an SI check, or periodically score policies against full signoff. The general lesson of tool-grounded RL is that the gap between a fast proxy and signoff quality is the central risk, and it must be measured, not assumed small.

## Going deeper

The agent-first EDA notes in this curriculum (tool-grounded reinforcement learning). Sutton and Barto, "Reinforcement Learning: An Introduction" (2nd ed., 2018), chapter 3, on reward design.
