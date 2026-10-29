# Routing a flow change request to a tier

Requests to change a flow arrive in plain language: make placement less congested near the memories, time this block with an extra margin, route this bus on upper layers. An agent that serves them must decide at which tier to implement each one: configuration, plug-in or core edit. The rule is to use the smallest safe tier, and the hard part is recognizing which tier a request really needs.

## Three tiers

Tier 1 is configuration: a value, a table or a mode switch, with no recompile. Setting a placement density target, choosing a derating table or adding a routing layer directive are Tier 1. Tier 2 is a plug-in: a small implementation at a defined extension point, such as a new cost function for the router, a new derate provider for the timer or a new optimization pass in a synthesis script, built and loaded without changing the core. Tier 3 is a core edit, a change to the engine's own source, which is rare and goes through a human-approved diff and full regression.

## Matching requests to tiers

A request maps to a tier by what it needs that the engine does not already offer. "Use 0.6 utilization in this region" fits an existing parameter: Tier 1. "Penalize routes that cross the analog keep-out by distance from its edge" needs new logic but fits the router's cost seam: Tier 2. "Add 3% late derate on data paths through high-threshold cells with depth under 5" needs a new rule at the timer's derate seam: Tier 2. "Support a new kind of placement legality for a novel cell architecture" changes what legal means: Tier 3. A request that seems to need Tier 3 often signals a missing seam, which is worth adding once so that the next such request is Tier 2.

## Proving a change at its tier

Each tier has its own proof. Tier 1 changes are proved by re-running the affected stages and comparing QoR with a baseline. Tier 2 changes also run the plug-in's golden tests and a differential check that only the intended behaviour moved. Tier 3 changes run the full regression across many designs, because the core affects everything. Choosing too high a tier wastes review time and adds risk; choosing too low a tier, for example hacking a parameter to imitate new logic, gives fragile results that break on the next design.

## Going deeper

The agent-first EDA notes in this curriculum, sections on change tiers and seams.
