# Structural pattern matching on a subject tree

Covering needs, at every subject node, the list of library patterns that match rooted there and the subject nodes each match leaves exposed as inputs. Structural matching computes this by walking a pattern and the subject tree together. It is the inner loop of tree covering, so its rules, and its blind spots, decide what the dynamic program can choose.

## Matching a pattern at a node

Pattern P matches rooted at subject node s if the two graphs agree node by node: a pattern INV must sit on a subject INV, a pattern NAND on a subject NAND, and a pattern leaf x may sit on any subject node, which then becomes an input of the match. Because NAND is commutative, its two pattern children may be paired with the subject children in either order. Example: NAND3's pattern NAND(x, INV(NAND(x, x))) matches at s = NAND(c, INV(NAND(a, b))) with leaves c, a, b. It does not match at NAND(c, NAND(a, b)), because the pattern needs an INV where the subject has a NAND.

## Counting matches

At every internal node INV or NAND2 matches trivially, since a single-node pattern only checks the node type. Larger cells match where the structure below allows. To count matches at a node, test each cell's pattern there; to count them over a tree, sum over all internal nodes. With a library of INV, NAND2, AND2 and NAND3, the node n = INV(NAND(a, b)) matches INV and AND2 (2 matches) and the node NAND(n, c) matches NAND2 and NAND3 (2 matches). The match lists are computed once, bottom-up, and stored for the covering step; with a fixed library each test is constant work.

## One match, several bindings

A pattern can match at the same node in more than one way. NAND3's pattern at NAND(INV(NAND(a, b)), INV(NAND(c, d))) binds either side as the INV(NAND) branch, giving leaves {INV(NAND(c, d)), a, b} or {INV(NAND(a, b)), c, d}. These bindings expose different subject nodes as inputs, so their total costs can differ, and a covering algorithm must consider each. A matcher that keeps only the first binding found can lose the optimum. The same issue appears at scale in Boolean matching, where one function can map onto a cell under several input permutations with different pin delays.

## Going deeper

Aho, Ganapathi and Tjiang, "Code generation using tree matching and dynamic programming", ACM TOPLAS 1989, the compiler ancestor of tree covering; Keutzer, "DAGON", DAC 1987.
