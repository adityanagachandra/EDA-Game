# Model-based OPC iteration

Model-based OPC is a numerical loop: simulate the print, measure edge placement errors, move mask fragments, repeat. Whether it converges, how fast, and how its errors are summarized are questions of fixed-point iteration. A tool designer who builds or tunes OPC engines decides the fragmentation, the feedback gain and the stopping rule, and those choices set runtime and quality across billions of fragments.

## Measuring edge placement error

Each fragment has a control site. The simulator computes the printed contour there, and the edge placement error is printed position minus target position, signed so that positive means the print lies outside the target. A polygon's quality is summarized by statistics over its fragments. For EPEs of 3, -2, 5 and -4 nm, the sum of squares is $9 + 4 + 25 + 16 = 54$, so the RMS EPE is $\sqrt{54/4} = 3.67$ nm, and the maximum magnitude is 5 nm. Both matter: RMS tracks overall fidelity, but a single bad fragment can still cause a bridge or an open, so tools also stop on maximum EPE.

## A feedback update for each fragment

Near its current position, moving a fragment by $b$ nm (at wafer scale) moves the printed edge by about $g b$, where $g$ is a local gain related to MEEF. OPC updates each fragment by $-f \cdot EPE$, with damping factor $f$. The new error is $EPE_{k+1} = EPE_k + g(-f\, EPE_k) = EPE_k(1 - f g)$. With $EPE_0 = 8$ nm, $g = 1.2$ and $f = 0.5$, the factor is $1 - 0.6 = 0.4$ and the error becomes 3.2 nm. If $f g = 1$ the correction is exact in one step; if $f g > 2$ the update overshoots and diverges.

## Convergence and stopping

The update is a fixed-point iteration with contraction factor $r = |1 - f g|$, so $|EPE_n| = r^n |EPE_0|$. The number of iterations to reach tolerance $\epsilon$ is the smallest $n$ with $r^n |EPE_0| \le \epsilon$. From 12 nm with $r = 0.4$ and $\epsilon = 0.5$ nm: 4.8, 1.92, 0.77, 0.31, so 4 iterations. In practice the gain differs by fragment and neighbouring fragments interact, so engines use conservative damping, per-fragment gains, and caps on iteration count. Each iteration re-simulates the image near every fragment, so cutting one iteration saves a large fraction of full-chip OPC runtime.

## Going deeper

Cobb, "Fast Optical and Process Proximity Correction Algorithms" (1998); Burden and Faires, "Numerical Analysis", the chapter on fixed-point iteration.
