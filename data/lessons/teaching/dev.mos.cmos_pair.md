# CMOS as a complementary device pair

CMOS pairs an NMOS pull-down network with a PMOS pull-up network so that, in any steady state, one of them is off. That complementarity gives near-zero static current, rail-to-rail outputs and robust noise margins, and it is why every standard cell, every timing arc and every power model in an EDA flow is built around it. Its costs are the PMOS's weaker mobility, short-circuit current during transitions and parasitic latch-up.

## Complementary networks

A static CMOS gate has a pull-up network of PMOS transistors between $V_{DD}$ and the output, and a pull-down network of NMOS transistors between the output and ground. The two are duals: series in one is parallel in the other. A NAND2 has two NMOS in series and two PMOS in parallel; a NOR2 the reverse.

For every input combination exactly one network conducts, so the output is driven firmly to a rail and no DC path exists from supply to ground apart from leakage. NMOS is used for pull-down because it passes a strong 0; PMOS for pull-up because it passes a strong 1. An NMOS pull-up would stop a threshold below $V_{DD}$.

The duality is what lets a cell generator or a switch-level simulator derive a gate's structure from its Boolean function, and what lets timing characterization reason about which transistors are on for each arc.

## Sizing the pair

Holes in bulk silicon have lower mobility than electrons, typically by a factor of 2 to 3, and the two devices may have different threshold magnitudes. For equal saturation current under the square law at full gate drive,
$$\mu_n W_n (V_{DD} - V_{tn})^2 = \mu_p W_p (V_{DD} - |V_{tp}|)^2 \;\Rightarrow\; \frac{W_p}{W_n} = \frac{\mu_n}{\mu_p}\left(\frac{V_{DD} - V_{tn}}{V_{DD} - |V_{tp}|}\right)^2.$$

Example: $\mu_n/\mu_p = 2$, $V_{DD} = 1.0$ V, $V_{tn} = 0.35$ V, $|V_{tp}| = 0.40$ V. $W_p/W_n = 2 \times (0.65/0.60)^2 = 2.35$.

Balanced sizing equalizes rise and fall delays and centres the switching threshold near $V_{DD}/2$. It costs PMOS area and input capacitance, so libraries often settle for partial balance. In FinFET processes, sidewall crystal orientation and strain narrow the n-to-p mobility gap, which changes these ratios and the cell architectures built on them.

## Energy, short-circuit current and latch-up

Charging a load $C$ from 0 to $V_{DD}$ draws charge $CV_{DD}$ from the supply at voltage $V_{DD}$, so the supply delivers $CV_{DD}^2$. Half is dissipated in the pull-up during charging and half is stored, then dissipated in the pull-down on discharge. Dynamic power is therefore $\alpha CV_{DD}^2 f$, where $\alpha$ is the switching activity.

During an input transition, while the input is between $V_{tn}$ and $V_{DD} - |V_{tp}|$, both networks conduct briefly. This short-circuit current grows with slow input slews, which is one reason timing signoff enforces maximum transition limits.

In bulk CMOS, the p+ source in an n-well, the well, the p-substrate and an n+ source form a pnpn thyristor. If triggered by a current spike, it can latch into a low-impedance supply-to-ground path. Well and substrate taps, guard rings and layout rules keep the parasitic bipolar gain below the latching threshold.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), chapters 1 and 5; Rabaey, Chandrakasan and Nikolić, "Digital Integrated Circuits" (2nd ed.), chapter 5.
