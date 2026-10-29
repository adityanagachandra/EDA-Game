# Compile-pass rate and pass@k

The first number anyone measures for a code-editing agent is how often its edits compile, and the second is how often at least one of several attempts passes. Both have precise definitions and a classic estimation trap. A tool designer evaluating models for a skill-editing harness should compute them correctly and know what they do not say.

## Compile-pass rate

For a benchmark of edit tasks, the compile-pass rate is the fraction of model edits that build successfully: if a model produces 200 edits and 154 build, the rate is 77%. It is a necessary condition, cheap to measure and a good early filter: a model that cannot produce compiling C++ against an existing codebase cannot drive a harness. But it rewards trivial edits as much as correct ones; an edit that compiles because it changed nothing passes. It is reported per task family, since compile rates for configuration edits and for new plug-ins differ widely.

## The unbiased pass@k estimator

pass@k is the probability that at least one of $k$ attempts at a task succeeds. Estimating it from exactly $k$ attempts per task is noisy, so one draws $n \ge k$ attempts, counts the $c$ successes, and uses the unbiased estimator
$$\widehat{\text{pass@}k} = 1 - \frac{\binom{n-c}{k}}{\binom{n}{k}},$$
the chance that a random $k$-subset of the $n$ attempts contains at least one success. With $n = 10$ and $c = 3$, pass@3 is $1 - \binom{7}{3}/\binom{10}{3} = 1 - 35/120 = 0.7083$. Plugging $\hat p = 0.3$ into $1 - (1 - \hat p)^k$ gives 0.657, a biased estimate.

## Averaging over tasks

pass@k is computed per task and then averaged; pooling attempts first gives a different and misleading number. Take three tasks with per-attempt success rates 0.9, 0.5 and 0.1. The pooled rate is 0.5, and $1 - 0.5^3 = 0.875$. Per task, pass@3 is 0.999, 0.875 and 0.271, with mean 0.715. Because $1 - (1-p)^k$ is concave, pooling overstates pass@k whenever tasks differ in difficulty. And pass@k with independent samples is not the same as success within $k$ harness attempts, which use feedback between attempts and usually do better.

## Going deeper

Chen et al., "Evaluating Large Language Models Trained on Code" (2021), which introduced the unbiased pass@k estimator.
