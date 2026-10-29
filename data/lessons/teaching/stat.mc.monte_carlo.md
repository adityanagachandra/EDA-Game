# Monte Carlo estimation

Monte Carlo is the reference method of variation analysis: draw process parameters from their distributions, simulate the circuit for each draw, and count. It needs no linearity, no Gaussian assumption and no model of the circuit beyond the simulator itself, which is why every faster method is validated against it. Its weakness is cost, and a tool designer must be able to state that cost exactly from the standard error.

## The estimator and its standard error

To estimate a yield $Y = P(\text{pass})$, draw $N$ independent samples, simulate each, and set $\hat{Y} = (\text{passes})/N$. Each sample is a Bernoulli trial, so $\hat{Y}$ is unbiased with variance $Y(1-Y)/N$, estimated by $SE = \sqrt{\hat{Y}(1 - \hat{Y})/N}$. With 1000 samples and 50 failures, $\hat{Y} = 0.95$ and $SE = \sqrt{0.95 \times 0.05/1000} = 0.0069$, so the yield is known to about 0.7 percentage points. The estimator does not care how many process parameters there are: the error depends on $N$ and $Y$ only, which is why Monte Carlo remains usable for circuits with thousands of variation sources.

## Confidence intervals

By the central limit theorem, $\hat{Y}$ is approximately Gaussian for large $N$, so an approximate 95 percent interval is $\hat{Y} \pm 1.96\,SE$. In the example that is $95.0 \pm 1.35$ percent. The approximation breaks down when the expected number of failures is small, roughly when $N(1-\hat{Y}) < 10$: the interval can extend past 100 percent, or collapse to zero width when no failures are seen. For those cases an exact binomial (Clopper-Pearson) or Wilson interval is used instead. A Monte Carlo report without an interval is incomplete, because two runs that differ only in seed will disagree within it.

## Cost scales as one over the square root of N

For a continuous output such as a delay, the mean estimate has $SE = s/\sqrt{N}$, where $s$ is the sample standard deviation. Halving the error needs four times the samples; a tenfold improvement needs a hundredfold. A pilot run gives $s$, and the total budget follows as $N = (s/SE_{\text{target}})^2$: a delay with $s = 3$ ps measured to $SE = 0.1$ ps needs 900 simulations. This square-root law is the reason the rest of this subdomain exists. Quasi-Monte Carlo changes the rate, importance sampling changes the variance per sample, and fast surrogates change the cost per sample, but each is judged against plain Monte Carlo at the same accuracy.

## Going deeper

Owen, "Monte Carlo Theory, Methods and Examples," chapters 1 and 2; Agresti and Coull, "Approximate is better than exact for interval estimation of binomial proportions," The American Statistician, 1998.
