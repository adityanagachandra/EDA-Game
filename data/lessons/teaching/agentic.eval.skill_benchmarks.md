# Benchmarks for agentic EDA

Choosing a model, tuning a harness or claiming progress all need a benchmark. Generic coding benchmarks are the wrong instrument for agent-first EDA because they test generating fresh code on small, self-contained problems. The task here is editing large existing engines correctly under a tier policy, and improving QoR on real designs. A good benchmark measures that, honestly.

## Why generic benchmarks under-separate models

Popular coding benchmarks pose short, self-contained problems with unit tests: write a function that sorts intervals. Strong models now solve most of them, so scores cluster near the top and differences shrink. They also measure the wrong skill. Editing a placer's move evaluator requires reading an existing interface, respecting invariants that are not in the prompt, producing a minimal diff and keeping performance. Two models with similar scores on fresh-code benchmarks can differ widely in compile-pass rate and skill-edit correctness on a real C++ engine.

## Building a benchmark from real engines

A skill-editing benchmark is a set of tasks, each a plain-language change request against a real engine (an open-source placer, router or timer such as those in OpenROAD and OpenSTA works well), with its tier label, the files in scope, golden tests, a differential specification of what should and should not move, and hidden tests the model never sees. Tasks span tiers and engines: configuration changes, new plug-ins at extension points, and a few legitimate core changes. Metrics are compile-pass rate, skill-edit correctness through the funnel, scope violations, attempts and wall time per success in the harness, and QoR deltas against baselines where the task is an optimization.

## Reporting results honestly

Benchmarks are easy to fool, including by accident. Tasks drawn from public repositories may appear in a model's training data, so held-out tasks written after the model's training cutoff, or on private engines, are needed. Results vary with sampling, so report the number of samples, pass@k with the unbiased estimator, and confidence intervals, averaged per task. QoR claims need baselines, seeds and the designs named. A claim such as "up to 30% better" with an unnamed baseline and no workload is not a benchmark result; it is a best case. The same discipline applies to evaluating ML-for-EDA models in general.

## Going deeper

Jimenez et al., "SWE-bench: Can Language Models Resolve Real-World GitHub Issues?" (2024), on benchmarks built from real repositories; the OpenROAD project's open designs and flow scripts.
