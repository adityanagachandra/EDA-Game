# Agentic IR-drop ECO loops

IR and EM closure is a natural target for agentic signoff loops: an agent reads the drop map, proposes strap, via or decap changes, runs them through the PDN and analysis engines, and iterates until the design meets its budgets, with a human approving the final ECO. The engineering questions are concrete: what the agent reads, which fixes it may propose, how each fix is verified, and which metrics the loop must watch so that fixing drop does not break something else.

## What the agent reads and proposes

The input is the signoff report: per-node or per-instance drop, the worst regions, the layers and branches that carry the current, and the instances whose current causes the hot spot. The agent classifies the problem. A static drop is a resistance problem: widen or add straps, add via stacks between rails and straps, or add supply connections. A dynamic droop with small static drop is a local charge problem: insert decap next to the switching cluster or spread high-current cells apart. The agent expresses each fix through the PDN and placement engines rather than by editing layout text, so the engines enforce legality and the database stays consistent.

## Choosing among candidates

A fast solver predicts each candidate's effect. If the worst drop is 70 mV against a 40 mV budget, the gap is 30 mV. Given candidates with predicted reductions of 25, 32, 38 and 30 mV and costs of 2, 5, 7 and 6 (tracks and area combined), the cheapest is insufficient, and the cheapest sufficient one is the 32 mV fix at cost 5. Predicted reductions are local sensitivities: two fixes that relieve the same current path overlap, so their combined effect is less than the sum. That is one reason every accepted candidate is re-verified by the accurate solver before the loop moves on.

## Verifying and recording the ECO

After each edit the loop re-runs IR analysis and also the checks the edit can disturb: EM on the paths that now carry more current, DRC on the changed area, routing congestion and timing. A strap that fixes drop can steal tracks from congested signals; decap in every free site leaks and blocks later timing buffers. Each ECO is logged with its trigger, the edit, before and after metrics and the engine versions, so it is reviewable and reversible. The fast estimate steers; the signoff solver decides, and a human approves the final change into the signed-off database.

## Going deeper

Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the chapter on power and ground routing.
