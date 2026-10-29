# Reflections at board impedance discontinuities

When a wave on a board trace meets a change of impedance, at a load, a via, a connector or a width change, part of it reflects. The reflection coefficient $\Gamma = (Z_L - Z_0)/(Z_L + Z_0)$ gives the fraction, and lattice bookkeeping of launched, reflected and transmitted waves predicts overshoot, ringing and settling. Board SI tools automate this with transmission-line simulation, but every termination decision rests on these three relations.

## The reflection coefficient

A wave of amplitude $V_i$ on a line of impedance $Z_0$ meets a load $Z_L$. Continuity of voltage and current at the junction gives a reflected wave $V_r = \Gamma V_i$ with $\Gamma = (Z_L - Z_0)/(Z_L + Z_0)$. Example: a 50 ohm trace into a 75 ohm load gives $\Gamma = 25/125 = 0.2$, a same-sign reflection; into 25 ohms, $\Gamma = -25/75 \approx -0.333$, an inverted one. An open end ($Z_L \to \infty$) gives $\Gamma = 1$, a short gives -1, and a matched load gives 0. CMOS inputs are nearly open, so an unterminated line reflects almost everything, and the reflection travels back to the source, where it reflects again with the source's own coefficient.

## Launch voltage and the first arrival

A driver with step $V_s$ and output resistance $R_s$ launches into a line as a voltage divider: $V_1 = V_s Z_0/(R_s + Z_0)$, because the line looks like $Z_0$ until any reflection returns. At the far end the voltage is the incident plus the reflected wave, $V_1(1 + \Gamma_L)$. Example: a 3.3 V driver with 20 ohms into a 50 ohm line launches $3.3 \cdot 50/70 = 2.357$ V; at a 1000 ohm load, $\Gamma_L = 0.905$, so the load jumps to about 4.49 V, an overshoot of 1.19 V above the rail. The reflection then returns to the source, where $\Gamma_S = (R_s - Z_0)/(R_s + Z_0) = -0.43$ inverts it, and the waveform rings down. A lattice diagram tracks these bounces until the net settles.

## Terminations and discontinuities

Terminations remove one of the two reflections. **Series** termination adds $R_t = Z_0 - R_s$ at the driver: the launched wave is half the swing, it doubles at the open receiver to the full swing, and the returning reflection is absorbed at the matched source. It suits point-to-point nets and costs no DC power. **Parallel** termination puts $Z_0$ at the far end, removing the load reflection but drawing DC current. Mid-line discontinuities, such as a connector, a via or a neck-down, both reflect and transmit: the transmitted fraction is $1 + \Gamma = 2 Z_2/(Z_1 + Z_2)$. A short discontinuity, electrically small compared with the rise time, acts as a lumped capacitance or inductance instead, which is why vias and pads matter more as edges get faster.

## Going deeper

Johnson and Graham, "High-Speed Digital Design" (1993), chapter 6, on terminations. Bogatin, "Signal and Power Integrity Simplified" (3rd ed., 2018), chapter 8.
