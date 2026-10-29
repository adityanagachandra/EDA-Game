# Gaussian-process posteriors

A Gaussian process is the standard surrogate in Bayesian optimization for analog sizing and flow tuning, because it gives a mean and a calibrated uncertainty at every point from a handful of expensive runs. Its posterior has closed-form mean and variance. Computing them for one and two observations shows exactly how data pulls the prediction and shrinks the uncertainty.

## The posterior mean

A GP prior has mean $m$ and kernel $k(x, x')$, for example the squared exponential $\sigma_f^2 \exp(-(x - x')^2/(2\ell^2))$. After one noisy observation $y_1$ at $x_1$ with noise variance $\sigma_n^2$, the posterior mean at $x^*$ is $\mu^* = m + \frac{k(x^*, x_1)}{k(x_1, x_1) + \sigma_n^2}(y_1 - m)$. With $m = 0$, $\sigma_f^2 = 1$, $\ell = 1$, $\sigma_n^2 = 0.1$, $y_1 = 2$ at $x_1 = 0$ and $x^* = 1$: $k(x^*, x_1) = e^{-0.5} = 0.607$, so $\mu^* = 0.607 \cdot 2/1.1 = 1.10$. Far from the data the kernel decays and the mean returns to the prior.

## The posterior variance

The posterior variance at $x^*$ is $\sigma^{*2} = k(x^*, x^*) - k(x^*, x_1)^2/(k(x_1, x_1) + \sigma_n^2)$. In the example, $1 - 0.368/1.1 = 0.666$. Two facts follow from the formula. The variance does not depend on the observed value $y_1$, only on where the data are, which is why GP-based active learning can plan where to sample before running anything. And with noise, the variance at an observed point does not fall to zero: at $x^* = x_1$ it is $1 - 1/1.1 = 0.09$, reflecting that a noisy run does not pin the function down exactly.

## Two observations

With observations $y$ at points $X$, the posterior mean is $\mu^* = m + k_*^\top (K + \sigma_n^2 I)^{-1}(y - m)$, where $K$ is the kernel matrix of the data and $k_*$ the kernel between the query and each point. With zero mean, no noise, $y = (1, 2)$ at $x = 0$ and 1, and $x^* = 0.5$: $K$ has ones on the diagonal and $e^{-0.5} = 0.607$ off it, $\alpha = K^{-1}y = (-0.337, 2.204)$ and $k_* = (0.882, 0.882)$, so $\mu^* = 1.648$. Correlated data points share credit through $K^{-1}$. Each solve costs $O(n^3)$ in the number of runs, harmless for the tens to hundreds of runs an analog loop affords.

## Going deeper

Rasmussen and Williams, "Gaussian Processes for Machine Learning" (2006), chapter 2.
