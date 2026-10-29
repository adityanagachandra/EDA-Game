# Noise analysis

Noise analysis is a small-signal frequency sweep with a twist: instead of one input, every resistor and transistor contributes a random noise source, and the simulator sums their powers at the output. Because the sources are uncorrelated, power densities add, not amplitudes. Simulators compute each source's transfer to the output with one adjoint solve per frequency, then report output noise and the equivalent noise referred to the input.

## Thermal noise of a resistor

A resistor at temperature $T$ has thermal noise modeled as a series voltage source with one-sided power spectral density $4kTR$ in V$^2$/Hz, or equivalently a parallel current source with density $4kT/R$. At 300 K, $4kT = 1.657 \times 10^{-20}$ J. For 1 k$\Omega$, $4kTR = 1.657 \times 10^{-17}$ V$^2$/Hz, and the density is its square root, 4.07 nV/$\sqrt{\text{Hz}}$. It scales as $\sqrt{R}$: 100 k$\Omega$ gives 40.7 nV/$\sqrt{\text{Hz}}$. Transistors add channel thermal noise, proportional to $g_m$, and flicker noise, whose density grows as $1/f$; the compact model supplies each density at the operating point.

## Summing independent noise at the output

Each source reaches the output through a transfer function $H_i(f)$, and uncorrelated sources add in power:
$$S_{out}(f) = \sum_i |H_i(f)|^2 S_i(f).$$
For a resistor contributing 4 nV/$\sqrt{\text{Hz}}$ through a gain of 10 and a transistor contributing 3 nV/$\sqrt{\text{Hz}}$ through a gain of 10, the output density is $\sqrt{40^2 + 30^2} = 50$ nV/$\sqrt{\text{Hz}}$, not 70. A simulator obtains all $H_i$ at once with an adjoint solve: one solve of $(G + j\omega C)^T \lambda = e_{out}$ per frequency gives the transfer from every possible source location. Dividing the output density by the signal gain gives the input-referred noise.

## Integrated kT/C noise

Integrating the output density over frequency gives the total rms noise. For a resistor across a capacitor, the density $4kTR$ shaped by the RC low-pass integrates to a variance of exactly $kT/C$, independent of $R$: a larger resistor is noisier but filters more, and the two effects cancel. At 300 K, 1 pF gives $\sqrt{kT/C} = \sqrt{4.14 \times 10^{-9}} = 64.4$ uV rms. This sets the noise floor of every sampling circuit and switched-capacitor stage, and it explains why sampling capacitors in high-resolution converters must be large. A noise analysis reproduces it only if its sweep covers the bandwidth well beyond the RC corner.

## Going deeper

Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., McGraw-Hill, 2017), chapter 7; Rohrer, Nagel, Meyer and Weber, "Computationally efficient electronic-circuit noise calculations" (IEEE JSSC, 1971).
