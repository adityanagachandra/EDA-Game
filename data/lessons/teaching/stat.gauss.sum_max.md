# Sums and maxima of Gaussian variables

Statistical static timing propagates arrival-time distributions through a timing graph with two operations: add a gate delay along an edge, and take the maximum where paths reconverge. The sum of Gaussians is Gaussian. The maximum is not, and the formulas that approximate it, due to Clark, are the computational heart of block-based statistical timing.

## Tightness probability

For jointly Gaussian arrival times $A \sim N(\mu_A, \sigma_A^2)$ and $B \sim N(\mu_B, \sigma_B^2)$ with correlation $\rho$, the difference $A - B$ is Gaussian with mean $\mu_A - \mu_B$ and standard deviation
$$\theta = \sqrt{\sigma_A^2 + \sigma_B^2 - 2\rho\sigma_A\sigma_B}.$$
The probability that $A$ is the larger, the **tightness probability**, is $T_A = P(A > B) = \Phi(\alpha)$ with $\alpha = (\mu_A - \mu_B)/\theta$. With $\mu_A = 110$, $\mu_B = 104$, $\sigma_A = 6$, $\sigma_B = 8$ ps and $\rho = 0$, $\theta = 10$, $\alpha = 0.6$ and $T_A = 0.726$. Correlation shrinks $\theta$, which sharpens the competition: two highly correlated paths with different means are almost never reordered.

## Clark's mean of the maximum

Clark (1961) gave the exact first two moments of $\max(A, B)$:
$$E[\max] = \mu_A\Phi(\alpha) + \mu_B\Phi(-\alpha) + \theta\varphi(\alpha),$$
$$E[\max^2] = (\mu_A^2 + \sigma_A^2)\Phi(\alpha) + (\mu_B^2 + \sigma_B^2)\Phi(-\alpha) + (\mu_A + \mu_B)\theta\varphi(\alpha).$$
In the example, $\varphi(0.6) = 0.333$, so $E[\max] = 110(0.726) + 104(0.274) + 10(0.333) = 111.7$ ps, above both means. The extra $\theta\varphi(\alpha)$ term is largest when the two inputs are evenly matched ($\alpha = 0$), where it equals $0.40\,\theta$: reconvergence of two equally critical paths pushes the arrival later than either alone. A deterministic timer that takes the max of means misses that shift entirely.

## The variance, and why max is not Gaussian

The variance of the max is $E[\max^2] - E[\max]^2$. A block-based timer then approximates $\max(A, B)$ by a Gaussian with these two moments, and it propagates that Gaussian onward. The true max is skewed: when the inputs are evenly matched the distribution has a heavier right tail than a Gaussian, and the error compounds over many reconvergent max operations. The moment match is exact in the first two moments, not in the tails, which is exactly where timing yield is decided. Engines refine this with skew-normal or quadratic canonical forms, or with Monte Carlo spot-checks of critical endpoints. The other approximation, re-expressing the max in the same variation basis, is the subject of canonical forms.

## Going deeper

Clark, "The greatest of a finite set of random variables," Operations Research, 1961; Visweswariah et al., "First-order incremental block-based statistical timing analysis," DAC 2004.
