# Timing-driven routing

For most nets the router minimises wirelength, but for timing-critical nets the shape of the routing tree decides delay. A Steiner tree that shares a long trunk between sinks is short in total but slows the critical sink with capacitance it does not need. Timing-driven routing evaluates topologies by delay, usually Elmore delay, and picks the one that serves the critical sink best.

## Elmore delay through a shared trunk

Model the net as a tree of distributed RC wires: resistance $r$ and capacitance $c$ per micron, a driver of resistance $R_d$ and pin capacitances at the sinks. The Elmore delay to a sink is $R_d C_{total}$ plus, for every wire segment on the path to the sink, $r l (c l / 2 + C_{below})$, where $C_{below}$ is all capacitance downstream of that segment.

In topology A a trunk of length $l_0$ runs to a Steiner point and branches of length $l_1$ and $l_2$ go to sinks 1 and 2. The trunk sees everything below it: $C_{below} = c(l_1 + l_2) + C_1 + C_2$. Example: $R_d = 350$ ohm, $r = 1$ ohm/um, $c = 0.2$ fF/um, $l_0 = l_1 = 150$ um, $l_2 = 250$ um, $C_1 = 10$ fF and $C_2 = 5$ fF give a delay to sink 1 of 64 ps.

## Delay with separate wires

In topology B the driver feeds two separate wires, of length $l_0 + l_1$ to sink 1 and $l_0 + l_2$ to sink 2. The path lengths are the same as in A, but now sink 1's wire carries only its own capacitance: the delay is $R_d C_{total} + r d_1 (c d_1 / 2 + C_1)$ with $d_1 = l_0 + l_1$.

B uses more wire, $l_0$ extra, so the driver term grows. Whether B is faster for sink 1 depends on the balance: a strong driver (small $R_d$) with resistive wire favours B, because isolating sink 1 from sink 2's branch matters more than the extra load on the driver; a weak driver favours A. Real routers compare such topologies with Elmore or a better delay model and keep the one that meets the required time.

## Choosing by criticality

Timing-driven routing chooses topology per net by slack. Nets with ample slack get the shortest tree, because wirelength is what loads the routing grid and the driver. Nets with negative or small slack get trees that shorten the path to their critical sinks, even at a wirelength cost: shorter source-sink paths, critical sinks moved closer to the driver in the tree, or wider and higher-layer wires.

Example: a router that picks the topology with the smaller delay to sink 1 might choose B at 56.9 ps over A at 63.7 ps. The same reasoning drives net weighting in timing-driven placement and the Prim-Dijkstra trade-off for tree construction; the router's job is to apply it where the timer says it matters.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.), on timing closure; J. Cong, A. B. Kahng and K.-S. Leung, "Efficient Algorithms for the Minimum Shortest Path Steiner Arborescence Problem" (1998).
