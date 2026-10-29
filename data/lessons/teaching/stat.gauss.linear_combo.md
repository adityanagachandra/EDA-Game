# Sums of correlated Gaussians

A path delay is a sum of stage delays, and a statistical timer must know its mean and spread. If the stage delays are Gaussian and linear in shared variation sources, the path delay is Gaussian too, and its variance follows from the covariance algebra alone. Getting the correlation terms right is the entire difference between a realistic path sigma and one that is off by a factor of several.

## Variance of a sum

For any two random variables, $\mathrm{Var}(X + Y) = \sigma_X^2 + \sigma_Y^2 + 2\rho\sigma_X\sigma_Y$, and the mean is simply $\mu_X + \mu_Y$. If $X$ and $Y$ are jointly Gaussian, the sum is exactly Gaussian with those moments. Two stages with $\sigma = 6$ and 8 ps have a sum sigma of 10 ps if independent, 14 ps if perfectly correlated, and $\sqrt{36 + 64 + 2(0.5)(48)} = 12.2$ ps at $\rho = 0.5$. A negative correlation, which arises when one parameter speeds one stage and slows another, shrinks the sum below the independent value. More generally, for $D = a^T X$ with covariance $\Sigma$, $\mathrm{Var}(D) = a^T \Sigma a$.

## A path of n correlated stages

For $n$ identical stages with sigma $\sigma$ and common pairwise correlation $\rho$, the covariance matrix has $\sigma^2$ on the diagonal and $\rho\sigma^2$ in all $n(n-1)$ off-diagonal places:
$$\mathrm{Var}(D) = n\sigma^2 + n(n-1)\rho\sigma^2.$$
For $n = 10$, $\sigma = 4$ ps and $\rho = 0.3$, the variance is $160 + 432 = 592$, and the path sigma is 24.3 ps. The independent answer would be 12.6 ps and the fully correlated answer 40 ps. Even a modest correlation dominates for long paths, because the off-diagonal terms grow as $n^2$ while the diagonal grows as $n$.

## Relative spread and logic depth

The ratio of spread to mean shows how variation scales with logic depth. With independent stages the mean grows as $n$ and the sigma as $\sqrt{n}$, so $\sigma/\mu$ falls as $1/\sqrt{n}$: deep paths average out random variation. With fully correlated stages both grow as $n$ and the ratio stays fixed. Ten independent stages of mean 25 ps and sigma 3 ps have a 3-sigma spread of $3 \times 9.5 = 28.5$ ps on a 250 ps mean, 11.4 percent; the correlated case is 36 percent. This is the reasoning behind depth-dependent derating (AOCV): a flat on-chip-variation derate that is right for a short path is pessimistic for a deep one when the variation is random. It is also why the random and correlated parts must be kept separate in the delay representation.

## Going deeper

Blaauw, Chopra, Srivastava and Scheffer, "Statistical timing analysis: from basic principles to state of the art," IEEE TCAD, 2008.
