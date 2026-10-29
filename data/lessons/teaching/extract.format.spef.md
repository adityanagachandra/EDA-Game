# SPEF

The Standard Parasitic Exchange Format (SPEF, IEEE 1481) is how an extractor hands parasitics to timers, signal-integrity tools and power analysis. It is an electrical netlist of resistors, capacitors and optionally inductors per net, with a header of units and a name map. A tool designer must know exactly what it records, because anything it does not record, such as geometry, is lost to every downstream tool, and anything it records ambiguously becomes an integration bug.

## The structure of a D_NET

A SPEF file opens with a header (`*SPEF`, design name, divider and delimiter characters, and the units `*T_UNIT`, `*C_UNIT`, `*R_UNIT`, `*L_UNIT`), then an optional `*NAME_MAP` and `*PORTS`, then one section per net. A detailed net section starts with `*D_NET name total_cap`. Its `*CONN` part lists the connection points: ports as `*P` entries and instance pins as `*I` entries, each with a direction (I, O or B). Its `*CAP` part lists capacitors: an entry with one node is a capacitor to ground, an entry with two nodes, one of them on another net, is a coupling capacitor. Its `*RES` part lists resistors between nodes. Internal nodes are named net, delimiter, index, for example n17:3. The section ends with `*END`. Every value is a bare number scaled by the header units.

## The header total

The number on the `*D_NET` line is the net's total capacitance: the sum of every entry in its `*CAP` section, ground and coupling alike. If the CAP section of n55 holds ground capacitors of 2.6, 2 and 1.7 fF and couplings of 0.9, 0.1 and 0.4 fF, the header carries 7.7 fF. Tools use the header as a quick load value and as a consistency check: a mismatch with the sum of entries signals a truncated or corrupted file. Reduced forms exist too. A `*R_NET` section carries a reduced model, and a net may be given only as a lumped total. The more a net is reduced, the less a timer can compute about its wire delay and slew to each sink.

## What SPEF does and does not carry

SPEF is per corner: resistances and capacitances depend on process and temperature, so a flow writes one file per extraction corner. It carries electrical elements, not geometry: wire layers and shapes are gone, apart from optional node coordinates, so no tool can re-extract from it. It annotates by name, so the netlist the timer loads must match the names the extractor saw; after an ECO that renames hierarchy, unmatched nets silently fall back to estimates. Coupling capacitors may be written from either net's side or both, so a reader building a coupled network must instantiate each physical capacitor once. And a writer that grounds couplings to save space keeps totals right but disables crosstalk analysis. DSPF, a SPICE-like detailed format, serves transistor-level simulation where SPEF serves STA and SI.

## Going deeper

IEEE Std 1481, the standard for integrated circuit delay and power calculation, which defines SPEF; the OpenSTA documentation on reading parasitics.
