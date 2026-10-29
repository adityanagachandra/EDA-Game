# Board autorouters versus chip routers

Board and chip routers share algorithms (maze search, rip-up and reroute, negotiated congestion) but face different problems. A chip router handles millions of nets on a uniform track grid with preferred directions per layer. A board router handles thousands of nets with irregular pads, large obstacles, few layers, any-angle geometry, and electrical constraints such as differential pairs and length matching. That is why board routing is often interactive, with autorouting filling in the remainder.

## Shape-based routing on few layers

Board routers are mostly **shape-based** (gridless): they route against the exact polygons of pads, keep-outs and existing copper, with 45-degree or any-angle segments, rather than on a uniform grid of tracks. With 2 to a few dozen layers and big blockages such as BGAs and connectors, a strict horizontal-vertical layer split wastes too much capacity, so each layer routes in any direction. Routing starts with **fanout and escape**: bringing every BGA ball out to a via and a channel, which largely determines the layer count. The dominant interactive tool is the **push-and-shove** router, which moves existing traces and vias aside within clearance rules as the designer draws, keeping the board DRC-clean at every step.

## Shared algorithms

The algorithms carry over. A **maze router** expands a wavefront over a multilayer graph where vias are layer-change edges with an explicit cost, so it prefers fewer vias at equal length. Batch board autorouters iterate **rip-up and reroute** with costs that rise in contested regions, the same negotiation idea as chip global routing, over many passes until conflicts vanish or progress stalls. Power and ground mostly use planes, so the router only fans pins out to them. **Topological** routers first decide which side of each obstacle a wire passes (its homotopy class) and then compute geometry, spreading and smoothing wires afterwards. Boards reach external routers through interchange files such as Specctra DSN, which carries the outline, components, nets, rules and fixed wiring.

## Completion is not correctness

An autorouter's headline metric is completion: every net connected with clearance rules met. That says nothing about impedance (width and reference plane), length matching, return paths over plane splits, crosstalk between parallel runs, or via stubs. A differential pair routed as two independent nets takes different paths, losing coupling and adding skew. So designers route critical nets by hand first and **lock** them, define **net classes** with width, clearance, via and pair-gap rules, and then autoroute the rest. An AI or algorithmic router that is judged on completion alone will learn to trade the unmeasured constraints for completion. The fix is to measure them: post-route SI checks, length reports and constraint DRC become part of the router's acceptance test.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapters 5 and 6, for the routing algorithms. The Freerouting project documentation, for an open-source board autorouter and the Specctra DSN flow.
