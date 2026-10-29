# Switchbox routing

A switchbox is a routing region with fixed pins on all four sides, typically where channels meet. Unlike a channel, it cannot grow tracks to fit its nets, so routability is a real question. Switchbox routing introduced the iterative, rip-up-based style of local routing that modern detailed routers use inside each clip.

## Channels versus switchboxes

A channel has pins on its top and bottom edges and open left and right ends; its height can grow until every net fits, so a channel router always succeeds by adding tracks. A switchbox has pins on all four sides, and its dimensions are fixed by its neighbours. Both bounds still apply: the nets crossing any vertical or horizontal cut line cannot exceed the tracks crossing it.

Example: a switchbox with 6 tracks and 7 nets crossing one vertical cut is unroutable as given, and no detailed router can fix that. Building-block layouts produced switchboxes at channel intersections, which is where they got their name.

## Routing inside a fixed box

Because constraints come from all sides, a single greedy pass rarely succeeds, so switchbox routers are iterative. They route the most constrained nets first, such as nets running straight across from one side to the opposite, then fit the rest, and rip up and reroute when a net is blocked. In a two-layer gridded box with one layer per direction, every bend is a via.

Density bounds are necessary but not sufficient: a box can respect every cut and still be unroutable because of how pins interleave. That is why switchbox routers report failure rather than adding tracks, and why the surrounding flow must be ready to change the box or its neighbours.

## Switchboxes in modern routers

Over-the-cell routing on many metal layers removed explicit channels from standard-cell designs, but the switchbox formulation survives. A detailed router's local repair works on a clip with its boundary connections fixed, which is exactly a switchbox, and the panels of track assignment are channels in all but name.

So a designer reads switchbox ideas as the logic of local repair: freeze the boundary, reroute inside, and accept that some local knots require changing the boundary (a larger clip or a different global route). FPGA routing has its own switch boxes, a programmable interconnect structure; the name is shared but the problem differs.

## Going deeper

Sherwani, "Algorithms for VLSI Physical Design Automation", on switchbox routing; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
