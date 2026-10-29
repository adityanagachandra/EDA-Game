# Design rules

Design rules are geometric constraints a layout must satisfy so that it can be manufactured with acceptable yield: minimum widths, spacings, enclosures, extensions, areas and densities. Each one encodes a limit of some process step. Routers, placers and cell generators must obey them by construction, and DRC engines check them, so their arithmetic is part of every physical-design algorithm.

## Width, space and pitch

The simplest rules are a minimum width per layer and a minimum space between shapes on that layer. Their sum is the tightest pitch: a layer with 18 nm minimum width and 22 nm minimum space supports a 40 nm line pitch. Width rules protect against opens and lithography failure of thin lines; space rules protect against shorts and bridges. Rules are checked on the polygons, with Manhattan geometry for most layers, so a DRC engine measures edge-to-edge distances between parallel edges and corner-to-corner distances where shapes approach diagonally. Routers turn the same numbers into a track grid.

## Enclosure and via landing

Enclosure rules require one layer to extend past another: metal must surround a via by a margin $e$ so that overlay error and size variation never leave the via partly outside the metal. A via of size $v$ in a line therefore forces line width at least $v + 2e$ across the line. If vias are 16 nm, enclosure is 4 nm and line spacing is 20 nm, a layer whose every line must accept a via has pitch $16 + 8 + 20 = 44$ nm, not the 36 nm that minimum width alone might suggest. Line-end enclosure along the line direction is often larger than side enclosure, because line ends pull back in printing.

## Width-dependent and conditional spacing

Wide metal needs more space: it dishes more in CMP, carries more current and prints with different bias. Spacing tables therefore select the required space by the width of the wider wire, and often by parallel run length. With a table of 20 nm up to 60 nm width, 28 nm up to 150 nm, and 40 nm above, a 100 nm wire beside a 20 nm wire must be 28 nm away edge to edge, so their centres are at least $50 + 28 + 10 = 88$ nm apart. Routers model such rules when they place wide power and clock wires, and many violations in practice come from conditional rules that a simple width-and-space model misses.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design", chapter on design rules and physical verification; the SKY130 PDK design-rule documentation as a public example.
