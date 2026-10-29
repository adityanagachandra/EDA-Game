# A match-cost seam for library-specific mapping

In an agent-first mapper, the covering algorithm is fixed, but how a match is priced varies by customer, library and objective. The design move is a seam where the dynamic program asks one question: what does this match cost in this context? Placing that seam well lets agents and users change mapping behaviour at the configuration or plug-in tier without touching cut enumeration, matching or legality.

## Core versus change surface in a mapper

The immutable core of a mapper guarantees a legal, equivalent cover: AIG construction, cut enumeration, canonicalization and Boolean matching, the covering DP and cover selection. The change surface is everything that only ranks legal choices: area and delay weights, per-cell penalties, don't-use lists, power terms, load estimates. A cost-provider interface separates the two. The core enumerates the legal matches at a node and calls the provider once, coarsely, with the batch; the provider returns a cost for each. Cost policy can then never break equivalence, because it never decides what is legal.

## What the seam's context must carry

A provider can only express policies its context supports. Useful context for a match: the cell and its drive strength and Vt class, the pin assignment from the matching transformation, input arrival times, the estimated load and required time at the node, the node's estimated fanout, and the mapping pass (delay, area flow, exact area). With that, a request such as "penalize high-Vt cells on paths with less than 50 ps slack" is a small plug-in that reads slack and Vt class. Leave out layout code, global state and anything nondeterministic, so plug-ins stay testable and reproducible.

## Keeping the seam honest

Three rules keep a cost seam safe. Costs rank and never forbid: legality constraints, such as forbidding matches across a clock-gating boundary, need a separate match-filter seam, because a huge cost still loses to nothing when no alternative exists. The interface is narrow and frozen, so plug-ins compile separately and can be swapped by configuration, for example a fast linear delay model during exploration and a Liberty-table model for the final pass. Every provider ships with golden tests: equivalence of the mapped netlist plus a differential QoR report showing that the matches that changed are the ones the policy targets.

## Going deeper

The OpenROAD and ABC source code for examples of cost hooks in mapping and sizing; Gamma, Helm, Johnson and Vlissides, "Design Patterns" (1994), for the strategy pattern behind provider seams.
