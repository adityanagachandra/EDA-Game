# Photonic ICs and electronic-photonic design

Photonic integrated circuits guide light in waveguides on silicon or other platforms, with modulators, couplers, ring resonators and photodetectors as their components. Electronic-photonic design automation co-designs them with the electronic drivers and receivers next to them. The engineering quantities are loss budgets in dB, resonance and free spectral range, and sensitivity to fabrication and temperature, which make photonic design closer to RF and analog than to digital.

## Link loss budgets

An optical link is a cascade of losses, so its budget is a sum in dB. Laser power in dBm minus coupler losses, waveguide propagation loss (dB/cm times length), modulator insertion loss and splitter or multiplexer losses gives the power at the photodetector. Launching 10 dBm through two grating couplers of 3 dB each, 1.5 cm of waveguide at 2 dB/cm and a 5 dB modulator leaves $10 - 6 - 3 - 5 = -4$ dBm. The receiver's sensitivity, set by photodiode responsivity and amplifier noise, must be below this with margin, much as an RF sensitivity budget works. A photonic design tool tracks these budgets per path and checks them across process and temperature, since coupler and waveguide losses vary.

## Ring resonators

A ring resonator couples to a bus waveguide and resonates when its round-trip length is a whole number of wavelengths in the guide. Neighbouring resonances are separated by the free spectral range $\mathrm{FSR} \approx \lambda^2/(n_g L)$, with group index $n_g$ and circumference $L$. A ring of radius 10 micrometres ($L = 62.8\ \mu$m) with $n_g = 4.2$ near 1550 nm has $\mathrm{FSR} = 1550^2/(4.2 \cdot 62832)\ \mathrm{nm} = 9.10$ nm. The resonance width sets the quality factor $Q = \lambda_0/\Delta\lambda$: a 50 pm linewidth at 1550 nm gives $Q = 31{,}000$. Rings are compact modulators, filters and wavelength multiplexers, but a high Q makes them very sensitive: nanometre-scale width variation or a few kelvin of temperature change shifts the resonance by a significant fraction of its width.

## What electronic-photonic design automation must do

Photonic layout is not Manhattan: waveguides are curves whose bend radius, width and spacing set loss and coupling, so layout tools generate curvilinear geometry from parameterized components and must check it against photonic design rules. Component behaviour is simulated by electromagnetic solvers (FDTD, eigenmode expansion) and abstracted into compact models, often S-parameter matrices per wavelength, which circuit-level photonic simulators then cascade. Co-design with electronics matters because modulator drivers, transimpedance amplifiers and thermal tuning heaters with their control loops determine link performance, so the flow needs mixed electronic-photonic simulation and joint layout and verification. Variation-aware analysis is essential: resonant devices are routinely trimmed or thermally tuned, and the tuning power is part of the design's energy budget.

## Going deeper

Chrostowski and Hochberg, "Silicon Photonics Design: From Devices to Systems" (2015).
