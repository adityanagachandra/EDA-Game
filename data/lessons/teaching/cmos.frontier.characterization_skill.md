# Characterization as an agent skill

Characterization is a natural first target for agent-driven engineering: it is compute-heavy, rule-bound and full of library-specific choices, and its outputs can be checked against SPICE. Packaging a characterizer as a skill, with a fixed measurement core, pluggable strategies, configuration and golden tests, is what lets an agent tune it safely for a new library without rewriting the parts that define what a Liberty number means.

## The four layers of a characterizer

Split the characterization engine as the four-layer model prescribes. The immutable core holds what defines correctness: deck generation from netlist and models, the measurement definitions (threshold crossings, the setup push-out criterion) and the Liberty writer. The strategy layer holds interchangeable policies: how to search for setup and hold (bisection, or a model-guided search), how to choose sensitizing side-input states, how to schedule simulations. The parameter layer holds per-library values: index grids, thresholds, corner lists, tolerances. Observation hooks report per-arc progress, convergence and measurement failures as runs complete. A request to change grid sizes then never touches code, and a new search policy never touches the definitions.

## Change tiers for library requests

Each customer request lands at the smallest safe tier. A change from 7 by 7 to 9 by 9 tables, a new corner, or different transition thresholds is Tier 1, configuration only. A constraint search that reuses neighbouring corners' results as starting guesses is Tier 2, a plug-in at the search extension point; it changes cost, not definitions, and must reproduce the same numbers within tolerance. A change to the setup criterion, from a 10% push-out to a pass-fail boundary, is Tier 3 even though the edit is one line, because it changes what every constraint value in the library means and shifts every downstream slack. The skill's manifest should say exactly which parameters are tunable and which definitions are fixed, so an agent selecting and configuring the skill never needs to read its source.

## The harness and its evidence

An agent editing the skill works inside the harness loop: propose a small change, build incrementally, run the golden tests, compare the produced library against a baseline, feed failures back, and present a passing change as a diff for human review. For a characterizer the golden tests are fixed cells, stimuli and expected tables checked against trusted SPICE within tolerance, and the baseline comparison is differential: every table of every cell, not only the goldens. A plug-in that is three times faster and passes the goldens but shifts one cell's hold table by 4 ps must be blocked until the shift is explained. Runs execute in a sandbox with no network, because foundry models and customer netlists are protected IP that must not leave the site.

## Going deeper

The agent-first EDA material in this curriculum (four-layer engine model, change tiers, harness loop); the Liberty Reference Manual for the quantities a characterizer must preserve.
