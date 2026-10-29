# The Laplacian system of quadratic placement

The linear system of quadratic placement has a precise structure: a weighted graph Laplacian over the movable cells plus diagonal terms from fixed connections. Knowing how to assemble it entry by entry, and why it is symmetric positive definite, is what a placer implementer needs to build the matrix, choose a solver and debug a singular system.

## Diagonal entries

Write the x cost as $\frac{1}{2}x^T A x - b^T x + \text{const}$. For movable cell $i$, the diagonal entry $A_{ii}$ is the total weight of all connections at $i$: to other movable cells and to fixed pins alike. If m1 connects to m2 with weight 2, to m3 with weight 3, and to fixed pins with weights 1 and 3, then $A_{11} = 2 + 3 + 1 + 3 = 9$. (Conventions differ by a factor of 2 depending on whether the cost is written with a $\frac{1}{2}$; the optimal x does not change.) A cell connected to nothing has a zero row, which makes the system singular, so such cells are removed or anchored.

## Off-diagonal entries and the right-hand side

For two movable cells $i \ne j$, $A_{ij} = -w_{ij}$, minus the total weight of connections between them, and zero if they are unconnected, which is why A is sparse: each row has as many nonzeros as the cell has movable neighbours. The right-hand side collects fixed pins: $b_i = \sum_f w_{if} x_f$ over the fixed connections of $i$. If m1 connects to fixed pins at x = 21 (weight 1) and x = 13 (weight 2), then $b_1 = 21 + 26 = 47$. The movable-to-movable part of A is exactly the weighted graph Laplacian of the movable cells; the fixed connections add a nonnegative diagonal matrix.

## Positive definiteness and what breaks it

A graph Laplacian is positive semidefinite with the all-ones vector in its null space: shifting every cell by the same amount does not change any distance. Adding positive diagonal terms from fixed connections removes that null space, provided every connected component of the movable graph touches at least one fixed pin. Then A is symmetric positive definite and the system has a unique solution, which conjugate gradient finds. If some component floats free, the matrix is singular and the solver drifts or fails; placers detect this and add weak anchors. A also tends to be ill-conditioned, because weakly anchored cells give small eigenvalues, which slows conjugate gradient and motivates preconditioning.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 4.3; Spielman's lecture notes on spectral graph theory, for Laplacian properties.
