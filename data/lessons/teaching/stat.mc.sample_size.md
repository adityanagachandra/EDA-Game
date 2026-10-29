# Monte Carlo sample size for a failure rate

Before launching a Monte Carlo run, a flow must decide how many samples to ask for. The answer depends on the failure rate itself, on the relative accuracy needed, and on whether the goal is to estimate a rate or merely to demonstrate it is below a bound. The formulas are short, and they reveal that plain Monte Carlo is hopeless for the failure rates memory design cares about.

## Samples for a relative error

For a failure probability $p$, the estimate $\hat{p}$ has standard error $\sqrt{p(1-p)/N}$, so its relative error is $\sqrt{(1-p)/(Np)}$. Requiring a relative half-width $\varepsilon$ at confidence level $z$ gives
$$N = \frac{z^2 (1 - p)}{p\,\varepsilon^2}.$$
For $p = 0.001$, $\varepsilon = 10$ percent and $z = 1.96$, $N = 3.84 \times 0.999/(0.001 \times 0.01) = 383{,}776$. The cost is inversely proportional to $p$: the rarer the failure, the more samples are needed, because most samples carry no information about the failure. A useful mental shortcut is that about $z^2/\varepsilon^2 \approx 400$ failures must be observed for 10 percent accuracy at 95 percent confidence.

## Demonstrating a bound with zero failures

Sometimes the goal is a claim, not an estimate: the failure rate is below $p_0$ with confidence $c$. If the true rate were $p_0$, the chance of seeing zero failures in $N$ samples is $(1 - p_0)^N$. Requiring that chance to be at most $1 - c$ gives
$$N \ge \frac{\ln(1 - c)}{\ln(1 - p_0)} \approx \frac{-\ln(1 - c)}{p_0}.$$
At 95 percent confidence, $-\ln 0.05 = 3.0$, which is the rule of three: zero failures in $3/p_0$ samples. To claim $p < 10^{-4}$ at 95 percent needs about 30,000 clean samples. Observing even one failure raises the requirement, and the bound then comes from the Poisson or binomial tail with that count.

## Why plain Monte Carlo fails at high sigma

Combine the formula with the sigma-to-probability conversion. A 4.5-sigma failure ($p = 3.4 \times 10^{-6}$) at 10 percent accuracy needs about $1.1 \times 10^8$ simulations. A 6-sigma bitcell ($p = 9.9 \times 10^{-10}$) needs about $3.9 \times 10^{11}$. At even a millisecond per transient simulation that is twelve years of CPU time for one cell at one corner. The conclusion is not that more compute will fix it: the estimator wastes almost every sample in the region where nothing fails. Rare-event methods concentrate samples where failures occur, and their savings are reported against exactly these baseline numbers.

## Going deeper

Hahn and Meeker, "Statistical Intervals: A Guide for Practitioners" (1991), on binomial and zero-failure intervals; Kanj, Joshi and Nassif, DAC 2006, for the SRAM motivation.
