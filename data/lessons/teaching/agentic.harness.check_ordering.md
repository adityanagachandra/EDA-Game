# Ordering checks by cost and failure rate

A harness runs several checks on every candidate change: build, unit tests, golden tests, a differential QoR run, perhaps a sanitizer build. It stops at the first failure. The order of the checks changes the expected cost of evaluating a candidate, sometimes by a large factor, and the best order follows from a short calculation that every harness designer should know.

## Expected cost of a staged pipeline

Let stage $i$ cost $c_i$ and pass with probability $q_i$, independently. The pipeline stops at the first failure, so stage 2 runs only if stage 1 passed, stage 3 only if both passed, and so on:
$$E[C] = c_1 + q_1 c_2 + q_1 q_2 c_3 + \cdots$$
Take a build of 30 s that passes 70% of the time, golden tests of 120 s that pass 60% of the time and a differential QoR run of 600 s. In that order, $E[C] = 30 + 0.7 \times 120 + 0.7 \times 0.6 \times 600 = 30 + 84 + 252 = 366$ s. If the differential run were first, every candidate would pay its full 600 s.

## The best order

For independent stages, the expected cost is minimized by running stages in increasing order of $c_i / (1 - q_i)$: cost per unit chance of failing. A cheap check that often fails should go first; an expensive check that rarely fails should go last. An exchange argument proves it: swapping adjacent stages $a$ then $b$ to $b$ then $a$ changes the cost by $c_b + q_b c_a - c_a - q_a c_b$, which is negative exactly when $c_b / (1 - q_b) < c_a / (1 - q_a)$. Some orders are forced: tests cannot run before the build, so the rule applies among checks that can run in either order.

## How much reordering saves

Suppose a lint check costs 5 s and fails 30% of the time, a sanitizer build costs 200 s and fails 10% of the time, and unit tests cost 40 s and fail 25% of the time. The ratios are 16.7, 2000 and 160, so the order is lint, unit tests, sanitizer: $5 + 0.7 \times 40 + 0.7 \times 0.75 \times 200 = 138$ s. Running the sanitizer first gives $200 + 0.9 \times 5 + 0.9 \times 0.7 \times 40 = 229.7$ s. The saving compounds over thousands of candidates. If stages share work, such as one build feeding several tests, or failures are correlated, the rule is a starting point and measured pass rates should decide.

## Going deeper

Smith, "Various optimizers for single-stage production" (1956), the origin of ratio rules for sequencing; any operations research text on sequential testing.
