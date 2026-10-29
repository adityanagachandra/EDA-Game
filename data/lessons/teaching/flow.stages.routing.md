# Global and detailed routing

Routing turns every net's connectivity into metal wires and vias on the available layers, without breaking design rules and without shorting nets together. It runs in two phases. Global routing assigns each net to a coarse path through a grid of routing regions, managing congestion. Detailed routing then draws exact tracks and vias inside each region. Routing failures are where earlier placement and floorplan decisions come due.

## Global routing and overflow

The global router divides the chip into a grid of rectangular regions called gcells. Each boundary between two adjacent gcells has a capacity: the number of tracks that cross it. Capacity per layer is the boundary length divided by the track pitch, and the layer's preferred direction decides which boundaries it serves. A gcell edge 10 $\mu m$ long with a pitch of 0.2 $\mu m$ has 50 tracks per layer, so three layers in that direction give 150. Each net's path adds one unit of demand to every boundary it crosses. Overflow on an edge is $\max(0, \text{demand} - \text{capacity})$, and the global router rips up and reroutes nets until total overflow is zero or as small as possible.

## Detailed routing and design rules

Detailed routing assigns each net segment to a specific track and layer and inserts vias between layers, obeying spacing, width, via enclosure and many more rules from the process. Modern processes have complex rules (minimum area, end-of-line spacing, multiple-patterning colouring), so detailed routers spend most of their effort on repairing design-rule violations. A design that routes with zero overflow can still fail detailed routing if local pin access is too dense. TritonRoute in OpenROAD is an open-source detailed router. Routing runtime is often dominated by the last few hundred violations, which need repeated local rip-up and reroute.

## Detours and what routing reveals

A routed net is usually longer than its HPWL because it detours around congestion and blockages. If a net with HPWL 120 $\mu m$ routes at 138 $\mu m$, it is 15% longer than the placement estimate, and its delay and coupling grow accordingly. Routing is where real wire geometry first exists, so it is also where crosstalk, antenna effects and actual wire delay appear. A design that met timing with placement estimates can fail after routing, and severe congestion may force a return to placement or even the floorplan.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapters 5 and 6; the OpenROAD FastRoute and TritonRoute documentation.
