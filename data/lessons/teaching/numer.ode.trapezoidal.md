# The trapezoidal rule for ODEs

The trapezoidal rule averages the slopes at both ends of a step. It is second-order accurate and A-stable, so it is SPICE's default transient method. Its weakness, a stiff error mode that flips sign every step instead of decaying, produces the familiar trapezoidal ringing, and every simulator designer must know when it appears and how to suppress it.

## Averaging the slopes at both ends

The trapezoidal rule integrates $y' = f$ over one step with the trapezoid quadrature:
$$y_{n+1} = y_n + \frac{h}{2}\big(f(y_n) + f(y_{n+1})\big).$$
For $y' = -ay + b$:
$$y_{n+1} = \frac{(1 - ah/2)\,y_n + hb}{1 + ah/2}.$$
With $a = 2$, $b = 4$, $y_0 = 5$, $h = 0.3$: $y_1 = (0.7 \cdot 5 + 1.2)/1.3 = 3.6154$. The exact solution is $2 + 3e^{-0.6} = 3.6464$, closer than backward Euler's 3.875.

The local error is $O(h^3)$ and the global error $O(h^2)$. Like backward Euler, the rule is implicit, so a nonlinear circuit needs a Newton solve per step. It also needs $f(y_n)$, which for a capacitor means the capacitor current at the previous point, stored as history.

## One trapezoidal step for an RC circuit

For the RC low-pass $\tau v' = u - v$ with constant input $u$ over the step:
$$v_{n+1} = \frac{(1 - h/2\tau)\,v_n + (h/\tau)\,u}{1 + h/2\tau}.$$
With $\tau = 1$ ns, $v_0 = 0$, $u = 1$ V and $h = 0.25$ ns: $v_1 = 0.25/1.125 = 0.2222$ V, against the exact 0.2212 V. The error is 0.001 V, twenty times smaller than backward Euler's at the same step.

That accuracy lets a simulator take larger steps for the same tolerance, which is why trapezoidal is the default in SPICE. The timestep controller's error estimate for this method scales like $h^3 v'''$, and the step size adapts accordingly.

## Trapezoidal ringing on stiff modes

On $y' = \lambda y$ the trapezoidal factor is $R(z) = (1 + z/2)/(1 - z/2)$, $z = h\lambda$. For a stiff mode with $z$ large and negative, $R \to -1$. Take a parasitic time constant of 2 ps under a 100 ps step: $z = -50$, $R = -24/26 = -0.923$. An error of 1 mV in that mode becomes -0.923, then 0.852, then -0.787 mV: it alternates sign and decays slowly, although the true mode died out within picoseconds.

This shows up as point-to-point oscillation in node voltages or capacitor currents, typically right after a sharp edge or a breakpoint. Remedies: take one or a few backward Euler steps after each breakpoint (SPICE does), use an L-stable method such as BDF2 or TR-BDF2, or cut the step locally. Ringing is a numerical artifact, not circuit behaviour, and a designer reading waveforms must recognise it.

## Going deeper

Hairer and Wanner, "Solving Ordinary Differential Equations II" (2nd ed., Springer, 1996), chapter IV.3; Bank et al., "Transient simulation of silicon devices and circuits", IEEE Trans. CAD (1985), on TR-BDF2.
