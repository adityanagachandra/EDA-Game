# Detailed placement by cell swapping

Detailed placement improves a legal placement with local moves that keep it legal. The simplest move swaps two cells of equal width; more general moves shift, reorder or reassign cells. Each move is accepted only if it lowers the cost, evaluated incrementally on the nets that touch the moved cells. Swaps are the building block of every detailed placer.

## Evaluating a swap

Swapping cells a and b exchanges their positions. If they have equal widths, the result is automatically legal. Only nets touching a or b change, so
$$\Delta W = \sum_{e \ni a \text{ or } b} \left(\mathrm{HPWL}_e^{after} - \mathrm{HPWL}_e^{before}\right).$$
With a at (8, 25) and b at (14, 7), a on one net with pins at (12, 14) and (3, 14), b on two nets, the deltas might be -2 for a's net and $+3$ and -11 for b's nets, a total of -10: the swap shortens wiring by 10 um and is accepted. A swap with $\Delta W \ge 0$ is rejected in a greedy detailed placer. Unequal widths make a swap need extra shifting of neighbours, which is why many swap moves are restricted to same-width cells or to cells with enough whitespace nearby.

## Shared nets and double counting

A net that contains both a and b must be evaluated with both cells moved at the same time. If the net's other pins are fixed, swapping a and b exchanges two of its pin positions and leaves the set of positions, and therefore the HPWL, unchanged. Evaluating the swap as "move a, then move b" with the other cell still in its old place would get this net wrong. Implementations collect the union of nets of both cells, deduplicate it, and recompute each once with both new positions. The same care applies to cells with two pins on one net, where the net appears twice in the cell's net list.

## Global swap and search strategy

Trying all pairs is quadratic, so detailed placers choose candidates smartly. Global swap, as in FastDP, computes for each cell its optimal region, the area where its nets would be shortest, and tries to swap it with a cell or a gap in that region. Vertical swap tries cells in the rows just above and below. The placer sweeps over all cells, accepts improving moves, and repeats passes until the gain per pass falls below a threshold. Because each move is evaluated incrementally in time proportional to the pins involved, a pass over a million cells takes seconds. Timing-driven variants add the change in slack or in weighted wirelength to the cost.

## Going deeper

Pan, Viswanathan and Chu, "An efficient and effective detailed placement algorithm" (ICCAD 2005), on FastDP; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 4.4.
