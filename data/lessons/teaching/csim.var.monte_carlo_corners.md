# Corners and Monte Carlo in circuit simulation

Variation enters circuit simulation in two ways. Corner analysis runs the circuit at a few deliberately chosen parameter sets; Monte Carlo draws many random parameter sets from the statistical models and estimates distributions. Both multiply the number of simulations, so knowing exactly how their costs grow, and how their accuracy depends on the number of runs, is essential for planning verification and for building faster statistical methods.

## Counting corner runs

A corner is one combination of process corner (TT, FF, SS, FS, SF), supply voltage, temperature and parasitic-extraction corner. Corners multiply: 5 process corners, 3 supplies, 3 temperatures and 3 extraction corners give $5 \times 3 \times 3 \times 3 = 135$ runs per testbench, before Monte Carlo or multiple testbenches. That product is why teams prune corner lists per metric, choosing the corners known to bound each specification, and why simulators and farms schedule corner runs in parallel. Each run is an independent simulation, so corner analysis parallelizes perfectly, but its total cost grows multiplicatively with every new axis.

## Samples for a yield estimate

Monte Carlo estimates yield as the fraction of samples that pass. Each sample is a Bernoulli trial, so the estimate's standard error is $\sqrt{Y(1 - Y)/N}$. For a target standard error se, $N = Y(1 - Y)/\text{se}^2$. With $Y$ near 0.9 and se $= 0.01$, $N = 0.09/10^{-4} = 900$ samples. The required count does not depend on the number of random parameters, which is Monte Carlo's great strength in high dimension, but it grows as $1/\text{se}^2$. Verifying a very high yield, such as a failure rate of $10^{-6}$ per SRAM cell, would need millions of samples, which motivates importance sampling and other rare-event methods.

## Standard error of a simulated mean

For a continuous metric such as offset or delay, the sample mean of $N$ runs has standard error $\sigma/\sqrt{N}$, where $\sigma$ is the sample standard deviation. With $\sigma = 5$ mV over 100 runs, the mean is known to within 0.5 mV (one standard error). Quadrupling the runs halves the error. Note the difference between the spread of the metric, $\sigma$, which describes the circuit, and the standard error, which describes the estimate. A designer usually needs the spread, for example to check a $3\sigma$ limit, and the standard error tells how many runs make that estimate trustworthy.

## Going deeper

Rubinstein and Kroese, "Simulation and the Monte Carlo Method" (3rd ed., Wiley, 2016); Najm, "Circuit Simulation" (Wiley, 2010), chapter on statistical analysis.
