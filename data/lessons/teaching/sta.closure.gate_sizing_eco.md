# Gate sizing in a timing ECO

Upsizing a gate on a failing path is the most common timing fix, and the one most often over-estimated. A larger gate drives its load faster, but its larger input capacitance slows the gate that drives it. A correct ECO evaluation computes both effects; a tool that counts only the first will promise slack it cannot deliver.

## The resized gate gets faster

Model a gate's delay as $d = d_0 + R\,C_L$, with intrinsic delay $d_0$, drive resistance $R$ and load $C_L$. A cell $k$ times larger has $k$ times wider transistors: its drive resistance falls to $R/k$ and its input capacitance rises to $k\,C_{in}$, while $d_0$ stays roughly the same. With $d_0 = 12$ ps, $R = 3$ kohm and $C_L = 14$ fF, the original delay is 54 ps; tripling the size gives $12 + 1 \times 14 = 26$ ps. The gain is $R\,C_L (1 - 1/k)$, largest when the gate drives a heavy load. When the load is small, intrinsic delay dominates and upsizing buys little.

## The driver gets slower

The gate's input capacitance is part of its driver's load. Growing it from $C_{in}$ to $k C_{in}$ adds $(k-1)C_{in}$ to that load, so the driver's delay rises by $R_{drv}(k-1)C_{in}$. With $R_{drv} = 1$ kohm, $C_{in} = 2$ fF and $k = 4$, the driver slows by 6 ps. Side loads on the driver's net are slowed by its new output slew too, and other paths through the driver may lose slack. This is the effect logical effort formalizes: each stage's delay depends on the ratio of its load to its own input capacitance, so sizing one gate shifts effort to its neighbor.

## The net change in slack

The path's slack changes by the sum of both effects with opposite signs: $S' = S - (\Delta d_G + \Delta d_{drv})$, where $\Delta d_G$ is negative. With an old slack of -16 ps, a gate speed-up of 33.75 ps and a driver slow-down of 1.5 ps, the new slack is 16.25 ps. The ECO tool also checks that the larger cell fits the placement, that its leakage and area are acceptable, and that hold on short paths through it still passes, since a faster gate speeds every path through it. Iterative sizers apply this one gate at a time in order of benefit, re-timing incrementally after each change. Lagrangian relaxation formulations size many gates at once.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort" (1999); Chen, Chu and Wong, "Fast and exact simultaneous gate and wire sizing by Lagrangian relaxation" (IEEE TCAD, 1999).
