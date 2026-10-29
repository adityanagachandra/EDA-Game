# Benchmarks and evaluation for ML in EDA

Benchmarks make ML-for-EDA claims comparable, and each one measures a narrow thing. VerilogEval and RTLLM evaluate RTL generation; datasets such as CircuitNet serve prediction tasks; contest suites serve placement and routing. Knowing what a benchmark checks, how its scores are computed and what it leaves out is what lets a tool designer read a leaderboard without being misled.

## What VerilogEval and RTLLM measure

VerilogEval adapts exercises from an online Verilog practice collection into problems with natural-language descriptions, reference solutions and testbenches; a sample counts as correct when simulation agrees with the reference. Scores are pass@k. RTLLM targets larger design-level descriptions and evaluates syntax, functionality and design quality, the PPA of the synthesized result. Both are simulation-based, so a sample passes when it matches the reference on the tested behavior. Neither tests editing an existing codebase, working with a large set of interfaces, or meeting a production PPA target, which are the tasks an agentic RTL tool faces.

## Reading a score

pass@1 estimates the probability that one sample to a random benchmark problem passes; pass@10 assumes a perfect selector among ten. A model at 90 percent pass@10 and 40 percent pass@1 is useful only with a checker that picks among samples. Two papers that use different decoding settings report different quantities under the same name. Weak testbenches inflate scores with false passes. Benchmarks of modest size give noisy estimates, so scores need the sample count and a confidence interval: a few points of difference between models may not be significant. With 150 problems, a pass@1 near 50 percent has a standard error of about 4 points from problem sampling alone.

## Building an honest evaluation

Fresh problems written after a model's training cutoff, with hidden references, separate ability from memorization. Benchmarks for agentic tools should measure the correctness of edits inside a larger codebase, regressions in unchanged behavior, and PPA after synthesis. Vendor claims such as productivity gains "up to" some factor are claims against an unstated baseline until the task set, metric and baseline are known. For a company deciding whether to deploy an RTL assistant, the best benchmark is an internal set of representative tasks that never left the company, scored by its own verification flow.

## Going deeper

Liu et al., "VerilogEval: Evaluating Large Language Models for Verilog Code Generation" (ICCAD 2023); Lu et al., "RTLLM: An Open-Source Benchmark for Design RTL Generation with Large Language Model" (ASP-DAC 2024).
