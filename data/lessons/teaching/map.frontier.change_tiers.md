# Change tiers for mapping requests

Requests to change a mapper arrive constantly: exclude these cells, raise K, penalize this Vt class, use our delay model, forbid matches across this boundary. Classifying each into configuration, plug-in or core edit decides who may make it, how it is verified and whether it needs a human. A well-designed mapper makes the common requests land in the first two tiers.

## Tier 1: configuration

A request that changes a value the engine already reads is Tier 1, with no recompile: a don't-use list of cells, the LUT size K, the number of priority cuts, the area-delay weight, a delay target, the number of area-recovery passes. Changing results is fine at this tier; what defines it is that no code changes. Configuration can still create infeasible problems: a don't-use list that removes every inverter leaves the library incomplete, so the engine should validate functional completeness of the remaining library before mapping.

## Tier 2: plug-ins at existing seams

A request that needs new logic, but only at an existing extension point, is Tier 2. A per-cell pin-access penalty from a customer table, a slack-dependent Vt penalty, or an in-house delay model behind the delay-model interface are all small implementations of an existing interface. They compile separately and leave the core untouched. Abusing a seam is a warning sign: implementing "forbid this match" as an enormous cost is not safe, because when no alternative exists the forbidden match still wins. That request needs a legality seam.

## Tier 3: core edits and seam design

A request that changes cut enumeration, matching, the covering DP or legality, with no seam for it, is Tier 3: a core edit with a human-approved diff and full regression. If an agent proposes a core edit for a request that a plug-in could serve, the harness steers it back to the smallest safe tier. Recurring Tier 3 requests of one kind mean a seam is missing: adding a match-filter extension point, for example, turns future legality requests into Tier 2. A seam is well placed when the common kind of change never needs a core edit, and every tier, including Tier 1, still passes equivalence and QoR golden tests.

## Going deeper

The agent_first_eda domain's treatment of change tiers and seams; the ABC source code, whose command options and plug-in points show which mapping changes are configuration today.
