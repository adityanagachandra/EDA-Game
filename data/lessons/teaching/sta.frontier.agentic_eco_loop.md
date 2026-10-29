# Agentic timing ECO loops

An agentic ECO loop reads timing reports, proposes fixes, applies them through the engines, verifies the result and iterates, with a human approving the final change. The engines stay deterministic and authoritative; the agent supplies search and judgment between them. Designing the loop means deciding what it may read, what it may propose, what it must verify, and where its authority ends.

## What the loop reads and proposes

The loop reads timing reports for every active scenario, merged slacks, the failing paths with their increments, and the physical context: placement density, congestion, whether a cell can grow in place. Its proposal vocabulary is netlist and clock changes: gate sizing, Vt swaps, buffer insertion or removal, and useful-skew adjustments. It does not propose constraint changes. Adding a false path or relaxing a clock makes the report pass without making the chip work, so such edits fall outside the ECO agent's authority and need separate human justification. Each proposal is small and local, so the timer can update incrementally and the next proposal sees the real state.

## What the loop verifies

After each change the loop re-times incrementally and accepts the change only if no check degrades: setup and hold together, in every scenario, along with design rules such as max transition and limits on area and leakage. Judging a move on setup alone causes oscillation, where a fix breaks hold and the hold fix breaks setup. If a fast surrogate ranks candidates, the chosen fix is confirmed with the signoff timer before acceptance, because the surrogate is a proxy with error. A differential timing report against the pre-ECO baseline shows every endpoint that moved, so side effects are visible. The loop's record of proposals, results and rejections becomes the evidence for review.

## Authority, containment and reuse

A human approves the final ECO with the loop's evidence attached; the loop works at engine speed but does not tape out on its own. It runs in a sandbox with resource limits and no network, because it executes generated actions on proprietary design data that must stay on site. When the agent finds a tactic that works, for example a sequence of Vt swaps and buffer moves for a class of hold problems, that know-how can be validated and promoted into a reusable skill, while the design itself stays inside the IP boundary. If an ECO passes in the implementation timer but fails at signoff, the lesson is to close on the signoff engine or keep margin for the correlation gap.
