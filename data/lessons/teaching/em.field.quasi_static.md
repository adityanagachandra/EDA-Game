# Quasi-static versus full-wave regimes

Lumped circuits, RC extraction and most on-chip models assume the quasi-static regime: fields change slowly compared with the time light takes to cross the structure. The test compares size with wavelength. A tool must know where the regime ends, because beyond it lumped models miss propagation, resonance and radiation. Lossy materials add a second time scale, the charge relaxation time.

## Wavelength in a dielectric

An electromagnetic wave in a medium with relative permittivity $\varepsilon_r$ (and $\mu_r = 1$) travels at $v = c/\sqrt{\varepsilon_r}$, so its wavelength at frequency $f$ is
$$\lambda = \frac{c}{f\sqrt{\varepsilon_r}}.$$
With $c = 3.00 \times 10^8$ m/s, the free-space wavelength is 300 mm at 1 GHz and 30 mm at 10 GHz.

Example: at 10 GHz in silicon dioxide, $\varepsilon_r = 3.9$, $\lambda = 30/1.975 = 15.2$ mm. In FR-4 with $\varepsilon_r = 4.3$ it is 14.5 mm. On-chip wires of a few millimetres are therefore small compared with a wavelength up to tens of gigahertz, while board traces and package planes of centimetres are not. Note that it is the dielectric's permittivity, not the conductor's, that sets the wavelength of a wave guided between conductors.

## The size-to-wavelength test

A common working rule treats a structure as quasi-static when its largest dimension $\ell$ is below about $\lambda/10$: the phase change across it is then under $36^\circ$ and lumped L, C and R capture its behavior. The threshold is a convention, chosen by the accuracy the team needs; some use $\lambda/20$.

Example: in FR-4 at 5 GHz, $\lambda = 300/(5 \cdot 2.074) = 28.9$ mm, so the lumped limit is 2.9 mm. Turned around, a 10 mm package trace in FR-4 stays lumped only up to $c/(10\ell\sqrt{\varepsilon_r}) = 1.45$ GHz. The rule should be applied at the highest significant frequency of the signal, the knee frequency, not at the clock frequency. When it fails, the structure becomes a transmission line or needs a full-wave model.

## Charge relaxation in a lossy substrate

A conducting medium neutralizes an excess charge in the charge relaxation time $\tau = \varepsilon/\sigma = \varepsilon_0\varepsilon_r\rho$. Below the frequency $1/(2\pi\tau)$ conduction current dominates and the medium acts as a conductor; above it displacement current dominates and it acts as a lossy dielectric.

Example: silicon with $\varepsilon_r = 11.7$ and resistivity 10 ohm cm (0.1 ohm m) has $\tau = 8.854 \times 10^{-12} \cdot 11.7 \cdot 0.1 = 10.4$ ps, so the crossover is about 15 GHz. A 1 ohm cm substrate moves it to about 150 GHz; a 1000 ohm cm high-resistivity substrate down to 150 MHz. That is why substrate coupling and on-chip inductor losses depend strongly on doping, and why substrate models need both resistance and capacitance between ports: neither a pure conductor nor a pure dielectric is right across the band.

## Going deeper

Haus and Melcher, "Electromagnetic Fields and Energy", chapters 3 and 7; Hall and Heck, "Advanced Signal Integrity for High-Speed Digital Designs", chapter 2.
