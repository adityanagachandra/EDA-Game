# Response surface methods

Response surface methodology fits a low-order polynomial to a few well-chosen simulations and then optimizes the polynomial instead of the circuit. It is the classical way to center a design, tune a process recipe, or build a cheap performance model over a region. The fitting is regression; the method is the loop of designing runs, fitting, locating a stationary point and moving the region.

## A quadratic through three coded runs

In one coded factor with runs at $x = -1, 0, +1$, the quadratic $y = b_0 + b_1x + b_2x^2$ passes exactly through the three responses $y_-, y_0, y_+$:
$$b_0 = y_0, \qquad b_1 = \frac{y_+ - y_-}{2}, \qquad b_2 = \frac{y_+ + y_- - 2y_0}{2}.$$
If a figure of merit is 52, 48 and 50 at $x = -1, 0, +1$, then $b_0 = 48$, $b_1 = -1$ and $b_2 = 3$. The stationary point solves $b_1 + 2b_2x = 0$: $x^* = -b_1/(2b_2) = 1/6 \approx 0.167$ in coded units. Since $b_2 > 0$ it is a minimum. Converting back to physical units uses the coding: $x = (p - p_{\text{centre}})/(\text{half range})$. Three runs determine three coefficients, so the fit has no residual and therefore no internal estimate of its own error; replicated centre runs and extra design points are what supply one in practice.

## The fitted optimum

Substituting $x^*$ gives the predicted optimum $y^* = b_0 - b_1^2/(4b_2) = 48 - 1/12 = 47.92$. Two checks keep this honest. First, the prediction is only as good as the quadratic approximation, so the optimum is confirmed with a real simulation at $x^*$. Second, a stationary point outside the coded region, $|x^*| > 1$, means the surface is extrapolating; the method then moves the experimental region toward the optimum (steepest descent on a first-order model) and refits, rather than trusting the far-away point. The sequence of moving the region and refitting is what makes RSM a search algorithm rather than a single fit.

## Two factors and the canonical analysis

With two factors the fitted surface is $y = b_0 + b_1x_1 + b_2x_2 + b_{11}x_1^2 + b_{22}x_2^2 + b_{12}x_1x_2$. Setting the gradient to zero gives a 2x2 linear system,
$$\begin{pmatrix} 2b_{11} & b_{12} \\ b_{12} & 2b_{22}\end{pmatrix}\begin{pmatrix} x_1 \\ x_2 \end{pmatrix} = -\begin{pmatrix} b_1 \\ b_2 \end{pmatrix},$$
solved by Cramer's rule. With $b_1 = -6$, $b_2 = 2$, $b_{11} = 3$, $b_{22} = 2$, $b_{12} = 1$ the system is $6x_1 + x_2 = 6$ and $x_1 + 4x_2 = -2$, with determinant $24 - 1 = 23$. Cramer's rule gives $x_1 = (6 \cdot 4 - 1 \cdot (-2))/23 = 26/23 = 1.13$ and $x_2 = (6 \cdot (-2) - 1 \cdot 6)/23 = -18/23 = -0.78$, with $x_1$ just outside the coded region. The eigenvalues of the Hessian classify the point as a minimum, maximum or saddle, and a near-zero eigenvalue reveals a ridge: a direction along which the response barely changes, often where the most robust design lies.

## Going deeper

Myers, Montgomery and Anderson-Cook, "Response Surface Methodology" (4th ed., 2016), chapters 2, 6 and 7.
