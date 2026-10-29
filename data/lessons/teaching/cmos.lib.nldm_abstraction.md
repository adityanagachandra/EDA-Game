# What NLDM preserves and abstracts

NLDM reduces a cell to a few numbers per grid point: one delay and one transition for a ramp input into a lumped capacitor. That abstraction made static timing fast and practical, and every one of its approximations is a place where advanced timers add machinery. Knowing exactly what is preserved and what is lost tells a tool designer where errors come from.

## What NLDM keeps and drops

At each grid point NLDM preserves the delay and output transition measured for one stimulus: a ramp input of the given transition time driving a pure capacitance equal to the given load. It drops the input waveform's shape, keeping only one transition number, the output waveform's shape, keeping only delay and slew, and the load's structure, keeping only total capacitance. It also drops the current the cell draws over time and the way the cell's input capacitance changes while it switches. Within those assumptions the table is exact at its grid points and interpolated between them. Everything else is an approximation layered on top by the timer.

## Resistive loads and effective capacitance

Real loads are RC networks: a driver sees the near end of a resistive wire, and the far capacitance is shielded by the wire's resistance. Using total capacitance as the NLDM load overestimates delay at the driver, because the driver's output moves faster than it would into a lumped capacitor of that size. Timers fix this with an effective capacitance: an iteration that finds the single capacitance drawing the same charge as the RC load up to a threshold crossing, looks up the table with it, and then computes the wire delay separately, for example with Elmore or higher-order moment methods. Ceff is itself an approximation, and its error grows as wire resistance rises relative to the driver's resistance, as it does at advanced nodes.

## Waveforms, receivers and the range of validity

Real waveforms are not ramps: they have tails and plateaus, especially through resistive wires and with crosstalk. A single slew number cannot distinguish two waveforms with the same 20%-to-80% time but different shapes, yet they produce different delays in the next stage. The receiving pin's capacitance also varies during switching, through the Miller effect, which a fixed pin capacitance misses. And outside the table's index range the timer must extrapolate, which is why max_transition and max_capacitance limits exist. Current-source models (CCS and ECSM) address the waveform and receiver problems by storing currents or voltage waveforms, and noise views add what signal-integrity analysis needs. NLDM remains the fast default.

## Going deeper

Qian, Pullela and Pillage, "Modeling the effective capacitance for the RC interconnect of CMOS gates" (IEEE TCAD, 1994); Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapters 3 and 4.
