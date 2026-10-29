# Clark's max of two Gaussian arrivals

Sums keep canonical forms Gaussian; maxima do not. At every merge pin a statistical timer must take the max of two random arrivals and approximate the result again by a Gaussian in canonical form. Clark's 1961 formulas give the exact mean and variance of the max of two jointly Gaussian variables, and block-based SSTA uses them, together with a tightness probability, to re-express the max.

## The spread of the difference

For arrivals $A$ and $B$ with means $\mu_A, \mu_B$, sigmas $\sigma_A, \sigma_B$ and correlation $\rho$, the key quantity is the standard deviation of their difference:
$$\theta = \sqrt{\sigma_A^2 + \sigma_B^2 - 2\rho\,\sigma_A\sigma_B}.$$
For $\sigma_A = 20$, $\sigma_B = 25$ ps and $\rho = 0.5$, $\theta^2 = 400 + 625 - 500 = 525$ and $\theta = 22.91$ ps. Correlation shrinks $\theta$: two perfectly correlated arrivals with equal sigma have $\theta = 0$, and their max is simply the one with the larger mean. The normalized separation is $\alpha = (\mu_A - \mu_B)/\theta$; a large $|\alpha|$ means one input almost always dominates.

## The tightness probability

The probability that $A$ is the later arrival is $P(A > B) = \Phi(\alpha)$, the standard normal distribution function at $\alpha$. It is called the tightness probability of $A$. For means 300 and 320 ps, sigmas 16 and 29 ps and $\rho = 0$, $\theta = 33.12$ ps, $\alpha = -0.604$ and $\Phi(\alpha) = 0.27$. In the canonical-form max, the sensitivities of the result are the tightness-weighted mix of the inputs' sensitivities: $c_i = \Phi(\alpha) a_i + \Phi(-\alpha) b_i$. The tightness also tells a statistical optimizer how much each input contributes to the output's variation, which is the statistical analogue of knowing which fanin is critical.

## The mean of the max

Clark's mean is
$$E[\max(A,B)] = \mu_A \Phi(\alpha) + \mu_B \Phi(-\alpha) + \theta\,\varphi(\alpha),$$
where $\varphi$ is the standard normal density. The last term is always positive, so the mean of the max exceeds the tightness-weighted mean of the inputs and, when the inputs overlap, exceeds both means. For means 520 and 510 ps, $\theta = 13.54$ ps, $\Phi(\alpha) = 0.77$ and $\varphi(\alpha) = 0.3037$, the mean is 521.81 ps. Clark also gives the second moment, from which the variance follows; the timer then fixes the random coefficient so that the canonical result has that variance. The max of a Gaussian pair is not Gaussian, and repeating this approximation at every merge accumulates error, especially when $|\alpha|$ is small.

## Going deeper

Clark, "The greatest of a finite set of random variables" (Operations Research, 1961); Visweswariah et al., "First-order incremental block-based statistical timing analysis" (DAC 2004).
