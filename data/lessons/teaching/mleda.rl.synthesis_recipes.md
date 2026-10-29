# Reinforcement learning for synthesis recipes

Logic synthesis in ABC runs a recipe: a sequence of transforms such as balance, rewrite, refactor and resub on an and-inverter graph. Which sequence works best depends on the circuit. RL and bandit methods search recipe space using the tool's own measurements as reward, and the choice of reward decides whether the found recipe helps after technology mapping.

## The recipe space

With $K$ transforms and recipes of length $L$, repeats allowed, there are $K^L$ recipes. Six transforms and length 10 give $6^{10} = 60{,}466{,}176$ recipes, far too many to try when each takes seconds to minutes on a large circuit. Transforms do not commute, and their effect depends on the current graph, so the problem is sequential: the state is the current AIG (or features of it), the action is the next transform, and the episode ends after $L$ steps. A learned policy generalizes across circuits only if its state features capture what makes a transform useful.

## A Q-learning update

Tabular or approximate Q-learning keeps $Q(s, a)$, the expected return of applying transform $a$ in state $s$. After applying rewrite with $Q(s, \text{rewrite}) = 0.5$ and receiving reward $r = 2$, with the best next-state value 1, learning rate 0.1 and $\gamma = 0.9$, the target is $2 + 0.9 \cdot 1 = 2.9$ and the update gives $0.5 + 0.1(2.9 - 0.5) = 0.74$. Real states are AIGs, so $Q$ is a neural network over graph features, and each reward needs a tool call, which keeps the number of training episodes small.

## Shaping the reward

A typical reward combines normalized improvements in AIG node count and logic depth, $r = w(A_{\mathrm{prev}} - A)/A_0 + (1 - w)(D_{\mathrm{prev}} - D)/D_0$. Both are proxies. AIG node count and depth are cheap to measure but do not determine mapped area and delay, which depend on the cell library and the mapper. A recipe that minimizes nodes can map worse than one that keeps useful structure. Sound evaluations report mapped area and delay, or post-synthesis timing, compare with standard recipes such as the ones ABC ships, and count the tool calls spent in search, since a tuned recipe found by brute force on one circuit is not a learned result.

## Going deeper

Hosny et al., "DRiLLS: Deep Reinforcement Learning for Logic Synthesis" (ASP-DAC 2020); Mishchenko, Chatterjee and Brayton, "DAG-Aware AIG Rewriting" (DAC 2006), for the transforms being sequenced.
