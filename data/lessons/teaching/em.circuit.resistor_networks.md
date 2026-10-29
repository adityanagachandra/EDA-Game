# Series, parallel and divider networks

Series and parallel reduction and the divider rules are KCL and KVL applied to the two simplest topologies. They are the mental shortcuts behind effective drive resistance in timing, resistor-network reduction in extraction, and every quick sanity check an engineer runs on a simulator's answer. They also expose the first real engineering effect: a load changes the very voltage it is connected to.

## Series and parallel reduction

Resistors in series carry the same current, so their voltages add and $R_{eq} = R_1 + R_2 + \cdots$. Resistors in parallel share the same voltage, so their currents add and conductances add: $1/R_{eq} = 1/R_1 + 1/R_2 + \cdots$. For two in parallel, $R_1 \parallel R_2 = R_1 R_2/(R_1 + R_2)$, always below the smaller one.

Example: 100 ohm in series with the parallel pair of 300 ohm and 600 ohm. The pair is $300 \cdot 600/900 = 200$ ohm, so the total is 300 ohm. Network reduction in parasitic extractors repeatedly applies these rules, plus more general eliminations, to shrink millions of resistor segments into a small model. The rules hold only for elements that genuinely share a current or a voltage; bridges and meshes need nodal analysis.

## Voltage and current dividers

Two resistors in series across a source split its voltage in proportion to resistance: $V_{out} = V_s R_2/(R_1 + R_2)$ across $R_2$. Two resistors in parallel fed by a current split it in inverse proportion: the current in $R_1$ is $I R_2/(R_1 + R_2)$, so the smaller resistor takes the larger share.

Example: a 1.8 V source across 200 ohm on top and 400 ohm on the bottom gives $1.8 \cdot 400/600 = 1.2$ V at the midpoint. A 30 mA current entering 100 ohm in parallel with 200 ohm sends $30 \cdot 200/300 = 20$ mA through the 100 ohm branch. Both formulas follow from one KVL or KCL equation, and both are the building blocks of reference generation, level shifting and the effective-resistance estimates timers use for drivers.

## Loading a divider

A divider's output formula assumes nothing draws current from the midpoint. Connect a load $R_L$ across the bottom resistor and the bottom becomes $R_2 \parallel R_L$, which is smaller, so the output drops: $V_{out} = V_s (R_2 \parallel R_L)/(R_1 + R_2 \parallel R_L)$.

Example: 1 kohm over 1 kohm from 2 V gives 1 V unloaded. A 1 kohm load makes the bottom 500 ohm and the output $2 \cdot 500/1500 = 0.667$ V, a third lower. The error is small only when $R_L \gg R_2$. The same effect appears whenever a model is characterized under one load and used under another: a cell's delay measured into a light load, a reference voltage that sags when a block draws current. Engines that cascade stages must therefore carry the load forward, not assume it.

## Going deeper

Agarwal and Lang, "Foundations of Analog and Digital Electronic Circuits", chapter 2; Horowitz and Hill, "The Art of Electronics" (3rd ed.), chapter 1.
