# Variance-aware surrogates in an agent loop

In an agent-first flow, fast engines run continuously and accurate engines are called only where they change a decision. A fast surrogate that returns only a prediction gives the agent no basis for that choice. A surrogate that also returns a calibrated predictive uncertainty does: the agent can compute the probability that the accurate engine would disagree, and spend accurate runs exactly where that probability is high.

## A prediction with an uncertainty

A variance-aware surrogate, such as a Gaussian process, an ensemble or a model with a learned variance head, returns for each query a mean $\hat{\mu}$ and a predictive standard deviation $\hat{\sigma}$. Treating the accurate engine's answer as $N(\hat{\mu}, \hat{\sigma}^2)$, the probability a path with predicted slack $\hat{\mu}$ actually violates is
$$P(\text{slack} < 0) = \Phi(-\hat{\mu}/\hat{\sigma}).$$
A prediction of $+12$ ps with $\hat{\sigma} = 8$ ps gives $\Phi(-1.5) = 6.7$ percent; the same $+12$ ps with $\hat{\sigma} = 2$ ps gives $10^{-9}$. The point predictions are identical and the decisions are opposite. This is the statistical content of the fast-versus-accurate loop: the variance, not the mean, says when to recurse.

## A recursion rule

A simple policy for the loop: trust the surrogate's verdict when $|\hat{\mu}| \ge k\hat{\sigma}$ and call the accurate engine otherwise. With $k = 2$ the worst accepted case has a one-sided error probability of $Q(2) = 2.3$ percent; with $k = 3$, 0.13 percent. Raising $k$ sends more candidates to the accurate engine and lowers the false-pass risk. The expected number of false passes among accepted candidates is $\sum \Phi(-\hat{\mu}_j/\hat{\sigma}_j)$ over those accepted as passing, which gives a budget the agent can track. In this design the loop's parameters, $k$ and the accurate-engine budget, are configuration, while the surrogate model is a selectable skill with its own validation record.

## Calibration is the precondition

Everything above assumes the predictive sigma is calibrated: errors larger than $2\hat{\sigma}$ should occur about 5 percent of the time. Surrogates are typically overconfident outside their training data, exactly where an optimizing agent drives them, because an optimizer seeks regions the surrogate rates as good. The checks are statistical: on held-out accurate runs, the standardized errors $(y - \hat{\mu})/\hat{\sigma}$ should look standard normal, with coverage of the $\pm 2\hat{\sigma}$ band near 95 percent. A golden test for a surrogate skill therefore asserts coverage, not only RMSE, and the loop should widen $\hat{\sigma}$ or force recursion for inputs far from the training distribution.

## Going deeper

Rasmussen and Williams, "Gaussian Processes for Machine Learning" (2006), chapters 2 and 5; Kuleshov, Fenner and Ermon, "Accurate uncertainties for deep learning using calibrated regression," ICML 2018.
