# Core versus change surface per engine

The four-layer model becomes useful only when it is drawn concretely for each engine. For a placer, a router, a timer and an extractor, the engine's owners decide once what must stay fixed because correctness depends on it, and what the model, a field engineer or a customer may tune. That line is expert work, and every later customization rests on it.

## Placer and router

In a **placer**, the core is legalization: cells on rows and sites, no overlaps, orientation and row rules, blockages respected. Algorithms such as Abacus guarantee this, and nothing outside the core may weaken it. The change surface is everything that trades quality: objective weights for wirelength, congestion and timing, the density target, the move set of detailed placement, and design-specific rules such as abutting regular cells of a datapath. In a **router**, the core is the DRC-correct path model and grid and layer legality: a routed path must respect spacing, width, via and direction rules. The change surface is the cost and penalty functions, the rip-up and reroute strategy, layer directives and net ordering. Two routers with different cost functions can both be correct; a router that skips a spacing check is not.

## Timer and extractor

In a **timer**, the core is graph propagation (arrival, required and slack) and the delay-model API, together with the correctness of path enumeration. The change surface is which corners and paths to report, pruning, fast or surrogate modes and the derating scheme. In an **extractor**, the core is the field-solver mathematics and netlist correctness: every conductor is accounted for and capacitances are physically consistent. The change surface is the accuracy-versus-speed mode, which nets and layers to extract at which fidelity, the pattern library used for rule-based extraction, and calibration against a reference solver. A request to report the 50 worst paths per group, or to use a coarser pattern library on non-critical nets, is a change-surface request.

## Requests that seem to cross the line

Some requests sound like core changes but are not. "Route clock nets with double width" is a layer directive or a non-default rule, not a change to legality. "Ignore timing on test-mode paths" is a constraint or reporting choice, not a change to propagation. Others sound like tuning but are core: "let the placer leave a few overlaps for the router to fix" breaks legalization; "skip coupling capacitance below 0.01 fF to go faster" changes what the extractor guarantees and must be handled as an explicit, reviewed accuracy mode rather than a silent threshold. The question to ask is whether the result would still meet the engine's correctness contract.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapters on legalization and detailed routing; the OpenROAD documentation for gpl, dpl and drt.
