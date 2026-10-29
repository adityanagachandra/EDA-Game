# Structural hashing of AIGs

And-inverter graphs (AIGs) represent logic with only two-input AND nodes and inverted edges. Structural hashing keeps them free of duplicate nodes: before creating an AND node, look up its pair of inputs in a hash table and reuse an existing node if there is one. It is the cheapest simplification in equivalence checking and synthesis, and in ABC it is the first thing done to any network, through the strash command.

## Hashing AND nodes

Each AIG node is $n = \text{AND}(p, q)$, where $p$ and $q$ are literals: a primary input or node, possibly complemented. Structural hashing normalizes the pair, for example by sorting the two literals, and uses it as a key. If the key is in the table, the new node is not created: its readers use the existing node. Otherwise a node is created and its key inserted. Take $n_1 = \text{AND}(a, b)$, $n_2 = \text{AND}(b, a)$, $n_3 = \text{AND}(n_1, \neg c)$ and $n_4 = \text{AND}(a, b)$. The key $\{a, b\}$ is created by $n_1$; $n_2$ and $n_4$ hit it and merge into $n_1$. Two AND nodes remain, $n_1$ and $n_3$. Complemented literals are different keys: $\text{AND}(a, \neg b)$ does not merge with $\text{AND}(a, b)$.

## Cascading merges

Hashing runs in topological order, and a merge can enable later merges. When a node is merged, every reader is redirected to its representative before the reader's own key is computed. Continue the example with $n_5 = \text{AND}(n_2, \neg c)$. Since $n_2$ was merged into $n_1$, the key of $n_5$ is $\{n_1, \neg c\}$, which is exactly the key of $n_3$, so $n_5$ merges too. A single pass in topological order catches every such cascade, because each node's inputs are final by the time it is processed. Count the result by walking the list in order, replacing inputs by their representatives, and adding a node only when its normalized key is new. Two circuits that are copies of each other up to input order of the gates collapse completely into one graph this way.

## Trivial simplifications

Hashing is usually combined with one-level rewriting rules applied before the lookup: $\text{AND}(x, x) = x$, $\text{AND}(x, \neg x) = 0$, $\text{AND}(x, 1) = x$ and $\text{AND}(x, 0) = 0$. A node removed this way becomes an alias for a literal or a constant, and its readers see that literal, which can trigger further simplifications and merges downstream. Take $n_1 = \text{AND}(a, b)$, $n_2 = \text{AND}(n_1, n_1)$, $n_3 = \text{AND}(n_2, \neg n_1)$ and $n_4 = \text{AND}(c, n_2)$. Then $n_2$ becomes $n_1$, $n_3$ becomes $\text{AND}(n_1, \neg n_1) = 0$, and $n_4$ becomes $\text{AND}(c, n_1)$, a new key: 2 AND nodes remain. Structural hashing only finds structural identity. Two nodes computing the same function with different structure, such as $\text{AND}(a, \text{AND}(b, c))$ and $\text{AND}(\text{AND}(a, b), c)$, stay distinct; finding those needs simulation and SAT sweeping.

## Going deeper

Kuehlmann, Paruthi, Krohm and Ganai, "Robust Boolean Reasoning for Equivalence Checking and Functional Property Verification" (IEEE TCAD, 2002); Mishchenko, Chatterjee and Brayton, "DAG-Aware AIG Rewriting" (DAC 2006); the ABC documentation for strash.
