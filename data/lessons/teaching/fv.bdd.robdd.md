# Reduced ordered BDDs

A reduced ordered BDD represents a Boolean function canonically for a fixed variable order: equal functions give identical graphs, so equivalence is a pointer comparison. That property made BDDs the first practical engine for equivalence checking and symbolic model checking, and they remain the workhorse wherever many operations are done on the same functions. Counting ROBDD nodes for small functions is how you build the intuition for which functions are cheap and which are not.

## Counting nodes from a formula

The ROBDD has one node for each distinct subfunction that arises by cofactoring along the order and that actually depends on the variable being tested. To count, cofactor the function variable by variable in the order, collect the distinct non-constant results, and assign each to the first variable of the order it depends on. Take $f = x_1 \oplus x_2 \oplus x_3$ with order $x_1 < x_2 < x_3$. Testing $x_1$ needs 1 node. Its cofactors are $x_2 \oplus x_3$ and its complement, 2 distinct functions, each depending on $x_2$: 2 nodes. Their cofactors by $x_2$ are $x_3$ and $\neg x_3$: 2 nodes. Total 5. For $f = (x_1 \vee x_2) \wedge x_3$, the cofactors by $x_1$ are $x_2 \wedge x_3$ and $x_3$; the first needs an $x_2$ node whose high child is the same $x_3$ node, so the total is 3.

## Counting nodes from a truth table

When a function is given by its on-set, the list of inputs where it is 1, split the list by the first variable of the order, then by the next, and compare the pieces. Let $f$ be 1 exactly on $x_1x_2x_3 \in \{000, 011, 101, 110\}$ with order $x_1 < x_2 < x_3$. Splitting by $x_1$: for $x_1 = 0$ the remaining on-set over $x_2x_3$ is $\{00, 11\}$, which is XNOR; for $x_1 = 1$ it is $\{01, 10\}$, which is XOR. They differ, so the root needs 2 children, each a node on $x_2$. Splitting each by $x_2$ gives $x_3$ or $\neg x_3$, and both pieces share those two $x_3$ nodes. Total $1 + 2 + 2 = 5$; the function is even parity, which needs 5 nodes in every order. When a split produces two equal pieces, the test is redundant and that node does not exist.

## Shared BDDs and canonicity

A BDD package stores all functions in one shared graph with a unique table that maps (variable, low, high) to an existing node. Building a node first looks it up, so any subfunction that occurs twice, in one function or in two, is stored once. Store $f = x_1 \wedge x_2$ and $g = x_2$ with order $x_1 < x_2$. The ROBDD of $f$ is an $x_1$ node whose high child is the $x_2$ node; that $x_2$ node is exactly the ROBDD of $g$. The shared graph holds 2 nodes, not $2 + 1 = 3$. Canonicity makes equivalence free: after building two circuits' outputs, they are equivalent exactly when their roots are the same node. Packages such as CUDD also use complemented edges, so a function and its negation share one graph; with them, node counts differ from the plain counts used in this course.

## Going deeper

Bryant, "Graph-Based Algorithms for Boolean Function Manipulation" (IEEE Transactions on Computers, 1986); Somenzi, the CUDD package documentation; Knuth, "The Art of Computer Programming", volume 4A, section 7.1.4.
