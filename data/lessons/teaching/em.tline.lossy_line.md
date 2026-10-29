# Attenuation on lossy lines

Real lines lose energy in conductor resistance and dielectric conductance. Attenuation grows with length and frequency, rounds edges, and closes the eye of high-speed links. Channel tools express loss in decibels, track how each mechanism scales with frequency, and budget it against what an equalizer can recover. The low-loss approximation makes those budgets easy to compute.

## The attenuation constant

With losses, a wave decays as $e^{-\alpha x}$ along the line. When $R' \ll \omega L'$ and $G' \ll \omega C'$, the low-loss approximation gives
$$\alpha \approx \frac{R'}{2Z_0} + \frac{G'Z_0}{2}$$
in nepers per unit length: a conductor term and a dielectric term. The characteristic impedance and velocity stay close to their lossless values.

Example: a 50 ohm line with $R' = 10$ ohm/m and negligible $G'$ has $\alpha = 0.1$ Np/m; over 0.5 m that is 0.05 Np or 0.434 dB. Raising $Z_0$ reduces the conductor term for a given $R'$ but raises the dielectric term, one reason 50 ohm is a reasonable compromise for many channels. $R'$ itself rises as $\sqrt f$ from the skin effect, and $G' = \omega C'\tan\delta$ rises as $f$.

## Nepers and decibels

Amplitude after a length $\ell$ is $V_{out} = V_{in}e^{-\alpha\ell}$; $\alpha\ell$ is the loss in nepers. Since $20\log_{10}e = 8.686$, one neper is 8.686 dB. Power decays as $e^{-2\alpha\ell}$, and in decibels the voltage and power losses are the same number.

Example: $\alpha = 1.5$ Np/m over 0.2 m gives $\alpha\ell = 0.3$ Np, so a 400 mV input leaves as $400e^{-0.3} = 296$ mV, a loss of 2.6 dB. Losses in dB add along a cascade of lines, connectors and vias, which is why channel budgets are written in dB at the Nyquist frequency. Mixing nepers and dB, or voltage and power dB, is a common source of factor-of-two errors in tool outputs.

## Scaling loss with frequency

Conductor loss in the skin-effect regime grows as $\sqrt f$ and dielectric loss as $f$. Given each part at a reference frequency $f_0$, the total at $f_1$ is
$$L(f_1) = L_c(f_0)\sqrt{f_1/f_0} + L_d(f_0)\,(f_1/f_0).$$

Example: a channel with 2 dB conductor loss and 1 dB dielectric loss at 4 GHz has, at 16 GHz, $2 \cdot 2 + 1 \cdot 4 = 8$ dB. The dielectric part grew from a third to half of the total. Doubling a link's data rate doubles its Nyquist frequency and so raises loss by more than $\sqrt 2$ but less than 2 times in dB. Channel tools also add surface-roughness and causal dielectric models, since a constant loss tangent with constant permittivity violates causality; these refinements change numbers but not the basic scaling.

## Going deeper

Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs", chapters 5 and 6; Pozar, "Microwave Engineering" (4th ed.), chapter 2.
