# Placing a seam and designing its context

A seam is useful only if it sits where customer-specific behavior actually varies and passes everything that behavior depends on. Too high, and the plug-in must reimplement core work; too low, and it is called too often; too little context, and every new scheme needs a core edit. Designing the seam is the expert work engine owners contribute to an agent-first architecture.

## Where customer-specific behavior lives

Put the seam where teams disagree but the surrounding mathematics does not. In a timer, propagation is settled and margining varies, so the seam is the derate. In a router, legality is settled and cost trade-offs vary, so the seam is the edge cost and the rip-up policy. In a placer, legalization is settled and preferences such as grouping, abutment and region affinity vary, so seams sit at the detailed-placement move evaluator and the global placement net weighting. In an extractor, the solver is settled and accuracy trade-offs vary, so the seam selects per-net fidelity. For each customer behavior, ask what question the core would have to ask to accommodate it.

## What context the seam must expose

The context is the set of facts a plug-in may key on. A derate seam that passes cell, Vt class, early or late, rise or fall, logic depth, distance, slew, corner, clock or data and path group can host flat OCV, AOCV, POCV and most custom margin stacks. Remove logic depth and distance and AOCV cannot be written; remove clock or data and clock-specific margins cannot. A router cost seam must pass the net's class (clock, analog, power), the layer, the edge's current and historical congestion and the region, or a plug-in cannot penalize vias near an analog block. Missing context leads either to a core edit or to a plug-in that reaches into global state, which the rules forbid.

## Signs a seam is misplaced

The request log reveals seams in the wrong place. Repeated Tier 3 edits for one kind of request mean the seam lacks context or does not exist. Plug-ins that duplicate core logic, such as a derate provider that recomputes path depth itself, mean the seam sits too low or passes too little. A seam called in the innermost loop with tiny payloads, so that profiling shows boundary overhead, means it should be raised to a batch. Plug-ins that must coordinate through shared state mean two seams were meant to be one.

## Going deeper

The agent-first EDA notes in this curriculum, section on seams and the derating example; the OpenROAD detailed router documentation on its cost model.
