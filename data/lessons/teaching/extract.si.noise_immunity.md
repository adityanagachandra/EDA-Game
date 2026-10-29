# Glitch peak with a holding driver and noise immunity

A real victim is not floating: its driver holds it at its rail through a finite resistance, and that resistance drains injected charge during the aggressor's edge. The glitch peak then depends on the ratio of the aggressor's rise time to the victim's time constant. Comparing that peak with what the receiving gate can tolerate, its noise immunity, is how glitch signoff decides whether a glitch is a violation.

## The holding-resistance model

Model the victim's driver as a resistance $R$ to its rail, with $C_g$ to ground and $C_c$ to an aggressor ramping from 0 to $V_{DD}$ in $t_r$. During the ramp the coupling injects a constant current $C_c V_{DD}/t_r$, and the victim obeys $(C_c + C_g)\, dv/dt + v/R = C_c V_{DD}/t_r$. Its solution rises toward $R C_c V_{DD}/t_r$ with time constant $\tau = R(C_c + C_g)$, and the peak, at the end of the ramp, is $V_p = V_{DD}\frac{R C_c}{t_r}(1 - e^{-t_r/\tau})$. For $R = 3.5$ kohm, $C_g = 2.1$ fF, $C_c = 0.8$ fF and a 0.9 V, 55 ps ramp, $\tau = 10.15$ ps and $V_p = 900 \times 0.0509 \times 0.996 = 45.6$ mV. After the ramp the victim decays back with time constant $\tau$.

## Noise margin

A glitch matters only if the receiving gate responds to it. The simplest test compares the peak with a flat threshold, a fraction of $V_{DD}$ derived from the receiver's DC noise margin, and reports the margin, threshold minus peak. For a 140.4 mV peak against a 25% threshold at 0.8 V (200 mV), the margin is 59.6 mV, a pass. Real noise-immunity curves depend on the glitch width: a pulse much shorter than the receiver's own delay is filtered and can exceed the DC threshold safely, while a wide pulse must stay below it. Glitches also propagate: a glitch that partly passes a gate can combine with others downstream, so noise analysis checks propagated noise at sequential and dynamic inputs, where a latched glitch becomes a functional failure.

## Rise time versus the bound

The ratio of the holding-driver peak to the floating bound $V_{DD} C_c/(C_c + C_g)$ is $(\tau/t_r)(1 - e^{-t_r/\tau})$, which depends only on $t_r/\tau$. For a fast aggressor ($t_r \ll \tau$) it approaches 1: the driver has no time to react. For a slow one ($t_r \gg \tau$) it approaches $\tau/t_r$: the driver bleeds charge as fast as it arrives. With $R = 3$ kohm, 1.9 fF to ground, 1.2 fF of coupling and a 20 ps edge, $\tau = 9.3$ ps and the peak is 41% of the 290 mV bound. The model tells a tool which fix works: a stronger holding driver lowers $\tau$ and the peak, slowing the aggressor lowers the peak, and reducing coupling lowers both the bound and the peak.

## Going deeper

Devgan, "Efficient coupled noise estimation for on-chip interconnects", ICCAD (1997); Shepard, Narayanan and Rosenberg, "Harmony: static noise analysis of deep submicron digital integrated circuits", IEEE Transactions on CAD (1999).
