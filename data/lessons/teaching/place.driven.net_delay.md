# Net delay as a function of placement

Placement fixes the length of every wire before any wire exists, and wire length sets wire delay. Estimating a net's delay from its placement, with HPWL or a tree length feeding the Elmore model, is how a timing-driven placer sees the consequences of its moves. The key fact is that distributed wire delay grows with the square of length.

## Wire delay from placed pins

For a two-pin net, the routed length is estimated by the Manhattan distance between the pins, which equals its HPWL and its Steiner tree. With resistance $r$ and capacitance $c$ per unit length, a uniform RC line of length $L$ has Elmore delay $rcL^2/2$, half the value of a lumped model, because each bit of resistance charges only the capacitance beyond it. A driver at (0, 370) and a sink at (860, 860) give $L = 1350$ um; with $r = 0.5$ ohm/um and $c = 0.15$ fF/um the wire delay is $0.5 \cdot 0.15 \cdot 1350^2/2 = 68344$ ohm-fF, 68.34 ps. For multi-pin nets a Steiner or spanning tree from the placement gives the topology, and the Elmore delay to each sink follows from the tree.

## Adding the driver and the load

The full stage delay includes the driver's output resistance $R_d$ charging the whole wire and the sink pin capacitance $C_L$:
$$D = R_d(cL + C_L) + rL\left(\frac{cL}{2} + C_L\right).$$
The first term is linear in $L$, the second has a quadratic part. For $L = 580$ um, $r = 1$ ohm/um, $c = 0.15$ fF/um, $R_d = 700$ ohm and $C_L = 7$ fF: the driver term is $700 \cdot 94 = 65800$ ohm-fF and the wire term $580 \cdot 50.5 = 29290$ ohm-fF, 95.09 ps in total. Short nets are dominated by the driver term, so cell sizing matters more than placement for them; long nets are dominated by the wire term, where placement matters most.

## What this means for a placer

Because wire delay scales with $L^2$, stretching a net by a factor of 2.5, for example when a congestion-driven move pushes a sink away, multiplies its wire delay by 6.25. A placer that trades a little wirelength on critical nets for a lot of wirelength elsewhere can therefore improve timing even while total HPWL rises. The quadratic law also explains why very long nets are not fixed by placement at all but by buffering, which breaks them into segments whose delays add linearly. Timing-driven placers run a fast timer on these estimates: wire parasitics from placement, cell delays from Liberty, and the timing graph from STA, refreshed every few iterations.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8 on timing closure; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), chapter 4 on wires.
