# Aging and lifetime guard bands

Aging is the slow drift of transistor parameters over a product's life, chiefly threshold voltage shifts from bias temperature instability and hot-carrier injection. Their effect on circuits is slower gates, so a design must carry a guard band in timing so that it still meets frequency at end of life. Package temperature drives the rate, which connects thermal design to timing signoff.

## Power-law threshold shift

Bias temperature instability shifts threshold voltage roughly as a power of stress time, $\Delta V_{th} = A t^n$, with $n$ often between about 0.1 and 0.25. So the shift at a later time is $\Delta V_{th}(t_2) = \Delta V_{th}(t_1)(t_2/t_1)^n$. A 15 mV shift after 1000 hours of stress, extrapolated to 10 years (87600 hours) with $n = 0.2$, gives $15 \times 87.6^{0.2} = 15 \times 2.45 = 36.7$ mV. The small exponent means most of the shift happens early, which is why burn-in captures a large part of it. BTI also partly recovers when stress is removed, so duty cycle and the measurement delay affect the observed value.

## Time to a shift limit

Inverting the power law gives the time to reach an allowed shift: $t = t_1(\Delta V_{lim}/\Delta V_1)^{1/n}$. For a 50 mV limit, the example's 15 mV at 1000 hours and $n = 0.2$ give $1000 \times (50/15)^5 = 411500$ hours, about 47 years. The exponent $1/n$ is large, so small changes in the measured shift or the exponent change the answer a lot, the same fragility as any long extrapolation. Temperature accelerates BTI through an Arrhenius-like factor in $A$, so a package that runs 15 K hotter reaches the limit much sooner. That is the coupling to thermal design.

## Delay guard bands

Threshold shift slows gates. In the alpha-power model, delay scales as $1/(V_{DD} - V_{th})^{\alpha}$, with $\alpha$ between 1 and 2 for short-channel devices. The relative delay after aging is $((V_{DD} - V_{th0})/(V_{DD} - V_{th0} - \Delta V_{th}))^{\alpha}$. At 0.75 V supply, 0.3 V threshold, 30 mV shift and $\alpha = 1.3$, this is $(0.45/0.42)^{1.3} = 1.094$, a 9.4% slowdown. Signoff adds this as an aging guard band, either a flat derate or aged cell libraries characterized at end-of-life conditions. A guard band that is too small fails in the field, and one that is too large wastes frequency on every part.

## Going deeper

Schroder and Babcock, "Negative bias temperature instability", Journal of Applied Physics (2003); Sakurai and Newton, "Alpha-power law MOSFET model", IEEE Journal of Solid-State Circuits (1990).
