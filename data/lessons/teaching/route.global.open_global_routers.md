# FastRoute, NCTU-GR and OpenROAD grt

The research ideas of global routing, FLUTE trees, pattern and monotonic routing, negotiated congestion with history and bounded maze search, are combined in open global routers such as FastRoute and NCTU-GR, and in OpenROAD's grt. Knowing how they are composed, and what they hand to detailed routing, is what you need to read or modify these tools.

## The FastRoute pipeline

FastRoute (Pan and Chu, with later versions by Xu, Zhang and Chu) builds each net's topology with FLUTE, reshapes the tree for congestion, decomposes it into two-pin connections and routes them by L and Z patterns. Later versions added monotonic routing and multi-source multi-sink maze routing, virtual capacities that steer negotiation, and via-aware routing with layer assignment. Maze routing is applied only to connections that patterns leave overflowed.

OpenROAD's global router, grt, is based on FastRoute. It reads LEF and DEF, builds the gcell grid with layer adjustments, routes, and writes route guides: rectangles per net and layer that the detailed router, drt, must stay inside. It also reports the congestion map that placement and the designer use.

## Negotiation in practice

Contest-era routers such as NTHU-Route 2.0 and NCTU-GR built on negotiation: every iteration reroutes the nets on overflowed edges with a cost that combines present congestion and history, and searches are bounded to a box around the net so they stay fast and avoid long detours. NCTU-GR 2.0 also routed nets in parallel, scheduling together nets whose regions do not overlap so their demand updates cannot collide.

The ISPD 2007 and 2008 benchmarks scored total overflow first, then wirelength, which shaped these designs: they spend most of their time clearing the last few overflowed edges. A useful mental model is pattern routing for the many easy nets, negotiation with maze search for the few hard ones.

## Using a global router in a flow

Many global routers route in 2D, on the projection with capacities summed over layers, and assign layers afterwards, because negotiation converges much faster on the smaller graph. Layer assignment then has to fit each 2D route into the per-layer capacities while minimising vias, which can reintroduce overflow on individual layers.

A global router is also a fast analysis engine. Routability-driven placement calls it, or a cheaper estimate, to find hot spots; a truncated run with fewer iterations and pattern routing only produces a usable congestion map early in a flow. If the final result has a little overflow concentrated in a few gcells, the right move is usually local: spread placement there or adjust capacity, rather than hoping detailed routing absorbs it.

## Going deeper

M. Pan and C. Chu, "FastRoute" (2006) and "FastRoute 2.0" (2007); Y. Xu, Y. Zhang and C. Chu, "FastRoute 4.0: Global Router with Efficient Via Minimization" (2009); the OpenROAD grt documentation.
