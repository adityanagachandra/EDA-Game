# And-inverter graphs

An and-inverter graph represents combinational logic with a single kind of node, the two-input AND, and inversion as an attribute of edges. It is the working representation of ABC and of most modern synthesis and verification engines, because one uniform node type makes hashing, simulation, cut enumeration and SAT encoding simple and fast.

## One node type with complemented edges

Every internal node computes $n = a \cdot b$ of its two fanins, and each edge may be complemented. Primary inputs are leaves, and outputs point at nodes (possibly complemented). Everything else follows from De Morgan: $a + b = \overline{\bar{a} \cdot \bar{b}}$ is one AND node with all three edges complemented, and NAND and NOR are one node each. Inverters cost nothing because they are bits on edges. A node is stored as two fanin literals, each a node index times two plus a complement bit, which is exactly the encoding of the AIGER file format. Converting a gate netlist is purely local, and the resulting AND count is the area measure ABC reports.

## AND counts of common gates

A two-input AND, OR, NAND or NOR is one node. A $k$-input AND or OR is $k - 1$ nodes. A two-input XOR needs three: $a \oplus b = \overline{\overline{a\bar{b}} \cdot \overline{\bar{a}b}}$, two ANDs for the minterms and one for their OR. XNOR is the same three with the output complemented. A multiplexer $s ? x : y = \overline{\overline{sx} \cdot \overline{\bar{s}y}}$ also needs three. A three-input XOR built as two XORs needs six. Converting the netlist $g_1 = XOR(b, c)$, $g_2 = XOR(g_1, c)$, $g_3 = NAND(a, g_2)$ without sharing gives $3 + 3 + 1 = 7$ AND nodes. XOR-heavy logic such as arithmetic and parity is expensive in AIGs, which motivates XOR-aware representations such as XAGs and MIGs.

## Levels as the delay measure

The level of a node is one more than the larger level of its fanins, with primary inputs at level 0, and the AIG depth is the largest output level. Inverters add nothing. A two-input AND, OR, NAND or NOR adds one level; an XOR, XNOR or MUX adds two, since all its inputs enter the first layer of ANDs. For $g_1 = AND(b, d)$, $g_2 = AND(g_1, d)$, $g_3 = MUX(c, g_2, a)$: $g_1$ is at 1, $g_2$ at 2, and $g_3$ at 4. AIG levels are a unit-delay model: they ignore fanout and the fact that a library gate may absorb several levels, but they correlate well enough with mapped delay to guide balancing and level-preserving rewrites.

## Going deeper

Kuehlmann, Paruthi, Krohm and Ganai, "Robust Boolean reasoning for equivalence checking and functional property verification", IEEE TCAD (2002); Biere, "The AIGER And-Inverter Graph (AIG) format"; the ABC documentation (Berkeley Logic Synthesis and Verification Group).
