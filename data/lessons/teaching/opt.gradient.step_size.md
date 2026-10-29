# Step size, smoothness and line search

The step size decides whether gradient descent converges, how fast, and how robustly. For smooth problems the right scale is set by the Lipschitz constant of the gradient, the largest curvature L. Fixed steps below $2/L$ converge; the best fixed step balances the stiffest and flattest directions; and backtracking line search finds a safe step automatically. Placers estimate L on the fly, and getting it wrong is a common cause of oscillating or stalled global placement.

## The stability limit two over L

For a quadratic with Hessian H, each eigen-direction contracts by $1 - \eta\lambda_i$ per step. Convergence in every direction needs $|1 - \eta\lambda_i| < 1$, so $\eta < 2/\lambda_{\max}$, where $L = \lambda_{\max}$ is also the Lipschitz constant of the gradient.

For a 2 by 2 symmetric $H = \begin{pmatrix} a & h \\ h & b \end{pmatrix}$ the eigenvalues are $\frac{a + b}{2} \pm \sqrt{\left(\frac{a - b}{2}\right)^2 + h^2}$. Example: $a = b = 6$ and $h = 2$ give $6 \pm 2$, eigenvalues 4 and 8. With $L = 8$ the limit is $\eta < 0.25$. A step of 0.3 makes the stiff direction's factor $1 - 2.4 = -1.4$, and the iterates blow up along it, even though the soft direction's factor $1 - 1.2 = -0.2$ is fine.

## The best fixed step and conditioning

With eigenvalues between $\mu$ (smallest) and $L$ (largest), the worst-case factor over all directions is $\max(|1 - \eta\mu|, |1 - \eta L|)$. It is minimized when the two are equal and opposite: $1 - \eta\mu = -(1 - \eta L)$, so
$$\eta^* = \frac{2}{\mu + L}, \qquad \rho = \frac{L - \mu}{L + \mu} = \frac{\kappa - 1}{\kappa + 1},$$
with condition number $\kappa = L/\mu$. For $\mu = 2$ and $L = 38$: $\eta^* = 0.05$ and $\rho = 0.9$, so each step removes only 10% of the error, and reducing the error a millionfold takes about 130 steps.

Placement systems are badly conditioned: a cell on a short heavy net and a cell on a long light net have curvatures differing by orders of magnitude. Preconditioning, dividing each coordinate's gradient by its own curvature estimate (for example its pin count plus a density term), is the standard fix and is used by ePlace and its descendants.

## Backtracking line search

When L is unknown, backtracking finds a step: start at $\eta = 1$ and halve until the Armijo condition holds,
$$f(x - \eta g) \le f(x) - c\,\eta\,\|g\|^2,$$
with a constant $0 < c < 1$, often $10^{-4}$ in solvers. It demands a decrease at least a fraction c of what the linear model predicts.

For $f(x) = \tfrac{a}{2}x^2$ the condition reduces to $(1 - \eta a)^2 \le 1 - 2c\eta a$, that is $\eta \le 2(1 - c)/a$. With $a = 11$ and $c = 0.4$ the bound is 0.109; halving from 1 gives 0.5, 0.25, 0.125, then 0.0625, accepted after 4 halvings. Each trial costs a function evaluation, which in placement means recomputing wirelength and density, so placers often predict the step from the previous iteration's curvature (as ePlace does with a Lipschitz estimate) and backtrack only when the prediction fails.

## Going deeper

Nocedal and Wright, "Numerical Optimization" (2nd ed., 2006), chapter 3 on line search; Nesterov, "Lectures on Convex Optimization" (2nd ed., 2018), chapter 2 on smooth convex functions.
