# Agentic package co-design loops

An agentic package co-design loop reads the reports of package, board, thermal and power-integrity tools, proposes changes such as ball reassignments, decap placement, die order or TIM choice, runs them through the engines and iterates until signoff passes, with a human approving the final change. It is the package-level counterpart of an agentic timing ECO loop, and it needs the same discipline about what it may change and how it proves its changes.

## What the loop reads and proposes

The loop reads structured reports: impedance profiles against the target, IR drop and current per ball and bump, escape-routing completion, per-die temperatures, warpage and joint strain. It maps each violation to candidate fixes from a library: add power balls under a hot region, move decoupling capacitors closer, swap die order in the stack, change the TIM, reassign signals to outer rows. Each proposal is a small, typed edit to the design database, such as a ball-map change or a stackup parameter, not free-form text. Typed edits can be validated before any engine runs, and they make the history of the loop reviewable.

## Seams the loop needs

The loop works only if engines expose seams for its changes. A ball-assignment interface lets it move a net between balls and re-run escape routing incrementally. A material and stackup interface lets it swap a TIM or substrate core without editing the solver. A power-map interface lets the chip team's power model update the thermal run. Through these, most fixes are Tier 1 configuration or Tier 2 plug-ins. A fix that requires editing an engine's core, say a new way to compute plane inductance, is out of scope for the loop and goes to engine developers. Each engine should also support incremental runs, because the loop's speed depends on not redoing the whole analysis after each small edit.

## Convergence and approval

Co-design fixes interact: moving power balls helps IR drop but may block escape routing, and a thicker TIM bond line from warpage undoes a thermal fix. The loop must therefore re-run every affected check after each change, not only the one it targeted, and track a global violation count so it can detect oscillation between two fixes. It should stop when signoff passes, when no proposal improves the count, or when a budget of iterations is spent. The final change goes to a human as a diff of the ball map, stackup and material choices with before-and-after reports. The human approves; the loop never signs off its own work.

## Going deeper

Swaminathan and Engin, "Power Integrity Modeling and Design for Semiconductors and Systems" (2007); the agentic domain concept on signoff and ECO loops.
