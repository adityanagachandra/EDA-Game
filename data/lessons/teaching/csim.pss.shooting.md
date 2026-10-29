# Shooting for periodic steady state

Many circuits are judged in periodic steady state: oscillators, mixers, switching converters, clocked comparators. A plain transient must wait for start-up transients to die out, which for a high-Q circuit can take thousands of periods. Shooting methods instead solve directly for the initial state that repeats after one period, using Newton on a boundary-value problem whose function evaluation is one period of transient integration.

## The periodic boundary condition

Let $\phi_T(x_0)$ be the state reached by integrating the circuit for one period $T$ from the state $x_0$. Periodic steady state requires $\phi_T(x_0) = x_0$. For a linear circuit driven by a periodic source, the one-period map is affine: $\phi_T(x_0) = a x_0 + b$, where $a$ describes how the initial state decays over a period and $b$ is the response to the drive. The fixed point is $x_0 = b/(1 - a)$. With $a = 0.75$ and $b = 1$ V, $x_0 = 4$ V: starting there, the circuit repeats exactly. For a vector state, $a$ is the state-transition matrix and $x_0 = (I - A)^{-1} b$.

## One shooting Newton step

For a nonlinear circuit, apply Newton to $r(x_0) = \phi_T(x_0) - x_0 = 0$. The Jacobian is $\partial\phi_T/\partial x_0 - I$, the sensitivity of the final state to the initial state minus the identity. That sensitivity comes almost free: it is the product of the per-step Jacobians the transient already factored, propagated along the period. With $x_0 = 0.7$ V, $\phi_T(x_0) = 0.65$ V and $\phi' = 0.15$, the residual is -0.05, the derivative -0.85, and $x_1 = 0.7 - (-0.05)/(-0.85) = 0.641$ V. Each Newton iteration costs one period of transient plus the sensitivity, and a few iterations usually suffice.

## Why brute-force transient is slow

In plain transient the initial-condition error decays by the factor $a$ per period. A circuit that keeps 99.5% of its error each period, typical of a high-Q resonator, needs $N$ periods with $0.995^N \le \varepsilon/e_0$. For $e_0 = 0.5$ V and $\varepsilon = 1$ mV, $N \ge \ln(0.002)/\ln(0.995) = 1240$ periods. Shooting reaches the same state in a handful of Newton iterations, each one period long, because Newton removes the slow mode directly instead of waiting for it to decay. For large circuits the dense sensitivity matrix is avoided with matrix-free Krylov methods, which only need its products with vectors.

## Going deeper

Aprille and Trick, "Steady-state analysis of nonlinear circuits with periodic inputs" (Proceedings of the IEEE, 1972); Kundert, White and Sangiovanni-Vincentelli, "Steady-State Methods for Simulating Analog and Microwave Circuits" (Kluwer, 1990).
