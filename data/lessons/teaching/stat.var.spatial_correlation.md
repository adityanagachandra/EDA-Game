# Spatially correlated variation

Within-die variation is not purely random from device to device. Lithography focus, etch loading and wafer-level gradients produce slow spatial trends, so neighbouring gates are more alike than distant ones. A statistical timer needs those correlations to get path delays and clock skews right: correlation along a path increases its spread, while correlation between two clock sinks shrinks their skew.

## A distance-based correlation model

A common model writes the spatial component $S(x)$ as a random field whose correlation depends only on distance $d$ between two points. The exponential model is $\rho(d) = e^{-d/\lambda}$, where the correlation length $\lambda$ is the distance at which correlation falls to $1/e \approx 0.37$. With $\lambda = 500\ \mu$m, two gates $200\ \mu$m apart have $\rho = e^{-0.4} = 0.67$, and gates 2 mm apart have $\rho = e^{-4} = 0.018$. Other kernels (Gaussian, spherical, linear with a cutoff) are also used; whatever the choice, the resulting covariance matrix must be positive semidefinite, which arbitrary hand-written correlation tables often are not. The exponential kernel satisfies that requirement in any dimension.

## Correlation of two gate delays

Combine the field with the other components: $D = \mu + G + S + R$ with variances $\sigma_G^2$, $\sigma_S^2$ and $\sigma_R^2$. For two gates at distance $d$, the global term is fully shared, the spatial term is shared with weight $\rho(d)$, and the random term not at all:
$$\mathrm{Cov}(D_1, D_2) = \sigma_G^2 + \sigma_S^2\,\rho(d), \qquad \mathrm{Corr} = \frac{\sigma_G^2 + \sigma_S^2\rho(d)}{\sigma_G^2 + \sigma_S^2 + \sigma_R^2}.$$
With $\sigma_G = 3$, $\sigma_S = 4$, $\sigma_R = 2$ ps and $\rho = 0.5$, the covariance is $9 + 8 = 17$, the variance is 29, and the correlation is 0.59. Even far apart, the correlation never drops below the global floor $9/29 = 0.31$.

## Skew between two sinks

Clock skew is a difference, so shared parts cancel. $\mathrm{Var}(D_1 - D_2) = 2\sigma_S^2(1 - \rho(d)) + 2\sigma_R^2$, and the global term drops out entirely. Close sinks have $\rho$ near 1 and see only random skew; distant sinks see the full spatial variance. This is the statistical reason clock-tree synthesis keeps the sinks of timing-critical pairs close in the tree and on the die, and why on-chip variation derates that grow with distance (as in AOCV tables) exist. An engine that ignores spatial correlation must pick one pessimistic number for every pair, which is too loose for close pairs and possibly optimistic for far ones.

## Going deeper

Xiong, Zolotov and He, "Robust extraction of spatial correlation," IEEE TCAD, 2007; Agarwal, Blaauw and Zolotov, "Statistical timing analysis for intra-die process variations with spatial correlations," ICCAD 2003.
