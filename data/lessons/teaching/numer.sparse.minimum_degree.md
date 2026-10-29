# Minimum degree ordering

Minimum degree is the symmetric version of Markowitz pivoting and the most widely used fill-reducing ordering. It eliminates the vertex with the fewest neighbours, updates the graph and repeats. Its approximate variant, AMD, orders matrices with millions of unknowns in seconds and is the default ordering inside KLU, CHOLMOD and many circuit simulators.

## Greedy elimination of the lowest-degree vertex

For a symmetric pattern, eliminating vertex $v$ of degree $d$ creates at most $d(d - 1)/2$ fill edges, so a small degree promises little fill. Minimum degree picks a vertex of least current degree, eliminates it, and repeats on the elimination graph. Ties need a rule; in these exercises the lowest label wins.

Example: edges 1-2, 1-3, 2-3, 2-4, 3-4, 4-5. Degrees: 1 has 2, 2 has 3, 3 has 3, 4 has 3, 5 has 1. Vertex 5 is eliminated first, at no cost. Then vertex 4 has degree 2 (neighbours 2 and 3) and ties with vertex 1, so vertex 1 goes next; its neighbours 2 and 3 are already joined, so again no fill.

## Tracing the ordering and its fill

Run the whole ordering: at each step record the chosen vertex and the fill edges its elimination creates. On a tree, minimum degree always eliminates a leaf (degree 1), so it creates no fill: trees factor perfectly. On a cycle of 6 vertices, the first elimination creates one fill edge, and so does each later one until three vertices remain, which form a clique: total fill 3.

Minimum degree is greedy, not optimal. There are graphs where it produces much more fill than the best order, and its result depends on tie-breaking. In practice it is excellent on irregular circuit graphs, while nested dissection does better on large regular meshes.

## Degree updates and approximate degrees

After eliminating $v$, a neighbour $u$ of $v$ becomes adjacent to all other neighbours of $v$. Its new degree is
$$\deg'(u) = |N(u) \cup N(v)| - 2,$$
subtracting $u$ itself and $v$. Example: $N(v) = \{a, u, w\}$ and $N(u) = \{v, x\}$. Then $N(u) \cup N(v) = \{a, u, v, w, x\}$, and $\deg'(u) = 5 - 2 = 3$ (neighbours $a, w, x$).

Recomputing exact degrees this way is the expensive part of minimum degree. AMD (Amestoy, Davis and Duff, 1996) represents eliminated vertices as elements in a quotient graph and computes cheap upper bounds on degree instead. The orderings are nearly as good, and the run time is close to linear in the number of nonzeros in practice.

## Going deeper

Amestoy, Davis and Duff, "An approximate minimum degree ordering algorithm", SIAM J. Matrix Anal. Appl. (1996); George and Liu, "The evolution of the minimum degree ordering algorithm", SIAM Review (1989).
