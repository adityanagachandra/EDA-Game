# The theory of arrays

Memories, register files and FIFOs are arrays: maps from addresses to data, read and written each cycle. The SMT theory of arrays models them with two operations, select and store, and a handful of axioms. It lets a verifier reason about a memory through the addresses actually used rather than its full size, which is the only way large memories fit into a formal proof.

## Read over write

select$(a, i)$ reads array $a$ at index $i$; store$(a, i, v)$ is the array that equals $a$ except that index $i$ holds $v$. The read-over-write axiom says select(store$(a, i, v)$, $j$) is $v$ if $i = j$ and select$(a, j)$ otherwise. With concrete addresses, peel the stores from the outside in. For $b$ = store(store(store$(a, 4, 92)$, 3, 69), 2, 14): select$(b, 3)$ skips the outer store (address 2), hits the store at address 3 and returns 69. select$(b, 4)$ returns 92. select$(b, 6)$ misses all three stores and reaches the initial array: the initial content $a[6]$, which is unknown, not 0. When the same address is written twice, the later (outer) store wins, because it is peeled first.

## Symbolic addresses

With symbolic indices the comparison may be unknown, and read-over-write produces a case split. Consider select(store(store$(a, i, v)$, $j$, $w$), $k$). Compare $k$ with the outer index $j$ first: if $k = j$ is known, the result is $w$. If $k \neq j$ is known, move inward and compare $k$ with $i$: $k = i$ gives $v$, $k \neq i$ gives select$(a, k)$. If $k \neq j$ is known but nothing is known about $k$ and $i$, the result is ite$(k = i, v, \text{select}(a, k))$. If $k = i$ is known but $k$ and $j$ are unrelated, the result is ite$(k = j, w, v)$, because a matching outer store still overrides the inner one. With nothing known the term is ite$(k = j, w, \text{ite}(k = i, v, \text{select}(a, k)))$. Each ite is a case the solver must explore, so the number of stores in a chain drives the cost.

## Arrays as memories in hardware proofs

The extensionality axiom states that two arrays equal at every index are equal, which lets a solver compare whole memories. Quantifier-free formulas over select and store are decidable: read-over-write case splits on the finitely many index terms that appear reduce them to equality reasoning, so the cost depends on the accesses, not on the memory size. In BMC of a design with a memory, each cycle's write adds a store, and each read becomes a select over the chain of earlier stores; word-level SMT encodings can keep a memory as one array term instead of bit-blasting every word. A further abstraction tracks only a few addresses and returns arbitrary data for the rest. That over-approximates the memory, so proofs carry over but counterexamples can be spurious and must be checked.

## Going deeper

McCarthy, "Towards a Mathematical Science of Computation" (IFIP Congress, 1962); Bradley, Manna and Sipma, "What's Decidable About Arrays?" (VMCAI 2006); Kroening and Strichman, "Decision Procedures" (2nd ed., 2016), chapter 7.
