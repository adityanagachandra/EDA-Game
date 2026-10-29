# Bounce diagrams and multiple reflections

When both ends of a line are mismatched, a step bounces back and forth, and the voltage at each end climbs in a staircase toward its final value. A bounce (lattice) diagram tracks each wave in time. It is the mental model behind ringing, settling-time and overshoot checks, and the exact solution a time-domain transmission-line model must reproduce, which makes it a natural golden test.

## Launching the first wave

At $t = 0$ the driver sees the line's characteristic impedance, not the load, so it launches
$$V_i = V_s\frac{Z_0}{R_s + Z_0}.$$
That wave travels to the load in one-way delay $T_d$. Two reflection coefficients then govern every later event: $\Gamma_L = (R_L - Z_0)/(R_L + Z_0)$ at the load and $\Gamma_S = (R_s - Z_0)/(R_s + Z_0)$ at the source.

Example: a 1 V step with $R_s = 25$ ohm on a 50 ohm line with $T_d = 500$ ps into $R_L = 100$ ohm. The launched wave is $50/75 = 0.667$ V; $\Gamma_L = 1/3$ and $\Gamma_S = -1/3$. Until $t = 2T_d$ the source end sits at 0.667 V regardless of the load. A driver that measures its own output during that window learns nothing about the far end.

## Tracking successive reflections

Each wave reaching the load reflects with $\Gamma_L$; each reflection reaching the source reflects with $\Gamma_S$. Arrivals at the load happen at $T_d, 3T_d, 5T_d, \ldots$, and the load voltage after the $k$-th arrival is
$$V_L^{(k)} = V_i(1 + \Gamma_L)\sum_{j=0}^{k-1}(\Gamma_L\Gamma_S)^j.$$
As $k \to \infty$ the geometric series sums to the DC divider value $V_s R_L/(R_s + R_L)$.

Example, continuing: $\Gamma_L\Gamma_S = -1/9$. After the first arrival the load is at $0.667 \cdot 1.333 = 0.889$ V; after the second, at $3T_d$, $0.889(1 - 1/9) = 0.790$ V; after the third, 0.801 V; the final value is $100/125 = 0.8$ V. A negative product $\Gamma_L\Gamma_S$ gives ringing around the final value; a positive one gives a monotonic staircase.

## The voltage at the source end

The source end changes at $0, 2T_d, 4T_d, \ldots$. At $2T_d$ the first reflection $\Gamma_L V_i$ arrives and is itself reflected with $\Gamma_S$, so just after $2T_d$
$$V_S = V_i(1 + \Gamma_L + \Gamma_L\Gamma_S).$$

Example, continuing: $V_S = 0.667(1 + 0.333 - 0.111) = 0.815$ V. With a matched source ($R_s = Z_0$, $\Gamma_S = 0$) the source end steps once to its final value at $2T_d$ and nothing more happens. With a series-terminated driver into an open load, the source end sits at half the swing for $2T_d$, a plateau that is harmless for a point-to-point net but corrupts any receiver placed in the middle of the line. That observation decides where series termination is allowed.

## Going deeper

Johnson and Graham, "High-Speed Digital Design: A Handbook of Black Magic", chapter 4; Hall, Hall and McCall, "High-Speed Digital System Design", chapter 2.
