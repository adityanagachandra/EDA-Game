# Extreme value tails

When the question is about the far tail of a distribution, the body of the distribution is irrelevant and often misleading. Extreme value theory says that, under mild conditions, exceedances over a high threshold follow a generalized Pareto distribution whatever the parent distribution. That lets an engine fit the tail from a modest number of tail samples and extrapolate to rates no simulation could reach, which is how statistical blockade turns a few hundred tail simulations into a 6-sigma estimate.

## The generalized Pareto tail

Choose a high threshold $u$ with exceedance probability $\zeta_u = P(X > u)$. The Pickands-Balkema-de Haan theorem says that for a wide class of distributions, the excess $Y = X - u$ given $X > u$ is approximately generalized Pareto with scale $\sigma$ and shape $\xi$:
$$P(X > u + y) = \zeta_u\left(1 + \frac{\xi y}{\sigma}\right)^{-1/\xi}.$$
With $u = 60$ ps, $\zeta_u = 0.01$, $\sigma = 5$ ps and $\xi = 0.2$, the chance of exceeding 80 ps is $0.01(1.8)^{-5} = 5.3 \times 10^{-4}$. The fit needs only samples above $u$, and the extrapolation is a closed form.

## The shape parameter

The shape $\xi$ classifies the tail. $\xi > 0$ is a heavy, power-law tail; $\xi < 0$ is a tail with a finite end point; $\xi = 0$ is the limit $P(X > u + y) = \zeta_u e^{-y/\sigma}$, an exponential tail. Gaussian and lognormal parents both fall in the $\xi = 0$ domain of attraction, so a fitted $\xi$ near zero is what a circuit metric driven by Gaussian parameters should show. With $\xi = 0$, $\zeta_u = 0.02$ and $\sigma = 4$ ps, exceeding $u + 12$ ps has probability $0.02e^{-3} = 9.96 \times 10^{-4}$. A fitted $\xi$ that is clearly positive is a warning: either the metric has a genuinely heavy tail, such as a near-metastable latch, or the threshold is too low.

## Return levels and threshold choice

Inverting the tail gives the level exceeded with a small probability $p$:
$$x_p = u + \frac{\sigma}{\xi}\left[\left(\frac{\zeta_u}{p}\right)^{\xi} - 1\right],$$
which is how a spec value at a target failure rate is read off. The threshold is the key engineering choice. Too low, and the asymptotic theory does not yet hold, biasing the fit; too high, and too few samples remain, inflating variance. Practice looks for a range of thresholds over which the fitted $\xi$ and the modified scale are stable. Extrapolating many orders of magnitude beyond the data magnifies small shape errors enormously, so a tail model should report an uncertainty band, not a single number.

## Going deeper

Coles, "An Introduction to Statistical Modeling of Extreme Values" (2001), chapter 4; Singhee and Rutenbar, "Statistical blockade," IEEE TCAD, 2009.
