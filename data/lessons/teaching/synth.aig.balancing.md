# AIG balancing

Balancing reduces AIG depth by exploiting associativity. A chain of ANDs, $((ab)c)d$, has depth 3; the tree $(ab)(cd)$ has depth 2 and the same function. ABC's balance command finds multi-input AND groups (supergates) and rebuilds each as a minimum-depth tree, taking the arrival levels of its inputs into account.

## Supergates and balanced trees

A **supergate** is a maximal tree of AND nodes connected by uncomplemented edges, each internal node having a single fanout. Collecting its leaves gives a multi-input AND, for example $x_1 x_2 \dots x_k$. AND is associative and commutative, so any binary tree over these leaves computes the same function. With all leaves at level 0, a chain has depth $k - 1$ and a balanced tree has depth $\lceil \log_2 k \rceil$: for $k = 5$, depth 4 becomes 3; for $k = 8$, depth 7 becomes 3. The node count stays $k - 1$, so balancing changes depth without changing area, as long as no shared internal node must be duplicated.

## Balancing with arrival levels

Leaves arrive at different levels. The minimum-depth tree is built greedily: repeatedly take the two leaves with the lowest levels, AND them, and insert the new node with level equal to their maximum plus one, until one node remains. This is Huffman's construction with $\max + 1$ in place of the sum, and it is optimal. For levels 2, 2, 4, 3, 3, 3, 0: combine 0 and 2 into 3, then 2 and 3 into 4, then 3 and 3 into 4, then 3 and 4 into 5, then 4 and 4 into 5, then 5 and 5 into 6. The output is at level 6. A late-arriving input ends up near the root, close to the output.

## Depth reduction and its costs

For a chain fed in the order given, with levels 2, 4, 1, 2, 2, the output is at $\max(\max(\max(\max(2,4)+1,1)+1,2)+1,2)+1 = 8$; balancing gives 5, a reduction of 3. Balancing is applied in topological order, so improving a supergate's output level also improves everything above it. It is cheap, linear time with a priority queue per supergate, and scripts call it between area-oriented passes because rewriting and refactoring can deepen the graph. Its limits: it only uses associativity and commutativity of AND, so it cannot shorten paths that need real restructuring, and where internal nodes are shared it either stops at them or duplicates logic.

## Going deeper

The ABC documentation for balance; Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting", DAC (2006), which uses balancing in its scripts; Huffman, "A method for the construction of minimum-redundancy codes", Proceedings of the IRE (1952).
