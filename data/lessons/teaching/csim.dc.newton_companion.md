# Newton-Raphson with companion models

Newton-Raphson on circuit equations has an elegant implementation: linearize every device at the current iterate, stamp the companion models, solve the resulting linear circuit, and repeat. Each linear solve is exactly one Newton step on the KCL residual. This is how SPICE computes every DC operating point, and the same loop runs inside each transient time step. Working one step by hand demystifies the iteration logs a simulator prints.

## Stamping the linearized diode

Take a source $V_s$ through resistor $R$ into node d, with a diode from d to ground. At iterate $v_k$ the diode becomes $g_d$ in parallel with $I_{eq}$ from d to ground. KCL at d, currents leaving: $(v - V_s)/R + g_d v + I_{eq} = 0$. The single matrix entry is $1/R + g_d$, and the right-hand side is $V_s/R - I_{eq}$ (the source and resistor form a Norton equivalent, and $I_{eq}$ flowing out of d moves to the right with a minus sign). With $R = 1$ k$\Omega$ and $g_d = 23.58$ mS, the entry is 24.58 mS. In a larger circuit each device adds its companion to whatever else is stamped.

## One Newton iterate

Solving the linear circuit gives the next iterate
$$v_{k+1} = \frac{V_s/R - I_{eq}}{1/R + g_d}.$$
For $V_s = 2$ V, $R = 10$ k$\Omega$, $I_S = 1$ fA and $v_k = 0.58$ V: $I_d = 0.0119$ mA, $g_d = 0.476$ mS, $I_{eq} = 0.0119 - 0.476 \cdot 0.58 = -0.264$ mA. Then $v_{k+1} = (0.2 + 0.264)/(0.1 + 0.476) = 0.806$ V. The step overshoots: the true solution is near 0.64 V, but the tangent at 0.58 V underestimates the diode current, so the linear circuit puts too much voltage on the diode. The next iterate's diode current would be about $e^{9}$ times larger, which is why simulators limit such steps.

## The KCL residual

Newton solves $f(v) = 0$, where $f$ is the KCL residual: the net current leaving each node through all elements, evaluated with the true nonlinear devices. For the diode circuit, $f(v_k) = (v_k - V_s)/R + I_d(v_k)$. The Newton step is $\Delta v = -f(v_k)/f'(v_k)$ with $f' = 1/R + g_d$, which is algebraically identical to solving the companion circuit. At $v_k = 0.65$ V with $V_s = 2$ V, $R = 4$ k$\Omega$ and $I_S = 10$ fA, $f = -0.3375 + 1.96 = 1.62$ mA: too much current leaves, so the voltage must fall. Simulators check both the update size and this residual before declaring convergence.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 4; Pillage, Rohrer and Visweswariah, "Electronic Circuit and System Simulation Methods" (McGraw-Hill, 1995), chapter on nonlinear DC analysis.
