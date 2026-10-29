# Dielectric loss and loss tangent

Real dielectrics absorb energy because their polarization lags the applied field. The loss tangent quantifies it. Dielectric loss grows linearly with frequency and dominates channel loss at high data rates, so material choice and loss models are central to signal-integrity tooling. In circuit terms it is a conductance in parallel with every capacitance.

## Complex permittivity and the loss tangent

A lossy dielectric has complex permittivity $\varepsilon = \varepsilon'(1 - j\tan\delta)$, where $\tan\delta = \varepsilon''/\varepsilon'$ is the loss tangent. A capacitor filled with it has admittance $j\omega C(1 - j\tan\delta) = \omega C\tan\delta + j\omega C$: an ideal capacitor in parallel with a conductance
$$G = \omega C\tan\delta.$$

Example: 10 pF with $\tan\delta = 0.01$ at 1 GHz gives $G = 2\pi \cdot 10^9 \cdot 10^{-11} \cdot 0.01 = 628$ uS, an equivalent parallel resistance of about 1.6 kohm. Typical loss tangents: FR-4 about 0.02, low-loss laminates 0.002 to 0.005, silicon dioxide much lower. On a transmission line the same physics appears as the shunt conductance $G'$ in the RLGC model, with $G' = \omega C'\tan\delta$.

## Dielectric attenuation on a line

For a TEM line in a homogeneous dielectric, the shunt conductance produces an attenuation constant
$$\alpha_d = \frac{\omega}{2}\sqrt{\mu\varepsilon}\tan\delta = \frac{\pi f\sqrt{\varepsilon_r}\tan\delta}{c}$$
in nepers per metre; multiply by 8.686 for dB per metre.

Example: FR-4 stripline, $\varepsilon_r = 4.3$, $\tan\delta = 0.02$, at 5 GHz: $\alpha_d = \pi \cdot 5 \times 10^9 \cdot 2.074 \cdot 0.02/3 \times 10^8 = 2.17$ Np/m, or 18.9 dB/m. Over 20 cm that is 3.8 dB. A low-loss laminate with $\tan\delta = 0.004$ cuts it to 0.75 dB. This term depends only on the material and frequency, not on trace width, so it cannot be reduced by geometry. For microstrip, only part of the field is in the dielectric, so an effective permittivity and a filling factor enter.

## Conductor versus dielectric loss

Conductor loss from the skin effect grows as $\sqrt f$; dielectric loss grows as $f$. If at a reference frequency $f_0$ the conductor contributes $a_c$ dB and the dielectric $a_d$ dB, the two are equal where $a_c\sqrt{f/f_0} = a_d f/f_0$, that is
$$f_\times = f_0\left(\frac{a_c}{a_d}\right)^2.$$

Example: at 2 GHz a channel loses 3 dB to conductors and 1.5 dB to the dielectric. The crossover is at $2 \cdot 2^2 = 8$ GHz; above it the dielectric dominates. Halving the loss tangent moves the crossover to 32 GHz. This scaling explains material roadmaps for high-speed links: wider traces fight conductor loss, low-loss laminates fight dielectric loss, and which matters depends on the Nyquist frequency of the link.

## Going deeper

Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs", chapter 5; Pozar, "Microwave Engineering" (4th ed.), chapter 2.
