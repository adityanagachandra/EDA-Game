# Golden tests for statistical engines

Every engine change in an agent-first flow is promoted only if golden and differential tests pass. For a deterministic engine, a golden test compares outputs exactly. For a statistical engine, such as a Monte Carlo yield estimator, an importance sampler or a statistical timer, the output is random, and the test itself becomes a hypothesis test. Designing it means choosing tolerances from the estimator's standard error and controlling the false-alarm rate of the whole suite.

## Tolerance bands from the standard error

A golden case with a known true value $p$ is run with $N$ samples and a fresh or fixed seed. A correct estimator lands within $z$ standard errors of the truth with probability about $1 - 2Q(z)$, where $SE = \sqrt{p(1-p)/N}$. With $p = 0.02$, $N = 10{,}000$ and $z = 4$, $SE = 0.0014$ and the band is $\pm 0.56$ percentage points. A band set by intuition, such as "within 1 percent", is either too loose to catch a real bias or so tight that correct changes fail. Bias detection power follows from the same numbers: a bias smaller than about one SE is essentially invisible at this $N$, so the golden case must use enough samples for the bias it must catch.

## The suite-wide false-alarm rate

A regression suite holds many statistical tests. If each has false-alarm probability $\alpha$ when the engine is correct and they are independent, the chance a correct change fails at least one is $1 - (1 - \alpha)^m$. With $m = 100$ tests at $\alpha = 0.01$ that is 63 percent: most correct changes would be rejected. A flaky suite is worse than a slow one, because the harness and the people reading it learn to rerun until green, and then real regressions pass. The suite's false-alarm rate is therefore a designed quantity, not an accident of each test's tolerance.

## Setting per-test thresholds

To hold the family-wise rate at $F$ across $m$ tests, the Bonferroni split gives each test $\alpha = F/m$, valid even if the tests are dependent. For a two-sided Gaussian test statistic the threshold is $z = Q^{-1}(\alpha/2)$. With $m = 200$ and $F = 0.01$, $\alpha = 5 \times 10^{-5}$ and $z = Q^{-1}(2.5 \times 10^{-5}) = 4.06$. Wider bands cost power, so suites compensate with more samples per golden case and with fixed seeds for the routine run, which make the run deterministic, plus a separate periodic coverage test across many seeds that checks the statistical claim itself. The harness reports the test statistic, not just pass or fail, so a drift toward the threshold is visible before it fails.

## Going deeper

Casella and Berger, "Statistical Inference" (2nd ed., 2002), chapter 8; the golden and differential testing material in this curriculum's agentic domain.
