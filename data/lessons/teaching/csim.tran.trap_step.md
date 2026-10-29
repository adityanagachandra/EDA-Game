# One trapezoidal transient step

The trapezoidal rule is SPICE's default integrator because it is second-order accurate and A-stable at the same cost per step as backward Euler. Its companion model needs one extra piece of history, the capacitor current, and that piece must be updated after every step. Working trapezoidal steps by hand shows both the accuracy gain and the bookkeeping that a simulator carries.

## The trapezoidal step equation

With the trapezoidal companion, KCL at an RC node driven by $V_s$ through $R$ is
$$\left(\frac{1}{R} + \frac{2C}{h}\right) v_{n+1} = \frac{V_s}{R} + \frac{2C}{h} v_n + i_n,$$
where $i_n$ is the capacitor current at $t_n$. With a constant source, $i_n = (V_s - v_n)/R$ by KCL. For $V_s = 2$ V, $R = 1$ k$\Omega$, $C = 1$ pF, $h = 1$ ns and $v_n = 0$: $i_n = 2$ mA, and $(1 + 2) v_{n+1} = 2 + 0 + 2$, so $v_{n+1} = 1.333$ V. The exact value at $t = \tau$ is 1.264 V; backward Euler would give 1.0 V. Note that $i_n$ just after the step is the full $V_s/R$.

## Updating the capacitor current

The next step needs $i_{n+1}$, computed from the companion formula $i_{n+1} = (2C/h)(v_{n+1} - v_n) - i_n$. In the example, $2 \times 1.333 - 2 = 0.667$ mA. KCL gives an independent check: the capacitor current must equal the resistor current, $(V_s - v_{n+1})/R = 0.667$ mA. In a linear circuit the two always agree, but in a simulator only the companion formula is used, because it is local to the element. If the stored current were computed carelessly, for example taken as zero after a source discontinuity, the error would persist and alternate in sign over later steps.

## Two trapezoidal steps

Repeating with $v_1 = 1.333$ V and $i_1 = 0.667$ mA gives $3 v_2 = 2 + 2.667 + 0.667$, so $v_2 = 1.778$ V. The exact value at $2\tau$ is 1.729 V. In general the trapezoidal rule multiplies the distance to the final value by $(1 - h/2\tau)/(1 + h/2\tau)$ per step, which matches $e^{-h/\tau}$ to second order in $h$: here $1/3$ against $e^{-1} = 0.368$. The local error is of order $h^3$ and the global error of order $h^2$, so halving the step cuts the error by four. When $h$ exceeds $2\tau$, however, the factor goes negative, which is the subject of trapezoidal ringing.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 5; Pillage, Rohrer and Visweswariah, "Electronic Circuit and System Simulation Methods" (McGraw-Hill, 1995).
