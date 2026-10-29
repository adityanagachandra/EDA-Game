# Reduction order and reproducibility

Floating-point addition is not associative, so a parallel sum gives slightly different answers depending on how it is split across threads. In an EDA engine those last-bit differences can change a Newton iteration count, a timestep decision or a placement move, making runs irreproducible. Designing reductions for determinism is part of building a trustworthy numerical kernel.

## Floating-point addition is not associative

IEEE single precision has a 24-bit significand, so the spacing of representable numbers near $2^{24} = 16777216$ is 2. Rounding is to nearest, ties to even. Then
$$(2^{24} + 1) + 1 = 2^{24} + 1 = 2^{24},$$
because each $+1$ lands exactly halfway and rounds to the even neighbour $2^{24}$. But $2^{24} + (1 + 1) = 2^{24} + 2 = 16777218$, exactly representable.

In double precision the same effect appears at $2^{53}$, and partially whenever magnitudes differ widely: adding a small current to a large one discards the small one's low-order bits. Order matters whenever terms of different magnitude are summed, which is every KCL residual and every dot product in a Krylov solver.

## Computing a sum in two orders

Sum $1, 1, 2^{24}, 1, 1$ in float32 left to right: $1 + 1 = 2$, $2 + 2^{24} = 16777218$, then $+1$ gives a tie between 16777218 and 16777220, resolved to even significand: 16777220; $+1$ again ties between 16777220 and 16777222 and rounds to 16777220. The result is 16777220, equal to the exact sum by luck of the ties.

With the order $2^{24}, 1, 1, 1, 1$: each $+1$ rounds back to $2^{24}$, and the result is 16777216, off by 4. A parallel reduction whose thread count or scheduling changes between runs takes different orders, so it can return either answer. GPU atomic additions make the order depend on timing and are inherently nondeterministic.

## Making parallel results reproducible

Error bounds depend on order too. Recursive left-to-right summation of $n$ terms has error at most about $(n - 1)u\sum|x_i|$; pairwise (tree) summation has at most about $\lceil\log_2 n\rceil u\sum|x_i|$. For $n = 10^6$ that is a factor $999999/20 = 50{,}000$ smaller worst-case bound.

Reproducibility needs a fixed order, not just accuracy. Common designs: a fixed reduction tree independent of thread count; deterministic partitioning of work with results combined in index order; compensated (Kahan) summation, which shrinks the dependence on order; and exact or reproducible accumulators that give bit-identical results regardless of order. Each costs some speed. Engines typically offer a deterministic mode for golden tests and signoff and a faster mode for exploration, and say clearly which is which.

## Going deeper

Higham, "Accuracy and Stability of Numerical Algorithms" (2nd ed., SIAM, 2002), chapter 4; Demmel and Nguyen, "Parallel reproducible summation", IEEE Trans. Computers (2015).
