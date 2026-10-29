# A domain brain model

The domain brain model is the learned component that knows the skill library: which skills exist, what they do, and above all when to pick each one and when to customize or recurse. It is trained on the record of selections and outcomes that the system accumulates. It turns each promoted skill and each run into better decisions next time, and it must do so without carrying any customer's design across the IP boundary.

## What the brain learns

The brain learns a mapping from situation to action. The situation is a summary of the design and the run: design traits (size, utilization, regularity, macro count, clock count), the current objective, and the metrics and their variance from the fast loop. The action is a skill choice, a parameter setting, a decision to customize, or a decision to recurse to an accurate engine. The training signal is the measured outcome: the QoR delta the choice produced, judged by the engines. Each record is a tuple of situation, action and outcome. Promoted skills add new actions; their validation runs add outcomes.

## Learning inside the boundary

The training data contains design-derived features, so where training happens matters. A shared base model can learn from what may leave a site: promoted skills, their manifests, their tests on vendor-held designs, and outcomes on public or synthetic designs. A per-customer model is fine-tuned on-site from that customer's own records and stays there. This layering requires a model whose weights the deployment controls, because fine-tuning and running a model inside a customer's environment is impossible with a model available only as a remote API.

## Failure modes of a learned selector

A learned selector inherits the classic problems of learning. It can **overfit**: trained mostly on one customer's designs, it picks the same datapath skill for designs that only look similar. It can suffer **distribution shift**: a new process node or design style falls outside everything it has seen, and its confidence is unearned. And it creates a **feedback loop**: it gets outcome data only for the skills it chooses, so a skill it never picks never gets the chance to prove itself. The remedies are held-out evaluation on unlike designs, an uncertainty estimate that triggers exploration or recursion, and a deliberate fraction of exploratory choices.

## Going deeper

Sutton and Barto, "Reinforcement Learning: An Introduction" (2nd ed., 2018), chapter 2 on exploration; the agent-first EDA notes in this curriculum, section on the skill flywheel.
