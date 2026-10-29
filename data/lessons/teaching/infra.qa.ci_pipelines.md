# Continuous integration for EDA tools

Continuous integration builds the tool and runs its tests on every proposed change, before it merges. For EDA tools the tests are expensive: a regression design takes minutes to hours, and the meaningful checks are QoR comparisons against baselines, not just pass or fail. A CI pipeline therefore has to be engineered like a parallel job: sharded across runners to meet a time budget, tiered so cheap checks run first, and robust against flaky tests that fail without a real bug.

## Wall time of a sharded test stage

A test stage with $N$ designs of $t$ minutes each, spread over $R$ runners, gives each runner at most $\lceil N/R \rceil$ designs. Each runner usually first builds or downloads the tool, a fixed setup time $S$. The stage wall time is
$$W = S + \left\lceil \frac{N}{R} \right\rceil t.$$
Example: 30 designs of 4 min on 6 runners with 8 min of setup: $8 + 5 \times 4 = 28$ min. On 12 runners: $8 + 3 \times 4 = 20$ min. Setup is paid on every runner in parallel, so it bounds the wall time from below and becomes the main cost as runners are added. Building the tool once and sharing the binary as a CI artifact turns per-runner setup into a download, and real suites balance shards by measured design run times rather than by count.

## Runners needed to meet a budget

Given a budget $B$ minutes, the stage fits if $S + \lceil N/R \rceil t \le B$, that is, if each runner gets at most $K = \lfloor (B - S)/t \rfloor$ designs. The fewest runners is then $R = \lceil N/K \rceil$.

Example: a pull-request check must finish in 30 min, with 45 designs of 5 min and 10 min of setup. $K = \lfloor 20/5 \rfloor = 4$ designs per runner, so $R = \lceil 45/4 \rceil = 12$ runners, which gives $10 + 4 \times 5 = 30$ min. Eleven runners would need 5 designs on some runner, 35 min. When no affordable runner count fits, the suite is tiered: a small, fast set of designs gates every pull request, and the full suite with QoR comparison runs nightly against the main branch.

## Flaky tests and retries

A flaky test fails sometimes on correct code: a timeout on a loaded runner, a race in the test harness, a nondeterministic engine. If $f$ flaky tests each fail independently with probability $q$, a correct commit passes all of them with probability $(1 - q)^f$. With $f = 5$ and $q = 0.05$, that is $0.95^5 \approx 77.38\%$: almost one correct commit in four is blocked.

Rerunning a failed flaky test once, and accepting a pass, cuts its failure probability to $q^2 = 0.0025$, so all five pass with probability $(1 - 0.0025)^5 \approx 98.76\%$. Retries are a stopgap with a cost: they hide nondeterminism that may be a real engine bug, and a test that passes on retry should be logged and fixed. For EDA engines, a QoR test that is flaky usually means the engine is not deterministic, which is exactly the property CI exists to protect.

## Going deeper

Humble and Farley, "Continuous Delivery" (2010); Luo et al., "An Empirical Analysis of Flaky Tests" (FSE 2014); the OpenROAD and OpenROAD-flow-scripts repositories' CI configuration as working examples.
