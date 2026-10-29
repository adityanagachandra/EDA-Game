# Fast and accurate engines across the flow

Every analysis in the flow exists at several accuracies: a fast placer and a full placer, a quick congestion estimate and a global router, a fast timer with lumped wires and a signoff timer with extracted parasitics. An agent exploring a design can run the fast versions continuously and recurse into the accurate ones only when needed. The arithmetic of that trade decides how many experiments fit in a day and when a fast result can be trusted.

## The cost of a modeling loop

Suppose an agent runs $N$ iterations of a fast loop, each taking $t_f$, and every $k$-th iteration it confirms the result with an accurate run taking $t_a$. The total time is $N t_f + \lfloor N/k \rfloor t_a$. With $N = 40$ iterations of 6 minutes and an accurate run of 90 minutes every 10 iterations, the loop takes $240 + 4 \cdot 90 = 600$ minutes. Running all 40 iterations accurately would take 3600 minutes, so the mixed loop is 6 times faster. The ratio $t_a / t_f$ and the confirmation interval $k$ set the speedup.

## When to recurse to the accurate engine

A fast engine is useful only if its error is bounded and known. If a fast timer is calibrated to be within $e$ picoseconds of signoff, then any endpoint whose fast slack is at least $e$ is safe, any endpoint below $-e$ surely fails, and endpoints with fast slack between $-e$ and $e$ are uncertain and must be checked by the accurate timer. With $e = 15$ ps and fast slacks of $-40, -12, 3, 9, 22, 60$ ps, three endpoints ($-12, 3, 9$) need the accurate check. The rest can be decided from the fast run.

## What the fast mode gives up

Fast modes trade fidelity for speed in specific ways: a fast placer may skip detailed legalization, a congestion estimator ignores detailed pin access, a fast timer lumps wire parasitics and skips crosstalk, and a surrogate model may be inaccurate outside its training distribution. Each trade is safe only within a calibrated regime. An agent must therefore track not only the metric but its variance and the regime: if the design moves into territory where the fast engine has not been correlated, such as a new layer stack or a very congested region, it must recurse to the accurate engine more often or recalibrate.

## Going deeper

The agent-first EDA notes in this curriculum, section on fast and accurate engines in a continuous loop.
