# Observation hooks

Observation hooks are the fourth layer of an agent-first engine: defined points where an agent reads live state and, at a few of them, injects behavior. They are how a fast modeling loop watches an engine converge, how it decides to stop early or switch strategy, and how it gathers the metrics it uses to select the next skill. Their design decides how much an agent can see and how much that seeing costs.

## Reading live state

A hook exposes a snapshot of engine state at a well-defined moment: every $K$ iterations of global placement, after each rip-up and reroute pass, after an incremental timing update. Useful state includes overflow and density for a placer, congestion and DRC counts for a router, worst and total negative slack for a timer, and IR drop for a power analyzer. The snapshot must be consistent, taken at a point where the engine's data structures agree with each other, and it must be read-only by default, so an observer cannot corrupt the run. Hooks report numbers and identifiers, not prose. A placer hook that returns iteration, HPWL, overflow and the ten densest bins lets an agent see that overflow has plateaued at 18% and stop wasting iterations.

## The cost of a hook

Observation is not free. If an engine iteration costs $t$ ms and a snapshot costs $c$ ms, a hook firing every $K$ iterations adds a fraction $c / (K t)$ to runtime. With a 40 ms placement iteration and a 120 ms density snapshot, firing every 10 iterations adds $120 / 400 = 30\%$, while every 50 iterations adds $6\%$. An agent that wants frequent snapshots should ask for cheap ones, such as scalar metrics, and request expensive ones, such as a full congestion map, only when the cheap metrics say something changed. The engine designer makes this possible by offering hooks at several levels of detail.

## Injecting behavior safely

Some hooks accept behavior: a callback that can stop the run, change a parameter between phases, or adjust net weights before the next iteration. These are powerful and must be bounded. Injection happens only at defined points between phases, never in the middle of the hot kernel. What may be injected is restricted to the change surface: an injected weight change is a Tier 1 change, an injected cost function is a Tier 2 plug-in, and no hook may alter legality. Every injection is logged with its values so the run stays reproducible, and the engine re-checks its invariants after an injected change, so a bad callback produces an error rather than a silently illegal result.

## Going deeper

The OpenROAD documentation for gpl, on its iteration reporting; Kahng et al., "VLSI Physical Design" (2nd ed., 2022), on global placement convergence.
