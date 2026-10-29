# Routing grids, gcells and tracks

Before any search, a router builds its resource model: tracks on each layer and gcells that group them. The number of usable tracks crossing a gcell boundary becomes that edge's capacity in global routing, and every congestion number downstream is measured against it. Getting capacity right, including blockages and reserved slack, is the first computation in a routing engine.

## Tracks from pitch and offset

Each routing layer carries parallel tracks at a fixed pitch $p$, set by the layer's minimum width plus minimum spacing, starting from an offset $o$ given in the technology or DEF track definitions. Vertical tracks sit at $x = o + kp$ for $k = 0, 1, 2, \dots$ A gcell column covering $[x_0, x_1)$ contains the tracks with $x_0 \le o + kp < x_1$, so the count is $\lfloor (x_1 - 1 - o)/p \rfloor - \lceil (x_0 - o)/p \rceil + 1$ for integer coordinates.

Example: $o = 70$ nm, $p = 140$ nm, gcell $[2800, 5600)$. The first track is at $k = \lceil 2730/140 \rceil = 20$ (x = 2870) and the last at $k = \lfloor 5529/140 \rfloor = 39$ (x = 5530), so the gcell holds 20 tracks. Gcells are usually sized to a fixed number of standard-cell row heights, so the track count per gcell differs from layer to layer.

## Blockages remove tracks

A macro, a power strap or a pre-routed wire makes some tracks unusable. A track is lost if a wire on it would violate spacing to the obstruction, that is, if its centerline lies within half the wire width plus the minimum spacing of the obstruction's edge. So a blockage over $[b_1, b_2]$ removes every track with $b_1 - m < x < b_2 + m$, where $m = w/2 + s$.

Example: tracks at $70 + 140k$ nm in a $[0, 2660)$ gcell (19 tracks), a blockage over $[1140, 1390]$ and $m = 35$ nm. The lost window is $(1105, 1425)$, which contains the tracks at 1190 and 1330, so the gcell keeps 17. Global routers compute these reductions per gcell edge and layer; a partial blockage reduces capacity proportionally rather than to zero.

## Capacity with reserved slack

The capacity of a gcell edge sums the usable tracks crossing it over all layers whose preferred direction runs across that boundary. Flows then reserve slack: a layer adjustment of 30% keeps only 70% of a layer's tracks, rounded down, leaving room for pin access, local detours and rules the coarse model cannot see. Lower layers usually get larger adjustments because pins crowd them.

Example: metal 2 has 10 tracks crossing an edge with a 40% adjustment, metal 4 has 13 with 60%. Capacity is $\lfloor 10 \cdot 0.6 \rfloor + \lfloor 13 \cdot 0.4 \rfloor = 6 + 5 = 11$. Too little adjustment and global routing reports a clean result that detailed routing cannot realise; too much and the global router detours nets for congestion that does not exist.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); the LEF/DEF reference for TRACKS and GCELLGRID statements; the OpenROAD grt documentation on layer adjustments.
