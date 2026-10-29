# RC circuits and time constants

A resistor charging a capacitor is the most important waveform in digital design: every gate driving a wire and every wire charging a gate input is, to first order, an RC stage. The time constant $\tau = RC$ sets how fast it settles. Timers, delay calculators and reduced-order interconnect models all start from this exponential, so its formulas must be second nature to a tool builder.

## The time constant

KCL at the capacitor of a series RC driven by $V_s$ gives the first-order ODE $RC\,dv/dt + v = V_s$. Its natural response decays as $e^{-t/\tau}$ with time constant $\tau = RC$. On chip, kilohms times femtofarads give picoseconds, which makes mental arithmetic easy.

Example: a 2 kohm effective driver resistance into a 30 fF load gives $\tau = 60$ ps. When other resistors connect to the capacitor's node, the right $R$ is the Thevenin resistance seen by the capacitor with the sources zeroed. A source with 4 kohm series resistance driving a node that also has 4 kohm to ground and 50 fF gives $R_{th} = 4 \parallel 4 = 2$ kohm and $\tau = 100$ ps, not 200 ps. Forgetting the parallel path is a classic error in hand estimates.

## Charging and discharging waveforms

For a step from 0 to $V$ at $t = 0$ with the capacitor initially discharged, the solution is $v(t) = V(1 - e^{-t/\tau})$. A capacitor charged to $V$ and discharged through $R$ follows $v(t) = Ve^{-t/\tau}$. More generally $v(t) = v_\infty + (v_0 - v_\infty)e^{-t/\tau}$: final value plus the initial gap decaying exponentially.

Example: $\tau = 50$ ps and a 1 V step. At 50 ps the voltage is $1 - e^{-1} = 0.632$ V; at 100 ps, $1 - e^{-2} = 0.865$ V; at 150 ps, 0.950 V. After five time constants the gap is under 1%. The same formula covers discharge with $v_\infty = 0$: a 1 V node discharging with $\tau = 50$ ps is at 0.368 V after 50 ps. A simulator uses this closed form to check its integrator on a single RC, the simplest golden test of a transient engine.

## Time to cross a threshold

Delay is the time for the output to reach a threshold. Setting $1 - e^{-t/\tau} = f$ for a fraction $f$ of the final value gives $t = \tau \ln\frac{1}{1-f}$. For $f = 0.5$ this is $0.693\tau$; for 63.2% it is exactly $\tau$; for 90% it is $2.303\tau$; for 99% it is $4.605\tau$.

Example: $\tau = 80$ ps and a receiver that switches at 70% of the supply: $t = 80 \ln(1/0.3) = 80 \cdot 1.204 = 96.3$ ps. The threshold choice changes the number a lot, which is why libraries define delay at a stated threshold, usually 50%. The logarithm also shows a trap: chasing the last few percent of settling is expensive, so a timer that waits for full settling would be badly pessimistic.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 4; Boyce and DiPrima, "Elementary Differential Equations", chapter 2 on first-order linear equations.
