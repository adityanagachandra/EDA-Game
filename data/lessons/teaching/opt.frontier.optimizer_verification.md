# Golden and differential tests for optimizers

When agents edit and tune optimization engines, every change needs proof that it does what it claims and nothing else. Optimizers offer unusually strong evidence: legality invariants, optimality certificates such as duality gaps and KKT residuals, and QoR baselines on golden designs. Knowing which check proves what, and how to aggregate QoR across designs, is what lets a harness accept an agent's change automatically or route it to a human.

## Aggregating QoR across golden designs

A new optimizer skill is compared with the baseline design by design as ratios, new over baseline, and the ratios are aggregated by their geometric mean. With wirelength ratios 0.95, 1.08 and 1.01, the geometric mean is $(0.95 \cdot 1.08 \cdot 1.01)^{1/3} = 1.036^{1/3} = 1.012$: a 1.2% loss overall, although one design improved by 5%.

Ratios normalize away design size, so a large design cannot dominate the summary as it would with absolute wirelength. The geometric mean treats a 10% gain and a 10% loss symmetrically in log terms, and it is the mean whose result does not depend on which method is the reference. The arithmetic mean of the same ratios, 1.013, is close here but drifts when ratios are spread. The summary is never enough on its own: per-design results must be shown, because an aggregate gain can hide one serious regression.

## Certificates and invariants

Some checks are pass or fail and independent of the solver's own claims. An LP-based legalizer's answer is optimal if a dual-feasible solution has the same objective: a cheap, independent proof. A convex QP sizing step is certified by KKT residuals, stationarity, feasibility and complementarity, below a tolerance. A MILP floorplanner with incumbent 1040 and best bound 1000 is provably within 4% of optimal.

Invariants come before QoR. Every placer output must be legal: no overlaps, cells on rows and sites, fixed cells unmoved, whatever its wirelength. Determinism is an invariant of the harness: an annealer that gives different results on two runs with the same seed breaks differential testing, because a QoR change can no longer be attributed to the code change. A change meant only to speed up branch-and-bound node selection must leave the optimal objective unchanged on every exactly solved golden instance, with fewer nodes or less time.

## Promotion decisions

A move-set skill that improves wirelength on 9 golden designs and worsens one by 6% is neither promoted blindly nor rejected outright. The regression must be explained and bounded, and the manifest may exclude that design class so selection stays conditional. New code carries more risk than a parameter, so a Tier 2 plug-in needs broader review and test scope than a Tier 1 weight change, even when both pass the same golden tests.

Differential checks prove intent on references: a placer change that moved only the intended cells by the intended amounts has shown that its effect matches its purpose on those designs, not that it is optimal or that every future design benefits. A Lagrangian-relaxation sizer that passes golden tests but shows a 15% duality gap on one design may be far from optimal there, a lead for investigation. And every check runs in a sandbox with time and memory limits, so a runaway change fails cleanly instead of hanging the harness.
