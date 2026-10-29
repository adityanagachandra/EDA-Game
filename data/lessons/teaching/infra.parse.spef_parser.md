# Parsing SPEF

SPEF (Standard Parasitic Exchange Format) carries the resistances and capacitances extracted from routed wires, net by net. Timers read it to compute wire delays, and it is often the largest file in the flow. A parser must handle the name map that compresses it, the unit header that gives its numbers meaning, and the per-net sections whose RC networks feed delay calculation.

## The name map compresses names

Hierarchical net names are long, and each appears many times in a SPEF: in the *D_NET line, in *CONN entries and in every node name of its RC network. The *NAME_MAP section assigns each name an index such as `*1234`, written once, and the body uses the index. Take the name u_core/u_alu/n_4521 (19 characters) appearing 12 times. Without the map that costs $12 \times 19 = 228$ bytes. With the map, references cost $12 \times 5 = 60$ bytes for the token *1234, and the map line costs $5 + 1 + 19 + 1 = 26$ bytes (index, space, name, newline), so the map saves $228 - 60 - 26 = 142$ bytes. A parser reads the map into a hash table first and resolves every token through it.

## Total capacitance of a D_NET

The header declares units: `*C_UNIT 1 PF`, `*R_UNIT 1 OHM`, `*T_UNIT 1 NS`. Each net section starts `*D_NET n17 <total cap>` followed by *CONN, *CAP and *RES. In *CAP, a line with one node is capacitance to ground, such as `1 n17:1 0.0012`; a line with two nodes is coupling capacitance, such as `3 n17:2 n940:1 0.0004`. The total on the *D_NET line includes both. With entries 0.0012, 0.0008 and coupling 0.0004 in pF, the total is $1.2 + 0.8 + 0.4 = 2.4$ fF. A coupling entry appears in each of the two nets' sections, so it counts toward each net's own total, and a crosstalk-aware timer keeps it as a coupling element rather than grounding it.

## Resistance along a parasitic path

The *RES section lists resistors as (index, node, node, value). Nodes are pins such as u1:ZN or internal points such as n17:3. The network is usually a tree rooted at the driver. With `*R_UNIT 1 KOHM`, the chain u1:ZN to n17:1 (0.012), n17:1 to n17:2 (0.030), n17:2 to u4:A (0.018), plus a branch n17:1 to u7:B (0.025), gives a driver-to-u4:A resistance of $12 + 30 + 18 = 60$ ohms; the branch to u7:B is not on that path. The parser builds an adjacency list over the nodes, then a delay calculator walks it, for example for Elmore delay. Lines can arrive in any order, so a reader must build the graph before walking it, and should flag loops or disconnected nodes.

## Going deeper

IEEE Std 1481 (the SPEF chapter); the OpenSTA source for its SPEF reader; the OpenRCX documentation in OpenROAD.
