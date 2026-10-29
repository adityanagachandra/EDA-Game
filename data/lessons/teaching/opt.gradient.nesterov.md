# Nesterov's accelerated gradient

Nesterov's accelerated gradient adds a look-ahead extrapolation to gradient descent and improves the worst-case convergence on smooth convex problems from order $1/k$ to order $1/k^2$. ePlace made it the standard engine of electrostatic global placement, and RePlAce in OpenROAD and DREAMPlace follow the same pattern. The method is two lines of arithmetic, but the order of those lines, extrapolate first and then take the gradient at the extrapolated point, is what makes it work.

## The look-ahead point

Nesterov's method keeps the last two iterates. It first extrapolates along the recent direction of motion,
$$y_k = x_k + \beta(x_k - x_{k-1}),$$
with momentum $\beta \in [0, 1)$, and then takes a gradient step from $y_k$. With $x_k = 17$, $x_{k-1} = 19$ and $\beta = 0.8$: $y_k = 17 + 0.8 \cdot (-2) = 15.4$. The iterate was moving left by 2, so the method guesses it will keep moving and looks 1.6 further left.

In a placer every cell coordinate gets its own look-ahead, computed in one vector operation, and the density and wirelength gradients are then evaluated at the extrapolated placement. The momentum coefficient is not a free constant in the classical method: it follows a schedule, $\beta_k = (a_k - 1)/a_{k+1}$ with $a_{k+1} = (1 + \sqrt{4a_k^2 + 1})/2$, which starts near 0 and approaches 1.

## The gradient step from the look-ahead

The new iterate is $x_{k+1} = y_k - \eta \nabla f(y_k)$: the gradient is evaluated at $y_k$, not at $x_k$. Heavy-ball momentum, by contrast, uses the gradient at $x_k$ and then adds momentum; the difference is small per step but decisive for the worst-case rate.

Example: $f(x) = 2(x - 7)^2$ so $f'(x) = 4(x - 7)$, with $\eta = 0.25$, $\beta = 0.3$, $x_k = 17$, $x_{k-1} = -3$. Look ahead: $y_k = 17 + 0.3 \cdot 20 = 23$. Gradient there: $4 \cdot 16 = 64$. Step: $x_{k+1} = 23 - 0.25 \cdot 64 = 7$. Here $\eta = 1/L$ with $L = 4$, so the step from $y_k$ lands exactly on the minimum of this one-dimensional quadratic. In many dimensions no single step does that, and the extrapolation is what carries the iterate quickly along the flat directions.

## Running the iteration

Start with $x_{-1} = x_0$, so the first look-ahead is $y_0 = x_0$ and the first step is plain gradient descent. Take $f(x) = 1.5(x + 1)^2$, $f'(x) = 3(x + 1)$, $\eta = 0.05$, $\beta = 0.9$, $x_0 = -2$. Then $x_1 = -2 - 0.05 \cdot (-3) = -1.85$. Next $y_1 = -1.85 + 0.9(-1.85 + 2) = -1.715$, $f'(y_1) = -2.145$, and $x_2 = -1.715 + 0.107 = -1.608$. Plain descent would have reached only -1.7225.

Acceleration has costs. The iterates are not monotone, since momentum can overshoot, so placers monitor the objective and restart momentum when it rises. The step must stay below $1/L$, and ePlace estimates L each iteration from the change in gradient between successive points, $\|\nabla f(y_k) - \nabla f(y_{k-1})\| / \|y_k - y_{k-1}\|$, which makes the method adaptive without a line search.

## Going deeper

Nesterov, "A method for solving the convex programming problem with convergence rate O(1/k^2)" (Soviet Mathematics Doklady, 1983); Lu et al., "ePlace: electrostatics-based placement using fast Fourier transform and Nesterov's method" (ACM TODAES, 2015); Cheng, Kahng, Kang and Wang, "RePlAce: advancing solution quality and routability validation in global placement" (IEEE TCAD, 2019).
