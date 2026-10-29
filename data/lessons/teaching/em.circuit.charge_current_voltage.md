# Charge, current, voltage and power

Every number a circuit simulator, timer or power tool reports is built from four quantities: charge, current, voltage and power. Current is charge in motion, voltage is energy per unit charge, and power is their product. A tool designer meets them everywhere: switching energy in power analysis, charge delivered by a driver in timing, dissipation in electromigration and IR-drop checks. Getting the units and the integrals right is the first step to trusting any engine's output.

## Current as the rate of charge flow

Current is the rate at which charge crosses a surface, $i(t) = dq/dt$. Turning that around, the charge delivered between two times is the integral of current, $q = \int i\,dt$: the area under the current waveform. A constant current $I$ for a time $T$ delivers $IT$. A current that ramps linearly from 0 to $I$ over $T$ delivers half of that, the area of a triangle, $IT/2$.

Units scale conveniently on chips: one milliampere flowing for one nanosecond moves one picocoulomb. Example: a driver whose current ramps from 0 to 20 mA over 4 ns delivers $20 \cdot 4/2 = 40$ pC. If that charge lands on a 10 pF load it raises the load by $Q/C = 4$ V. A simulator integrates exactly this relation at every node and every time step.

## Voltage, energy and power

Voltage is energy per unit charge: moving charge $q$ through a potential difference $V$ exchanges energy $qV$. Power is the rate of that exchange, $p = dw/dt = v\,i$. For a steady supply voltage and current, energy over an interval is simply $E = VIt$.

Example: a block draws a steady 25 mA from a 1.2 V supply for 10 ns. Its power is $1.2 \cdot 25 = 30$ mW, and the energy taken is $30 \text{ mW} \times 10 \text{ ns} = 300$ pJ (one milliwatt for one nanosecond is one picojoule). Power-analysis engines sum exactly these products over every supply pin and time window. The sign convention matters: with current entering the positive terminal of an element, positive $vi$ means the element absorbs power, and negative means it delivers power, as a source does.

## Power dissipated in a resistor

For a resistor, Ohm's law $v = iR$ turns $p = vi$ into two equivalent forms: $p = i^2 R$ and $p = v^2/R$. Which one to use depends on what is known. If a fixed current flows through a resistor, dissipation grows with $R$; if a fixed voltage is applied across it, dissipation falls with $R$.

Example: 1.5 V across a 300 ohm resistor dissipates $1.5^2/300 = 7.5$ mW. A current of 10 mA through the same resistor dissipates $(0.01)^2 \cdot 300 = 30$ mW. On chip, $i^2R$ heating of wires drives self-heating and electromigration limits, and $v^2/R$ dissipation is the static cost of parallel terminations and resistive loads. A reliability checker therefore needs current densities, not just voltages, from the circuit engine.

## Going deeper

Alexander and Sadiku, "Fundamentals of Electric Circuits", chapter 1; Agarwal and Lang, "Foundations of Analog and Digital Electronic Circuits", chapter 1.
