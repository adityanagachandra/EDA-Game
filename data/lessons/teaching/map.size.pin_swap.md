# Pin swapping on symmetric inputs

Inputs of a NAND or NOR are logically interchangeable, but not electrically: the transistor nearest the output switches fastest. Pin swapping reassigns signals among symmetric pins so that late signals use fast pins. It changes neither function nor area, so it is one of the cheapest timing fixes after mapping, and the matcher's symmetry information makes it automatic.

## Why symmetric pins differ

In a NAND3 the three nMOS transistors are in series. A signal on the transistor nearest the output switches when the internal nodes below it are already discharged, so it sees the least capacitance to discharge and has the smallest delay; the transistor nearest ground is slowest. Liberty records this as different timing arcs per pin. A NAND3 might show pin A at 15 ps, B at 19 ps and C at 20 ps. The output settles at $\max_i (a_i + d_{\text{pin}(i)})$, so the assignment of signals to pins matters whenever arrivals differ.

## The optimal assignment

To minimize $\max_i (a_i + d_{\sigma(i)})$, sort signals by arrival and pins by delay, and give the latest signal the fastest pin, the next latest the next fastest, and so on. An exchange argument proves it: if a later signal sits on a slower pin than an earlier one, swapping them cannot increase the maximum. Example: arrivals s = 12, t = 36, u = 36 ps and pins A = 12, B = 15, C = 18 ps. Assign t to A, u to B, s to C: max(48, 51, 30) = 51 ps.

## Gains and where symmetry comes from

The saving is the current settle time minus the sorted one: with arrivals 16, 39, 6 ps on pins of 14, 15, 16 ps in that order, the current output time is 54 ps and the sorted one 53 ps. Gains are small per gate but free, and they add along a path. Pins are swappable only within a symmetry group: all inputs of a NAND3, but within an AOI22 single pins swap only inside {A1, A2} or inside {B1, B2}, while the two pairs can be exchanged as wholes. Boolean matching computes these groups as a side effect, since symmetric variables are exactly those whose swap leaves the cut function unchanged, and the same information lets the mapper pick the best transformation at match time.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed., 2010), chapter on series transistor ordering.
