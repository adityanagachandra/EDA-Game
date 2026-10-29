# Coverage closure

Coverage closure is the work of driving a coverage model to completion, or to an agreed, justified stopping point. It is where verification effort concentrates late in a project: analyzing holes, deciding which are reachable, changing stimulus, excluding what is impossible, and keeping regressions short enough to rerun. It is an optimization loop around the simulator, which makes it a natural target for automation and for agents.

## Holes, exclusions and plateaus

A **hole** is a bin or cross that no test has hit. Coverage is measured over the merged databases of all runs, so a hole means the whole regression missed it. Each hole needs a decision. If the scenario is reachable and matters, change the stimulus. If it is unreachable in this configuration, exclude it with a documented, reviewed waiver, so the denominator stays honest. If it is impossible by design, an illegal bin turns any hit into an error. A **plateau**, where coverage stops rising as seeds are added, signals that the remaining holes are improbable under the current constraints. Fast 100% coverage on a new block is also a warning: the model is probably too coarse.

## Ranking tests with greedy set cover

After thousands of runs, most tests add nothing new. **Test ranking** finds a small subset that preserves the merged coverage, so that it can be rerun on every change. Picking the minimum subset is the set cover problem, which is NP-hard, so tools use the greedy algorithm: pick the test that covers the most bins, then repeatedly pick the test that adds the most not-yet-covered bins, until no test adds anything. Example: T1 covers {1, 2, 3, 4}, T2 covers {3, 4, 5}, T3 covers {5, 6}, T4 covers {1, 6}. Greedy picks T1 (4 new), then T3 (2 new: 5 and 6); T2 and T4 now add nothing. Two tests keep all six bins. Greedy is guaranteed within a factor $H_n$ of the optimum and is fast.

## The coverage-driven loop

Coverage-driven verification is a loop: run the regression, merge coverage, analyze holes, change stimulus, repeat. The change depends on why a hole is hard. A hole needing a full FIFO during a flush calls for a sequence that throttles the reader, then randomizes around that state. A hole needing a rare field value calls for `dist` weights. A hole needing a long specific sequence may need a directed test. Each change is verified the same way: the hole closes, nothing else regresses, and checking stays intact. Sign-off then rests on a coverage model reviewed against the specification, with every hole hit or waived, combined with code coverage and assertion results rather than any single number.

## Going deeper

Piziali, "Functional Verification Coverage Measurement and Analysis" (2004); Wile, Goss and Roesner, "Comprehensive Functional Verification" (2005), chapter on coverage.
