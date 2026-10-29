# The M1, M2 and M3 moves

Simulated annealing over normalized Polish expressions needs a move set that can reach every slicing floorplan, changes the floorplan a little at a time and keeps expressions valid. Wong and Liu's three moves, M1, M2 and M3, do exactly that. Implementing them correctly, including the legality check M3 needs, is a standard exercise in designing a search space for a combinatorial optimiser.

## M1: swap adjacent operands

M1 picks two operands that are adjacent in the sequence of operands, ignoring operators between them, and swaps them. In 6 3 H 2 4 V 5 1 H V H, operands 4 and 5 are adjacent in the operand order 6, 3, 2, 4, 5, 1, so M1 gives 6 3 H 2 5 V 4 1 H V H. Every operator stays where it was, so the tree shape is unchanged and only two leaves trade places; balloting and normalization are automatically preserved. The effect on the floorplan is two blocks exchanging positions, which changes area when their shapes differ and wirelength when their connections differ.

## M2: complement a chain

A chain is a maximal run of consecutive operators, such as V H V. M2 picks one chain and complements every operator in it, turning V into H and H into V. In 1 4 V 3 V 2 5 H V the chains are V, V and H V; complementing the second gives 1 4 V 3 H 2 5 H V. Because a normalized chain alternates, its complement alternates too, so the result is normalized; operand and operator counts do not change, so balloting holds. Geometrically, M2 turns side-by-side arrangements into stacked ones at several levels of the tree at once, a larger move than M1. Complementing only part of a chain could create two equal adjacent operators, which is why the whole chain flips.

## M3: swap an operand and an operator, with a check

M3 swaps an adjacent operand and operator, which changes the tree's shape and is the only move that lets the search reach different tree topologies. It can break validity in two ways. Moving an operator left can violate balloting, leaving a prefix with as many operators as operands, as in 1 H 4 2 3 H 5 H V, which starts with an operator after one operand. It can also create two identical adjacent operators, as in 1 4 H 2 3 H H 5 V. The annealer therefore checks the candidate: count operands and operators along the prefix and compare neighbours, both linear time, and the check can be made constant time by tracking counts at the swap position. Rejected candidates are simply redrawn. Together the three moves connect the whole space of normalized expressions.

## Going deeper

Wong and Liu, "A new algorithm for floorplan design" (DAC 1986); Sarrafzadeh and Wong, "An Introduction to VLSI Physical Design" (1996), chapter on floorplanning.
