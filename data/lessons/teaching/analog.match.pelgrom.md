# Mismatch and Pelgrom's law

Two transistors drawn identically side by side are never identical in silicon. Random dopant fluctuation, line-edge roughness and oxide variations give each device its own threshold voltage and current factor. Pelgrom's law says how large these local differences are and how they shrink with area. It sets the area, and therefore the capacitance and power, of every precision analog block, and it is the statistical model behind Monte Carlo mismatch simulation.

## Local mismatch shrinks with area

Random local variation comes from many independent microscopic sources spread over the gate area. Averaging over more of them reduces the variance in inverse proportion to area. Pelgrom's model states $\sigma(\Delta V_T) = A_{VT}/\sqrt{WL}$ for the threshold difference of a pair, with a process constant $A_{VT}$ typically between about 1 and 5 mV$\cdot\mu$m in modern processes. A similar law holds for the current factor, $\sigma(\Delta\beta/\beta) = A_\beta/\sqrt{WL}$. For $A_{VT} = 3$ mV$\cdot\mu$m and a device of $W = 4\ \mu$m by $L = 1\ \mu$m, $\sigma(\Delta V_T) = 3/2 = 1.5$ mV. The original model adds a distance term for long-range gradients, but at the close spacing of matched pairs the area term dominates, and gradients are handled by layout instead.

## Area for a target sigma

Inverting the law gives the area needed for a matching target: $WL = (A_{VT}/\sigma)^2$. For $A_{VT} = 3$ mV$\cdot\mu$m and a target of 0.5 mV, each device needs $36\ \mu\mathrm{m}^2$. The square is the important part: halving sigma quadruples area. Area brings capacitance, which lowers bandwidth or raises the power needed to drive it. This quadratic cost is the fundamental reason precision analog is expensive, and why designers prefer calibration, chopping or auto-zeroing when they need sub-millivolt offsets. A sizing tool must therefore include mismatch in its formulation: an optimizer that minimizes area without an offset constraint will shrink exactly the devices that set the offset.

## Using the statistics

Mismatch is a zero-mean Gaussian random variable for each device pair. A comparator offset spec of 5 mV with $\sigma = 1.5$ mV is a $3.3\sigma$ limit, met by about 99.9% of instances. Specs are written as a multiple of sigma for a yield target, typically 3 sigma for a block that appears once and more for one replicated thousands of times, such as a comparator in a flash ADC or a cell in an in-memory array. Because $\sigma \propto 1/\sqrt{WL}$, a pair's sigma at a new area follows directly: quadrupling the area halves sigma. Simulators implement the model by drawing a random threshold shift for each device with the Pelgrom sigma for its geometry, which is how Monte Carlo mismatch analysis works.

## Going deeper

Pelgrom, Duinmaijer and Welbers, "Matching properties of MOS transistors" (IEEE JSSC, 1989); Kinget, "Device mismatch and tradeoffs in the design of analog circuits" (IEEE JSSC, 2005).
