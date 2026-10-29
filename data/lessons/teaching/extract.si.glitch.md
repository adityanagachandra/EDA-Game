# Crosstalk glitches

A quiet net is not safe from its neighbours. When an aggressor switches, charge injected through the coupling capacitor bumps the victim's voltage, a glitch. If the glitch is large enough and reaches a sequential element or a dynamic node, it can flip a stored value, a functional failure rather than a timing one. The first estimate of a glitch's size is a capacitive divider, an upper bound that every noise analysis starts from.

## The capacitive divider

Suppose the victim's driver is too weak to respond during a fast aggressor edge, so the victim floats. The aggressor's swing $V_{DD}$ divides across the coupling and the victim's other capacitance: $V_p = V_{DD}\, C_c/(C_c + C_g)$, by charge conservation on the floating node. A victim with 7.9 fF to ground and 1.4 fF of coupling to a 1 V aggressor sees $1.4/9.3 = 150.5$ mV. Real victims are held by a driver that bleeds the injected charge away, so the real peak is lower; the divider is the bound, reached when the aggressor edge is much faster than the victim's holding time constant. Static noise analysis uses it as a cheap filter: nets whose bound is below the noise threshold need no further work.

## Several aggressors

With several aggressors, every coupling is part of the divider's total. Switching aggressors push charge in; quiet ones act as ground. If aggressor A (1.5 fF) rises while B (2 fF) stays quiet on a victim with 7 fF to ground, $V_p = 1.5/(1.5 + 2 + 7) = 142.9$ mV, lower than if B were not there, because B's coupling loads the victim. If both rise together, the switching coupling is 3.5 fF and $V_p = 3.5/10.5$ of the swing. Glitch analysis therefore needs to know which aggressors can switch in the same direction at the same time, which comes from timing windows and logic; assuming all aggressors align is pessimistic but safe.

## Shielding

Inserting a grounded shield wire between victim and aggressor converts much of the coupling into capacitance to the shield, which is ground. The victim's total capacitance stays about the same, but the switching coupling drops. For a victim with 7.2 fF to ground and 1.6 fF of coupling to a 1 V aggressor, the bound is 181.8 mV; if a shield converts 60% of the coupling, the remaining 0.64 fF gives 72.7 mV, a 109 mV improvement. Shielding costs a routing track and adds ground capacitance, which slows the victim slightly. That is why it is reserved for clocks, asynchronous resets and other nets where a glitch is fatal, while ordinary nets are fixed with spacing or by upsizing the victim's holding driver.

## Going deeper

Shepard and Narayanan, "Noise in deep submicron digital design", ICCAD (1996); Devgan, "Efficient coupled noise estimation for on-chip interconnects", ICCAD (1997).
