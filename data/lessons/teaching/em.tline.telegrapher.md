# Transmission lines and the telegrapher's equations

When a wire's length is comparable to the distance a signal travels during its edge, the wire must be modeled as a transmission line: distributed series inductance and resistance, shunt capacitance and conductance per unit length. The telegrapher's equations govern voltage and current on it and give two numbers that summarize a lossless line, its characteristic impedance and its wave velocity. Every signal-integrity engine starts from them.

## The distributed RLGC model

Cut a uniform line into short sections of length $dx$: each has series $R'dx$ and $L'dx$ and shunt $G'dx$ and $C'dx$. KVL and KCL on one section, in the limit, give the telegrapher's equations
$$\frac{\partial v}{\partial x} = -R'i - L'\frac{\partial i}{\partial t}, \qquad \frac{\partial i}{\partial x} = -G'v - C'\frac{\partial v}{\partial t}.$$
For a lossless line ($R' = G' = 0$) they combine into the wave equation $\partial^2 v/\partial x^2 = L'C'\,\partial^2 v/\partial t^2$, whose solutions are waves travelling in each direction at speed $1/\sqrt{L'C'}$.

The one-way delay of a length $\ell$ is $\ell\sqrt{L'C'} = \sqrt{L_{tot}C_{tot}}$. Example: a 0.2 m line with a total of 66 nH and 26 pF has a delay of $\sqrt{66 \times 10^{-9} \cdot 26 \times 10^{-12}} = 1.31$ ns. The per-unit-length values, not the totals, define the line.

## Characteristic impedance and wave velocity

A wave travelling in one direction has a fixed ratio of voltage to current, the characteristic impedance; for a lossless line
$$Z_0 = \sqrt{L'/C'}, \qquad v = 1/\sqrt{L'C'}.$$
$Z_0$ is not a resistance that dissipates power; it is the load a long line presents to a driver before any reflection returns.

Example: $L' = 300$ nH/m and $C' = 120$ pF/m give $Z_0 = \sqrt{2500} = 50$ ohm and $v = 1/\sqrt{3.6 \times 10^{-17}} = 1.67 \times 10^8$ m/s, which is 167 mm/ns, a little over half the speed of light. Widening a trace raises $C'$ and lowers $L'$, so $Z_0$ falls; moving it farther from its reference plane does the opposite. Impedance control on boards and packages is the art of holding $L'/C'$ constant along a route.

## Effective dielectric constant

Since $v = c/\sqrt{\varepsilon_{eff}}$ and $v = 1/\sqrt{L'C'}$,
$$\varepsilon_{eff} = c^2 L'C'.$$
For a stripline fully embedded in one dielectric, $\varepsilon_{eff} = \varepsilon_r$. For a microstrip, part of the field is in air, so $\varepsilon_{eff}$ lies between 1 and $\varepsilon_r$.

Example: $L' = 300$ nH/m and $C' = 120$ pF/m give $\varepsilon_{eff} = (3.00 \times 10^8)^2 \cdot 3.6 \times 10^{-17} = 3.24$. A microstrip on FR-4 ($\varepsilon_r = 4.3$) often has $\varepsilon_{eff}$ near 3.2 to 3.4. An extracted line with $\varepsilon_{eff} < 1$ is unphysical and signals a solver or unit error, a cheap sanity check. Because microstrip's even and odd coupled modes see different mixtures of air and dielectric, they travel at different speeds, which causes far-end crosstalk.

## Going deeper

Ramo, Whinnery and Van Duzer, "Fields and Waves in Communication Electronics" (3rd ed.), chapter 5; Bogatin, "Signal and Power Integrity, Simplified" (3rd ed.), chapter 7.
