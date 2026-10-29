# E-graphs and equality saturation

Destructive rewriting commits to each rewrite as it is applied, so the order of rewrites decides what can be found later, the phase-ordering problem. Equality saturation instead keeps every version of the expression in an e-graph, a compact data structure for sets of equivalent terms, applies all rewrites until nothing new appears or a budget runs out, and extracts the best term at the end. It is increasingly used for datapath optimization.

## E-classes and e-nodes

An **e-graph** consists of e-classes, each a set of equivalent e-nodes. An e-node is an operator applied to child e-classes, not child terms: the e-node $C_0 \ll C_5$ means "an element of class $C_0$ shifted by an element of class $C_5$". Because children are classes, one e-graph represents exponentially many terms: choosing one e-node in each relevant class spells out one term. Two invariants are maintained. Hashconsing: no two identical e-nodes exist. Congruence: if two classes are merged, parent e-nodes that become identical are merged too. Union-find tracks which classes are the same.

## Saturation

Each iteration searches for matches of every rewrite rule's left-hand side, adds the corresponding right-hand side, and merges its class with the matched class; then it rebuilds congruence. Nothing is deleted. With rules $x \cdot 2 \to x \ll 1$ and $x \ll 1 \to x + x$, the class of $a \cdot 2$ ends up containing $a \cdot 2$, $a \ll 1$ and $a + a$. If an iteration adds nothing, the e-graph is **saturated**: it contains every term reachable by the rules. Rules such as associativity and commutativity can grow the e-graph exponentially, so runs are bounded by node count, iterations or time, and rule schedulers throttle explosive rules.

## Extraction and soundness

After saturation, **extraction** chooses one e-node per needed class to minimize a cost. With tree cost (each use counted separately) a bottom-up dynamic program is exact; with DAG cost (shared subterms counted once) extraction is NP-hard and is solved with integer programming or heuristics. Correctness rests entirely on the rules: a single unsound rule makes unequal terms equivalent, and every term in the affected classes is suspect. Hardware rules need bit-width and signedness conditions, because modular arithmetic breaks identities that hold over integers. Extracted designs are therefore checked by equivalence checking against the original.

## Going deeper

Tate, Stepp, Tatlock and Lerner, "Equality saturation: a new approach to optimization", POPL (2009); Willsey et al., "egg: fast and extensible equality saturation", POPL (2021).
