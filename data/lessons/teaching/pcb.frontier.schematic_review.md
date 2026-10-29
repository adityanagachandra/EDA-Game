# Schematic-review agents

Design review of schematics is skilled, tedious work: checking every power pin's decoupling, every regulator's capacitors against its datasheet, every bus's pull-ups, every level crossing between voltage domains. A schematic-review agent combines deterministic checks with model reasoning over datasheets. It is useful only if its findings are grounded, precise and reviewable. That depends on the representation it reads, the evidence it attaches, and an evaluation on seeded defects.

## What the agent checks and from what

Split review into what rules can decide and what needs judgment. ERC-style predicates (undriven inputs, conflicting outputs, unpowered power pins, single-pin nets) are exact and cheap, so the agent calls the rule engine rather than reasoning about them. The model's value is in checks that need part knowledge: a regulator's output capacitor below its datasheet minimum, a missing pull-up on an open-drain bus, a 3.3 V output into a 1.8 V input, an enable pin left floating. The agent reads the **netlist with pin types and part parameters** through tool calls, never a picture of the schematic, so its queries are exact. Each finding carries evidence: the net or part, the rule or datasheet section relied on, and the observed values.

## Precision and recall

Evaluation uses reference schematics with **seeded defects**. If the agent raises $F$ findings, of which $T$ match seeded defects, its **precision** is $T/F$ and its **recall** is $T/S$ for $S$ seeded defects. Example: 25 seeded defects, 30 findings, 20 true: precision $20/30 \approx 66.7\%$, recall $20/25 = 80\%$. Both matter: low recall misses real bugs, and low precision buries them, since engineers stop reading reports that are mostly noise. Raising a confidence threshold usually trades recall for precision. The evaluation set must be held out from whatever data is used to tune the agent's prompts and skills, or measured recall inflates without real improvement.

## Grounding and the review loop

The characteristic failure is a confident citation of a requirement the datasheet does not contain. The harness catches it by **grounding**: each cited requirement is matched against the retrieved datasheet text before the finding is shown, and unmatched findings are dropped or flagged. Fixes reach the design only as proposed schematic edits in a reviewable diff that the designer approves, never applied silently. Findings confirmed at review or at bring-up become new seeded defects or new deterministic rules, the skill flywheel applied to review: a check that started as model reasoning, once validated, can be promoted into the rule engine where it runs exactly and cheaply on every design.

## Going deeper

The agent-first EDA notes in this curriculum (the harness loop and the skill flywheel). Manning, Raghavan and Schutze, "Introduction to Information Retrieval" (2008), chapter 8, on precision and recall.
