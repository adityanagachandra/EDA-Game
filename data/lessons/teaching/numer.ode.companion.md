# Companion models of integration steps

A simulator does not integrate capacitor and inductor equations separately from the circuit. It discretizes each dynamic element with the integration formula and replaces it by a resistive companion: a conductance plus a current source. Each timestep then becomes a DC solve of a resistive network with the same pattern, which is how integration and sparse linear algebra fit together.

## A capacitor becomes a conductance and a current source

A capacitor obeys $i = C\,dv/dt$. Backward Euler replaces the derivative by $(v_{n+1} - v_n)/h$:
$$i_{n+1} = \frac{C}{h}v_{n+1} - \frac{C}{h}v_n.$$
That is a conductance $G_{eq} = C/h$ in parallel with a history current source of value $(C/h)\,v_n$ that depends only on the previous step. With $C = 100$ fF and $h = 20$ ps, $G_{eq} = 5$ mS (one femtofarad per picosecond is one millisiemens). If $v_n = 0.8$ V, the history source is 4 mA.

The companion is stamped into the MNA matrix like any resistor and current source. Small steps make $C/h$ large, which strengthens the diagonal and makes each timestep's linear system better conditioned than the DC one.

## Trapezoidal companion with the history current

The trapezoidal rule averages the capacitor current: $(i_{n+1} + i_n)/2 = C(v_{n+1} - v_n)/h$, so
$$i_{n+1} = \frac{2C}{h}v_{n+1} - \Big[\frac{2C}{h}v_n + i_n\Big].$$
The conductance doubles to $2C/h$, and the history source needs the previous capacitor current $i_n$ as well as $v_n$. With $C = 100$ fF, $h = 20$ ps, $v_n = 0.8$ V and $i_n = -0.5$ mA: $G_{eq} = 10$ mS and the history current is $8 - 0.5 = 7.5$ mA.

The simulator must store $i_n$ for every capacitor and update it after each accepted step. At a breakpoint, an inconsistent $i_n$ can seed trapezoidal ringing, one reason simulators restart with a backward Euler step there.

## Inductor companions and the solve per step

An inductor obeys $v = L\,di/dt$. Backward Euler gives $i_{n+1} = i_n + (h/L)\,v_{n+1}$: a conductance $h/L$ in parallel with a current source $i_n$. With $L = 2$ nH, $h = 10$ ps, $i_n = 1$ mA and a solved voltage of 0.4 V, $h/L = 5$ mS and $i_{n+1} = 1 + 2 = 3$ mA.

Putting it together: a 1 V source drives 1 kΩ into a node with a 100 fF capacitor at $v_n = 0.2$ V, $h = 20$ ps. The companion is 5 mS with a 1 mA history source. KCL gives $(1 - v)/1 + 1 = 5v$ (in mA), so $v = 2/6 = 0.333$ V. Each timestep is exactly this: stamp companions, solve the resistive network (with Newton if devices are nonlinear), update histories.

## Going deeper

Vlach and Singhal, "Computer Methods for Circuit Analysis and Design" (2nd ed., 1994), chapter on numerical integration; Najm, "Circuit Simulation" (Wiley, 2010), chapter 4.
