# Coupling capacitance and the Miller factor

Neighbouring wires couple through the capacitance between their facing sidewalls. Coupling capacitance changes delay, because a switching neighbour changes the voltage across it, and it injects noise into quiet or floating nets. Timers model the delay effect with Miller factors and noise analysis models the glitch, so extraction must report coupling separately from ground capacitance.

## Sidewall coupling between neighbours

Two parallel wires of thickness $T$ running side by side over length $\ell$ with spacing $S$ have, to first order, a parallel-plate coupling capacitance between their sidewalls,
$$C_c \approx \frac{\varepsilon_0\varepsilon_r T\ell}{S}.$$
This underestimates the true value, because fringing over the tops adds more, but it captures the scaling: coupling grows with aspect ratio $T/S$.

Example: $T = 0.2$ um, $S = 0.1$ um, $\ell = 500$ um, $\varepsilon_r = 2.7$: $C_c = 8.854 \times 10^{-3} \cdot 2.7 \cdot 0.2 \cdot 500/0.1 = 23.9$ fF. Modern minimum-pitch wires are taller than their spacing, so coupling to neighbours is often the largest single component of a wire's capacitance, exceeding the capacitance to the layers above and below. Spacing a wire out, or shielding it with grounded neighbours, are the standard fixes.

## The Miller factor

The charge a driver must deliver to a coupling capacitor depends on what the neighbour does. If the neighbour holds still, the voltage across $C_c$ changes by the victim's swing: factor 1. If it switches in the same direction at the same time, the voltage across $C_c$ does not change: factor 0. If it switches the opposite way, the change is two swings: factor 2. The effective load is
$$C_{eff} = C_g + k\,C_c, \quad k \in \{0, 1, 2\}.$$

Example: $C_g = 30$ fF and $C_c = 20$ fF give 30, 50 or 70 fF. Delay can therefore vary by more than a factor of two depending on neighbour activity. Signoff timers compute switching windows to decide which aggressors can align, rather than assuming the worst factor everywhere; the factor-of-two model itself is a simplification of slew-dependent coupling.

## Capacitive crosstalk on a floating victim

A dynamic node that is precharged and then left floating has no driver to hold it. When a coupled neighbour swings by $V$, charge conservation on the victim makes it follow the capacitive divider:
$$\Delta V = V\frac{C_c}{C_c + C_g}.$$

Example: a floating node with 15 fF to ground and 5 fF of coupling to a neighbour swinging 1 V jumps by $1000 \cdot 5/20 = 250$ mV, enough to flip a dynamic gate's evaluation. A driven victim recovers, so its glitch is smaller and decays with the driver's time constant, but the divider is the worst-case bound. Static noise analysis computes these glitches, propagates them through gates, and compares them with noise margins. It needs coupling capacitances per aggressor, which is why extractors keep coupled netlists rather than lumping everything to ground.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 9 on interconnect; Shepard, Narayanan and Rose, "Harmony: static noise analysis of deep submicron digital integrated circuits", IEEE Transactions on Computer-Aided Design, 1999.
