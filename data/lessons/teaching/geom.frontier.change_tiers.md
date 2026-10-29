# Change tiers for geometry-engine requests

Requests to change a geometry engine arrive constantly: a new process value, a new rule shape, a performance problem, a bug. An agent that edits the engine should land each request at the smallest safe tier: configuration, a plug-in at an existing seam, or, rarely and under review, the core. Classifying requests correctly, and knowing what evidence each tier needs, is the practical skill. Repeated core edits for one kind of request are a sign that a seam is missing.

## Tier 1: rule decks and parameters

Most requests change values that the engine already reads: a spacing from 40 to 38 nm, density windows of 50 um, a choice between built-in metrics, a worker count, a tile size. These are Tier 1: no code changes, no recompile, and the risk is limited to the value being wrong. The evidence needed is that the intended check results change and nothing else does, which a rerun on reference layouts shows. A parallelism setting is Tier 1 with a special obligation: results must not depend on it, which golden tests verify by running at several worker counts. An agent should always look for a Tier 1 solution first, because most customer differences are differences in numbers, not in algorithms.

## Tier 2: plug-ins at a seam

When a request needs new logic that fits an existing extension point, it is Tier 2. A spacing rule that depends on whether both wires share a net and on their summed widths becomes a new spacing provider, given the pair's widths, nets and run length. A fill policy that avoids analog regions subtracts keep-outs with core Booleans. A new sweep ordering for cache behaviour plugs into an ordering interface. The core is untouched, and the evidence is golden tests plus a differential run: for a new rule, only the expected new markers appear; for a performance policy, the output is identical to the baseline with zero XOR. The harness can then promote the plug-in into the library, where other agents can select it.

## Tier 3: touching the kernel

Tier 3 changes what the core represents or computes: curvilinear primitives for photonics, wider coordinates for very large dies, a fix for a Boolean bug at a vertex shared by three polygons. These go through a human-approved diff and full regression, and a bug fix always adds the failing input as a permanent golden case. Patching around a kernel bug in a plug-in hides it from every other caller and is the wrong answer. Even changes that look small, such as widening the coordinate type, are Tier 3 because they alter overflow bounds, memory layout and every interface that carries coordinates. When the same kind of request keeps reaching Tier 3, the varying part should move behind a new seam so the next one is Tier 2.

## Going deeper

The four-layer engine model and change tiers taught in this curriculum's agent-first EDA domain; Fowler, "Refactoring" (2nd ed., 2018), on introducing seams for change.
