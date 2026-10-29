# Quantum tunneling through thin barriers

Quantum mechanics lets an electron cross a barrier it classically cannot surmount, with a probability that falls exponentially with barrier thickness. Gate leakage through thin dielectrics, band-to-band tunneling in heavily doped junctions and gate-induced drain leakage are all tunneling. Because the dependence is exponential, a fraction of a nanometre of oxide changes leakage by an order of magnitude, and every leakage model in a compact model is built around that sensitivity.

## The decay constant inside a barrier

Inside a region where the potential energy $U$ exceeds the electron's energy $E$, the Schrödinger equation has decaying solutions $\psi \propto e^{-\kappa x}$ with
$$\kappa = \frac{\sqrt{2m^*(U - E)}}{\hbar}.$$
A convenient constant is $\hbar^2/(2m_0) = 0.0381$ eV nm$^2$, so $\kappa = \sqrt{(m^*/m_0)(U - E)/0.0381}$ nm$^{-1}$.

Example: an electron in silicon facing the SiO$_2$ conduction-band barrier of about 3.1 eV, with tunneling mass $0.5\,m_0$: $\kappa = \sqrt{0.5 \times 3.1/0.0381} = 6.38$ nm$^{-1}$. The wavefunction falls by $e$ every 0.16 nm.

The decay constant grows with barrier height and effective mass. A lower barrier, as in many high-k dielectrics, gives a smaller $\kappa$, which is a cost that must be paid back with thickness. Electron and hole barriers differ too, which is why NMOS and PMOS gate leakage differ in the same stack.

## Transmission through a barrier

For a barrier of thickness $d$, the WKB approximation gives the transmission probability
$$T \approx e^{-2\int \kappa\, dx} = e^{-2\kappa d}$$
for a rectangular barrier. The factor 2 appears because probability is the square of the wavefunction amplitude.

Example: $\kappa = 6.4$ nm$^{-1}$ and $d = 1.2$ nm. Then $2\kappa d = 15.36$ and $\log_{10} T = -15.36/2.303 = -6.67$, so $T \approx 2 \times 10^{-7}$.

This tiny probability multiplies an enormous attempt flux: an inversion layer holds about $10^{13}$ electrons per square centimetre striking the interface very often. The resulting direct tunneling current reached amperes per square centimetre for SiO$_2$ near 1.2 nm, which is what ended pure SiO$_2$ scaling. Under strong bias the barrier becomes triangular and the integral changes shape (Fowler-Nordheim tunneling), but the exponential thickness dependence stays.

## Sensitivity to thickness

For two barriers that differ only in thickness, the prefactors cancel in the ratio:
$$\frac{J(d - \Delta d)}{J(d)} = e^{2\kappa\,\Delta d}.$$

Example: $\kappa = 6.4$ nm$^{-1}$. Thinning by 0.2 nm multiplies the current by $e^{2.56} = 12.9$. Thinning by one atomic layer, roughly 0.3 nm, multiplies it by about 46.

This sensitivity has three consequences for EDA. First, gate leakage is a strong function of oxide-thickness variation, so leakage corners and statistical models must include thickness. Second, compact gate-current models, such as those in BSIM4 and BSIM-CMG, carry an exponential thickness dependence with fitted parameters, and an extraction error in effective thickness shows up as a large leakage error. Third, the switch to high-k dielectrics was driven by this formula: a material with higher permittivity can be physically thicker at the same capacitance, and thickness is in the exponent.

## Going deeper

Griffiths, "Introduction to Quantum Mechanics", chapter on the WKB approximation; Lo, Buchanan, Taur and Wang, "Quantum-mechanical modeling of electron tunneling current from the inversion layer of ultra-thin-oxide nMOSFETs", IEEE Electron Device Letters (1997).
