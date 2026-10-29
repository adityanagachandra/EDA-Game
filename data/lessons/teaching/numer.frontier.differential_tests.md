# Differential tests for solver changes

Every change to a numerical engine, by an agent or a person, must be proven before promotion. For solvers the proof has a specific shape: backward-error and residual checks that are independent of conditioning, forward comparisons whose tolerances scale with the condition number, and golden circuits whose answers are known. Getting these right is what makes automated solver edits safe.

## Residual and backward-error checks

The normwise backward error of a computed solution $\hat x$ is
$$\eta = \frac{\|b - A\hat x\|}{\|A\|\,\|\hat x\| + \|b\|}.$$
It is the smallest relative perturbation of $A$ and $b$ for which $\hat x$ is exact, and it does not depend on $\kappa$. A backward-stable direct solver in double precision gives $\eta$ near $10^{-16}$ times a modest factor; an iterative solver gives $\eta$ near its stopping tolerance.

Example: $A = \begin{pmatrix} 4 & 1 \\ 1 & 3 \end{pmatrix}$, $b = (6, 7)$, exact $x = (1, 2)$, computed $\hat x = (1.0001, 1.9999)$. Then $r = b - A\hat x = (-0.0003, 0.0002)$, $\|r\|_\infty = 3 \times 10^{-4}$, $\|A\|_\infty = 5$, $\|\hat x\|_\infty = 1.9999$, $\|b\|_\infty = 7$: $\eta = 3 \times 10^{-4}/17.0 = 1.8 \times 10^{-5}$. For a direct solver that is a failure.

## Tolerances that respect conditioning

Two correct solvers can return different answers. If each meets a relative residual tolerance $\tau$, each forward error is bounded by about $\kappa\tau$, so their solutions may differ by up to about $2\kappa\tau$. With $\kappa = 10^7$ and $\tau = 10^{-10}$, a relative difference of $10^{-3}$ is consistent; a difference of $10^{-1}$ is not, and one solver is wrong.

So a differential test comparing solutions must scale its tolerance with a condition estimate, or it will flag correct changes on ill-conditioned cases and miss real bugs on well-conditioned ones. A residual or backward-error check needs no such scaling, which is why it is the first gate. Waveform comparisons follow the same logic: compare at common sample times, with reltol and abstol bands, because adaptive timestep selection differs between versions.

## Golden circuits and what must not move

A golden suite for a solver mixes cases with known answers and cases that stress numerics: an RC step response against the analytic exponential; a diode DC point against an independent high-precision solve; an IR-drop grid with an independently verified worst-node drop; deliberately ill-conditioned matrices; stiff and oscillatory transients to catch damping or ringing changes.

Beyond accuracy, the differential report records cost as QoR: Newton iterations, timesteps accepted and rejected, Krylov iterations, factorization fill and memory. A plug-in that matches every waveform but doubles Newton iterations is a regression to review, not a pass. Small designs are not enough: promotion needs the QoR baseline on representative large designs, where scaling problems in fill or iteration count appear.

## Going deeper

Higham, "Accuracy and Stability of Numerical Algorithms" (2nd ed., SIAM, 2002), chapters 1 and 7 on backward error; Rigal and Gaches, "On the compatibility of a given solution with the data of a linear system", J. ACM (1967).
