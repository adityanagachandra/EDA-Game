# Fast and accurate solver modes

Every numerical engine can trade accuracy for speed: looser tolerances, lower precision, coarser meshes, larger timesteps, reduced models. An agent running a continuous modeling loop needs to know exactly what each fast setting buys and costs, how to tell whether a fast result can rank candidates, and when a result must be promoted to the accurate engine.

## Loose tolerances and what they buy

Many fast settings have predictable effects. CG iterations grow like $\ln(2/\varepsilon)$ at fixed $\kappa$, so relaxing the tolerance from $10^{-10}$ to $10^{-4}$ cuts iterations to about $\ln(2 \times 10^4)/\ln(2 \times 10^{10}) = 9.90/23.72 = 0.42$ of the original. Trapezoidal LTE scales like $h^3$, so loosening the timestep tolerance tenfold grows steps by $10^{1/3} = 2.15$ and cuts step count by about that factor. FP32 halves memory traffic for bandwidth-bound kernels, nearly doubling SpMV speed, but loses digits on ill-conditioned matrices.

Coarsening a thermal mesh by 2 in 3D cuts unknowns by 8 and, at the same tolerance, CG iterations by about 2. Each of these is a parameter or policy setting, Tier 1, so an agent can switch modes without touching the core.

## Error budgets for ranking candidates

An agent mostly uses fast results to rank: is layout A better than layout B? If each fast estimate is within $\pm e$ of the accurate value, A's true value can be $e$ higher and B's $e$ lower, so the ranking is trustworthy only if the gap exceeds $2e$.

Example: fast worst-case IR drops of 42.1 mV for A and 44.0 mV for B, with $e = 0.8$ mV. The gap 1.9 mV exceeds 1.6 mV, so A really is better. With $e = 1.2$ mV the combined band 2.4 mV exceeds the gap, and the agent must recurse. In practice $e$ is measured: run fast and accurate modes on a sample of candidates of the same design class and record the error distribution, or the rank correlation between the two.

## Promoting a fast result to signoff

Promotion means re-running the chosen candidate in accurate mode: FP64, tight tolerances, full mesh, deterministic reductions, and deciding on those numbers alone. Fast numbers are never reported as signoff results, even with a disclaimer.

Three checks guard the loop. First, the fast mode's error band must have been measured on designs like the current one; a reduced thermal model built for one floorplan can be badly wrong on another. Second, when fast results feed a learning agent's reward, periodically recompute rewards with the accurate engine: an agent that learns layouts the fast mode likes but the accurate engine does not is reward hacking the proxy. Third, keep the fast mode deterministic within itself, so that rankings are reproducible and differences between candidates are real.

## Going deeper

Fast and accurate engines, and tool-grounded reinforcement learning, in docs/eda/agent-first-concepts.md of this repository; Carson and Higham, "Accelerating the solution of linear systems by iterative refinement in three precisions", SIAM J. Sci. Comput. (2018).
