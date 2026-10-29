# Linearized diode companion model

Nonlinear devices enter MNA through linearization. At each Newton iterate, a diode is replaced by its tangent line: a conductance in parallel with a current source. This companion model has the same stamp shape as a resistor plus a current source, so the linear solver never sees a nonlinear element. Every SPICE device model, from the diode to BSIM, ends in the same move: compute currents and their derivatives at the present voltages and stamp them.

## Conductance at the iterate

The ideal diode current is $I(V) = I_S(e^{V/V_T} - 1)$, with $V_T \approx 25$ mV near room temperature. Its derivative at the iterate $V_k$ is the small-signal conductance
$$g_d = \frac{dI}{dV}\Big|_{V_k} = \frac{I_S}{V_T} e^{V_k/V_T}.$$
For $I_S = 10$ fA and $V_k = 0.6$ V, $e^{24} = 2.65 \times 10^{10}$, so $I_d = 0.265$ mA and $g_d = 10^{-14}/0.025 \times 2.65 \times 10^{10} = 10.6$ mS. Note that $g_d \approx (I_d + I_S)/V_T$: the conductance is the current divided by the thermal voltage, which is a quick check. Because both grow by a factor of $e$ for every 25 mV, a 60 mV error in the iterate changes $g_d$ about elevenfold.

## The equivalent current source

The tangent line at $V_k$ is $I \approx I_d(V_k) + g_d(V - V_k) = g_d V + I_{eq}$, with
$$I_{eq} = I_d(V_k) - g_d V_k.$$
The diode is replaced by $g_d$ in parallel with a current source $I_{eq}$ flowing from anode to cathode. In the example, $I_{eq} = 0.265 - 10.6 \times 0.6 = -6.09$ mA: a large negative value, because the tangent of a steep exponential crosses the voltage axis far to the right of the origin. The conductance stamps like a resistor; $I_{eq}$ stamps like a current source, entering the right-hand side with sign set by its direction. The same pattern holds for every two-terminal nonlinearity.

## What the linear model predicts

The companion model is exact only at $V_k$. At $V_k + \Delta V$ it predicts $I_d(V_k) + g_d \Delta V$, while the true current is $I_d(V_k) e^{\Delta V/V_T}$. With $\Delta V = 40$ mV, the tangent predicts $I_d(1 + 1.6) = 2.6 I_d$ and the diode delivers $e^{1.6} I_d = 4.95 I_d$. The exponential outruns its tangent above $V_k$ and the tangent overshoots below it, going negative where the diode cannot. This is why a Newton step built from the companion overshoots when the iterate is low and why simulators limit junction voltages: the linear model is trustworthy only within a few thermal voltages of the point where it was built.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 4; Nagel, "SPICE2" (Berkeley ERL Memo M520, 1975).
