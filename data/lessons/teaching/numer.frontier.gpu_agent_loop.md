# GPU solvers inside agent loops

Agents exploring power-grid straps, decap placements or thermal floorplans issue hundreds of similar solves. GPU solvers make that loop fast, but only if the loop is organized around how GPUs work: batched, warm-started, bandwidth-bound solves, with a clear rule for when a fast answer is good enough and when to call the accurate engine.

## Batching many small solves

Each independent solve pays fixed costs: symbolic analysis, data transfer, kernel launches. When candidates share a sparsity pattern, those costs can be paid once. With setup $s$ and numeric work $t$ per candidate, $k$ unbatched solves cost $k(s + t)$ and a batch costs about $s + kt$.

Example: $k = 100$ variants of a PDN, $s = 20$ ms, $t = 4$ ms. Unbatched: 2400 ms. Batched: 420 ms, a factor of 5.7. Batching also fills the GPU: a single small solve leaves most of the device idle, while a batch of a hundred keeps it busy. Warm starts help further: consecutive candidates differ slightly, so starting CG from the previous solution makes the initial residual small and cuts iterations, often by a large factor.

## Budgeting candidate evaluations

An agent loop has a time budget. With total time $T$, fixed overhead $F$ (loading the design, building the base matrix), fast-mode cost $t_f$ per candidate, and the requirement that the best $m$ candidates be re-run in the accurate solver at $t_a$ each, the number of fast candidates that fit is
$$N = \left\lfloor\frac{T - F - mt_a}{t_f}\right\rfloor.$$
With $T = 30$ minutes, $F = 60$ s, $t_f = 2$ s, $m = 5$ and $t_a = 120$ s: $N = \lfloor(1800 - 60 - 600)/2\rfloor = 570$.

Forgetting to reserve the accurate re-runs is the common planning error: the loop finds a winner and then has no time to confirm it. The budget also shows the leverage of the fast mode: halving $t_f$ doubles the exploration.

## When to recurse to the accurate solver

A fast mode (FP32, loose tolerance, coarse mesh, reduced model) has an error band $e$. The agent recurses to the accurate engine when the fast answer cannot decide: two candidates closer than $2e$; a candidate within $e$ of a constraint limit, such as 49 mV against a 50 mV IR-drop limit with a 2 mV band; or a design unlike the one on which the band was measured, since error bands of fast modes are empirical and must be re-measured on a new design class.

And it always recurses before signoff. The fast mode guides exploration; the accurate engine, with validated tolerances in its deterministic mode, produces the numbers a human approves. Using fast-mode numbers as the reward for a learning agent invites reward hacking: the agent finds layouts the fast mode scores well and the accurate solver does not.

## Going deeper

Fast and accurate engines in a modeling loop, docs/eda/agent-first-concepts.md of this repository; Williams, Waterman and Patterson, "Roofline: an insightful visual performance model for multicore architectures", Communications of the ACM (2009).
