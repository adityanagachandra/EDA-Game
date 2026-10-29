# Fast placement inside a modeling loop

Placement exists in fast and accurate forms: a GPU global placer with loose stopping rules and RUDY congestion, versus full global and detailed placement followed by global routing. An agent exploring floorplans or parameters runs the fast version in a loop, uses its metrics and their error to decide which candidates deserve the accurate version, and bases final decisions only on accurate results.

## What the fast placer trades away

A fast placer stops at a higher overflow, uses coarser bins, skips or shortens detailed placement and estimates congestion with RUDY instead of a router. It returns HPWL, overflow and a congestion map in a fraction of the time, with a known error relative to the accurate flow. That error must be measured, not assumed: run both on a sample of designs and record the spread of the ratio, for example HPWL within about 3%. The fast placer must also be deterministic and seedable, so that differences between candidates come from the candidates rather than from noise, and its error should be re-measured whenever the design family changes.

## Deciding what to run accurately

Given fast estimates and an error band of $\pm e$, a candidate can be discarded only if even its optimistic value, estimate times $(1 - e)$, is worse than the best candidate's pessimistic value, best estimate times $(1 + e)$. With estimates 934, 910, 905, 990, 1059, 952, 1048 and 944 mm and $e = 2\%$, the best pessimistic value is 923.1 and only candidates whose optimistic value is at most that survive: 905, 910 and 934, three candidates for the accurate placer. If many candidates survive, the error band is too wide for the decision and the agent needs a better fast model or more accurate runs. Metric variance across seeds plays the same role as the error band.

## Pitfalls of the loop

An agent rewarded on fast-placer metrics can learn to exploit the fast model's blind spots, for example floorplans whose RUDY estimate looks clean but whose real routing congests near macros, a form of reward hacking. Guards include adding routability terms that correlate with the accurate flow, periodic accurate audits of chosen candidates, and making final decisions only from accurate results. Speed shifts the bottleneck: when global placement takes a minute, I/O, legalization and evaluation dominate, so the loop is only as fast as its slowest stage. The fast placer, its error model and the selection rule together form a skill whose manifest should state when its estimates can be trusted.

## Going deeper

The agent-first EDA notes in this curriculum, section on fast and accurate engines in a continuous loop; Lin et al., "DREAMPlace" (DAC 2019), on GPU placement speed.
