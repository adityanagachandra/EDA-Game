# Timing-driven routing

A timing-driven router blends delay and congestion in each node's cost, weighted by how critical the connection is. Critical connections pay mostly for delay and take fast routes even if congested; non-critical ones pay mostly for congestion and detour. VPR's router uses this cost with criticalities refreshed by static timing analysis between iterations.

## Criticality

After a timing analysis, the router assigns each source-sink connection $\mathrm{Crit}(i,j) = \max\left(0, \mathrm{MaxCrit} - \mathrm{slack}(i,j)/D_{max}\right)$, with $\mathrm{MaxCrit}$ just below 1, for example 0.99. A connection on the critical path gets 0.99, not 1, so a little congestion pressure always remains and even the most critical net can be pushed off an overused wire. With $D_{max} = 2850$ ps and slack 810 ps, $\mathrm{Crit} = 0.99 - 0.284 = 0.706$. VPR can also raise criticality to an exponent to sharpen the distinction. Criticalities are updated once per iteration from the current routing, so they follow the router as delays change.

## The node cost

The timing-driven cost of node $n$ for connection $(i, j)$ is $\mathrm{Crit}(i,j) \cdot \mathrm{delay}(n) + (1 - \mathrm{Crit}(i,j)) \cdot b(n) h(n) p(n)$, with delays normalized to the scale of the congestion term. With $\mathrm{Crit} = 0.9624$, a node delay of 2.25, $b = 1$, $h = 3$ and $p = 1$, the cost is $0.9624 \times 2.25 + 0.0376 \times 3 = 2.278$. The delay term is the Elmore delay the node adds to the partial route, which depends on the upstream resistance of the tree being built. The congestion term is PathFinder's, with history applied multiplicatively in VPR's form. Because both are priced in the same search, one wavefront handles both objectives.

## Choosing between fast and free paths

Given two paths, the router compares $\mathrm{Crit} \cdot \mathrm{delay} + (1 - \mathrm{Crit}) \cdot \mathrm{congestion}$ for each. A connection with slack 1230 ps out of $D_{max} = 1980$ ps has $\mathrm{Crit} = 0.369$. Path 1, delay 5 and congestion 19, costs $0.369 \times 5 + 0.631 \times 19 = 13.84$; path 2, delay 13 and congestion 5, costs 7.95, so the connection detours on path 2 and frees path 1 for critical nets. The effect over many nets is that slack is spent where it exists. The risk is instability: as routes change, criticalities change, and a connection can flip between paths, which is why criticalities are refreshed per iteration rather than per net.

## Going deeper

Betz, Rose and Marquardt (1999), chapter 6; Swartz, Betz and Rose, FPGA (1998); the VPR router documentation for max_criticality and criticality_exp.
