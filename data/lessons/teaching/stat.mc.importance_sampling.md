# Importance sampling

Importance sampling is the workhorse of rare-event estimation in circuit variability. Instead of sampling the true process distribution and waiting for a failure, it samples a distribution that fails often and corrects each sample with a weight, keeping the estimate unbiased. With a good proposal, a 6-sigma failure rate can be estimated with thousands of simulations instead of hundreds of billions.

## The likelihood-ratio weight

Let $f$ be the true density of the parameters $x$ and $I(x)$ the failure indicator, so $p = E_f[I(X)]$. For any proposal density $g$ that is positive wherever $f I$ is,
$$p = E_g\!\left[I(X)\,\frac{f(X)}{g(X)}\right], \qquad \hat{p} = \frac{1}{N}\sum_{i=1}^N I(x_i)\,w(x_i), \quad w = f/g,$$
with $x_i$ drawn from $g$. For a standard normal $f$ and the mean-shifted proposal $g = N(\mu, 1)$, the weight simplifies to
$$w(x) = \frac{\varphi(x)}{\varphi(x - \mu)} = e^{-\mu x + \mu^2/2}.$$
With $\mu = 4$ and $x = 4.2$, $w = e^{-16.8 + 8} = e^{-8.8} = 1.5 \times 10^{-4}$. Failures are now common, but each counts only for its true, tiny probability.

## Computing an estimate

The estimator sums the weights of failing samples and divides by the total number of samples drawn, failures and passes alike. To estimate $P(X > 4)$ with $\mu = 4$, draw from $N(4, 1)$; about half the samples exceed 4. Four samples at $3.6, 4.3, 4.9, 4.1$ give failing weights $e^{-9.2} = 1.01 \times 10^{-4}$, $e^{-11.6} = 9.2 \times 10^{-6}$ and $e^{-8.4} = 2.25 \times 10^{-4}$, summing to $3.35 \times 10^{-4}$; dividing by 4 gives $8.4 \times 10^{-5}$, against the exact $3.17 \times 10^{-5}$. With thousands of samples the estimate converges, and its standard error comes from the sample variance of $I(x_i)w(x_i)$. Forgetting the weights, or dividing by the number of failures, are the classic bugs.

## Several dimensions and the choice of shift

With $d$ independent standard normal parameters and a proposal shifted to mean vector $\mu$, the weight is $w(x) = \exp(-\mu \cdot x + \|\mu\|^2/2)$. The variance of $\hat{p}$ depends entirely on the proposal. Centering it on the most probable failure point, the closest failure to the origin, makes the weights of typical failures roughly uniform, and the variance falls by orders of magnitude. A poorly placed shift does the opposite: if the proposal undersamples part of the failure region, the estimate converges slowly to the right answer while reporting a deceptively small sample variance. In high dimensions, weights degenerate unless the shift lies only along the directions that matter, which is the motivation for locating the failure point first.

## Going deeper

Owen, "Monte Carlo Theory, Methods and Examples," chapter 9; Kanj, Joshi and Nassif, "Mixture importance sampling and its application to the analysis of SRAM designs in the presence of rare failure events," DAC 2006.
