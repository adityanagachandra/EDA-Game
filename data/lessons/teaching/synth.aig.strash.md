# Structural hashing

Structural hashing makes an AIG canonical up to structure: no two AND nodes have the same pair of fanins. Every node creation goes through a hash table, so duplicate logic is merged the moment it is built, and trivial cases such as $x \cdot x$ and $x \cdot \bar{x}$ are folded. It is the reason building, rewriting and copying AIGs is cheap, and the reason they stay compact.

## The hash table on fanin pairs

To create $AND(a, b)$, first normalize the pair of literals by sorting, since AND is commutative, then look it up in a hash table keyed by the pair. If an entry exists, return the existing node; otherwise create a node and insert it. Given $n_1 = AND(\bar{a}, d)$ and then $n_2 = AND(d, \bar{a})$, the second lookup finds $n_1$, so only one node exists. This makes the AIG **structurally hashed**: every node is unique for its fanins. ABC's strash command rebuilds a network through this table, and every later transformation creates nodes the same way, so the invariant is maintained everywhere.

## Cascading merges

Process nodes in topological order and translate each fanin through the map of earlier merges before hashing. If $n_2$ merged into $n_1$, then a later node $AND(n_2, c)$ becomes $AND(n_1, c)$ and may itself match an existing $AND(n_1, c)$. In the list $n_1 = AND(a, b)$, $n_2 = AND(n_1, a)$, $n_3 = AND(b, a)$, $n_4 = AND(n_2, \bar{d})$, $n_5 = AND(a, n_1)$, $n_6 = AND(n_1, a)$, $n_7 = AND(\bar{d}, n_2)$: $n_3$ merges with $n_1$, $n_5$ and $n_6$ with $n_2$, and $n_7$ with $n_4$, leaving 3 nodes. One duplicate deep in a cone can collapse whole cones above it.

## Trivial rules and what hashing cannot see

Node creation also applies one-level simplifications: $x \cdot x = x$, $x \cdot \bar{x} = 0$, $x \cdot 0 = 0$ and $x \cdot 1 = x$. These return an existing literal or a constant instead of creating a node, and constants then propagate through later nodes. Structural hashing is not functional canonicity: $a(b + c)$ and $ab + ac$ compute the same function with different structures and both survive. Finding such equivalences needs simulation and SAT, which is what FRAIGing adds. Hashing keeps the AIG free of exact structural duplicates at constant cost per node, and it makes copying a network through a transformation automatically remove the duplicates the transformation creates.

## Going deeper

Kuehlmann and Krohm, "Equivalence checking using cuts and heaps", DAC (1997); Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting: a fresh look at combinational logic synthesis", DAC (2006).
