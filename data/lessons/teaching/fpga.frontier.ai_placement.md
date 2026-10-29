# Learned policies inside open placers

Learned components are entering open FPGA placers: policies that choose which move to try in an annealer, models that predict routing congestion from a placement, and learned delay estimators. Used well, they speed up search inside an engine whose core still guarantees legality; used badly, they optimize a proxy and miss the real objective. The design question is where the learned piece sits and what grounds its reward.

## Learned policies behind seams

A move-selection policy decides which move type to propose next (a random swap, a move toward the median of a block's nets, a swap within a small window). The placer's core still checks legality, evaluates the exact cost change and applies the Metropolis rule, so a bad policy wastes moves but cannot corrupt the placement. This makes FPGA placement a good setting for tool-grounded reinforcement learning: the placer computes exact incremental costs in microseconds, giving a dense reward, namely the cost reduction each proposed move produced. A congestion predictor consulted by an analytical placer's cost function sits behind a similar seam. In both cases the learned piece is strategy, and legality and the real routing stay in their engines.

## Proxies and reward hacking

A reward based on bounding-box cost alone can teach a policy to make placements whose wirelength improves while routed timing gets worse, because routed delay was never in the reward. That is reward hacking on a proxy. The remedies are to include the real objective, to check samples periodically with the accurate engines, and to treat placement-stage estimates as fast proxies that need correction. Generalization is the other risk: a policy trained on small benchmarks must be shown to match or beat the default policy at equal run time on held-out designs of the target size. A policy that sometimes yields an unroutable result must pass the same routability and timing checks as any placer, with a fallback to the default when it fails.

## Judging and deploying learned placers

A claim that a learned placer beats VPR is meaningful only with the same architecture, benchmarks, routed timing and run-time budget for both. A learned move-type policy is easier to deploy than an end-to-end learned placer, because it plugs into an existing seam and inherits the engine's guarantees, so its failure modes are bounded. Replacing the annealer's acceptance test with a learned classifier, by contrast, changes the algorithm's convergence behavior and is a core change needing human review and full regression. When a policy is trained on a customer's designs on site, what returns to the shared library is the validated policy or its know-how, not the designs: the IP boundary is part of the architecture.

## Going deeper

Murray et al., "VTR 8", ACM TRETS (2020), for VPR's placement options.
