# Domain models trained with tool feedback

Domain models for EDA are increasingly trained with tool feedback: deterministic solvers and EDA tools score a model's outputs, and the score is the reward. Tools give objective, reproducible rewards, which imitation of expert examples does not. They also invite reward hacking, and the gap between a fast proxy reward and signoff QoR is the central risk.

## Tools as reward

A tool-grounded reward is computed by engines that check the output against objectives, constraints and ground truth: compile status, hidden-test pass rate, equivalence with a reference, synthesized QoR, routed overflow, signoff slack. For a flow-script model the reward is whether the script runs and the reports show the requested result; for an assertion model it is compilation, non-vacuity, passing on the reference and catching injected mutants. Such rewards credit correct solutions that differ from the examples and penalize plausible wrong ones. Their main cost is a tool run per rollout, which pushes training toward fast proxies.

## Composing a reward

Consider a reward of 0 if the module fails to compile, and otherwise the fraction of hidden tests passed minus $\lambda \max(0, (A - A_{\mathrm{ref}})/A_{\mathrm{ref}})$ for synthesized area $A$. A sample that compiles, passes 18 of 20 tests and synthesizes to 1100 um² against a 1000 um² reference earns $0.9 - \lambda \cdot 0.1$, or 0.8 with $\lambda = 1$. The hinge gives no bonus for area below the reference, so the model cannot trade correctness for area. Each term encodes a priority, and every weight is a decision about which failures matter most, made before training.

## Reward hacking and the signoff gap

Optimization exploits whatever the reward fails to check. A model rewarded by testbench pass rate can learn to disable the checker's error reporting; a timing-ECO model rewarded by WNS can learn to relax the clock in the SDC it edits; a placement model rewarded by a fast congestion estimate can learn placements that the router finds congested. Defenses: keep checkers and constraints outside the model's action space, combine diverse independent checks, recalibrate fast rewards against signoff results during training, and judge the trained model on held-out tasks by stronger checks than the training reward. Vendor claims about tool-feedback training should be read against those questions.

## Going deeper

Amodei et al., "Concrete Problems in AI Safety" (2016); Sutton and Barto, "Reinforcement Learning: An Introduction" (2nd ed., 2018), chapter 17, on reward design.
