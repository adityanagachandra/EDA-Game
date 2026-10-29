# Cut size as a Laplacian quadratic form

The quadratic form $x^T L x$ is one formula with two EDA meanings. Evaluated on a $\pm 1$ vector it is four times the cut of a bisection; evaluated on real coordinates it is quadratic wirelength. That identity is why spectral partitioning and quadratic placement are the same mathematics, and why relaxing a discrete cut problem to real vectors gives a solvable eigenproblem.

## The sum-of-squares identity

For any real vector $x$,

$$x^T L x = \sum_{(i,j) \in E} w_{ij}(x_i - x_j)^2.$$

Expanding the right side gives $\sum_i x_i^2 D_{ii} - 2\sum_{(i,j)} w_{ij} x_i x_j$, which is $x^T D x - x^T A x$. The edge form is how tools evaluate it: one pass over the edges, without building $L$.

For edges $v_1$-$v_2$ (3), $v_1$-$v_3$ (3), $v_2$-$v_4$ (4), $v_3$-$v_4$ (3) and $x = (4, 2, 2, 1)$: $3 \cdot 4 + 3 \cdot 4 + 4 \cdot 1 + 3 \cdot 1 = 31$. Each term is non-negative, which is the direct proof that $L$ is positive semidefinite.

## Bisections as plus-or-minus-one vectors

Encode a bisection by $x_i = +1$ on side $A$ and $x_i = -1$ on side $B$. An uncut edge has $x_i - x_j = 0$; a cut edge has $(x_i - x_j)^2 = 4$. So

$$x^T L x = 4 \cdot \text{cut}(A, B).$$

A bisection with cut weight 24 gives $x^T L x = 96$. With a 0/1 indicator vector instead, $x^T L x$ equals the cut directly. Minimum bisection is therefore the minimization of $x^T L x$ over $\pm 1$ vectors with $\sum_i x_i = 0$. That discrete problem is NP-hard; relaxing $x$ to real values with $x \perp \mathbf{1}$ and $\|x\|^2 = n$ makes it an eigenvalue problem solved by the Fiedler vector.

## The same form in placement

Put cells at positions $x_i$ on a line and expand nets as weighted cliques. Then $x^T L x$ is the quadratic wirelength $\sum_e \sum_{i<j \in e} w_e (x_i - x_j)^2$. Cells at $x = 8, 2, 1, 6$ on a 4-pin net with weight $1/3$ contribute $\big(36 + 49 + 4 + 1 + 16 + 25\big)/3 = 131/3$.

Quadratic placers minimize this with I/O pins fixed, by solving $L_{mm} x_m = b$ where $L_{mm}$ is the movable block and $b$ collects the fixed-pin terms. The squared distance overweights long nets compared with linear wirelength, which is why placers reweight or switch to log-sum-exp or weighted-average smooth wirelength models. But the Laplacian structure, sparse, symmetric and positive definite, is what makes the quadratic version fast enough to start from.

## Going deeper

Hall, "An r-dimensional quadratic placement algorithm" (Management Science, 1970); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 4.
