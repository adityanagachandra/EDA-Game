# RLC circuits

Add inductance to an RC stage and the step response can ring. Series RLC behavior appears in package and bond-wire models, power-delivery networks, clock lines and long on-chip wires at high frequency. Its two parameters, the natural frequency and the damping ratio, decide whether an edge settles cleanly or overshoots, and they are the first thing a signal-integrity engine reports.

## Natural frequency and damping ratio

A series RLC with output across the capacitor obeys $LC\,v'' + RC\,v' + v = V_s$. The characteristic equation $LCs^2 + RCs + 1 = 0$ has roots set by two parameters: the undamped natural frequency $\omega_0 = 1/\sqrt{LC}$ and the damping ratio
$$\zeta = \frac{R}{2}\sqrt{\frac{C}{L}} = \frac{R}{2\sqrt{L/C}}.$$
The quantity $\sqrt{L/C}$ is the circuit's characteristic resistance.

Example: $L = 4$ nH and $C = 1$ pF give $\sqrt{L/C} = \sqrt{4000} = 63.25$ ohm and $\omega_0 = 1/\sqrt{4 \times 10^{-21}} = 1.58 \times 10^{10}$ rad/s, so $f_0 = 2.52$ GHz. With $R = 20$ ohm, $\zeta = 20/126.5 = 0.158$, strongly underdamped. Nanohenries over picofarads give $\sqrt{1000\,L/C}$ in ohms. Note that damping grows with $R$ in a series circuit; in a parallel RLC the roles flip and a larger $R$ means less damping.

## Critical damping

The roots are real and distinct when $\zeta > 1$ (overdamped: slow, no overshoot), repeated when $\zeta = 1$ (critically damped: fastest without overshoot) and complex when $\zeta < 1$ (underdamped: ringing at $\omega_d = \omega_0\sqrt{1-\zeta^2}$). Setting $\zeta = 1$ gives the critical resistance
$$R_{crit} = 2\sqrt{L/C}.$$

Example: a 2 nH package lead driving 2 pF has $\sqrt{L/C} = 31.6$ ohm, so $R_{crit} = 63.2$ ohm. A driver with less output resistance than that rings on this load. This is the lumped cousin of source termination on a transmission line, where the series resistance is chosen to match the line's characteristic impedance; in both cases the knob is a resistance compared with $\sqrt{L/C}$.

## Overshoot of an underdamped step

For $\zeta < 1$ the capacitor voltage after a unit step overshoots its final value. The peak occurs at $t_p = \pi/\omega_d$ and the fractional overshoot is
$$M_p = e^{-\zeta\pi/\sqrt{1-\zeta^2}}.$$
It depends only on $\zeta$: about 73% at $\zeta = 0.1$, 37% at $\zeta = 0.3$, 16% at $\zeta = 0.5$ and 4.3% at $\zeta = 0.707$.

Example: $L = 10$ nH, $C = 10$ pF and $R = 12.65$ ohm give $\sqrt{L/C} = 31.6$ ohm and $\zeta = 0.2$, so the overshoot is $e^{-0.2\pi/0.98} = 52.7$%. A 1 V edge peaks near 1.53 V, which can stress gate oxides and inject noise. Overshoot and ringing also cost timing, because the signal can cross the receiver threshold more than once. Signal-integrity checkers flag nets by exactly these metrics.

## Going deeper

Agarwal and Lang, "Foundations of Analog and Digital Electronic Circuits", chapter 12; Bakoglu, "Circuits, Interconnections, and Packaging for VLSI", chapter 6.
