# Hierarchy, top-down and bottom-up design

A chip is never designed as one flat list of billions of cells. It is a tree of modules: the top instantiates subsystems, which instantiate blocks, which instantiate library cells. Hierarchy lets a team divide the work, lets a tool analyze a block once and reuse the result, and is the structure that top-down planning and bottom-up assembly operate on. Its cost is that optimization across a boundary becomes harder.

## Definitions, instances and flattening

A module definition is written once. An instance is one use of it. The instantiation structure is a rooted tree (more exactly a DAG of definitions, unrolled into a tree of instances). Flattening replaces every instance by a copy of its contents until only leaf cells remain. To count the flattened cells, multiply along each root-to-leaf path. If the top holds 4 instances of block A and 3 of block B, A holds 10 cells of its own plus 2 instances of C, C holds 25 cells and B holds 40 cells, then the flattened design has $4(10 + 2 \cdot 25) + 3 \cdot 40 = 360$ cells.

## Top-down and bottom-up

In top-down design the architect splits the system into blocks, gives each block a budget (area, timing at its ports, power) and refines each block independently. In bottom-up design, proven blocks (cells, memories, IP cores) are built and characterized first and then assembled. Real projects do both: the floorplan and budgets come top-down, while library cells, memories and hard IP come bottom-up. The meeting point is the block boundary, where top-down budgets must match bottom-up characterizations. A budget that the block cannot meet sends the plan back up the tree.

## What hierarchy buys and costs

A hierarchical tool stores and analyzes each definition once. In the example, the definitions hold $10 + 25 + 40 = 75$ cells against 360 flattened, almost five times fewer. Timing tools use extracted block models (interface logic models, extracted timing models) for the same reason. The cost is lost optimization: a synthesizer that keeps boundaries cannot merge logic across them, and a placer that keeps a block in a fixed region cannot interleave it with its neighbours. Flows therefore flatten selectively: small blocks are flattened for quality, large ones kept hierarchical for capacity and team ownership.

## Going deeper

Weste and Harris, "CMOS VLSI Design: A Circuits and Systems Perspective" (4th ed., 2010), chapter 1 on hierarchy, regularity, modularity and locality.
