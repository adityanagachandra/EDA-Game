# Slicing floorplans

A slicing floorplan is one that can be produced by repeatedly cutting a rectangle into two with a straight horizontal or vertical line. Its structure is a binary tree whose leaves are blocks and whose internal nodes are cuts, which makes area evaluation a single bottom-up pass and orientation optimisation polynomial. That convenience is why slicing structures were the first practical floorplan representation.

## Cuts and the slicing tree

Start with the chip rectangle and cut it with a vertical line into a left and right part, or with a horizontal line into a bottom and top part, then cut each part again, until each piece holds one block. The record of cuts is a slicing tree: each internal node is a V (vertical) or H (horizontal) cut, its children are the two sub-floorplans, and the leaves are blocks. A vertical cut puts children side by side; a horizontal cut stacks them. Not every floorplan is slicing: the smallest non-slicing arrangement is the wheel, or pinwheel, of five blocks spiralling around a centre, where no straight line crosses the whole chip without cutting a block.

## Evaluating the area bottom-up

For fixed block shapes, each node's dimensions follow from its children's. A vertical cut adds widths and takes the larger height: $(w_l + w_r, \max(h_l, h_r))$. A horizontal cut takes the larger width and adds heights: $(\max(w_l, w_r), h_l + h_r)$. A post-order traversal computes the root's dimensions in time linear in the number of blocks. Example: blocks A (5 by 6), B (3 by 9), C (3 by 10). Node r2 is a vertical cut of B and A: $(8, 9)$. The root is a horizontal cut with C below and r2 above: $(\max(3, 8), 10 + 9) = (8, 19)$, area 152. The blocks cover 87, so 65 um$^2$ is whitespace.

## Why slicing structures are useful and limited

Slicing trees make several hard problems easy. Area evaluation is linear time. Choosing the best orientation or shape of every block, given each block's candidate shapes, is polynomial for slicing floorplans (Stockmeyer's algorithm, by combining shape curves bottom-up), while it is NP-hard for general floorplans. Incremental evaluation after a local change only recomputes the path to the root. The cost is expressiveness: some optimal packings are non-slicing, and forcing a slicing structure can waste area, especially with many differently shaped blocks. Slicing structures also map naturally onto hierarchical channel routing, which is part of why they persisted.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 3.4; Stockmeyer, "Optimal orientations of cells in slicing floorplan designs" (Information and Control, 1983).
