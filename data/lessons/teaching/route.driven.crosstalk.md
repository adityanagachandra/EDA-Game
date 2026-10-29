# Crosstalk-aware routing

Wires on adjacent tracks are coupled through their sidewall capacitance. When a neighbour switches, it injects charge into the victim, changing its delay and possibly causing a glitch. At advanced nodes coupling is a large share of total wire capacitance, so routers control it by spacing, shielding and ordering wires, and timers model its effect on delay.

## Coupling capacitance

To first order, the coupling capacitance between two parallel wires grows with their parallel run length and falls with their spacing: $C_c \approx c_0 L s_0 / s$, where $c_0$ is the coupling per micron at the minimum spacing $s_0$. Real extraction uses tables or field solvers, and the dependence on spacing is weaker than $1/s$ once fringing dominates, but the trend is what a router uses.

Example: $c_0 = 0.05$ fF/um over a 1100 um parallel run at minimum spacing gives 55 fF of coupling, likely more than the victim's ground capacitance. Long parallel runs on adjacent tracks are therefore the crosstalk hazard routers look for.

## Miller factors and the delay window

The effect of $C_c$ on the victim depends on what the aggressor does. If it is quiet, $C_c$ acts like a capacitor to ground (Miller factor 1). If it switches in the same direction, the voltage across $C_c$ does not change and it disappears (factor 0). If it switches in the opposite direction, the voltage change doubles (factor 2). So the victim's effective load ranges from $C_g$ to $C_g + 2C_c$.

With a simple driver model, delay = $R_d C_{eff}$, so the delay varies by $R_d \cdot 2C_c$ between the best and worst case. Example: 160 fF of coupling behind a 250 ohm driver gives a window of $250 \cdot 320$ fF = 80 ps. Timers apply these factors per aggressor, using timing windows to decide which aggressors can switch together.

## Spacing and shielding

A crosstalk-aware router reduces coupling on critical nets. Spacing a victim one extra track away roughly halves $C_c$ in the first-order model, cutting the worst-case delay by $R_d \cdot 2 \cdot \Delta C_c$. Example: moving from 3 to 6 times minimum spacing over 2000 um with $c_0 = 0.06$ fF/um halves $C_c$ from 40 to 20 fF, saving $600 \cdot 2 \cdot 20$ fF = 24 ps behind a 600 ohm driver.

Shielding goes further: a grounded wire on the neighbouring track never switches, so it fixes the Miller factor at 1 and blocks the aggressor entirely, at the cost of a track. Routers also limit parallel run length between known aggressor and victim pairs, order nets within a panel so that nets with overlapping switching windows are not neighbours, and use non-default rules with double spacing for clocks.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.), on timing closure.
