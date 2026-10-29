# Reinforcement learning for placement

The best-known RL application in EDA treats macro placement as a game: an agent places macros one at a time on a grid, a conventional method places the standard cells, and a proxy cost scores the result. Understanding exactly what the agent decides, what the reward measures and what it leaves out is the basis for judging any claim about the method.

## Macro placement as a sequence

The state is the netlist plus the partial placement; the action is the grid cell where the next macro goes, typically ordered by size; the episode ends when all macros are placed. Illegal actions, positions that overlap a placed macro or leave the canvas, are masked so the policy can only choose legal cells. Standard cells, which number in the millions, are not placed by the agent; a separate conventional placer, such as a force-directed method, places them for evaluating the reward. A graph neural network encodes the netlist and the placement for the policy and value networks, which lets one model be pretrained across blocks and fine-tuned on a new one.

## The proxy cost

Routing and signoff timing are too slow to run for every episode, so the reward is the negative of a proxy cost, $\mathrm{WL} + \lambda_c C + \lambda_d D$, combining normalized wirelength, congestion and density. A placement with $\mathrm{WL} = 0.8$, $C = 1.2$ and $D = 0.9$ and weights $\lambda_c = \lambda_d = 0.5$ costs $0.8 + 0.6 + 0.45 = 1.85$. The weights set which failure the agent avoids most. No term is routed timing or power, so a lower proxy cost is a hypothesis about a better chip, to be checked after routing.

## Judging an RL placement result

Three questions decide whether an RL placement result is meaningful. First, does a strong non-learned optimizer, such as simulated annealing on the same proxy and grid with matched compute, reach the same cost? If so, the gain comes from the formulation, not from learning. Second, does the proxy track post-route metrics? A lower proxy with worse routed timing shows the agent exploited the proxy. Third, what did the agent receive: pretraining data, initial placement information from another tool, compute? In the four-layer engine model the policy sits in the strategy layer, proposing macro positions that the core legalizes and the downstream flow evaluates.

## Going deeper

Mirhoseini et al., "A graph placement methodology for fast chip design," Nature (2021); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4, on macro placement.
