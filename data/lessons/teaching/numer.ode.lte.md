# Local truncation error

The local truncation error (LTE) of an integration step is the error it would make starting from exact data. Simulators estimate it at every step to decide whether to accept the step and how large the next one should be. Its formula for each method, and how to estimate the derivatives it needs from computed points, is the heart of adaptive timestepping.

## Taylor expansion of one step

Expand the exact solution: $x(t_n + h) = x(t_n) + hx'(t_{n+1}) - \frac{h^2}{2}x''(\xi)$ around the end point. Backward Euler keeps only the first two terms, so its LTE is
$$\text{LTE}_{BE} \approx \frac{h^2}{2}\,x''.$$
With $|x''| = 4$ V/ns$^2$ and $h = 0.02$ ns, $\text{LTE} = 0.0002 \cdot 4 = 0.8$ mV.

The LTE scales like $h^{p+1}$ for a method of order $p$: halving the step cuts the backward Euler LTE by 4. Errors accumulate over about $T/h$ steps, so the global error is one order lower, $O(h^p)$. The LTE is local and estimable, which is why controllers work with it.

## Estimating the error from divided differences

The simulator does not know $x''$, so it estimates it from accepted points. With values $x_0, x_1, x_2$ at times $t_0, t_1, t_2$, the second divided difference is
$$x[t_0, t_1, t_2] = \frac{(x_2 - x_1)/(t_2 - t_1) - (x_1 - x_0)/(t_1 - t_0)}{t_2 - t_0},$$
and $x'' \approx 2\,x[t_0, t_1, t_2]$. Example: $0.20, 0.30, 0.38$ V at spacings 0.1 and 0.1 ns. First differences 1.0 and 0.8 V/ns; the divided difference is $-0.2/0.2 = -1$; $x'' \approx -2$ V/ns$^2$. A next backward Euler step of 0.1 ns would have LTE $\approx 0.005 \cdot 2 = 10$ mV.

SPICE applies this to capacitor charges and inductor fluxes, scaling by a trtol factor (default 7) that allows for the estimate's pessimism.

## Error constants of BE, trapezoidal and BDF2

Each method has its own leading term:
$$\text{LTE}_{BE} \approx \frac{h^2}{2}x'', \quad \text{LTE}_{TR} \approx -\frac{h^3}{12}x''', \quad \text{LTE}_{BDF2} \approx -\frac{2h^3}{9}x'''.$$
With $|x'''| = 5$ V/ns$^3$ and $h = 0.1$ ns: trapezoidal gives $0.001 \cdot 5/12 = 0.42$ mV, BDF2 gives $0.001 \cdot 10/9 = 1.11$ mV. Both are second order, but BDF2's constant is about 2.7 times larger, so for equal accuracy it needs a step about $2.7^{1/3} = 1.39$ times smaller.

The third derivative is estimated from a third divided difference: $x''' \approx 6\,x[t_0, t_1, t_2, t_3]$, which needs four accepted points. Mixing up these constants, or using a second-derivative estimate for a second-order method, makes the controller systematically too timid or too bold.

## Going deeper

Hairer, Norsett and Wanner, "Solving Ordinary Differential Equations I" (2nd ed., Springer, 1993), chapters II.4 and III.2; Nagel, "SPICE2" (UC Berkeley ERL Memo M520, 1975), on timestep control.
