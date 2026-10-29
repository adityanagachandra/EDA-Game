# Incremental HPWL after a cell move

Detailed placers and annealers evaluate thousands of candidate moves per second, and each one needs the change in wirelength. Recomputing every net is far too slow. Incremental HPWL updates only the nets touching the moved cell, and with careful bookkeeping only the boxes that actually change, which is the difference between a usable and an unusable local optimiser.

## Only incident nets change

When cell $c$ moves from $p$ to $p'$, only the nets that contain one of $c$'s pins can change length. The change in total wirelength is
$$\Delta W = \sum_{e \ni c} w_e\left(\mathrm{HPWL}_e(p') - \mathrm{HPWL}_e(p)\right).$$
A placer stores, for each cell, the list of its nets (the netlist hypergraph's incidence lists), so this sum touches a handful of nets regardless of design size. Example: $c$ moves from (20, 20) to (30, 25) and is on one net with other pins at (10, 15) and (25, 40). The box was 15 by 25 (HPWL 40) and becomes 20 by 25 (HPWL 45), so $\Delta W = 5$.

## Recomputing a box after a move

For each affected net, the new box can often be found without a full rescan. If the cell moves outward, the box simply grows to include the new position. If the cell was the unique extreme pin on some side and moves inward, the box may shrink, and the new extreme must be found among the other pins. Keeping for each net its minimum and maximum together with the count of pins at each extreme makes most updates constant time: the rescan is needed only when the count of pins at an extreme drops to zero. For a three-net example where a cell leaves the left edge of one box, enters the interior of a second and pushes out the top of a third, the three deltas are computed separately and summed.

## Cost and correctness in practice

Incremental evaluation turns a move from $O(\text{pins in the design})$ to $O(\text{pins on the cell's nets})$, which for typical cells is a few dozen. High-fanout nets break this: a cell on a 5000-pin net makes every move expensive, so placers often skip or cap such nets in detailed placement. Two classic bugs are double-counting a net that contains two pins of the moved cell, and evaluating a swap of two cells as two independent moves when they share a net: the shared net must be evaluated with both cells moved. Differential tests against a full recomputation on random moves catch both.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4 on detailed placement; the OpenROAD dpl source for its incremental cost evaluation.
