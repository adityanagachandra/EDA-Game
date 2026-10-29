# BDF and Gear integration

Backward differentiation formulas (BDF), introduced to circuit simulation by Gear, approximate the derivative at the new time point from several past points. They are the standard alternative to trapezoidal in SPICE (method=gear) because they damp stiff modes, at the cost of needing more history and, above second order, losing A-stability. The second-order formula, BDF2, has a companion model as simple as backward Euler's.

## The BDF2 derivative formula

BDF2 differentiates the quadratic through $x_{n-1}$, $x_n$ and $x_{n+1}$. At constant step $h$,
$$x'_{n+1} \approx \frac{3x_{n+1} - 4x_n + x_{n-1}}{2h}.$$
Applied to a capacitor, $i_{n+1} = C x'_{n+1}$, so the coefficient of $v_{n+1}$ is the companion conductance $G_{eq} = 3C/(2h)$, one and a half times the backward-Euler value. For 5 pF at $h = 0.1$ ns, $G_{eq} = 75$ mS. BDF of order 1 is backward Euler; orders 1 and 2 are A-stable and L-stable, while orders 3 to 6 are only stable in a wedge of the left half plane, which is usually acceptable for circuits.

## The BDF2 history current

The remaining terms form the history: $i_{n+1} = G_{eq} v_{n+1} - I_{eq}$ with
$$I_{eq} = \frac{C}{2h}(4v_n - v_{n-1}).$$
For 4 pF at $h = 0.2$ ns with $v_n = 0.8$ V and $v_{n-1} = 0.95$ V, $C/(2h) = 10$ mS and $4v_n - v_{n-1} = 2.25$ V, so $I_{eq} = 22.5$ mA. The history needs two past voltages but no past current. Because it needs two points, BDF2 cannot start itself: simulators take a backward-Euler first step. With variable steps the coefficients change with the step ratio, and simulators store the past points and recompute them.

## One BDF2 step on an RC node

For a source $V_s$ through $R$ into $C$, KCL gives $(1/R + 3C/(2h)) v_{n+1} = V_s/R + I_{eq}$. With $V_s = 3.5$ V, $R = 0.5$ k$\Omega$, $C = 0.5$ pF, $h = 0.1$ ns, $v_n = 1.8$ V and $v_{n-1} = 1.65$ V: $1/R = 2$ mS, $G_{eq} = 7.5$ mS, $I_{eq} = 2.5 \times 5.55 = 13.875$ mA, so $v_{n+1} = (7 + 13.875)/9.5 = 2.197$ V. For fast modes BDF2's amplification factor tends to 0, so there is no ringing; for genuine oscillations it adds artificial damping, so a high-Q oscillator simulated with Gear can decay that would not decay with trapezoidal.

## Going deeper

Gear, "Numerical Initial Value Problems in Ordinary Differential Equations" (Prentice Hall, 1971); Brayton, Gustavson and Hachtel, "A new efficient algorithm for solving differential-algebraic systems using implicit backward differentiation formulas" (Proceedings of the IEEE, 1972).
