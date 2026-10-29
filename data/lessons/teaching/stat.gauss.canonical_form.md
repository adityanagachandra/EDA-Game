# Canonical first-order delay forms

A block-based statistical timer must store every arrival time in a representation that is closed under addition, approximately closed under max, and that remembers correlation. The first-order canonical form does this: a mean, sensitivities to a shared set of independent variation sources from principal component analysis, and one independent random term. It is the data structure inside the statistical timers built on this idea.

## The representation

An arrival time is written
$$A = a_0 + \sum_{i=1}^{m} a_i\,\Delta X_i + a_{r}\,\Delta R_A,$$
where the $\Delta X_i$ are independent standard normal principal components shared by the whole design, and $\Delta R_A$ is a standard normal private to this quantity. The variance is $\sigma_A^2 = \sum_i a_i^2 + a_r^2$. With $A = 100 + 3\Delta X_1 - 2\Delta X_2 + 2\Delta R_A$, the sigma is $\sqrt{9 + 4 + 4} = 4.12$ ps. Sums are exact: add the means and the sensitivities to each shared component, and combine the private terms in variance, $\sqrt{a_r^2 + b_r^2}$, because they are independent.

## Correlation comes for free

Two canonical forms covary only through their shared components. The private terms are independent of each other, so
$$\mathrm{Cov}(A, B) = \sum_i a_i b_i, \qquad \rho_{AB} = \frac{\sum_i a_i b_i}{\sigma_A\sigma_B}.$$
For $A$ above and $B = 95 + 1\Delta X_1 + 4\Delta X_2 + 3\Delta R_B$, $\mathrm{Cov} = 3 - 8 = -5$ and $\sigma_B = \sqrt{26}$, so $\rho = -5/(4.12 \times 5.10) = -0.24$. No covariance matrix between arrival times is ever stored: correlation is recovered on demand from a dot product. This is why the representation scales: memory per node grows with the number of retained components, not with the number of nodes it might be correlated with.

## Re-expressing the max

Clark's formulas give the mean and variance of $\max(A, B)$, but the timer needs the result in canonical form again. The standard construction weights the sensitivities by tightness probability:
$$c_i = T_A a_i + (1 - T_A) b_i, \qquad c_0 = E[\max(A,B)],$$
and sets the private coefficient so the total variance matches Clark's variance. With $T_A = 0.7$, $a_1 = 3$ and $b_1 = 1$, the max has $c_1 = 2.4$. The result is a linear approximation of a nonlinear operation: exact for the first two moments and approximately right for the correlation with other quantities. Every max in the graph adds approximation error, which is why statistical timers validate against Monte Carlo on representative designs.

## Going deeper

Visweswariah, Ravindran, Kalafala, Walker and Narayan, "First-order incremental block-based statistical timing analysis," DAC 2004; Chang and Sapatnekar, "Statistical timing analysis under spatial correlations," IEEE TCAD, 2005.
