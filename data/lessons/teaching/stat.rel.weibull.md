# The Weibull distribution and lifetime

A chip that works at test can still fail in the field: oxides break down, wires electromigrate, transistors drift under bias. Reliability engineering models the time to failure as a random variable, and the Weibull distribution is its default. Its two parameters separate the time scale of a mechanism from the shape of its hazard, which is what tells an engineer whether failures are early defects, random events or wear-out.

## Survival and failure fraction

The Weibull distribution with scale $\eta$ and shape $\beta$ has reliability (survival) function and CDF
$$R(t) = e^{-(t/\eta)^\beta}, \qquad F(t) = 1 - R(t).$$
At $t = \eta$, $R = e^{-1} = 36.8$ percent for every $\beta$, which makes $\eta$ the characteristic life: 63.2 percent of the population has failed. With $\eta = 20$ years and $\beta = 2$, the fraction surviving 10 years is $e^{-0.25} = 77.9$ percent. The exponential distribution is the case $\beta = 1$. Weibull is also an extreme value distribution: it is the limiting law for the minimum of many independent positive lifetimes, which fits a chip that fails when its weakest site does.

## Time to a failure fraction

Reliability targets are stated as the time by which a small fraction may fail, for example 0.1 percent in 10 years. Solving $F(t) = q$ gives
$$t_q = \eta\,(-\ln(1 - q))^{1/\beta}.$$
With $\eta = 20$ years and $\beta = 2$, $t_{0.001} = 20\sqrt{0.0010005} = 0.63$ years. The early-failure fraction is extremely sensitive to $\beta$: with $\beta = 4$ the same $\eta$ gives $t_{0.001} = 20(0.0010005)^{1/4} = 3.56$ years. For small $q$, $t_q \approx \eta\,q^{1/\beta}$, so a lower shape parameter, meaning a wider spread of lifetimes, pulls the first failures much earlier. Product qualification is about this low tail, not the mean.

## Mean time to failure

The mean of a Weibull lifetime is
$$MTTF = \eta\,\Gamma\!\left(1 + \frac{1}{\beta}\right).$$
For $\beta = 2$, $\Gamma(1.5) = 0.886$, so with $\eta = 20$ years the MTTF is 17.7 years; for $\beta = 1$ the MTTF equals $\eta$; for large $\beta$ it approaches $\eta$ from below. The mean is a poor summary for qualification: two populations with the same MTTF can differ by orders of magnitude in their 0.1 percent life. Reliability engines therefore propagate the full distribution, or at least $(\eta, \beta)$, and report a percentile life at the use condition.

## Going deeper

Meeker and Escobar, "Statistical Methods for Reliability Data" (1998), chapter 4; JEDEC JEP122, "Failure Mechanisms and Models for Semiconductor Devices."
