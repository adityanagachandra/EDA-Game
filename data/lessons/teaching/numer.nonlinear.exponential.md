# Newton on exponential device laws

Semiconductor devices obey exponential laws: diode and bipolar currents grow like $e^{V/V_T}$, and MOSFETs do in subthreshold. Exponentials make Newton's method both very fast near the answer and very fragile far from it. A tool designer must understand exactly how the iteration behaves on a diode to design the limiting and damping that keep simulators from overflowing or stalling.

## The diode-resistor equation

A source $V_s$ drives a resistor $R$ into a diode to ground. KCL at the diode node gives
$$f(V) = \frac{V - V_s}{R} + I_s\left(e^{V/V_T} - 1\right) = 0,$$
with saturation current $I_s$ (around $10^{-14}$ A for a small junction) and thermal voltage $V_T = kT/q$, about 25.85 mV at 300 K (problems here round it to 25 mV).

Take $V_s = 5$ V, $R = 1$ kΩ, $I_s = 10^{-14}$ A, $V_T = 25$ mV, and guess $V_0 = 0.65$ V. The resistor term is $(0.65 - 5)/1 = -4.35$ mA. The diode term is $10^{-14}e^{26}$ A $= 1.957 \times 10^{-3}$ A, about 1.957 mA. The residual is $f(V_0) = -2.393$ mA: the resistor supplies more current than the diode takes, so the solution lies above 0.65 V.

## One Newton update from an initial guess

The derivative is $f'(V) = 1/R + (I_s/V_T)e^{V/V_T}$, the resistor conductance plus the diode's small-signal conductance $g_d \approx I_D/V_T$. At $V_0 = 0.65$ V: $g_d = 1.957\text{ mA}/25\text{ mV} = 78.3$ mS, so $f' = 1 + 78.3 = 79.3$ mS. Then
$$V_1 = V_0 - \frac{f(V_0)}{f'(V_0)} = 0.65 + \frac{2.393}{79.3} = 0.680\ \text{V}.$$
Geometrically, Newton replaces the diode by its companion model at $V_0$: a conductance $g_d$ in parallel with a current source, and solves the resulting linear circuit. Every SPICE iteration does exactly this for every nonlinear device.

## Overshoot and exponential sensitivity

Current multiplies by $e^{\Delta V/V_T}$ for a voltage change $\Delta V$: 60 mV gives $e^{2.4} = 11$, about a factor of ten, and 0.5 V gives $e^{20} \approx 5 \times 10^8$. From a guess below the solution, where $g_d$ is tiny, the linearized diode looks almost open, and Newton jumps toward $V_s$, perhaps to several volts, where $e^{V/V_T}$ overflows or yields absurd currents.

From a guess far above the solution the opposite happens: the tangent is extremely steep, and each Newton step falls by only about $V_T$, so convergence crawls. Both failures are why SPICE limits junction voltage changes per iteration (pnjlim uses a logarithmic step for large changes), why global damping exists, and why device models clamp their exponentials beyond a critical voltage.

## Going deeper

Nagel, "SPICE2: A computer program to simulate semiconductor circuits", UC Berkeley ERL Memo M520 (1975); Najm, "Circuit Simulation" (Wiley, 2010), chapter 3.
