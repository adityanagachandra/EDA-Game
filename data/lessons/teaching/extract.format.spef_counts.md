# Reading a SPEF net section

Reading SPEF fluently means predicting its contents from a net's topology and computing timing quantities directly from its entries. Counting entries checks a writer or a reducer; totalling capacitance checks the header; and walking the resistor tree with the units in the header yields an Elmore delay. These are exactly the operations a parasitic reader and a delay calculator perform, and doing them by hand catches unit and topology mistakes that otherwise surface as baffling timing results.

## Counting entries from topology

A net whose resistor network is a tree on $p$ pins and $m$ internal nodes has $p + m$ nodes and therefore $p + m - 1$ resistors in its `*RES` section, since a tree has one fewer edge than nodes. If the extractor puts one ground capacitor on every internal node and none on pins, and the net has $c$ couplings, the `*CAP` section has $m + c$ entries. The `*CONN` section lists only the $p$ connection points. For one driver, one sink and five internal nodes, RES has 6 entries; with no couplings CAP has 5; CONN has 2. Meshes, created by multiple routes or loops, add resistors beyond the tree count, and a reader must handle them, because Elmore formulas assume a tree.

## Totals and units

Header units scale every value. With `*C_UNIT 1 FF`, a CAP section listing 1.3, 0.8 and 1.4 fF to ground and couplings of 0.3, 0.4 and 0.9 fF totals 5.1 fF. With `*C_UNIT 1 PF`, the same physical values appear as 0.0013, 0.0008 and so on, and the reader must multiply by 1000 to get femtofarads. Unit mistakes are a thousandfold error, large enough to be obvious in a delay, but a mix of files with different units in one flow can hide them in averages. A robust reader converts everything to internal units at parse time and never carries file units further, and a robust test compares the parsed totals against the header line for every net.

## Elmore delay from the entries

With `*R_UNIT 1 KOHM` and `*C_UNIT 1 FF`, resistance times capacitance comes out directly in picoseconds. For a chain driven by u1:Z through nodes n93:1, n93:2, n93:3 to sink u2:A, with resistors 0.13, 0.16, 0.03 and 0.07 kohm and capacitors 3.6, 0.8, 4 and 3.8 fF on the downstream nodes, the Elmore delay to the sink sums each resistor times all capacitance downstream of it: $0.13 \times 12.2 + 0.16 \times 8.6 + 0.03 \times 7.8 + 0.07 \times 3.8 = 3.462$ ps. Coupling capacitors are treated as grounded for this, the Miller factor 1 case. A delay calculator does this walk once per net and sink, in time linear in the number of resistors, by accumulating downstream capacitance in a reverse topological pass.

## Going deeper

IEEE Std 1481 (SPEF); Elmore, "The transient response of damped linear networks with particular regard to wideband amplifiers", Journal of Applied Physics (1948).
