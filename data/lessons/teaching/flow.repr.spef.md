# SPEF parasitic exchange

The Standard Parasitic Exchange Format (SPEF), defined in IEEE 1481, carries the resistances and capacitances that extraction computes from routed layout. Static timing and power tools read it to compute wire delays, slews and crosstalk. A SPEF file is often the largest file in a signoff run, because it describes every routed net as an RC network, and its accuracy sets the accuracy of signoff timing.

## Header and units

A SPEF file begins with a header naming the design, the extraction tool and the units: `*T_UNIT 1 NS`, `*C_UNIT 1 FF`, `*R_UNIT 1 OHM`. All later values are multiples of those units, so a capacitance written as 0.0034 under `*C_UNIT 1 PF` is 3.4 fF. A name map follows, assigning short indices such as `*12` to long hierarchical net names to shrink the file. Getting the units wrong scales every parasitic by a factor of a thousand and is a classic integration bug. Readers should always convert to internal units immediately on parsing.

## A detailed net

Each net appears as a `*D_NET` block whose first line gives the net name and its total capacitance. Inside, `*CONN` lists the pins the net connects, `*CAP` lists capacitances, and `*RES` lists resistors between internal nodes. A `*CAP` entry with one node is a capacitance to ground; an entry with two nodes, one of them on another net, is a coupling capacitance. The total on the `*D_NET` line is the sum of the net's ground and coupling capacitances. For entries of 1.2, 0.8 and 2.5 fF to ground and 0.6 fF coupling, the total is 5.1 fF.

## Resistance paths and use in timing

The `*RES` entries form a resistor tree (or mesh) from the driver pin to each sink. The resistance from driver to a given sink is the sum along the unique tree path: segments of 12, 30 and 18 ohm give 60 ohm. A timer builds the RC network from these entries and computes delay with the Elmore approximation or a higher-order reduction. The coupling entries let it compute crosstalk delay, since the effective coupling capacitance depends on whether the neighbour switches in the same or the opposite direction. A reduced SPEF that lumps each net into one capacitance is smaller but loses exactly this information.

## Going deeper

IEEE Std 1481, whose SPEF chapter defines the format; the OpenSTA documentation for read_spef.
