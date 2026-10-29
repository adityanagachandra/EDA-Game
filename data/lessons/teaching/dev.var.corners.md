# Process corners and statistical device models

A process corner is a single set of model parameters representing a deliberately extreme process, such as slow NMOS with slow PMOS. Corners are how most timing signoff still handles global variation, while statistical models and Monte Carlo handle the rest. Knowing what corners represent, how independent sources combine, and where corner methods are pessimistic or blind is essential for anyone writing a timer or a characterization flow.

## What corners represent

Foundries ship corner model cards named by NMOS and PMOS speed: TT (typical), SS, FF, and the skewed corners SF (slow NMOS, fast PMOS) and FS. Each is a global shift of device parameters, typically placed at a few sigma of the die-to-die distribution. A signoff corner combines a process corner with a supply voltage and temperature (PVT), plus interconnect corners for parasitics.

Skewed corners matter for circuits whose function depends on n-to-p strength ratios: SRAM write and read margins, level shifters, and ratioed logic.

Temperature does not monotonically slow devices. Mobility falls with temperature, slowing them, but threshold also falls, speeding them. At low supply voltage, where overdrive is small, the threshold effect can dominate, so delay decreases with temperature: temperature inversion. The slow corner at low voltage may then be the cold one, which is why signoff covers both temperature extremes.

## Combining independent sigmas

When independent Gaussian sources affect one parameter, their variances add:
$$\sigma_{total} = \sqrt{\sigma_{global}^2 + \sigma_{local}^2}.$$

Example: a threshold with nominal 400 mV, global sigma 20 mV and local sigma 15 mV. $\sigma_{total} = 25$ mV, so the $+3\sigma$ point for one device on a random die is 475 mV. Adding the sigmas linearly would give 505 mV, treating the sources as fully correlated and overstating the extreme.

The same arithmetic explains why corners are pessimistic for paths. A corner shifts global parameters to $3\sigma$; stacking worst-case local variation on every gate of a path at the same time adds a second set of correlated extremes that rarely co-occur. Statistical timing separates the two: corners for global variation, LVF-based POCV for local variation, with local sigmas combined in quadrature along the path.

## Statistical models and Monte Carlo

A statistical compact model carries distributions for selected parameters (threshold, mobility, effective length, oxide thickness), split into global and local parts with their correlations. Monte Carlo SPICE samples them to estimate the distribution of a circuit metric.

Library characterization uses Monte Carlo per arc, slew and load to produce LVF tables of delay sigma due to local variation. Analog designers use it to estimate offsets and yields.

Plain Monte Carlo has a sampling limit. A 6-sigma event has probability around $10^{-9}$, so a thousand samples with no failure say almost nothing beyond about 3 sigma. Rare-event methods (importance sampling, statistical blockade, scaled-sigma sampling and others) concentrate samples near the failure region, and SRAM yield analysis relies on them.

For a tool designer, Monte Carlo is also a cost problem: thousands of SPICE runs per cell arc. Surrogate models and sensitivity-based methods that predict sigma from a few runs are an active way to cut it.

## Going deeper

Orshansky, Nassif and Boning, "Design for Manufacturability and Statistical Design" (2008); the Liberty Variation Format (LVF) documentation in the Liberty reference manual.
