# Constraint-driven layout agents

An analog layout agent turns designer intent ("match these, keep that away from the output") into constraints, drives a layout engine and checks the result. Applying the four-layer model and change tiers to the layout engine shows where each request belongs, where the seam for team-specific matching rules should sit, and what must be proved before any change is accepted.

## Layers and change tiers of a layout engine

The immutable core of an analog layout engine is design-rule legality and connectivity correctness of the generated geometry: whatever else changes, the GDS must be DRC clean and match the netlist. Strategy and policy are the placement and routing algorithms and the pattern generators. Parameters are weights (symmetry versus wirelength), keep-out distances, shielding lists. Observation hooks expose congestion, symmetry violations and extracted parasitics. A request such as "keep the input pair 20 micrometres from the output stage" is Tier 1: a value in the constraint file. A new interdigitation rule chosen by finger count is Tier 2: a plug-in at the pattern-generator seam. Changing legalization is Tier 3, rare and human-approved. A well-placed seam makes Tier 3 almost never necessary.

## The constraint-provider seam

Matching policy is where teams differ most, so it belongs behind a narrow interface. A constraint provider receives the netlist with device types, sizes and terminal connectivity, plus designer annotations such as net roles, and returns symmetry pairs, matched groups with patterns, proximity and keep-out constraints, and sensitive-net rules. The placer and router consume only that output. A plug-in can then recognize a team's differential-pair conventions or apply a house common-centroid style without touching the engine. The seam must expose enough context to recognize structures (topology and sizes, not just device counts), and must be deterministic, so the same netlist always yields the same constraints and differential tests remain meaningful.

## Verification, fast estimates and the flywheel

Before a new policy plug-in reaches a human, the harness runs DRC, LVS and post-layout simulation on golden blocks and compares them with baseline layouts. A differential check must show that only the intended blocks or metrics changed and that no golden block lost a spec. Inside its search loop an agent may use a fast parasitic estimator, but it calls full extraction and simulation before accepting a layout and whenever an estimate is near a spec boundary. Rewards must include matching and post-layout specs; "minimize area subject to DRC" produces compact layouts that break performance. A pattern that works on one block can be validated and promoted into the library of pattern generators with a manifest advertising, for example, supported ratios and which gradients it cancels, so other agents can select it.

## Going deeper

The agent_first_eda lessons on change tiers, seams and extension points, and the skill flywheel.
