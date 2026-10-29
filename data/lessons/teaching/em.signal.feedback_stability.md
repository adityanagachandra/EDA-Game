# Feedback, loop gain and phase margin

Feedback makes gain precise and insensitive to device variation, but it can also make a system oscillate. The loop gain decides both. Regulators in power delivery, phase-locked loops, amplifiers and even convergence loops inside optimization engines are feedback systems, and phase margin is the standard measure of how far each is from instability.

## Closed-loop gain

An amplifier with forward gain $A$ and feedback factor $\beta$ in a negative-feedback loop has closed-loop gain
$$A_{cl} = \frac{A}{1 + A\beta}.$$
The product $T = A\beta$ is the loop gain. When $T \gg 1$, $A_{cl} \approx 1/\beta$: the result depends on the feedback network, usually precise passives, not on the amplifier.

Example: $A = 1000$ and $\beta = 0.1$ give $A_{cl} = 1000/101 = 9.901$, within 1% of the ideal 10. If $A$ drops by half to 500, $A_{cl}$ becomes $500/51 = 9.804$, a change of only 1%: feedback divided the sensitivity by $1 + T$. That desensitization is the reason feedback is everywhere, and the cost is that the loop must be stable.

## Crossover of a one-pole loop

The crossover frequency $f_c$ is where $|T(f)| = 1$. For a single-pole loop $T(f) = T_0/(1 + jf/f_p)$,
$$f_c = f_p\sqrt{T_0^2 - 1} \approx T_0 f_p.$$
Above its pole the loop gain falls at 20 dB per decade, so the gain-bandwidth product $T_0 f_p$ is constant and equals the crossover.

Example: $T_0 = 200$ and $f_p = 5$ kHz give $f_c \approx 1$ MHz. A one-pole loop's phase never passes $-90^\circ$, so it always has at least $90^\circ$ of phase margin and cannot oscillate. Crossover sets the loop's speed: closed-loop bandwidth is about $f_c$, and a regulator with a 1 MHz crossover cannot respond to load steps faster than about a microsecond, which is why on-die decap must cover the faster events.

## Phase margin of a two-pole loop

Phase margin is $180^\circ$ plus the loop's phase at crossover: $PM = 180^\circ + \angle T(f_c)$. For $T(f) = T_0/((1 + jf/f_1)(1 + jf/f_2))$,
$$PM = 180^\circ - \arctan(f_c/f_1) - \arctan(f_c/f_2),$$
with $f_c$ found by solving $|T(f_c)| = 1$, usually numerically.

Example: $T_0 = 1000$, $f_1 = 1$ kHz and $f_2 = 1$ MHz. The one-pole estimate puts $f_c$ at 1 MHz, right at the second pole, where the extra pole cuts $|T|$ by $\sqrt 2$; solving exactly gives $f_c = 786$ kHz and $PM = 180 - 89.9 - 38.2 = 51.9^\circ$. Below about $45^\circ$ the closed-loop step response rings noticeably; at zero the loop oscillates. Compensation moves poles to buy margin, trading away bandwidth. The same check applies to any loop an engine closes, including control loops in adaptive tools.

## Going deeper

Franklin, Powell and Emami-Naeini, "Feedback Control of Dynamic Systems", chapters 4 and 6; Gray, Hurst, Lewis and Meyer, "Analysis and Design of Analog Integrated Circuits", chapter 9.
