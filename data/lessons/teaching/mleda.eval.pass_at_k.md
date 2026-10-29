# Pass at k from n samples

Code-generation benchmarks report pass@k: the probability that at least one of $k$ sampled solutions to a problem passes its tests, averaged over problems. Estimating it from $k$ samples alone is noisy, so evaluations draw $n \ge k$ samples, count the $c$ that pass, and use an unbiased estimator. Reading pass@k correctly means knowing what the $k$ tries assume.

## Pass at 1

With $n$ samples per problem and $c$ passing, pass@1 for that problem is $c/n$, the probability that one random sample passes. The benchmark's pass@1 is the average over problems. Five problems with 100 samples each and 83, 14, 35, 65 and 92 passing give pass@1 = 57.8 percent. Decoding settings matter: greedy decoding gives one deterministic sample per problem, while sampled pass@1 at temperature 0.8 is an expectation over a distribution. Both are legitimate and they are not interchangeable, so a reported number must state the decoding. A per-problem fraction near zero or one is estimated well; one near one half needs the most samples.

## The unbiased estimator

For one problem, pass@k is the probability that a random subset of $k$ of the $n$ samples contains at least one passing sample: $1 - \binom{n-c}{k}/\binom{n}{k}$. With $n = 20$, $c = 2$ and $k = 5$, $\binom{18}{5} = 8568$ and $\binom{20}{5} = 15{,}504$, so pass@5 = $1 - 8568/15{,}504 = 0.447$. The plug-in estimate $1 - (1 - c/n)^k = 0.410$ is biased low for this purpose. When $n - c < k$ every subset contains a pass and the estimate is 1. Compute the ratio as a running product to avoid overflow for large $n$.

## Averaging and interpreting

The benchmark's pass@k averages the per-problem values. With $n = 20$, $k = 5$ and 0, 2 and 10 passing samples on three problems, the per-problem values are 0, 0.447 and 0.984, and pass@5 is 47.7 percent, against a pass@1 of 20 percent. Pass@k with $k > 1$ assumes an oracle that picks the passing sample. In a tool that oracle must be a real check, a testbench or an equivalence proof, and it costs a tool run per sample. A large gap between pass@k and pass@1 means a checker is essential, not that the model is reliable.

## Going deeper

Chen et al., "Evaluating Large Language Models Trained on Code" (2021), section 2.1, for the unbiased estimator.
