# Agents for board constraint setup and placement

A large share of board-design effort goes into setting up constraints (net classes, matching groups, impedance targets, keep-outs) from datasheets and design guides, and into interactive placement. Both are natural agent tasks, because they are reading, translating and iterating against checkable rules. Applying the four-layer engine model and the change tiers to a board placer and its constraint system shows where an agent may act, what it must not touch, and what evidence should gate its changes.

## The four layers of a board placer

Split a board placement and constraint engine as the four-layer model prescribes. The **immutable core** is legality: courtyard non-overlap, keep-outs, board outline, the clearance model. The **policy** layer holds the search: move sets, ordering, the cost function's structure, the router's rip-up strategy. **Parameters** are values: the weight between ratsnest length and a thermal margin, a matching tolerance, an impedance target. **Observation hooks** expose live state: unrouted nets, congestion by region, length-match errors, thermal estimates. An agent reads the hooks, edits parameters and selects policies; it never edits the legality core. Keeping the core small and the rule interface general is what lets most customer requests avoid a recompile.

## Change tiers for constraint requests

Customer requests land in tiers. "Match every DQ in a byte lane within 5 ps of its strobe" is **Tier 1**: a value in an existing length-match table. "No signal trace within 3 mm of any switching inductor's outline, on any layer" is **Tier 2** if the rule engine lacks that rule kind: a plug-in checker at the rule extension point, written, compiled and tested by the harness without touching the core. "Change how the legalizer resolves courtyard overlaps" is **Tier 3**, a core edit needing human approval and full regression. An agent that extracts constraints from a datasheet produces Tier 1 values, but each is a claim: it must be checked against the source table and shown to a human, since one wrong tolerance yields a board that passes every check and fails in the lab.

## Evidence for accepting agent changes

An agent's placement is accepted on evidence from deterministic tools, not on its own metric. A placement that cuts ratsnest length by 20% must also be courtyard and keep-out clean, meet critical-net constraints after a trial route, and be no worse than the baseline on thermal and assembly checks. Agent edits to constraints should be small, versioned and diffable, so a reviewer sees exactly which rules changed and can revert one. Fast and accurate engines combine in the loop: quick placement and a congestion estimate prune candidates, and a full autoroute with DRC and SI confirms the finalists. Some changes are out of bounds: relaxing a clearance rule to reach completion changes what legal means, so the harness must block it and escalate to a human.

## Going deeper

The agent-first EDA notes in this curriculum (the four-layer model and change tiers). Vendor documentation of board constraint managers, for the rule types an agent would configure.
