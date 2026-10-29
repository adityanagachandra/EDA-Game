# LU factorization for circuit matrices

A transient simulation of a large circuit solves a sparse linear system at every Newton iteration of every timestep, often millions of times. Sparse LU factorization is that solve. Circuit matrices have a fixed pattern, so the work splits into a symbolic phase done once and a numeric refactorization repeated many times; designing for that split is what makes solvers such as KLU fast.

## Eliminating a nodal matrix row by row

LU factorization writes $A = LU$ with $L$ unit lower triangular and $U$ upper triangular. Step $k$ eliminates column $k$ below the pivot: for each row $i > k$, the multiplier is $l_{ik} = a_{ik}/a_{kk}$, and row $i$ becomes row $i$ minus $l_{ik}$ times row $k$.

Take
$$A = \begin{pmatrix} 4 & -1 & -2 \\ -1 & 5 & -2 \\ -2 & -2 & 6 \end{pmatrix}.$$
Step 1 has multipliers $l_{21} = -1/4 = -0.25$ and $l_{31} = -2/4 = -0.5$. These are stored in $L$; the updated rows become rows of $U$ and the next active submatrix. A nodal matrix that is diagonally dominant, as most resistive networks are, has every multiplier at most 1 in magnitude without any row exchange.

## The Schur complement update

After step 1 the trailing block is the Schur complement $S = A_{22} - a_{21}a_{12}^T/a_{11}$. Continuing the example:
$u_{22} = 5 - (-0.25)(-1) = 4.75$, $u_{23} = -2 - (-0.25)(-2) = -2.5$, updated $a_{32} = -2 - (-0.5)(-1) = -2.5$, updated $a_{33} = 6 - (-0.5)(-2) = 5$.
Step 2: $l_{32} = -2.5/4.75 = -0.5263$ and
$$u_{33} = 5 - (-0.5263)(-2.5) = 3.6842.$$
The determinant is the product of pivots, $4 \cdot 4.75 \cdot 3.6842 = 70$. Once $L$ and $U$ exist, each new right-hand side costs only a forward and a back substitution, $O(\text{nnz}(L) + \text{nnz}(U))$ work. That is why one factorization per Newton iteration serves every solve inside it.

## Symbolic structure before numeric values

The pattern of an MNA matrix depends only on the netlist, not on device operating points, so it is the same at every Newton iteration and timestep. A sparse LU therefore runs a symbolic phase once: choose an ordering and compute the patterns of $L$ and $U$. Every later numeric factorization reuses that analysis and allocates nothing.

With a symmetric pattern and no pivoting, the factors hold $\text{nnz}(A) + 2 \cdot \text{fill}$ entries, counting the diagonal once in $U$. A 5-vertex pattern with 4 off-diagonal pairs (nnz 13) and 2 fill edges stores 17 entries.

KLU, designed for circuit matrices, adds a block triangular form step that splits the matrix into independent irreducible blocks, orders each with AMD or COLAMD, and refactors with a left-looking algorithm that keeps the pivot sequence when it stays stable.

## Going deeper

Davis and Palamadai Natarajan, "Algorithm 907: KLU, a direct sparse solver for circuit simulation problems", ACM TOMS (2010); Davis, "Direct Methods for Sparse Linear Systems" (SIAM, 2006), chapters 4 to 6.
