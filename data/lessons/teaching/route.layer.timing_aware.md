# Timing-aware layer assignment

Upper metal layers are thicker and wider, so their resistance per micron is far lower. Putting a long, timing-critical net on them cuts its delay, but reaching them costs a via stack at each end, which adds resistance and uses scarce capacity. Timing-aware layer assignment decides which nets earn the upper layers.

## Delay on a low layer

Model the net as a driver of resistance $R_d$ driving a distributed wire of length $L$ with resistance $r$ and capacitance $c$ per micron into a load $C_L$. Its Elmore delay is
$$T = R_d(cL + C_L) + rL\left(\frac{cL}{2} + C_L\right).$$
The first term charges everything through the driver; the second is the wire's own distributed RC, which grows with $L^2$.

Example: $R_d = 550$ ohm, $c = 0.25$ fF/um, $C_L = 10$ fF, $r = 2$ ohm/um and $L = 200$ um give $550 \cdot 60 + 400 \cdot 35 = 47000$ ohm-fF, which is 47 ps. One ohm times one femtofarad is one femtosecond.

## Delay on an upper layer

Moving the wire to an upper layer replaces $r$ with a much smaller $r'$ but adds a via stack of $k$ vias of resistance $R_v$ at each end. With the start stack in series with the driver and the end stack in front of the load,
$$T' = (R_d + kR_v)(cL + C_L) + r'L\left(\frac{cL}{2} + C_L\right) + kR_v C_L.$$
Capacitance per micron is similar on both layers in this simple model.

Example: $R_d = 850$ ohm, two 30-ohm vias at each end, $c = 0.25$ fF/um, $C_L = 6$ fF, $L = 950$ um, $r' = 0.2$ ohm/um against $r = 3$ ohm/um. The upper layer gives about 246 ps against about 563 ps on the low layer: for long wires the $L^2$ term dominates and the upper layer wins by a wide margin.

## Break-even length

The upper layer saves $(r - r')L(cL/2 + C_L)$ and pays $kR_v(cL + 2C_L)$. The saving grows like $L^2$ and the penalty like $L$, so there is a single break-even length: shorter nets are faster on low layers, longer ones on high layers. Scanning lengths (or solving the quadratic) finds it; for typical numbers it is a few hundred microns.

Upper-layer capacity is scarce, so assigners rank nets by how much slack they need and how much an upper layer would buy, and give the top layers to the nets with the biggest gain per track. Short nets with small timing benefit and large via cost are pushed down. The same reasoning drives non-default rules: wider wires on critical nets lower $r$ at the price of tracks.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); S. Hu, Z. Li and C. J. Alpert, "A Polynomial Time Approximation Scheme for Timing Constrained Minimum Cost Layer Assignment" (2008).
