# Truncation by singular values

Krylov methods match moments; singular-value methods rank directions by importance and throw away the weak ones with a computable error bound. Truncated SVD, proper orthogonal decomposition and balanced truncation all follow this pattern. They are used to compress thermal models, post-process Krylov models and build reduced models from simulation snapshots.

## Singular values rank the important directions

Any matrix factors as $A = U\Sigma V^T$ with singular values $\sigma_1 \ge \sigma_2 \ge \dots \ge 0$. Collect simulated temperature fields as columns of a snapshot matrix; its leading left singular vectors are the dominant spatial patterns, and the energy captured by the first $k$ is
$$\frac{\sum_{i \le k}\sigma_i^2}{\sum_i\sigma_i^2}.$$
With singular values 10, 4, 1, 0.5: total energy $100 + 16 + 1 + 0.25 = 117.25$; one mode captures 85.3%, two capture 98.9%.

Proper orthogonal decomposition (POD) uses those vectors as a reduced basis and projects the thermal equations onto them, like a Krylov basis. Rapid decay of singular values is what makes a problem compressible.

## The Eckart-Young error of a truncation

Keeping the first $k$ singular triplets gives $A_k$, the best rank-$k$ approximation (Eckart-Young-Mirsky), with errors
$$\|A - A_k\|_2 = \sigma_{k+1}, \qquad \|A - A_k\|_F = \sqrt{\textstyle\sum_{i > k}\sigma_i^2}.$$
For singular values 10, 4, 1, 0.5 and $k = 2$: the 2-norm error is 1, relative error $1/10 = 0.1$, and the Frobenius error is $\sqrt{1.25} = 1.118$.

No other rank-2 matrix does better in either norm. That optimality is what makes SVD truncation a reference against which other compression is measured, but it is optimal for approximating a matrix, not a dynamical system's input-output map, which is what balanced truncation addresses.

## Choosing the reduced order for a tolerance

Balanced truncation works with the Hankel singular values of a stable system, which measure how controllable and observable each state is. Truncating to $k$ states gives a stable reduced model with
$$\|H - H_k\|_\infty \le 2\sum_{i > k}\sigma_i.$$
With Hankel singular values 5, 2, 0.6, 0.15, 0.05 and tolerance 0.5: keeping 2 states gives a bound of $2(0.8) = 1.6$; keeping 3 gives $2(0.2) = 0.4 \le 0.5$. So $k = 3$.

The bound is a priori, which a Krylov method lacks. The catch is cost: balanced truncation needs two Lyapunov equations, $O(n^3)$ in dense form, so it is applied to small or medium models, often after a Krylov method has already reduced a large network to a few hundred states.

## Going deeper

Antoulas, "Approximation of Large-Scale Dynamical Systems" (SIAM, 2005), chapters 3 and 7; Moore, "Principal component analysis in linear systems: controllability, observability, and model reduction", IEEE Trans. Automatic Control (1981).
