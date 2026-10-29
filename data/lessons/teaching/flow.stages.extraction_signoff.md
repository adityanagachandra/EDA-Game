# Extraction, timing and signoff

Signoff is the set of final checks that must all pass before a design is sent to manufacturing. Parasitic extraction computes the resistance and capacitance of every routed wire; static timing analysis uses them to check every path at every corner; power, IR-drop and electromigration analysis check the supply; design-rule and layout-versus-schematic checks verify the geometry. A signoff check is meant to be trusted, so signoff tools favour accuracy over speed.

## Parasitic extraction

Once wires exist, an extractor computes each net's resistance from wire length, width and sheet resistance, and its capacitance to ground and to neighbouring nets from the geometry and the dielectric stack. The result is an RC network per net, written as SPEF. Extraction is pattern-based for speed in most of the design and uses field solvers where accuracy matters. Coupling capacitance is the important new information: it makes a net's delay depend on whether its neighbours switch at the same time, which is crosstalk. Extracted coupling also feeds signal-integrity checks for glitches on quiet nets.

## Corners, modes and scenarios

A chip must work across process variation, supply voltage and temperature (PVT), and with interconnect at its own extremes (RC corners such as worst capacitance or worst RC). It must also work in every functional mode: normal operation, test mode, low-power modes. Signoff timing runs each combination as a scenario. With 3 modes, 4 PVT corners and 3 RC corners, there are $3 \cdot 4 \cdot 3 = 36$ scenarios. Tools prune dominated scenarios, but large designs still sign off dozens, which is why multi-scenario timing must be fast and incremental.

## Signoff slack with real parasitics

With extracted parasitics, setup slack at a flop is $T - (t_{cq} + t_{logic} + t_{wire}) - t_{setup}$, adjusted for skew and uncertainty. Wire delay estimated before routing is often too low. If a path with a 1000 ps clock had $t_{cq} = 60$, $t_{logic} = 700$, an estimated $t_{wire} = 120$ and $t_{setup} = 40$ ps, its slack was 80 ps. If extraction shows wire delay 50% higher, 180 ps, slack becomes 20 ps. Paths that fail at signoff are fixed with engineering change orders (ECOs): small, local edits such as resizing a gate or inserting a buffer, followed by incremental re-routing and re-extraction.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009); the OpenSTA and OpenRCX documentation.
