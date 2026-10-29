# Drift current and mobility

Drift is current driven by an electric field, and mobility is the proportionality constant between field and carrier velocity. Mobility sets the resistance of every doped region and the drive current of every transistor; it is the parameter compact models spend the most effort describing. A tool designer meets it in sheet-resistance values for extraction, in $\mu C_{ox}$ for transistor sizing, and in the scattering models of TCAD.

## Velocity proportional to field

In a field $E$ a carrier accelerates, scatters off phonons, impurities and interfaces, and reaches an average drift velocity. At low fields this velocity is proportional to the field:
$$v_d = \mu E.$$
The mobility $\mu$ is in cm$^2$/(V s). Bulk silicon at low doping has electron mobility around 1400 and hole mobility around 450 cm$^2$/(V s); inversion layers in MOSFETs are lower because of interface scattering.

Example: $\mu_n = 600$ cm$^2$/(V s) and $E = 2000$ V/cm give $v_d = 1.2 \times 10^6$ cm/s.

The linear relation breaks down at high fields. Above roughly $10^4$ V/cm in silicon, velocity approaches a saturation value near $10^7$ cm/s. A 30 nm channel with 0.5 V across it averages $1.7 \times 10^5$ V/cm, so modern transistors operate far from the low-field regime, which is why velocity saturation enters every short-channel model.

## Conductivity and resistivity

The drift current density is charge density times velocity: $J = q n v_d = q n \mu E$ for electrons, plus the hole term $q p \mu_p E$. The conductivity is therefore
$$\sigma = q(n\mu_n + p\mu_p), \qquad \rho = 1/\sigma.$$
In doped material the minority term is negligible.

Example: n-type silicon with $n = 10^{16}$ cm$^{-3}$ and $\mu_n = 1200$ cm$^2$/(V s): $\sigma = 1.602 \times 10^{-19} \times 10^{16} \times 1200 = 1.92$ S/cm and $\rho = 0.52$ ohm cm.

For a layer of thickness $t$ the sheet resistance is $R_s = \rho/t$ in ohms per square. Extraction tools and parasitic extractors use sheet resistance for diffusion, well and polysilicon resistors: a resistor of $N$ squares has resistance $N R_s$. Because mobility itself falls as doping rises, resistivity falls more slowly than $1/N_D$.

## Combining scattering mechanisms

Each scattering mechanism acts independently and contributes its own scattering rate. Since mobility is inversely proportional to the scattering rate, the rates add as reciprocals of mobility. This is Matthiessen's rule:
$$\frac{1}{\mu} = \frac{1}{\mu_{lattice}} + \frac{1}{\mu_{impurity}} + \frac{1}{\mu_{surface}} + \cdots$$

Example: lattice scattering alone gives 1400 cm$^2$/(V s) and impurity scattering alone gives 800 cm$^2$/(V s). Then $1/\mu = 1/1400 + 1/800$, so $\mu = 509$ cm$^2$/(V s), lower than either and closest to the smaller.

Lattice scattering grows with temperature (roughly $\mu \propto T^{-3/2}$ in lightly doped silicon), so transistors slow down when hot at high supply voltage. In MOSFET inversion layers a third term, surface-roughness scattering, grows with the vertical field. Compact models combine these through effective-mobility expressions with fitted coefficients, and TCAD uses calibrated mobility models of the same structure.

## Going deeper

Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), section 1.5; Lundstrom, "Fundamentals of Carrier Transport" (2nd ed.), chapter 3.
