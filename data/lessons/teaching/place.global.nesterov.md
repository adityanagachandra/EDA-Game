# Nesterov's method with Lipschitz step estimates

ePlace solves the penalised placement problem with Nesterov's accelerated gradient method, whose momentum gives faster convergence than plain gradient descent on smooth problems. Because the Lipschitz constant of the placement objective is unknown and changes as $\lambda$ grows, ePlace predicts the step length from the last two iterates. This combination is the optimiser inside RePlAce, OpenROAD's gpl and DREAMPlace.

## Step length from a local Lipschitz estimate

For a function with $L$-Lipschitz gradient, a safe gradient step is $1/L$. ePlace estimates $L$ from the latest two points:
$$\alpha_k = \frac{\|v_k - v_{k-1}\|}{\|\nabla f(v_k) - \nabla f(v_{k-1})\|}.$$
If the position moved by (4, 6), norm 7.2111, and the gradient changed by (-2, 11), norm 11.1803, then $\alpha_k = 0.645$. This is a secant estimate of the inverse curvature, related to Barzilai-Borwein steps. ePlace adds backtracking: if the step computed after moving predicts a much smaller step than the one used, the step is recomputed, so a bad estimate does not destabilise the run. Vectors here are all cell coordinates stacked, so the norms are over millions of entries.

## The momentum sequence

Nesterov's method maintains a scalar sequence $a_0 = 1$, $a_{k+1} = (1 + \sqrt{4a_k^2 + 1})/2$, giving 1, 1.618, 2.194, 2.750, 3.295 and so on, growing roughly like $k/2$. The momentum coefficient is $(a_k - 1)/a_{k+1}$: 0 at the first step, 0.2817 at the second, and approaching 1 later. From $a_k = 3.2949$, $a_{k+1} = 3.8326$ and the coefficient is 0.5988. Increasing momentum is what accelerates convergence on smooth convex problems from $O(1/k)$ to $O(1/k^2)$. Placement is not convex, but the acceleration still helps, and the schedule can be restarted if the objective rises. The sequence depends only on the iteration count, not on the problem, so it costs nothing to compute; what makes the method work is pairing it with a good step length at every iteration.

## One full Nesterov update

Each iteration takes a gradient step from the look-ahead point and then extrapolates:
$$u_{k+1} = v_k - \alpha_k\nabla f(v_k), \qquad v_{k+1} = u_{k+1} + \frac{a_k - 1}{a_{k+1}}(u_{k+1} - u_k).$$
For one coordinate with $u_k = 75$, $v_k = 70$, gradient 6, $\alpha = 0.1$ and $a_k = 2.7486$: $u_{k+1} = 69.4$, $a_{k+1} = 3.2937$, the coefficient is 0.5309 and $v_{k+1} = 69.4 + 0.5307(69.4 - 75) = 66.427$. The gradient is evaluated at the extrapolated point $v$, not at $u$. Note the order of work: the step is taken from the look-ahead point, and the momentum term then reuses the difference between the new and old gradient-step points, so the method stores two position vectors rather than one. Per iteration the cost is one wirelength gradient, one density solve and a few vector operations, all parallel across cells.

## Going deeper

Nesterov, "A method of solving a convex programming problem with convergence rate O(1/k^2)" (Soviet Math. Doklady, 1983); Lu et al., "ePlace" (ACM TODAES, 2015), section on Nesterov's method with Lipschitz prediction.
