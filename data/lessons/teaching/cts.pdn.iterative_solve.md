# Solving large power grids iteratively

A full-chip grid has far too many nodes for dense elimination, and even sparse direct factorisation runs out of memory on the largest designs. IR-drop engines therefore use iterative methods, or sparse direct solvers with careful ordering, exploiting the fact that $G$ is sparse, symmetric and positive definite. Learning the simple iterations by hand shows what the production solvers accelerate.

## Jacobi sweeps

Rewrite row $i$ of $G d = i$ as $d_i = (i_i - \sum_{j \ne i} G_{ij} d_j)/G_{ii}$. Jacobi updates every node from the previous sweep's values. On a three-node rail with supplies at both ends, resistances 0.5, 2, 2 and $1.5\ \Omega$ and loads 14, 10 and 5 mA, starting from zero: sweep 1 gives $(5.6, 10, 4.29)$ mV, sweep 2 gives $(7.6, 14.94, 8.57)$. Each sweep is one sparse matrix-vector product, perfectly parallel. Jacobi converges for this kind of matrix but slowly, because a correction travels only one node per sweep.

## Gauss-Seidel sweeps

Gauss-Seidel uses each new value as soon as it is computed. On the rail with resistances 2, 2, 1.5 and $1\ \Omega$ and loads 15, 5 and 7 mA, one sweep gives $d_1 = 15/1 = 15$, then $d_2 = (5 + 0.5 \cdot 15)/1.1667 = 10.71$, then $d_3 = (7 + 0.6667 \cdot 10.71)/1.6667 = 8.49$ mV. It usually converges about twice as fast as Jacobi on such matrices, at the cost of sequential dependence. Both are smoothers: they remove local error quickly and global error slowly.

## Scale and the solvers that handle it

An $n \times n$ grid layer has $n^2$ nodes and $2n(n - 1)$ branches, so $G$ has $n^2 + 4n(n - 1)$ non-zeros, about five per row: 2.66 million for $n = 730$. Production tools use preconditioned conjugate gradient, which works because $G$ is symmetric positive definite, with incomplete Cholesky or multigrid preconditioners that attack the slow global error. Sparse direct solvers with fill-reducing ordering are used where memory allows, and they make repeated solves with new currents cheap. Random-walk methods estimate single nodes without solving the whole grid, useful for incremental checks.

## Going deeper

Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., 2003); Kahng et al., "VLSI Physical Design" (2nd ed., 2022).
