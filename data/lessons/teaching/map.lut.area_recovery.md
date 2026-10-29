# Required depths and LUT area recovery

A depth-optimal LUT mapping is fast but large. Area recovery keeps the depth and shrinks the cover by letting nodes with slack use deeper but cheaper cuts. The mechanism is the one static timing uses: propagate required times backward, compute slack, and spend it. For a mapper designer this is where most of the LUT count is won.

## Required depths

Fix a target depth D, at least the optimal depth. Each primary output gets required depth D. Then traverse the selected cover in reverse topological order: for a LUT rooted at r with required depth R(r), each gate leaf l of its cut gets R(l) = min(R(l), R(r) - 1). A node feeding several LUTs takes the tightest requirement. Example: n8 is an output with D = 4; its LUT has leaves n1 and n2, so both get 3; n2's LUT has leaf n1, which then gets min(3, 2) = 2. The pass is linear in the size of the cover.

## Slack and where it comes from

The arrival depth of a node is its depth in the current mapping, one more than its deepest LUT input. Its **slack** is required depth minus arrival depth. Nodes on a critical path have slack 0; a node off the critical paths, or any node when D exceeds the optimum, has positive slack. Setting D one above the optimal depth gives every node at least one level of slack, which often cuts the LUT count noticeably: a deliberate depth-area trade-off that designers request when timing has margin. Slack is a property of the current cover and must be recomputed after every change.

## Spending slack on cheaper cuts

In an area pass, visit nodes in topological order and give each the cut with the smallest area flow, or exact area, among those whose depth does not exceed its required depth. A cut of depth 5 is forbidden at a node required by depth 4 even if its area flow of 4.0 is the best; the best legal cut might have area flow 5.9. After each pass the cover is reselected and required depths recomputed, because changing one node's cut changes which nodes are needed. Several passes, area flow first and exact area after, typically converge in two or three iterations.

## Going deeper

Mishchenko, Chatterjee and Brayton, "Improvements to technology mapping for LUT-based FPGAs", IEEE TCAD 2007; Chen and Cong, "DAOmap: a depth-optimal area optimization mapping algorithm for FPGA designs", ICCAD 2004.
