# The four-layer engine model

An engine that agents, field engineers and customers will customize must separate what may change from what may not. The four-layer model does this for every engine, implementation or analysis: an immutable core, a strategy or policy layer, parameters, and observation hooks. Each layer has its own cost of change, from human-approved regression to a configuration edit. Knowing which layer a request lands in is the first decision in every customization.

## Four layers from hardest to easiest to change

The **immutable core** holds correctness, legality, invariants and the numerical kernel: overlap-free legalization in a placer, the DRC-correct path model in a router, graph propagation in a timer, the field-solver mathematics in an extractor. Changing it needs human approval and full regression. The **strategy or policy** layer holds the heuristics that decide how the engine searches: cost functions, move sets, orderings, rip-up strategies. Policies are selected from a library of compiled implementations, and new ones are added as plug-ins. The **parameters** layer holds weights, thresholds and modes, set by configuration with no recompile. The **observation hooks** are where an agent reads live state, such as timing, congestion, density or IR drop, and where it may inject behavior at defined points.

## Classifying a change

Most requests land in parameters. "Weight timing more heavily in placement" changes a weight. "Target 70% density instead of 65%" changes a threshold. "Try a different move set in detailed placement" selects or adds a policy. "Stop global placement when overflow falls below 10% and report congestion every 50 iterations" uses hooks. "Allow two cells to overlap by half a site" attacks legality, so it is core, and it should almost always be refused. The test is what must stay true afterwards: if the change could make a result illegal or wrong rather than merely better or worse, it touches the core.

## The design rule: a small core

The rule is to keep the core small and rarely touched and to push everything that varies out into policy, parameters and hooks, so that the common change never recompiles the core. A core that is too large forces routine requests into reviewed core edits, which is slow and risky. A core that is too small exposes legality to plug-ins that can break it. The right line follows correctness: whatever guarantees a legal, correct result belongs in the core; whatever trades one quality metric against another belongs outside it. An engine whose change log shows many core edits for routine tuning has its line in the wrong place.

## Going deeper

The agent-first EDA notes in this curriculum, section on the four-layer engine model; Gamma et al., "Design Patterns" (1994), on the strategy pattern.
