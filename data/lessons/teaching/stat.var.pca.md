# Principal components of correlated variation

A statistical timer that tracks a separate random variable for every grid cell of a die, all correlated with each other, would carry thousands of dependent variables through every addition and max. Principal component analysis replaces them with a smaller set of independent variables. This is the step that turns a spatial correlation model into the canonical delay forms used by block-based statistical timing.

## Diagonalizing the covariance

Partition the die into $n$ grid cells and let $\Sigma$ be the $n \times n$ covariance matrix of a parameter across them. It is symmetric and positive semidefinite, so the spectral theorem gives $\Sigma = V \Lambda V^T$ with orthonormal eigenvectors $v_k$ and eigenvalues $\lambda_k \ge 0$. Define independent standard normal variables $z_k$; then $x = V \Lambda^{1/2} z$ has exactly covariance $\Sigma$. Each cell's variation becomes a weighted sum of the same few $z_k$. For a two-cell matrix with entries $a$, $b$ on the diagonal and $c$ off it, the eigenvalues are $\frac{a+b}{2} \pm \sqrt{(\frac{a-b}{2})^2 + c^2}$; with $a = 10$, $b = 6$, $c = 4$ they are $8 \pm 4.47$, so 12.47 and 3.53.

## Variance explained and truncation

The total variance is the trace, $\sum_k \lambda_k$. The share explained by the first $m$ components is $\sum_{k \le m}\lambda_k / \mathrm{tr}\,\Sigma$. In the example, the first component carries $12.47/16 = 78$ percent. Spatial correlation makes the eigenvalues fall quickly, so a timer can keep the components that explain, say, 95 percent of the variance and drop the rest, or fold the remainder into each cell's independent random term so total variance is preserved. That truncation is the main speed lever of statistical timing: cost grows with the number of retained components, not with the number of grid cells.

## The equicorrelated case

A clean limit shows what the components mean. If every one of $n$ cells has unit variance and every pair has correlation $\rho$, the matrix is $(1-\rho)I + \rho\,\mathbf{1}\mathbf{1}^T$. Its eigenvalues are $1 + (n-1)\rho$ once, for the all-ones direction, and $1 - \rho$ repeated $n-1$ times. The first component is a die-wide shift, which is exactly global variation; it explains $(1 + (n-1)\rho)/n$ of the variance, which tends to $\rho$ as $n$ grows. With $n = 16$ and $\rho = 0.5$ that is $8.5/16 = 53$ percent. Real spatial matrices behave in between: a dominant global component, then smooth gradient-like components, then a long flat tail that is close to independent noise.

## Going deeper

Chang and Sapatnekar, "Statistical timing analysis under spatial correlations," IEEE TCAD, 2005; Jolliffe, "Principal Component Analysis" (2nd ed., 2002).
