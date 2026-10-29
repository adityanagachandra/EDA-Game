# The skin effect

At high frequency, current in a conductor crowds toward its surface, within about one skin depth. The effective cross-section shrinks, so AC resistance rises with the square root of frequency. Skin effect sets the conductor loss of every package trace and board channel and the frequency-dependent resistance of on-chip clock and power wires, so extraction and channel models must include it.

## Skin depth

Inside a good conductor, a field at angular frequency $\omega$ decays with depth as $e^{-x/\delta}$, with skin depth
$$\delta = \sqrt{\frac{2}{\omega\mu\sigma}} = \sqrt{\frac{\rho}{\pi f\mu}}.$$
Current density follows the same decay, so most current flows within a few skin depths of the surface.

Example: copper with $\rho = 1.72 \times 10^{-8}$ ohm m and $\mu = \mu_0$ has $\delta = 2.09$ um at 1 GHz, 6.6 um at 100 MHz and 0.66 um at 10 GHz; at 60 Hz it is about 8.5 mm. Skin depth falls as $1/\sqrt f$. On-chip wires are a few hundred nanometres thick, so skin effect matters on chip only for thick top-metal lines at many gigahertz, while package and board traces of tens of micrometres feel it from a few hundred megahertz.

## AC resistance of a round wire

When $\delta$ is much smaller than the conductor's dimensions, current flows in a shell about $\delta$ thick at the surface. For a round wire of radius $a \gg \delta$, the AC resistance per unit length is about $\rho/(2\pi a\delta)$, against $R_{dc} = \rho/(\pi a^2)$, so
$$\frac{R_{ac}}{R_{dc}} \approx \frac{a}{2\delta}.$$

Example: a copper bond wire of radius 12.5 um at 1 GHz has $\delta = 2.09$ um and $R_{ac}/R_{dc} \approx 3.0$; a 50 um radius wire has a ratio of 12. Flat traces behave similarly, with current concentrating on the surfaces closest to the return path. The approximation fails when $\delta$ is comparable to the conductor size, where the exact Bessel-function solution or a filament model is needed; extractors subdivide conductors into filaments to compute frequency-dependent R and L.

## Square-root frequency scaling

Deep in the skin-effect regime, $R_{ac} \propto 1/\delta \propto \sqrt f$. So a resistance measured or simulated at one frequency scales to another as
$$R_{ac}(f_2) = R_{ac}(f_1)\sqrt{f_2/f_1}.$$
Quadrupling frequency doubles the resistance.

Example: a trace with 1.5 ohm of AC resistance at 2 GHz has 3 ohm at 8 GHz and about 4.7 ohm at 20 GHz. Channel loss from conductors therefore grows as $\sqrt f$ in dB, while dielectric loss grows as $f$, which is why dielectric loss takes over at high frequency. Real copper surfaces are rough, and roughness raises loss further once $\delta$ falls to the roughness scale; channel models add a roughness correction for this. Causal models must also include the matching frequency-dependent internal inductance, which falls as $1/\sqrt f$.

## Going deeper

Ramo, Whinnery and Van Duzer, "Fields and Waves in Communication Electronics" (3rd ed.), chapter 5; Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs", chapter 5.
