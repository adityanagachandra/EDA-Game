# The body effect

When the source of a MOSFET sits above its body, the threshold rises. This body effect matters in every series stack of transistors, in pass gates and in body-biasing schemes that trade speed for leakage. Timing characterization captures it implicitly through stacked-cell simulations, and compact models carry it as an explicit term, so a tool designer should know its size and shape.

## Threshold shift with source-body bias

With a source-to-body reverse bias $V_{SB}$, the surface must bend further, to $2\phi_F + V_{SB}$, before the channel inverts relative to the source. The depletion region deepens and holds more charge, so
$$V_t = V_{t0} + \gamma\left(\sqrt{2\phi_F + V_{SB}} - \sqrt{2\phi_F}\right),$$
with body-effect coefficient $\gamma$ in V$^{1/2}$.

Example: $\gamma = 0.4$ V$^{1/2}$, $2\phi_F = 0.8$ V, $V_{SB} = 0.5$ V. $\Delta V_t = 0.4 \times (\sqrt{1.3} - \sqrt{0.8}) = 0.4 \times (1.1402 - 0.8944) = 98.3$ mV.

The shift is sublinear in $V_{SB}$: the first few hundred millivolts of body bias move the threshold most. Forward body bias (negative $V_{SB}$) lowers $V_t$, which adaptive body-bias schemes use to speed up slow dies.

## Body effect in transistor stacks

In a two-high NMOS stack, such as the pull-down of a NAND2, the lower transistor's source is at ground but the upper transistor's source is the intermediate node. When that node rises during a transition, the upper device sees $V_{SB} > 0$ and a higher threshold, while its gate drive $V_{GS}$ also shrinks. Both effects weaken it.

Example: $V_{t0} = 0.35$ V, $\gamma = 0.3$ V$^{1/2}$, $2\phi_F = 0.85$ V, intermediate node at 0.4 V. $V_t = 0.35 + 0.3 \times (\sqrt{1.25} - \sqrt{0.85}) = 0.35 + 0.0588 = 0.409$ V.

This is one reason stacked gates are slower than an inverter of the same width and why cells size stacked transistors up. It also helps leakage: the stack effect, where two off devices in series leak far less than one, relies partly on this raised threshold. FinFET and FD-SOI devices have a much weaker body effect, because their thin bodies have little depletion charge for the substrate to modulate.

## The body-effect coefficient

The coefficient follows from the depletion charge term of the threshold equation:
$$\gamma = \frac{\sqrt{2q\varepsilon_{Si}N_A}}{C_{ox}}.$$

Example: $N_A = 10^{17}$ cm$^{-3}$ and $t_{ox} = 3$ nm. $\sqrt{2 \times 1.602 \times 10^{-19} \times 1.036 \times 10^{-12} \times 10^{17}} = 1.82 \times 10^{-7}$, and $C_{ox} = 1.15 \times 10^{-6}$ F/cm$^2$, so $\gamma = 0.158$ V$^{1/2}$.

Heavier doping and thicker oxide strengthen the body effect; thin oxides weaken it. This explains a trend a tool designer will see in models: older bulk processes had $\gamma$ around 0.4 to 0.6 V$^{1/2}$, while thin-oxide bulk and especially thin-body devices have much smaller values. Body biasing as a post-silicon tuning knob is therefore most effective in FD-SOI, where the back gate couples through the thin buried oxide instead of through depletion charge.

## Going deeper

Tsividis and McAndrew, "Operation and Modeling of the MOS Transistor" (3rd ed.), chapter 4; Rabaey, Chandrakasan and Nikolić, "Digital Integrated Circuits" (2nd ed.), chapter 3.
