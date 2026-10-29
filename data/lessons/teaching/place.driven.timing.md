# Timing-driven placement

A wirelength-driven placer treats all nets alike, so critical paths end up no shorter than any other. Timing-driven placement fixes this by weighting nets according to their timing criticality, computed from static timing analysis on the current placement. Net weighting is simple, works with every placement algorithm, and is what OpenROAD's gpl and RePlAce use in timing-driven mode.

## Criticality and net weights

After a timing run, each net gets its worst slack $s$, the minimum slack over paths through it. A criticality maps slack to [0, 1]: one common form is $c = \min(1, \max(0, -s/|WNS|))$, which is 0 for nets with non-negative slack and 1 for nets on the worst path. The net weight is then $w = 1 + \alpha c^\beta$, with $\alpha$ setting the maximum boost and $\beta$ concentrating it on the most critical nets. With WNS = -125 ps, $\alpha = 2$ and $\beta = 3$, a net with slack -55 ps has $c = 0.44$ and $w = 1 + 2 \cdot 0.44^3 = 1.1704$. Many variants exist; the shape matters less than updating the weights as timing changes.

## Weighted wirelength and its effect

The placer then minimises $\sum_e w_e \mathrm{WL}_e$. In quadratic placement, a heavier spring pulls harder: a movable cell between a fixed pin at x = 3 on a critical net and a fixed pin at x = 74 on a non-critical net sits at the midpoint 38.5 when weights are equal, and at $(1.5625 \cdot 3 + 74)/2.5625 = 30.71$ when the critical net weighs 1.5625, a shift of 7.79 um toward the critical pin. In nonlinear placement the weighted gradient does the same. The effect is to shorten critical nets at the expense of non-critical ones, ideally without much change in total wirelength, since most nets have slack to spare.

## Keeping weights current

Criticalities change as the placement changes, so timing-driven placers re-run timing periodically, for example at fixed overflow milestones, and update weights, often blending new with old weights to avoid oscillation: a net heavily weighted becomes short, gains slack, loses weight and drifts long again. Early in global placement, wire estimates are poor and timing is not yet meaningful, so the first updates come after cells have spread. Path-based methods, which constrain the delay of specific critical paths directly, are more precise but harder to scale; net weighting is the robust default. The weighting function, its parameters and the update schedule are change-surface parameters, while the timer it calls is a separate engine behind an interface.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 8.3 on timing-driven placement; the OpenROAD documentation for gpl's timing-driven mode.
