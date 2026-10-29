# Model order reduction by Krylov projection

Krylov projection builds reduced interconnect models that match moments without ever forming them explicitly. An orthonormal basis of a Krylov subspace, computed by Arnoldi or Lanczos, is used to project the full system matrices onto a few dimensions. PRIMA, the standard passive reduction method for RLC interconnect, is built this way.

## Projecting onto a Krylov subspace

Write $A = -G^{-1}C$ and $r = G^{-1}b$; the moment vectors are $A^kr$. The Krylov subspace
$$K_q(A, r) = \text{span}\{r, Ar, \dots, A^{q-1}r\}$$
contains the first $q$ moment vectors. With $V$ an $n \times q$ orthonormal basis of it, the reduced model is
$$G_r = V^TGV, \quad C_r = V^TCV, \quad b_r = V^Tb, \quad l_r = V^Tl,$$
a $q$-state system. Example of one Krylov step: with $A = \begin{pmatrix} -2 & 1 \\ 1 & -3 \end{pmatrix}$ and $r = (1, 1)$, $Ar = (-1, -2)$, and $K_2 = \text{span}\{(1, 1), (-1, -2)\}$.

Each new basis vector costs one solve with the factored $G$ and an orthogonalization, so building a model with $q = 20$ for a million-node net costs about twenty sparse solves.

## How many moments a projection matches

A one-sided projection whose basis spans $K_q$ matches the first $q$ moments of $H(s)$ at the expansion point. Two-sided methods, which use a left Krylov space from $l$ as well (Padé via Lanczos, PVL, by Feldmann and Freund), match $2q$ moments with a $q$-state model, the Padé optimum.

Expansion points need not be zero. Rational Krylov methods use vectors $(G + s_iC)^{-1}\dots$ at several points $s_i$; a one-sided basis with $q_1$ vectors at $s_1$ and $q_2$ at $s_2$ matches $q_1$ moments about $s_1$ and $q_2$ about $s_2$. Choosing points across the band of interest gives accuracy to high frequencies, which matters for clock and RF nets. Multiport nets use block Krylov spaces, growing the model by one vector per port per block step.

## Arnoldi orthogonalization and passivity

Computing $r, Ar, A^2r, \dots$ explicitly is the power method: the vectors align with the dominant eigenvector, the basis becomes numerically rank-deficient, and higher moments are lost. That is AWE's failure mode. Arnoldi orthogonalizes each new vector against the previous ones before normalizing, so the basis stays well conditioned and the same moments are matched stably.

PRIMA (Odabasioglu, Celik and Pileggi, 1998) adds a structural guarantee. For RLC networks in a suitable MNA form, $G + G^T$ and $C$ are positive semidefinite. The congruence $V^T(\cdot)V$ preserves semidefiniteness, so the reduced model is passive: it cannot generate energy, and connecting it to any passive circuit keeps the simulation stable. Two-sided Padé models match more moments but can lose passivity, which is why PRIMA-style congruence is the default in extraction-to-simulation flows.

## Going deeper

Odabasioglu, Celik and Pileggi, "PRIMA: passive reduced-order interconnect macromodeling algorithm", IEEE Trans. CAD (1998); Feldmann and Freund, "Efficient linear circuit analysis by Padé approximation via the Lanczos process", IEEE Trans. CAD (1995).
