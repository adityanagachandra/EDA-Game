# Sensitivity analysis

Before spending a Monte Carlo budget, an engine wants to know which variation sources matter. Sensitivity analysis answers that from derivatives: linearize the output about nominal, weight each derivative by its parameter's sigma, and the result is both a fast estimate of the output's spread and a ranking of where that spread comes from. It is the cheapest form of uncertainty quantification, and it tells a designer or an optimizer which device to fix.

## First-order variance propagation

Let a delay $D$ depend on independent parameters $p_i$ with standard deviations $\sigma_i$. A first-order Taylor expansion about nominal gives $D \approx D_0 + \sum_i s_i\,\Delta p_i$ with sensitivities $s_i = \partial D/\partial p_i$, and therefore
$$\sigma_D^2 \approx \sum_i s_i^2 \sigma_i^2.$$
With $s_1 = 1.5$ ps/mV for threshold voltage ($\sigma_1 = 10$ mV) and $s_2 = 2$ ps/nm for gate length ($\sigma_2 = 4$ nm), $\sigma_D^2 = 225 + 64 = 289$ and $\sigma_D = 17$ ps. If the parameters are correlated, the sum becomes $s^T \Sigma s$. The approximation is exact for a linear response with Gaussian inputs, and degrades as the response curves over the range of the variation.

## Ranking the contributors

Each term $s_i^2\sigma_i^2$ is the variance that parameter $i$ contributes, so its share is $s_i^2\sigma_i^2/\sigma_D^2$. In the example, threshold voltage contributes $225/289 = 78$ percent and gate length 22 percent. Note what this ranks: not the raw derivative, which depends on units, and not the sigma alone, but their product. A parameter with a large derivative and a tiny spread can matter less than one with a modest derivative and a large spread. The sign of the sensitivity does not affect the share, though it matters for correlated parameters and for which direction a worst case pushes. This ranking is what tells an optimizer whether to upsize a device to reduce its mismatch or to change the topology.

## Estimating sensitivities by simulation

A SPICE engine rarely provides analytic derivatives with respect to every process parameter, so sensitivities come from perturbed simulations. The central difference $s \approx (D(p + h) - D(p - h))/(2h)$ has error of order $h^2$, against order $h$ for a one-sided difference, at the cost of one extra simulation per parameter. If raising $V_{th}$ by 5 mV gives 82.4 ps and lowering it gives 77.6 ps, $s = 4.8/10 = 0.48$ ps/mV. The step $h$ must be large enough that simulator noise and convergence tolerances do not swamp the difference, and small enough that curvature does not; around one sigma is a common choice. For $d$ parameters a central-difference sensitivity analysis costs $2d$ runs, far fewer than Monte Carlo for small $d$; adjoint sensitivity analysis gets all derivatives for roughly the cost of one extra solve.

## Going deeper

Saltelli et al., "Global Sensitivity Analysis: The Primer" (2008), chapter 1; the sensitivity analysis sections of the ngspice manual.
