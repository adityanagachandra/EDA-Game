# Doping and the mass-action law

Doping is how a process turns silicon into a device: donors such as phosphorus and arsenic add electrons, acceptors such as boron add holes. The mass-action law, $np = n_i^2$ in equilibrium, ties the two carrier populations together, so knowing the doping fixes both. Every TCAD initial guess, every junction model and every minority-carrier leakage estimate starts from this one relation.

## Donors and minority holes

A donor atom has one more valence electron than silicon. At room temperature essentially all donors are ionized, so in n-type material with donor density $N_D \gg n_i$ the electron density is $n \approx N_D$. The product of the two Boltzmann expressions for $n$ and $p$ does not depend on where the Fermi level is, so in equilibrium
$$np = n_i^2$$
holds regardless of doping. The minority hole density is therefore $p = n_i^2/N_D$.

Example: with $n_i = 10^{10}$ cm$^{-3}$ and $N_D = 10^{17}$ cm$^{-3}$, $p = 10^{20}/10^{17} = 1000$ cm$^{-3}$. Doping raised the electron density by seven orders of magnitude and suppressed holes by the same factor. With the older value $n_i = 1.5 \times 10^{10}$ the answer is 2250 cm$^{-3}$, so always use the $n_i$ stated for the problem.

## Acceptors and minority electrons

An acceptor atom such as boron has one fewer valence electron, so it captures an electron from the valence band and leaves a hole. In p-type material with $N_A \gg n_i$, the hole density is $p \approx N_A$ and the minority electron density is $n = n_i^2/N_A$.

Example: $N_A = 4 \times 10^{15}$ cm$^{-3}$ and $n_i = 10^{10}$ cm$^{-3}$ give $n = 10^{20}/(4 \times 10^{15}) = 25{,}000$ cm$^{-3}$.

Minority carriers are what an NMOS inverts: a p-type body becomes an electron channel when the gate pulls the surface electron density above the bulk hole density. Minority densities also set the diffusion component of junction leakage, which scales as $n_i^2/N$. A lightly doped side of a junction therefore dominates its saturation current, a fact compact models encode through junction parameters.

## Compensated doping

Real device regions often contain both donors and acceptors: a source implant on top of a well, or a halo pocket inside a channel. The net charge decides the type. If $N_D > N_A$ the region is n-type with $n \approx N_D - N_A$ (when this difference is much larger than $n_i$), and the holes follow from mass action, $p = n_i^2/(N_D - N_A)$.

Example: $N_D = 4 \times 10^{17}$, $N_A = 10^{17}$ cm$^{-3}$ and $n_i = 10^{10}$ cm$^{-3}$. The net is $3 \times 10^{17}$, so $p = 10^{20}/(3 \times 10^{17}) = 333$ cm$^{-3}$.

Compensation does not cancel the dopants physically. Both kinds of ionized atom remain and both scatter carriers, so a compensated region has lower mobility than one with the same net doping from a single species. TCAD mobility models take total ionized impurity density as an input for exactly this reason.

## Going deeper

Pierret, "Semiconductor Device Fundamentals", chapter 2; Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), section 2.1.
