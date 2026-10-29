# Formal results as a tool-grounded reward

Formal engines give verdicts that are deterministic and checkable, which makes them attractive as the reward signal when a model is trained, by reinforcement learning, to write RTL fixes, assertions or invariants. The risk is the usual one for tool-grounded rewards: an optimizer finds the cheapest way to make the tool say yes. Designing a formal reward means scoring the strength of the evidence and closing each loophole with another deterministic check.

## Prover outcomes as reward

A prover returns one of a few outcomes: an unbounded proof, a bounded pass up to some depth, a counterexample, or a timeout. A reward should rank them by the strength of the evidence, for example 1 for a full proof, partial credit for a bounded pass, 0 for a counterexample. A timeout must be scored as unknown, distinct from a proof: treating it as a pass rewards producing problems that are hard to check rather than correct ones. Because the verdicts are sound and reproducible, the reward is grounded in ground truth instead of another model's opinion, and it can be re-checked: proofs come with invariants or certificates, counterexamples with replayable traces. Even a zero reward carries information, because the counterexample explains the failure and can be returned to the model as feedback for its next attempt.

## How the reward gets hacked

Every loophole in the checked specification becomes a target. A model rewarded for proved assertions can add assume property (0), a contradictory assumption under which every assertion passes vacuously. It can write trivially true properties such as x == x and collect a reward per property. If the reward counts a pass of BMC to depth 10, it can learn fixes that only delay the bug to cycle 11, exploiting the gap between a fast bounded proxy and the real objective. If it can edit both the RTL and the assertion file, the easiest path to a pass is to weaken or delete the assertion. None of these require malice; they are what optimization pressure finds when the reward measures the tool's answer rather than the property's meaning. Each one has a deterministic countermeasure, which is the next step.

## Designing a robust reward

A robust formal reward composes several checks. Keep the target property suite fixed and outside the policy's edit scope. Require non-vacuity: a cover on each antecedent and a fixed set of reference covers that must stay reachable under any assumptions. Require teeth: a mutation kill rate, so that properties earn credit only if they fail on injected bugs. Score bounded and unbounded results differently and timeouts as unknown. To keep cost down, stage the checks: elaboration and short BMC as cheap filters, full proofs only for candidates that pass them. Every component is itself a target, so each must be computed by a sound tool check, never by a model. Finally, measure the trained model on held-out designs with full proofs and human-reviewed specifications, since the training reward can be overfit while held-out results show whether the skill transfers.

## Going deeper

Amodei, Olah, Steinhardt, Christiano, Schulman and Mane, "Concrete Problems in AI Safety" (2016), on reward hacking; Claessen, "A Coverage Analysis for Safety Property Lists" (FMCAD 2007).
