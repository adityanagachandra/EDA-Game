# Capacitor companion models

Transient analysis turns every capacitor into a resistor and a current source for the duration of one time step. This companion model comes straight from the integration formula: discretize $i = C\, dv/dt$, solve for the current at the new time point, and read off a conductance and a history term. The stamp is then identical to a DC resistor plus current source, so the same Newton and LU code serves every time step.

## Backward Euler capacitor

Backward Euler approximates the derivative at the new point by $dv/dt \approx (v_{n+1} - v_n)/h$. Then
$$i_{n+1} = \frac{C}{h} v_{n+1} - \frac{C}{h} v_n = G_{eq} v_{n+1} - I_{eq},$$
with $G_{eq} = C/h$ and history source $I_{eq} = (C/h) v_n$. The capacitor becomes a conductance $C/h$ in parallel with a current source that remembers the previous voltage. Units pair neatly: picofarads over nanoseconds is millisiemens. A 2 pF capacitor at $h = 0.1$ ns gives $G_{eq} = 20$ mS, and if it held 1.2 V, $I_{eq} = 24$ mA. Small steps make $G_{eq}$ large, which is why short steps behave like holding the capacitor voltage fixed.

## Trapezoidal capacitor

The trapezoidal rule averages the derivative at both ends: $v_{n+1} = v_n + \tfrac{h}{2C}(i_n + i_{n+1})$. Solving for the new current,
$$i_{n+1} = \frac{2C}{h} v_{n+1} - \left(\frac{2C}{h} v_n + i_n\right),$$
so $G_{eq} = 2C/h$ and $I_{eq} = (2C/h) v_n + i_n$. The history now needs the previous current as well as the voltage, which the simulator must store per capacitor. For 1 pF at $h = 0.5$ ns, $G_{eq} = 4$ mS; with $v_n = 0.8$ V and $i_n = -0.3$ mA, $I_{eq} = 3.2 - 0.3 = 2.9$ mA. Trapezoidal is second-order accurate, compared with first order for backward Euler.

## Recovering the capacitor current

After the step's Newton solve gives $v_{n+1}$, the trapezoidal companion needs $i_{n+1}$ for the next step. It follows from the same formula: $i_{n+1} = (2C/h)(v_{n+1} - v_n) - i_n$. For 5 pF at $h = 0.2$ ns, $v_n = 1.1$ V, $i_n = -1.5$ mA and $v_{n+1} = 1.6$ V: $50 \times 0.5 + 1.5 = 26.5$ mA. Note the subtraction of $i_n$: an error in the stored current propagates to every later step and, because the factor multiplying it is -1, it alternates in sign rather than decaying. This is the root of trapezoidal ringing, and it is why simulators reset the history at breakpoints.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 5; Pillage, Rohrer and Visweswariah, "Electronic Circuit and System Simulation Methods" (McGraw-Hill, 1995), chapter on transient analysis.
