# Parameterized complexity

Parameterized complexity measures running time in two quantities: the input size $n$ and a parameter $k$ such as the solution size or the treewidth. A problem is fixed-parameter tractable (FPT) when it can be solved in time $f(k) \cdot n^{O(1)}$, exponential only in the parameter. Many NP-hard EDA problems have small natural parameters, few faults to explain, few cuts, small treewidth of a netlist, so an FPT algorithm can be exact and fast on the instances that matter. Vertex cover parameterized by cover size is the canonical example.

## Bounded search trees

To decide whether a graph has a vertex cover of size at most $k$, pick any uncovered edge $(u, v)$. Some endpoint must be in the cover, so branch: take $u$ or take $v$, each lowering the budget by 1. The search tree has depth at most $k$ and at most $2^k$ leaves, with polynomial work per node, so the running time is $O(2^k n)$: for $k = 20$ that is about a million leaves, regardless of whether $n$ is a thousand or a million.

Smarter branching lowers the base. Branching on a vertex $v$ of degree at least 3, either take $v$ or take all of its neighbours, gives leaf counts $L(k) = L(k-1) + L(k-3)$, which grows like $1.4656^k$; the best known algorithms reach about $1.2738^k$. By contrast, clique parameterized by $k$ is W[1]-complete and believed not to be FPT, so the parameter alone does not make a problem easy.

## Kernelization: forced vertices

A kernel is a polynomial-time reduction to an equivalent instance whose size depends only on $k$. Buss's rule for vertex cover: a vertex of degree greater than $k$ must be in any cover of size at most $k$, because otherwise all of its more than $k$ neighbours would be. So take it, delete its edges, and lower the budget to $k - 1$. Repeat with the new budget until no vertex qualifies.

On a star with center 1 and 8 leaves plus one extra edge (3,4), with $k = 2$, vertex 1 has degree 8 > 2 and is forced. The budget drops to 1, and the remaining edge (3,4) has endpoints of degree 1, so the rule stops after forcing one vertex.

## The kernel size bound

After the rule stops with budget $k'$, every vertex has degree at most $k'$, so $k'$ vertices can cover at most $k'^2$ edges. If more than $k'^2$ edges remain, the answer is no immediately; otherwise the remaining graph has at most $k'^2$ edges and is solved by the search tree.

In the example $k' = 1$ and one edge remains, within the bound $1^2 = 1$, and taking vertex 3 finishes the cover $\{1, 3\}$. Kernelization is preprocessing with a guarantee, the theoretical version of what SAT preprocessors and logic-synthesis sweeps do in practice: shrink the instance by rules that are always safe before spending exponential effort.

## Going deeper

Cygan, Fomin, Kowalik, Lokshtanov, Marx, Pilipczuk, Pilipczuk and Saurabh, "Parameterized Algorithms" (2015), chapters 2 and 3. Chen, Kanj and Xia, "Improved upper bounds for vertex cover", Theoretical Computer Science (2010).
