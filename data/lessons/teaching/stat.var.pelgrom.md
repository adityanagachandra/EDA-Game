# Pelgrom mismatch scaling

The local random part of process variation has a remarkably simple law. Pelgrom, Duinmaijer and Welbers showed in 1989 that the mismatch between two identically drawn, closely spaced transistors shrinks with the square root of their gate area. That law sets the area of every matched analog pair, every sense amplifier and every SRAM bitcell, and it is the number a Monte Carlo model uses for each device's local sigma.

## The area law

Random mismatch comes from many small independent events under the gate: individual dopant atoms, roughness along the gate edge, grain boundaries. The parameter of one device averages those events over its area, so its random part has variance inversely proportional to area. For the threshold-voltage difference of a pair,
$$\sigma(\Delta V_T) = \frac{A_{V_T}}{\sqrt{W L}},$$
where $A_{V_T}$ is the Pelgrom coefficient of the process, in mV$\cdot\mu$m. With $A_{V_T} = 2.5$ mV$\cdot\mu$m and a device of $W = 1\ \mu$m, $L = 0.1\ \mu$m, the area is $0.1\ \mu$m$^2$ and $\sigma(\Delta V_T) = 2.5/0.316 = 7.9$ mV. A second term for distance appears in the original paper, describing gradients across the die, but for closely spaced pairs the area term dominates.

## Sizing for a target mismatch

Inverting the law gives the area a designer must pay: $WL = (A_{V_T}/\sigma)^2$. To bring a pair from 7.9 mV down to 4 mV with $A_{V_T} = 2.5$ mV$\cdot\mu$m needs $WL = (2.5/4)^2 = 0.39\ \mu$m$^2$, about four times the area. That quadratic cost is the central trade-off of analog layout: halving offset quadruples area, capacitance and power. It is also why SRAM bitcells, built from the smallest devices on the die, have the worst mismatch and the hardest statistical requirements. A sizing tool that minimizes area under a mismatch spec can use the law directly as a constraint.

## Pair sigma versus device sigma

The law is stated for the difference of two devices. A Monte Carlo model instead draws an independent local deviation for each device. The difference of two independent deviations with sigma $\sigma_1$ has variance $2\sigma_1^2$, so the per-device sigma is $\sigma_1 = \sigma(\Delta V_T)/\sqrt{2}$. Using the pair number per device doubles the simulated mismatch variance and overestimates every failure rate driven by it. The same factor appears in reverse when a test chip measures pairs and a modeling engineer extracts device statistics. Foundry mismatch models usually specify which convention their coefficients use; an engine that ingests them must convert explicitly.

## Going deeper

Pelgrom, Duinmaijer and Welbers, "Matching properties of MOS transistors," IEEE Journal of Solid-State Circuits, 1989; Kinget, "Device mismatch and tradeoffs in the design of analog circuits," IEEE JSSC, 2005.
