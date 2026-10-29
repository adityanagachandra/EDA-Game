# Learning synthesis recipes

The order of ABC commands matters, and the best script depends on the design. That makes recipe selection a sequential decision problem: an agent observes the AIG, picks the next command, and receives a reward computed by the tools. Framed carefully, it is a clean case of tool-grounded reinforcement learning; framed carelessly, it rewards gaming a proxy.

## A tool-grounded reward

The reward comes from the tools, not from a model's opinion. A simple one weighs area and depth against a baseline: $R = w \frac{A_0 - A}{A_0} + (1 - w) \frac{L_0 - L}{L_0}$, with $A$ the AND count and $L$ the level count after the recipe, and $A_0$, $L_0$ the baseline. With $A_0 = 1753$, $L_0 = 39$ and $w = 0.5$, a recipe reaching 1412 nodes and 27 levels earns $0.5 \cdot 0.195 + 0.5 \cdot 0.308 = 0.251$. The weight decides whether area or depth wins. Better rewards use mapped area and delay from the mapper and timer, which cost more to compute but are what users care about.

## State, actions and policy

The state is a summary of the current AIG: node count, levels, balance statistics, and the recent command history. The actions are ABC commands (or short macro sequences). A policy maps state to the next command, trained by policy-gradient or actor-critic methods over many designs, with the episode ending after a fixed number of steps. Closed-loop policies, which observe the AIG after each command, can react to the actual effect of each pass, while open-loop scripts commit in advance. Held-out designs test whether the policy generalizes rather than memorizing what works on the training circuits.

## Where it fits and how it fails

In the four-layer model the learned policy is strategy: it orders compiled, verified transformations and never touches the AIG core, so a bad policy can waste time but cannot break correctness, provided every result passes equivalence checking. The main failure is reward hacking: if the reward is AIG node count, the agent may learn long tails of zero-cost moves or structures that count well but map badly. The defenses are to compute the reward from mapped QoR at least on a schedule, to penalize run time, and to report mapped area, delay, run time and equivalence for every adopted recipe.

## Going deeper

Hosny, Hashemi, Shalan and Reda, "DRiLLS: deep reinforcement learning for logic synthesis", ASP-DAC (2020); Sutton and Barto, "Reinforcement Learning: An Introduction" (2nd ed., 2018), on policy gradients.
