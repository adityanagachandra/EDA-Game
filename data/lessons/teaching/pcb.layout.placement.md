# Board placement and ratsnest length

Board placement positions a few hundred to a few thousand heterogeneous packages under mechanical, thermal, electrical and assembly constraints. Unlike chip placement, it is mostly interactive, guided by grouping from the schematic and by a live ratsnest: each net drawn as a minimum spanning tree over its pins. Total ratsnest length is the board placer's wirelength proxy. Knowing what it measures, how it is computed and what it misses is the core of placement quality.

## What drives board placement

Placement starts with **fixed** parts: connectors, mounting holes, switches, displays and anything else the enclosure dictates. Then come the parts that set the board's critical structure: the processor or FPGA, memory placed to keep its bus short, the power stages near the inputs they feed, and the clock sources. **Decoupling capacitors** go as close as possible to the power pins they serve, because their high-frequency benefit is limited by the loop inductance of the connection. Cross-probing from the schematic gathers each functional block so its local nets stay short. Courtyard checks run continuously, since courtyards include lands and an assembly margin. The placer is the designer, with the tool supplying constraints, metrics and checks.

## Ratsnest as a minimum spanning tree

The ratsnest shows unrouted connectivity. Each net with pins $p_1, \dots, p_k$ is drawn as a **minimum spanning tree** over the pins, usually with Manhattan or Euclidean distance, recomputed as parts move. Prim's algorithm on the complete graph of $k$ pins takes $O(k^2)$, which is fine for typical board nets. Example with Manhattan distance: pins listed as (0, 0), (10, 8), (10, 0) and (2, 9) give tree edges of 10, 8 and 9 for a length of 27, while chaining the pins in listed order gives $18 + 8 + 17 = 43$ and a star from the first pin gives $18 + 10 + 11 = 39$. The bounding-box half-perimeter, 10 plus 9 = 19, only lower-bounds the routed length. Summed over nets, ratsnest length is a quick, incremental placement score.

## What the ratsnest misses

Total ratsnest length weights every net equally and ignores everything but distance. A placement that saves 15% of total length but moves a DDR memory farther from its controller is worse, because those nets carry tight timing and SI limits. It ignores **noise coupling**: a switching regulator's inductor beside an analog input. It ignores **thermal** interaction: two hot parts side by side raise each other's junction temperature through shared copper. It ignores **assembly**: on a board reflowed twice, heavy parts on the bottom fall off in the second pass. A useful automatic placer must therefore weight nets, carry keep-out and proximity rules, and score thermal and assembly constraints, which is why board placement remains largely interactive.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4, on wirelength models. Montrose, "Printed Circuit Board Design Techniques for EMC Compliance" (2nd ed., 2000), on placement for noise.
