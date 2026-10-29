# The placer's core and change surface

An agent-first placer is built so that most change requests never touch its hardest code. The four-layer model splits it into an immutable core (legalization, the slot model, the numerical kernels), strategy (move sets, spreading methods, net models), parameters (weights, targets, schedules) and observation hooks (live metrics for agents). Knowing where each placer feature belongs decides how cheaply and safely it can change.

## What belongs in the placer's core

The core holds the invariants and the kernels every result depends on: the row and site slot model, the legality checker (no overlap, site alignment, rail orientation, regions), the legalizer that guarantees those properties, and the numerical kernels of global placement such as the density accumulation, the Poisson solve and the wirelength operators' interfaces. These change only for a new technology or a new kind of cell architecture, need human review, and run the full regression. Everything downstream, from CTS to routing to signoff, assumes the core's output is legal, so a bug there corrupts the flow. Keeping it small, deterministic and heavily tested is the design goal.

## Strategy, parameters and hooks

Strategy is how the placer searches: the net model (clique, star, B2B), the smooth wirelength (LSE, WA), the spreading method, the detailed-placement move set (swaps, reordering, matching), and rule plug-ins that add costs. New strategies are added as plug-ins at defined seams. Parameters are values the engine already uses: target density, overflow stopping point, timing and congestion weights, $\gamma$ schedules, inflation caps. Changing them needs no recompile. Hooks expose live state during a run, such as overflow, HPWL, density and congestion maps every N iterations, so an agent can watch, stop early or adjust parameters. Most real requests are parameter changes.

## Routing a request to a layer

A request is placed by asking what it needs that the engine does not already offer. "Use a target density of 0.65 in this region" fits an existing parameter: Tier 1. "Keep the cells of each bit-slice abutted" needs new logic that scores arrangements but does not change legality: a Tier 2 plug-in. "Let two cells overlap to save area" breaks the legality invariant: refuse or escalate. "Support rows with a new rail pattern" changes what legal means: a rare Tier 3 core change. If similar requests keep forcing core edits, a seam is missing, and adding it once turns future requests into plug-ins. Seams should take batches of cells, pass plain data and hold no hidden state, so the core stays fast and runs reproducible.

## Going deeper

The agent-first EDA notes in this curriculum, sections on the four-layer engine model and change tiers; the OpenROAD gpl and dpl documentation, for the parameters a real placer exposes.
