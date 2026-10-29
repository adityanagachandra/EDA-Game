# Power planning: rings, stripes and rails

Every cell needs VDD and VSS, so the floorplan includes a power distribution network: rings around the core and macros, stripes on upper metal layers across the core, and rails along every standard-cell row. Its geometry consumes routing resource and limits where cells and wires can go, so the placer must know it, and its density sets the IR drop the design will see.

## Rings, stripes and rails

A power grid starts with rings of VDD and VSS around the core and often around each macro. Stripes, also called straps, run across the core on upper layers, in pairs of VDD and VSS at a regular pitch, alternating direction on successive layers to form a mesh. At the bottom, each row has rails on metal 1 along its top and bottom edges, which the cells touch by abutment; vias connect rails to the stripes above. The stripes are placed from an offset at a fixed pitch: a 1040 um core with a first stripe pair at 19 um and a 25 um pitch gets $\lfloor (1040 - 19)/25 \rfloor + 1 = 41$ pairs. OpenROAD's pdngen builds such grids from a specification of layers, widths, pitches and offsets.

## The routing cost of the grid

Power wiring occupies tracks that signal routing cannot use. If every 20 um of a layer carries a VDD and a VSS stripe each 1.5 um wide, 3 um of every 20 um, 15%, is gone. Dense grids reduce IR drop and electromigration risk but steal routing capacity, especially on the layers that carry long signal wires. Stripes also cut the space available near macros and may block pin access if placed carelessly, which is why power planning follows macro placement. The placer sees the grid as routing blockage in congestion estimates, and sometimes as placement blockage where vias stacks land on cell pins, so a placer that ignores it will underestimate congestion.

## IR drop along a rail

A row's rail is a resistor carrying current to the cells along it. With the rail fed at both ends by stripes, total current $I$ drawn uniformly along its length and rail resistance $R$ between the stripes, the voltage drop is largest at the midpoint, where it equals $IR/8$. A 110 um rail at 0.1 ohm/um has $R = 11$ ohm, and 5 mA gives $5 \times 11/8 = 6.875$ mV. Because $R$ grows with stripe pitch and the current fed by one segment also grows with pitch, the drop grows with the square of the pitch: halving the pitch quarters it. That quadratic law is how a planner sizes stripe pitch for a drop budget, before a full IR analysis confirms it. Dense logic placed in one region raises local current, so placement and power planning interact.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 3.8 on power and ground routing; the OpenROAD documentation for pdngen.
