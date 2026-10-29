# Preferred directions and directional capacity

Each metal layer has a preferred direction, and directions alternate up the stack. The convention keeps wires on one layer parallel, which simplifies spacing and lithography and makes track grids regular. It also means each bend in a route becomes a via, and each gcell edge draws its capacity only from the layers running across it. Directional capacity is the second half of the resource model.

## Directional capacity

If layer $M_i$ is horizontal, its tracks run east-west and carry wires across the boundaries between horizontally adjacent gcells. So the capacity of an east-west crossing is the sum of tracks over the horizontal routing layers only, and the capacity of a north-south crossing sums the vertical layers. Metal 1 is often reserved for cell-internal wiring and pins, so routing capacity starts at metal 2.

Example: a stack with M3, M5, M7 and M9 horizontal holding 13, 6, 12 and 5 tracks per gcell gives a horizontal edge capacity of 36. The vertical layers contribute nothing to that edge unless the router allows wrong-way wiring. Upper layers usually have wider pitch, hence fewer tracks, but lower resistance, which matters later for timing-driven layer assignment.

## Filling layers bottom-up

A simple layer assigner places wires on the lowest layer of the right direction first and moves up only when a layer is full. Low layers are close to the pins, so using them saves vias; high layers are reserved for long and critical nets. Given a demand of $D$ crossings and per-layer capacities $c_1, c_2, \dots$ for the horizontal layers in order, the highest layer used is the first $k$ with $c_1 + \dots + c_k \ge D$.

Example: horizontal layers M3 with 8 tracks and M5 with 14. A demand of 22 crossings fills M3 and needs all 14 of M5, so M5 is the highest layer used. If the demand exceeds the total over all horizontal layers, the edge overflows no matter how layers are assigned, which is why global routing can work on the 2D projection with summed capacities.

## Wrong-way wiring

Some routers allow short wires against a layer's preferred direction, at a cost. A wrong-way step uses a track that crosses other tracks of the same layer, so it blocks several of them; routers charge a multiple $f > 1$ of the normal step cost. The alternative to a wrong-way jog of $n$ steps is a via down to a layer of the other direction, $n$ normal steps and a via back up.

With step cost 1, via cost $c_v$ and wrong-way factor $f$, the jog is cheaper when $fn < n + 2c_v$, that is $n < 2c_v/(f - 1)$. With $c_v = 4$ and $f = 3$ the bound is 4, so jogs of up to 3 steps stay on the layer. Advanced nodes often forbid wrong-way wiring on the lowest layers entirely, because their patterning only prints lines in one direction.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.), on routing grids and layer stacks.
