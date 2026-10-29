# Timing yield from a Gaussian delay

Statistical timing changes the question from whether a path passes to how many dies pass. If the circuit delay is approximately Gaussian with mean $\mu$ and sigma $\sigma$, the fraction of dies meeting a clock period $T$ is a normal probability. Turning this around gives the period needed for a yield target, the statistical counterpart of a signoff margin.

## The z-score of a period

Standardize the period against the delay distribution: $z = (T - \mu)/\sigma$. It counts how many sigma of margin the period leaves above the mean delay. With $\mu = 800$ ps, $\sigma = 39$ ps and $T = 840$ ps, $z = 40/39 = 1.03$. A negative $z$ means the period is shorter than the mean delay, so most dies fail. The delay distribution here is the distribution of the circuit's worst delay, the statistical max over all endpoints, which SSTA approximates in canonical form; a single path's distribution would overstate the yield, since any failing path fails the die.

## Yield at a given period

The yield is $Y = P(D \le T) = \Phi(z)$. For $\mu = 720$ ps, $\sigma = 37$ ps and $T = 712.6$ ps, $z = -0.2$ and $\Phi(-0.2) = 0.4207$, so 42.1 percent of dies meet the period. Some anchors from the standard normal table: $\Phi(1) = 0.8413$, $\Phi(2) = 0.9772$ and $\Phi(3) = 0.99865$. Yield is not linear in margin: going from 0 to 1 sigma raises yield from 50 to 84 percent, while going from 2 to 3 sigma adds only about 2 percentage points. This is the trade-off a design team makes when it chooses the number of sigma for signoff.

## The period for a target yield

Invert the relation: the shortest period meeting a yield target $p$ is $T = \mu + z_p \sigma$, where $z_p$ is the standard normal quantile at $p$, for example 1.645 for 95 percent, 2.326 for 99 percent and 3.090 for 99.9 percent. With $\mu = 1000$ ps and $\sigma = 13$ ps, 99.9 percent yield needs $1000 + 3.090 \times 13 = 1040.2$ ps. The Gaussian tail is an approximation: the true delay distribution of a circuit, a max over many correlated paths, is skewed to the right, so high-yield quantiles from a Gaussian fit can be optimistic. Monte Carlo on a sample of dies is the usual check.

## Going deeper

Sapatnekar, "Timing" (2004), chapter on statistical timing; Blaauw, Chopra, Srivastava and Scheffer, "Statistical timing analysis: from basic principles to state of the art" (IEEE TCAD, 2008).
