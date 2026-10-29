# Reflections and impedance matching

A wave on a line of impedance $Z_0$ that meets a different impedance partly reflects. The reflection coefficient sets how much, with what sign, and determines the voltage the load actually sees. Matching the load to the line eliminates reflections, and return loss and VSWR are the standard ways tools and datasheets report how good a match is.

## The reflection coefficient

At a load $Z_L$ terminating a line $Z_0$, the reflected wave is $\Gamma$ times the incident wave, with
$$\Gamma = \frac{Z_L - Z_0}{Z_L + Z_0}.$$
It ranges from -1 (short) through 0 (matched) to $+1$ (open) for resistive loads, and is complex for reactive ones.

Example: a 50 ohm line into 100 ohm gives $\Gamma = 50/150 = 0.333$; into 25 ohm it gives $-25/75 = -0.333$. The same formula applies at the source end with $R_s$ in place of $Z_L$, and at any discontinuity along a route, such as a narrow neck or a via, using the impedances on either side. A signal-integrity checker computes $\Gamma$ at every impedance change along a net to locate the worst discontinuities.

## The voltage at the load

A source of voltage $V_s$ and resistance $R_s$ driving a long line sees only $Z_0$ until a reflection returns, so it launches $V_i = V_s Z_0/(R_s + Z_0)$. When that wave reaches the load, the load voltage is the incident plus the reflected wave:
$$V_L = V_i(1 + \Gamma_L).$$

Example: a 1 V step with 25 ohm source resistance on a 50 ohm line launches 0.667 V. Into a 100 ohm load, $\Gamma_L = 0.333$ and the load sees $0.667 \cdot 1.333 = 0.889$ V on the first arrival. Into an open, $\Gamma_L = 1$ and the load sees 1.333 V, an overshoot above the source voltage. The current changes by $(1 - \Gamma_L)$ times the incident current, so a short doubles the current while the voltage collapses.

## Return loss and VSWR

In sinusoidal steady state the incident and reflected waves form a standing wave whose maximum and minimum amplitudes have the ratio
$$\text{VSWR} = \frac{1 + |\Gamma|}{1 - |\Gamma|}.$$
Return loss is the reflected power ratio expressed positively, $RL = -20\log_{10}|\Gamma|$ dB.

Example: a 50 ohm line into 100 ohm has $|\Gamma| = 1/3$, a VSWR of 2.0 and a return loss of 9.54 dB. A perfect match has infinite return loss and a VSWR of 1. A 20 dB return loss means $|\Gamma| = 0.1$: one percent of the power is reflected. For resistive loads the VSWR is simply the ratio of the larger impedance to the smaller one. Specifications for connectors and packages are written in return loss, and $|S_{11}|$ in dB is its negative.

## Going deeper

Pozar, "Microwave Engineering" (4th ed.), chapter 2; Johnson and Graham, "High-Speed Digital Design: A Handbook of Black Magic", chapter 4.
