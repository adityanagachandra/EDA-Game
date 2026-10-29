# Conditioning of circuit matrices

Even a perfectly stable solver cannot recover digits that the problem itself does not determine. The condition number of the system matrix measures that sensitivity. Circuit matrices mix conductances spanning twenty orders of magnitude, from switch resistances to leakage and gmin, so conditioning sets how many digits of a node voltage a simulator can actually trust.

## Condition number of a conductance matrix

For a nonsingular matrix, $\kappa(G) = \|G\|\,\|G^{-1}\|$. In the infinity norm, $\|G\|_\infty$ is the largest absolute row sum. Take two nodes with 10 mS from node 1 to ground, 5 mS between them, and 1 mS from node 2 to ground:
$$G = \begin{pmatrix} 15 & -5 \\ -5 & 6 \end{pmatrix}\ \text{mS}, \quad \det G = 65.$$
$\|G\|_\infty = 20$. The inverse is $\frac{1}{65}\begin{pmatrix} 6 & 5 \\ 5 & 15 \end{pmatrix}$, with $\|G^{-1}\|_\infty = 20/65 = 0.3077$. So $\kappa_\infty = 6.15$, a very well conditioned system.

Now let both ground conductances be 0.001 mS, so the pair is held to ground only by leakage: $G = \begin{pmatrix} 5.001 & -5 \\ -5 & 5.001 \end{pmatrix}$, $\det G \approx 0.01$ and $\kappa_\infty \approx 10^4$. Nodes that reach ground only through tiny conductances are nearly floating, and $G$ is nearly singular.

## From residual to error bound

A computed solution $\hat x$ has residual $r = b - G\hat x$, which the solver can measure. The error it cannot measure is bounded by
$$\frac{\|x - \hat x\|}{\|x\|} \le \kappa(G)\,\frac{\|r\|}{\|b\|}.$$
With $\kappa = 2 \times 10^6$ and relative residual $10^{-12}$, the relative error is at most $2 \times 10^{-6}$, that is 2 parts per million.

This is why a small residual is not proof of an accurate answer. An iterative solver stopped at relative residual $10^{-8}$ on a matrix with $\kappa = 10^7$ may have no correct digits. The bound is a worst case; actual errors are often smaller, but tool designers set tolerances against it, and condition estimators (as in LAPACK's xGECON, which costs only a few solves) make it computable.

## Conductance spread, gmin and digits lost

A rule of thumb: solving with a matrix of condition number $\kappa$ loses about $\log_{10}\kappa$ decimal digits from the roughly 16 that double precision carries ($u \approx 1.1 \times 10^{-16}$). A circuit with a 1 kS closed switch and a $10^{-12}$ S gmin leak has a spread near $10^{15}$ before scaling, leaving about one digit.

Remedies are structural. Row and column scaling (equilibration) removes artificial spread caused by units. Gmin, typically $10^{-12}$ S, is added from nodes to ground to keep $G$ nonsingular when nodes float, but it is a small perturbation of the circuit and should not dominate any real conductance. Ideal switches modelled as tiny and huge resistances are better replaced by formulations that change the topology. When digits matter (precision analog), tools report the condition estimate so an engineer knows when to distrust a result.

## Going deeper

Higham, "Accuracy and Stability of Numerical Algorithms" (2nd ed., SIAM, 2002), chapter 7; Golub and Van Loan, "Matrix Computations" (4th ed., 2013), section 2.6.
