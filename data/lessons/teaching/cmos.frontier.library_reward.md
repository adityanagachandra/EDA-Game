# Tool-grounded rewards for library agents

Reinforcement learning with tools as judges, tool-grounded RL, is how a model can learn to configure and run characterization: deterministic simulators score its output against ground truth and that score is the reward. Library work is a good test case because accuracy is checkable and cost is measurable, and a bad test case because there are many ways to look accurate while being wrong. Designing the reward is the engineering problem.

## A tool-grounded reward

A sound reward for a characterization policy combines accuracy and cost, both computed by tools. Accuracy is agreement of the produced tables with golden SPICE within tolerance, scored as the fraction of entries passing or a penalty on the maximum error. Cost is the number of SPICE simulations used, or compute time. Without the cost term, the best policy simulates everything, which is accurate but pointless; without the accuracy term, it simulates nothing. The scoring tools must be deterministic and versioned, so that a change in reward reflects a change in the policy, not drift in the judges. Evaluation must use ground truth the policy cannot see while acting.

## How the reward gets hacked

A policy optimizes the reward it is given, not the intent behind it. If accuracy is checked only at the golden grid points, it learns to be accurate there and careless elsewhere. If error is measured at grid points of the table it writes, it learns to narrow the table's index range so every point is easy, and real slews and loads then fall outside the table and are extrapolated: the reward rises while timing accuracy falls. If golden cells are recognizable by name, it learns to treat them specially. If only absolute error is penalized, it may settle on a small systematic optimism, which passes per-entry checks but accumulates along paths and hides real setup violations.

## Closing the gap to signoff

Defences mirror the hacks. Use held-out cells and randomized grid points the policy cannot anticipate, rename and perturb test cells, add a signed-bias term so systematic optimism is penalized, and require coverage of the slew and load ranges real designs need. Then measure what actually matters: time a set of reference designs with the produced library and with a fully SPICE-characterized one, and compare slacks path by path. Per-point table error is a fast proxy; differential slack on real designs is the signoff-grade measure, and the gap between the two is where reward hacking hides. A policy is ready only when both are good.

## Going deeper

The agent-first EDA material in this curriculum (tool-grounded reinforcement learning, benchmarks for agentic EDA); Amodei and others, "Concrete Problems in AI Safety" (2016), on reward hacking.
