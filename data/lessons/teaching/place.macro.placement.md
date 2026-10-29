# Macro placement

Macros (memories, register files, analog and IP blocks) are tens to thousands of times larger than standard cells, and their positions shape everything else: where connected logic can go, where routing channels form and how the power grid reaches them. Macro placement is usually done during floorplanning, by engineers, scripts or dedicated placers, and then frozen before standard-cell placement.

## What a good macro placement looks like

Good macro placements share a few traits. Macros sit near the core boundary or are clustered, leaving a large contiguous region for standard cells, rather than scattered across the middle where they would fragment rows and force detours. Pins face the logic they connect to. Macros that exchange wide buses sit close together, following the design's dataflow. Each macro carries a halo, a keep-out margin where no standard cell may sit, which leaves room for pin access and power connections. Narrow channels between macros are avoided, or blocked for placement, because cells placed there are hard to route to. Orientation matters as much as position: flipping a memory can turn its pin edge from the die boundary toward the core.

## Sizing the channel between macros

When macros must face each other, the channel between them must carry the nets that cross it. With $N$ nets crossing, $L$ routing layers in the crossing direction and track pitch $p$, the channel needs at least $\lceil N/L \rceil$ tracks, or $\lceil N/L \rceil p$ of width, plus a halo $h$ on each side: width $\ge \lceil N/L \rceil p + 2h$. For 220 nets, 3 layers, a 0.4 um pitch and 10 um halos, that is $74 \times 0.4 + 20 = 49.6$ um. Real estimates add margin for power stripes, blockages and imperfect track use. A channel that is too narrow shows up later as routing overflow, which the router cannot fix without moving the macros.

## How macro placement is automated and judged

Automated macro placers include annealing over block representations, analytical mixed-size placement that places macros and cells together and then legalizes and fixes the macros, and learning-based approaches. Whatever the method, macro placement is optimised through a cheap proxy (wirelength, density and congestion estimates), and its real quality appears only after standard cells are placed and routed. A placement with a good proxy score can still be congested if the proxy misses routing demand. Flows therefore evaluate candidate macro placements by running the rest of placement, at least approximately, and comparing final metrics. OpenROAD provides macro placement through engines such as its RTL-MP hierarchical macro placer, operating on the same database as global placement.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 3; the OpenROAD documentation for macro placement.
