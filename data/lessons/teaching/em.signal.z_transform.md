# The Z transform and discrete-time stability

Simulators do not solve differential equations; they solve difference equations produced by an integration method. The Z transform is the tool for analyzing those: a discrete-time system is stable when its poles lie inside the unit circle. Mapping each integration method's pole for a simple decay $dv/dt = -v/\tau$ explains why some methods blow up with large steps and others ring.

## Difference equations and the Z transform

The Z transform $X(z) = \sum_n x[n] z^{-n}$ turns a delay by one sample into multiplication by $z^{-1}$. A first-order difference equation $y[n] = a\,y[n-1] + b\,x[n]$ becomes $Y(z)(1 - az^{-1}) = bX(z)$, so $H(z) = b/(1 - az^{-1})$ with a pole at $z = a$. The system is stable when $|a| < 1$.

Example: $a = 0.8$, $b = 0.2$, a unit step input and $y[-1] = 0$. Iterating: $y[0] = 0.2$, $y[1] = 0.36$, $y[2] = 0.488$, $y[3] = 0.590$. It converges to $b/(1 - a) = 1$, the DC gain $H(1)$. With $a = -0.8$ the output overshoots and undershoots alternately as it converges, the discrete signature of a pole on the negative real axis, and with $|a| > 1$ it grows without bound. In discrete time the unit circle plays the role the imaginary axis plays in continuous time: the continuous pole $s$ maps to $z = e^{sh}$ for a sampling step $h$, so the left half-plane maps inside the circle.

## Forward Euler and the stability limit

Forward Euler on $dv/dt = -v/\tau$ gives $v_{n+1} = v_n - (h/\tau)v_n = (1 - h/\tau)v_n$, a pole at $z = 1 - h/\tau$. The true solution always decays, but the numerical one decays only if $|1 - h/\tau| < 1$, which requires
$$h < 2\tau.$$
For $\tau < h < 2\tau$ the pole is negative and the solution oscillates in sign while decaying; beyond $2\tau$ it grows.

Example: $\tau = 10$ ps. A step of 5 ps gives $z = 0.5$, fine; 15 ps gives $z = -0.5$, oscillating; 25 ps gives $z = -1.5$, blowing up. Circuits contain time constants spread over many decades (stiff systems), and the smallest one would force a tiny step everywhere. That is why SPICE-class simulators use implicit methods.

## Backward Euler and trapezoidal poles

Backward Euler evaluates the derivative at the new point: $v_{n+1} = v_n - (h/\tau)v_{n+1}$, so $z = 1/(1 + h/\tau)$, inside the unit circle for every positive $h$. It is unconditionally stable but adds numerical damping: large steps over-damp oscillations. The trapezoidal rule averages both ends:
$$z = \frac{1 - h/(2\tau)}{1 + h/(2\tau)},$$
also stable for every step and second-order accurate.

Example: $\tau = 10$ ps and $h = 40$ ps. Backward Euler gives $z = 0.2$, a sensible fast decay. Trapezoidal gives $z = (1 - 2)/(1 + 2) = -0.333$: stable, but the sign alternates from step to step, the trapezoidal ringing seen in SPICE waveforms. As $h/\tau \to \infty$ its pole approaches -1. Simulators therefore mix methods, such as Gear or a backward-Euler step after a discontinuity.

## Going deeper

Najm, "Circuit Simulation", chapter 4; Oppenheim and Schafer, "Discrete-Time Signal Processing", chapter 3; Pillage, Rohrer and Visweswariah, "Electronic Circuit and System Simulation Methods", chapter 4.
