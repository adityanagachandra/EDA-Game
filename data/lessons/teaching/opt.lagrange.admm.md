# ADMM and splitting

The alternating direction method of multipliers splits a problem into pieces that each own some variables, solves the pieces separately, and drives shared copies to agree through scaled dual variables. It is the augmented Lagrangian method with the joint minimization replaced by one pass of block coordinate minimization. In EDA it suits partitioned or distributed formulations: placement split across regions that share boundary cells, or problems that combine an easy smooth term with an easy constraint set.

## The x-update

Consensus form: minimize $f(x) + g(z)$ subject to $x = z$, where each partition owns one copy. The scaled ADMM iteration is
$$x \leftarrow \arg\min_x f(x) + \tfrac{\rho}{2}(x - z + u)^2,\quad z \leftarrow \arg\min_z g(z) + \tfrac{\rho}{2}(x - z + u)^2,\quad u \leftarrow u + x - z.$$
Take $f = (x - a)^2$ and $g = (z - b)^2$: partition 1 wants the boundary cell at a, partition 2 at b. The x-update solves $2(x - a) + \rho(x - z + u) = 0$:
$$x = \frac{2a + \rho(z - u)}{2 + \rho}.$$
With $a = 36$, $\rho = 1$, $z_0 = 28$ and $u_0 = 0$: $x_1 = (72 + 28)/3 = 33.33$. Partition 1 solves only its own small problem, pulled partly toward the other copy.

## The z-update with the new x

The second block uses the x just computed, which is what makes ADMM a Gauss-Seidel pass rather than a Jacobi pass. For $g = (z - b)^2$, the z-update solves $2(z - b) - \rho(x - z + u) = 0$:
$$z = \frac{2b + \rho(x + u)}{2 + \rho}.$$
Example: $a = 27$, $b = 5$, $\rho = 4$, $z_0 = 4$, $u_0 = 0$. First $x_1 = (54 + 16)/6 = 11.667$. Then $z_1 = (10 + 46.667)/6 = 9.444$. Each copy has moved toward the other, but they still disagree by 2.222.

In a distributed placer the two updates run on different partitions, perhaps on different machines, exchanging only the shared boundary coordinates. When f and g are complicated but separately easy, for example a smooth wirelength term and a projection onto legal rows, ADMM lets each be handled by its own specialized solver.

## The dual update and convergence

The scaled dual accumulates the disagreement: $u \leftarrow u + x - z$. Example: $a = 30$, $b = 22$, $\rho = 1$, $z_0 = 32$. Then $x_1 = (60 + 32)/3 = 30.667$, $z_1 = (44 + 30.667)/3 = 24.889$, and $u_1 = 0 + 30.667 - 24.889 = 5.778$. In the next iteration u shifts both copies toward agreement; at convergence $x = z = (a + b)/2 = 26$, and $\rho u$ equals the multiplier of the consensus constraint.

ADMM converges for convex f and g under mild conditions, for any $\rho > 0$, but often slowly to high accuracy; it is usually run to moderate accuracy and polished by another method. The two residuals, $x - z$ (primal) and $\rho$ times the change in z (dual), are monitored, and $\rho$ is adjusted to balance them. For nonconvex pieces, as in legalization, ADMM is a heuristic that works well in practice without a general guarantee.

## Going deeper

Boyd, Parikh, Chu, Peleato and Eckstein, "Distributed optimization and statistical learning via the alternating direction method of multipliers" (Foundations and Trends in Machine Learning, 2011).
