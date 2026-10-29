# Poles, zeros and stability

The poles of a transfer function are the natural frequencies of the system; the zeros are the frequencies it blocks. Pole locations decide stability, speed, damping and ringing. Model order reduction keeps the dominant poles, stability checks look for poles in the right half-plane, and passivity checks of extracted models start from the same picture.

## Poles as roots of the denominator

Write $H(s) = N(s)/D(s)$. The poles are the roots of $D(s)$ and the zeros the roots of $N(s)$. A real pole at $s = -p$ contributes a mode $e^{-pt}$ to the response, so its time constant is $1/p$. The pole closest to the origin decays slowest and dominates the response: it is the dominant pole.

Example: $D(s) = s^2 + 7s + 10$ with $s$ in units of $10^9$ per second factors as $(s + 2)(s + 5)$, poles at -2 and -5. The dominant pole -2 has a time constant of 0.5 ns. For $D(s) = s^2 + 4s + 13$ the roots are complex, $-2 \pm 3j$: the real part -2 sets the decay rate of the envelope and the imaginary part sets the ringing frequency, 3 Grad/s.

## Stability from pole locations

A linear system is stable when every pole has a negative real part: every natural mode then decays. A pole with positive real part grows without bound; a pole on the imaginary axis oscillates forever. For a quadratic factor $s^2 + bs + c$, if $c < 0$ the roots have opposite signs, one in the right half-plane; if $c > 0$ both roots lie in the half-plane where the sign of their real part equals the sign of $-b$.

Example: $(s + 3)(s^2 - 2s + 5)$ has one pole at -3 and two at $1 \pm 2j$: two right-half-plane poles, unstable. A passive RLC network can never have right-half-plane poles, so if a reduced-order interconnect model shows one, the reduction itself broke something. That check catches a real failure mode of moment-matching methods.

## Damping from pole positions

A complex pole pair $s = -\sigma \pm j\omega_d$ corresponds to a second-order factor $s^2 + 2\zeta\omega_n s + \omega_n^2$. The undamped natural frequency is the poles' distance from the origin, $\omega_n = \sqrt{\sigma^2 + \omega_d^2}$, and the damping ratio is $\zeta = \sigma/\omega_n$, the cosine of the angle the pole makes with the negative real axis.

Example: poles at $-3 \pm 4j$ (units of $10^9$ per second) give $\omega_n = 5$ and $\zeta = 0.6$, so the step overshoot is $e^{-0.6\pi/0.8} = 9.5$%. Poles at $-1 \pm 4j$ give $\omega_n = 4.12$ and $\zeta = 0.24$, much more ringing. Moving poles toward the imaginary axis lowers $\zeta$. Engineers read a pole plot geometrically: radius is speed, angle is damping.

## Going deeper

Franklin, Powell and Emami-Naeini, "Feedback Control of Dynamic Systems", chapter 3; Oppenheim and Willsky, "Signals and Systems" (2nd ed.), chapter 9.
