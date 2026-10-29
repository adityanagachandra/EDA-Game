# Global versus detailed routing

Every production router splits the work in two. Global routing decides which coarse regions each net passes through, balancing demand against capacity over the whole chip. Detailed routing then draws exact wires and vias on tracks inside those regions, obeying every design rule. The split is what makes full-chip routing tractable, and the interface between the two stages is where many routability failures are born.

## What each stage decides

Global routing works on a coarse grid of gcells (global routing cells), each covering many tracks. Its graph has one vertex per gcell and an edge between adjacent gcells whose capacity is the number of usable tracks crossing their boundary. The output is, for each net, a set of gcells and layers: in OpenROAD these are written as route guides, rectangles per layer that bound where the net may go.

Detailed routing works on the full track grid inside those guides. It decides exact tracks, wire geometry, via positions and how each pin is reached, and it is responsible for DRC correctness: spacing, end-of-line, minimum area, via enclosure and the rest of the rule deck. Global routing reasons about resource allocation; detailed routing reasons about geometry.

## Why the coarse model can mislead

Gcell capacity is an estimate. It counts tracks crossing a boundary, but it cannot see that pins inside a gcell block some tracks, that a partial blockage cuts others, or that a detailed route needs local detours to reach a pin. So a global solution with zero overflow can still fail in detailed routing where pin density is high.

Flows defend against this in two ways. Layer adjustments shave a share of each layer's tracks off the capacity before global routing, leaving slack for detours and pin access. And the congestion map (overflow per edge) is read as a forecast: regions of high overflow are where detailed routing will struggle. Guides also trade flexibility for speed. A detailed router confined to its guides searches a small region and keeps global balance, but it cannot escape a guide that turns out to be locally blocked.

## Consequences for the flow

Placement estimates wirelength with half-perimeter bounds, a lower bound for each net in isolation. Global routing under capacity limits adds detours, so a large gap between placed and globally routed wirelength signals congestion. That is why routability-driven placement calls a fast global router or congestion estimate and spreads cells where it finds hot spots.

Runtime is dominated by detailed routing, which works on the fine grid and iterates search-and-repair until violations vanish. Gcell size is a real knob: halving it resolves local hot spots better but multiplies the vertex count by about four, and each edge then holds fewer tracks, so its capacity estimate becomes noisier.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); the OpenROAD documentation for the grt and drt modules.
