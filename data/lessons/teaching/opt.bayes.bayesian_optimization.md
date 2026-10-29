# Bayesian optimization

Bayesian optimization finds a good input for an expensive black-box function in few evaluations. It fits a surrogate, usually a Gaussian process, to the runs done so far, and picks the next run by maximizing an acquisition function that trades predicted value against uncertainty. In EDA it tunes analog transistor sizes with SPICE in the loop and flow parameters with full implementation runs, where each evaluation costs minutes to hours and no gradient is available.

## The upper confidence bound

For maximization, the UCB acquisition scores a candidate by $\mu(x) + \kappa\sigma(x)$, the predicted mean plus $\kappa$ predictive standard deviations. With $\kappa = 2.5$ and candidates x1 (mean 2.1, sd 1.4), x2 (1.3, 0.9), x3 (0.3, 1.7) and x4 (3.6, 0.4), the scores are 5.6, 3.55, 4.55 and 4.6. The loop runs x1, whose mean is not the best but whose uncertainty is large.

$\kappa$ sets the balance. Small $\kappa$ exploits, sampling where the mean is high; large $\kappa$ explores, sampling where little is known. Theory gives schedules for $\kappa$ that grow slowly with the iteration count and guarantee convergence, but in practice a constant between 1 and 3 is common. For minimization the score is $\mu - \kappa\sigma$, and the loop picks the smallest.

## Expected improvement

For minimization with best observed value $f^*$, expected improvement is $\text{EI}(x) = \mathbb{E}[\max(0, f^* - f(x))]$. Under the Gaussian posterior it has a closed form:
$$\text{EI} = (f^* - \mu)\Phi(z) + \sigma\varphi(z), \qquad z = \frac{f^* - \mu}{\sigma},$$
with $\Phi$ and $\varphi$ the standard normal distribution and density. With $f^* = 2.0$, $\mu = 2.3$ and $\sigma = 0.6$: $z = -0.5$, $\Phi = 0.3085$, $\varphi = 0.3521$, so $\text{EI} = -0.3 \cdot 0.3085 + 0.6 \cdot 0.3521 = 0.119$.

EI is positive even when the mean is worse than the incumbent, because the uncertainty might hide an improvement. The first term rewards a good mean and the second rewards uncertainty, so EI balances exploitation and exploration without a tuning constant. Probability of improvement, $\Phi(z)$, ignores how large the improvement is and tends to over-exploit.

## Exploit or explore

Compare a safe candidate P (mean 1.4, sd 0.3) with an uncertain one Q (mean 2.7, sd 0.8) when the incumbent is $f^* = 1.9$, minimizing. For P: $z = 1.667$, $\Phi = 0.952$, $\varphi = 0.0994$, $\text{EI} = 0.5 \cdot 0.952 + 0.3 \cdot 0.0994 = 0.506$. For Q: $z = -1$, $\Phi = 0.159$, $\varphi = 0.242$, $\text{EI} = -0.8 \cdot 0.159 + 0.8 \cdot 0.242 = 0.067$. P wins. With Q's sd at 2.5, its EI would be 0.648 and it would win: wide spread can outweigh a worse mean.

The loop alternates fitting and acquisition until the budget runs out, and reports the best evaluated point, not the surrogate's predicted optimum. Maximizing the acquisition is itself an optimization, cheap because the surrogate is cheap, done by multi-start gradient methods or by scoring many random candidates.

## Going deeper

Jones, Schonlau and Welch, "Efficient global optimization of expensive black-box functions" (Journal of Global Optimization, 1998); Srinivas, Krause, Kakade and Seeger, "Gaussian process optimization in the bandit setting" (ICML, 2010); Frazier, "A tutorial on Bayesian optimization" (arXiv, 2018).
