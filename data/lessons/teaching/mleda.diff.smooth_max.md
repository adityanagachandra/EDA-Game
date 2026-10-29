# Smooth maximum for differentiable timing

Static timing takes a max at every gate input: the output arrives after the latest input. The max has a gradient only through the latest input, which stalls gradient-based optimization. Replacing it by the log-sum-exp soft maximum makes arrival times differentiable in every delay, at a controlled cost in accuracy. The temperature $\gamma$ sets the trade.

## Log-sum-exp

The soft maximum is $\mathrm{LSE}_\gamma(a) = \gamma \ln \sum_i e^{a_i/\gamma}$. Compute it stably by factoring out the largest value: $\mathrm{LSE}_\gamma(a) = a_{\max} + \gamma \ln \sum_i e^{(a_i - a_{\max})/\gamma}$, which avoids overflow for arrival times of hundreds of picoseconds. For arrivals 50 and 48 ps and $\gamma = 2$ ps, $\mathrm{LSE} = 50 + 2\ln(1 + e^{-1}) = 50.63$ ps. As $\gamma \to 0$ it approaches the max; as $\gamma$ grows it approaches a scaled sum. Inputs far below the max contribute almost nothing, so LSE behaves like max where one input clearly dominates.

## How much it overestimates

LSE always lies between the true max and the max plus $\gamma \ln n$ for $n$ inputs. In the example the gap is 0.63 ps, within the bound $2 \ln 2 = 1.39$ ps. The gap is largest when the inputs tie and vanishes when one dominates. Because LSE overestimates, a smoothed timer is pessimistic: it reports slightly later arrivals than STA, and the error accumulates along a path, one max per gate. A differentiable timer therefore uses a small $\gamma$, reports true STA numbers for evaluation, and uses the smooth values only for the gradients that steer optimization.

## Gradients through the soft max

The derivative of LSE with respect to input $i$ is the softmax weight $e^{a_i/\gamma}/\sum_j e^{a_j/\gamma}$. For arrivals 50 and 48 ps and $\gamma = 2$, the later input gets weight $1/(1 + e^{-1}) = 0.731$ and the other 0.269. The weights sum to one, so the gradient of an arrival is shared among inputs by how close they are to critical. With a hard max the latest input would get all of it and a near-critical input none, so an optimizer would fix one path and expose the next. The soft gradient lets it improve the near-critical paths together.

## Going deeper

Guo and Lin, "Differentiable-Timing-Driven Global Placement" (DAC 2022); Boyd and Vandenberghe, "Convex Optimization" (2004), section 3.1.5, on log-sum-exp.
