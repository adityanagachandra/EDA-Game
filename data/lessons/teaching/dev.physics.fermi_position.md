# Fermi-level position from doping

Where the Fermi level sits relative to the intrinsic level and the band edges is the compact way to describe a doped region. The offset $E_F - E_i$ feeds directly into the threshold voltage (through $2\phi_F$), the built-in potential of junctions and the work-function difference between gate and channel. Converting between doping and Fermi-level position, in both directions, is a routine step in reading TCAD output and in building compact models.

## Fermi level from doping

Write the carrier densities relative to the intrinsic level $E_i$:
$$n = n_i e^{(E_F - E_i)/kT}, \qquad p = n_i e^{(E_i - E_F)/kT}.$$
In n-type material $n \approx N_D$, so $E_F - E_i = kT \ln(N_D/n_i)$. In p-type material $E_i - E_F = kT \ln(N_A/n_i)$.

Example at 300 K ($kT = 25.85$ meV, $n_i = 10^{10}$ cm$^{-3}$): $N_A = 10^{17}$ cm$^{-3}$ gives $E_i - E_F = 25.85 \times \ln(10^7) = 25.85 \times 16.118 = 416.7$ meV. The quantity $\phi_F = (E_i - E_F)/q = 0.4167$ V is the Fermi potential of an NMOS body; strong inversion begins when the surface potential reaches $2\phi_F$.

Each factor of ten in doping moves $E_F$ by $kT \ln 10 = 59.5$ meV at 300 K, a number worth remembering. For an n-type region at $10^{18}$ cm$^{-3}$ the offset is $25.85 \times \ln(10^8) = 476$ meV, already close to the conduction band, since $E_c - E_i$ is about 560 meV. The logarithm makes the Fermi position a gentle function of doping: two orders of magnitude in doping move it by only about 120 meV, while the carrier densities themselves change a hundredfold.

## Doping from a Fermi level

TCAD and band-diagram plots report potentials, and a model builder often has to go the other way: from a Fermi-level offset back to a carrier density. Invert the same relation: $n = n_i e^{(E_F - E_i)/kT}$.

Example: a simulator shows $E_F$ 470 meV above $E_i$ in an n-type region at 300 K. Then $n = 10^{10} \times e^{470/25.85} = 10^{10} \times e^{18.18} = 7.9 \times 10^{17}$ cm$^{-3}$.

Because the dependence is exponential, small errors in potential become large errors in density: 10 meV is a factor of $e^{10/25.85} = 1.47$. That is why device solvers often use potentials or quasi-Fermi potentials as unknowns. It also explains why a threshold extraction that is off by 20 mV implies an apparent doping error of more than a factor of two, if the error is attributed to doping.

## Distance from the band edge

The intrinsic level is a convenient reference, but contact and degeneracy questions are about distance to a band edge. Using the effective density of states, $n = N_c e^{-(E_c - E_F)/kT}$, so
$$E_c - E_F = kT \ln(N_c/N_D).$$

Example: $N_D = 10^{17}$ cm$^{-3}$ and $N_c = 2.8 \times 10^{19}$ cm$^{-3}$ at 300 K give $E_c - E_F = 25.85 \times \ln(280) = 145.7$ meV.

As $N_D$ approaches $N_c$, the offset approaches zero and then would become negative: the Fermi level enters the band. The Boltzmann formula then no longer applies, and the region is degenerate. Source and drain regions above about $10^{19}$ cm$^{-3}$ are in this regime, so simulators switch to Fermi-Dirac statistics and include bandgap narrowing. A compact-model designer meets this through the contact and series-resistance terms rather than through the channel equations.

## Going deeper

Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), section 2.1; Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), section 1.4.
