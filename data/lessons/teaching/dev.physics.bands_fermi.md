# Energy bands and the Fermi level

Every device model an EDA tool uses, from a TCAD mesh solve to a BSIM card, rests on one picture: electrons in a crystal occupy bands of allowed energies, separated by a forbidden gap, and the Fermi level says how full those bands are. Threshold voltages, junction potentials and leakage currents are all statements about where the Fermi level sits relative to the band edges. A tool designer who can read a band diagram can read what a device simulator is actually computing.

## Bands, gaps and carriers

In a crystal the discrete levels of isolated atoms broaden into bands. In silicon at 0 K the valence band is completely full and the conduction band completely empty, separated by a band gap $E_g \approx 1.12$ eV at room temperature. A metal, by contrast, has a partly filled band, so it conducts at any temperature. An insulator such as SiO$_2$ has a gap near 9 eV, too wide for thermal excitation.

At room temperature a few electrons are thermally excited across silicon's gap. Each leaves behind an empty valence state, a hole, which moves like a positive carrier. Near a band edge the energy is close to parabolic, $E \approx \hbar^2k^2/(2m^*)$, and the curvature defines the effective mass $m^*$. That one number carries the lattice's influence into mobility, density of states and tunneling. Silicon's gap is indirect, so radiative recombination is weak, which is why silicon is a poor light emitter.

## The Fermi-Dirac occupancy

The probability that an allowed state at energy $E$ holds an electron in equilibrium is the Fermi-Dirac function
$$f(E) = \frac{1}{1 + e^{(E - E_F)/kT}}.$$
At $E = E_F$ the occupancy is exactly one half. The function falls from 1 to 0 over a few $kT$ around $E_F$, and $kT$ is 25.85 meV at 300 K.

Example: at 300 K a state 50 meV above $E_F$ has $(E - E_F)/kT = 50/25.85 = 1.934$, so $f = 1/(1 + e^{1.934}) = 0.126$. The probability that the same state is empty, $1 - f = 0.874$, is the hole occupancy. When $E - E_F$ is several $kT$, $f \approx e^{-(E - E_F)/kT}$: this Boltzmann approximation is what makes carrier densities exponential in potential, and what device solvers exploit. Note the symmetry: a state the same distance below $E_F$ is occupied with probability $1 - f$, so the occupancy curve is antisymmetric about one half. Raising the temperature widens the transition region but leaves the half-filled point at $E_F$.

## What the Fermi level does in a device

In thermal equilibrium the Fermi level is flat across an entire device: a single $E_F$ for the source, channel, drain and body. Any gradient in $E_F$ would drive a net current, which equilibrium forbids. The bands bend instead, and that bending is the built-in potential of a junction or the surface potential of a MOS capacitor.

Under bias the electron and hole populations each have their own quasi-Fermi level. Their separation across a forward-biased junction is $qV$, and their gradients drive the electron and hole currents. TCAD reports these quasi-Fermi potentials directly.

In undoped silicon $E_F$ sits near midgap, offset by $(kT/2)\ln(N_v/N_c)$. In heavily doped source and drain regions, around $10^{20}$ cm$^{-3}$, $E_F$ moves into the band. Then the Boltzmann approximation fails, Fermi-Dirac statistics are needed, and bandgap narrowing appears, which is why simulators expose a degenerate-statistics switch.

## Going deeper

Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), chapter 1; Pierret, "Semiconductor Device Fundamentals", chapters 2 and 4.
