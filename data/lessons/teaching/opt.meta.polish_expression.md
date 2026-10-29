# Slicing floorplans and annealing moves

A slicing floorplan recursively cuts the chip into rectangles, and a Polish expression writes that cutting tree in postfix form. Wong and Liu showed how to anneal over normalized Polish expressions with three simple moves, and their method remains the textbook example of a representation that makes annealing easy: every expression is a legal floorplan, area is computed by a linear-time stack evaluation, and each move is a small local edit.

## Evaluating a Polish expression

Operands are blocks with width and height. An operator combines the two sub-floorplans on top of the stack. With the convention used here, V places them side by side (widths add, height is the larger) and H stacks them (heights add, width is the larger). Rotations and soft blocks are handled by keeping shape curves, but fixed shapes show the idea.

Example: blocks 1 is 4 x 6, 2 is 5 x 5, 3 is 3 x 6, 4 is 3 x 3, 5 is 4 x 5, and expression 4 2 3 V 1 V 5 H H. Push 4, 2, 3. V combines 2 and 3 into 8 x 6. V with 1 gives 12 x 6. Push 5; H combines 12 x 6 with 5 to 12 x 11. H combines 4 (3 x 3) with that to 12 x 14. Area 168. The evaluation is one pass, linear in the number of blocks.

## Annealing moves on the expression

Wong and Liu's moves keep the expression valid. M1 swaps two adjacent operands, changing which blocks are neighbours. M2 complements a chain of consecutive operators, turning every H into V and vice versa within it. M3 swaps an adjacent operand and operator, which changes the tree shape and must preserve the balloting property (at every prefix, operands outnumber operators) and normalization (no two identical operators in a row).

Example: in 4 5 3 H H 2 1 H V, M1 swapping operands 3 and 2 gives 4 5 2 H H 3 1 H V. Re-evaluating with the stack gives the new area, and the annealer accepts or rejects by the Metropolis rule. Because a move edits only a short span, the bounding boxes above it in the tree are the only ones that change, which lets an implementation update area incrementally.

## Accepting a move and the representation's limits

After a move the cost change is $\Delta = A_{new} - A_{old}$ (plus any wirelength term), and an uphill move is accepted with probability $e^{-\Delta/T}$. If a move raises area from 180 to 195 at $T = 12$, it is accepted with probability $e^{-1.25} = 0.287$.

Normalized Polish expressions correspond one to one with slicing floorplans, so the annealer does not waste moves on duplicates. Their limit is that only slicing structures are representable; a pinwheel of five blocks is not a slicing floorplan. Non-slicing representations such as sequence pairs and B*-trees lift that limit at the cost of more expensive evaluation (longest paths in constraint graphs, or tree packing). The general lesson for engine design: choose a representation in which every point is legal and cheap to evaluate, then let the search operate on the representation.

## Going deeper

Wong and Liu, "A new algorithm for floorplan design" (DAC, 1986); Murata, Fujiyoshi, Nakatake and Kajitani, "Rectangle-packing-based module placement" (ICCAD, 1995), for sequence pairs; Chang, Chang, Wu and Wu, "B*-trees: a new representation for non-slicing floorplans" (DAC, 2000).
