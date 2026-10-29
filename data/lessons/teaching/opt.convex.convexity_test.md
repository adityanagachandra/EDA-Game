# Testing convexity with the Hessian

Convexity is the property that makes continuous optimization trustworthy: every local minimum of a convex function over a convex set is global, so descent methods cannot get trapped. Analytical placers, continuous sizers and LP and QP subproblems are designed to be convex or nearly so. The fastest test for a smooth function is its Hessian, and for EDA cost terms the composition rules, max of convex is convex, log-sum-exp is convex, posynomials are convex in log space, settle most cases at a glance.

## Classifying a quadratic by its Hessian

A twice-differentiable function is convex exactly when its Hessian is positive semidefinite everywhere. For a quadratic $f(x, y) = ax^2 + bxy + cy^2 + dx$ the Hessian is constant, $H = \begin{pmatrix} 2a & b \\ b & 2c \end{pmatrix}$, so one check covers the whole plane. Linear terms never affect convexity.

For a 2 by 2 symmetric matrix: if $\det H > 0$ and $2a > 0$, H is positive definite and f is strictly convex; if $\det H > 0$ and $2a < 0$, negative definite and f strictly concave; if $\det H < 0$, the eigenvalues have opposite signs and f is a saddle, neither convex nor concave; if $\det H = 0$, H is semidefinite, and f is convex but not strictly when the diagonal is nonnegative. Example: $f = -x^2 - 4xy - 3y^2 - x$ has $H = \begin{pmatrix} -2 & -4 \\ -4 & -6 \end{pmatrix}$ with determinant $12 - 16 = -4$: a saddle.

## The determinant test in practice

The determinant of $H = \begin{pmatrix} 2a & b \\ b & 2c \end{pmatrix}$ is $4ac - b^2$. The common slips are forgetting the factor 2 on the diagonal, which gives $ac - b^2$, or doubling the cross term, which gives $4ac - 4b^2$. Example: $f = 2x^2 - 5xy + 3y^2 - 5x$ has $H = \begin{pmatrix} 4 & -5 \\ -5 & 6 \end{pmatrix}$ and determinant $24 - 25 = -1$, so although both squared terms have positive coefficients the cross term makes f a saddle.

This matters for placement cost models built from many terms. A sum of convex terms is convex, but a single term with a large coupling coefficient can break it. In higher dimensions the test generalizes: a symmetric matrix is positive definite exactly when all leading principal minors are positive (Sylvester's criterion), and positive semidefinite when all principal minors are nonnegative; numerically, a Cholesky factorization succeeds exactly for positive definite matrices.

## Convexity of EDA cost terms

Most cost terms are checked by rules, not Hessians. HPWL, $\max_i x_i - \min_i x_i$, is convex: a max of linear functions is convex, and $-\min_i x_i = \max_i(-x_i)$ is too. It is not smooth, which is why placers use log-sum-exp, which is convex and smooth. Quadratic wirelength $\sum w_{ij}(x_i - x_j)^2$ is convex, but its Hessian is the graph Laplacian, singular because shifting all cells changes nothing; fixed pads make it definite.

A posynomial delay, such as $\sum_k c_k \prod_j x_j^{a_{jk}}$ with $c_k > 0$, is not convex in x in general, but substituting $x_j = e^{u_j}$ and taking the log gives log-sum-exp of affine functions, which is convex: the basis of geometric-programming sizing. Density overflow with sharp bin edges is nonconvex and nonsmooth. Pointwise max and nonnegative sums preserve convexity; products, differences and minima of convex functions in general do not.

## Going deeper

Boyd and Vandenberghe, "Convex Optimization" (2004), chapter 3 on convex functions and operations that preserve convexity; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 4, for the wirelength and density models used by analytical placers.
