# Device reduction and property comparison

Layouts draw devices differently from schematics: one schematic transistor becomes four parallel fingers, a resistor becomes a serpentine of segments. Before comparison, LVS reduces the layout netlist to logical devices by merging parallel (and optionally series) devices, then compares device properties against tolerances. These steps decide whether a correct layout matches at all.

## Parallel merging

Devices are merged in parallel when they have the same type, the same gate net, the same pair of source and drain nets in either order, and the same L. The merged device's W is the sum. A list of nine layout transistors might reduce to four logical devices: three groups of fingers plus one device that shares a gate but has a different L and stays separate. Merging is grouping by a key (type, gate, unordered source/drain pair, L), which is a hash-table pass over the devices. Without it, a four-finger layout would never match a single schematic device.

## The merged device's width

Once grouped, the merged device carries $W = \sum_i W_i$ over its fingers, and its L is the common L. Fingers of 200, 310 and 150 nm give W = 660 nm. The schematic often states $W$ per finger and a finger count instead; the tool must compare like with like, either total W or W per finger and count, as the deck's property rules say. Source and drain interchangeability matters: a finger with source and drain swapped is still parallel.

## Resistors and property tolerance

A layout resistor's value comes from its geometry: $R = R_s \sum_i L_i / W_i$, the sheet resistance times the number of squares, plus contact and bend corrections the deck defines. A poly resistor with $R_s = 200$ ohm per square drawn as segments 12 by 1 and 8 by 0.5 um has $12 + 16 = 28$ squares, so 5600 ohm. Against a schematic value of 5400 ohm, the difference is 3.7 percent. LVS reports a property error when the difference exceeds the deck's tolerance. Property checks run after topology matches, because only matched devices can be compared.
