# Placement

Placement assigns every standard cell a legal location in the rows of the floorplan. Its objective is short wires, because wirelength drives delay, power and routability, subject to density limits and timing. Modern placers are analytical: they solve a smooth optimization for approximate positions, then legalize them into rows and sites and refine locally. OpenROAD's global placer and the RePlAce algorithm are open-source examples.

## Half-perimeter wirelength

A placer cannot route every candidate placement, so it estimates each net's wire with the half-perimeter wirelength (HPWL): the half-perimeter of the smallest axis-aligned box that contains all its pins,
$$\mathrm{HPWL} = (x_{max} - x_{min}) + (y_{max} - y_{min}).$$
For pins at $(2, 3)$, $(7, 1)$ and $(4, 6)$ the box spans $x$ from 2 to 7 and $y$ from 1 to 6, so HPWL $= 5 + 5 = 10$. For two- and three-pin nets HPWL equals the shortest rectilinear Steiner tree length; for larger nets it is a lower bound. It is cheap and is the standard placement objective.

## Total wirelength and incremental moves

The placer's objective sums HPWL over all nets, often with timing-critical nets weighted higher. Because a cell touches only its own nets, moving one cell changes only those nets' terms. A placer evaluating a move recomputes HPWL for just the affected nets and takes the difference, which is how detailed placement can try millions of swaps cheaply. If a move shrinks one net's box by 4 and stretches another's by 1, the total improves by 3. The same locality makes HPWL easy to update incrementally inside an analytical placer, where the smooth approximations of HPWL are differentiated with respect to each cell's coordinates.

## Density, legality and what goes wrong

Minimizing wirelength alone collapses all cells onto one point. So global placement adds a density constraint that spreads cells, legalization snaps them into rows without overlap, and detailed placement polishes locally. A placement that packs cells too tightly in one region creates routing congestion: more wires want to cross a region than there are tracks. A placement that ignores timing makes critical nets long. Both failures surface only in later stages, routing and timing analysis, which is why placers now include congestion and timing estimates in their objective. Placement quality is therefore judged after routing, not by HPWL alone.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4; Cheng et al., "RePlAce: Advancing Solution Quality and Routability Validation in Global Placement", IEEE Transactions on CAD (2019).
