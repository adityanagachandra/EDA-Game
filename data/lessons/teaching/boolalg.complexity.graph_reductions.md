# Vertex cover, independent set and clique

Vertex cover, independent set and clique are three faces of one problem, and they appear throughout EDA: vertex cover in test-point selection and in breaking cycles, independent sets in conflict-free scheduling and in register binding, cliques in compatibility-based resource sharing. They are related by complementation, which turns any algorithm or hardness result for one into one for the others, and vertex cover is shown NP-hard by a classic gadget reduction from 3-SAT.

## Covers and independent sets are complements

A vertex cover is a set $C$ of vertices touching every edge. An independent set is a set $S$ with no edge inside it. A set $S$ is independent exactly when $V \setminus S$ is a vertex cover: an edge inside $S$ is precisely an edge the complement misses. So the maximum independent set has size $|V|$ minus the minimum vertex cover.

A graph with 26 vertices whose maximum independent set has 2 vertices has a minimum vertex cover of 24. The equivalence is exact for optimal solutions but not for approximations: a cover within a factor 2 of optimal gives an independent set with no comparable guarantee, which is why vertex cover is approximable and independent set is not.

## 3-SAT to vertex cover

The Garey-Johnson reduction builds, for a 3-CNF with $n$ variables and $m$ clauses, one edge per variable between vertices $x_i$ and $\neg x_i$, one triangle per clause, and an edge from each triangle corner to the literal vertex it represents. The graph has $2n + 3m$ vertices and $n + 6m$ edges (one per variable, three triangle edges and three links per clause), and the target is $k = n + 2m$.

Any cover needs at least one endpoint per variable edge and at least two corners per triangle, so $n + 2m$ is the minimum possible. A cover of exactly that size picks one literal per variable (an assignment) and leaves one corner per triangle uncovered by the triangle itself; that corner's link must be covered by its literal vertex, so the literal is true. For $n = 2$ and $m = 1$: 7 vertices, 8 edges, $k = 4$.

## Clique through the complement graph

A clique in $G$ is an independent set in the complement graph $\overline{G}$, which has the same vertices and exactly the edges $G$ lacks: $\binom{|V|}{2} - |E|$ of them. A graph with 5 vertices and 4 edges has a complement with $10 - 4 = 6$ edges. Building the complement takes polynomial time, so clique and independent set are polynomially equivalent, and both are NP-complete.

Complementation can blow up a sparse graph into a dense one, which matters in practice even though it is polynomial: a netlist conflict graph with $10^6$ vertices and a few million edges has a complement with about $5 \times 10^{11}$ edges. Tools therefore choose the formulation whose graph is sparse, and solve maximum independent set on the conflict graph rather than maximum clique on the compatibility graph, or vice versa.

## Going deeper

Garey and Johnson, "Computers and Intractability" (1979), section 3.1.3. Karp, "Reducibility among combinatorial problems" (1972).
