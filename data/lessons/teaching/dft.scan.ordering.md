# Scan-chain ordering

Which flop follows which in a chain has no effect on test quality, so the order is free to minimize the chain's routing. After placement, the DFT tool reorders each chain to shorten its wirelength and congestion while keeping its endpoints and clock-domain constraints. The problem is a travelling-salesman path, and its heuristics are the same ones placement and routing engines use.

## Wirelength of an order

A chain runs from a scan-in point, through every flop in its order, to a scan-out point. With Manhattan routing, its wirelength is the sum of $|x_i - x_j| + |y_i - y_j|$ over consecutive pairs. For scan-in at (0, 0), scan-out at (8, 11) and the order (11, 8), (9, 2), (4, 10), (10, 2), the hops are 19, 8, 13, 14 and 11, for 65. An order taken from the netlist, or from synthesis, ignores placement and can zigzag across the die, which is why reordering after placement routinely saves most of the scan wiring.

## Greedy and exact ordering

Finding the shortest path through all flops between fixed endpoints is a **Hamiltonian path** version of the travelling-salesman problem, NP-hard in general. A **nearest-neighbour** heuristic starts at scan-in and repeatedly moves to the closest unvisited flop, then ends at scan-out; it is fast but can leave long final hops. For a handful of flops, enumerating all $n!$ orders finds the optimum ($5! = 120$ orders). Production tools partition the die into regions, order within each region with greedy plus local improvement (2-opt swaps), and connect regions, keeping runtime linear in practice.

## Constraints and consequences

Reordering must keep each chain within its clock domain or insert lockup latches at crossings, keep the chain lengths balanced, and respect hold timing on the shift path: two flops placed side by side with a short wire and clock skew can violate hold during shift, which is fixed with buffers. Ordering also interacts with compression and diagnosis: chain assignment determines which flops share a compactor group, and a chain's physical path matters when localizing a broken chain. Reordering is only legal after placement and before routing, and the new order must be written back to the scan definition the ATPG uses.

## Going deeper

Wang, Wu and Wen (2006), chapter 2; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), on wirelength estimation.
