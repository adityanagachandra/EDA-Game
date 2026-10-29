# Propagation delay and electrical length

Signals on interconnect travel at a finite speed set by the dielectric. Propagation delay per unit length, the ratio of line delay to rise time, and electrical length in degrees decide whether an interconnect can be modeled as lumped, whether reflections matter, and how much phase a structure adds at a given frequency. These are the first checks a signal-integrity tool runs on each net.

## Delay per unit length

In a homogeneous dielectric the delay per unit length is $t_{pd} = \sqrt{\varepsilon_r}/c$. Since $1/c = 3.336$ ps/mm,
$$t_{pd} = 3.336\sqrt{\varepsilon_r}\ \text{ps/mm}.$$
For a microstrip, use $\varepsilon_{eff}$ instead.

Example: stripline in FR-4 ($\varepsilon_r = 4.3$) has $t_{pd} = 3.336 \cdot 2.074 = 6.92$ ps/mm, about 176 ps per inch; a 100 mm route has a 692 ps delay. In SiO2 ($\varepsilon_r = 3.9$) the speed of light limit is 6.59 ps/mm, but on-chip wires are much slower than this because their resistance makes them diffusive RC lines. Time-of-flight delay is the floor on interconnect delay, reached only by low-resistance lines such as package and board traces and wide top-metal on chip.

## When a wire is a transmission line

A line behaves like a lumped capacitance when reflections settle well within the signal's edge, and like a transmission line when its one-way delay $T_d$ is a significant fraction of the rise time $t_r$. Teams adopt a threshold rule, commonly $T_d > t_r/6$ (others use $t_r/4$ or $t_r/2$ for different accuracy); the exact fraction is a convention, so a tool should make it a parameter.

Example: with the one-sixth rule, a 60 ps rise time allows 10 ps of delay; in FR-4 at 6.92 ps/mm that is 1.45 mm. Below that length lumped C is adequate; above it, reflections and termination must be modeled. The criterion depends on the rise time, not the clock frequency, so a slow clock with fast drivers still needs transmission-line treatment. The same reasoning on chip uses the inductive regime check, comparing wire resistance with $2Z_0$.

## Electrical length in degrees

At a single frequency $f$, a line of delay $T_d$ shifts phase by $\theta = 360^\circ f T_d$, equivalently $\beta\ell$ with $\beta = 2\pi/\lambda$. One full wavelength is $360^\circ$ and a quarter wavelength is $90^\circ$.

Example: 50 mm of FR-4 stripline has $T_d = 346$ ps. At 2 GHz the period is 500 ps, so $\theta = 360 \cdot 2 \cdot 0.346 = 249^\circ$, about 0.69 wavelengths. Electrical length is what enters the input-impedance formula, stub resonances and quarter-wave transformers. A stub of $90^\circ$ at the Nyquist frequency of a link, such as an unused via barrel, reflects strongly and carves a notch in the channel response, which is why back-drilling via stubs matters at high data rates.

## Going deeper

Johnson and Graham, "High-Speed Digital Design: A Handbook of Black Magic", chapters 2 and 4; Bogatin, "Signal and Power Integrity, Simplified" (3rd ed.), chapter 7.
