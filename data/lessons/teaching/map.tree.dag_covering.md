# DAG covering and logic duplication

Tree partitioning forces every multi-fanout node to be a gate output. DAG covering drops that rule: a match may cross a fanout point, which means the logic below it is duplicated into every match that crosses. Duplication can cut delay and sometimes area. Understanding when it pays is the bridge from DAGON to the cut-based mappers in ABC.

## The forced cut and its cost

Let s = NAND(a, b) feed m consumers y_i = NAND(INV(s), c_i). A tree mapper roots a tree at s, so s is built as a NAND2, and each consumer's tree NAND(INV(s), c_i), with s as a leaf, can only be covered by INV plus NAND2. With INV = 1, NAND2 = 3 and m = 3, the area is 3 + 3(1 + 3) = 15. The consumers' structure NAND(INV(NAND(a, b)), c) is exactly a NAND3, but the partition hides it because s is a leaf.

## Duplication and its accounting

Allow matches to cross s. Each consumer can now be a NAND3 over a, b and c_i, so s is computed inside every NAND3 and never built alone. With NAND3 = 4 the area is 3 x 4 = 12, better than 15. If NAND3 cost 6, duplication would give 18 and lose. Mixing does not help in this example: keeping s for even one consumer pays its whole NAND2, so the best is the cheaper of the two pure choices. In general, the accounting problem is that a shared node's cost is paid once if built, but its absorption is paid in every crossing match, which is why exact area DAG covering is NP-hard.

## Delay is where duplication shines

For delay, duplication is safe and polynomial. In the example the tree cover's path is NAND2, INV, NAND2; the duplicated cover's is one NAND3. With delays 2, 1, 2 and 3, arrival drops from 5 to 3. More generally, under load-independent delays the best arrival at a node depends only on its transitive fanin, so a mapper can compute it at every node of the DAG as if each node had its own private copy of its fanin cone, then build only the copies it needs. Area is recovered afterwards on nodes with slack. Load-dependent delay complicates this, because duplication changes fanout loads.

## Going deeper

Kukimoto, Brayton and Sawkar, "Delay-optimal technology mapping by DAG covering", DAC 1998; Mishchenko, Chatterjee and Brayton, "Improvements to technology mapping for LUT-based FPGAs", IEEE TCAD 2007.
