# Power, performance and area

Power, performance and area (PPA) are the three numbers by which a chip implementation, and the tools that produced it, are judged. Performance is usually the clock frequency the design meets, power is the energy it burns per second, and area is the silicon it occupies, which sets cost. Every stage of the flow trades among them, and every EDA engine reports its quality of results (QoR) in these terms.

## Dynamic and static power

CMOS logic burns power when it switches and when it leaks. Dynamic power is
$$P_{dyn} = \alpha C V^2 f,$$
where $\alpha$ is the activity factor (the average fraction of the capacitance that switches per cycle), $C$ the switched capacitance, $V$ the supply voltage and $f$ the clock frequency. With $\alpha = 0.15$, $C = 500$ pF, $V = 0.8$ V and $f = 1$ GHz, $P_{dyn} = 0.15 \cdot 500 \times 10^{-12} \cdot 0.64 \cdot 10^9 = 0.048$ W, or 48 mW. Static (leakage) power flows even without switching and grows with temperature and with lower threshold voltages. Short-circuit power during transitions is usually a smaller third term.

## Voltage, frequency and the quadratic lever

Because $V$ enters squared, supply voltage is the strongest power lever. Lowering $V$ from 0.8 to 0.7 V at fixed frequency scales dynamic power by $(0.7/0.8)^2 = 0.766$, a 23% saving. But lower voltage makes gates slower, so the maximum frequency also falls. Dynamic voltage and frequency scaling (DVFS) exploits this at run time: when full speed is not needed, lowering $f$ and $V$ together reduces power roughly with the cube of the scaling. Implementation tools fight for the same lever indirectly: every picosecond of timing margin they recover can be spent on a lower voltage.

## Performance from the critical path

A synchronous design runs at most at the frequency set by its slowest register-to-register path. If that path has clock-to-output delay $t_{cq}$, logic and wire delay $t_{d}$ and setup time $t_{su}$, then $T_{min} = t_{cq} + t_d + t_{su}$ and $f_{max} = 1/T_{min}$. With $60 + 780 + 40 = 880$ ps, $f_{max} = 1/0.88\ \text{ns} \approx 1136.4$ MHz. Area then sets cost through dies per wafer and yield. The three pull against each other: upsizing gates raises performance and costs area and power; adding pipeline registers raises frequency and costs area, power and latency.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), chapter 5 on power; Weste and Harris, "CMOS VLSI Design" (4th ed.), chapter 5.
