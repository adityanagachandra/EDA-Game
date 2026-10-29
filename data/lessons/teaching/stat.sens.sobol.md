# Variance-based sensitivity indices

First-order sensitivities assume the response is linear. When parameters interact, for example when a mismatch matters only at low supply, derivatives at nominal under-report importance. Sobol indices measure importance directly as shares of output variance, with no linearity assumption, and they separate a parameter's own effect from its interactions. They are the standard global sensitivity measure in uncertainty quantification.

## First-order indices

For independent inputs, the variance of $Y = f(X_1, \ldots, X_d)$ decomposes into contributions of single inputs, pairs, triples and so on. The first-order index of $X_i$ is
$$S_i = \frac{\mathrm{Var}(E[Y \mid X_i])}{\mathrm{Var}(Y)},$$
the fraction of variance removed, on average, by learning $X_i$ alone. For an additive linear model $Y = \sum a_i X_i$, $E[Y \mid X_i] = a_i X_i$ (plus a constant), so $S_i = a_i^2\sigma_i^2/\mathrm{Var}(Y)$, which is exactly the first-order variance share. With $Y = 3X_1 + 2X_2 + X_3$ and unit sigmas, $S_1 = 9/14 = 0.64$. For an additive model the first-order indices sum to 1.

## Interactions

Take $Y = a_1X_1 + a_2X_2 + cX_1X_2$ with independent zero-mean Gaussian inputs. The product term has variance $c^2\sigma_1^2\sigma_2^2$ and is uncorrelated with both linear terms, so $\mathrm{Var}(Y) = a_1^2\sigma_1^2 + a_2^2\sigma_2^2 + c^2\sigma_1^2\sigma_2^2$. Conditioning on $X_1$ averages $X_2$ to zero, so $E[Y \mid X_1] = a_1X_1$ and $S_1 = a_1^2\sigma_1^2/\mathrm{Var}(Y)$. With $a_1 = a_2 = 1$, $c = 2$ and unit sigmas, $\mathrm{Var}(Y) = 6$ and $S_1 = S_2 = 1/6$. The first-order indices now sum to only $1/3$; the remaining $4/6$ is interaction variance that no single input explains. A derivative at nominal would see $\partial Y/\partial X_1 = 1$ and miss the interaction entirely.

## Total indices and estimation cost

The total index $S_{Ti} = 1 - \mathrm{Var}(E[Y \mid X_{\sim i}])/\mathrm{Var}(Y)$ counts every term involving $X_i$: its first-order effect plus all its interactions. In the example, $S_{T1} = (1 + 4)/6 = 0.83$. A parameter with $S_{Ti}$ near zero can be fixed at nominal without changing the output variance, which is the practical use: shrinking a 200-parameter Monte Carlo to the 15 parameters that matter. Estimating the indices by sampling (Saltelli or Jansen estimators) costs about $N(d + 2)$ model runs for $d$ inputs, so for circuits they are usually computed on a fitted surrogate, such as a polynomial chaos expansion, whose coefficients give the indices in closed form.

## Going deeper

Saltelli et al., "Global Sensitivity Analysis: The Primer" (2008), chapter 4; Sobol, "Global sensitivity indices for nonlinear mathematical models and their Monte Carlo estimates," Mathematics and Computers in Simulation, 2001.
