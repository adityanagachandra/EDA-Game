# Change tiers for timing requests

Requests to change a timer arrive constantly: a new sigma count, a different derate scheme, a margin rule from silicon correlation, a reporting change. An agent-first engine handles each at the smallest safe tier. Configuration changes need no recompile, plug-ins add behavior at a defined extension point, and core edits are rare and human-approved. Classifying requests correctly is what keeps the timer's core stable while it adapts quickly.

## Tier 1: configuration

A Tier 1 change alters a value, a table or a choice among shipped components, with no code change. For a timer this covers choosing the derate provider, setting the number of sigma for POCV, setting global early and late derates, pointing at AOCV or LVF table files, setting per-path-group margins, and choosing which corners and how many paths to report. Running signoff at 3.5 sigma instead of 3, switching from flat OCV to the shipped AOCV provider with the team's own tables, or adding a 10 ps margin to one path group are all Tier 1. These changes are reversible and need no rebuild, so most requests should land here if the parameters were designed well.

## Tier 2: a plug-in at the seam

A Tier 2 change adds a small implementation behind an existing extension point. A rule such as three percent more late derate on HVT data paths with logic depth under 5, falling through to AOCV for clocks, needs code, because no parameter combines Vt class, depth and clock-versus-data, but the seam already passes those fields, so a new provider implements it without touching the core. A margin layer that wraps the AOCV result with a silicon-correlation adjustment is another example. The harness builds the plug-in incrementally, runs the golden tests and a differential timing report, and promotes it only when exactly the intended paths moved.

## Tier 3: a core edit

A Tier 3 change edits the core: the propagation math, the slack conventions or the delay-model API. Changing how arrivals combine at a merge pin, or adding a field to the derate context, is Tier 3. Such a change affects every result, so it goes through a reviewable diff, full regression and human approval, even when the diff is small. If a request needs a field the seam does not carry, the honest path is one Tier 3 change to extend the context, after which similar requests become Tier 2. A seam is in the right place when Tier 3 edits for that kind of request essentially never happen; repeated core edits for derating mean the seam exposes too little.
