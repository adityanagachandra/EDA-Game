# Agentic DRC fix loops

Signoff DRC on a large block can report thousands of violations after routing. An agentic fix loop reads the report, classifies each violation, routes it to the engine that owns the geometry, reruns the checks and iterates, with a human approving the final ECO. The engineering lies in classification and routing: most violations should go back to an engine as a small, safe change, not be patched by hand.

## Routing violations to the right fixer

Each violation has an owner. Spacing, shorts, minimum area and via enclosure on routed metal belong to the detailed router, which can rerun the affected region with the marker as a blockage or penalty: a Tier 1 action. Antenna violations go to the antenna-repair step, density windows to the fill step, and anything inside a standard-cell or macro layout is escalated, because the router treats those layouts as fixed. A report with 35 spacing (none inside cells), 10 short, 22 minimum-area and 4 enclosure violations sends 71 to the router. The policy is a table, so the agent's routing step is cheap and auditable.

## Classifying violations

Good classification looks for patterns before fixing instances. The same minimum-area violation on hundreds of nets after a router update points to a missing or wrong rule in the router's configuration: one Tier 1 config fix, not hundreds of patches. Markers for a rule the foundry removed point to a stale deck version. A violation inside vendor IP is outside the change surface. Rule names, layers, cell context and repetition across nets are the features a classifier reads, and all of them come from the results database, not from the deck's protected text.

## Closing the loop safely

An agent's change is accepted only when full signoff DRC and LVS pass and a differential run shows no new markers: DRC-clean and LVS-clean are the golden tests. Fixes go through the engines that own the data, never through direct polygon edits to the stream, which would make the design database, extraction and timing disagree with the layout. Reroutes change parasitics, so extraction and timing rerun before signoff. The loop is bounded: if the violation count stops falling (120, 40, 38, 41, 39), the agent stops and escalates the rest.
