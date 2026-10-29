# Posynomials and geometric programming in sizing

Delays under the Elmore model are posynomials in the gate and wire sizes: sums of positive terms that are products of powers of the variables. Minimizing a posynomial subject to posynomial constraints is a geometric program, which becomes convex after a logarithmic change of variables. That is why continuous gate and wire sizing has a unique global optimum, a fact exploited since TILOS, and why the optimum balances competing terms in a predictable way.

## Balancing a size against its load

Scaling a gate by x divides its drive resistance by x but multiplies its input capacitance by x, which loads the previous stage. A one-gate model is $D(x) = ax + b/x$ with $a, b > 0$: the first term is the upstream delay increase, the second the gate's own delay. Setting $D'(x) = a - b/x^2 = 0$ gives
$$x^* = \sqrt{b/a}.$$
With $a = 1.5$ and $b = 54$, $x^* = \sqrt{36} = 6$.

The second derivative $2b/x^3$ is positive for $x > 0$, so this stationary point is a minimum. The sizing intuition: grow a gate until the extra load it presents upstream costs as much as the drive it gains. A larger downstream load b or a cheaper input capacitance a both push the optimal size up as a square root, not linearly.

## Equal terms at the optimum

At the optimum of $ax + b/x$ the two terms are equal: $a x^* = \sqrt{ab} = b/x^*$. So the minimum delay is
$$D^* = 2\sqrt{ab}.$$
With $a = 4.5$ and $b = 46$: $D^* = 2\sqrt{207} = 28.77$ ps, reached at $x^* = 3.197$.

This equal-terms property is the arithmetic-geometric mean inequality: $ax + b/x \ge 2\sqrt{ax \cdot b/x} = 2\sqrt{ab}$, with equality exactly when the terms are equal. It is the one-dimensional case of geometric-programming duality, and it explains logical-effort results such as equal stage efforts along an optimally sized chain. Equality conditions of this kind give closed-form answers wherever a posynomial has as many terms as there are variables plus one.

## Log transform and multivariable posynomials

Substituting $x = e^u$ turns each monomial $c\, x^{p} y^{q}$ into $c\, e^{pu + qv}$, an exponential of an affine function. A posynomial becomes a sum of such exponentials, which is convex in $(u, v)$, and so is its logarithm. Every stationary point is therefore global, and interior-point GP solvers handle thousands of gates.

For $D(x, y) = ax + by + c/(xy)$ the optimum again equalizes the terms: setting both partial derivatives to zero gives $ax = c/(xy)$ and $by = c/(xy)$. Call the common value t. Multiplying the three terms gives $t^3 = abc$, so $D^* = 3(abc)^{1/3}$. With $a = b = 6$ and $c = 51$: $t^3 = 1836$, $t = 12.245$ and $D^* = 36.73$ ps. Real libraries offer discrete sizes, so the continuous GP solution is rounded and repaired, and its value serves as a lower bound on what discrete sizing can reach.

## Going deeper

Fishburn and Dunlop, "TILOS: a posynomial programming approach to transistor sizing" (ICCAD, 1985); Boyd, Kim, Vandenberghe and Hassibi, "A tutorial on geometric programming" (Optimization and Engineering, 2007).
