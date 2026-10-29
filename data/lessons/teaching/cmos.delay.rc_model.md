# RC delay models of a gate

The RC model replaces each conducting transistor by an effective resistance and each node by a capacitance to ground, so a gate's delay becomes an RC time constant. It is the model behind logical effort, early sizing in synthesis and placement, buffer insertion and every quick delay estimate a tool makes before it has characterized tables. A tool designer must know its numbers, its units and the factor that converts a time constant into a 50% delay.

## Effective resistance and capacitance per unit width

A transistor's drive current is proportional to its width, so its effective on-resistance scales as $1/W$. Take a unit nMOS of width 1 with resistance $R$. A pMOS of the same width is weaker because holes are slower, so its resistance is $\gamma R$ with $\gamma$ around 2 in many processes. Capacitance scales the other way: gate capacitance and drain diffusion capacitance are both proportional to width, often taken as roughly equal per unit width. A gate's output node therefore carries its own diffusion capacitance $C_{par}$ plus the load $C_{load}$ of the wire and the inputs it drives. Doubling a transistor's width halves its resistance and doubles its capacitance, which is the trade-off every sizing algorithm optimizes.

## The delay of a gate driving a load

With one conducting path of resistance $R_{eff}$ discharging or charging the output, the Elmore delay is
$$t = R_{eff}(C_{par} + C_{load}).$$
Units are convenient: one kilohm times one femtofarad is one picosecond. Example: a pull-down of 6 kOhm with 3 fF of parasitic and 17 fF of load gives $6 \times 20 = 120$ ps. Rising and falling edges use different paths: a falling output goes through the nMOS network and a rising one through the pMOS network. An inverter with an nMOS of width 1 and a pMOS of width 3, with $\gamma = 2.5$ and unit resistance 8 kOhm, has a pull-up of $2.5 \times 8 / 3 = 6.67$ kOhm, so a 40 fF total load rises in about 267 ps.

## From time constant to 50% delay

A single RC driven by a step reaches a fraction $1 - e^{-t/RC}$ of its swing at time $t$. The 50% point is at $t = RC \ln 2 \approx 0.69 RC$, and the 10%-to-90% rise time is $RC \ln 9 \approx 2.2 RC$. A model that reports $RC$ as the delay is therefore pessimistic for a single pole by about 1.45 times; in practice the effective resistance is fitted to measured delays so that the constant is absorbed into $R$. The model's real limits are elsewhere: it assumes a step input and a fixed resistance, while real inputs ramp and real transistors change regime during the transition. That is why the RC model is used for optimization and estimation, and characterized lookup tables for signoff.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), sections 4.3 and 4.4; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), section 5.4.
