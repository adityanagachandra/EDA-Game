# Metal-only ECOs with spare cells

A metal-only ECO is a small physical design problem with tight constraints: find spare cells of the right types close to where the fix belongs, wire them in with new metal, and change as few mask layers as possible. Each step is computable: nearest-spare search, wirelength of the new nets, and the mask count the change implies. Together they decide whether a late bug fix costs days or a respin.

## Finding the nearest spare

The fix needs particular cell types near a target point, typically the centroid of the pins it must connect. Among unused spare cells of each needed type, the router prefers the closest by Manhattan distance, $|x - x_0| + |y - y_0|$, since routing is rectilinear. If an ECO needs a NOR2 at (320, 270) um and NOR2 spares sit at (240, 160), (360, 140) and (500, 350), their distances are 190, 170 and 260 um, so the second is chosen. Spatial indexes, such as a k-d tree or a grid of bins per cell type, make this query fast on chips with many thousands of spares. A spare too far away adds wire delay that can defeat a timing fix, so spare density is a design-time trade between area and ECO reach.

## Wirelength of the new net

Once a spare is chosen, it must be wired to its sinks. Half-perimeter wirelength (HPWL), the sum of the bounding box's width and height over all pins, estimates the wire needed. A spare inverter at (360, 0) driving sinks at (110, 390), (110, 380) and (340, 260) has a bounding box 250 um wide and 390 um tall, HPWL 640 um. That wire adds resistance and capacitance, so the ECO's timing must be checked with the new net's extracted parasitics, not with the ideal delay of the inserted cell. Several fixes competing for the same spares become a small assignment problem: minimize total HPWL subject to each fix getting a cell of the right type, solvable as a bipartite matching.

## Masks changed

The cost of a metal-only ECO is the set of masks that must be remade. Every metal layer the change touches needs new masks, and so does every via layer between them. Multiple patterning multiplies this: a layer printed with double patterning has two masks. If an ECO rewires M1 through M5 with V1 through V4, and M1 to M4 are double patterned, the masks are $4 \times 2 + 1 = 9$ metal plus 4 via, 13 in total. Confining the fix to upper layers, for example using spares whose inputs are tied off on an upper layer so that reconnection only needs upper metal, cuts the mask count and cost sharply. This is why spare-cell design decides which layers carry the tie-offs.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapters on placement and routing; the OpenROAD documentation on its ECO and spare-cell commands.
