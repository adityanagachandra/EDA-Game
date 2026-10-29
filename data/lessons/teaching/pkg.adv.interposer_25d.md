# 2.5D integration on interposers

In 2.5D integration several dies sit side by side on an interposer, a thin slab carrying fine wiring between them, which in turn sits on a package substrate. A silicon interposer uses back-end-of-line style wiring at about a micrometre pitch and through-silicon vias to reach the substrate. It gives die-to-die wiring thousands of times denser than a laminate, at the cost of a large silicon part that may exceed the lithography reticle.

## Interposer area and the reticle

An interposer must hold its dies, the gaps between them and an edge margin. For $k$ dies of width $W$ and height $H$ in a row with gap $g$ and margin $m$, the area is $(kW + (k-1)g + 2m)(H + 2m)$. A scanner exposes at most a 26 mm by 33 mm field, 858 mm^2. Two 25 by 30 mm dies with a 0.2 mm gap and 1 mm margin need $52.2 \times 32 = 1670.4$ mm^2, so at least $\lceil 1670.4/858 \rceil = 2$ reticle fields, stitched together. Stitching joins fields with matched wiring at the boundary; it is routine for interposers but adds design rules at the seams. Large interposers motivate the alternatives: organic interposers, and small silicon bridges.

## Wires across the die edge

The value of an interposer is wire density at the die edge. With wiring pitch $p$ and $N$ routing layers, the edge of length $E$ supports $N \lfloor E/p \rfloor$ wires. At 0.8 um pitch on 3 layers over 8 mm that is $3 \times 10000 = 30000$ wires. A laminate substrate at 20 um pitch on the same edge would carry 400 per layer. Interposer wires are short, usually a few millimetres, and are driven by simple unterminated drivers, which keeps energy per bit low. Their resistance, a few ohms per millimetre at sub-micrometre widths, limits reach to a few millimetres before repeaters or slower signalling are needed.

## The bump bottleneck

Wires on the interposer must land on the die through microbumps, and the bumps, not the wires, are usually the limit. If signal bumps occupy a band of depth $D$ along an edge of length $E$ at pitch $p_b$, the band holds $\lfloor E/p_b \rfloor \lfloor D/p_b \rfloor$ bumps. At 40 um pitch over 8 mm and a 1 mm band, that is $200 \times 25 = 5000$, far fewer than the 30000 wires the interposer could carry. The usable count is the minimum of the two. Raising bandwidth therefore means shrinking bump pitch, deepening the band at the cost of die area, or raising the data rate per wire. This is the motivation for hybrid bonding and for interface standards that fix bump maps.

## Going deeper

Lau, "Semiconductor Advanced Packaging" (2021), chapters on 2.5D; Kahng and others on interposer design in the IEEE Heterogeneous Integration Roadmap.
