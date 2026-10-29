# Hazard rate and FIT

Reliability is quoted in rates: how often units fail per hour of operation among those still working. The hazard rate makes that precise, its shape over time is the bathtub curve, and the industry's unit for it is the FIT. Converting between Weibull parameters, test results, FIT and mean time to failure is daily work in qualification and in system-level reliability budgets.

## The hazard function

The hazard rate is the conditional failure rate of survivors, $h(t) = f(t)/R(t)$. For a Weibull lifetime,
$$h(t) = \frac{\beta}{\eta}\left(\frac{t}{\eta}\right)^{\beta - 1}.$$
$\beta < 1$ gives a falling hazard (infant mortality, weak units failing early), $\beta = 1$ a constant hazard (random failures, the memoryless exponential), and $\beta > 1$ a rising hazard (wear-out). The three phases together form the bathtub curve. With $\eta = 10^7$ hours, $\beta = 2$ and $t = 87{,}600$ hours (10 years), $h = 2 \times 10^{-7} \times 0.00876 = 1.75 \times 10^{-9}$ per hour. Burn-in exists to run the infant-mortality phase off before shipment.

## FIT from a life test

One FIT is one failure per $10^9$ device-hours. The hazard above is $1.75 \times 10^{-9}$ per hour, which is 1.75 FIT. From a test, the point estimate of a constant failure rate is failures divided by accumulated device-hours: if 500 devices run 1000 hours at a stress with acceleration factor 50 and 2 fail, the equivalent field exposure is $500 \times 1000 \times 50 = 2.5 \times 10^7$ device-hours, and the rate is $2/(2.5 \times 10^7) = 8 \times 10^{-8}$ per hour, or 80 FIT. Qualification reports an upper confidence bound, typically from the chi-square distribution with $2r + 2$ degrees of freedom for $r$ failures, which is always larger than the point estimate and is finite even with zero failures.

## From FIT to mean time to failure

For a constant failure rate $\lambda$, the exponential lifetime has $MTTF = 1/\lambda$. At 100 FIT, $MTTF = 10^9/100 = 10^7$ hours, about 1140 years. That does not mean a chip lasts a millennium: it means that in a population of 10,000 chips about one fails every 1000 hours. FIT rates of components add for a system in series, so a board with 200 components at 50 FIT each runs at 10,000 FIT, an MTTF of about 11 years. Converting a FIT budget to an MTTF is only meaningful in the constant-hazard phase; in wear-out, the hazard and hence the effective FIT grow with time.

## Going deeper

Meeker and Escobar, "Statistical Methods for Reliability Data" (1998), chapters 1 and 7; JEDEC JESD85, "Methods for Calculating Failure Rates in Units of FITs."
