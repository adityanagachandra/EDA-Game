# Rows, sites, rails and utilization

Before a placer moves a single cell, the floorplan fixes the placement grid: how many rows fit in the core, how many sites each row holds, where the rails go and what fraction of the sites cells will fill. This arithmetic sets capacity, density targets and power-rail structure, and placers, legalizers and density models all compute it.

## Rows in a core

A core area of height $H$ holds $\lfloor H / h_{row} \rfloor$ complete rows of height $h_{row}$, the library's cell height, stacked from the bottom edge. A 54,100 nm tall core with a 270 nm row height holds $\lfloor 200.37 \rfloor = 200$ rows; the leftover 100 nm is unusable for cells. Floorplanners therefore size cores in whole multiples of the row height. Each row is oriented: in alternate rows cells are flipped vertically (orientation N and FS in LEF/DEF), so that the VSS rail of one row coincides with the VSS rail of the next and the VDD rails likewise.

## Sites, capacity and rails

Each row is divided into sites of the width defined by the library's SITE. A row of width $W$ holds $\lfloor W / w_{site} \rfloor$ sites, and the core holds rows times sites per row. A 40,000 nm wide core with 54 nm sites holds 740 sites per row; with 200 rows that is 148,000 sites. Every cell occupies a whole number of consecutive sites. Because flipped rows share their boundary rails, $N$ rows need $N + 1$ horizontal rails: one at each internal boundary plus one at the top and bottom edges. With 200 rows that is 201 rails, alternating VSS and VDD. These rails are then strapped to the power grid at regular intervals.

## Utilization and why it is not 100%

Utilization is the fraction of sites occupied by cells. If cells fill 99,160 of 148,000 sites, utilization is 67.0%. Designs rarely target more than about 70% to 80% in the core, because the router needs room to reach pins, the timing optimizer needs space to upsize cells and insert buffers, and physical-only cells (taps, endcaps, decaps) take sites too. Placers model this with density targets: global placement spreads cells so that every bin stays below a target density, and the legalizer then snaps them to sites. A utilization number alone hides local hot spots, which is why congestion and density maps per bin matter more than the global figure.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapters 3 and 4; the LEF/DEF Language Reference, ROW and SITE statements.
