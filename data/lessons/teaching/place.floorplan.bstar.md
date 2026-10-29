# B*-trees and contour packing

A B*-tree encodes a compacted floorplan as an ordered binary tree: left children sit to the right of their parent, right children sit above it at the same x. Packing is a single pre-order traversal with a contour that tracks the skyline of placed blocks. Its linear-time packing and simple moves made it a popular representation for annealing floorplanners and for analog placement with symmetry constraints.

## Left children right, right children up

The root block sits at the origin. For a node b with parent p: if b is p's left child, b is placed immediately to the right of p, at $x_b = x_p + w_p$; if b is p's right child, b is placed above p at the same x, $x_b = x_p$. So the x coordinates follow directly from the tree: walk from the root and add the parent's width on each left edge. With A = 2 by 4 at the root, B as A's left child (x = 2), D as B's left child (x = 2 + 5 = 7), and C as B's right child (x = 2), every block's x is known before any y is computed. The tree represents admissible placements, those compacted to the left and bottom.

## Dropping blocks onto a contour

The y coordinate is found by dropping the block from above onto what is already placed: $y_b$ is the highest top edge among placed blocks that overlap b's x interval. Blocks are placed in pre-order (node, left subtree, right subtree), which guarantees that the blocks a node can rest on have already been placed. A contour data structure, a doubly linked list of horizontal segments describing the current skyline, answers each query and is updated by splicing in the new block's top edge. With amortised constant work per block, the whole packing is $O(n)$. Example: A = 8 by 2 at the origin, its left child B = 2 by 5 at x = 8, and its right child C = 4 by 1 at x = 0, which lands on A at y = 2.

## Moves and uses

Annealing over B*-trees uses three moves: rotate a block, move a node to another position in the tree, and swap two nodes. Each move changes only part of the packing, though a fresh pre-order pass is simple and linear. Compared with sequence pairs, B*-trees represent only admissible placements, compacted left and down, which loses nothing for area minimisation; compared with Polish expressions, they also represent non-slicing floorplans. Extensions handle symmetry groups for analog layout, where mirrored device pairs must share an axis, and fixed-outline constraints. Because packing is cheap, B*-tree annealers can run millions of moves on designs with tens to hundreds of blocks.

## Going deeper

Chang, Chang, Wu and Wu, "B*-trees: a new representation for non-slicing floorplans" (DAC 2000); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 3.
