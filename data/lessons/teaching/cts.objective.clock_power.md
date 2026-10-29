# Clock network power

The clock is the busiest net on a chip: it toggles every cycle, reaches every flop and carries the longest wires. In many designs the clock network is a large share of dynamic power, often cited as a third or more. Every CTS decision, from topology to buffer count, is also a power decision, and clock gating exists because of this term. A tool designer needs the clock-power formula at their fingertips to judge any proposed network.

## Activity one

Dynamic power is $P = \alpha C V^2 f$, where $\alpha$ is the probability of a 0-to-1 transition per cycle, $C$ the switched capacitance, $V$ the supply and $f$ the clock frequency. A clock rises once every cycle, so $\alpha = 1$ and $P = C V^2 f$. The units work out neatly: picofarads times volts squared times gigahertz give milliwatts. A network of 200 pF at 0.8 V and 2 GHz burns $200 \cdot 0.64 \cdot 2 = 256$ mW. That is why clock capacitance is a first-class objective in CTS, next to skew.

## Where the capacitance lives

Clock capacitance has three parts: wire, sink pins and the buffers themselves. A tree with 120 mm of wire at 0.21 pF/mm, 21 000 flop pins of 0.8 fF and 330 buffers of 2 fF switched capacitance has $25.2 + 16.8 + 0.66 = 42.66$ pF. At 0.75 GHz and 0.95 V this is 28.9 mW. The breakdown tells the designer where to act: wire capacitance responds to topology and wirelength, pin capacitance only to gating, and buffer capacitance to sizing and stage count. A mesh, for example, attacks skew by adding wire and so pays mostly in the first term.

## The clock's share of total power

Data nets switch less often. With 240 pF of clock at activity 1 and 750 pF of data at average activity 0.12, the effective switched capacitances are 240 and 90 pF; the common factor $V^2 f$ cancels, so the clock takes $240/330 = 72.7\%$ of dynamic power. Shares like this explain the economics of clock gating: removing clock toggles from idle registers attacks the largest term. Activity factors for data come from simulation or from probabilistic propagation, which the power-analysis concepts cover.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits: A Design Perspective" (2nd ed., 2003), the chapter on power; Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter.
