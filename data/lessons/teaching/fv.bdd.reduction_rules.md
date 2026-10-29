# BDD reduction rules

A binary decision diagram starts life as a decision tree: test one variable per level, branch on its value, and read 0 or 1 at the leaves. Two reduction rules turn the tree into a compact, canonical graph: merge isomorphic subgraphs and delete tests whose two outcomes are the same. Knowing exactly what each rule removes is how a tool designer reasons about BDD size, memory and the cost of building BDDs for real logic.

## From decision tree to diagram

Fix a variable order, say $x_1 < x_2 < \dots < x_n$. The complete decision tree tests $x_1$ at the root, $x_2$ at both children, and so on, so level $i$ holds $2^i$ nodes for $i = 0, \dots, n - 1$ and the bottom holds $2^n$ leaves, one per assignment, each labelled with the function's value. The number of nonterminal (decision) nodes is $1 + 2 + \dots + 2^{n-1} = 2^n - 1$. For $n = 3$ that is 7 decision nodes and 8 leaves; for $n = 20$ it is over a million nodes, no matter how simple the function is. A tree is just a truth table drawn as a graph, so it is useless for large functions; all the compression comes from the reductions that follow. In this course node counts always mean nonterminal nodes, with the two terminals 0 and 1 not counted.

## Merging isomorphic subgraphs

The first rule merges any two nodes that test the same variable and have the same children, bottom-up; all leaves with the same value merge into one terminal. Applied alone, it gives the quasi-reduced BDD, in which every path still tests every variable. Its level $i$ has one node per distinct cofactor of $f$ by the first $i$ variables of the order. Take $f = (x_1 \wedge x_2) \vee x_3$ with order $x_1 < x_2 < x_3$. Level 0 has $f$ itself. Level 1 has the cofactors $f|_{x_1=0} = x_3$ and $f|_{x_1=1} = x_2 \vee x_3$: 2 nodes. Level 2 has the cofactors by $x_1$ and $x_2$: $x_3$, $x_3$, $x_3$ and 1, which are 2 distinct functions: 2 nodes. The quasi-reduced BDD has $1 + 2 + 2 = 5$ nonterminal nodes, against 7 in the tree.

## Removing redundant tests

The second rule deletes a node whose low and high children are the same node, redirecting its parents to that child: the test does not matter. In the quasi-reduced BDD of $(x_1 \wedge x_2) \vee x_3$, the level-1 node for $x_3$ tests $x_2$ but both its children are the $x_3$ node, and the level-2 node for the constant 1 tests $x_3$ with both children equal to the terminal 1. Both are removed, leaving 3 nodes, one per variable: the reduced ordered BDD (ROBDD). Applying both rules until neither applies gives a unique result for each function and order, which is Bryant's canonicity theorem: two functions are equal exactly when their ROBDDs are the same graph. So the redundant-test rule removed $5 - 3 = 2$ nodes. Packages never build the tree or the quasi-reduced form; they apply both rules on the fly through a unique table, so every node is created already reduced.

## Going deeper

Bryant, "Graph-Based Algorithms for Boolean Function Manipulation" (IEEE Transactions on Computers, 1986); Knuth, "The Art of Computer Programming", volume 4A, section 7.1.4.
