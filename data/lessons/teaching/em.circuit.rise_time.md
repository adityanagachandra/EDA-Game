# Rise time and 50 percent delay of an RC stage

Timing libraries characterize every cell by two numbers per transition: the delay to the 50% point and the output transition time, usually measured between 10% and 90% (or 20% and 80%). For a single RC stage both are fixed multiples of $\tau$, and the rise time also sets the stage's bandwidth. These relations let a tool move between time-domain and frequency-domain views of the same edge.

## The 50 percent delay

For a step into a single RC stage, the output reaches 50% when $1 - e^{-t/\tau} = 0.5$, so $t_{50} = \tau\ln 2 \approx 0.693\tau$. This is the delay a timer would report for the stage if its input were an ideal step and the 50% points were used for delay measurement.

Example: a 4 kohm effective driver into a 25 fF load has $\tau = 100$ ps and $t_{50} = 69.3$ ps. Real inputs have finite slew, so library delays depend on input transition as well as load; the 0.69RC figure is the step-input limit. When a resistance and capacitance are distributed along a wire rather than lumped, the coefficient changes, which is the next concept's subject.

## The 10 to 90 percent rise time

The 10% point is at $\tau\ln(10/9) = 0.105\tau$ and the 90% point at $\tau\ln 10 = 2.303\tau$. Their difference is the 10-90 rise time, $t_r = \tau\ln 9 \approx 2.197\tau$. For a 20-80 measurement it is $\tau\ln 4 = 1.386\tau$.

Example: $R = 1$ kohm and $C = 200$ fF give $\tau = 200$ ps and $t_r = 439$ ps. The slew a timer propagates matters as much as the delay, because a slow input edge increases the next gate's delay and short-circuit current. Different libraries use different thresholds, so a delay calculator must convert slews with factors like $\ln 9/\ln 4 = 1.585$ between 10-90 and 20-80 conventions when it mixes data.

## Rise time and bandwidth

A single-pole low-pass $H(s) = 1/(1 + s\tau)$ has its 3 dB bandwidth at $f_{3dB} = 1/(2\pi\tau)$. Combining with $t_r = \tau\ln 9$ gives
$$t_r f_{3dB} = \frac{\ln 9}{2\pi} \approx 0.35.$$
So a rise time in picoseconds converts to bandwidth in gigahertz by $f_{3dB} \approx 350/t_r$.

Example: a 100 ps rise time corresponds to about 3.5 GHz of bandwidth. The rule ties together what a scope can resolve, what an interconnect model must represent, and how fast a buffer can drive a line. It is exact only for one real pole; systems with several poles or resonances have a different constant. A model fitted only up to 1 GHz cannot be trusted to reproduce a 50 ps edge, which needs about 7 GHz.

## Going deeper

Johnson and Graham, "High-Speed Digital Design: A Handbook of Black Magic", chapter 1; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 4.
