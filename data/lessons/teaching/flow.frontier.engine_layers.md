# The four-layer model across flow engines

An agent that changes EDA tools on request needs to know which part of an engine a request touches. The four-layer engine model answers that for every stage of the flow: an immutable core of correctness and legality, a strategy layer of heuristics, a parameter layer of settings, and observation hooks where an agent reads state. Applied across the flow, it shows that most customer requests never need to touch an engine's core.

## The same four layers in every engine

Every flow engine splits the same way. A placer's core is legalization: non-overlapping cells on legal rows and sites. Its strategy is the objective and move set; its parameters are density targets and net weights; its hooks expose density and congestion maps. A router's core is the DRC-correct path model on a legal grid; its strategy is the cost function and rip-up order; its parameters are layer directives and penalties. A timer's core is graph propagation and the delay-model interface; its strategy includes path pruning and derating schemes; its parameters are corners and margins. A synthesizer's core is function-preserving transformation; its strategy is the sequence of optimization passes; its parameters are effort levels and targets.

## Classifying a request

To classify a request, ask what it changes. Changing a correctness rule (what counts as a legal placement, how arrival times combine) is a core change. Changing how the engine searches (a new cost term, a new move, a new rip-up order) is a strategy change. Changing a number or mode (target density 0.65 instead of 0.7, an extra corner) is a parameter change. Reading state during a run (export the congestion map every iteration) is a hook. A request to weight timing-critical nets three times more in placement is a parameter; a request to add a new congestion-aware cost term is a strategy plug-in; a request to allow overlapping cells is a core change that should be refused or escalated.

## Why the core must stay small

Each layer has a different cost of change. A parameter change needs no rebuild and is checked by re-running the flow. A strategy change needs a build of one plug-in and golden tests. A core change needs human review and full regression, because every design depends on it. Engines designed so that variation lives outside the core let an agent satisfy most requests quickly and safely. Engines whose heuristics are tangled into the core force core edits for routine requests, which is slow and risky. The layer boundaries are therefore an architecture decision with direct consequences for how automatable a flow is.

## Going deeper

The agent-first EDA notes in this curriculum (the four-layer engine model and change tiers); the OpenROAD source tree, whose engines expose many parameters through Tcl.
