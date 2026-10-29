# Acquisition functions: EI and UCB

Bayesian optimization chooses the next expensive run by maximizing an acquisition function computed from the GP posterior. Expected improvement and confidence bounds are the common choices. Each balances exploiting points with good predicted values against exploring points with high uncertainty, and the balance decides how many runs the search needs.

## Expected improvement

When minimizing, with best observed value $f^*$ and a Gaussian posterior with mean $\mu$ and standard deviation $\sigma$ at a candidate, the expected improvement is $\mathrm{EI} = (f^* - \mu)\Phi(z) + \sigma\varphi(z)$ with $z = (f^* - \mu)/\sigma$. For $f^* = 70$ mW, $\mu = 67$ mW and $\sigma = 5$ mW, $z = 0.6$, $\Phi(z) = 0.726$ and $\varphi(z) = 0.333$, so $\mathrm{EI} = 3 \cdot 0.726 + 5 \cdot 0.333 = 3.84$ mW. EI is the expectation of $\max(0, f^* - f)$, so it is never negative, and it is positive even where the mean is worse than $f^*$, as long as there is uncertainty.

## Confidence bounds

The lower confidence bound for minimization, $\mu - \kappa\sigma$, scores a candidate by an optimistic estimate of its value; the next run goes to the candidate with the smallest bound. With $\kappa = 1$, a candidate with mean 64 and standard deviation 4 has bound 60, and one with mean 75 and standard deviation 9 has bound 66, so the first wins. With $\kappa = 3$ the bounds are 52 and 48 and the uncertain one wins. $\kappa$ sets exploration explicitly. Maximization uses the upper bound $\mu + \kappa\sigma$ (UCB); confusing the two directions sends the search toward the worst candidates.

## Comparing candidates

EI trades mean against uncertainty automatically. With $f^* = 65$ uW, a candidate with mean 64 and standard deviation 3 has EI 1.76, one with mean 66 and standard deviation 4 has EI 1.15, and one with mean 61 and standard deviation 2 has EI 4.02. The third, with the best mean, wins here, but a candidate with a worse mean can win when its standard deviation is large enough. EI tends to exploit once the posterior is confident, so practical loops add exploration, run batches of candidates in parallel, and stop on a budget rather than when EI becomes small, since a miscalibrated GP can report small EI everywhere while missing the optimum.

## Going deeper

Frazier, "A Tutorial on Bayesian Optimization" (2018); Jones, Schonlau and Welch, "Efficient Global Optimization of Expensive Black-Box Functions," Journal of Global Optimization (1998).
