# Seams and change tiers in a simulator

A simulator is an engine like a placer or a timer, and the four-layer model applies to it: an immutable core, policies, parameters and observation hooks. Placing each feature in the right layer decides which customer requests are configuration, which are plug-ins, and which would touch the semantics that every result depends on. The seams must also be shaped so that extension does not destroy the compiled speed that makes a modern simulator useful.

## The four layers of a simulator

The **core** is what defines correctness: the scheduling semantics of the language standard (the regions of a time step and their order), the four-state value system and its operator tables, and the event queue's invariants. The **policy** layer holds alternative strategies on top of that core: X-propagation modes (standard semantics versus more pessimistic merging), logic partitioning and ordering heuristics for multithreading, waveform compression schemes. **Parameters** change behavior without code: random seeds, thread counts, optimization levels, plusargs, dump scope. **Observation hooks** let tools read and act on live state: VPI callbacks on value changes and time, coverage and assertion monitors, waveform dumpers. An agent reads state through hooks and changes behavior through parameters and policies.

## Classifying change requests

Requests map onto tiers. Dumping waveforms for one subsystem only is **Tier 1**: a dump-scope setting. A custom X-propagation policy, say one that merges both data inputs of a mux whose select is X, is **Tier 2** if the simulator has an X-propagation seam: a new policy implementing a defined interface, with the value system untouched. A new coverage metric collected during simulation is also Tier 2, a monitor attached at the hooks. Reordering the Active and NBA regions for speed is **Tier 3**, a core edit: it changes language semantics, so it needs human approval and the full regression, and in practice it should be refused. A seam is placed well when the common requests in this list never need Tier 3.

## Designing seams that keep compiled speed

A compiled simulator is fast because the evaluation of the design is straight-line, inlined code. A seam that calls a plug-in on every signal event would destroy that. Seams are therefore coarse: a policy is consulted when code is generated or once per operation type, callbacks fire per cycle or per batch of changes, and plug-ins receive plain data structures. An X-propagation seam needs enough context to decide: the operator, the operand values including their X bits, and where the operation sits (an if condition, a case selector, a mux select). Plug-ins must be deterministic and seedable, with no hidden global state, so that golden and differential tests can tell a regression from noise. A new policy is promoted only after differential runs against the existing simulator on designs that stress ordering and X behavior.
