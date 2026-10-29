# Series and parallel termination

Termination is how designers stop reflections: make the impedance seen by the wave equal $Z_0$ at the source, at the load, or both. Each scheme trades signal quality against power, component count and topology. A layout or signal-integrity tool that recommends terminations must know these trade-offs, because the right answer depends on driver strength, receiver count and the power budget.

## Series termination at the source

A series resistor $R_{ser}$ at the driver makes $R_{drv} + R_{ser} = Z_0$, so $\Gamma_S = 0$. The driver launches half the swing; at a high-impedance receiver $\Gamma_L = +1$, so the wave doubles to the full swing; the reflection returns and is absorbed at the matched source. The net settles after one round trip.

Example: a driver of 18 ohm output resistance on a 50 ohm line needs $R_{ser} = 32$ ohm. Series termination draws no DC power and needs one resistor per net. Its limits: it suits point-to-point nets only, since receivers along the line see the half-swing plateau for up to $2T_d$, and the driver's output resistance varies with process, voltage, temperature and edge direction, so the match is approximate. Many I/O drivers implement on-die calibrated series termination for this reason.

## Parallel termination at the load

A resistor $R_t = Z_0$ from the receiver to ground (or to a termination voltage) makes $\Gamma_L = 0$: the first wave is absorbed and the receiver sees a clean edge, wherever the receivers are along the line, as long as they are lightly loaded. The cost is DC current whenever the line is held away from the termination voltage.

Example: a 50 ohm line terminated by 50 ohm to ground and held at 1 V dissipates $1^2/50 = 20$ mW in the terminator while high; at 1.8 V it would be 64.8 mW. The driver must also source that current and is loaded by $Z_0$, so its swing falls to the divider $Z_0/(R_{drv} + Z_0)$ of the supply. Across a wide bus that power is significant, which is why low-swing signaling and terminating to a mid-rail voltage are common.

## Thevenin termination

A Thevenin termination places $R_1$ from the receiver to the supply $V_{dd}$ and $R_2$ to ground. To the line it looks like $R_1 \parallel R_2$ in series with a voltage $V_t = V_{dd}R_2/(R_1 + R_2)$. Choosing $R_1 \parallel R_2 = Z_0$ and a desired $V_t$ gives
$$R_1 = \frac{Z_0V_{dd}}{V_t}, \qquad R_2 = \frac{Z_0V_{dd}}{V_{dd} - V_t}.$$

Example: $Z_0 = 50$ ohm, $V_{dd} = 3.3$ V and $V_t = 1.0$ V give $R_1 = 165$ ohm and $R_2 = 71.7$ ohm; check: $165 \parallel 71.7 = 50$ ohm. Biasing the line at mid-swing halves the worst-case drive current compared with terminating to ground, but the divider draws standing current from the supply even when the line is idle. An optimizer choosing among termination schemes must weigh that static current against signal quality.

## Going deeper

Johnson and Graham, "High-Speed Digital Design: A Handbook of Black Magic", chapter 6; Bogatin, "Signal and Power Integrity, Simplified" (3rd ed.), chapter 8.
