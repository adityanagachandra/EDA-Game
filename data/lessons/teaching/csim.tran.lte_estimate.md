# Estimating local truncation error

A simulator cannot compare its result with the exact solution, so it estimates the local truncation error (LTE) of each step from the solution itself. The estimate uses divided differences of recent accepted points to approximate the derivative that appears in the method's error term. This estimate drives every step-size decision in transient analysis.

## Derivatives from divided differences

Divided differences of accepted samples estimate derivatives. For equally spaced points $v_0, v_1, v_2$, the second divided difference is $(v_2 - 2v_1 + v_0)/(2h^2)$, and twice it estimates the second derivative:
$$x'' \approx \frac{v_2 - 2v_1 + v_0}{h^2}.$$
Similarly $x''' \approx (v_3 - 3v_2 + 3v_1 - v_0)/h^3$ from four points. For samples 1.0, 1.2, 1.3 V at $h = 0.5$ ns, $x'' \approx (1.3 - 2.4 + 1.0)/0.25 = -0.4$ V/ns$^2$. With unequal steps the general divided-difference recursion is used. Because these are differences of nearly equal numbers, they need the samples to many digits, and they are meaningless across a discontinuity.

## Backward Euler error estimate

Backward Euler's local error is $\tfrac{h^2}{2} x''$ in leading order. Substituting the divided-difference estimate,
$$\text{LTE}_{BE} \approx \frac{h^2}{2} \cdot \frac{v_2 - 2v_1 + v_0}{h^2} = \frac{v_2 - 2v_1 + v_0}{2}.$$
The step size cancels: the error of one step is half the second difference of the samples. For samples 0.1535, 0.2212, 0.2835 V, the second difference is -0.0054 V, and the LTE estimate is 2.7 mV. In a simulator the estimate is formed for every node or every reactive element, and the largest relative to its tolerance decides the step. A first-order method's estimate needs only three points, so it is available one step after a restart.

## Trapezoidal error estimate

The trapezoidal rule's local error is $\tfrac{h^3}{12} x'''$. With the third-difference estimate,
$$\text{LTE}_{TR} \approx \frac{v_3 - 3v_2 + 3v_1 - v_0}{12}.$$
For samples 0.786939, 1.186861, 1.454936, 1.634633 V the third difference is 0.04347 V and the estimate is 3.62 mV. SPICE implements this through charges: it estimates the derivative of capacitor charge from divided differences and compares the implied error with a tolerance built from reltol, abstol and chgtol, scaled by trtol, which defaults to 7 because the raw estimate is known to be pessimistic. A higher-order method needs higher differences, one more past point per order.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 5; Nagel, "SPICE2" (Berkeley ERL Memo M520, 1975), on timestep control.
