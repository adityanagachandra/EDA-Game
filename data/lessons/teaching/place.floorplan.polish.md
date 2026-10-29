# Normalized Polish expressions

A Polish expression is the post-order traversal of a slicing tree, written as a string of block names and cut operators. Wong and Liu showed that normalized Polish expressions correspond one-to-one with slicing floorplans, giving simulated annealing a compact, unambiguous search space with cheap moves. Evaluating one is a stack computation that every floorplanner of this family runs millions of times.

## From tree to postfix string

Write a slicing tree in post-order: left subtree, right subtree, then the operator. Operands are blocks; operators are V (side by side, left operand on the left) and H (stacked, left operand below). The tree with root H, left child block 1 and right child a V cut of blocks 2 and 3 becomes 1 2 3 V H. Every valid expression has the balloting property: reading left to right, every prefix contains more operands than operators, and the whole string has $n$ operands and $n - 1$ operators. A slicing floorplan can have several trees, because a run of the same cut can be grouped in different ways, so several expressions can describe one floorplan.

## Evaluating with a stack

Read the expression left to right. Push each block's shape. On an operator, pop the right operand, then the left, combine them (V adds widths and takes the larger height; H takes the larger width and adds heights) and push the result. The last shape on the stack is the floorplan. For 1 2 3 V H with blocks 1 = 7 by 6, 2 = 7 by 2 and 3 = 10 by 9: push 1, 2 and 3; V pops 3 and 2 and pushes (17, 9); H pops (17, 9) and 1 and pushes $(\max(7, 17), 6 + 9) = (17, 15)$. The area is 255, and with 146 um$^2$ of blocks, 109 um$^2$ is whitespace. The cost is linear in the expression length.

## Normalization makes the encoding unique

A Polish expression is normalized if it has no two identical operators in a row, such as V V or H H. Normalized expressions correspond to skewed slicing trees, where no node has a right child with the same operator, and Wong and Liu proved that each slicing floorplan has exactly one normalized expression. That removes redundancy from the search space: an annealer exploring normalized expressions never wastes moves on different strings for the same floorplan. The moves M1, M2 and M3 are designed to keep expressions normalized and balloting, or to be checked and rejected when they would not. The number of distinct slicing structures still grows exponentially with $n$, so search rather than enumeration is required.

## Going deeper

Wong and Liu, "A new algorithm for floorplan design" (DAC 1986); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 3.5.
