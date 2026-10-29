# Selecting the LUT cover from best cuts

Labeling gives every node a best cut, but most nodes never become LUTs. The cover is selected backward from the outputs: only nodes that some chosen LUT needs as an input are implemented. This step decides the LUT count and shows where logic is duplicated, which is the main cost of depth-optimal mapping.

## Backward selection

Start with the primary outputs in a worklist. Pop a node, implement it as a LUT over its best cut, and push every cut leaf that is a gate and not yet implemented. Primary inputs need nothing. When the worklist empties, the implemented nodes are the cover. Example: outputs n4 and n6, best cuts n4: {n1, n3}, n6: {a, b, n3}, n1: {a, b}, n3: {c, d}. Selection implements n4, then n1 and n3, then n6, whose leaves are already handled: 4 LUTs. Gates never pushed, such as an n2 absorbed inside n3's LUT, cost nothing.

## Counting LUTs and duplication

The LUT count is the number of implemented nodes. A LUT rooted at r over cut C absorbs every gate on paths from C to r, excluding the leaves. A gate can lie inside the cones of two or more LUTs; then it is computed twice, which is **duplication**. To count duplicated gates, list each selected LUT's absorbed gates and count those that appear more than once. Depth-optimal mapping duplicates freely, since each node's best cut was chosen for its own depth without regard to sharing, so a network of 10 gates may map into 6 LUTs that together absorb 14 gate copies.

## Depth of the selected cover

With unit delay, a LUT's depth is one more than the deepest LUT among its inputs, and primary inputs have depth 0. When every node uses its depth-optimal cut, each implemented node sits at its label, so the network depth is the largest label among the outputs. Selection does not change depth; it only decides area. That separation is what area recovery exploits: depth is fixed by labels on the critical paths, while off-critical nodes may switch to cuts that reduce the LUT count or the duplication. Any such switch must keep each LUT within the required depth of its consumers.

## Going deeper

Cong and Ding, "FlowMap", IEEE TCAD 1994, section on mapping generation; Betz, Rose and Marquardt, "Architecture and CAD for Deep-Submicron FPGAs" (1999), for the LUT-mapping stage of the FPGA flow.
