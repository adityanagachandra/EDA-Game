# Dynamic power

Dynamic power is the energy spent charging and discharging capacitance as signals switch. In most active digital blocks it is the largest power component, and its formula $P = \alpha C V^2 f$ drives the main power-saving techniques: voltage scaling, clock gating and activity reduction. Power-analysis tools compute it net by net from capacitances and switching activity, so a tool designer must know exactly what each factor means.

## Energy per transition

Charging a capacitance $C$ from 0 to $V_{DD}$ through a pull-up draws charge $Q = C V_{DD}$ from the supply, so the supply delivers energy $Q V_{DD} = C V_{DD}^2$. Half of it, $\tfrac{1}{2} C V_{DD}^2$, ends up stored on the capacitor and half is dissipated in the pull-up resistance, whatever that resistance is. When the node later discharges, the stored half is dissipated in the pull-down, and the supply delivers nothing. So a full cycle, one rise and one fall, costs $C V_{DD}^2$ from the supply, all of it drawn at the rising edge. The result does not depend on the transistor sizes or the speed of the transition, only on $C$ and $V_{DD}$.

## The alpha C V squared f formula

Let $\alpha$ be the probability that a node makes a 0-to-1 transition in a clock cycle, $C$ its capacitance and $f$ the clock frequency. Its average power is
$$P = \alpha C V_{DD}^2 f.$$
For a whole block, sum over nodes, or use the total switched capacitance $\alpha C$. Example: 400 pF of total capacitance, $\alpha = 0.1$, $V_{DD} = 0.8$ V and $f = 2$ GHz give $0.1 \times 400\text{ pF} \times 0.64 \times 2\text{ GHz} = 51.2$ mW, since pF times GHz is $10^{-3}$. A clock net toggles every cycle, so its $\alpha$ is 1, which is why clock networks often consume a large share of dynamic power. Some texts define activity as toggles per cycle, counting both edges, and write $\tfrac{1}{2}\alpha C V^2 f$; the two conventions describe the same power, so always check which one a tool or library uses.

## Scaling levers

The formula shows where power can be cut. Voltage is squared, so it is the strongest lever: lowering $V_{DD}$ from 1.0 V to 0.8 V cuts dynamic power by 36% at the same frequency. Because a lower supply also slows gates, frequency usually drops with it, and if $f$ scales in proportion to $V$, power falls roughly as $V^3$: the same change then saves about 49%. Capacitance is reduced by smaller cells and shorter wires, which is why placement quality affects power. Activity is reduced by clock gating (stopping the clock to idle registers), operand isolation and better encoding. Dynamic voltage and frequency scaling exploits the $V^2 f$ dependence at run time. Each lever has a cost elsewhere: delay, area or design complexity.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), section 5.5; Chandrakasan, Sheng and Brodersen, "Low-power CMOS digital design" (IEEE JSSC, 1992).
