# Technology mapping by tree covering

Technology mapping chooses library cells to implement an optimized network. The classic formulation decomposes the network into a subject graph of NAND2 and inverters, describes each library cell as a pattern graph in the same primitives, and covers the subject graph with patterns at minimum cost. On trees, dynamic programming finds the optimal cover in linear time.

## Subject graphs and patterns

The subject graph expresses the logic using only two-input NANDs and inverters. Every library cell is given as one or more pattern trees in the same primitives: INV is $INV(x)$, NAND2 is $NAND(x, y)$, AND2 is $INV(NAND(x, y))$, OR2 is $NAND(INV(x), INV(y))$, NOR2 is $INV(NAND(INV(x), INV(y)))$, and NAND3 is $NAND(x, INV(NAND(y, z)))$. A **cover** assigns cells to the subject graph so that every node is covered by exactly one pattern match and every match's inputs are outputs of other matches or primary inputs. With only INV and NAND2 in the library the cover is forced: each node is its own match, and the area is the sum of the base cell areas.

## Minimum-area cover by dynamic programming

For a tree, visit nodes bottom up. At each node try every pattern that matches with its root there; the cost of that choice is the cell area plus the best costs already computed for the subtrees at the pattern's inputs. Keep the minimum. For $NAND(INV(a), NAND(NAND(b, c), d))$, OR2 would need two inverted children and does not match at the root; NAND2 at the root costs 3 plus INV(a) at 2 plus the cheapest cover of $NAND(NAND(b, c), d)$, which is two NAND2s at 6, for 11 in total. Matching both input orders of each NAND handles commutativity. Each node is visited once with a fixed pattern set, so the algorithm is linear in the tree size and optimal for the tree.

## Delay covers and the limits of trees

For delay, the same dynamic program uses cost = cell delay plus the latest arrival among the pattern's inputs, giving the minimum arrival time at the root under a load-independent delay model. Real networks are DAGs, not trees. The classic approach (DAGON, then the mapper in MIS) cuts the DAG at multiple-fanout points into a forest of trees and maps each tree optimally, which loses optimization across tree boundaries. Load-dependent delay breaks the optimal substructure. Modern mappers, such as ABC's map and amap, work directly on the AIG with cuts and supergates and use structural choices to reduce the bias of a single subject graph.

## Going deeper

Keutzer, "DAGON: technology binding and local optimization by DAG matching", DAC (1987); De Micheli (1994), the chapter on cell-library binding.
