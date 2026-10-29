# Stage delay: cell arc plus net arc

A timing path alternates cell arcs and net arcs. A stage is one driving cell plus the net it drives up to a sink pin, and its delay is what a delay calculator computes in one call. The cell's delay depends on the load the net presents, and the net's delay depends on its own resistance and the sink pin's capacitance, so the two arcs of a stage are computed together.

## Cell delay from drive resistance and load

A simple, widely used model treats the driving cell as an intrinsic delay plus a resistor charging its load: $d_{cell} = d_0 + R_{drv} C_{load}$. The load is everything the output must charge: the wire capacitance plus every sink pin's input capacitance. With $d_0 = 20$ ps, $R_{drv} = 0.5$ kohm, a wire of 8 fF and a 4 fF sink pin, $d_{cell} = 20 + 0.5 \times 12 = 26$ ps; one kilohm times one femtofarad is one picosecond. A Liberty table refines this with slew dependence and curvature, but the linear model captures the essential trade-off: a stronger cell, smaller $R_{drv}$, is faster for the same load but presents more input capacitance to its own driver.

## Adding the net arc

The net arc runs from the driver pin to the sink pin. Model the wire as a pi network: half the wire capacitance at each end and the wire resistance $R_w$ between. The near half is charged by the driver and belongs to the cell's load; the wire resistance charges the far half plus the sink pin, so its Elmore delay is $d_{net} = R_w (C_w/2 + C_{pin})$. With $R_w = 0.3$ kohm, $C_w = 6$ fF and $C_{pin} = 1$ fF, $d_{net} = 0.3 \times 4 = 1.2$ ps. If the cell delay is 37 ps, the stage takes 38.2 ps. Short nets have tiny net delays compared with cell delays; long nets, or nets on resistive lower metal layers, do not.

## Summing stages into a path arrival

A path's arrival is its start time plus the delays of its stages in order, each stage's sink pin being the next stage's input. With a start at 30 ps and two stages of 12.45 and 22.35 ps, the arrival is 64.8 ps. Every stage couples to its neighbors. Upsizing a gate lowers its own delay but raises the load of the stage before it; a long wire raises both the driver's delay and the net delay. Delay calculation also passes slew forward: each stage's output slew becomes the next stage's input slew. That is why timers evaluate stages in level order, and why an optimizer that changes one cell must update the delays of the stage that drives it as well as its own.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapters 4 and 5; Sutherland, Sproull and Harris, "Logical Effort" (1999), on the drive-resistance view of a gate.
