# Spectral bisection with the Fiedler vector

Spectral bisection relaxes the discrete cut problem to an eigenproblem: the eigenvector of the second-smallest Laplacian eigenvalue, the Fiedler vector, embeds the vertices on a line so that strongly connected vertices sit close together. Splitting that line gives a partition. It is a global method, seeing the whole graph at once, which complements local FM, and its eigenvalue is a certificate of how well connected the graph is.

## Algebraic connectivity

Order the Laplacian eigenvalues $0 = \lambda_1 \le \lambda_2 \le \dots \le \lambda_n$. The second, $\lambda_2$, is the **algebraic connectivity**: it is positive exactly when the graph is connected, and small when the graph has a sparse cut. By the Courant-Fischer theorem, $\lambda_2 = \min_{x \perp \mathbf{1},\, x \ne 0} x^T L x / x^T x$.

For a three-vertex path with weights $w_1$ and $w_2$, the characteristic polynomial is $\lambda(\lambda^2 - 2(w_1 + w_2)\lambda + 3w_1 w_2)$, so $\lambda_2 = (w_1 + w_2) - \sqrt{w_1^2 - w_1 w_2 + w_2^2}$. With $w_1 = 4$ and $w_2 = 2$, $\lambda_2 = 6 - \sqrt{12} \approx 2.54$. A weak link anywhere pulls $\lambda_2$ toward 0.

## Splitting the Fiedler vector

The eigenvector $f$ of $\lambda_2$ is the Fiedler vector. Sort the vertices by $f_i$ and split: by sign, or at the median for an exact bisection, or at the position along the sorted order with the best ratio cut. Vertices in a tightly connected cluster get similar entries, because a large difference across a heavy edge would cost $w_{ij}(f_i - f_j)^2$ in the Rayleigh quotient.

For a graph of two dense groups joined by a few light edges, the Fiedler vector is positive on one group and negative on the other, and the median split cuts only the light edges. Computing it uses a sparse eigensolver, such as Lanczos iteration, in time roughly linear in the edges per iteration. Spectral bisection was a leading method in the early 1990s, before multilevel FM proved better on netlists.

## Rayleigh quotients as bounds

Any vector $x \perp \mathbf{1}$ gives an upper bound $\lambda_2 \le x^T L x / x^T x$. For $x = (-2, 3, 1, -1, -1)$ on a small graph with $x^T L x = 67$ and $x^T x = 16$, the quotient is about 4.19, so $\lambda_2 \le 4.19$. Iterative eigensolvers improve such a vector until the quotient converges to $\lambda_2$.

In the other direction, $\lambda_2$ bounds every cut from below: $\lambda_2 \le \text{RatioCut}(A, B)$ for every bisection, and Cheeger-type inequalities relate it to the best conductance. So a spectral computation can certify that no very sparse cut exists. Its weaknesses for netlists are the clique distortion of hyperedges and the loss of information when only one eigenvector is used; using several eigenvectors as coordinates gives a multi-dimensional embedding, an idea that returns in analytical placement.

## Going deeper

Fiedler, "Algebraic connectivity of graphs" (Czechoslovak Mathematical Journal, 1973); Pothen, Simon and Liou, "Partitioning sparse matrices with eigenvectors of graphs" (SIAM J. Matrix Analysis and Applications, 1990); Hagen and Kahng (IEEE TCAD, 1992).
