# Frozen extension ABIs, versioning and rollback

A seam is a contract between the core and every plug-in written against it. If the contract shifts with each core release, plug-ins break, promoted skills stop loading and customer sites cannot upgrade. Freezing the extension ABI, versioning every change and keeping rollback make plug-ins durable and every change reversible, which is what lets a deployment accept agent-written code at all.

## Freezing the extension ABI

An API is the source-level contract; an ABI is the binary one: symbol names, calling convention, the layout of structures passed across the boundary and the layout of virtual tables. A plug-in compiled as a shared library and loaded by the core depends on the ABI. Adding a field in the middle of a context structure, reordering virtual methods or changing a parameter type breaks every compiled plug-in, even if the source still compiles. Freezing the ABI means these never change within a major version. New capability is added compatibly: a new versioned interface beside the old one, or a structure that carries its own size and version so the core knows which fields a plug-in understands.

## Versioning every change

Every change, of whatever tier, is versioned: a configuration change gets a new configuration revision, a plug-in gets a version and a manifest update, and a core change gets a release. Each version records what changed, the request that caused it, the evidence that admitted it (golden, differential and QoR results) and who approved it. A run's provenance then names the exact versions it used. When a result looks wrong a month later, the team can tell which configuration and which plug-in versions produced it and reproduce the run exactly.

## Rollback and compatibility

Rollback is the ability to return to the previous known-good state quickly. For Tier 1 it means restoring the prior configuration revision; for Tier 2, unloading the new plug-in and reloading the old one; for Tier 3, reverting a release, which is why core changes need more evidence. Rollback works only if versions are immutable and old artifacts are kept, and only if a newer core can still load older plug-ins. A plug-in that declares the ABI version it was built for lets the core refuse an incompatible plug-in with a clear error, instead of crashing or silently computing wrong derates.

## Going deeper

Drepper, "How To Write Shared Libraries" (2011), on ABI stability and symbol versioning.
