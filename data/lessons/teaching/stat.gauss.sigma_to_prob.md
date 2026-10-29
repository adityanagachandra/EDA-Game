# Sigma level and failure probability

Designers state margins in sigmas: a 3-sigma corner, a 6-sigma bitcell. A tool designer must be able to convert that language into probabilities and back, instantly and in both directions, because every rare-event target, array-yield budget and sample-size estimate in variation analysis is set in sigmas and checked in probabilities. The conversion is the Gaussian tail function, applied with care about which tail and how many independent trials.

## The one-sided Gaussian tail

If a margin $M$ is Gaussian and a cell fails when $M$ falls more than $k$ standard deviations below its mean, the failure probability is the upper tail $Q(k) = 1 - \Phi(k) = \tfrac{1}{2}\mathrm{erfc}(k/\sqrt{2})$. Useful anchors: $Q(3) = 1.35 \times 10^{-3}$ (1350 ppm), $Q(4) = 3.17 \times 10^{-5}$ (31.7 ppm), $Q(5) = 2.87 \times 10^{-7}$, $Q(6) = 9.87 \times 10^{-10}$. Each additional sigma costs one to three orders of magnitude, faster as $k$ grows. A two-sided spec doubles the tail. Software should compute $Q$ with erfc rather than as $1 - \Phi(k)$, which loses every significant digit to cancellation beyond about 8 sigma in double precision. The quality-engineering "six sigma" convention of 3.4 ppm assumes a 1.5-sigma mean shift; it is not $Q(6)$.

## Expected failures in an array

A memory multiplies a tiny per-cell probability by a huge count. An array of $N$ independent cells, each failing with probability $p$, contains on average $Np$ failing cells. A 16 Mb macro holds $16 \times 2^{20} = 16{,}777{,}216$ cells; at 4.5 sigma, $p = 3.40 \times 10^{-6}$ and the macro expects $16{,}777{,}216 \times 3.40 \times 10^{-6} \approx 57$ bad cells, far too many to ship without repair. At 6 sigma it expects 0.017. The count of failing cells is binomial and, for small $p$, close to Poisson with mean $Np$, which is the starting point for the redundancy analysis in the yield subdomain.

## From array yield back to a cell sigma

Run the conversion backwards to set a design target. If an array of $N$ cells without repair must work with probability $Y$, every cell must pass: $(1-p)^N = Y$, so $p = 1 - Y^{1/N} \approx (1 - Y)/N$. Then the required margin is $k = Q^{-1}(p)$. For a 64 Mb array ($N = 67{,}108{,}864$) and $Y = 0.99$, $p \approx 1.5 \times 10^{-10}$ and $k \approx 6.3$. That is where the familiar 6-sigma bitcell requirement comes from. The inverse must also be computed carefully: $Y^{1/N}$ is within $10^{-10}$ of 1, so compute $p$ as $-\mathrm{expm1}(\ln Y / N)$, then invert $Q$ by a bracketed root-finder.

## Going deeper

Abramowitz and Stegun, "Handbook of Mathematical Functions," section 26.2 on the normal distribution; Kanj, Joshi and Nassif, "Mixture importance sampling and its application to the analysis of SRAM designs in the presence of rare failure events," DAC 2006.
