# Transfer functions and the Laplace domain

The Laplace transform turns a linear circuit's differential equations into algebra in the complex frequency $s$. Each element becomes an impedance in $s$, and the ratio of output to input transforms is the transfer function $H(s)$. Model order reduction, AWE, Elmore delay, stability analysis and the initial and final value of any response all come from $H(s)$, which makes it the central object of interconnect analysis.

## Impedances in the s-domain

With zero initial conditions, a resistor's impedance is $R$, an inductor's is $sL$ and a capacitor's is $1/(sC)$. Circuits are then solved with the same divider, nodal and Thevenin rules as for resistors, but with rational functions of $s$. The DC gain is $H(0)$: capacitors become open circuits ($1/(sC) \to \infty$) and inductors become shorts ($sL \to 0$).

Example: $R_1 = 3$ kohm in series from input to output, and $R_2 = 1$ kohm in parallel with $C = 100$ fF to ground. At $s = 0$ the capacitor is open, so $H(0) = R_2/(R_1 + R_2) = 0.25$. Setting $s = j\omega$ recovers the phasor transfer function, so a Laplace-domain model gives both the transient and the frequency response.

## The pole of an RC divider

Continuing the network, the shunt arm is $Z_2 = R_2/(1 + sR_2C)$, and
$$H(s) = \frac{Z_2}{R_1 + Z_2} = \frac{R_2/(R_1 + R_2)}{1 + sC(R_1 \parallel R_2)}.$$
The single pole sits at $s = -1/\tau$ with $\tau = C(R_1 \parallel R_2)$: the capacitor sees the Thevenin resistance of the rest of the circuit.

Example: with $R_1 = 3$ kohm, $R_2 = 1$ kohm and $C = 100$ fF, $R_1 \parallel R_2 = 0.75$ kohm and $\tau = 75$ ps. The step response is $0.25V(1 - e^{-t/75\text{ ps}})$. Reading time constants off a transfer function's denominator is what a reduced-order model does: it keeps a few dominant poles and discards the rest. Here there is only one capacitor, so one pole is exact; a network with $n$ independent capacitors has $n$ poles.

## Initial and final values of a step response

For a step input $V/s$, the output transform is $Y(s) = H(s)V/s$. The final value theorem gives $y(\infty) = \lim_{s\to 0} sY(s) = H(0)V$ when the system is stable, and the initial value theorem gives $y(0^+) = \lim_{s\to\infty} sY(s) = H(\infty)V$.

Example: a probe attenuator has $R_1 \parallel C_1$ in the series arm and $R_2 \parallel C_2$ to ground. At high frequency the capacitors dominate, so $y(0^+) = V C_1/(C_1 + C_2)$; at DC the resistors dominate, so $y(\infty) = V R_2/(R_1 + R_2)$. With $C_1 = 10$ fF, $C_2 = 30$ fF, $R_1 = 3$ kohm and $R_2 = 1$ kohm, a 1 V step starts at 0.25 V and settles at 0.25 V: a flat response, because $R_1C_1 = R_2C_2$. Mismatched values give an overshooting or sagging edge, which is how a scope probe is compensated.

## Going deeper

Oppenheim and Willsky, "Signals and Systems" (2nd ed.), chapter 9; Pillage, Rohrer and Visweswariah, "Electronic Circuit and System Simulation Methods", chapters 2 and 3.
