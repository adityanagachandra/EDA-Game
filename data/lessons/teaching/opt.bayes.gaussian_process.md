# Gaussian process surrogates

A Gaussian process is a probability model over functions, used as the surrogate in Bayesian optimization. It predicts both a mean and an uncertainty at any untried input, from a handful of expensive evaluations such as SPICE runs or full place-and-route jobs. Everything follows from a kernel that says how strongly nearby inputs are correlated, and from the conditioning formulas for jointly Gaussian variables.

## The kernel as covariance

A zero-mean Gaussian process models $f(x)$ so that any finite set of values is jointly Gaussian with covariance $\text{Cov}(f(x), f(x')) = k(x, x')$. The squared-exponential kernel is
$$k(x, x') = \sigma^2 \exp\!\left(-\frac{(x - x')^2}{2\ell^2}\right),$$
with signal variance $\sigma^2$ and length scale $\ell$. With $\sigma^2 = 2$ and $\ell = 0.5$: $k(1.3, 1.6) = 2\exp(-0.09/0.5) = 2e^{-0.18} = 1.671$.

The length scale sets how far knowledge travels. At a distance of one length scale the correlation is $e^{-0.5} = 0.61$; at three, about 0.011. A flow knob such as target utilization might have a length scale of a few percent, a transistor width a fraction of a micron. Other kernels, such as Matern, give rougher functions, and per-input length scales (automatic relevance determination) reveal which knobs matter.

## The posterior mean from one observation

Observe $y_1$ at $x_1$, possibly with noise variance $\sigma_n^2$. Conditioning the joint Gaussian of $(f(x), y_1)$ gives the posterior mean
$$\mu(x) = \frac{k(x, x_1)}{k(x_1, x_1) + \sigma_n^2}\, y_1.$$
With $\sigma^2 = 3$, $\ell = 1.5$, a noise-free observation $y_1 = 0.9$ at $x_1 = 1.6$, and $x = 1.2$: $k = 3\exp(-0.16/4.5) = 2.895$, so $\mu = (2.895/3) \cdot 0.9 = 0.869$.

The prediction is a weighted copy of the observation, weighted by correlation. Far from the data it falls back to the prior mean, here 0. With n observations the formula becomes $\mu(x) = k_x^T (K + \sigma_n^2 I)^{-1} y$, where K is the n by n kernel matrix. Solving that system costs $O(n^3)$, which is fine for the hundreds of runs typical of flow tuning but not for millions of points.

## The posterior variance

The posterior variance is
$$\sigma^2(x) = k(x, x) - \frac{k(x, x_1)^2}{k(x_1, x_1) + \sigma_n^2}.$$
With $\sigma^2 = 3$, $\ell = 0.5$, noise variance 0.5, $x_1 = 1.7$ and $x = 1.6$: $k(x, x_1) = 3e^{-0.02} = 2.941$, so $\sigma^2(x) = 3 - 8.648/3.5 = 0.529$. Without noise it would be $3 - 8.648/3 = 0.117$, and exactly 0 at $x_1$ itself.

The variance depends only on where the data are, not on the measured values. It is small near evaluated points and returns to the prior variance far away, which is the signal Bayesian optimization uses to explore. The noise term matters in EDA because tool runs vary with seeds: with $\sigma_n^2 = 0$ the surrogate interpolates seed noise exactly and becomes overconfident. Kernel hyperparameters ($\sigma^2$, $\ell$, $\sigma_n^2$) are fitted by maximizing the marginal likelihood of the data.

## Going deeper

Rasmussen and Williams, "Gaussian Processes for Machine Learning" (2006), chapters 2 and 5.
