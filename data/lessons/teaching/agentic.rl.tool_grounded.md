# Tool-grounded reinforcement learning

EDA has something most domains lack: deterministic engines that can score an answer exactly. A placer's legality checker, a router's DRC engine, a timer, an equivalence checker and a compiler all return verdicts and numbers that do not depend on anyone's opinion. Tool-grounded reinforcement learning uses those engines as the reward for training or tuning a model. It is powerful because the reward is objective, and risky because every reward is a proxy for what the designer actually wants.

## Engines as the reward function

In reinforcement learning, a policy produces actions and receives rewards, and training adjusts the policy to raise expected reward. Policy-gradient methods do this by increasing the probability of actions that earned above-average reward. In tool-grounded RL the action is a model output, such as a code edit, a placement of macros, a choice of skill or a set of flow parameters, and the reward is computed by running engines on it: did it compile, are the goldens passed, is the result legal, how much did WNS or wirelength improve against the baseline? The model never grades itself. Reinforcement learning for macro placement is a well-known example, with a reward built from proxy wirelength, density and congestion.

## What can be rewarded

Rewards come in kinds. Binary verdicts: the edit compiles, the goldens pass, the layout is DRC-clean, two netlists are equivalent. Continuous metrics: wirelength, slack, overflow, power, runtime. Composite rewards gate on verdicts and then score metrics. Verdicts are sparse, since an edit that fails to compile earns nothing and teaches little, while metrics are dense but easier to game. The strongest reward signals in this domain are those a checker can establish exactly, such as legality and equivalence, because they cannot be argued with. QoR metrics are informative but always partial.

## Where the reward is computed

Computing a reward means running engines, so cost and location matter. A rollout that needs a full place-and-route and signoff run takes hours, so training uses fast engines for most rewards and accurate engines for periodic evaluation, which creates a gap between what is optimized and what is wanted. Rewards that involve customer designs must be computed on-site, inside the sandbox and the IP boundary, so training on customer data happens where the data is. And reward engines must be deterministic and versioned, or the same action earns different rewards and training chases noise.

## Going deeper

Sutton and Barto, "Reinforcement Learning: An Introduction" (2nd ed., 2018), chapter 13 on policy gradients; Mirhoseini et al., "A graph placement methodology for fast chip design" (Nature, 2021).
