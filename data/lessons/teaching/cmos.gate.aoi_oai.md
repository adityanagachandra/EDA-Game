# AOI and OAI compound gates

Compound AND-OR-invert (AOI) and OR-AND-invert (OAI) cells implement a two-level function in one static stage. They are some of the most valuable cells in a library: technology mappers in synthesis tools such as ABC use them to collapse several gates into one, saving area, delay and power. A tool writer must know the naming convention cold, because cell matching, library parsing and equivalence checking all depend on reading a name or a Liberty function correctly.

## Reading the name

An AOI cell ANDs groups of inputs, ORs the group results, and inverts. An OAI cell ORs the groups, ANDs the results, and inverts. The digits after the letters give the group sizes, with 1 meaning a single input that enters the outer operator directly. With inputs lettered in group order:

- AOI21: $F = \overline{AB + C}$.
- AOI22: $F = \overline{AB + CD}$.
- AOI211: $F = \overline{AB + C + D}$.
- OAI21: $F = \overline{(A + B)C}$.
- OAI22: $F = \overline{(A + B)(C + D)}$.

The cells AO21 and OA21, without the I, are the non-inverting versions: the same stage followed by an inverter. A Liberty file states each cell's function explicitly, for example function : "!((A1&A2)|B)", so tools never have to trust the name; but humans and scripts do.

## Inputs, transistors and networks

Because AOI and OAI functions are single-stage static gates, their structure follows the general rules. The number of inputs is the sum of the digits, and the transistor count is twice that: an AOI221 has 5 inputs and 10 transistors. The networks come from the expression: in an AOI22 the pull-down network is two series pairs in parallel, P(S(A, B), S(C, D)), and the pull-up network is the dual, two parallel pairs in series, S(P(A, B), P(C, D)). An OAI is the mirror image: its pull-down has parallel groups in series. A useful symmetry: the pull-up of an AOI has the shape of the pull-down of the corresponding OAI, which is why both families exist with similar areas and why their pMOS stacks set their speed.

## Why mappers love compound cells

Consider $\overline{AB + CD}$. Built from simple gates it is two NAND2s feeding an AND2 (a NAND2 and an inverter): three stages and 14 transistors. As an AOI22 it is one stage of 8 transistors with no internal nodes to route. The single-stage version has a higher logical effort per input than a NAND2 (with a pMOS-to-nMOS ratio of 2, each AOI22 input has $g = 2$), but one stage of delay instead of two or three, and fewer wires for the placer and router. Technology mapping in ABC or Yosys enumerates cuts of the logic network and matches each cut's function against library cells, so compound cells give the mapper more ways to cover a function cheaply.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), section 9.2; Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting" and the ABC documentation on technology mapping.
