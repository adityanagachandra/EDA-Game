# Change tiers: config, plug-in, core

Every request to change an engine lands in one of three tiers, and the harness treats them very differently. Tier 1 is configuration, Tier 2 is a plug-in at an extension point and Tier 3 is an edit to the core. The model's first job on any request is to find the smallest safe tier that can express it, because that choice sets the cost, the risk and who must approve.

## Tier 1: configuration

A Tier 1 change edits a value or a table and needs no recompile. Selecting an existing skill, changing a placement density target from 0.65 to 0.70, raising the timing weight, choosing the POCV provider at 3.5 sigma instead of 3, pointing the AOCV provider at a customer's tables, telling the router to prefer metal 5 for a bus, or asking the timer for 20 paths per endpoint instead of 5 are all Tier 1. They are reversible by restoring the old configuration, they cannot break legality because the core validates every parameter against its declared range, and most requests land here. A configuration change is still reviewed and versioned, but it needs no build and only a QoR check.

## Tier 2: plug-in at an extension point

A Tier 2 change adds a small implementation at a defined seam when no parameter can express the behavior. Examples: a detailed-placement rule that abuts the cells of one bit slice; a router cost function that penalizes vias near a sensitive analog block; a derate provider that adds 3% late derate on HVT data paths with logic depth under 5. The core is untouched, the plug-in builds as its own unit in seconds, and the harness runs golden and differential tests before a human reviews the diff. If the plug-in proves generally useful, it is later promoted into the compiled skill library.

## Tier 3: core edits

A Tier 3 change edits the core: legalization, the DRC path model, timing propagation, the field-solver kernel, or an extension-point interface itself. It is rare and needs a human-approved diff and full regression. Some Tier 3 requests are legitimate, such as extending a seam's context with a field no plug-in can otherwise see; then one reviewed core edit turns many future requests into Tier 2. Others should be refused, such as letting cells overlap. Two rules follow. Always try the lower tiers first: a request that sounds like a core change ("route clocks differently") is often a parameter. And if one kind of request keeps needing Tier 3, the seam is in the wrong place.

## Going deeper

The agent-first EDA notes in this curriculum, sections on change tiers and the derating example.
