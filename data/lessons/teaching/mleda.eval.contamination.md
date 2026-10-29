# Benchmark contamination

A benchmark score estimates ability on new problems only if the model has not seen the benchmark. Language models trained on web-scale code may have seen public RTL exercises and their posted solutions, and learned EDA predictors can see their test designs through variants. Contamination turns a test of generalization into a test of recall, and detecting it takes behavioral evidence, not assurances.

## How contamination happens

Contamination means benchmark problems, their solutions or close variants are in the training data. RTL benchmarks built from public tutorial exercises are exposed: the exercises and many solutions are posted online, where large crawls collect them. A model that scores far higher on a public set than on fresh problems of the same kind is showing that signature. Learned EDA predictors have the same failure in another form: if the test design also appears in training under other flow settings, the model recognizes it instead of generalizing. Excluding the official benchmark files does not exclude reposted copies.

## Detecting memorization

Behavioral tests are the evidence. Give the model a prefix of a reference solution and see whether it completes it verbatim. Rewrite problems in new words and rename signals while keeping the hardware behavior; a sharp drop in score on equivalent problems points to memorization. Text-overlap checks between training data and solutions catch exact copies but miss renamed variants and flag generic boilerplate, so they are a screen, not a verdict. Distillation adds a path: fine-tuning on synthetic Verilog from a larger model can carry that model's memorized benchmark solutions into a dataset that looks clean.

## Contamination-resistant practice

The durable defense is a private, regularly refreshed problem set with hidden references, evaluated through a submission process, so problems are never published. Repeated tuning against a benchmark is contamination by another route: adjusting prompts or agent settings while watching the test score turns the test set into a training set, even without gradient updates. A company evaluating an assistant for its own designs should use internal tasks that never left the company, which are uncontaminated by construction and represent the real distribution of work. Reporting the date each problem was written lets readers compare it with a model's training cutoff.

## Going deeper

Liu et al., "VerilogEval" (ICCAD 2023), discussion of problem sources; Chen et al., "Evaluating Large Language Models Trained on Code" (2021), on hand-written held-out problems.
