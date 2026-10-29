# Priority cuts

Exhaustive enumeration stores every k-feasible cut of every node, which costs memory and time that explode with k. Priority cuts keep only the C best cuts per node, ranked by the current mapping objective, and build each node's candidates from its fanins' priority cuts. The method is the basis of ABC's if mapper and makes 6-input LUT mapping of large designs practical.

## Ranking candidates by depth

For depth-oriented mapping, the depth of a cut is one more than the largest best-depth among its leaves, with primary inputs at depth 0 and unit LUT delay. A node's best depth is the smallest depth over its cuts. With leaf depths n1: 2, n2: 2, n3: 2, the cut {n2, n3} has depth 3 and {a, b, n4} with n4 at depth 3 has depth 4. The mapper computes each candidate's depth from stored leaf depths in constant time per leaf, so ranking a node's candidates is cheap.

## Keeping the best C cuts

At node n, merge the stored cuts of the fanins, each fanin offering its C priority cuts plus its trivial cut. Rank the resulting candidates by the objective and keep the top C. In delay mode the key is depth, with tie-breakers such as cut size and area flow; in area-recovery passes the primary key becomes area flow or exact area. For candidates C1 (depth 2, area flow 2.5), C2 (depth 2, area flow 1.7), C3 (depth 3, area flow 1.0) and C = 2, delay mode keeps C1 and C2, while area-flow mode keeps C3 and C2. The same node can therefore store different cuts in different passes.

## Why bounded cut sets work

With C cuts per fanin, a two-input node evaluates $(C + 1)^2$ unions, independent of how many cuts exhaustive enumeration would have produced; C = 8 gives 81 candidates. Memory is linear in the number of nodes times C. The price is completeness: a cut that is not among the top C at some node can never appear in its fanouts' candidates, so the mapping is no longer guaranteed depth-optimal. In practice small C loses little, because the best cuts of a node tend to be built from the best cuts of its fanins, and repeated mapping passes with different rankings recover much of the remaining quality.

## Going deeper

Mishchenko, Cho, Chatterjee and Brayton, "Combinational and sequential mapping with priority cuts", ICCAD 2007; the ABC source code for the if command.
