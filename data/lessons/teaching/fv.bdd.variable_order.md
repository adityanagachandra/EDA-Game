# BDD variable ordering

The same function can have a BDD with a handful of nodes under one variable order and an exponential number under another. Choosing the order is therefore the main lever a BDD-based tool has, and it is a hard optimization problem. The classic example, a sum of products of variable pairs, shows the gap exactly and explains the rule of thumb that related variables belong close together in the order.

## A function that rewards interleaving

Let $f = x_1 y_1 + x_2 y_2 + \dots + x_n y_n$. Under the interleaved order $x_1 < y_1 < x_2 < y_2 < \dots$, the BDD tests $x_1$; if it is 0, the first product is false and the path moves on to the next pair; if it is 1, it tests $y_1$, which either makes $f$ true or moves on to the next pair. At every point the graph only needs to remember which pair it is on, so each pair contributes one $x_i$ node and one $y_i$ node: $2n$ nonterminal nodes in total. For $n = 3$ that is 6 nodes, and for $n = 100$ it is 200. The width of every level is 1. Information flows through the order locally, which is the property good orders have in common.

## The same function under a bad order

Under the separated order $x_1 < \dots < x_n < y_1 < \dots < y_n$, the BDD must read all the $x$ variables before any $y$. When it reaches $y_1$, the remaining function is $\bigvee_{i \in S} y_i$, where $S$ is the set of indices with $x_i = 1$, and different sets give different functions. So the $x$ levels form a complete tree: level $i$ needs $2^i$ nodes, $2^n - 1$ in all. The $y$ levels need one node for each distinct function $\bigvee_{i \in S'} y_i$ with $S'$ nonempty, organized so that level $y_j$ holds $2^{n-j}$ nodes, another $2^n - 1$. The total is $2^{n+1} - 2$: 14 for $n = 3$, against 6 interleaved, and over $2 \cdot 10^{30}$ for $n = 100$. The function did not change; only the amount of state the order forces the diagram to carry did.

## Comparing orders and finding good ones

For a general function, compare orders by building the ROBDD under each and counting nodes, and expect differences in both directions. Finding an optimal order is NP-hard, and even improving a given order is NP-complete, so packages use heuristics. Static heuristics order variables from the circuit structure, for example by a depth-first traversal from the outputs so that inputs feeding the same gates sit together. Dynamic reordering by sifting, due to Rudell, moves one variable at a time through every position, keeps the best, and repeats; swapping two adjacent levels is a local operation, which makes sifting affordable. Reordering is triggered when the node count grows past a threshold. A tool designer should remember that some functions, such as the middle bits of a multiplier, have exponential BDDs under every order, so no reordering can save them.

## Going deeper

Bryant, "Graph-Based Algorithms for Boolean Function Manipulation" (1986); Rudell, "Dynamic Variable Ordering for Ordered Binary Decision Diagrams" (ICCAD 1993); Bollig and Wegener, "Improving the Variable Ordering of OBDDs Is NP-Complete" (IEEE Transactions on Computers, 1996).
