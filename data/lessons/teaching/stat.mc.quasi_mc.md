# Quasi-Monte Carlo sequences

Plain Monte Carlo spends its samples at random, so they clump and leave gaps. Quasi-Monte Carlo replaces random points with deterministic low-discrepancy sequences that fill the space evenly. For smooth, effectively low-dimensional integrands such as the mean of a delay or a moderately rare yield, the error can fall close to $1/N$ instead of $1/\sqrt{N}$, a large saving in circuit simulations.

## The van der Corput sequence

The building block is the radical inverse. Write the integer $n$ in base $b$, $n = \sum_k d_k b^k$, and mirror its digits about the radix point: $\phi_b(n) = \sum_k d_k b^{-k-1}$. In base 2, $n = 1, 2, 3, 4, 5, 6$ are $1, 10, 11, 100, 101, 110$ and map to $0.1, 0.01, 0.11, 0.001, 0.101, 0.011$ in binary, that is $0.5, 0.25, 0.75, 0.125, 0.625, 0.375$. Each new point lands in the largest remaining gap. After $2^m$ points, every interval $[j/2^m, (j+1)/2^m)$ holds exactly one point, which random sampling achieves only by luck. This is what low discrepancy means: the fraction of points in any subinterval tracks the length of that subinterval far more closely than it does for random points.

## Halton points in several dimensions

A Halton sequence uses a different prime base per coordinate: point $n$ is $(\phi_2(n), \phi_3(n), \phi_5(n), \ldots)$. For $n = 7$: in base 2, $111 \to 0.111_2 = 0.875$; in base 3, $7 = 21_3 \to 0.12_3 = 1/3 + 2/9 = 0.556$; in base 5, $7 = 12_5 \to 0.21_5 = 2/5 + 1/25 = 0.44$. The points are mapped to Gaussian process parameters through the inverse normal CDF. In high dimensions, plain Halton coordinates with large bases are strongly correlated for small $n$, so practical engines use Sobol sequences with good direction numbers, or scramble the digits (randomized QMC), which also restores an unbiased error estimate from independent replicates.

## Convergence rate and its limits

The Koksma-Hlawka inequality bounds the QMC error by the variation of the integrand times the discrepancy of the points, which for low-discrepancy sequences is $O((\log N)^d/N)$. Ignoring the log factor, cutting the error by a factor $r$ needs $r$ times the points, against $r^2$ for plain Monte Carlo: a tenfold improvement costs 10x instead of 100x. The gain shrinks when the integrand is discontinuous, as a pass/fail indicator is, and when many dimensions matter equally. QMC is therefore strongest for moments of smooth outputs and for problems whose variance is concentrated in a few principal components, and weakest for very rare failure indicators, where importance sampling is the better tool.

## Going deeper

Niederreiter, "Random Number Generation and Quasi-Monte Carlo Methods" (1992); Owen, "Monte Carlo Theory, Methods and Examples," the chapters on quasi-Monte Carlo and randomized QMC.
