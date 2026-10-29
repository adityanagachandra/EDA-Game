# Tree covering in technology mapping

Technology mapping covers a subject graph of simple gates with library cells. When the subject graph is a tree, minimum-area covering is solved exactly by dynamic programming: for every node, the best cover of the subtree rooted there is the best pattern match at that node plus the best covers of the subtrees the pattern leaves exposed. Keutzer's DAGON cast this as tree covering, and the same DP with delay labels underlies classic mappers and the cut-based mappers in ABC.

## One DP step at a node

Decompose the logic into NAND2 and INV nodes. Each library cell is a pattern over these: INV matches one INV node, NAND2 one NAND node, AND2 is an INV fed by a NAND2, and BUF is an INV fed by an INV. For a chain listed from the input, define $\text{best}(i)$ as the cheapest cover of nodes 1 to i with node i the root of a pattern.

At an INV node i, two options exist: a 1-node INV, costing $\text{best}(i-1) + c_{INV}$, or a 2-node pattern covering i and i-1, costing $\text{best}(i-2)$ plus the area of AND2 (if node i-1 is a NAND2) or BUF (if it is an INV). Example: all nodes INV, areas INV 1 and BUF 3, with $\text{best}(3) = 3$ and $\text{best}(4) = 4$. Then $\text{best}(5) = \min(4 + 1, 3 + 3) = 5$.

## Covering a chain

Run the step from the input end to the output. Chain INV, NAND2, INV, NAND2 with areas INV 2, NAND2 2, AND2 3, BUF 3. Node 1 (INV): 2. Node 2 (NAND2): only the NAND2 pattern matches, $2 + 2 = 4$. Node 3 (INV fed by NAND2): $\min(4 + 2, 2 + 3) = 5$, by AND2 covering nodes 2 and 3. Node 4 (NAND2): $5 + 2 = 7$. Minimum area 7: INV, AND2, NAND2.

The algorithm is linear in the number of nodes times the number of patterns that can match at a node. The result is optimal for the given subject tree and library, but the subject tree is not unique: a different decomposition of the same function into NAND2 and INV can map better, which is why mappers try several decompositions or add redundant choice nodes. Recording the chosen pattern at each node lets the mapper emit the netlist by a final top-down traversal.

## Trees, subtrees and the DAG problem

On a tree, subtrees share no nodes, so the cost of a pattern is its area plus the independent best costs of the subtrees at its inputs. Take an output INV fed by a NAND2 whose two inputs are chains with best costs 5 and 5, and areas INV 1, NAND2 3, AND2 4. Covering the NAND2 alone costs $3 + 5 + 5 = 13$, so the INV on top costs 14; AND2 covering both costs $4 + 5 + 5 = 14$. Either gives 14.

Real netlists are DAGs: a node with several fanouts is shared. Tree mappers cut the DAG into trees at multi-fanout points and map each optimally, which loses optimization across the cuts. DAG covering for minimum area is NP-hard, while depth-optimal mapping of DAGs into lookup tables is polynomial, as FlowMap showed. Modern mappers such as those in ABC enumerate cuts on the DAG and choose among them with area-recovery heuristics.

## Going deeper

Keutzer, "DAGON: technology binding and local optimization by DAG matching" (DAC, 1987); Cong and Ding, "FlowMap: an optimal technology mapping algorithm for delay optimization in lookup-table based FPGA designs" (IEEE TCAD, 1994); Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting" (DAC, 2006), and the ABC documentation.
