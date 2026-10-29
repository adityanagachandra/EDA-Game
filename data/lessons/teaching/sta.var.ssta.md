# Statistical static timing analysis

POCV treats each path's random variation statistically but keeps corners for correlated variation. Block-based statistical timing goes further: arrival times themselves become random variables that share global variation sources, propagated through the graph with statistical sum and max. The standard representation is the first-order canonical form, which lets sum and max stay in closed form and tracks correlation between paths through shared sources.

## The canonical form and its sigma

A canonical first-order form writes a delay or arrival as
$$A = a_0 + \sum_i a_i X_i + a_r R,$$
where $a_0$ is the mean, each $X_i$ is a shared global variation source such as a channel-length or threshold-voltage component, $R$ is a private random term, and all are independent standard normals. The coefficients are sensitivities in ps per sigma. Its variance is $\sum_i a_i^2 + a_r^2$. For $A = 150 + 9X_1 + 9X_2 + 8R$, $\sigma_A = \sqrt{81 + 81 + 64} = 15.03$ ps. The sources often come from a principal-component decomposition of spatially correlated variation, which makes them independent by construction.

## Adding canonical forms

Propagating through an arc adds its delay to the arrival. Sums of canonical forms are canonical: means add, sensitivities to each shared source add, and the private random terms, being independent, combine by root-sum-square. For $A = 300 + 9X_1 + 9X_2 + 3R$ and an arc $D = 430 - 3X_1 + 5X_2 + 10R'$, the sum is $730 + 6X_1 + 14X_2 + 10.44R''$, with $\sigma = \sqrt{36 + 196 + 109} = 18.47$ ps. Note that $\sigma_{A+D}$ is not $\sigma_A + \sigma_D$, nor their root-sum-square: the shared $X_1$ terms partly cancel and the $X_2$ terms add coherently. That is exactly the correlation a corner-based analysis cannot express.

## Correlation through shared sources

Two arrivals in canonical form are correlated only through their shared sources: $\text{cov}(A, B) = \sum_i a_i b_i$, and $\rho = \text{cov}/(\sigma_A \sigma_B)$. For $A = 340 + 11X_1 + 10X_2 + 4R$ and $B = 240 + 6X_1 - 3X_2 + 8R'$, the covariance is $66 - 30 = 36$ and $\rho = 36 / (15.40 \times 10.44) = 0.224$. This correlation matters at every merge: the max of two highly correlated arrivals is close to the larger one, while the max of two independent arrivals is pushed noticeably later. It also matters at the end, where slacks of many endpoints share sources, so the chip's timing yield is not the product of per-path yields.

## Going deeper

Visweswariah et al., "First-order incremental block-based statistical timing analysis" (DAC 2004); Chang and Sapatnekar, "Statistical timing analysis considering spatial correlations using a single PERT-like traversal" (ICCAD 2003); Sapatnekar, "Timing" (2004).
