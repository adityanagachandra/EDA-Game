# Power distribution grids

The power distribution network (PDN) carries current from the package bumps or pads to every cell. On chip it is a hierarchy: rings around the core and macros, a mesh of wide straps on the upper metal layers, via stacks down through the stack, and the standard-cell rails on the lowest metal that each row's cells abut. OpenROAD builds this structure with its pdngen tool. Its design is a direct trade between voltage integrity and routing resource, since every strap is a track lost to signals.

## Rings, straps, rails and vias

A **ring** of VDD and VSS wires surrounds the core or a macro and collects current from the IO or bump connections. **Straps** cross the core at a regular pitch on upper layers, usually alternating direction on adjacent layers to form a mesh. **Rails** run along every standard-cell row on metal 1 or 2, shared by the cells above and below. **Via stacks** connect each strap crossing down to the rails. The resistance of a strap segment is $R = R_s L/W$, with $R_s$ the sheet resistance in ohms per square: a 240 µm, 3.5 µm wide strap at 0.03 $\Omega$/sq is 68.6 squares, or 2.06 ohms.

## The routing cost

Each strap pair occupies tracks on its layer. If a VDD and a VSS strap, each 2 µm wide, repeat every 30 µm, they take $4/30 = 13.3\%$ of that layer's width, and the vias they drop down block tracks on the layers in between. Raising strap width or density lowers IR drop but reduces routing capacity and can create congestion. PDN planning therefore chooses widths and pitches per layer: wide straps on thick top metals where resistance is low and routing is sparse, narrower ones lower down. In LEF and DEF the PDN appears as special nets, which the signal router must avoid.

## Parallel paths

Many straps share the current, so the grid's effective resistance falls with strap count. Ten straps of 10 ohms each, tied at their ends, give 1 ohm end to end. In a full mesh the straps are tied at every crossing, and the drop at a cell depends on its distance to the nearest supply connections through the whole network, which is why IR analysis solves the mesh as a linear system rather than by series and parallel reduction. Bumps spread over the die, as in flip-chip packages, shorten every path; peripheral pads force current to cross the whole die and give much larger drops at the centre.

## Going deeper

Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the chapter on power and ground routing; the OpenROAD documentation for pdngen.
