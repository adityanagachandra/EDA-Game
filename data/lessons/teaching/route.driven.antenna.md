# Antenna-aware routing

During fabrication, metal is patterned one layer at a time with plasma etching, which deposits charge on exposed conductors. A long wire connected to a transistor gate, with no path to a diffusion yet, can collect enough charge to damage the thin gate oxide. Antenna rules limit the ratio of such metal to gate area, and routers must obey them or insert fixes.

## The partial antenna ratio

When layer $M_i$ is etched, the pieces of $M_i$ that are connected to a gate but not yet to any diffusion collect charge. The partial antenna ratio for $M_i$ is the area of that metal divided by the total gate area it connects to. (Rule decks may use top area or sidewall area, depending on the layer.)

Example: a net driving three gates of width 0.3 um and length 0.1 um has gate area 0.09 um$^2$. A gate-connected metal-2 piece 180 um long and 0.05 um wide has area 9 um$^2$, so the partial ratio is 100. The connection to diffusion usually comes through a driver's output, which may be reached only on a higher layer, so lower-layer pieces are exposed while they are etched.

## The longest legal wire

A rule caps the partial ratio at $R_{max}$ on each layer. For a wire of width $w$ attached to gate area $A_g$, the longest legal gate-connected length on that layer is $L_{max} = R_{max} A_g / w$. Example: $R_{max} = 700$, three gates of 0.2 um by 0.05 um (0.03 um$^2$) and a 0.1 um wire give $L_{max} = 210$ um.

The router can avoid violations by construction: when a net's route from its driver to a gate is long on a lower layer, it can climb to a higher layer near the gate. That piece on the lower layer is then disconnected from the gate while the lower layer is etched, and is connected only when the upper layer is built. This layer hop is called a jumper. Alternatively a reverse-biased diode near the gate gives the charge a path to the substrate.

## Cumulative ratios

Some rules use a cumulative antenna ratio: the sum of the partial ratios of all layers up to the current one. This captures charge collected over several etch steps. Example: gate-connected pieces of 20, 75 and 65 um on metals 1 to 3, all 0.05 um wide, connected to 0.1 um$^2$ of gate, give a cumulative ratio of $8 / 0.1 = 80$ at metal 3.

Antenna-aware routing handles this during search, by pricing long lower-layer runs toward gates, and afterwards, with a repair step that inserts jumpers or diodes on violating nets. OpenROAD includes an antenna checker and can repair violations by inserting diodes. Fixes cost something: jumpers add vias, diodes add capacitance and area.

## Going deeper

The LEF reference on ANTENNA properties (ANTENNAAREARATIO and ANTENNACUMAREARATIO); the OpenROAD documentation for antenna checking and repair.
