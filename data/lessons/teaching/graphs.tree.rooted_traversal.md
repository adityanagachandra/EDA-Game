# Pre-order and post-order traversal of a rooted tree

Routed nets, clock trees, Steiner trees and buffer trees are rooted trees: a driver at the root, sinks at the leaves. Almost every computation on them is a traversal. Post-order (children before parent) accumulates quantities that flow upward, such as downstream capacitance; pre-order (parent before children) pushes quantities downward, such as delay from the root. Elmore delay, van Ginneken buffering and clock skew analysis are all built from these two passes.

## Pre-order: parent before children

In a rooted tree every node except the root has one parent. A **pre-order** traversal visits a node, then recursively visits each child subtree in a fixed child order. The root comes first, and each node appears before all of its descendants.

Take root 1 with children 2 and 3, node 2 with children 4 and 5, and node 4 with children 6 and 7, visiting children in increasing label order. Pre-order is 1, 2, 4, 6, 7, 5, 3: from 2 we finish its whole subtree (4, 6, 7, then 5) before reaching 3. The sixth node visited is 5. Pre-order is the right order whenever a node's value depends on its parent's, because the parent is always finished first.

## Post-order: children before parent

A **post-order** traversal recursively visits every child subtree first and visits the node last. Leaves come first, the root last, and every node appears after all of its descendants.

For the same tree, post-order is 6, 7, 4, 5, 2, 3, 1. Node 4 waits for 6 and 7, node 2 waits for 4 and 5, and the root waits for everything. The fourth node visited is 5. Post-order is the right order whenever a node's value is a combination of its children's values, because those values are ready when the node is reached. Both traversals take time linear in the number of nodes, and an explicit stack replaces recursion for trees deep enough to overflow the call stack.

## Two passes on an RC tree

Elmore delay on an RC tree uses both passes. A post-order pass computes the downstream capacitance $C_{down}(v) = C_v + \sum_{c \in \text{children}(v)} C_{down}(c)$. A pre-order pass then computes delay $T(v) = T(\text{parent}(v)) + R_{edge(v)}\,C_{down}(v)$.

With capacitances $C_2 = 9$, $C_4 = 8$ and $C_6 = 3$ fF, where 4 and 6 are the children of 2 and both are leaves, $C_{down}(2) = 9 + 8 + 3 = 20$ fF. The same pattern appears in van Ginneken's buffer insertion (candidate solutions merge bottom-up, then the chosen buffers are placed top-down) and in clock-tree synthesis (subtree delays merge bottom-up, then the tree is embedded top-down). Recognizing which quantity flows which way is the design step.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), chapter 12; Rubinstein, Penfield and Horowitz, "Signal delay in RC tree networks" (IEEE TCAD, 1983).
