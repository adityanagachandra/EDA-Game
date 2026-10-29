# Fast fault-sampled coverage in an agent loop

An agent exploring DFT configurations cannot afford full fault simulation of millions of faults for every candidate. A fast mode simulates a uniform random sample of faults and reports the sample coverage with a confidence interval. Knowing the interval tells the agent when the fast answer is good enough and when it must recurse to the accurate engine.

## Sampled coverage

If $n$ faults are drawn uniformly from the full list and $d$ are detected, $\hat{p} = d/n$ is an unbiased estimate of the full coverage. Its cost is proportional to $n$, not to the full fault count: sampling 10,000 of 2 million faults is 200 times cheaper. With 9864 of 10,000 detected, the fast mode reports 98.64 percent. Sampling must be uniform over the same fault list (collapsed or not) that the full coverage uses, or the estimate is biased.

## Confidence

By the central limit theorem, $\hat{p}$ is approximately normal with standard deviation $\sqrt{p(1-p)/n}$, so the 95 percent confidence half-width is
$$h = 1.96\sqrt{\frac{\hat{p}(1-\hat{p})}{n}}.$$
At $\hat{p} = 0.9$ and $n = 1000$, $h = 1.86$ percentage points; at $n = 10{,}000$ and $\hat{p} = 0.986$, about 0.23 points. Coverage near 100 percent has small variance, which helps, but sign-off decisions are often about differences of a few hundredths of a point.

## When to recurse

To reach half-width $h$, the sample size must be
$$n = \left\lceil \frac{1.96^2 \, p(1-p)}{h^2} \right\rceil.$$
At $p = 0.942$ and $h = 0.05$ points (0.0005), that is about 840,000 faults, approaching the full list. So the agent uses sampling to discard clearly worse configurations, and recurses to full fault simulation when two candidates differ by less than about twice the half-width, or when the result feeds sign-off. Sampled estimates are also a reward-hacking risk if an agent can choose the sample, so the harness, not the agent, draws the seed.

## Going deeper

The agent-first EDA notes in this curriculum (fast and accurate engines); standard texts on survey sampling for the binomial confidence interval.
