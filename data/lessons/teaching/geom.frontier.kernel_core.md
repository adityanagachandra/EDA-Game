# The geometry kernel as a protected core

Every DRC, extraction, fill and LVS engine rests on a small geometry kernel: exact Booleans, sizing, edge and polygon predicates on integer coordinates. In an agent-first engine that kernel is the immutable core, and everything that varies between customers and runs (strategies, rule values, metrics, tiling) sits in layers around it. Deciding what belongs in the core, and how the seams around it are shaped, determines whether agents can customize the engine safely and quickly.

## What belongs in the geometry core

The core holds what must be correct for every caller: Boolean operations on polygon sets, sizing, merging, the orientation and incircle predicates, coordinate transforms, and the invariants on their outputs (simple, non-overlapping, correctly oriented rings with all vertices on the grid). A bug here is a silent sign-off escape in every rule that uses the operation, so the core is small, exhaustively tested and changed only with human approval and full regression. Anything that is a choice rather than a fact stays out: which rule to run, which metric a rule uses, how a layer is tiled. A useful test is to ask whether two reasonable customers could want different behaviour; if so, it is not core.

## Policies, parameters and hooks around it

Around the core sit three layers. Policies are strategies selected from a library or added as plug-ins: global sweep versus tiled processing, the order in which rules run, a fill-placement heuristic, a candidate-pair filter. Parameters are values with no code change: rule distances, density windows, tile size and halo, worker counts. Observation hooks expose live state, such as per-tile violation counts, density maps and progress, which an agent reads to decide what to run next. Example: a team wanting fill to prefer regions away from critical nets writes a placement policy that calls the core's NOT and sizing operations; it never edits the Boolean kernel, so a mistake in it can misplace fill but cannot corrupt a Boolean result elsewhere.

## Keeping the hot loop inside the core

Geometry engines spend their time in tight inner loops over edges and candidate pairs. If an extension point sits inside that loop, say a virtual callback per edge pair, call overhead and lost inlining can slow a check by an order of magnitude. Seams must be coarse: the core hands a plug-in a batch of candidates as plain arrays and receives a batch of results, or calls a policy once per tile. Plain data, a frozen interface and no hidden global state let the core change internally without breaking plug-ins and keep plug-ins from mutating kernel state. A second implementation of the core, such as a GPU back end, is held to the same standard: identical results on the same input, proved by differential tests.

## Going deeper

The Boost.Polygon and KLayout documentation, for examples of a compact geometry kernel with layered APIs; Parnas, "On the criteria to be used in decomposing systems into modules" (Communications of the ACM, 1972).
