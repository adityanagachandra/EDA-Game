# Pick-and-place data and throughput

Placement machines pick parts from feeders and place them at coordinates from the centroid (pick-and-place) file: reference designator, X, Y, rotation and side. The board tool generates this file from the placed footprints. Throughput depends on machine speed in components per hour, how placements are split across machines in a line, and how cleanly the CAD data maps to the machine's own part library. Each is a small computation with large consequences on the factory floor.

## Machine time from CPH

Placement machines are rated in **components per hour** (CPH). Time for a build is total placements divided by the sustained CPH: $t = P \cdot B / \text{CPH}$. Example: 800 placements per board, 250 boards and 25,000 CPH give $200{,}000 / 25{,}000 = 8$ hours. Rated CPH is measured on ideal small chips; large parts, fine-pitch ICs, nozzle changes and tray-fed parts all slow a real job, so planners use a sustained figure per job. Setup time (loading feeders for every BOM line, verifying them) adds a fixed cost per build that does not shrink with quantity, which is why small builds cost more per board and why reducing BOM lines matters.

## Balancing a placement line

A line usually has several machines in series: a fast chip shooter for small passives and a slower fine-pitch placer for ICs. Each board passes through every machine, so the line's cycle time per board is set by the slowest machine: $t = \max_k P_k / \text{CPH}_k$. Example: 400 parts at 30,000 CPH take 48 s, and 100 parts at 8,000 CPH take 45 s, so the cycle is 48 s per board. Moving placements that both machines can handle from the busier to the idler machine balances the line and shortens the cycle. This is a small scheduling problem, a partition of placements across machines subject to which machine can place which part, and line-balancing software solves it from the centroid file and the BOM.

## Centroid data and rotations

The centroid file lists each part's position and rotation in the board's coordinate system, referenced to the footprint's origin, which should be the part's center. Rotation conventions differ: the CAD library's zero orientation for a package may not match the machine's feeder orientation, so the machine applies an offset per package: machine angle = (CAD angle + offset) mod 360. Example: a part at 270 degrees with a +180 degree package offset is placed at 90 degrees. A wrong offset rotates every instance of that package, which is why first-article inspection checks polarity on each part type. Bottom-side parts add a mirror transformation whose sign convention varies by tool, a classic source of rotated parts. Standard library orientation rules exist to make this mapping predictable.

## Going deeper

IPC-7351, section on zero component orientation. IPC-2581, on placement data in a neutral format.
