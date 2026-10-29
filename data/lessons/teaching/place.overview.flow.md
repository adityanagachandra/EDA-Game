# The physical design flow

Physical design turns a gate-level netlist into geometry. Floorplanning fixes the chip outline, the core, the rows, the macros and the I/O; placement gives every standard cell a legal location; clock tree synthesis buffers the clock to the placed flip-flops; routing draws the wires. Placement sits in the middle, and a tool designer cares because its decisions set the wirelength, timing, congestion and power that every later stage inherits and can only partly repair.

## Floorplan, place, clock, route

The physical flow runs in a fixed order because each step needs the previous one's output. Floorplanning decides the die and core size, cuts the core into rows of sites, places the large macros (memories, analog blocks), assigns I/O pins and plans the power grid. Placement then works inside that region: it reads the netlist, the cell abstracts and the floorplan, and writes a location and orientation for every standard cell. Clock tree synthesis needs those locations, because it builds a buffered tree to the flip-flops where they now sit. Routing needs everything fixed, because it draws real metal between real pins. In OpenROAD these are separate engines over one database: floorplan and pin placement, global placement (gpl), detailed placement (dpl), CTS and global and detailed routing.

## Global, legal, detailed

Placement itself has three phases. Global placement treats cells as points or small blobs, minimises an estimate of wirelength and spreads cells so that no region is over-full, but tolerates overlap. Legalization snaps every cell into a row, aligns it to the site grid and removes all overlap with as little movement as possible. Detailed placement then improves the legal result with local moves: swapping cells, reordering short windows of a row, or matching groups of cells to slots. The split exists because the global problem is a huge continuous optimisation that is tractable only when legality is relaxed, while legality is a discrete, local property that is cheap to restore once cells are roughly in place. Wirelength is the main objective because it is fast to evaluate and correlates with wire capacitance, delay, power and routing demand.

## Why placement mistakes propagate

Placement is NP-hard even for wirelength alone, since it contains problems such as graph partitioning and minimum linear arrangement, so every placer is a heuristic and its errors flow downstream. Cells packed too tightly around a group of highly connected logic leave the router too few tracks, and the router reports overflow that it cannot fix without moving cells. Critical nets stretched by a wirelength-only objective fail timing, and CTS and routing can recover only a little with buffers and layer choices. A poor floorplan constrains the placer before it starts: if macros cut the rows into slivers or point their pins away from their logic, no placer can produce short nets. A tool designer therefore builds feedback into the placer: density targets, timing weights and congestion estimates that approximate what later stages will see.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapters 1 and 4; the OpenROAD project documentation for its flow.
