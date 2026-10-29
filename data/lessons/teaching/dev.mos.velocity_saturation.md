# Velocity saturation

In a short channel the lateral field is so strong that carrier velocity stops growing with field and saturates near $10^7$ cm/s. The drain current then grows roughly linearly with overdrive instead of quadratically. This changes how delay depends on supply voltage, which is why timing tools and power-performance models use the alpha-power law instead of the square law, and why compact models carry velocity-saturation parameters.

## The critical field

At low field, $v = \mu E$. At high field carriers lose energy to optical phonons as fast as they gain it, and velocity saturates at $v_{sat} \approx 10^7$ cm/s for electrons in silicon. A simple model uses a critical field
$$E_c = \frac{v_{sat}}{\mu}$$
where the two regimes meet. Different texts include a factor of 2 in this definition, so a model card's convention must be checked.

Example: effective channel mobility 250 cm$^2$/(V s) gives $E_c = 10^7/250 = 4 \times 10^4$ V/cm $= 4$ V/µm. Across a 30 nm channel, an average field of $E_c$ needs only $4 \times 0.03 = 0.12$ V.

So in any modern channel the carriers near the drain travel at or near saturated velocity for most of the operating range. Long-channel equations, which assume $v = \mu E$ everywhere, overestimate the current.

## Current in the saturated-velocity limit

If carriers move at $v_{sat}$ everywhere, the current per width is charge times velocity:
$$\frac{I_D}{W} = C_{ox}v_{sat}(V_{GS} - V_t).$$
It is linear in overdrive and independent of channel length.

Example: EOT $= 1$ nm ($C_{ox} = 3.45 \times 10^{-6}$ F/cm$^2$), $v_{sat} = 10^7$ cm/s, overdrive 0.5 V. $I_D/W = 17.3$ A/cm $= 1727$ µA/µm. Measured devices deliver noticeably less, because real channels are only partly velocity saturated and have series resistance.

This limit explains why shrinking $L$ stopped buying proportional current. Mobility enhancement still helps, since it moves carriers to high velocity sooner. In very short channels, carriers can briefly exceed $v_{sat}$ (velocity overshoot), and the limit is better described as injection velocity at the source.

## The alpha-power law

Sakurai and Newton's alpha-power law interpolates between the square law and the velocity-saturated limit:
$$I_{D,sat} \propto (V_{GS} - V_t)^{\alpha}, \qquad 1 \le \alpha \le 2.$$
Short-channel devices have $\alpha$ around 1.2 to 1.5. Gate delay, roughly $C V_{DD}/I$, then scales as
$$t_d \propto \frac{V_{DD}}{(V_{DD} - V_t)^{\alpha}}.$$

Example: $\alpha = 1.3$, $V_t = 0.3$ V. Lowering $V_{DD}$ from 1.0 V to 0.8 V changes delay by $\frac{0.8/0.5^{1.3}}{1.0/0.7^{1.3}} = 1.24$.

This formula is how voltage-scaling studies and early power-performance models estimate the delay cost of lowering supply. Near threshold it fails, because subthreshold conduction takes over and delay grows exponentially. Timing libraries avoid the formula entirely by characterizing at each voltage, but voltage-interpolation schemes between characterized points rely on its smooth shape.

## Going deeper

Sakurai and Newton, "Alpha-power law MOSFET model and its applications to CMOS inverter delay", IEEE Journal of Solid-State Circuits (1990); Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), section 3.2.
