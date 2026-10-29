# PDK-aware agents and change tiers

An agent that edits or drives physical-design engines must know the process: rules, layers, models and their versions. A PDK-aware agent reads those facts from the installed PDK through tools instead of recalling them, and it routes each process-related change to the smallest safe tier of the engine. This concept applies the four-layer engine model and change tiers to the interface between EDA engines and process data.

## Reading process facts instead of recalling them

Process values are version-specific and confidential, so a language model's memory of "typical" values is worthless and dangerous. A PDK-aware agent answers process questions by calling tools that parse the installed technology files, rule decks and model cards and return typed results: the M2 minimum spacing, the via enclosure table, the corners available. Skill manifests state which processes and PDK versions each skill supports, so selection is a cheap text match. Before acting, the agent checks provenance: if a design was built on version 1.2 and the installed PDK is 1.4, it reports the mismatch and diffs the rule and model changes first.

## Which tier a process change lands in

Most process-driven changes should be data. A new value for an existing density limit, a new layer in the stack or a new corner is Tier 1: the engines read it from updated PDK files, and if that ever needs a recompile, the seam is in the wrong place. A new kind of rule the engine cannot express, such as a spacing that depends on two neighbours' colours and run length, is Tier 2: a new rule evaluator plugged into the router's or DRC engine's rule-check extension point, leaving the path model and legality core untouched. A change to grid legality or the core geometry kernel is Tier 3: rare, reviewed by a human as a diff, and run through full regression.

## Guarding against invented and gamed process data

Two failure modes need structural guards. First, an agent may state a process value that is not in the deck; the harness should require every process number in an output to cite the file and rule it came from, and a checker verifies each citation. Second, an optimizing agent may reach a clean DRC by weakening the oracle, for example by disabling a check; rule decks and signoff settings must be read-only to optimization agents. When a rule is ambiguous, the right move is to run both readings as checks on test clips, report the difference, and escalate to a human, since only the foundry can resolve it.

## Going deeper

The OpenROAD and KLayout documentation on technology files and DRC decks, as examples of the files a PDK-aware agent must parse; the agent-first EDA anchors on change tiers and skill manifests.
