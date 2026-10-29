# Inductor companion models

Inductors get companion models too, by discretizing $v = L\, di/dt$ instead of $i = C\, dv/dt$. The result is again a conductance with a history source, now expressing the new inductor current in terms of the new voltage. Simulators that keep inductor currents as branch unknowns stamp the equivalent branch equation instead, but the coefficients are the same, and computing them by hand fixes the duality between the two reactive elements.

## Backward Euler inductor

Backward Euler applied to $di/dt = v/L$ gives $i_{n+1} = i_n + (h/L) v_{n+1}$. Read as a companion, the inductor is a conductance $G_{eq} = h/L$ in parallel with a current source $I_{eq} = i_n$ that carries the previous current, with $i_{n+1} = G_{eq} v_{n+1} + I_{eq}$. Units: nanoseconds over microhenries is millisiemens. A 1 uH inductor at $h = 0.5$ ns has $G_{eq} = 0.5$ mS. If it carries 2 mA and the step solve finds 0.4 V across it, the new current is $2 + 0.5 \times 0.4 = 2.2$ mA. Small steps make $G_{eq}$ small: over a short step an inductor holds its current, as a capacitor holds its voltage.

## Trapezoidal inductor conductance

The trapezoidal rule averages the voltage over the step: $i_{n+1} = i_n + \tfrac{h}{2L}(v_n + v_{n+1})$. The coefficient of $v_{n+1}$ is the companion conductance,
$$G_{eq} = \frac{h}{2L},$$
half the backward Euler value for the same step. A 5 uH inductor at $h = 0.2$ ns gives $G_{eq} = 0.02$ mS. In an MNA implementation that keeps the inductor current $i$ as an unknown, the branch row becomes $v_+ - v_- - \tfrac{2L}{h} i_{n+1} = -\tfrac{2L}{h} i_n - v_n$, which carries the same information: the diagonal of the branch row is now $-2L/h$ instead of the zero it has at DC.

## The trapezoidal history current

The remaining terms of the trapezoidal update form the history source:
$$I_{eq} = i_n + \frac{h}{2L} v_n.$$
It needs both the previous current and the previous voltage. For 0.2 uH at $h = 0.2$ ns, with 2.5 mA and 1.9 V at $t_n$: $G_{eq} = 0.5$ mS and $I_{eq} = 2.5 + 0.5 \times 1.9 = 3.45$ mA. As with the capacitor, the trapezoidal history carries the previous step's derivative information, so a sudden voltage step across an inductor makes the computed current oscillate from step to step instead of settling. Backward Euler's history, the current alone, has no such memory, which is why it damps such transients.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 5; Vlach and Singhal, "Computer Methods for Circuit Analysis and Design" (2nd ed., 1994).
