# Tiled parallel geometry with halos

The standard way to run geometry checks on many cores, many machines or a GPU is to cut the die into tiles and process them independently. Correctness requires each worker to see a halo around its tile at least as wide as the largest interaction distance, and to report only results it owns. The halo makes the work redundant at the borders, and the tile size trades load balance against that redundancy. These are the numbers a fast-mode engine and an agent scheduling it must reason about.

## Cutting the die into tiles

A $W \times H$ die cut into $T \times T$ tiles has $\lceil W/T \rceil \cdot \lceil H/T \rceil$ tiles, the last row and column possibly partial. Example: a $2350 \times 1180$ um die with 100 um tiles has $24 \times 12 = 288$ tiles. Each tile is an independent task, so a scheduler can spread them over threads, machines or GPU thread blocks. Load balance depends on the data, not the area: a tile over a memory array may hold a hundred times the shapes of a tile over a routing channel. Engines therefore make tiles small enough that there are many more tiles than workers, or size tiles adaptively by shape count, and they schedule the heaviest tiles first.

## Halos and ownership

A spacing check near a tile border needs shapes from the neighbouring tile. So each worker reads its tile expanded by a halo of width $h$ at least equal to the largest interaction distance of the rules it runs. A shape is therefore read by every worker whose expanded tile it overlaps. Example: tiles of 40 um with a 5 um halo, and a shape spanning $[78, 83) \times [100, 110)$: in x it meets expanded tiles 1 and 2, in y tile 2 only, so two workers read it. To avoid double reporting, each violation is owned by one tile, for example the tile containing the lower-left corner of its marker. Too small a halo silently misses violations at tile borders, which is why the halo is derived from the rule deck, never guessed.

## Halo overhead and the fast-accurate loop

Each worker processes $(T + 2h)^2$ of area for $T^2$ of owned area, an overhead of $((T + 2h)^2 - T^2)/T^2$. Example: 200 um tiles with a 5 um halo process 10.25 percent extra; 50 um tiles with the same halo process 44 percent extra. Smaller tiles balance load better and fit GPU memory, but the halo is fixed by the rules, so overhead grows as tiles shrink. This trade-off is a parameter an agent can tune in a fast mode: large tiles and a subset of rules for quick feedback during an ECO loop, then the full deck for sign-off. A GPU or tiled back end must still produce exactly the results of the serial engine, which a differential test on reference layouts verifies.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on physical verification; the KLayout documentation on tiling mode in DRC scripts.
