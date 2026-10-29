# BDF2 and Gear methods

Backward differentiation formulas (BDF), known in circuit simulation as Gear's methods, generalize backward Euler by using several past points. BDF2 is second order, A-stable and, unlike the trapezoidal rule, damps stiff modes. SPICE's Gear option and many DAE solvers use variable-order BDF, so its formula, start-up and limits are core knowledge.

## Differentiating an interpolant through past points

BDF methods fit a polynomial through the newest points and differentiate it at the new time. With equal steps, the quadratic through $y_{n-1}, y_n, y_{n+1}$ has derivative at $t_{n+1}$
$$y'(t_{n+1}) \approx \frac{3y_{n+1} - 4y_n + y_{n-1}}{2h}.$$
Check on $y = t^2$ at $t = 2$ with $h = 1$: $(3 \cdot 4 - 4 \cdot 1 + 0)/2 = 4$, exactly $y'(2)$, because the formula is exact for quadratics.

BDF$k$ uses $k + 1$ points and is order $k$; BDF1 is backward Euler. Setting this derivative equal to $f(y_{n+1})$ gives an implicit equation, solved by Newton in a circuit simulator exactly as for backward Euler.

## One BDF2 step on a linear test equation

BDF2 with equal steps is
$$y_{n+1} = \frac{4y_n - y_{n-1} + 2h\,f(y_{n+1})}{3}.$$
For $y' = -ay + b$:
$$y_{n+1} = \frac{4y_n - y_{n-1} + 2hb}{3 + 2ah}.$$
With $a = 2$, $b = 4$, $h = 0.3$, $y_{n-1} = 5$, $y_n = 3.65$: $y_{n+1} = (14.6 - 5 + 2.4)/4.2 = 2.857$. The exact solution through those points continues toward 2.

On $y' = \lambda y$ the BDF2 amplification tends to 0 as $h\lambda \to -\infty$, so stiff components are damped, not reflected as with the trapezoidal rule. The price: its error constant is larger than trapezoidal's, so for the same accuracy it needs somewhat smaller steps.

## Starting the method and its order limits

BDF2 needs two past values, so the first step uses backward Euler. From $y_0 = 5$ with $a = 2$, $b = 4$, $h = 0.3$: backward Euler gives $y_1 = 6.2/1.6 = 3.875$; BDF2 then gives $y_2 = (15.5 - 5 + 2.4)/4.2 = 3.071$. Variable-step BDF needs coefficients recomputed from unequal spacings, and order changes need history, which is why BDF codes store several past solutions.

Two theorems limit the family. BDF methods of order 1 to 6 are zero-stable, and order 7 and above are not, so codes stop at 6. Dahlquist's second barrier says an A-stable linear multistep method has order at most 2, so only BDF1 and BDF2 are A-stable; BDF3 to BDF6 lose stability near the imaginary axis, which hurts on weakly damped oscillators. Simulators therefore cap the order for oscillatory circuits.

## Going deeper

Gear, "Numerical Initial Value Problems in Ordinary Differential Equations" (Prentice-Hall, 1971); Hairer and Wanner, "Solving Ordinary Differential Equations II" (2nd ed., Springer, 1996), chapter V.
