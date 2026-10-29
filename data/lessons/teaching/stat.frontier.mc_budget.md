# Allocating a Monte Carlo budget

An agent that runs variation analysis has a fixed simulation budget and many places to spend it: several corners, several process splits, several candidate designs. Classical sampling theory already says how to split a budget to minimize estimator variance. Encoding that in the agent's policy, rather than splitting evenly or by intuition, is one of the clearest ways a statistical engine makes an agent loop cheaper.

## Neyman allocation

Suppose the quantity of interest is a population mean over strata $h$ (process splits, operating modes) with weights $W_h$ and within-stratum standard deviations $\sigma_h$. The stratified estimator $\sum_h W_h\bar{y}_h$ has variance $\sum_h W_h^2\sigma_h^2/n_h$. Minimizing it subject to $\sum n_h = n$ gives Neyman allocation:
$$n_h = n\,\frac{W_h\sigma_h}{\sum_j W_j\sigma_j}.$$
With weights 0.5, 0.3, 0.2, sigmas 2, 4, 8 nA and $n = 1000$, the products are 1, 1.2 and 1.6, so the allocation is 263, 316 and 421 samples. The smallest stratum gets the most samples because it is the noisiest. Equal allocation would waste samples on the quiet stratum.

## Unequal simulation costs

Simulation costs differ: a split simulated at a slow corner with a long transient may cost several times more per sample. With per-sample costs $c_h$ and a total budget $C$, minimizing the variance subject to $\sum c_h n_h = C$ gives
$$n_h = C\,\frac{W_h\sigma_h/\sqrt{c_h}}{\sum_j W_j\sigma_j\sqrt{c_j}}.$$
Expensive strata get fewer samples, by the square root of their cost. The policy needs $\sigma_h$ before it has run, so the agent runs a small pilot in each stratum, estimates $\sigma_h$, allocates the remainder, and re-allocates as estimates firm up. This is a design decision for the harness: the allocation rule is a policy-layer skill, the per-stratum sigma and cost are observations from hooks, and the total budget is configuration.

## Sequential stopping

For a single failure rate, the agent should stop when the answer is precise enough, not after a fixed count. After a pilot of $n_0$ samples with $\hat{p}$ failures per sample, a target standard error $SE$ requires about $N = \hat{p}(1 - \hat{p})/SE^2$ samples in total, so $N - n_0$ more. A pilot of 500 samples with 15 failures ($\hat{p} = 0.03$) and a target $SE = 0.003$ needs $0.0291/0.000009 = 3234$ in total, 2734 more. Stopping rules must be honest: stopping as soon as an estimate crosses a spec, and continuing otherwise, biases the result. The policy fixes the precision target in advance, and the reported interval accounts for the sequential procedure.

## Going deeper

Cochran, "Sampling Techniques" (3rd ed., 1977), chapter 5; Owen, "Monte Carlo Theory, Methods and Examples," the chapter on variance reduction.
