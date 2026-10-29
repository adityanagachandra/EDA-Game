# Decision diagrams overview

A reduced ordered binary decision diagram (ROBDD) is the Shannon expansion tree of a function with identical subfunctions shared and pointless tests removed, under a fixed variable order. For a given order it is canonical: two functions are equal exactly when their diagrams are the same node. That property made BDDs the engine of symbolic model checking and early equivalence checking. Their weakness is size, which depends on the variable order and is exponential for some functions under every order. Formal verification studies the algorithms in depth; this concept is the representation itself.

## From the Shannon tree to a reduced diagram

Expand $f$ on $x_1$, then each cofactor on $x_2$, and so on, giving a binary tree whose leaves are constants. Two reduction rules shrink it: merge nodes that represent the same subfunction (hash-consing through a unique table), and delete any node whose two children are the same, since that test does not matter. Applying both until nothing changes gives the ROBDD.

For $f(a, b, c) = ab + \overline{c}$ with order $a < b < c$: the $a$ node has children $f_{\overline{a}} = \overline{c}$ and $f_a = b + \overline{c}$. The node for $b + \overline{c}$ tests $b$ and has children $\overline{c}$ and 1. The node for $\overline{c}$ is shared by both paths. The diagram has 3 internal nodes ($a$, $b$, $c$) and the two terminals.

## Counting reduced nodes level by level

The node count can be computed without drawing. At level $i$, fix the variables above it in every possible way and collect the distinct subfunctions obtained; the ones that depend on $x_i$ are exactly the nodes labeled $x_i$. Summing over levels gives the internal node count.

For the example: level $a$ has the one function $f$, which depends on $a$, 1 node. Level $b$ sees $\overline{c}$ and $b + \overline{c}$; only the second depends on $b$, 1 node. Level $c$ sees $\overline{c}$ and 1; only $\overline{c}$ depends on $c$, 1 node. Total 3. This count is the number of distinct cofactors that depend on the next variable, so it ties BDD size directly to Shannon expansion.

## The order decides the size

Take $f = a_1 b_1 + a_2 b_2 + a_3 b_3$. With the interleaved order $a_1 < b_1 < a_2 < b_2 < a_3 < b_3$ the diagram has 6 internal nodes, two per pair. With $a_1 < a_2 < a_3 < b_1 < b_2 < b_3$, after reading the $a$ variables the diagram must remember which subset of them is 1, and it grows to $2^{4} - 2 = 14$ internal nodes; for $n$ pairs it is $2n$ against $2^{n+1} - 2$.

Finding the best order is NP-hard, so packages such as CUDD reorder dynamically by sifting. Some functions defeat every order: Bryant proved that the middle output bit of an $n$-bit multiplier needs exponentially many nodes under any order. That limit is why SAT displaced BDDs in most equivalence checking, while BDDs remain valuable for small, structured functions and for counting and quantification. Variants trade canonicity rules for compactness: zero-suppressed diagrams for sparse sets, multi-terminal diagrams for integer-valued functions.

## Going deeper

Bryant, "Graph-Based Algorithms for Boolean Function Manipulation", IEEE Transactions on Computers (1986); Bryant, "On the complexity of VLSI implementations and graph representations of Boolean functions with application to integer multiplication", IEEE Transactions on Computers (1991). Knuth, "The Art of Computer Programming", volume 4A, section 7.1.4.
