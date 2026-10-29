# The graph Laplacian

The Laplacian matrix $L = D - A$ turns a weighted netlist graph into linear algebra. Its quadratic form measures cut size and quadratic wirelength, its second eigenvector drives spectral partitioning, and its linear systems are what analytical placers solve millions of times. For a tool designer, the Laplacian is the bridge between connectivity and numerical optimization, and building it correctly from nets is the first step.

## Degree minus adjacency

For a weighted undirected graph, the adjacency matrix $A$ has $A_{ij} = w_{ij}$, the weight of edge $i$-$j$ (0 if absent), and the degree matrix $D$ is diagonal with $D_{ii} = \sum_j w_{ij}$, the weighted degree. The Laplacian is $L = D - A$: weighted degrees on the diagonal, minus the edge weights off it.

For edges $v_1$-$v_3$ (6), $v_2$-$v_3$ (6) and $v_3$-$v_4$ (6), the diagonal entry $L_{33}$ is $6 + 6 + 6 = 18$, and $L_{34} = -6$. Each row sums to zero, so $L\mathbf{1} = 0$. In a netlist, $L$ is as sparse as the graph, with one off-diagonal pair per edge.

## Properties that tools rely on

$L$ is symmetric and positive semidefinite, because $x^T L x = \sum_{(i,j)} w_{ij}(x_i - x_j)^2 \ge 0$ for every vector $x$. Its smallest eigenvalue is 0, with eigenvector $\mathbf{1}$, and the multiplicity of 0 equals the number of connected components. So a partitioner can detect a disconnected netlist from the spectrum, and a placer knows that $L$ alone is singular.

That singularity is why quadratic placement needs fixed pins: with some vertices fixed, the remaining block of $L$ restricted to movable cells is positive definite (for a connected graph), so the linear system has a unique solution and conjugate gradient converges. The diagonal dominance of $L$ also makes simple preconditioners such as Jacobi effective.

## Laplacians of netlists

A netlist has nets, not edges, so its Laplacian depends on the net model. Under the clique model with edge weight $2/k$, a cell on a $k$-pin net gets $k - 1$ clique edges of weight $2/k$, so the net adds $2(k-1)/k$ to the cell's diagonal entry. A cell on one 4-pin net and nothing else has diagonal entry $3 \cdot 2/4 = 1.5$.

Changing the model changes $L$: under weight $1/(k-1)$ each net adds exactly 1 to the diagonal of each of its cells, so the diagonal is just the cell degree. Star models add a vertex per net, giving a larger but sparser Laplacian. Placers rebuild $L$ every iteration under the Bound2Bound model, whose weights depend on the current positions, so building $L$ fast from the compressed netlist matters.

## Going deeper

Chung, "Spectral Graph Theory" (1997), chapter 1; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 4 on quadratic placement.
