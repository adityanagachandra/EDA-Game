# Hash-consing and structural hashing

Hash-consing is the rule that a data structure never stores two copies of the same node: before creating a node, look up its key, the operator and its children's identifiers, in a unique table, and return the existing node if one is there. Applied to and-inverter graphs it is called structural hashing, the strash command in ABC; applied to BDDs it is the unique table that makes them canonical. It turns structural equality into pointer equality and keeps memory proportional to the number of distinct subterms.

## One node per distinct structure

Each node is keyed by its operator and the identifiers of its operands. An AIG has a single operator, two-input AND, with complement flags on edges, so a key is a pair of signed node identifiers. The builder computes the key, probes the hash table, and either returns the hit or allocates a new node and inserts it.

Build $g_1 = a \wedge b$, $g_2 = a \wedge b$, $g_3 = g_1 \wedge c$. The second call finds the key $(a, b)$ and returns $g_1$, so the table holds 2 nodes for 3 requests. Every later reference to $g_2$ is a reference to $g_1$; there is nothing to merge afterwards because the duplicate never existed.

## Canonical keys

Lookup only finds duplicates that produce the same key, so keys must be canonical. AND is commutative, so the operands are sorted by identifier (and polarity) before hashing: $b \wedge a$ and $a \wedge b$ get the same key. Duplicates also propagate: if $g_2$ was merged into $g_1$, then $g_4 = g_2 \wedge c$ is keyed through $g_1$ and hits $g_3$.

With $g_1 = a \wedge b$, $g_2 = b \wedge a$, $g_3 = g_1 \wedge c$ and $g_4 = c \wedge g_2$, sorted keys give $g_2 = g_1$ and then $g_4 = g_3$: 2 nodes for 4 definitions. Without canonical keys a commuted copy would survive and every algorithm that relies on pointer equality would treat two identical functions as different.

## Trivial simplification before hashing

Builders normalize before probing: $x \wedge x = x$, $x \wedge \overline{x} = 0$, $x \wedge 1 = x$ and $x \wedge 0 = 0$. Such gates never create nodes; they return an existing signal or a constant, and constants then propagate into the gates that use them.

So $g_5 = \overline{a} \wedge \overline{a}$ returns $\overline{a}$, and $g_6 = g_5 \wedge a$ becomes $\overline{a} \wedge a = 0$, and any later AND with $g_6$ is 0 too. Structural hashing is not functional equivalence: $a \wedge (b \wedge c)$ and $(a \wedge b) \wedge c$ are different structures. Detecting that needs simulation and SAT (FRAIGing), or a canonical form such as a BDD, where hash-consing plus reduction rules does make equal functions share one node.

## Going deeper

Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting", DAC (2006). Kuehlmann, Paruthi, Krohm and Ganai, "Robust Boolean reasoning for equivalence checking and functional property verification", IEEE TCAD (2002).
