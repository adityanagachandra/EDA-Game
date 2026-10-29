# Quantum confinement in thin bodies

When a channel is only a few nanometres thick, electrons behave as waves trapped in a box: their allowed energies split into subbands, and the lowest one sits above the bulk band edge. This quantum confinement raises the threshold, moves the charge away from the interface and makes threshold sensitive to body thickness. Every thin-body device, from FD-SOI to nanosheets, needs a compact model that captures it.

## Ground-state confinement energy

Model the thin body as an infinite square well of width $t$. The allowed energies above the band edge are
$$E_n = \frac{n^2\pi^2\hbar^2}{2m^*t^2}, \qquad n = 1, 2, \ldots$$
Using $\hbar^2/(2m_0) = 0.0381$ eV nm$^2$, $E_1 = 0.376/((m^*/m_0)\,t^2)$ eV with $t$ in nm.

Example: $t = 5$ nm and confinement mass $0.916\,m_0$ (silicon's heavy longitudinal mass, relevant for one valley pair on a (100) surface) give $E_1 = 0.376/(0.916 \times 25) = 16.4$ meV. With $m^* = 0.19\,m_0$, $E_1 = 79$ meV.

The infinite well overestimates the energy somewhat, since real barriers are finite and fields tilt the well. Still, the scaling, $E_1 \propto 1/(m^* t^2)$, is correct and is what matters for design: halving the body thickness quadruples the confinement energy.

## Subband spacing

The second level of the well sits at $E_2 = 4E_1$, so the spacing between the first two subbands is
$$E_2 - E_1 = 3E_1.$$

Example: $t = 5$ nm with $m^* = 0.19\,m_0$ gives $E_1 = 79$ meV and a spacing of 237 meV, about nine times $kT$ at room temperature. Almost all electrons then occupy the lowest subband (or lowest few, counting silicon's several valleys with different masses).

When the spacing is many $kT$, transport is effectively two-dimensional, and the valley with the lightest confinement mass in the transport direction may not be the lowest in energy. This changes effective mobility and density of states. Strong confinement also redistributes electrons among silicon's conduction valleys, which is one way crystal orientation and strain affect thin-body mobility. Device simulators handle this with Schrödinger-Poisson solvers or density-gradient corrections, and compact models fold it into effective parameters.

## Threshold and band-gap effects

Confinement raises the lowest electron state and lowers the highest hole state, so the effective band gap widens by
$$\Delta E_g = E_{1,e} + E_{1,h}.$$

Example: a 4 nm body with electron mass $0.916\,m_0$ and hole mass $0.29\,m_0$. $E_{1,e} = 0.376/(0.916 \times 16) = 25.7$ meV and $E_{1,h} = 0.376/(0.29 \times 16) = 81.0$ meV, a widening of 107 meV.

For a transistor the electron shift raises the threshold by roughly $E_{1,e}/q$, and because $E_1 \propto 1/t^2$, threshold depends on body thickness: $dE_1/dt = -2E_1/t$. At $t = 4$ nm with $E_1 = 25.7$ meV, a 0.5 nm thickness variation shifts it by about 6 meV, and much more for lighter masses or thinner bodies. Confinement also pushes the charge centroid away from the interface, adding to effective oxide thickness. Compact models for FinFETs and nanosheets include both effects.

## Going deeper

Davies, "The Physics of Low-Dimensional Semiconductors" (1998), chapter 4; Omura and others, "Quantum-mechanical effects on the threshold voltage of ultrathin-SOI nMOSFETs", IEEE Electron Device Letters (1993).
