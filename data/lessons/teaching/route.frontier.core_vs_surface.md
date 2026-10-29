# Router core versus change surface

An agent that customises a router must know which parts it may change freely and which it must never touch. The four-layer engine model gives the split: an immutable core that guarantees legal output, strategies that decide how to search, parameters that tune them, and hooks that expose live state. For a router, legality is the core and almost everything a customer asks for is surface.

## The four layers of a router

The core of a router is its DRC-correct path model and grid and layer legality: what geometry is allowed, how shorts and spacing are checked, how vias connect layers. Changing it needs human approval and full regression. The strategy layer holds the cost functions, the rip-up and reroute policy, net ordering and the search algorithm, selected from a library or added as plug-ins. Parameters are numbers such as via cost, history increment, present-factor growth, layer adjustments and iteration limits. Observation hooks expose the congestion map, violation markers and per-net state to an agent.

Example: an agent reading the congestion map after every global-routing iteration uses a hook; changing the history increment is a parameter; adding a new rip-up policy is a strategy plug-in; allowing off-grid wires is a core change.

## Why legality stays in the core

If each cost plug-in had to check shorts and spacing, one buggy plug-in could make the router emit an unmanufacturable layout. So legality is enforced by the core regardless of which strategy is active, and the seams between core and strategy carry contracts the core can enforce: costs must be finite and non-negative (search optimality depends on it), deterministic for the same inputs, and computed only from the context the seam passes.

The seam must expose enough context to make useful policies possible: the edge, its layer and direction, demand, capacity and history, the net and its criticality. A seam that hides criticality cannot host a timing-aware cost; one that exposes the whole database as text invites plug-ins that depend on internals and break when the core changes.

## Performance and determinism at the seam

Router kernels are hot loops: a maze search expands millions of states. A plug-in called through a virtual function per edge expansion can slow the search several times and stop the compiler from optimising it. Seams should therefore be coarse: a plug-in computes a cost map per region or per iteration, which the kernel reads as a table.

Plug-ins must also be deterministic with no hidden global state, so the same design routes the same way twice; otherwise golden and differential tests are meaningless. A customer preference such as making vias near an analog block ten times more expensive is a strategy or parameter change, not a core edit. A change to what geometry is legal is the rare core change, and is treated as one.

## Going deeper

The OpenROAD drt and grt source structure as a case study.
