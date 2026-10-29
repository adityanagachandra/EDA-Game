# Incomplete LU and incomplete Cholesky

Incomplete LU and incomplete Cholesky factorizations run Gaussian elimination but throw away most of the fill. What remains is a sparse approximate factorization, cheap to store and apply, that captures the matrix's local coupling. They are the default general-purpose preconditioners for power-grid CG and circuit GMRES solvers.

## Factor but keep the original pattern

ILU(0) computes $L$ and $U$ by the usual elimination updates $a_{ij} \leftarrow a_{ij} - l_{ik}u_{kj}$, but performs an update only if $(i, j)$ is a nonzero position of $A$. Every fill entry is dropped. For SPD matrices the symmetric version is IC(0), incomplete Cholesky.

The factors have exactly the pattern of $A$, so storage and application cost match one SpMV each. The preconditioner is $M = \tilde L\tilde U$, which equals $A$ on the original pattern but differs from it by the dropped fill. Example: in a 3-by-3 matrix with $a_{23} = a_{32} = 0$ but $a_{21}, a_{13}, a_{31}, a_{12}$ nonzero, eliminating row 1 would create fill at $(2, 3)$ and $(3, 2)$; ILU(0) discards both.

## The entry ILU(0) drops

Take
$$A = \begin{pmatrix} 4 & -1 & -1 \\ -1 & 4 & 0 \\ -1 & 0 & 4 \end{pmatrix}.$$
Exact elimination of column 1: $l_{21} = l_{31} = -0.25$; the $(2, 3)$ entry becomes $0 - (-0.25)(-1) = -0.25$, fill; $u_{22} = 4 - 0.25 = 3.75$; updated $a_{33} = 3.75$; $l_{32} = -0.25/3.75 = -0.0667$; exact $u_{33} = 3.75 - (-0.0667)(-0.25) = 3.7333$.

ILU(0) drops the fill at $(2, 3)$ and $(3, 2)$, so $l_{32} = 0$ and $u_{33} = 3.75$. The difference, 0.0167, is the error the preconditioner makes; the Krylov method absorbs it. For diagonally dominant M-matrices like power grids, the dropped entries are small and IC(0) always exists.

## Levels of fill and drop tolerances

Two knobs trade accuracy for cost. ILU($k$) assigns each entry a level: original entries level 0, and an entry created from entries of levels $a$ and $b$ gets level $a + b + 1$; entries above level $k$ are dropped. ILU(1) keeps first-generation fill, usually a large improvement on grids.

ILUT drops by size instead: a computed entry in row $i$ is discarded if $|x| < \tau\|a_{i*}\|_2$, and often only the largest $p$ entries per row are kept. With row norm 10 and $\tau = 0.01$, fill entries of magnitude 0.05 are dropped and 0.2 kept. Stronger preconditioners cut iterations but cost memory and setup, and incomplete factorizations of indefinite or non-M-matrices can break down with zero or negative pivots, which is why robust codes add diagonal shifts.

## Going deeper

Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., SIAM, 2003), chapter 10; Meijerink and van der Vorst, "An iterative solution method for linear systems of which the coefficient matrix is a symmetric M-matrix", Math. Comp. (1977).
