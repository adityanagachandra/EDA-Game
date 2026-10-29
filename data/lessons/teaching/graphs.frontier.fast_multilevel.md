# Fast multilevel modes in an agent's modeling loop

An agent exploring floorplans or clusterings needs to evaluate many candidates quickly, so it runs a fast mode of the partitioner in a continuous loop and recurses to the accurate mode only when the decision demands it. Building a useful fast mode, and a policy that knows when to trust it, is a concrete instance of the fast-versus-accurate engine idea.

## What a fast mode gives up

A multilevel partitioner has natural speed knobs, all of them parameters or policies: one V-cycle instead of several, a larger coarsening stop size, fewer initial-partitioning tries, and lighter or no FM refinement at the finest levels, where most runtime goes. Skipping refinement at the two finest levels might cut runtime severalfold and raise the cut by a few percent.

What a fast mode must not give up is the core: the same cut definition, the same balance rule and exact bookkeeping. If the fast mode measured a different objective or allowed illegal partitions, its scores would not predict the accurate engine's results, and the agent would optimize the wrong thing. A fast mode built from parameters is a Tier 1 change; removing the balance check for speed is not a fast mode, it is a bug.

## Deciding whether to trust it

The agent judges the fast mode by its bias and its variance. Bias: compare its mean cut over several seeds with a recent accurate reference. Variance: the spread across seeds. A policy might keep the fast mode only if its mean is within 5 percent of the reference and its seed spread is at most 30 nets.

Five fast runs of 953, 953, 943, 943, 951 against a reference of 884 have mean 948.6, a gap of $(948.6 - 884)/884 \approx 7.3$ percent, and a spread of 10. The spread passes but the gap fails, so the loop recurses to the accurate mode. For ranking candidates, what matters most is rank correlation between fast and accurate results: a biased but well-ordered fast mode can still choose well.

## Budgeting the loop

Fast and accurate runs share one compute budget. With a 2-second fast mode, a 60-second accurate mode, 50 candidates and 600 seconds, screening all 50 fast (100 seconds) and running the best 8 accurately (480 seconds) uses the budget well. Running accurate on 10 random candidates wastes the screen; trusting the fast winner alone risks picking a candidate whose fast score was within the proxy's error.

Two rules follow. When candidates differ by less than the fast mode's typical error, treat them as tied and let the accurate engine separate them. And recurse to the accurate mode before committing a decision that placement and signoff will build on, while periodically re-running accurate references so the fast mode's bias is recalibrated as the design changes. Learned surrogates fit the same loop, with the added risk of being confidently wrong on unfamiliar designs.

## Going deeper

Karypis and Kumar's hMETIS manual on its speed and quality options; the KaHyPar documentation on its preset configurations.
