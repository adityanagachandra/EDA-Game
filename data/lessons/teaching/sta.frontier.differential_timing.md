# Differential timing reports as golden tests

How do you prove that a change to a timer did what it claimed and nothing else? Not by checking that the worst slack improved, and not by checking that the code compiles. A differential timing report compares every endpoint's slack on a reference design before and after the change, against an explicit expectation. It is the golden test that gates every derate provider, delay-calculator tweak or agent-proposed edit.

## Which endpoints moved

Run the baseline and the changed timer on the same reference design, constraints and parasitics, and list each endpoint's slack in both. An endpoint moved if its slack changed by more than a tolerance $\varepsilon$, chosen above numerical noise, for example 0.5 ps. If endpoints change as $67 \to 63$, $78 \to 78$, $32 \to 30$, $47 \to 40$, $67 \to 67$, $55 \to 52$, $46 \to 42$ and $18 \to 18$ ps, five of eight moved. Determinism matters here: if two runs of the unchanged timer differ, the report is meaningless, which is why seams must be deterministic and seedable. Large reference designs make the test sensitive; small targeted designs make failures easy to diagnose.

## Unexpected movers

The change comes with an expectation: for a derate provider, the set of endpoints whose paths it should affect and by how much. Every endpoint outside that set that moved is an unexpected mover, and one is enough to fail the test. In a provider that should touch only paths through HVT data cells, a moved endpoint whose path has no HVT cell points to a bug, such as the rule firing on clock arcs, or a side effect through shared logic that the expectation missed. The report should show each unexpected mover with its path, so the engineer or agent can tell a real bug from an incomplete expectation.

## Intended endpoints and the promotion rule

Endpoints in the intended set must move by the intended amount. If the provider should lower their slack by exactly 3 ps, an intended endpoint that moved by 4 ps is also a violation. The total failure count is the number of unexpected movers plus the number of intended endpoints off by more than the tolerance; with no unexpected movers and two intended endpoints at -4 instead of -3 ps, the count is 2. The change is promoted only when the count is zero. The same check protects agentic flows: an agent's edit to an engine passes the harness only when its differential report matches the claim it made, which turns a plausible change into a verified one.
