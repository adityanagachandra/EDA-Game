# Tree covering by dynamic programming

On a tree, minimum-cost covering has optimal substructure: once a match is chosen at a node, the subtrees under its inputs are independent problems. DAGON exploits this with a bottom-up dynamic program that is exact for each tree and linear in its size. It is the reference algorithm every later mapper is compared with, and the template for delay-oriented and load-aware variants.

## The recurrence

Let cost(v) be the minimum area of a cover of the subtree rooted at v, with cost 0 for primary inputs. For an internal node, consider every match m rooted at v and every binding of it:
$$\text{cost}(v) = \min_{m} \Big( \text{area}(m) + \sum_{u \in \text{inputs}(m)} \text{cost}(u) \Big).$$
The inputs are the subject nodes bound to the pattern's leaves. Because the tree's subtrees under different inputs share no nodes, the minimum for each is chosen independently and the sum is exact. Record the winning match at each node so the cover can be rebuilt.

## Running it bottom-up

Visit nodes in post-order. Example with INV = 2, NAND2 = 3, AND2 = 4, NAND3 = 4: tree n1 = NAND(a, b), n2 = INV(n1), n3 = NAND(n2, c). cost(n1) = 3 (NAND2). cost(n2) = min(INV + cost(n1), AND2 + cost(a) + cost(b)) = min(5, 4) = 4. cost(n3) = min(NAND2 + cost(n2) + cost(c), NAND3 + cost(a) + cost(b) + cost(c)) = min(7, 4) = 4. The optimal cover is one NAND3. Then rebuild top-down: take the recorded match at the root and recurse into its inputs only, so n1 and n2 are never built as gates.

## Why greedy fails and what is exact

A greedy mapper that takes the biggest match at the root and recurses can lose: a large match may expose inputs whose own covers are expensive, while a smaller match leaves inputs that a complex cell covers cheaply. The DP never commits before seeing the cost below. Its optimality is narrow, though: it is exact for one tree, one subject graph and one pattern set, with costs that add. Cutting the DAG into trees, fixing the decomposition and ignoring load all happen outside the DP. Run time is linear in the tree size for a fixed library, which is why it remained practical for decades.

## Going deeper

Keutzer, "DAGON", DAC 1987; De Micheli, "Synthesis and Optimization of Digital Circuits", section 10.3; Cormen et al., "Introduction to Algorithms", the chapter on dynamic programming for the optimal-substructure argument.
