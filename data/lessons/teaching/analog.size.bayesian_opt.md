# Bayesian optimization for sizing

Bayesian optimization is the standard sample-efficient method for simulation-based sizing. It fits a Gaussian-process surrogate to the simulations done so far and chooses the next SPICE run by an acquisition function that balances exploiting promising regions against exploring uncertain ones. Each choice is cheap compared with a simulation, so the method spends compute deciding where to look.

## Probability of improvement

After each simulation the Gaussian process gives, at any untried sizing, a predicted mean $\mu$ and standard deviation $\sigma$ of the performance. When minimizing power with best feasible value $f^*$ so far, the probability that a candidate improves on it is $\mathrm{PI} = \Phi((f^* - \mu)/\sigma)$, with $\Phi$ the standard normal CDF. For $f^* = 1.2$ mW, $\mu = 1.1$ and $\sigma = 0.2$, $z = 0.5$ and $\mathrm{PI} = 0.691$. PI favours candidates predicted just below the incumbent with high confidence, so it tends to exploit and can stall near a local optimum; practical sizers prefer expected improvement or confidence bounds.

## Expected improvement

Expected improvement weights each possible improvement by its size: $\mathrm{EI} = (f^* - \mu)\Phi(z) + \sigma\varphi(z)$ with $z = (f^* - \mu)/\sigma$ and $\varphi$ the standard normal density. The first term rewards a good mean, the second rewards uncertainty. For $f^* = 1.2$, $\mu = 1.25$, $\sigma = 0.1$: $z = -0.5$, $\Phi = 0.3085$, $\varphi = 0.3521$, so $\mathrm{EI} = -0.05 \cdot 0.3085 + 0.1 \cdot 0.3521 = 0.0198$ mW. A candidate predicted to be worse than the incumbent still has positive EI because it might be better. Constraints are handled by multiplying EI by the surrogate's probability that every spec is met, which concentrates the search on likely feasible regions.

## Confidence bounds and practical limits

The lower confidence bound $\mu - \kappa\sigma$ (for minimization) gives a tunable trade-off: larger $\kappa$ explores more. With $\kappa = 2$, a candidate at $\mu = 1.3$, $\sigma = 0.2$ has bound 0.9 and beats one at $\mu = 1.0$, $\sigma = 0.02$ with bound 0.96, though its mean is worse. The method's limits shape how sizers use it. Gaussian-process fitting scales cubically with the number of samples and struggles beyond a few dozen dimensions, so sizers use sparse or local models, trust regions, or batch acquisition to run several SPICE jobs in parallel. Corners and yield multiply the cost of each evaluation, so the surrogate is often built on worst-case values. The final answer is always re-simulated, since the surrogate is a guide, not a judge.

## Going deeper

Shahriari et al., "Taking the human out of the loop: a review of Bayesian optimization" (Proceedings of the IEEE, 2016); Lyu et al., "An efficient Bayesian optimization approach for automated optimization of analog circuits" (IEEE TCAS-I, 2018).
