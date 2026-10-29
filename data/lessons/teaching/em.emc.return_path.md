# Return current and the path of least impedance

Every signal current returns to its source. Which path the return takes depends on frequency: at DC it follows the least resistance, at high frequency the least inductance, which usually means directly under the signal. Return paths determine loop inductance, crosstalk and emissions, so layout checkers flag broken reference planes and power-integrity tools model where return current flows.

## Low frequency: the path of least resistance

At DC and low frequency, inductive reactance is negligible and parallel return paths share current inversely to their resistance, like any current divider. The fraction taking path 1 is
$$\frac{I_1}{I} = \frac{R_2}{R_1 + R_2}.$$

Example: a return can flow through a narrow strip of plane under the trace with 20 milliohm, or spread across a wide region of the plane with 5 milliohm. At DC only $5/25 = 20$% takes the strip under the trace; most of the current spreads out wherever the resistance is lowest. So at low frequency the return current does not follow the signal trace, and the loop it forms can be large. That is harmless for DC drop analysis but makes low-frequency magnetic coupling and ground loops possible.

## High frequency: the path of least inductance

At high frequency the impedance of each path is dominated by $j\omega L$, so current divides inversely to inductance (ignoring mutual coupling between the paths, which in reality strengthens the effect):
$$\frac{I_1}{I} = \frac{L_2}{L_1 + L_2}.$$
The lowest-inductance path is the one that minimizes the loop area, which for a trace over a plane is the strip of plane directly beneath it.

Example: if the path under the trace has 1 nH of loop inductance and the wide detour 9 nH, the path under the trace carries $9/10 = 90$% of the return current at high frequency. A slot in the plane under the trace removes that path and forces the current around the slot, enlarging the loop. Layout rule checks for plane splits, and power-integrity tools that model return current, exist because of this.

## The crossover frequency

A return path changes from resistive to inductive behavior where its reactance equals its resistance, $\omega L = R$:
$$f = \frac{R}{2\pi L}.$$
Below this frequency the current follows resistance; above it, inductance.

Example: a plane path with 10 milliohm and 2 nH crosses over at $0.01/(2\pi \cdot 2 \times 10^{-9}) = 796$ kHz. For typical board planes the crossover falls in the tens of kilohertz to around a megahertz, so for digital edges with gigahertz content the return current follows the trace almost exactly. This is why ground-plane integrity matters far more than DC ground resistance for signal integrity, and why a single analysis frequency is not enough for power-integrity models that span DC to gigahertz.

## Going deeper

Ott, "Electromagnetic Compatibility Engineering", chapter 3; Bogatin, "Signal and Power Integrity, Simplified" (3rd ed.), chapter 6; Paul, "Inductance: Loop and Partial", chapter 1.
