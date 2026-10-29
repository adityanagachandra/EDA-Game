# Intrinsic carrier concentration

The intrinsic carrier concentration $n_i$ is the electron density of pure silicon in equilibrium. It sets the scale of every minority-carrier density, every junction leakage current and every built-in potential, and it rises steeply with temperature. Leakage models in timing and power libraries inherit that steep temperature dependence from $n_i$, so a tool designer should be able to compute it and predict how it moves.

## Intrinsic density from the band gap

In intrinsic material the electron and hole densities are equal, $n = p = n_i$. Integrating the Fermi-Dirac occupancy against the density of states in the Boltzmann limit gives
$$n = N_c e^{-(E_c - E_F)/kT}, \qquad p = N_v e^{-(E_F - E_v)/kT},$$
where $N_c$ and $N_v$ are the effective densities of states of the conduction and valence bands. Multiplying them removes $E_F$:
$$n_i = \sqrt{N_c N_v}\, e^{-E_g/(2kT)}.$$

Example: $N_c = 2.8 \times 10^{19}$, $N_v = 1.04 \times 10^{19}$ cm$^{-3}$ and $E_g = 1.12$ eV at 300 K. Then $\sqrt{N_c N_v} = 1.71 \times 10^{19}$ and $E_g/(2kT) = 21.66$, so $n_i \approx 6.7 \times 10^{9}$ cm$^{-3}$. The commonly quoted silicon values of $1.0 \times 10^{10}$ or $1.5 \times 10^{10}$ cm$^{-3}$ come from different choices of $N_c$, $N_v$ and gap, which is why a problem must state the value it uses.

## Intrinsic density at another temperature

Both effective densities of states scale as $T^{3/2}$, and $kT$ in the exponent grows linearly with $T$. To compute $n_i$ at a new temperature, rescale $N_c$ and $N_v$ by $(T/300)^{3/2}$ and recompute the exponential with $kT = 8.617 \times 10^{-5}\,T$ eV. The band gap also shrinks slightly with temperature (about 0.05 eV between 300 K and 400 K in silicon), and careful models include it; simple estimates often hold it fixed.

Example with a fixed 1.12 eV gap: at 400 K, $kT = 0.03447$ eV and $E_g/(2kT) = 16.25$. The prefactor grows by $(400/300)^{3/2} = 1.54$. Starting from $\sqrt{N_c N_v} = 1.71 \times 10^{19}$ at 300 K, $n_i(400) = 1.71 \times 10^{19} \times 1.54 \times e^{-16.25} \approx 2.3 \times 10^{12}$ cm$^{-3}$, more than 300 times the room-temperature value. Going down in temperature works the same way: at 250 K the exponent grows and $n_i$ falls by roughly two orders of magnitude, which is why cold-corner junction leakage is negligible.

## How fast the intrinsic density grows

For leakage modeling the ratio matters more than the absolute value. Dividing two evaluations of the formula,
$$\frac{n_i(T_2)}{n_i(T_1)} = \left(\frac{T_2}{T_1}\right)^{3/2} \exp\left[\frac{E_g}{2k}\left(\frac{1}{T_1} - \frac{1}{T_2}\right)\right].$$
The exponential term dominates. From 300 K to 350 K with $E_g = 1.12$ eV, the exponent is $6498 \times (1/300 - 1/350) = 3.094$, giving $e^{3.094} = 22.1$, times $(350/300)^{3/2} = 1.26$, so $n_i$ grows by about 28.

Junction diffusion leakage scales as $n_i^2$ and so as $e^{-E_g/kT}$, roughly doubling every 7 to 10 kelvin near room temperature. Generation current in depletion regions scales as $n_i$. This is why the leakage tables in a Liberty file are characterized per temperature corner, and why a hot corner can show orders of magnitude more junction leakage than a cold one.

## Going deeper

Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), section 1.4; Green, "Intrinsic concentration, effective densities of states, and effective mass in silicon", Journal of Applied Physics (1990).
