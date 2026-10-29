# Small pivots and element growth

Gaussian elimination is only as stable as its pivots. A tiny pivot, which in a circuit can come from a nearly floating node, a minimum gmin conductance or a transistor switched almost off, produces huge multipliers that swamp the rest of the matrix in rounding error. Recognizing and preventing that growth is the reason every robust solver pivots.

## A tiny pivot creates a huge multiplier

Take
$$A = \begin{pmatrix} 10^{-6} & 1 \\ 2 & 3 \end{pmatrix}.$$
Without row exchange the multiplier is $l_{21} = 2/10^{-6} = 2 \times 10^6$. The update gives $u_{22} = 3 - 2 \times 10^6 \cdot 1 = -1999997$.

In floating point the information in the original entry 3 is now carried in the low-order digits of a number near $-2 \times 10^6$. With a pivot of $10^{-17}$ instead, $3 - 2 \times 10^{17}$ rounds to exactly $-2 \times 10^{17}$ in double precision, and the entry 3 is lost completely. The computed factors then describe a different matrix, and the solution can be wrong in every digit even though the original matrix is perfectly well conditioned.

## Growth factor of an elimination

The growth factor measures how much elimination inflates entries:
$$\rho = \frac{\max_{i,j,k} |a^{(k)}_{ij}|}{\max_{i,j} |a_{ij}|},$$
and in practice is computed as the largest entry of $U$ over the largest entry of $A$. Backward error analysis shows that LU computes the exact factors of $A + E$ with $\|E\|$ bounded by a modest multiple of $\rho\, u\, \|A\|$, with $u \approx 1.1 \times 10^{-16}$ the unit roundoff. A large $\rho$ means a large backward error.

For the matrix above, $\max|U| = 1999997$ and $\max|A| = 3$, so $\rho \approx 6.7 \times 10^5$: about six decimal digits of backward accuracy are lost before conditioning even enters.

## Swapping rows restores stability

Partial pivoting exchanges rows so the pivot is the largest entry in its column. For the example, swap the rows:
$$\begin{pmatrix} 2 & 3 \\ 10^{-6} & 1 \end{pmatrix}, \quad l_{21} = 5 \times 10^{-7}, \quad u_{22} = 1 - 5 \times 10^{-7} \cdot 3 = 0.9999985.$$
Now every multiplier is at most 1 in magnitude, nothing grows, and the factors are accurate to roundoff.

Partial pivoting bounds the growth factor by $2^{n-1}$ in theory; in practice it stays small. For sparse matrices, though, the largest pivot may sit in a dense row and cause heavy fill, so circuit solvers relax the rule to threshold pivoting. Diagonally dominant nodal matrices need no pivoting at all, which is one reason gmin and careful MNA stamping help.

## Going deeper

Higham, "Accuracy and Stability of Numerical Algorithms" (2nd ed., SIAM, 2002), chapter 9; Trefethen and Bau, "Numerical Linear Algebra" (SIAM, 1997), lectures 20 to 22.
